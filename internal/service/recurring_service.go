package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type RecurringService struct {
	repo   domain.RecurringEntryRepository
	jeRepo domain.JournalRepository
}

func NewRecurringService(repo domain.RecurringEntryRepository, jeRepo domain.JournalRepository) *RecurringService {
	return &RecurringService{repo: repo, jeRepo: jeRepo}
}

func (s *RecurringService) Create(ctx context.Context, entry *domain.RecurringEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("REC-%d", time.Now().UnixNano())
	}
	if entry.NextRunDate == "" {
		entry.NextRunDate = time.Now().Format("2006-01-02")
	}
	return s.repo.Create(ctx, entry)
}

func (s *RecurringService) GetByID(ctx context.Context, id string) (*domain.RecurringEntry, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RecurringService) List(ctx context.Context, companyID string) ([]domain.RecurringEntry, error) {
	return s.repo.List(ctx, companyID)
}

func (s *RecurringService) Update(ctx context.Context, entry *domain.RecurringEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	entry.UpdatedAt = time.Now()
	return s.repo.Update(ctx, entry)
}

func (s *RecurringService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *RecurringService) RunNow(ctx context.Context, id string, userID string) (*domain.JournalEntry, error) {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !entry.IsActive {
		return nil, domain.ErrRecurringNotActive
	}

	now := time.Now()
	je := &domain.JournalEntry{
		CompanyID:   entry.CompanyID,
		EntryDate:   now,
		Description: fmt.Sprintf("RECURRING:%s", entry.TemplateName),
		Status:      domain.JournalEntryDraft,
		CreatedBy:   userID,
		Lines:       make([]domain.JournalLine, 0, len(entry.Lines)),
	}
	for i, rl := range entry.Lines {
		je.Lines = append(je.Lines, domain.JournalLine{
			AccountCode:  rl.AccountCode,
			DebitAmount:  rl.DebitAmount,
			CreditAmount: rl.CreditAmount,
			Description:  rl.Description,
			LineNumber:   i + 1,
		})
	}
	if err := s.jeRepo.Create(ctx, je); err != nil {
		return nil, err
	}

	next := advanceDate(entry.Frequency, entry.DayOfMonth, now)
	_ = s.repo.UpdateNextRunDate(ctx, id, next.Format("2006-01-02"))

	return je, nil
}

func (s *RecurringService) ProcessDue(ctx context.Context, userID string) (int, error) {
	today := time.Now().Format("2006-01-02")
	due, err := s.repo.GetDueEntries(ctx, today)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range due {
		if _, err := s.RunNow(ctx, entry.ID, userID); err != nil {
			continue
		}
		count++
	}
	return count, nil
}

func advanceDate(freq domain.RecurringFrequency, dayOfMonth int, from time.Time) time.Time {
	switch freq {
	case domain.RecurringMonthly:
		return time.Date(from.Year(), from.Month()+1, dayOfMonth, 0, 0, 0, 0, time.UTC)
	case domain.RecurringQuarterly:
		return time.Date(from.Year(), from.Month()+3, dayOfMonth, 0, 0, 0, 0, time.UTC)
	case domain.RecurringYearly:
		return time.Date(from.Year()+1, from.Month(), dayOfMonth, 0, 0, 0, 0, time.UTC)
	default:
		return from.AddDate(0, 1, 0)
	}
}
