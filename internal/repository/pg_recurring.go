package repository

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type pgRecurringEntry struct {
	ID           string `gorm:"column:id;primaryKey"`
	CompanyID    string `gorm:"column:company_id"`
	TemplateName string `gorm:"column:template_name"`
	Description  string `gorm:"column:description"`
	Frequency    string `gorm:"column:frequency"`
	DayOfMonth   int    `gorm:"column:day_of_month"`
	IsActive     bool   `gorm:"column:is_active"`
	NextRunDate  string `gorm:"column:next_run_date"`
	LinesJSON    string `gorm:"column:lines_json"`
	CreatedBy    string `gorm:"column:created_by"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (pgRecurringEntry) TableName() string { return "recurring_entries" }

type PGRecurringEntryRepo struct {
	db *gorm.DB
}

func NewPGRecurringEntryRepo(db *gorm.DB) *PGRecurringEntryRepo {
	return &PGRecurringEntryRepo{db: db}
}

func (r *PGRecurringEntryRepo) Create(ctx context.Context, entry *domain.RecurringEntry) error {
	m := toPGRecurringEntry(entry)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *PGRecurringEntryRepo) GetByID(ctx context.Context, id string) (*domain.RecurringEntry, error) {
	var m pgRecurringEntry
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrRecurringNotFound
	}
	return toDomainRecurringEntry(&m), nil
}

func (r *PGRecurringEntryRepo) List(ctx context.Context, companyID string) ([]domain.RecurringEntry, error) {
	var rows []pgRecurringEntry
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.RecurringEntry, len(rows))
	for i, m := range rows {
		out[i] = *toDomainRecurringEntry(&m)
	}
	return out, nil
}

func (r *PGRecurringEntryRepo) Update(ctx context.Context, entry *domain.RecurringEntry) error {
	m := toPGRecurringEntry(entry)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *PGRecurringEntryRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&pgRecurringEntry{}, "id = ?", id).Error
}

func (r *PGRecurringEntryRepo) UpdateNextRunDate(ctx context.Context, id, nextDate string) error {
	return r.db.WithContext(ctx).Model(&pgRecurringEntry{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"next_run_date": nextDate,
			"updated_at":    time.Now(),
		}).Error
}

func (r *PGRecurringEntryRepo) GetDueEntries(ctx context.Context, today string) ([]domain.RecurringEntry, error) {
	var rows []pgRecurringEntry
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND next_run_date != '' AND next_run_date <= ?", true, today).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.RecurringEntry, len(rows))
	for i, m := range rows {
		out[i] = *toDomainRecurringEntry(&m)
	}
	return out, nil
}

func toPGRecurringEntry(e *domain.RecurringEntry) *pgRecurringEntry {
	b, _ := json.Marshal(e.Lines)
	return &pgRecurringEntry{
		ID:           e.ID,
		CompanyID:    e.CompanyID,
		TemplateName: e.TemplateName,
		Description:  e.Description,
		Frequency:    string(e.Frequency),
		DayOfMonth:   e.DayOfMonth,
		IsActive:     e.IsActive,
		NextRunDate:  e.NextRunDate,
		LinesJSON:    string(b),
		CreatedBy:    e.CreatedBy,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func toDomainRecurringEntry(m *pgRecurringEntry) *domain.RecurringEntry {
	e := &domain.RecurringEntry{
		ID:           m.ID,
		CompanyID:    m.CompanyID,
		TemplateName: m.TemplateName,
		Description:  m.Description,
		Frequency:    domain.RecurringFrequency(m.Frequency),
		DayOfMonth:   m.DayOfMonth,
		IsActive:     m.IsActive,
		NextRunDate:  m.NextRunDate,
		CreatedBy:    m.CreatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
	if m.LinesJSON != "" {
		_ = json.Unmarshal([]byte(m.LinesJSON), &e.Lines)
	}
	return e
}
