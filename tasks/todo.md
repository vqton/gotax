# Payroll Module — Task List

**Generated:** 2026-08-04
**Total Tasks:** 40
**Estimated Duration:** 16 weeks

---

## Phase 1: Foundation & Calculation Engine (Weeks 1-4)

### T1.1: Payroll Migration Schema

**Description:** Create PostgreSQL migration for all payroll-related tables: employee_payroll_info, dependants, salary_components, payroll_periods, payroll_runs, timekeeping_records, leave_requests, leave_balances, payslips, payroll_config.

**Acceptance criteria:**
- [ ] Migration file `000019_payroll_schema.up.sql` created
- [ ] All tables use UUID primary keys with `gen_random_uuid()`
- [ ] All tables use `CREATE TABLE IF NOT EXISTS` for idempotency
- [ ] Foreign keys to employees table properly defined
- [ ] Indexes on company_id, employee_id, period_id
- [ ] Down migration `000019_payroll_schema.down.sql` created
- [ ] `go run .` starts without error on PG backend

**Verification:**
- [ ] `migrate -path migrations -database "$DATABASE_URL" up` succeeds
- [ ] All tables created: `\dt` shows new tables
- [ ] Down migration drops all tables cleanly

**Dependencies:** None
**Files:**
- `migrations/000019_payroll_schema.up.sql`
- `migrations/000019_payroll_schema.down.sql`

**Estimated scope:** M (3-5 files)

---

### T1.2: Payroll Domain Models

**Description:** Create domain model structs for all payroll entities in `internal/domain/models_payroll.go`. Follow existing pattern (all models in `package domain`).

**Acceptance criteria:**
- [ ] `EmployeePayrollInfo` struct with all fields from spec
- [ ] `Dependant` struct
- [ ] `SalaryComponent` struct
- [ ] `PayrollPeriod` struct with status enum
- [ ] `PayrollRun` struct with all calculation fields
- [ ] `Payslip` struct
- [ ] JSON tags on all fields
- [ ] Validate struct tags on required fields

**Verification:**
- [ ] `go build ./...` compiles
- [ ] `go vet ./...` passes

**Dependencies:** T1.1
**Files:**
- `internal/domain/models_payroll.go`

**Estimated scope:** S (1-2 files)

---

### T1.3: Payroll Enums and Constants

**Description:** Define all enum types and constants for payroll module: PayrollStatus, ContractType, SalaryType, LeaveType, etc.

**Acceptance criteria:**
- [ ] `PayrollStatus` enum: DRAFT, PROCESSING, APPROVED, PAID, CLOSED
- [ ] `ContractType` enum: INDEFINITE, DEFINITE, PROBATION
- [ ] `SalaryType` enum: TIME_BASED, PIECE_RATE, COMMISSION, COEFFICIENT
- [ ] `LeaveType` enum: ANNUAL, SICK, MATERNITY, PATERNITY, UNPAID, COMPENSATORY
- [ ] Insurance rate constants (SI: 8%/17.5%, HI: 1.5%/3%, UI: 1%/1%)
- [ ] PIT bracket constants (5 brackets with rates and constants)
- [ ] Regional minimum wage constants

**Verification:**
- [ ] All constants compile
- [ ] No duplicate enum values

**Dependencies:** T1.2
**Files:**
- `internal/domain/models_payroll.go` (same file or separate)

**Estimated scope:** XS (single file addition)

---

### T1.4: Payroll Repository Interfaces

**Description:** Define repository interfaces for all payroll entities in `internal/domain/interfaces.go` (or new `interfaces_payroll.go`).

**Acceptance criteria:**
- [ ] `PayrollInfoRepository` interface: Create, Get, Update, ListByCompany, GetByEmployee
- [ ] `DependantRepository` interface: Create, Update, Delete, ListByEmployee
- [ ] `SalaryComponentRepository` interface: Create, Update, ListByCompany
- [ ] `PayrollPeriodRepository` interface: Create, Update, Get, ListByCompany
- [ ] `PayrollRunRepository` interface: Create, Update, ListByPeriod, Get
- [ ] `TimekeepingRepository` interface: Create, Update, Delete, ListByEmployeePeriod, BulkCreate
- [ ] `LeaveRequestRepository` interface: Create, Update, ListByEmployee, ListPending
- [ ] `LeaveBalanceRepository` interface: Get, Update, ListByEmployee
- [ ] `PayslipRepository` interface: Create, Get, ListByPeriod

**Verification:**
- [ ] `go build ./...` compiles
- [ ] Interfaces are minimal (no implementation details)

**Dependencies:** T1.2
**Files:**
- `internal/domain/interfaces_payroll.go`

**Estimated scope:** S (1-2 files)

---

### T1.5: Payroll Repository Implementations (PG + Memory)

**Description:** Implement PostgreSQL and in-memory repositories for all payroll entities.

**Acceptance criteria:**
- [ ] `PGPayrollInfoRepo` implements all interface methods
- [ ] `MemoryPayrollInfoRepo` implements all interface methods
- [ ] Same pattern for all 8 repository interfaces
- [ ] PG repos use GORM (`*gorm.DB`)
- [ ] Memory repos use `sync.RWMutex` + maps
- [ ] Memory repos copy structs before mutation

**Verification:**
- [ ] `go build ./...` compiles
- [ ] `go test ./internal/repository/ -count=1` passes

**Dependencies:** T1.1, T1.4
**Files:**
- `internal/repository/pg_payroll.go`
- `internal/repository/memory_payroll.go`

**Estimated scope:** L (5-8 files)

---

### Checkpoint 1: Foundation Complete

- [ ] Migration runs cleanly on PG
- [ ] All domain models compile
- [ ] All repository interfaces defined
- [ ] Both PG and Memory repos compile
- [ ] `go vet ./...` passes
- [ ] **REVIEW GATE: Human reviews schema + models before proceeding**

---

### T1.6: Salary Calculation Engine — Income Components

**Description:** Implement gross salary calculation: base salary, OT pay, night shift pay, leave pay, holiday pay, allowances, bonuses.

**Acceptance criteria:**
- [ ] `CalculateBaseSalary(employee, period)` — handles coefficient and fixed types
- [ ] `CalculateOTPay(timekeeping, hourlyRate, dayType)` — 150%/200%/300% multipliers
- [ ] `CalculateNightShiftPay(timekeeping, hourlyRate)` — 30% premium
- [ ] `CalculateLeavePay(baseSalary, leaveDays)` — daily rate calculation
- [ ] `CalculateHolidayPay(baseSalary, holidayDays)` — daily rate calculation
- [ ] `CalculateAllowances(employee)` — sum of all allowance components
- [ ] All functions are pure (no side effects, testable)

**Verification:**
- [ ] Unit tests for each function
- [ ] Test cases: normal, edge (zero days, max OT), boundary values

**Dependencies:** T1.2, T1.3
**Files:**
- `internal/service/payroll_calculator.go`
- `internal/service/payroll_calculator_test.go`

**Estimated scope:** M (3-5 files)

---

### T1.7: Salary Calculation Engine — Insurance

**Description:** Implement SI/HI/UI contribution calculations for both employee and employer.

**Acceptance criteria:**
- [ ] `CalculateEmployeeSI(insuranceBase, isForeign)` — 8% (capped at 20× base salary)
- [ ] `CalculateEmployeeHI(insuranceBase)` — 1.5% (capped at 20× base salary)
- [ ] `CalculateEmployeeUI(uiBase, isForeign)` — 1% Vietnamese only (capped at 20× regional min)
- [ ] `CalculateEmployerSI(insuranceBase, isForeign)` — 17.5% (14% + 3% + 0.5%)
- [ ] `CalculateEmployerHI(insuranceBase)` — 3%
- [ ] `CalculateEmployerUI(uiBase, isForeign)` — 1% Vietnamese only
- [ ] `CalculateInsuranceBase(salary, allowances, cap)` — MIN(income, cap)
- [ ] `CalculateUIBase(salary, allowances, regionalMin, cap)` — MIN(income, cap)
- [ ] Trade union: 1% capped at VND 253,000

**Verification:**
- [ ] Unit tests match official rates from Law 41/2024
- [ ] Test Vietnamese vs foreign employees
- [ ] Test cap boundaries (below, at, above cap)

**Dependencies:** T1.6
**Files:**
- `internal/service/payroll_calculator.go` (add to existing)
- `internal/service/payroll_calculator_test.go` (add to existing)

**Estimated scope:** M (3-5 files)

---

### T1.8: Salary Calculation Engine — PIT

**Description:** Implement PIT progressive 5-bracket calculation with new deductions from Law 109/2025/QH15.

**Acceptance criteria:**
- [ ] `CalculatePIT(taxableIncome)` — 5-bracket progressive
- [ ] `CalculateTaxableIncome(gross, insurance, personalDeduction, dependantDeduction, exemptIncome)` — before PIT
- [ ] Personal deduction: VND 15,500,000/month
- [ ] Per dependant deduction: VND 6,200,000/month
- [ ] OT income exempt from PIT (flag-based)
- [ ] Night shift income exempt from PIT (flag-based)
- [ ] Unused leave income exempt from PIT (flag-based)
- [ ] Quick-deduction shortcut implemented (bracket × rate - constant)
- [ ] Non-resident: flat 20%

**Verification:**
- [ ] Unit tests match HAPRI calculator (https://www.hapri.org/tools-financial/personal-income-tax-calc)
- [ ] Test all 5 brackets
- [ ] Test zero/negative taxable income
- [ ] Test with 0, 1, 2, 3 dependants
- [ ] Test non-resident flat rate

**Dependencies:** T1.7
**Files:**
- `internal/service/payroll_calculator.go` (add to existing)
- `internal/service/payroll_calculator_test.go` (add to existing)

**Estimated scope:** M (3-5 files)

---

### T1.9: Gross-to-Net Pipeline

**Description:** Implement complete gross-to-net calculation pipeline that orchestrates all calculator functions.

**Acceptance criteria:**
- [ ] `CalculatePayrollRun(employee, timekeeping, config)` — full pipeline
- [ ] Returns complete `PayrollRun` with all fields populated
- [ ] Handles all edge cases: zero OT, no dependants, foreign employee, etc.
- [ ] Logs calculation steps for audit trail
- [ ] Returns warnings (salary below min wage, OT exceeds limit, etc.)

**Verification:**
- [ ] Integration test: create employee → run calculation → verify all fields
- [ ] Test with 10 employees, verify totals match manual calculation
- [ ] Performance: <100ms for 100 employees

**Dependencies:** T1.6, T1.7, T1.8
**Files:**
- `internal/service/payroll_calculator.go` (add to existing)
- `internal/service/payroll_calculator_test.go` (add to existing)

**Estimated scope:** M (3-5 files)

---

### T1.10: Net-to-Gross Reverse Calculation

**Description:** Implement reverse calculation: given target net pay, find required gross salary.

**Acceptance criteria:**
- [ ] `CalculateNetToGross(targetNet, employee, config)` — iterative approach
- [ ] Converges within 10 iterations (tolerance: 1 VND)
- [ ] Returns gross salary that produces target net pay
- [ ] Handles edge cases (target too low, too high)

**Verification:**
- [ ] Round-trip test: gross → net → gross produces original gross
- [ ] Test with 5 different salary levels

**Dependencies:** T1.9
**Files:**
- `internal/service/payroll_calculator.go` (add to existing)
- `internal/service/payroll_calculator_test.go` (add to existing)

**Estimated scope:** S (1-2 files)

---

### Checkpoint 2: Calculation Engine Complete

- [ ] All calculator functions implemented
- [ ] Unit test coverage >90% for calculator
- [ ] Gross-to-net matches HAPRI calculator for 50 test cases
- [ ] Performance: <100ms for 100 employees
- [ ] **REVIEW GATE: Chief Accountant verifies calculation accuracy**

---

## Phase 2: Timekeeping & Attendance (Weeks 5-8)

### T2.1: Timekeeping Domain Models

**Description:** Add timekeeping-specific fields and models (if not already in T1.2).

**Acceptance criteria:**
- [ ] `TimekeepingRecord` struct complete
- [ ] `LeaveRequest` struct with status workflow
- [ ] `LeaveBalance` struct with entitled/used/remaining
- [ ] `Holiday` struct for company-specific holidays

**Verification:**
- [ ] `go build ./...` compiles

**Dependencies:** T1.2
**Files:**
- `internal/domain/models_payroll.go` (add to existing)

**Estimated scope:** XS

---

### T2.2: Timekeeping CSV Import

**Description:** Implement CSV/Excel import for timekeeping data.

**Acceptance criteria:**
- [ ] `POST /api/v1/payroll/timekeeping/import` endpoint
- [ ] Accepts CSV with columns: employee_code, date, clock_in, clock_out, ot_hours, night_hours, leave_type
- [ ] Validates format, employee existence, date range
- [ ] Returns import summary (imported count, error count, error details)
- [ ] Supports bulk insert (efficient for 500+ records)
- [ ] Handles duplicates (skip or overwrite option)

**Verification:**
- [ ] Import 100 records, verify all created
- [ ] Import with errors, verify error log
- [ ] Import duplicate, verify handling

**Dependencies:** T1.5, T2.1
**Files:**
- `internal/handler/payroll_handler.go`
- `internal/handler/payroll_handler_test.go`
- `internal/service/payroll_service.go`

**Estimated scope:** M (3-5 files)

---

### T2.3: Overtime Calculation

**Description:** Implement overtime calculation with rate multipliers and limit validation.

**Acceptance criteria:**
- [ ] OT rate multiplier: weekday 150%, rest day 200%, holiday 300%
- [ ] Night OT: additional 20% on top of OT rate
- [ ] Validate: max 40h/month, max 200h/year (300 for specific industries)
- [ ] Warning when OT approaches limit (e.g., >35h/month)
- [ ] OT calculation: hours × hourly_rate × multiplier

**Verification:**
- [ ] Test each day type (weekday, rest day, holiday)
- [ ] Test night OT cumulative
- [ ] Test limit validation (at limit, over limit)
- [ ] Test hourly rate calculation: base_salary / 26 / 8

**Dependencies:** T1.6
**Files:**
- `internal/service/payroll_calculator.go` (add to existing)
- `internal/service/payroll_calculator_test.go` (add to existing)

**Estimated scope:** S

---

### T2.4: Leave Management

**Description:** Implement leave type configuration, balance tracking, and leave pay calculation.

**Acceptance criteria:**
- [ ] Leave types: ANNUAL (12 days), SICK (75% pay), MATERNITY (180 days), PATERNITY (5-14 days), UNPAID, COMPENSATORY
- [ ] Leave balance: entitled, used, remaining, carried over
- [ ] Annual leave: 12 days/year, +1 day per 5 years of service
- [ ] Leave pay calculation: (base_salary / 26) × leave_days
- [ ] Sick leave: requires certificate, paid at 75%
- [ ] Leave request workflow: PENDING → APPROVED/REJECTED

**Verification:**
- [ ] Create leave request, approve, verify balance updated
- [ ] Test annual leave entitlement with different service years
- [ ] Test sick leave pay at 75%
- [ ] Test leave pay calculation accuracy

**Dependencies:** T1.5, T2.1
**Files:**
- `internal/service/payroll_service.go`
- `internal/service/payroll_service_test.go`
- `internal/handler/payroll_handler.go`

**Estimated scope:** M (3-5 files)

---

### T2.5: Holiday Calendar

**Description:** Implement public holiday calendar with company-specific holidays.

**Acceptance criteria:**
- [ ] Pre-loaded 2026 Vietnamese public holidays
- [ ] CRUD API for company-specific holidays
- [ ] Holiday pay calculation: (base_salary / 26) × holiday_days
- [ ] Holidays flagged in timekeeping records
- [ ] Holiday detection for OT rate selection

**Verification:**
- [ ] 2026 holidays pre-loaded correctly
- [ ] Add company holiday, verify it affects payroll
- [ ] Test holiday pay calculation

**Dependencies:** T1.5
**Files:**
- `internal/service/payroll_service.go`
- `internal/handler/payroll_handler.go`

**Estimated scope:** S

---

### Checkpoint 3: Timekeeping Complete

- [ ] CSV import works for 500+ records
- [ ] OT calculation matches manual calculation
- [ ] Leave balance tracking accurate
- [ ] Holiday calendar loaded
- [ ] **REVIEW GATE: BA verifies timekeeping rules against Labour Code 2019**

---

## Phase 3: Payroll Processing & GL (Weeks 9-12)

### T3.1: Payroll Period Management

**Description:** Implement payroll period CRUD with status workflow.

**Acceptance criteria:**
- [ ] Create period (year/month) — auto-check for duplicates
- [ ] Period statuses: DRAFT → PROCESSING → APPROVED → PAID → CLOSED
- [ ] Status transitions enforced (can't skip steps)
- [ ] Period lock after approval (no modifications)
- [ ] List periods with summary stats

**Verification:**
- [ ] Create period, verify status is DRAFT
- [ ] Transition through all statuses
- [ ] Try to modify approved period, verify rejection

**Dependencies:** T1.5
**Files:**
- `internal/service/payroll_service.go`
- `internal/handler/payroll_handler.go`
- `internal/handler/payroll_handler_test.go`

**Estimated scope:** M (3-5 files)

---

### T3.2: Payroll Run Engine

**Description:** Implement payroll run calculation for all employees in a period.

**Acceptance criteria:**
- [ ] `POST /api/v1/payroll/periods/:id/calculate` endpoint
- [ ] Iterates all active employees for the period
- [ ] Calls `CalculatePayrollRun()` for each employee
- [ ] Creates `PayrollRun` records for each employee
- [ ] Generates period summary (total gross, deductions, net pay)
- [ ] Returns warnings for each employee (salary below min, OT limit, etc.)
- [ ] Performance: <30s for 100 employees

**Verification:**
- [ ] Calculate payroll for 10 employees, verify all runs created
- [ ] Verify totals match sum of individual runs
- [ ] Verify warnings generated correctly
- [ ] Benchmark: 100 employees in <30s

**Dependencies:** T1.9, T3.1
**Files:**
- `internal/service/payroll_service.go`
- `internal/service/payroll_service_test.go`

**Estimated scope:** M (3-5 files)

---

### T3.3: Payslip Generation (PDF)

**Description:** Implement payslip PDF generation using maroto/v2.

**Acceptance criteria:**
- [ ] Generate PDF payslip per employee
- [ ] Include all fields: income breakdown, deductions, net pay, employer costs
- [ ] Vietnamese language support
- [ ] Company logo placeholder
- [ ] Payslip number: `{company_code}-{year}-{month}-{sequence}`
- [ ] Store PDF path in payslips table

**Verification:**
- [ ] Generate payslip for test employee
- [ ] Verify PDF opens and displays correctly
- [ ] Verify all amounts match payroll run

**Dependencies:** T1.5, T3.2
**Files:**
- `internal/service/payslip_service.go`
- `internal/service/payslip_service_test.go`

**Estimated scope:** M (3-5 files)

---

### T3.4: GL Journal Entry Integration

**Description:** Implement automatic GL journal entry posting when payroll is approved.

**Acceptance criteria:**
- [ ] On approval, generate journal entries:
  - Dr. 6421 Salary Expense, Dr. 6422 OT Expense, Cr. 3331 Salary Payable
  - Dr. 6424 Insurance Expense, Cr. 3332/3333/3334 (ER insurance)
  - Cr. 3335/3336/3337/3338 (EE insurance + PIT)
- [ ] Use existing `Service.PostJournalEntry()` from GL module
- [ ] Reference: `PAYROLL-{year}-{month}`
- [ ] Validate GL period is open before posting
- [ ] Rollback on posting failure

**Verification:**
- [ ] Approve payroll, verify journal entries in GL
- [ ] Verify GL account balances match payroll summary
- [ ] Test with closed GL period, verify rejection

**Dependencies:** T3.2, GL module
**Files:**
- `internal/service/payroll_service.go`
- `internal/service/payroll_service_test.go`

**Estimated scope:** L (5-8 files)

---

### T3.5: Payroll Approval Workflow

**Description:** Implement multi-level approval: preparer → reviewer → approver.

**Acceptance criteria:**
- [ ] Submit for review: status → PROCESSING
- [ ] Approve: status → APPROVED, GL posted, payslips generated
- [ ] Reject: status → DRAFT, with comments
- [ ] Lock period after approval (no re-processing)
- [ ] Audit trail: who, when, what action

**Verification:**
- [ ] Full workflow: create → calculate → submit → approve
- [ ] Test rejection flow
- [ ] Test lock after approval

**Dependencies:** T3.2, T3.3, T3.4
**Files:**
- `internal/service/payroll_service.go`
- `internal/handler/payroll_handler.go`
- `internal/handler/payroll_handler_test.go`

**Estimated scope:** M (3-5 files)

---

### Checkpoint 4: Payroll Processing Complete

- [ ] Full payroll run produces correct payslips
- [ ] GL journal entries posted correctly
- [ ] Approval workflow functions
- [ ] PDF payslips generate correctly
- [ ] **REVIEW GATE: Chief Accountant verifies GL entries and payslips**

---

## Phase 4: Declarations & Polish (Weeks 13-16)

### T4.1: D02-TS Social Insurance Declaration

**Description:** Generate D02-TS XML for social insurance registration/adjustment.

**Acceptance criteria:**
- [ ] Auto-detect: new employees (registration), salary changes (adjustment), terminated (reduction)
- [ ] Generate D02-TS list with all required fields
- [ ] XML generation in BHXH electronic format
- [ ] Digital signature support (via existing xmldsig)
- [ ] Download XML for I-VAN submission

**Verification:**
- [ ] Generate D02-TS for 5 employees (2 new, 2 adjustments, 1 termination)
- [ ] Verify XML format matches BHXH spec
- [ ] Verify all fields populated correctly

**Dependencies:** T1.5, T3.2
**Files:**
- `internal/service/declaration_service.go`
- `internal/service/declaration_service_test.go`
- `internal/handler/declaration_handler.go`

**Estimated scope:** M (3-5 files)

---

### T4.2: TK3-TS Employer Registration

**Description:** Generate TK3-TS for employer SI/HI registration.

**Acceptance criteria:**
- [ ] One-time registration form
- [ ] Include: enterprise name, tax code, address, employee count
- [ ] XML generation for electronic submission
- [ ] Store unit code after BHXH assigns

**Verification:**
- [ ] Generate TK3-TS for test company
- [ ] Verify XML format

**Dependencies:** T1.5
**Files:**
- `internal/service/declaration_service.go` (add to existing)
- `internal/handler/declaration_handler.go` (add to existing)

**Estimated scope:** S

---

### T4.3: 05/KK-TNCN Quarterly PIT Declaration

**Description:** Generate quarterly PIT declaration form 05/KK-TNCN per Decision 1109/QD-BTC.

**Acceptance criteria:**
- [ ] Aggregate quarterly payroll data per employee
- [ ] Total income, insurance deductions, PIT withheld
- [ ] Exempt income (OT, night shift, leave)
- [ ] Generate Form 05/KK-TNCN XML
- [ ] Support quarterly filing (from Q2/2026)

**Verification:**
- [ ] Generate 05/KK-TNCN for Q3/2026 (3 months data)
- [ ] Verify totals match monthly payroll sums
- [ ] Verify XML format for eTax submission

**Dependencies:** T3.2, T1.8
**Files:**
- `internal/service/declaration_service.go` (add to existing)
- `internal/service/declaration_service_test.go` (add to existing)

**Estimated scope:** M (3-5 files)

---

### T4.4: Year-End Tax Finalization

**Description:** Support annual PIT finalization for employees.

**Acceptance criteria:**
- [ ] Aggregate 12 months of payroll data
- [ ] Calculate annual tax liability
- [ ] Compare with monthly withholdings
- [ ] Determine refund/additional payment
- [ ] Generate annual tax finalization report

**Verification:**
- [ ] Process year-end for test employee
- [ ] Verify annual calculation matches sum of monthly

**Dependencies:** T4.3
**Files:**
- `internal/service/declaration_service.go` (add to existing)
- `internal/service/declaration_service_test.go` (add to existing)

**Estimated scope:** M (3-5 files)

---

### T4.5: Employee Self-Service Portal

**Description:** Implement employee self-service for payslips and leave.

**Acceptance criteria:**
- [ ] `GET /api/v1/payroll/payslips` — list own payslips
- [ ] `GET /api/v1/payroll/payslips/:id` — view payslip detail
- [ ] `GET /api/v1/payroll/payslips/:id/pdf` — download PDF
- [ ] `GET /api/v1/payroll/leave/balance` — view leave balances
- [ ] `POST /api/v1/payroll/leave/request` — submit leave request
- [ ] Auth: employee can only see own data

**Verification:**
- [ ] Employee logs in, views payslip
- [ ] Employee downloads PDF
- [ ] Employee submits leave request
- [ ] Employee cannot see other employee's data

**Dependencies:** T3.3, T2.4
**Files:**
- `internal/handler/payroll_ess_handler.go`
- `internal/handler/payroll_ess_handler_test.go`

**Estimated scope:** M (3-5 files)

---

### T4.6: Bank Salary File Generation

**Description:** Generate bank salary file for salary disbursement.

**Acceptance criteria:**
- [ ] CSV format with: sequence, account number, name, amount, reference
- [ ] Support multiple bank formats (VCB, CTG, VTB)
- [ ] Include only approved payroll runs
- [ ] Total amount matches payroll summary

**Verification:**
- [ ] Generate bank file, verify format
- [ ] Verify total matches payroll net pay

**Dependencies:** T3.2
**Files:**
- `internal/service/payroll_service.go` (add to existing)
- `internal/handler/payroll_handler.go` (add to existing)

**Estimated scope:** S

---

### T4.7: Payroll Reports

**Description:** Implement payroll reporting: summary, cost allocation, insurance, PIT, overtime.

**Acceptance criteria:**
- [ ] `GET /api/v1/payroll/reports/summary` — period summary
- [ ] `GET /api/v1/payroll/reports/cost-by-department` — cost allocation
- [ ] `GET /api/v1/payroll/reports/insurance-summary` — SI/HI/UI totals
- [ ] `GET /api/v1/payroll/reports/pit-summary` — PIT by employee
- [ ] `GET /api/v1/payroll/reports/overtime` — OT hours by employee
- [ ] `GET /api/v1/payroll/reports/leave-balance` — leave balances

**Verification:**
- [ ] Generate each report type
- [ ] Verify data matches payroll runs

**Dependencies:** T3.2
**Files:**
- `internal/handler/payroll_handler.go` (add to existing)
- `internal/handler/payroll_handler_test.go` (add to existing)

**Estimated scope:** M (3-5 files)

---

### T4.8: Payroll Configuration API

**Description:** Implement configuration management for rates, caps, holidays.

**Acceptance criteria:**
- [ ] `GET /api/v1/payroll/config/rates` — current rates
- [ ] `PUT /api/v1/payroll/config/rates` — update rates
- [ ] `GET /api/v1/payroll/config/holidays` — list holidays
- [ ] `POST /api/v1/payroll/config/holidays` — add holiday
- [ ] `GET /api/v1/payroll/config/regions` — list regions with min wages
- [ ] Rate changes effective from date

**Verification:**
- [ ] Update SI rate, verify new calculation uses updated rate
- [ ] Add company holiday, verify it affects payroll

**Dependencies:** T1.5
**Files:**
- `internal/handler/payroll_handler.go` (add to existing)
- `internal/handler/payroll_handler_test.go` (add to existing)

**Estimated scope:** S

---

### T4.9: Handler Tests (Comprehensive)

**Description:** Write comprehensive handler tests for all payroll endpoints.

**Acceptance criteria:**
- [ ] Test all CRUD operations
- [ ] Test calculation endpoints
- [ ] Test declaration generation
- [ ] Test error cases (unauthorized, invalid data, missing employee)
- [ ] Test with in-memory repos (no DB dependency)

**Verification:**
- [ ] `go test ./internal/handler/ -count=1 -run TestPayroll` passes
- [ ] Coverage >80% for payroll handlers

**Dependencies:** T4.1-T4.8
**Files:**
- `internal/handler/payroll_handler_test.go`

**Estimated scope:** L (5-8 files)

---

### T4.10: Integration Tests & Documentation

**Description:** End-to-end integration tests and API documentation.

**Acceptance criteria:**
- [ ] Full flow test: create employee → import timekeeping → calculate → approve → GL → payslip
- [ ] Declaration generation test
- [ ] Swagger annotations for all endpoints
- [ ] `swag init` generates updated docs

**Verification:**
- [ ] `go test -count=1 ./...` all pass
- [ ] `swag init --parseDependency --parseInternal` succeeds
- [ ] Swagger UI shows payroll endpoints

**Dependencies:** T4.9
**Files:**
- `internal/handler/payroll_handler.go` (add swagger annotations)
- `docs/docs.go` (regenerated)
- `docs/swagger.json` (regenerated)
- `docs/swagger.yaml` (regenerated)

**Estimated scope:** M (3-5 files)

---

### Checkpoint 5: Declarations Complete

- [ ] D02-TS generates correctly
- [ ] 05/KK-TNCN generates correctly
- [ ] Employee self-service works
- [ ] Bank file generates correctly
- [ ] **REVIEW GATE: Tax advisor verifies declaration forms**

---

### Checkpoint 6: PROD Ready

- [ ] All 40 tasks complete
- [ ] All tests pass
- [ ] Swagger docs updated
- [ ] AGENTS.md updated with Payroll status = PROD
- [ ] **FINAL GATE: Chief Accountant + BA sign-off**

---

## Task Summary

| Phase | Tasks | Duration | Key Deliverable |
|-------|-------|----------|-----------------|
| 1 | T1.1-T1.10 | Weeks 1-4 | Gross-to-net calculation engine |
| 2 | T2.1-T2.5 | Weeks 5-8 | Timekeeping import + OT + leave |
| 3 | T3.1-T3.5 | Weeks 9-12 | Payroll processing + GL integration |
| 4 | T4.1-T4.10 | Weeks 13-16 | Declarations + self-service + polish |
| **Total** | **40 tasks** | **16 weeks** | **PROD-ready payroll module** |
