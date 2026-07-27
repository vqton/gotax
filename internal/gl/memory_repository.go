package gl

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

type MemoryAccountRepo struct {
	mu       sync.RWMutex
	accounts map[string]*Account
}

func NewMemoryAccountRepo() *MemoryAccountRepo {
	return &MemoryAccountRepo{accounts: make(map[string]*Account)}
}

func (r *MemoryAccountRepo) Create(_ context.Context, account *Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[account.Code] = account
	return nil
}

func (r *MemoryAccountRepo) GetByCode(_ context.Context, code string) (*Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	acc, ok := r.accounts[code]
	if !ok {
		return nil, ErrAccountNotFound
	}
	return acc, nil
}

func (r *MemoryAccountRepo) GetAll(_ context.Context, activeOnly bool) ([]Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Account
	for _, acc := range r.accounts {
		if activeOnly && !acc.IsActive {
			continue
		}
		result = append(result, *acc)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})
	return result, nil
}

func (r *MemoryAccountRepo) Update(_ context.Context, account *Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[account.Code] = account
	return nil
}

func (r *MemoryAccountRepo) Delete(_ context.Context, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.accounts, code)
	return nil
}

func (r *MemoryAccountRepo) GetChildren(_ context.Context, parentCode string) ([]Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var children []Account
	for _, acc := range r.accounts {
		if acc.ParentCode == parentCode {
			children = append(children, *acc)
		}
	}
	return children, nil
}

type MemoryJournalRepo struct {
	mu       sync.RWMutex
	journals map[string]*JournalEntry
	seq      int
}

func NewMemoryJournalRepo() *MemoryJournalRepo {
	return &MemoryJournalRepo{journals: make(map[string]*JournalEntry)}
}

func (r *MemoryJournalRepo) nextID() string {
	r.seq++
	return fmt.Sprintf("JE-%04d", r.seq)
}

func (r *MemoryJournalRepo) Create(_ context.Context, entry *JournalEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry.ID = r.nextID()
	entry.CreatedAt = time.Now()
	if entry.Status == "" {
		entry.Status = JournalEntryDraft
	}
	if entry.CompanyID == "" {
		entry.CompanyID = "00000000-0000-0000-0000-000000000000"
	}
	if entry.CurrencyCode == "" {
		entry.CurrencyCode = "VND"
	}
	if entry.ExchangeRate == 0 {
		entry.ExchangeRate = 1
	}
	for i := range entry.Lines {
		entry.Lines[i].EntryID = entry.ID
		entry.Lines[i].ID = fmt.Sprintf("%s-L%02d", entry.ID, i+1)
	}
	r.journals[entry.ID] = entry
	return nil
}

func (r *MemoryJournalRepo) GetByID(_ context.Context, id string) (*JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.journals[id]
	if !ok {
		return nil, ErrJournalNotFound
	}
	return entry, nil
}

func (r *MemoryJournalRepo) GetByPeriod(_ context.Context, periodID string) ([]JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []JournalEntry
	for _, e := range r.journals {
		if e.PeriodID == periodID {
			result = append(result, *e)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].EntryDate.Before(result[j].EntryDate)
	})
	return result, nil
}

func (r *MemoryJournalRepo) GetByDateRange(_ context.Context, from, to time.Time) ([]JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []JournalEntry
	for _, e := range r.journals {
		if (e.EntryDate.Equal(from) || e.EntryDate.After(from)) &&
			(e.EntryDate.Equal(to) || e.EntryDate.Before(to)) {
			result = append(result, *e)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].EntryDate.Before(result[j].EntryDate)
	})
	return result, nil
}

func (r *MemoryJournalRepo) GetByStatus(_ context.Context, status JournalEntryStatus) ([]JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []JournalEntry
	for _, e := range r.journals {
		if e.Status == status {
			result = append(result, *e)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].EntryDate.Before(result[j].EntryDate)
	})
	return result, nil
}

func (r *MemoryJournalRepo) GetByVoucherType(_ context.Context, voucherType VoucherType) ([]JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []JournalEntry
	for _, e := range r.journals {
		if e.VoucherType == voucherType {
			result = append(result, *e)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].EntryDate.Before(result[j].EntryDate)
	})
	return result, nil
}

func (r *MemoryJournalRepo) UpdateStatus(_ context.Context, id string, status JournalEntryStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.journals[id]
	if !ok {
		return ErrJournalNotFound
	}
	entry.Status = status
	if status == JournalEntryPosted {
		now := time.Now()
		entry.PostedAt = &now
	}
	return nil
}

func (r *MemoryJournalRepo) Update(_ context.Context, entry *JournalEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.journals[entry.ID] = entry
	return nil
}

func (r *MemoryJournalRepo) Approve(_ context.Context, id, approvedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.journals[id]
	if !ok {
		return ErrJournalNotFound
	}
	entry.Status = JournalEntryApproved
	entry.ApprovedBy = approvedBy
	now := time.Now()
	entry.ApprovedAt = &now
	return nil
}

func (r *MemoryJournalRepo) Review(_ context.Context, id, reviewedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.journals[id]
	if !ok {
		return ErrJournalNotFound
	}
	entry.Status = JournalEntryReviewing
	entry.ReviewedBy = reviewedBy
	return nil
}

func (r *MemoryJournalRepo) GetLinesByEntryID(_ context.Context, entryID string) ([]JournalLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.journals[entryID]
	if !ok {
		return nil, ErrJournalNotFound
	}
	return entry.Lines, nil
}

func (r *MemoryJournalRepo) GetBalance(_ context.Context, accountCode, periodID string) (*AccountBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	balance := &AccountBalance{
		AccountCode: accountCode,
		PeriodID:    periodID,
	}
	for _, entry := range r.journals {
		if entry.PeriodID != periodID || entry.Status != JournalEntryPosted {
			continue
		}
		for _, line := range entry.Lines {
			if line.AccountCode == accountCode {
				balance.PeriodDebit += line.DebitAmount
				balance.PeriodCredit += line.CreditAmount
			}
		}
	}
	return balance, nil
}

func (r *MemoryJournalRepo) GetTrialBalance(_ context.Context, periodID string) ([]AccountBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	accBalances := make(map[string]*AccountBalance)
	for _, entry := range r.journals {
		if entry.PeriodID != periodID || entry.Status != JournalEntryPosted {
			continue
		}
		for _, line := range entry.Lines {
			b, ok := accBalances[line.AccountCode]
			if !ok {
				b = &AccountBalance{
					AccountCode: line.AccountCode,
					PeriodID:    periodID,
				}
				accBalances[line.AccountCode] = b
			}
			b.PeriodDebit += line.DebitAmount
			b.PeriodCredit += line.CreditAmount
		}
	}
	var result []AccountBalance
	for _, b := range accBalances {
		if b.PeriodDebit != 0 || b.PeriodCredit != 0 {
			result = append(result, *b)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].AccountCode < result[j].AccountCode
	})
	return result, nil
}

func (r *MemoryJournalRepo) GetFinancialStatement(_ context.Context, periodID string, accountTypes []AccountType) ([]AccountBalance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	accBalances := make(map[string]*AccountBalance)
	typeMap := make(map[string]AccountType)
	for _, entry := range r.journals {
		if entry.PeriodID != periodID || entry.Status != JournalEntryPosted {
			continue
		}
		for _, line := range entry.Lines {
			b, ok := accBalances[line.AccountCode]
			if !ok {
				b = &AccountBalance{
					AccountCode: line.AccountCode,
					PeriodID:    periodID,
				}
				accBalances[line.AccountCode] = b
			}
			b.PeriodDebit += line.DebitAmount
			b.PeriodCredit += line.CreditAmount
		}
	}
	var result []AccountBalance
	for code, b := range accBalances {
		actType, found := typeMap[code]
		if !found {
			continue
		}
		for _, t := range accountTypes {
			if actType == t {
				result = append(result, *b)
				break
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].AccountCode < result[j].AccountCode
	})
	return result, nil
}

type MemoryPeriodRepo struct {
	mu      sync.RWMutex
	periods map[string]*Period
}

func NewMemoryPeriodRepo() *MemoryPeriodRepo {
	return &MemoryPeriodRepo{periods: make(map[string]*Period)}
}

func (r *MemoryPeriodRepo) Create(_ context.Context, period *Period) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	period.ID = fmt.Sprintf("P-%d-%02d", period.Year, period.Month)
	r.periods[period.ID] = period
	return nil
}

func (r *MemoryPeriodRepo) GetByID(_ context.Context, id string) (*Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.periods[id]
	if !ok {
		return nil, ErrPeriodNotFound
	}
	return p, nil
}

func (r *MemoryPeriodRepo) GetByYearMonth(_ context.Context, year, month int) (*Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id := fmt.Sprintf("P-%d-%02d", year, month)
	p, ok := r.periods[id]
	if !ok {
		return nil, ErrPeriodNotFound
	}
	return p, nil
}

func (r *MemoryPeriodRepo) GetAll(_ context.Context) ([]Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Period
	for _, p := range r.periods {
		result = append(result, *p)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Year != result[j].Year {
			return result[i].Year < result[j].Year
		}
		return result[i].Month < result[j].Month
	})
	return result, nil
}

func (r *MemoryPeriodRepo) UpdateStatus(_ context.Context, id string, status PeriodStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.periods[id]
	if !ok {
		return ErrPeriodNotFound
	}
	p.Status = status
	return nil
}

func (r *MemoryPeriodRepo) GetOpenPeriod(_ context.Context) (*Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.periods {
		if p.Status == PeriodOpen {
			return p, nil
		}
	}
	return nil, ErrPeriodNotFound
}

type MemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{users: make(map[string]*User)}
}

func (r *MemoryUserRepo) Create(_ context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = fmt.Sprintf("U-%d", len(r.users)+1)
	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepo) GetByID(_ context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *MemoryUserRepo) GetByUsername(_ context.Context, username string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *MemoryUserRepo) GetAll(_ context.Context) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []User
	for _, u := range r.users {
		result = append(result, *u)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Username < result[j].Username
	})
	return result, nil
}

func (r *MemoryUserRepo) Update(_ context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, id)
	return nil
}

type MemoryAuditLogRepo struct {
	mu    sync.RWMutex
	logs  []AuditEntry
	seq   int
}

func NewMemoryAuditLogRepo() *MemoryAuditLogRepo {
	return &MemoryAuditLogRepo{logs: make([]AuditEntry, 0)}
}

func (r *MemoryAuditLogRepo) Create(_ context.Context, entry *AuditEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	entry.ID = fmt.Sprintf("AUD-%04d", r.seq)
	entry.CreatedAt = time.Now()
	r.logs = append(r.logs, *entry)
	return nil
}

func (r *MemoryAuditLogRepo) GetByEntity(_ context.Context, entityType, entityID string) ([]AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []AuditEntry
	for _, e := range r.logs {
		if e.EntityType == entityType && e.EntityID == entityID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *MemoryAuditLogRepo) GetByUser(_ context.Context, userID string) ([]AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []AuditEntry
	for _, e := range r.logs {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *MemoryAuditLogRepo) GetByDateRange(_ context.Context, from, to time.Time) ([]AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []AuditEntry
	for _, e := range r.logs {
		if (e.CreatedAt.Equal(from) || e.CreatedAt.After(from)) &&
			(e.CreatedAt.Equal(to) || e.CreatedAt.Before(to)) {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *MemoryAuditLogRepo) GetAll(_ context.Context, limit int) ([]AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if limit <= 0 || limit > len(r.logs) {
		limit = len(r.logs)
	}
	result := make([]AuditEntry, limit)
	for i := 0; i < limit; i++ {
		result[i] = r.logs[len(r.logs)-1-i]
	}
	return result, nil
}

type MemoryExchangeRateRepo struct {
	mu    sync.RWMutex
	rates map[string]*ExchangeRate
	seq   int
}

func NewMemoryExchangeRateRepo() *MemoryExchangeRateRepo {
	return &MemoryExchangeRateRepo{rates: make(map[string]*ExchangeRate)}
}

func (r *MemoryExchangeRateRepo) Create(_ context.Context, rate *ExchangeRate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	rate.ID = fmt.Sprintf("XR-%04d", r.seq)
	key := fmt.Sprintf("%s|%s", rate.CurrencyCode, rate.RateDate.Format("2006-01-02"))
	r.rates[key] = rate
	return nil
}

func (r *MemoryExchangeRateRepo) GetByCurrencyAndDate(_ context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s|%s", currencyCode, rateDate.Format("2006-01-02"))
	rate, ok := r.rates[key]
	if !ok {
		return nil, ErrRateNotFound
	}
	return rate, nil
}

func (r *MemoryExchangeRateRepo) GetByDateRange(_ context.Context, from, to time.Time) ([]ExchangeRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ExchangeRate
	for _, rate := range r.rates {
		if (rate.RateDate.Equal(from) || rate.RateDate.After(from)) &&
			(rate.RateDate.Equal(to) || rate.RateDate.Before(to)) {
			result = append(result, *rate)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RateDate.Before(result[j].RateDate)
	})
	return result, nil
}

func (r *MemoryExchangeRateRepo) GetAll(_ context.Context) ([]ExchangeRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ExchangeRate
	for _, rate := range r.rates {
		result = append(result, *rate)
	}
	return result, nil
}

func (r *MemoryExchangeRateRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for key, rate := range r.rates {
		if rate.ID == id {
			delete(r.rates, key)
			return nil
		}
	}
	return ErrRateNotFound
}

type MemoryClosingTemplateRepo struct {
	mu        sync.RWMutex
	templates map[string]*ClosingTemplate
	seq       int
}

func NewMemoryClosingTemplateRepo() *MemoryClosingTemplateRepo {
	return &MemoryClosingTemplateRepo{templates: make(map[string]*ClosingTemplate)}
}

func (r *MemoryClosingTemplateRepo) Create(_ context.Context, template *ClosingTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	template.ID = fmt.Sprintf("CT-%04d", r.seq)
	r.templates[template.ID] = template
	return nil
}

func (r *MemoryClosingTemplateRepo) GetByID(_ context.Context, id string) (*ClosingTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[id]
	if !ok {
		return nil, ErrTemplateNotFound
	}
	return t, nil
}

func (r *MemoryClosingTemplateRepo) GetAll(_ context.Context) ([]ClosingTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ClosingTemplate
	for _, t := range r.templates {
		result = append(result, *t)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SequenceOrder < result[j].SequenceOrder
	})
	return result, nil
}

func (r *MemoryClosingTemplateRepo) Update(_ context.Context, template *ClosingTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[template.ID] = template
	return nil
}

func (r *MemoryClosingTemplateRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.templates, id)
	return nil
}

var _ AccountRepository = (*MemoryAccountRepo)(nil)
var _ JournalRepository = (*MemoryJournalRepo)(nil)
var _ PeriodRepository = (*MemoryPeriodRepo)(nil)
var _ UserRepository = (*MemoryUserRepo)(nil)
var _ AuditLogRepository = (*MemoryAuditLogRepo)(nil)
var _ ExchangeRateRepository = (*MemoryExchangeRateRepo)(nil)
var _ ClosingTemplateRepository = (*MemoryClosingTemplateRepo)(nil)
