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
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"time"

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
	"gotax/internal/einvoice"
	"gotax/internal/gdt"
	"gotax/internal/handler"
	gotaxi18n "gotax/internal/i18n"
	"gotax/internal/logger"
	"gotax/internal/repository"
	"gotax/internal/service"
	"gotax/internal/web"
	"gotax/internal/xmldsig"
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
// Static assets: no-cache forces revalidation so browser never serves stale
// compiled CSS/JS after a Tailwind rebuild or JS edit (last build at 1287819
// was served from Brave disk cache, breaking the COA layout).
assets := r.Group("/assets")
assets.Use(func(c *gin.Context) { c.Header("Cache-Control", "no-cache"); c.Next() })
assets.Static("", "./web/static")
payroll := r.Group("/payroll")
payroll.Use(func(c *gin.Context) { c.Header("Cache-Control", "no-cache"); c.Next() })
payroll.Static("", "./web/payroll")
// /app served by internal/web catch-all: converted pages render templates,
// unconverted pages fall back to static files.

r.GET("/ping", func(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GoTax GL Server"})
})
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

r.GET("/", func(c *gin.Context) {
	c.Redirect(http.StatusFound, "/app/dashboard.html")
})
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
		exportSvc := service.NewExportService(jeRepo, accRepo, perRepo)
		h := handler.NewHandlerWithExport(svc, exportSvc)

		companyRepo := repository.NewPGCompanyRepo(gormDB)
		companySvc := service.NewCompanyService(companyRepo)
		companyH := handler.NewCompanyHandler(companySvc)
		deptH := handler.NewDepartmentHandler(companySvc)

		authMW := handler.AuthMiddleware()
		adminMW := handler.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleChiefAccountant)

		taxRepo := repository.NewPGTaxRepo(gormDB)
		taxSvc := service.NewTaxService(taxRepo, jeRepo, companyRepo, newGDTClient(), newEInvoiceSigner())
		taxDeclSvc := service.NewTaxDeclarationService(taxRepo)
		taxH := handler.NewTaxHandler(taxSvc, taxDeclSvc, auditRepo)
		cashH := handler.NewCashHandler(svc)
		bankRepo := repository.NewPGBankRepo(gormDB)
		bankSvc := service.NewBankService(bankRepo)
		bankH := handler.NewBankHandler(bankSvc, companySvc)
		bankImportRepo := repository.NewPGBankImportRepo(gormDB)
		bankImportSvc := service.NewBankImportService(bankImportRepo)
		bankImportH := handler.NewBankImportHandler(bankImportSvc)
		invoiceBookRepo := repository.NewPGInvoiceBookRepo(gormDB)
		invoiceBookSvc := service.NewInvoiceBookService(invoiceBookRepo)
		invoiceBookH := handler.NewInvoiceBookHandler(invoiceBookSvc)

	supplierRepo := repository.NewPGSupplierRepo(gormDB)
	purchaseOrderRepo := repository.NewPGPurchaseRepo(gormDB)
	grnRepo := repository.NewPGGRNRepo(gormDB)
	supplierInvoiceRepo := repository.NewPGSupplierInvoiceRepo(gormDB)
	apTxnRepo := repository.NewPGAPTransactionRepo(gormDB)
	costAllocRepo := repository.NewPGCostAllocationRepo(gormDB)
	provisionRepo := repository.NewPGDoubtfulDebtProvisionRepo(gormDB)
	requisitionRepo := repository.NewPGRequisitionRepo(gormDB)
	fxRevaluationRepo := repository.NewPGFXRevaluationRepo(gormDB)
	purchaseSvc := service.NewPurchaseService(supplierRepo, purchaseOrderRepo, grnRepo, supplierInvoiceRepo, apTxnRepo, costAllocRepo, provisionRepo, requisitionRepo, fxRevaluationRepo, svc)
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
	ccdcCatRepo := repository.NewPGCCDCCategoryRepo(gormDB)
	ccdcItemRepo := repository.NewPGToolEquipmentRepo(gormDB)
	ccdcSvc := service.NewCCDCService(ccdcCatRepo, ccdcItemRepo)
	ccdcH := handler.NewCCDCHandler(ccdcSvc)
	pwRepo := repository.NewPGPayrollRepo(gormDB)
	pwSvc := service.NewPayrollService(pwRepo, companyRepo)
	pwH := handler.NewPayrollHandler(pwSvc)
	recRepo := repository.NewPGRecurringEntryRepo(gormDB)
	recSvc := service.NewRecurringService(recRepo, jeRepo)
	recH := handler.NewRecurringHandler(recSvc)
	budRepo := repository.NewPGBudgetRepo(gormDB)
	budSvc := service.NewBudgetService(budRepo, jeRepo)
	budH := handler.NewBudgetHandler(budSvc)
	notifRepo := repository.NewPGNotificationRepo(gormDB)
	notifSvc := service.NewNotificationService(notifRepo)
	notifH := handler.NewNotificationHandler(notifSvc)
	costRepo := repository.NewPGCostCenterRepo(gormDB)
	costSvc := service.NewCostCenterService(costRepo)
	costH := handler.NewCostCenterHandler(costSvc)
	keeperRepo := repository.NewPGWarehouseKeeperRepo(gormDB)
	keeperSvc := service.NewWarehouseKeeperService(keeperRepo, whRepo, itemRepo)
	keeperH := handler.NewWarehouseKeeperHandler(keeperSvc)
	costObjectRepo := repository.NewPGCostObjectRepo(gormDB)
	costPoolRepo := repository.NewPGCostPoolRepo(gormDB)
	costPoolLineRepo := repository.NewPGCostPoolLineRepo(gormDB)
	costObjectSvc := service.NewCostObjectService(costObjectRepo)
	costPoolSvc := service.NewCostPoolService(costPoolRepo, costPoolLineRepo)
	costObjectH := handler.NewCostObjectHandler(costObjectSvc)
	costPoolH := handler.NewCostPoolHandler(costPoolSvc)
	costingPeriodRepo := repository.NewPGCostingPeriodRepo(gormDB)
	costingResultRepo := repository.NewPGCostingResultRepo(gormDB)
	costingResultLineRepo := repository.NewPGCostingResultLineRepo(gormDB)
	costingPeriodSvc := service.NewCostingPeriodService(costingPeriodRepo)
	costingEngine := service.NewCostingEngine(costingPeriodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, costingResultRepo, costingResultLineRepo)
	costingJESvc := service.NewCostingJEService(costingPeriodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, costingResultRepo, costingResultLineRepo, svc)
	costReportSvc := service.NewCostReportService(costingResultRepo, costingResultLineRepo, costObjectRepo)
	costingH := handler.NewCostingHandler(costingPeriodSvc, costingEngine, costingJESvc)
	costReportH := handler.NewCostReportHandler(costReportSvc)
	handler.RegisterCostReportRoutes(r, costReportH, authMW)
	systemOptRepo := repository.NewPGSystemOptionRepo(gormDB)
	systemOptSvc := service.NewSystemOptionService(systemOptRepo)
	systemOptH := handler.NewSystemOptionHandler(systemOptSvc)
	numRuleRepo := repository.NewPGNumberingRuleRepo(gormDB)
	numRuleSvc := service.NewNumberingRuleService(numRuleRepo)
	numRuleH := handler.NewNumberingRuleHandler(numRuleSvc)
	handler.RegisterSystemOptionRoutes(r, systemOptH, authMW)
	handler.RegisterNumberingRuleRoutes(r, numRuleH, authMW)
	fiscalYearSvc := service.NewFiscalYearService(perRepo, systemOptRepo)
	fiscalYearH := handler.NewFiscalYearHandler(fiscalYearSvc)
	handler.RegisterFiscalYearRoutes(r, fiscalYearH, authMW)
	reportOptSvc := service.NewReportOptionService(systemOptRepo)
	reportOptH := handler.NewReportOptionHandler(reportOptSvc)
	handler.RegisterReportOptionRoutes(r, reportOptH, authMW)
	backupRepo := repository.NewPGBackupRepo(gormDB)
	backupSvc := service.NewBackupService(backupRepo, cfg.DatabaseURL, "/tmp/gotax-backups")
	backupH := handler.NewBackupHandler(backupSvc)
	handler.RegisterBackupRoutes(r, backupH, authMW)
	contractRepo := repository.NewPGContractRepo(gormDB)
	contractPaymentRepo := repository.NewPGContractPaymentRepo(gormDB)
	contractSvc := service.NewContractService(contractRepo, contractPaymentRepo)
	contractH := handler.NewContractHandler(contractSvc)
	handler.RegisterContractRoutes(r, contractH, authMW)
	handler.RegisterRoutesWithCompany(r, h, companyH, taxH, cashH, bankH, purchaseH, saleH, whH, faH, pwH, recH, budH, ccdcH, costH, keeperH, costObjectH, costPoolH, costingH, authMW, adminMW)
	handler.RegisterDepartmentRoutes(r, deptH, authMW)
	handler.RegisterBankImportRoutes(r, bankImportH, authMW)
	handler.RegisterInvoiceBookRoutes(r, invoiceBookH, authMW)
	handler.RegisterNotificationRoutes(r, notifH, authMW)
	priceListRepo := repository.NewPGPriceListRepo(gormDB)
	priceListSvc := service.NewPriceListService(priceListRepo)
	priceListH := handler.NewPriceListHandler(priceListSvc)
	handler.RegisterPriceListRoutes(r, priceListH, authMW)
	finAnalysisSvc := service.NewFinancialAnalysisService(jeRepo, budRepo, companyRepo)
	finAnalysisH := handler.NewFinancialAnalysisHandler(finAnalysisSvc)
	handler.RegisterFinancialAnalysisRoutes(r, finAnalysisH, authMW)
	zap.L().Info("GoTax GL server (PG) starting", zap.String("port", cfg.ServerPort))
		webSrv, err := web.NewServer([]string{"dashboard", "users", "journal-entries"})
		if err != nil {
			zap.L().Fatal("init web templates", zap.Error(err))
		}
		webDeps := web.Deps{Svc: svc}
		webSrv.RegisterPages(r, web.NewPages(webDeps), webSrv.NewActions(webDeps))
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

	// Seed admin user (password: Admin@123456!)
	if existing, _ := userRepo.GetByUsername(ctx, "admin"); existing == nil {
		admin := &domain.User{
			Username: "admin",
			FullName: "System Admin",
			Email:    "admin@gotax.vn",
			Role:     domain.UserRoleAdmin,
			IsActive: true,
		}
		if err := svc.CreateUser(ctx, admin, "Admin@123456!"); err != nil {
			zap.L().Warn("seed admin user", zap.Error(err))
		} else {
			zap.L().Info("admin user seeded", zap.String("username", "admin"))
		}
	}

	exportSvc := service.NewExportService(jeRepo, accRepo, perRepo)
	h := handler.NewHandlerWithExport(svc, exportSvc)

	companyRepo := repository.NewMemoryCompanyRepo()
	companySvc := service.NewCompanyService(companyRepo)
	companyH := handler.NewCompanyHandler(companySvc)
	deptHMem := handler.NewDepartmentHandler(companySvc)

	authMW := handler.AuthMiddleware()
	adminMW := handler.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleChiefAccountant)

	taxRepo := repository.NewMemoryTaxRepo()
	taxSvc := service.NewTaxService(taxRepo, jeRepo, companyRepo, newGDTClient(), newEInvoiceSigner())
	taxDeclSvc := service.NewTaxDeclarationService(taxRepo)
	taxH := handler.NewTaxHandler(taxSvc, taxDeclSvc, auditRepo)
	cashH := handler.NewCashHandler(svc)
	bankRepo := repository.NewMemoryBankRepo()
	bankSvc := service.NewBankService(bankRepo)
	bankH := handler.NewBankHandler(bankSvc, companySvc)
	bankImportRepo := repository.NewMemoryBankImportRepo()
	bankImportSvc := service.NewBankImportService(bankImportRepo)
	bankImportHMem := handler.NewBankImportHandler(bankImportSvc)
	invoiceBookRepoMem := repository.NewMemoryInvoiceBookRepo()
	invoiceBookSvcMem := service.NewInvoiceBookService(invoiceBookRepoMem)
	invoiceBookHMem := handler.NewInvoiceBookHandler(invoiceBookSvcMem)

	memPurchaseRepo := repository.NewMemoryPurchaseRepo()
	purchaseSvc := service.NewPurchaseService(memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, memPurchaseRepo, svc)
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
	memPWRepo := repository.NewMemoryPayrollRepo()
	pwSvcMem := service.NewPayrollService(memPWRepo, companyRepo)
	pwHMem := handler.NewPayrollHandler(pwSvcMem)
	memRecRepo := repository.NewMemoryRecurringEntryRepo()
	recSvcMem := service.NewRecurringService(memRecRepo, jeRepo)
	recHMem := handler.NewRecurringHandler(recSvcMem)
	memBudRepo := repository.NewMemoryBudgetRepo()
	budSvcMem := service.NewBudgetService(memBudRepo, jeRepo)
	budHMem := handler.NewBudgetHandler(budSvcMem)
	memCCDCCatRepo := repository.NewMemoryCCDCCategoryRepo()
	memCCDCItemRepo := repository.NewMemoryCCDCItemRepo()
	ccdcSvcMem := service.NewCCDCService(memCCDCCatRepo, memCCDCItemRepo)
	ccdcHMem := handler.NewCCDCHandler(ccdcSvcMem)
	memNotifRepo := repository.NewMemoryNotificationRepo()
	notifSvcMem := service.NewNotificationService(memNotifRepo)
	notifHMem := handler.NewNotificationHandler(notifSvcMem)
	memCostRepo := repository.NewMemoryCostCenterRepo()
	costSvcMem := service.NewCostCenterService(memCostRepo)
	costHMem := handler.NewCostCenterHandler(costSvcMem)
	memKeeperRepo := repository.NewMemoryWarehouseKeeperRepo()
	keeperSvcMem := service.NewWarehouseKeeperService(memKeeperRepo, memWhRepo, memItemRepo)
	keeperHMem := handler.NewWarehouseKeeperHandler(keeperSvcMem)
	costObjectRepoMem := repository.NewMemoryCostObjectRepo()
	costPoolRepoMem := repository.NewMemoryCostPoolRepo()
	costPoolLineRepoMem := repository.NewMemoryCostPoolLineRepo()
	costObjectSvcMem := service.NewCostObjectService(costObjectRepoMem)
	costPoolSvcMem := service.NewCostPoolService(costPoolRepoMem, costPoolLineRepoMem)
	costObjectHMem := handler.NewCostObjectHandler(costObjectSvcMem)
	costPoolHMem := handler.NewCostPoolHandler(costPoolSvcMem)
	costingPeriodRepoMem := repository.NewMemoryCostingPeriodRepo()
	costingResultRepoMem := repository.NewMemoryCostingResultRepo()
	costingResultLineRepoMem := repository.NewMemoryCostingResultLineRepo()
	costingPeriodSvcMem := service.NewCostingPeriodService(costingPeriodRepoMem)
	costingEngineMem := service.NewCostingEngine(costingPeriodRepoMem, costObjectRepoMem, costPoolRepoMem, costPoolLineRepoMem, costingResultRepoMem, costingResultLineRepoMem)
	costingJESvcMem := service.NewCostingJEService(costingPeriodRepoMem, costObjectRepoMem, costPoolRepoMem, costPoolLineRepoMem, costingResultRepoMem, costingResultLineRepoMem, svc)
	costReportSvcMem := service.NewCostReportService(costingResultRepoMem, costingResultLineRepoMem, costObjectRepoMem)
	costingHMem := handler.NewCostingHandler(costingPeriodSvcMem, costingEngineMem, costingJESvcMem)
	costReportHMem := handler.NewCostReportHandler(costReportSvcMem)
	handler.RegisterCostReportRoutes(r, costReportHMem, authMW)
	systemOptRepoMem := repository.NewMemorySystemOptionRepo()
	systemOptSvcMem := service.NewSystemOptionService(systemOptRepoMem)
	systemOptHMem := handler.NewSystemOptionHandler(systemOptSvcMem)
	numRuleRepoMem := repository.NewMemoryNumberingRuleRepo()
	numRuleSvcMem := service.NewNumberingRuleService(numRuleRepoMem)
	numRuleHMem := handler.NewNumberingRuleHandler(numRuleSvcMem)
	handler.RegisterSystemOptionRoutes(r, systemOptHMem, authMW)
	handler.RegisterNumberingRuleRoutes(r, numRuleHMem, authMW)
	fiscalYearSvcMem := service.NewFiscalYearService(perRepo, systemOptRepoMem)
	fiscalYearHMem := handler.NewFiscalYearHandler(fiscalYearSvcMem)
	handler.RegisterFiscalYearRoutes(r, fiscalYearHMem, authMW)
	reportOptSvcMem := service.NewReportOptionService(systemOptRepoMem)
	reportOptHMem := handler.NewReportOptionHandler(reportOptSvcMem)
	handler.RegisterReportOptionRoutes(r, reportOptHMem, authMW)
	backupRepoMem := repository.NewMemoryBackupRepo()
	backupSvcMem := service.NewBackupService(backupRepoMem, "", "")
	backupHMem := handler.NewBackupHandler(backupSvcMem)
	handler.RegisterBackupRoutes(r, backupHMem, authMW)
	contractRepoMem := repository.NewMemoryContractRepo()
	contractPaymentRepoMem := repository.NewMemoryContractPaymentRepo()
	contractSvcMem := service.NewContractService(contractRepoMem, contractPaymentRepoMem)
	contractHMem := handler.NewContractHandler(contractSvcMem)
	handler.RegisterContractRoutes(r, contractHMem, authMW)
	handler.RegisterRoutesWithCompany(r, h, companyH, taxH, cashH, bankH, purchaseH, saleH, whH, faHMem, pwHMem, recHMem, budHMem, ccdcHMem, costHMem, keeperHMem, costObjectHMem, costPoolHMem, costingHMem, authMW, adminMW)
	handler.RegisterDepartmentRoutes(r, deptHMem, authMW)
	handler.RegisterBankImportRoutes(r, bankImportHMem, authMW)
	handler.RegisterInvoiceBookRoutes(r, invoiceBookHMem, authMW)
	handler.RegisterNotificationRoutes(r, notifHMem, authMW)
	priceListRepoMem := repository.NewMemoryPriceListRepo()
	priceListSvcMem := service.NewPriceListService(priceListRepoMem)
	priceListHMem := handler.NewPriceListHandler(priceListSvcMem)
	handler.RegisterPriceListRoutes(r, priceListHMem, authMW)
	finAnalysisSvcMem := service.NewFinancialAnalysisService(jeRepo, memBudRepo, companyRepo)
	finAnalysisHMem := handler.NewFinancialAnalysisHandler(finAnalysisSvcMem)
	handler.RegisterFinancialAnalysisRoutes(r, finAnalysisHMem, authMW)
	zap.L().Info("GoTax GL server (CA) starting", zap.String("port", cfg.ServerPort))
	webSrv, err := web.NewServer([]string{"dashboard", "users", "journal-entries"})
	if err != nil {
		zap.L().Fatal("init web templates", zap.Error(err))
	}
	webDeps := web.Deps{Svc: svc}
	webSrv.RegisterPages(r, web.NewPages(webDeps), webSrv.NewActions(webDeps))
	r.Run(cfg.ServerPort)
}

// newGDTClient returns a GDT API client when GDT_BASE_URL is set, else nil
// (e-invoice issue/status calls then return ErrGDTUnavailable).
func newGDTClient() service.GDTClient {
	base := os.Getenv("GDT_BASE_URL")
	if base == "" {
		return nil
	}
	opts := []gdt.Option{gdt.WithToken(os.Getenv("GDT_TOKEN"))}
	// mTLS: optional client cert for GDT authentication
	certPath := os.Getenv("GDT_CLIENT_CERT")
	if certPath != "" {
		opts = append(opts, gdt.WithClientCert(certPath, os.Getenv("GDT_CA_BUNDLE")))
	}
	opts = append(opts, gdt.WithLogger(func(msg string, args ...any) {
		zap.S().Named("gdt").Infof(msg, args...)
	}))
	c, err := gdt.New(base, opts...)
	if err != nil {
		zap.L().Fatal("invalid GDT_BASE_URL", zap.Error(err))
	}
	return c
}

// newEInvoiceSigner builds the TXML signer. Requires EINVOICE_SIGNING_KEY
// (RSA private key PEM). Without it, an ephemeral key is generated at
// startup — signatures then cannot be verified after a restart, so set it
// in production.
func newEInvoiceSigner() service.TXMLSigner {
	var (
		key *rsa.PrivateKey
		err error
	)
	pemBytes := []byte(os.Getenv("EINVOICE_SIGNING_KEY"))
	if len(pemBytes) > 0 {
		key, err = xmldsig.ParsePrivateKeyPEM(pemBytes)
		if err != nil {
			zap.L().Fatal("EINVOICE_SIGNING_KEY is not a valid RSA private key", zap.Error(err))
		}
	} else {
		key, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			zap.L().Fatal("generate ephemeral e-invoice signing key", zap.Error(err))
		}
		zap.L().Warn("EINVOICE_SIGNING_KEY unset — using ephemeral signing key; signatures invalid after restart")
	}
	serial := os.Getenv("EINVOICE_CERT_SERIAL")
	if serial == "" {
		serial = "GOTAX-DEV"
	}
	return einvoice.NewPEMSigner(key, serial, time.Now)
}
