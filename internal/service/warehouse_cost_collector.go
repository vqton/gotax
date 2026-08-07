package service

import (
	"context"
	"time"

	"gotax/internal/domain"
)

// WarehouseCostCollector implements CostDataCollector for warehouse material issuances.
type WarehouseCostCollector struct {
	invTxnRepo domain.InventoryTransactionRepository
}

func NewWarehouseCostCollector(invTxnRepo domain.InventoryTransactionRepository) *WarehouseCostCollector {
	return &WarehouseCostCollector{invTxnRepo: invTxnRepo}
}

func (c *WarehouseCostCollector) CollectMaterialCosts(ctx context.Context, companyID, periodID string) ([]domain.CostPoolLineInput, error) {
	return nil, nil
}

func (c *WarehouseCostCollector) CollectLaborCosts(ctx context.Context, companyID, periodID string) ([]domain.CostPoolLineInput, error) {
	return nil, nil
}

func (c *WarehouseCostCollector) CollectOverheadCosts(ctx context.Context, companyID, periodID string) ([]domain.CostPoolLineInput, error) {
	return nil, nil
}

// CollectMaterialCostsByDateRange collects material issuances between from and to dates.
func (c *WarehouseCostCollector) CollectMaterialCostsByDateRange(ctx context.Context, companyID string, from, to time.Time) ([]domain.CostPoolLineInput, error) {
	// Query all inventory transactions for the company
	txns, _, err := c.invTxnRepo.ListInventoryTransactions(ctx, companyID, "", "", 0, 10000)
	if err != nil {
		return nil, err
	}

	var lines []domain.CostPoolLineInput
	for _, txn := range txns {
		// Filter by date range and transaction type (ISSUE = material issuance)
		if txn.TransType == domain.TransIssue &&
			txn.CreatedAt.After(from) && txn.CreatedAt.Before(to) {
			lines = append(lines, domain.CostPoolLineInput{
				SourceID:    txn.ID,
				Description: txn.Notes,
				Amount:      txn.TotalCost,
			})
		}
	}

	return lines, nil
}
