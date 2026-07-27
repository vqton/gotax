package service

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"

	"gotax/internal/repository"
)

func setupImportService(t *testing.T) (Service, context.Context) {
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

	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo,
		approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo)
	return svc, context.Background()
}

func makeExcel(t *testing.T, rows [][]string) []byte {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()
	for i, row := range rows {
		for j, cell := range row {
			col, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", col, i+1), cell)
		}
	}
	var buf bytes.Buffer
	_, err := f.WriteTo(&buf)
	require.NoError(t, err)
	return buf.Bytes()
}

func TestImportOpeningBalances_Success(t *testing.T) {
	svc, ctx := setupImportService(t)
	data := makeExcel(t, [][]string{
		{"account_code", "debit_amount", "credit_amount"},
		{"1111", "100000000", ""},
		{"2111", "", "100000000"},
	})

	result, err := svc.ImportOpeningBalances(ctx, data, "C001", "P-2025-01", "U001")
	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 2, result.Success)
	assert.Len(t, result.BalanceIDs, 2)
	assert.Empty(t, result.Errors)
}

func TestImportOpeningBalances_WithVietnameseHeaders(t *testing.T) {
	svc, ctx := setupImportService(t)
	data := makeExcel(t, [][]string{
		{"ma tk", "no", "co"},
		{"1111", "50000000", ""},
		{"2111", "", "50000000"},
	})

	result, err := svc.ImportOpeningBalances(ctx, data, "C001", "P-2025-01", "U001")
	require.NoError(t, err)
	assert.Equal(t, 2, result.Success)
}

func TestImportOpeningBalances_PartialErrors(t *testing.T) {
	svc, ctx := setupImportService(t)
	data := makeExcel(t, [][]string{
		{"account_code", "debit_amount", "credit_amount"},
		{"1111", "100000000", ""},
		{"", "50000", ""},
		{"2111", "", "50000000"},
	})

	result, err := svc.ImportOpeningBalances(ctx, data, "C001", "P-2025-01", "U001")
	require.NoError(t, err)
	assert.Equal(t, 3, result.Total)
	assert.Equal(t, 2, result.Success)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0].Error, "required")
}

func TestImportOpeningBalances_EmptyFile(t *testing.T) {
	svc, ctx := setupImportService(t)
	data := makeExcel(t, [][]string{
		{"account_code", "debit_amount"},
	})

	result, err := svc.ImportOpeningBalances(ctx, data, "C001", "P-2025-01", "U001")
	require.NoError(t, err)
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 0, result.Success)
}

func TestImportOpeningBalances_InvalidExcel(t *testing.T) {
	svc, ctx := setupImportService(t)
	_, err := svc.ImportOpeningBalances(ctx, []byte("not an excel file"), "C001", "P-2025-01", "U001")
	assert.Error(t, err)
}

func TestImportOpeningBalances_WithAllColumns(t *testing.T) {
	svc, ctx := setupImportService(t)
	data := makeExcel(t, [][]string{
		{"account_code", "currency_code", "original_amount", "debit_amount", "credit_amount", "exchange_rate", "source_type", "reason"},
		{"1111", "USD", "2000", "2000", "", "25400", "MANUAL", "import USD cash"},
	})

	result, err := svc.ImportOpeningBalances(ctx, data, "C001", "P-2025-01", "U001")
	require.NoError(t, err)
	assert.Equal(t, 1, result.Success)
}

func TestImportOpeningBalances_BothAmountsZero(t *testing.T) {
	svc, ctx := setupImportService(t)
	data := makeExcel(t, [][]string{
		{"account_code", "debit_amount"},
		{"1111", "0"},
	})

	result, err := svc.ImportOpeningBalances(ctx, data, "C001", "P-2025-01", "U001")
	require.NoError(t, err)
	assert.Equal(t, 0, result.Success)
	assert.Len(t, result.Errors, 1)
}
