# Implementation Roadmap — Admin/Config Modules

**Document ID:** GOTAX-ROADMAP-ADMIN-001  
**Version:** 1.0  
**Date:** 2026-08-07  
**Total Estimated Effort:** 82 days (16.4 weeks)

---

## Phase 1: System Foundation (Weeks 1-3)

### Goal
Unlock all other modules' configurability. Without this, nothing else can be properly configured.

### 1.1 System Options (FR-SYS-001)
**Effort:** 10 days  
**Priority:** P0 — Critical

| Day | Task | Deliverable |
|-----|------|-------------|
| 1 | Create migration `000035_system_options.up.sql` | `system_options` table |
| 2 | Create domain model `models_system.go` | `SystemOption` struct |
| 3 | Create interface `SystemOptionRepo` in `interfaces.go` | Repository interface |
| 4 | Implement `pg_system_option.go` | PostgreSQL repo |
| 5 | Implement `memory_system_option.go` | Memory repo |
| 6 | Create `system_option_service.go` | Service with defaults |
| 7 | Create `system_option_handler.go` | HTTP handlers |
| 8 | Register routes + wire in `main.go` | Integration |
| 9 | Write handler tests | `system_option_handler_test.go` |
| 10 | Update swagger + verify | Documentation |

**Default Values (Vietnamese Standards):**
```json
{
  "accounting_standard": "circular99",
  "fiscal_year_start_month": 1,
  "base_currency": "VND",
  "costing_method": "weighted_avg",
  "thousands_separator": ".",
  "decimal_separator": ",",
  "decimal_places": 0,
  "negative_display": "parentheses"
}
```

---

### 1.2 Voucher Numbering Rules (FR-NUM-001)
**Effort:** 6 days  
**Priority:** P0 — Critical

| Day | Task | Deliverable |
|-----|------|-------------|
| 1 | Create migration `000036_numbering_rules.up.sql` | `numbering_rules` table |
| 2 | Create domain model + interface | `NumberingRule` struct |
| 3 | Implement PG + Memory repos | Repository layer |
| 4 | Create service with auto-increment logic | `numbering_rule_service.go` |
| 5 | Create handler + register routes | `numbering_rule_handler.go` |
| 6 | Write tests | `numbering_rule_handler_test.go` |

**Default Rules (MISA standard):**
| Voucher Type | Prefix | Length | Scope | Reset |
|-------------|--------|--------|-------|-------|
| Journal Entry | BG | 5 | company | never |
| Purchase | MB | 5 | company | never |
| Sale | BH | 5 | company | never |
| Cash Receipt | TN | 5 | company | never |
| Cash Payment | CN | 5 | company | never |
| Goods Receipt | NK | 5 | company | never |
| Goods Issue | XK | 5 | company | never |

---

### 1.3 Fiscal Year Configuration (FR-FY-001)
**Effort:** 4 days  
**Priority:** P0 — Critical

| Day | Task | Deliverable |
|-----|------|-------------|
| 1 | Extend `periods` table with fiscal_year fields | Migration |
| 2 | Add fiscal year config to SystemOption | Integration |
| 3 | Auto-generate periods on FY creation | Service logic |
| 4 | Write tests | Verification |

**Note:** Periods table already exists. This adds fiscal year metadata and auto-generation.

---

### 1.4 Number Format (FR-SYS-004)
**Effort:** 3 days  
**Priority:** P0 — Critical

| Day | Task | Deliverable |
|-----|------|-------------|
| 1 | Add number format options to SystemOption | Defaults |
| 2 | Create `format.go` utility package | Formatting functions |
| 3 | Write unit tests | `format_test.go` |

**Utility Functions:**
```go
func FormatNumber(value float64, opts NumberFormat) string
func ParseNumber(s string, opts NumberFormat) (float64, error)
```

---

## Phase 2: Missing Business Modules (Weeks 4-6)

### Goal
Complete the 19-module MISA parity.

### 2.1 Contracts (FR-CON-001)
**Effort:** 9 days  
**Priority:** P1 — High

| Day | Task | Deliverable |
|-----|------|-------------|
| 1 | Create migration `000037_contracts.up.sql` | `contracts` + `contract_payments` tables |
| 2 | Create domain models | `Contract`, `ContractPayment` structs |
| 3 | Create interfaces | Repository interfaces |
| 4 | Implement PG repo | `pg_contract.go` |
| 5 | Implement Memory repo | `memory_contract.go` |
| 6 | Create service | `contract_service.go` |
| 7 | Create handler + routes | `contract_handler.go` |
| 8 | Write handler tests | `contract_handler_test.go` |
| 9 | Wire in main.go + verify | Integration |

**API Endpoints:**
```
GET    /api/v1/contracts
POST   /api/v1/contracts
GET    /api/v1/contracts/:id
PUT    /api/v1/contracts/:id
DELETE /api/v1/contracts/:id
GET    /api/v1/contracts/:id/payments
POST   /api/v1/contracts/:id/payments
```

---

### 2.2 Loan Agreements (FR-LOAN-001)
**Effort:** 9 days  
**Priority:** P1 — High

| Day | Task | Deliverable |
|-----|------|-------------|
| 1 | Create migration `000038_loan_agreements.up.sql` | Tables |
| 2 | Create domain models | `LoanAgreement`, `LoanPayment` structs |
| 3 | Create interfaces | Repository interfaces |
| 4 | Implement PG repo | `pg_loan.go` |
| 5 | Implement Memory repo | `memory_loan.go` |
| 6 | Create service with amortization | `loan_service.go` |
| 7 | Create handler + routes | `loan_handler.go` |
| 8 | Write handler tests | `loan_handler_test.go` |
| 9 | Wire in main.go + verify | Integration |

**Amortization Formula (Circular 99):**
```
Monthly Interest = Outstanding Balance × Annual Rate / 12
Principal Payment = Fixed Payment - Interest
New Balance = Old Balance - Principal Payment
```

**GL Accounts:**
- 131: Loans made to others (asset)
- 331: Loans from banks (liability)

---

## Phase 3: Integration (Weeks 7-9)

### Goal
Enable real-world production use with external integrations.

### 3.1 E-Banking Integration (FR-EBK-001)
**Effort:** 10 days  
**Priority:** P2 — Medium

| Day | Task | Deliverable |
|-----|------|-------------|
| 1-2 | Design bank statement import format | CSV/MT940 parser |
| 3-4 | Implement parsers for top 5 banks | VCB, BIDV, CTG, VTB, ACB |
| 5-6 | Create import service + handler | `bank_import_service.go` |
| 7-8 | Create reconciliation logic | Match to existing entries |
| 9-10 | Write tests + verify | Integration tests |

**Bank Statement Format (CSV):**
```csv
日期,Diễn giải,Số tiền,Còn lại,Loại giao dịch
2026-08-01,"Payment to supplier",5000000,95000000,"Debit"
2026-08-02,"Customer payment",10000000,105000000,"Credit"
```

---

### 3.2 E-Tax Filing (FR-ETX-001)
**Effort:** 10 days  
**Priority:** P2 — Medium

| Day | Task | Deliverable |
|-----|------|-------------|
| 1-2 | Extend GDT client for full filing | `gdt_client.go` |
| 3-4 | Implement VAT declaration (Mẫu 01) | XML generation |
| 5-6 | Implement CIT declaration (Mẫu 01) | XML generation |
| 7-8 | Implement PIT declaration | XML generation |
| 9-10 | Write tests + verify | Integration tests |

**Tax Forms (Decree 123/2020):**
- VAT: Mẫu 01/GTGT (monthly), 02/GTGT (quarterly)
- CIT: 01/TNDN (quarterly), 03/TNDN (annual)
- PIT: 05/TNCN (monthly)

---

### 3.3 Digital Signature (FR-DSG-001)
**Effort:** 5 days  
**Priority:** P2 — Medium

| Day | Task | Deliverable |
|-----|------|-------------|
| 1-2 | Design signature provider interface | `SignatureProvider` interface |
| 3-4 | Implement VNPT-CA provider | Provider implementation |
| 5 | Write tests + docs | Verification |

**Providers (Vietnam):**
- VNPT-CA (most common)
- Viettel-CA
- BKAV-CA
- FPT-CA

---

## Phase 4: Polish (Weeks 10-12)

### Goal
Improve UX and operational readiness.

### 4.1 Backup & Restore (FR-BKP-001)
**Effort:** 5 days  
**Priority:** P1 — High

| Day | Task | Deliverable |
|-----|------|-------------|
| 1 | Create migration + domain model | Tables |
| 2 | Implement backup service (pg_dump) | `backup_service.go` |
| 3 | Implement restore service | Restore logic |
| 4 | Create handler + routes | `backup_handler.go` |
| 5 | Write tests | Verification |

---

### 4.2 Report Customization (FR-RPT-001)
**Effort:** 5 days  
**Priority:** P1 — High

| Day | Task | Deliverable |
|-----|------|-------------|
| 1 | Add report options to SystemOption | Config integration |
| 2 | Extend PDF generator with company info | `maroto` updates |
| 3 | Add logo upload endpoint | File handling |
| 4 | Add report header/footer templates | Template system |
| 5 | Write tests | Verification |

---

### 4.3 Multi-branch Enhancements
**Effort:** 6 days  
**Priority:** P2 — Medium

| Day | Task | Deliverable |
|-----|------|-------------|
| 1-2 | Branch-level options support | SystemOption branch scope |
| 3-4 | Branch numbering rules | Per-branch sequences |
| 5-6 | Branch report filtering | Report by branch |

---

## Summary

| Phase | Weeks | Days | Modules |
|-------|-------|------|---------|
| 1: System Foundation | 1-3 | 23 | System Options, Numbering, Fiscal Year, Number Format |
| 2: Business Modules | 4-6 | 18 | Contracts, Loan Agreements |
| 3: Integration | 7-9 | 25 | E-Banking, E-Tax, Digital Signature |
| 4: Polish | 10-12 | 16 | Backup, Report Customization, Multi-branch |
| **TOTAL** | **12** | **82** | **12 modules** |

---

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Scope creep on system options | High | Stick to MISA parity, no custom features |
| Bank statement format variations | Medium | Start with top 5 banks, add others later |
| GDT API changes | Medium | Abstract interface, easy to update |
| Backup file size limits | Low | Compress, set reasonable limits |
| Multi-branch complexity | Medium | Phase 4, can defer if needed |

---

## Dependencies

| Dependency | Type | Impact |
|-----------|------|--------|
| PostgreSQL | External | Required for all features |
| GDT API | External | Required for E-Tax |
| Bank APIs | External | Required for E-Banking |
| Digital cert providers | External | Required for signatures |
| GoTax core modules | Internal | All 17 modules must be stable first |

---

## Next Steps

1. **Review & approve** this roadmap with stakeholders
2. **Start Phase 1** with System Options (FR-SYS-001)
3. **Create detailed specs** for each module before implementation
4. **Set up CI/CD** for automated testing
5. **Begin implementation** on 2026-08-11 (Monday)
