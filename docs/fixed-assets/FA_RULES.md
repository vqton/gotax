# Fixed Asset Module — Business Rules & Compliance

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** VAS 03, VAS 04, Circular 45/2013/TT-BTC (→ TT 147/2024), Circular 99/2025/TT-BTC, Circular 200/2014/TT-BTC, Circular 133/2016/TT-BTC, Decree 123/2020/ND-CP, IAS 16, IAS 23, IAS 36, IAS 38

---

## Table of Contents

1. [Accounting Rules](#1-accounting-rules-r-01--r-08)
2. [Depreciation Rules](#2-depreciation-rules-r-09--r-14)
3. [Validation Rules](#3-validation-rules-r-15--r-20)
4. [Compliance Rules](#4-compliance-rules-r-21--r-26)
5. [Security Rules](#5-security-rules-r-27--r-30)

---

## 1. Accounting Rules (R-01 – R-08)

---

### R-01: FA Recognition Threshold

| Field | Value |
|-------|-------|
| **Title** | Minimum Recognition Threshold |
| **Description** | Asset qualifies as Fixed Asset only if: (a) original cost ≥ 30,000,000 VND, (b) useful life > 12 months, (c) future economic benefits probable, (d) cost reliably measurable. All 4 conditions must be met simultaneously per VAS 03 para 06 and current MOF guidance. Assets below threshold → record as tools/supplies (TK 153) or expense immediately (TK 623/627/641/642). |
| **Rationale** | Prevents capitalization of low-value items. Ensures consistent FA recognition across enterprises. Aligns with tax audit expectations (CIT inspection typically checks threshold compliance). |
| **Regulatory Reference** | VAS 03 para 06(b)–(d), IAS 16.7, Circular 45/2013 Art 2, Circular 99/2025 Account 211 guidance |
| **Validation Logic** | ```
IF asset.original_cost < 30_000_000 OR asset.useful_life_months <= 12
THEN REJECT "Asset does not meet FA recognition criteria"
IF NOT (future_economic_benefits_probable AND cost_reliably_measurable)
THEN REJECT "VAS 03 para 06 conditions not satisfied"
```
| **Error Message** | `"FA_R01: Original cost must be ≥ 30,000,000 VND and useful life > 12 months per VAS 03"` |

---

### R-02: Original Cost Composition

| Field | Value |
|-------|-------|
| **Title** | Original Cost Determination |
| **Description** | Original cost = purchase price (ex-VAT, after trade discounts) + import duties + non-reclaimable taxes + transport costs + installation costs + professional fees + testing costs + other directly attributable costs to bring asset to working condition. Excludes: administrative overhead, training costs, abnormal waste, internal profits on self-construction. |
| **Rationale** | Ensures complete and consistent capitalization. Directly attributable costs per VAS 03 and IAS 16 must be included; non-qualifying costs must not inflate asset value. Incorrect cost basis misstates depreciation expense and CIT liability. |
| **Regulatory Reference** | VAS 03 para 12–16, IAS 16.16–17, TT 45/2013 Art 4, Circular 99/2025 Account 211 |
| **Validation Logic** | ```
cost_components = sum(
    purchase_price_ex_vat,
    import_duties,
    non_reclaimable_taxes,
    transport,
    installation,
    professional_fees,
    testing
)
excluded = sum(admin_overhead, training, abnormal_waste, internal_profit)
Assert original_cost = cost_components - trade_discounts
Assert excluded NOT IN original_cost
```
| **Error Message** | `"FA_R02: Original cost must include all directly attributable costs per VAS 03 para 12-16"` |

---

### R-03: Capitalizable Borrowing Costs

| Field | Value |
|-------|-------|
| **Title** | Borrowing Cost Capitalization |
| **Description** | Borrowing costs directly attributable to acquisition, construction, or production of a qualifying FA shall be capitalized as part of original cost. Qualifying asset: one that necessarily takes > 12 months to get ready for intended use (e.g., self-constructed buildings, large machinery installation). Capitalization starts when expenditures + borrowing costs incurred and activities commence. Suspended when construction interrupted (extended periods). Ceases when substantially ready for use. |
| **Rationale** | Aligns with IAS 23 (borrowing costs) and VAS 16. Avoids expensing finance costs during construction. Correct treatment required for CIT — capitalized borrowing costs depreciated over asset life rather than expensed immediately. |
| **Regulatory Reference** | IAS 23.8, IAS 23.18–25, VAS 16, Circular 99/2025 Account 241 guidance |
| **Validation Logic** | ```
IF asset.source == "CONSTRUCTION" OR "CIP"
AND estimated_construction_months > 12
AND borrowing_cost_incurred == TRUE
THEN capitalize_borrowing_cost == TRUE
Capitalization period = MAX(activity_start_date, borrowing_start)
Capitalization end = MIN(substantially_ready_date, borrowing_end)
Suspended_periods = periods where construction_interrupted > 3_months
```
| **Error Message** | `"FA_R03: Borrowing costs for qualifying FA (construction >12mo) must be capitalized per IAS 23"` |

---

### R-04: Capital vs. Expense — Maintenance & Improvement

| Field | Value |
|-------|-------|
| **Title** | Capital vs. Expense Distinction |
| **Description** | Subsequent costs categorized as: (a) **Routine maintenance/repairs** → expense immediately (Dr 627/641/642, Cr 111/112/331); (b) **Periodic major maintenance (overhaul)** → capitalize as prepaid expense via TK 242, amortize over period until next overhaul; (c) **Improvements/upgrades** → add to FA original cost (increase TK 211) if they extend useful life, increase capacity, improve product quality, or reduce operating costs. Improvements condition: future economic benefits increase beyond original assessed standard. |
| **Rationale** | Proper classification prevents misstatement of assets vs expenses. Overcapitalization inflates assets and understates expenses; undercapitalization does the opposite. Both distort financial ratios and CIT calculations. TT 45/2013 Art 9 and VAS 03 para 23 provide clear criteria. |
| **Regulatory Reference** | VAS 03 para 23, IAS 16.12–14, TT 45/2013 Art 9, Circular 99/2025 Account 242/2414 guidance |
| **Validation Logic** | ```
IF transaction_type == "REPAIR"
AND NOT (
    extends_useful_life
    OR increases_capacity
    OR improves_quality
    OR reduces_op_cost
)
THEN classify_as_expense = TRUE
ELSE IF extends_useful_life OR increases_capacity OR improves_quality OR reduces_op_cost
THEN classify_as_capital_improvement = TRUE
ELSE IF periodic_major_overhaul
THEN classify_as_prepaid = TRUE (TK 242)
```
| **Error Message** | `"FA_R04: Improvements extending useful life/capacity must be capitalized; repairs expensed per VAS 03 para 23"` |

---

### R-05: Component Depreciation

| Field | Value |
|-------|-------|
| **Title** | Component Depreciation (Significant Parts) |
| **Description** | If a FA comprises significant components with different useful lives or depreciation patterns, each component shall be depreciated separately. Components identified as significant when cost ≥ 20% of total FA cost (configurable threshold). Component approach required for: buildings (roof, HVAC, elevator, electrical), machinery (motor, frame, control system), vehicles (engine, body). Each component tracked as sub-record under parent FA with own useful life, method, and depreciation schedule. |
| **Rationale** | Required per IAS 16.43–47. Prevents distortion where one part wears out faster than the rest. Without component depreciation, replacing a roof (20yr life) on a building (50yr life) would be expensed — wrong. Component approach correctly capitalizes replacement and depreciates over remaining useful life. |
| **Regulatory Reference** | IAS 16.43–47, VAS 03 para 06–07 (implied), Circular 99/2025 Account 211 sub-accounts |
| **Validation Logic** | ```
FOR EACH component IN asset.components:
    IF component.cost / asset.original_cost >= component_significance_threshold
    THEN component.must_depreciate_separately = TRUE
    component.useful_life = min(component.technical_life, asset.remaining_life）
    component.depreciation = component.cost / component.useful_life
    NEXT
Validate sum(component.original_cost) <= asset.original_cost
```
| **Error Message** | `"FA_R05: Significant components (cost ≥20% of total) must have separate useful life and depreciation per IAS 16.43"` |

---

### R-06: Depreciation Method Consistency

| Field | Value |
|-------|-------|
| **Title** | Depreciation Method Uniformity |
| **Description** | Assets of the same type (same category node at level 3) in the same company must use the same depreciation method unless exceptional circumstances are documented and approved by chief accountant. Exception: a specific asset may use a different method if its pattern of economic benefit consumption differs demonstrably from others in the same category (requires justification memo). Method changes are not allowed arbitrarily — only when pattern of benefits changes materially. |
| **Rationale** | Prevents cosmetic manipulation of depreciation expense. Consistent method for same asset class is fundamental accounting principle per IAS 16 and VAS 03. Frequent method changes undermine comparability. Tax authorities (CIT inspection) scrutinize method inconsistency. |
| **Regulatory Reference** | IAS 16.62–66, VAS 03 para 32, TT 45/2013 Appendix 2 |
| **Validation Logic** | ```
category_assets = GET assets WHERE category_id = asset.category_id AND company_id = asset.company_id
category_methods = DISTINCT(category_assets.depreciation_method)
IF len(category_methods) > 1
AND NOT EXISTS(approved_exception_for_asset)
THEN FLAG "Inconsistent depreciation method within category"
```
| **Error Message** | `"FA_R06: All assets of same category must use same depreciation method (chief accountant exception required)"` |

---

### R-07: Residual Value Review

| Field | Value |
|-------|-------|
| **Title** | Annual Residual Value Assessment |
| **Description** | Residual value must be reviewed at least annually (at each fiscal year-end). If expected residual value changes significantly (≥ 20% from current estimate, configurable), depreciation must be recalculated prospectively over remaining useful life. Residual value may be revised to zero if no third-party commitment to purchase and no active market exists for the aged asset. Intangible FA residual value shall be zero unless a third party has committed to purchase or an active market exists (VAS 04 para 62). |
| **Rationale** | Residual value significantly affects depreciation charge. An asset with 30M VND residual value on a 100M cost depreciates 70M; if residual drops to 0, additional 30M must be depreciated. Annual review prevents accumulated error. Tax authorities expect documented residual value assumptions. |
| **Regulatory Reference** | IAS 16.51, IAS 16.77, VAS 03 para 13, VAS 04 para 62 |
| **Validation Logic** | ```
IF current_period == year_end
AND ABS(asset.residual_value - new_residual_estimate) / MAX(asset.residual_value, 1) >= 0.20
THEN recalculate_depreciation_prospectively(asset.id, new_residual_estimate)
```
| **Error Message** | `"FA_R07: Residual value must be reviewed annually per IAS 16.51. Change ≥20% triggers prospective recalculation"` |

---

### R-08: Impairment Indicators and Reversal

| Field | Value |
|-------|-------|
| **Title** | Impairment Assessment & Reversal |
| **Description** | At each reporting date, assess whether impairment indicators exist: (a) external: market value decline > 20%, technological obsolescence, regulatory changes, interest rate changes affecting discount rate; (b) internal: physical damage, asset idled/restructured, economic performance below budget. If indicators present → test impairment: recoverable amount = max(fair value − costs to sell, value in use). If carrying amount > recoverable amount → recognize impairment loss (Dr 811, Cr 2294). Impairment reversal permitted if indicators reverse (except goodwill). Reversal limited to carrying amount had no impairment been recognized (depreciated historical cost). |
| **Rationale** | Assets must not be carried above recoverable amount. Impairment reflects genuine economic loss. Reversal prevents permanent overstatement. IAS 36 is the authoritative standard; VAS references impairment via Circular 228/2009 for listed companies. VT tax law does not allow impairment deduction for CIT — creates temporary difference (deferred tax asset). |
| **Regulatory Reference** | IAS 36.9–14 (indicators), IAS 36.59–64 (reversal), IAS 36.18 (recoverable amount), Circular 228/2009/TT-BTC, Circular 99/2025 Account 2294 |
| **Validation Logic** | ```
indicators = check_external_indicators(asset) OR check_internal_indicators(asset)
IF indicators:
    recoverable = max(fair_value - selling_costs, value_in_use)
    IF recoverable < asset.carrying_amount:
        impairment = asset.carrying_amount - recoverable
        create_impairment_entry(asset.id, impairment)
    ELSE IF previous_impairment AND indicators_reversed:
        reversal = min(
            impairment_recognized,
            depreciated_historical_cost - current_carrying_amount
        )
        IF reversal > 0: reverse_impairment(asset.id, reversal)
```
| **Error Message** | `"FA_R08: Impairment indicators triggered — recoverable amount test required per IAS 36"` |

---

## 2. Depreciation Rules (R-09 – R-14)

---

### R-09: Straight-Line Method Calculation

| Field | Value |
|-------|-------|
| **Title** | Straight-Line Depreciation |
| **Description** | Monthly depreciation = (OriginalCost − ResidualValue) ÷ UsefulLifeMonths. When carrying amount reaches residual value, depreciation stops. Straight-line is the default method under VAS 03 and the most commonly used in Vietnamese enterprises. Must be applied consistently over the asset's useful life unless method change is justified. |
| **Rationale** | Simplest, most auditable method. Assumes economic benefits consumed evenly over time. Default method in TT 45/2013 Appendix 2. |
| **Regulatory Reference** | TT 45/2013 Appendix 2 Method 1, VAS 03 para 32, IAS 16.62 |
| **Validation Logic** | ```
depreciable_amount = asset.original_cost - asset.residual_value
monthly_depreciation = depreciable_amount / asset.useful_life_months
// First period: prorata (R-12)
// Last period: adjustment (R-13)
```
| **Error Message** | `"FA_R09: Straight-line: monthly depreciation = (cost − residual) ÷ useful life months"` |

---

### R-10: Declining Balance Method

| Field | Value |
|-------|-------|
| **Title** | Declining Balance with Conversion Rule |
| **Description** | Annual depreciation rate = (1 ÷ UsefulLifeYears) × DecliningFactor. DecliningFactor = 2.0 for useful life > 4 years, 1.5 for ≤ 4 years (TT 45/2013). Monthly depreciation = CarryingAmount × AnnualRate ÷ 12. **Conversion rule**: when straight-line depreciation on remaining carrying amount exceeds declining balance depreciation, switch to straight-line for remaining life (optional per IAS 16, recommended for IFRS convergence). Declining balance only permitted for assets meeting high-tech/innovation criteria under Vietnamese regulations. |
| **Rationale** | Matches depreciation to higher economic benefit in early years. Conversion rule ensures full depreciation over useful life (pure declining balance never reaches zero). TT 45 limited to qualifying assets — must check eligibility. |
| **Regulatory Reference** | TT 45/2013 Appendix 2 Method 2, IAS 16.62, VAS 03 para 32, Circular 99/2025 |
| **Validation Logic** | ```
years = asset.useful_life_months / 12
factor = IF years > 4 THEN 2.0 ELSE 1.5
annual_rate = (1.0 / years) * factor
monthly_rate = annual_rate / 12
monthly_depreciation = asset.carrying_amount * monthly_rate
// Conversion check:
sl_depr = (asset.carrying_amount - asset.residual_value) / remaining_months
db_depr = asset.carrying_amount * monthly_rate
IF sl_depr > db_depr: switch_to_straight_line(asset.id)
```
| **Error Message** | `"FA_R10: Declining balance factor = 2.0 (>4yr) or 1.5 (≤4yr). Convert to SL when SL > DB per IAS 16"` |

---

### R-11: Production-Based Method

| Field | Value |
|-------|-------|
| **Title** | Production (Units-of-Activity) Depreciation |
| **Description** | Depreciation per unit = (OriginalCost − ResidualValue) ÷ TotalEstimatedProduction. Monthly depreciation = DepreciationPerUnit × ActualProductionInMonth. TotalEstimatedProduction must be set at acquisition and reviewed annually. Changes in estimate (updated total capacity) → applied prospectively. Method appropriate for mining equipment, power generators, vehicles (km-based), and machinery where wear correlates with usage rather than time. |
| **Rationale** | Most accurate when asset usage varies significantly between periods. Matches depreciation expense to actual economic benefit consumption per VAS 03 and IAS 16. |
| **Regulatory Reference** | TT 45/2013 Appendix 2 Method 3, IAS 16.62, VAS 03 para 32 |
| **Validation Logic** | ```
depr_per_unit = (asset.original_cost - asset.residual_value) / asset.total_estimated_production
monthly_depr = depr_per_unit * actual_production_in_month
IF actual_production > asset.total_estimated_production:
    monthly_depr = asset.carrying_amount - asset.residual_value  // fully depreciate
```
| **Error Message** | `"FA_R11: Production-based: monthly dep = (cost − residual) × (actual production / total estimated)"` |

---

### R-12: First Period Prorata Temporis

| Field | Value |
|-------|-------|
| **Title** | First-Period Partial Depreciation (Prorata Temporis) |
| **Description** | In the first period (month) depreciation starts, calculate depreciation proportionally based on days remaining in the month from the depreciation start date. Formula: MonthlyDepreciation × (DaysRemaining ÷ TotalDaysInMonth). Depreciation start date = date asset is available for use (acquisition date or later, if commissioning period required). If start date = 1st of month → full month charged. |
| **Rationale** | Fair allocation — asset not available for full month, so charging full month overstates expense. Aligns with VAS 03 para 30 (depreciation starts when asset available for use) and IAS 16.55. Many Vietnamese enterprises (especially SMEs) charge full month regardless — this rule supports the more precise prorata approach. |
| **Regulatory Reference** | VAS 03 para 30, IAS 16.55, TT 45/2013 Art 10 |
| **Validation Logic** | ```
start_date = asset.depreciation_start_date
days_in_month = days_in_month(start_date.year, start_date.month)
days_remaining = days_in_month - start_date.day + 1
prorata_factor = days_remaining / days_in_month
first_depreciation = monthly_depreciation * prorata_factor
```
| **Error Message** | `"FA_R12: First period depreciation prorated: (monthly amount) × (days remaining / total days in month)"` |

---

### R-13: Last Period Adjustment

| Field | Value |
|-------|-------|
| **Title** | Last-Period Remaining Balance Adjustment |
| **Description** | In the final depreciation period (when remaining useful life ends or earlier if disposed), adjust depreciation amount to ensure carrying amount equals residual value (typically 0). Formula: LastPeriodAmount = CarryingAmount(period start) − ResidualValue. This may differ from the standard monthly amount. No negative depreciation allowed (cannot create negative expense). |
| **Rationale** | Ensures carrying amount converges precisely to residual value. Without this adjustment, rounding differences accumulate over long useful lives and carrying amount does not reach residual value. |
| **Regulatory Reference** | IAS 16.50–51, VAS 03 para 31, standard accounting practice |
| **Validation Logic** | ```
IF (period_index + 1) >= asset.useful_life_months  // last or second-to-last
OR asset.carrying_amount - monthly_depr < asset.residual_value:
    last_depreciation = asset.carrying_amount - asset.residual_value
    last_depreciation = MAX(last_depreciation, 0)
```
| **Error Message** | `"FA_R13: Last period depreciation adjusted to remaining balance: carrying amount − residual value"` |

---

### R-14: Depreciation Suspension Rules

| Field | Value |
|-------|-------|
| **Title** | Suspension and Resumption of Depreciation |
| **Description** | Depreciation may be suspended when FA is: (a) idle/not in use for ≥ 1 month (e.g., machinery under repair, vehicle in long-term shop, building under renovation), (b) temporarily out of service. Suspension limits: (1) maximum suspension period = 12 consecutive months per TT 45/2013 Art 10; (2) no depreciation charged during suspension; (3) upon resumption, remaining useful life may be recalculated (extend by suspension duration or keep original end date — chosen at suspension creation, requires chief accountant approval); (4) cannot suspend if asset is FULLY_DEPRECIATED or DISPOSED; (5) cannot suspend retroactively (suspension effective from current or future period only). |
| **Rationale** | Depreciation should not be charged when asset is not generating economic benefits. Suspension > 12 months suggests asset may be impaired (trigger R-08 assessment). Prevents companies from "suspending" depreciation to manipulate earnings. |
| **Regulatory Reference** | TT 45/2013 Art 10, VAS 03 para 30, TT 147/2024 (replacement) |
| **Validation Logic** | ```
IF asset.status NOT IN ("DEPRECIATING"):
    REJECT "Only DEPRECIATING assets can be suspended"
IF suspension_period_exists AND suspension_period > 12_months:
    REJECT "Max suspension is 12 consecutive months per TT 45 Art 10"
suspended_months = count_consecutive_suspension_months(asset.id)
IF suspended_months >= 12:
    TRIGGER impairment_assessment(asset.id)
resume.remaining_life = IF resume.option == "extend"
    THEN asset.remaining_useful_life + suspended_months
    ELSE asset.remaining_useful_life
```
| **Error Message** | `"FA_R14: Suspension max 12 months. Cannot suspend FULLY_DEPRECIATED/DISPOSED assets. No retroactive suspension"` |

---

## 3. Validation Rules (R-15 – R-20)

---

### R-15: Asset Code Uniqueness

| Field | Value |
|-------|-------|
| **Title** | Asset Code Unique Within Company |
| **Description** | FA code must be unique within the same company. Code format: alphanumeric, max 50 characters, configurable per company (auto-generate or manual entry). If auto-generate, default format: `FA-{YYYY}-{NNNNNN}` (year + 6-digit sequential). Code cannot be changed after asset status changes from DRAFT. |
| **Rationale** | Asset code is the primary identifier for FA register, tax filings, and audit trail. Duplicate codes cause reconciliation failures and misidentification. Uniqueness constraint enforced at database level (UNIQUE(company_id, code) on fixed_assets table). |
| **Regulatory Reference** | VAS 03 para 39 (disclosure — FA by class), standard accounting practice |
| **Validation Logic** | ```
existing = SELECT COUNT(*) FROM fixed_assets
    WHERE company_id = @company_id AND code = @code AND id != @asset_id
IF existing > 0: REJECT
```
| **Error Message** | `"FA_R15: Asset code '{{code}}' already exists in company {{company_id}}"` |

---

### R-16: Category Level Depth ≤ 3

| Field | Value |
|-------|-------|
| **Title** | FA Category Hierarchy Depth Limit |
| **Description** | FA category hierarchy is limited to 3 levels: Level 1 = Group (e.g., Buildings, Machinery, Vehicles), Level 2 = Type (e.g., Industrial Buildings, Office Buildings), Level 3 = Category/Category Detail (e.g., Steel-frame Factory, Concrete Office). Circular reference detection required (category cannot be its own ancestor). Parent category must exist and have level = current_level − 1. |
| **Rationale** | Standard enterprise FA classification has 3 levels. Deeper hierarchies add complexity without benefit. Circular references cause infinite loops in tree traversal. |
| **Regulatory Reference** | TT 45/2013 Art 3 (6 FA groups), IAS 16.37 (class of PP&E), standard ERP practice (MISA, Fast, Bravo) |
| **Validation Logic** | ```
IF NEW.level > 3:
    REJECT "Category depth max is 3 levels"
IF NEW.parent_id IS NOT NULL:
    parent = SELECT level FROM fixed_asset_categories WHERE id = @parent_id
    IF parent.level != (NEW.level - 1):
        REJECT "Parent category level must be (current level - 1)"
// Circular reference check:
ancestors = []
current = NEW.parent_id
WHILE current IS NOT NULL:
    IF current IN ancestors: REJECT "Circular reference detected"
    ancestors.append(current)
    current = SELECT parent_id FROM fixed_asset_categories WHERE id = @current
```
| **Error Message** | `"FA_R16: Category hierarchy max 3 levels. Circular references prohibited"` |

---

### R-17: Disposal State Constraint

| Field | Value |
|-------|-------|
| **Title** | Disposal Requires DEPRECIATING or FULLY_DEPRECIATED State |
| **Description** | Disposal (dispose/sell) is only permitted when asset status is DEPRECIATING or FULLY_DEPRECIATED. Assets in DRAFT, CANCELLED, ACTIVE (not yet depreciating), or SUSPENDED states cannot be disposed. Exception: asset in SUSPENDED state may be disposed only with chief accountant override (requires documented justification). |
| **Rationale** | An asset that has never been depreciated (ACTIVE without depreciation start) cannot be disposed — this indicates incorrect lifecycle handling. Disposal from SUSPENDED state is valid (asset may never resume) but requires approval. Prevents erroneous disposal before asset is properly registered and depreciated. |
| **Regulatory Reference** | VAS 03 para 37–38, TT 45/2013 Art 11 |
| **Validation Logic** | ```
IF asset.status NOT IN ("DEPRECIATING", "FULLY_DEPRECIATED"):
    IF asset.status == "SUSPENDED" AND request.has_chief_accountant_override:
        ALLOW
    ELSE:
        REJECT
```
| **Error Message** | `"FA_R17: Disposal only allowed for DEPRECIATING or FULLY_DEPRECIATED assets. SUSPENDED assets require chief accountant override"` |

---

### R-18: Deletion Constraint — Depreciation Entries Exist

| Field | Value |
|-------|-------|
| **Title** | Cannot Delete Asset with Depreciation History |
| **Description** | An asset with any posted depreciation entries cannot be deleted. Non-posted (draft) entries block deletion as well. The asset must be disposed (status = DISPOSED or SOLD) rather than deleted. DRAFT assets with zero depreciation entries may be hard-deleted. ACTIVE assets (depreciation not started) may be cancelled (status = CANCELLED) but not hard-deleted. |
| **Rationale** | Depreciation entries are financial postings — deleting the asset would orphan GL references and break audit trail. Vietnamese accounting law requires 10-year record retention. Soft-delete (dispose) preserves history. |
| **Regulatory Reference** | Accounting Law Article 12 (document retention), VAS 03 para 37–38, standard audit requirements |
| **Validation Logic** | ```
depr_count = SELECT COUNT(*) FROM depreciation_entries WHERE fixed_asset_id = @id
IF depr_count > 0:
    REJECT "Asset has depreciation history — use disposal rather than deletion"
IF asset.status == "DRAFT":
    ALLOW hard_delete
ELSE IF asset.status == "ACTIVE":
    ALLOW cancel (status = CANCELLED), NOT hard_delete
ELSE:
    REJECT "Asset must be disposed, not deleted"
```
| **Error Message** | `"FA_R18: Asset with depreciation entries cannot be deleted. Use disposal workflow"` |

---

### R-19: Allocation Percentage Sum = 100%

| Field | Value |
|-------|-------|
| **Title** | Multi-Department Allocation Constraint |
| **Description** | When an FA is allocated across multiple departments, the sum of allocation percentages must equal exactly 100%. Each allocation record must have allocation_pct between 1% and 100%. A department may appear at most once per asset. If no allocations exist, the full depreciation is charged to the department specified on the asset's expense_account_id. |
| **Rationale** | Allocations represent how depreciation expense is distributed. Sum ≠ 100% means expense is either over- or under-allocated, causing incorrect department P&L. |
| **Regulatory Reference** | Circular 99/2025 Account 6274/6414/6424, management accounting best practice |
| **Validation Logic** | ```
total_pct = SELECT SUM(allocation_pct) WHERE fixed_asset_id = @id
IF total_pct IS NOT NULL AND total_pct != 100.00:
    REJECT "Allocation percentages must sum to 100% (current: " + total_pct + ")"
FOR EACH allocation:
    IF allocation.allocation_pct < 1 OR allocation.allocation_pct > 100:
        REJECT "Each allocation percentage must be between 1% and 100%"
```
| **Error Message** | `"FA_R19: Allocation percentages must sum to 100% (currently {{total}}%)"` |

---

### R-20: Adjustment Requires Chief Accountant Approval

| Field | Value |
|-------|-------|
| **Title** | FA Adjustment Approval Workflow |
| **Description** | Any adjustment to FA value, useful life, depreciation method, or residual value requires chief accountant approval before posting. Adjustments include: value increase (capitalized improvement, revaluation), value decrease (impairment, partial disposal), useful life change, method change, residual value change. Adjustment flow: FA accountant submits → chief accountant reviews → approved/rejected → posted. Rejected adjustments require mandatory reason. No retroactive adjustments (effective from current period forward, prospective). |
| **Rationale** | FA adjustments affect financial statements and tax calculations. Unauthorized adjustments could manipulate earnings. Two-step approval provides control. Prospective application per IAS 16.65–66 prevents restating prior periods. |
| **Regulatory Reference** | IAS 16.65–66 (prospective), VAS 03 para 33–34, TT 45/2013 Art 9, internal control standards |
| **Validation Logic** | ```
IF request.role NOT IN ("chief_accountant", "admin"):
    REJECT "Only chief accountant can approve FA adjustments"
// Validate prospective only
IF adjustment.effective_date < current_period:
    REJECT "Adjustments must be prospective, not retroactive"
// Method change
IF adjustment.field == "depreciation_method":
    doc = GET justification_document(asset.id)
    IF NOT doc: REJECT "Method change requires documented justification"
```
| **Error Message** | `"FA_R20: FA adjustment requires chief accountant approval. Adjustments are prospective only"` |

---

## 4. Compliance Rules (R-21 – R-26)

---

### R-21: FA Register Form (TT 200/2014)

| Field | Value |
|-------|-------|
| **Title** | FA Register Format Compliance |
| **Description** | The FA Register (Sổ Tài sản cố định) must follow the mandatory form prescribed in Circular 200/2014/TT-BTC (Appendix — Sổ TSCĐ). Required fields: asset code, name, category, acquisition date, original cost, accumulated depreciation, carrying amount, depreciation method, useful life, department/location. The register must be printable as a unified document covering all FA by category. Each page must show: company name, report title, period, page number, totals carried forward. |
| **Rationale** | Mandatory form for all Vietnamese enterprises. Tax authorities and auditors expect the standard format during inspection. Non-standard format may result in non-compliance finding. |
| **Regulatory Reference** | TT 200/2014 Appendix — FA Register form, TT 133/2016 Appendix (SME simplified form), VAS 03 para 39–40 |
| **Validation Logic** | ```
report = build_fa_register(company_id, period)
FOR EACH asset IN report:
    ASSERT asset.code IS NOT NULL
    ASSERT asset.original_cost IS NOT NULL
    ASSERT asset.carrying_amount == asset.original_cost - asset.accumulated_depreciation
report.pages.each(page_number += 1, show_totals_carried_forward)
```
| **Error Message** | `"FA_R21: FA Register must follow TT 200/2014 mandatory form — all required fields present"` |

---

### R-22: Depreciation Schedule Printable per VAS 03

| Field | Value |
|-------|-------|
| **Title** | Depreciation Schedule Disclosure |
| **Description** | Depreciation schedule must be printable per period showing: asset code, name, opening carrying amount, depreciation amount for period, accumulated depreciation to date, closing carrying amount. Schedule must show subtotals by FA category and grand total. Used for VAS 03 para 39 disclosure (reconciliation of carrying amount at beginning and end of period). Must support export to PDF and Excel for FS preparation. |
| **Rationale** | VAS 03 para 39 requires reconciliation disclosure: opening balance + additions − disposals − depreciation − impairment = closing balance. The depreciation schedule is the primary source document. |
| **Regulatory Reference** | VAS 03 para 39(d), IAS 16.73(e), Circular 99/2025 BCTC forms |
| **Validation Logic** | ```
schedule = get_depreciation_schedule(company_id, period_start, period_end)
total_opening = SUM(schedule.opening_carrying_amount)
total_depreciation = SUM(schedule.depreciation_amount)
total_closing = SUM(schedule.closing_carrying_amount)
ASSERT total_closing == total_opening - total_depreciation + additions - disposals
```
| **Error Message** | `"FA_R22: Depreciation schedule must reconcile per VAS 03 para 39: opening − depreciation ± movements = closing"` |

---

### R-23: FA Movement Report for Annual FS

| Field | Value |
|-------|-------|
| **Title** | FA Movement Schedule (VAS 03 Para 55-60) |
| **Description** | FA Movement Report (Báo cáo tăng giảm TSCĐ) is required for annual financial statement disclosures. Must show per FA category: opening balance (cost + accumulated dep + carrying amount), additions (cost), disposals (cost + accumulated dep), depreciation for period, other changes (revaluation, impairment), closing balance. Used for VAS 03 para 39 reconciliation disclosure. Must be signed by preparer and chief accountant. |
| **Rationale** | VAS 03 mandatorily requires movement schedule. Without this report, annual FS are incomplete. Auditors will qualify opinion if FA movement not disclosed. |
| **Regulatory Reference** | VAS 03 para 39 (paras 55-60 for detailed guidance), IAS 16.73, TT 200/2014 FA report forms |
| **Validation Logic** | ```
movement = build_movement_report(company_id, fiscal_year)
FOR EACH category:
    ASSERT movement.opening_cost + movement.additions - movement.disposals == movement.closing_cost
    ASSERT movement.opening_depr + movement.depreciation - movement.disposals_depr == movement.closing_depr
    ASSERT movement.closing_cost - movement.closing_depr == movement.closing_carrying
```
| **Error Message** | `"FA_R23: FA Movement Report required for annual FS per VAS 03 para 39. Must reconcile opening ↔ closing"` |

---

### R-24: E-Invoice Required for FA Sale (Decree 123/2020)

| Field | Value |
|-------|-------|
| **Title** | E-Invoice Mandatory on FA Disposal with Sale Proceeds |
| **Description** | When a FA is sold (not liquidated/scrapped), an e-invoice must be issued per Decree 123/2020 Art 4 and Art 9. The e-invoice must include: seller/buyer tax codes, asset description, quantity (1 unit), unit price (sale price), total amount, VAT rate (typically 10% for used FA), VAT amount. Invoice timing = time of ownership/use rights transfer to buyer (Art 9.1). System must allow generating e-invoice from disposal transaction. Exception: FA contributed as capital to establish enterprise → use capital contribution minutes + asset delivery records (no e-invoice per Art 9.2(e)). |
| **Rationale** | Decree 123 mandates e-invoices for all goods/services sales, including FA. Tax authority requires e-invoice data transmission. Selling without e-invoice = tax evasion. |
| **Regulatory Reference** | Decree 123/2020 Art 4, Art 9.1, Art 9.2(e), Art 12; Circular 99/2025 Account 711/811 |
| **Validation Logic** | ```
IF disposal.type == "SALE":
    e_invoice = generate_e_invoice(
        buyer_tax_code = disposal.buyer_tax_code,
        asset_description = asset.name,
        quantity = 1,
        unit_price = disposal.sale_price,
        vat_rate = determine_vat_rate(asset),
        invoice_type = "FA_SALE"
    )
    IF NOT e_invoice:
        REJECT "E-invoice required for FA sale per Decree 123 Art 4"
ELSE IF disposal.type == "LIQUIDATION" OR "DONATION":
    // No e-invoice required, use liquidation record
```
| **Error Message** | `"FA_R24: FA sale requires e-invoice per Decree 123/2020 Art 4. Capital contribution uses minutes + delivery record per Art 9.2(e)"` |

---

### R-25: FA Purchase Invoice for VAT Deduction

| Field | Value |
|-------|-------|
| **Title** | Purchase Invoice Requirements for FA Acquisition |
| **Description** | For FA purchased from suppliers, the invoice must satisfy following conditions for VAT input deduction: (a) valid e-invoice from supplier per Decree 123; (b) invoice shows supplier name, tax code, asset description, serial number (if applicable), quantity, unit price, total amount, VAT rate, VAT amount; (c) for invoices > 20,000,000 VND — payment must be made via bank transfer (non-cash payment evidence required per Circular 99 guidance); (d) VAT rate must match applicable rate for asset type (typically 10% for FA, 8% for certain equipment per VAT reduction policies). |
| **Rationale** | VAT input on FA purchase is a significant recoverable amount. Incorrect invoice → denied VAT deduction → increased cost. Decree 123 and VAT law require strict compliance. Payment > 20M via bank is mandatory for deduction. |
| **Regulatory Reference** | Decree 123/2020 Art 12, Law on VAT 13/2024/QH15, Circular 99/2025 Account 1332, Decree 123 non-cash payment requirement |
| **Validation Logic** | ```
IF asset.source == "PURCHASE" AND asset.original_cost > 20_000_000:
    IF NOT invoice.bank_payment_evidence:
        FLAG "VAT deduction risk: payment >20M requires bank transfer evidence"
IF invoice.vat_rate NOT IN (applicable_rates_for_asset):
    REJECT "Invoice VAT rate does not match applicable rate for asset type"
```
| **Error Message** | `"FA_R25: FA purchase invoice >20M VND requires bank payment for VAT deduction. Validate VAT rate matches asset type"` |

---

### R-26: CIT Adjustment for Depreciation Differences

| Field | Value |
|-------|-------|
| **Title** | Tax vs. Accounting Depreciation CIT Adjustment |
| **Description** | When tax depreciation (per CIT law, Circular 78/2014/TT-BTC, Circular 151/2014) differs from accounting depreciation (per VAS 03/TT 45), a CIT adjustment is required at tax declaration time. Typical differences: (a) useful life shorter for tax → tax depreciation > accounting → temporary difference (deferred tax liability); (b) useful life longer for tax → accounting depreciation > tax → temporary difference (deferred tax asset); (c) method difference (e.g., declining balance for book, straight-line for tax); (d) impairment recognized for accounting but not deductible for tax (permanent difference if never deductible). System must track: accounting depreciation, tax depreciation, difference, temporary/permanent classification. |
| **Rationale** | Vietnamese CIT law has its own depreciation rules that may differ from accounting. Enterprises must file CIT based on tax depreciation. The difference creates deferred tax. Non-compliance → incorrect CIT payment → penalties. |
| **Regulatory Reference** | Circular 78/2014/TT-BTC Art 4 (depreciation for CIT), Circular 151/2014, Law on CIT 14/2008/QH12 (amended), IAS 12 (deferred tax), VAS 17 (deferred tax) |
| **Validation Logic** | ```
accounting_depr = get_period_depreciation(asset.id, period)
tax_depr = calculate_tax_depreciation(asset.id, period, tax_useful_life, tax_method)
diff = accounting_depr - tax_depr
IF diff > 0:
    // Accounting > Tax → deferred tax asset (future deductible)
    deferred_tax = diff * current_cit_rate
    IF difference_type == "TEMPORARY":
        record_deferred_tax_asset(deferred_tax)
    ELSE:
        // Permanent difference (e.g., impairment) → never deductible
        add_to_cit_adjustment_schedule(asset.id, diff, "PERMANENT")
ELSE IF diff < 0:
    // Tax > Accounting → deferred tax liability
    record_deferred_tax_liability(abs(diff) * current_cit_rate)
```
| **Error Message** | `"FA_R26: CIT adjustment required if tax depreciation differs from accounting. Accrual >= 0, IFRS < 0, but our code uses GAAP"` |

---

## 5. Security Rules (R-27 – R-30)

---

### R-27: Role-Based Access Control

| Field | Value |
|-------|-------|
| **Title** | FA Module Role Permissions |
| **Description** | Access to FA operations is controlled by user role. Minimum permission matrix: |

| Operation | Viewer | FA Accountant | Chief Accountant | Admin |
|-----------|--------|---------------|------------------|-------|
| View FA register | ✓ | ✓ | ✓ | ✓ |
| View FA details | ✓ | ✓ | ✓ | ✓ |
| View reports | ✓ | ✓ | ✓ | ✓ |
| Create FA (DRAFT) | ✗ | ✓ | ✓ | ✓ |
| Edit FA (DRAFT/ACTIVE) | ✗ | ✓ | ✓ | ✓ |
| Activate FA | ✗ | ✗ | ✓ | ✓ |
| Set allocation | ✗ | ✓ | ✓ | ✓ |
| Transfer FA | ✗ | ✓ | ✓ | ✓ |
| Suspend/resume dep | ✗ | ✓ | ✓ | ✓ |
| Adjust FA value/life | ✗ | ✗ | ✓ | ✓ |
| Approve adjustments | ✗ | ✗ | ✓ | ✓ |
| Dispose/sell FA | ✗ | ✗ | ✓ | ✓ |
| Run depreciation | ✗ | ✗ | ✓ | ✓ |
| Post depreciation to GL | ✗ | ✗ | ✓ | ✓ |
| Unpost depreciation | ✗ | ✗ | ✗ | ✓ |
| Create FA category | ✗ | ✗ | ✓ | ✓ |
| Import/export FA | ✗ | ✗ | ✗ | ✓ |
| Override blocked ops | ✗ | ✗ | ✓ | ✓ |

**Rationale:** Segregation of duties prevents unauthorized FA transactions. FA accountant manages daily ops; chief accountant controls approvals and GL impact; admin handles system-level changes. |
| **Regulatory Reference** | Accounting Law Article 12 (internal control), Circular 99/2025, standard enterprise internal control |
| **Validation Logic** | ```
permissions = GET role_permissions(current_user.role)
IF request.operation NOT IN permissions:
    RETURN 403 "Insufficient permissions"
```
| **Error Message** | `"FA_R27: {{role}} does not have permission for {{operation}}. Required role: {{required_role}}"` |

---

### R-28: Audit Trail

| Field | Value |
|-------|-------|
| **Title** | FA Transaction Audit Trail |
| **Description** | Every FA transaction must be logged immutably in the `fixed_asset_transactions` table. Required fields per transaction: timestamp (auto), user ID, transaction type, asset ID, old value, new value, description, GL journal reference (if applicable). Immutable: once created, transactions cannot be modified or deleted (only compensating entries allowed). The following events are always logged: acquisition, activation, depreciation run, adjustment (each field change logged separately), transfer (from/to department), suspension, resumption, disposal, sale, revaluation, impairment, CIP transfer, allocation change, category change, cost component change. |
| **Rationale** | Audit trail is legally required by Accounting Law Article 12. Immutability ensures evidence integrity. Each transaction type maps to a distinct FATransactionType constant in the data model (ACQUISITION, ADJUSTMENT, TRANSFER, DISPOSAL, SALE, REVALUATION, IMPAIRMENT, CIP_TRANSFER, DEPRECIATION). |
| **Regulatory Reference** | Accounting Law 2015 Article 12 (document retention), IAS 16 (disclosure), VAS 03 para 39, Circular 99/2025 |
| **Validation Logic** | ```
transaction = create_transaction(
    fixed_asset_id = asset.id,
    transaction_type = determine_type(operation),
    transaction_date = NOW(),
    amount = amount_change,
    old_value = OLD(affected_field),
    new_value = NEW(affected_field),
    description = human_readable_summary,
    created_by = current_user.id
)
ASSERT transaction.id IS NOT NULL
ASSERT transaction IS NOT UPDATEABLE, NOT DELETABLE
```
| **Error Message** | `"FA_R28: Every FA transaction must be logged immutably in fixed_asset_transactions"` |

---

### R-29: Period Lock — Cannot Post Depreciation to Closed Period

| Field | Value |
|-------|-------|
| **Title** | Depreciation Period Lock Protection |
| **Description** | Depreciation entries cannot be created or posted for a period that is locked or closed. A period is closed when: (a) period status = CLOSED in the periods table; (b) period-end checklist completed; (c) financial statements for that period finalized. Exception: only admin-level role can unlock a period for corrections (must be logged and documented). Depreciation run validates all periods in range before processing. Cannot create depreciation entry with period_id referencing a closed period. |
| **Rationale** | Posting to closed periods undermines period-end integrity. GL balances would change after FS release. Accounting law requires period cutoff. |
| **Regulatory Reference** | Circular 99/2025 period management, Circular 200/2014 period closing procedures, Accounting Law |
| **Validation Logic** | ```
period = SELECT status FROM periods WHERE id = @period_id
IF period.status IN ("LOCKED", "CLOSED"):
    REJECT "Cannot post depreciation to " + period.status + " period"
// Batch validation
FOR EACH period_id IN request.period_ids:
    period_status = SELECT status FROM periods WHERE id = @period_id
    IF period_status NOT IN ("OPEN", "ADJUSTMENT"):
        REJECT "Period " + period_id + " is " + period_status + " — unavailable for depreciation"
```
| **Error Message** | `"FA_R29: Cannot post depreciation to {{status}} period. Open or unlocked periods only"` |

---

### R-30: Two-Step Approval for Adjustments and Disposals

| Field | Value |
|-------|-------|
| **Title** | Two-Step Approval Workflow |
| **Description** | The following FA operations require a 2-step approval process (submit → approve → execute):

| Operation | Step 1 (Submit) | Step 2 (Approve) |
|-----------|----------------|------------------|
| FA value adjustment | FA Accountant | Chief Accountant |
| Useful life change | FA Accountant | Chief Accountant |
| Depreciation method change | FA Accountant | Chief Accountant |
| Above-threshold disposal (> 500M VND carrying amount) | FA Accountant + Dept Head | Chief Accountant + Admin |
| Impairment recognition | FA Accountant | Chief Accountant |
| Revaluation | Chief Accountant | Admin/CFO |

Approval properties: (a) approval must be digital (not paper) — recorded in system with timestamp; (b) rejected requests must include rejection reason; (c) approved requests auto-execute; (d) approval can be delegated to deputy chief accountant for max 30 consecutive days; (e) approval history viewable from asset transaction timeline. |
| **Rationale** | High-value FA changes materially impact financial statements. Two-step approval provides segregation of duties and prevents unauthorized changes. Above-threshold disposals require additional oversight (CFO/Admin) for enterprise risk management. |
| **Regulatory Reference** | Accounting Law Article 12 (internal control), Enterprise Law, standard Vietnamese enterprise control procedures (MISA, Fast, Bravo observation) |
| **Validation Logic** | ```
IF operation in ["ADJUSTMENT", "DISPOSAL", "SALE"]:
    IF NOT approval.step_1_completed:
        REJECT "Step 1 submission required before approval"
    IF operation in ["ADJUSTMENT_VALUE", "ADJUSTMENT_LIFE", "ADJUSTMENT_METHOD"]:
        IF request.approver_role != "chief_accountant":
            REJECT "Two-step approval: step 2 requires chief accountant"
    ELSE IF operation in ["DISPOSAL", "SALE"]
        AND asset.carrying_amount > 500_000_000:
        IF request.approver_role NOT IN ("chief_accountant", "admin"):
            REJECT "Two-step approval: disposal >500M requires chief accountant + admin"
```
| **Error Message** | `"FA_R30: Two-step approval required for {{operation}}. Step 2 must be completed by {{required_role}}"` |

---

## Appendix A: Rule Index

| ID | Title | Category | Severity |
|----|-------|----------|----------|
| R-01 | FA Recognition Threshold | Accounting | BLOCKING |
| R-02 | Original Cost Composition | Accounting | BLOCKING |
| R-03 | Capitalizable Borrowing Costs | Accounting | BLOCKING |
| R-04 | Capital vs. Expense | Accounting | BLOCKING |
| R-05 | Component Depreciation | Accounting | WARNING |
| R-06 | Depreciation Method Consistency | Accounting | WARNING |
| R-07 | Residual Value Review | Accounting | WARNING |
| R-08 | Impairment Indicators | Accounting | WARNING |
| R-09 | Straight-Line Calculation | Depreciation | BLOCKING |
| R-10 | Declining Balance | Depreciation | BLOCKING |
| R-11 | Production-Based Method | Depreciation | BLOCKING |
| R-12 | First Period Prorata | Depreciation | BLOCKING |
| R-13 | Last Period Adjustment | Depreciation | BLOCKING |
| R-14 | Suspension Rules | Depreciation | BLOCKING |
| R-15 | Asset Code Uniqueness | Validation | BLOCKING |
| R-16 | Category Level Depth ≤ 3 | Validation | BLOCKING |
| R-17 | Disposal State Constraint | Validation | BLOCKING |
| R-18 | Deletion Constraint | Validation | BLOCKING |
| R-19 | Allocation Sum = 100% | Validation | BLOCKING |
| R-20 | Adjustment Approval | Validation | BLOCKING |
| R-21 | FA Register Form | Compliance | BLOCKING |
| R-22 | Depreciation Schedule | Compliance | BLOCKING |
| R-23 | FA Movement Report | Compliance | BLOCKING |
| R-24 | E-Invoice for FA Sale | Compliance | BLOCKING |
| R-25 | Purchase Invoice VAT | Compliance | WARNING |
| R-26 | CIT Adjustment | Compliance | WARNING |
| R-27 | Role-Based Access | Security | BLOCKING |
| R-28 | Audit Trail | Security | BLOCKING |
| R-29 | Period Lock | Security | BLOCKING |
| R-30 | Two-Step Approval | Security | BLOCKING |

Severity legend:
- **BLOCKING** → operation prevented if rule violated
- **WARNING** → operation proceeds with flag/report, violation logged

---

## Appendix B: Regulatory Reference Map

| Rule | VAS 03 | VAS 04 | TT 45/147 | TT 99 | Decree 123 | IAS 16/36 |
|------|--------|--------|-----------|-------|------------|-----------|
| R-01 | ✓ | ✓ | ✓ | ✓ | | IAS 16 |
| R-02 | ✓ | | ✓ | ✓ | | IAS 16 |
| R-03 | | | | ✓ | | IAS 23 |
| R-04 | ✓ | | ✓ | ✓ | | IAS 16 |
| R-05 | ✓ | | | | | IAS 16 |
| R-06 | ✓ | | ✓ | | | IAS 16 |
| R-07 | ✓ | ✓ | | | | IAS 16 |
| R-08 | | | | ✓ | | IAS 36 |
| R-09 | ✓ | | ✓ | | | IAS 16 |
| R-10 | ✓ | | ✓ | | | IAS 16 |
| R-11 | ✓ | | ✓ | | | IAS 16 |
| R-12 | ✓ | | ✓ | | | IAS 16 |
| R-13 | ✓ | | | | | IAS 16 |
| R-14 | ✓ | | ✓ | | | |
| R-15 | ✓ | | | | | |
| R-16 | ✓ | | ✓ | | | |
| R-17 | ✓ | | ✓ | | | |
| R-18 | ✓ | | | | | |
| R-19 | | | | ✓ | | |
| R-20 | ✓ | | ✓ | | | IAS 16 |
| R-21 | ✓ | | | | | |
| R-22 | ✓ | | | | | |
| R-23 | ✓ | | | | | |
| R-24 | | | | | ✓ | |
| R-25 | | | | ✓ | ✓ | |
| R-26 | | | | ✓ | | IAS 12 |
| R-27 | | | | | | |
| R-28 | ✓ | | | | | |
| R-29 | | | | ✓ | | |
| R-30 | | | | | | |

---

*Document generated for GoTax Fixed Asset Module implementation. All rules reviewed against Vietnamese Accounting Law, VAS 03/04, TT 45/2013 (→ TT 147/2024), TT 99/2025, Decree 123/2020, and IFRS convergence standards.*
