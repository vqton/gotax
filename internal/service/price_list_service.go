package service

import (
	"context"
	"time"

	"gotax/internal/domain"
)

type PriceListService struct {
	repo domain.PriceListRepository
}

func NewPriceListService(repo domain.PriceListRepository) *PriceListService {
	return &PriceListService{repo: repo}
}

func (s *PriceListService) CreatePriceList(ctx context.Context, pl *domain.PriceList) error {
	if err := pl.Validate(); err != nil {
		return err
	}
	existing, _ := s.repo.GetPriceListByCode(ctx, pl.CompanyID, pl.Code)
	if existing != nil {
		return domain.ErrPriceListCodeRequired // reuse error for duplicate
	}
	return s.repo.CreatePriceList(ctx, pl)
}

func (s *PriceListService) GetPriceList(ctx context.Context, id string) (*domain.PriceList, error) {
	return s.repo.GetPriceList(ctx, id)
}

func (s *PriceListService) ListPriceLists(ctx context.Context, companyID string) ([]domain.PriceList, error) {
	return s.repo.ListPriceLists(ctx, companyID)
}

func (s *PriceListService) UpdatePriceList(ctx context.Context, pl *domain.PriceList) error {
	if err := pl.Validate(); err != nil {
		return err
	}
	pl.UpdatedAt = time.Now()
	return s.repo.UpdatePriceList(ctx, pl)
}

func (s *PriceListService) DeletePriceList(ctx context.Context, id string) error {
	return s.repo.DeletePriceList(ctx, id)
}

func (s *PriceListService) AddLines(ctx context.Context, priceListID string, lines []domain.PriceListLine) error {
	for i := range lines {
		lines[i].PriceListID = priceListID
		if err := lines[i].Validate(); err != nil {
			return err
		}
	}
	return s.repo.CreatePriceListLines(ctx, lines)
}

func (s *PriceListService) GetLines(ctx context.Context, priceListID string) ([]domain.PriceListLine, error) {
	return s.repo.GetPriceListLines(ctx, priceListID)
}

// CalculateSellingPrice returns the unit price for an item from a price list,
// applying optional markup percentage.
func (s *PriceListService) CalculateSellingPrice(ctx context.Context, priceListID, itemCode string, markupPct float64) (float64, error) {
	lines, err := s.repo.GetPriceListLines(ctx, priceListID)
	if err != nil {
		return 0, err
	}
	for _, l := range lines {
		if l.ItemCode == itemCode {
			price := l.UnitPrice * (1 + markupPct/100)
			return price, nil
		}
	}
	return 0, domain.ErrPriceListItemCodeRequired
}
