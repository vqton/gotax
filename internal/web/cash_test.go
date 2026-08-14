package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
)

// TestToastHeaderASCII guards against header mojibake: browsers decode
// response headers as ISO-8859-1, so HX-Trigger must be pure ASCII JSON with
// non-ASCII escaped as \uXXXX. The text must JSON.parse back to proper UTF-8.
func TestToastHeaderASCII(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	Toast(c, "success", "Đã tạo phiếu thu R-2026-0001.")

	header := c.Writer.Header().Get("HX-Trigger")
	require.NotEmpty(t, header)
	for _, b := range []byte(header) {
		if b > 127 {
			t.Fatalf("HX-Trigger header contains non-ASCII byte %#x: %s", b, header)
		}
	}
	var payload struct {
		Toast struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"toast"`
	}
	require.NoError(t, json.Unmarshal([]byte(header), &payload))
	assert.Equal(t, "success", payload.Toast.Type)
	assert.Equal(t, "Đã tạo phiếu thu R-2026-0001.", payload.Toast.Text)
}

func seedReceipt(t *testing.T, cashRepo interface{ CreateReceipt(context.Context, *domain.CashReceipt) error }) *domain.CashReceipt {
	t.Helper()
	r := &domain.CashReceipt{
		ID:              "CR-SEED-1",
		CompanyID:       "CMP001",
		VoucherNo:       "R-2026-0001",
		VoucherDate:     "2026-08-01",
		CashAccountID:   "1111",
		CounterpartName: "Công ty ABC",
		CounterpartType: domain.CounterpartCustomer,
		Currency:        "VND",
		ExchangeRate:    1,
		Amount:          5000000,
		AmountVND:       5000000,
		DebitAccountID:  "1111",
		CreditAccountID: "1311",
		Reason:          "Khách hàng thanh toán nợ",
		ReceiptType:     domain.ReceiptCustomerPayment,
		Status:          domain.CashDraft,
	}
	require.NoError(t, cashRepo.CreateReceipt(context.Background(), r))
	return r
}

func TestCashReceiptsPageRender(t *testing.T) {
	r, _, _, _, cashRepo := setupSvc(t)
	seedReceipt(t, cashRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/cash-receipts.html?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Phiếu thu")
	assert.Contains(t, body, "R-2026-0001")
	assert.Contains(t, body, "5.000.000")
	assert.NotContains(t, body, "x-data")
}

func TestCashReceiptsCreateAction(t *testing.T) {
	r, _, _, _, cashRepo := setupSvc(t)

	form := url.Values{
		"voucher_date":     {"2026-08-02"},
		"receipt_type":     {"CUSTOMER_PAYMENT"},
		"reason":           {"Thu tiền bán hàng"},
		"counterpart_name": {"Công ty XYZ"},
		"amount":           {"1500000"},
		"cash_account_id":  {"1111"},
		"debit_account_id": {"5111"},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/cash-receipts/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "R-2026-0001")
	assert.Contains(t, body, "1.500.000")
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")

	receipts, _, err := cashRepo.ListReceipts(context.Background(), domain.CashReceiptFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	require.Len(t, receipts, 1)
	assert.Equal(t, "R-2026-0001", receipts[0].VoucherNo)
	assert.Equal(t, domain.CashDraft, receipts[0].Status)
	assert.NotEmpty(t, receipts[0].ID)
	// credit mirrors debit (legacy behavior) so Validate passes.
	assert.Equal(t, "5111", receipts[0].CreditAccountID)
}

func TestCashReceiptsCreateValidationError(t *testing.T) {
	r, _, _, _, cashRepo := setupSvc(t)

	form := url.Values{
		"voucher_date":     {"2026-08-02"},
		"reason":           {"Thu tiền"},
		"amount":           {"0"}, // invalid: must be > 0
		"cash_account_id":  {"1111"},
		"debit_account_id": {"5111"},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/cash-receipts/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
	receipts, _, err := cashRepo.ListReceipts(context.Background(), domain.CashReceiptFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	assert.Empty(t, receipts)
}

func TestCashReceiptsLifecycleActions(t *testing.T) {
	r, _, _, _, cashRepo := setupSvc(t)
	seedReceipt(t, cashRepo)

	post := func(action string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/app/cash-receipts/"+action+"?company_id=CMP001", strings.NewReader("id=CR-SEED-1"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)
		return w
	}

	// DRAFT → submit
	w := post("submit")
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	receipt, err := cashRepo.GetReceipt(context.Background(), "CR-SEED-1")
	require.NoError(t, err)
	assert.Equal(t, domain.CashSubmitted, receipt.Status)

	// SUBMITTED → approve
	w = post("approve")
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	receipt, err = cashRepo.GetReceipt(context.Background(), "CR-SEED-1")
	require.NoError(t, err)
	assert.Equal(t, domain.CashApproved, receipt.Status)

	// APPROVED → post (creates GL journal entry)
	w = post("post")
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	receipt, err = cashRepo.GetReceipt(context.Background(), "CR-SEED-1")
	require.NoError(t, err)
	assert.Equal(t, domain.CashPosted, receipt.Status)
	assert.NotEmpty(t, receipt.PostedAt)
	assert.NotEmpty(t, receipt.GLJournalID)
}

func TestCashReceiptsPostFromDraftRejected(t *testing.T) {
	r, _, _, _, cashRepo := setupSvc(t)
	seedReceipt(t, cashRepo)

	// Posting straight from DRAFT is an invalid transition — error toast.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/cash-receipts/post?company_id=CMP001", strings.NewReader("id=CR-SEED-1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")

	receipt, err := cashRepo.GetReceipt(context.Background(), "CR-SEED-1")
	require.NoError(t, err)
	assert.Equal(t, domain.CashDraft, receipt.Status)
}
