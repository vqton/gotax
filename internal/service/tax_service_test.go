package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func newTaxTestSvc() (*taxService, domain.TaxRepository) {
	repo := repository.NewMemoryTaxRepo()
	return NewTaxService(repo).(*taxService), repo
}

// ─── A1: Rate Resolver ──────────────────────────────────────────────────

func TestResolveRate_ByRateCodeSuffix(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 20,
		EffectiveFrom: "2025-01-01", IsActive: true,
	}))

	rate, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "STANDARD", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 20.0, rate.RateValue)
}

func TestResolveRate_ByApplicableToField(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "VAT8", TaxType: domain.TaxTypeVAT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 8,
		ApplicableTo: "REDUCED", EffectiveFrom: "2025-07-01", IsActive: true,
	}))

	rate, err := svc.resolveRate(ctx, domain.TaxTypeVAT, "REDUCED", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 8.0, rate.RateValue)
}

func TestResolveRate_RespectsEffectiveWindow(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	// legacy rate, superseded on 2026-01-01
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD_LEGACY", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 20,
		EffectiveFrom: "2020-01-01", EffectiveTo: "2025-12-31", IsActive: true,
	}))
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 17,
		EffectiveFrom: "2026-01-01", IsActive: true,
	}))

	rate, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "STANDARD", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 17.0, rate.RateValue)

	legacy, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "STANDARD", "2025-06-01")
	require.NoError(t, err)
	assert.Equal(t, 20.0, legacy.RateValue)
}

func TestResolveRate_SkipsInactive(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 25,
		EffectiveFrom: "2025-01-01", IsActive: false,
	}))

	rate, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "STANDARD", "2026-06-01")
	require.NoError(t, err)
	// statutory fallback for CIT STANDARD
	assert.Equal(t, 20.0, rate.RateValue)
}

func TestResolveRate_StatutoryFallback(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()

	rate, err := svc.resolveRate(ctx, domain.TaxTypeVAT, "STANDARD", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 10.0, rate.RateValue)

	cit, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "SMALL", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 17.0, cit.RateValue)
}
