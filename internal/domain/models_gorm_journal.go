package domain

import "time"

type JournalEntryGORM struct {
	ID              string    `gorm:"primaryKey;size:20" json:"id"`
	CompanyID       string    `gorm:"column:company_id;size:20;index" json:"companyId"`
	EntryNumber     string    `gorm:"column:entry_number;not null;size:30;uniqueIndex:idx_entry_number" json:"entryNumber"`
	VoucherType     string    `gorm:"column:voucher_type;not null;size:20" json:"voucherType"`
	EntryDate       time.Time `gorm:"column:entry_date;not null;type:date;index" json:"entryDate"`
	AccountingDate  time.Time `gorm:"column:accounting_date;not null;type:date" json:"accountingDate"`
	PeriodID        string    `gorm:"column:period_id;size:20;index" json:"periodId"`
	Description     string    `gorm:"column:description;type:text" json:"description"`
	Status          string    `gorm:"column:status;not null;size:20;index" json:"status"`
	CurrencyCode    string    `gorm:"column:currency_code;size:3" json:"currencyCode"`
	ExchangeRate    float64   `gorm:"column:exchange_rate;default:1" json:"exchangeRate"`
	CreatedBy       string    `gorm:"column:created_by;size:36;index;default:null" json:"createdBy"`
	ReviewedBy      string    `gorm:"column:reviewed_by;size:36;default:null" json:"reviewedBy"`
	ApprovedBy      string    `gorm:"column:approved_by;size:36;default:null" json:"approvedBy"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	PostedAt        *time.Time `gorm:"column:posted_at" json:"postedAt"`
	ApprovedAt      *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	Lines           []JournalLineGORM `gorm:"foreignKey:EntryID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (JournalEntryGORM) TableName() string { return "journal_entries" }

type JournalLineGORM struct {
	ID            string    `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	EntryID       string    `gorm:"column:entry_id;not null;index:idx_entry_lines,unique;constraint:OnDelete:CASCADE" json:"entryId"`
	LineNumber    int       `gorm:"column:line_number;not null;index:idx_entry_lines,unique" json:"lineNumber"`
	AccountCode   string    `gorm:"column:account_code;not null;size:20;index" json:"accountCode"`
	DebitAmount   float64   `gorm:"column:debit_amount;default:0" json:"debitAmount"`
	CreditAmount  float64   `gorm:"column:credit_amount;default:0" json:"creditAmount"`
	Description   string    `gorm:"column:description;type:text" json:"description"`
	CurrencyCode  string    `gorm:"column:currency_code;size:3" json:"currencyCode"`
	ForeignAmount float64   `gorm:"column:foreign_amount;default:0" json:"foreignAmount"`
	ExchangeRate  float64   `gorm:"column:exchange_rate;default:1" json:"exchangeRate"`
	ObjectID      string    `gorm:"column:object_id;size:36" json:"objectId"`
	ProjectID     string    `gorm:"column:project_id;size:36" json:"projectId"`
	ContractID    string    `gorm:"column:contract_id;size:36" json:"contractId"`
	CostItemID    string    `gorm:"column:cost_item_id;size:36" json:"costItemId"`
	DepartmentID  string    `gorm:"column:department_id;size:36" json:"departmentId"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (JournalLineGORM) TableName() string { return "journal_lines" }
