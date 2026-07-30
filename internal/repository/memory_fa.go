package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"gotax/internal/domain"
)

type MemoryFARepo struct {
	mu                sync.RWMutex
	categories        map[string]*domain.FixedAssetCategory
	catByCode         map[string]string // companyID:code -> id
	assets            map[string]*domain.FixedAsset
	assetByCode       map[string]string
	depreciation      map[string]*domain.DepreciationEntry
	deprByAsset       map[string][]string
	transactions      map[string]*domain.FixedAssetTransaction
	txnByAsset        map[string][]string
	allocations       map[string][]domain.FixedAssetAllocation
	inventoryPlans    map[string]*domain.FixedAssetInventoryPlan
	inventoryResults  map[string]*domain.FixedAssetInventoryResult
	resultByPlan      map[string][]string
	catCounter        int
	assetCounter      int
	deprCounter       int
	txnCounter        int
	allocCounter      int
	planCounter       int
	resultCounter     int
}

func NewMemoryFARepo() *MemoryFARepo {
	return &MemoryFARepo{
		categories:       make(map[string]*domain.FixedAssetCategory),
		catByCode:        make(map[string]string),
		assets:           make(map[string]*domain.FixedAsset),
		assetByCode:      make(map[string]string),
		depreciation:     make(map[string]*domain.DepreciationEntry),
		deprByAsset:      make(map[string][]string),
		transactions:     make(map[string]*domain.FixedAssetTransaction),
		txnByAsset:       make(map[string][]string),
		allocations:      make(map[string][]domain.FixedAssetAllocation),
		inventoryPlans:   make(map[string]*domain.FixedAssetInventoryPlan),
		inventoryResults: make(map[string]*domain.FixedAssetInventoryResult),
		resultByPlan:     make(map[string][]string),
	}
}

func faUUID(prefix string, counter *int) string {
	*counter++
	return fmt.Sprintf("FA-%s-%d-%d", prefix, time.Now().UnixNano(), *counter)
}

// ─── Categories ──────────────────────────────────────────────────────

func (r *MemoryFARepo) CreateCategory(_ context.Context, c *domain.FixedAssetCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := c.CompanyID + ":" + c.Code
	if _, ok := r.catByCode[key]; ok {
		return domain.ErrFACategoryCodeExists
	}
	cp := *c
	cp.ID = faUUID("cat", &r.catCounter)
	r.categories[cp.ID] = &cp
	r.catByCode[key] = cp.ID
	c.ID = cp.ID
	return nil
}

func (r *MemoryFARepo) GetCategoryByID(_ context.Context, id string) (*domain.FixedAssetCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.categories[id]
	if !ok {
		return nil, domain.ErrFACategoryNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemoryFARepo) GetCategoryByCode(_ context.Context, companyID, code string) (*domain.FixedAssetCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := companyID + ":" + code
	id, ok := r.catByCode[key]
	if !ok {
		return nil, domain.ErrFACategoryNotFound
	}
	c, ok := r.categories[id]
	if !ok {
		return nil, domain.ErrFACategoryNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemoryFARepo) ListCategories(_ context.Context, filter domain.FACategoryFilter) ([]domain.FixedAssetCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.FixedAssetCategory
	for _, c := range r.categories {
		if c.CompanyID != filter.CompanyID {
			continue
		}
		if filter.ParentID != nil && (c.ParentID == nil || *c.ParentID != *filter.ParentID) {
			continue
		}
		if filter.Level != nil && c.Level != *filter.Level {
			continue
		}
		out = append(out, *c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, nil
}

func (r *MemoryFARepo) UpdateCategory(_ context.Context, c *domain.FixedAssetCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.categories[c.ID]; !ok {
		return domain.ErrFACategoryNotFound
	}
	cp := *c
	r.categories[c.ID] = &cp
	return nil
}

func (r *MemoryFARepo) DeleteCategory(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.categories[id]
	if !ok {
		return domain.ErrFACategoryNotFound
	}
	key := c.CompanyID + ":" + c.Code
	delete(r.catByCode, key)
	delete(r.categories, id)
	return nil
}

// ─── Fixed Assets ────────────────────────────────────────────────────

func (r *MemoryFARepo) CreateAsset(_ context.Context, a *domain.FixedAsset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := a.CompanyID + ":" + a.Code
	if _, ok := r.assetByCode[key]; ok {
		return domain.ErrFACodeExists
	}
	cp := *a
	cp.ID = faUUID("fa", &r.assetCounter)
	cp.CarryingAmount = cp.CalcCarryingAmount()
	r.assets[cp.ID] = &cp
	r.assetByCode[key] = cp.ID
	a.ID = cp.ID
	a.CarryingAmount = cp.CarryingAmount
	return nil
}

func (r *MemoryFARepo) GetAssetByID(_ context.Context, id string) (*domain.FixedAsset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.assets[id]
	if !ok {
		return nil, domain.ErrFANotFound
	}
	cp := *a
	return &cp, nil
}

func (r *MemoryFARepo) GetAssetByCode(_ context.Context, companyID, code string) (*domain.FixedAsset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := companyID + ":" + code
	id, ok := r.assetByCode[key]
	if !ok {
		return nil, domain.ErrFANotFound
	}
	a, ok := r.assets[id]
	if !ok {
		return nil, domain.ErrFANotFound
	}
	cp := *a
	return &cp, nil
}

func (r *MemoryFARepo) ListAssets(_ context.Context, filter domain.FAListFilter) ([]domain.FixedAsset, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.FixedAsset
	for _, a := range r.assets {
		if a.CompanyID != filter.CompanyID {
			continue
		}
		if filter.Status != nil && a.Status != *filter.Status {
			continue
		}
		if filter.CategoryID != nil && a.CategoryID != *filter.CategoryID {
			continue
		}
		if filter.DepartmentID != nil && a.DepartmentID != *filter.DepartmentID {
			continue
		}
		if filter.Keyword != "" {
			kw := strings.ToLower(filter.Keyword)
			if !strings.Contains(strings.ToLower(a.Code), kw) && !strings.Contains(strings.ToLower(a.Name), kw) {
				continue
			}
		}
		out = append(out, *a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })

	total := len(out)
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	start := filter.Offset
	if start > len(out) {
		start = len(out)
	}
	end := start + filter.Limit
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], total, nil
}

func (r *MemoryFARepo) UpdateAsset(_ context.Context, a *domain.FixedAsset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.assets[a.ID]; !ok {
		return domain.ErrFANotFound
	}
	cp := *a
	cp.CarryingAmount = cp.CalcCarryingAmount()
	r.assets[a.ID] = &cp
	return nil
}

func (r *MemoryFARepo) DeleteAsset(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.assets[id]
	if !ok {
		return domain.ErrFANotFound
	}
	key := a.CompanyID + ":" + a.Code
	delete(r.assetByCode, key)
	delete(r.assets, id)
	delete(r.deprByAsset, id)
	return nil
}

// ─── Depreciation ────────────────────────────────────────────────────

func (r *MemoryFARepo) CreateDepreciationEntry(_ context.Context, e *domain.DepreciationEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *e
	cp.ID = faUUID("depr", &r.deprCounter)
	r.depreciation[cp.ID] = &cp
	r.deprByAsset[cp.FixedAssetID] = append(r.deprByAsset[cp.FixedAssetID], cp.ID)
	e.ID = cp.ID
	return nil
}

func (r *MemoryFARepo) GetDepreciationEntry(_ context.Context, id string) (*domain.DepreciationEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.depreciation[id]
	if !ok {
		return nil, domain.ErrFANotFound
	}
	cp := *e
	return &cp, nil
}

func (r *MemoryFARepo) ListDepreciationByAsset(_ context.Context, assetID string) ([]domain.DepreciationEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.deprByAsset[assetID]
	out := make([]domain.DepreciationEntry, 0, len(ids))
	for _, id := range ids {
		if e, ok := r.depreciation[id]; ok {
			out = append(out, *e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].PeriodYear != out[j].PeriodYear {
			return out[i].PeriodYear < out[j].PeriodYear
		}
		return out[i].PeriodMonth < out[j].PeriodMonth
	})
	return out, nil
}

func (r *MemoryFARepo) ListDepreciationByPeriod(_ context.Context, periodID string) ([]domain.DepreciationEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.DepreciationEntry
	for _, e := range r.depreciation {
		if e.PeriodID == periodID {
			out = append(out, *e)
		}
	}
	return out, nil
}

func (r *MemoryFARepo) DepreciationExistsForPeriod(_ context.Context, assetID, periodID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, e := range r.depreciation {
		if e.FixedAssetID == assetID && e.PeriodID == periodID {
			return true, nil
		}
	}
	return false, nil
}

func (r *MemoryFARepo) DeleteDepreciationEntry(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.depreciation[id]
	if !ok {
		return nil
	}
	delete(r.depreciation, id)
	ids := r.deprByAsset[e.FixedAssetID]
	for i, v := range ids {
		if v == id {
			r.deprByAsset[e.FixedAssetID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

// ─── Transactions ────────────────────────────────────────────────────

func (r *MemoryFARepo) CreateTransaction(_ context.Context, t *domain.FixedAssetTransaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	cp.ID = faUUID("txn", &r.txnCounter)
	r.transactions[cp.ID] = &cp
	r.txnByAsset[cp.FixedAssetID] = append(r.txnByAsset[cp.FixedAssetID], cp.ID)
	t.ID = cp.ID
	return nil
}

func (r *MemoryFARepo) ListTransactionsByAsset(_ context.Context, assetID string) ([]domain.FixedAssetTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.txnByAsset[assetID]
	out := make([]domain.FixedAssetTransaction, 0, len(ids))
	for _, id := range ids {
		if t, ok := r.transactions[id]; ok {
			out = append(out, *t)
		}
	}
	return out, nil
}

// ─── Allocations ─────────────────────────────────────────────────────

func (r *MemoryFARepo) SetAllocations(_ context.Context, assetID string, allocs []domain.FixedAssetAllocation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := make([]domain.FixedAssetAllocation, len(allocs))
	for i := range allocs {
		cp[i] = allocs[i]
		cp[i].ID = faUUID("alloc", &r.allocCounter)
		cp[i].FixedAssetID = assetID
	}
	r.allocations[assetID] = cp
	return nil
}

func (r *MemoryFARepo) GetAllocations(_ context.Context, assetID string) ([]domain.FixedAssetAllocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	allocs, ok := r.allocations[assetID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.FixedAssetAllocation, len(allocs))
	copy(out, allocs)
	return out, nil
}

// ─── Inventory Plans ─────────────────────────────────────────────────

func (r *MemoryFARepo) CreateInventoryPlan(_ context.Context, p *domain.FixedAssetInventoryPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	cp.ID = faUUID("plan", &r.planCounter)
	r.inventoryPlans[cp.ID] = &cp
	p.ID = cp.ID
	return nil
}

func (r *MemoryFARepo) GetInventoryPlan(_ context.Context, id string) (*domain.FixedAssetInventoryPlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.inventoryPlans[id]
	if !ok {
		return nil, domain.ErrFANotFound
	}
	cp := *p
	return &cp, nil
}

func (r *MemoryFARepo) ListInventoryPlans(_ context.Context, companyID string) ([]domain.FixedAssetInventoryPlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.FixedAssetInventoryPlan
	for _, p := range r.inventoryPlans {
		if p.CompanyID == companyID {
			out = append(out, *p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PlanDate > out[j].PlanDate })
	return out, nil
}

func (r *MemoryFARepo) UpdateInventoryPlan(_ context.Context, p *domain.FixedAssetInventoryPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.inventoryPlans[p.ID]; !ok {
		return domain.ErrFANotFound
	}
	cp := *p
	r.inventoryPlans[p.ID] = &cp
	return nil
}

// ─── Inventory Results ───────────────────────────────────────────────

func (r *MemoryFARepo) CreateInventoryResult(_ context.Context, res *domain.FixedAssetInventoryResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *res
	cp.ID = faUUID("res", &r.resultCounter)
	r.inventoryResults[cp.ID] = &cp
	r.resultByPlan[cp.PlanID] = append(r.resultByPlan[cp.PlanID], cp.ID)
	res.ID = cp.ID
	return nil
}

func (r *MemoryFARepo) GetInventoryResultsByPlan(_ context.Context, planID string) ([]domain.FixedAssetInventoryResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.resultByPlan[planID]
	out := make([]domain.FixedAssetInventoryResult, 0, len(ids))
	for _, id := range ids {
		if res, ok := r.inventoryResults[id]; ok {
			out = append(out, *res)
		}
	}
	return out, nil
}
