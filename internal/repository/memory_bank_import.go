package repository

import (
	"context"
	"sync"
	"time"

	"gotax/internal/domain"
)

type MemoryBankImportRepo struct {
	mu           sync.RWMutex
	imports      map[string]*domain.BankImport
	transactions map[string][]domain.BankTransaction
}

func NewMemoryBankImportRepo() *MemoryBankImportRepo {
	return &MemoryBankImportRepo{
		imports:      make(map[string]*domain.BankImport),
		transactions: make(map[string][]domain.BankTransaction),
	}
}

func (r *MemoryBankImportRepo) CreateImport(_ context.Context, imp *domain.BankImport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.imports[imp.ID] = imp
	return nil
}

func (r *MemoryBankImportRepo) GetImport(_ context.Context, id string) (*domain.BankImport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	imp, ok := r.imports[id]
	if !ok {
		return nil, domain.ErrImportNotFound
	}
	out := *imp
	return &out, nil
}

func (r *MemoryBankImportRepo) ListImports(_ context.Context, companyID string) ([]domain.BankImport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.BankImport
	for _, imp := range r.imports {
		if imp.CompanyID == companyID {
			out = append(out, *imp)
		}
	}
	return out, nil
}

func (r *MemoryBankImportRepo) UpdateImportStatus(_ context.Context, id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	imp, ok := r.imports[id]
	if !ok {
		return domain.ErrImportNotFound
	}
	imp.Status = status
	return nil
}

func (r *MemoryBankImportRepo) CreateTransactions(_ context.Context, txns []domain.BankTransaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, txn := range txns {
		r.transactions[txn.Reference] = append(r.transactions[txn.Reference], txn)
	}
	return nil
}

func (r *MemoryBankImportRepo) GetTransactionsByImport(_ context.Context, importID string) ([]domain.BankTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.transactions[importID], nil
}

// Ensure time.Time is imported for potential future use.
var _ = time.Now
