package repository

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
	"gorm.io/gorm"
)

func pgGenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type PGInvoiceBookRepo struct {
	db *gorm.DB
}

func NewPGInvoiceBookRepo(db *gorm.DB) *PGInvoiceBookRepo {
	return &PGInvoiceBookRepo{db: db}
}

func (r *PGInvoiceBookRepo) CreateBook(ctx context.Context, book *domain.InvoiceBook) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *PGInvoiceBookRepo) GetBook(ctx context.Context, id string) (*domain.InvoiceBook, error) {
	var book domain.InvoiceBook
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&book).Error; err != nil {
		return nil, domain.ErrInvoiceBookNotFound
	}
	return &book, nil
}

func (r *PGInvoiceBookRepo) ListBooks(ctx context.Context, companyID string) ([]domain.InvoiceBook, error) {
	var books []domain.InvoiceBook
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *PGInvoiceBookRepo) UpdateBook(ctx context.Context, book *domain.InvoiceBook) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *PGInvoiceBookRepo) DeleteBook(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.InvoiceBook{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return domain.ErrInvoiceBookNotFound
	}
	return result.Error
}

func (r *PGInvoiceBookRepo) AllocateNumber(ctx context.Context, bookID string) (*domain.InvoiceNumber, error) {
	var book domain.InvoiceBook
	if err := r.db.WithContext(ctx).Where("id = ?", bookID).First(&book).Error; err != nil {
		return nil, domain.ErrInvoiceBookNotFound
	}
	if book.Status != "active" {
		return nil, domain.ErrInvoiceBookNotActive
	}
	if book.NextNumber > book.ToNumber {
		return nil, domain.ErrInvoiceBookFull
	}
	num := &domain.InvoiceNumber{
		ID:     pgGenerateID(),
		BookID: bookID,
		Number: book.NextNumber,
		Status: domain.InvNumIssued,
	}
	if err := r.db.WithContext(ctx).Create(num).Error; err != nil {
		return nil, err
	}
	book.NextNumber++
	book.UsedCount++
	if err := r.db.WithContext(ctx).Save(&book).Error; err != nil {
		return nil, err
	}
	return num, nil
}

func (r *PGInvoiceBookRepo) ReleaseNumber(ctx context.Context, numberID string) error {
	var num domain.InvoiceNumber
	if err := r.db.WithContext(ctx).Where("id = ?", numberID).First(&num).Error; err != nil {
		return domain.ErrInvoiceNumberNotFound
	}
	if num.Status != domain.InvNumIssued {
		return domain.ErrInvoiceNumberNotAvailable
	}
	num.Status = domain.InvNumAvailable
	num.InvoiceID = ""
	num.IssuedAt = nil
	return r.db.WithContext(ctx).Save(&num).Error
}

func (r *PGInvoiceBookRepo) GetNumbersByBook(ctx context.Context, bookID string) ([]domain.InvoiceNumber, error) {
	var nums []domain.InvoiceNumber
	if err := r.db.WithContext(ctx).Where("book_id = ?", bookID).Order("number ASC").Find(&nums).Error; err != nil {
		return nil, err
	}
	return nums, nil
}

func (r *PGInvoiceBookRepo) GetNumberByID(ctx context.Context, id string) (*domain.InvoiceNumber, error) {
	var num domain.InvoiceNumber
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&num).Error; err != nil {
		return nil, domain.ErrInvoiceNumberNotFound
	}
	return &num, nil
}

func (r *PGInvoiceBookRepo) MarkMissing(ctx context.Context, numberID string, reason string) error {
	var num domain.InvoiceNumber
	if err := r.db.WithContext(ctx).Where("id = ?", numberID).First(&num).Error; err != nil {
		return domain.ErrInvoiceNumberNotFound
	}
	num.Status = domain.InvNumMissing
	num.MissingReason = reason
	return r.db.WithContext(ctx).Save(&num).Error
}
