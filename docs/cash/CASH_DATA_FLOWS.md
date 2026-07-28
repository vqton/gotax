# Cash Module — Data Flows

---

## 1. Cash Receipt Data Flow

```
User                    Handler               Service                Repository           GL
 │                        │                      │                      │                   │
 │ POST /receipts         │                      │                      │                   │
 │───────────────────────>│                      │                      │                   │
 │                        │ Validate request     │                      │                   │
 │                        │─────────────────┐    │                      │                   │
 │                        │ Validate rules   │    │                      │                   │
 │                        │<────────────────┘    │                      │                   │
 │                        │ CreateReceipt(dto)   │                      │                   │
 │                        │─────────────────────>│                      │                   │
 │                        │                      │ Validate business    │                   │
 │                        │                      │ rules                │                   │
 │                        │                      │ Generate voucher #   │                   │
 │                        │                      │─────────────────┐    │                   │
 │                        │                      │ CreateReceipt() │    │                   │
 │                        │                      │<────────────────┘    │                   │
 │                        │                      │ CreateReceipt        │                   │
 │                        │                      │─────────────────────>│                   │
 │                        │                      │                      │ INSERT cash_      │
 │                        │                      │                      │ receipts          │
 │                        │                      │                      │─────────────────┐ │
 │                        │                      │ Success              │                 │ │
 │                        │                      │<─────────────────────│<────────────────┘ │
 │                        │ Return receipt       │                      │                   │
 │                        │<─────────────────────│                      │                   │
 │ Receipt created        │                      │                      │                   │
 │<───────────────────────│                      │                      │                   │
                                                                                              
                                                                                              
User                    Handler               Service                Repository           GL
 │ POST /receipts/:id    │                      │                      │                   │
 │ /post                 │                      │                      │                   │
 │───────────────────────>│                      │                      │                   │
 │                        │ PostReceipt(id)      │                      │                   │
 │                        │─────────────────────>│                      │                   │
 │                        │                      │ GetReceipt(id)       │                   │
 │                        │                      │─────────────────────>│                   │
 │                        │                      │<─ Receipt            │                   │
 │                        │                      │                      │                   │
 │                        │                      │ Validate: approved   │                   │
 │                        │                      │ CreateJournalEntry() │                   │
 │                        │                      │─────────────────────────────────────────>│
 │                        │                      │                      │                   │
 │                        │                      │                      │ INSERT journal_   │
 │                        │                      │                      │ entries + lines   │
 │                        │                      │                      │<──────────────────│
 │                        │                      │                      │                   │
 │                        │                      │ UpdateReceipt()      │                   │
 │                        │                      │ (status=posted,      │                   │
 │                        │                      │  gl_journal_id)      │                   │
 │                        │                      │─────────────────────>│                   │
 │                        │                      │                      │ UPDATE status     │
 │                        │                      │                      │─────────────────┐ │
 │                        │                      │ UpdateCashBook()     │                 │ │
 │                        │                      │─────────────────────>│                 │ │
 │                        │                      │                      │ Recalc balance   │ │
 │                        │                      │ Success              │<────────────────┘ │
 │ Receipt posted         │<─────────────────────│<─────────────────────│                   │
│<───────────────────────│                      │                      │                   │
```

---

## 2. Cash Payment Data Flow

```
Similar to receipt but with additional:
- Balance sufficiency check (optional)
- Approval threshold check
  - If amount > config → route to Director
  - Else → route to Chief Accountant
- WHT calculation (if applicable)
```

---

## 3. Cash Book Report Data Flow

```
User                    Handler               Service                Repository
 │                        │                      │                      │
 │ GET /reports/          │                      │                      │
 │ cash-book?currency=    │                      │                      │
 │ VND&from=2026-01-01    │                      │                      │
 │ &to=2026-01-31         │                      │                      │
 │───────────────────────>│                      │                      │
 │                        │ GetCashBook(filter)  │                      │
 │                        │─────────────────────>│                      │
 │                        │                      │ GetOpeningBalance()  │
 │                        │                      │─────────────────────>│
 │                        │                      │<─ opening_balance    │
 │                        │                      │                      │
 │                        │                      │ GetTransactions()    │
 │                        │                      │ (receipts+payments   │
 │                        │                      │  in date range)      │
 │                        │                      │─────────────────────>│
 │                        │                      │<─ []receipts,         │
 │                        │                      │    []payments        │
 │                        │                      │                      │
 │                        │                      │ Build cash book:     │
 │                        │                      │ - Sort by date       │
 │                        │                      │ - Running balance    │
 │                        │                      │ - Calculate totals   │
 │                        │                      │─────────────────┐    │
 │                        │                      │ Verify: closing  │    │
 │                        │                      │ = opening +      │    │
 │                        │                      │ receipts -       │    │
 │                        │                      │ payments         │    │
 │                        │ Return CashBook      │<────────────────┘    │
 │                        │<─────────────────────│                      │
 │ Cash book report       │                      │                      │
 │<───────────────────────│                      │                      │
```

---

## 4. Cash Inventory Data Flow

```
User                    Handler               Service                Repository
 │                        │                      │                      │
 │ POST /cash/inventory   │                      │                      │
 │ {cash_account_id,      │                      │                      │
 │  denominations:        │                      │                      │
 │  [{denom, count},...]} │                      │                      │
 │───────────────────────>│                      │                      │
 │                        │ CreateInventory()    │                      │
 │                        │─────────────────────>│                      │
 │                        │                      │ GetBalance()         │
 │                        │                      │ (book_balance)       │
 │                        │                      │─────────────────────>│
 │                        │                      │<─ book_balance       │
 │                        │                      │                      │
 │                        │                      │ Calculate actual:    │
 │                        │                      │ Σ(denom × count)     │
 │                        │                      │                      │
 │                        │                      │ diff = actual - book │
 │                        │                      │                      │
 │                        │                      │ if diff ≠ 0:        │
 │                        │                      │   determine type     │
 │                        │                      │   (excess/shortage)  │
 │                        │                      │                      │
 │                        │                      │ StoreInventory()     │
 │                        │                      │─────────────────────>│
 │                        │ Return inventory     │                      │
 │                        │<─────────────────────│                      │
 │ Inventory result       │                      │                      │
 │<───────────────────────│                      │                      │
```

---

## 5. Cash Flow Statement (B03-DN) Data Flow

```
User                    Handler               Service                Repository
 │                        │                      │                      │
 │ GET /reports/          │                      │                      │
 │ cash-flow?period=      │                      │                      │
 │ 2026-Q1&method=direct  │                      │                      │
 │───────────────────────>│                      │                      │
 │                        │ GetCashFlow(filter)  │                      │
 │                        │─────────────────────>│                      │
 │                        │                      │ ┌─────────────────┐  │
 │                        │                      │ │ GL Account      │  │
 │                        │                      │ │ Mapping:        │  │
 │                        │                      │ │ Operating:      │  │
 │                        │                      │ │  511, 515,      │  │
 │                        │                      │ │  632, 641, 642  │  │
 │                        │                      │ │ Investing:      │  │
 │                        │                      │ │  211, 221, 222  │  │
 │                        │                      │ │ Financing:      │  │
 │                        │                      │ │  341, 411, 421  │  │
 │                        │                      │ └─────────────────┘  │
 │                        │                      │                      │
 │                        │                      │ Aggregate GL data   │
 │                        │                      │ by cash flow        │
 │                        │                      │ category            │
 │                        │                      │─────────────────────>│
 │                        │                      │<─ aggregated data    │
 │                        │                      │                      │
 │                        │                      │ Build B03-DN report  │
 │                        │                      │ Calculate:           │
 │                        │                      │ Opening cash         │
 │                        │                      │ + Operating CF       │
 │                        │                      │ + Investing CF       │
 │                        │                      │ + Financing CF       │
 │                        │                      │ = Closing cash       │
 │                        │                      │─────────────────┐    │
 │                        │                      │ Verify: closing │    │
 │                        │                      │ = GL 111+112    │    │
 │                        │                      │ +113 balance    │    │
 │                        │ Return B03-DN        │<────────────────┘    │
 │                        │<─────────────────────│                      │
 │ B03-DN report          │                      │                      │
 │<───────────────────────│                      │                      │
```

---

## 6. Multi-Currency Revaluation Data Flow

```
Period-end:
1. Get all foreign currency cash balances (1112, 1122)
2. Get period-end exchange rates from configured source
3. For each currency:
   Calculate: revalued_amount = original_fcy × period_end_rate
   diff = revalued_amount - book_amount_vnd
4. If diff > 0 (gain): Credit 515 (Financial income)
   If diff < 0 (loss): Debit 635 (Financial expense)
5. Create adjusting JournalEntry
6. Update cash book with adjustment
```

---

## 7. System Integration Data Flow

```
AR Module ──────> Cash Receipt (auto-create from AR payment)
                    │
AP Module ──────> Cash Payment (auto-create from AP payment)
                    │
Bank Module ────> Cash Transfer (bank withdrawal/deposit)
                    │
                    ▼
              GL Posting
                    │
                    ▼
           ┌────────────────┐
           │ Cash Book      │
           │ (Real-time)    │
           └────────┬───────┘
                    │
         ┌──────────┴──────────┐
         ▼                     ▼
┌────────────────┐  ┌────────────────┐
│ B03-DN         │  │ Tax Reports    │
│ Cash Flow      │  │ (non-cash      │
│ Statement      │  │  payment       │
│                │  │  thresholds)   │
└────────────────┘  └────────────────┘
```
