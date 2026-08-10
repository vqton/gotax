package service

import (
	"context"
	"fmt"
	"testing"

	"gotax/internal/domain"
)

type mockBankImportRepo struct {
	imports    map[string]*domain.BankImport
	transactions map[string][]domain.BankTransaction
}

func newMockBankImportRepo() *mockBankImportRepo {
	return &mockBankImportRepo{
		imports:      make(map[string]*domain.BankImport),
		transactions: make(map[string][]domain.BankTransaction),
	}
}

func (r *mockBankImportRepo) CreateImport(_ context.Context, imp *domain.BankImport) error {
	r.imports[imp.ID] = imp
	return nil
}

func (r *mockBankImportRepo) GetImport(_ context.Context, id string) (*domain.BankImport, error) {
	imp, ok := r.imports[id]
	if !ok {
		return nil, fmt.Errorf("import not found")
	}
	return imp, nil
}

func (r *mockBankImportRepo) ListImports(_ context.Context, companyID string) ([]domain.BankImport, error) {
	var out []domain.BankImport
	for _, imp := range r.imports {
		if imp.CompanyID == companyID {
			out = append(out, *imp)
		}
	}
	return out, nil
}

func (r *mockBankImportRepo) UpdateImportStatus(_ context.Context, id string, status string) error {
	imp, ok := r.imports[id]
	if !ok {
		return fmt.Errorf("import not found")
	}
	imp.Status = status
	return nil
}

func (r *mockBankImportRepo) CreateTransactions(_ context.Context, txns []domain.BankTransaction) error {
	for _, txn := range txns {
		r.transactions[txn.Reference] = append(r.transactions[txn.Reference], txn)
	}
	return nil
}

func (r *mockBankImportRepo) GetTransactionsByImport(_ context.Context, importID string) ([]domain.BankTransaction, error) {
	return r.transactions[importID], nil
}

func TestBankImportService_Import(t *testing.T) {
	repo := newMockBankImportRepo()
	svc := NewBankImportService(repo)
	csvData := []byte(`Date,TransactionID,Description,Debit,Credit,Reference
2026-08-10,TXN001,Payment,1000000,,REF001
2026-08-11,TXN002,Receipt,,500000,REF002`)
	imp, err := svc.Import(context.Background(), "CMP001", "VCB", "test.csv", csvData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if imp.Status != "completed" {
		t.Errorf("expected status completed, got %s", imp.Status)
	}
	if imp.TransactionCount != 2 {
		t.Errorf("expected 2 transactions, got %d", imp.TransactionCount)
	}
}

func TestBankImportService_InvalidBankCode(t *testing.T) {
	repo := newMockBankImportRepo()
	svc := NewBankImportService(repo)
	csvData := []byte(`Date,TransactionID,Description,Debit,Credit,Reference`)
	_, err := svc.Import(context.Background(), "CMP001", "INVALID", "test.csv", csvData)
	if err == nil {
		t.Error("expected error for invalid bank code")
	}
}

func TestBankImportService_ListImports(t *testing.T) {
	repo := newMockBankImportRepo()
	svc := NewBankImportService(repo)
	repo.imports["IMP1"] = &domain.BankImport{ID: "IMP1", CompanyID: "CMP001", Status: "completed"}
	repo.imports["IMP2"] = &domain.BankImport{ID: "IMP2", CompanyID: "CMP001", Status: "pending"}
	repo.imports["IMP3"] = &domain.BankImport{ID: "IMP3", CompanyID: "CMP002", Status: "completed"}
	imports, err := svc.ListImports(context.Background(), "CMP001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(imports) != 2 {
		t.Errorf("expected 2 imports, got %d", len(imports))
	}
}
