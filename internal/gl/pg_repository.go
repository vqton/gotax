package gl

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Account ───────────────────────────────────────────────────────────────

type PGAccountRepo struct {
	pool *pgxpool.Pool
}

func NewPGAccountRepo(pool *pgxpool.Pool) *PGAccountRepo {
	return &PGAccountRepo{pool: pool}
}

func (r *PGAccountRepo) Create(ctx context.Context, a *Account) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO accounts (code, name, name2, type, parent_code, is_active, is_foreign, detail_by, is_parent, note)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		a.Code, a.Name, nullStr(a.Name2), a.Type, nullStr(a.ParentCode),
		a.IsActive, a.IsForeign, nullStr(string(a.DetailBy)), a.IsParent, nullStr(a.Note))
	return err
}

func (r *PGAccountRepo) GetByCode(ctx context.Context, code string) (*Account, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT code, name, COALESCE(name2,''), type, COALESCE(parent_code,''), is_active, is_foreign,
		        COALESCE(detail_by,''), is_parent, 0, COALESCE(note,'')
		 FROM accounts WHERE code=$1`, code)
	return scanAccount(row)
}

func (r *PGAccountRepo) GetAll(ctx context.Context, activeOnly bool) ([]Account, error) {
	var query string
	if activeOnly {
		query = `SELECT code, name, COALESCE(name2,''), type, COALESCE(parent_code,''), is_active, is_foreign,
		                COALESCE(detail_by,''), is_parent, 0, COALESCE(note,'')
		         FROM accounts WHERE is_active=true ORDER BY code`
	} else {
		query = `SELECT code, name, COALESCE(name2,''), type, COALESCE(parent_code,''), is_active, is_foreign,
		                COALESCE(detail_by,''), is_parent, 0, COALESCE(note,'')
		         FROM accounts ORDER BY code`
	}
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAccounts(rows)
}

func (r *PGAccountRepo) Update(ctx context.Context, a *Account) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE accounts SET name=$1, name2=$2, type=$3, parent_code=$4, is_active=$5, is_foreign=$6,
		 detail_by=$7, is_parent=$8, note=$9, updated_at=NOW() WHERE code=$10`,
		a.Name, nullStr(a.Name2), a.Type, nullStr(a.ParentCode),
		a.IsActive, a.IsForeign, nullStr(string(a.DetailBy)), a.IsParent, nullStr(a.Note), a.Code)
	return err
}

func (r *PGAccountRepo) Delete(ctx context.Context, code string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM accounts WHERE code=$1`, code)
	return err
}

func (r *PGAccountRepo) GetChildren(ctx context.Context, parentCode string) ([]Account, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT code, name, COALESCE(name2,''), type, COALESCE(parent_code,''), is_active, is_foreign,
		        COALESCE(detail_by,''), is_parent, 0, COALESCE(note,'')
		 FROM accounts WHERE parent_code=$1 ORDER BY code`, parentCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAccounts(rows)
}

// ─── Journal Entry ─────────────────────────────────────────────────────────

type PGJournalRepo struct {
	pool *pgxpool.Pool
}

func NewPGJournalRepo(pool *pgxpool.Pool) *PGJournalRepo {
	return &PGJournalRepo{pool: pool}
}

func (r *PGJournalRepo) Create(ctx context.Context, e *JournalEntry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO journal_entries (id, company_id, entry_number, voucher_type, entry_date, accounting_date,
		 period_id, description, status, currency_code, exchange_rate, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		 RETURNING id, created_at`,
		e.ID, e.CompanyID, e.EntryNumber, e.VoucherType, e.EntryDate, e.AccountingDate,
		e.PeriodID, e.Description, e.Status, e.CurrencyCode, e.ExchangeRate, e.CreatedBy,
	).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert entry: %w", err)
	}

	for _, line := range e.Lines {
		line.EntryID = e.ID
		_, err = tx.Exec(ctx,
			`INSERT INTO journal_lines (entry_id, line_number, account_code, debit_amount, credit_amount,
			 description, currency_code, foreign_amount, exchange_rate, object_id, project_id, contract_id, cost_item_id, department_id)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
			e.ID, line.LineNumber, line.AccountCode, line.DebitAmount, line.CreditAmount,
			nullStr(line.Description), nullStr(line.CurrencyCode), nullF64(line.ForeignAmount), line.ExchangeRate,
			nullStr(line.ObjectID), nullStr(line.ProjectID), nullStr(line.ContractID),
			nullStr(line.CostItemID), nullStr(line.DepartmentID))
		if err != nil {
			return fmt.Errorf("insert line: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PGJournalRepo) GetByID(ctx context.Context, id string) (*JournalEntry, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, COALESCE(company_id,'00000000-0000-0000-0000-000000000000'), entry_number, voucher_type, entry_date, accounting_date,
		        period_id, description, status, currency_code, exchange_rate, created_by,
		        COALESCE(reviewed_by,''), COALESCE(approved_by,''), created_at, posted_at, approved_at
		 FROM journal_entries WHERE id=$1`, id)
	entry, err := scanJournalEntry(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrJournalNotFound
		}
		return nil, err
	}
	lines, err := r.GetLinesByEntryID(ctx, id)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines
	return entry, nil
}

func (r *PGJournalRepo) GetByPeriod(ctx context.Context, periodID string) ([]JournalEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(company_id,'00000000-0000-0000-0000-000000000000'), entry_number, voucher_type, entry_date, accounting_date,
		        period_id, description, status, currency_code, exchange_rate, created_by,
		        COALESCE(reviewed_by,''), COALESCE(approved_by,''), created_at, posted_at, approved_at
		 FROM journal_entries WHERE period_id=$1 ORDER BY entry_date`, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJournalEntries(rows)
}

func (r *PGJournalRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]JournalEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(company_id,'00000000-0000-0000-0000-000000000000'), entry_number, voucher_type, entry_date, accounting_date,
		        period_id, description, status, currency_code, exchange_rate, created_by,
		        COALESCE(reviewed_by,''), COALESCE(approved_by,''), created_at, posted_at, approved_at
		 FROM journal_entries WHERE entry_date>=$1 AND entry_date<=$2 ORDER BY entry_date`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJournalEntries(rows)
}

func (r *PGJournalRepo) GetByStatus(ctx context.Context, status JournalEntryStatus) ([]JournalEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(company_id,'00000000-0000-0000-0000-000000000000'), entry_number, voucher_type, entry_date, accounting_date,
		        period_id, description, status, currency_code, exchange_rate, created_by,
		        COALESCE(reviewed_by,''), COALESCE(approved_by,''), created_at, posted_at, approved_at
		 FROM journal_entries WHERE status=$1 ORDER BY entry_date`, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJournalEntries(rows)
}

func (r *PGJournalRepo) GetByVoucherType(ctx context.Context, voucherType VoucherType) ([]JournalEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(company_id,'00000000-0000-0000-0000-000000000000'), entry_number, voucher_type, entry_date, accounting_date,
		        period_id, description, status, currency_code, exchange_rate, created_by,
		        COALESCE(reviewed_by,''), COALESCE(approved_by,''), created_at, posted_at, approved_at
		 FROM journal_entries WHERE voucher_type=$1 ORDER BY entry_date`, string(voucherType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJournalEntries(rows)
}

func (r *PGJournalRepo) UpdateStatus(ctx context.Context, id string, status JournalEntryStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE journal_entries SET status=$1 WHERE id=$2`, string(status), id)
	return err
}

func (r *PGJournalRepo) Update(ctx context.Context, e *JournalEntry) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE journal_entries SET entry_number=$1, voucher_type=$2, entry_date=$3, accounting_date=$4,
		 period_id=$5, description=$6, status=$7, currency_code=$8, exchange_rate=$9,
		 posted_at=$10, approved_at=$11, reviewed_by=$12, approved_by=$13
		 WHERE id=$14`,
		e.EntryNumber, e.VoucherType, e.EntryDate, e.AccountingDate,
		e.PeriodID, e.Description, e.Status, e.CurrencyCode, e.ExchangeRate,
		e.PostedAt, e.ApprovedAt, nullStr(e.ReviewedBy), nullStr(e.ApprovedBy), e.ID)
	return err
}

func (r *PGJournalRepo) Approve(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE journal_entries SET status='APPROVED', approved_by=$1, approved_at=NOW(), updated_at=NOW()
		 WHERE id=$2`, approvedBy, id)
	return err
}

func (r *PGJournalRepo) Review(ctx context.Context, id, reviewedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE journal_entries SET status='REVIEWING', reviewed_by=$1, updated_at=NOW()
		 WHERE id=$2`, reviewedBy, id)
	return err
}

func (r *PGJournalRepo) GetLinesByEntryID(ctx context.Context, entryID string) ([]JournalLine, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, entry_id, line_number, account_code, debit_amount, credit_amount,
		        COALESCE(description,''), COALESCE(currency_code,''), COALESCE(foreign_amount,0),
		        COALESCE(exchange_rate,0), COALESCE(object_id,''), COALESCE(project_id,''),
		        COALESCE(contract_id,''), COALESCE(cost_item_id,''), COALESCE(department_id,'')
		 FROM journal_lines WHERE entry_id=$1 ORDER BY line_number`, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJournalLines(rows)
}

func (r *PGJournalRepo) GetBalance(ctx context.Context, accountCode, periodID string) (*AccountBalance, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT $1 AS account_code, $2 AS period_id,
		        COALESCE(SUM(debit_amount),0) AS period_debit,
		        COALESCE(SUM(credit_amount),0) AS period_credit
		 FROM journal_lines l
		 JOIN journal_entries e ON l.entry_id = e.id
		 WHERE l.account_code=$1 AND e.period_id=$2 AND e.status='POSTED'`,
		accountCode, periodID)

	b := &AccountBalance{AccountCode: accountCode, PeriodID: periodID}
	err := row.Scan(&b.PeriodDebit, &b.PeriodCredit)
	return b, err
}

func (r *PGJournalRepo) GetTrialBalance(ctx context.Context, periodID string) ([]AccountBalance, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT l.account_code, $1 AS period_id,
		        COALESCE(SUM(l.debit_amount),0) AS period_debit,
		        COALESCE(SUM(l.credit_amount),0) AS period_credit
		 FROM journal_lines l
		 JOIN journal_entries e ON l.entry_id = e.id
		 WHERE e.period_id=$1 AND e.status='POSTED'
		 GROUP BY l.account_code
		 ORDER BY l.account_code`, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var balances []AccountBalance
	for rows.Next() {
		var b AccountBalance
		if err := rows.Scan(&b.AccountCode, &b.PeriodID, &b.PeriodDebit, &b.PeriodCredit); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

func (r *PGJournalRepo) GetFinancialStatement(ctx context.Context, periodID string, accountTypes []AccountType) ([]AccountBalance, error) {
	if len(accountTypes) == 0 {
		return nil, nil
	}
	types := make([]string, len(accountTypes))
	for i, t := range accountTypes {
		types[i] = string(t)
	}
	rows, err := r.pool.Query(ctx,
		`SELECT l.account_code, $1 AS period_id, a.type,
		        COALESCE(SUM(l.debit_amount),0) AS period_debit,
		        COALESCE(SUM(l.credit_amount),0) AS period_credit
		 FROM journal_lines l
		 JOIN journal_entries e ON l.entry_id = e.id
		 JOIN accounts a ON l.account_code = a.code
		 WHERE e.period_id=$1 AND e.status='POSTED' AND a.type = ANY($2)
		 GROUP BY l.account_code, a.type
		 ORDER BY l.account_code`, periodID, types)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var balances []AccountBalance
	for rows.Next() {
		var b AccountBalance
		if err := rows.Scan(&b.AccountCode, &b.PeriodID, &b.AccountType, &b.PeriodDebit, &b.PeriodCredit); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

// ─── Period ────────────────────────────────────────────────────────────────

type PGPeriodRepo struct {
	pool *pgxpool.Pool
}

func NewPGPeriodRepo(pool *pgxpool.Pool) *PGPeriodRepo {
	return &PGPeriodRepo{pool: pool}
}

func (r *PGPeriodRepo) Create(ctx context.Context, p *Period) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO periods (id, year, month, start_date, end_date, status)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		p.ID, p.Year, p.Month, p.StartDate, p.EndDate, string(p.Status))
	return err
}

func (r *PGPeriodRepo) GetByID(ctx context.Context, id string) (*Period, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, year, month, start_date, end_date, status FROM periods WHERE id=$1`, id)
	p := &Period{}
	err := row.Scan(&p.ID, &p.Year, &p.Month, &p.StartDate, &p.EndDate, &p.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrPeriodNotFound
		}
		return nil, err
	}
	return p, nil
}

func (r *PGPeriodRepo) GetByYearMonth(ctx context.Context, year, month int) (*Period, error) {
	id := fmt.Sprintf("P-%d-%02d", year, month)
	return r.GetByID(ctx, id)
}

func (r *PGPeriodRepo) GetAll(ctx context.Context) ([]Period, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, year, month, start_date, end_date, status FROM periods ORDER BY year, month`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var periods []Period
	for rows.Next() {
		var p Period
		if err := rows.Scan(&p.ID, &p.Year, &p.Month, &p.StartDate, &p.EndDate, &p.Status); err != nil {
			return nil, err
		}
		periods = append(periods, p)
	}
	return periods, nil
}

func (r *PGPeriodRepo) UpdateStatus(ctx context.Context, id string, status PeriodStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE periods SET status=$1, updated_at=NOW() WHERE id=$2`, string(status), id)
	return err
}

func (r *PGPeriodRepo) GetOpenPeriod(ctx context.Context) (*Period, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, year, month, start_date, end_date, status FROM periods WHERE status='OPEN' LIMIT 1`)
	p := &Period{}
	err := row.Scan(&p.ID, &p.Year, &p.Month, &p.StartDate, &p.EndDate, &p.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrPeriodNotFound
		}
		return nil, err
	}
	return p, nil
}

// ─── User ──────────────────────────────────────────────────────────────────

type PGUserRepo struct {
	pool *pgxpool.Pool
}

func NewPGUserRepo(pool *pgxpool.Pool) *PGUserRepo {
	return &PGUserRepo{pool: pool}
}

func (r *PGUserRepo) Create(ctx context.Context, u *User) error {
	if u.ID == "" {
		u.ID = fmt.Sprintf("U-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, full_name, email, role, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())`,
		u.ID, u.Username, u.PasswordHash, u.FullName, nullStr(u.Email), string(u.Role), u.IsActive)
	return err
}

func (r *PGUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, full_name, COALESCE(email,''), role, is_active, created_at, updated_at
		 FROM users WHERE id=$1`, id)
	return scanUser(row)
}

func (r *PGUserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, full_name, COALESCE(email,''), role, is_active, created_at, updated_at
		 FROM users WHERE username=$1`, username)
	u, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *PGUserRepo) GetAll(ctx context.Context) ([]User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, username, password_hash, full_name, COALESCE(email,''), role, is_active, created_at, updated_at
		 FROM users ORDER BY username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		u, err := scanUserRaw(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, nil
}

func (r *PGUserRepo) Update(ctx context.Context, u *User) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET full_name=$1, email=$2, role=$3, is_active=$4, updated_at=NOW() WHERE id=$5`,
		u.FullName, nullStr(u.Email), string(u.Role), u.IsActive, u.ID)
	return err
}

func (r *PGUserRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

// ─── Audit Log ─────────────────────────────────────────────────────────────

type PGAuditLogRepo struct {
	pool *pgxpool.Pool
}

func NewPGAuditLogRepo(pool *pgxpool.Pool) *PGAuditLogRepo {
	return &PGAuditLogRepo{pool: pool}
}

func (r *PGAuditLogRepo) Create(ctx context.Context, e *AuditEntry) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_log (user_id, username, ip_address, action, entity_type, entity_id, old_value, new_value)
		 VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8::jsonb)`,
		nullStr(e.UserID), e.Username, nullStr(e.IPAddress), string(e.Action),
		e.EntityType, nullStr(e.EntityID), nullStr(e.OldValue), nullStr(e.NewValue))
	return err
}

func (r *PGAuditLogRepo) GetByEntity(ctx context.Context, entityType, entityID string) ([]AuditEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, COALESCE(username,''), COALESCE(ip_address,''), action, entity_type,
		        COALESCE(entity_id,''), COALESCE(old_value,''), COALESCE(new_value,''), created_at
		 FROM audit_log WHERE entity_type=$1 AND entity_id=$2 ORDER BY created_at DESC`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAuditEntries(rows)
}

func (r *PGAuditLogRepo) GetByUser(ctx context.Context, userID string) ([]AuditEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, COALESCE(username,''), COALESCE(ip_address,''), action, entity_type,
		        COALESCE(entity_id,''), COALESCE(old_value,''), COALESCE(new_value,''), created_at
		 FROM audit_log WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAuditEntries(rows)
}

func (r *PGAuditLogRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]AuditEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, COALESCE(username,''), COALESCE(ip_address,''), action, entity_type,
		        COALESCE(entity_id,''), COALESCE(old_value,''), COALESCE(new_value,''), created_at
		 FROM audit_log WHERE created_at>=$1 AND created_at<=$2 ORDER BY created_at DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAuditEntries(rows)
}

func (r *PGAuditLogRepo) GetAll(ctx context.Context, limit int) ([]AuditEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, COALESCE(username,''), COALESCE(ip_address,''), action, entity_type,
		        COALESCE(entity_id,''), COALESCE(old_value,''), COALESCE(new_value,''), created_at
		 FROM audit_log ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAuditEntries(rows)
}

// ─── Exchange Rate ─────────────────────────────────────────────────────────

type PGExchangeRateRepo struct {
	pool *pgxpool.Pool
}

func NewPGExchangeRateRepo(pool *pgxpool.Pool) *PGExchangeRateRepo {
	return &PGExchangeRateRepo{pool: pool}
}

func (r *PGExchangeRateRepo) Create(ctx context.Context, rate *ExchangeRate) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO exchange_rates (currency_code, rate_date, buy_rate, sell_rate, average_rate, source)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		rate.CurrencyCode, rate.RateDate, nullF64(rate.BuyRate), nullF64(rate.SellRate),
		rate.AverageRate, nullStr(rate.Source))
	return err
}

func (r *PGExchangeRateRepo) GetByCurrencyAndDate(ctx context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, currency_code, rate_date, COALESCE(buy_rate,0), COALESCE(sell_rate,0), average_rate,
		        COALESCE(source,''), created_at
		 FROM exchange_rates WHERE currency_code=$1 AND rate_date=$2`, currencyCode, rateDate)
	xr := &ExchangeRate{}
	err := row.Scan(&xr.ID, &xr.CurrencyCode, &xr.RateDate, &xr.BuyRate, &xr.SellRate, &xr.AverageRate, &xr.Source, &xr.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrRateNotFound
		}
		return nil, err
	}
	return xr, nil
}

func (r *PGExchangeRateRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]ExchangeRate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, currency_code, rate_date, COALESCE(buy_rate,0), COALESCE(sell_rate,0), average_rate,
		        COALESCE(source,''), created_at
		 FROM exchange_rates WHERE rate_date>=$1 AND rate_date<=$2 ORDER BY rate_date, currency_code`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []ExchangeRate
	for rows.Next() {
		var xr ExchangeRate
		if err := rows.Scan(&xr.ID, &xr.CurrencyCode, &xr.RateDate, &xr.BuyRate, &xr.SellRate, &xr.AverageRate, &xr.Source, &xr.CreatedAt); err != nil {
			return nil, err
		}
		rates = append(rates, xr)
	}
	return rates, nil
}

func (r *PGExchangeRateRepo) GetAll(ctx context.Context) ([]ExchangeRate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, currency_code, rate_date, COALESCE(buy_rate,0), COALESCE(sell_rate,0), average_rate,
		        COALESCE(source,''), created_at
		 FROM exchange_rates ORDER BY rate_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []ExchangeRate
	for rows.Next() {
		var xr ExchangeRate
		if err := rows.Scan(&xr.ID, &xr.CurrencyCode, &xr.RateDate, &xr.BuyRate, &xr.SellRate, &xr.AverageRate, &xr.Source, &xr.CreatedAt); err != nil {
			return nil, err
		}
		rates = append(rates, xr)
	}
	return rates, nil
}

func (r *PGExchangeRateRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM exchange_rates WHERE id=$1`, id)
	return err
}

// ─── Closing Template ──────────────────────────────────────────────────────

type PGClosingTemplateRepo struct {
	pool *pgxpool.Pool
}

func NewPGClosingTemplateRepo(pool *pgxpool.Pool) *PGClosingTemplateRepo {
	return &PGClosingTemplateRepo{pool: pool}
}

func (r *PGClosingTemplateRepo) Create(ctx context.Context, t *ClosingTemplate) error {
	if t.ID == "" {
		t.ID = fmt.Sprintf("CT-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO closing_templates (id, name, description, sequence_order, is_active) VALUES ($1,$2,$3,$4,$5)`,
		t.ID, t.Name, nullStr(t.Description), t.SequenceOrder, t.IsActive)
	return err
}

func (r *PGClosingTemplateRepo) GetByID(ctx context.Context, id string) (*ClosingTemplate, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, COALESCE(description,''), sequence_order, is_active, created_at
		 FROM closing_templates WHERE id=$1`, id)
	ct := &ClosingTemplate{}
	err := row.Scan(&ct.ID, &ct.Name, &ct.Description, &ct.SequenceOrder, &ct.IsActive, &ct.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	return ct, nil
}

func (r *PGClosingTemplateRepo) GetAll(ctx context.Context) ([]ClosingTemplate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, COALESCE(description,''), sequence_order, is_active, created_at
		 FROM closing_templates ORDER BY sequence_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var templates []ClosingTemplate
	for rows.Next() {
		var ct ClosingTemplate
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.Description, &ct.SequenceOrder, &ct.IsActive, &ct.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, ct)
	}
	return templates, nil
}

func (r *PGClosingTemplateRepo) Update(ctx context.Context, t *ClosingTemplate) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE closing_templates SET name=$1, description=$2, sequence_order=$3, is_active=$4 WHERE id=$5`,
		t.Name, nullStr(t.Description), t.SequenceOrder, t.IsActive, t.ID)
	return err
}

func (r *PGClosingTemplateRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM closing_templates WHERE id=$1`, id)
	return err
}

// ─── Scan helpers ──────────────────────────────────────────────────────────

type scannable interface {
	Scan(dest ...interface{}) error
}

type rowScanner interface {
	scannable
}

func scanAccount(row rowScanner) (*Account, error) {
	a := &Account{}
	var detailBy, parentCode, name2, note string
	err := row.Scan(&a.Code, &a.Name, &name2, &a.Type, &parentCode, &a.IsActive, &a.IsForeign,
		&detailBy, &a.IsParent, &a.ArrearsDays, &note)
	if err != nil {
		return nil, err
	}
	a.Name2 = name2
	a.ParentCode = parentCode
	a.DetailBy = DetailBy(detailBy)
	a.Note = note
	return a, nil
}

func scanAccounts(rows pgx.Rows) ([]Account, error) {
	defer rows.Close()
	var accounts []Account
	for rows.Next() {
		a, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, *a)
	}
	return accounts, nil
}

func scanJournalEntry(row rowScanner) (*JournalEntry, error) {
	e := &JournalEntry{}
	var reviewedBy, approvedBy string
	err := row.Scan(&e.ID, &e.CompanyID, &e.EntryNumber, &e.VoucherType, &e.EntryDate, &e.AccountingDate,
		&e.PeriodID, &e.Description, &e.Status, &e.CurrencyCode, &e.ExchangeRate,
		&e.CreatedBy, &reviewedBy, &approvedBy, &e.CreatedAt, &e.PostedAt, &e.ApprovedAt)
	if err != nil {
		return nil, err
	}
	if reviewedBy != "" {
		e.ReviewedBy = reviewedBy
	}
	if approvedBy != "" {
		e.ApprovedBy = approvedBy
	}
	return e, nil
}

func scanJournalEntries(rows pgx.Rows) ([]JournalEntry, error) {
	defer rows.Close()
	var entries []JournalEntry
	for rows.Next() {
		e, err := scanJournalEntry(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, *e)
	}
	return entries, nil
}

func scanJournalLines(rows pgx.Rows) ([]JournalLine, error) {
	defer rows.Close()
	var lines []JournalLine
	for rows.Next() {
		var l JournalLine
		err := rows.Scan(&l.ID, &l.EntryID, &l.LineNumber, &l.AccountCode,
			&l.DebitAmount, &l.CreditAmount, &l.Description, &l.CurrencyCode,
			&l.ForeignAmount, &l.ExchangeRate, &l.ObjectID, &l.ProjectID,
			&l.ContractID, &l.CostItemID, &l.DepartmentID)
		if err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}
	return lines, nil
}

func scanUser(row rowScanner) (*User, error) {
	u := &User{}
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func scanUserRaw(rows pgx.Rows) (*User, error) {
	return scanUser(rows)
}

func scanAuditEntries(rows pgx.Rows) ([]AuditEntry, error) {
	defer rows.Close()
	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		err := rows.Scan(&e.ID, &e.UserID, &e.Username, &e.IPAddress, &e.Action, &e.EntityType,
			&e.EntityID, &e.OldValue, &e.NewValue, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// ─── Null helpers ──────────────────────────────────────────────────────────

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullF64(f float64) interface{} {
	if f == 0 {
		return nil
	}
	return f
}
