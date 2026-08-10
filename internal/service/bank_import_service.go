package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
	"gotax/internal/einvoice"
)

type BankImportService struct {
	repo       domain.BankImportRepository
	parsers    map[string]domain.BankCSVParser
}

func NewBankImportService(repo domain.BankImportRepository) *BankImportService {
	svc := &BankImportService{
		repo:    repo,
		parsers: make(map[string]domain.BankCSVParser),
	}
	svc.registerParser(einvoice.NewVCBParser())
	svc.registerParser(einvoice.NewBIDVParser())
	svc.registerParser(einvoice.NewCTGParser())
	svc.registerParser(einvoice.NewVTBParser())
	svc.registerParser(einvoice.NewACBParser())
	return svc
}

func (s *BankImportService) registerParser(parser domain.BankCSVParser) {
	s.parsers[parser.GetBankCode()] = parser
}

func (s *BankImportService) Import(ctx context.Context, companyID, bankCode, fileName string, data []byte) (*domain.BankImport, error) {
	parser, ok := s.parsers[bankCode]
	if !ok {
		return nil, fmt.Errorf("unsupported bank code: %s", bankCode)
	}
	if err := parser.Validate(data); err != nil {
		return nil, fmt.Errorf("invalid CSV data: %w", err)
	}
	txns, err := parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}
	imp := &domain.BankImport{
		ID:               fmt.Sprintf("IMP-%d", time.Now().UnixNano()),
		CompanyID:        companyID,
		BankCode:         bankCode,
		FileName:         fileName,
		Status:           "pending",
		TransactionCount: len(txns),
		ImportedAt:       time.Now(),
		CreatedAt:        time.Now(),
	}
	if err := s.repo.CreateImport(ctx, imp); err != nil {
		return nil, fmt.Errorf("failed to create import: %w", err)
	}
	for i := range txns {
		txns[i].Reference = imp.ID
	}
	if err := s.repo.CreateTransactions(ctx, txns); err != nil {
		return nil, fmt.Errorf("failed to create transactions: %w", err)
	}
	if err := s.repo.UpdateImportStatus(ctx, imp.ID, "completed"); err != nil {
		return nil, fmt.Errorf("failed to update import status: %w", err)
	}
	imp.Status = "completed"
	return imp, nil
}

func (s *BankImportService) GetImport(ctx context.Context, id string) (*domain.BankImport, error) {
	return s.repo.GetImport(ctx, id)
}

func (s *BankImportService) ListImports(ctx context.Context, companyID string) ([]domain.BankImport, error) {
	return s.repo.ListImports(ctx, companyID)
}

func (s *BankImportService) GetTransactions(ctx context.Context, importID string) ([]domain.BankTransaction, error) {
	return s.repo.GetTransactionsByImport(ctx, importID)
}
