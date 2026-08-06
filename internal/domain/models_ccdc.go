package domain

type CCDCStatus string

const (
	CCDCActive   CCDCStatus = "ACTIVE"
	CCDCDisposed CCDCStatus = "DISPOSED"
	CCDCLost     CCDCStatus = "LOST"
)

type ToolEquipment struct {
	ID           string     `json:"id"`
	CompanyID    string     `json:"company_id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Category     string     `json:"category,omitempty"`
	PurchaseDate string     `json:"purchase_date"`
	PurchaseCost float64    `json:"purchase_cost"`
	CurrentCost  float64    `json:"current_cost"`
	WarehouseID  string     `json:"warehouse_id,omitempty"`
	Status       CCDCStatus `json:"status"`
	Description  string     `json:"description,omitempty"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
}

type ToolEquipmentCategory struct {
	ID          string `json:"id"`
	CompanyID   string `json:"company_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}
