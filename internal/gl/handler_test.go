package gl

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockService implements Service interface.
type mockService struct {
	createAccountFn         func(ctx context.Context, account *Account) error
	getAccountFn            func(ctx context.Context, code string) (*Account, error)
	getAllAccountsFn        func(ctx context.Context, activeOnly bool) ([]Account, error)
	updateAccountFn         func(ctx context.Context, account *Account) error
	deleteAccountFn         func(ctx context.Context, code string) error
	createEntryFn           func(ctx context.Context, entry *JournalEntry, userID string) error
	submitForReviewFn       func(ctx context.Context, id, userID string) error
	reviewEntryFn           func(ctx context.Context, id, reviewerID string) error
	approveEntryFn          func(ctx context.Context, id, approverID string) error
	postEntryFn             func(ctx context.Context, id string) error
	cancelEntryFn           func(ctx context.Context, id string) error
	getEntryByIDFn          func(ctx context.Context, id string) (*JournalEntry, error)
	getEntriesByDateRangeFn func(ctx context.Context, from, to time.Time) ([]JournalEntry, error)
	getEntriesByStatusFn    func(ctx context.Context, status JournalEntryStatus) ([]JournalEntry, error)
	trialBalanceFn          func(ctx context.Context, year, month int) ([]AccountBalance, error)
	balanceSheetFn          func(ctx context.Context, year, month int) ([]AccountBalance, error)
	incomeStatementFn       func(ctx context.Context, year, month int) ([]AccountBalance, error)
	createPeriodFn          func(ctx context.Context, period *Period) error
	getPeriodFn             func(ctx context.Context, id string) (*Period, error)
	getPeriodByYearMonthFn  func(ctx context.Context, year, month int) (*Period, error)
	getAllPeriodsFn         func(ctx context.Context) ([]Period, error)
	closePeriodFn           func(ctx context.Context, id string) error
	reopenPeriodFn          func(ctx context.Context, id string) error
	createExchangeRateFn    func(ctx context.Context, rate *ExchangeRate) error
	getExchangeRateFn       func(ctx context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error)
	listExchangeRatesFn     func(ctx context.Context) ([]ExchangeRate, error)
	getAuditLogFn           func(ctx context.Context, limit int) ([]AuditEntry, error)
	getAuditLogByEntityFn   func(ctx context.Context, entityType, entityID string) ([]AuditEntry, error)
	createUserFn            func(ctx context.Context, user *User, password string) error
	getUserFn               func(ctx context.Context, id string) (*User, error)
	listUsersFn             func(ctx context.Context) ([]User, error)
	authenticateFn          func(ctx context.Context, username, password string) (*User, error)
	// COA methods
	freezeAccountFn            func(ctx context.Context, code, reason string) error
	unfreezeAccountFn          func(ctx context.Context, code, reason string) error
	getAccountBalanceFn        func(ctx context.Context, code, periodID string) (*AccountBalance, error)
	getAccountBalanceDrillDownFn func(ctx context.Context, code, periodID string) ([]JournalEntry, error)
	getAccountUsageFn          func(ctx context.Context, code string) (*AccountUsage, error)
	createApprovalRequestFn    func(ctx context.Context, req *ApprovalRequest) error
	approveRequestFn           func(ctx context.Context, id, reviewerID, note string) error
	rejectRequestFn            func(ctx context.Context, id, reviewerID, note string) error
	getApprovalRequestsFn      func(ctx context.Context, status ApprovalStatus) ([]ApprovalRequest, error)
	createAccountVersionFn     func(ctx context.Context, reason string) (*AccountVersion, error)
	getVersionFn               func(ctx context.Context, versionNumber string) (*AccountVersion, error)
	listVersionsFn             func(ctx context.Context) ([]AccountVersion, error)
	compareVersionsFn          func(ctx context.Context, v1, v2 string) (*VersionDiff, error)
	createAccountAnalysisFn    func(ctx context.Context, analysis *AccountAnalysis) error
	getAccountAnalysisFn       func(ctx context.Context, accountCode string) (*AccountAnalysis, error)
	updateAccountAnalysisFn    func(ctx context.Context, analysis *AccountAnalysis) error
	createAccountMappingFn     func(ctx context.Context, mapping *AccountMapping) error
	getMappingByOldCodeFn      func(ctx context.Context, sourceRegime, oldCode string) (*AccountMapping, error)
	listMappingsFn             func(ctx context.Context, sourceRegime, targetRegime string) ([]AccountMapping, error)
	createIFRSMappingFn        func(ctx context.Context, mapping *IFRSMapping) error
	getIFRSMappingFn           func(ctx context.Context, vasCode string) (*IFRSMapping, error)
	listIFRSMappingsFn          func(ctx context.Context) ([]IFRSMapping, error)
	// Auth methods
	loginFn           func(ctx context.Context, username, password, ip string) (*AuthResult, error)
	verify2FAFn       func(ctx context.Context, tempToken, code string) (*AuthResult, error)
	refreshTokenFn    func(ctx context.Context, refreshTokenStr string) (*TokenPair, error)
	logoutFn          func(ctx context.Context, userID, refreshTokenStr string) error
	logoutAllFn       func(ctx context.Context, userID string) error
	changePasswordFn  func(ctx context.Context, userID, currentPassword, newPassword string) error
	forgotPasswordFn  func(ctx context.Context, email string) error
	resetPasswordFn   func(ctx context.Context, token, newPassword string) error
	setupTOTPFn       func(ctx context.Context, userID string) (*TOTPSetup, error)
	confirmTOTPFn     func(ctx context.Context, userID, code string) error
	disableTOTPFn     func(ctx context.Context, userID, currentPassword, code string) error
	generateBackupCodesFn func(ctx context.Context, userID string) ([]string, error)
	listSessionsFn    func(ctx context.Context, userID string) ([]Session, error)
	revokeSessionFn   func(ctx context.Context, userID, sessionID string) error
}

func (m *mockService) CreateAccount(ctx context.Context, account *Account) error {
	return m.createAccountFn(ctx, account)
}
func (m *mockService) GetAccount(ctx context.Context, code string) (*Account, error) {
	return m.getAccountFn(ctx, code)
}
func (m *mockService) GetAllAccounts(ctx context.Context, activeOnly bool) ([]Account, error) {
	return m.getAllAccountsFn(ctx, activeOnly)
}
func (m *mockService) UpdateAccount(ctx context.Context, account *Account) error {
	return m.updateAccountFn(ctx, account)
}
func (m *mockService) DeleteAccount(ctx context.Context, code string) error {
	return m.deleteAccountFn(ctx, code)
}
func (m *mockService) CreateEntry(ctx context.Context, entry *JournalEntry, userID string) error {
	return m.createEntryFn(ctx, entry, userID)
}
func (m *mockService) SubmitForReview(ctx context.Context, id, userID string) error {
	return m.submitForReviewFn(ctx, id, userID)
}
func (m *mockService) ReviewEntry(ctx context.Context, id, reviewerID string) error {
	return m.reviewEntryFn(ctx, id, reviewerID)
}
func (m *mockService) ApproveEntry(ctx context.Context, id, approverID string) error {
	return m.approveEntryFn(ctx, id, approverID)
}
func (m *mockService) PostEntry(ctx context.Context, id string) error { return m.postEntryFn(ctx, id) }
func (m *mockService) CancelEntry(ctx context.Context, id string) error {
	return m.cancelEntryFn(ctx, id)
}
func (m *mockService) GetEntryByID(ctx context.Context, id string) (*JournalEntry, error) {
	return m.getEntryByIDFn(ctx, id)
}
func (m *mockService) GetEntriesByDateRange(ctx context.Context, from, to time.Time) ([]JournalEntry, error) {
	return m.getEntriesByDateRangeFn(ctx, from, to)
}
func (m *mockService) GetEntriesByStatus(ctx context.Context, status JournalEntryStatus) ([]JournalEntry, error) {
	return m.getEntriesByStatusFn(ctx, status)
}
func (m *mockService) TrialBalance(ctx context.Context, year, month int) ([]AccountBalance, error) {
	return m.trialBalanceFn(ctx, year, month)
}
func (m *mockService) BalanceSheet(ctx context.Context, year, month int) ([]AccountBalance, error) {
	return m.balanceSheetFn(ctx, year, month)
}
func (m *mockService) IncomeStatement(ctx context.Context, year, month int) ([]AccountBalance, error) {
	return m.incomeStatementFn(ctx, year, month)
}
func (m *mockService) CreatePeriod(ctx context.Context, period *Period) error {
	return m.createPeriodFn(ctx, period)
}
func (m *mockService) GetPeriod(ctx context.Context, id string) (*Period, error) {
	return m.getPeriodFn(ctx, id)
}
func (m *mockService) GetPeriodByYearMonth(ctx context.Context, year, month int) (*Period, error) {
	return m.getPeriodByYearMonthFn(ctx, year, month)
}
func (m *mockService) GetAllPeriods(ctx context.Context) ([]Period, error) {
	return m.getAllPeriodsFn(ctx)
}
func (m *mockService) ClosePeriod(ctx context.Context, id string) error {
	return m.closePeriodFn(ctx, id)
}
func (m *mockService) ReopenPeriod(ctx context.Context, id string) error {
	return m.reopenPeriodFn(ctx, id)
}
func (m *mockService) CreateExchangeRate(ctx context.Context, rate *ExchangeRate) error {
	return m.createExchangeRateFn(ctx, rate)
}
func (m *mockService) GetExchangeRate(ctx context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error) {
	return m.getExchangeRateFn(ctx, currencyCode, rateDate)
}
func (m *mockService) ListExchangeRates(ctx context.Context) ([]ExchangeRate, error) {
	return m.listExchangeRatesFn(ctx)
}
func (m *mockService) GetAuditLog(ctx context.Context, limit int) ([]AuditEntry, error) {
	return m.getAuditLogFn(ctx, limit)
}
func (m *mockService) GetAuditLogByEntity(ctx context.Context, entityType, entityID string) ([]AuditEntry, error) {
	return m.getAuditLogByEntityFn(ctx, entityType, entityID)
}
func (m *mockService) CreateUser(ctx context.Context, user *User, password string) error {
	return m.createUserFn(ctx, user, password)
}
func (m *mockService) GetUser(ctx context.Context, id string) (*User, error) {
	return m.getUserFn(ctx, id)
}
func (m *mockService) ListUsers(ctx context.Context) ([]User, error) { return m.listUsersFn(ctx) }
func (m *mockService) Authenticate(ctx context.Context, username, password string) (*User, error) {
	return m.authenticateFn(ctx, username, password)
}

// ─── COA delegate methods ────────────────────────────────────────

func (m *mockService) FreezeAccount(ctx context.Context, code, reason string) error {
	return m.freezeAccountFn(ctx, code, reason)
}
func (m *mockService) UnfreezeAccount(ctx context.Context, code, reason string) error {
	return m.unfreezeAccountFn(ctx, code, reason)
}
func (m *mockService) GetAccountBalance(ctx context.Context, code, periodID string) (*AccountBalance, error) {
	return m.getAccountBalanceFn(ctx, code, periodID)
}
func (m *mockService) GetAccountBalanceDrillDown(ctx context.Context, code, periodID string) ([]JournalEntry, error) {
	return m.getAccountBalanceDrillDownFn(ctx, code, periodID)
}
func (m *mockService) GetAccountUsage(ctx context.Context, code string) (*AccountUsage, error) {
	return m.getAccountUsageFn(ctx, code)
}
func (m *mockService) CreateApprovalRequest(ctx context.Context, req *ApprovalRequest) error {
	return m.createApprovalRequestFn(ctx, req)
}
func (m *mockService) ApproveRequest(ctx context.Context, id, reviewerID, note string) error {
	return m.approveRequestFn(ctx, id, reviewerID, note)
}
func (m *mockService) RejectRequest(ctx context.Context, id, reviewerID, note string) error {
	return m.rejectRequestFn(ctx, id, reviewerID, note)
}
func (m *mockService) GetApprovalRequests(ctx context.Context, status ApprovalStatus) ([]ApprovalRequest, error) {
	return m.getApprovalRequestsFn(ctx, status)
}
func (m *mockService) CreateAccountVersion(ctx context.Context, reason string) (*AccountVersion, error) {
	return m.createAccountVersionFn(ctx, reason)
}
func (m *mockService) GetVersion(ctx context.Context, versionNumber string) (*AccountVersion, error) {
	return m.getVersionFn(ctx, versionNumber)
}
func (m *mockService) ListVersions(ctx context.Context) ([]AccountVersion, error) {
	return m.listVersionsFn(ctx)
}
func (m *mockService) CompareVersions(ctx context.Context, v1, v2 string) (*VersionDiff, error) {
	return m.compareVersionsFn(ctx, v1, v2)
}
func (m *mockService) CreateAccountAnalysis(ctx context.Context, analysis *AccountAnalysis) error {
	return m.createAccountAnalysisFn(ctx, analysis)
}
func (m *mockService) GetAccountAnalysis(ctx context.Context, accountCode string) (*AccountAnalysis, error) {
	return m.getAccountAnalysisFn(ctx, accountCode)
}
func (m *mockService) UpdateAccountAnalysis(ctx context.Context, analysis *AccountAnalysis) error {
	return m.updateAccountAnalysisFn(ctx, analysis)
}
func (m *mockService) CreateAccountMapping(ctx context.Context, mapping *AccountMapping) error {
	return m.createAccountMappingFn(ctx, mapping)
}
func (m *mockService) GetMappingByOldCode(ctx context.Context, sourceRegime, oldCode string) (*AccountMapping, error) {
	return m.getMappingByOldCodeFn(ctx, sourceRegime, oldCode)
}
func (m *mockService) ListMappings(ctx context.Context, sourceRegime, targetRegime string) ([]AccountMapping, error) {
	return m.listMappingsFn(ctx, sourceRegime, targetRegime)
}
func (m *mockService) CreateIFRSMapping(ctx context.Context, mapping *IFRSMapping) error {
	return m.createIFRSMappingFn(ctx, mapping)
}
func (m *mockService) GetIFRSMapping(ctx context.Context, vasCode string) (*IFRSMapping, error) {
	return m.getIFRSMappingFn(ctx, vasCode)
}
func (m *mockService) ListIFRSMappings(ctx context.Context) ([]IFRSMapping, error) {
	return m.listIFRSMappingsFn(ctx)
}

func (m *mockService) Login(ctx context.Context, username, password, ip string) (*AuthResult, error) {
	return m.loginFn(ctx, username, password, ip)
}
func (m *mockService) Verify2FA(ctx context.Context, tempToken, code string) (*AuthResult, error) {
	return m.verify2FAFn(ctx, tempToken, code)
}
func (m *mockService) RefreshToken(ctx context.Context, refreshTokenStr string) (*TokenPair, error) {
	return m.refreshTokenFn(ctx, refreshTokenStr)
}
func (m *mockService) Logout(ctx context.Context, userID, refreshTokenStr string) error {
	return m.logoutFn(ctx, userID, refreshTokenStr)
}
func (m *mockService) LogoutAll(ctx context.Context, userID string) error {
	return m.logoutAllFn(ctx, userID)
}
func (m *mockService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	return m.changePasswordFn(ctx, userID, currentPassword, newPassword)
}
func (m *mockService) ForgotPassword(ctx context.Context, email string) error {
	return m.forgotPasswordFn(ctx, email)
}
func (m *mockService) ResetPassword(ctx context.Context, token, newPassword string) error {
	return m.resetPasswordFn(ctx, token, newPassword)
}
func (m *mockService) SetupTOTP(ctx context.Context, userID string) (*TOTPSetup, error) {
	return m.setupTOTPFn(ctx, userID)
}
func (m *mockService) ConfirmTOTP(ctx context.Context, userID, code string) error {
	return m.confirmTOTPFn(ctx, userID, code)
}
func (m *mockService) DisableTOTP(ctx context.Context, userID, currentPassword, code string) error {
	return m.disableTOTPFn(ctx, userID, currentPassword, code)
}
func (m *mockService) GenerateBackupCodes(ctx context.Context, userID string) ([]string, error) {
	return m.generateBackupCodesFn(ctx, userID)
}
func (m *mockService) ListSessions(ctx context.Context, userID string) ([]Session, error) {
	return m.listSessionsFn(ctx, userID)
}
func (m *mockService) RevokeSession(ctx context.Context, userID, sessionID string) error {
	return m.revokeSessionFn(ctx, userID, sessionID)
}

func setupTestRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockAuth := func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
		c.Next()
	}
	mockAdmin := func(c *gin.Context) {
		c.Next()
	}
	RegisterRoutes(r, h, mockAuth, mockAdmin)
	return r
}

// ---------------------------------------------------------------------------
// Account tests
// ---------------------------------------------------------------------------

func TestCreateAccountHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create account", func(t *testing.T) {
		ms.createAccountFn = func(_ context.Context, acc *Account) error {
			assert.Equal(t, "1111", acc.Code)
			return nil
		}

		body := `{"code":"1111","name":"TM VND","type":"ASSET","is_active":true}`
		req := httptest.NewRequest("POST", "/api/v1/accounts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "account created", resp["message"])
	})

	t.Run("fail - bad request body", func(t *testing.T) {
		body := `{invalid json}`
		req := httptest.NewRequest("POST", "/api/v1/accounts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail - duplicate code", func(t *testing.T) {
		ms.createAccountFn = func(_ context.Context, _ *Account) error {
			return ErrAccountCodeExists
		}

		body := `{"code":"1111","name":"TM VND","type":"ASSET","is_active":true}`
		req := httptest.NewRequest("POST", "/api/v1/accounts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestGetAccountHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get account", func(t *testing.T) {
		ms.getAccountFn = func(_ context.Context, code string) (*Account, error) {
			return &Account{Code: "1111", Name: "TM VND", Type: AccountTypeAsset, IsActive: true}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/accounts/1111", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var acc Account
		json.Unmarshal(w.Body.Bytes(), &acc)
		assert.Equal(t, "1111", acc.Code)
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.getAccountFn = func(_ context.Context, code string) (*Account, error) {
			return nil, ErrAccountNotFound
		}

		req := httptest.NewRequest("GET", "/api/v1/accounts/9999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestListAccountsHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success list accounts", func(t *testing.T) {
		ms.getAllAccountsFn = func(_ context.Context, activeOnly bool) ([]Account, error) {
			return []Account{
				{Code: "1111", Name: "TM VND", Type: AccountTypeAsset, IsActive: true},
				{Code: "5111", Name: "DT BH", Type: AccountTypeRevenue, IsActive: true},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/accounts?active=true", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var accounts []Account
		json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.Len(t, accounts, 2)
	})
}

func TestUpdateAccountHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success update account", func(t *testing.T) {
		ms.updateAccountFn = func(_ context.Context, acc *Account) error {
			return nil
		}

		body := `{"code":"1111","name":"Updated Name","type":"ASSET","is_active":true}`
		req := httptest.NewRequest("PUT", "/api/v1/accounts/1111", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - code mismatch", func(t *testing.T) {
		body := `{"code":"1111","name":"Test","type":"ASSET","is_active":true}`
		req := httptest.NewRequest("PUT", "/api/v1/accounts/2222", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteAccountHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success delete account", func(t *testing.T) {
		ms.deleteAccountFn = func(_ context.Context, code string) error {
			return nil
		}

		req := httptest.NewRequest("DELETE", "/api/v1/accounts/1112", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.deleteAccountFn = func(_ context.Context, code string) error {
			return ErrAccountNotFound
		}

		req := httptest.NewRequest("DELETE", "/api/v1/accounts/9999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("fail - has children", func(t *testing.T) {
		ms.deleteAccountFn = func(_ context.Context, code string) error {
			return ErrAccountHasChildren
		}

		req := httptest.NewRequest("DELETE", "/api/v1/accounts/111", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

// ---------------------------------------------------------------------------
// Journal Entry tests
// ---------------------------------------------------------------------------

func TestCreateEntryHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create entry", func(t *testing.T) {
		ms.createEntryFn = func(_ context.Context, entry *JournalEntry, userID string) error {
			entry.ID = "JE-001"
			return nil
		}

		body := `{"entry_date":"2026-07-24T00:00:00Z","description":"Mua hang","lines":[{"account_code":"1561","debit_amount":10000000,"credit_amount":0},{"account_code":"1111","debit_amount":0,"credit_amount":10000000}]}`
		req := httptest.NewRequest("POST", "/api/v1/journal-entries", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "entry created", resp["message"])
	})

	t.Run("fail - bad request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/journal-entries", strings.NewReader(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail - account not found", func(t *testing.T) {
		ms.createEntryFn = func(_ context.Context, _ *JournalEntry, _ string) error {
			return ErrAccountNotFound
		}

		body := `{"entry_date":"2026-07-24T00:00:00Z","description":"Test","lines":[{"account_code":"9999","debit_amount":1000,"credit_amount":0}]}`
		req := httptest.NewRequest("POST", "/api/v1/journal-entries", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail - inactive account", func(t *testing.T) {
		ms.createEntryFn = func(_ context.Context, _ *JournalEntry, _ string) error {
			return ErrAccountInactive
		}

		body := `{"entry_date":"2026-07-24T00:00:00Z","description":"Test","lines":[{"account_code":"1111","debit_amount":1000,"credit_amount":0}]}`
		req := httptest.NewRequest("POST", "/api/v1/journal-entries", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSubmitEntryHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success submit entry", func(t *testing.T) {
		ms.submitForReviewFn = func(_ context.Context, id, userID string) error {
			return nil
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/submit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "entry submitted for review", resp["message"])
	})

	t.Run("fail - invalid state", func(t *testing.T) {
		ms.submitForReviewFn = func(_ context.Context, _, _ string) error {
			return errors.New("entry must be DRAFT to submit for review")
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/submit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestReviewEntryHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success review entry", func(t *testing.T) {
		ms.reviewEntryFn = func(_ context.Context, id, reviewerID string) error {
			return nil
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/review", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "entry reviewed", resp["message"])
	})

	t.Run("fail - already reviewed", func(t *testing.T) {
		ms.reviewEntryFn = func(_ context.Context, _, _ string) error {
			return ErrJournalAlreadyReviewed
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/review", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestApproveEntryHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success approve entry", func(t *testing.T) {
		ms.approveEntryFn = func(_ context.Context, id, approverID string) error {
			return nil
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "entry approved", resp["message"])
	})

	t.Run("fail - already approved", func(t *testing.T) {
		ms.approveEntryFn = func(_ context.Context, _, _ string) error {
			return ErrJournalAlreadyApproved
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPostJournalEntryHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success post entry", func(t *testing.T) {
		ms.postEntryFn = func(_ context.Context, id string) error {
			return nil
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/post", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "entry posted", resp["message"])
	})

	t.Run("fail - already posted", func(t *testing.T) {
		ms.postEntryFn = func(_ context.Context, id string) error {
			return ErrJournalAlreadyPosted
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/post", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("fail - already cancelled", func(t *testing.T) {
		ms.postEntryFn = func(_ context.Context, id string) error {
			return ErrJournalAlreadyCancelled
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/post", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("fail - period closed", func(t *testing.T) {
		ms.postEntryFn = func(_ context.Context, id string) error {
			return ErrJournalPeriodClosed
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/post", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail - period not found", func(t *testing.T) {
		ms.postEntryFn = func(_ context.Context, id string) error {
			return ErrPeriodNotFound
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/post", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCancelJournalEntryHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success cancel entry", func(t *testing.T) {
		ms.cancelEntryFn = func(_ context.Context, id string) error {
			return nil
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "entry cancelled", resp["message"])
	})

	t.Run("fail - already cancelled", func(t *testing.T) {
		ms.cancelEntryFn = func(_ context.Context, id string) error {
			return ErrJournalAlreadyCancelled
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-001/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.cancelEntryFn = func(_ context.Context, id string) error {
			return ErrJournalNotFound
		}

		req := httptest.NewRequest("POST", "/api/v1/journal-entries/JE-999/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetJournalEntryHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get entry", func(t *testing.T) {
		ms.getEntryByIDFn = func(_ context.Context, id string) (*JournalEntry, error) {
			return &JournalEntry{ID: "JE-001", Description: "Test", Status: JournalEntryPosted}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/journal-entries/JE-001", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var entry JournalEntry
		json.Unmarshal(w.Body.Bytes(), &entry)
		assert.Equal(t, "JE-001", entry.ID)
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.getEntryByIDFn = func(_ context.Context, id string) (*JournalEntry, error) {
			return nil, ErrJournalNotFound
		}

		req := httptest.NewRequest("GET", "/api/v1/journal-entries/JE-999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetJournalEntriesByDateRangeHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get entries by date range", func(t *testing.T) {
		ms.getEntriesByDateRangeFn = func(_ context.Context, from, to time.Time) ([]JournalEntry, error) {
			return []JournalEntry{
				{ID: "JE-001", Description: "Entry 1"},
				{ID: "JE-002", Description: "Entry 2"},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/journal-entries?from=2026-07-01&to=2026-07-31", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var entries []JournalEntry
		json.Unmarshal(w.Body.Bytes(), &entries)
		assert.Len(t, entries, 2)
	})

	t.Run("fail - missing from date", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/journal-entries?to=2026-07-31", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// ---------------------------------------------------------------------------
// Report tests
// ---------------------------------------------------------------------------

func TestTrialBalanceHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success trial balance", func(t *testing.T) {
		ms.trialBalanceFn = func(_ context.Context, year, month int) ([]AccountBalance, error) {
			return []AccountBalance{
				{AccountCode: "1111", AccountType: AccountTypeAsset, PeriodDebit: 1000000, PeriodCredit: 300000, ClosingBalance: 700000},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/reports/trial-balance?year=2026&month=7", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result []AccountBalance
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Len(t, result, 1)
	})

	t.Run("fail - missing params", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reports/trial-balance", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestBalanceSheetHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success balance sheet", func(t *testing.T) {
		ms.balanceSheetFn = func(_ context.Context, year, month int) ([]AccountBalance, error) {
			return []AccountBalance{
				{AccountCode: "1111", AccountType: AccountTypeAsset, ClosingBalance: 5000000},
				{AccountCode: "3111", AccountType: AccountTypeEquity, ClosingBalance: -5000000},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/reports/balance-sheet?year=2026&month=7", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result []AccountBalance
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Len(t, result, 2)
	})

	t.Run("fail - missing params", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reports/balance-sheet", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestIncomeStatementHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success income statement", func(t *testing.T) {
		ms.incomeStatementFn = func(_ context.Context, year, month int) ([]AccountBalance, error) {
			return []AccountBalance{
				{AccountCode: "5111", AccountType: AccountTypeRevenue, ClosingBalance: -20000000},
				{AccountCode: "6411", AccountType: AccountTypeExpense, ClosingBalance: 15000000},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/reports/income-statement?year=2026&month=7", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result []AccountBalance
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Len(t, result, 2)
	})

	t.Run("fail - missing params", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reports/income-statement", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// ---------------------------------------------------------------------------
// Period tests
// ---------------------------------------------------------------------------

func TestCreatePeriodHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create period", func(t *testing.T) {
		ms.createPeriodFn = func(_ context.Context, p *Period) error {
			return nil
		}
		body := `{"year":2026,"month":7,"start_date":"2026-07-01T00:00:00Z","end_date":"2026-07-31T00:00:00Z","status":"OPEN"}`
		req := httptest.NewRequest("POST", "/api/v1/periods", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("fail - bad request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/periods", strings.NewReader(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetPeriodHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get period", func(t *testing.T) {
		ms.getPeriodFn = func(_ context.Context, id string) (*Period, error) {
			return &Period{ID: "P-2026-07", Year: 2026, Month: 7, Status: PeriodOpen}, nil
		}
		req := httptest.NewRequest("GET", "/api/v1/periods/P-2026-07", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.getPeriodFn = func(_ context.Context, id string) (*Period, error) {
			return nil, ErrPeriodNotFound
		}
		req := httptest.NewRequest("GET", "/api/v1/periods/P-2099-01", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestListPeriodsHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success list periods", func(t *testing.T) {
		ms.getAllPeriodsFn = func(_ context.Context) ([]Period, error) {
			return []Period{{ID: "P-2026-07", Year: 2026, Month: 7, Status: PeriodOpen}}, nil
		}
		req := httptest.NewRequest("GET", "/api/v1/periods", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestClosePeriodHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success close period", func(t *testing.T) {
		ms.closePeriodFn = func(_ context.Context, id string) error {
			return nil
		}

		req := httptest.NewRequest("POST", "/api/v1/periods/P-2026-07/close", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "period closed", resp["message"])
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.closePeriodFn = func(_ context.Context, id string) error {
			return ErrPeriodNotFound
		}

		req := httptest.NewRequest("POST", "/api/v1/periods/P-2099-01/close", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("fail - already closed", func(t *testing.T) {
		ms.closePeriodFn = func(_ context.Context, id string) error {
			return ErrPeriodAlreadyClosed
		}

		req := httptest.NewRequest("POST", "/api/v1/periods/P-2026-07/close", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestReopenPeriodHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success reopen period", func(t *testing.T) {
		ms.reopenPeriodFn = func(_ context.Context, id string) error {
			return nil
		}

		req := httptest.NewRequest("POST", "/api/v1/periods/P-2026-07/reopen", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "period reopened", resp["message"])
	})

	t.Run("fail - has posted entries", func(t *testing.T) {
		ms.reopenPeriodFn = func(_ context.Context, id string) error {
			return errors.New("period has posted entries")
		}

		req := httptest.NewRequest("POST", "/api/v1/periods/P-2026-07/reopen", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// ---------------------------------------------------------------------------
// Exchange Rate tests
// ---------------------------------------------------------------------------

func TestCreateExchangeRateHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create exchange rate", func(t *testing.T) {
		ms.createExchangeRateFn = func(_ context.Context, rate *ExchangeRate) error {
			return nil
		}

		body := `{"currency_code":"USD","rate_date":"2026-07-24T00:00:00Z","buy_rate":25400,"sell_rate":25600,"average_rate":25500}`
		req := httptest.NewRequest("POST", "/api/v1/exchange-rates", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "exchange rate created", resp["message"])
	})

	t.Run("fail - bad request", func(t *testing.T) {
		ms.createExchangeRateFn = func(_ context.Context, rate *ExchangeRate) error {
			return errors.New("invalid currency code")
		}

		body := `{"currency_code":"US","rate_date":"2026-07-24T00:00:00Z","average_rate":25500}`
		req := httptest.NewRequest("POST", "/api/v1/exchange-rates", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListExchangeRatesHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success list exchange rates", func(t *testing.T) {
		ms.listExchangeRatesFn = func(_ context.Context) ([]ExchangeRate, error) {
			return []ExchangeRate{
				{CurrencyCode: "USD", AverageRate: 25500},
				{CurrencyCode: "EUR", AverageRate: 27800},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/exchange-rates", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var rates []ExchangeRate
		json.Unmarshal(w.Body.Bytes(), &rates)
		assert.Len(t, rates, 2)
	})
}

// ---------------------------------------------------------------------------
// Audit Log tests
// ---------------------------------------------------------------------------

func TestGetAuditLogHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get audit log", func(t *testing.T) {
		ms.getAuditLogFn = func(_ context.Context, limit int) ([]AuditEntry, error) {
			return []AuditEntry{
				{ID: "AUD-001", Username: "admin", Action: AuditActionCreate, EntityType: "journal_entry"},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/audit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var logs []AuditEntry
		json.Unmarshal(w.Body.Bytes(), &logs)
		assert.Len(t, logs, 1)
	})

	t.Run("success with limit param", func(t *testing.T) {
		ms.getAuditLogFn = func(_ context.Context, limit int) ([]AuditEntry, error) {
			assert.Equal(t, 10, limit)
			return []AuditEntry{}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/audit?limit=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetAuditLogByEntityHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get audit log by entity", func(t *testing.T) {
		ms.getAuditLogByEntityFn = func(_ context.Context, entityType, entityID string) ([]AuditEntry, error) {
			return []AuditEntry{
				{ID: "AUD-001", EntityType: "journal_entry", EntityID: "JE-001"},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/audit/entity?entity_type=journal_entry&entity_id=JE-001", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var logs []AuditEntry
		json.Unmarshal(w.Body.Bytes(), &logs)
		assert.Len(t, logs, 1)
	})

	t.Run("fail - missing params", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/audit/entity", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// ---------------------------------------------------------------------------
// User tests
// ---------------------------------------------------------------------------

func TestCreateUserHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create user", func(t *testing.T) {
		ms.createUserFn = func(_ context.Context, user *User, password string) error {
			user.ID = "USR-001"
			return nil
		}

		body := `{"user":{"username":"newuser","full_name":"New User","role":"accountant","is_active":true},"password":"secret123"}`
		req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "user created", resp["message"])
	})

	t.Run("fail - bad request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail - duplicate username", func(t *testing.T) {
		ms.createUserFn = func(_ context.Context, _ *User, _ string) error {
			return ErrUsernameExists
		}

		body := `{"user":{"username":"existing","full_name":"Existing","role":"accountant","is_active":true},"password":"secret123"}`
		req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListUsersHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success list users", func(t *testing.T) {
		ms.listUsersFn = func(_ context.Context) ([]User, error) {
			return []User{
				{ID: "USR-001", Username: "admin", FullName: "Admin User", Role: UserRoleAdmin},
				{ID: "USR-002", Username: "accountant", FullName: "Accountant User", Role: UserRoleAccountant},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var users []User
		json.Unmarshal(w.Body.Bytes(), &users)
		assert.Len(t, users, 2)
	})
}

func TestGetUserHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get user", func(t *testing.T) {
		ms.getUserFn = func(_ context.Context, id string) (*User, error) {
			return &User{ID: "USR-001", Username: "admin", FullName: "Admin User", Role: UserRoleAdmin}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/users/USR-001", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		json.Unmarshal(w.Body.Bytes(), &user)
		assert.Equal(t, "USR-001", user.ID)
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.getUserFn = func(_ context.Context, id string) (*User, error) {
			return nil, ErrUserNotFound
		}

		req := httptest.NewRequest("GET", "/api/v1/users/USR-999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// ---------------------------------------------------------------------------
// Auth / Login tests
// ---------------------------------------------------------------------------

func TestLoginHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success login", func(t *testing.T) {
		ms.loginFn = func(_ context.Context, username, password, ip string) (*AuthResult, error) {
			return &AuthResult{
				AccessToken:  "test-access-token",
				RefreshToken: "test-refresh-token",
				ExpiresIn:    900,
				TokenType:    "Bearer",
				User:         &User{ID: "USR-001", Username: "admin", Role: UserRoleAdmin},
			}, nil
		}

		body := `{"username":"admin","password":"secret"}`
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp["access_token"])
	})

	t.Run("fail - bad request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail - invalid credentials", func(t *testing.T) {
		ms.loginFn = func(_ context.Context, username, password, ip string) (*AuthResult, error) {
			return nil, ErrInvalidCredentials
		}

		body := `{"username":"admin","password":"wrong"}`
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// ---------------------------------------------------------------------------
// Current User tests
// ---------------------------------------------------------------------------

func TestGetCurrentUserHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get current user", func(t *testing.T) {
		ms.getUserFn = func(_ context.Context, id string) (*User, error) {
			return &User{ID: "test-user", Username: "testuser", Role: UserRoleAdmin}, nil
		}

		req := httptest.NewRequest("GET", "/api/v1/me", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		json.Unmarshal(w.Body.Bytes(), &user)
		assert.Equal(t, "test-user", user.ID)
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.getUserFn = func(_ context.Context, id string) (*User, error) {
			return nil, ErrUserNotFound
		}

		req := httptest.NewRequest("GET", "/api/v1/me", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// ─── COA: Freeze/Unfreeze Handler Tests ──────────────────────────

func TestFreezeAccountHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success freeze", func(t *testing.T) {
		ms.freezeAccountFn = func(_ context.Context, code, reason string) error {
			assert.Equal(t, "1111", code)
			assert.Equal(t, "audit hold", reason)
			return nil
		}
		body := `{"reason":"audit hold"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/accounts/1111/freeze", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - already frozen", func(t *testing.T) {
		ms.freezeAccountFn = func(_ context.Context, _, _ string) error {
			return ErrAccountAlreadyFrozen
		}
		body := `{"reason":"test"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/accounts/1111/freeze", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("fail - not found", func(t *testing.T) {
		ms.freezeAccountFn = func(_ context.Context, _, _ string) error {
			return ErrAccountNotFound
		}
		body := `{"reason":"test"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/accounts/9999/freeze", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUnfreezeAccountHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success unfreeze", func(t *testing.T) {
		ms.unfreezeAccountFn = func(_ context.Context, code, reason string) error {
			return nil
		}
		body := `{"reason":"resolved"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/accounts/1111/unfreeze", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - not frozen", func(t *testing.T) {
		ms.unfreezeAccountFn = func(_ context.Context, _, _ string) error {
			return ErrAccountNotFrozen
		}
		body := `{"reason":"test"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/accounts/1111/unfreeze", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

// ─── COA: Balance/Drill-down Handler Tests ───────────────────────

func TestGetAccountBalanceHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get balance", func(t *testing.T) {
		ms.getAccountBalanceFn = func(_ context.Context, code, periodID string) (*AccountBalance, error) {
			return &AccountBalance{AccountCode: "1111", ClosingBalance: 500000}, nil
		}
		req := httptest.NewRequest("GET", "/api/v1/coa/accounts/1111/balance?period_id=P-2026-07", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - missing period_id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/coa/accounts/1111/balance", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAccountUsageHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get usage", func(t *testing.T) {
		ms.getAccountUsageFn = func(_ context.Context, code string) (*AccountUsage, error) {
			return &AccountUsage{AccountCode: "1111", EntryCount: 5}, nil
		}
		req := httptest.NewRequest("GET", "/api/v1/coa/accounts/1111/usage", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// ─── COA: Approval Handler Tests ─────────────────────────────────

func TestCreateApprovalRequestHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create approval", func(t *testing.T) {
		ms.createApprovalRequestFn = func(_ context.Context, req *ApprovalRequest) error {
			return nil
		}
		body := `{"entity_type":"ACCOUNT","entity_id":"1111","request_type":"FREEZE","reason":"audit"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/approvals", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestApproveRequestHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success approve", func(t *testing.T) {
		ms.approveRequestFn = func(_ context.Context, id, reviewerID, note string) error {
			return nil
		}
		body := `{"note":"approved"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/approvals/APPR-001/approve", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - self approval", func(t *testing.T) {
		ms.approveRequestFn = func(_ context.Context, _, _, _ string) error {
			return ErrSelfApproval
		}
		body := `{"note":"ok"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/approvals/APPR-001/approve", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("fail - already processed", func(t *testing.T) {
		ms.approveRequestFn = func(_ context.Context, _, _, _ string) error {
			return ErrApprovalAlreadyProcessed
		}
		body := `{"note":"ok"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/approvals/APPR-001/approve", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestRejectRequestHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success reject", func(t *testing.T) {
		ms.rejectRequestFn = func(_ context.Context, id, reviewerID, note string) error {
			return nil
		}
		body := `{"note":"needs revision"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/approvals/APPR-001/reject", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - missing note", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/coa/approvals/APPR-001/reject", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// ─── COA: Versioning Handler Tests ───────────────────────────────

func TestCreateAccountVersionHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create version", func(t *testing.T) {
		ms.createAccountVersionFn = func(_ context.Context, reason string) (*AccountVersion, error) {
			return &AccountVersion{VersionNumber: "v1", ChangeReason: reason}, nil
		}
		body := `{"reason":"year-end snapshot"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/versions", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestCompareVersionsHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success compare", func(t *testing.T) {
		ms.compareVersionsFn = func(_ context.Context, v1, v2 string) (*VersionDiff, error) {
			return &VersionDiff{VersionFrom: v1, VersionTo: v2}, nil
		}
		req := httptest.NewRequest("GET", "/api/v1/coa/versions/compare?from=v1&to=v2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail - missing params", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/coa/versions/compare", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// ─── COA: Mapping Handler Tests ──────────────────────────────────

func TestCreateAccountMappingHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create mapping", func(t *testing.T) {
		ms.createAccountMappingFn = func(_ context.Context, m *AccountMapping) error {
			return nil
		}
		body := `{"source_regime":"TT200","target_regime":"TT99","old_code":"611","new_code":"632","mapping_type":"DIRECT"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/mappings", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("fail - duplicate", func(t *testing.T) {
		ms.createAccountMappingFn = func(_ context.Context, m *AccountMapping) error {
			return ErrMappingExists
		}
		body := `{"source_regime":"TT200","target_regime":"TT99","old_code":"611","new_code":"632","mapping_type":"DIRECT"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/mappings", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

// ─── COA: IFRS Mapping Handler Tests ─────────────────────────────

func TestCreateIFRSMappingHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success create IFRS mapping", func(t *testing.T) {
		ms.createIFRSMappingFn = func(_ context.Context, m *IFRSMapping) error {
			return nil
		}
		body := `{"vas_code":"1111","ifrs_code":"IFRS-1100","adjustment_type":"RECLASSIFY"}`
		req := httptest.NewRequest("POST", "/api/v1/coa/ifrs", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestGetIFRSMappingHandler(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms)
	r := setupTestRouter(h)

	t.Run("success get IFRS mapping", func(t *testing.T) {
		ms.getIFRSMappingFn = func(_ context.Context, vasCode string) (*IFRSMapping, error) {
			return &IFRSMapping{VASCode: "1111", IFRSCode: "IFRS-1100"}, nil
		}
		req := httptest.NewRequest("GET", "/api/v1/coa/ifrs/1111", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
