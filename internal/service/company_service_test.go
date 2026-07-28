package service

import (
	"context"
	"gotax/internal/domain"
	"gotax/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCompanySvc(t *testing.T) (CompanyService, context.Context) {
	t.Helper()
	return NewCompanyService(repository.NewMemoryCompanyRepo()), context.Background()
}

func seedCompany(t *testing.T, svc CompanyService, ctx context.Context) *domain.Company {
	t.Helper()
	c := &domain.Company{
		TaxCode:          "1234567890",
		LegalNameVN:      "Test Co",
		LegalForm:        domain.LegalFormJSC,
		AccountingRegime: domain.RegimeTT99,
	}
	require.NoError(t, svc.CreateCompany(ctx, c))
	return c
}

// ─── Company ─────────────────────────────────────────────────────────────

func TestCreateCompany_Defaults(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := &domain.Company{
		TaxCode:          "1234567890",
		LegalNameVN:      "Test Co",
		LegalForm:        domain.LegalFormJSC,
		AccountingRegime: domain.RegimeTT99,
	}
	err := svc.CreateCompany(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, domain.CompanyStatusActive, c.Status)
	assert.Equal(t, "VND", c.DefaultCurrency)
	assert.Equal(t, 1, c.FiscalYearStartMonth)
	assert.NotEmpty(t, c.ID)
}

func TestCreateCompany_InvalidTaxCode(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := &domain.Company{
		TaxCode:          "12345",
		LegalNameVN:      "Test Co",
		LegalForm:        domain.LegalFormJSC,
		AccountingRegime: domain.RegimeTT99,
	}
	err := svc.CreateCompany(ctx, c)
	assert.ErrorIs(t, err, domain.ErrInvalidTaxCodeFormat)
}

func TestCreateCompany_DuplicateTaxCode(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	seedCompany(t, svc, ctx)
	c := &domain.Company{
		TaxCode:          "1234567890",
		LegalNameVN:      "Dup Co",
		LegalForm:        domain.LegalFormJSC,
		AccountingRegime: domain.RegimeTT99,
	}
	err := svc.CreateCompany(ctx, c)
	assert.ErrorIs(t, err, domain.ErrCompanyTaxCodeExists)
}

func TestCreateCompany_13DigitTaxCode(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := &domain.Company{
		TaxCode:          "1234567890-001",
		LegalNameVN:      "Sub Co",
		LegalForm:        domain.LegalFormJSC,
		AccountingRegime: domain.RegimeTT99,
	}
	err := svc.CreateCompany(ctx, c)
	require.NoError(t, err)
	assert.NotEmpty(t, c.ID)
}

func TestGetCompany(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	got, err := svc.GetCompany(ctx, c.ID)
	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
}

func TestGetCompany_NotFound(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	_, err := svc.GetCompany(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrCompanyNotFound)
}

func TestUpdateCompany(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	c.LegalNameVN = "Updated Co"
	err := svc.UpdateCompany(ctx, c)
	require.NoError(t, err)

	got, _ := svc.GetCompany(ctx, c.ID)
	assert.Equal(t, "Updated Co", got.LegalNameVN)
}

func TestUpdateCompany_Inactive(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	require.NoError(t, svc.DeactivateCompany(ctx, c.ID, "test"))

	c.LegalNameVN = "Updated Co"
	err := svc.UpdateCompany(ctx, c)
	assert.ErrorIs(t, err, domain.ErrCompanyInactive)
}

func TestDeactivateCompany(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	err := svc.DeactivateCompany(ctx, c.ID, "dissolved")
	require.NoError(t, err)

	got, _ := svc.GetCompany(ctx, c.ID)
	assert.Equal(t, domain.CompanyStatusDissolved, got.Status)
}

func TestDeactivateCompany_NotFound(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	err := svc.DeactivateCompany(ctx, "nonexistent", "test")
	assert.ErrorIs(t, err, domain.ErrCompanyNotFound)
}

// ─── Tax Code Lookup ──────────────────────────────────────────────────────

func TestGetCompanyByTaxCode(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	got, err := svc.GetCompanyByTaxCode(ctx, "", c.TaxCode)
	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
}

func TestGetCompanyByTaxCode_InvalidFormat(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	_, err := svc.GetCompanyByTaxCode(ctx, "", "bad")
	assert.ErrorIs(t, err, domain.ErrInvalidTaxCode)
}

// ─── Hierarchy ────────────────────────────────────────────────────────────

func TestGetCompanyHierarchy_Single(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	h, err := svc.GetCompanyHierarchy(ctx, c.ID)
	require.NoError(t, err)
	assert.Len(t, h, 1)
	assert.Equal(t, c.ID, h[0].ID)
}

func TestGetCompanyHierarchy_ParentChild(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	parent := seedCompany(t, svc, ctx)

	child := &domain.Company{
		TaxCode:           "1234567890-001",
		LegalNameVN:       "Child Co",
		LegalForm:         domain.LegalFormJSC,
		AccountingRegime:  domain.RegimeTT99,
		ParentCompanyID:   parent.ID,
	}
	require.NoError(t, svc.CreateCompany(ctx, child))

	h, err := svc.GetCompanyHierarchy(ctx, parent.ID)
	require.NoError(t, err)
	assert.Len(t, h, 2)
}

func TestGetCompanyHierarchy_Recursive(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	grandparent := seedCompany(t, svc, ctx)

	parent := &domain.Company{
		TaxCode:           "1234567890-001",
		LegalNameVN:       "Parent Co",
		LegalForm:         domain.LegalFormJSC,
		AccountingRegime:  domain.RegimeTT99,
		ParentCompanyID:   grandparent.ID,
	}
	require.NoError(t, svc.CreateCompany(ctx, parent))

	child := &domain.Company{
		TaxCode:           "1234567890-002",
		LegalNameVN:       "Child Co",
		LegalForm:         domain.LegalFormJSC,
		AccountingRegime:  domain.RegimeTT99,
		ParentCompanyID:   parent.ID,
	}
	require.NoError(t, svc.CreateCompany(ctx, child))

	h, err := svc.GetCompanyHierarchy(ctx, grandparent.ID)
	require.NoError(t, err)
	assert.Len(t, h, 3)
}

// ─── Switch ───────────────────────────────────────────────────────────────

func TestSwitchCompany(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	ctxOut, err := svc.SwitchCompany(ctx, "user-1", c.ID)
	require.NoError(t, err)
	assert.Equal(t, c.ID, ctxOut.CompanyID)
	assert.Equal(t, c.LegalNameVN, ctxOut.LegalNameVN)
	assert.Equal(t, c.TaxCode, ctxOut.TaxCode)
	assert.NotEmpty(t, ctxOut.Permissions)
}

func TestSwitchCompany_Inactive(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)
	require.NoError(t, svc.DeactivateCompany(ctx, c.ID, "test"))

	_, err := svc.SwitchCompany(ctx, "user-1", c.ID)
	assert.ErrorIs(t, err, domain.ErrCompanyInactive)
}

// ─── Branch ───────────────────────────────────────────────────────────────

func TestCreateBranch(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	b := &domain.CompanyBranch{
		CompanyID:    c.ID,
		BranchName:   "Hanoi Branch",
		BranchTaxCode: "1234567890-001",
	}
	err := svc.CreateBranch(ctx, b)
	require.NoError(t, err)
	assert.NotEmpty(t, b.ID)
}

func TestListBranches(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	b := &domain.CompanyBranch{CompanyID: c.ID, BranchName: "Hanoi", BranchTaxCode: "1234567890-001"}
	require.NoError(t, svc.CreateBranch(ctx, b))

	branches, err := svc.ListBranches(ctx, c.ID)
	require.NoError(t, err)
	assert.Len(t, branches, 1)
}

// ─── Fiscal Year ──────────────────────────────────────────────────────────

func TestCreateFiscalYear(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	fy := &domain.FiscalYear{CompanyID: c.ID, Year: 2026, StartMonth: 1}
	err := svc.CreateFiscalYear(ctx, fy)
	require.NoError(t, err)
	assert.NotEmpty(t, fy.ID)
}

// ─── Period Lifecycle ─────────────────────────────────────────────────────

func TestGeneratePeriods(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	fy := &domain.FiscalYear{CompanyID: c.ID, Year: 2026, StartMonth: 1}
	require.NoError(t, svc.CreateFiscalYear(ctx, fy))

	periods, err := svc.GeneratePeriods(ctx, fy)
	require.NoError(t, err)
	assert.Len(t, periods, 12)
	assert.Equal(t, domain.PeriodV2Open, periods[0].Status)
	for i := 1; i < 12; i++ {
		assert.Equal(t, domain.PeriodV2Future, periods[i].Status)
	}
}

func TestClosePeriod(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	fy := &domain.FiscalYear{CompanyID: c.ID, Year: 2026, StartMonth: 1}
	require.NoError(t, svc.CreateFiscalYear(ctx, fy))
	periods, _ := svc.GeneratePeriods(ctx, fy)

	err := svc.ClosePeriod(ctx, c.ID, periods[0].ID, "user-1")
	require.NoError(t, err)

	// No period open after close (only 1st was OPEN, rest are FUTURE)
	_, err = svc.GetCurrentPeriod(ctx, c.ID)
	assert.ErrorIs(t, err, domain.ErrPeriodV2NotFound)
}

func TestClosePeriod_WrongCompany(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c1 := seedCompany(t, svc, ctx)
	c2 := &domain.Company{TaxCode: "0987654321", LegalNameVN: "Other Co", LegalForm: domain.LegalFormJSC, AccountingRegime: domain.RegimeTT99}
	require.NoError(t, svc.CreateCompany(ctx, c2))

	fy := &domain.FiscalYear{CompanyID: c1.ID, Year: 2026, StartMonth: 1}
	require.NoError(t, svc.CreateFiscalYear(ctx, fy))
	periods, _ := svc.GeneratePeriods(ctx, fy)

	err := svc.ClosePeriod(ctx, c2.ID, periods[0].ID, "user-1")
	assert.ErrorIs(t, err, domain.ErrAccessDenied)
}

func TestReopenPeriod(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	fy := &domain.FiscalYear{CompanyID: c.ID, Year: 2026, StartMonth: 1}
	require.NoError(t, svc.CreateFiscalYear(ctx, fy))
	periods, _ := svc.GeneratePeriods(ctx, fy)
	require.NoError(t, svc.ClosePeriod(ctx, c.ID, periods[0].ID, "user-1"))

	err := svc.ReopenPeriod(ctx, c.ID, periods[0].ID, "user-1")
	require.NoError(t, err)

	reopened, _ := svc.GetCurrentPeriod(ctx, c.ID)
	assert.Equal(t, domain.PeriodV2Open, reopened.Status)
}

func TestPermanentClosePeriod(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	fy := &domain.FiscalYear{CompanyID: c.ID, Year: 2026, StartMonth: 1}
	require.NoError(t, svc.CreateFiscalYear(ctx, fy))
	periods, _ := svc.GeneratePeriods(ctx, fy)

	err := svc.PermanentClosePeriod(ctx, c.ID, periods[0].ID, "user-1")
	require.NoError(t, err)

	err = svc.ReopenPeriod(ctx, c.ID, periods[0].ID, "user-1")
	assert.ErrorIs(t, err, domain.ErrPeriodPermanentlyClosed)
}

// ─── Department ───────────────────────────────────────────────────────────

func TestCreateDepartment(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	d := &domain.Department{CompanyID: c.ID, Code: "ACC", Name: "Accounting"}
	err := svc.CreateDepartment(ctx, d)
	require.NoError(t, err)
	assert.NotEmpty(t, d.ID)
}

// ─── Employee ─────────────────────────────────────────────────────────────

func TestCreateEmployee(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	e := &domain.Employee{CompanyID: c.ID, EmployeeCode: "E001", FullName: "Nguyen Van A"}
	err := svc.CreateEmployee(ctx, e)
	require.NoError(t, err)
	assert.NotEmpty(t, e.ID)
}

func TestDeactivateEmployee(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	e := &domain.Employee{CompanyID: c.ID, EmployeeCode: "E001", FullName: "Nguyen Van A"}
	require.NoError(t, svc.CreateEmployee(ctx, e))
	require.NoError(t, svc.DeactivateEmployee(ctx, e.ID))

	got, _ := svc.GetEmployee(ctx, e.ID)
	assert.Equal(t, domain.EmployeeTerminated, got.Status)
}

// ─── Bank Account ─────────────────────────────────────────────────────────

func TestCreateBankAccount(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	ba := &domain.CompanyBankAccount{
		CompanyID:     c.ID,
		BankName:      "Vietcombank",
		AccountNumber: "1234567890",
		AccountHolder: "Test Co",
		Currency:      "VND",
	}
	err := svc.CreateBankAccount(ctx, ba)
	require.NoError(t, err)
	assert.NotEmpty(t, ba.ID)
}

// ─── E-Invoice Pattern ────────────────────────────────────────────────────

func TestRegisterEInvoicePattern(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	inv := &domain.EInvoicePattern{
		CompanyID:   c.ID,
		PatternCode: "01GTKT",
		Serial:      "AA/25E",
	}
	err := svc.RegisterEInvoicePattern(ctx, inv)
	require.NoError(t, err)
	assert.NotEmpty(t, inv.ID)
}

// ─── Digital Signature ────────────────────────────────────────────────────

func TestRegisterDigitalSignature(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	sig := &domain.DigitalSignature{
		CompanyID:    c.ID,
		SerialNumber: "SIG-001",
		ValidFrom:    "2026-01-01",
		ValidTo:      "2027-01-01",
	}
	err := svc.RegisterDigitalSignature(ctx, sig)
	require.NoError(t, err)
	assert.NotEmpty(t, sig.ID)
}

// ─── Integration Profile ──────────────────────────────────────────────────

func TestCreateIntegrationProfile(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	p := &domain.IntegrationProfile{
		CompanyID:       c.ID,
		IntegrationType: domain.IntegrationGDT,
		EndpointURL:     "https://thuedientu.gdt.gov.vn",
	}
	err := svc.CreateIntegrationProfile(ctx, p)
	require.NoError(t, err)
	assert.NotEmpty(t, p.ID)
}

func TestCreateIntegrationProfile_InvalidType(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	p := &domain.IntegrationProfile{
		CompanyID:       c.ID,
		IntegrationType: "BAD_TYPE",
	}
	err := svc.CreateIntegrationProfile(ctx, p)
	assert.Error(t, err)
}

func TestTestIntegration_NotFound(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	err := svc.TestIntegration(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrIntegrationNotFound)
}

func TestTestIntegration_Success(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	c := seedCompany(t, svc, ctx)

	p := &domain.IntegrationProfile{
		CompanyID:       c.ID,
		IntegrationType: domain.IntegrationGDT,
		EndpointURL:     "https://thuedientu.gdt.gov.vn",
	}
	require.NoError(t, svc.CreateIntegrationProfile(ctx, p))

	err := svc.TestIntegration(ctx, p.ID)
	require.NoError(t, err)
}

// ─── List ─────────────────────────────────────────────────────────────────

func TestListCompanies(t *testing.T) {
	svc, ctx := newCompanySvc(t)
	seedCompany(t, svc, ctx)

	cs, err := svc.ListCompanies(ctx, "")
	require.NoError(t, err)
	assert.Len(t, cs, 1)
}
