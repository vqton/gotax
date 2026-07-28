package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupBankTest(t *testing.T) (*gin.Engine, *service.BankService, context.Context) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	bankRepo := repository.NewMemoryBankRepo()
	bankSvc := service.NewBankService(bankRepo)
	bankH := NewBankHandler(bankSvc, nil)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterBankRoutes(r, bankH, noopMW)

	return r, bankSvc, context.Background()
}

// ─── Statements ────────────────────────────────────────────────────────

func TestImportStatement(t *testing.T) {
	r, _, _ := setupBankTest(t)

	body := `{"statement":{"company_id":"CMP001","bank_account_id":"BA001","statement_date":"2026-07-01","from_date":"2026-07-01","to_date":"2026-07-31","opening_balance":1000000,"closing_balance":2000000,"currency":"VND","import_method":"MANUAL"},"lines":[{"transaction_date":"2026-07-01","description":"test txn","debit_amount":0,"credit_amount":1000000}]}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/statements/import?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var stmt domain.BankStatement
	err := json.Unmarshal(w.Body.Bytes(), &stmt)
	require.NoError(t, err)
	assert.NotEmpty(t, stmt.ID)
}

func TestGetStatement(t *testing.T) {
	r, svc, ctx := setupBankTest(t)
	stmt := &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-07", OpeningBalance: 1000000, ClosingBalance: 2000000, Currency: "VND", Status: domain.BankStatementImported}
	require.NoError(t, svc.ImportStatement(ctx, stmt, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bank/statements/"+stmt.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestListStatements(t *testing.T) {
	r, svc, ctx := setupBankTest(t)
	svc.ImportStatement(ctx, &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-07", OpeningBalance: 1000000, ClosingBalance: 2000000, Currency: "VND"}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bank/statements?company_id=CMP001&bank_account_id=BA001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp struct {
		Data  []domain.BankStatement `json:"data"`
		Total int                    `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1, resp.Total)
}

func TestDeleteStatement(t *testing.T) {
	r, svc, ctx := setupBankTest(t)
	stmt := &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-07", OpeningBalance: 1000000, ClosingBalance: 2000000, Currency: "VND", Status: domain.BankStatementImported}
	require.NoError(t, svc.ImportStatement(ctx, stmt, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/bank/statements/"+stmt.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestGetStatementLines(t *testing.T) {
	r, svc, ctx := setupBankTest(t)
	stmt := &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-07", OpeningBalance: 1000000, ClosingBalance: 2000000, Currency: "VND"}
	lines := []domain.BankStatementLine{
		{TransactionDate: "2026-07-01", Description: "txn1", DebitAmount: 0, CreditAmount: 500000},
	}
	require.NoError(t, svc.ImportStatement(ctx, stmt, lines))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bank/statements/"+stmt.ID+"/lines", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var got []domain.BankStatementLine
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, 1, len(got))
}

// ─── Payment Orders ────────────────────────────────────────────────────

func TestCreatePaymentOrder(t *testing.T) {
	r, _, _ := setupBankTest(t)

	body := `{"company_id":"CMP001","payment_date":"2026-07-15","amount":5000000,"currency":"VND","beneficiary_name":"ABC Corp","beneficiary_acc":"123456789","beneficiary_bank":"VCB","from_bank_acc_id":"BA001","payment_content":"Invoice payment"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/payment-orders?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var po domain.PaymentOrder
	json.Unmarshal(w.Body.Bytes(), &po)
	assert.NotEmpty(t, po.ID)
	assert.Equal(t, domain.PODraft, po.Status)
}

func TestPaymentOrderWorkflow(t *testing.T) {
	r, svc, ctx := setupBankTest(t)
	po := &domain.PaymentOrder{CompanyID: "CMP001", PaymentDate: "2026-07-15", Amount: 5000000, Currency: "VND", BeneficiaryName: "ABC Corp", BeneficiaryAcc: "123456789", BeneficiaryBank: "VCB", FromBankAccID: "BA001", PaymentContent: "test", CreatedBy: "user1"}
	require.NoError(t, svc.CreatePaymentOrder(ctx, po))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/payment-orders/"+po.ID+"/submit", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/bank/payment-orders/"+po.ID+"/approve", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/bank/payment-orders/"+po.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var got domain.PaymentOrder
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, domain.POApproved, got.Status)
}

func TestRejectPaymentOrder(t *testing.T) {
	r, svc, ctx := setupBankTest(t)
	po := &domain.PaymentOrder{CompanyID: "CMP001", PaymentDate: "2026-07-15", Amount: 5000000, Currency: "VND", BeneficiaryName: "ABC Corp", BeneficiaryAcc: "123456789", BeneficiaryBank: "VCB", FromBankAccID: "BA001", PaymentContent: "test", CreatedBy: "user1"}
	require.NoError(t, svc.CreatePaymentOrder(ctx, po))
	require.NoError(t, svc.SubmitPaymentOrder(ctx, po.ID))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/payment-orders/"+po.ID+"/reject", strings.NewReader(`{"reason":"insufficient funds"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Loans ─────────────────────────────────────────────────────────────

func TestCreateLoan(t *testing.T) {
	r, _, _ := setupBankTest(t)

	body := `{"company_id":"CMP001","bank_account_id":"BA001","contract_no":"LN-001","loan_type":"SHORT_TERM","principal_amount":100000000,"currency":"VND","interest_rate":8.5,"interest_method":"FIXED","start_date":"2026-01-01","maturity_date":"2026-12-31","repayment_method":"ANNUITY","repayment_freq":"MONTHLY"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/loans?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var l domain.LoanAgreement
	json.Unmarshal(w.Body.Bytes(), &l)
	assert.NotEmpty(t, l.ID)
	assert.Equal(t, domain.LoanActive, l.Status)
}

func TestDisburseLoan(t *testing.T) {
	r, svc, ctx := setupBankTest(t)
	l := &domain.LoanAgreement{CompanyID: "CMP001", BankAccountID: "BA001", ContractNo: "LN-001", PrincipalAmount: 100000000, Currency: "VND", InterestRate: 8.5, StartDate: "2026-01-01", MaturityDate: "2026-12-31"}
	require.NoError(t, svc.CreateLoan(ctx, l))

	body := `{"amount":50000000,"disbursement_date":"2026-01-15"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/loans/"+l.ID+"/disburse", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

// ─── Term Deposits ─────────────────────────────────────────────────────

func TestCreateDeposit(t *testing.T) {
	r, _, _ := setupBankTest(t)

	body := `{"company_id":"CMP001","bank_account_id":"BA001","deposit_no":"TD-001","amount":50000000,"currency":"VND","interest_rate":6.5,"term_days":90,"start_date":"2026-07-01","maturity_date":"2026-09-28"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/term-deposits?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var d domain.TermDeposit
	json.Unmarshal(w.Body.Bytes(), &d)
	assert.NotEmpty(t, d.ID)
}

func TestMatureDeposit(t *testing.T) {
	r, svc, ctx := setupBankTest(t)
	d := &domain.TermDeposit{CompanyID: "CMP001", BankAccountID: "BA001", DepositNo: "TD-001", Amount: 50000000, Currency: "VND", InterestRate: 6.5, TermDays: 90, StartDate: "2026-07-01", MaturityDate: "2026-09-28", Status: domain.DepositActive}
	require.NoError(t, svc.CreateDeposit(ctx, d))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/term-deposits/"+d.ID+"/mature", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// ─── Reconciliation ────────────────────────────────────────────────────

func TestStartReconciliation(t *testing.T) {
	r, _, _ := setupBankTest(t)

	body := `{"company_id":"CMP001","bank_account_id":"BA001","statement_id":"STMT001","from_date":"2026-07-01","to_date":"2026-07-31","opening_balance":1000000,"closing_balance":2000000,"statement_balance":2000000}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/reconciliations?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

// ─── Batches ───────────────────────────────────────────────────────────

func TestCreateBatch(t *testing.T) {
	r, _, _ := setupBankTest(t)

	body := `{"batch":{"company_id":"CMP001","batch_name":"BATCH-001","batch_date":"2026-07-15","total_amount":10000000,"currency":"VND","order_count":2},"order_ids":["PO1","PO2"]}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/bank/batches?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

// ─── Reports ─────────────────────────────────────────────────────────

func TestGetBankLedger(t *testing.T) {
	r, _, _ := setupBankTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/bank-ledger?company_id=CMP001&bank_account_id=BA001&from_date=2026-01-01&to_date=2026-12-31", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGetBankBalance(t *testing.T) {
	r, _, _ := setupBankTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/bank-balance?company_id=CMP001&bank_account_id=BA001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
