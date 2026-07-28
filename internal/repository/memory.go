package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"gotax/internal/domain"
)

// ─── Account ────────────────────────────────────────────────────────

type MemoryAccountRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.Account
}

func NewMemoryAccountRepo() *MemoryAccountRepo {
	return &MemoryAccountRepo{data: make(map[string]*domain.Account)}
}

func (r *MemoryAccountRepo) Accounts() map[string]*domain.Account { return r.data }

func (r *MemoryAccountRepo) Create(_ context.Context, a *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[a.Code]; ok {
		return domain.ErrAccountCodeExists
	}
	cp := *a
	r.data[a.Code] = &cp
	return nil
}

func (r *MemoryAccountRepo) GetByCode(_ context.Context, code string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.data[code]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}
	cp := *a
	return &cp, nil
}

func (r *MemoryAccountRepo) GetAll(_ context.Context, activeOnly bool) ([]domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Account
	for _, a := range r.data {
		if activeOnly && !a.IsActive {
			continue
		}
		out = append(out, *a)
	}
	return out, nil
}

func (r *MemoryAccountRepo) Update(_ context.Context, a *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[a.Code]; !ok {
		return domain.ErrAccountNotFound
	}
	cp := *a
	r.data[a.Code] = &cp
	return nil
}

func (r *MemoryAccountRepo) Delete(_ context.Context, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[code]; !ok {
		return domain.ErrAccountNotFound
	}
	delete(r.data, code)
	return nil
}

func (r *MemoryAccountRepo) GetChildren(_ context.Context, parentCode string) ([]domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Account
	for _, a := range r.data {
		if a.ParentCode == parentCode {
			out = append(out, *a)
		}
	}
	return out, nil
}

// ─── Journal ────────────────────────────────────────────────────────

type MemoryJournalRepo struct {
	mu       sync.RWMutex
	data     map[string]*domain.JournalEntry
	accounts map[string]*domain.Account
	counter  int
}

func NewMemoryJournalRepo() *MemoryJournalRepo {
	return &MemoryJournalRepo{data: make(map[string]*domain.JournalEntry)}
}

func (r *MemoryJournalRepo) SetAccounts(accs map[string]*domain.Account) {
	r.accounts = accs
}

func (r *MemoryJournalRepo) nextID() string {
	r.counter++
	return time.Now().Format("20060102150405") + "-" + formatInt(r.counter)
}

func formatInt(n int) string {
	if n < 10 {
		return "00" + string(rune('0'+n))
	}
	if n < 100 {
		return "0" + string(rune('0'+n/10)) + string(rune('0'+n%10))
	}
	return string(rune('0' + n/100)) + string(rune('0'+(n/10)%10)) + string(rune('0'+n%10))
}

func (r *MemoryJournalRepo) Create(_ context.Context, e *domain.JournalEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e.ID == "" {
		e.ID = r.nextID()
	}
	cp := *e
	cp.Lines = make([]domain.JournalLine, len(e.Lines))
	copy(cp.Lines, e.Lines)
	r.data[e.ID] = &cp
	return nil
}

func (r *MemoryJournalRepo) GetByID(_ context.Context, id string) (*domain.JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.data[id]
	if !ok {
		return nil, domain.ErrJournalNotFound
	}
	return e, nil
}

func (r *MemoryJournalRepo) GetByPeriod(_ context.Context, periodID string) ([]domain.JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.JournalEntry
	for _, e := range r.data {
		if e.PeriodID == periodID {
			out = append(out, *e)
		}
	}
	return out, nil
}

func (r *MemoryJournalRepo) GetByDateRange(_ context.Context, from, to time.Time) ([]domain.JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.JournalEntry
	for _, e := range r.data {
		if !e.EntryDate.Before(from) && !e.EntryDate.After(to) {
			out = append(out, *e)
		}
	}
	return out, nil
}

func (r *MemoryJournalRepo) GetByStatus(_ context.Context, status domain.JournalEntryStatus) ([]domain.JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.JournalEntry
	for _, e := range r.data {
		if e.Status == status {
			out = append(out, *e)
		}
	}
	return out, nil
}

func (r *MemoryJournalRepo) GetByVoucherType(_ context.Context, vt domain.VoucherType) ([]domain.JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.JournalEntry
	for _, e := range r.data {
		if e.VoucherType == vt {
			out = append(out, *e)
		}
	}
	return out, nil
}

func (r *MemoryJournalRepo) UpdateStatus(_ context.Context, id string, status domain.JournalEntryStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.data[id]
	if !ok {
		return domain.ErrJournalNotFound
	}
	e.Status = status
	return nil
}

func (r *MemoryJournalRepo) Update(_ context.Context, e *domain.JournalEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[e.ID]; !ok {
		return domain.ErrJournalNotFound
	}
	cp := *e
	cp.Lines = make([]domain.JournalLine, len(e.Lines))
	copy(cp.Lines, e.Lines)
	r.data[e.ID] = &cp
	return nil
}

func (r *MemoryJournalRepo) Approve(_ context.Context, id, approvedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.data[id]
	if !ok {
		return domain.ErrJournalNotFound
	}
	e.Status = domain.JournalEntryApproved
	e.ApprovedBy = approvedBy
	now := time.Now()
	e.ApprovedAt = &now
	return nil
}

func (r *MemoryJournalRepo) Review(_ context.Context, id, reviewedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.data[id]
	if !ok {
		return domain.ErrJournalNotFound
	}
	e.Status = domain.JournalEntryReviewing
	e.ReviewedBy = reviewedBy
	return nil
}

func (r *MemoryJournalRepo) GetLinesByEntryID(_ context.Context, entryID string) ([]domain.JournalLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.data[entryID]
	if !ok {
		return nil, domain.ErrJournalNotFound
	}
	out := make([]domain.JournalLine, len(e.Lines))
	copy(out, e.Lines)
	return out, nil
}

func (r *MemoryJournalRepo) GetBalance(_ context.Context, accountCode, periodID string) (*domain.AccountBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b := &domain.AccountBalance{AccountCode: accountCode, PeriodID: periodID}
	for _, e := range r.data {
		if e.PeriodID != periodID || e.Status != domain.JournalEntryPosted {
			continue
		}
		for _, l := range e.Lines {
			if l.AccountCode == accountCode {
				b.PeriodDebit += l.DebitAmount
				b.PeriodCredit += l.CreditAmount
			}
		}
	}
	b.Calculate()
	return b, nil
}

func (r *MemoryJournalRepo) GetTrialBalance(_ context.Context, periodID string) ([]domain.AccountBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	balances := make(map[string]*domain.AccountBalance)
	for _, e := range r.data {
		if e.PeriodID != periodID || e.Status != domain.JournalEntryPosted {
			continue
		}
		for _, l := range e.Lines {
			b, ok := balances[l.AccountCode]
			if !ok {
				b = &domain.AccountBalance{AccountCode: l.AccountCode, PeriodID: periodID}
				if r.accounts != nil {
					if a, exists := r.accounts[l.AccountCode]; exists {
						b.AccountType = a.Type
					}
				}
				balances[l.AccountCode] = b
			}
			b.PeriodDebit += l.DebitAmount
			b.PeriodCredit += l.CreditAmount
		}
	}
	var out []domain.AccountBalance
	for _, b := range balances {
		b.Calculate()
		out = append(out, *b)
	}
	return out, nil
}

func (r *MemoryJournalRepo) GetFinancialStatement(_ context.Context, periodID string, types []domain.AccountType) ([]domain.AccountBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make(map[string]*domain.AccountBalance)
	for _, e := range r.data {
		if e.PeriodID != periodID || e.Status != domain.JournalEntryPosted {
			continue
		}
		for _, l := range e.Lines {
			b, ok := all[l.AccountCode]
			if !ok {
				b = &domain.AccountBalance{AccountCode: l.AccountCode, PeriodID: periodID}
				if r.accounts != nil {
					if a, exists := r.accounts[l.AccountCode]; exists {
						b.AccountType = a.Type
					}
				}
				all[l.AccountCode] = b
			}
			b.PeriodDebit += l.DebitAmount
			b.PeriodCredit += l.CreditAmount
		}
	}
	typeSet := make(map[domain.AccountType]bool)
	for _, t := range types {
		typeSet[t] = true
	}
	var out []domain.AccountBalance
	for _, b := range all {
		if typeSet[b.AccountType] {
			b.Calculate()
			out = append(out, *b)
		}
	}
	return out, nil
}

func (r *MemoryJournalRepo) GetAccountUsage(_ context.Context, accountCode string) (*domain.AccountUsage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u := &domain.AccountUsage{AccountCode: accountCode}
	for _, e := range r.data {
		if e.Status != domain.JournalEntryPosted {
			continue
		}
		for _, l := range e.Lines {
			if l.AccountCode != accountCode {
				continue
			}
			u.EntryCount++
			u.TotalDebit += l.DebitAmount
			u.TotalCredit += l.CreditAmount
			ds := e.EntryDate.Format("2006-01-02")
			if ds > u.LastUsedDate {
				u.LastUsedDate = ds
			}
		}
	}
	return u, nil
}

func (r *MemoryJournalRepo) GetPostedEntriesByAccount(_ context.Context, periodID, accountCode string) ([]domain.JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.JournalEntry
	for _, e := range r.data {
		if e.PeriodID != periodID || e.Status != domain.JournalEntryPosted {
			continue
		}
		found := false
		for _, l := range e.Lines {
			if l.AccountCode == accountCode {
				found = true
				break
			}
		}
		if found {
			out = append(out, *e)
		}
	}
	return out, nil
}

// ─── Period ─────────────────────────────────────────────────────────

type MemoryPeriodRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.Period
}

func NewMemoryPeriodRepo() *MemoryPeriodRepo {
	return &MemoryPeriodRepo{data: make(map[string]*domain.Period)}
}

func (r *MemoryPeriodRepo) Create(_ context.Context, p *domain.Period) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryPeriodRepo) GetByID(_ context.Context, id string) (*domain.Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.data[id]
	if !ok {
		return nil, domain.ErrPeriodNotFound
	}
	return p, nil
}

func (r *MemoryPeriodRepo) GetByYearMonth(_ context.Context, year, month int) (*domain.Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.data {
		if p.Year == year && p.Month == month {
			return p, nil
		}
	}
	return nil, domain.ErrPeriodNotFound
}

func (r *MemoryPeriodRepo) GetAll(_ context.Context) ([]domain.Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Period
	for _, p := range r.data {
		out = append(out, *p)
	}
	return out, nil
}

func (r *MemoryPeriodRepo) UpdateStatus(_ context.Context, id string, status domain.PeriodStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.data[id]
	if !ok {
		return domain.ErrPeriodNotFound
	}
	p.Status = status
	return nil
}

func (r *MemoryPeriodRepo) GetOpenPeriod(_ context.Context) (*domain.Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.data {
		if p.Status == domain.PeriodOpen {
			return p, nil
		}
	}
	return nil, domain.ErrPeriodNotFound
}

// ─── User ───────────────────────────────────────────────────────────

type MemoryUserRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.User
}

func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{data: make(map[string]*domain.User)}
}

func (r *MemoryUserRepo) Create(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *u
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	u.ID = cp.ID
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryUserRepo) GetByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.data[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (r *MemoryUserRepo) GetByUsername(_ context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.data {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (r *MemoryUserRepo) GetAll(_ context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.User
	for _, u := range r.data {
		out = append(out, *u)
	}
	return out, nil
}

func (r *MemoryUserRepo) Update(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[u.ID]; !ok {
		return domain.ErrUserNotFound
	}
	cp := *u
	r.data[u.ID] = &cp
	return nil
}

func (r *MemoryUserRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrUserNotFound
	}
	delete(r.data, id)
	return nil
}

// ─── RefreshToken ───────────────────────────────────────────────────

type MemoryRefreshTokenRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.RefreshToken
}

func NewMemoryRefreshTokenRepo() *MemoryRefreshTokenRepo {
	return &MemoryRefreshTokenRepo{data: make(map[string]*domain.RefreshToken)}
}

func (r *MemoryRefreshTokenRepo) Create(_ context.Context, t *domain.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	t.ID = cp.ID
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryRefreshTokenRepo) GetByID(_ context.Context, id string) (*domain.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return t, nil
}

func (r *MemoryRefreshTokenRepo) GetByHash(_ context.Context, hash string) (*domain.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.data {
		if t.TokenHash == hash {
			return t, nil
		}
	}
	return nil, domain.ErrInvalidRefreshToken
}

func (r *MemoryRefreshTokenRepo) GetByUserID(_ context.Context, userID string) ([]domain.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.RefreshToken
	for _, t := range r.data {
		if t.UserID == userID {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (r *MemoryRefreshTokenRepo) Revoke(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.data[id]
	if !ok {
		return domain.ErrInvalidRefreshToken
	}
	now := time.Now()
	t.RevokedAt = &now
	return nil
}

func (r *MemoryRefreshTokenRepo) RevokeAllByUserID(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for _, t := range r.data {
		if t.UserID == userID {
			t.RevokedAt = &now
		}
	}
	return nil
}

// ─── PasswordResetToken ─────────────────────────────────────────────

type MemoryPasswordResetTokenRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.PasswordResetToken
}

func NewMemoryPasswordResetTokenRepo() *MemoryPasswordResetTokenRepo {
	return &MemoryPasswordResetTokenRepo{data: make(map[string]*domain.PasswordResetToken)}
}

func (r *MemoryPasswordResetTokenRepo) Create(_ context.Context, t *domain.PasswordResetToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryPasswordResetTokenRepo) GetByID(_ context.Context, id string) (*domain.PasswordResetToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return t, nil
}

func (r *MemoryPasswordResetTokenRepo) MarkUsed(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.data[id]
	if !ok {
		return domain.ErrUserNotFound
	}
	now := time.Now()
	t.UsedAt = &now
	return nil
}

// ─── AuditLog ───────────────────────────────────────────────────────

type MemoryAuditLogRepo struct {
	mu   sync.RWMutex
	data []domain.AuditEntry
}

func NewMemoryAuditLogRepo() *MemoryAuditLogRepo {
	return &MemoryAuditLogRepo{}
}

func (r *MemoryAuditLogRepo) Create(_ context.Context, e *domain.AuditEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *e
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	r.data = append(r.data, cp)
	return nil
}

func (r *MemoryAuditLogRepo) GetByEntity(_ context.Context, entityType, entityID string) ([]domain.AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.AuditEntry
	for _, e := range r.data {
		if e.EntityType == entityType && e.EntityID == entityID {
			out = append(out, e)
		}
	}
	return out, nil
}

func (r *MemoryAuditLogRepo) GetByUser(_ context.Context, userID string) ([]domain.AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.AuditEntry
	for _, e := range r.data {
		if e.UserID == userID {
			out = append(out, e)
		}
	}
	return out, nil
}

func (r *MemoryAuditLogRepo) GetByDateRange(_ context.Context, from, to time.Time) ([]domain.AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.AuditEntry
	for _, e := range r.data {
		if !e.CreatedAt.Before(from) && !e.CreatedAt.After(to) {
			out = append(out, e)
		}
	}
	return out, nil
}

func (r *MemoryAuditLogRepo) GetAll(_ context.Context, limit int) ([]domain.AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n := limit
	if n > len(r.data) {
		n = len(r.data)
	}
	out := make([]domain.AuditEntry, n)
	copy(out, r.data[len(r.data)-n:])
	return out, nil
}

// ─── ExchangeRate ───────────────────────────────────────────────────

type MemoryExchangeRateRepo struct {
	mu   sync.RWMutex
	data []domain.ExchangeRate
}

func NewMemoryExchangeRateRepo() *MemoryExchangeRateRepo {
	return &MemoryExchangeRateRepo{}
}

func (r *MemoryExchangeRateRepo) Create(_ context.Context, rate *domain.ExchangeRate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *rate
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	r.data = append(r.data, cp)
	return nil
}

func (r *MemoryExchangeRateRepo) GetByCurrencyAndDate(_ context.Context, currencyCode string, rateDate time.Time) (*domain.ExchangeRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rate := range r.data {
		if rate.CurrencyCode == currencyCode && rate.RateDate.Equal(rateDate) {
			return &rate, nil
		}
	}
	return nil, domain.ErrRateNotFound
}

func (r *MemoryExchangeRateRepo) GetByDateRange(_ context.Context, from, to time.Time) ([]domain.ExchangeRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ExchangeRate
	for _, rate := range r.data {
		if !rate.RateDate.Before(from) && !rate.RateDate.After(to) {
			out = append(out, rate)
		}
	}
	return out, nil
}

func (r *MemoryExchangeRateRepo) GetAll(_ context.Context) ([]domain.ExchangeRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.ExchangeRate, len(r.data))
	copy(out, r.data)
	return out, nil
}

func (r *MemoryExchangeRateRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, rate := range r.data {
		if rate.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return domain.ErrRateNotFound
}

// ─── ClosingTemplate ────────────────────────────────────────────────

type MemoryClosingTemplateRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.ClosingTemplate
}

func NewMemoryClosingTemplateRepo() *MemoryClosingTemplateRepo {
	return &MemoryClosingTemplateRepo{data: make(map[string]*domain.ClosingTemplate)}
}

func (r *MemoryClosingTemplateRepo) Create(_ context.Context, t *domain.ClosingTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryClosingTemplateRepo) GetByID(_ context.Context, id string) (*domain.ClosingTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok {
		return nil, domain.ErrTemplateNotFound
	}
	return t, nil
}

func (r *MemoryClosingTemplateRepo) GetAll(_ context.Context) ([]domain.ClosingTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ClosingTemplate
	for _, t := range r.data {
		out = append(out, *t)
	}
	return out, nil
}

func (r *MemoryClosingTemplateRepo) Update(_ context.Context, t *domain.ClosingTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[t.ID]; !ok {
		return domain.ErrTemplateNotFound
	}
	cp := *t
	r.data[t.ID] = &cp
	return nil
}

func (r *MemoryClosingTemplateRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrTemplateNotFound
	}
	delete(r.data, id)
	return nil
}

// ─── Approval ───────────────────────────────────────────────────────

type MemoryApprovalRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.ApprovalRequest
}

func NewMemoryApprovalRepo() *MemoryApprovalRepo {
	return &MemoryApprovalRepo{data: make(map[string]*domain.ApprovalRequest)}
}

func (r *MemoryApprovalRepo) Create(_ context.Context, req *domain.ApprovalRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *req
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	req.ID = cp.ID
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryApprovalRepo) GetByID(_ context.Context, id string) (*domain.ApprovalRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	req, ok := r.data[id]
	if !ok {
		return nil, domain.ErrApprovalNotFound
	}
	return req, nil
}

func (r *MemoryApprovalRepo) GetByStatus(_ context.Context, status domain.ApprovalStatus) ([]domain.ApprovalRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ApprovalRequest
	for _, req := range r.data {
		if req.Status == status {
			out = append(out, *req)
		}
	}
	return out, nil
}

func (r *MemoryApprovalRepo) GetByEntity(_ context.Context, entityType, entityID string) ([]domain.ApprovalRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ApprovalRequest
	for _, req := range r.data {
		if req.EntityType == entityType && req.EntityID == entityID {
			out = append(out, *req)
		}
	}
	return out, nil
}

func (r *MemoryApprovalRepo) UpdateStatus(_ context.Context, id string, status domain.ApprovalStatus, reviewedBy, reviewNote string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	req, ok := r.data[id]
	if !ok {
		return domain.ErrApprovalNotFound
	}
	req.Status = status
	req.ReviewedBy = reviewedBy
	req.ReviewNote = reviewNote
	now := time.Now()
	req.ReviewedAt = &now
	return nil
}

func (r *MemoryApprovalRepo) GetAll(_ context.Context) ([]domain.ApprovalRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.ApprovalRequest
	for _, req := range r.data {
		out = append(out, *req)
	}
	return out, nil
}

// ─── AccountVersion ─────────────────────────────────────────────────

type MemoryAccountVersionRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.AccountVersion
}

func NewMemoryAccountVersionRepo() *MemoryAccountVersionRepo {
	return &MemoryAccountVersionRepo{data: make(map[string]*domain.AccountVersion)}
}

func (r *MemoryAccountVersionRepo) Create(_ context.Context, v *domain.AccountVersion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *v
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryAccountVersionRepo) GetByVersionNumber(_ context.Context, vn string) (*domain.AccountVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.data {
		if v.VersionNumber == vn {
			return v, nil
		}
	}
	return nil, domain.ErrVersionNotFound
}

func (r *MemoryAccountVersionRepo) GetLatest(_ context.Context) (*domain.AccountVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var latest *domain.AccountVersion
	for _, v := range r.data {
		if latest == nil || v.CreatedAt.After(latest.CreatedAt) {
			latest = v
		}
	}
	if latest == nil {
		return nil, domain.ErrVersionNotFound
	}
	return latest, nil
}

func (r *MemoryAccountVersionRepo) GetAll(_ context.Context) ([]domain.AccountVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.AccountVersion
	for _, v := range r.data {
		out = append(out, *v)
	}
	return out, nil
}

// ─── AccountMapping ─────────────────────────────────────────────────

type MemoryAccountMappingRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.AccountMapping
}

func NewMemoryAccountMappingRepo() *MemoryAccountMappingRepo {
	return &MemoryAccountMappingRepo{data: make(map[string]*domain.AccountMapping)}
}

func (r *MemoryAccountMappingRepo) Create(_ context.Context, m *domain.AccountMapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *m
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryAccountMappingRepo) GetByOldCode(_ context.Context, sourceRegime, oldCode string) (*domain.AccountMapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.data {
		if m.SourceRegime == sourceRegime && m.OldCode == oldCode {
			return m, nil
		}
	}
	return nil, domain.ErrMappingNotFound
}

func (r *MemoryAccountMappingRepo) GetByRegime(_ context.Context, sourceRegime, targetRegime string) ([]domain.AccountMapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.AccountMapping
	for _, m := range r.data {
		if m.SourceRegime == sourceRegime && m.TargetRegime == targetRegime {
			out = append(out, *m)
		}
	}
	return out, nil
}

func (r *MemoryAccountMappingRepo) GetAll(_ context.Context) ([]domain.AccountMapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.AccountMapping
	for _, m := range r.data {
		out = append(out, *m)
	}
	return out, nil
}

// ─── AccountAnalysis ────────────────────────────────────────────────

type MemoryAccountAnalysisRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.AccountAnalysis
}

func NewMemoryAccountAnalysisRepo() *MemoryAccountAnalysisRepo {
	return &MemoryAccountAnalysisRepo{data: make(map[string]*domain.AccountAnalysis)}
}

func (r *MemoryAccountAnalysisRepo) Create(_ context.Context, a *domain.AccountAnalysis) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *a
	r.data[cp.AccountCode] = &cp
	return nil
}

func (r *MemoryAccountAnalysisRepo) GetByAccount(_ context.Context, accountCode string) (*domain.AccountAnalysis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.data[accountCode]
	if !ok {
		return nil, domain.ErrAnalysisNotFound
	}
	return a, nil
}

func (r *MemoryAccountAnalysisRepo) Update(_ context.Context, a *domain.AccountAnalysis) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[a.AccountCode]; !ok {
		return domain.ErrAnalysisNotFound
	}
	cp := *a
	r.data[a.AccountCode] = &cp
	return nil
}

// ─── IFRSMapping ────────────────────────────────────────────────────

type MemoryIFRSMappingRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.IFRSMapping
}

func NewMemoryIFRSMappingRepo() *MemoryIFRSMappingRepo {
	return &MemoryIFRSMappingRepo{data: make(map[string]*domain.IFRSMapping)}
}

func (r *MemoryIFRSMappingRepo) Create(_ context.Context, m *domain.IFRSMapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *m
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	r.data[cp.ID] = &cp
	return nil
}

func (r *MemoryIFRSMappingRepo) GetByVASCode(_ context.Context, vasCode string) (*domain.IFRSMapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.data {
		if m.VASCode == vasCode {
			return m, nil
		}
	}
	return nil, domain.ErrIFRSMappingNotFound
}

func (r *MemoryIFRSMappingRepo) GetAll(_ context.Context) ([]domain.IFRSMapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.IFRSMapping
	for _, m := range r.data {
		out = append(out, *m)
	}
	return out, nil
}

func (r *MemoryIFRSMappingRepo) Update(_ context.Context, m *domain.IFRSMapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[m.ID]; !ok {
		return domain.ErrIFRSMappingNotFound
	}
	cp := *m
	r.data[m.ID] = &cp
	return nil
}

// ─── Tax ──────────────────────────────────────────────────────────────────

type MemoryTaxRepo struct {
	mu          sync.RWMutex
	declarations map[string]*domain.TaxDeclaration
	rates       map[string]*domain.TaxRate
	payments    map[string]*domain.TaxPayment
	invoices    map[string]*domain.EInvoice
	calendars   map[string]*domain.TaxCalendar
	alerts      map[string]*domain.TaxAlert
	auditCases  map[string]*domain.TaxAuditCase
}

func NewMemoryTaxRepo() *MemoryTaxRepo {
	return &MemoryTaxRepo{
		declarations: make(map[string]*domain.TaxDeclaration),
		rates:       make(map[string]*domain.TaxRate),
		payments:    make(map[string]*domain.TaxPayment),
		invoices:    make(map[string]*domain.EInvoice),
		calendars:   make(map[string]*domain.TaxCalendar),
		alerts:      make(map[string]*domain.TaxAlert),
		auditCases:  make(map[string]*domain.TaxAuditCase),
	}
}

// ── Declarations ────────────────────────────────────────────────────────

func (r *MemoryTaxRepo) CreateDeclaration(_ context.Context, d *domain.TaxDeclaration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *d
	if cp.ID == "" {
		cp.ID = fmt.Sprintf("T-%d", time.Now().UnixNano())
	}
	if cp.Lines != nil {
		cp.Lines = make([]domain.TaxDeclarationLine, len(d.Lines))
		copy(cp.Lines, d.Lines)
	}
	r.declarations[cp.ID] = &cp
	d.ID = cp.ID
	return nil
}

func (r *MemoryTaxRepo) GetDeclarationByID(_ context.Context, id string) (*domain.TaxDeclaration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.declarations[id]
	if !ok {
		return nil, domain.ErrDeclarationNotFound
	}
	cp := *d
	if d.Lines != nil {
		cp.Lines = make([]domain.TaxDeclarationLine, len(d.Lines))
		copy(cp.Lines, d.Lines)
	}
	return &cp, nil
}

func (r *MemoryTaxRepo) GetDeclarations(_ context.Context, filter domain.TaxDeclarationFilter) ([]domain.TaxDeclaration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.TaxDeclaration
	for _, d := range r.declarations {
		if filter.CompanyID != "" && d.CompanyID != filter.CompanyID { continue }
		if filter.DeclarationType != "" && d.DeclarationType != filter.DeclarationType { continue }
		if filter.Status != "" && d.Status != filter.Status { continue }
		if filter.PeriodYear != 0 && d.TaxPeriod.PeriodYear != filter.PeriodYear { continue }
		if filter.PeriodNumber != 0 && d.TaxPeriod.PeriodNumber != filter.PeriodNumber { continue }
		out = append(out, *d)
	}
	return out, nil
}

func (r *MemoryTaxRepo) UpdateDeclaration(_ context.Context, d *domain.TaxDeclaration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.declarations[d.ID]; !ok {
		return domain.ErrDeclarationNotFound
	}
	cp := *d
	r.declarations[d.ID] = &cp
	return nil
}

func (r *MemoryTaxRepo) UpdateDeclarationStatus(_ context.Context, id string, status domain.DeclarationStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.declarations[id]
	if !ok {
		return domain.ErrDeclarationNotFound
	}
	d.Status = status
	return nil
}

// ── Tax Rates ───────────────────────────────────────────────────────────

func (r *MemoryTaxRepo) CreateRate(_ context.Context, rate *domain.TaxRate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, x := range r.rates {
		if x.RateCode == rate.RateCode {
			return domain.ErrTaxRateCodeExists
		}
	}
	cp := *rate
	if cp.ID == "" {
		cp.ID = fmt.Sprintf("T-%d", time.Now().UnixNano())
	}
	r.rates[cp.ID] = &cp
	rate.ID = cp.ID
	return nil
}

func (r *MemoryTaxRepo) GetRateByID(_ context.Context, id string) (*domain.TaxRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rate, ok := r.rates[id]
	if !ok {
		return nil, domain.ErrTaxRateNotFound
	}
	cp := *rate
	return &cp, nil
}

func (r *MemoryTaxRepo) GetRates(_ context.Context, filter domain.TaxRateFilter) ([]domain.TaxRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.TaxRate
	for _, rate := range r.rates {
		if filter.TaxType != "" && rate.TaxType != filter.TaxType { continue }
		if filter.IsActive != nil && rate.IsActive != *filter.IsActive { continue }
		out = append(out, *rate)
	}
	return out, nil
}

func (r *MemoryTaxRepo) UpdateRate(_ context.Context, rate *domain.TaxRate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rates[rate.ID]; !ok {
		return domain.ErrTaxRateNotFound
	}
	cp := *rate
	r.rates[rate.ID] = &cp
	return nil
}

// ── Tax Payments ────────────────────────────────────────────────────────

func (r *MemoryTaxRepo) CreatePayment(_ context.Context, p *domain.TaxPayment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	if cp.ID == "" {
		cp.ID = fmt.Sprintf("T-%d", time.Now().UnixNano())
	}
	r.payments[cp.ID] = &cp
	p.ID = cp.ID
	return nil
}

func (r *MemoryTaxRepo) GetPaymentByID(_ context.Context, id string) (*domain.TaxPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.payments[id]
	if !ok {
		return nil, domain.ErrTaxPaymentNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *MemoryTaxRepo) GetPayments(_ context.Context, filter domain.PaymentFilter) ([]domain.TaxPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.TaxPayment
	for _, p := range r.payments {
		if filter.CompanyID != "" && p.CompanyID != filter.CompanyID { continue }
		if filter.TaxType != "" && p.TaxType != filter.TaxType { continue }
		if filter.Status != "" && p.Status != filter.Status { continue }
		if filter.PeriodYear != 0 && p.PeriodYear != filter.PeriodYear { continue }
		if filter.PeriodNumber != 0 && p.PeriodNumber != filter.PeriodNumber { continue }
		out = append(out, *p)
	}
	return out, nil
}

func (r *MemoryTaxRepo) UpdatePayment(_ context.Context, p *domain.TaxPayment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.payments[p.ID]; !ok {
		return domain.ErrTaxPaymentNotFound
	}
	cp := *p
	r.payments[p.ID] = &cp
	return nil
}

// ── E-Invoices ──────────────────────────────────────────────────────────

func (r *MemoryTaxRepo) CreateEInvoice(_ context.Context, inv *domain.EInvoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *inv
	if cp.ID == "" {
		cp.ID = fmt.Sprintf("T-%d", time.Now().UnixNano())
	}
	if cp.Lines != nil {
		cp.Lines = make([]domain.EInvoiceLine, len(inv.Lines))
		copy(cp.Lines, inv.Lines)
	}
	r.invoices[cp.ID] = &cp
	inv.ID = cp.ID
	return nil
}

func (r *MemoryTaxRepo) GetEInvoiceByID(_ context.Context, id string) (*domain.EInvoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inv, ok := r.invoices[id]
	if !ok {
		return nil, domain.ErrInvoiceNotFound
	}
	cp := *inv
	return &cp, nil
}

func (r *MemoryTaxRepo) GetEInvoices(_ context.Context, filter domain.EInvoiceFilter) ([]domain.EInvoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.EInvoice
	for _, inv := range r.invoices {
		if filter.CompanyID != "" && inv.CompanyID != filter.CompanyID { continue }
		if filter.Status != "" && inv.Status != filter.Status { continue }
		out = append(out, *inv)
	}
	return out, nil
}

func (r *MemoryTaxRepo) UpdateEInvoice(_ context.Context, inv *domain.EInvoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.invoices[inv.ID]; !ok {
		return domain.ErrInvoiceNotFound
	}
	cp := *inv
	r.invoices[inv.ID] = &cp
	return nil
}

func (r *MemoryTaxRepo) UpdateEInvoiceStatus(_ context.Context, id string, status domain.EInvLifecycleStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return domain.ErrInvoiceNotFound
	}
	inv.Status = status
	return nil
}

// ── Tax Calendar ────────────────────────────────────────────────────────

func (r *MemoryTaxRepo) CreateCalendarEntry(_ context.Context, c *domain.TaxCalendar) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *c
	if cp.ID == "" {
		cp.ID = fmt.Sprintf("T-%d", time.Now().UnixNano())
	}
	r.calendars[cp.ID] = &cp
	c.ID = cp.ID
	return nil
}

func (r *MemoryTaxRepo) GetCalendarEntryByID(_ context.Context, id string) (*domain.TaxCalendar, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.calendars[id]
	if !ok {
		return nil, domain.ErrCalendarNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MemoryTaxRepo) GetCalendarByPeriod(_ context.Context, companyID string, periodYear, periodNumber int) ([]domain.TaxCalendar, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.TaxCalendar
	for _, c := range r.calendars {
		if c.CompanyID == companyID && c.PeriodYear == periodYear && c.PeriodNumber == periodNumber {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (r *MemoryTaxRepo) GetCalendarByCompany(_ context.Context, companyID string) ([]domain.TaxCalendar, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.TaxCalendar
	for _, c := range r.calendars {
		if c.CompanyID == companyID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (r *MemoryTaxRepo) UpdateCalendarStatus(_ context.Context, id string, status domain.CalendarStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.calendars[id]
	if !ok {
		return domain.ErrCalendarNotFound
	}
	c.Status = status
	return nil
}

// ── Alerts ──────────────────────────────────────────────────────────────

func (r *MemoryTaxRepo) CreateAlert(_ context.Context, a *domain.TaxAlert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *a
	if cp.ID == "" {
		cp.ID = fmt.Sprintf("T-%d", time.Now().UnixNano())
	}
	r.alerts[cp.ID] = &cp
	a.ID = cp.ID
	return nil
}

func (r *MemoryTaxRepo) GetAlertByID(_ context.Context, id string) (*domain.TaxAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.alerts[id]
	if !ok {
		return nil, domain.ErrTaxAlertNotFound
	}
	cp := *a
	return &cp, nil
}

func (r *MemoryTaxRepo) GetAlerts(_ context.Context, companyID string, limit int) ([]domain.TaxAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.TaxAlert
	for _, a := range r.alerts {
		if companyID != "" && a.CompanyID != companyID { continue }
		out = append(out, *a)
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// ── Audit Cases ─────────────────────────────────────────────────────────

func (r *MemoryTaxRepo) CreateAuditCase(_ context.Context, a *domain.TaxAuditCase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *a
	if cp.ID == "" {
		cp.ID = fmt.Sprintf("T-%d", time.Now().UnixNano())
	}
	r.auditCases[cp.ID] = &cp
	a.ID = cp.ID
	return nil
}

func (r *MemoryTaxRepo) GetAuditCaseByID(_ context.Context, id string) (*domain.TaxAuditCase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.auditCases[id]
	if !ok {
		return nil, domain.ErrAuditCaseNotFound
	}
	cp := *a
	return &cp, nil
}

func (r *MemoryTaxRepo) GetAuditCases(_ context.Context, companyID string) ([]domain.TaxAuditCase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.TaxAuditCase
	for _, a := range r.auditCases {
		if companyID != "" && a.CompanyID != companyID { continue }
		out = append(out, *a)
	}
	return out, nil
}

func (r *MemoryTaxRepo) UpdateAuditCase(_ context.Context, a *domain.TaxAuditCase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.auditCases[a.ID]; !ok {
		return domain.ErrAuditCaseNotFound
	}
	cp := *a
	r.auditCases[a.ID] = &cp
	return nil
}

// ─── Cash ────────────────────────────────────────────────────────

type MemoryCashRepo struct {
	mu          sync.RWMutex
	receipts    map[string]*domain.CashReceipt
	payments    map[string]*domain.CashPayment
	transfers   map[string]*domain.CashTransfer
	funds       map[string]*domain.PettyCashFund
	inventories map[string]*domain.CashInventory
	advances    map[string]*domain.AdvanceRequest
	settlements map[string]*domain.AdvanceSettlement
	receiptCnt  int
	paymentCnt  int
	advCnt      int
}

func NewMemoryCashRepo() *MemoryCashRepo {
	return &MemoryCashRepo{
		receipts:    make(map[string]*domain.CashReceipt),
		payments:    make(map[string]*domain.CashPayment),
		transfers:   make(map[string]*domain.CashTransfer),
		funds:       make(map[string]*domain.PettyCashFund),
		inventories: make(map[string]*domain.CashInventory),
		advances:    make(map[string]*domain.AdvanceRequest),
		settlements: make(map[string]*domain.AdvanceSettlement),
	}
}

func (r *MemoryCashRepo) nextReceiptID() string {
	r.receiptCnt++
	return time.Now().Format("20060102150405") + "-CR" + formatInt(r.receiptCnt)
}

func (r *MemoryCashRepo) nextPaymentID() string {
	r.paymentCnt++
	return time.Now().Format("20060102150405") + "-CP" + formatInt(r.paymentCnt)
}

func (r *MemoryCashRepo) newID() string {
	r.advCnt++
	return fmt.Sprintf("%d", r.advCnt) + uuidSuffix()
}

func uuidSuffix() string {
	return time.Now().Format("20060102150405")
}

// ── Receipts ─────────────────────────────────────────────────────

func (r *MemoryCashRepo) CreateReceipt(_ context.Context, e *domain.CashReceipt) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *e
	if cp.ID == "" {
		cp.ID = r.nextReceiptID()
	}
	r.receipts[cp.ID] = &cp
	e.ID = cp.ID
	return nil
}

func (r *MemoryCashRepo) GetReceipt(_ context.Context, id string) (*domain.CashReceipt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.receipts[id]
	if !ok {
		return nil, domain.ErrCashReceiptNotFound
	}
	cp := *e
	return &cp, nil
}

func (r *MemoryCashRepo) ListReceipts(_ context.Context, filter domain.CashReceiptFilter) ([]domain.CashReceipt, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var matched []domain.CashReceipt
	for _, e := range r.receipts {
		if filter.CompanyID != "" && e.CompanyID != filter.CompanyID {
			continue
		}
		if filter.ReceiptType != "" && e.ReceiptType != filter.ReceiptType {
			continue
		}
		if filter.Currency != "" && e.Currency != filter.Currency {
			continue
		}
		if filter.Status != "" && e.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && e.VoucherDate < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && e.VoucherDate > filter.ToDate {
			continue
		}
		matched = append(matched, *e)
	}
	total := len(matched)
	if filter.Offset > 0 && filter.Offset < len(matched) {
		matched = matched[filter.Offset:]
	}
	if filter.Limit > 0 && filter.Limit < len(matched) {
		matched = matched[:filter.Limit]
	}
	return matched, total, nil
}

func (r *MemoryCashRepo) UpdateReceipt(_ context.Context, e *domain.CashReceipt) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.receipts[e.ID]; !ok {
		return domain.ErrCashReceiptNotFound
	}
	cp := *e
	r.receipts[e.ID] = &cp
	return nil
}

func (r *MemoryCashRepo) DeleteReceipt(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.receipts[id]; !ok {
		return domain.ErrCashReceiptNotFound
	}
	delete(r.receipts, id)
	return nil
}

func (r *MemoryCashRepo) LastReceiptNo(_ context.Context, companyID, year string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var last string
	prefix := "R-" + year + "-"
	for _, e := range r.receipts {
		if e.CompanyID == companyID && len(e.VoucherNo) >= len(prefix) && e.VoucherNo[:len(prefix)] == prefix {
			if e.VoucherNo > last {
				last = e.VoucherNo
			}
		}
	}
	return last, nil
}

// ── Payments ─────────────────────────────────────────────────────

func (r *MemoryCashRepo) CreatePayment(_ context.Context, p *domain.CashPayment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	if cp.ID == "" {
		cp.ID = r.nextPaymentID()
	}
	r.payments[cp.ID] = &cp
	p.ID = cp.ID
	return nil
}

func (r *MemoryCashRepo) GetPayment(_ context.Context, id string) (*domain.CashPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.payments[id]
	if !ok {
		return nil, domain.ErrCashPaymentNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *MemoryCashRepo) ListPayments(_ context.Context, filter domain.CashPaymentFilter) ([]domain.CashPayment, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var matched []domain.CashPayment
	for _, p := range r.payments {
		if filter.CompanyID != "" && p.CompanyID != filter.CompanyID {
			continue
		}
		if filter.PaymentType != "" && p.PaymentType != filter.PaymentType {
			continue
		}
		if filter.Currency != "" && p.Currency != filter.Currency {
			continue
		}
		if filter.Status != "" && p.Status != filter.Status {
			continue
		}
		if filter.FromDate != "" && p.VoucherDate < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && p.VoucherDate > filter.ToDate {
			continue
		}
		matched = append(matched, *p)
	}
	total := len(matched)
	if filter.Offset > 0 && filter.Offset < len(matched) {
		matched = matched[filter.Offset:]
	}
	if filter.Limit > 0 && filter.Limit < len(matched) {
		matched = matched[:filter.Limit]
	}
	return matched, total, nil
}

func (r *MemoryCashRepo) UpdatePayment(_ context.Context, p *domain.CashPayment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.payments[p.ID]; !ok {
		return domain.ErrCashPaymentNotFound
	}
	cp := *p
	r.payments[p.ID] = &cp
	return nil
}

func (r *MemoryCashRepo) DeletePayment(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.payments[id]; !ok {
		return domain.ErrCashPaymentNotFound
	}
	delete(r.payments, id)
	return nil
}

func (r *MemoryCashRepo) LastPaymentNo(_ context.Context, companyID, year string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var last string
	prefix := "P-" + year + "-"
	for _, p := range r.payments {
		if p.CompanyID == companyID && len(p.VoucherNo) >= len(prefix) && p.VoucherNo[:len(prefix)] == prefix {
			if p.VoucherNo > last {
				last = p.VoucherNo
			}
		}
	}
	return last, nil
}

// ── Transfers ────────────────────────────────────────────────────

func (r *MemoryCashRepo) CreateTransfer(_ context.Context, t *domain.CashTransfer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405") + "-CT" + formatInt(r.receiptCnt)
		r.receiptCnt++
	}
	r.transfers[cp.ID] = &cp
	t.ID = cp.ID
	return nil
}

func (r *MemoryCashRepo) GetTransfer(_ context.Context, id string) (*domain.CashTransfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.transfers[id]
	if !ok {
		return nil, domain.ErrCashTransferNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *MemoryCashRepo) ListTransfers(_ context.Context, companyID string) ([]domain.CashTransfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CashTransfer
	for _, t := range r.transfers {
		if companyID != "" && t.CompanyID != companyID {
			continue
		}
		out = append(out, *t)
	}
	return out, nil
}

// ── Cash Book ────────────────────────────────────────────────────

func (r *MemoryCashRepo) GetCashBook(_ context.Context, companyID, currency, accountID, fromDate, toDate string) (*domain.CashBook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var openingBalance float64
	for _, e := range r.receipts {
		if e.CompanyID != companyID || e.Status != domain.CashPosted {
			continue
		}
		if currency != "" && e.Currency != currency {
			continue
		}
		if accountID != "" && e.CashAccountID != accountID {
			continue
		}
		if e.VoucherDate < fromDate {
			openingBalance += e.AmountVND
		}
	}
	for _, p := range r.payments {
		if p.CompanyID != companyID || p.Status != domain.CashPosted {
			continue
		}
		if currency != "" && p.Currency != currency {
			continue
		}
		if accountID != "" && p.CashAccountID != accountID {
			continue
		}
		if p.VoucherDate < fromDate {
			openingBalance -= p.AmountVND
		}
	}

	var entries []domain.CashBookEntry
	var totalReceipts, totalPayments float64

	for _, e := range r.receipts {
		if e.CompanyID != companyID || e.Status != domain.CashPosted {
			continue
		}
		if currency != "" && e.Currency != currency {
			continue
		}
		if accountID != "" && e.CashAccountID != accountID {
			continue
		}
		if e.VoucherDate < fromDate || e.VoucherDate > toDate {
			continue
		}
		totalReceipts += e.AmountVND
		entries = append(entries, domain.CashBookEntry{
			VoucherDate:   e.VoucherDate,
			VoucherType:   "RECEIPT",
			VoucherNo:     e.VoucherNo,
			Description:   e.Reason,
			ReceiptAmount: e.AmountVND,
			RefID:         e.ID,
		})
	}
	for _, p := range r.payments {
		if p.CompanyID != companyID || p.Status != domain.CashPosted {
			continue
		}
		if currency != "" && p.Currency != currency {
			continue
		}
		if accountID != "" && p.CashAccountID != accountID {
			continue
		}
		if p.VoucherDate < fromDate || p.VoucherDate > toDate {
			continue
		}
		totalPayments += p.AmountVND
		entries = append(entries, domain.CashBookEntry{
			VoucherDate:   p.VoucherDate,
			VoucherType:   "PAYMENT",
			VoucherNo:     p.VoucherNo,
			Description:   p.Reason,
			PaymentAmount: p.AmountVND,
			RefID:         p.ID,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].VoucherDate < entries[j].VoucherDate
	})

	runningBalance := openingBalance
	for i := range entries {
		entries[i].LineNo = i + 1
		runningBalance += entries[i].ReceiptAmount - entries[i].PaymentAmount
		entries[i].RunningBalance = runningBalance
	}

	closingBalance := openingBalance + totalReceipts - totalPayments

	return &domain.CashBook{
		CompanyID:      companyID,
		Currency:       currency,
		AccountID:      accountID,
		FromDate:       fromDate,
		ToDate:         toDate,
		OpeningBalance: openingBalance,
		TotalReceipts:  totalReceipts,
		TotalPayments:  totalPayments,
		ClosingBalance: closingBalance,
		Entries:        entries,
	}, nil
}

func (r *MemoryCashRepo) GetBalance(_ context.Context, companyID, accountID string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var balance float64
	for _, e := range r.receipts {
		if e.CompanyID == companyID && e.CashAccountID == accountID && e.Status == domain.CashPosted {
			balance += e.AmountVND
		}
	}
	for _, p := range r.payments {
		if p.CompanyID == companyID && p.CashAccountID == accountID && p.Status == domain.CashPosted {
			balance -= p.AmountVND
		}
	}
	return balance, nil
}

// ── Petty Cash Funds ─────────────────────────────────────────────

func (r *MemoryCashRepo) CreatePettyCashFund(_ context.Context, f *domain.PettyCashFund) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *f
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	r.funds[cp.ID] = &cp
	f.ID = cp.ID
	return nil
}

func (r *MemoryCashRepo) GetPettyCashFund(_ context.Context, id string) (*domain.PettyCashFund, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.funds[id]
	if !ok {
		return nil, domain.ErrPettyCashFundNotFound
	}
	cp := *f
	return &cp, nil
}

func (r *MemoryCashRepo) ListPettyCashFunds(_ context.Context, companyID string) ([]domain.PettyCashFund, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.PettyCashFund
	for _, f := range r.funds {
		if companyID != "" && f.CompanyID != companyID {
			continue
		}
		out = append(out, *f)
	}
	return out, nil
}

func (r *MemoryCashRepo) UpdatePettyCashFund(_ context.Context, f *domain.PettyCashFund) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.funds[f.ID]; !ok {
		return domain.ErrPettyCashFundNotFound
	}
	cp := *f
	r.funds[f.ID] = &cp
	return nil
}

// ── Cash Inventory ───────────────────────────────────────────────

func (r *MemoryCashRepo) CreateInventory(_ context.Context, inv *domain.CashInventory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *inv
	if cp.ID == "" {
		cp.ID = time.Now().Format("20060102150405")
	}
	if cp.Denominations != nil {
		cp.Denominations = make([]domain.DenominationDetail, len(inv.Denominations))
		copy(cp.Denominations, inv.Denominations)
	}
	r.inventories[cp.ID] = &cp
	inv.ID = cp.ID
	return nil
}

func (r *MemoryCashRepo) GetInventory(_ context.Context, id string) (*domain.CashInventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inv, ok := r.inventories[id]
	if !ok {
		return nil, domain.ErrCashInventoryNotFound
	}
	cp := *inv
	if inv.Denominations != nil {
		cp.Denominations = make([]domain.DenominationDetail, len(inv.Denominations))
		copy(cp.Denominations, inv.Denominations)
	}
	return &cp, nil
}

func (r *MemoryCashRepo) ListInventories(_ context.Context, companyID string) ([]domain.CashInventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.CashInventory
	for _, inv := range r.inventories {
		if companyID != "" && inv.CompanyID != companyID {
			continue
		}
		out = append(out, *inv)
	}
	return out, nil
}

func (r *MemoryCashRepo) UpdateInventory(_ context.Context, inv *domain.CashInventory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.inventories[inv.ID]; !ok {
		return domain.ErrCashInventoryNotFound
	}
	cp := *inv
	if inv.Denominations != nil {
		cp.Denominations = make([]domain.DenominationDetail, len(inv.Denominations))
		copy(cp.Denominations, inv.Denominations)
	}
	r.inventories[inv.ID] = &cp
	return nil
}

// ─── Advance Request / Settlement (Memory) ─────────────────────────────

func (r *MemoryCashRepo) CreateAdvance(_ context.Context, a *domain.AdvanceRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if a.ID == "" {
		a.ID = fmt.Sprintf("ADV-%s", r.newID())
	}
	cp := *a
	r.advances[a.ID] = &cp
	a.ID = cp.ID
	return nil
}

func (r *MemoryCashRepo) GetAdvance(_ context.Context, id string) (*domain.AdvanceRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.advances[id]
	if !ok {
		return nil, domain.ErrAdvanceNotFound
	}
	cp := *a
	return &cp, nil
}

func (r *MemoryCashRepo) ListAdvances(_ context.Context, companyID string) ([]domain.AdvanceRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.AdvanceRequest
	for _, a := range r.advances {
		if companyID != "" && a.CompanyID != companyID {
			continue
		}
		out = append(out, *a)
	}
	return out, nil
}

func (r *MemoryCashRepo) UpdateAdvance(_ context.Context, a *domain.AdvanceRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.advances[a.ID]; !ok {
		return domain.ErrAdvanceNotFound
	}
	cp := *a
	r.advances[a.ID] = &cp
	return nil
}

func (r *MemoryCashRepo) ListAdvancesByStatus(_ context.Context, companyID string, status domain.AdvanceStatus) ([]domain.AdvanceRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.AdvanceRequest
	for _, a := range r.advances {
		if companyID != "" && a.CompanyID != companyID {
			continue
		}
		if a.Status != status {
			continue
		}
		out = append(out, *a)
	}
	return out, nil
}

func (r *MemoryCashRepo) CreateAdvanceSettlement(_ context.Context, s *domain.AdvanceSettlement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s.ID == "" {
		s.ID = fmt.Sprintf("ADVS-%s", r.newID())
	}
	cp := *s
	r.settlements[s.ID] = &cp
	s.ID = cp.ID
	return nil
}
