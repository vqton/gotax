package service

import (
	"context"
	"testing"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func TestPopulateVAT01(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	decl, err := svc.PopulateVAT01(context.Background(), "comp1", 2026, 7, 100_000_000, 10_000_000, 3_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decl.DeclarationType != domain.DeclTypeGTGT01 {
		t.Errorf("expected DeclTypeGTGT01, got %v", decl.DeclarationType)
	}
	if len(decl.Lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(decl.Lines))
	}

	// Line 10: revenue
	if decl.Lines[0].Amount != 100_000_000 {
		t.Errorf("expected revenue 100M, got %f", decl.Lines[0].Amount)
	}
	// Line 30: outputVAT - inputVAT = 7M
	if decl.Lines[3].Amount != 7_000_000 {
		t.Errorf("expected vatPayable 7M, got %f", decl.Lines[3].Amount)
	}
}

func TestPopulateVAT01_NegativeInputVAT(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	decl, err := svc.PopulateVAT01(context.Background(), "comp1", 2026, 7, 50_000_000, 2_000_000, 5_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// vatPayable should be 0 (input > output)
	if decl.Lines[3].Amount != 0 {
		t.Errorf("expected vatPayable 0, got %f", decl.Lines[3].Amount)
	}
}

func TestPopulateCIT03_3Tier(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	// 20B taxable income → 15% on first 2B + 17% on next 16B + 20% on remaining 2B
	decl, err := svc.PopulateCIT03(context.Background(), "comp1", 2026, 30_000_000_000, 10_000_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decl.DeclarationType != domain.DeclTypeTNDN03 {
		t.Errorf("expected DeclTypeTNDN03, got %v", decl.DeclarationType)
	}
	if len(decl.Lines) != 7 {
		t.Fatalf("expected 7 lines, got %d", len(decl.Lines))
	}

	// Line 30: taxable income = 20B
	if decl.Lines[2].Amount != 20_000_000_000 {
		t.Errorf("expected taxable income 20B, got %f", decl.Lines[2].Amount)
	}
	// Line 31a: 15% on 2B = 300M
	if decl.Lines[3].Amount != 300_000_000 {
		t.Errorf("expected tax15 300M, got %f", decl.Lines[3].Amount)
	}
	// Line 31b: 17% on 16B = 2.72B
	if decl.Lines[4].Amount != 2_720_000_000 {
		t.Errorf("expected tax17 2.72B, got %f", decl.Lines[4].Amount)
	}
	// Line 31c: 20% on 2B = 400M
	if decl.Lines[5].Amount != 400_000_000 {
		t.Errorf("expected tax20 400M, got %f", decl.Lines[5].Amount)
	}
	// Line 40: total = 300M + 2.72B + 400M = 3.42B
	if decl.Lines[6].Amount != 3_420_000_000 {
		t.Errorf("expected totalTax 3.42B, got %f", decl.Lines[6].Amount)
	}
}

func TestPopulateCIT03_SmallIncome(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	// 1B taxable income → 15% only
	decl, err := svc.PopulateCIT03(context.Background(), "comp1", 2026, 5_000_000_000, 4_000_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Line 31a: 15% on 1B = 150M
	if decl.Lines[3].Amount != 150_000_000 {
		t.Errorf("expected tax15 150M, got %f", decl.Lines[3].Amount)
	}
	// Line 31b: 0
	if decl.Lines[4].Amount != 0 {
		t.Errorf("expected tax17 0, got %f", decl.Lines[4].Amount)
	}
	// Line 31c: 0
	if decl.Lines[5].Amount != 0 {
		t.Errorf("expected tax20 0, got %f", decl.Lines[5].Amount)
	}
}

func TestPopulateCIT03_NegativeIncome(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	decl, err := svc.PopulateCIT03(context.Background(), "comp1", 2026, 5_000_000_000, 10_000_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// taxable income = 0
	if decl.Lines[2].Amount != 0 {
		t.Errorf("expected taxable income 0, got %f", decl.Lines[2].Amount)
	}
	if decl.Lines[6].Amount != 0 {
		t.Errorf("expected totalTax 0, got %f", decl.Lines[6].Amount)
	}
}

func TestPopulatePIT05(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	decl, err := svc.PopulatePIT05(context.Background(), "comp1", 2026, 7, 50_000_000, 20_000_000, 5_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decl.DeclarationType != domain.DeclTypeKKTNCN {
		t.Errorf("expected DeclTypeKKTNCN, got %v", decl.DeclarationType)
	}
	// Line 30: taxable income = 50M - 20M - 5M = 25M
	if decl.Lines[3].Amount != 25_000_000 {
		t.Errorf("expected taxable income 25M, got %f", decl.Lines[3].Amount)
	}
}

func TestPopulatePIT05_NegativeIncome(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	decl, err := svc.PopulatePIT05(context.Background(), "comp1", 2026, 7, 10_000_000, 20_000_000, 5_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decl.Lines[3].Amount != 0 {
		t.Errorf("expected taxable income 0, got %f", decl.Lines[3].Amount)
	}
}

func TestCreateDeclaration_Validation(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	// Missing company ID
	err := svc.CreateDeclaration(context.Background(), &domain.TaxDeclaration{})
	if err != domain.ErrCompanyIDRequired {
		t.Errorf("expected ErrCompanyIDRequired, got %v", err)
	}

	// Invalid declaration type
	err = svc.CreateDeclaration(context.Background(), &domain.TaxDeclaration{
		CompanyID:       "comp1",
		DeclarationType: "INVALID",
	})
	if err != domain.ErrDeclarationTypeInvalid {
		t.Errorf("expected ErrDeclarationTypeInvalid, got %v", err)
	}
}

func TestCreateDeclaration_Success(t *testing.T) {
	repo := repository.NewMemoryTaxRepo()
	svc := NewTaxDeclarationService(repo)

	err := svc.CreateDeclaration(context.Background(), &domain.TaxDeclaration{
		CompanyID:       "comp1",
		DeclarationType: domain.DeclTypeGTGT01,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify persisted
	filter := domain.TaxDeclarationFilter{CompanyID: "comp1"}
	decls, err := repo.GetDeclarations(context.Background(), filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(decls) != 1 {
		t.Fatalf("expected 1 declaration, got %d", len(decls))
	}
}
