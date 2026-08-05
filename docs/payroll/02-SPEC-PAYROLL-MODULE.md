# Payroll Module — Detailed Specification

**Document ID:** SPEC-PAYROLL-001
**Version:** 1.0
**Date:** 2026-08-04

---

## 1. Domain Model

### 1.1 Extended Employee Model

```go
// Extends existing Employee in models_company.go
type EmployeePayrollInfo struct {
    ID                 string    `json:"id"`
    EmployeeID         string    `json:"employee_id"`         // FK to Employee
    ContractType       string    `json:"contract_type"`       // INDEFINITE, DEFINITE, PROBATION
    ContractStartDate  string    `json:"contract_start_date"`
    ContractEndDate    string    `json:"contract_end_date,omitempty"`
    SalaryType         string    `json:"salary_type"`         // TIME_BASED, PIECE_RATE, COMMISSION, COEFFICIENT
    BaseSalary         float64   `json:"base_salary"`         // VND/month
    SalaryCoefficient  float64   `json:"salary_coefficient"` // e.g. 2.34
    PositionAllowance  float64   `json:"position_allowance"`
    ResponsibilityAllowance float64 `json:"responsibility_allowance"`
    SeniorityAllowance float64   `json:"seniority_allowance"`
    OtherAllowances    float64   `json:"other_allowances"`
    InsuranceBaseSalary float64  `json:"insurance_base_salary"` // may differ from contract salary
    Region             string    `json:"region"`               // I, II, III, IV
    IsForeignEmployee  bool      `json:"is_foreign_employee"`
    IsTradeUnionMember bool      `json:"is_trade_union_member"`
    IsHighTechTalent   bool      `json:"is_high_tech_talent"` // 5-year PIT exemption
    BankAccountNo      string    `json:"bank_account_no"`
    BankCode           string    `json:"bank_code"`
    EffectiveDate      string    `json:"effective_date"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
}
```

### 1.2 Dependant Model

```go
type Dependant struct {
    ID           string `json:"id"`
    EmployeeID   string `json:"employee_id"`
    FullName     string `json:"full_name"`
    Relationship string `json:"relationship"` // CHILD, SPOUSE, PARENT
    DateOfBirth  string `json:"date_of_birth"`
    TaxCode      string `json:"tax_code,omitempty"`
    IsDisabled   bool   `json:"is_disabled"`
    IsActive     bool   `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
}
```

### 1.3 Salary Structure Models

```go
type SalaryComponent struct {
    ID          string  `json:"id"`
    CompanyID   string  `json:"company_id"`
    Code        string  `json:"code"`        // e.g. "BS", "TA", "PC_CN"
    Name        string  `json:"name"`        // e.g. "Lương cơ bản", "Phụ cấp ăn trưa"
    Type        string  `json:"type"`        // INCOME, DEDUCTION, EMPLOYER_COST
    Calculation string  `json:"calculation"` // FIXED, PERCENTAGE, FORMULA
    Formula     string  `json:"formula,omitempty"` // e.g. "base_salary * coefficient"
    DefaultValue float64 `json:"default_value,omitempty"`
    IsTaxable   bool    `json:"is_taxable"`
    IsInsurable bool    `json:"is_insurable"` // counts toward insurance base
    IsActive    bool    `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
}

type SalaryTemplate struct {
    ID          string  `json:"id"`
    CompanyID   string  `json:"company_id"`
    Name        string  `json:"name"`        // e.g. "Office Staff", "Factory Worker"
    Components  []SalaryTemplateComponent `json:"components"`
    CreatedAt   time.Time `json:"created_at"`
}

type SalaryTemplateComponent struct {
    ID               string  `json:"id"`
    SalaryTemplateID string  `json:"salary_template_id"`
    ComponentID      string  `json:"component_id"`
    DefaultValue     float64 `json:"default_value"`
    Formula          string  `json:"formula,omitempty"`
    Order            int     `json:"order"`
}
```

### 1.4 Payroll Run Models

```go
type PayrollPeriod struct {
    ID          string           `json:"id"`
    CompanyID   string           `json:"company_id"`
    Year        int              `json:"year"`
    Month       int              `json:"month"`
    Status      PayrollStatus    `json:"status"` // DRAFT, PROCESSING, APPROVED, PAID, CLOSED
    PreparedBy  string           `json:"prepared_by"`
    PreparedAt  *time.Time       `json:"prepared_at,omitempty"`
    ReviewedBy  string           `json:"reviewed_by,omitempty"`
    ReviewedAt  *time.Time       `json:"reviewed_at,omitempty"`
    ApprovedBy  string           `json:"approved_by,omitempty"`
    ApprovedAt  *time.Time       `json:"approved_at,omitempty"`
    PaidAt      *time.Time       `json:"paid_at,omitempty"`
    ClosedAt    *time.Time       `json:"closed_at,omitempty"`
    CreatedAt   time.Time        `json:"created_at"`
    UpdatedAt   time.Time        `json:"updated_at"`
}

type PayrollRun struct {
    ID              string          `json:"id"`
    PeriodID        string          `json:"period_id"`
    EmployeeID      string          `json:"employee_id"`
    CompanyID       string          `json:"company_id"`
    
    // Input data
    WorkingDays     int             `json:"working_days"`      // days worked
    OTHours         float64         `json:"ot_hours"`          // overtime hours
    NightShiftHours float64         `json:"night_shift_hours"`
    LeaveDays       float64         `json:"leave_days"`        // paid leave days
    UnpaidLeaveDays float64         `json:"unpaid_leave_days"`
    AbsentDays      int             `json:"absent_days"`       // absent without leave
    
    // Income
    BaseSalary      float64         `json:"base_salary"`
    OTpay           float64         `json:"ot_pay"`
    NightShiftPay   float64         `json:"night_shift_pay"`
    LeavePay        float64         `json:"leave_pay"`
    HolidayPay      float64         `json:"holiday_pay"`
    Allowances      float64         `json:"allowances"`        // total allowances
    Bonuses         float64         `json:"bonuses"`
    OtherIncome     float64         `json:"other_income"`
    GrossSalary     float64         `json:"gross_salary"`
    
    // Employee deductions
    SIDeduction     float64         `json:"si_deduction"`      // 8% employee
    HIDeduction     float64         `json:"hi_deduction"`      // 1.5% employee
    UIDeduction     float64         `json:"ui_deduction"`      // 1% employee
    TradeUnionDues  float64         `json:"trade_union_dues"`  // 1% (members)
    PITAmount       float64         `json:"pit_amount"`
    OtherDeductions float64         `json:"other_deductions"`
    TotalDeductions float64         `json:"total_deductions"`
    
    // Net pay
    NetPay          float64         `json:"net_pay"`
    
    // Employer costs
    EmployerSI      float64         `json:"employer_si"`       // 17.5%
    EmployerHI      float64         `json:"employer_hi"`       // 3%
    EmployerUI      float64         `json:"employer_ui"`       // 1%
    EmployerTradeUnion float64      `json:"employer_trade_union"` // 1%
    TotalEmployerCost float64       `json:"total_employer_cost"`
    
    // Insurance base (for declaration)
    InsuranceBase   float64         `json:"insurance_base"`    // capped at 20× base salary
    UIBase          float64         `json:"ui_base"`           // capped at 20× regional min
    
    // Status
    Status          string          `json:"status"` // DRAFT, CALCULATED, APPROVED, PAID
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
}

type Payslip struct {
    ID          string    `json:"id"`
    RunID       string    `json:"run_id"`
    EmployeeID  string    `json:"employee_id"`
    PeriodID    string    `json:"period_id"`
    PayslipNo   string    `json:"payslip_no"`
    PDFPath     string    `json:"pdf_path,omitempty"`
    SentAt      *time.Time `json:"sent_at,omitempty"`
    ViewedAt    *time.Time `json:"viewed_at,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### 1.5 Timekeeping Models

```go
type TimekeepingRecord struct {
    ID          string    `json:"id"`
    EmployeeID  string    `json:"employee_id"`
    CompanyID   string    `json:"company_id"`
    Date        string    `json:"date"`        // YYYY-MM-DD
    ClockIn     string    `json:"clock_in"`    // HH:MM
    ClockOut    string    `json:"clock_out"`   // HH:MM
    HoursWorked float64   `json:"hours_worked"`
    OTHours     float64   `json:"ot_hours"`
    NightHours  float64   `json:"night_hours"`
    IsHoliday   bool      `json:"is_holiday"`
    IsRestDay   bool      `json:"is_rest_day"`
    LeaveType   string    `json:"leave_type,omitempty"` // ANNUAL, SICK, MATERNITY, UNPAID, COMPENSATORY
    Notes       string    `json:"notes,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

type LeaveRequest struct {
    ID          string    `json:"id"`
    EmployeeID  string    `json:"employee_id"`
    CompanyID   string    `json:"company_id"`
    LeaveType   string    `json:"leave_type"`
    StartDate   string    `json:"start_date"`
    EndDate     string    `json:"end_date"`
    Days        float64   `json:"days"`
    Reason      string    `json:"reason"`
    Status      string    `json:"status"` // PENDING, APPROVED, REJECTED
    ApprovedBy  string    `json:"approved_by,omitempty"`
    ApprovedAt  *time.Time `json:"approved_at,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

type LeaveBalance struct {
    ID            string  `json:"id"`
    EmployeeID    string  `json:"employee_id"`
    Year          int     `json:"year"`
    LeaveType     string  `json:"leave_type"`
    Entitled      float64 `json:"entitled"`      // total entitled days
    Used          float64 `json:"used"`           // days used
    Remaining     float64 `json:"remaining"`      // entitled - used
    CarriedOver   float64 `json:"carried_over"`   // from previous year
}
```

---

## 2. Calculation Formulas

### 2.1 Gross-to-Net Calculation

```
STEP 1: Calculate Income Components
  base_salary = Employee.BaseSalary (or coefficient × reference salary)
  ot_pay = Σ(ot_hours × hourly_rate × ot_multiplier)
    where hourly_rate = base_salary / 26 / 8
    ot_multiplier = 1.5 (weekday) | 2.0 (rest day) | 3.0 (holiday)
  night_shift_pay = Σ(night_hours × hourly_rate × 0.30)
  leave_pay = (base_salary / 26) × paid_leave_days
  holiday_pay = (base_salary / 26) × holiday_days
  allowances = Σ(allowance_components)
  bonuses = Σ(bonus_components)
  other_income = Σ(other_income_components)
  
  gross_salary = base_salary + ot_pay + night_shift_pay + leave_pay 
                 + holiday_pay + allowances + bonuses + other_income

STEP 2: Calculate Insurance Base
  insurance_base = MIN(contract_salary + regular_allowances, 20 × base_salary)
  ui_base = MIN(contract_salary + regular_allowances, 20 × regional_min_wage)

STEP 3: Calculate Employee Insurance Deductions
  si_deduction = insurance_base × 0.08
  hi_deduction = insurance_base × 0.015
  ui_deduction = ui_base × 0.01  (Vietnamese employees only)
  trade_union = insurance_base × 0.01 (members only, capped at 253,000)

STEP 4: Calculate Taxable Income
  total_insurance = si_deduction + hi_deduction + ui_deduction
  personal_deduction = 15,500,000
  dependant_deduction = 6,200,000 × number_of_dependants
  
  taxable_income = gross_salary - total_insurance - personal_deduction 
                   - dependant_deduction - ot_income_exempt 
                   - night_shift_income_exempt - unused_leave_exempt

STEP 5: Calculate PIT
  ⚠️ TRANSITION: Apply different brackets depending on pay period.
  - Jan-Jun 2026: OLD 7-bracket schedule (Law 103/2016/QH13)
  - Jul-Dec 2026: NEW 5-bracket schedule (Law 109/2025/QH15)

  NEW SCHEDULE (from 01 Jul 2026):
  if taxable_income <= 0:
    pit = 0
  elif taxable_income <= 10,000,000:
    pit = taxable_income × 0.05
  elif taxable_income <= 30,000,000:
    pit = taxable_income × 0.10 - 500,000
  elif taxable_income <= 60,000,000:
    pit = taxable_income × 0.20 - 3,500,000
  elif taxable_income <= 100,000,000:
    pit = taxable_income × 0.30 - 9,500,000
  else:
    pit = taxable_income × 0.35 - 14,500,000

  OLD SCHEDULE (01 Jan – 30 Jun 2026):
  if taxable_income <= 0:
    pit = 0
  elif taxable_income <= 5,000,000:
    pit = taxable_income × 0.05
  elif taxable_income <= 10,000,000:
    pit = taxable_income × 0.10 - 250,000
  elif taxable_income <= 18,000,000:
    pit = taxable_income × 0.15 - 750,000
  elif taxable_income <= 32,000,000:
    pit = taxable_income × 0.20 - 1,650,000
  elif taxable_income <= 52,000,000:
    pit = taxable_income × 0.25 - 3,250,000
  elif taxable_income <= 80,000,000:
    pit = taxable_income × 0.30 - 5,850,000
  else:
    pit = taxable_income × 0.35 - 9,850,000

STEP 6: Calculate Net Pay
  total_deductions = si_deduction + hi_deduction + ui_deduction 
                     + trade_union + pit + other_deductions
  net_pay = gross_salary - total_deductions

STEP 7: Calculate Employer Costs
  employer_si = insurance_base × 0.175
  employer_hi = insurance_base × 0.03
  employer_ui = ui_base × 0.01  (Vietnamese employees only)
  employer_trade_union = insurance_base × 0.01 (members only)
  total_employer_cost = gross_salary + employer_si + employer_hi 
                        + employer_ui + employer_trade_union
```

### 2.2 Net-to-Gross Calculation (Reverse)

```
Given: net_pay target, known deductions
Iterative approach (binary search or Newton's method):
  1. Estimate gross_salary
  2. Run gross-to-net calculation
  3. Compare result with target net_pay
  4. Adjust estimate and repeat until converged (tolerance: 1 VND)
```

### 2.3 Overtime Multiplier Rules

```go
func OTMultiplier(dayType string, isNight bool) float64 {
    base := 1.0
    switch dayType {
    case "WEEKDAY":
        base = 1.5
    case "REST_DAY":
        base = 2.0
    case "HOLIDAY":
        base = 3.0
    }
    if isNight {
        // Night premium: +30% of hourly rate
        // Night OT: additional +20% of the OT rate
        base += 0.20  // additional 20% for night OT
    }
    return base
}
```

---

## 3. Use Cases

### UC-01: Process Monthly Payroll

**Actor:** Payroll Clerk
**Precondition:** Employee data exists, timekeeping data imported, period is DRAFT
**Postcondition:** Payroll calculated, journal entries created, payslips generated

#### Happy Path
1. Payroll Clerk selects payroll period (month/year)
2. System displays list of active employees for the period
3. Clerk imports timekeeping data (CSV/Excel) or confirms manual entry
4. System validates timekeeping data (no gaps, no overlaps, within limits)
5. Clerk clicks "Calculate Payroll"
6. System calculates for each employee:
   - Income components (base, OT, night shift, leave, allowances)
   - Insurance deductions (SI, HI, UI)
   - PIT calculation with new deductions
   - Net pay
   - Employer costs
7. System displays payroll summary (total cost, total deductions, total net pay)
8. Clerk reviews and clicks "Submit for Review"
9. System generates journal entries (Dr Salary Expense, Cr Payables)
10. System generates payslips (PDF)
11. System status → PROCESSING

#### Alternative Paths
- **A1: Clerk needs to adjust individual employee payroll**
  1. Clerk clicks employee name
  2. System shows individual payroll detail
  3. Clerk modifies amounts (with reason)
  4. System recalculates totals
  5. Adjustment logged in audit trail

- **A2: New employee hired mid-month**
  1. System auto-prorates salary based on hire date
  2. Insurance registration for new employee included in D02-TS

- **A3: Employee terminated mid-month**
  1. System calculates final pay including:
     - Prorated salary
     - Unused annual leave payout
     - Severance pay (if applicable)
     - Final insurance contributions

#### Exception Paths
- **E1: Timekeeping data has errors**
  1. System validates and highlights errors
  2. Clerk must correct before proceeding
  3. Common errors: missing clock-out, OT exceeds 40h/month, negative hours

- **E2: Employee missing insurance base salary**
  1. System flags employee
  2. Clerk must set insurance base before calculation

- **E3: Regional minimum wage violation**
  1. System checks if employee salary < regional minimum
  2. System flags violation and suggests minimum wage supplement

---

### UC-02: Import Timekeeping Data

**Actor:** Payroll Clerk
**Precondition:** Timekeeping data available from external source
**Postcondition:** Timekeeping records created for payroll period

#### Happy Path
1. Clerk navigates to Timekeeping → Import
2. Clerk selects file (CSV/Excel)
3. System validates file format
4. System displays preview of records
5. Clerk confirms import
6. System creates timekeeping records
7. System shows import summary (X records imported, Y errors)

#### Alternative Paths
- **A1: Partial import (some records have errors)**
  1. System imports valid records
  2. System shows error log for invalid records
  3. Clerk corrects and re-imports errors

#### Exception Paths
- **E1: File format invalid**
  1. System rejects file
  2. System shows expected format template

- **E2: Duplicate records**
  1. System detects duplicates
  2. Clerk chooses: skip, overwrite, or merge

---

### UC-03: Generate D02-TS Social Insurance Declaration

**Actor:** Payroll Clerk
**Precondition:** Payroll approved, employees registered
**Postcondition:** D02-TS XML generated for electronic submission

#### Happy Path
1. Clerk navigates to Declarations → D02-TS
2. System auto-detects:
   - New employees (need registration)
   - Salary changes (need adjustment)
   - Terminated employees (need deregistration)
3. System generates D02-TS list
4. Clerk reviews and approves
5. System generates XML for I-VAN submission
6. Clerk downloads or submits electronically

#### Alternative Paths
- **A1: Manual adjustment to D02-TS**
  1. Clerk edits employee entries
  2. System recalculates totals
  3. Change logged in audit trail

#### Exception Paths
- **E1: Employee missing SI number**
  1. System flags employee
  2. Clerk must obtain SI number before submission

---

### UC-04: Generate 05/KK-TNCN Quarterly PIT Declaration

**Actor:** Payroll Clerk
**Precondition:** Payroll processed for quarter
**Postcondition:** 05/KK-TNCN form generated

#### Happy Path
1. Clerk navigates to Declarations → 05/KK-TNCN
2. System aggregates quarterly payroll data:
   - Total income per employee
   - Total insurance deductions
   - Total PIT withheld
   - Exempt income (OT, night shift, leave)
3. System generates Form 05/KK-TNCN
4. Clerk reviews and approves
5. System generates XML for eTax submission
6. Clerk downloads or submits electronically

#### Alternative Paths
- **A1: Employee has multiple income sources**
  1. System consolidates all income sources
  2. System calculates annual tax finalization

#### Exception Paths
- **E1: Mid-quarter employee addition**
  1. System includes employee from start date
  2. Declaration reflects partial period

---

### UC-05: Approve Payroll

**Actor:** Chief Accountant / Finance Manager
**Precondition:** Payroll submitted for review
**Postcondition:** Payroll approved, ready for payment

#### Happy Path
1. Approver receives notification
2. Approver reviews payroll summary
3. Approver drills down into individual employees
4. Approver approves payroll
5. System:
   - Locks payroll period (no changes)
   - Posts journal entries to GL
   - Generates payslips
   - Updates leave balances
   - Status → APPROVED

#### Alternative Paths
- **A1: Approver requests changes**
  1. Approver adds comments
  2. System returns payroll to PREPARER
  3. Status → REVISION_REQUIRED

#### Exception Paths
- **E1: Payroll period already closed in GL**
  1. System rejects approval
  2. System shows error: "GL period closed"

---

### UC-06: Employee Views Payslip (Self-Service)

**Actor:** Employee
**Precondition:** Payroll approved, employee authenticated
**Postcondition:** Employee views/downloads payslip

#### Happy Path
1. Employee logs in to self-service portal
2. Employee navigates to Payslips
3. System shows list of available payslips
4. Employee selects month
5. System displays payslip detail
6. Employee downloads PDF

#### Alternative Paths
- **A1: Employee disputes payslip amount**
  1. Employee clicks "Dispute"
  2. Employee enters dispute reason
  3. System sends notification to Payroll Clerk
  4. Clerk investigates and responds

---

## 4. Process Flows

### 4.1 Monthly Payroll Process

```
┌─────────────────────────────────────────────────────────────┐
│                    MONTHLY PAYROLL CYCLE                      │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Day 1-5: Data Collection                                    │
│  ├── Import timekeeping data                                  │
│  ├── Import attendance data                                   │
│  ├── Collect leave requests                                   │
│  ├── Update employee changes (hire/terminate/transfer)        │
│  └── Update salary changes                                    │
│                                                               │
│  Day 5-10: Calculation                                       │
│  ├── Validate timekeeping data                                │
│  ├── Calculate gross-to-net for all employees                 │
│  ├── Calculate insurance contributions                        │
│  ├── Calculate PIT withholding                                │
│  ├── Generate payroll summary                                 │
│  └── Review and adjust (if needed)                            │
│                                                               │
│  Day 10-15: Approval                                         │
│  ├── Payroll Clerk submits for review                         │
│  ├── Chief Accountant reviews                                 │
│  ├── Finance Manager approves                                 │
│  ├── System posts journal entries to GL                       │
│  └── System generates payslips                                │
│                                                               │
│  Day 15-20: Payment                                          │
│  ├── Generate bank salary file                                │
│  ├── Upload to banking system                                 │
│  ├── Confirm payments                                         │
│  └── Record payment in GL                                     │
│                                                               │
│  Day 20-25: Declarations                                      │
│  ├── Generate D02-TS (if new employees)                       │
│  ├── Generate 05/KK-TNCN (quarterly)                          │
│  ├── Submit to social insurance authority                     │
│  └── Submit to tax authority                                   │
│                                                               │
│  Day 25-30: Closing                                           │
│  ├── Update leave balances                                    │
│  ├── Archive payslips                                         │
│  ├── Close payroll period                                     │
│  └── Generate month-end reports                               │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 Gross-to-Net Calculation Flow

```
┌─────────────────────────────────────────────────────────────┐
│              GROSS-TO-NET CALCULATION FLOW                    │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Input: Employee Record + Timekeeping + Salary Structure      │
│                                                               │
│  ┌──────────────────┐                                         │
│  │ Calculate Income  │                                         │
│  │ Components        │                                         │
│  └────────┬─────────┘                                         │
│           │                                                    │
│           ▼                                                    │
│  ┌──────────────────┐     ┌──────────────────┐                │
│  │ Base Salary       │────▶│ OT Pay           │                │
│  │ (coefficient ×    │     │ (hours × rate ×  │                │
│  │  reference)       │     │  multiplier)     │                │
│  └────────┬─────────┘     └────────┬─────────┘                │
│           │                        │                           │
│           ▼                        ▼                           │
│  ┌──────────────────┐     ┌──────────────────┐                │
│  │ Night Shift Pay   │────▶│ Leave Pay        │                │
│  │ (hours × 30%)    │     │ (daily × days)   │                │
│  └────────┬─────────┘     └────────┬─────────┘                │
│           │                        │                           │
│           ▼                        ▼                           │
│  ┌──────────────────┐     ┌──────────────────┐                │
│  │ Allowances        │────▶│ Bonuses          │                │
│  └────────┬─────────┘     └────────┬─────────┘                │
│           │                        │                           │
│           ▼                        ▼                           │
│  ┌──────────────────────────────────────────┐                  │
│  │           GROSS SALARY                    │                  │
│  └──────────────────┬───────────────────────┘                  │
│                     │                                          │
│  ┌──────────────────▼───────────────────────┐                  │
│  │        Insurance Base Calculation         │                  │
│  │  insurance_base = MIN(salary, 20 × ref)  │                  │
│  └──────────────────┬───────────────────────┘                  │
│                     │                                          │
│  ┌──────────────────▼───────────────────────┐                  │
│  │     Employee Insurance Deductions         │                  │
│  │  SI: 8% | HI: 1.5% | UI: 1%             │                  │
│  └──────────────────┬───────────────────────┘                  │
│                     │                                          │
│  ┌──────────────────▼───────────────────────┐                  │
│  │         PIT Calculation                   │                  │
│  │  taxable = gross - insurance - deductions │                  │
│  │  Apply 5-bracket progressive schedule     │                  │
│  └──────────────────┬───────────────────────┘                  │
│                     │                                          │
│  ┌──────────────────▼───────────────────────┐                  │
│  │           NET PAY                         │                  │
│  │  net = gross - all deductions             │                  │
│  └──────────────────────────────────────────┘                  │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### 4.3 GL Journal Entry Flow

```
┌─────────────────────────────────────────────────────────────┐
│              PAYROLL JOURNAL ENTRIES                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Entry 1: Salary Expense Accrual                              │
│  ┌─────────────────────────────────────────┐                  │
│  │ Dr. 6421 - Salary Expense    XXX        │                  │
│  │ Dr. 6422 - OT Expense        XXX        │                  │
│  │ Dr. 6423 - Allowance Expense XXX        │                  │
│  │     Cr. 3331 - Salary Payable    XXX    │                  │
│  └─────────────────────────────────────────┘                  │
│                                                               │
│  Entry 2: Insurance Expense (Employer)                        │
│  ┌─────────────────────────────────────────┐                  │
│  │ Dr. 6424 - Insurance Expense  XXX       │                  │
│  │     Cr. 3332 - SI Payable (ER)   XXX    │                  │
│  │     Cr. 3333 - HI Payable (ER)   XXX    │                  │
│  │     Cr. 3334 - UI Payable (ER)   XXX    │                  │
│  └─────────────────────────────────────────┘                  │
│                                                               │
│  Entry 3: Employee Deductions                                 │
│  ┌─────────────────────────────────────────┐                  │
│  │     Cr. 3335 - SI Payable (EE)   XXX    │                  │
│  │     Cr. 3336 - HI Payable (EE)   XXX    │                  │
│  │     Cr. 3337 - UI Payable (EE)   XXX    │                  │
│  │     Cr. 3338 - PIT Payable       XXX    │                  │
│  │ (offset against Salary Payable)          │                  │
│  └─────────────────────────────────────────┘                  │
│                                                               │
│  Entry 4: Payment                                            │
│  ┌─────────────────────────────────────────┐                  │
│  │ Dr. 3331 - Salary Payable       XXX     │                  │
│  │     Cr. 1111 - Bank Account      XXX    │                  │
│  └─────────────────────────────────────────┘                  │
│                                                               │
│  Entry 5: Insurance Payment                                   │
│  ┌─────────────────────────────────────────┐                  │
│  │ Dr. 3332 - SI Payable (ER)      XXX     │                  │
│  │ Dr. 3335 - SI Payable (EE)      XXX     │                  │
│  │ Dr. 3333 - HI Payable (ER)      XXX     │                  │
│  │ Dr. 3336 - HI Payable (EE)      XXX     │                  │
│  │ Dr. 3334 - UI Payable (ER)      XXX     │                  │
│  │ Dr. 3337 - UI Payable (EE)      XXX     │                  │
│  │     Cr. 1111 - Bank Account      XXX    │                  │
│  └─────────────────────────────────────────┘                  │
│                                                               │
│  Entry 6: PIT Payment                                         │
│  ┌─────────────────────────────────────────┐                  │
│  │ Dr. 3338 - PIT Payable          XXX     │                  │
│  │     Cr. 1111 - Bank Account      XXX    │                  │
│  └─────────────────────────────────────────┘                  │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## 5. Business Rules

### 5.1 Salary Rules

| Rule ID | Rule | Law Reference |
|---------|------|---------------|
| SR-01 | Employee salary must not be below regional minimum wage | Decree 293/2025/ND-CP |
| SR-02 | Insurance base salary capped at 20 × base salary | Law 41/2024/QH15 Art. 31 |
| SR-03 | UI base salary capped at 20 × regional minimum wage | Law 74/2025/QH15 Art. 34 |
| SR-04 | Overtime pay not subject to insurance contributions | Labour Code 2019 Art. 98 |
| SR-05 | Night shift premium: minimum 30% of hourly rate | Labour Code 2019 Art. 97 |
| SR-06 | OT on weekdays: minimum 150% | Labour Code 2019 Art. 98 |
| SR-07 | OT on rest days: minimum 200% | Labour Code 2019 Art. 98 |
| SR-08 | OT on holidays: minimum 300% | Labour Code 2019 Art. 98 |
| SR-09 | Max OT: 40 hours/month, 200 hours/year | Labour Code 2019 Art. 107 |
| SR-10 | Total work hours cannot exceed 12 hours/day | Labour Code 2019 Art. 105 |

### 5.2 Insurance Rules

| Rule ID | Rule | Law Reference |
|---------|------|---------------|
| IR-01 | Vietnamese employees: SI 8% + HI 1.5% + UI 1% = 10.5% | Law 41/2024, Law 74/2025, Law 51/2024 |
| IR-02 | Foreign employees (12+ months): SI 8% + HI 1.5% = 9.5% (no UI) | Law 74/2025 Art. 2 |
| IR-03 | Employer Vietnamese: SI 17.5% + HI 3% + UI 1% = 21.5% | Law 41/2024, Law 74/2025 |
| IR-04 | Employer Foreign: SI 17.5% + HI 3% = 20.5% | Law 41/2024 |
| IR-05 | Insurance base cannot be below minimum wage | Decree 158/2025/ND-CP |
| IR-06 | Insurance registration within 30 days of contract | Law 41/2024 Art. 30 |
| IR-07 | Late payment penalty: 0.03% per day | Law 41/2024 Art. 38 |

### 5.3 PIT Rules

| Rule ID | Rule | Law Reference |
|---------|------|---------------|
| PR-01 | Resident employees: progressive schedule (5-bracket from Jul, 7-bracket Jan-Jun) | Law 109/2025 Art. 9, Law 103/2016 |
| PR-02 | Non-resident employees: flat 20% | Law 109/2025 Art. 9 |
| PR-03 | Personal deduction: VND 15,500,000/month (from 01 Jan 2026) | Resolution 110/2025 |
| PR-04 | Dependant deduction: VND 6,200,000/month each (from 01 Jan 2026) | Resolution 110/2025 |
| PR-05 | OT income fully exempt from PIT (from 01 Jul 2026 only) | Law 109/2025 Art. 4 |
| PR-06 | Night shift income fully exempt (from 01 Jul 2026 only) | Law 109/2025 Art. 4 |
| PR-07 | Unused leave income fully exempt (from 01 Jul 2026 only) | Law 109/2025 Art. 4 |
| PR-08 | High-tech talent: 5-year PIT exemption | Law 109/2025 Art. 5 |
| PR-09 | Supplementary pension deduction: up to VND 3M/month | Decree 253/2026 |
| PR-10 | Quarterly PIT declaration (from Q2/2026) | Resolution 66.16/NQ-CP |
| PR-11 | Annual PIT finalization required | Law 109/2025 Art. 51 |

### 5.4 Leave Rules

| Rule ID | Rule | Law Reference |
|---------|------|---------------|
| LR-01 | Annual leave: 12 days/year (after 12 months) | Labour Code 2019 Art. 113 |
| LR-02 | Annual leave increases 1 day per 5 years of service | Labour Code 2019 Art. 113 |
| LR-03 | Sick leave: paid at 75% of salary (with certificate) | Law 41/2024 Art. 25 |
| LR-04 | Maternity leave: 6 months (birth), 2 months (adoption) | Law 41/2024 Art. 39 |
| LR-05 | Paternity leave: 5-14 days (depending on birth type) | Law 41/2024 Art. 39 |
| LR-06 | Unused annual leave: paid out at termination | Labour Code 2019 Art. 114 |

---

## 6. Data Flow Diagrams

### 6.1 Payroll Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        DATA FLOW DIAGRAM                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  External Sources              GoTax System                       │
│  ─────────────────            ─────────────────                   │
│                                                                   │
│  ┌──────────────┐            ┌──────────────────┐                 │
│  │ Timekeeping  │───CSV───▶│ Timekeeping       │                 │
│  │ Device/App   │            │ Import Handler    │                 │
│  └──────────────┘            └────────┬─────────┘                 │
│                                       │                           │
│  ┌──────────────┐                     ▼                           │
│  │ Leave        │───API───▶┌──────────────────┐                 │
│  │ Requests     │            │ Timekeeping      │                 │
│  └──────────────┘            │ Module           │                 │
│                              └────────┬─────────┘                 │
│  ┌──────────────┐                     │                           │
│  │ Employee     │───DB────▶┌────────▼─────────┐                 │
│  │ Master Data  │            │ Payroll Engine   │                 │
│  └──────────────┘            │ (Calculation)    │                 │
│                              └────────┬─────────┘                 │
│  ┌──────────────┐                     │                           │
│  │ Salary       │───DB────▶│                  │                 │
│  │ Structure    │            │                  │                 │
│  └──────────────┘            │                  │                 │
│                              └────────┬─────────┘                 │
│                                       │                           │
│                    ┌──────────────────┼──────────────────┐        │
│                    │                  │                  │        │
│                    ▼                  ▼                  ▼        │
│  ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐  │
│  │ Payslips (PDF)   │ │ Journal Entries  │ │ Declarations     │  │
│  │ Email/Self-Svc   │ │ GL Integration   │ │ D02-TS, 05/KK    │  │
│  └──────────────────┘ └──────────────────┘ └──────────────────┘  │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 7. API Endpoints (Proposed)

### 7.1 Employee Payroll Info

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/payroll/employees/:id` | Get employee payroll info |
| PUT | `/api/v1/payroll/employees/:id` | Update employee payroll info |
| GET | `/api/v1/payroll/employees/:id/dependants` | List dependants |
| POST | `/api/v1/payroll/employees/:id/dependants` | Add dependant |
| PUT | `/api/v1/payroll/employees/:id/dependants/:did` | Update dependant |
| DELETE | `/api/v1/payroll/employees/:id/dependants/:did` | Remove dependant |

### 7.2 Salary Structure

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/payroll/components` | List salary components |
| POST | `/api/v1/payroll/components` | Create component |
| PUT | `/api/v1/payroll/components/:id` | Update component |
| GET | `/api/v1/payroll/templates` | List salary templates |
| POST | `/api/v1/payroll/templates` | Create template |
| PUT | `/api/v1/payroll/templates/:id` | Update template |

### 7.3 Timekeeping

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/payroll/timekeeping/import` | Import timekeeping CSV |
| GET | `/api/v1/payroll/timekeeping` | List timekeeping records |
| PUT | `/api/v1/payroll/timekeeping/:id` | Update record |
| DELETE | `/api/v1/payroll/timekeeping/:id` | Delete record |

### 7.4 Payroll Processing

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/payroll/periods` | List payroll periods |
| POST | `/api/v1/payroll/periods` | Create period |
| POST | `/api/v1/payroll/periods/:id/calculate` | Calculate payroll |
| POST | `/api/v1/payroll/periods/:id/submit` | Submit for review |
| POST | `/api/v1/payroll/periods/:id/approve` | Approve payroll |
| POST | `/api/v1/payroll/periods/:id/reject` | Reject payroll |
| GET | `/api/v1/payroll/periods/:id/runs` | List payroll runs |
| PUT | `/api/v1/payroll/runs/:id` | Adjust individual run |
| GET | `/api/v1/payroll/periods/:id/summary` | Get payroll summary |

### 7.5 Payslips

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/payroll/payslips` | List payslips |
| GET | `/api/v1/payroll/payslips/:id` | Get payslip detail |
| GET | `/api/v1/payroll/payslips/:id/pdf` | Download payslip PDF |
| POST | `/api/v1/payroll/payslips/:id/send` | Send payslip email |

### 7.6 Declarations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/payroll/declarations/d02-ts` | Generate D02-TS |
| GET | `/api/v1/payroll/declarations/05-kk-tncn` | Generate 05/KK-TNCN |
| GET | `/api/v1/payroll/declarations/tk3-ts` | Generate TK3-TS |
| GET | `/api/v1/payroll/declarations/:id/xml` | Download declaration XML |

### 7.7 Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/payroll/reports/summary` | Payroll summary report |
| GET | `/api/v1/payroll/reports/cost-by-department` | Cost allocation |
| GET | `/api/v1/payroll/reports/insurance-summary` | Insurance summary |
| GET | `/api/v1/payroll/reports/pit-summary` | PIT summary |
| GET | `/api/v1/payroll/reports/overtime` | Overtime report |
| GET | `/api/v1/payroll/reports/leave-balance` | Leave balance |

### 7.8 Configuration

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/payroll/config/rates` | Get current rates |
| PUT | `/api/v1/payroll/config/rates` | Update rates |
| GET | `/api/v1/payroll/config/holidays` | List holidays |
| POST | `/api/v1/payroll/config/holidays` | Add holiday |
| GET | `/api/v1/payroll/config/regions` | List regions with min wages |

---

## 8. Database Schema (PostgreSQL)

### 8.1 New Tables

```sql
-- Employee payroll info (extends employees table)
CREATE TABLE employee_payroll_info (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    contract_type VARCHAR(20) NOT NULL DEFAULT 'INDEFINITE',
    contract_start_date DATE,
    contract_end_date DATE,
    salary_type VARCHAR(20) NOT NULL DEFAULT 'TIME_BASED',
    base_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    salary_coefficient NUMERIC(5,2) DEFAULT 0,
    position_allowance NUMERIC(15,2) DEFAULT 0,
    responsibility_allowance NUMERIC(15,2) DEFAULT 0,
    seniority_allowance NUMERIC(15,2) DEFAULT 0,
    other_allowances NUMERIC(15,2) DEFAULT 0,
    insurance_base_salary NUMERIC(15,2) DEFAULT 0,
    region VARCHAR(5) NOT NULL DEFAULT 'I',
    is_foreign_employee BOOLEAN DEFAULT FALSE,
    is_trade_union_member BOOLEAN DEFAULT FALSE,
    is_high_tech_talent BOOLEAN DEFAULT FALSE,
    bank_account_no VARCHAR(50),
    bank_code VARCHAR(20),
    effective_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(employee_id, effective_date)
);

-- Dependants
CREATE TABLE dependants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    full_name VARCHAR(200) NOT NULL,
    relationship VARCHAR(20) NOT NULL,
    date_of_birth DATE,
    tax_code VARCHAR(20),
    is_disabled BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Salary components
CREATE TABLE salary_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(20) NOT NULL, -- INCOME, DEDUCTION, EMPLOYER_COST
    calculation VARCHAR(20) NOT NULL, -- FIXED, PERCENTAGE, FORMULA
    formula TEXT,
    default_value NUMERIC(15,2) DEFAULT 0,
    is_taxable BOOLEAN DEFAULT TRUE,
    is_insurable BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, code)
);

-- Payroll periods
CREATE TABLE payroll_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    year INT NOT NULL,
    month INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    prepared_by VARCHAR(100),
    prepared_at TIMESTAMPTZ,
    reviewed_by VARCHAR(100),
    reviewed_at TIMESTAMPTZ,
    approved_by VARCHAR(100),
    approved_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, year, month)
);

-- Payroll runs (per employee per period)
CREATE TABLE payroll_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_id UUID NOT NULL REFERENCES payroll_periods(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    company_id UUID NOT NULL,
    
    -- Input
    working_days INT DEFAULT 0,
    ot_hours NUMERIC(5,2) DEFAULT 0,
    night_shift_hours NUMERIC(5,2) DEFAULT 0,
    leave_days NUMERIC(5,2) DEFAULT 0,
    unpaid_leave_days NUMERIC(5,2) DEFAULT 0,
    absent_days INT DEFAULT 0,
    
    -- Income
    base_salary NUMERIC(15,2) DEFAULT 0,
    ot_pay NUMERIC(15,2) DEFAULT 0,
    night_shift_pay NUMERIC(15,2) DEFAULT 0,
    leave_pay NUMERIC(15,2) DEFAULT 0,
    holiday_pay NUMERIC(15,2) DEFAULT 0,
    allowances NUMERIC(15,2) DEFAULT 0,
    bonuses NUMERIC(15,2) DEFAULT 0,
    other_income NUMERIC(15,2) DEFAULT 0,
    gross_salary NUMERIC(15,2) DEFAULT 0,
    
    -- Employee deductions
    si_deduction NUMERIC(15,2) DEFAULT 0,
    hi_deduction NUMERIC(15,2) DEFAULT 0,
    ui_deduction NUMERIC(15,2) DEFAULT 0,
    trade_union_dues NUMERIC(15,2) DEFAULT 0,
    pit_amount NUMERIC(15,2) DEFAULT 0,
    other_deductions NUMERIC(15,2) DEFAULT 0,
    total_deductions NUMERIC(15,2) DEFAULT 0,
    
    -- Net
    net_pay NUMERIC(15,2) DEFAULT 0,
    
    -- Employer costs
    employer_si NUMERIC(15,2) DEFAULT 0,
    employer_hi NUMERIC(15,2) DEFAULT 0,
    employer_ui NUMERIC(15,2) DEFAULT 0,
    employer_trade_union NUMERIC(15,2) DEFAULT 0,
    total_employer_cost NUMERIC(15,2) DEFAULT 0,
    
    -- Insurance base
    insurance_base NUMERIC(15,2) DEFAULT 0,
    ui_base NUMERIC(15,2) DEFAULT 0,
    
    -- Status
    status VARCHAR(20) DEFAULT 'DRAFT',
    adjustment_reason TEXT,
    adjusted_by VARCHAR(100),
    adjusted_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(period_id, employee_id)
);

-- Timekeeping records
CREATE TABLE timekeeping_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    company_id UUID NOT NULL,
    date DATE NOT NULL,
    clock_in TIME,
    clock_out TIME,
    hours_worked NUMERIC(5,2) DEFAULT 0,
    ot_hours NUMERIC(5,2) DEFAULT 0,
    night_hours NUMERIC(5,2) DEFAULT 0,
    is_holiday BOOLEAN DEFAULT FALSE,
    is_rest_day BOOLEAN DEFAULT FALSE,
    leave_type VARCHAR(20),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(employee_id, date)
);

-- Leave requests
CREATE TABLE leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    company_id UUID NOT NULL,
    leave_type VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days NUMERIC(5,2) NOT NULL,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'PENDING',
    approved_by VARCHAR(100),
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Leave balances
CREATE TABLE leave_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    year INT NOT NULL,
    leave_type VARCHAR(20) NOT NULL,
    entitled NUMERIC(5,2) DEFAULT 0,
    used NUMERIC(5,2) DEFAULT 0,
    remaining NUMERIC(5,2) DEFAULT 0,
    carried_over NUMERIC(5,2) DEFAULT 0,
    UNIQUE(employee_id, year, leave_type)
);

-- Payslips
CREATE TABLE payslips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES payroll_runs(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    period_id UUID NOT NULL REFERENCES payroll_periods(id),
    payslip_no VARCHAR(50) NOT NULL,
    pdf_path VARCHAR(500),
    sent_at TIMESTAMPTZ,
    viewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Payroll configuration (rates, caps)
CREATE TABLE payroll_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    config_key VARCHAR(100) NOT NULL,
    config_value TEXT NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, config_key, effective_from)
);
```

---

## 9. Testing Strategy

### 9.1 Unit Tests (Calculation Engine)

| Test Case | Input | Expected Output |
|-----------|-------|-----------------|
| TC-01: Basic gross-to-net | Gross: 20,000,000, 0 dependants | Net: ~17,677,500 |
| TC-02: With dependants | Gross: 20,000,000, 2 dependants | Net: ~18,392,500 |
| TC-03: High salary (cap) | Gross: 60,000,000, 0 dependants | Insurance base capped at 50,600,000 |
| TC-04: Foreign employee | Gross: 30,000,000, foreign | No UI deduction |
| TC-05: OT exempt | Gross: 20M, OT: 2M | OT exempt from PIT |
| TC-06: Night shift exempt | Gross: 20M, Night: 1M | Night exempt from PIT |
| TC-07: Negative taxable | Gross: 16,000,000, 2 dependants | Taxable < 0, PIT = 0 |
| TC-08: Bracket 5 | Gross: 150,000,000 | PIT at 35% bracket |
| TC-09: Minimum wage check | Salary: 4,000,000, Region I | Flag: below minimum |
| TC-10: Trade union | Member, base: 50,600,000 | Dues: 253,000 (capped) |

### 9.2 Integration Tests

| Test Case | Description |
|-----------|-------------|
| IT-01 | Full payroll cycle: import → calculate → approve → GL posting |
| IT-02 | Payslip generation and PDF output |
| IT-03 | D02-TS declaration generation |
| IT-04 | 05/KK-TNCN quarterly declaration |
| IT-05 | Leave balance update after payroll |
| IT-06 | Period close prevents re-processing |

### 9.3 Compliance Tests

| Test Case | Description |
|-----------|-------------|
| CT-01 | PIT calculation matches official calculator (HAPRI) |
| CT-02 | Insurance rates match Law 41/2024 + Decree 161/2026 |
| CT-03 | Overtime limits enforced (40h/month, 200h/year) |
| CT-04 | Regional minimum wage compliance |
| CT-05 | Foreign employee rates correctly applied |

---

## 10. User Journey Maps

### 10.1 Payroll Clerk Journey

```
┌─────────────────────────────────────────────────────────────┐
│           PAYROLL CLERK — MONTHLY JOURNEY                     │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Week 1: Data Collection                                     │
│  ┌──────────────────────────────────────────────────────┐    │
│  │ ☐ Login to GoTax                                     │    │
│  │ ☐ Navigate to Payroll → Timekeeping                  │    │
│  │ ☐ Import timekeeping CSV from attendance system       │    │
│  │ ☐ Review imported data for errors                     │    │
│  │ ☐ Update employee changes (new hires, terminations)   │    │
│  │ ☐ Update salary changes (promotions, adjustments)     │    │
│  │ ☐ Collect pending leave requests                      │    │
│  └──────────────────────────────────────────────────────┘    │
│                                                               │
│  Week 2: Calculation                                         │
│  ┌──────────────────────────────────────────────────────┐    │
│  │ ☐ Navigate to Payroll → Periods                       │    │
│  │ ☐ Select month (e.g., July 2026)                      │    │
│  │ ☐ Click "Calculate Payroll"                           │    │
│  │ ☐ Review calculation results                          │    │
│  │ ☐ Check flagged employees (warnings)                  │    │
│  │ ☐ Make adjustments if needed                          │    │
│  │ ☐ Submit for review                                   │    │
│  └──────────────────────────────────────────────────────┘    │
│                                                               │
│  Week 3: Payment                                             │
│  ┌──────────────────────────────────────────────────────┐    │
│  │ ☐ After approval: generate bank salary file           │    │
│  │ ☐ Upload to bank portal                               │    │
│  │ ☐ Confirm payment completion                          │    │
│  │ ☐ Record payment date in system                       │    │
│  └──────────────────────────────────────────────────────┘    │
│                                                               │
│  Week 4: Declarations                                        │
│  ┌──────────────────────────────────────────────────────┐    │
│  │ ☐ Generate D02-TS (if new employees)                  │    │
│  │ ☐ Generate 05/KK-TNCN (quarterly)                     │    │
│  │ ☐ Submit declarations electronically                  │    │
│  │ ☐ Archive declaration receipts                        │    │
│  └──────────────────────────────────────────────────────┘    │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### 10.2 Employee Self-Service Journey

```
┌─────────────────────────────────────────────────────────────┐
│           EMPLOYEE — PAYSLIP JOURNEY                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────────────────────────────────────────────┐    │
│  │ ☐ Login to GoTax self-service portal                  │    │
│  │ ☐ Navigate to My Payslips                             │    │
│  │ ☐ View current month payslip                          │    │
│  │ ☐ Check gross salary, deductions, net pay             │    │
│  │ ☐ Download PDF for records                            │    │
│  │ ☐ View historical payslips                            │    │
│  │ ☐ Check leave balances                                │    │
│  │ ☐ Submit leave request (if needed)                    │    │
│  └──────────────────────────────────────────────────────┘    │
│                                                               │
│  Dispute Flow:                                                │
│  ┌──────────────────────────────────────────────────────┐    │
│  │ ☐ Click "Dispute" on payslip                          │    │
│  │ ☐ Enter reason for dispute                            │    │
│  │ ☐ Submit dispute                                      │    │
│  │ ☐ Receive notification when resolved                  │    │
│  └──────────────────────────────────────────────────────┘    │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## 11. Migration Strategy

### 11.1 Data Migration from External Systems

| Source | Migration Method | Priority |
|--------|-----------------|----------|
| Excel payroll data | CSV import template | HIGH |
| MISA AMIS | API export → CSV → import | MEDIUM |
| Fast Accounting | API export → CSV → import | MEDIUM |
| Manual records | Web form entry | HIGH |

### 11.2 Historical Data

- Import last 12 months of payslips for reference
- Import current leave balances
- Import current insurance registration data
- Import employee payroll info (salary, allowances, dependants)

---

## 12. Appendix

### 12.1 Glossary

| Term | Vietnamese | Definition |
|------|------------|------------|
| SI | BHXH | Social Insurance |
| HI | BHYT | Health Insurance |
| UI | BHTN | Unemployment Insurance |
| PIT | TNCN | Personal Income Tax |
| OT | Làm thêm giờ | Overtime |
| Gross | Tổng thu nhập | Total income before deductions |
| Net | Thực lĩnh | Amount received after deductions |
| D02-TS | Danh sách tham gia BHXH | SI registration list |
| TK3-TS | Tờ khai đơn vị | Employer SI declaration |
| 05/KK-TNCN | Tờ khai thuế TNCN | PIT declaration form |

### 12.2 Legal References

1. Luật Bảo hiểm xã hội số 41/2024/QH15
2. Luật Việc làm số 74/2025/QH15
3. Luật Bảo hiểm y tế số 51/2024/QH15
4. Luật Thuế thu nhập cá nhân số 109/2025/QH15
5. Nghị định 161/2026/NĐ-CP (mức lương cơ sở)
6. Nghị định 293/2025/NĐ-CP (lương tối thiểu vùng)
7. Nghị định 253/2026/NĐ-CP (hướng dẫn Luật TNCN)
8. Nghị định 158/2025/NĐ-CP (hướng dẫn Luật BHXH)
9. Quyết định 1109/QĐ-BTC (mẫu 05/KK-TNCN)
10. Quyết định 366/QĐ-BHXH (biểu mẫu BHXH 2026)
11. Thông tư 99/2025/TT-BTC (chế độ kế toán doanh nghiệp)
12. Luật Lao động 2019 (ghovertime, phép nghỉ)
