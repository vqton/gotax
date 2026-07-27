package repository

import (
	"context"
	"math"
	"sync"
	"time"

	"gotax/internal/domain"
)

type MemoryOpeningBalanceRepo struct {
	mu         sync.RWMutex
	balances   map[string]*domain.OpeningBalance
	details    map[string][]domain.OpeningBalanceDetail // balanceID -> details
	cfLogs     map[string]*domain.CarryForwardLog
	mappings   map[string]*domain.Circular99Mapping
	migrations map[string]*domain.BalanceMigration
	counter    int
}

func NewMemoryOpeningBalanceRepo() *MemoryOpeningBalanceRepo {
	return &MemoryOpeningBalanceRepo{
		balances:   make(map[string]*domain.OpeningBalance),
		details:    make(map[string][]domain.OpeningBalanceDetail),
		cfLogs:     make(map[string]*domain.CarryForwardLog),
		mappings:   make(map[string]*domain.Circular99Mapping),
		migrations: make(map[string]*domain.BalanceMigration),
	}
}

func (r *MemoryOpeningBalanceRepo) nextID(prefix string) string {
	r.counter++
	return prefix + time.Now().Format("20060102150405") + formatInt(r.counter)
}

func (r *MemoryOpeningBalanceRepo) Create(_ context.Context, ob *domain.OpeningBalance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.balances {
		if existing.CompanyID == ob.CompanyID && existing.PeriodID == ob.PeriodID &&
			existing.AccountCode == ob.AccountCode && existing.CurrencyCode == ob.CurrencyCode &&
			existing.Status == domain.OBStatusApproved {
			return domain.ErrOpeningBalanceExists
		}
	}
	cp := *ob
	if cp.ID == "" {
		cp.ID = r.nextID("OB")
	}
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	if cp.Status == "" {
		cp.Status = domain.OBStatusDraft
	}
	if cp.SourceType == "" {
		cp.SourceType = "MANUAL"
	}
	if cp.CurrencyCode == "" {
		cp.CurrencyCode = "VND"
	}
	if cp.OriginalAmount == 0 && cp.DebitAmount > 0 {
		cp.OriginalAmount = cp.DebitAmount
	} else if cp.OriginalAmount == 0 && cp.CreditAmount > 0 {
		cp.OriginalAmount = cp.CreditAmount
	}
	r.balances[cp.ID] = &cp
	ob.ID = cp.ID
	return nil
}

func (r *MemoryOpeningBalanceRepo) GetByID(_ context.Context, id string) (*domain.OpeningBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ob, ok := r.balances[id]
	if !ok {
		return nil, domain.ErrOpeningBalanceNotFound
	}
	cp := *ob
	return &cp, nil
}

func (r *MemoryOpeningBalanceRepo) List(_ context.Context, filter domain.OBListFilter) ([]domain.OpeningBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.OpeningBalance
	for _, ob := range r.balances {
		if ob.CompanyID != filter.CompanyID {
			continue
		}
		if filter.PeriodID != "" && ob.PeriodID != filter.PeriodID {
			continue
		}
		if filter.Status != "" && ob.Status != filter.Status {
			continue
		}
		if filter.AccountCode != "" && ob.AccountCode != filter.AccountCode {
			continue
		}
		if filter.Currency != "" && ob.CurrencyCode != filter.Currency {
			continue
		}
		if filter.SourceType != "" && ob.SourceType != filter.SourceType {
			continue
		}
		out = append(out, *ob)
	}
	if out == nil {
		return []domain.OpeningBalance{}, nil
	}
	return out, nil
}

func (r *MemoryOpeningBalanceRepo) GetByAccount(_ context.Context, companyID, periodID, accountCode string) (*domain.OpeningBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ob := range r.balances {
		if ob.CompanyID == companyID && ob.PeriodID == periodID &&
			ob.AccountCode == accountCode && ob.Status == domain.OBStatusApproved {
			cp := *ob
			return &cp, nil
		}
	}
	return nil, domain.ErrOpeningBalanceNotFound
}

func (r *MemoryOpeningBalanceRepo) Update(_ context.Context, ob *domain.OpeningBalance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.balances[ob.ID]
	if !ok {
		return domain.ErrOpeningBalanceNotFound
	}
	if existing.Status == domain.OBStatusApproved || existing.Status == domain.OBStatusCorrected {
		return domain.ErrOpeningBalanceImmutable
	}
	cp := *ob
	cp.UpdatedAt = time.Now()
	r.balances[ob.ID] = &cp
	return nil
}

func (r *MemoryOpeningBalanceRepo) UpdateStatus(_ context.Context, id string, status domain.OpeningBalanceStatus, approvedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ob, ok := r.balances[id]
	if !ok {
		return domain.ErrOpeningBalanceNotFound
	}
	if ob.Status.IsTerminal() && status != domain.OBStatusCorrected {
		return domain.ErrOpeningBalanceImmutable
	}
	ob.Status = status
	ob.UpdatedAt = time.Now()
	if status == domain.OBStatusApproved {
		now := time.Now()
		ob.ApprovedAt = &now
		ob.ApprovedBy = approvedBy
	}
	return nil
}

func (r *MemoryOpeningBalanceRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ob, ok := r.balances[id]
	if !ok {
		return domain.ErrOpeningBalanceNotFound
	}
	if ob.Status == domain.OBStatusApproved || ob.Status == domain.OBStatusCorrected {
		return domain.ErrOpeningBalanceImmutable
	}
	delete(r.balances, id)
	delete(r.details, id)
	return nil
}

func (r *MemoryOpeningBalanceRepo) BulkCreate(_ context.Context, balances []domain.OpeningBalance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range balances {
		cp := balances[i]
		if cp.ID == "" {
			r.counter++
			cp.ID = r.nextID("OB")
		}
		cp.CreatedAt = time.Now()
		cp.UpdatedAt = cp.CreatedAt
		if cp.Status == "" {
			cp.Status = domain.OBStatusDraft
		}
		if cp.SourceType == "" {
			cp.SourceType = "MANUAL"
		}
		if cp.CurrencyCode == "" {
			cp.CurrencyCode = "VND"
		}
		r.balances[cp.ID] = &cp
		balances[i].ID = cp.ID
	}
	return nil
}

func (r *MemoryOpeningBalanceRepo) BulkUpdateStatus(_ context.Context, ids []string, status domain.OpeningBalanceStatus, approvedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for _, id := range ids {
		ob, ok := r.balances[id]
		if !ok {
			return domain.ErrOpeningBalanceNotFound
		}
		ob.Status = status
		ob.UpdatedAt = now
		if status == domain.OBStatusApproved {
			ob.ApprovedAt = &now
			ob.ApprovedBy = approvedBy
		}
	}
	return nil
}

func (r *MemoryOpeningBalanceRepo) CreateDetail(_ context.Context, d *domain.OpeningBalanceDetail) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *d
	if cp.ID == "" {
		r.counter++
		cp.ID = r.nextID("OBD")
	}
	cp.CreatedAt = time.Now()
	r.details[cp.OpeningBalanceID] = append(r.details[cp.OpeningBalanceID], cp)
	d.ID = cp.ID
	return nil
}

func (r *MemoryOpeningBalanceRepo) GetDetails(_ context.Context, balanceID string) ([]domain.OpeningBalanceDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	details, ok := r.details[balanceID]
	if !ok || len(details) == 0 {
		return []domain.OpeningBalanceDetail{}, nil
	}
	out := make([]domain.OpeningBalanceDetail, len(details))
	copy(out, details)
	return out, nil
}

func (r *MemoryOpeningBalanceRepo) DeleteDetail(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for balanceID, dets := range r.details {
		for i, d := range dets {
			if d.ID == id {
				r.details[balanceID] = append(dets[:i], dets[i+1:]...)
				return nil
			}
		}
	}
	return domain.ErrOpeningBalanceNotFound
}

func (r *MemoryOpeningBalanceRepo) DeleteDetails(_ context.Context, balanceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.details, balanceID)
	return nil
}

func (r *MemoryOpeningBalanceRepo) GetTotals(_ context.Context, companyID, periodID string) (totalDebit, totalCredit float64, _ error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ob := range r.balances {
		if ob.CompanyID == companyID && ob.PeriodID == periodID && ob.Status == domain.OBStatusApproved {
			totalDebit += ob.DebitAmount
			totalCredit += ob.CreditAmount
		}
	}
	totalDebit = math.Round(totalDebit*100) / 100
	totalCredit = math.Round(totalCredit*100) / 100
	return
}

func (r *MemoryOpeningBalanceRepo) ValidateBalanced(_ context.Context, companyID, periodID string) (bool, error) {
	d, c, err := r.GetTotals(nil, companyID, periodID)
	if err != nil {
		return false, err
	}
	return math.Abs(d-c) < 0.01, nil
}

func (r *MemoryOpeningBalanceRepo) CreateCarryForwardLog(_ context.Context, log *domain.CarryForwardLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *log
	if cp.ID == "" {
		r.counter++
		cp.ID = r.nextID("CF")
	}
	cp.ExecutedAt = time.Now()
	r.cfLogs[cp.ID] = &cp
	log.ID = cp.ID
	return nil
}

func (r *MemoryOpeningBalanceRepo) GetCarryForwardLogs(_ context.Context, companyID string) ([]domain.CarryForwardLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CarryForwardLog
	for _, l := range r.cfLogs {
		if l.CompanyID == companyID {
			out = append(out, *l)
		}
	}
	if out == nil {
		return []domain.CarryForwardLog{}, nil
	}
	return out, nil
}

func (r *MemoryOpeningBalanceRepo) GetCarryForwardLogByID(_ context.Context, id string) (*domain.CarryForwardLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.cfLogs[id]
	if !ok {
		return nil, domain.ErrOpeningBalanceNotFound
	}
	cp := *l
	return &cp, nil
}

func (r *MemoryOpeningBalanceRepo) CreateCircular99Mapping(_ context.Context, m *domain.Circular99Mapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *m
	if cp.ID == "" {
		r.counter++
		cp.ID = r.nextID("C99")
	}
	r.mappings[cp.ID] = &cp
	m.ID = cp.ID
	return nil
}

func (r *MemoryOpeningBalanceRepo) ListCircular99Mappings(_ context.Context) ([]domain.Circular99Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Circular99Mapping
	for _, m := range r.mappings {
		out = append(out, *m)
	}
	if out == nil {
		return []domain.Circular99Mapping{}, nil
	}
	return out, nil
}

func (r *MemoryOpeningBalanceRepo) GetCircular99MappingByOldCode(_ context.Context, oldCode string) (*domain.Circular99Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.mappings {
		if m.OldAccountCode == oldCode {
			cp := *m
			return &cp, nil
		}
	}
	return nil, domain.ErrCircular99MappingNotFound
}

func (r *MemoryOpeningBalanceRepo) CreateMigration(_ context.Context, m *domain.BalanceMigration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *m
	if cp.ID == "" {
		r.counter++
		cp.ID = r.nextID("MIG")
	}
	cp.CreatedAt = time.Now()
	r.migrations[cp.ID] = &cp
	m.ID = cp.ID
	return nil
}

func (r *MemoryOpeningBalanceRepo) GetMigrationByID(_ context.Context, id string) (*domain.BalanceMigration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.migrations[id]
	if !ok {
		return nil, domain.ErrOpeningBalanceNotFound
	}
	cp := *m
	return &cp, nil
}

func (r *MemoryOpeningBalanceRepo) ListMigrations(_ context.Context, companyID string) ([]domain.BalanceMigration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.BalanceMigration
	for _, m := range r.migrations {
		if m.CompanyID == companyID {
			out = append(out, *m)
		}
	}
	if out == nil {
		return []domain.BalanceMigration{}, nil
	}
	return out, nil
}
