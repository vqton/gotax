package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gotax/internal/domain"
)

// ─── Contract ──────────────────────────────────────────────────────

type MemoryContractRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.Contract
}

func NewMemoryContractRepo() *MemoryContractRepo {
	return &MemoryContractRepo{data: make(map[string]*domain.Contract)}
}

func (r *MemoryContractRepo) Create(_ context.Context, c *domain.Contract) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.CompanyID == c.CompanyID && existing.Code == c.Code {
			return domain.ErrContractExists
		}
	}
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	if c.ID == "" {
		c.ID = fmt.Sprintf("CON-%d", time.Now().UnixNano())
	}
	if c.Status == "" {
		c.Status = "draft"
	}
	c.CreatedAt = now
	c.UpdatedAt = now
	cp := *c
	r.data[c.ID] = &cp
	return nil
}

func (r *MemoryContractRepo) GetByID(_ context.Context, id string) (*domain.Contract, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.data[id]
	if !ok {
		return nil, domain.ErrContractNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemoryContractRepo) GetByCode(_ context.Context, companyID, code string) (*domain.Contract, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.data {
		if c.CompanyID == companyID && c.Code == code {
			cp := *c
			return &cp, nil
		}
	}
	return nil, domain.ErrContractNotFound
}

func (r *MemoryContractRepo) List(_ context.Context, companyID string) ([]domain.Contract, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Contract
	for _, c := range r.data {
		if c.CompanyID == companyID {
			cp := *c
			out = append(out, cp)
		}
	}
	return out, nil
}

func (r *MemoryContractRepo) Update(_ context.Context, c *domain.Contract) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[c.ID]; !ok {
		return domain.ErrContractNotFound
	}
	c.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	cp := *c
	r.data[c.ID] = &cp
	return nil
}

func (r *MemoryContractRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrContractNotFound
	}
	delete(r.data, id)
	return nil
}

// ─── ContractPayment ───────────────────────────────────────────────

type MemoryContractPaymentRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.ContractPayment
}

func NewMemoryContractPaymentRepo() *MemoryContractPaymentRepo {
	return &MemoryContractPaymentRepo{data: make(map[string]*domain.ContractPayment)}
}

func (r *MemoryContractPaymentRepo) Create(_ context.Context, p *domain.ContractPayment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	if p.ID == "" {
		p.ID = fmt.Sprintf("CP-%d", time.Now().UnixNano())
	}
	p.CreatedAt = now
	cp := *p
	r.data[p.ID] = &cp
	return nil
}

func (r *MemoryContractPaymentRepo) ListByContract(_ context.Context, contractID string) ([]domain.ContractPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ContractPayment
	for _, p := range r.data {
		if p.ContractID == contractID {
			cp := *p
			out = append(out, cp)
		}
	}
	return out, nil
}

func (r *MemoryContractPaymentRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrContractNotFound
	}
	delete(r.data, id)
	return nil
}
