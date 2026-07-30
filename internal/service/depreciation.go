package service

import (
	"context"
	"fmt"
	"math"

	"gotax/internal/domain"
)

type DepreciationEngine struct {
	faRepo domain.FARepository
}

func NewDepreciationEngine(faRepo domain.FARepository) *DepreciationEngine {
	return &DepreciationEngine{faRepo: faRepo}
}

type DepreciationResult struct {
	AssetID    string
	Amount     float64
	AccAfter   float64
	CarryAfter float64
}

func (e *DepreciationEngine) Calculate(ctx context.Context, asset *domain.FixedAsset, year, month int) (*DepreciationResult, error) {
	switch asset.DepreciationMethod {
	case domain.DepStraightLine:
		return e.straightLine(asset, year, month)
	case domain.DepDecliningBalance:
		return e.decliningBalance(asset, year, month)
	case domain.DepProductionBased:
		return nil, fmt.Errorf("production-based depreciation requires production volume input")
	default:
		return nil, domain.ErrFADepreciationMethodInvalid
	}
}

func (e *DepreciationEngine) RunForAsset(ctx context.Context, asset *domain.FixedAsset, input domain.FARunDepreciationInput) (*DepreciationResult, error) {
	exists, err := e.faRepo.DepreciationExistsForPeriod(ctx, asset.ID, input.PeriodID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrFADepreciationExists
	}

	result, err := e.Calculate(ctx, asset, input.Year, input.Month)
	if err != nil {
		return nil, err
	}

	entry := &domain.DepreciationEntry{
		CompanyID:           input.CompanyID,
		FixedAssetID:        asset.ID,
		PeriodID:            input.PeriodID,
		PeriodYear:          input.Year,
		PeriodMonth:         input.Month,
		DepreciationAmount:  result.Amount,
		AccumulatedAfter:    result.AccAfter,
		CarryingAmountAfter: result.CarryAfter,
		CreatedBy:           input.CreatedBy,
	}

	if err := e.faRepo.CreateDepreciationEntry(ctx, entry); err != nil {
		return nil, err
	}

	asset.AccumulatedDepreciation = result.AccAfter
	asset.CarryingAmount = result.CarryAfter
	if result.CarryAfter <= 0 {
		asset.Status = domain.FAFullyDepr
	} else {
		asset.Status = domain.FADepreciating
	}
	if err := e.faRepo.UpdateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return result, nil
}

func (e *DepreciationEngine) RunForCompany(ctx context.Context, companyID string, input domain.FARunDepreciationInput) ([]DepreciationResult, error) {
	assets, _, err := e.faRepo.ListAssets(ctx, domain.FAListFilter{
		CompanyID: companyID,
	})
	if err != nil {
		return nil, err
	}

	var results []DepreciationResult
	for i := range assets {
		a := assets[i]
		if a.Status != domain.FAActive && a.Status != domain.FADepreciating {
			continue
		}
		if a.DepreciationStartDate != nil {
			if !isMonthWithinRange(a.DepreciationStartDate, a.DepreciationEndDate, input.Year, input.Month) {
				continue
			}
		}
		r, err := e.RunForAsset(ctx, &a, input)
		if err != nil {
			if err == domain.ErrFADepreciationExists {
				continue
			}
			return nil, err
		}
		results = append(results, *r)
	}
	return results, nil
}

// ─── Straight-Line ───────────────────────────────────────────────────

func (e *DepreciationEngine) straightLine(asset *domain.FixedAsset, year, month int) (*DepreciationResult, error) {
	depreciableAmount := asset.OriginalCost - asset.ResidualValue
	if depreciableAmount <= 0 {
		return &DepreciationResult{
			AssetID:    asset.ID,
			Amount:     0,
			AccAfter:   asset.AccumulatedDepreciation,
			CarryAfter: asset.CarryingAmount,
		}, nil
	}

	monthlyDepr := depreciableAmount / float64(asset.UsefulLifeMonths)
	monthlyDepr = math.Round(monthlyDepr*100) / 100

	accAfter := asset.AccumulatedDepreciation + monthlyDepr
	carryAfter := asset.OriginalCost - accAfter

	if carryAfter < 0 {
		monthlyDepr = asset.OriginalCost - asset.AccumulatedDepreciation - asset.ResidualValue
		if monthlyDepr < 0 {
			monthlyDepr = 0
		}
		accAfter = asset.OriginalCost - asset.ResidualValue
		carryAfter = asset.ResidualValue
	}

	return &DepreciationResult{
		AssetID:    asset.ID,
		Amount:     monthlyDepr,
		AccAfter:   accAfter,
		CarryAfter: carryAfter,
	}, nil
}

// ─── Declining Balance ───────────────────────────────────────────────

func (e *DepreciationEngine) decliningBalance(asset *domain.FixedAsset, year, month int) (*DepreciationResult, error) {
	depreciableAmount := asset.OriginalCost - asset.ResidualValue
	if depreciableAmount <= 0 {
		return &DepreciationResult{
			AssetID:    asset.ID,
			Amount:     0,
			AccAfter:   asset.AccumulatedDepreciation,
			CarryAfter: asset.CarryingAmount,
		}, nil
	}

	annualRate := 2.0 / float64(asset.UsefulLifeMonths) * 12.0
	if annualRate > 1 {
		annualRate = 1
	}
	monthlyRate := annualRate / 12.0

	monthlyDepr := math.Round(asset.CarryingAmount*monthlyRate*100) / 100

	accAfter := asset.AccumulatedDepreciation + monthlyDepr
	carryAfter := asset.OriginalCost - accAfter

	if carryAfter < asset.ResidualValue {
		monthlyDepr = asset.CarryingAmount - asset.ResidualValue
		if monthlyDepr < 0 {
			monthlyDepr = 0
		}
		accAfter = asset.OriginalCost - asset.ResidualValue
		carryAfter = asset.ResidualValue
	}

	return &DepreciationResult{
		AssetID:    asset.ID,
		Amount:     monthlyDepr,
		AccAfter:   accAfter,
		CarryAfter: carryAfter,
	}, nil
}

func isMonthWithinRange(start, end *string, year, month int) bool {
	if start == nil || *start == "" {
		return true
	}
	target := year*100 + month
	var startYM, endYM, startM, endM int
	fmt.Sscanf(*start, "%4d-%2d", &startYM, &startM)
	startYM = startYM*100 + startM
	if target < startYM {
		return false
	}
	if end != nil && *end != "" {
		fmt.Sscanf(*end, "%4d-%2d", &endYM, &endM)
		endYM = endYM*100 + endM
		if target > endYM {
			return false
		}
	}
	return true
}
