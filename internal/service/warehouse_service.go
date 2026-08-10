package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type JournalEntryService interface {
	CreateEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error
	PostEntry(ctx context.Context, id string) error
}

type WarehouseService struct {
	whRepo   domain.WarehouseRepository
	catRepo  domain.ItemCategoryRepository
	itemRepo domain.ItemRepository
	balRepo  domain.StockBalanceRepository
	txnRepo  domain.InventoryTransactionRepository
	trfRepo  domain.StockTransferRepository
	adjRepo  domain.StockAdjustmentRepository
	takeRepo domain.StockTakeRepository
	valRepo  domain.InventoryValuationRunRepository
	grnRepo  domain.GRNRepository
	jeSvc    JournalEntryService
	now      func() time.Time
}

func NewWarehouseService(
	whRepo domain.WarehouseRepository,
	catRepo domain.ItemCategoryRepository,
	itemRepo domain.ItemRepository,
	balRepo domain.StockBalanceRepository,
	txnRepo domain.InventoryTransactionRepository,
	trfRepo domain.StockTransferRepository,
	adjRepo domain.StockAdjustmentRepository,
	takeRepo domain.StockTakeRepository,
	valRepo domain.InventoryValuationRunRepository,
	grnRepo domain.GRNRepository,
	jeSvc JournalEntryService,
) *WarehouseService {
	return &WarehouseService{
		whRepo: whRepo, catRepo: catRepo, itemRepo: itemRepo,
		balRepo: balRepo, txnRepo: txnRepo, trfRepo: trfRepo,
		adjRepo: adjRepo, takeRepo: takeRepo, valRepo: valRepo,
		grnRepo: grnRepo, jeSvc: jeSvc,
		now: time.Now,
	}
}

func periodFromTime(t time.Time) string {
	return fmt.Sprintf("%d%02d", t.Year(), t.Month())
}

func (s *WarehouseService) adjustStockBalance(ctx context.Context, companyID, warehouseID, itemID, period string, qtyDelta, costDelta float64, transType domain.TransactionType, refType, refID, createdBy string) error {
	bal, err := s.balRepo.FindStockBalance(ctx, companyID, warehouseID, itemID, period)
	if err != nil {
		newID := fmt.Sprintf("BAL-%s-%s-%s-%s", companyID, warehouseID, itemID, period)
		bal = &domain.StockBalance{
			ID: newID, CompanyID: companyID, WarehouseID: warehouseID,
			ItemID: itemID, Period: period,
		}
	}
	oldQty := bal.Quantity

	bal.Quantity += qtyDelta
	bal.TotalCost += costDelta
	if bal.Quantity > 0 {
		bal.UnitCost = bal.TotalCost / bal.Quantity
	} else if bal.Quantity == 0 {
		bal.UnitCost = 0
	}
	now := s.now()
	bal.LastTransactionAt = now
	bal.UpdatedAt = now
	if err := s.balRepo.UpsertStockBalance(ctx, bal); err != nil {
		return err
	}

	txn := &domain.InventoryTransaction{
		CompanyID: companyID, WarehouseID: warehouseID, ItemID: itemID,
		TransType: transType, RefType: refType, RefID: refID,
		QtyBefore: oldQty, Quantity: qtyDelta, QtyAfter: bal.Quantity,
		UnitCost: bal.UnitCost, TotalCost: costDelta,
		CreatedBy: createdBy, CreatedAt: now,
	}
	return s.txnRepo.CreateInventoryTransaction(ctx, txn)
}

func (s *WarehouseService) postGLEntry(ctx context.Context, companyID, description string, entryDate time.Time, lines []domain.JournalLine, voucherType domain.VoucherType, userID string) error {
	if s.jeSvc == nil {
		return nil
	}
	entry := &domain.JournalEntry{
		CompanyID: companyID, EntryDate: entryDate,
		Description: description, VoucherType: voucherType,
		CurrencyCode: "VND", ExchangeRate: 1, Lines: lines,
	}
	if err := s.jeSvc.CreateEntry(ctx, entry, userID); err != nil {
		return err
	}
	return s.jeSvc.PostEntry(ctx, entry.ID)
}

const defaultInventoryAccount = "152"
const transitAccount = "151"
const cogsAccount = "632"
const surplusAccount = "3381"
const otherExpenseAccount = "811"
const otherIncomeAccount = "711"
const adminExpenseAccount = "642"
const pettyCashAccount = "111"

// ─── Warehouse ────────────────────────────────────────────────────────

func (s *WarehouseService) CreateWarehouse(ctx context.Context, w *domain.Warehouse) error {
	if err := w.Validate(); err != nil { return err }
	existing, _ := s.whRepo.GetWarehouseByCode(ctx, w.CompanyID, w.Code)
	if existing != nil { return domain.ErrWarehouseCodeExists }
	w.CreatedAt = s.now()
	w.UpdatedAt = w.CreatedAt
	return s.whRepo.CreateWarehouse(ctx, w)
}

func (s *WarehouseService) GetWarehouse(ctx context.Context, id string) (*domain.Warehouse, error) {
	return s.whRepo.GetWarehouseByID(ctx, id)
}

func (s *WarehouseService) ListWarehouses(ctx context.Context, companyID string) ([]domain.Warehouse, error) {
	return s.whRepo.ListWarehouses(ctx, companyID)
}

func (s *WarehouseService) UpdateWarehouse(ctx context.Context, w *domain.Warehouse) error {
	existing, err := s.whRepo.GetWarehouseByID(ctx, w.ID)
	if err != nil { return err }
	if existing.Code != w.Code {
		dup, _ := s.whRepo.GetWarehouseByCode(ctx, w.CompanyID, w.Code)
		if dup != nil { return domain.ErrWarehouseCodeExists }
	}
	w.UpdatedAt = s.now()
	return s.whRepo.UpdateWarehouse(ctx, w)
}

func (s *WarehouseService) DeleteWarehouse(ctx context.Context, id string) error {
	return s.whRepo.DeleteWarehouse(ctx, id)
}

// ─── Item Category ────────────────────────────────────────────────────

func (s *WarehouseService) CreateCategory(ctx context.Context, c *domain.ItemCategory) error {
	if err := c.Validate(); err != nil { return err }
	existing, _ := s.catRepo.GetCategoryByCode(ctx, c.CompanyID, c.Code)
	if existing != nil { return domain.ErrCategoryCodeExists }
	c.CreatedAt = s.now()
	c.UpdatedAt = c.CreatedAt
	return s.catRepo.CreateCategory(ctx, c)
}

func (s *WarehouseService) GetCategory(ctx context.Context, id string) (*domain.ItemCategory, error) {
	return s.catRepo.GetCategoryByID(ctx, id)
}

func (s *WarehouseService) ListCategories(ctx context.Context, companyID string) ([]domain.ItemCategory, error) {
	return s.catRepo.ListCategories(ctx, companyID)
}

func (s *WarehouseService) UpdateCategory(ctx context.Context, c *domain.ItemCategory) error {
	existing, err := s.catRepo.GetCategoryByID(ctx, c.ID)
	if err != nil { return err }
	if existing.Code != c.Code {
		dup, _ := s.catRepo.GetCategoryByCode(ctx, c.CompanyID, c.Code)
		if dup != nil { return domain.ErrCategoryCodeExists }
	}
	c.UpdatedAt = s.now()
	return s.catRepo.UpdateCategory(ctx, c)
}

func (s *WarehouseService) DeleteCategory(ctx context.Context, id string) error {
	return s.catRepo.DeleteCategory(ctx, id)
}

// ─── Item ─────────────────────────────────────────────────────────────

func (s *WarehouseService) CreateItem(ctx context.Context, i *domain.Item) error {
	if err := i.Validate(); err != nil { return err }
	existing, _ := s.itemRepo.GetItemByCode(ctx, i.CompanyID, i.Code)
	if existing != nil { return domain.ErrItemCodeExists }
	i.CreatedAt = s.now()
	i.UpdatedAt = i.CreatedAt
	return s.itemRepo.CreateItem(ctx, i)
}

func (s *WarehouseService) GetItem(ctx context.Context, id string) (*domain.Item, error) {
	return s.itemRepo.GetItemByID(ctx, id)
}

func (s *WarehouseService) ListItems(ctx context.Context, companyID string) ([]domain.Item, error) {
	return s.itemRepo.ListItems(ctx, companyID)
}

func (s *WarehouseService) UpdateItem(ctx context.Context, i *domain.Item) error {
	existing, err := s.itemRepo.GetItemByID(ctx, i.ID)
	if err != nil { return err }
	if existing.Code != i.Code {
		dup, _ := s.itemRepo.GetItemByCode(ctx, i.CompanyID, i.Code)
		if dup != nil { return domain.ErrItemCodeExists }
	}
	i.UpdatedAt = s.now()
	return s.itemRepo.UpdateItem(ctx, i)
}

func (s *WarehouseService) DeleteItem(ctx context.Context, id string) error {
	return s.itemRepo.DeleteItem(ctx, id)
}

// ─── Stock Balance ────────────────────────────────────────────────────

func (s *WarehouseService) GetStockBalance(ctx context.Context, id string) (*domain.StockBalance, error) {
	return s.balRepo.GetStockBalanceByID(ctx, id)
}

func (s *WarehouseService) FindStockBalance(ctx context.Context, companyID, warehouseID, itemID, period string) (*domain.StockBalance, error) {
	return s.balRepo.FindStockBalance(ctx, companyID, warehouseID, itemID, period)
}

func (s *WarehouseService) ListStockBalances(ctx context.Context, companyID, warehouseID string) ([]domain.StockBalance, error) {
	return s.balRepo.ListStockBalances(ctx, companyID, warehouseID)
}

func (s *WarehouseService) UpsertStockBalance(ctx context.Context, b *domain.StockBalance) error {
	return s.balRepo.UpsertStockBalance(ctx, b)
}

// ─── Inventory Transaction ────────────────────────────────────────────

func (s *WarehouseService) ListInventoryTransactions(ctx context.Context, companyID, warehouseID, itemID string, offset, limit int) ([]domain.InventoryTransaction, int, error) {
	return s.txnRepo.ListInventoryTransactions(ctx, companyID, warehouseID, itemID, offset, limit)
}

func (s *WarehouseService) CreateInventoryTransaction(ctx context.Context, t *domain.InventoryTransaction) error {
	return s.txnRepo.CreateInventoryTransaction(ctx, t)
}

// ─── Stock Transfer ───────────────────────────────────────────────────

func (s *WarehouseService) CreateStockTransfer(ctx context.Context, t *domain.StockTransfer) error {
	if err := t.Validate(); err != nil { return err }
	t.Status = domain.TransStatusDraft
	t.CreatedAt = s.now()
	t.UpdatedAt = t.CreatedAt
	return s.trfRepo.CreateStockTransfer(ctx, t)
}

func (s *WarehouseService) GetStockTransfer(ctx context.Context, id string) (*domain.StockTransfer, error) {
	trf, err := s.trfRepo.GetStockTransferByID(ctx, id)
	if err != nil { return nil, err }
	items, err := s.trfRepo.GetTransferItems(ctx, id)
	if err != nil { return nil, err }
	trf.Items = items
	return trf, nil
}

func (s *WarehouseService) ListStockTransfers(ctx context.Context, companyID string) ([]domain.StockTransfer, error) {
	return s.trfRepo.ListStockTransfers(ctx, companyID)
}

func (s *WarehouseService) UpdateStockTransfer(ctx context.Context, t *domain.StockTransfer) error {
	existing, err := s.trfRepo.GetStockTransferByID(ctx, t.ID)
	if err != nil { return err }
	if existing.Status != domain.TransStatusDraft {
		return domain.ErrTransferInvalidStatus
	}
	t.Status = existing.Status
	t.UpdatedAt = s.now()
	return s.trfRepo.UpdateStockTransfer(ctx, t)
}

func (s *WarehouseService) SubmitStockTransfer(ctx context.Context, id string) error {
	t, err := s.trfRepo.GetStockTransferByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TransStatusPending) {
		return domain.ErrTransferInvalidStatus
	}
	t.Status = domain.TransStatusPending
	t.UpdatedAt = s.now()
	return s.trfRepo.UpdateStockTransfer(ctx, t)
}

func (s *WarehouseService) ApproveStockTransfer(ctx context.Context, id, approvedBy string) error {
	t, err := s.trfRepo.GetStockTransferByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TransStatusApproved) {
		return domain.ErrTransferInvalidStatus
	}
	t.Status = domain.TransStatusApproved
	t.ApprovedBy = approvedBy
	t.ApprovedAt = s.now().Format("2006-01-02 15:04:05")
	t.UpdatedAt = s.now()
	return s.trfRepo.UpdateStockTransfer(ctx, t)
}

func (s *WarehouseService) TransferStockTransfer(ctx context.Context, id string) error {
	t, err := s.trfRepo.GetStockTransferByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TransStatusTransferred) {
		return domain.ErrTransferInvalidStatus
	}
	items, err := s.trfRepo.GetTransferItems(ctx, id)
	if err != nil { return fmt.Errorf("get items: %w", err) }
	now := s.now()
	period := periodFromTime(now)

	for _, item := range items {
		cost := item.Quantity * item.UnitCost
		if err := s.adjustStockBalance(ctx, t.CompanyID, t.FromWarehouseID, item.ItemID, period,
			-item.Quantity, -cost, domain.TransTransferOut, "stock_transfer", id, ""); err != nil {
			return fmt.Errorf("adjust source balance: %w", err)
		}
	}

	if s.jeSvc != nil {
		var lines []domain.JournalLine
		totalCost := 0.0
		for _, item := range items {
			cost := item.Quantity * item.UnitCost
			totalCost += cost
			lines = append(lines, domain.JournalLine{
				AccountCode: transitAccount, DebitAmount: totalCost,
				Description: fmt.Sprintf("Transfer %s → %s item %s qty %.0f", t.FromWarehouseID, t.ToWarehouseID, item.ItemID, item.Quantity),
			}, domain.JournalLine{
				AccountCode: defaultInventoryAccount, CreditAmount: totalCost,
			})
		}
		if len(lines) > 0 {
			lines[0].DebitAmount = totalCost
			lines = append(lines, domain.JournalLine{
				AccountCode: defaultInventoryAccount, CreditAmount: totalCost,
				Description: fmt.Sprintf("Transfer out WH %s", t.FromWarehouseID),
			})
			if err := s.postGLEntry(ctx, t.CompanyID,
				fmt.Sprintf("Stock transfer %s: %s → %s", t.TransferNumber, t.FromWarehouseID, t.ToWarehouseID),
				now, lines, domain.VoucherTypeInventoryIssue, ""); err != nil {
				return fmt.Errorf("post GL: %w", err)
			}
		}
	}

	t.Status = domain.TransStatusTransferred
	t.UpdatedAt = now
	return s.trfRepo.UpdateStockTransfer(ctx, t)
}

func (s *WarehouseService) CompleteStockTransfer(ctx context.Context, id, completedBy string) error {
	t, err := s.trfRepo.GetStockTransferByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TransStatusCompleted) {
		return domain.ErrTransferInvalidStatus
	}
	items, err := s.trfRepo.GetTransferItems(ctx, id)
	if err != nil { return fmt.Errorf("get items: %w", err) }
	now := s.now()
	period := periodFromTime(now)

	for _, item := range items {
		cost := item.Quantity * item.UnitCost
		if err := s.adjustStockBalance(ctx, t.CompanyID, t.ToWarehouseID, item.ItemID, period,
			item.Quantity, cost, domain.TransTransferIn, "stock_transfer", id, ""); err != nil {
			return fmt.Errorf("adjust dest balance: %w", err)
		}
	}

	if s.jeSvc != nil {
		totalCost := 0.0
		for _, item := range items {
			totalCost += item.Quantity * item.UnitCost
		}
		lines := []domain.JournalLine{
			{AccountCode: defaultInventoryAccount, DebitAmount: totalCost,
				Description: fmt.Sprintf("Receive transfer %s at WH %s", t.TransferNumber, t.ToWarehouseID)},
			{AccountCode: transitAccount, CreditAmount: totalCost,
				Description: fmt.Sprintf("Clear transit for %s", t.TransferNumber)},
		}
		if err := s.postGLEntry(ctx, t.CompanyID,
			fmt.Sprintf("Transfer complete %s → %s", t.FromWarehouseID, t.ToWarehouseID),
			now, lines, domain.VoucherTypeInventoryReceipt, ""); err != nil {
			return fmt.Errorf("post GL: %w", err)
		}
	}

	t.Status = domain.TransStatusCompleted
	t.CompletedBy = completedBy
	t.CompletedAt = now.Format("2006-01-02 15:04:05")
	t.UpdatedAt = now
	return s.trfRepo.UpdateStockTransfer(ctx, t)
}

func (s *WarehouseService) CancelStockTransfer(ctx context.Context, id, reason string) error {
	t, err := s.trfRepo.GetStockTransferByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TransStatusCancelled) {
		return domain.ErrTransferInvalidStatus
	}
	t.Status = domain.TransStatusCancelled
	t.CancelledReason = reason
	t.UpdatedAt = s.now()
	return s.trfRepo.UpdateStockTransfer(ctx, t)
}

// ─── Stock Adjustment ─────────────────────────────────────────────────

func (s *WarehouseService) CreateStockAdjustment(ctx context.Context, a *domain.StockAdjustment) error {
	if err := a.Validate(); err != nil { return err }
	a.Status = domain.AdjStatusDraft
	a.CreatedAt = s.now()
	a.UpdatedAt = a.CreatedAt
	return s.adjRepo.CreateStockAdjustment(ctx, a)
}

func (s *WarehouseService) GetStockAdjustment(ctx context.Context, id string) (*domain.StockAdjustment, error) {
	a, err := s.adjRepo.GetStockAdjustmentByID(ctx, id)
	if err != nil { return nil, err }
	items, err := s.adjRepo.GetAdjustmentItems(ctx, id)
	if err != nil { return nil, err }
	a.Items = items
	return a, nil
}

func (s *WarehouseService) ListStockAdjustments(ctx context.Context, companyID string) ([]domain.StockAdjustment, error) {
	return s.adjRepo.ListStockAdjustments(ctx, companyID)
}

func (s *WarehouseService) UpdateStockAdjustment(ctx context.Context, a *domain.StockAdjustment) error {
	existing, err := s.adjRepo.GetStockAdjustmentByID(ctx, a.ID)
	if err != nil { return err }
	if existing.Status != domain.AdjStatusDraft {
		return domain.ErrAdjInvalidStatus
	}
	a.Status = existing.Status
	a.UpdatedAt = s.now()
	return s.adjRepo.UpdateStockAdjustment(ctx, a)
}

func (s *WarehouseService) SubmitStockAdjustment(ctx context.Context, id string) error {
	a, err := s.adjRepo.GetStockAdjustmentByID(ctx, id)
	if err != nil { return err }
	if !a.Status.ValidTransition(domain.AdjStatusPending) {
		return domain.ErrAdjInvalidStatus
	}
	a.Status = domain.AdjStatusPending
	a.UpdatedAt = s.now()
	return s.adjRepo.UpdateStockAdjustment(ctx, a)
}

func (s *WarehouseService) ApproveStockAdjustment(ctx context.Context, id, approvedBy string) error {
	a, err := s.adjRepo.GetStockAdjustmentByID(ctx, id)
	if err != nil { return err }
	if !a.Status.ValidTransition(domain.AdjStatusApproved) {
		return domain.ErrAdjInvalidStatus
	}
	a.Status = domain.AdjStatusApproved
	a.ApprovedBy = approvedBy
	a.ApprovedAt = s.now().Format("2006-01-02 15:04:05")
	a.UpdatedAt = s.now()
	return s.adjRepo.UpdateStockAdjustment(ctx, a)
}

func (s *WarehouseService) PostStockAdjustment(ctx context.Context, id string) error {
	a, err := s.adjRepo.GetStockAdjustmentByID(ctx, id)
	if err != nil { return err }
	if !a.Status.ValidTransition(domain.AdjStatusPosted) {
		return domain.ErrAdjInvalidStatus
	}
	items, err := s.adjRepo.GetAdjustmentItems(ctx, id)
	if err != nil { return fmt.Errorf("get items: %w", err) }
	now := s.now()
	period := periodFromTime(now)

	var glLines []domain.JournalLine
	totalDebit := 0.0
	totalCredit := 0.0

	for _, item := range items {
		qtyDelta := item.QtyAfter - item.QtyBefore
		if qtyDelta == 0 {
			continue
		}
		costDelta := qtyDelta * item.UnitCost
		transType := domain.TransAdjIn
		if qtyDelta < 0 {
			transType = domain.TransAdjOut
		}

		if err := s.adjustStockBalance(ctx, a.CompanyID, a.WarehouseID, item.ItemID, period,
			qtyDelta, costDelta, transType, "stock_adjustment", id, a.CreatedBy); err != nil {
			return fmt.Errorf("adjust balance: %w", err)
		}

		if s.jeSvc != nil {
			if qtyDelta > 0 {
				glLines = append(glLines, domain.JournalLine{
					AccountCode: defaultInventoryAccount, DebitAmount: costDelta,
					Description: fmt.Sprintf("Adj increase item %s qty +%.0f", item.ItemID, qtyDelta),
				})
				debited := costDelta
				if totalDebit > 0 {
					totalCredit += debited
					glLines = append(glLines, domain.JournalLine{
						AccountCode: surplusAccount, CreditAmount: debited,
						Description: item.Reason,
					})
				} else {
					totalDebit += debited
				}
			} else {
				loss := -costDelta
				glLines = append(glLines, domain.JournalLine{
					AccountCode: cogsAccount, DebitAmount: loss,
					Description: fmt.Sprintf("Adj decrease item %s qty %.0f", item.ItemID, -qtyDelta),
				})
				totalCredit += loss
			}
		}
	}

	if s.jeSvc != nil && len(glLines) > 0 {
		if totalDebit > 0 && totalCredit == 0 {
			glLines = append(glLines, domain.JournalLine{
				AccountCode: surplusAccount, CreditAmount: totalDebit,
				Description: "Surplus awaiting resolution",
			})
		} else if totalCredit > 0 {
			// Add inventory credit for decreases
			for i, line := range glLines {
				if line.DebitAmount > 0 && line.AccountCode == cogsAccount {
					glLines[i].DebitAmount = totalCredit
				}
			}
			glLines = append(glLines, domain.JournalLine{
				AccountCode: defaultInventoryAccount, CreditAmount: totalCredit,
				Description: "Adjust decrease inventory",
			})
		}
		// Rebuild balanced lines
		balancedLines := rebuildBalancedLines(glLines)
		if err := s.postGLEntry(ctx, a.CompanyID,
			fmt.Sprintf("Stock adjustment %s: %s", a.AdjustmentNumber, string(a.AdjType)),
			now, balancedLines, domain.VoucherTypeOther, a.CreatedBy); err != nil {
			return fmt.Errorf("post GL: %w", err)
		}
	}

	a.Status = domain.AdjStatusPosted
	a.PostedAt = now.Format("2006-01-02 15:04:05")
	a.UpdatedAt = now
	return s.adjRepo.UpdateStockAdjustment(ctx, a)
}

func rebuildBalancedLines(lines []domain.JournalLine) []domain.JournalLine {
	var debitTotal, creditTotal float64
	for _, l := range lines {
		debitTotal += l.DebitAmount
		creditTotal += l.CreditAmount
	}
	if debitTotal == creditTotal {
		return lines
	}
	if debitTotal > creditTotal {
		lines = append(lines, domain.JournalLine{
			AccountCode: surplusAccount, CreditAmount: debitTotal - creditTotal,
			Description: "Balancing entry",
		})
	} else {
		lines = append(lines, domain.JournalLine{
			AccountCode: defaultInventoryAccount, DebitAmount: creditTotal - debitTotal,
			Description: "Balancing entry",
		})
	}
	return lines
}

func (s *WarehouseService) RejectStockAdjustment(ctx context.Context, id, reason string) error {
	a, err := s.adjRepo.GetStockAdjustmentByID(ctx, id)
	if err != nil { return err }
	if !a.Status.ValidTransition(domain.AdjStatusRejected) {
		return domain.ErrAdjInvalidStatus
	}
	a.Status = domain.AdjStatusRejected
	a.RejectedReason = reason
	a.UpdatedAt = s.now()
	return s.adjRepo.UpdateStockAdjustment(ctx, a)
}

// ─── Stock Take ───────────────────────────────────────────────────────

func (s *WarehouseService) CreateStockTake(ctx context.Context, t *domain.StockTake) error {
	if err := t.Validate(); err != nil { return err }
	t.Status = domain.TakeStatusPlanning
	t.CreatedAt = s.now()
	t.UpdatedAt = t.CreatedAt
	return s.takeRepo.CreateStockTake(ctx, t)
}

func (s *WarehouseService) GetStockTake(ctx context.Context, id string) (*domain.StockTake, error) {
	tk, err := s.takeRepo.GetStockTakeByID(ctx, id)
	if err != nil { return nil, err }
	items, err := s.takeRepo.GetTakeItems(ctx, id)
	if err != nil { return nil, err }
	tk.Items = items
	return tk, nil
}

func (s *WarehouseService) ListStockTakes(ctx context.Context, companyID string) ([]domain.StockTake, error) {
	return s.takeRepo.ListStockTakes(ctx, companyID)
}

func (s *WarehouseService) UpdateStockTake(ctx context.Context, t *domain.StockTake) error {
	existing, err := s.takeRepo.GetStockTakeByID(ctx, t.ID)
	if err != nil { return err }
	if existing.Status != domain.TakeStatusPlanning && existing.Status != domain.TakeStatusInProgress {
		return domain.ErrTakeInvalidStatus
	}
	t.Status = existing.Status
	t.UpdatedAt = s.now()
	return s.takeRepo.UpdateStockTake(ctx, t)
}

func (s *WarehouseService) StartStockTake(ctx context.Context, id string) error {
	t, err := s.takeRepo.GetStockTakeByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TakeStatusInProgress) {
		return domain.ErrTakeInvalidStatus
	}
	t.Status = domain.TakeStatusInProgress
	t.UpdatedAt = s.now()
	return s.takeRepo.UpdateStockTake(ctx, t)
}

func (s *WarehouseService) CompleteStockTake(ctx context.Context, id string) error {
	t, err := s.takeRepo.GetStockTakeByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TakeStatusCompleted) {
		return domain.ErrTakeInvalidStatus
	}
	t.Status = domain.TakeStatusCompleted
	t.UpdatedAt = s.now()
	return s.takeRepo.UpdateStockTake(ctx, t)
}

func (s *WarehouseService) VerifyStockTake(ctx context.Context, id, verifiedBy string) error {
	t, err := s.takeRepo.GetStockTakeByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TakeStatusVerified) {
		return domain.ErrTakeInvalidStatus
	}
	t.Status = domain.TakeStatusVerified
	t.VerifiedBy = verifiedBy
	t.VerifiedAt = s.now().Format("2006-01-02 15:04:05")
	t.UpdatedAt = s.now()
	return s.takeRepo.UpdateStockTake(ctx, t)
}

func (s *WarehouseService) PostStockTake(ctx context.Context, id string) error {
	t, err := s.takeRepo.GetStockTakeByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TakeStatusPosted) {
		return domain.ErrTakeInvalidStatus
	}
	items, err := s.takeRepo.GetTakeItems(ctx, id)
	if err != nil { return fmt.Errorf("get items: %w", err) }
	now := s.now()
	period := periodFromTime(now)

	var glLines []domain.JournalLine
	totalDebit := 0.0
	totalCredit := 0.0

	for _, item := range items {
		variance := item.ActualQty - item.ExpectedQty
		if variance == 0 {
			continue
		}
		cost := variance * item.UnitCost

		transType := domain.TransTakeVariance
		if err := s.adjustStockBalance(ctx, t.CompanyID, t.WarehouseID, item.ItemID, period,
			variance, cost, transType, "stock_take", id, ""); err != nil {
			return fmt.Errorf("adjust balance: %w", err)
		}

		if s.jeSvc != nil {
			if variance > 0 {
				glLines = append(glLines, domain.JournalLine{
					AccountCode: defaultInventoryAccount, DebitAmount: cost,
					Description: fmt.Sprintf("Take gain item %s qty +%.0f", item.ItemID, variance),
				})
				totalDebit += cost
			} else {
				loss := -cost
				glLines = append(glLines, domain.JournalLine{
					AccountCode: cogsAccount, DebitAmount: loss,
					Description: fmt.Sprintf("Take loss item %s qty %.0f", item.ItemID, -variance),
				})
				totalCredit += loss
			}
		}
	}

	if s.jeSvc != nil && len(glLines) > 0 {
		if totalDebit > 0 {
			glLines = append(glLines, domain.JournalLine{
				AccountCode: cogsAccount, CreditAmount: totalDebit,
				Description: "Take gain offset",
			})
		}
		if totalCredit > 0 {
			glLines = append(glLines, domain.JournalLine{
				AccountCode: defaultInventoryAccount, CreditAmount: totalCredit,
				Description: "Take loss offset",
			})
		}
		balancedLines := rebuildBalancedLines(glLines)
		if err := s.postGLEntry(ctx, t.CompanyID,
			fmt.Sprintf("Stock take %s variance posted", t.TakeNumber),
			now, balancedLines, domain.VoucherTypeOther, ""); err != nil {
			return fmt.Errorf("post GL: %w", err)
		}
	}

	t.Status = domain.TakeStatusPosted
	t.PostedAt = now.Format("2006-01-02 15:04:05")
	t.UpdatedAt = now
	return s.takeRepo.UpdateStockTake(ctx, t)
}

// ─── Inventory Valuation Run ──────────────────────────────────────────

func (s *WarehouseService) CreateValuationRun(ctx context.Context, v *domain.InventoryValuationRun) error {
	if err := v.Validate(); err != nil { return err }
	v.Status = domain.ValRunPending
	v.CreatedAt = s.now()
	return s.valRepo.CreateValuationRun(ctx, v)
}

func (s *WarehouseService) GetValuationRun(ctx context.Context, id string) (*domain.InventoryValuationRun, error) {
	return s.valRepo.GetValuationRunByID(ctx, id)
}

func (s *WarehouseService) ListValuationRuns(ctx context.Context, companyID string) ([]domain.InventoryValuationRun, error) {
	return s.valRepo.ListValuationRuns(ctx, companyID)
}

func (s *WarehouseService) RunValuation(ctx context.Context, id string) error {
	v, err := s.valRepo.GetValuationRunByID(ctx, id)
	if err != nil { return err }
	if v.Status != domain.ValRunPending { return fmt.Errorf("valuation run already started or completed") }
	v.Status = domain.ValRunRunning
	if err := s.valRepo.UpdateValuationRun(ctx, v); err != nil { return err }
	v.Status = domain.ValRunCompleted
	v.CompletedAt = s.now().Format("2006-01-02 15:04:05")
	return s.valRepo.UpdateValuationRun(ctx, v)
}

// ─── GRN (Goods Receipt Note) ─────────────────────────────────────────

func (s *WarehouseService) CreateGRN(ctx context.Context, g *domain.GRN) error {
	if err := g.Validate(); err != nil { return err }
	g.Status = 	domain.GRNDraft
	now := s.now()
	g.CreatedAt = now
	g.UpdatedAt = now
	if err := s.grnRepo.CreateGRN(ctx, g); err != nil {
		return err
	}
	for i := range g.Lines {
		g.Lines[i].GRNID = g.ID
	}
	if len(g.Lines) > 0 {
		if err := s.grnRepo.CreateGRNLines(ctx, g.Lines); err != nil {
			return fmt.Errorf("create grn lines: %w", err)
		}
	}
	return nil
}

func (s *WarehouseService) GetGRN(ctx context.Context, id string) (*domain.GRN, error) {
	g, err := s.grnRepo.GetGRN(ctx, id)
	if err != nil { return nil, err }
	items, err := s.grnRepo.GetGRNLines(ctx, id)
	if err != nil { return nil, err }
	g.Lines = items
	return g, nil
}

func (s *WarehouseService) ListGRNs(ctx context.Context, companyID string) ([]domain.GRN, error) {
	results, _, err := s.grnRepo.ListGRNs(ctx, domain.GRNFilter{CompanyID: companyID})
	return results, err
}

func (s *WarehouseService) PostGRN(ctx context.Context, id string) error {
	g, err := s.grnRepo.GetGRN(ctx, id)
	if err != nil { return err }
	if !g.Status.ValidTransition(domain.GRNPosted) {
		return domain.ErrGRNStatusInvalid
	}
	items, err := s.grnRepo.GetGRNLines(ctx, id)
	if err != nil { return fmt.Errorf("get items: %w", err) }
	now := s.now()
	period := periodFromTime(now)
	totalCost := 0.0

	for _, item := range items {
		receive := item.QuantityReceived - item.QuantityRejected
		if receive <= 0 {
			continue
		}
		cost := float64(receive) * item.UnitPrice
		totalCost += cost
		if err := s.adjustStockBalance(ctx, g.CompanyID, g.WarehouseID, item.ItemID, period,
			float64(receive), cost, domain.TransReceipt, "grn", id, g.CreatedBy); err != nil {
			return fmt.Errorf("adjust balance: %w", err)
		}
	}

	if s.jeSvc != nil && totalCost > 0 {
		lines := []domain.JournalLine{
			{AccountCode: defaultInventoryAccount, DebitAmount: totalCost,
				Description: fmt.Sprintf("GRN %s goods receipt", g.GRNNumber)},
			{AccountCode: pettyCashAccount, CreditAmount: totalCost,
				Description: fmt.Sprintf("AP temp for GRN %s", g.GRNNumber)},
		}
		if err := s.postGLEntry(ctx, g.CompanyID,
			fmt.Sprintf("Goods receipt GRN %s", g.GRNNumber),
			now, lines, domain.VoucherTypeInventoryReceipt, g.CreatedBy); err != nil {
			return fmt.Errorf("post GL: %w", err)
		}
	}

	g.Status = domain.GRNPosted
	g.PostedAt = now
	g.UpdatedAt = now
	return s.grnRepo.UpdateGRN(ctx, g)
}

func (s *WarehouseService) CancelGRN(ctx context.Context, id, reason string) error {
	g, err := s.grnRepo.GetGRN(ctx, id)
	if err != nil { return err }
	if !	g.Status.ValidTransition(domain.GRNCancelled) {
		return domain.ErrGRNStatusInvalid
	}
	if g.Status == domain.GRNPosted {
		items, err := s.grnRepo.GetGRNLines(ctx, id)
		if err != nil { return fmt.Errorf("get items: %w", err) }
		now := s.now()
		period := periodFromTime(now)
		totalCost := 0.0

		for _, item := range items {
			receive := item.QuantityReceived - item.QuantityRejected
			if receive <= 0 {
				continue
			}
			cost := float64(receive) * item.UnitPrice
			totalCost += cost
			if err := s.adjustStockBalance(ctx, g.CompanyID, g.WarehouseID, item.ItemID, period,
				-float64(receive), -cost, domain.TransIssue, "grn_cancel", id, ""); err != nil {
				return fmt.Errorf("reverse balance: %w", err)
			}
		}

		if s.jeSvc != nil && totalCost > 0 {
			lines := []domain.JournalLine{
				{AccountCode: pettyCashAccount, DebitAmount: totalCost,
					Description: fmt.Sprintf("Reverse GRN %s", g.GRNNumber)},
				{AccountCode: defaultInventoryAccount, CreditAmount: totalCost,
					Description: fmt.Sprintf("Reverse GRN %s inventory", g.GRNNumber)},
			}
			if err := s.postGLEntry(ctx, g.CompanyID,
				fmt.Sprintf("Cancel GRN %s: %s", g.GRNNumber, reason),
				now, lines, domain.VoucherTypeInventoryIssue, ""); err != nil {
				return fmt.Errorf("reverse GL: %w", err)
			}
		}
	}

	g.Status = domain.GRNCancelled
	g.CancelledReason = reason
	g.UpdatedAt = s.now()
	return s.grnRepo.UpdateGRN(ctx, g)
}

// ─── Stock Level Warnings ─────────────────────────────────────────────

type StockWarning struct {
	ItemID      string  `json:"item_id"`
	ItemCode    string  `json:"item_code"`
	ItemName    string  `json:"item_name"`
	WarehouseID string  `json:"warehouse_id"`
	CurrentQty  float64 `json:"current_qty"`
	MinStock    float64 `json:"min_stock"`
	MaxStock    float64 `json:"max_stock"`
	WarningType string  `json:"warning_type"` // "BELOW_MIN" or "ABOVE_MAX"
}

// CheckStockLevel checks if a single item's stock level triggers a warning.
func (s *WarehouseService) CheckStockLevel(ctx context.Context, companyID, warehouseID, itemID, period string) (*StockWarning, error) {
	item, err := s.itemRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	bal, err := s.balRepo.FindStockBalance(ctx, companyID, warehouseID, itemID, period)
	if err != nil {
		// No balance = zero stock
		if item.MinStock > 0 {
			return &StockWarning{
				ItemID: itemID, ItemCode: item.Code, ItemName: item.Name,
				WarehouseID: warehouseID, CurrentQty: 0,
				MinStock: item.MinStock, MaxStock: item.MaxStock,
				WarningType: "BELOW_MIN",
			}, nil
		}
		return nil, nil
	}
	if item.MinStock > 0 && bal.Quantity < item.MinStock {
		return &StockWarning{
			ItemID: itemID, ItemCode: item.Code, ItemName: item.Name,
			WarehouseID: warehouseID, CurrentQty: bal.Quantity,
			MinStock: item.MinStock, MaxStock: item.MaxStock,
			WarningType: "BELOW_MIN",
		}, nil
	}
	if item.MaxStock > 0 && bal.Quantity > item.MaxStock {
		return &StockWarning{
			ItemID: itemID, ItemCode: item.Code, ItemName: item.Name,
			WarehouseID: warehouseID, CurrentQty: bal.Quantity,
			MinStock: item.MinStock, MaxStock: item.MaxStock,
			WarningType: "ABOVE_MAX",
		}, nil
	}
	return nil, nil
}

// GetStockWarnings returns all items in a warehouse with stock level warnings.
func (s *WarehouseService) GetStockWarnings(ctx context.Context, companyID, warehouseID, period string) ([]StockWarning, error) {
	items, err := s.itemRepo.ListItems(ctx, companyID)
	if err != nil {
		return nil, err
	}
	var warnings []StockWarning
	for _, item := range items {
		if !item.IsActive {
			continue
		}
		w, err := s.CheckStockLevel(ctx, companyID, warehouseID, item.ID, period)
		if err != nil {
			continue
		}
		if w != nil {
			warnings = append(warnings, *w)
		}
	}
	return warnings, nil
}
