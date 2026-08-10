package domain

import "errors"

// ─── Payroll Errors ─────────────────────────────────────────────

var (
	ErrPayrollPeriodExists   = errors.New("payroll period already exists")
	ErrPayrollPeriodNotFound = errors.New("payroll period not found")
	ErrPayrollPeriodNotDraft = errors.New("period is not in draft status")
	ErrPayrollRunNotFound    = errors.New("payroll run not found")
	ErrPayrollRunNotDraft    = errors.New("run is not in draft status")
	ErrPayrollRunAlreadyApproved = errors.New("run already approved")
	ErrPayrollEmployeeNotFound   = errors.New("employee payroll info not found")
	ErrPayrollDependantNotFound  = errors.New("dependant not found")
	ErrPayrollLeaveNotFound      = errors.New("leave request not found")
	ErrPayrollLeaveAlreadyProcessed = errors.New("leave request already processed")
	ErrPayrollPayslipNotFound    = errors.New("payslip not found")
	ErrPayrollInvalidPeriod      = errors.New("invalid payroll period")
	ErrPayrollNoEmployees        = errors.New("no employees to process")
	ErrPayrollConfigNotFound     = errors.New("payroll config not found")
	ErrPayrollComponentNotFound  = errors.New("salary component not found")
	ErrPayrollComponentExists    = errors.New("salary component code already exists")
	ErrPayrollTemplateNotFound   = errors.New("salary template not found")
	ErrPayrollTemplateExists     = errors.New("salary template name already exists")
	ErrPayrollHolidayNotFound    = errors.New("holiday not found")
	ErrPayrollHolidayExists      = errors.New("holiday already exists")
	ErrPayrollTimekeepingNotFound = errors.New("timekeeping record not found")
	ErrPayrollSalaryGradeNotFound = errors.New("salary grade not found")
	ErrPayrollSalaryGradeExists   = errors.New("salary grade code already exists")
	ErrPayrollSalaryScaleNotFound = errors.New("salary scale not found")
	ErrPayrollSalaryScaleExists   = errors.New("salary scale code already exists")
	ErrPayrollESGNotFound         = errors.New("employee salary grade not found")
	ErrPayrollCompanyRequired     = errors.New("company_id is required")
	ErrPayrollLinesRequired       = errors.New("at least one allocation item is required")
)
