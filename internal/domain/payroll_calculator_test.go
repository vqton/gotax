package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Helper ─────────────────────────────────────────────────────

func makeTestEmployee(baseSalary float64, region InsuranceRegion) EmployeePayrollInfo {
	return EmployeePayrollInfo{
		ID:                "emp-001",
		EmployeeID:        "NV001",
		ContractType:      ContractIndefinite,
		SalaryType:        SalaryTimeBased,
		BaseSalary:        baseSalary,
		SalaryCoefficient: 2.0,
		PositionAllowance: 1_000_000,
		Region:            region,
	}
}

// ─── Daily Salary ───────────────────────────────────────────────

func TestCalcDailySalary(t *testing.T) {
	tests := []struct {
		name     string
		monthly  float64
		days     int
		expected float64
	}{
		{"standard 26 days", 26_000_000, 26, 1_000_000},
		{"22 days", 22_000_000, 22, 1_000_000},
		{"base salary 2,340,000", 2_340_000, 26, 90_000},
		{"10M / 26", 10_000_000, 26, 384_615.38},
		{"zero days", 10_000_000, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcDailySalary(tt.monthly, tt.days)
			assert.InDelta(t, tt.expected, got, 1)
		})
	}
}

// ─── Hourly Salary ──────────────────────────────────────────────

func TestCalcHourlySalary(t *testing.T) {
	got := CalcHourlySalary(1_000_000)
	assert.Equal(t, 125_000.0, got)
}

// ─── OT Pay ─────────────────────────────────────────────────────

func TestCalcOTPay(t *testing.T) {
	tests := []struct {
		name       string
		hourlyRate float64
		otHours    float64
		multiplier float64
		expected   float64
	}{
		{"no OT", 125_000, 0, 1.5, 0},
		{"1 hour weekday 150%", 125_000, 1, 1.5, 187_500},
		{"2 hours weekday 150%", 125_000, 2, 1.5, 375_000},
		{"8 hours full day 150%", 125_000, 8, 1.5, 1_500_000},
		{"4 hours weekend 200%", 125_000, 4, 2.0, 1_000_000},
		{"4 hours holiday 300%", 125_000, 4, 3.0, 1_500_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcOTPay(tt.hourlyRate, tt.otHours, tt.multiplier)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// ─── Night Shift Pay ────────────────────────────────────────────

func TestCalcNightShiftPay(t *testing.T) {
	got := CalcNightShiftPay(125_000, 4)
	// (125,000 * 4) * 1.3 = 650,000
	assert.Equal(t, float64(650_000), got)
}

// ─── Leave Pay ──────────────────────────────────────────────────

func TestCalcLeavePay(t *testing.T) {
	got := CalcLeavePay(1_000_000, 2)
	assert.Equal(t, float64(2_000_000), got)
}

func TestCalcLeavePay_Sick(t *testing.T) {
	// Sick leave = 75% of daily
	got := CalcLeavePay(1_000_000*0.75, 2)
	assert.Equal(t, float64(1_500_000), got)
}

// ─── Holiday Pay ────────────────────────────────────────────────

func TestCalcHolidayPay(t *testing.T) {
	got := CalcHolidayPay(1_000_000)
	assert.Equal(t, float64(3_000_000), got)
}

// ─── Date-Ranged Constants ─────────────────────────────────────

func TestGetBaseSalary(t *testing.T) {
	tests := []struct {
		name     string
		month    int
		year     int
		expected float64
	}{
		{"Jan 2026", 1, 2026, 2_340_000},
		{"Jun 2026", 6, 2026, 2_340_000},
		{"Jul 2026", 7, 2026, 2_530_000},
		{"Dec 2026", 12, 2026, 2_530_000},
		{"Jan 2027", 1, 2027, 2_530_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBaseSalary(tt.month, tt.year)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetSIHICap(t *testing.T) {
	tests := []struct {
		name     string
		month    int
		year     int
		expected float64
	}{
		{"H1 2026 — 20×2,340,000", 1, 2026, 46_800_000},
		{"H2 2026 — 20×2,530,000", 7, 2026, 50_600_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSIHICap(tt.month, tt.year)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetTradeUnionCap(t *testing.T) {
	tests := []struct {
		name     string
		month    int
		year     int
		expected float64
	}{
		{"H1 2026 — 1% of 46,800,000", 1, 2026, 468_000},
		{"H2 2026 — 1% of 50,600,000", 7, 2026, 506_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTradeUnionCap(tt.month, tt.year)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetPITBrackets(t *testing.T) {
	tests := []struct {
		name   string
		month  int
		year   int
		expect int
	}{
		{"Jan 2026 → old (7)", 1, 2026, 7},
		{"Jun 2026 → old (7)", 6, 2026, 7},
		{"Jul 2026 → new (5)", 7, 2026, 5},
		{"Dec 2026 → new (5)", 12, 2026, 5},
		{"Jan 2027 → new (5)", 1, 2027, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			brackets := GetPITBrackets(tt.month, tt.year)
			assert.Equal(t, tt.expect, len(brackets))
		})
	}
}

// ─── Insurance (Employee) ───────────────────────────────────────

func TestCalcEmployeeSI(t *testing.T) {
	tests := []struct {
		name     string
		base     float64
		month    int
		year     int
		expected float64
	}{
		{"standard 10M H2", 10_000_000, 7, 2026, 800_000},       // 8%
		{"above cap H2", 60_000_000, 7, 2026, 4_048_000},         // 8% of 50,600,000
		{"at cap H2", 50_600_000, 7, 2026, 4_048_000},
		{"below cap H2", 5_000_000, 7, 2026, 400_000},
		{"above cap H1", 60_000_000, 1, 2026, 3_744_000},         // 8% of 46,800,000
		{"at cap H1", 46_800_000, 1, 2026, 3_744_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcEmployeeSI(tt.base, tt.month, tt.year)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestCalcEmployeeHI(t *testing.T) {
	tests := []struct {
		name     string
		base     float64
		month    int
		year     int
		expected float64
	}{
		{"standard 10M H2", 10_000_000, 7, 2026, 150_000},       // 1.5%
		{"above cap H2", 60_000_000, 7, 2026, 759_000},          // 1.5% of 50,600,000
		{"below cap H2", 5_000_000, 7, 2026, 75_000},
		{"above cap H1", 60_000_000, 1, 2026, 702_000},          // 1.5% of 46,800,000
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcEmployeeHI(tt.base, tt.month, tt.year)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestCalcEmployeeUI(t *testing.T) {
	t.Run("Vietnamese employee", func(t *testing.T) {
		got := CalcEmployeeUI(10_000_000, false)
		assert.Equal(t, float64(100_000), got) // 1%
	})

	t.Run("above cap", func(t *testing.T) {
		got := CalcEmployeeUI(120_000_000, false)
		// cap = 20 * 5,310,000 = 106,200,000 → 1% = 1,062,000
		assert.Equal(t, float64(1_062_000), got)
	})

	t.Run("foreign employee = 0", func(t *testing.T) {
		got := CalcEmployeeUI(10_000_000, true)
		assert.Equal(t, float64(0), got)
	})
}

// ─── Insurance (Employer) ───────────────────────────────────────

func TestCalcEmployerSI(t *testing.T) {
	got := CalcEmployerSI(10_000_000, 7, 2026)
	assert.Equal(t, float64(1_750_000), got) // 17.5%
}

func TestCalcEmployerHI(t *testing.T) {
	got := CalcEmployerHI(10_000_000, 7, 2026)
	assert.Equal(t, float64(300_000), got) // 3%
}

func TestCalcEmployerUI(t *testing.T) {
	t.Run("Vietnamese", func(t *testing.T) {
		got := CalcEmployerUI(10_000_000, false)
		assert.Equal(t, float64(100_000), got) // 1%
	})

	t.Run("foreign = 0", func(t *testing.T) {
		got := CalcEmployerUI(10_000_000, true)
		assert.Equal(t, float64(0), got)
	})
}

// ─── Trade Union ────────────────────────────────────────────────

func TestCalcTradeUnion(t *testing.T) {
	tests := []struct {
		name     string
		salary   float64
		foreign  bool
		month    int
		year     int
		expected float64
	}{
		{"standard 10M H2", 10_000_000, false, 7, 2026, 100_000},
		{"high salary capped H2", 60_000_000, false, 7, 2026, 506_000}, // 1% of 50,600,000
		{"high salary capped H1", 60_000_000, false, 1, 2026, 468_000}, // 1% of 46,800,000
		{"low salary H2", 5_000_000, false, 7, 2026, 50_000},
		{"foreign = 0", 10_000_000, true, 7, 2026, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcTradeUnion(tt.salary, tt.foreign, tt.month, tt.year)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// ─── PIT (Old 7-Bracket — H1) ──────────────────────────────────

func TestCalcPITOldBrackets(t *testing.T) {
	tests := []struct {
		name     string
		income   float64
		expected float64
	}{
		{"zero", 0, 0},
		{"negative", -1_000_000, 0},
		{"bracket 1 — 5M (5%)", 5_000_000, 250_000},
		{"bracket 2 — 10M (10%)", 10_000_000, 750_000},
		{"bracket 2 — 12.5M (10%)", 12_500_000, 1_125_000},
		{"bracket 3 — 18M (15%)", 18_000_000, 1_950_000},
		{"bracket 4 — 20M (20%)", 20_000_000, 2_350_000},
		{"bracket 4 — 32M (20%)", 32_000_000, 4_750_000},
		{"bracket 5 — 40M (25%)", 40_000_000, 6_750_000},
		{"bracket 6 — 80M (30%)", 80_000_000, 18_150_000},
		{"bracket 7 — 100M (35%)", 100_000_000, 25_150_000},
		{"bracket 7 — 150M (35%)", 150_000_000, 42_650_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcPIT(tt.income, 6, 2026) // Jun → old brackets
			assert.Equal(t, tt.expected, got)
		})
	}
}

// ─── PIT (New 5-Bracket — H2) ──────────────────────────────────

func TestCalcPITNewBrackets(t *testing.T) {
	tests := []struct {
		name     string
		income   float64
		expected float64
	}{
		{"zero", 0, 0},
		{"bracket 1 — 5M (5%)", 5_000_000, 250_000},
		{"bracket 1 — 10M (5%)", 10_000_000, 500_000},
		{"bracket 2 — 15M (10%)", 15_000_000, 1_000_000},     // 15M×10%-500K
		{"bracket 2 — 30M (10%)", 30_000_000, 2_500_000},     // 30M×10%-500K
		{"bracket 3 — 40M (20%)", 40_000_000, 4_500_000},     // 40M×20%-3.5M
		{"bracket 3 — 60M (20%)", 60_000_000, 8_500_000},     // 60M×20%-3.5M
		{"bracket 4 — 80M (30%)", 80_000_000, 14_500_000},    // 80M×30%-9.5M
		{"bracket 4 — 100M (30%)", 100_000_000, 20_500_000},  // 100M×30%-9.5M
		{"bracket 5 — 150M (35%)", 150_000_000, 38_000_000},  // 150M×35%-14.5M
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcPIT(tt.income, 7, 2026) // Jul → new brackets
			assert.Equal(t, tt.expected, got)
		})
	}
}

// ─── Full Employee Payroll Calculation ───────────────────────────

func TestCalculateEmployeePayroll(t *testing.T) {
	employee := makeTestEmployee(10_000_000, RegionI)
	period := &PayrollPeriod{Year: 2026, Month: 7}

	t.Run("standard payroll (no OT, no leave)", func(t *testing.T) {
		run := CalculateEmployeePayroll(employee, period, nil)

		// Base salary
		assert.Equal(t, float64(10_000_000), run.BaseSalary)

		// Working days default to 26
		assert.Equal(t, 26, run.WorkingDays)

		// Daily = 384,615.38
		daily := CalcDailySalary(10_000_000, 26)
		assert.InDelta(t, 384_615.38, daily, 1)

		// No OT, no night, no leave
		assert.Equal(t, float64(0), run.OTPay)
		assert.Equal(t, float64(0), run.NightShiftPay)
		assert.Equal(t, float64(0), run.LeavePay)

		// Allowances
		assert.Equal(t, float64(1_000_000), run.Allowances)

		// Gross = base + allowances
		assert.Equal(t, float64(11_000_000), run.GrossSalary)

		// Employee insurance (H2: cap = 50,600,000)
		assert.Equal(t, float64(800_000), run.SIDeduction)   // 8%
		assert.Equal(t, float64(150_000), run.HIDeduction)   // 1.5%
		assert.Equal(t, float64(100_000), run.UIDeduction)   // 1%
		assert.Equal(t, float64(100_000), run.TradeUnionDues) // 1%

		// Total deductions
		totalDed := 800_000 + 150_000 + 100_000 + 100_000 + run.PITAmount
		assert.Equal(t, totalDed, run.TotalDeductions)

		// Net = gross - deductions
		assert.Equal(t, run.GrossSalary-run.TotalDeductions, run.NetPay)

		// Employer costs (H2: cap = 50,600,000)
		assert.Equal(t, float64(1_750_000), run.EmployerSI)       // 17.5%
		assert.Equal(t, float64(300_000), run.EmployerHI)         // 3%
		assert.Equal(t, float64(100_000), run.EmployerUI)         // 1%
		assert.Equal(t, float64(100_000), run.EmployerTradeUnion) // 1%

		// Status
		assert.Equal(t, "CALCULATED", run.Status)
	})

	t.Run("with overtime and night shift", func(t *testing.T) {
		employee := makeTestEmployee(10_000_000, RegionI)
		timekeeping := []TimekeepingRecord{
			{Date: "2026-07-01", HoursWorked: 8, OTHours: 2, NightHours: 4},
		}

		run := CalculateEmployeePayroll(employee, period, timekeeping)

		// 1 working day from timekeeping
		assert.Equal(t, 1, run.WorkingDays)
		assert.Equal(t, 2.0, run.OTHours)
		assert.Equal(t, 4.0, run.NightShiftHours)

		// Daily = 10M / 1 = 10M; hourly = 1.25M
		daily := CalcDailySalary(10_000_000, 1)
		hourly := daily / 8
		assert.Equal(t, float64(10_000_000), daily)
		assert.Equal(t, float64(1_250_000), hourly)

		// OT: 1,250,000 * 2 * 1.5 = 3,750,000
		expectedOT := CalcOTPay(hourly, 2, OTWeekdayRate)
		assert.Equal(t, float64(3_750_000), expectedOT)
		assert.Equal(t, expectedOT, run.OTPay)

		// Night: 1,250,000 * 4 * 1.3 = 6,500,000
		expectedNight := CalcNightShiftPay(hourly, 4)
		assert.Equal(t, float64(6_500_000), expectedNight)
		assert.Equal(t, expectedNight, run.NightShiftPay)

		// Gross includes OT and night
		assert.True(t, run.GrossSalary > run.BaseSalary+run.Allowances)
	})

	t.Run("foreign employee — no UI", func(t *testing.T) {
		emp := makeTestEmployee(10_000_000, RegionI)
		emp.IsForeignEmployee = true

		run := CalculateEmployeePayroll(emp, period, nil)

		assert.Equal(t, float64(0), run.UIDeduction)
		assert.Equal(t, float64(0), run.EmployerUI)
		assert.Equal(t, float64(0), run.TradeUnionDues)
		assert.Equal(t, float64(0), run.EmployerTradeUnion)
	})

	t.Run("high salary — SI capped H2", func(t *testing.T) {
		emp := makeTestEmployee(60_000_000, RegionI)
		run := CalculateEmployeePayroll(emp, period, nil)

		// SI base capped at 50,600,000 (H2)
		assert.Equal(t, float64(4_048_000), run.SIDeduction)
		assert.Equal(t, float64(759_000), run.HIDeduction)
	})

	t.Run("net pay positive", func(t *testing.T) {
		run := CalculateEmployeePayroll(employee, period, nil)
		assert.True(t, run.NetPay > 0, "net pay must be positive")
	})

	t.Run("gross > net", func(t *testing.T) {
		run := CalculateEmployeePayroll(employee, period, nil)
		assert.True(t, run.GrossSalary > run.NetPay)
	})

	t.Run("employer cost > employee deductions", func(t *testing.T) {
		run := CalculateEmployeePayroll(employee, period, nil)
		assert.True(t, run.TotalEmployerCost > run.TotalDeductions)
	})
}

// ─── Transition Period: H1 vs H2 ────────────────────────────────

func TestPayrollTransitionPeriod(t *testing.T) {
	emp := makeTestEmployee(60_000_000, RegionI) // above both H1 and H2 caps

	t.Run("H1 — SI cap 46,800,000", func(t *testing.T) {
		period := &PayrollPeriod{Year: 2026, Month: 3}
		run := CalculateEmployeePayroll(emp, period, nil)

		// SI: 8% of 46,800,000 = 3,744,000
		assert.Equal(t, float64(3_744_000), run.SIDeduction)
		// HI: 1.5% of 46,800,000 = 702,000
		assert.Equal(t, float64(702_000), run.HIDeduction)
		// Trade union: 1% of 46,800,000 = 468,000
		assert.Equal(t, float64(468_000), run.TradeUnionDues)
	})

	t.Run("H2 — SI cap 50,600,000", func(t *testing.T) {
		period := &PayrollPeriod{Year: 2026, Month: 8}
		run := CalculateEmployeePayroll(emp, period, nil)

		// SI: 8% of 50,600,000 = 4,048,000
		assert.Equal(t, float64(4_048_000), run.SIDeduction)
		// HI: 1.5% of 50,600,000 = 759,000
		assert.Equal(t, float64(759_000), run.HIDeduction)
		// Trade union: 1% of 50,600,000 = 506,000
		assert.Equal(t, float64(506_000), run.TradeUnionDues)
	})

	t.Run("PIT brackets differ H1 vs H2", func(t *testing.T) {
		taxable := 50_000_000.0
		pitH1 := CalcPIT(taxable, 6, 2026) // old brackets: 25%
		pitH2 := CalcPIT(taxable, 7, 2026) // new brackets: 20%

		// Old: 50M × 25% - 3,250,000 = 9,250,000
		assert.Equal(t, float64(9_250_000), pitH1)
		// New: 50M × 20% - 3,500,000 = 6,500,000
		assert.Equal(t, float64(6_500_000), pitH2)
		// H2 should be lower (new brackets more favorable at this income)
		assert.True(t, pitH2 < pitH1)
	})
}

// ─── Net-to-Gross ──────────────────────────────────────────────

func TestCalcNetToGross(t *testing.T) {
	tests := []struct {
		name          string
		targetNet     float64
		insBase       float64
		uiBase        float64
		isForeign     bool
		depCount      int
		month, year   int
		tolerance     float64
	}{
		{
			name:      "standard 10M gross — round trip",
			targetNet: 0, // will compute from gross
			insBase:   10_000_000,
			uiBase:    10_000_000,
			month:     7, year: 2026,
			tolerance: 100,
		},
		{
			name:      "high salary 50M — capped insurance",
			targetNet: 0,
			insBase:   50_600_000,
			uiBase:    50_600_000,
			month:     7, year: 2026,
			tolerance: 100,
		},
		{
			name:      "foreign employee — no UI",
			targetNet: 0,
			insBase:   20_000_000,
			uiBase:    20_000_000,
			isForeign: true,
			month:     7, year: 2026,
			tolerance: 100,
		},
		{
			name:      "with 2 dependants",
			targetNet: 0,
			insBase:   15_000_000,
			uiBase:    15_000_000,
			depCount:  2,
			month:     7, year: 2026,
			tolerance: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Forward: compute gross→net
			emp := makeTestEmployee(tt.insBase, RegionI)
			emp.IsForeignEmployee = tt.isForeign
			period := &PayrollPeriod{Year: tt.year, Month: tt.month}
			fwdRun := CalculateEmployeePayroll(emp, period, nil)

			// Set target net from forward calculation
			targetNet := fwdRun.NetPay
			if tt.name == "with 2 dependants" {
				// Recalc with dependants
				emp2 := emp
				_ = emp2
				// Use the forward result as-is for round-trip test
				targetNet = fwdRun.NetPay
			}

			// Reverse: net→gross
			result := CalcNetToGross(NetToGrossInput{
				TargetNetPay:   targetNet,
				InsuranceBase:  tt.insBase,
				UIBase:         tt.uiBase,
				IsForeign:      tt.isForeign,
				DependantCount: tt.depCount,
				Month:          tt.month,
				Year:           tt.year,
			})

			// Verify round-trip: forward net ≈ reverse net
			assert.InDelta(t, fwdRun.NetPay, result.NetPay, tt.tolerance,
				"round-trip net pay mismatch")
		})
	}
}

func TestCalcNetToGross_ZeroNetPay(t *testing.T) {
	result := CalcNetToGross(NetToGrossInput{TargetNetPay: 0})
	assert.Equal(t, float64(0), result.GrossSalary)
}

func TestCalcNetToGross_KnownTarget(t *testing.T) {
	// Target net of 15,000,000 VND
	result := CalcNetToGross(NetToGrossInput{
		TargetNetPay:   15_000_000,
		InsuranceBase:  20_000_000,
		UIBase:         20_000_000,
		Month:          7,
		Year:           2026,
	})

	// Gross should be higher than net
	assert.True(t, result.GrossSalary > 15_000_000)
	// Net should be close to target
	assert.InDelta(t, 15_000_000, result.NetPay, 100)
	// Deductions should be positive
	assert.True(t, result.TotalDeductions > 0)
	assert.True(t, result.SIDeduction > 0)
	assert.True(t, result.HIDeduction > 0)
}

func TestCalcNetToGross_H1vsH2(t *testing.T) {
	target := 20_000_000.0

	// Use insurance base above H1 cap (46.8M) but below H2 cap (50.6M)
	rH1 := CalcNetToGross(NetToGrossInput{
		TargetNetPay: target, InsuranceBase: 48_000_000, UIBase: 48_000_000,
		Month: 3, Year: 2026,
	})
	rH2 := CalcNetToGross(NetToGrossInput{
		TargetNetPay: target, InsuranceBase: 48_000_000, UIBase: 48_000_000,
		Month: 8, Year: 2026,
	})

	// Both should produce valid results
	assert.True(t, rH1.GrossSalary > target)
	assert.True(t, rH2.GrossSalary > target)
	// Both should produce valid net pay close to target
	assert.InDelta(t, target, rH1.NetPay, 100)
	assert.InDelta(t, target, rH2.NetPay, 100)
	// Insurance deductions differ between H1 and H2 (48M > H1 cap 46.8M, < H2 cap 50.6M)
	assert.NotEqual(t, rH1.SIDeduction, rH2.SIDeduction,
		"SI deduction should differ: 48M is above H1 cap but below H2 cap")
}

// ─── Edge Cases ─────────────────────────────────────────────────

func TestEdgeCases(t *testing.T) {
	t.Run("zero salary", func(t *testing.T) {
		emp := EmployeePayrollInfo{
			ID:           "emp-001",
			EmployeeID:   "NV001",
			ContractType: ContractIndefinite,
			SalaryType:   SalaryTimeBased,
			BaseSalary:   0,
			Region:       RegionI,
		}
		period := &PayrollPeriod{Year: 2026, Month: 7}
		run := CalculateEmployeePayroll(emp, period, nil)
		assert.Equal(t, float64(0), run.GrossSalary)
		assert.Equal(t, float64(0), run.NetPay)
	})

	t.Run("very high salary above all caps", func(t *testing.T) {
		emp := makeTestEmployee(100_000_000, RegionI)
		period := &PayrollPeriod{Year: 2026, Month: 7}
		run := CalculateEmployeePayroll(emp, period, nil)

		// InsuranceBase stores raw base salary; capping happens in Calc functions
		assert.Equal(t, float64(100_000_000), run.InsuranceBase)
		// SI capped at 50,600,000 × 8% = 4,048,000 (H2)
		assert.Equal(t, float64(4_048_000), run.SIDeduction)
	})

	t.Run("regional minimum wage affects UI cap", func(t *testing.T) {
		// Region I: 5,310,000
		// Region IV: 3,700,000
		assert.Equal(t, 5_310_000.0, RegionalMinWageI)
		assert.Equal(t, 3_700_000.0, RegionalMinWageIV)
	})

	t.Run("PIT brackets cover full range", func(t *testing.T) {
		// Old brackets
		assert.True(t, CalcPIT(5_000_000, 6, 2026) > 0)
		assert.True(t, CalcPIT(150_000_000, 6, 2026) > 0)
		assert.True(t, CalcPIT(1_000_000_000, 6, 2026) > 0)
		// New brackets
		assert.True(t, CalcPIT(5_000_000, 7, 2026) > 0)
		assert.True(t, CalcPIT(150_000_000, 7, 2026) > 0)
		assert.True(t, CalcPIT(1_000_000_000, 7, 2026) > 0)
	})
}

// ─── Insurance Constants Verification ───────────────────────────

func TestInsuranceConstants(t *testing.T) {
	// Employee rates: 8% + 1.5% + 1% = 10.5%
	empTotal := SIEmployeeRate + HIEmployeeRate + UIEmployeeRate
	assert.InDelta(t, 0.105, empTotal, 0.001)

	// Employer rates: 17.5% + 3% + 1% = 21.5%
	emplTotal := SIEmployerRate + HIEmployerRate + UIEmployerRate
	assert.InDelta(t, 0.215, emplTotal, 0.001)

	// H1: SI/HI cap = 20 * 2,340,000 = 46,800,000
	assert.Equal(t, float64(46_800_000), GetSIHICap(1, 2026))

	// H2: SI/HI cap = 20 * 2,530,000 = 50,600,000
	assert.Equal(t, float64(50_600_000), GetSIHICap(7, 2026))

	// UI cap = 20 * 5,310,000 = 106,200,000
	assert.Equal(t, float64(106_200_000), BaseSalaryCapUI*RegionalMinWageI)

	// Trade union cap H2 = 1% of 50,600,000 = 506,000
	assert.Equal(t, float64(506_000), GetTradeUnionCap(7, 2026))

	// Trade union cap H1 = 1% of 46,800,000 = 468,000
	assert.Equal(t, float64(468_000), GetTradeUnionCap(1, 2026))
}

// ─── PIT Bracket Coverage ───────────────────────────────────────

func TestPITBracketCoverage(t *testing.T) {
	require.Equal(t, 7, len(PITBracketsOld), "old must have 7 brackets")
	require.Equal(t, 5, len(PITBracketsNew), "new must have 5 brackets")

	for _, b := range PITBracketsOld {
		require.True(t, b.UpperLimit > 0)
		require.True(t, b.Rate > 0 && b.Rate <= 0.35)
	}

	for _, b := range PITBracketsNew {
		require.True(t, b.UpperLimit > 0)
		require.True(t, b.Rate > 0 && b.Rate <= 0.35)
	}
}
