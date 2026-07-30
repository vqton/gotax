# Fixed Asset Module — Use Cases

**Version:** 1.0
**Date:** 2026-07-30
**Author:** BA Lead + Chief Accountant (20+ yrs each)

**Regulatory references:**
- Circular 45/2013/TT-BTC (TT 45) — FA management, depreciation, useful life table
- Circular 147/2024/TT-BTC (TT 147, supersedes TT 45) — updated FA regime
- VAS 03 — Tangible Fixed Assets
- VAS 04 — Intangible Fixed Assets
- Circular 200/2014/TT-BTC (TT 200) — Enterprise accounting regime
- Circular 133/2016/TT-BTC (TT 133) — SME accounting regime
- Circular 99/2025/TT-BTC (TT 99) — Chart of Accounts
- Decree 123/2020/ND-CP — E-invoices
- IAS 16 — Property, Plant & Equipment
- IAS 36 — Impairment of Assets

---

## UC-01: Register New Fixed Asset (Acquisition)

| Element | Detail |
|---------|--------|
| **ID** | UC-01 |
| **Title** | Register New Fixed Asset (Acquisition) |
| **Actor** | FA Accountant (Kế toán TSCĐ) |
| **Description** | Register a new fixed asset from any acquisition source. Asset enters the FA register, is valued at cost, and depreciation begins upon activation. Original cost must be ≥ 30M VND with useful life > 1 year per VAS 03 / TT 45. |
| **Sources** | Purchase (mua), Self-construction (tự xây dựng), Finance lease (thuê tài chính), Donation (biếu tặng), Capital contribution (góp vốn), Exchange (trao đổi) |
| **Preconditions** | FA categories configured. Company selected. User has FA create permission. Supplier/customer master data exists (if purchase source). |
| **Postconditions** | FA created with DRAFT or ACTIVE status. GL posted for acquisition. FA register updated. Depreciation start date set if activated. |

### Happy Path
1. FA Accountant navigates to Fixed Assets → Register New FA
2. System displays acquisition source types: Purchase, Self-construction, Finance lease, Donation, Capital contribution, Exchange
3. User selects source type (e.g., Purchase)
4. System loads category tree (Group → Type → Category, 3 levels)
5. User selects category; system auto-fills defaults from category:
   - Asset account (211x/212/213x)
   - Depreciation account (2141/2142/2143)
   - Expense account (6274/6414/6424)
   - Default useful life (months)
   - Default depreciation method
6. User fills FA details:
   - Code (auto-generated or manual: TS-YYYYMM-XXXX)
   - Name, serial number, manufacturer, model, year, country of origin
   - Acquisition date, original cost (≥ 30M VND), residual value
   - Useful life (months), depreciation method, depreciation start date
   - Department, location, user responsible
   - Supplier, contract number, e-invoice reference (linked from purchase module)
7. User fills allocation: department(s) and percentage (sum = 100%), expense account per department
8. User attaches supporting documents: purchase invoice, contract, delivery receipt
9. System validates:
   - Original cost ≥ 30M VND (per TT 45 §4 / TT 147 §5)
   - Useful life > 12 months
   - Category mandatory, department mandatory
   - E-invoice linked for VAT deduction (Decree 123)
   - Allocation percentages sum to 100%
10. User saves as DRAFT (no GL posting yet)
11. User reviews and clicks Submit for Approval
12. Approval workflow triggers (if configured: chief accountant approval)
13. Chief Accountant approves
14. System activates FA (status → ACTIVE) and sets depreciation start date
15. System posts acquisition GL entry:
    - Debit 211x (original cost)
    - Debit 1332 (VAT input, if purchase with VAT)
    - Credit 331/111/112 (total payment)
16. System records transaction in `fixed_asset_transactions` (type: ACQUISITION)
17. System returns FA detail with generated ID and posted GL journal ID

### Alternative Paths
- **A1: Acquisition sources** — Different source = different default GL posting:
  - Self-construction: Debit 211x, Credit 2412 (CIP transferred on completion)
  - Finance lease: Debit 212, Credit 341 (lease liability); separate depreciation per VAS 06 / IFRS 16
  - Donation: Debit 211x, Credit 711 (donation income, no cash outflow)
  - Capital contribution: Debit 211x, Credit 411 (equity contribution)
  - Exchange: Debit 211x (fair value), Debit/credit 811/711 for gain/loss on exchange
- **A2: Auto-create from purchase invoice** — Purchase module sends supplier invoice with FA-type items → system prompts: "Create FA from this invoice?" → pre-fills cost, supplier, invoice ref → user completes remaining fields
- **A3: Direct activation (skip draft)** — User with chief accountant role can activate directly, skipping DRAFT and approval steps
- **A4: CIP route (see UC-09)** — Costs accumulated in 2411 first, then transferred to 211 on completion certificate
- **A5: Batch registration** — User uploads spreadsheet with multiple FA → system validates each row → shows error summary → creates all valid records

### Exception Paths
- **E1: Cost < 30M VND** — System rejects: "Asset does not meet FA recognition threshold. Consider recording as expense (tool/LVP)." Per TT 45 §4: assets under 30M must be expensed immediately.
- **E2: Missing e-invoice for VAT deduction** — System warns: "No valid e-invoice linked. VAT input (1332) cannot be claimed per Decree 123." Allows save with warning, but VAT claimed separately.
- **E3: Duplicate FA code** — System blocks: "FA code TS-2026-0001 already exists in this company."
- **E4: Category not configured** — System rejects: "FA category required. Configure categories first."
- **E5: Allocation percentages ≠ 100%** — System rejects: "Department allocation must sum to 100%. Current total: X%."
- **E6: Useful life beyond category max** — Warning: "Useful life exceeds category default by X%. Approver must confirm."

### Business Rules
- BR-FA-01: FA recognition: original cost ≥ 30M VND, useful life > 12 months (TT 45 §4)
- BR-FA-02: FA code unique per company
- BR-FA-03: Acquisition VAT on FA (1332) requires valid e-invoice (Decree 123 §4)
- BR-FA-04: Donated FA valued at fair market price with appraisal report
- BR-FA-05: Capital-contributed FA valued at contributing party's agreed price (must be independent)
- BR-FA-06: Exchange FA valued at fair value of asset received or given up, whichever is more reliable (VAS 03 §26)

---

## UC-02: Run Monthly Depreciation

| Element | Detail |
|---------|--------|
| **ID** | UC-02 |
| **Title** | Run Monthly Depreciation |
| **Actor** | Chief Accountant / FA Accountant (Kế toán trưởng / Kế toán TSCĐ) |
| **Description** | Calculate and post monthly depreciation for all active depreciating assets. Supports straight-line, declining balance, and production-based methods per Circular 45/2013. Generates GL entries debiting expense accounts and crediting accumulated depreciation. |
| **Preconditions** | Depreciation period (month/year) is open in GL. At least one FA in DEPRECIATING status. No prior depreciation run for this period (or prior run unposted). |
| **Postconditions** | Depreciation entries calculated for all eligible assets. GL posted. Accumulated depreciation and carrying amounts updated. |

### Happy Path
1. Chief Accountant navigates to Fixed Assets → Depreciation → Run Depreciation
2. System displays company, period selection (month, year)
3. User selects period (e.g., July 2026)
4. System validates period is open (not closed in GL)
5. System checks previous depreciation for this period:
   - If already calculated but unposted: shows preview, ask "Post now?"
   - If already posted: blocks (see E1)
6. User clicks "Calculate"
7. System loads all FA with status = DEPRECIATING and depreciation_start_date ≤ period end
8. For each asset, system calculates monthly depreciation:
   - **Straight-line:** `(OriginalCost - ResidualValue) / UsefulLifeMonths`
   - **Declining balance:** `CarryingAmount × ((1 / UsefulLifeYears) × Factor) / 12`
   - **Production-based:** `(OriginalCost - ResidualValue) / TotalEstimatedProduction × ActualMonthlyProduction`
9. First-period prorata applied: `MonthlyAmount × (DaysRemainingInMonth / TotalDaysInMonth)`
   - Depreciation start date is mid-month → prorata from start date to month end
10. Last-period adjustment: final entry = remaining carrying amount − residual value
11. System generates depreciation preview table:
    - Asset code, name, method, original cost, carrying amount before, monthly amount, accumulated after, carrying amount after
    - Grouped by expense account (6274, 6414, 6424)
    - Total depreciation per expense account
12. User reviews preview
13. User clicks "Post to GL"
14. System validates: sufficient balance in expense accounts (warning only)
15. System creates GL journal entry:
    - Debit 6274 / 6414 / 6424 (by department allocation)
    - Credit 2141 / 2142 / 2143
16. For multi-department assets, system splits depreciation per allocation percentage:
    - Department A (60%): Debit 6274, Credit 2141 = 60% of monthly amount
    - Department B (40%): Debit 6414, Credit 2141 = 40% of monthly amount
17. System updates each asset:
    - `accumulated_depreciation` += monthly amount
    - `carrying_amount` -= monthly amount
18. System records `depreciation_entries` per asset per period (unique constraint)
19. System records `fixed_asset_transactions` (type: DEPRECIATION)
20. System marks period depreciation as POSTED
21. System returns journal entry ID and summary: "Depreciation for Jul 2026 posted. Assets: 245, Total: 1,245,678,900 VND"

### Alternative Paths
- **A1: Preview only (no post)** — User clicks "Calculate" then "Save as Draft" (entries saved as unposted). Can resume later with "Post pending entries."
- **A2: Suspended assets skip** — Assets in SUSPENDED status are excluded from calculation. System notifies: "5 assets suspended, skipped."
- **A3: Fully depreciated assets** — Assets where carrying amount = residual value are skipped. System notifies: "12 assets fully depreciated, skipped."
- **A4: First-period prorata** — New asset activated mid-month: depreciation amount prorated from activation date to month end. Example: activated 15-Jul, monthly dep = 10M, prorata = 10M × 16/31 = 5,161,290.
- **A5: Last-period adjustment** — Final month: depreciation = remaining carrying amount − residual value (not standard monthly amount). Ensures carrying amount = residual value at end of useful life.
- **A6: Production-based — actual quantity entry** — User prompted to enter actual production quantity for each production-based asset before calculation.
- **A7: Unpost depreciation** — Admin may unpost a period's depreciation (reversal GL entry created). Requires reason. Cannot unpost if GL period is closed.

### Exception Paths
- **E1: Period already posted** — Error: "Depreciation for Jul 2026 already posted on 2026-07-31 by user X. Use 'Unpost' to recalculate."
- **E2: Asset with zero carrying amount** — Error: "Asset TS-2026-0100 (Machine A) has carrying amount = 0. Verify data or mark as fully depreciated."
- **E3: Conflicting allocation percentages** — Error: "Asset TS-2026-0100 allocation sum is 95% (not 100%). Update allocation before running depreciation."
- **E4: No depreciating assets for period** — Warning: "No assets in DEPRECIATING status for Jul 2026." Proceed without generating entries.
- **E5: GL period closed** — Block: "Period Jul 2026 is CLOSED in GL. Open period or select a different month."
- **E6: Depreciation start date after period end** — Asset ignored with notification: "Asset TS-2026-0150 starts depreciation Aug 2026, not included."

### Business Rules
- BR-DEP-01: Depreciation runs per asset per period (unique fixed_asset_id + period_id)
- BR-DEP-02: Straight-line: equal monthly amounts except first/last period (VAS 03 §12)
- BR-DEP-03: Declining balance: factor = 2 (double-declining) or configurable, switch to SL optional
- BR-DEP-04: Production-based: requires actual monthly production input
- BR-DEP-05: Depreciation starts month after activation (or prorata from activation date, configurable)
- BR-DEP-06: FA under construction (CIP) does not depreciate until transferred to 211 (VAS 03 §9)
- BR-DEP-07: Depreciation method change = accounting estimate change (prospective, VAS 03 §31)
- BR-DEP-08: Suspended assets: no depreciation during suspension period (TT 45 §10)

---

## UC-03: Transfer Fixed Asset Between Departments

| Element | Detail |
|---------|--------|
| **ID** | UC-03 |
| **Title** | Transfer Fixed Asset Between Departments |
| **Actor** | FA Accountant (Kế toán TSCĐ) |
| **Description** | Move an FA from one department to another. Asset account (211x) unchanged — only depreciation allocation changes. Mid-period transfers apply prorata allocation for source and target departments. |
| **Preconditions** | Asset exists, status = ACTIVE or DEPRECIATING. Target department exists. User has transfer permission. |
| **Postconditions** | Asset department/location updated. Allocation records recalculated. Depreciation entries for current period updated if already calculated. Transaction recorded. |

### Happy Path
1. FA Accountant navigates to Fixed Assets → Transfer
2. User selects asset (search by code, name, or scan barcode)
3. System displays current FA detail:
   - Current department, location, user responsible
   - Current allocation percentages (if multi-department)
   - Depreciation method, monthly amount, carrying amount
4. User selects transfer type: Full transfer (whole asset) or Partial transfer (allocation change)
5. If full transfer: user selects new department, new location (optional), new responsible user (optional)
   - System auto-sets allocation to 100% for new department
6. If partial transfer: user adjusts allocation percentages across departments
   - Must sum to 100%
7. User sets effective date (default = today)
8. System validates:
   - Asset status is ACTIVE or DEPRECIATING (not DISPOSED, SOLD, SUSPENDED)
   - Effective date not in closed GL period
   - Allocation percentages sum to 100%
9. System calculates prorata if mid-period transfer (A1)
10. User reviews transfer summary
11. User confirms
12. System updates asset record: department_id, location, updated_by, updated_at
13. System updates or creates allocation records
14. System records `fixed_asset_transactions` (type: TRANSFER):
    - old_value = previous department ID, new_value = new department ID
    - Old allocation percentages, new allocation percentages
15. If depreciation already calculated for current period:
    - System recalculates entries with new allocation (if unposted)
    - If posted: blocks or forces unpost first
16. System returns transfer confirmation with updated asset detail

### Alternative Paths
- **A1: Mid-period transfer (prorata)** — Effective date is mid-month:
  - Source department depreciated for days in possession: `MonthlyAmount × (1..effective_date-1) / TotalDays`
  - Target department depreciated for remaining days: `MonthlyAmount × (effective_date..month_end) / TotalDays`
  - System splits the single monthly entry into two allocation entries per department
  - Example: Transfer 15-Jul, asset monthly dep = 10M, full allocation to Dept A before, Dept B after:
    - Dept A depreciation (1 Jul - 14 Jul): 10M × 14/31 = 4,516,129
    - Dept B depreciation (15 Jul - 31 Jul): 10M × 17/31 = 5,483,871
- **A2: Transfer with location change only** — Same department, new physical location or responsible user → update location/user, no allocation change.
- **A3: Batch transfer** — User selects multiple assets → assigns new department to all → effective date same → bulk confirm.
- **A4: Transfer with no depreciation impact** — If effective date = start of period or period not yet calculated → simply update allocation for future periods.

### Exception Paths
- **E1: Asset disposed** — Error: "Asset TS-2026-0100 is DISPOSED. Cannot transfer disposed assets."
- **E2: Asset suspended** — Error: "Asset TS-2026-0100 is SUSPENDED. Resume depreciation before transferring." (Or allow with warning that allocation will apply on resume.)
- **E3: Allocation mismatch** — Error: "Total allocation must equal 100%. Current input: 80%."
- **E4: Effective date in closed period** — Block: "Effective date falls in a closed period. Select an open period."
- **E5: Depreciation already posted** — Warning: "Depreciation for current period already posted. Unpost before recalculating allocation."
- **E6: Source and target department same** — No-op: "Asset already assigned to this department. Nothing to update."

### Business Rules
- BR-TRF-01: Transfer does not change asset account (211x) — only allocation changes
- BR-TRF-02: Mid-period transfer = prorata depreciation for source and target (VAS 03 §12 prorata principle)
- BR-TRF-03: Transfer recorded in transaction history for audit trail
- BR-TRF-04: Transfer does not reset acquisition date, useful life, or depreciation method

---

## UC-04: Dispose / Sell Fixed Asset

| Element | Detail |
|---------|--------|
| **ID** | UC-04 |
| **Title** | Dispose / Sell Fixed Asset |
| **Actor** | Chief Accountant (Kế toán trưởng) |
| **Description** | Derecognize a fixed asset through sale, liquidation, donation, or return to lessor. Calculate gain/loss = Net proceeds − Carrying amount. Post GL entries to remove asset and accumulated depreciation. Generate e-invoice if sale. |
| **Disposal types** | Sale (bán — customer invoice + e-invoice), Liquidation (thanh lý — scrap/zero proceeds), Donation (biếu tặng — charitable), Return to lessor (trả lại bên cho thuê — finance lease) |
| **Preconditions** | Asset exists, status is DEPRECIATING or FULLY_DEPRECIATED. User has chief accountant role. GL period is open. |
| **Postconditions** | Asset status = DISPOSED or SOLD. Carrying amount removed from register. Gain/loss recorded. GL posted. E-invoice created (if sale). Transaction recorded. |

### Happy Path
1. Chief Accountant navigates to Fixed Assets → Dispose
2. User selects asset (search by code, name)
3. System displays current FA detail:
   - Original cost, accumulated depreciation, carrying amount
   - Current status, department, location
   - If FULLY_DEPRECIATED: carrying amount = residual value (often 0 or nominal)
4. User selects disposal type: Sell
5. User enters disposal details:
   - Disposal date (effective)
   - Customer (from customer master, if sale)
   - Proceeds amount (sale price before VAT)
   - VAT amount (if applicable, VAT output on disposal per Decree 123)
   - Disposal reason / notes
6. If proceeds > 0: system suggests creating e-invoice (integration with e-invoice module)
   - User optionally generates sales e-invoice with FA details
7. System calculates gain/loss:
   - Gain = Proceeds − Carrying amount (if proceeds > carrying amount)
   - Loss = Carrying amount − Proceeds (if carrying amount > proceeds)
   - Zero gain/loss if carrying amount = proceeds
8. System shows preview:
   - Debit 2141/2143 (accumulated depreciation = total accumulated)
   - Debit 111/112/131 (proceeds received/receivable, VAT separate)
   - Debit 811 (loss, if carrying amount > proceeds) — or —
   - Credit 711 (gain, if proceeds > carrying amount)
   - Credit 211x (original cost removed)
   - Credit 3331 (VAT output on sale, if applicable)
9. User reviews disposal summary
10. User clicks "Confirm Disposal"
11. System validates:
    - Asset not already disposed/sold
    - Disposal date in open period
    - No unposted depreciation for periods before disposal date
12. System posts GL entries (preview above)
13. System updates asset: status = SOLD (if proceeds > 0) or DISPOSED (if proceeds = 0 or liquidation)
14. System records `fixed_asset_transactions` (type: DISPOSAL or SALE)
15. System generates e-invoice (if sale, linked to disposal ID)
16. System returns confirmation: "Asset TS-2026-0100 disposed. Gain: 50,000,000 VND. GL: JE-2026-08912. E-invoice: 01GTGT-2026-12345."

### Alternative Paths
- **A1: Liquidation (zero proceeds)** — Proceeds = 0. No e-invoice needed. GL: Debit 811 (full carrying amount as loss), Debit 2141, Credit 211x. Disposal reason: "Thanh lý do hết hạn sử dụng / hư hỏng."
- **A2: Donation** — Proceeds = 0 but treated as donation. GL: Debit 811 (carrying amount as donation expense), Debit 2141, Credit 211x. VAT output may apply (donation subject to VAT per Decree 123).
- **A3: Returns to lessor (finance lease)** — Remove 212 asset and 341 liability. GL: Debit 341 (lease liability), Debit 2142, Credit 212. Gain/loss if liability ≠ asset carrying amount.
- **A4: Partial disposal** — User disposes a component of a multi-component asset:
  - Enter % or value of component being disposed
  - System calculates proportional original cost and accumulated depreciation
  - Remaining portion continues as separate FA or stays as adjusted original asset
  - New `carrying_amount` = old carrying amount − disposed component carrying amount
- **A5: Disposal before full depreciation** — System auto-calculates un-depreciated carrying amount as loss (if proceeds insufficient to cover). Asset removed before end of useful life.
- **A6: Batch disposal** — Select multiple assets → set disposal type, proceeds, date → bulk confirm → system processes each separately.

### Exception Paths
- **E1: Asset already disposed** — Error: "Asset TS-2026-0100 was disposed on 2026-06-15 (status: DISPOSED). Cannot dispose again."
- **E2: Carrying amount < proceeds (gain scenario)** — Acceptable flow (gain is normal). System flags gain amount for tax treatment (gain subject to CIT).
- **E3: Disposal date in closed GL period** — Block: "Cannot dispose in closed period. Select open period."
- **E4: Missing GL period for mid-year disposal** — Warning: "Disposal date is mid-month. Run depreciation to date before disposing." (System should auto-prompt or auto-calculate partial-month depreciation.)
- **E5: E-invoice generation fails** — Dispose proceeds recorded but e-invoice pending. Status: SOLD (PENDING_EINVOICE). User can retry e-invoice generation.
- **E6: Negative proceeds** — Reject: "Proceeds must be ≥ 0."
- **E7: No customer selected for sale** — Reject: "Customer required for sale disposal type."

### Business Rules
- BR-DSP-01: Disposal gain/loss = Proceeds − Carrying amount at disposal date (VAS 03 §31)
- BR-DSP-02: Asset removed from FA register on disposal date (carrying amount = 0 post-disposal)
- BR-DSP-03: Partial disposal splits asset into disposed component and remaining asset (VAS 03 component approach)
- BR-DSP-04: Sale of FA requires e-invoice per Decree 123 (KT-01GTGT output VAT)
- BR-DSP-05: Disposal loss (811) deductible for CIT; disposal gain (711) taxable
- BR-DSP-06: Donation of FA must comply with charitable donation CIT rules (limited to 10% of taxable income)

---

## UC-05: Adjust Fixed Asset

| Element | Detail |
|---------|--------|
| **ID** | UC-05 |
| **Title** | Adjust Fixed Asset |
| **Actor** | Chief Accountant (Kế toán trưởng) |
| **Description** | Modify FA parameters mid-life: value changes (revaluation, impairment, capitalized improvements), useful life change, depreciation method change, residual value change. All adjustments are prospective per VAS 03 (accounting estimate). Some adjustments require independent appraisal (revaluation). |
| **Adjustment types** | Value increase (revaluation or capitalized improvement), Value decrease (impairment or partial derecognition), Useful life change, Depreciation method change, Residual value change |
| **Preconditions** | Asset exists, status is DEPRECIATING or ACTIVE. User has chief accountant role. Adjustment reason documented. |
| **Postconditions** | Asset parameters updated. Remaining depreciation recalculated. GL adjustment entry posted. Approval trail recorded (if required). |

### Happy Path
1. Chief Accountant navigates to Fixed Assets → Select Asset → Adjust
2. System displays current FA parameters:
   - Original cost, accumulated depreciation, carrying amount
   - Useful life (remaining months), depreciation method
   - Residual value, depreciation start/end dates
   - Full transaction history
3. User selects adjustment type: Value Increase (Revaluation)
4. User enters new value: fair market value per independent appraisal
5. User attaches appraisal report document
6. System compares new value to current carrying amount
7. System calculates adjustment:
   - New carrying amount = appraised fair value
   - Revaluation surplus = new carrying amount − old carrying amount (if increase)
8. System shows preview:
   - Debit 211x (increase amount)
   - Credit 411 (revaluation surplus, equity)
   - Remaining useful life unchanged (or user may adjust separately)
   - Recalculated monthly depreciation after increase
9. User reviews adjustment
10. User confirms
11. System posts approval request (if configured)
12. On approval, system posts GL entry:
    - If increase: Debit 211x, Credit 411
    - If decrease (impairment): Debit 811 (P&L), Credit 211x
13. System updates asset: original cost (if cost model change), carrying amount, accumulated depreciation (may be restated proportionally)
14. System recalculates remaining depreciation schedule prospectively:
    - `NewMonthlyDep = (NewCarryingAmount − ResidualValue) / RemainingMonths`
15. System records `fixed_asset_transactions` (type: ADJUSTMENT, REVALUATION, or IMPAIRMENT)
16. System returns confirmation with recalculated schedule

### Alternative Paths
- **A1: Useful life change** — User selects "Change useful life" → enters new total useful life or remaining months → system validates within legal limits → recalculates depreciation over remaining life. Example: asset originally 120 months, after 36 months user extends to 180 months → remaining 144 months, depreciation recalculated.
- **A2: Depreciation method change** — User selects "Change method" → selects new method (e.g., Declining balance → Straight-line) → provides reason → system validates per VAS 03 §31 (prospective) → recalculates remaining depreciation with new method.
- **A3: Residual value change** — User updates residual value → system recalculates depreciation over remaining life: `(CarryingAmount − NewResidualValue) / RemainingMonths`.
- **A4: Capitalized improvement** — User adds cost of improvement (e.g., engine upgrade) → Debit 2414 (CIP improvement) while work in progress → upon completion → Debit 211x, Credit 2414 — increases original cost — depreciation recalculated over remaining life (or extended life).
- **A5: Multiple adjustments on same asset** — System creates adjustment history version. Each adjustment creates a new depreciation schedule from the adjustment date forward. Prior periods unchanged.
- **A6: Retrospective adjustment (error correction)** — If prior period depreciation was calculated with wrong parameters due to data entry error → system handles as error correction per VAS 03 §32: restate prior periods, post correction entry.

### Exception Paths
- **E1: Adjustment creates negative carrying amount** — Error: "Proposed adjustment would create negative carrying amount (-10,000,000 VND). Cannot reduce value below zero."
- **E2: Value increase exceeds original cost unreasonably** — Warning: "Revaluation increases carrying amount by 300% of original cost. Justification required."
- **E3: Asset disposed or sold** — Block: "Cannot adjust DISPOSED or SOLD assets."
- **E4: Asset fully depreciated** — Warning: "Asset is fully depreciated (carrying amount = residual value). Adjustment may only change residual value or trigger revaluation."
- **E5: Appraisal report missing for revaluation** — Block: "Revaluation requires independent appraisal report per VAS 03."
- **E6: Method change to unsupported method** — Error: "Method 'SUM_OF_YEARS_DIGITS' not supported. Supported: STRAIGHT_LINE, DECLINING_BALANCE, PRODUCTION_BASED."
- **E7: Useful life change violates regulatory limits** — Warning: "New useful life (240 months) exceeds maximum per TT 45 schedule for this category (120 months). Approver override required."

### Business Rules
- BR-ADJ-01: Depreciation method and useful life changes are prospective only (VAS 03 §31)
- BR-ADJ-02: Revaluation requires fair market value evidence (appraisal or market data)
- BR-ADJ-03: Impairment recognized when recoverable amount < carrying amount (IAS 36, VAS 03 §29)
- BR-ADJ-04: Capitalized improvements increase original cost if they extend useful life, increase capacity, or improve quality (VAS 03 §18)
- BR-ADJ-05: Routine repairs expensed immediately, not capitalized (VAS 03 §19)
- BR-ADJ-06: Revaluation surplus (411) cannot be distributed as dividends until realized through depreciation or disposal

---

## UC-06: FA Physical Inventory

| Element | Detail |
|---------|--------|
| **ID** | UC-06 |
| **Title** | FA Physical Inventory (Kiểm kê TSCĐ) |
| **Actor** | Internal Auditor (Kiểm toán nội bộ) / FA Accountant (Kế toán TSCĐ) |
| **Description** | Periodic physical verification of all fixed assets. Compare physical existence and condition against FA register. Record discrepancies: missing, damaged, unregistered assets. Adjust FA register based on verified results. Mandatory annually per Vietnamese accounting law (Luật Kế toán 2015 §36). |
| **Preconditions** | FA register populated with active assets. Company selected. User has inventory permission. |
| **Postconditions** | Inventory plan completed with results. Discrepancies recorded. Register adjusted if required. Audit trail established. |

### Happy Path
1. Internal Auditor navigates to Fixed Assets → Physical Inventory → Create Plan
2. User enters plan details:
   - Plan name / reference (e.g., "Kiểm kê TSCĐ cuối năm 2026")
   - Inventory date (scheduled date for physical count)
   - Scope: all assets or by department/location/category
   - Notes / instructions
3. User saves plan as DRAFT
4. System generates FA inventory list for field use:
   - Asset code, name, serial number, location, department, responsible user
   - Current status (ACTIVE, DEPRECIATING, SUSPENDED, FULLY_DEPRECIATED)
   - Space for physical count results (found/missing/damaged)
5. User prints or exports FA tags/list (Excel/PDF) for physical count
6. User conducts physical count in field
7. User returns to system, navigates to plan → Start Inventory (status → IN_PROGRESS)
8. User records results per asset:
   - For each asset on list: confirm location, condition (Good/Damaged/Needs Repair)
   - Mark as FOUND, MISSING, or DAMAGED
   - Enter actual location (if different from expected)
   - Enter notes
9. All assets accounted for — no discrepancies found
10. User clicks "Complete Inventory"
11. System validates: all assets in scope have results recorded
12. System sets plan status → COMPLETED
13. System generates inventory report:
    - Total assets in scope: X
    - Verified OK: Y
    - Discrepancies: 0
    - Resolution: None required
14. System records completion timestamp and user

### Alternative Paths
- **A1: Discrepancies found (missing asset)** — Asset in register but not found physically:
  - User marks discrepancy type: MISSING
  - System triggers investigation workflow
  - Internal Auditor investigates, reports finding (stolen, misplaced, disposed without record)
  - Resolution options:
    - Asset truly missing → adjust to DISPOSED status with reason, GL adjustment (carrying amount → 811 loss)
    - Temporarily misplaced → mark as INVESTIGATING, re-check in next cycle
- **A2: Discrepancies found (damaged asset)** — Asset found but damaged:
  - User marks discrepancy type: DAMAGED, describes damage extent
  - System assesses whether carrying amount recoverable
  - If impairment indicated → initiate impairment workflow (see UC-05)
  - If repairable → create maintenance request, no register change
- **A3: Unregistered asset found** — Asset exists physically but not in register:
  - User clicks "Add Unregistered Asset" during inventory
  - System opens quick FA creation form (minimal fields)
  - User records: code (temp), name, estimated value, location, condition
  - Asset added to inventory results as UNREGISTERED
  - Post-inventory: user completes full FA registration, source = "Phát hiện qua kiểm kê"
- **A4: Partial inventory** — Scope limited to specific department/location. Only assets in scope are verified. System tracks coverage % of total FA register.
- **A5: Inventory counted over multiple days** — Plan status IN_PROGRESS for duration. Results recorded incrementally. System timestamps each result entry.

### Exception Paths
- **E1: Asset exists in register but status = DISPOSED/SOLD** — Excluded from inventory scope (no physical verification expected). System auto-marks as NOT_APPLICABLE.
- **E2: Inventory list not generated** — Block: "Generate inventory list before starting physical count."
- **E3: Results incomplete at completion** — Error: "Results recorded for 95/120 assets. Complete all entries before finalizing."
- **E4: Duplicate asset physically found (two tags, same asset)** — User marks one as DUPLICATE → system prompts to investigate and merge records.
- **E5: Third-party leased asset counted** — Asset belongs to lessor (finance lease): system marks as LEASED_ASSET, notes for disclosure purposes.

### Business Rules
- BR-INV-01: Physical inventory mandatory at least once per fiscal year (Luật Kế toán 2015 §36)
- BR-INV-02: Inventory scope: all ACTIVE, DEPRECIATING, SUSPENDED, FULLY_DEPRECIATED assets
- BR-INV-03: Missing FA: carrying amount expensed to 811, status → DISPOSED
- BR-INV-04: Unregistered FA: must be registered at estimated fair value per VAS 03
- BR-INV-05: Damaged FA: impairment assessment required within 30 days (TT 45 §12)
- BR-INV-06: Inventory results must be approved by Chief Accountant and stored for audit

---

## UC-07: Import Opening Balance for FA

| Element | Detail |
|---------|--------|
| **ID** | UC-07 |
| **Title** | Import Opening Balance for FA |
| **Actor** | Chief Accountant (Kế toán trưởng) |
| **Description** | Import legacy FA data during system migration. Upload FA register from old system (Excel/CSV). Validate each asset against FA recognition criteria. Create opening GL balance for FA accounts. Must reconcile with opening GL balances. |
| **Preconditions** | Company FA accounts (211x, 214x) have opening balances in GL. Company configured for migration period. User has admin/chief role. |
| **Postconditions** | FA imported with correct status, accumulated depreciation, and carrying amount. FA register matches GL opening balance. Opening balance audit log created. |

### Happy Path
1. Chief Accountant navigates to Fixed Assets → Import → Opening Balance
2. System displays download template link (Excel or CSV with columns)
3. User downloads template — template includes:
   - Standard fields: code, name, category, acquisition date, original cost, accumulated depreciation, residual value, useful life, depreciation method, department, location, serial number, supplier, notes
   - Validation rules embedded (data validation in Excel)
4. User fills template with legacy FA data from old system
5. User uploads completed file
6. System validates each row:
   - Required fields: code, name, category, acquisition date, original cost, useful life
   - Original cost ≥ 30M VND (warning only for opening balance — legacy assets may have been below threshold)
   - Useful life > 12 months
   - Accumulated depreciation ≤ original cost (cannot exceed)
   - Residual value ≤ original cost
   - Carrying amount = original cost − accumulated depreciation ≥ 0
   - Valid category, department (must exist or be created)
   - Duplicate code check
7. System shows validation summary:
   - Total rows: 250
   - Passed: 245
   - Warnings: 3 (cost < 30M for old assets, legacy)
   - Errors: 2 (duplicate code, missing category)
8. User fixes errors, re-uploads, or manually skips failed rows
9. User proceeds to preview
10. System displays import preview: assets grouped by category, total original cost, total accumulated depreciation
11. System calculates: Total carrying amount (imported) vs. GL opening balance (211x − 214x)
    - Displays comparison: "FA import total: 125,678,900,000 VND. GL opening balance (211x-214x): 125,678,900,000 VND. Match: OK."
12. User confirms import
13. System creates FA records with status depending on accumulated depreciation:
    - accumulated_depreciation = 0 → ACTIVE
    - 0 < accumulated_depreciation < original_cost → DEPRECIATING
    - accumulated_depreciation = original_cost → FULLY_DEPRECIATED
14. System does NOT post GL entries (opening balance already in GL from migration)
15. System records `fixed_asset_transactions` (type: ACQUISITION) with note "Opening balance import"
16. System sets depreciation_start_date based on:
    - If fully depreciated: depreciation_start_date = acquisition_date (historical)
    - If still depreciating: depreciation_start_date = earliest prior month needed for schedule
17. System calculates past depreciation history entries if needed (optional, for schedule completeness)
18. System returns import confirmation: "245 FA imported. Total cost: 125B. Accumulated dep: 85B. Carrying amount: 40B."

### Alternative Paths
- **A1: Legacy FA with different depreciation method** — Old system used Sum-of-Years-Digits (not supported). User maps each legacy method to a supported method:
  - Sum-of-Years-Digits → Straight-line (with remaining life recalculated)
  - Accelerated tax method → adjust remaining depreciation to standard method per VAS 03
- **A2: Legacy FA with partially elapsed useful life** — System calculates remaining life = total life − months already depreciated. Depreciation schedule regenerated from import date forward.
- **A3: Assets with no original cost (fully depreciated legacy)** — Original cost recorded historically, accumulated depreciation = original cost, residual value = 0. Status = FULLY_DEPRECIATED.
- **A4: Incremental import** — Company uses system for new FA but has not imported all legacy assets. Can import additional assets in batches. System checks for duplicates each time.

### Exception Paths
- **E1: GL balance mismatch** — Error: "FA import total carrying amount (42,000,000,000 VND) does not match GL opening balance for 211x − 214x (40,000,000,000 VND). Difference: 2,000,000,000 VND."
  - Options:
    - Investigate (show GL account balances vs. FA import summary by account)
    - Post adjustment GL entry to match (Debit 211x, Credit 421 — retained earnings)
    - Force import with mismatch (logged for audit)
- **E2: Duplicate codes across batches** — Error: "Asset code TS-LEGACY-0100 already exists from prior import batch."
- **E3: Category code does not exist** — Error: "Category 'VN006' not found. Create category or map to existing."
- **E4: Accumulated depreciation exceeds original cost** — Error: "Row 15: accumulated depreciation (120M) > original cost (100M). Correct data."
- **E5: File format invalid** — Error: "Could not parse file. Supported formats: .xlsx, .xls, .csv. Ensure column headers match template."

### Business Rules
- BR-IMP-01: Opening balance import does not post GL entries (GL opening balances entered separately)
- BR-IMP-02: Legacy FA with original cost < 30M may be imported with special approval flag
- BR-IMP-03: FA import must reconcile with GL opening balance before go-live
- BR-IMP-04: Import creates audit log: who, when, file hash, row count
- BR-IMP-05: Previously depreciated assets should have depreciation history preserved for schedule
- BR-IMP-06: One-time operation (company goes live once) — no recurring import for opening balance

---

## UC-08: Suspend / Resume FA Depreciation

| Element | Detail |
|---------|--------|
| **ID** | UC-08 |
| **Title** | Suspend / Resume FA Depreciation |
| **Actor** | FA Accountant (Kế toán TSCĐ) |
| **Description** | Temporarily stop depreciation for an asset that is idle, under repair, or not in use. Depreciation resumes when asset returns to service. Per TT 45 §10, suspension ≤ 12 months. Upon resumption, remaining useful life is recalculated to extend depreciation period by suspension duration. |
| **Preconditions** | Asset exists, status = DEPRECIATING. User has FA manage permission. |
| **Postconditions** | Asset status = SUSPENDED (no depreciation during suspension). Upon resume: status → DEPRECIATING, remaining life extended, depreciation recalculated. |

### Happy Path (Suspend)
1. FA Accountant navigates to Fixed Assets → Select Asset → Suspend
2. System displays current FA status (must be DEPRECIATING)
3. User enters suspension details:
   - Effective date (default = today, cannot be in past closed period)
   - Return date (expected, optional)
   - Reason (drop-down or free text):
     - Under repair/renovation (đang sửa chữa, nâng cấp)
     - Idle/temporarily not in use (không sử dụng tạm thời)
     - Awaiting disposal decision (chờ quyết định thanh lý)
     - Other
   - Supporting document (optional)
4. System validates:
   - Asset status = DEPRECIATING
   - Effective date in open period
   - Asset not already suspended
   - If expected return date provided: suspension period ≤ 12 months (TT 45)
5. System shows confirmation:
   - Asset will stop depreciation from effective date
   - No GL entry (suspension is informational — no accounting entry)
   - Remaining useful life at suspension: X months
6. User confirms
7. System updates asset: status → SUSPENDED
8. System records suspension in asset history
9. System calculates suspension end date = suspension date + months of allowed suspension (max 12 months)
10. System stops generating depreciation entries from effective date (if depreciation run, asset skipped)
11. System shows: "Asset suspended from 2026-07-15. No depreciation during suspension."

### Happy Path (Resume)
1. FA Accountant navigates to Fixed Assets → Select Asset → Resume
2. System displays current FA status (must be SUSPENDED)
3. User enters resumption details:
   - Effective date (default = today, or date asset returned to service)
   - Notes
4. System validates:
   - Asset status = SUSPENDED
   - Effective date in open period
   - Suspension duration calculated
5. System calculates:
   - Suspension period: X months (from suspend date to resume date)
   - Remaining useful life BEFORE suspension: Y months
   - New remaining useful life: Y + X months (extension by suspension duration per TT 45 §10)
   - New depreciation end date = old end date + X months
   - Recalculated monthly depreciation: `CarryingAmountAtResume / NewRemainingMonths`
6. System shows preview:
   - Asset: TS-2026-0100
   - Suspension period: 3 months (15-Jul-2026 to 15-Oct-2026)
   - Old remaining life: 60 months, New remaining life: 63 months
   - Old monthly dep: 10,000,000, New monthly dep: 9,523,810 (slightly less per month)
7. User confirms resume
8. System updates asset:
   - status → DEPRECIATING
   - `depreciation_end_date` = new end date
   - (carrying amount unchanged — depreciation did not run during suspension)
   - Extended useful life = used months + new remaining life
9. System recalculates future depreciation schedule
10. System records resume in asset history
11. System shows: "Asset resumed from 2026-10-15. Depreciation recalculated. New monthly: 9,523,810 VND."

### Alternative Paths
- **A1: Suspension exceeds 12 months** — Per TT 45 §10, suspension > 12 months requires special treatment:
  - System warns: "Suspension period exceeds 12 months. Assets suspended > 12 months may need reclassification."
  - If user confirms: re-evaluate whether asset should be disposed instead of suspended
  - Record asset as long-term idle with audit flag
- **A2: Resume with no remaining life** — If asset was near end of useful life when suspended → on resume, remaining months may be very few → system handles naturally by recalculating over remaining short period.
- **A3: Asset put back in service before expected return date** — No issue. Simply resume on actual return date. Suspension duration adjusted.
- **A4: Multiple suspend/resume cycles** — Asset may be suspended and resumed multiple times over its life. System tracks each cycle with dates and duration. Each resume extends remaining life.
- **A5: Suspended asset transferred** — May be allowed or blocked per policy. If allowed: transfer while suspended, allocation applies on resume.

### Exception Paths
- **E1: Asset already suspended** — Error: "Asset TS-2026-0100 is already SUSPENDED since 2026-06-01."
- **E2: Asset disposed or sold** — Error: "Cannot suspend DISPOSED or SOLD assets."
- **E3: Asset fully depreciated** — Warning: "Asset is FULLY_DEPRECIATED. Suspend has no effect (no depreciation remaining)."
- **E4: Resume date before suspend date** — Error: "Resume date cannot precede suspend date."
- **E5: Effective date in closed period** — Block: "Effective date falls in a closed GL period."
- **E6: No depreciation calculated before suspend** — System warns: "Current month depreciation not yet calculated. If running after suspend, asset will be excluded."

### Business Rules
- BR-SUS-01: Depreciation stops on suspension date, resumes on resume date (no pro-rata for partial-month suspension — full month if active at month start)
- BR-SUS-02: Suspension period extends remaining useful life by same duration (TT 45 §10)
- BR-SUS-03: Maximum 12 months continuous suspension (TT 45 §10)
- BR-SUS-04: No GL entry for suspension — only informational status change
- BR-SUS-05: Assets idle > 12 months should be evaluated for impairment or reclassification
- BR-SUS-06: Asset under repair > 12 months → consider capitalizing repair as improvement

---

## UC-09: CIP to FA Transfer (Construction-in-Progress)

| Element | Detail |
|---------|--------|
| **ID** | UC-09 |
| **Title** | CIP to FA Transfer (Chuyển TSCĐ từ XDCB) |
| **Actor** | Chief Accountant (Kế toán trưởng) |
| **Description** | Transfer accumulated costs from Construction-in-Progress account (2411/2412) to Fixed Asset account (211x) upon project completion. Requires completion certificate (Biên bản nghiệm thu). Asset starts depreciation from transfer date. |
| **Preconditions** | CIP account (2411/2412) has accumulated costs for the project. Completion certificate exists. GL period is open. User has chief accountant role. |
| **Postconditions** | CIP balance transferred to FA account. Fixed asset created in FA register. Depreciation starts per asset settings. GL entries posted for transfer. |

### Happy Path
1. Chief Accountant navigates to Fixed Assets → CIP to FA
2. System displays list of CIP projects with accumulated costs in 2411/2412 accounts:
   - Project code, name, total accumulated cost, start date, status
3. User selects project (e.g., "Xây dựng nhà xưởng mới — 2412")
4. System displays CIP detail:
   - Account 2412 total: 25,678,900,000 VND
   - Breakdown by cost type: materials (152), labor (622), subcontractor (331)
   - Accumulated cost over time chart
5. User enters transfer details:
   - Completion date (date of handover)
   - Completion certificate reference number
   - FA category (e.g., Buildings — 2111)
   - FA code and name (auto-suggested: "Nhà xưởng mới")
   - Department responsible
   - Location
   - Useful life (default from category)
   - Depreciation method (default from category)
   - Residual value
   - Depreciation start date (default = completion date or next month)
6. User attaches completion certificate document
7. User selects which CIP costs to transfer:
   - Option: Transfer all costs (full project completion)
   - Option: Draft asset details and verify costs
8. System validates:
   - CIP account balance > 0 (not zero or negative)
   - Completion certificate reference provided
   - FA recognition: total cost ≥ 30M VND, useful life > 12 months
   - No prior CIP-to-FA transfer for this project (prevents double-transfer)
9. System shows preview:
   - Debit 2111 (FA account): 25,678,900,000 VND
   - Credit 2412 (CIP account): 25,678,900,000 VND
   - New FA created: "Nhà xưởng mới" — original cost 25,678,900,000 VND
   - Monthly depreciation (straight-line, 20 years): ~107,000,000 VND
10. User reviews and confirms
11. System posts GL entries:
    - Debit 2111, Credit 2412 (CIP → FA)
12. System creates FixedAsset record:
    - Status: DEPRECIATING (transferred asset starts depreciating)
    - Source: CONSTRUCTION
    - Original cost = total CIP cost transferred
    - Other parameters as specified
13. System records `fixed_asset_transactions` (type: CIP_TRANSFER)
14. System sets up depreciation schedule starting from specified start date
15. System returns confirmation: "CIP project 'Nhà xưởng mới' transferred to FA TS-2026-0200. Original cost: 25.68B. Monthly depreciation: 107M."

### Alternative Paths
- **A1: Partial completion (phase-by-phase transfer)** — Project completed in phases:
  - User selects partial CIP cost amount for this phase
  - System transfers selected amount to FA
  - Remaining CIP balance stays in 2412 for future phases
  - Each phase creates a separate FA or adds to existing FA (component approach)
  - GL: Debit 2111 (partial), Credit 2412 (partial)
  - Example: Building completed 3 of 5 floors → transfer 60% of costs as FA
- **A2: Cost overrun / actual cost exceeds estimate** — System transfers actual costs regardless of original budget. Warning if overrun > 20% of original estimate.
- **A3: Multiple FA from one CIP project** — User splits CIP costs across multiple new FA:
  - E.g., Building: 20B → FA-001, Equipment installed: 5B → FA-002
  - User manually allocates cost breakdown from CIP to each FA
  - System verifies sum of allocated costs = total CIP cost
- **A4: CIP costs include VAT not deducted yet** — System segregates VAT input (1332) from capitalized cost:
  - Capitalized cost = CIP total − VAT amount (if vAT deductible)
  - VAT portion: Debit 1332 (if eligible)
  - If VAT not yet determined: include full cost, adjust later
- **A5: CIP account 2411 vs 2412** — System supports both:
  - 2411: Procurement of FA (mua sắm) — typically single asset
  - 2412: Construction (xây dựng cơ bản) — typically large projects

### Exception Paths
- **E1: CIP account balance is zero** — Error: "CIP account 2412 has zero balance. No costs to transfer."
- **E2: Completion certificate missing** — Block: "Completion certificate reference required. Attach certificate document."
- **E3: CIP already fully transferred** — Error: "Project already transferred to FA on 2026-03-15. Remaining balance: 0 VND."
- **E4: CIP costs exclude necessary components** — Warning: "Total CIP costs (25B) do not include directly attributable costs (transport, installation, testing). These should be added per VAS 03 §8."
- **E5: Useful life or category not set** — Block: "Select FA category and useful life for the new asset."
- **E6: CIP project period not closed in GL** — Warning: "CIP account may have additional costs arriving after transfer. Consider closing the project."

### Business Rules
- BR-CIP-01: FA created from CIP capitalized at total accumulated cost (VAS 03 §8)
- BR-CIP-02: Directly attributable costs: purchase price, import duties, transport, installation, testing, professional fees
- BR-CIP-03: CIP costs exclude: training cost, administration overhead, abnormal waste (VAS 03 §9)
- BR-CIP-04: CIP → FA transfer requires completion certificate (Biên bản nghiệm thu)
- BR-CIP-05: Depreciation starts when asset is available for use, not when CIP is fully paid (VAS 03 §9)
- BR-CIP-06: Partial completion = partial transfer when phase is independently operational (VAS 03 §10)

---

## UC-10: Generate FA Reports

| Element | Detail |
|---------|--------|
| **ID** | UC-10 |
| **Title** | Generate FA Reports |
| **Actor** | Chief Accountant (Kế toán trưởng), CFO (Giám đốc tài chính), External Auditor (Kiểm toán viên) |
| **Description** | Generate standard FA reports per TT 200 / VAS 03 disclosure requirements. Report types: FA Register, Depreciation Schedule, FA Movement, Inventory Report, Aging Report, Department Summary. Exportable to PDF and Excel for submission and audit. |
| **Preconditions** | FA register populated. Depreciation runs completed for relevant periods. Company selected. User has report view permission. |
| **Postconditions** | Report generated in requested format. Temporary report file cached. Download initiated. |

### Happy Path
1. User navigates to Reports → Fixed Assets
2. System displays report type selection:
   - FA Register (Sổ TSCĐ — C03-TS form per TT 200)
   - Depreciation Schedule (Bảng tính khấu hao)
   - FA Movement Report (Báo cáo tăng giảm TSCĐ)
   - FA Inventory Report (Báo cáo kiểm kê TSCĐ)
   - FA Aging Report (TSCĐ theo thời gian sử dụng)
   - FA by Department (TSCĐ theo bộ phận)
3. User selects report type: FA Register
4. User sets filters:
   - Period (month/year): e.g., Jul 2026
   - Status: All / Active / Depreciating / Fully Depreciated (optional)
   - Category: All / by category (optional)
   - Department: All / by department (optional)
   - Include zero-balance assets: Yes/No
5. User clicks "Generate"
6. System loads data for selected period:
   - All FA meeting filter criteria
   - For each FA: code, name, category, acquisition date, original cost, accumulated depreciation to date, carrying amount, monthly depreciation, department, status
7. System calculates totals:
   - Total original cost, total accumulated depreciation, total carrying amount
   - Reconciliation: opening balance + additions − disposals = closing balance
8. System renders report in preview pane
9. User reviews report
10. User clicks "Export"
11. User selects format: PDF (with company logo, signature lines) or Excel (with sortable columns)
12. System generates file, returns download link
13. User downloads report

### Report Types Detail

**FA Register (Sổ TSCĐ — C03-TS)**
- Per TT 200/2014 form C03-TS
- Columns: FA code, name, date of creation, location, original cost, accumulated depreciation, carrying amount, monthly depreciation, notes
- One row per asset, grouped by category
- Header: company name, address, tax code, reporting period
- Footer: prepared by, chief accountant, signature lines

**Depreciation Schedule (Bảng tính khấu hao)**
- Per month/period
- Columns: FA code, name, original cost, carrying amount at start, monthly depreciation, accumulated after, carrying amount after
- Total depreciation per expense account (6274, 6414, 6424)
- Per TT 200 appendix

**FA Movement Report (Báo cáo tăng giảm TSCĐ)**
- Reconciliation format:
  - Opening balance (01/01): original cost, accumulated dep, carrying amount
  - Additions during period: by source (purchase, construction, donation, etc.)
  - Disposals during period: by type (sale, liquidation, donation)
  - Other adjustments: revaluation, impairment, transfers
  - Closing balance (31/12)
- Required for VAS 03 financial statement disclosures

**FA Aging Report**
- Assets grouped by years in service:
  - < 1 year
  - 1-3 years
  - 3-5 years
  - 5-10 years
  - > 10 years
- Helps identify near-fully-depreciated assets, assets past useful life
- Useful for CFO capex planning

**FA by Department**
- Total original cost and carrying amount per department
- Depreciation expense per department for period
- Allows department head cost center analysis

### Alternative Paths
- **A1: Large data set (pagination)** — Report exceeds 5000 rows → system warns: "This report contains 12,500 assets. Generating may take time."
  - Option 1: Generate full report with paginated preview (100 rows per page)
  - Option 2: Chunked export (system writes to temp file, streams download)
  - Option 3: Apply more filters to narrow scope
- **A2: Period not yet run for depreciation** — System warns: "Depreciation for Jul 2026 not yet calculated. Report will show data up to previous period."
  - Allow user to proceed with available data
  - Option: Run depreciation first
- **A3: Cross-period movement report** — User selects year-to-date period (e.g., Jan-Jul 2026):
  - System aggregates movements across selected months
  - Opening balance from start of range, closing balance at end of range
- **A4: Print-ready PDF with signature blocks** — PDF includes:
  - Company header (logo, name, tax code)
  - Prepared by / Chief Accountant signature lines
  - Date of printing
  - "Kính gửi: ..." addressing (for official submissions)

### Exception Paths
- **E1: No assets for selected period** — Warning: "No FA found for the selected filters and period. Adjust filter criteria."
- **E2: Filter returns zero results** — Notification: "No assets match current filters (status = SOLD, department = IT). Expand criteria."
- **E3: Selected period has no depreciation data** — Depreciation reports: "No depreciation entries for Jul 2026. Run depreciation first."
- **E4: Export fails (large file)** — Error: "Report generation exceeded memory limit. Narrow filters or use chunked export."
- **E5: Report template not configured** — Error: "Report template C03-TS not configured. Contact system administrator."
- **E6: FA register unbalanced** — Warning: "Total carrying amount of FA register (40.5B) does not match GL balance (40.2B). Difference: 300M. Report flags for reconciliation."

### Business Rules
- BR-RPT-01: FA Register must show opening balance + movements = closing balance (reconciliation)
- BR-RPT-02: Depreciation Schedule must match total depreciation posted to GL expense accounts
- BR-RPT-03: FA Movement Report required for annual financial statement disclosures (VAS 03)
- BR-RPT-04: All reports must show company identification (name, tax code, address)
- BR-RPT-05: Reports must include "Prepared by" and "Chief Accountant" name/signature fields
- BR-RPT-06: Historical reports frozen at period close — subsequent adjustments do not change closed-period reports

---

## Appendix: Status Transition Matrix

| Transition | From | To | Actor | GL Impact |
|------------|------|----|-------|-----------|
| Register | DRAFT | ACTIVE | Chief Acc. | Debit 211x, Credit 331/111/112 |
| Start depreciation | ACTIVE | DEPRECIATING | System | First depreciation run |
| Suspend | DEPRECIATING | SUSPENDED | FA Acc. | None |
| Resume | SUSPENDED | DEPRECIATING | FA Acc. | None (recalc schedule) |
| Fully depreciate | DEPRECIATING | FULLY_DEPR | System | Carrying amount = residual value |
| Dispose | DEPRECIATING | DISPOSED | Chief Acc. | Debit 2141, Credit 211x, +/- 811/711 |
| Dispose (fully depr) | FULLY_DEPR | DISPOSED | Chief Acc. | Debit 2141, Credit 211x |
| Sell | DEPRECIATING | SOLD | Chief Acc. | As disposal + proceeds entry |
| Sell (fully depr) | FULLY_DEPR | SOLD | Chief Acc. | As disposal + proceeds entry |
| Cancel | DRAFT | CANCELLED | FA Acc. | None (no GL posted) |
| CIP transfer | (CIP) | DEPRECIATING | Chief Acc. | Debit 211x, Credit 2411/2412 |

## Appendix: Regulatory Reference Map

| Regulation | FA Module Impact |
|------------|------------------|
| TT 45/2013 §4 | FA recognition threshold (30M VND, 12 months) |
| TT 45/2013 §10 | Depreciation suspension (max 12 months) |
| TT 45/2013 §11 | Straight-line, declining balance, production-based methods |
| TT 45/2013 §13 | Useful life table by FA category |
| TT 147/2024 | Updated useful life schedules, recognition rules (supersedes TT 45) |
| VAS 03 §8 | Capitalized cost components |
| VAS 03 §9 | CIP — depreciation starts when asset available for use |
| VAS 03 §12 | Depreciation methods — systematic allocation basis |
| VAS 03 §18 | Capitalized improvements vs. routine repairs |
| VAS 03 §29 | Impairment — recoverable amount assessment |
| VAS 03 §31 | Derecognition — gain/loss = proceeds − carrying amount |
| TT 200/2014 | FA register form C03-TS, movement report |
| TT 99/2025 | COA accounts 211-214, 241, 2294, 6274, 6414, 6424 |
| Decree 123/2020 | E-invoice for FA acquisition and disposal |
| IAS 16 | Component depreciation, revaluation model |
| IAS 36 | Impairment indicators, reversal |
| Luật Kế toán 2015 §36 | Mandatory annual physical inventory |
