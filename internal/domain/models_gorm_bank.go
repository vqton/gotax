package domain

import "time"

type PaymentOrderGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID    string    `gorm:"column:company_id;not null;size:36;index:idx_po_company,unique" json:"companyId"`
	PayeeName    string    `gorm:"column:payee_name;not null;size:255" json:"payeeName"`
	PayeeTaxCode *string   `gorm:"column:payee_tax_code;size:20" json:"payeeTaxCode"`
	PayeeBank    *string   `gorm:"column:payee_bank;size:100" json:"payeeBank"`
	PayeeAccount *string   `gorm:"column:payee_account;size:30" json:"payeeAccount"`
	Amount       float64   `gorm:"column:amount;not null" json:"amount"`
	Currency     string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	Priority     string    `gorm:"column:priority;not null;size:10;default:'NORMAL'" json:"priority"`
	Purpose      *string   `gorm:"column:purpose;type:text" json:"purpose"`
	Status       string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	DueDate      time.Time `gorm:"column:due_date;not null;type:date" json:"dueDate"`
	PaidAt       *time.Time `gorm:"column:paid_at" json:"paidAt"`
	CreatedBy    string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	BatchID      *string   `gorm:"column:batch_id;size:36;index" json:"batchId"`
}

func (PaymentOrderGORM) TableName() string { return "payment_orders" }

type PaymentOrderBatchGORM struct {
	ID          string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID   string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	BatchNumber string    `gorm:"column:batch_number;not null;size:30;uniqueIndex" json:"batchNumber"`
	TotalAmount float64   `gorm:"column:total_amount;not null" json:"totalAmount"`
	OrderCount  int       `gorm:"column:order_count;not null" json:"orderCount"`
	Status      string    `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	SubmittedAt *time.Time `gorm:"column:submitted_at" json:"submittedAt"`
	SubmittedBy *string   `gorm:"column:submitted_by;size:36" json:"submittedBy"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (PaymentOrderBatchGORM) TableName() string { return "payment_order_batches" }

type LoanAgreementGORM struct {
	ID                string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string    `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	LenderName        string    `gorm:"column:lender_name;not null;size:255" json:"lenderName"`
	LenderTaxCode     *string   `gorm:"column:lender_tax_code;size:20" json:"lenderTaxCode"`
	Principal         float64   `gorm:"column:principal;not null" json:"principal"`
	Currency          string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	InterestRate      float64   `gorm:"column:interest_rate;not null" json:"interestRate"`
	StartDate         time.Time `gorm:"column:start_date;not null;type:date" json:"startDate"`
	MaturityDate      time.Time `gorm:"column:maturity_date;not null;type:date" json:"maturityDate"`
	TermMonths        int       `gorm:"column:term_months;not null" json:"termMonths"`
	RepaymentType     string    `gorm:"column:repayment_type;not null;size:30" json:"repaymentType"`
	Collateral        *string   `gorm:"column:collateral;type:text" json:"collateral"`
	Purpose           *string   `gorm:"column:purpose;type:text" json:"purpose"`
	Status            string    `gorm:"column:status;not null;size:20;default:'ACTIVE';index" json:"status"`
	GLJECreated       bool      `gorm:"column:gl_je_created;default:false" json:"glJeCreated"`
	CreatedBy         string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Disbursements     []LoanDisbursementGORM `gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE" json:"disbursements,omitempty"`
	Repayments        []LoanRepaymentGORM    `gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE" json:"repayments,omitempty"`
}

func (LoanAgreementGORM) TableName() string { return "loan_agreements" }

type LoanDisbursementGORM struct {
	ID           string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	LoanID       string    `gorm:"column:loan_id;not null;size:36;index" json:"loanId"`
	DisburseDate time.Time `gorm:"column:disburse_date;not null;type:date" json:"disburseDate"`
	Amount       float64   `gorm:"column:amount;not null" json:"amount"`
	BankAccount  *string   `gorm:"column:bank_account;size:50" json:"bankAccount"`
	BankRef      *string   `gorm:"column:bank_ref;size:100" json:"bankRef"`
	Note         *string   `gorm:"column:note;type:text" json:"note"`
	CreatedBy    string    `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (LoanDisbursementGORM) TableName() string { return "loan_disbursements" }

type LoanRepaymentGORM struct {
	ID           string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	LoanID       string     `gorm:"column:loan_id;not null;size:36;index" json:"loanId"`
	RepayDate    time.Time  `gorm:"column:repay_date;not null;type:date" json:"repayDate"`
	Principal    float64    `gorm:"column:principal;not null;default:0" json:"principal"`
	Interest     float64    `gorm:"column:interest;not null;default:0" json:"interest"`
	TotalAmount  float64    `gorm:"column:total_amount;not null" json:"totalAmount"`
	BankRef      *string    `gorm:"column:bank_ref;size:100" json:"bankRef"`
	CreatedBy    string     `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (LoanRepaymentGORM) TableName() string { return "loan_repayments" }

type TermDepositGORM struct {
	ID                string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	BankName          string     `gorm:"column:bank_name;not null;size:255" json:"bankName"`
	CertificateNumber *string    `gorm:"column:certificate_number;size:100;uniqueIndex" json:"certificateNumber"`
	Principal         float64    `gorm:"column:principal;not null" json:"principal"`
	Currency          string     `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	InterestRate      float64    `gorm:"column:interest_rate;not null" json:"interestRate"`
	StartDate         time.Time  `gorm:"column:start_date;not null;type:date" json:"startDate"`
	MaturityDate      time.Time  `gorm:"column:maturity_date;not null;type:date" json:"maturityDate"`
	TermDays          int        `gorm:"column:term_days;not null" json:"termDays"`
	AutoRenew         bool       `gorm:"column:auto_renew;default:false" json:"autoRenew"`
	Status            string     `gorm:"column:status;not null;size:20;default:'ACTIVE';index" json:"status"`
	MaturedAt         *time.Time `gorm:"column:matured_at" json:"maturedAt"`
	CreatedBy         string     `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (TermDepositGORM) TableName() string { return "term_deposits" }

// ─── Bank Statements ────────────────────────────────────────────────────

type BankStatementGORM struct {
	ID              string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID       string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	BankAccountID   string     `gorm:"column:bank_account_id;size:36;index" json:"bankAccountId"`
	StatementDate   time.Time  `gorm:"column:statement_date;not null;type:date" json:"statementDate"`
	StatementRef    string     `gorm:"column:statement_ref;size:50" json:"statementRef"`
	FromDate        *time.Time `gorm:"column:from_date;type:date" json:"fromDate"`
	ToDate          *time.Time `gorm:"column:to_date;type:date" json:"toDate"`
	LineCount       int        `gorm:"column:line_count;default:0" json:"lineCount"`
	OpeningBalance  float64    `gorm:"column:opening_balance;not null" json:"openingBalance"`
	ClosingBalance  float64    `gorm:"column:closing_balance;not null" json:"closingBalance"`
	TotalDebits     float64    `gorm:"column:total_debits;not null;default:0" json:"totalDebits"`
	TotalCredits    float64    `gorm:"column:total_credits;not null;default:0" json:"totalCredits"`
	Currency        string     `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	ImportMethod    string     `gorm:"column:import_method;size:20" json:"importMethod"`
	RawFileName     *string    `gorm:"column:raw_file_name;size:255" json:"rawFileName"`
	RawFileHash     *string    `gorm:"column:raw_file_hash;size:64" json:"rawFileHash"`
	ImportedBy      string     `gorm:"column:imported_by;size:36" json:"importedBy"`
	ImportedAt      *time.Time `gorm:"column:imported_at" json:"importedAt"`
	Status          string     `gorm:"column:status;not null;size:20;default:'DRAFT'" json:"status"`
	Notes           *string    `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy       string     `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (BankStatementGORM) TableName() string { return "bank_statements" }

type BankStatementLineGORM struct {
	ID                string      `gorm:"column:id;primaryKey;size:36" json:"id"`
	StatementID       string      `gorm:"column:statement_id;not null;size:36;index" json:"statementId"`
	LineNumber        int         `gorm:"column:line_number;not null" json:"lineNumber"`
	TransactionDate   time.Time   `gorm:"column:transaction_date;not null;type:date" json:"transactionDate"`
	ValueDate         *time.Time  `gorm:"column:value_date;type:date" json:"valueDate"`
	Reference         *string     `gorm:"column:reference;size:100" json:"reference"`
	Description       *string     `gorm:"column:description;type:text" json:"description"`
	DebitAmount       float64     `gorm:"column:debit_amount;not null;default:0" json:"debitAmount"`
	CreditAmount      float64     `gorm:"column:credit_amount;not null;default:0" json:"creditAmount"`
	BalanceAfter      float64     `gorm:"column:balance_after;default:0" json:"balanceAfter"`
	ReferenceNo       *string     `gorm:"column:reference_no;size:100" json:"referenceNo"`
	BankRef           *string     `gorm:"column:bank_ref;size:100" json:"bankRef"`
	Counterparty      *string     `gorm:"column:counterparty;size:255" json:"counterparty"`
	CounterpartyAcc   *string     `gorm:"column:counterparty_acc;size:50" json:"counterpartyAcc"`
	CounterpartyBank  *string     `gorm:"column:counterparty_bank;size:100" json:"counterpartyBank"`
	RawData           *string     `gorm:"column:raw_data;type:text" json:"rawData"`
	MatchStatus       string      `gorm:"column:match_status;size:20;default:'PENDING'" json:"matchStatus"`
	MatchedLineID     *string     `gorm:"column:matched_line_id;size:36" json:"matchedLineId"`
	MatchedAt         *time.Time  `gorm:"column:matched_at" json:"matchedAt"`
	MatchedBy         *string     `gorm:"column:matched_by;size:36" json:"matchedBy"`
	Currency          string      `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	FXRate            float64     `gorm:"column:fx_rate;not null;default:1" json:"fxRate"`
	LocalDebit        float64     `gorm:"column:local_debit;default:0" json:"localDebit"`
	LocalCredit       float64     `gorm:"column:local_credit;default:0" json:"localCredit"`
	Reconciled        bool        `gorm:"column:reconciled;not null;default:false" json:"reconciled"`
	Matched           bool        `gorm:"column:matched;not null;default:false" json:"matched"`
	MatchRef          *string     `gorm:"column:match_ref;size:100" json:"matchRef"`
	Notes             *string     `gorm:"column:notes;type:text" json:"notes"`
	CreatedAt         time.Time   `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time   `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (BankStatementLineGORM) TableName() string { return "bank_statement_lines" }

// ─── Bank Reconciliations ────────────────────────────────────────────────

type BankReconciliationGORM struct {
	ID               string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID        string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	BankAccountID    string     `gorm:"column:bank_account_id;not null;size:36;index" json:"bankAccountId"`
	StatementID      string     `gorm:"column:statement_id;size:36;index" json:"statementId"`
	ReconDate        time.Time  `gorm:"column:recon_date;not null;type:date" json:"reconDate"`
	FromDate         *time.Time `gorm:"column:from_date;type:date" json:"fromDate"`
	ToDate           *time.Time `gorm:"column:to_date;type:date" json:"toDate"`
	OpeningBalance   float64    `gorm:"column:opening_balance;not null;default:0" json:"openingBalance"`
	ClosingBalance   float64    `gorm:"column:closing_balance;not null;default:0" json:"closingBalance"`
	StatementBalance float64    `gorm:"column:statement_balance;not null" json:"statementBalance"`
	BookBalance      float64    `gorm:"column:book_balance;not null" json:"bookBalance"`
	Difference       float64    `gorm:"column:difference;not null" json:"difference"`
	MatchedLines     int        `gorm:"column:matched_lines;default:0" json:"matchedLines"`
	UnmatchedLines   int        `gorm:"column:unmatched_lines;default:0" json:"unmatchedLines"`
	WriteOffAmount   float64    `gorm:"column:write_off_amount;default:0" json:"writeOffAmount"`
	CompletedBy      string     `gorm:"column:completed_by;size:36" json:"completedBy"`
	ReversedAt       *time.Time `gorm:"column:reversed_at" json:"reversedAt"`
	Status           string     `gorm:"column:status;not null;size:20;default:'DRAFT'" json:"status"`
	CompletedAt      *time.Time `gorm:"column:completed_at" json:"completedAt"`
	Notes            *string    `gorm:"column:notes;type:text" json:"notes"`
	CreatedBy        string     `gorm:"column:created_by;size:36" json:"createdBy"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (BankReconciliationGORM) TableName() string { return "bank_reconciliations" }

type BankReconciliationMatchGORM struct {
	ID               string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	ReconciliationID string    `gorm:"column:reconciliation_id;not null;size:36;index" json:"reconciliationId"`
	StatementLineID  string    `gorm:"column:statement_line_id;size:36;index" json:"statementLineId"`
	TransactionType  string    `gorm:"column:transaction_type;size:30" json:"transactionType"`
	TransactionID    string    `gorm:"column:transaction_id;size:36" json:"transactionId"`
	TransactionRef   string    `gorm:"column:transaction_ref;size:100" json:"transactionRef"`
	MatchMethod      string    `gorm:"column:match_method;size:30" json:"matchMethod"`
	Confidence       float64   `gorm:"column:confidence;default:0" json:"confidence"`
	Amount           float64   `gorm:"column:amount;not null" json:"amount"`
	MatchedAt        time.Time `gorm:"column:matched_at;autoCreateTime" json:"matchedAt"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (BankReconciliationMatchGORM) TableName() string { return "bank_reconciliation_matches" }

// ─── Bank Ledger Entries ────────────────────────────────────────────────

type BankLedgerEntryGORM struct {
	ID            string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID     string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	BankAccountID string     `gorm:"column:bank_account_id;not null;size:36;index" json:"bankAccountId"`
	EntryDate     time.Time  `gorm:"column:entry_date;not null;type:date;index" json:"entryDate"`
	Reference     string     `gorm:"column:reference;size:100" json:"reference"`
	Description   string     `gorm:"column:description;type:text" json:"description"`
	DebitAmount   float64    `gorm:"column:debit_amount;not null;default:0" json:"debitAmount"`
	CreditAmount  float64    `gorm:"column:credit_amount;not null;default:0" json:"creditAmount"`
	Balance       float64    `gorm:"column:balance;not null" json:"balance"`
	Currency      string     `gorm:"column:currency;size:3;default:VND" json:"currency"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (BankLedgerEntryGORM) TableName() string { return "bank_ledger_entries" }
