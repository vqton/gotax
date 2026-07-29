package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gotax/internal/domain"
)

type PGOpeningBalanceRepo struct {
	pool *pgxpool.Pool
}

func NewPGOpeningBalanceRepo(pool *pgxpool.Pool) *PGOpeningBalanceRepo {
	return &PGOpeningBalanceRepo{pool: pool}
}

func (r *PGOpeningBalanceRepo) Create(ctx context.Context, ob *domain.OpeningBalance) error {
	if ob.ID == "" {
		ob.ID = "OB" + time.Now().Format("20060102150405.000000000")
	}
	if ob.OriginalAmount == 0 {
		if ob.DebitAmount > 0 {
			ob.OriginalAmount = ob.DebitAmount
		} else if ob.CreditAmount > 0 {
			ob.OriginalAmount = ob.CreditAmount
		}
	}
	now := time.Now()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO opening_balances (id,company_id,period_id,fiscal_year_id,account_code,
		 currency_code,original_amount,debit_amount,credit_amount,exchange_rate,status,source_type,
		 batch_id,reason,created_by,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
		ob.ID, ob.CompanyID, ob.PeriodID, nullStr(ob.FiscalYearID), ob.AccountCode,
		ob.CurrencyCode, ob.OriginalAmount, ob.DebitAmount, ob.CreditAmount, nullF64(ob.ExchangeRate),
		string(ob.Status), ob.SourceType, nullStr(ob.BatchID), nullStr(ob.Reason), ob.CreatedBy,
		now, now)
	return err
}

func (r *PGOpeningBalanceRepo) GetByID(ctx context.Context, id string) (*domain.OpeningBalance, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,period_id,COALESCE(fiscal_year_id,''),account_code,
		 currency_code,original_amount,debit_amount,credit_amount,COALESCE(exchange_rate,0),
		 status,source_type,COALESCE(batch_id,''),COALESCE(reason,''),
		 COALESCE(approved_by,''),approved_at,
		 COALESCE(corrected_by,''),COALESCE(correction_of,''),COALESCE(correction_reason,''),
		 created_by,created_at,updated_at
		 FROM opening_balances WHERE id=$1`, id)
	return scanOpeningBalance(row)
}

func (r *PGOpeningBalanceRepo) List(ctx context.Context, filter domain.OBListFilter) ([]domain.OpeningBalance, error) {
	query := `SELECT id,company_id,period_id,COALESCE(fiscal_year_id,''),account_code,
		 currency_code,original_amount,debit_amount,credit_amount,COALESCE(exchange_rate,0),
		 status,source_type,COALESCE(batch_id,''),COALESCE(reason,''),
		 COALESCE(approved_by,''),approved_at,
		 COALESCE(corrected_by,''),COALESCE(correction_of,''),COALESCE(correction_reason,''),
		 created_by,created_at,updated_at
		 FROM opening_balances WHERE company_id=$1`
	args := []any{filter.CompanyID}
	argIdx := 2
	if filter.PeriodID != "" {
		query += ` AND period_id=$` + formatInt(argIdx); argIdx++
		args = append(args, filter.PeriodID)
	}
	if filter.Status != "" {
		query += ` AND status=$` + formatInt(argIdx); argIdx++
		args = append(args, string(filter.Status))
	}
	if filter.AccountCode != "" {
		query += ` AND account_code=$` + formatInt(argIdx); argIdx++
		args = append(args, filter.AccountCode)
	}
	if filter.Currency != "" {
		query += ` AND currency_code=$` + formatInt(argIdx); argIdx++
		args = append(args, filter.Currency)
	}
	query += ` ORDER BY account_code`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.OpeningBalance
	for rows.Next() {
		var b domain.OpeningBalance
		if err := scanOpeningBalanceRow(rows, &b); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	if out == nil {
		return []domain.OpeningBalance{}, nil
	}
	return out, nil
}

func (r *PGOpeningBalanceRepo) GetByAccount(ctx context.Context, companyID, periodID, accountCode string) (*domain.OpeningBalance, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,period_id,COALESCE(fiscal_year_id,''),account_code,
		 currency_code,original_amount,debit_amount,credit_amount,COALESCE(exchange_rate,0),
		 status,source_type,COALESCE(batch_id,''),COALESCE(reason,''),
		 COALESCE(approved_by,''),approved_at,
		 COALESCE(corrected_by,''),COALESCE(correction_of,''),COALESCE(correction_reason,''),
		 created_by,created_at,updated_at
		 FROM opening_balances WHERE company_id=$1 AND period_id=$2 AND account_code=$3 AND status='APPROVED'`,
		companyID, periodID, accountCode)
	return scanOpeningBalance(row)
}

func (r *PGOpeningBalanceRepo) Update(ctx context.Context, ob *domain.OpeningBalance) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE opening_balances SET account_code=$1,currency_code=$2,original_amount=$3,
		 debit_amount=$4,credit_amount=$5,exchange_rate=$6,source_type=$7,batch_id=$8,
		 reason=$9,fiscal_year_id=$10,updated_at=NOW() WHERE id=$11 AND status NOT IN ('APPROVED','CORRECTED')`,
		ob.AccountCode, ob.CurrencyCode, ob.OriginalAmount,
		ob.DebitAmount, ob.CreditAmount, nullF64(ob.ExchangeRate), ob.SourceType,
		nullStr(ob.BatchID), nullStr(ob.Reason), nullStr(ob.FiscalYearID), ob.ID)
	return err
}

func (r *PGOpeningBalanceRepo) UpdateStatus(ctx context.Context, id string, status domain.OpeningBalanceStatus, approvedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE opening_balances SET status=$1,approved_by=CASE WHEN $1='APPROVED' THEN $2 ELSE approved_by END,approved_at=CASE WHEN $1='APPROVED' THEN NOW() ELSE approved_at END,updated_at=NOW() WHERE id=$3`,
		string(status), nullStr(approvedBy), id)
	return err
}

func (r *PGOpeningBalanceRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM opening_balances WHERE id=$1 AND status NOT IN ('APPROVED','CORRECTED')`, id)
	return err
}

func (r *PGOpeningBalanceRepo) BulkCreate(ctx context.Context, balances []domain.OpeningBalance) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	now := time.Now()
	for i := range balances {
		if balances[i].ID == "" {
			balances[i].ID = "OB" + time.Now().Format("20060102150405.000000000")
		}
		if balances[i].OriginalAmount == 0 {
			if balances[i].DebitAmount > 0 {
				balances[i].OriginalAmount = balances[i].DebitAmount
			} else if balances[i].CreditAmount > 0 {
				balances[i].OriginalAmount = balances[i].CreditAmount
			}
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO opening_balances (id,company_id,period_id,fiscal_year_id,account_code,
			 currency_code,original_amount,debit_amount,credit_amount,exchange_rate,status,source_type,
			 batch_id,reason,created_by,created_at,updated_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
			balances[i].ID, balances[i].CompanyID, balances[i].PeriodID, nullStr(balances[i].FiscalYearID),
			balances[i].AccountCode, balances[i].CurrencyCode, balances[i].OriginalAmount,
			balances[i].DebitAmount, balances[i].CreditAmount, nullF64(balances[i].ExchangeRate),
			string(balances[i].Status), balances[i].SourceType, nullStr(balances[i].BatchID),
			nullStr(balances[i].Reason), balances[i].CreatedBy, now, now); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PGOpeningBalanceRepo) BulkUpdateStatus(ctx context.Context, ids []string, status domain.OpeningBalanceStatus, approvedBy string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for _, id := range ids {
		if _, err := tx.Exec(ctx,
			`UPDATE opening_balances SET status=$1,approved_by=CASE WHEN $1='APPROVED' THEN $2 ELSE approved_by END,approved_at=CASE WHEN $1='APPROVED' THEN NOW() ELSE approved_at END,updated_at=NOW() WHERE id=$3`,
			string(status), nullStr(approvedBy), id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PGOpeningBalanceRepo) CreateDetail(ctx context.Context, d *domain.OpeningBalanceDetail) error {
	if d.ID == "" {
		d.ID = "OBD" + time.Now().Format("20060102150405.000000000")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO opening_balance_details (id,opening_balance_id,entity_type,entity_id,entity_name,
		 debit_amount,credit_amount,quantity,unit_price,original_cost,acc_depreciation,
		 counterpart_account,note)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		d.ID, d.OpeningBalanceID, string(d.EntityType), d.EntityID, nullStr(d.EntityName),
		d.DebitAmount, d.CreditAmount, 		nullF64(d.Quantity), nullF64(d.UnitPrice),
		nullF64(d.OriginalCost), nullF64(d.AccDepreciation),
		nullStr(d.CounterpartAccount), nullStr(d.Note))
	return err
}

func (r *PGOpeningBalanceRepo) GetDetails(ctx context.Context, balanceID string) ([]domain.OpeningBalanceDetail, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,opening_balance_id,entity_type,entity_id,COALESCE(entity_name,''),
		 debit_amount,credit_amount,COALESCE(quantity,0),COALESCE(unit_price,0),
		 COALESCE(original_cost,0),COALESCE(acc_depreciation,0),
		 COALESCE(counterpart_account,''),COALESCE(note,''),created_at
		 FROM opening_balance_details WHERE opening_balance_id=$1
		 ORDER BY entity_type,entity_id`, balanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.OpeningBalanceDetail
	for rows.Next() {
		var d domain.OpeningBalanceDetail
		if err := rows.Scan(&d.ID, &d.OpeningBalanceID, &d.EntityType, &d.EntityID, &d.EntityName,
			&d.DebitAmount, &d.CreditAmount, &d.Quantity, &d.UnitPrice,
			&d.OriginalCost, &d.AccDepreciation, &d.CounterpartAccount, &d.Note, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	if out == nil {
		return []domain.OpeningBalanceDetail{}, nil
	}
	return out, nil
}

func (r *PGOpeningBalanceRepo) DeleteDetail(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM opening_balance_details WHERE id=$1`, id)
	return err
}

func (r *PGOpeningBalanceRepo) DeleteDetails(ctx context.Context, balanceID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM opening_balance_details WHERE opening_balance_id=$1`, balanceID)
	return err
}

func (r *PGOpeningBalanceRepo) GetTotals(ctx context.Context, companyID, periodID string) (float64, float64, error) {
	var totalDebit, totalCredit float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(debit_amount),0), COALESCE(SUM(credit_amount),0)
		 FROM opening_balances WHERE company_id=$1 AND period_id=$2 AND status='APPROVED'`,
		companyID, periodID).Scan(&totalDebit, &totalCredit)
	return totalDebit, totalCredit, err
}

func (r *PGOpeningBalanceRepo) ValidateBalanced(ctx context.Context, companyID, periodID string) (bool, error) {
	var diff float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(ABS(SUM(debit_amount)-SUM(credit_amount)),0)
		 FROM opening_balances WHERE company_id=$1 AND period_id=$2 AND status='APPROVED'`,
		companyID, periodID).Scan(&diff)
	return diff < 0.01, err
}

func (r *PGOpeningBalanceRepo) CreateCarryForwardLog(ctx context.Context, log *domain.CarryForwardLog) error {
	if log.ID == "" {
		log.ID = "CF" + time.Now().Format("20060102150405.000000000")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO carry_forward_logs (id,company_id,from_period_id,to_period_id,
		 from_fiscal_year,to_fiscal_year,account_count,total_debit,total_credit,status,executed_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		log.ID, log.CompanyID, log.FromPeriodID, log.ToPeriodID,
		log.FromFiscalYear, log.ToFiscalYear, log.AccountCount, log.TotalDebit, log.TotalCredit,
		log.Status, log.ExecutedBy)
	return err
}

func (r *PGOpeningBalanceRepo) GetCarryForwardLogs(ctx context.Context, companyID string) ([]domain.CarryForwardLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,from_period_id,to_period_id,from_fiscal_year,to_fiscal_year,
		 account_count,total_debit,total_credit,closing_entry_ids,status,executed_by,executed_at
		 FROM carry_forward_logs WHERE company_id=$1 ORDER BY executed_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.CarryForwardLog
	for rows.Next() {
		var l domain.CarryForwardLog
		if err := rows.Scan(&l.ID, &l.CompanyID, &l.FromPeriodID, &l.ToPeriodID,
			&l.FromFiscalYear, &l.ToFiscalYear, &l.AccountCount, &l.TotalDebit, &l.TotalCredit,
			&l.ClosingEntryIDs, &l.Status, &l.ExecutedBy, &l.ExecutedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	if out == nil {
		return []domain.CarryForwardLog{}, nil
	}
	return out, nil
}

func (r *PGOpeningBalanceRepo) GetCarryForwardLogByID(ctx context.Context, id string) (*domain.CarryForwardLog, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,from_period_id,to_period_id,from_fiscal_year,to_fiscal_year,
		 account_count,total_debit,total_credit,closing_entry_ids,status,executed_by,executed_at
		 FROM carry_forward_logs WHERE id=$1`, id)
	var l domain.CarryForwardLog
	err := row.Scan(&l.ID, &l.CompanyID, &l.FromPeriodID, &l.ToPeriodID,
		&l.FromFiscalYear, &l.ToFiscalYear, &l.AccountCount, &l.TotalDebit, &l.TotalCredit,
		&l.ClosingEntryIDs, &l.Status, &l.ExecutedBy, &l.ExecutedAt)
	if err != nil {
		return nil, domain.ErrOpeningBalanceNotFound
	}
	return &l, nil
}

func (r *PGOpeningBalanceRepo) CreateCircular99Mapping(ctx context.Context, m *domain.Circular99Mapping) error {
	if m.ID == "" {
		m.ID = "C99" + time.Now().Format("20060102150405.000000000")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO circular99_mappings (id,old_account_code,new_account_code,mapping_type,
		 split_ratio,counterpart_account,effective_date,note,is_active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		m.ID, m.OldAccountCode, m.NewAccountCode, m.MappingType,
		nullF64(m.SplitRatio), nullStr(m.CounterpartAccount), m.EffectiveDate, nullStr(m.Note), m.IsActive)
	return err
}

func (r *PGOpeningBalanceRepo) ListCircular99Mappings(ctx context.Context) ([]domain.Circular99Mapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,old_account_code,new_account_code,mapping_type,COALESCE(split_ratio,0),
		 COALESCE(counterpart_account,''),effective_date,COALESCE(note,''),is_active
		 FROM circular99_mappings WHERE is_active=true ORDER BY old_account_code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Circular99Mapping
	for rows.Next() {
		var m domain.Circular99Mapping
		if err := rows.Scan(&m.ID, &m.OldAccountCode, &m.NewAccountCode, &m.MappingType,
			&m.SplitRatio, &m.CounterpartAccount, &m.EffectiveDate, &m.Note, &m.IsActive); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if out == nil {
		return []domain.Circular99Mapping{}, nil
	}
	return out, nil
}

func (r *PGOpeningBalanceRepo) GetCircular99MappingByOldCode(ctx context.Context, oldCode string) (*domain.Circular99Mapping, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,old_account_code,new_account_code,mapping_type,COALESCE(split_ratio,0),
		 COALESCE(counterpart_account,''),effective_date,COALESCE(note,''),is_active
		 FROM circular99_mappings WHERE old_account_code=$1 AND is_active=true`, oldCode)
	var m domain.Circular99Mapping
	err := row.Scan(&m.ID, &m.OldAccountCode, &m.NewAccountCode, &m.MappingType,
		&m.SplitRatio, &m.CounterpartAccount, &m.EffectiveDate, &m.Note, &m.IsActive)
	if err != nil {
		return nil, domain.ErrCircular99MappingNotFound
	}
	return &m, nil
}

func (r *PGOpeningBalanceRepo) CreateMigration(ctx context.Context, m *domain.BalanceMigration) error {
	if m.ID == "" {
		m.ID = "MIG" + time.Now().Format("20060102150405.000000000")
	}
	_, err := r.pool.Exec(ctx,
		`INSERT INTO balance_migrations (id,company_id,from_regime,to_regime,execution_date,status,
		 source_balance_id,target_balance_id,journal_entry_id,summary,executed_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		m.ID, m.CompanyID, m.FromRegime, m.ToRegime, m.ExecutionDate, m.Status,
		nullStr(m.SourceBalanceID), nullStr(m.TargetBalanceID), nullStr(m.JournalEntryID),
		nullStr(m.Summary), m.ExecutedBy)
	return err
}

func (r *PGOpeningBalanceRepo) GetMigrationByID(ctx context.Context, id string) (*domain.BalanceMigration, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id,company_id,from_regime,to_regime,execution_date,status,
		 COALESCE(source_balance_id,''),COALESCE(target_balance_id,''),COALESCE(journal_entry_id,''),
		 COALESCE(summary,''),executed_by,created_at,executed_at
		 FROM balance_migrations WHERE id=$1`, id)
	var m domain.BalanceMigration
	err := row.Scan(&m.ID, &m.CompanyID, &m.FromRegime, &m.ToRegime, &m.ExecutionDate, &m.Status,
		&m.SourceBalanceID, &m.TargetBalanceID, &m.JournalEntryID,
		&m.Summary, &m.ExecutedBy, &m.CreatedAt, &m.ExecutedAt)
	if err != nil {
		return nil, domain.ErrOpeningBalanceNotFound
	}
	return &m, nil
}

func (r *PGOpeningBalanceRepo) ListMigrations(ctx context.Context, companyID string) ([]domain.BalanceMigration, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,from_regime,to_regime,execution_date,status,
		 COALESCE(source_balance_id,''),COALESCE(target_balance_id,''),COALESCE(journal_entry_id,''),
		 COALESCE(summary,''),executed_by,created_at,executed_at
		 FROM balance_migrations WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.BalanceMigration
	for rows.Next() {
		var m domain.BalanceMigration
		if err := rows.Scan(&m.ID, &m.CompanyID, &m.FromRegime, &m.ToRegime, &m.ExecutionDate, &m.Status,
			&m.SourceBalanceID, &m.TargetBalanceID, &m.JournalEntryID,
			&m.Summary, &m.ExecutedBy, &m.CreatedAt, &m.ExecutedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if out == nil {
		return []domain.BalanceMigration{}, nil
	}
	return out, nil
}

// ─── Scanners ──────────────────────────────────────────────────────

func scanOpeningBalance(row pgx.Row) (*domain.OpeningBalance, error) {
	var b domain.OpeningBalance
	err := row.Scan(&b.ID, &b.CompanyID, &b.PeriodID, &b.FiscalYearID, &b.AccountCode,
		&b.CurrencyCode, &b.OriginalAmount, &b.DebitAmount, &b.CreditAmount, &b.ExchangeRate,
		&b.Status, &b.SourceType, &b.BatchID, &b.Reason,
		&b.ApprovedBy, &b.ApprovedAt,
		&b.CorrectedBy, &b.CorrectionOf, &b.CorrectionReason,
		&b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, domain.ErrOpeningBalanceNotFound
	}
	return &b, nil
}

func scanOpeningBalanceRow(row pgx.Row, b *domain.OpeningBalance) error {
	return row.Scan(&b.ID, &b.CompanyID, &b.PeriodID, &b.FiscalYearID, &b.AccountCode,
		&b.CurrencyCode, &b.OriginalAmount, &b.DebitAmount, &b.CreditAmount, &b.ExchangeRate,
		&b.Status, &b.SourceType, &b.BatchID, &b.Reason,
		&b.ApprovedBy, &b.ApprovedAt,
		&b.CorrectedBy, &b.CorrectionOf, &b.CorrectionReason,
		&b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
}
