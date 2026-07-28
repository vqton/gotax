# Bank Module — Workflows

---

## WF-1: Monthly Bank Closing Workflow

```
┌─────────────────────────────────────────────────────────────────────┐
│  Accountant                          Chief Accountant               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  [1] Import statement (CSV/MT940)                                   │
│       │                                                             │
│       ▼                                                             │
│  [2] System validates & stores lines                                │
│       │                                                             │
│       ▼                                                             │
│  [3] Run auto-match                                                  │
│       │                                                             │
│       ▼                                                             │
│  [4] Review matches (green=auto, yellow=needs review)               │
│       │                                                             │
│       ▼                                                             │
│  [5] Manual match remaining items ────► [6] Review reconciliation   │
│       │                                      │                      │
│       ▼                                      ▼                      │
│  [7] Handle exceptions                      [8] Approve?           │
│   ├─ Write-off (bank fee)                   ├─ Yes ─► [10] Complete │
│   ├─ Missing entry → create                 └─ No ──► [9] Return    │
│   └─ Investigate → flag                         │                   │
│       │                                         ▼                   │
│       └──────────────────────────────────── [5] Revise              │
│                                                                     │
│  [11] Print S08-DN (Bank Ledger)                                    │
│  [12] Print S09-DN (Reconciliation Report)                          │
│  [13] Archive reconciliation                                        │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## WF-2: Payment Order Lifecycle

```
┌─────────────────────────────────────────────────────────────────────┐
│  Maker (Accountant)               Checker (Chief Accountant)        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  [1] Create payment order (or batch)                                │
│       │                                                             │
│       ▼                                                             │
│  [2] Validate: balance, duplicates, limits                          │
│       │                                                             │
│       ▼                                                             │
│  [3] Submit for approval ──────────► [4] Review payment             │
│       │                                      │                      │
│       ▼                                      ▼                      │
│  [5] Wait...                              [6] Approve?             │
│                                            ├─ Yes ─► [7] Approved   │
│                                            └─ No ──► [8] Rejected   │
│                                                      │              │
│       ◄──────────────────────────────────────────────┘              │
│       │                                                             │
│       ▼                                                             │
│  [9] Print UNC form (bank template)                                 │
│       │                                                             │
│       ▼                                                             │
│ [10] Submit to bank (manual/electronic)                             │
│       │                                                             │
│       ▼                                                             │
│ [11] Confirm payment in system                                      │
│       │                                                             │
│       ▼                                                             │
│ [12] System posts GL entry                                          │
│       │                                                             │
│       ▼                                                             │
│ [13] Payment marked CONFIRMED                                       │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## WF-3: Auto-Match Engine Flow

```
┌─────────────────────────────────────────────────────┐
│  Input: BankStatementLine + GL Transactions         │
│  Output: Matches + Confidence Scores                 │
├─────────────────────────────────────────────────────┤
│                                                     │
│  [1] Load statement lines for period                │
│  [2] Load GL entries for bank account (period)      │
│       │                                             │
│       ▼                                             │
│  [3] Apply matching rules sequentially              │
│       │                                             │
│       ├── R1: Exact amount + same date              │
│       │   └─ Match if |stmt.amount - gl.amount| = 0 │
│       │       AND stmt.date = gl.posted_date        │
│       │                                             │
│       ├── R2: Exact amount + 1 day tolerance        │
│       │   └─ Match if |date_diff| ≤ 1               │
│       │                                             │
│       ├── R3: Reference match                       │
│       │   └─ Match if bank_ref == gl_voucher_no     │
│       │                                             │
│       ├── R4: Fuzzy reference match                 │
│       │   └─ Match if bank_ref CONTAINS gl_ref      │
│       │       OR gl_ref CONTAINS bank_ref           │
│       │                                             │
│       ├── R5: Counterparty + amount                 │
│       │   └─ Match if fuzzy(name) AND amt match     │
│       │                                             │
│       └── R6: Rounding tolerance                    │
│           └─ Match if |diff| < 500                   │
│                                                     │
│  [4] Deduplicate: keep highest confidence match     │
│  [5] Return match results                           │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## WF-4: Loan Lifecycle

```
┌─────────────────────────────────────────────────────┐
│  Treasury Manager                                    │
├─────────────────────────────────────────────────────┤
│                                                     │
│  [1] Create loan agreement                           │
│   ├─ Type: Short-term / Long-term / Overdraft        │
│   ├─ Amount, interest rate, term, repayment method   │
│   └─ Status: ACTIVE                                  │
│                                                     │
│  [2] Disbursement                                    │
│   ├─ Record disbursement date & amount               │
│   ├─ System posts: Nợ 112 / Có 341                   │
│   └─ Updates outstanding_balance                     │
│                                                     │
│  [3] Repayment Schedule (auto-generated)             │
│   ├─ Annuity: fixed payment each period              │
│   ├─ Straight-line: fixed principal + variable int   │
│   └─ Bullet: interest only + principal at maturity   │
│                                                     │
│  [4] Monthly Repayment (repeating)                   │
│   ├─ Create payment order for installment            │
│   ├─ System posts: Nợ 341 + 635 / Có 112             │
│   └─ Updates outstanding_balance                     │
│       │                                              │
│       ▼                                              │
│  [5] Maturity                                         │
│   ├─ If balance = 0 → FULLY_PAID                     │
│   └─ If balance > 0 → OVERDUE flag                   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## WF-5: Term Deposit Lifecycle

```
┌─────────────────────────────────────────────────────┐
│  Treasury Manager                                    │
├─────────────────────────────────────────────────────┤
│                                                     │
│  [1] Create term deposit                             │
│   ├─ Amount, term (30/60/90/180/365 days), rate     │
│   ├─ Auto-renewal flag                               │
│   └─ System posts: Nợ 1281 / Có 112                  │
│                                                     │
│  [2] Monitor (holding period)                        │
│   ├─ No transactions during term                     │
│   └─ Interest accrues (off-system, calculated)       │
│                                                     │
│  [3] Maturity Date                                   │
│   ├─ System calculates interest:                     │
│   │   Interest = Amount × Rate × Term / 365          │
│   ├─ Process maturity:                               │
│   │   Nợ 112 / Có 1281 (principal)                   │
│   │   Nợ 112 / Có 515 (interest)                     │
│   └─ Status → MATURED                                │
│       │                                              │
│       ▼                                              │
│  [4] Auto-renewal?                                   │
│   ├─ Yes → Create new deposit with same amount       │
│   │        + current interest rate                   │
│   │        Status → ACTIVE                           │
│   └─ No → Done (CLOSED)                              │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## WF-6: FCY Bank Account Revaluation (Period-End)

```
┌─────────────────────────────────────────────────────┐
│  Chief Accountant                                    │
├─────────────────────────────────────────────────────┤
│                                                     │
│  [1] Navigate to Period-End → FX Revaluation          │
│       │                                              │
│       ▼                                              │
│  [2] Select FCY bank accounts (1122)                  │
│       │                                              │
│       ▼                                              │
│  [3] System fetches exchange rate                     │
│   ├─ Source: configured bank's average transfer rate  │
│   └─ Prompt if rate unavailable                      │
│       │                                              │
│       ▼                                              │
│  [4] System calculates gain/loss                      │
│   ├─ GL balance at book rate                          │
│   ├─ GL balance at current rate                       │
│   └─ Difference → gain (515) or loss (635)            │
│       │                                              │
│       ▼                                              │
│  [5] Preview revaluation entries                      │
│       │                                              │
│       ▼                                              │
│  [6] Confirm → System posts                           │
│   ├─ Gain: Nợ 1122 / Có 515                          │
│   └─ Loss: Nợ 635 / Có 1122                          │
│       │                                              │
│       ▼                                              │
│  [7] Log revaluation history                          │
│                                                     │
└─────────────────────────────────────────────────────┘
```
