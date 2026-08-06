package service

import (
	"context"
	"time"

	"gotax/internal/domain"
)

type CCDCServiceInterface interface {
	CreateCategory(ctx context.Context, c *domain.ToolEquipmentCategory) error
	GetCategory(ctx context.Context, id string) (*domain.ToolEquipmentCategory, error)
	ListCategories(ctx context.Context, companyID string) ([]domain.ToolEquipmentCategory, error)
	UpdateCategory(ctx context.Context, c *domain.ToolEquipmentCategory) error
	DeleteCategory(ctx context.Context, id string) error

	Create(ctx context.Context, t *domain.ToolEquipment) error
	GetByID(ctx context.Context, id string) (*domain.ToolEquipment, error)
	List(ctx context.Context, companyID string) ([]domain.ToolEquipment, error)
	Update(ctx context.Context, t *domain.ToolEquipment) error
	Delete(ctx context.Context, id string) error
}

type ccdcService struct {
	catRepo domain.ToolEquipmentCategoryRepository
	itemRepo domain.ToolEquipmentRepository
	now     func() time.Time
}

func NewCCDCService(catRepo domain.ToolEquipmentCategoryRepository, itemRepo domain.ToolEquipmentRepository) CCDCServiceInterface {
	return &ccdcService{
		catRepo: catRepo,
		itemRepo: itemRepo,
		now:     time.Now,
	}
}

// ─── Categories ──────────────────────────────────────────────────────

func (s *ccdcService) CreateCategory(ctx context.Context, c *domain.ToolEquipmentCategory) error {
	if c.Code == "" {
		return domain.ErrCCDCCategoryCodeRequired
	}
	if c.Name == "" {
		return domain.ErrCCDCCategoryNameRequired
	}
	now := s.now()
	c.CreatedAt = now.Format(time.RFC3339)
	return s.catRepo.Create(ctx, c)
}

func (s *ccdcService) GetCategory(ctx context.Context, id string) (*domain.ToolEquipmentCategory, error) {
	return s.catRepo.GetByID(ctx, id)
}

func (s *ccdcService) ListCategories(ctx context.Context, companyID string) ([]domain.ToolEquipmentCategory, error) {
	return s.catRepo.List(ctx, companyID)
}

func (s *ccdcService) UpdateCategory(ctx context.Context, c *domain.ToolEquipmentCategory) error {
	existing, err := s.catRepo.GetByID(ctx, c.ID)
	if err != nil {
		return err
	}
	c.CompanyID = existing.CompanyID
	c.CreatedAt = existing.CreatedAt
	return s.catRepo.Update(ctx, c)
}

func (s *ccdcService) DeleteCategory(ctx context.Context, id string) error {
	return s.catRepo.Delete(ctx, id)
}

// ─── Items ───────────────────────────────────────────────────────────

func (s *ccdcService) Create(ctx context.Context, t *domain.ToolEquipment) error {
	if t.Code == "" {
		return domain.ErrCCDCItemCodeRequired
	}
	if t.Name == "" {
		return domain.ErrCCDCItemNameRequired
	}
	now := s.now()
	t.CreatedAt = now.Format(time.RFC3339)
	t.UpdatedAt = now.Format(time.RFC3339)
	if t.Status == "" {
		t.Status = domain.CCDCActive
	}
	if t.CurrentCost == 0 {
		t.CurrentCost = t.PurchaseCost
	}
	return s.itemRepo.Create(ctx, t)
}

func (s *ccdcService) GetByID(ctx context.Context, id string) (*domain.ToolEquipment, error) {
	return s.itemRepo.GetByID(ctx, id)
}

func (s *ccdcService) List(ctx context.Context, companyID string) ([]domain.ToolEquipment, error) {
	return s.itemRepo.List(ctx, companyID)
}

func (s *ccdcService) Update(ctx context.Context, t *domain.ToolEquipment) error {
	existing, err := s.itemRepo.GetByID(ctx, t.ID)
	if err != nil {
		return err
	}
	t.CompanyID = existing.CompanyID
	t.CreatedAt = existing.CreatedAt
	t.UpdatedAt = s.now().Format(time.RFC3339)
	return s.itemRepo.Update(ctx, t)
}

func (s *ccdcService) Delete(ctx context.Context, id string) error {
	return s.itemRepo.Delete(ctx, id)
}
