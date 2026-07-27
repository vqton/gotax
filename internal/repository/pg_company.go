package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gotax/internal/domain"
)

type PGCompanyRepo struct {
	pool *pgxpool.Pool
}

func NewPGCompanyRepo(pool *pgxpool.Pool) *PGCompanyRepo {
	return &PGCompanyRepo{pool: pool}
}

// ─── Company ────────────────────────────────────────────────────────

func (r *PGCompanyRepo) Create(ctx context.Context, c *domain.Company) error {
	if c.ID == "" {
		c.ID = "CMP" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO companies (id,tenant_id,legal_name_vn,legal_name_en,short_name,legal_form,tax_code,
		 business_reg_no,business_reg_date,business_reg_place,reg_address,reg_province,reg_district,
		 head_office_address,head_office_province,head_office_district,phone,email,website,
		 legal_rep_name,legal_rep_title,legal_rep_id_number,chief_accountant,chief_accountant_email,
		 tax_office_code,tax_office_name,accounting_regime,fiscal_year_start_month,default_currency,
		 secondary_currency,company_type,company_size,status,parent_company_id,logo_url,settings)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36)`,
		c.ID, c.TenantID, c.LegalNameVN, nullStr(c.LegalNameEN), nullStr(c.ShortName), c.LegalForm, c.TaxCode,
		nullStr(c.BusinessRegNo), nullStr(c.BusinessRegDate), nullStr(c.BusinessRegPlace),
		c.RegAddress, nullStr(c.RegProvince), nullStr(c.RegDistrict),
		nullStr(c.HeadOfficeAddress), nullStr(c.HeadOfficeProvince), nullStr(c.HeadOfficeDistrict),
		nullStr(c.Phone), nullStr(c.Email), nullStr(c.Website),
		nullStr(c.LegalRepName), nullStr(c.LegalRepTitle), nullStr(c.LegalRepIDNumber),
		nullStr(c.ChiefAccountant), nullStr(c.ChiefAccountantEmail),
		nullStr(c.TaxOfficeCode), nullStr(c.TaxOfficeName),
		c.AccountingRegime, c.FiscalYearStartMonth, c.DefaultCurrency,
		nullStr(c.SecondaryCurrency), nullStr(string(c.CompanyType)), nullStr(string(c.CompanySize)),
		c.Status, nullStr(c.ParentCompanyID), nullStr(c.LogoURL), c.Settings)
	return err
}

func (r *PGCompanyRepo) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,tenant_id,legal_name_vn,COALESCE(legal_name_en,''),COALESCE(short_name,''),
		 legal_form,tax_code,COALESCE(business_reg_no,''),COALESCE(business_reg_date,''),
		 COALESCE(business_reg_place,''),reg_address,COALESCE(reg_province,''),COALESCE(reg_district,''),
		 COALESCE(head_office_address,''),COALESCE(head_office_province,''),COALESCE(head_office_district,''),
		 COALESCE(phone,''),COALESCE(email,''),COALESCE(website,''),
		 COALESCE(legal_rep_name,''),COALESCE(legal_rep_title,''),COALESCE(legal_rep_id_number,''),
		 COALESCE(chief_accountant,''),COALESCE(chief_accountant_email,''),
		 COALESCE(tax_office_code,''),COALESCE(tax_office_name,''),
		 accounting_regime,fiscal_year_start_month,default_currency,
		 COALESCE(secondary_currency,''),COALESCE(company_type,''),COALESCE(company_size,''),
		 status,COALESCE(parent_company_id,''),COALESCE(logo_url,''),COALESCE(settings,'{}'),
		 created_at,updated_at
		 FROM companies WHERE id=$1`, id)
	return scanCompany(row)
}

func (r *PGCompanyRepo) GetByTaxCode(ctx context.Context, tenantID, taxCode string) (*domain.Company, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,tenant_id,legal_name_vn,COALESCE(legal_name_en,''),COALESCE(short_name,''),
		 legal_form,tax_code,COALESCE(business_reg_no,''),COALESCE(business_reg_date,''),
		 COALESCE(business_reg_place,''),reg_address,COALESCE(reg_province,''),COALESCE(reg_district,''),
		 COALESCE(head_office_address,''),COALESCE(head_office_province,''),COALESCE(head_office_district,''),
		 COALESCE(phone,''),COALESCE(email,''),COALESCE(website,''),
		 COALESCE(legal_rep_name,''),COALESCE(legal_rep_title,''),COALESCE(legal_rep_id_number,''),
		 COALESCE(chief_accountant,''),COALESCE(chief_accountant_email,''),
		 COALESCE(tax_office_code,''),COALESCE(tax_office_name,''),
		 accounting_regime,fiscal_year_start_month,default_currency,
		 COALESCE(secondary_currency,''),COALESCE(company_type,''),COALESCE(company_size,''),
		 status,COALESCE(parent_company_id,''),COALESCE(logo_url,''),COALESCE(settings,'{}'),
		 created_at,updated_at
		 FROM companies WHERE tenant_id=$1 AND tax_code=$2`, tenantID, taxCode)
	return scanCompany(row)
}

func (r *PGCompanyRepo) GetAll(ctx context.Context, tenantID string) ([]domain.Company, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,tenant_id,legal_name_vn,COALESCE(legal_name_en,''),COALESCE(short_name,''),
		 legal_form,tax_code,COALESCE(business_reg_no,''),COALESCE(business_reg_date,''),
		 COALESCE(business_reg_place,''),reg_address,COALESCE(reg_province,''),COALESCE(reg_district,''),
		 COALESCE(head_office_address,''),COALESCE(head_office_province,''),COALESCE(head_office_district,''),
		 COALESCE(phone,''),COALESCE(email,''),COALESCE(website,''),
		 COALESCE(legal_rep_name,''),COALESCE(legal_rep_title,''),COALESCE(legal_rep_id_number,''),
		 COALESCE(chief_accountant,''),COALESCE(chief_accountant_email,''),
		 COALESCE(tax_office_code,''),COALESCE(tax_office_name,''),
		 accounting_regime,fiscal_year_start_month,default_currency,
		 COALESCE(secondary_currency,''),COALESCE(company_type,''),COALESCE(company_size,''),
		 status,COALESCE(parent_company_id,''),COALESCE(logo_url,''),COALESCE(settings,'{}'),
		 created_at,updated_at
		 FROM companies WHERE tenant_id=$1 ORDER BY legal_name_vn`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCompanies(rows)
}

func (r *PGCompanyRepo) Update(ctx context.Context, c *domain.Company) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE companies SET legal_name_vn=$1,legal_name_en=$2,short_name=$3,legal_form=$4,tax_code=$5,
		 business_reg_no=$6,business_reg_date=$7,business_reg_place=$8,reg_address=$9,
		 reg_province=$10,reg_district=$11,head_office_address=$12,head_office_province=$13,
		 head_office_district=$14,phone=$15,email=$16,website=$17,legal_rep_name=$18,
		 legal_rep_title=$19,legal_rep_id_number=$20,chief_accountant=$21,chief_accountant_email=$22,
		 tax_office_code=$23,tax_office_name=$24,accounting_regime=$25,fiscal_year_start_month=$26,
		 default_currency=$27,secondary_currency=$28,company_type=$29,company_size=$30,
		 parent_company_id=$31,logo_url=$32,settings=$33,updated_at=NOW()
		 WHERE id=$34`,
		c.LegalNameVN, nullStr(c.LegalNameEN), nullStr(c.ShortName), c.LegalForm, c.TaxCode,
		nullStr(c.BusinessRegNo), nullStr(c.BusinessRegDate), nullStr(c.BusinessRegPlace),
		c.RegAddress, nullStr(c.RegProvince), nullStr(c.RegDistrict),
		nullStr(c.HeadOfficeAddress), nullStr(c.HeadOfficeProvince), nullStr(c.HeadOfficeDistrict),
		nullStr(c.Phone), nullStr(c.Email), nullStr(c.Website),
		nullStr(c.LegalRepName), nullStr(c.LegalRepTitle), nullStr(c.LegalRepIDNumber),
		nullStr(c.ChiefAccountant), nullStr(c.ChiefAccountantEmail),
		nullStr(c.TaxOfficeCode), nullStr(c.TaxOfficeName),
		c.AccountingRegime, c.FiscalYearStartMonth, c.DefaultCurrency,
		nullStr(c.SecondaryCurrency), nullStr(string(c.CompanyType)), nullStr(string(c.CompanySize)),
		nullStr(c.ParentCompanyID), nullStr(c.LogoURL), c.Settings, c.ID)
	return err
}

func (r *PGCompanyRepo) Deactivate(ctx context.Context, id, reason string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE companies SET status='DISSOLVED',updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *PGCompanyRepo) GetHierarchy(ctx context.Context, companyID string) ([]domain.Company, error) {
	rows, err := r.pool.Query(ctx,
		`WITH RECURSIVE tree AS (
		 SELECT id,tenant_id,legal_name_vn,COALESCE(legal_name_en,''),COALESCE(short_name,''),
		  legal_form,tax_code,COALESCE(business_reg_no,''),COALESCE(business_reg_date,''),
		  COALESCE(business_reg_place,''),reg_address,COALESCE(reg_province,''),COALESCE(reg_district,''),
		  COALESCE(head_office_address,''),COALESCE(head_office_province,''),COALESCE(head_office_district,''),
		  COALESCE(phone,''),COALESCE(email,''),COALESCE(website,''),
		  COALESCE(legal_rep_name,''),COALESCE(legal_rep_title,''),COALESCE(legal_rep_id_number,''),
		  COALESCE(chief_accountant,''),COALESCE(chief_accountant_email,''),
		  COALESCE(tax_office_code,''),COALESCE(tax_office_name,''),
		  accounting_regime,fiscal_year_start_month,default_currency,
		  COALESCE(secondary_currency,''),COALESCE(company_type,''),COALESCE(company_size,''),
		  status,COALESCE(parent_company_id,''),COALESCE(logo_url,''),COALESCE(settings,'{}'),
		  created_at,updated_at
		 FROM companies WHERE id=$1
		 UNION
		 SELECT c.id,c.tenant_id,c.legal_name_vn,COALESCE(c.legal_name_en,''),COALESCE(c.short_name,''),
		  c.legal_form,c.tax_code,COALESCE(c.business_reg_no,''),COALESCE(c.business_reg_date,''),
		  COALESCE(c.business_reg_place,''),c.reg_address,COALESCE(c.reg_province,''),COALESCE(c.reg_district,''),
		  COALESCE(c.head_office_address,''),COALESCE(c.head_office_province,''),COALESCE(c.head_office_district,''),
		  COALESCE(c.phone,''),COALESCE(c.email,''),COALESCE(c.website,''),
		  COALESCE(c.legal_rep_name,''),COALESCE(c.legal_rep_title,''),COALESCE(c.legal_rep_id_number,''),
		  COALESCE(c.chief_accountant,''),COALESCE(c.chief_accountant_email,''),
		  COALESCE(c.tax_office_code,''),COALESCE(c.tax_office_name,''),
		  c.accounting_regime,c.fiscal_year_start_month,c.default_currency,
		  COALESCE(c.secondary_currency,''),COALESCE(c.company_type,''),COALESCE(c.company_size,''),
		  c.status,COALESCE(c.parent_company_id,''),COALESCE(c.logo_url,''),COALESCE(c.settings,'{}'),
		  c.created_at,c.updated_at
		 FROM companies c INNER JOIN tree t ON c.parent_company_id=t.id)
		SELECT id,tenant_id,legal_name_vn,COALESCE(legal_name_en,''),COALESCE(short_name,''),
		 legal_form,tax_code,COALESCE(business_reg_no,''),COALESCE(business_reg_date,''),
		 COALESCE(business_reg_place,''),reg_address,COALESCE(reg_province,''),COALESCE(reg_district,''),
		 COALESCE(head_office_address,''),COALESCE(head_office_province,''),COALESCE(head_office_district,''),
		 COALESCE(phone,''),COALESCE(email,''),COALESCE(website,''),
		 COALESCE(legal_rep_name,''),COALESCE(legal_rep_title,''),COALESCE(legal_rep_id_number,''),
		 COALESCE(chief_accountant,''),COALESCE(chief_accountant_email,''),
		 COALESCE(tax_office_code,''),COALESCE(tax_office_name,''),
		 accounting_regime,fiscal_year_start_month,default_currency,
		 COALESCE(secondary_currency,''),COALESCE(company_type,''),COALESCE(company_size,''),
		 status,COALESCE(parent_company_id,''),COALESCE(logo_url,''),COALESCE(settings,'{}'),
		 created_at,updated_at
		 FROM tree`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCompanies(rows)
}

// ─── Branch ─────────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateBranch(ctx context.Context, b *domain.CompanyBranch) error {
	if b.ID == "" {
		b.ID = "BR" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO company_branches (id,company_id,branch_name,branch_tax_code,branch_type,address,phone,email,manager_name,status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		b.ID, b.CompanyID, b.BranchName, b.BranchTaxCode, b.BranchType,
		nullStr(b.Address), nullStr(b.Phone), nullStr(b.Email), nullStr(b.ManagerName), b.Status)
	return err
}

func (r *PGCompanyRepo) GetBranchByID(ctx context.Context, id string) (*domain.CompanyBranch, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,branch_name,branch_tax_code,branch_type,
		 COALESCE(address,''),COALESCE(phone,''),COALESCE(email,''),COALESCE(manager_name,''),
		 status,created_at,updated_at
		 FROM company_branches WHERE id=$1`, id)
	return scanBranch(row)
}

func (r *PGCompanyRepo) GetBranchesByCompany(ctx context.Context, companyID string) ([]domain.CompanyBranch, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,branch_name,branch_tax_code,branch_type,
		 COALESCE(address,''),COALESCE(phone,''),COALESCE(email,''),COALESCE(manager_name,''),
		 status,created_at,updated_at
		 FROM company_branches WHERE company_id=$1 ORDER BY branch_name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBranches(rows)
}

func (r *PGCompanyRepo) UpdateBranch(ctx context.Context, b *domain.CompanyBranch) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE company_branches SET branch_name=$1,branch_tax_code=$2,branch_type=$3,
		 address=$4,phone=$5,email=$6,manager_name=$7,updated_at=NOW() WHERE id=$8`,
		b.BranchName, b.BranchTaxCode, b.BranchType,
		nullStr(b.Address), nullStr(b.Phone), nullStr(b.Email), nullStr(b.ManagerName), b.ID)
	return err
}

func (r *PGCompanyRepo) DeactivateBranch(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE company_branches SET status='INACTIVE',updated_at=NOW() WHERE id=$1`, id)
	return err
}

// ─── Fiscal Year ────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateFiscalYear(ctx context.Context, fy *domain.FiscalYear) error {
	if fy.ID == "" {
		fy.ID = "FY" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO fiscal_years (id,company_id,year,start_month,is_short_year,start_date,end_date,status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		fy.ID, fy.CompanyID, fy.Year, fy.StartMonth, fy.IsShortYear,
		nullStr(fy.StartDate), nullStr(fy.EndDate), fy.Status)
	return err
}

func (r *PGCompanyRepo) GetFiscalYearByID(ctx context.Context, id string) (*domain.FiscalYear, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,year,start_month,is_short_year,COALESCE(start_date,''),COALESCE(end_date,''),
		 status,created_at FROM fiscal_years WHERE id=$1`, id)
	return scanFiscalYear(row)
}

func (r *PGCompanyRepo) GetFiscalYearsByCompany(ctx context.Context, companyID string) ([]domain.FiscalYear, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,year,start_month,is_short_year,COALESCE(start_date,''),COALESCE(end_date,''),
		 status,created_at FROM fiscal_years WHERE company_id=$1 ORDER BY year DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFiscalYears(rows)
}

func (r *PGCompanyRepo) GetFiscalYearByYear(ctx context.Context, companyID string, year int) (*domain.FiscalYear, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,year,start_month,is_short_year,COALESCE(start_date,''),COALESCE(end_date,''),
		 status,created_at FROM fiscal_years WHERE company_id=$1 AND year=$2`, companyID, year)
	return scanFiscalYear(row)
}

// ─── Period V2 ──────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreatePeriod(ctx context.Context, p *domain.PeriodV2) error {
	if p.ID == "" {
		p.ID = "P2" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO periods_v2 (id,company_id,fiscal_year_id,period_type,period_number,label,
		 start_date,end_date,status,opened_at,closed_at,closed_by,reopened_count)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		p.ID, p.CompanyID, p.FiscalYearID, p.PeriodType, p.PeriodNumber, p.Label,
		p.StartDate, p.EndDate, p.Status,
		nullStr(p.OpenedAt), nullStr(p.ClosedAt), nullStr(p.ClosedBy), p.ReopenedCount)
	return err
}

func (r *PGCompanyRepo) GetPeriodByID(ctx context.Context, id string) (*domain.PeriodV2, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,fiscal_year_id,period_type,period_number,label,start_date,end_date,
		 status,COALESCE(opened_at,''),COALESCE(closed_at,''),COALESCE(closed_by,''),
		 reopened_count FROM periods_v2 WHERE id=$1`, id)
	return scanPeriodV2(row)
}

func (r *PGCompanyRepo) GetPeriodsByFiscalYear(ctx context.Context, fiscalYearID string) ([]domain.PeriodV2, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,fiscal_year_id,period_type,period_number,label,start_date,end_date,
		 status,COALESCE(opened_at,''),COALESCE(closed_at,''),COALESCE(closed_by,''),
		 reopened_count FROM periods_v2 WHERE fiscal_year_id=$1 ORDER BY period_number`, fiscalYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPeriodsV2(rows)
}

func (r *PGCompanyRepo) GetPeriodsByCompany(ctx context.Context, companyID string) ([]domain.PeriodV2, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.id,p.company_id,p.fiscal_year_id,p.period_type,p.period_number,p.label,p.start_date,p.end_date,
		 p.status,COALESCE(p.opened_at,''),COALESCE(p.closed_at,''),COALESCE(p.closed_by,''),
		 p.reopened_count FROM periods_v2 p
		 JOIN fiscal_years fy ON fy.id=p.fiscal_year_id
		 WHERE p.company_id=$1 ORDER BY fy.year DESC, p.period_number`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPeriodsV2(rows)
}

func (r *PGCompanyRepo) GetOpenPeriod(ctx context.Context, companyID string) (*domain.PeriodV2, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,fiscal_year_id,period_type,period_number,label,start_date,end_date,
		 status,COALESCE(opened_at,''),COALESCE(closed_at,''),COALESCE(closed_by,''),
		 reopened_count FROM periods_v2 WHERE company_id=$1 AND status='OPEN' LIMIT 1`, companyID)
	return scanPeriodV2(row)
}

func (r *PGCompanyRepo) UpdatePeriodStatus(ctx context.Context, id string, status domain.PeriodStatusV2, closedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE periods_v2 SET status=$1,closed_at=$2,closed_by=$3 WHERE id=$4`,
		status, time.Now().Format(time.RFC3339), closedBy, id)
	return err
}

func (r *PGCompanyRepo) IncrementReopenCount(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE periods_v2 SET reopened_count=reopened_count+1 WHERE id=$1`, id)
	return err
}

// ─── Department ─────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateDepartment(ctx context.Context, d *domain.Department) error {
	if d.ID == "" {
		d.ID = "DEPT" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO departments (id,company_id,code,name,parent_id,manager_id,status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		d.ID, d.CompanyID, d.Code, d.Name, nullStr(d.ParentID), nullStr(d.ManagerID), d.Status)
	return err
}

func (r *PGCompanyRepo) GetDepartmentByID(ctx context.Context, id string) (*domain.Department, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,code,name,COALESCE(parent_id,''),COALESCE(manager_id,''),
		 status,created_at,updated_at FROM departments WHERE id=$1`, id)
	return scanDepartment(row)
}

func (r *PGCompanyRepo) GetDepartmentsByCompany(ctx context.Context, companyID string) ([]domain.Department, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,code,name,COALESCE(parent_id,''),COALESCE(manager_id,''),
		 status,created_at,updated_at FROM departments WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDepartments(rows)
}

func (r *PGCompanyRepo) UpdateDepartment(ctx context.Context, d *domain.Department) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE departments SET code=$1,name=$2,parent_id=$3,manager_id=$4,updated_at=NOW() WHERE id=$5`,
		d.Code, d.Name, nullStr(d.ParentID), nullStr(d.ManagerID), d.ID)
	return err
}

// ─── Employee ───────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateEmployee(ctx context.Context, e *domain.Employee) error {
	if e.ID == "" {
		e.ID = "EMP" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO employees (id,company_id,employee_code,full_name,title,email,phone,department_id,
		 personal_tax_code,social_insurance_no,bank_account_no,user_id,status,hire_date,termination_date)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		e.ID, e.CompanyID, e.EmployeeCode, e.FullName,
		nullStr(e.Title), nullStr(e.Email), nullStr(e.Phone), nullStr(e.DepartmentID),
		nullStr(e.PersonalTaxCode), nullStr(e.SocialInsuranceNo), nullStr(e.BankAccountNo),
		nullStr(e.UserID), e.Status, nullStr(e.HireDate), nullStr(e.TerminationDate))
	return err
}

func (r *PGCompanyRepo) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,employee_code,full_name,COALESCE(title,''),COALESCE(email,''),COALESCE(phone,''),
		 COALESCE(department_id,''),COALESCE(personal_tax_code,''),COALESCE(social_insurance_no,''),
		 COALESCE(bank_account_no,''),COALESCE(user_id,''),status,COALESCE(hire_date,''),
		 COALESCE(termination_date,''),created_at,updated_at FROM employees WHERE id=$1`, id)
	return scanEmployee(row)
}

func (r *PGCompanyRepo) GetEmployeesByCompany(ctx context.Context, companyID string) ([]domain.Employee, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,employee_code,full_name,COALESCE(title,''),COALESCE(email,''),COALESCE(phone,''),
		 COALESCE(department_id,''),COALESCE(personal_tax_code,''),COALESCE(social_insurance_no,''),
		 COALESCE(bank_account_no,''),COALESCE(user_id,''),status,COALESCE(hire_date,''),
		 COALESCE(termination_date,''),created_at,updated_at FROM employees WHERE company_id=$1 ORDER BY full_name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEmployees(rows)
}

func (r *PGCompanyRepo) GetEmployeeByCode(ctx context.Context, companyID, code string) (*domain.Employee, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,employee_code,full_name,COALESCE(title,''),COALESCE(email,''),COALESCE(phone,''),
		 COALESCE(department_id,''),COALESCE(personal_tax_code,''),COALESCE(social_insurance_no,''),
		 COALESCE(bank_account_no,''),COALESCE(user_id,''),status,COALESCE(hire_date,''),
		 COALESCE(termination_date,''),created_at,updated_at
		 FROM employees WHERE company_id=$1 AND employee_code=$2`, companyID, code)
	return scanEmployee(row)
}

func (r *PGCompanyRepo) UpdateEmployee(ctx context.Context, e *domain.Employee) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE employees SET employee_code=$1,full_name=$2,title=$3,email=$4,phone=$5,department_id=$6,
		 personal_tax_code=$7,social_insurance_no=$8,bank_account_no=$9,user_id=$10,
		 hire_date=$11,termination_date=$12,updated_at=NOW() WHERE id=$13`,
		e.EmployeeCode, e.FullName, nullStr(e.Title), nullStr(e.Email), nullStr(e.Phone),
		nullStr(e.DepartmentID), nullStr(e.PersonalTaxCode), nullStr(e.SocialInsuranceNo),
		nullStr(e.BankAccountNo), nullStr(e.UserID), nullStr(e.HireDate), nullStr(e.TerminationDate), e.ID)
	return err
}

func (r *PGCompanyRepo) DeactivateEmployee(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE employees SET status='TERMINATED',updated_at=NOW() WHERE id=$1`, id)
	return err
}

// ─── Bank Account ───────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateBankAccount(ctx context.Context, ba *domain.CompanyBankAccount) error {
	if ba.ID == "" {
		ba.ID = "BA" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO company_bank_accounts (id,company_id,bank_code,bank_name,branch_name,account_number,
		 account_holder,currency,is_default,is_verified,status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		ba.ID, ba.CompanyID, nullStr(ba.BankCode), ba.BankName, nullStr(ba.BranchName),
		ba.AccountNumber, ba.AccountHolder, ba.Currency, ba.IsDefault, ba.IsVerified, ba.Status)
	return err
}

func (r *PGCompanyRepo) GetBankAccountByID(ctx context.Context, id string) (*domain.CompanyBankAccount, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,COALESCE(bank_code,''),bank_name,COALESCE(branch_name,''),
		 account_number,account_holder,currency,is_default,is_verified,status,created_at,updated_at
		 FROM company_bank_accounts WHERE id=$1`, id)
	return scanBankAccount(row)
}

func (r *PGCompanyRepo) GetBankAccountsByCompany(ctx context.Context, companyID string) ([]domain.CompanyBankAccount, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,COALESCE(bank_code,''),bank_name,COALESCE(branch_name,''),
		 account_number,account_holder,currency,is_default,is_verified,status,created_at,updated_at
		 FROM company_bank_accounts WHERE company_id=$1 ORDER BY is_default DESC, bank_name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBankAccounts(rows)
}

func (r *PGCompanyRepo) UpdateBankAccount(ctx context.Context, ba *domain.CompanyBankAccount) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE company_bank_accounts SET bank_code=$1,bank_name=$2,branch_name=$3,account_number=$4,
		 account_holder=$5,currency=$6,is_default=$7,is_verified=$8,updated_at=NOW() WHERE id=$9`,
		nullStr(ba.BankCode), ba.BankName, nullStr(ba.BranchName), ba.AccountNumber,
		ba.AccountHolder, ba.Currency, ba.IsDefault, ba.IsVerified, ba.ID)
	return err
}

func (r *PGCompanyRepo) DeactivateBankAccount(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE company_bank_accounts SET status='CLOSED',updated_at=NOW() WHERE id=$1`, id)
	return err
}

// ─── E-Invoice ──────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateEInvoicePattern(ctx context.Context, inv *domain.EInvoicePattern) error {
	if inv.ID == "" {
		inv.ID = "INV" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO e_invoice_patterns (id,company_id,pattern_code,serial,form,invoice_type,status,gdt_status,description)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		inv.ID, inv.CompanyID, inv.PatternCode, inv.Serial,
		nullStr(inv.Form), nullStr(inv.InvoiceType), inv.Status,
		nullStr(inv.GDTStatus), nullStr(inv.Description))
	return err
}

func (r *PGCompanyRepo) GetEInvoicePatternByID(ctx context.Context, id string) (*domain.EInvoicePattern, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,pattern_code,serial,COALESCE(form,''),COALESCE(invoice_type,''),
		 status,COALESCE(gdt_status,''),COALESCE(description,''),created_at,updated_at
		 FROM e_invoice_patterns WHERE id=$1`, id)
	return scanEInvoicePattern(row)
}

func (r *PGCompanyRepo) GetEInvoicePatternsByCompany(ctx context.Context, companyID string) ([]domain.EInvoicePattern, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,pattern_code,serial,COALESCE(form,''),COALESCE(invoice_type,''),
		 status,COALESCE(gdt_status,''),COALESCE(description,''),created_at,updated_at
		 FROM e_invoice_patterns WHERE company_id=$1 ORDER BY pattern_code`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEInvoicePatterns(rows)
}

func (r *PGCompanyRepo) UpdateEInvoicePattern(ctx context.Context, inv *domain.EInvoicePattern) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE e_invoice_patterns SET pattern_code=$1,serial=$2,form=$3,invoice_type=$4,
		 status=$5,gdt_status=$6,description=$7,updated_at=NOW() WHERE id=$8`,
		inv.PatternCode, inv.Serial, nullStr(inv.Form), nullStr(inv.InvoiceType),
		inv.Status, nullStr(inv.GDTStatus), nullStr(inv.Description), inv.ID)
	return err
}

// ─── Digital Signature ──────────────────────────────────────────────

func (r *PGCompanyRepo) CreateDigitalSignature(ctx context.Context, sig *domain.DigitalSignature) error {
	if sig.ID == "" {
		sig.ID = "SIG" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO digital_signatures (id,company_id,signature_type,provider,serial_number,owner_name,
		 certificate_subject,certificate_issuer,valid_from,valid_to,status,is_default)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		sig.ID, sig.CompanyID, sig.SignatureType, nullStr(sig.Provider), sig.SerialNumber,
		nullStr(sig.OwnerName), nullStr(sig.CertificateSubject), nullStr(sig.CertificateIssuer),
		sig.ValidFrom, sig.ValidTo, sig.Status, sig.IsDefault)
	return err
}

func (r *PGCompanyRepo) GetDigitalSignatureByID(ctx context.Context, id string) (*domain.DigitalSignature, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,signature_type,COALESCE(provider,''),serial_number,COALESCE(owner_name,''),
		 COALESCE(certificate_subject,''),COALESCE(certificate_issuer,''),valid_from,valid_to,
		 status,is_default,created_at,updated_at
		 FROM digital_signatures WHERE id=$1`, id)
	return scanDigitalSignature(row)
}

func (r *PGCompanyRepo) GetDigitalSignaturesByCompany(ctx context.Context, companyID string) ([]domain.DigitalSignature, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,signature_type,COALESCE(provider,''),serial_number,COALESCE(owner_name,''),
		 COALESCE(certificate_subject,''),COALESCE(certificate_issuer,''),valid_from,valid_to,
		 status,is_default,created_at,updated_at
		 FROM digital_signatures WHERE company_id=$1 ORDER BY is_default DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDigitalSignatures(rows)
}

func (r *PGCompanyRepo) UpdateDigitalSignature(ctx context.Context, sig *domain.DigitalSignature) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE digital_signatures SET signature_type=$1,provider=$2,serial_number=$3,owner_name=$4,
		 certificate_subject=$5,certificate_issuer=$6,valid_from=$7,valid_to=$8,status=$9,
		 is_default=$10,updated_at=NOW() WHERE id=$11`,
		sig.SignatureType, nullStr(sig.Provider), sig.SerialNumber, nullStr(sig.OwnerName),
		nullStr(sig.CertificateSubject), nullStr(sig.CertificateIssuer),
		sig.ValidFrom, sig.ValidTo, sig.Status, sig.IsDefault, sig.ID)
	return err
}

// ─── Integration ────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateIntegrationProfile(ctx context.Context, prof *domain.IntegrationProfile) error {
	if prof.ID == "" {
		prof.ID = "INT" + time.Now().Format("20060102150405")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO integration_profiles (id,company_id,integration_type,endpoint_url,status,
		 last_connected_at,last_error_at,last_error_msg,config)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		prof.ID, prof.CompanyID, prof.IntegrationType, nullStr(prof.EndpointURL), prof.Status,
		nullStr(prof.LastConnectedAt), nullStr(prof.LastErrorAt), nullStr(prof.LastErrorMsg), prof.Config)
	return err
}

func (r *PGCompanyRepo) GetIntegrationProfileByID(ctx context.Context, id string) (*domain.IntegrationProfile, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,integration_type,COALESCE(endpoint_url,''),status,
		 COALESCE(last_connected_at,''),COALESCE(last_error_at,''),COALESCE(last_error_msg,''),
		 COALESCE(config,'{}'),created_at,updated_at
		 FROM integration_profiles WHERE id=$1`, id)
	return scanIntegrationProfile(row)
}

func (r *PGCompanyRepo) GetIntegrationProfilesByCompany(ctx context.Context, companyID string) ([]domain.IntegrationProfile, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,integration_type,COALESCE(endpoint_url,''),status,
		 COALESCE(last_connected_at,''),COALESCE(last_error_at,''),COALESCE(last_error_msg,''),
		 COALESCE(config,'{}'),created_at,updated_at
		 FROM integration_profiles WHERE company_id=$1 ORDER BY integration_type`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanIntegrationProfiles(rows)
}

func (r *PGCompanyRepo) GetIntegrationByType(ctx context.Context, companyID string, itype domain.IntegrationType) (*domain.IntegrationProfile, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,integration_type,COALESCE(endpoint_url,''),status,
		 COALESCE(last_connected_at,''),COALESCE(last_error_at,''),COALESCE(last_error_msg,''),
		 COALESCE(config,'{}'),created_at,updated_at
		 FROM integration_profiles WHERE company_id=$1 AND integration_type=$2`, companyID, itype)
	return scanIntegrationProfile(row)
}

func (r *PGCompanyRepo) UpdateIntegrationProfile(ctx context.Context, prof *domain.IntegrationProfile) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE integration_profiles SET integration_type=$1,endpoint_url=$2,status=$3,
		 last_connected_at=$4,last_error_at=$5,last_error_msg=$6,config=$7,updated_at=NOW() WHERE id=$8`,
		prof.IntegrationType, nullStr(prof.EndpointURL), prof.Status,
		nullStr(prof.LastConnectedAt), nullStr(prof.LastErrorAt), nullStr(prof.LastErrorMsg), prof.Config, prof.ID)
	return err
}

// ─── Scanners ───────────────────────────────────────────────────────

func scanCompany(row scannable) (*domain.Company, error) {
	c := &domain.Company{}
	err := row.Scan(&c.ID, &c.TenantID, &c.LegalNameVN, &c.LegalNameEN, &c.ShortName,
		&c.LegalForm, &c.TaxCode, &c.BusinessRegNo, &c.BusinessRegDate, &c.BusinessRegPlace,
		&c.RegAddress, &c.RegProvince, &c.RegDistrict,
		&c.HeadOfficeAddress, &c.HeadOfficeProvince, &c.HeadOfficeDistrict,
		&c.Phone, &c.Email, &c.Website,
		&c.LegalRepName, &c.LegalRepTitle, &c.LegalRepIDNumber,
		&c.ChiefAccountant, &c.ChiefAccountantEmail,
		&c.TaxOfficeCode, &c.TaxOfficeName,
		&c.AccountingRegime, &c.FiscalYearStartMonth, &c.DefaultCurrency,
		&c.SecondaryCurrency, &c.CompanyType, &c.CompanySize,
		&c.Status, &c.ParentCompanyID, &c.LogoURL, &c.Settings,
		&c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func scanCompanies(rows pgx.Rows) ([]domain.Company, error) {
	defer rows.Close()
	var companies []domain.Company
	for rows.Next() {
		c, err := scanCompany(rows)
		if err != nil {
			return nil, err
		}
		companies = append(companies, *c)
	}
	return companies, nil
}

func scanBranch(row scannable) (*domain.CompanyBranch, error) {
	b := &domain.CompanyBranch{}
	err := row.Scan(&b.ID, &b.CompanyID, &b.BranchName, &b.BranchTaxCode, &b.BranchType,
		&b.Address, &b.Phone, &b.Email, &b.ManagerName, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func scanBranches(rows pgx.Rows) ([]domain.CompanyBranch, error) {
	defer rows.Close()
	var branches []domain.CompanyBranch
	for rows.Next() {
		b, err := scanBranch(rows)
		if err != nil {
			return nil, err
		}
		branches = append(branches, *b)
	}
	return branches, nil
}

func scanFiscalYear(row scannable) (*domain.FiscalYear, error) {
	f := &domain.FiscalYear{}
	err := row.Scan(&f.ID, &f.CompanyID, &f.Year, &f.StartMonth, &f.IsShortYear,
		&f.StartDate, &f.EndDate, &f.Status, &f.CreatedAt)
	return f, err
}

func scanFiscalYears(rows pgx.Rows) ([]domain.FiscalYear, error) {
	defer rows.Close()
	var years []domain.FiscalYear
	for rows.Next() {
		f, err := scanFiscalYear(rows)
		if err != nil {
			return nil, err
		}
		years = append(years, *f)
	}
	return years, nil
}

func scanPeriodV2(row scannable) (*domain.PeriodV2, error) {
	p := &domain.PeriodV2{}
	err := row.Scan(&p.ID, &p.CompanyID, &p.FiscalYearID, &p.PeriodType, &p.PeriodNumber, &p.Label,
		&p.StartDate, &p.EndDate, &p.Status, &p.OpenedAt, &p.ClosedAt, &p.ClosedBy, &p.ReopenedCount)
	return p, err
}

func scanPeriodsV2(rows pgx.Rows) ([]domain.PeriodV2, error) {
	defer rows.Close()
	var periods []domain.PeriodV2
	for rows.Next() {
		p, err := scanPeriodV2(rows)
		if err != nil {
			return nil, err
		}
		periods = append(periods, *p)
	}
	return periods, nil
}

func scanDepartment(row scannable) (*domain.Department, error) {
	d := &domain.Department{}
	err := row.Scan(&d.ID, &d.CompanyID, &d.Code, &d.Name, &d.ParentID, &d.ManagerID,
		&d.Status, &d.CreatedAt, &d.UpdatedAt)
	return d, err
}

func scanDepartments(rows pgx.Rows) ([]domain.Department, error) {
	defer rows.Close()
	var depts []domain.Department
	for rows.Next() {
		d, err := scanDepartment(rows)
		if err != nil {
			return nil, err
		}
		depts = append(depts, *d)
	}
	return depts, nil
}

func scanEmployee(row scannable) (*domain.Employee, error) {
	e := &domain.Employee{}
	err := row.Scan(&e.ID, &e.CompanyID, &e.EmployeeCode, &e.FullName,
		&e.Title, &e.Email, &e.Phone, &e.DepartmentID,
		&e.PersonalTaxCode, &e.SocialInsuranceNo, &e.BankAccountNo, &e.UserID,
		&e.Status, &e.HireDate, &e.TerminationDate, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func scanEmployees(rows pgx.Rows) ([]domain.Employee, error) {
	defer rows.Close()
	var emps []domain.Employee
	for rows.Next() {
		e, err := scanEmployee(rows)
		if err != nil {
			return nil, err
		}
		emps = append(emps, *e)
	}
	return emps, nil
}

func scanBankAccount(row scannable) (*domain.CompanyBankAccount, error) {
	ba := &domain.CompanyBankAccount{}
	err := row.Scan(&ba.ID, &ba.CompanyID, &ba.BankCode, &ba.BankName, &ba.BranchName,
		&ba.AccountNumber, &ba.AccountHolder, &ba.Currency, &ba.IsDefault, &ba.IsVerified,
		&ba.Status, &ba.CreatedAt, &ba.UpdatedAt)
	return ba, err
}

func scanBankAccounts(rows pgx.Rows) ([]domain.CompanyBankAccount, error) {
	defer rows.Close()
	var accounts []domain.CompanyBankAccount
	for rows.Next() {
		ba, err := scanBankAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, *ba)
	}
	return accounts, nil
}

func scanEInvoicePattern(row scannable) (*domain.EInvoicePattern, error) {
	inv := &domain.EInvoicePattern{}
	err := row.Scan(&inv.ID, &inv.CompanyID, &inv.PatternCode, &inv.Serial,
		&inv.Form, &inv.InvoiceType, &inv.Status, &inv.GDTStatus, &inv.Description,
		&inv.CreatedAt, &inv.UpdatedAt)
	return inv, err
}

func scanEInvoicePatterns(rows pgx.Rows) ([]domain.EInvoicePattern, error) {
	defer rows.Close()
	var patterns []domain.EInvoicePattern
	for rows.Next() {
		inv, err := scanEInvoicePattern(rows)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, *inv)
	}
	return patterns, nil
}

func scanDigitalSignature(row scannable) (*domain.DigitalSignature, error) {
	sig := &domain.DigitalSignature{}
	err := row.Scan(&sig.ID, &sig.CompanyID, &sig.SignatureType, &sig.Provider, &sig.SerialNumber,
		&sig.OwnerName, &sig.CertificateSubject, &sig.CertificateIssuer,
		&sig.ValidFrom, &sig.ValidTo, &sig.Status, &sig.IsDefault, &sig.CreatedAt, &sig.UpdatedAt)
	return sig, err
}

func scanDigitalSignatures(rows pgx.Rows) ([]domain.DigitalSignature, error) {
	defer rows.Close()
	var sigs []domain.DigitalSignature
	for rows.Next() {
		sig, err := scanDigitalSignature(rows)
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, *sig)
	}
	return sigs, nil
}

func scanIntegrationProfile(row scannable) (*domain.IntegrationProfile, error) {
	prof := &domain.IntegrationProfile{}
	err := row.Scan(&prof.ID, &prof.CompanyID, &prof.IntegrationType, &prof.EndpointURL,
		&prof.Status, &prof.LastConnectedAt, &prof.LastErrorAt, &prof.LastErrorMsg,
		&prof.Config, &prof.CreatedAt, &prof.UpdatedAt)
	return prof, err
}

func scanIntegrationProfiles(rows pgx.Rows) ([]domain.IntegrationProfile, error) {
	defer rows.Close()
	var profs []domain.IntegrationProfile
	for rows.Next() {
		prof, err := scanIntegrationProfile(rows)
		if err != nil {
			return nil, err
		}
		profs = append(profs, *prof)
	}
	return profs, nil
}
