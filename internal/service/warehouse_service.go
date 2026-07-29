package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

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
) *WarehouseService {
	return &WarehouseService{
		whRepo: whRepo, catRepo: catRepo, itemRepo: itemRepo,
		balRepo: balRepo, txnRepo: txnRepo, trfRepo: trfRepo,
		adjRepo: adjRepo, takeRepo: takeRepo, valRepo: valRepo,
		now: time.Now,
	}
}

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
	return s.trfRepo.GetStockTransferByID(ctx, id)
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
	t.Status = domain.TransStatusTransferred
	t.UpdatedAt = s.now()
	return s.trfRepo.UpdateStockTransfer(ctx, t)
}

func (s *WarehouseService) CompleteStockTransfer(ctx context.Context, id, completedBy string) error {
	t, err := s.trfRepo.GetStockTransferByID(ctx, id)
	if err != nil { return err }
	if !t.Status.ValidTransition(domain.TransStatusCompleted) {
		return domain.ErrTransferInvalidStatus
	}
	t.Status = domain.TransStatusCompleted
	t.CompletedBy = completedBy
	t.CompletedAt = s.now().Format("2006-01-02 15:04:05")
	t.UpdatedAt = s.now()
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
	return s.adjRepo.GetStockAdjustmentByID(ctx, id)
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
	a.Status = domain.AdjStatusPosted
	a.PostedAt = s.now().Format("2006-01-02 15:04:05")
	a.UpdatedAt = s.now()
	return s.adjRepo.UpdateStockAdjustment(ctx, a)
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
	return s.takeRepo.GetStockTakeByID(ctx, id)
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
	t.Status = domain.TakeStatusPosted
	t.PostedAt = s.now().Format("2006-01-02 15:04:05")
	t.UpdatedAt = s.now()
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
