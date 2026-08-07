package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type CostObjectService struct {
	repo domain.CostObjectRepository
}

func NewCostObjectService(repo domain.CostObjectRepository) *CostObjectService {
	return &CostObjectService{repo: repo}
}

func (s *CostObjectService) Create(ctx context.Context, co *domain.CostObject) error {
	if co.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if co.Code == "" {
		return domain.ErrCostObjectCodeRequired
	}
	if co.Name == "" {
		return domain.ErrCostObjectNameRequired
	}
	if !isValidCostObjectType(co.Type) {
		return domain.ErrCostObjectTypeInvalid
	}
	if !isValidCostingMethod(co.CostingMethod) {
		return domain.ErrCostObjectMethodInvalid
	}
	existing, err := s.repo.GetByCode(ctx, co.CompanyID, co.Code)
	if err == nil && existing != nil {
		return domain.ErrCostObjectCodeExists
	}
	if co.ID == "" {
		co.ID = fmt.Sprintf("COBJ-%d", time.Now().UnixNano())
	}
	now := time.Now().Format("2006-01-02T15:04:05Z")
	co.CreatedAt = now
	co.UpdatedAt = now
	co.IsActive = true
	return s.repo.Create(ctx, co)
}

func (s *CostObjectService) GetByID(ctx context.Context, id string) (*domain.CostObject, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CostObjectService) List(ctx context.Context, companyID string) ([]domain.CostObject, error) {
	return s.repo.List(ctx, companyID)
}

func (s *CostObjectService) Update(ctx context.Context, co *domain.CostObject) error {
	if co.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if co.Code == "" {
		return domain.ErrCostObjectCodeRequired
	}
	if co.Name == "" {
		return domain.ErrCostObjectNameRequired
	}
	if !isValidCostObjectType(co.Type) {
		return domain.ErrCostObjectTypeInvalid
	}
	if !isValidCostingMethod(co.CostingMethod) {
		return domain.ErrCostObjectMethodInvalid
	}
	existing, err := s.repo.GetByID(ctx, co.ID)
	if err != nil {
		return err
	}
	if existing.Code != co.Code {
		dup, err := s.repo.GetByCode(ctx, co.CompanyID, co.Code)
		if err == nil && dup != nil {
			return domain.ErrCostObjectCodeExists
		}
	}
	co.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")
	co.CreatedAt = existing.CreatedAt
	co.IsActive = existing.IsActive
	return s.repo.Update(ctx, co)
}

func (s *CostObjectService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func isValidCostObjectType(t domain.CostObjectType) bool {
	switch t {
	case domain.CostObjectTypeProduct, domain.CostObjectTypeService,
		domain.CostObjectTypeProject, domain.CostObjectTypeOrder,
		domain.CostObjectTypeBatch:
		return true
	default:
		return false
	}
}

func isValidCostingMethod(m domain.CostingMethod) bool {
	switch m {
	case domain.CostingMethodSimple, domain.CostingMethodCoefficient,
		domain.CostingMethodProportion, domain.CostingMethodStandard,
		domain.CostingMethodProcess, domain.CostingMethodByProduct:
		return true
	default:
		return false
	}
}
