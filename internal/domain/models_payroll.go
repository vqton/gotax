package domain

import "time"

// ─── Payroll Enums ──────────────────────────────────────────────

type PayrollStatus string

const (
	PayrollDraft      PayrollStatus = "DRAFT"
	PayrollProcessing PayrollStatus = "PROCESSING"
	PayrollReviewing  PayrollStatus = "REVIEWING"
	PayrollApproved   PayrollStatus = "APPROVED"
	PayrollPaid       PayrollStatus = "PAID"
	PayrollClosed     PayrollStatus = "CLOSED"
)

type ContractType string

const (
	ContractIndefinite ContractType = "INDEFINITE"
	ContractDefinite   ContractType = "DEFINITE"
	ContractProbation  ContractType = "PROBATION"
)

type SalaryType string

const (
	SalaryTimeBased   SalaryType = "TIME_BASED"
	SalaryPieceRate   SalaryType = "PIECE_RATE"
	SalaryCommission  SalaryType = "COMMISSION"
	SalaryCoefficient SalaryType = "COEFFICIENT"
)

type LeaveType string

const (
	LeaveAnnual       LeaveType = "ANNUAL"
	LeaveSick         LeaveType = "SICK"
	LeaveMaternity    LeaveType = "MATERNITY"
	LeavePaternity    LeaveType = "PATERNITY"
	LeaveUnpaid       LeaveType = "UNPAID"
	LeaveCompensatory LeaveType = "COMPENSATORY"
	LeaveMarriage     LeaveType = "MARRIAGE"
	LeaveBereavement  LeaveType = "BEREAVEMENT"
)

type LeaveRequestStatus string

const (
	LeavePending  LeaveRequestStatus = "PENDING"
	LeaveApproved LeaveRequestStatus = "APPROVED"
	LeaveRejected LeaveRequestStatus = "REJECTED"
)

type InsuranceRegion string

const (
	RegionI   InsuranceRegion = "I"
	RegionII  InsuranceRegion = "II"
	RegionIII InsuranceRegion = "III"
	RegionIV  InsuranceRegion = "IV"
)

// ─── Payroll Domain Models ─────────────────────────────────────

// EmployeePayrollInfo extends Employee with payroll-specific data.
type EmployeePayrollInfo struct {
	ID                    string          `json:"id"`
	EmployeeID            string          `json:"employee_id"`
	ContractType          ContractType    `json:"contract_type"`
	ContractStartDate     string          `json:"contract_start_date"`
	ContractEndDate       string          `json:"contract_end_date,omitempty"`
	SalaryType            SalaryType      `json:"salary_type"`
	BaseSalary            float64         `json:"base_salary"`
	SalaryCoefficient     float64         `json:"salary_coefficient"`
	PositionAllowance     float64         `json:"position_allowance"`
	ResponsibilityAllowance float64       `json:"responsibility_allowance"`
	SeniorityAllowance    float64         `json:"seniority_allowance"`
	OtherAllowances       float64         `json:"other_allowances"`
	InsuranceBaseSalary   float64         `json:"insurance_base_salary"`
	Region                InsuranceRegion `json:"region"`
	IsForeignEmployee     bool            `json:"is_foreign_employee"`
	IsTradeUnionMember    bool            `json:"is_trade_union_member"`
	IsHighTechTalent      bool            `json:"is_high_tech_talent"`
	BankAccountNo         string          `json:"bank_account_no"`
	BankCode              string          `json:"bank_code"`
	EffectiveDate         string          `json:"effective_date"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

// Dependant represents an employee's dependant for PIT deduction.
type Dependant struct {
	ID           string    `json:"id"`
	EmployeeID   string    `json:"employee_id"`
	FullName     string    `json:"full_name"`
	Relationship string    `json:"relationship"`
	DateOfBirth  string    `json:"date_of_birth"`
	TaxCode      string    `json:"tax_code,omitempty"`
	IsDisabled   bool      `json:"is_disabled"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// PayrollPeriod represents a payroll processing period.
type PayrollPeriod struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"company_id"`
	Year        int            `json:"year"`
	Month       int            `json:"month"`
	Status      PayrollStatus  `json:"status"`
	PreparedBy  string         `json:"prepared_by,omitempty"`
	PreparedAt  *time.Time     `json:"prepared_at,omitempty"`
	ReviewedBy  string         `json:"reviewed_by,omitempty"`
	ReviewedAt  *time.Time     `json:"reviewed_at,omitempty"`
	ApprovedBy  string         `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time     `json:"approved_at,omitempty"`
	PaidAt      *time.Time     `json:"paid_at,omitempty"`
	ClosedAt    *time.Time     `json:"closed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// PayrollRun represents a single employee's payroll calculation for a period.
type PayrollRun struct {
	ID                string    `json:"id"`
	PeriodID          string    `json:"period_id"`
	EmployeeID        string    `json:"employee_id"`
	CompanyID         string    `json:"company_id"`

	// Input
	WorkingDays     int     `json:"working_days"`
	OTHours         float64 `json:"ot_hours"`
	NightShiftHours float64 `json:"night_shift_hours"`
	LeaveDays       float64 `json:"leave_days"`
	UnpaidLeaveDays float64 `json:"unpaid_leave_days"`
	AbsentDays      int     `json:"absent_days"`

	// Income
	BaseSalary      float64 `json:"base_salary"`
	OTPay           float64 `json:"ot_pay"`
	NightShiftPay   float64 `json:"night_shift_pay"`
	LeavePay        float64 `json:"leave_pay"`
	HolidayPay      float64 `json:"holiday_pay"`
	Allowances      float64 `json:"allowances"`
	Bonuses         float64 `json:"bonuses"`
	OtherIncome     float64 `json:"other_income"`
	GrossSalary     float64 `json:"gross_salary"`

	// Employee deductions
	SIDeduction     float64 `json:"si_deduction"`
	HIDeduction     float64 `json:"hi_deduction"`
	UIDeduction     float64 `json:"ui_deduction"`
	TradeUnionDues  float64 `json:"trade_union_dues"`
	PITAmount       float64 `json:"pit_amount"`
	OtherDeductions float64 `json:"other_deductions"`
	TotalDeductions float64 `json:"total_deductions"`

	// Net pay
	NetPay float64 `json:"net_pay"`

	// Employer costs
	EmployerSI          float64 `json:"employer_si"`
	EmployerHI          float64 `json:"employer_hi"`
	EmployerUI          float64 `json:"employer_ui"`
	EmployerTradeUnion  float64 `json:"employer_trade_union"`
	TotalEmployerCost   float64 `json:"total_employer_cost"`

	// Insurance base
	InsuranceBase float64 `json:"insurance_base"`
	UIBase        float64 `json:"ui_base"`

	// Status
	Status        string     `json:"status"`
	AdjustmentReason string  `json:"adjustment_reason,omitempty"`
	AdjustedBy    string     `json:"adjusted_by,omitempty"`
	AdjustedAt    *time.Time `json:"adjusted_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// TimekeepingRecord represents daily attendance data.
type TimekeepingRecord struct {
	ID          string    `json:"id"`
	EmployeeID  string    `json:"employee_id"`
	CompanyID   string    `json:"company_id"`
	Date        string    `json:"date"`
	ClockIn     string    `json:"clock_in,omitempty"`
	ClockOut    string    `json:"clock_out,omitempty"`
	HoursWorked float64   `json:"hours_worked"`
	OTHours     float64   `json:"ot_hours"`
	NightHours  float64   `json:"night_hours"`
	IsHoliday   bool      `json:"is_holiday"`
	IsRestDay   bool      `json:"is_rest_day"`
	LeaveType   LeaveType `json:"leave_type,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// LeaveRequest represents an employee leave request.
type LeaveRequest struct {
	ID         string            `json:"id"`
	EmployeeID string            `json:"employee_id"`
	CompanyID  string            `json:"company_id"`
	LeaveType  LeaveType         `json:"leave_type"`
	StartDate  string            `json:"start_date"`
	EndDate    string            `json:"end_date"`
	Days       float64           `json:"days"`
	Reason     string            `json:"reason"`
	Status     LeaveRequestStatus `json:"status"`
	ApprovedBy string            `json:"approved_by,omitempty"`
	ApprovedAt *time.Time        `json:"approved_at,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

// LeaveBalance tracks employee leave entitlements.
type LeaveBalance struct {
	ID          string    `json:"id"`
	EmployeeID  string    `json:"employee_id"`
	Year        int       `json:"year"`
	LeaveType   LeaveType `json:"leave_type"`
	Entitled    float64   `json:"entitled"`
	Used        float64   `json:"used"`
	Remaining   float64   `json:"remaining"`
	CarriedOver float64   `json:"carried_over"`
}

// Payslip represents an employee payslip.
type Payslip struct {
	ID         string    `json:"id"`
	RunID      string    `json:"run_id"`
	EmployeeID string    `json:"employee_id"`
	PeriodID   string    `json:"period_id"`
	PayslipNo  string    `json:"payslip_no"`
	PDFPath    string    `json:"pdf_path,omitempty"`
	SentAt     *time.Time `json:"sent_at,omitempty"`
	ViewedAt   *time.Time `json:"viewed_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// PayrollConfig stores configurable payroll parameters.
type PayrollConfig struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	ConfigKey   string    `json:"config_key"`
	ConfigValue string    `json:"config_value"`
	EffectiveFrom string  `json:"effective_from"`
	EffectiveTo   string  `json:"effective_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// PayrollSummary aggregates payroll data for a period.
type PayrollSummary struct {
	PeriodID          string  `json:"period_id"`
	EmployeeCount     int     `json:"employee_count"`
	TotalGross        float64 `json:"total_gross"`
	TotalDeductions   float64 `json:"total_deductions"`
	TotalNetPay       float64 `json:"total_net_pay"`
	TotalEmployerCost float64 `json:"total_employer_cost"`
	TotalSI           float64 `json:"total_si"`
	TotalHI           float64 `json:"total_hi"`
	TotalUI           float64 `json:"total_ui"`
	TotalPIT          float64 `json:"total_pit"`
}

// SalaryComponent defines a pay/deduction component.
type SalaryComponent struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`       // INCOME, DEDUCTION, EMPLOYER_COST
	Calculation  string    `json:"calculation"` // FIXED, PERCENTAGE, FORMULA
	Formula      string    `json:"formula,omitempty"`
	DefaultValue float64   `json:"default_value"`
	IsTaxable    bool      `json:"is_taxable"`
	IsInsurable  bool      `json:"is_insurable"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// SalaryTemplate groups components for a role/department.
type SalaryTemplate struct {
	ID          string                    `json:"id"`
	CompanyID   string                    `json:"company_id"`
	Name        string                    `json:"name"`
	Components  []SalaryTemplateComponent `json:"components,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
}

// SalaryTemplateComponent links a template to a component with defaults.
type SalaryTemplateComponent struct {
	ID               string  `json:"id"`
	SalaryTemplateID string  `json:"salary_template_id"`
	ComponentID      string  `json:"component_id"`
	DefaultValue     float64 `json:"default_value"`
	Formula          string  `json:"formula,omitempty"`
	Order            int     `json:"order"`
}

// InsuranceSummary aggregates insurance data for a period.
type InsuranceSummary struct {
	PeriodID          string  `json:"period_id"`
	EmployeeCount     int     `json:"employee_count"`
	TotalEmployeeSI   float64 `json:"total_employee_si"`
	TotalEmployeeHI   float64 `json:"total_employee_hi"`
	TotalEmployeeUI   float64 `json:"total_employee_ui"`
	TotalEmployerSI   float64 `json:"total_employer_si"`
	TotalEmployerHI   float64 `json:"total_employer_hi"`
	TotalEmployerUI   float64 `json:"total_employer_ui"`
	TotalSI           float64 `json:"total_si"`
	TotalHI           float64 `json:"total_hi"`
	TotalUI           float64 `json:"total_ui"`
}

// PITSummary aggregates PIT data for a period.
type PITSummary struct {
	PeriodID            string  `json:"period_id"`
	EmployeeCount       int     `json:"employee_count"`
	EmployeesWithPIT    int     `json:"employees_with_pit"`
	TotalPIT            float64 `json:"total_pit"`
	TotalTaxableIncome  float64 `json:"total_taxable_income"`
}

// OvertimeSummary aggregates overtime data for a period.
type OvertimeSummary struct {
	PeriodID        string  `json:"period_id"`
	EmployeesWithOT int     `json:"employees_with_ot"`
	TotalOTHours    float64 `json:"total_ot_hours"`
	TotalOTPay      float64 `json:"total_ot_pay"`
	TotalNightHours float64 `json:"total_night_hours"`
	TotalNightPay   float64 `json:"total_night_pay"`
}

// LeaveBalanceReport shows leave balance per employee per type.
type LeaveBalanceReport struct {
	EmployeeID string  `json:"employee_id"`
	Year       int     `json:"year"`
	LeaveType  string  `json:"leave_type"`
	Entitled   float64 `json:"entitled"`
	Used       float64 `json:"used"`
	Remaining  float64 `json:"remaining"`
}

// PayrollHoliday represents a public holiday.
type PayrollHoliday struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Name      string    `json:"name"`
	Date      string    `json:"date"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
}

// ─── Severance Pay ─────────────────────────────────────────────

// TerminationReason classifies why the contract ended.
type TerminationReason string

const (
	TerminationRedundancy   TerminationReason = "REDUNDANCY"
	TerminationRestructure  TerminationReason = "RESTRUCTURE"
	TerminationPerformance  TerminationReason = "PERFORMANCE"
	TerminationExpiration   TerminationReason = "EXPIRATION"
	TerminationMutual       TerminationReason = "MUTUAL"
	TerminationHealth       TerminationReason = "HEALTH"
	TerminationRetirement   TerminationReason = "RETIREMENT"
	TerminationGrossMisconduct TerminationReason = "GROSS_MISCONDUCT"
)

// SeveranceInput contains data for severance pay calculation.
type SeveranceInput struct {
	AvgSalary6Months float64          `json:"avg_salary_6_months"` // Average salary of last 6 months
	YearsOfService   float64          `json:"years_of_service"`    // Total years (can be fractional)
	Reason           TerminationReason `json:"reason"`
}

// SeveranceResult contains the calculated severance pay.
type SeveranceResult struct {
	GrossSeverance   float64 `json:"gross_severance"`   // 0.5 × avg × years
	PITAmount        float64 `json:"pit_amount"`        // PIT on severance
	NetSeverance     float64 `json:"net_severance"`     // After PIT
	YearsOfService   float64 `json:"years_of_service"`
	Reason           TerminationReason `json:"reason"`
}
