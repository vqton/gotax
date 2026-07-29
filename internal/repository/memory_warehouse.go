package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"gotax/internal/domain"
)

var whSeq int64

func whID(prefix string) string {
	whSeq++
	return prefix + time.Now().Format("20060102150405") + fmt.Sprintf("%03d", whSeq%1000)
}

// ─── Warehouse ─────────────────────────────────────────────────────────

type MemoryWarehouseRepo struct {
	mu       sync.RWMutex
	data     map[string]*domain.Warehouse
	byCode   map[string]map[string]*domain.Warehouse
}

func NewMemoryWarehouseRepo() *MemoryWarehouseRepo {
	return &MemoryWarehouseRepo{
		data:   make(map[string]*domain.Warehouse),
		byCode: make(map[string]map[string]*domain.Warehouse),
	}
}

func (r *MemoryWarehouseRepo) CreateWarehouse(_ context.Context, w *domain.Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *w
	if cp.ID == "" { cp.ID = whID("WH-") }
	if r.byCode[cp.CompanyID] == nil { r.byCode[cp.CompanyID] = make(map[string]*domain.Warehouse) }
	if prev, ok := r.byCode[cp.CompanyID][cp.Code]; ok && prev != nil { return domain.ErrWarehouseCodeExists }
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	r.data[cp.ID] = &cp
	r.byCode[cp.CompanyID][cp.Code] = r.data[cp.ID]
	w.ID = cp.ID
	return nil
}

func (r *MemoryWarehouseRepo) GetWarehouseByID(_ context.Context, id string) (*domain.Warehouse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.data[id]
	if !ok { return nil, domain.ErrWarehouseNotFound }
	cp := *w
	return &cp, nil
}

func (r *MemoryWarehouseRepo) ListWarehouses(_ context.Context, companyID string) ([]domain.Warehouse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Warehouse
	for _, w := range r.data {
		if w.CompanyID != companyID { continue }
		out = append(out, *w)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, nil
}

func (r *MemoryWarehouseRepo) UpdateWarehouse(_ context.Context, w *domain.Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.data[w.ID]
	if !ok { return domain.ErrWarehouseNotFound }
	if r.byCode[w.CompanyID] == nil { r.byCode[w.CompanyID] = make(map[string]*domain.Warehouse) }
	if existing.Code != w.Code {
		if prev, ok := r.byCode[w.CompanyID][w.Code]; ok && prev != nil { return domain.ErrWarehouseCodeExists }
		delete(r.byCode[w.CompanyID], existing.Code)
	}
	cp := *w
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.data[w.ID] = &cp
	r.byCode[w.CompanyID][w.Code] = r.data[w.ID]
	return nil
}

func (r *MemoryWarehouseRepo) DeleteWarehouse(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.data[id]
	if !ok { return domain.ErrWarehouseNotFound }
	delete(r.data, id)
	if m, ok := r.byCode[w.CompanyID]; ok { delete(m, w.Code) }
	return nil
}

func (r *MemoryWarehouseRepo) GetWarehouseByCode(_ context.Context, companyID, code string) (*domain.Warehouse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.byCode[companyID]
	if !ok { return nil, domain.ErrWarehouseNotFound }
	w, ok := cm[code]
	if !ok { return nil, domain.ErrWarehouseNotFound }
	cp := *w
	return &cp, nil
}

// ─── Item Category ─────────────────────────────────────────────────────

type MemoryItemCategoryRepo struct {
	mu       sync.RWMutex
	data     map[string]*domain.ItemCategory
	byCode   map[string]map[string]*domain.ItemCategory
}

func NewMemoryItemCategoryRepo() *MemoryItemCategoryRepo {
	return &MemoryItemCategoryRepo{
		data:   make(map[string]*domain.ItemCategory),
		byCode: make(map[string]map[string]*domain.ItemCategory),
	}
}

func (r *MemoryItemCategoryRepo) CreateCategory(_ context.Context, c *domain.ItemCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *c
	if cp.ID == "" { cp.ID = whID("CAT-") }
	if r.byCode[cp.CompanyID] == nil { r.byCode[cp.CompanyID] = make(map[string]*domain.ItemCategory) }
	if prev, ok := r.byCode[cp.CompanyID][cp.Code]; ok && prev != nil { return domain.ErrCategoryCodeExists }
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	r.data[cp.ID] = &cp
	r.byCode[cp.CompanyID][cp.Code] = r.data[cp.ID]
	c.ID = cp.ID
	return nil
}

func (r *MemoryItemCategoryRepo) GetCategoryByID(_ context.Context, id string) (*domain.ItemCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.data[id]
	if !ok { return nil, domain.ErrCategoryNotFound }
	cp := *c
	return &cp, nil
}

func (r *MemoryItemCategoryRepo) ListCategories(_ context.Context, companyID string) ([]domain.ItemCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ItemCategory
	for _, c := range r.data {
		if c.CompanyID != companyID { continue }
		out = append(out, *c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, nil
}

func (r *MemoryItemCategoryRepo) UpdateCategory(_ context.Context, c *domain.ItemCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.data[c.ID]
	if !ok { return domain.ErrCategoryNotFound }
	if r.byCode[c.CompanyID] == nil { r.byCode[c.CompanyID] = make(map[string]*domain.ItemCategory) }
	if existing.Code != c.Code {
		if prev, ok := r.byCode[c.CompanyID][c.Code]; ok && prev != nil { return domain.ErrCategoryCodeExists }
		delete(r.byCode[c.CompanyID], existing.Code)
	}
	cp := *c
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.data[c.ID] = &cp
	r.byCode[c.CompanyID][c.Code] = r.data[c.ID]
	return nil
}

func (r *MemoryItemCategoryRepo) DeleteCategory(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.data[id]
	if !ok { return domain.ErrCategoryNotFound }
	delete(r.data, id)
	if m, ok := r.byCode[c.CompanyID]; ok { delete(m, c.Code) }
	return nil
}

func (r *MemoryItemCategoryRepo) GetCategoryByCode(_ context.Context, companyID, code string) (*domain.ItemCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.byCode[companyID]
	if !ok { return nil, domain.ErrCategoryNotFound }
	c, ok := cm[code]
	if !ok { return nil, domain.ErrCategoryNotFound }
	cp := *c
	return &cp, nil
}

// ─── Item ──────────────────────────────────────────────────────────────

type MemoryItemRepo struct {
	mu       sync.RWMutex
	data     map[string]*domain.Item
	byCode   map[string]map[string]*domain.Item
}

func NewMemoryItemRepo() *MemoryItemRepo {
	return &MemoryItemRepo{
		data:   make(map[string]*domain.Item),
		byCode: make(map[string]map[string]*domain.Item),
	}
}

func (r *MemoryItemRepo) CreateItem(_ context.Context, i *domain.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *i
	if cp.ID == "" { cp.ID = whID("ITEM-") }
	if r.byCode[cp.CompanyID] == nil { r.byCode[cp.CompanyID] = make(map[string]*domain.Item) }
	if prev, ok := r.byCode[cp.CompanyID][cp.Code]; ok && prev != nil { return domain.ErrItemCodeExists }
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	r.data[cp.ID] = &cp
	r.byCode[cp.CompanyID][cp.Code] = r.data[cp.ID]
	i.ID = cp.ID
	return nil
}

func (r *MemoryItemRepo) GetItemByID(_ context.Context, id string) (*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	i, ok := r.data[id]
	if !ok { return nil, domain.ErrItemNotFound }
	cp := *i
	return &cp, nil
}

func (r *MemoryItemRepo) ListItems(_ context.Context, companyID string) ([]domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Item
	for _, i := range r.data {
		if i.CompanyID != companyID { continue }
		out = append(out, *i)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, nil
}

func (r *MemoryItemRepo) UpdateItem(_ context.Context, i *domain.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.data[i.ID]
	if !ok { return domain.ErrItemNotFound }
	if r.byCode[i.CompanyID] == nil { r.byCode[i.CompanyID] = make(map[string]*domain.Item) }
	if existing.Code != i.Code {
		if prev, ok := r.byCode[i.CompanyID][i.Code]; ok && prev != nil { return domain.ErrItemCodeExists }
		delete(r.byCode[i.CompanyID], existing.Code)
	}
	cp := *i
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.data[i.ID] = &cp
	r.byCode[i.CompanyID][i.Code] = r.data[i.ID]
	return nil
}

func (r *MemoryItemRepo) DeleteItem(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	i, ok := r.data[id]
	if !ok { return domain.ErrItemNotFound }
	delete(r.data, id)
	if m, ok := r.byCode[i.CompanyID]; ok { delete(m, i.Code) }
	return nil
}

func (r *MemoryItemRepo) GetItemByCode(_ context.Context, companyID, code string) (*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cm, ok := r.byCode[companyID]
	if !ok { return nil, domain.ErrItemNotFound }
	i, ok := cm[code]
	if !ok { return nil, domain.ErrItemNotFound }
	cp := *i
	return &cp, nil
}

// ─── Stock Balance ─────────────────────────────────────────────────────

type MemoryStockBalanceRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.StockBalance
}

func NewMemoryStockBalanceRepo() *MemoryStockBalanceRepo {
	return &MemoryStockBalanceRepo{data: make(map[string]*domain.StockBalance)}
}

func (r *MemoryStockBalanceRepo) CreateStockBalance(_ context.Context, b *domain.StockBalance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *b
	if cp.ID == "" { cp.ID = whID("BAL-") }
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	r.data[cp.ID] = &cp
	b.ID = cp.ID
	return nil
}

func (r *MemoryStockBalanceRepo) GetStockBalanceByID(_ context.Context, id string) (*domain.StockBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.data[id]
	if !ok { return nil, domain.ErrBalanceNotFound }
	cp := *b
	return &cp, nil
}

func (r *MemoryStockBalanceRepo) FindStockBalance(_ context.Context, companyID, warehouseID, itemID, period string) (*domain.StockBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, b := range r.data {
		if b.CompanyID == companyID && b.WarehouseID == warehouseID && b.ItemID == itemID && b.Period == period {
			cp := *b
			return &cp, nil
		}
	}
	return nil, domain.ErrBalanceNotFound
}

func (r *MemoryStockBalanceRepo) ListStockBalances(_ context.Context, companyID, warehouseID string) ([]domain.StockBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.StockBalance
	for _, b := range r.data {
		if b.CompanyID != companyID { continue }
		if warehouseID != "" && b.WarehouseID != warehouseID { continue }
		out = append(out, *b)
	}
	return out, nil
}

func (r *MemoryStockBalanceRepo) UpdateStockBalance(_ context.Context, b *domain.StockBalance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.data[b.ID]
	if !ok { return domain.ErrBalanceNotFound }
	cp := *b
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.data[b.ID] = &cp
	return nil
}

func (r *MemoryStockBalanceRepo) UpsertStockBalance(_ context.Context, b *domain.StockBalance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for _, existing := range r.data {
		if existing.CompanyID == b.CompanyID && existing.WarehouseID == b.WarehouseID && existing.ItemID == b.ItemID && existing.Period == b.Period {
			existing.Quantity = b.Quantity
			existing.UnitCost = b.UnitCost
			existing.TotalCost = b.TotalCost
			existing.LastTransactionAt = b.LastTransactionAt
			existing.UpdatedAt = now
			b.ID = existing.ID
			return nil
		}
	}
	cp := *b
	if cp.ID == "" { cp.ID = whID("BAL-") }
	cp.CreatedAt, cp.UpdatedAt = now, now
	r.data[cp.ID] = &cp
	b.ID = cp.ID
	return nil
}

// ─── Inventory Transaction ─────────────────────────────────────────────

type MemoryInventoryTransactionRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.InventoryTransaction
}

func NewMemoryInventoryTransactionRepo() *MemoryInventoryTransactionRepo {
	return &MemoryInventoryTransactionRepo{data: make(map[string]*domain.InventoryTransaction)}
}

func (r *MemoryInventoryTransactionRepo) CreateInventoryTransaction(_ context.Context, t *domain.InventoryTransaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" { cp.ID = whID("TXN-") }
	cp.CreatedAt = time.Now()
	r.data[cp.ID] = &cp
	t.ID = cp.ID
	return nil
}

func (r *MemoryInventoryTransactionRepo) GetInventoryTransactionByID(_ context.Context, id string) (*domain.InventoryTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok { return nil, domain.ErrTransNotFound }
	cp := *t
	return &cp, nil
}

func (r *MemoryInventoryTransactionRepo) ListInventoryTransactions(_ context.Context, companyID, warehouseID, itemID string, offset, limit int) ([]domain.InventoryTransaction, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []domain.InventoryTransaction
	for _, t := range r.data {
		if t.CompanyID != companyID { continue }
		if warehouseID != "" && t.WarehouseID != warehouseID { continue }
		if itemID != "" && t.ItemID != itemID { continue }
		filtered = append(filtered, *t)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].CreatedAt.After(filtered[j].CreatedAt) })
	total := len(filtered)
	start, end := offset, offset+limit
	if limit <= 0 { start, end = 0, total }
	if start > total { return []domain.InventoryTransaction{}, total, nil }
	if end > total { end = total }
	return filtered[start:end], total, nil
}

// ─── Stock Transfer ────────────────────────────────────────────────────

type MemoryStockTransferRepo struct {
	mu       sync.RWMutex
	data     map[string]*domain.StockTransfer
	items    map[string][]domain.TransferItem
}

func NewMemoryStockTransferRepo() *MemoryStockTransferRepo {
	return &MemoryStockTransferRepo{
		data:  make(map[string]*domain.StockTransfer),
		items: make(map[string][]domain.TransferItem),
	}
}

func (r *MemoryStockTransferRepo) CreateStockTransfer(_ context.Context, t *domain.StockTransfer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" { cp.ID = whID("TRF-") }
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	cp.Items = nil
	r.data[cp.ID] = &cp
	t.ID = cp.ID
	if len(t.Items) > 0 {
		its := make([]domain.TransferItem, len(t.Items))
		for i := range t.Items {
			its[i] = t.Items[i]
			its[i].ID = whID("TFI-")
			its[i].TransferID = cp.ID
		}
		r.items[cp.ID] = its
	}
	return nil
}

func (r *MemoryStockTransferRepo) GetStockTransferByID(_ context.Context, id string) (*domain.StockTransfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok { return nil, domain.ErrTransferNotFound }
	cp := *t
	cp.Items = r.items[id]
	return &cp, nil
}

func (r *MemoryStockTransferRepo) ListStockTransfers(_ context.Context, companyID string) ([]domain.StockTransfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.StockTransfer
	for _, t := range r.data {
		if t.CompanyID != companyID { continue }
		out = append(out, *t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (r *MemoryStockTransferRepo) UpdateStockTransfer(_ context.Context, t *domain.StockTransfer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.data[t.ID]
	if !ok { return domain.ErrTransferNotFound }
	cp := *t
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	cp.Items = nil
	r.data[t.ID] = &cp
	if len(t.Items) > 0 {
		its := make([]domain.TransferItem, len(t.Items))
		for i := range t.Items {
			its[i] = t.Items[i]
			its[i].TransferID = t.ID
		}
		r.items[t.ID] = its
	}
	return nil
}

func (r *MemoryStockTransferRepo) UpdateStockTransferStatus(_ context.Context, id string, status domain.TransferStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.data[id]
	if !ok { return domain.ErrTransferNotFound }
	t.Status = status
	t.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryStockTransferRepo) GetTransferItems(_ context.Context, transferID string) ([]domain.TransferItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	its, ok := r.items[transferID]
	if !ok { return []domain.TransferItem{}, nil }
	out := make([]domain.TransferItem, len(its))
	copy(out, its)
	return out, nil
}

func (r *MemoryStockTransferRepo) CreateTransferItem(_ context.Context, item *domain.TransferItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = whID("TFI-")
	r.items[item.TransferID] = append(r.items[item.TransferID], *item)
	return nil
}

// ─── Stock Adjustment ──────────────────────────────────────────────────

type MemoryStockAdjustmentRepo struct {
	mu       sync.RWMutex
	data     map[string]*domain.StockAdjustment
	items    map[string][]domain.AdjItem
}

func NewMemoryStockAdjustmentRepo() *MemoryStockAdjustmentRepo {
	return &MemoryStockAdjustmentRepo{
		data:  make(map[string]*domain.StockAdjustment),
		items: make(map[string][]domain.AdjItem),
	}
}

func (r *MemoryStockAdjustmentRepo) CreateStockAdjustment(_ context.Context, a *domain.StockAdjustment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *a
	if cp.ID == "" { cp.ID = whID("ADJ-") }
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	cp.Items = nil
	r.data[cp.ID] = &cp
	a.ID = cp.ID
	if len(a.Items) > 0 {
		its := make([]domain.AdjItem, len(a.Items))
		for i := range a.Items {
			its[i] = a.Items[i]
			its[i].ID = whID("AJI-")
			its[i].AdjustmentID = cp.ID
		}
		r.items[cp.ID] = its
	}
	return nil
}

func (r *MemoryStockAdjustmentRepo) GetStockAdjustmentByID(_ context.Context, id string) (*domain.StockAdjustment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.data[id]
	if !ok { return nil, domain.ErrAdjNotFound }
	cp := *a
	cp.Items = r.items[id]
	return &cp, nil
}

func (r *MemoryStockAdjustmentRepo) ListStockAdjustments(_ context.Context, companyID string) ([]domain.StockAdjustment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.StockAdjustment
	for _, a := range r.data {
		if a.CompanyID != companyID { continue }
		out = append(out, *a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (r *MemoryStockAdjustmentRepo) UpdateStockAdjustment(_ context.Context, a *domain.StockAdjustment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.data[a.ID]
	if !ok { return domain.ErrAdjNotFound }
	cp := *a
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	cp.Items = nil
	r.data[a.ID] = &cp
	if len(a.Items) > 0 {
		its := make([]domain.AdjItem, len(a.Items))
		for i := range a.Items {
			its[i] = a.Items[i]
			its[i].AdjustmentID = a.ID
		}
		r.items[a.ID] = its
	}
	return nil
}

func (r *MemoryStockAdjustmentRepo) UpdateStockAdjustmentStatus(_ context.Context, id string, status domain.AdjStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.data[id]
	if !ok { return domain.ErrAdjNotFound }
	a.Status = status
	a.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryStockAdjustmentRepo) GetAdjustmentItems(_ context.Context, adjustmentID string) ([]domain.AdjItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	its, ok := r.items[adjustmentID]
	if !ok { return []domain.AdjItem{}, nil }
	out := make([]domain.AdjItem, len(its))
	copy(out, its)
	return out, nil
}

func (r *MemoryStockAdjustmentRepo) CreateAdjustmentItem(_ context.Context, item *domain.AdjItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = whID("AJI-")
	r.items[item.AdjustmentID] = append(r.items[item.AdjustmentID], *item)
	return nil
}

// ─── Stock Take ────────────────────────────────────────────────────────

type MemoryStockTakeRepo struct {
	mu       sync.RWMutex
	data     map[string]*domain.StockTake
	items    map[string][]domain.TakeItem
}

func NewMemoryStockTakeRepo() *MemoryStockTakeRepo {
	return &MemoryStockTakeRepo{
		data:  make(map[string]*domain.StockTake),
		items: make(map[string][]domain.TakeItem),
	}
}

func (r *MemoryStockTakeRepo) CreateStockTake(_ context.Context, t *domain.StockTake) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" { cp.ID = whID("TAKE-") }
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	cp.Items = nil
	r.data[cp.ID] = &cp
	t.ID = cp.ID
	if len(t.Items) > 0 {
		its := make([]domain.TakeItem, len(t.Items))
		for i := range t.Items {
			its[i] = t.Items[i]
			its[i].ID = whID("TKI-")
			its[i].TakeID = cp.ID
		}
		r.items[cp.ID] = its
	}
	return nil
}

func (r *MemoryStockTakeRepo) GetStockTakeByID(_ context.Context, id string) (*domain.StockTake, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok { return nil, domain.ErrTakeNotFound }
	cp := *t
	cp.Items = r.items[id]
	return &cp, nil
}

func (r *MemoryStockTakeRepo) ListStockTakes(_ context.Context, companyID string) ([]domain.StockTake, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.StockTake
	for _, t := range r.data {
		if t.CompanyID != companyID { continue }
		out = append(out, *t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (r *MemoryStockTakeRepo) UpdateStockTake(_ context.Context, t *domain.StockTake) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.data[t.ID]
	if !ok { return domain.ErrTakeNotFound }
	cp := *t
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	cp.Items = nil
	r.data[t.ID] = &cp
	if len(t.Items) > 0 {
		its := make([]domain.TakeItem, len(t.Items))
		for i := range t.Items {
			its[i] = t.Items[i]
			its[i].TakeID = t.ID
		}
		r.items[t.ID] = its
	}
	return nil
}

func (r *MemoryStockTakeRepo) UpdateStockTakeStatus(_ context.Context, id string, status domain.TakeStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.data[id]
	if !ok { return domain.ErrTakeNotFound }
	t.Status = status
	t.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryStockTakeRepo) GetTakeItems(_ context.Context, takeID string) ([]domain.TakeItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	its, ok := r.items[takeID]
	if !ok { return []domain.TakeItem{}, nil }
	out := make([]domain.TakeItem, len(its))
	copy(out, its)
	return out, nil
}

func (r *MemoryStockTakeRepo) CreateTakeItem(_ context.Context, item *domain.TakeItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = whID("TKI-")
	r.items[item.TakeID] = append(r.items[item.TakeID], *item)
	return nil
}

// ─── Inventory Valuation Run ───────────────────────────────────────────

type MemoryInventoryValuationRunRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.InventoryValuationRun
}

func NewMemoryInventoryValuationRunRepo() *MemoryInventoryValuationRunRepo {
	return &MemoryInventoryValuationRunRepo{data: make(map[string]*domain.InventoryValuationRun)}
}

func (r *MemoryInventoryValuationRunRepo) CreateValuationRun(_ context.Context, v *domain.InventoryValuationRun) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *v
	if cp.ID == "" { cp.ID = whID("VAL-") }
	cp.CreatedAt = time.Now()
	r.data[cp.ID] = &cp
	v.ID = cp.ID
	return nil
}

func (r *MemoryInventoryValuationRunRepo) GetValuationRunByID(_ context.Context, id string) (*domain.InventoryValuationRun, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.data[id]
	if !ok { return nil, domain.ErrValRunNotFound }
	cp := *v
	return &cp, nil
}

func (r *MemoryInventoryValuationRunRepo) ListValuationRuns(_ context.Context, companyID string) ([]domain.InventoryValuationRun, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.InventoryValuationRun
	for _, v := range r.data {
		if v.CompanyID != companyID { continue }
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (r *MemoryInventoryValuationRunRepo) UpdateValuationRun(_ context.Context, v *domain.InventoryValuationRun) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[v.ID]; !ok { return domain.ErrValRunNotFound }
	cp := *v
	r.data[v.ID] = &cp
	return nil
}
