package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

// ─── Cash ────────────────────────────────────────────────────────────────────

type PGCashRepo struct {
	db *gorm.DB
}

func NewPGCashRepo(db *gorm.DB) *PGCashRepo {
	return &PGCashRepo{db: db}
}

// ─── Conversion helpers ──────────────────────────────────────────────────────

func gormReceiptToDomain(m *domain.CashReceiptGORM) *domain.CashReceipt {
	return &domain.CashReceipt{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		VoucherNo:       m.VoucherNo,
		VoucherDate:     safeTimeStr(m.VoucherDate),
		CashAccountID:   m.CashAccountID,
		CounterpartID:   safeStr(m.CounterpartID),
		CounterpartName: safeStr(m.CounterpartName),
		CounterpartType: domain.CounterpartType(m.CounterpartType),
		Currency:        m.Currency,
		ExchangeRate:    m.ExchangeRate,
		Amount:          m.Amount,
		AmountVND:       m.AmountVND,
		DebitAccountID:  m.DebitAccountID,
		CreditAccountID: m.CreditAccountID,
		Reason:          safeStr(m.Reason),
		ReceiptType:     domain.ReceiptType(m.ReceiptType),
		Status:          domain.CashStatus(m.Status),
		ApprovedBy:      safeStr(m.ApprovedBy),
		ApprovedAt:      safeTimePtrStr(m.ApprovedAt),
		PostedBy:        safeStr(m.PostedBy),
		PostedAt:        safeTimePtrStr(m.PostedAt),
		GLJournalID:     safeStr(m.GLJournalID),
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       m.UpdatedAt.Format(time.RFC3339),
	}
}

func domainToGormReceipt(cr *domain.CashReceipt) *domain.CashReceiptGORM {
	return &domain.CashReceiptGORM{
		ID:              cr.ID,
		CompanyID:       cr.CompanyID,
		VoucherNo:       cr.VoucherNo,
		VoucherDate:     parseDate(cr.VoucherDate),
		CashAccountID:   cr.CashAccountID,
		CounterpartID:   nullStrG(cr.CounterpartID),
		CounterpartName: nullStrG(cr.CounterpartName),
		CounterpartType: string(cr.CounterpartType),
		Currency:        cr.Currency,
		ExchangeRate:    cr.ExchangeRate,
		Amount:          cr.Amount,
		AmountVND:       cr.AmountVND,
		DebitAccountID:  cr.DebitAccountID,
		CreditAccountID: cr.CreditAccountID,
		Reason:          nullStrG(cr.Reason),
		ReceiptType:     string(cr.ReceiptType),
		Status:          string(cr.Status),
		ApprovedBy:      nullStrG(cr.ApprovedBy),
		ApprovedAt:      timePtr(parseDate(cr.ApprovedAt)),
		PostedBy:        nullStrG(cr.PostedBy),
		PostedAt:        timePtr(parseDate(cr.PostedAt)),
		GLJournalID:     nullStrG(cr.GLJournalID),
		CreatedBy:       cr.CreatedBy,
	}
}

func gormPaymentToDomain(m *domain.CashPaymentGORM) *domain.CashPayment {
	return &domain.CashPayment{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		VoucherNo:       m.VoucherNo,
		VoucherDate:     safeTimeStr(m.VoucherDate),
		CashAccountID:   m.CashAccountID,
		PayeeID:         safeStr(m.PayeeID),
		PayeeName:       safeStr(m.PayeeName),
		PayeeType:       domain.CounterpartType(m.PayeeType),
		Currency:        m.Currency,
		ExchangeRate:    m.ExchangeRate,
		Amount:          m.Amount,
		AmountVND:       m.AmountVND,
		DebitAccountID:  m.DebitAccountID,
		CreditAccountID: m.CreditAccountID,
		Reason:          safeStr(m.Reason),
		PaymentType:     domain.PaymentType(m.PaymentType),
		Status:          domain.CashStatus(m.Status),
		ApprovedBy:      safeStr(m.ApprovedBy),
		ApprovedAt:      safeTimePtrStr(m.ApprovedAt),
		PostedBy:        safeStr(m.PostedBy),
		PostedAt:        safeTimePtrStr(m.PostedAt),
		GLJournalID:     safeStr(m.GLJournalID),
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       m.UpdatedAt.Format(time.RFC3339),
	}
}

func domainToGormPayment(p *domain.CashPayment) *domain.CashPaymentGORM {
	return &domain.CashPaymentGORM{
		ID:              p.ID,
		CompanyID:       p.CompanyID,
		VoucherNo:       p.VoucherNo,
		VoucherDate:     parseDate(p.VoucherDate),
		CashAccountID:   p.CashAccountID,
		PayeeID:         nullStrG(p.PayeeID),
		PayeeName:       nullStrG(p.PayeeName),
		PayeeType:       string(p.PayeeType),
		Currency:        p.Currency,
		ExchangeRate:    p.ExchangeRate,
		Amount:          p.Amount,
		AmountVND:       p.AmountVND,
		DebitAccountID:  p.DebitAccountID,
		CreditAccountID: p.CreditAccountID,
		Reason:          nullStrG(p.Reason),
		PaymentType:     string(p.PaymentType),
		Status:          string(p.Status),
		ApprovedBy:      nullStrG(p.ApprovedBy),
		ApprovedAt:      timePtr(parseDate(p.ApprovedAt)),
		PostedBy:        nullStrG(p.PostedBy),
		PostedAt:        timePtr(parseDate(p.PostedAt)),
		GLJournalID:     nullStrG(p.GLJournalID),
		CreatedBy:       p.CreatedBy,
	}
}

func gormTransferToDomain(m *domain.CashTransferGORM) *domain.CashTransfer {
	return &domain.CashTransfer{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		TransferDate:    safeTimeStr(m.TransferDate),
		FromAccountID:   m.FromAccountID,
		ToAccountID:     m.ToAccountID,
		Amount:          m.Amount,
		Currency:        m.Currency,
		ExchangeRate:    m.ExchangeRate,
		Reason:          safeStr(m.Reason),
		TransferType:    domain.TransferType(m.TransferType),
		Status:          domain.CashStatus(m.Status),
		SourceVoucherID: safeStr(m.SourceVoucherID),
		DestVoucherID:   safeStr(m.DestVoucherID),
		CreatedAt:       m.CreatedAt.Format(time.RFC3339),
		PostedAt:        safeTimePtrStr(m.PostedAt),
	}
}

func domainToGormTransfer(t *domain.CashTransfer) *domain.CashTransferGORM {
	return &domain.CashTransferGORM{
		ID:              t.ID,
		CompanyID:       t.CompanyID,
		TransferDate:    parseDate(t.TransferDate),
		FromAccountID:   t.FromAccountID,
		ToAccountID:     t.ToAccountID,
		Amount:          t.Amount,
		Currency:        t.Currency,
		ExchangeRate:    t.ExchangeRate,
		Reason:          nullStrG(t.Reason),
		TransferType:    string(t.TransferType),
		Status:          string(t.Status),
		SourceVoucherID: nullStrG(t.SourceVoucherID),
		DestVoucherID:   nullStrG(t.DestVoucherID),
		PostedAt:        timePtr(parseDate(t.PostedAt)),
	}
}

func gormPettyCashToDomain(m *domain.PettyCashFundGORM) *domain.PettyCashFund {
	return &domain.PettyCashFund{
		ID:             m.ID,
		CompanyID:      m.CompanyID,
		FundCode:       m.Custodian,
		FundName:       "",
		CustodianID:    m.Custodian,
		InitialAmount:  m.FundAmount,
		CurrentBalance: m.FundAmount,
		Currency:       m.Currency,
		Status:         domain.PettyCashStatus(m.Status),
		CreatedAt:      m.CreatedAt.Format(time.RFC3339),
	}
}

func domainToGormPettyCash(f *domain.PettyCashFund) *domain.PettyCashFundGORM {
	return &domain.PettyCashFundGORM{
		ID:        f.ID,
		CompanyID: f.CompanyID,
		Custodian: f.CustodianID,
		FundAmount: f.InitialAmount,
		Currency:  f.Currency,
		Status:    string(f.Status),
	}
}

func gormInventoryToDomain(m *domain.CashInventoryGORM) *domain.CashInventory {
	return &domain.CashInventory{
		ID:            m.ID,
		CompanyID:     m.CompanyID,
		InventoryDate: m.InventoryDate.Format("2006-01-02"),
		CashAccountID: m.Cashier,
		BookBalance:   m.TotalCash,
		ActualBalance: m.CountedAmount,
		Difference:    m.Difference,
		Status:        domain.CashInventoryStatus(m.Status),
		ApprovedBy:    safeStr(m.CreatedBy),
		CreatedAt:     m.CreatedAt.Format(time.RFC3339),
	}
}

func domainToGormInventory(inv *domain.CashInventory) *domain.CashInventoryGORM {
	invDate, _ := time.Parse("2006-01-02", inv.InventoryDate)
	return &domain.CashInventoryGORM{
		ID:            inv.ID,
		CompanyID:     inv.CompanyID,
		InventoryDate: invDate,
		Cashier:       inv.CashAccountID,
		TotalCash:     inv.BookBalance,
		CountedAmount: inv.ActualBalance,
		Difference:    inv.Difference,
		Status:        string(inv.Status),
		CreatedBy:     nullStrG(inv.ApprovedBy),
	}
}

func gormAdvanceToDomain(m *domain.AdvanceRequestGORM) *domain.AdvanceRequest {
	return &domain.AdvanceRequest{
		ID:            m.ID,
		CompanyID:     m.CompanyID,
		RequestorID:   m.EmployeeID,
		RequestorName: m.Purpose,
		Amount:        m.Amount,
		AmountVND:     m.Amount,
		Currency:      m.Currency,
		Purpose:       m.Purpose,
		Status:        domain.AdvanceStatus(m.Status),
		ApprovedBy:    m.ApprovedBy,
		ApprovedAt:    m.ApprovedAt,
		GLJournalID:   m.GLJEID,
		CreatedBy:     m.CreatedBy,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func domainToGormAdvance(a *domain.AdvanceRequest) *domain.AdvanceRequestGORM {
	return &domain.AdvanceRequestGORM{
		ID:          a.ID,
		CompanyID:   a.CompanyID,
		RequestNo:   a.ID,
		EmployeeID:  a.RequestorID,
		Amount:      a.Amount,
		Currency:    a.Currency,
		Purpose:     a.Purpose,
		Status:      string(a.Status),
		ApprovedBy:  a.ApprovedBy,
		ApprovedAt:  a.ApprovedAt,
		GLJEID:      a.GLJournalID,
		CreatedBy:   a.CreatedBy,
	}
}

func gormAdvanceSettlementToDomain(m *domain.AdvanceSettlementGORM) *domain.AdvanceSettlement {
	return &domain.AdvanceSettlement{
		ID:              m.ID,
		AdvanceID:       m.AdvanceID,
		CompanyID:       "",
		TotalSpent:      m.ActualSpent,
		RemainingAmount: m.RefundAmount,
		Currency:        "",
		Notes:           m.ReceiptRef,
		Status:          "COMPLETED",
		CreatedAt:       m.CreatedAt,
	}
}

func domainToGormAdvanceSettlement(s *domain.AdvanceSettlement) *domain.AdvanceSettlementGORM {
	return &domain.AdvanceSettlementGORM{
		ID:            s.ID,
		AdvanceID:     s.AdvanceID,
		SettlementDate: time.Now(),
		SettledAmount: s.TotalSpent,
		ActualSpent:   s.TotalSpent,
		RefundAmount:  s.RemainingAmount,
		ReceiptRef:    s.Notes,
		SettledBy:     "",
	}
}

// ─── Cash Receipt ────────────────────────────────────────────────────────────

func (r *PGCashRepo) CreateReceipt(ctx context.Context, cr *domain.CashReceipt) error {
	m := domainToGormReceipt(cr)
	if cr.ID == "" {
		m.ID = fmt.Sprintf("CR-%d", time.Now().UnixNano())
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	cr.ID = m.ID
	return nil
}

func (r *PGCashRepo) GetReceipt(ctx context.Context, id string) (*domain.CashReceipt, error) {
	var m domain.CashReceiptGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCashReceiptNotFound
	}
	return gormReceiptToDomain(&m), nil
}

func (r *PGCashRepo) ListReceipts(ctx context.Context, filter domain.CashReceiptFilter) ([]domain.CashReceipt, int, error) {
	q := r.db.WithContext(ctx).Model(&domain.CashReceiptGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.ReceiptType != "" {
		q = q.Where("receipt_type = ?", string(filter.ReceiptType))
	}
	if filter.Currency != "" {
		q = q.Where("currency = ?", filter.Currency)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		q = q.Where("voucher_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("voucher_date <= ?", filter.ToDate)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	var models []domain.CashReceiptGORM
	if err := q.Order("voucher_date DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CashReceipt, len(models))
	for i := range models {
		out[i] = *gormReceiptToDomain(&models[i])
	}
	return out, int(total), nil
}

func (r *PGCashRepo) UpdateReceipt(ctx context.Context, cr *domain.CashReceipt) error {
	m := domainToGormReceipt(cr)
	return r.db.WithContext(ctx).Model(&domain.CashReceiptGORM{}).Where("id = ?", cr.ID).Updates(map[string]interface{}{
		"voucher_no":       m.VoucherNo,
		"voucher_date":     m.VoucherDate,
		"cash_account_id":  m.CashAccountID,
		"counterpart_id":   m.CounterpartID,
		"counterpart_name": m.CounterpartName,
		"counterpart_type": m.CounterpartType,
		"currency":         m.Currency,
		"exchange_rate":    m.ExchangeRate,
		"amount":           m.Amount,
		"amount_vnd":       m.AmountVND,
		"debit_account_id": m.DebitAccountID,
		"credit_account_id": m.CreditAccountID,
		"reason":           m.Reason,
		"receipt_type":     m.ReceiptType,
		"status":           m.Status,
		"approved_by":      m.ApprovedBy,
		"approved_at":      m.ApprovedAt,
		"posted_by":        m.PostedBy,
		"posted_at":        m.PostedAt,
		"gl_journal_id":    m.GLJournalID,
		"updated_at":       time.Now(),
	}).Error
}

func (r *PGCashRepo) DeleteReceipt(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.CashReceiptGORM{}).Error
}

func (r *PGCashRepo) LastReceiptNo(ctx context.Context, companyID, year string) (string, error) {
	var last string
	err := r.db.WithContext(ctx).Raw(
		`SELECT voucher_no FROM cash_receipts
		 WHERE company_id = ? AND voucher_date LIKE ? || '%'
		 ORDER BY voucher_no DESC LIMIT 1`, companyID, year).Scan(&last).Error
	if err != nil {
		return "", nil
	}
	return last, nil
}

// ─── Cash Payment ────────────────────────────────────────────────────────────

func (r *PGCashRepo) CreatePayment(ctx context.Context, p *domain.CashPayment) error {
	m := domainToGormPayment(p)
	if p.ID == "" {
		m.ID = fmt.Sprintf("CP-%d", time.Now().UnixNano())
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	p.ID = m.ID
	return nil
}

func (r *PGCashRepo) GetPayment(ctx context.Context, id string) (*domain.CashPayment, error) {
	var m domain.CashPaymentGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCashPaymentNotFound
	}
	return gormPaymentToDomain(&m), nil
}

func (r *PGCashRepo) ListPayments(ctx context.Context, filter domain.CashPaymentFilter) ([]domain.CashPayment, int, error) {
	q := r.db.WithContext(ctx).Model(&domain.CashPaymentGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.PaymentType != "" {
		q = q.Where("payment_type = ?", string(filter.PaymentType))
	}
	if filter.Currency != "" {
		q = q.Where("currency = ?", filter.Currency)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.FromDate != "" {
		q = q.Where("voucher_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("voucher_date <= ?", filter.ToDate)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	var models []domain.CashPaymentGORM
	if err := q.Order("voucher_date DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CashPayment, len(models))
	for i := range models {
		out[i] = *gormPaymentToDomain(&models[i])
	}
	return out, int(total), nil
}

func (r *PGCashRepo) UpdatePayment(ctx context.Context, p *domain.CashPayment) error {
	m := domainToGormPayment(p)
	return r.db.WithContext(ctx).Model(&domain.CashPaymentGORM{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
		"voucher_no":       m.VoucherNo,
		"voucher_date":     m.VoucherDate,
		"cash_account_id":  m.CashAccountID,
		"payee_id":         m.PayeeID,
		"payee_name":       m.PayeeName,
		"payee_type":       m.PayeeType,
		"currency":         m.Currency,
		"exchange_rate":    m.ExchangeRate,
		"amount":           m.Amount,
		"amount_vnd":       m.AmountVND,
		"debit_account_id": m.DebitAccountID,
		"credit_account_id": m.CreditAccountID,
		"reason":           m.Reason,
		"payment_type":     m.PaymentType,
		"status":           m.Status,
		"approved_by":      m.ApprovedBy,
		"approved_at":      m.ApprovedAt,
		"posted_by":        m.PostedBy,
		"posted_at":        m.PostedAt,
		"gl_journal_id":    m.GLJournalID,
		"updated_at":       time.Now(),
	}).Error
}

func (r *PGCashRepo) DeletePayment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.CashPaymentGORM{}).Error
}

func (r *PGCashRepo) LastPaymentNo(ctx context.Context, companyID, year string) (string, error) {
	var last string
	err := r.db.WithContext(ctx).Raw(
		`SELECT voucher_no FROM cash_payments
		 WHERE company_id = ? AND voucher_date LIKE ? || '%'
		 ORDER BY voucher_no DESC LIMIT 1`, companyID, year).Scan(&last).Error
	if err != nil {
		return "", nil
	}
	return last, nil
}

// ─── Cash Transfer ───────────────────────────────────────────────────────────

func (r *PGCashRepo) CreateTransfer(ctx context.Context, t *domain.CashTransfer) error {
	m := domainToGormTransfer(t)
	if t.ID == "" {
		m.ID = fmt.Sprintf("CTF-%d", time.Now().UnixNano())
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	t.ID = m.ID
	return nil
}

func (r *PGCashRepo) GetTransfer(ctx context.Context, id string) (*domain.CashTransfer, error) {
	var m domain.CashTransferGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCashTransferNotFound
	}
	return gormTransferToDomain(&m), nil
}

func (r *PGCashRepo) ListTransfers(ctx context.Context, companyID string) ([]domain.CashTransfer, error) {
	var models []domain.CashTransferGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("transfer_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CashTransfer, len(models))
	for i := range models {
		out[i] = *gormTransferToDomain(&models[i])
	}
	return out, nil
}

// ─── Cash Book / Balance ─────────────────────────────────────────────────────

func (r *PGCashRepo) GetCashBook(ctx context.Context, companyID, currency, accountID, fromDate, toDate string) (*domain.CashBook, error) {
	var opening float64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE((
			SELECT SUM(CASE WHEN rpt.type='RECEIPT' THEN rpt.amt ELSE -rpt.amt END)
			FROM (
				SELECT 'RECEIPT' AS type, amount_vnd AS amt, voucher_date
				FROM cash_receipts
				WHERE company_id = ? AND currency = ? AND cash_account_id = ? AND status = 'POSTED' AND voucher_date < ?
				UNION ALL
				SELECT 'PAYMENT', amount_vnd, voucher_date
				FROM cash_payments
				WHERE company_id = ? AND currency = ? AND cash_account_id = ? AND status = 'POSTED' AND voucher_date < ?
			) rpt
		),0)`, companyID, currency, accountID, fromDate, companyID, currency, accountID, fromDate).Scan(&opening).Error
	if err != nil {
		return nil, fmt.Errorf("opening balance: %w", err)
	}

	type cbRow struct {
		VoucherDate   string
		VoucherType   string
		VoucherNo     string
		Reason        string
		ReceiptAmount float64
		PaymentAmount float64
		ID            string
	}
	var rows []cbRow
	err = r.db.WithContext(ctx).Raw(`
		SELECT voucher_date, 'RECEIPT' AS voucher_type, voucher_no, reason,
		       amount_vnd, 0 AS payment_amount, id
		FROM cash_receipts
		WHERE company_id = ? AND currency = ? AND cash_account_id = ? AND status = 'POSTED'
		  AND voucher_date >= ? AND voucher_date <= ?
		UNION ALL
		SELECT voucher_date, 'PAYMENT', voucher_no, reason,
		       0, amount_vnd, id
		FROM cash_payments
		WHERE company_id = ? AND currency = ? AND cash_account_id = ? AND status = 'POSTED'
		  AND voucher_date >= ? AND voucher_date <= ?
		ORDER BY voucher_date, voucher_type`,
		companyID, currency, accountID, fromDate, toDate, companyID, currency, accountID, fromDate, toDate).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("cash book entries: %w", err)
	}

	cb := &domain.CashBook{
		CompanyID:      companyID,
		Currency:       currency,
		AccountID:      accountID,
		FromDate:       fromDate,
		ToDate:         toDate,
		OpeningBalance: opening,
	}
	running := opening
	for i, r := range rows {
		e := domain.CashBookEntry{
			LineNo:        i + 1,
			VoucherDate:   r.VoucherDate,
			VoucherType:   r.VoucherType,
			VoucherNo:     r.VoucherNo,
			Description:   r.Reason,
			ReceiptAmount: r.ReceiptAmount,
			PaymentAmount: r.PaymentAmount,
			RefID:         r.ID,
		}
		running += r.ReceiptAmount - r.PaymentAmount
		e.RunningBalance = running
		cb.Entries = append(cb.Entries, e)
	}
	for _, e := range cb.Entries {
		cb.TotalReceipts += e.ReceiptAmount
		cb.TotalPayments += e.PaymentAmount
	}
	cb.ClosingBalance = opening + cb.TotalReceipts - cb.TotalPayments
	return cb, nil
}

func (r *PGCashRepo) GetBalance(ctx context.Context, companyID, accountID string) (float64, error) {
	var balance float64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE((
			SELECT SUM(amount_vnd) FROM cash_receipts
			WHERE company_id = ? AND cash_account_id = ? AND status = 'POSTED'
		),0) - COALESCE((
			SELECT SUM(amount_vnd) FROM cash_payments
			WHERE company_id = ? AND cash_account_id = ? AND status = 'POSTED'
		),0)`, companyID, accountID, companyID, accountID).Scan(&balance).Error
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// ─── Petty Cash Fund ─────────────────────────────────────────────────────────

func (r *PGCashRepo) CreatePettyCashFund(ctx context.Context, f *domain.PettyCashFund) error {
	m := domainToGormPettyCash(f)
	if f.ID == "" {
		m.ID = fmt.Sprintf("PCF-%d", time.Now().UnixNano())
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	f.ID = m.ID
	return nil
}

func (r *PGCashRepo) GetPettyCashFund(ctx context.Context, id string) (*domain.PettyCashFund, error) {
	var m domain.PettyCashFundGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrPettyCashFundNotFound
	}
	return gormPettyCashToDomain(&m), nil
}

func (r *PGCashRepo) ListPettyCashFunds(ctx context.Context, companyID string) ([]domain.PettyCashFund, error) {
	var models []domain.PettyCashFundGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.PettyCashFund, len(models))
	for i := range models {
		out[i] = *gormPettyCashToDomain(&models[i])
	}
	return out, nil
}

func (r *PGCashRepo) UpdatePettyCashFund(ctx context.Context, f *domain.PettyCashFund) error {
	return r.db.WithContext(ctx).Model(&domain.PettyCashFundGORM{}).Where("id = ?", f.ID).Updates(map[string]interface{}{
		"custodian":   f.CustodianID,
		"fund_amount": f.InitialAmount,
		"currency":    f.Currency,
		"status":      string(f.Status),
	}).Error
}

// ─── Cash Inventory ──────────────────────────────────────────────────────────

func (r *PGCashRepo) CreateInventory(ctx context.Context, inv *domain.CashInventory) error {
	m := domainToGormInventory(inv)
	if inv.ID == "" {
		m.ID = fmt.Sprintf("CI-%d", time.Now().UnixNano())
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		for i, d := range inv.Denominations {
			if err := tx.Exec(
				`INSERT INTO cash_inventory_details (id, inventory_id, denomination, count, subtotal, sort_order)
				 VALUES (?,?,?,?,?,?)`,
				fmt.Sprintf("%s-D%d", m.ID, i), m.ID, d.Denomination, d.Count, d.Subtotal, i,
			).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	inv.ID = m.ID
	return nil
}

func (r *PGCashRepo) GetInventory(ctx context.Context, id string) (*domain.CashInventory, error) {
	var m domain.CashInventoryGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrCashInventoryNotFound
	}
	inv := gormInventoryToDomain(&m)

	type denomRow struct {
		Denomination float64
		Count        int
		Subtotal     float64
	}
	var dens []denomRow
	if err := r.db.WithContext(ctx).Raw(
		`SELECT denomination, count, subtotal FROM cash_inventory_details
		 WHERE inventory_id = ? ORDER BY sort_order`, id).Scan(&dens).Error; err != nil {
		return nil, err
	}
	for _, d := range dens {
		inv.Denominations = append(inv.Denominations, domain.DenominationDetail{
			Denomination: d.Denomination,
			Count:        d.Count,
			Subtotal:     d.Subtotal,
		})
	}
	return inv, nil
}

func (r *PGCashRepo) ListInventories(ctx context.Context, companyID string) ([]domain.CashInventory, error) {
	var models []domain.CashInventoryGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("inventory_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CashInventory, len(models))
	for i := range models {
		out[i] = *gormInventoryToDomain(&models[i])
	}
	return out, nil
}

func (r *PGCashRepo) UpdateInventory(ctx context.Context, inv *domain.CashInventory) error {
	m := domainToGormInventory(inv)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.CashInventoryGORM{}).Where("id = ?", inv.ID).Updates(map[string]interface{}{
			"inventory_date": m.InventoryDate,
			"cash_account_id": m.Cashier,
			"currency":       "",
			"book_balance":   m.TotalCash,
			"actual_balance": m.CountedAmount,
			"difference":     m.Difference,
			"status":         m.Status,
		}).Error; err != nil {
			return err
		}
		if err := tx.Exec(`DELETE FROM cash_inventory_details WHERE inventory_id = ?`, inv.ID).Error; err != nil {
			return err
		}
		for i, d := range inv.Denominations {
			if err := tx.Exec(
				`INSERT INTO cash_inventory_details (id, inventory_id, denomination, count, subtotal, sort_order)
				 VALUES (?,?,?,?,?,?)`,
				fmt.Sprintf("%s-D%d", inv.ID, i), inv.ID, d.Denomination, d.Count, d.Subtotal, i,
			).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ─── Advance Request / Settlement ────────────────────────────────────────

func (r *PGCashRepo) CreateAdvance(ctx context.Context, a *domain.AdvanceRequest) error {
	m := domainToGormAdvance(a)
	if a.ID == "" {
		m.ID = fmt.Sprintf("ADV-%d", time.Now().UnixNano())
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	a.ID = m.ID
	return nil
}

func (r *PGCashRepo) GetAdvance(ctx context.Context, id string) (*domain.AdvanceRequest, error) {
	var m domain.AdvanceRequestGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrAdvanceNotFound
	}
	return gormAdvanceToDomain(&m), nil
}

func (r *PGCashRepo) ListAdvances(ctx context.Context, companyID string) ([]domain.AdvanceRequest, error) {
	var models []domain.AdvanceRequestGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.AdvanceRequest, len(models))
	for i := range models {
		out[i] = *gormAdvanceToDomain(&models[i])
	}
	return out, nil
}

func (r *PGCashRepo) UpdateAdvance(ctx context.Context, a *domain.AdvanceRequest) error {
	m := domainToGormAdvance(a)
	return r.db.WithContext(ctx).Model(&domain.AdvanceRequestGORM{}).Where("id = ?", a.ID).Updates(map[string]interface{}{
		"employee_id": m.EmployeeID,
		"amount":      m.Amount,
		"currency":    m.Currency,
		"purpose":     m.Purpose,
		"status":      m.Status,
		"approved_by": m.ApprovedBy,
		"approved_at": m.ApprovedAt,
		"gl_je_id":    m.GLJEID,
		"updated_at":  time.Now(),
	}).Error
}

func (r *PGCashRepo) ListAdvancesByStatus(ctx context.Context, companyID string, status domain.AdvanceStatus) ([]domain.AdvanceRequest, error) {
	var models []domain.AdvanceRequestGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND status = ?", companyID, string(status)).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.AdvanceRequest, len(models))
	for i := range models {
		out[i] = *gormAdvanceToDomain(&models[i])
	}
	return out, nil
}

func (r *PGCashRepo) CreateAdvanceSettlement(ctx context.Context, s *domain.AdvanceSettlement) error {
	m := domainToGormAdvanceSettlement(s)
	if s.ID == "" {
		m.ID = fmt.Sprintf("ADVS-%d", time.Now().UnixNano())
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	s.ID = m.ID
	return nil
}
