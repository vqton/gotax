package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func RegisterRoutes(r *gin.Engine, h *Handler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	const (
		resAccounts  = "accounts"
		resEntries   = "journal-entries"
		resReports   = "reports"
		resOB        = "opening-balances"
		resCF        = "carry-forward"
		resC99       = "circular99-mappings"
		resMig       = "balance-migrations"
		resPeriods   = "periods"
		resRates     = "exchange-rates"
		resAudit     = "audit"
		resUsers     = "users"
		resCOA       = "coa"
		resCOAAcc    = resCOA + "/accounts"
		resApprovals = resCOA + "/approvals"
		resVersions  = resCOA + "/versions"
		resMappings  = resCOA + "/mappings"
		resIFRS      = resCOA + "/ifrs"
	)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.POST("/totp/verify", h.Verify2FA)
	}

	v1 := r.Group("/api/v1", authMW)
	{
		accounts := v1.Group("/" + resAccounts)
		{
			accounts.POST("", h.CreateAccount)
			accounts.GET("", h.ListAccounts)
			accounts.GET("/:code", h.GetAccount)
			accounts.PUT("/:code", h.UpdateAccount)
			accounts.DELETE("/:code", adminMW, h.DeleteAccount)
		}

		entries := v1.Group("/" + resEntries)
		{
			entries.POST("", h.CreateEntry)
			entries.GET("", h.GetJournalEntries)
			entries.GET("/:id", h.GetJournalEntry)
			entries.POST("/:id/submit", h.SubmitEntry)
			entries.POST("/:id/review", h.ReviewEntry)
			entries.POST("/:id/approve", h.ApproveEntry)
			entries.POST("/:id/post", h.PostJournalEntry)
			entries.POST("/:id/cancel", h.CancelJournalEntry)
		}

		reports := v1.Group("/" + resReports)
		{
			reports.GET("/trial-balance", h.TrialBalance)
			reports.GET("/balance-sheet", h.BalanceSheet)
			reports.GET("/income-statement", h.IncomeStatement)
		}

		ob := v1.Group("/" + resOB)
		{
			ob.POST("", h.CreateOpeningBalance)
			ob.GET("", h.ListOpeningBalances)
			ob.GET("/:id", h.GetOpeningBalance)
			ob.PUT("/:id", h.UpdateOpeningBalance)
			ob.DELETE("/:id", h.DeleteOpeningBalance)
			ob.POST("/:id/submit", h.SubmitOpeningBalance)
			ob.POST("/:id/approve", h.ApproveOpeningBalance)
			ob.POST("/:id/correct", h.CorrectOpeningBalance)
			ob.POST("/:id/details", h.CreateOpeningBalanceDetail)
			ob.GET("/:id/details", h.GetOpeningBalanceDetails)
			ob.DELETE("/:id/details/:detailId", h.DeleteOpeningBalanceDetail)
			ob.GET("/totals", h.GetOpeningBalanceTotals)
			ob.POST("/import", h.ImportOpeningBalances)
			ob.GET("/report", h.DownloadOpeningBalancePDF)
		}

		cf := v1.Group("/" + resCF)
		{
			cf.POST("", adminMW, h.CarryForward)
			cf.GET("", h.GetCarryForwardLogs)
			cf.GET("/:id", h.GetCarryForwardLogByID)
		}

		c99 := v1.Group("/" + resC99)
		{
			c99.POST("", h.CreateCircular99Mapping)
			c99.GET("", h.ListCircular99Mappings)
			c99.GET("/:oldCode", h.GetCircular99MappingByOldCode)
		}

		mig := v1.Group("/" + resMig)
		{
			mig.POST("", h.CreateBalanceMigration)
			mig.GET("", h.ListBalanceMigrations)
			mig.GET("/:id", h.GetBalanceMigrationByID)
		}

		periods := v1.Group("/" + resPeriods)
		{
			periods.POST("", adminMW, h.CreatePeriod)
			periods.GET("", h.ListPeriods)
			periods.GET("/:id", h.GetPeriod)
			periods.POST("/:id/close", adminMW, h.ClosePeriod)
			periods.POST("/:id/reopen", adminMW, h.ReopenPeriod)
		}

		rates := v1.Group("/" + resRates)
		{
			rates.POST("", h.CreateExchangeRate)
			rates.GET("", h.ListExchangeRates)
		}

		audit := v1.Group("/" + resAudit, adminMW)
		{
			audit.GET("", h.GetAuditLog)
			audit.GET("/entity", h.GetAuditLogByEntity)
		}

		users := v1.Group("/" + resUsers, adminMW)
		{
			users.POST("", h.CreateUser)
			users.GET("", h.ListUsers)
			users.GET("/:id", h.GetUser)
		}

		coa := v1.Group("/" + resCOA)
		{
			acc := coa.Group("/accounts")
			{
				acc.POST("/:code/freeze", h.FreezeAccount)
				acc.POST("/:code/unfreeze", h.UnfreezeAccount)
				acc.GET("/:code/balance", h.GetAccountBalance)
				acc.GET("/:code/drill-down", h.GetAccountBalanceDrillDown)
				acc.GET("/:code/usage", h.GetAccountUsage)
				acc.POST("/:code/analysis", h.CreateAccountAnalysis)
				acc.GET("/:code/analysis", h.GetAccountAnalysis)
				acc.PUT("/:code/analysis", h.UpdateAccountAnalysis)
			}
			approvals := coa.Group("/approvals")
			{
				approvals.POST("", h.CreateApprovalRequest)
				approvals.GET("", h.ListApprovalRequests)
				approvals.POST("/:id/approve", h.ApproveRequest)
				approvals.POST("/:id/reject", h.RejectRequest)
			}
			versions := coa.Group("/versions")
			{
				versions.POST("", h.CreateAccountVersion)
				versions.GET("", h.ListVersions)
				versions.GET("/compare", h.CompareVersions)
				versions.GET("/:versionNumber", h.GetVersion)
			}
			mappings := coa.Group("/mappings")
			{
				mappings.POST("", h.CreateAccountMapping)
				mappings.GET("", h.ListMappings)
				mappings.GET("/:oldCode", h.GetMappingByOldCode)
			}
			ifrs := coa.Group("/ifrs")
			{
				ifrs.POST("", h.CreateIFRSMapping)
				ifrs.GET("", h.ListIFRSMappings)
				ifrs.GET("/:vasCode", h.GetIFRSMapping)
			}
		}

		authed := v1.Group("/auth")
		{
			authed.POST("/change-password", h.ChangePassword)
			authed.POST("/logout", h.Logout)
			authed.POST("/logout-all", h.LogoutAll)
			authed.POST("/totp/setup", h.SetupTOTP)
			authed.POST("/totp/confirm", h.ConfirmTOTP)
			authed.POST("/totp/disable", h.DisableTOTP)
			authed.POST("/backup-codes", h.GenerateBackupCodes)
			authed.GET("/sessions", h.ListSessions)
			authed.DELETE("/sessions/:id", h.RevokeSession)
		}

		v1.GET("/me", h.GetCurrentUser)
	}
}

func RegisterRoutesWithCompany(r *gin.Engine, h *Handler, ch *CompanyHandler, th *TaxHandler, cashH *CashHandler, bankH *BankHandler, purchaseH *PurchaseHandler, saleH *SaleHandler, whH *WarehouseHandler, faH *FAHandler, pwH *PayrollHandler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	RegisterRoutes(r, h, authMW, adminMW)
	RegisterCompanyRoutes(r, ch, authMW, adminMW)
	RegisterTaxRoutes(r, th, authMW)
	RegisterCashRoutes(r, cashH, authMW)
	RegisterBankRoutes(r, bankH, authMW)
	RegisterPurchaseRoutes(r, purchaseH, authMW)
	RegisterSaleRoutes(r, saleH, authMW)
	RegisterWarehouseRoutes(r, whH, authMW)
	RegisterFixedAssetRoutes(r, faH, authMW)
	RegisterPayrollRoutes(r, pwH, authMW)
}

// ─── Auth ──────────────────────────────────────────────────────────

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	ip := c.ClientIP()
	result, err := h.svc.Login(c.Request.Context(), req.Username, req.Password, ip)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRateLimited):
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrAccountLocked):
			c.JSON(http.StatusLocked, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrAccountInactive):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		}
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
		return
	}
	pair, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pair)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.svc.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "if email exists, reset link sent"})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.svc.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}

func (h *Handler) SetupTOTP(c *gin.Context) {
	userID := GetUserID(c)
	setup, err := h.svc.SetupTOTP(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, setup)
}

func (h *Handler) ConfirmTOTP(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.ConfirmTOTP(c.Request.Context(), userID, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "TOTP enabled"})
}

func (h *Handler) Verify2FA(c *gin.Context) {
	var req struct {
		TempToken string `json:"temp_token"`
		Code      string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	result, err := h.svc.Verify2FA(c.Request.Context(), req.TempToken, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) DisableTOTP(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password"`
		Code            string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.DisableTOTP(c.Request.Context(), userID, req.CurrentPassword, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "TOTP disabled"})
}

func (h *Handler) GenerateBackupCodes(c *gin.Context) {
	userID := GetUserID(c)
	codes, err := h.svc.GenerateBackupCodes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"backup_codes": codes})
}

func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.Logout(c.Request.Context(), userID, req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) LogoutAll(c *gin.Context) {
	userID := GetUserID(c)
	h.svc.LogoutAll(c.Request.Context(), userID)
	c.JSON(http.StatusOK, gin.H{"message": "logged out from all devices"})
}

func (h *Handler) ListSessions(c *gin.Context) {
	userID := GetUserID(c)
	sessions, err := h.svc.ListSessions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func (h *Handler) RevokeSession(c *gin.Context) {
	sessionID := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.RevokeSession(c.Request.Context(), userID, sessionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "session revoked"})
}

// ─── Accounts ──────────────────────────────────────────────────────

func (h *Handler) CreateAccount(c *gin.Context) {
	var account domain.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAccount(c.Request.Context(), &account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}

func (h *Handler) GetAccount(c *gin.Context) {
	code := c.Param("code")
	account, err := h.svc.GetAccount(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (h *Handler) ListAccounts(c *gin.Context) {
	activeOnly := c.DefaultQuery("active_only", "false") == "true"
	accounts, err := h.svc.GetAllAccounts(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	code := c.Param("code")
	var account domain.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	account.Code = code
	if err := h.svc.UpdateAccount(c.Request.Context(), &account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	code := c.Param("code")
	if err := h.svc.DeleteAccount(c.Request.Context(), code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}

// ─── Journal Entries ───────────────────────────────────────────────

func (h *Handler) CreateEntry(c *gin.Context) {
	var entry domain.JournalEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.CreateEntry(c.Request.Context(), &entry, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *Handler) SubmitEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.SubmitForReview(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submitted for review"})
}

func (h *Handler) ReviewEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.ReviewEntry(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry reviewed"})
}

func (h *Handler) ApproveEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.ApproveEntry(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry approved"})
}

func (h *Handler) PostJournalEntry(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.PostEntry(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry posted"})
}

func (h *Handler) CancelJournalEntry(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.CancelEntry(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry cancelled"})
}

func (h *Handler) GetJournalEntry(c *gin.Context) {
	id := c.Param("id")
	entry, err := h.svc.GetEntryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *Handler) GetJournalEntries(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	status := c.Query("status")
	if fromStr != "" && toStr != "" {
		from, err1 := time.Parse("2006-01-02", fromStr)
		to, err2 := time.Parse("2006-01-02", toStr)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		entries, err := h.svc.GetEntriesByDateRange(c.Request.Context(), from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, entries)
		return
	}
	if status != "" {
		entries, err := h.svc.GetEntriesByStatus(c.Request.Context(), domain.JournalEntryStatus(status))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, entries)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "specify from/to dates or status"})
}

// ─── Reports ───────────────────────────────────────────────────────

func (h *Handler) TrialBalance(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	balances, err := h.svc.TrialBalance(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

func (h *Handler) BalanceSheet(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	balances, err := h.svc.BalanceSheet(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

func (h *Handler) IncomeStatement(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	balances, err := h.svc.IncomeStatement(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

// ─── Periods ───────────────────────────────────────────────────────

func (h *Handler) CreatePeriod(c *gin.Context) {
	var period domain.Period
	if err := c.ShouldBindJSON(&period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreatePeriod(c.Request.Context(), &period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, period)
}

func (h *Handler) GetPeriod(c *gin.Context) {
	id := c.Param("id")
	period, err := h.svc.GetPeriod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, period)
}

func (h *Handler) ListPeriods(c *gin.Context) {
	periods, err := h.svc.GetAllPeriods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, periods)
}

func (h *Handler) ClosePeriod(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ClosePeriod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period closed"})
}

func (h *Handler) ReopenPeriod(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ReopenPeriod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period reopened"})
}

// ─── Exchange Rates ────────────────────────────────────────────────

func (h *Handler) CreateExchangeRate(c *gin.Context) {
	var rate domain.ExchangeRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateExchangeRate(c.Request.Context(), &rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rate)
}

func (h *Handler) ListExchangeRates(c *gin.Context) {
	rates, err := h.svc.ListExchangeRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rates)
}

// ─── Audit ─────────────────────────────────────────────────────────

func (h *Handler) GetAuditLog(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	entries, err := h.svc.GetAuditLog(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetAuditLogByEntity(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")
	if entityType == "" || entityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_type and entity_id are required"})
		return
	}
	entries, err := h.svc.GetAuditLogByEntity(c.Request.Context(), entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

// ─── Users ─────────────────────────────────────────────────────────

func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		User     *domain.User `json:"user"`
		Password string       `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateUser(c.Request.Context(), req.User, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req.User)
}

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := GetUserID(c)
	user, err := h.svc.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ─── COA: Freeze ───────────────────────────────────────────────────

func (h *Handler) FreezeAccount(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}
	if err := h.svc.FreezeAccount(c.Request.Context(), code, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account frozen"})
}

func (h *Handler) UnfreezeAccount(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}
	if err := h.svc.UnfreezeAccount(c.Request.Context(), code, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account unfrozen"})
}

// ─── COA: Balance ──────────────────────────────────────────────────

func (h *Handler) GetAccountBalance(c *gin.Context) {
	code := c.Param("code")
	periodID := c.Query("period_id")
	if periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_id is required"})
		return
	}
	balance, err := h.svc.GetAccountBalance(c.Request.Context(), code, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *Handler) GetAccountBalanceDrillDown(c *gin.Context) {
	code := c.Param("code")
	periodID := c.Query("period_id")
	if periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_id is required"})
		return
	}
	entries, err := h.svc.GetAccountBalanceDrillDown(c.Request.Context(), code, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetAccountUsage(c *gin.Context) {
	code := c.Param("code")
	usage, err := h.svc.GetAccountUsage(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usage)
}

// ─── COA: Approval ─────────────────────────────────────────────────

func (h *Handler) CreateApprovalRequest(c *gin.Context) {
	var req domain.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	req.RequestedBy = GetUserID(c)
	if err := h.svc.CreateApprovalRequest(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) ApproveRequest(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note string `json:"note"`
	}
	c.ShouldBindJSON(&req)
	userID := GetUserID(c)
	if err := h.svc.ApproveRequest(c.Request.Context(), id, userID, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request approved"})
}

func (h *Handler) RejectRequest(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "note is required"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.RejectRequest(c.Request.Context(), id, userID, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request rejected"})
}

func (h *Handler) ListApprovalRequests(c *gin.Context) {
	status := domain.ApprovalStatus(c.Query("status"))
	requests, err := h.svc.GetApprovalRequests(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// ─── COA: Versions ─────────────────────────────────────────────────

func (h *Handler) CreateAccountVersion(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}
	ver, err := h.svc.CreateAccountVersion(c.Request.Context(), req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ver)
}

func (h *Handler) GetVersion(c *gin.Context) {
	versionNumber := c.Param("versionNumber")
	ver, err := h.svc.GetVersion(c.Request.Context(), versionNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ver)
}

func (h *Handler) ListVersions(c *gin.Context) {
	versions, err := h.svc.ListVersions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func (h *Handler) CompareVersions(c *gin.Context) {
	v1 := c.Query("v1")
	v2 := c.Query("v2")
	if v1 == "" || v2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "v1 and v2 query params required"})
		return
	}
	diff, err := h.svc.CompareVersions(c.Request.Context(), v1, v2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, diff)
}

// ─── COA: Analysis ─────────────────────────────────────────────────

func (h *Handler) CreateAccountAnalysis(c *gin.Context) {
	code := c.Param("code")
	var analysis domain.AccountAnalysis
	if err := c.ShouldBindJSON(&analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	analysis.AccountCode = code
	if err := h.svc.CreateAccountAnalysis(c.Request.Context(), &analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, analysis)
}

func (h *Handler) GetAccountAnalysis(c *gin.Context) {
	code := c.Param("code")
	analysis, err := h.svc.GetAccountAnalysis(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analysis)
}

func (h *Handler) UpdateAccountAnalysis(c *gin.Context) {
	code := c.Param("code")
	var analysis domain.AccountAnalysis
	if err := c.ShouldBindJSON(&analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	analysis.AccountCode = code
	if err := h.svc.UpdateAccountAnalysis(c.Request.Context(), &analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analysis)
}

// ─── COA: Mappings ─────────────────────────────────────────────────

func (h *Handler) CreateAccountMapping(c *gin.Context) {
	var mapping domain.AccountMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAccountMapping(c.Request.Context(), &mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mapping)
}

func (h *Handler) GetMappingByOldCode(c *gin.Context) {
	oldCode := c.Param("oldCode")
	sourceRegime := c.Query("source_regime")
	if sourceRegime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_regime query param required"})
		return
	}
	mapping, err := h.svc.GetMappingByOldCode(c.Request.Context(), sourceRegime, oldCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mapping)
}

func (h *Handler) ListMappings(c *gin.Context) {
	sourceRegime := c.Query("source_regime")
	targetRegime := c.Query("target_regime")
	mappings, err := h.svc.ListMappings(c.Request.Context(), sourceRegime, targetRegime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mappings)
}

// ─── COA: IFRS ─────────────────────────────────────────────────────

func (h *Handler) CreateIFRSMapping(c *gin.Context) {
	var mapping domain.IFRSMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateIFRSMapping(c.Request.Context(), &mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mapping)
}

func (h *Handler) GetIFRSMapping(c *gin.Context) {
	vasCode := c.Param("vasCode")
	mapping, err := h.svc.GetIFRSMapping(c.Request.Context(), vasCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mapping)
}

func (h *Handler) ListIFRSMappings(c *gin.Context) {
	mappings, err := h.svc.ListIFRSMappings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mappings)
}
