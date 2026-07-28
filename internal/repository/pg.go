package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gotax/internal/domain"
)

// ─── Account ─────────────────────────────────────────────────────────────────

type PGAccountRepo struct {
	pool *pgxpool.Pool
}

func NewPGAccountRepo(pool *pgxpool.Pool) *PGAccountRepo {
	return &PGAccountRepo{pool: pool}
}

// ─── Refresh Token ───────────────────────────────────────────────────────────

type PGRefreshTokenRepo struct {
	pool *pgxpool.Pool
}

func NewPGRefreshTokenRepo(pool *pgxpool.Pool) *PGRefreshTokenRepo {
	return &PGRefreshTokenRepo{pool: pool}
}

func (r *PGRefreshTokenRepo) Create(ctx context.Context, t *domain.RefreshToken) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, device_info, ip_address, expires_at, created_at, revoked_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		t.ID, t.UserID, t.TokenHash, nullStr(t.DeviceInfo), nullStr(t.IPAddress),
		t.ExpiresAt, t.CreatedAt, nullTimePtr(t.RevokedAt))
	return err
}

func (r *PGRefreshTokenRepo) GetByID(ctx context.Context, id string) (*domain.RefreshToken, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, COALESCE(device_info,''), COALESCE(ip_address,''), expires_at, created_at, revoked_at
		 FROM refresh_tokens WHERE id=$1`, id)
	var t domain.RefreshToken
	var revokedAt *time.Time
	err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.DeviceInfo, &t.IPAddress, &t.ExpiresAt, &t.CreatedAt, &revokedAt)
	if err != nil {
		return nil, err
	}
	t.RevokedAt = revokedAt
	return &t, nil
}

func (r *PGRefreshTokenRepo) GetByUserID(ctx context.Context, userID string) ([]domain.RefreshToken, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, token_hash, COALESCE(device_info,''), COALESCE(ip_address,''), expires_at, created_at, revoked_at
		 FROM refresh_tokens WHERE user_id=$1 ORDER BY created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRefreshTokens(rows)
}

func (r *PGRefreshTokenRepo) GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, COALESCE(device_info,''), COALESCE(ip_address,''), expires_at, created_at, revoked_at
		 FROM refresh_tokens WHERE token_hash=$1`, hash)
	var t domain.RefreshToken
	var revokedAt *time.Time
	err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.DeviceInfo, &t.IPAddress, &t.ExpiresAt, &t.CreatedAt, &revokedAt)
	if err != nil {
		return nil, err
	}
	t.RevokedAt = revokedAt
	return &t, nil
}

func (r *PGRefreshTokenRepo) Revoke(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *PGRefreshTokenRepo) RevokeAllByUserID(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE user_id=$1 AND revoked_at IS NULL`, userID)
	return err
}

func scanRefreshTokens(rows pgx.Rows) ([]domain.RefreshToken, error) {
	var out []domain.RefreshToken
	for rows.Next() {
		var t domain.RefreshToken
		var revokedAt *time.Time
		if err := rows.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.DeviceInfo, &t.IPAddress, &t.ExpiresAt, &t.CreatedAt, &revokedAt); err != nil {
			return nil, err
		}
		t.RevokedAt = revokedAt
		out = append(out, t)
	}
	return out, nil
}

// ─── Password Reset Token ───────────────────────────────────────────────────

type PGPasswordResetTokenRepo struct {
	pool *pgxpool.Pool
}

func NewPGPasswordResetTokenRepo(pool *pgxpool.Pool) *PGPasswordResetTokenRepo {
	return &PGPasswordResetTokenRepo{pool: pool}
}

func (r *PGPasswordResetTokenRepo) Create(ctx context.Context, t *domain.PasswordResetToken) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, created_at, used_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		t.ID, t.UserID, t.TokenHash, t.ExpiresAt, t.CreatedAt, nullTimePtr(t.UsedAt))
	return err
}

func (r *PGPasswordResetTokenRepo) GetByID(ctx context.Context, id string) (*domain.PasswordResetToken, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at, used_at
		 FROM password_reset_tokens WHERE id=$1`, id)
	var t domain.PasswordResetToken
	var usedAt *time.Time
	err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt, &usedAt)
	if err != nil {
		return nil, err
	}
	t.UsedAt = usedAt
	return &t, nil
}

func (r *PGPasswordResetTokenRepo) MarkUsed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE password_reset_tokens SET used_at=NOW() WHERE id=$1`, id)
	return err
}

// ─── Account ─────────────────────────────────────────────────────────────────

func (r *PGAccountRepo) Create(ctx context.Context, a *domain.Account) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO accounts (code, name, name2, type, parent_code, is_active, is_foreign, detail_by, is_parent, note)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		a.Code, a.Name, nullStr(a.Name2), a.Type, nullStr(a.ParentCode),
		a.IsActive, a.IsForeign, nullStr(string(a.DetailBy)), a.IsParent, nullStr(a.Note))
	return err
}

func (r *PGAccountRepo) GetByCode(ctx context.Context, code string) (*domain.Account, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT code, name, COALESCE(name2,''), type, COALESCE(parent_code,''), is_active, is_foreign,
		        COALESCE(detail_by,''), is_parent, status, COALESCE(freeze_reason,''), arrears_days, COALESCE(note,'')
		 FROM accounts WHERE code=$1`, code)
	return scanAccount(row)
}

func (r *PGAccountRepo) GetAll(ctx context.Context, activeOnly bool) ([]domain.Account, error) {
	var query string
	if activeOnly {
		query = `SELECT code, name, COALESCE(name2,''), type, COALESCE(parent_code,''), is_active, is_foreign,
		                COALESCE(detail_by,''), is_parent, status, COALESCE(freeze_reason,''), arrears_days, COALESCE(note,'')
		         FROM accounts WHERE is_active=true ORDER BY code`
	} else {
		query = `SELECT code, name, COALESCE(name2,''), type, COALESCE(parent_code,''), is_active, is_foreign,
		                COALESCE(detail_by,''), is_parent, status, COALESCE(freeze_reason,''), arrears_days, COALESCE(note,'')
		         FROM accounts ORDER BY code`
	}
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAccounts(rows)
}

func (r *PGAccountRepo) Update(ctx context.Context, a *domain.Account) error {
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

func (r *PGAccountRepo) GetChildren(ctx context.Context, parentCode string) ([]domain.Account, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT code, name, COALESCE(name2,''), type, COALESCE(parent_code,''), is_active, is_foreign,
		        COALESCE(detail_by,''), is_parent, status, COALESCE(freeze_reason,''), arrears_days, COALESCE(note,'')
		 FROM accounts WHERE parent_code=$1 ORDER BY code`, parentCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAccounts(rows)
}

// ─── Journal Entry ───────────────────────────────────────────────────────────

type PGJournalRepo struct {
	pool *pgxpool.Pool
}

func NewPGJournalRepo(pool *pgxpool.Pool) *PGJournalRepo {
	return &PGJournalRepo{pool: pool}
}

func (r *PGJournalRepo) Create(ctx context.Context, e *domain.JournalEntry) error {
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

func (r *PGJournalRepo) GetByID(ctx context.Context, id string) (*domain.JournalEntry, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, COALESCE(company_id,'00000000-0000-0000-0000-000000000000'), entry_number, voucher_type, entry_date, accounting_date,
		        period_id, description, status, currency_code, exchange_rate, created_by,
		        COALESCE(reviewed_by,''), COALESCE(approved_by,''), created_at, posted_at, approved_at
		 FROM journal_entries WHERE id=$1`, id)
	entry, err := scanJournalEntry(row)
	if err != nil {
		return nil, err
	}
	lines, err := r.GetLinesByEntryID(ctx, id)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines
	return entry, nil
}

func (r *PGJournalRepo) GetByPeriod(ctx context.Context, periodID string) ([]domain.JournalEntry, error) {
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

func (r *PGJournalRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]domain.JournalEntry, error) {
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

func (r *PGJournalRepo) GetByStatus(ctx context.Context, status domain.JournalEntryStatus) ([]domain.JournalEntry, error) {
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

func (r *PGJournalRepo) GetByVoucherType(ctx context.Context, vt domain.VoucherType) ([]domain.JournalEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(company_id,'00000000-0000-0000-0000-000000000000'), entry_number, voucher_type, entry_date, accounting_date,
		        period_id, description, status, currency_code, exchange_rate, created_by,
		        COALESCE(reviewed_by,''), COALESCE(approved_by,''), created_at, posted_at, approved_at
		 FROM journal_entries WHERE voucher_type=$1 ORDER BY entry_date`, string(vt))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJournalEntries(rows)
}

func (r *PGJournalRepo) UpdateStatus(ctx context.Context, id string, status domain.JournalEntryStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE journal_entries SET status=$1 WHERE id=$2`, string(status), id)
	return err
}

func (r *PGJournalRepo) Update(ctx context.Context, e *domain.JournalEntry) error {
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

func (r *PGJournalRepo) GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error) {
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

func (r *PGJournalRepo) GetBalance(ctx context.Context, accountCode, periodID string) (*domain.AccountBalance, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(debit_amount),0) AS period_debit,
		        COALESCE(SUM(credit_amount),0) AS period_credit
		 FROM journal_lines l
		 JOIN journal_entries e ON l.entry_id = e.id
		 WHERE l.account_code=$1 AND e.period_id=$2 AND e.status='POSTED'`,
		accountCode, periodID)

	b := &domain.AccountBalance{AccountCode: accountCode, PeriodID: periodID}
	err := row.Scan(&b.PeriodDebit, &b.PeriodCredit)
	return b, err
}

func (r *PGJournalRepo) GetTrialBalance(ctx context.Context, periodID string) ([]domain.AccountBalance, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT l.account_code, $1 AS period_id, a.type,
		        COALESCE(SUM(l.debit_amount),0) AS period_debit,
		        COALESCE(SUM(l.credit_amount),0) AS period_credit
		 FROM journal_lines l
		 JOIN journal_entries e ON l.entry_id = e.id
		 JOIN accounts a ON l.account_code = a.code
		 WHERE e.period_id=$1 AND e.status='POSTED'
		 GROUP BY l.account_code, a.type
		 ORDER BY l.account_code`, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var balances []domain.AccountBalance
	for rows.Next() {
		var b domain.AccountBalance
		if err := rows.Scan(&b.AccountCode, &b.PeriodID, &b.AccountType, &b.PeriodDebit, &b.PeriodCredit); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

func (r *PGJournalRepo) GetFinancialStatement(ctx context.Context, periodID string, accountTypes []domain.AccountType) ([]domain.AccountBalance, error) {
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
	var balances []domain.AccountBalance
	for rows.Next() {
		var b domain.AccountBalance
		if err := rows.Scan(&b.AccountCode, &b.PeriodID, &b.AccountType, &b.PeriodDebit, &b.PeriodCredit); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

func (r *PGJournalRepo) GetAccountUsage(ctx context.Context, accountCode string) (*domain.AccountUsage, error) {
	var u domain.AccountUsage
	u.AccountCode = accountCode
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(COUNT(*),0), COALESCE(SUM(l.debit_amount),0), COALESCE(SUM(l.credit_amount),0),
		        COALESCE(MAX(e.entry_date)::TEXT,'')
		 FROM journal_lines l
		 JOIN journal_entries e ON l.entry_id = e.id
		 WHERE l.account_code=$1 AND e.status='POSTED'`, accountCode).Scan(
		&u.EntryCount, &u.TotalDebit, &u.TotalCredit, &u.LastUsedDate)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PGJournalRepo) GetPostedEntriesByAccount(ctx context.Context, periodID, accountCode string) ([]domain.JournalEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT e.id, e.company_id, e.entry_number, e.voucher_type, e.entry_date, e.accounting_date,
		        e.period_id, e.description, e.status, e.currency_code, e.exchange_rate,
		        e.created_by, COALESCE(e.reviewed_by,''), COALESCE(e.approved_by,''),
		        e.created_at, e.posted_at, e.approved_at
		 FROM journal_entries e
		 JOIN journal_lines l ON l.entry_id = e.id
		 WHERE e.period_id=$1 AND e.status='POSTED' AND l.account_code=$2
		 ORDER BY e.entry_date`, periodID, accountCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJournalEntries(rows)
}

// ─── Period ──────────────────────────────────────────────────────────────────

type PGPeriodRepo struct {
	pool *pgxpool.Pool
}

func NewPGPeriodRepo(pool *pgxpool.Pool) *PGPeriodRepo {
	return &PGPeriodRepo{pool: pool}
}

func (r *PGPeriodRepo) Create(ctx context.Context, p *domain.Period) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO periods (id, year, month, start_date, end_date, status)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		p.ID, p.Year, p.Month, p.StartDate, p.EndDate, string(p.Status))
	return err
}

func (r *PGPeriodRepo) GetByID(ctx context.Context, id string) (*domain.Period, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, year, month, start_date, end_date, status FROM periods WHERE id=$1`, id)
	p := &domain.Period{}
	err := row.Scan(&p.ID, &p.Year, &p.Month, &p.StartDate, &p.EndDate, &p.Status)
	if err != nil {
		return nil, domain.ErrPeriodNotFound
	}
	return p, nil
}

func (r *PGPeriodRepo) GetByYearMonth(ctx context.Context, year, month int) (*domain.Period, error) {
	return r.GetByID(ctx, fmt.Sprintf("P-%d-%02d", year, month))
}

func (r *PGPeriodRepo) GetAll(ctx context.Context) ([]domain.Period, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, year, month, start_date, end_date, status FROM periods ORDER BY year, month`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var periods []domain.Period
	for rows.Next() {
		var p domain.Period
		if err := rows.Scan(&p.ID, &p.Year, &p.Month, &p.StartDate, &p.EndDate, &p.Status); err != nil {
			return nil, err
		}
		periods = append(periods, p)
	}
	return periods, nil
}

func (r *PGPeriodRepo) UpdateStatus(ctx context.Context, id string, status domain.PeriodStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE periods SET status=$1, updated_at=NOW() WHERE id=$2`, string(status), id)
	return err
}

func (r *PGPeriodRepo) GetOpenPeriod(ctx context.Context) (*domain.Period, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, year, month, start_date, end_date, status FROM periods WHERE status='OPEN' LIMIT 1`)
	p := &domain.Period{}
	err := row.Scan(&p.ID, &p.Year, &p.Month, &p.StartDate, &p.EndDate, &p.Status)
	if err != nil {
		return nil, domain.ErrPeriodNotFound
	}
	return p, nil
}

// ─── User ────────────────────────────────────────────────────────────────────

type PGUserRepo struct {
	pool *pgxpool.Pool
}

func NewPGUserRepo(pool *pgxpool.Pool) *PGUserRepo {
	return &PGUserRepo{pool: pool}
}

func (r *PGUserRepo) Create(ctx context.Context, u *domain.User) error {
	if u.ID == "" {
		u.ID = fmt.Sprintf("U-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, full_name, email, role, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())`,
		u.ID, u.Username, u.PasswordHash, u.FullName, nullStr(u.Email), string(u.Role), u.IsActive)
	return err
}

func (r *PGUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, full_name, COALESCE(email,''), role, is_active, created_at, updated_at
		 FROM users WHERE id=$1`, id)
	return scanUser(row)
}

func (r *PGUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, full_name, COALESCE(email,''), role, is_active, created_at, updated_at
		 FROM users WHERE username=$1`, username)
	u, err := scanUser(row)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (r *PGUserRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, username, password_hash, full_name, COALESCE(email,''), role, is_active, created_at, updated_at
		 FROM users ORDER BY username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, nil
}

func (r *PGUserRepo) Update(ctx context.Context, u *domain.User) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET full_name=$1, email=$2, role=$3, is_active=$4, updated_at=NOW() WHERE id=$5`,
		u.FullName, nullStr(u.Email), string(u.Role), u.IsActive, u.ID)
	return err
}

func (r *PGUserRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

// ─── Audit Log ───────────────────────────────────────────────────────────────

type PGAuditLogRepo struct {
	pool *pgxpool.Pool
}

func NewPGAuditLogRepo(pool *pgxpool.Pool) *PGAuditLogRepo {
	return &PGAuditLogRepo{pool: pool}
}

func (r *PGAuditLogRepo) Create(ctx context.Context, e *domain.AuditEntry) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_log (user_id, username, ip_address, action, entity_type, entity_id, old_value, new_value)
		 VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8::jsonb)`,
		nullStr(e.UserID), e.Username, nullStr(e.IPAddress), string(e.Action),
		e.EntityType, nullStr(e.EntityID), nullStr(e.OldValue), nullStr(e.NewValue))
	return err
}

func (r *PGAuditLogRepo) GetByEntity(ctx context.Context, entityType, entityID string) ([]domain.AuditEntry, error) {
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

func (r *PGAuditLogRepo) GetByUser(ctx context.Context, userID string) ([]domain.AuditEntry, error) {
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

func (r *PGAuditLogRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]domain.AuditEntry, error) {
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

func (r *PGAuditLogRepo) GetAll(ctx context.Context, limit int) ([]domain.AuditEntry, error) {
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

// ─── Exchange Rate ───────────────────────────────────────────────────────────

type PGExchangeRateRepo struct {
	pool *pgxpool.Pool
}

func NewPGExchangeRateRepo(pool *pgxpool.Pool) *PGExchangeRateRepo {
	return &PGExchangeRateRepo{pool: pool}
}

func (r *PGExchangeRateRepo) Create(ctx context.Context, rate *domain.ExchangeRate) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO exchange_rates (currency_code, rate_date, buy_rate, sell_rate, average_rate, source)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		rate.CurrencyCode, rate.RateDate, nullF64(rate.BuyRate), nullF64(rate.SellRate),
		rate.AverageRate, nullStr(rate.Source))
	return err
}

func (r *PGExchangeRateRepo) GetByCurrencyAndDate(ctx context.Context, currencyCode string, rateDate time.Time) (*domain.ExchangeRate, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, currency_code, rate_date, COALESCE(buy_rate,0), COALESCE(sell_rate,0), average_rate,
		        COALESCE(source,''), created_at
		 FROM exchange_rates WHERE currency_code=$1 AND rate_date=$2`, currencyCode, rateDate)
	xr := &domain.ExchangeRate{}
	err := row.Scan(&xr.ID, &xr.CurrencyCode, &xr.RateDate, &xr.BuyRate, &xr.SellRate, &xr.AverageRate, &xr.Source, &xr.CreatedAt)
	if err != nil {
		return nil, domain.ErrRateNotFound
	}
	return xr, nil
}

func (r *PGExchangeRateRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]domain.ExchangeRate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, currency_code, rate_date, COALESCE(buy_rate,0), COALESCE(sell_rate,0), average_rate,
		        COALESCE(source,''), created_at
		 FROM exchange_rates WHERE rate_date>=$1 AND rate_date<=$2 ORDER BY rate_date, currency_code`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []domain.ExchangeRate
	for rows.Next() {
		var xr domain.ExchangeRate
		if err := rows.Scan(&xr.ID, &xr.CurrencyCode, &xr.RateDate, &xr.BuyRate, &xr.SellRate, &xr.AverageRate, &xr.Source, &xr.CreatedAt); err != nil {
			return nil, err
		}
		rates = append(rates, xr)
	}
	return rates, nil
}

func (r *PGExchangeRateRepo) GetAll(ctx context.Context) ([]domain.ExchangeRate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, currency_code, rate_date, COALESCE(buy_rate,0), COALESCE(sell_rate,0), average_rate,
		        COALESCE(source,''), created_at
		 FROM exchange_rates ORDER BY rate_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []domain.ExchangeRate
	for rows.Next() {
		var xr domain.ExchangeRate
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

// ─── Closing Template ────────────────────────────────────────────────────────

type PGClosingTemplateRepo struct {
	pool *pgxpool.Pool
}

func NewPGClosingTemplateRepo(pool *pgxpool.Pool) *PGClosingTemplateRepo {
	return &PGClosingTemplateRepo{pool: pool}
}

func (r *PGClosingTemplateRepo) Create(ctx context.Context, t *domain.ClosingTemplate) error {
	if t.ID == "" {
		t.ID = fmt.Sprintf("CT-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO closing_templates (id, name, description, sequence_order, is_active) VALUES ($1,$2,$3,$4,$5)`,
		t.ID, t.Name, nullStr(t.Description), t.SequenceOrder, t.IsActive)
	return err
}

func (r *PGClosingTemplateRepo) GetByID(ctx context.Context, id string) (*domain.ClosingTemplate, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, COALESCE(description,''), sequence_order, is_active, created_at
		 FROM closing_templates WHERE id=$1`, id)
	ct := &domain.ClosingTemplate{}
	err := row.Scan(&ct.ID, &ct.Name, &ct.Description, &ct.SequenceOrder, &ct.IsActive, &ct.CreatedAt)
	if err != nil {
		return nil, domain.ErrTemplateNotFound
	}
	return ct, nil
}

func (r *PGClosingTemplateRepo) GetAll(ctx context.Context) ([]domain.ClosingTemplate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, COALESCE(description,''), sequence_order, is_active, created_at
		 FROM closing_templates ORDER BY sequence_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var templates []domain.ClosingTemplate
	for rows.Next() {
		var ct domain.ClosingTemplate
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.Description, &ct.SequenceOrder, &ct.IsActive, &ct.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, ct)
	}
	return templates, nil
}

func (r *PGClosingTemplateRepo) Update(ctx context.Context, t *domain.ClosingTemplate) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE closing_templates SET name=$1, description=$2, sequence_order=$3, is_active=$4 WHERE id=$5`,
		t.Name, nullStr(t.Description), t.SequenceOrder, t.IsActive, t.ID)
	return err
}

func (r *PGClosingTemplateRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM closing_templates WHERE id=$1`, id)
	return err
}

// ─── Scan helpers ────────────────────────────────────────────────────────────

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanAccount(row scannable) (*domain.Account, error) {
	a := &domain.Account{}
	var detailBy, parentCode, name2, freezeReason, note string
	err := row.Scan(&a.Code, &a.Name, &name2, &a.Type, &parentCode, &a.IsActive, &a.IsForeign,
		&detailBy, &a.IsParent, &a.Status, &freezeReason, &a.ArrearsDays, &note)
	if err != nil {
		return nil, err
	}
	a.Name2 = name2
	a.ParentCode = parentCode
	a.DetailBy = domain.DetailBy(detailBy)
	a.FreezeReason = freezeReason
	a.Note = note
	return a, nil
}

func scanAccounts(rows pgx.Rows) ([]domain.Account, error) {
	defer rows.Close()
	var accounts []domain.Account
	for rows.Next() {
		a, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, *a)
	}
	return accounts, nil
}

func scanJournalEntry(row scannable) (*domain.JournalEntry, error) {
	e := &domain.JournalEntry{}
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

func scanJournalEntries(rows pgx.Rows) ([]domain.JournalEntry, error) {
	defer rows.Close()
	var entries []domain.JournalEntry
	for rows.Next() {
		e, err := scanJournalEntry(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, *e)
	}
	return entries, nil
}

func scanJournalLines(rows pgx.Rows) ([]domain.JournalLine, error) {
	defer rows.Close()
	var lines []domain.JournalLine
	for rows.Next() {
		var l domain.JournalLine
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

func scanUser(row scannable) (*domain.User, error) {
	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func scanAuditEntries(rows pgx.Rows) ([]domain.AuditEntry, error) {
	defer rows.Close()
	var entries []domain.AuditEntry
	for rows.Next() {
		var e domain.AuditEntry
		err := rows.Scan(&e.ID, &e.UserID, &e.Username, &e.IPAddress, &e.Action, &e.EntityType,
			&e.EntityID, &e.OldValue, &e.NewValue, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

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

func nullTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}

func nullTimePtr(t *time.Time) interface{} {
	if t == nil || t.IsZero() {
		return nil
	}
	return *t
}

func nullFloat(f float64) interface{} {
	if f == 0 {
		return nil
	}
	return f
}

// ─── Tax (PG) ──────────────────────────────────────────────────────────────

type PGTaxRepo struct {
	pool *pgxpool.Pool
}

func NewPGTaxRepo(pool *pgxpool.Pool) *PGTaxRepo {
	return &PGTaxRepo{pool: pool}
}

func (r *PGTaxRepo) CreateDeclaration(ctx context.Context, d *domain.TaxDeclaration) error {
	if d.ID == "" {
		d.ID = fmt.Sprintf("TD-%d", time.Now().UnixNano())
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`INSERT INTO tax_declarations (id, company_id, declaration_type, period_year, period_number, period_type,
		 status, adjustment_type, version, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,1,$9,NOW(),NOW())`,
		d.ID, d.CompanyID, string(d.DeclarationType),
		d.TaxPeriod.PeriodYear, d.TaxPeriod.PeriodNumber, string(d.TaxPeriod.PeriodType),
		string(d.Status), string(d.AdjustmentType), d.CreatedBy)
	if err != nil { return err }
	for i, line := range d.Lines {
		lineID := fmt.Sprintf("%s-L%d", d.ID, i)
		_, err = tx.Exec(ctx,
			`INSERT INTO tax_declaration_lines (id, declaration_id, line_code, line_name, amount, source_type, sort_order)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			lineID, d.ID, line.LineCode, line.LineName, line.Amount, string(line.SourceType), line.SortOrder)
		if err != nil { return err }
	}
	return tx.Commit(ctx)
}

func (r *PGTaxRepo) GetDeclarationByID(ctx context.Context, id string) (*domain.TaxDeclaration, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, declaration_type, period_year, period_number, period_type,
		        status, COALESCE(submitted_at::TEXT,''), COALESCE(submitted_by,''),
		        COALESCE(acknowledged_at::TEXT,''), COALESCE(acknowledgement_ref,''),
		        COALESCE(gdt_response_xml,''), COALESCE(declaration_xml,''),
		        COALESCE(previous_declaration_id,''), adjustment_type, version,
		        created_at::TEXT, COALESCE(created_by,''), updated_at::TEXT
		 FROM tax_declarations WHERE id=$1`, id)
	d := &domain.TaxDeclaration{}
	err := row.Scan(&d.ID, &d.CompanyID, &d.DeclarationType,
		&d.TaxPeriod.PeriodYear, &d.TaxPeriod.PeriodNumber, &d.TaxPeriod.PeriodType,
		&d.Status, &d.SubmittedAt, &d.SubmittedBy,
		&d.AcknowledgedAt, &d.AcknowledgementRef,
		&d.GDTResponseXML, &d.DeclarationXML,
		&d.PreviousDeclID, &d.AdjustmentType, &d.Version,
		&d.CreatedAt, &d.CreatedBy, &d.UpdatedAt)
	if err != nil { return nil, domain.ErrDeclarationNotFound }
	rows, err := r.pool.Query(ctx,
		`SELECT id, declaration_id, line_code, COALESCE(line_name,''), amount, source_type, sort_order
		 FROM tax_declaration_lines WHERE declaration_id=$1 ORDER BY sort_order`, id)
	if err != nil { return nil, err }
	defer rows.Close()
	for rows.Next() {
		var line domain.TaxDeclarationLine
		if err := rows.Scan(&line.ID, &line.DeclarationID, &line.LineCode, &line.LineName,
			&line.Amount, &line.SourceType, &line.SortOrder); err != nil { return nil, err }
		d.Lines = append(d.Lines, line)
	}
	return d, nil
}

func (r *PGTaxRepo) GetDeclarations(ctx context.Context, filter domain.TaxDeclarationFilter) ([]domain.TaxDeclaration, error) {
	query := `SELECT id, company_id, declaration_type, period_year, period_number, period_type,
	                 status, COALESCE(submitted_at::TEXT,''), COALESCE(submitted_by,''),
	                 COALESCE(acknowledged_at::TEXT,''), COALESCE(acknowledgement_ref,''),
	                 COALESCE(gdt_response_xml,''), COALESCE(declaration_xml,''),
	                 COALESCE(previous_declaration_id,''), adjustment_type, version,
	                 created_at::TEXT, COALESCE(created_by,''), updated_at::TEXT
	          FROM tax_declarations WHERE 1=1`
	var args []interface{}
	n := 1
	if filter.CompanyID != "" { query += fmt.Sprintf(" AND company_id=$%d", n); args = append(args, filter.CompanyID); n++ }
	if filter.DeclarationType != "" { query += fmt.Sprintf(" AND declaration_type=$%d", n); args = append(args, string(filter.DeclarationType)); n++ }
	if filter.Status != "" { query += fmt.Sprintf(" AND status=$%d", n); args = append(args, string(filter.Status)); n++ }
	if filter.PeriodYear != 0 { query += fmt.Sprintf(" AND period_year=$%d", n); args = append(args, filter.PeriodYear); n++ }
	if filter.PeriodNumber != 0 { query += fmt.Sprintf(" AND period_number=$%d", n); args = append(args, filter.PeriodNumber); n++ }
	query += " ORDER BY created_at DESC"
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TaxDeclaration
	for rows.Next() {
		var d domain.TaxDeclaration
		if err := rows.Scan(&d.ID, &d.CompanyID, &d.DeclarationType,
			&d.TaxPeriod.PeriodYear, &d.TaxPeriod.PeriodNumber, &d.TaxPeriod.PeriodType,
			&d.Status, &d.SubmittedAt, &d.SubmittedBy,
			&d.AcknowledgedAt, &d.AcknowledgementRef,
			&d.GDTResponseXML, &d.DeclarationXML,
			&d.PreviousDeclID, &d.AdjustmentType, &d.Version,
			&d.CreatedAt, &d.CreatedBy, &d.UpdatedAt); err != nil { return nil, err }
		out = append(out, d)
	}
	return out, nil
}

func (r *PGTaxRepo) UpdateDeclaration(ctx context.Context, d *domain.TaxDeclaration) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tax_declarations SET status=$1, declaration_xml=$2, adjustment_type=$3, version=version+1, updated_at=NOW()
		 WHERE id=$4`,
		string(d.Status), nullStr(d.DeclarationXML), string(d.AdjustmentType), d.ID)
	return err
}

func (r *PGTaxRepo) UpdateDeclarationStatus(ctx context.Context, id string, status domain.DeclarationStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tax_declarations SET status=$1, updated_at=NOW() WHERE id=$2`,
		string(status), id)
	return err
}

func (r *PGTaxRepo) CreateRate(ctx context.Context, rate *domain.TaxRate) error {
	if rate.ID == "" {
		rate.ID = fmt.Sprintf("TR-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tax_rates (id, tax_type, rate_code, rate_name, rate_type, rate_value, effective_from, effective_to, is_active, legal_ref)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		rate.ID, string(rate.TaxType), rate.RateCode, rate.RateName, string(rate.RateType),
		nullF64(rate.RateValue), rate.EffectiveFrom, nullStr(rate.EffectiveTo), rate.IsActive, nullStr(rate.LegalRef))
	return err
}

func (r *PGTaxRepo) GetRateByID(ctx context.Context, id string) (*domain.TaxRate, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tax_type, rate_code, rate_name, rate_type, COALESCE(rate_value,0),
		        effective_from, COALESCE(effective_to,''), is_active, COALESCE(legal_ref,''), created_at
		 FROM tax_rates WHERE id=$1`, id)
	rate := &domain.TaxRate{}
	err := row.Scan(&rate.ID, &rate.TaxType, &rate.RateCode, &rate.RateName, &rate.RateType,
		&rate.RateValue, &rate.EffectiveFrom, &rate.EffectiveTo, &rate.IsActive, &rate.LegalRef, &rate.CreatedAt)
	if err != nil { return nil, domain.ErrTaxRateNotFound }
	return rate, nil
}

func (r *PGTaxRepo) GetRates(ctx context.Context, filter domain.TaxRateFilter) ([]domain.TaxRate, error) {
	query := `SELECT id, tax_type, rate_code, rate_name, rate_type, COALESCE(rate_value,0),
	                 effective_from, COALESCE(effective_to,''), is_active, COALESCE(legal_ref,''), created_at
	          FROM tax_rates WHERE 1=1`
	var args []interface{}
	n := 1
	if filter.TaxType != "" { query += fmt.Sprintf(" AND tax_type=$%d", n); args = append(args, string(filter.TaxType)); n++ }
	if filter.IsActive != nil { query += fmt.Sprintf(" AND is_active=$%d", n); args = append(args, *filter.IsActive); n++ }
	query += " ORDER BY rate_code"
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TaxRate
	for rows.Next() {
		var rate domain.TaxRate
		if err := rows.Scan(&rate.ID, &rate.TaxType, &rate.RateCode, &rate.RateName, &rate.RateType,
			&rate.RateValue, &rate.EffectiveFrom, &rate.EffectiveTo, &rate.IsActive, &rate.LegalRef, &rate.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rate)
	}
	return out, nil
}

func (r *PGTaxRepo) UpdateRate(ctx context.Context, rate *domain.TaxRate) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tax_rates SET rate_name=$1, rate_type=$2, rate_value=$3, is_active=$4, effective_to=$5 WHERE id=$6`,
		rate.RateName, string(rate.RateType), nullF64(rate.RateValue), rate.IsActive, nullStr(rate.EffectiveTo), rate.ID)
	return err
}

func (r *PGTaxRepo) CreatePayment(ctx context.Context, p *domain.TaxPayment) error {
	if p.ID == "" {
		p.ID = fmt.Sprintf("TP-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tax_payments (id, company_id, declaration_id, tax_type, period_year, period_number,
		 declared_amount, paid_amount, due_date, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW())`,
		p.ID, p.CompanyID, nullStr(p.DeclarationID), string(p.TaxType),
		p.PeriodYear, p.PeriodNumber, p.DeclaredAmount, p.PaidAmount, p.DueDate, string(p.Status))
	return err
}

func (r *PGTaxRepo) GetPaymentByID(ctx context.Context, id string) (*domain.TaxPayment, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, COALESCE(declaration_id,''), tax_type, period_year, period_number,
		        declared_amount, paid_amount, COALESCE(payment_date::TEXT,''), due_date,
		        COALESCE(payment_ref,''), COALESCE(payment_method,''), status, COALESCE(late_days,0),
		        COALESCE(late_interest,0), COALESCE(notes,''), created_at::TEXT
		 FROM tax_payments WHERE id=$1`, id)
	p := &domain.TaxPayment{}
	err := row.Scan(&p.ID, &p.CompanyID, &p.DeclarationID, &p.TaxType,
		&p.PeriodYear, &p.PeriodNumber, &p.DeclaredAmount, &p.PaidAmount,
		&p.PaymentDate, &p.DueDate, &p.PaymentRef, &p.PaymentMethod, &p.Status,
		&p.LateDays, &p.LateInterest, &p.Notes, &p.CreatedAt)
	if err != nil { return nil, domain.ErrTaxPaymentNotFound }
	return p, nil
}

func (r *PGTaxRepo) GetPayments(ctx context.Context, filter domain.PaymentFilter) ([]domain.TaxPayment, error) {
	query := `SELECT id, company_id, COALESCE(declaration_id,''), tax_type, period_year, period_number,
	                 declared_amount, paid_amount, COALESCE(payment_date::TEXT,''), due_date,
	                 COALESCE(payment_ref,''), COALESCE(payment_method,''), status, COALESCE(late_days,0),
	                 COALESCE(late_interest,0), COALESCE(notes,''), created_at::TEXT
	          FROM tax_payments WHERE 1=1`
	var args []interface{}
	n := 1
	if filter.CompanyID != "" { query += fmt.Sprintf(" AND company_id=$%d", n); args = append(args, filter.CompanyID); n++ }
	if filter.TaxType != "" { query += fmt.Sprintf(" AND tax_type=$%d", n); args = append(args, string(filter.TaxType)); n++ }
	if filter.Status != "" { query += fmt.Sprintf(" AND status=$%d", n); args = append(args, string(filter.Status)); n++ }
	if filter.PeriodYear != 0 { query += fmt.Sprintf(" AND period_year=$%d", n); args = append(args, filter.PeriodYear); n++ }
	if filter.PeriodNumber != 0 { query += fmt.Sprintf(" AND period_number=$%d", n); args = append(args, filter.PeriodNumber); n++ }
	query += " ORDER BY due_date"
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TaxPayment
	for rows.Next() {
		var p domain.TaxPayment
		if err := rows.Scan(&p.ID, &p.CompanyID, &p.DeclarationID, &p.TaxType,
			&p.PeriodYear, &p.PeriodNumber, &p.DeclaredAmount, &p.PaidAmount,
			&p.PaymentDate, &p.DueDate, &p.PaymentRef, &p.PaymentMethod, &p.Status,
			&p.LateDays, &p.LateInterest, &p.Notes, &p.CreatedAt); err != nil { return nil, err }
		out = append(out, p)
	}
	return out, nil
}

func (r *PGTaxRepo) UpdatePayment(ctx context.Context, p *domain.TaxPayment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tax_payments SET paid_amount=$1, payment_date=$2, payment_ref=$3, status=$4, late_days=$5, late_interest=$6, notes=$7 WHERE id=$8`,
		p.PaidAmount, nullStr(p.PaymentDate), nullStr(p.PaymentRef), string(p.Status), p.LateDays, p.LateInterest, nullStr(p.Notes), p.ID)
	return err
}

func (r *PGTaxRepo) CreateEInvoice(ctx context.Context, inv *domain.EInvoice) error {
	if inv.ID == "" {
		inv.ID = fmt.Sprintf("INV-%d", time.Now().UnixNano())
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`INSERT INTO e_invoices (id, company_id, pattern, serial, invoice_type, buyer_name, buyer_tax_code,
		 buyer_address, buyer_email, currency_code, exchange_rate, subtotal, vat_amount, grand_total,
		 issue_date, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,NOW())`,
		inv.ID, inv.CompanyID, inv.Pattern, inv.Serial, string(inv.InvoiceType),
		inv.BuyerName, nullStr(inv.BuyerTaxCode), nullStr(inv.BuyerAddress), nullStr(inv.BuyerEmail),
		inv.CurrencyCode, inv.ExchangeRate, inv.Subtotal, inv.VATAmount, inv.GrandTotal,
		inv.IssueDate, string(inv.Status))
	if err != nil { return err }
	for i, line := range inv.Lines {
		lineID := fmt.Sprintf("%s-L%d", inv.ID, i)
		_, err = tx.Exec(ctx,
			`INSERT INTO e_invoice_lines (id, e_invoice_id, line_number, description, unit, quantity, unit_price, line_total, vat_rate, vat_amount)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			lineID, inv.ID, line.LineNumber, line.Description, nullStr(line.Unit),
			line.Quantity, line.UnitPrice, line.LineTotal, line.VATRate, line.VATAmount)
		if err != nil { return err }
	}
	return tx.Commit(ctx)
}

func (r *PGTaxRepo) GetEInvoiceByID(ctx context.Context, id string) (*domain.EInvoice, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, pattern, serial, COALESCE(invoice_number,0), invoice_type, COALESCE(gdt_transaction_id,''),
		        buyer_name, COALESCE(buyer_tax_code,''), COALESCE(buyer_address,''), COALESCE(buyer_email,''),
		        currency_code, exchange_rate, subtotal, vat_amount, grand_total,
		        COALESCE(xml_body,''), COALESCE(signed_xml,''), issue_date, COALESCE(signing_date::TEXT,''),
		        COALESCE(digital_signature_id,''), COALESCE(journal_entry_id,''), status,
		        COALESCE(cancelled_at::TEXT,''), COALESCE(cancel_reason,''), COALESCE(original_invoice_id,''),
		        COALESCE(gdt_response,''), created_at::TEXT
		 FROM e_invoices WHERE id=$1`, id)
	inv := &domain.EInvoice{}
	err := row.Scan(&inv.ID, &inv.CompanyID, &inv.Pattern, &inv.Serial, &inv.InvoiceNumber,
		&inv.InvoiceType, &inv.GDTTransactionID,
		&inv.BuyerName, &inv.BuyerTaxCode, &inv.BuyerAddress, &inv.BuyerEmail,
		&inv.CurrencyCode, &inv.ExchangeRate, &inv.Subtotal, &inv.VATAmount, &inv.GrandTotal,
		&inv.XMLBody, &inv.SignedXML, &inv.IssueDate, &inv.SigningDate,
		&inv.DigitalSignatureID, &inv.JournalEntryID, &inv.Status,
		&inv.CancelledAt, &inv.CancelReason, &inv.OriginalInvoiceID,
		&inv.GDTResponse, &inv.CreatedAt)
	if err != nil { return nil, domain.ErrInvoiceNotFound }
	return inv, nil
}

func (r *PGTaxRepo) GetEInvoices(ctx context.Context, filter domain.EInvoiceFilter) ([]domain.EInvoice, error) {
	query := `SELECT id, company_id, pattern, serial, COALESCE(invoice_number,0), invoice_type, COALESCE(gdt_transaction_id,''),
	                 buyer_name, COALESCE(buyer_tax_code,''), COALESCE(buyer_address,''), COALESCE(buyer_email,''),
	                 currency_code, exchange_rate, subtotal, vat_amount, grand_total,
	                 COALESCE(xml_body,''), COALESCE(signed_xml,''), issue_date, COALESCE(signing_date::TEXT,''),
	                 COALESCE(digital_signature_id,''), COALESCE(journal_entry_id,''), status,
	                 COALESCE(cancelled_at::TEXT,''), COALESCE(cancel_reason,''), COALESCE(original_invoice_id,''),
	                 COALESCE(gdt_response,''), created_at::TEXT
	          FROM e_invoices WHERE 1=1`
	var args []interface{}
	n := 1
	if filter.CompanyID != "" { query += fmt.Sprintf(" AND company_id=$%d", n); args = append(args, filter.CompanyID); n++ }
	if filter.Status != "" { query += fmt.Sprintf(" AND status=$%d", n); args = append(args, string(filter.Status)); n++ }
	if filter.FromDate != "" { query += fmt.Sprintf(" AND issue_date>=$%d", n); args = append(args, filter.FromDate); n++ }
	if filter.ToDate != "" { query += fmt.Sprintf(" AND issue_date<=$%d", n); args = append(args, filter.ToDate); n++ }
	query += " ORDER BY created_at DESC"
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.EInvoice
	for rows.Next() {
		var inv domain.EInvoice
		if err := rows.Scan(&inv.ID, &inv.CompanyID, &inv.Pattern, &inv.Serial, &inv.InvoiceNumber,
			&inv.InvoiceType, &inv.GDTTransactionID,
			&inv.BuyerName, &inv.BuyerTaxCode, &inv.BuyerAddress, &inv.BuyerEmail,
			&inv.CurrencyCode, &inv.ExchangeRate, &inv.Subtotal, &inv.VATAmount, &inv.GrandTotal,
			&inv.XMLBody, &inv.SignedXML, &inv.IssueDate, &inv.SigningDate,
			&inv.DigitalSignatureID, &inv.JournalEntryID, &inv.Status,
			&inv.CancelledAt, &inv.CancelReason, &inv.OriginalInvoiceID,
			&inv.GDTResponse, &inv.CreatedAt); err != nil { return nil, err }
		out = append(out, inv)
	}
	return out, nil
}

func (r *PGTaxRepo) UpdateEInvoice(ctx context.Context, inv *domain.EInvoice) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE e_invoices SET buyer_name=$1, buyer_address=$2, buyer_email=$3, gdt_response=$4 WHERE id=$5`,
		inv.BuyerName, nullStr(inv.BuyerAddress), nullStr(inv.BuyerEmail), nullStr(inv.GDTResponse), inv.ID)
	return err
}

func (r *PGTaxRepo) UpdateEInvoiceStatus(ctx context.Context, id string, status domain.EInvLifecycleStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE e_invoices SET status=$1 WHERE id=$2`, string(status), id)
	return err
}

func (r *PGTaxRepo) CreateCalendarEntry(ctx context.Context, c *domain.TaxCalendar) error {
	if c.ID == "" {
		c.ID = fmt.Sprintf("CAL-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tax_calendar (id, company_id, tax_type, period_type, period_year, period_number,
		 start_date, end_date, declaration_due, payment_due, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW())`,
		c.ID, c.CompanyID, string(c.TaxType), string(c.PeriodType),
		c.PeriodYear, c.PeriodNumber, c.StartDate, c.EndDate, c.DeclarationDue, nullStr(c.PaymentDue), string(c.Status))
	return err
}

func (r *PGTaxRepo) GetCalendarEntryByID(ctx context.Context, id string) (*domain.TaxCalendar, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, tax_type, period_type, period_year, period_number,
		        start_date, end_date, declaration_due, COALESCE(payment_due,''), status, created_at::TEXT
		 FROM tax_calendar WHERE id=$1`, id)
	c := &domain.TaxCalendar{}
	err := row.Scan(&c.ID, &c.CompanyID, &c.TaxType, &c.PeriodType,
		&c.PeriodYear, &c.PeriodNumber, &c.StartDate, &c.EndDate, &c.DeclarationDue, &c.PaymentDue, &c.Status, &c.CreatedAt)
	if err != nil { return nil, domain.ErrCalendarNotFound }
	return c, nil
}

func (r *PGTaxRepo) GetCalendarByPeriod(ctx context.Context, companyID string, periodYear, periodNumber int) ([]domain.TaxCalendar, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, tax_type, period_type, period_year, period_number,
		        start_date, end_date, declaration_due, COALESCE(payment_due,''), status, created_at::TEXT
		 FROM tax_calendar WHERE company_id=$1 AND period_year=$2 AND period_number=$3 ORDER BY tax_type`,
		companyID, periodYear, periodNumber)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TaxCalendar
	for rows.Next() {
		var c domain.TaxCalendar
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.TaxType, &c.PeriodType,
			&c.PeriodYear, &c.PeriodNumber, &c.StartDate, &c.EndDate, &c.DeclarationDue, &c.PaymentDue, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *PGTaxRepo) GetCalendarByCompany(ctx context.Context, companyID string) ([]domain.TaxCalendar, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, tax_type, period_type, period_year, period_number,
		        start_date, end_date, declaration_due, COALESCE(payment_due,''), status, created_at::TEXT
		 FROM tax_calendar WHERE company_id=$1 ORDER BY period_year DESC, period_number DESC`,
		companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TaxCalendar
	for rows.Next() {
		var c domain.TaxCalendar
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.TaxType, &c.PeriodType,
			&c.PeriodYear, &c.PeriodNumber, &c.StartDate, &c.EndDate, &c.DeclarationDue, &c.PaymentDue, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *PGTaxRepo) UpdateCalendarStatus(ctx context.Context, id string, status domain.CalendarStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE tax_calendar SET status=$1 WHERE id=$2`, string(status), id)
	return err
}

func (r *PGTaxRepo) CreateAlert(ctx context.Context, a *domain.TaxAlert) error {
	if a.ID == "" {
		a.ID = fmt.Sprintf("ALERT-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tax_alerts (id, company_id, calendar_id, alert_type, channel, message, sent_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW())`,
		a.ID, a.CompanyID, nullStr(a.CalendarID), string(a.AlertType), string(a.Channel), a.Message)
	return err
}

func (r *PGTaxRepo) GetAlertByID(ctx context.Context, id string) (*domain.TaxAlert, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, COALESCE(calendar_id,''), alert_type, channel, message,
		        sent_at::TEXT, COALESCE(acknowledged_at::TEXT,''), COALESCE(acknowledged_by,'')
		 FROM tax_alerts WHERE id=$1`, id)
	a := &domain.TaxAlert{}
	err := row.Scan(&a.ID, &a.CompanyID, &a.CalendarID, &a.AlertType, &a.Channel, &a.Message,
		&a.SentAt, &a.AcknowledgedAt, &a.AcknowledgedBy)
	if err != nil { return nil, domain.ErrTaxAlertNotFound }
	return a, nil
}

func (r *PGTaxRepo) GetAlerts(ctx context.Context, companyID string, limit int) ([]domain.TaxAlert, error) {
	query := `SELECT id, company_id, COALESCE(calendar_id,''), alert_type, channel, message,
	                 sent_at::TEXT, COALESCE(acknowledged_at::TEXT,''), COALESCE(acknowledged_by,'')
	          FROM tax_alerts WHERE company_id=$1 ORDER BY sent_at DESC`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TaxAlert
	for rows.Next() {
		var a domain.TaxAlert
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.CalendarID, &a.AlertType, &a.Channel, &a.Message,
			&a.SentAt, &a.AcknowledgedAt, &a.AcknowledgedBy); err != nil { return nil, err }
		out = append(out, a)
	}
	return out, nil
}

func (r *PGTaxRepo) CreateAuditCase(ctx context.Context, a *domain.TaxAuditCase) error {
	if a.ID == "" {
		a.ID = fmt.Sprintf("AUDIT-%d", time.Now().UnixNano())
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tax_audit_cases (id, company_id, audit_period_start, audit_period_end,
		 audit_decision_number, auditor_name, auditor_contact, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())`,
		a.ID, a.CompanyID, a.AuditPeriodStart, a.AuditPeriodEnd,
		a.AuditDecNumber, a.AuditorName, nullStr(a.AuditorContact), string(a.Status))
	return err
}

func (r *PGTaxRepo) GetAuditCaseByID(ctx context.Context, id string) (*domain.TaxAuditCase, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, company_id, audit_period_start, audit_period_end,
		        audit_decision_number, auditor_name, COALESCE(auditor_contact,''), status,
		        COALESCE(findings,''), COALESCE(penalty_amount,0), created_at::TEXT, COALESCE(closed_at::TEXT,'')
		 FROM tax_audit_cases WHERE id=$1`, id)
	a := &domain.TaxAuditCase{}
	err := row.Scan(&a.ID, &a.CompanyID, &a.AuditPeriodStart, &a.AuditPeriodEnd,
		&a.AuditDecNumber, &a.AuditorName, &a.AuditorContact, &a.Status,
		&a.Findings, &a.PenaltyAmount, &a.CreatedAt, &a.ClosedAt)
	if err != nil { return nil, domain.ErrAuditCaseNotFound }
	return a, nil
}

func (r *PGTaxRepo) GetAuditCases(ctx context.Context, companyID string) ([]domain.TaxAuditCase, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, audit_period_start, audit_period_end,
		        audit_decision_number, auditor_name, COALESCE(auditor_contact,''), status,
		        COALESCE(findings,''), COALESCE(penalty_amount,0), created_at::TEXT, COALESCE(closed_at::TEXT,'')
		 FROM tax_audit_cases WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TaxAuditCase
	for rows.Next() {
		var a domain.TaxAuditCase
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.AuditPeriodStart, &a.AuditPeriodEnd,
			&a.AuditDecNumber, &a.AuditorName, &a.AuditorContact, &a.Status,
			&a.Findings, &a.PenaltyAmount, &a.CreatedAt, &a.ClosedAt); err != nil { return nil, err }
		out = append(out, a)
	}
	return out, nil
}

func (r *PGTaxRepo) UpdateAuditCase(ctx context.Context, a *domain.TaxAuditCase) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tax_audit_cases SET status=$1, findings=$2, penalty_amount=$3, closed_at=$4 WHERE id=$5`,
		string(a.Status), nullStr(a.Findings), a.PenaltyAmount, nullStr(a.ClosedAt), a.ID)
	return err
}

// ─── COA: ApprovalRequest ─────────────────────────────────────────────────────

type PGApprovalRepo struct {
	pool *pgxpool.Pool
}

func NewPGApprovalRepo(pool *pgxpool.Pool) *PGApprovalRepo {
	return &PGApprovalRepo{pool: pool}
}

func (r *PGApprovalRepo) Create(ctx context.Context, req *domain.ApprovalRequest) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO approval_requests (id, tenant_id, entity_type, entity_id, request_type, old_value, new_value, reason, status, requested_by, reviewed_by, review_note, expires_at, created_at, reviewed_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		req.ID, nullStr(req.TenantID), req.EntityType, req.EntityID, req.RequestType,
		nullStr(req.OldValue), req.NewValue, req.Reason, string(req.Status),
		req.RequestedBy, nullStr(req.ReviewedBy), nullStr(req.ReviewNote),
		nullTime(req.ExpiresAt), req.CreatedAt, nullTimePtr(req.ReviewedAt))
	return err
}

func (r *PGApprovalRepo) GetByID(ctx context.Context, id string) (*domain.ApprovalRequest, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, COALESCE(tenant_id,''), entity_type, entity_id, request_type, COALESCE(old_value,''), new_value, reason, status, requested_by, COALESCE(reviewed_by,''), COALESCE(review_note,''), expires_at, created_at, reviewed_at
		 FROM approval_requests WHERE id=$1`, id)
	var req domain.ApprovalRequest
	var reviewedAt *time.Time
	err := row.Scan(&req.ID, &req.TenantID, &req.EntityType, &req.EntityID, &req.RequestType,
		&req.OldValue, &req.NewValue, &req.Reason, &req.Status, &req.RequestedBy,
		&req.ReviewedBy, &req.ReviewNote, &req.ExpiresAt, &req.CreatedAt, &reviewedAt)
	if err != nil {
		return nil, err
	}
	req.ReviewedAt = reviewedAt
	return &req, nil
}

func (r *PGApprovalRepo) GetByStatus(ctx context.Context, status domain.ApprovalStatus) ([]domain.ApprovalRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(tenant_id,''), entity_type, entity_id, request_type, COALESCE(old_value,''), new_value, reason, status, requested_by, COALESCE(reviewed_by,''), COALESCE(review_note,''), expires_at, created_at, reviewed_at
		 FROM approval_requests WHERE status=$1 ORDER BY created_at`, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApprovalRequests(rows)
}

func (r *PGApprovalRepo) GetByEntity(ctx context.Context, entityType, entityID string) ([]domain.ApprovalRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(tenant_id,''), entity_type, entity_id, request_type, COALESCE(old_value,''), new_value, reason, status, requested_by, COALESCE(reviewed_by,''), COALESCE(review_note,''), expires_at, created_at, reviewed_at
		 FROM approval_requests WHERE entity_type=$1 AND entity_id=$2 ORDER BY created_at`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApprovalRequests(rows)
}

func (r *PGApprovalRepo) UpdateStatus(ctx context.Context, id string, status domain.ApprovalStatus, reviewedBy, reviewNote string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE approval_requests SET status=$1, reviewed_by=$2, review_note=$3, reviewed_at=NOW() WHERE id=$4`,
		string(status), nullStr(reviewedBy), nullStr(reviewNote), id)
	return err
}

func (r *PGApprovalRepo) GetAll(ctx context.Context) ([]domain.ApprovalRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(tenant_id,''), entity_type, entity_id, request_type, COALESCE(old_value,''), new_value, reason, status, requested_by, COALESCE(reviewed_by,''), COALESCE(review_note,''), expires_at, created_at, reviewed_at
		 FROM approval_requests ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApprovalRequests(rows)
}

func scanApprovalRequests(rows pgx.Rows) ([]domain.ApprovalRequest, error) {
	var out []domain.ApprovalRequest
	for rows.Next() {
		var req domain.ApprovalRequest
		var reviewedAt *time.Time
		if err := rows.Scan(&req.ID, &req.TenantID, &req.EntityType, &req.EntityID, &req.RequestType,
			&req.OldValue, &req.NewValue, &req.Reason, &req.Status, &req.RequestedBy,
			&req.ReviewedBy, &req.ReviewNote, &req.ExpiresAt, &req.CreatedAt, &reviewedAt); err != nil {
			return nil, err
		}
		req.ReviewedAt = reviewedAt
		out = append(out, req)
	}
	return out, nil
}

// ─── COA: AccountVersion ─────────────────────────────────────────────────────

type PGAccountVersionRepo struct {
	pool *pgxpool.Pool
}

func NewPGAccountVersionRepo(pool *pgxpool.Pool) *PGAccountVersionRepo {
	return &PGAccountVersionRepo{pool: pool}
}

func (r *PGAccountVersionRepo) Create(ctx context.Context, v *domain.AccountVersion) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO account_versions (id, version_number, snapshot, change_summary, change_reason, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		v.ID, v.VersionNumber, v.Snapshot, nullStr(v.ChangeSummary), nullStr(v.ChangeReason), nullStr(v.CreatedBy), v.CreatedAt)
	return err
}

func (r *PGAccountVersionRepo) GetByVersionNumber(ctx context.Context, versionNumber string) (*domain.AccountVersion, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, version_number, snapshot, COALESCE(change_summary,''), COALESCE(change_reason,''), COALESCE(created_by,''), created_at
		 FROM account_versions WHERE version_number=$1`, versionNumber)
	var v domain.AccountVersion
	err := row.Scan(&v.ID, &v.VersionNumber, &v.Snapshot, &v.ChangeSummary, &v.ChangeReason, &v.CreatedBy, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PGAccountVersionRepo) GetLatest(ctx context.Context) (*domain.AccountVersion, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, version_number, snapshot, COALESCE(change_summary,''), COALESCE(change_reason,''), COALESCE(created_by,''), created_at
		 FROM account_versions ORDER BY version_number DESC LIMIT 1`)
	var v domain.AccountVersion
	err := row.Scan(&v.ID, &v.VersionNumber, &v.Snapshot, &v.ChangeSummary, &v.ChangeReason, &v.CreatedBy, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PGAccountVersionRepo) GetAll(ctx context.Context) ([]domain.AccountVersion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, version_number, snapshot, COALESCE(change_summary,''), COALESCE(change_reason,''), COALESCE(created_by,''), created_at
		 FROM account_versions ORDER BY version_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.AccountVersion
	for rows.Next() {
		var v domain.AccountVersion
		if err := rows.Scan(&v.ID, &v.VersionNumber, &v.Snapshot, &v.ChangeSummary, &v.ChangeReason, &v.CreatedBy, &v.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

// ─── COA: AccountMapping ─────────────────────────────────────────────────────

type PGAccountMappingRepo struct {
	pool *pgxpool.Pool
}

func NewPGAccountMappingRepo(pool *pgxpool.Pool) *PGAccountMappingRepo {
	return &PGAccountMappingRepo{pool: pool}
}

func (r *PGAccountMappingRepo) Create(ctx context.Context, m *domain.AccountMapping) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO account_mappings (id, source_regime, target_regime, old_code, new_code, mapping_type, split_ratio, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())`,
		m.ID, m.SourceRegime, m.TargetRegime, m.OldCode, m.NewCode, m.MappingType, nullFloat(m.SplitRatio))
	return err
}

func (r *PGAccountMappingRepo) GetByOldCode(ctx context.Context, sourceRegime, oldCode string) (*domain.AccountMapping, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, source_regime, target_regime, old_code, new_code, mapping_type, COALESCE(split_ratio,0), created_at
		 FROM account_mappings WHERE source_regime=$1 AND old_code=$2`, sourceRegime, oldCode)
	var m domain.AccountMapping
	err := row.Scan(&m.ID, &m.SourceRegime, &m.TargetRegime, &m.OldCode, &m.NewCode, &m.MappingType, &m.SplitRatio, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *PGAccountMappingRepo) GetByRegime(ctx context.Context, sourceRegime, targetRegime string) ([]domain.AccountMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, source_regime, target_regime, old_code, new_code, mapping_type, COALESCE(split_ratio,0), created_at
		 FROM account_mappings WHERE source_regime=$1 AND target_regime=$2 ORDER BY old_code`, sourceRegime, targetRegime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAccountMappings(rows)
}

func (r *PGAccountMappingRepo) GetAll(ctx context.Context) ([]domain.AccountMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, source_regime, target_regime, old_code, new_code, mapping_type, COALESCE(split_ratio,0), created_at
		 FROM account_mappings ORDER BY source_regime, old_code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAccountMappings(rows)
}

func scanAccountMappings(rows pgx.Rows) ([]domain.AccountMapping, error) {
	var out []domain.AccountMapping
	for rows.Next() {
		var m domain.AccountMapping
		if err := rows.Scan(&m.ID, &m.SourceRegime, &m.TargetRegime, &m.OldCode, &m.NewCode, &m.MappingType, &m.SplitRatio, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

// ─── COA: AccountAnalysis ────────────────────────────────────────────────────

type PGAccountAnalysisRepo struct {
	pool *pgxpool.Pool
}

func NewPGAccountAnalysisRepo(pool *pgxpool.Pool) *PGAccountAnalysisRepo {
	return &PGAccountAnalysisRepo{pool: pool}
}

func (r *PGAccountAnalysisRepo) Create(ctx context.Context, a *domain.AccountAnalysis) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO account_analysis (account_code, cost_center_id, profit_center_id, department_id, project_id, custom_dim1, custom_dim2)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		a.AccountCode, nullStr(a.CostCenterID), nullStr(a.ProfitCenterID), nullStr(a.DepartmentID),
		nullStr(a.ProjectID), nullStr(a.CustomDim1), nullStr(a.CustomDim2))
	return err
}

func (r *PGAccountAnalysisRepo) GetByAccount(ctx context.Context, accountCode string) (*domain.AccountAnalysis, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT account_code, COALESCE(cost_center_id,''), COALESCE(profit_center_id,''), COALESCE(department_id,''), COALESCE(project_id,''), COALESCE(custom_dim1,''), COALESCE(custom_dim2,'')
		 FROM account_analysis WHERE account_code=$1`, accountCode)
	var a domain.AccountAnalysis
	err := row.Scan(&a.AccountCode, &a.CostCenterID, &a.ProfitCenterID, &a.DepartmentID, &a.ProjectID, &a.CustomDim1, &a.CustomDim2)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *PGAccountAnalysisRepo) Update(ctx context.Context, a *domain.AccountAnalysis) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE account_analysis SET cost_center_id=$1, profit_center_id=$2, department_id=$3, project_id=$4, custom_dim1=$5, custom_dim2=$6
		 WHERE account_code=$7`,
		nullStr(a.CostCenterID), nullStr(a.ProfitCenterID), nullStr(a.DepartmentID), nullStr(a.ProjectID),
		nullStr(a.CustomDim1), nullStr(a.CustomDim2), a.AccountCode)
	return err
}

// ─── COA: IFRSMapping ────────────────────────────────────────────────────────

type PGIFRSMappingRepo struct {
	pool *pgxpool.Pool
}

func NewPGIFRSMappingRepo(pool *pgxpool.Pool) *PGIFRSMappingRepo {
	return &PGIFRSMappingRepo{pool: pool}
}

func (r *PGIFRSMappingRepo) Create(ctx context.Context, m *domain.IFRSMapping) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO ifrs_mappings (id, vas_code, ifrs_code, ifrs_name, reclassification_rule, adjustment_type, is_active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		m.ID, m.VASCode, m.IFRSCode, nullStr(m.IFRSName), nullStr(m.ReclassificationRule),
		nullStr(m.AdjustmentType), m.IsActive)
	return err
}

func (r *PGIFRSMappingRepo) GetByVASCode(ctx context.Context, vasCode string) (*domain.IFRSMapping, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, vas_code, ifrs_code, COALESCE(ifrs_name,''), COALESCE(reclassification_rule,''), COALESCE(adjustment_type,''), is_active
		 FROM ifrs_mappings WHERE vas_code=$1`, vasCode)
	var m domain.IFRSMapping
	err := row.Scan(&m.ID, &m.VASCode, &m.IFRSCode, &m.IFRSName, &m.ReclassificationRule, &m.AdjustmentType, &m.IsActive)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *PGIFRSMappingRepo) GetAll(ctx context.Context) ([]domain.IFRSMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, vas_code, ifrs_code, COALESCE(ifrs_name,''), COALESCE(reclassification_rule,''), COALESCE(adjustment_type,''), is_active
		 FROM ifrs_mappings ORDER BY vas_code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.IFRSMapping
	for rows.Next() {
		var m domain.IFRSMapping
		if err := rows.Scan(&m.ID, &m.VASCode, &m.IFRSCode, &m.IFRSName, &m.ReclassificationRule, &m.AdjustmentType, &m.IsActive); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func (r *PGIFRSMappingRepo) Update(ctx context.Context, m *domain.IFRSMapping) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE ifrs_mappings SET vas_code=$1, ifrs_code=$2, ifrs_name=$3, reclassification_rule=$4, adjustment_type=$5, is_active=$6 WHERE id=$7`,
		m.VASCode, m.IFRSCode, nullStr(m.IFRSName), nullStr(m.ReclassificationRule), nullStr(m.AdjustmentType), m.IsActive, m.ID)
	return err
}

