package service

import (
	"context"
	"fmt"

	"gotax/internal/domain"
)

type ReportOptionService struct {
	optRepo domain.SystemOptionRepository
}

func NewReportOptionService(optRepo domain.SystemOptionRepository) *ReportOptionService {
	return &ReportOptionService{optRepo: optRepo}
}

// ReportOptions holds company-specific report formatting.
type ReportOptions struct {
	CompanyName       string `json:"company_name"`
	CompanyAddress    string `json:"company_address"`
	CompanyTaxCode    string `json:"company_tax_code"`
	CompanyPhone      string `json:"company_phone"`
	CompanyEmail      string `json:"company_email"`
	CompanyWebsite    string `json:"company_website"`
	ReportFontFamily  string `json:"report_font_family"`
	ReportFontSize    string `json:"report_font_size"`
	Alignment         string `json:"alignment"` // left, center, right
	RepeatOnEachPage  bool   `json:"repeat_on_each_page"`
}

// Get reads report options from system options.
func (s *ReportOptionService) Get(ctx context.Context, companyID string) (*ReportOptions, error) {
	opts, err := s.optRepo.GetByCategory(ctx, companyID, "report")
	if err != nil {
		return nil, err
	}
	ro := &ReportOptions{
		Alignment:        "left",
		ReportFontFamily: "Arial",
		ReportFontSize:   "10",
	}
	for _, opt := range opts {
		switch opt.Key {
		case "company_name":
			ro.CompanyName = opt.Value
		case "company_address":
			ro.CompanyAddress = opt.Value
		case "company_tax_code":
			ro.CompanyTaxCode = opt.Value
		case "company_phone":
			ro.CompanyPhone = opt.Value
		case "company_email":
			ro.CompanyEmail = opt.Value
		case "company_website":
			ro.CompanyWebsite = opt.Value
		case "report_font_family":
			ro.ReportFontFamily = opt.Value
		case "report_font_size":
			ro.ReportFontSize = opt.Value
		case "alignment":
			ro.Alignment = opt.Value
		case "repeat_on_each_page":
			ro.RepeatOnEachPage = opt.Value == "true"
		}
	}
	return ro, nil
}

// Update saves report options to system options.
func (s *ReportOptionService) Update(ctx context.Context, companyID string, ro *ReportOptions) error {
	if companyID == "" {
		return fmt.Errorf("company_id is required")
	}
	entries := map[string]string{
		"company_name":       ro.CompanyName,
		"company_address":    ro.CompanyAddress,
		"company_tax_code":   ro.CompanyTaxCode,
		"company_phone":      ro.CompanyPhone,
		"company_email":      ro.CompanyEmail,
		"company_website":    ro.CompanyWebsite,
		"report_font_family": ro.ReportFontFamily,
		"report_font_size":   ro.ReportFontSize,
		"alignment":          ro.Alignment,
		"repeat_on_each_page": fmt.Sprintf("%v", ro.RepeatOnEachPage),
	}
	for key, value := range entries {
		opt := &domain.SystemOption{
			CompanyID: companyID,
			Category:  "report",
			Key:       key,
			Value:     value,
		}
		if err := s.optRepo.Upsert(ctx, opt); err != nil {
			return fmt.Errorf("save %s: %w", key, err)
		}
	}
	return nil
}
