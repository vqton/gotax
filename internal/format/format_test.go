package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatNumber_Vietnamese(t *testing.T) {
	f := Default()
	assert.Equal(t, "1.234.568", FormatNumber(1234567.89, f))
	assert.Equal(t, "0", FormatNumber(0, f))
	assert.Equal(t, "(1.234.568)", FormatNumber(-1234567.89, f))
}

func TestFormatNumber_Decimal(t *testing.T) {
	f := NumberFormat{
		ThousandsSeparator: ".",
		DecimalSeparator:    ",",
		DecimalPlaces:       2,
		NegativeDisplay:     "minus",
	}
	assert.Equal(t, "1.234.567,89", FormatNumber(1234567.89, f))
	assert.Equal(t, "-1.234.567,89", FormatNumber(-1234567.89, f))
}

func TestFormatNumber_US(t *testing.T) {
	f := NumberFormat{
		ThousandsSeparator: ",",
		DecimalSeparator:    ".",
		DecimalPlaces:       2,
		NegativeDisplay:     "minus",
	}
	assert.Equal(t, "1,234,567.89", FormatNumber(1234567.89, f))
}

func TestParseNumber_Vietnamese(t *testing.T) {
	f := Default()
	val, err := ParseNumber("1.234.568", f)
	require.NoError(t, err)
	assert.Equal(t, float64(1234568), val)

	val, err = ParseNumber("(1.234.568)", f)
	require.NoError(t, err)
	assert.Equal(t, float64(-1234568), val)
}

func TestParseNumber_Decimal(t *testing.T) {
	f := NumberFormat{
		ThousandsSeparator: ".",
		DecimalSeparator:    ",",
		DecimalPlaces:       2,
		NegativeDisplay:     "minus",
	}
	val, err := ParseNumber("1.234.567,89", f)
	require.NoError(t, err)
	assert.InDelta(t, 1234567.89, val, 0.01)
}
