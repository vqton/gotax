package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type PGFARepo struct{ db *gorm.DB }

func NewPGFARepo(db *gorm.DB) *PGFARepo { return &PGFARepo{db} }

// ─── Categories ──────────────────────────────────────────────────────

func (r *PGFARepo) CreateCategory(ctx context.Context, c *domain.FixedAssetCategory) error {
	m := gormFACategoryToGORM(c)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGFARepo) GetCategoryByID(ctx context.Context, id string) (*domain.FixedAssetCategory, error) {
	var m domain.FixedAssetCategoryGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormFACategoryToDomain(&m), nil
}

func (r *PGFARepo) GetCategoryByCode(ctx context.Context, companyID, code string) (*domain.FixedAssetCategory, error) {
	var m domain.FixedAssetCategoryGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		return nil, err
	}
	return gormFACategoryToDomain(&m), nil
}

func (r *PGFARepo) ListCategories(ctx context.Context, filter domain.FACategoryFilter) ([]domain.FixedAssetCategory, error) {
	q := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.ParentID != nil {
		q = q.Where("parent_id = ?", *filter.ParentID)
	}
	if filter.Level != nil {
		q = q.Where("level = ?", *filter.Level)
	}
	var models []domain.FixedAssetCategoryGORM
	if err := q.Order("code").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.FixedAssetCategory, len(models))
	for i := range models {
		out[i] = *gormFACategoryToDomain(&models[i])
	}
	return out, nil
}

var faCategoryUpdateCols = []string{
	"company_id", "code", "name", "parent_id", "level",
	"default_useful_life_months", "default_depreciation_method",
	"asset_account_id", "depreciation_account_id", "expense_account_id",
}

func (r *PGFARepo) UpdateCategory(ctx context.Context, c *domain.FixedAssetCategory) error {
	m := gormFACategoryToGORM(c)
	return r.db.WithContext(ctx).Model(&domain.FixedAssetCategoryGORM{}).
		Select(faCategoryUpdateCols).Where("id = ?", c.ID).Updates(m).Error
}

func (r *PGFARepo) DeleteCategory(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.FixedAssetCategoryGORM{}).Error
}

// ─── Fixed Assets ────────────────────────────────────────────────────

func (r *PGFARepo) CreateAsset(ctx context.Context, a *domain.FixedAsset) error {
	m := gormFAssetToGORM(a)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGFARepo) GetAssetByID(ctx context.Context, id string) (*domain.FixedAsset, error) {
	var m domain.FixedAssetGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormFAssetToDomain(&m), nil
}

func (r *PGFARepo) GetAssetByCode(ctx context.Context, companyID, code string) (*domain.FixedAsset, error) {
	var m domain.FixedAssetGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		return nil, err
	}
	return gormFAssetToDomain(&m), nil
}

func (r *PGFARepo) ListAssets(ctx context.Context, filter domain.FAListFilter) ([]domain.FixedAsset, int, error) {
	q := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.CategoryID != nil {
		q = q.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.DepartmentID != nil {
		q = q.Where("department_id = ?", *filter.DepartmentID)
	}
	if filter.Keyword != "" {
		q = q.Where("(code ILIKE ? OR name ILIKE ?)", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	var total int64
	q.Model(&domain.FixedAssetGORM{}).Count(&total)
	var models []domain.FixedAssetGORM
	if err := q.Order("code").Limit(filter.Limit).Offset(filter.Offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.FixedAsset, len(models))
	for i := range models {
		out[i] = *gormFAssetToDomain(&models[i])
	}
	return out, int(total), nil
}

var faAssetUpdateCols = []string{
	"company_id", "code", "name", "category_id", "status", "acquisition_date",
	"original_cost", "accumulated_depreciation", "residual_value", "carrying_amount",
	"useful_life_months", "depreciation_method", "depreciation_start_date",
	"depreciation_end_date", "department_id", "location", "user_id", "supplier_id",
	"contract_no", "invoice_id", "serial_no", "manufacturer", "manufacture_year",
	"country_of_origin", "technical_specs", "notes", "source", "cip_account_id",
	"asset_account_id", "depreciation_account_id", "expense_account_id", "updated_by",
}

func (r *PGFARepo) UpdateAsset(ctx context.Context, a *domain.FixedAsset) error {
	m := gormFAssetToGORM(a)
	return r.db.WithContext(ctx).Model(&domain.FixedAssetGORM{}).
		Select(faAssetUpdateCols).Where("id = ?", a.ID).Updates(m).Error
}

func (r *PGFARepo) DeleteAsset(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.FixedAssetGORM{}).Error
}

// ─── Depreciation ────────────────────────────────────────────────────

func (r *PGFARepo) CreateDepreciationEntry(ctx context.Context, e *domain.DepreciationEntry) error {
	m := gormDeprEntryToGORM(e)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGFARepo) GetDepreciationEntry(ctx context.Context, id string) (*domain.DepreciationEntry, error) {
	var m domain.DepreciationEntryGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormDeprEntryToDomain(&m), nil
}

func (r *PGFARepo) ListDepreciationByAsset(ctx context.Context, assetID string) ([]domain.DepreciationEntry, error) {
	var models []domain.DepreciationEntryGORM
	if err := r.db.WithContext(ctx).Where("fixed_asset_id = ?", assetID).Order("period_year, period_month").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.DepreciationEntry, len(models))
	for i := range models {
		out[i] = *gormDeprEntryToDomain(&models[i])
	}
	return out, nil
}

func (r *PGFARepo) ListDepreciationByPeriod(ctx context.Context, periodID string) ([]domain.DepreciationEntry, error) {
	var models []domain.DepreciationEntryGORM
	if err := r.db.WithContext(ctx).Where("period_id = ?", periodID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.DepreciationEntry, len(models))
	for i := range models {
		out[i] = *gormDeprEntryToDomain(&models[i])
	}
	return out, nil
}

func (r *PGFARepo) DepreciationExistsForPeriod(ctx context.Context, assetID, periodID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.DepreciationEntryGORM{}).
		Where("fixed_asset_id = ? AND period_id = ?", assetID, periodID).
		Count(&count).Error
	return count > 0, err
}

func (r *PGFARepo) DeleteDepreciationEntry(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.DepreciationEntryGORM{}).Error
}

// ─── Transactions ────────────────────────────────────────────────────

func (r *PGFARepo) CreateTransaction(ctx context.Context, t *domain.FixedAssetTransaction) error {
	m := gormFATransactionToGORM(t)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGFARepo) ListTransactionsByAsset(ctx context.Context, assetID string) ([]domain.FixedAssetTransaction, error) {
	var models []domain.FixedAssetTransactionGORM
	if err := r.db.WithContext(ctx).Where("fixed_asset_id = ?", assetID).Order("transaction_date, created_at").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.FixedAssetTransaction, len(models))
	for i := range models {
		out[i] = *gormFATransactionToDomain(&models[i])
	}
	return out, nil
}

// ─── Allocations ─────────────────────────────────────────────────────

func (r *PGFARepo) SetAllocations(ctx context.Context, assetID string, allocs []domain.FixedAssetAllocation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("fixed_asset_id = ?", assetID).Delete(&domain.FixedAssetAllocationGORM{}).Error; err != nil {
			return err
		}
		for i := range allocs {
			m := gormFAAllocationToGORM(&allocs[i])
			if err := tx.Create(m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PGFARepo) GetAllocations(ctx context.Context, assetID string) ([]domain.FixedAssetAllocation, error) {
	var models []domain.FixedAssetAllocationGORM
	if err := r.db.WithContext(ctx).Where("fixed_asset_id = ?", assetID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.FixedAssetAllocation, len(models))
	for i := range models {
		out[i] = *gormFAAllocationToDomain(&models[i])
	}
	return out, nil
}

// ─── Inventory Plans ─────────────────────────────────────────────────

func (r *PGFARepo) CreateInventoryPlan(ctx context.Context, p *domain.FixedAssetInventoryPlan) error {
	m := gormFAIPlanToGORM(p)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGFARepo) GetInventoryPlan(ctx context.Context, id string) (*domain.FixedAssetInventoryPlan, error) {
	var m domain.FixedAssetInventoryPlanGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormFAIPlanToDomain(&m), nil
}

func (r *PGFARepo) ListInventoryPlans(ctx context.Context, companyID string) ([]domain.FixedAssetInventoryPlan, error) {
	var models []domain.FixedAssetInventoryPlanGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("plan_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.FixedAssetInventoryPlan, len(models))
	for i := range models {
		out[i] = *gormFAIPlanToDomain(&models[i])
	}
	return out, nil
}

func (r *PGFARepo) UpdateInventoryPlan(ctx context.Context, p *domain.FixedAssetInventoryPlan) error {
	m := gormFAIPlanToGORM(p)
	return r.db.WithContext(ctx).Model(&domain.FixedAssetInventoryPlanGORM{}).
		Select("company_id", "plan_date", "status", "notes").Where("id = ?", p.ID).Updates(m).Error
}

// ─── Inventory Results ───────────────────────────────────────────────

func (r *PGFARepo) CreateInventoryResult(ctx context.Context, res *domain.FixedAssetInventoryResult) error {
	m := gormFAIResultToGORM(res)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGFARepo) GetInventoryResultsByPlan(ctx context.Context, planID string) ([]domain.FixedAssetInventoryResult, error) {
	var models []domain.FixedAssetInventoryResultGORM
	if err := r.db.WithContext(ctx).Where("plan_id = ?", planID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.FixedAssetInventoryResult, len(models))
	for i := range models {
		out[i] = *gormFAIResultToDomain(&models[i])
	}
	return out, nil
}

// ─── Conversion helpers ─────────────────────────────────────────────

func gormFACategoryToDomain(m *domain.FixedAssetCategoryGORM) *domain.FixedAssetCategory {
	return &domain.FixedAssetCategory{
		ID:                          m.ID,
		CompanyID:                   m.CompanyID,
		Code:                        m.Code,
		Name:                        m.Name,
		ParentID:                    m.ParentID,
		Level:                       m.Level,
		DefaultUsefulLifeMonths:     m.DefaultUsefulLifeMonths,
		DefaultDepreciationMethod:   domain.DepreciationMethod(m.DefaultDepreciationMethod),
		AssetAccountID:              m.AssetAccountID,
		DepreciationAccountID:       m.DepreciationAccountID,
		ExpenseAccountID:            m.ExpenseAccountID,
		CreatedAt:                   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                   m.UpdatedAt.Format(time.RFC3339),
	}
}

func gormFACategoryToGORM(d *domain.FixedAssetCategory) *domain.FixedAssetCategoryGORM {
	return &domain.FixedAssetCategoryGORM{
		ID:                          d.ID,
		CompanyID:                   d.CompanyID,
		Code:                        d.Code,
		Name:                        d.Name,
		ParentID:                    d.ParentID,
		Level:                       d.Level,
		DefaultUsefulLifeMonths:     d.DefaultUsefulLifeMonths,
		DefaultDepreciationMethod:   string(d.DefaultDepreciationMethod),
		AssetAccountID:              d.AssetAccountID,
		DepreciationAccountID:       d.DepreciationAccountID,
		ExpenseAccountID:            d.ExpenseAccountID,
	}
}

func gormFAssetToDomain(m *domain.FixedAssetGORM) *domain.FixedAsset {
	return &domain.FixedAsset{
		ID:                      m.ID,
		CompanyID:               m.CompanyID,
		Code:                    m.Code,
		Name:                    m.Name,
		CategoryID:              m.CategoryID,
		Status:                  domain.FixedAssetStatus(m.Status),
		AcquisitionDate:         m.AcquisitionDate.Format(time.DateOnly),
		OriginalCost:            m.OriginalCost,
		AccumulatedDepreciation: m.AccumulatedDepreciation,
		ResidualValue:           m.ResidualValue,
		CarryingAmount:          m.CarryingAmount,
		UsefulLifeMonths:        m.UsefulLifeMonths,
		DepreciationMethod:      domain.DepreciationMethod(m.DepreciationMethod),
		DepreciationStartDate:   timePtrToStringPtr(m.DepreciationStartDate),
		DepreciationEndDate:     timePtrToStringPtr(m.DepreciationEndDate),
		DepartmentID:            m.DepartmentID,
		Location:                m.Location,
		UserID:                  safeStr(m.UserID),
		SupplierID:              safeStr(m.SupplierID),
		ContractNo:              safeStr(m.ContractNo),
		InvoiceID:               safeStr(m.InvoiceID),
		SerialNo:                safeStr(m.SerialNo),
		Manufacturer:            safeStr(m.Manufacturer),
		ManufactureYear:         m.ManufactureYear,
		CountryOfOrigin:         safeStr(m.CountryOfOrigin),
		TechnicalSpecs:          safeStr(m.TechnicalSpecs),
		Notes:                   safeStr(m.Notes),
		Source:                  domain.FASource(m.Source),
		CIPAccountID:            safeStr(m.CIPAccountID),
		AssetAccountID:          m.AssetAccountID,
		DepreciationAccountID:   m.DepreciationAccountID,
		ExpenseAccountID:        m.ExpenseAccountID,
		CreatedAt:               m.CreatedAt.Format(time.RFC3339),
		CreatedBy:               m.CreatedBy,
		UpdatedAt:               m.UpdatedAt.Format(time.RFC3339),
		UpdatedBy:               m.UpdatedBy,
	}
}

func gormFAssetToGORM(d *domain.FixedAsset) *domain.FixedAssetGORM {
	return &domain.FixedAssetGORM{
		ID:                      d.ID,
		CompanyID:               d.CompanyID,
		Code:                    d.Code,
		Name:                    d.Name,
		CategoryID:              d.CategoryID,
		Status:                  string(d.Status),
		AcquisitionDate:         parseDate(d.AcquisitionDate),
		OriginalCost:            d.OriginalCost,
		AccumulatedDepreciation: d.AccumulatedDepreciation,
		ResidualValue:           d.ResidualValue,
		CarryingAmount:          d.CarryingAmount,
		UsefulLifeMonths:        d.UsefulLifeMonths,
		DepreciationMethod:      string(d.DepreciationMethod),
		DepreciationStartDate:   stringPtrToTimePtr(d.DepreciationStartDate),
		DepreciationEndDate:     stringPtrToTimePtr(d.DepreciationEndDate),
		DepartmentID:            d.DepartmentID,
		Location:                d.Location,
		UserID:                  strPtr(d.UserID),
		SupplierID:              strPtr(d.SupplierID),
		ContractNo:              strPtr(d.ContractNo),
		InvoiceID:               strPtr(d.InvoiceID),
		SerialNo:                strPtr(d.SerialNo),
		Manufacturer:            strPtr(d.Manufacturer),
		ManufactureYear:         d.ManufactureYear,
		CountryOfOrigin:         strPtr(d.CountryOfOrigin),
		TechnicalSpecs:          strPtr(d.TechnicalSpecs),
		Notes:                   strPtr(d.Notes),
		Source:                  string(d.Source),
		CIPAccountID:            strPtr(d.CIPAccountID),
		AssetAccountID:          d.AssetAccountID,
		DepreciationAccountID:   d.DepreciationAccountID,
		ExpenseAccountID:        d.ExpenseAccountID,
		CreatedBy:               d.CreatedBy,
		UpdatedBy:               d.UpdatedBy,
	}
}

func gormDeprEntryToDomain(m *domain.DepreciationEntryGORM) *domain.DepreciationEntry {
	entry := &domain.DepreciationEntry{
		ID:                  m.ID,
		CompanyID:           m.CompanyID,
		FixedAssetID:        m.FixedAssetID,
		PeriodID:            m.PeriodID,
		PeriodYear:          m.PeriodYear,
		PeriodMonth:         m.PeriodMonth,
		DepreciationAmount:  m.DepreciationAmount,
		AccumulatedAfter:    m.AccumulatedAfter,
		CarryingAmountAfter: m.CarryingAmountAfter,
		GLPosted:            m.GLPosted,
		CreatedBy:           m.CreatedBy,
		CreatedAt:           m.CreatedAt.Format(time.RFC3339),
	}
	if m.GLJournalEntryID != nil {
		entry.GLJournalEntryID = *m.GLJournalEntryID
	}
	return entry
}

func gormDeprEntryToGORM(d *domain.DepreciationEntry) *domain.DepreciationEntryGORM {
	return &domain.DepreciationEntryGORM{
		ID:                  d.ID,
		CompanyID:           d.CompanyID,
		FixedAssetID:        d.FixedAssetID,
		PeriodID:            d.PeriodID,
		PeriodYear:          d.PeriodYear,
		PeriodMonth:         d.PeriodMonth,
		DepreciationAmount:  d.DepreciationAmount,
		AccumulatedAfter:    d.AccumulatedAfter,
		CarryingAmountAfter: d.CarryingAmountAfter,
		GLPosted:            d.GLPosted,
		GLJournalEntryID:    strPtr(d.GLJournalEntryID),
		CreatedBy:           d.CreatedBy,
	}
}

func gormFATransactionToDomain(m *domain.FixedAssetTransactionGORM) *domain.FixedAssetTransaction {
	return &domain.FixedAssetTransaction{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		FixedAssetID:    m.FixedAssetID,
		TransactionType: domain.FATransactionType(m.TransactionType),
		TransactionDate: m.TransactionDate.Format(time.DateOnly),
		Amount:          m.Amount,
		OldValue:        m.OldValue,
		NewValue:        m.NewValue,
		Description:     safeStr(m.Description),
		GLJournalID:     safeStr(m.GLJournalID),
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt.Format(time.RFC3339),
	}
}

func gormFATransactionToGORM(d *domain.FixedAssetTransaction) *domain.FixedAssetTransactionGORM {
	return &domain.FixedAssetTransactionGORM{
		ID:              d.ID,
		CompanyID:       d.CompanyID,
		FixedAssetID:    d.FixedAssetID,
		TransactionType: string(d.TransactionType),
		TransactionDate: parseDate(d.TransactionDate),
		Amount:          d.Amount,
		OldValue:        d.OldValue,
		NewValue:        d.NewValue,
		Description:     strPtr(d.Description),
		GLJournalID:     strPtr(d.GLJournalID),
		CreatedBy:       d.CreatedBy,
	}
}

func gormFAAllocationToDomain(m *domain.FixedAssetAllocationGORM) *domain.FixedAssetAllocation {
	return &domain.FixedAssetAllocation{
		ID:               m.ID,
		FixedAssetID:     m.FixedAssetID,
		DepartmentID:     m.DepartmentID,
		AllocationPct:    m.AllocationPct,
		ExpenseAccountID: m.ExpenseAccountID,
	}
}

func gormFAAllocationToGORM(d *domain.FixedAssetAllocation) *domain.FixedAssetAllocationGORM {
	return &domain.FixedAssetAllocationGORM{
		ID:               d.ID,
		FixedAssetID:     d.FixedAssetID,
		DepartmentID:     d.DepartmentID,
		AllocationPct:    d.AllocationPct,
		ExpenseAccountID: d.ExpenseAccountID,
	}
}

func gormFAIPlanToDomain(m *domain.FixedAssetInventoryPlanGORM) *domain.FixedAssetInventoryPlan {
	return &domain.FixedAssetInventoryPlan{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		PlanDate:  m.PlanDate.Format(time.DateOnly),
		Status:    m.Status,
		Notes:     safeStr(m.Notes),
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
}

func gormFAIPlanToGORM(d *domain.FixedAssetInventoryPlan) *domain.FixedAssetInventoryPlanGORM {
	return &domain.FixedAssetInventoryPlanGORM{
		ID:        d.ID,
		CompanyID: d.CompanyID,
		PlanDate:  parseDate(d.PlanDate),
		Status:    d.Status,
		Notes:     strPtr(d.Notes),
		CreatedBy: d.CreatedBy,
	}
}

func gormFAIResultToDomain(m *domain.FixedAssetInventoryResultGORM) *domain.FixedAssetInventoryResult {
	return &domain.FixedAssetInventoryResult{
		ID:               m.ID,
		PlanID:           m.PlanID,
		FixedAssetID:     m.FixedAssetID,
		ExpectedLocation: safeStr(m.ExpectedLocation),
		ActualLocation:   safeStr(m.ActualLocation),
		ExpectedStatus:   safeStr(m.ExpectedStatus),
		ActualStatus:     safeStr(m.ActualStatus),
		Discrepancy:      m.Discrepancy,
		Notes:            safeStr(m.Notes),
	}
}

func gormFAIResultToGORM(d *domain.FixedAssetInventoryResult) *domain.FixedAssetInventoryResultGORM {
	return &domain.FixedAssetInventoryResultGORM{
		ID:               d.ID,
		PlanID:           d.PlanID,
		FixedAssetID:     d.FixedAssetID,
		ExpectedLocation: strPtr(d.ExpectedLocation),
		ActualLocation:   strPtr(d.ActualLocation),
		ExpectedStatus:   strPtr(d.ExpectedStatus),
		ActualStatus:     strPtr(d.ActualStatus),
		Discrepancy:      d.Discrepancy,
		Notes:            strPtr(d.Notes),
	}
}

// ─── FA-specific helpers ────────────────────────────────────────────

func timePtrToStringPtr(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	s := t.Format(time.DateOnly)
	return &s
}

func stringPtrToTimePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse(time.DateOnly, *s)
	if err != nil {
		return nil
	}
	return &t
}
