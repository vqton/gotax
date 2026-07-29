# AR Module — Task Tracker

## Phase 0: Critical Fix

### Task 1: Fix AR Aging Report — Bucket by DueDate

**Desc:** `GetARAgingReport` puts all amount in Bucket0. Fix: query invoices with DueDate, calc days overdue from report date, group into buckets (current, 1-30, 31-60, 61-90, 91-120, 120+).

**AC:**
- [ ] Aging report shows correct bucket per invoice BalanceDue
- [ ] Current (<0 days overdue) in Bucket0
- [ ] 1-30 days overdue in Bucket30
- [ ] 31-60 days overdue in Bucket60
- [ ] 61-90 days overdue in Bucket90
- [ ] 120+ days overdue in Bucket120
- [ ] Existing tests updated

**Verification:** `./go test -v -run TestARAging ./internal/service/` (or handler equivalent)

**Files:** `internal/service/sale_service.go` (GetARAgingReport), `internal/domain/models_sale.go` (ARAgingReport struct?), possibly `internal/repository/pg_sale.go` + `memory_sale.go`

**Scope:** S (2-3 files)

---

## Phase 1: Core AR (P0)

### Task 2: GL Auto-Posting on Invoice Post

**Desc:** When invoice transitions ISSUED→POSTED, auto-create journal entry: Dr 131 (total_amount), Cr revenue_account (subtotal), Cr 3331 (tax_amount). Use existing `service.Service.CreateEntry()`.

**AC:**
- [ ] PostInvoice creates JournalEntry via GL service
- [ ] GL entry matches: Dr AR account (131), Cr revenue account (511X), Cr VAT account (3331)
- [ ] `gl_posted=true`, `gl_posted_at` set
- [ ] Error in GL post → rollback invoice status
- [ ] Test: posted invoice has matching GL entry

**Files:** `internal/service/sale_service.go`, `internal/service/sale_service_test.go`, DI wiring in `main.go`

**Scope:** M (3-4 files)

---

### Task 3: GL Auto-Posting on Receipt Post

**Desc:** When receipt transitions DRAFT→POSTED, auto-create journal entry: Dr 111/112 (amount), Cr 131 (amount). Handle allocation-level discount (Dr 5213).

**AC:**
- [ ] PostReceipt creates JournalEntry: Dr cash/bank, Cr AR
- [ ] Early payment discount → Dr 5213
- [ ] Error in GL post → rollback receipt status
- [ ] Test: posted receipt has matching GL entry

**Files:** `internal/service/sale_service.go`, `internal/service/sale_service_test.go`, `main.go`

**Scope:** M (3-4 files)

---

### Task 4: GL Auto-Posting on Credit Note Post

**Desc:** When credit note transitions DRAFT→POSTED, auto-create reversing journal entry: Dr revenue_account (subtotal), Dr 3331 (VAT), Cr 131 (total_amount). If goods returned → also reverse COGS (Dr inventory Cr 632).

**AC:**
- [ ] PostCN creates JournalEntry: Dr 511X/3331, Cr 131
- [ ] If return with goods → Dr inventory account Cr 632
- [ ] `gl_posted=true`, `gl_posted_at` set
- [ ] Test: posted CN reduces AR balance

**Files:** `internal/service/sale_service.go`, `internal/service/sale_service_test.go`, `main.go`

**Scope:** M (3-4 files)

---

### Task 5: Customer Statement Endpoint

**Desc:** Generate customer statement of account: opening balance, transactions (invoices, receipts, credit notes) in period, closing balance. PDF/JSON output.

**AC:**
- [ ] GET /api/v1/sales/customer-statement/:id returns statement
- [ ] Opening = closing of prior period
- [ ] All period transactions listed
- [ ] Closing = opening + invoiced - received - credit notes
- [ ] Multi-currency support
- [ ] Test: statement matches invoice+receipt+CN totals

**Files:** `internal/handler/sale_handler.go`, `internal/service/sale_service.go`, `internal/domain/models_sale.go` (statement model), route registration

**Scope:** M (4-5 files)

**Checkpoint: Phase 1**
- [ ] Invoice posted → GL entry exists → statement shows it
- [ ] Receipt posted → GL entry exists → statement shows it
- [ ] Credit note posted → GL entry exists → statement shows it
- [ ] All Phase 1 tests pass

---

## Phase 2: E-Invoice Pipeline (P0)

### Task 6: E-Invoice Pipeline Interface + TXML Generator

**Desc:** Define e-invoice pipeline interfaces (TXMLGenerator, Signer, GDTSubmitter). Implement TXML generator per Decree 254 format — seller, buyer, line items, totals, VAT breakdown.

**AC:**
- [ ] TXMLGenerator interface defined
- [ ] Default implementation generates valid TXML XML
- [ ] XML includes: seller info, buyer info, line items, totals, VAT by rate
- [ ] Amount in words (Vietnamese)
- [ ] QR code URL format correct
- [ ] Test: generated XML validates against Decree 254 schema

**Files:** `internal/einvoice/` (new package), `internal/domain/interfaces.go` (if interfaces go there)

**Scope:** M (3-5 files in new package)

---

### Task 7: E-Invoice Digital Signer

**Desc:** Implement XMLDSig signing using company's DigitalSignature (cert + private key). Wraps TXML with digital signature per GDT requirements.

**AC:**
- [ ] Signer signs TXML with RSA-SHA256
- [ ] Signature embeds certificate info
- [ ] Signed XML valid against GDT schema
- [ ] Error if cert expired or key invalid
- [ ] Test: signed XML passes verification

**Files:** `internal/einvoice/signer.go`, `internal/domain/models_company.go` (DigitalSignature usage)

**Scope:** M (2-3 files)

---

### Task 8: GDT API Client

**Desc:** GDT REST API client for e-invoice operations: submit invoice, cancel invoice, replace invoice, query status. Uses existing IntegrationProfile for credentials.

**AC:**
- [ ] Submit signed TXML → receive invoice code
- [ ] Cancel invoice by code
- [ ] Replace invoice
- [ ] Query invoice status
- [ ] Rate limit + retry handling
- [ ] Mock GDT for tests
- [ ] Test: integration test with mock GDT

**Files:** `internal/einvoice/gdt_client.go`, `internal/domain/models_company.go` (IntegrationProfile)

**Scope:** L (4-6 files including mock + test)

---

### Task 9: E-Invoice Auto-Pipeline Hook

**Desc:** Hook e-invoice pipeline into PostInvoice: on ISSUED→POSTED, auto-generate TXML → sign → submit to GDT → store code. Configurable auto vs manual.

**AC:**
- [ ] PostInvoice triggers pipeline if auto enabled
- [ ] Pipeline: TXML generate → sign → GDT submit → store code
- [ ] Status transitions: SIGNED → SUBMITTED → CODED
- [ ] Failure → manual retry endpoint
- [ ] Config: auto-pipeline on/off per company
- [ ] Test: mock pipeline verifies state transitions

**Files:** `internal/service/sale_service.go`, `internal/einvoice/pipeline.go`, `main.go`

**Scope:** M (3-5 files)

**Checkpoint: Phase 2**
- [ ] TXML generates valid XML
- [ ] Signing produces verifiable signature
- [ ] GDT mock returns invoice code
- [ ] Full pipeline end-to-end with mock GDT
- [ ] All Phase 1 + 2 tests pass

---

## Phase 3: Collection Management (P1)

### Task 10: Dunning Engine

**Desc:** Dunning levels (L1-L4) with configurable triggers, auto-email/SMS, collection notes, promise-to-pay tracking.

**AC:**
- [ ] Dunning level config (days overdue, amount threshold)
- [ ] Auto-queue invoices for dunning daily batch
- [ ] L1-L4 escalation logic
- [ ] Promise-to-pay record (suspends dunning)
- [ ] Dunning letter template + send
- [ ] Test: dunning batch processes correct invoices

**Files:** `internal/service/collection_service.go` (new), `internal/handler/collection_handler.go` (new), `internal/domain/models_collection.go` (new), routes

**Scope:** L (6-8 files)

---

### Task 11: Bad Debt Provision + Write-Off

**Desc:** Provision calc per VAS 17 (R8): by aging bucket %. Write-off with director approval workflow. Off-balance-sheet tracking.

**AC:**
- [ ] Bad debt provision calc (month-end batch)
- [ ] Provision % configurable per bucket
- [ ] Journal entry: Dr 642 Cr 229 (provision)
- [ ] Write-off with approval (segregation of duty)
- [ ] Write-off: Dr 229 Cr 131 (if provisioned) or Dr 642 Cr 131
- [ ] Off-balance-sheet tracking account
- [ ] Recovery: Dr 131 Cr 711 + Dr bank Cr 131

**Files:** `internal/service/collection_service.go`, `internal/domain/models_collection.go`, `internal/handler/sale_handler.go`

**Scope:** L (5-7 files)

---

### Task 12: Prepayment/Deposit Workflow

**Desc:** Full prepayment lifecycle: receive prepayment → optional deposit e-invoice → offset against final invoice → remaining balance invoice.

**AC:**
- [ ] Receipt with type=prepayment creates credit AR balance
- [ ] Optional deposit e-invoice issuance
- [ ] Offset prepayment against invoice (clearing entry)
- [ ] Final invoice = total - deposit
- [ ] Test: prepayment offsets correctly

**Files:** `internal/service/sale_service.go`, `internal/handler/sale_handler.go`, `internal/domain/models_sale.go`

**Scope:** M (4-5 files)

---

### Task 13: FX Revaluation

**Desc:** At receipt time: calc realized FX gain/loss between invoice rate and receipt rate. At month-end: revalue outstanding foreign currency AR at period-end rate.

**AC:**
- [ ] Realized FX at receipt allocation (Dr 635/Cr 515)
- [ ] Month-end unrealized FX revaluation batch
- [ ] Reval journal: Dr/Cr 131, Cr/Dr 515/635
- [ ] Test: FX gain/loss posted correctly

**Files:** `internal/service/sale_service.go`, `internal/service/fx_service.go` (new), `internal/domain/models.go` (ExchangeRate usage)

**Scope:** M (4-5 files)

**Checkpoint: Phase 3**
- [ ] Dunning levels escalate correctly
- [ ] Bad debt write-off flow end-to-end
- [ ] Prepayment offsets against invoice
- [ ] FX gain/loss posts correctly
- [ ] All tests pass

---

## Phase 4: Controls & Monitoring (P1/P2)

### Task 14: Credit Limit Enforcement

**Desc:** Check at SO confirm + invoice create: outstanding AR + new amount ≤ credit limit. Configurable action (warn/block).

**AC:**
- [ ] Credit limit check at SO confirm
- [ ] Credit limit check at invoice create
- [ ] Config: warn vs block per company
- [ ] Manager override with reason
- [ ] Test: exceeding limit blocked

**Files:** `internal/service/sale_service.go`

**Scope:** S (2-3 files)

---

### Task 15: AR-GL Month-End Reconciliation

**Desc:** Report comparing AR sub-ledger (Σ customer invoice balances) vs GL 131 balance. List variances by invoice.

**AC:**
- [ ] Report: sub-ledger total → GL 131 balance → variance
- [ ] Variance drill-down by invoice
- [ ] Test: reconciled when GL matches sub-ledger

**Files:** `internal/handler/sale_handler.go`, `internal/service/report_service.go` (or sale_service.go)

**Scope:** M (3-4 files)

---

### Task 16: DSO + AR KPI Dashboard

**Desc:** Days Sales Outstanding calculation + AR KPI endpoints: total AR, overdue %, collection rate, avg days to pay.

**AC:**
- [ ] DSO = (Avg AR / Total Revenue) × days in period
- [ ] Total AR outstanding
- [ ] Overdue % of total AR
- [ ] Collection rate (received / invoiced) per period
- [ ] Test: DSO matches manual calc

**Files:** `internal/handler/sale_handler.go`, `internal/service/sale_service.go`

**Scope:** M (3-4 files)

---

### Task 17: Off-Balance-Sheet Tracking

**Desc:** Written-off AR tracked off-balance-sheet (monitoring account) for 10 years per Accounting Law. Separate table + reconciliation.

**AC:**
- [ ] Write-off moves AR to off-balance-sheet table
- [ ] Reconciliation: written-off vs off-balance-sheet total
- [ ] Recovery: remove from off-balance-sheet, record cash
- [ ] Test: write-off tracked correctly

**Files:** `internal/domain/models_collection.go`, `internal/repository/pg_sale.go`, `internal/repository/memory_sale.go`, migration

**Scope:** M (4-5 files)

**Checkpoint: Phase 4**
- [ ] Credit limit blocks over-limit orders
- [ ] AR-GL reconciliation report zero variance after GL posting
- [ ] DSO calc reasonable
- [ ] Written-off AR tracked off-balance-sheet
- [ ] All 17 tasks done, all tests pass

---

## Progress

- [ ] Phase 0: Fix AR Aging (P0)
- [ ] Phase 1: Core AR (P0)
- [ ] Phase 2: E-Invoice Pipeline (P0)
- [ ] Phase 3: Collection Management (P1)
- [ ] Phase 4: Controls & Monitoring (P1/P2)
