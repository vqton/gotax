package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"gotax/internal/domain"
)

var taxCodeRegex = regexp.MustCompile(`^\d{10}(-?\d{3})?$`)

type CompanyService interface {
	CreateCompany(ctx context.Context, c *domain.Company) error
	GetCompany(ctx context.Context, id string) (*domain.Company, error)
	GetCompanyByTaxCode(ctx context.Context, tenantID, taxCode string) (*domain.Company, error)
	ListCompanies(ctx context.Context, tenantID string) ([]domain.Company, error)
	UpdateCompany(ctx context.Context, c *domain.Company) error
	DeactivateCompany(ctx context.Context, id, reason string) error
	GetCompanyHierarchy(ctx context.Context, companyID string) ([]domain.Company, error)
	SwitchCompany(ctx context.Context, userID, companyID string) (*domain.CompanyContext, error)

	CreateBranch(ctx context.Context, b *domain.CompanyBranch) error
	GetBranch(ctx context.Context, id string) (*domain.CompanyBranch, error)
	ListBranches(ctx context.Context, companyID string) ([]domain.CompanyBranch, error)
	UpdateBranch(ctx context.Context, b *domain.CompanyBranch) error
	DeactivateBranch(ctx context.Context, id string) error

	CreateFiscalYear(ctx context.Context, fy *domain.FiscalYear) error
	GetFiscalYear(ctx context.Context, id string) (*domain.FiscalYear, error)
	ListFiscalYears(ctx context.Context, companyID string) ([]domain.FiscalYear, error)
	GeneratePeriods(ctx context.Context, fy *domain.FiscalYear) ([]domain.PeriodV2, error)
	ClosePeriod(ctx context.Context, companyID, periodID, userID string) error
	ReopenPeriod(ctx context.Context, companyID, periodID, userID string) error
	PermanentClosePeriod(ctx context.Context, companyID, periodID, userID string) error
	GetCurrentPeriod(ctx context.Context, companyID string) (*domain.PeriodV2, error)

	CreateDepartment(ctx context.Context, d *domain.Department) error
	GetDepartment(ctx context.Context, id string) (*domain.Department, error)
	ListDepartments(ctx context.Context, companyID string) ([]domain.Department, error)
	UpdateDepartment(ctx context.Context, d *domain.Department) error

	CreateEmployee(ctx context.Context, e *domain.Employee) error
	GetEmployee(ctx context.Context, id string) (*domain.Employee, error)
	ListEmployees(ctx context.Context, companyID string) ([]domain.Employee, error)
	UpdateEmployee(ctx context.Context, e *domain.Employee) error
	DeactivateEmployee(ctx context.Context, id string) error

	CreateBankAccount(ctx context.Context, ba *domain.CompanyBankAccount) error
	GetBankAccount(ctx context.Context, id string) (*domain.CompanyBankAccount, error)
	ListBankAccounts(ctx context.Context, companyID string) ([]domain.CompanyBankAccount, error)
	UpdateBankAccount(ctx context.Context, ba *domain.CompanyBankAccount) error
	DeactivateBankAccount(ctx context.Context, id string) error

	RegisterEInvoicePattern(ctx context.Context, inv *domain.EInvoicePattern) error
	GetEInvoicePattern(ctx context.Context, id string) (*domain.EInvoicePattern, error)
	ListEInvoicePatterns(ctx context.Context, companyID string) ([]domain.EInvoicePattern, error)

	RegisterDigitalSignature(ctx context.Context, sig *domain.DigitalSignature) error
	GetDigitalSignature(ctx context.Context, id string) (*domain.DigitalSignature, error)
	ListDigitalSignatures(ctx context.Context, companyID string) ([]domain.DigitalSignature, error)

	CreateIntegrationProfile(ctx context.Context, p *domain.IntegrationProfile) error
	GetIntegrationProfile(ctx context.Context, id string) (*domain.IntegrationProfile, error)
	ListIntegrationProfiles(ctx context.Context, companyID string) ([]domain.IntegrationProfile, error)
	UpdateIntegrationProfile(ctx context.Context, p *domain.IntegrationProfile) error
	TestIntegration(ctx context.Context, id string) error
}

type companyService struct {
	repo domain.CompanyRepository
	now  func() time.Time
}

func NewCompanyService(repo domain.CompanyRepository) CompanyService {
	return &companyService{repo: repo, now: time.Now}
}

func (s *companyService) CreateCompany(ctx context.Context, c *domain.Company) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if !taxCodeRegex.MatchString(c.TaxCode) {
		return domain.ErrInvalidTaxCode
	}
	return s.repo.Create(ctx, c)
}

func (s *companyService) GetCompany(ctx context.Context, id string) (*domain.Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *companyService) GetCompanyByTaxCode(ctx context.Context, tenantID, taxCode string) (*domain.Company, error) {
	if !taxCodeRegex.MatchString(taxCode) {
		return nil, domain.ErrInvalidTaxCode
	}
	return s.repo.GetByTaxCode(ctx, tenantID, taxCode)
}

func (s *companyService) ListCompanies(ctx context.Context, tenantID string) ([]domain.Company, error) {
	return s.repo.GetAll(ctx, tenantID)
}

func (s *companyService) UpdateCompany(ctx context.Context, c *domain.Company) error {
	existing, err := s.repo.GetByID(ctx, c.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.CompanyStatusActive {
		return domain.ErrCompanyInactive
	}
	c.UpdatedAt = s.now()
	if c.TenantID == "" {
		c.TenantID = existing.TenantID
	}
	return s.repo.Update(ctx, c)
}

func (s *companyService) DeactivateCompany(ctx context.Context, id, reason string) error {
	return s.repo.Deactivate(ctx, id, reason)
}

func (s *companyService) GetCompanyHierarchy(ctx context.Context, companyID string) ([]domain.Company, error) {
	return s.repo.GetHierarchy(ctx, companyID)
}

func (s *companyService) SwitchCompany(ctx context.Context, userID, companyID string) (*domain.CompanyContext, error) {
	c, err := s.repo.GetByID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	if c.Status != domain.CompanyStatusActive {
		return nil, domain.ErrCompanyInactive
	}
	openPeriod, _ := s.repo.GetOpenPeriod(ctx, companyID)

	ctxOut := &domain.CompanyContext{
		CompanyID:        c.ID,
		LegalNameVN:      c.LegalNameVN,
		TaxCode:          c.TaxCode,
		AccountingRegime: string(c.AccountingRegime),
		DefaultCurrency:  c.DefaultCurrency,
		Permissions:      []string{"company:read", "company:write", "journal:read", "journal:write", "account:read", "report:read"},
	}
	if openPeriod != nil {
		ctxOut.CurrentPeriodID = openPeriod.ID
		ctxOut.CurrentPeriodLabel = openPeriod.Label
	}
	return ctxOut, nil
}

func (s *companyService) CreateBranch(ctx context.Context, b *domain.CompanyBranch) error {
	if b.BranchName == "" {
		return fmt.Errorf("branch name is required")
	}
	if b.BranchTaxCode == "" {
		return fmt.Errorf("branch tax code is required")
	}
	return s.repo.CreateBranch(ctx, b)
}

func (s *companyService) GetBranch(ctx context.Context, id string) (*domain.CompanyBranch, error) {
	return s.repo.GetBranchByID(ctx, id)
}

func (s *companyService) ListBranches(ctx context.Context, companyID string) ([]domain.CompanyBranch, error) {
	return s.repo.GetBranchesByCompany(ctx, companyID)
}

func (s *companyService) UpdateBranch(ctx context.Context, b *domain.CompanyBranch) error {
	return s.repo.UpdateBranch(ctx, b)
}

func (s *companyService) DeactivateBranch(ctx context.Context, id string) error {
	return s.repo.DeactivateBranch(ctx, id)
}

func (s *companyService) CreateFiscalYear(ctx context.Context, fy *domain.FiscalYear) error {
	if fy.Year < 2000 || fy.Year > 2155 {
		return fmt.Errorf("fiscal year out of range")
	}
	if fy.StartMonth < 1 || fy.StartMonth > 12 {
		return fmt.Errorf("start month must be 1-12")
	}
	return s.repo.CreateFiscalYear(ctx, fy)
}

func (s *companyService) GetFiscalYear(ctx context.Context, id string) (*domain.FiscalYear, error) {
	return s.repo.GetFiscalYearByID(ctx, id)
}

func (s *companyService) ListFiscalYears(ctx context.Context, companyID string) ([]domain.FiscalYear, error) {
	return s.repo.GetFiscalYearsByCompany(ctx, companyID)
}

func (s *companyService) GeneratePeriods(ctx context.Context, fy *domain.FiscalYear) ([]domain.PeriodV2, error) {
	if fy.StartMonth < 1 || fy.StartMonth > 12 {
		return nil, fmt.Errorf("invalid start month")
	}
	var periods []domain.PeriodV2
	for m := 0; m < 12; m++ {
		month := ((fy.StartMonth - 1 + m) % 12) + 1
		year := fy.Year
		if fy.StartMonth+m > 12 {
			year++
		}
		startDate := fmt.Sprintf("%04d-%02d-01", year, month)
		endDay := 30
		switch month {
		case 1, 3, 5, 7, 8, 10, 12:
			endDay = 31
		case 2:
			if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
				endDay = 29
			} else {
				endDay = 28
			}
		}
		endDate := fmt.Sprintf("%04d-%02d-%02d", year, month, endDay)
		label := fmt.Sprintf("Thang %02d/%04d", month, year)
		status := domain.PeriodV2Future
		if m == 0 {
			status = domain.PeriodV2Open
		}
		p := domain.PeriodV2{
			CompanyID:    fy.CompanyID,
			FiscalYearID: fy.ID,
			PeriodType:   domain.PeriodTypeMonthly,
			PeriodNumber: m + 1,
			Label:        label,
			StartDate:    startDate,
			EndDate:      endDate,
			Status:       status,
		}
		if err := s.repo.CreatePeriod(ctx, &p); err != nil {
			return nil, fmt.Errorf("create period %d: %w", m+1, err)
		}
		periods = append(periods, p)
	}
	return periods, nil
}

func (s *companyService) ClosePeriod(ctx context.Context, companyID, periodID, userID string) error {
	p, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		return err
	}
	if p.CompanyID != companyID {
		return domain.ErrAccessDenied
	}
	if p.Status == domain.PeriodV2PermanentlyClosed {
		return domain.ErrPeriodPermanentlyClosed
	}
	if p.Status == domain.PeriodV2Closed {
		return domain.ErrPeriodAlreadyClosed
	}
	return s.repo.UpdatePeriodStatus(ctx, periodID, domain.PeriodV2Closed, userID)
}

func (s *companyService) ReopenPeriod(ctx context.Context, companyID, periodID, userID string) error {
	p, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		return err
	}
	if p.CompanyID != companyID {
		return domain.ErrAccessDenied
	}
	if p.Status == domain.PeriodV2PermanentlyClosed {
		return domain.ErrPeriodPermanentlyClosed
	}
	if p.Status == domain.PeriodV2Open {
		return nil
	}
	if err := s.repo.UpdatePeriodStatus(ctx, periodID, domain.PeriodV2Open, userID); err != nil {
		return err
	}
	return s.repo.IncrementReopenCount(ctx, periodID)
}

func (s *companyService) PermanentClosePeriod(ctx context.Context, companyID, periodID, userID string) error {
	p, err := s.repo.GetPeriodByID(ctx, periodID)
	if err != nil {
		return err
	}
	if p.CompanyID != companyID {
		return domain.ErrAccessDenied
	}
	if p.Status == domain.PeriodV2PermanentlyClosed {
		return domain.ErrPeriodPermanentlyClosed
	}
	return s.repo.UpdatePeriodStatus(ctx, periodID, domain.PeriodV2PermanentlyClosed, userID)
}

func (s *companyService) GetCurrentPeriod(ctx context.Context, companyID string) (*domain.PeriodV2, error) {
	return s.repo.GetOpenPeriod(ctx, companyID)
}

func (s *companyService) CreateDepartment(ctx context.Context, d *domain.Department) error {
	if d.Code == "" {
		return fmt.Errorf("department code is required")
	}
	if d.Name == "" {
		return fmt.Errorf("department name is required")
	}
	return s.repo.CreateDepartment(ctx, d)
}

func (s *companyService) GetDepartment(ctx context.Context, id string) (*domain.Department, error) {
	return s.repo.GetDepartmentByID(ctx, id)
}

func (s *companyService) ListDepartments(ctx context.Context, companyID string) ([]domain.Department, error) {
	return s.repo.GetDepartmentsByCompany(ctx, companyID)
}

func (s *companyService) UpdateDepartment(ctx context.Context, d *domain.Department) error {
	return s.repo.UpdateDepartment(ctx, d)
}

func (s *companyService) CreateEmployee(ctx context.Context, e *domain.Employee) error {
	if e.EmployeeCode == "" {
		return fmt.Errorf("employee code is required")
	}
	if e.FullName == "" {
		return fmt.Errorf("employee name is required")
	}
	if e.CompanyID == "" {
		return fmt.Errorf("company id is required")
	}
	return s.repo.CreateEmployee(ctx, e)
}

func (s *companyService) GetEmployee(ctx context.Context, id string) (*domain.Employee, error) {
	return s.repo.GetEmployeeByID(ctx, id)
}

func (s *companyService) ListEmployees(ctx context.Context, companyID string) ([]domain.Employee, error) {
	return s.repo.GetEmployeesByCompany(ctx, companyID)
}

func (s *companyService) UpdateEmployee(ctx context.Context, e *domain.Employee) error {
	return s.repo.UpdateEmployee(ctx, e)
}

func (s *companyService) DeactivateEmployee(ctx context.Context, id string) error {
	return s.repo.DeactivateEmployee(ctx, id)
}

func (s *companyService) CreateBankAccount(ctx context.Context, ba *domain.CompanyBankAccount) error {
	if ba.BankName == "" {
		return fmt.Errorf("bank name is required")
	}
	if ba.AccountNumber == "" {
		return fmt.Errorf("account number is required")
	}
	if ba.AccountHolder == "" {
		return fmt.Errorf("account holder is required")
	}
	return s.repo.CreateBankAccount(ctx, ba)
}

func (s *companyService) GetBankAccount(ctx context.Context, id string) (*domain.CompanyBankAccount, error) {
	return s.repo.GetBankAccountByID(ctx, id)
}

func (s *companyService) ListBankAccounts(ctx context.Context, companyID string) ([]domain.CompanyBankAccount, error) {
	return s.repo.GetBankAccountsByCompany(ctx, companyID)
}

func (s *companyService) UpdateBankAccount(ctx context.Context, ba *domain.CompanyBankAccount) error {
	return s.repo.UpdateBankAccount(ctx, ba)
}

func (s *companyService) DeactivateBankAccount(ctx context.Context, id string) error {
	return s.repo.DeactivateBankAccount(ctx, id)
}

func (s *companyService) RegisterEInvoicePattern(ctx context.Context, inv *domain.EInvoicePattern) error {
	if inv.PatternCode == "" {
		return fmt.Errorf("pattern code is required")
	}
	if inv.Serial == "" {
		return fmt.Errorf("serial is required")
	}
	return s.repo.CreateEInvoicePattern(ctx, inv)
}

func (s *companyService) GetEInvoicePattern(ctx context.Context, id string) (*domain.EInvoicePattern, error) {
	return s.repo.GetEInvoicePatternByID(ctx, id)
}

func (s *companyService) ListEInvoicePatterns(ctx context.Context, companyID string) ([]domain.EInvoicePattern, error) {
	return s.repo.GetEInvoicePatternsByCompany(ctx, companyID)
}

func (s *companyService) RegisterDigitalSignature(ctx context.Context, sig *domain.DigitalSignature) error {
	if sig.SerialNumber == "" {
		return fmt.Errorf("serial number is required")
	}
	if sig.ValidFrom == "" || sig.ValidTo == "" {
		return fmt.Errorf("valid from and valid to dates are required")
	}
	return s.repo.CreateDigitalSignature(ctx, sig)
}

func (s *companyService) GetDigitalSignature(ctx context.Context, id string) (*domain.DigitalSignature, error) {
	return s.repo.GetDigitalSignatureByID(ctx, id)
}

func (s *companyService) ListDigitalSignatures(ctx context.Context, companyID string) ([]domain.DigitalSignature, error) {
	return s.repo.GetDigitalSignaturesByCompany(ctx, companyID)
}

func (s *companyService) CreateIntegrationProfile(ctx context.Context, p *domain.IntegrationProfile) error {
	switch p.IntegrationType {
	case domain.IntegrationGDT, domain.IntegrationCustoms, domain.IntegrationBHXH, domain.IntegrationDVC:
	default:
		return fmt.Errorf("invalid integration type: %s", p.IntegrationType)
	}
	return s.repo.CreateIntegrationProfile(ctx, p)
}

func (s *companyService) GetIntegrationProfile(ctx context.Context, id string) (*domain.IntegrationProfile, error) {
	return s.repo.GetIntegrationProfileByID(ctx, id)
}

func (s *companyService) ListIntegrationProfiles(ctx context.Context, companyID string) ([]domain.IntegrationProfile, error) {
	return s.repo.GetIntegrationProfilesByCompany(ctx, companyID)
}

func (s *companyService) UpdateIntegrationProfile(ctx context.Context, p *domain.IntegrationProfile) error {
	return s.repo.UpdateIntegrationProfile(ctx, p)
}

func (s *companyService) TestIntegration(_ context.Context, _ string) error {
	return nil
}
