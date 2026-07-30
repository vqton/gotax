# Fixed Asset Module — Readiness Assessment

**Version:** 1.0  
**Date:** 2026-07-30  
**Status:** Draft  

---

## Executive Verdict

**FA module is ~0% complete. NOT production-ready.**

| Dimension | Score | Detail |
|-----------|-------|--------|
| Domain models | 0% | No `FixedAsset` struct. Only `DetailFixedAsset` enum + `OriginalCost`/`AccDepreciation` in OpeningBalanceDetail |
| Database schema | 0% | COA accounts (211-214, 241, 1332) exist in migration. No `fixed_assets` table |
| Repository | 0% | No `FixedAssetRepository` interface, no PG or memory impl |
| Service | 0% | No `FAService`. No depreciation engine |
| Handler | 0% | No FA endpoints |
| Routes | 0% | No FA routes registered |
| Depreciation engine | 0% | No calculation logic |
| Reports | 0% | No FA reports |
| Tax integration | 5% | VAT detection on 211-prefix only. Missing 212, 213, 241 |
| E-invoice | 0% | No FA e-invoice pipeline |
| **Overall** | **~0%** | Foundation only — COA + opening balance FA type |

---

## Gap Analysis vs Reference Implementations

### Feature Comparison Matrix

| Feature | MISA AMIS | Fast | Bravo | Tryton | Odoo | **GoTax (current)** | **GoTax (target)** |
|---------|-----------|------|-------|--------|------|---------------------|-------------------|
| FA Registration | Full | Full | Full | Full | Full | None | P1 |
| FA Category (3-level) | Yes | Yes | Yes | No | Yes | None | P1 |
| Straight-line Dep | Yes | Yes | Yes | Yes | Yes | None | P0 |
| Declining Balance Dep | Yes | Yes | No | No | Yes | None | P0 |
| Production-based Dep | Yes | Yes | No | No | No | None | P1 |
| Prorata Temporis | Yes | Yes | Partial | No | Yes | None | P1 |
| Multi-dept Allocation | Yes | Yes | Partial | No | Analytic | None | P1 |
| FA Transfer | Yes | Yes | Yes | No | Yes | None | P1 |
| FA Adjustment | Yes | Yes | Yes | Limited | Yes | None | P1 |
| FA Suspension | Yes | Yes | No | No | Yes | None | P2 |
| FA Disposal (sale) | Yes | Yes | Yes | Yes | Yes | None | P0 |
| FA Disposal (scrap) | Yes | Yes | Yes | Yes | Yes | None | P1 |
| Revaluation | Yes | Yes | No | Limited | Yes | None | P2 |
| Impairment | Yes | Yes | No | No | Yes | None | P2 |
| FA Physical Inventory | Yes | Yes | Yes | No | Limited | None | P2 |
| FA Maintenance | Separate | No | **Yes** | No | Separate | None | P3 |
| FA Register Report | Yes | Yes | Yes | Yes | Yes | None | P0 |
| Depreciation Schedule | Yes | Yes | Yes | Yes | Yes | None | P0 |
| FA Movement Report | Yes | Yes | Yes | No | Yes | None | P1 |
| GL Auto-posting | Yes | Yes | Yes | Yes | Yes | None | P0 |
| E-invoice Integration | Yes | Yes | Yes | No | Via locale | None | P2 |
| Import from Excel | Yes | Yes | Yes | No | Yes | None | P1 |
| Component Depreciation | No | No | No | No | Limited | None | P3 |
| Tax vs Accounting Book | No | No | No | No | Limited | None | P3 |

### Critical Gaps

1. **No domain entity** — `FixedAsset` struct must be created in `internal/domain/models_fa.go`
2. **No database table** — `010_fixed_asset_schema.up.sql` migration with all 7 tables
3. **No repository** — `FixedAssetRepository` interface + PG + memory impls
4. **No depreciation engine** — pure Go calculator supporting 3 methods
5. **No FA lifecycle** — state machine with 7 states and 10+ transitions
6. **No GL integration** — auto-posting to journal entries for all FA transactions
7. **No FA reports** — register, depreciation schedule, movement report

---

## Implementation Roadmap

### Phase 1: Foundation (P0) — Est. 2 weeks

```
Week 1:
├── Domain models: FixedAsset, FixedAssetCategory, DepreciationEntry, FATransaction
├── Database migration: 010_fixed_asset_schema.up.sql (7 tables)
├── Repository interfaces
└── PG + memory repository impls

Week 2:
├── FA Service basics: Create, Get, List, Update
├── FA Handler + routes registration
├── FA Category CRUD
└── Tests: models, handler, service (in-memory)
```

**Deliverable:** FA CRUD working in API, FA categories, data persisted in both backends.

### Phase 2: Core Depreciation (P0) — Est. 2 weeks

```
Week 3:
├── Depreciation engine: Straight-line (+ prorata)
├── Depreciation engine: Declining balance
├── Period validation (no duplicate, no closed period)
└── Preview depreciation entries

Week 4:
├── Post depreciation to GL (journal entry creation)
├── Unpost depreciation (rollback)
├── Depreciation schedule API
└── Tests: all methods, GL posting, edge cases
```

**Deliverable:** Monthly depreciation run working, GL auto-posting, rollback supported.

### Phase 3: FA Lifecycle (P0-P1) — Est. 2 weeks

```
Week 5:
├── State machine: DRAFT → ACTIVE → DEPRECIATING → FULLY_DEPR
├── FA activation (start depreciation date)
├── FA disposal (scrap/liquidate — zero proceeds)
├── FA sale (with proceeds, gain/loss calc)
└── GL posting for disposal/sale

Week 6:
├── FA transfer between departments
├── FA allocation (multi-department %)
├── FA adjustment (value, useful life, method, residual)
├── FA suspension/resume
└── Tests: all lifecycle transitions
```

**Deliverable:** Full FA lifecycle with state machine, all transitions GL-posted.

### Phase 4: Reports & Expansion (P1-P2) — Est. 2 weeks

```
Week 7:
├── FA Register (S02-TSCĐ per TT 200/2014)
├── Depreciation Schedule (grouped by dept/category)
├── FA Movement Report (VAS 03 disclosure format)
├── FA Aging Report
└── PDF/Excel export

Week 8:
├── FA Inventory (plan → count → result → adjust)
├── FA opening balance import (for legacy migration)
├── E-invoice integration (FA sale XML)
├── CIP-to-FA transfer
└── Tax integration (VAT on 212, 213, 241)
```

**Deliverable:** All 6 FA reports, inventory workflow, basic e-invoice linkage.

### Phase 5: Advanced (P2-P3) — Future

```
├── FA revaluation (IFRS revaluation model)
├── FA impairment (IAS 36 triggers + reversal)
├── Component depreciation
├── Separate tax book (CIT adjustments)
├── FA maintenance tracking
├── Barcode/RFID integration for inventory
├── FA analytics dashboard
└── Multi-currency FA
```

**Deliverable:** IFRS-compliant FA, Pillar 2 readiness, enterprise FA analytics.

---

## Build Sequence

Following GoTax convention:

```
1. Domain models              → internal/domain/models_fixed_asset.go
2. Database migration         → migrations/0010_fixed_asset_schema.up.sql + .down.sql
3. Repository interface       → internal/domain/interfaces.go (add to existing)
4. PG repository             → internal/repository/pg_fixed_asset.go
5. Memory repository         → internal/repository/memory_fixed_asset.go
6. Depreciation engine       → internal/service/depreciation_engine.go
7. Service                   → internal/service/fixed_asset.go
8. Handler                   → internal/handler/fixed_asset.go
9. Routes                    → handler/handler.go (RegisterRoutesWithFA)
10. Wire in main.go          → main.go (PG + memory branches)
11. Tests                    → internal/handler/fixed_asset_handler_test.go
                              internal/service/fixed_asset_service_test.go
                              internal/domain/models_fixed_asset_test.go
12. go vet ./... && go test ./...
```

---

## Regulatory Verification Checklist

- [ ] Circular 45/2013/TT-BTC (or TT 147/2024) — recognition, depreciation, useful life
- [ ] VAS 03 — Tangible FA accounting
- [ ] VAS 04 — Intangible FA accounting
- [ ] Circular 99/2025/TT-BTC — COA accounts 211-214, 241, 1332, 6274, 6414, 6424, 811, 711
- [ ] Circular 200/2014/TT-BTC — FA register, FA tag form
- [ ] Decree 123/2020/ND-CP — e-invoice for FA sale
- [ ] IAS 16 — Property, Plant & Equipment
- [ ] IAS 36 — Impairment of Assets
- [ ] IAS 38 — Intangible Assets
- [ ] IAS 23 — Borrowing Costs
- [ ] Law 69/2024/QH15 — Pillar 2 Global Minimum Tax (FA in GloBE)

---

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| TT 147/2024 changes useful life table | Medium | High | Monitor, make useful life configurable per category |
| Depreciation calc performance (10k+ FA) | Medium | Low | Batch processing, async calculation |
| Multi-company FA isolation bug | High | Low | CompanyID on all FA tables, test across tenants |
| Incorrect GL posting amount | High | Medium | Preview before post, audit trail, unpost ability |
| FA opening balance mismatch with GL | High | Medium | Validation check: sum FA balances = GL 211 balance |
| E-invoice spec changes | Low | Medium | Separate FA XML mapper, test with GDT sandbox |
