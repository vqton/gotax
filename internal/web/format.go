package web

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

func formatVND(v float64) string {
	s := fmt.Sprintf("%.0f", v)
	out := make([]byte, 0, len(s)+len(s)/3)
	for i, ch := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, byte(ch))
	}
	return string(out)
}

func fmtDate(s any) string {
	switch t := s.(type) {
	case time.Time:
		if t.IsZero() {
			return "—"
		}
		return t.Format("02/01/2006")
	case *time.Time:
		if t == nil || t.IsZero() {
			return "—"
		}
		return t.Format("02/01/2006")
	case string:
		if t == "" {
			return "—"
		}
		if d, err := time.Parse("2006-01-02", t); err == nil {
			return d.Format("02/01/2006")
		}
		if d, err := time.Parse(time.RFC3339, t); err == nil {
			return d.Format("02/01/2006")
		}
		return t
	default:
		return "—"
	}
}

func fmtDateTime(s any) string {
	switch t := s.(type) {
	case time.Time:
		if t.IsZero() {
			return "—"
		}
		return t.Format("02/01/2006 15:04")
	case *time.Time:
		if t == nil || t.IsZero() {
			return "—"
		}
		return t.Format("02/01/2006 15:04")
	case string:
		if t == "" {
			return "—"
		}
		if d, err := time.Parse(time.RFC3339, t); err == nil {
			return d.Format("02/01/2006 15:04")
		}
		return t
	default:
		return "—"
	}
}

// statusBadge renders a colored pill for a journal entry status. Returns
// template.HTML (trusted — status value is matched against a fixed map).
func statusBadge(status string) template.HTML {
	cfg := map[string][2]string{
		"draft":     {"bg-gray-100 text-gray-700", "Nháp"},
		"pending":   {"bg-amber-50 text-amber-700", "Chờ duyệt"},
		"reviewed":  {"bg-blue-50 text-blue-700", "Đã rà soát"},
		"approved":  {"bg-indigo-50 text-indigo-700", "Đã duyệt"},
		"posted":    {"bg-green-50 text-green-700", "Đã ghi sổ"},
		"cancelled": {"bg-red-50 text-red-700", "Hủy"},
	}
	classes, label := "bg-gray-100 text-gray-700", status
	if v, ok := cfg[strings.ToLower(status)]; ok {
		classes, label = v[0], v[1]
	}
	return template.HTML(`<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ` + classes + `">` + label + `</span>`)
}
