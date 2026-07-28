package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupBankService(t *testing.T) (*BankService, context.Context) {
	t.Helper()
	bankRepo := repository.NewMemoryBankRepo()
	svc := NewBankService(bankRepo)
	return svc, context.Background()
}

// ─── Statements ────────────────────────────────────────────────────────

func TestImportStatement_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	stmt := &domain.BankStatement{
		CompanyID: "CMP001", BankAccountID: "BA001",
		StatementDate: "2026-07", OpeningBalance: 0, ClosingBalance: 1000000,
		Currency: "VND", Status: domain.BankStatementImported,
	}
	lines := []domain.BankStatementLine{
		{TransactionDate: "2026-07-01", Description: "payment received", DebitAmount: 0, CreditAmount: 1000000},
	}
	err := svc.ImportStatement(ctx, stmt, lines)
	require.NoError(t, err)
	assert.NotEmpty(t, stmt.ID)
}

func TestGetStatement_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	stmt := &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-07", Currency: "VND"}
	require.NoError(t, svc.ImportStatement(ctx, stmt, nil))

	got, err := svc.GetStatement(ctx, stmt.ID)
	require.NoError(t, err)
	assert.Equal(t, "CMP001", got.CompanyID)
}

func TestGetStatement_NotFound(t *testing.T) {
	svc, ctx := setupBankService(t)
	_, err := svc.GetStatement(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrBankStatementNotFound)
}

func TestListStatements(t *testing.T) {
	svc, ctx := setupBankService(t)
	svc.ImportStatement(ctx, &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-07", Currency: "VND"}, nil)
	svc.ImportStatement(ctx, &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-06", Currency: "VND"}, nil)

	items, total, err := svc.ListStatements(ctx, "CMP001", "BA001", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, 2, len(items))
}

func TestDeleteStatement_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	stmt := &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-07", Currency: "VND", Status: domain.BankStatementImported}
	require.NoError(t, svc.ImportStatement(ctx, stmt, nil))

	err := svc.DeleteStatement(ctx, stmt.ID)
	assert.NoError(t, err)
}

func TestDeleteStatement_NotFound(t *testing.T) {
	svc, ctx := setupBankService(t)
	err := svc.DeleteStatement(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestGetStatementLines(t *testing.T) {
	svc, ctx := setupBankService(t)
	stmt := &domain.BankStatement{CompanyID: "CMP001", BankAccountID: "BA001", StatementDate: "2026-07", Currency: "VND"}
	lines := []domain.BankStatementLine{
		{TransactionDate: "2026-07-01", Description: "line1", DebitAmount: 0, CreditAmount: 500000},
		{TransactionDate: "2026-07-02", Description: "line2", DebitAmount: 0, CreditAmount: 300000},
	}
	require.NoError(t, svc.ImportStatement(ctx, stmt, lines))

	got, err := svc.GetStatementLines(ctx, stmt.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, len(got))
}

// ─── Payment Orders ────────────────────────────────────────────────────

func TestCreatePaymentOrder_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	po := &domain.PaymentOrder{
		CompanyID: "CMP001", PaymentDate: "2026-07-15", Amount: 5000000,
		Currency: "VND", BeneficiaryName: "ABC Corp", BeneficiaryAcc: "123456789",
		BeneficiaryBank: "VCB", FromBankAccID: "BA001", PaymentContent: "payment",
	}
	err := svc.CreatePaymentOrder(ctx, po)
	require.NoError(t, err)
	assert.NotEmpty(t, po.ID)
	assert.Equal(t, domain.PODraft, po.Status)
}

func TestCreatePaymentOrder_Validation(t *testing.T) {
	svc, ctx := setupBankService(t)
	po := &domain.PaymentOrder{}
	err := svc.CreatePaymentOrder(ctx, po)
	assert.Error(t, err)
}

func TestGetPaymentOrder_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	po := &domain.PaymentOrder{CompanyID: "CMP001", PaymentDate: "2026-07-15", Amount: 5000000, Currency: "VND", BeneficiaryName: "ABC Corp", BeneficiaryAcc: "123456789", BeneficiaryBank: "VCB", FromBankAccID: "BA001"}
	require.NoError(t, svc.CreatePaymentOrder(ctx, po))

	got, err := svc.GetPaymentOrder(ctx, po.ID)
	require.NoError(t, err)
	assert.Equal(t, float64(5000000), got.Amount)
}

func TestPaymentOrder_Workflow(t *testing.T) {
	svc, ctx := setupBankService(t)
	po := &domain.PaymentOrder{CompanyID: "CMP001", PaymentDate: "2026-07-15", Amount: 5000000, Currency: "VND", BeneficiaryName: "ABC Corp", BeneficiaryAcc: "123456789", BeneficiaryBank: "VCB", FromBankAccID: "BA001", PaymentContent: "test", CreatedBy: "user1"}
	require.NoError(t, svc.CreatePaymentOrder(ctx, po))

	err := svc.SubmitPaymentOrder(ctx, po.ID)
	assert.NoError(t, err)

	err = svc.ApprovePaymentOrder(ctx, po.ID, "approver1")
	assert.NoError(t, err)

	got, _ := svc.GetPaymentOrder(ctx, po.ID)
	assert.Equal(t, domain.POApproved, got.Status)
}

func TestPaymentOrder_CannotSelfApprove(t *testing.T) {
	svc, ctx := setupBankService(t)
	po := &domain.PaymentOrder{CompanyID: "CMP001", PaymentDate: "2026-07-15", Amount: 5000000, Currency: "VND", BeneficiaryName: "ABC Corp", BeneficiaryAcc: "123456789", BeneficiaryBank: "VCB", FromBankAccID: "BA001", CreatedBy: "user1"}
	require.NoError(t, svc.CreatePaymentOrder(ctx, po))
	require.NoError(t, svc.SubmitPaymentOrder(ctx, po.ID))

	err := svc.ApprovePaymentOrder(ctx, po.ID, "user1")
	assert.ErrorIs(t, err, domain.ErrCannotSelfApprovePayment)
}

func TestPaymentOrder_Reject(t *testing.T) {
	svc, ctx := setupBankService(t)
	po := &domain.PaymentOrder{CompanyID: "CMP001", PaymentDate: "2026-07-15", Amount: 5000000, Currency: "VND", BeneficiaryName: "ABC Corp", BeneficiaryAcc: "123456789", BeneficiaryBank: "VCB", FromBankAccID: "BA001", CreatedBy: "user1"}
	require.NoError(t, svc.CreatePaymentOrder(ctx, po))
	require.NoError(t, svc.SubmitPaymentOrder(ctx, po.ID))

	err := svc.RejectPaymentOrder(ctx, po.ID, "insufficient funds")
	assert.NoError(t, err)

	got, _ := svc.GetPaymentOrder(ctx, po.ID)
	assert.Equal(t, domain.PORejected, got.Status)
}

func TestListPaymentOrders(t *testing.T) {
	svc, ctx := setupBankService(t)
	po1 := &domain.PaymentOrder{CompanyID: "CMP001", PaymentDate: "2026-07-15", Amount: 5000000, Currency: "VND", BeneficiaryName: "A Corp", BeneficiaryAcc: "111", BeneficiaryBank: "VCB", FromBankAccID: "BA001"}
	po2 := &domain.PaymentOrder{CompanyID: "CMP001", PaymentDate: "2026-07-16", Amount: 3000000, Currency: "VND", BeneficiaryName: "B Corp", BeneficiaryAcc: "222", BeneficiaryBank: "VCB", FromBankAccID: "BA001"}
	require.NoError(t, svc.CreatePaymentOrder(ctx, po1))
	require.NoError(t, svc.CreatePaymentOrder(ctx, po2))

	items, total, err := svc.ListPaymentOrders(ctx, domain.PaymentOrderFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, 2, len(items))
}

// ─── Batches ───────────────────────────────────────────────────────────

func TestCreateBatch_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	b := &domain.PaymentOrderBatch{CompanyID: "CMP001", BatchName: "BATCH-001", BatchDate: "2026-07-15", TotalAmount: 5000000, Currency: "VND", OrderCount: 2}
	err := svc.CreateBatch(ctx, b, []string{"PO1", "PO2"})
	require.NoError(t, err)
	assert.NotEmpty(t, b.ID)
}

func TestSubmitBatch(t *testing.T) {
	svc, ctx := setupBankService(t)
	b := &domain.PaymentOrderBatch{CompanyID: "CMP001", BatchName: "BATCH-001", TotalAmount: 5000000, Currency: "VND"}
	require.NoError(t, svc.CreateBatch(ctx, b, nil))

	err := svc.SubmitBatch(ctx, b.ID)
	assert.NoError(t, err)

	got, _ := svc.GetBatch(ctx, b.ID)
	assert.Equal(t, domain.BatchSubmitted, got.Status)
}

// ─── Loans ─────────────────────────────────────────────────────────────

func TestCreateLoan_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	l := &domain.LoanAgreement{
		CompanyID: "CMP001", BankAccountID: "BA001", ContractNo: "LN-001",
		PrincipalAmount: 100000000, Currency: "VND", InterestRate: 8.5,
		StartDate: "2026-01-01", MaturityDate: "2026-12-31",
	}
	err := svc.CreateLoan(ctx, l)
	require.NoError(t, err)
	assert.NotEmpty(t, l.ID)
	assert.Equal(t, domain.LoanActive, l.Status)
	assert.Equal(t, float64(100000000), l.OutstandingBalance)
}

func TestCreateLoan_Validation(t *testing.T) {
	svc, ctx := setupBankService(t)
	l := &domain.LoanAgreement{}
	err := svc.CreateLoan(ctx, l)
	assert.Error(t, err)
}

func TestGetLoan_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	l := &domain.LoanAgreement{CompanyID: "CMP001", BankAccountID: "BA001", ContractNo: "LN-001", PrincipalAmount: 100000000, Currency: "VND", InterestRate: 8.5, StartDate: "2026-01-01", MaturityDate: "2026-12-31"}
	require.NoError(t, svc.CreateLoan(ctx, l))

	got, err := svc.GetLoan(ctx, l.ID)
	require.NoError(t, err)
	assert.Equal(t, "LN-001", got.ContractNo)
}

func TestDisburseLoan_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	l := &domain.LoanAgreement{CompanyID: "CMP001", BankAccountID: "BA001", ContractNo: "LN-001", PrincipalAmount: 100000000, Currency: "VND", InterestRate: 8.5, StartDate: "2026-01-01", MaturityDate: "2026-12-31"}
	require.NoError(t, svc.CreateLoan(ctx, l))

	d := &domain.LoanDisbursement{LoanID: l.ID, DisbursementDate: "2026-01-15", Amount: 50000000}
	err := svc.DisburseLoan(ctx, d)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.ID)
}

func TestDisburseLoan_OverLimit(t *testing.T) {
	svc, ctx := setupBankService(t)
	l := &domain.LoanAgreement{CompanyID: "CMP001", BankAccountID: "BA001", ContractNo: "LN-001", PrincipalAmount: 100000000, Currency: "VND", InterestRate: 8.5, StartDate: "2026-01-01", MaturityDate: "2026-12-31"}
	require.NoError(t, svc.CreateLoan(ctx, l))

	d := &domain.LoanDisbursement{LoanID: l.ID, DisbursementDate: "2026-01-15", Amount: 150000000}
	err := svc.DisburseLoan(ctx, d)
	assert.ErrorIs(t, err, domain.ErrLoanDisbursementOverLimit)
}

func TestRepayLoan_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	l := &domain.LoanAgreement{CompanyID: "CMP001", BankAccountID: "BA001", ContractNo: "LN-001", PrincipalAmount: 100000000, Currency: "VND", InterestRate: 8.5, StartDate: "2026-01-01", MaturityDate: "2026-12-31"}
	require.NoError(t, svc.CreateLoan(ctx, l))

	rp := &domain.LoanRepayment{LoanID: l.ID, RepaymentDate: "2026-02-15", PrincipalAmount: 50000000, InterestAmount: 708333, TotalAmount: 50708333}
	err := svc.MakeRepayment(ctx, rp)
	assert.NoError(t, err)

	got, _ := svc.GetLoan(ctx, l.ID)
	assert.Equal(t, float64(50000000), got.OutstandingBalance)
}

func TestRepayLoan_FullyPaid(t *testing.T) {
	svc, ctx := setupBankService(t)
	l := &domain.LoanAgreement{CompanyID: "CMP001", BankAccountID: "BA001", ContractNo: "LN-001", PrincipalAmount: 100000000, Currency: "VND", InterestRate: 8.5, StartDate: "2026-01-01", MaturityDate: "2026-12-31"}
	require.NoError(t, svc.CreateLoan(ctx, l))

	rp := &domain.LoanRepayment{LoanID: l.ID, RepaymentDate: "2026-02-15", PrincipalAmount: 100000000, InterestAmount: 1416667, TotalAmount: 101416667}
	err := svc.MakeRepayment(ctx, rp)
	assert.NoError(t, err)

	got, _ := svc.GetLoan(ctx, l.ID)
	assert.Equal(t, float64(0), got.OutstandingBalance)
	assert.Equal(t, domain.LoanFullyPaid, got.Status)
}

// ─── Term Deposits ─────────────────────────────────────────────────────

func TestCreateDeposit_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	d := &domain.TermDeposit{
		CompanyID: "CMP001", BankAccountID: "BA001", DepositNo: "TD-001",
		Amount: 50000000, Currency: "VND", InterestRate: 6.5, TermDays: 90,
		StartDate: "2026-07-01", MaturityDate: "2026-09-28",
	}
	err := svc.CreateDeposit(ctx, d)
	require.NoError(t, err)
	assert.NotEmpty(t, d.ID)
	assert.Equal(t, domain.DepositActive, d.Status)
}

func TestCreateDeposit_MinTerm(t *testing.T) {
	svc, ctx := setupBankService(t)
	d := &domain.TermDeposit{
		CompanyID: "CMP001", BankAccountID: "BA001", DepositNo: "TD-002",
		Amount: 50000000, Currency: "VND", InterestRate: 6.5, TermDays: 3,
		StartDate: "2026-07-01", MaturityDate: "2026-07-04",
	}
	err := svc.CreateDeposit(ctx, d)
	assert.ErrorIs(t, err, domain.ErrDepositMinTerm)
}

func TestGetDeposit_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	d := &domain.TermDeposit{CompanyID: "CMP001", BankAccountID: "BA001", DepositNo: "TD-001", Amount: 50000000, Currency: "VND", InterestRate: 6.5, TermDays: 90, StartDate: "2026-07-01", MaturityDate: "2026-09-28"}
	require.NoError(t, svc.CreateDeposit(ctx, d))

	got, err := svc.GetDeposit(ctx, d.ID)
	require.NoError(t, err)
	assert.Equal(t, float64(50000000), got.Amount)
}

func TestMatureDeposit_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	d := &domain.TermDeposit{CompanyID: "CMP001", BankAccountID: "BA001", DepositNo: "TD-001", Amount: 50000000, Currency: "VND", InterestRate: 6.5, TermDays: 90, StartDate: "2026-07-01", MaturityDate: "2026-09-28"}
	require.NoError(t, svc.CreateDeposit(ctx, d))

	err := svc.MatureDeposit(ctx, d.ID)
	assert.NoError(t, err)

	got, _ := svc.GetDeposit(ctx, d.ID)
	assert.Equal(t, domain.DepositMatured, got.Status)
}

func TestMatureDeposit_AlreadyMatured(t *testing.T) {
	svc, ctx := setupBankService(t)
	d := &domain.TermDeposit{CompanyID: "CMP001", BankAccountID: "BA001", DepositNo: "TD-001", Amount: 50000000, Currency: "VND", InterestRate: 6.5, TermDays: 90, StartDate: "2026-07-01", MaturityDate: "2026-09-28"}
	require.NoError(t, svc.CreateDeposit(ctx, d))
	require.NoError(t, svc.MatureDeposit(ctx, d.ID))

	err := svc.MatureDeposit(ctx, d.ID)
	assert.ErrorIs(t, err, domain.ErrDepositAlreadyMatured)
}

// ─── Reconciliation ────────────────────────────────────────────────────

func TestStartReconciliation_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	rc := &domain.BankReconciliation{
		CompanyID: "CMP001", BankAccountID: "BA001", StatementID: "STMT001",
		FromDate: "2026-07-01", ToDate: "2026-07-31",
		OpeningBalance: 1000000, ClosingBalance: 2000000, StatementBalance: 2000000,
	}
	err := svc.StartReconciliation(ctx, rc)
	require.NoError(t, err)
	assert.NotEmpty(t, rc.ID)
	assert.Equal(t, domain.ReconInProgress, rc.Status)
}

func TestCompleteReconciliation_Success(t *testing.T) {
	svc, ctx := setupBankService(t)
	rc := &domain.BankReconciliation{
		CompanyID: "CMP001", BankAccountID: "BA001",
		FromDate: "2026-07-01", ToDate: "2026-07-31",
		OpeningBalance: 1000000, ClosingBalance: 2000000, StatementBalance: 2000000,
		Difference: 0,
	}
	require.NoError(t, svc.StartReconciliation(ctx, rc))

	err := svc.CompleteReconciliation(ctx, rc.ID, "approver1")
	assert.NoError(t, err)

	got, _ := svc.GetReconciliation(ctx, rc.ID)
	assert.Equal(t, domain.ReconCompleted, got.Status)
}

func TestCompleteReconciliation_DifferenceNotZero(t *testing.T) {
	svc, ctx := setupBankService(t)
	rc := &domain.BankReconciliation{
		CompanyID: "CMP001", BankAccountID: "BA001",
		FromDate: "2026-07-01", ToDate: "2026-07-31",
		OpeningBalance: 1000000, ClosingBalance: 2000000, StatementBalance: 2100000,
		Difference: 100000,
	}
	require.NoError(t, svc.StartReconciliation(ctx, rc))

	err := svc.CompleteReconciliation(ctx, rc.ID, "approver1")
	assert.ErrorIs(t, err, domain.ErrReconDifferenceNotZero)
}

func TestAddMatch(t *testing.T) {
	svc, ctx := setupBankService(t)
	rc := &domain.BankReconciliation{CompanyID: "CMP001", BankAccountID: "BA001", FromDate: "2026-07-01", ToDate: "2026-07-31"}
	require.NoError(t, svc.StartReconciliation(ctx, rc))

	m := &domain.BankReconciliationMatch{
		ReconciliationID: rc.ID, StatementLineID: "SL001",
		TransactionID: "JE001", MatchMethod: "AUTO", Confidence: 1.0,
	}
	err := svc.AddMatch(ctx, m)
	assert.NoError(t, err)

	matches, err := svc.GetMatches(ctx, rc.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, len(matches))
}

// ─── Reports ─────────────────────────────────────────────────────────

func TestGetBankLedger(t *testing.T) {
	svc, ctx := setupBankService(t)
	ledger, err := svc.GetBankLedger(ctx, "CMP001", "BA001", "2026-01-01", "2026-12-31")
	require.NoError(t, err)
	assert.Equal(t, "CMP001", ledger.CompanyID)
	assert.Equal(t, "BA001", ledger.BankAccountID)
	assert.Equal(t, float64(0), ledger.OpeningBalance)
}

func TestGetBalance(t *testing.T) {
	svc, ctx := setupBankService(t)
	balance, err := svc.GetBalance(ctx, "CMP001", "BA001")
	require.NoError(t, err)
	assert.Equal(t, float64(0), balance)
}
