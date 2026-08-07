package domain

import (
	"context"
	"time"
)

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByCode(ctx context.Context, code string) (*Account, error)
	GetAll(ctx context.Context, activeOnly bool) ([]Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, code string) error
	GetChildren(ctx context.Context, parentCode string) ([]Account, error)
}

type JournalRepository interface {
	Create(ctx context.Context, entry *JournalEntry) error
	GetByID(ctx context.Context, id string) (*JournalEntry, error)
	GetByPeriod(ctx context.Context, periodID string) ([]JournalEntry, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]JournalEntry, error)
	GetByStatus(ctx context.Context, status JournalEntryStatus) ([]JournalEntry, error)
	GetByVoucherType(ctx context.Context, voucherType VoucherType) ([]JournalEntry, error)
	UpdateStatus(ctx context.Context, id string, status JournalEntryStatus) error
	Update(ctx context.Context, entry *JournalEntry) error
	Approve(ctx context.Context, id, approvedBy string) error
	Review(ctx context.Context, id, reviewedBy string) error
	GetLinesByEntryID(ctx context.Context, entryID string) ([]JournalLine, error)
	GetBalance(ctx context.Context, accountCode string, periodID string) (*AccountBalance, error)
	GetTrialBalance(ctx context.Context, periodID string) ([]AccountBalance, error)
	GetFinancialStatement(ctx context.Context, periodID string, accountTypes []AccountType) ([]AccountBalance, error)
	GetAccountUsage(ctx context.Context, accountCode string) (*AccountUsage, error)
	GetPostedEntriesByAccount(ctx context.Context, periodID, accountCode string) ([]JournalEntry, error)
}

type PeriodRepository interface {
	Create(ctx context.Context, period *Period) error
	GetByID(ctx context.Context, id string) (*Period, error)
	GetByYearMonth(ctx context.Context, year, month int) (*Period, error)
	GetAll(ctx context.Context) ([]Period, error)
	UpdateStatus(ctx context.Context, id string, status PeriodStatus) error
	GetOpenPeriod(ctx context.Context) (*Period, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByID(ctx context.Context, id string) (*RefreshToken, error)
	GetByUserID(ctx context.Context, userID string) ([]RefreshToken, error)
	GetByHash(ctx context.Context, hash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllByUserID(ctx context.Context, userID string) error
}

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *PasswordResetToken) error
	GetByID(ctx context.Context, id string) (*PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
}

type AuditLogRepository interface {
	Create(ctx context.Context, entry *AuditEntry) error
	GetByEntity(ctx context.Context, entityType, entityID string) ([]AuditEntry, error)
	GetByUser(ctx context.Context, userID string) ([]AuditEntry, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]AuditEntry, error)
	GetAll(ctx context.Context, limit int) ([]AuditEntry, error)
}

type ExchangeRateRepository interface {
	Create(ctx context.Context, rate *ExchangeRate) error
	GetByCurrencyAndDate(ctx context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]ExchangeRate, error)
	GetAll(ctx context.Context) ([]ExchangeRate, error)
	Delete(ctx context.Context, id string) error
}

type ClosingTemplateRepository interface {
	Create(ctx context.Context, template *ClosingTemplate) error
	GetByID(ctx context.Context, id string) (*ClosingTemplate, error)
	GetAll(ctx context.Context) ([]ClosingTemplate, error)
	Update(ctx context.Context, template *ClosingTemplate) error
	Delete(ctx context.Context, id string) error
}

type ApprovalRepository interface {
	Create(ctx context.Context, req *ApprovalRequest) error
	GetByID(ctx context.Context, id string) (*ApprovalRequest, error)
	GetByStatus(ctx context.Context, status ApprovalStatus) ([]ApprovalRequest, error)
	GetByEntity(ctx context.Context, entityType, entityID string) ([]ApprovalRequest, error)
	UpdateStatus(ctx context.Context, id string, status ApprovalStatus, reviewedBy, reviewNote string) error
	GetAll(ctx context.Context) ([]ApprovalRequest, error)
}

type AccountVersionRepository interface {
	Create(ctx context.Context, version *AccountVersion) error
	GetByVersionNumber(ctx context.Context, versionNumber string) (*AccountVersion, error)
	GetLatest(ctx context.Context) (*AccountVersion, error)
	GetAll(ctx context.Context) ([]AccountVersion, error)
}

type AccountMappingRepository interface {
	Create(ctx context.Context, mapping *AccountMapping) error
	GetByOldCode(ctx context.Context, sourceRegime, oldCode string) (*AccountMapping, error)
	GetByRegime(ctx context.Context, sourceRegime, targetRegime string) ([]AccountMapping, error)
	GetAll(ctx context.Context) ([]AccountMapping, error)
}

type AccountAnalysisRepository interface {
	Create(ctx context.Context, analysis *AccountAnalysis) error
	GetByAccount(ctx context.Context, accountCode string) (*AccountAnalysis, error)
	Update(ctx context.Context, analysis *AccountAnalysis) error
}

type IFRSMappingRepository interface {
	Create(ctx context.Context, mapping *IFRSMapping) error
	GetByVASCode(ctx context.Context, vasCode string) (*IFRSMapping, error)
	GetAll(ctx context.Context) ([]IFRSMapping, error)
	Update(ctx context.Context, mapping *IFRSMapping) error
}

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	GetByID(ctx context.Context, id string) (*Company, error)
	GetByTaxCode(ctx context.Context, tenantID, taxCode string) (*Company, error)
	GetAll(ctx context.Context, tenantID string) ([]Company, error)
	Update(ctx context.Context, company *Company) error
	Deactivate(ctx context.Context, id, reason string) error
	GetHierarchy(ctx context.Context, companyID string) ([]Company, error)

	CreateBranch(ctx context.Context, branch *CompanyBranch) error
	GetBranchByID(ctx context.Context, id string) (*CompanyBranch, error)
	GetBranchesByCompany(ctx context.Context, companyID string) ([]CompanyBranch, error)
	UpdateBranch(ctx context.Context, branch *CompanyBranch) error
	DeactivateBranch(ctx context.Context, id string) error

	CreateFiscalYear(ctx context.Context, fy *FiscalYear) error
	GetFiscalYearByID(ctx context.Context, id string) (*FiscalYear, error)
	GetFiscalYearsByCompany(ctx context.Context, companyID string) ([]FiscalYear, error)
	GetFiscalYearByYear(ctx context.Context, companyID string, year int) (*FiscalYear, error)

	CreatePeriod(ctx context.Context, period *PeriodV2) error
	GetPeriodByID(ctx context.Context, id string) (*PeriodV2, error)
	GetPeriodsByFiscalYear(ctx context.Context, fiscalYearID string) ([]PeriodV2, error)
	GetPeriodsByCompany(ctx context.Context, companyID string) ([]PeriodV2, error)
	GetOpenPeriod(ctx context.Context, companyID string) (*PeriodV2, error)
	UpdatePeriodStatus(ctx context.Context, id string, status PeriodStatusV2, closedBy string) error
	IncrementReopenCount(ctx context.Context, id string) error

	CreateDepartment(ctx context.Context, dept *Department) error
	GetDepartmentByID(ctx context.Context, id string) (*Department, error)
	GetDepartmentsByCompany(ctx context.Context, companyID string) ([]Department, error)
	UpdateDepartment(ctx context.Context, dept *Department) error

	CreateEmployee(ctx context.Context, emp *Employee) error
	GetEmployeeByID(ctx context.Context, id string) (*Employee, error)
	GetEmployeesByCompany(ctx context.Context, companyID string) ([]Employee, error)
	GetEmployeeByCode(ctx context.Context, companyID, code string) (*Employee, error)
	UpdateEmployee(ctx context.Context, emp *Employee) error
	DeactivateEmployee(ctx context.Context, id string) error

	CreateBankAccount(ctx context.Context, ba *CompanyBankAccount) error
	GetBankAccountByID(ctx context.Context, id string) (*CompanyBankAccount, error)
	GetBankAccountsByCompany(ctx context.Context, companyID string) ([]CompanyBankAccount, error)
	UpdateBankAccount(ctx context.Context, ba *CompanyBankAccount) error
	DeactivateBankAccount(ctx context.Context, id string) error

	CreateEInvoicePattern(ctx context.Context, inv *EInvoicePattern) error
	GetEInvoicePatternByID(ctx context.Context, id string) (*EInvoicePattern, error)
	GetEInvoicePatternsByCompany(ctx context.Context, companyID string) ([]EInvoicePattern, error)
	UpdateEInvoicePattern(ctx context.Context, inv *EInvoicePattern) error

	CreateDigitalSignature(ctx context.Context, sig *DigitalSignature) error
	GetDigitalSignatureByID(ctx context.Context, id string) (*DigitalSignature, error)
	GetDigitalSignaturesByCompany(ctx context.Context, companyID string) ([]DigitalSignature, error)
	UpdateDigitalSignature(ctx context.Context, sig *DigitalSignature) error

	CreateIntegrationProfile(ctx context.Context, prof *IntegrationProfile) error
	GetIntegrationProfileByID(ctx context.Context, id string) (*IntegrationProfile, error)
	GetIntegrationProfilesByCompany(ctx context.Context, companyID string) ([]IntegrationProfile, error)
	GetIntegrationByType(ctx context.Context, companyID string, itype IntegrationType) (*IntegrationProfile, error)
	UpdateIntegrationProfile(ctx context.Context, prof *IntegrationProfile) error
}

// ─── Tax Repository ─────────────────────────────────────────────

type TaxRepository interface {
	CreateDeclaration(ctx context.Context, d *TaxDeclaration) error
	GetDeclarationByID(ctx context.Context, id string) (*TaxDeclaration, error)
	GetDeclarations(ctx context.Context, filter TaxDeclarationFilter) ([]TaxDeclaration, error)
	UpdateDeclaration(ctx context.Context, d *TaxDeclaration) error
	UpdateDeclarationStatus(ctx context.Context, id string, status DeclarationStatus) error

	CreateRate(ctx context.Context, r *TaxRate) error
	GetRateByID(ctx context.Context, id string) (*TaxRate, error)
	GetRates(ctx context.Context, filter TaxRateFilter) ([]TaxRate, error)
	UpdateRate(ctx context.Context, r *TaxRate) error

	CreatePayment(ctx context.Context, p *TaxPayment) error
	GetPaymentByID(ctx context.Context, id string) (*TaxPayment, error)
	GetPayments(ctx context.Context, filter PaymentFilter) ([]TaxPayment, error)
	UpdatePayment(ctx context.Context, p *TaxPayment) error

	CreateEInvoice(ctx context.Context, inv *EInvoice) error
	GetEInvoiceByID(ctx context.Context, id string) (*EInvoice, error)
	GetEInvoices(ctx context.Context, filter EInvoiceFilter) ([]EInvoice, error)
	UpdateEInvoice(ctx context.Context, inv *EInvoice) error
	UpdateEInvoiceStatus(ctx context.Context, id string, status EInvLifecycleStatus) error

	CreateCalendarEntry(ctx context.Context, c *TaxCalendar) error
	GetCalendarEntryByID(ctx context.Context, id string) (*TaxCalendar, error)
	GetCalendarByPeriod(ctx context.Context, companyID string, periodYear, periodNumber int) ([]TaxCalendar, error)
	GetCalendarByCompany(ctx context.Context, companyID string) ([]TaxCalendar, error)
	UpdateCalendarStatus(ctx context.Context, id string, status CalendarStatus) error

	CreateAlert(ctx context.Context, a *TaxAlert) error
	GetAlertByID(ctx context.Context, id string) (*TaxAlert, error)
	GetAlerts(ctx context.Context, companyID string, limit int) ([]TaxAlert, error)

	CreateAuditCase(ctx context.Context, a *TaxAuditCase) error
	GetAuditCaseByID(ctx context.Context, id string) (*TaxAuditCase, error)
	GetAuditCases(ctx context.Context, companyID string) ([]TaxAuditCase, error)
	UpdateAuditCase(ctx context.Context, a *TaxAuditCase) error

	CreateCITLoss(ctx context.Context, loss *CITLossCarryForward) error
	GetActiveCITLosses(ctx context.Context, companyID string, beforeYear int) ([]CITLossCarryForward, error)
	UpdateCITLoss(ctx context.Context, loss *CITLossCarryForward) error
}

// ─── Opening Balance Repository ───────────────────────────────────

type OpeningBalanceRepository interface {
	Create(ctx context.Context, ob *OpeningBalance) error
	GetByID(ctx context.Context, id string) (*OpeningBalance, error)
	List(ctx context.Context, filter OBListFilter) ([]OpeningBalance, error)
	GetByAccount(ctx context.Context, companyID, periodID, accountCode string) (*OpeningBalance, error)
	Update(ctx context.Context, ob *OpeningBalance) error
	UpdateStatus(ctx context.Context, id string, status OpeningBalanceStatus, approvedBy string) error
	Delete(ctx context.Context, id string) error

	BulkCreate(ctx context.Context, balances []OpeningBalance) error
	BulkUpdateStatus(ctx context.Context, ids []string, status OpeningBalanceStatus, approvedBy string) error

	CreateDetail(ctx context.Context, d *OpeningBalanceDetail) error
	GetDetails(ctx context.Context, balanceID string) ([]OpeningBalanceDetail, error)
	DeleteDetail(ctx context.Context, id string) error
	DeleteDetails(ctx context.Context, balanceID string) error

	GetTotals(ctx context.Context, companyID, periodID string) (totalDebit, totalCredit float64, err error)
	ValidateBalanced(ctx context.Context, companyID, periodID string) (bool, error)

	CreateCarryForwardLog(ctx context.Context, log *CarryForwardLog) error
	GetCarryForwardLogs(ctx context.Context, companyID string) ([]CarryForwardLog, error)
	GetCarryForwardLogByID(ctx context.Context, id string) (*CarryForwardLog, error)

	CreateCircular99Mapping(ctx context.Context, m *Circular99Mapping) error
	ListCircular99Mappings(ctx context.Context) ([]Circular99Mapping, error)
	GetCircular99MappingByOldCode(ctx context.Context, oldCode string) (*Circular99Mapping, error)

	CreateMigration(ctx context.Context, m *BalanceMigration) error
	GetMigrationByID(ctx context.Context, id string) (*BalanceMigration, error)
	ListMigrations(ctx context.Context, companyID string) ([]BalanceMigration, error)
}

// ─── Bank Repository ───────────────────────────────────────────────────

type BankRepository interface {
	// Statements
	CreateStatement(ctx context.Context, s *BankStatement) error
	GetStatement(ctx context.Context, id string) (*BankStatement, error)
	ListStatements(ctx context.Context, companyID, bankAccountID string, limit, offset int) ([]BankStatement, int, error)
	DeleteStatement(ctx context.Context, id string) error

	CreateStatementLines(ctx context.Context, lines []BankStatementLine) error
	GetStatementLines(ctx context.Context, statementID string) ([]BankStatementLine, error)
	GetStatementLinesByStatus(ctx context.Context, statementID string, status MatchStatus) ([]BankStatementLine, error)
	UpdateStatementLineMatch(ctx context.Context, lineID string, matchStatus MatchStatus, matchedLineID, matchedBy string) error

	// Reconciliation
	CreateReconciliation(ctx context.Context, r *BankReconciliation) error
	GetReconciliation(ctx context.Context, id string) (*BankReconciliation, error)
	ListReconciliations(ctx context.Context, companyID, bankAccountID string) ([]BankReconciliation, error)
	UpdateReconciliation(ctx context.Context, r *BankReconciliation) error

	CreateReconciliationMatch(ctx context.Context, m *BankReconciliationMatch) error
	GetReconciliationMatches(ctx context.Context, reconID string) ([]BankReconciliationMatch, error)
	DeleteReconciliationMatch(ctx context.Context, id string) error

	// Payment Orders
	CreatePaymentOrder(ctx context.Context, po *PaymentOrder) error
	GetPaymentOrder(ctx context.Context, id string) (*PaymentOrder, error)
	ListPaymentOrders(ctx context.Context, filter PaymentOrderFilter) ([]PaymentOrder, int, error)
	UpdatePaymentOrder(ctx context.Context, po *PaymentOrder) error

	CreatePaymentOrderBatch(ctx context.Context, b *PaymentOrderBatch) error
	GetPaymentOrderBatch(ctx context.Context, id string) (*PaymentOrderBatch, error)
	ListPaymentOrderBatches(ctx context.Context, companyID string) ([]PaymentOrderBatch, error)
	UpdatePaymentOrderBatch(ctx context.Context, b *PaymentOrderBatch) error
	AddOrdersToBatch(ctx context.Context, batchID string, orderIDs []string) error
	GetBatchOrderIDs(ctx context.Context, batchID string) ([]string, error)

	// Loans
	CreateLoan(ctx context.Context, l *LoanAgreement) error
	GetLoan(ctx context.Context, id string) (*LoanAgreement, error)
	ListLoans(ctx context.Context, filter LoanFilter) ([]LoanAgreement, error)
	UpdateLoan(ctx context.Context, l *LoanAgreement) error

	CreateDisbursement(ctx context.Context, d *LoanDisbursement) error
	GetDisbursements(ctx context.Context, loanID string) ([]LoanDisbursement, error)

	CreateRepayment(ctx context.Context, r *LoanRepayment) error
	GetRepayments(ctx context.Context, loanID string) ([]LoanRepayment, error)
	UpdateRepayment(ctx context.Context, r *LoanRepayment) error

	// Term Deposits
	CreateDeposit(ctx context.Context, d *TermDeposit) error
	GetDeposit(ctx context.Context, id string) (*TermDeposit, error)
	ListDeposits(ctx context.Context, companyID string) ([]TermDeposit, error)
	UpdateDeposit(ctx context.Context, d *TermDeposit) error

	// Reports
	GetBankLedger(ctx context.Context, companyID, bankAccountID, fromDate, toDate string) (*BankLedger, error)
	GetBalance(ctx context.Context, companyID, bankAccountID string) (float64, error)
}

// ─── Cash Repository ───────────────────────────────────────────────────

type CashRepository interface {
	CreateReceipt(ctx context.Context, r *CashReceipt) error
	GetReceipt(ctx context.Context, id string) (*CashReceipt, error)
	ListReceipts(ctx context.Context, filter CashReceiptFilter) ([]CashReceipt, int, error)
	UpdateReceipt(ctx context.Context, r *CashReceipt) error
	DeleteReceipt(ctx context.Context, id string) error
	LastReceiptNo(ctx context.Context, companyID, year string) (string, error)

	CreatePayment(ctx context.Context, p *CashPayment) error
	GetPayment(ctx context.Context, id string) (*CashPayment, error)
	ListPayments(ctx context.Context, filter CashPaymentFilter) ([]CashPayment, int, error)
	UpdatePayment(ctx context.Context, p *CashPayment) error
	DeletePayment(ctx context.Context, id string) error
	LastPaymentNo(ctx context.Context, companyID, year string) (string, error)

	CreateTransfer(ctx context.Context, t *CashTransfer) error
	GetTransfer(ctx context.Context, id string) (*CashTransfer, error)
	ListTransfers(ctx context.Context, companyID string) ([]CashTransfer, error)

	GetCashBook(ctx context.Context, companyID, currency, accountID, fromDate, toDate string) (*CashBook, error)
	GetBalance(ctx context.Context, companyID, accountID string) (float64, error)

	CreatePettyCashFund(ctx context.Context, f *PettyCashFund) error
	GetPettyCashFund(ctx context.Context, id string) (*PettyCashFund, error)
	ListPettyCashFunds(ctx context.Context, companyID string) ([]PettyCashFund, error)
	UpdatePettyCashFund(ctx context.Context, f *PettyCashFund) error

	CreateInventory(ctx context.Context, inv *CashInventory) error
	GetInventory(ctx context.Context, id string) (*CashInventory, error)
	ListInventories(ctx context.Context, companyID string) ([]CashInventory, error)
	UpdateInventory(ctx context.Context, inv *CashInventory) error

	CreateAdvance(ctx context.Context, a *AdvanceRequest) error
	GetAdvance(ctx context.Context, id string) (*AdvanceRequest, error)
	ListAdvances(ctx context.Context, companyID string) ([]AdvanceRequest, error)
	UpdateAdvance(ctx context.Context, a *AdvanceRequest) error
	ListAdvancesByStatus(ctx context.Context, companyID string, status AdvanceStatus) ([]AdvanceRequest, error)
	CreateAdvanceSettlement(ctx context.Context, s *AdvanceSettlement) error
}

// ─── Purchase Repository ───────────────────────────────────────────

type SupplierRepository interface {
	CreateSupplier(ctx context.Context, s *Supplier) error
	GetSupplier(ctx context.Context, id string) (*Supplier, error)
	GetSupplierByCode(ctx context.Context, companyID, code string) (*Supplier, error)
	ListSuppliers(ctx context.Context, filter PurchaseOrderFilter) ([]Supplier, int, error)
	ListSuppliersByIDs(ctx context.Context, ids []string) ([]Supplier, error)
	UpdateSupplier(ctx context.Context, s *Supplier) error
	DeleteSupplier(ctx context.Context, id string) error
}

type PurchaseOrderRepository interface {
	CreatePO(ctx context.Context, po *PurchaseOrder) error
	GetPO(ctx context.Context, id string) (*PurchaseOrder, error)
	GetPOByNumber(ctx context.Context, companyID, poNumber string) (*PurchaseOrder, error)
	ListPOs(ctx context.Context, filter PurchaseOrderFilter) ([]PurchaseOrder, int, error)
	UpdatePO(ctx context.Context, po *PurchaseOrder) error
	UpdatePOStatus(ctx context.Context, id string, status POStatus) error
	ApprovePO(ctx context.Context, id string, approvedBy string, approvedAt time.Time) error
	CancelPO(ctx context.Context, id string, cancelReason string) error
	GetPOLines(ctx context.Context, poID string) ([]POItem, error)
	CreatePOLines(ctx context.Context, items []POItem) error
	UpdatePOLines(ctx context.Context, items []POItem) error
	NextPONumber(ctx context.Context, companyID, yyyymm string) (string, error)
}

type GRNRepository interface {
	CreateGRN(ctx context.Context, g *GRN) error
	GetGRN(ctx context.Context, id string) (*GRN, error)
	GetGRNByNumber(ctx context.Context, companyID, grnNumber string) (*GRN, error)
	ListGRNs(ctx context.Context, filter GRNFilter) ([]GRN, int, error)
	UpdateGRN(ctx context.Context, g *GRN) error
	UpdateGRNStatus(ctx context.Context, id string, status GRNStatus) error
	GetGRNLines(ctx context.Context, grnID string) ([]GRNItem, error)
	CreateGRNLines(ctx context.Context, items []GRNItem) error
	UpdateGRNLines(ctx context.Context, items []GRNItem) error
	NextGRNNumber(ctx context.Context, companyID, yyyymm string) (string, error)
}

type SupplierInvoiceRepository interface {
	CreateInvoice(ctx context.Context, inv *SupplierInvoice) error
	GetInvoice(ctx context.Context, id string) (*SupplierInvoice, error)
	GetInvoiceByNumber(ctx context.Context, companyID, invoiceNumber string) (*SupplierInvoice, error)
	ListInvoices(ctx context.Context, filter SupplierInvoiceFilter) ([]SupplierInvoice, int, error)
	UpdateInvoice(ctx context.Context, inv *SupplierInvoice) error
	UpdateInvoiceStatus(ctx context.Context, id string, status InvoiceStatus) error
	PostInvoice(ctx context.Context, id string, postedAt time.Time) error
	SetInvoiceGLPosted(ctx context.Context, id string, postedAt time.Time) error
	GetInvoiceLines(ctx context.Context, invoiceID string) ([]SupplierInvoiceLine, error)
	CreateInvoiceLines(ctx context.Context, items []SupplierInvoiceLine) error
	UpdateInvoiceLines(ctx context.Context, items []SupplierInvoiceLine) error
	ReduceInvoiceBalance(ctx context.Context, id string, amount float64) error
}

type APTransactionRepository interface {
	CreateAPTransaction(ctx context.Context, t *APTransaction) error
	GetAPTransaction(ctx context.Context, id string) (*APTransaction, error)
	ListAPTransactionsBySupplier(ctx context.Context, companyID, supplierID string) ([]APTransaction, error)
	ListAPTransactions(ctx context.Context, companyID string, offset, limit int) ([]APTransaction, int, error)
}

type CostAllocationRepository interface {
	CreateCostAllocation(ctx context.Context, c *CostAllocation) error
	GetCostAllocation(ctx context.Context, id string) (*CostAllocation, error)
	ListCostAllocationsByInvoice(ctx context.Context, invoiceID string) ([]CostAllocation, error)
}

type DoubtfulDebtProvisionRepository interface {
	CreateProvision(ctx context.Context, p *DoubtfulDebtProvision) error
	CreateProvisionLines(ctx context.Context, lines []DoubtfulDebtProvisionLine) error
	GetProvision(ctx context.Context, id string) (*DoubtfulDebtProvision, error)
	GetProvisionLines(ctx context.Context, provisionID string) ([]DoubtfulDebtProvisionLine, error)
	ListProvisions(ctx context.Context, companyID string, limit, offset int) ([]DoubtfulDebtProvision, int, error)
}

type FXRevaluationRepository interface {
	CreateRevaluation(ctx context.Context, r *FXRevaluation) error
	CreateRevaluationLines(ctx context.Context, lines []FXRevaluationLine) error
	GetRevaluation(ctx context.Context, id string) (*FXRevaluation, error)
	GetRevaluationLines(ctx context.Context, revaluationID string) ([]FXRevaluationLine, error)
	ListRevaluations(ctx context.Context, companyID string, limit, offset int) ([]FXRevaluation, int, error)
	UpdateRevaluationStatus(ctx context.Context, id string, status FXRevaluationStatus) error
	SetRevaluationGLPosted(ctx context.Context, id string, postedAt time.Time) error
}

type RequisitionRepository interface {
	CreateRequisition(ctx context.Context, r *PurchaseRequisition) error
	CreateRequisitionLines(ctx context.Context, lines []RequisitionItem) error
	GetRequisition(ctx context.Context, id string) (*PurchaseRequisition, error)
	GetRequisitionLines(ctx context.Context, requisitionID string) ([]RequisitionItem, error)
	ListRequisitions(ctx context.Context, filter RequisitionFilter) ([]PurchaseRequisition, int, error)
	UpdateRequisition(ctx context.Context, r *PurchaseRequisition) error
	UpdateRequisitionStatus(ctx context.Context, id string, status RequisitionStatus, approvedBy string, approvedAt time.Time) error
	RejectRequisition(ctx context.Context, id string, reason string) error
	DeleteRequisition(ctx context.Context, id string) error
}

// ─── Sale Repository ────────────────────────────────────────────

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, c *Customer) error
	GetCustomer(ctx context.Context, id string) (*Customer, error)
	GetCustomerByCode(ctx context.Context, companyID, code string) (*Customer, error)
	ListCustomers(ctx context.Context, companyID string) ([]Customer, error)
	UpdateCustomer(ctx context.Context, c *Customer) error
	DeleteCustomer(ctx context.Context, id string) error
}

type SaleOrderRepository interface {
	CreateSO(ctx context.Context, so *SalesOrder) error
	GetSO(ctx context.Context, id string) (*SalesOrder, error)
	GetSOByNumber(ctx context.Context, companyID, soNumber string) (*SalesOrder, error)
	ListSOs(ctx context.Context, filter SalesOrderFilter) ([]SalesOrder, int, error)
	UpdateSO(ctx context.Context, so *SalesOrder) error
	UpdateSOStatus(ctx context.Context, id string, status SOStatus) error
	ApproveSO(ctx context.Context, id, approvedBy string, approvedAt time.Time) error
	CancelSO(ctx context.Context, id, cancelReason string) error
	GetSOLines(ctx context.Context, soID string) ([]SOLine, error)
	CreateSOLines(ctx context.Context, items []SOLine) error
	UpdateSOLines(ctx context.Context, items []SOLine) error
	NextSONumber(ctx context.Context, companyID, yyyymm string) (string, error)
}

type DeliveryNoteRepository interface {
	CreateDN(ctx context.Context, dn *DeliveryNote) error
	GetDN(ctx context.Context, id string) (*DeliveryNote, error)
	GetDNByNumber(ctx context.Context, companyID, dnNumber string) (*DeliveryNote, error)
	ListDNs(ctx context.Context, filter DeliveryNoteFilter) ([]DeliveryNote, int, error)
	UpdateDN(ctx context.Context, dn *DeliveryNote) error
	UpdateDNStatus(ctx context.Context, id string, status DNStatus) error
	SetDNGLPosted(ctx context.Context, id string, postedAt time.Time) error
	GetDNLines(ctx context.Context, dnID string) ([]DNLine, error)
	CreateDNLines(ctx context.Context, items []DNLine) error
	UpdateDNLines(ctx context.Context, items []DNLine) error
	NextDNNumber(ctx context.Context, companyID, yyyymm string) (string, error)
}

type CustomerInvoiceRepository interface {
	CreateInvoice(ctx context.Context, inv *CustomerInvoice) error
	GetInvoice(ctx context.Context, id string) (*CustomerInvoice, error)
	GetInvoiceByNumber(ctx context.Context, companyID, invoiceNumber string) (*CustomerInvoice, error)
	ListInvoices(ctx context.Context, filter CustomerInvoiceFilter) ([]CustomerInvoice, int, error)
	UpdateInvoice(ctx context.Context, inv *CustomerInvoice) error
	UpdateInvoiceStatus(ctx context.Context, id string, status SaleInvoiceStatus) error
	PostInvoice(ctx context.Context, id string, postedAt time.Time) error
	SetInvoiceGLPosted(ctx context.Context, id string, postedAt time.Time) error
	AllocateToInvoice(ctx context.Context, invoiceID string, amount float64) error
	ReduceInvoiceBalance(ctx context.Context, invoiceID string, amount float64) error
	GetInvoiceLines(ctx context.Context, invoiceID string) ([]InvLine, error)
	CreateInvoiceLines(ctx context.Context, items []InvLine) error
	UpdateInvoiceLines(ctx context.Context, items []InvLine) error
	NextInvNumber(ctx context.Context, companyID, yyyymm string) (string, error)
}

type CustomerReceiptRepository interface {
	CreateReceipt(ctx context.Context, r *CustomerReceipt) error
	GetReceipt(ctx context.Context, id string) (*CustomerReceipt, error)
	GetReceiptByNumber(ctx context.Context, companyID, receiptNumber string) (*CustomerReceipt, error)
	ListReceipts(ctx context.Context, filter ReceiptFilter) ([]CustomerReceipt, int, error)
	UpdateReceipt(ctx context.Context, r *CustomerReceipt) error
	UpdateReceiptStatus(ctx context.Context, id string, status ReceiptStatus) error
	SetReceiptGLPosted(ctx context.Context, id string, postedAt time.Time) error
	CreateReceiptAllocations(ctx context.Context, allocs []RcpAllocation) error
	GetReceiptAllocations(ctx context.Context, receiptID string) ([]RcpAllocation, error)
}

type CreditNoteRepository interface {
	CreateCN(ctx context.Context, cn *CreditNote) error
	GetCN(ctx context.Context, id string) (*CreditNote, error)
	GetCNByNumber(ctx context.Context, companyID, cnNumber string) (*CreditNote, error)
	ListCNs(ctx context.Context, filter CreditNoteFilter) ([]CreditNote, int, error)
	UpdateCN(ctx context.Context, cn *CreditNote) error
	UpdateCNStatus(ctx context.Context, id string, status CNStatus) error
	PostCN(ctx context.Context, id string, postedAt time.Time) error
	SetCNGLPosted(ctx context.Context, id string, postedAt time.Time) error
	GetCNLines(ctx context.Context, cnID string) ([]CNLine, error)
	CreateCNLines(ctx context.Context, items []CNLine) error
}

type ARTransactionRepository interface {
	CreateARTransaction(ctx context.Context, t *ARTransaction) error
	GetARTransaction(ctx context.Context, id string) (*ARTransaction, error)
	ListARTransactions(ctx context.Context, companyID, customerID string) ([]ARTransaction, error)
	ListARTransactionsAll(ctx context.Context, companyID string, offset, limit int) ([]ARTransaction, int, error)
}

type SalesQuotationRepository interface {
	CreateSQ(ctx context.Context, sq *SalesQuotation) error
	GetSQ(ctx context.Context, id string) (*SalesQuotation, error)
	ListSQs(ctx context.Context, companyID string) ([]SalesQuotation, error)
	UpdateSQ(ctx context.Context, sq *SalesQuotation) error
}

// ─── Warehouse Repository ──────────────────────────────────────────

type WarehouseRepository interface {
	CreateWarehouse(ctx context.Context, w *Warehouse) error
	GetWarehouseByID(ctx context.Context, id string) (*Warehouse, error)
	ListWarehouses(ctx context.Context, companyID string) ([]Warehouse, error)
	UpdateWarehouse(ctx context.Context, w *Warehouse) error
	DeleteWarehouse(ctx context.Context, id string) error
	GetWarehouseByCode(ctx context.Context, companyID, code string) (*Warehouse, error)
}

type ItemCategoryRepository interface {
	CreateCategory(ctx context.Context, c *ItemCategory) error
	GetCategoryByID(ctx context.Context, id string) (*ItemCategory, error)
	ListCategories(ctx context.Context, companyID string) ([]ItemCategory, error)
	UpdateCategory(ctx context.Context, c *ItemCategory) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryByCode(ctx context.Context, companyID, code string) (*ItemCategory, error)
}

type ItemRepository interface {
	CreateItem(ctx context.Context, i *Item) error
	GetItemByID(ctx context.Context, id string) (*Item, error)
	ListItems(ctx context.Context, companyID string) ([]Item, error)
	UpdateItem(ctx context.Context, i *Item) error
	DeleteItem(ctx context.Context, id string) error
	GetItemByCode(ctx context.Context, companyID, code string) (*Item, error)
}

type StockBalanceRepository interface {
	CreateStockBalance(ctx context.Context, b *StockBalance) error
	GetStockBalanceByID(ctx context.Context, id string) (*StockBalance, error)
	FindStockBalance(ctx context.Context, companyID, warehouseID, itemID, period string) (*StockBalance, error)
	ListStockBalances(ctx context.Context, companyID, warehouseID string) ([]StockBalance, error)
	UpdateStockBalance(ctx context.Context, b *StockBalance) error
	UpsertStockBalance(ctx context.Context, b *StockBalance) error
}

type InventoryTransactionRepository interface {
	CreateInventoryTransaction(ctx context.Context, t *InventoryTransaction) error
	GetInventoryTransactionByID(ctx context.Context, id string) (*InventoryTransaction, error)
	ListInventoryTransactions(ctx context.Context, companyID, warehouseID, itemID string, offset, limit int) ([]InventoryTransaction, int, error)
}

type StockTransferRepository interface {
	CreateStockTransfer(ctx context.Context, t *StockTransfer) error
	GetStockTransferByID(ctx context.Context, id string) (*StockTransfer, error)
	ListStockTransfers(ctx context.Context, companyID string) ([]StockTransfer, error)
	UpdateStockTransfer(ctx context.Context, t *StockTransfer) error
	UpdateStockTransferStatus(ctx context.Context, id string, status TransferStatus) error
	GetTransferItems(ctx context.Context, transferID string) ([]TransferItem, error)
	CreateTransferItem(ctx context.Context, item *TransferItem) error
}

type StockAdjustmentRepository interface {
	CreateStockAdjustment(ctx context.Context, a *StockAdjustment) error
	GetStockAdjustmentByID(ctx context.Context, id string) (*StockAdjustment, error)
	ListStockAdjustments(ctx context.Context, companyID string) ([]StockAdjustment, error)
	UpdateStockAdjustment(ctx context.Context, a *StockAdjustment) error
	UpdateStockAdjustmentStatus(ctx context.Context, id string, status AdjStatus) error
	GetAdjustmentItems(ctx context.Context, adjustmentID string) ([]AdjItem, error)
	CreateAdjustmentItem(ctx context.Context, item *AdjItem) error
}

type StockTakeRepository interface {
	CreateStockTake(ctx context.Context, t *StockTake) error
	GetStockTakeByID(ctx context.Context, id string) (*StockTake, error)
	ListStockTakes(ctx context.Context, companyID string) ([]StockTake, error)
	UpdateStockTake(ctx context.Context, t *StockTake) error
	UpdateStockTakeStatus(ctx context.Context, id string, status TakeStatus) error
	GetTakeItems(ctx context.Context, takeID string) ([]TakeItem, error)
	CreateTakeItem(ctx context.Context, item *TakeItem) error
}

type InventoryValuationRunRepository interface {
	CreateValuationRun(ctx context.Context, v *InventoryValuationRun) error
	GetValuationRunByID(ctx context.Context, id string) (*InventoryValuationRun, error)
	ListValuationRuns(ctx context.Context, companyID string) ([]InventoryValuationRun, error)
	UpdateValuationRun(ctx context.Context, v *InventoryValuationRun) error
}

// ─── Cost Center Repository ──────────────────────────────────────

type CostCenterRepository interface {
	Create(ctx context.Context, cc *CostCenter) error
	GetByID(ctx context.Context, id string) (*CostCenter, error)
	List(ctx context.Context, companyID string) ([]CostCenter, error)
	Update(ctx context.Context, cc *CostCenter) error
	Delete(ctx context.Context, id string) error
	GetByCode(ctx context.Context, companyID, code string) (*CostCenter, error)
}

// ─── Recurring Entry Repository ─────────────────────────────────────

// ─── Notification Repository ──────────────────────────────────────

type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	GetByID(ctx context.Context, id string) (*Notification, error)
	ListByUser(ctx context.Context, companyID, userID string, limit int) ([]Notification, error)
	UnreadCount(ctx context.Context, companyID, userID string) (int, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, companyID, userID string) error
	Delete(ctx context.Context, id string) error
}

type RecurringEntryRepository interface {
	Create(ctx context.Context, entry *RecurringEntry) error
	GetByID(ctx context.Context, id string) (*RecurringEntry, error)
	List(ctx context.Context, companyID string) ([]RecurringEntry, error)
	Update(ctx context.Context, entry *RecurringEntry) error
	Delete(ctx context.Context, id string) error
	UpdateNextRunDate(ctx context.Context, id, nextDate string) error
	GetDueEntries(ctx context.Context, today string) ([]RecurringEntry, error)
}

// ─── Tools & Equipment (CCDC) Repository ─────────────────────────────

type ToolEquipmentRepository interface {
	Create(ctx context.Context, t *ToolEquipment) error
	GetByID(ctx context.Context, id string) (*ToolEquipment, error)
	List(ctx context.Context, companyID string) ([]ToolEquipment, error)
	Update(ctx context.Context, t *ToolEquipment) error
	Delete(ctx context.Context, id string) error
	GetByCode(ctx context.Context, companyID, code string) (*ToolEquipment, error)
}

type ToolEquipmentCategoryRepository interface {
	Create(ctx context.Context, c *ToolEquipmentCategory) error
	GetByID(ctx context.Context, id string) (*ToolEquipmentCategory, error)
	List(ctx context.Context, companyID string) ([]ToolEquipmentCategory, error)
	Update(ctx context.Context, c *ToolEquipmentCategory) error
	Delete(ctx context.Context, id string) error
}

// ─── Fixed Asset Repository ─────────────────────────────────────────

type FARepository interface {
	CreateCategory(ctx context.Context, c *FixedAssetCategory) error
	GetCategoryByID(ctx context.Context, id string) (*FixedAssetCategory, error)
	GetCategoryByCode(ctx context.Context, companyID, code string) (*FixedAssetCategory, error)
	ListCategories(ctx context.Context, filter FACategoryFilter) ([]FixedAssetCategory, error)
	UpdateCategory(ctx context.Context, c *FixedAssetCategory) error
	DeleteCategory(ctx context.Context, id string) error

	CreateAsset(ctx context.Context, a *FixedAsset) error
	GetAssetByID(ctx context.Context, id string) (*FixedAsset, error)
	GetAssetByCode(ctx context.Context, companyID, code string) (*FixedAsset, error)
	ListAssets(ctx context.Context, filter FAListFilter) ([]FixedAsset, int, error)
	UpdateAsset(ctx context.Context, a *FixedAsset) error
	DeleteAsset(ctx context.Context, id string) error

	CreateDepreciationEntry(ctx context.Context, e *DepreciationEntry) error
	GetDepreciationEntry(ctx context.Context, id string) (*DepreciationEntry, error)
	ListDepreciationByAsset(ctx context.Context, assetID string) ([]DepreciationEntry, error)
	ListDepreciationByPeriod(ctx context.Context, periodID string) ([]DepreciationEntry, error)
	DepreciationExistsForPeriod(ctx context.Context, assetID, periodID string) (bool, error)
	DeleteDepreciationEntry(ctx context.Context, id string) error

	CreateTransaction(ctx context.Context, t *FixedAssetTransaction) error
	ListTransactionsByAsset(ctx context.Context, assetID string) ([]FixedAssetTransaction, error)

	SetAllocations(ctx context.Context, assetID string, allocs []FixedAssetAllocation) error
	GetAllocations(ctx context.Context, assetID string) ([]FixedAssetAllocation, error)

	CreateInventoryPlan(ctx context.Context, p *FixedAssetInventoryPlan) error
	GetInventoryPlan(ctx context.Context, id string) (*FixedAssetInventoryPlan, error)
	ListInventoryPlans(ctx context.Context, companyID string) ([]FixedAssetInventoryPlan, error)
	UpdateInventoryPlan(ctx context.Context, p *FixedAssetInventoryPlan) error

	CreateInventoryResult(ctx context.Context, r *FixedAssetInventoryResult) error
	GetInventoryResultsByPlan(ctx context.Context, planID string) ([]FixedAssetInventoryResult, error)
}

// ─── Budget Repository ────────────────────────────────────────────

type BudgetRepository interface {
	Create(ctx context.Context, b *Budget) error
	GetByID(ctx context.Context, id string) (*Budget, error)
	List(ctx context.Context, companyID string, year int) ([]Budget, error)
	Update(ctx context.Context, b *Budget) error
	Delete(ctx context.Context, id string) error
	Upsert(ctx context.Context, b *Budget) error
}

// ─── Warehouse Keeper Repository ─────────────────────────────────

type WarehouseKeeperRepository interface {
	// Assignment
	CreateAssignment(ctx context.Context, a *WarehouseKeeperAssignment) error
	GetAssignment(ctx context.Context, id string) (*WarehouseKeeperAssignment, error)
	ListAssignments(ctx context.Context, companyID string) ([]WarehouseKeeperAssignment, error)
	GetActiveAssignment(ctx context.Context, companyID, warehouseID string, date time.Time) (*WarehouseKeeperAssignment, error)
	UpdateAssignment(ctx context.Context, a *WarehouseKeeperAssignment) error
	DeleteAssignment(ctx context.Context, id string) error

	// Stock Ledger
	CreateLedgerEntry(ctx context.Context, e *StockLedgerEntry) error
	GetLedgerEntry(ctx context.Context, id string) (*StockLedgerEntry, error)
	ListLedgerEntries(ctx context.Context, filter LedgerFilter) ([]StockLedgerEntry, int, error)
	UnrecordLedgerEntry(ctx context.Context, id string, unrecordedBy string, reason string) error
	GetLedgerBalance(ctx context.Context, companyID, warehouseID, itemID string) (float64, error)

	// Reconciliation
	GetReconciliationReport(ctx context.Context, companyID, warehouseID string, from, to time.Time) ([]KeeperReconciliationItem, error)

	// Stock Card
	GetStockCard(ctx context.Context, companyID, warehouseID, itemID string, period string) (*StockCard, error)

	// Keeper Reports
	GetKeeperInventorySummary(ctx context.Context, companyID, warehouseID string) ([]KeeperInventorySummaryItem, error)
}

// ─── Cost Accounting Repositories ─────────────────────────────────

type CostObjectRepository interface {
	Create(ctx context.Context, co *CostObject) error
	GetByID(ctx context.Context, id string) (*CostObject, error)
	GetByCode(ctx context.Context, companyID, code string) (*CostObject, error)
	List(ctx context.Context, companyID string) ([]CostObject, error)
	Update(ctx context.Context, co *CostObject) error
	Delete(ctx context.Context, id string) error
}

type CostPoolRepository interface {
	Create(ctx context.Context, cp *CostPool) error
	GetByID(ctx context.Context, id string) (*CostPool, error)
	ListByPeriod(ctx context.Context, companyID, periodID string) ([]CostPool, error)
	Update(ctx context.Context, cp *CostPool) error
	Delete(ctx context.Context, id string) error
}

type CostPoolLineRepository interface {
	Create(ctx context.Context, line *CostPoolLine) error
	ListByPool(ctx context.Context, poolID string) ([]CostPoolLine, error)
	DeleteByPool(ctx context.Context, poolID string) error
}

type CostingPeriodRepository interface {
	Create(ctx context.Context, cp *CostingPeriod) error
	GetByID(ctx context.Context, id string) (*CostingPeriod, error)
	GetByYearMonth(ctx context.Context, companyID string, year, month int) (*CostingPeriod, error)
	List(ctx context.Context, companyID string) ([]CostingPeriod, error)
	Update(ctx context.Context, cp *CostingPeriod) error
}

type CostingResultRepository interface {
	Create(ctx context.Context, cr *CostingResult) error
	GetByID(ctx context.Context, id string) (*CostingResult, error)
	ListByPeriod(ctx context.Context, companyID, periodID string) ([]CostingResult, error)
	Update(ctx context.Context, cr *CostingResult) error
	Delete(ctx context.Context, id string) error
}

type CostingResultLineRepository interface {
	Create(ctx context.Context, line *CostingResultLine) error
	ListByResult(ctx context.Context, resultID string) ([]CostingResultLine, error)
	DeleteByResult(ctx context.Context, resultID string) error
}

// CostDataCollector provides cost data from operational modules for costing period collection.
type CostDataCollector interface {
	CollectMaterialCosts(ctx context.Context, companyID, periodID string) ([]CostPoolLineInput, error)
	CollectLaborCosts(ctx context.Context, companyID, periodID string) ([]CostPoolLineInput, error)
	CollectOverheadCosts(ctx context.Context, companyID, periodID string) ([]CostPoolLineInput, error)
}

type CostPoolLineInput struct {
	SourceID     string
	Description  string
	Amount       float64
	CostCenterID string
}
