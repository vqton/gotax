package repository

import (
	"context"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type pgNotificationRepo struct {
	db *gorm.DB
}

func NewPGNotificationRepo(db *gorm.DB) domain.NotificationRepository {
	return &pgNotificationRepo{db: db}
}

type notifModel struct {
	ID         string `gorm:"column:id;type:varchar(36);primaryKey"`
	CompanyID  string `gorm:"column:company_id;type:varchar(36);not null"`
	UserID     string `gorm:"column:user_id;type:varchar(36);not null"`
	EntityType string `gorm:"column:entity_type;type:varchar(50)"`
	EntityID   string `gorm:"column:entity_id;type:varchar(36)"`
	Type       string `gorm:"column:type;type:varchar(20);not null;default:INFO"`
	Title      string `gorm:"column:title;type:varchar(255);not null"`
	Message    string `gorm:"column:message;type:text;not null"`
	Link       string `gorm:"column:link;type:varchar(500)"`
	ReadAt     *string `gorm:"column:read_at"`
	CreatedAt  string `gorm:"column:created_at;not null"`
}

func (notifModel) TableName() string { return "notifications" }

func toNotifDomain(m notifModel) *domain.Notification {
	n := &domain.Notification{
		ID:         m.ID,
		CompanyID:  m.CompanyID,
		UserID:     m.UserID,
		EntityType: m.EntityType,
		EntityID:   m.EntityID,
		Type:       domain.NotificationType(m.Type),
		Title:      m.Title,
		Message:    m.Message,
		Link:       m.Link,
		CreatedAt:  m.CreatedAt,
	}
	if m.ReadAt != nil {
		n.ReadAt = *m.ReadAt
	}
	return n
}

func fromNotifDomain(n *domain.Notification) notifModel {
	m := notifModel{
		ID:         n.ID,
		CompanyID:  n.CompanyID,
		UserID:     n.UserID,
		EntityType: n.EntityType,
		EntityID:   n.EntityID,
		Type:       string(n.Type),
		Title:      n.Title,
		Message:    n.Message,
		Link:       n.Link,
		CreatedAt:  n.CreatedAt,
	}
	if n.ReadAt != "" {
		m.ReadAt = &n.ReadAt
	}
	return m
}

func (r *pgNotificationRepo) Create(ctx context.Context, n *domain.Notification) error {
	m := fromNotifDomain(n)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *pgNotificationRepo) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	var m notifModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return toNotifDomain(m), nil
}

func (r *pgNotificationRepo) ListByUser(ctx context.Context, companyID, userID string, limit int) ([]domain.Notification, error) {
	var models []notifModel
	q := r.db.WithContext(ctx).Where("company_id = ? AND user_id = ?", companyID, userID).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]domain.Notification, len(models))
	for i, m := range models {
		result[i] = *toNotifDomain(m)
	}
	return result, nil
}

func (r *pgNotificationRepo) UnreadCount(ctx context.Context, companyID, userID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&notifModel{}).
		Where("company_id = ? AND user_id = ? AND read_at IS NULL", companyID, userID).
		Count(&count).Error
	return int(count), err
}

func (r *pgNotificationRepo) MarkRead(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&notifModel{}).Where("id = ?", id).Update("read_at", gorm.Expr("NOW()")).Error
}

func (r *pgNotificationRepo) MarkAllRead(ctx context.Context, companyID, userID string) error {
	return r.db.WithContext(ctx).Model(&notifModel{}).
		Where("company_id = ? AND user_id = ? AND read_at IS NULL", companyID, userID).
		Update("read_at", gorm.Expr("NOW()")).Error
}

func (r *pgNotificationRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&notifModel{}).Error
}
