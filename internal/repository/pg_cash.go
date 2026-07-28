package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gotax/internal/domain"
)

// ─── Cash ────────────────────────────────────────────────────────────────────

type PGCashRepo struct {
	pool *pgxpool.Pool
}

func NewPGCashRepo(pool *pgxpool.Pool) *PGCashRepo {
	return &PGCashRepo{pool: pool}
}

// ─── Cash Receipt ────────────────────────────────────────────────────────────

func (r *PGCashRepo) CreateReceipt(ctx context.Context, cr *domain.CashReceipt) error {
	if cr.ID == "" {
		cr.ID = fmt.Sprintf("CR-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO cash_receipts (id, company_id, voucher_no, voucher_date, cash_account_id,
		 counterpart_id, counterpart_name, counterpart_type, currency, exchange_rate,
		 amount, amount_vnd, debit_account_id, credit_account_id, reason,
		 receipt_type, status, approved_by, approved_at, posted_by, posted_at, gl_journal_id,
		 created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,NOW(),NOW())`,
		cr.ID, cr.CompanyID, cr.VoucherNo, cr.VoucherDate, cr.CashAccountID,
		nullStr(cr.CounterpartID), nullStr(cr.CounterpartName), string(cr.CounterpartType),
		cr.Currency, cr.ExchangeRate, cr.Amount, cr.AmountVND,
		cr.DebitAccountID, cr.CreditAccountID, cr.Reason,
		string(cr.ReceiptType), string(cr.Status),
		nullStr(cr.ApprovedBy), nullStr(cr.ApprovedAt),
		nullStr(cr.PostedBy), nullStr(cr.PostedAt), nullStr(cr.GLJournalID),
		cr.CreatedBy)
	return err
}

func (r *PGCashRepo) GetReceipt(ctx context.Context, id string) (*domain.CashReceipt, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, voucher_no, voucher_date, cash_account_id,
		        COALESCE(counterpart_id,''), COALESCE(counterpart_name,''), counterpart_type,
		        currency, exchange_rate, amount, amount_vnd,
		        debit_account_id, credit_account_id, reason,
		        receipt_type, status,
		        COALESCE(approved_by,''), COALESCE(approved_at,''),
		        COALESCE(posted_by,''), COALESCE(posted_at,''), COALESCE(gl_journal_id,''),
		        created_by, created_at, updated_at
		 FROM cash_receipts WHERE id=$1`, id)
	cr := &domain.CashReceipt{}
	err := row.Scan(&cr.ID, &cr.CompanyID, &cr.VoucherNo, &cr.VoucherDate, &cr.CashAccountID,
		&cr.CounterpartID, &cr.CounterpartName, &cr.CounterpartType,
		&cr.Currency, &cr.ExchangeRate, &cr.Amount, &cr.AmountVND,
		&cr.DebitAccountID, &cr.CreditAccountID, &cr.Reason,
		&cr.ReceiptType, &cr.Status,
		&cr.ApprovedBy, &cr.ApprovedAt, &cr.PostedBy, &cr.PostedAt, &cr.GLJournalID,
		&cr.CreatedBy, &cr.CreatedAt, &cr.UpdatedAt)
	if err != nil {
		return nil, domain.ErrCashReceiptNotFound
	}
	return cr, nil
}

func scanCashReceipts(rows pgx.Rows) ([]domain.CashReceipt, error) {
	defer rows.Close()
	var out []domain.CashReceipt
	for rows.Next() {
		var cr domain.CashReceipt
		err := rows.Scan(&cr.ID, &cr.CompanyID, &cr.VoucherNo, &cr.VoucherDate, &cr.CashAccountID,
			&cr.CounterpartID, &cr.CounterpartName, &cr.CounterpartType,
			&cr.Currency, &cr.ExchangeRate, &cr.Amount, &cr.AmountVND,
			&cr.DebitAccountID, &cr.CreditAccountID, &cr.Reason,
			&cr.ReceiptType, &cr.Status,
			&cr.ApprovedBy, &cr.ApprovedAt, &cr.PostedBy, &cr.PostedAt, &cr.GLJournalID,
			&cr.CreatedBy, &cr.CreatedAt, &cr.UpdatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, cr)
	}
	return out, nil
}

func (r *PGCashRepo) ListReceipts(ctx context.Context, filter domain.CashReceiptFilter) ([]domain.CashReceipt, int, error) {
	where := " WHERE company_id=$1"
	args := []interface{}{filter.CompanyID}
	n := 2
	if filter.ReceiptType != "" {
		where += fmt.Sprintf(" AND receipt_type=$%d", n)
		args = append(args, string(filter.ReceiptType))
		n++
	}
	if filter.Currency != "" {
		where += fmt.Sprintf(" AND currency=$%d", n)
		args = append(args, filter.Currency)
		n++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, string(filter.Status))
		n++
	}
	if filter.FromDate != "" {
		where += fmt.Sprintf(" AND voucher_date>=$%d", n)
		args = append(args, filter.FromDate)
		n++
	}
	if filter.ToDate != "" {
		where += fmt.Sprintf(" AND voucher_date<=$%d", n)
		args = append(args, filter.ToDate)
		n++
	}

	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM cash_receipts"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, company_id, voucher_no, voucher_date, cash_account_id,
	                 COALESCE(counterpart_id,''), COALESCE(counterpart_name,''), counterpart_type,
	                 currency, exchange_rate, amount, amount_vnd,
	                 debit_account_id, credit_account_id, reason,
	                 receipt_type, status,
	                 COALESCE(approved_by,''), COALESCE(approved_at,''),
	                 COALESCE(posted_by,''), COALESCE(posted_at,''), COALESCE(gl_journal_id,''),
	                 created_by, created_at, updated_at
	          FROM cash_receipts` + where + " ORDER BY voucher_date DESC"

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	query += fmt.Sprintf(" LIMIT $%d", n)
	args = append(args, limit)
	n++
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, filter.Offset)
		n++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	out, err := scanCashReceipts(rows)
	if err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *PGCashRepo) UpdateReceipt(ctx context.Context, cr *domain.CashReceipt) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE cash_receipts SET voucher_no=$1, voucher_date=$2, cash_account_id=$3,
		 counterpart_id=$4, counterpart_name=$5, counterpart_type=$6,
		 currency=$7, exchange_rate=$8, amount=$9, amount_vnd=$10,
		 debit_account_id=$11, credit_account_id=$12, reason=$13,
		 receipt_type=$14, status=$15,
		 approved_by=$16, approved_at=$17, posted_by=$18, posted_at=$19, gl_journal_id=$20,
		 updated_at=NOW() WHERE id=$21`,
		cr.VoucherNo, cr.VoucherDate, cr.CashAccountID,
		nullStr(cr.CounterpartID), nullStr(cr.CounterpartName), string(cr.CounterpartType),
		cr.Currency, cr.ExchangeRate, cr.Amount, cr.AmountVND,
		cr.DebitAccountID, cr.CreditAccountID, cr.Reason,
		string(cr.ReceiptType), string(cr.Status),
		nullStr(cr.ApprovedBy), nullStr(cr.ApprovedAt),
		nullStr(cr.PostedBy), nullStr(cr.PostedAt), nullStr(cr.GLJournalID),
		cr.ID)
	return err
}

func (r *PGCashRepo) DeleteReceipt(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM cash_receipts WHERE id=$1`, id)
	return err
}

func (r *PGCashRepo) LastReceiptNo(ctx context.Context, companyID, year string) (string, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT voucher_no FROM cash_receipts
		 WHERE company_id=$1 AND voucher_date LIKE $2||'%'
		 ORDER BY voucher_no DESC LIMIT 1`, companyID, year)
	var last string
	err := row.Scan(&last)
	if err != nil {
		return "", nil // no prior receipt is OK
	}
	return last, nil
}

// ─── Cash Payment ────────────────────────────────────────────────────────────

func (r *PGCashRepo) CreatePayment(ctx context.Context, p *domain.CashPayment) error {
	if p.ID == "" {
		p.ID = fmt.Sprintf("CP-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO cash_payments (id, company_id, voucher_no, voucher_date, cash_account_id,
		 payee_id, payee_name, payee_type, currency, exchange_rate,
		 amount, amount_vnd, debit_account_id, credit_account_id, reason,
		 payment_type, status, approved_by, approved_at, posted_by, posted_at, gl_journal_id,
		 created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,NOW(),NOW())`,
		p.ID, p.CompanyID, p.VoucherNo, p.VoucherDate, p.CashAccountID,
		nullStr(p.PayeeID), nullStr(p.PayeeName), string(p.PayeeType),
		p.Currency, p.ExchangeRate, p.Amount, p.AmountVND,
		p.DebitAccountID, p.CreditAccountID, p.Reason,
		string(p.PaymentType), string(p.Status),
		nullStr(p.ApprovedBy), nullStr(p.ApprovedAt),
		nullStr(p.PostedBy), nullStr(p.PostedAt), nullStr(p.GLJournalID),
		p.CreatedBy)
	return err
}

func (r *PGCashRepo) GetPayment(ctx context.Context, id string) (*domain.CashPayment, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, voucher_no, voucher_date, cash_account_id,
		        COALESCE(payee_id,''), COALESCE(payee_name,''), payee_type,
		        currency, exchange_rate, amount, amount_vnd,
		        debit_account_id, credit_account_id, reason,
		        payment_type, status,
		        COALESCE(approved_by,''), COALESCE(approved_at,''),
		        COALESCE(posted_by,''), COALESCE(posted_at,''), COALESCE(gl_journal_id,''),
		        created_by, created_at, updated_at
		 FROM cash_payments WHERE id=$1`, id)
	p := &domain.CashPayment{}
	err := row.Scan(&p.ID, &p.CompanyID, &p.VoucherNo, &p.VoucherDate, &p.CashAccountID,
		&p.PayeeID, &p.PayeeName, &p.PayeeType,
		&p.Currency, &p.ExchangeRate, &p.Amount, &p.AmountVND,
		&p.DebitAccountID, &p.CreditAccountID, &p.Reason,
		&p.PaymentType, &p.Status,
		&p.ApprovedBy, &p.ApprovedAt, &p.PostedBy, &p.PostedAt, &p.GLJournalID,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, domain.ErrCashPaymentNotFound
	}
	return p, nil
}

func scanCashPayments(rows pgx.Rows) ([]domain.CashPayment, error) {
	defer rows.Close()
	var out []domain.CashPayment
	for rows.Next() {
		var p domain.CashPayment
		err := rows.Scan(&p.ID, &p.CompanyID, &p.VoucherNo, &p.VoucherDate, &p.CashAccountID,
			&p.PayeeID, &p.PayeeName, &p.PayeeType,
			&p.Currency, &p.ExchangeRate, &p.Amount, &p.AmountVND,
			&p.DebitAccountID, &p.CreditAccountID, &p.Reason,
			&p.PaymentType, &p.Status,
			&p.ApprovedBy, &p.ApprovedAt, &p.PostedBy, &p.PostedAt, &p.GLJournalID,
			&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *PGCashRepo) ListPayments(ctx context.Context, filter domain.CashPaymentFilter) ([]domain.CashPayment, int, error) {
	where := " WHERE company_id=$1"
	args := []interface{}{filter.CompanyID}
	n := 2
	if filter.PaymentType != "" {
		where += fmt.Sprintf(" AND payment_type=$%d", n)
		args = append(args, string(filter.PaymentType))
		n++
	}
	if filter.Currency != "" {
		where += fmt.Sprintf(" AND currency=$%d", n)
		args = append(args, filter.Currency)
		n++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, string(filter.Status))
		n++
	}
	if filter.FromDate != "" {
		where += fmt.Sprintf(" AND voucher_date>=$%d", n)
		args = append(args, filter.FromDate)
		n++
	}
	if filter.ToDate != "" {
		where += fmt.Sprintf(" AND voucher_date<=$%d", n)
		args = append(args, filter.ToDate)
		n++
	}

	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM cash_payments"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, company_id, voucher_no, voucher_date, cash_account_id,
	                 COALESCE(payee_id,''), COALESCE(payee_name,''), payee_type,
	                 currency, exchange_rate, amount, amount_vnd,
	                 debit_account_id, credit_account_id, reason,
	                 payment_type, status,
	                 COALESCE(approved_by,''), COALESCE(approved_at,''),
	                 COALESCE(posted_by,''), COALESCE(posted_at,''), COALESCE(gl_journal_id,''),
	                 created_by, created_at, updated_at
	          FROM cash_payments` + where + " ORDER BY voucher_date DESC"

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	query += fmt.Sprintf(" LIMIT $%d", n)
	args = append(args, limit)
	n++
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, filter.Offset)
		n++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	out, err := scanCashPayments(rows)
	if err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *PGCashRepo) UpdatePayment(ctx context.Context, p *domain.CashPayment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE cash_payments SET voucher_no=$1, voucher_date=$2, cash_account_id=$3,
		 payee_id=$4, payee_name=$5, payee_type=$6,
		 currency=$7, exchange_rate=$8, amount=$9, amount_vnd=$10,
		 debit_account_id=$11, credit_account_id=$12, reason=$13,
		 payment_type=$14, status=$15,
		 approved_by=$16, approved_at=$17, posted_by=$18, posted_at=$19, gl_journal_id=$20,
		 updated_at=NOW() WHERE id=$21`,
		p.VoucherNo, p.VoucherDate, p.CashAccountID,
		nullStr(p.PayeeID), nullStr(p.PayeeName), string(p.PayeeType),
		p.Currency, p.ExchangeRate, p.Amount, p.AmountVND,
		p.DebitAccountID, p.CreditAccountID, p.Reason,
		string(p.PaymentType), string(p.Status),
		nullStr(p.ApprovedBy), nullStr(p.ApprovedAt),
		nullStr(p.PostedBy), nullStr(p.PostedAt), nullStr(p.GLJournalID),
		p.ID)
	return err
}

func (r *PGCashRepo) DeletePayment(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM cash_payments WHERE id=$1`, id)
	return err
}

func (r *PGCashRepo) LastPaymentNo(ctx context.Context, companyID, year string) (string, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT voucher_no FROM cash_payments
		 WHERE company_id=$1 AND voucher_date LIKE $2||'%'
		 ORDER BY voucher_no DESC LIMIT 1`, companyID, year)
	var last string
	err := row.Scan(&last)
	if err != nil {
		return "", nil
	}
	return last, nil
}

// ─── Cash Transfer ───────────────────────────────────────────────────────────

func (r *PGCashRepo) CreateTransfer(ctx context.Context, t *domain.CashTransfer) error {
	if t.ID == "" {
		t.ID = fmt.Sprintf("CTF-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO cash_transfers (id, company_id, transfer_date, from_account_id, to_account_id,
		 amount, currency, exchange_rate, reason, transfer_type, status,
		 source_voucher_id, dest_voucher_id, created_at, posted_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,NOW(),$14)`,
		t.ID, t.CompanyID, t.TransferDate, t.FromAccountID, t.ToAccountID,
		t.Amount, t.Currency, t.ExchangeRate, t.Reason, string(t.TransferType),
		string(t.Status), nullStr(t.SourceVoucherID), nullStr(t.DestVoucherID),
		nullStr(t.PostedAt))
	return err
}

func (r *PGCashRepo) GetTransfer(ctx context.Context, id string) (*domain.CashTransfer, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, transfer_date, from_account_id, to_account_id,
		        amount, currency, exchange_rate, reason, transfer_type, status,
		        COALESCE(source_voucher_id,''), COALESCE(dest_voucher_id,''),
		        created_at, COALESCE(posted_at,'')
		 FROM cash_transfers WHERE id=$1`, id)
	t := &domain.CashTransfer{}
	err := row.Scan(&t.ID, &t.CompanyID, &t.TransferDate, &t.FromAccountID, &t.ToAccountID,
		&t.Amount, &t.Currency, &t.ExchangeRate, &t.Reason, &t.TransferType, &t.Status,
		&t.SourceVoucherID, &t.DestVoucherID, &t.CreatedAt, &t.PostedAt)
	if err != nil {
		return nil, domain.ErrCashTransferNotFound
	}
	return t, nil
}

func scanCashTransfers(rows pgx.Rows) ([]domain.CashTransfer, error) {
	defer rows.Close()
	var out []domain.CashTransfer
	for rows.Next() {
		var t domain.CashTransfer
		err := rows.Scan(&t.ID, &t.CompanyID, &t.TransferDate, &t.FromAccountID, &t.ToAccountID,
			&t.Amount, &t.Currency, &t.ExchangeRate, &t.Reason, &t.TransferType, &t.Status,
			&t.SourceVoucherID, &t.DestVoucherID, &t.CreatedAt, &t.PostedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *PGCashRepo) ListTransfers(ctx context.Context, companyID string) ([]domain.CashTransfer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, transfer_date, from_account_id, to_account_id,
		        amount, currency, exchange_rate, reason, transfer_type, status,
		        COALESCE(source_voucher_id,''), COALESCE(dest_voucher_id,''),
		        created_at, COALESCE(posted_at,'')
		 FROM cash_transfers WHERE company_id=$1 ORDER BY transfer_date DESC`, companyID)
	if err != nil {
		return nil, err
	}
	return scanCashTransfers(rows)
}

// ─── Cash Book / Balance ─────────────────────────────────────────────────────

func (r *PGCashRepo) GetCashBook(ctx context.Context, companyID, currency, accountID, fromDate, toDate string) (*domain.CashBook, error) {
	// opening balance: sum of posted receipts minus payments before fromDate
	var opening float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE((
			SELECT SUM(CASE WHEN rpt.type='RECEIPT' THEN rpt.amt ELSE -rpt.amt END)
			FROM (
				SELECT 'RECEIPT' AS type, amount_vnd AS amt, voucher_date
				FROM cash_receipts
				WHERE company_id=$1 AND currency=$2 AND cash_account_id=$3 AND status='POSTED' AND voucher_date<$4
				UNION ALL
				SELECT 'PAYMENT', amount_vnd, voucher_date
				FROM cash_payments
				WHERE company_id=$1 AND currency=$2 AND cash_account_id=$3 AND status='POSTED' AND voucher_date<$4
			) rpt
		),0)`, companyID, currency, accountID, fromDate).Scan(&opening)
	if err != nil {
		return nil, fmt.Errorf("opening balance: %w", err)
	}

	// entries in date range
	rows, err := r.pool.Query(ctx,
		`SELECT voucher_date, 'RECEIPT' AS voucher_type, voucher_no, reason,
		        amount_vnd, 0 AS payment_amount, id
		 FROM cash_receipts
		 WHERE company_id=$1 AND currency=$2 AND cash_account_id=$3 AND status='POSTED'
		   AND voucher_date>=$4 AND voucher_date<=$5
		 UNION ALL
		 SELECT voucher_date, 'PAYMENT', voucher_no, reason,
		        0, amount_vnd, id
		 FROM cash_payments
		 WHERE company_id=$1 AND currency=$2 AND cash_account_id=$3 AND status='POSTED'
		   AND voucher_date>=$4 AND voucher_date<=$5
		 ORDER BY voucher_date, voucher_type`,
		companyID, currency, accountID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("cash book entries: %w", err)
	}
	defer rows.Close()

	cb := &domain.CashBook{
		CompanyID:      companyID,
		Currency:       currency,
		AccountID:      accountID,
		FromDate:       fromDate,
		ToDate:         toDate,
		OpeningBalance: opening,
	}
	running := opening
	var lineNo int
	for rows.Next() {
		lineNo++
		var e domain.CashBookEntry
		var vtype, vno, desc, refID string
		var rcpAmt, payAmt float64
		if err := rows.Scan(&e.VoucherDate, &vtype, &vno, &desc, &rcpAmt, &payAmt, &refID); err != nil {
			return nil, err
		}
		e.LineNo = lineNo
		e.VoucherType = vtype
		e.VoucherNo = vno
		e.Description = desc
		e.ReceiptAmount = rcpAmt
		e.PaymentAmount = payAmt
		running += rcpAmt - payAmt
		e.RunningBalance = running
		e.RefID = refID
		cb.Entries = append(cb.Entries, e)
	}

	cb.TotalReceipts = 0
	cb.TotalPayments = 0
	for _, e := range cb.Entries {
		cb.TotalReceipts += e.ReceiptAmount
		cb.TotalPayments += e.PaymentAmount
	}
	cb.ClosingBalance = opening + cb.TotalReceipts - cb.TotalPayments

	return cb, nil
}

func (r *PGCashRepo) GetBalance(ctx context.Context, companyID, accountID string) (float64, error) {
	var balance float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE((
			SELECT SUM(amount_vnd) FROM cash_receipts
			WHERE company_id=$1 AND cash_account_id=$2 AND status='POSTED'
		),0) - COALESCE((
			SELECT SUM(amount_vnd) FROM cash_payments
			WHERE company_id=$1 AND cash_account_id=$2 AND status='POSTED'
		),0)`, companyID, accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// ─── Petty Cash Fund ─────────────────────────────────────────────────────────

func (r *PGCashRepo) CreatePettyCashFund(ctx context.Context, f *domain.PettyCashFund) error {
	if f.ID == "" {
		f.ID = fmt.Sprintf("PCF-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO petty_cash_funds (id, company_id, fund_code, fund_name, custodian_id,
		 initial_amount, current_balance, currency, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW())`,
		f.ID, f.CompanyID, f.FundCode, f.FundName, f.CustodianID,
		f.InitialAmount, f.CurrentBalance, f.Currency, string(f.Status))
	return err
}

func (r *PGCashRepo) GetPettyCashFund(ctx context.Context, id string) (*domain.PettyCashFund, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, fund_code, fund_name, custodian_id,
		        initial_amount, current_balance, currency, status, created_at
		 FROM petty_cash_funds WHERE id=$1`, id)
	f := &domain.PettyCashFund{}
	err := row.Scan(&f.ID, &f.CompanyID, &f.FundCode, &f.FundName, &f.CustodianID,
		&f.InitialAmount, &f.CurrentBalance, &f.Currency, &f.Status, &f.CreatedAt)
	if err != nil {
		return nil, domain.ErrPettyCashFundNotFound
	}
	return f, nil
}

func scanPettyCashFunds(rows pgx.Rows) ([]domain.PettyCashFund, error) {
	defer rows.Close()
	var out []domain.PettyCashFund
	for rows.Next() {
		var f domain.PettyCashFund
		err := rows.Scan(&f.ID, &f.CompanyID, &f.FundCode, &f.FundName, &f.CustodianID,
			&f.InitialAmount, &f.CurrentBalance, &f.Currency, &f.Status, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

func (r *PGCashRepo) ListPettyCashFunds(ctx context.Context, companyID string) ([]domain.PettyCashFund, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, fund_code, fund_name, custodian_id,
		        initial_amount, current_balance, currency, status, created_at
		 FROM petty_cash_funds WHERE company_id=$1 ORDER BY fund_code`, companyID)
	if err != nil {
		return nil, err
	}
	return scanPettyCashFunds(rows)
}

func (r *PGCashRepo) UpdatePettyCashFund(ctx context.Context, f *domain.PettyCashFund) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE petty_cash_funds SET fund_code=$1, fund_name=$2, custodian_id=$3,
		 initial_amount=$4, current_balance=$5, currency=$6, status=$7
		 WHERE id=$8`,
		f.FundCode, f.FundName, f.CustodianID,
		f.InitialAmount, f.CurrentBalance, f.Currency, string(f.Status), f.ID)
	return err
}

// ─── Cash Inventory ──────────────────────────────────────────────────────────

func (r *PGCashRepo) CreateInventory(ctx context.Context, inv *domain.CashInventory) error {
	if inv.ID == "" {
		inv.ID = fmt.Sprintf("CI-%d", time.Now().UnixNano())
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO cash_inventories (id, company_id, inventory_date, cash_account_id,
		 currency, book_balance, actual_balance, difference, difference_type,
		 reason, status, approved_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW())`,
		inv.ID, inv.CompanyID, inv.InventoryDate, inv.CashAccountID,
		inv.Currency, inv.BookBalance, inv.ActualBalance, inv.Difference, inv.DifferenceType,
		nullStr(inv.Reason), string(inv.Status), nullStr(inv.ApprovedBy))
	if err != nil {
		return err
	}

	for i, d := range inv.Denominations {
		_, err = tx.Exec(ctx,
			`INSERT INTO cash_inventory_details (id, inventory_id, denomination, count, subtotal, sort_order)
			 VALUES ($1,$2,$3,$4,$5,$6)`,
			fmt.Sprintf("%s-D%d", inv.ID, i), inv.ID, d.Denomination, d.Count, d.Subtotal, i)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PGCashRepo) GetInventory(ctx context.Context, id string) (*domain.CashInventory, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, inventory_date, cash_account_id,
		        currency, book_balance, actual_balance, difference, difference_type,
		        COALESCE(reason,''), status, COALESCE(approved_by,''), created_at
		 FROM cash_inventories WHERE id=$1`, id)
	inv := &domain.CashInventory{}
	err := row.Scan(&inv.ID, &inv.CompanyID, &inv.InventoryDate, &inv.CashAccountID,
		&inv.Currency, &inv.BookBalance, &inv.ActualBalance, &inv.Difference, &inv.DifferenceType,
		&inv.Reason, &inv.Status, &inv.ApprovedBy, &inv.CreatedAt)
	if err != nil {
		return nil, domain.ErrCashInventoryNotFound
	}

	// load denominations
	drows, err := r.pool.Query(ctx,
		`SELECT denomination, count, subtotal FROM cash_inventory_details
		 WHERE inventory_id=$1 ORDER BY sort_order`, id)
	if err != nil {
		return nil, err
	}
	defer drows.Close()
	for drows.Next() {
		var d domain.DenominationDetail
		if err := drows.Scan(&d.Denomination, &d.Count, &d.Subtotal); err != nil {
			return nil, err
		}
		inv.Denominations = append(inv.Denominations, d)
	}

	return inv, nil
}

func scanCashInventories(rows pgx.Rows) ([]domain.CashInventory, error) {
	defer rows.Close()
	var out []domain.CashInventory
	for rows.Next() {
		var inv domain.CashInventory
		err := rows.Scan(&inv.ID, &inv.CompanyID, &inv.InventoryDate, &inv.CashAccountID,
			&inv.Currency, &inv.BookBalance, &inv.ActualBalance, &inv.Difference, &inv.DifferenceType,
			&inv.Reason, &inv.Status, &inv.ApprovedBy, &inv.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, nil
}

func (r *PGCashRepo) ListInventories(ctx context.Context, companyID string) ([]domain.CashInventory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, inventory_date, cash_account_id,
		        currency, book_balance, actual_balance, difference, difference_type,
		        COALESCE(reason,''), status, COALESCE(approved_by,''), created_at
		 FROM cash_inventories WHERE company_id=$1 ORDER BY inventory_date DESC`, companyID)
	if err != nil {
		return nil, err
	}
	return scanCashInventories(rows)
}

func (r *PGCashRepo) UpdateInventory(ctx context.Context, inv *domain.CashInventory) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE cash_inventories SET inventory_date=$1, cash_account_id=$2,
		 currency=$3, book_balance=$4, actual_balance=$5, difference=$6,
		 difference_type=$7, reason=$8, status=$9, approved_by=$10
		 WHERE id=$11`,
		inv.InventoryDate, inv.CashAccountID,
		inv.Currency, inv.BookBalance, inv.ActualBalance, inv.Difference, inv.DifferenceType,
		nullStr(inv.Reason), string(inv.Status), nullStr(inv.ApprovedBy), inv.ID)
	if err != nil {
		return err
	}

	// replace denominations: delete old, insert new
	_, err = tx.Exec(ctx, `DELETE FROM cash_inventory_details WHERE inventory_id=$1`, inv.ID)
	if err != nil {
		return err
	}
	for i, d := range inv.Denominations {
		_, err = tx.Exec(ctx,
			`INSERT INTO cash_inventory_details (id, inventory_id, denomination, count, subtotal, sort_order)
			 VALUES ($1,$2,$3,$4,$5,$6)`,
			fmt.Sprintf("%s-D%d", inv.ID, i), inv.ID, d.Denomination, d.Count, d.Subtotal, i)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// ─── Advance Request / Settlement (PG) ────────────────────────────────

func (r *PGCashRepo) CreateAdvance(ctx context.Context, a *domain.AdvanceRequest) error {
	if a.ID == "" {
		a.ID = fmt.Sprintf("ADV-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO advance_requests (id, company_id, requestor_id, requestor_name,
		 amount, amount_vnd, currency, exchange_rate, purpose, status,
		 created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		a.ID, a.CompanyID, a.RequestorID, nullStr(a.RequestorName),
		a.Amount, a.AmountVND, a.Currency, a.ExchangeRate, a.Purpose, string(a.Status),
		a.CreatedBy, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *PGCashRepo) GetAdvance(ctx context.Context, id string) (*domain.AdvanceRequest, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, requestor_id, COALESCE(requestor_name,''),
		 amount, amount_vnd, currency, exchange_rate, purpose, status,
		 COALESCE(approved_by,''), COALESCE(approved_at,''),
		 COALESCE(paid_by,''), COALESCE(paid_at,''),
		 COALESCE(gl_journal_id,''),
		 created_by, created_at, updated_at
		 FROM advance_requests WHERE id=$1`, id)
	a := &domain.AdvanceRequest{}
	err := row.Scan(&a.ID, &a.CompanyID, &a.RequestorID, &a.RequestorName,
		&a.Amount, &a.AmountVND, &a.Currency, &a.ExchangeRate, &a.Purpose, &a.Status,
		&a.ApprovedBy, &a.ApprovedAt, &a.PaidBy, &a.PaidAt, &a.GLJournalID,
		&a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, domain.ErrAdvanceNotFound
	}
	return a, nil
}

func (r *PGCashRepo) ListAdvances(ctx context.Context, companyID string) ([]domain.AdvanceRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, requestor_id, COALESCE(requestor_name,''),
		 amount, amount_vnd, currency, exchange_rate, purpose, status,
		 COALESCE(approved_by,''), COALESCE(approved_at,''),
		 COALESCE(paid_by,''), COALESCE(paid_at,''),
		 COALESCE(gl_journal_id,''),
		 created_by, created_at, updated_at
		 FROM advance_requests WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	return scanAdvances(rows)
}

func scanAdvances(rows pgx.Rows) ([]domain.AdvanceRequest, error) {
	defer rows.Close()
	var out []domain.AdvanceRequest
	for rows.Next() {
		var a domain.AdvanceRequest
		err := rows.Scan(&a.ID, &a.CompanyID, &a.RequestorID, &a.RequestorName,
			&a.Amount, &a.AmountVND, &a.Currency, &a.ExchangeRate, &a.Purpose, &a.Status,
			&a.ApprovedBy, &a.ApprovedAt, &a.PaidBy, &a.PaidAt, &a.GLJournalID,
			&a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (r *PGCashRepo) UpdateAdvance(ctx context.Context, a *domain.AdvanceRequest) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE advance_requests SET requestor_id=$1, requestor_name=$2,
		 amount=$3, amount_vnd=$4, currency=$5, exchange_rate=$6,
		 purpose=$7, status=$8, approved_by=$9, approved_at=$10,
		 paid_by=$11, paid_at=$12, gl_journal_id=$13, updated_at=$14
		 WHERE id=$15`,
		a.RequestorID, nullStr(a.RequestorName),
		a.Amount, a.AmountVND, a.Currency, a.ExchangeRate,
		a.Purpose, string(a.Status), nullStr(a.ApprovedBy), nullStr(a.ApprovedAt),
		nullStr(a.PaidBy), nullStr(a.PaidAt), nullStr(a.GLJournalID), a.UpdatedAt,
		a.ID)
	return err
}

func (r *PGCashRepo) ListAdvancesByStatus(ctx context.Context, companyID string, status domain.AdvanceStatus) ([]domain.AdvanceRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, requestor_id, COALESCE(requestor_name,''),
		 amount, amount_vnd, currency, exchange_rate, purpose, status,
		 COALESCE(approved_by,''), COALESCE(approved_at,''),
		 COALESCE(paid_by,''), COALESCE(paid_at,''),
		 COALESCE(gl_journal_id,''),
		 created_by, created_at, updated_at
		 FROM advance_requests WHERE company_id=$1 AND status=$2 ORDER BY created_at DESC`, companyID, string(status))
	if err != nil {
		return nil, err
	}
	return scanAdvances(rows)
}

func (r *PGCashRepo) CreateAdvanceSettlement(ctx context.Context, s *domain.AdvanceSettlement) error {
	if s.ID == "" {
		s.ID = fmt.Sprintf("ADVS-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO advance_settlements (id, advance_id, company_id,
		 total_spent, remaining_amount, currency, notes, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		s.ID, s.AdvanceID, s.CompanyID,
		s.TotalSpent, s.RemainingAmount, s.Currency, nullStr(s.Notes), s.Status, s.CreatedAt)
	return err
}
