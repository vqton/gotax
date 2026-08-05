package payroll

import (
	"strings"
	"testing"

	"gotax/internal/domain"
)

func TestParseTimekeepingCSV_Success(t *testing.T) {
	csv := `employee_code,date,clock_in,clock_out,ot_hours,night_hours,leave_type,notes
EMP001,2026-01-15,08:00,17:00,2,1,,overtime
EMP002,2026-01-15,08:00,17:00,0,0,ANNUAL,sick`
	records, err := ParseTimekeepingCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].EmployeeID != "EMP001" {
		t.Errorf("expected EMP001, got %s", records[0].EmployeeID)
	}
	if records[0].OTHours != 2 {
		t.Errorf("expected OT 2h, got %.1f", records[0].OTHours)
	}
	if records[0].NightHours != 1 {
		t.Errorf("expected night 1h, got %.1f", records[0].NightHours)
	}
	if records[1].LeaveType != domain.LeaveAnnual {
		t.Errorf("expected ANNUAL, got %s", records[1].LeaveType)
	}
}

func TestParseTimekeepingCSV_Empty(t *testing.T) {
	csv := `employee_code,date,clock_in,clock_out,ot_hours,night_hours,leave_type,notes`
	records, err := ParseTimekeepingCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected 0 records, got %d", len(records))
	}
}
