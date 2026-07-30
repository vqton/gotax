# Fixed Asset Module — Business Requirements Document (BRD)

**Version:** 1.0  
**Date:** 2026-07-30  
**Author:** BA Lead + Chief Accountant (20+ yrs)  
**Status:** Draft  

---

## 1. Executive Summary

GoTax lacks a Fixed Asset (FA) module. Current FA support is limited to:
- COA accounts for FA (211-214, 1332, 2411, 2414) in SQL migration
- `DetailFixedAsset` entity type in Opening Balance module
- `OriginalCost` / `AccDepreciation` fields on `OpeningBalanceDetail`
- VAT input detection on 211-prefix accounts in tax calculation

No FA lifecycle, no depreciation engine, no disposal, no FA reports. **FA module is ~0% complete.**

This document specifies a complete FA module compliant with:
- Circular 45/2013/TT-BTC (FA management & depreciation) — or its replacement TT 147/2024
- VAS 03 — Tangible Fixed Assets
- VAS 04 — Intangible Fixed Assets
- Circular 99/2025/TT-BTC — Chart of Accounts
- Circular 200/2014/TT-BTC & 133/2016/TT-BTC — Enterprise accounting regime
- Decree 123/2020/ND-CP — E-invoices
- IAS 16 — Property, Plant & Equipment (IFRS convergence)
- IAS 36 — Impairment of Assets
- IAS 38 — Intangible Assets

---

## 2. Business Context

### 2.1 Target Users

| Actor | Role | Needs |
|-------|------|-------|
| Chief Accountant (Kế toán trưởng) | FA policy, depreciation method approval, report sign-off | Compliance, accurate depreciation, audit trail |
| FA Accountant (Kế toán TSCĐ) | Daily FA operations: registration, depreciation, disposal | Efficient processing, schedule automation |
| CFO / Financial Manager | FA investment decisions, budget planning | FA value overview, aging, utilization reports |
| Tax Accountant (Kế toán thuế) | VAT declaration, CIT adjustment for FA | Correct tax treatment per transaction |
| Internal Auditor (Kiểm toán nội bộ) | FA existence verification, impairment review | FA inventory, audit trail |
| External Auditor | FA schedule substantiation, depreciation testing | Full FA register, movement schedule |

### 2.2 Business Problems Solved

1. **Manual depreciation**: Companies calculate depreciation in Excel, error-prone, no audit trail
2. **Tax non-compliance**: Incorrect useful life or method selection leads to CIT adjustments
3. **FA tracking**: No centralized FA register — assets lost, double-counted, or unreported
4. **Partial disposal**: Difficult to handle partial sale/scrap of multi-component assets
5. **Revaluation/impairment**: No workflow for value adjustments per VAS 03
6. **FA opening balance**: Hard to migrate from legacy systems — manual journal entries only
7. **Maintenance integration**: FA repair costs not linked to asset lifecycle
8. **E-invoice linkage**: FA purchase/sale not automatically reflected in e-invoice pipeline

### 2.3 Regulatory Matrix

| Regulation | Key Requirements | FA Module Impact |
|-----------|-----------------|------------------|
| VAS 03 / TT 45 | 4 recognition criteria, cost model, 3 dep methods, useful life table | Core 211/212 asset lifecycle |
| VAS 04 / TT 45 | Intangible FA criteria, development phase, ≤20 yr amortization | 213 asset lifecycle |
| TT 99/2025 | FA accounts (211-214, 241, 1332, 811, 711) | GL posting mapping |
| TT 200/2014 | FA register form, depreciation schedule, FA movement report | Report templates |
| Decree 123/2020 | E-invoice for FA sale, purchase invoice for FA acquisition | E-invoice integration |
| IAS 16 | Component depreciation, revaluation model, residual value review | Optional IFRS features |
| IAS 36 | Impairment indicators, recoverable amount, impairment reversal | Impairment workflow |
| Law 69/2024/QH15 (Pillar 2) | Global minimum tax — FA in GloBE income calc | FA register for GMT |

---

## 3. Functional Requirements

### FR-01: FA Registration (Acquisition)

| ID | FR-01 |
|----|-------|
| **Title** | Fixed Asset Registration |
| **Priority** | P0 — Critical |
| **Description** | User registers new FA with full details. Sources: purchase, self-construction, finance lease, donation, capital contribution, exchange |
| **Data captured** | FA code, name, category, acquisition date, original cost, useful life, residual value, depreciation method, department, location, supplier, contract, e-invoice ref |
| **Validation** | Original cost ≥ 30M VND, useful life > 1 year, category required, department required |
| **GL posting** | Debit 211x, Credit 331/111/112; VAT Debit 1332, Credit 331/111/112 |
| **CIP route** | Debit 2411 (accumulation), then transfer to 211x on completion |

### FR-02: FA Classification

| ID | FR-02 |
|----|-------|
| **Title** | FA Category Management |
| **Priority** | P0 — Critical |
| **Description** | Hierarchical FA categories. 3 levels: Group → Type → Category. Each level defines default depreciation parameters |
| **Category data** | Code, name, parent, default useful life, default depreciation method, default GL accounts (asset, depreciation, expense) |
| **Validation** | Circular reference prevention, level depth ≤ 3 |

### FR-03: FA Depreciation Calculation

| ID | FR-03 |
|----|-------|
| **Title** | Periodic Depreciation |
| **Priority** | P0 — Critical |
| **Description** | Auto-calculate depreciation per period for all active FA. Generate GL posting entries |
| **Methods supported** | Straight-line, Declining balance, Production-based (per Circular 45), Straight-line with prorata temporis (first period partial) |
| **Frequency** | Monthly (default), configurable per asset |
| **Prorata** | First period: days remaining in month / total days. Last period: remaining balance |
| **GL posting** | Debit 6274/6414/6424, Credit 214x |
| **Multi-department** | Allocate depreciation % across departments, each posts to own expense account |

### FR-04: FA Transfer / Movement

| ID | FR-04 |
|----|-------|
| **Title** | FA Transfer |
| **Priority** | P1 — High |
| **Description** | Transfer FA between departments, locations, or users. Mid-period transfer recalculates depreciation proportionally |
| **Validation** | Asset must be ACTIVE or DEPRECIATING. Cannot transfer DISPOSED or SOLD assets |
| **GL posting** | No GL change (same 211 account). Only depreciation allocation adjusts |

### FR-05: FA Adjustment

| ID | FR-05 |
|----|-------|
| **Title** | FA Adjustment |
| **Priority** | P1 — High |
| **Description** | Adjust FA value, useful life, depreciation method, or residual value. Recalculate remaining depreciation |
| **Types** | Value increase (revaluation/capitalized improvement), value decrease (impairment/partial disposal), useful life change, method change, residual value change |
| **Validation** | Adjustment requires approval (chief accountant). Cannot adjust DISPOSED assets |
| **GL posting** | Value increase: Debit 211x, Credit 711 (gain) or 411 (revaluation surplus). Value decrease: Debit 811, Credit 211x |

### FR-06: FA Disposal / Derecognition

| ID | FR-06 |
|----|-------|
| **Title** | FA Disposal |
| **Priority** | P0 — Critical |
| **Description** | Dispose FA through sale, liquidation, donation, or other. Calculate gain/loss. Generate GL entries |
| **Types** | Sell (customer invoice + e-invoice), Liquidation (scrap), Donation (charitable contribution), Return to lessor (finance lease) |
| **Calculation** | Gain/Loss = Net proceeds - Carrying amount (Original cost - Accumulated depreciation) |
| **GL posting** | Debit 214x (accumulated dep), Debit 111/112/131 (proceeds), Debit 811 (loss) or Credit 711 (gain); Credit 211x (original cost) |
| **Validation** | Asset must not be fully depreciated and already disposed. One disposal per asset |

### FR-07: FA Suspension / Resume

| ID | FR-07 |
|----|-------|
| **Title** | FA Depreciation Suspension |
| **Priority** | P2 — Medium |
| **Description** | Temporarily suspend depreciation (e.g., asset under repair, idle). Resume with recalculated schedule |
| **Rules** | Suspension max 12 months per TT 45. Recalculate remaining useful life after resume |
| **GL posting** | Stop depreciation during suspension. No GL entry |

### FR-08: FA Inventory / Physical Verification

| ID | FR-08 |
|----|-------|
| **Title** | FA Inventory |
| **Priority** | P2 — Medium |
| **Description** | Periodic physical FA verification. Record discrepancies, adjust records |
| **Process** | Create inventory plan → Print FA tags/list → Physical count → Record results → Adjust discrepancies |
| **Discrepancy handling** | Missing: investigate, adjust record. Found unregistered: add to register. Damaged: impairment/revaluation |

### FR-09: FA Impairment

| ID | FR-09 |
|----|-------|
| **Title** | FA Impairment Assessment |
| **Priority** | P2 — Medium |
| **Description** | Annual impairment assessment per IAS 36 / VAS guidance. Record impairment loss |
| **Indicators** | Market value decline, physical damage, technological obsolescence, change in use, regulatory changes |
| **GL posting** | Debit 811, Credit 2294 (Impairment provision). Impairment reversal: Credit 811, Debit 2294 |

### FR-10: FA Revaluation

| ID | FR-10 |
|----|-------|
| **Title** | FA Revaluation |
| **Priority** | P2 — Medium |
| **Description** | Revalue FA to fair value per VAS 03 / IAS 16 revaluation model. Supported for listed IFRS entities |
| **Frequency** | Regular intervals (annually for volatile assets, 3-5 years for stable) |
| **GL posting** | Increase: Debit 211x, Credit 411 (revaluation surplus). Decrease: Debit 811 (P&L loss), Credit 211x |

### FR-11: FA Reports

| ID | FR-11 |
|----|-------|
| **Title** | FA Reporting |
| **Priority** | P0 — Critical |
| **Required reports** | |
| | FA Register (Sổ TSCĐ) per TT 200/2014 form |
| | Depreciation Schedule (Bảng tính khấu hao) per period |
| | FA Increase/Decrease Report (Báo cáo tăng/giảm TSCĐ) |
| | FA Movement Schedule per VAS 03 disclosure |
| | FA Inventory Report (Báo cáo kiểm kê TSCĐ) |
| | FA by Department/Category summary |
| | FA Aging Report (TSCĐ quá cũ / sắp hết khấu hao) |
| | FA Utilization Report (utilized vs idle) |
| **Export** | PDF, Excel, CSV |

### FR-12: FA E-Invoice Integration

| ID | FR-12 |
|----|-------|
| **Title** | FA E-Invoice Linkage |
| **Priority** | P2 — Medium |
| **Description** | Link FA acquisition to purchase e-invoice. Link FA sale to sales e-invoice. Generate XML per GDT spec |
| **Acquisition** | Purchase invoice → asset auto-creation (optional). Manual linking |
| **Disposal** | Sales invoice → auto-calculate gain/loss, generate GL. Manual linking |
| **XML** | Populate FA-specific fields in e-invoice XML (asset description, serial number, etc.) |

---

## 4. Non-Functional Requirements

| NFR | Requirement |
|-----|------------|
| NFR-01 | FA depreciation calc for 10,000+ assets completes in < 30 seconds |
| NFR-02 | FA register supports 100,000+ assets |
| NFR-03 | All FA transactions have full audit trail (who, when, what, old value, new value) |
| NFR-04 | Depreciation calculation supports undo/rollback (unpost depreciation) |
| NFR-05 | Multi-company, multi-tenant (FA belong to company context) |
| NFR-06 | FA data export to GDT-required format for tax audit |
| NFR-07 | FA module supports Vietnamese language UI + report templates |

---

## 5. COA Integration Map

### 5.1 FA GL Accounts (from TT 99/2025)

| Account | Name | Type | Used For |
|---------|------|------|----------|
| 211 | TSCĐ hữu hình | ASSET | Tangible FA at cost |
| 2111-2118 | Sub-accounts | ASSET | By FA type |
| 212 | TSCĐ thuê tài chính | ASSET | Finance lease FA |
| 213 | TSCĐ vô hình | ASSET | Intangible FA |
| 2131-2135 | Sub-accounts | ASSET | By intangible type |
| 2141 | HM TSCĐ hữu hình | CONTRA_ASSET | Accum. dep — tangible |
| 2142 | HM TSCĐ thuê TC | CONTRA_ASSET | Accum. dep — lease |
| 2143 | HM TSCĐ vô hình | CONTRA_ASSET | Accum. dep — intangible |
| 2294 | Dự phòng giảm giá TSCĐ | CONTRA_ASSET | Impairment provision |
| 2411 | Mua sắm TSCĐ | ASSET | CIP — procurement |
| 2412 | Xây dựng cơ bản | ASSET | CIP — construction |
| 2414 | Nâng cấp cải tạo | ASSET | CIP — improvements |

### 5.2 Depreciation Expense Accounts

| Account | Name | Allocation |
|---------|------|------------|
| 6274 | CP SXC — Khấu hao TSCĐ | Manufacturing overhead |
| 6414 | CP bán hàng — Khấu hao TSCĐ | Selling expense |
| 6424 | CP QLDN — Khấu hao TSCĐ | Admin expense |

### 5.3 Transaction → GL Posting

| Transaction | Debit | Credit | Amount |
|------------|-------|--------|--------|
| FA acquisition (purchase) | 211x | 331/111/112 | Original cost |
| VAT on FA purchase | 1332 | 331/111/112 | VAT amount |
| CIP accumulation | 2411 | 331/111/112 | CIP costs |
| CIP → FA (completion) | 211x | 2411 | Total CIP cost |
| Monthly depreciation | 6274/6414/6424 | 2141/2143 | Depreciation amount |
| Accum. up to disposal | 2141/2143 | 211x/213x | Accumulated dep |
| Disposal — book value | 811 (if loss) | 211x/213x | Carrying amount |
| Disposal — proceeds | 111/112/131 | 711 (if gain) | Net proceeds |
| Impairment loss | 811 | 2294 | Impairment amount |
| Revaluation increase | 211x | 411 | Revaluation surplus |
| Upgrade capitalized | 2414 → 211x | 331 | Improvement cost |

---

## 6. FA Lifecycle State Machine

```
                    ┌─────────────────────────────────────────────┐
                    │                                             │
                    v                                             │
              ┌──────────┐    acquisition    ┌──────────┐        │
              │  DRAFT   │ ────────────────→ │  ACTIVE  │        │
              └──────────┘                   └──────────┘        │
                    │                              │              │
                    │ cancel                       │ start dep    │
                    v                              v              │
              ┌──────────┐                   ┌──────────────┐    │
              │ CANCELLED│                   │ DEPRECIATING │    │
              └──────────┘                   └──────────────┘    │
                                                  │    │          │
                     ┌────────────────────────────┘    │          │
                     │                                 │          │
                     v                                 v          │
              ┌──────────────┐              ┌──────────────┐      │
              │  SUSPENDED   │              │ FULLY_DEPR   │      │
              │ (no dep)     │              │ (book val=0) │      │
              └──────────────┘              └──────────────┘      │
                     │                                 │          │
                     │ resume                          │ disposal │
                     v                                 v          │
              ┌──────────────┐              ┌──────────────┐      │
              │ DEPRECIATING │              │  DISPOSED    │      │
              └──────────────┘              │  / SOLD      │      │
                                            └──────────────┘      │
                                                  │                │
                                                  └────────────────┘
```
Transitions: DRAFT → CANCELLED, DRAFT → ACTIVE → DEPRECIATING → SUSPENDED → DEPRECIATING → FULLY_DEPR → DISPOSED/SOLD

---

## 7. Integration Points

| Module | Integration |
|--------|-------------|
| Purchase | Supplier invoice → FA auto-creation. CIP tracking from PO |
| GL | Auto-posting of depreciation, acquisition, disposal. Account balance updates |
| Tax | VAT input detection for FA acquisition. CIT adjustment for depreciation method differences |
| Cash/Bank | Payment order for FA acquisition, disposal proceeds |
| E-Invoice | FA purchase e-invoice → asset linking. FA sale → disposal e-invoice |
| Opening Balance | FA opening balance import (legacy migration) |
| Company | FA belongs to company context. Multi-company FA register |
| User/Auth | FA operation permissions (view, create, approve, dispose) |

---

## 8. Assumptions & Constraints

- FA module assumes VND as base currency. Multi-currency FA deferred to v2
- Revaluation model (FR-10) is IFRS-only — Vietnamese firms use cost model per VAS 03
- Maintenance workflow (FR-07) is manual tracking — full CMMS integration deferred
- Tax depreciation book is same as accounting book (no separate tax book in v1)
- FA inventory (FR-08) uses barcode/RFID tracking — deferred to v2

---

## 9. Open Questions

1. Full text of TT 147/2024 (replacing TT 45/2013) — need confirmation of updated useful life tables and recognition thresholds
2. CIT treatment for revaluation surplus — current Circular guidance vs IFRS treatment
3. Finance lease FA recognition — IFRS 16 vs VAS 06 convergence status
4. FA module standalone vs. integrated with AMIS Tài sản (if GoTax offers separate FA product)
