# Bank Module — Use Cases

---

## UC-1: Import Bank Statement

**ID:** UC-01
**Actor:** Accountant
**Precondition:** Bank account exists, GL entries posted
**Trigger:** Monthly bank statement received from bank

### Happy Path
1. Accountant navigates to Bank → Statement Import
2. Selects bank account
3. Uploads CSV/MT940 file from bank
4. System parses file, validates format, detects duplicates
5. System displays preview: statement date range, opening/closing balance, line count
6. Accountant confirms import
7. System stores statement lines with PENDING match status
8. System logs import audit: user, timestamp, file hash

### Alternative Paths
- **A1: Manual entry** — Accountant enters lines manually via form
- **A2: Duplicate statement** — System detects existing statement for same period, prompts to replace or cancel
- **A3: Balance mismatch** — Opening balance ≠ last statement closing balance. Warns but allows import

### Exception Paths
- **E1: Unparseable file** — System returns specific parsing errors (wrong format, encoding issues)
- **E2: Network error** — Import fails, no partial data saved
- **E3: Empty statement** — No lines found in file

---

## UC-2: Reconcile Bank Account

**ID:** UC-02
**Actor:** Accountant / Chief Accountant
**Precondition:** Statement imported, GL entries exist for period
**Trigger:** End of month reconciliation requirement

### Happy Path
1. Accountant opens Bank → Reconciliation
2. Selects bank account, date range (month)
3. System displays statement lines (left) vs GL entries (right)
4. Accountant reviews auto-matched items (highlighted green)
5. Accountant manually matches remaining items by selecting statement line + GL entry
6. For unmatched items, accountant investigates (creates write-off or flags exception)
7. System shows: opening balance, closing balance, matched total, difference
8. Difference = 0 (all items matched or written off)
9. Accountant completes reconciliation
10. System locks period for this bank account

### Alternative Paths
- **A1: Partial reconciliation** — Difference ≠ 0, some items still unmatched. Accountant can save but not complete
- **A2: Write-off needed** — Bank fee / interest not in GL. Accountant creates write-off entry
- **A3: Reverse reconciliation** — Chief Accountant reverses completed reconciliation to fix errors

### Exception Paths
- **E1: Period already closed in GL** — Cannot reconcile for closed period
- **E2: Statement already reconciled** — Cannot reconcile same statement twice
- **E3: Balance discontinuity** — Missing earlier statement, gap in transaction history

---

## UC-3: Create and Submit Payment Order

**ID:** UC-03
**Actor:** Accountant (maker) → Chief Accountant (checker)
**Precondition:** Sufficient bank balance, approved invoices
**Trigger:** Supplier payment due / salary payment date

### Happy Path
1. Maker creates payment order: beneficiary, amount, bank account, content
2. System validates: balance check, duplicate check
3. Maker submits for approval
4. Checker receives approval notification
5. Checker reviews payment details
6. Checker approves
7. Maker prints UNC form (bank-specific template)
8. Maker submits to bank (manually or via electronic banking)
9. Maker confirms payment in system
10. System posts GL entry (Nợ 331/642, Có 112)

### Alternative Paths
- **A1: Batch payment** — Maker creates batch of 50+ supplier payments, submits batch for approval
- **A2: Reject** — Checker rejects with reason, maker revises
- **A3: Failed payment** — Bank returns error, maker investigates and re-submits

### Exception Paths
- **E1: Insufficient balance** — System warns, but allows if company policy permits
- **E2: Duplicate payment** — System detects same amount + same beneficiary in last 7 days
- **E3: Bank holiday** — Payment date falls on non-business day

---

## UC-4: Manage Loan Agreement

**ID:** UC-04
**Actor:** Treasury Manager
**Precondition:** Loan contract signed with bank
**Trigger:** New loan drawdown / monthly repayment

### Happy Path (Disbursement)
1. Treasury creates loan agreement: contract no, amount, interest rate, term
2. System records loan with ACTIVE status
3. Treasury records disbursement: date, amount, receiving bank account
4. System posts GL: Nợ 112, Có 341
5. System auto-generates repayment schedule

### Happy Path (Repayment)
1. Treasury opens loan agreement
2. Views repayment schedule (auto-generated)
3. Records payment: principal + interest amounts
4. System auto-updates outstanding balance
5. System posts GL: Nợ 341 (principal) + Nợ 635 (interest), Có 112

### Alternative Paths
- **A1: Early repayment** — Treasury records early repayment, system recalculates interest
- **A2: Overdraft** — Revolving credit facility, no fixed schedule

### Exception Paths
- **E1: Overdue** — System flags overdue repayment, sends alert
- **E2: Loan restructure** — Terms modified, system recalculates schedule

---

## UC-5: Manage Term Deposit

**ID:** UC-05
**Actor:** Treasury Manager
**Precondition:** Excess cash available
**Trigger:** Investment decision

### Happy Path
1. Treasury creates term deposit: amount, term, interest rate, bank
2. System records ACTIVE deposit
3. System posts GL: Nợ 1281, Có 112
4. At maturity, system auto-calculates interest
5. Treasury processes maturity
6. System posts GL: Nợ 112 (principal + interest), Có 1281 (principal), Có 515 (interest)
7. System marks deposit as MATURED

### Alternative Paths
- **A1: Auto-renewal** — If auto-renewal enabled, system creates new deposit with current rate
- **A2: Early withdrawal** — Treasury closes before maturity (with penalty)

### Exception Paths
- **E1: Interest rate change** — System uses contracted rate for fixed deposits
- **E2: Bank fails to credit** — Treasury follows up manually

---

## UC-6: Generate Bank Deposit Ledger (S08-DN)

**ID:** UC-06
**Actor:** Accountant
**Precondition:** Transactions posted for period
**Trigger:** Month-end reporting

### Happy Path
1. Accountant opens Reports → Bank Deposit Ledger
2. Selects bank account, date range
3. System generates S08-DN report
4. Report shows: opening balance, all transactions (date, ref, description, debit, credit), running balance, closing balance
5. Accountant exports to PDF/Excel/Print

### Alternative Paths
- **A1: Multi-currency** — Separate reports per currency
- **A2: Combined bank report** — All bank accounts aggregated

---

## UC-7: Generate Bank Reconciliation Report (S09-DN)

**ID:** UC-07
**Actor:** Chief Accountant
**Precondition:** Reconciliation completed
**Trigger:** Month-end / audit preparation

### Happy Path
1. CA opens Reports → Bank Reconciliation
2. Selects reconciliation period
3. System generates S09-DN: opening balance per bank vs GL, detailed matching, difference, closing balance
4. Report includes: all matched items (statement line + GL entry), unmatched items, write-offs
5. CA exports for auditor

---

## UC-8: Foreign Currency Bank Account Revaluation

**ID:** UC-08
**Actor:** Chief Accountant
**Precondition:** Foreign currency bank accounts have balance
**Trigger:** Period-end (monthly/quarterly/yearly)

### Happy Path
1. CA runs period-end FX revaluation
2. System gets exchange rate from configured bank (average transfer rate)
3. System calculates FX gain/loss for each FCY bank account
4. System previews revaluation entries
5. CA confirms
6. System posts: Nợ 1122 / Có 515 (gain) OR Nợ 635 / Có 1122 (loss)

### Exception Paths
- **E1: Rate not configured** — System prompts for manual rate entry
- **E2: Reversal needed** — If previous revaluation was posted, reversed first
