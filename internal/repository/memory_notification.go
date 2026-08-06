package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"gotax/internal/domain"
)

type MemoryNotificationRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.Notification
}

func NewMemoryNotificationRepo() domain.NotificationRepository {
	return &MemoryNotificationRepo{data: make(map[string]*domain.Notification)}
}

func (r *MemoryNotificationRepo) Create(_ context.Context, n *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *n
	r.data[n.ID] = &cp
	return nil
}

func (r *MemoryNotificationRepo) GetByID(_ context.Context, id string) (*domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.data[id]
	if !ok {
		return nil, domain.ErrNotificationNotFound
	}
	cp := *n
	return &cp, nil
}

func (r *MemoryNotificationRepo) ListByUser(_ context.Context, companyID, userID string, limit int) ([]domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Notification
	for _, n := range r.data {
		if n.CompanyID == companyID && n.UserID == userID {
			result = append(result, *n)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (r *MemoryNotificationRepo) UnreadCount(_ context.Context, companyID, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, n := range r.data {
		if n.CompanyID == companyID && n.UserID == userID && n.ReadAt == "" {
			count++
		}
	}
	return count, nil
}

func (r *MemoryNotificationRepo) MarkRead(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.data[id]
	if !ok {
		return domain.ErrNotificationNotFound
	}
	now := time.Now().UTC().Format(time.RFC3339)
	n.ReadAt = now
	return nil
}

func (r *MemoryNotificationRepo) MarkAllRead(_ context.Context, companyID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC().Format(time.RFC3339)
	for _, n := range r.data {
		if n.CompanyID == companyID && n.UserID == userID && n.ReadAt == "" {
			n.ReadAt = now
		}
	}
	return nil
}

func (r *MemoryNotificationRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrNotificationNotFound
	}
	delete(r.data, id)
	return nil
}
