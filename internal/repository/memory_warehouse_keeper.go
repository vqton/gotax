package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"gotax/internal/domain"
)

var keeperSeq int64

func keeperID(prefix string) string {
	keeperSeq++
	return prefix + time.Now().Format("20060102150405") + fmt.Sprintf("%03d", keeperSeq%1000)
}

type MemoryWarehouseKeeperRepo struct {
	mu              sync.RWMutex
	assignments     map[string]*domain.WarehouseKeeperAssignment
	ledger          map[string]*domain.StockLedgerEntry
	configs         map[string]*domain.WarehouseKeeperConfig
}

func NewMemoryWarehouseKeeperRepo() *MemoryWarehouseKeeperRepo {
	return &MemoryWarehouseKeeperRepo{
		assignments: make(map[string]*domain.WarehouseKeeperAssignment),
		ledger:      make(map[string]*domain.StockLedgerEntry),
		configs:     make(map[string]*domain.WarehouseKeeperConfig),
	}
}

// ─── Assignment ─────────────────────────────────────────────────────────────

func (r *MemoryWarehouseKeeperRepo) CreateAssignment(_ context.Context, a *domain.WarehouseKeeperAssignment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *a
	if cp.ID == "" {
		cp.ID = keeperID("KEEPER-")
	}
	now := time.Now()
	cp.CreatedAt, cp.UpdatedAt = now, now
	cp.IsActive = true
	r.assignments[cp.ID] = &cp
	a.ID = cp.ID
	a.CreatedAt = cp.CreatedAt
	a.UpdatedAt = cp.UpdatedAt
	return nil
}

func (r *MemoryWarehouseKeeperRepo) GetAssignment(_ context.Context, id string) (*domain.WarehouseKeeperAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.assignments[id]
	if !ok {
		return nil, domain.ErrWarehouseNotFound
	}
	cp := *a
	return &cp, nil
}

func (r *MemoryWarehouseKeeperRepo) ListAssignments(_ context.Context, companyID string) ([]domain.WarehouseKeeperAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.WarehouseKeeperAssignment
	for _, a := range r.assignments {
		if a.CompanyID != companyID {
			continue
		}
		out = append(out, *a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].EffectiveFrom.After(out[j].EffectiveFrom)
	})
	return out, nil
}

func (r *MemoryWarehouseKeeperRepo) GetActiveAssignment(_ context.Context, companyID, warehouseID string, date time.Time) (*domain.WarehouseKeeperAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.assignments {
		if a.CompanyID != companyID || a.WarehouseID != warehouseID || !a.IsActive {
			continue
		}
		if !a.EffectiveFrom.After(date) && (a.EffectiveTo == nil || !a.EffectiveTo.Before(date)) {
			cp := *a
			return &cp, nil
		}
	}
	return nil, domain.ErrWarehouseNotFound
}

func (r *MemoryWarehouseKeeperRepo) UpdateAssignment(_ context.Context, a *domain.WarehouseKeeperAssignment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.assignments[a.ID]
	if !ok {
		return domain.ErrWarehouseNotFound
	}
	cp := *a
	cp.UpdatedAt = time.Now()
	cp.CreatedAt = existing.CreatedAt
	r.assignments[a.ID] = &cp
	return nil
}

func (r *MemoryWarehouseKeeperRepo) DeleteAssignment(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.assignments[id]; !ok {
		return domain.ErrWarehouseNotFound
	}
	delete(r.assignments, id)
	return nil
}

// ─── Stock Ledger ───────────────────────────────────────────────────────────

func (r *MemoryWarehouseKeeperRepo) CreateLedgerEntry(_ context.Context, e *domain.StockLedgerEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *e
	if cp.ID == "" {
		cp.ID = keeperID("LEDGER-")
	}
	now := time.Now()
	cp.RecordedAt, cp.CreatedAt = now, now
	cp.Status = domain.LedgerStatusRecorded
	r.ledger[cp.ID] = &cp
	e.ID = cp.ID
	e.RecordedAt = cp.RecordedAt
	e.CreatedAt = cp.CreatedAt
	return nil
}

func (r *MemoryWarehouseKeeperRepo) GetLedgerEntry(_ context.Context, id string) (*domain.StockLedgerEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.ledger[id]
	if !ok {
		return nil, domain.ErrWarehouseNotFound
	}
	cp := *e
	return &cp, nil
}

func (r *MemoryWarehouseKeeperRepo) ListLedgerEntries(_ context.Context, filter domain.LedgerFilter) ([]domain.StockLedgerEntry, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var matched []domain.StockLedgerEntry
	for _, e := range r.ledger {
		if e.CompanyID != filter.CompanyID || e.WarehouseID != filter.WarehouseID {
			continue
		}
		if !filter.From.IsZero() && e.EntryDate.Before(filter.From) {
			continue
		}
		if !filter.To.IsZero() && e.EntryDate.After(filter.To) {
			continue
		}
		if filter.ItemID != "" && e.ItemID != filter.ItemID {
			continue
		}
		if filter.VoucherType != "" && string(e.VoucherType) != filter.VoucherType {
			continue
		}
		if filter.Status != "" && string(e.Status) != filter.Status {
			continue
		}
		matched = append(matched, *e)
	}
	sort.Slice(matched, func(i, j int) bool {
		if !matched[i].EntryDate.Equal(matched[j].EntryDate) {
			return matched[i].EntryDate.After(matched[j].EntryDate)
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	start := (page - 1) * pageSize
	if start >= total {
		return nil, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (r *MemoryWarehouseKeeperRepo) UnrecordLedgerEntry(_ context.Context, id string, unrecordedBy string, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.ledger[id]
	if !ok {
		return domain.ErrWarehouseNotFound
	}
	now := time.Now()
	e.Status = domain.LedgerStatusUnrecorded
	e.UnrecordedBy = unrecordedBy
	e.UnrecordedAt = &now
	e.UnrecordReason = reason
	return nil
}

func (r *MemoryWarehouseKeeperRepo) GetLedgerBalance(_ context.Context, companyID, warehouseID, itemID string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var maxBalance float64
	for _, e := range r.ledger {
		if e.CompanyID == companyID && e.WarehouseID == warehouseID && e.ItemID == itemID && e.Status == domain.LedgerStatusRecorded {
			if e.BalanceQty > maxBalance {
				maxBalance = e.BalanceQty
			}
		}
	}
	return maxBalance, nil
}

// ─── Reconciliation ─────────────────────────────────────────────────────────

func (r *MemoryWarehouseKeeperRepo) GetReconciliationReport(_ context.Context, companyID, warehouseID string, from, to time.Time) ([]domain.KeeperReconciliationItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type itemAgg struct {
		balance  float64
		maxDate  time.Time
	}
	ledgerAgg := make(map[string]*itemAgg)
	for _, e := range r.ledger {
		if e.CompanyID != companyID || e.WarehouseID != warehouseID || e.Status != domain.LedgerStatusRecorded {
			continue
		}
		if !e.EntryDate.Before(from) && !e.EntryDate.After(to) {
			agg, ok := ledgerAgg[e.ItemID]
			if !ok {
				agg = &itemAgg{}
				ledgerAgg[e.ItemID] = agg
			}
			if e.BalanceQty > agg.balance {
				agg.balance = e.BalanceQty
			}
			if e.EntryDate.After(agg.maxDate) {
				agg.maxDate = e.EntryDate
			}
		}
	}

	var out []domain.KeeperReconciliationItem
	for itemID, agg := range ledgerAgg {
		item := domain.KeeperReconciliationItem{
			ItemID:      itemID,
			WarehouseID: warehouseID,
			KeeperQty:   agg.balance,
			LastUpdated: agg.maxDate,
		}
		item.VarianceQty = item.KeeperQty - item.AccountingQty
		item.VarianceValue = item.VarianceQty * item.UnitCost
		out = append(out, item)
	}
	return out, nil
}

// ─── Stock Card ─────────────────────────────────────────────────────────────

func (r *MemoryWarehouseKeeperRepo) GetStockCard(_ context.Context, companyID, warehouseID, itemID string, period string) (*domain.StockCard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	periodStart := period + "-01"
	var lines []domain.StockCardLine
	for _, e := range r.ledger {
		if e.CompanyID != companyID || e.WarehouseID != warehouseID || e.ItemID != itemID || e.Status != domain.LedgerStatusRecorded {
			continue
		}
		if e.EntryDate.Format("2006-01") != period {
			continue
		}
		lines = append(lines, domain.StockCardLine{
			Date:        e.EntryDate,
			VoucherNo:   e.VoucherNo,
			Description: e.Description,
			ReceiptQty:  e.ReceiptQty,
			IssueQty:    e.IssueQty,
			BalanceQty:  e.BalanceQty,
		})
	}
	sort.Slice(lines, func(i, j int) bool {
		if !lines[i].Date.Equal(lines[j].Date) {
			return lines[i].Date.Before(lines[j].Date)
		}
		return lines[i].VoucherNo < lines[j].VoucherNo
	})

	card := &domain.StockCard{
		WarehouseID: warehouseID,
		ItemID:      itemID,
		Period:      period,
		Lines:       lines,
	}
	_ = periodStart
	if len(lines) > 0 {
		card.OpeningBalance = lines[0].BalanceQty - lines[0].ReceiptQty
		card.ClosingBalance = lines[len(lines)-1].BalanceQty
	}
	return card, nil
}

// ─── Keeper Reports ─────────────────────────────────────────────────────────

func (r *MemoryWarehouseKeeperRepo) GetKeeperInventorySummary(_ context.Context, companyID, warehouseID string) ([]domain.KeeperInventorySummaryItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type itemAgg struct {
		qty      float64
		unitCost float64
		totalVal float64
		maxDate  time.Time
	}
	aggs := make(map[string]*itemAgg)
	for _, e := range r.ledger {
		if e.CompanyID != companyID || e.WarehouseID != warehouseID || e.Status != domain.LedgerStatusRecorded {
			continue
		}
		agg, ok := aggs[e.ItemID]
		if !ok {
			agg = &itemAgg{}
			aggs[e.ItemID] = agg
		}
		if e.BalanceQty > agg.qty {
			agg.qty = e.BalanceQty
			agg.unitCost = e.UnitCost
			agg.totalVal = e.TotalValue
		}
		if e.EntryDate.After(agg.maxDate) {
			agg.maxDate = e.EntryDate
		}
	}

	var out []domain.KeeperInventorySummaryItem
	for itemID, agg := range aggs {
		out = append(out, domain.KeeperInventorySummaryItem{
			ItemID:      itemID,
			Quantity:    agg.qty,
			UnitCost:    agg.unitCost,
			TotalValue:  agg.totalVal,
			LastUpdated: agg.maxDate,
		})
	}
	return out, nil
}
