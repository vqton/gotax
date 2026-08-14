package domain

import "time"

// CompanyGORM mirrors migrations/000003_company_schema.up.sql (companies table).
type CompanyGORM struct {
	ID                   string    `gorm:"column:id;primaryKey;size:20"`
	TenantID             string    `gorm:"column:tenant_id;not null;size:50"`
	TaxCode              string    `gorm:"column:tax_code;not null;size:20"`
	Name                 string    `gorm:"column:legal_name_vn;not null;size:255"`
	NameEn               *string   `gorm:"column:legal_name_en;size:255"`
	ShortName            *string   `gorm:"column:short_name;size:100"`
	LegalForm            string    `gorm:"column:legal_form;not null;size:20"`
	Address              string    `gorm:"column:reg_address;not null;size:255"`
	City                 *string   `gorm:"column:reg_province;size:100"`
	Province             *string   `gorm:"column:reg_district;size:100"`
	Phone                *string   `gorm:"column:phone;size:30"`
	Email                *string   `gorm:"column:email;size:100"`
	Website              *string   `gorm:"column:website;size:255"`
	LegalRepName         *string   `gorm:"column:legal_rep_name;size:100"`
	LegalRepTitle        *string   `gorm:"column:legal_rep_title;size:100"`
	LegalRepID           *string   `gorm:"column:legal_rep_id_number;size:30"`
	ChiefAccountant      *string   `gorm:"column:chief_accountant;size:100"`
	ChiefAccountantEmail *string   `gorm:"column:chief_accountant_email;size:100"`
	TaxOfficeCode        *string   `gorm:"column:tax_office_code;size:20"`
	TaxOffice            *string   `gorm:"column:tax_office_name;size:255"`
	AccountingRegime     string    `gorm:"column:accounting_regime;not null;size:10"`
	FiscalYearStartMonth int       `gorm:"column:fiscal_year_start_month;not null;default:1"`
	DefaultCurrency      string    `gorm:"column:default_currency;not null;size:3;default:VND"`
	Status               string    `gorm:"column:status;not null;size:15;default:ACTIVE"`
	ParentCompanyID      *string   `gorm:"column:parent_company_id;size:20"`
	Settings             *string   `gorm:"column:settings;type:jsonb"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Branches             []CompanyBranchGORM   `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"branches,omitempty"`
	FiscalYears          []FiscalYearGORM      `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"fiscalYears,omitempty"`
}

func (CompanyGORM) TableName() string { return "companies" }

type CompanyBranchGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID    string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	Code         string    `gorm:"column:code;not null;size:20;uniqueIndex:idx_branch_company_code,unique" json:"code"`
	Name         string    `gorm:"column:name;not null;size:255" json:"name"`
	Address      *string   `gorm:"column:address;type:text" json:"address"`
	City         *string   `gorm:"column:city;size:100" json:"city"`
	Province     *string   `gorm:"column:province;size:100" json:"province"`
	Phone        *string   `gorm:"column:phone;size:20" json:"phone"`
	ManagerName  *string   `gorm:"column:manager_name;size:255" json:"managerName"`
	IsActive     bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	DeactivatedAt *time.Time `gorm:"column:deactivated_at" json:"deactivatedAt"`
	DeactivationReason *string `gorm:"column:deactivation_reason;type:text" json:"deactivationReason"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (CompanyBranchGORM) TableName() string { return "company_branches" }

type FiscalYearGORM struct {
	ID         string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID  string    `gorm:"column:company_id;not null;size:36;index:idx_fy_company_year,unique" json:"companyId"`
	Year       int       `gorm:"column:year;not null;index:idx_fy_company_year,unique" json:"year"`
	StartDate  time.Time `gorm:"column:start_date;not null;type:date" json:"startDate"`
	EndDate    time.Time `gorm:"column:end_date;not null;type:date" json:"endDate"`
	IsClosed   bool      `gorm:"column:is_closed;default:false;index" json:"isClosed"`
	ClosedAt   *time.Time `gorm:"column:closed_at" json:"closedAt"`
	ClosedBy   *string   `gorm:"column:closed_by;size:36" json:"closedBy"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Periods    []PeriodV2GORM `gorm:"foreignKey:FiscalYearID;constraint:OnDelete:CASCADE" json:"periods,omitempty"`
}

func (FiscalYearGORM) TableName() string { return "fiscal_years" }

type PeriodV2GORM struct {
	ID             string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	FiscalYearID   string     `gorm:"column:fiscal_year_id;not null;size:36;index:idx_period_fy,unique" json:"fiscalYearId"`
	PeriodNumber   int        `gorm:"column:period_number;not null;index:idx_period_fy,unique" json:"periodNumber"`
	StartDate      time.Time  `gorm:"column:start_date;not null;type:date" json:"startDate"`
	EndDate        time.Time  `gorm:"column:end_date;not null;type:date" json:"endDate"`
	Status         string     `gorm:"column:status;not null;size:20;default:'OPEN';index" json:"status"`
	ClosedAt       *time.Time `gorm:"column:closed_at" json:"closedAt"`
	ClosedBy       *string    `gorm:"column:closed_by;size:36" json:"closedBy"`
	ReopenCount    int        `gorm:"column:reopen_count;default:0" json:"reopenCount"`
	LastReopenedAt *time.Time `gorm:"column:last_reopened_at" json:"lastReopenedAt"`
	LastReopenedBy *string    `gorm:"column:last_reopened_by;size:36" json:"lastReopenedBy"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (PeriodV2GORM) TableName() string { return "period_v2" }

type DepartmentGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	Code        string    `gorm:"column:code;not null;size:20;uniqueIndex:idx_dept_company_code,unique" json:"code"`
	Name        string    `gorm:"column:name;not null;size:255" json:"name"`
	Description *string   `gorm:"column:description;type:text" json:"description"`
	ManagerID   *string   `gorm:"column:manager_id;size:36" json:"managerId"`
	IsActive    bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (DepartmentGORM) TableName() string { return "departments" }

type EmployeeGORM struct {
	ID        string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID string    `gorm:"column:company_id;not null;size:36;index:idx_emp_company_code,unique" json:"companyId"`
	Code      string    `gorm:"column:code;not null;size:30;index:idx_emp_company_code,unique" json:"code"`
	FullName  string    `gorm:"column:full_name;not null;size:255" json:"fullName"`
	IDNumber  *string   `gorm:"column:id_number;size:20;uniqueIndex" json:"idNumber"`
	Email     *string   `gorm:"column:email;size:255" json:"email"`
	Phone     *string   `gorm:"column:phone;size:20" json:"phone"`
	DeptID    *string   `gorm:"column:dept_id;size:36;index" json:"deptId"`
	Position  *string   `gorm:"column:position;size:100" json:"position"`
	JoinDate  *time.Time `gorm:"column:join_date;type:date" json:"joinDate"`
	IsActive  bool      `gorm:"column:is_active;default:true;index" json:"isActive"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (EmployeeGORM) TableName() string { return "employees" }
