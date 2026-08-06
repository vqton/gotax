package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── Severance Pay Tests ────────────────────────────────────────

func TestSeverancePay_Basic(t *testing.T) {
	// 10M avg, 3 years → 0.5 × 10M × 3 = 15M gross
	result := CalcSeverancePay(SeveranceInput{
		AvgSalary6Months: 10_000_000,
		YearsOfService:   3,
		Reason:           TerminationRedundancy,
	})
	assert.Equal(t, 15_000_000.0, result.GrossSeverance)
	assert.Equal(t, 3.0, result.YearsOfService)
	assert.Equal(t, TerminationRedundancy, result.Reason)
	assert.True(t, result.NetSeverance > 0, "net should be positive")
}

func TestSeverancePay_FractionalYears(t *testing.T) {
	// 15M avg, 2.5 years → 0.5 × 15M × 2.5 = 18.75M
	result := CalcSeverancePay(SeveranceInput{
		AvgSalary6Months: 15_000_000,
		YearsOfService:   2.5,
		Reason:           TerminationRestructure,
	})
	assert.Equal(t, 18_750_000.0, result.GrossSeverance)
}

func TestSeverancePay_GrossMisconduct_Zero(t *testing.T) {
	result := CalcSeverancePay(SeveranceInput{
		AvgSalary6Months: 10_000_000,
		YearsOfService:   5,
		Reason:           TerminationGrossMisconduct,
	})
	assert.Equal(t, 0.0, result.GrossSeverance)
	assert.Equal(t, 0.0, result.NetSeverance)
}

func TestSeverancePay_PITApplied(t *testing.T) {
	// High salary → significant PIT
	result := CalcSeverancePay(SeveranceInput{
		AvgSalary6Months: 50_000_000,
		YearsOfService:   5,
		Reason:           TerminationRedundancy,
	})
	// 0.5 × 50M × 5 = 125M gross
	assert.Equal(t, 125_000_000.0, result.GrossSeverance)
	assert.True(t, result.PITAmount > 0, "high severance should have PIT")
	assert.Equal(t, result.GrossSeverance-result.PITAmount, result.NetSeverance)
}

func TestSeverancePay_ZeroYears(t *testing.T) {
	result := CalcSeverancePay(SeveranceInput{
		AvgSalary6Months: 10_000_000,
		YearsOfService:   0,
		Reason:           TerminationExpiration,
	})
	assert.Equal(t, 0.0, result.GrossSeverance)
	assert.Equal(t, 0.0, result.NetSeverance)
}

func TestSeverancePay_AllReasons(t *testing.T) {
	reasons := []TerminationReason{
		TerminationRedundancy, TerminationRestructure, TerminationPerformance,
		TerminationExpiration, TerminationMutual, TerminationHealth, TerminationRetirement,
	}
	for _, reason := range reasons {
		result := CalcSeverancePay(SeveranceInput{
			AvgSalary6Months: 10_000_000,
			YearsOfService:   2,
			Reason:           reason,
		})
		assert.Equal(t, 10_000_000.0, result.GrossSeverance, "reason: %s", reason)
	}
}
