package repository

import (
	"context"

	"gotax/internal/domain"
	"gorm.io/gorm"
)

type PGBankImportRepo struct {
	db *gorm.DB
}

func NewPGBankImportRepo(db *gorm.DB) *PGBankImportRepo {
	return &PGBankImportRepo{db: db}
}

func (r *PGBankImportRepo) CreateImport(ctx context.Context, imp *domain.BankImport) error {
	return r.db.WithContext(ctx).Create(imp).Error
}

func (r *PGBankImportRepo) GetImport(ctx context.Context, id string) (*domain.BankImport, error) {
	var imp domain.BankImport
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&imp).Error; err != nil {
		return nil, err
	}
	return &imp, nil
}

func (r *PGBankImportRepo) ListImports(ctx context.Context, companyID string) ([]domain.BankImport, error) {
	var items []domain.BankImport
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PGBankImportRepo) UpdateImportStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.BankImport{}).Where("id = ?", id).Update("status", status).Error
}

func (r *PGBankImportRepo) CreateTransactions(ctx context.Context, txns []domain.BankTransaction) error {
	if len(txns) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&txns).Error
}

func (r *PGBankImportRepo) GetTransactionsByImport(ctx context.Context, importID string) ([]domain.BankTransaction, error) {
	var txns []domain.BankTransaction
	if err := r.db.WithContext(ctx).Where("reference = ?", importID).Find(&txns).Error; err != nil {
		return nil, err
	}
	return txns, nil
}
