package domain

import "time"

type SupplierGORM struct {
	ID                string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string    `gorm:"column:company_id;not null;size:36;index:idx_suppliers_company" json:"companyId"`
	Code              string    `gorm:"column:code;not null;size:20;uniqueIndex:idx_suppliers_company_code" json:"code"`
	Name              string    `gorm:"column:name;not null;size:255" json:"name"`
	TaxCode           string    `gorm:"column:tax_code;not null;size:20" json:"taxCode"`
	Address           string    `gorm:"column:address;not null;default:''" json:"address"`
	Phone             string    `gorm:"column:phone;not null;default:''" json:"phone"`
	Email             string    `gorm:"column:email;not null;default:''" json:"email"`
	BankAccountName   string    `gorm:"column:bank_account_name;not null;default:''" json:"bankAccountName"`
	BankAccountNumber string    `gorm:"column:bank_account_number;not null;default:''" json:"bankAccountNumber"`
	BankName          string    `gorm:"column:bank_name;not null;default:''" json:"bankName"`
	PaymentTerms      string    `gorm:"column:payment_terms;not null;default:net30" json:"paymentTerms"`
	CreditLimit       float64   `gorm:"column:credit_limit;not null;default:0" json:"creditLimit"`
	Currency          string    `gorm:"column:currency;not null;default:VND" json:"currency"`
	SupplierType      string    `gorm:"column:supplier_type;not null;default:domestic" json:"supplierType"`
	Status            string    `gorm:"column:status;not null;default:ACTIVE" json:"status"`
	Notes             string    `gorm:"column:notes;not null;default:''" json:"notes"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (SupplierGORM) TableName() string { return "suppliers" }

type PurchaseOrderGORM struct {
	ID              string           `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string           `gorm:"column:company_id;not null;size:36;index:idx_po_company" json:"companyId"`
	PONumber        string           `gorm:"column:po_number;not null;size:30;index:idx_po_number" json:"poNumber"`
	SupplierID      string           `gorm:"column:supplier_id;not null;size:36;index:idx_po_supplier" json:"supplierId"`
	RequisitionID   string           `gorm:"column:requisition_id;not null;default:''" json:"requisitionId"`
	OrderDate       string           `gorm:"column:order_date;not null" json:"orderDate"`
	ExpectedDate    string           `gorm:"column:expected_date;not null;default:''" json:"expectedDate"`
	Currency        string           `gorm:"column:currency;not null;default:VND" json:"currency"`
	ExchangeRate    float64          `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	PaymentTerms    string           `gorm:"column:payment_terms;not null;default:''" json:"paymentTerms"`
	DeliveryTerms   string           `gorm:"column:delivery_terms;not null;default:''" json:"deliveryTerms"`
	Subtotal        float64          `gorm:"column:subtotal;not null;default:0" json:"subtotal"`
	DiscountAmount  float64          `gorm:"column:discount_amount;not null;default:0" json:"discountAmount"`
	TaxAmount       float64          `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	TotalAmount     float64          `gorm:"column:total_amount;not null;default:0" json:"totalAmount"`
	Status          string           `gorm:"column:status;not null;default:DRAFT;index:idx_po_status" json:"status"`
	ApprovedBy      string           `gorm:"column:approved_by;not null;default:''" json:"approvedBy"`
	ApprovedAt      string           `gorm:"column:approved_at;not null;default:''" json:"approvedAt"`
	CancelledReason string           `gorm:"column:cancelled_reason;not null;default:''" json:"cancelledReason"`
	Notes           string           `gorm:"column:notes;not null;default:''" json:"notes"`
	CreatedBy       string           `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt       time.Time        `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time        `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines           []POItemGORM     `gorm:"foreignKey:POID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (PurchaseOrderGORM) TableName() string { return "purchase_orders" }

type POItemGORM struct {
	ID            string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	POID          string    `gorm:"column:po_id;not null;size:36;index:idx_po_lines_po" json:"poId"`
	LineNumber    int       `gorm:"column:line_number;not null" json:"lineNumber"`
	ItemCode      string    `gorm:"column:item_code;not null;default:''" json:"itemCode"`
	ItemName      string    `gorm:"column:item_name;not null" json:"itemName"`
	Unit          string    `gorm:"column:unit;not null" json:"unit"`
	Quantity      float64   `gorm:"column:quantity;not null;default:0" json:"quantity"`
	UnitPrice     float64   `gorm:"column:unit_price;not null;default:0" json:"unitPrice"`
	DiscountPct   float64   `gorm:"column:discount_pct;not null;default:0" json:"discountPct"`
	VATRate       float64   `gorm:"column:vat_rate;not null;default:0" json:"vatRate"`
	VATType       string    `gorm:"column:vat_type;not null;default:VAT_0" json:"vatType"`
	AccountID     string    `gorm:"column:account_id;not null" json:"accountId"`
	VATAccountID  string    `gorm:"column:vat_account_id;not null" json:"vatAccountId"`
	LineTotal     float64   `gorm:"column:line_total;not null;default:0" json:"lineTotal"`
	LineVATAmount float64   `gorm:"column:line_vat_amount;not null;default:0" json:"lineVatAmount"`
	ReceivedQty   float64   `gorm:"column:received_qty;not null;default:0" json:"receivedQty"`
	InvoicedQty   float64   `gorm:"column:invoiced_qty;not null;default:0" json:"invoicedQty"`
}

func (POItemGORM) TableName() string { return "po_lines" }

type GRNGORM struct {
	ID          string        `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string        `gorm:"column:company_id;not null;size:36;index:idx_grn_company" json:"companyId"`
	GRNNumber   string        `gorm:"column:grn_number;not null;size:30;index:idx_grn_number" json:"grnNumber"`
	POID        string        `gorm:"column:po_id;not null;size:36;index:idx_grn_po" json:"poId"`
	ReturnOfGRNID string      `gorm:"column:return_of_grn_id;not null;default:'';index:idx_grn_return_of" json:"returnOfGrnId"`
	ReceiptDate string        `gorm:"column:receipt_date;not null" json:"receiptDate"`
	Warehouse   string        `gorm:"column:warehouse;not null;default:''" json:"warehouse"`
	Status      string        `gorm:"column:status;not null;default:DRAFT;index:idx_grn_status" json:"status"`
	Notes       string        `gorm:"column:notes;not null;default:''" json:"notes"`
	CreatedBy   string        `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt   time.Time     `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	Lines       []GRNItemGORM `gorm:"foreignKey:GRNID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (GRNGORM) TableName() string { return "goods_receipt_notes" }

type GRNItemGORM struct {
	ID               string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	GRNID            string    `gorm:"column:grn_id;not null;size:36;index:idx_grn_lines_grn" json:"grnId"`
	POLineID         string    `gorm:"column:po_line_id;not null" json:"poLineId"`
	ItemCode         string    `gorm:"column:item_code;not null;default:''" json:"itemCode"`
	ItemName         string    `gorm:"column:item_name;not null" json:"itemName"`
	Unit             string    `gorm:"column:unit;not null;default:''" json:"unit"`
	QuantityReceived float64   `gorm:"column:quantity_received;not null;default:0" json:"quantityReceived"`
	QuantityRejected float64   `gorm:"column:quantity_rejected;not null;default:0" json:"quantityRejected"`
	UnitPrice        float64   `gorm:"column:unit_price;not null;default:0" json:"unitPrice"`
	LineTotal        float64   `gorm:"column:line_total;not null;default:0" json:"lineTotal"`
}

func (GRNItemGORM) TableName() string { return "grn_lines" }

type SupplierInvoiceGORM struct {
	ID                string                  `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string                  `gorm:"column:company_id;not null;size:36;index:idx_inv_company" json:"companyId"`
	InvoiceNumber     string                  `gorm:"column:invoice_number;not null;size:30;index:idx_inv_number" json:"invoiceNumber"`
	InvoiceDate       string                  `gorm:"column:invoice_date;not null" json:"invoiceDate"`
	POID              string                  `gorm:"column:po_id;not null;default:''" json:"poId"`
	GRNID             string                  `gorm:"column:grn_id;not null;default:''" json:"grnId"`
	SupplierID        string                  `gorm:"column:supplier_id;not null;size:36;index:idx_inv_supplier" json:"supplierId"`
	SupplierName      string                  `gorm:"column:supplier_name;not null" json:"supplierName"`
	SupplierTaxCode   string                  `gorm:"column:supplier_tax_code;not null" json:"supplierTaxCode"`
	InvoiceType       string                  `gorm:"column:invoice_type;not null;default:domestic" json:"invoiceType"`
	OriginalInvoiceID string                  `gorm:"column:original_invoice_id;not null;default:'';index:idx_inv_original" json:"originalInvoiceId"`
	ImportDuty        float64                 `gorm:"column:import_duty;not null;default:0" json:"importDuty"`
	ImportVAT         float64                 `gorm:"column:import_vat;not null;default:0" json:"importVat"`
	CustomsDeclarationNumber string           `gorm:"column:customs_declaration_number;not null;default:'';index:idx_inv_customs" json:"customsDeclarationNumber"`
	Currency          string                  `gorm:"column:currency;not null;default:VND" json:"currency"`
	ExchangeRate      float64                 `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	Subtotal          float64                 `gorm:"column:subtotal;not null;default:0" json:"subtotal"`
	DiscountAmount    float64                 `gorm:"column:discount_amount;not null;default:0" json:"discountAmount"`
	TaxAmount         float64                 `gorm:"column:tax_amount;not null;default:0" json:"taxAmount"`
	TotalAmount       float64                 `gorm:"column:total_amount;not null;default:0" json:"totalAmount"`
	AmountPaid        float64                 `gorm:"column:amount_paid;not null;default:0" json:"amountPaid"`
	BalanceDue        float64                 `gorm:"column:balance_due;not null;default:0" json:"balanceDue"`
	DueDate           string                  `gorm:"column:due_date;not null;default:''" json:"dueDate"`
	VATDeductionStatus string                 `gorm:"column:vat_deduction_status;not null;default:pending" json:"vatDeductionStatus"`
	EInvoiceData      string                  `gorm:"column:e_invoice_data;not null;default:''" json:"eInvoiceData"`
	EInvoiceCode      string                  `gorm:"column:e_invoice_code;not null;default:''" json:"eInvoiceCode"`
	Status            string                  `gorm:"column:status;not null;default:DRAFT;index:idx_inv_status" json:"status"`
	GLPosted          bool                    `gorm:"column:gl_posted;not null;default:false" json:"glPosted"`
	GLPostedAt        string                  `gorm:"column:gl_posted_at;not null;default:''" json:"glPostedAt"`
	Notes             string                  `gorm:"column:notes;not null;default:''" json:"notes"`
	CreatedBy         string                  `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt         time.Time               `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	Lines             []SupplierInvoiceLineGORM `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (SupplierInvoiceGORM) TableName() string { return "supplier_invoices" }

type SupplierInvoiceLineGORM struct {
	ID            string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	InvoiceID     string    `gorm:"column:invoice_id;not null;size:36;index:idx_inv_lines_inv" json:"invoiceId"`
	POLineID      string    `gorm:"column:po_line_id;not null;default:''" json:"poLineId"`
	GRNLineID     string    `gorm:"column:grn_line_id;not null;default:''" json:"grnLineId"`
	LineNumber    int       `gorm:"column:line_number;not null" json:"lineNumber"`
	ItemCode      string    `gorm:"column:item_code;not null;default:''" json:"itemCode"`
	ItemName      string    `gorm:"column:item_name;not null" json:"itemName"`
	Unit          string    `gorm:"column:unit;not null;default:''" json:"unit"`
	Quantity      float64   `gorm:"column:quantity;not null;default:0" json:"quantity"`
	UnitPrice     float64   `gorm:"column:unit_price;not null;default:0" json:"unitPrice"`
	VATRate       float64   `gorm:"column:vat_rate;not null;default:0" json:"vatRate"`
	VATType       string    `gorm:"column:vat_type;not null;default:VAT_0" json:"vatType"`
	LineTotal     float64   `gorm:"column:line_total;not null;default:0" json:"lineTotal"`
	LineVATAmount float64   `gorm:"column:line_vat_amount;not null;default:0" json:"lineVatAmount"`
	AccountID     string    `gorm:"column:account_id;not null" json:"accountId"`
	VATAccountID  string    `gorm:"column:vat_account_id;not null" json:"vatAccountId"`
}

func (SupplierInvoiceLineGORM) TableName() string { return "invoice_lines" }

type APTransactionGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36;index:idx_apt_company" json:"companyId"`
	SupplierID      string    `gorm:"column:supplier_id;not null;size:36;index:idx_apt_supplier" json:"supplierId"`
	InvoiceID       string    `gorm:"column:invoice_id;not null;default:'';index:idx_apt_invoice" json:"invoiceId"`
	TransactionType string    `gorm:"column:transaction_type;not null" json:"transactionType"`
	TransactionDate string    `gorm:"column:transaction_date;not null" json:"transactionDate"`
	Amount          float64   `gorm:"column:amount;not null;default:0" json:"amount"`
	Currency        string    `gorm:"column:currency;not null;default:VND" json:"currency"`
	ReferenceType   string    `gorm:"column:reference_type;not null;default:''" json:"referenceType"`
	ReferenceID     string    `gorm:"column:reference_id;not null;default:''" json:"referenceId"`
	Notes           string    `gorm:"column:notes;not null;default:''" json:"notes"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (APTransactionGORM) TableName() string { return "ap_transactions" }

type CostAllocationGORM struct {
	ID               string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID        string    `gorm:"column:company_id;not null;size:36" json:"companyId"`
	InvoiceID        string    `gorm:"column:invoice_id;not null;default:'';index:idx_costalloc_invoice" json:"invoiceId"`
	CostType         string    `gorm:"column:cost_type;not null;default:''" json:"costType"`
	CostAmount       float64   `gorm:"column:cost_amount;not null;default:0" json:"costAmount"`
	AllocationMethod string    `gorm:"column:allocation_method;not null;default:''" json:"allocationMethod"`
	AllocatedLines   string    `gorm:"column:allocated_lines;not null;default:''" json:"allocatedLines"`
	Notes            string    `gorm:"column:notes;not null;default:''" json:"notes"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (CostAllocationGORM) TableName() string { return "cost_allocations" }

type DoubtfulDebtProvisionGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36;index:idx_ddp_company" json:"companyId"`
	AsOfDate        time.Time `gorm:"column:as_of_date;not null;type:date;index:idx_ddp_date" json:"asOfDate"`
	TotalOutstanding float64  `gorm:"column:total_outstanding;not null;default:0" json:"totalOutstanding"`
	TotalProvision  float64   `gorm:"column:total_provision;not null;default:0" json:"totalProvision"`
	Status          string    `gorm:"column:status;not null;size:20;default:DRAFT;index:idx_ddp_status" json:"status"`
	CreatedBy       string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (DoubtfulDebtProvisionGORM) TableName() string { return "doubtful_debt_provisions" }

type DoubtfulDebtProvisionLineGORM struct {
	ID               string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	ProvisionID      string  `gorm:"column:provision_id;not null;size:36;index:idx_ddpl_provision" json:"provisionId"`
	SupplierID       string  `gorm:"column:supplier_id;not null;size:36;index:idx_ddpl_supplier" json:"supplierId"`
	SupplierName     string  `gorm:"column:supplier_name;not null;size:255" json:"supplierName"`
	TaxCode          string  `gorm:"column:tax_code;size:50;default:''" json:"taxCode"`
	OutstandingAmount float64 `gorm:"column:outstanding_amount;not null;default:0" json:"outstandingAmount"`
	AgeMonths        int     `gorm:"column:age_months;not null;default:0" json:"ageMonths"`
	RatePct          float64 `gorm:"column:rate_pct;not null;default:0" json:"ratePct"`
	ProvisionAmount  float64 `gorm:"column:provision_amount;not null;default:0" json:"provisionAmount"`
}

func (DoubtfulDebtProvisionLineGORM) TableName() string { return "doubtful_debt_provision_lines" }

type RequisitionGORM struct {
	ID                string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string    `gorm:"column:company_id;not null;size:36;index:idx_req_company" json:"companyId"`
	RequisitionNumber string    `gorm:"column:requisition_number;not null;size:50;index:idx_req_number" json:"requisitionNumber"`
	RequesterID       string    `gorm:"column:requester_id;not null;size:36" json:"requesterId"`
	RequesterName     string    `gorm:"column:requester_name;size:255;default:''" json:"requesterName"`
	DepartmentID      string    `gorm:"column:department_id;size:36;default:''" json:"departmentId"`
	NeedByDate        *time.Time `gorm:"column:need_by_date;type:date" json:"needByDate"`
	Priority          string    `gorm:"column:priority;size:20;default:''" json:"priority"`
	Reason            string    `gorm:"column:reason;type:text" json:"reason"`
	Status            string    `gorm:"column:status;not null;size:20;default:DRAFT;index:idx_req_status" json:"status"`
	TotalEstimated    float64   `gorm:"column:total_estimated;not null;default:0" json:"totalEstimated"`
	ApprovedBy        string    `gorm:"column:approved_by;size:36;default:''" json:"approvedBy"`
	ApprovedAt        *time.Time `gorm:"column:approved_at;type:timestamptz" json:"approvedAt"`
	RejectedReason    string    `gorm:"column:rejected_reason;type:text" json:"rejectedReason"`
	CreatedBy         string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (RequisitionGORM) TableName() string { return "purchase_requisitions" }

type RequisitionItemGORM struct {
	ID             string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	RequisitionID  string  `gorm:"column:requisition_id;not null;size:36;index:idx_reql_req" json:"requisitionId"`
	LineNumber     int     `gorm:"column:line_number;not null;default:0" json:"lineNumber"`
	ItemCode       string  `gorm:"column:item_code;size:50;default:''" json:"itemCode"`
	ItemName       string  `gorm:"column:item_name;not null;size:255" json:"itemName"`
	Unit           string  `gorm:"column:unit;size:20;default:''" json:"unit"`
	Quantity       float64 `gorm:"column:quantity;not null;default:0" json:"quantity"`
	EstimatedPrice float64 `gorm:"column:estimated_price;not null;default:0" json:"estimatedPrice"`
	EstimatedTotal float64 `gorm:"column:estimated_total;not null;default:0" json:"estimatedTotal"`
	AccountID      string  `gorm:"column:account_id;not null;size:36" json:"accountId"`
}

func (RequisitionItemGORM) TableName() string { return "requisition_lines" }

type FXRevaluationGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36;index:idx_fxr_company" json:"companyId"`
	RevaluationDate string    `gorm:"column:revaluation_date;not null;type:date" json:"revaluationDate"`
	Status          string    `gorm:"column:status;not null;size:20;default:DRAFT;index:idx_fxr_status" json:"status"`
	TotalGain       float64   `gorm:"column:total_gain;not null;default:0" json:"totalGain"`
	TotalLoss       float64   `gorm:"column:total_loss;not null;default:0" json:"totalLoss"`
	GLPosted        bool      `gorm:"column:gl_posted;not null;default:false" json:"glPosted"`
	GLPostedAt      string    `gorm:"column:gl_posted_at;default:''" json:"glPostedAt"`
	CreatedBy       string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (FXRevaluationGORM) TableName() string { return "fx_revaluations" }

type FXRevaluationLineGORM struct {
	ID              string  `gorm:"column:id;primaryKey;size:36" json:"id"`
	RevaluationID   string  `gorm:"column:revaluation_id;not null;size:36;index:idx_fxrl_revaluation" json:"revaluationId"`
	InvoiceID       string  `gorm:"column:invoice_id;not null;size:36;index:idx_fxrl_invoice" json:"invoiceId"`
	InvoiceNumber   string  `gorm:"column:invoice_number;not null;size:30" json:"invoiceNumber"`
	SupplierID      string  `gorm:"column:supplier_id;not null;size:36" json:"supplierId"`
	SupplierName    string  `gorm:"column:supplier_name;not null;size:255" json:"supplierName"`
	Currency        string  `gorm:"column:currency;not null;size:3" json:"currency"`
	BalanceDue      float64 `gorm:"column:balance_due;not null;default:0" json:"balanceDue"`
	OriginalRate    float64 `gorm:"column:original_rate;not null;default:0" json:"originalRate"`
	RevaluationRate float64 `gorm:"column:revaluation_rate;not null;default:0" json:"revaluationRate"`
	FxGain          float64 `gorm:"column:fx_gain;not null;default:0" json:"fxGain"`
	FxLoss          float64 `gorm:"column:fx_loss;not null;default:0" json:"fxLoss"`
}

func (FXRevaluationLineGORM) TableName() string { return "fx_revaluation_lines" }
