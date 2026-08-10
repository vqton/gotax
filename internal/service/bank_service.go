package service

import (
	"context"
	"fmt"
	"gotax/internal/domain"
	"time"
)

type BankService struct {
	bankRepo domain.BankRepository
}

func NewBankService(bankRepo domain.BankRepository) *BankService {
	return &BankService{bankRepo: bankRepo}
}

// ─── Statements ────────────────────────────────────────────────────────

func (s *BankService) ImportStatement(ctx context.Context, stmt *domain.BankStatement, lines []domain.BankStatementLine) error {
	stmt.ID = newID()
	for i := range lines {
		lines[i].ID = newID()
		lines[i].StatementID = stmt.ID
		lines[i].MatchStatus = domain.MatchPending
	}
	if err := s.bankRepo.CreateStatement(ctx, stmt); err != nil {
		return err
	}
	return s.bankRepo.CreateStatementLines(ctx, lines)
}

func (s *BankService) GetStatement(ctx context.Context, id string) (*domain.BankStatement, error) {
	return s.bankRepo.GetStatement(ctx, id)
}

func (s *BankService) ListStatements(ctx context.Context, companyID, bankAccountID string, limit, offset int) ([]domain.BankStatement, int, error) {
	return s.bankRepo.ListStatements(ctx, companyID, bankAccountID, limit, offset)
}

func (s *BankService) DeleteStatement(ctx context.Context, id string) error {
	return s.bankRepo.DeleteStatement(ctx, id)
}

func (s *BankService) GetStatementLines(ctx context.Context, statementID string) ([]domain.BankStatementLine, error) {
	return s.bankRepo.GetStatementLines(ctx, statementID)
}

// ─── Reconciliation ─────────────────────────────────────────────────────

func (s *BankService) StartReconciliation(ctx context.Context, rc *domain.BankReconciliation) error {
	rc.ID = newID()
	rc.Status = domain.ReconInProgress
	return s.bankRepo.CreateReconciliation(ctx, rc)
}

func (s *BankService) GetReconciliation(ctx context.Context, id string) (*domain.BankReconciliation, error) {
	return s.bankRepo.GetReconciliation(ctx, id)
}

func (s *BankService) ListReconciliations(ctx context.Context, companyID, bankAccountID string) ([]domain.BankReconciliation, error) {
	return s.bankRepo.ListReconciliations(ctx, companyID, bankAccountID)
}

func (s *BankService) CompleteReconciliation(ctx context.Context, id, completedBy string) error {
	rc, err := s.bankRepo.GetReconciliation(ctx, id)
	if err != nil {
		return err
	}
	if rc.Status == domain.ReconCompleted {
		return domain.ErrReconAlreadyCompleted
	}
	rc.Difference = rc.ClosingBalance - rc.StatementBalance
	if rc.Difference != 0 {
		return domain.ErrReconDifferenceNotZero
	}
	rc.Status = domain.ReconCompleted
	rc.CompletedBy = completedBy
	rc.CompletedAt = time.Now().UTC().Format(time.RFC3339)
	return s.bankRepo.UpdateReconciliation(ctx, rc)
}

func (s *BankService) AddMatch(ctx context.Context, m *domain.BankReconciliationMatch) error {
	m.ID = newID()
	return s.bankRepo.CreateReconciliationMatch(ctx, m)
}

func (s *BankService) GetMatches(ctx context.Context, reconID string) ([]domain.BankReconciliationMatch, error) {
	return s.bankRepo.GetReconciliationMatches(ctx, reconID)
}

func (s *BankService) RemoveMatch(ctx context.Context, id string) error {
	return s.bankRepo.DeleteReconciliationMatch(ctx, id)
}

// ─── Payment Orders ────────────────────────────────────────────────────

func (s *BankService) CreatePaymentOrder(ctx context.Context, po *domain.PaymentOrder) error {
	if err := po.Validate(); err != nil {
		return err
	}
	po.ID = newID()
	po.Status = domain.PODraft
	return s.bankRepo.CreatePaymentOrder(ctx, po)
}

func (s *BankService) GetPaymentOrder(ctx context.Context, id string) (*domain.PaymentOrder, error) {
	return s.bankRepo.GetPaymentOrder(ctx, id)
}

func (s *BankService) ListPaymentOrders(ctx context.Context, filter domain.PaymentOrderFilter) ([]domain.PaymentOrder, int, error) {
	return s.bankRepo.ListPaymentOrders(ctx, filter)
}

func (s *BankService) UpdatePaymentOrder(ctx context.Context, po *domain.PaymentOrder) error {
	existing, err := s.bankRepo.GetPaymentOrder(ctx, po.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.PODraft {
		return domain.ErrPaymentOrderNotDraft
	}
	if err := po.Validate(); err != nil {
		return err
	}
	return s.bankRepo.UpdatePaymentOrder(ctx, po)
}

func (s *BankService) SubmitPaymentOrder(ctx context.Context, id string) error {
	po, err := s.bankRepo.GetPaymentOrder(ctx, id)
	if err != nil {
		return err
	}
	if !po.Status.ValidTransition(domain.POPendingApproval) {
		return domain.ErrPaymentOrderNotDraft
	}
	po.Status = domain.POPendingApproval
	return s.bankRepo.UpdatePaymentOrder(ctx, po)
}

func (s *BankService) ApprovePaymentOrder(ctx context.Context, id, approvedBy string) error {
	po, err := s.bankRepo.GetPaymentOrder(ctx, id)
	if err != nil {
		return err
	}
	if po.CreatedBy == approvedBy {
		return domain.ErrCannotSelfApprovePayment
	}
	if !po.Status.ValidTransition(domain.POApproved) {
		return domain.ErrPaymentOrderNotDraft
	}
	po.Status = domain.POApproved
	po.ApprovedBy = approvedBy
	po.ApprovedAt = time.Now().UTC().Format(time.RFC3339)
	return s.bankRepo.UpdatePaymentOrder(ctx, po)
}

func (s *BankService) RejectPaymentOrder(ctx context.Context, id, reason string) error {
	po, err := s.bankRepo.GetPaymentOrder(ctx, id)
	if err != nil {
		return err
	}
	if !po.Status.ValidTransition(domain.PORejected) {
		return domain.ErrPaymentOrderNotDraft
	}
	po.Status = domain.PORejected
	po.FailureReason = reason
	return s.bankRepo.UpdatePaymentOrder(ctx, po)
}

// ─── Payment Batches ────────────────────────────────────────────────────

func (s *BankService) CreateBatch(ctx context.Context, b *domain.PaymentOrderBatch, orderIDs []string) error {
	b.ID = newID()
	b.Status = domain.BatchDraft
	if err := s.bankRepo.CreatePaymentOrderBatch(ctx, b); err != nil {
		return err
	}
	return s.bankRepo.AddOrdersToBatch(ctx, b.ID, orderIDs)
}

func (s *BankService) GetBatch(ctx context.Context, id string) (*domain.PaymentOrderBatch, error) {
	return s.bankRepo.GetPaymentOrderBatch(ctx, id)
}

func (s *BankService) ListBatches(ctx context.Context, companyID string) ([]domain.PaymentOrderBatch, error) {
	return s.bankRepo.ListPaymentOrderBatches(ctx, companyID)
}

func (s *BankService) SubmitBatch(ctx context.Context, id string) error {
	b, err := s.bankRepo.GetPaymentOrderBatch(ctx, id)
	if err != nil {
		return err
	}
	if b.Status != domain.BatchDraft {
		return domain.ErrBatchAlreadySubmitted
	}
	b.Status = domain.BatchSubmitted
	b.SubmittedAt = time.Now().UTC().Format(time.RFC3339)
	return s.bankRepo.UpdatePaymentOrderBatch(ctx, b)
}

func (s *BankService) GetBatchOrders(ctx context.Context, batchID string) ([]string, error) {
	return s.bankRepo.GetBatchOrderIDs(ctx, batchID)
}

// ─── Loans ─────────────────────────────────────────────────────────────

func (s *BankService) CreateLoan(ctx context.Context, l *domain.LoanAgreement) error {
	if err := l.Validate(); err != nil {
		return err
	}
	l.ID = newID()
	l.Status = domain.LoanActive
	l.OutstandingBalance = l.PrincipalAmount
	return s.bankRepo.CreateLoan(ctx, l)
}

func (s *BankService) GetLoan(ctx context.Context, id string) (*domain.LoanAgreement, error) {
	return s.bankRepo.GetLoan(ctx, id)
}

func (s *BankService) ListLoans(ctx context.Context, filter domain.LoanFilter) ([]domain.LoanAgreement, error) {
	return s.bankRepo.ListLoans(ctx, filter)
}

func (s *BankService) DisburseLoan(ctx context.Context, disbursement *domain.LoanDisbursement) error {
	loan, err := s.bankRepo.GetLoan(ctx, disbursement.LoanID)
	if err != nil {
		return err
	}
	newTotal := loan.DisbursedAmount + disbursement.Amount
	if newTotal > loan.PrincipalAmount {
		return domain.ErrLoanDisbursementOverLimit
	}
	disbursement.ID = newID()
	if err := s.bankRepo.CreateDisbursement(ctx, disbursement); err != nil {
		return err
	}
	loan.DisbursedAmount = newTotal
	return s.bankRepo.UpdateLoan(ctx, loan)
}

func (s *BankService) GetDisbursements(ctx context.Context, loanID string) ([]domain.LoanDisbursement, error) {
	return s.bankRepo.GetDisbursements(ctx, loanID)
}

func (s *BankService) MakeRepayment(ctx context.Context, rp *domain.LoanRepayment) error {
	loan, err := s.bankRepo.GetLoan(ctx, rp.LoanID)
	if err != nil {
		return err
	}
	rp.ID = newID()
	rp.Status = domain.RepayPaid
	if err := s.bankRepo.CreateRepayment(ctx, rp); err != nil {
		return err
	}
	loan.OutstandingBalance -= rp.PrincipalAmount
	if loan.OutstandingBalance <= 0 {
		loan.OutstandingBalance = 0
		loan.Status = domain.LoanFullyPaid
	}
	return s.bankRepo.UpdateLoan(ctx, loan)
}

func (s *BankService) GetRepayments(ctx context.Context, loanID string) ([]domain.LoanRepayment, error) {
	return s.bankRepo.GetRepayments(ctx, loanID)
}

// ─── Term Deposits ─────────────────────────────────────────────────────

func (s *BankService) CreateDeposit(ctx context.Context, d *domain.TermDeposit) error {
	if err := d.Validate(); err != nil {
		return err
	}
	d.ID = newID()
	d.Status = domain.DepositActive
	return s.bankRepo.CreateDeposit(ctx, d)
}

func (s *BankService) GetDeposit(ctx context.Context, id string) (*domain.TermDeposit, error) {
	return s.bankRepo.GetDeposit(ctx, id)
}

func (s *BankService) ListDeposits(ctx context.Context, companyID string) ([]domain.TermDeposit, error) {
	return s.bankRepo.ListDeposits(ctx, companyID)
}

func (s *BankService) MatureDeposit(ctx context.Context, id string) error {
	d, err := s.bankRepo.GetDeposit(ctx, id)
	if err != nil {
		return err
	}
	if d.Status == domain.DepositMatured {
		return domain.ErrDepositAlreadyMatured
	}
	d.Status = domain.DepositMatured
	d.MaturedAt = time.Now().UTC().Format(time.RFC3339)
	return s.bankRepo.UpdateDeposit(ctx, d)
}

// ─── Reports ─────────────────────────────────────────────────────────

func (s *BankService) GetBankLedger(ctx context.Context, companyID, bankAccountID, fromDate, toDate string) (*domain.BankLedger, error) {
	return s.bankRepo.GetBankLedger(ctx, companyID, bankAccountID, fromDate, toDate)
}

func (s *BankService) GetBalance(ctx context.Context, companyID, bankAccountID string) (float64, error) {
	return s.bankRepo.GetBalance(ctx, companyID, bankAccountID)
}

// ─── Bank API Client Interface ─────────────────────────────────────

// BankTransaction represents a transaction fetched from a bank API.
type BankTransaction struct {
	TransactionID string
	Date          string
	Description   string
	Amount        float64
	Balance       float64
	Reference     string
	Counterparty  string
	AccountNumber string
}

// BankAPIClient is the interface for real-time bank API adapters.
type BankAPIClient interface {
	Name() string
	FetchTransactions(ctx context.Context, accountNumber, fromDate, toDate string) ([]BankTransaction, error)
}

// ─── Sync + Auto-Match ─────────────────────────────────────────────

func (s *BankService) SyncTransactions(ctx context.Context, companyID, bankAccountID, fromDate, toDate string, client BankAPIClient) (int, error) {
	txns, err := client.FetchTransactions(ctx, bankAccountID, fromDate, toDate)
	if err != nil {
		return 0, err
	}
	stmt := &domain.BankStatement{
		CompanyID:     companyID,
		BankAccountID: bankAccountID,
		StatementDate: toDate,
		FromDate:      fromDate,
		ToDate:        toDate,
		ImportMethod:  "API:" + client.Name(),
	}
	var lines []domain.BankStatementLine
	for _, t := range txns {
		credit, debit := 0.0, 0.0
		if t.Amount >= 0 {
			credit = t.Amount
		} else {
			debit = -t.Amount
		}
		lines = append(lines, domain.BankStatementLine{
			TransactionDate: t.Date,
			Description:     t.Description,
			DebitAmount:     debit,
			CreditAmount:    credit,
			ReferenceNo:     t.Reference,
			Counterparty:    t.Counterparty,
			BalanceAfter:    t.Balance,
		})
	}
	if err := s.ImportStatement(ctx, stmt, lines); err != nil {
		return 0, err
	}
	return len(lines), nil
}

// MatchResult represents the result of an auto-match operation.
type MatchResult struct {
	StatementLineID string
	Confidence      float64 // 0.0 - 1.0
	MatchType       string  // "EXACT", "FUZZY_AMOUNT", "FUZZY_REFERENCE"
}

// AutoMatch matches pending statement lines by reference + amount heuristics.
func (s *BankService) AutoMatch(ctx context.Context, statementID string) ([]MatchResult, error) {
	lines, err := s.bankRepo.GetStatementLines(ctx, statementID)
	if err != nil {
		return nil, err
	}
	var results []MatchResult
	for _, line := range lines {
		if line.MatchStatus != domain.MatchPending {
			continue
		}
		if line.ReferenceNo != "" && line.Counterparty != "" {
			results = append(results, MatchResult{
				StatementLineID: line.ID,
				Confidence:      0.8,
				MatchType:       "FUZZY_REFERENCE",
			})
		} else if line.DebitAmount > 0 || line.CreditAmount > 0 {
			results = append(results, MatchResult{
				StatementLineID: line.ID,
				Confidence:      0.6,
				MatchType:       "AMOUNT_ONLY",
			})
		}
	}
	return results, nil
}

func newID() string {
	return fmt.Sprintf("B-%d", time.Now().UnixNano())
}
