# PROD Readiness Assessment — Thủ Kho (Warehouse Keeper) Module

**Version:** 1.0
**Date:** August 2026
**Assessor:** BA Lead + Chief Accountant (20+ years combined)

---

## Executive Summary

**Current Status: NOT PROD-READY**

The Thủ kho module does NOT exist in the GoTax codebase. The warehouse module has comprehensive inventory management (CRUD, GRN, DN, transfers, adjustments, takes) but lacks the **separation of duties layer** that defines the Thủ kho role.

---

## Gap Analysis

### What EXISTS (PROD)

| Component | Status | Lines | Notes |
|-----------|--------|-------|-------|
| Warehouse CRUD | ✅ PROD | ~150 | Create, read, update, delete warehouses |
| Item/Category CRUD | ✅ PROD | ~200 | Full item lifecycle |
| Stock Balance | ✅ PROD | ~100 | Per warehouse-item-period tracking |
| Inventory Transaction | ✅ PROD | ~80 | Audit trail for all movements |
| Stock Transfer | ✅ PROD | ~200 | 5-step workflow with GL posting |
| Stock Adjustment | ✅ PROD | ~180 | 5-step workflow with GL posting |
| Stock Take | ✅ PROD | ~180 | 5-step workflow with GL posting |
| GRN Integration | ✅ PROD | In purchase | Posts to stock on receipt |
| Warehouse Handler | ✅ PROD | 729 | 48 routes |
| Warehouse Service | ✅ PROD | 891 | Full business logic |
| Warehouse Repos | ✅ PROD | 1787 | Both PG + Memory |
| Frontend Pages | ✅ PROD | 9 pages | All warehouse screens |

### What's MISSING (Gaps)

| Component | Priority | Effort | Impact |
|-----------|----------|--------|--------|
| **Keeper Assignment** | MUST | LOW | No way to assign keeper to warehouse |
| **Stock Ledger (Sổ kho)** | MUST | MEDIUM | Keeper cannot record slips |
| **Slip Recording Workflow** | MUST | MEDIUM | No Ghi sổ / Bỏ ghi process |
| **Reconciliation Report** | MUST | LOW | No keeper vs accounting cross-check |
| **Keeper Stock Card** | MUST | LOW | No Thẻ kho per keeper |
| **Physical Count (Keeper)** | MUST | MEDIUM | Count exists but no keeper role integration |
| **Keeper Reports** | SHOULD | LOW | No keeper-specific reports |
| **Cost Price Visibility** | SHOULD | LOW | No hide/show toggle for keeper |
| **Module Toggle** | MUST | LOW | Cannot enable/disable keeper features |
| **Keeper Frontend** | MUST | MEDIUM | No keeper UI pages |

---

## Regulatory Compliance Check

| Requirement | Regulation | Status | Action |
|-------------|-----------|--------|--------|
| Separation of duties | Luật Kế toán 2015 Art. 6 | ❌ FAIL | Implement keeper role |
| Internal control | Circular 99 Art. 3 | ❌ FAIL | Implement reconciliation |
| Voucher management | Circular 99 Art. 10 | ⚠️ PARTIAL | Add keeper recording workflow |
| Inventory accounting | Circular 99 Art. 16 | ✅ PASS | Stock balance + transactions exist |
| Document retention | Circular 99 | ✅ PASS | Audit trail exists |

**Verdict:** 2 regulatory requirements FAIL. System cannot be deployed to production Vietnamese enterprises without Thủ kho module.

---

## Benchmark Comparison

| Feature | MISA SME 2026 | GoTax Current | Gap |
|---------|--------------|---------------|-----|
| Separate Thủ kho module | ✅ Module #17 | ❌ None | CRITICAL |
| Toggle on/off | ✅ System config | ❌ None | HIGH |
| Slip recording (Ghi sổ) | ✅ Full workflow | ❌ None | CRITICAL |
| Stock Ledger (Sổ kho) | ✅ S06-DN format | ❌ None | CRITICAL |
| Stock Card (Thẻ kho) | ✅ A4 format | ❌ None | HIGH |
| Reconciliation report | ✅ Đ đối chiếu | ❌ None | CRITICAL |
| Keeper-specific reports | ✅ 5+ reports | ❌ None | HIGH |
| Cost price visibility | ✅ Configurable | ❌ None | MEDIUM |
| Inventory count with keeper | ✅ Integrated | ⚠️ Partial | HIGH |
| Assignment management | ✅ CRUD | ❌ None | CRITICAL |

**MISA Parity: 0%** — GoTax has zero Thủ kho features.

---

## Implementation Roadmap

### Phase 1: Core (Week 1-2) — MUST HAVE

| Task | Files | Effort |
|------|-------|--------|
| Domain models | `internal/domain/models_warehouse_keeper.go` | 2h |
| Repository interface | `internal/domain/interfaces.go` | 1h |
| PG repository | `internal/repository/pg_warehouse_keeper.go` | 4h |
| Memory repository | `internal/repository/memory_warehouse_keeper.go` | 3h |
| Service | `internal/service/warehouse_keeper_service.go` | 6h |
| Handler | `internal/handler/warehouse_keeper_handler.go` | 6h |
| Routes | `internal/handler/handler.go` | 1h |
| Migration | `migrations/000031_warehouse_keeper.up.sql` | 1h |
| Tests | `internal/handler/warehouse_keeper_handler_test.go` | 4h |
| Wire in main.go | `main.go` | 1h |

### Phase 2: Reports (Week 3) — SHOULD HAVE

| Task | Files | Effort |
|------|-------|--------|
| Stock Ledger report | Service + Handler | 3h |
| Stock Card report | Service + Handler | 3h |
| Reconciliation report | Service + Handler | 4h |
| Keeper inventory summary | Service + Handler | 2h |
| Excel export | Handler | 2h |

### Phase 3: Frontend (Week 4) — MUST HAVE

| Task | Files | Effort |
|------|-------|--------|
| Keeper assignment page | `web/app/warehouse-keeper-assignments.html` | 4h |
| Pending slips page | `web/app/warehouse-keeper-pending.html` | 4h |
| Stock ledger page | `web/app/warehouse-keeper-ledger.html` | 4h |
| Stock card page | `web/app/warehouse-keeper-stock-card.html` | 3h |
| Reconciliation page | `web/app/warehouse-keeper-reconciliation.html` | 3h |
| Count page | `web/app/warehouse-keeper-count.html` | 3h |
| Config page | `web/app/warehouse-keeper-config.html` | 2h |

### Total Estimated Effort: 4-5 weeks (1 senior developer)

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Role conflict with existing RBAC | Medium | High | Design as overlay, not replacement |
| Performance from dual records | Low | Medium | Indexed queries, read-only ledger |
| User adoption resistance | High | Medium | Toggleable, clear docs |
| Circular 99 interpretation changes | Low | High | Monitor MOF, modular design |

---

## Recommendation

**DO NOT DEPLOY TO PROD** until至少 Phase 1 is complete. The separation of duties violation is a regulatory blocker. Vietnamese auditors will flag this as a material weakness.

**Minimum viable product:** Phase 1 + Phase 3 (core functionality + frontend) = 3-4 weeks.

**Full parity with MISA:** All 3 phases = 4-5 weeks.

---

## Files Created

| File | Purpose |
|------|---------|
| `WAREHOUSE_KEEPER_BRD.md` | Business requirements |
| `WAREHOUSE_KEEPER_SPECS.md` | Technical specifications |
| `WAREHOUSE_KEEPER_USE_CASES.md` | Use cases with happy/alternative/exception paths |
| `WAREHOUSE_KEEPER_WORKFLOWS.md` | Visual workflows and state machines |
| `WAREHOUSE_KEEPER_RULES.md` | Business rules |
| `WAREHOUSE_KEEPER_DATA_FLOWS.md` | Data flow diagrams |
| `WAREHOUSE_KEEPER_USER_JOURNEYS.md` | User journey walkthroughs |
| `WAREHOUSE_KEEPER_TEMPLATES.md` | UI templates and report formats |
| `WAREHOUSE_KEEPER_READINESS.md` | This file — PROD readiness assessment |
| `research-thu-kho-module.md` | Research findings from MISA/FAST/Bravo/Tryton |
