# Cost Accounting Module — Task List

## Phase 1: Foundation

- [ ] Task 1: Domain models — CostObject, CostPool, CostingPeriod, CostingResult
  - Acceptance: All structs defined in `internal/domain/models_cost.go` with JSON tags, validate tags
  - Verify: `go build ./...`
  - Files: `internal/domain/models_cost.go`, `internal/domain/interfaces.go`

- [ ] Task 2: Migration SQL — new tables + CostCenter extension
  - Acceptance: `000032_cost_accounting.up.sql` creates all tables with `IF NOT EXISTS`
  - Verify: Migration runs clean on PG
  - Files: `migrations/000032_cost_accounting.up.sql`, `migrations/000032_cost_accounting.down.sql`

- [ ] Task 3: Repository interfaces in domain/interfaces.go
  - Acceptance: CostObjectRepository, CostPoolRepository, CostingPeriodRepository, CostingResultRepository defined
  - Verify: `go build ./...`
  - Files: `internal/domain/interfaces.go`

- [ ] Task 4: PG + Memory repo implementations
  - Acceptance: All 4 repos implemented with CRUD operations
  - Verify: `go test -count=1 ./internal/repository/...`
  - Files: `internal/repository/pg_cost_accounting.go`, `internal/repository/memory_cost_accounting.go`

- [ ] Task 5: CostObject service (CRUD + validation)
  - Acceptance: Create/Get/List/Update/Delete with company_id scoping
  - Verify: `go test -count=1 ./internal/service/...`
  - Files: `internal/service/cost_object_service.go`

- [ ] Task 6: CostPool service (CRUD + cost collection)
  - Acceptance: Create pools, collect costs from GL accounts, list by period
  - Verify: `go test -count=1 ./internal/service/...`
  - Files: `internal/service/cost_pool_service.go`

- [ ] Task 7: CostObject handler + routes
  - Acceptance: REST endpoints at `/api/v1/cost-objects/`
  - Verify: `go test -count=1 ./internal/handler/...`
  - Files: `internal/handler/cost_object_handler.go`

- [ ] Task 8: CostPool handler + routes
  - Acceptance: REST endpoints at `/api/v1/cost-pools/`
  - Verify: `go test -count=1 ./internal/handler/...`
  - Files: `internal/handler/cost_pool_handler.go`

- [ ] Task 9: Handler tests for CostObject + CostPool
  - Acceptance: 10+ tests per handler covering CRUD + error cases
  - Verify: `go test -count=1 ./internal/handler/ -run TestCost`
  - Files: `internal/handler/cost_object_handler_test.go`, `internal/handler/cost_pool_handler_test.go`

## Phase 2: Costing Engine

- [ ] Task 10: Simple costing method (giản đơn)
  - Acceptance: Calculate unit cost = total costs / total output quantity
  - Verify: Unit test with known inputs produces correct unit cost
  - Files: `internal/service/costing_engine.go`

- [ ] Task 11: Coefficient method (hệ số)
  - Acceptance: Allocate costs using coefficient = actual/plan per cost item
  - Verify: Unit test with coefficient examples from Circular 99
  - Files: `internal/service/costing_engine.go`

- [ ] Task 12: Proportion method (tỷ lệ)
  - Acceptance: Allocate using ratio = actual costs / planned costs
  - Verify: Unit test with proportion examples
  - Files: `internal/service/costing_engine.go`

- [ ] Task 13: Costing period management (open/close)
  - Acceptance: Open period, validate no duplicate, close with cost calculation
  - Verify: `go test -count=1 ./internal/service/...`
  - Files: `internal/service/costing_period_service.go`

- [ ] Task 14: GL journal entry generation from costing results
  - Acceptance: Auto-create journal entries: Dr TK154, Cr TK621/622/627; Dr TK155/632, Cr TK154
  - Verify: Journal entries match Circular 99 templates
  - Files: `internal/service/costing_journal.go`

- [ ] Task 15: Integration hooks — warehouse material issuance → cost pool
  - Acceptance: Material issuance journal entries auto-post to TK621 cost pool
  - Verify: Integration test with warehouse + costing
  - Files: `internal/service/costing_integration.go`

- [ ] Task 16: Integration hooks — payroll labor → cost pool
  - Acceptance: Payroll costs auto-post to TK622 cost pool
  - Verify: Integration test
  - Files: `internal/service/costing_integration.go`

- [ ] Task 17: Integration hooks — FA depreciation → cost pool
  - Acceptance: Depreciation allocates to TK627 cost pool
  - Verify: Integration test
  - Files: `internal/service/costing_integration.go`

- [ ] Task 18: Costing handler + routes
  - Acceptance: POST `/api/v1/costing/calculate`, GET `/api/v1/costing/results`
  - Verify: `go test -count=1 ./internal/handler/...`
  - Files: `internal/handler/costing_handler.go`

- [ ] Task 19: Tests for all 3 methods
  - Acceptance: 5+ tests per method with known inputs/outputs
  - Verify: `go test -count=1 ./internal/service/ -run TestCosting`
  - Files: `internal/service/costing_engine_test.go`

## Phase 3: Advanced Methods

- [ ] Task 20: Standard/norm costing method (định mức)
  - Acceptance: Compare actual vs standard, calculate variance
  - Verify: Unit test
  - Files: `internal/service/costing_engine.go`

- [ ] Task 21: Process costing method (phân bước liên tục)
  - Acceptance: Multi-step cost allocation across production stages
  - Verify: Unit test with 3-step example
  - Files: `internal/service/costing_engine.go`

- [ ] Task 22: WIP valuation entity + logic
  - Acceptance: Evaluate WIP using 4 methods (direct materials, main materials, equivalent units, standard costs)
  - Verify: Unit test
  - Files: `internal/service/wip_valuation.go`

- [ ] Task 23: WIP valuation handler + routes
  - Acceptance: POST `/api/v1/costing/wip-valuation`
  - Verify: `go test -count=1 ./internal/handler/...`
  - Files: `internal/handler/costing_handler.go`

- [ ] Task 24: Tests for advanced methods + WIP
  - Acceptance: 5+ tests for standard, process, WIP methods
  - Verify: `go test -count=1 ./internal/service/ -run TestCosting`
  - Files: `internal/service/costing_engine_test.go`

## Phase 4: Reports + Polish

- [ ] Task 25: By-product exclusion method (loại trừ SP phụ)
  - Acceptance: Deduct by-product value from main product cost
  - Verify: Unit test
  - Files: `internal/service/costing_engine.go`

- [ ] Task 26: Cost calculation sheet report (Thẻ tính giá thành)
  - Acceptance: PDF/JSON report matching Circular 99 template
  - Verify: Report output matches expected format
  - Files: `internal/service/cost_report.go`

- [ ] Task 27: Cost summary by object report
  - Acceptance: Group costs by cost object, show breakdown
  - Verify: Report output
  - Files: `internal/service/cost_report.go`

- [ ] Task 28: WIP valuation report
  - Acceptance: Show WIP balances by cost object
  - Verify: Report output
  - Files: `internal/service/cost_report.go`

- [ ] Task 29: Cost variance analysis report
  - Acceptance: Actual vs standard variance by cost item
  - Verify: Report output
  - Files: `internal/service/cost_report.go`

- [ ] Task 30: Frontend pages (Alpine.js)
  - Acceptance: Cost objects, cost pools, costing run, reports pages
  - Verify: Pages load and display data
  - Files: `web/app/cost-objects.html`, `web/app/cost-pools.html`, `web/app/costing.html`

- [ ] Task 31: Update AGENTS.md + sidebar nav
  - Acceptance: Cost Accounting module documented, sidebar updated
  - Verify: Manual check
  - Files: `AGENTS.md`, `web/static/js/app.js`

- [ ] Task 32: Integration tests
  - Acceptance: End-to-end test: create cost object → collect costs → run costing → verify GL entries
  - Verify: `go test -count=1 ./...`
  - Files: `internal/handler/costing_integration_test.go`
