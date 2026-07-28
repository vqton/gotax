package domain

type BankStatementStatus string
const (BankStatementImported BankStatementStatus="IMPORTED"; BankStatementProcessing BankStatementStatus="PROCESSING"; BankStatementReconciled BankStatementStatus="RECONCILED"; BankStatementArchived BankStatementStatus="ARCHIVED")
type MatchStatus string
const (MatchPending MatchStatus="PENDING"; MatchMatched MatchStatus="MATCHED"; MatchUnmatched MatchStatus="UNMATCHED"; MatchWrittenOff MatchStatus="WRITTEN_OFF")

type BankStatement struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BankAccountID string `json:"bank_account_id"`
	StatementDate string `json:"statement_date"`
	FromDate string `json:"from_date"`
	ToDate string `json:"to_date"`
	OpeningBalance float64 `json:"opening_balance"`
	ClosingBalance float64 `json:"closing_balance"`
	TotalDebits float64 `json:"total_debits"`
	TotalCredits float64 `json:"total_credits"`
	Currency string `json:"currency"`
	ImportMethod string `json:"import_method"`
	RawFileName string `json:"raw_file_name,omitempty"`
	RawFileHash string `json:"raw_file_hash,omitempty"`
	LineCount int `json:"line_count"`
	Status BankStatementStatus `json:"status"`
	ImportedBy string `json:"imported_by"`
	ImportedAt string `json:"imported_at"`
	Notes string `json:"notes,omitempty"`
}

type BankStatementLine struct {
	ID string `json:"id"`
	StatementID string `json:"statement_id"`
	TransactionDate string `json:"transaction_date"`
	ValueDate string `json:"value_date,omitempty"`
	Description string `json:"description"`
	DebitAmount float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	BalanceAfter float64 `json:"balance_after,omitempty"`
	ReferenceNo string `json:"reference_no,omitempty"`
	BankRef string `json:"bank_ref,omitempty"`
	Counterparty string `json:"counterparty,omitempty"`
	CounterpartyAcc string `json:"counterparty_acc,omitempty"`
	CounterpartyBank string `json:"counterparty_bank,omitempty"`
	RawData string `json:"raw_data,omitempty"`
	MatchStatus MatchStatus `json:"match_status"`
	MatchedLineID string `json:"matched_line_id,omitempty"`
	MatchedAt string `json:"matched_at,omitempty"`
	MatchedBy string `json:"matched_by,omitempty"`
	CreatedAt string `json:"created_at"`
}

type ReconStatus string
const (ReconInProgress ReconStatus="IN_PROGRESS"; ReconCompleted ReconStatus="COMPLETED"; ReconReversed ReconStatus="REVERSED")

type BankReconciliation struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BankAccountID string `json:"bank_account_id"`
	StatementID string `json:"statement_id,omitempty"`
	FromDate string `json:"from_date"`
	ToDate string `json:"to_date"`
	OpeningBalance float64 `json:"opening_balance"`
	ClosingBalance float64 `json:"closing_balance"`
	StatementBalance float64 `json:"statement_balance"`
	Difference float64 `json:"difference"`
	Status ReconStatus `json:"status"`
	MatchedLines int `json:"matched_lines"`
	UnmatchedLines int `json:"unmatched_lines"`
	WriteOffAmount float64 `json:"write_off_amount"`
	CompletedBy string `json:"completed_by,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
	ReversedAt string `json:"reversed_at,omitempty"`
	Notes string `json:"notes,omitempty"`
	CreatedAt string `json:"created_at"`
}

type BankReconciliationMatch struct {
	ID string `json:"id"`
	ReconciliationID string `json:"reconciliation_id"`
	StatementLineID string `json:"statement_line_id"`
	TransactionType string `json:"transaction_type"`
	TransactionID string `json:"transaction_id"`
	TransactionRef string `json:"transaction_ref,omitempty"`
	MatchMethod string `json:"match_method"`
	Confidence float64 `json:"confidence"`
	CreatedAt string `json:"created_at"`
}

type PaymentOrderType string
const (PaymentOrderSupplier PaymentOrderType="SUPPLIER"; PaymentOrderSalary PaymentOrderType="SALARY"; PaymentOrderTax PaymentOrderType="TAX"; PaymentOrderLoan PaymentOrderType="LOAN"; PaymentOrderInternal PaymentOrderType="INTERNAL"; PaymentOrderOther PaymentOrderType="OTHER")
type PaymentOrderStatus string
const (PODraft PaymentOrderStatus="DRAFT"; POPendingApproval PaymentOrderStatus="PENDING_APPROVAL"; POApproved PaymentOrderStatus="APPROVED"; PORejected PaymentOrderStatus="REJECTED"; POSubmitted PaymentOrderStatus="SUBMITTED"; POConfirmed PaymentOrderStatus="CONFIRMED"; POFailed PaymentOrderStatus="FAILED"; POCancelled PaymentOrderStatus="CANCELLED")

func (s PaymentOrderStatus) ValidTransition(next PaymentOrderStatus) bool {
	switch s {
	case PODraft: return next==POPendingApproval||next==POCancelled
	case POPendingApproval: return next==POApproved||next==PORejected
	case POApproved: return next==POSubmitted||next==PODraft
	case PORejected: return next==PODraft
	case POSubmitted: return next==POConfirmed||next==POFailed
	case POConfirmed, POFailed, POCancelled: return false
	default: return false
	}
}

type PaymentOrder struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	PaymentDate string `json:"payment_date"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency"`
	ExchangeRate float64 `json:"exchange_rate"`
	BeneficiaryName string `json:"beneficiary_name"`
	BeneficiaryAcc string `json:"beneficiary_acc"`
	BeneficiaryBank string `json:"beneficiary_bank"`
	BeneficiaryBranch string `json:"beneficiary_branch,omitempty"`
	BeneficiaryCode string `json:"beneficiary_code,omitempty"`
	FromBankAccID string `json:"from_bank_acc_id"`
	PaymentContent string `json:"payment_content,omitempty"`
	Urgent bool `json:"urgent"`
	PaymentType PaymentOrderType `json:"payment_type"`
	ReferenceOrders []string `json:"reference_orders,omitempty"`
	Status PaymentOrderStatus `json:"status"`
	CreatedBy string `json:"created_by"`
	ApprovedBy string `json:"approved_by,omitempty"`
	ApprovedAt string `json:"approved_at,omitempty"`
	SubmittedAt string `json:"submitted_at,omitempty"`
	BankRef string `json:"bank_ref,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	PrintCount int `json:"print_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (po *PaymentOrder) Validate() error {
	if po.Amount<=0 { return ErrPaymentAmountRequired }
	if po.BeneficiaryName=="" { return ErrBeneficiaryNameRequired }
	if po.BeneficiaryAcc=="" { return ErrBeneficiaryAccRequired }
	if po.BeneficiaryBank=="" { return ErrBeneficiaryBankRequired }
	if po.FromBankAccID=="" { return ErrFromBankAccRequired }
	if po.PaymentDate=="" { return ErrPaymentDateRequired }
	return nil
}

type BatchStatus string
const (BatchDraft BatchStatus="DRAFT"; BatchSubmitted BatchStatus="SUBMITTED"; BatchConfirmed BatchStatus="CONFIRMED"; BatchPartial BatchStatus="PARTIAL"; BatchFailed BatchStatus="FAILED")

type PaymentOrderBatch struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BatchName string `json:"batch_name"`
	BatchDate string `json:"batch_date"`
	TotalAmount float64 `json:"total_amount"`
	Currency string `json:"currency"`
	OrderCount int `json:"order_count"`
	Status BatchStatus `json:"status"`
	CreatedBy string `json:"created_by"`
	SubmittedAt string `json:"submitted_at,omitempty"`
	BankRef string `json:"bank_ref,omitempty"`
	CreatedAt string `json:"created_at"`
}

type LoanType string
const (LoanShortTerm LoanType="SHORT_TERM"; LoanLongTerm LoanType="LONG_TERM"; LoanOverdraft LoanType="OVERDRAFT")
type LoanStatus string
const (LoanActive LoanStatus="ACTIVE"; LoanFullyPaid LoanStatus="FULLY_PAID"; LoanOverdue LoanStatus="OVERDUE"; LoanRestructured LoanStatus="RESTRUCTURED"; LoanWrittenOff LoanStatus="WRITTEN_OFF")
type RepaymentMethod string
const (RepayAnnuity RepaymentMethod="ANNUITY"; RepayStraightLine RepaymentMethod="STRAIGHT_LINE"; RepayBullet RepaymentMethod="BULLET")
type RepaymentFreq string
const (RepayMonthly RepaymentFreq="MONTHLY"; RepayQuarterly RepaymentFreq="QUARTERLY"; RepayAnnually RepaymentFreq="ANNUAL"; RepayMaturity RepaymentFreq="MATURITY")
type InterestMethod string
const (InterestFixed InterestMethod="FIXED"; InterestFloating InterestMethod="FLOATING")

type LoanAgreement struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BankAccountID string `json:"bank_account_id"`
	ContractNo string `json:"contract_no"`
	LoanType LoanType `json:"loan_type"`
	PrincipalAmount float64 `json:"principal_amount"`
	Currency string `json:"currency"`
	InterestRate float64 `json:"interest_rate"`
	InterestMethod InterestMethod `json:"interest_method"`
	BaseRate float64 `json:"base_rate,omitempty"`
	MarginRate float64 `json:"margin_rate,omitempty"`
	DisbursedAmount float64 `json:"disbursed_amount"`
	OutstandingBalance float64 `json:"outstanding_balance"`
	StartDate string `json:"start_date"`
	MaturityDate string `json:"maturity_date"`
	RepaymentMethod RepaymentMethod `json:"repayment_method"`
	RepaymentFreq RepaymentFreq `json:"repayment_freq"`
	Status LoanStatus `json:"status"`
	Notes string `json:"notes,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (l *LoanAgreement) Validate() error {
	if l.ContractNo=="" { return ErrContractNoRequired }
	if l.PrincipalAmount<=0 { return ErrLoanPrincipalRequired }
	if l.InterestRate<=0 { return ErrInterestRateRequired }
	if l.StartDate==""||l.MaturityDate=="" { return ErrLoanDateRequired }
	return nil
}

type LoanDisbursement struct {
	ID string `json:"id"`
	LoanID string `json:"loan_id"`
	DisbursementDate string `json:"disbursement_date"`
	Amount float64 `json:"amount"`
	ToBankAccountID string `json:"to_bank_account_id,omitempty"`
	ReferenceNo string `json:"reference_no,omitempty"`
	Notes string `json:"notes,omitempty"`
	CreatedAt string `json:"created_at"`
}

type RepaymentStatus string
const (RepayScheduled RepaymentStatus="SCHEDULED"; RepayPaid RepaymentStatus="PAID"; RepayPartial RepaymentStatus="PARTIAL"; RepayOverdue RepaymentStatus="OVERDUE"; RepayWaived RepaymentStatus="WAIVED")

type LoanRepayment struct {
	ID string `json:"id"`
	LoanID string `json:"loan_id"`
	RepaymentDate string `json:"repayment_date"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestAmount float64 `json:"interest_amount"`
	FeeAmount float64 `json:"fee_amount"`
	TotalAmount float64 `json:"total_amount"`
	PaymentOrderID string `json:"payment_order_id,omitempty"`
	Status RepaymentStatus `json:"status"`
	Notes string `json:"notes,omitempty"`
	CreatedAt string `json:"created_at"`
}

type DepositStatus string
const (DepositActive DepositStatus="ACTIVE"; DepositMatured DepositStatus="MATURED"; DepositRenewed DepositStatus="RENEWED"; DepositClosed DepositStatus="CLOSED")

type TermDeposit struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BankAccountID string `json:"bank_account_id"`
	DepositNo string `json:"deposit_no"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency"`
	InterestRate float64 `json:"interest_rate"`
	TermDays int `json:"term_days"`
	StartDate string `json:"start_date"`
	MaturityDate string `json:"maturity_date"`
	InterestAtMaturity float64 `json:"interest_at_maturity"`
	AutoRenewal bool `json:"auto_renewal"`
	RenewalTermDays int `json:"renewal_term_days,omitempty"`
	Status DepositStatus `json:"status"`
	Notes string `json:"notes,omitempty"`
	CreatedAt string `json:"created_at"`
	MaturedAt string `json:"matured_at,omitempty"`
}

func (d *TermDeposit) Validate() error {
	if d.Amount<=0 { return ErrDepositAmountRequired }
	if d.InterestRate<=0 { return ErrInterestRateRequired }
	if d.TermDays<7 { return ErrDepositMinTerm }
	if d.StartDate==""||d.MaturityDate=="" { return ErrDepositDateRequired }
	return nil
}

type PaymentOrderFilter struct {
	CompanyID string
	Status PaymentOrderStatus
	PaymentType PaymentOrderType
	FromDate string
	ToDate string
	Offset int
	Limit int
}

type LoanFilter struct {
	CompanyID string
	Status LoanStatus
	LoanType LoanType
}

type BankLedgerEntry struct {
	LineNo int `json:"line_no"`
	TransactionDate string `json:"transaction_date"`
	VoucherNo string `json:"voucher_no,omitempty"`
	Description string `json:"description"`
	DebitAmount float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	RunningBalance float64 `json:"running_balance"`
	RefID string `json:"ref_id"`
	CounterpartAccount string `json:"counterpart_account,omitempty"`
}

type BankLedger struct {
	CompanyID string `json:"company_id"`
	BankAccountID string `json:"bank_account_id"`
	BankAccountNo string `json:"bank_account_no"`
	Currency string `json:"currency"`
	FromDate string `json:"from_date"`
	ToDate string `json:"to_date"`
	OpeningBalance float64 `json:"opening_balance"`
	TotalDebits float64 `json:"total_debits"`
	TotalCredits float64 `json:"total_credits"`
	ClosingBalance float64 `json:"closing_balance"`
	Entries []BankLedgerEntry `json:"entries"`
}

type ReconciliationReport struct {
	ReconciliationID string `json:"reconciliation_id"`
	CompanyID string `json:"company_id"`
	BankAccountID string `json:"bank_account_id"`
	StatementID string `json:"statement_id"`
	FromDate string `json:"from_date"`
	ToDate string `json:"to_date"`
	BookOpening float64 `json:"book_opening"`
	BookClosing float64 `json:"book_closing"`
	StatementOpening float64 `json:"statement_opening"`
	StatementClosing float64 `json:"statement_closing"`
	Difference float64 `json:"difference"`
	Adjustments []ReconAdjustment `json:"adjustments"`
	MatchedItems []MatchedItem `json:"matched_items"`
	UnmatchedItems []UnmatchedItem `json:"unmatched_items"`
}

type ReconAdjustment struct {
	Description string `json:"description"`
	DebitAmount float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
}

type MatchedItem struct {
	StatementDate string `json:"statement_date"`
	StatementDesc string `json:"statement_desc"`
	StatementAmount float64 `json:"statement_amount"`
	TransactionType string `json:"transaction_type"`
	TransactionRef string `json:"transaction_ref"`
	TransactionAmount float64 `json:"transaction_amount"`
	MatchMethod string `json:"match_method"`
}

type UnmatchedItem struct {
	StatementDate string `json:"statement_date"`
	StatementDesc string `json:"statement_desc"`
	DebitAmount float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	Reason string `json:"reason"`
}

type RepaymentScheduleItem struct {
	Period int `json:"period"`
	DueDate string `json:"due_date"`
	OpeningBalance float64 `json:"opening_balance"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestAmount float64 `json:"interest_amount"`
	TotalAmount float64 `json:"total_amount"`
	ClosingBalance float64 `json:"closing_balance"`
	Status RepaymentStatus `json:"status"`
}
