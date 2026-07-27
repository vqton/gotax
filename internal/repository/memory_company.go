package repository

import (
	"context"
	"sync"
	"time"

	"gotax/internal/domain"
)

type MemoryCompanyRepo struct {
	mu        sync.RWMutex
	companies map[string]*domain.Company
	branches  map[string]*domain.CompanyBranch
	fys       map[string]*domain.FiscalYear
	periods   map[string]*domain.PeriodV2
	depts     map[string]*domain.Department
	employees map[string]*domain.Employee
	bankAccs  map[string]*domain.CompanyBankAccount
	eInvoices map[string]*domain.EInvoicePattern
	sigs      map[string]*domain.DigitalSignature
	integrs   map[string]*domain.IntegrationProfile
	counter   int
}

func NewMemoryCompanyRepo() *MemoryCompanyRepo {
	return &MemoryCompanyRepo{
		companies: make(map[string]*domain.Company),
		branches:  make(map[string]*domain.CompanyBranch),
		fys:       make(map[string]*domain.FiscalYear),
		periods:   make(map[string]*domain.PeriodV2),
		depts:     make(map[string]*domain.Department),
		employees: make(map[string]*domain.Employee),
		bankAccs:  make(map[string]*domain.CompanyBankAccount),
		eInvoices: make(map[string]*domain.EInvoicePattern),
		sigs:      make(map[string]*domain.DigitalSignature),
		integrs:   make(map[string]*domain.IntegrationProfile),
	}
}

func (r *MemoryCompanyRepo) nextID(prefix string) string {
	r.counter++
	return prefix + time.Now().Format("20060102150405") + formatInt(r.counter)
}

// ─── Company ────────────────────────────────────────────────────────

func (r *MemoryCompanyRepo) Create(_ context.Context, c *domain.Company) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.companies {
		if existing.TaxCode == c.TaxCode && existing.TenantID == c.TenantID {
			return domain.ErrCompanyTaxCodeExists
		}
	}
	cp := *c
	if cp.ID == "" {
		cp.ID = r.nextID("CMP")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	r.companies[cp.ID] = &cp
	c.ID = cp.ID
	return nil
}

func (r *MemoryCompanyRepo) GetByID(_ context.Context, id string) (*domain.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.companies[id]
	if !ok {
		return nil, domain.ErrCompanyNotFound
	}
	return c, nil
}

func (r *MemoryCompanyRepo) GetByTaxCode(_ context.Context, tenantID, taxCode string) (*domain.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.companies {
		if c.TaxCode == taxCode && c.TenantID == tenantID {
			return c, nil
		}
	}
	return nil, domain.ErrCompanyNotFound
}

func (r *MemoryCompanyRepo) GetAll(_ context.Context, tenantID string) ([]domain.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Company
	for _, c := range r.companies {
		if c.TenantID == tenantID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) Update(_ context.Context, c *domain.Company) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.companies[c.ID]
	if !ok {
		return domain.ErrCompanyNotFound
	}
	cp := *c
	cp.CreatedAt = existing.CreatedAt
	r.companies[c.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) Deactivate(_ context.Context, id, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.companies[id]
	if !ok {
		return domain.ErrCompanyNotFound
	}
	c.Status = domain.CompanyStatusActive
	c.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryCompanyRepo) GetHierarchy(_ context.Context, companyID string) ([]domain.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Company
	for _, c := range r.companies {
		if c.ParentCompanyID == companyID || c.ID == companyID {
			out = append(out, *c)
		}
	}
	return out, nil
}

// ─── Branch ─────────────────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreateBranch(_ context.Context, b *domain.CompanyBranch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *b
	if cp.ID == "" {
		cp.ID = r.nextID("BR")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	r.branches[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetBranchByID(_ context.Context, id string) (*domain.CompanyBranch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.branches[id]
	if !ok {
		return nil, domain.ErrBranchNotFound
	}
	return b, nil
}

func (r *MemoryCompanyRepo) GetBranchesByCompany(_ context.Context, companyID string) ([]domain.CompanyBranch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CompanyBranch
	for _, b := range r.branches {
		if b.CompanyID == companyID {
			out = append(out, *b)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) UpdateBranch(_ context.Context, b *domain.CompanyBranch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.branches[b.ID]
	if !ok {
		return domain.ErrBranchNotFound
	}
	cp := *b
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.branches[b.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) DeactivateBranch(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.branches[id]
	if !ok {
		return domain.ErrBranchNotFound
	}
	b.Status = "INACTIVE"
	b.UpdatedAt = time.Now()
	return nil
}

// ─── Fiscal Year ────────────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreateFiscalYear(_ context.Context, fy *domain.FiscalYear) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *fy
	if cp.ID == "" {
		cp.ID = r.nextID("FY")
	}
	cp.CreatedAt = time.Now()
	r.fys[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetFiscalYearByID(_ context.Context, id string) (*domain.FiscalYear, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fy, ok := r.fys[id]
	if !ok {
		return nil, domain.ErrFiscalYearNotFound
	}
	return fy, nil
}

func (r *MemoryCompanyRepo) GetFiscalYearsByCompany(_ context.Context, companyID string) ([]domain.FiscalYear, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.FiscalYear
	for _, fy := range r.fys {
		if fy.CompanyID == companyID {
			out = append(out, *fy)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) GetFiscalYearByYear(_ context.Context, companyID string, year int) (*domain.FiscalYear, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, fy := range r.fys {
		if fy.CompanyID == companyID && fy.Year == year {
			return fy, nil
		}
	}
	return nil, domain.ErrFiscalYearNotFound
}

// ─── Period V2 ──────────────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreatePeriod(_ context.Context, p *domain.PeriodV2) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	if cp.ID == "" {
		cp.ID = r.nextID("PER")
	}
	r.periods[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetPeriodByID(_ context.Context, id string) (*domain.PeriodV2, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.periods[id]
	if !ok {
		return nil, domain.ErrPeriodV2NotFound
	}
	return p, nil
}

func (r *MemoryCompanyRepo) GetPeriodsByFiscalYear(_ context.Context, fiscalYearID string) ([]domain.PeriodV2, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.PeriodV2
	for _, p := range r.periods {
		if p.FiscalYearID == fiscalYearID {
			out = append(out, *p)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) GetPeriodsByCompany(_ context.Context, companyID string) ([]domain.PeriodV2, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.PeriodV2
	for _, p := range r.periods {
		if p.CompanyID == companyID {
			out = append(out, *p)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) GetOpenPeriod(_ context.Context, companyID string) (*domain.PeriodV2, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.periods {
		if p.CompanyID == companyID && p.Status == domain.PeriodV2Open {
			return p, nil
		}
	}
	return nil, domain.ErrPeriodV2NotFound
}

func (r *MemoryCompanyRepo) UpdatePeriodStatus(_ context.Context, id string, status domain.PeriodStatusV2, closedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.periods[id]
	if !ok {
		return domain.ErrPeriodV2NotFound
	}
	p.Status = status
	if status == domain.PeriodV2Closed || status == domain.PeriodV2PermanentlyClosed {
		p.ClosedAt = time.Now().Format("2006-01-02 15:04:05")
		p.ClosedBy = closedBy
	}
	return nil
}

func (r *MemoryCompanyRepo) IncrementReopenCount(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.periods[id]
	if !ok {
		return domain.ErrPeriodV2NotFound
	}
	p.ReopenedCount++
	return nil
}

// ─── Department ─────────────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreateDepartment(_ context.Context, d *domain.Department) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *d
	if cp.ID == "" {
		cp.ID = r.nextID("DEPT")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	r.depts[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetDepartmentByID(_ context.Context, id string) (*domain.Department, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.depts[id]
	if !ok {
		return nil, domain.ErrDepartmentNotFound
	}
	return d, nil
}

func (r *MemoryCompanyRepo) GetDepartmentsByCompany(_ context.Context, companyID string) ([]domain.Department, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Department
	for _, d := range r.depts {
		if d.CompanyID == companyID {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) UpdateDepartment(_ context.Context, d *domain.Department) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.depts[d.ID]
	if !ok {
		return domain.ErrDepartmentNotFound
	}
	cp := *d
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.depts[d.ID] = &cp
	return nil
}

// ─── Employee ───────────────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreateEmployee(_ context.Context, e *domain.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *e
	if cp.ID == "" {
		cp.ID = r.nextID("EMP")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	r.employees[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetEmployeeByID(_ context.Context, id string) (*domain.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.employees[id]
	if !ok {
		return nil, domain.ErrEmployeeNotFound
	}
	return e, nil
}

func (r *MemoryCompanyRepo) GetEmployeesByCompany(_ context.Context, companyID string) ([]domain.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Employee
	for _, e := range r.employees {
		if e.CompanyID == companyID {
			out = append(out, *e)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) GetEmployeeByCode(_ context.Context, companyID, code string) (*domain.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, e := range r.employees {
		if e.CompanyID == companyID && e.EmployeeCode == code {
			return e, nil
		}
	}
	return nil, domain.ErrEmployeeNotFound
}

func (r *MemoryCompanyRepo) UpdateEmployee(_ context.Context, e *domain.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.employees[e.ID]
	if !ok {
		return domain.ErrEmployeeNotFound
	}
	cp := *e
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.employees[e.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) DeactivateEmployee(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.employees[id]
	if !ok {
		return domain.ErrEmployeeNotFound
	}
	e.Status = domain.EmployeeTerminated
	e.UpdatedAt = time.Now()
	return nil
}

// ─── Bank Account ───────────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreateBankAccount(_ context.Context, ba *domain.CompanyBankAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *ba
	if cp.ID == "" {
		cp.ID = r.nextID("BA")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	r.bankAccs[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetBankAccountByID(_ context.Context, id string) (*domain.CompanyBankAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ba, ok := r.bankAccs[id]
	if !ok {
		return nil, domain.ErrBankAccountNotFound
	}
	return ba, nil
}

func (r *MemoryCompanyRepo) GetBankAccountsByCompany(_ context.Context, companyID string) ([]domain.CompanyBankAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CompanyBankAccount
	for _, ba := range r.bankAccs {
		if ba.CompanyID == companyID {
			out = append(out, *ba)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) UpdateBankAccount(_ context.Context, ba *domain.CompanyBankAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.bankAccs[ba.ID]
	if !ok {
		return domain.ErrBankAccountNotFound
	}
	cp := *ba
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.bankAccs[ba.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) DeactivateBankAccount(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ba, ok := r.bankAccs[id]
	if !ok {
		return domain.ErrBankAccountNotFound
	}
	ba.Status = domain.BankAccountClosed
	ba.UpdatedAt = time.Now()
	return nil
}

// ─── E-Invoice Pattern ──────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreateEInvoicePattern(_ context.Context, inv *domain.EInvoicePattern) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *inv
	if cp.ID == "" {
		cp.ID = r.nextID("INV")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	r.eInvoices[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetEInvoicePatternByID(_ context.Context, id string) (*domain.EInvoicePattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inv, ok := r.eInvoices[id]
	if !ok {
		return nil, domain.ErrEInvoiceNotFound
	}
	return inv, nil
}

func (r *MemoryCompanyRepo) GetEInvoicePatternsByCompany(_ context.Context, companyID string) ([]domain.EInvoicePattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.EInvoicePattern
	for _, inv := range r.eInvoices {
		if inv.CompanyID == companyID {
			out = append(out, *inv)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) UpdateEInvoicePattern(_ context.Context, inv *domain.EInvoicePattern) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.eInvoices[inv.ID]
	if !ok {
		return domain.ErrEInvoiceNotFound
	}
	cp := *inv
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.eInvoices[inv.ID] = &cp
	return nil
}

// ─── Digital Signature ──────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreateDigitalSignature(_ context.Context, sig *domain.DigitalSignature) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *sig
	if cp.ID == "" {
		cp.ID = r.nextID("SIG")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	r.sigs[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetDigitalSignatureByID(_ context.Context, id string) (*domain.DigitalSignature, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sig, ok := r.sigs[id]
	if !ok {
		return nil, domain.ErrSignatureNotFound
	}
	return sig, nil
}

func (r *MemoryCompanyRepo) GetDigitalSignaturesByCompany(_ context.Context, companyID string) ([]domain.DigitalSignature, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.DigitalSignature
	for _, sig := range r.sigs {
		if sig.CompanyID == companyID {
			out = append(out, *sig)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) UpdateDigitalSignature(_ context.Context, sig *domain.DigitalSignature) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.sigs[sig.ID]
	if !ok {
		return domain.ErrSignatureNotFound
	}
	cp := *sig
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.sigs[sig.ID] = &cp
	return nil
}

// ─── Integration Profile ────────────────────────────────────────────

func (r *MemoryCompanyRepo) CreateIntegrationProfile(_ context.Context, prof *domain.IntegrationProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *prof
	if cp.ID == "" {
		cp.ID = r.nextID("INT")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	r.integrs[cp.ID] = &cp
	return nil
}

func (r *MemoryCompanyRepo) GetIntegrationProfileByID(_ context.Context, id string) (*domain.IntegrationProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	prof, ok := r.integrs[id]
	if !ok {
		return nil, domain.ErrIntegrationNotFound
	}
	return prof, nil
}

func (r *MemoryCompanyRepo) GetIntegrationProfilesByCompany(_ context.Context, companyID string) ([]domain.IntegrationProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.IntegrationProfile
	for _, prof := range r.integrs {
		if prof.CompanyID == companyID {
			out = append(out, *prof)
		}
	}
	return out, nil
}

func (r *MemoryCompanyRepo) GetIntegrationByType(_ context.Context, companyID string, itype domain.IntegrationType) (*domain.IntegrationProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, prof := range r.integrs {
		if prof.CompanyID == companyID && prof.IntegrationType == itype {
			return prof, nil
		}
	}
	return nil, domain.ErrIntegrationNotFound
}

func (r *MemoryCompanyRepo) UpdateIntegrationProfile(_ context.Context, prof *domain.IntegrationProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.integrs[prof.ID]
	if !ok {
		return domain.ErrIntegrationNotFound
	}
	cp := *prof
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.integrs[prof.ID] = &cp
	return nil
}
