package gl

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func RegisterRoutes(r *gin.Engine, h *Handler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	r.POST("/api/v1/auth/login", h.Login)

	v1 := r.Group("/api/v1", authMW)
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("", h.CreateAccount)
			accounts.GET("", h.ListAccounts)
			accounts.GET("/:code", h.GetAccount)
			accounts.PUT("/:code", h.UpdateAccount)
			accounts.DELETE("/:code", adminMW, h.DeleteAccount)
		}

		entries := v1.Group("/journal-entries")
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

		reports := v1.Group("/reports")
		{
			reports.GET("/trial-balance", h.TrialBalance)
			reports.GET("/balance-sheet", h.BalanceSheet)
			reports.GET("/income-statement", h.IncomeStatement)
		}

		periods := v1.Group("/periods")
		{
			periods.POST("", adminMW, h.CreatePeriod)
			periods.GET("", h.ListPeriods)
			periods.GET("/:id", h.GetPeriod)
			periods.POST("/:id/close", adminMW, h.ClosePeriod)
			periods.POST("/:id/reopen", adminMW, h.ReopenPeriod)
		}

		rates := v1.Group("/exchange-rates")
		{
			rates.POST("", h.CreateExchangeRate)
			rates.GET("", h.ListExchangeRates)
		}

		audit := v1.Group("/audit")
		{
			audit.GET("", adminMW, h.GetAuditLog)
			audit.GET("/entity", adminMW, h.GetAuditLogByEntity)
		}

		users := v1.Group("/users")
		{
			users.POST("", adminMW, h.CreateUser)
			users.GET("", adminMW, h.ListUsers)
			users.GET("/:id", h.GetUser)
		}

		coa := v1.Group("/coa")
		{
			accounts := coa.Group("/accounts")
			{
				accounts.POST("/:code/freeze", h.FreezeAccount)
				accounts.POST("/:code/unfreeze", h.UnfreezeAccount)
				accounts.GET("/:code/balance", h.GetAccountBalance)
				accounts.GET("/:code/drill-down", h.GetAccountBalanceDrillDown)
				accounts.GET("/:code/usage", h.GetAccountUsage)
				accounts.POST("/:code/analysis", h.CreateAccountAnalysis)
				accounts.GET("/:code/analysis", h.GetAccountAnalysis)
				accounts.PUT("/:code/analysis", h.UpdateAccountAnalysis)
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

		v1.GET("/me", h.GetCurrentUser)
	}
}

// Login authenticates user and returns JWT token
// @Summary      Login
// @Description  Authenticate with username/password, receive JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body object{username=string,password=string} true "Credentials"
// @Success      200 {object} map[string]interface{} "token + user"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Failure      401 {object} map[string]interface{} "invalid credentials"
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.svc.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// CreateAccount creates a new account
// @Summary      Create account
// @Description  Create a new account with code and name
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body Account true "Account payload"
// @Success      201 {object} map[string]interface{} "account created"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Failure      409 {object} map[string]interface{} "code conflict"
// @Router       /accounts [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var acc Account
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAccount(c.Request.Context(), &acc); err != nil {
		switch {
		case errors.Is(err, ErrAccountCodeExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "account created", "data": acc})
}

// GetAccount retrieves account by code
// @Summary      Get account
// @Description  Get a single account by its code
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        code path string true "Account code"
// @Success      200 {object} map[string]interface{} "account found"
// @Failure      404 {object} map[string]interface{} "not found"
// @Router       /accounts/{code} [get]
func (h *Handler) GetAccount(c *gin.Context) {
	code := c.Param("code")
	acc, err := h.svc.GetAccount(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}
	c.JSON(http.StatusOK, acc)
}

// ListAccounts lists all accounts
// @Summary      List accounts
// @Description  List all accounts, optionally filter by active status
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        active query bool false "Filter by active status"
// @Success      200 {object} map[string]interface{} "accounts list"
// @Failure      500 {object} map[string]interface{} "server error"
// @Router       /accounts [get]
func (h *Handler) ListAccounts(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "false") == "true"
	accounts, err := h.svc.GetAllAccounts(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

// UpdateAccount updates an existing account
// @Summary      Update account
// @Description  Update account details by code
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        code path string true "Account code"
// @Param        request body Account true "Account payload"
// @Success      200 {object} map[string]interface{} "account updated"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Failure      404 {object} map[string]interface{} "not found"
// @Router       /accounts/{code} [put]
func (h *Handler) UpdateAccount(c *gin.Context) {
	code := c.Param("code")
	var acc Account
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if acc.Code != code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code mismatch"})
		return
	}
	if err := h.svc.UpdateAccount(c.Request.Context(), &acc); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account updated"})
}

// DeleteAccount deletes an account (Admin only)
// @Summary      Delete account
// @Description  Delete account by code. Admin only.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        code path string true "Account code"
// @Success      200 {object} map[string]interface{} "account deleted"
// @Failure      404 {object} map[string]interface{} "not found"
// @Failure      409 {object} map[string]interface{} "has children"
// @Router       /accounts/{code} [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	code := c.Param("code")
	if err := h.svc.DeleteAccount(c.Request.Context(), code); err != nil {
		switch {
		case errors.Is(err, ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrAccountHasChildren):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}

// CreateEntry creates a new journal entry
// @Summary      Create entry
// @Description  Create a new journal entry with lines
// @Tags         journal-entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body JournalEntry true "Journal entry payload"
// @Success      201 {object} map[string]interface{} "entry created"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /journal-entries [post]
func (h *Handler) CreateEntry(c *gin.Context) {
	var entry JournalEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.CreateEntry(c.Request.Context(), &entry, userID); err != nil {
		switch {
		case errors.Is(err, ErrAccountNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrAccountInactive):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "entry created", "data": entry})
}

// SubmitEntry submits entry for review
// @Summary      Submit entry
// @Description  Submit a journal entry for review
// @Tags         journal-entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Entry ID"
// @Success      200 {object} map[string]interface{} "submitted for review"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /journal-entries/{id}/submit [post]
func (h *Handler) SubmitEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.SubmitForReview(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry submitted for review"})
}

// ReviewEntry reviews a submitted entry
// @Summary      Review entry
// @Description  Review a journal entry that has been submitted
// @Tags         journal-entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Entry ID"
// @Success      200 {object} map[string]interface{} "entry reviewed"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /journal-entries/{id}/review [post]
func (h *Handler) ReviewEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.ReviewEntry(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry reviewed"})
}

// ApproveEntry approves a reviewed entry
// @Summary      Approve entry
// @Description  Approve a journal entry that has been reviewed
// @Tags         journal-entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Entry ID"
// @Success      200 {object} map[string]interface{} "entry approved"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /journal-entries/{id}/approve [post]
func (h *Handler) ApproveEntry(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := h.svc.ApproveEntry(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry approved"})
}

// PostJournalEntry posts an approved entry to ledger
// @Summary      Post entry
// @Description  Post an approved journal entry to the general ledger
// @Tags         journal-entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Entry ID"
// @Success      200 {object} map[string]interface{} "entry posted"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Failure      409 {object} map[string]interface{} "conflict"
// @Router       /journal-entries/{id}/post [post]
func (h *Handler) PostJournalEntry(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.PostEntry(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, ErrJournalAlreadyPosted):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, ErrJournalAlreadyCancelled):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, ErrPeriodNotFound), errors.Is(err, ErrJournalPeriodClosed):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry posted"})
}

// GetJournalEntry retrieves a journal entry by ID
// @Summary      Get journal entry
// @Description  Get a single journal entry by its ID
// @Tags         journal-entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Entry ID"
// @Success      200 {object} map[string]interface{} "entry found"
// @Failure      404 {object} map[string]interface{} "not found"
// @Router       /journal-entries/{id} [get]
func (h *Handler) GetJournalEntry(c *gin.Context) {
	id := c.Param("id")
	entry, err := h.svc.GetEntryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "journal entry not found"})
		return
	}
	c.JSON(http.StatusOK, entry)
}

// GetJournalEntries retrieves journal entries by date range
// @Summary      List journal entries
// @Description  Get all journal entries within a date range
// @Tags         journal-entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        from query string true "Start date (YYYY-MM-DD)"
// @Param        to query string true "End date (YYYY-MM-DD)"
// @Success      200 {object} map[string]interface{} "entries list"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /journal-entries [get]
func (h *Handler) GetJournalEntries(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params required"})
		return
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format, use YYYY-MM-DD"})
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format, use YYYY-MM-DD"})
		return
	}
	entries, err := h.svc.GetEntriesByDateRange(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

// CancelJournalEntry cancels a journal entry
// @Summary      Cancel entry
// @Description  Cancel a journal entry by ID
// @Tags         journal-entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Entry ID"
// @Success      200 {object} map[string]interface{} "entry cancelled"
// @Failure      404 {object} map[string]interface{} "not found"
// @Failure      409 {object} map[string]interface{} "already cancelled"
// @Router       /journal-entries/{id}/cancel [post]
func (h *Handler) CancelJournalEntry(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.CancelEntry(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, ErrJournalNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrJournalAlreadyCancelled):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "entry cancelled"})
}

// TrialBalance generates trial balance report
// @Summary      Trial balance
// @Description  Generate trial balance for a given year/month
// @Tags         reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        year query int true "Fiscal year"
// @Param        month query int true "Fiscal month (1-12)"
// @Success      200 {object} map[string]interface{} "trial balance"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /reports/trial-balance [get]
func (h *Handler) TrialBalance(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month query params required"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month (1-12)"})
		return
	}
	balances, err := h.svc.TrialBalance(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

// BalanceSheet generates balance sheet report
// @Summary      Balance sheet
// @Description  Generate balance sheet for a given year/month
// @Tags         reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        year query int true "Fiscal year"
// @Param        month query int true "Fiscal month (1-12)"
// @Success      200 {object} map[string]interface{} "balance sheet"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /reports/balance-sheet [get]
func (h *Handler) BalanceSheet(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month query params required"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month (1-12)"})
		return
	}
	balances, err := h.svc.BalanceSheet(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

// IncomeStatement generates income statement report
// @Summary      Income statement
// @Description  Generate income statement for a given year/month
// @Tags         reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        year query int true "Fiscal year"
// @Param        month query int true "Fiscal month (1-12)"
// @Success      200 {object} map[string]interface{} "income statement"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /reports/income-statement [get]
func (h *Handler) IncomeStatement(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month query params required"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month (1-12)"})
		return
	}
	balances, err := h.svc.IncomeStatement(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

// CreatePeriod creates a new accounting period (Admin only)
// @Summary      Create period
// @Description  Create a new accounting period. Admin only.
// @Tags         periods
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body Period true "Period payload"
// @Success      201 {object} map[string]interface{} "period created"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Failure      500 {object} map[string]interface{} "server error"
// @Router       /periods [post]
func (h *Handler) CreatePeriod(c *gin.Context) {
	var period Period
	if err := c.ShouldBindJSON(&period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreatePeriod(c.Request.Context(), &period); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "period created", "data": period})
}

// GetPeriod retrieves period by ID
// @Summary      Get period
// @Description  Get a single accounting period by ID
// @Tags         periods
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Period ID"
// @Success      200 {object} map[string]interface{} "period found"
// @Failure      404 {object} map[string]interface{} "not found"
// @Router       /periods/{id} [get]
func (h *Handler) GetPeriod(c *gin.Context) {
	id := c.Param("id")
	period, err := h.svc.GetPeriod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "period not found"})
		return
	}
	c.JSON(http.StatusOK, period)
}

// ListPeriods lists all accounting periods
// @Summary      List periods
// @Description  List all accounting periods
// @Tags         periods
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "periods list"
// @Failure      500 {object} map[string]interface{} "server error"
// @Router       /periods [get]
func (h *Handler) ListPeriods(c *gin.Context) {
	periods, err := h.svc.GetAllPeriods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, periods)
}

// ClosePeriod closes an accounting period (Admin only)
// @Summary      Close period
// @Description  Close an accounting period by ID. Admin only.
// @Tags         periods
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Period ID"
// @Success      200 {object} map[string]interface{} "period closed"
// @Failure      404 {object} map[string]interface{} "not found"
// @Failure      409 {object} map[string]interface{} "already closed"
// @Router       /periods/{id}/close [post]
func (h *Handler) ClosePeriod(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ClosePeriod(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, ErrPeriodNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrPeriodAlreadyClosed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period closed"})
}

// ReopenPeriod reopens a closed period (Admin only)
// @Summary      Reopen period
// @Description  Reopen a closed accounting period by ID. Admin only.
// @Tags         periods
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Period ID"
// @Success      200 {object} map[string]interface{} "period reopened"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /periods/{id}/reopen [post]
func (h *Handler) ReopenPeriod(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ReopenPeriod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period reopened"})
}

// CreateExchangeRate creates a new exchange rate
// @Summary      Create exchange rate
// @Description  Create a new exchange rate entry
// @Tags         exchange-rates
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body ExchangeRate true "Exchange rate payload"
// @Success      201 {object} map[string]interface{} "rate created"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /exchange-rates [post]
func (h *Handler) CreateExchangeRate(c *gin.Context) {
	var rate ExchangeRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateExchangeRate(c.Request.Context(), &rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "exchange rate created", "data": rate})
}

// ListExchangeRates lists all exchange rates
// @Summary      List exchange rates
// @Description  List all exchange rates
// @Tags         exchange-rates
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "rates list"
// @Failure      500 {object} map[string]interface{} "server error"
// @Router       /exchange-rates [get]
func (h *Handler) ListExchangeRates(c *gin.Context) {
	rates, err := h.svc.ListExchangeRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rates)
}

// GetAuditLog retrieves audit log entries (Admin only)
// @Summary      Get audit log
// @Description  Get audit log entries with optional limit. Admin only.
// @Tags         audit
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Max entries (default 50)"
// @Success      200 {object} map[string]interface{} "audit log"
// @Failure      500 {object} map[string]interface{} "server error"
// @Router       /audit [get]
func (h *Handler) GetAuditLog(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	logs, err := h.svc.GetAuditLog(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GetAuditLogByEntity retrieves audit log filtered by entity (Admin only)
// @Summary      Get audit log by entity
// @Description  Get audit log entries filtered by entity type and ID. Admin only.
// @Tags         audit
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        entity_type query string true "Entity type"
// @Param        entity_id query string true "Entity ID"
// @Success      200 {object} map[string]interface{} "audit log"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Failure      500 {object} map[string]interface{} "server error"
// @Router       /audit/entity [get]
func (h *Handler) GetAuditLogByEntity(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")
	if entityType == "" || entityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_type and entity_id query params required"})
		return
	}
	logs, err := h.svc.GetAuditLogByEntity(c.Request.Context(), entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// CreateUser creates a new user (Admin only)
// @Summary      Create user
// @Description  Create a new user with password. Admin only.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object{user=User,password=string} true "User + password"
// @Success      201 {object} map[string]interface{} "user created"
// @Failure      400 {object} map[string]interface{} "bad request"
// @Router       /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		User     User   `json:"user"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateUser(c.Request.Context(), &req.User, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created", "data": req.User})
}

// ListUsers lists all users (Admin only)
// @Summary      List users
// @Description  List all registered users. Admin only.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "users list"
// @Failure      500 {object} map[string]interface{} "server error"
// @Router       /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser retrieves user by ID
// @Summary      Get user
// @Description  Get a single user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{} "user found"
// @Failure      404 {object} map[string]interface{} "not found"
// @Router       /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetCurrentUser retrieves the authenticated user's profile
// @Summary      Get current user
// @Description  Get profile of the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "user profile"
// @Failure      404 {object} map[string]interface{} "not found"
// @Router       /me [get]
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := GetUserID(c)
	user, err := h.svc.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ─── COA: Freeze / Unfreeze ──────────────────────────────────────

func (h *Handler) FreezeAccount(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason required"})
		return
	}
	if err := h.svc.FreezeAccount(c.Request.Context(), code, req.Reason); err != nil {
		switch {
		case errors.Is(err, ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrAccountAlreadyFrozen):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account frozen"})
}

func (h *Handler) UnfreezeAccount(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason required"})
		return
	}
	if err := h.svc.UnfreezeAccount(c.Request.Context(), code, req.Reason); err != nil {
		switch {
		case errors.Is(err, ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrAccountNotFrozen):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account unfrozen"})
}

// ─── COA: Account Balance & Drill-down ───────────────────────────

func (h *Handler) GetAccountBalance(c *gin.Context) {
	code := c.Param("code")
	periodID := c.Query("period_id")
	if periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_id query param required"})
		return
	}
	bal, err := h.svc.GetAccountBalance(c.Request.Context(), code, periodID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bal)
}

func (h *Handler) GetAccountBalanceDrillDown(c *gin.Context) {
	code := c.Param("code")
	periodID := c.Query("period_id")
	if periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_id query param required"})
		return
	}
	entries, err := h.svc.GetAccountBalanceDrillDown(c.Request.Context(), code, periodID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetAccountUsage(c *gin.Context) {
	code := c.Param("code")
	usage, err := h.svc.GetAccountUsage(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usage)
}

// ─── COA: Approval Requests ──────────────────────────────────────

func (h *Handler) CreateApprovalRequest(c *gin.Context) {
	var req ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	req.RequestedBy = GetUserID(c)
	if err := h.svc.CreateApprovalRequest(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "approval request created", "data": req})
}

func (h *Handler) ApproveRequest(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err == nil {
		_ = req
	}
	reviewerID := GetUserID(c)
	if err := h.svc.ApproveRequest(c.Request.Context(), id, reviewerID, req.Note); err != nil {
		switch {
		case errors.Is(err, ErrApprovalNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrSelfApproval):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, ErrApprovalAlreadyProcessed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "rejection note required"})
		return
	}
	reviewerID := GetUserID(c)
	if err := h.svc.RejectRequest(c.Request.Context(), id, reviewerID, req.Note); err != nil {
		switch {
		case errors.Is(err, ErrApprovalNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrApprovalAlreadyProcessed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request rejected"})
}

func (h *Handler) ListApprovalRequests(c *gin.Context) {
	status := ApprovalStatus(c.Query("status"))
	requests, err := h.svc.GetApprovalRequests(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// ─── COA: Versioning ─────────────────────────────────────────────

func (h *Handler) CreateAccountVersion(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason required"})
		return
	}
	ver, err := h.svc.CreateAccountVersion(c.Request.Context(), req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "version created", "data": ver})
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
	v1 := c.Query("from")
	v2 := c.Query("to")
	if v1 == "" || v2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params required"})
		return
	}
	diff, err := h.svc.CompareVersions(c.Request.Context(), v1, v2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, diff)
}

// ─── COA: Account Analysis ───────────────────────────────────────

func (h *Handler) CreateAccountAnalysis(c *gin.Context) {
	var analysis AccountAnalysis
	if err := c.ShouldBindJSON(&analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAccountAnalysis(c.Request.Context(), &analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "analysis created", "data": analysis})
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
	var analysis AccountAnalysis
	if err := c.ShouldBindJSON(&analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.UpdateAccountAnalysis(c.Request.Context(), &analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "analysis updated", "data": analysis})
}

// ─── COA: Account Mappings ───────────────────────────────────────

func (h *Handler) CreateAccountMapping(c *gin.Context) {
	var mapping AccountMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAccountMapping(c.Request.Context(), &mapping); err != nil {
		switch {
		case errors.Is(err, ErrMappingExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "mapping created", "data": mapping})
}

func (h *Handler) GetMappingByOldCode(c *gin.Context) {
	sourceRegime := c.Query("source_regime")
	oldCode := c.Param("oldCode")
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

// ─── COA: IFRS Mappings ──────────────────────────────────────────

func (h *Handler) CreateIFRSMapping(c *gin.Context) {
	var mapping IFRSMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateIFRSMapping(c.Request.Context(), &mapping); err != nil {
		switch {
		case errors.Is(err, ErrIFRSMappingExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "IFRS mapping created", "data": mapping})
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
