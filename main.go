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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "GoTax GL Server"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		svc := service.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil, nil, nil, obRepo, cashRepo)
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

		handler.RegisterRoutesWithCompany(r, h, companyH, taxH, cashH, authMW, adminMW)
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

	handler.RegisterRoutesWithCompany(r, h, companyH, taxH, cashH, authMW, adminMW)
	log.Println("GoTax GL server (CA) starting on :8080")
	r.Run(":8080")
}
