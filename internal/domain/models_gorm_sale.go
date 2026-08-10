package domain

import "time"

// GORM persistence structs for the Sale (O2C) module.
// Column names mirror domain models (models_sale.go) — do not drift.
// Tables created by migrations/000017_sale_schema.{up,down}.sql.

type CustomerGORM struct {
	ID                string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string    `gorm:"column:company_id;not null;size:36;index:idx_customer_company_code,unique" json:"companyId"`
	Code              string    `gorm:"column:code;not null;size:20;index:idx_customer_company_code,unique" json:"code"`
	Name              string    `gorm:"column:name;not null;size:255" json:"name"`
	TaxCode           string    `gorm:"column:tax_code;not null;size:20" json:"taxCode"`
	Address           string    `gorm:"column:address;type:text" json:"address"`
	Phone             string    `gorm:"column:phone;size:20" json:"phone"`
	Email             string    `gorm:"column:email;size:255" json:"email"`
	BankAccountName   string    `gorm:"column:bank_account_name;size:255" json:"bankAccountName"`
	BankAccountNumber string    `gorm:"column:bank_account_number;size:50" json:"bankAccountNumber"`
	BankName          string    `gorm:"column:bank_name;size:255" json:"bankName"`
	PaymentTerms      string    `gorm:"column:payment_terms;size:50" json:"paymentTerms"`
	CreditLimit       float64   `gorm:"column:credit_limit;not null;default:0" json:"creditLimit"`
	Currency          string    `gorm:"column:currency;size:3;default:VND" json:"currency"`
	CustomerType      string    `gorm:"column:customer_type;size:20;default:domestic" json:"customerType"`
	CustomerGroup     string    `gorm:"column:customer_group;size:20" json:"customerGroup"`
	PriceListID       string    `gorm:"column:price_list_id;size:36" json:"priceListId"`
	Status            string    `gorm:"column:status;size:20;default:ACTIVE;index" json:"status"`
	Notes             string    `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy         string    `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (CustomerGORM) TableName() string { return "customers" }

type SalesOrderGORM struct {
	ID              string       `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string       `gorm:"column:company_id;not null;size:36;index:idx_so_company_number,unique" json:"companyId"`
	SONumber        string       `gorm:"column:so_number;not null;size:30;index:idx_so_company_number,unique" json:"soNumber"`
	QuotationID     string       `gorm:"column:quotation_id;size:36" json:"quotationId"`
	CustomerID      string       `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	OrderDate       time.Time    `gorm:"column:order_date;not null;type:date;index" json:"orderDate"`
	ExpectedDate    *time.Time   `gorm:"column:expected_date;type:date" json:"expectedDate"`
	Currency        string       `gorm:"column:currency;size:3;default:VND" json:"currency"`
	ExchangeRate    float64      `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	PaymentTerms    string       `gorm:"column:payment_terms;size:50" json:"paymentTerms"`
	DeliveryTerms   string       `gorm:"column:delivery_terms;size:50" json:"deliveryTerms"`
	ShippingAddress string       `gorm:"column:shipping_address;type:text" json:"shippingAddress"`
	Subtotal        float64      `gorm:"column:subtotal;not null" json:"subtotal"`
	DiscountAmount  float64      `gorm:"column:discount_amount;not null;default:0" json:"discountAmount"`
	TaxAmount       float64      `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	TotalAmount     float64      `gorm:"column:total_amount;not null" json:"totalAmount"`
	Status          string       `gorm:"column:status;size:20;default:DRAFT;index" json:"status"`
	ApprovedBy      string       `gorm:"column:approved_by;size:36" json:"approvedBy"`
	ApprovedAt      *time.Time   `gorm:"column:approved_at" json:"approvedAt"`
	CancelledReason string       `gorm:"column:cancelled_reason;type:text" json:"cancelledReason"`
	Notes           string       `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy       string       `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt       time.Time    `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines           []SOLineGORM `gorm:"foreignKey:SOID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (SalesOrderGORM) TableName() string { return "sales_orders" }

type SOLineGORM struct {
	ID              string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	SOID            string  `gorm:"column:so_id;not null;size:36;index:idx_soline_so" json:"soId"`
	LineNumber      int     `gorm:"column:line_number;not null" json:"lineNumber"`
	ItemCode        string  `gorm:"column:item_code;size:50" json:"itemCode"`
	ItemName        string  `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Unit            string  `gorm:"column:unit;not null;size:20" json:"unit"`
	Quantity        float64 `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice       float64 `gorm:"column:unit_price;not null" json:"unitPrice"`
	DiscountPct     float64 `gorm:"column:discount_pct;not null;default:0" json:"discountPct"`
	VATRate         float64 `gorm:"column:vat_rate;not null;default:0" json:"vatRate"`
	VATType         string  `gorm:"column:vat_type;size:10;default:VAT10" json:"vatType"`
	RevenueAccount  string  `gorm:"column:revenue_account_id;not null;size:20" json:"revenueAccount"`
	VATAccountID    string  `gorm:"column:vat_account_id;not null;size:20" json:"vatAccount"`
	LineTotal       float64 `gorm:"column:line_total;not null" json:"lineTotal"`
	LineVATAmount   float64 `gorm:"column:line_vat_amount;not null;default:0" json:"lineVatAmount"`
	DeliveredQty    float64 `gorm:"column:delivered_qty;not null;default:0" json:"deliveredQty"`
	InvoicedQty     float64 `gorm:"column:invoiced_qty;not null;default:0" json:"invoicedQty"`
}

func (SOLineGORM) TableName() string { return "so_lines" }

type DeliveryNoteGORM struct {
	ID              string       `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string       `gorm:"column:company_id;not null;size:36;index:idx_dn_company_number,unique" json:"companyId"`
	DNNumber        string       `gorm:"column:dn_number;not null;size:30;index:idx_dn_company_number,unique" json:"dnNumber"`
	SOID            string       `gorm:"column:so_id;not null;size:36;index" json:"soId"`
	DeliveryDate    time.Time    `gorm:"column:delivery_date;not null;type:date;index" json:"deliveryDate"`
	Warehouse       string       `gorm:"column:warehouse;size:50" json:"warehouse"`
	ShippingMethod  string       `gorm:"column:shipping_method;size:50" json:"shippingMethod"`
	CarrierName     string       `gorm:"column:carrier_name;size:100" json:"carrierName"`
	TrackingNumber  string       `gorm:"column:tracking_number;size:100" json:"trackingNumber"`
	DeliveryAddress string       `gorm:"column:delivery_address;type:text" json:"deliveryAddress"`
	Status          string       `gorm:"column:status;size:20;default:DRAFT;index" json:"status"`
	Notes           string       `gorm:"column:notes;type:text" json:"notes"`
	TolerancePercent float64    `gorm:"column:tolerance_percent;not null;default:5" json:"tolerancePercent"`
	GLPosted        bool         `gorm:"column:gl_posted;not null;default:false" json:"glPosted"`
	GLPostedAt      *time.Time   `gorm:"column:gl_posted_at" json:"glPostedAt"`
	CreatedBy       string       `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt       time.Time    `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines           []DNLineGORM `gorm:"foreignKey:DNID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (DeliveryNoteGORM) TableName() string { return "delivery_notes" }

type DNLineGORM struct {
	ID           string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	DNID         string  `gorm:"column:dn_id;not null;size:36;index:idx_dnline_dn" json:"dnId"`
	SOLineID     string  `gorm:"column:so_line_id;size:36" json:"soLineId"`
	ItemCode     string  `gorm:"column:item_code;size:50" json:"itemCode"`
	ItemName     string  `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Unit         string  `gorm:"column:unit;size:20" json:"unit"`
	QtyDelivered float64 `gorm:"column:qty_delivered;not null" json:"qtyDelivered"`
	QtyReturned  float64 `gorm:"column:qty_returned;not null;default:0" json:"qtyReturned"`
	UnitPrice    float64 `gorm:"column:unit_price;not null" json:"unitPrice"`
	LineTotal    float64 `gorm:"column:line_total;not null" json:"lineTotal"`
	CostPrice    float64 `gorm:"column:cost_price;not null;default:0" json:"costPrice"`
}

func (DNLineGORM) TableName() string { return "dn_lines" }

type PriceListGORM struct {
	ID          string           `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string           `gorm:"column:company_id;not null;size:36;index:idx_pl_company_code,unique" json:"companyId"`
	Code        string           `gorm:"column:code;not null;size:20;index:idx_pl_company_code,unique" json:"code"`
	Name        string           `gorm:"column:name;not null;size:255" json:"name"`
	Description string           `gorm:"column:description;type:text" json:"description"`
	Currency    string           `gorm:"column:currency;size:3;default:VND" json:"currency"`
	IsActive    bool             `gorm:"column:is_active;not null;default:true" json:"isActive"`
	CreatedAt   time.Time        `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines       []PriceListLineGORM `gorm:"foreignKey:PriceListID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (PriceListGORM) TableName() string { return "price_lists" }

type PriceListLineGORM struct {
	ID            string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	PriceListID   string  `gorm:"column:price_list_id;not null;size:36;index:idx_pll_pl" json:"priceListId"`
	ItemCode      string  `gorm:"column:item_code;not null;size:50" json:"itemCode"`
	ItemName      string  `gorm:"column:item_name;size:255" json:"itemName"`
	Unit          string  `gorm:"column:unit;size:20" json:"unit"`
	UnitPrice     float64 `gorm:"column:unit_price;not null" json:"unitPrice"`
	VATRate       float64 `gorm:"column:vat_rate;not null;default:0" json:"vatRate"`
	MinQuantity   float64 `gorm:"column:min_quantity;not null;default:0" json:"minQuantity"`
	EffectiveFrom string  `gorm:"column:effective_from;size:10" json:"effectiveFrom"`
	EffectiveTo   string  `gorm:"column:effective_to;size:10" json:"effectiveTo"`
}

func (PriceListLineGORM) TableName() string { return "price_list_lines" }
