# GoTax Company Module — Workflows & Data Flows

## Version: 1.0 | Date: 2026-07-27

---

## 1. Company Registration Workflow

```
[Admin] → Register Company form
    │
    ├─→ Validate MST format (regex)
    ├─→ Validate MSDN (optional DPI lookup)
    ├─→ Validate legal form against Enterprise Law
    │
    ├─→ [Validation OK]
    │     ├─→ Create company (ACTIVE)
    │     ├─→ Load COA template (based on regime: TT99/TT133/TT58)
    │     ├─→ Generate fiscal periods (12 months + quarterly)
    │     ├─→ Create default admin user-company assignment
    │     ├─→ Log audit (CREATE, COMPANY, MST)
    │     └─→ Return 201 + full company profile
    │
    └─→ [Validation Failed]
          └─→ Return 400/409 with specific error
```

### Data Flow: Company Creation
```
POST /api/v1/companies
  Headers: Authorization: Bearer <jwt>
  Body: {
    "legal_name_vn": "Cong ty TNHH ABC",
    "legal_name_en": "ABC Company Limited",
    "legal_form": "LLC_1MEMBER",
    "tax_code": "0123456789",
    "business_reg_no": "0101234567",
    "reg_address": "So 1, Pho A, Quan B, HN",
    "head_office_address": "So 1, Pho A, Quan B, HN",
    "phone": "+84 24 1234 5678",
    "email": "info@abc.vn",
    "legal_rep_name": "Nguyen Van A",
    "legal_rep_title": "Director",
    "chief_accountant_name": "Tran Thi B",
    "chief_accountant_email": "accountant@abc.vn",
    "accounting_regime": "TT99",
    "fiscal_year_start_month": 1,
    "currency_code": "VND",
    "tax_office_code": "01"
  }

  System Flow:
  1. Extract tenant from JWT context
  2. Validate MST unique within tenant
  3. Validate MST format: ^[0-9]{10}(-[0-9]{3})?$
  4. Validate legal_form enum
  5. BEGIN TRANSACTION
  6. INSERT INTO companies (...)
  7. INSERT INTO fiscal_years (company_id, start_month, periods)
  8. INSERT INTO periods (company_id, fiscal_year, month, status)
  9. COPY COA template WHERE regime = 'TT99' → chart_of_accounts
  10. INSERT INTO user_company_assignments (user_id, company_id)
  11. COMMIT
  12. Log: { action: "CREATE", entity: "COMPANY", entity_id: new.id, metadata: {mst} }

  Response 201:
  {
    "id": "uuid",
    "tax_code": "0123456789",
    "status": "ACTIVE",
    "accounting_regime": "TT99",
    "fiscal_year": { "start_month": 1, "periods": [...] },
    "created_at": "2026-07-27T10:00:00Z"
  }
```

---

## 2. Multi-Company Hierarchy Workflow

```
[Parent Company: Cong ty TNHH ABC - MST: 0123456789]
    │
    ├── [Branch 1: Chi nhanh HN - MST: 0123456789-001]
    │     ├── Department: Sales
    │     ├── Department: Operations
    │     └── Employees: []
    │
    ├── [Branch 2: Chi nhanh HCM - MST: 0123456789-002]
    │     ├── Department: Sales
    │     ├── Department: Support
    │     └── Employees: []
    │
    ├── [Subsidiary: Cong ty TNHH ABC Logistics - MST: 0987654321]
    │     ├── (Separate legal entity, 50%+ owned)
    │     ├── Own bank accounts, e-invoice patterns
    │     ├── Fiscal year aligns with parent
    │     └── Ready for consolidation
    │
    └── [Representative Office: VP Dai dien Da Nang]
          └── (No separate legal entity, limited operations)
```

### Data Isolation: Row-Level Security
```sql
-- Every table has company_id
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id),
    entry_number INTEGER NOT NULL,
    -- ... other fields
    UNIQUE(company_id, entry_number)
);

-- RLS Policy (if using PostgreSQL RLS)
ALTER TABLE journal_entries ENABLE ROW LEVEL SECURITY;
CREATE POLICY company_isolation ON journal_entries
    USING (company_id = current_setting('app.current_company_id')::UUID);

-- Application-level: every service method filters by company_id
func (s *Service) CreateJournalEntry(ctx context.Context, companyID uuid.UUID, entry JournalEntry) error {
    // companyID extracted from JWT context, never from request body
    // Prevents cross-company injection
}
```

### Company Context Flow
```
Client                          API Gateway                     Service                    Database
  │                                 │                              │                          │
  │  POST /api/v1/journal/entries   │                              │                          │
  │  Authorization: Bearer <jwt>    │                              │                          │
  │────────────────────────────────▶│                              │                          │
  │                                 │  JWT validation              │                          │
  │                                 │  Extract: user_id,           │                          │
  │                                 │    current_company_id,       │                          │
  │                                 │    permissions               │                          │
  │                                 │                              │                          │
  │                                 │  Forward request + context   │                          │
  │                                 │─────────────────────────────▶│                          │
  │                                 │                              │  Validate user has         │
  │                                 │                              │  access to company_id      │
  │                                 │                              │──────────────────────────▶│
  │                                 │                              │  INSERT with company_id    │
  │                                 │                              │◀──────────────────────────│
  │                                 │◀─────────────────────────────│                          │
  │◀────────────────────────────────│                              │                          │
  │  201 Created                    │                              │                          │
```

---

## 3. Fiscal Year & Period Lifecycle

```
[Company Created]
    │
    ▼
[Fiscal Year Config: start_month=1, months=12]
    │
    ▼
[Period Generation]
    ├── Period-01: Jan 2026 (OPEN)
    ├── Period-02: Feb 2026 (FUTURE)
    ├── ...
    ├── Period-12: Dec 2026 (FUTURE)
    ├── Quarter-01: Q1-2026 (OPEN)
    ├── Quarter-02: Q2-2026 (FUTURE)
    ├── Quarter-03: Q3-2026 (FUTURE)
    └── Quarter-04: Q4-2026 (FUTURE)

[Monthly Close Cycle]
    Period N (OPEN) → Journal entries posted → Validate balances
    → Validate sub-ledgers → [Optional: Revaluation] → CLOSE Period N
    → AUTO OPEN Period N+1

[Annual Close Cycle]
    Period-12 CLOSE → Validate entire year → Permanently Close
    → Generate annual FS → Fiscal Year CLOSED
    → AUTO create next fiscal year
```

### Period State Machine
```
    OPEN ──────────► CLOSED ──────────► PERMANENTLY_CLOSED
     │                                      ▲
     │                                      │
     └──────────── (reopen) ────────────────┘
          (only if NOT permanently closed)

Transitions:
  OPEN → CLOSED:             All entries balanced, sub-ledgers reconciled
  CLOSED → OPEN:             Reopen (requires admin approval, audit logged)
  CLOSED → PERMANENTLY_CLOSED: Annual closure, no return
```

---

## 4. E-Invoice Registration Flow

```
[Chief Accountant]
    │
    ├─→ Enter invoice pattern (Mau so HD)
    ├─→ Enter serial (Ky hieu)
    ├─→ Select form (Kieu HD)
    │
    ├─→ [GDT Integration Active?]
    │     ├─ YES → Push registration to GDT portal
    │     │         ├─ GDT OK → Status = ACTIVE
    │     │         └─ GDT FAIL → Status = REGISTERED (pending manual)
    │     └─ NO  → Status = REGISTERED (no auto-sync)
    │
    └─→ Log audit
    └─→ Return 201
```

### E-Invoice Pattern Data Model
```json
{
  "id": "uuid",
  "company_id": "uuid",
  "pattern_code": "01GTKT0/001",
  "serial": "AA/26E",
  "form": "Hoa don GTGT",
  "type": "GTGT",
  "status": "ACTIVE",
  "gdt_registered_at": "2026-01-15T08:00:00Z",
  "description": "Hoa don GTGT ban ra"
}
```

---

## 5. Digital Signature Workflow

```
[Chief Accountant]
    │
    ├─→ [Select Type]
    │     ├─→ USB Token
    │     │     ├─→ Insert USB → Read certificate
    │     │     ├─→ Validate: issuer, subject, valid_from/to
    │     │     └─→ Register with company
    │     │
    │     └─→ Remote HSM
    │           ├─→ Select provider (ViettelCA/VNPTCA/CloudCA)
    │           ├─→ Enter API credentials
    │           ├─→ Test connection
    │           └─→ Register with company
    │
    ├─→ [Expiry Alert Set: 30 days before expiry]
    │
    └─→ Status = ACTIVE
```

---

## 6. Company Switch Context Flow

```
Client                              API                            Database
  │                                  │                                │
  │ POST /api/v1/companies/switch    │                                │
  │ { "company_id": "uuid" }         │                                │
  │─────────────────────────────────▶│                                │
  │                                  │  Verify user → company access  │
  │                                  │───────────────────────────────▶│
  │                                  │◀───────────────────────────────│
  │                                  │                                │
  │                                  │  Generate new JWT with         │
  │                                  │  updated current_company       │
  │                                  │                                │
  │◀─────────────────────────────────│                                │
  │ 200 { "token": "new.jwt.here",   │                                │
  │         "company": { ... } }     │                                │
```

---

## 7. Data Flow: Company Isolation Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    GoTax API Gateway                      │
│  ┌─────────────────────────────────────────────────────┐ │
│  │ JWT Middleware                                       │ │
│  │  • Validate token                                    │ │
│  │  • Extract: tenant_id, user_id, company_id, roles    │ │
│  │  • Inject into context                               │ │
│  └─────────────────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    Company Service                        │
│                                                          │
│  • Company CRUD (scoped by tenant)                       │
│  • Fiscal year management (scoped by company)            │
│  • Period management (scoped by company)                 │
│  • Employee management (scoped by company)               │
│  • Department management (scoped by company)             │
│  • Bank account management (scoped by company)           │
│  • E-invoice registration (scoped by company)            │
│  • Digital signature management (scoped by company)      │
│  • Integration profile management (scoped by company)    │
│                                                          │
│  ALL service methods accept company_id from context       │
│  NEVER from request body (prevents injection)            │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    Database Layer                         │
│                                                          │
│  companies          → tenant_id (multi-tenant)            │
│  fiscal_years       → company_id (FK)                    │
│  periods            → company_id (FK)                    │
│  employees          → company_id (FK)                    │
│  departments        → company_id (FK)                    │
│  bank_accounts      → company_id (FK)                    │
│  einvoice_patterns  → company_id (FK)                    │
│  digital_signatures → company_id (FK)                    │
│  integration_profiles → company_id (FK)                  │
│                                                          │
│  RLS: company_id = current_setting('app.company_id')      │
│  INDEX: (company_id, ...) on all tables                  │
└─────────────────────────────────────────────────────────┘
```

---

## 8. User Journey: Enterprise Setup

### Day 1: Admin registers holding company
```
Admin → Login → Company Registration → Enter holding company info
→ Select TT99 regime → Set fiscal year → System loads 71 L1 accounts
→ Company created, first admin user assigned
```

### Day 2: Admin registers subsidiaries
```
Admin → Add Subsidiary → Enter sub-company info
→ Link to holding as parent → System copies regime/COA template
→ Optional: different fiscal year config
```

### Day 3: Chief Accountant configures details
```
Chief Accountant → Bank accounts → Add 3 bank accounts (VCB, CTG, TCB)
→ E-invoice → Register 2 patterns (GTGT, export)
→ Digital signature → Register USB token
→ Departments → Create 5 departments
→ Employees → Import 20 employees via CSV
```

### Day 4: Integration setup
```
Chief Accountant → GDT profile → Enter HTKK credentials
→ Test connection → Success → Status = CONNECTED
→ Optional: BHXH, Customs, DVC profiles
```

### Week 2: Fiscal year ready
```
System auto-generates 12 monthly periods
Accountant opens period-01 → Ready for journal entries
```

---

## 9. Error Handling Matrix

| Scenario | HTTP Status | Error Code | User Action |
|----------|-------------|------------|-------------|
| MST duplicate | 409 | COMPANY_MST_DUPLICATE | Use correct MST |
| MST format invalid | 400 | COMPANY_MST_INVALID | Fix MST format |
| Company not found | 404 | COMPANY_NOT_FOUND | Verify company ID |
| No access to company | 403 | COMPANY_ACCESS_DENIED | Request access |
| Company inactive | 400 | COMPANY_INACTIVE | Reactivate company |
| Period already closed | 409 | PERIOD_ALREADY_CLOSED | Select open period |
| Period permanently closed | 400 | PERIOD_PERMANENTLY_CLOSED | Cannot reopen |
| Fiscal year change with entries | 400 | FISCAL_YEAR_HAS_ENTRIES | Wait for year boundary |
| E-invoice pattern duplicate | 409 | EINVOICE_PATTERN_DUPLICATE | Use different pattern |
| Certificate expired | 400 | SIGNATURE_EXPIRED | Renew certificate |
| GDT connection failed | 400 | GDT_CONNECTION_FAILED | Check credentials/network |
| Employee code duplicate | 409 | EMPLOYEE_CODE_DUPLICATE | Use unique code |
| Department not found | 400 | DEPARTMENT_NOT_FOUND | Create department first |

---

## 10. Sequence Diagram: Annual Closing

```
Accountant          GoTax               Period-11          Period-12          GL Module
   │                  │                    │                  │                  │
   │ Close Period-11  │                    │                  │                  │
   │─────────────────▶│                    │                  │                  │
   │                  │ Validate entries   │                  │                  │
   │                  │───────────────────▶│                  │                  │
   │                  │ All balanced       │                  │                  │
   │                  │◀───────────────────│                  │                  │
   │                  │                    │                  │                  │
   │                  │ CLOSE Period-11    │                  │                  │
   │                  │───────────────────▶│                  │                  │
   │                  │ OPEN Period-12     │                  │                  │
   │                  │─────────────────────────────────────▶│                  │
   │                  │                    │                  │                  │
   │ 200 OK           │                    │                  │                  │
   │◀─────────────────│                    │                  │                  │
   │                  │                    │                  │                  │
   │ (Work continues in Period-12)         │                  │                  │
   │                  │                    │                  │                  │
   │ Close Year       │                    │                  │                  │
   │─────────────────▶│                    │                  │                  │
   │                  │ Close Period-12    │                  │                  │
   │                  │─────────────────────────────────────▶│                  │
   │                  │                    │                  │ Generate FS      │
   │                  │                    │                  │───────────────▶  │
   │                  │                    │                  │                  │
   │                  │ PERMANENT_LOCK     │                  │                  │
   │                  │─────────────────────────────────────▶│                  │
   │                  │                    │                  │                  │
   │                  │ AUTO Create FY-2027                   │                  │
   │                  │ (Periods generated, Period-01 OPEN)    │                  │
   │                  │                    │                  │                  │
   │ 200 OK + FS data │                    │                  │                  │
   │◀─────────────────│                    │                  │                  │
```
