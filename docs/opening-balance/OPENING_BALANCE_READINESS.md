# Opening Balance Module — Production Readiness Assessment

**Version:** 1.0 | **Date:** 2026-07-27
**Author:** BA Lead + Chief Accountant (20+ yrs)
**Standard:** Circular 99/2025/TT-BTC, Circular 200/2014/TT-BTC, IAS 1, VAS 29

---

## Executive Verdict: NOT PRODUCTION READY

| Dimension | Score | Status |
|-----------|-------|--------|
| Data Model | 1/5 | Skeleton exists, no population logic |
| Service Layer | 0/5 | No OpeningBalanceService |
| Repository | 0/5 | No OpeningBalanceRepository (Memory or PG) |
| API Endpoints | 0/5 | No opening balance CRUD endpoints |
| Validation | 0/5 | No opening balance validation rules |
| Audit Trail | 0/5 | No audit logging for balance changes |
| Multi-currency | 0/5 | No foreign currency opening balance |
| Detail Tracking | 0/5 | No per-object/customer/supplier opening balance |
| Approval Workflow | 0/5 | No 4-eyes approval for balance changes |
| Import/Export | 0/5 | No Excel/CSV import of opening balances |
| Reporting | 0/5 | Opening balances always 0 in reports |
| Carry-Forward | 0/5 | No automatic carry-forward between periods |
| Circular 99 Compliance | 0/5 | Fails mandatory balance transfer requirements |
| **Overall** | **0.7/5** | **CRITICAL GAP — DO NOT DEPLOY** |

---

## Detailed Gap Analysis

### 1. Current Codebase State

```go
// domain/models.go:177-192
type AccountBalance struct {
    OpenBalanceDebit  float64  // NEVER POPULATED
    OpenBalanceCredit float64  // NEVER POPULATED
    PeriodDebit       float64  // Populated from journal lines
    PeriodCredit      float64  // Populated from journal lines
    // ...
}
func (b *AccountBalance) Calculate() {
    // Open* fields always 0 → closing = period only
    b.TotalDebit = b.OpenBalanceDebit + b.PeriodDebit
    b.ClosingBalance = b.TotalDebit - b.TotalCredit
}
```

Both MemoryJournalRepo (memory.go:258) and PGJournalRepo (pg.go:248) `GetBalance` methods:
- Only aggregate journal lines for the requested period
- Never look up opening balance from previous period
- Never store or retrieve opening balance from any table
- Return `OpenBalanceDebit = 0`, `OpenBalanceCredit = 0` always

### 2. What Exists vs What's Needed

| Component | Exists | Missing |
|-----------|--------|---------|
| `AccountBalance.OpenBalanceDebit` | Model field only | Population ↔ carry-forward |
| `AccountBalance.OpenBalanceCredit` | Model field only | Population ↔ carry-forward |
| `OpeningBalance` entity | No | Full struct with period, currency, detail tracking |
| `OpeningBalanceRepository` interface | No | Create, GetByPeriod, Update, Lock |
| `MemoryOpeningBalanceRepo` | No | In-memory implementation |
| `PGOpeningBalanceRepo` | No | PG implementation + migration |
| `OpeningBalanceService` | No | Business rules, validation, approval |
| `Handler` endpoints | No | CRUD, import, export, approve |
| Route registration | No | `/api/v1/opening-balances/*` |
| Migration SQL | No | `opening_balances` table DDL |

### 3. Functional Gaps (20+ years Chief Accountant perspective)

**Gap A — No initial setup wizard:**
When a company starts using GoTax mid-year (e.g., April 2026), there is no way to enter the cumulative balances from previous months. Accountants must create individual journal entries for each account balance — a 200+ line entry that must be hand-balanced. MISA, Fast, and BravoERP all provide dedicated opening balance screens with Excel import and auto-balance checking.

**Gap B — No opening balance modification workflow:**
If an opening balance error is discovered after posting (e.g., Account 1111 should be 125M but was entered as 125.5M), there's no controlled correction workflow. MISA requires chief accountant approval for any opening balance change, with full audit trail.

**Gap C — No fiscal year carry-forward:**
At year-end, revenue/expense accounts must be zeroed to Retained Earnings (421). Asset/Liability opening balances must carry forward to the new year. Current system has no `CloseYear` or `CarryForwardOpeningBalances` method. The BRD_GL_MODULE.md BR-03-7 states: "Opening balance carry-forward mandatory at year start" — this is unimplemented.

**Gap D — No multi-currency opening balance:**
Vietnamese enterprises commonly maintain USD, EUR bank accounts and receivables. MISA and Fast allow entering opening balance per currency with exchange rate. GoTax has no multi-currency opening balance support.

**Gap E — No receivable/payable detail opening:**
Account 131 (Receivables) and 331 (Payables) require per-customer/supplier breakdown. MISA's "Vào số dư công nợ đầu kỳ" function allows entering opening balances per customer with aging. GoTax has nothing.

**Gap F — No inventory opening balance:**
Account 152 (Raw materials), 153 (Tools), 155 (Finished goods), 156 (Merchandise) require quantity + unit price breakdown. Fast provides separate "Tồn kho đầu kỳ" screens. GoTax has no inventory opening balance.

**Gap G — No fixed asset opening balance:**
Account 211 (Fixed assets), 213 (Intangible assets) require original cost, accumulated depreciation, remaining life. MISA provides dedicated Fixed Asset opening screens. GoTax has none.

**Gap H — No opening balance audit trail:**
Every balance change must be logged per VSA (Vietnamese Standards on Auditing) and Circular 99 Article 12 (Accounting Books). GoTax has no audit trail for balance modifications.

**Gap I — No opening balance freeze/lock:**
Once a period is closed, its opening balances must be immutable. Current system has no balance freeze mechanism.

**Gap J — No Circular 99 transitional compliance:**
Circular 99 Section 30 (Transitional Provisions) and KPMG/PWC guidance require:
- Mandatory account balance transfers on 31 Dec 2025 (441→4118, 466→4118, 138→2281, etc.)
- Validate opening balances in new COA structure by 1 Jan 2026
- Retrospective adjustment method per VAS 29
None of this is supported.

### 4. Regulatory Compliance Gaps

| Regulation | Requirement | Status |
|-----------|-------------|--------|
| Circular 99 Art 13 | Opening/closing accounting books with correct balances | FAIL |
| Circular 99 Appendix IV | Financial statement forms with opening balance column | FAIL |
| VAS 01 | Opening balance must equal previous period's closing | FAIL |
| VAS 29 | Retrospective adjustments for accounting policy changes | FAIL |
| IAS 1.40A | Third balance sheet for retrospective restatement | FAIL |
| Decree 123/2020 | Tax declaration opening balances must match GL | FAIL |

### 5. Competitive Comparison

| Feature | GoTax | MISA SME 2026 | Fast Business Online | BravoERP | Tryton |
|---------|-------|---------------|---------------------|----------|--------|
| Account OB entry | ✗ | ✓ | ✓ | ✓ | ✓ (manual JE) |
| Multi-currency OB | ✗ | ✓ | ✓ | ✓ | ✓ |
| Receivable/payable detail | ✗ | ✓ | ✓ | ✓ | ✓ |
| Inventory OB | ✗ | ✓ | ✓ | ✓ | via module |
| Fixed asset OB | ✗ | ✓ | ✓ | ✓ | via module |
| Excel import | ✗ | ✓ | ✓ | ✓ | ✓ |
| Auto-balance check | ✗ | ✓ | ✓ | ✓ | journal balance |
| Approval workflow | ✗ | ✓ | ✓ | partial | ✗ |
| Audit trail | ✗ | ✓ | ✓ | ✓ | ✓ |
| Fiscal year carry-forward | ✗ | auto | auto | auto | wizard |
| Closing adjustment period | ✗ | ✓ | ✓ | ✓ | ✓ |
| Circular 99 transitional mapping | ✗ | planned | planned | partial | N/A |

### 6. Recommendation

**DO NOT deploy GoTax with any real company data** until the Opening Balance module is fully implemented. The system cannot produce accurate financial statements because all opening balances are zero.

**Minimum viable implementation order:**
1. OpeningBalance data model + migration (2 days)
2. OpeningBalanceRepository (both backends) (1 day)
3. OpeningBalanceService (CRUD, validation, balancing) (2 days)
4. Handler + API endpoints (1 day)
5. Balance carry-forward + report integration (2 days)
6. Excel import (1 day)
7. Receivable/payable detail OB (2 days)
8. Approval workflow (1 day)
9. Audit trail (0.5 day)
10. Circular 99 transitional support (2 days)

**Total estimate: ~15 days for full production readiness**

### 7. Risk Assessment

| Risk | Severity | Probability | Mitigation |
|------|----------|-------------|------------|
| Financial statements wrong | CRITICAL | 100% | Block deployment |
| Tax filing errors | CRITICAL | 100% | Block deployment |
| Audit failure | HIGH | 100% | Block deployment |
| Data corruption on carry-forward | HIGH | 80% | Design review needed |
| User data entry errors | MEDIUM | 60% | Auto-balance validation |

---

## Appendix: Benchmarking Sources

- MISA SME 2026: helpact.misa.vn — full opening balance module
- Fast Business Online: fbohelp.fast.com.vn — 8 separate opening balance functions
- BravoERP: partner documentation — opening balance wizard
- Tryton: docs.tryton.org — "Fill opening balance" using Situation Journal
- KPMG Vietnam: Circular 99 transitional guidance
- PWC Vietnam: VAS 29 retrospective adjustment guidance
- EY Vietnam: Opening balance validation procedures
- Indochina Link: Circular 99 implementation roadmap
