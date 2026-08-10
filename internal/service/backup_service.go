package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gotax/internal/domain"
)

type BackupService struct {
	repo      domain.BackupRepository
	databaseURL string
	backupDir   string
}

func NewBackupService(repo domain.BackupRepository, databaseURL, backupDir string) *BackupService {
	return &BackupService{repo: repo, databaseURL: databaseURL, backupDir: backupDir}
}

// CreateBackup performs a pg_dump and records the backup metadata.
func (s *BackupService) CreateBackup(ctx context.Context, companyID, createdBy string) (*domain.BackupRecord, error) {
	if companyID == "" {
		return nil, fmt.Errorf("company_id is required")
	}
	if s.databaseURL == "" {
		return nil, fmt.Errorf("database_url not configured")
	}

	if s.backupDir == "" {
		s.backupDir = "/tmp/gotax-backups"
	}
	os.MkdirAll(s.backupDir, 0755)

	filename := fmt.Sprintf("gotax_%s_%s.sql", companyID, time.Now().UTC().Format("20060102_150405"))
	filepath := filepath.Join(s.backupDir, filename)

	cmd := exec.CommandContext(ctx, "pg_dump", s.databaseURL, "-f", filepath, "--no-owner", "--no-privileges")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pg_dump failed: %w", err)
	}

	stat, _ := os.Stat(filepath)
	var fileSize int64
	if stat != nil {
		fileSize = stat.Size()
	}

	b := &domain.BackupRecord{
		CompanyID:  companyID,
		Filename:   filename,
		FileSize:   fileSize,
		BackupType: "manual",
		Status:     "completed",
		CreatedBy:  createdBy,
	}
	if err := s.repo.Create(ctx, b); err != nil {
		return nil, fmt.Errorf("record backup: %w", err)
	}
	return b, nil
}

func (s *BackupService) List(ctx context.Context, companyID string) ([]domain.BackupRecord, error) {
	return s.repo.List(ctx, companyID)
}

func (s *BackupService) GetByID(ctx context.Context, id string) (*domain.BackupRecord, error) {
	return s.repo.GetByID(ctx, id)
}
