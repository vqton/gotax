package payroll

import (
	"strings"
	"testing"

	"gotax/internal/domain"
)

func TestGenerateD02TS(t *testing.T) {
	employees := []domain.EmployeePayrollInfo{
		{EmployeeID: "EMP001", BaseSalary: 15000000, InsuranceBaseSalary: 15000000, ContractStartDate: "2024-01-01"},
	}
	runs := []domain.PayrollRun{
		{EmployeeID: "EMP001", SIDeduction: 1200000, HIDeduction: 225000, UIDeduction: 150000, EmployerSI: 2625000, EmployerHI: 450000, EmployerUI: 150000},
	}
	xmlBytes, err := GenerateD02TS("Công ty ABC", "0123456789", 2026, 1, employees, runs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(xmlBytes)
	if !strings.Contains(s, "D02TS") {
		t.Error("expected D02TS element")
	}
	if !strings.Contains(s, "EMP001") {
		t.Error("expected employee code")
	}
}

func TestGenerate05KKTNCN(t *testing.T) {
	employees := []domain.EmployeePayrollInfo{
		{EmployeeID: "EMP001", BaseSalary: 15000000, InsuranceBaseSalary: 15000000, ContractStartDate: "2024-01-01"},
	}
	runs := []domain.PayrollRun{
		{EmployeeID: "EMP001", GrossSalary: 15000000, SIDeduction: 1200000, HIDeduction: 225000, UIDeduction: 150000, TradeUnionDues: 150000, PITAmount: 500000},
	}
	xmlBytes, err := Generate05KKTNCN("0123456789", "Công ty ABC", "Q1/2026", employees, runs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(xmlBytes)
	if !strings.Contains(s, "KK_TNCN") {
		t.Error("expected KK_TNCN element")
	}
	if !strings.Contains(s, "500000") {
		t.Error("expected PIT amount")
	}
}

func TestGenerateTK3TS(t *testing.T) {
	employees := []domain.EmployeePayrollInfo{
		{EmployeeID: "EMP001", BaseSalary: 15000000, ContractStartDate: "2024-01-01"},
	}
	xmlBytes, err := GenerateTK3TS("Công ty ABC", "0123456789", "123 Đường ABC", employees)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(xmlBytes)
	if !strings.Contains(s, "TK3TS") {
		t.Error("expected TK3TS element")
	}
	if !strings.Contains(s, "EMP001") {
		t.Error("expected employee code")
	}
}

func TestGenerateD02TS_EmptyEmployees(t *testing.T) {
	xmlBytes, err := GenerateD02TS("Công ty ABC", "0123456789", 2026, 1, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(xmlBytes), "D02TS") {
		t.Error("expected valid XML even with no employees")
	}
}

func TestParseTimekeepingCSV(t *testing.T) {
	csv := `employee_code,date,clock_in,clock_out,ot_hours,night_hours,leave_type,notes
EMP001,2026-01-15,08:00,17:00,2,1,,overtime
EMP002,2026-01-15,08:00,17:00,0,0,ANNUAL,annual leave`
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
}

func TestParseTimekeepingCSV_MissingColumn(t *testing.T) {
	csv := `employee_code,clock_in,clock_out
EMP001,08:00,17:00`
	_, err := ParseTimekeepingCSV(strings.NewReader(csv))
	if err == nil {
		t.Fatal("expected error for missing date column")
	}
}
