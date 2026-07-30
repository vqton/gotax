# IFRS Standards Reference — Fixed Assets (Property, Plant & Equipment)

> Compiled from IFRS Foundation official publications. Covers 8 standards relevant to FA module implementation.

---

## 1. IAS 16 — Property, Plant and Equipment

### Scope
Tangible items held for use in production/supply of goods/services, rental, or administrative purposes, expected use >1 period. Includes bearer plants.

### Recognition Criteria
Recognise as asset **if and only if**:
- Probable future economic benefits will flow to entity
- Cost can be measured reliably

No materiality threshold in standard — entity applies its own materiality policy.

### Initial Measurement — Cost

**Cost includes:**
- Purchase price (incl import duties, non-refundable purchase taxes) **after** deducting trade discounts and rebates
- Directly attributable costs to bring asset to location/condition for intended use:
  - Site preparation
  - Initial delivery and handling
  - Installation and assembly
  - Professional fees (architects, engineers)
  - Testing costs (less net proceeds from test production — 2020 amendment prohibits deducting pre-use sale proceeds)
- Estimated costs of dismantling/removal and site restoration (decommissioning provision per IAS 37)

**Excluded from cost:**
- Administrative and general overhead
- Opening/new facility costs
- Training costs
- Abnormal costs (waste, labour disputes, delays)

### Subsequent Measurement — Two Models (policy choice per class)

| Model | Treatment |
|-------|-----------|
| **Cost model** | Carrying amount = cost less accumulated depreciation less accumulated impairment |
| **Revaluation model** | Carrying amount = fair value at revaluation date less subsequent depreciation/impairment. Revaluations must be regular (annual or every 3-5 years depending on volatility). |

**Revaluation accounting:**
- Increase → OCI (revaluation surplus in equity), unless reversing previous decrease in P&L
- Decrease → P&L, unless reversing previous surplus in OCI
- Surplus transferred to retained earnings on derecognition (or as asset is used, via excess of depreciation based on revalued amount over historical cost depreciation)

### Depreciation

**Must depreciate** (land is not depreciated, but buildings are).

**Component approach** (IAS 16.43-49):
- Each **significant part** of an item of PPE with different useful life/benefit pattern must be depreciated separately
- If a component cost is significant relative to total cost, depreciate separately
- Examples: aircraft engines, furnace linings, building roofs, aircraft seating
- On component replacement, derecognise old component (whether or not separately depreciated) and recognise new one

**Depreciation method** must reflect pattern of consumption of future economic benefits:

| Method | Description |
|--------|-------------|
| **Straight-line** | Equal charge over useful life |
| **Declining balance** | Higher charge early, lower later |
| **Units of production** | Based on actual usage/output |
| **Sum-of-years-digits** | Accelerated, based on remaining life fraction |

- **Revenue-based depreciation prohibited** (amended May 2014) — revenue reflects output, not consumption of economic benefits
- Method reviewed at least each financial year-end
- Residual value and useful life reviewed at least each financial year-end

**Depreciation begins** when asset is available for use (not when brought into use).

**Depreciation ceases** at earlier of: derecognition date or classification as held for sale (IFRS 5).

### Derecognition

Derecognise on:
- Disposal (sale, finance lease)
- No future economic benefits expected from use/disposal

**Gain/loss on derecognition:**
- Difference between net disposal proceeds and carrying amount
- Recognised in P&L (not classified as revenue per IFRS 15 — May 2014 amendment)
- Under revaluation model, transfer remaining surplus to retained earnings (not through P&L)

**Exchange of assets:**
- Commercial substance test — if transaction has commercial substance, measure at fair value
- If no commercial substance or fair value not measurable, use carrying amount of asset given up

### Disclosures

Per class of PPE:
- Measurement basis (cost or revaluation)
- Depreciation method used
- Useful lives or depreciation rates
- Gross carrying amount and accumulated depreciation at beginning/end
- Reconciliation of carrying amount (additions, disposals, acquisitions, revaluations, impairment, depreciation, net FX differences)
- Restrictions on title / assets pledged as security
- Contractual commitments to acquire PPE
- Compensation from third parties for impaired/lost items (recognised in P&L when receivable)
- Revaluation: effective date, whether independent valuer, carrying amount if cost model, surplus movement
- Temporarily idle PPE, fully depreciated but still in use, retirement from active use

---

## 2. IAS 36 — Impairment of Assets

### Scope
Applies to PPE, intangible assets, goodwill, right-of-use assets, investment property at cost. Exceptions: inventories, deferred tax, financial assets (IFRS 9), investment property at FV, biological assets, insurance contract assets.

### Core Principle
Asset must not be carried above its **recoverable amount** (higher of fair value less costs to sell and value in use).

### When to Test

| Asset type | Test required |
|------------|---------------|
| Goodwill, indefinite-life intangibles, not-yet-available intangibles | **Annual** (plus on indication) |
| All other assets (PPE, finite-life intangibles, right-of-use) | Only when **impairment indicators** exist |

### Impairment Indicators (External)

**External:**
- Market value decline > expected
- Significant adverse technological, market, economic, or legal changes
- Interest rate increases affecting discount rate
- Entity net asset > market capitalisation

**Internal:**
- Obsolescence or physical damage
- Asset becomes idle, held for disposal, or restructuring
- Worse-than-expected economic performance (actual cash flows < budget)

### Recoverable Amount

`Recoverable Amount = max(Fair Value Less Costs to Sell, Value in Use)`

**Fair value less costs to sell (FVLCS):**
- Arm's length sale between knowledgeable willing parties
- Less incremental costs directly attributable to disposal (legal, stamp duty, removal)

**Value in use (VIU):**
- Present value of expected future cash flows from:
  - Continuing use of asset
  - Disposal at end of useful life
- Cash flow projections based on most recent budgets (max 5 years, extrapolate with stable/declining growth rate)
- Discount rate = pre-tax, reflects time value and asset-specific risks

### Cash-Generating Units (CGUs)
When individual asset VIU cannot be determined, group into smallest identifiable group of assets generating independent cash inflows (CGU). Goodwill allocated to CGUs for impairment testing.

### Recognition
Impairment loss = carrying amount − recoverable amount
- **Individual asset:** Reduce carrying amount. Loss recognised immediately in P&L (or OCI if revalued per IAS 16/38).
- **CGU:** First reduce goodwill allocated to CGU, then pro-rata to other assets (but no asset below zero).

### Depreciation Revisions
After impairment, adjust depreciation to allocate revised carrying amount over remaining useful life.

### Reversal

| Asset | Reversal allowed? |
|-------|-------------------|
| Goodwill | **Never** |
| Other assets | Yes, if circumstances causing impairment have resolved |
| Individual asset | Increase carrying amount, but not above what it would have been without impairment |
| CGU | Increase assets pro-rata (excluding goodwill) |
| Recognition | Immediately in P&L (or OCI if revalued) |

### Disclosures
- Impairment losses and reversals by class/segment (P&L line item)
- Events/circumstances leading to impairment
- Details for individual material impairments: recoverable amount basis (FVLCS vs VIU)
- CGU details for goodwill: carrying amount, discount rate, growth rate, sensitivity

---

## 3. IAS 38 — Intangible Assets

### Scope
Identifiable non-monetary asset without physical substance. Identifiable = separable OR arising from contractual/other legal rights.

Examples: software, licences, trademarks, patents, films, copyrights, import quotas. Excludes goodwill (IFRS 3), financial assets, insurance contract assets, exploration/evaluation assets.

### Recognition
Recognise **if and only if**:
- Probable future economic benefits
- Cost measured reliably
- Meets definition (identifiable, control, no physical substance)

**Internally generated:**
- **Research phase** — all expense (must recognise as expense, cannot capitalise later)
- **Development phase** — capitalise if ALL criteria met:
  - Technical feasibility of completion
  - Intention to complete and use/sell
  - Ability to use/sell
  - Probable future economic benefits (demonstrate market exists)
  - Adequate technical/financial resources to complete
  - Expenditure reliably measurable
- **Never capitalise:** internally generated brands, mastheads, publishing titles, customer lists, goodwill

### Initial Measurement — Cost

**Separate acquisition:**
- Purchase price (import duties, non-refundable taxes, after trade discounts)
- Directly attributable costs (employee benefits, professional fees, testing)
- **Not included:** new admin costs, training, inefficiencies, initial operating losses

**Business combination:** Fair value at acquisition date

**Government grant:** Fair value (or nominal value + directly attributable costs per IAS 20)

**Exchange:** Fair value (unless no commercial substance or fair value unreliable)

### Subsequent Measurement

| Model | Treatment |
|-------|-----------|
| **Cost model** | Cost less accumulated amortisation less impairment (default, overwhelmingly most common) |
| **Revaluation model** | Rare — only if active market exists (uncommon for most intangibles). Fair value less amortisation less impairment. Surplus in OCI. |

### Amortisation

**Finite useful life:** Amortise over useful life. Method reflects consumption pattern. Residual value assumed zero unless:
- Third party committed to purchase
- Active market exists and residual value determinable

**Indefinite useful life:** No amortisation. Annual impairment test. Review annually for finite/indefinite classification.

**Methods allowed:** Straight-line (most common), declining balance, units of production. Revenue-based amortisation **presumed inappropriate** (except limited circumstance where revenue and consumption are highly correlated — 2014 amendment).

**Review:**
- Useful life and amortisation method reviewed annually
- Amortisation period/method changed as change in estimate (IAS 8)

### Derecognition
- Disposal
- No future economic benefits expected
- Gain/loss in P&L

### Disclosures
Per class of intangible:
- Useful lives/amortisation rates, amortisation methods
- Gross carrying amount and accumulated amortisation/impairment
- Reconciliation (additions, disposals, reclassifications, impairment, amortisation, FX)
- Separately: internally generated vs acquired
- Indefinite-life intangibles: carrying amount and rationale
- Research and development expense recognised in P&L
- Revaluation details (if applicable)
- Contractual commitments

---

## 4. IAS 23 — Borrowing Costs

### Core Rule
Borrowing costs directly attributable to acquisition, construction, or production of a **qualifying asset** must be **capitalised** as part of cost. All other borrowing costs expensed.

### Qualifying Asset
Asset that necessarily takes a **substantial period of time** to get ready for intended use or sale. Examples:
- Manufacturing plants
- Power generation facilities
- Investment properties
- Intangible assets (development phase)
- Bearer plants
- Inventories requiring substantial time (e.g., ships, aircraft)

**Does NOT apply to:**
- Assets ready for use at acquisition
- Inventories routinely manufactured over short period

### Capitalisation Mechanics

**Specific borrowings:**
- Capitalise actual borrowing costs less any investment income on temporary investment of those borrowings

**General borrowings (pool):**
- Capitalise: capitalisation rate × expenditures on qualifying asset
- Capitalisation rate = weighted average of borrowing costs on general borrowings outstanding
- Cannot exceed actual borrowing costs incurred

### Capitalisation Period

**Begins when ALL of:**
1. Expenditures for asset being incurred
2. Borrowing costs being incurred
3. Activities necessary to prepare asset for use/sale in progress

**Suspended during:**
- Extended periods where active development is interrupted (not brief/necessary delays)
- Example: suspension while waiting for regulatory approval (if active development stops)

**Ceases when:**
- Substantially all activities to prepare asset for intended use/sale are complete
- Capitalise until asset substantially complete — even if not yet brought into use

### Disclosures
- Amount of borrowing costs capitalised during period
- Capitalisation rate used

---

## 5. IFRS 5 — Non-current Assets Held for Sale

### Classification Criteria
Non-current asset (or disposal group) classified as **held for sale** if:
- Carrying amount will be recovered **principally through sale** rather than continuing use
- Asset must be available for **immediate sale** in present condition
- Sale is **highly probable** within **one year** from classification date
- Management committed to plan, active programme to find buyer initiated
- Asset being actively marketed at reasonable price
- Unlikely plan will be withdrawn (sale expected to meet one-year condition)

Also applies to **held for distribution to owners** (e.g., demerger).

### Measurement
- Measure at **lower of carrying amount and fair value less costs to sell**
- Immediately before classification, measure asset per applicable standard (depreciation, impairment)
- Recognise impairment loss for any write-down to FV less costs
- **Depreciation ceases** on classification as held for sale
- Subsequent reversal of impairment allowed (FV less costs to sell increases)

### Presentation
- Present held-for-sale assets **separately** on balance sheet
- Liabilities of disposal group presented separately (not offset)
- No reclassification or retrospective adjustment
- Discontinued operations presented separately in income statement (single amount: post-tax profit/loss + gain/loss on disposal)

### Disclosures
- Description of asset/disposal group
- Description of facts/circumstances of sale
- Gain/loss recognised and income statement line
- Segment (if applicable)
- If of sale no longer highly probable: change in classification, reason, and effect

---

## 6. IFRS 16 — Leases (Right-of-Use Assets)

### Core Model — Lessee
Single accounting model: recognise **right-of-use (ROU) asset** and **lease liability** for all leases >12 months (unless low-value underlying asset).

### Recognition
At lease commencement date.

### Initial Measurement — ROU Asset

**Cost includes:**
1. Initial amount of lease liability (see below)
2. Lease payments made at/before commencement date **less** lease incentives received
3. Initial direct costs incurred by lessee
4. Estimated dismantling, removal, site restoration costs (per IAS 37)

### Initial Measurement — Lease Liability

Present value of **lease payments not yet paid**:
- Fixed payments (less incentives)
- Variable payments based on index/rate (using initial index/rate)
- Residual value guarantees expected payable
- Purchase option exercise price (if reasonably certain)
- Termination penalty payments (if exercising termination is reasonably certain)

**Discount rate:** Implicit in lease (if determinable) → incremental borrowing rate

### Subsequent Measurement — ROU Asset

**Cost model (default):**
- Cost less accumulated depreciation less accumulated impairment
- Depreciate over **shorter of useful life and lease term**
- Component approach applies per IAS 16 (significant parts of ROU asset depreciated separately)

**Revaluation model** — allowed if entity applies revaluation model to that class of underlying asset per IAS 16.

**Impairment** — IAS 36 applies (test ROU asset for impairment, treat as element of CGU).

### Subsequent Measurement — Lease Liability
- Increase carrying amount to reflect interest (effective interest method per IFRS 9)
- Decrease to reflect lease payments made
- Remeasure when: lease modification, change in lease term, change in purchase option assessment, change in residual value guarantee, change in variable payment index/rate

### Lessor Accounting
- **Finance lease** — derecognise asset, recognise net investment in lease (different from lessee)
- **Operating lease** — lessor continues to recognise asset, depreciation per IAS 16

### Disclosures (Lessee)
- Depreciation charge for ROU assets (by class)
- Interest expense on lease liabilities
- Short-term/low-value lease expense
- Variable lease payments not included in liability
- Income from subleasing ROU assets
- Total cash outflow for leases
- Additions to ROU assets
- Maturity analysis of lease liabilities (per IFRS 7)
- Gains/losses on sale and leaseback

---

## 7. IAS 20 — Government Grants

### Scope
Government grants = transfers of resources in return for past/future compliance with conditions. Includes non-monetary grants at fair value.

### Recognition
Recognise **only when** there is **reasonable assurance** that:
1. Entity will comply with conditions attached
2. Grant will be received

Recognise in P&L on **systematic basis** over periods in which related costs are recognised.

### Grants Related to Assets (Relevant to FA Module)

**Two permitted methods:**

| Method | Presentation |
|--------|-------------|
| **Deferred income** | Grant set up as deferred income (liability) and released to P&L systematically over asset useful life |
| **Net-of-cost (deduction)** | Grant deducted from carrying amount of asset. Depreciation charged on reduced carrying amount |

Both methods result in same P&L effect over asset life.

### Grants Related to Income
- Presented as credit in P&L (separately or under "Other income") — OR
- Deducted from related expense

### Non-monetary Grants
Measure at **fair value** (alternative: nominal value).

### Repayment
Accounted for as change in accounting estimate (IAS 8):
- If deferred income: reduce deferred income balance
- If net-of-cost: increase carrying amount of asset
- Additional depreciation recognised if repayment exceeds remaining deferred income

### Disclosures
- Accounting policy adopted (deferred income vs net-of-cost)
- Nature and extent of grants recognised
- Unfulfilled conditions/other contingencies
- Other forms of government assistance

---

## 8. IFRS for SMEs — Section 17 (Property, Plant and Equipment)

> 2025 edition effective for periods beginning on or after 1 Jan 2027. Summarised below.

### Scope
Same as IAS 16: tangible items for production/supply, rental, administrative, expected use >1 period.

### Simplified from Full IFRS

**Recognition:**
- Same criteria as IAS 16 (probable benefits, reliable cost measurement)
- No materiality threshold defined

**Initial measurement:**
- Same cost components as IAS 16 (purchase price, directly attributable costs, decommissioning)
- Also includes borrowing costs — **option** (not required): SMEs may elect to expense OR capitalise directly attributable borrowing costs per Section 15

**Subsequent measurement — Simplified choice per class:**

| Option | Treatment |
|--------|-----------|
| **Cost model** (default) | Cost less accumulated depreciation less impairment |
| **Revaluation model** | Permitted only if fair value can be measured reliably by reference to **active market** (much more restrictive than IAS 16) |

**Revaluation model limitations:**
- Revaluations must be made with sufficient regularity
- Only when fair value reliably measurable via active market
- Most PPE does not have active market — cost model is effectively the only option for most SMEs
- Increase to OCI, decrease to P&L (same as IAS 16)

### Depreciation
- Component approach **not required** (simplification from IAS 16)
- However, significant parts with different useful lives **should** be depreciated separately if practicable
- Depreciation methods: straight-line, declining balance, units of production (same as IAS 16)
- Revenue-based depreciation **prohibited** (consistent with IAS 16)
- Residual value, useful life, depreciation method reviewed only if indication of change (not annually)

### Borrowing Costs (Section 15)
- Policy choice: expense all OR capitalise qualifying borrowing costs (same definition of qualifying asset per IAS 23)
- If elected, capitalisation method follows full IFRS 23

### Derecognition
- Same as IAS 16: on disposal or when no future benefits expected
- Gain/loss in P&L

### Impairment (Section 20 — different from full IFRS)

Simplified impairment model:
- **Impairment indicator** triggers test (not annual for indefinite-life assets — SMEs unlikely to have them)
- Recoverable amount = higher of **fair value less costs to sell** and **value in use**
- But VIU calculated as undiscounted cash flows (not discounted — significant simplification)
- If impaired, write down to recoverable amount
- Reversal required if conditions improve (except goodwill reversal prohibited)

### Disclosures (significantly fewer than IAS 16)
- Depreciation method and useful lives/rates
- Gross carrying amount and accumulated depreciation at beginning/end
- Reconciliation (additions, disposals, revaluations, impairment, depreciation)
- Revaluation details: effective date, whether independent valuer
- Restrictions on title / assets pledged
- Contractual commitments
- Impairment losses and reversals recognised in P&L

### Key Simplifications Summary vs Full IFRS

| Area | Full IFRS (IAS 16/36) | IFRS for SMEs (Section 17) |
|------|----------------------|---------------------------|
| Component depreciation | Required | Not required (encouraged if significant) |
| Revaluation model | Any asset class with fair value | Only if active market exists |
| Residual value review | Annual | When indication of change |
| Useful life review | Annual | When indication of change |
| Depreciation method review | Annual | When indication of change |
| Value in use calculation | Discounted cash flows | Undiscounted cash flows |
| Annual impairment test | Goodwill/indefinite intangibles | Not required |
| Borrowing costs | Must capitalise (no option) | Policy choice to capitalise or expense |
| Disclosures | Extensive (15+ items) | Moderate (~8 items) |

---

## Cross-Cutting Data Model Requirements for FA Module

### Asset Master Data
| Field | Standard | Notes |
|-------|----------|-------|
| Asset class/type | IAS 16 | Per-class disclosure policy |
| Component flag | IAS 16 | Component approach |
| Cost model vs revaluation | IAS 16 | Policy choice per class |
| Useful life | IAS 16 | Annual review |
| Residual value | IAS 16 | Annual review |
| Depreciation method | IAS 16 | Annual review; revenue method prohibited |
| Location/cost centre | — | Entity-specific |
| Impairment indicators | IAS 36 | External + internal indicators |
| Grant details | IAS 20 | Deferred income vs net-of-cost |
| Lease flag | IFRS 16 | ROU asset vs owned |
| Held-for-sale flag | IFRS 5 | Triggers depreciation stop |

### Depreciation Engine Requirements
- SL, DB, UoP, SYD methods
- Component-level separate depreciation
- Partial-period conventions (mid-month, mid-quarter, full-month)
- Change in estimate (method, life, residual value) — prospective
- Impairment → adjust depreciable base prospectively
- Held-for-sale → stop depreciation
- Revaluation → adjust depreciation from revaluation date

### Impairment Testing Workflow
- Trigger-based (external + internal indicators)
- Annual for goodwill/indefinite-life
- Calculate FVLCS (quoted price, recent transaction, DCF)
- Calculate VIU (discounted cash flows)
- CGU grouping
- Goodwill allocation to CGUs
- Reversal tracking (not for goodwill)

### Disclosure Report Requirements
- Movement schedule (opening → additions → disposals → revaluation → impairment → depreciation → FX → closing)
- Cost vs accumulated depreciation vs net book value
- Per class of asset
- Separately for owned vs ROU vs gran-aquired

---

*Compiled 2026-07 from IFRS Foundation publications. Standards referenced: IAS 16, IAS 36, IAS 38, IAS 23, IFRS 5, IFRS 16, IAS 20, IFRS for SMEs Section 17 (2025 edition). For exact wording consult official IFRS Red Book.*
