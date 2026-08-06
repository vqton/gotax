package domain

type Budget struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	AccountCode string  `json:"account_code"`
	PeriodYear  int     `json:"period_year"`
	PeriodMonth int     `json:"period_month"`
	Budgeted    float64 `json:"budgeted"`
	Actual      float64 `json:"actual"`
	Variance    float64 `json:"variance"`
	Notes       string  `json:"notes,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
