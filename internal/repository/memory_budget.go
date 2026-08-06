package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gotax/internal/domain"
)

type MemoryBudgetRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.Budget
}

func NewMemoryBudgetRepo() *MemoryBudgetRepo {
	return &MemoryBudgetRepo{data: make(map[string]*domain.Budget)}
}

func (r *MemoryBudgetRepo) Create(_ context.Context, b *domain.Budget) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.CompanyID == b.CompanyID &&
			existing.AccountCode == b.AccountCode &&
			existing.PeriodYear == b.PeriodYear &&
			existing.PeriodMonth == b.PeriodMonth {
			return domain.ErrBudgetExists
		}
	}
	cp := *b
	r.data[b.ID] = &cp
	return nil
}

func (r *MemoryBudgetRepo) GetByID(_ context.Context, id string) (*domain.Budget, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.data[id]
	if !ok {
		return nil, domain.ErrBudgetNotFound
	}
	cp := *b
	return &cp, nil
}

func (r *MemoryBudgetRepo) List(_ context.Context, companyID string, year int) ([]domain.Budget, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Budget
	for _, b := range r.data {
		if b.CompanyID != companyID {
			continue
		}
		if year > 0 && b.PeriodYear != year {
			continue
		}
		cp := *b
		out = append(out, cp)
	}
	return out, nil
}

func (r *MemoryBudgetRepo) Update(_ context.Context, b *domain.Budget) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[b.ID]; !ok {
		return domain.ErrBudgetNotFound
	}
	cp := *b
	r.data[b.ID] = &cp
	return nil
}

func (r *MemoryBudgetRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrBudgetNotFound
	}
	delete(r.data, id)
	return nil
}

func (r *MemoryBudgetRepo) Upsert(_ context.Context, b *domain.Budget) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.CompanyID == b.CompanyID &&
			existing.AccountCode == b.AccountCode &&
			existing.PeriodYear == b.PeriodYear &&
			existing.PeriodMonth == b.PeriodMonth {
			existing.Budgeted = b.Budgeted
			existing.Actual = b.Actual
			existing.Variance = b.Variance
			existing.Notes = b.Notes
			existing.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")
			b.ID = existing.ID
			return nil
		}
	}
	if b.ID == "" {
		b.ID = fmt.Sprintf("BUD-%d", time.Now().UnixNano())
	}
	cp := *b
	r.data[b.ID] = &cp
	return nil
}
