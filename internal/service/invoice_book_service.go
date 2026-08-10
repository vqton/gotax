package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

func invGenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type InvoiceBookService struct {
	repo domain.InvoiceBookRepository
}

func NewInvoiceBookService(repo domain.InvoiceBookRepository) *InvoiceBookService {
	return &InvoiceBookService{repo: repo}
}

func (s *InvoiceBookService) CreateBook(ctx context.Context, book *domain.InvoiceBook) error {
	if book.Name == "" {
		return fmt.Errorf("book name is required")
	}
	if book.Pattern == "" {
		return fmt.Errorf("pattern is required")
	}
	if book.FromNumber >= book.ToNumber {
		return domain.ErrInvoiceBookRangeInvalid
	}
	if book.ID == "" {
		book.ID = invGenerateID()
	}
	book.NextNumber = book.FromNumber
	book.UsedCount = 0
	if book.Status == "" {
		book.Status = "active"
	}
	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now
	return s.repo.CreateBook(ctx, book)
}

func (s *InvoiceBookService) GetBook(ctx context.Context, id string) (*domain.InvoiceBook, error) {
	return s.repo.GetBook(ctx, id)
}

func (s *InvoiceBookService) ListBooks(ctx context.Context, companyID string) ([]domain.InvoiceBook, error) {
	return s.repo.ListBooks(ctx, companyID)
}

func (s *InvoiceBookService) UpdateBook(ctx context.Context, book *domain.InvoiceBook) error {
	existing, err := s.repo.GetBook(ctx, book.ID)
	if err != nil {
		return err
	}
	if existing.Status == "closed" {
		return fmt.Errorf("cannot update closed book")
	}
	book.UpdatedAt = time.Now()
	return s.repo.UpdateBook(ctx, book)
}

func (s *InvoiceBookService) DeleteBook(ctx context.Context, id string) error {
	book, err := s.repo.GetBook(ctx, id)
	if err != nil {
		return err
	}
	if book.UsedCount > 0 {
		return fmt.Errorf("cannot delete book with issued numbers")
	}
	return s.repo.DeleteBook(ctx, id)
}

func (s *InvoiceBookService) AllocateNumber(ctx context.Context, bookID string) (*domain.InvoiceNumber, error) {
	return s.repo.AllocateNumber(ctx, bookID)
}

func (s *InvoiceBookService) ReleaseNumber(ctx context.Context, numberID string) error {
	return s.repo.ReleaseNumber(ctx, numberID)
}

func (s *InvoiceBookService) GetNumbersByBook(ctx context.Context, bookID string) ([]domain.InvoiceNumber, error) {
	return s.repo.GetNumbersByBook(ctx, bookID)
}

func (s *InvoiceBookService) MarkMissing(ctx context.Context, numberID string, reason string) error {
	if reason == "" {
		return fmt.Errorf("missing reason is required")
	}
	return s.repo.MarkMissing(ctx, numberID, reason)
}

func (s *InvoiceBookService) GetAvailableCount(ctx context.Context, bookID string) (int, error) {
	book, err := s.repo.GetBook(ctx, bookID)
	if err != nil {
		return 0, err
	}
	return book.ToNumber - book.NextNumber + 1, nil
}
