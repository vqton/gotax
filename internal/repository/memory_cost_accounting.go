package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gotax/internal/domain"
)

// ─── CostObject ────────────────────────────────────────────────────────────

type MemoryCostObjectRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.CostObject
}

func NewMemoryCostObjectRepo() *MemoryCostObjectRepo {
	return &MemoryCostObjectRepo{data: make(map[string]*domain.CostObject)}
}

func (r *MemoryCostObjectRepo) Create(_ context.Context, co *domain.CostObject) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.CompanyID == co.CompanyID && existing.Code == co.Code {
			return domain.ErrCostObjectCodeExists
		}
	}
	cp := *co
	r.data[co.ID] = &cp
	return nil
}

func (r *MemoryCostObjectRepo) GetByID(_ context.Context, id string) (*domain.CostObject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	co, ok := r.data[id]
	if !ok {
		return nil, domain.ErrCostObjectNotFound
	}
	cp := *co
	return &cp, nil
}

func (r *MemoryCostObjectRepo) GetByCode(_ context.Context, companyID, code string) (*domain.CostObject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, co := range r.data {
		if co.CompanyID == companyID && co.Code == code {
			cp := *co
			return &cp, nil
		}
	}
	return nil, domain.ErrCostObjectNotFound
}

func (r *MemoryCostObjectRepo) List(_ context.Context, companyID string) ([]domain.CostObject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CostObject
	for _, co := range r.data {
		if co.CompanyID != companyID {
			continue
		}
		cp := *co
		out = append(out, cp)
	}
	return out, nil
}

func (r *MemoryCostObjectRepo) Update(_ context.Context, co *domain.CostObject) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[co.ID]; !ok {
		return domain.ErrCostObjectNotFound
	}
	cp := *co
	r.data[co.ID] = &cp
	return nil
}

func (r *MemoryCostObjectRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrCostObjectNotFound
	}
	delete(r.data, id)
	return nil
}

// ─── CostPool ──────────────────────────────────────────────────────────────

type MemoryCostPoolRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.CostPool
}

func NewMemoryCostPoolRepo() *MemoryCostPoolRepo {
	return &MemoryCostPoolRepo{data: make(map[string]*domain.CostPool)}
}

func (r *MemoryCostPoolRepo) Create(_ context.Context, cp *domain.CostPool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := *cp
	r.data[cp.ID] = &c
	return nil
}

func (r *MemoryCostPoolRepo) GetByID(_ context.Context, id string) (*domain.CostPool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp, ok := r.data[id]
	if !ok {
		return nil, domain.ErrCostPoolNotFound
	}
	c := *cp
	return &c, nil
}

func (r *MemoryCostPoolRepo) ListByPeriod(_ context.Context, companyID, periodID string) ([]domain.CostPool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CostPool
	for _, cp := range r.data {
		if cp.CompanyID != companyID || cp.PeriodID != periodID {
			continue
		}
		c := *cp
		out = append(out, c)
	}
	return out, nil
}

func (r *MemoryCostPoolRepo) Update(_ context.Context, cp *domain.CostPool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[cp.ID]; !ok {
		return domain.ErrCostPoolNotFound
	}
	c := *cp
	r.data[cp.ID] = &c
	return nil
}

func (r *MemoryCostPoolRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrCostPoolNotFound
	}
	delete(r.data, id)
	return nil
}

// ─── CostPoolLine ──────────────────────────────────────────────────────────

type MemoryCostPoolLineRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.CostPoolLine
}

func NewMemoryCostPoolLineRepo() *MemoryCostPoolLineRepo {
	return &MemoryCostPoolLineRepo{data: make(map[string]*domain.CostPoolLine)}
}

func (r *MemoryCostPoolLineRepo) Create(_ context.Context, line *domain.CostPoolLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := *line
	r.data[line.ID] = &c
	return nil
}

func (r *MemoryCostPoolLineRepo) ListByPool(_ context.Context, poolID string) ([]domain.CostPoolLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CostPoolLine
	for _, l := range r.data {
		if l.PoolID != poolID {
			continue
		}
		cp := *l
		out = append(out, cp)
	}
	return out, nil
}

func (r *MemoryCostPoolLineRepo) DeleteByPool(_ context.Context, poolID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, l := range r.data {
		if l.PoolID == poolID {
			delete(r.data, id)
		}
	}
	return nil
}

// ─── CostingPeriod ─────────────────────────────────────────────────────────

type MemoryCostingPeriodRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.CostingPeriod
}

func NewMemoryCostingPeriodRepo() *MemoryCostingPeriodRepo {
	return &MemoryCostingPeriodRepo{data: make(map[string]*domain.CostingPeriod)}
}

func (r *MemoryCostingPeriodRepo) Create(_ context.Context, cp *domain.CostingPeriod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.CompanyID == cp.CompanyID && existing.Year == cp.Year && existing.Month == cp.Month {
			return domain.ErrCostingPeriodExists
		}
	}
	c := *cp
	r.data[cp.ID] = &c
	return nil
}

func (r *MemoryCostingPeriodRepo) GetByID(_ context.Context, id string) (*domain.CostingPeriod, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp, ok := r.data[id]
	if !ok {
		return nil, domain.ErrCostingPeriodNotFound
	}
	c := *cp
	return &c, nil
}

func (r *MemoryCostingPeriodRepo) GetByYearMonth(_ context.Context, companyID string, year, month int) (*domain.CostingPeriod, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, cp := range r.data {
		if cp.CompanyID == companyID && cp.Year == year && cp.Month == month {
			c := *cp
			return &c, nil
		}
	}
	return nil, domain.ErrCostingPeriodNotFound
}

func (r *MemoryCostingPeriodRepo) List(_ context.Context, companyID string) ([]domain.CostingPeriod, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CostingPeriod
	for _, cp := range r.data {
		if cp.CompanyID != companyID {
			continue
		}
		c := *cp
		out = append(out, c)
	}
	return out, nil
}

func (r *MemoryCostingPeriodRepo) Update(_ context.Context, cp *domain.CostingPeriod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[cp.ID]; !ok {
		return domain.ErrCostingPeriodNotFound
	}
	c := *cp
	r.data[cp.ID] = &c
	return nil
}

// ─── CostingResult ─────────────────────────────────────────────────────────

type MemoryCostingResultRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.CostingResult
}

func NewMemoryCostingResultRepo() *MemoryCostingResultRepo {
	return &MemoryCostingResultRepo{data: make(map[string]*domain.CostingResult)}
}

func (r *MemoryCostingResultRepo) Create(_ context.Context, cr *domain.CostingResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := *cr
	r.data[cr.ID] = &c
	return nil
}

func (r *MemoryCostingResultRepo) GetByID(_ context.Context, id string) (*domain.CostingResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cr, ok := r.data[id]
	if !ok {
		return nil, domain.ErrCostingResultNotFound
	}
	c := *cr
	return &c, nil
}

func (r *MemoryCostingResultRepo) ListByPeriod(_ context.Context, companyID, periodID string) ([]domain.CostingResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CostingResult
	for _, cr := range r.data {
		if cr.CompanyID != companyID || cr.PeriodID != periodID {
			continue
		}
		c := *cr
		out = append(out, c)
	}
	return out, nil
}

func (r *MemoryCostingResultRepo) Update(_ context.Context, cr *domain.CostingResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[cr.ID]; !ok {
		return domain.ErrCostingResultNotFound
	}
	c := *cr
	r.data[cr.ID] = &c
	return nil
}

func (r *MemoryCostingResultRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrCostingResultNotFound
	}
	delete(r.data, id)
	return nil
}

// ─── CostingResultLine ─────────────────────────────────────────────────────

type MemoryCostingResultLineRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.CostingResultLine
}

func NewMemoryCostingResultLineRepo() *MemoryCostingResultLineRepo {
	return &MemoryCostingResultLineRepo{data: make(map[string]*domain.CostingResultLine)}
}

func (r *MemoryCostingResultLineRepo) Create(_ context.Context, line *domain.CostingResultLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := *line
	if c.ID == "" {
		c.ID = fmt.Sprintf("CRL-%d", time.Now().UnixNano())
	}
	r.data[c.ID] = &c
	line.ID = c.ID
	return nil
}

func (r *MemoryCostingResultLineRepo) ListByResult(_ context.Context, resultID string) ([]domain.CostingResultLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CostingResultLine
	for _, l := range r.data {
		if l.ResultID != resultID {
			continue
		}
		cp := *l
		out = append(out, cp)
	}
	return out, nil
}

func (r *MemoryCostingResultLineRepo) DeleteByResult(_ context.Context, resultID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, l := range r.data {
		if l.ResultID == resultID {
			delete(r.data, id)
		}
	}
	return nil
}
