package service

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	"gotax/internal/domain"
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
	CalculatePIT(ctx context.Context, companyID string, period domain.TaxPeriod, employeeIDs []string) (*domain.PITResult, error)
}

type taxService struct {
	repo domain.TaxRepository
	now  func() time.Time
}

func NewTaxService(repo domain.TaxRepository) TaxServiceInterface {
	return &taxService{
		repo: repo,
		now:  time.Now,
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
	return s.repo.UpdateDeclaration(ctx, d)
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

func (s *taxService) IssueEInvoice(ctx context.Context, id string) error {
	inv, err := s.repo.GetEInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status != domain.EInvStatusDRAFT && inv.Status != domain.EInvStatusVALIDATED {
		return domain.ErrInvoiceStatusInvalid
	}
	inv.Status = domain.EInvStatusISSUED
	inv.IssueDate = s.now().Format("2006-01-02")
	return s.repo.UpdateEInvoice(ctx, inv)
}

func (s *taxService) CancelEInvoice(ctx context.Context, id, reason string) error {
	inv, err := s.repo.GetEInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	if !inv.Status.CanCancel() {
		return domain.ErrInvoiceStatusInvalid
	}
	inv.Status = domain.EInvStatusCANCELLED
	inv.CancelledAt = s.now().Format(time.RFC3339)
	inv.CancelReason = reason
	return s.repo.UpdateEInvoice(ctx, inv)
}

// ─── Tax Calendar ──────────────────────────────────────────────────────

func (s *taxService) CreateCalendarEntry(ctx context.Context, c *domain.TaxCalendar) error {
	if c.CompanyID == "" { return domain.ErrCompanyIDRequired }
	if c.DeclarationDue == "" { return domain.ErrValidationFailed }
	if c.Status == "" { c.Status = domain.CalStatusPENDING }
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
	if a.CompanyID == "" { return domain.ErrCompanyIDRequired }
	if a.Message == "" { return domain.ErrValidationFailed }
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
	if a.CompanyID == "" { return domain.ErrCompanyIDRequired }
	if a.AuditPeriodStart == "" { return domain.ErrValidationFailed }
	if a.Status == "" { a.Status = domain.AuditCaseOPEN }
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

var vatOutputAccounts = []string{"5111", "5112", "5113", "711"}
var vatInputAccounts = []string{"152", "153", "156", "611", "621", "627", "641", "642"}
var vatOutputRate = 0.10
var vatInputRate = 0.10
var vatInputFARate = 0.10

func (s *taxService) CalculateVAT(_ context.Context, companyID string, period domain.TaxPeriod, entries []domain.JournalEntry) (*domain.VATResult, error) {
	if companyID == "" {
		return nil, domain.ErrCompanyIDRequired
	}
	result := &domain.VATResult{
		CompanyID: companyID,
		Period:    period,
	}
	for _, entry := range entries {
		for _, line := range entry.Lines {
			if line.CreditAmount > 0 {
				for _, acct := range vatOutputAccounts {
					if line.AccountCode == acct {
						result.OutputVAT += line.CreditAmount * vatOutputRate
						result.SalesTotal += line.CreditAmount
					}
				}
			}
			if line.DebitAmount > 0 {
				isInput := false
				for _, acct := range vatInputAccounts {
					if line.AccountCode == acct {
						result.InputVAT += line.DebitAmount * vatInputRate
						result.PurchaseTotal += line.DebitAmount
						isInput = true
						break
					}
				}
				if !isInput && len(line.AccountCode) >= 3 && line.AccountCode[:3] == "211" {
					result.InputVATFA += line.DebitAmount * vatInputFARate
					result.PurchaseTotal += line.DebitAmount
				}
			}
		}
	}
	result.TotalInputVAT = result.InputVAT + result.InputVATFA
	if result.OutputVAT > result.TotalInputVAT {
		result.VATPayable = math.Round((result.OutputVAT - result.TotalInputVAT)*100) / 100
	} else {
		result.VATRefundable = math.Round((result.TotalInputVAT - result.OutputVAT)*100) / 100
	}
	result.OutputVAT = math.Round(result.OutputVAT*100) / 100
	result.InputVAT = math.Round(result.InputVAT*100) / 100
	result.TotalInputVAT = math.Round(result.TotalInputVAT*100) / 100
	return result, nil
}

var citRate = 0.20

func (s *taxService) CalculateCIT(_ context.Context, companyID string, year int, entries []domain.JournalEntry) (*domain.CITResult, error) {
	if companyID == "" {
		return nil, domain.ErrCompanyIDRequired
	}
	result := &domain.CITResult{
		CompanyID:  companyID,
		PeriodYear: year,
	}
	for _, entry := range entries {
		for _, line := range entry.Lines {
			acct := line.AccountCode
			if len(acct) >= 3 {
				switch acct[0] {
				case '5', '6', '7':
					if acct[1] >= '1' && acct[1] <= '9' {
						// revenue or expense
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
	result.CITPayable = math.Round(result.TaxableIncome*citRate*100) / 100
	result.CITFinal = result.CITPayable
	return result, nil
}

func (s *taxService) CalculatePIT(_ context.Context, _ string, _ domain.TaxPeriod, _ []string) (*domain.PITResult, error) {
	return &domain.PITResult{}, nil
}
