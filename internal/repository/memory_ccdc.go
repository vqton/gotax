package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"gotax/internal/domain"
)

// ─── Category Repo ───────────────────────────────────────────────────

type MemoryCCDCCategoryRepo struct {
	mu         sync.RWMutex
	categories map[string]*domain.ToolEquipmentCategory
	counter    int
}

func NewMemoryCCDCCategoryRepo() *MemoryCCDCCategoryRepo {
	return &MemoryCCDCCategoryRepo{
		categories: make(map[string]*domain.ToolEquipmentCategory),
	}
}

func (r *MemoryCCDCCategoryRepo) Create(_ context.Context, c *domain.ToolEquipmentCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	cp := *c
	cp.ID = fmt.Sprintf("CCDC-CAT-%d-%d", time.Now().UnixNano(), r.counter)
	r.categories[cp.ID] = &cp
	c.ID = cp.ID
	return nil
}

func (r *MemoryCCDCCategoryRepo) GetByID(_ context.Context, id string) (*domain.ToolEquipmentCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.categories[id]
	if !ok {
		return nil, domain.ErrCCDCCategoryNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemoryCCDCCategoryRepo) List(_ context.Context, companyID string) ([]domain.ToolEquipmentCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ToolEquipmentCategory
	for _, c := range r.categories {
		if c.CompanyID == companyID {
			out = append(out, *c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, nil
}

func (r *MemoryCCDCCategoryRepo) Update(_ context.Context, c *domain.ToolEquipmentCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.categories[c.ID]; !ok {
		return domain.ErrCCDCCategoryNotFound
	}
	cp := *c
	r.categories[c.ID] = &cp
	return nil
}

func (r *MemoryCCDCCategoryRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.categories[id]; !ok {
		return domain.ErrCCDCCategoryNotFound
	}
	delete(r.categories, id)
	return nil
}

// ─── Item Repo ───────────────────────────────────────────────────────

type MemoryToolEquipmentRepo struct {
	mu         sync.RWMutex
	items      map[string]*domain.ToolEquipment
	itemByCode map[string]string
	counter    int
}

func NewMemoryCCDCItemRepo() *MemoryToolEquipmentRepo {
	return &MemoryToolEquipmentRepo{
		items:      make(map[string]*domain.ToolEquipment),
		itemByCode: make(map[string]string),
	}
}

func (r *MemoryToolEquipmentRepo) Create(_ context.Context, t *domain.ToolEquipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := t.CompanyID + ":" + t.Code
	if _, ok := r.itemByCode[key]; ok {
		return domain.ErrCCDCCodeExists
	}
	r.counter++
	cp := *t
	cp.ID = fmt.Sprintf("CCDC-ITEM-%d-%d", time.Now().UnixNano(), r.counter)
	r.items[cp.ID] = &cp
	r.itemByCode[key] = cp.ID
	t.ID = cp.ID
	return nil
}

func (r *MemoryToolEquipmentRepo) GetByID(_ context.Context, id string) (*domain.ToolEquipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	if !ok {
		return nil, domain.ErrCCDCNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *MemoryToolEquipmentRepo) List(_ context.Context, companyID string) ([]domain.ToolEquipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ToolEquipment
	for _, t := range r.items {
		if t.CompanyID == companyID {
			out = append(out, *t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, nil
}

func (r *MemoryToolEquipmentRepo) Update(_ context.Context, t *domain.ToolEquipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[t.ID]; !ok {
		return domain.ErrCCDCNotFound
	}
	cp := *t
	r.items[t.ID] = &cp
	return nil
}

func (r *MemoryToolEquipmentRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok {
		return domain.ErrCCDCNotFound
	}
	key := item.CompanyID + ":" + item.Code
	delete(r.itemByCode, key)
	delete(r.items, id)
	return nil
}

func (r *MemoryToolEquipmentRepo) GetByCode(_ context.Context, companyID, code string) (*domain.ToolEquipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := companyID + ":" + code
	id, ok := r.itemByCode[key]
	if !ok {
		return nil, domain.ErrCCDCNotFound
	}
	t, ok := r.items[id]
	if !ok {
		return nil, domain.ErrCCDCNotFound
	}
	cp := *t
	return &cp, nil
}
