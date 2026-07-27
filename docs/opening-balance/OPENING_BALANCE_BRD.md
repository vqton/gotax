# Business Requirements Document — Opening Balance Module

**Version:** 1.0 | **Date:** 2026-07-27
**Author:** BA Lead + Chief Accountant (combined 40+ years)
**Standard:** Circular 99/2025/TT-BTC (effective 1 Jan 2026)
**Supersedes:** Circular 200/2014/TT-BTC

---

## 1. Executive Summary

Opening Balance (Số dư đầu kỳ / SDĐK) is the foundational data set of any General Ledger system. It represents the cumulative balances of all accounts at the start of a fiscal period. Without accurate opening balances, every financial statement, tax declaration, and management report produced by GoTax is incorrect.

Current GoTax implementation has zero opening balance functionality. The `AccountBalance` struct has `OpenBalanceDebit`/`OpenBalanceCredit` fields but they are never populated — always zero. This BRD defines the complete business requirements for a production-grade Opening Balance module compliant with Vietnamese regulations (Circular 99), IFRS (IAS 1), and competitive with MISA/Fast/BravoERP.

## 2. Scope

### In Scope
- Opening balance entry for all account types (balance sheet + off-balance-sheet)
- Multi-currency opening balances
- Receivable/payable detail opening (per customer/supplier/employee)
- Inventory opening (quantity + value)
- Fixed asset opening (cost + accumulated depreciation)
- Excel/CSV bulk import with validation
- Opening balance approval workflow (4-eyes)
- Fiscal year carry-forward
- Opening balance audit trail
- Circular 99 transitional balance mapping
- Integration with trial balance, balance sheet, income statement reports

### Out of Scope
- WIP (Work in Progress) opening balance → separate module
- Cost center/project opening balance → separate module
- Budget opening balance → separate module

## 3. Business Context

### 3.1 When Opening Balances Are Needed

| Scenario | Description | Frequency |
|----------|-------------|-----------|
| **New implementation** | Company starts using GoTax for first time | Once |
| **New fiscal year** | Balances carry forward from prior year | Annual |
| **System migration** | Switching from another accounting software | Once |
| **Merger/acquisition** | New entity added mid-year | Per event |
| **Correction** | Prior period error discovered | Per event |
| **Circular transition** | TT200 → TT99 account mapping (2025->2026) | Once (Dec 2025) |

### 3.2 User Personas

| Role | Responsibilities | Permissions |
|------|-----------------|-------------|
| **Accountant** | Enter opening balances, import from Excel | read, write, import |
| **Chief Accountant** | Review, approve, correct, freeze | approve, freeze, correct |
| **Admin** | System setup, company configuration | all |
| **Auditor** | View opening balance history | read-only |

### 3.3 Regulatory Requirements

**Circular 99/2025/TT-BTC, Article 13 — Opening and Closing of Accounting Books:**
> "Accounting books must be opened at the beginning of each accounting period. The opening balances of the current period must equal the closing balances of the immediately preceding period. If an enterprise applies a new accounting policy, the opening balances must be adjusted retrospectively in accordance with VAS 29."

**Circular 99, Article 30 — Transitional Provisions:**
> "By 31 December 2025, enterprises must execute mandatory account balance transfers as specified in Appendix II. Opening balances as of 1 January 2026 must comply with the new Chart of Accounts structure."

**VAS 01 — Accounting Standard No. 01:**
> "Financial statements must present corresponding figures for the previous period. Opening balances must be consistent with closing balances of the previous period."

**IFRS IAS 1.40A:**
> "When an entity applies an accounting policy retrospectively or makes a retrospective restatement, it must present a third statement of financial position as at the beginning of the preceding period."

## 4. Functional Requirements

### FR-01: Opening Balance Entry by Account
The system must allow entry of debit/credit opening balance for any active account.

**Detail:**
- Input: account code, debit amount, credit amount, currency, date
- Only leaf accounts (non-parent) can receive opening balances
- Opening balance date = accounting start date (configurable per company)
- Auto-populate from previous period closing balance when available
- Manual override with chief accountant approval

### FR-02: Opening Balance by Detail Entity
The system must support opening balance breakdown by detail entity:

| Account Type | Detail Entity | Example |
|-------------|---------------|---------|
| 112 (Bank) | Bank account | 1121-VCB: 500M, 1121-CTG: 300M |
| 131 (Receivables) | Customer | KH001: 50M, KH002: 30M |
| 331 (Payables) | Supplier | NCC001: 20M, NCC002: 15M |
| 141 (Advances) | Employee | NV001: 5M, NV002: 3M |
| 334 (Salary) | Employee | Per employee |
| 152 (Materials) | Item + quantity | VL001: 100kg x 50k = 5M |
| 211 (Fixed Assets) | Asset item | TSCD001: cost 500M, dep 200M |

### FR-03: Multi-Currency Opening Balance
The system must allow opening balance entry per currency:
- VND (default)
- USD, EUR, JPY, etc.
- Exchange rate at opening date
- Auto-compute VND equivalent

### FR-04: Opening Balance Validation
Automatic validation on save:

| Rule | Error | Severity |
|------|-------|----------|
| Total debit = Total credit | "Opening balance not balanced" | BLOCKER |
| Each line: debit XOR credit > 0 | "Amount required" | BLOCKER |
| Account must exist and be active | "Account not found" | BLOCKER |
| Account must be leaf (non-parent) | "Parent account cannot have balance" | BLOCKER |
| Account type matches normal balance | "Account 1111 (ASSET) should have debit balance" | WARNING |
| Currency exists and is active | "Currency not configured" | BLOCKER |
| Opening date must be within open period | "Period not open" | BLOCKER |

### FR-05: Opening Balance Approval Workflow
```
Accountant enters OB → System validates → Chief Accountant reviews → Approve/Reject
                                            ↓
                                   If rejected: return with reason
```
- Approval required for: new company setup, modification of approved balance, balance > threshold (configurable)
- No approval for: first-time entry during initial setup (chief accountant can bypass)
- 4-eyes principle: approver ≠ creator

### FR-06: Fiscal Year Carry-Forward
At fiscal year close:
1. Zero out revenue (5xx, 7xx, 8xx) and expense (6xx, 8xx) to Account 421 (Retained Earnings)
2. Carry forward asset (1xx, 2xx), liability (3xx), equity (4xx) closing balances as opening balances for new year
3. Generate audit trail for each carried-forward balance

### FR-07: Opening Balance Correction Workflow
If error discovered after approval:
1. Create correction request with old/new values and reason
2. Chief Accountant reviews
3. If approved: system creates reversal entry + corrected balance
4. Full audit trail preserved

### FR-08: Excel/CSV Import
- Template: account code, debit, credit, currency, detail columns
- Encoding: UTF-8, VNI, TCVN3 auto-detect
- Preview before import (X rows valid, Y errors)
- Transactional: all-or-nothing within each batch
- Error report generation

### FR-09: Circular 99 Transitional Support
- Built-in mapping table: Circular 200 accounts → Circular 99 accounts
- Auto-generate transfer journal entries for:
  - 441 (Capital construction) → 4118 (Other capital)
  - 466 (Fixed asset forming fund) → 4118
  - 138 (BCC investment) → 2281
  - 2413 (Extraordinary repair) → 2414
- Migration validation report

### FR-10: Opening Balance Audit Trail
Every balance change must log:
- Who (user ID, name)
- When (timestamp)
- What (account, old debit, new debit, old credit, new credit)
- Why (reason, reference)
- Approval reference (if applicable)

## 5. Non-Functional Requirements

| Requirement | Specification |
|-------------|---------------|
| Performance | 10,000 opening balance lines in < 3 seconds |
| Concurrency | Multiple accountants can enter different accounts simultaneously |
| Data integrity | No partial saves; all-or-nothing per batch |
| Audit readiness | 7-year retention per VSA |
| Availability | 99.9% uptime |
| Security | Role-based access; sensitive accounts require 2FA |

## 6. Integration Points

| Module | Integration |
|--------|-------------|
| Trial Balance | Opening balances flow into trial balance report |
| Balance Sheet | Opening balances appear in B01-DN column "Đầu kỳ" |
| Income Statement | Prior period adjustments affect retained earnings |
| Tax Declaration | Opening balances affect VAT/CIT carry-forward |
| Audit Log | All OB changes logged to audit_entries table |
| Company Settings | Opening balance date, auto-approval threshold |

## 7. Glossary

| Vietnamese | English | Definition |
|-----------|---------|------------|
| Số dư đầu kỳ (SDĐK) | Opening Balance | Balance of an account at period start |
| Số dư cuối kỳ (SDCK) | Closing Balance | Balance of an account at period end |
| Kết chuyển | Carry-Forward | Transfer closing balance → next period opening |
| Nhập số dư ban đầu | Initial Balance Entry | First-time entry when starting system |
| Điều chỉnh hồi tố | Retrospective Adjustment | Correction applied to prior period balances |
| Số dư Nợ | Debit Balance | Normal balance for asset/expense accounts |
| Số dư Có | Credit Balance | Normal balance for liability/equity/revenue accounts |
| Chênh lệch | Variance/OFF | Difference when total debit ≠ total credit |

## 8. Acceptance Criteria

1. Accountant can enter opening balance for any leaf account via UI
2. System rejects unbalanced opening balances (debit ≠ credit)
3. Import 5000 opening balance lines from Excel in < 5 seconds
4. Opening balance appears in Trial Balance report
5. Opening balance appears in Balance Sheet (B01-DN)
6. Fiscal year carry-forward generates correct opening balances for new year
7. Approval workflow prevents unapproved balance changes
8. Audit trail captures every balance modification
9. Circular 99 transitional mapping generates correct transfer entries
10. Multi-currency opening balances display correctly in reports
