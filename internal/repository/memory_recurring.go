package repository

import (
	"context"
	"sync"
	"time"

	"gotax/internal/domain"
)

type MemoryRecurringEntryRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.RecurringEntry
}

func NewMemoryRecurringEntryRepo() *MemoryRecurringEntryRepo {
	return &MemoryRecurringEntryRepo{data: make(map[string]*domain.RecurringEntry)}
}

func (r *MemoryRecurringEntryRepo) Create(_ context.Context, entry *domain.RecurringEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *entry
	r.data[entry.ID] = &cp
	return nil
}

func (r *MemoryRecurringEntryRepo) GetByID(_ context.Context, id string) (*domain.RecurringEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.data[id]
	if !ok {
		return nil, domain.ErrRecurringNotFound
	}
	cp := *e
	return &cp, nil
}

func (r *MemoryRecurringEntryRepo) List(_ context.Context, companyID string) ([]domain.RecurringEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.RecurringEntry
	for _, e := range r.data {
		if e.CompanyID != companyID {
			continue
		}
		cp := *e
		out = append(out, cp)
	}
	return out, nil
}

func (r *MemoryRecurringEntryRepo) Update(_ context.Context, entry *domain.RecurringEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[entry.ID]; !ok {
		return domain.ErrRecurringNotFound
	}
	cp := *entry
	r.data[entry.ID] = &cp
	return nil
}

func (r *MemoryRecurringEntryRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrRecurringNotFound
	}
	delete(r.data, id)
	return nil
}

func (r *MemoryRecurringEntryRepo) UpdateNextRunDate(_ context.Context, id, nextDate string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.data[id]
	if !ok {
		return domain.ErrRecurringNotFound
	}
	cp := *e
	cp.NextRunDate = nextDate
	cp.UpdatedAt = time.Now()
	r.data[id] = &cp
	return nil
}

func (r *MemoryRecurringEntryRepo) GetDueEntries(_ context.Context, today string) ([]domain.RecurringEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.RecurringEntry
	for _, e := range r.data {
		if e.IsActive && e.NextRunDate != "" && e.NextRunDate <= today {
			cp := *e
			out = append(out, cp)
		}
	}
	return out, nil
}
