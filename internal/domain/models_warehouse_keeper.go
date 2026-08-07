package domain

import "time"

// ─── Keeper Enums ───────────────────────────────────────────────────────────

type KeeperRole string
const (
	KeeperRoleKeeper   KeeperRole = "keeper"
	KeeperRoleManager  KeeperRole = "manager"
)

type LedgerVoucherType string
const (
	VoucherReceipt        LedgerVoucherType = "receipt"
	VoucherIssue          LedgerVoucherType = "issue"
	VoucherTransferIn     LedgerVoucherType = "transfer_in"
	VoucherTransferOut    LedgerVoucherType = "transfer_out"
	VoucherAdjustmentIn   LedgerVoucherType = "adjustment_in"
	VoucherAdjustmentOut  LedgerVoucherType = "adjustment_out"
	VoucherCountGain      LedgerVoucherType = "count_gain"
	VoucherCountLoss      LedgerVoucherType = "count_loss"
)

type LedgerEntryStatus string
const (
	LedgerStatusRecorded   LedgerEntryStatus = "recorded"
	LedgerStatusUnrecorded LedgerEntryStatus = "unrecorded"
)

type KeeperCountStatus string
const (
	KeeperCountPendingKeeper    KeeperCountStatus = "pending_keeper"
	KeeperCountPendingAccountant KeeperCountStatus = "pending_accountant"
	KeeperCountPendingManager   KeeperCountStatus = "pending_manager"
	KeeperCountCompleted        KeeperCountStatus = "completed"
)

// ─── Warehouse Keeper Assignment ────────────────────────────────────────────

type WarehouseKeeperAssignment struct {
	ID            string      `json:"id"`
	CompanyID     string      `json:"company_id" validate:"required"`
	WarehouseID   string      `json:"warehouse_id" validate:"required"`
	UserID        string      `json:"user_id" validate:"required"`
	Role          KeeperRole  `json:"role" validate:"required,oneof=keeper manager"`
	EffectiveFrom time.Time   `json:"effective_from" validate:"required"`
	EffectiveTo   *time.Time  `json:"effective_to,omitempty"`
	IsActive      bool        `json:"is_active"`
	CreatedBy     string      `json:"created_by"`
	CreatedAt     time.Time   `json:"created_at,omitempty"`
	UpdatedAt     time.Time   `json:"updated_at,omitempty"`
}

func (a *WarehouseKeeperAssignment) Validate() error {
	if a.CompanyID == "" { return ErrWarehouseNotFound } // TODO: proper errors
	if a.WarehouseID == "" { return ErrWarehouseNotFound }
	if a.UserID == "" { return ErrWarehouseNotFound }
	if a.Role != KeeperRoleKeeper && a.Role != KeeperRoleManager { return ErrWarehouseNotFound }
	if a.EffectiveTo != nil && a.EffectiveTo.Before(a.EffectiveFrom) { return ErrWarehouseNotFound }
	return nil
}

// ─── Stock Ledger Entry ─────────────────────────────────────────────────────

type StockLedgerEntry struct {
	ID             string            `json:"id"`
	CompanyID      string            `json:"company_id" validate:"required"`
	WarehouseID    string            `json:"warehouse_id" validate:"required"`
	ItemID         string            `json:"item_id" validate:"required"`
	EntryDate      time.Time         `json:"entry_date" validate:"required"`
	VoucherType    LedgerVoucherType `json:"voucher_type" validate:"required"`
	VoucherNo      string            `json:"voucher_no,omitempty"`
	VoucherRefID   string            `json:"voucher_ref_id,omitempty"`
	Description    string            `json:"description,omitempty"`
	ReceiptQty     float64           `json:"receipt_qty" validate:"gte=0"`
	IssueQty       float64           `json:"issue_qty" validate:"gte=0"`
	BalanceQty     float64           `json:"balance_qty"`
	UnitCost       float64           `json:"unit_cost,omitempty"`
	TotalValue     float64           `json:"total_value,omitempty"`
	RecordedBy     string            `json:"recorded_by"`
	RecordedAt     time.Time         `json:"recorded_at,omitempty"`
	UnrecordedBy   string            `json:"unrecorded_by,omitempty"`
	UnrecordedAt   *time.Time        `json:"unrecorded_at,omitempty"`
	UnrecordReason string            `json:"unrecord_reason,omitempty"`
	Status         LedgerEntryStatus `json:"status"`
	CreatedAt      time.Time         `json:"created_at,omitempty"`
}

// ─── Keeper Inventory Count ─────────────────────────────────────────────────

type KeeperInventoryCount struct {
	ID                     string            `json:"id"`
	StockTakeID            string            `json:"stock_take_id" validate:"required"`
	KeeperID               string            `json:"keeper_id" validate:"required"`
	BookQtyHidden          bool              `json:"book_qty_hidden"`
	KeeperCountDate        *time.Time        `json:"keeper_count_date,omitempty"`
	KeeperSignature        string            `json:"keeper_signature,omitempty"`
	AccountantSignature    string            `json:"accountant_signature,omitempty"`
	ManagerSignature       string            `json:"manager_signature,omitempty"`
	Status                 KeeperCountStatus `json:"status"`
	CreatedAt              time.Time         `json:"created_at,omitempty"`
	UpdatedAt              time.Time         `json:"updated_at,omitempty"`
}

// ─── Warehouse Keeper Config ────────────────────────────────────────────────

type WarehouseKeeperConfig struct {
	ID                      string    `json:"id"`
	CompanyID               string    `json:"company_id" validate:"required"`
	ModuleEnabled           bool      `json:"module_enabled"`
	CostPriceHiddenFromKeeper bool    `json:"cost_price_hidden_from_keeper"`
	AutoRecordOnGRN         bool      `json:"auto_record_on_grn"`
	CreatedAt               time.Time `json:"created_at,omitempty"`
	UpdatedAt               time.Time `json:"updated_at,omitempty"`
}

// ─── Virtual Entities (not persisted) ───────────────────────────────────────

type KeeperReconciliationItem struct {
	ItemID        string  `json:"item_id"`
	ItemCode      string  `json:"item_code"`
	ItemName      string  `json:"item_name"`
	WarehouseID   string  `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	KeeperQty     float64 `json:"keeper_qty"`
	AccountingQty float64 `json:"accounting_qty"`
	VarianceQty   float64 `json:"variance_qty"`
	UnitCost      float64 `json:"unit_cost"`
	VarianceValue float64 `json:"variance_value"`
	LastUpdated   time.Time `json:"last_updated"`
}

type StockCardLine struct {
	Date        time.Time `json:"date"`
	VoucherNo   string    `json:"voucher_no"`
	Description string    `json:"description"`
	ReceiptQty  float64   `json:"receipt_qty"`
	IssueQty    float64   `json:"issue_qty"`
	BalanceQty  float64   `json:"balance_qty"`
}

type StockCard struct {
	WarehouseID   string         `json:"warehouse_id"`
	WarehouseName string         `json:"warehouse_name"`
	ItemID        string         `json:"item_id"`
	ItemCode      string         `json:"item_code"`
	ItemName      string         `json:"item_name"`
	Period        string         `json:"period"`
	OpeningBalance float64       `json:"opening_balance"`
	Lines         []StockCardLine `json:"lines"`
	ClosingBalance float64       `json:"closing_balance"`
}

type KeeperInventorySummaryItem struct {
	ItemID        string  `json:"item_id"`
	ItemCode      string  `json:"item_code"`
	ItemName      string  `json:"item_name"`
	Unit          string  `json:"unit"`
	Quantity      float64 `json:"quantity"`
	UnitCost      float64 `json:"unit_cost"`
	TotalValue    float64 `json:"total_value"`
	LastUpdated   time.Time `json:"last_updated"`
}

type LedgerFilter struct {
	CompanyID   string
	WarehouseID string
	ItemID      string   // optional
	From        time.Time
	To          time.Time
	VoucherType string   // optional
	Status      string   // optional (recorded/unrecorded)
	Page        int
	PageSize    int
}
