package domain

import "time"

type SupplierGORM struct {
	ID                string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string    `gorm:"column:company_id;not null;size:36;index:idx_supp_company_code,unique" json:"companyId"`
	Code              string    `gorm:"column:supplier_code;not null;size:30;index:idx_supp_company_code,unique" json:"code"`
	Name              string    `gorm:"column:supplier_name;not null;size:200" json:"name"`
	TaxCode           *string   `gorm:"column:tax_code;size:20" json:"taxCode"`
	Address           *string   `gorm:"column:address;type:text" json:"address"`
	Phone             *string   `gorm:"column:phone;size:30" json:"phone"`
	Email             *string   `gorm:"column:email;size:100" json:"email"`
	ContactPerson     *string   `gorm:"column:contact_person;size:100" json:"contactPerson"`
	BankAccountName   *string   `gorm:"column:bank_account_name;size:100" json:"bankAccountName"`
	BankAccountNumber *string   `gorm:"column:bank_account_number;size:50" json:"bankAccountNumber"`
	BankName          *string   `gorm:"column:bank_name;size:100" json:"bankName"`
	PaymentTerms      *string   `gorm:"column:payment_terms;size:100" json:"paymentTerms"`
	CreditLimit       *float64  `gorm:"column:credit_limit" json:"creditLimit"`
	Currency          string    `gorm:"column:currency;size:3;default:VND" json:"currency"`
	IsActive          bool      `gorm:"column:is_active;default:true" json:"isActive"`
	Notes             *string   `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy         string    `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (SupplierGORM) TableName() string { return "suppliers" }

type PurchaseOrderGORM struct {
	ID              string              `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string              `gorm:"column:company_id;not null;size:36;index:idx_po_company_number,unique" json:"companyId"`
	PONumber        string              `gorm:"column:po_number;not null;size:30;index:idx_po_company_number,unique" json:"poNumber"`
	OrderDate       time.Time           `gorm:"column:po_date;not null;type:date;index" json:"orderDate"`
	SupplierID      string              `gorm:"column:supplier_id;not null;size:36;index" json:"supplierId"`
	DepartmentID    *string             `gorm:"column:department_id;size:36" json:"departmentId"`
	Requester       *string             `gorm:"column:requester;size:100" json:"requester"`
	Subtotal        float64             `gorm:"column:subtotal;not null" json:"subtotal"`
	TaxAmount       float64             `gorm:"column:tax_amount;default:0" json:"taxAmount"`
	TotalAmount     float64             `gorm:"column:grand_total;not null" json:"totalAmount"`
	Currency        string              `gorm:"column:currency;size:3;default:VND" json:"currency"`
	ExchangeRate    float64             `gorm:"column:exchange_rate;default:1" json:"exchangeRate"`
	LocalTotal      float64             `gorm:"column:local_total;default:0" json:"localTotal"`
	Status          string              `gorm:"column:status;size:20;default:'DRAFT';index" json:"status"`
	ApprovedBy      *string             `gorm:"column:approved_by;size:36" json:"approvedBy"`
	ApprovedAt      *time.Time          `gorm:"column:approved_at" json:"approvedAt"`
	CancelledBy     *string             `gorm:"column:cancelled_by;size:36" json:"cancelledBy"`
	CancelledAt     *time.Time          `gorm:"column:cancelled_at" json:"cancelledAt"`
	CancelReason    *string             `gorm:"column:cancel_reason;type:text" json:"cancelReason"`
	PaymentTerms    *string             `gorm:"column:payment_terms;size:100" json:"paymentTerms"`
	DeliveryDate    *time.Time          `gorm:"column:delivery_date;type:date" json:"deliveryDate"`
	DeliveryAddress *string             `gorm:"column:delivery_address;type:text" json:"deliveryAddress"`
	Notes           *string             `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy       string              `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt       time.Time           `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time           `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines           []POItemGORM        `gorm:"foreignKey:POID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (PurchaseOrderGORM) TableName() string { return "purchase_orders" }

type POItemGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	POID        string    `gorm:"column:po_id;not null;size:36;index" json:"poId"`
	LineNumber  int       `gorm:"column:line_number;not null" json:"lineNumber"`
	ItemCode    string    `gorm:"column:item_code;size:50" json:"itemCode"`
	ItemID      *string   `gorm:"column:item_id;size:36;index" json:"itemId"`
	ItemName    string    `gorm:"column:item_name;not null;size:300" json:"itemName"`
	Quantity    float64   `gorm:"column:quantity;not null" json:"quantity"`
	Unit        string    `gorm:"column:unit;size:20" json:"unit"`
	UnitPrice   float64   `gorm:"column:unit_price;not null" json:"unitPrice"`
	Discount    *float64  `gorm:"column:discount" json:"discount"`
	TaxRate     *float64  `gorm:"column:tax_rate;default:0" json:"taxRate"`
	TaxAmount   float64   `gorm:"column:tax_amount;default:0" json:"taxAmount"`
	LineTotal   float64   `gorm:"column:line_total;not null" json:"lineTotal"`
	LocalTotal  float64   `gorm:"column:local_total;default:0" json:"localTotal"`
	WarehouseID *string   `gorm:"column:warehouse_id;size:36" json:"warehouseId"`
	AccountID   string    `gorm:"column:account_id;size:36" json:"accountId"`
	Description *string   `gorm:"column:description;type:text" json:"description"`
	ReceivedQty float64   `gorm:"column:received_qty;default:0" json:"receivedQty"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (POItemGORM) TableName() string { return "po_items" }

type GRNGORM struct {
	ID          string          `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string          `gorm:"column:company_id;not null;size:36;index:idx_grn_company_number,unique" json:"companyId"`
	GRNNumber   string          `gorm:"column:grn_number;not null;size:30;index:idx_grn_company_number,unique" json:"grnNumber"`
	GRNDate     time.Time       `gorm:"column:grn_date;not null;type:date;index" json:"grnDate"`
	POID        *string         `gorm:"column:po_id;size:36;index" json:"poId"`
	SupplierID  string          `gorm:"column:supplier_id;size:36;index" json:"supplierId"`
	WarehouseID *string         `gorm:"column:warehouse_id;size:36" json:"warehouseId"`
	ReferenceNo *string         `gorm:"column:reference_no;size:100" json:"referenceNo"`
	Status      string          `gorm:"column:status;size:20;default:'DRAFT';index" json:"status"`
	GLPosted    bool            `gorm:"column:gl_posted;default:false" json:"glPosted"`
	Notes       *string         `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy   string          `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines       []GRNItemGORM   `gorm:"foreignKey:GRNID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (GRNGORM) TableName() string { return "grn_receipts" }

type GRNItemGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	GRNID       string    `gorm:"column:grn_id;not null;size:36;index" json:"grnId"`
	POLineID    *string   `gorm:"column:po_line_id;size:36" json:"poLineId"`
	LineNumber  int       `gorm:"column:line_number;not null" json:"lineNumber"`
	ItemCode    string    `gorm:"column:item_code;size:50" json:"itemCode"`
	ItemID      *string   `gorm:"column:item_id;size:36;index" json:"itemId"`
	ItemName    string    `gorm:"column:item_name;not null;size:300" json:"itemName"`
	QtyShipped  float64   `gorm:"column:qty_shipped;not null" json:"qtyShipped"`
	Quantity    float64   `gorm:"column:quantity;not null" json:"quantity"`
	Unit        string    `gorm:"column:unit;size:20" json:"unit"`
	UnitPrice   float64   `gorm:"column:unit_price;not null" json:"unitPrice"`
	TaxRate     float64   `gorm:"column:tax_rate;default:0" json:"taxRate"`
	TaxAmount   float64   `gorm:"column:tax_amount;default:0" json:"taxAmount"`
	LineTotal   float64   `gorm:"column:line_total;not null" json:"lineTotal"`
	LocalTotal  float64   `gorm:"column:local_total;default:0" json:"localTotal"`
	WarehouseID *string   `gorm:"column:warehouse_id;size:36" json:"warehouseId"`
	AccountID   string    `gorm:"column:account_id;size:36" json:"accountId"`
	Description *string   `gorm:"column:description;type:text" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (GRNItemGORM) TableName() string { return "grn_items" }

type SupplierInvoiceGORM struct {
	ID             string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID      string    `gorm:"column:company_id;not null;size:36;index:idx_sinv_company_number,unique" json:"companyId"`
	InvoiceNumber  string    `gorm:"column:invoice_number;not null;size:30;index:idx_sinv_company_number,unique" json:"invoiceNumber"`
	InvoiceDate    time.Time `gorm:"column:invoice_date;not null;type:date;index" json:"invoiceDate"`
	SupplierID     string    `gorm:"column:supplier_id;not null;size:36;index" json:"supplierId"`
	ReferenceNo    *string   `gorm:"column:reference_no;size:100" json:"referenceNo"`
	Subtotal       float64   `gorm:"column:subtotal;not null" json:"subtotal"`
	TaxAmount      float64   `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	TotalAmount    float64   `gorm:"column:grand_total;not null" json:"totalAmount"`
	Currency       string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	DueDate        *time.Time `gorm:"column:due_date;type:date" json:"dueDate"`
	Status         string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	PostedAt       *time.Time `gorm:"column:posted_at" json:"postedAt"`
	GLPosted       bool      `gorm:"column:gl_posted;default:false;index" json:"glPosted"`
	GLJEID         *string   `gorm:"column:gl_je_id;size:36" json:"glJeId"`
	CreatedBy      string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines          []SupplierInvoiceLineGORM `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (SupplierInvoiceGORM) TableName() string { return "supplier_invoices" }

type SupplierInvoiceLineGORM struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceID  string    `gorm:"column:invoice_id;not null;size:36;index:idx_sinvline_inv,unique" json:"invoiceId"`
	LineNumber int       `gorm:"column:line_number;not null;index:idx_sinvline_inv,unique" json:"lineNumber"`
	ItemCode   string    `gorm:"column:item_code;not null;size:50" json:"itemCode"`
	ItemName   string    `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Quantity   float64   `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice  float64   `gorm:"column:unit_price;not null" json:"unitPrice"`
	LineTotal  float64   `gorm:"column:line_total;not null" json:"lineTotal"`
	TaxRate    *float64  `gorm:"column:tax_rate" json:"taxRate"`
	TaxAmount  *float64  `gorm:"column:tax_amount" json:"taxAmount"`
	AccountCode *string  `gorm:"column:account_code;size:20" json:"accountCode"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (SupplierInvoiceLineGORM) TableName() string { return "supplier_invoice_lines" }

type APTransactionGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	SupplierID  string    `gorm:"column:supplier_id;not null;size:36;index" json:"supplierId"`
	Amount      float64   `gorm:"column:amount;not null" json:"amount"`
	PaidAmount  float64   `gorm:"column:paid_amount;not null;default:0" json:"paidAmount"`
	Currency    string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	Status      string    `gorm:"column:status;not null;size:20;default:'OPEN';index" json:"status"`
	DueDate     *time.Time `gorm:"column:due_date;type:date" json:"dueDate"`
	CreatedBy   string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (APTransactionGORM) TableName() string { return "ap_transactions" }

type CostAllocationGORM struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceLineID string    `gorm:"column:invoice_line_id;not null;size:36;index" json:"invoiceLineId"`
	CostCenter    *string   `gorm:"column:cost_center;size:100" json:"costCenter"`
	DeptID        *string   `gorm:"column:dept_id;size:36" json:"deptId"`
	AllocPct      float64   `gorm:"column:alloc_pct;not null" json:"allocPct"`
	AllocAmount   float64   `gorm:"column:alloc_amount;not null" json:"allocAmount"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (CostAllocationGORM) TableName() string { return "cost_allocations" }
