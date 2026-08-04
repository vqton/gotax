package domain

import "math"

// ─── Insurance Constants ────────────────────────────────────────

const (
	SIEmployeeRate  = 0.08  // 8%
	HIEmployeeRate  = 0.015 // 1.5%
	UIEmployeeRate  = 0.01  // 1%
	SIEmployerRate  = 0.175 // 17.5%
	HIEmployerRate  = 0.03  // 3%
	UIEmployerRate  = 0.01  // 1%
	TradeUnionRate  = 0.01  // 1%
	BaseSalaryCapSI = 20    // SI/HI cap = 20 × base salary
	BaseSalaryCapUI = 20    // UI cap = 20 × regional min wage
	TradeUnionCap   = 253_000.0

	BaseSalary2026     = 2_530_000.0 // Decree 161/2026
	RegionalMinWageI   = 5_310_000.0  // Decree 293/2025
	RegionalMinWageII  = 4_680_000.0
	RegionalMinWageIII = 4_160_000.0
	RegionalMinWageIV  = 3_990_000.0

	WorkingHoursPerDay = 8
)

// OT / leave multipliers
const (
	OTWeekdayRate     = 1.5 // 150%
	OTWeekendRate     = 2.0 // 200%
	OTHolidayRate     = 3.0 // 300%
	NightShiftPremium = 0.3 // 30%
	HolidayPayRate    = 3.0 // 300%
)

// PITBracket defines a progressive tax bracket.
type PITBracket struct {
	UpperLimit float64
	Rate       float64
	QuickDed   float64
}

// PITBrackets — monthly progressive tax (Decision 253/2026).
// 7 brackets: 5% → 10% → 15% → 20% → 25% → 30% → 35%
var PITBrackets = []PITBracket{
	{5_000_000, 0.05, 0},
	{10_000_000, 0.10, 250_000},
	{18_000_000, 0.15, 750_000},
	{32_000_000, 0.20, 1_650_000},
	{52_000_000, 0.25, 3_250_000},
	{80_000_000, 0.30, 5_850_000},
	{999_999_999_999, 0.35, 9_850_000},
}

// ─── Helpers ────────────────────────────────────────────────────

func round2(v float64) float64 { return math.Round(v*100) / 100 }
func round0(v float64) float64 { return math.Round(v) }
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// ─── Insurance (Employee) ──────────────────────────────────────

func CalcEmployeeSI(base float64) float64 {
	cap := BaseSalaryCapSI * BaseSalary2026
	return round0(minFloat(base, cap) * SIEmployeeRate)
}

func CalcEmployeeHI(base float64) float64 {
	cap := BaseSalaryCapSI * BaseSalary2026
	return round0(minFloat(base, cap) * HIEmployeeRate)
}

func CalcEmployeeUI(base float64, isForeign bool) float64 {
	if isForeign {
		return 0
	}
	cap := BaseSalaryCapUI * RegionalMinWageI
	return round0(minFloat(base, cap) * UIEmployeeRate)
}

// ─── Insurance (Employer) ──────────────────────────────────────

func CalcEmployerSI(base float64) float64 {
	cap := BaseSalaryCapSI * BaseSalary2026
	return round0(minFloat(base, cap) * SIEmployerRate)
}

func CalcEmployerHI(base float64) float64 {
	cap := BaseSalaryCapSI * BaseSalary2026
	return round0(minFloat(base, cap) * HIEmployerRate)
}

func CalcEmployerUI(base float64, isForeign bool) float64 {
	if isForeign {
		return 0
	}
	cap := BaseSalaryCapUI * RegionalMinWageI
	return round0(minFloat(base, cap) * UIEmployerRate)
}

// ─── Trade Union ────────────────────────────────────────────────

func CalcTradeUnion(base float64, isForeign bool) float64 {
	if isForeign {
		return 0
	}
	return round0(minFloat(base*TradeUnionRate, TradeUnionCap))
}

// ─── PIT ────────────────────────────────────────────────────────

func CalcPIT(taxableIncome float64) float64 {
	if taxableIncome <= 0 {
		return 0
	}
	for _, b := range PITBrackets {
		if taxableIncome <= b.UpperLimit {
			return round0(taxableIncome*b.Rate - b.QuickDed)
		}
	}
	return round0(taxableIncome*0.35 - 9_850_000)
}

// ─── Salary Helpers ─────────────────────────────────────────────

func CalcDailySalary(monthlySalary float64, workingDays int) float64 {
	if workingDays == 0 {
		return 0
	}
	return round2(monthlySalary / float64(workingDays))
}

func CalcHourlySalary(dailySalary float64) float64 {
	return round2(dailySalary / WorkingHoursPerDay)
}

func CalcOTPay(hourlyRate, otHours, multiplier float64) float64 {
	return round0(hourlyRate * otHours * multiplier)
}

func CalcNightShiftPay(hourlyRate, nightHours float64) float64 {
	return round0(hourlyRate * nightHours * (1 + NightShiftPremium))
}

func CalcLeavePay(dailySalary, days float64) float64 {
	return round0(dailySalary * days)
}

func CalcHolidayPay(dailySalary float64) float64 {
	return round0(dailySalary * HolidayPayRate)
}

// ─── Full Employee Payroll ──────────────────────────────────────

// CalculateEmployeePayroll computes gross → deductions → net for one employee.
func CalculateEmployeePayroll(info EmployeePayrollInfo, period *PayrollPeriod, timekeeping []TimekeepingRecord) PayrollRun {
	run := PayrollRun{}

	// Count working days from timekeeping
	workingDays := 0
	otHours := 0.0
	nightShiftHours := 0.0
	leaveDays := 0.0
	absentDays := 0

	for _, tk := range timekeeping {
		if tk.IsRestDay || tk.IsHoliday {
			continue
		}
		if tk.LeaveType != "" {
			leaveDays++
			continue
		}
		if tk.HoursWorked == 0 {
			absentDays++
			continue
		}
		workingDays++
		otHours += tk.OTHours
		nightShiftHours += tk.NightHours
	}

	if workingDays == 0 {
		workingDays = 26
	}

	run.WorkingDays = workingDays
	run.OTHours = otHours
	run.NightShiftHours = nightShiftHours
	run.LeaveDays = leaveDays
	run.AbsentDays = absentDays
	run.BaseSalary = info.BaseSalary

	dailySalary := CalcDailySalary(info.BaseSalary, workingDays)
	hourlyRate := CalcHourlySalary(dailySalary)

	run.OTPay = CalcOTPay(hourlyRate, otHours, OTWeekdayRate)
	run.NightShiftPay = CalcNightShiftPay(hourlyRate, nightShiftHours)
	run.LeavePay = CalcLeavePay(dailySalary, leaveDays)
	run.Allowances = info.PositionAllowance + info.ResponsibilityAllowance +
		info.SeniorityAllowance + info.OtherAllowances

	run.GrossSalary = run.BaseSalary + run.OTPay + run.NightShiftPay +
		run.LeavePay + run.HolidayPay + run.Allowances + run.Bonuses + run.OtherIncome

	// Insurance base
	run.InsuranceBase = info.InsuranceBaseSalary
	if run.InsuranceBase == 0 {
		run.InsuranceBase = info.BaseSalary
	}
	run.UIBase = info.BaseSalary

	// Employee deductions
	run.SIDeduction = CalcEmployeeSI(run.InsuranceBase)
	run.HIDeduction = CalcEmployeeHI(run.InsuranceBase)
	run.UIDeduction = CalcEmployeeUI(run.UIBase, info.IsForeignEmployee)
	run.TradeUnionDues = CalcTradeUnion(info.BaseSalary, info.IsForeignEmployee)

	taxable := run.GrossSalary - run.SIDeduction - run.HIDeduction - run.UIDeduction - run.TradeUnionDues
	run.PITAmount = CalcPIT(taxable)

	run.TotalDeductions = run.SIDeduction + run.HIDeduction + run.UIDeduction +
		run.TradeUnionDues + run.PITAmount + run.OtherDeductions
	run.NetPay = run.GrossSalary - run.TotalDeductions

	// Employer costs
	run.EmployerSI = CalcEmployerSI(run.InsuranceBase)
	run.EmployerHI = CalcEmployerHI(run.InsuranceBase)
	run.EmployerUI = CalcEmployerUI(run.UIBase, info.IsForeignEmployee)
	run.EmployerTradeUnion = CalcTradeUnion(info.BaseSalary, info.IsForeignEmployee)
	run.TotalEmployerCost = run.EmployerSI + run.EmployerHI + run.EmployerUI + run.EmployerTradeUnion

	run.Status = "CALCULATED"
	return run
}
