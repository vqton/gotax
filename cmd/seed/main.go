package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

const seedAdminPassword = "Admin@123456!"

func main() {
	userRepo := repository.NewMemoryUserRepo()

	svc := service.NewService(
		mustRepo(repository.NewMemoryAccountRepo),
		mustRepo(repository.NewMemoryJournalRepo),
		mustRepo(repository.NewMemoryPeriodRepo),
		userRepo,
		mustRepo(repository.NewMemoryAuditLogRepo),
		mustRepo(repository.NewMemoryExchangeRateRepo),
		mustRepo(repository.NewMemoryClosingTemplateRepo),
		mustRepo(repository.NewMemoryApprovalRepo),
		mustRepo(repository.NewMemoryAccountVersionRepo),
		mustRepo(repository.NewMemoryAccountMappingRepo),
		mustRepo(repository.NewMemoryAccountAnalysisRepo),
		mustRepo(repository.NewMemoryIFRSMappingRepo),
		mustRepo(repository.NewMemoryRefreshTokenRepo),
		mustRepo(repository.NewMemoryPasswordResetTokenRepo),
		mustRepo(repository.NewMemoryOpeningBalanceRepo),
		mustRepo(repository.NewMemoryCashRepo),
	)

	ctx := context.Background()
	existing, _ := userRepo.GetByUsername(ctx, "admin")
	if existing != nil {
		if err := svc.AdminResetPassword(ctx, existing.ID, seedAdminPassword); err != nil {
			log.Fatalf("refresh admin password: %v", err)
		}
		log.Println("admin password refreshed (in-memory) — username: admin | password: Admin@123456!")
		os.Exit(0)
	}

	u := &domain.User{
		Username:  "admin",
		FullName:  "System Admin",
		Email:     "admin@gotax.vn",
		Role:      domain.UserRoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := svc.CreateUser(ctx, u, seedAdminPassword); err != nil {
		log.Fatalf("create admin: %v", err)
	}
	log.Println("admin created (in-memory) — username: admin | password: Admin@123456!")
}

func mustRepo[T any](fn func() T) T {
	return fn()
}
