package domain

import "time"

// ─── Enums ────────────────────────────────────────────────────────────────

type ValuationMethod string
const (ValuationWeightedAvg ValuationMethod="weighted_average"; ValuationFIFO ValuationMethod="fifo"; ValuationSpecificID ValuationMethod="specific_id"; ValuationStandard ValuationMethod="standard")

type TransactionType string
const (TransReceipt TransactionType="RECEIPT"; TransIssue TransactionType="ISSUE"; TransTransferIn TransactionType="TRANSFER_IN"; TransTransferOut TransactionType="TRANSFER_OUT"; TransAdjIn TransactionType="ADJ_IN"; TransAdjOut TransactionType="ADJ_OUT"; TransTakeVariance TransactionType="TAKE_VARIANCE")

type TransferStatus string
const (TransStatusDraft TransferStatus="DRAFT"; TransStatusPending TransferStatus="PENDING"; TransStatusApproved TransferStatus="APPROVED"; TransStatusTransferred TransferStatus="TRANSFERRED"; TransStatusCompleted TransferStatus="COMPLETED"; TransStatusCancelled TransferStatus="CANCELLED")

func (s TransferStatus) ValidTransition(next TransferStatus) bool {
	switch s {
	case TransStatusDraft: return next == TransStatusPending || next == TransStatusCancelled
	case TransStatusPending: return next == TransStatusApproved || next == TransStatusCancelled
	case TransStatusApproved: return next == TransStatusTransferred
	case TransStatusTransferred: return next == TransStatusCompleted
	case TransStatusCompleted, TransStatusCancelled: return false
	default: return false
	}
}

type AdjType string
const (AdjIncrease AdjType="INCREASE"; AdjDecrease AdjType="DECREASE")

type AdjStatus string
const (AdjStatusDraft AdjStatus="DRAFT"; AdjStatusPending AdjStatus="PENDING"; AdjStatusApproved AdjStatus="APPROVED"; AdjStatusPosted AdjStatus="POSTED"; AdjStatusRejected AdjStatus="REJECTED")

func (s AdjStatus) ValidTransition(next AdjStatus) bool {
	switch s {
	case AdjStatusDraft: return next == AdjStatusPending || next == AdjStatusRejected
	case AdjStatusPending: return next == AdjStatusApproved || next == AdjStatusRejected
	case AdjStatusApproved: return next == AdjStatusPosted
	case AdjStatusPosted, AdjStatusRejected: return false
	default: return false
	}
}

type TakeStatus string
const (TakeStatusPlanning TakeStatus="PLANNING"; TakeStatusInProgress TakeStatus="IN_PROGRESS"; TakeStatusCompleted TakeStatus="COMPLETED"; TakeStatusVerified TakeStatus="VERIFIED"; TakeStatusPosted TakeStatus="POSTED")

func (s TakeStatus) ValidTransition(next TakeStatus) bool {
	switch s {
	case TakeStatusPlanning: return next == TakeStatusInProgress
	case TakeStatusInProgress: return next == TakeStatusCompleted
	case TakeStatusCompleted: return next == TakeStatusVerified
	case TakeStatusVerified: return next == TakeStatusPosted
	case TakeStatusPosted: return false
	default: return false
	}
}

type ValuationRunStatus string
const (ValRunPending ValuationRunStatus="PENDING"; ValRunRunning ValuationRunStatus="RUNNING"; ValRunCompleted ValuationRunStatus="COMPLETED"; ValRunFailed ValuationRunStatus="FAILED")

// ─── Warehouse ────────────────────────────────────────────────────────────

type Warehouse struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   string    `json:"address,omitempty"`
	Manager   string    `json:"manager,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (w *Warehouse) Validate() error {
	if w.Code == "" { return ErrWarehouseCodeRequired }
	if w.Name == "" { return ErrWarehouseNameRequired }
	return nil
}

// ─── Item Category ────────────────────────────────────────────────────────

type ItemCategory struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ParentID    string    `json:"parent_id,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (c *ItemCategory) Validate() error {
	if c.Code == "" { return ErrCategoryCodeRequired }
	if c.Name == "" { return ErrCategoryNameRequired }
	return nil
}

// ─── Item ─────────────────────────────────────────────────────────────────

type Item struct {
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	Code            string          `json:"code"`
	Name            string          `json:"name"`
	CategoryID      string          `json:"category_id,omitempty"`
	Unit            string          `json:"unit"`
	BasePrice       float64         `json:"base_price"`
	MinStock        float64         `json:"min_stock"`
	MaxStock        float64         `json:"max_stock"`
	ValuationMethod ValuationMethod `json:"valuation_method,omitempty"`
	TaxRate         float64         `json:"tax_rate"`
	IsActive        bool            `json:"is_active"`
	Notes           string          `json:"notes,omitempty"`
	CreatedAt       time.Time       `json:"created_at,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at,omitempty"`
}

func (i *Item) Validate() error {
	if i.Code == "" { return ErrItemCodeRequired }
	if i.Name == "" { return ErrItemNameRequired }
	if i.Unit == "" { return ErrItemUnitRequired }
	if i.ValuationMethod == "" { i.ValuationMethod = ValuationWeightedAvg }
	switch i.ValuationMethod {
	case ValuationWeightedAvg, ValuationFIFO, ValuationSpecificID, ValuationStandard:
	default: return ErrValuationMethodInvalid
	}
	return nil
}

// ─── Stock Balance ────────────────────────────────────────────────────────

type StockBalance struct {
	ID                 string    `json:"id"`
	CompanyID          string    `json:"company_id"`
	WarehouseID        string    `json:"warehouse_id"`
	ItemID             string    `json:"item_id"`
	Period             string    `json:"period"`
	Quantity           float64   `json:"quantity"`
	UnitCost           float64   `json:"unit_cost"`
	TotalCost          float64   `json:"total_cost"`
	LastTransactionAt  time.Time `json:"last_transaction_at,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// ─── Inventory Transaction ───────────────────────────────────────────────

type InventoryTransaction struct {
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	WarehouseID     string          `json:"warehouse_id"`
	ItemID          string          `json:"item_id"`
	TransType       TransactionType `json:"trans_type"`
	RefType         string          `json:"ref_type,omitempty"`
	RefID           string          `json:"ref_id,omitempty"`
	QtyBefore       float64         `json:"qty_before"`
	Quantity        float64         `json:"quantity"`
	QtyAfter        float64         `json:"qty_after"`
	UnitCost        float64         `json:"unit_cost"`
	TotalCost       float64         `json:"total_cost"`
	CreatedBy       string          `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at,omitempty"`
	Notes           string          `json:"notes,omitempty"`
}

// ─── Stock Transfer ──────────────────────────────────────────────────────

type StockTransfer struct {
	ID              string         `json:"id"`
	CompanyID       string         `json:"company_id"`
	TransferNumber  string         `json:"transfer_number"`
	FromWarehouseID string         `json:"from_warehouse_id"`
	ToWarehouseID   string         `json:"to_warehouse_id"`
	Status          TransferStatus `json:"status"`
	TransferDate    string         `json:"transfer_date"`
	CreatedBy       string         `json:"created_by"`
	ApprovedBy      string         `json:"approved_by,omitempty"`
	ApprovedAt      string         `json:"approved_at,omitempty"`
	CompletedBy     string         `json:"completed_by,omitempty"`
	CompletedAt     string         `json:"completed_at,omitempty"`
	CancelledReason string         `json:"cancelled_reason,omitempty"`
	Notes           string         `json:"notes,omitempty"`
	CreatedAt       time.Time      `json:"created_at,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty"`
	Items           []TransferItem `json:"items,omitempty"`
}

type TransferItem struct {
	ID        string  `json:"id"`
	TransferID string `json:"transfer_id"`
	ItemID    string  `json:"item_id"`
	Quantity  float64 `json:"quantity"`
	UnitCost  float64 `json:"unit_cost"`
}

func (t *StockTransfer) Validate() error {
	if t.TransferNumber == "" { return ErrTransferNumberRequired }
	if t.FromWarehouseID == "" { return ErrTransferFromWHRequired }
	if t.ToWarehouseID == "" { return ErrTransferToWHRequired }
	if t.FromWarehouseID == t.ToWarehouseID { return ErrTransferSameWH }
	if len(t.Items) == 0 { return ErrTransferItemsRequired }
	if t.Status == "" { t.Status = TransStatusDraft }
	return nil
}

// ─── Stock Adjustment ────────────────────────────────────────────────────

type StockAdjustment struct {
	ID                string          `json:"id"`
	CompanyID         string          `json:"company_id"`
	WarehouseID       string          `json:"warehouse_id"`
	AdjustmentNumber  string          `json:"adjustment_number"`
	AdjType           AdjType         `json:"adj_type"`
	Reason            string          `json:"reason,omitempty"`
	Status            AdjStatus       `json:"status"`
	CreatedBy         string          `json:"created_by"`
	ApprovedBy        string          `json:"approved_by,omitempty"`
	ApprovedAt        string          `json:"approved_at,omitempty"`
	PostedAt          string          `json:"posted_at,omitempty"`
	RejectedReason    string          `json:"rejected_reason,omitempty"`
	Notes             string          `json:"notes,omitempty"`
	CreatedAt         time.Time       `json:"created_at,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at,omitempty"`
	Items             []AdjItem       `json:"items,omitempty"`
}

type AdjItem struct {
	ID             string  `json:"id"`
	AdjustmentID   string  `json:"adjustment_id"`
	ItemID         string  `json:"item_id"`
	QtyBefore      float64 `json:"qty_before"`
	QtyAfter       float64 `json:"qty_after"`
	UnitCost       float64 `json:"unit_cost"`
	Reason         string  `json:"reason,omitempty"`
}

func (a *StockAdjustment) Validate() error {
	if a.AdjustmentNumber == "" { return ErrAdjNumberRequired }
	if a.WarehouseID == "" { return ErrAdjWarehouseRequired }
	if len(a.Items) == 0 { return ErrAdjItemsRequired }
	if a.AdjType == "" { return ErrAdjTypeRequired }
	switch a.AdjType {
	case AdjIncrease, AdjDecrease:
	default: return ErrAdjTypeInvalid
	}
	if a.Status == "" { a.Status = AdjStatusDraft }
	return nil
}

// ─── Stock Take ──────────────────────────────────────────────────────────

type StockTake struct {
	ID          string      `json:"id"`
	CompanyID   string      `json:"company_id"`
	WarehouseID string      `json:"warehouse_id"`
	TakeNumber  string      `json:"take_number"`
	Status      TakeStatus  `json:"status"`
	TakeDate    string      `json:"take_date"`
	CreatedBy   string      `json:"created_by"`
	VerifiedBy  string      `json:"verified_by,omitempty"`
	VerifiedAt  string      `json:"verified_at,omitempty"`
	PostedAt    string      `json:"posted_at,omitempty"`
	Notes       string      `json:"notes,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
	Items       []TakeItem  `json:"items,omitempty"`
}

type TakeItem struct {
	ID               string  `json:"id"`
	TakeID           string  `json:"take_id"`
	ItemID           string  `json:"item_id"`
	ExpectedQty      float64 `json:"expected_qty"`
	ActualQty        float64 `json:"actual_qty"`
	UnitCost         float64 `json:"unit_cost"`
	Variance         float64 `json:"variance"`
	Notes            string  `json:"notes,omitempty"`
}

func (t *StockTake) Validate() error {
	if t.TakeNumber == "" { return ErrTakeNumberRequired }
	if t.WarehouseID == "" { return ErrTakeWarehouseRequired }
	if t.TakeDate == "" { return ErrTakeDateRequired }
	if len(t.Items) == 0 { return ErrTakeItemsRequired }
	if t.Status == "" { t.Status = TakeStatusPlanning }
	return nil
}

// ─── Inventory Valuation Run ─────────────────────────────────────────────

type InventoryValuationRun struct {
	ID            string            `json:"id"`
	CompanyID     string            `json:"company_id"`
	ValuationDate string            `json:"valuation_date"`
	Method        ValuationMethod   `json:"method"`
	Status        ValuationRunStatus `json:"status"`
	CreatedBy     string            `json:"created_by"`
	CompletedAt   string            `json:"completed_at,omitempty"`
	ErrorLog      string            `json:"error_log,omitempty"`
	Notes         string            `json:"notes,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
}

func (v *InventoryValuationRun) Validate() error {
	if v.ValuationDate == "" { return ErrValRunDateRequired }
	if v.Method == "" { v.Method = ValuationWeightedAvg }
	switch v.Method {
	case ValuationWeightedAvg, ValuationFIFO, ValuationSpecificID, ValuationStandard:
	default: return ErrValuationMethodInvalid
	}
	if v.Status == "" { v.Status = ValRunPending }
	return nil
}
