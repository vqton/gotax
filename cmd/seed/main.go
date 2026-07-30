package main

import (
	"context"
	"log"
	"os"
	"time"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func main() {
	userRepo := repository.NewMemoryUserRepo()

	svc := service.NewService(
		repository.NewMemoryAccountRepo(),
		repository.NewMemoryJournalRepo(),
		repository.NewMemoryPeriodRepo(),
		userRepo,
		repository.NewMemoryAuditLogRepo(),
		repository.NewMemoryExchangeRateRepo(),
		repository.NewMemoryClosingTemplateRepo(),
		repository.NewMemoryApprovalRepo(),
		repository.NewMemoryAccountVersionRepo(),
		repository.NewMemoryAccountMappingRepo(),
		repository.NewMemoryAccountAnalysisRepo(),
		repository.NewMemoryIFRSMappingRepo(),
		repository.NewMemoryRefreshTokenRepo(),
		repository.NewMemoryPasswordResetTokenRepo(),
		repository.NewMemoryOpeningBalanceRepo(),
		repository.NewMemoryCashRepo(),
	)

	ctx := context.Background()
	existing, _ := userRepo.GetByUsername(ctx, "admin")
	if existing != nil {
		log.Println("admin already exists, skipping")
		os.Exit(0)
	}

	u := &domain.User{
		Username:    "admin",
		FullName:    "System Admin",
		Email:       "admin@gotax.vn",
		Role:        domain.UserRoleAdmin,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := svc.CreateUser(ctx, u, "Admin@123456!"); err != nil {
		log.Fatalf("seed failed: %v", err)
	}
	log.Println("seeded admin — username: admin | password: Admin@123456")
}
