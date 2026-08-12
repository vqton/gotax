package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type PGAccountRepo struct{ db *gorm.DB }
type PGJournalRepo struct{ db *gorm.DB }
type PGPeriodRepo struct{ db *gorm.DB }
type PGUserRepo struct{ db *gorm.DB }
type PGRefreshTokenRepo struct{ db *gorm.DB }
type PGPasswordResetTokenRepo struct{ db *gorm.DB }
type PGAuditLogRepo struct{ db *gorm.DB }
type PGExchangeRateRepo struct{ db *gorm.DB }
type PGClosingTemplateRepo struct{ db *gorm.DB }

func NewPGAccountRepo(db *gorm.DB) *PGAccountRepo          { return &PGAccountRepo{db} }
func NewPGJournalRepo(db *gorm.DB) *PGJournalRepo           { return &PGJournalRepo{db} }
func NewPGPeriodRepo(db *gorm.DB) *PGPeriodRepo             { return &PGPeriodRepo{db} }
func NewPGUserRepo(db *gorm.DB) *PGUserRepo                { return &PGUserRepo{db} }
func NewPGRefreshTokenRepo(db *gorm.DB) *PGRefreshTokenRepo { return &PGRefreshTokenRepo{db} }
func NewPGPasswordResetTokenRepo(db *gorm.DB) *PGPasswordResetTokenRepo { return &PGPasswordResetTokenRepo{db} }
func NewPGAuditLogRepo(db *gorm.DB) *PGAuditLogRepo       { return &PGAuditLogRepo{db} }
func NewPGExchangeRateRepo(db *gorm.DB) *PGExchangeRateRepo  { return &PGExchangeRateRepo{db} }
func NewPGClosingTemplateRepo(db *gorm.DB) *PGClosingTemplateRepo { return &PGClosingTemplateRepo{db} }


// account methods reuse domain.Account; use Raw or model scans.

func (r *PGAccountRepo) Create(ctx context.Context, a *domain.Account) error {
	// map insert so empty parent_code/detail_by stay NULL (GORM writes ''
	// for zero strings, violating the FK and CHECK constraints)
	data := map[string]interface{}{
		"code": a.Code, "name": a.Name, "type": string(a.Type),
		"is_active": a.IsActive, "is_foreign": a.IsForeign,
		"is_parent": a.IsParent, "status": string(a.Status),
		"arrears_days": a.ArrearsDays,
	}
	if a.Name2 != "" { data["name2"] = a.Name2 }
	if a.ParentCode != "" { data["parent_code"] = a.ParentCode }
	if a.DetailBy != "" { data["detail_by"] = string(a.DetailBy) }
	if a.FreezeReason != "" { data["freeze_reason"] = a.FreezeReason }
	if a.Note != "" { data["note"] = a.Note }
	return r.db.WithContext(ctx).Model(&domain.AccountGORM{}).Create(data).Error
}

func (r *PGAccountRepo) GetByCode(ctx context.Context, code string) (*domain.Account, error) {
	var m domain.AccountGORM
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}
	return gormAccountToDomain(&m), nil
}

func (r *PGAccountRepo) GetAll(ctx context.Context, activeOnly bool) ([]domain.Account, error) {
	q := r.db.WithContext(ctx).Order("code")
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}
	var models []domain.AccountGORM
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Account, len(models))
	for i := range models {
		out[i] = *gormAccountToDomain(&models[i])
	}
	return out, nil
}

func (r *PGAccountRepo) Update(ctx context.Context, a *domain.Account) error {
	// empty parent_code/detail_by omitted so they stay NULL (see Create)
	data := map[string]interface{}{
		"name": a.Name, "name2": a.Name2, "type": string(a.Type),
		"is_active": a.IsActive, "is_foreign": a.IsForeign,
		"is_parent": a.IsParent, "note": a.Note,
		"status": string(a.Status), "freeze_reason": a.FreezeReason,
	}
	if a.ParentCode != "" { data["parent_code"] = a.ParentCode }
	if a.DetailBy != "" { data["detail_by"] = string(a.DetailBy) }
	return r.db.WithContext(ctx).Model(&domain.AccountGORM{}).Where("code = ?", a.Code).Updates(data).Error
}

func (r *PGAccountRepo) Delete(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).Where("code = ?", code).Delete(&domain.AccountGORM{}).Error
}

func (r *PGAccountRepo) GetChildren(ctx context.Context, parentCode string) ([]domain.Account, error) {
	var models []domain.AccountGORM
	if err := r.db.WithContext(ctx).Where("parent_code = ?", parentCode).Order("code").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Account, len(models))
	for i := range models {
		out[i] = *gormAccountToDomain(&models[i])
	}
	return out, nil
}

func gormAccountToDomain(m *domain.AccountGORM) *domain.Account {
	return &domain.Account{
		Code: m.Code, Name: m.Name, Name2: m.Name2, Type: domain.AccountType(m.Type),
		ParentCode: m.ParentCode, IsActive: m.IsActive, IsForeign: m.IsForeign,
		DetailBy: domain.DetailBy(m.DetailBy), IsParent: m.IsParent, Status: domain.AccountStatus(m.Status),
		FreezeReason: m.FreezeReason, ArrearsDays: m.ArrearsDays, Note: m.Note,
	}
}

// ─── Journal Entry ───────────────────────────────────────────────────────────

func (r *PGJournalRepo) Create(ctx context.Context, e *domain.JournalEntry) error {
	m := domain.JournalEntryGORM{
		CompanyID: e.CompanyID, EntryNumber: e.EntryNumber, VoucherType: string(e.VoucherType),
		EntryDate: e.EntryDate, AccountingDate: e.AccountingDate, PeriodID: e.PeriodID,
		Description: e.Description, Status: string(e.Status), CurrencyCode: e.CurrencyCode,
		ExchangeRate: e.ExchangeRate, CreatedBy: e.CreatedBy,
	}
	lines := make([]domain.JournalLineGORM, len(e.Lines))
	for i, l := range e.Lines {
		lines[i] = domain.JournalLineGORM{
			LineNumber: l.LineNumber, AccountCode: l.AccountCode,
			DebitAmount: l.DebitAmount, CreditAmount: l.CreditAmount, Description: l.Description,
			CurrencyCode: l.CurrencyCode, ForeignAmount: l.ForeignAmount, ExchangeRate: l.ExchangeRate,
			ObjectID: l.ObjectID, ProjectID: l.ProjectID, ContractID: l.ContractID,
			CostItemID: l.CostItemID, DepartmentID: l.DepartmentID,
		}
	}
	m.Lines = lines
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	e.ID = fmt.Sprintf("JE-%d", m.ID)
	e.CreatedAt = m.CreatedAt
	return nil
}

func (r *PGJournalRepo) GetByID(ctx context.Context, id string) (*domain.JournalEntry, error) {
	var m domain.JournalEntryGORM
	if err := r.db.WithContext(ctx).Preload("Lines", func(db *gorm.DB) *gorm.DB {
		return db.Order("line_number")
	}).Where("id = ?", id).Or("entry_number = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormJournalEntryToDomain(&m), nil
}

func (r *PGJournalRepo) GetByPeriod(ctx context.Context, periodID string) ([]domain.JournalEntry, error) {
	var models []domain.JournalEntryGORM
	if err := r.db.WithContext(ctx).Where("period_id = ?", periodID).Order("entry_date").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormJournalEntriesToDomain(models), nil
}

func (r *PGJournalRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]domain.JournalEntry, error) {
	var models []domain.JournalEntryGORM
	if err := r.db.WithContext(ctx).Where("entry_date >= ? AND entry_date <= ?", from, to).Order("entry_date").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormJournalEntriesToDomain(models), nil
}

func (r *PGJournalRepo) GetByStatus(ctx context.Context, status domain.JournalEntryStatus) ([]domain.JournalEntry, error) {
	var models []domain.JournalEntryGORM
	if err := r.db.WithContext(ctx).Where("status = ?", string(status)).Order("entry_date").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormJournalEntriesToDomain(models), nil
}

func (r *PGJournalRepo) GetByVoucherType(ctx context.Context, vt domain.VoucherType) ([]domain.JournalEntry, error) {
	var models []domain.JournalEntryGORM
	if err := r.db.WithContext(ctx).Where("voucher_type = ?", string(vt)).Order("entry_date").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormJournalEntriesToDomain(models), nil
}

func (r *PGJournalRepo) GetAll(ctx context.Context) ([]domain.JournalEntry, error) {
	var models []domain.JournalEntryGORM
	if err := r.db.WithContext(ctx).Order("entry_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormJournalEntriesToDomain(models), nil
}

func (r *PGJournalRepo) UpdateStatus(ctx context.Context, id string, status domain.JournalEntryStatus) error {
	return r.db.WithContext(ctx).Model(&domain.JournalEntryGORM{}).Where("CAST(id AS TEXT) = ?", id).Or("entry_number = ?", id).Update("status", string(status)).Error
}

func (r *PGJournalRepo) Update(ctx context.Context, e *domain.JournalEntry) error {
	return r.db.WithContext(ctx).Model(&domain.JournalEntryGORM{}).Where("entry_number = ?", e.EntryNumber).Updates(map[string]interface{}{
		"voucher_type": string(e.VoucherType), "entry_date": e.EntryDate,
		"accounting_date": e.AccountingDate, "period_id": e.PeriodID,
		"description": e.Description, "status": string(e.Status),
		"currency_code": e.CurrencyCode, "exchange_rate": e.ExchangeRate,
	}).Error
}

func (r *PGJournalRepo) Approve(ctx context.Context, id, approvedBy string) error {
	return r.db.WithContext(ctx).Model(&domain.JournalEntryGORM{}).Where("CAST(id AS TEXT) = ?", id).Or("entry_number = ?", id).Updates(map[string]interface{}{
		"status": "APPROVED", "approved_by": approvedBy, "approved_at": time.Now(),
	}).Error
}

func (r *PGJournalRepo) Review(ctx context.Context, id, reviewedBy string) error {
	return r.db.WithContext(ctx).Model(&domain.JournalEntryGORM{}).Where("CAST(id AS TEXT) = ?", id).Or("entry_number = ?", id).Updates(map[string]interface{}{
		"status": "REVIEWING", "reviewed_by": reviewedBy,
	}).Error
}

func (r *PGJournalRepo) GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error) {
	var lines []domain.JournalLineGORM
	if err := r.db.WithContext(ctx).Where("entry_id = ?", entryID).Order("line_number").Find(&lines).Error; err != nil {
		return nil, err
	}
	out := make([]domain.JournalLine, len(lines))
	for i := range lines {
		out[i] = domain.JournalLine{
			ID: fmt.Sprint(lines[i].ID), EntryID: fmt.Sprint(lines[i].EntryID),
			LineNumber: lines[i].LineNumber, AccountCode: lines[i].AccountCode,
			DebitAmount: lines[i].DebitAmount, CreditAmount: lines[i].CreditAmount,
			Description: lines[i].Description, CurrencyCode: lines[i].CurrencyCode,
			ForeignAmount: lines[i].ForeignAmount, ExchangeRate: lines[i].ExchangeRate,
			ObjectID: lines[i].ObjectID, ProjectID: lines[i].ProjectID,
			ContractID: lines[i].ContractID, CostItemID: lines[i].CostItemID,
			DepartmentID: lines[i].DepartmentID,
		}
	}
	return out, nil
}

func (r *PGJournalRepo) GetBalance(ctx context.Context, accountCode, periodID string) (*domain.AccountBalance, error) {
	type row struct {
		PeriodDebit  float64
		PeriodCredit float64
	}
	var rw row
	err := r.db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(l.debit_amount),0) AS period_debit, COALESCE(SUM(l.credit_amount),0) AS period_credit FROM journal_lines l JOIN journal_entries e ON l.entry_id = CAST(e.id AS TEXT) WHERE l.account_code = ? AND e.period_id = ? AND e.status = 'POSTED'`, accountCode, periodID).Scan(&rw).Error
	if err != nil {
		return nil, err
	}
	return &domain.AccountBalance{
		AccountCode:  accountCode,
		PeriodID:     periodID,
		PeriodDebit:  rw.PeriodDebit,
		PeriodCredit: rw.PeriodCredit,
	}, nil
}

func (r *PGJournalRepo) GetTrialBalance(ctx context.Context, periodID string) ([]domain.AccountBalance, error) {
	var balances []domain.AccountBalance
	err := r.db.WithContext(ctx).Raw(`SELECT l.account_code AS account_code, ? AS period_id, a.type AS account_type, COALESCE(SUM(l.debit_amount),0) AS period_debit, COALESCE(SUM(l.credit_amount),0) AS period_credit FROM journal_lines l JOIN journal_entries e ON l.entry_id = CAST(e.id AS TEXT) JOIN accounts a ON l.account_code = a.code WHERE e.period_id = ? AND e.status = 'POSTED' GROUP BY l.account_code, a.type ORDER BY l.account_code`, periodID, periodID).Scan(&balances).Error
	if err != nil {
		return nil, err
	}
	return balances, nil
}

func (r *PGJournalRepo) GetFinancialStatement(ctx context.Context, periodID string, accountTypes []domain.AccountType) ([]domain.AccountBalance, error) {
	if len(accountTypes) == 0 {
		return nil, nil
	}
	types := make([]string, len(accountTypes))
	for i, t := range accountTypes {
		types[i] = string(t)
	}
	var balances []domain.AccountBalance
	err := r.db.WithContext(ctx).Raw(`SELECT l.account_code, ? AS period_id, a.type AS account_type, COALESCE(SUM(l.debit_amount),0) AS period_debit, COALESCE(SUM(l.credit_amount),0) AS period_credit FROM journal_lines l JOIN journal_entries e ON l.entry_id = CAST(e.id AS TEXT) JOIN accounts a ON l.account_code = a.code WHERE e.period_id = ? AND e.status = 'POSTED' AND a.type = ANY(?) GROUP BY l.account_code, a.type ORDER BY l.account_code`, periodID, periodID, types).Scan(&balances).Error
	return balances, err
}

func (r *PGJournalRepo) GetAccountUsage(ctx context.Context, accountCode string) (*domain.AccountUsage, error) {
	var u domain.AccountUsage
	u.AccountCode = accountCode
	err := r.db.WithContext(ctx).Raw(`SELECT COALESCE(COUNT(*),0) AS entry_count, COALESCE(SUM(l.debit_amount),0) AS total_debit, COALESCE(SUM(l.credit_amount),0) AS total_credit, COALESCE(MAX(e.entry_date::TEXT),'') AS last_used_date FROM journal_lines l JOIN journal_entries e ON l.entry_id = CAST(e.id AS TEXT) WHERE l.account_code = ? AND e.status = 'POSTED'`, accountCode).Scan(&u).Error
	return &u, err
}

func (r *PGJournalRepo) GetPostedEntriesByAccount(ctx context.Context, periodID, accountCode string) ([]domain.JournalEntry, error) {
	var models []domain.JournalEntryGORM
	err := r.db.WithContext(ctx).Raw(`SELECT DISTINCT e.* FROM journal_entries e JOIN journal_lines l ON l.entry_id = CAST(e.id AS TEXT) WHERE e.period_id = ? AND e.status = 'POSTED' AND l.account_code = ? ORDER BY e.entry_date`, periodID, accountCode).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	return gormJournalEntriesToDomain(models), nil
}

func gormJournalEntryToDomain(m *domain.JournalEntryGORM) *domain.JournalEntry {
	e := &domain.JournalEntry{
		ID: fmt.Sprint(m.ID), CompanyID: m.CompanyID, EntryNumber: m.EntryNumber,
		VoucherType: domain.VoucherType(m.VoucherType), EntryDate: m.EntryDate,
		AccountingDate: m.AccountingDate, PeriodID: m.PeriodID, Description: m.Description,
		Status: domain.JournalEntryStatus(m.Status), CurrencyCode: m.CurrencyCode,
		ExchangeRate: m.ExchangeRate, CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt, PostedAt: m.PostedAt, ApprovedAt: m.ApprovedAt,
		ReviewedBy: m.ReviewedBy, ApprovedBy: m.ApprovedBy,
	}
	if m.Lines != nil {
		e.Lines = make([]domain.JournalLine, len(m.Lines))
		for i, l := range m.Lines {
			e.Lines[i] = domain.JournalLine{
				ID: fmt.Sprint(l.ID), EntryID: fmt.Sprint(l.EntryID),
				LineNumber: l.LineNumber, AccountCode: l.AccountCode,
				DebitAmount: l.DebitAmount, CreditAmount: l.CreditAmount,
				Description: l.Description, CurrencyCode: l.CurrencyCode,
				ForeignAmount: l.ForeignAmount, ExchangeRate: l.ExchangeRate,
				ObjectID: l.ObjectID, ProjectID: l.ProjectID, ContractID: l.ContractID,
				CostItemID: l.CostItemID, DepartmentID: l.DepartmentID,
			}
		}
	}
	return e
}

func gormJournalEntriesToDomain(models []domain.JournalEntryGORM) []domain.JournalEntry {
	out := make([]domain.JournalEntry, len(models))
	for i := range models {
		out[i] = *gormJournalEntryToDomain(&models[i])
	}
	return out
}

// ─── Period ──────────────────────────────────────────────────────────────────

func (r *PGPeriodRepo) Create(ctx context.Context, p *domain.Period) error {
	return r.db.WithContext(ctx).Create(&domain.PeriodGORM{
		ID: p.ID, Year: p.Year, Month: p.Month,
		StartDate: p.StartDate, EndDate: p.EndDate, Status: string(p.Status),
	}).Error
}

func (r *PGPeriodRepo) GetByID(ctx context.Context, id string) (*domain.Period, error) {
	var m domain.PeriodGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrPeriodNotFound
	}
	return &domain.Period{ID: m.ID, Year: m.Year, Month: m.Month, StartDate: m.StartDate, EndDate: m.EndDate, Status: domain.PeriodStatus(m.Status)}, nil
}

func (r *PGPeriodRepo) GetByYearMonth(ctx context.Context, year, month int) (*domain.Period, error) {
	return r.GetByID(ctx, fmt.Sprintf("P-%d-%02d", year, month))
}

func (r *PGPeriodRepo) GetAll(ctx context.Context) ([]domain.Period, error) {
	var models []domain.PeriodGORM
	if err := r.db.WithContext(ctx).Order("year, month").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Period, len(models))
	for i := range models {
		out[i] = domain.Period{ID: models[i].ID, Year: models[i].Year, Month: models[i].Month, StartDate: models[i].StartDate, EndDate: models[i].EndDate, Status: domain.PeriodStatus(models[i].Status)}
	}
	return out, nil
}

func (r *PGPeriodRepo) UpdateStatus(ctx context.Context, id string, status domain.PeriodStatus) error {
	return r.db.WithContext(ctx).Model(&domain.PeriodGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGPeriodRepo) GetOpenPeriod(ctx context.Context) (*domain.Period, error) {
	var m domain.PeriodGORM
	if err := r.db.WithContext(ctx).Where("status = 'OPEN'").First(&m).Error; err != nil {
		return nil, domain.ErrPeriodNotFound
	}
	return &domain.Period{ID: m.ID, Year: m.Year, Month: m.Month, StartDate: m.StartDate, EndDate: m.EndDate, Status: domain.PeriodStatus(m.Status)}, nil
}

// ─── User ────────────────────────────────────────────────────────────────────

func (r *PGUserRepo) Create(ctx context.Context, u *domain.User) error {
	if u.ID == "" {
		u.ID = fmt.Sprintf("U-%d", time.Now().UnixNano())
	}
	return r.db.WithContext(ctx).Create(&domain.UserGORM{
		ID: u.ID, Username: u.Username, PasswordHash: u.PasswordHash,
		FullName: u.FullName, Email: u.Email, Role: string(u.Role), IsActive: u.IsActive,
	}).Error
}

func (r *PGUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var m domain.UserGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormUserToDomain(&m), nil
}

func (r *PGUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m domain.UserGORM
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		return nil, domain.ErrUserNotFound
	}
	return gormUserToDomain(&m), nil
}

func (r *PGUserRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	var models []domain.UserGORM
	if err := r.db.WithContext(ctx).Order("username").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.User, len(models))
	for i := range models {
		out[i] = *gormUserToDomain(&models[i])
	}
	return out, nil
}

func (r *PGUserRepo) Update(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Model(&domain.UserGORM{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
		"full_name": u.FullName, "email": u.Email, "role": string(u.Role), "is_active": u.IsActive,
	}).Error
}

func (r *PGUserRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.UserGORM{}).Error
}

func gormUserToDomain(m *domain.UserGORM) *domain.User {
	return &domain.User{
		ID: m.ID, Username: m.Username, PasswordHash: m.PasswordHash,
		FullName: m.FullName, Email: m.Email, Role: domain.UserRole(m.Role),
		IsActive: m.IsActive, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

// ─── Refresh Token ───────────────────────────────────────────────────────────

func (r *PGRefreshTokenRepo) Create(ctx context.Context, t *domain.RefreshToken) error {
	if t.ID == "" {
		t.ID = "RT" + time.Now().Format("20060102150405.000000000")
	}
	return r.db.WithContext(ctx).Create(&domain.RefreshTokenGORM{
		ID: t.ID, UserID: t.UserID, TokenHash: t.TokenHash,
		DeviceInfo: nullStrG(t.DeviceInfo), IPAddress: nullStrG(t.IPAddress),
		ExpiresAt: t.ExpiresAt, RevokedAt: t.RevokedAt,
	}).Error
}

func (r *PGRefreshTokenRepo) GetByID(ctx context.Context, id string) (*domain.RefreshToken, error) {
	var m domain.RefreshTokenGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return gormRefreshTokenToDomain(&m), nil
}

func (r *PGRefreshTokenRepo) GetByUserID(ctx context.Context, userID string) ([]domain.RefreshToken, error) {
	var models []domain.RefreshTokenGORM
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormRefreshTokensToDomain(models), nil
}

func (r *PGRefreshTokenRepo) GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var m domain.RefreshTokenGORM
	if err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&m).Error; err != nil {
		return nil, err
	}
	return gormRefreshTokenToDomain(&m), nil
}

func (r *PGRefreshTokenRepo) Revoke(ctx context.Context, id string) error {
	n := time.Now()
	return r.db.WithContext(ctx).Model(&domain.RefreshTokenGORM{}).Where("id = ?", id).Update("revoked_at", &n).Error
}

func (r *PGRefreshTokenRepo) RevokeAllByUserID(ctx context.Context, userID string) error {
	n := time.Now()
	return r.db.WithContext(ctx).Model(&domain.RefreshTokenGORM{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", &n).Error
}

func gormRefreshTokenToDomain(m *domain.RefreshTokenGORM) *domain.RefreshToken {
	d := &domain.RefreshToken{
		ID: m.ID, UserID: m.UserID, TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt, CreatedAt: m.CreatedAt, RevokedAt: m.RevokedAt,
	}
	if m.DeviceInfo != nil {
		d.DeviceInfo = *m.DeviceInfo
	}
	if m.IPAddress != nil {
		d.IPAddress = *m.IPAddress
	}
	return d
}

func gormRefreshTokensToDomain(models []domain.RefreshTokenGORM) []domain.RefreshToken {
	out := make([]domain.RefreshToken, len(models))
	for i := range models {
		out[i] = *gormRefreshTokenToDomain(&models[i])
	}
	return out
}

// ─── Password Reset Token ───────────────────────────────────────────────────

func (r *PGPasswordResetTokenRepo) Create(ctx context.Context, t *domain.PasswordResetToken) error {
	if t.ID == "" {
		t.ID = "PRT" + time.Now().Format("20060102150405.000000000")
	}
	return r.db.WithContext(ctx).Create(&domain.PasswordResetTokenGORM{
		ID: t.ID, UserID: t.UserID, TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt, UsedAt: t.UsedAt,
	}).Error
}

func (r *PGPasswordResetTokenRepo) GetByID(ctx context.Context, id string) (*domain.PasswordResetToken, error) {
	var m domain.PasswordResetTokenGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &domain.PasswordResetToken{
		ID: m.ID, UserID: m.UserID, TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt, CreatedAt: m.CreatedAt, UsedAt: m.UsedAt,
	}, nil
}

func (r *PGPasswordResetTokenRepo) MarkUsed(ctx context.Context, id string) error {
	n := time.Now()
	return r.db.WithContext(ctx).Model(&domain.PasswordResetTokenGORM{}).Where("id = ?", id).Update("used_at", &n).Error
}

// ─── Audit Log ───────────────────────────────────────────────────────────────

func (r *PGAuditLogRepo) Create(ctx context.Context, e *domain.AuditEntry) error {
	return r.db.WithContext(ctx).Create(&domain.AuditEntryGORM{
		UserID: nullStrG(e.UserID), Username: e.Username, IPAddress: nullStrG(e.IPAddress),
		Action: string(e.Action), EntityType: e.EntityType, EntityID: nullStrG(e.EntityID),
		OldValue: nullStrG(e.OldValue), NewValue: nullStrG(e.NewValue),
	}).Error
}

func (r *PGAuditLogRepo) GetByEntity(ctx context.Context, entityType, entityID string) ([]domain.AuditEntry, error) {
	var models []domain.AuditEntryGORM
	if err := r.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormAuditToDomain(models), nil
}

func (r *PGAuditLogRepo) GetByUser(ctx context.Context, userID string) ([]domain.AuditEntry, error) {
	var models []domain.AuditEntryGORM
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormAuditToDomain(models), nil
}

func (r *PGAuditLogRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]domain.AuditEntry, error) {
	var models []domain.AuditEntryGORM
	if err := r.db.WithContext(ctx).Where("created_at >= ? AND created_at <= ?", from, to).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormAuditToDomain(models), nil
}

func (r *PGAuditLogRepo) GetAll(ctx context.Context, limit int) ([]domain.AuditEntry, error) {
	var models []domain.AuditEntryGORM
	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}
	return gormAuditToDomain(models), nil
}

func gormAuditToDomain(models []domain.AuditEntryGORM) []domain.AuditEntry {
	out := make([]domain.AuditEntry, len(models))
	for i := range models {
		var ue string
		if models[i].UserID != nil {
			ue = *models[i].UserID
		}
		var ip string
		if models[i].IPAddress != nil {
			ip = *models[i].IPAddress
		}
		var eid string
		if models[i].EntityID != nil {
			eid = *models[i].EntityID
		}
		var ov string
		if models[i].OldValue != nil {
			ov = *models[i].OldValue
		}
		var nv string
		if models[i].NewValue != nil {
			nv = *models[i].NewValue
		}
		out[i] = domain.AuditEntry{
			ID: fmt.Sprint(models[i].ID), UserID: ue, Username: models[i].Username,
			IPAddress: ip, Action: domain.AuditAction(models[i].Action),
			EntityType: models[i].EntityType, EntityID: eid,
			OldValue: ov, NewValue: nv, CreatedAt: models[i].CreatedAt,
		}
	}
	return out
}

// ─── Exchange Rate ───────────────────────────────────────────────────────────

func (r *PGExchangeRateRepo) Create(ctx context.Context, rate *domain.ExchangeRate) error {
	return r.db.WithContext(ctx).Create(&domain.ExchangeRateGORM{
		CurrencyCode: rate.CurrencyCode, RateDate: rate.RateDate,
		BuyRate: float64Ptr(rate.BuyRate), SellRate: float64Ptr(rate.SellRate),
		AverageRate: rate.AverageRate, Source: nullStrG(rate.Source),
	}).Error
}

func (r *PGExchangeRateRepo) GetByCurrencyAndDate(ctx context.Context, currencyCode string, rateDate time.Time) (*domain.ExchangeRate, error) {
	var m domain.ExchangeRateGORM
	if err := r.db.WithContext(ctx).Where("currency_code = ? AND rate_date = ?", currencyCode, rateDate).First(&m).Error; err != nil {
		return nil, domain.ErrRateNotFound
	}
	return gormExchangeRateToDomain(&m), nil
}

func (r *PGExchangeRateRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]domain.ExchangeRate, error) {
	var models []domain.ExchangeRateGORM
	if err := r.db.WithContext(ctx).Where("rate_date >= ? AND rate_date <= ?", from, to).Order("rate_date, currency_code").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormExchangeRatesToDomain(models), nil
}

func (r *PGExchangeRateRepo) GetAll(ctx context.Context) ([]domain.ExchangeRate, error) {
	var models []domain.ExchangeRateGORM
	if err := r.db.WithContext(ctx).Order("rate_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	return gormExchangeRatesToDomain(models), nil
}

func (r *PGExchangeRateRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("CAST(id AS TEXT) = ?", id).Delete(&domain.ExchangeRateGORM{}).Error
}

func gormExchangeRateToDomain(m *domain.ExchangeRateGORM) *domain.ExchangeRate {
	return &domain.ExchangeRate{
		ID: fmt.Sprint(m.ID), CurrencyCode: m.CurrencyCode, RateDate: m.RateDate,
		BuyRate: safeFloat64(m.BuyRate), SellRate: safeFloat64(m.SellRate),
		AverageRate: m.AverageRate, Source: safeStr(m.Source),
		CreatedAt: m.CreatedAt,
	}
}

func gormExchangeRatesToDomain(models []domain.ExchangeRateGORM) []domain.ExchangeRate {
	out := make([]domain.ExchangeRate, len(models))
	for i := range models {
		out[i] = *gormExchangeRateToDomain(&models[i])
	}
	return out
}

// ─── Closing Template ────────────────────────────────────────────────────────

func (r *PGClosingTemplateRepo) Create(ctx context.Context, t *domain.ClosingTemplate) error {
	if t.ID == "" {
		t.ID = fmt.Sprintf("CT-%d", time.Now().UnixNano())
	}
	return r.db.WithContext(ctx).Create(&domain.ClosingTemplateGORM{
		Name: t.Name, Description: nullStrG(t.Description),
		SequenceOrder: t.SequenceOrder, IsActive: t.IsActive,
	}).Error
}

func (r *PGClosingTemplateRepo) GetByID(ctx context.Context, id string) (*domain.ClosingTemplate, error) {
	var m domain.ClosingTemplateGORM
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, domain.ErrTemplateNotFound
	}
	return &domain.ClosingTemplate{
		ID: fmt.Sprint(m.ID), Name: m.Name, Description: "",
		SequenceOrder: m.SequenceOrder, IsActive: m.IsActive, CreatedAt: m.CreatedAt,
	}, nil
}

func (r *PGClosingTemplateRepo) GetAll(ctx context.Context) ([]domain.ClosingTemplate, error) {
	var models []domain.ClosingTemplateGORM
	if err := r.db.WithContext(ctx).Order("sequence_order").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ClosingTemplate, len(models))
	for i := range models {
		out[i] = domain.ClosingTemplate{
			ID: fmt.Sprint(models[i].ID), Name: models[i].Name, Description: "",
			SequenceOrder: models[i].SequenceOrder, IsActive: models[i].IsActive,
			CreatedAt: models[i].CreatedAt,
		}
	}
	return out, nil
}

func (r *PGClosingTemplateRepo) Update(ctx context.Context, t *domain.ClosingTemplate) error {
	return r.db.WithContext(ctx).Model(&domain.ClosingTemplateGORM{}).Where("CAST(id AS TEXT) = ?", t.ID).Updates(map[string]interface{}{
		"name": t.Name, "description": t.Description,
		"sequence_order": t.SequenceOrder, "is_active": t.IsActive,
	}).Error
}

func (r *PGClosingTemplateRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.ClosingTemplateGORM{}, id).Error
}

// ─── Tax ──────────────────────────────────────────────────────────────

type PGTaxRepo struct{ db *gorm.DB }
func NewPGTaxRepo(db *gorm.DB) *PGTaxRepo { return &PGTaxRepo{db} }

func (r *PGTaxRepo) CreateDeclaration(ctx context.Context, d *domain.TaxDeclaration) error {
	if d.ID == "" {
		d.ID = fmt.Sprintf("TD-%d", time.Now().UnixNano())
	}
	m := domain.TaxDeclarationGORM{
		ID: d.ID, CompanyID: d.CompanyID, DeclarationType: string(d.DeclarationType),
		Status: string(d.Status), AdjustmentType: string(d.AdjustmentType),
		Version: d.Version, PreviousDeclID: nullStrG(d.PreviousDeclID),
		CreatedBy: d.CreatedBy,
	}
	m.TaxPeriod = domain.TaxPeriodGORM{
		PeriodYear: d.TaxPeriod.PeriodYear, PeriodNumber: d.TaxPeriod.PeriodNumber,
		PeriodType: string(d.TaxPeriod.PeriodType),
	}
	lines := make([]domain.TaxDeclarationLineGORM, len(d.Lines))
	for i, l := range d.Lines {
		lines[i] = domain.TaxDeclarationLineGORM{
			LineCode: l.LineCode, LineName: l.LineName, Amount: l.Amount,
			SourceType: string(l.SourceType), SortOrder: l.SortOrder,
		}
	}
	m.Lines = lines
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *PGTaxRepo) GetDeclarationByID(ctx context.Context, id string) (*domain.TaxDeclaration, error) {
	var m domain.TaxDeclarationGORM
	if err := r.db.WithContext(ctx).Preload("Lines", func(d *gorm.DB) *gorm.DB { return d.Order("sort_order") }).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrDeclarationNotFound
	}
	return taxDeclarationOut(&m), nil
}

func taxDeclarationOut(m *domain.TaxDeclarationGORM) *domain.TaxDeclaration {
	d := &domain.TaxDeclaration{
		ID: m.ID, CompanyID: m.CompanyID, DeclarationType: domain.DeclarationType(m.DeclarationType),
		Status: domain.DeclarationStatus(m.Status), AdjustmentType: domain.AdjustmentType(m.AdjustmentType),
		Version: m.Version, PreviousDeclID: safeStr(m.PreviousDeclID),
		SubmittedAt: safeTimePtrStr(m.SubmittedAt), SubmittedBy: safeStr(m.SubmittedBy),
		AcknowledgedAt: safeTimePtrStr(m.AcknowledgedAt), AcknowledgementRef: safeStr(m.AcknowledgementRef),
		GDTSubmissionID: safeStr(m.GDTSubmissionID),
		DeclarationXML: safeStr(m.DeclarationXML), GDTResponseXML: safeStr(m.GDTResponseXML),
		CreatedBy: m.CreatedBy, CreatedAt: safeTimeStr(m.CreatedAt), UpdatedAt: safeTimeStr(m.UpdatedAt),
	}
	d.TaxPeriod.PeriodYear = m.TaxPeriod.PeriodYear
	d.TaxPeriod.PeriodNumber = m.TaxPeriod.PeriodNumber
	d.TaxPeriod.PeriodType = domain.PeriodTypeV2(m.TaxPeriod.PeriodType)
	if m.Lines != nil {
		d.Lines = make([]domain.TaxDeclarationLine, len(m.Lines))
		for i, l := range m.Lines {
			d.Lines[i] = domain.TaxDeclarationLine{
				ID: l.ID, DeclarationID: l.DeclarationID, LineCode: l.LineCode,
				LineName: l.LineName, Amount: l.Amount,
				SourceType: domain.SourceType(l.SourceType), SortOrder: l.SortOrder,
			}
		}
	}
	return d
}

func (r *PGTaxRepo) GetDeclarations(ctx context.Context, filter domain.TaxDeclarationFilter) ([]domain.TaxDeclaration, error) {
	q := r.db.WithContext(ctx).Model(&domain.TaxDeclarationGORM{})
	if filter.CompanyID != "" { q = q.Where("company_id = ?", filter.CompanyID) }
	if filter.DeclarationType != "" { q = q.Where("declaration_type = ?", string(filter.DeclarationType)) }
	if filter.Status != "" { q = q.Where("status = ?", string(filter.Status)) }
	if filter.PeriodYear != 0 { q = q.Where("period_year = ?", filter.PeriodYear) }
	if filter.PeriodNumber != 0 { q = q.Where("period_number = ?", filter.PeriodNumber) }
	var models []domain.TaxDeclarationGORM
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.TaxDeclaration, len(models))
	for i := range models { out[i] = *taxDeclarationOut(&models[i]) }
	return out, nil
}

func (r *PGTaxRepo) UpdateDeclaration(ctx context.Context, d *domain.TaxDeclaration) error {
	return r.db.WithContext(ctx).Model(&domain.TaxDeclarationGORM{}).Where("id = ?", d.ID).Updates(map[string]interface{}{
		"status": string(d.Status), "declaration_xml": d.DeclarationXML,
		"adjustment_type": string(d.AdjustmentType),
		"submitted_at":    timePtr(parseDateTime(d.SubmittedAt)), "submitted_by": nullStrG(d.SubmittedBy),
		"acknowledged_at": timePtr(parseDateTime(d.AcknowledgedAt)), "acknowledgement_ref": nullStrG(d.AcknowledgementRef),
		"gdt_submission_id": nullStrG(d.GDTSubmissionID), "gdt_response_xml": d.GDTResponseXML,
	}).Error
}

func (r *PGTaxRepo) UpdateDeclarationStatus(ctx context.Context, id string, status domain.DeclarationStatus) error {
	return r.db.WithContext(ctx).Model(&domain.TaxDeclarationGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGTaxRepo) CreateRate(ctx context.Context, rate *domain.TaxRate) error {
	if rate.ID == "" { rate.ID = fmt.Sprintf("TR-%d", time.Now().UnixNano()) }
	return r.db.WithContext(ctx).Create(&domain.TaxRateGORM{
		ID: rate.ID, TaxType: string(rate.TaxType), RateCode: rate.RateCode,
		RateName: rate.RateName, RateType: string(rate.RateType),
		EffectiveFrom: parseDate(rate.EffectiveFrom), IsActive: rate.IsActive,
		IncentiveReducPct: float64Ptr(rate.IncentiveReducPct),
		LegalRef: nullStrG(rate.LegalRef),
	}).Error
}

func (r *PGTaxRepo) GetRateByID(ctx context.Context, id string) (*domain.TaxRate, error) {
	var m domain.TaxRateGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil { return nil, domain.ErrTaxRateNotFound }
	return &domain.TaxRate{
		ID: m.ID, TaxType: domain.TaxType(m.TaxType), RateCode: m.RateCode,
		RateName: m.RateName, RateType: domain.RateType(m.RateType),
		RateValue: safeFloat64(m.RateValue), EffectiveFrom: safeTimeStr(m.EffectiveFrom), IsActive: m.IsActive,
		IncentiveReducPct: safeFloat64(m.IncentiveReducPct),
		LegalRef: safeStr(m.LegalRef), CreatedAt: safeTimeStr(m.CreatedAt),
	}, nil
}

func (r *PGTaxRepo) GetRates(ctx context.Context, filter domain.TaxRateFilter) ([]domain.TaxRate, error) {
	q := r.db.WithContext(ctx).Model(&domain.TaxRateGORM{})
	if filter.TaxType != "" { q = q.Where("tax_type = ?", string(filter.TaxType)) }
	if filter.IsActive != nil { q = q.Where("is_active = ?", *filter.IsActive) }
	var models []domain.TaxRateGORM
	if err := q.Order("rate_code").Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.TaxRate, len(models))
	for i := range models {
		out[i] = domain.TaxRate{
			ID: models[i].ID, TaxType: domain.TaxType(models[i].TaxType),
			RateCode: models[i].RateCode, RateName: models[i].RateName,
			RateType: domain.RateType(models[i].RateType),
			RateValue: safeFloat64(models[i].RateValue),
			EffectiveFrom: safeTimeStr(models[i].EffectiveFrom), IsActive: models[i].IsActive,
			IncentiveReducPct: safeFloat64(models[i].IncentiveReducPct),
			CreatedAt: safeTimeStr(models[i].CreatedAt),
		}
	}
	return out, nil
}

func (r *PGTaxRepo) UpdateRate(ctx context.Context, rate *domain.TaxRate) error {
	return r.db.WithContext(ctx).Model(&domain.TaxRateGORM{}).Where("id = ?", rate.ID).Updates(map[string]interface{}{
		"rate_name": rate.RateName, "rate_type": string(rate.RateType),
		"is_active": rate.IsActive,
	}).Error
}

func (r *PGTaxRepo) CreatePayment(ctx context.Context, p *domain.TaxPayment) error {
	if p.ID == "" { p.ID = fmt.Sprintf("TP-%d", time.Now().UnixNano()) }
	return r.db.WithContext(ctx).Create(&domain.TaxPaymentGORM{
		ID: p.ID, CompanyID: p.CompanyID, DeclarationID: nullStrG(p.DeclarationID),
		TaxType: string(p.TaxType), PeriodYear: p.PeriodYear, PeriodNumber: p.PeriodNumber,
		DeclaredAmount: p.DeclaredAmount, PaidAmount: p.PaidAmount, DueDate: parseDate(p.DueDate),
		Status: string(p.Status), GLJournalID: nullStrG(p.GLJournalID),
	}).Error
}

func (r *PGTaxRepo) GetPaymentByID(ctx context.Context, id string) (*domain.TaxPayment, error) {
	var m domain.TaxPaymentGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil { return nil, domain.ErrTaxPaymentNotFound }
	return gormTaxPaymentToDomain(&m), nil
}

func (r *PGTaxRepo) GetPayments(ctx context.Context, filter domain.PaymentFilter) ([]domain.TaxPayment, error) {
	q := r.db.WithContext(ctx).Model(&domain.TaxPaymentGORM{})
	if filter.CompanyID != "" { q = q.Where("company_id = ?", filter.CompanyID) }
	if filter.TaxType != "" { q = q.Where("tax_type = ?", string(filter.TaxType)) }
	if filter.Status != "" { q = q.Where("status = ?", string(filter.Status)) }
	if filter.PeriodYear != 0 { q = q.Where("period_year = ?", filter.PeriodYear) }
	if filter.PeriodNumber != 0 { q = q.Where("period_number = ?", filter.PeriodNumber) }
	var models []domain.TaxPaymentGORM
	if err := q.Order("due_date").Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.TaxPayment, len(models))
	for i := range models { out[i] = *gormTaxPaymentToDomain(&models[i]) }
	return out, nil
}

func (r *PGTaxRepo) UpdatePayment(ctx context.Context, p *domain.TaxPayment) error {
	return r.db.WithContext(ctx).Model(&domain.TaxPaymentGORM{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
		"paid_amount": p.PaidAmount, "payment_ref": p.PaymentRef, "payment_date": nullStrG(p.PaymentDate),
		"payment_method": nullStrG(string(p.PaymentMethod)),
		"status": string(p.Status), "late_days": p.LateDays, "late_interest": p.LateInterest,
		"gl_journal_id": nullStrG(p.GLJournalID), "notes": p.Notes,
	}).Error
}

func gormTaxPaymentToDomain(m *domain.TaxPaymentGORM) *domain.TaxPayment {
	return &domain.TaxPayment{
		ID: m.ID, CompanyID: m.CompanyID, DeclarationID: safeStr(m.DeclarationID),
		TaxType: domain.TaxType(m.TaxType), PeriodYear: m.PeriodYear, PeriodNumber: m.PeriodNumber,
		DeclaredAmount: m.DeclaredAmount, PaidAmount: m.PaidAmount,
		PaymentDate: safeStr(m.PaymentDate), DueDate: safeTimeStr(m.DueDate),
		PaymentRef: safeStr(m.PaymentRef),
		PaymentMethod: domain.PaymentMethod(safeStr(m.PaymentMethod)),
		Status: domain.PaymentStatus(m.Status),
		LateDays: m.LateDays, LateInterest: m.LateInterest,
		GLJournalID: safeStr(m.GLJournalID),
		Notes: safeStr(m.Notes), CreatedAt: safeTimeStr(m.CreatedAt),
	}
}

func (r *PGTaxRepo) CreateEInvoice(ctx context.Context, inv *domain.EInvoice) error {
	if inv.ID == "" { inv.ID = fmt.Sprintf("INV-%d", time.Now().UnixNano()) }
	m := domain.EInvoiceGORM{
		ID: inv.ID, CompanyID: inv.CompanyID, Pattern: inv.Pattern, Serial: inv.Serial,
		InvoiceType: string(inv.InvoiceType), BuyerName: inv.BuyerName,
		BuyerTaxCode: nullStrG(inv.BuyerTaxCode), BuyerAddress: nullStrG(inv.BuyerAddress),
		BuyerEmail: nullStrG(inv.BuyerEmail), CurrencyCode: inv.CurrencyCode,
		ExchangeRate: inv.ExchangeRate, Subtotal: inv.Subtotal, VATAmount: inv.VATAmount,
		GrandTotal: inv.GrandTotal, IssueDate: parseDate(inv.IssueDate), Status: string(inv.Status),
	}
	lines := make([]domain.EInvoiceLineGORM, len(inv.Lines))
	for i, l := range inv.Lines {
		lines[i] = domain.EInvoiceLineGORM{LineNumber: l.LineNumber, Description: l.Description,
			Unit: nullStrG(l.Unit), Quantity: l.Quantity, UnitPrice: l.UnitPrice,
			LineTotal: l.LineTotal, VATRate: l.VATRate, VATAmount: l.VATAmount}
	}
	m.Lines = lines
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *PGTaxRepo) GetEInvoiceByID(ctx context.Context, id string) (*domain.EInvoice, error) {
	var m domain.EInvoiceGORM
	if err := r.db.WithContext(ctx).Preload("Lines", func(d *gorm.DB) *gorm.DB { return d.Order("line_number") }).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrInvoiceNotFound
	}
	return gormEInvoiceToDomain(&m), nil
}

func (r *PGTaxRepo) GetEInvoices(ctx context.Context, filter domain.EInvoiceFilter) ([]domain.EInvoice, error) {
	q := r.db.WithContext(ctx).Model(&domain.EInvoiceGORM{})
	if filter.CompanyID != "" { q = q.Where("company_id = ?", filter.CompanyID) }
	if filter.Status != "" { q = q.Where("status = ?", string(filter.Status)) }
	var models []domain.EInvoiceGORM
	if err := q.Preload("Lines").Order("issue_date DESC").Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.EInvoice, len(models))
	for i := range models { out[i] = *gormEInvoiceToDomain(&models[i]) }
	return out, nil
}

func gormEInvoiceToDomain(m *domain.EInvoiceGORM) *domain.EInvoice {
	e := &domain.EInvoice{
		ID: m.ID, CompanyID: m.CompanyID, Pattern: m.Pattern, Serial: m.Serial,
		InvoiceNumber: int64(m.InvoiceNumber), InvoiceType: domain.EInvoiceType(m.InvoiceType),
		BuyerName: m.BuyerName, CurrencyCode: m.CurrencyCode, ExchangeRate: m.ExchangeRate,
		Subtotal: m.Subtotal, VATAmount: m.VATAmount, GrandTotal: m.GrandTotal,
		IssueDate: safeTimeStr(m.IssueDate), Status: domain.EInvLifecycleStatus(m.Status),
		XMLBody: safeStr(m.XMLBody), SignedXML: safeStr(m.SignedXML), JournalEntryID: safeStr(m.JournalEntryID),
		CancelReason: safeStr(m.CancelReason), GDTTransactionID: safeStr(m.GDTResponse),
		CreatedAt: safeTimeStr(m.CreatedAt),
	}
	if m.BuyerTaxCode != nil { e.BuyerTaxCode = *m.BuyerTaxCode }
	if m.BuyerAddress != nil { e.BuyerAddress = *m.BuyerAddress }
	if m.BuyerEmail != nil { e.BuyerEmail = *m.BuyerEmail }
	if m.Lines != nil {
		e.Lines = make([]domain.EInvoiceLine, len(m.Lines))
		for i, l := range m.Lines {
			e.Lines[i] = domain.EInvoiceLine{
				LineNumber: l.LineNumber, Description: l.Description,
				Unit: "", Quantity: l.Quantity, UnitPrice: l.UnitPrice,
				LineTotal: l.LineTotal, VATRate: l.VATRate, VATAmount: l.VATAmount,
			}
			if l.Unit != nil { e.Lines[i].Unit = *l.Unit }
		}
	}
	return e
}

func (r *PGTaxRepo) UpdateEInvoice(ctx context.Context, inv *domain.EInvoice) error {
	return r.db.WithContext(ctx).Model(&domain.EInvoiceGORM{}).Where("id = ?", inv.ID).Updates(map[string]interface{}{
		"status": string(inv.Status), "xml_body": inv.XMLBody, "signed_xml": inv.SignedXML,
	}).Error
}

func (r *PGTaxRepo) UpdateEInvoiceStatus(ctx context.Context, id string, status domain.EInvLifecycleStatus) error {
	return r.db.WithContext(ctx).Model(&domain.EInvoiceGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGTaxRepo) CreateCalendarEntry(ctx context.Context, c *domain.TaxCalendar) error {
	m := domain.TaxCalendarGORM{
		ID: c.ID, CompanyID: c.CompanyID, TaxType: string(c.TaxType),
		PeriodType: string(c.PeriodType), PeriodYear: c.PeriodYear, PeriodNumber: c.PeriodNumber,
		StartDate: timePtr(parseDate(c.StartDate)), EndDate: timePtr(parseDate(c.EndDate)),
		DueDate: parseDate(c.DeclarationDue), DeclarationDue: parseDate(c.DeclarationDue),
		PaymentDue: timePtr(parseDate(c.PaymentDue)), Status: string(c.Status),
		DeclarationID: strPtr(c.DeclarationID),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *PGTaxRepo) GetCalendarEntryByID(ctx context.Context, id string) (*domain.TaxCalendar, error) {
	var m domain.TaxCalendarGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil { return nil, err }
	return &domain.TaxCalendar{
		ID: m.ID, CompanyID: m.CompanyID, TaxType: domain.TaxType(m.TaxType),
		PeriodType: domain.PeriodTypeV2(m.PeriodType), PeriodYear: m.PeriodYear, PeriodNumber: m.PeriodNumber,
		StartDate: safeTimePtrStr(m.StartDate), EndDate: safeTimePtrStr(m.EndDate),
		DeclarationDue: safeTimeStr(m.DeclarationDue), PaymentDue: safeTimePtrStr(m.PaymentDue),
		Status: domain.CalendarStatus(m.Status), DeclarationID: safeStr(m.DeclarationID),
		CreatedAt: safeTimeStr(m.CreatedAt),
	}, nil
}

func (r *PGTaxRepo) GetCalendarByPeriod(ctx context.Context, companyID string, periodYear, periodNumber int) ([]domain.TaxCalendar, error) {
	var models []domain.TaxCalendarGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND period_year = ? AND period_number = ?", companyID, periodYear, periodNumber).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.TaxCalendar, len(models))
	for i := range models {
		out[i] = domain.TaxCalendar{
			ID: models[i].ID, CompanyID: models[i].CompanyID, TaxType: domain.TaxType(models[i].TaxType),
			PeriodType: domain.PeriodTypeV2(models[i].PeriodType), PeriodYear: models[i].PeriodYear, PeriodNumber: models[i].PeriodNumber,
			StartDate: safeTimePtrStr(models[i].StartDate), EndDate: safeTimePtrStr(models[i].EndDate),
			DeclarationDue: safeTimeStr(models[i].DeclarationDue), PaymentDue: safeTimePtrStr(models[i].PaymentDue),
			Status: domain.CalendarStatus(models[i].Status), DeclarationID: safeStr(models[i].DeclarationID),
			CreatedAt: safeTimeStr(models[i].CreatedAt),
		}
	}
	return out, nil
}

func (r *PGTaxRepo) GetCalendarByCompany(ctx context.Context, companyID string) ([]domain.TaxCalendar, error) {
	var models []domain.TaxCalendarGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("period_year, period_number").Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.TaxCalendar, len(models))
	for i := range models {
		out[i] = domain.TaxCalendar{
			ID: models[i].ID, CompanyID: models[i].CompanyID, TaxType: domain.TaxType(models[i].TaxType),
			PeriodType: domain.PeriodTypeV2(models[i].PeriodType), PeriodYear: models[i].PeriodYear, PeriodNumber: models[i].PeriodNumber,
			StartDate: safeTimePtrStr(models[i].StartDate), EndDate: safeTimePtrStr(models[i].EndDate),
			DeclarationDue: safeTimeStr(models[i].DeclarationDue), PaymentDue: safeTimePtrStr(models[i].PaymentDue),
			Status: domain.CalendarStatus(models[i].Status), DeclarationID: safeStr(models[i].DeclarationID),
			CreatedAt: safeTimeStr(models[i].CreatedAt),
		}
	}
	return out, nil
}

func (r *PGTaxRepo) UpdateCalendarStatus(ctx context.Context, id string, status domain.CalendarStatus) error {
	return r.db.WithContext(ctx).Model(&domain.TaxCalendarGORM{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *PGTaxRepo) CreateAlert(ctx context.Context, a *domain.TaxAlert) error {
	m := domain.TaxAlertGORM{
		ID: a.ID, CompanyID: a.CompanyID, CalendarID: strPtr(a.CalendarID),
		AlertType: string(a.AlertType), Channel: string(a.Channel), Message: a.Message,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *PGTaxRepo) GetAlertByID(ctx context.Context, id string) (*domain.TaxAlert, error) {
	var m domain.TaxAlertGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil { return nil, err }
	return &domain.TaxAlert{
		ID: m.ID, CompanyID: m.CompanyID, CalendarID: safeStr(m.CalendarID),
		AlertType: domain.AlertType(m.AlertType), Channel: domain.AlertChannel(m.Channel),
		Message: m.Message, SentAt: m.CreatedAt.Format(time.RFC3339),
		AcknowledgedAt: safeTimePtrRFC3339(m.AcknowledgedAt), AcknowledgedBy: safeStr(m.AcknowledgedBy),
	}, nil
}

func (r *PGTaxRepo) GetAlerts(ctx context.Context, companyID string, limit int) ([]domain.TaxAlert, error) {
	var models []domain.TaxAlertGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Limit(limit).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.TaxAlert, len(models))
	for i := range models {
		out[i] = domain.TaxAlert{
			ID: models[i].ID, CompanyID: models[i].CompanyID, CalendarID: safeStr(models[i].CalendarID),
			AlertType: domain.AlertType(models[i].AlertType), Channel: domain.AlertChannel(models[i].Channel),
			Message: models[i].Message, SentAt: models[i].CreatedAt.Format(time.RFC3339),
			AcknowledgedAt: safeTimePtrRFC3339(models[i].AcknowledgedAt), AcknowledgedBy: safeStr(models[i].AcknowledgedBy),
		}
	}
	return out, nil
}

func (r *PGTaxRepo) CreateAuditCase(ctx context.Context, a *domain.TaxAuditCase) error {
	m := domain.TaxAuditCaseGORM{
		ID: a.ID, CompanyID: a.CompanyID, CaseNumber: a.AuditDecNumber,
		AuditType: "GENERAL", Status: string(a.Status),
		OpenDate: parseDate(a.AuditPeriodStart),
		AuditorName: nullStrG(a.AuditorName), Notes: nullStrG(a.Findings),
		AuditPeriodStart: timePtr(parseDate(a.AuditPeriodStart)),
		AuditPeriodEnd: timePtr(parseDate(a.AuditPeriodEnd)),
		AuditDecisionNumber: nullStrG(a.AuditDecNumber),
		AuditorContact: nullStrG(a.AuditorContact),
		Findings: nullStrG(a.Findings),
		PenaltyAmount: a.PenaltyAmount,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *PGTaxRepo) GetAuditCaseByID(ctx context.Context, id string) (*domain.TaxAuditCase, error) {
	var m domain.TaxAuditCaseGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil { return nil, err }
	return &domain.TaxAuditCase{
		ID: m.ID, CompanyID: m.CompanyID, AuditDecNumber: safeStr(m.AuditDecisionNumber),
		AuditPeriodStart: safeTimePtrStr(m.AuditPeriodStart), AuditPeriodEnd: safeTimePtrStr(m.AuditPeriodEnd),
		AuditorName: safeStr(m.AuditorName), AuditorContact: safeStr(m.AuditorContact),
		Status: domain.AuditCaseStatus(m.Status),
		Findings: safeStr(m.Findings), PenaltyAmount: m.PenaltyAmount,
		CreatedAt: safeTimeStr(m.CreatedAt), ClosedAt: safeTimePtrStr(m.CloseDate),
	}, nil
}

func (r *PGTaxRepo) GetAuditCases(ctx context.Context, companyID string) ([]domain.TaxAuditCase, error) {
	var models []domain.TaxAuditCaseGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.TaxAuditCase, len(models))
	for i := range models {
		out[i] = domain.TaxAuditCase{
			ID: models[i].ID, CompanyID: models[i].CompanyID, AuditDecNumber: safeStr(models[i].AuditDecisionNumber),
			AuditPeriodStart: safeTimePtrStr(models[i].AuditPeriodStart), AuditPeriodEnd: safeTimePtrStr(models[i].AuditPeriodEnd),
			AuditorName: safeStr(models[i].AuditorName), AuditorContact: safeStr(models[i].AuditorContact),
			Status: domain.AuditCaseStatus(models[i].Status),
			Findings: safeStr(models[i].Findings), PenaltyAmount: models[i].PenaltyAmount,
			CreatedAt: safeTimeStr(models[i].CreatedAt), ClosedAt: safeTimePtrStr(models[i].CloseDate),
		}
	}
	return out, nil
}

func (r *PGTaxRepo) UpdateAuditCase(ctx context.Context, a *domain.TaxAuditCase) error {
	return r.db.WithContext(ctx).Model(&domain.TaxAuditCaseGORM{}).Where("id = ?", a.ID).Updates(map[string]interface{}{
		"status": string(a.Status), "close_date": stringToTimePtr(a.ClosedAt),
		"findings": a.Findings, "penalty_amount": a.PenaltyAmount,
		"auditor_contact": nullStrG(a.AuditorContact),
	}).Error
}

func (r *PGTaxRepo) CreateCITLoss(ctx context.Context, loss *domain.CITLossCarryForward) error {
	g := domain.CITLossCarryForwardGORM{
		ID: loss.ID, CompanyID: loss.CompanyID, LossYear: loss.LossYear,
		LossAmount: loss.LossAmount, UsedAmount: loss.UsedAmount, ExpiryYear: loss.ExpiryYear,
	}
	return r.db.WithContext(ctx).Create(&g).Error
}

func (r *PGTaxRepo) GetActiveCITLosses(ctx context.Context, companyID string, beforeYear int) ([]domain.CITLossCarryForward, error) {
	var rows []domain.CITLossCarryForwardGORM
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND loss_year < ? AND expiry_year >= ? AND used_amount < loss_amount",
			companyID, beforeYear, beforeYear).
		Order("loss_year ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	var result []domain.CITLossCarryForward
	for _, row := range rows {
		result = append(result, domain.CITLossCarryForward{
			ID: row.ID, CompanyID: row.CompanyID, LossYear: row.LossYear,
			LossAmount: row.LossAmount, UsedAmount: row.UsedAmount, ExpiryYear: row.ExpiryYear,
		})
	}
	return result, nil
}

func (r *PGTaxRepo) UpdateCITLoss(ctx context.Context, loss *domain.CITLossCarryForward) error {
	return r.db.WithContext(ctx).Model(&domain.CITLossCarryForwardGORM{}).
		Where("id = ?", loss.ID).
		Update("used_amount", loss.UsedAmount).Error
}

// ─── Approval ─────────────────────────────────────────────────────────

type PGApprovalRepo struct{ db *gorm.DB }
func NewPGApprovalRepo(db *gorm.DB) *PGApprovalRepo { return &PGApprovalRepo{db} }

func (r *PGApprovalRepo) Create(ctx context.Context, req *domain.ApprovalRequest) error {
	return r.db.WithContext(ctx).Create(&domain.ApprovalRequestGORM{
		EntityType: req.EntityType, EntityID: req.EntityID, RequestedBy: req.RequestedBy,
		Status: string(req.Status),
	}).Error
}

func (r *PGApprovalRepo) GetByID(ctx context.Context, id string) (*domain.ApprovalRequest, error) {
	var m domain.ApprovalRequestGORM
	if err := r.db.WithContext(ctx).Where("CAST(id AS TEXT) = ?", id).First(&m).Error; err != nil { return nil, err }
	return &domain.ApprovalRequest{ID: fmt.Sprint(m.ID), EntityType: m.EntityType, EntityID: m.EntityID, RequestedBy: m.RequestedBy, Status: domain.ApprovalStatus(m.Status), CreatedAt: m.RequestedAt}, nil
}

func (r *PGApprovalRepo) GetByStatus(ctx context.Context, status domain.ApprovalStatus) ([]domain.ApprovalRequest, error) {
	var models []domain.ApprovalRequestGORM
	if err := r.db.WithContext(ctx).Where("status = ?", string(status)).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.ApprovalRequest, len(models))
	for i := range models { out[i] = domain.ApprovalRequest{ID: fmt.Sprint(models[i].ID), EntityType: models[i].EntityType, EntityID: models[i].EntityID, RequestedBy: models[i].RequestedBy, Status: domain.ApprovalStatus(models[i].Status), CreatedAt: models[i].RequestedAt} }
	return out, nil
}

func (r *PGApprovalRepo) GetByEntity(ctx context.Context, entityType, entityID string) ([]domain.ApprovalRequest, error) {
	var models []domain.ApprovalRequestGORM
	if err := r.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityID).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.ApprovalRequest, len(models))
	for i := range models { out[i] = domain.ApprovalRequest{ID: fmt.Sprint(models[i].ID), EntityType: models[i].EntityType, EntityID: models[i].EntityID, RequestedBy: models[i].RequestedBy, Status: domain.ApprovalStatus(models[i].Status), CreatedAt: models[i].RequestedAt} }
	return out, nil
}

func (r *PGApprovalRepo) UpdateStatus(ctx context.Context, id string, status domain.ApprovalStatus, reviewedBy, reviewNote string) error {
	n := time.Now()
	return r.db.WithContext(ctx).Model(&domain.ApprovalRequestGORM{}).Where("CAST(id AS TEXT) = ?", id).Updates(map[string]interface{}{
		"status": string(status), "reviewed_by": reviewedBy, "review_note": reviewNote, "reviewed_at": &n,
	}).Error
}

func (r *PGApprovalRepo) GetAll(ctx context.Context) ([]domain.ApprovalRequest, error) {
	var models []domain.ApprovalRequestGORM
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.ApprovalRequest, len(models))
	for i := range models { out[i] = domain.ApprovalRequest{ID: fmt.Sprint(models[i].ID), EntityType: models[i].EntityType, EntityID: models[i].EntityID, RequestedBy: models[i].RequestedBy, Status: domain.ApprovalStatus(models[i].Status), CreatedAt: models[i].RequestedAt} }
	return out, nil
}

// ─── Account Version ──────────────────────────────────────────────────────

type PGAccountVersionRepo struct{ db *gorm.DB }
func NewPGAccountVersionRepo(db *gorm.DB) *PGAccountVersionRepo { return &PGAccountVersionRepo{db} }

func (r *PGAccountVersionRepo) Create(ctx context.Context, v *domain.AccountVersion) error {
	return r.db.WithContext(ctx).Create(&domain.AccountVersionGORM{
		VersionNumber: v.VersionNumber, Name: v.Snapshot, Description: strPtr(v.ChangeSummary),
		CreatedBy: v.CreatedBy,
	}).Error
}

func (r *PGAccountVersionRepo) GetByVersionNumber(ctx context.Context, versionNumber string) (*domain.AccountVersion, error) {
	var m domain.AccountVersionGORM
	if err := r.db.WithContext(ctx).Where("version_number = ?", versionNumber).First(&m).Error; err != nil { return nil, err }
	return &domain.AccountVersion{ID: fmt.Sprint(m.ID), VersionNumber: m.VersionNumber, Snapshot: m.Name, ChangeSummary: safeStr(m.Description), CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt}, nil
}

func (r *PGAccountVersionRepo) GetLatest(ctx context.Context) (*domain.AccountVersion, error) {
	var m domain.AccountVersionGORM
	if err := r.db.WithContext(ctx).Order("effective_from DESC").First(&m).Error; err != nil { return nil, err }
	return &domain.AccountVersion{ID: fmt.Sprint(m.ID), VersionNumber: m.VersionNumber, Snapshot: m.Name, ChangeSummary: safeStr(m.Description), CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt}, nil
}

func (r *PGAccountVersionRepo) GetAll(ctx context.Context) ([]domain.AccountVersion, error) {
	var models []domain.AccountVersionGORM
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.AccountVersion, len(models))
	for i := range models { out[i] = domain.AccountVersion{ID: fmt.Sprint(models[i].ID), VersionNumber: models[i].VersionNumber, Snapshot: models[i].Name, ChangeSummary: safeStr(models[i].Description), CreatedBy: models[i].CreatedBy, CreatedAt: models[i].CreatedAt} }
	return out, nil
}

// ─── Account Mapping ─────────────────────────────────────────────────────

type PGAccountMappingRepo struct{ db *gorm.DB }
func NewPGAccountMappingRepo(db *gorm.DB) *PGAccountMappingRepo { return &PGAccountMappingRepo{db} }

func (r *PGAccountMappingRepo) Create(ctx context.Context, m *domain.AccountMapping) error {
	return r.db.WithContext(ctx).Create(&domain.AccountMappingGORM{
		SourceRegime: m.SourceRegime, TargetRegime: m.TargetRegime,
		OldAccountCode: m.OldCode, NewAccountCode: m.NewCode,
		MappingType: m.MappingType,
	}).Error
}

func (r *PGAccountMappingRepo) GetByOldCode(ctx context.Context, sourceRegime, oldCode string) (*domain.AccountMapping, error) {
	var m domain.AccountMappingGORM
	if err := r.db.WithContext(ctx).Where("source_regime = ? AND old_account_code = ?", sourceRegime, oldCode).First(&m).Error; err != nil { return nil, err }
	return &domain.AccountMapping{SourceRegime: m.SourceRegime, TargetRegime: m.TargetRegime, OldCode: m.OldAccountCode, NewCode: m.NewAccountCode, MappingType: m.MappingType}, nil
}

func (r *PGAccountMappingRepo) GetByRegime(ctx context.Context, sourceRegime, targetRegime string) ([]domain.AccountMapping, error) {
	var models []domain.AccountMappingGORM
	if err := r.db.WithContext(ctx).Where("source_regime = ? AND target_regime = ?", sourceRegime, targetRegime).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.AccountMapping, len(models))
	for i := range models { out[i] = domain.AccountMapping{SourceRegime: models[i].SourceRegime, TargetRegime: models[i].TargetRegime, OldCode: models[i].OldAccountCode, NewCode: models[i].NewAccountCode, MappingType: models[i].MappingType} }
	return out, nil
}

func (r *PGAccountMappingRepo) GetAll(ctx context.Context) ([]domain.AccountMapping, error) {
	var models []domain.AccountMappingGORM
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.AccountMapping, len(models))
	for i := range models { out[i] = domain.AccountMapping{SourceRegime: models[i].SourceRegime, TargetRegime: models[i].TargetRegime, OldCode: models[i].OldAccountCode, NewCode: models[i].NewAccountCode, MappingType: models[i].MappingType} }
	return out, nil
}

// ─── Account Analysis ─────────────────────────────────────────────────────

type PGAccountAnalysisRepo struct{ db *gorm.DB }
func NewPGAccountAnalysisRepo(db *gorm.DB) *PGAccountAnalysisRepo { return &PGAccountAnalysisRepo{db} }

func (r *PGAccountAnalysisRepo) Create(ctx context.Context, a *domain.AccountAnalysis) error {
	return r.db.WithContext(ctx).Create(&domain.AccountAnalysisGORM{
		AccountCode: a.AccountCode,
	}).Error
}

func (r *PGAccountAnalysisRepo) GetByAccount(ctx context.Context, accountCode string) (*domain.AccountAnalysis, error) {
	var m domain.AccountAnalysisGORM
	if err := r.db.WithContext(ctx).Where("account_code = ?", accountCode).First(&m).Error; err != nil { return nil, err }
	return &domain.AccountAnalysis{AccountCode: m.AccountCode}, nil
}

func (r *PGAccountAnalysisRepo) Update(ctx context.Context, a *domain.AccountAnalysis) error {
	return r.db.WithContext(ctx).Model(&domain.AccountAnalysisGORM{}).Where("account_code = ?", a.AccountCode).Updates(map[string]interface{}{
		"account_code": a.AccountCode,
	}).Error
}

// ─── IFRS Mapping ─────────────────────────────────────────────────────────

type PGIFRSMappingRepo struct{ db *gorm.DB }
func NewPGIFRSMappingRepo(db *gorm.DB) *PGIFRSMappingRepo { return &PGIFRSMappingRepo{db} }

func (r *PGIFRSMappingRepo) Create(ctx context.Context, m *domain.IFRSMapping) error {
	return r.db.WithContext(ctx).Create(&domain.IFRSMappingGORM{
		VASCode: m.VASCode, IFRS: m.IFRSCode,
	}).Error
}

func (r *PGIFRSMappingRepo) GetByVASCode(ctx context.Context, vasCode string) (*domain.IFRSMapping, error) {
	var m domain.IFRSMappingGORM
	if err := r.db.WithContext(ctx).Where("vas_code = ?", vasCode).First(&m).Error; err != nil { return nil, err }
	return &domain.IFRSMapping{VASCode: m.VASCode, IFRSCode: m.IFRS}, nil
}

func (r *PGIFRSMappingRepo) GetAll(ctx context.Context) ([]domain.IFRSMapping, error) {
	var models []domain.IFRSMappingGORM
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil { return nil, err }
	out := make([]domain.IFRSMapping, len(models))
	for i := range models { out[i] = domain.IFRSMapping{VASCode: models[i].VASCode, IFRSCode: models[i].IFRS} }
	return out, nil
}

func (r *PGIFRSMappingRepo) Update(ctx context.Context, m *domain.IFRSMapping) error {
	return r.db.WithContext(ctx).Model(&domain.IFRSMappingGORM{}).Where("vas_code = ?", m.VASCode).Updates(map[string]interface{}{
		"ifrs_code": m.IFRSCode,
	}).Error
}
