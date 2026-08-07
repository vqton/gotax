package domain

import "time"

type WarehouseKeeperAssignmentGORM struct {
	ID            string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID     string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	WarehouseID   string     `gorm:"column:warehouse_id;not null;size:36;index" json:"warehouseId"`
	UserID        string     `gorm:"column:user_id;not null;size:36;index" json:"userId"`
	Role          string     `gorm:"column:role;not null;size:20;default:'keeper'" json:"role"`
	EffectiveFrom time.Time  `gorm:"column:effective_from;not null" json:"effectiveFrom"`
	EffectiveTo   *time.Time `gorm:"column:effective_to" json:"effectiveTo"`
	IsActive      bool       `gorm:"column:is_active;not null;default:true" json:"isActive"`
	CreatedBy     string     `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (WarehouseKeeperAssignmentGORM) TableName() string { return "warehouse_keeper_assignments" }

type StockLedgerEntryGORM struct {
	ID             string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID      string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	WarehouseID    string     `gorm:"column:warehouse_id;not null;size:36;index" json:"warehouseId"`
	ItemID         string     `gorm:"column:item_id;not null;size:36;index" json:"itemId"`
	EntryDate      time.Time  `gorm:"column:entry_date;not null" json:"entryDate"`
	VoucherType    string     `gorm:"column:voucher_type;not null;size:30" json:"voucherType"`
	VoucherNo      *string    `gorm:"column:voucher_no;size:50" json:"voucherNo"`
	VoucherRefID   *string    `gorm:"column:voucher_ref_id;size:36" json:"voucherRefId"`
	Description    *string    `gorm:"column:description;type:text" json:"description"`
	ReceiptQty     float64    `gorm:"column:receipt_qty;not null;default:0" json:"receiptQty"`
	IssueQty       float64    `gorm:"column:issue_qty;not null;default:0" json:"issueQty"`
	BalanceQty     float64    `gorm:"column:balance_qty;not null;default:0" json:"balanceQty"`
	UnitCost       *float64   `gorm:"column:unit_cost" json:"unitCost"`
	TotalValue     *float64   `gorm:"column:total_value" json:"totalValue"`
	RecordedBy     string     `gorm:"column:recorded_by;not null;size:36" json:"recordedBy"`
	RecordedAt     time.Time  `gorm:"column:recorded_at;autoCreateTime" json:"recordedAt"`
	UnrecordedBy   *string    `gorm:"column:unrecorded_by;size:36" json:"unrecordedBy"`
	UnrecordedAt   *time.Time `gorm:"column:unrecorded_at" json:"unrecordedAt"`
	UnrecordReason *string    `gorm:"column:unrecord_reason;type:text" json:"unrecordReason"`
	Status         string     `gorm:"column:status;not null;size:20;default:'recorded'" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (StockLedgerEntryGORM) TableName() string { return "stock_ledger_entries" }

type KeeperInventoryCountGORM struct {
	ID                  string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	StockTakeID         string     `gorm:"column:stock_take_id;not null;size:36;index" json:"stockTakeId"`
	KeeperID            string     `gorm:"column:keeper_id;not null;size:36;index" json:"keeperId"`
	BookQtyHidden       bool       `gorm:"column:book_qty_hidden;not null;default:false" json:"bookQtyHidden"`
	KeeperCountDate     *time.Time `gorm:"column:keeper_count_date" json:"keeperCountDate"`
	KeeperSignature     *string    `gorm:"column:keeper_signature;size:255" json:"keeperSignature"`
	AccountantSignature *string    `gorm:"column:accountant_signature;size:255" json:"accountantSignature"`
	ManagerSignature    *string    `gorm:"column:manager_signature;size:255" json:"managerSignature"`
	Status              string     `gorm:"column:status;not null;size:30;default:'pending_keeper'" json:"status"`
	CreatedAt           time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (KeeperInventoryCountGORM) TableName() string { return "keeper_inventory_counts" }

type WarehouseKeeperConfigGORM struct {
	ID                        string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID                 string    `gorm:"column:company_id;not null;size:36;uniqueIndex" json:"companyId"`
	ModuleEnabled             bool      `gorm:"column:module_enabled;not null;default:true" json:"moduleEnabled"`
	CostPriceHiddenFromKeeper bool      `gorm:"column:cost_price_hidden_from_keeper;not null;default:false" json:"costPriceHiddenFromKeeper"`
	AutoRecordOnGRN           bool      `gorm:"column:auto_record_on_grn;not null;default:false" json:"autoRecordOnGrn"`
	CreatedAt                 time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt                 time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (WarehouseKeeperConfigGORM) TableName() string { return "warehouse_keeper_config" }
