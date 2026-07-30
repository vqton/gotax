package domain

import "time"

// ─── Cash Receipt ────────────────────────────────────────────────────────────

type CashReceiptGORM struct {
	ID              string     `gorm:"column:id;primaryKey;size:50" json:"id"`
	CompanyID       string     `gorm:"column:company_id;not null;size:50;index" json:"companyId"`
	VoucherNo       string     `gorm:"column:voucher_no;not null;size:20;uniqueIndex:idx_cr_company_voucher" json:"voucherNo"`
	VoucherDate     time.Time  `gorm:"column:voucher_date;not null;type:date" json:"voucherDate"`
	CashAccountID   string     `gorm:"column:cash_account_id;not null;size:20" json:"cashAccountId"`
	CounterpartID   *string    `gorm:"column:counterpart_id;size:50" json:"counterpartId"`
	CounterpartName *string    `gorm:"column:counterpart_name;size:255" json:"counterpartName"`
	CounterpartType string     `gorm:"column:counterpart_type;not null;size:20;default:'OTHER'" json:"counterpartType"`
	Currency        string     `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	ExchangeRate    float64    `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	Amount          float64    `gorm:"column:amount;not null" json:"amount"`
	AmountVND       float64    `gorm:"column:amount_vnd;not null" json:"amountVnd"`
	DebitAccountID  string     `gorm:"column:debit_account_id;not null;size:20" json:"debitAccountId"`
	CreditAccountID string     `gorm:"column:credit_account_id;not null;size:20" json:"creditAccountId"`
	Reason          *string    `gorm:"column:reason;type:text" json:"reason"`
	ReceiptType     string     `gorm:"column:receipt_type;not null;size:30;default:'OTHER'" json:"receiptType"`
	Status          string     `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	ApprovedBy      *string    `gorm:"column:approved_by;size:50" json:"approvedBy"`
	ApprovedAt      *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	PostedBy        *string    `gorm:"column:posted_by;size:50" json:"postedBy"`
	PostedAt        *time.Time `gorm:"column:posted_at" json:"postedAt"`
	GLJournalID     *string    `gorm:"column:gl_journal_id;size:50" json:"glJournalId"`
	CreatedBy       string     `gorm:"column:created_by;not null;size:50" json:"createdBy"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (CashReceiptGORM) TableName() string { return "cash_receipts" }

type CashReceiptLineGORM struct {
	ID            string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	ReceiptID     string    `gorm:"column:receipt_id;not null;size:36;index" json:"receiptId"`
	AccountID     string    `gorm:"column:account_id;size:36;index" json:"accountId"`
	Amount        float64   `gorm:"column:amount;not null" json:"amount"`
	LocalAmount   float64   `gorm:"column:local_amount;default:0" json:"localAmount"`
	Description   *string   `gorm:"column:description;type:text" json:"description"`
	AllocationRef *string   `gorm:"column:allocation_ref;size:100" json:"allocationRef"`
	LineOrder     int       `gorm:"column:line_order;default:0" json:"lineOrder"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (CashReceiptLineGORM) TableName() string { return "cash_receipt_lines" }

// ─── Cash Payment ────────────────────────────────────────────────────────────

type CashPaymentGORM struct {
	ID              string     `gorm:"column:id;primaryKey;size:50" json:"id"`
	CompanyID       string     `gorm:"column:company_id;not null;size:50;index" json:"companyId"`
	VoucherNo       string     `gorm:"column:voucher_no;not null;size:20;uniqueIndex" json:"voucherNo"`
	VoucherDate     time.Time  `gorm:"column:voucher_date;not null;type:date" json:"voucherDate"`
	CashAccountID   string     `gorm:"column:cash_account_id;not null;size:20" json:"cashAccountId"`
	PayeeID         *string    `gorm:"column:payee_id;size:50" json:"payeeId"`
	PayeeName       *string    `gorm:"column:payee_name;size:255" json:"payeeName"`
	PayeeType       string     `gorm:"column:payee_type;not null;size:20;default:'OTHER'" json:"payeeType"`
	Currency        string     `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	ExchangeRate    float64    `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	Amount          float64    `gorm:"column:amount;not null" json:"amount"`
	AmountVND       float64    `gorm:"column:amount_vnd;not null" json:"amountVnd"`
	DebitAccountID  string     `gorm:"column:debit_account_id;not null;size:20" json:"debitAccountId"`
	CreditAccountID string     `gorm:"column:credit_account_id;not null;size:20" json:"creditAccountId"`
	Reason          *string    `gorm:"column:reason;type:text" json:"reason"`
	PaymentType     string     `gorm:"column:payment_type;not null;size:30;default:'OTHER'" json:"paymentType"`
	Status          string     `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	ApprovedBy      *string    `gorm:"column:approved_by;size:50" json:"approvedBy"`
	ApprovedAt      *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	PostedBy        *string    `gorm:"column:posted_by;size:50" json:"postedBy"`
	PostedAt        *time.Time `gorm:"column:posted_at" json:"postedAt"`
	GLJournalID     *string    `gorm:"column:gl_journal_id;size:50" json:"glJournalId"`
	CreatedBy       string     `gorm:"column:created_by;not null;size:50" json:"createdBy"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (CashPaymentGORM) TableName() string { return "cash_payments" }

type CashPaymentLineGORM struct {
	ID            string    `gorm:"column:id;primaryKey;size:36" json:"id"`
	PaymentID     string    `gorm:"column:payment_id;not null;size:36;index" json:"paymentId"`
	AccountID     string    `gorm:"column:account_id;size:36;index" json:"accountId"`
	Amount        float64   `gorm:"column:amount;not null" json:"amount"`
	LocalAmount   float64   `gorm:"column:local_amount;default:0" json:"localAmount"`
	Description   *string   `gorm:"column:description;type:text" json:"description"`
	AllocationRef *string   `gorm:"column:allocation_ref;size:100" json:"allocationRef"`
	LineOrder     int       `gorm:"column:line_order;default:0" json:"lineOrder"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (CashPaymentLineGORM) TableName() string { return "cash_payment_lines" }

// ─── Cash Transfer ───────────────────────────────────────────────────────────

type CashTransferGORM struct {
	ID              string     `gorm:"column:id;primaryKey;size:50" json:"id"`
	CompanyID       string     `gorm:"column:company_id;not null;size:50;index" json:"companyId"`
	TransferDate    time.Time  `gorm:"column:transfer_date;not null;type:date" json:"transferDate"`
	FromAccountID   string     `gorm:"column:from_account_id;not null;size:20" json:"fromAccountId"`
	ToAccountID     string     `gorm:"column:to_account_id;not null;size:20" json:"toAccountId"`
	Amount          float64    `gorm:"column:amount;not null" json:"amount"`
	Currency        string     `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	ExchangeRate    float64    `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	Reason          *string    `gorm:"column:reason;type:text" json:"reason"`
	TransferType    string     `gorm:"column:transfer_type;not null;size:30" json:"transferType"`
	Status          string     `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	SourceVoucherID *string    `gorm:"column:source_voucher_id;size:50" json:"sourceVoucherId"`
	DestVoucherID   *string    `gorm:"column:dest_voucher_id;size:50" json:"destVoucherId"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	PostedAt        *time.Time `gorm:"column:posted_at" json:"postedAt"`
}

func (CashTransferGORM) TableName() string { return "cash_transfers" }

// ─── Petty Cash Fund ────────────────────────────────────────────────────────

type PettyCashFundGORM struct {
	ID             string    `gorm:"column:id;primaryKey;size:50" json:"id"`
	CompanyID      string    `gorm:"column:company_id;not null;size:50;index" json:"companyId"`
	FundCode       string    `gorm:"column:fund_code;not null;size:20" json:"fundCode"`
	FundName       string    `gorm:"column:fund_name;not null;size:255" json:"fundName"`
	Custodian      string    `gorm:"column:custodian_id;not null;size:50" json:"custodianId"`
	FundAmount     float64   `gorm:"column:initial_amount;not null" json:"initialAmount"`
	CurrentBalance float64   `gorm:"column:current_balance;not null" json:"currentBalance"`
	Currency       string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	Status         string    `gorm:"column:status;not null;size:20;default:'ACTIVE';index" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (PettyCashFundGORM) TableName() string { return "petty_cash_funds" }

// ─── Cash Inventory ─────────────────────────────────────────────────────────

type CashInventoryGORM struct {
	ID            string     `gorm:"column:id;primaryKey;size:50" json:"id"`
	CompanyID     string     `gorm:"column:company_id;not null;size:50;index" json:"companyId"`
	InventoryDate time.Time  `gorm:"column:inventory_date;not null;type:date;index" json:"inventoryDate"`
	Cashier       string     `gorm:"column:cash_account_id;not null;size:20" json:"cashAccountId"`
	Currency      string     `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	TotalCash     float64    `gorm:"column:book_balance;not null" json:"bookBalance"`
	CountedAmount float64    `gorm:"column:actual_balance;not null" json:"actualBalance"`
	Difference    float64    `gorm:"column:difference;not null;default:0" json:"difference"`
	DiffType      string     `gorm:"column:difference_type;not null;size:10;default:'none'" json:"diffType"`
	Reason        *string    `gorm:"column:reason;type:text" json:"reason"`
	Status        string     `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	CreatedBy     *string    `gorm:"column:approved_by;size:50" json:"approvedBy"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (CashInventoryGORM) TableName() string { return "cash_inventories" }

// ─── Advance Request / Settlement ───────────────────────────────────────────

type AdvanceRequestGORM struct {
	ID            string     `gorm:"column:id;primaryKey;size:50" json:"id"`
	CompanyID     string     `gorm:"column:company_id;not null;size:50;index" json:"companyId"`
	RequestNo     string     `gorm:"-" json:"requestNo"`
	EmployeeID    string     `gorm:"column:requestor_id;not null;size:50" json:"requestorId"`
	RequestorName string     `gorm:"column:requestor_name;not null;default:''" json:"requestorName"`
	Amount        float64    `gorm:"column:amount;not null" json:"amount"`
	AmountVND     float64    `gorm:"column:amount_vnd;not null;default:0" json:"amountVnd"`
	Currency      string     `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	ExchangeRate  float64    `gorm:"column:exchange_rate;not null;default:1" json:"exchangeRate"`
	Purpose       string     `gorm:"column:purpose;not null;type:text" json:"purpose"`
	Status        string     `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	ApprovedBy    string     `gorm:"column:approved_by;not null;default:''" json:"approvedBy"`
	ApprovedAt    string     `gorm:"column:approved_at;not null;default:''" json:"approvedAt"`
	PaidBy        string     `gorm:"column:paid_by;not null;default:''" json:"paidBy"`
	PaidAt        string     `gorm:"column:paid_at;not null;default:''" json:"paidAt"`
	GLJEID        string     `gorm:"column:gl_journal_id;not null;default:''" json:"glJournalId"`
	CreatedBy     string     `gorm:"column:created_by;not null;size:50" json:"createdBy"`
	CreatedAt     string     `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt     string     `gorm:"column:updated_at;not null" json:"updatedAt"`
	Settlement    *AdvanceSettlementGORM `gorm:"foreignKey:AdvanceID;constraint:OnDelete:SET NULL" json:"settlement,omitempty"`
}

func (AdvanceRequestGORM) TableName() string { return "advance_requests" }

type AdvanceSettlementGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:50" json:"id"`
	AdvanceID       string    `gorm:"column:advance_id;not null;size:50;uniqueIndex" json:"advanceId"`
	CompanyID       string    `gorm:"column:company_id;not null;size:50" json:"companyId"`
	SettlementDate  time.Time `gorm:"-" json:"settlementDate"`
	SettledAmount   float64   `gorm:"-" json:"settledAmount"`
	ActualSpent     float64   `gorm:"column:total_spent;not null;default:0" json:"totalSpent"`
	RefundAmount    float64   `gorm:"column:remaining_amount;not null;default:0" json:"remainingAmount"`
	Currency        string    `gorm:"column:currency;not null;size:3;default:VND" json:"currency"`
	ReceiptRef      string    `gorm:"column:notes;not null;default:''" json:"notes"`
	Status          string    `gorm:"column:status;not null;size:20;default:'DRAFT'" json:"status"`
	SettledBy       string    `gorm:"-" json:"settledBy"`
	CreatedAt       string    `gorm:"column:created_at;not null" json:"createdAt"`
}

func (AdvanceSettlementGORM) TableName() string { return "advance_settlements" }
