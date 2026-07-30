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
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "gotax/docs"
	"gotax/internal/auth"
	"gotax/internal/config"
	"gotax/internal/db"
	"gotax/internal/domain"
	"gotax/internal/handler"
	gotaxi18n "gotax/internal/i18n"
	"gotax/internal/logger"
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

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	zapLog, err := logger.New(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init logger: %v\n", err)
		os.Exit(1)
	}
	defer zapLog.Sync()

	zap.ReplaceGlobals(zapLog)

	auth.SetJWTSecret(cfg.JWTSecret)

	r := gin.New()
	r.Use(logger.GinMiddleware(zapLog))
	r.Use(gin.Recovery())
	r.SetTrustedProxies(cfg.TrustedProxies)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60 * 1000000000,
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

	if cfg.DatabaseURL != "" {
		zap.L().Info("using PostgreSQL backend")
		gormCfg := db.GormConfig{
			DSN:                cfg.DatabaseURL,
			MaxOpenConns:       cfg.GORMMaxOpenConns,
			MaxIdleConns:       cfg.GORMMaxIdleConns,
			ConnMaxLifetimeMS:  cfg.GORMConnMaxLifetimeMS,
		}
		gormDB, err := db.NewGorm(ctx, gormCfg)
		if err != nil {
			zap.L().Fatal("connect to PostgreSQL", zap.Error(err))
		}
		defer db.CloseGorm(gormDB)

		if err := db.RunGolangMigrate(cfg.DatabaseURL); err != nil {
			zap.L().Fatal("run migrations", zap.Error(err))
		}

		accRepo := repository.NewPGAccountRepo(gormDB)
		jeRepo := repository.NewPGJournalRepo(gormDB)
		perRepo := repository.NewPGPeriodRepo(gormDB)
		userRepo := repository.NewPGUserRepo(gormDB)
		auditRepo := repository.NewPGAuditLogRepo(gormDB)
		rateRepo := repository.NewPGExchangeRateRepo(gormDB)
		templateRepo := repository.NewPGClosingTemplateRepo(gormDB)

		obRepo := repository.NewPGOpeningBalanceRepo(gormDB)
		cashRepo := repository.NewPGCashRepo(gormDB)
		approvalRepo := repository.NewPGApprovalRepo(gormDB)
		faRepo := repository.NewPGFARepo(gormDB)
		versionRepo := repository.NewPGAccountVersionRepo(gormDB)
		mappingRepo := repository.NewPGAccountMappingRepo(gormDB)
		analysisRepo := repository.NewPGAccountAnalysisRepo(gormDB)
		ifrsRepo := repository.NewPGIFRSMappingRepo(gormDB)
		refreshRepo := repository.NewPGRefreshTokenRepo(gormDB)
		resetRepo := repository.NewPGPasswordResetTokenRepo(gormDB)
		svc := service.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)
		h := handler.NewHandler(svc)

		companyRepo := repository.NewPGCompanyRepo(gormDB)
		companySvc := service.NewCompanyService(companyRepo)
		companyH := handler.NewCompanyHandler(companySvc)

		authMW := handler.AuthMiddleware()
		adminMW := handler.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleChiefAccountant)

		taxRepo := repository.NewPGTaxRepo(gormDB)
		taxSvc := service.NewTaxService(taxRepo)
		taxH := handler.NewTaxHandler(taxSvc)
		cashH := handler.NewCashHandler(svc)
		bankRepo := repository.NewPGBankRepo(gormDB)
		bankSvc := service.NewBankService(bankRepo)
		bankH := handler.NewBankHandler(bankSvc, companySvc)

	supplierRepo := repository.NewPGSupplierRepo(gormDB)
	purchaseOrderRepo := repository.NewPGPurchaseRepo(gormDB)
	grnRepo := repository.NewPGGRNRepo(gormDB)
	supplierInvoiceRepo := repository.NewPGSupplierInvoiceRepo(gormDB)
	apTxnRepo := repository.NewPGAPTransactionRepo(gormDB)
	costAllocRepo := repository.NewPGCostAllocationRepo(gormDB)
	purchaseSvc := service.NewPurchaseService(supplierRepo, purchaseOrderRepo, grnRepo, supplierInvoiceRepo, apTxnRepo, costAllocRepo)
	purchaseH := handler.NewPurchaseHandler(purchaseSvc)
	customerRepo := repository.NewPGCustomerRepo(gormDB)
	saleOrderRepo := repository.NewPGSaleOrderRepo(gormDB)
	dnRepo := repository.NewPGDeliveryNoteRepo(gormDB)
	custInvRepo := repository.NewPGCustomerInvoiceRepo(gormDB)
	custReceiptRepo := repository.NewPGCustomerReceiptRepo(gormDB)
	creditNoteRepo := repository.NewPGCreditNoteRepo(gormDB)
	arTxnRepo := repository.NewPGARTransactionRepo(gormDB)
	salesQuotRepo := repository.NewPGSalesQuotationRepo(gormDB)
	saleSvc := service.NewSaleService(customerRepo, saleOrderRepo, dnRepo, custInvRepo, custReceiptRepo, creditNoteRepo, arTxnRepo, salesQuotRepo, svc)
	saleH := handler.NewSaleHandler(saleSvc)
	whRepo := repository.NewPGWarehouseRepo(gormDB)
	catRepo := repository.NewPGItemCategoryRepo(gormDB)
	itemRepo := repository.NewPGItemRepo(gormDB)
	balRepo := repository.NewPGStockBalanceRepo(gormDB)
	invTxnRepo := repository.NewPGInventoryTransactionRepo(gormDB)
	trfRepo := repository.NewPGStockTransferRepo(gormDB)
	adjRepo := repository.NewPGStockAdjustmentRepo(gormDB)
	takeRepo := repository.NewPGStockTakeRepo(gormDB)
	valRepo := repository.NewPGInventoryValuationRunRepo(gormDB)
	whSvc := service.NewWarehouseService(whRepo, catRepo, itemRepo, balRepo, invTxnRepo, trfRepo, adjRepo, takeRepo, valRepo, grnRepo, svc)
		whH := handler.NewWarehouseHandler(whSvc)
	faSvc := service.NewFAService(faRepo)
	faH := handler.NewFAHandler(faSvc)
	handler.RegisterRoutesWithCompany(r, h, companyH, taxH, cashH, bankH, purchaseH, saleH, whH, faH, authMW, adminMW)
	zap.L().Info("GoTax GL server (PG) starting", zap.String("port", cfg.ServerPort))
		r.Run(cfg.ServerPort)
		return
	}

	zap.L().Info("using in-memory backend")
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
	memFARepo := repository.NewMemoryFARepo()
	faSvcMem := service.NewFAService(memFARepo)
	faHMem := handler.NewFAHandler(faSvcMem)
	handler.RegisterRoutesWithCompany(r, h, companyH, taxH, cashH, bankH, purchaseH, saleH, whH, faHMem, authMW, adminMW)
	zap.L().Info("GoTax GL server (CA) starting", zap.String("port", cfg.ServerPort))
	r.Run(cfg.ServerPort)
}
