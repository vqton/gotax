package domain

import (
	"fmt"
	"math"
	"time"
)

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

// ─── Payroll Journal Entry Generation ───────────────────────────

// PayrollJEInput holds data for generating payroll journal entries.
type PayrollJEInput struct {
	PeriodID    string
	CompanyID   string
	PeriodMonth int
	PeriodYear  int
	Runs        []PayrollRun
}

// PayrollJournalEntries holds the generated journal entries for a payroll period.
type PayrollJournalEntries struct {
	Entries []*JournalEntry
}

// GeneratePayrollJournalEntries creates journal entries from approved payroll runs.
// Account codes per Circular 99/2025/TT-BTC:
//
//	6421 — Salary expense
//	6422 — Overtime expense
//	6423 — Allowance expense
//	6424 — Insurance expense (employer)
//	3331 — Salary payable
//	3332 — SI payable (employer)
//	3333 — HI payable (employer)
//	3334 — UI payable (employer)
//	3335 — SI payable (employee)
//	3336 — HI payable (employee)
//	3337 — UI payable (employee)
//	3338 — PIT payable
//	3339 — Trade union payable
func GeneratePayrollJournalEntries(input PayrollJEInput) PayrollJournalEntries {
	if len(input.Runs) == 0 {
		return PayrollJournalEntries{}
	}

	var totalBaseSalary, totalOTPay, totalNightPay, totalLeavePay, totalHolidayPay float64
	var totalAllowances, totalBonuses, totalOtherIncome float64
	var totalSIDeduction, totalHIDeduction, totalUIDeduction float64
	var totalTradeUnion, totalPIT float64
	var totalEmployerSI, totalEmployerHI, totalEmployerUI float64
	var totalEmployerTU float64

	for _, r := range input.Runs {
		totalBaseSalary += r.BaseSalary
		totalOTPay += r.OTPay
		totalNightPay += r.NightShiftPay
		totalLeavePay += r.LeavePay
		totalHolidayPay += r.HolidayPay
		totalAllowances += r.Allowances
		totalBonuses += r.Bonuses
		totalOtherIncome += r.OtherIncome
		totalSIDeduction += r.SIDeduction
		totalHIDeduction += r.HIDeduction
		totalUIDeduction += r.UIDeduction
		totalTradeUnion += r.TradeUnionDues
		totalPIT += r.PITAmount
		totalEmployerSI += r.EmployerSI
		totalEmployerHI += r.EmployerHI
		totalEmployerUI += r.EmployerUI
		totalEmployerTU += r.EmployerTradeUnion
	}

	totalGross := totalBaseSalary + totalOTPay + totalNightPay + totalLeavePay +
		totalHolidayPay + totalAllowances + totalBonuses + totalOtherIncome

	desc := fmt.Sprintf("Payroll %02d/%d", input.PeriodMonth, input.PeriodYear)

	lines1 := []JournalLine{
		{LineNumber: 1, AccountCode: "6421", DebitAmount: totalBaseSalary, Description: "Lương cơ bản"},
		{LineNumber: 2, AccountCode: "6422", DebitAmount: totalOTPay + totalNightPay, Description: "Lương tăng ca ca đêm"},
		{LineNumber: 3, AccountCode: "6423", DebitAmount: totalAllowances + totalHolidayPay + totalBonuses + totalOtherIncome, Description: "Phụ cấp & thưởng"},
		{LineNumber: 4, AccountCode: "3331", CreditAmount: totalGross, Description: "Phải trả người lao động"},
	}

	lines2 := []JournalLine{
		{LineNumber: 1, AccountCode: "6424", DebitAmount: totalEmployerSI + totalEmployerHI + totalEmployerUI + totalEmployerTU, Description: "Chi phí bảo hiểm NSDLĐ"},
		{LineNumber: 2, AccountCode: "3332", CreditAmount: totalEmployerSI, Description: "BHXH phần sử dụng lao động"},
		{LineNumber: 3, AccountCode: "3333", CreditAmount: totalEmployerHI, Description: "BHYT phần sử dụng lao động"},
		{LineNumber: 4, AccountCode: "3334", CreditAmount: totalEmployerUI, Description: "BHTN phần sử dụng lao động"},
		{LineNumber: 5, AccountCode: "3339", CreditAmount: totalEmployerTU, Description: "Công đoàn phần sử dụng lao động"},
	}

	// Entry 3: Employee deductions (offset against salary payable)
	lines3 := []JournalLine{
		{LineNumber: 1, AccountCode: "3331", DebitAmount: totalSIDeduction + totalHIDeduction + totalUIDeduction + totalTradeUnion + totalPIT, Description: "Trừ tiền lương NLĐ"},
		{LineNumber: 2, AccountCode: "3335", CreditAmount: totalSIDeduction, Description: "BHXH phần NLĐ"},
		{LineNumber: 3, AccountCode: "3336", CreditAmount: totalHIDeduction, Description: "BHYT phần NLĐ"},
		{LineNumber: 4, AccountCode: "3337", CreditAmount: totalUIDeduction, Description: "BHTN phần NLĐ"},
		{LineNumber: 5, AccountCode: "3339", CreditAmount: totalTradeUnion, Description: "Công đoàn phần NLĐ"},
		{LineNumber: 6, AccountCode: "3338", CreditAmount: totalPIT, Description: "Thuế TNCN"},
	}

	salaryEntry := &JournalEntry{
		CompanyID:   input.CompanyID,
		EntryNumber: fmt.Sprintf("PC-%02d%04d", input.PeriodMonth, input.PeriodYear),
		EntryDate:   time.Now(),
		Description: desc + " — Lương và phân bổ",
		Status:      JournalEntryDraft,
		Lines:       lines1,
	}

	insuranceEntry := &JournalEntry{
		CompanyID:   input.CompanyID,
		EntryNumber: fmt.Sprintf("PC-BH-%02d%04d", input.PeriodMonth, input.PeriodYear),
		EntryDate:   time.Now(),
		Description: desc + " — Bảo hiểm NSDLĐ",
		Status:      JournalEntryDraft,
		Lines:       lines2,
	}

	deductionEntry := &JournalEntry{
		CompanyID:   input.CompanyID,
		EntryNumber: fmt.Sprintf("PC-TRU-%02d%04d", input.PeriodMonth, input.PeriodYear),
		EntryDate:   time.Now(),
		Description: desc + " — Trừ tiền lương NLĐ",
		Status:      JournalEntryDraft,
		Lines:       lines3,
	}

	return PayrollJournalEntries{
		Entries: []*JournalEntry{salaryEntry, insuranceEntry, deductionEntry},
	}
}

// ─── Net-to-Gross (Reverse Calculation) ─────────────────────────

// NetToGrossInput holds known values for reverse-calculating gross salary.
type NetToGrossInput struct {
	TargetNetPay    float64 `json:"target_net_pay"`
	InsuranceBase   float64 `json:"insurance_base"`
	UIBase          float64 `json:"ui_base"`
	IsForeign       bool    `json:"is_foreign"`
	DependantCount  int     `json:"dependant_count"`
	OtherDeductions float64 `json:"other_deductions"`
	Month           int     `json:"month"`
	Year            int     `json:"year"`
}

// NetToGrossResult holds the reverse-calculated gross and breakdown.
type NetToGrossResult struct {
	GrossSalary     float64 `json:"gross_salary"`
	BaseSalary      float64 `json:"base_salary"`
	SIDeduction     float64 `json:"si_deduction"`
	HIDeduction     float64 `json:"hi_deduction"`
	UIDeduction     float64 `json:"ui_deduction"`
	TradeUnionDues  float64 `json:"trade_union_dues"`
	PITAmount       float64 `json:"pit_amount"`
	TotalDeductions float64 `json:"total_deductions"`
	NetPay          float64 `json:"net_pay"`
}

// CalcNetToGross reverse-engineers gross salary from target net pay using binary search.
// Tolerance: 1 VND.
func CalcNetToGross(input NetToGrossInput) NetToGrossResult {
	if input.TargetNetPay <= 0 {
		return NetToGrossResult{}
	}
	if input.InsuranceBase == 0 {
		input.InsuranceBase = input.TargetNetPay // fallback
	}
	if input.UIBase == 0 {
		input.UIBase = input.InsuranceBase
	}

	// Binary search for gross salary
	low, high := 0.0, input.TargetNetPay*5 // upper bound heuristic
	var best GrossNetPair

	for i := 0; i < 100; i++ { // max iterations
		mid := (low + high) / 2
		gn := calcGrossNet(mid, input)
		diff := gn.Net - input.TargetNetPay

		if math.Abs(diff) < 1 { // within 1 VND
			best = gn
			break
		}

		if diff < 0 {
			low = mid
		} else {
			high = mid
		}
		best = gn
	}

	// Recalculate with final best gross for accurate breakdown
	result := calcGrossNet(best.Gross, input)
	return NetToGrossResult{
		GrossSalary:     round0(result.Gross),
		BaseSalary:      round0(result.Gross),
		SIDeduction:     round0(result.SI),
		HIDeduction:     round0(result.HI),
		UIDeduction:     round0(result.UI),
		TradeUnionDues:  round0(result.TU),
		PITAmount:       round0(result.PIT),
		TotalDeductions: round0(result.SI + result.HI + result.UI + result.TU + result.PIT + input.OtherDeductions),
		NetPay:          round0(result.Net),
	}
}

type GrossNetPair struct {
	Gross float64
	Net   float64
	SI    float64
	HI    float64
	UI    float64
	TU    float64
	PIT   float64
}

func calcGrossNet(gross float64, input NetToGrossInput) GrossNetPair {
	si := CalcEmployeeSI(input.InsuranceBase, input.Month, input.Year)
	hi := CalcEmployeeHI(input.InsuranceBase, input.Month, input.Year)
	ui := CalcEmployeeUI(input.UIBase, input.IsForeign)
	tu := CalcTradeUnion(input.InsuranceBase, input.IsForeign, input.Month, input.Year)

	taxable := gross - si - hi - ui - tu - PersonalDeductionMonthly - float64(input.DependantCount)*DependantDeductionMonthly
	pit := CalcPIT(taxable, input.Month, input.Year)

	totalDed := si + hi + ui + tu + pit + input.OtherDeductions
	net := gross - totalDed

	return GrossNetPair{Gross: gross, Net: net, SI: si, HI: hi, UI: ui, TU: tu, PIT: pit}
}

// ─── 13th-Month Salary ──────────────────────────────────────────

// ThirteenthMonthInput holds inputs for 13th-month salary calculation.
type ThirteenthMonthInput struct {
	BaseSalary      float64 `json:"base_salary"`
	MonthsWorked    int     `json:"months_worked"`     // 1-12
	Allowances      float64 `json:"allowances"`         // monthly allowances total
	IsForeign       bool    `json:"is_foreign"`
	DependantCount  int     `json:"dependant_count"`
	InsuranceBase   float64 `json:"insurance_base"`
	UIBase          float64 `json:"ui_base"`
	Month           int     `json:"month"`              // payout month (for PIT bracket selection)
	Year            int     `json:"year"`
}

// ThirteenthMonthResult holds the 13th-month salary breakdown.
type ThirteenthMonthResult struct {
	GrossAmount     float64 `json:"gross_amount"`
	SIDeduction     float64 `json:"si_deduction"`
	HIDeduction     float64 `json:"hi_deduction"`
	UIDeduction     float64 `json:"ui_deduction"`
	TradeUnionDues  float64 `json:"trade_union_dues"`
	PITAmount       float64 `json:"pit_amount"`
	TotalDeductions float64 `json:"total_deductions"`
	NetPay          float64 `json:"net_pay"`
	EmployerSI      float64 `json:"employer_si"`
	EmployerHI      float64 `json:"employer_hi"`
	EmployerUI      float64 `json:"employer_ui"`
	TotalEmployerCost float64 `json:"total_employer_cost"`
}

// CalcThirteenthMonth computes 13th-month salary with full deductions.
// Per Vietnamese Labour Code Art. 103: employees with <12 months get proportional amount.
func CalcThirteenthMonth(input ThirteenthMonthInput) ThirteenthMonthResult {
	if input.MonthsWorked <= 0 || input.BaseSalary <= 0 {
		return ThirteenthMonthResult{}
	}
	if input.MonthsWorked > 12 {
		input.MonthsWorked = 12
	}

	// Gross = (months worked / 12) × base salary
	gross := round0(input.BaseSalary * float64(input.MonthsWorked) / 12.0)

	// Insurance base for 13th month = proportional insurance base
	insBase := input.InsuranceBase
	if insBase == 0 {
		insBase = input.BaseSalary
	}
	monthInsBase := insBase * float64(input.MonthsWorked) / 12.0

	uiBase := input.UIBase
	if uiBase == 0 {
		uiBase = input.BaseSalary
	}
	monthUIBase := uiBase * float64(input.MonthsWorked) / 12.0

	// Employee deductions
	si := CalcEmployeeSI(monthInsBase, input.Month, input.Year)
	hi := CalcEmployeeHI(monthInsBase, input.Month, input.Year)
	ui := CalcEmployeeUI(monthUIBase, input.IsForeign)
	tu := CalcTradeUnion(monthInsBase, input.IsForeign, input.Month, input.Year)

	// PIT: taxable = gross - insurance - personal deduction
	// 13th month is separate from regular salary for PIT purposes
	taxable := gross - si - hi - ui - tu - PersonalDeductionMonthly
	pit := CalcPIT(taxable, input.Month, input.Year)

	totalDed := si + hi + ui + tu + pit
	net := gross - totalDed

	// Employer costs
	emplSI := CalcEmployerSI(monthInsBase, input.Month, input.Year)
	emplHI := CalcEmployerHI(monthInsBase, input.Month, input.Year)
	emplUI := CalcEmployerUI(monthUIBase, input.IsForeign)

	return ThirteenthMonthResult{
		GrossAmount:       gross,
		SIDeduction:       si,
		HIDeduction:       hi,
		UIDeduction:       ui,
		TradeUnionDues:    tu,
		PITAmount:         pit,
		TotalDeductions:   totalDed,
		NetPay:            net,
		EmployerSI:        emplSI,
		EmployerHI:        emplHI,
		EmployerUI:        emplUI,
		TotalEmployerCost: emplSI + emplHI + emplUI,
	}
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

// ─── Severance Pay Calculator ───────────────────────────────────

// CalcSeverancePay computes severance per Labor Code 2019 Art. 46.
// Rule: 0.5 month salary × years of service, based on avg of last 6 months.
// No SI/HI/UI deductions on severance. PIT applies on gross.
// Gross misconduct = no severance.
func CalcSeverancePay(input SeveranceInput) SeveranceResult {
	result := SeveranceResult{
		YearsOfService: input.YearsOfService,
		Reason:         input.Reason,
	}

	// Gross misconduct = zero severance
	if input.Reason == TerminationGrossMisconduct {
		return result
	}

	// Gross severance: 0.5 × avg × years
	result.GrossSeverance = math.Round(0.5*input.AvgSalary6Months*input.YearsOfService*100) / 100
	if result.GrossSeverance < 0 {
		result.GrossSeverance = 0
	}

	// PIT on severance (use current month's brackets)
	now := time.Now()
	result.PITAmount = CalcPIT(result.GrossSeverance, int(now.Month()), now.Year())
	result.NetSeverance = result.GrossSeverance - result.PITAmount

	return result
}
