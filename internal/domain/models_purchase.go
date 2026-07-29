package domain

import "time"

type SupplierStatus string
const (SupplierActive SupplierStatus="ACTIVE"; SupplierSuspended SupplierStatus="SUSPENDED"; SupplierBlacklisted SupplierStatus="BLACKLISTED")

type SupplierType string
const (SupplierTypeDomestic SupplierType="domestic"; SupplierTypeImport SupplierType="import"; SupplierTypeBoth SupplierType="both")

type PaymentTerms string
const (PaymentNet30 PaymentTerms="net30"; PaymentNet60 PaymentTerms="net60"; PaymentNet90 PaymentTerms="net90"; PaymentCOD PaymentTerms="COD")

type Supplier struct {
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
	SupplierType      SupplierType   `json:"supplier_type,omitempty"`
	Status            SupplierStatus `json:"status"`
	Notes             string         `json:"notes,omitempty"`
	CreatedAt         time.Time      `json:"created_at,omitempty"`
	UpdatedAt         time.Time      `json:"updated_at,omitempty"`
}

func (s *Supplier) Validate() error {
	if s.Code == "" { return ErrSupplierCodeRequired }
	if s.Name == "" { return ErrSupplierNameRequired }
	if s.TaxCode == "" { return ErrSupplierTaxCodeRequired }
	if s.Currency == "" { s.Currency = "VND" }
	if s.Status == "" { s.Status = SupplierActive }
	switch s.Status {
	case SupplierActive, SupplierSuspended, SupplierBlacklisted:
	default: return ErrSupplierStatusInvalid
	}
	return nil
}

type POStatus string
const (POStatusDraft POStatus="DRAFT"; POStatusApproved POStatus="APPROVED"; POStatusSent POStatus="SENT"; POStatusPartial POStatus="PARTIAL"; POStatusReceived POStatus="RECEIVED"; POStatusCancelled POStatus="CANCELLED"; POStatusClosed POStatus="CLOSED")

func (s POStatus) ValidTransition(next POStatus) bool {
	switch s {
	case POStatusDraft: return next == POStatusApproved || next == POStatusCancelled
	case POStatusApproved: return next == POStatusSent || next == POStatusDraft || next == POStatusCancelled
	case POStatusSent: return next == POStatusPartial || next == POStatusReceived || next == POStatusClosed || next == POStatusCancelled
	case POStatusPartial: return next == POStatusReceived || next == POStatusClosed || next == POStatusCancelled
	case POStatusReceived, POStatusClosed, POStatusCancelled: return false
	default: return false
	}
}

type VATType string
const (VAT0 VATType="VAT_0"; VAT5 VATType="VAT_5"; VAT8 VATType="VAT_8"; VAT10 VATType="VAT_10"; VATNonTaxable VATType="NT")

type PurchaseOrder struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	PONumber        string    `json:"po_number"`
	SupplierID      string    `json:"supplier_id"`
	RequisitionID   string    `json:"requisition_id,omitempty"`
	OrderDate       time.Time `json:"order_date"`
	ExpectedDate    *time.Time `json:"expected_date,omitempty"`
	Currency        string    `json:"currency"`
	ExchangeRate    float64   `json:"exchange_rate"`
	PaymentTerms    string    `json:"payment_terms,omitempty"`
	DeliveryTerms   string    `json:"delivery_terms,omitempty"`
	Subtotal        float64   `json:"subtotal"`
	DiscountAmount  float64   `json:"discount_amount"`
	TaxAmount       float64   `json:"tax_amount"`
	TotalAmount     float64   `json:"total_amount"`
	Status          POStatus  `json:"status"`
	ApprovedBy      string    `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	CancelledReason string    `json:"cancelled_reason,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	Lines           []POItem  `json:"lines"`
}

func (po *PurchaseOrder) Validate() error {
	if po.PONumber == "" { return ErrPONumberRequired }
	if po.SupplierID == "" { return ErrPOSupplierRequired }
	if po.OrderDate.IsZero() { return ErrPODateRequired }
	if po.Currency == "" { po.Currency = "VND" }
	if po.Status == "" { po.Status = POStatusDraft }
	switch po.Status {
	case POStatusDraft, POStatusApproved, POStatusSent, POStatusPartial, POStatusReceived, POStatusCancelled, POStatusClosed:
	default: return ErrPOStatusInvalid
	}
	if len(po.Lines) == 0 { return ErrPOLinesRequired }
	return nil
}

func (po *PurchaseOrder) CalculateTotals() {
	var subtotal, discount, tax float64
	for i := range po.Lines {
		l := &po.Lines[i]
		l.LineTotal = l.Quantity * l.UnitPrice
		if l.DiscountPct > 0 {
			l.LineTotal -= l.LineTotal * l.DiscountPct / 100
		}
		l.LineVATAmount = l.LineTotal * l.VATRate / 100
		subtotal += l.LineTotal
		discount += l.LineTotal * l.DiscountPct / 100
		tax += l.LineVATAmount
	}
	po.Subtotal = subtotal
	po.DiscountAmount = discount
	po.TaxAmount = tax
	po.TotalAmount = subtotal + tax
}

type POItem struct {
	ID           string  `json:"id"`
	POID         string  `json:"po_id"`
	LineNumber   int     `json:"line_number"`
	ItemCode     string  `json:"item_code,omitempty"`
	ItemName     string  `json:"item_name"`
	Unit         string  `json:"unit"`
	Quantity     float64 `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	DiscountPct  float64 `json:"discount_pct"`
	VATRate      float64 `json:"vat_rate"`
	VATType      VATType `json:"vat_type"`
	AccountID    string  `json:"account_id"`
	VATAccountID string  `json:"vat_account_id"`
	LineTotal    float64 `json:"line_total"`
	LineVATAmount float64 `json:"line_vat_amount"`
	ReceivedQty  float64 `json:"received_qty"`
	InvoicedQty  float64 `json:"invoiced_qty"`
}

func (l *POItem) Validate() error {
	if l.ItemName == "" { return ErrPOItemNameRequired }
	if l.Unit == "" { return ErrPOItemUnitRequired }
	if l.Quantity <= 0 { return ErrPOItemQuantityRequired }
	if l.UnitPrice < 0 { return ErrPOItemPriceRequired }
	if l.AccountID == "" { return ErrPOItemAccountRequired }
	if l.VATAccountID == "" { return ErrPOItemVATAccountRequired }
	return nil
}

type GRNStatus string
const (GRNDraft GRNStatus="DRAFT"; GRNPosted GRNStatus="POSTED"; GRNCancelled GRNStatus="CANCELLED")

type GRN struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	GRNNumber   string     `json:"grn_number"`
	POID        string     `json:"po_id"`
	ReceiptDate time.Time  `json:"receipt_date"`
	Warehouse   string     `json:"warehouse,omitempty"`
	Status      GRNStatus  `json:"status"`
	Notes       string     `json:"notes,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	Lines       []GRNItem  `json:"lines"`
}

func (g *GRN) Validate() error {
	if g.GRNNumber == "" { return ErrGRNNumberRequired }
	if g.POID == "" { return ErrGRNPORequired }
	if g.ReceiptDate.IsZero() { return ErrGRNDateRequired }
	if g.Status == "" { g.Status = GRNDraft }
	switch g.Status {
	case GRNDraft, GRNPosted, GRNCancelled:
	default: return ErrGRNStatusInvalid
	}
	if len(g.Lines) == 0 { return ErrGRNLinesRequired }
	return nil
}

type GRNItem struct {
	ID               string  `json:"id"`
	GRNID            string  `json:"grn_id"`
	POLineID         string  `json:"po_line_id"`
	ItemCode         string  `json:"item_code,omitempty"`
	ItemName         string  `json:"item_name"`
	Unit             string  `json:"unit"`
	QuantityReceived float64 `json:"quantity_received"`
	QuantityRejected float64 `json:"quantity_rejected"`
	UnitPrice        float64 `json:"unit_price"`
	LineTotal        float64 `json:"line_total"`
}

func (l *GRNItem) Validate() error {
	if l.ItemName == "" { return ErrGRNItemNameRequired }
	if l.POLineID == "" { return ErrGRNItemPOLineRequired }
	if l.QuantityReceived < 0 { return ErrGRNItemQtyRequired }
	return nil
}

type InvoiceStatus string
const (InvoiceDraft InvoiceStatus="DRAFT"; InvoiceVerified InvoiceStatus="VERIFIED"; InvoicePosted InvoiceStatus="POSTED"; InvoicePaid InvoiceStatus="PAID"; InvoiceCancelled InvoiceStatus="CANCELLED")

type VATDeductionStatus string
const (VATPending VATDeductionStatus="pending"; VATClaimed VATDeductionStatus="claimed"; VATRejected VATDeductionStatus="rejected")

type SupplierInvoice struct {
	ID                   string              `json:"id"`
	CompanyID            string              `json:"company_id"`
	InvoiceNumber        string              `json:"invoice_number"`
	InvoiceDate          time.Time           `json:"invoice_date"`
	POID                 string              `json:"po_id,omitempty"`
	GRNID                string              `json:"grn_id,omitempty"`
	SupplierID           string              `json:"supplier_id"`
	SupplierName         string              `json:"supplier_name"`
	SupplierTaxCode      string              `json:"supplier_tax_code"`
	InvoiceType          string              `json:"invoice_type"`
	Currency             string              `json:"currency"`
	ExchangeRate         float64             `json:"exchange_rate"`
	Subtotal             float64             `json:"subtotal"`
	DiscountAmount       float64             `json:"discount_amount"`
	TaxAmount            float64             `json:"tax_amount"`
	TotalAmount          float64             `json:"total_amount"`
	AmountPaid           float64             `json:"amount_paid"`
	BalanceDue           float64             `json:"balance_due"`
	DueDate              *time.Time          `json:"due_date,omitempty"`
	VATDeductionStatus   VATDeductionStatus  `json:"vat_deduction_status"`
	EInvoiceData         string              `json:"e_invoice_data,omitempty"`
	EInvoiceCode         string              `json:"e_invoice_code,omitempty"`
	Status               InvoiceStatus       `json:"status"`
	GLPosted             bool                `json:"gl_posted"`
	GLPostedAt           *time.Time          `json:"gl_posted_at,omitempty"`
	Notes                string              `json:"notes,omitempty"`
	CreatedBy            string              `json:"created_by"`
	CreatedAt            time.Time           `json:"created_at,omitempty"`
	Lines                []SupplierInvoiceLine `json:"lines"`
}

func (inv *SupplierInvoice) Validate() error {
	if inv.InvoiceNumber == "" { return ErrInvoiceNumberRequired }
	if inv.SupplierID == "" { return ErrInvoiceSupplierRequired }
	if inv.SupplierName == "" { return ErrInvoiceSupplierNameRequired }
	if inv.SupplierTaxCode == "" { return ErrInvoiceSupplierTaxCodeRequired }
	if inv.InvoiceDate.IsZero() { return ErrInvoiceDateRequired }
	if inv.Currency == "" { inv.Currency = "VND" }
	if inv.Status == "" { inv.Status = InvoiceDraft }
	switch inv.Status {
	case InvoiceDraft, InvoiceVerified, InvoicePosted, InvoicePaid, InvoiceCancelled:
	default: return ErrInvoiceStatusInvalidPurchase
	}
	if len(inv.Lines) == 0 { return ErrInvoiceLinesRequired }
	return nil
}

type SupplierInvoiceLine struct {
	ID           string  `json:"id"`
	InvoiceID    string  `json:"invoice_id"`
	POLineID     string  `json:"po_line_id,omitempty"`
	GRNLineID    string  `json:"grn_line_id,omitempty"`
	ItemCode     string  `json:"item_code,omitempty"`
	ItemName     string  `json:"item_name"`
	Unit         string  `json:"unit"`
	Quantity     float64 `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	VATRate      float64 `json:"vat_rate"`
	VATType      VATType `json:"vat_type"`
	LineTotal    float64 `json:"line_total"`
	LineVATAmount float64 `json:"line_vat_amount"`
	AccountID    string  `json:"account_id"`
	VATAccountID string  `json:"vat_account_id"`
}

func (l *SupplierInvoiceLine) Validate() error {
	if l.ItemName == "" { return ErrInvoiceItemNameRequired }
	if l.Quantity <= 0 { return ErrInvoiceItemQtyRequired }
	if l.UnitPrice < 0 { return ErrInvoiceItemPriceRequired }
	if l.AccountID == "" { return ErrInvoiceItemAccountRequired }
	if l.VATAccountID == "" { return ErrInvoiceItemVATAccountRequired }
	return nil
}

type APTransactionType string
const (APTransInvoice APTransactionType="invoice"; APTransCreditNote APTransactionType="credit_note"; APTransPayment APTransactionType="payment"; APTransPrepayment APTransactionType="prepayment"; APTransOffset APTransactionType="offset")

type APTransaction struct {
	ID              string            `json:"id"`
	CompanyID       string            `json:"company_id"`
	SupplierID      string            `json:"supplier_id"`
	InvoiceID       string            `json:"invoice_id,omitempty"`
	TransactionType APTransactionType `json:"transaction_type"`
	TransactionDate time.Time         `json:"transaction_date"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	ReferenceType   string            `json:"reference_type,omitempty"`
	ReferenceID     string            `json:"reference_id,omitempty"`
	Notes           string            `json:"notes,omitempty"`
	CreatedAt       time.Time         `json:"created_at,omitempty"`
}

func (t *APTransaction) Validate() error {
	if t.SupplierID == "" { return ErrAPTransSupplierRequired }
	if t.TransactionDate.IsZero() { return ErrAPTransDateRequired }
	if t.Amount == 0 { return ErrAPTransAmountRequired }
	if t.Currency == "" { t.Currency = "VND" }
	switch t.TransactionType {
	case APTransInvoice, APTransCreditNote, APTransPayment, APTransPrepayment, APTransOffset:
	default: return ErrAPTransTypeInvalid
	}
	return nil
}

type CostAllocationType string
const (CostTransport CostAllocationType="transport"; CostInsurance CostAllocationType="insurance"; CostCustoms CostAllocationType="customs"; CostInspection CostAllocationType="inspection")

type CostAllocationMethod string
const (CostAllocByQty CostAllocationMethod="by_qty"; CostAllocByValue CostAllocationMethod="by_value"; CostAllocByWeight CostAllocationMethod="by_weight"; CostAllocByVolume CostAllocationMethod="by_volume")

type CostAllocation struct {
	ID               string               `json:"id"`
	CompanyID        string               `json:"company_id"`
	InvoiceID        string               `json:"invoice_id"`
	CostType         CostAllocationType   `json:"cost_type"`
	CostAmount       float64              `json:"cost_amount"`
	AllocationMethod CostAllocationMethod `json:"allocation_method"`
	AllocatedLines   string               `json:"allocated_lines"`
	Notes            string               `json:"notes,omitempty"`
}

func (c *CostAllocation) Validate() error {
	if c.InvoiceID == "" { return ErrCostAllocInvoiceRequired }
	if c.CostAmount <= 0 { return ErrCostAllocAmountRequired }
	switch c.CostType {
	case CostTransport, CostInsurance, CostCustoms, CostInspection:
	default: return ErrCostAllocTypeInvalid
	}
	switch c.AllocationMethod {
	case CostAllocByQty, CostAllocByValue, CostAllocByWeight, CostAllocByVolume:
	default: return ErrCostAllocMethodInvalid
	}
	return nil
}

type PurchaseOrderFilter struct {
	CompanyID  string
	SupplierID string
	Status     POStatus
	FromDate   string
	ToDate     string
	Offset     int
	Limit      int
}

type SupplierInvoiceFilter struct {
	CompanyID  string
	SupplierID string
	Status     InvoiceStatus
	FromDate   string
	ToDate     string
	Offset     int
	Limit      int
}

type GRNFilter struct {
	CompanyID string
	POID      string
	Status    GRNStatus
	FromDate  string
	ToDate    string
	Offset    int
	Limit     int
}

type APAgingBucket struct {
	Bucket0  float64 `json:"bucket_0"`
	Bucket30 float64 `json:"bucket_30"`
	Bucket60 float64 `json:"bucket_60"`
	Bucket90 float64 `json:"bucket_90"`
	Total    float64 `json:"total"`
}

type APAgingReport struct {
	SupplierID   string          `json:"supplier_id"`
	SupplierName string          `json:"supplier_name"`
	TaxCode      string          `json:"tax_code"`
	Buckets      APAgingBucket   `json:"buckets"`
}

type APSummary struct {
	SupplierID      string  `json:"supplier_id"`
	SupplierName    string  `json:"supplier_name"`
	TaxCode         string  `json:"tax_code"`
	TotalInvoiced   float64 `json:"total_invoiced"`
	TotalPaid       float64 `json:"total_paid"`
	Outstanding     float64 `json:"outstanding"`
	Currency        string  `json:"currency"`
}
