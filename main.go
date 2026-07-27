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
	"gotax/internal/gl"
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
	gl.SetJWTSecret(secret)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "GoTax GL Server"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		log.Println("using PostgreSQL backend")
		cfg := gl.DefaultPGConfig()
		cfg.DSN = dsn

		pool, err := gl.NewPool(ctx, cfg)
		if err != nil {
			log.Fatalf("connect to PostgreSQL: %v", err)
		}
		defer pool.Close()

		if err := gl.RunMigrations(ctx, pool); err != nil {
			log.Fatalf("run migrations: %v", err)
		}

		accRepo := gl.NewPGAccountRepo(pool)
		jeRepo := gl.NewPGJournalRepo(pool)
		perRepo := gl.NewPGPeriodRepo(pool)
		userRepo := gl.NewPGUserRepo(pool)
		auditRepo := gl.NewPGAuditLogRepo(pool)
		rateRepo := gl.NewPGExchangeRateRepo(pool)
		templateRepo := gl.NewPGClosingTemplateRepo(pool)

		svc := gl.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil, nil, nil)
		handler := gl.NewHandler(svc)

		authMW := gl.AuthMiddleware()
		adminMW := gl.RoleMiddleware(gl.UserRoleAdmin, gl.UserRoleChiefAccountant)

		gl.RegisterRoutes(r, handler, authMW, adminMW)
		log.Println("GoTax GL server (PG) starting on :8080")
		r.Run(":8080")
		return
	}

	log.Println("using in-memory backend (no DATABASE_URL set)")
	accRepo := gl.NewMemoryAccountRepo()
	jeRepo := gl.NewMemoryJournalRepo()
	jeRepo.SetAccounts(accRepo.Accounts())
	perRepo := gl.NewMemoryPeriodRepo()
	userRepo := gl.NewMemoryUserRepo()
	auditRepo := gl.NewMemoryAuditLogRepo()
	rateRepo := gl.NewMemoryExchangeRateRepo()
	templateRepo := gl.NewMemoryClosingTemplateRepo()
	refreshRepo := gl.NewMemoryRefreshTokenRepo()
	resetRepo := gl.NewMemoryPasswordResetTokenRepo()

	svc := gl.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil, refreshRepo, resetRepo)
	handler := gl.NewHandler(svc)

	authMW := gl.AuthMiddleware()
	adminMW := gl.RoleMiddleware(gl.UserRoleAdmin, gl.UserRoleChiefAccountant)

	gl.RegisterRoutes(r, handler, authMW, adminMW)
	log.Println("GoTax GL server (mem) starting on :8080")
	r.Run(":8080")
}
