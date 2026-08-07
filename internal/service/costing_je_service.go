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
	collector        domain.CostDataCollector
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

func (s *CostingJEService) SetCollector(c domain.CostDataCollector) {
	s.collector = c
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

// GenerateCOGSEntry: Dr 632, Cr 155
func (s *CostingJEService) GenerateCOGSEntry(ctx context.Context, companyID, periodID, objectID string, quantity float64) error {
	results, err := s.resultRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return err
	}

	var result *domain.CostingResult
	for i, r := range results {
		if r.CostObjectID == objectID {
			result = &results[i]
			break
		}
	}
	if result == nil {
		return domain.ErrCostingResultNotFound
	}

	if result.Status != "FINAL" {
		return fmt.Errorf("costing result not finalized")
	}

	if result.UnitCost <= 0 || quantity <= 0 {
		return nil
	}

	amount := result.UnitCost * quantity

	// Account 632 — Cost of goods sold, Circular 99/2025, Appendix 01
	// Account 155 — Finished goods inventory, Circular 99/2025, Appendix 01
	entry := buildCostEntry(
		s.genJEID(), companyID,
		fmt.Sprintf("COGS-%s-%s", periodID, objectID),
		fmt.Sprintf("COGS for %s - %.0f units", objectID, quantity),
		"632", "155",
		amount,
	)

	if err := s.jeCreator.CreateEntry(ctx, entry, "system"); err != nil {
		return err
	}

	result.COGSAmt += amount
	result.UpdatedAt = nowTimestamp()
	return s.resultRepo.Update(ctx, result)
}

// CollectMaterialCosts: aggregate warehouse issuances into cost pool
func (s *CostingJEService) CollectMaterialCosts(ctx context.Context, companyID, periodID string, lines []domain.CostPoolLineInput) error {
	return s.collectCosts(ctx, companyID, periodID, "621", "Direct materials", "DIRECT_MATERIAL", lines)
}

// CollectLaborCosts: aggregate payroll direct labor into cost pool
func (s *CostingJEService) CollectLaborCosts(ctx context.Context, companyID, periodID string, lines []domain.CostPoolLineInput) error {
	return s.collectCosts(ctx, companyID, periodID, "622", "Direct labor", "DIRECT_LABOR", lines)
}

// AutoCollectCosts uses the CostDataCollector to automatically gather costs from warehouse/payroll
func (s *CostingJEService) AutoCollectCosts(ctx context.Context, companyID, periodID string) error {
	if s.collector == nil {
		return nil
	}

	matLines, err := s.collector.CollectMaterialCosts(ctx, companyID, periodID)
	if err != nil {
		return err
	}
	if err := s.CollectMaterialCosts(ctx, companyID, periodID, matLines); err != nil {
		return err
	}

	labLines, err := s.collector.CollectLaborCosts(ctx, companyID, periodID)
	if err != nil {
		return err
	}
	return s.CollectLaborCosts(ctx, companyID, periodID, labLines)
}

func (s *CostingJEService) collectCosts(ctx context.Context, companyID, periodID, glAccountCode, poolName, poolType string, lines []domain.CostPoolLineInput) error {
	if len(lines) == 0 {
		return nil
	}

	var totalAmount float64
	for _, l := range lines {
		totalAmount += l.Amount
	}
	if totalAmount <= 0 {
		return nil
	}

	pools, err := s.costPoolRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return err
	}

	var pool *domain.CostPool
	for _, p := range pools {
		if p.GLAccountCode == glAccountCode {
			cp := p
			pool = &cp
			break
		}
	}

	if pool == nil {
		pool = &domain.CostPool{
			ID:            costingGenID("CP"),
			CompanyID:     companyID,
			PeriodID:      periodID,
			GLAccountCode: glAccountCode,
			PoolType:      poolType,
			Name:          poolName,
			Status:        domain.CostPoolStatusOpen,
			CreatedAt:     nowTimestamp(),
			UpdatedAt:     nowTimestamp(),
		}
		if err := s.costPoolRepo.Create(ctx, pool); err != nil {
			return err
		}
	}

	for _, l := range lines {
		line := &domain.CostPoolLine{
			ID:           costingGenID("CPL"),
			PoolID:       pool.ID,
			SourceType:   "MANUAL",
			SourceID:     l.SourceID,
			Description:  l.Description,
			Amount:       l.Amount,
			CostCenterID: l.CostCenterID,
			CreatedAt:    nowTimestamp(),
		}
		if err := s.costPoolLineRepo.Create(ctx, line); err != nil {
			return err
		}
	}

	pool.TotalAmount += totalAmount
	pool.UpdatedAt = nowTimestamp()
	if err := s.costPoolRepo.Update(ctx, pool); err != nil {
		return err
	}

	// Account 621/622 — Cost pools per Circular 99/2025, Appendix 01
	// Account 154 — Work-in-progress, Circular 99/2025, Appendix 01
	entry := buildCostEntry(
		s.genJEID(), companyID,
		fmt.Sprintf("COLLECT-%s-%s", periodID, glAccountCode),
		fmt.Sprintf("Collect costs: %s", poolName),
		"154", glAccountCode,
		totalAmount,
	)

	return s.jeCreator.CreateEntry(ctx, entry, "system")
}

// ReopenPeriod: reverse close, generate reversal entries, set period back to OPEN
func (s *CostingJEService) ReopenPeriod(ctx context.Context, companyID, periodID string) error {
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		return err
	}
	if period.Status != "CLOSED" {
		return fmt.Errorf("period is not closed")
	}

	results, err := s.resultRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return err
	}

	for i := range results {
		r := &results[i]
		if r.Status == "FINAL" {
			// Reverse WIP transfer: Dr 154, Cr 155
			reversal := buildCostEntry(
				s.genJEID(), companyID,
				fmt.Sprintf("REV-WIP-%s-%s", periodID, r.CostObjectID),
				fmt.Sprintf("Reversal WIP transfer - %s", r.CostObjectID),
				"154", "155",
				r.TotalCost,
			)
			if err := s.jeCreator.CreateEntry(ctx, reversal, "system"); err != nil {
				return err
			}

			// Reverse COGS: Dr 155, Cr 632
			if r.COGSAmt > 0 {
				cogsReversal := buildCostEntry(
					s.genJEID(), companyID,
					fmt.Sprintf("REV-COGS-%s-%s", periodID, r.CostObjectID),
					fmt.Sprintf("Reversal COGS - %s", r.CostObjectID),
					"155", "632",
					r.COGSAmt,
				)
				if err := s.jeCreator.CreateEntry(ctx, cogsReversal, "system"); err != nil {
					return err
				}
			}

			r.COGSAmt = 0
			r.UpdatedAt = nowTimestamp()
			if err := s.resultRepo.Update(ctx, r); err != nil {
				return err
			}
		}
	}

	pools, err := s.costPoolRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return err
	}

	for _, pool := range pools {
		if pool.TotalAmount > 0 {
			reversal := buildCostEntry(
				s.genJEID(), companyID,
				fmt.Sprintf("REV-COST-%s-%s", periodID, pool.GLAccountCode),
				fmt.Sprintf("Reversal cost allocation - %s", pool.Name),
				pool.GLAccountCode, "154",
				pool.TotalAmount,
			)
			if err := s.jeCreator.CreateEntry(ctx, reversal, "system"); err != nil {
				return err
			}
		}
	}

	period.Status = "OPEN"
	period.ClosedBy = ""
	period.ClosedAt = ""

	return s.periodRepo.Update(ctx, period)
}
