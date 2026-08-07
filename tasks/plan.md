# Implementation Plan: Cost Accounting (Giá thành) Module

## Overview
Build the Cost Accounting module for GoTax to support 6 costing methods per Circular 99/2025/TT-BTC. The module collects direct/indirect costs into pools, allocates them to cost objects, calculates unit costs, and values WIP. Currently only Cost Centers CRUD exists.

## Architecture Decisions
- Extend existing CostCenter entity (add `cost_type` field) rather than replacing
- New entities: CostObject, CostPool, CostAllocation, CostingPeriod, CostingResult, WIPValuation
- Cost engine runs as service-layer batch job at period-end
- GL journal entries auto-generated from costing results
- All 6 methods implement same interface; method selection is config per CostObject

## Phase 1: Foundation (Weeks 1-4)
- [ ] Task 1: Domain models — CostObject, CostPool, CostingPeriod, CostingResult
- [ ] Task 2: Migration SQL — new tables + CostCenter extension
- [ ] Task 3: Repository interfaces in domain/interfaces.go
- [ ] Task 4: PG + Memory repo implementations
- [ ] Task 5: CostObject service (CRUD + validation)
- [ ] Task 6: CostPool service (CRUD + cost collection)
- [ ] Task 7: CostObject handler + routes
- [ ] Task 8: CostPool handler + routes
- [ ] Task 9: Handler tests for CostObject + CostPool

### Checkpoint: Phase 1
- [ ] Cost objects and cost pools can be created/updated/listed
- [ ] Tests pass, builds clean

## Phase 2: Costing Engine (Weeks 5-8)
- [ ] Task 10: Simple costing method (giản đơn)
- [ ] Task 11: Coefficient method (hệ số)
- [ ] Task 12: Proportion method (tỷ lệ)
- [ ] Task 13: Costing period management (open/close)
- [ ] Task 14: GL journal entry generation from costing results
- [ ] Task 15: Integration hooks — warehouse material issuance → cost pool
- [ ] Task 16: Integration hooks — payroll labor → cost pool
- [ ] Task 17: Integration hooks — FA depreciation → cost pool
- [ ] Task 18: Costing handler + routes
- [ ] Task 19: Tests for all 3 methods

### Checkpoint: Phase 2
- [ ] Can run simple/ coefficient/ proportion costing for a period
- [ ] GL entries auto-generated
- [ ] Tests pass

## Phase 3: Advanced Methods (Weeks 9-12)
- [ ] Task 20: Standard/norm costing method (định mức)
- [ ] Task 21: Process costing method (phân bước liên tục)
- [ ] Task 22: WIP valuation entity + logic
- [ ] Task 23: WIP valuation handler + routes
- [ ] Task 24: Tests for advanced methods + WIP

### Checkpoint: Phase 3
- [ ] All 5 methods operational
- [ ] WIP valuation works

## Phase 4: Reports + Polish (Weeks 13-16)
- [ ] Task 25: By-product exclusion method (loại trừ SP phụ)
- [ ] Task 26: Cost calculation sheet report (Thẻ tính giá thành)
- [ ] Task 27: Cost summary by object report
- [ ] Task 28: WIP valuation report
- [ ] Task 29: Cost variance analysis report
- [ ] Task 30: Frontend pages (Alpine.js)
- [ ] Task 31: Update AGENTS.md + sidebar nav
- [ ] Task 32: Integration tests

### Checkpoint: Phase 4
- [ ] All 6 methods operational
- [ ] All reports generated
- [ ] Frontend pages functional
- [ ] Full test suite passes

## Risks and Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| Complex allocation logic | High | Start with simple method, add complexity incrementally |
| GL entry correctness | High | Manual verification against MISA examples |
| Performance with large datasets | Medium | Batch processing, database indexes |
| Circular 99 changes | Low | Modular design allows method updates |

## Open Questions
- Should costing run synchronously or as background job?
- How to handle mid-period cost corrections?
- Multi-currency costing support needed?
