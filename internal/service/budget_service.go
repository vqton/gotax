package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type BudgetService struct {
	repo   domain.BudgetRepository
	jeRepo domain.JournalRepository
}

func NewBudgetService(repo domain.BudgetRepository, jeRepo domain.JournalRepository) *BudgetService {
	return &BudgetService{repo: repo, jeRepo: jeRepo}
}

func (s *BudgetService) Create(ctx context.Context, b *domain.Budget) error {
	if b.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if b.AccountCode == "" {
		return fmt.Errorf("account_code is required")
	}
	if b.PeriodYear < 2000 || b.PeriodYear > 2100 {
		return fmt.Errorf("period_year must be 2000-2100")
	}
	if b.PeriodMonth < 1 || b.PeriodMonth > 12 {
		return fmt.Errorf("period_month must be 1-12")
	}
	b.Variance = b.Actual - b.Budgeted
	if b.ID == "" {
		b.ID = fmt.Sprintf("BUD-%d", time.Now().UnixNano())
	}
	now := time.Now().Format("2006-01-02T15:04:05Z")
	b.CreatedAt = now
	b.UpdatedAt = now
	return s.repo.Create(ctx, b)
}

func (s *BudgetService) GetByID(ctx context.Context, id string) (*domain.Budget, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BudgetService) List(ctx context.Context, companyID string, year int) ([]domain.Budget, error) {
	return s.repo.List(ctx, companyID, year)
}

func (s *BudgetService) Update(ctx context.Context, b *domain.Budget) error {
	if b.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if b.AccountCode == "" {
		return fmt.Errorf("account_code is required")
	}
	b.Variance = b.Actual - b.Budgeted
	b.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")
	return s.repo.Update(ctx, b)
}

func (s *BudgetService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *BudgetService) Upsert(ctx context.Context, b *domain.Budget) error {
	if b.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if b.AccountCode == "" {
		return fmt.Errorf("account_code is required")
	}
	b.Variance = b.Actual - b.Budgeted
	now := time.Now().Format("2006-01-02T15:04:05Z")
	if b.CreatedAt == "" {
		b.CreatedAt = now
	}
	b.UpdatedAt = now
	return s.repo.Upsert(ctx, b)
}

func (s *BudgetService) BulkUpsert(ctx context.Context, budgets []*domain.Budget) (int, error) {
	count := 0
	for _, b := range budgets {
		if err := s.Upsert(ctx, b); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

type BudgetVarianceItem struct {
	AccountCode string  `json:"account_code"`
	PeriodYear  int     `json:"period_year"`
	PeriodMonth int     `json:"period_month"`
	Budgeted    float64 `json:"budgeted"`
	Actual      float64 `json:"actual"`
	Variance    float64 `json:"variance"`
	Notes       string  `json:"notes,omitempty"`
}

type BudgetVarianceReport struct {
	CompanyID string              `json:"company_id"`
	Year      int                 `json:"year"`
	Month     int                 `json:"month"`
	Items     []BudgetVarianceItem `json:"items"`
	TotalBudgeted float64          `json:"total_budgeted"`
	TotalActual   float64          `json:"total_actual"`
	TotalVariance float64          `json:"total_variance"`
}

func (s *BudgetService) VarianceReport(ctx context.Context, companyID string, year, month int) (*BudgetVarianceReport, error) {
	budgets, err := s.repo.List(ctx, companyID, year)
	if err != nil {
		return nil, err
	}
	report := &BudgetVarianceReport{
		CompanyID: companyID,
		Year:      year,
		Month:     month,
	}
	for _, b := range budgets {
		if month > 0 && b.PeriodMonth != month {
			continue
		}
		v := b.Actual - b.Budgeted
		report.Items = append(report.Items, BudgetVarianceItem{
			AccountCode: b.AccountCode,
			PeriodYear:  b.PeriodYear,
			PeriodMonth: b.PeriodMonth,
			Budgeted:    b.Budgeted,
			Actual:      b.Actual,
			Variance:    v,
			Notes:       b.Notes,
		})
		report.TotalBudgeted += b.Budgeted
		report.TotalActual += b.Actual
		report.TotalVariance += v
	}
	return report, nil
}

func (s *BudgetService) SyncActuals(ctx context.Context, companyID string, year, month int) (int, error) {
	budgets, err := s.repo.List(ctx, companyID, year)
	if err != nil {
		return 0, err
	}
	updated := 0
	for i := range budgets {
		if budgets[i].PeriodMonth != month {
			continue
		}
		budgets[i].Actual = 0
		budgets[i].Variance = budgets[i].Actual - budgets[i].Budgeted
		budgets[i].UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")
		if err := s.repo.Update(ctx, &budgets[i]); err != nil {
			continue
		}
		updated++
	}
	return updated, nil
}
