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
	cls := map[string]string{
		"draft":     "badge-draft",
		"reviewing": "badge-reviewing",
		"approved":  "badge-approved",
		"posted":    "badge-posted",
		"cancelled": "badge-cancelled",
	}
	labels := map[string]string{
		"draft":     "Nháp",
		"reviewing": "Chờ duyệt",
		"approved":  "Đã duyệt",
		"posted":    "Đã ghi sổ",
		"cancelled": "Hủy",
	}
	key := strings.ToLower(status)
	class, label := cls[key], labels[key]
	if class == "" {
		class = "badge-default"
	}
	if label == "" {
		label = status
	}
	return template.HTML(`<span class="badge ` + class + `">` + label + `</span>`)
}

// accountStatusBadge renders a colored pill for an account status
// (ACTIVE/FROZEN/INACTIVE). Fixed map — safe to emit as trusted HTML.
func accountStatusBadge(status string) template.HTML {
	cls := map[string]string{
		"active":   "badge-success",
		"frozen":   "badge-warning",
		"inactive": "badge-default",
	}
	labels := map[string]string{
		"active":   "Hoạt động",
		"frozen":   "Ngừng TD",
		"inactive": "Không dùng",
	}
	key := strings.ToLower(status)
	class, label := cls[key], labels[key]
	if class == "" {
		class = "badge-default"
	}
	if label == "" {
		label = status
	}
	return template.HTML(`<span class="badge ` + class + `">` + label + `</span>`)
}
