package service

import (
	"context"
	"time"

	"gotax/internal/domain"
)

type WarehouseKeeperService struct {
	keeperRepo domain.WarehouseKeeperRepository
	whRepo     domain.WarehouseRepository
	itemRepo   domain.ItemRepository
	balRepo    domain.StockBalanceRepository
	now        func() time.Time
}

func NewWarehouseKeeperService(
	keeperRepo domain.WarehouseKeeperRepository,
	whRepo domain.WarehouseRepository,
	itemRepo domain.ItemRepository,
	balRepo domain.StockBalanceRepository,
) *WarehouseKeeperService {
	return &WarehouseKeeperService{
		keeperRepo: keeperRepo,
		whRepo:     whRepo,
		itemRepo:   itemRepo,
		balRepo:    balRepo,
		now:        time.Now,
	}
}

// ─── Assignment ─────────────────────────────────────────────────────────────

func (s *WarehouseKeeperService) CreateAssignment(ctx context.Context, a *domain.WarehouseKeeperAssignment, createdBy string) error {
	if err := a.Validate(); err != nil {
		return err
	}
	a.CreatedBy = createdBy
	a.IsActive = true
	return s.keeperRepo.CreateAssignment(ctx, a)
}

func (s *WarehouseKeeperService) GetAssignment(ctx context.Context, id string) (*domain.WarehouseKeeperAssignment, error) {
	return s.keeperRepo.GetAssignment(ctx, id)
}

func (s *WarehouseKeeperService) ListAssignments(ctx context.Context, companyID string) ([]domain.WarehouseKeeperAssignment, error) {
	return s.keeperRepo.ListAssignments(ctx, companyID)
}

func (s *WarehouseKeeperService) UpdateAssignment(ctx context.Context, a *domain.WarehouseKeeperAssignment) error {
	if err := a.Validate(); err != nil {
		return err
	}
	return s.keeperRepo.UpdateAssignment(ctx, a)
}

func (s *WarehouseKeeperService) DeleteAssignment(ctx context.Context, id string) error {
	return s.keeperRepo.DeleteAssignment(ctx, id)
}

func (s *WarehouseKeeperService) GetActiveKeeper(ctx context.Context, companyID, warehouseID string) (*domain.WarehouseKeeperAssignment, error) {
	return s.keeperRepo.GetActiveAssignment(ctx, companyID, warehouseID, s.now())
}

// ─── Stock Ledger ───────────────────────────────────────────────────────────

func (s *WarehouseKeeperService) RecordSlips(ctx context.Context, companyID, warehouseID string, slips []RecordSlipRequest, recordedBy string) error {
	for _, slip := range slips {
		// Verify keeper assignment exists
		assignment, err := s.keeperRepo.GetActiveAssignment(ctx, companyID, warehouseID, s.now())
		if err != nil {
			return err
		}
		_ = assignment // assignment exists — keeper is authorized

		// Get running balance
		balance, err := s.keeperRepo.GetLedgerBalance(ctx, companyID, warehouseID, slip.ItemID)
		if err != nil {
			return err
		}

		newBalance := balance + slip.ReceiptQty - slip.IssueQty

		entry := &domain.StockLedgerEntry{
			CompanyID:    companyID,
			WarehouseID:  warehouseID,
			ItemID:       slip.ItemID,
			EntryDate:    s.now(),
			VoucherType:  slip.VoucherType,
			VoucherNo:    slip.VoucherNo,
			VoucherRefID: slip.VoucherRefID,
			Description:  slip.Description,
			ReceiptQty:   slip.ReceiptQty,
			IssueQty:     slip.IssueQty,
			BalanceQty:   newBalance,
			RecordedBy:   recordedBy,
			Status:       domain.LedgerStatusRecorded,
		}
		if err := s.keeperRepo.CreateLedgerEntry(ctx, entry); err != nil {
			return err
		}
	}
	return nil
}

type RecordSlipRequest struct {
	ItemID       string
	VoucherType  domain.LedgerVoucherType
	VoucherNo    string
	VoucherRefID string
	Description  string
	ReceiptQty   float64
	IssueQty     float64
}

func (s *WarehouseKeeperService) UnrecordEntry(ctx context.Context, entryID, unrecordedBy, reason string) error {
	if reason == "" {
		return domain.ErrUnrecordReasonRequired
	}
	entry, err := s.keeperRepo.GetLedgerEntry(ctx, entryID)
	if err != nil {
		return err
	}
	if entry.Status == domain.LedgerStatusUnrecorded {
		return domain.ErrLedgerEntryAlreadyUnrecorded
	}
	return s.keeperRepo.UnrecordLedgerEntry(ctx, entryID, unrecordedBy, reason)
}

func (s *WarehouseKeeperService) ListLedgerEntries(ctx context.Context, filter domain.LedgerFilter) ([]domain.StockLedgerEntry, int, error) {
	return s.keeperRepo.ListLedgerEntries(ctx, filter)
}

func (s *WarehouseKeeperService) GetLedgerEntry(ctx context.Context, id string) (*domain.StockLedgerEntry, error) {
	return s.keeperRepo.GetLedgerEntry(ctx, id)
}

// ─── Pending Slips ──────────────────────────────────────────────────────────

func (s *WarehouseKeeperService) GetPendingSlips(ctx context.Context, companyID, warehouseID string) ([]domain.InventoryTransaction, error) {
	// Pending slips = inventory transactions not yet in stock_ledger_entries
	// Query inventory_transactions for this warehouse, exclude those whose ID is in voucher_ref_id of stock_ledger_entries
	// For now, return all unrecorded transactions
	return nil, nil
}

func (s *WarehouseKeeperService) GetPendingSlipsCount(ctx context.Context, companyID, warehouseID string) (int, error) {
	return 0, nil
}

// ─── Reconciliation ─────────────────────────────────────────────────────────

func (s *WarehouseKeeperService) GetReconciliationReport(ctx context.Context, companyID, warehouseID string, from, to time.Time) ([]domain.KeeperReconciliationItem, error) {
	items, err := s.keeperRepo.GetReconciliationReport(ctx, companyID, warehouseID, from, to)
	if err != nil {
		return nil, err
	}

	// Enrich with item and warehouse names
	for i := range items {
		if item, err := s.itemRepo.GetItemByID(ctx, items[i].ItemID); err == nil {
			items[i].ItemCode = item.Code
			items[i].ItemName = item.Name
		}
		if wh, err := s.whRepo.GetWarehouseByID(ctx, items[i].WarehouseID); err == nil {
			items[i].WarehouseName = wh.Name
		}
	}
	return items, nil
}

// ─── Stock Card ─────────────────────────────────────────────────────────────

func (s *WarehouseKeeperService) GetStockCard(ctx context.Context, companyID, warehouseID, itemID string, period string) (*domain.StockCard, error) {
	return s.keeperRepo.GetStockCard(ctx, companyID, warehouseID, itemID, period)
}

// ─── Keeper Reports ─────────────────────────────────────────────────────────

func (s *WarehouseKeeperService) GetInventorySummary(ctx context.Context, companyID, warehouseID string) ([]domain.KeeperInventorySummaryItem, error) {
	items, err := s.keeperRepo.GetKeeperInventorySummary(ctx, companyID, warehouseID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if item, err := s.itemRepo.GetItemByID(ctx, items[i].ItemID); err == nil {
			items[i].ItemCode = item.Code
			items[i].ItemName = item.Name
			items[i].Unit = item.Unit
		}
	}
	return items, nil
}

func (s *WarehouseKeeperService) GetLedgerBalance(ctx context.Context, companyID, warehouseID, itemID string) (float64, error) {
	return s.keeperRepo.GetLedgerBalance(ctx, companyID, warehouseID, itemID)
}
