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
	Code              string         `json:"code" validate:"required"`
	Name              string         `json:"name" validate:"required"`
	TaxCode           string         `json:"tax_code" validate:"required"`
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
	Status            SupplierStatus `json:"status" validate:"suppstatus"`
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
	PONumber        string    `json:"po_number" validate:"required"`
	SupplierID      string    `json:"supplier_id" validate:"required"`
	RequisitionID   string    `json:"requisition_id,omitempty"`
	OrderDate       time.Time `json:"order_date" validate:"required"`
	ExpectedDate    *time.Time `json:"expected_date,omitempty"`
	Currency        string    `json:"currency"`
	ExchangeRate    float64   `json:"exchange_rate"`
	PaymentTerms    string    `json:"payment_terms,omitempty"`
	DeliveryTerms   string    `json:"delivery_terms,omitempty"`
	Subtotal        float64   `json:"subtotal"`
	DiscountAmount  float64   `json:"discount_amount"`
	TaxAmount       float64   `json:"tax_amount"`
	TotalAmount     float64   `json:"total_amount"`
	Status          POStatus  `json:"status" validate:"postatus"`
	ApprovedBy      string    `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	CancelledReason string    `json:"cancelled_reason,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	Lines           []POItem  `json:"lines" validate:"min=1"`
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
		rawTotal := l.Quantity * l.UnitPrice
		l.LineTotal = rawTotal
		if l.DiscountPct > 0 {
			l.LineTotal = rawTotal - rawTotal*l.DiscountPct/100
			discount += rawTotal * l.DiscountPct / 100
		}
		l.LineVATAmount = l.LineTotal * l.VATRate / 100
		subtotal += l.LineTotal
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
	ItemName     string  `json:"item_name" validate:"required"`
	Unit         string  `json:"unit" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"gt=0"`
	UnitPrice    float64 `json:"unit_price" validate:"gte=0"`
	DiscountPct  float64 `json:"discount_pct"`
	VATRate      float64 `json:"vat_rate"`
	VATType      VATType `json:"vat_type"`
	AccountID    string  `json:"account_id" validate:"required"`
	VATAccountID string  `json:"vat_account_id" validate:"required"`
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

func (s GRNStatus) ValidTransition(next GRNStatus) bool {
	switch s {
	case GRNDraft: return next == GRNPosted || next == GRNCancelled
	case GRNPosted: return next == GRNCancelled
	case GRNCancelled: return false
	default: return false
	}
}

type GRN struct {
	ID              string     `json:"id"`
	CompanyID       string     `json:"company_id"`
	GRNNumber       string     `json:"grn_number" validate:"required"`
	POID            string     `json:"po_id,omitempty"`
	WarehouseID     string     `json:"warehouse_id,omitempty"`
	ReturnOfGRNID   string     `json:"return_of_grn_id,omitempty"`
	ReceiptDate     time.Time  `json:"receipt_date" validate:"required"`
	Warehouse       string     `json:"warehouse,omitempty"`
	Status          GRNStatus  `json:"status" validate:"grnstatus"`
	Notes           string     `json:"notes,omitempty"`
	CreatedBy       string     `json:"created_by"`
	PostedAt        time.Time  `json:"posted_at,omitempty"`
	CancelledReason string     `json:"cancelled_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at,omitempty"`
	Lines           []GRNItem  `json:"lines" validate:"min=1"`
}

func (g *GRN) Validate() error {
	if g.GRNNumber == "" { return ErrGRNNumberRequired }
	if g.POID == "" && g.WarehouseID == "" { return ErrGRNPORequired }
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
	POLineID         string  `json:"po_line_id,omitempty" validate:"required"`
	ItemID           string  `json:"item_id,omitempty"`
	ItemCode         string  `json:"item_code,omitempty"`
	ItemName         string  `json:"item_name" validate:"required"`
	Unit             string  `json:"unit"`
	QuantityReceived float64 `json:"quantity_received" validate:"gte=0"`
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

const (InvoiceTypeInvoice="invoice"; InvoiceTypeCreditNote="credit_note"; InvoiceTypeImport="import")

type VATDeductionStatus string
const (VATPending VATDeductionStatus="pending"; VATClaimed VATDeductionStatus="claimed"; VATRejected VATDeductionStatus="rejected")

type SupplierInvoice struct {
	ID                   string              `json:"id"`
	CompanyID            string              `json:"company_id"`
	InvoiceNumber        string              `json:"invoice_number" validate:"required"`
	InvoiceDate          time.Time           `json:"invoice_date" validate:"required"`
	POID                 string              `json:"po_id,omitempty"`
	GRNID                string              `json:"grn_id,omitempty"`
	SupplierID           string              `json:"supplier_id" validate:"required"`
	SupplierName         string              `json:"supplier_name" validate:"required"`
	SupplierTaxCode      string              `json:"supplier_tax_code" validate:"required"`
	InvoiceType          string              `json:"invoice_type"`
	OriginalInvoiceID    string              `json:"original_invoice_id,omitempty"`
	ImportDuty           float64             `json:"import_duty,omitempty"`
	ImportVAT            float64             `json:"import_vat,omitempty"`
	CustomsDeclarationNumber string          `json:"customs_declaration_number,omitempty"`
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
	Status               InvoiceStatus       `json:"status" validate:"invstatus"`
	GLPosted             bool                `json:"gl_posted"`
	GLPostedAt           *time.Time          `json:"gl_posted_at,omitempty"`
	Notes                string              `json:"notes,omitempty"`
	CreatedBy            string              `json:"created_by"`
	CreatedAt            time.Time           `json:"created_at,omitempty"`
	UpdatedAt            time.Time           `json:"updated_at,omitempty"`
	Lines                []SupplierInvoiceLine `json:"lines" validate:"min=1"`
}

func (inv *SupplierInvoice) Validate() error {
	if inv.InvoiceNumber == "" { return ErrInvoiceNumberRequired }
	if inv.SupplierID == "" { return ErrInvoiceSupplierRequired }
	if inv.SupplierName == "" { return ErrInvoiceSupplierNameRequired }
	if inv.SupplierTaxCode == "" { return ErrInvoiceSupplierTaxCodeRequired }
	if inv.InvoiceDate.IsZero() { return ErrInvoiceDateRequired }
	if inv.InvoiceType == InvoiceTypeCreditNote && inv.OriginalInvoiceID == "" { return ErrCreditNoteOriginalRequired }
	if inv.Currency == "" { inv.Currency = "VND" }
	if inv.Status == "" { inv.Status = InvoiceDraft }
	switch inv.Status {
	case InvoiceDraft, InvoiceVerified, InvoicePosted, InvoicePaid, InvoiceCancelled:
	default: return ErrInvoiceStatusInvalidPurchase
	}
	if len(inv.Lines) == 0 { return ErrInvoiceLinesRequired }
	return nil
}

func (inv *SupplierInvoice) CalculateTotals() {
	subtotal, tax := 0.0, 0.0
	for i := range inv.Lines {
		l := &inv.Lines[i]
		raw := l.Quantity * l.UnitPrice
		l.LineTotal = raw
		l.LineVATAmount = raw * l.VATRate / 100
		subtotal += l.LineTotal
		tax += l.LineVATAmount
	}
	inv.Subtotal = subtotal
	inv.TaxAmount = tax
	inv.TotalAmount = subtotal + tax
}

type SupplierInvoiceLine struct {
	ID           string  `json:"id"`
	InvoiceID    string  `json:"invoice_id"`
	POLineID     string  `json:"po_line_id,omitempty"`
	GRNLineID    string  `json:"grn_line_id,omitempty"`
	ItemCode     string  `json:"item_code,omitempty"`
	ItemName     string  `json:"item_name" validate:"required"`
	Unit         string  `json:"unit"`
	Quantity     float64 `json:"quantity" validate:"gt=0"`
	UnitPrice    float64 `json:"unit_price" validate:"gte=0"`
	VATRate      float64 `json:"vat_rate"`
	VATType      VATType `json:"vat_type"`
	LineTotal    float64 `json:"line_total"`
	LineVATAmount float64 `json:"line_vat_amount"`
	AccountID    string  `json:"account_id" validate:"required"`
	VATAccountID string  `json:"vat_account_id" validate:"required"`
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
	SupplierID      string            `json:"supplier_id" validate:"required"`
	InvoiceID       string            `json:"invoice_id,omitempty"`
	TransactionType APTransactionType `json:"transaction_type" validate:"apttype"`
	TransactionDate time.Time         `json:"transaction_date" validate:"required"`
	Amount          float64           `json:"amount" validate:"ne=0"`
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
	InvoiceID        string               `json:"invoice_id" validate:"required"`
	CostType         CostAllocationType   `json:"cost_type" validate:"costtype"`
	CostAmount       float64              `json:"cost_amount" validate:"gt=0"`
	AllocationMethod CostAllocationMethod `json:"allocation_method" validate:"allocmethod"`
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
	CompanyID        string
	SupplierID       string
	OriginalInvoiceID string
	Status           InvoiceStatus
	FromDate         string
	ToDate           string
	Offset           int
	Limit            int
}

type GRNFilter struct {
	CompanyID    string
	POID         string
	ReturnOfGRNID string
	Status       GRNStatus
	FromDate     string
	ToDate       string
	Offset       int
	Limit        int
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

type ProvisionStatus string
const (
	ProvisionDraft  ProvisionStatus = "DRAFT"
	ProvisionPosted ProvisionStatus = "POSTED"
)

// Doubtful debt provisioning per Circular 99/2025/TT-BTC (replaces 48/2019):
// 6mo-1yr=30%, 1-2yr=50%, 2-3yr=70%, 3yr+=100%, applied to outstanding
// supplier prepayments (advances) aged by months overdue from the transaction date.
type DoubtfulDebtProvision struct {
	ID             string                        `json:"id"`
	CompanyID      string                        `json:"company_id"`
	AsOfDate       string                        `json:"as_of_date"`
	TotalOutstanding float64                     `json:"total_outstanding"`
	TotalProvision float64                       `json:"total_provision"`
	Status         ProvisionStatus               `json:"status"`
	Lines          []DoubtfulDebtProvisionLine   `json:"lines,omitempty"`
	CreatedBy      string                        `json:"created_by"`
	CreatedAt      time.Time                     `json:"created_at,omitempty"`
}

type DoubtfulDebtProvisionLine struct {
	ID               string  `json:"id"`
	ProvisionID      string  `json:"provision_id,omitempty"`
	SupplierID       string  `json:"supplier_id"`
	SupplierName     string  `json:"supplier_name"`
	TaxCode          string  `json:"tax_code"`
	OutstandingAmount float64 `json:"outstanding_amount"`
	AgeMonths        int     `json:"age_months"`
	RatePct          float64 `json:"rate_pct"`
	ProvisionAmount  float64 `json:"provision_amount"`
}

func (p *DoubtfulDebtProvision) Validate() error {
	if p.AsOfDate == "" {
		return ErrProvisionDateRequired
	}
	if len(p.Lines) == 0 {
		return ErrProvisionNoLines
	}
	return nil
}

// DoubtfulDebtRate returns the Circular 99 provisioning rate for an age in months.
func DoubtfulDebtRate(ageMonths int) float64 {
	switch {
	case ageMonths >= 36:
		return 1.0
	case ageMonths >= 24:
		return 0.70
	case ageMonths >= 12:
		return 0.50
	case ageMonths >= 6:
		return 0.30
	default:
		return 0
	}
}

// ─── Regulatory Reports (Circular 99) ───────────────────────────────────

type PurchaseLedgerRow struct {
	Date         string  `json:"date"`
	DocNumber    string  `json:"doc_number"`
	Description  string  `json:"description"`
	SupplierID   string  `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	Increase     float64 `json:"increase"`
	Decrease     float64 `json:"decrease"`
}

type PurchaseLedgerReport struct {
	CompanyID string              `json:"company_id"`
	FromDate  string              `json:"from_date"`
	ToDate    string              `json:"to_date"`
	Opening   float64             `json:"opening"`
	Increase  float64             `json:"increase"`
	Decrease  float64             `json:"decrease"`
	Closing   float64             `json:"closing"`
	Rows      []PurchaseLedgerRow `json:"rows"`
}

type SupplierLedgerRow struct {
	Date        string  `json:"date"`
	DocNumber   string  `json:"doc_number"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Balance     float64 `json:"balance"`
}

type SupplierLedgerReport struct {
	SupplierID   string              `json:"supplier_id"`
	SupplierName string              `json:"supplier_name"`
	TaxCode      string              `json:"tax_code"`
	FromDate     string              `json:"from_date"`
	ToDate       string              `json:"to_date"`
	Opening      float64             `json:"opening"`
	Closing      float64             `json:"closing"`
	Rows         []SupplierLedgerRow `json:"rows"`
}

type GoodsPurchaseRow struct {
	ItemName  string  `json:"item_name"`
	Unit      string  `json:"unit"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	LineTotal float64 `json:"line_total"`
	VATRate   float64 `json:"vat_rate"`
	VATAmount float64 `json:"vat_amount"`
	AccountID string  `json:"account_id"`
}

type GoodsPurchaseReport struct {
	CompanyID    string              `json:"company_id"`
	FromDate     string              `json:"from_date"`
	ToDate       string              `json:"to_date"`
	TotalQuantity float64            `json:"total_quantity"`
	TotalAmount  float64             `json:"total_amount"`
	TotalVAT     float64             `json:"total_vat"`
	Rows         []GoodsPurchaseRow  `json:"rows"`
}

type VATInputRow struct {
	InvoiceNumber   string  `json:"invoice_number"`
	InvoiceDate     string  `json:"invoice_date"`
	SupplierName    string  `json:"supplier_name"`
	SupplierTaxCode string  `json:"supplier_tax_code"`
	VATRate         float64 `json:"vat_rate"`
	Subtotal        float64 `json:"subtotal"`
	VATAmount       float64 `json:"vat_amount"`
	DeductionStatus string  `json:"deduction_status"`
}

type VATInputReport struct {
	CompanyID    string        `json:"company_id"`
	FromDate     string        `json:"from_date"`
	ToDate       string        `json:"to_date"`
	TotalSubtotal float64      `json:"total_subtotal"`
	TotalVAT     float64       `json:"total_vat"`
	Rows         []VATInputRow `json:"rows"`
}

type UninvoicedReceiptRow struct {
	GRNNumber    string  `json:"grn_number"`
	ReceiptDate  string  `json:"receipt_date"`
	SupplierName string  `json:"supplier_name"`
	POID         string  `json:"po_id"`
	ItemName     string  `json:"item_name"`
	Unit         string  `json:"unit"`
	Quantity     float64 `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	LineTotal    float64 `json:"line_total"`
}

type RequisitionStatus string
const (ReqDraft RequisitionStatus="DRAFT"; ReqPending RequisitionStatus="PENDING"; ReqApproved RequisitionStatus="APPROVED"; ReqRejected RequisitionStatus="REJECTED"; ReqOrdered RequisitionStatus="ORDERED")

func (s RequisitionStatus) ValidTransition(next RequisitionStatus) bool {
	switch s {
	case ReqDraft: return next == ReqPending
	case ReqPending: return next == ReqApproved || next == ReqRejected
	case ReqApproved: return next == ReqOrdered
	case ReqRejected, ReqOrdered: return false
	}
	return false
}

type PurchaseRequisition struct {
	ID                string            `json:"id"`
	CompanyID         string            `json:"company_id"`
	RequisitionNumber string            `json:"requisition_number" validate:"required"`
	RequesterID       string            `json:"requester_id" validate:"required"`
	RequesterName     string            `json:"requester_name,omitempty"`
	DepartmentID      string            `json:"department_id,omitempty"`
	NeedByDate        *time.Time        `json:"need_by_date,omitempty"`
	Priority          string            `json:"priority,omitempty"`
	Reason            string            `json:"reason,omitempty"`
	Status            RequisitionStatus `json:"status" validate:"reqstatus"`
	TotalEstimated    float64           `json:"total_estimated"`
	ApprovedBy        string            `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time        `json:"approved_at,omitempty"`
	RejectedReason    string            `json:"rejected_reason,omitempty"`
	CreatedBy         string            `json:"created_by"`
	CreatedAt         time.Time         `json:"created_at,omitempty"`
	UpdatedAt         time.Time         `json:"updated_at,omitempty"`
	Lines             []RequisitionItem `json:"lines" validate:"min=1"`
}

func (r *PurchaseRequisition) Validate() error {
	if r.RequisitionNumber == "" { return ErrRequisitionNumberRequired }
	if r.RequesterID == "" { return ErrRequisitionRequesterRequired }
	if r.Status == "" { r.Status = ReqDraft }
	switch r.Status {
	case ReqDraft, ReqPending, ReqApproved, ReqRejected, ReqOrdered:
	default: return ErrRequisitionStatusInvalid
	}
	if len(r.Lines) == 0 { return ErrRequisitionLinesRequired }
	return nil
}

func (r *PurchaseRequisition) CalculateTotals() {
	total := 0.0
	for i := range r.Lines {
		l := &r.Lines[i]
		l.EstimatedTotal = l.Quantity * l.EstimatedPrice
		total += l.EstimatedTotal
	}
	r.TotalEstimated = total
}

type RequisitionItem struct {
	ID             string  `json:"id"`
	RequisitionID  string  `json:"requisition_id,omitempty"`
	LineNumber     int     `json:"line_number"`
	ItemCode       string  `json:"item_code,omitempty"`
	ItemName       string  `json:"item_name" validate:"required"`
	Unit           string  `json:"unit"`
	Quantity       float64 `json:"quantity" validate:"gt=0"`
	EstimatedPrice float64 `json:"estimated_price" validate:"gte=0"`
	EstimatedTotal float64 `json:"estimated_total"`
	AccountID      string  `json:"account_id" validate:"required"`
}

func (l *RequisitionItem) Validate() error {
	if l.ItemName == "" { return ErrRequisitionItemNameRequired }
	if l.Quantity <= 0 { return ErrRequisitionItemQtyRequired }
	if l.EstimatedPrice < 0 { return ErrRequisitionItemPriceRequired }
	if l.AccountID == "" { return ErrRequisitionItemAccountRequired }
	return nil
}

type RequisitionFilter struct {
	CompanyID   string
	Status      RequisitionStatus
	RequesterID string
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       int
	Offset      int
}

// ─── FX Revaluation (P2-4) ───────────────────────────────────────────────

type FXRevaluationStatus string
const (FXRevalDraft FXRevaluationStatus="DRAFT"; FXRevalPosted FXRevaluationStatus="POSTED")

type FXRevaluation struct {
	ID              string               `json:"id"`
	CompanyID       string               `json:"company_id"`
	RevaluationDate time.Time            `json:"revaluation_date" validate:"required"`
	Status          FXRevaluationStatus  `json:"status" validate:"fxrstatus"`
	TotalGain       float64              `json:"total_gain"`
	TotalLoss       float64              `json:"total_loss"`
	GLPosted        bool                 `json:"gl_posted"`
	GLPostedAt      *time.Time           `json:"gl_posted_at,omitempty"`
	CreatedBy       string               `json:"created_by"`
	CreatedAt       time.Time            `json:"created_at,omitempty"`
	Lines           []FXRevaluationLine  `json:"lines" validate:"min=1"`
}

func (r *FXRevaluation) Validate() error {
	if r.RevaluationDate.IsZero() { return ErrFXRevaluationDateRequired }
	if r.Status == "" { r.Status = FXRevalDraft }
	switch r.Status {
	case FXRevalDraft, FXRevalPosted:
	default: return ErrFXRevaluationStatusInvalid
	}
	if len(r.Lines) == 0 { return ErrFXRevaluationLinesRequired }
	return nil
}

type FXRevaluationLine struct {
	ID              string  `json:"id"`
	RevaluationID   string  `json:"revaluation_id"`
	InvoiceID       string  `json:"invoice_id" validate:"required"`
	InvoiceNumber   string  `json:"invoice_number" validate:"required"`
	SupplierID      string  `json:"supplier_id"`
	SupplierName    string  `json:"supplier_name"`
	Currency        string  `json:"currency"`
	BalanceDue      float64 `json:"balance_due"`
	OriginalRate    float64 `json:"original_rate"`
	RevaluationRate float64 `json:"revaluation_rate"`
	FxGain          float64 `json:"fx_gain"`
	FxLoss          float64 `json:"fx_loss"`
}
