package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGBankRepo struct {
	pool *pgxpool.Pool
}

func NewPGBankRepo(pool *pgxpool.Pool) *PGBankRepo {
	return &PGBankRepo{pool: pool}
}

// ─── Statements ──────────────────────────────────────────────────────────

func (r *PGBankRepo) CreateStatement(ctx context.Context, s *domain.BankStatement) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO bank_statements (id, company_id, bank_account_id,
			statement_date, from_date, to_date, opening_balance, closing_balance,
			total_credits, total_debits, line_count, currency, status,
			import_method, raw_file_name, raw_file_hash, imported_by, imported_at, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		RETURNING id`,
		s.ID, s.CompanyID, s.BankAccountID,
		s.StatementDate, s.FromDate, s.ToDate, s.OpeningBalance, s.ClosingBalance,
		s.TotalCredits, s.TotalDebits, s.LineCount, s.Currency, s.Status,
		s.ImportMethod, s.RawFileName, s.RawFileHash, s.ImportedBy, now, s.Notes,
	).Scan(&s.ID)
}

func scanStatement(row pgx.CollectableRow) (domain.BankStatement, error) {
	var s domain.BankStatement
	err := row.Scan(&s.ID, &s.CompanyID, &s.BankAccountID,
		&s.StatementDate, &s.FromDate, &s.ToDate, &s.OpeningBalance, &s.ClosingBalance,
		&s.TotalCredits, &s.TotalDebits, &s.LineCount, &s.Currency, &s.Status,
		&s.ImportMethod, &s.RawFileName, &s.RawFileHash, &s.ImportedBy, &s.ImportedAt, &s.Notes)
	return s, err
}

func (r *PGBankRepo) GetStatement(ctx context.Context, id string) (*domain.BankStatement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, bank_account_id,
			statement_date, from_date, to_date, opening_balance, closing_balance,
			total_credits, total_debits, line_count, currency, status,
			import_method, raw_file_name, raw_file_hash, imported_by, imported_at, notes
		FROM bank_statements WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	s, err := pgx.CollectOneRow(rows, scanStatement)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrBankStatementNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *PGBankRepo) ListStatements(ctx context.Context, companyID, bankAccountID string, limit, offset int) ([]domain.BankStatement, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM bank_statements WHERE company_id=$1 AND bank_account_id=$2`,
		companyID, bankAccountID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, bank_account_id,
			statement_date, from_date, to_date, opening_balance, closing_balance,
			total_credits, total_debits, line_count, currency, status,
			import_method, raw_file_name, raw_file_hash, imported_by, imported_at, notes
		FROM bank_statements
		WHERE company_id=$1 AND bank_account_id=$2
		ORDER BY statement_date DESC
		LIMIT $3 OFFSET $4`, companyID, bankAccountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, scanStatement)
	if err != nil {
		return nil, 0, err
	}
	if items == nil {
		items = []domain.BankStatement{}
	}
	return items, total, nil
}

func (r *PGBankRepo) DeleteStatement(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM bank_statements WHERE id=$1 AND status=$2`,
		id, domain.BankStatementImported)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrBankStatementNotFound
	}
	return nil
}

// ─── Statement Lines ──────────────────────────────────────────────────────

func (r *PGBankRepo) CreateStatementLines(ctx context.Context, lines []domain.BankStatementLine) error {
	if len(lines) == 0 {
		return nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	rows := make([][]any, len(lines))
	for i, l := range lines {
		rows[i] = []any{l.ID, l.StatementID, l.TransactionDate,
			l.ValueDate, l.Description, l.DebitAmount, l.CreditAmount,
			l.BalanceAfter, l.ReferenceNo, l.BankRef, l.Counterparty,
			l.CounterpartyAcc, l.CounterpartyBank, l.RawData,
			l.MatchStatus, l.MatchedLineID, l.MatchedAt, l.MatchedBy, now}
	}
	_, err := r.pool.CopyFrom(ctx,
		pgx.Identifier{"bank_statement_lines"},
		[]string{"id", "statement_id", "transaction_date", "value_date",
			"description", "debit_amount", "credit_amount", "balance_after",
			"reference_no", "bank_ref", "counterparty", "counterparty_acc",
			"counterparty_bank", "raw_data", "match_status", "matched_line_id",
			"matched_at", "matched_by", "created_at"},
		pgx.CopyFromRows(rows))
	return err
}

func scanStatementLine(row pgx.CollectableRow) (domain.BankStatementLine, error) {
	var l domain.BankStatementLine
	err := row.Scan(&l.ID, &l.StatementID, &l.TransactionDate,
		&l.ValueDate, &l.Description, &l.DebitAmount, &l.CreditAmount,
		&l.BalanceAfter, &l.ReferenceNo, &l.BankRef, &l.Counterparty,
		&l.CounterpartyAcc, &l.CounterpartyBank, &l.RawData,
		&l.MatchStatus, &l.MatchedLineID, &l.MatchedAt, &l.MatchedBy, &l.CreatedAt)
	return l, err
}

func (r *PGBankRepo) GetStatementLines(ctx context.Context, statementID string) ([]domain.BankStatementLine, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, statement_id, transaction_date, value_date,
			description, debit_amount, credit_amount, balance_after,
			reference_no, bank_ref, counterparty, counterparty_acc,
			counterparty_bank, raw_data, match_status, matched_line_id,
			matched_at, matched_by, created_at
		FROM bank_statement_lines WHERE statement_id=$1 ORDER BY transaction_date`, statementID)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, scanStatementLine)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.BankStatementLine{}
	}
	return items, nil
}

func (r *PGBankRepo) GetStatementLinesByStatus(ctx context.Context, statementID string, status domain.MatchStatus) ([]domain.BankStatementLine, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, statement_id, transaction_date, value_date,
			description, debit_amount, credit_amount, balance_after,
			reference_no, bank_ref, counterparty, counterparty_acc,
			counterparty_bank, raw_data, match_status, matched_line_id,
			matched_at, matched_by, created_at
		FROM bank_statement_lines
		WHERE statement_id=$1 AND match_status=$2
		ORDER BY transaction_date`, statementID, status)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, scanStatementLine)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.BankStatementLine{}
	}
	return items, nil
}

func (r *PGBankRepo) UpdateStatementLineMatch(ctx context.Context, lineID string, matchStatus domain.MatchStatus, matchedLineID, matchedBy string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	tag, err := r.pool.Exec(ctx,
		`UPDATE bank_statement_lines
		SET match_status=$1, matched_line_id=$2, matched_by=$3, matched_at=$4
		WHERE id=$5`,
		matchStatus, matchedLineID, matchedBy, now, lineID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrStatementLineNotFound
	}
	return nil
}

// ─── Reconciliation ────────────────────────────────────────────────────────

func (r *PGBankRepo) CreateReconciliation(ctx context.Context, rc *domain.BankReconciliation) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO bank_reconciliations (id, company_id, bank_account_id, statement_id,
			from_date, to_date, opening_balance, closing_balance, statement_balance,
			difference, status, matched_lines, unmatched_lines, write_off_amount,
			completed_by, completed_at, notes, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		RETURNING id`,
		rc.ID, rc.CompanyID, rc.BankAccountID, rc.StatementID,
		rc.FromDate, rc.ToDate, rc.OpeningBalance, rc.ClosingBalance, rc.StatementBalance,
		rc.Difference, rc.Status, rc.MatchedLines, rc.UnmatchedLines, rc.WriteOffAmount,
		rc.CompletedBy, rc.CompletedAt, rc.Notes, now,
	).Scan(&rc.ID)
}

func scanRecon(row pgx.CollectableRow) (domain.BankReconciliation, error) {
	var rc domain.BankReconciliation
	err := row.Scan(&rc.ID, &rc.CompanyID, &rc.BankAccountID, &rc.StatementID,
		&rc.FromDate, &rc.ToDate, &rc.OpeningBalance, &rc.ClosingBalance, &rc.StatementBalance,
		&rc.Difference, &rc.Status, &rc.MatchedLines, &rc.UnmatchedLines, &rc.WriteOffAmount,
		&rc.CompletedBy, &rc.CompletedAt, &rc.ReversedAt, &rc.Notes, &rc.CreatedAt)
	return rc, err
}

func (r *PGBankRepo) GetReconciliation(ctx context.Context, id string) (*domain.BankReconciliation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, bank_account_id, statement_id,
			from_date, to_date, opening_balance, closing_balance, statement_balance,
			difference, status, matched_lines, unmatched_lines, write_off_amount,
			completed_by, completed_at, reversed_at, notes, created_at
		FROM bank_reconciliations WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	rc, err := pgx.CollectOneRow(rows, scanRecon)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrReconciliationNotFound
		}
		return nil, err
	}
	return &rc, nil
}

func (r *PGBankRepo) ListReconciliations(ctx context.Context, companyID, bankAccountID string) ([]domain.BankReconciliation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, bank_account_id, statement_id,
			from_date, to_date, opening_balance, closing_balance, statement_balance,
			difference, status, matched_lines, unmatched_lines, write_off_amount,
			completed_by, completed_at, reversed_at, notes, created_at
		FROM bank_reconciliations
		WHERE company_id=$1 AND bank_account_id=$2
		ORDER BY from_date DESC`, companyID, bankAccountID)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, scanRecon)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.BankReconciliation{}
	}
	return items, nil
}

func (r *PGBankRepo) UpdateReconciliation(ctx context.Context, rc *domain.BankReconciliation) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE bank_reconciliations
		SET closing_balance=$1, statement_balance=$2, difference=$3, status=$4,
			matched_lines=$5, unmatched_lines=$6, write_off_amount=$7,
			completed_by=$8, completed_at=$9, notes=$10
		WHERE id=$11`,
		rc.ClosingBalance, rc.StatementBalance, rc.Difference, rc.Status,
		rc.MatchedLines, rc.UnmatchedLines, rc.WriteOffAmount,
		rc.CompletedBy, rc.CompletedAt, rc.Notes, rc.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrReconciliationNotFound
	}
	return nil
}

// ─── Reconciliation Matches ───────────────────────────────────────────────

func (r *PGBankRepo) CreateReconciliationMatch(ctx context.Context, m *domain.BankReconciliationMatch) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO bank_reconciliation_matches (id, reconciliation_id, statement_line_id,
			transaction_type, transaction_id, transaction_ref, match_method, confidence, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		m.ID, m.ReconciliationID, m.StatementLineID,
		m.TransactionType, m.TransactionID, m.TransactionRef, m.MatchMethod, m.Confidence, now,
	).Scan(&m.ID)
}

func (r *PGBankRepo) GetReconciliationMatches(ctx context.Context, reconID string) ([]domain.BankReconciliationMatch, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, reconciliation_id, statement_line_id,
			transaction_type, transaction_id, transaction_ref, match_method, confidence, created_at
		FROM bank_reconciliation_matches WHERE reconciliation_id=$1`, reconID)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.BankReconciliationMatch])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.BankReconciliationMatch{}
	}
	return items, nil
}

func (r *PGBankRepo) DeleteReconciliationMatch(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM bank_reconciliation_matches WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrReconciliationNotFound
	}
	return nil
}

// ─── Payment Orders ──────────────────────────────────────────────────────

func (r *PGBankRepo) CreatePaymentOrder(ctx context.Context, po *domain.PaymentOrder) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO payment_orders (id, company_id, payment_date, amount, currency,
			exchange_rate, beneficiary_name, beneficiary_acc, beneficiary_bank,
			beneficiary_branch, beneficiary_code, from_bank_acc_id, payment_content,
			urgent, payment_type, status, created_by, approved_by, approved_at,
			submitted_at, bank_ref, failure_reason, error_code, print_count,
			created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26) RETURNING id`,
		po.ID, po.CompanyID, po.PaymentDate, po.Amount, po.Currency,
		po.ExchangeRate, po.BeneficiaryName, po.BeneficiaryAcc, po.BeneficiaryBank,
		po.BeneficiaryBranch, po.BeneficiaryCode, po.FromBankAccID, po.PaymentContent,
		po.Urgent, po.PaymentType, po.Status, po.CreatedBy,
		po.ApprovedBy, po.ApprovedAt, po.SubmittedAt, po.BankRef,
		po.FailureReason, po.ErrorCode, po.PrintCount, now, now,
	).Scan(&po.ID)
}

func scanPaymentOrder(row pgx.CollectableRow) (domain.PaymentOrder, error) {
	var po domain.PaymentOrder
	err := row.Scan(&po.ID, &po.CompanyID, &po.PaymentDate, &po.Amount, &po.Currency,
		&po.ExchangeRate, &po.BeneficiaryName, &po.BeneficiaryAcc, &po.BeneficiaryBank,
		&po.BeneficiaryBranch, &po.BeneficiaryCode, &po.FromBankAccID, &po.PaymentContent,
		&po.Urgent, &po.PaymentType, &po.Status, &po.CreatedBy, &po.ApprovedBy,
		&po.ApprovedAt, &po.SubmittedAt, &po.BankRef, &po.FailureReason, &po.ErrorCode,
		&po.PrintCount, &po.CreatedAt, &po.UpdatedAt)
	return po, err
}

func (r *PGBankRepo) GetPaymentOrder(ctx context.Context, id string) (*domain.PaymentOrder, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, payment_date, amount, currency, exchange_rate,
			beneficiary_name, beneficiary_acc, beneficiary_bank, beneficiary_branch,
			beneficiary_code, from_bank_acc_id, payment_content, urgent, payment_type,
			status, created_by, approved_by, approved_at, submitted_at, bank_ref,
			failure_reason, error_code, print_count, created_at, updated_at
		FROM payment_orders WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	po, err := pgx.CollectOneRow(rows, scanPaymentOrder)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrPaymentOrderNotFound
		}
		return nil, err
	}
	return &po, nil
}

func (r *PGBankRepo) ListPaymentOrders(ctx context.Context, filter domain.PaymentOrderFilter) ([]domain.PaymentOrder, int, error) {
	where := []string{"company_id=$1"}
	args := []any{filter.CompanyID}
	argIdx := 2

	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status=$%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.PaymentType != "" {
		where = append(where, fmt.Sprintf("payment_type=$%d", argIdx))
		args = append(args, filter.PaymentType)
		argIdx++
	}
	if filter.FromDate != "" {
		where = append(where, fmt.Sprintf("payment_date>=$%d", argIdx))
		args = append(args, filter.FromDate)
		argIdx++
	}
	if filter.ToDate != "" {
		where = append(where, fmt.Sprintf("payment_date<=$%d", argIdx))
		args = append(args, filter.ToDate)
		argIdx++
	}

	wClause := strings.Join(where, " AND ")

	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payment_orders WHERE `+wClause, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	sql := `SELECT id, company_id, payment_date, amount, currency, exchange_rate,
		beneficiary_name, beneficiary_acc, beneficiary_bank, beneficiary_branch,
		beneficiary_code, from_bank_acc_id, payment_content, urgent, payment_type,
		status, created_by, approved_by, approved_at, submitted_at, bank_ref,
		failure_reason, error_code, print_count, created_at, updated_at
	FROM payment_orders WHERE ` + wClause + ` ORDER BY payment_date DESC`

	args = append(args, limit, offset)
	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, scanPaymentOrder)
	if err != nil {
		return nil, 0, err
	}
	if items == nil {
		items = []domain.PaymentOrder{}
	}
	return items, total, nil
}

func (r *PGBankRepo) UpdatePaymentOrder(ctx context.Context, po *domain.PaymentOrder) error {
	now := time.Now().UTC().Format(time.RFC3339)
	tag, err := r.pool.Exec(ctx,
		`UPDATE payment_orders
		SET payment_date=$1, amount=$2, currency=$3, exchange_rate=$4,
			beneficiary_name=$5, beneficiary_acc=$6, beneficiary_bank=$7,
			beneficiary_branch=$8, beneficiary_code=$9, payment_content=$10,
			urgent=$11, payment_type=$12, status=$13, approved_by=$14,
			approved_at=$15, submitted_at=$16, bank_ref=$17, failure_reason=$18,
			error_code=$19, updated_at=$20
		WHERE id=$21`,
		po.PaymentDate, po.Amount, po.Currency, po.ExchangeRate,
		po.BeneficiaryName, po.BeneficiaryAcc, po.BeneficiaryBank,
		po.BeneficiaryBranch, po.BeneficiaryCode, po.PaymentContent,
		po.Urgent, po.PaymentType, po.Status, po.ApprovedBy,
		po.ApprovedAt, po.SubmittedAt, po.BankRef, po.FailureReason,
		po.ErrorCode, now, po.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPaymentOrderNotFound
	}
	return nil
}

// ─── Payment Order Batches ───────────────────────────────────────────────

func (r *PGBankRepo) CreatePaymentOrderBatch(ctx context.Context, b *domain.PaymentOrderBatch) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO payment_order_batches (id, company_id, batch_name, batch_date,
			total_amount, currency, order_count, status, created_by, submitted_at, bank_ref, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`,
		b.ID, b.CompanyID, b.BatchName, b.BatchDate,
		b.TotalAmount, b.Currency, b.OrderCount, b.Status, b.CreatedBy,
		b.SubmittedAt, b.BankRef, now,
	).Scan(&b.ID)
}

func scanPaymentBatch(row pgx.CollectableRow) (domain.PaymentOrderBatch, error) {
	var b domain.PaymentOrderBatch
	err := row.Scan(&b.ID, &b.CompanyID, &b.BatchName, &b.BatchDate,
		&b.TotalAmount, &b.Currency, &b.OrderCount, &b.Status, &b.CreatedBy,
		&b.SubmittedAt, &b.BankRef, &b.CreatedAt)
	return b, err
}

func (r *PGBankRepo) GetPaymentOrderBatch(ctx context.Context, id string) (*domain.PaymentOrderBatch, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, batch_name, batch_date, total_amount, currency,
			order_count, status, created_by, submitted_at, bank_ref, created_at
		FROM payment_order_batches WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	b, err := pgx.CollectOneRow(rows, scanPaymentBatch)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrPaymentBatchNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *PGBankRepo) ListPaymentOrderBatches(ctx context.Context, companyID string) ([]domain.PaymentOrderBatch, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, batch_name, batch_date, total_amount, currency,
			order_count, status, created_by, submitted_at, bank_ref, created_at
		FROM payment_order_batches WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, scanPaymentBatch)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.PaymentOrderBatch{}
	}
	return items, nil
}

func (r *PGBankRepo) UpdatePaymentOrderBatch(ctx context.Context, b *domain.PaymentOrderBatch) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE payment_order_batches
		SET total_amount=$1, order_count=$2, status=$3, submitted_at=$4, bank_ref=$5
		WHERE id=$6`,
		b.TotalAmount, b.OrderCount, b.Status, b.SubmittedAt, b.BankRef, b.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPaymentBatchNotFound
	}
	return nil
}

func (r *PGBankRepo) AddOrdersToBatch(ctx context.Context, batchID string, orderIDs []string) error {
	if len(orderIDs) == 0 {
		return nil
	}
	rows := make([][]any, len(orderIDs))
	for i, oid := range orderIDs {
		rows[i] = []any{newUUID(), batchID, oid}
	}
	_, err := r.pool.CopyFrom(ctx,
		pgx.Identifier{"payment_order_batch_items"},
		[]string{"id", "batch_id", "order_id"},
		pgx.CopyFromRows(rows))
	return err
}

func (r *PGBankRepo) GetBatchOrderIDs(ctx context.Context, batchID string) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT order_id FROM payment_order_batch_items WHERE batch_id=$1 ORDER BY created_at`, batchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, rows.Err()
}

// ─── Loans ──────────────────────────────────────────────────────────────

func (r *PGBankRepo) CreateLoan(ctx context.Context, l *domain.LoanAgreement) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO loan_agreements (id, company_id, bank_account_id, contract_no,
			loan_type, principal_amount, currency, interest_rate, interest_method,
			base_rate, margin_rate, disbursed_amount, outstanding_balance,
			start_date, maturity_date, repayment_method, repayment_freq, status,
			notes, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21) RETURNING id`,
		l.ID, l.CompanyID, l.BankAccountID, l.ContractNo,
		l.LoanType, l.PrincipalAmount, l.Currency, l.InterestRate, l.InterestMethod,
		l.BaseRate, l.MarginRate, l.DisbursedAmount, l.OutstandingBalance,
		l.StartDate, l.MaturityDate, l.RepaymentMethod, l.RepaymentFreq, l.Status,
		l.Notes, now, now,
	).Scan(&l.ID)
}

func scanLoan(row pgx.CollectableRow) (domain.LoanAgreement, error) {
	var l domain.LoanAgreement
	err := row.Scan(&l.ID, &l.CompanyID, &l.BankAccountID, &l.ContractNo,
		&l.LoanType, &l.PrincipalAmount, &l.Currency, &l.InterestRate, &l.InterestMethod,
		&l.BaseRate, &l.MarginRate, &l.DisbursedAmount, &l.OutstandingBalance,
		&l.StartDate, &l.MaturityDate, &l.RepaymentMethod, &l.RepaymentFreq, &l.Status,
		&l.Notes, &l.CreatedAt, &l.UpdatedAt)
	return l, err
}

func (r *PGBankRepo) GetLoan(ctx context.Context, id string) (*domain.LoanAgreement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, bank_account_id, contract_no,
			loan_type, principal_amount, currency, interest_rate, interest_method,
			base_rate, margin_rate, disbursed_amount, outstanding_balance,
			start_date, maturity_date, repayment_method, repayment_freq, status,
			notes, created_at, updated_at
		FROM loan_agreements WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	l, err := pgx.CollectOneRow(rows, scanLoan)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrLoanAgreementNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (r *PGBankRepo) ListLoans(ctx context.Context, filter domain.LoanFilter) ([]domain.LoanAgreement, error) {
	where := []string{"company_id=$1"}
	args := []any{filter.CompanyID}
	argIdx := 2

	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status=$%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.LoanType != "" {
		where = append(where, fmt.Sprintf("loan_type=$%d", argIdx))
		args = append(args, filter.LoanType)
		argIdx++
	}

	wClause := strings.Join(where, " AND ")
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, bank_account_id, contract_no,
			loan_type, principal_amount, currency, interest_rate, interest_method,
			base_rate, margin_rate, disbursed_amount, outstanding_balance,
			start_date, maturity_date, repayment_method, repayment_freq, status,
			notes, created_at, updated_at
		FROM loan_agreements WHERE `+wClause+` ORDER BY start_date DESC`, args...)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, scanLoan)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.LoanAgreement{}
	}
	return items, nil
}

func (r *PGBankRepo) UpdateLoan(ctx context.Context, l *domain.LoanAgreement) error {
	now := time.Now().UTC().Format(time.RFC3339)
	tag, err := r.pool.Exec(ctx,
		`UPDATE loan_agreements
		SET disbursed_amount=$1, outstanding_balance=$2, status=$3, notes=$4, updated_at=$5
		WHERE id=$6`,
		l.DisbursedAmount, l.OutstandingBalance, l.Status, l.Notes, now, l.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrLoanAgreementNotFound
	}
	return nil
}

// ─── Disbursements ──────────────────────────────────────────────────────

func (r *PGBankRepo) CreateDisbursement(ctx context.Context, d *domain.LoanDisbursement) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO loan_disbursements (id, loan_id, disbursement_date, amount,
			to_bank_account_id, reference_no, notes, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		d.ID, d.LoanID, d.DisbursementDate, d.Amount,
		d.ToBankAccountID, d.ReferenceNo, d.Notes, now,
	).Scan(&d.ID)
}

func (r *PGBankRepo) GetDisbursements(ctx context.Context, loanID string) ([]domain.LoanDisbursement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, loan_id, disbursement_date, amount,
			to_bank_account_id, reference_no, notes, created_at
		FROM loan_disbursements WHERE loan_id=$1 ORDER BY disbursement_date`, loanID)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.LoanDisbursement])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.LoanDisbursement{}
	}
	return items, nil
}

// ─── Repayments ─────────────────────────────────────────────────────────

func (r *PGBankRepo) CreateRepayment(ctx context.Context, rp *domain.LoanRepayment) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO loan_repayments (id, loan_id, repayment_date, principal_amount,
			interest_amount, fee_amount, total_amount, payment_order_id, status,
			notes, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id`,
		rp.ID, rp.LoanID, rp.RepaymentDate, rp.PrincipalAmount,
		rp.InterestAmount, rp.FeeAmount, rp.TotalAmount, rp.PaymentOrderID, rp.Status,
		rp.Notes, now,
	).Scan(&rp.ID)
}

func (r *PGBankRepo) GetRepayments(ctx context.Context, loanID string) ([]domain.LoanRepayment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, loan_id, repayment_date, principal_amount,
			interest_amount, fee_amount, total_amount, payment_order_id, status,
			notes, created_at
		FROM loan_repayments WHERE loan_id=$1 ORDER BY repayment_date`, loanID)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.LoanRepayment])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.LoanRepayment{}
	}
	return items, nil
}

func (r *PGBankRepo) UpdateRepayment(ctx context.Context, rp *domain.LoanRepayment) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE loan_repayments
		SET status=$1, notes=$2
		WHERE id=$3`,
		rp.Status, rp.Notes, rp.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrLoanRepaymentNotFound
	}
	return nil
}

// ─── Term Deposits ──────────────────────────────────────────────────────

func (r *PGBankRepo) CreateDeposit(ctx context.Context, d *domain.TermDeposit) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO term_deposits (id, company_id, bank_account_id, deposit_no,
			amount, currency, interest_rate, term_days, start_date, maturity_date,
			interest_at_maturity, auto_renewal, renewal_term_days, status, notes, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING id`,
		d.ID, d.CompanyID, d.BankAccountID, d.DepositNo,
		d.Amount, d.Currency, d.InterestRate, d.TermDays, d.StartDate, d.MaturityDate,
		d.InterestAtMaturity, d.AutoRenewal, d.RenewalTermDays, d.Status, d.Notes, now,
	).Scan(&d.ID)
}

func (r *PGBankRepo) GetDeposit(ctx context.Context, id string) (*domain.TermDeposit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, bank_account_id, deposit_no, amount, currency,
			interest_rate, term_days, start_date, maturity_date, interest_at_maturity,
			auto_renewal, renewal_term_days, status, notes, created_at, matured_at
		FROM term_deposits WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	var d domain.TermDeposit
	err = rows.Scan(&d.ID, &d.CompanyID, &d.BankAccountID, &d.DepositNo,
		&d.Amount, &d.Currency, &d.InterestRate, &d.TermDays, &d.StartDate, &d.MaturityDate,
		&d.InterestAtMaturity, &d.AutoRenewal, &d.RenewalTermDays, &d.Status, &d.Notes,
		&d.CreatedAt, &d.MaturedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrTermDepositNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *PGBankRepo) ListDeposits(ctx context.Context, companyID string) ([]domain.TermDeposit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, bank_account_id, deposit_no, amount, currency,
			interest_rate, term_days, start_date, maturity_date, interest_at_maturity,
			auto_renewal, renewal_term_days, status, notes, created_at, matured_at
		FROM term_deposits WHERE company_id=$1 ORDER BY start_date DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.TermDeposit
	for rows.Next() {
		var d domain.TermDeposit
		err := rows.Scan(&d.ID, &d.CompanyID, &d.BankAccountID, &d.DepositNo,
			&d.Amount, &d.Currency, &d.InterestRate, &d.TermDays, &d.StartDate, &d.MaturityDate,
			&d.InterestAtMaturity, &d.AutoRenewal, &d.RenewalTermDays, &d.Status, &d.Notes,
			&d.CreatedAt, &d.MaturedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.TermDeposit{}
	}
	return items, nil
}

func (r *PGBankRepo) UpdateDeposit(ctx context.Context, d *domain.TermDeposit) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE term_deposits
		SET interest_at_maturity=$1, status=$2, notes=$3, matured_at=$4
		WHERE id=$5`,
		d.InterestAtMaturity, d.Status, d.Notes, d.MaturedAt, d.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrTermDepositNotFound
	}
	return nil
}

// ─── Reports ────────────────────────────────────────────────────────────

func (r *PGBankRepo) GetBankLedger(ctx context.Context, companyID, bankAccountID, fromDate, toDate string) (*domain.BankLedger, error) {
	ledger := &domain.BankLedger{
		CompanyID:     companyID,
		BankAccountID: bankAccountID,
		FromDate:      fromDate,
		ToDate:        toDate,
	}

	var opening float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN txn_type='credit' THEN amount ELSE -amount END), 0)
		FROM (
			SELECT 'credit' AS txn_type, amount FROM payment_orders
			WHERE company_id=$1 AND from_bank_acc_id=$2 AND payment_date < $3 AND status IN ('CONFIRMED','SUBMITTED')
			UNION ALL
			SELECT 'debit' AS txn_type, amount FROM cash_receipts
			WHERE company_id=$1 AND cash_account_id=$2 AND receipt_date < $3 AND status='POSTED'
			UNION ALL
			SELECT 'debit' AS txn_type, amount FROM cash_payments
			WHERE company_id=$1 AND cash_account_id=$2 AND payment_date < $3 AND status='POSTED'
		) t`,
		companyID, bankAccountID, fromDate).Scan(&opening)
	if err != nil {
		return nil, err
	}
	ledger.OpeningBalance = math.Round(opening*100) / 100

	rows, err := r.pool.Query(ctx,
		`SELECT transaction_date, reference, description, debit_amount, credit_amount, balance, ref_id
		FROM (
			SELECT payment_date AS transaction_date, id AS ref_id,
				beneficiary_name || ' - ' || payment_content AS description,
				0 AS debit_amount, amount AS credit_amount, 0 AS balance,
				'' AS reference
			FROM payment_orders
			WHERE company_id=$1 AND from_bank_acc_id=$2 AND payment_date>=$3 AND payment_date<=$4 AND status IN ('CONFIRMED','SUBMITTED')
			UNION ALL
			SELECT receipt_date, id, counterparty_name || ' - ' || reason AS description,
				amount, 0, 0, voucher_no
			FROM cash_receipts
			WHERE company_id=$1 AND cash_account_id=$2 AND receipt_date>=$3 AND receipt_date<=$4 AND status='POSTED'
			UNION ALL
			SELECT payment_date, id, payee_name || ' - ' || reason AS description,
				0, amount, 0, voucher_no
			FROM cash_payments
			WHERE company_id=$1 AND cash_account_id=$2 AND payment_date>=$3 AND payment_date<=$4 AND status='POSTED'
		) t ORDER BY transaction_date, reference`,
		companyID, bankAccountID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.BankLedgerEntry
	var running float64 = ledger.OpeningBalance
	var totalDebits, totalCredits float64
	lineNo := 1
	for rows.Next() {
		var e domain.BankLedgerEntry
		var ref string
		if err := rows.Scan(&e.TransactionDate, &e.RefID, &e.Description,
			&e.DebitAmount, &e.CreditAmount, &e.RunningBalance, &ref); err != nil {
			return nil, err
		}
		e.LineNo = lineNo
		e.VoucherNo = ref
		running += e.DebitAmount - e.CreditAmount
		e.RunningBalance = math.Round(running*100) / 100
		totalDebits += e.DebitAmount
		totalCredits += e.CreditAmount
		entries = append(entries, e)
		lineNo++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if entries == nil {
		entries = []domain.BankLedgerEntry{}
	}
	ledger.Entries = entries
	ledger.TotalDebits = math.Round(totalDebits*100) / 100
	ledger.TotalCredits = math.Round(totalCredits*100) / 100
	ledger.ClosingBalance = math.Round(running*100) / 100

	return ledger, nil
}

func (r *PGBankRepo) GetBalance(ctx context.Context, companyID, bankAccountID string) (float64, error) {
	var total float64
	err := r.pool.QueryRow(ctx,
		`		SELECT COALESCE((
			SELECT SUM(amount) FROM cash_receipts
			WHERE company_id=$1 AND cash_account_id=$2 AND status='POSTED'
		), 0) - COALESCE((
			SELECT SUM(amount) FROM (
				SELECT amount FROM payment_orders WHERE company_id=$1 AND from_bank_acc_id=$2 AND status IN ('CONFIRMED','SUBMITTED')
				UNION ALL
				SELECT amount FROM cash_payments WHERE company_id=$1 AND cash_account_id=$2 AND status='POSTED'
			) t
		), 0)`,
		companyID, bankAccountID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return math.Round(total*100) / 100, nil
}
