# Business Requirements Document — Admin/Config Modules

**Document ID:** GOTAX-BRD-ADMIN-001  
**Version:** 1.0  
**Date:** 2026-08-07  
**Author:** Chief Accountant (BA Lead)  
**Status:** Draft

---

## 1. Business Context

GoTax is a Vietnamese tax-compliant General Ledger API targeting MISA SME 2026 parity. Current backend covers 17/19 core modules but admin/config depth is ~40%. Without proper system configuration, GoTax cannot be deployed to production Vietnamese enterprises.

### 1.1 Business Problem
- Companies cannot configure accounting policies per their needs
- No voucher numbering automation (manual numbering = audit risk)
- No backup/restore (data loss risk)
- No fiscal year flexibility (hardcoded = unusable for non-calendar fiscal years)
- No report customization (unprofessional output)

### 1.2 Success Criteria
- 100% of MISA SME 2026 admin/config features available via API
- All configurations persist in database (not env vars)
- Configurations are company-scoped (multi-tenant)
- Default values match Vietnamese accounting standards

---

## 2. Functional Requirements

### 2.1 System Options (FR-SYS-001)

**Priority:** P0 — Critical  
**Description:** Centralized configuration panel for all company settings

#### 2.1.1 Personal Options (Per-User)
| ID | Requirement | Priority |
|----|------------|----------|
| FR-SYS-001a | UI theme preference (light/dark) | P2 |
| FR-SYS-001b | Default search filter mode (starts-with/contains) | P1 |
| FR-SYS-001c | Search field preference (code/name/both) | P1 |
| FR-SYS-001d | Email address for notifications | P2 |

#### 2.1.2 Global Options (Per-Company)
| ID | Requirement | Priority |
|----|------------|----------|
| FR-SYS-002a | Accounting standard (Circular 99/2025 vs older) | P0 |
| FR-SYS-002b | Fiscal year start month (1-12) | P0 |
| FR-SYS-002c | Base currency (default VND) | P0 |
| FR-SYS-002d | Multi-currency enabled (yes/no) | P1 |
| FR-SYS-002e | Accounting start date | P0 |
| FR-SYS-002f | Branch-level cost allocation | P1 |
| FR-SYS-002g | Shared vs separate catalogs per branch | P1 |

#### 2.1.3 Inventory Options
| ID | Requirement | Priority |
|----|------------|----------|
| FR-SYS-003a | Costing method (FIFO/Weighted Avg/Specific ID) | P0 |
| FR-SYS-003b | Branch-level costing (shared/separate) | P1 |
| FR-SYS-003c | Negative inventory allowed | P1 |

#### 2.1.4 Number Format Options
| ID | Requirement | Priority |
|----|------------|----------|
| FR-SYS-004a | Thousands separator (period/comma/space) | P0 |
| FR-SYS-004b | Decimal separator (period/comma) | P0 |
| FR-SYS-004c | Decimal places (0-4) | P0 |
| FR-SYS-004d | Negative number display (parentheses/minus) | P1 |

#### 2.1.5 Report Options
| ID | Requirement | Priority |
|----|------------|----------|
| FR-SYS-005a | Company name on reports | P1 |
| FR-SYS-005b | Company logo on reports | P1 |
| FR-SYS-005c | Report font family | P2 |
| FR-SYS-005d | Report font size | P2 |
| FR-SYS-005e | Company info alignment (left/center/right) | P2 |
| FR-SYS-005f | Repeat company info on each page | P2 |

---

### 2.2 Voucher Numbering Rules (FR-NUM-001)

**Priority:** P0 — Critical  
**Description:** Auto-numbering for all voucher types with configurable rules

| ID | Requirement | Priority |
|----|------------|----------|
| FR-NUM-001a | Numbering rule per voucher type | P0 |
| FR-NUM-001b | Prefix (e.g., "BH" for sale, "MB" for purchase) | P0 |
| FR-NUM-001c | Suffix (optional) | P1 |
| FR-NUM-001d | Number length (default 5 digits) | P0 |
| FR-NUM-001e | Auto-increment per branch or company-wide | P0 |
| FR-NUM-001f | Reset rules (never / yearly / monthly) | P1 |
| FR-NUM-001g | Next number preview | P2 |

**Voucher Types (MISA standard):**
- BG: Bút toán ghi sổ (Journal entry)
- MB: Mua hàng (Purchase)
- BH: Bán hàng (Sale)
- TN: Thu tiền (Cash receipt)
- CN: Chi tiền (Cash payment)
- NK: Nhập kho (Goods receipt)
- XK: Xuất kho (Goods issue)
- CC: Công cụ dụng cụ (Tools)
- TSCĐ: Tài sản cố định (Fixed asset)
- Lương: Tiền lương (Payroll)

---

### 2.3 Fiscal Year Configuration (FR-FY-001)

**Priority:** P0 — Critical  
**Description:** Configurable fiscal year and accounting periods

| ID | Requirement | Priority |
|----|------------|----------|
| FR-FY-001a | Fiscal year start month (1-12) | P0 |
| FR-FY-001b | Auto-generate 12 monthly periods on FY creation | P0 |
| FR-FY-001c | Period status (OPEN/CLOSED/CLOSING) | P0 |
| FR-FY-001d | Period close requires all entries posted | P0 |
| FR-FY-001e | Period reopen with reversal | P1 |
| FR-FY-001f | Year-end close with carry-forward | P1 |

---

### 2.4 Backup & Restore (FR-BKP-001)

**Priority:** P1 — High  
**Description:** Data backup and restore capabilities

| ID | Requirement | Priority |
|----|------------|----------|
| FR-BKP-001a | Manual backup to local file | P1 |
| FR-BKP-001b | Scheduled backup (daily/weekly) | P2 |
| FR-BKP-001c | Restore from backup file | P1 |
| FR-BKP-001d | Backup includes all company data | P1 |
| FR-BKP-001e | Backup metadata (timestamp, size, version) | P2 |

---

### 2.5 Contracts (FR-CON-001)

**Priority:** P1 — High  
**Description:** Contract tracking for project-based businesses

| ID | Requirement | Priority |
|----|------------|----------|
| FR-CON-001a | Contract register (code, name, type, value) | P1 |
| FR-CON-001b | Contract types (sales, purchase, service, construction) | P1 |
| FR-CON-001c | Link to journal entries | P1 |
| FR-CON-001d | Contract status (draft/active/completed/cancelled) | P1 |
| FR-CON-001e | Start/end dates | P1 |
| FR-CON-001f | Contract value and payment schedule | P1 |
| FR-CON-001g | Cost tracking per contract | P1 |

---

### 2.6 Loan Agreements (FR-LOAN-001)

**Priority:** P1 — High  
**Description:** Loan principal and interest tracking

| ID | Requirement | Priority |
|----|------------|----------|
| FR-LOAN-001a | Loan register (code, lender, amount, rate) | P1 |
| FR-LOAN-001b | Loan types (short-term, long-term, overdraft) | P1 |
| FR-LOAN-001c | Interest calculation (simple/compound) | P1 |
| FR-LOAN-001d | Payment schedule generation | P1 |
| FR-LOAN-001e | Principal/interest breakdown per payment | P1 |
| FR-LOAN-001f | Link to GL accounts (131, 331) | P1 |
| FR-LOAN-001g | Outstanding balance tracking | P1 |

---

### 2.7 Report Customization (FR-RPT-001)

**Priority:** P1 — High  
**Description:** Company-specific report formatting

| ID | Requirement | Priority |
|----|------------|----------|
| FR-RPT-001a | Company info on report headers | P1 |
| FR-RPT-001b | Logo upload and display | P1 |
| FR-RPT-001c | Report title customization | P2 |
| FR-RPT-001d | Footer with signature lines | P2 |
| FR-RPT-001e | Paper size (A4/A5) | P2 |

---

## 3. Non-Functional Requirements

| ID | Requirement | Priority |
|----|------------|----------|
| NFR-001 | All configs stored in PostgreSQL (not env vars) | P0 |
| NFR-002 | Company-scoped (multi-tenant isolation) | P0 |
| NFR-003 | Default values match Vietnamese standards | P0 |
| NFR-004 | Config changes logged in audit trail | P1 |
| NFR-005 | Config API response < 100ms | P1 |
| NFR-006 | Backup file encrypted at rest | P2 |

---

## 4. API Design

### 4.1 System Options
```
GET    /api/v1/system-options              — Get all options
PUT    /api/v1/system-options              — Update options
GET    /api/v1/system-options/:category    — Get by category
```

### 4.2 Voucher Numbering
```
GET    /api/v1/numbering-rules              — List all rules
POST   /api/v1/numbering-rules              — Create rule
PUT    /api/v1/numbering-rules/:id          — Update rule
GET    /api/v1/numbering-rules/:type/next   — Get next number
```

### 4.3 Backup
```
POST   /api/v1/backups                     — Create backup
GET    /api/v1/backups                     — List backups
POST   /api/v1/backups/:id/restore         — Restore backup
GET    /api/v1/backups/:id/download         — Download backup file
```

### 4.4 Contracts
```
GET    /api/v1/contracts                   — List contracts
POST   /api/v1/contracts                   — Create contract
GET    /api/v1/contracts/:id               — Get contract
PUT    /api/v1/contracts/:id               — Update contract
DELETE /api/v1/contracts/:id               — Delete contract
GET    /api/v1/contracts/:id/entries       — Get linked entries
```

### 4.5 Loan Agreements
```
GET    /api/v1/loans                       — List loans
POST   /api/v1/loans                       — Create loan
GET    /api/v1/loans/:id                   — Get loan
PUT    /api/v1/loans/:id                   — Update loan
GET    /api/v1/loans/:id/schedule          — Get payment schedule
POST   /api/v1/loans/:id/payments          — Record payment
```

---

## 5. Database Schema

### 5.1 System Options Table
```sql
CREATE TABLE system_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    category VARCHAR(50) NOT NULL, -- 'personal', 'global', 'inventory', 'number_format', 'report'
    option_key VARCHAR(100) NOT NULL,
    option_value JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, category, option_key)
);
```

### 5.2 Numbering Rules Table
```sql
CREATE TABLE numbering_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    voucher_type VARCHAR(50) NOT NULL,
    prefix VARCHAR(20),
    suffix VARCHAR(20),
    number_length INT DEFAULT 5,
    scope VARCHAR(20) DEFAULT 'company', -- 'company' or 'branch'
    reset_rule VARCHAR(20) DEFAULT 'never', -- 'never', 'yearly', 'monthly'
    current_number INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, voucher_type)
);
```

### 5.3 Contracts Table
```sql
CREATE TABLE contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    contract_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    value NUMERIC(19,4),
    start_date DATE,
    end_date DATE,
    counterparty_name VARCHAR(255),
    counterparty_tax_code VARCHAR(20),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, code)
);
```

### 5.4 Loan Agreements Table
```sql
CREATE TABLE loan_agreements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    code VARCHAR(50) NOT NULL,
    lender_name VARCHAR(255) NOT NULL,
    loan_type VARCHAR(50) NOT NULL,
    principal NUMERIC(19,4) NOT NULL,
    annual_rate NUMERIC(8,4) NOT NULL,
    start_date DATE NOT NULL,
    maturity_date DATE NOT NULL,
    payment_frequency VARCHAR(20) DEFAULT 'monthly',
    gl_account_code VARCHAR(20) DEFAULT '331',
    status VARCHAR(20) DEFAULT 'active',
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, code)
);
```

### 5.5 Backups Table
```sql
CREATE TABLE backups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    filename VARCHAR(255) NOT NULL,
    file_size BIGINT,
    backup_type VARCHAR(20) DEFAULT 'manual',
    status VARCHAR(20) DEFAULT 'completed',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 6. Open Questions

1. Should system options be per-company or per-branch? (MISA supports both)
2. What's the maximum backup file size we should support?
3. Do we need Google Drive integration for backup, or is local file sufficient?
4. For contracts, do we need milestone tracking or just basic status?
5. For loans, do we need amortization schedule generation?

---

## 7. Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| BA Lead | | | |
| Chief Accountant | | | |
| Tech Lead | | | |
