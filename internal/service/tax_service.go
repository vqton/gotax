package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"gotax/internal/domain"
	"gotax/internal/einvoice"
	"gotax/internal/gdt"
	"gotax/internal/xmldsig"
)

// defaultRateValues holds statutory fallback rates used when the rate table
// has no matching active rate. Keeps calculation deterministic on fresh
// installs (empty table). Keys are applicableTo selectors.
var defaultRateValues = map[domain.TaxType]map[string]float64{
	domain.TaxTypeVAT: {"STANDARD": 10, "REDUCED": 8, "ZERO": 0},
	domain.TaxTypeCIT: {"STANDARD": 20, "SMALL": 17, "MICRO": 15},
	domain.TaxTypePIT: {"STANDARD": 20},
}

// resolveRate returns the active rate for taxType as of onDate (ISO date).
// The applicableTo key is matched against the RateCode suffix (e.g. "SMALL"
// matches "CIT_SMALL") or the ApplicableTo field. RateCode suffix matching
// keeps the resolver backend-agnostic (PG mapping drops ApplicableTo).
// Falls back to a STANDARD-keyed rate, then any active rate, then the
// statutory default.
func (s *taxService) resolveRate(ctx context.Context, taxType domain.TaxType, applicableTo, onDate string) (domain.TaxRate, error) {
	rates, err := s.repo.GetRates(ctx, domain.TaxRateFilter{TaxType: taxType})
	if err != nil {
		return domain.TaxRate{}, err
	}
	key := strings.ToUpper(strings.TrimSpace(applicableTo))
	var active []domain.TaxRate
	for _, r := range rates {
		if r.IsActive && rateAppliesOn(r, onDate) {
			active = append(active, r)
		}
	}
	for _, r := range active {
		if matchesRateKey(r, key) {
			return r, nil
		}
	}
	for _, r := range active {
		if matchesRateKey(r, "STANDARD") {
			return r, nil
		}
	}
	sort.Slice(active, func(i, j int) bool { return active[i].RateCode < active[j].RateCode })
	if len(active) > 0 {
		return active[0], nil
	}
	val, ok := defaultRateValues[taxType][key]
	if !ok {
		val = defaultRateValues[taxType]["STANDARD"]
	}
	return domain.TaxRate{RateType: domain.RateTypePERCENTAGE, RateValue: val, IsActive: true}, nil
}

func matchesRateKey(r domain.TaxRate, key string) bool {
	if key == "" {
		return false
	}
	if r.ApplicableTo != "" && strings.EqualFold(strings.TrimSpace(r.ApplicableTo), key) {
		return true
	}
	return strings.HasSuffix(strings.ToUpper(r.RateCode), key)
}

// rateAppliesOn reports whether the rate is effective on the ISO date.
// Empty onDate means "always". Lexicographic compare is valid for ISO dates.
func rateAppliesOn(r domain.TaxRate, onDate string) bool {
	if onDate == "" {
		return true
	}
	if r.EffectiveFrom != "" && r.EffectiveFrom > onDate {
		return false
	}
	if r.EffectiveTo != "" && r.EffectiveTo < onDate {
		return false
	}
	return true
}

type TaxServiceInterface interface {
	// Declarations
	CreateDeclaration(ctx context.Context, d *domain.TaxDeclaration) error
	GetDeclaration(ctx context.Context, id string) (*domain.TaxDeclaration, error)
	ListDeclarations(ctx context.Context, filter domain.TaxDeclarationFilter) ([]domain.TaxDeclaration, error)
	SubmitDeclaration(ctx context.Context, id, userID string) error
	AcknowledgeDeclaration(ctx context.Context, id, ref string) error
	RejectDeclaration(ctx context.Context, id, reason string) error
	AmendDeclaration(ctx context.Context, id string, lines []domain.TaxDeclarationLine) (*domain.TaxDeclaration, error)
	CancelDeclaration(ctx context.Context, id string) error

	// Tax Rates
	CreateRate(ctx context.Context, rate *domain.TaxRate) error
	GetRate(ctx context.Context, id string) (*domain.TaxRate, error)
	ListRates(ctx context.Context, filter domain.TaxRateFilter) ([]domain.TaxRate, error)
	UpdateRate(ctx context.Context, rate *domain.TaxRate) error

	// Payments
	CreatePayment(ctx context.Context, p *domain.TaxPayment) error
	GetPayment(ctx context.Context, id string) (*domain.TaxPayment, error)
	ListPayments(ctx context.Context, filter domain.PaymentFilter) ([]domain.TaxPayment, error)
	RecordPayment(ctx context.Context, id string, amount float64, date, ref string) error

	// E-Invoices
	CreateEInvoice(ctx context.Context, inv *domain.EInvoice) error
	GetEInvoice(ctx context.Context, id string) (*domain.EInvoice, error)
	ListEInvoices(ctx context.Context, filter domain.EInvoiceFilter) ([]domain.EInvoice, error)
	IssueEInvoice(ctx context.Context, id string) error
	CancelEInvoice(ctx context.Context, id, reason string) error

	// Tax Calendar
	CreateCalendarEntry(ctx context.Context, c *domain.TaxCalendar) error
	GetCalendarEntry(ctx context.Context, id string) (*domain.TaxCalendar, error)
	GetCalendarByPeriod(ctx context.Context, companyID string, periodYear, periodNumber int) ([]domain.TaxCalendar, error)
	GetCalendarByCompany(ctx context.Context, companyID string) ([]domain.TaxCalendar, error)

	// Alerts
	CreateAlert(ctx context.Context, a *domain.TaxAlert) error
	GetAlert(ctx context.Context, id string) (*domain.TaxAlert, error)
	ListAlerts(ctx context.Context, companyID string, limit int) ([]domain.TaxAlert, error)

	// Audit Cases
	CreateAuditCase(ctx context.Context, a *domain.TaxAuditCase) error
	GetAuditCase(ctx context.Context, id string) (*domain.TaxAuditCase, error)
	ListAuditCases(ctx context.Context, companyID string) ([]domain.TaxAuditCase, error)
	CloseAuditCase(ctx context.Context, id string, findings string, penalty float64) error

	// Calculations
	CalculateVAT(ctx context.Context, companyID string, period domain.TaxPeriod, entries []domain.JournalEntry) (*domain.VATResult, error)
	CalculateCIT(ctx context.Context, companyID string, year int, entries []domain.JournalEntry) (*domain.CITResult, error)
	CalculatePIT(ctx context.Context, companyID string, period domain.TaxPeriod, employees []domain.PITEmployeeInput) (*domain.PITResult, error)

	// Declaration Engine
	GenerateDeclaration(ctx context.Context, companyID string, declType domain.DeclarationType, period domain.TaxPeriod, userID string) (*domain.TaxDeclaration, error)

	// E-Invoice Issuance
	CheckInvoiceStatus(ctx context.Context, id string) error
}

// GDTClient is the GDT e-invoice API surface used by the service.
type GDTClient interface {
	SubmitInvoice(ctx context.Context, invoiceXML, certID string) (*gdt.SubmitResponse, error)
	GetInvoiceStatus(ctx context.Context, transactionID string) (*gdt.StatusResponse, error)
	CancelInvoice(ctx context.Context, transactionID, reason string) error
}

// TXMLSigner signs canonicalized GDT TXML (BK:ChuKySo/DuLieuKy).
type TXMLSigner interface {
	SignTXML(xmlBody, signatureID string) (string, error)
}

type taxService struct {
	repo   domain.TaxRepository
	jeRepo domain.JournalRepository
	gdt    GDTClient
	signer TXMLSigner
	now    func() time.Time
}

func NewTaxService(repo domain.TaxRepository, jeRepo domain.JournalRepository, gdt GDTClient, signer TXMLSigner) TaxServiceInterface {
	return &taxService{
		repo:   repo,
		jeRepo: jeRepo,
		gdt:    gdt,
		signer: signer,
		now:    time.Now,
	}
}

// ─── Declarations ──────────────────────────────────────────────────────

func (s *taxService) CreateDeclaration(ctx context.Context, d *domain.TaxDeclaration) error {
	if err := d.Validate(); err != nil {
		return err
	}
	return s.repo.CreateDeclaration(ctx, d)
}

func (s *taxService) GetDeclaration(ctx context.Context, id string) (*domain.TaxDeclaration, error) {
	return s.repo.GetDeclarationByID(ctx, id)
}

func (s *taxService) ListDeclarations(ctx context.Context, filter domain.TaxDeclarationFilter) ([]domain.TaxDeclaration, error) {
	return s.repo.GetDeclarations(ctx, filter)
}

func (s *taxService) SubmitDeclaration(ctx context.Context, id, userID string) error {
	d, err := s.repo.GetDeclarationByID(ctx, id)
	if err != nil {
		return err
	}
	if !d.CanSubmit() {
		return domain.ErrDeclarationNotEditable
	}
	d.Status = domain.DeclStatusSUBMITTED
	d.SubmittedAt = s.now().Format(time.RFC3339)
	d.SubmittedBy = userID
	return s.repo.UpdateDeclaration(ctx, d)
}

func (s *taxService) AcknowledgeDeclaration(ctx context.Context, id, ref string) error {
	d, err := s.repo.GetDeclarationByID(ctx, id)
	if err != nil {
		return err
	}
	if d.Status != domain.DeclStatusSUBMITTED {
		return domain.ErrDeclarationNotEditable
	}
	d.Status = domain.DeclStatusACKNOWLEDGED
	d.AcknowledgedAt = s.now().Format(time.RFC3339)
	d.AcknowledgementRef = ref
	if err := s.repo.UpdateDeclaration(ctx, d); err != nil {
		return err
	}
	// A6: Declaration→payment automation. Acknowledged payable declarations
	// auto-create a PENDING TaxPayment with statutory due date.
	return s.createPaymentForDeclaration(ctx, d)
}

func (s *taxService) RejectDeclaration(ctx context.Context, id, reason string) error {
	d, err := s.repo.GetDeclarationByID(ctx, id)
	if err != nil {
		return err
	}
	if d.Status != domain.DeclStatusSUBMITTED && d.Status != domain.DeclStatusVALIDATED {
		return domain.ErrDeclarationNotEditable
	}
	d.Status = domain.DeclStatusREJECTED
	return s.repo.UpdateDeclaration(ctx, d)
}

func (s *taxService) AmendDeclaration(ctx context.Context, id string, lines []domain.TaxDeclarationLine) (*domain.TaxDeclaration, error) {
	d, err := s.repo.GetDeclarationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !d.CanAmend() {
		return nil, domain.ErrDeclarationNotEditable
	}
	amended := &domain.TaxDeclaration{
		CompanyID:       d.CompanyID,
		DeclarationType: d.DeclarationType,
		TaxPeriod:       d.TaxPeriod,
		Status:          domain.DeclStatusDRAFT,
		AdjustmentType:  domain.AdjTypeAMENDMENT,
		PreviousDeclID:  d.ID,
		Lines:           lines,
		CreatedBy:       d.CreatedBy,
		Version:         1,
	}
	if err := amended.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.CreateDeclaration(ctx, amended); err != nil {
		return nil, err
	}
	d.Status = domain.DeclStatusAMENDED
	_ = s.repo.UpdateDeclaration(ctx, d)
	return amended, nil
}

func (s *taxService) CancelDeclaration(ctx context.Context, id string) error {
	d, err := s.repo.GetDeclarationByID(ctx, id)
	if err != nil {
		return err
	}
	if d.Status.IsTerminal() {
		return domain.ErrDeclarationNotEditable
	}
	d.Status = domain.DeclStatusCANCELLED
	return s.repo.UpdateDeclaration(ctx, d)
}

// ─── Tax Rates ─────────────────────────────────────────────────────────

func (s *taxService) CreateRate(ctx context.Context, rate *domain.TaxRate) error {
	if rate.RateCode == "" {
		return domain.ErrTaxRateCodeExists
	}
	if rate.EffectiveFrom == "" {
		return domain.ErrValidationFailed
	}
	return s.repo.CreateRate(ctx, rate)
}

func (s *taxService) GetRate(ctx context.Context, id string) (*domain.TaxRate, error) {
	return s.repo.GetRateByID(ctx, id)
}

func (s *taxService) ListRates(ctx context.Context, filter domain.TaxRateFilter) ([]domain.TaxRate, error) {
	return s.repo.GetRates(ctx, filter)
}

func (s *taxService) UpdateRate(ctx context.Context, rate *domain.TaxRate) error {
	return s.repo.UpdateRate(ctx, rate)
}

// ─── Payments ──────────────────────────────────────────────────────────

func (s *taxService) CreatePayment(ctx context.Context, p *domain.TaxPayment) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return s.repo.CreatePayment(ctx, p)
}

func (s *taxService) GetPayment(ctx context.Context, id string) (*domain.TaxPayment, error) {
	return s.repo.GetPaymentByID(ctx, id)
}

func (s *taxService) ListPayments(ctx context.Context, filter domain.PaymentFilter) ([]domain.TaxPayment, error) {
	return s.repo.GetPayments(ctx, filter)
}

func (s *taxService) RecordPayment(ctx context.Context, id string, amount float64, date, ref string) error {
	p, err := s.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return err
	}
	p.PaidAmount = amount
	p.PaymentDate = date
	p.PaymentRef = ref
	p.Status = domain.PayStatusPAID
	if p.PaidAmount < p.DeclaredAmount {
		p.Status = domain.PayStatusPARTIAL
	}
	p.CalculateLateInterest()
	return s.repo.UpdatePayment(ctx, p)
}

// ─── E-Invoices ────────────────────────────────────────────────────────

func (s *taxService) CreateEInvoice(ctx context.Context, inv *domain.EInvoice) error {
	if err := inv.Validate(); err != nil {
		return err
	}
	return s.repo.CreateEInvoice(ctx, inv)
}

func (s *taxService) GetEInvoice(ctx context.Context, id string) (*domain.EInvoice, error) {
	return s.repo.GetEInvoiceByID(ctx, id)
}

func (s *taxService) ListEInvoices(ctx context.Context, filter domain.EInvoiceFilter) ([]domain.EInvoice, error) {
	return s.repo.GetEInvoices(ctx, filter)
}

// IssueEInvoice runs the issuance pipeline: generate TXML from the invoice,
// sign it (xmldsig), submit to GDT, record the transaction. Transitions
// DRAFT/VALIDATED → SIGNED → SUBMITTED. Requires a wired signer + GDT
// client, else ErrGDTUnavailable.
func (s *taxService) IssueEInvoice(ctx context.Context, id string) error {
	inv, err := s.repo.GetEInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status != domain.EInvStatusDRAFT && inv.Status != domain.EInvStatusVALIDATED {
		return domain.ErrInvoiceStatusInvalid
	}
	if s.signer == nil || s.gdt == nil {
		return domain.ErrGDTUnavailable
	}
	body, err := einvoice.GenerateTXML(inv)
	if err != nil {
		return fmt.Errorf("issue e-invoice: %w", err)
	}
	inv.XMLBody = string(body)
	signed, err := s.signer.SignTXML(inv.XMLBody, inv.DigitalSignatureID)
	if err != nil {
		return err
	}
	inv.SignedXML = signed
	inv.SigningDate = s.now().Format(time.RFC3339)
	inv.Status = domain.EInvStatusSIGNED
	if err := s.repo.UpdateEInvoice(ctx, inv); err != nil {
		return err
	}
	resp, err := s.gdt.SubmitInvoice(ctx, signed, inv.DigitalSignatureID)
	if err != nil {
		return mapGDTErr(err)
	}
	inv.GDTTransactionID = resp.TransactionID
	inv.Status = domain.EInvStatusSUBMITTED
	return s.repo.UpdateEInvoice(ctx, inv)
}

// CheckInvoiceStatus polls GDT for the invoice's current status and
// advances the local lifecycle: ACKNOWLEDGED → ISSUED, REJECTED → VALIDATED
// (editable, re-issuable).
func (s *taxService) CheckInvoiceStatus(ctx context.Context, id string) error {
	inv, err := s.repo.GetEInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status != domain.EInvStatusSUBMITTED {
		return domain.ErrInvoiceStatusInvalid
	}
	if s.gdt == nil {
		return domain.ErrGDTUnavailable
	}
	st, err := s.gdt.GetInvoiceStatus(ctx, inv.GDTTransactionID)
	if err != nil {
		return mapGDTErr(err)
	}
	switch st.Status {
	case "ACKNOWLEDGED":
		inv.Status = domain.EInvStatusISSUED
	case "REJECTED", "INVALID":
		inv.Status = domain.EInvStatusVALIDATED
	default:
		return nil // still processing — stay SUBMITTED
	}
	return s.repo.UpdateEInvoice(ctx, inv)
}

func mapGDTErr(err error) error {
	switch {
	case errors.Is(err, gdt.ErrUpstream), errors.Is(err, gdt.ErrNotFound):
		return domain.ErrGDTUnavailable
	case errors.Is(err, gdt.ErrUnauthorized):
		return domain.ErrGDTUnauthorized
	case errors.Is(err, gdt.ErrInvalidRequest):
		return domain.ErrGDTRejected
	default:
		return err
	}
}

// pemSigner signs canonicalized TXML with an RSA private key, embedding the
// signature into BK:ChuKySo before BK:KyThuat closes (TAX_TEMPLATES §5).
type pemSigner struct {
	key    *rsa.PrivateKey
	serial string
	now    func() time.Time
}

// NewPEMSigner builds a TXMLSigner from an RSA private key. serial is the
// digital certificate serial placed in BK:SerialNumber.
func NewPEMSigner(key *rsa.PrivateKey, serial string, now func() time.Time) TXMLSigner {
	return &pemSigner{key: key, serial: serial, now: now}
}

func (p *pemSigner) SignTXML(xmlBody, signatureID string) (string, error) {
	canon, err := xmldsig.Canonicalize([]byte(xmlBody))
	if err != nil {
		return "", fmt.Errorf("sign TXML: %w", err)
	}
	sig, err := xmldsig.SignBase64(p.key, canon)
	if err != nil {
		return "", fmt.Errorf("sign TXML: %w", err)
	}
	stamp := p.now().Format("2006-01-02T15:04:05-07:00")
	chuku := "<BK:ChuKySo><BK:SerialNumber>" + p.serial +
		"</BK:SerialNumber><BK:ThoiDiemKy>" + stamp +
		"</BK:ThoiDiemKy><BK:DuLieuKy>" + sig +
		"</BK:DuLieuKy></BK:ChuKySo>"
	signed := strings.Replace(string(canon), "</BK:KyThuat>", chuku+"</BK:KyThuat>", 1)
	if signed == string(canon) {
		return "", fmt.Errorf("sign TXML: BK:KyThuat not found")
	}
	return signed, nil
}

func (s *taxService) CancelEInvoice(ctx context.Context, id, reason string) error {
	inv, err := s.repo.GetEInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	if !inv.Status.CanCancel() {
		return domain.ErrInvoiceStatusInvalid
	}
	if inv.GDTTransactionID != "" && s.gdt != nil {
		if err := s.gdt.CancelInvoice(ctx, inv.GDTTransactionID, reason); err != nil {
			return mapGDTErr(err)
		}
	}
	inv.Status = domain.EInvStatusCANCELLED
	inv.CancelledAt = s.now().Format(time.RFC3339)
	inv.CancelReason = reason
	return s.repo.UpdateEInvoice(ctx, inv)
}

// ─── Tax Calendar ──────────────────────────────────────────────────────

func (s *taxService) CreateCalendarEntry(ctx context.Context, c *domain.TaxCalendar) error {
	if c.CompanyID == "" {
		return domain.ErrCompanyIDRequired
	}
	if c.DeclarationDue == "" {
		return domain.ErrValidationFailed
	}
	if c.Status == "" {
		c.Status = domain.CalStatusPENDING
	}
	return s.repo.CreateCalendarEntry(ctx, c)
}

func (s *taxService) GetCalendarEntry(ctx context.Context, id string) (*domain.TaxCalendar, error) {
	return s.repo.GetCalendarEntryByID(ctx, id)
}

func (s *taxService) GetCalendarByPeriod(ctx context.Context, companyID string, periodYear, periodNumber int) ([]domain.TaxCalendar, error) {
	return s.repo.GetCalendarByPeriod(ctx, companyID, periodYear, periodNumber)
}

func (s *taxService) GetCalendarByCompany(ctx context.Context, companyID string) ([]domain.TaxCalendar, error) {
	return s.repo.GetCalendarByCompany(ctx, companyID)
}

// ─── Alerts ────────────────────────────────────────────────────────────

func (s *taxService) CreateAlert(ctx context.Context, a *domain.TaxAlert) error {
	if a.CompanyID == "" {
		return domain.ErrCompanyIDRequired
	}
	if a.Message == "" {
		return domain.ErrValidationFailed
	}
	return s.repo.CreateAlert(ctx, a)
}

func (s *taxService) GetAlert(ctx context.Context, id string) (*domain.TaxAlert, error) {
	return s.repo.GetAlertByID(ctx, id)
}

func (s *taxService) ListAlerts(ctx context.Context, companyID string, limit int) ([]domain.TaxAlert, error) {
	return s.repo.GetAlerts(ctx, companyID, limit)
}

// ─── Audit Cases ───────────────────────────────────────────────────────

func (s *taxService) CreateAuditCase(ctx context.Context, a *domain.TaxAuditCase) error {
	if a.CompanyID == "" {
		return domain.ErrCompanyIDRequired
	}
	if a.AuditPeriodStart == "" {
		return domain.ErrValidationFailed
	}
	if a.Status == "" {
		a.Status = domain.AuditCaseOPEN
	}
	return s.repo.CreateAuditCase(ctx, a)
}

func (s *taxService) GetAuditCase(ctx context.Context, id string) (*domain.TaxAuditCase, error) {
	return s.repo.GetAuditCaseByID(ctx, id)
}

func (s *taxService) ListAuditCases(ctx context.Context, companyID string) ([]domain.TaxAuditCase, error) {
	return s.repo.GetAuditCases(ctx, companyID)
}

func (s *taxService) CloseAuditCase(ctx context.Context, id string, findings string, penalty float64) error {
	a, err := s.repo.GetAuditCaseByID(ctx, id)
	if err != nil {
		return err
	}
	a.Status = domain.AuditCaseCLOSED
	a.Findings = findings
	a.PenaltyAmount = penalty
	a.ClosedAt = s.now().Format(time.RFC3339)
	return s.repo.UpdateAuditCase(ctx, a)
}

// ─── Calculations ──────────────────────────────────────────────────────

// VAT account classification. Explicit VAT accounts (3331x output, 133x input)
// carry the bookkeeper-posted amounts and win over rate-based estimation.
var vatOutputAccounts = []string{"33311", "33312", "33313"}
var vatInputAccounts = []string{"1331", "1332"}
var vatRevenueAccounts = []string{"5111", "5112", "5113", "711"}
var vatExpenseAccounts = []string{"152", "153", "156", "611", "621", "627", "641", "642"}

func (s *taxService) CalculateVAT(ctx context.Context, companyID string, period domain.TaxPeriod, entries []domain.JournalEntry) (*domain.VATResult, error) {
	if companyID == "" {
		return nil, domain.ErrCompanyIDRequired
	}
	standard, err := s.resolveRate(ctx, domain.TaxTypeVAT, "STANDARD", periodEndDate(period))
	if err != nil {
		return nil, err
	}
	result := &domain.VATResult{
		CompanyID: companyID,
		Period:    period,
	}
	explicitOutput := false
	explicitInput := false
	for _, entry := range entries {
		for _, line := range entry.Lines {
			if line.CreditAmount > 0 {
				if contains(vatOutputAccounts, line.AccountCode) {
					result.OutputVAT += line.CreditAmount
					explicitOutput = true
				}
				if contains(vatRevenueAccounts, line.AccountCode) {
					result.SalesTotal += line.CreditAmount
				}
			}
			if line.DebitAmount > 0 {
				if contains(vatInputAccounts, line.AccountCode) {
					result.InputVAT += line.DebitAmount
					explicitInput = true
				}
				if contains(vatExpenseAccounts, line.AccountCode) || isFixedAssetAccount(line.AccountCode) {
					result.PurchaseTotal += line.DebitAmount
				}
			}
		}
	}
	if !explicitOutput {
		result.OutputVAT = result.SalesTotal * standard.RateValue / 100
	}
	if !explicitInput {
		rate := standard.RateValue / 100
		var inputBase, faBase float64
		for _, entry := range entries {
			for _, line := range entry.Lines {
				if line.DebitAmount <= 0 {
					continue
				}
				if contains(vatExpenseAccounts, line.AccountCode) {
					inputBase += line.DebitAmount
				}
				if isFixedAssetAccount(line.AccountCode) {
					faBase += line.DebitAmount
				}
			}
		}
		result.InputVAT = inputBase * rate
		result.InputVATFA = faBase * rate
	}
	result.TotalInputVAT = result.InputVAT + result.InputVATFA
	if result.OutputVAT > result.TotalInputVAT {
		result.VATPayable = round2(result.OutputVAT - result.TotalInputVAT)
	} else {
		result.VATRefundable = round2(result.TotalInputVAT - result.OutputVAT)
	}
	result.OutputVAT = round2(result.OutputVAT)
	result.InputVAT = round2(result.InputVAT)
	result.InputVATFA = round2(result.InputVATFA)
	result.TotalInputVAT = round2(result.TotalInputVAT)
	result.SalesTotal = round2(result.SalesTotal)
	result.PurchaseTotal = round2(result.PurchaseTotal)
	return result, nil
}

func isFixedAssetAccount(code string) bool {
	return len(code) >= 3 && code[:3] == "211"
}

// periodEndDate returns the ISO end date of the tax period (used for rate
// resolution). Monthly/quarterly/annual mapping; per-occurrence defaults to
// December of the period year.
func periodEndDate(p domain.TaxPeriod) string {
	n := p.PeriodNumber
	switch p.PeriodType {
	case domain.PeriodTypeQuarterly:
		n *= 3
	case domain.PeriodTypeAnnual:
		n = 12
	case domain.PeriodTypePerOccurrence:
		n = 12
	}
	if n < 1 {
		n = 1
	}
	if n > 12 {
		n = 12
	}
	date := time.Date(p.PeriodYear, time.Month(n)+1, 0, 0, 0, 0, 0, time.UTC)
	return date.Format("2006-01-02")
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

// CIT size thresholds per CIT Law 67/2025/QH15 Art. 11(2): MICRO < 3B,
// SMALL 3-50B, STANDARD >= 50B annual revenue. Size derived from the
// revenue in the provided entries (company master lookup is future work).
const (
	citMicroRevenueLimit = 3e9
	citSmallRevenueLimit = 50e9
)

func citSizeKey(revenue float64) string {
	switch {
	case revenue < citMicroRevenueLimit:
		return "MICRO"
	case revenue < citSmallRevenueLimit:
		return "SMALL"
	default:
		return "STANDARD"
	}
}

func (s *taxService) CalculateCIT(ctx context.Context, companyID string, year int, entries []domain.JournalEntry) (*domain.CITResult, error) {
	if companyID == "" {
		return nil, domain.ErrCompanyIDRequired
	}
	result := &domain.CITResult{
		CompanyID:  companyID,
		PeriodYear: year,
		PeriodType: "ANNUAL",
	}
	for _, entry := range entries {
		for _, line := range entry.Lines {
			acct := line.AccountCode
			if len(acct) >= 3 {
				switch acct[0] {
				case '5', '6', '7':
					if acct[1] >= '1' && acct[1] <= '9' {
						if acct[0] == '5' || acct[0] == '7' {
							result.Revenue += line.CreditAmount
						} else {
							result.Expenses += line.DebitAmount
						}
					}
				case '8':
					if acct[1] == '2' {
						result.NonDeductible += line.DebitAmount
					}
				}
			}
		}
	}
	result.TaxableIncome = result.Revenue - result.Expenses + result.NonDeductible
	if result.TaxableIncome < 0 {
		result.TaxableIncome = 0
	}
	rate, err := s.resolveRate(ctx, domain.TaxTypeCIT, citSizeKey(result.Revenue), fmt.Sprintf("%04d-12-31", year))
	if err != nil {
		return nil, err
	}
	result.TaxRate = rate.RateValue
	result.CITPayable = round2(result.TaxableIncome * result.TaxRate / 100)
	result.CITFinal = result.CITPayable
	return result, nil
}

// PIT progressive brackets per PIT Law Art. 7 (monthly taxable income, VND).
// tax = taxableIncome * rate% - reduction.
var pitBrackets = []struct {
	upper     float64
	rate      float64
	reduction float64
}{
	{5e6, 5, 0},
	{10e6, 10, 250000},
	{18e6, 15, 750000},
	{32e6, 20, 1650000},
	{52e6, 25, 3250000},
	{80e6, 30, 5850000},
	{math.Inf(1), 35, 9850000},
}

const (
	pitPersonalDeduction  = 11e6  // PIT-03: 11M/month resident
	pitDependantDeduction = 4.4e6 // PIT-04: 4.4M/month per dependant
	pitSocialRate         = 0.08  // PIT-05
	pitHealthRate         = 0.015 // PIT-06
	pitUnemploymentRate   = 0.01  // PIT-07
)

func progressivePIT(monthlyTaxable float64) float64 {
	if monthlyTaxable <= 0 {
		return 0
	}
	for _, b := range pitBrackets {
		if monthlyTaxable < b.upper {
			return round2(monthlyTaxable*b.rate/100 - b.reduction)
		}
	}
	return round2(monthlyTaxable*35/100 - 9850000)
}

func (s *taxService) CalculatePIT(ctx context.Context, companyID string, period domain.TaxPeriod, employees []domain.PITEmployeeInput) (*domain.PITResult, error) {
	if companyID == "" {
		return nil, domain.ErrCompanyIDRequired
	}
	nonResidentRate, err := s.resolveRate(ctx, domain.TaxTypePIT, "STANDARD", periodEndDate(period))
	if err != nil {
		return nil, err
	}
	result := &domain.PITResult{CompanyID: companyID, Period: period}
	for _, e := range employees {
		months := e.Months
		if months < 1 {
			months = 1
		}
		gross := e.GrossMonthly * float64(months)
		result.TotalGross += gross
		result.EmployeeCount++
		if e.NonResident {
			result.TotalPIT += gross * nonResidentRate.RateValue / 100
			continue
		}
		insuranceBase := e.SocialInsuranceBase
		if insuranceBase <= 0 {
			insuranceBase = e.GrossMonthly
		}
		insurance := (pitSocialRate + pitHealthRate + pitUnemploymentRate) * insuranceBase
		deductions := insurance + pitPersonalDeduction + pitDependantDeduction*float64(e.Dependants)
		result.TotalDeductions += deductions * float64(months)
		result.TotalPIT += progressivePIT(e.GrossMonthly-deductions) * float64(months)
	}
	result.TotalGross = round2(result.TotalGross)
	result.TotalDeductions = round2(result.TotalDeductions)
	result.TotalPIT = round2(result.TotalPIT)
	return result, nil
}

// ─── A5: Declaration Engine ─────────────────────────────────────────────
// Generates GTGT01 / TNDN03 declarations from posted journal entries.
// Unsupported types (payroll-dependent, e.g. KK_TNCN) are rejected.

func (s *taxService) GenerateDeclaration(ctx context.Context, companyID string, declType domain.DeclarationType, period domain.TaxPeriod, userID string) (*domain.TaxDeclaration, error) {
	if declType != domain.DeclTypeGTGT01 && declType != domain.DeclTypeTNDN03 {
		return nil, domain.ErrDeclarationTypeInvalid
	}
	existing, err := s.repo.GetDeclarations(ctx, domain.TaxDeclarationFilter{
		CompanyID: companyID, DeclarationType: declType,
		PeriodYear: period.PeriodYear, PeriodNumber: period.PeriodNumber,
	})
	if err != nil {
		return nil, err
	}
	for _, d := range existing {
		if d.TaxPeriod.PeriodType == period.PeriodType && d.Status != domain.DeclStatusCANCELLED {
			return nil, domain.ErrDuplicateDeclaration
		}
	}
	from, to := periodDateRange(period)
	journals, err := s.jeRepo.GetByDateRange(ctx, from, to)
	if err != nil {
		return nil, err
	}
	var posted []domain.JournalEntry
	for _, je := range journals {
		if je.CompanyID == companyID && je.Status == domain.JournalEntryPosted {
			posted = append(posted, je)
		}
	}
	decl := &domain.TaxDeclaration{
		CompanyID:       companyID,
		DeclarationType: declType,
		TaxPeriod:       period,
		Status:          domain.DeclStatusVALIDATED,
		AdjustmentType:  domain.AdjTypeNONE,
		Version:         1,
		CreatedBy:       userID,
		CreatedAt:       s.now().Format(time.RFC3339),
	}
	switch declType {
	case domain.DeclTypeGTGT01:
		res, err := s.CalculateVAT(ctx, companyID, period, posted)
		if err != nil {
			return nil, err
		}
		decl.Lines = vatDeclarationLines(res, posted)
	case domain.DeclTypeTNDN03:
		res, err := s.CalculateCIT(ctx, companyID, period.PeriodYear, posted)
		if err != nil {
			return nil, err
		}
		decl.Lines = citDeclarationLines(res, posted)
	}
	if err := validateDeclarationRules(decl); err != nil {
		return nil, err
	}
	if err := decl.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.CreateDeclaration(ctx, decl); err != nil {
		return nil, err
	}
	return decl, nil
}

func periodDateRange(p domain.TaxPeriod) (from, to time.Time) {
	switch p.PeriodType {
	case domain.PeriodTypeQuarterly:
		startMonth := time.Month((p.PeriodNumber-1)*3 + 1)
		from = time.Date(p.PeriodYear, startMonth, 1, 0, 0, 0, 0, time.UTC)
		to = time.Date(p.PeriodYear, startMonth+2, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
	case domain.PeriodTypeAnnual:
		from = time.Date(p.PeriodYear, 1, 1, 0, 0, 0, 0, time.UTC)
		to = time.Date(p.PeriodYear, 12, 31, 23, 59, 59, 0, time.UTC)
	default: // monthly
		from = time.Date(p.PeriodYear, time.Month(p.PeriodNumber), 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(0, 1, -1)
	}
	return from, to
}

func entryIDs(entries []domain.JournalEntry) []string {
	ids := make([]string, 0, len(entries))
	for _, je := range entries {
		ids = append(ids, je.ID)
	}
	return ids
}

// GTGT01 lines per TAX_RULES §2.1: [16]=[14]+[15], [23]=[21]+[22],
// [30]=[23]-[16], XOR [31]/[32].
func vatDeclarationLines(res *domain.VATResult, entries []domain.JournalEntry) []domain.TaxDeclarationLine {
	payable, refundable := 0.0, 0.0
	if res.VATPayable > 0 {
		payable = res.VATPayable
	} else {
		refundable = res.VATRefundable
	}
	src := domain.SrcTypeFROM_LEDGER
	ids := entryIDs(entries)
	return []domain.TaxDeclarationLine{
		{LineCode: "14", LineName: "VAT đầu vào được khấu trừ (hàng hóa, dịch vụ)", Amount: res.InputVAT, SourceType: src, SourceEntryIDs: ids, SortOrder: 14},
		{LineCode: "15", LineName: "VAT đầu vào được khấu trừ (TSCĐ)", Amount: res.InputVATFA, SourceType: src, SourceEntryIDs: ids, SortOrder: 15},
		{LineCode: "16", LineName: "Tổng VAT đầu vào được khấu trừ", Amount: res.TotalInputVAT, SourceType: src, SourceEntryIDs: ids, SortOrder: 16},
		{LineCode: "21", LineName: "VAT đầu ra (nội địa)", Amount: res.OutputVAT, SourceType: src, SourceEntryIDs: ids, SortOrder: 21},
		{LineCode: "22", LineName: "VAT đầu ra (xuất khẩu)", Amount: 0, SourceType: src, SourceEntryIDs: ids, SortOrder: 22},
		{LineCode: "23", LineName: "Tổng VAT đầu ra", Amount: res.OutputVAT, SourceType: src, SourceEntryIDs: ids, SortOrder: 23},
		{LineCode: "30", LineName: "VAT phải nộp/đề nghị hoàn ([23]-[16])", Amount: res.OutputVAT - res.TotalInputVAT, SourceType: src, SourceEntryIDs: ids, SortOrder: 30},
		{LineCode: "31", LineName: "Thuế GTGT phải nộp", Amount: payable, SourceType: src, SourceEntryIDs: ids, SortOrder: 31},
		{LineCode: "32", LineName: "Thuế GTGT đề nghị hoàn", Amount: refundable, SourceType: src, SourceEntryIDs: ids, SortOrder: 32},
	}
}

// TNDN03 lines per TAX_RULES §2.1: [04]>=[06], [14]<=[12]*[13].
func citDeclarationLines(res *domain.CITResult, entries []domain.JournalEntry) []domain.TaxDeclarationLine {
	src := domain.SrcTypeFROM_LEDGER
	ids := entryIDs(entries)
	return []domain.TaxDeclarationLine{
		{LineCode: "04", LineName: "Tổng doanh thu", Amount: res.Revenue, SourceType: src, SourceEntryIDs: ids, SortOrder: 4},
		{LineCode: "06", LineName: "Doanh thu tính thuế", Amount: res.Revenue, SourceType: src, SourceEntryIDs: ids, SortOrder: 6},
		{LineCode: "12", LineName: "Thu nhập chịu thuế", Amount: res.TaxableIncome, SourceType: src, SourceEntryIDs: ids, SortOrder: 12},
		{LineCode: "13", LineName: "Thuế suất (%)", Amount: res.TaxRate, SourceType: src, SourceEntryIDs: ids, SortOrder: 13},
		{LineCode: "14", LineName: "Thuế TNDN phải nộp", Amount: res.CITPayable, SourceType: src, SourceEntryIDs: ids, SortOrder: 14},
	}
}

func lineAmount(lines []domain.TaxDeclarationLine, code string) float64 {
	for _, l := range lines {
		if l.LineCode == code {
			return l.Amount
		}
	}
	return 0
}

const vatEpsilon = 1.0 // 1 VND rounding tolerance

func validateDeclarationRules(decl *domain.TaxDeclaration) error {
	switch decl.DeclarationType {
	case domain.DeclTypeGTGT01:
		// [16], [23], [30] are algebraic sums — [30] may be negative when
		// input VAT exceeds output VAT (refundable direction).
		if lineAmount(decl.Lines, "16") != lineAmount(decl.Lines, "14")+lineAmount(decl.Lines, "15") {
			return domain.ErrValidationFailed
		}
		if lineAmount(decl.Lines, "23") != lineAmount(decl.Lines, "21")+lineAmount(decl.Lines, "22") {
			return domain.ErrValidationFailed
		}
		if lineAmount(decl.Lines, "30") != lineAmount(decl.Lines, "23")-lineAmount(decl.Lines, "16") {
			return domain.ErrValidationFailed
		}
		if lineAmount(decl.Lines, "31") < 0 || lineAmount(decl.Lines, "32") < 0 {
			return domain.ErrValidationFailed
		}
		if lineAmount(decl.Lines, "31") > 0 && lineAmount(decl.Lines, "32") > 0 {
			return domain.ErrValidationFailed
		}
	case domain.DeclTypeTNDN03:
		for _, l := range decl.Lines {
			if l.Amount < 0 {
				return domain.ErrValidationFailed
			}
		}
		if lineAmount(decl.Lines, "04") < lineAmount(decl.Lines, "06") {
			return domain.ErrValidationFailed
		}
		if lineAmount(decl.Lines, "14") > lineAmount(decl.Lines, "12")*lineAmount(decl.Lines, "13")/100+vatEpsilon {
			return domain.ErrValidationFailed
		}
	}
	return nil
}

// ─── A6: Payment Automation ─────────────────────────────────────────────

func (s *taxService) createPaymentForDeclaration(ctx context.Context, d *domain.TaxDeclaration) error {
	payable := lineAmount(d.Lines, "31")
	if d.DeclarationType == domain.DeclTypeTNDN03 {
		payable = lineAmount(d.Lines, "14")
	}
	if payable <= 0 {
		return nil // refundable or zero declaration — no payable
	}
	existing, err := s.repo.GetPayments(ctx, domain.PaymentFilter{
		CompanyID: d.CompanyID, TaxType: declarationTaxType(d.DeclarationType),
		PeriodYear: d.TaxPeriod.PeriodYear, PeriodNumber: d.TaxPeriod.PeriodNumber,
	})
	if err != nil {
		return err
	}
	for _, p := range existing {
		if p.DeclarationID == d.ID {
			return nil // already created
		}
	}
	p := &domain.TaxPayment{
		CompanyID:      d.CompanyID,
		DeclarationID:  d.ID,
		TaxType:        declarationTaxType(d.DeclarationType),
		PeriodYear:     d.TaxPeriod.PeriodYear,
		PeriodNumber:   d.TaxPeriod.PeriodNumber,
		DeclaredAmount: payable,
		DueDate:        paymentDueDate(d),
		Status:         domain.PayStatusPENDING,
		Notes:          "auto-generated on declaration acknowledgement",
		CreatedAt:      s.now().Format(time.RFC3339),
	}
	return s.repo.CreatePayment(ctx, p)
}

func declarationTaxType(dt domain.DeclarationType) domain.TaxType {
	switch {
	case dt == domain.DeclTypeGTGT01 || dt == domain.DeclTypeGTGT02 ||
		dt == domain.DeclTypeGTGT03 || dt == domain.DeclTypeGTGT04 ||
		dt == domain.DeclTypeGTGT05:
		return domain.TaxTypeVAT
	case dt == domain.DeclTypeTNDN03 || dt == domain.DeclTypeTNDN04 ||
		dt == domain.DeclTypeTNDN02 || dt == domain.DeclTypeTNDN05 ||
		dt == domain.DeclTypeTNDN06:
		return domain.TaxTypeCIT
	case dt == domain.DeclTypeKKTNCN || dt == domain.DeclTypeQTTTNCN:
		return domain.TaxTypePIT
	}
	return domain.TaxTypeVAT
}

// Statutory deadlines: VAT monthly = 20th next month (VAT-12); VAT/CIT
// quarterly = 30th of month after quarter end (VAT-12, CIT-08); CIT annual
// = 31-Mar next year (CIT-09).
func paymentDueDate(d *domain.TaxDeclaration) string {
	y, n := d.TaxPeriod.PeriodYear, d.TaxPeriod.PeriodNumber
	switch d.TaxPeriod.PeriodType {
	case domain.PeriodTypeQuarterly:
		endMonth := time.Month(n * 3)
		return time.Date(y, endMonth+1, 30, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	case domain.PeriodTypeAnnual:
		return time.Date(y+1, 3, 31, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	default: // monthly
		return time.Date(y, time.Month(n)+1, 20, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}
}
