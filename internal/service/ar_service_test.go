package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func timePtr(t time.Time) *time.Time { return &t }

func setupSaleSvc(t *testing.T) (*SaleService, context.Context) {
	t.Helper()
	repo := repository.NewMemorySaleRepo()
	svc := NewSaleService(repo, repo, repo, repo, repo, repo, repo, repo, nil)
	return svc, context.Background()
}

func seedInvoice(t *testing.T, svc *SaleService, ctx context.Context, inv *domain.CustomerInvoice) {
	t.Helper()
	err := svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)
	// transition to posted so it appears in AR
	err = svc.invRepo.PostInvoice(ctx, inv.ID, time.Now().UTC())
	require.NoError(t, err)
}

func TestGetARAgingReport_BucketsByDueDate(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "C001", Name: "TestCo",
		TaxCode: "1234567890", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(24 * time.Hour)

	tests := []struct {
		name     string
		dueDate  time.Time
		amount   float64
		bucket30 float64
		bucket60 float64
		bucket90 float64
	}{
		{"current", now.Add(10 * 24 * time.Hour), 100, 0, 0, 0},
		{"overdue_30", now.Add(-15 * 24 * time.Hour), 200, 200, 0, 0},
		{"overdue_60", now.Add(-45 * 24 * time.Hour), 300, 0, 300, 0},
		{"overdue_90", now.Add(-75 * 24 * time.Hour), 400, 0, 0, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := &domain.CustomerInvoice{
				CompanyID:       "c1",
				InvoiceNumber:   "INV-" + tt.name,
				InvoiceDate:     now,
				CustomerID:      cust.ID,
				CustomerName:    "TestCo",
				CustomerTaxCode: "1234567890",
				CustomerAddress: "addr",
				InvoiceType:     "domestic",
				DueDate:         &tt.dueDate,
				Currency:        "VND",
				Status:          domain.SInvDraft,
				Lines: []domain.InvLine{{
					ItemName: "item", Quantity: 1, UnitPrice: tt.amount,
					VATRate: 10, VATType: "VAT_10",
					RevenueAccount: "5111", VATAccountID: "3331",
				}},
			}
			seedInvoice(t, svc, ctx, inv)
		})
	}

	reports, err := svc.GetARAgingReport(ctx, "c1")
	require.NoError(t, err)
	require.Len(t, reports, 1)

	r := reports[0]
	assert.Equal(t, cust.ID, r.CustomerID)
	assert.Equal(t, 110.0, r.Buckets.Bucket0, "current invoice should be in Bucket0")
	assert.Equal(t, 220.0, r.Buckets.Bucket30, "15d overdue should be in Bucket30")
	assert.Equal(t, 330.0, r.Buckets.Bucket60, "45d overdue should be in Bucket60")
	assert.Equal(t, 440.0, r.Buckets.Bucket90, "75d overdue should be in Bucket90")
	assert.Equal(t, 1100.0, r.Buckets.Total, "total should be sum of all invoices (incl VAT)")
}

func TestGetARAgingReport_NoInvoices(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	reports, err := svc.GetARAgingReport(ctx, "c1")
	require.NoError(t, err)
	assert.Empty(t, reports)
}

func TestGetARAgingReport_PaidInvoiceExcluded(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "C002", Name: "PaidCo",
		TaxCode: "1234567891", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(24 * time.Hour)
	future := now.Add(30 * 24 * time.Hour)
	dueDate := now.Add(-10 * 24 * time.Hour)

	// unpaid invoice
	inv1 := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-unpaid",
		InvoiceDate: now, CustomerID: cust.ID,
		CustomerName: "PaidCo", CustomerTaxCode: "1234567891",
		CustomerAddress: "addr", InvoiceType: "domestic",
		DueDate: &dueDate, Currency: "VND",
		Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "item", Quantity: 1, UnitPrice: 500,
			VATRate: 10, VATType: "VAT_10", RevenueAccount: "5111", VATAccountID: "3331"}},
	}
	seedInvoice(t, svc, ctx, inv1)

	// fully paid invoice — due in future, fully received
	inv2 := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-paid",
		InvoiceDate: now, CustomerID: cust.ID,
		CustomerName: "PaidCo", CustomerTaxCode: "1234567891",
		CustomerAddress: "addr", InvoiceType: "domestic",
		DueDate: &future, Currency: "VND", AmountReceived: 550,
		Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "item", Quantity: 1, UnitPrice: 500,
			VATRate: 10, VATType: "VAT_10", RevenueAccount: "5111", VATAccountID: "3331"}},
	}
	seedInvoice(t, svc, ctx, inv2)

	reports, err := svc.GetARAgingReport(ctx, "c1")
	require.NoError(t, err)
	require.Len(t, reports, 1)

	// only unpaid invoice appears; BalanceDue = 550
	assert.Equal(t, 550.0, reports[0].Buckets.Total)
}

func TestGetCustomerStatement(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "STMT01", Name: "StmtCo",
		TaxCode: "999", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	now := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)

	inv1 := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-STMT-1",
		InvoiceDate: now.Add(-10 * 24 * time.Hour),
		DueDate:     timePtr(now.Add(20 * 24 * time.Hour)),
		CustomerID: cust.ID, CustomerName: "StmtCo",
		CustomerTaxCode: "999", CustomerAddress: "addr",
		InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "svc", Quantity: 1, UnitPrice: 200,
			VATRate: 10, VATType: "VAT_10",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "test-user",
	}
	svc.CreateInvoice(ctx, inv1)
	svc.invRepo.PostInvoice(ctx, inv1.ID, now.Add(-10*24*time.Hour))

	inv2 := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-STMT-2",
		InvoiceDate: now.Add(-5 * 24 * time.Hour),
		DueDate:     timePtr(now.Add(25 * 24 * time.Hour)),
		CustomerID: cust.ID, CustomerName: "StmtCo",
		CustomerTaxCode: "999", CustomerAddress: "addr",
		InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "svc2", Quantity: 1, UnitPrice: 300,
			VATRate: 10, VATType: "VAT_10",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "test-user",
	}
	svc.CreateInvoice(ctx, inv2)
	svc.invRepo.PostInvoice(ctx, inv2.ID, now.Add(-5*24*time.Hour))

	rcpt := &domain.CustomerReceipt{
		CompanyID: "c1", ReceiptNumber: "RCP-STMT",
		CustomerID: cust.ID, ReceiptDate: now.Add(-1 * 24 * time.Hour),
		PaymentMethod: "bank_transfer", Currency: "VND",
		Amount: 150, Status: domain.RcpDraft,
	}
	svc.CreateReceipt(ctx, rcpt)
	svc.rcptRepo.UpdateReceiptStatus(ctx, rcpt.ID, domain.RcpPosted)

	fromDate := now.Add(-30 * 24 * time.Hour)
	toDate := now.Add(30 * 24 * time.Hour)

	stmt, err := svc.GetCustomerStatement(ctx, cust.ID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	require.NoError(t, err)
	require.NotNil(t, stmt)
	assert.Equal(t, cust.ID, stmt.Customer.ID)
	assert.Equal(t, 0.0, stmt.OpeningBal)

	require.GreaterOrEqual(t, len(stmt.Lines), 3)

	// last line's balance = closing balance
	lastLine := stmt.Lines[len(stmt.Lines)-1]
	assert.Equal(t, 400.0, lastLine.Balance) // 220 + 330 - 150 = 400
	assert.InDelta(t, 400, stmt.ClosingBal, 0.001)
}

func TestCreditLimit_BlocksInvoiceOverLimit(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID:   "c1",
		Code:        "CRED01",
		Name:        "CreditCo",
		TaxCode:     "999",
		CreditLimit: 500,
		Currency:    "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	// invoice within limit (300 < 500) — should pass
	inv1 := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-CR-1",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "CreditCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "item", Quantity: 1, UnitPrice: 300,
			RevenueAccount: "5111", VATRate: 0, VATType: "NON_TAX",
			VATAccountID: "3331"}},
		CreatedBy: "test-user",
	}
	err = svc.CreateInvoice(ctx, inv1)
	require.NoError(t, err)

	// second invoice would push total to 700 > 500 — should block
	inv2 := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-CR-2",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "CreditCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "item", Quantity: 1, UnitPrice: 400,
			RevenueAccount: "5111", VATRate: 0, VATType: "NON_TAX",
			VATAccountID: "3331"}},
		CreatedBy: "test-user",
	}
	err = svc.CreateInvoice(ctx, inv2)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrCreditLimitExceeded)
}

func TestCreditLimit_AllowsInvoiceAtLimit(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "CRED02", Name: "EdgeCo",
		TaxCode: "999", CreditLimit: 250, Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	// existing invoice: 150
	inv1 := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-EDGE-1",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "EdgeCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "a", Quantity: 1, UnitPrice: 150,
			RevenueAccount: "5111", VATRate: 0, VATType: "NON_TAX",
			VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = svc.CreateInvoice(ctx, inv1)
	require.NoError(t, err)

	// second invoice: 100, total = 250, exactly at limit — should pass
	inv2 := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-EDGE-2",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "EdgeCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "b", Quantity: 1, UnitPrice: 100,
			RevenueAccount: "5111", VATRate: 0, VATType: "NON_TAX",
			VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = svc.CreateInvoice(ctx, inv2)
	require.NoError(t, err)
}

func TestCreditLimit_BlocksSOApproveOverLimit(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "CRED03", Name: "SOCo",
		TaxCode: "999", CreditLimit: 200, Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	// existing invoice: 150
	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-SO-OVR",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "SOCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "a", Quantity: 1, UnitPrice: 150,
			RevenueAccount: "5111", VATRate: 0, VATType: "NON_TAX",
			VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)

	// SO for 100 would make total 250 > 200 — block
	so := &domain.SalesOrder{
		CompanyID: "c1", SONumber: "SO-CR-1",
		OrderDate: time.Now(), CustomerID: cust.ID,
		Currency: "VND", Status: domain.SODraft,
		Lines: []domain.SOLine{{ItemName: "svc", Quantity: 1, UnitPrice: 100}},
		CreatedBy: "user",
	}
	err = svc.CreateSO(ctx, so)
	require.NoError(t, err)

	err = svc.ApproveSO(ctx, so.ID, "approver")
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrCreditLimitExceeded)
}

func TestCreditLimit_NoLimitMeansNoCheck(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "CRED04", Name: "NoLimitCo",
		TaxCode: "999", CreditLimit: 0, Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-NOLIMIT",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "NoLimitCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "big", Quantity: 1, UnitPrice: 999999,
			RevenueAccount: "5111", VATRate: 0, VATType: "NON_TAX",
			VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)
}

func TestPrepayment_CreatesAndOffsets(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "PRE01", Name: "PrepayCo",
		TaxCode: "999", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-PRE",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "PrepayCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "svc", Quantity: 1, UnitPrice: 300,
			VATRate: 10, VATType: "VAT_10",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)

	// create prepayment receipt for 100
	rcpt := &domain.CustomerReceipt{
		CompanyID: "c1", ReceiptNumber: "PRE-RCP-001",
		CustomerID: cust.ID, ReceiptDate: time.Now(),
		PaymentMethod: "bank_transfer", Currency: "VND",
		Amount: 100, Status: domain.RcpDraft,
	}
	err = svc.CreateReceipt(ctx, rcpt)
	require.NoError(t, err)

	// post receipt
	err = svc.PostReceipt(ctx, rcpt.ID)
	require.NoError(t, err)

	// allocate 100 to invoice
	err = svc.AllocateReceipt(ctx, rcpt.ID, inv.ID, 100)
	require.NoError(t, err)

	// verify invoice BalanceDue decreased
	updatedInv, err := svc.GetInvoice(ctx, inv.ID)
	require.NoError(t, err)
	assert.Equal(t, 330.0-100.0, updatedInv.BalanceDue)

	// verify receipt UnallocatedAmount
	updatedRcpt, err := svc.GetReceipt(ctx, rcpt.ID)
	require.NoError(t, err)
	assert.Equal(t, 0.0, updatedRcpt.UnallocatedAmount)
}

func TestPrepayment_AllocationCannotExceedInvoice(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "PRE02", Name: "OverCo",
		TaxCode: "999", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-PRE-2",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "OverCo", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "svc", Quantity: 1, UnitPrice: 200,
			VATRate: 10, VATType: "VAT_10",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)

	rcpt := &domain.CustomerReceipt{
		CompanyID: "c1", ReceiptNumber: "PRE-RCP-002",
		CustomerID: cust.ID, ReceiptDate: time.Now(),
		PaymentMethod: "bank_transfer", Currency: "VND",
		Amount: 500, Status: domain.RcpDraft,
	}
	err = svc.CreateReceipt(ctx, rcpt)
	require.NoError(t, err)
	err = svc.PostReceipt(ctx, rcpt.ID)
	require.NoError(t, err)

	// try to allocate 300 to invoice that only has 220 balance
	err = svc.AllocateReceipt(ctx, rcpt.ID, inv.ID, 300)
	require.Error(t, err)
}

func TestPrepayment_AllocCannotExceedUnallocated(t *testing.T) {
	svc, ctx := setupSaleSvc(t)

	cust := &domain.Customer{
		CompanyID: "c1", Code: "PRE03", Name: "OverRcp",
		TaxCode: "999", Currency: "VND",
	}
	err := svc.CreateCustomer(ctx, cust)
	require.NoError(t, err)

	inv := &domain.CustomerInvoice{
		CompanyID: "c1", InvoiceNumber: "INV-PRE-3",
		InvoiceDate: time.Now(), CustomerID: cust.ID,
		CustomerName: "OverRcp", CustomerTaxCode: "999",
		CustomerAddress: "addr", InvoiceType: "domestic",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "svc", Quantity: 1, UnitPrice: 500,
			VATRate: 0, VATType: "NON_TAX",
			RevenueAccount: "5111", VATAccountID: "3331"}},
		CreatedBy: "user",
	}
	err = svc.CreateInvoice(ctx, inv)
	require.NoError(t, err)

	rcpt := &domain.CustomerReceipt{
		CompanyID: "c1", ReceiptNumber: "PRE-RCP-003",
		CustomerID: cust.ID, ReceiptDate: time.Now(),
		PaymentMethod: "bank_transfer", Currency: "VND",
		Amount: 50, Status: domain.RcpDraft,
	}
	err = svc.CreateReceipt(ctx, rcpt)
	require.NoError(t, err)
	err = svc.PostReceipt(ctx, rcpt.ID)
	require.NoError(t, err)

	// try to allocate 100 from receipt that only has 50
	err = svc.AllocateReceipt(ctx, rcpt.ID, inv.ID, 100)
	require.Error(t, err)
}

// ─── S2: AR txn auto-population ────────────────────────────────────────

func seedCust(t *testing.T, svc *SaleService, ctx context.Context, id, companyID string) {
	t.Helper()
	cust := &domain.Customer{
		ID: id, CompanyID: companyID, Code: "C-" + id, Name: "TestCo",
		TaxCode: "1234567890", Currency: "VND",
	}
	require.NoError(t, svc.CreateCustomer(ctx, cust))
}

func TestPostInvoice_CreatesARTransaction(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")

	inv := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-2026-0001",
		CustomerID: "c1", InvoiceDate: time.Now().UTC(),
		CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{
			RevenueAccount: "5111", VATAccountID: "33311",
			Quantity: 10, UnitPrice: 200, VATRate: 0,
		}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))
	require.NoError(t, svc.PostInvoice(ctx, inv.ID))

	arTxns, err := svc.artRepo.ListARTransactions(ctx, "co1", "c1")
	require.NoError(t, err)
	require.NotEmpty(t, arTxns)

	found := false
	for _, txn := range arTxns {
		if txn.TransactionType == domain.ARTransInvoice && txn.InvoiceID == inv.ID {
			assert.Equal(t, 2000.0, txn.Amount)
			assert.Equal(t, "VND", txn.Currency)
			assert.Equal(t, "customer_invoices", txn.ReferenceType)
			found = true
			break
		}
	}
	assert.True(t, found, "expected AR transaction for invoice")
}

func TestPostReceipt_CreatesARTransaction(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")

	rcpt := &domain.CustomerReceipt{
		CompanyID: "co1", ReceiptNumber: "RC-2026-0001",
		CustomerID: "c1", ReceiptDate: time.Now().UTC(),
		PaymentMethod: "bank_transfer", Currency: "VND",
		Amount: 500, Status: domain.RcpDraft,
	}
	require.NoError(t, svc.CreateReceipt(ctx, rcpt))
	require.NoError(t, svc.PostReceipt(ctx, rcpt.ID))

	arTxns, err := svc.artRepo.ListARTransactions(ctx, "co1", "c1")
	require.NoError(t, err)
	require.NotEmpty(t, arTxns)

	found := false
	for _, txn := range arTxns {
		if txn.TransactionType == domain.ARTransReceipt && txn.ReferenceID == rcpt.ID {
			assert.Equal(t, 500.0, txn.Amount)
			assert.Equal(t, "VND", txn.Currency)
			found = true
			break
		}
	}
	assert.True(t, found, "expected AR transaction for receipt")
}

func TestPostCN_CreatesARTransaction(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")

	inv := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-2026-0010",
		CustomerID: "c1", InvoiceDate: time.Now().UTC(),
		CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		Currency: "VND", Status: domain.SInvDraft,
		Lines: []domain.InvLine{{
			RevenueAccount: "5111", VATAccountID: "33311",
			Quantity: 10, UnitPrice: 500, VATRate: 0,
		}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))
	require.NoError(t, svc.invRepo.PostInvoice(ctx, inv.ID, time.Now().UTC()))

	cn := &domain.CreditNote{
		CompanyID: "co1", CNNumber: "CN-2026-0001",
		OriginalInvoiceID: inv.ID, CustomerID: "c1",
		ReturnDate: time.Now().UTC(), ReturnType: domain.RetPartial,
		Status: domain.CNDraft,
		Lines: []domain.CNLine{{
			ItemName: "Widget", Unit: "pcs",
			Quantity: 2, UnitPrice: 500, VATRate: 0,
		}},
	}
	require.NoError(t, svc.CreateCN(ctx, cn))
	require.NoError(t, svc.PostCN(ctx, cn.ID))

	arTxns, err := svc.artRepo.ListARTransactions(ctx, "co1", "c1")
	require.NoError(t, err)
	require.NotEmpty(t, arTxns)

	found := false
	for _, txn := range arTxns {
		if txn.TransactionType == domain.ARTransCreditNote && txn.ReferenceID == cn.ID {
			assert.Equal(t, 1000.0, txn.Amount)
			assert.Equal(t, inv.ID, txn.InvoiceID)
			found = true
			break
		}
	}
	assert.True(t, found, "expected AR transaction for credit note")
}
