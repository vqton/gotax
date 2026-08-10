package domain

import (
	"context"
	"time"
)

// AuditChange captures field-level changes for detailed audit trail.
type AuditChange struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	FieldName  string    `json:"field_name"`
	OldValue   string    `json:"old_value"`
	NewValue   string    `json:"new_value"`
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuditChangeRepository interface {
	Create(ctx context.Context, c *AuditChange) error
	GetByEntity(ctx context.Context, entityType, entityID string) ([]AuditChange, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]AuditChange, error)
}
