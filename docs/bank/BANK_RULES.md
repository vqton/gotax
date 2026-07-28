# Bank Module — Business Rules

---

## R1: Bank Account Rules

| ID | Rule | Description |
|----|------|-------------|
| R1.1 | Unique account per bank | One account number per bank code per company (soft unique) |
| R1.2 | One default per currency | Exactly one ACTIVE bank account can be DEFAULT per currency |
| R1.3 | Account number validation | Min 6, max 20 digits (Vietnamese standard) |
| R1.4 | Required fields | Bank name, account number, account holder, currency |
| R1.5 | Status transition | ACTIVE → CLOSED (no reactivation) |
| R1.6 | Closed account restriction | No new transactions on CLOSED accounts |
| R1.7 | FK constraint | CompanyBankAccount.CompanyID must reference existing Company |
| R1.8 | Audit logging | All mutations logged in audit log |

## R2: Statement Import Rules

| ID | Rule | Description |
|----|------|-------------|
| R2.1 | Unique statement period | One statement per bank account per date range |
| R2.2 | Balance continuity | Next statement opening balance = previous statement closing balance |
| R2.3 | Line balance check | Closing balance = opening balance + credits - debits |
| R2.4 | Duplicate line detection | Same bank_ref + amount + date within 7 days = duplicate |
| R2.5 | Required line fields | Transaction date, description, at least one of {debit, credit} |
| R2.6 | Currency consistency | All lines in statement must match statement currency |
| R2.7 | Import audit | Raw file hash stored, raw content in each line (raw_data field) |
| R2.8 | Fiscal year check | Statement date must be within open fiscal year period |

## R3: Auto-Match Rules

| ID | Rule | Priority | Confidence | Type |
|----|------|----------|------------|------|
| R3.1 | Exact amount + same date | 1 | 0.95 | High |
| R3.2 | Exact amount + ±1 day | 2 | 0.90 | High |
| R3.3 | Reference number match | 3 | 0.95 | High |
| R3.4 | Fuzzy reference match | 4 | 0.80 | Medium |
| R3.5 | Counterparty fuzzy + amount | 5 | 0.75 | Medium |
| R3.6 | Rounding tolerance (< 500 VND) | 6 | 0.50 | Low |
| R3.7 | Customer AR payment | 7 | 0.80 | Medium |
| R3.8 | Supplier AP payment | 8 | 0.80 | Medium |
| R3.9 | Salary batch match | 9 | 0.70 | Medium |

### Match Deduplication
- Each statement line can match exactly ONE GL entry (1:1)
- Each GL entry can match exactly ONE statement line (1:1)
- Higher confidence wins when multiple matches found
- Auto-matches must be confirmed by user (soft match → confirm)

## R4: Reconciliation Rules

| ID | Rule | Description |
|----|------|-------------|
| R4.1 | Match completeness | Reconciliation can complete only when difference = 0 |
| R4.2 | Write-off limit | Max write-off per reconciliation: 1,000,000 VND (configurable) |
| R4.3 | Write-off GL posting | Write-off → Nợ 635 / Có 112 for bank fees |
| R4.4 | Period lock | Completed reconciliation locks bank account for period |
| R4.5 | Reverse authority | Only Chief Accountant (or admin) can reverse reconciliation |
| R4.6 | Reverse effect | Reversal restores PENDING match status for all matched lines |
| R4.7 | Reconciliation sequence | Must reconcile in chronological order (date ASC) |
| R4.8 | Audit trail | All reconciliation actions logged |

## R5: Payment Order Rules

| ID | Rule | Description |
|----|------|-------------|
| R5.1 | Maker-checker principle | Creator ≠ Approver. Different users required |
| R5.2 | Balance check warning | Warn if total > 80% of current bank balance (configurable) |
| R5.3 | Duplicate detection | Flag if same amount + same beneficiary within 7 days |
| R5.4 | Max single payment | 5 billion VND (configurable per company) |
| R5.5 | Status transitions | DRAFT → PENDING_APPROVAL → APPROVED → SUBMITTED → CONFIRMED |
| R5.6 | Cancel allowed | Can cancel only if status ≤ APPROVED |
| R5.7 | Approval authority | Payment > threshold (100M) requires Chief Accountant + Director |
| R5.8 | Batch consistency | All orders in batch must have same currency + from bank account |
| R5.9 | Weekend handling | Payment date automatically adjusted to next business day if weekend |
| R5.10 | Beneficiary validation | Beneficiary account must be non-empty, 6-20 digits |

## R6: Loan Rules

| ID | Rule | Description |
|----|------|-------------|
| R6.1 | Disbursement ≤ principal | Sum of disbursements cannot exceed principal amount |
| R6.2 | Interest calculation | Simple interest: principal × rate × days/365 |
| R6.3 | Annuity formula | PMT = P × r(1+r)^n / ((1+r)^n - 1) |
| R6.4 | Overdue interest | Overdue interest rate = contract rate × 1.5 (common practice) |
| R6.5 | Restructure limits | Max 2 restructures per loan |
| R6.6 | Balance check | Outstanding balance cannot go below 0 |
| R6.7 | Currency matching | Loan currency must match bank account currency |

## R7: Term Deposit Rules

| ID | Rule | Description |
|----|------|-------------|
| R7.1 | Min amount | 50,000,000 VND (or equivalent FCY) |
| R7.2 | Min term | 7 days |
| R7.3 | Interest at maturity | Amount × rate × term_days / 365 |
| R7.4 | Early withdrawal penalty | Interest rate reduced to demand deposit rate (0.1-0.5%) |
| R7.5 | Auto-renewal rate | At maturity, new rate = current offered rate |
| R7.6 | GL posting at maturity | Principal → 112, Interest → 515 |

## R8: Revaluation Rules

| ID | Rule | Description |
|----|------|-------------|
| R8.1 | Rate source | Average TTM rate of bank where account is opened (Circular 99 Art. 6) |
| R8.2 | Revaluation frequency | At minimum: year-end. Recommended: monthly/quarterly |
| R8.3 | Gain/Loss recognition | Gain → 515 (Financial income), Loss → 635 (Financial expense) |
| R8.4 | Reverse on next period | Previous revaluation reversed, new revaluation posted |
| R8.5 | Materiality threshold | Ignore if gain/loss < 100,000 VND (configurable) |

## R9: General Ledger Posting Rules

| ID | Transaction | Debit | Credit |
|----|-------------|-------|--------|
| R9.1 | Customer payment received | 112 | 131 |
| R9.2 | Supplier payment | 331 | 112 |
| R9.3 | Bank transfer (to same bank) | 112 (target) | 112 (source) |
| R9.4 | Bank withdrawal to cash | 111 | 112 |
| R9.5 | Cash deposit to bank | 112 | 111 |
| R9.6 | Loan disbursement | 112 | 341 |
| R9.7 | Loan repayment (principal) | 341 | 112 |
| R9.8 | Loan repayment (interest) | 635 | 112 |
| R9.9 | Term deposit opening | 1281 | 112 |
| R9.10 | Term deposit maturity (principal) | 112 | 1281 |
| R9.11 | Term deposit maturity (interest) | 112 | 515 |
| R9.12 | Bank fee | 635 | 112 |
| R9.13 | FX revaluation gain | 1122 | 515 |
| R9.14 | FX revaluation loss | 635 | 1122 |
| R9.15 | Tax payment via bank | 333 | 112 |

## R10: Security & Compliance Rules

| ID | Rule | Description |
|----|------|-------------|
| R10.1 | Sensitive field encryption | Bank account number encrypted at rest (AES-256-GCM) |
| R10.2 | Payment approval | 2-person rule for all payment orders > 0 VND |
| R10.3 | Audit log | All CREATE/UPDATE/DELETE on bank entities logged |
| R10.4 | Data retention | 10 years (Circular 99 Art. 41, Law on Accounting Art. 13) |
| R10.5 | IP restriction | Payment submission restricted to whitelisted IPs (configurable) |
| R10.6 | Rate limiting | Max 100 payment creations per hour per user (configurable) |
| R10.7 | Anti-money laundering | Flag payments > 500M VND or suspicious patterns |
| R10.8 | Data isolation | Bank data scoped to company (multi-tenant) |
