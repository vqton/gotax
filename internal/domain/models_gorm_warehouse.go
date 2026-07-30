package domain

import "time"

type WarehouseGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index:idx_warehouse_company_code,unique" json:"companyId"`
	Code        string    `gorm:"column:code;not null;size:20;index:idx_warehouse_company_code,unique" json:"code"`
	Name        string    `gorm:"column:name;not null;size:255" json:"name"`
	Location    *string   `gorm:"column:location;type:text" json:"location"`
	ManagerID   *string   `gorm:"column:manager_id;size:36" json:"managerId"`
	IsActive    bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (WarehouseGORM) TableName() string { return "warehouses" }

type ItemCategoryGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index:idx_itemcat_company_code,unique" json:"companyId"`
	Code        string    `gorm:"column:code;not null;size:20;index:idx_itemcat_company_code,unique" json:"code"`
	Name        string    `gorm:"column:name;not null;size:255" json:"name"`
	ParentID    *string   `gorm:"column:parent_id;size:36;index" json:"parentId"`
	IsActive    bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ItemCategoryGORM) TableName() string { return "item_categories" }

type ItemGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID    string    `gorm:"column:company_id;not null;size:36;index:idx_item_company_code,unique" json:"companyId"`
	Code         string    `gorm:"column:code;not null;size:30;index:idx_item_company_code,unique" json:"code"`
	Name         string    `gorm:"column:name;not null;size:255" json:"name"`
	CategoryID   *string   `gorm:"column:category_id;size:36;index" json:"categoryId"`
	Unit         *string   `gorm:"column:unit;size:20" json:"unit"`
	CostMethod   string    `gorm:"column:cost_method;not null;size:20;default:'WEIGHTED_AVG'" json:"costMethod"`
	StandardCost *float64  `gorm:"column:standard_cost" json:"standardCost"`
	IsActive     bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ItemGORM) TableName() string { return "items" }

type StockBalanceGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index:idx_stock_balance,unique" json:"companyId"`
	WarehouseID string    `gorm:"column:warehouse_id;not null;size:36;index:idx_stock_balance,unique" json:"warehouseId"`
	ItemID      string    `gorm:"column:item_id;not null;size:36;index:idx_stock_balance,unique" json:"itemId"`
	Period      string    `gorm:"column:period;not null;size:10;index:idx_stock_balance,unique" json:"period"`
	QtyOnHand   float64   `gorm:"column:qty_on_hand;not null;default:0" json:"qtyOnHand"`
	QtyReserved float64   `gorm:"column:qty_reserved;not null;default:0" json:"qtyReserved"`
	QtyAvailable float64  `gorm:"column:qty_available;not null;default:0" json:"qtyAvailable"`
	UnitCost    float64   `gorm:"column:unit_cost;not null;default:0" json:"unitCost"`
	TotalValue  float64   `gorm:"column:total_value;not null;default:0" json:"totalValue"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (StockBalanceGORM) TableName() string { return "stock_balances" }

type InventoryTransactionGORM struct {
	ID            string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID     string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	WarehouseID   string    `gorm:"column:warehouse_id;not null;size:36;index" json:"warehouseId"`
	ItemID        string    `gorm:"column:item_id;not null;size:36;index:idx_invtxn_item_period,unique" json:"itemId"`
	Period        string    `gorm:"column:period;not null;size:10;index:idx_invtxn_item_period,unique" json:"period"`
	MovementType  string    `gorm:"column:movement_type;not null;size:30;index" json:"movementType"`
	Qty           float64   `gorm:"column:qty;not null" json:"qty"`
	UnitCost      float64   `gorm:"column:unit_cost;not null" json:"unitCost"`
	TotalValue    float64   `gorm:"column:total_value;not null" json:"totalValue"`
	RefDocID      *string   `gorm:"column:ref_doc_id;size:36" json:"refDocId"`
	RefDocType    *string   `gorm:"column:ref_doc_type;size:30" json:"refDocType"`
	Note          *string   `gorm:"column:note;type:text" json:"note"`
	CreatedBy     string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (InventoryTransactionGORM) TableName() string { return "inventory_transactions" }
