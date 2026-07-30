package domain

import "time"

type StockTransferGORM struct {
	ID             string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID      string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	FromWarehouseID string  `gorm:"column:from_warehouse_id;not null;size:36;index" json:"fromWarehouseId"`
	ToWarehouseID   string  `gorm:"column:to_warehouse_id;not null;size:36;index" json:"toWarehouseId"`
	Status         string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	ShippedAt      *time.Time `gorm:"column:shipped_at" json:"shippedAt"`
	ReceivedAt     *time.Time `gorm:"column:received_at" json:"receivedAt"`
	CreatedBy      string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Items          []TransferItemGORM `gorm:"foreignKey:TransferID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (StockTransferGORM) TableName() string { return "stock_transfers" }

type TransferItemGORM struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TransferID  string    `gorm:"column:transfer_id;not null;size:36;index:idx_trfitem_trf,unique" json:"transferId"`
	LineNumber  int       `gorm:"column:line_number;not null;index:idx_trfitem_trf,unique" json:"lineNumber"`
	ItemID      string    `gorm:"column:item_id;not null;size:36" json:"itemId"`
	Qty         float64   `gorm:"column:qty;not null" json:"qty"`
	ReceivedQty *float64  `gorm:"column:received_qty" json:"receivedQty"`
	Note        *string   `gorm:"column:note;type:text" json:"note"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (TransferItemGORM) TableName() string { return "transfer_items" }

type StockAdjustmentGORM struct {
	ID             string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID      string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	WarehouseID    string    `gorm:"column:warehouse_id;not null;size:36;index" json:"warehouseId"`
	AdjustmentDate time.Time `gorm:"column:adjustment_date;not null;type:date;index" json:"adjustmentDate"`
	Reason         string    `gorm:"column:reason;not null;size:100" json:"reason"`
	Status         string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	ApprovedAt     *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	ApprovedBy     *string   `gorm:"column:approved_by;size:36" json:"approvedBy"`
	CreatedBy      string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Items          []AdjItemGORM `gorm:"foreignKey:AdjustmentID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (StockAdjustmentGORM) TableName() string { return "stock_adjustments" }

type AdjItemGORM struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AdjustmentID  string    `gorm:"column:adjustment_id;not null;size:36;index:idx_adjitem_adj,unique" json:"adjustmentId"`
	LineNumber    int       `gorm:"column:line_number;not null;index:idx_adjitem_adj,unique" json:"lineNumber"`
	ItemID        string    `gorm:"column:item_id;not null;size:36" json:"itemId"`
	QtyOnHand     float64   `gorm:"column:qty_on_hand;not null" json:"qtyOnHand"`
	QtyCounted    float64   `gorm:"column:qty_counted;not null" json:"qtyCounted"`
	QtyAdjustment float64   `gorm:"column:qty_adjustment;not null" json:"qtyAdjustment"`
	Reason        *string   `gorm:"column:reason;type:text" json:"reason"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (AdjItemGORM) TableName() string { return "adj_items" }

type StockTakeGORM struct {
	ID             string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID      string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	WarehouseID    string    `gorm:"column:warehouse_id;not null;size:36;index" json:"warehouseId"`
	StockTakeDate  time.Time `gorm:"column:stock_take_date;not null;type:date;index" json:"stockTakeDate"`
	Status         string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	PostedAt       *time.Time `gorm:"column:posted_at" json:"postedAt"`
	PostedBy       *string   `gorm:"column:posted_by;size:36" json:"postedBy"`
	CreatedBy      string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Items          []TakeItemGORM `gorm:"foreignKey:StockTakeID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (StockTakeGORM) TableName() string { return "stock_takes" }

type TakeItemGORM struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StockTakeID  string    `gorm:"column:stock_take_id;not null;size:36;index:idx_takeitem_st,unique" json:"stockTakeId"`
	LineNumber   int       `gorm:"column:line_number;not null;index:idx_takeitem_st,unique" json:"lineNumber"`
	ItemID       string    `gorm:"column:item_id;not null;size:36" json:"itemId"`
	QtyBook      float64   `gorm:"column:qty_book;not null" json:"qtyBook"`
	QtyCounted   float64   `gorm:"column:qty_counted;not null" json:"qtyCounted"`
	QtyVariance  float64   `gorm:"column:qty_variance;not null" json:"qtyVariance"`
	Note         *string   `gorm:"column:note;type:text" json:"note"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (TakeItemGORM) TableName() string { return "take_items" }

type InventoryValuationRunGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	RunDate     time.Time `gorm:"column:run_date;not null;type:date;index" json:"runDate"`
	Period      string    `gorm:"column:period;not null;size:10" json:"period"`
	Method      string    `gorm:"column:method;not null;size:20" json:"method"`
	TotalValue  float64   `gorm:"column:total_value;not null" json:"totalValue"`
	Status      string    `gorm:"column:status;not null;size:20;default:'COMPLETED';index" json:"status"`
	ExecutedBy  string    `gorm:"column:executed_by;not null;size:36" json:"executedBy"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (InventoryValuationRunGORM) TableName() string { return "inventory_valuation_runs" }
