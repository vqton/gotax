# FA Module — Workflows & Processes

**Version:** 1.0  
**Date:** 2026-07-30  
**Author:** BA Lead + Chief Accountant (20+ yrs each)  
**Regulations:** TT 45/2013/TT-BTC (TT 147/2024), VAS 03, VAS 04, TT 99/2025/TT-BTC, TT 200/2014/TT-BTC, Decree 123/2020/ND-CP

---

## Workflow Legend

Each workflow uses a structured table with the following columns:

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| Sequential step # | Person/role performing step | What they do | How system responds/automates | Business rules & validation checks | Documents created & GL entries posted |

---

## WF-01: Procure-to-Register (Mua sắm → Đăng ký TSCĐ)

End-to-end FA acquisition via purchase. Supports direct capitalization (asset ready on receipt) and CIP route (asset requires installation/assembly).

**Flow:**
```
Requisition → PO → Receive asset → Supplier invoice → FA registration → CIP or direct capitalization
                │
                ├── Direct:     Debit 211x, Credit 331; VAT Debit 1332
                └── CIP route:  Debit 2411, Credit 331 → (install) → Debit 211x, Credit 2411
```

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| 1.1 | Requester (user) | Create purchase requisition for FA (specify asset criteria, estimated cost, category, department) | PR form with FA flag; estimate ≥ 30M VND triggers FA classification | Original cost estimate ≥ 30M; useful life > 12 months per VAS 03 para 06; department not null | Draft PR in system |
| 1.2 | Approver (manager) | Approve PR with FA intent | Route to purchasing; mark PR as FA-procurement | PR total within budget; FA flag validated; category code exists in FA category tree | Approved PR → PO ready |
| 1.3 | Purchaser | Convert PR → PO; send to supplier | Generate PO number; populate FA fields (category hint, dep method default from category) | Supplier exists (tax code validated); terms match PR; PO amount ≠ 0 | PO issued; PO date recorded |
| 1.4 | Supplier | Deliver asset + issue e-invoice (XML per Decree 123 Art 12) | GDT webhook receives e-invoice XML; system parses and links to PO | Invoice matches PO (item, qty, price ± tolerance%); tax code valid; digital signature verified | E-invoice received; raw XML stored |
| 1.5 | Warehouse / Recipient | Receive asset physically; complete delivery receipt | Create receipt record: confirm condition, serial number, location; gateway for FA draft | QR/barcode scan matches PO item; serial number uniqueness check (per company) | Delivery receipt (Biên bản giao nhận); asset tag printed |
| 1.6 | FA Accountant | Create FA draft record from receipt + invoice | Pre-populate from PO/invoice: name, cost, supplier, serial; FA accountant enters: category, useful life, dep method, residual value, expense account | Cost = invoice ex-VAT + direct costs (transport, install, testing) per TT 45 Art 4; useful life within TT 45 Appendix 1 range; residual value ≥ 0 and ≤ 10% original cost (default); dep method per category default | FA draft (status = DRAFT); FA tracking card (Thẻ TSCĐ) |
| 1.7 | Chief Accountant | Approve FA registration | **Direct route:** post GL: Dr 211x, Cr 331 (or 111/112); Dr 1332, Cr 331. **CIP route:** Dr 2411, Cr 331; Dr 1332, Cr 331 | All VAS 03 recognition criteria satisfied (Art 2 TT 45); approval threshold if cost > board limit; CIP flag determines GL path | FA (status = ACTIVE → DEPRECIATING); GL: Dr 211x/2411, Cr 331 |
| 1.8 | FA Accountant | Print FA code label; affix to asset; file FA tracking card | Generate QR/barcode; update FA register (Sổ TSCĐ) | FA track card completeness: all mandatory fields filled | Physical asset tagged; FA register updated |

**GL Posting — Direct Purchase (FA ready for use immediately):**

| Account | Debit | Credit | Amount | Note |
|---------|-------|--------|--------|------|
| 211x | 500,000,000 | | Original cost | Asset account per category (2111-2118) |
| 1332 | 50,000,000 | | VAT input | VAT 10% on FA purchase |
| 331 | | 550,000,000 | Total | Supplier payable (or 111/112 if paid direct) |

**GL Posting — CIP Route (asset requires installation):**

| Step | Account | Debit | Credit | Amount |
|------|---------|-------|--------|--------|
| Reception | 2411 | 500,000,000 | | CIP accumulation |
| | 1332 | 50,000,000 | | VAT input |
| | 331 | | 550,000,000 | Supplier payable |
| Installation costs | 2411 | 20,000,000 | | Installation labor/materials |
| | 111/331 | | 20,000,000 | Cash/payable |
| Capitalization (step 1.7) | 211x | 520,000,000 | | Total CIP cost |
| | 2411 | | 520,000,000 | CIP cleared |

**Regulatory References:**
- VAS 03 para 06-07, 12-27 — Recognition, classification, original cost
- TT 45/2013/TT-BTC Art 2-5 — Recognition criteria, cost determination, purchase cost
- TT 99/2025/TT-BTC — TK 211, 241, 1332, 331
- Decree 123/2020/ND-CP Art 4, 9, 12 — E-invoice requirements
- TT 200/2014/TT-BTC — FA tracking card form (Mẫu số 02-TSCĐ)

---

## WF-02: Self-Construction to FA (Tự xây dựng → TSCĐ)

Capital works: construction of buildings, assembly of production lines, development of intangible assets (software per VAS 04 para 40).

**Flow:**
```
Construction planning → CIP accumulation (2412) → Completion certificate → Transfer to FA (211x) → Start depreciation
```

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| 2.1 | Project Manager | Submit construction plan with budget, timeline, completion criteria | Create CIP project record (TK 2412); link to company, department | Budget approved by board (if > threshold); timeline ≤ 36 months (default); completion criteria defined | CIP project (status = IN_PROGRESS); budget recorded |
| 2.2 | Accountant | Record construction costs as incurred (materials, labor, contractors) | Dr 2412, Cr 331/141/152/334 for each cost item | Cost documented (contract, delivery receipt, timesheet); VAT separable (1332 vs 1331); cost within budget line item | Cost entries on CIP: Dr 2412, Cr 331/141/334 |
| 2.3 | Project Manager | Submit completion report; request acceptance committee | System generates completion checklist; notifies acceptance committee members | All milestones met; total cost ≤ budget + 10% tolerance; quality inspection passed | Completion report; acceptance meeting scheduled |
| 2.4 | Acceptance Committee | Inspect asset; sign acceptance minutes (Biên bản nghiệm thu) | Record committee decision; document defects (if any) | Asset meets design specs; safety compliance verified; legal permits obtained (building permit, fire safety) | Acceptance minutes signed |
| 2.5 | Chief Accountant | Finalize CIP total cost; approve transfer to FA | Calculate final CIP cost = sum(2412) + borrowing costs (VAS 16) if capitalized; create FA record from CIP | Total cost confirms all direct costs included per VAS 03 para 14; borrowing cost capitalized per eligible period; CIP variance analyzed | FA record created from CIP |
| 2.6 | System | Transfer CIP to FA | Dr 211x (full CIP cost), Cr 2412; set dep start date = completion date; calculate depreciation from next month (or prorata from completion) | Carrying amount > 0; 2412 balance = 0 after transfer; FA status → ACTIVE → DEPRECIATING | GL: Dr 211x, Cr 2412; FA register updated |
| 2.7 | FA Accountant | Verify FA record; confirm depreciation start | Display FA tracking card for review; system schedules first depreciation | Useful life appropriate for asset type (building: 25-50 yr, machinery: 7-15 yr); dep method correct per category | FA operational; monthly dep schedule active |

**GL Posting — Construction Accumulation:**

| Account | Debit | Credit | Amount | Description |
|---------|-------|--------|--------|-------------|
| 2412 | 100,000,000 | | | Contractor invoice (phase 1) |
| 1331/1332 | 10,000,000 | | | VAT (if deductible) |
| 331 | | 110,000,000 | | Contractor payable |
| 2412 | 50,000,000 | | | Materials issued from warehouse |
| 152 | | 50,000,000 | | Raw materials consumed |
| 2412 | 30,000,000 | | | Direct labor |
| 334 | | 30,000,000 | | Employee salary payable |

**GL Posting — Completion Transfer:**

| Account | Debit | Credit | Amount |
|---------|-------|--------|--------|
| 2111 (Building) | 180,000,000 | | Total CIP cost |
| 2412 | | 180,000,000 | CIP cleared |
| 2111 | 20,000,000 | | Capitalized borrowing cost |
| 242 (if elected) | | | Or defer borrowing cost |

**Regulatory References:**
- VAS 03 para 14-15 — Self-constructed asset cost
- VAS 04 para 38-40 — Research vs development phase (intangible)
- TT 45/2013/TT-BTC Art 5 — Self-constructed FA original cost
- TT 99/2025/TT-BTC — TK 2412, 211x, 334, 152
- VAS 16 — Borrowing costs capitalization

---

## WF-03: Monthly Depreciation Run (Tính khấu hao hàng tháng)

Periodic batch process. FA Accountant triggers → engine calculates → Chief Accountant confirms → GL posted.

**Flow:**
```
Trigger → Select period → Engine runs → Preview results → Chief Accountant confirms → Post to GL → Update FA balances
```

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| 3.1 | FA Accountant | Open monthly depreciation screen; select period (year + month) | Display active FA count, prior month total depreciation, current month estimate | Period not closed; prior period depreciation posted (if any); current period not already calculated fully | Period selected; dashboard shows summary stats |
| 3.2 | FA Accountant | Click "Calculate Depreciation" | Engine iterates all FA with status = DEPRECIATING (or ACTIVE with dep start date ≤ period end) | Per asset: start date ≤ period end; end date ≥ period start (or null); not already calculated for this period | Depreciation entries (status = CALCULATED, gl_posted = false) |
| 3.3 | System | Depreciation engine runs per method | **Straight-line:** monthly = (cost - residual) / useful_life_months. **Declining balance:** (carrying × rate / 12). **Production:** (cost - residual) / total_units × actual_units. First period prorata = monthly × days_remaining / total_days | Method config per asset; prorata temporis for first period per TT 45 Art 10; last period adjusted to residual value; carrying amount never negative | Preview: list of FA with calculated amounts, total depreciation |
| 3.4 | FA Accountant | Review depreciation preview | Sortable grid: FA code, name, original cost, accumulated, calculated amount, expense account. Validate totals against prior month | Run variance check: month-over-month diff > 10% should flag; new FA acquisitions show first depreciation; disposed FA excluded | Preview confirmed or adjustments flagged |
| 3.5 | Chief Accountant | Confirm and approve depreciation run | System locks calculation; create GL journal entry lines | Total depreciation = sum of all calculated entries; no duplicate periods; all entries carry correct expense account per FA allocation | Depreciation entries (status = APPROVED) |
| 3.6 | System | Post depreciation to GL | Create compound journal entry: Dr 6274/6414/6424 (per FA expense account), Cr 2141/2142/2143 (per FA dep account). Update FA accumulated_depreciation, carrying_amount | GL posting rules: debit = credit; account exists and is active; period is open; no unposted entries remain | GL journal entry posted; FA balances updated; dep entries marked gl_posted = true |
| 3.7 | System | Update FA depreciation schedule | Extend depreciation schedule with new entry; check if FA now fully depreciated (carrying = residual or 0) | If accumulated + current ≥ original cost - residual: set status = FULLY_DEPRECIATED, end_date = period end | FA schedule updated; FULLY_DEPRECIATED flags set |
| 3.8 | FA Accountant | Verify GL posting; export depreciation schedule | View GL journal entry; run FA depreciation report; reconcile total dep vs GL account balance | Total Dr 2141/2142/2143 (period change) = total depreciation posted; matches sum of individual FA entries | Depreciation schedule exported (PDF/Excel); GL reconciled |

**Multi-Department Allocation:**

If FA has allocations across departments (fixed_asset_allocations), each allocation generates a separate GL line:

| FA | Dept | % | Expense Account | Dep Amount |
|----|------|---|-----------------|------------|
| Molding Machine 01 | Production A | 60% | 6274 | 6,000,000 |
| Molding Machine 01 | Production B | 40% | 6274 | 4,000,000 |

**GL Posting Example (single dept):**

| Account | Debit | Credit | Amount |
|---------|-------|--------|--------|
| 6274 (Production overhead — depreciation) | 10,000,000 | | |
| 6414 (Selling expense — depreciation) | 3,000,000 | | |
| 6424 (Admin expense — depreciation) | 2,000,000 | | |
| 2141 (Accum. dep — tangible FA) | | 15,000,000 | Total |

**Regulatory References:**
- TT 45/2013/TT-BTC Art 10-13 — Depreciation methods, calculation, first year/departure year
- VAS 03 para 29-36 — Depreciation: depreciable amount, method, useful life, component approach
- TT 99/2025/TT-BTC — TK 6274, 6414, 6424, 2141/2142/2143
- Circular 200/2014/TT-BTC — Depreciation schedule form

---

## WF-04: FA Adjustment/Revaluation (Điều chỉnh/Đánh giá lại TSCĐ)

Changes to FA value, useful life, depreciation method, or residual value. Requires Chief Accountant approval.

**Types:**
- **Value increase** (revaluation surplus, capitalized improvement) → Dr 211x, Cr 411/711
- **Value decrease** (impairment, partial disposal) → Dr 811, Cr 211x
- **Useful life change** → recalc depreciation prospectively
- **Method change** → recalc prospectively (VAS 03 para 34)
- **Residual value change** → recalc prospectively

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| 4.1 | Chief Accountant | Select FA; choose adjustment type | Display current FA details: cost, accumulated dep, carrying amount, useful life, method, residual value | FA status must be ACTIVE, DEPRECIATING, or FULLY_DEPRECIATED (not DISPOSED/SOLD) | Adjustment form initialized |
| 4.2 | Chief Accountant | Enter new values: new cost, new useful life, new method, new residual value, effective date, reason | System calculates differences: Δ cost, Δ useful life (months), Δ method, Δ residual | New useful life within TT 45 Appendix 1 range (if extended); method change valid (once per asset per TT 45); reason document referenced | Adjustment draft (status = PENDING) |
| 4.3 | System | Pre-calculate impact | Compute remaining depreciation schedule with new parameters; compare old remaining vs new remaining; show net P&L impact (if value change) | Remaining useful life > 0; carrying amount after adjustment ≥ new residual value; if useful life changed, document reason | Impact preview shown to Chief Accountant |
| 4.4 | Chief Accountant | Approve or reject adjustment | If approve: apply adjustment, post GL. If reject: discard | Approve requires digital signature (if company policy); rejection requires reason | Adjustment record created (status = APPROVED/REJECTED) |
| 4.5 | System | Post adjustment to GL | **Value increase:** Dr 211x (Δ cost), Cr 411 (revaluation surplus) or Cr 711 (if reversing prior impairment). **Value decrease:** Dr 811, Cr 211x. **No-value changes** (life/method/residual): no GL, update FA parameters only | Δ cost posted to correct GL account per adjustment type; 411 balance tracked per asset for future revaluation | GL posted; FA updated; depreciation recalculated from effective date |
| 4.6 | System | Update depreciation schedule | For value changes: recalculate remaining depreciation. For life/method/residual: recalc from effective date. For method change: remaining = carrying - residual, spread over remaining months | New schedule starts from effective date; no gap or overlap with old schedule | Depreciation schedule regenerated; FA transaction recorded |

**GL Posting — Value Increase (Revaluation Surplus):**

| Account | Debit | Credit | Amount | Note |
|---------|-------|--------|--------|------|
| 211x | 100,000,000 | | | Increase to original cost |
| 411 | | 100,000,000 | | Revaluation surplus (equity) |

**GL Posting — Value Decrease (Impairment):**

| Account | Debit | Credit | Amount |
|---------|-------|--------|--------|
| 811 | 50,000,000 | | Impairment loss (P&L) |
| 211x | | 50,000,000 | Direct cost reduction |
| 811 | 50,000,000 | | |
| 2294 | | 50,000,000 | Or via provision account (per TT 99) |

**Useful Life Change — Prospective Recalculation:**

```
Old: Original cost 500M, Residual 0, Life 120 months, Remaining 60 months
     Monthly dep = 4,166,667
New: Useful life extended to 180 months total (120 + 60 additional)
     Remaining useful life = 120 months
     Monthly dep (prospective) = Carrying amount / Remaining months
                                = (500M - 60×4.166M) / 120
                                = 250M / 120 = 2,083,333
```

**Regulatory References:**
- VAS 03 para 28 — Revaluation (only when mandated by law)
- VAS 03 para 33-34 — Useful life and method review at year-end
- VAS 04 para 65 — Intangible FA useful life and method review
- TT 45/2013/TT-BTC Art 9 — Change of useful life (max 1 change per asset)
- TT 99/2025/TT-BTC — TK 411, 711, 811, 2294, 211x
- IAS 16 para 31-42 — Revaluation model (IFRS entities)

---

## WF-05: FA Disposal (Thanh lý/Nhượng bán TSCĐ)

Remove FA from register through sale, liquidation, donation, or return to lessor. Calculate gain/loss.

**Types:**
- **Sale** (Nhượng bán) — Third-party buyer, e-invoice required, proceeds received
- **Liquidation** (Thanh lý) — Scrap/zero value, no proceeds or minimal
- **Donation** (Biếu tặng) — Charitable contribution, tax implications
- **Return to lessor** — Finance lease termination

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| 5.1 | Chief Accountant | Select FA; choose disposal type; enter disposal date, reason, counterparty | Display FA: original cost, accumulated dep, carrying amount, current location, department | FA status = DEPRECIATING or FULLY_DEPRECIATED; not already DISPOSED/SOLD; no pending inventory count | Disposal form (type, date, reason) |
| 5.2 | Chief Accountant | Enter disposal proceeds (net of selling costs) and payment terms | For sale: enter sale price, buyer info, payment due date. For liquidation: enter scrap value (0 if none). For donation: enter 0 | Proceeds ≥ 0; if sale price < fair value > 10%, require supporting reason; counterparty name required | Disposal draft (status = PENDING) |
| 5.3 | System | Calculate gain/loss | Gain/Loss = Proceeds - Carrying amount. Carrying amount = Original cost - Accumulated depreciation | Formula matches VAS 03 para 37; result is gain (positive) or loss (negative) | Disposal impact preview: gain/loss amount, proposed GL |
| 5.4 | Chief Accountant | Approve disposal | If approve: lock FA for disposal; trigger e-invoice creation (if sale). If reject: return to draft | Approval threshold check: if carrying amount > board limit, require higher approval | Disposal approved; e-invoice workflow triggered (for sale) |
| 5.5 | System | Post disposal GL | **Sale:** Dr 2141 (accumulated dep), Dr 111/112/131 (proceeds/AR), Dr 811 (if loss) or Cr 711 (if gain), Cr 211x (original cost). **Liquidation:** Dr 2141, Dr 811, Cr 211x | Debit = Credit; FA removed from FA register; accumulated dep cleared; GL posted to correct period | GL posted; FA status = SOLD or DISPOSED |
| 5.6 | System | Remove FA from depreciation schedule | Mark end_date = disposal date; no future depreciation calculated; FA removed from active depreciation run | Periods after disposal date: no entries in depreciation_entries | FA depreciable life ended |
| 5.7 | Tax Accountant | Handle VAT/CIT implications | For sale: e-invoice issued per Decree 123. For liquidation: internal record (no e-invoice). For donation: tax adjustment for CIT non-deductible portion | Sale price ≥ 0 → e-invoice required. Donation: 811 cost is tax non-deductible (CIT adjustment) | E-invoice sent (sale); CIT adjustment memo (donation) |
| 5.8 | FA Accountant | Archive FA tracking card; remove physical tag (if scrapped) | FA record kept in system (status = SOLD/DISPOSED) for audit trail; optionally print disposal report | FA register shows disposal; disposal record (Biên bản thanh lý) printed and signed | Disposal record archived; FA removed from active register |

**GL Posting — Sale (Gain):**

| Account | Debit | Credit | Amount | Description |
|---------|-------|--------|--------|-------------|
| 2141 | 400,000,000 | | | Remove accumulated dep |
| 111/112 | 200,000,000 | | | Proceeds received |
| 211x | | 500,000,000 | | Remove original cost |
| 711 | | 100,000,000 | | Gain (200M - (500M-400M)) |

**GL Posting — Sale (Loss):**

| Account | Debit | Credit | Amount |
|---------|-------|--------|--------|
| 2141 | 400,000,000 | | |
| 111/112 | 50,000,000 | | |
| 811 | 50,000,000 | | Loss |
| 211x | | 500,000,000 | |

**GL Posting — Liquidation (Zero Proceeds):**

| Account | Debit | Credit | Amount |
|---------|-------|--------|--------|
| 2141 | 400,000,000 | | |
| 811 | 100,000,000 | | Carrying amount = loss |
| 211x | | 500,000,000 | |

**Regulatory References:**
- VAS 03 para 37-38 — Derecognition, gain/loss calculation
- VAS 04 para 68-69 — Intangible FA derecognition
- TT 45/2013/TT-BTC Art 11 — Disposal/liquidation procedures
- TT 99/2025/TT-BTC — TK 211x, 2141, 111/112, 131, 711, 811
- Decree 123/2020/ND-CP Art 4, 9 — E-invoice for FA sale
- Circular 200/2014/TT-BTC — FA liquidation record form (Mẫu số 04-TSCĐ)

---

## WF-06: FA Transfer/Allocation (Điều chuyển/Phân bổ TSCĐ)

Transfer FA between departments or change allocation percentages for multi-department assets.

**Flow:**
```
Transfer request → Department head approve → Update department/location → Recalculate depreciation allocation
```

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| 6.1 | Department Head (source) | Request FA transfer: select FA, target department, target location, effective date, reason | Display FA current assignment; pre-fill from FA record | FA status = DEPRECIATING or FULLY_DEPRECIATED; target department exists and is active; effective date ≥ today (or first of next period) | Transfer request (status = PENDING) |
| 6.2 | Department Head (target) | Confirm receipt intent; review asset condition | Accept or reject transfer; if reject, provide reason | Target department head assigned in company org chart; asset condition notes recorded | Transfer confirmed or rejected |
| 6.3 | Chief Accountant | Approve transfer | If cross-company transfer (different legal entity): GL required. If same entity (same tax code): no GL, only department change | Same entity: Dr/Cr both 211x → no net GL impact. Different entity: Dr 211x (target), Cr 211x (source) with inter-company clearing (136/336) | Transfer approved |
| 6.4 | System | Update FA department, location, allocation | Change department_id, location fields. If allocation existed, update allocation records. If effective date mid-period, prorate depreciation | Mid-period transfer: source department dep for this period = full month × (days before transfer / total days); target department = rest | FA record updated; allocation recalculated |
| 6.5 | System | Adjust depreciation allocation prospectively | From effective date, future depreciation posts to new expense accounts (per new department's default dep account) | If multi-department allocation: split ratios per allocation table; sum = 100% | Depreciation schedule updated |
| 6.6 | FA Accountant | Verify transfer; update physical location tag | Print new location label; update FA tracking card with transfer history | Physical asset moved (if applicable); transfer record (Biên bản bàn giao) signed by both departments | Transfer record archived; FA register reflects new department |

**Multi-Department Allocation Model:**

```
Fixed Asset: "CNC Milling Machine" (Original cost 1.2B)

Default allocation (single dept):
  Production Dept:       100% → 6274

After transfer to shared allocation:
  Production Dept A:     60% → 6274
  Production Dept B:     30% → 6274
  Admin Dept:            10% → 6424

Depreciation posting each month:
  Dr 6274 (Dept A)       3,600,000  (60% × 6M)
  Dr 6274 (Dept B)       1,800,000  (30% × 6M)
  Dr 6424 (Admin)          600,000  (10% × 6M)
  Cr 2141                           6,000,000
```

**GL Posting — Cross-Company Transfer (same group, different legal entity):**

| Account | Debit | Credit | Amount |
|---------|-------|--------|--------|
| 136 (receivable from target entity) | 500,000,000 | | Transfer price (book value or agreed) |
| 2141 | 400,000,000 | | Accumulated dep transferred |
| 211x | | 500,000,000 | Remove from source |
| 711 | | 400,000,000 | Or gain if transfer price > book value |

**Regulatory References:**
- TT 99/2025/TT-BTC — TK 136, 336 (inter-company), 211x
- VAS 03 para 37 — Transfer is not derecognition (same entity)
- Decree 123/2020/ND-CP Art 9.2(g) — No e-invoice for internal transfers
- TT 200/2014/TT-BTC — FA transfer record form (Biên bản bàn giao)

---

## WF-07: FA Physical Inventory (Kiểm kê TSCĐ)

Periodic physical verification of FA existence, location, and condition. Mandatory at year-end; can be done quarterly or ad-hoc.

**Flow:**
```
Plan → Print tags/list → Physical count → Record results → Discrepancy handling → Adjustment
```

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| 7.1 | Chief Accountant | Create inventory plan: scope (all FA or by department/location), date, assign counters | Generate inventory plan record; list FA in scope with expected location and status | Plan date not in closed period; scope defined (department, category, or all); counters assigned (minimum 2 per location per internal control) | Inventory plan (status = DRAFT); FA checklist generated |
| 7.2 | FA Accountant | Print inventory tags/lists; prepare counting sheets | Export FA list by location: code, name, serial, expected location, barcode/QR | Each FA on list has unique barcode; tags printed in duplicate (one on asset, one as control) | Physical FA tags attached to assets; counting sheets ready |
| 7.3 | Inventory Team (2+ counters) | Physical count: scan/check each FA tag; record actual location and condition | Mobile app or paper form: scan barcode → confirm location → select condition (Good/Damaged/Not Found) | Every FA in scope accounted for (found or noted missing); unregistered assets discovered are tagged as "UNREGISTERED" | Count results recorded per FA |
| 7.4 | Inventory Team | Record discrepancies | For each discrepancy: enter actual vs expected. System auto-classifies discrepancy type: OK, MISSING, DAMAGED, UNREGISTERED | Missing: physical not found, tag present but asset absent. Damaged: physical found but not in usable condition. Unregistered: asset found but not in FA register | Inventory result records created; discrepancies identified |
| 7.5 | Chief Accountant | Review discrepancy report | Report shows: total inventoried, found, missing, damaged, unregistered. Variances: by department, by category | Materiality threshold: items > 10M VND individually reportable; total variance % vs FA register analyzed | Discrepancy report signed |
| 7.6 | Chief Accountant | Decide disposition for each discrepancy | **MISSING:** investigate; if confirmed lost → disposal (WF-05). **DAMAGED:** assess impairment (WF-04). **UNREGISTERED:** create FA record (WF-01 with source = INVENTORY). **OK:** no action | Missing: requires police report (if theft) or investigation memo. Damaged: impairment % assessed by technical team. Unregistered: cost determined by valuation or historical records | Disposition decisions recorded |
| 7.7 | FA Accountant | Execute dispositions | For each decision: initiate disposal, impairment, or registration workflow | Disposal: WF-05 steps. Impairment: WF-04 steps. Registration: WF-01 steps (source = inventory find) | FA register updated; discrepancies resolved |
| 7.8 | Chief Accountant | Sign inventory report; archive results | Inventory report generated: plan, results, discrepancies, dispositions, signatures | Report signed by Chief Accountant and Inventory Team Leader; filed for audit trail | Inventory completed; report archived (7 years per Law on Accounting) |

**Discrepancy Types:**

| Code | Name | Condition | Action |
|------|------|-----------|--------|
| OK | Khớp đúng | Asset found at expected location, condition as recorded | No action |
| MISSING | Thiếu TSCĐ | Asset not found at expected location | Investigate → disposal (missing) |
| DAMAGED | Hư hỏng | Asset found but damaged/cannot operate | Impairment assessment → repair or disposal |
| UNREGISTERED | Chưa đăng ký | Asset found but not in FA register | Create FA registration → retroactive depreciation (if applicable) |

**GL Posting — Inventory Adjustments:**

No direct GL from inventory results. GL impact occurs through resulting workflows:

| Discrepancy | Workflow | GL Entry |
|-------------|----------|----------|
| MISSING (confirmed lost) | Disposal (WF-05) | Dr 2141, Dr 811, Cr 211x |
| DAMAGED (impairment) | Adjustment (WF-04) | Dr 811, Cr 211x (or 2294) |
| UNREGISTERED (add to register) | Registration (WF-01) | Dr 211x, Cr 711 (or 421 if error correction) |

**Regulatory References:**
- TT 45/2013/TT-BTC Art 6 — FA management, physical inventory requirement
- Law on Accounting 2015 Art 36 — Physical inventory requirement
- VAS 03 para 39-40 — Disclosure: impairment, disposals
- Circular 200/2014/TT-BTC — FA inventory record form (Mẫu số 05-TSCĐ)
- TT 99/2025/TT-BTC — TK 2294, 811, 711

---

## WF-08: FA End-of-Year Closing (Cuối năm — TSCĐ)

Year-end procedures for FA module: reconcile, review, adjust, disclose.

**Flow:**
```
Reconcile FA ↔ GL → Verify depreciation → Impairment assessment → FA movement schedule → Notes to FS
```

| Step | Actor | Action | System | Validation | Output/GL |
|------|-------|--------|--------|------------|-----------|
| 8.1 | FA Accountant | Reconcile FA sub-ledger vs GL balances | Report: total original cost (FA register) vs GL 211x balance; total accumulated dep (FA register) vs GL 214x balance; total depreciation expense vs GL 6274+6414+6424 | FA register balance = GL account balance (difference = 0); period-end rates applied for FX (if foreign currency asset); FA count = number of active FA records | Reconciliation report |
| 8.2 | FA Accountant | Verify all depreciation posted for the year | System checks: for each month Jan-Dec, depreciation_entries.gl_posted = true for all DEPRECIATING FA; no gaps in schedule | Every month has depreciation run; all runs posted to GL; no FA missed (active but missing dep entry → investigate) | Depreciation verification report |
| 8.3 | FA Accountant | Review FA additions for the year | List all FA registered during year: source, original cost, capitalization date, useful life, dep method. Verify against purchase/construction records | Additions reconciled to actual purchases/CIP completions; capitalization dates within correct year; useful life and method within policy | FA additions schedule |
| 8.4 | FA Accountant | Review FA disposals for the year | List all FA disposed/sold during year: original cost, accumulated dep, proceeds, gain/loss. Verify disposal records signed | Disposals reconciled to sales/liquidation records; e-invoices issued for sales (if required); gain/loss correctly calculated | FA disposals schedule |
| 8.5 | Chief Accountant | Perform impairment assessment per VAS 03 / IAS 36 | For each FA with indicators (physical damage, obsolescence, market decline), compare carrying amount vs recoverable amount | Indicators per IAS 36 para 12: market value decline, physical damage, technological obsolescence, change in use, regulatory changes, net asset value > market cap | Impairment assessment memo; impairment entries (if needed) |
| 8.6 | Chief Accountant | Review useful life and residual value | For each FA category: compare actual useful life vs original estimate; check residual value still reasonable | If significant change: adjust prospectively per VAS 03 para 33-34; change documented with reason | Useful life review memo |
| 8.7 | System | Generate FA movement schedule (Bảng tăng giảm TSCĐ) | Per VAS 03 para 39-40 disclosure format: opening balance, additions, disposals, revaluation, impairment, depreciation, closing balance (separately for cost, accumulated dep, carrying amount) | By class of FA: 2111-2118, 213, 212. Three columns: original cost, accumulated dep, carrying amount. Movement column for each activity | FA movement schedule (TT 200 form B01-TSCĐ) |
| 8.8 | FA Accountant | Generate Notes to Financial Statements — FA section | Per VAS 03 para 40: measurement basis, depreciation method, useful lives/rates, gross carrying amount reconciliation, impairment details, commitments, pledged assets | Data matches FA movement schedule; all disclosures per VAS 03 checklist; comparative prior year figures included | Notes to FS — FA section complete |
| 8.9 | Chief Accountant | Sign off FA closing | Review all reports; sign movement schedule, impairment memo, reconciliation; authorize FA closing | All FA transactions recorded; all depreciation posted; no unreconciled differences; signed reports filed | FA year-end closing complete |

**FA Movement Schedule (Bảng tăng giảm TSCĐ — per VAS 03 para 39-40):**

| Item | 2111 Buildings | 2112 Machinery | 2113 Vehicles | 2114 Equipment | 213 Software | Total |
|------|----------------|----------------|---------------|----------------|--------------|-------|
| **Original cost** | | | | | | |
| Opening balance | 50,000M | 30,000M | 10,000M | 5,000M | 2,000M | 97,000M |
| Additions | 10,000M | 5,000M | 2,000M | 1,000M | 500M | 18,500M |
| Disposals | (2,000M) | (1,000M) | (500M) | (200M) | (100M) | (3,800M) |
| Revaluation | 3,000M | — | — | — | — | 3,000M |
| Closing balance | 61,000M | 34,000M | 11,500M | 5,800M | 2,400M | 114,700M |
| **Accumulated depreciation** | | | | | | |
| Opening balance | (15,000M) | (18,000M) | (6,000M) | (3,000M) | (800M) | (42,800M) |
| Depreciation for year | (2,000M) | (3,000M) | (1,500M) | (1,000M) | (400M) | (7,900M) |
| Disposals | 800M | 600M | 300M | 150M | 80M | 1,930M |
| Closing balance | (16,200M) | (20,400M) | (7,200M) | (3,850M) | (1,120M) | (48,770M) |
| **Carrying amount** | | | | | | |
| Opening | 35,000M | 12,000M | 4,000M | 2,000M | 1,200M | 54,200M |
| Closing | 44,800M | 13,600M | 4,300M | 1,950M | 1,280M | 65,930M |

**GL Reconciliation — Year-End Check:**

```
SELECT: SUM(original_cost) FROM fixed_assets WHERE status NOT IN ('DISPOSED','SOLD')
  SHOULD EQUAL: GL balance of accounts 211x + 212 + 213x (sum of all FA asset accounts)

SELECT: SUM(accumulated_depreciation) FROM fixed_assets
  SHOULD EQUAL: GL balance of accounts 2141 + 2142 + 2143

SELECT: SUM(depreciation_amount) FROM depreciation_entries WHERE period_year = 2026 AND gl_posted = true
  SHOULD EQUAL: GL period change in 214x (year-end balance - year-start balance) + disposals accumulated dep removed
```

**Regulatory References:**
- VAS 03 para 39-40 — Disclosure: movement schedule, depreciation method, useful lives, impairment
- VAS 04 para 70 — Intangible FA disclosure
- TT 99/2025/TT-BTC — TK 211x, 214x, 2294
- TT 200/2014/TT-BTC — Financial statement forms (B01-TSCĐ, B03-TSCĐ)
- Circular 228/2009/TT-BTC — Impairment guidance for listed companies
- IAS 36 — Impairment of assets (full guidance)
- IAS 16 para 73-79 — FA disclosure (IFRS)

---

## Regulatory Reference Summary

| Workflow | TT 45/2013 / TT 147/2024 | VAS 03 | VAS 04 | TT 99/2025 | TT 200/2014 | Decree 123/2020 |
|----------|--------------------------|--------|--------|------------|-------------|-----------------|
| WF-01 Procure-to-Register | Art 2-5, App 1 | Para 06-07, 12-27 | — | TK 211x, 2411, 1332, 331 | Mẫu 02-TSCĐ | Art 4, 9, 12 |
| WF-02 Self-Construction | Art 5 | Para 14-15 | Para 38-40 | TK 2412, 211x | — | — |
| WF-03 Monthly Depreciation | Art 10-13, App 2 | Para 29-36 | Para 60-65 | TK 6274, 6414, 6424, 214x | Mẫu bảng KH TSCĐ | — |
| WF-04 Adjustment/Revaluation | Art 9 | Para 28, 33-34 | Para 65 | TK 411, 711, 811, 2294 | — | — |
| WF-05 Disposal | Art 11 | Para 37-38 | Para 68-69 | TK 111/112, 131, 711, 811 | Mẫu 04-TSCĐ | Art 4, 9 |
| WF-06 Transfer/Allocation | — | Para 37 | — | TK 136, 336 | Mẫu bàn giao | Art 9.2(g) |
| WF-07 Physical Inventory | Art 6 | — | — | TK 2294, 811 | Mẫu 05-TSCĐ | — |
| WF-08 Year-End Closing | — | Para 39-40 | Para 70 | TK 211x-214x | B01-TSCĐ, B03-TSCĐ | — |

---

## FA Document Templates (per Vietnamese Accounting Regime)

| Template | Form Code | Workflow | Description |
|----------|-----------|----------|-------------|
| FA Delivery Receipt | Biên bản giao nhận TSCĐ | WF-01 | Signed when asset received from supplier |
| FA Tracking Card | Thẻ TSCĐ (Mẫu 02-TSCĐ) | WF-01 | Perpetual record per individual asset |
| CIP Completion Certificate | Biên bản nghiệm thu | WF-02 | Signed by acceptance committee |
| Depreciation Schedule | Bảng tính khấu hao TSCĐ | WF-03 | Monthly calculation output |
| FA Adjustment Record | Biên bản điều chỉnh TSCĐ | WF-04 | Documented changes to FA parameters |
| FA Liquidation Record | Biên bản thanh lý TSCĐ (Mẫu 04-TSCĐ) | WF-05 | Disposal documentation |
| FA Transfer Record | Biên bản bàn giao TSCĐ | WF-06 | Department-to-department transfer |
| FA Inventory Record | Biên bản kiểm kê TSCĐ (Mẫu 05-TSCĐ) | WF-07 | Physical count results |
| FA Movement Schedule | Bảng tăng giảm TSCĐ (B01-TSCĐ) | WF-08 | Year-end movement report |

---

## State Transitions Per Workflow

```
WF-01: DRAFT → ACTIVE → DEPRECIATING
WF-02: DRAFT → ACTIVE → DEPRECIATING
WF-03: DEPRECIATING → DEPRECIATING (no state change, or → FULLY_DEPRECIATED)
WF-04: DEPRECIATING / FULLY_DEPRECIATED → DEPRECIATING (recalc)
WF-05: DEPRECIATING / FULLY_DEPRECIATED → DISPOSED / SOLD
WF-06: DEPRECIATING / FULLY_DEPRECIATED → DEPRECIATING / FULLY_DEPRECIATED (no state change)
WF-07: No state change (information gathering only; resulting actions follow WF-01/04/05)
WF-08: No state change (reporting only)
```
