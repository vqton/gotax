package domain

import "time"

// SystemOption stores per-company configuration key-value pairs.
type SystemOption struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	Category  string `json:"category"`  // personal, global, inventory, number_format, report
	Key       string `json:"key"`
	Value     string `json:"value"` // JSON string
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NumberingRule defines auto-numbering for a voucher type.
type NumberingRule struct {
	ID           string `json:"id"`
	CompanyID    string `json:"company_id"`
	VoucherType  string `json:"voucher_type"`
	Prefix       string `json:"prefix,omitempty"`
	Suffix       string `json:"suffix,omitempty"`
	NumberLength int    `json:"number_length"` // default 5
	Scope        string `json:"scope"`         // company or branch
	ResetRule    string `json:"reset_rule"`    // never, yearly, monthly
	CurrentNum   int    `json:"current_num"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// BackupRecord tracks data backup metadata.
type BackupRecord struct {
	ID         string `json:"id"`
	CompanyID  string `json:"company_id"`
	Filename   string `json:"filename"`
	FileSize   int64  `json:"file_size"`
	BackupType string `json:"backup_type"` // manual, scheduled
	Status     string `json:"status"`     // completed, failed
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
}

func nowString() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}
