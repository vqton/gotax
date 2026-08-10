package repository

import (
	"context"
	"gotax/internal/domain"
	"time"

	"gorm.io/gorm"
)

type PGPriceListRepo struct {
	db *gorm.DB
}

func NewPGPriceListRepo(db *gorm.DB) *PGPriceListRepo {
	return &PGPriceListRepo{db: db}
}

func priceListToGORM(p *domain.PriceList) *domain.PriceListGORM {
	g := &domain.PriceListGORM{
		ID:          p.ID,
		CompanyID:   p.CompanyID,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		Currency:    p.Currency,
		IsActive:    p.IsActive,
	}
	if !p.CreatedAt.IsZero() { g.CreatedAt = p.CreatedAt }
	if !p.UpdatedAt.IsZero() { g.UpdatedAt = p.UpdatedAt }
	return g
}

func priceListFromGORM(g *domain.PriceListGORM) *domain.PriceList {
	p := &domain.PriceList{
		ID:          g.ID,
		CompanyID:   g.CompanyID,
		Code:        g.Code,
		Name:        g.Name,
		Description: g.Description,
		Currency:    g.Currency,
		IsActive:    g.IsActive,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
	return p
}

func priceListLineToGORM(l *domain.PriceListLine) *domain.PriceListLineGORM {
	return &domain.PriceListLineGORM{
		ID:            l.ID,
		PriceListID:   l.PriceListID,
		ItemCode:      l.ItemCode,
		ItemName:      l.ItemName,
		Unit:          l.Unit,
		UnitPrice:     l.UnitPrice,
		VATRate:       l.VATRate,
		MinQuantity:   l.MinQuantity,
		EffectiveFrom: l.EffectiveFrom,
		EffectiveTo:   l.EffectiveTo,
	}
}

func priceListLineFromGORM(g *domain.PriceListLineGORM) domain.PriceListLine {
	return domain.PriceListLine{
		ID:            g.ID,
		PriceListID:   g.PriceListID,
		ItemCode:      g.ItemCode,
		ItemName:      g.ItemName,
		Unit:          g.Unit,
		UnitPrice:     g.UnitPrice,
		VATRate:       g.VATRate,
		MinQuantity:   g.MinQuantity,
		EffectiveFrom: g.EffectiveFrom,
		EffectiveTo:   g.EffectiveTo,
	}
}

func (r *PGPriceListRepo) CreatePriceList(ctx context.Context, p *domain.PriceList) error {
	g := priceListToGORM(p)
	now := time.Now()
	g.CreatedAt, g.UpdatedAt = now, now
	return r.db.WithContext(ctx).Create(g).Error
}

func (r *PGPriceListRepo) GetPriceList(ctx context.Context, id string) (*domain.PriceList, error) {
	var g domain.PriceListGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&g, "id = ?", id).Error; err != nil {
		return nil, err
	}
	p := priceListFromGORM(&g)
	if len(g.Lines) > 0 {
		p.Lines = make([]domain.PriceListLine, len(g.Lines))
		for i, l := range g.Lines {
			p.Lines[i] = priceListLineFromGORM(&l)
		}
	}
	return p, nil
}

func (r *PGPriceListRepo) GetPriceListByCode(ctx context.Context, companyID, code string) (*domain.PriceList, error) {
	var g domain.PriceListGORM
	if err := r.db.WithContext(ctx).Preload("Lines").First(&g, "company_id = ? AND code = ?", companyID, code).Error; err != nil {
		return nil, err
	}
	p := priceListFromGORM(&g)
	if len(g.Lines) > 0 {
		p.Lines = make([]domain.PriceListLine, len(g.Lines))
		for i, l := range g.Lines {
			p.Lines[i] = priceListLineFromGORM(&l)
		}
	}
	return p, nil
}

func (r *PGPriceListRepo) ListPriceLists(ctx context.Context, companyID string) ([]domain.PriceList, error) {
	var rows []domain.PriceListGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.PriceList, len(rows))
	for i, g := range rows {
		out[i] = *priceListFromGORM(&g)
	}
	return out, nil
}

func (r *PGPriceListRepo) UpdatePriceList(ctx context.Context, p *domain.PriceList) error {
	g := priceListToGORM(p)
	g.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(g).Error
}

func (r *PGPriceListRepo) DeletePriceList(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("price_list_id = ?", id).Delete(&domain.PriceListLineGORM{}).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.PriceListGORM{}, "id = ?", id).Error
	})
}

func (r *PGPriceListRepo) CreatePriceListLines(ctx context.Context, lines []domain.PriceListLine) error {
	if len(lines) == 0 { return nil }
	gormLines := make([]domain.PriceListLineGORM, len(lines))
	for i, l := range lines {
		gormLines[i] = *priceListLineToGORM(&l)
	}
	return r.db.WithContext(ctx).Create(&gormLines).Error
}

func (r *PGPriceListRepo) GetPriceListLines(ctx context.Context, priceListID string) ([]domain.PriceListLine, error) {
	var rows []domain.PriceListLineGORM
	if err := r.db.WithContext(ctx).Where("price_list_id = ?", priceListID).Order("item_code").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.PriceListLine, len(rows))
	for i, g := range rows {
		out[i] = priceListLineFromGORM(&g)
	}
	return out, nil
}

func (r *PGPriceListRepo) UpdatePriceListLine(ctx context.Context, l *domain.PriceListLine) error {
	g := priceListLineToGORM(l)
	return r.db.WithContext(ctx).Save(g).Error
}

func (r *PGPriceListRepo) DeletePriceListLines(ctx context.Context, priceListID string) error {
	return r.db.WithContext(ctx).Where("price_list_id = ?", priceListID).Delete(&domain.PriceListLineGORM{}).Error
}
