package domain

type DeclarationType string
const (
	DeclTypeGTGT01    DeclarationType = "GTGT01"
	DeclTypeGTGT02    DeclarationType = "GTGT02"
	DeclTypeGTGT03    DeclarationType = "GTGT03"
	DeclTypeGTGT04    DeclarationType = "GTGT04"
	DeclTypeGTGT05    DeclarationType = "GTGT05"
	DeclTypeTNDN03    DeclarationType = "TNDN03"
	DeclTypeTNDN04    DeclarationType = "TNDN04"
	DeclTypeTNDN02    DeclarationType = "TNDN02"
	DeclTypeTNDN05    DeclarationType = "TNDN05"
	DeclTypeTNDN06    DeclarationType = "TNDN06"
	DeclTypeKKTNCN    DeclarationType = "KK_TNCN"
	DeclTypeQTTTNCN   DeclarationType = "QTT_TNCN"
	DeclTypeTTDB01    DeclarationType = "TTDB01"
	DeclTypeBVMT01    DeclarationType = "BVMT01"
	DeclTypeNTNN01    DeclarationType = "NTNN01"
	DeclTypeNTNN02    DeclarationType = "NTNN02"
	DeclTypeNTNN03    DeclarationType = "NTNN03"
)
func (dt DeclarationType) Valid() bool {
	switch dt {
	case DeclTypeGTGT01, DeclTypeGTGT02, DeclTypeGTGT03, DeclTypeGTGT04, DeclTypeGTGT05,
		DeclTypeTNDN03, DeclTypeTNDN04, DeclTypeTNDN02, DeclTypeTNDN05, DeclTypeTNDN06,
		DeclTypeKKTNCN, DeclTypeQTTTNCN,
		DeclTypeTTDB01, DeclTypeBVMT01,
		DeclTypeNTNN01, DeclTypeNTNN02, DeclTypeNTNN03:
		return true
	}
	return false
}

type DeclarationStatus string
const (
	DeclStatusDRAFT         DeclarationStatus = "DRAFT"
	DeclStatusVALIDATED     DeclarationStatus = "VALIDATED"
	DeclStatusSUBMITTED     DeclarationStatus = "SUBMITTED"
	DeclStatusACKNOWLEDGED  DeclarationStatus = "ACKNOWLEDGED"
	DeclStatusREJECTED      DeclarationStatus = "REJECTED"
	DeclStatusAMENDED       DeclarationStatus = "AMENDED"
	DeclStatusCANCELLED     DeclarationStatus = "CANCELLED"
	DeclStatusFROZEN        DeclarationStatus = "FROZEN"
)
func (s DeclarationStatus) Valid() bool {
	switch s {
	case DeclStatusDRAFT, DeclStatusVALIDATED, DeclStatusSUBMITTED,
		DeclStatusACKNOWLEDGED, DeclStatusREJECTED,
		DeclStatusAMENDED, DeclStatusCANCELLED, DeclStatusFROZEN:
		return true
	}
	return false
}
func (s DeclarationStatus) IsTerminal() bool {
	return s == DeclStatusACKNOWLEDGED || s == DeclStatusCANCELLED || s == DeclStatusFROZEN
}
func (s DeclarationStatus) CanSubmit() bool {
	return s == DeclStatusDRAFT || s == DeclStatusVALIDATED || s == DeclStatusREJECTED
}

type TaxType string
const (
	TaxTypeVAT      TaxType = "VAT"
	TaxTypeCIT      TaxType = "CIT"
	TaxTypePIT      TaxType = "PIT"
	TaxTypeTTDB     TaxType = "TTDB"
	TaxTypeBVMT     TaxType = "BVMT"
	TaxTypeFCT      TaxType = "FCT"
	TaxTypeRESOURCE TaxType = "RESOURCE"
	TaxTypeLAND     TaxType = "LAND"
)

type RateType string
const (
	RateTypePERCENTAGE    RateType = "PERCENTAGE"
	RateTypeFIXED         RateType = "FIXED"
	RateTypePROGRESSIVE   RateType = "PROGRESSIVE"
)

type EInvoiceType string
const (
	EInvTypeORIGINAL        EInvoiceType = "ORIGINAL"
	EInvTypeADJUSTMENT      EInvoiceType = "ADJUSTMENT"
	EInvTypeREPLACEMENT     EInvoiceType = "REPLACEMENT"
	EInvTypeCANCELLATION_NOTE EInvoiceType = "CANCELLATION_NOTE"
)

type EInvLifecycleStatus string
const (
	EInvStatusDRAFT      EInvLifecycleStatus = "DRAFT"
	EInvStatusSIGNED     EInvLifecycleStatus = "SIGNED"
	EInvStatusSUBMITTED  EInvLifecycleStatus = "SUBMITTED"
	EInvStatusVALIDATED  EInvLifecycleStatus = "VALIDATED"
	EInvStatusISSUED     EInvLifecycleStatus = "ISSUED"
	EInvStatusCANCELLED  EInvLifecycleStatus = "CANCELLED"
	EInvStatusREPLACED   EInvLifecycleStatus = "REPLACED"
)
func (s EInvLifecycleStatus) Valid() bool {
	switch s {
	case EInvStatusDRAFT, EInvStatusSIGNED, EInvStatusSUBMITTED,
		EInvStatusVALIDATED, EInvStatusISSUED,
		EInvStatusCANCELLED, EInvStatusREPLACED:
		return true
	}
	return false
}
// CanCancel reports whether the invoice may be cancelled. Per spec §2.2 only
// an ISSUED invoice can be cancelled (Circular 78/2021: cancel/replace
// applies to invoices already carrying the CQT code).
func (s EInvLifecycleStatus) CanCancel() bool {
	return s == EInvStatusISSUED
}

type CalendarStatus string
const (
	CalStatusPENDING   CalendarStatus = "PENDING"
	CalStatusSUBMITTED CalendarStatus = "SUBMITTED"
	CalStatusPAID      CalendarStatus = "PAID"
	CalStatusMISSED    CalendarStatus = "MISSED"
	CalStatusOVERDUE   CalendarStatus = "OVERDUE"
)

type PaymentMethod string
const (
	PayMethodEFT          PaymentMethod = "EFT"
	PayMethodBANK_TRANSFER PaymentMethod = "BANK_TRANSFER"
	PayMethodCASH         PaymentMethod = "CASH"
	PayMethodCQ           PaymentMethod = "CQ"
)

type PaymentStatus string
const (
	PayStatusPENDING  PaymentStatus = "PENDING"
	PayStatusPAID     PaymentStatus = "PAID"
	PayStatusPARTIAL  PaymentStatus = "PARTIAL"
	PayStatusOVERPAID PaymentStatus = "OVERPAID"
)

type AdjustmentType string
const (
	AdjTypeNONE       AdjustmentType = "NONE"
	AdjTypeAMENDMENT  AdjustmentType = "AMENDMENT"
	AdjTypeADDITIONAL AdjustmentType = "ADDITIONAL"
)

type SourceType string
const (
	SrcTypeAUTO       SourceType = "AUTO_CALCULATED"
	SrcTypeMANUAL     SourceType = "MANUAL_ENTRY"
	SrcTypeFROM_LEDGER SourceType = "FROM_LEDGER"
)

type AlertType string
const (
	AlertTypeINFO      AlertType = "INFO"
	AlertTypeWARNING   AlertType = "WARNING"
	AlertTypeCRITICAL  AlertType = "CRITICAL"
	AlertTypeDUE_TODAY AlertType = "DUE_TODAY"
)

type AlertChannel string
const (
	AlertChanEMAIL AlertChannel = "EMAIL"
	AlertChanINAPP AlertChannel = "IN_APP"
	AlertChanSMS   AlertChannel = "SMS"
	AlertChanALL   AlertChannel = "ALL"
)

type AuditCaseStatus string
const (
	AuditCaseOPEN       AuditCaseStatus = "OPEN"
	AuditCaseINPROGRESS AuditCaseStatus = "IN_PROGRESS"
	AuditCaseCLOSED     AuditCaseStatus = "CLOSED"
)

type TaxPeriod struct {
	PeriodType   PeriodTypeV2 `json:"period_type"`
	PeriodYear   int         `json:"period_year"`
	PeriodNumber int         `json:"period_number"`
}

type ProgressiveBracket struct {
	MinAmount   float64 `json:"min_amount"`
	MaxAmount   float64 `json:"max_amount"`
	RatePercent float64 `json:"rate_percent"`
	FlatAmount  float64 `json:"flat_amount,omitempty"`
}

type TaxRate struct {
	ID              string               `json:"id"`
	TaxType         TaxType              `json:"tax_type"`
	RateCode        string               `json:"rate_code"`
	RateName        string               `json:"rate_name"`
	RateType        RateType             `json:"rate_type"`
	RateValue       float64              `json:"rate_value,omitempty"`
	Brackets        []ProgressiveBracket `json:"brackets,omitempty"`
	EffectiveFrom   string               `json:"effective_from"`
	EffectiveTo     string               `json:"effective_to,omitempty"`
	IsActive        bool                 `json:"is_active"`
	ApplicableTo    string               `json:"applicable_to,omitempty"`
	LegalRef        string               `json:"legal_ref,omitempty"`
	CreatedAt       string               `json:"created_at,omitempty"`
}

type DeclarationSignature struct {
	SignatureID string `json:"signature_id"`
	SignedAt    string `json:"signed_at"`
	SignedBy    string `json:"signed_by"`
	SignatureData string `json:"signature_data,omitempty"`
}

type TaxDeclaration struct {
	ID                  string               `json:"id"`
	CompanyID           string               `json:"company_id"`
	DeclarationType     DeclarationType      `json:"declaration_type"`
	TaxPeriod           TaxPeriod            `json:"tax_period"`
	Status              DeclarationStatus    `json:"status"`
	SubmittedAt         string               `json:"submitted_at,omitempty"`
	SubmittedBy         string               `json:"submitted_by,omitempty"`
	AcknowledgedAt      string               `json:"acknowledged_at,omitempty"`
	AcknowledgementRef  string               `json:"acknowledgement_ref,omitempty"`
	GDTSubmissionID     string               `json:"gdt_submission_id,omitempty"`
	GDTResponseXML      string               `json:"gdt_response_xml,omitempty"`
	DeclarationXML      string               `json:"declaration_xml,omitempty"`
	PreviousDeclID      string               `json:"previous_declaration_id,omitempty"`
	AdjustmentType      AdjustmentType       `json:"adjustment_type"`
	Lines               []TaxDeclarationLine `json:"lines,omitempty"`
	Signatures          []DeclarationSignature `json:"signatures,omitempty"`
	Version             int                  `json:"version"`
	CreatedAt           string               `json:"created_at"`
	CreatedBy           string               `json:"created_by"`
	UpdatedAt           string               `json:"updated_at"`
}

func (d *TaxDeclaration) Validate() error {
	if d.CompanyID == "" { return ErrCompanyIDRequired }
	if !d.DeclarationType.Valid() { return ErrDeclarationTypeInvalid }
	if d.TaxPeriod.PeriodYear < 2000 || d.TaxPeriod.PeriodYear > 2100 { return ErrPeriodYearOutOfRange }
	if d.TaxPeriod.PeriodNumber < 1 { return ErrPeriodNumberInvalid }
	switch d.TaxPeriod.PeriodType {
	case PeriodTypeMonthly:
		if d.TaxPeriod.PeriodNumber > 12 { return ErrPeriodNumberInvalid }
	case PeriodTypeQuarterly:
		if d.TaxPeriod.PeriodNumber > 4 { return ErrPeriodNumberInvalid }
	}
	if d.Status == "" { d.Status = DeclStatusDRAFT }
	if !d.Status.Valid() { return ErrDeclarationStatusInvalid }
	return nil
}

func (d *TaxDeclaration) CanSubmit() bool {
	return d.Status.CanSubmit()
}

func (d *TaxDeclaration) CanAmend() bool {
	return d.Status == DeclStatusACKNOWLEDGED
}

type TaxDeclarationLine struct {
	ID              string     `json:"id,omitempty"`
	DeclarationID   string     `json:"declaration_id,omitempty"`
	LineCode        string     `json:"line_code"`
	LineName        string     `json:"line_name"`
	Amount          float64    `json:"amount"`
	SourceType      SourceType `json:"source_type"`
	SourceAccount   string     `json:"source_account,omitempty"`
	SourceEntryIDs  []string   `json:"source_entry_ids,omitempty"`
	Note            string     `json:"note,omitempty"`
	SortOrder       int        `json:"sort_order"`
}

type TaxPayment struct {
	ID              string        `json:"id"`
	CompanyID       string        `json:"company_id"`
	DeclarationID   string        `json:"declaration_id,omitempty"`
	TaxType         TaxType       `json:"tax_type"`
	PeriodYear      int           `json:"period_year"`
	PeriodNumber    int           `json:"period_number"`
	DeclaredAmount  float64       `json:"declared_amount"`
	PaidAmount      float64       `json:"paid_amount"`
	PaymentDate     string        `json:"payment_date,omitempty"`
	DueDate         string        `json:"due_date"`
	PaymentRef      string        `json:"payment_ref,omitempty"`
	PaymentMethod   PaymentMethod `json:"payment_method,omitempty"`
	Status          PaymentStatus `json:"status"`
	LateDays        int           `json:"late_days,omitempty"`
	LateInterest    float64       `json:"late_interest,omitempty"`
	Notes           string        `json:"notes,omitempty"`
	CreatedAt       string        `json:"created_at"`
}

func (p *TaxPayment) Validate() error {
	if p.CompanyID == "" { return ErrCompanyIDRequired }
	if p.DeclaredAmount <= 0 { return ErrPaymentAmountRequired }
	if p.DueDate == "" { return ErrPaymentDueDateRequired }
	if p.Status == "" { p.Status = PayStatusPENDING }
	return nil
}

func (p *TaxPayment) CalculateLateInterest() {
	if p.PaidAmount >= p.DeclaredAmount {
		p.LateDays = 0
		p.LateInterest = 0
		return
	}
	underpaid := p.DeclaredAmount - p.PaidAmount
	p.LateInterest = underpaid * 0.0003 * float64(p.LateDays)
}

type EInvoiceLine struct {
	ID          string  `json:"id,omitempty"`
	EInvoiceID  string  `json:"e_invoice_id,omitempty"`
	LineNumber  int     `json:"line_number"`
	Description string  `json:"description"`
	Unit        string  `json:"unit,omitempty"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	LineTotal   float64 `json:"line_total"`
	VATRate     float64 `json:"vat_rate"`
	VATAmount   float64 `json:"vat_amount"`
}

type EInvoice struct {
	ID                string         `json:"id"`
	CompanyID         string         `json:"company_id"`
	Pattern           string         `json:"pattern"`
	Serial            string         `json:"serial"`
	InvoiceNumber     int64          `json:"invoice_number,omitempty"`
	InvoiceType       EInvoiceType   `json:"invoice_type"`
	GDTTransactionID  string         `json:"gdt_transaction_id,omitempty"`
	BuyerName         string         `json:"buyer_name"`
	BuyerTaxCode      string         `json:"buyer_tax_code,omitempty"`
	BuyerAddress      string         `json:"buyer_address,omitempty"`
	BuyerEmail        string         `json:"buyer_email,omitempty"`
	CurrencyCode      string         `json:"currency_code"`
	ExchangeRate      float64        `json:"exchange_rate,omitempty"`
	Lines             []EInvoiceLine `json:"lines"`
	Subtotal          float64        `json:"subtotal"`
	VATAmount         float64        `json:"vat_amount"`
	GrandTotal        float64        `json:"grand_total"`
	XMLBody           string         `json:"xml_body,omitempty"`
	SignedXML         string         `json:"signed_xml,omitempty"`
	IssueDate         string         `json:"issue_date"`
	SigningDate       string         `json:"signing_date,omitempty"`
	DigitalSignatureID string        `json:"digital_signature_id,omitempty"`
	JournalEntryID    string         `json:"journal_entry_id,omitempty"`
	Status            EInvLifecycleStatus `json:"status"`
	CancelledAt       string         `json:"cancelled_at,omitempty"`
	CancelReason      string         `json:"cancel_reason,omitempty"`
	OriginalInvoiceID string         `json:"original_invoice_id,omitempty"`
	GDTResponse       string         `json:"gdt_response,omitempty"`
	CreatedAt         string         `json:"created_at"`
}

func (inv *EInvoice) Validate() error {
	if inv.CompanyID == "" { return ErrCompanyIDRequired }
	if inv.BuyerName == "" { return ErrBuyerNameRequired }
	if inv.Pattern == "" { return ErrInvoicePatternRequired }
	if len(inv.Lines) == 0 { return ErrInvoiceNoLines }
	totalVAT := 0.0
	totalLines := 0.0
	for _, line := range inv.Lines {
		computedVAT := line.LineTotal * line.VATRate / 100.0
		if computedVAT != line.VATAmount {
			return ErrInvoiceVATMismatch
		}
		totalVAT += line.VATAmount
		totalLines += line.LineTotal
	}
	if totalLines != inv.Subtotal { return ErrInvoiceSubtotalMismatch }
	if totalVAT != inv.VATAmount { return ErrInvoiceVATTotalMismatch }
	if inv.Subtotal + inv.VATAmount != inv.GrandTotal { return ErrInvoiceGrandTotalMismatch }
	if inv.Status == "" { inv.Status = EInvStatusDRAFT }
	if !inv.Status.Valid() { return ErrInvoiceStatusInvalid }
	if inv.CurrencyCode == "" { inv.CurrencyCode = "VND" }
	switch inv.InvoiceType {
	case EInvTypeADJUSTMENT, EInvTypeREPLACEMENT, EInvTypeCANCELLATION_NOTE:
		if inv.OriginalInvoiceID == "" { return ErrOriginalInvoiceRequired }
	}
	return nil
}

type TaxCalendar struct {
	ID              string         `json:"id"`
	CompanyID       string         `json:"company_id"`
	TaxType         TaxType        `json:"tax_type"`
	PeriodType      PeriodTypeV2   `json:"period_type"`
	PeriodYear      int            `json:"period_year"`
	PeriodNumber    int            `json:"period_number"`
	StartDate       string         `json:"start_date"`
	EndDate         string         `json:"end_date"`
	DeclarationDue  string         `json:"declaration_due"`
	PaymentDue      string         `json:"payment_due,omitempty"`
	Status          CalendarStatus `json:"status"`
	DeclarationID   string         `json:"declaration_id,omitempty"`
	CreatedAt       string         `json:"created_at"`
}

type TaxAlert struct {
	ID            string       `json:"id"`
	CompanyID     string       `json:"company_id"`
	CalendarID    string       `json:"calendar_id,omitempty"`
	AlertType     AlertType    `json:"alert_type"`
	Channel       AlertChannel `json:"channel"`
	Message       string       `json:"message"`
	SentAt        string       `json:"sent_at"`
	AcknowledgedAt string      `json:"acknowledged_at,omitempty"`
	AcknowledgedBy string      `json:"acknowledged_by,omitempty"`
}

type TaxAuditCase struct {
	ID               string         `json:"id"`
	CompanyID        string         `json:"company_id"`
	AuditPeriodStart string         `json:"audit_period_start"`
	AuditPeriodEnd   string         `json:"audit_period_end"`
	AuditDecNumber   string         `json:"audit_decision_number"`
	AuditorName      string         `json:"auditor_name"`
	AuditorContact   string         `json:"auditor_contact,omitempty"`
	Status           AuditCaseStatus `json:"status"`
	Findings         string         `json:"findings,omitempty"`
	PenaltyAmount    float64        `json:"penalty_amount,omitempty"`
	CreatedAt        string         `json:"created_at"`
	ClosedAt         string         `json:"closed_at,omitempty"`
}

type VATResult struct {
	CompanyID        string  `json:"company_id"`
	Period           TaxPeriod `json:"period"`
	OutputVAT        float64 `json:"output_vat"`
	InputVAT         float64 `json:"input_vat"`
	InputVATFA       float64 `json:"input_vat_fa"`
	TotalInputVAT    float64 `json:"total_input_vat"`
	VATPayable       float64 `json:"vat_payable"`
	VATRefundable    float64 `json:"vat_refundable"`
	PurchaseTotal    float64 `json:"purchase_total"`
	SalesTotal       float64 `json:"sales_total"`
}

type CITResult struct {
	CompanyID        string    `json:"company_id"`
	PeriodYear       int       `json:"period_year"`
	PeriodType       string    `json:"period_type"`
	Revenue          float64   `json:"revenue"`
	Expenses         float64   `json:"expenses"`
	NonDeductible    float64   `json:"non_deductible"`
	TaxableIncome    float64   `json:"taxable_income"`
	TaxRate          float64   `json:"tax_rate"`
	CITPayable       float64   `json:"cit_payable"`
	IncentiveReduc   float64   `json:"incentive_reduction"`
	CITFinal         float64   `json:"cit_final"`
	Provisionals     float64   `json:"provisionals"`
	BalanceDue       float64   `json:"balance_due"`
	Refund           float64   `json:"refund"`
	LateInterest     float64   `json:"late_interest"`
}

type PITResult struct {
	CompanyID       string  `json:"company_id"`
	Period          TaxPeriod `json:"period"`
	EmployeeCount   int     `json:"employee_count"`
	TotalGross      float64 `json:"total_gross"`
	TotalDeductions float64 `json:"total_deductions"`
	TotalPIT        float64 `json:"total_pit"`
}

// PITEmployeeInput carries per-employee payroll data for PIT calculation.
// NonResident defaults to false (resident): the safe default applies the
// progressive schedule with deductions.
type PITEmployeeInput struct {
	EmployeeID         string  `json:"employee_id"`
	GrossMonthly       float64 `json:"gross_monthly"`          // gross salary per month (VND)
	Dependants         int     `json:"dependants"`             // registered dependants
	Months             int     `json:"months"`                 // months covered by the period (1 = monthly)
	NonResident        bool    `json:"non_resident,omitempty"` // flat 20%, no deductions
	SocialInsuranceBase float64 `json:"social_insurance_base,omitempty"` // capped base; defaults to gross
}

type TaxDeclarationFilter struct {
	CompanyID       string           `json:"company_id,omitempty"`
	DeclarationType DeclarationType  `json:"declaration_type,omitempty"`
	Status          DeclarationStatus `json:"status,omitempty"`
	PeriodYear      int              `json:"period_year,omitempty"`
	PeriodNumber    int              `json:"period_number,omitempty"`
	FromDate        string           `json:"from_date,omitempty"`
	ToDate          string           `json:"to_date,omitempty"`
}

type TaxRateFilter struct {
	TaxType     TaxType `json:"tax_type,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	EffectiveOn string  `json:"effective_on,omitempty"`
}

type PaymentFilter struct {
	CompanyID    string        `json:"company_id,omitempty"`
	TaxType      TaxType       `json:"tax_type,omitempty"`
	Status       PaymentStatus `json:"status,omitempty"`
	PeriodYear   int           `json:"period_year,omitempty"`
	PeriodNumber int           `json:"period_number,omitempty"`
}

type EInvoiceFilter struct {
	CompanyID string         `json:"company_id,omitempty"`
	Status    EInvLifecycleStatus `json:"status,omitempty"`
	FromDate  string         `json:"from_date,omitempty"`
	ToDate    string         `json:"to_date,omitempty"`
}

type EInvoiceInput struct {
	CompanyID         string         `json:"company_id"`
	Pattern           string         `json:"pattern"`
	Serial            string         `json:"serial"`
	InvoiceType       EInvoiceType   `json:"invoice_type"`
	BuyerName         string         `json:"buyer_name"`
	BuyerTaxCode      string         `json:"buyer_tax_code,omitempty"`
	BuyerAddress      string         `json:"buyer_address,omitempty"`
	BuyerEmail        string         `json:"buyer_email,omitempty"`
	CurrencyCode      string         `json:"currency_code,omitempty"`
	ExchangeRate      float64        `json:"exchange_rate,omitempty"`
	Lines             []EInvoiceLine `json:"lines"`
	DigitalSignatureID string        `json:"digital_signature_id,omitempty"`
	IssueDate         string         `json:"issue_date"`
}

// VATRateReconciliation is the per-VAT-rate breakdown of issued e-invoices
// for a period (BR-VAT-06: e-invoice data must reconcile with GTGT01).
type VATRateReconciliation struct {
	Rate          float64 `json:"rate"`
	InvoiceCount  int     `json:"invoice_count"`
	InvoiceAmount float64 `json:"invoice_amount"`
	InvoiceVAT    float64 `json:"invoice_vat"`
}

// VATReconciliationResult compares issued e-invoices against the GTGT01
// declaration's output VAT ([23]) for the same company + period.
type VATReconciliationResult struct {
	CompanyID        string                    `json:"company_id"`
	Period           TaxPeriod                 `json:"period"`
	DeclarationID    string                    `json:"declaration_id,omitempty"`
	DeclarationTotal float64                   `json:"declaration_total"`
	InvoiceTotal     float64                   `json:"invoice_total"`
	Variance         float64                   `json:"variance"`
	Matched          bool                      `json:"matched"`
	InvoiceCount     int                       `json:"invoice_count"`
	ByRate           []VATRateReconciliation   `json:"by_rate"`
	ExcludedInvoices []string                  `json:"excluded_invoices,omitempty"`
	Notes            []string                  `json:"notes,omitempty"`
}
