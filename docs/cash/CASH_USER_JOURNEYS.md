# Cash Module — User Journeys

---

## Journey 1: Daily Cashier (Thủ quỹ)

**Role:** Thủ quỹ — manages physical cash, maintains cash book
**Frequency:** Daily

```
07:30 — Arrive, check opening balance
         Login → Cash Dashboard → View opening balance
         
08:00 — Process customer cash payment
         Cashier receives notification: new approved Phiếu thu
         Verify: identity of payer, amount matches voucher
         Count cash → Confirm receipt → Sign voucher
         System updates cash book: +amount

10:30 — Process supplier payment
         Cashier receives notification: new approved Phiếu chi
         Verify: identity of payee, amount
         Count cash → Pay → Get signature → Confirm
         System updates cash book: -amount

12:00 — Lunch break
         Cash locked in safe

14:00 — Bank withdrawal
         Giám đốc approved: rút 50M từ NH
         Cashier goes to bank → withdraws → returns
         Confirms receipt: +50M in cash book

16:30 — End-of-day cash count
         Count physical cash by denomination
         Enter count into system
         System compares: book vs actual
         
17:00 — If match:
         Sign daily cash book → Print → File
         If mismatch:
         Recount → If still mismatch → Escalate

17:30 — Lock cash in safe → Logout
```

**Pain points:**
- Manual counting error-prone
- Book-actual discrepancy needs quick resolution
- No denomination breakdown in current system

**Delight:**
- Real-time balance visibility
- Auto-generated cash book
- Quick discrepancy detection

---

## Journey 2: Chief Accountant (Kế toán trưởng)

**Role:** Kế toán trưởng — approves, reviews, reports
**Frequency:** Daily/Weekly/Monthly

```
08:00 — Review pending approvals
         Dashboard: 12 pending receipts, 5 pending payments
         Review each: amount, reason, attachments
         Approve (< 10M): batch approve
         Review closely (> 10M): check supporting docs

09:30 — Handle large payment request
         Payment: 150M to supplier XYZ
         Check: contract, invoice, delivery receipt
         Forward to Giám đốc for final approval

14:00 — Monthly review
         Open Cash Flow Report → Select current month
         Check: cash balance trend, large transactions
         Export for management meeting

16:00 — Audit preparation
         Open Inventory Report → Last 3 months
         Check: all inventories completed, no unresolved discrepancies
         Sign off

Monthly:
- Review B03-DN Cash Flow Statement
- Approve foreign currency revaluation
- Review bank reconciliation
- Sign monthly cash book
```

**Pain points:**
- Too many small approvals
- Difficult to spot anomalies across large dataset
- Manual reconciliation between cash book and GL

**Delight:**
- Approval dashboard with filters
- Anomaly detection alerts
- Auto-GL reconciliation

---

## Journey 3: AP Accountant (Kế toán thanh toán)

**Role:** Creates payment vouchers
**Frequency:** Daily

```
09:00 — Process supplier invoices due
         Open AP Aging → Select due invoices
         For each: verify goods receipt, contract, invoice
         Select "Pay by Cash" → Create Phiếu chi
         System auto-fills: supplier, amount, account 331
         Assign correct expense account (debit)
         Attach invoice PDFs
         Submit for approval

11:00 — Handle urgent payment
         Đột xuất: pay customs tax 15M in cash
         Create payment: type=Tax, account 333
         Submit with priority flag

14:00 — Petty cash settlement
         Employee returns from business trip
         Review expenses: hotel (5M), transport (2M), meals (1.5M)
         Total: 8.5M, advance was 10M → return 1.5M
         Create CashReceipt for refund
         Link to petty cash fund
```

**Pain points:**
- Manual entry of repetitive suppliers
- Attachment management (scanning, uploading)

**Delight:**
- Auto-suggest account codes from past transactions
- Bulk payment creation
- Integration with AP module

---

## Journey 4: Auditor (Kiểm toán)

**Role:** External/Internal auditor
**Frequency:** Quarterly/Annually

```
Day 1 — Request cash documents
         Request: cash book (S07-DN), all vouchers month 1-6
         System generates: PDF cash book + voucher list
         Sample selection: 20 receipts, 20 payments

Day 2 — Vouchers testing
         For each sample:
         - Check: voucher number sequence (no gaps)
         - Check: signatures present (creator, approver, cashier, receiver)
         - Check: attachment completeness
         - Check: GL posting correct
         System provides audit trail

Day 3 — Cash inventory verification
         Surprise cash count
         System generates: inventory report
         Compare: audited count vs system balance
         Verify discrepancy handling

Day 4 — Cash flow statement review
         B03-DN: verify direct/indirect method
         Check: opening + net cash flow = closing
         Cross-check: closing = GL 111+112+113
```

**Delight:**
- Full digital audit trail
- Voucher sequence completeness check
- Instant report generation
