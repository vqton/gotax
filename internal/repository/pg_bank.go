package repository

import (
	"context"
	"math"
	"time"

	"gorm.io/gorm"

	"gotax/internal/domain"
)

type PGBankRepo struct {
	db *gorm.DB
}

func NewPGBankRepo(db *gorm.DB) *PGBankRepo {
	return &PGBankRepo{db: db}
}

// ─── Conversion helpers ────────────────────────────────────────────────

func gormBankStatementToDomain(m *domain.BankStatementGORM) *domain.BankStatement {
	return &domain.BankStatement{
		ID:            m.ID,
		CompanyID:     m.CompanyID,
		BankAccountID: m.BankAccountID,
		StatementDate: safeTimeStr(m.StatementDate),
		FromDate:      safeTimePtrStr(m.FromDate),
		ToDate:        safeTimePtrStr(m.ToDate),
		OpeningBalance: m.OpeningBalance,
		ClosingBalance: m.ClosingBalance,
		TotalCredits:  m.TotalCredits,
		TotalDebits:   m.TotalDebits,
		LineCount:     m.LineCount,
		Currency:      m.Currency,
		ImportMethod:  m.ImportMethod,
		RawFileName:   safeStr(m.RawFileName),
		RawFileHash:   safeStr(m.RawFileHash),
		Status:        domain.BankStatementStatus(m.Status),
		ImportedBy:    m.ImportedBy,
		ImportedAt:    safeTimePtrRFC3339(m.ImportedAt),
		Notes:         safeStr(m.Notes),
	}
}

func domainToGormBankStatement(s *domain.BankStatement) *domain.BankStatementGORM {
	return &domain.BankStatementGORM{
		ID:            s.ID,
		CompanyID:     s.CompanyID,
		BankAccountID: s.BankAccountID,
		StatementDate: parseDate(s.StatementDate),
		FromDate:      timePtr(parseDate(s.FromDate)),
		ToDate:        timePtr(parseDate(s.ToDate)),
		OpeningBalance: s.OpeningBalance,
		ClosingBalance: s.ClosingBalance,
		TotalCredits:  s.TotalCredits,
		TotalDebits:   s.TotalDebits,
		LineCount:     s.LineCount,
		Currency:      s.Currency,
		ImportMethod:  s.ImportMethod,
		RawFileName:   strPtr(s.RawFileName),
		RawFileHash:   strPtr(s.RawFileHash),
		Status:        string(s.Status),
		ImportedBy:    s.ImportedBy,
		ImportedAt:    timePtr(parseDateTime(s.ImportedAt)),
		Notes:         strPtr(s.Notes),
	}
}

func gormStatementLineToDomain(l *domain.BankStatementLineGORM) *domain.BankStatementLine {
	return &domain.BankStatementLine{
		ID:             l.ID,
		StatementID:    l.StatementID,
		TransactionDate: safeTimeStr(l.TransactionDate),
		ValueDate:      safeTimePtrStr(l.ValueDate),
		Description:    safeStr(l.Description),
		DebitAmount:    l.DebitAmount,
		CreditAmount:   l.CreditAmount,
		BalanceAfter:   l.BalanceAfter,
		ReferenceNo:    safeStr(l.ReferenceNo),
		BankRef:        safeStr(l.BankRef),
		Counterparty:   safeStr(l.Counterparty),
		CounterpartyAcc: safeStr(l.CounterpartyAcc),
		CounterpartyBank: safeStr(l.CounterpartyBank),
		RawData:        safeStr(l.RawData),
		MatchStatus:    domain.MatchStatus(l.MatchStatus),
		MatchedLineID:  safeStr(l.MatchedLineID),
		MatchedAt:      safeTimePtrRFC3339(l.MatchedAt),
		MatchedBy:      safeStr(l.MatchedBy),
		CreatedAt:      safeTimePtrRFC3339(&l.CreatedAt),
	}
}

func domainToGormStatementLine(l *domain.BankStatementLine) *domain.BankStatementLineGORM {
	return &domain.BankStatementLineGORM{
		ID:             l.ID,
		StatementID:    l.StatementID,
		TransactionDate: parseDate(l.TransactionDate),
		ValueDate:      timePtr(parseDate(l.ValueDate)),
		Description:    strPtr(l.Description),
		DebitAmount:    l.DebitAmount,
		CreditAmount:   l.CreditAmount,
		BalanceAfter:   l.BalanceAfter,
		ReferenceNo:    strPtr(l.ReferenceNo),
		BankRef:        strPtr(l.BankRef),
		Counterparty:   strPtr(l.Counterparty),
		CounterpartyAcc: strPtr(l.CounterpartyAcc),
		CounterpartyBank: strPtr(l.CounterpartyBank),
		RawData:        strPtr(l.RawData),
		MatchStatus:    string(l.MatchStatus),
		MatchedLineID:  strPtr(l.MatchedLineID),
		MatchedAt:      stringToTimePtr(l.MatchedAt),
		MatchedBy:      strPtr(l.MatchedBy),
		CreatedAt:      time.Now(),
	}
}

func gormReconToDomain(m *domain.BankReconciliationGORM) *domain.BankReconciliation {
	return &domain.BankReconciliation{
		ID:              m.ID,
		CompanyID:       m.CompanyID,
		BankAccountID:   m.BankAccountID,
		StatementID:     m.StatementID,
		FromDate:        safeTimePtrStr(m.FromDate),
		ToDate:          safeTimePtrStr(m.ToDate),
		OpeningBalance:  m.OpeningBalance,
		ClosingBalance:  m.ClosingBalance,
		StatementBalance: m.StatementBalance,
		Difference:      m.Difference,
		Status:          domain.ReconStatus(m.Status),
		MatchedLines:    m.MatchedLines,
		UnmatchedLines:  m.UnmatchedLines,
		WriteOffAmount:  m.WriteOffAmount,
		CompletedBy:     m.CompletedBy,
		CompletedAt:     safeTimePtrRFC3339(m.CompletedAt),
		ReversedAt:      safeTimePtrRFC3339(m.ReversedAt),
		Notes:           safeStr(m.Notes),
		CreatedAt:       safeTimePtrRFC3339(&m.CreatedAt),
	}
}

func domainToGormRecon(rc *domain.BankReconciliation) *domain.BankReconciliationGORM {
	return &domain.BankReconciliationGORM{
		ID:              rc.ID,
		CompanyID:       rc.CompanyID,
		BankAccountID:   rc.BankAccountID,
		StatementID:     rc.StatementID,
		FromDate:        timePtr(parseDate(rc.FromDate)),
		ToDate:          timePtr(parseDate(rc.ToDate)),
		OpeningBalance:  rc.OpeningBalance,
		ClosingBalance:  rc.ClosingBalance,
		StatementBalance: rc.StatementBalance,
		Difference:      rc.Difference,
		Status:          string(rc.Status),
		MatchedLines:    rc.MatchedLines,
		UnmatchedLines:  rc.UnmatchedLines,
		WriteOffAmount:  rc.WriteOffAmount,
		CompletedBy:     rc.CompletedBy,
		CompletedAt:     stringToTimePtr(rc.CompletedAt),
		ReversedAt:      stringToTimePtr(rc.ReversedAt),
		Notes:           strPtr(rc.Notes),
		CreatedAt:       time.Now(),
	}
}

func gormReconMatchToDomain(m *domain.BankReconciliationMatchGORM) *domain.BankReconciliationMatch {
	return &domain.BankReconciliationMatch{
		ID:              m.ID,
		ReconciliationID: m.ReconciliationID,
		StatementLineID: m.StatementLineID,
		TransactionType: m.TransactionType,
		TransactionID:   m.TransactionID,
		TransactionRef:  m.TransactionRef,
		MatchMethod:     m.MatchMethod,
		Confidence:      m.Confidence,
		CreatedAt:       safeTimePtrRFC3339(&m.CreatedAt),
	}
}

func domainToGormReconMatch(m *domain.BankReconciliationMatch) *domain.BankReconciliationMatchGORM {
	return &domain.BankReconciliationMatchGORM{
		ID:              m.ID,
		ReconciliationID: m.ReconciliationID,
		StatementLineID: m.StatementLineID,
		TransactionType: m.TransactionType,
		TransactionID:   m.TransactionID,
		TransactionRef:  m.TransactionRef,
		MatchMethod:     m.MatchMethod,
		Confidence:      m.Confidence,
		CreatedAt:       time.Now(),
	}
}

func gormPaymentOrderToDomain(m *domain.PaymentOrderGORM) *domain.PaymentOrder {
	return &domain.PaymentOrder{
		ID:             m.ID,
		CompanyID:      m.CompanyID,
		PaymentDate:    m.DueDate.Format("2006-01-02"),
		Amount:         m.Amount,
		Currency:       m.Currency,
		BeneficiaryName: m.PayeeName,
		PaymentContent: safeStr(m.Purpose),
		Status:         domain.PaymentOrderStatus(m.Status),
		CreatedBy:      m.CreatedBy,
		CreatedAt:      m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      m.UpdatedAt.Format(time.RFC3339),
	}
}

func domainToGormPaymentOrder(po *domain.PaymentOrder) *domain.PaymentOrderGORM {
	dueDate, _ := time.Parse("2006-01-02", po.PaymentDate)
	now := time.Now()
	return &domain.PaymentOrderGORM{
		ID:        po.ID,
		CompanyID: po.CompanyID,
		PayeeName: po.BeneficiaryName,
		PayeeBank: nullStrG(po.BeneficiaryBank),
		PayeeAccount: nullStrG(po.BeneficiaryAcc),
		Amount:    po.Amount,
		Currency:  po.Currency,
		Status:    string(po.Status),
		DueDate:   dueDate,
		Purpose:   nullStrG(po.PaymentContent),
		CreatedBy: po.CreatedBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func gormPaymentBatchToDomain(m *domain.PaymentOrderBatchGORM) *domain.PaymentOrderBatch {
	return &domain.PaymentOrderBatch{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		BatchName:   m.BatchNumber,
		TotalAmount: m.TotalAmount,
		OrderCount:  m.OrderCount,
		Status:      domain.BatchStatus(m.Status),
		CreatedBy:   safeStr(m.SubmittedBy),
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
	}
}

func domainToGormPaymentBatch(b *domain.PaymentOrderBatch) *domain.PaymentOrderBatchGORM {
	now := time.Now()
	return &domain.PaymentOrderBatchGORM{
		ID:          b.ID,
		CompanyID:   b.CompanyID,
		BatchNumber: b.BatchName,
		TotalAmount: b.TotalAmount,
		OrderCount:  b.OrderCount,
		Status:      string(b.Status),
		SubmittedBy: nullStrG(b.CreatedBy),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func gormLoanToDomain(m *domain.LoanAgreementGORM) *domain.LoanAgreement {
	return &domain.LoanAgreement{
		ID:                m.ID,
		CompanyID:         m.CompanyID,
		BankAccountID:     m.LenderName,
		PrincipalAmount:   m.Principal,
		Currency:          m.Currency,
		InterestRate:      m.InterestRate,
		StartDate:         m.StartDate.Format("2006-01-02"),
		MaturityDate:      m.MaturityDate.Format("2006-01-02"),
		Status:            domain.LoanStatus(m.Status),
		CreatedAt:         m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         m.UpdatedAt.Format(time.RFC3339),
	}
}

func domainToGormLoan(l *domain.LoanAgreement) *domain.LoanAgreementGORM {
	startDate, _ := time.Parse("2006-01-02", l.StartDate)
	maturityDate, _ := time.Parse("2006-01-02", l.MaturityDate)
	now := time.Now()
	return &domain.LoanAgreementGORM{
		ID:           l.ID,
		CompanyID:    l.CompanyID,
		LenderName:   l.BankAccountID,
		Principal:    l.PrincipalAmount,
		Currency:     l.Currency,
		InterestRate: l.InterestRate,
		StartDate:    startDate,
		MaturityDate: maturityDate,
		Status:       string(l.Status),
		CreatedBy:    "",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func gormDisbursementToDomain(m *domain.LoanDisbursementGORM) *domain.LoanDisbursement {
	return &domain.LoanDisbursement{
		ID:           m.ID,
		LoanID:       m.LoanID,
		DisbursementDate: m.DisburseDate.Format("2006-01-02"),
		Amount:       m.Amount,
		ReferenceNo:  safeStr(m.BankRef),
		Notes:        safeStr(m.Note),
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
	}
}

func domainToGormDisbursement(d *domain.LoanDisbursement) *domain.LoanDisbursementGORM {
	disburseDate, _ := time.Parse("2006-01-02", d.DisbursementDate)
	return &domain.LoanDisbursementGORM{
		ID:        d.ID,
		LoanID:    d.LoanID,
		DisburseDate: disburseDate,
		Amount:    d.Amount,
		BankRef:   nullStrG(d.ReferenceNo),
		Note:      nullStrG(d.Notes),
		CreatedAt: time.Now(),
	}
}

func gormRepaymentToDomain(m *domain.LoanRepaymentGORM) *domain.LoanRepayment {
	return &domain.LoanRepayment{
		ID:             m.ID,
		LoanID:         m.LoanID,
		RepaymentDate:  m.RepayDate.Format("2006-01-02"),
		PrincipalAmount: m.Principal,
		InterestAmount: m.Interest,
		TotalAmount:    m.TotalAmount,
		CreatedAt:      m.CreatedAt.Format(time.RFC3339),
	}
}

func domainToGormRepayment(rp *domain.LoanRepayment) *domain.LoanRepaymentGORM {
	repayDate, _ := time.Parse("2006-01-02", rp.RepaymentDate)
	return &domain.LoanRepaymentGORM{
		ID:          rp.ID,
		LoanID:      rp.LoanID,
		RepayDate:   repayDate,
		Principal:   rp.PrincipalAmount,
		Interest:    rp.InterestAmount,
		TotalAmount: rp.TotalAmount,
		CreatedBy:   "",
		CreatedAt:   time.Now(),
	}
}

func gormDepositToDomain(m *domain.TermDepositGORM) *domain.TermDeposit {
	return &domain.TermDeposit{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		DepositNo:   safeStr(m.CertificateNumber),
		Amount:      m.Principal,
		Currency:    m.Currency,
		InterestRate: m.InterestRate,
		TermDays:    m.TermDays,
		StartDate:   m.StartDate.Format("2006-01-02"),
		MaturityDate: m.MaturityDate.Format("2006-01-02"),
		AutoRenewal: m.AutoRenew,
		Status:      domain.DepositStatus(m.Status),
		MaturedAt:   safeTimePtrRFC3339(m.MaturedAt),
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
	}
}

func domainToGormDeposit(d *domain.TermDeposit) *domain.TermDepositGORM {
	startDate, _ := time.Parse("2006-01-02", d.StartDate)
	maturityDate, _ := time.Parse("2006-01-02", d.MaturityDate)
	now := time.Now()
	return &domain.TermDepositGORM{
		ID:          d.ID,
		CompanyID:   d.CompanyID,
		Principal:   d.Amount,
		Currency:    d.Currency,
		InterestRate: d.InterestRate,
		TermDays:    d.TermDays,
		StartDate:   startDate,
		MaturityDate: maturityDate,
		AutoRenew:   d.AutoRenewal,
		Status:      string(d.Status),
		MaturedAt:   stringToTimePtr(d.MaturedAt),
		CreatedBy:   "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ─── Helpers ───────────────────────────────────────────────────────────

func stringToTimePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			return nil
		}
		return &t
	}
	return &t
}

// ─── Statements ────────────────────────────────────────────────────────

func (r *PGBankRepo) CreateStatement(ctx context.Context, s *domain.BankStatement) error {
	m := domainToGormBankStatement(s)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	s.ID = m.ID
	return nil
}

func (r *PGBankRepo) GetStatement(ctx context.Context, id string) (*domain.BankStatement, error) {
	var m domain.BankStatementGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrBankStatementNotFound
	}
	return gormBankStatementToDomain(&m), nil
}

func (r *PGBankRepo) ListStatements(ctx context.Context, companyID, bankAccountID string, limit, offset int) ([]domain.BankStatement, int, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.BankStatementGORM{}).
		Where("company_id = ? AND bank_account_id = ?", companyID, bankAccountID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []domain.BankStatementGORM
	if err := r.db.WithContext(ctx).
		Where("company_id = ? AND bank_account_id = ?", companyID, bankAccountID).
		Order("statement_date DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.BankStatement, len(models))
	for i := range models {
		out[i] = *gormBankStatementToDomain(&models[i])
	}
	if out == nil {
		out = []domain.BankStatement{}
	}
	return out, int(total), nil
}

func (r *PGBankRepo) DeleteStatement(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Where("id = ? AND status = ?", id, domain.BankStatementImported).Delete(&domain.BankStatementGORM{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrBankStatementNotFound
	}
	return nil
}

// ─── Statement Lines ────────────────────────────────────────────────────

func (r *PGBankRepo) CreateStatementLines(ctx context.Context, lines []domain.BankStatementLine) error {
	if len(lines) == 0 {
		return nil
	}
	models := make([]domain.BankStatementLineGORM, len(lines))
	for i := range lines {
		m := domainToGormStatementLine(&lines[i])
		models[i] = *m
	}
	return r.db.WithContext(ctx).Create(&models).Error
}

func (r *PGBankRepo) GetStatementLines(ctx context.Context, statementID string) ([]domain.BankStatementLine, error) {
	var models []domain.BankStatementLineGORM
	if err := r.db.WithContext(ctx).Where("statement_id = ?", statementID).Order("transaction_date").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.BankStatementLine, len(models))
	for i := range models {
		out[i] = *gormStatementLineToDomain(&models[i])
	}
	if out == nil {
		out = []domain.BankStatementLine{}
	}
	return out, nil
}

func (r *PGBankRepo) GetStatementLinesByStatus(ctx context.Context, statementID string, status domain.MatchStatus) ([]domain.BankStatementLine, error) {
	var models []domain.BankStatementLineGORM
	if err := r.db.WithContext(ctx).Where("statement_id = ? AND match_status = ?", statementID, string(status)).Order("transaction_date").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.BankStatementLine, len(models))
	for i := range models {
		out[i] = *gormStatementLineToDomain(&models[i])
	}
	if out == nil {
		out = []domain.BankStatementLine{}
	}
	return out, nil
}

func (r *PGBankRepo) UpdateStatementLineMatch(ctx context.Context, lineID string, matchStatus domain.MatchStatus, matchedLineID, matchedBy string) error {
	res := r.db.WithContext(ctx).Model(&domain.BankStatementLineGORM{}).Where("id = ?", lineID).Updates(map[string]interface{}{
		"match_status":   string(matchStatus),
		"matched_line_id": matchedLineID,
		"matched_by":     matchedBy,
		"matched_at":     time.Now().Format(time.RFC3339),
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrStatementLineNotFound
	}
	return nil
}

// ─── Reconciliation ──────────────────────────────────────────────────────

func (r *PGBankRepo) CreateReconciliation(ctx context.Context, rc *domain.BankReconciliation) error {
	m := domainToGormRecon(rc)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	rc.ID = m.ID
	return nil
}

func (r *PGBankRepo) GetReconciliation(ctx context.Context, id string) (*domain.BankReconciliation, error) {
	var m domain.BankReconciliationGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrReconciliationNotFound
	}
	return gormReconToDomain(&m), nil
}

func (r *PGBankRepo) ListReconciliations(ctx context.Context, companyID, bankAccountID string) ([]domain.BankReconciliation, error) {
	var models []domain.BankReconciliationGORM
	if err := r.db.WithContext(ctx).Where("company_id = ? AND bank_account_id = ?", companyID, bankAccountID).Order("from_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.BankReconciliation, len(models))
	for i := range models {
		out[i] = *gormReconToDomain(&models[i])
	}
	if out == nil {
		out = []domain.BankReconciliation{}
	}
	return out, nil
}

func (r *PGBankRepo) UpdateReconciliation(ctx context.Context, rc *domain.BankReconciliation) error {
	res := r.db.WithContext(ctx).Model(&domain.BankReconciliationGORM{}).Where("id = ?", rc.ID).Updates(map[string]interface{}{
		"closing_balance":  rc.ClosingBalance,
		"statement_balance": rc.StatementBalance,
		"difference":       rc.Difference,
		"status":           string(rc.Status),
		"matched_lines":    rc.MatchedLines,
		"unmatched_lines":  rc.UnmatchedLines,
		"write_off_amount": rc.WriteOffAmount,
		"completed_by":     rc.CompletedBy,
		"completed_at":     rc.CompletedAt,
		"notes":            rc.Notes,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrReconciliationNotFound
	}
	return nil
}

// ─── Reconciliation Matches ─────────────────────────────────────────────

func (r *PGBankRepo) CreateReconciliationMatch(ctx context.Context, m *domain.BankReconciliationMatch) error {
	gm := domainToGormReconMatch(m)
	if err := r.db.WithContext(ctx).Create(gm).Error; err != nil {
		return err
	}
	m.ID = gm.ID
	return nil
}

func (r *PGBankRepo) GetReconciliationMatches(ctx context.Context, reconID string) ([]domain.BankReconciliationMatch, error) {
	var models []domain.BankReconciliationMatchGORM
	if err := r.db.WithContext(ctx).Where("reconciliation_id = ?", reconID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.BankReconciliationMatch, len(models))
	for i := range models {
		out[i] = *gormReconMatchToDomain(&models[i])
	}
	if out == nil {
		out = []domain.BankReconciliationMatch{}
	}
	return out, nil
}

func (r *PGBankRepo) DeleteReconciliationMatch(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.BankReconciliationMatchGORM{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrReconciliationNotFound
	}
	return nil
}

// ─── Payment Orders ────────────────────────────────────────────────────

func (r *PGBankRepo) CreatePaymentOrder(ctx context.Context, po *domain.PaymentOrder) error {
	m := domainToGormPaymentOrder(po)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	po.ID = m.ID
	return nil
}

func (r *PGBankRepo) GetPaymentOrder(ctx context.Context, id string) (*domain.PaymentOrder, error) {
	var m domain.PaymentOrderGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrPaymentOrderNotFound
	}
	return gormPaymentOrderToDomain(&m), nil
}

func (r *PGBankRepo) ListPaymentOrders(ctx context.Context, filter domain.PaymentOrderFilter) ([]domain.PaymentOrder, int, error) {
	q := r.db.WithContext(ctx).Model(&domain.PaymentOrderGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.PaymentType != "" {
		q = q.Where("payment_type = ?", string(filter.PaymentType))
	}
	if filter.FromDate != "" {
		q = q.Where("due_date >= ?", filter.FromDate)
	}
	if filter.ToDate != "" {
		q = q.Where("due_date <= ?", filter.ToDate)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
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

	var models []domain.PaymentOrderGORM
	if err := q.Order("due_date DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.PaymentOrder, len(models))
	for i := range models {
		out[i] = *gormPaymentOrderToDomain(&models[i])
	}
	if out == nil {
		out = []domain.PaymentOrder{}
	}
	return out, int(total), nil
}

func (r *PGBankRepo) UpdatePaymentOrder(ctx context.Context, po *domain.PaymentOrder) error {
	m := domainToGormPaymentOrder(po)
	return r.db.WithContext(ctx).Model(&domain.PaymentOrderGORM{}).Where("id = ?", po.ID).Updates(map[string]interface{}{
		"payee_name":    m.PayeeName,
		"payee_bank":    m.PayeeBank,
		"payee_account": m.PayeeAccount,
		"amount":        m.Amount,
		"currency":      m.Currency,
		"status":        m.Status,
		"due_date":      m.DueDate,
		"purpose":       m.Purpose,
		"updated_at":    time.Now(),
	}).Error
}

// ─── Payment Order Batches ─────────────────────────────────────────────

func (r *PGBankRepo) CreatePaymentOrderBatch(ctx context.Context, b *domain.PaymentOrderBatch) error {
	m := domainToGormPaymentBatch(b)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	b.ID = m.ID
	return nil
}

func (r *PGBankRepo) GetPaymentOrderBatch(ctx context.Context, id string) (*domain.PaymentOrderBatch, error) {
	var m domain.PaymentOrderBatchGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrPaymentBatchNotFound
	}
	return gormPaymentBatchToDomain(&m), nil
}

func (r *PGBankRepo) ListPaymentOrderBatches(ctx context.Context, companyID string) ([]domain.PaymentOrderBatch, error) {
	var models []domain.PaymentOrderBatchGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.PaymentOrderBatch, len(models))
	for i := range models {
		out[i] = *gormPaymentBatchToDomain(&models[i])
	}
	if out == nil {
		out = []domain.PaymentOrderBatch{}
	}
	return out, nil
}

func (r *PGBankRepo) UpdatePaymentOrderBatch(ctx context.Context, b *domain.PaymentOrderBatch) error {
	return r.db.WithContext(ctx).Model(&domain.PaymentOrderBatchGORM{}).Where("id = ?", b.ID).Updates(map[string]interface{}{
		"total_amount": b.TotalAmount,
		"order_count":  b.OrderCount,
		"status":       string(b.Status),
		"submitted_at": b.SubmittedAt,
		"bank_ref":     b.BankRef,
		"updated_at":   time.Now(),
	}).Error
}

func (r *PGBankRepo) AddOrdersToBatch(ctx context.Context, batchID string, orderIDs []string) error {
	if len(orderIDs) == 0 {
		return nil
	}
	type batchItem struct {
		ID      string `gorm:"column:id;size:36"`
		BatchID string `gorm:"column:batch_id;size:36"`
		OrderID string `gorm:"column:order_id;size:36"`
	}
	items := make([]batchItem, len(orderIDs))
	for i, oid := range orderIDs {
		items[i] = batchItem{ID: newUUID(), BatchID: batchID, OrderID: oid}
	}
	return r.db.WithContext(ctx).Table("payment_order_batch_items").Create(&items).Error
}

func (r *PGBankRepo) GetBatchOrderIDs(ctx context.Context, batchID string) ([]string, error) {
	type result struct {
		OrderID string
	}
	var rows []result
	if err := r.db.WithContext(ctx).Table("payment_order_batch_items").
		Select("order_id").Where("batch_id = ?", batchID).
		Order("created_at").Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, len(rows))
	for i, r := range rows {
		ids[i] = r.OrderID
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}

// ─── Loans ─────────────────────────────────────────────────────────────

func (r *PGBankRepo) CreateLoan(ctx context.Context, l *domain.LoanAgreement) error {
	m := domainToGormLoan(l)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	l.ID = m.ID
	return nil
}

func (r *PGBankRepo) GetLoan(ctx context.Context, id string) (*domain.LoanAgreement, error) {
	var m domain.LoanAgreementGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrLoanAgreementNotFound
	}
	return gormLoanToDomain(&m), nil
}

func (r *PGBankRepo) ListLoans(ctx context.Context, filter domain.LoanFilter) ([]domain.LoanAgreement, error) {
	q := r.db.WithContext(ctx).Model(&domain.LoanAgreementGORM{}).Where("company_id = ?", filter.CompanyID)
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.LoanType != "" {
		q = q.Where("loan_type = ?", string(filter.LoanType))
	}
	var models []domain.LoanAgreementGORM
	if err := q.Order("start_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.LoanAgreement, len(models))
	for i := range models {
		out[i] = *gormLoanToDomain(&models[i])
	}
	if out == nil {
		out = []domain.LoanAgreement{}
	}
	return out, nil
}

func (r *PGBankRepo) UpdateLoan(ctx context.Context, l *domain.LoanAgreement) error {
	return r.db.WithContext(ctx).Model(&domain.LoanAgreementGORM{}).Where("id = ?", l.ID).Updates(map[string]interface{}{
		"principal":   l.PrincipalAmount,
		"status":      string(l.Status),
		"notes":       l.Notes,
		"updated_at":  time.Now(),
	}).Error
}

// ─── Disbursements ────────────────────────────────────────────────────

func (r *PGBankRepo) CreateDisbursement(ctx context.Context, d *domain.LoanDisbursement) error {
	m := domainToGormDisbursement(d)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	d.ID = m.ID
	return nil
}

func (r *PGBankRepo) GetDisbursements(ctx context.Context, loanID string) ([]domain.LoanDisbursement, error) {
	var models []domain.LoanDisbursementGORM
	if err := r.db.WithContext(ctx).Where("loan_id = ?", loanID).Order("disburse_date").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.LoanDisbursement, len(models))
	for i := range models {
		out[i] = *gormDisbursementToDomain(&models[i])
	}
	if out == nil {
		out = []domain.LoanDisbursement{}
	}
	return out, nil
}

// ─── Repayments ───────────────────────────────────────────────────────

func (r *PGBankRepo) CreateRepayment(ctx context.Context, rp *domain.LoanRepayment) error {
	m := domainToGormRepayment(rp)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	rp.ID = m.ID
	return nil
}

func (r *PGBankRepo) GetRepayments(ctx context.Context, loanID string) ([]domain.LoanRepayment, error) {
	var models []domain.LoanRepaymentGORM
	if err := r.db.WithContext(ctx).Where("loan_id = ?", loanID).Order("repay_date").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.LoanRepayment, len(models))
	for i := range models {
		out[i] = *gormRepaymentToDomain(&models[i])
	}
	if out == nil {
		out = []domain.LoanRepayment{}
	}
	return out, nil
}

func (r *PGBankRepo) UpdateRepayment(ctx context.Context, rp *domain.LoanRepayment) error {
	return r.db.WithContext(ctx).Model(&domain.LoanRepaymentGORM{}).Where("id = ?", rp.ID).Updates(map[string]interface{}{
		"status": string(rp.Status),
		"notes":  rp.Notes,
	}).Error
}

// ─── Term Deposits ────────────────────────────────────────────────────

func (r *PGBankRepo) CreateDeposit(ctx context.Context, d *domain.TermDeposit) error {
	m := domainToGormDeposit(d)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	d.ID = m.ID
	return nil
}

func (r *PGBankRepo) GetDeposit(ctx context.Context, id string) (*domain.TermDeposit, error) {
	var m domain.TermDepositGORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, domain.ErrTermDepositNotFound
	}
	return gormDepositToDomain(&m), nil
}

func (r *PGBankRepo) ListDeposits(ctx context.Context, companyID string) ([]domain.TermDeposit, error) {
	var models []domain.TermDepositGORM
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Order("start_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.TermDeposit, len(models))
	for i := range models {
		out[i] = *gormDepositToDomain(&models[i])
	}
	if out == nil {
		out = []domain.TermDeposit{}
	}
	return out, nil
}

func (r *PGBankRepo) UpdateDeposit(ctx context.Context, d *domain.TermDeposit) error {
	return r.db.WithContext(ctx).Model(&domain.TermDepositGORM{}).Where("id = ?", d.ID).Updates(map[string]interface{}{
		"status":             string(d.Status),
		"notes":              d.Notes,
		"matured_at":         stringToTimePtr(d.MaturedAt),
		"interest_at_maturity": d.InterestAtMaturity,
	}).Error
}

// ─── Reports ──────────────────────────────────────────────────────────

func (r *PGBankRepo) GetBankLedger(ctx context.Context, companyID, bankAccountID, fromDate, toDate string) (*domain.BankLedger, error) {
	ledger := &domain.BankLedger{
		CompanyID:     companyID,
		BankAccountID: bankAccountID,
		FromDate:      fromDate,
		ToDate:        toDate,
	}

	var opening float64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(CASE WHEN txn_type='credit' THEN amount ELSE -amount END), 0)
		FROM (
			SELECT 'credit' AS txn_type, amount FROM payment_orders
			WHERE company_id = ? AND from_bank_acc_id = ? AND payment_date < ? AND status IN ('CONFIRMED','SUBMITTED')
			UNION ALL
			SELECT 'debit' AS txn_type, amount FROM cash_receipts
			WHERE company_id = ? AND cash_account_id = ? AND receipt_date < ? AND status = 'POSTED'
			UNION ALL
			SELECT 'debit' AS txn_type, amount FROM cash_payments
			WHERE company_id = ? AND cash_account_id = ? AND payment_date < ? AND status = 'POSTED'
		) t`, companyID, bankAccountID, fromDate, companyID, bankAccountID, fromDate, companyID, bankAccountID, fromDate).Scan(&opening).Error
	if err != nil {
		return nil, err
	}
	ledger.OpeningBalance = math.Round(opening*100) / 100

	type entryRow struct {
		TransactionDate string
		RefID           string
		Description     string
		DebitAmount     float64
		CreditAmount    float64
		Reference       string
	}
	var rows []entryRow
	query := `
		SELECT transaction_date, reference, description, debit_amount, credit_amount, ref_id
		FROM (
			SELECT payment_date AS transaction_date, id AS ref_id,
				beneficiary_name || ' - ' || payment_content AS description,
				0 AS debit_amount, amount AS credit_amount, '' AS reference
			FROM payment_orders
			WHERE company_id = ? AND from_bank_acc_id = ? AND payment_date >= ? AND payment_date <= ? AND status IN ('CONFIRMED','SUBMITTED')
			UNION ALL
			SELECT receipt_date, id, counterparty_name || ' - ' || reason AS description,
				amount, 0, voucher_no
			FROM cash_receipts
			WHERE company_id = ? AND cash_account_id = ? AND receipt_date >= ? AND receipt_date <= ? AND status = 'POSTED'
			UNION ALL
			SELECT payment_date, id, payee_name || ' - ' || reason AS description,
				0, amount, voucher_no
			FROM cash_payments
			WHERE company_id = ? AND cash_account_id = ? AND payment_date >= ? AND payment_date <= ? AND status = 'POSTED'
		) t ORDER BY transaction_date, reference`
	if err := r.db.WithContext(ctx).Raw(query, companyID, bankAccountID, fromDate, toDate, companyID, bankAccountID, fromDate, toDate, companyID, bankAccountID, fromDate, toDate).Scan(&rows).Error; err != nil {
		return nil, err
	}

	entries := make([]domain.BankLedgerEntry, len(rows))
	var running float64 = ledger.OpeningBalance
	var totalDebits, totalCredits float64
	for i, row := range rows {
		entries[i] = domain.BankLedgerEntry{
			LineNo:          i + 1,
			TransactionDate: row.TransactionDate,
			VoucherNo:       row.Reference,
			Description:     row.Description,
			DebitAmount:     row.DebitAmount,
			CreditAmount:    row.CreditAmount,
			RefID:           row.RefID,
		}
		running += row.DebitAmount - row.CreditAmount
		entries[i].RunningBalance = math.Round(running*100) / 100
		totalDebits += row.DebitAmount
		totalCredits += row.CreditAmount
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
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE((
			SELECT SUM(amount) FROM cash_receipts
			WHERE company_id = ? AND cash_account_id = ? AND status = 'POSTED'
		), 0) - COALESCE((
			SELECT SUM(amount) FROM (
				SELECT amount FROM payment_orders WHERE company_id = ? AND from_bank_acc_id = ? AND status IN ('CONFIRMED','SUBMITTED')
				UNION ALL
				SELECT amount FROM cash_payments WHERE company_id = ? AND cash_account_id = ? AND status = 'POSTED'
			) t
		), 0)`, companyID, bankAccountID, companyID, bankAccountID, companyID, bankAccountID).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return math.Round(total*100) / 100, nil
}
