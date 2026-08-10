package service

import (
	"context"
	"fmt"
	"time"

	"gotax/internal/domain"
)

// ─── SystemOptionService ───────────────────────────────────────────

type SystemOptionService struct {
	repo domain.SystemOptionRepository
}

func NewSystemOptionService(repo domain.SystemOptionRepository) *SystemOptionService {
	return &SystemOptionService{repo: repo}
}

func (s *SystemOptionService) Upsert(ctx context.Context, opt *domain.SystemOption) error {
	if opt.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if opt.Category == "" {
		return fmt.Errorf("category is required")
	}
	if opt.Key == "" {
		return fmt.Errorf("key is required")
	}
	return s.repo.Upsert(ctx, opt)
}

func (s *SystemOptionService) GetByCategory(ctx context.Context, companyID, category string) ([]domain.SystemOption, error) {
	if companyID == "" {
		return nil, fmt.Errorf("company_id is required")
	}
	return s.repo.GetByCategory(ctx, companyID, category)
}

func (s *SystemOptionService) GetAll(ctx context.Context, companyID string) ([]domain.SystemOption, error) {
	if companyID == "" {
		return nil, fmt.Errorf("company_id is required")
	}
	return s.repo.GetAll(ctx, companyID)
}

func (s *SystemOptionService) Get(ctx context.Context, companyID, category, key string) (*domain.SystemOption, error) {
	return s.repo.Get(ctx, companyID, category, key)
}

func (s *SystemOptionService) Delete(ctx context.Context, companyID, category, key string) error {
	return s.repo.Delete(ctx, companyID, category, key)
}

// InitDefaults populates default Vietnamese accounting options if not already set.
func (s *SystemOptionService) InitDefaults(ctx context.Context, companyID string) error {
	defaults := map[string]map[string]string{
		"global": {
			"accounting_standard":    "circular99",
			"fiscal_year_start":      "1",
			"base_currency":          "VND",
			"multi_currency_enabled": "false",
			"accounting_start_date":  time.Now().UTC().Format("2006-01-02"),
		},
		"inventory": {
			"costing_method":      "weighted_avg",
			"branch_costing":      "shared",
			"negative_inventory":  "false",
		},
		"number_format": {
			"thousands_separator": ".",
			"decimal_separator":   ",",
			"decimal_places":      "0",
			"negative_display":    "parentheses",
		},
	}
	for category, keys := range defaults {
		for key, value := range keys {
			existing, err := s.repo.Get(ctx, companyID, category, key)
			if err == nil && existing != nil {
				continue // already set
			}
			opt := &domain.SystemOption{
				CompanyID: companyID,
				Category:  category,
				Key:       key,
				Value:     value,
			}
			if err := s.repo.Upsert(ctx, opt); err != nil {
				return fmt.Errorf("init default %s.%s: %w", category, key, err)
			}
		}
	}
	return nil
}

// ─── NumberingRuleService ──────────────────────────────────────────

type NumberingRuleService struct {
	repo domain.NumberingRuleRepository
}

func NewNumberingRuleService(repo domain.NumberingRuleRepository) *NumberingRuleService {
	return &NumberingRuleService{repo: repo}
}

func (s *NumberingRuleService) Create(ctx context.Context, rule *domain.NumberingRule) error {
	if rule.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	if rule.VoucherType == "" {
		return fmt.Errorf("voucher_type is required")
	}
	return s.repo.Create(ctx, rule)
}

func (s *NumberingRuleService) GetByID(ctx context.Context, id string) (*domain.NumberingRule, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *NumberingRuleService) List(ctx context.Context, companyID string) ([]domain.NumberingRule, error) {
	return s.repo.List(ctx, companyID)
}

func (s *NumberingRuleService) Update(ctx context.Context, rule *domain.NumberingRule) error {
	return s.repo.Update(ctx, rule)
}

func (s *NumberingRuleService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetNextNumber atomically increments and returns the next number for a voucher type.
func (s *NumberingRuleService) GetNextNumber(ctx context.Context, companyID, voucherType string) (int, error) {
	return s.repo.IncrementAndGet(ctx, companyID, voucherType)
}

// FormatNumber returns a zero-padded number string per the rule's number_length.
func (s *NumberingRuleService) FormatNumber(ctx context.Context, companyID, voucherType string) (string, error) {
	rule, err := s.repo.GetByVoucherType(ctx, companyID, voucherType)
	if err != nil {
		return "", err
	}
	num, err := s.repo.IncrementAndGet(ctx, companyID, voucherType)
	if err != nil {
		return "", err
	}
	padded := fmt.Sprintf("%0*d", rule.NumberLength, num)
	return rule.Prefix + padded + rule.Suffix, nil
}
