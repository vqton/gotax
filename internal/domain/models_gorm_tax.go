package domain

import "time"

type TaxDeclarationGORM struct {
	ID                  string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID           string    `gorm:"column:company_id;not null;size:36;index:idx_taxdecl_company_type,unique" json:"companyId"`
	DeclarationType     string    `gorm:"column:declaration_type;not null;size:20;index:idx_taxdecl_company_type,unique" json:"declarationType"`
	TaxPeriod           TaxPeriodGORM `gorm:"embedded;embeddedPrefix:tax_period_" json:"taxPeriod"`
	Status              string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	AdjustmentType      string    `gorm:"column:adjustment_type;not null;size:20;default:'REGULAR'" json:"adjustmentType"`
	Version             int       `gorm:"column:version;not null;default:1" json:"version"`
	SubmittedAt         *time.Time `gorm:"column:submitted_at" json:"submittedAt"`
	SubmittedBy         *string   `gorm:"column:submitted_by;size:36" json:"submittedBy"`
	AcknowledgedAt      *time.Time `gorm:"column:acknowledged_at" json:"acknowledgedAt"`
	AcknowledgementRef  *string   `gorm:"column:acknowledgement_ref;size:100" json:"acknowledgementRef"`
	DeclarationXML      *string   `gorm:"column:declaration_xml;type:text" json:"-"`
	GDTResponseXML      *string   `gorm:"column:gdt_response_xml;type:text" json:"-"`
	PreviousDeclID      *string   `gorm:"column:previous_declaration_id;size:36" json:"previousDeclarationId"`
	CreatedBy           string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Lines               []TaxDeclarationLineGORM `gorm:"foreignKey:DeclarationID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (TaxDeclarationGORM) TableName() string { return "tax_declarations" }

type TaxPeriodGORM struct {
	PeriodYear  int    `gorm:"column:period_year;not null" json:"periodYear"`
	PeriodNumber int   `gorm:"column:period_number;not null" json:"periodNumber"`
	PeriodType  string `gorm:"column:period_type;not null;size:10" json:"periodType"`
}

type TaxDeclarationLineGORM struct {
	ID            string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	DeclarationID string    `gorm:"column:declaration_id;not null;size:36;index:idx_taxdecl_line,unique" json:"declarationId"`
	LineCode      string    `gorm:"column:line_code;not null;size:20;index:idx_taxdecl_line,unique" json:"lineCode"`
	LineName      string    `gorm:"column:line_name;not null;size:255" json:"lineName"`
	Amount        float64   `gorm:"column:amount;not null" json:"amount"`
	SourceType    string    `gorm:"column:source_type;not null;size:30" json:"sourceType"`
	SortOrder     int       `gorm:"column:sort_order;not null" json:"sortOrder"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (TaxDeclarationLineGORM) TableName() string { return "tax_declaration_lines" }

type TaxRateGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	TaxType      string    `gorm:"column:tax_type;not null;size:20;index" json:"taxType"`
	RateCode     string    `gorm:"column:rate_code;not null;size:20;uniqueIndex" json:"rateCode"`
	RateName     string    `gorm:"column:rate_name;not null;size:255" json:"rateName"`
	RateType     string    `gorm:"column:rate_type;not null;size:30" json:"rateType"`
	RateValue    *float64  `gorm:"column:rate_value" json:"rateValue"`
	EffectiveFrom time.Time `gorm:"column:effective_from;not null;type:date" json:"effectiveFrom"`
	EffectiveTo  *string   `gorm:"column:effective_to;type:date" json:"effectiveTo"`
	IsActive     bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	LegalRef     *string   `gorm:"column:legal_ref;size:255" json:"legalRef"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (TaxRateGORM) TableName() string { return "tax_rates" }

type TaxPaymentGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	DeclarationID   *string   `gorm:"column:declaration_id;size:36;index" json:"declarationId"`
	TaxType         string    `gorm:"column:tax_type;not null;size:20" json:"taxType"`
	PeriodYear      int       `gorm:"column:period_year;not null" json:"periodYear"`
	PeriodNumber    int       `gorm:"column:period_number;not null" json:"periodNumber"`
	DeclaredAmount  float64   `gorm:"column:declared_amount;not null" json:"declaredAmount"`
	PaidAmount      float64   `gorm:"column:paid_amount;not null;default:0" json:"paidAmount"`
	PaymentDate     *string   `gorm:"column:payment_date;type:date" json:"paymentDate"`
	DueDate         time.Time `gorm:"column:due_date;not null;type:date" json:"dueDate"`
	PaymentRef      *string   `gorm:"column:payment_ref;size:100" json:"paymentRef"`
	PaymentMethod   *string   `gorm:"column:payment_method;size:50" json:"paymentMethod"`
	Status          string    `gorm:"column:status;not null;size:20;default:'PENDING';index" json:"status"`
	LateDays        int       `gorm:"column:late_days;default:0" json:"lateDays"`
	LateInterest    float64   `gorm:"column:late_interest;default:0" json:"lateInterest"`
	Notes           *string   `gorm:"column:notes;type:text" json:"notes"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (TaxPaymentGORM) TableName() string { return "tax_payments" }
