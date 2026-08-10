package format

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// NumberFormat holds formatting options for numeric display.
type NumberFormat struct {
	ThousandsSeparator string // ".", ",", " "
	DecimalSeparator    string // ".", ","
	DecimalPlaces       int    // 0-4
	NegativeDisplay     string // "parentheses" or "minus"
}

// Default returns Vietnamese number format (1.234.567,89).
func Default() NumberFormat {
	return NumberFormat{
		ThousandsSeparator: ".",
		DecimalSeparator:    ",",
		DecimalPlaces:       0,
		NegativeDisplay:     "parentheses",
	}
}

// FormatNumber formats a float64 per the given NumberFormat.
func FormatNumber(value float64, f NumberFormat) string {
	negative := value < 0
	abs := math.Abs(value)

	// Round to decimal places
	pow := math.Pow(10, float64(f.DecimalPlaces))
	abs = math.Round(abs*pow) / pow

	// Split integer and decimal
	intPart := int64(abs)
	decPart := abs - float64(intPart)

	// Format integer with thousands separator
	intStr := formatIntWithSeparator(intPart, f.ThousandsSeparator)

	// Format decimal part
	var result string
	if f.DecimalPlaces > 0 {
		decStr := fmt.Sprintf("%.*f", f.DecimalPlaces, decPart)
		// Remove leading "0." from decStr
		if len(decStr) > 2 {
			decStr = decStr[2:]
		}
		result = intStr + f.DecimalSeparator + decStr
	} else {
		result = intStr
	}

	if negative {
		switch f.NegativeDisplay {
		case "parentheses":
			return "(" + result + ")"
		default:
			return "-" + result
		}
	}
	return result
}

func formatIntWithSeparator(n int64, sep string) string {
	if sep == "" {
		return strconv.FormatInt(n, 10)
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	var result []string
	for len(s) > 3 {
		result = append([]string{s[len(s)-3:]}, result...)
		s = s[:len(s)-3]
	}
	if s != "" {
		result = append([]string{s}, result...)
	}
	return sign + strings.Join(result, sep)
}

// ParseNumber parses a formatted number string back to float64.
func ParseNumber(s string, f NumberFormat) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	negative := false
	// Handle parentheses negative
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		negative = true
		s = s[1 : len(s)-1]
	}

	// Remove thousands separator
	s = strings.ReplaceAll(s, f.ThousandsSeparator, "")

	// Replace decimal separator with dot
	s = strings.ReplaceAll(s, f.DecimalSeparator, ".")

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}

	if negative {
		val = -val
	}
	return val, nil
}
