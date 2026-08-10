package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupPrintService(t *testing.T) (Service, context.Context) {
	t.Helper()
	accRepo := repository.NewMemoryAccountRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	jeRepo.SetAccounts(accRepo.Accounts())
	perRepo := repository.NewMemoryPeriodRepo()
	userRepo := repository.NewMemoryUserRepo()
	auditRepo := repository.NewMemoryAuditLogRepo()
	rateRepo := repository.NewMemoryExchangeRateRepo()
	templateRepo := repository.NewMemoryClosingTemplateRepo()
	approvalRepo := repository.NewMemoryApprovalRepo()
	versionRepo := repository.NewMemoryAccountVersionRepo()
	mappingRepo := repository.NewMemoryAccountMappingRepo()
	analysisRepo := repository.NewMemoryAccountAnalysisRepo()
	ifrsRepo := repository.NewMemoryIFRSMappingRepo()
	refreshRepo := repository.NewMemoryRefreshTokenRepo()
	resetRepo := repository.NewMemoryPasswordResetTokenRepo()
	obRepo := repository.NewMemoryOpeningBalanceRepo()
	cashRepo := repository.NewMemoryCashRepo()

	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo,
		approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)
	return svc, context.Background()
}

func TestGeneratePrintPDF_Receipt(t *testing.T) {
	svc, ctx := setupPrintService(t)

	// Create a cash receipt
	receipt := &domain.CashReceipt{
		CompanyID:       "C001",
		VoucherNo:       "PT-001",
		VoucherDate:     "2026-08-10",
		CashAccountID:   "1111",
		CounterpartName: "Nguyen Van A",
		CounterpartType: domain.CounterpartCustomer,
		Currency:        "VND",
		ExchangeRate:    1,
		Amount:          5_000_000,
		AmountVND:       5_000_000,
		DebitAccountID:  "1111",
		CreditAccountID: "5111",
		Reason:          "Thu tien ban hang",
		ReceiptType:     domain.ReceiptSales,
		Status:          domain.CashPosted,
		CreatedBy:       "U001",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, receipt))

	// Generate PDF
	data, err := svc.GeneratePrintPDF(ctx, PrintTypeReceiptVoucher, receipt.ID)
	require.NoError(t, err)
	assert.True(t, len(data) > 100) // PDF has content
	assert.Contains(t, string(data[:20]), "%PDF") // Starts with PDF header
}

func TestGeneratePrintPDF_Payment(t *testing.T) {
	svc, ctx := setupPrintService(t)

	// Create a cash payment
	payment := &domain.CashPayment{
		CompanyID:     "C001",
		VoucherNo:     "PC-001",
		VoucherDate:   "2026-08-10",
		CashAccountID: "1111",
		PayeeName:     "Nguyen Van B",
		PayeeType:     domain.CounterpartSupplier,
		Currency:      "VND",
		ExchangeRate:  1,
		Amount:        3_000_000,
		AmountVND:     3_000_000,
		DebitAccountID:  "6321",
		CreditAccountID: "1111",
		Reason:        "Chi mua hang",
		PaymentType:   domain.PaymentSupplier,
		Status:        domain.CashPosted,
		CreatedBy:     "U001",
	}
	require.NoError(t, svc.CreateCashPayment(ctx, payment))

	// Generate PDF
	data, err := svc.GeneratePrintPDF(ctx, PrintTypePaymentVoucher, payment.ID)
	require.NoError(t, err)
	assert.True(t, len(data) > 100)
	assert.Contains(t, string(data[:20]), "%PDF")
}

func TestGeneratePrintPDF_InvalidType(t *testing.T) {
	svc, ctx := setupPrintService(t)

	_, err := svc.GeneratePrintPDF(ctx, "invalid", "some-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported print type")
}

func TestGeneratePrintPDF_NotFound(t *testing.T) {
	svc, ctx := setupPrintService(t)

	_, err := svc.GeneratePrintPDF(ctx, PrintTypeReceiptVoucher, "nonexistent")
	assert.Error(t, err)
}
