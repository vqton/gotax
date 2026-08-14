package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

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

func seedReceipt(t *testing.T, cashRepo interface {
	CreateReceipt(context.Context, *domain.CashReceipt) error
}) *domain.CashReceipt {
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
	r, _, _, _, cashRepo, _, _ := setupSvc(t)
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
	r, _, _, _, cashRepo, _, _ := setupSvc(t)

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
	r, _, _, _, cashRepo, _, _ := setupSvc(t)

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
	r, _, _, _, cashRepo, _, _ := setupSvc(t)
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
	r, _, _, _, cashRepo, _, _ := setupSvc(t)
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

// ─── Cash Payments ─────────────────────────────────────────────────

func seedPayment(t *testing.T, cashRepo interface {
	CreatePayment(context.Context, *domain.CashPayment) error
}) *domain.CashPayment {
	t.Helper()
	p := &domain.CashPayment{
		ID:              "CP-SEED-1",
		CompanyID:       "CMP001",
		VoucherNo:       "P-2026-0001",
		VoucherDate:     "2026-08-03",
		CashAccountID:   "1111",
		PayeeName:       "Công ty ABC",
		PayeeType:       domain.CounterpartSupplier,
		Currency:        "VND",
		ExchangeRate:    1,
		Amount:          3000000,
		AmountVND:       3000000,
		DebitAccountID:  "3311",
		CreditAccountID: "1111",
		Reason:          "Thanh toán tiền hàng cho nhà cung cấp",
		PaymentType:     domain.PaymentSupplier,
		Status:          domain.CashDraft,
	}
	require.NoError(t, cashRepo.CreatePayment(context.Background(), p))
	return p
}

func TestCashPaymentsPageRender(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)
	seedPayment(t, cashRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/cash-payments.html?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Phiếu chi")
	assert.Contains(t, body, "P-2026-0001")
	assert.Contains(t, body, "3.000.000")
	assert.NotContains(t, body, "x-data")
}

func TestCashPaymentsCreateAction(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)

	form := url.Values{
		"voucher_date":     {"2026-08-04"},
		"payment_type":     {"SUPPLIER"},
		"reason":           {"Chi thanh toán NCC"},
		"payee_name":       {"Công ty XYZ"},
		"amount":           {"1500000"},
		"cash_account_id":  {"1111"},
		"debit_account_id": {"3311"},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/cash-payments/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "P-2026-0001")
	assert.Contains(t, body, "1.500.000")
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")

	payments, _, err := cashRepo.ListPayments(context.Background(), domain.CashPaymentFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	require.Len(t, payments, 1)
	assert.Equal(t, "P-2026-0001", payments[0].VoucherNo)
	assert.Equal(t, domain.CashDraft, payments[0].Status)
	assert.NotEmpty(t, payments[0].ID)
	// credit mirrors debit (legacy behavior) so Validate passes.
	assert.Equal(t, "3311", payments[0].CreditAccountID)
}

func TestCashPaymentsCreateValidationError(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)

	form := url.Values{
		"voucher_date":     {"2026-08-04"},
		"reason":           {"Chi tiền"},
		"amount":           {"0"}, // invalid: must be > 0
		"cash_account_id":  {"1111"},
		"debit_account_id": {"3311"},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/cash-payments/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
	payments, _, err := cashRepo.ListPayments(context.Background(), domain.CashPaymentFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	assert.Empty(t, payments)
}

func TestCashPaymentsLifecycleActions(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)
	seedPayment(t, cashRepo)
	// Fund the cash account: posting a payment requires sufficient balance.
	fund := &domain.CashReceipt{
		ID:              "CR-SEED-FUND",
		CompanyID:       "CMP001",
		VoucherNo:       "R-2026-0002",
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
		Reason:          "Quỹ đầu kỳ",
		ReceiptType:     domain.ReceiptCustomerPayment,
		Status:          domain.CashPosted,
	}
	require.NoError(t, cashRepo.CreateReceipt(context.Background(), fund))

	post := func(action string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/app/cash-payments/"+action+"?company_id=CMP001", strings.NewReader("id=CP-SEED-1"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)
		return w
	}

	// DRAFT → submit
	w := post("submit")
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	payment, err := cashRepo.GetPayment(context.Background(), "CP-SEED-1")
	require.NoError(t, err)
	assert.Equal(t, domain.CashSubmitted, payment.Status)

	// SUBMITTED → approve
	w = post("approve")
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	payment, err = cashRepo.GetPayment(context.Background(), "CP-SEED-1")
	require.NoError(t, err)
	assert.Equal(t, domain.CashApproved, payment.Status)

	// APPROVED → post (creates GL journal entry)
	w = post("post")
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	payment, err = cashRepo.GetPayment(context.Background(), "CP-SEED-1")
	require.NoError(t, err)
	assert.Equal(t, domain.CashPosted, payment.Status)
	assert.NotEmpty(t, payment.PostedAt)
	assert.NotEmpty(t, payment.GLJournalID)
}

func TestCashPaymentsPostFromDraftRejected(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)
	seedPayment(t, cashRepo)

	// Posting straight from DRAFT is an invalid transition — error toast.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/cash-payments/post?company_id=CMP001", strings.NewReader("id=CP-SEED-1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")

	payment, err := cashRepo.GetPayment(context.Background(), "CP-SEED-1")
	require.NoError(t, err)
	assert.Equal(t, domain.CashDraft, payment.Status)
}

// ─── Cash Transfers ───────────────────────────────────────────────

func TestCashTransfersPageRender(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)
	seedTransfer(t, cashRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/cash-transfers.html?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Chuyển quỹ")
	assert.Contains(t, body, "TRF-SEED-1")
	assert.Contains(t, body, "10.000.000")
	assert.NotContains(t, body, "x-data")
}

func seedTransfer(t *testing.T, cashRepo interface {
	CreateTransfer(context.Context, *domain.CashTransfer) error
}) *domain.CashTransfer {
	t.Helper()
	tf := &domain.CashTransfer{
		ID:            "TRF-SEED-1",
		CompanyID:     "CMP001",
		TransferDate:  "2026-08-05",
		FromAccountID: "1121",
		ToAccountID:   "1111",
		Amount:        10000000,
		Currency:      "VND",
		ExchangeRate:  1,
		Reason:        "Rút tiền NH về quỹ",
		TransferType:  domain.TransferBankWithdrawal,
		Status:        domain.CashPosted,
	}
	require.NoError(t, cashRepo.CreateTransfer(context.Background(), tf))
	return tf
}

func TestCashTransfersCreateAction(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)

	form := url.Values{
		"transfer_date":   {"2026-08-06"},
		"from_account_id": {"1121"},
		"to_account_id":   {"1111"},
		"amount":          {"5000000"},
		"reason":          {"Chuyển quỹ tiền mặt"},
		"transfer_type":   {"BANK_WITHDRAWAL"},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/cash-transfers/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "5.000.000")
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")

	transfers, err := cashRepo.ListTransfers(context.Background(), "CMP001")
	require.NoError(t, err)
	require.Len(t, transfers, 1)
	assert.Equal(t, domain.CashPosted, transfers[0].Status)

	// Transfer creates receipt + payment vouchers (both posted, same journal).
	receipts, _, err := cashRepo.ListReceipts(context.Background(), domain.CashReceiptFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	require.Len(t, receipts, 1)
	assert.Equal(t, "1111", receipts[0].CashAccountID)
	assert.Equal(t, domain.CashPosted, receipts[0].Status)
	assert.NotEmpty(t, receipts[0].VoucherNo)
	assert.NotEmpty(t, receipts[0].GLJournalID)

	payments, _, err := cashRepo.ListPayments(context.Background(), domain.CashPaymentFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	require.Len(t, payments, 1)
	assert.Equal(t, "1121", payments[0].CashAccountID)
	assert.Equal(t, domain.CashPosted, payments[0].Status)
	assert.NotEmpty(t, payments[0].VoucherNo)
	assert.Equal(t, receipts[0].GLJournalID, payments[0].GLJournalID)
}

func TestCashTransfersCreateValidationError(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)

	form := url.Values{
		"transfer_date":   {"2026-08-06"},
		"from_account_id": {"1121"},
		"to_account_id":   {"1111"},
		"amount":          {"0"}, // invalid: must be > 0
		"reason":          {"Chuyển quỹ"},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/cash-transfers/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
	transfers, err := cashRepo.ListTransfers(context.Background(), "CMP001")
	require.NoError(t, err)
	assert.Empty(t, transfers)
}

// ─── Cash Book Report ─────────────────────────────────────────────

func TestCashBookPageRenderAndFilter(t *testing.T) {
	r, _, _, _, cashRepo, _, _ := setupSvc(t)

	// Posted docs within the current month (default filter period).
	now := time.Now()
	in := now.Format("2006-01-02")
	rcpt := &domain.CashReceipt{
		ID: "CR-SEED-BOOK", CompanyID: "CMP001", VoucherNo: "R-2026-0099",
		VoucherDate: in, CashAccountID: "1111", Currency: "VND", ExchangeRate: 1,
		Amount: 5000000, AmountVND: 5000000, DebitAccountID: "1111", CreditAccountID: "5111",
		Reason: "Thu tiền bán hàng", Status: domain.CashPosted,
	}
	require.NoError(t, cashRepo.CreateReceipt(context.Background(), rcpt))
	pay := &domain.CashPayment{
		ID: "CP-SEED-BOOK", CompanyID: "CMP001", VoucherNo: "P-2026-0099",
		VoucherDate: in, CashAccountID: "1111", Currency: "VND", ExchangeRate: 1,
		Amount: 2000000, AmountVND: 2000000, DebitAccountID: "3311", CreditAccountID: "1111",
		Reason: "Chi thanh toán NCC", Status: domain.CashPosted,
	}
	require.NoError(t, cashRepo.CreatePayment(context.Background(), pay))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/cash-book.html?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Sổ quỹ tiền mặt")
	assert.Contains(t, body, "5.000.000") // total receipts
	assert.Contains(t, body, "2.000.000") // total payments
	assert.Contains(t, body, "3.000.000") // closing balance
	assert.Contains(t, body, "R-2026-0099")
	assert.NotContains(t, body, "x-data")

	// Filter action re-renders the report fragment with explicit dates.
	params := url.Values{
		"company_id": {"CMP001"}, "account_id": {"1111"},
		"currency": {"VND"}, "from_date": {"2026-08-01"}, "to_date": {"2026-08-31"},
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/app/cash-book/filter?"+params.Encode(), strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body = w.Body.String()
	assert.Contains(t, body, "Số dư đầu kỳ")
	assert.Contains(t, body, "5.000.000")
}

func TestCashFlowPageRenderAndFilter(t *testing.T) {
	r, _, perRepo, _, _, jeRepo, _ := setupSvc(t)

	// Period + one posted entry touching cash (operating inflow).
	now := time.Now()
	perRepo.Create(context.Background(), &domain.Period{
		ID: "P-202608", Year: now.Year(), Month: int(now.Month()),
		StartDate: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(now.Year(), now.Month(), 28, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	})
	require.NoError(t, jeRepo.Create(context.Background(), &domain.JournalEntry{
		ID: "je-cashflow", CompanyID: "CMP001", VoucherType: domain.VoucherTypeOther,
		EntryDate: time.Date(now.Year(), now.Month(), 5, 0, 0, 0, 0, time.UTC),
		PeriodID:  "P-202608", Status: domain.JournalEntryPosted,
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 5000000},
			{LineNumber: 2, AccountCode: "5111", CreditAmount: 5000000},
		},
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/cash-flow.html?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Lưu chuyển tiền tệ")
	assert.Contains(t, body, "BẢNG LƯU CHUYỂN TIỀN TỆ")
	assert.Contains(t, body, "I. Lưu chuyển tiền từ hoạt động kinh doanh")
	assert.Contains(t, body, "5.000.000") // operating inflow
	assert.Contains(t, body, "Tăng (giảm) ròng tiền và tương đương tiền")
	assert.NotContains(t, body, "x-data")

	// Filter action re-renders with explicit year/month.
	params := url.Values{"year": {"2026"}, "month": {"8"}}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/app/cash-flow/filter?"+params.Encode(), strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body = w.Body.String()
	assert.Contains(t, body, "Kỳ: Tháng 8/2026")
	assert.Contains(t, body, "5.000.000")
}
