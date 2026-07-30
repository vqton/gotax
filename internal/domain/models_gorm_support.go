package domain

import "time"

type PeriodGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	Year        int       `gorm:"column:year;not null;index:idx_period_year_month,unique" json:"year"`
	Month       int       `gorm:"column:month;not null;index:idx_period_year_month,unique" json:"month"`
	StartDate   time.Time `gorm:"column:start_date;not null;type:date" json:"startDate"`
	EndDate     time.Time `gorm:"column:end_date;not null;type:date" json:"endDate"`
	Status      string    `gorm:"column:status;not null;size:20;default:'OPEN';index" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (PeriodGORM) TableName() string { return "periods" }

type UserGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	Username     string    `gorm:"column:username;uniqueIndex;not null;size:100" json:"username"`
	PasswordHash string    `gorm:"column:password_hash;not null;size:255" json:"-"`
	FullName     string    `gorm:"column:full_name;not null;size:255" json:"fullName"`
	Email        string    `gorm:"column:email;size:255" json:"email"`
	Role         string    `gorm:"column:role;not null;size:20;index" json:"role"`
	IsActive     bool      `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (UserGORM) TableName() string { return "users" }

type RefreshTokenGORM struct {
	ID         string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	UserID     string     `gorm:"column:user_id;not null;size:36;index" json:"userId"`
	TokenHash  string     `gorm:"column:token_hash;not null;size:64;uniqueIndex" json:"-"`
	DeviceInfo *string    `gorm:"column:device_info;size:255" json:"deviceInfo"`
	IPAddress  *string    `gorm:"column:ip_address;size:45" json:"ipAddress"`
	ExpiresAt  time.Time  `gorm:"column:expires_at;not null;index" json:"expiresAt"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	RevokedAt  *time.Time `gorm:"column:revoked_at" json:"revokedAt"`
}

func (RefreshTokenGORM) TableName() string { return "refresh_tokens" }

type PasswordResetTokenGORM struct {
	ID        string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	UserID    string     `gorm:"column:user_id;not null;size:36;index" json:"userId"`
	TokenHash string     `gorm:"column:token_hash;not null;size:64;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null" json:"expiresAt"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"usedAt"`
}

func (PasswordResetTokenGORM) TableName() string { return "password_reset_tokens" }

type AuditEntryGORM struct {
	ID         string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	UserID     *string   `gorm:"column:user_id;size:36;index" json:"userId"`
	Username   string    `gorm:"column:username;size:100" json:"username"`
	IPAddress  *string   `gorm:"column:ip_address;size:45" json:"ipAddress"`
	Action     string    `gorm:"column:action;not null;size:30;index" json:"action"`
	EntityType string    `gorm:"column:entity_type;not null;size:50;index" json:"entityType"`
	EntityID   *string   `gorm:"column:entity_id;size:36;index" json:"entityId"`
	OldValue   *string   `gorm:"column:old_value;type:jsonb" json:"oldValue"`
	NewValue   *string   `gorm:"column:new_value;type:jsonb" json:"newValue"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index" json:"createdAt"`
}

func (AuditEntryGORM) TableName() string { return "audit_log" }

type ExchangeRateGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CurrencyCode string    `gorm:"column:currency_code;not null;size:3;uniqueIndex:idx_rate_currency_date,unique" json:"currencyCode"`
	RateDate     time.Time `gorm:"column:rate_date;not null;type:date;uniqueIndex:idx_rate_currency_date,unique" json:"rateDate"`
	BuyRate      *float64  `gorm:"column:buy_rate" json:"buyRate"`
	SellRate     *float64  `gorm:"column:sell_rate" json:"sellRate"`
	AverageRate  float64   `gorm:"column:average_rate;not null" json:"averageRate"`
	Source       *string   `gorm:"column:source;size:100" json:"source"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (ExchangeRateGORM) TableName() string { return "exchange_rates" }

type ClosingTemplateGORM struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"column:name;not null;size:255;uniqueIndex" json:"name"`
	Description   *string   `gorm:"column:description;type:text" json:"description"`
	SequenceOrder int       `gorm:"column:sequence_order;not null;uniqueIndex" json:"sequenceOrder"`
	IsActive      bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (ClosingTemplateGORM) TableName() string { return "closing_templates" }
