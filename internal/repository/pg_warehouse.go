package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

// ─── Warehouse ───────────────────────────────────────────────────────────

type PGWarehouseRepo struct{ db *gorm.DB }

func NewPGWarehouseRepo(db *gorm.DB) *PGWarehouseRepo { return &PGWarehouseRepo{db} }

func (r *PGWarehouseRepo) CreateWarehouse(ctx context.Context, w *domain.Warehouse) error {
	m := warehouseToGORM(w)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	w.ID = m.ID
	w.CreatedAt = m.CreatedAt
	w.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *PGWarehouseRepo) GetWarehouseByID(ctx context.Context, id string) (*domain.Warehouse, error) {
	var m domain.WarehouseGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWarehouseNotFound
		}
		return nil, err
	}
	return gormWHToDomain(&m), nil
}

func (r *PGWarehouseRepo) ListWarehouses(ctx context.Context, companyID string) ([]domain.Warehouse, error) {
	var models []domain.WarehouseGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Warehouse, len(models))
	for i := range models {
		out[i] = *gormWHToDomain(&models[i])
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateWarehouse(ctx context.Context, w *domain.Warehouse) error {
	return r.db.WithContext(ctx).Model(&domain.WarehouseGORM{}).Where("id = ?", w.ID).Updates(map[string]interface{}{
		"code":       w.Code,
		"name":       w.Name,
		"location":   nullStrG(w.Address),
		"manager_id": nullStrG(w.Manager),
		"is_active":  w.IsActive,
	}).Error
}

func (r *PGWarehouseRepo) DeleteWarehouse(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.WarehouseGORM{}).Error
}

func (r *PGWarehouseRepo) GetWarehouseByCode(ctx context.Context, companyID, code string) (*domain.Warehouse, error) {
	var m domain.WarehouseGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWarehouseNotFound
		}
		return nil, err
	}
	return gormWHToDomain(&m), nil
}

func warehouseToGORM(w *domain.Warehouse) domain.WarehouseGORM {
	return domain.WarehouseGORM{
		ID:        w.ID,
		CompanyID: w.CompanyID,
		Code:      w.Code,
		Name:      w.Name,
		Location:  nullStrG(w.Address),
		ManagerID: nullStrG(w.Manager),
		IsActive:  w.IsActive,
	}
}

func gormWHToDomain(m *domain.WarehouseGORM) *domain.Warehouse {
	w := &domain.Warehouse{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Code:      m.Code,
		Name:      m.Name,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Location != nil {
		w.Address = *m.Location
	}
	if m.ManagerID != nil {
		w.Manager = *m.ManagerID
	}
	return w
}

// ─── Item Category ───────────────────────────────────────────────────────

type PGItemCategoryRepo struct{ db *gorm.DB }

func NewPGItemCategoryRepo(db *gorm.DB) *PGItemCategoryRepo { return &PGItemCategoryRepo{db} }

func (r *PGItemCategoryRepo) CreateCategory(ctx context.Context, c *domain.ItemCategory) error {
	m := itemCategoryToGORM(c)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	c.ID = m.ID
	c.CreatedAt = m.CreatedAt
	c.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *PGItemCategoryRepo) GetCategoryByID(ctx context.Context, id string) (*domain.ItemCategory, error) {
	var m domain.ItemCategoryGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}
	return gormCategoryToDomain(&m), nil
}

func (r *PGItemCategoryRepo) ListCategories(ctx context.Context, companyID string) ([]domain.ItemCategory, error) {
	var models []domain.ItemCategoryGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ItemCategory, len(models))
	for i := range models {
		out[i] = *gormCategoryToDomain(&models[i])
	}
	return out, nil
}

func (r *PGItemCategoryRepo) UpdateCategory(ctx context.Context, c *domain.ItemCategory) error {
	return r.db.WithContext(ctx).Model(&domain.ItemCategoryGORM{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
		"code":      c.Code,
		"name":      c.Name,
		"parent_id": nullStrG(c.ParentID),
		"is_active": c.IsActive,
	}).Error
}

func (r *PGItemCategoryRepo) DeleteCategory(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.ItemCategoryGORM{}).Error
}

func (r *PGItemCategoryRepo) GetCategoryByCode(ctx context.Context, companyID, code string) (*domain.ItemCategory, error) {
	var m domain.ItemCategoryGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}
	return gormCategoryToDomain(&m), nil
}

func itemCategoryToGORM(c *domain.ItemCategory) domain.ItemCategoryGORM {
	return domain.ItemCategoryGORM{
		ID:        c.ID,
		CompanyID: c.CompanyID,
		Code:      c.Code,
		Name:      c.Name,
		ParentID:  nullStrG(c.ParentID),
		IsActive:  c.IsActive,
	}
}

func gormCategoryToDomain(m *domain.ItemCategoryGORM) *domain.ItemCategory {
	c := &domain.ItemCategory{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Code:      m.Code,
		Name:      m.Name,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.ParentID != nil {
		c.ParentID = *m.ParentID
	}
	return c
}

// ─── Item ────────────────────────────────────────────────────────────────

type PGItemRepo struct{ db *gorm.DB }

func NewPGItemRepo(db *gorm.DB) *PGItemRepo { return &PGItemRepo{db} }

func (r *PGItemRepo) CreateItem(ctx context.Context, i *domain.Item) error {
	m := itemToGORM(i)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	i.ID = m.ID
	i.CreatedAt = m.CreatedAt
	i.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *PGItemRepo) GetItemByID(ctx context.Context, id string) (*domain.Item, error) {
	var m domain.ItemGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrItemNotFound
		}
		return nil, err
	}
	return gormItemToDomain(&m), nil
}

func (r *PGItemRepo) ListItems(ctx context.Context, companyID string) ([]domain.Item, error) {
	var models []domain.ItemGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Item, len(models))
	for i := range models {
		out[i] = *gormItemToDomain(&models[i])
	}
	return out, nil
}

func (r *PGItemRepo) UpdateItem(ctx context.Context, i *domain.Item) error {
	return r.db.WithContext(ctx).Model(&domain.ItemGORM{}).Where("id = ?", i.ID).Updates(map[string]interface{}{
		"code":          i.Code,
		"name":          i.Name,
		"category_id":   nullStrG(i.CategoryID),
		"unit":          nullStrG(i.Unit),
		"cost_method":   string(i.ValuationMethod),
		"is_active":     i.IsActive,
	}).Error
}

func (r *PGItemRepo) DeleteItem(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.ItemGORM{}).Error
}

func (r *PGItemRepo) GetItemByCode(ctx context.Context, companyID, code string) (*domain.Item, error) {
	var m domain.ItemGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrItemNotFound
		}
		return nil, err
	}
	return gormItemToDomain(&m), nil
}

func itemToGORM(i *domain.Item) domain.ItemGORM {
	return domain.ItemGORM{
		ID:          i.ID,
		CompanyID:   i.CompanyID,
		Code:        i.Code,
		Name:        i.Name,
		CategoryID:  nullStrG(i.CategoryID),
		Unit:        nullStrG(i.Unit),
		CostMethod:  string(i.ValuationMethod),
		IsActive:    i.IsActive,
	}
}

func gormItemToDomain(m *domain.ItemGORM) *domain.Item {
	i := &domain.Item{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		Code:            m.Code,
		Name:            m.Name,
		ValuationMethod: domain.ValuationMethod(m.CostMethod),
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
	if m.CategoryID != nil {
		i.CategoryID = *m.CategoryID
	}
	if m.Unit != nil {
		i.Unit = *m.Unit
	}
	return i
}

// ─── Stock Balance ───────────────────────────────────────────────────────

type PGStockBalanceRepo struct{ db *gorm.DB }

func NewPGStockBalanceRepo(db *gorm.DB) *PGStockBalanceRepo { return &PGStockBalanceRepo{db} }

func (r *PGStockBalanceRepo) CreateStockBalance(ctx context.Context, b *domain.StockBalance) error {
	m := stockBalanceToGORM(b)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	b.ID = m.ID
	b.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *PGStockBalanceRepo) GetStockBalanceByID(ctx context.Context, id string) (*domain.StockBalance, error) {
	var m domain.StockBalanceGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrBalanceNotFound
		}
		return nil, err
	}
	return gormStockBalanceToDomain(&m), nil
}

func (r *PGStockBalanceRepo) FindStockBalance(ctx context.Context, companyID, warehouseID, itemID, period string) (*domain.StockBalance, error) {
	var m domain.StockBalanceGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND warehouse_id = ? AND item_id = ? AND period = ?",
		companyID, warehouseID, itemID, period).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrBalanceNotFound
		}
		return nil, err
	}
	return gormStockBalanceToDomain(&m), nil
}

func (r *PGStockBalanceRepo) ListStockBalances(ctx context.Context, companyID, warehouseID string) ([]domain.StockBalance, error) {
	q := r.db.WithContext(ctx).Model(&domain.StockBalanceGORM{}).Where("company_id = ?", companyID)
	if warehouseID != "" {
		q = q.Where("warehouse_id = ?", warehouseID)
	}
	var models []domain.StockBalanceGORM
	if err := q.Order("item_id, period").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.StockBalance, len(models))
	for i := range models {
		out[i] = *gormStockBalanceToDomain(&models[i])
	}
	return out, nil
}

func (r *PGStockBalanceRepo) UpdateStockBalance(ctx context.Context, b *domain.StockBalance) error {
	return r.db.WithContext(ctx).Model(&domain.StockBalanceGORM{}).Where("id = ?", b.ID).Updates(map[string]interface{}{
		"qty_on_hand": b.Quantity,
		"unit_cost":   b.UnitCost,
		"total_value": b.TotalCost,
	}).Error
}

func (r *PGStockBalanceRepo) UpsertStockBalance(ctx context.Context, b *domain.StockBalance) error {
	m := domain.StockBalanceGORM{}
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND warehouse_id = ? AND item_id = ? AND period = ?",
			b.CompanyID, b.WarehouseID, b.ItemID, b.Period).
		Assign(map[string]interface{}{
			"qty_on_hand": b.Quantity,
			"unit_cost":   b.UnitCost,
			"total_value": b.TotalCost,
		}).
		FirstOrCreate(&m).Error
	if err != nil {
		return err
	}
	b.ID = m.ID
	return nil
}

func stockBalanceToGORM(b *domain.StockBalance) domain.StockBalanceGORM {
	return domain.StockBalanceGORM{
		ID:          b.ID,
		CompanyID:   b.CompanyID,
		WarehouseID: b.WarehouseID,
		ItemID:      b.ItemID,
		Period:      b.Period,
		QtyOnHand:   b.Quantity,
		UnitCost:    b.UnitCost,
		TotalValue:  b.TotalCost,
	}
}

func gormStockBalanceToDomain(m *domain.StockBalanceGORM) *domain.StockBalance {
	return &domain.StockBalance{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		WarehouseID: m.WarehouseID,
		ItemID:      m.ItemID,
		Period:      m.Period,
		Quantity:    m.QtyOnHand,
		UnitCost:    m.UnitCost,
		TotalCost:   m.TotalValue,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ─── Inventory Transaction ──────────────────────────────────────────────

type PGInventoryTransactionRepo struct{ db *gorm.DB }

func NewPGInventoryTransactionRepo(db *gorm.DB) *PGInventoryTransactionRepo {
	return &PGInventoryTransactionRepo{db}
}

func (r *PGInventoryTransactionRepo) CreateInventoryTransaction(ctx context.Context, t *domain.InventoryTransaction) error {
	m := invTxnToGORM(t)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	t.ID = m.ID
	t.CreatedAt = m.CreatedAt
	return nil
}

func (r *PGInventoryTransactionRepo) GetInventoryTransactionByID(ctx context.Context, id string) (*domain.InventoryTransaction, error) {
	var m domain.InventoryTransactionGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTransNotFound
		}
		return nil, err
	}
	return gormInvTxnToDomain(&m), nil
}

func (r *PGInventoryTransactionRepo) ListInventoryTransactions(ctx context.Context, companyID, warehouseID, itemID string, offset, limit int) ([]domain.InventoryTransaction, int, error) {
	q := r.db.WithContext(ctx).Model(&domain.InventoryTransactionGORM{}).Where("company_id = ?", companyID)
	if warehouseID != "" {
		q = q.Where("warehouse_id = ?", warehouseID)
	}
	if itemID != "" {
		q = q.Where("item_id = ?", itemID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 50
	}
	var models []domain.InventoryTransactionGORM
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.InventoryTransaction, len(models))
	for i := range models {
		out[i] = *gormInvTxnToDomain(&models[i])
	}
	return out, int(total), nil
}

func invTxnToGORM(t *domain.InventoryTransaction) domain.InventoryTransactionGORM {
	return domain.InventoryTransactionGORM{
		ID:           t.ID,
		CompanyID:    t.CompanyID,
		WarehouseID:  t.WarehouseID,
		ItemID:       t.ItemID,
		MovementType: string(t.TransType),
		Qty:          t.Quantity,
		UnitCost:     t.UnitCost,
		TotalValue:   t.TotalCost,
		RefDocID:     nullStrG(t.RefID),
		RefDocType:   nullStrG(t.RefType),
		Note:         nullStrG(t.Notes),
		CreatedBy:    t.CreatedBy,
	}
}

func gormInvTxnToDomain(m *domain.InventoryTransactionGORM) *domain.InventoryTransaction {
	t := &domain.InventoryTransaction{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		WarehouseID: m.WarehouseID,
		ItemID:    m.ItemID,
		TransType: domain.TransactionType(m.MovementType),
		Quantity:  m.Qty,
		UnitCost:  m.UnitCost,
		TotalCost: m.TotalValue,
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt,
	}
	if m.RefDocID != nil {
		t.RefID = *m.RefDocID
	}
	if m.RefDocType != nil {
		t.RefType = *m.RefDocType
	}
	if m.Note != nil {
		t.Notes = *m.Note
	}
	return t
}

// ─── Stock Transfer ──────────────────────────────────────────────────────

type PGStockTransferRepo struct{ db *gorm.DB }

func NewPGStockTransferRepo(db *gorm.DB) *PGStockTransferRepo { return &PGStockTransferRepo{db} }

func (r *PGStockTransferRepo) CreateStockTransfer(ctx context.Context, t *domain.StockTransfer) error {
	m := stockTransferToGORM(t)
	items := make([]domain.TransferItemGORM, len(t.Items))
	for i, it := range t.Items {
		items[i] = transferItemToGORM(&it)
	}
	m.Items = items
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	t.ID = m.ID
	t.CreatedAt = m.CreatedAt
	t.UpdatedAt = m.UpdatedAt
	for i := range t.Items {
		if i < len(m.Items) {
			t.Items[i].ID = fmt.Sprint(m.Items[i].ID)
		}
	}
	return nil
}

func (r *PGStockTransferRepo) GetStockTransferByID(ctx context.Context, id string) (*domain.StockTransfer, error) {
	var m domain.StockTransferGORM
	if err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTransferNotFound
		}
		return nil, err
	}
	return gormStockTransferToDomain(&m), nil
}

func (r *PGStockTransferRepo) ListStockTransfers(ctx context.Context, companyID string) ([]domain.StockTransfer, error) {
	var models []domain.StockTransferGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.StockTransfer, len(models))
	for i := range models {
		out[i] = *gormStockTransferToDomain(&models[i])
	}
	return out, nil
}

func (r *PGStockTransferRepo) UpdateStockTransfer(ctx context.Context, t *domain.StockTransfer) error {
	return r.db.WithContext(ctx).Model(&domain.StockTransferGORM{}).Where("id = ?", t.ID).Updates(map[string]interface{}{
		"from_warehouse_id": t.FromWarehouseID,
		"to_warehouse_id":   t.ToWarehouseID,
		"status":            string(t.Status),
	}).Error
}

func (r *PGStockTransferRepo) UpdateStockTransferStatus(ctx context.Context, id string, status domain.TransferStatus) error {
	return r.db.WithContext(ctx).Model(&domain.StockTransferGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGStockTransferRepo) GetTransferItems(ctx context.Context, transferID string) ([]domain.TransferItem, error) {
	var models []domain.TransferItemGORM
	if err := r.db.WithContext(ctx).Where("transfer_id = ?", transferID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.TransferItem, len(models))
	for i := range models {
		out[i] = *gormTransferItemToDomain(&models[i])
	}
	return out, nil
}

func (r *PGStockTransferRepo) CreateTransferItem(ctx context.Context, it *domain.TransferItem) error {
	m := transferItemToGORM(it)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	it.ID = fmt.Sprint(m.ID)
	return nil
}

func stockTransferToGORM(t *domain.StockTransfer) domain.StockTransferGORM {
	return domain.StockTransferGORM{
		ID:              t.ID,
		CompanyID:       t.CompanyID,
		FromWarehouseID: t.FromWarehouseID,
		ToWarehouseID:   t.ToWarehouseID,
		Status:          string(t.Status),
		CreatedBy:       t.CreatedBy,
	}
}

func gormStockTransferToDomain(m *domain.StockTransferGORM) *domain.StockTransfer {
	t := &domain.StockTransfer{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		FromWarehouseID: m.FromWarehouseID,
		ToWarehouseID:   m.ToWarehouseID,
		Status:          domain.TransferStatus(m.Status),
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
	if m.Items != nil {
		t.Items = make([]domain.TransferItem, len(m.Items))
		for i, it := range m.Items {
			t.Items[i] = *gormTransferItemToDomain(&it)
		}
	}
	return t
}

func transferItemToGORM(it *domain.TransferItem) domain.TransferItemGORM {
	return domain.TransferItemGORM{
		TransferID: it.TransferID,
		ItemID:     it.ItemID,
		Qty:        it.Quantity,
	}
}

func gormTransferItemToDomain(m *domain.TransferItemGORM) *domain.TransferItem {
	it := &domain.TransferItem{
		ID:         fmt.Sprint(m.ID),
		TransferID: m.TransferID,
		ItemID:     m.ItemID,
		Quantity:   m.Qty,
	}
	return it
}

// ─── Stock Adjustment ───────────────────────────────────────────────────

type PGStockAdjustmentRepo struct{ db *gorm.DB }

func NewPGStockAdjustmentRepo(db *gorm.DB) *PGStockAdjustmentRepo { return &PGStockAdjustmentRepo{db} }

func (r *PGStockAdjustmentRepo) CreateStockAdjustment(ctx context.Context, a *domain.StockAdjustment) error {
	m := stockAdjToGORM(a)
	items := make([]domain.AdjItemGORM, len(a.Items))
	for i, it := range a.Items {
		items[i] = adjItemToGORM(&it)
	}
	m.Items = items
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	a.ID = m.ID
	a.CreatedAt = m.CreatedAt
	a.UpdatedAt = m.UpdatedAt
	for i := range a.Items {
		if i < len(m.Items) {
			a.Items[i].ID = fmt.Sprint(m.Items[i].ID)
		}
	}
	return nil
}

func (r *PGStockAdjustmentRepo) GetStockAdjustmentByID(ctx context.Context, id string) (*domain.StockAdjustment, error) {
	var m domain.StockAdjustmentGORM
	if err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAdjNotFound
		}
		return nil, err
	}
	return gormStockAdjToDomain(&m), nil
}

func (r *PGStockAdjustmentRepo) ListStockAdjustments(ctx context.Context, companyID string) ([]domain.StockAdjustment, error) {
	var models []domain.StockAdjustmentGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.StockAdjustment, len(models))
	for i := range models {
		out[i] = *gormStockAdjToDomain(&models[i])
	}
	return out, nil
}

func (r *PGStockAdjustmentRepo) UpdateStockAdjustment(ctx context.Context, a *domain.StockAdjustment) error {
	return r.db.WithContext(ctx).Model(&domain.StockAdjustmentGORM{}).Where("id = ?", a.ID).Updates(map[string]interface{}{
		"warehouse_id": a.WarehouseID,
		"reason":       a.Reason,
		"status":       string(a.Status),
	}).Error
}

func (r *PGStockAdjustmentRepo) UpdateStockAdjustmentStatus(ctx context.Context, id string, status domain.AdjStatus) error {
	return r.db.WithContext(ctx).Model(&domain.StockAdjustmentGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGStockAdjustmentRepo) GetAdjustmentItems(ctx context.Context, adjustmentID string) ([]domain.AdjItem, error) {
	var models []domain.AdjItemGORM
	if err := r.db.WithContext(ctx).Where("adjustment_id = ?", adjustmentID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.AdjItem, len(models))
	for i := range models {
		out[i] = *gormAdjItemToDomain(&models[i])
	}
	return out, nil
}

func (r *PGStockAdjustmentRepo) CreateAdjustmentItem(ctx context.Context, it *domain.AdjItem) error {
	m := adjItemToGORM(it)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	it.ID = fmt.Sprint(m.ID)
	return nil
}

func stockAdjToGORM(a *domain.StockAdjustment) domain.StockAdjustmentGORM {
	reason := a.Reason
	if reason == "" {
		reason = "adjustment"
	}
	return domain.StockAdjustmentGORM{
		ID:        a.ID,
		CompanyID: a.CompanyID,
		WarehouseID: a.WarehouseID,
		Reason:    reason,
		Status:    string(a.Status),
		CreatedBy: a.CreatedBy,
	}
}

func gormStockAdjToDomain(m *domain.StockAdjustmentGORM) *domain.StockAdjustment {
	a := &domain.StockAdjustment{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		WarehouseID: m.WarehouseID,
		Reason:    m.Reason,
		Status:    domain.AdjStatus(m.Status),
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Items != nil {
		a.Items = make([]domain.AdjItem, len(m.Items))
		for i, it := range m.Items {
			a.Items[i] = *gormAdjItemToDomain(&it)
		}
	}
	return a
}

func adjItemToGORM(it *domain.AdjItem) domain.AdjItemGORM {
	return domain.AdjItemGORM{
		AdjustmentID:  it.AdjustmentID,
		ItemID:        it.ItemID,
		QtyOnHand:     it.QtyBefore,
		QtyCounted:    it.QtyAfter,
		QtyAdjustment: it.QtyAfter - it.QtyBefore,
		Reason:        nullStrG(it.Reason),
	}
}

func gormAdjItemToDomain(m *domain.AdjItemGORM) *domain.AdjItem {
	it := &domain.AdjItem{
		ID:           fmt.Sprint(m.ID),
		AdjustmentID: m.AdjustmentID,
		ItemID:       m.ItemID,
		QtyBefore:    m.QtyOnHand,
		QtyAfter:     m.QtyCounted,
	}
	if m.Reason != nil {
		it.Reason = *m.Reason
	}
	return it
}

// ─── Stock Take ─────────────────────────────────────────────────────────

type PGStockTakeRepo struct{ db *gorm.DB }

func NewPGStockTakeRepo(db *gorm.DB) *PGStockTakeRepo { return &PGStockTakeRepo{db} }

func (r *PGStockTakeRepo) CreateStockTake(ctx context.Context, t *domain.StockTake) error {
	m := stockTakeToGORM(t)
	items := make([]domain.TakeItemGORM, len(t.Items))
	for i, it := range t.Items {
		items[i] = takeItemToGORM(&it)
	}
	m.Items = items
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	t.ID = m.ID
	t.CreatedAt = m.CreatedAt
	t.UpdatedAt = m.UpdatedAt
	for i := range t.Items {
		if i < len(m.Items) {
			t.Items[i].ID = fmt.Sprint(m.Items[i].ID)
		}
	}
	return nil
}

func (r *PGStockTakeRepo) GetStockTakeByID(ctx context.Context, id string) (*domain.StockTake, error) {
	var m domain.StockTakeGORM
	if err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTakeNotFound
		}
		return nil, err
	}
	return gormStockTakeToDomain(&m), nil
}

func (r *PGStockTakeRepo) ListStockTakes(ctx context.Context, companyID string) ([]domain.StockTake, error) {
	var models []domain.StockTakeGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.StockTake, len(models))
	for i := range models {
		out[i] = *gormStockTakeToDomain(&models[i])
	}
	return out, nil
}

func (r *PGStockTakeRepo) UpdateStockTake(ctx context.Context, t *domain.StockTake) error {
	return r.db.WithContext(ctx).Model(&domain.StockTakeGORM{}).Where("id = ?", t.ID).Updates(map[string]interface{}{
		"warehouse_id": t.WarehouseID,
		"status":       string(t.Status),
	}).Error
}

func (r *PGStockTakeRepo) UpdateStockTakeStatus(ctx context.Context, id string, status domain.TakeStatus) error {
	return r.db.WithContext(ctx).Model(&domain.StockTakeGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGStockTakeRepo) GetTakeItems(ctx context.Context, takeID string) ([]domain.TakeItem, error) {
	var models []domain.TakeItemGORM
	if err := r.db.WithContext(ctx).Where("stock_take_id = ?", takeID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.TakeItem, len(models))
	for i := range models {
		out[i] = *gormTakeItemToDomain(&models[i])
	}
	return out, nil
}

func (r *PGStockTakeRepo) CreateTakeItem(ctx context.Context, it *domain.TakeItem) error {
	m := takeItemToGORM(it)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	it.ID = fmt.Sprint(m.ID)
	return nil
}

func stockTakeToGORM(t *domain.StockTake) domain.StockTakeGORM {
	return domain.StockTakeGORM{
		ID:        t.ID,
		CompanyID: t.CompanyID,
		WarehouseID: t.WarehouseID,
		Status:    string(t.Status),
		CreatedBy: t.CreatedBy,
	}
}

func gormStockTakeToDomain(m *domain.StockTakeGORM) *domain.StockTake {
	t := &domain.StockTake{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		WarehouseID: m.WarehouseID,
		Status:    domain.TakeStatus(m.Status),
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Items != nil {
		t.Items = make([]domain.TakeItem, len(m.Items))
		for i, it := range m.Items {
			t.Items[i] = *gormTakeItemToDomain(&it)
		}
	}
	return t
}

func takeItemToGORM(it *domain.TakeItem) domain.TakeItemGORM {
	return domain.TakeItemGORM{
		StockTakeID: it.TakeID,
		ItemID:      it.ItemID,
		QtyBook:     it.ExpectedQty,
		QtyCounted:  it.ActualQty,
		QtyVariance: it.Variance,
		Note:        nullStrG(it.Notes),
	}
}

func gormTakeItemToDomain(m *domain.TakeItemGORM) *domain.TakeItem {
	it := &domain.TakeItem{
		ID:          fmt.Sprint(m.ID),
		TakeID:      m.StockTakeID,
		ItemID:      m.ItemID,
		ExpectedQty: m.QtyBook,
		ActualQty:   m.QtyCounted,
		Variance:    m.QtyVariance,
	}
	if m.Note != nil {
		it.Notes = *m.Note
	}
	return it
}

// ─── Inventory Valuation Run ─────────────────────────────────────────────

type PGInventoryValuationRunRepo struct{ db *gorm.DB }

func NewPGInventoryValuationRunRepo(db *gorm.DB) *PGInventoryValuationRunRepo {
	return &PGInventoryValuationRunRepo{db}
}

func (r *PGInventoryValuationRunRepo) CreateValuationRun(ctx context.Context, v *domain.InventoryValuationRun) error {
	m := valRunToGORM(v)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	v.ID = m.ID
	v.CreatedAt = m.CreatedAt
	return nil
}

func (r *PGInventoryValuationRunRepo) GetValuationRunByID(ctx context.Context, id string) (*domain.InventoryValuationRun, error) {
	var m domain.InventoryValuationRunGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrValRunNotFound
		}
		return nil, err
	}
	return gormValRunToDomain(&m), nil
}

func (r *PGInventoryValuationRunRepo) ListValuationRuns(ctx context.Context, companyID string) ([]domain.InventoryValuationRun, error) {
	var models []domain.InventoryValuationRunGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.InventoryValuationRun, len(models))
	for i := range models {
		out[i] = *gormValRunToDomain(&models[i])
	}
	return out, nil
}

func (r *PGInventoryValuationRunRepo) UpdateValuationRun(ctx context.Context, v *domain.InventoryValuationRun) error {
	return r.db.WithContext(ctx).Model(&domain.InventoryValuationRunGORM{}).Where("id = ?", v.ID).Updates(map[string]interface{}{
		"method":   string(v.Method),
		"status":   string(v.Status),
	}).Error
}

func valRunToGORM(v *domain.InventoryValuationRun) domain.InventoryValuationRunGORM {
	return domain.InventoryValuationRunGORM{
		ID:         v.ID,
		CompanyID:  v.CompanyID,
		Method:     string(v.Method),
		Status:     string(v.Status),
		ExecutedBy: v.CreatedBy,
	}
}

func gormValRunToDomain(m *domain.InventoryValuationRunGORM) *domain.InventoryValuationRun {
	return &domain.InventoryValuationRun{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Method:    domain.ValuationMethod(m.Method),
		Status:    domain.ValuationRunStatus(m.Status),
		CreatedBy: m.ExecutedBy,
		CreatedAt: m.CreatedAt,
	}
}
