package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type CostPoolService struct {
	poolRepo domain.CostPoolRepository
	lineRepo domain.CostPoolLineRepository
}

func NewCostPoolService(poolRepo domain.CostPoolRepository, lineRepo domain.CostPoolLineRepository) *CostPoolService {
	return &CostPoolService{poolRepo: poolRepo, lineRepo: lineRepo}
}

func (s *CostPoolService) Create(ctx context.Context, cp *domain.CostPool) error {
	if cp.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if cp.PeriodID == "" {
		return fmt.Errorf("period_id is required")
	}
	if cp.GLAccountCode == "" {
		return domain.ErrCostPoolAccountRequired
	}
	if !isValidGLAccountCode(cp.GLAccountCode) {
		return domain.ErrCostPoolAccountRequired
	}
	if cp.Name == "" {
		return fmt.Errorf("name is required")
	}
	if cp.ID == "" {
		cp.ID = fmt.Sprintf("CPOOL-%d", time.Now().UnixNano())
	}
	now := time.Now().Format("2006-01-02T15:04:05Z")
	cp.CreatedAt = now
	cp.UpdatedAt = now
	cp.Status = domain.CostPoolStatusOpen
	cp.TotalAmount = 0
	return s.poolRepo.Create(ctx, cp)
}

func (s *CostPoolService) GetByID(ctx context.Context, id string) (*domain.CostPool, error) {
	return s.poolRepo.GetByID(ctx, id)
}

func (s *CostPoolService) ListByPeriod(ctx context.Context, companyID, periodID string) ([]domain.CostPool, error) {
	return s.poolRepo.ListByPeriod(ctx, companyID, periodID)
}

func (s *CostPoolService) AddLine(ctx context.Context, line *domain.CostPoolLine) error {
	if line.PoolID == "" {
		return fmt.Errorf("pool_id is required")
	}
	if line.SourceType == "" {
		return fmt.Errorf("source_type is required")
	}
	if !isValidSourceType(line.SourceType) {
		return fmt.Errorf("invalid source_type")
	}
	if line.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	pool, err := s.poolRepo.GetByID(ctx, line.PoolID)
	if err != nil {
		return err
	}
	if pool.Status == domain.CostPoolStatusClosed {
		return domain.ErrCostPoolAlreadyClosed
	}
	if line.ID == "" {
		line.ID = fmt.Sprintf("CPL-%d", time.Now().UnixNano())
	}
	line.CreatedAt = time.Now().Format("2006-01-02T15:04:05Z")
	if err := s.lineRepo.Create(ctx, line); err != nil {
		return err
	}
	pool.TotalAmount += line.Amount
	pool.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")
	return s.poolRepo.Update(ctx, pool)
}

func (s *CostPoolService) ListLines(ctx context.Context, poolID string) ([]domain.CostPoolLine, error) {
	return s.lineRepo.ListByPool(ctx, poolID)
}

func (s *CostPoolService) ClosePool(ctx context.Context, poolID string) error {
	pool, err := s.poolRepo.GetByID(ctx, poolID)
	if err != nil {
		return err
	}
	if pool.Status == domain.CostPoolStatusClosed {
		return domain.ErrCostPoolAlreadyClosed
	}
	pool.Status = domain.CostPoolStatusClosed
	pool.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")
	return s.poolRepo.Update(ctx, pool)
}

func (s *CostPoolService) Delete(ctx context.Context, id string) error {
	pool, err := s.poolRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if pool.Status == domain.CostPoolStatusClosed {
		return domain.ErrCostPoolAlreadyClosed
	}
	if err := s.lineRepo.DeleteByPool(ctx, id); err != nil {
		return err
	}
	return s.poolRepo.Delete(ctx, id)
}

func isValidGLAccountCode(code string) bool {
	switch code {
	case "621", "622", "623", "627":
		return true
	default:
		return false
	}
}

func isValidSourceType(t string) bool {
	switch t {
	case "JOURNAL", "PAYROLL", "DEPRECIATION", "MANUAL":
		return true
	default:
		return false
	}
}
