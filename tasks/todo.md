# AR Module — Task Tracker v2

## ✅ Phase 0: Critical Fix (COMPLETE)

### Task 1: Fix AR Aging Report — Bucket by DueDate

**Desc:** `GetARAgingReport` puts all amount in Bucket0. Fix: query invoices by DueDate, calc days overdue, group into proper buckets.

**AC:**
- [x] Aging report shows correct bucket per invoice BalanceDue
- [x] Current (<0 days overdue) in Bucket0
- [x] 1-30 days overdue in Bucket30, 31-60→Bucket60, 61-90→Bucket90, 91-120→Bucket120, 120+→Bucket120
- [x] Paid invoices excluded
- [x] 3 tests: validates buckets, no invoices, paid excluded

**Verification:** `go test -v -run TestARAging ./internal/service/`

**Files:** `internal/service/sale_service.go`, `internal/service/ar_service_test.go`

---

## ✅ Phase 1: Core AR (COMPLETE)

### Task 2: GL Auto-Posting on Invoice Post

**Desc:** PostInvoice creates JournalEntry via GL.CreatePostedEntry: Dr 131 (total), Cr revenue accounts (subtotal by line), Cr 3331 (VAT). Set gl_posted=true.

**AC:**
- [x] PostInvoice creates JournalEntry via GL in service layer (not in repo)
- [x] GL entry: Dr 131 (AR), Cr revenue account(s), Cr 3331 (VAT)
- [x] `gl_posted=true`, `gl_posted_at` set
- [x] nil GL → no crash, GLPosted stays false
- [x] Test: posted invoice has matching GL entry

**Files:** `internal/service/sale_service.go`, `internal/service/service.go` (CreatePostedEntry), `internal/domain/interfaces.go` (Set*GLPosted methods), `internal/repository/*`

---

### Task 3: GL Auto-Posting on Receipt Post

**Desc:** PostReceipt creates JournalEntry: Dr 1111 (cash) or 1121 (bank), Cr 131 (AR) by amount.

**AC:**
- [x] PostReceipt creates JournalEntry: Dr cash/bank, Cr 131
- [x] Payment method mapping: cash→1111, bank_transfer/cheque/credit_card→1121
- [x] `gl_posted=true`, `gl_posted_at` set
- [x] Test: bank transfer Dr 1121, cash receipt Dr 1111

**Files:** `internal/service/sale_service.go`, `internal/domain/models_sale.go` (added GLPosted to Receipt)

---

### Task 4: GL Auto-Posting on Credit Note Post

**Desc:** PostCN creates reversing JournalEntry: Dr revenue account (from original inv), Dr 3331 (VAT), Cr 131 (total).

**AC:**
- [x] PostCN creates JournalEntry: Dr 511X, Dr 3331, Cr 131
- [x] Revenue account resolved from original invoice lines
- [x] `gl_posted=true`, `gl_posted_at` set
- [x] Test: CN reduces AR balance

**Files:** `internal/service/sale_service.go`, `internal/repository/*` (stripped GLPosted from PostCN)

---

### Task 5: Customer Statement Endpoint

**Desc:** Generate customer statement: opening balance, all transactions (invoices/receipts/CNs) in period with running balance, closing balance.

**AC:**
- [x] GET `/api/v1/sale/ar/statement?customer_id=X` returns statement
- [x] All period transactions listed (invoices, receipts, credit notes)
- [x] Running balance per line item
- [x] Closing balance = opening + all debits - all credits
- [x] Handler test: statement matches invoice+receipt totals

**Files:** `internal/handler/sale_handler.go`, `internal/service/sale_service.go`, `internal/domain/models_sale.go` (CustomerStatement structs), `internal/handler/sale_handler_test.go`

---

## ⏸️ Phase 2: E-Invoice Pipeline (DEFERRED)

Deferred per PI direction. Full design in `docs/sale/AR_WORKFLOW.md`. Retained for reference.

### Task 6: E-Invoice Pipeline Interface + TXML Generator
### Task 7: E-Invoice Digital Signer
### Task 8: GDT API Client
### Task 9: E-Invoice Auto-Pipeline Hook

---

## 🔲 Phase 3: Collection Management (P1)

### Task 6 (new numbering): Prepayment/Deposit Workflow

**Desc:** Full prepayment lifecycle. Customer pays deposit → optional deposit e-invoice (deferred) → offset against final invoice. Extend CustomerReceipt with ReceiptType field.

**Foundations exist:**
- `ARTransPrepayment` type constant defined
- `CustomerReceipt` model exists with allocation support
- Receipt allocation/reconciliation workflow exists

**New:**
- `CustomerReceipt.ReceiptType` field — values: standard, prepayment, refund
- Receipt create validates: prepayment has no invoice allocation
- Deposit offset endpoint: apply prepayment to one or more invoices
- Refund workflow: receipt with type=refund, negative amount, links to original prepayment
- Refund check: cannot refund more than remaining deposit balance

**AC:**
- [ ] CustomerReceipt has ReceiptType (standard/prepayment/refund)
- [ ] CreateReceipt validates prepayment has no invoice allocation (unallocated)
- [ ] POST `/api/v1/sale/receipts/:id/offset` — apply prepayment to invoice
- [ ] Offset creates ARTransOffset + reduces invoice BalanceDue
- [ ] Refund receipt type: negative amount, links to original prepayment
- [ ] Test: prepayment offsets correctly reduce invoice balance

**Implementation order:**
1. Add ReceiptType to CustomerReceipt model + Validate()
2. Add OffsetReceipt to SaleService + repos
3. Add handler endpoint + route
4. Write tests

**Files:**
- `internal/domain/models_sale.go` — ReceiptType field, constants
- `internal/domain/interfaces.go` — maybe new offset methods
- `internal/repository/memory_sale.go` — OffsetReceipt impl
- `internal/repository/pg_sale.go` — OffsetReceipt impl
- `internal/service/sale_service.go` — OffsetReceipt service method
- `internal/handler/sale_handler.go` — offset endpoint
- `internal/handler/sale_handler_test.go` — handler test

**Scope:** M (4-5 files)

---

### Task 7: FX Revaluation

**Desc:** Realized FX at receipt allocation time + unrealized FX at month-end for outstanding foreign-currency AR.

**Foundations exist:**
- ExchangeRate CRUD exists in GL service
- CustomerInvoice.Currency + ExchangeRate fields exist
- CustomerReceipt.Currency + ExchangeRate fields exist
- Account 515 (Doanh thu tai chinh) and 635 (Chi phi tai chinh) exist

**Realized FX (at receipt allocation):**
- Compare invoice exchange rate vs receipt exchange rate
- If receipt rate > invoice rate → Dr 635 (FX loss) Cr 131 (reduce AR)
- If receipt rate < invoice rate → Dr 131 (increase AR) Cr 515 (FX gain)
- Post alongside the main receipt GL entry

**Unrealized FX (month-end batch):**
- Get all open (non-zero BalanceDue) foreign-currency invoices
- Compute reval at period-end rate from ExchangeRate table
- Post net adjustment: Dr/Cr 131, Cr/Dr 515/635
- Store reval journal entry reference for audit trail

**AC:**
- [ ] Receipt allocation calculates realized FX gain/loss when currencies match but rates differ
- [ ] Realized FX posted: Dr 635/Cr 131 (loss) or Dr 131/Cr 515 (gain)
- [ ] Month-end batch endpoint: revalue all open FC invoices
- [ ] Unrealized FX posted as net adjustment entry
- [ ] Test: FX gain/loss posted correctly for simple case

**Implementation order:**
1. Add FX calc to receipt allocation logic in sale_service.go
2. Add month-end batch endpoint
3. Write tests

**Files:**
- `internal/service/sale_service.go` — FX in allocation + month-end
- `internal/handler/sale_handler.go` — month-end endpoint
- `internal/service/ar_service_test.go` — FX tests
- `internal/handler/sale_handler_test.go` — handler test

**Scope:** M (4-5 files)

---

### Task 8: Dunning Engine

**Desc:** Multi-level dunning with configurable triggers, reminder generation, promise-to-pay tracking, dunning history.

**Foundations exist:**
- AR aging report calculates overdue buckets
- Customer model has email field for reminders
- Invoice and receipt statuses support lifecycle

**New models:** `internal/domain/models_collection.go`

```go
type DunningLevel string
const (
    DunningL1 DunningLevel = "L1" // 1-30 days overdue — gentle reminder
    DunningL2 DunningLevel = "L2" // 31-60 days — formal notice
    DunningL3 DunningLevel = "L3" // 61-90 days — final notice
    DunningL4 DunningLevel = "L4" // 90+ days — legal action warning
)

type DunningConfig struct {
    CompanyID       string
    Level           DunningLevel
    DaysOverdueMin  int
    DaysOverdueMax  int
    AmountThreshold float64        // 0 = no threshold
    Active          bool
}

type DunningQueue struct {
    ID              string
    CompanyID       string
    InvoiceID       string
    CustomerID      string
    CurrentLevel    DunningLevel
    LastReminderAt  *time.Time
    NextReminderAt  *time.Time
    PTPDate         *time.Time     // promise-to-pay date (suspends dunning)
    PTPNotes        string
    Status          DunningStatus  // pending/escalated/resolved/suspended
    CreatedAt       time.Time
    ResolvedAt      *time.Time
}

type DunningHistory struct {
    ID          string
    QueueID     string
    InvoiceID   string
    CustomerID  string
    Level       DunningLevel
    Action      string           // reminder_sent/ptp_recorded/escalated/resolved
    Notes       string
    CreatedAt   time.Time
}
```

**Service methods:**
- `GenerateDunningQueue(ctx, companyID)` — batch: scan overdue invoices, upsert dunning queue
- `SendReminders(ctx, companyID, level)` — generate reminders for queue items due
- `RecordPTP(ctx, queueID, date, notes)` — record promise-to-pay, suspend dunning
- `EscalateLevel(ctx, queueID)` — manually escalate dunning level
- `ResolveQueue(ctx, queueID)` — mark resolved (invoice paid)
- `GetDunningQueue(ctx, companyID, level)` — list queue items
- `GetDunningHistory(ctx, invoiceID)` — dunning history

**Dunning batch logic:**
1. Full scan of open (non-paid, non-cancelled) invoices past DueDate
2. Compute days overdue from report date
3. Match to DunningConfig by days range + threshold
4. Upsert DunningQueue (update level, next reminder date)
5. Auto-escalate: if current level < computed level, increment with history

**AC:**
- [ ] DunningConfig CRUD (create/list/update per company)
- [ ] GenerateDunningQueue: scans open invoices, creates/updates queue
- [ ] Correct level assignment based on days overdue
- [ ] Promise-to-pay records (suspends dunning for that invoice)
- [ ] Escalation history tracked in DunningHistory
- [ ] Queue resolved when invoice fully paid
- [ ] Test: dunning batch processes correct invoices at correct level

**Implementation order:**
1. Models + repo interface + memory repo
2. PG schema + PG repo
3. CollectionService + DunningConfig CRUD
4. GenerateDunningQueue batch + test
5. SendReminders (tracking only, no actual email/SMS)
6. PTP + escalation + resolution
7. Handler + routes
8. Handler tests

**Files:**
- `internal/domain/models_collection.go` — new file
- `internal/domain/interfaces.go` — DunningRepository interface
- `internal/repository/memory_sale.go` — MemoryDunningRepo
- `internal/repository/pg_sale.go` — PGDunningRepo
- `migrations/008_dunning.sql` — dunning tables
- `internal/service/collection_service.go` — new file
- `internal/handler/sale_handler.go` — add to SaleHandler
- `internal/service/ar_service_test.go` — dunning tests
- `internal/handler/sale_handler_test.go` — handler tests

**Scope:** L (8-10 files)

---

### Task 9: Bad Debt Provision + Write-Off

**Desc:** Provision calc per VAS 17 by aging bucket %. Write-off workflow with approval. Recovery tracking.

**Foundations exist:**
- Account 229 (Du phong phai thu kho doi) seeded in tests
- Account 642 (Chi phi quan ly) exists
- Account 711 (Doanh thu khac) exists
- AR aging report gives buckets
- Customer has CustomerStatus (ACTIVE/SUSPENDED/BLACKLISTED)

**Provision calc (VAS 17, R8):**
```
- 0-6 months overdue:     0% (no provision) or optional
- 6-12 months overdue:    30% of balance
- 1-2 years overdue:      50%
- 2-3 years overdue:      70%
- >3 years overdue:      100%
```
Configurable percentages per company. Batch calc at month-end.

**Write-off criteria:**
- Customer bankrupt/dissolved/liquidated
- Debtor deceased
- 3+ years overdue with full provision
- Legal proceedings concluded with no recovery
- Must be approved (segregation: accountant proposes, chief accountant approves)

**Write-off entry:** Dr 229 (provision) Cr 131 (AR) — if insufficient provision, Dr 642 for remainder.
Off-balance-sheet: record in monitoring account (Dr 009 "Debt written off" — per Circular 99).

**Recovery (unexpected payment after write-off):**
- Reinstate receivable: Dr 131 Cr 711
- Record payment: Dr bank Cr 131
- Remove from off-balance-sheet: Cr 009

**AC:**
- [ ] ProvisionConfig CRUD (bucket %, active toggle per company)
- [ ] `CalculateProvision(ctx, companyID, asOfDate)` — computes provision per customer
- [ ] `PostProvision(ctx, companyID, asOfDate)` — posts Dr 642 Cr 229
- [ ] `WriteOff(ctx, invoiceIDs, reason, approvedBy)` — writes off AR
- [ ] Write-off: Dr 229 Cr 131 (+ Dr 642 if insufficient), Cr 009 if tracked
- [ ] `RecoverWriteOff(ctx, invoiceID, amount)` — reinstatement + payment
- [ ] Test: provision calc matches manual calc

**Implementation order:**
1. Models + repo interface
2. PG schema + PG repo
3. Provision calc + post in CollectionService
4. Write-off + recovery in CollectionService
5. Handler + routes
6. Tests

**Files:**
- `internal/domain/models_collection.go` — ProvisionConfig, BadDebtWriteOff models
- `internal/domain/interfaces.go` — new repository interfaces
- `internal/repository/memory_sale.go` — memory impl
- `internal/repository/pg_sale.go` — PG impl
- `internal/service/collection_service.go` — provision/writeoff methods
- `internal/handler/sale_handler.go` — endpoints
- `internal/service/ar_service_test.go` — tests
- `internal/handler/sale_handler_test.go` — handler tests

**Scope:** L (6-8 files)

---

## 🔲 Phase 4: Controls & Monitoring (P1/P2)

### Task 10: Credit Limit Enforcement

**Desc:** Check at SO confirm + invoice create. Configurable action (warn vs block) per company.

**Foundations exist:**
- Customer.CreditLimit field exists (float64)
- SO lifecycle has Confirm → Approved → Confirmed transition
- Invoice Create checks customer exists

**Logic:**
```
outstandingAR = sum(BalanceDue) for customer's open invoices
if outstandingAR + newAmount > CreditLimit
    if action == "block" → return error
    if action == "warn" → add warning to response
```

**Override:** manager role can pass `?override=true` with reason to bypass limit.

**AC:**
- [ ] SO confirm checks credit limit (outstanding + SO total ≤ limit)
- [ ] Invoice create checks credit limit
- [ ] Config: warn vs block per customer or company
- [ ] Manager override: `?override=true&reason=...` bypasses block
- [ ] Test: exceeding limit returns configured action

**Files:**
- `internal/service/sale_service.go` — add check to ConfirmSO, CreateInvoice
- `internal/handler/sale_handler.go` — handle override param
- `internal/service/ar_service_test.go` — credit limit tests

**Scope:** S (2-3 files)

---

### Task 11: AR-GL Month-End Reconciliation

**Desc:** Report comparing AR sub-ledger total vs GL 131 balance. List variances by customer.

**Logic:**
```
subledger_total = sum(CustomerInvoice.BalanceDue) where status in (posted,paid)
gl_balance = GL account 131 ending balance for period
variance = gl_balance - subledger_total
if variance != 0 → drill-down by customer
```

**AC:**
- [ ] GET `/api/v1/sale/ar/reconciliation?period=X` returns report
- [ ] Report: sub-ledger total, GL 131 balance, variance
- [ ] Variance drill-down by customer
- [ ] Test: reconciled when GL matches sub-ledger (variance = 0)

**Files:**
- `internal/handler/sale_handler.go` — endpoint
- `internal/service/sale_service.go` — ARSubledgerTotal calc
- `internal/service/ar_service_test.go` — test

**Scope:** M (3-4 files)

---

### Task 12: DSO + AR KPI Dashboard

**Desc:** Days Sales Outstanding + key AR metrics.

**Metrics:**
- DSO = (Avg AR / Total Revenue) × days in period
- Avg AR = (opening AR + closing AR) / 2
- Total AR outstanding (current)
- Overdue % = overdue AR / total AR × 100
- Collection Rate = receipts / (invoices - credit notes) in period
- Avg Days to Pay = avg of payment date - invoice date

**AC:**
- [ ] DSO calc matches manual formula
- [ ] Total AR outstanding returned
- [ ] Overdue % correct
- [ ] Collection rate per period
- [ ] Test: DSO reasonable (~30-45 for standard terms)

**Files:**
- `internal/handler/sale_handler.go` — endpoint
- `internal/service/sale_service.go` — KPI calc methods
- `internal/service/ar_service_test.go` — tests

**Scope:** M (3-4 files)

---

### Task 13: Off-Balance-Sheet Tracking

**Desc:** Written-off AR tracked in monitoring account (account 009 "Debt written off") for 10 years per Accounting Law.

**Foundations exist:**
- Write-off posts Cr 009 (off-BS debit)
- Recovery removes from 009 (Cr 009)
- Query: total written-off, by year, by customer

**Note:** Off-balance-sheet accounts in Vietnamese accounting use "account 009" (tai khoan 009 "No kho kho doi da xu ly"). This is tracked as a separate table, not a GL account with Dr/Cr.

**Model:**
```go
type OffBalanceSheetItem struct {
    ID              string
    CompanyID       string
    InvoiceID       string
    CustomerID      string
    WrittenOffAt    time.Time
    WrittenOffAmount float64
    RecoveredAmount  float64
    RecoveredAt     *time.Time
    Status          string    // outstanding/recovered/partial
    Notes           string
}
```

**AC:**
- [ ] Write-off creates OffBalanceSheetItem
- [ ] Recovery updates item (recovered_amount, recovered_at)
- [ ] GET `/api/v1/sale/ar/written-off` lists off-balance-sheet items
- [ ] Query: filter by year, customer, status
- [ ] Test: write-off + recovery tracked correctly

**Files:**
- `internal/domain/models_collection.go` — OffBalanceSheetItem
- `internal/domain/interfaces.go` — OffBalanceSheetRepository
- `internal/repository/memory_sale.go` — memory impl
- `internal/repository/pg_sale.go` — PG impl
- `internal/handler/sale_handler.go` — endpoint
- `internal/service/collection_service.go` — methods
- `internal/service/ar_service_test.go` — tests

**Scope:** M (4-5 files)

---

## Progress

- [x] **Phase 0:** Fix AR Aging (COMPLETE)
- [x] **Phase 1:** Core AR (COMPLETE)
- [ ] **Phase 2:** E-Invoice Pipeline (DEFERRED)
- [ ] **Phase 3:** Collection Management (P1)
- [ ] **Phase 4:** Controls & Monitoring (P1/P2)

## Verification Gate

- [ ] All Phases 0-1: `go test ./... && go vet ./...` — PASS
- [ ] Each Phase 3/4 task: tests at service + handler level
- [ ] Every GL posting: Dr = Cr (double-entry invariant)
- [ ] No orphaned test data between runs
- [ ] Codegraph index up to date after implementation

---

## ✅ Purchase P2-5: GDT E-Invoice XML (COMPLETE — c3068c8, 8a484a5, 5074d93, 6895977)

### Slice 5.1 — einvoice package (Parse + Generate)

**Desc:** New pure package `internal/einvoice`. XML structs mirror PURCHASE_TEMPLATES §8. `Parse(raw) → *domain.SupplierInvoice`, `Generate(inv) → []byte`. GL defaults 152/1331, VAT rate→type map.

**AC:**
- [x] Parse valid VND invoice → SupplierInvoice (fields + lines + totals)
- [x] Parse FC invoice with ExchangeRate
- [x] Parse credit note (InvoiceType)
- [x] Parse error: malformed XML, missing required fields
- [x] Generate→Parse round-trip
- [x] VAT map: 0/5/8/10/-1

**Verification:** `go vet ./internal/einvoice && go test -count=1 ./internal/einvoice/`

### Slice 5.2 — Service wiring

**Desc:** `ReceiveEInvoiceXML(ctx, companyID, raw)` — parse → auto-create supplier → dedupe → create draft invoice → persist raw XML. `GenerateEInvoiceXML(ctx, id)` — load → XML.

**AC:**
- [x] ReceiveEInvoiceXML creates draft invoice with supplier auto-created
- [x] Raw XML stored in EInvoiceData
- [x] Duplicate invoice number rejected
- [x] GenerateEInvoiceXML returns parseable XML for posted invoice
- [x] Generate on missing invoice → ErrInvoiceNotFound

**Verification:** `go test -count=1 -run 'ReceiveEInvoiceXML|GenerateEInvoiceXML|EInvoice' ./internal/service/`

### Slice 5.3 — Handler + routes

**Desc:** `POST /api/v1/purchase/invoices/e-invoice` (raw XML body, company_id query) → 201. `GET /api/v1/purchase/invoices/:id/e-invoice` → 200 XML.

**AC:**
- [x] POST valid XML → 201, invoice returned, supplier auto-created
- [x] POST malformed XML → 400
- [x] GET → 200, content-type application/xml, round-trips through Parse
- [x] GET missing → 404

**Verification:** `go vet ./... && go test -count=1 ./...`

### Slice 5.4 — Docs + closeout

**Desc:** Readiness 0%→100%, AGENTS.md module table row, plan.md checkpoint.

**AC:**
- [x] PURCHASE_READINESS.md e-invoice row updated
- [x] PURCHASE_ANALYSIS_SUMMARY.md updated
- [x] AGENTS.md purchase row readiness updated
- [x] plan.md checkpoint P2-5 ✅, P2 closed

**Verification:** `git log` shows per-slice commits; final commit includes docs

---

## 🔲 Tax Core Foundations — Round A: Calculation + Declaration (NEXT)

### A1: Rate resolver
- [ ] `resolveRate(taxType, applicableTo, onDate)` helper in tax service — reads TaxRepository.GetRates, picks active rate for date
- [ ] Fallback: STANDARD rate if no match; error if no rate at all
- [ ] Test: active/inactive, effective date window, fallback
- [ ] `go vet ./... && go test -count=1 ./...`

### A2: VAT engine rewrite
- [ ] CalculateVAT: explicit 33311 (output) / 1331/1332 (input) accounts first
- [ ] Fallback: rate-table × revenue/expense for accounts without explicit VAT line
- [ ] 8% reduced rate honored via rate table (VAT-02)
- [ ] Tests: payable, refundable, 8% reduced, FA input, zero input
- [ ] `go vet ./... && go test -count=1 ./...` + commit

### A3: CIT engine rewrite
- [ ] Rate by company size: STANDARD 20 / SMALL 17 / MICRO 15 (rate table keys)
- [ ] Provisional vs final; 80% rule flag (CIT-10)
- [ ] Tests: standard/small/micro, loss (taxable=0), provisional<80% flag
- [ ] `go vet ./... && go test -count=1 ./...` + commit

### A4: PIT engine
- [ ] PITEmployeeInput struct (gross, dependants, residence, months)
- [ ] Progressive brackets 5-35% from rate table (PROGRESSIVE type)
- [ ] Deductions: personal 11M, dependant 4.4M, social 8%, health 1.5%, unemployment 1%
- [ ] CalculatePIT signature change + handler route body update + tests
- [ ] `go vet ./... && go test -count=1 ./...` + commit

### A5: Declaration engine
- [ ] taxService gains journalRepo + periodRepo deps (constructor + main.go + tests)
- [ ] GenerateDeclaration: pull posted journals for period → VAT/CIT engine → GTGT01/TNDN03 lines
- [ ] Cross-validation per TAX_RULES §2.1 ([30]=[23]-[16], XOR 31/32)
- [ ] Duplicate period+type rejection
- [ ] Tests: generation, validation fail, duplicate
- [ ] `go vet ./... && go test -count=1 ./...` + commit

### A6: Declaration→payment automation
- [ ] CreatePaymentFromDeclaration: amount from lines, due date, status PENDING
- [ ] Late flag: due date passed → calendar OVERDUE + alert
- [ ] Tests: payment created, late flag
- [ ] `go vet ./... && go test -count=1 ./...` + commit — Round A DONE

## Verification Gate
- [ ] Round A: `go vet ./... && go test -count=1 ./...` green
- [ ] Each slice: RED test → GREEN impl → commit
- [ ] No migrations added Round A

## ✅ Tax Round A: Calculation + Declaration Core (COMPLETE — c6c3a86)

### A1: Rate resolver (3e84268)
- [x] `resolveRate(taxType, applicableTo, onDate)` helper — rate table lookup, statutory fallback, suffix matching on RateCode
- [x] Tests: by-rate-code-suffix, active-window, empty-table-fallback

### A2: Rate-table VAT engine (a06e4db)
- [x] `CalculateVAT` reads explicit VAT accounts first (Cr 33311 output, Dr 1331/1332 input), rate × revenue fallback
- [x] OutputVAT/InputVAT/InputVATFA/payable/refundable from rate table
- [x] Tests per rule incl. zero-rate export, partial-credit purchase

### A3: Size-based CIT engine (9d85d2b)
- [x] `CalculateCIT` — revenue → size tier (MICRO 15% / SMALL 17% / STANDARD 20%) from rate table
- [x] Revenue - expenses + other income = taxable income
- [x] Handler test updated (100M → MICRO 15%)

### A4: Real PIT engine (c3398fd)
- [x] `CalculatePIT` progressive brackets (5-35%), resident/non-resident, dependants, insurance caps
- [x] **Breaking change approved**: `[]PITEmployeeInput` replaces `[]string employeeIDs` — interface + handler body
- [x] Tests: progressive bracket, non-resident flat 20%, insurance cap

### A5: Declaration engine (f57dd9d)
- [x] `GenerateDeclaration` — posted journals by period date range → CalculateVAT/CIT → form lines (FROM_LEDGER source) → cross-validation (TAX_RULES §2.1) → VALIDATED
- [x] GTGT01 lines 14/15/16/21/22/23/30/31/32; TNDN03 lines 04/06/12/13/14
- [x] Duplicate period+type → 409 ErrDuplicateDeclaration; KK_TNCN → 400
- [x] Zero declaration (no posted journals) supported — VAT-13
- [x] Company-scoped journal filter (cross-tenant fix)
- [x] `POST /api/v1/tax/declarations/generate`; NewTaxService + JournalRepository

### A6: Declaration→payment automation (c6c3a86)
- [x] AcknowledgeDeclaration auto-creates PENDING TaxPayment (payable lines [31]/[14])
- [x] Statutory due dates: VAT monthly 20th next month / quarterly 30th, CIT annual 31-Mar
- [x] Refundable/zero → no payment; idempotent per declaration
- [x] Validation fix: GTGT01 [30] algebraic sum may be negative
