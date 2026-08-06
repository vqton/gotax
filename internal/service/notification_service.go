package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type NotificationService struct {
	repo domain.NotificationRepository
}

func NewNotificationService(repo domain.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Create(ctx context.Context, n *domain.Notification) error {
	if n.CreatedAt == "" {
		n.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if n.Type == "" {
		n.Type = domain.NotifTypeINFO
	}
	if n.ID == "" {
		n.ID = fmt.Sprintf("n-%d", time.Now().UnixNano())
	}
	return s.repo.Create(ctx, n)
}

func (s *NotificationService) List(ctx context.Context, companyID, userID string, limit int) ([]domain.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.ListByUser(ctx, companyID, userID, limit)
}

func (s *NotificationService) UnreadCount(ctx context.Context, companyID, userID string) (int, error) {
	return s.repo.UnreadCount(ctx, companyID, userID)
}

func (s *NotificationService) MarkRead(ctx context.Context, id string) error {
	return s.repo.MarkRead(ctx, id)
}

func (s *NotificationService) MarkAllRead(ctx context.Context, companyID, userID string) error {
	return s.repo.MarkAllRead(ctx, companyID, userID)
}

func (s *NotificationService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
