package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

const usage = `gotax-reset-password — admin force-reset a user password

Usage:
  go run ./cmd/reset-password --username <name> --new-password <pass>
  go run ./cmd/reset-password --id <user-id>    --new-password <pass>

Either --username or --id is required.

Password rules:
  Min 12 chars, at least one upper, lower, digit, special char
`

func main() {
	username := flag.String("username", "", "Login name")
	userID := flag.String("id", "", "User ID (UUID)")
	newPw := flag.String("new-password", "", "New password")
	flag.Parse()

	if (username == nil || *username == "") && (userID == nil || *userID == "") {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	if newPw == nil || *newPw == "" {
		fmt.Fprint(os.Stderr, "error: --new-password is required\n"+usage)
		os.Exit(2)
	}
	if err := domain.ValidatePassword(*newPw); err != nil {
		log.Fatalf("weak password: %v", err)
	}

	svc := mustService()
	ctx := context.Background()

	var targetID string
	if userID != nil && *userID != "" {
		targetID = *userID
	} else {
		users, err := svc.ListUsers(ctx)
		if err != nil {
			log.Fatalf("list users: %v", err)
		}
		for _, u := range users {
			if u.Username == *username {
				targetID = u.ID
				break
			}
		}
		if targetID == "" {
			log.Fatalf("user %q not found", *username)
		}
	}

	if err := svc.AdminResetPassword(ctx, targetID, *newPw); err != nil {
		log.Fatalf("reset password failed: %v", err)
	}
	fmt.Printf("OK — password reset for %s (id=%s)\n",
		alias(*username, *userID), targetID)
}

func alias(uname, uid string) string {
	if uname != "" {
		return uname
	}
	return uid
}

func mustService() service.Service {
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
	return service.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)
}
