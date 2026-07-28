# Cash Module — Use Cases

---

## UC-01: Process Customer Cash Receipt

**Actor:** Kế toán thanh toán (AP/AR Accountant)
**Trigger:** Customer pays debt in cash

### Happy Path
1. Accountant selects "Create Receipt"
2. Selects receipt type = "Customer Payment"
3. Selects customer from AR list (auto-fills name)
4. Enters amount, currency, exchange rate (if foreign)
5. System auto-selects Debit = 1111/1112, prompts for Credit account
6. Enters reason, attaches supporting documents
7. Saves as Draft
8. Submits for approval
9. Kế toán trưởng reviews, approves
10. System updates status to Approved
11. Accountant posts — system creates JournalEntry (Nợ 111/Có 131)
12. System updates cash book balance, updates AR aging
13. Print Phiếu thu (Mẫu 01-TT)

### Alternative Paths
- **A1: Cash sales (not AR):** Credit = 511 (Revenue) + 3331 (VAT)
- **A2: Advance refund:** Credit = 141 (Advance)
- **A3: Bank withdrawal:** Receipt type = "Bank Withdrawal", auto-links to CashTransfer

### Exception Paths
- **E1: Duplicate voucher number:** System generates unique number, reject duplicate
- **E2: Insufficient permissions:** User without submit right cannot submit
- **E3: Posted receipt cannot be edited:** Allow only reversal
- **E4: Invalid account:** Debit account not cash → reject

---

## UC-02: Process Supplier Cash Payment

**Actor:** Kế toán thanh toán
**Trigger:** Company pays supplier in cash

### Happy Path
1. Accountant selects "Create Payment"
2. Selects payment type = "Supplier Payment"
3. Selects supplier from AP list (auto-fills)
4. Selects invoice to pay (optional, auto-fills amount)
5. Enters amount, currency, exchange rate
6. System validates: balance sufficient (if check enabled)
7. Debit = 331 (AP), Credit = 1111/1112
8. Enters reason, attaches invoice copy
9. Saves as Draft → Submits → Kế toán trưởng approves
10. Posts — JournalEntry (Nợ 331/Có 111)
11. Updates cash book, AP aging
12. Print Phiếu chi (Mẫu 02-TT)

### Alternative Paths
- **A1: Expense payment:** Debit = 641/642 (expense)
- **A2: Salary payment:** Debit = 334 (payable to employees)
- **A3: Tax payment:** Debit = 333 (tax payable)
- **A4: Advance to employee:** Debit = 141 (tạm ứng)

### Exception Paths
- **E1: Insufficient cash balance:** Warn but allow override (configurable)
- **E2: Exceeds approval threshold:** Route to Giám đốc approval
- **E3: Payment > invoice balance:** Warn, allow partial payment

---

## UC-03: Cash Book Report

**Actor:** Thủ quỹ (Cashier), Kế toán trưởng
**Trigger:** End of day / On-demand

### Happy Path
1. User navigates to Reports → Cash Book
2. Selects currency (VND/USD/EUR), date range
3. System displays: opening balance + all transactions (receipts/payments) + closing balance
4. Running balance after each line
5. Closing balance matches GL account 111x balance
6. User exports PDF or Excel

### Exception Paths
- **E1: Balance mismatch:** System detects discrepancy between cash book and GL → flag for reconciliation
- **E2: No transactions:** Show zero balance report

---

## UC-04: Bank Transfer

**Actor:** Kế toán thanh toán
**Trigger:** Need to withdraw cash from bank / deposit cash to bank

### Happy Path (Bank Withdrawal)
1. Accountant selects "Create Transfer"
2. Transfer type = "Bank Withdrawal"
3. From: 1121 (bank), To: 1111 (cash)
4. Enters amount, reason
5. System auto-creates: CashReceipt (Nợ 111/Có 112) and links transfer
6. Posts both: GL entry
7. Updates cash book + bank balance

### Alternative Paths
- **A1: Bank deposit:** Transfer type = "Bank Deposit", From: 1111, To: 1121
- **A2: Currency conversion:** From: 1111 VND, To: 1112 USD (with exchange rate)

---

## UC-05: Petty Cash Management

**Actor:** Kế toán trưởng, Thủ quỹ
**Trigger:** Need to establish/manage petty cash fund

### Happy Path (Establish Fund)
1. Chief Accountant creates petty cash fund
2. Assigns custodian (employee), initial amount, purpose
3. System creates fund with initial balance
4. Cash transferred to custodian: CashPayment (Nợ 141/Có 111)

### Happy Path (Settlement)
1. Custodian submits expense receipts
2. Accountant reviews, approves
3. System settles: adjusts fund balance
4. If amount < initial → top-up created
5. If fund closed → remaining balance returned (CashReceipt Nợ 111/Có 141)

---

## UC-06: Cash Inventory

**Actor:** Ban kiểm kê (Inventory committee: accountant + cashier)
**Trigger:** Period-end / Surprise / Cashier handover

### Happy Path
1. Committee initiates inventory
2. Cashier prepares: records all receipts/payments, calculates book balance
3. Committee counts cash by denomination
4. Enters denomination details into system
5. System calculates: actual total, compares with book
6. Result: balanced (difference = 0)
7. All parties sign, inventory completed

### Exception Paths
- **E1: Cash excess (actual > book):** Record as Debit 111/Credit 338 (waiting for owner)
- **E2: Cash shortage (actual < book):** Investigate → Debit 138 (receivable from responsible) or Debit 334 (deduct from salary)
- **E3: Significant discrepancy:** Escalate to Giám đốc for decision

---

## UC-07: Cash Flow Statement (B03-DN)

**Actor:** Kế toán trưởng
**Trigger:** Period-end reporting

### Happy Path
1. User navigates to Reports → Cash Flow Statement
2. Selects period (month/quarter/year)
3. Selects method (Direct / Indirect)
4. System calculates:
   - Cash flows from operating activities
   - Cash flows from investing activities
   - Cash flows from financing activities
5. Net increase/decrease in cash = sum
6. Opening cash + Net = Closing cash
7. Verifies: Closing cash = GL 111+112+113 balance
8. Exports B03-DN format

---

## UC-08: Foreign Currency Cash Revaluation

**Actor:** Kế toán trưởng
**Trigger:** Period-end (monthly/quarterly)

### Happy Path
1. User selects "Revaluate Foreign Currency Cash"
2. System gets period-end exchange rates
3. Calculates: Book balance vs Market rate balance
4. FX Gain → Credit 515
5. FX Loss → Debit 635
6. Creates adjustment JournalEntry
7. Updates cash book
