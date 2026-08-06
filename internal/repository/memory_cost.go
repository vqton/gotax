package repository

import (
	"context"
	"sync"

	"gotax/internal/domain"
)

type MemoryCostCenterRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.CostCenter
}

func NewMemoryCostCenterRepo() *MemoryCostCenterRepo {
	return &MemoryCostCenterRepo{data: make(map[string]*domain.CostCenter)}
}

func (r *MemoryCostCenterRepo) Create(_ context.Context, cc *domain.CostCenter) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.CompanyID == cc.CompanyID && existing.Code == cc.Code {
			return domain.ErrCostCenterCodeExists
		}
	}
	cp := *cc
	r.data[cc.ID] = &cp
	return nil
}

func (r *MemoryCostCenterRepo) GetByID(_ context.Context, id string) (*domain.CostCenter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cc, ok := r.data[id]
	if !ok {
		return nil, domain.ErrCostCenterNotFound
	}
	cp := *cc
	return &cp, nil
}

func (r *MemoryCostCenterRepo) List(_ context.Context, companyID string) ([]domain.CostCenter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CostCenter
	for _, cc := range r.data {
		if cc.CompanyID != companyID {
			continue
		}
		cp := *cc
		out = append(out, cp)
	}
	return out, nil
}

func (r *MemoryCostCenterRepo) Update(_ context.Context, cc *domain.CostCenter) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[cc.ID]; !ok {
		return domain.ErrCostCenterNotFound
	}
	cp := *cc
	r.data[cc.ID] = &cp
	return nil
}

func (r *MemoryCostCenterRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrCostCenterNotFound
	}
	delete(r.data, id)
	return nil
}

func (r *MemoryCostCenterRepo) GetByCode(_ context.Context, companyID, code string) (*domain.CostCenter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, cc := range r.data {
		if cc.CompanyID == companyID && cc.Code == code {
			cp := *cc
			return &cp, nil
		}
	}
	return nil, domain.ErrCostCenterNotFound
}
