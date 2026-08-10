package domain

// Contract represents a business contract (hợp đồng).
type Contract struct {
	ID                string  `json:"id"`
	CompanyID         string  `json:"company_id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	ContractType      string  `json:"contract_type"` // SALES, PURCHASE, SERVICE, CONSTRUCTION
	Status            string  `json:"status"`        // draft, active, completed, cancelled
	Value             float64 `json:"value"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
	CounterpartyName  string  `json:"counterparty_name"`
	CounterpartyTaxCode string `json:"counterparty_tax_code"`
	Notes             string  `json:"notes,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// ContractPayment tracks payments against a contract.
type ContractPayment struct {
	ID          string  `json:"id"`
	ContractID  string  `json:"contract_id"`
	PaymentDate string  `json:"payment_date"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
}
