package repository

import (
	"time"
)

// GORM → Domain helpers

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func safeInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func safeBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func safeTimeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateOnly)
}

func safeTimePtrStr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(time.DateOnly)
}

func safeTimePtrRFC3339(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// Domain → GORM helpers

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func parseDate(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseDateTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func timePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// nullStrG keeps existing definition (string → *string)
func nullStrG(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
