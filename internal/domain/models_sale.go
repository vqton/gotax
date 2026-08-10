package domain

import "time"

type CustomerStatus string
const (CustomerActive CustomerStatus="ACTIVE"; CustomerSuspended CustomerStatus="SUSPENDED"; CustomerBlacklisted CustomerStatus="BLACKLISTED")

type CustomerType string
const (CustomerDomestic CustomerType="domestic"; CustomerExport CustomerType="export"; CustomerBoth CustomerType="both")

type CustomerGroup string
const (CustGroupRetail CustomerGroup="retail"; CustGroupWholesale CustomerGroup="wholesale"; CustGroupDistributor CustomerGroup="distributor"; CustGroupAgent CustomerGroup="agent")

type SOStatus string
const (SODraft SOStatus="DRAFT"; SOApproved SOStatus="APPROVED"; SOConfirmed SOStatus="CONFIRMED"; SOProcessing SOStatus="PROCESSING"; SODelivered SOStatus="DELIVERED"; SOInvoiced SOStatus="INVOICED"; SOCancelled SOStatus="CANCELLED"; SOClosed SOStatus="CLOSED")
func (s SOStatus) ValidTransition(next SOStatus) bool {
	switch s {
	case SODraft: return next==SOApproved||next==SOCancelled
	case SOApproved: return next==SOConfirmed||next==SODraft||next==SOCancelled
	case SOConfirmed: return next==SOProcessing||next==SODelivered||next==SOCancelled
	case SOProcessing: return next==SODelivered||next==SOCancelled
	case SODelivered: return next==SOInvoiced||next==SOClosed||next==SOCancelled
	case SOInvoiced, SOClosed, SOCancelled: return false
	default: return false
	}
}

type DNStatus string
const (DNDraft DNStatus="DRAFT"; DNPosted DNStatus="POSTED"; DNCancelled DNStatus="CANCELLED")

const DefaultTolerancePct = 5.0

type SaleInvoiceStatus string
const (SInvDraft SaleInvoiceStatus="DRAFT"; SInvSigned SaleInvoiceStatus="SIGNED"; SInvSubmitted SaleInvoiceStatus="SUBMITTED"; SInvCoded SaleInvoiceStatus="CODED"; SInvIssued SaleInvoiceStatus="ISSUED"; SInvPosted SaleInvoiceStatus="POSTED"; SInvPaid SaleInvoiceStatus="PAID"; SInvCancelled SaleInvoiceStatus="CANCELLED"; SInvReplaced SaleInvoiceStatus="REPLACED")

type InvEStatus string
const (EInvPending InvEStatus="pending"; EInvSigned InvEStatus="signed"; EInvSubmitted InvEStatus="submitted"; EInvCoded InvEStatus="coded"; EInvIssued InvEStatus="issued"; EInvCancelled InvEStatus="cancelled"; EInvReplaced InvEStatus="replaced")

type CNStatus string
const (CNDraft CNStatus="DRAFT"; CNSigned CNStatus="SIGNED"; CNSubmitted CNStatus="SUBMITTED"; CNPosted CNStatus="POSTED"; CNCancelled CNStatus="CANCELLED")

type ReceiptStatus string
const (RcpDraft ReceiptStatus="DRAFT"; RcpPosted ReceiptStatus="POSTED"; RcpReconciled ReceiptStatus="RECONCILED"; RcpCancelled ReceiptStatus="CANCELLED")

type ARTransType string
const (ARTransInvoice ARTransType="invoice"; ARTransCreditNote ARTransType="credit_note"; ARTransReceipt ARTransType="receipt"; ARTransPrepayment ARTransType="prepayment"; ARTransOffset ARTransType="offset")

type SaleAdjustmentType string
const (SAdjIncrease SaleAdjustmentType="increase"; SAdjDecrease SaleAdjustmentType="decrease"; SAdjReplacement SaleAdjustmentType="replacement")

type ReturnType string
const (RetFull ReturnType="full_return"; RetPartial ReturnType="partial_return"; RetPriceAdj ReturnType="price_adjustment")

// ── Customer ──────────────────────────────────────────────────────────

type Customer struct {
	ID                string         `json:"id"`
	CompanyID         string         `json:"company_id"`
	Code              string         `json:"code"`
	Name              string         `json:"name"`
	TaxCode           string         `json:"tax_code"`
	Address           string         `json:"address,omitempty"`
	Phone             string         `json:"phone,omitempty"`
	Email             string         `json:"email,omitempty"`
	BankAccountName   string         `json:"bank_account_name,omitempty"`
	BankAccountNumber string         `json:"bank_account_number,omitempty"`
	BankName          string         `json:"bank_name,omitempty"`
	PaymentTerms      PaymentTerms   `json:"payment_terms,omitempty"`
	CreditLimit       float64        `json:"credit_limit"`
	Currency          string         `json:"currency"`
	CustomerType      CustomerType   `json:"customer_type,omitempty"`
	CustomerGroup     CustomerGroup  `json:"customer_group,omitempty"`
	PriceListID       string         `json:"price_list_id,omitempty"`
	Status            CustomerStatus `json:"status"`
	Notes             string         `json:"notes,omitempty"`
	CreatedBy         string         `json:"created_by,omitempty"`
	CreatedAt         time.Time      `json:"created_at,omitempty"`
	UpdatedAt         time.Time      `json:"updated_at,omitempty"`
}

func (c *Customer) Validate() error {
	if c.Code=="" { return ErrCustomerCodeRequired }
	if c.Name=="" { return ErrCustomerNameRequired }
	if c.TaxCode=="" { return ErrCustomerTaxCodeRequired }
	if c.Currency=="" { c.Currency="VND" }
	if c.Status=="" { c.Status=CustomerActive }
	switch c.Status {
	case CustomerActive,CustomerSuspended,CustomerBlacklisted:
	default: return ErrCustomerStatusInvalid
	}
	return nil
}

// ── Sales Order ───────────────────────────────────────────────────────

type SalesOrder struct {
	ID              string       `json:"id"`
	CompanyID       string       `json:"company_id"`
	SONumber        string       `json:"so_number"`
	QuotationID     string       `json:"quotation_id,omitempty"`
	CustomerID      string       `json:"customer_id"`
	OrderDate       time.Time    `json:"order_date"`
	ExpectedDate    *time.Time   `json:"expected_date,omitempty"`
	Currency        string       `json:"currency"`
	ExchangeRate    float64      `json:"exchange_rate"`
	PaymentTerms    string       `json:"payment_terms,omitempty"`
	DeliveryTerms   string       `json:"delivery_terms,omitempty"`
	ShippingAddress string       `json:"shipping_address,omitempty"`
	Subtotal        float64      `json:"subtotal"`
	DiscountAmount  float64      `json:"discount_amount"`
	TaxAmount       float64      `json:"tax_amount"`
	TotalAmount     float64      `json:"total_amount"`
	Status          SOStatus     `json:"status"`
	ApprovedBy      string       `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time   `json:"approved_at,omitempty"`
	CancelledReason string       `json:"cancelled_reason,omitempty"`
	Notes           string       `json:"notes,omitempty"`
	CreatedBy       string       `json:"created_by"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	UpdatedAt       time.Time    `json:"updated_at,omitempty"`
	Lines           []SOLine     `json:"lines"`
}

func (so *SalesOrder) Validate() error {
	if so.SONumber=="" { return ErrSONumberRequired }
	if so.CustomerID=="" { return ErrSOCustomerRequired }
	if so.OrderDate.IsZero() { return ErrSODateRequired }
	if so.Currency=="" { so.Currency="VND" }
	if so.Status=="" { so.Status=SODraft }
	switch so.Status {
	case SODraft,SOApproved,SOConfirmed,SOProcessing,SODelivered,SOInvoiced,SOCancelled,SOClosed:
	default: return ErrSOStatusInvalid
	}
	if len(so.Lines)==0 { return ErrSOLinesRequired }
	return nil
}

func (so *SalesOrder) CalculateTotals() {
	var subtotal,tax,discount float64
	for i:=range so.Lines {
		l:=&so.Lines[i]; raw:=l.Quantity*l.UnitPrice
		l.LineTotal=raw
		if l.DiscountPct>0 { l.LineTotal=raw-raw*l.DiscountPct/100; discount+=raw*l.DiscountPct/100 }
		l.LineVATAmount=l.LineTotal*l.VATRate/100
		subtotal+=l.LineTotal; tax+=l.LineVATAmount
	}
	so.Subtotal=subtotal; so.DiscountAmount=discount; so.TaxAmount=tax; so.TotalAmount=subtotal+tax
}

type SOLine struct {
	ID              string  `json:"id"`
	SOID            string  `json:"so_id"`
	LineNumber      int     `json:"line_number"`
	ItemCode        string  `json:"item_code,omitempty"`
	ItemName        string  `json:"item_name"`
	Unit            string  `json:"unit"`
	Quantity        float64 `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	DiscountPct     float64 `json:"discount_pct"`
	VATRate         float64 `json:"vat_rate"`
	VATType         VATType `json:"vat_type"`
	RevenueAccount  string  `json:"revenue_account_id"`
	VATAccountID    string  `json:"vat_account_id"`
	LineTotal       float64 `json:"line_total"`
	LineVATAmount   float64 `json:"line_vat_amount"`
	DeliveredQty    float64 `json:"delivered_qty"`
	InvoicedQty     float64 `json:"invoiced_qty"`
}

func (l *SOLine) Validate() error {
	if l.ItemName=="" { return ErrSOItemNameRequired }
	if l.Unit=="" { return ErrSOItemUnitRequired }
	if l.Quantity<=0 { return ErrSOItemQtyRequired }
	if l.UnitPrice<0 { return ErrSOItemPriceRequired }
	if l.RevenueAccount=="" { return ErrSOItemAccountRequired }
	if l.VATAccountID=="" { return ErrSOItemVATAccountRequired }
	return nil
}

// ── Delivery Note ─────────────────────────────────────────────────────

type DeliveryNote struct {
	ID              string       `json:"id"`
	CompanyID       string       `json:"company_id"`
	DNNumber        string       `json:"dn_number"`
	SOID            string       `json:"so_id"`
	DeliveryDate    time.Time    `json:"delivery_date"`
	Warehouse       string       `json:"warehouse,omitempty"`
	ShippingMethod  string       `json:"shipping_method,omitempty"`
	CarrierName     string       `json:"carrier_name,omitempty"`
	TrackingNumber  string       `json:"tracking_number,omitempty"`
	DeliveryAddress string       `json:"delivery_address,omitempty"`
	Status           DNStatus     `json:"status"`
	Notes            string       `json:"notes,omitempty"`
	TolerancePercent float64      `json:"tolerance_percent,omitempty"`
	GLPosted         bool         `json:"gl_posted"`
	GLPostedAt       *time.Time   `json:"gl_posted_at,omitempty"`
	CreatedBy        string       `json:"created_by"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	UpdatedAt       time.Time    `json:"updated_at,omitempty"`
	Lines           []DNLine     `json:"lines"`
}

func (d *DeliveryNote) Validate() error {
	if d.DNNumber=="" { return ErrDNNumberRequired }
	if d.SOID=="" { return ErrDNSORequired }
	if d.DeliveryDate.IsZero() { return ErrDNDateRequired }
	if d.Status=="" { d.Status=DNDraft }
	switch d.Status {
	case DNDraft,DNPosted,DNCancelled:
	default: return ErrDNStatusInvalid
	}
	if len(d.Lines)==0 { return ErrDNLinesRequired }
	return nil
}

type DNLine struct {
	ID                string  `json:"id"`
	DNID              string  `json:"dn_id"`
	SOLineID          string  `json:"so_line_id"`
	ItemCode          string  `json:"item_code,omitempty"`
	ItemName          string  `json:"item_name"`
	Unit              string  `json:"unit"`
	QtyDelivered      float64 `json:"quantity_delivered"`
	QtyReturned       float64 `json:"quantity_returned"`
	UnitPrice         float64 `json:"unit_price"`
	LineTotal         float64 `json:"line_total"`
	CostPrice         float64 `json:"cost_price,omitempty"`
	LotNumber         string  `json:"lot_number,omitempty"`
	ExpiryDate        string  `json:"expiry_date,omitempty"`
}

func (l *DNLine) Validate() error {
	if l.ItemName=="" { return ErrDNItemNameRequired }
	if l.SOLineID=="" { return ErrDNItemSOLineRequired }
	if l.QtyDelivered<=0 { return ErrDNItemQtyRequired }
	return nil
}

// ── Customer Invoice ──────────────────────────────────────────────────

type CustomerInvoice struct {
	ID                  string         `json:"id"`
	CompanyID           string         `json:"company_id"`
	InvoiceNumber       string         `json:"invoice_number"`
	InvoiceDate         time.Time      `json:"invoice_date"`
	SOID                string         `json:"so_id,omitempty"`
	DNID                string         `json:"dn_id,omitempty"`
	CustomerID          string         `json:"customer_id"`
	CustomerName        string         `json:"customer_name"`
	CustomerTaxCode     string         `json:"customer_tax_code"`
	CustomerAddress     string         `json:"customer_address"`
	InvoiceType         string         `json:"invoice_type"`
	Currency            string         `json:"currency"`
	ExchangeRate        float64        `json:"exchange_rate"`
	Subtotal            float64        `json:"subtotal"`
	DiscountAmount      float64        `json:"discount_amount"`
	TaxAmount           float64        `json:"tax_amount"`
	TotalAmount         float64        `json:"total_amount"`
	AmountReceived      float64        `json:"amount_received"`
	BalanceDue          float64        `json:"balance_due"`
	DueDate             *time.Time     `json:"due_date,omitempty"`
	InvoiceNote         string         `json:"invoice_note,omitempty"`
	EInvoiceData        string         `json:"e_invoice_data,omitempty"`
	EInvoiceCode        string         `json:"e_invoice_code,omitempty"`
	EInvStatus          InvEStatus     `json:"e_invoice_status"`
	DigitalSignatureID  string         `json:"digital_signature_id,omitempty"`
	SignedData          string         `json:"signed_data,omitempty"`
	GDTResponse         string         `json:"gdt_response,omitempty"`
	OriginalInvoiceID   string         `json:"original_invoice_id,omitempty"`
	AdjustmentType      SaleAdjustmentType `json:"adjustment_type,omitempty"`
	Status              SaleInvoiceStatus  `json:"status"`
	GLPosted            bool           `json:"gl_posted"`
	GLPostedAt          *time.Time     `json:"gl_posted_at,omitempty"`
	Notes           string       `json:"notes,omitempty"`
	CreatedBy       string       `json:"created_by"`
	CreatedAt           time.Time      `json:"created_at,omitempty"`
	UpdatedAt           time.Time      `json:"updated_at,omitempty"`
	Lines               []InvLine      `json:"lines"`
}

func (inv *CustomerInvoice) Validate() error {
	if inv.InvoiceNumber=="" { return ErrInvNumberRequired }
	if inv.CustomerID=="" { return ErrInvCustomerRequired }
	if inv.CustomerName=="" { return ErrInvCustomerNameRequired }
	if inv.CustomerTaxCode=="" { return ErrInvCustomerTaxCodeRequired }
	if inv.InvoiceDate.IsZero() { return ErrInvDateRequired }
	if inv.Currency=="" { inv.Currency="VND" }
	if inv.Status=="" { inv.Status=SInvDraft }
	switch inv.Status {
	case SInvDraft,SInvSigned,SInvSubmitted,SInvCoded,SInvIssued,SInvPosted,SInvPaid,SInvCancelled,SInvReplaced:
	default: return ErrInvStatusInvalid
	}
	if len(inv.Lines)==0 { return ErrInvLinesRequired }
	return nil
}

func (inv *CustomerInvoice) CalculateTotals() {
	var subtotal,tax float64
	for i:=range inv.Lines {
		l:=&inv.Lines[i]; raw:=l.Quantity*l.UnitPrice; l.LineTotal=raw
		if l.DiscountPct>0 { l.LineTotal=raw-raw*l.DiscountPct/100 }
		l.LineVATAmount=l.LineTotal*l.VATRate/100
		subtotal+=l.LineTotal; tax+=l.LineVATAmount
	}
	inv.Subtotal=subtotal; inv.TaxAmount=tax; inv.TotalAmount=subtotal+tax; inv.BalanceDue=subtotal+tax-inv.AmountReceived
}

type InvLine struct {
	ID             string  `json:"id"`
	InvoiceID      string  `json:"invoice_id"`
	SOLineID       string  `json:"so_line_id,omitempty"`
	DNLineID       string  `json:"dn_line_id,omitempty"`
	ItemCode       string  `json:"item_code,omitempty"`
	ItemName       string  `json:"item_name"`
	Unit           string  `json:"unit"`
	Quantity       float64 `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	DiscountPct    float64 `json:"discount_pct"`
	VATRate        float64 `json:"vat_rate"`
	VATType        VATType `json:"vat_type"`
	LineTotal      float64 `json:"line_total"`
	LineVATAmount  float64 `json:"line_vat_amount"`
	RevenueAccount string  `json:"revenue_account_id"`
	VATAccountID   string  `json:"vat_account_id"`
}

func (l *InvLine) Validate() error {
	if l.ItemName=="" { return ErrInvItemNameRequired }
	if l.Quantity<=0 { return ErrInvItemQtyRequired }
	if l.UnitPrice<0 { return ErrInvItemPriceRequired }
	if l.RevenueAccount=="" { return ErrInvItemAccountRequired }
	if l.VATAccountID=="" { return ErrInvItemVATAccountRequired }
	return nil
}

// ── Customer Receipt ──────────────────────────────────────────────────

type CustomerReceipt struct {
	ID                string         `json:"id"`
	CompanyID         string         `json:"company_id"`
	ReceiptNumber     string         `json:"receipt_number"`
	CustomerID        string         `json:"customer_id"`
	ReceiptDate       time.Time      `json:"receipt_date"`
	PaymentMethod     string         `json:"payment_method"`
	BankAccountID     string         `json:"bank_account_id,omitempty"`
	Currency          string         `json:"currency"`
	ExchangeRate      float64        `json:"exchange_rate"`
	Amount            float64        `json:"amount"`
	UnallocatedAmount float64        `json:"unallocated_amount"`
	Reference         string         `json:"reference,omitempty"`
	Notes             string         `json:"notes,omitempty"`
	Status            ReceiptStatus  `json:"status"`
	GLPosted          bool           `json:"gl_posted"`
	GLPostedAt        *time.Time     `json:"gl_posted_at,omitempty"`
	CreatedBy         string         `json:"created_by"`
	CreatedAt         time.Time      `json:"created_at,omitempty"`
	UpdatedAt         time.Time      `json:"updated_at,omitempty"`
	Allocations       []RcpAllocation `json:"allocations"`
}

func (r *CustomerReceipt) Validate() error {
	if r.ReceiptNumber=="" { return ErrRcpNumberRequired }
	if r.CustomerID=="" { return ErrRcpCustomerRequired }
	if r.ReceiptDate.IsZero() { return ErrRcpDateRequired }
	if r.Amount<=0 { return ErrRcpAmountRequired }
	if r.Currency=="" { r.Currency="VND" }
	if r.Status=="" { r.Status=RcpDraft }
	switch r.Status {
	case RcpDraft,RcpPosted,RcpReconciled,RcpCancelled:
	default: return ErrRcpStatusInvalid
	}
	return nil
}

type RcpAllocation struct {
	ID              string  `json:"id"`
	ReceiptID       string  `json:"receipt_id"`
	InvoiceID       string  `json:"invoice_id"`
	AllocatedAmount float64 `json:"allocated_amount"`
	DiscountAmount  float64 `json:"discount_amount,omitempty"`
}

// ── Credit Note ───────────────────────────────────────────────────────

type CreditNote struct {
	ID                string         `json:"id"`
	CompanyID         string         `json:"company_id"`
	CNNumber          string         `json:"cn_number"`
	OriginalInvoiceID string         `json:"original_invoice_id"`
	CustomerID        string         `json:"customer_id"`
	ReturnDate        time.Time      `json:"return_date"`
	ReturnReason      string         `json:"return_reason"`
	ReturnType        ReturnType     `json:"return_type"`
	DNID              string         `json:"dn_id,omitempty"`
	Subtotal          float64        `json:"subtotal"`
	TaxAmount         float64        `json:"tax_amount"`
	TotalAmount       float64        `json:"total_amount"`
	EInvoiceData      string         `json:"e_invoice_data,omitempty"`
	EInvoiceCode      string         `json:"e_invoice_code,omitempty"`
	Status            CNStatus       `json:"status"`
	GLPosted          bool           `json:"gl_posted"`
	GLPostedAt        *time.Time     `json:"gl_posted_at,omitempty"`
	Notes             string         `json:"notes,omitempty"`
	CreatedBy         string         `json:"created_by"`
	CreatedAt         time.Time      `json:"created_at,omitempty"`
	UpdatedAt         time.Time      `json:"updated_at,omitempty"`
	Lines             []CNLine       `json:"lines"`
}

func (cn *CreditNote) Validate() error {
	if cn.CNNumber=="" { return ErrCNNumberRequired }
	if cn.OriginalInvoiceID=="" { return ErrCNOriginalInvRequired }
	if cn.CustomerID=="" { return ErrCNCustomerRequired }
	if cn.ReturnDate.IsZero() { return ErrCNDateRequired }
	if cn.Status=="" { cn.Status=CNDraft }
	switch cn.Status {
	case CNDraft,CNSigned,CNSubmitted,CNPosted,CNCancelled:
	default: return ErrCNStatusInvalid
	}
	if len(cn.Lines)==0 { return ErrCNLinesRequired }
	return nil
}

func (cn *CreditNote) CalculateTotals() {
	var subtotal,tax float64
	for i:=range cn.Lines {
		l:=&cn.Lines[i]; raw:=l.Quantity*l.UnitPrice; l.LineTotal=raw
		l.LineVATAmount=l.LineTotal*l.VATRate/100
		subtotal+=l.LineTotal; tax+=l.LineVATAmount
	}
	cn.Subtotal=subtotal; cn.TaxAmount=tax; cn.TotalAmount=subtotal+tax
}

type CNLine struct {
	ID             string  `json:"id"`
	CNID           string  `json:"cn_id"`
	InvLineID      string  `json:"invoice_line_id"`
	ItemName       string  `json:"item_name"`
	Unit           string  `json:"unit"`
	Quantity       float64 `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	VATRate        float64 `json:"vat_rate"`
	LineTotal      float64 `json:"line_total"`
	LineVATAmount  float64 `json:"line_vat_amount"`
}

// ── AR Transaction ────────────────────────────────────────────────────

type ARTransaction struct {
	ID              string      `json:"id"`
	CompanyID       string      `json:"company_id"`
	CustomerID      string      `json:"customer_id"`
	InvoiceID       string      `json:"invoice_id,omitempty"`
	TransactionType ARTransType `json:"transaction_type"`
	TransactionDate time.Time   `json:"transaction_date"`
	Amount          float64     `json:"amount"`
	Currency        string      `json:"currency"`
	ReferenceType   string      `json:"reference_type,omitempty"`
	ReferenceID     string      `json:"reference_id,omitempty"`
	Notes           string      `json:"notes,omitempty"`
	CreatedAt       time.Time   `json:"created_at,omitempty"`
}

func (t *ARTransaction) Validate() error {
	if t.CustomerID=="" { return ErrARTransCustomerRequired }
	if t.TransactionDate.IsZero() { return ErrARTransDateRequired }
	if t.Amount==0 { return ErrARTransAmountRequired }
	if t.Currency=="" { t.Currency="VND" }
	switch t.TransactionType {
	case ARTransInvoice,ARTransCreditNote,ARTransReceipt,ARTransPrepayment,ARTransOffset:
	default: return ErrARTransTypeInvalid
	}
	return nil
}

// ── Sales Quotation (P1) ──────────────────────────────────────────────

type SalesQuotation struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	QNNumber    string    `json:"qn_number"`
	CustomerID  string    `json:"customer_id"`
	ValidUntil  *time.Time `json:"valid_until,omitempty"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// ── Filters ───────────────────────────────────────────────────────────

type SalesOrderFilter struct {
	CompanyID  string
	CustomerID string
	Status     SOStatus
	FromDate   string
	ToDate     string
	Offset     int
	Limit      int
}

type DeliveryNoteFilter struct {
	CompanyID string
	SOID      string
	Status    DNStatus
	FromDate  string
	ToDate    string
	Offset    int
	Limit     int
}

type CustomerInvoiceFilter struct {
	CompanyID  string
	CustomerID string
	Status     SaleInvoiceStatus
	FromDate   string
	ToDate     string
	Offset     int
	Limit      int
}

type ReceiptFilter struct {
	CompanyID  string
	CustomerID string
	Status     ReceiptStatus
	FromDate   string
	ToDate     string
	Offset     int
	Limit      int
}

type CreditNoteFilter struct {
	CompanyID  string
	CustomerID string
	Status     CNStatus
	FromDate   string
	ToDate     string
	Offset     int
	Limit      int
}

// ── Report Structs ────────────────────────────────────────────────────

type ARAgingBucket struct {
	Bucket0  float64 `json:"bucket_0"`
	Bucket30 float64 `json:"bucket_30"`
	Bucket60 float64 `json:"bucket_60"`
	Bucket90 float64 `json:"bucket_90"`
	Bucket120 float64 `json:"bucket_120"`
	Total    float64 `json:"total"`
}

type ARAgingReport struct {
	CustomerID   string        `json:"customer_id"`
	CustomerName string        `json:"customer_name"`
	TaxCode      string        `json:"tax_code"`
	Buckets      ARAgingBucket `json:"buckets"`
}

type ARSummary struct {
	CustomerID    string  `json:"customer_id"`
	CustomerName  string  `json:"customer_name"`
	TaxCode       string  `json:"tax_code"`
	TotalInvoiced float64 `json:"total_invoiced"`
	TotalReceived float64 `json:"total_received"`
	Outstanding   float64 `json:"outstanding"`
	Currency      string  `json:"currency"`
}

type CustomerStatementLine struct {
	Date        time.Time `json:"date"`
	RefType     string    `json:"ref_type"`
	RefNumber   string    `json:"ref_number"`
	Description string    `json:"description"`
	Debit       float64   `json:"debit"`
	Credit      float64   `json:"credit"`
	Balance     float64   `json:"balance"`
}

type CustomerStatement struct {
	Customer     Customer                `json:"customer"`
	FromDate     string                  `json:"from_date"`
	ToDate       string                  `json:"to_date"`
	OpeningBal   float64                 `json:"opening_balance"`
	ClosingBal   float64                 `json:"closing_balance"`
	Lines        []CustomerStatementLine `json:"lines"`
}

type ARGLReconciliation struct {
	PeriodID       string                  `json:"period_id"`
	PeriodLabel    string                  `json:"period_label"`
	SubledgerTotal float64                 `json:"subledger_total"`
	GLBalance      float64                 `json:"gl_balance"`
	Variance       float64                 `json:"variance"`
	Details        []ARGLReconDetail       `json:"details,omitempty"`
}

type ARGLReconDetail struct {
	CustomerID     string  `json:"customer_id"`
	CustomerName   string  `json:"customer_name"`
	BalanceDue     float64 `json:"balance_due"`
}

// ── FR-9 Reports ──────────────────────────────────────────────────────

// SalesLedgerRow — S01-BH sales ledger line (per customer per period).
type SalesLedgerRow struct {
	Date        string  `json:"date"`
	Ref         string  `json:"ref"`
	Description string  `json:"description"`
	Revenue     float64 `json:"revenue"` // 5111
	VAT         float64 `json:"vat"`     // 3331
	Total       float64 `json:"total"`
}

type SalesLedgerReport struct {
	CompanyID    string           `json:"company_id"`
	CustomerID   string           `json:"customer_id"`
	CustomerName string           `json:"customer_name"`
	FromDate     string           `json:"from_date"`
	ToDate       string           `json:"to_date"`
	Rows         []SalesLedgerRow `json:"rows"`
	TotalRevenue float64          `json:"total_revenue"`
	TotalVAT     float64          `json:"total_vat"`
	Total        float64          `json:"total"`
}

// CustomerLedgerRow — S02-BH customer detail ledger line (131).
type CustomerLedgerRow struct {
	Date        string  `json:"date"`
	Ref         string  `json:"ref"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Balance     float64 `json:"balance"`
}

type CustomerLedgerReport struct {
	CompanyID      string              `json:"company_id"`
	CustomerID     string              `json:"customer_id"`
	CustomerName   string              `json:"customer_name"`
	FromDate       string              `json:"from_date"`
	ToDate         string              `json:"to_date"`
	OpeningBalance float64             `json:"opening_balance"`
	Rows           []CustomerLedgerRow `json:"rows"`
	ClosingBalance float64             `json:"closing_balance"`
}

// GoodsLedgerRow — S03-BH goods sales ledger line (per item per period).
type GoodsLedgerRow struct {
	Date      string  `json:"date"`
	Ref       string  `json:"ref"`
	ItemCode  string  `json:"item_code,omitempty"`
	ItemName  string  `json:"item_name"`
	Unit      string  `json:"unit"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Revenue   float64 `json:"revenue"`
	VAT       float64 `json:"vat"`
	Total     float64 `json:"total"`
}

type GoodsLedgerReport struct {
	CompanyID   string          `json:"company_id"`
	FromDate    string          `json:"from_date"`
	ToDate      string          `json:"to_date"`
	Rows        []GoodsLedgerRow `json:"rows"`
	TotalQty    float64         `json:"total_quantity"`
	TotalRevenue float64        `json:"total_revenue"`
	TotalVAT    float64         `json:"total_vat"`
}

// VATOutputRow — VAT output tracking per rate (3331).
type VATOutputRow struct {
	VatRate      float64 `json:"vat_rate"`
	Subtotal     float64 `json:"subtotal"`
	VatAmount    float64 `json:"vat_amount"`
	InvoiceCount int     `json:"invoice_count"`
}

type VATOutputReport struct {
	CompanyID     string         `json:"company_id"`
	FromDate      string         `json:"from_date"`
	ToDate        string         `json:"to_date"`
	Rows          []VATOutputRow `json:"rows"`
	TotalSubtotal float64        `json:"total_subtotal"`
	TotalVAT      float64        `json:"total_vat"`
}

// UnbilledDeliveryRow — posted delivery with no covering invoice.
type UnbilledDeliveryRow struct {
	DNID         string  `json:"dn_id"`
	DNNumber     string  `json:"dn_number"`
	DeliveryDate string  `json:"delivery_date"`
	SOID         string  `json:"so_id"`
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	ItemCode     string  `json:"item_code,omitempty"`
	ItemName     string  `json:"item_name"`
	QtyDelivered float64 `json:"quantity_delivered"`
	Amount       float64 `json:"amount"`
}

type UnbilledDeliveryReport struct {
	CompanyID string                `json:"company_id"`
	Rows      []UnbilledDeliveryRow `json:"rows"`
}

// ─── Price Lists ─────────────────────────────────────────────────────

type PriceList struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Currency    string    `json:"currency"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Lines       []PriceListLine `json:"lines,omitempty"`
}

func (p *PriceList) Validate() error {
	if p.Code == "" { return ErrPriceListCodeRequired }
	if p.Name == "" { return ErrPriceListNameRequired }
	if p.Currency == "" { p.Currency = "VND" }
	return nil
}

type PriceListLine struct {
	ID           string  `json:"id"`
	PriceListID  string  `json:"price_list_id"`
	ItemCode     string  `json:"item_code"`
	ItemName     string  `json:"item_name"`
	Unit         string  `json:"unit"`
	UnitPrice    float64 `json:"unit_price"`
	VATRate      float64 `json:"vat_rate"`
	MinQuantity  float64 `json:"min_quantity,omitempty"`
	EffectiveFrom string `json:"effective_from,omitempty"`
	EffectiveTo   string `json:"effective_to,omitempty"`
}

func (l *PriceListLine) Validate() error {
	if l.ItemCode == "" { return ErrPriceListItemCodeRequired }
	if l.UnitPrice < 0 { return ErrPriceListInvalidPrice }
	return nil
}
