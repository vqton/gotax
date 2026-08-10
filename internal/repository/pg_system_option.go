package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

// ─── SystemOption ──────────────────────────────────────────────────

type pgSystemOption struct {
	ID        string    `gorm:"column:id;primaryKey"`
	CompanyID string    `gorm:"column:company_id"`
	Category  string    `gorm:"column:category"`
	Key       string    `gorm:"column:key"`
	Value     string    `gorm:"column:value"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (pgSystemOption) TableName() string { return "system_options" }

type PGSystemOptionRepo struct{ db *gorm.DB }

func NewPGSystemOptionRepo(db *gorm.DB) *PGSystemOptionRepo {
	return &PGSystemOptionRepo{db: db}
}

func toPGSystemOption(opt *domain.SystemOption) *pgSystemOption {
	var created, updated time.Time
	if opt.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", opt.CreatedAt)
	} else {
		created = time.Now()
	}
	if opt.UpdatedAt != "" {
		updated, _ = time.Parse("2006-01-02T15:04:05Z", opt.UpdatedAt)
	} else {
		updated = time.Now()
	}
	return &pgSystemOption{
		ID:        opt.ID,
		CompanyID: opt.CompanyID,
		Category:  opt.Category,
		Key:       opt.Key,
		Value:     opt.Value,
		CreatedAt: created,
		UpdatedAt: updated,
	}
}

func toDomainSystemOption(m *pgSystemOption) *domain.SystemOption {
	return &domain.SystemOption{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Category:  m.Category,
		Key:       m.Key,
		Value:     m.Value,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (r *PGSystemOptionRepo) Get(ctx context.Context, companyID, category, key string) (*domain.SystemOption, error) {
	var m pgSystemOption
	err := r.db.WithContext(ctx).Where("company_id = ? AND category = ? AND key = ?", companyID, category, key).First(&m).Error
	if err != nil {
		return nil, domain.ErrSystemOptionNotFound
	}
	return toDomainSystemOption(&m), nil
}

func (r *PGSystemOptionRepo) GetByCategory(ctx context.Context, companyID, category string) ([]domain.SystemOption, error) {
	var rows []pgSystemOption
	if err := r.db.WithContext(ctx).Where("company_id = ? AND category = ?", companyID, category).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.SystemOption, len(rows))
	for i, row := range rows {
		out[i] = *toDomainSystemOption(&row)
	}
	return out, nil
}

func (r *PGSystemOptionRepo) GetAll(ctx context.Context, companyID string) ([]domain.SystemOption, error) {
	var rows []pgSystemOption
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.SystemOption, len(rows))
	for i, row := range rows {
		out[i] = *toDomainSystemOption(&row)
	}
	return out, nil
}

func (r *PGSystemOptionRepo) Upsert(ctx context.Context, opt *domain.SystemOption) error {
	m := toPGSystemOption(opt)
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *PGSystemOptionRepo) Delete(ctx context.Context, companyID, category, key string) error {
	result := r.db.WithContext(ctx).Where("company_id = ? AND category = ? AND key = ?", companyID, category, key).Delete(&pgSystemOption{})
	if result.RowsAffected == 0 {
		return domain.ErrSystemOptionNotFound
	}
	return result.Error
}

// ─── NumberingRule ─────────────────────────────────────────────────

type pgNumberingRule struct {
	ID           string    `gorm:"column:id;primaryKey"`
	CompanyID    string    `gorm:"column:company_id"`
	VoucherType  string    `gorm:"column:voucher_type"`
	Prefix       string    `gorm:"column:prefix"`
	Suffix       string    `gorm:"column:suffix"`
	NumberLength int       `gorm:"column:number_length"`
	Scope        string    `gorm:"column:scope"`
	ResetRule    string    `gorm:"column:reset_rule"`
	CurrentNum   int       `gorm:"column:current_num"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (pgNumberingRule) TableName() string { return "numbering_rules" }

type PGNumberingRuleRepo struct{ db *gorm.DB }

func NewPGNumberingRuleRepo(db *gorm.DB) *PGNumberingRuleRepo {
	return &PGNumberingRuleRepo{db: db}
}

func toPGNumberingRule(rule *domain.NumberingRule) *pgNumberingRule {
	var created, updated time.Time
	if rule.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", rule.CreatedAt)
	} else {
		created = time.Now()
	}
	if rule.UpdatedAt != "" {
		updated, _ = time.Parse("2006-01-02T15:04:05Z", rule.UpdatedAt)
	} else {
		updated = time.Now()
	}
	return &pgNumberingRule{
		ID:           rule.ID,
		CompanyID:    rule.CompanyID,
		VoucherType:  rule.VoucherType,
		Prefix:       rule.Prefix,
		Suffix:       rule.Suffix,
		NumberLength: rule.NumberLength,
		Scope:        rule.Scope,
		ResetRule:    rule.ResetRule,
		CurrentNum:   rule.CurrentNum,
		CreatedAt:    created,
		UpdatedAt:    updated,
	}
}

func toDomainNumberingRule(m *pgNumberingRule) *domain.NumberingRule {
	return &domain.NumberingRule{
		ID:           m.ID,
		CompanyID:    m.CompanyID,
		VoucherType:  m.VoucherType,
		Prefix:       m.Prefix,
		Suffix:       m.Suffix,
		NumberLength: m.NumberLength,
		Scope:        m.Scope,
		ResetRule:    m.ResetRule,
		CurrentNum:   m.CurrentNum,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (r *PGNumberingRuleRepo) Create(ctx context.Context, rule *domain.NumberingRule) error {
	m := toPGNumberingRule(rule)
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if m.NumberLength == 0 {
		m.NumberLength = 5
	}
	if m.Scope == "" {
		m.Scope = "company"
	}
	if m.ResetRule == "" {
		m.ResetRule = "never"
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGNumberingRuleRepo) GetByID(ctx context.Context, id string) (*domain.NumberingRule, error) {
	var m pgNumberingRule
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, domain.ErrNumberingRuleNotFound
	}
	return toDomainNumberingRule(&m), nil
}

func (r *PGNumberingRuleRepo) GetByVoucherType(ctx context.Context, companyID, voucherType string) (*domain.NumberingRule, error) {
	var m pgNumberingRule
	err := r.db.WithContext(ctx).Where("company_id = ? AND voucher_type = ?", companyID, voucherType).First(&m).Error
	if err != nil {
		return nil, domain.ErrNumberingRuleNotFound
	}
	return toDomainNumberingRule(&m), nil
}

func (r *PGNumberingRuleRepo) List(ctx context.Context, companyID string) ([]domain.NumberingRule, error) {
	var rows []pgNumberingRule
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.NumberingRule, len(rows))
	for i, row := range rows {
		out[i] = *toDomainNumberingRule(&row)
	}
	return out, nil
}

func (r *PGNumberingRuleRepo) Update(ctx context.Context, rule *domain.NumberingRule) error {
	m := toPGNumberingRule(rule)
	m.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *PGNumberingRuleRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&pgNumberingRule{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return domain.ErrNumberingRuleNotFound
	}
	return result.Error
}

func (r *PGNumberingRuleRepo) IncrementAndGet(ctx context.Context, companyID, voucherType string) (int, error) {
	var m pgNumberingRule
	err := r.db.WithContext(ctx).Where("company_id = ? AND voucher_type = ?", companyID, voucherType).First(&m).Error
	if err != nil {
		return 0, domain.ErrNumberingRuleNotFound
	}
	m.CurrentNum++
	m.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return 0, err
	}
	return m.CurrentNum, nil
}

// ─── Backup ────────────────────────────────────────────────────────

type pgBackupRecord struct {
	ID         string    `gorm:"column:id;primaryKey"`
	CompanyID  string    `gorm:"column:company_id"`
	Filename   string    `gorm:"column:filename"`
	FileSize   int64     `gorm:"column:file_size"`
	BackupType string    `gorm:"column:backup_type"`
	Status     string    `gorm:"column:status"`
	CreatedBy  *string   `gorm:"column:created_by"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (pgBackupRecord) TableName() string { return "backup_records" }

type PGBackupRepo struct{ db *gorm.DB }

func NewPGBackupRepo(db *gorm.DB) *PGBackupRepo {
	return &PGBackupRepo{db: db}
}

func toPGBackup(b *domain.BackupRecord) *pgBackupRecord {
	var created time.Time
	if b.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", b.CreatedAt)
	} else {
		created = time.Now()
	}
	m := &pgBackupRecord{
		ID:         b.ID,
		CompanyID:  b.CompanyID,
		Filename:   b.Filename,
		FileSize:   b.FileSize,
		BackupType: b.BackupType,
		Status:     b.Status,
		CreatedAt:  created,
	}
	if b.CreatedBy != "" {
		m.CreatedBy = &b.CreatedBy
	}
	return m
}

func toDomainBackup(m *pgBackupRecord) *domain.BackupRecord {
	b := &domain.BackupRecord{
		ID:         m.ID,
		CompanyID:  m.CompanyID,
		Filename:   m.Filename,
		FileSize:   m.FileSize,
		BackupType: m.BackupType,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if m.CreatedBy != nil {
		b.CreatedBy = *m.CreatedBy
	}
	return b
}

func (r *PGBackupRepo) Create(ctx context.Context, b *domain.BackupRecord) error {
	m := toPGBackup(b)
	m.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGBackupRepo) GetByID(ctx context.Context, id string) (*domain.BackupRecord, error) {
	var m pgBackupRecord
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, domain.ErrBackupNotFound
	}
	return toDomainBackup(&m), nil
}

func (r *PGBackupRepo) List(ctx context.Context, companyID string) ([]domain.BackupRecord, error) {
	var rows []pgBackupRecord
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.BackupRecord, len(rows))
	for i, row := range rows {
		out[i] = *toDomainBackup(&row)
	}
	return out, nil
}

func (r *PGBackupRepo) UpdateStatus(ctx context.Context, id, status string) error {
	result := r.db.WithContext(ctx).Model(&pgBackupRecord{}).Where("id = ?", id).Update("status", status)
	if result.RowsAffected == 0 {
		return domain.ErrBackupNotFound
	}
	return result.Error
}
