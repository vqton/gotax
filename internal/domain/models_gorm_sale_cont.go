package domain

import "time"

// GORM persistence structs for the Sale (O2C) module — continued.
// Column names mirror domain models (models_sale.go) — do not drift.
// Tables created by migrations/000017_sale_schema.{up,down}.sql.

type CustomerInvoiceGORM struct {
	ID                 string         `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID          string         `gorm:"column:company_id;not null;size:36;index:idx_cinv_company_number,unique" json:"companyId"`
	InvoiceNumber      string         `gorm:"column:invoice_number;not null;size:30;index:idx_cinv_company_number,unique" json:"invoiceNumber"`
	InvoiceDate        time.Time      `gorm:"column:invoice_date;not null;type:date;index" json:"invoiceDate"`
	SOID               string         `gorm:"column:so_id;size:36" json:"soId"`
	DNID               string         `gorm:"column:dn_id;size:36" json:"dnId"`
	CustomerID         string         `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	CustomerName       string         `gorm:"column:customer_name;not null;size:255" json:"customerName"`
	CustomerTaxCode    string         `gorm:"column:customer_tax_code;not null;size:20" json:"customerTaxCode"`
	CustomerAddress    string         `gorm:"column:customer_address;type:text" json:"customerAddress"`
	InvoiceType        string         `gorm:"column:invoice_type;size:20" json:"invoiceType"`
	Currency           string         `gorm:"column:currency;size:3;default:VND" json:"currency"`
	ExchangeRate       float64        `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	Subtotal           float64        `gorm:"column:subtotal;not null" json:"subtotal"`
	DiscountAmount     float64        `gorm:"column:discount_amount;not null;default:0" json:"discountAmount"`
	TaxAmount          float64        `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	TotalAmount        float64        `gorm:"column:total_amount;not null" json:"totalAmount"`
	AmountReceived     float64        `gorm:"column:amount_received;not null;default:0" json:"amountReceived"`
	BalanceDue         float64        `gorm:"column:balance_due;not null;default:0" json:"balanceDue"`
	DueDate            *time.Time     `gorm:"column:due_date;type:date" json:"dueDate"`
	InvoiceNote        string         `gorm:"column:invoice_note;type:text" json:"invoiceNote"`
	EInvoiceData       string         `gorm:"column:e_invoice_data;type:text" json:"eInvoiceData"`
	EInvoiceCode       string         `gorm:"column:e_invoice_code;size:50" json:"eInvoiceCode"`
	EInvStatus         string         `gorm:"column:e_invoice_status;size:20;default:pending" json:"eInvStatus"`
	DigitalSignatureID string         `gorm:"column:digital_signature_id;size:36" json:"digitalSignatureId"`
	SignedData         string         `gorm:"column:signed_data;type:text" json:"signedData"`
	GDTResponse        string         `gorm:"column:gdt_response;type:text" json:"gdtResponse"`
	OriginalInvoiceID  string         `gorm:"column:original_invoice_id;size:36" json:"originalInvoiceId"`
	AdjustmentType     string         `gorm:"column:adjustment_type;size:20" json:"adjustmentType"`
	Status             string         `gorm:"column:status;size:20;default:DRAFT;index" json:"status"`
	GLPosted           bool           `gorm:"column:gl_posted;not null;default:false" json:"glPosted"`
	GLPostedAt         *time.Time     `gorm:"column:gl_posted_at" json:"glPostedAt"`
	Notes              string         `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy          string         `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines              []InvLineGORM  `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (CustomerInvoiceGORM) TableName() string { return "customer_invoices" }

type InvLineGORM struct {
	ID             string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	InvoiceID      string  `gorm:"column:invoice_id;not null;size:36;index:idx_invline_inv" json:"invoiceId"`
	SOLineID       string  `gorm:"column:so_line_id;size:36" json:"soLineId"`
	DNLineID       string  `gorm:"column:dn_line_id;size:36" json:"dnLineId"`
	ItemCode       string  `gorm:"column:item_code;size:50" json:"itemCode"`
	ItemName       string  `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Unit           string  `gorm:"column:unit;size:20" json:"unit"`
	Quantity       float64 `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice      float64 `gorm:"column:unit_price;not null" json:"unitPrice"`
	DiscountPct    float64 `gorm:"column:discount_pct;not null;default:0" json:"discountPct"`
	VATRate        float64 `gorm:"column:vat_rate;not null;default:0" json:"vatRate"`
	VATType        string  `gorm:"column:vat_type;size:10;default:VAT10" json:"vatType"`
	LineTotal      float64 `gorm:"column:line_total;not null" json:"lineTotal"`
	LineVATAmount  float64 `gorm:"column:line_vat_amount;not null;default:0" json:"lineVatAmount"`
	RevenueAccount string  `gorm:"column:revenue_account_id;not null;size:20" json:"revenueAccount"`
	VATAccountID   string  `gorm:"column:vat_account_id;not null;size:20" json:"vatAccount"`
}

func (InvLineGORM) TableName() string { return "invoice_lines" }

type CustomerReceiptGORM struct {
	ID                string              `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string              `gorm:"column:company_id;not null;size:36;index:idx_rcpt_company_number,unique" json:"companyId"`
	ReceiptNumber     string              `gorm:"column:receipt_number;not null;size:30;index:idx_rcpt_company_number,unique" json:"receiptNumber"`
	CustomerID        string              `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	ReceiptDate       time.Time           `gorm:"column:receipt_date;not null;type:date;index" json:"receiptDate"`
	PaymentMethod     string              `gorm:"column:payment_method;size:30" json:"paymentMethod"`
	BankAccountID     string              `gorm:"column:bank_account_id;size:36" json:"bankAccountId"`
	Currency          string              `gorm:"column:currency;size:3;default:VND" json:"currency"`
	ExchangeRate      float64             `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	Amount            float64             `gorm:"column:amount;not null" json:"amount"`
	UnallocatedAmount float64             `gorm:"column:unallocated_amount;not null;default:0" json:"unallocatedAmount"`
	Reference         string              `gorm:"column:reference;size:100" json:"reference"`
	Notes             string              `gorm:"column:notes;type:text" json:"notes"`
	Status            string              `gorm:"column:status;size:20;default:DRAFT;index" json:"status"`
	GLPosted          bool                `gorm:"column:gl_posted;not null;default:false" json:"glPosted"`
	GLPostedAt        *time.Time          `gorm:"column:gl_posted_at" json:"glPostedAt"`
	CreatedBy         string              `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt         time.Time           `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time           `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Allocations       []RcpAllocationGORM `gorm:"foreignKey:ReceiptID;constraint:OnDelete:CASCADE" json:"allocations,omitempty"`
}

func (CustomerReceiptGORM) TableName() string { return "customer_receipts" }

type RcpAllocationGORM struct {
	ID              string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	ReceiptID       string  `gorm:"column:receipt_id;not null;size:36;index:idx_rcp_alloc_rcpt" json:"receiptId"`
	InvoiceID       string  `gorm:"column:invoice_id;not null;size:36" json:"invoiceId"`
	AllocatedAmount float64 `gorm:"column:allocated_amount;not null" json:"allocatedAmount"`
	DiscountAmount  float64 `gorm:"column:discount_amount;not null;default:0" json:"discountAmount"`
}

func (RcpAllocationGORM) TableName() string { return "receipt_allocations" }

type CreditNoteGORM struct {
	ID                string        `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string        `gorm:"column:company_id;not null;size:36;index:idx_cn_company_number,unique" json:"companyId"`
	CNNumber          string        `gorm:"column:cn_number;not null;size:30;index:idx_cn_company_number,unique" json:"cnNumber"`
	OriginalInvoiceID string        `gorm:"column:original_invoice_id;not null;size:36;index" json:"originalInvoiceId"`
	CustomerID        string        `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	ReturnDate        time.Time     `gorm:"column:return_date;not null;type:date;index" json:"returnDate"`
	ReturnReason      string        `gorm:"column:return_reason;type:text" json:"returnReason"`
	ReturnType        string        `gorm:"column:return_type;size:20;default:partial_return" json:"returnType"`
	DNID              string        `gorm:"column:dn_id;size:36" json:"dnId"`
	Subtotal          float64       `gorm:"column:subtotal;not null" json:"subtotal"`
	TaxAmount         float64       `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	TotalAmount       float64       `gorm:"column:total_amount;not null" json:"totalAmount"`
	EInvoiceData      string        `gorm:"column:e_invoice_data;type:text" json:"eInvoiceData"`
	EInvoiceCode      string        `gorm:"column:e_invoice_code;size:50" json:"eInvoiceCode"`
	Status            string        `gorm:"column:status;size:20;default:DRAFT;index" json:"status"`
	GLPosted          bool          `gorm:"column:gl_posted;not null;default:false" json:"glPosted"`
	GLPostedAt        *time.Time    `gorm:"column:gl_posted_at" json:"glPostedAt"`
	Notes             string        `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy         string        `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt         time.Time     `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines             []CNLineGORM  `gorm:"foreignKey:CNID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (CreditNoteGORM) TableName() string { return "credit_notes" }

type CNLineGORM struct {
	ID            string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	CNID          string  `gorm:"column:cn_id;not null;size:36;index:idx_cnline_cn" json:"cnId"`
	InvLineID     string  `gorm:"column:invoice_line_id;size:36" json:"invLineId"`
	ItemName      string  `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Unit          string  `gorm:"column:unit;size:20" json:"unit"`
	Quantity      float64 `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice     float64 `gorm:"column:unit_price;not null" json:"unitPrice"`
	VATRate       float64 `gorm:"column:vat_rate;not null;default:0" json:"vatRate"`
	LineTotal     float64 `gorm:"column:line_total;not null" json:"lineTotal"`
	LineVATAmount float64 `gorm:"column:line_vat_amount;not null;default:0" json:"lineVatAmount"`
}

func (CNLineGORM) TableName() string { return "cn_lines" }

type ARTransactionGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	CustomerID      string    `gorm:"column:customer_id;not null;size:36;index" json:"customerId"`
	InvoiceID       string    `gorm:"column:invoice_id;size:36" json:"invoiceId"`
	TransactionType string    `gorm:"column:transaction_type;not null;size:20" json:"transactionType"`
	TransactionDate time.Time `gorm:"column:transaction_date;not null;type:date;index" json:"transactionDate"`
	Amount          float64   `gorm:"column:amount;not null" json:"amount"`
	Currency        string    `gorm:"column:currency;size:3;default:VND" json:"currency"`
	ReferenceType   string    `gorm:"column:reference_type;size:50" json:"referenceType"`
	ReferenceID     string    `gorm:"column:reference_id;size:36" json:"referenceId"`
	Notes           string    `gorm:"column:notes;type:text" json:"notes"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (ARTransactionGORM) TableName() string { return "ar_transactions" }

type SalesQuotationGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	QNNumber    string    `gorm:"column:qn_number;not null;size:30;uniqueIndex:idx_sq_company,unique" json:"qnNumber"`
	CustomerID  string    `gorm:"column:customer_id;size:36" json:"customerId"`
	ValidUntil  *time.Time `gorm:"column:valid_until;type:date" json:"validUntil"`
	Status      string    `gorm:"column:status;size:20;default:DRAFT" json:"status"`
	TotalAmount float64   `gorm:"column:total_amount;not null" json:"totalAmount"`
	CreatedBy   string    `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (SalesQuotationGORM) TableName() string { return "sales_quotations" }
