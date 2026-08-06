package domain

import "time"

type RecurringFrequency string

const (
	RecurringMonthly   RecurringFrequency = "MONTHLY"
	RecurringQuarterly RecurringFrequency = "QUARTERLY"
	RecurringYearly    RecurringFrequency = "YEARLY"
)

type RecurringLine struct {
	AccountCode  string  `json:"account_code"`
	DebitAmount  float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	Description  string  `json:"description,omitempty"`
}

type RecurringEntry struct {
	ID           string              `json:"id"`
	CompanyID    string              `json:"company_id"`
	TemplateName string              `json:"template_name"`
	Description  string              `json:"description,omitempty"`
	Frequency    RecurringFrequency  `json:"frequency"`
	DayOfMonth   int                 `json:"day_of_month"`
	IsActive     bool                `json:"is_active"`
	NextRunDate  string              `json:"next_run_date"`
	Lines        []RecurringLine     `json:"lines"`
	CreatedBy    string              `json:"created_by,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func (r *RecurringEntry) Validate() error {
	if r.CompanyID == "" {
		return ErrRecurringCompanyRequired
	}
	if r.TemplateName == "" {
		return ErrRecurringTemplateNameRequired
	}
	switch r.Frequency {
	case RecurringMonthly, RecurringQuarterly, RecurringYearly:
	default:
		return ErrRecurringFrequencyInvalid
	}
	if r.DayOfMonth < 1 || r.DayOfMonth > 31 {
		return ErrRecurringDayOfMonthInvalid
	}
	if len(r.Lines) == 0 {
		return ErrRecurringLinesRequired
	}
	return nil
}
