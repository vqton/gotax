package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type FiscalYearService struct {
	periodRepo domain.PeriodRepository
	optRepo    domain.SystemOptionRepository
}

func NewFiscalYearService(periodRepo domain.PeriodRepository, optRepo domain.SystemOptionRepository) *FiscalYearService {
	return &FiscalYearService{periodRepo: periodRepo, optRepo: optRepo}
}

// CreateYear generates 12 monthly periods for a fiscal year starting at startMonth.
func (s *FiscalYearService) CreateYear(ctx context.Context, companyID string, year, startMonth int) ([]domain.Period, error) {
	if year < 2000 || year > 2100 {
		return nil, fmt.Errorf("year out of range")
	}
	if startMonth < 1 || startMonth > 12 {
		return nil, fmt.Errorf("start_month must be 1-12")
	}

	var periods []domain.Period
	for i := 0; i < 12; i++ {
		month := startMonth + i
		periodYear := year
		if month > 12 {
			month -= 12
			periodYear++
		}

		startDate := time.Date(periodYear, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, -1)

		p := &domain.Period{
			ID:        fmt.Sprintf("PER-%d-%02d", periodYear, month),
			Year:      periodYear,
			Month:     month,
			StartDate: startDate,
			EndDate:   endDate,
			Status:    domain.PeriodOpen,
		}

		if err := s.periodRepo.Create(ctx, p); err != nil {
			return periods, fmt.Errorf("create period %d/%02d: %w", periodYear, month, err)
		}
		periods = append(periods, *p)
	}
	return periods, nil
}

// CreateYearFromOptions reads the fiscal year start month from system options.
func (s *FiscalYearService) CreateYearFromOptions(ctx context.Context, companyID string, year int) ([]domain.Period, error) {
	startMonth := 1 // default January
	opt, err := s.optRepo.Get(ctx, companyID, "global", "fiscal_year_start")
	if err == nil && opt != nil {
		switch opt.Value {
		case "1":
			startMonth = 1
		case "4":
			startMonth = 4
		case "7":
			startMonth = 7
		case "10":
			startMonth = 10
		}
	}
	return s.CreateYear(ctx, companyID, year, startMonth)
}

func (s *FiscalYearService) ListPeriods(ctx context.Context, companyID string) ([]domain.Period, error) {
	return s.periodRepo.GetAll(ctx)
}
