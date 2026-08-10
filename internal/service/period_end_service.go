package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type PeriodEndService struct {
	faSvc DepreciationRunner
	pSvc  *PurchaseService
	svc   *service
}

type DepreciationRunner interface {
	RunDepreciation(ctx context.Context, input domain.FARunDepreciationInput) ([]DepreciationResult, error)
}

func NewPeriodEndService(faSvc DepreciationRunner, pSvc *PurchaseService, svc *service) *PeriodEndService {
	return &PeriodEndService{faSvc: faSvc, pSvc: pSvc, svc: svc}
}

type PeriodEndResult struct {
	DepreciationCount int    `json:"depreciation_count"`
	FXRevaluationID   string `json:"fx_revaluation_id,omitempty"`
	PeriodClosed      bool   `json:"period_closed"`
}

func (s *PeriodEndService) RunPeriodEnd(ctx context.Context, companyID, periodID, userID string, year, month int) (*PeriodEndResult, error) {
	result := &PeriodEndResult{}

	// 1. Run depreciation for all fixed assets
	input := domain.FARunDepreciationInput{
		CompanyID: companyID,
		PeriodID:  periodID,
		Year:      year,
		Month:     month,
		CreatedBy: userID,
	}
	results, err := s.faSvc.RunDepreciation(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("depreciation failed: %w", err)
	}
	result.DepreciationCount = len(results)

	// 2. FX revaluation for AP
	reval, err := s.pSvc.RevalueAP(ctx, companyID, time.Now())
	if err == nil && reval != nil {
		result.FXRevaluationID = reval.ID
	}

	// 3. Close period
	if err := s.svc.ClosePeriod(ctx, periodID); err != nil {
		return nil, fmt.Errorf("period close failed: %w", err)
	}
	result.PeriodClosed = true

	return result, nil
}
