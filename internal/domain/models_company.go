package domain

import "time"

type CompanyStatus string
const (CompanyStatusActive CompanyStatus="ACTIVE"; CompanyStatusSuspended CompanyStatus="SUSPENDED"; CompanyStatusDissolved CompanyStatus="DISSOLVED"; CompanyStatusMerged CompanyStatus="MERGED")

type LegalForm string
const (LegalFormLLC1Member LegalForm="LLC_1MEMBER"; LegalFormLLC2Members LegalForm="LLC_2MEMBERS"; LegalFormJSC LegalForm="JSC"; LegalFormSoleProp LegalForm="SOLE_PROP"; LegalFormPartnership LegalForm="PARTNERSHIP"; LegalFormFO LegalForm="FO"; LegalFormRO LegalForm="RO")

type CompanyType string
const (CompanyTypeManufacturing CompanyType="MANUFACTURING"; CompanyTypeTrading CompanyType="TRADING"; CompanyTypeService CompanyType="SERVICE"; CompanyTypeConstruction CompanyType="CONSTRUCTION"; CompanyTypeAgriculture CompanyType="AGRICULTURE"; CompanyTypeOther CompanyType="OTHER")

type CompanySize string
const (CompanySizeMicro CompanySize="MICRO"; CompanySizeSmall CompanySize="SMALL"; CompanySizeMedium CompanySize="MEDIUM"; CompanySizeLarge CompanySize="LARGE")

type AccountingRegime string
const (RegimeTT99 AccountingRegime="TT99"; RegimeTT133 AccountingRegime="TT133"; RegimeTT58 AccountingRegime="TT58")

type BranchType string
const (BranchTypeBranch BranchType="BRANCH"; BranchTypeRepOffice BranchType="REP_OFFICE"; BranchTypeDependentUnit BranchType="DEPENDENT_UNIT")

type PeriodStatusV2 string
const (PeriodV2Open PeriodStatusV2="OPEN"; PeriodV2Closed PeriodStatusV2="CLOSED"; PeriodV2PermanentlyClosed PeriodStatusV2="PERMANENTLY_CLOSED"; PeriodV2Future PeriodStatusV2="FUTURE")

type EmployeeStatus string
const (EmployeeActive EmployeeStatus="ACTIVE"; EmployeeLeave EmployeeStatus="LEAVE"; EmployeeTerminated EmployeeStatus="TERMINATED")

type BankAccountStatus string
const (BankAccountActive BankAccountStatus="ACTIVE"; BankAccountClosed BankAccountStatus="CLOSED")

type EInvoiceStatus string
const (EInvoiceRegistered EInvoiceStatus="REGISTERED"; EInvoiceActive EInvoiceStatus="ACTIVE"; EInvoiceCancelled EInvoiceStatus="CANCELLED")

type SignatureType string
const (SignatureUSBToken SignatureType="USB_TOKEN"; SignatureRemoteHSM SignatureType="REMOTE_HSM")

type SignatureStatus string
const (SignatureActive SignatureStatus="ACTIVE"; SignatureExpired SignatureStatus="EXPIRED"; SignatureRevoked SignatureStatus="REVOKED")

type IntegrationType string
const (IntegrationGDT IntegrationType="GDT"; IntegrationCustoms IntegrationType="CUSTOMS"; IntegrationBHXH IntegrationType="BHXH"; IntegrationDVC IntegrationType="DVC")

type IntegrationStatus string
const (IntegrationConnected IntegrationStatus="CONNECTED"; IntegrationDisconnected IntegrationStatus="DISCONNECTED"; IntegrationError IntegrationStatus="ERROR")

type PeriodTypeV2 string
const (PeriodTypeMonthly PeriodTypeV2="MONTHLY"; PeriodTypeQuarterly PeriodTypeV2="QUARTERLY"; PeriodTypeAnnual PeriodTypeV2="ANNUAL"; PeriodTypePerOccurrence PeriodTypeV2="PER_OCCURRENCE")

type Company struct {
	ID string `json:"id"`
	TenantID string `json:"tenant_id,omitempty"`
	LegalNameVN string `json:"legal_name_vn"`
	LegalNameEN string `json:"legal_name_en,omitempty"`
	ShortName string `json:"short_name,omitempty"`
	LegalForm LegalForm `json:"legal_form"`
	TaxCode string `json:"tax_code"`
	BusinessRegNo string `json:"business_reg_no,omitempty"`
	BusinessRegDate string `json:"business_reg_date,omitempty"`
	BusinessRegPlace string `json:"business_reg_place,omitempty"`
	RegAddress string `json:"reg_address"`
	RegProvince string `json:"reg_province,omitempty"`
	RegDistrict string `json:"reg_district,omitempty"`
	HeadOfficeAddress string `json:"head_office_address,omitempty"`
	HeadOfficeProvince string `json:"head_office_province,omitempty"`
	HeadOfficeDistrict string `json:"head_office_district,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	Website string `json:"website,omitempty"`
	LegalRepName string `json:"legal_rep_name,omitempty"`
	LegalRepTitle string `json:"legal_rep_title,omitempty"`
	LegalRepIDNumber string `json:"legal_rep_id_number,omitempty"`
	ChiefAccountant string `json:"chief_accountant,omitempty"`
	ChiefAccountantEmail string `json:"chief_accountant_email,omitempty"`
	TaxOfficeCode string `json:"tax_office_code,omitempty"`
	TaxOfficeName string `json:"tax_office_name,omitempty"`
	AccountingRegime AccountingRegime `json:"accounting_regime"`
	FiscalYearStartMonth int `json:"fiscal_year_start_month"`
	DefaultCurrency string `json:"default_currency"`
	SecondaryCurrency string `json:"secondary_currency,omitempty"`
	CompanyType CompanyType `json:"company_type,omitempty"`
	CompanySize CompanySize `json:"company_size,omitempty"`
	Status CompanyStatus `json:"status"`
	ParentCompanyID string `json:"parent_company_id,omitempty"`
	LogoURL string `json:"logo_url,omitempty"`
	Settings map[string]any `json:"settings,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Company) Validate() error {
	if c.TaxCode=="" { return ErrCompanyCodeRequired }
	if len(c.TaxCode)!=10&&len(c.TaxCode)!=13&&len(c.TaxCode)!=14 { return ErrInvalidTaxCodeFormat }
	if c.LegalNameVN=="" { return ErrCompanyNameRequired }
	switch c.LegalForm{case LegalFormLLC1Member,LegalFormLLC2Members,LegalFormJSC,LegalFormSoleProp,LegalFormPartnership,LegalFormFO,LegalFormRO:default: return ErrCompanyInvalidLegalForm}
	switch c.AccountingRegime{case RegimeTT99,RegimeTT133,RegimeTT58:default: return ErrCompanyInvalidRegime}
	if c.Status=="" { c.Status=CompanyStatusActive }
	if c.DefaultCurrency=="" { c.DefaultCurrency="VND" }
	if c.FiscalYearStartMonth==0 { c.FiscalYearStartMonth=1 }
	return nil
}

type CompanyBranch struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BranchName string `json:"branch_name"`
	BranchTaxCode string `json:"branch_tax_code"`
	BranchType BranchType `json:"branch_type"`
	Address string `json:"address,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	ManagerName string `json:"manager_name,omitempty"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type FiscalYear struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	Year int `json:"year"`
	StartMonth int `json:"start_month"`
	IsShortYear bool `json:"is_short_year"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Status string `json:"status"`
	Periods []PeriodV2 `json:"periods,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type PeriodV2 struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	FiscalYearID string `json:"fiscal_year_id"`
	PeriodType PeriodTypeV2 `json:"period_type"`
	PeriodNumber int `json:"period_number"`
	Label string `json:"label"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Status PeriodStatusV2 `json:"status"`
	OpenedAt string `json:"opened_at,omitempty"`
	ClosedAt string `json:"closed_at,omitempty"`
	ClosedBy string `json:"closed_by,omitempty"`
	ReopenedCount int `json:"reopened_count,omitempty"`
}
type Department struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	Code string `json:"code"`
	Name string `json:"name"`
	ParentID string `json:"parent_id,omitempty"`
	ManagerID string `json:"manager_id,omitempty"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Employee struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	EmployeeCode string `json:"employee_code"`
	FullName string `json:"full_name"`
	Title string `json:"title,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	DepartmentID string `json:"department_id,omitempty"`
	PersonalTaxCode string `json:"personal_tax_code,omitempty"`
	SocialInsuranceNo string `json:"social_insurance_no,omitempty"`
	BankAccountNo string `json:"bank_account_no,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Status EmployeeStatus `json:"status"`
	HireDate string `json:"hire_date,omitempty"`
	TerminationDate string `json:"termination_date,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CompanyBankAccount struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BankCode string `json:"bank_code,omitempty"`
	BankName string `json:"bank_name"`
	BranchName string `json:"branch_name,omitempty"`
	AccountNumber string `json:"account_number"`
	AccountHolder string `json:"account_holder"`
	Currency string `json:"currency"`
	IsDefault bool `json:"is_default"`
	IsVerified bool `json:"is_verified"`
	Status BankAccountStatus `json:"status"`
	IBAN string `json:"iban,omitempty"`
	SWIFTCode string `json:"swift_code,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EInvoicePattern struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	PatternCode string `json:"pattern_code"`
	Serial string `json:"serial"`
	Form string `json:"form,omitempty"`
	InvoiceType string `json:"invoice_type,omitempty"`
	Status EInvoiceStatus `json:"status"`
	GDTStatus string `json:"gdt_status,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type DigitalSignature struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	SignatureType SignatureType `json:"signature_type"`
	Provider string `json:"provider,omitempty"`
	SerialNumber string `json:"serial_number"`
	OwnerName string `json:"owner_name,omitempty"`
	CertificateSubject string `json:"certificate_subject,omitempty"`
	CertificateIssuer string `json:"certificate_issuer,omitempty"`
	ValidFrom string `json:"valid_from"`
	ValidTo string `json:"valid_to"`
	Status SignatureStatus `json:"status"`
	IsDefault bool `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type IntegrationProfile struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	IntegrationType IntegrationType `json:"integration_type"`
	EndpointURL string `json:"endpoint_url,omitempty"`
	Status IntegrationStatus `json:"status"`
	LastConnectedAt string `json:"last_connected_at,omitempty"`
	LastErrorAt string `json:"last_error_at,omitempty"`
	LastErrorMsg string `json:"last_error_msg,omitempty"`
	Config map[string]any `json:"config,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CompanyContext struct {
	CompanyID string `json:"company_id"`
	LegalNameVN string `json:"legal_name_vn"`
	TaxCode string `json:"tax_code"`
	AccountingRegime string `json:"accounting_regime"`
	DefaultCurrency string `json:"default_currency"`
	FiscalYear int `json:"fiscal_year"`
	CurrentPeriodID string `json:"current_period_id"`
	CurrentPeriodLabel string `json:"current_period_label"`
	Permissions []string `json:"permissions"`
}
