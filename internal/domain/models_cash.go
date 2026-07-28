package domain

import "time"

type CashStatus string
const (CashDraft CashStatus="DRAFT"; CashSubmitted CashStatus="SUBMITTED"; CashApproved CashStatus="APPROVED"; CashRejected CashStatus="REJECTED"; CashPosted CashStatus="POSTED")

func (s CashStatus) ValidTransition(next CashStatus) bool {
	switch s{
	case CashDraft: return next==CashSubmitted||next==CashRejected
	case CashSubmitted: return next==CashApproved||next==CashRejected
	case CashApproved: return next==CashPosted||next==CashDraft
	case CashRejected: return next==CashDraft
	case CashPosted: return false
	default: return false
	}
}

type ReceiptType string
const (ReceiptCustomerPayment ReceiptType="CUSTOMER_PAYMENT"; ReceiptLoanRecovery ReceiptType="LOAN_RECOVERY"; ReceiptBankWithdrawal ReceiptType="BANK_WITHDRAWAL"; ReceiptAdvanceRefund ReceiptType="ADVANCE_REFUND"; ReceiptSales ReceiptType="SALES"; ReceiptOther ReceiptType="OTHER")

type PaymentType string
const (PaymentSupplier PaymentType="SUPPLIER"; PaymentSalary PaymentType="SALARY"; PaymentExpense PaymentType="EXPENSE"; PaymentBankDeposit PaymentType="BANK_DEPOSIT"; PaymentAdvance PaymentType="ADVANCE"; PaymentTax PaymentType="TAX"; PaymentOther PaymentType="OTHER")

type TransferType string
const (TransferBankWithdrawal TransferType="BANK_WITHDRAWAL"; TransferBankDeposit TransferType="BANK_DEPOSIT"; TransferCurrencyConversion TransferType="CURRENCY_CONVERSION")

type CounterpartType string
const (CounterpartCustomer CounterpartType="CUSTOMER"; CounterpartSupplier CounterpartType="SUPPLIER"; CounterpartEmployee CounterpartType="EMPLOYEE"; CounterpartOther CounterpartType="OTHER")

type CashReceipt struct {
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	VoucherNo       string          `json:"voucher_no"`
	VoucherDate     string          `json:"voucher_date"`
	CashAccountID   string          `json:"cash_account_id"`
	CounterpartID   string          `json:"counterpart_id,omitempty"`
	CounterpartName string          `json:"counterpart_name,omitempty"`
	CounterpartType CounterpartType `json:"counterpart_type"`
	Currency        string          `json:"currency"`
	ExchangeRate    float64         `json:"exchange_rate"`
	Amount          float64         `json:"amount"`
	AmountVND       float64         `json:"amount_vnd"`
	DebitAccountID  string          `json:"debit_account_id"`
	CreditAccountID string          `json:"credit_account_id"`
	Reason          string          `json:"reason"`
	ReceiptType     ReceiptType     `json:"receipt_type"`
	Status          CashStatus      `json:"status"`
	ApprovedBy      string          `json:"approved_by,omitempty"`
	ApprovedAt      string          `json:"approved_at,omitempty"`
	PostedBy        string          `json:"posted_by,omitempty"`
	PostedAt        string          `json:"posted_at,omitempty"`
	GLJournalID     string          `json:"gl_journal_id,omitempty"`
	CreatedBy       string          `json:"created_by"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (r *CashReceipt) Validate() error {
	if r.Amount<=0 { return ErrCashAmountRequired }
	if r.Currency=="" { r.Currency="VND" }
	if r.Currency!="VND"&&r.ExchangeRate<=0 { return ErrExchangeRateRequired }
	if r.VoucherDate=="" { return ErrJournalEntryInvalidDate }
	if _, err := time.Parse("2006-01-02", r.VoucherDate); err != nil { return ErrJournalEntryInvalidDate }
	if r.Reason=="" { return ErrJournalEntryNoDescription }
	if r.CashAccountID==""||r.DebitAccountID==""||r.CreditAccountID=="" { return ErrCashAccountInvalid }
	return nil
}

type CashPayment struct {
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	VoucherNo       string          `json:"voucher_no"`
	VoucherDate     string          `json:"voucher_date"`
	CashAccountID   string          `json:"cash_account_id"`
	PayeeID         string          `json:"payee_id,omitempty"`
	PayeeName       string          `json:"payee_name,omitempty"`
	PayeeType       CounterpartType `json:"payee_type"`
	Currency        string          `json:"currency"`
	ExchangeRate    float64         `json:"exchange_rate"`
	Amount          float64         `json:"amount"`
	AmountVND       float64         `json:"amount_vnd"`
	DebitAccountID  string          `json:"debit_account_id"`
	CreditAccountID string          `json:"credit_account_id"`
	Reason          string          `json:"reason"`
	PaymentType     PaymentType     `json:"payment_type"`
	Status          CashStatus      `json:"status"`
	ApprovedBy      string          `json:"approved_by,omitempty"`
	ApprovedAt      string          `json:"approved_at,omitempty"`
	PostedBy        string          `json:"posted_by,omitempty"`
	PostedAt        string          `json:"posted_at,omitempty"`
	GLJournalID     string          `json:"gl_journal_id,omitempty"`
	CreatedBy       string          `json:"created_by"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (p *CashPayment) Validate() error {
	if p.Amount<=0 { return ErrCashAmountRequired }
	if p.Currency=="" { p.Currency="VND" }
	if p.Currency!="VND"&&p.ExchangeRate<=0 { return ErrExchangeRateRequired }
	if p.VoucherDate=="" { return ErrJournalEntryInvalidDate }
	if _, err := time.Parse("2006-01-02", p.VoucherDate); err != nil { return ErrJournalEntryInvalidDate }
	if p.Reason=="" { return ErrJournalEntryNoDescription }
	if p.CashAccountID==""||p.DebitAccountID==""||p.CreditAccountID=="" { return ErrCashAccountInvalid }
	return nil
}

type CashTransfer struct {
	ID              string       `json:"id"`
	CompanyID       string       `json:"company_id"`
	TransferDate    string       `json:"transfer_date"`
	FromAccountID   string       `json:"from_account_id"`
	ToAccountID     string       `json:"to_account_id"`
	Amount          float64      `json:"amount"`
	Currency        string       `json:"currency"`
	ExchangeRate    float64      `json:"exchange_rate"`
	Reason          string       `json:"reason"`
	TransferType    TransferType `json:"transfer_type"`
	Status          CashStatus   `json:"status"`
	SourceVoucherID string       `json:"source_voucher_id,omitempty"`
	DestVoucherID   string       `json:"dest_voucher_id,omitempty"`
	CreatedAt       string       `json:"created_at"`
	PostedAt        string       `json:"posted_at,omitempty"`
}

func (t *CashTransfer) Validate() error {
	if t.Amount <= 0 { return ErrCashAmountRequired }
	if t.FromAccountID == "" || t.ToAccountID == "" { return ErrCashAccountInvalid }
	if t.Reason == "" { return ErrJournalEntryNoDescription }
	if t.Currency == "" { t.Currency = "VND" }
	if t.Currency != "VND" && t.ExchangeRate <= 0 { return ErrExchangeRateRequired }
	return nil
}

type CashBookEntry struct {
	LineNo        int     `json:"line_no"`
	VoucherDate   string  `json:"voucher_date"`
	VoucherType   string  `json:"voucher_type"`
	VoucherNo     string  `json:"voucher_no"`
	Description   string  `json:"description"`
	ReceiptAmount float64 `json:"receipt_amount"`
	PaymentAmount float64 `json:"payment_amount"`
	RunningBalance float64 `json:"running_balance"`
	RefID              string  `json:"ref_id"`
	CounterpartAccount string  `json:"counterpart_account,omitempty"`
}

type CashBook struct {
	CompanyID      string          `json:"company_id"`
	Currency       string          `json:"currency"`
	AccountID      string          `json:"account_id"`
	FromDate       string          `json:"from_date"`
	ToDate         string          `json:"to_date"`
	OpeningBalance float64         `json:"opening_balance"`
	TotalReceipts  float64         `json:"total_receipts"`
	TotalPayments  float64         `json:"total_payments"`
	ClosingBalance float64         `json:"closing_balance"`
	Entries        []CashBookEntry `json:"entries"`
}

type PettyCashStatus string
const (PettyCashActive PettyCashStatus="ACTIVE"; PettyCashFrozen PettyCashStatus="FROZEN"; PettyCashClosed PettyCashStatus="CLOSED")

type PettyCashFund struct {
	ID             string          `json:"id"`
	CompanyID      string          `json:"company_id"`
	FundCode       string          `json:"fund_code"`
	FundName       string          `json:"fund_name"`
	CustodianID    string          `json:"custodian_id"`
	InitialAmount  float64         `json:"initial_amount"`
	CurrentBalance float64         `json:"current_balance"`
	Currency       string          `json:"currency"`
	Status         PettyCashStatus `json:"status"`
	CreatedAt      string          `json:"created_at"`
}

type CashInventoryStatus string
const (CashInventoryDraft CashInventoryStatus="DRAFT"; CashInventoryCompleted CashInventoryStatus="COMPLETED")

type DenominationDetail struct {
	Denomination float64 `json:"denomination"`
	Count        int     `json:"count"`
	Subtotal     float64 `json:"subtotal"`
}

type CashInventory struct {
	ID              string               `json:"id"`
	CompanyID       string               `json:"company_id"`
	InventoryDate   string               `json:"inventory_date"`
	CashAccountID   string               `json:"cash_account_id"`
	Currency        string               `json:"currency"`
	BookBalance     float64              `json:"book_balance"`
	ActualBalance   float64              `json:"actual_balance"`
	Difference      float64              `json:"difference"`
	DifferenceType  string               `json:"difference_type"`
	Denominations   []DenominationDetail `json:"denominations,omitempty"`
	Reason          string               `json:"reason,omitempty"`
	Status          CashInventoryStatus  `json:"status"`
	ApprovedBy      string               `json:"approved_by,omitempty"`
	CreatedAt       string               `json:"created_at"`
}

type CashReceiptFilter struct {
	CompanyID  string      `json:"company_id"`
	ReceiptType ReceiptType `json:"receipt_type,omitempty"`
	Currency   string      `json:"currency,omitempty"`
	Status     CashStatus  `json:"status,omitempty"`
	FromDate   string      `json:"from_date,omitempty"`
	ToDate     string      `json:"to_date,omitempty"`
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
}

type CashPaymentFilter struct {
	CompanyID   string      `json:"company_id"`
	PaymentType PaymentType `json:"payment_type,omitempty"`
	Currency    string      `json:"currency,omitempty"`
	Status      CashStatus  `json:"status,omitempty"`
	FromDate    string      `json:"from_date,omitempty"`
	ToDate      string      `json:"to_date,omitempty"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
}

type AdvanceStatus string
const (AdvanceDraft AdvanceStatus="DRAFT"; AdvanceSubmitted AdvanceStatus="SUBMITTED"; AdvanceApproved AdvanceStatus="APPROVED"; AdvanceRejected AdvanceStatus="REJECTED"; AdvancePaid AdvanceStatus="PAID"; AdvanceSettled AdvanceStatus="SETTLED")

func (s AdvanceStatus) ValidTransition(next AdvanceStatus) bool {
	switch s{
	case AdvanceDraft: return next==AdvanceSubmitted||next==AdvanceApproved||next==AdvanceRejected
	case AdvanceSubmitted: return next==AdvanceApproved||next==AdvanceRejected
	case AdvanceApproved: return next==AdvancePaid||next==AdvanceDraft
	case AdvanceRejected: return next==AdvanceDraft
	case AdvancePaid: return next==AdvanceSettled
	case AdvanceSettled: return false
	default: return false
	}
}

type AdvanceRequest struct {
	ID           string        `json:"id"`
	CompanyID    string        `json:"company_id"`
	RequestorID  string        `json:"requestor_id"`
	RequestorName string       `json:"requestor_name,omitempty"`
	Amount       float64       `json:"amount"`
	AmountVND    float64       `json:"amount_vnd"`
	Currency     string        `json:"currency"`
	ExchangeRate float64       `json:"exchange_rate"`
	Purpose      string        `json:"purpose"`
	Status       AdvanceStatus `json:"status"`
	ApprovedBy   string        `json:"approved_by,omitempty"`
	ApprovedAt   string        `json:"approved_at,omitempty"`
	PaidBy       string        `json:"paid_by,omitempty"`
	PaidAt       string        `json:"paid_at,omitempty"`
	GLJournalID  string        `json:"gl_journal_id,omitempty"`
	CreatedBy    string        `json:"created_by"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
}

func (a *AdvanceRequest) Validate() error {
	if a.Amount <= 0 { return ErrCashAmountRequired }
	if a.Purpose == "" { return ErrJournalEntryNoDescription }
	if a.Currency == "" { a.Currency = "VND" }
	if a.Currency != "VND" && a.ExchangeRate <= 0 { return ErrExchangeRateRequired }
	return nil
}

type AdvanceSettlement struct {
	ID              string  `json:"id"`
	AdvanceID       string  `json:"advance_id"`
	CompanyID       string  `json:"company_id"`
	TotalSpent      float64 `json:"total_spent"`
	RemainingAmount float64 `json:"remaining_amount"`
	Currency        string  `json:"currency"`
	Notes           string  `json:"notes,omitempty"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}
