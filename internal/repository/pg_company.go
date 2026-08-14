package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type PGCompanyRepo struct {
	db *gorm.DB
}

func NewPGCompanyRepo(db *gorm.DB) *PGCompanyRepo {
	return &PGCompanyRepo{db: db}
}

// ─── Company ────────────────────────────────────────────────────────

func (r *PGCompanyRepo) Create(ctx context.Context, c *domain.Company) error {
	if c.ID == "" {
		c.ID = "CMP" + time.Now().Format("20060102150405")
	}
	status := string(c.Status)
	if status == "" {
		status = string(domain.CompanyStatusActive)
	}
	legalForm := string(c.LegalForm)
	if legalForm == "" {
		legalForm = "LLC_1MEMBER"
	}
	regime := string(c.AccountingRegime)
	if regime == "" {
		regime = "TT99"
	}
	m := &domain.CompanyGORM{
		ID:                   c.ID,
		TenantID:             c.TenantID,
		Name:                 c.LegalNameVN,
		TaxCode:              c.TaxCode,
		Address:              c.RegAddress,
		Phone:                strPtr(c.Phone),
		Email:                strPtr(c.Email),
		Website:              strPtr(c.Website),
		LegalRepName:         strPtr(c.LegalRepName),
		LegalRepTitle:        strPtr(c.LegalRepTitle),
		LegalRepID:           strPtr(c.LegalRepIDNumber),
		ChiefAccountant:      strPtr(c.ChiefAccountant),
		ChiefAccountantEmail: strPtr(c.ChiefAccountantEmail),
		TaxOfficeCode:        strPtr(c.TaxOfficeCode),
		TaxOffice:            strPtr(c.TaxOfficeName),
		LegalForm:            legalForm,
		AccountingRegime:     regime,
		FiscalYearStartMonth: c.FiscalYearStartMonth,
		DefaultCurrency:      c.DefaultCurrency,
		Status:               status,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if c.LegalNameEN != "" {
		m.NameEn = strPtr(c.LegalNameEN)
	}
	if c.ShortName != "" {
		m.ShortName = strPtr(c.ShortName)
	}
	if c.RegProvince != "" {
		m.City = strPtr(c.RegProvince)
	}
	if c.RegDistrict != "" {
		m.Province = strPtr(c.RegDistrict)
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	var m domain.CompanyGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormCompanyToDomain(&m), nil
}

func (r *PGCompanyRepo) GetByTaxCode(ctx context.Context, tenantID, taxCode string) (*domain.Company, error) {
	var m domain.CompanyGORM
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND tax_code = ?", tenantID, taxCode).First(&m).Error; err != nil {
		return nil, err
	}
	return gormCompanyToDomain(&m), nil
}

func (r *PGCompanyRepo) GetAll(ctx context.Context, tenantID string) ([]domain.Company, error) {
	var ms []domain.CompanyGORM
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("legal_name_vn ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Company, len(ms))
	for i := range ms {
		out[i] = *gormCompanyToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) Update(ctx context.Context, c *domain.Company) error {
	updates := map[string]any{
		"legal_name_vn":        c.LegalNameVN,
		"tax_code":             c.TaxCode,
		"reg_address":          c.RegAddress,
		"phone":                c.Phone,
		"email":                c.Email,
		"website":              c.Website,
		"legal_rep_name":       c.LegalRepName,
		"legal_rep_id_number":  c.LegalRepIDNumber,
		"chief_accountant":     c.ChiefAccountant,
		"tax_office_name":      c.TaxOfficeName,
		"legal_form":           string(c.LegalForm),
		"accounting_regime":    string(c.AccountingRegime),
		"status":               string(c.Status),
		"legal_rep_title":      c.LegalRepTitle,
		"short_name":           c.ShortName,
		"fiscal_year_start_month": c.FiscalYearStartMonth,
		"default_currency":     c.DefaultCurrency,
		"updated_at":           time.Now(),
	}
	if c.LegalNameEN != "" {
		updates["legal_name_en"] = c.LegalNameEN
	}
	if c.RegProvince != "" {
		updates["reg_province"] = c.RegProvince
	}
	if c.RegDistrict != "" {
		updates["reg_district"] = c.RegDistrict
	}
	return r.db.WithContext(ctx).Model(&domain.CompanyGORM{}).Where("id = ?", c.ID).Updates(updates).Error
}

func (r *PGCompanyRepo) Deactivate(ctx context.Context, id, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.CompanyGORM{}).Where("id = ?", id).Updates(map[string]any{
		"status":     string(domain.CompanyStatusDissolved),
		"updated_at": now,
	}).Error
}

func (r *PGCompanyRepo) GetHierarchy(ctx context.Context, companyID string) ([]domain.Company, error) {
	var ms []domain.CompanyGORM
	if err := r.db.WithContext(ctx).Raw(
		`WITH RECURSIVE tree AS (
		 SELECT * FROM companies WHERE id = ?
		 UNION
		 SELECT c.* FROM companies c INNER JOIN tree t ON c.parent_company_id = t.id)
		SELECT * FROM tree`, companyID).Scan(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Company, len(ms))
	for i := range ms {
		out[i] = *gormCompanyToDomain(&ms[i])
	}
	return out, nil
}

// ─── Branch ─────────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateBranch(ctx context.Context, b *domain.CompanyBranch) error {
	if b.ID == "" {
		b.ID = "BR" + time.Now().Format("20060102150405")
	}
	m := &domain.CompanyBranchGORM{
		ID:        b.ID,
		CompanyID: b.CompanyID,
		Code:      b.BranchTaxCode,
		Name:      b.BranchName,
		Address:   strPtr(b.Address),
		Phone:     strPtr(b.Phone),
		ManagerName: strPtr(b.ManagerName),
		IsActive:  true,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetBranchByID(ctx context.Context, id string) (*domain.CompanyBranch, error) {
	var m domain.CompanyBranchGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormBranchToDomain(&m), nil
}

func (r *PGCompanyRepo) GetBranchesByCompany(ctx context.Context, companyID string) ([]domain.CompanyBranch, error) {
	var ms []domain.CompanyBranchGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("name ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CompanyBranch, len(ms))
	for i := range ms {
		out[i] = *gormBranchToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) UpdateBranch(ctx context.Context, b *domain.CompanyBranch) error {
	return r.db.WithContext(ctx).Model(&domain.CompanyBranchGORM{}).Where("id = ?", b.ID).Updates(map[string]any{
		"code":       b.BranchTaxCode,
		"name":       b.BranchName,
		"address":    b.Address,
		"phone":      b.Phone,
		"manager_name": b.ManagerName,
		"updated_at": time.Now(),
	}).Error
}

func (r *PGCompanyRepo) DeactivateBranch(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.CompanyBranchGORM{}).Where("id = ?", id).Updates(map[string]any{
		"is_active":  false,
		"deactivated_at": &now,
		"updated_at": now,
	}).Error
}

// ─── Fiscal Year ────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateFiscalYear(ctx context.Context, fy *domain.FiscalYear) error {
	if fy.ID == "" {
		fy.ID = "FY" + time.Now().Format("20060102150405")
	}
	sd, _ := time.Parse("2006-01-02", fy.StartDate)
	ed, _ := time.Parse("2006-01-02", fy.EndDate)
	m := &domain.FiscalYearGORM{
		ID:        fy.ID,
		CompanyID: fy.CompanyID,
		Year:      fy.Year,
		StartDate: sd,
		EndDate:   ed,
		IsClosed:  fy.Status == "CLOSED",
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetFiscalYearByID(ctx context.Context, id string) (*domain.FiscalYear, error) {
	var m domain.FiscalYearGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormFiscalYearToDomain(&m), nil
}

func (r *PGCompanyRepo) GetFiscalYearsByCompany(ctx context.Context, companyID string) ([]domain.FiscalYear, error) {
	var ms []domain.FiscalYearGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("year DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.FiscalYear, len(ms))
	for i := range ms {
		out[i] = *gormFiscalYearToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) GetFiscalYearByYear(ctx context.Context, companyID string, year int) (*domain.FiscalYear, error) {
	var m domain.FiscalYearGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND year = ?", companyID, year).First(&m).Error; err != nil {
		return nil, err
	}
	return gormFiscalYearToDomain(&m), nil
}

// ─── Period V2 ──────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreatePeriod(ctx context.Context, p *domain.PeriodV2) error {
	if p.ID == "" {
		p.ID = "P2" + time.Now().Format("20060102150405")
	}
	sd, _ := time.Parse("2006-01-02", p.StartDate)
	ed, _ := time.Parse("2006-01-02", p.EndDate)
	m := &domain.PeriodV2GORM{
		ID:           p.ID,
		FiscalYearID: p.FiscalYearID,
		PeriodNumber: p.PeriodNumber,
		StartDate:    sd,
		EndDate:      ed,
		Status:       string(p.Status),
		ReopenCount:  p.ReopenedCount,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetPeriodByID(ctx context.Context, id string) (*domain.PeriodV2, error) {
	var m domain.PeriodV2GORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormPeriodToDomain(&m), nil
}

func (r *PGCompanyRepo) GetPeriodsByFiscalYear(ctx context.Context, fiscalYearID string) ([]domain.PeriodV2, error) {
	var ms []domain.PeriodV2GORM
	if err := r.db.WithContext(ctx).Where("fiscal_year_id = ?", fiscalYearID).Order("period_number ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.PeriodV2, len(ms))
	for i := range ms {
		out[i] = *gormPeriodToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) GetPeriodsByCompany(ctx context.Context, companyID string) ([]domain.PeriodV2, error) {
	var ms []domain.PeriodV2GORM
	if err := r.db.WithContext(ctx).
		Joins("JOIN fiscal_years fy ON fy.id = period_v2.fiscal_year_id").
		Where("fy.company_id = ?", companyID).
		Order("fy.year DESC, period_v2.period_number ASC").
		Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.PeriodV2, len(ms))
	for i := range ms {
		out[i] = *gormPeriodToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) GetOpenPeriod(ctx context.Context, companyID string) (*domain.PeriodV2, error) {
	var m domain.PeriodV2GORM
	if err := r.db.WithContext(ctx).
		Joins("JOIN fiscal_years fy ON fy.id = period_v2.fiscal_year_id").
		Where("fy.company_id = ? AND period_v2.status = 'OPEN'", companyID).
		First(&m).Error; err != nil {
		return nil, err
	}
	return gormPeriodToDomain(&m), nil
}

func (r *PGCompanyRepo) UpdatePeriodStatus(ctx context.Context, id string, status domain.PeriodStatusV2, closedBy string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.PeriodV2GORM{}).Where("id = ?", id).Updates(map[string]any{
		"status":    string(status),
		"closed_at": &now,
		"closed_by": closedBy,
	}).Error
}

func (r *PGCompanyRepo) IncrementReopenCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.PeriodV2GORM{}).Where("id = ?", id).
		UpdateColumn("reopen_count", gorm.Expr("reopen_count + 1")).Error
}

// ─── Department ─────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateDepartment(ctx context.Context, d *domain.Department) error {
	if d.ID == "" {
		d.ID = "DEPT" + time.Now().Format("20060102150405")
	}
	m := &domain.DepartmentGORM{
		ID:        d.ID,
		CompanyID: d.CompanyID,
		Code:      d.Code,
		Name:      d.Name,
		ManagerID: strPtr(d.ManagerID),
		IsActive:  true,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetDepartmentByID(ctx context.Context, id string) (*domain.Department, error) {
	var m domain.DepartmentGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormDepartmentToDomain(&m), nil
}

func (r *PGCompanyRepo) GetDepartmentsByCompany(ctx context.Context, companyID string) ([]domain.Department, error) {
	var ms []domain.DepartmentGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("name ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Department, len(ms))
	for i := range ms {
		out[i] = *gormDepartmentToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) UpdateDepartment(ctx context.Context, d *domain.Department) error {
	return r.db.WithContext(ctx).Model(&domain.DepartmentGORM{}).Where("id = ?", d.ID).Updates(map[string]any{
		"code":       d.Code,
		"name":       d.Name,
		"manager_id": d.ManagerID,
		"updated_at": time.Now(),
	}).Error
}

func (r *PGCompanyRepo) DeactivateDepartment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.DepartmentGORM{}).Where("id = ?", id).Updates(map[string]any{
		"is_active":  false,
		"updated_at": time.Now(),
	}).Error
}

// ─── Employee ───────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateEmployee(ctx context.Context, e *domain.Employee) error {
	if e.ID == "" {
		e.ID = "EMP" + time.Now().Format("20060102150405")
	}
	m := &domain.EmployeeGORM{
		ID:        e.ID,
		CompanyID: e.CompanyID,
		Code:      e.EmployeeCode,
		FullName:  e.FullName,
		Email:     strPtr(e.Email),
		Phone:     strPtr(e.Phone),
		DeptID:    strPtr(e.DepartmentID),
		Position:  strPtr(e.Title),
		IsActive:  e.Status != domain.EmployeeTerminated,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	var m domain.EmployeeGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormEmployeeToDomain(&m), nil
}

func (r *PGCompanyRepo) GetEmployeesByCompany(ctx context.Context, companyID string) ([]domain.Employee, error) {
	var ms []domain.EmployeeGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("full_name ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Employee, len(ms))
	for i := range ms {
		out[i] = *gormEmployeeToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) GetEmployeeByCode(ctx context.Context, companyID, code string) (*domain.Employee, error) {
	var m domain.EmployeeGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		return nil, err
	}
	return gormEmployeeToDomain(&m), nil
}

func (r *PGCompanyRepo) UpdateEmployee(ctx context.Context, e *domain.Employee) error {
	return r.db.WithContext(ctx).Model(&domain.EmployeeGORM{}).Where("id = ?", e.ID).Updates(map[string]any{
		"code":       e.EmployeeCode,
		"full_name":  e.FullName,
		"email":      e.Email,
		"phone":      e.Phone,
		"dept_id":    e.DepartmentID,
		"position":   e.Title,
		"updated_at": time.Now(),
	}).Error
}

func (r *PGCompanyRepo) DeactivateEmployee(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.EmployeeGORM{}).Where("id = ?", id).Updates(map[string]any{
		"is_active":  false,
		"updated_at": time.Now(),
	}).Error
}

// ─── Bank Account ───────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateBankAccount(ctx context.Context, ba *domain.CompanyBankAccount) error {
	if ba.ID == "" {
		ba.ID = "BA" + time.Now().Format("20060102150405")
	}
	m := &domain.CompanyBankAccountGORM{
		ID:            ba.ID,
		CompanyID:     ba.CompanyID,
		BankName:      ba.BankName,
		BankCode:      strPtr(ba.BankCode),
		Branch:        strPtr(ba.BranchName),
		AccountNumber: ba.AccountNumber,
		AccountName:   ba.AccountHolder,
		Currency:      ba.Currency,
		IsPrimary:     ba.IsDefault,
		IsActive:      ba.Status != domain.BankAccountClosed,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetBankAccountByID(ctx context.Context, id string) (*domain.CompanyBankAccount, error) {
	var m domain.CompanyBankAccountGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormBankAccountToDomain(&m), nil
}

func (r *PGCompanyRepo) GetBankAccountsByCompany(ctx context.Context, companyID string) ([]domain.CompanyBankAccount, error) {
	var ms []domain.CompanyBankAccountGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("is_primary DESC, bank_name ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CompanyBankAccount, len(ms))
	for i := range ms {
		out[i] = *gormBankAccountToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) UpdateBankAccount(ctx context.Context, ba *domain.CompanyBankAccount) error {
	return r.db.WithContext(ctx).Model(&domain.CompanyBankAccountGORM{}).Where("id = ?", ba.ID).Updates(map[string]any{
		"bank_name":      ba.BankName,
		"bank_code":      ba.BankCode,
		"branch":         ba.BranchName,
		"account_number": ba.AccountNumber,
		"account_name":   ba.AccountHolder,
		"currency":       ba.Currency,
		"is_primary":     ba.IsDefault,
		"updated_at":     time.Now(),
	}).Error
}

func (r *PGCompanyRepo) DeactivateBankAccount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.CompanyBankAccountGORM{}).Where("id = ?", id).Updates(map[string]any{
		"is_active":  false,
		"updated_at": time.Now(),
	}).Error
}

// ─── E-Invoice ──────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateEInvoicePattern(ctx context.Context, inv *domain.EInvoicePattern) error {
	if inv.ID == "" {
		inv.ID = "INV" + time.Now().Format("20060102150405")
	}
	m := &domain.EInvoicePatternGORM{
		ID:          inv.ID,
		CompanyID:   inv.CompanyID,
		Pattern:     inv.PatternCode,
		Serial:      inv.Serial,
		InvoiceType: inv.InvoiceType,
		IsActive:    inv.Status != domain.EInvoiceCancelled,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetEInvoicePatternByID(ctx context.Context, id string) (*domain.EInvoicePattern, error) {
	var m domain.EInvoicePatternGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormEInvoicePatternToDomain(&m), nil
}

func (r *PGCompanyRepo) GetEInvoicePatternsByCompany(ctx context.Context, companyID string) ([]domain.EInvoicePattern, error) {
	var ms []domain.EInvoicePatternGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("pattern ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.EInvoicePattern, len(ms))
	for i := range ms {
		out[i] = *gormEInvoicePatternToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) UpdateEInvoicePattern(ctx context.Context, inv *domain.EInvoicePattern) error {
	return r.db.WithContext(ctx).Model(&domain.EInvoicePatternGORM{}).Where("id = ?", inv.ID).Updates(map[string]any{
		"pattern":      inv.PatternCode,
		"serial":       inv.Serial,
		"invoice_type": inv.InvoiceType,
		"updated_at":   time.Now(),
	}).Error
}

// ─── Digital Signature ──────────────────────────────────────────────

func (r *PGCompanyRepo) CreateDigitalSignature(ctx context.Context, sig *domain.DigitalSignature) error {
	if sig.ID == "" {
		sig.ID = "SIG" + time.Now().Format("20060102150405")
	}
	vf, _ := time.Parse("2006-01-02", sig.ValidFrom)
	vt, _ := time.Parse("2006-01-02", sig.ValidTo)
	m := &domain.DigitalSignatureGORM{
		ID:           sig.ID,
		CompanyID:    sig.CompanyID,
		Name:         sig.SerialNumber,
		SerialNumber: sig.SerialNumber,
		Issuer:       strPtr(sig.CertificateIssuer),
		ValidFrom:    vf,
		ValidTo:      vt,
		IsActive:     sig.Status == domain.SignatureActive,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetDigitalSignatureByID(ctx context.Context, id string) (*domain.DigitalSignature, error) {
	var m domain.DigitalSignatureGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormDigitalSignatureToDomain(&m), nil
}

func (r *PGCompanyRepo) GetDigitalSignaturesByCompany(ctx context.Context, companyID string) ([]domain.DigitalSignature, error) {
	var ms []domain.DigitalSignatureGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("is_active DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.DigitalSignature, len(ms))
	for i := range ms {
		out[i] = *gormDigitalSignatureToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) UpdateDigitalSignature(ctx context.Context, sig *domain.DigitalSignature) error {
	vf, _ := time.Parse("2006-01-02", sig.ValidFrom)
	vt, _ := time.Parse("2006-01-02", sig.ValidTo)
	return r.db.WithContext(ctx).Model(&domain.DigitalSignatureGORM{}).Where("id = ?", sig.ID).Updates(map[string]any{
		"serial_number": sig.SerialNumber,
		"issuer":        sig.CertificateIssuer,
		"valid_from":    vf,
		"valid_to":      vt,
		"updated_at":    time.Now(),
	}).Error
}

// ─── Integration ────────────────────────────────────────────────────

func (r *PGCompanyRepo) CreateIntegrationProfile(ctx context.Context, prof *domain.IntegrationProfile) error {
	if prof.ID == "" {
		prof.ID = "INT" + time.Now().Format("20060102150405")
	}
	m := &domain.IntegrationProfileGORM{
		ID:        prof.ID,
		CompanyID: prof.CompanyID,
		Type:      string(prof.IntegrationType),
		Name:      string(prof.IntegrationType),
		IsActive:  prof.Status == domain.IntegrationConnected,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCompanyRepo) GetIntegrationProfileByID(ctx context.Context, id string) (*domain.IntegrationProfile, error) {
	var m domain.IntegrationProfileGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormIntegrationToDomain(&m), nil
}

func (r *PGCompanyRepo) GetIntegrationProfilesByCompany(ctx context.Context, companyID string) ([]domain.IntegrationProfile, error) {
	var ms []domain.IntegrationProfileGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("type ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.IntegrationProfile, len(ms))
	for i := range ms {
		out[i] = *gormIntegrationToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGCompanyRepo) GetIntegrationByType(ctx context.Context, companyID string, itype domain.IntegrationType) (*domain.IntegrationProfile, error) {
	var m domain.IntegrationProfileGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND type = ?", companyID, string(itype)).First(&m).Error; err != nil {
		return nil, err
	}
	return gormIntegrationToDomain(&m), nil
}

func (r *PGCompanyRepo) UpdateIntegrationProfile(ctx context.Context, prof *domain.IntegrationProfile) error {
	return r.db.WithContext(ctx).Model(&domain.IntegrationProfileGORM{}).Where("id = ?", prof.ID).Updates(map[string]any{
		"type":       string(prof.IntegrationType),
		"is_active":  prof.Status == domain.IntegrationConnected,
		"updated_at": time.Now(),
	}).Error
}

// ─── Converters ─────────────────────────────────────────────────────

func gormCompanyToDomain(m *domain.CompanyGORM) *domain.Company {
	d := &domain.Company{
		ID:        m.ID,
		TenantID:  m.TenantID,
		LegalNameVN: m.Name,
		TaxCode:   m.TaxCode,
		RegAddress: m.Address,
		Status:    domain.CompanyStatusActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.NameEn != nil {
		d.LegalNameEN = *m.NameEn
	}
	if m.City != nil {
		d.RegProvince = *m.City
	}
	if m.Phone != nil {
		d.Phone = *m.Phone
	}
	if m.Email != nil {
		d.Email = *m.Email
	}
	if m.Website != nil {
		d.Website = *m.Website
	}
	if m.LegalRepName != nil {
		d.LegalRepName = *m.LegalRepName
	}
	if m.LegalRepID != nil {
		d.LegalRepIDNumber = *m.LegalRepID
	}
	if m.ChiefAccountant != nil {
		d.ChiefAccountant = *m.ChiefAccountant
	}
	if m.TaxOffice != nil {
		d.TaxOfficeName = *m.TaxOffice
	}
	switch domain.CompanyStatus(m.Status) {
	case domain.CompanyStatusSuspended:
		d.Status = domain.CompanyStatusSuspended
	case domain.CompanyStatusDissolved:
		d.Status = domain.CompanyStatusDissolved
	case domain.CompanyStatusMerged:
		d.Status = domain.CompanyStatusMerged
	default:
		d.Status = domain.CompanyStatusActive
	}
	d.LegalForm = domain.LegalForm(m.LegalForm)
	d.AccountingRegime = domain.AccountingRegime(m.AccountingRegime)
	d.FiscalYearStartMonth = m.FiscalYearStartMonth
	d.DefaultCurrency = m.DefaultCurrency
	if m.ShortName != nil {
		d.ShortName = *m.ShortName
	}
	if m.Province != nil {
		d.RegDistrict = *m.Province
	}
	if m.LegalRepTitle != nil {
		d.LegalRepTitle = *m.LegalRepTitle
	}
	if m.TaxOfficeCode != nil {
		d.TaxOfficeCode = *m.TaxOfficeCode
	}
	if m.ChiefAccountantEmail != nil {
		d.ChiefAccountantEmail = *m.ChiefAccountantEmail
	}
	return d
}

func gormBranchToDomain(m *domain.CompanyBranchGORM) *domain.CompanyBranch {
	d := &domain.CompanyBranch{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		BranchName: m.Name,
		BranchTaxCode: m.Code,
		Status:    "ACTIVE",
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Address:   "",
		Phone:     "",
	}
	if !m.IsActive {
		d.Status = "INACTIVE"
	}
	if m.Address != nil {
		d.Address = *m.Address
	}
	if m.Phone != nil {
		d.Phone = *m.Phone
	}
	if m.ManagerName != nil {
		d.ManagerName = *m.ManagerName
	}
	return d
}

func gormFiscalYearToDomain(m *domain.FiscalYearGORM) *domain.FiscalYear {
	d := &domain.FiscalYear{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Year:      m.Year,
		StartDate: m.StartDate.Format("2006-01-02"),
		EndDate:   m.EndDate.Format("2006-01-02"),
		Status:    "OPEN",
		CreatedAt: m.CreatedAt,
	}
	if m.IsClosed {
		d.Status = "CLOSED"
	}
	return d
}

func gormPeriodToDomain(m *domain.PeriodV2GORM) *domain.PeriodV2 {
	d := &domain.PeriodV2{
		ID:            m.ID,
		FiscalYearID:  m.FiscalYearID,
		PeriodNumber:  m.PeriodNumber,
		Label:         fmt.Sprintf("Period %d", m.PeriodNumber),
		StartDate:     m.StartDate.Format("2006-01-02"),
		EndDate:       m.EndDate.Format("2006-01-02"),
		Status:        domain.PeriodStatusV2(m.Status),
		ReopenedCount: m.ReopenCount,
	}
	if m.ClosedAt != nil {
		d.ClosedAt = m.ClosedAt.Format(time.RFC3339)
	}
	if m.ClosedBy != nil {
		d.ClosedBy = *m.ClosedBy
	}
	return d
}

func gormDepartmentToDomain(m *domain.DepartmentGORM) *domain.Department {
	d := &domain.Department{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Code:      m.Code,
		Name:      m.Name,
		Status:    "ACTIVE",
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if !m.IsActive {
		d.Status = "INACTIVE"
	}
	if m.ManagerID != nil {
		d.ManagerID = *m.ManagerID
	}
	return d
}

func gormEmployeeToDomain(m *domain.EmployeeGORM) *domain.Employee {
	d := &domain.Employee{
		ID:           m.ID,
		CompanyID:    m.CompanyID,
		EmployeeCode: m.Code,
		FullName:     m.FullName,
		Status:       domain.EmployeeActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
	if !m.IsActive {
		d.Status = domain.EmployeeTerminated
	}
	if m.Email != nil {
		d.Email = *m.Email
	}
	if m.Phone != nil {
		d.Phone = *m.Phone
	}
	if m.DeptID != nil {
		d.DepartmentID = *m.DeptID
	}
	if m.Position != nil {
		d.Title = *m.Position
	}
	if m.IDNumber != nil {
		d.PersonalTaxCode = *m.IDNumber
	}
	return d
}

func gormBankAccountToDomain(m *domain.CompanyBankAccountGORM) *domain.CompanyBankAccount {
	d := &domain.CompanyBankAccount{
		ID:            m.ID,
		CompanyID:     m.CompanyID,
		BankName:      m.BankName,
		AccountNumber: m.AccountNumber,
		AccountHolder: m.AccountName,
		Currency:      m.Currency,
		IsDefault:     m.IsPrimary,
		Status:        domain.BankAccountActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	if !m.IsActive {
		d.Status = domain.BankAccountClosed
	}
	if m.BankCode != nil {
		d.BankCode = *m.BankCode
	}
	if m.Branch != nil {
		d.BranchName = *m.Branch
	}
	return d
}

func gormEInvoicePatternToDomain(m *domain.EInvoicePatternGORM) *domain.EInvoicePattern {
	d := &domain.EInvoicePattern{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		PatternCode: m.Pattern,
		Serial:      m.Serial,
		InvoiceType: m.InvoiceType,
		Status:      domain.EInvoiceActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if !m.IsActive {
		d.Status = domain.EInvoiceCancelled
	}
	return d
}

func gormDigitalSignatureToDomain(m *domain.DigitalSignatureGORM) *domain.DigitalSignature {
	d := &domain.DigitalSignature{
		ID:           m.ID,
		CompanyID:    m.CompanyID,
		SerialNumber: m.SerialNumber,
		ValidFrom:    m.ValidFrom.Format("2006-01-02"),
		ValidTo:      m.ValidTo.Format("2006-01-02"),
		Status:       domain.SignatureActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
	if !m.IsActive {
		d.Status = domain.SignatureExpired
	}
	if m.Issuer != nil {
		d.CertificateIssuer = *m.Issuer
	}
	return d
}

func gormIntegrationToDomain(m *domain.IntegrationProfileGORM) *domain.IntegrationProfile {
	d := &domain.IntegrationProfile{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		IntegrationType: domain.IntegrationType(m.Type),
		Status:          domain.IntegrationConnected,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
	if !m.IsActive {
		d.Status = domain.IntegrationDisconnected
	}
	return d
}

// ─── Helpers ────────────────────────────────────────────────────────

