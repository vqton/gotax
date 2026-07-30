// GoTax GL Module API
//
//     Schemes: http
//     Host: localhost:8080
//     BasePath: /api/v1
//     Version: 1.0.0
//     License: MIT
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     SecurityDefinitions:
//     bearerAuth:
//       type: apiKey
//       in: header
//       name: Authorization
//
// swagger:meta
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gotax/docs"
	"gotax/internal/auth"
	"gotax/internal/db"
	"gotax/internal/domain"
	"gotax/internal/handler"
	gotaxi18n "gotax/internal/i18n"
	"gotax/internal/repository"
	"gotax/internal/service"
)

// @title           GoTax GL API
// @version         1.0.0
// @description     General Ledger API — Circular 99/2025/TT-BTC compliant
// @termsOfService  https://gotax.vn/terms

// @contact.name   API Support
// @contact.url    https://gotax.vn/support
// @contact.email  support@gotax.vn

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	ctx := context.Background()

	secret := os.Getenv("JWT_SECRET")
	auth.SetJWTSecret(secret)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
r.LoadHTMLGlob("web/auth/*.html")
r.Static("/assets", "./web/static")

r.GET("/ping", func(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GoTax GL Server"})
})
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

r.GET("/login", func(c *gin.Context) {
	c.HTML(200, "login.html", nil)
})
r.GET("/2fa", func(c *gin.Context) {
	c.HTML(200, "2fa.html", nil)
})
r.GET("/forgot-password", func(c *gin.Context) {
	c.HTML(200, "forgot-password.html", nil)
})
r.GET("/reset-password", func(c *gin.Context) {
	c.HTML(200, "reset-password.html", nil)
})

	i18nL := gotaxi18n.MustNew()
	r.Use(handler.I18nMiddleware(i18nL))

	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		log.Println("using PostgreSQL backend (Clean Architecture)")
		cfg := db.DefaultPGConfig()
		cfg.DSN = dsn

		pool, err := db.NewPool(ctx, cfg)
		if err != nil {
			log.Fatalf("connect to PostgreSQL: %v", err)
		}
		defer pool.Close()

		if err := db.RunMigrations(ctx, pool); err != nil {
			log.Fatalf("run migrations: %v", err)
		}

		accRepo := repository.NewPGAccountRepo(pool)
		jeRepo := repository.NewPGJournalRepo(pool)
		perRepo := repository.NewPGPeriodRepo(pool)
		userRepo := repository.NewPGUserRepo(pool)
		auditRepo := repository.NewPGAuditLogRepo(pool)
		rateRepo := repository.NewPGExchangeRateRepo(pool)
		templateRepo := repository.NewPGClosingTemplateRepo(pool)

		obRepo := repository.NewPGOpeningBalanceRepo(pool)
		cashRepo := repository.NewPGCashRepo(pool)
		approvalRepo := repository.NewPGApprovalRepo(pool)
		versionRepo := repository.NewPGAccountVersionRepo(pool)
		mappingRepo := repository.NewPGAccountMappingRepo(pool)
		analysisRepo := repository.NewPGAccountAnalysisRepo(pool)
		ifrsRepo := repository.NewPGIFRSMappingRepo(pool)
		refreshRepo := repository.NewPGRefreshTokenRepo(pool)
		resetRepo := repository.NewPGPasswordResetTokenRepo(pool)
		svc := service.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)
		h := handler.NewHandler(svc)

		companyRepo := repository.NewPGCompanyRepo(pool)
		companySvc := service.NewCompanyService(companyRepo)
		companyH := handler.NewCompanyHandler(companySvc)

		authMW := handler.AuthMiddleware()
		adminMW := handler.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleChiefAccountant)

		taxRepo := repository.NewPGTaxRepo(pool)
		taxSvc := service.NewTaxService(taxRepo)
		taxH := handler.NewTaxHandler(taxSvc)
		cashH := handler.NewCashHandler(svc)
		bankRepo := repository.NewPGBankRepo(pool)
		bankSvc := service.NewBankService(bankRepo)
		bankH := handler.NewBankHandler(bankSvc, companySvc)

	purchaseRepo := repository.NewPGPurchaseRepo(pool)
	purchaseSvc := service.NewPurchaseService(purchaseRepo, purchaseRepo, purchaseRepo, purchaseRepo, purchaseRepo, purchaseRepo)
	purchaseH := handler.NewPurchaseHandler(purchaseSvc)
	saleRepo := repository.NewPGSaleRepo(pool)
	saleSvc := service.NewSaleService(saleRepo, saleRepo, saleRepo, saleRepo, saleRepo, saleRepo, saleRepo, saleRepo, svc)
	saleH := handler.NewSaleHandler(saleSvc)
	whRepo := repository.NewPGWarehouseRepo(pool)
	whSvc := service.NewWarehouseService(whRepo, whRepo, whRepo, whRepo, whRepo, whRepo, whRepo, whRepo, whRepo, purchaseRepo, svc)
	whH := handler.NewWarehouseHandler(whSvc)
	handler.RegisterRoutesWithCompany(r, h, companyH, taxH, cashH, bankH, purchaseH, saleH, whH, authMW, adminMW)
	log.Println("GoTax GL server (PG) starting on :8080")
		r.Run(":8080")
		return
	}

	log.Println("using in-memory backend (Clean Architecture)")
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
	cashRepo := repository.NewMemoryCashRepo()
	svc := service.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)
	h := handler.NewHandler(svc)

	companyRepo := repository.NewMemoryCompanyRepo()
	companySvc := service.NewCompanyService(companyRepo)
	companyH := handler.NewCompanyHandler(companySvc)

	authMW := handler.AuthMiddleware()
	adminMW := handler.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleChiefAccountant)

	taxRepo := repository.NewMemoryTaxRepo()
	taxSvc := service.NewTaxService(taxRepo)
	taxH := handler.NewTaxHandler(taxSvc)
	cashH := handler.NewCashHandler(svc)
	bankRepo := repository.NewMemoryBankRepo()
	bankSvc := service.NewBankService(bankRepo)
	bankH := handler.NewBankHandler(bankSvc, companySvc)

	memPurchaseRepo := repository.NewMemoryPurchaseRepo()
	purchaseSvc := service.NewPurchaseService(memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo)
	purchaseH := handler.NewPurchaseHandler(purchaseSvc)
	saleRepo := repository.NewMemorySaleRepo()
	saleSvc := service.NewSaleService(saleRepo, saleRepo, saleRepo, saleRepo, saleRepo, saleRepo, saleRepo, saleRepo, svc)
	saleH := handler.NewSaleHandler(saleSvc)
	memWhRepo := repository.NewMemoryWarehouseRepo()
	memCatRepo := repository.NewMemoryItemCategoryRepo()
	memItemRepo := repository.NewMemoryItemRepo()
	memBalRepo := repository.NewMemoryStockBalanceRepo()
	memTxnRepo := repository.NewMemoryInventoryTransactionRepo()
	memTrfRepo := repository.NewMemoryStockTransferRepo()
	memAdjRepo := repository.NewMemoryStockAdjustmentRepo()
	memTakeRepo := repository.NewMemoryStockTakeRepo()
	memValRepo := repository.NewMemoryInventoryValuationRunRepo()
	whSvc := service.NewWarehouseService(memWhRepo, memCatRepo, memItemRepo, memBalRepo, memTxnRepo, memTrfRepo, memAdjRepo, memTakeRepo, memValRepo, memPurchaseRepo, svc)
	whH := handler.NewWarehouseHandler(whSvc)
	handler.RegisterRoutesWithCompany(r, h, companyH, taxH, cashH, bankH, purchaseH, saleH, whH, authMW, adminMW)
	log.Println("GoTax GL server (CA) starting on :8080")
	r.Run(":8080")
}
