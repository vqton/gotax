package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── Retroactive Pay Tests ──────────────────────────────────────

func TestRetroactivePay_MidMonthIncrease(t *testing.T) {
	// 10M → 12M on day 16 of 30-day month
	result := CalcRetroactivePay(RetroactivePayInput{
		OldBaseSalary: 10_000_000,
		NewBaseSalary: 12_000_000,
		EffectiveDay:  16,
		DaysInMonth:   30,
	})
	// Old: 333,333.33 × 15 = 5,000,000
	// New: 400,000.00 × 15 = 6,000,000
	// Total: 11,000,000
	// Retro: 11M - 10M = 1,000,000
	assert.Equal(t, 15, result.DaysAtOld)
	assert.Equal(t, 15, result.DaysAtNew)
	assert.Equal(t, 11_000_000.0, result.TotalAmount)
	assert.Equal(t, 1_000_000.0, result.RetroAmount)
}

func TestRetroactivePay_BeginningOfMonth(t *testing.T) {
	// Change on day 1 = full month at new rate
	result := CalcRetroactivePay(RetroactivePayInput{
		OldBaseSalary: 10_000_000,
		NewBaseSalary: 12_000_000,
		EffectiveDay:  1,
		DaysInMonth:   30,
	})
	assert.Equal(t, 0, result.DaysAtOld)
	assert.Equal(t, 30, result.DaysAtNew)
	assert.Equal(t, 12_000_000.0, result.TotalAmount)
	assert.Equal(t, 2_000_000.0, result.RetroAmount)
}

func TestRetroactivePay_EndOfMonth(t *testing.T) {
	// Change on day 31 = full month at old rate
	result := CalcRetroactivePay(RetroactivePayInput{
		OldBaseSalary: 10_000_000,
		NewBaseSalary: 12_000_000,
		EffectiveDay:  31,
		DaysInMonth:   30,
	})
	assert.Equal(t, 30, result.DaysAtOld)
	assert.Equal(t, 0, result.DaysAtNew)
	assert.Equal(t, 10_000_000.0, result.TotalAmount)
	assert.Equal(t, 0.0, result.RetroAmount)
}

func TestRetroactivePay_NoChange(t *testing.T) {
	result := CalcRetroactivePay(RetroactivePayInput{
		OldBaseSalary: 10_000_000,
		NewBaseSalary: 10_000_000,
		EffectiveDay:  16,
		DaysInMonth:   30,
	})
	assert.Equal(t, 0.0, result.RetroAmount)
	assert.Equal(t, 10_000_000.0, result.TotalAmount)
}

func TestRetroactivePay_February(t *testing.T) {
	// 28-day month
	result := CalcRetroactivePay(RetroactivePayInput{
		OldBaseSalary: 14_000_000,
		NewBaseSalary: 16_000_000,
		EffectiveDay:  15,
		DaysInMonth:   28,
	})
	assert.Equal(t, 14, result.DaysAtOld)
	assert.Equal(t, 14, result.DaysAtNew)
	assert.True(t, result.RetroAmount > 0, "should have positive retro")
}

func TestRetroactivePay_SalaryDecrease(t *testing.T) {
	// Pay cut: 12M → 10M on day 16
	result := CalcRetroactivePay(RetroactivePayInput{
		OldBaseSalary: 12_000_000,
		NewBaseSalary: 10_000_000,
		EffectiveDay:  16,
		DaysInMonth:   30,
	})
	assert.True(t, result.RetroAmount < 0, "pay cut should have negative retro")
	assert.Equal(t, -1_000_000.0, result.RetroAmount)
}

func TestRetroactivePay_Defaults(t *testing.T) {
	// Zero daysInMonth defaults to 30
	result := CalcRetroactivePay(RetroactivePayInput{
		OldBaseSalary: 10_000_000,
		NewBaseSalary: 12_000_000,
		EffectiveDay:  16,
		DaysInMonth:   0,
	})
	assert.Equal(t, 30, result.DaysAtOld+result.DaysAtNew)
}
