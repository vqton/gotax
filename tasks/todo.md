# Task List: GoTax Pending Features

## Phase 0: Fix Immediate Blocker

### Task 0.1: Rename purchase interfaces to prefixed names
- **Description:** Fix Go 1.26 interface naming collision by renaming interfaces to prefixed names
- **Acceptance criteria:**
  - [ ] Interfaces renamed to `CreateSupplier`, `GetPO`, `CreateInvoice`, `GetCostAllocation`
  - [ ] No naming collision errors
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** None
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** S (1-2 files)

### Task 0.2: Align memory_purchase.go to new interface names
- **Description:** Update memory repository to implement renamed interfaces
- **Acceptance criteria:**
  - [ ] Memory repo implements all renamed interfaces
  - [ ] All methods have correct signatures
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 0.1
- **Files likely touched:**
  - `internal/repository/memory_purchase.go`
- **Estimated scope:** S (1-2 files)

### Task 0.3: Align pg_purchase.go to new interface names
- **Description:** Update PostgreSQL repository to implement renamed interfaces
- **Acceptance criteria:**
  - [ ] PG repo implements all renamed interfaces
  - [ ] All methods have correct signatures
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 0.1
- **Files likely touched:**
  - `internal/repository/pg_purchase.go`
- **Estimated scope:** S (1-2 files)

### Task 0.4: Align purchase_service.go call sites
- **Description:** Update service to use renamed interface methods
- **Acceptance criteria:**
  - [ ] Service calls correct interface methods
  - [ ] No compilation errors
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 0.1, 0.2, 0.3
- **Files likely touched:**
  - `internal/service/purchase_service.go`
- **Estimated scope:** S (1-2 files)

### Task 0.5: Verify `go build ./...` passes
- **Description:** Confirm build is green after interface fixes
- **Acceptance criteria:**
  - [ ] `go build ./...` completes without errors
  - [ ] `go vet ./...` passes
- **Verification:**
  - [ ] Build commands succeed
- **Dependencies:** Task 0.1, 0.2, 0.3, 0.4
- **Files likely touched:** None
- **Estimated scope:** XS

### Task 0.6: Write purchase handler tests (TDD)
- **Description:** Write handler tests for purchase module using TDD approach
- **Acceptance criteria:**
  - [ ] Tests for all purchase endpoints
- [ ] Tests use in-memory repos
- [ ] Tests follow existing patterns
- **Verification:**
  - [ ] `go test -v -run TestPurchase ./internal/handler/` passes
- **Dependencies:** Task 0.5
- **Files likely touched:**
  - `internal/handler/purchase_handler_test.go`
- **Estimated scope:** M (3-5 files)

### Task 0.7: Run full test suite
- **Description:** Run complete test suite to ensure no regressions
- **Acceptance criteria:**
  - [ ] `go test -count=1 ./...` passes
  - [ ] No test failures
- **Verification:**
  - [ ] Test command succeeds
- **Dependencies:** Task 0.6
- **Files likely touched:** None
- **Estimated scope:** XS

## Checkpoint: Phase 0
- [ ] Build green
- [ ] Purchase module tests pass
- [ ] Code review

---

## Phase 1: Foundation Modules

### Module 1: Number Format

### Task 1.1: Create `internal/format/format.go` with `FormatNumber`
- **Description:** Create number formatting utility
- **Acceptance criteria:**
  - [ ] `FormatNumber` function implemented
  - [ ] Handles Vietnamese number format (dots for thousands)
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Phase 0 complete
- **Files likely touched:**
  - `internal/format/format.go`
- **Estimated scope:** S

### Task 1.2: Implement `ParseNumber`
- **Description:** Implement number parsing from Vietnamese format
- **Acceptance criteria:**
  - [ ] `ParseNumber` function implemented
  - [ ] Handles Vietnamese number format
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 1.1
- **Files likely touched:**
  - `internal/format/format.go`
- **Estimated scope:** S

### Task 1.3: Write unit tests
- **Description:** Write comprehensive tests for number formatting
- **Acceptance criteria:**
  - [ ] Tests for `FormatNumber`
  - [ ] Tests for `ParseNumber`
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -v ./internal/format/` passes
- **Dependencies:** Task 1.1, 1.2
- **Files likely touched:**
  - `internal/format/format_test.go`
- **Estimated scope:** S

### Task 1.4: Integrate with existing handlers
- **Description:** Integrate number formatting with existing handlers
- **Acceptance criteria:**
  - [ ] Handlers use new formatting functions
  - [ ] No breaking changes
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 1.3
- **Files likely touched:**
  - Various handler files
- **Estimated scope:** M

---

### Module 2: System Options

### Task 2.1: Create migration `000035_system_options.up.sql` + `.down.sql`
- **Description:** Create database migration for system options
- **Acceptance criteria:**
  - [ ] Migration files created
  - [ ] Uses `CREATE TABLE IF NOT EXISTS`
  - [ ] Includes indexes
- **Verification:**
  - [ ] Migration files syntactically correct
- **Dependencies:** Task 1.4
- **Files likely touched:**
  - `migrations/000035_system_options.up.sql`
  - `migrations/000035_system_options.down.sql`
- **Estimated scope:** S

### Task 2.2: Create `models_system.go` with `SystemOption` struct
- **Description:** Create domain model for system options
- **Acceptance criteria:**
  - [ ] `SystemOption` struct defined
  - [ ] Validation tags added
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.1
- **Files likely touched:**
  - `internal/domain/models_system.go`
- **Estimated scope:** S

### Task 2.3: Create `SystemOptionRepo` interface in `interfaces.go`
- **Description:** Create repository interface for system options
- **Acceptance criteria:**
  - [ ] Interface defined with CRUD methods
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.2
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** XS

### Task 2.4: Define error variables in `errors.go`
- **Description:** Define error variables for system options
- **Acceptance criteria:**
  - [ ] Error variables defined
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.2
- **Files likely touched:**
  - `internal/domain/errors.go`
- **Estimated scope:** XS

### Task 2.5: Implement `pg_system_option.go`
- **Description:** Implement PostgreSQL repository for system options
- **Acceptance criteria:**
  - [ ] All interface methods implemented
  - [ ] Uses GORM
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.3, 2.4
- **Files likely touched:**
  - `internal/repository/pg_system_option.go`
- **Estimated scope:** M

### Task 2.6: Implement `memory_system_option.go`
- **Description:** Implement in-memory repository for system options
- **Acceptance criteria:**
  - [ ] All interface methods implemented
  - [ ] Uses sync.RWMutex
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.3, 2.4
- **Files likely touched:**
  - `internal/repository/memory_system_option.go`
- **Estimated scope:** M

### Task 2.7: Create `system_option_service.go`
- **Description:** Create service for system options
- **Acceptance criteria:**
  - [ ] Service struct defined
  - [ ] Constructor implemented
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.5, 2.6
- **Files likely touched:**
  - `internal/service/system_option_service.go`
- **Estimated scope:** S

### Task 2.8: Implement Get/Set by category
- **Description:** Implement Get and Set methods for system options
- **Acceptance criteria:**
  - [ ] `GetByCategory` method implemented
  - [ ] `SetByCategory` method implemented
  - [ ] Validation logic included
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.7
- **Files likely touched:**
  - `internal/service/system_option_service.go`
- **Estimated scope:** S

### Task 2.9: Implement defaults initialization
- **Description:** Implement default system options initialization
- **Acceptance criteria:**
  - [ ] `InitializeDefaults` method implemented
  - [ ] Default values defined
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.8
- **Files likely touched:**
  - `internal/service/system_option_service.go`
- **Estimated scope:** S

### Task 2.10: Create `system_option_handler.go`
- **Description:** Create HTTP handler for system options
- **Acceptance criteria:**
  - [ ] Handler struct defined
  - [ ] Constructor implemented
  - [ ] Routes registered
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.9
- **Files likely touched:**
  - `internal/handler/system_option_handler.go`
- **Estimated scope:** M

### Task 2.11: Register routes in `handler.go`
- **Description:** Register system option routes in main handler
- **Acceptance criteria:**
  - [ ] Routes added to `RegisterRoutesWithCompany`
  - [ ] Handler parameter added
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.10
- **Files likely touched:**
  - `internal/handler/handler.go`
- **Estimated scope:** XS

### Task 2.12: Wire repos/services in `main.go` (PG + memory)
- **Description:** Wire system option dependencies in main.go
- **Acceptance criteria:**
  - [ ] PG branch wired
  - [ ] Memory branch wired
  - [ ] Service initialized
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 2.11
- **Files likely touched:**
  - `main.go`
- **Estimated scope:** S

### Task 2.13: Write service tests
- **Description:** Write tests for system option service
- **Acceptance criteria:**
  - [ ] Tests for Get/Set methods
  - [ ] Tests for defaults initialization
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -v ./internal/service/` passes
- **Dependencies:** Task 2.12
- **Files likely touched:**
  - `internal/service/system_option_service_test.go`
- **Estimated scope:** M

### Task 2.14: Write handler tests
- **Description:** Write tests for system option handler
- **Acceptance criteria:**
  - [ ] Tests for all endpoints
  - [ ] Tests use in-memory repos
  - [ ] Tests follow existing patterns
- **Verification:**
  - [ ] `go test -v ./internal/handler/` passes
- **Dependencies:** Task 2.13
- **Files likely touched:**
  - `internal/handler/system_option_handler_test.go`
- **Estimated scope:** M

---

### Module 3: Voucher Numbering

### Task 3.1: Create migration `000036_numbering_rules.up.sql` + `.down.sql`
- **Description:** Create database migration for numbering rules
- **Acceptance criteria:**
  - [ ] Migration files created
  - [ ] Uses `CREATE TABLE IF NOT EXISTS`
  - [ ] Includes indexes
- **Verification:**
  - [ ] Migration files syntactically correct
- **Dependencies:** Task 2.14
- **Files likely touched:**
  - `migrations/000036_numbering_rules.up.sql`
  - `migrations/000036_numbering_rules.down.sql`
- **Estimated scope:** S

### Task 3.2: Create `models_numbering.go`
- **Description:** Create domain model for numbering rules
- **Acceptance criteria:**
  - [ ] `NumberingRule` struct defined
  - [ ] Validation tags added
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.1
- **Files likely touched:**
  - `internal/domain/models_numbering.go`
- **Estimated scope:** S

### Task 3.3: Create `NumberingRuleRepo` interface
- **Description:** Create repository interface for numbering rules
- **Acceptance criteria:**
  - [ ] Interface defined with CRUD methods
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.2
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** XS

### Task 3.4: Define errors
- **Description:** Define error variables for numbering rules
- **Acceptance criteria:**
  - [ ] Error variables defined
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.2
- **Files likely touched:**
  - `internal/domain/errors.go`
- **Estimated scope:** XS

### Task 3.5: Implement PG repo
- **Description:** Implement PostgreSQL repository for numbering rules
- **Acceptance criteria:**
  - [ ] All interface methods implemented
  - [ ] Uses GORM
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.3, 3.4
- **Files likely touched:**
  - `internal/repository/pg_numbering.go`
- **Estimated scope:** M

### Task 3.6: Implement Memory repo
- **Description:** Implement in-memory repository for numbering rules
- **Acceptance criteria:**
  - [ ] All interface methods implemented
  - [ ] Uses sync.RWMutex
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.3, 3.4
- **Files likely touched:**
  - `internal/repository/memory_numbering.go`
- **Estimated scope:** M

### Task 3.7: Create `numbering_rule_service.go`
- **Description:** Create service for numbering rules
- **Acceptance criteria:**
  - [ ] Service struct defined
  - [ ] Constructor implemented
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.5, 3.6
- **Files likely touched:**
  - `internal/service/numbering_rule_service.go`
- **Estimated scope:** S

### Task 3.8: Implement `GetNextNumber` with atomic increment
- **Description:** Implement atomic number generation
- **Acceptance criteria:**
  - [ ] `GetNextNumber` method implemented
  - [ ] Atomic increment logic
  - [ ] Format: PREFIX-YYYYMM-XXXX
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.7
- **Files likely touched:**
  - `internal/service/numbering_rule_service.go`
- **Estimated scope:** M

### Task 3.9: Create handler
- **Description:** Create HTTP handler for numbering rules
- **Acceptance criteria:**
  - [ ] Handler struct defined
  - [ ] Constructor implemented
  - [ ] Routes registered
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.8
- **Files likely touched:**
  - `internal/handler/numbering_handler.go`
- **Estimated scope:** M

### Task 3.10: Wire in main.go
- **Description:** Wire numbering rule dependencies in main.go
- **Acceptance criteria:**
  - [ ] PG branch wired
  - [ ] Memory branch wired
  - [ ] Service initialized
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 3.9
- **Files likely touched:**
  - `main.go`
- **Estimated scope:** S

### Task 3.11: Write tests
- **Description:** Write tests for numbering rule service and handler
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 3.10
- **Files likely touched:**
  - `internal/service/numbering_rule_service_test.go`
  - `internal/handler/numbering_handler_test.go`
- **Estimated scope:** M

## Checkpoint: Phase 1
- [ ] All foundation modules working
- [ ] Tests pass
- [ ] Code review

---

## Phase 2: Core Business Modules

### Module 4: Fiscal Year

### Task 4.1: Extend `periods` table with fiscal_year fields
- **Description:** Add fiscal year fields to periods table
- **Acceptance criteria:**
  - [ ] Migration created
  - [ ] New columns added
  - [ ] Backward compatible
- **Verification:**
  - [ ] Migration applies successfully
- **Dependencies:** Phase 1 complete
- **Files likely touched:**
  - `migrations/000040_fiscal_year.up.sql`
  - `migrations/000040_fiscal_year.down.sql`
- **Estimated scope:** S

### Task 4.2: Add fiscal year config to SystemOption
- **Description:** Add fiscal year configuration to system options
- **Acceptance criteria:**
  - [ ] Default values defined
  - [ ] Validation logic added
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 4.1
- **Files likely touched:**
  - `internal/service/system_option_service.go`
- **Estimated scope:** S

### Task 4.3: Auto-generate 12 periods on FY creation
- **Description:** Implement automatic period generation
- **Acceptance criteria:**
  - [ ] 12 periods generated
  - [ ] Correct date ranges
  - [ ] Status set to OPEN
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 4.2
- **Files likely touched:**
  - `internal/service/period_service.go`
- **Estimated scope:** M

### Task 4.4: Integrate with existing period handler
- **Description:** Integrate fiscal year with existing period endpoints
- **Acceptance criteria:**
  - [ ] New endpoints added
  - [ ] Existing endpoints updated
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 4.3
- **Files likely touched:**
  - `internal/handler/handler.go`
- **Estimated scope:** S

### Task 4.5: Write tests
- **Description:** Write tests for fiscal year functionality
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 4.4
- **Files likely touched:**
  - `internal/service/period_service_test.go`
  - `internal/handler/handler_test.go`
- **Estimated scope:** M

---

### Module 5: Contracts

### Task 5.1: Create migration `000037_contracts.up.sql` + `.down.sql`
- **Description:** Create database migration for contracts
- **Acceptance criteria:**
  - [ ] Migration files created
  - [ ] Uses `CREATE TABLE IF NOT EXISTS`
  - [ ] Includes indexes
- **Verification:**
  - [ ] Migration files syntactically correct
- **Dependencies:** Task 4.5
- **Files likely touched:**
  - `migrations/000037_contracts.up.sql`
  - `migrations/000037_contracts.down.sql`
- **Estimated scope:** S

### Task 5.2: Create `models_contract.go`
- **Description:** Create domain model for contracts
- **Acceptance criteria:**
  - [ ] `Contract` struct defined
  - [ ] Validation tags added
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.1
- **Files likely touched:**
  - `internal/domain/models_contract.go`
- **Estimated scope:** S

### Task 5.3: Create `ContractRepo` interface
- **Description:** Create repository interface for contracts
- **Acceptance criteria:**
  - [ ] Interface defined with CRUD methods
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.2
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** XS

### Task 5.4: Define errors
- **Description:** Define error variables for contracts
- **Acceptance criteria:**
  - [ ] Error variables defined
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.2
- **Files likely touched:**
  - `internal/domain/errors.go`
- **Estimated scope:** XS

### Task 5.5: Implement PG repo
- **Description:** Implement PostgreSQL repository for contracts
- **Acceptance criteria:**
  - [ ] All interface methods implemented
  - [ ] Uses GORM
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.3, 5.4
- **Files likely touched:**
  - `internal/repository/pg_contract.go`
- **Estimated scope:** M

### Task 5.6: Implement Memory repo
- **Description:** Implement in-memory repository for contracts
- **Acceptance criteria:**
  - [ ] All interface methods implemented
  - [ ] Uses sync.RWMutex
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.3, 5.4
- **Files likely touched:**
  - `internal/repository/memory_contract.go`
- **Estimated scope:** M

### Task 5.7: Create `contract_service.go`
- **Description:** Create service for contracts
- **Acceptance criteria:**
  - [ ] Service struct defined
  - [ ] Constructor implemented
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.5, 5.6
- **Files likely touched:**
  - `internal/service/contract_service.go`
- **Estimated scope:** S

### Task 5.8: Implement CRUD + status transitions
- **Description:** Implement contract CRUD and status transitions
- **Acceptance criteria:**
  - [ ] Create, Read, Update, Delete methods
  - [ ] Status transition logic
  - [ ] Validation rules
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.7
- **Files likely touched:**
  - `internal/service/contract_service.go`
- **Estimated scope:** M

### Task 5.9: Create handler
- **Description:** Create HTTP handler for contracts
- **Acceptance criteria:**
  - [ ] Handler struct defined
  - [ ] Constructor implemented
  - [ ] Routes registered
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.8
- **Files likely touched:**
  - `internal/handler/contract_handler.go`
- **Estimated scope:** M

### Task 5.10: Register routes
- **Description:** Register contract routes in main handler
- **Acceptance criteria:**
  - [ ] Routes added to `RegisterRoutesWithCompany`
  - [ ] Handler parameter added
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.9
- **Files likely touched:**
  - `internal/handler/handler.go`
- **Estimated scope:** XS

### Task 5.11: Wire in main.go
- **Description:** Wire contract dependencies in main.go
- **Acceptance criteria:**
  - [ ] PG branch wired
  - [ ] Memory branch wired
  - [ ] Service initialized
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 5.10
- **Files likely touched:**
  - `main.go`
- **Estimated scope:** S

### Task 5.12: Write tests
- **Description:** Write tests for contract service and handler
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 5.11
- **Files likely touched:**
  - `internal/service/contract_service_test.go`
  - `internal/handler/contract_handler_test.go`
- **Estimated scope:** M

---

### Module 6: Loan Agreements

### Task 6.1: Create migration `000038_loan_agreements.up.sql` + `.down.sql`
- **Description:** Create database migration for loan agreements
- **Acceptance criteria:**
  - [ ] Migration files created
  - [ ] Uses `CREATE TABLE IF NOT EXISTS`
  - [ ] Includes indexes
- **Verification:**
  - [ ] Migration files syntactically correct
- **Dependencies:** Task 5.12
- **Files likely touched:**
  - `migrations/000038_loan_agreements.up.sql`
  - `migrations/000038_loan_agreements.down.sql`
- **Estimated scope:** S

### Task 6.2: Create `models_loan.go`
- **Description:** Create domain model for loan agreements
- **Acceptance criteria:**
  - [ ] `LoanAgreement` struct defined
  - [ ] Validation tags added
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.1
- **Files likely touched:**
  - `internal/domain/models_loan.go`
- **Estimated scope:** S

### Task 6.3: Create `LoanAgreementRepo` interface
- **Description:** Create repository interface for loan agreements
- **Acceptance criteria:**
  - [ ] Interface defined with CRUD methods
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.2
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** XS

### Task 6.4: Define errors
- **Description:** Define error variables for loan agreements
- **Acceptance criteria:**
  - [ ] Error variables defined
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.2
- **Files likely touched:**
  - `internal/domain/errors.go`
- **Estimated scope:** XS

### Task 6.5: Implement PG repo
- **Description:** Implement PostgreSQL repository for loan agreements
- **Acceptance criteria:**
  - [ ] All interface methods implemented
  - [ ] Uses GORM
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.3, 6.4
- **Files likely touched:**
  - `internal/repository/pg_loan.go`
- **Estimated scope:** M

### Task 6.6: Implement Memory repo
- **Description:** Implement in-memory repository for loan agreements
- **Acceptance criteria:**
  - [ ] All interface methods implemented
  - [ ] Uses sync.RWMutex
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.3, 6.4
- **Files likely touched:**
  - `internal/repository/memory_loan.go`
- **Estimated scope:** M

### Task 6.7: Create `loan_service.go`
- **Description:** Create service for loan agreements
- **Acceptance criteria:**
  - [ ] Service struct defined
  - [ ] Constructor implemented
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.5, 6.6
- **Files likely touched:**
  - `internal/service/loan_service.go`
- **Estimated scope:** S

### Task 6.8: Implement amortization schedule
- **Description:** Implement loan amortization schedule calculation
- **Acceptance criteria:**
  - [ ] Amortization schedule generated
  - [ ] Correct payment calculations
  - [ ] Interest calculations
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 6.7
- **Files likely touched:**
  - `internal/service/loan_service.go`
- **Estimated scope:** L

### Task 6.9: Implement payment recording
- **Description:** Implement loan payment recording
- **Acceptance criteria:**
  - [ ] Payment recording method
  - [ ] Balance updates
  - [ ] Status transitions
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.8
- **Files likely touched:**
  - `internal/service/loan_service.go`
- **Estimated scope:** M

### Task 6.10: Create handler
- **Description:** Create HTTP handler for loan agreements
- **Acceptance criteria:**
  - [ ] Handler struct defined
  - [ ] Constructor implemented
  - [ ] Routes registered
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.9
- **Files likely touched:**
  - `internal/handler/loan_handler.go`
- **Estimated scope:** M

### Task 6.11: Register routes
- **Description:** Register loan agreement routes in main handler
- **Acceptance criteria:**
  - [ ] Routes added to `RegisterRoutesWithCompany`
  - [ ] Handler parameter added
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.10
- **Files likely touched:**
  - `internal/handler/handler.go`
- **Estimated scope:** XS

### Task 6.12: Wire in main.go
- **Description:** Wire loan agreement dependencies in main.go
- **Acceptance criteria:**
  - [ ] PG branch wired
  - [ ] Memory branch wired
  - [ ] Service initialized
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 6.11
- **Files likely touched:**
  - `main.go`
- **Estimated scope:** S

### Task 6.13: Write tests
- **Description:** Write tests for loan agreement service and handler
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 6.12
- **Files likely touched:**
  - `internal/service/loan_service_test.go`
  - `internal/handler/loan_handler_test.go`
- **Estimated scope:** M

## Checkpoint: Phase 2
- [ ] Core business modules working
- [ ] Tests pass
- [ ] Code review

---

## Phase 3: Integration Modules

### Module 7: E-Banking

### Task 7.1: Design CSV parser interface
- **Description:** Design interface for bank CSV parsers
- **Acceptance criteria:**
  - [ ] Interface defined
  - [ ] Common methods identified
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Phase 2 complete
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** S

### Task 7.2: Implement VCB parser
- **Description:** Implement Vietcombank CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles VCB format
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 7.1
- **Files likely touched:**
  - `internal/einvoice/parser_vcb.go`
- **Estimated scope:** M

### Task 7.3: Implement BIDV parser
- **Description:** Implement BIDV CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles BIDV format
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 7.1
- **Files likely touched:**
  - `internal/einvoice/parser_bidv.go`
- **Estimated scope:** M

### Task 7.4: Implement CTG parser
- **Description:** Implement VietinBank CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles CTG format
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 7.1
- **Files likely touched:**
  - `internal/einvoice/parser_ctg.go`
- **Estimated scope:** M

### Task 7.5: Implement VTB parser
- **Description:** Implement VPBank CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles VTB format
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 7.1
- **Files likely touched:**
  - `internal/einvoice/parser_vtb.go`
- **Estimated scope:** M

### Task 7.6: Implement ACB parser
- **Description:** Implement ACB CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles ACB format
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 7.1
- **Files likely touched:**
  - `internal/einvoice/parser_acb.go`
- **Estimated scope:** M

### Task 7.7: Create `bank_import_service.go`
- **Description:** Create service for bank import
- **Acceptance criteria:**
  - [ ] Service struct defined
  - [ ] Constructor implemented
  - [ ] Parser selection logic
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 7.2-7.6
- **Files likely touched:**
  - `internal/service/bank_import_service.go`
- **Estimated scope:** M

### Task 7.8: Implement reconciliation logic
- **Description:** Implement bank reconciliation logic
- **Acceptance criteria:**
  - [ ] Matching algorithm
  - [ ] Tolerance handling
  - [ ] Status updates
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 7.7
- **Files likely touched:**
  - `internal/service/bank_import_service.go`
- **Estimated scope:** L

### Task 7.9: Create handler
- **Description:** Create HTTP handler for bank import
- **Acceptance criteria:**
  - [ ] Handler struct defined
  - [ ] Constructor implemented
  - [ ] Routes registered
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 7.8
- **Files likely touched:**
  - `internal/handler/bank_import_handler.go`
- **Estimated scope:** M

### Task 7.10: Write tests
- **Description:** Write tests for bank import service and handler
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Parser tests
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 7.9
- **Files likely touched:**
  - `internal/service/bank_import_service_test.go`
  - `internal/handler/bank_import_handler_test.go`
- **Estimated scope:** M

---

### Module 8: E-Tax Filing

### Task 8.1: VAT declaration XML (Mẫu 01)
- **Description:** Implement VAT declaration XML generation
- **Acceptance criteria:**
  - [ ] XML generation implemented
  - [ ] Follows Mẫu 01 format
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 7.10
- **Files likely touched:**
  - `internal/htkk/vat_declaration.go`
- **Estimated scope:** M

### Task 8.2: CIT declaration XML (Mẫu 01)
- **Description:** Implement CIT declaration XML generation
- **Acceptance criteria:**
  - [ ] XML generation implemented
  - [ ] Follows Mẫu 01 format
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 8.1
- **Files likely touched:**
  - `internal/htkk/cit_declaration.go`
- **Estimated scope:** M

### Task 8.3: PIT declaration XML
- **Description:** Implement PIT declaration XML generation
- **Acceptance criteria:**
  - [ ] XML generation implemented
  - [ ] Follows PIT format
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 8.2
- **Files likely touched:**
  - `internal/htkk/pit_declaration.go`
- **Estimated scope:** M

### Task 8.4: Create `tax_filing_service.go`
- **Description:** Create service for tax filing
- **Acceptance criteria:**
  - [ ] Service struct defined
  - [ ] Constructor implemented
  - [ ] XML generation logic
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 8.3
- **Files likely touched:**
  - `internal/service/tax_filing_service.go`
- **Estimated scope:** M

### Task 8.5: Integrate with GDT client
- **Description:** Integrate tax filing with GDT client
- **Acceptance criteria:**
  - [ ] GDT client integration
  - [ ] Submission logic
  - [ ] Status tracking
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 8.4
- **Files likely touched:**
  - `internal/service/tax_filing_service.go`
- **Estimated scope:** M

### Task 8.6: Create handler
- **Description:** Create HTTP handler for tax filing
- **Acceptance criteria:**
  - [ ] Handler struct defined
  - [ ] Constructor implemented
  - [ ] Routes registered
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 8.5
- **Files likely touched:**
  - `internal/handler/tax_filing_handler.go`
- **Estimated scope:** M

### Task 8.7: Write tests
- **Description:** Write tests for tax filing service and handler
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] XML generation tests
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 8.6
- **Files likely touched:**
  - `internal/service/tax_filing_service_test.go`
  - `internal/handler/tax_filing_handler_test.go`
- **Estimated scope:** M

---

### Module 9: Digital Signature

### Task 9.1: Design `SignatureProvider` interface
- **Description:** Design interface for digital signature providers
- **Acceptance criteria:**
  - [ ] Interface defined
  - [ ] Common methods identified
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 8.7
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** S

### Task 9.2: Implement VNPT-CA provider
- **Description:** Implement VNPT-CA signature provider
- **Acceptance criteria:**
  - [ ] Provider implemented
  - [ ] Handles VNPT-CA protocol
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 9.1
- **Files likely touched:**
  - `internal/xmldsig/vnpt_ca.go`
- **Estimated scope:** L

### Task 9.3: Implement mock provider for testing
- **Description:** Implement mock signature provider
- **Acceptance criteria:**
  - [ ] Mock provider implemented
  - [ ] Returns predictable results
  - [ ] Used in tests
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 9.1
- **Files likely touched:**
  - `internal/xmldsig/mock_provider.go`
- **Estimated scope:** S

### Task 9.4: Create `signature_service.go`
- **Description:** Create service for digital signatures
- **Acceptance criteria:**
  - [ ] Service struct defined
  - [ ] Constructor implemented
  - [ ] Provider selection logic
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 9.2, 9.3
- **Files likely touched:**
  - `internal/service/signature_service.go`
- **Estimated scope:** M

### Task 9.5: Write tests
- **Description:** Write tests for signature service
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Provider tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 9.4
- **Files likely touched:**
  - `internal/service/signature_service_test.go`
- **Estimated scope:** M

## Checkpoint: Phase 3
- [ ] Integration modules working
- [ ] Tests pass
- [ ] Code review

---

## Phase 4: Polish Modules

### Module 10: Backup & Restore

### Task 10.1: Create migration `000039_backups.up.sql`
- **Description:** Create database migration for backup tracking
- **Acceptance criteria:**
  - [ ] Migration file created
  - [ ] Uses `CREATE TABLE IF NOT EXISTS`
  - [ ] Includes indexes
- **Verification:**
  - [ ] Migration file syntactically correct
- **Dependencies:** Phase 3 complete
- **Files likely touched:**
  - `migrations/000039_backups.up.sql`
- **Estimated scope:** S

### Task 10.2: Create `backup_service.go`
- **Description:** Create service for backup and restore
- **Acceptance criteria:**
  - [ ] Service struct defined
  - [ ] Constructor implemented
  - [ ] Follows existing patterns
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 10.1
- **Files likely touched:**
  - `internal/service/backup_service.go`
- **Estimated scope:** M

### Task 10.3: Implement pg_dump wrapper
- **Description:** Implement pg_dump wrapper for backups
- **Acceptance criteria:**
  - [ ] pg_dump execution
  - [ ] File handling
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 10.2
- **Files likely touched:**
  - `internal/service/backup_service.go`
- **Estimated scope:** M

### Task 10.4: Implement restore logic
- **Description:** Implement database restore logic
- **Acceptance criteria:**
  - [ ] Restore execution
  - [ ] Validation
  - [ ] Error handling
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 10.3
- **Files likely touched:**
  - `internal/service/backup_service.go`
- **Estimated scope:** M

### Task 10.5: Create handler
- **Description:** Create HTTP handler for backup and restore
- **Acceptance criteria:**
  - [ ] Handler struct defined
  - [ ] Constructor implemented
  - [ ] Routes registered
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 10.4
- **Files likely touched:**
  - `internal/handler/backup_handler.go`
- **Estimated scope:** M

### Task 10.6: Write tests
- **Description:** Write tests for backup service and handler
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 10.5
- **Files likely touched:**
  - `internal/service/backup_service_test.go`
  - `internal/handler/backup_handler_test.go`
- **Estimated scope:** M

---

### Module 11: Report Customization

### Task 11.1: Add report options to SystemOption
- **Description:** Add report customization options to system options
- **Acceptance criteria:**
  - [ ] Report options defined
  - [ ] Default values set
  - [ ] Validation logic
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 10.6
- **Files likely touched:**
  - `internal/service/system_option_service.go`
- **Estimated scope:** S

### Task 11.2: Extend PDF generator
- **Description:** Extend PDF generator for customization
- **Acceptance criteria:**
  - [ ] Customization options
  - [ ] Logo support
  - [ ] Layout options
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 11.1
- **Files likely touched:**
  - `internal/service/report_service.go`
- **Estimated scope:** M

### Task 11.3: Add logo upload
- **Description:** Implement logo upload for reports
- **Acceptance criteria:**
  - [ ] Upload endpoint
  - [ ] File storage
  - [ ] Retrieval logic
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 11.2
- **Files likely touched:**
  - `internal/handler/report_handler.go`
- **Estimated scope:** M

### Task 11.4: Write tests
- **Description:** Write tests for report customization
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 11.3
- **Files likely touched:**
  - `internal/service/report_service_test.go`
  - `internal/handler/report_handler_test.go`
- **Estimated scope:** M

---

### Module 12: Multi-branch

### Task 12.1: Branch-level options support
- **Description:** Add branch-level system options
- **Acceptance criteria:**
  - [ ] Branch options schema
  - [ ] Override logic
  - [ ] Inheritance from company
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 11.4
- **Files likely touched:**
  - `internal/service/system_option_service.go`
- **Estimated scope:** M

### Task 12.2: Per-branch numbering rules
- **Description:** Implement per-branch numbering rules
- **Acceptance criteria:**
  - [ ] Branch-specific rules
  - [ ] Fallback to company rules
  - [ ] Atomic generation
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 12.1
- **Files likely touched:**
  - `internal/service/numbering_rule_service.go`
- **Estimated scope:** M

### Task 12.3: Branch report filtering
- **Description:** Add branch filtering to reports
- **Acceptance criteria:**
  - [ ] Branch parameter support
  - [ ] Filtering logic
  - [ ] Performance optimization
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 12.2
- **Files likely touched:**
  - `internal/service/report_service.go`
- **Estimated scope:** M

### Task 12.4: Write tests
- **Description:** Write tests for multi-branch functionality
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 12.3
- **Files likely touched:**
  - `internal/service/system_option_service_test.go`
  - `internal/service/numbering_rule_service_test.go`
  - `internal/service/report_service_test.go`
- **Estimated scope:** M

## Checkpoint: Phase 4
- [ ] All modules complete
- [ ] Full test suite passing
- [ ] Final code review

---

## Summary

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| Phase 0 | 7 | 1 day |
| Phase 1 | 29 | 5 days |
| Phase 2 | 24 | 8 days |
| Phase 3 | 22 | 7 days |
| Phase 4 | 14 | 4 days |
| **Total** | **96** | **25 days** |

## Verification

After each phase:
1. `go build ./...` passes
2. `go vet ./...` passes
3. `go test -count=1 ./...` passes
4. Code review completed
5. Manual verification of key workflows
