package domain

type NotificationType string

const (
	NotifTypeINFO    NotificationType = "INFO"
	NotifTypeWARNING NotificationType = "WARNING"
	NotifTypeERROR   NotificationType = "ERROR"
	NotifTypeDUE     NotificationType = "DUE"
)

type Notification struct {
	ID         string           `json:"id"`
	CompanyID  string           `json:"company_id"`
	UserID     string           `json:"user_id"`
	EntityType string           `json:"entity_type,omitempty"`
	EntityID   string           `json:"entity_id,omitempty"`
	Type       NotificationType `json:"type"`
	Title      string           `json:"title"`
	Message    string           `json:"message"`
	Link       string           `json:"link,omitempty"`
	ReadAt     string           `json:"read_at,omitempty"`
	CreatedAt  string           `json:"created_at"`
}
