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

const usage = `gotax-set-role — change user role

Usage:
  go run ./cmd/set-role --username <name> --role <role>
  go run ./cmd/set-role --id <user-id>    --role <role>

Either --username or --id is required.

Roles:
  admin            full system access
  chief_accountant GL approve, close periods, issue reports
  accountant       CRUD journals, COA, reports — no admin
  viewer           read-only
`

func main() {
	username := flag.String("username", "", "Login name of user")
	userID := flag.String("id", "", "User ID (UUID)")
	role := flag.String("role", "", "admin | chief_accountant | accountant | viewer")
	flag.Parse()

	if *username == "" && *userID == "" {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	if *role == "" {
		fmt.Fprint(os.Stderr, "error: --role is required\n"+usage)
		os.Exit(2)
	}
	ur, ok := validRole(*role)
	if !ok {
		fmt.Fprintf(os.Stderr, "invalid role %q\n", *role)
		os.Exit(2)
	}

	svc := mustService()
	ctx := context.Background()

	var targetID string
	if *userID != "" {
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

	old, _ := svc.GetUser(ctx, targetID)
	if err := svc.UpdateUserRole(ctx, targetID, ur); err != nil {
		log.Fatalf("update role failed: %v", err)
	}
	oldRoleStr := ""
	if old != nil {
		oldRoleStr = string(old.Role)
	}
	fmt.Printf("OK — user %s role %s -> %s\n",
		alias(*username, *userID),
		oldRoleStr, ur)
}

func validRole(s string) (domain.UserRole, bool) {
	switch domain.UserRole(s) {
	case domain.UserRoleAdmin, domain.UserRoleChiefAccountant, domain.UserRoleAccountant, domain.UserRoleViewer:
		return domain.UserRole(s), true
	}
	return "", false
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
