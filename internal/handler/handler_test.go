package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/auth"
	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupTest(t *testing.T) (*gin.Engine, service.Service, context.Context) {
	t.Helper()
	gin.SetMode(gin.TestMode)

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

	svc := service.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo,
		approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo)

	h := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterRoutes(r, h, noopMW, noopMW)

	return r, svc, context.Background()
}

// ─── Auth ───────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	r, svc, ctx := setupTest(t)
	auth.SetJWTSecret("test-secret-handler")
	svc.CreateUser(ctx, &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}, "ValidPass123!")

	body := `{"username":"admin","password":"ValidPass123!"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var result domain.AuthResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result.AccessToken)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	r, svc, ctx := setupTest(t)
	auth.SetJWTSecret("test-secret-handler-2")
	svc.CreateUser(ctx, &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}, "ValidPass123!")

	body := `{"username":"admin","password":"wrong"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestLogin_BadRequest(t *testing.T) {
	r, _, _ := setupTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`not-json`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestRefreshToken(t *testing.T) {
	r, svc, ctx := setupTest(t)
	auth.SetJWTSecret("test-secret-refresh")
	svc.CreateUser(ctx, &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}, "ValidPass123!")

	// login first to get refresh token
	loginBody := `{"username":"admin","password":"ValidPass123!"}`
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(loginBody))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	require.Equal(t, 200, w1.Code)
	var loginResult domain.AuthResult
	json.Unmarshal(w1.Body.Bytes(), &loginResult)

	body := `{"refresh_token":"` + loginResult.RefreshToken + `"}`
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)
}

// ─── Accounts (auth'd) ──────────────────────────────────────────────

func TestCreateAccount(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"code":"1111","name":"Cash","type":"ASSET"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var acct domain.Account
	json.Unmarshal(w.Body.Bytes(), &acct)
	assert.Equal(t, "1111", acct.Code)
}

func TestCreateAccount_Duplicate(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"code":"1111","name":"Cash","type":"ASSET"}`
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/accounts", strings.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	require.Equal(t, 201, w1.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/accounts", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 400, w2.Code)
}

func TestCreateAccount_ValidationError(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"code":"","name":"","type":"ASSET"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListAccounts(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, IsActive: true})
	svc.CreateAccount(ctx, &domain.Account{Code: "2111", Name: "Payables", Type: domain.AccountTypeLiability, IsActive: true})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accounts", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var accounts []domain.Account
	json.Unmarshal(w.Body.Bytes(), &accounts)
	assert.Len(t, accounts, 2)
}

func TestGetAccount(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accounts/1111", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var acct domain.Account
	json.Unmarshal(w.Body.Bytes(), &acct)
	assert.Equal(t, "Cash", acct.Name)
}

func TestGetAccount_NotFound(t *testing.T) {
	r, _, _ := setupTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/accounts/9999", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestUpdateAccount(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset})

	body := `{"name":"Cash Updated","type":"ASSET"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/accounts/1111", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var acct domain.Account
	json.Unmarshal(w.Body.Bytes(), &acct)
	assert.Equal(t, "Cash Updated", acct.Name)
}

func TestDeleteAccount(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/accounts/1111", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Journal Entries ────────────────────────────────────────────────

func TestCreateEntry(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true})

	body := `{"entry_date":"2025-06-15T00:00:00Z","description":"test","lines":[{"line_number":1,"account_code":"1111","debit_amount":100,"credit_amount":0},{"line_number":2,"account_code":"5111","debit_amount":0,"credit_amount":100}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/journal-entries", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestCreateEntry_ValidationError(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"entry_date":"","description":"","lines":[]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/journal-entries", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

// ─── Periods ────────────────────────────────────────────────────────

func TestCreatePeriod(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"id":"P-2025-06","year":2025,"month":6,"start_date":"2025-06-01T00:00:00Z","end_date":"2025-06-30T00:00:00Z","status":"OPEN"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/periods", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestListPeriods(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/periods", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var periods []domain.Period
	json.Unmarshal(w.Body.Bytes(), &periods)
	assert.Len(t, periods, 1)
}

// ─── Exchange Rates ─────────────────────────────────────────────────

func TestCreateExchangeRate(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"currency_code":"USD","average_rate":23000,"rate_date":"2025-06-15T00:00:00Z"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/exchange-rates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

// ─── Users ──────────────────────────────────────────────────────────

func TestCreateUser(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"user":{"username":"newuser","full_name":"New User","role":"accountant"},"password":"ValidPass123!"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

// ─── Reports ────────────────────────────────────────────────────────

func TestReports_RequireParams(t *testing.T) {
	r, _, _ := setupTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/trial-balance", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/reports/balance-sheet?year=2025&month=0", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 400, w2.Code)
}

func TestTrialBalance(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	})
	je := &domain.JournalEntry{
		EntryDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		PeriodID:  "P-2025-06", Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 500, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 500},
		},
	}
	svc.CreateEntry(ctx, je, "u")
	svc.SubmitForReview(ctx, je.ID, "u")
	svc.ApproveEntry(ctx, je.ID, "a")
	svc.PostEntry(ctx, je.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/trial-balance?year=2025&month=6", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── COA Operations ─────────────────────────────────────────────────

func TestFreezeAccount(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset})

	body := `{"reason":"audit"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/coa/accounts/1111/freeze", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestFreezeAccount_MissingReason(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/coa/accounts/1111/freeze", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAccountBalance(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	})
	je := &domain.JournalEntry{
		EntryDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		PeriodID:  "P-2025-06", Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 500, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 500},
		},
	}
	svc.CreateEntry(ctx, je, "u")
	svc.SubmitForReview(ctx, je.ID, "u")
	svc.ApproveEntry(ctx, je.ID, "a")
	svc.PostEntry(ctx, je.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/coa/accounts/1111/balance?period_id=P-2025-06", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestAccountBalance_MissingPeriod(t *testing.T) {
	r, _, _ := setupTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/coa/accounts/1111/balance", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

// ─── Audit ──────────────────────────────────────────────────────────

func TestAuditLogByEntity_RequiresParams(t *testing.T) {
	r, _, _ := setupTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/audit/entity", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestAuditLog(t *testing.T) {
	r, _, _ := setupTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/audit?limit=10", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── COA: Approval ──────────────────────────────────────────────────

func TestCreateApprovalRequest(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"entity_type":"account","entity_id":"1111","request_type":"change","reason":"restructure"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/coa/approvals", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

// ─── COA: Version ──────────────────────────────────────────────────

func TestCreateAccountVersion(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset})

	body := `{"reason":"snapshot"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/coa/versions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

// ─── COA: Mapping ───────────────────────────────────────────────────

func TestCreateAccountMapping(t *testing.T) {
	r, _, _ := setupTest(t)
	body := `{"source_regime":"VAS","target_regime":"IFRS","old_code":"1111","new_code":"IFRS-1111","mapping_type":"DIRECT"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/coa/mappings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

// ─── TOTP ───────────────────────────────────────────────────────────

func TestSetupTOTP(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateUser(ctx, &domain.User{ID: "test-user", Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}, "ValidPass123!")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/totp/setup", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// ─── GetJournalEntries ──────────────────────────────────────────────

func TestGetJournalEntries_RequiresParams(t *testing.T) {
	r, _, _ := setupTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/journal-entries", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetJournalEntries_ByStatus(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreateEntry(ctx, &domain.JournalEntry{
		EntryDate: time.Now(), Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}, "u")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/journal-entries?status=DRAFT", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetJournalEntries_ByDateRange(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true})
	svc.CreateEntry(ctx, &domain.JournalEntry{
		EntryDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC), Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}, "u")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/journal-entries?from=2025-01-01&to=2025-12-31", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetJournalEntries_InvalidDate(t *testing.T) {
	r, _, _ := setupTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/journal-entries?from=bad&to=bad", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

// ─── GetCurrentUser ─────────────────────────────────────────────────

func TestGetCurrentUser(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateUser(ctx, &domain.User{ID: "test-user", Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}, "ValidPass123!")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var u domain.User
	json.Unmarshal(w.Body.Bytes(), &u)
	assert.Equal(t, "admin", u.Username)
}

// ─── Sessions ───────────────────────────────────────────────────────

func TestListSessions(t *testing.T) {
	r, svc, ctx := setupTest(t)
	auth.SetJWTSecret("test-secret-sessions-h")
	svc.CreateUser(ctx, &domain.User{ID: "test-user", Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}, "ValidPass123!")
	svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/auth/sessions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// ─── ChangePassword ─────────────────────────────────────────────────

func TestChangePassword(t *testing.T) {
	r, svc, ctx := setupTest(t)
	svc.CreateUser(ctx, &domain.User{ID: "test-user", Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}, "ValidPass123!")

	body := `{"current_password":"ValidPass123!","new_password":"NewValidPass456!"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/change-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
