package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type CostingPeriodService struct {
	periodRepo domain.CostingPeriodRepository
}

func NewCostingPeriodService(periodRepo domain.CostingPeriodRepository) *CostingPeriodService {
	return &CostingPeriodService{periodRepo: periodRepo}
}

func (s *CostingPeriodService) Create(ctx context.Context, cp *domain.CostingPeriod) error {
	if cp.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if cp.Year < 2000 || cp.Year > 2100 {
		return fmt.Errorf("year must be between 2000 and 2100")
	}
	if cp.Month < 1 || cp.Month > 12 {
		return fmt.Errorf("month must be between 1 and 12")
	}

	existing, err := s.periodRepo.GetByYearMonth(ctx, cp.CompanyID, cp.Year, cp.Month)
	if err == nil && existing != nil {
		return domain.ErrCostingPeriodExists
	}

	if cp.ID == "" {
		cp.ID = fmt.Sprintf("CP-%d", time.Now().UnixNano())
	}
	cp.Status = "OPEN"
	now := time.Now().Format("2006-01-02T15:04:05Z")
	cp.CreatedAt = now

	return s.periodRepo.Create(ctx, cp)
}

func (s *CostingPeriodService) GetByID(ctx context.Context, id string) (*domain.CostingPeriod, error) {
	return s.periodRepo.GetByID(ctx, id)
}

func (s *CostingPeriodService) List(ctx context.Context, companyID string) ([]domain.CostingPeriod, error) {
	if companyID == "" {
		return nil, fmt.Errorf("company_id is required")
	}
	return s.periodRepo.List(ctx, companyID)
}

func (s *CostingPeriodService) Close(ctx context.Context, id, closedBy string) error {
	cp, err := s.periodRepo.GetByID(ctx, id)
	if err != nil {
		return domain.ErrCostingPeriodNotFound
	}
	if cp.Status == "CLOSED" {
		return domain.ErrCostingPeriodAlreadyClosed
	}
	if closedBy == "" {
		return fmt.Errorf("closed_by is required")
	}

	cp.Status = "CLOSED"
	cp.ClosedBy = closedBy
	now := time.Now().Format("2006-01-02T15:04:05Z")
	cp.ClosedAt = now

	return s.periodRepo.Update(ctx, cp)
}
