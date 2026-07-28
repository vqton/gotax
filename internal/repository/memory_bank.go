package repository

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"math"
	"sort"
	"sync"
	"time"
)

var uuidSeed int64

func newUUID() string {
	uuidSeed++
	return fmt.Sprintf("B-%d-%d", time.Now().UnixNano(), uuidSeed)
}

type MemoryBankRepo struct {
	mu                  sync.RWMutex
	statements          map[string]*domain.BankStatement
	statementLines      map[string]*domain.BankStatementLine
	statementLineByStmt map[string][]string
	reconciliations     map[string]*domain.BankReconciliation
	reconMatches        map[string]*domain.BankReconciliationMatch
	reconMatchByRecon   map[string][]string
	paymentOrders       map[string]*domain.PaymentOrder
	paymentBatches      map[string]*domain.PaymentOrderBatch
	batchOrders         map[string][]string
	loans               map[string]*domain.LoanAgreement
	disbursements       map[string]*domain.LoanDisbursement
	disbByLoan          map[string][]string
	repayments          map[string]*domain.LoanRepayment
	repayByLoan         map[string][]string
	deposits            map[string]*domain.TermDeposit
}

func NewMemoryBankRepo() *MemoryBankRepo {
	return &MemoryBankRepo{
		statements:          make(map[string]*domain.BankStatement),
		statementLines:      make(map[string]*domain.BankStatementLine),
		statementLineByStmt: make(map[string][]string),
		reconciliations:     make(map[string]*domain.BankReconciliation),
		reconMatches:        make(map[string]*domain.BankReconciliationMatch),
		reconMatchByRecon:   make(map[string][]string),
		paymentOrders:       make(map[string]*domain.PaymentOrder),
		paymentBatches:      make(map[string]*domain.PaymentOrderBatch),
		batchOrders:         make(map[string][]string),
		loans:               make(map[string]*domain.LoanAgreement),
		disbursements:       make(map[string]*domain.LoanDisbursement),
		disbByLoan:          make(map[string][]string),
		repayments:          make(map[string]*domain.LoanRepayment),
		repayByLoan:         make(map[string][]string),
		deposits:            make(map[string]*domain.TermDeposit),
	}
}

// ─── BankStatement ─────────────────────────────────────────────────

func (r *MemoryBankRepo) CreateStatement(_ context.Context, s *domain.BankStatement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *s
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.statements[cp.ID] = &cp
	s.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetStatement(_ context.Context, id string) (*domain.BankStatement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.statements[id]
	if !ok {
		return nil, domain.ErrBankStatementNotFound
	}
	cp := *s
	return &cp, nil
}

func (r *MemoryBankRepo) ListStatements(_ context.Context, companyID, bankAccountID string, limit, offset int) ([]domain.BankStatement, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var matched []domain.BankStatement
	for _, s := range r.statements {
		if companyID != "" && s.CompanyID != companyID {
			continue
		}
		if bankAccountID != "" && s.BankAccountID != bankAccountID {
			continue
		}
		matched = append(matched, *s)
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].StatementDate > matched[j].StatementDate
	})
	total := len(matched)
	if offset > 0 && offset < len(matched) {
		matched = matched[offset:]
	}
	if limit > 0 && limit < len(matched) {
		matched = matched[:limit]
	}
	return matched, total, nil
}

func (r *MemoryBankRepo) DeleteStatement(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.statements[id]; !ok {
		return domain.ErrBankStatementNotFound
	}
	delete(r.statements, id)
	if ids, ok := r.statementLineByStmt[id]; ok {
		for _, lid := range ids {
			delete(r.statementLines, lid)
		}
		delete(r.statementLineByStmt, id)
	}
	return nil
}

// ─── BankStatementLine ─────────────────────────────────────────────

func (r *MemoryBankRepo) CreateStatementLines(_ context.Context, lines []domain.BankStatementLine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range lines {
		cp := lines[i]
		if cp.ID == "" {
			cp.ID = newUUID()
		}
		lines[i].ID = cp.ID
		r.statementLines[cp.ID] = &cp
		r.statementLineByStmt[cp.StatementID] = append(r.statementLineByStmt[cp.StatementID], cp.ID)
	}
	return nil
}

func (r *MemoryBankRepo) GetStatementLines(_ context.Context, statementID string) ([]domain.BankStatementLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids, ok := r.statementLineByStmt[statementID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.BankStatementLine, 0, len(ids))
	for _, id := range ids {
		if line, ok := r.statementLines[id]; ok {
			out = append(out, *line)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].TransactionDate < out[j].TransactionDate
	})
	return out, nil
}

func (r *MemoryBankRepo) GetStatementLinesByStatus(_ context.Context, statementID string, status domain.MatchStatus) ([]domain.BankStatementLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids, ok := r.statementLineByStmt[statementID]
	if !ok {
		return nil, nil
	}
	var out []domain.BankStatementLine
	for _, id := range ids {
		line, ok := r.statementLines[id]
		if !ok || line.MatchStatus != status {
			continue
		}
		out = append(out, *line)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].TransactionDate < out[j].TransactionDate
	})
	return out, nil
}

func (r *MemoryBankRepo) UpdateStatementLineMatch(_ context.Context, lineID string, matchStatus domain.MatchStatus, matchedLineID, matchedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	line, ok := r.statementLines[lineID]
	if !ok {
		return domain.ErrStatementLineNotFound
	}
	line.MatchStatus = matchStatus
	line.MatchedLineID = matchedLineID
	line.MatchedBy = matchedBy
	line.MatchedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

// ─── BankReconciliation ────────────────────────────────────────────

func (r *MemoryBankRepo) CreateReconciliation(_ context.Context, rc *domain.BankReconciliation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *rc
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.reconciliations[cp.ID] = &cp
	rc.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetReconciliation(_ context.Context, id string) (*domain.BankReconciliation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rc, ok := r.reconciliations[id]
	if !ok {
		return nil, domain.ErrReconciliationNotFound
	}
	cp := *rc
	return &cp, nil
}

func (r *MemoryBankRepo) ListReconciliations(_ context.Context, companyID, bankAccountID string) ([]domain.BankReconciliation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.BankReconciliation
	for _, rc := range r.reconciliations {
		if companyID != "" && rc.CompanyID != companyID {
			continue
		}
		if bankAccountID != "" && rc.BankAccountID != bankAccountID {
			continue
		}
		out = append(out, *rc)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt > out[j].CreatedAt
	})
	return out, nil
}

func (r *MemoryBankRepo) UpdateReconciliation(_ context.Context, rc *domain.BankReconciliation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.reconciliations[rc.ID]; !ok {
		return domain.ErrReconciliationNotFound
	}
	cp := *rc
	r.reconciliations[rc.ID] = &cp
	return nil
}

// ─── BankReconciliationMatch ───────────────────────────────────────

func (r *MemoryBankRepo) CreateReconciliationMatch(_ context.Context, m *domain.BankReconciliationMatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *m
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.reconMatches[cp.ID] = &cp
	r.reconMatchByRecon[cp.ReconciliationID] = append(r.reconMatchByRecon[cp.ReconciliationID], cp.ID)
	m.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetReconciliationMatches(_ context.Context, reconID string) ([]domain.BankReconciliationMatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids, ok := r.reconMatchByRecon[reconID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.BankReconciliationMatch, 0, len(ids))
	for _, id := range ids {
		if m, ok := r.reconMatches[id]; ok {
			out = append(out, *m)
		}
	}
	return out, nil
}

func (r *MemoryBankRepo) DeleteReconciliationMatch(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.reconMatches[id]
	if !ok {
		return domain.ErrReconciliationNotFound
	}
	reconID := m.ReconciliationID
	delete(r.reconMatches, id)
	if ids, ok := r.reconMatchByRecon[reconID]; ok {
		filtered := ids[:0]
		for _, rid := range ids {
			if rid != id {
				filtered = append(filtered, rid)
			}
		}
		r.reconMatchByRecon[reconID] = filtered
	}
	return nil
}

// ─── PaymentOrder ──────────────────────────────────────────────────

func (r *MemoryBankRepo) CreatePaymentOrder(_ context.Context, po *domain.PaymentOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *po
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.paymentOrders[cp.ID] = &cp
	po.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetPaymentOrder(_ context.Context, id string) (*domain.PaymentOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	po, ok := r.paymentOrders[id]
	if !ok {
		return nil, domain.ErrPaymentOrderNotFound
	}
	cp := *po
	return &cp, nil
}

func (r *MemoryBankRepo) ListPaymentOrders(_ context.Context, filter domain.PaymentOrderFilter) ([]domain.PaymentOrder, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var matched []domain.PaymentOrder
	for _, po := range r.paymentOrders {
		if filter.CompanyID != "" && po.CompanyID != filter.CompanyID {
			continue
		}
		if filter.Status != "" && po.Status != filter.Status {
			continue
		}
		if filter.PaymentType != "" && po.PaymentType != filter.PaymentType {
			continue
		}
		if filter.FromDate != "" && po.PaymentDate < filter.FromDate {
			continue
		}
		if filter.ToDate != "" && po.PaymentDate > filter.ToDate {
			continue
		}
		matched = append(matched, *po)
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].PaymentDate < matched[j].PaymentDate
	})
	total := len(matched)
	offset := filter.Offset
	limit := filter.Limit
	if offset > 0 && offset < len(matched) {
		matched = matched[offset:]
	}
	if limit > 0 && limit < len(matched) {
		matched = matched[:limit]
	}
	return matched, total, nil
}

func (r *MemoryBankRepo) UpdatePaymentOrder(_ context.Context, po *domain.PaymentOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.paymentOrders[po.ID]; !ok {
		return domain.ErrPaymentOrderNotFound
	}
	cp := *po
	r.paymentOrders[po.ID] = &cp
	return nil
}

// ─── PaymentOrderBatch ─────────────────────────────────────────────

func (r *MemoryBankRepo) CreatePaymentOrderBatch(_ context.Context, b *domain.PaymentOrderBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *b
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.paymentBatches[cp.ID] = &cp
	b.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetPaymentOrderBatch(_ context.Context, id string) (*domain.PaymentOrderBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.paymentBatches[id]
	if !ok {
		return nil, domain.ErrPaymentBatchNotFound
	}
	cp := *b
	return &cp, nil
}

func (r *MemoryBankRepo) ListPaymentOrderBatches(_ context.Context, companyID string) ([]domain.PaymentOrderBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.PaymentOrderBatch
	for _, b := range r.paymentBatches {
		if companyID != "" && b.CompanyID != companyID {
			continue
		}
		out = append(out, *b)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].BatchDate > out[j].BatchDate
	})
	return out, nil
}

func (r *MemoryBankRepo) UpdatePaymentOrderBatch(_ context.Context, b *domain.PaymentOrderBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.paymentBatches[b.ID]; !ok {
		return domain.ErrPaymentBatchNotFound
	}
	cp := *b
	r.paymentBatches[b.ID] = &cp
	return nil
}

func (r *MemoryBankRepo) AddOrdersToBatch(_ context.Context, batchID string, orderIDs []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.paymentBatches[batchID]; !ok {
		return domain.ErrPaymentBatchNotFound
	}
	r.batchOrders[batchID] = append(r.batchOrders[batchID], orderIDs...)
	return nil
}

func (r *MemoryBankRepo) GetBatchOrderIDs(_ context.Context, batchID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.paymentBatches[batchID]; !ok {
		return nil, domain.ErrPaymentBatchNotFound
	}
	ids := r.batchOrders[batchID]
	out := make([]string, len(ids))
	copy(out, ids)
	return out, nil
}

// ─── LoanAgreement ─────────────────────────────────────────────────

func (r *MemoryBankRepo) CreateLoan(_ context.Context, l *domain.LoanAgreement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *l
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.loans[cp.ID] = &cp
	l.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetLoan(_ context.Context, id string) (*domain.LoanAgreement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.loans[id]
	if !ok {
		return nil, domain.ErrLoanAgreementNotFound
	}
	cp := *l
	return &cp, nil
}

func (r *MemoryBankRepo) ListLoans(_ context.Context, filter domain.LoanFilter) ([]domain.LoanAgreement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.LoanAgreement
	for _, l := range r.loans {
		if filter.CompanyID != "" && l.CompanyID != filter.CompanyID {
			continue
		}
		if filter.Status != "" && l.Status != filter.Status {
			continue
		}
		if filter.LoanType != "" && l.LoanType != filter.LoanType {
			continue
		}
		out = append(out, *l)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartDate > out[j].StartDate
	})
	return out, nil
}

func (r *MemoryBankRepo) UpdateLoan(_ context.Context, l *domain.LoanAgreement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.loans[l.ID]; !ok {
		return domain.ErrLoanAgreementNotFound
	}
	cp := *l
	r.loans[l.ID] = &cp
	return nil
}

// ─── LoanDisbursement ──────────────────────────────────────────────

func (r *MemoryBankRepo) CreateDisbursement(_ context.Context, d *domain.LoanDisbursement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *d
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.disbursements[cp.ID] = &cp
	r.disbByLoan[cp.LoanID] = append(r.disbByLoan[cp.LoanID], cp.ID)
	d.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetDisbursements(_ context.Context, loanID string) ([]domain.LoanDisbursement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids, ok := r.disbByLoan[loanID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.LoanDisbursement, 0, len(ids))
	for _, id := range ids {
		if d, ok := r.disbursements[id]; ok {
			out = append(out, *d)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].DisbursementDate < out[j].DisbursementDate
	})
	return out, nil
}

// ─── LoanRepayment ─────────────────────────────────────────────────

func (r *MemoryBankRepo) CreateRepayment(_ context.Context, rp *domain.LoanRepayment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *rp
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.repayments[cp.ID] = &cp
	r.repayByLoan[cp.LoanID] = append(r.repayByLoan[cp.LoanID], cp.ID)
	rp.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetRepayments(_ context.Context, loanID string) ([]domain.LoanRepayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids, ok := r.repayByLoan[loanID]
	if !ok {
		return nil, nil
	}
	out := make([]domain.LoanRepayment, 0, len(ids))
	for _, id := range ids {
		if rp, ok := r.repayments[id]; ok {
			out = append(out, *rp)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].RepaymentDate < out[j].RepaymentDate
	})
	return out, nil
}

func (r *MemoryBankRepo) UpdateRepayment(_ context.Context, rp *domain.LoanRepayment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.repayments[rp.ID]; !ok {
		return domain.ErrLoanRepaymentNotFound
	}
	cp := *rp
	r.repayments[rp.ID] = &cp
	return nil
}

// ─── TermDeposit ───────────────────────────────────────────────────

func (r *MemoryBankRepo) CreateDeposit(_ context.Context, d *domain.TermDeposit) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *d
	if cp.ID == "" {
		cp.ID = newUUID()
	}
	r.deposits[cp.ID] = &cp
	d.ID = cp.ID
	return nil
}

func (r *MemoryBankRepo) GetDeposit(_ context.Context, id string) (*domain.TermDeposit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.deposits[id]
	if !ok {
		return nil, domain.ErrTermDepositNotFound
	}
	cp := *d
	return &cp, nil
}

func (r *MemoryBankRepo) ListDeposits(_ context.Context, companyID string) ([]domain.TermDeposit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.TermDeposit
	for _, d := range r.deposits {
		if companyID != "" && d.CompanyID != companyID {
			continue
		}
		out = append(out, *d)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartDate > out[j].StartDate
	})
	return out, nil
}

func (r *MemoryBankRepo) UpdateDeposit(_ context.Context, d *domain.TermDeposit) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.deposits[d.ID]; !ok {
		return domain.ErrTermDepositNotFound
	}
	cp := *d
	r.deposits[d.ID] = &cp
	return nil
}

// ─── Reports ───────────────────────────────────────────────────────

func (r *MemoryBankRepo) GetBankLedger(_ context.Context, companyID, bankAccountID, fromDate, toDate string) (*domain.BankLedger, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type txEntry struct {
		date        string
		desc        string
		debit       float64
		credit      float64
		refID       string
		counterpart string
	}

	var allTx []txEntry

	for _, po := range r.paymentOrders {
		if po.CompanyID != companyID || po.FromBankAccID != bankAccountID {
			continue
		}
		desc := po.PaymentContent
		if desc == "" {
			desc = "Payment " + string(po.PaymentType)
		}
		allTx = append(allTx, txEntry{
			date:        po.PaymentDate,
			desc:        desc,
			debit:       po.Amount,
			refID:       po.ID,
			counterpart: po.BeneficiaryName,
		})
	}

	for _, d := range r.disbursements {
		loan, ok := r.loans[d.LoanID]
		if !ok || loan.CompanyID != companyID {
			continue
		}
		if d.ToBankAccountID != bankAccountID && loan.BankAccountID != bankAccountID {
			continue
		}
		desc := d.Notes
		if desc == "" {
			desc = "Loan disbursement - " + loan.ContractNo
		}
		allTx = append(allTx, txEntry{
			date:   d.DisbursementDate,
			desc:   desc,
			credit: d.Amount,
			refID:  d.ID,
		})
	}

	sort.Slice(allTx, func(i, j int) bool {
		return allTx[i].date < allTx[j].date
	})

	var openingBalance float64
	for _, t := range allTx {
		if t.date < fromDate {
			openingBalance += t.credit - t.debit
		}
	}

	var entries []domain.BankLedgerEntry
	var totalDebits, totalCredits float64
	runningBalance := openingBalance
	lineNo := 1
	for _, t := range allTx {
		if t.date < fromDate || t.date > toDate {
			continue
		}
		runningBalance += t.credit - t.debit
		entries = append(entries, domain.BankLedgerEntry{
			LineNo:             lineNo,
			TransactionDate:    t.date,
			Description:        t.desc,
			DebitAmount:        t.debit,
			CreditAmount:       t.credit,
			RunningBalance:     runningBalance,
			RefID:              t.refID,
			CounterpartAccount: t.counterpart,
		})
		lineNo++
		totalDebits += t.debit
		totalCredits += t.credit
	}

	closingBalance := openingBalance + totalCredits - totalDebits

	return &domain.BankLedger{
		CompanyID:      companyID,
		BankAccountID:  bankAccountID,
		FromDate:       fromDate,
		ToDate:         toDate,
		OpeningBalance: openingBalance,
		TotalDebits:    totalDebits,
		TotalCredits:   totalCredits,
		ClosingBalance: closingBalance,
		Entries:        entries,
	}, nil
}

func (r *MemoryBankRepo) GetBalance(_ context.Context, companyID, bankAccountID string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var balance float64

	for _, po := range r.paymentOrders {
		if po.CompanyID == companyID && po.FromBankAccID == bankAccountID {
			balance -= po.Amount
		}
	}

	for _, d := range r.disbursements {
		loan, ok := r.loans[d.LoanID]
		if !ok || loan.CompanyID != companyID {
			continue
		}
		if d.ToBankAccountID != bankAccountID && loan.BankAccountID != bankAccountID {
			continue
		}
		balance += d.Amount
	}

	return math.Round(balance*100) / 100, nil
}
