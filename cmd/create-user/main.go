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

const usage = `gotax-create-user — create a new user

Usage:
  go run ./cmd/create-user --username <name> --password <pass> --role <role> [--full-name "..."] [--email "..."]

Required flags:
  --username   Login name (unique, no spaces)
  --password   Min 12 chars with upper, lower, digit, special
  --role       One of: admin | chief_accountant | accountant | viewer

Optional flags:
  --full-name  Display name
  --email      Email address (used for password reset flow)

Roles map to backend enum UserRole*:
  admin            full system access
  chief_accountant GL approve, close periods, issue reports
  accountant       CRUD journals, COA, reports — no admin
  viewer           read-only
`

func main() {
	username := flag.String("username", "", "Login name")
	password := flag.String("password", "", "Account password")
	role := flag.String("role", "", "admin | chief_accountant | accountant | viewer")
	fullName := flag.String("full-name", "", "Display name")
	email := flag.String("email", "", "Email address")
	flag.Parse()

	if *username == "" || *password == "" || *role == "" {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	ur, ok := validRole(*role)
	if !ok {
		fmt.Fprintf(os.Stderr, "invalid role %q — must be admin | chief_accountant | accountant | viewer\n", *role)
		os.Exit(2)
	}

	svc := mustService()
	ctx := context.Background()

	u := &domain.User{
		Username: *username,
		FullName: *fullName,
		Email:    *email,
		Role:     ur,
		IsActive: true,
	}
	if err := svc.CreateUser(ctx, u, *password); err != nil {
		log.Fatalf("create user failed: %v", err)
	}
	fmt.Printf("OK — user %s (id=%s) role=%s created\n", u.Username, u.ID, ur)
}

func validRole(s string) (domain.UserRole, bool) {
	switch domain.UserRole(s) {
	case domain.UserRoleAdmin, domain.UserRoleChiefAccountant, domain.UserRoleAccountant, domain.UserRoleViewer:
		return domain.UserRole(s), true
	}
	return "", false
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
