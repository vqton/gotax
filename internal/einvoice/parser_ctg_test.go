package einvoice

import (
	"testing"
)

func TestCTGParser_GetBankCode(t *testing.T) {
	parser := NewCTGParser()
	if parser.GetBankCode() != "CTG" {
		t.Errorf("expected CTG, got %s", parser.GetBankCode())
	}
}

func TestCTGParser_Validate(t *testing.T) {
	parser := NewCTGParser()
	validCSV := []byte(`Date,TransactionID,Description,Debit,Credit,Reference
2026-08-10,TXN001,Payment,1000000,,REF001
2026-08-11,TXN002,Receipt,,500000,REF002`)
	if err := parser.Validate(validCSV); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	invalidCSV := []byte(`Date,TransactionID`)
	if err := parser.Validate(invalidCSV); err == nil {
		t.Error("expected error for invalid CSV")
	}
}

func TestCTGParser_Parse(t *testing.T) {
	parser := NewCTGParser()
	csvData := []byte(`Date,TransactionID,Description,Debit,Credit,Reference
2026-08-10,TXN001,Payment,1000000,,REF001
2026-08-11,TXN002,Receipt,,500000,REF002`)
	txns, err := parser.Parse(csvData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txns) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(txns))
	}
	if txns[0].TransactionDate != "2026-08-10" {
		t.Errorf("expected date 2026-08-10, got %s", txns[0].TransactionDate)
	}
	if txns[0].TransactionID != "TXN001" {
		t.Errorf("expected TXN001, got %s", txns[0].TransactionID)
	}
	if txns[0].Debit != 1000000 {
		t.Errorf("expected debit 1000000, got %f", txns[0].Debit)
	}
	if txns[0].Credit != 0 {
		t.Errorf("expected credit 0, got %f", txns[0].Credit)
	}
}
