package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type PGWarehouseKeeperRepo struct{ db *gorm.DB }

func NewPGWarehouseKeeperRepo(db *gorm.DB) *PGWarehouseKeeperRepo {
	return &PGWarehouseKeeperRepo{db}
}

// ─── Assignment ─────────────────────────────────────────────────────────────

func (r *PGWarehouseKeeperRepo) CreateAssignment(ctx context.Context, a *domain.WarehouseKeeperAssignment) error {
	m := assignmentToGORM(a)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	a.ID = m.ID
	a.CreatedAt = m.CreatedAt
	a.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *PGWarehouseKeeperRepo) GetAssignment(ctx context.Context, id string) (*domain.WarehouseKeeperAssignment, error) {
	var m domain.WarehouseKeeperAssignmentGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWarehouseNotFound
		}
		return nil, err
	}
	return gormAssignmentToDomain(&m), nil
}

func (r *PGWarehouseKeeperRepo) ListAssignments(ctx context.Context, companyID string) ([]domain.WarehouseKeeperAssignment, error) {
	var models []domain.WarehouseKeeperAssignmentGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("effective_from DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.WarehouseKeeperAssignment, len(models))
	for i := range models {
		out[i] = *gormAssignmentToDomain(&models[i])
	}
	return out, nil
}

func (r *PGWarehouseKeeperRepo) GetActiveAssignment(ctx context.Context, companyID, warehouseID string, date time.Time) (*domain.WarehouseKeeperAssignment, error) {
	var m domain.WarehouseKeeperAssignmentGORM
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND warehouse_id = ? AND is_active = ? AND effective_from <= ? AND (effective_to IS NULL OR effective_to >= ?)",
			companyID, warehouseID, true, date, date).
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWarehouseNotFound
		}
		return nil, err
	}
	return gormAssignmentToDomain(&m), nil
}

func (r *PGWarehouseKeeperRepo) UpdateAssignment(ctx context.Context, a *domain.WarehouseKeeperAssignment) error {
	return r.db.WithContext(ctx).Model(&domain.WarehouseKeeperAssignmentGORM{}).Where("id = ?", a.ID).Updates(map[string]interface{}{
		"warehouse_id":   a.WarehouseID,
		"user_id":        a.UserID,
		"role":           a.Role,
		"effective_from": a.EffectiveFrom,
		"effective_to":   a.EffectiveTo,
		"is_active":      a.IsActive,
	}).Error
}

func (r *PGWarehouseKeeperRepo) DeleteAssignment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.WarehouseKeeperAssignmentGORM{}).Error
}

// ─── Stock Ledger ───────────────────────────────────────────────────────────

func (r *PGWarehouseKeeperRepo) CreateLedgerEntry(ctx context.Context, e *domain.StockLedgerEntry) error {
	m := ledgerToGORM(e)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	e.ID = m.ID
	e.RecordedAt = m.RecordedAt
	e.CreatedAt = m.CreatedAt
	return nil
}

func (r *PGWarehouseKeeperRepo) GetLedgerEntry(ctx context.Context, id string) (*domain.StockLedgerEntry, error) {
	var m domain.StockLedgerEntryGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWarehouseNotFound
		}
		return nil, err
	}
	return gormLedgerToDomain(&m), nil
}

func (r *PGWarehouseKeeperRepo) ListLedgerEntries(ctx context.Context, filter domain.LedgerFilter) ([]domain.StockLedgerEntry, int, error) {
	db := r.db.WithContext(ctx).Model(&domain.StockLedgerEntryGORM{}).
		Where("company_id = ? AND warehouse_id = ?", filter.CompanyID, filter.WarehouseID)

	if !filter.From.IsZero() {
		db = db.Where("entry_date >= ?", filter.From)
	}
	if !filter.To.IsZero() {
		db = db.Where("entry_date <= ?", filter.To)
	}
	if filter.ItemID != "" {
		db = db.Where("item_id = ?", filter.ItemID)
	}
	if filter.VoucherType != "" {
		db = db.Where("voucher_type = ?", filter.VoucherType)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	var models []domain.StockLedgerEntryGORM
	if err := db.Order("entry_date DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	out := make([]domain.StockLedgerEntry, len(models))
	for i := range models {
		out[i] = *gormLedgerToDomain(&models[i])
	}
	return out, int(total), nil
}

func (r *PGWarehouseKeeperRepo) UnrecordLedgerEntry(ctx context.Context, id string, unrecordedBy string, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.StockLedgerEntryGORM{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":          domain.LedgerStatusUnrecorded,
		"unrecorded_by":   unrecordedBy,
		"unrecorded_at":   now,
		"unrecord_reason": reason,
	}).Error
}

func (r *PGWarehouseKeeperRepo) GetLedgerBalance(ctx context.Context, companyID, warehouseID, itemID string) (float64, error) {
	var result struct {
		Balance float64
	}
	if err := r.db.WithContext(ctx).
		Table("stock_ledger_entries").
		Select("COALESCE(MAX(balance_qty), 0)").
		Where("company_id = ? AND warehouse_id = ? AND item_id = ? AND status = ?",
			companyID, warehouseID, itemID, domain.LedgerStatusRecorded).
		Scan(&result).Error; err != nil {
		return 0, err
	}
	return result.Balance, nil
}

// ─── Reconciliation ─────────────────────────────────────────────────────────

func (r *PGWarehouseKeeperRepo) GetReconciliationReport(ctx context.Context, companyID, warehouseID string, from, to time.Time) ([]domain.KeeperReconciliationItem, error) {
	type ledgerAgg struct {
		ItemID   string
		Balance  float64
		MaxDate  time.Time
	}
	var ledgerItems []ledgerAgg
	if err := r.db.WithContext(ctx).
		Table("stock_ledger_entries").
		Select("item_id, MAX(balance_qty) as balance, MAX(entry_date) as max_date").
		Where("company_id = ? AND warehouse_id = ? AND entry_date >= ? AND entry_date <= ? AND status = ?",
			companyID, warehouseID, from, to, domain.LedgerStatusRecorded).
		Group("item_id").
		Scan(&ledgerItems).Error; err != nil {
		return nil, err
	}

	type balanceAgg struct {
		ItemID    string
		QtyOnHand float64
		UnitCost  float64
	}
	var balItems []balanceAgg
	if err := r.db.WithContext(ctx).
		Table("stock_balances").
		Select("item_id, qty_on_hand, unit_cost").
		Where("company_id = ? AND warehouse_id = ?", companyID, warehouseID).
		Scan(&balItems).Error; err != nil {
		return nil, err
	}

	balMap := make(map[string]balanceAgg)
	for _, b := range balItems {
		balMap[b.ItemID] = b
	}

	var out []domain.KeeperReconciliationItem
	for _, li := range ledgerItems {
		item := domain.KeeperReconciliationItem{
			ItemID:        li.ItemID,
			WarehouseID:   warehouseID,
			KeeperQty:     li.Balance,
			LastUpdated:   li.MaxDate,
		}
		if b, ok := balMap[li.ItemID]; ok {
			item.AccountingQty = b.QtyOnHand
			item.UnitCost = b.UnitCost
		}
		item.VarianceQty = item.KeeperQty - item.AccountingQty
		item.VarianceValue = item.VarianceQty * item.UnitCost
		out = append(out, item)
	}
	return out, nil
}

// ─── Stock Card ─────────────────────────────────────────────────────────────

func (r *PGWarehouseKeeperRepo) GetStockCard(ctx context.Context, companyID, warehouseID, itemID string, period string) (*domain.StockCard, error) {
	periodStart := period + "-01"
	var lines []domain.StockCardLine
	if err := r.db.WithContext(ctx).
		Table("stock_ledger_entries").
		Select("entry_date as date, voucher_no, description, receipt_qty, issue_qty, balance_qty").
		Where("company_id = ? AND warehouse_id = ? AND item_id = ? AND entry_date >= ? AND status = ?",
			companyID, warehouseID, itemID, periodStart, domain.LedgerStatusRecorded).
		Order("entry_date ASC, created_at ASC").
		Scan(&lines).Error; err != nil {
		return nil, err
	}

	card := &domain.StockCard{
		WarehouseID: warehouseID,
		ItemID:      itemID,
		Period:      period,
		Lines:       lines,
	}

	if len(lines) > 0 {
		card.OpeningBalance = lines[0].BalanceQty - lines[0].ReceiptQty
		card.ClosingBalance = lines[len(lines)-1].BalanceQty
	}

	return card, nil
}

// ─── Keeper Reports ─────────────────────────────────────────────────────────

func (r *PGWarehouseKeeperRepo) GetKeeperInventorySummary(ctx context.Context, companyID, warehouseID string) ([]domain.KeeperInventorySummaryItem, error) {
	type agg struct {
		ItemID    string
		Quantity  float64
		UnitCost  float64
		TotalVal  float64
		MaxDate   time.Time
	}
	var items []agg
	if err := r.db.WithContext(ctx).
		Table("stock_ledger_entries").
		Select("item_id, MAX(balance_qty) as quantity, MAX(unit_cost) as unit_cost, MAX(total_value) as total_val, MAX(entry_date) as max_date").
		Where("company_id = ? AND warehouse_id = ? AND status = ?", companyID, warehouseID, domain.LedgerStatusRecorded).
		Group("item_id").
		Scan(&items).Error; err != nil {
		return nil, err
	}

	out := make([]domain.KeeperInventorySummaryItem, len(items))
	for i, it := range items {
		out[i] = domain.KeeperInventorySummaryItem{
			ItemID:     it.ItemID,
			Quantity:   it.Quantity,
			UnitCost:   it.UnitCost,
			TotalValue: it.TotalVal,
			LastUpdated: it.MaxDate,
		}
	}
	return out, nil
}

// ─── GORM converters ────────────────────────────────────────────────────────

func assignmentToGORM(a *domain.WarehouseKeeperAssignment) domain.WarehouseKeeperAssignmentGORM {
	return domain.WarehouseKeeperAssignmentGORM{
		ID:            a.ID,
		CompanyID:     a.CompanyID,
		WarehouseID:   a.WarehouseID,
		UserID:        a.UserID,
		Role:          string(a.Role),
		EffectiveFrom: a.EffectiveFrom,
		EffectiveTo:   a.EffectiveTo,
		IsActive:      a.IsActive,
		CreatedBy:     a.CreatedBy,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func gormAssignmentToDomain(m *domain.WarehouseKeeperAssignmentGORM) *domain.WarehouseKeeperAssignment {
	return &domain.WarehouseKeeperAssignment{
		ID:            m.ID,
		CompanyID:     m.CompanyID,
		WarehouseID:   m.WarehouseID,
		UserID:        m.UserID,
		Role:          domain.KeeperRole(m.Role),
		EffectiveFrom: m.EffectiveFrom,
		EffectiveTo:   m.EffectiveTo,
		IsActive:      m.IsActive,
		CreatedBy:     m.CreatedBy,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ledgerToGORM(e *domain.StockLedgerEntry) domain.StockLedgerEntryGORM {
	m := domain.StockLedgerEntryGORM{
		ID:           e.ID,
		CompanyID:    e.CompanyID,
		WarehouseID:  e.WarehouseID,
		ItemID:       e.ItemID,
		EntryDate:    e.EntryDate,
		VoucherType:  string(e.VoucherType),
		ReceiptQty:   e.ReceiptQty,
		IssueQty:     e.IssueQty,
		BalanceQty:   e.BalanceQty,
		RecordedBy:   e.RecordedBy,
		RecordedAt:   e.RecordedAt,
		Status:       string(e.Status),
		CreatedAt:    e.CreatedAt,
	}
	if e.VoucherNo != "" {
		m.VoucherNo = &e.VoucherNo
	}
	if e.VoucherRefID != "" {
		m.VoucherRefID = &e.VoucherRefID
	}
	if e.Description != "" {
		m.Description = &e.Description
	}
	if e.UnitCost != 0 {
		m.UnitCost = &e.UnitCost
	}
	if e.TotalValue != 0 {
		m.TotalValue = &e.TotalValue
	}
	if e.UnrecordedBy != "" {
		m.UnrecordedBy = &e.UnrecordedBy
	}
	if e.UnrecordedAt != nil {
		m.UnrecordedAt = e.UnrecordedAt
	}
	if e.UnrecordReason != "" {
		m.UnrecordReason = &e.UnrecordReason
	}
	return m
}

func gormLedgerToDomain(m *domain.StockLedgerEntryGORM) *domain.StockLedgerEntry {
	e := &domain.StockLedgerEntry{
		ID:           m.ID,
		CompanyID:    m.CompanyID,
		WarehouseID:  m.WarehouseID,
		ItemID:       m.ItemID,
		EntryDate:    m.EntryDate,
		VoucherType:  domain.LedgerVoucherType(m.VoucherType),
		ReceiptQty:   m.ReceiptQty,
		IssueQty:     m.IssueQty,
		BalanceQty:   m.BalanceQty,
		RecordedBy:   m.RecordedBy,
		RecordedAt:   m.RecordedAt,
		Status:       domain.LedgerEntryStatus(m.Status),
		CreatedAt:    m.CreatedAt,
	}
	if m.VoucherNo != nil {
		e.VoucherNo = *m.VoucherNo
	}
	if m.VoucherRefID != nil {
		e.VoucherRefID = *m.VoucherRefID
	}
	if m.Description != nil {
		e.Description = *m.Description
	}
	if m.UnitCost != nil {
		e.UnitCost = *m.UnitCost
	}
	if m.TotalValue != nil {
		e.TotalValue = *m.TotalValue
	}
	if m.UnrecordedBy != nil {
		e.UnrecordedBy = *m.UnrecordedBy
	}
	if m.UnrecordedAt != nil {
		e.UnrecordedAt = m.UnrecordedAt
	}
	if m.UnrecordReason != nil {
		e.UnrecordReason = *m.UnrecordReason
	}
	return e
}
