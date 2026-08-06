package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type pgBudget struct {
	ID          string    `gorm:"column:id;primaryKey"`
	CompanyID   string    `gorm:"column:company_id"`
	AccountCode string    `gorm:"column:account_code"`
	PeriodYear  int       `gorm:"column:period_year"`
	PeriodMonth int       `gorm:"column:period_month"`
	Budgeted    float64   `gorm:"column:budgeted"`
	Actual      float64   `gorm:"column:actual"`
	Variance    float64   `gorm:"column:variance"`
	Notes       string    `gorm:"column:notes"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (pgBudget) TableName() string { return "budgets" }

type PGBudgetRepo struct{ db *gorm.DB }

func NewPGBudgetRepo(db *gorm.DB) *PGBudgetRepo {
	return &PGBudgetRepo{db: db}
}

func (r *PGBudgetRepo) Create(ctx context.Context, b *domain.Budget) error {
	m := toPGBudget(b)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *PGBudgetRepo) GetByID(ctx context.Context, id string) (*domain.Budget, error) {
	var m pgBudget
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrBudgetNotFound
	}
	return toDomainBudget(&m), nil
}

func (r *PGBudgetRepo) List(ctx context.Context, companyID string, year int) ([]domain.Budget, error) {
	var rows []pgBudget
	q := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("account_code, period_month")
	if year > 0 {
		q = q.Where("period_year = ?", year)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Budget, len(rows))
	for i, m := range rows {
		out[i] = *toDomainBudget(&m)
	}
	return out, nil
}

func (r *PGBudgetRepo) Update(ctx context.Context, b *domain.Budget) error {
	m := toPGBudget(b)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *PGBudgetRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&pgBudget{}, "id = ?", id).Error
}

func (r *PGBudgetRepo) Upsert(ctx context.Context, b *domain.Budget) error {
	m := toPGBudget(b)
	return r.db.WithContext(ctx).
		Where("company_id = ? AND account_code = ? AND period_year = ? AND period_month = ?",
			m.CompanyID, m.AccountCode, m.PeriodYear, m.PeriodMonth).
		Assign(map[string]interface{}{
			"budgeted":  m.Budgeted,
			"actual":    m.Actual,
			"variance":  m.Variance,
			"notes":     m.Notes,
			"updated_at": time.Now(),
		}).
		FirstOrCreate(&m).Error
}

func toPGBudget(b *domain.Budget) *pgBudget {
	var created, updated time.Time
	if b.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", b.CreatedAt)
	} else {
		created = time.Now()
	}
	if b.UpdatedAt != "" {
		updated, _ = time.Parse("2006-01-02T15:04:05Z", b.UpdatedAt)
	} else {
		updated = time.Now()
	}
	return &pgBudget{
		ID:          b.ID,
		CompanyID:   b.CompanyID,
		AccountCode: b.AccountCode,
		PeriodYear:  b.PeriodYear,
		PeriodMonth: b.PeriodMonth,
		Budgeted:    b.Budgeted,
		Actual:      b.Actual,
		Variance:    b.Variance,
		Notes:       b.Notes,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}
}

func toDomainBudget(m *pgBudget) *domain.Budget {
	return &domain.Budget{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		AccountCode: m.AccountCode,
		PeriodYear:  m.PeriodYear,
		PeriodMonth: m.PeriodMonth,
		Budgeted:    m.Budgeted,
		Actual:      m.Actual,
		Variance:    m.Variance,
		Notes:       m.Notes,
		CreatedAt:   m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
