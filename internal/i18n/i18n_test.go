package i18n

import (
	"testing"
)

func TestNew_LoadsAllFiles(t *testing.T) {
	l, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if l == nil {
		t.Fatal("New() returned nil")
	}
}

func TestLocalize_English(t *testing.T) {
	l := MustNew()
	localizer := l.ForLocale("en")
	msg := Localize(localizer, "ob_not_found")
	if msg != "Opening balance not found" {
		t.Errorf("got %q, want %q", msg, "Opening balance not found")
	}
}

func TestLocalize_Vietnamese(t *testing.T) {
	l := MustNew()
	localizer := l.ForLocale("vi")
	msg := Localize(localizer, "ob_not_found")
	if msg != "Không tìm thấy số dư đầu kỳ" {
		t.Errorf("got %q, want Vietnamese", msg)
	}
}

func TestLocalize_Fallback(t *testing.T) {
	l := MustNew()
	localizer := l.ForLocale("en")
	msg := Localize(localizer, "nonexistent_id")
	if msg != "nonexistent_id" {
		t.Errorf("got %q, want %q", msg, "nonexistent_id")
	}
}

func TestLocalize_NilLocalizer(t *testing.T) {
	msg := Localize(nil, "ob_not_found")
	if msg != "ob_not_found" {
		t.Errorf("got %q, want %q", msg, "ob_not_found")
	}
}

func TestLocalize_Template(t *testing.T) {
	l := MustNew()
	localizer := l.ForLocale("en")
	msg := Localize(localizer, "import_success", map[string]any{"Success": 5, "Total": 10})
	if msg != "Imported 5 of 10 opening balances" {
		t.Errorf("got %q", msg)
	}
}

func TestForLocale_FrenchFallsbackToEnglish(t *testing.T) {
	l := MustNew()
	localizer := l.ForLocale("fr")
	msg := Localize(localizer, "ob_not_found")
	if msg != "Opening balance not found" {
		t.Errorf("got %q, want English fallback", msg)
	}
}
