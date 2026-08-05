package domain

import "math"

// ─── Insurance Constants ────────────────────────────────────────

const (
	SIEmployeeRate = 0.08  // 8%
	HIEmployeeRate = 0.015 // 1.5%
	UIEmployeeRate = 0.01  // 1%
	SIEmployerRate = 0.175 // 17.5%
	HIEmployerRate = 0.03  // 3%
	UIEmployerRate = 0.01  // 1%
	TradeUnionRate = 0.01  // 1%

	BaseSalaryCapSI = 20 // SI/HI cap = 20 × base salary
	BaseSalaryCapUI = 20 // UI cap = 20 × regional min wage

	WorkingHoursPerDay = 8
)

// ─── Date-Ranged Salary Constants (Decree 161/2026, Decree 293/2025) ──

// ReferenceLevel (mức tham chiếu) used for SI/HI cap from 01 Jan 2026.
const ReferenceLevel2026H1 = 2_340_000.0

// BaseSalary (mức lương cơ sở) from 01 Jul 2026 per Decree 161/2026/NĐ-CP.
const BaseSalary2026H2 = 2_530_000.0

// GetBaseSalary returns the applicable base salary for a given month/year.
func GetBaseSalary(month, year int) float64 {
	if year > 2026 || (year == 2026 && month >= 7) {
		return BaseSalary2026H2
	}
	return ReferenceLevel2026H1
}

// Regional minimum wages from Decree 293/2025/NĐ-CP (effective 01 Jan 2026).
const (
	RegionalMinWageI   = 5_310_000.0
	RegionalMinWageII  = 4_730_000.0
	RegionalMinWageIII = 4_140_000.0
	RegionalMinWageIV  = 3_700_000.0
)

// GetSIHICap returns the SI/HI contribution cap for a given month/year.
func GetSIHICap(month, year int) float64 {
	return BaseSalaryCapSI * GetBaseSalary(month, year)
}

// GetTradeUnionCap returns the trade union dues cap for a given month/year.
func GetTradeUnionCap(month, year int) float64 {
	return 0.01 * GetSIHICap(month, year)
}

// ─── OT / leave multipliers ────────────────────────────────────

const (
	OTWeekdayRate     = 1.5 // 150%
	OTWeekendRate     = 2.0 // 200%
	OTHolidayRate     = 3.0 // 300%
	NightShiftPremium = 0.3 // 30%
	HolidayPayRate    = 3.0 // 300%
)

// ─── PIT Brackets ──────────────────────────────────────────────

// PITBracket defines a progressive tax bracket.
type PITBracket struct {
	UpperLimit float64
	Rate       float64
	QuickDed   float64
}

// PITBracketsOld — 7-bracket schedule (Law 103/2016/QH13, effective 01 Jul 2020).
// Applies from 01 Jan to 30 Jun 2026.
var PITBracketsOld = []PITBracket{
	{5_000_000, 0.05, 0},
	{10_000_000, 0.10, 250_000},
	{18_000_000, 0.15, 750_000},
	{32_000_000, 0.20, 1_650_000},
	{52_000_000, 0.25, 3_250_000},
	{80_000_000, 0.30, 5_850_000},
	{999_999_999_999, 0.35, 9_850_000},
}

// PITBracketsNew — 5-bracket schedule (Law 109/2025/QH15, effective 01 Jul 2026).
// Applies from 01 Jul 2026 onward.
var PITBracketsNew = []PITBracket{
	{10_000_000, 0.05, 0},
	{30_000_000, 0.10, 500_000},
	{60_000_000, 0.20, 3_500_000},
	{100_000_000, 0.30, 9_500_000},
	{999_999_999_999, 0.35, 14_500_000},
}

// GetPITBrackets returns the applicable brackets for a given month/year.
func GetPITBrackets(month, year int) []PITBracket {
	if year > 2026 || (year == 2026 && month >= 7) {
		return PITBracketsNew
	}
	return PITBracketsOld
}

// PITBrackets is an alias for backward compatibility — uses old 7-bracket schedule.
// Deprecated: use GetPITBrackets(month, year) instead.
var PITBrackets = PITBracketsOld

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

func CalcEmployeeSI(base float64, month, year int) float64 {
	cap := GetSIHICap(month, year)
	return round0(minFloat(base, cap) * SIEmployeeRate)
}

func CalcEmployeeHI(base float64, month, year int) float64 {
	cap := GetSIHICap(month, year)
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

func CalcEmployerSI(base float64, month, year int) float64 {
	cap := GetSIHICap(month, year)
	return round0(minFloat(base, cap) * SIEmployerRate)
}

func CalcEmployerHI(base float64, month, year int) float64 {
	cap := GetSIHICap(month, year)
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

func CalcTradeUnion(base float64, isForeign bool, month, year int) float64 {
	if isForeign {
		return 0
	}
	cap := GetTradeUnionCap(month, year)
	return round0(minFloat(base*TradeUnionRate, cap))
}

// ─── PIT ────────────────────────────────────────────────────────

const (
	// PersonalDeductionMonthly is the monthly personal deduction for PIT (Resolution 110/2025/UBTVQH15, effective 01/01/2026).
	PersonalDeductionMonthly = 15_500_000.0
	// DependantDeductionMonthly is the monthly dependent deduction for PIT (Resolution 110/2025/UBTVQH15, effective 01/01/2026).
	DependantDeductionMonthly = 6_200_000.0
)

func CalcPIT(taxableIncome float64, month, year int) float64 {
	if taxableIncome <= 0 {
		return 0
	}
	brackets := GetPITBrackets(month, year)
	for _, b := range brackets {
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

	// Period date for date-ranged calculations
	month := period.Month
	year := period.Year

	// Employee deductions
	run.SIDeduction = CalcEmployeeSI(run.InsuranceBase, month, year)
	run.HIDeduction = CalcEmployeeHI(run.InsuranceBase, month, year)
	run.UIDeduction = CalcEmployeeUI(run.UIBase, info.IsForeignEmployee)
	run.TradeUnionDues = CalcTradeUnion(info.BaseSalary, info.IsForeignEmployee, month, year)

	taxable := run.GrossSalary - run.SIDeduction - run.HIDeduction - run.UIDeduction - run.TradeUnionDues - PersonalDeductionMonthly
	run.PITAmount = CalcPIT(taxable, month, year)

	run.TotalDeductions = run.SIDeduction + run.HIDeduction + run.UIDeduction +
		run.TradeUnionDues + run.PITAmount + run.OtherDeductions
	run.NetPay = run.GrossSalary - run.TotalDeductions

	// Employer costs
	run.EmployerSI = CalcEmployerSI(run.InsuranceBase, month, year)
	run.EmployerHI = CalcEmployerHI(run.InsuranceBase, month, year)
	run.EmployerUI = CalcEmployerUI(run.UIBase, info.IsForeignEmployee)
	run.EmployerTradeUnion = CalcTradeUnion(info.BaseSalary, info.IsForeignEmployee, month, year)
	run.TotalEmployerCost = run.EmployerSI + run.EmployerHI + run.EmployerUI + run.EmployerTradeUnion

	run.Status = "CALCULATED"
	return run
}
