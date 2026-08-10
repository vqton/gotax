package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gotax/internal/domain"
)

// ─── SystemOption ──────────────────────────────────────────────────

type MemorySystemOptionRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.SystemOption // key: "companyID:category:key"
}

func NewMemorySystemOptionRepo() *MemorySystemOptionRepo {
	return &MemorySystemOptionRepo{data: make(map[string]*domain.SystemOption)}
}

func optKey(companyID, category, key string) string {
	return companyID + ":" + category + ":" + key
}

func (r *MemorySystemOptionRepo) Get(_ context.Context, companyID, category, key string) (*domain.SystemOption, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	opt, ok := r.data[optKey(companyID, category, key)]
	if !ok {
		return nil, domain.ErrSystemOptionNotFound
	}
	cp := *opt
	return &cp, nil
}

func (r *MemorySystemOptionRepo) GetByCategory(_ context.Context, companyID, category string) ([]domain.SystemOption, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.SystemOption
	for _, opt := range r.data {
		if opt.CompanyID == companyID && opt.Category == category {
			cp := *opt
			out = append(out, cp)
		}
	}
	return out, nil
}

func (r *MemorySystemOptionRepo) GetAll(_ context.Context, companyID string) ([]domain.SystemOption, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.SystemOption
	for _, opt := range r.data {
		if opt.CompanyID == companyID {
			cp := *opt
			out = append(out, cp)
		}
	}
	return out, nil
}

func (r *MemorySystemOptionRepo) Upsert(_ context.Context, opt *domain.SystemOption) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := optKey(opt.CompanyID, opt.Category, opt.Key)
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	if existing, ok := r.data[k]; ok {
		existing.Value = opt.Value
		existing.UpdatedAt = now
		opt.ID = existing.ID
		opt.CreatedAt = existing.CreatedAt
		opt.UpdatedAt = now
	} else {
		if opt.ID == "" {
			opt.ID = fmt.Sprintf("SO-%d", time.Now().UnixNano())
		}
		opt.CreatedAt = now
		opt.UpdatedAt = now
		cp := *opt
		r.data[k] = &cp
	}
	return nil
}

func (r *MemorySystemOptionRepo) Delete(_ context.Context, companyID, category, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := optKey(companyID, category, key)
	if _, ok := r.data[k]; !ok {
		return domain.ErrSystemOptionNotFound
	}
	delete(r.data, k)
	return nil
}

// ─── NumberingRule ─────────────────────────────────────────────────

type MemoryNumberingRuleRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.NumberingRule
}

func NewMemoryNumberingRuleRepo() *MemoryNumberingRuleRepo {
	return &MemoryNumberingRuleRepo{data: make(map[string]*domain.NumberingRule)}
}

func (r *MemoryNumberingRuleRepo) Create(_ context.Context, rule *domain.NumberingRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.CompanyID == rule.CompanyID && existing.VoucherType == rule.VoucherType {
			return domain.ErrNumberingRuleExists
		}
	}
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	if rule.ID == "" {
		rule.ID = fmt.Sprintf("NR-%d", time.Now().UnixNano())
	}
	if rule.NumberLength == 0 {
		rule.NumberLength = 5
	}
	if rule.Scope == "" {
		rule.Scope = "company"
	}
	if rule.ResetRule == "" {
		rule.ResetRule = "never"
	}
	rule.CreatedAt = now
	rule.UpdatedAt = now
	cp := *rule
	r.data[rule.ID] = &cp
	return nil
}

func (r *MemoryNumberingRuleRepo) GetByID(_ context.Context, id string) (*domain.NumberingRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.data[id]
	if !ok {
		return nil, domain.ErrNumberingRuleNotFound
	}
	cp := *rule
	return &cp, nil
}

func (r *MemoryNumberingRuleRepo) GetByVoucherType(_ context.Context, companyID, voucherType string) (*domain.NumberingRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rule := range r.data {
		if rule.CompanyID == companyID && rule.VoucherType == voucherType {
			cp := *rule
			return &cp, nil
		}
	}
	return nil, domain.ErrNumberingRuleNotFound
}

func (r *MemoryNumberingRuleRepo) List(_ context.Context, companyID string) ([]domain.NumberingRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.NumberingRule
	for _, rule := range r.data {
		if rule.CompanyID == companyID {
			cp := *rule
			out = append(out, cp)
		}
	}
	return out, nil
}

func (r *MemoryNumberingRuleRepo) Update(_ context.Context, rule *domain.NumberingRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[rule.ID]; !ok {
		return domain.ErrNumberingRuleNotFound
	}
	rule.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	cp := *rule
	r.data[rule.ID] = &cp
	return nil
}

func (r *MemoryNumberingRuleRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return domain.ErrNumberingRuleNotFound
	}
	delete(r.data, id)
	return nil
}

func (r *MemoryNumberingRuleRepo) IncrementAndGet(_ context.Context, companyID, voucherType string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rule := range r.data {
		if rule.CompanyID == companyID && rule.VoucherType == voucherType {
			rule.CurrentNum++
			rule.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			return rule.CurrentNum, nil
		}
	}
	return 0, domain.ErrNumberingRuleNotFound
}

// ─── Backup ────────────────────────────────────────────────────────

type MemoryBackupRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.BackupRecord
}

func NewMemoryBackupRepo() *MemoryBackupRepo {
	return &MemoryBackupRepo{data: make(map[string]*domain.BackupRecord)}
}

func (r *MemoryBackupRepo) Create(_ context.Context, b *domain.BackupRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if b.ID == "" {
		b.ID = fmt.Sprintf("BK-%d", time.Now().UnixNano())
	}
	b.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	cp := *b
	r.data[b.ID] = &cp
	return nil
}

func (r *MemoryBackupRepo) GetByID(_ context.Context, id string) (*domain.BackupRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.data[id]
	if !ok {
		return nil, domain.ErrBackupNotFound
	}
	cp := *b
	return &cp, nil
}

func (r *MemoryBackupRepo) List(_ context.Context, companyID string) ([]domain.BackupRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.BackupRecord
	for _, b := range r.data {
		if b.CompanyID == companyID {
			cp := *b
			out = append(out, cp)
		}
	}
	return out, nil
}

func (r *MemoryBackupRepo) UpdateStatus(_ context.Context, id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.data[id]
	if !ok {
		return domain.ErrBackupNotFound
	}
	b.Status = status
	return nil
}
