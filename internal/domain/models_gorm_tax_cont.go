package domain

import "time"

type EInvoiceGORM struct {
	ID                  string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID           string    `gorm:"column:company_id;not null;size:36;index:idx_einv_company,unique" json:"companyId"`
	Pattern             string    `gorm:"column:pattern;not null;size:10;index:idx_einv_company,unique" json:"pattern"`
	Serial              string    `gorm:"column:serial;not null;size:20" json:"serial"`
	InvoiceNumber       int       `gorm:"column:invoice_number;default:0" json:"invoiceNumber"`
	InvoiceType         string    `gorm:"column:invoice_type;not null;size:10" json:"invoiceType"`
	BuyerName           string    `gorm:"column:buyer_name;not null;size:255" json:"buyerName"`
	BuyerTaxCode        *string   `gorm:"column:buyer_tax_code;size:20" json:"buyerTaxCode"`
	BuyerAddress        *string   `gorm:"column:buyer_address;type:text" json:"buyerAddress"`
	BuyerEmail          *string   `gorm:"column:buyer_email;size:255" json:"buyerEmail"`
	CurrencyCode        string    `gorm:"column:currency_code;not null;size:3;default:VND" json:"currencyCode"`
	ExchangeRate        float64   `gorm:"column:exchange_rate;default:1" json:"exchangeRate"`
	Subtotal            float64   `gorm:"column:subtotal;not null" json:"subtotal"`
	VATAmount           float64   `gorm:"column:vat_amount;not null" json:"vatAmount"`
	GrandTotal          float64   `gorm:"column:grand_total;not null" json:"grandTotal"`
	IssueDate           time.Time `gorm:"column:issue_date;not null;type:date" json:"issueDate"`
	Status              string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	XMLBody             *string   `gorm:"column:xml_body;type:text" json:"-"`
	SignedXML           *string   `gorm:"column:signed_xml;type:text" json:"-"`
	SigningDate         *string   `gorm:"column:signing_date" json:"signingDate"`
	DigitalSignatureID  *string   `gorm:"column:digital_signature_id;size:36" json:"digitalSignatureId"`
	JournalEntryID      *string   `gorm:"column:journal_entry_id;size:36;index" json:"journalEntryId"`
	CancelledAt         *string   `gorm:"column:cancelled_at" json:"cancelledAt"`
	CancelReason        *string   `gorm:"column:cancel_reason;type:text" json:"cancelReason"`
	OriginalInvoiceID   *string   `gorm:"column:original_invoice_id;size:36" json:"originalInvoiceId"`
	GDTResponse         *string   `gorm:"column:gdt_response;type:text" json:"-"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	Lines               []EInvoiceLineGORM `gorm:"foreignKey:EInvoiceID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (EInvoiceGORM) TableName() string { return "e_invoices" }

type EInvoiceLineGORM struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EInvoiceID  string    `gorm:"column:e_invoice_id;not null;size:36;index:idx_einvline_inv,unique" json:"eInvoiceId"`
	LineNumber  int       `gorm:"column:line_number;not null;index:idx_einvline_inv,unique" json:"lineNumber"`
	Description string    `gorm:"column:description;not null;type:text" json:"description"`
	Unit        *string   `gorm:"column:unit;size:20" json:"unit"`
	Quantity    float64   `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice   float64   `gorm:"column:unit_price;not null" json:"unitPrice"`
	LineTotal   float64   `gorm:"column:line_total;not null" json:"lineTotal"`
	VATRate     float64   `gorm:"column:vat_rate;not null" json:"vatRate"`
	VATAmount   float64   `gorm:"column:vat_amount;not null" json:"vatAmount"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (EInvoiceLineGORM) TableName() string { return "e_invoice_lines" }

type TaxCalendarGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID    string    `gorm:"column:company_id;not null;size:36;index:idx_taxcal_company_period,unique" json:"companyId"`
	TaxType      string    `gorm:"column:tax_type;not null;size:20;index:idx_taxcal_company_period,unique" json:"taxType"`
	PeriodYear   int       `gorm:"column:period_year;not null;index:idx_taxcal_company_period,unique" json:"periodYear"`
	PeriodNumber int       `gorm:"column:period_number;not null;index:idx_taxcal_company_period,unique" json:"periodNumber"`
	DueDate      time.Time `gorm:"column:due_date;not null;type:date" json:"dueDate"`
	Status       string    `gorm:"column:status;not null;size:20;default:'PENDING';index" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (TaxCalendarGORM) TableName() string { return "tax_calendar" }

type TaxAlertGORM struct {
	ID         string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID  string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	AlertType  string    `gorm:"column:alert_type;not null;size:30;index" json:"alertType"`
	Message    string    `gorm:"column:message;not null;type:text" json:"message"`
	IsRead     bool      `gorm:"column:is_read;default:false;index" json:"isRead"`
	DueDate    *time.Time `gorm:"column:due_date;type:date" json:"dueDate"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (TaxAlertGORM) TableName() string { return "tax_alerts" }

type TaxAuditCaseGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	CaseNumber  string    `gorm:"column:case_number;not null;size:30;uniqueIndex" json:"caseNumber"`
	AuditType   string    `gorm:"column:audit_type;not null;size:30" json:"auditType"`
	Status      string    `gorm:"column:status;not null;size:20;default:'OPEN';index" json:"status"`
	OpenDate    time.Time `gorm:"column:open_date;not null;type:date" json:"openDate"`
	CloseDate   *time.Time `gorm:"column:close_date;type:date" json:"closeDate"`
	AuditorName *string   `gorm:"column:auditor_name;size:255" json:"auditorName"`
	Notes       *string   `gorm:"column:notes;type:text" json:"notes"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (TaxAuditCaseGORM) TableName() string { return "tax_audit_cases" }
