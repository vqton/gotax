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
