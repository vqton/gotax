package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type CostCenterService struct {
	repo domain.CostCenterRepository
}

func NewCostCenterService(repo domain.CostCenterRepository) *CostCenterService {
	return &CostCenterService{repo: repo}
}

func (s *CostCenterService) Create(ctx context.Context, cc *domain.CostCenter) error {
	if cc.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if cc.Code == "" {
		return fmt.Errorf("code is required")
	}
	if cc.Name == "" {
		return fmt.Errorf("name is required")
	}
	existing, err := s.repo.GetByCode(ctx, cc.CompanyID, cc.Code)
	if err == nil && existing != nil {
		return domain.ErrCostCenterCodeExists
	}
	if cc.ID == "" {
		cc.ID = fmt.Sprintf("CC-%d", time.Now().UnixNano())
	}
	now := time.Now().Format("2006-01-02T15:04:05Z")
	cc.CreatedAt = now
	cc.UpdatedAt = now
	cc.IsActive = true
	return s.repo.Create(ctx, cc)
}

func (s *CostCenterService) GetByID(ctx context.Context, id string) (*domain.CostCenter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CostCenterService) List(ctx context.Context, companyID string) ([]domain.CostCenter, error) {
	return s.repo.List(ctx, companyID)
}

func (s *CostCenterService) Update(ctx context.Context, cc *domain.CostCenter) error {
	if cc.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if cc.Code == "" {
		return fmt.Errorf("code is required")
	}
	if cc.Name == "" {
		return fmt.Errorf("name is required")
	}
	existing, err := s.repo.GetByID(ctx, cc.ID)
	if err != nil {
		return err
	}
	if existing.Code != cc.Code {
		dup, err := s.repo.GetByCode(ctx, cc.CompanyID, cc.Code)
		if err == nil && dup != nil {
			return domain.ErrCostCenterCodeExists
		}
	}
	cc.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")
	cc.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, cc)
}

func (s *CostCenterService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *CostCenterService) ListHierarchy(ctx context.Context, companyID string) ([]domain.CostCenter, error) {
	return s.repo.List(ctx, companyID)
}
