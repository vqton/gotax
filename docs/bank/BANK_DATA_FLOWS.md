# Bank Module — Data Flows

---

## DF-1: Bank Statement Import

```
┌─────────┐    ┌──────────────┐    ┌───────────────┐    ┌──────────────┐
│  Bank   │    │  CSV/MT940   │    │  Import       │    │  PostgreSQL  │
│  Portal │────▶│  File        │────▶│  Engine       │────▶│  Database    │
└─────────┘    └──────────────┘    └───────────────┘    └──────────────┘
                                        │
                                        ▼
                                  ┌──────────────┐    ┌──────────────┐
                                  │  Parser      │────▶│  Validator   │
                                  └──────────────┘    └──────────────┘
                                        │                    │
                                        ▼                    ▼
                                  ┌──────────────┐    ┌──────────────┐
                                  │  Line        │    │  Duplicate   │
                                  │  Normalizer  │    │  Detector    │
                                  └──────────────┘    └──────────────┘
                                        │                    │
                                        └────────┬───────────┘
                                                 │
                                                 ▼
                                          ┌──────────────┐
                                          │  Write to    │
                                          │  DB + Log    │
                                          └──────────────┘

Data:
- Raw file → bank_statements (header metadata)
- Each line → bank_statement_lines (normalized)
- Raw line → bank_statement_lines.raw_data for audit
```

## DF-2: Auto-Reconciliation

```
┌────────────────┐    ┌──────────────┐    ┌────────────────┐
│  bank_statement│    │  Match       │    │  GL Journal    │
│  _lines        │────▶│  Engine      │◀────│  Entries       │
└────────────────┘    └──────────────┘    └────────────────┘
       │                      │
       ▼                      ▼
┌────────────────┐    ┌──────────────┐
│  Balance       │    │  Rule        │
│  Verifier      │    │  Processor   │
└────────────────┘    └──────────────┘
       │                      │
       └──────────┬───────────┘
                  │
                  ▼
          ┌───────────────┐    ┌────────────────┐
          │  Match        │────▶│  bank_recon-   │
          │  Deduplicator │    │  ciliation_    │
          └───────────────┘    │  matches       │
                               └────────────────┘
                  │
                  ▼
          ┌───────────────┐
          │  Notify User  │
          │  (matches)    │
          └───────────────┘

Search scope for match engine:
- GL entries: journal_entries + cash_receipts + cash_payments + payment_orders
- Date range: statement period ± 2 days
- Amount: exact match first, then tolerance
```

## DF-3: Payment Order → GL Posting

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  User        │    │  Payment     │    │  Approval    │
│  Input       │────▶│  Validation  │────▶│  Workflow    │
└──────────────┘    └──────────────┘    └──────────────┘
                        │                      │
                        ▼                      ▼
                  ┌──────────────┐    ┌──────────────┐
                  │  Balance     │    │  Auth Check  │
                  │  Check       │    │  (maker≠     │
                  └──────────────┘    │  checker)    │
                                      └──────────────┘
                        │
                        ▼
                  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
                  │  GL Posting  │────▶│  Journal     │────▶│  Update      │
                  │  Engine      │    │  Entry       │    │  PO Status   │
                  └──────────────┘    │  (created)   │    │  → CONFIRMED │
                                      └──────────────┘    └──────────────┘
                        │
                        ▼
                  ┌──────────────┐
                  │  Response    │
                  │  to Client   │
                  └──────────────┘

GL Entry created:
Debit: 331 (AP) / 642 (expense) / 334 (salary) / 333 (tax)
Credit: 112x (bank account)
```

## DF-4: Loan Disbursement Flow

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Loan        │    │  Disbursement│    │  Balance     │
│  Agreement   │────▶│  Record      │────▶│  Update      │
└──────────────┘    └──────────────┘    └──────────────┘
                        │                      │
                        ▼                      ▼
                  ┌──────────────┐    ┌──────────────┐
                  │  GL Posting  │    │  Repayment   │
                  │  (112/341)   │    │  Schedule    │
                  └──────────────┘    │  Generate    │
                                      └──────────────┘

Data updates:
- loan_agreements.disbursed_amount += amount
- loan_agreements.outstanding_balance += amount
- loan_disbursements (new row)
- journal_entry: Nợ 112 / Có 341
- repayment_schedule generated if first disbursement
```

## DF-5: Term Deposit Maturity

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  System      │    │  Interest    │    │  Auto-       │
│  Cron/Daily  │────▶│  Calculator  │────▶│  Renewal     │
└──────────────┘    └──────────────┘    │  Check       │
                                         └──────────────┘
                        │                      │
                        ▼                      ▼
                  ┌──────────────┐    ┌──────────────┐
                  │  Status      │    │  Renewal     │
                  │  MATURED     │    │  Processing  │
                  └──────────────┘    └──────────────┘
                        │
                        ▼
                  ┌──────────────┐
                  │  GL Posting  │
                  │  112/1281/515 │
                  └──────────────┘

Auto-renewal posting:
1. Maturity: Nợ 112 / Có 1281 (principal) + Nợ 112 / Có 515 (interest)
2. New deposit (if auto-renewal): Nợ 1281 / Có 112 (new amount)
```

## DF-6: FX Revaluation at Period-End

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Fetch       │    │  Get Bank    │    │  Calculate   │
│  FCY Bank    │────▶│  TTM Rate    │────▶│  Gain/Loss   │
│  Accounts    │    │  (Circular   │    │  per Account │
└──────────────┘    │  99 Art. 6)  │    └──────────────┘
                    └──────────────┘         │
                                             ▼
                                       ┌──────────────┐
                                       │  Materiality │
                                       │  Check       │
                                       └──────────────┘
                                             │
                                             ▼
                                       ┌──────────────┐
                                       │  Post GL     │
                                       │  Revaluation │
                                       └──────────────┘

Calculation:
Book balance in VND = FCY amount × book rate
Current balance in VND = FCY amount × current TTM rate
Gain/Loss = Current - Book

Materiality: if |Gain/Loss| < 100,000 VND → skip (configurable)
```

## DF-7: Data Model Relationships

```
Company 1──N CompanyBankAccount 1──N BankStatement 1──N BankStatementLine
                                          │
                                          │ 1
                                          ▼
Company 1──N BankReconciliation 1──N BankReconciliationMatch
                                          │
                            ┌─────────────┴─────────────┐
                            │                           │
                            ▼                           ▼
                      BankStatementLine          GL Journal Entry
                                                 Cash Receipt
                                                 Cash Payment
                                                 Payment Order
                                                 AR Receipt
                                                 AP Payment

Company 1──N PaymentOrder N──1 CompanyBankAccount
Company 1──N PaymentOrderBatch N──N PaymentOrder

Company 1──N LoanAgreement 1──N LoanDisbursement
                       1──N LoanRepayment N──1 PaymentOrder

Company 1──N TermDeposit N──1 CompanyBankAccount
```
