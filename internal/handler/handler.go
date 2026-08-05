package handler

import (
	"github.com/gin-gonic/gin"

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
