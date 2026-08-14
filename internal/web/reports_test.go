package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func seedReportData(t *testing.T, perRepo *repository.MemoryPeriodRepo, jeRepo *repository.MemoryJournalRepo) {
	t.Helper()
	now := time.Now()
	perRepo.Create(context.Background(), &domain.Period{
		ID: "P-202608", Year: now.Year(), Month: int(now.Month()),
		StartDate: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(now.Year(), now.Month(), 28, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	})
	require.NoError(t, jeRepo.Create(context.Background(), &domain.JournalEntry{
		ID: "je-report", CompanyID: "CMP001", EntryNumber: "202608-001",
		VoucherType: domain.VoucherTypeReceipt,
		EntryDate:   time.Date(now.Year(), now.Month(), 5, 0, 0, 0, 0, time.UTC),
		PeriodID:    "P-202608", Status: domain.JournalEntryPosted,
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 5000000},
			{LineNumber: 2, AccountCode: "5111", CreditAmount: 5000000},
		},
	}))
}

func TestJournalExportPageRender(t *testing.T) {
	r, _, _, _, _, _, _ := setupSvc(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/journal-export.html?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Nhật ký chung")
	assert.Contains(t, body, "Xuất file Excel")
	assert.Contains(t, body, `id="je-year"`)
	assert.Contains(t, body, `id="je-month"`)
	assert.Contains(t, body, `/app/journal-export/download`)
	assert.NotContains(t, body, "x-data")
}

func TestJournalExportDownloadAction(t *testing.T) {
	r, _, perRepo, _, _, jeRepo, _ := setupSvc(t)
	seedReportData(t, perRepo, jeRepo)

	params := strings.NewReader("year=2026&month=8")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/journal-export/download", params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "spreadsheetml")
	assert.Contains(t, w.Header().Get("Content-Disposition"), "chung-tu-2026-08.xlsx")
	// Valid xlsx is a zip archive.
	assert.Equal(t, "PK", w.Body.String()[:2])
}

func TestTrialBalancePageRenderAndFilter(t *testing.T) {
	r, _, perRepo, _, _, jeRepo, accRepo := setupSvc(t)
	seedReportData(t, perRepo, jeRepo)
	accRepo.Create(context.Background(), &domain.Account{Code: "1111", Name: "Tiền mặt", Type: domain.AccountTypeAsset})
	accRepo.Create(context.Background(), &domain.Account{Code: "5111", Name: "Doanh thu bán hàng", Type: domain.AccountTypeRevenue})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/trial-balance.html?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Bảng cân đối phát sinh")
	assert.Contains(t, body, "1111")
	assert.Contains(t, body, "5111")
	assert.Contains(t, body, "Tiền mặt")
	assert.Contains(t, body, "Doanh thu bán hàng")
	assert.Contains(t, body, "5.000.000")
	assert.Contains(t, body, "Tổng cộng")
	assert.NotContains(t, body, "x-data")

	// Filter action re-renders the report fragment.
	params := strings.NewReader("year=2026&month=8")
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/app/trial-balance/filter", params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body = w.Body.String()
	assert.Contains(t, body, "Tháng 8/2026")
	assert.Contains(t, body, "5.000.000")
}
