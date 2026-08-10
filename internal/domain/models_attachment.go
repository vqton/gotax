package domain

import (
	"context"
	"time"
)

type Attachment struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	MimeType   string    `json:"mime_type"`
	FileSize   int64     `json:"file_size"`
	UploadedBy string    `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type AttachmentRepository interface {
	Create(ctx context.Context, a *Attachment) error
	GetByID(ctx context.Context, id string) (*Attachment, error)
	ListByEntity(ctx context.Context, entityType, entityID string) ([]Attachment, error)
	Delete(ctx context.Context, id string) error
}
