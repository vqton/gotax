package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type PGOpeningBalanceRepo struct {
	db *gorm.DB
}

func NewPGOpeningBalanceRepo(db *gorm.DB) *PGOpeningBalanceRepo {
	return &PGOpeningBalanceRepo{db: db}
}

// ─── Opening Balance ────────────────────────────────────────────────

func (r *PGOpeningBalanceRepo) Create(ctx context.Context, ob *domain.OpeningBalance) error {
	if ob.ID == "" {
		ob.ID = "OB" + time.Now().Format("20060102150405.000000000")
	}
	if ob.OriginalAmount == 0 {
		if ob.DebitAmount > 0 {
			ob.OriginalAmount = ob.DebitAmount
		} else if ob.CreditAmount > 0 {
			ob.OriginalAmount = ob.CreditAmount
		}
	}
	m := &domain.OpeningBalanceGORM{
		ID:             ob.ID,
		CompanyID:      ob.CompanyID,
		PeriodID:       ob.PeriodID,
		AccountCode:    ob.AccountCode,
		CurrencyCode:   ob.CurrencyCode,
		OriginalAmount: ob.OriginalAmount,
		DebitAmount:    ob.DebitAmount,
		CreditAmount:   ob.CreditAmount,
		ExchangeRate:   ob.ExchangeRate,
		Status:         string(ob.Status),
		SourceType:     ob.SourceType,
		BatchID:        strPtr(ob.BatchID),
		Reason:         strPtr(ob.Reason),
		CreatedBy:      ob.CreatedBy,
	}
	if ob.FiscalYearID != "" {
		m.FiscalYearID = strPtr(ob.FiscalYearID)
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGOpeningBalanceRepo) GetByID(ctx context.Context, id string) (*domain.OpeningBalance, error) {
	var m domain.OpeningBalanceGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormOBToDomain(&m), nil
}

func (r *PGOpeningBalanceRepo) List(ctx context.Context, filter domain.OBListFilter) ([]domain.OpeningBalance, error) {
	q := r.db.WithContext(ctx).Where("company_id = ?", filter.CompanyID)
	if filter.PeriodID != "" {
		q = q.Where("period_id = ?", filter.PeriodID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.AccountCode != "" {
		q = q.Where("account_code = ?", filter.AccountCode)
	}
	if filter.Currency != "" {
		q = q.Where("currency_code = ?", filter.Currency)
	}
	var ms []domain.OpeningBalanceGORM
	if err := q.Order("account_code ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.OpeningBalance, len(ms))
	for i := range ms {
		out[i] = *gormOBToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGOpeningBalanceRepo) GetByAccount(ctx context.Context, companyID, periodID, accountCode string) (*domain.OpeningBalance, error) {
	var m domain.OpeningBalanceGORM
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND period_id = ? AND account_code = ? AND status = 'APPROVED'", companyID, periodID, accountCode).
		First(&m).Error; err != nil {
		return nil, err
	}
	return gormOBToDomain(&m), nil
}

func (r *PGOpeningBalanceRepo) Update(ctx context.Context, ob *domain.OpeningBalance) error {
	return r.db.WithContext(ctx).Model(&domain.OpeningBalanceGORM{}).
		Where("id = ? AND status NOT IN ('APPROVED','CORRECTED')", ob.ID).
		Updates(map[string]any{
			"account_code":  ob.AccountCode,
			"currency_code": ob.CurrencyCode,
			"original_amount": ob.OriginalAmount,
			"debit_amount":  ob.DebitAmount,
			"credit_amount": ob.CreditAmount,
			"exchange_rate": ob.ExchangeRate,
			"source_type":   ob.SourceType,
			"batch_id":      ob.BatchID,
			"reason":        ob.Reason,
			"fiscal_year_id": ob.FiscalYearID,
			"updated_at":    time.Now(),
		}).Error
}

func (r *PGOpeningBalanceRepo) UpdateStatus(ctx context.Context, id string, status domain.OpeningBalanceStatus, approvedBy string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.OpeningBalanceGORM{}).Where("id = ?", id).
		Updates(map[string]any{
			"status":     string(status),
			"approved_by": gorm.Expr("CASE WHEN ? = 'APPROVED' THEN ? ELSE approved_by END", string(status), approvedBy),
			"approved_at": gorm.Expr("CASE WHEN ? = 'APPROVED' THEN ? ELSE approved_at END", string(status), now),
			"updated_at": now,
		}).Error
}

func (r *PGOpeningBalanceRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND status NOT IN ('APPROVED','CORRECTED')", id).
		Delete(&domain.OpeningBalanceGORM{}).Error
}

func (r *PGOpeningBalanceRepo) BulkCreate(ctx context.Context, balances []domain.OpeningBalance) error {
	now := time.Now()
	ms := make([]domain.OpeningBalanceGORM, len(balances))
	for i, ob := range balances {
		if ob.ID == "" {
			ob.ID = "OB" + time.Now().Format("20060102150405.000000000")
		}
		if ob.OriginalAmount == 0 {
			if ob.DebitAmount > 0 {
				ob.OriginalAmount = ob.DebitAmount
			} else if ob.CreditAmount > 0 {
				ob.OriginalAmount = ob.CreditAmount
			}
		}
		ms[i] = domain.OpeningBalanceGORM{
			ID:             ob.ID,
			CompanyID:      ob.CompanyID,
			PeriodID:       ob.PeriodID,
			FiscalYearID:   strPtr(ob.FiscalYearID),
			AccountCode:    ob.AccountCode,
			CurrencyCode:   ob.CurrencyCode,
			OriginalAmount: ob.OriginalAmount,
			DebitAmount:    ob.DebitAmount,
			CreditAmount:   ob.CreditAmount,
			ExchangeRate:   ob.ExchangeRate,
			Status:         string(ob.Status),
			SourceType:     ob.SourceType,
			BatchID:        strPtr(ob.BatchID),
			Reason:         strPtr(ob.Reason),
			CreatedBy:      ob.CreatedBy,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
	return r.db.WithContext(ctx).Create(&ms).Error
}

func (r *PGOpeningBalanceRepo) BulkUpdateStatus(ctx context.Context, ids []string, status domain.OpeningBalanceStatus, approvedBy string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.OpeningBalanceGORM{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"status":     string(status),
			"approved_by": gorm.Expr("CASE WHEN ? = 'APPROVED' THEN ? ELSE approved_by END", string(status), approvedBy),
			"approved_at": gorm.Expr("CASE WHEN ? = 'APPROVED' THEN ? ELSE approved_at END", string(status), now),
			"updated_at": now,
		}).Error
}

func (r *PGOpeningBalanceRepo) CreateDetail(ctx context.Context, d *domain.OpeningBalanceDetail) error {
	if d.ID == "" {
		d.ID = "OBD" + time.Now().Format("20060102150405.000000000")
	}
	m := &domain.OpeningBalanceDetailGORM{
		ID:                d.ID,
		OpeningBalanceID:  d.OpeningBalanceID,
		EntityType:        string(d.EntityType),
		EntityID:          d.EntityID,
		EntityName:        strPtr(d.EntityName),
		DebitAmount:       d.DebitAmount,
		CreditAmount:      d.CreditAmount,
		Quantity:          f64Ptr(d.Quantity),
		UnitPrice:         f64Ptr(d.UnitPrice),
		OriginalCost:      f64Ptr(d.OriginalCost),
		AccDepreciation:   f64Ptr(d.AccDepreciation),
		CounterpartAccount: strPtr(d.CounterpartAccount),
		Note:              strPtr(d.Note),
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGOpeningBalanceRepo) GetDetails(ctx context.Context, balanceID string) ([]domain.OpeningBalanceDetail, error) {
	var ms []domain.OpeningBalanceDetailGORM
	if err := r.db.WithContext(ctx).
		Where("opening_balance_id = ?", balanceID).
		Order("entity_type, entity_id ASC").
		Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.OpeningBalanceDetail, len(ms))
	for i := range ms {
		out[i] = *gormOBDetailToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGOpeningBalanceRepo) DeleteDetail(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.OpeningBalanceDetailGORM{}).Error
}

func (r *PGOpeningBalanceRepo) DeleteDetails(ctx context.Context, balanceID string) error {
	return r.db.WithContext(ctx).Where("opening_balance_id = ?", balanceID).Delete(&domain.OpeningBalanceDetailGORM{}).Error
}

func (r *PGOpeningBalanceRepo) GetTotals(ctx context.Context, companyID, periodID string) (float64, float64, error) {
	var result struct {
		TotalDebit  float64
		TotalCredit float64
	}
	err := r.db.WithContext(ctx).Raw(
		`SELECT COALESCE(SUM(debit_amount),0) AS total_debit, COALESCE(SUM(credit_amount),0) AS total_credit
		 FROM opening_balances WHERE company_id = ? AND period_id = ? AND status = 'APPROVED'`,
		companyID, periodID).Scan(&result).Error
	return result.TotalDebit, result.TotalCredit, err
}

func (r *PGOpeningBalanceRepo) ValidateBalanced(ctx context.Context, companyID, periodID string) (bool, error) {
	var diff float64
	err := r.db.WithContext(ctx).Raw(
		`SELECT COALESCE(ABS(SUM(debit_amount)-SUM(credit_amount)),0)
		 FROM opening_balances WHERE company_id = ? AND period_id = ? AND status = 'APPROVED'`,
		companyID, periodID).Scan(&diff).Error
	return diff < 0.01, err
}

// ─── Carry Forward Log ──────────────────────────────────────────────

func (r *PGOpeningBalanceRepo) CreateCarryForwardLog(ctx context.Context, log *domain.CarryForwardLog) error {
	if log.ID == "" {
		log.ID = "CF" + time.Now().Format("20060102150405.000000000")
	}
	m := &domain.CarryForwardLogGORM{
		ID:             log.ID,
		CompanyID:      log.CompanyID,
		FromPeriodID:   log.FromPeriodID,
		ToPeriodID:     log.ToPeriodID,
		FromFiscalYear: log.FromFiscalYear,
		ToFiscalYear:   log.ToFiscalYear,
		AccountCount:   log.AccountCount,
		TotalDebit:     log.TotalDebit,
		TotalCredit:    log.TotalCredit,
		Status:         log.Status,
		ExecutedBy:     log.ExecutedBy,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGOpeningBalanceRepo) GetCarryForwardLogs(ctx context.Context, companyID string) ([]domain.CarryForwardLog, error) {
	var ms []domain.CarryForwardLogGORM
	if err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("executed_at DESC").
		Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CarryForwardLog, len(ms))
	for i := range ms {
		out[i] = *gormCFLogToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGOpeningBalanceRepo) GetCarryForwardLogByID(ctx context.Context, id string) (*domain.CarryForwardLog, error) {
	var m domain.CarryForwardLogGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrOpeningBalanceNotFound
	}
	return gormCFLogToDomain(&m), nil
}

// ─── Circular 99 Mapping ───────────────────────────────────────────

func (r *PGOpeningBalanceRepo) CreateCircular99Mapping(ctx context.Context, m *domain.Circular99Mapping) error {
	if m.ID == "" {
		m.ID = "C99" + time.Now().Format("20060102150405.000000000")
	}
	ed, _ := time.Parse("2006-01-02", m.EffectiveDate)
	gormM := &domain.Circular99MappingGORM{
		ID:                 m.ID,
		OldAccountCode:     m.OldAccountCode,
		NewAccountCode:     m.NewAccountCode,
		MappingType:        m.MappingType,
		SplitRatio:         f64Ptr(m.SplitRatio),
		CounterpartAccount: strPtr(m.CounterpartAccount),
		EffectiveDate:      ed,
		Note:               strPtr(m.Note),
		IsActive:           m.IsActive,
	}
	return r.db.WithContext(ctx).Create(gormM).Error
}

func (r *PGOpeningBalanceRepo) ListCircular99Mappings(ctx context.Context) ([]domain.Circular99Mapping, error) {
	var ms []domain.Circular99MappingGORM
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("old_account_code ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Circular99Mapping, len(ms))
	for i := range ms {
		out[i] = *gormC99MappingToDomain(&ms[i])
	}
	return out, nil
}

func (r *PGOpeningBalanceRepo) GetCircular99MappingByOldCode(ctx context.Context, oldCode string) (*domain.Circular99Mapping, error) {
	var m domain.Circular99MappingGORM
	if err := r.db.WithContext(ctx).Where("old_account_code = ? AND is_active = ?", oldCode, true).First(&m).Error; err != nil {
		return nil, domain.ErrCircular99MappingNotFound
	}
	return gormC99MappingToDomain(&m), nil
}

// ─── Balance Migration ─────────────────────────────────────────────

func (r *PGOpeningBalanceRepo) CreateMigration(ctx context.Context, m *domain.BalanceMigration) error {
	if m.ID == "" {
		m.ID = "MIG" + time.Now().Format("20060102150405.000000000")
	}
	ed, _ := time.Parse("2006-01-02", m.ExecutionDate)
	gormM := &domain.BalanceMigrationGORM{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		FromRegime:      m.FromRegime,
		ToRegime:        m.ToRegime,
		ExecutionDate:   ed,
		Status:          m.Status,
		SourceBalanceID: strPtr(m.SourceBalanceID),
		TargetBalanceID: strPtr(m.TargetBalanceID),
		JournalEntryID:  strPtr(m.JournalEntryID),
		Summary:         strPtr(m.Summary),
		ExecutedBy:      m.ExecutedBy,
	}
	return r.db.WithContext(ctx).Create(gormM).Error
}

func (r *PGOpeningBalanceRepo) GetMigrationByID(ctx context.Context, id string) (*domain.BalanceMigration, error) {
	var m domain.BalanceMigrationGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrOpeningBalanceNotFound
	}
	return gormMigrationToDomain(&m), nil
}

func (r *PGOpeningBalanceRepo) ListMigrations(ctx context.Context, companyID string) ([]domain.BalanceMigration, error) {
	var ms []domain.BalanceMigrationGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.BalanceMigration, len(ms))
	for i := range ms {
		out[i] = *gormMigrationToDomain(&ms[i])
	}
	return out, nil
}

// ─── Converters ─────────────────────────────────────────────────────

func gormOBToDomain(m *domain.OpeningBalanceGORM) *domain.OpeningBalance {
	d := &domain.OpeningBalance{
		ID:             m.ID,
		CompanyID:      m.CompanyID,
		PeriodID:       m.PeriodID,
		AccountCode:    m.AccountCode,
		CurrencyCode:   m.CurrencyCode,
		OriginalAmount: m.OriginalAmount,
		DebitAmount:    m.DebitAmount,
		CreditAmount:   m.CreditAmount,
		ExchangeRate:   m.ExchangeRate,
		Status:         domain.OpeningBalanceStatus(m.Status),
		SourceType:     m.SourceType,
		CreatedBy:      m.CreatedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	if m.FiscalYearID != nil {
		d.FiscalYearID = *m.FiscalYearID
	}
	if m.BatchID != nil {
		d.BatchID = *m.BatchID
	}
	if m.Reason != nil {
		d.Reason = *m.Reason
	}
	if m.ApprovedBy != nil {
		d.ApprovedBy = *m.ApprovedBy
	}
	d.ApprovedAt = m.ApprovedAt
	if m.CorrectedBy != nil {
		d.CorrectedBy = *m.CorrectedBy
	}
	if m.CorrectionOf != nil {
		d.CorrectionOf = *m.CorrectionOf
	}
	if m.CorrectionReason != nil {
		d.CorrectionReason = *m.CorrectionReason
	}
	return d
}

func gormOBDetailToDomain(m *domain.OpeningBalanceDetailGORM) *domain.OpeningBalanceDetail {
	d := &domain.OpeningBalanceDetail{
		ID:               m.ID,
		OpeningBalanceID: m.OpeningBalanceID,
		EntityType:       domain.DetailEntityType(m.EntityType),
		EntityID:         m.EntityID,
		DebitAmount:      m.DebitAmount,
		CreditAmount:     m.CreditAmount,
		CreatedAt:        m.CreatedAt,
	}
	if m.EntityName != nil {
		d.EntityName = *m.EntityName
	}
	if m.Quantity != nil {
		d.Quantity = *m.Quantity
	}
	if m.UnitPrice != nil {
		d.UnitPrice = *m.UnitPrice
	}
	if m.OriginalCost != nil {
		d.OriginalCost = *m.OriginalCost
	}
	if m.AccDepreciation != nil {
		d.AccDepreciation = *m.AccDepreciation
	}
	if m.CounterpartAccount != nil {
		d.CounterpartAccount = *m.CounterpartAccount
	}
	if m.Note != nil {
		d.Note = *m.Note
	}
	return d
}

func gormCFLogToDomain(m *domain.CarryForwardLogGORM) *domain.CarryForwardLog {
	d := &domain.CarryForwardLog{
		ID:             m.ID,
		CompanyID:      m.CompanyID,
		FromPeriodID:   m.FromPeriodID,
		ToPeriodID:     m.ToPeriodID,
		FromFiscalYear: m.FromFiscalYear,
		ToFiscalYear:   m.ToFiscalYear,
		AccountCount:   m.AccountCount,
		TotalDebit:     m.TotalDebit,
		TotalCredit:    m.TotalCredit,
		Status:         m.Status,
		ExecutedBy:     m.ExecutedBy,
		ExecutedAt:     m.ExecutedAt,
	}
	return d
}

func gormC99MappingToDomain(m *domain.Circular99MappingGORM) *domain.Circular99Mapping {
	d := &domain.Circular99Mapping{
		ID:             m.ID,
		OldAccountCode: m.OldAccountCode,
		NewAccountCode: m.NewAccountCode,
		MappingType:    m.MappingType,
		EffectiveDate:  m.EffectiveDate.Format("2006-01-02"),
		IsActive:       m.IsActive,
	}
	if m.SplitRatio != nil {
		d.SplitRatio = *m.SplitRatio
	}
	if m.CounterpartAccount != nil {
		d.CounterpartAccount = *m.CounterpartAccount
	}
	if m.Note != nil {
		d.Note = *m.Note
	}
	return d
}

func gormMigrationToDomain(m *domain.BalanceMigrationGORM) *domain.BalanceMigration {
	d := &domain.BalanceMigration{
		ID:            m.ID,
		CompanyID:     m.CompanyID,
		FromRegime:    m.FromRegime,
		ToRegime:      m.ToRegime,
		ExecutionDate: m.ExecutionDate.Format("2006-01-02"),
		Status:        m.Status,
		ExecutedBy:    m.ExecutedBy,
		CreatedAt:     m.CreatedAt,
		ExecutedAt:    m.ExecutedAt,
	}
	if m.SourceBalanceID != nil {
		d.SourceBalanceID = *m.SourceBalanceID
	}
	if m.TargetBalanceID != nil {
		d.TargetBalanceID = *m.TargetBalanceID
	}
	if m.JournalEntryID != nil {
		d.JournalEntryID = *m.JournalEntryID
	}
	if m.Summary != nil {
		d.Summary = *m.Summary
	}
	return d
}

// ─── Helpers ────────────────────────────────────────────────────────

func f64Ptr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}
