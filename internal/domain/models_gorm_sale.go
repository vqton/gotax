package domain

import "time"

type CustomerGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID    string    `gorm:"column:company_id;not null;size:36;index:idx_customer_company_code,unique" json:"companyId"`
	Code         string    `gorm:"column:code;not null;size:20;index:idx_customer_company_code,unique" json:"code"`
	Name         string    `gorm:"column:name;not null;size:255" json:"name"`
	TaxCode      *string   `gorm:"column:tax_code;size:20;uniqueIndex" json:"taxCode"`
	Address      *string   `gorm:"column:address;type:text" json:"address"`
	Phone        *string   `gorm:"column:phone;size:20" json:"phone"`
	Email        *string   `gorm:"column:email;size:255" json:"email"`
	ContactPerson *string  `gorm:"column:contact_person;size:255" json:"contactPerson"`
	PaymentTerms *string   `gorm:"column:payment_terms;size:50" json:"paymentTerms"`
	CreditLimit  *float64  `gorm:"column:credit_limit" json:"creditLimit"`
	IsActive     bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (CustomerGORM) TableName() string { return "customers" }

type SalesOrderGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36;index:idx_so_company_number,unique" json:"companyId"`
	SONumber        string    `gorm:"column:so_number;not null;size:30;index:idx_so_company_number,unique" json:"soNumber"`
	OrderDate       time.Time `gorm:"column:order_date;not null;type:date;index" json:"orderDate"`
	CustomerID      string    `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	Subtotal        float64   `gorm:"column:subtotal;not null" json:"subtotal"`
	TaxAmount       float64   `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	TotalAmount     float64   `gorm:"column:total_amount;not null" json:"totalAmount"`
	Currency        string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	DeliveryDate    *time.Time `gorm:"column:delivery_date;type:date" json:"deliveryDate"`
	DeliveryAddress *string   `gorm:"column:delivery_address;type:text" json:"deliveryAddress"`
	PaymentTerms    *string   `gorm:"column:payment_terms;size:50" json:"paymentTerms"`
	Status          string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	ApprovedAt      *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	ApprovedBy      *string   `gorm:"column:approved_by;size:36" json:"approvedBy"`
	CreatedBy       string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines           []SOLineGORM `gorm:"foreignKey:SOID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (SalesOrderGORM) TableName() string { return "sales_orders" }

type SOLineGORM struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SOID       string    `gorm:"column:so_id;not null;size:36;index:idx_soline_so,unique" json:"soId"`
	LineNumber int       `gorm:"column:line_number;not null;index:idx_soline_so,unique" json:"lineNumber"`
	ItemCode   string    `gorm:"column:item_code;not null;size:50" json:"itemCode"`
	ItemName   string    `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Quantity   float64   `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice  float64   `gorm:"column:unit_price;not null" json:"unitPrice"`
	LineTotal  float64   `gorm:"column:line_total;not null" json:"lineTotal"`
	Discount   *float64  `gorm:"column:discount" json:"discount"`
	TaxRate    *float64  `gorm:"column:tax_rate" json:"taxRate"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (SOLineGORM) TableName() string { return "so_lines" }

type DeliveryNoteGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36;index:idx_dn_company,unique" json:"companyId"`
	DNNumber        string    `gorm:"column:dn_number;not null;size:30;index:idx_dn_company,unique" json:"dnNumber"`
	DNDate          time.Time `gorm:"column:dn_date;not null;type:date;index" json:"dnDate"`
	SOID            *string   `gorm:"column:so_id;size:36;index" json:"soId"`
	CustomerID      string    `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	Status          string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	TotalItems      int       `gorm:"column:total_items;not null;default:0" json:"totalItems"`
	CreatedBy       string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines           []DNLineGORM `gorm:"foreignKey:DNID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (DeliveryNoteGORM) TableName() string { return "delivery_notes" }

type DNLineGORM struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DNID        string    `gorm:"column:dn_id;not null;size:36;index:idx_dnline_dn,unique" json:"dnId"`
	LineNumber  int       `gorm:"column:line_number;not null;index:idx_dnline_dn,unique" json:"lineNumber"`
	SOLineID    *string   `gorm:"column:so_line_id;size:36" json:"soLineId"`
	ItemCode    string    `gorm:"column:item_code;not null;size:50" json:"itemCode"`
	QtyShipped  float64   `gorm:"column:qty_shipped;not null" json:"qtyShipped"`
	Unit        *string   `gorm:"column:unit;size:20" json:"unit"`
	Note        *string   `gorm:"column:note;type:text" json:"note"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (DNLineGORM) TableName() string { return "dn_lines" }
