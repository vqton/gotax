package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type ccdcCatGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36"`
	Code        string    `gorm:"column:code;not null;size:50"`
	Name        string    `gorm:"column:name;not null;size:255"`
	Description *string   `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ccdcCatGORM) TableName() string { return "tool_equipment_categories" }

type ccdcItemGORM struct {
	ID           string     `gorm:"column:id;primaryKey;size:36"`
	CompanyID    string     `gorm:"column:company_id;not null;size:36"`
	Code         string     `gorm:"column:code;not null;size:50"`
	Name         string     `gorm:"column:name;not null;size:255"`
	Category     *string    `gorm:"column:category;size:255"`
	PurchaseDate *time.Time `gorm:"column:purchase_date;type:date"`
	PurchaseCost float64    `gorm:"column:purchase_cost;default:0"`
	CurrentCost  float64    `gorm:"column:current_cost;default:0"`
	WarehouseID  *string    `gorm:"column:warehouse_id;size:36"`
	Status       string     `gorm:"column:status;not null;size:20;default:ACTIVE"`
	Description  *string    `gorm:"column:description;type:text"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (ccdcItemGORM) TableName() string { return "tool_equipment" }

// ─── Category Repo ───────────────────────────────────────────────────

type PGCCDCCategoryRepo struct{ db *gorm.DB }

func NewPGCCDCCategoryRepo(db *gorm.DB) *PGCCDCCategoryRepo {
	return &PGCCDCCategoryRepo{db}
}

func (r *PGCCDCCategoryRepo) Create(ctx context.Context, c *domain.ToolEquipmentCategory) error {
	m := ccdcCatToGORM(c)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGCCDCCategoryRepo) GetByID(ctx context.Context, id string) (*domain.ToolEquipmentCategory, error) {
	var m ccdcCatGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return ccdcCatToDomain(&m), nil
}

func (r *PGCCDCCategoryRepo) List(ctx context.Context, companyID string) ([]domain.ToolEquipmentCategory, error) {
	var models []ccdcCatGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ToolEquipmentCategory, len(models))
	for i := range models {
		out[i] = *ccdcCatToDomain(&models[i])
	}
	return out, nil
}

func (r *PGCCDCCategoryRepo) Update(ctx context.Context, c *domain.ToolEquipmentCategory) error {
	m := ccdcCatToGORM(c)
	return r.db.WithContext(ctx).Model(&ccdcCatGORM{}).
		Select("company_id", "code", "name", "description").Where("id = ?", c.ID).Updates(m).Error
}

func (r *PGCCDCCategoryRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ccdcCatGORM{}).Error
}

// ─── Item Repo ───────────────────────────────────────────────────────

type PGToolEquipmentRepo struct{ db *gorm.DB }

func NewPGToolEquipmentRepo(db *gorm.DB) *PGToolEquipmentRepo {
	return &PGToolEquipmentRepo{db}
}

func (r *PGToolEquipmentRepo) Create(ctx context.Context, t *domain.ToolEquipment) error {
	m := ccdcItemToGORM(t)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PGToolEquipmentRepo) GetByID(ctx context.Context, id string) (*domain.ToolEquipment, error) {
	var m ccdcItemGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return ccdcItemToDomain(&m), nil
}

func (r *PGToolEquipmentRepo) List(ctx context.Context, companyID string) ([]domain.ToolEquipment, error) {
	var models []ccdcItemGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("code").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ToolEquipment, len(models))
	for i := range models {
		out[i] = *ccdcItemToDomain(&models[i])
	}
	return out, nil
}

func (r *PGToolEquipmentRepo) Update(ctx context.Context, t *domain.ToolEquipment) error {
	m := ccdcItemToGORM(t)
	return r.db.WithContext(ctx).Model(&ccdcItemGORM{}).
		Select("company_id", "code", "name", "category", "purchase_date",
			"purchase_cost", "current_cost", "warehouse_id", "status", "description", "updated_at").
		Where("id = ?", t.ID).Updates(m).Error
}

func (r *PGToolEquipmentRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ccdcItemGORM{}).Error
}

func (r *PGToolEquipmentRepo) GetByCode(ctx context.Context, companyID, code string) (*domain.ToolEquipment, error) {
	var m ccdcItemGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND code = ?", companyID, code).First(&m).Error; err != nil {
		return nil, err
	}
	return ccdcItemToDomain(&m), nil
}

// ─── Conversion helpers ─────────────────────────────────────────────

func ccdcCatToDomain(m *ccdcCatGORM) *domain.ToolEquipmentCategory {
	return &domain.ToolEquipmentCategory{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		Code:        m.Code,
		Name:        m.Name,
		Description: safeStr(m.Description),
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
	}
}

func ccdcCatToGORM(d *domain.ToolEquipmentCategory) *ccdcCatGORM {
	return &ccdcCatGORM{
		ID:          d.ID,
		CompanyID:   d.CompanyID,
		Code:        d.Code,
		Name:        d.Name,
		Description: strPtr(d.Description),
	}
}

func ccdcItemToDomain(m *ccdcItemGORM) *domain.ToolEquipment {
	return &domain.ToolEquipment{
		ID:           m.ID,
		CompanyID:    m.CompanyID,
		Code:         m.Code,
		Name:         m.Name,
		Category:     safeStr(m.Category),
		PurchaseDate: safeTimePtrStr(m.PurchaseDate),
		PurchaseCost: m.PurchaseCost,
		CurrentCost:  m.CurrentCost,
		WarehouseID:  safeStr(m.WarehouseID),
		Status:       domain.CCDCStatus(m.Status),
		Description:  safeStr(m.Description),
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}

func ccdcItemToGORM(d *domain.ToolEquipment) *ccdcItemGORM {
	return &ccdcItemGORM{
		ID:           d.ID,
		CompanyID:    d.CompanyID,
		Code:         d.Code,
		Name:         d.Name,
		Category:     strPtr(d.Category),
		PurchaseDate: timePtr(parseDate(d.PurchaseDate)),
		PurchaseCost: d.PurchaseCost,
		CurrentCost:  d.CurrentCost,
		WarehouseID:  strPtr(d.WarehouseID),
		Status:       string(d.Status),
		Description:  strPtr(d.Description),
	}
}
