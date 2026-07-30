package domain

import "time"

type CompanyBankAccountGORM struct {
	ID           string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID    string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	BankName     string     `gorm:"column:bank_name;not null;size:255" json:"bankName"`
	BankCode     *string    `gorm:"column:bank_code;size:20" json:"bankCode"`
	Branch       *string    `gorm:"column:branch;size:255" json:"branch"`
	AccountNumber string    `gorm:"column:account_number;not null;size:30;uniqueIndex:idx_bank_account_company,unique" json:"accountNumber"`
	AccountName  string     `gorm:"column:account_name;not null;size:255" json:"accountName"`
	Currency     string     `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	IsPrimary    bool       `gorm:"column:is_primary;default:false;index" json:"isPrimary"`
	IsActive     bool       `gorm:"column:is_active;default:true;index" json:"isActive"`
	OpenedAt     *time.Time `gorm:"column:opened_at;type:date" json:"openedAt"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (CompanyBankAccountGORM) TableName() string { return "company_bank_accounts" }

type EInvoicePatternGORM struct {
	ID         string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID  string    `gorm:"column:company_id;not null;size:36;index:idx_einv_company_pattern,unique" json:"companyId"`
	Pattern    string    `gorm:"column:pattern;not null;size:20;index:idx_einv_company_pattern,unique" json:"pattern"`
	Serial     string    `gorm:"column:serial;not null;size:20" json:"serial"`
	InvoiceType string   `gorm:"column:invoice_type;not null;size:10" json:"invoiceType"`
	IsActive   bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	ExpiredAt  *time.Time `gorm:"column:expired_at;type:date" json:"expiredAt"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (EInvoicePatternGORM) TableName() string { return "einvoice_patterns" }

type DigitalSignatureGORM struct {
	ID              string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	Name            string     `gorm:"column:name;not null;size:255" json:"name"`
	SerialNumber    string     `gorm:"column:serial_number;not null;size:100" json:"serialNumber"`
	Issuer          *string    `gorm:"column:issuer;size:255" json:"issuer"`
	ValidFrom       time.Time  `gorm:"column:valid_from;not null;type:date" json:"validFrom"`
	ValidTo         time.Time  `gorm:"column:valid_to;not null;type:date" json:"validTo"`
	CertificateData *string    `gorm:"column:certificate_data;type:text" json:"-"`
	IsActive        bool       `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (DigitalSignatureGORM) TableName() string { return "digital_signatures" }

type IntegrationProfileGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index:idx_integration_company_type,unique" json:"companyId"`
	Type        string    `gorm:"column:type;not null;size:30;index:idx_integration_company_type,unique" json:"type"`
	Name        string    `gorm:"column:name;not null;size:255" json:"name"`
	IsActive    bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	Config      *string   `gorm:"column:config;type:jsonb" json:"-"`
	LastSyncAt  *time.Time `gorm:"column:last_sync_at" json:"lastSyncAt"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (IntegrationProfileGORM) TableName() string { return "integration_profiles" }
