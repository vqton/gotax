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
	GetInvoiceLines(ctx context.Context, invoiceID string) ([]SupplierInvoiceLine, error)
	CreateInvoiceLines(ctx context.Context, items []SupplierInvoiceLine) error
	UpdateInvoiceLines(ctx context.Context, items []SupplierInvoiceLine) error
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
