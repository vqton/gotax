package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gotax/internal/domain"
)

type MemoryInvoiceBookRepo struct {
	mu      sync.RWMutex
	books   map[string]*domain.InvoiceBook
	numbers map[string]*domain.InvoiceNumber
}

func NewMemoryInvoiceBookRepo() *MemoryInvoiceBookRepo {
	return &MemoryInvoiceBookRepo{
		books:   make(map[string]*domain.InvoiceBook),
		numbers: make(map[string]*domain.InvoiceNumber),
	}
}

func (r *MemoryInvoiceBookRepo) CreateBook(_ context.Context, book *domain.InvoiceBook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, b := range r.books {
		if b.CompanyID == book.CompanyID && b.Pattern == book.Pattern && b.Serial == book.Serial {
			return domain.ErrInvoiceBookExists
		}
	}
	r.books[book.ID] = book
	return nil
}

func (r *MemoryInvoiceBookRepo) GetBook(_ context.Context, id string) (*domain.InvoiceBook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	book, ok := r.books[id]
	if !ok {
		return nil, domain.ErrInvoiceBookNotFound
	}
	out := *book
	return &out, nil
}

func (r *MemoryInvoiceBookRepo) ListBooks(_ context.Context, companyID string) ([]domain.InvoiceBook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.InvoiceBook
	for _, b := range r.books {
		if b.CompanyID == companyID {
			out = append(out, *b)
		}
	}
	return out, nil
}

func (r *MemoryInvoiceBookRepo) UpdateBook(_ context.Context, book *domain.InvoiceBook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[book.ID]; !ok {
		return domain.ErrInvoiceBookNotFound
	}
	book.UpdatedAt = time.Now()
	r.books[book.ID] = book
	return nil
}

func (r *MemoryInvoiceBookRepo) DeleteBook(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; !ok {
		return domain.ErrInvoiceBookNotFound
	}
	delete(r.books, id)
	return nil
}

func (r *MemoryInvoiceBookRepo) AllocateNumber(_ context.Context, bookID string) (*domain.InvoiceNumber, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	book, ok := r.books[bookID]
	if !ok {
		return nil, domain.ErrInvoiceBookNotFound
	}
	if book.Status != "active" {
		return nil, domain.ErrInvoiceBookNotActive
	}
	if book.NextNumber > book.ToNumber {
		return nil, domain.ErrInvoiceBookFull
	}
	now := time.Now()
	num := &domain.InvoiceNumber{
		ID:        fmt.Sprintf("IN-%d", now.UnixNano()),
		BookID:    bookID,
		Number:    book.NextNumber,
		Status:    domain.InvNumIssued,
		IssuedAt:  &now,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.numbers[num.ID] = num
	book.NextNumber++
	book.UsedCount++
	return num, nil
}

func (r *MemoryInvoiceBookRepo) ReleaseNumber(_ context.Context, numberID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	num, ok := r.numbers[numberID]
	if !ok {
		return domain.ErrInvoiceNumberNotFound
	}
	if num.Status != domain.InvNumIssued {
		return domain.ErrInvoiceNumberNotAvailable
	}
	num.Status = domain.InvNumAvailable
	num.InvoiceID = ""
	num.IssuedAt = nil
	num.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryInvoiceBookRepo) GetNumbersByBook(_ context.Context, bookID string) ([]domain.InvoiceNumber, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.InvoiceNumber
	for _, n := range r.numbers {
		if n.BookID == bookID {
			out = append(out, *n)
		}
	}
	return out, nil
}

func (r *MemoryInvoiceBookRepo) GetNumberByID(_ context.Context, id string) (*domain.InvoiceNumber, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	num, ok := r.numbers[id]
	if !ok {
		return nil, domain.ErrInvoiceNumberNotFound
	}
	out := *num
	return &out, nil
}

func (r *MemoryInvoiceBookRepo) MarkMissing(_ context.Context, numberID string, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	num, ok := r.numbers[numberID]
	if !ok {
		return domain.ErrInvoiceNumberNotFound
	}
	num.Status = domain.InvNumMissing
	num.MissingReason = reason
	num.UpdatedAt = time.Now()
	return nil
}
