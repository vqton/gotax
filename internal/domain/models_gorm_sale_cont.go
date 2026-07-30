package domain

import "time"

type CustomerInvoiceGORM struct {
	ID             string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID      string    `gorm:"column:company_id;not null;size:36;index:idx_cinv_company,unique" json:"companyId"`
	InvoiceNumber  string    `gorm:"column:invoice_number;not null;size:30;index:idx_cinv_company,unique" json:"invoiceNumber"`
	InvoiceDate    time.Time `gorm:"column:invoice_date;not null;type:date;index" json:"invoiceDate"`
	CustomerID     string    `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	ReferenceNo    *string   `gorm:"column:reference_no;size:100" json:"referenceNo"`
	Subtotal       float64   `gorm:"column:subtotal;not null" json:"subtotal"`
	TaxAmount      float64   `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	GrandTotal     float64   `gorm:"column:grand_total;not null" json:"grandTotal"`
	Currency       string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	DueDate        *time.Time `gorm:"column:due_date;type:date" json:"dueDate"`
	Status         string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	PostedAt       *time.Time `gorm:"column:posted_at" json:"postedAt"`
	GLPosted       bool      `gorm:"column:gl_posted;default:false;index" json:"glPosted"`
	GLJEID         *string   `gorm:"column:gl_je_id;size:36" json:"glJeId"`
	CreatedBy      string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines          []InvLineGORM `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (CustomerInvoiceGORM) TableName() string { return "customer_invoices" }

type InvLineGORM struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceID   string    `gorm:"column:invoice_id;not null;size:36;index:idx_invline_inv,unique" json:"invoiceId"`
	LineNumber  int       `gorm:"column:line_number;not null;index:idx_invline_inv,unique" json:"lineNumber"`
	ItemCode    string    `gorm:"column:item_code;not null;size:50" json:"itemCode"`
	ItemName    string    `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Description *string   `gorm:"column:description;type:text" json:"description"`
	Quantity    float64   `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice   float64   `gorm:"column:unit_price;not null" json:"unitPrice"`
	Discount    *float64  `gorm:"column:discount" json:"discount"`
	LineTotal   float64   `gorm:"column:line_total;not null" json:"lineTotal"`
	TaxRate     *float64  `gorm:"column:tax_rate" json:"taxRate"`
	TaxAmount   *float64  `gorm:"column:tax_amount" json:"taxAmount"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (InvLineGORM) TableName() string { return "invoice_lines" }

type CustomerReceiptGORM struct {
	ID             string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID      string    `gorm:"column:company_id;not null;size:36;index:idx_rcpt_company_number,unique" json:"companyId"`
	ReceiptNumber  string    `gorm:"column:receipt_number;not null;size:30;index:idx_rcpt_company_number,unique" json:"receiptNumber"`
	ReceiptDate    time.Time `gorm:"column:receipt_date;not null;type:date;index" json:"receiptDate"`
	CustomerID     string    `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	PaymentMethod  string    `gorm:"column:payment_method;not null;size:30" json:"paymentMethod"`
	TotalReceived  float64   `gorm:"column:total_received;not null" json:"totalReceived"`
	RefNo          *string   `gorm:"column:ref_no;size:100" json:"refNo"`
	Status         string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	GLPosted       bool      `gorm:"column:gl_posted;default:false;index" json:"glPosted"`
	GLJEID         *string   `gorm:"column:gl_je_id;size:36" json:"glJeId"`
	CreatedBy      string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Allocations    []RcpAllocationGORM `gorm:"foreignKey:ReceiptID;constraint:OnDelete:CASCADE" json:"allocations,omitempty"`
}

func (CustomerReceiptGORM) TableName() string { return "customer_receipts" }

type RcpAllocationGORM struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ReceiptID   string    `gorm:"column:receipt_id;not null;size:36;index:idx_rcp_alloc_rcpt,unique" json:"receiptId"`
	InvoiceID   string    `gorm:"column:invoice_id;not null;size:36;index:idx_rcp_alloc_rcpt,unique" json:"invoiceId"`
	Allocated   float64   `gorm:"column:allocated;not null" json:"allocated"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (RcpAllocationGORM) TableName() string { return "receipt_allocations" }

type CreditNoteGORM struct {
	ID            string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID     string    `gorm:"column:company_id;not null;size:36;index:idx_cn_company,unique" json:"companyId"`
	CNNumber      string    `gorm:"column:cn_number;not null;size:30;index:idx_cn_company,unique" json:"cnNumber"`
	CNDate        time.Time `gorm:"column:cn_date;not null;type:date;index" json:"cnDate"`
	InvoiceID     *string   `gorm:"column:invoice_id;size:36;index" json:"invoiceId"`
	CustomerID    string    `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	Reason        *string   `gorm:"column:reason;type:text" json:"reason"`
	Subtotal      float64   `gorm:"column:subtotal;not null" json:"subtotal"`
	TaxAmount     float64   `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	GrandTotal    float64   `gorm:"column:grand_total;not null" json:"grandTotal"`
	Status        string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	PostedAt      *time.Time `gorm:"column:posted_at" json:"postedAt"`
	GLPosted      bool      `gorm:"column:gl_posted;default:false;index" json:"glPosted"`
	GLJEID        *string   `gorm:"column:gl_je_id;size:36" json:"glJeId"`
	CreatedBy     string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines         []CNLineGORM `gorm:"foreignKey:CNID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (CreditNoteGORM) TableName() string { return "credit_notes" }

type CNLineGORM struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CNID       string    `gorm:"column:cn_id;not null;size:36;index:idx_cnline_cn,unique" json:"cnId"`
	LineNumber int       `gorm:"column:line_number;not null;index:idx_cnline_cn,unique" json:"lineNumber"`
	ItemCode   string    `gorm:"column:item_code;not null;size:50" json:"itemCode"`
	ItemName   string    `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Qty        float64   `gorm:"column:qty;not null" json:"qty"`
	UnitPrice  float64   `gorm:"column:unit_price;not null" json:"unitPrice"`
	LineTotal  float64   `gorm:"column:line_total;not null" json:"lineTotal"`
	Reason     *string   `gorm:"column:reason;type:text" json:"reason"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (CNLineGORM) TableName() string { return "cn_lines" }

type ARTransactionGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	CustomerID  string    `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	Amount      float64   `gorm:"column:amount;not null" json:"amount"`
	PaidAmount  float64   `gorm:"column:paid_amount;not null;default:0" json:"paidAmount"`
	Currency    string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	Status      string    `gorm:"column:status;not null;size:20;default:'OPEN';index" json:"status"`
	DueDate     *time.Time `gorm:"column:due_date;type:date" json:"dueDate"`
	CreatedBy   string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ARTransactionGORM) TableName() string { return "ar_transactions" }

type SalesQuotationGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	QuoteNumber     string    `gorm:"column:quote_number;not null;size:30;uniqueIndex:idx_sq_company,unique" json:"quoteNumber"`
	QuoteDate       time.Time `gorm:"column:quote_date;not null;type:date;index" json:"quoteDate"`
	CustomerID      string    `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	ValidUntil      time.Time `gorm:"column:valid_until;not null;type:date" json:"validUntil"`
	Subtotal        float64   `gorm:"column:subtotal;not null" json:"subtotal"`
	TotalAmount     float64   `gorm:"column:total_amount;not null" json:"totalAmount"`
	Currency        string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	Status          string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	ConvertedToSOID *string   `gorm:"column:converted_to_so_id;size:36" json:"convertedToSoId"`
	CreatedBy       string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (SalesQuotationGORM) TableName() string { return "sales_quotations" }
