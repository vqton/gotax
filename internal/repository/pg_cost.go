package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type pgCostCenter struct {
	ID          string     `gorm:"column:id;primaryKey"`
	CompanyID   string     `gorm:"column:company_id"`
	Code        string     `gorm:"column:code"`
	Name        string     `gorm:"column:name"`
	ParentID    *string    `gorm:"column:parent_id"`
	Description *string    `gorm:"column:description"`
	IsActive    bool       `gorm:"column:is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (pgCostCenter) TableName() string { return "cost_centers" }

type PGCostCenterRepo struct{ db *gorm.DB }

func NewPGCostCenterRepo(db *gorm.DB) *PGCostCenterRepo {
	return &PGCostCenterRepo{db: db}
}

func toPGCC(cc *domain.CostCenter) *pgCostCenter {
	var created, updated time.Time
	if cc.CreatedAt != "" {
		created, _ = time.Parse("2006-01-02T15:04:05Z", cc.CreatedAt)
	} else {
		created = time.Now()
	}
	if cc.UpdatedAt != "" {
		updated, _ = time.Parse("2006-01-02T15:04:05Z", cc.UpdatedAt)
	} else {
		updated = time.Now()
	}
	m := &pgCostCenter{
		ID:        cc.ID,
		CompanyID: cc.CompanyID,
		Code:      cc.Code,
		Name:      cc.Name,
		IsActive:  cc.IsActive,
		CreatedAt: created,
		UpdatedAt: updated,
	}
	if cc.ParentID != "" {
		m.ParentID = &cc.ParentID
	}
	if cc.Description != "" {
		m.Description = &cc.Description
	}
	return m
}

func toDomainCC(m *pgCostCenter) *domain.CostCenter {
	cc := &domain.CostCenter{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Code:      m.Code,
		Name:      m.Name,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if m.ParentID != nil {
		cc.ParentID = *m.ParentID
	}
	if m.Description != nil {
		cc.Description = *m.Description
	}
	return cc
}

func (r *PGCostCenterRepo) Create(ctx context.Context, cc *domain.CostCenter) error {
	return r.db.WithContext(ctx).Create(toPGCC(cc)).Error
}

func (r *PGCostCenterRepo) GetByID(ctx context.Context, id string) (*domain.CostCenter, error) {
	var m pgCostCenter
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCostCenterNotFound
	}
	return toDomainCC(&m), nil
}

func (r *PGCostCenterRepo) List(ctx context.Context, companyID string) ([]domain.CostCenter, error) {
	var rows []pgCostCenter
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CostCenter, len(rows))
	for i := range rows {
		out[i] = *toDomainCC(&rows[i])
	}
	return out, nil
}

func (r *PGCostCenterRepo) Update(ctx context.Context, cc *domain.CostCenter) error {
	return r.db.WithContext(ctx).Save(toPGCC(cc)).Error
}

func (r *PGCostCenterRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&pgCostCenter{}, "id = ?", id).Error
}

func (r *PGCostCenterRepo) GetByCode(ctx context.Context, companyID, code string) (*domain.CostCenter, error) {
	var m pgCostCenter
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		return nil, domain.ErrCostCenterNotFound
	}
	return toDomainCC(&m), nil
}
