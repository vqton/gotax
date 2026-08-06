package domain

type CostCenter struct {
	ID          string `json:"id"`
	CompanyID   string `json:"company_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	ParentID    string `json:"parent_id,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
