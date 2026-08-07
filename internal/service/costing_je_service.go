package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type JECreator interface {
	CreateEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error
}

type CostingJEService struct {
	periodRepo       domain.CostingPeriodRepository
	costObjectRepo   domain.CostObjectRepository
	costPoolRepo     domain.CostPoolRepository
	costPoolLineRepo domain.CostPoolLineRepository
	resultRepo       domain.CostingResultRepository
	resultLineRepo   domain.CostingResultLineRepository
	jeCreator        JECreator
}

func NewCostingJEService(
	periodRepo domain.CostingPeriodRepository,
	costObjectRepo domain.CostObjectRepository,
	costPoolRepo domain.CostPoolRepository,
	costPoolLineRepo domain.CostPoolLineRepository,
	resultRepo domain.CostingResultRepository,
	resultLineRepo domain.CostingResultLineRepository,
	jeCreator JECreator,
) *CostingJEService {
	return &CostingJEService{
		periodRepo:       periodRepo,
		costObjectRepo:   costObjectRepo,
		costPoolRepo:     costPoolRepo,
		costPoolLineRepo: costPoolLineRepo,
		resultRepo:       resultRepo,
		resultLineRepo:   resultLineRepo,
		jeCreator:        jeCreator,
	}
}

func (s *CostingJEService) genJEID() string {
	return fmt.Sprintf("JE-%d", time.Now().UnixNano())
}

func nowTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05Z")
}

func buildCostEntry(id, companyID, entryNumber, description, accountDr, accountCr string, amount float64) *domain.JournalEntry {
	return &domain.JournalEntry{
		ID:           id,
		CompanyID:    companyID,
		EntryNumber:  entryNumber,
		EntryDate:    time.Now(),
		Description:  description,
		Status:       domain.JournalEntryDraft,
		CurrencyCode: "VND",
		ExchangeRate: 1,
		Lines: []domain.JournalLine{
			{
				ID:           id + "-D",
				LineNumber:   1,
				AccountCode:  accountDr,
				DebitAmount:  amount,
				CreditAmount: 0,
				Description:  fmt.Sprintf("Dr %s", accountDr),
			},
			{
				ID:           id + "-C",
				LineNumber:   2,
				AccountCode:  accountCr,
				DebitAmount:  0,
				CreditAmount: amount,
				Description:  fmt.Sprintf("Cr %s", accountCr),
			},
		},
	}
}

// GenerateCostPoolEntries: Dr 154 (WIP), Cr 62x (cost pool account)
func (s *CostingJEService) GenerateCostPoolEntries(ctx context.Context, companyID, periodID string) error {
	pools, err := s.costPoolRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return err
	}

	for _, pool := range pools {
		if pool.TotalAmount <= 0 {
			continue
		}

		jeID := s.genJEID()
		entry := buildCostEntry(
			jeID, companyID,
			fmt.Sprintf("COST-%s-%s", periodID, pool.GLAccountCode),
			fmt.Sprintf("Cost allocation: %s to WIP", pool.Name),
			"154", pool.GLAccountCode,
			pool.TotalAmount,
		)

		if err := s.jeCreator.CreateEntry(ctx, entry, "system"); err != nil {
			return err
		}
	}

	return nil
}

// GenerateWIPTransferEntry: Dr 155 (Finished goods), Cr 154 (WIP)
func (s *CostingJEService) GenerateWIPTransferEntry(ctx context.Context, companyID, periodID, objectID string) error {
	results, err := s.resultRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return err
	}

	var result domain.CostingResult
	found := false
	for _, r := range results {
		if r.CostObjectID == objectID {
			result = r
			found = true
			break
		}
	}
	if !found {
		return domain.ErrCostingResultNotFound
	}

	if result.Status != "FINAL" {
		return fmt.Errorf("costing result not finalized")
	}

	if result.TotalCost <= 0 {
		return nil
	}

	entry := buildCostEntry(
		s.genJEID(), companyID,
		fmt.Sprintf("WIP-%s-%s", periodID, objectID),
		fmt.Sprintf("Transfer WIP to finished goods - %s", objectID),
		"155", "154",
		result.TotalCost,
	)

	return s.jeCreator.CreateEntry(ctx, entry, "system")
}

// ClosePeriod: generate cost pool entries → WIP transfers → close period
func (s *CostingJEService) ClosePeriod(ctx context.Context, companyID, periodID, closedBy string) error {
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		return err
	}
	if period.Status == "CLOSED" {
		return domain.ErrCostingPeriodAlreadyClosed
	}

	if err := s.GenerateCostPoolEntries(ctx, companyID, periodID); err != nil {
		return err
	}

	results, err := s.resultRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return err
	}
	for _, r := range results {
		if r.Status == "FINAL" && r.TotalCost > 0 {
			if err := s.GenerateWIPTransferEntry(ctx, companyID, periodID, r.CostObjectID); err != nil {
				return err
			}
		}
	}

	period.Status = "CLOSED"
	period.ClosedBy = closedBy
	period.ClosedAt = nowTimestamp()

	return s.periodRepo.Update(ctx, period)
}

func (s *CostingJEService) CarryForwardOpeningBalance(ctx context.Context, companyID, currentPeriodID string) error {
	currentPeriod, err := s.periodRepo.GetByID(ctx, currentPeriodID)
	if err != nil {
		return err
	}

	periods, err := s.periodRepo.List(ctx, companyID)
	if err != nil {
		return err
	}

	var prevPeriod *domain.CostingPeriod
	for _, p := range periods {
		if p.ID == currentPeriodID {
			continue
		}
		if p.Year < currentPeriod.Year || (p.Year == currentPeriod.Year && p.Month < currentPeriod.Month) {
			if prevPeriod == nil || p.Year > prevPeriod.Year || (p.Year == prevPeriod.Year && p.Month > prevPeriod.Month) {
				prev := p
				prevPeriod = &prev
			}
		}
	}

	if prevPeriod == nil {
		return nil
	}

	prevResults, err := s.resultRepo.ListByPeriod(ctx, companyID, prevPeriod.ID)
	if err != nil {
		return err
	}

	for _, r := range prevResults {
		if r.WIPEnd > 0 {
			ts := nowTimestamp()
			openingResult := &domain.CostingResult{
				ID:            costingGenID("CR"),
				CompanyID:     companyID,
				PeriodID:      currentPeriodID,
				CostObjectID:  r.CostObjectID,
				CostingMethod: r.CostingMethod,
				WIPBegin:      r.WIPEnd,
				OutputQuantity: 0,
				Status:        "DRAFT",
				CreatedAt:     ts,
				UpdatedAt:     ts,
			}
			if err := s.resultRepo.Create(ctx, openingResult); err != nil {
				return err
			}
		}
	}

	return nil
}
