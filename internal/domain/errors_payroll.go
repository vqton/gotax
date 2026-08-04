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
)
