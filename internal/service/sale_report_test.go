package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
)

// ─── FR-9 Reports ─────────────────────────────────────────────────────

func TestGetSalesLedgerReport_RowsAndTotals(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")
	now := time.Date(2026, 7, 18, 10, 0, 0, 0, time.UTC)

	inv := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-0001", InvoiceDate: now,
		CustomerID: "c1", CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines: []domain.InvLine{{
			ItemName: "Office furniture", Quantity: 10, UnitPrice: 1000,
			VATRate: 10, VATType: "VAT_10", RevenueAccount: "5111", VATAccountID: "3331",
		}},
	}
	seedInvoice(t, svc, ctx, inv) // posted

	cn := &domain.CreditNote{
		CompanyID: "co1", CNNumber: "CN-0001", OriginalInvoiceID: inv.ID,
		CustomerID: "c1", ReturnDate: now.Add(24 * time.Hour), ReturnReason: "defective",
		Status: domain.CNDraft,
		Lines:  []domain.CNLine{{ItemName: "defective chairs", Quantity: 2, UnitPrice: 500}},
	}
	require.NoError(t, svc.CreateCN(ctx, cn))
	require.NoError(t, svc.cnRepo.PostCN(ctx, cn.ID, time.Now().UTC()))

	rpt, err := svc.GetSalesLedgerReport(ctx, "co1", "c1", "2026-07-01", "2026-07-31")
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 2)
	assert.Equal(t, "INV-0001", rpt.Rows[0].Ref)
	assert.Equal(t, 10000.0, rpt.Rows[0].Revenue)
	assert.Equal(t, 1000.0, rpt.Rows[0].VAT)
	assert.Equal(t, 11000.0, rpt.Rows[0].Total)
	assert.Equal(t, "CN-0001", rpt.Rows[1].Ref)
	assert.Equal(t, -1000.0, rpt.Rows[1].Revenue)
	assert.Equal(t, 0.0, rpt.Rows[1].VAT)
	assert.Equal(t, -1000.0, rpt.Rows[1].Total)
	assert.Equal(t, 9000.0, rpt.TotalRevenue)
	assert.Equal(t, 1000.0, rpt.TotalVAT)
	assert.Equal(t, 10000.0, rpt.Total)
}

func TestGetCustomerLedgerReport_RunningBalance(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")
	now := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)

	// prior-period invoice → opening balance
	prior := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-PRIOR", InvoiceDate: now.Add(-30 * 24 * time.Hour),
		CustomerID: "c1", CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines:  []domain.InvLine{{ItemName: "item", Quantity: 1, UnitPrice: 100, VATRate: 0, RevenueAccount: "5111", VATAccountID: "3331"}},
	}
	seedInvoice(t, svc, ctx, prior)

	inv := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-0001", InvoiceDate: now,
		CustomerID: "c1", CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines:  []domain.InvLine{{ItemName: "item", Quantity: 1, UnitPrice: 100, VATRate: 0, RevenueAccount: "5111", VATAccountID: "3331"}},
	}
	seedInvoice(t, svc, ctx, inv) // 100 debit in period

	rcpt := &domain.CustomerReceipt{
		CompanyID: "co1", ReceiptNumber: "RCP-0001", CustomerID: "c1",
		ReceiptDate: now.Add(24 * time.Hour), PaymentMethod: "CASH",
		Amount: 50, UnallocatedAmount: 0, Status: domain.RcpDraft,
	}
	require.NoError(t, svc.CreateReceipt(ctx, rcpt))
	require.NoError(t, svc.PostReceipt(ctx, rcpt.ID))

	rpt, err := svc.GetCustomerLedgerReport(ctx, "co1", "c1", "2026-07-01", "2026-07-31")
	require.NoError(t, err)
	assert.Equal(t, 100.0, rpt.OpeningBalance)
	require.Len(t, rpt.Rows, 2)
	assert.Equal(t, "INV-0001", rpt.Rows[0].Ref)
	assert.Equal(t, 200.0, rpt.Rows[0].Balance) // 100 opening + 100 debit
	assert.Equal(t, "RCP-0001", rpt.Rows[1].Ref)
	assert.Equal(t, 150.0, rpt.Rows[1].Balance)
	assert.Equal(t, 150.0, rpt.ClosingBalance)
}

func TestGetGoodsSalesLedgerReport_PerItem(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")
	now := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)

	inv := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-0001", InvoiceDate: now,
		CustomerID: "c1", CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines: []domain.InvLine{
			{ItemCode: "A", ItemName: "Widget", Unit: "pcs", Quantity: 2, UnitPrice: 100,
				VATRate: 10, RevenueAccount: "5111", VATAccountID: "3331"},
			{ItemCode: "B", ItemName: "Gadget", Unit: "pcs", Quantity: 1, UnitPrice: 50,
				VATRate: 5, RevenueAccount: "5111", VATAccountID: "3331"},
		},
	}
	seedInvoice(t, svc, ctx, inv)

	// draft invoice excluded
	draft := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-DRAFT", InvoiceDate: now,
		CustomerID: "c1", CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines:  []domain.InvLine{{ItemName: "Widget", Unit: "pcs", Quantity: 1, UnitPrice: 100, VATRate: 10, RevenueAccount: "5111", VATAccountID: "3331"}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, draft))

	rpt, err := svc.GetGoodsSalesLedgerReport(ctx, "co1", "", "")
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 2)
	assert.Equal(t, "Gadget", rpt.Rows[0].ItemName) // sorted by item
	assert.Equal(t, "Widget", rpt.Rows[1].ItemName)
	assert.Equal(t, 3.0, rpt.TotalQty)
	assert.Equal(t, 250.0, rpt.TotalRevenue)
	assert.Equal(t, 22.5, rpt.TotalVAT)
}

func TestGetVATOutputReport_PerRate(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")
	now := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)

	for i, l := range []struct {
		rate  float64
		price float64
	}{
		{10, 100}, {10, 50}, {5, 200},
	} {
		inv := &domain.CustomerInvoice{
			CompanyID: "co1", InvoiceNumber: "INV-00" + string(rune('1'+i)), InvoiceDate: now,
			CustomerID: "c1", CustomerName: "TestCo", CustomerTaxCode: "1234567890",
			CustomerAddress: "addr", InvoiceType: "domestic", Currency: "VND",
			Status: domain.SInvDraft,
			Lines: []domain.InvLine{{ItemName: "item", Unit: "pcs", Quantity: 1, UnitPrice: l.price,
				VATRate: l.rate, RevenueAccount: "5111", VATAccountID: "3331"}},
		}
		seedInvoice(t, svc, ctx, inv)
	}

	rpt, err := svc.GetVATOutputReport(ctx, "co1", "", "")
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 2)
	assert.Equal(t, 5.0, rpt.Rows[0].VatRate) // sorted ascending
	assert.Equal(t, 200.0, rpt.Rows[0].Subtotal)
	assert.Equal(t, 10.0, rpt.Rows[0].VatAmount)
	assert.Equal(t, 1, rpt.Rows[0].InvoiceCount)
	assert.Equal(t, 10.0, rpt.Rows[1].VatRate)
	assert.Equal(t, 150.0, rpt.Rows[1].Subtotal)
	assert.Equal(t, 15.0, rpt.Rows[1].VatAmount)
	assert.Equal(t, 2, rpt.Rows[1].InvoiceCount)
	assert.Equal(t, 350.0, rpt.TotalSubtotal)
	assert.Equal(t, 25.0, rpt.TotalVAT)
}

func TestGetUnbilledDeliveriesReport(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")

	// unbilled: posted DN, no invoice
	so1 := seedSO(t, svc, ctx, "SO-UNBILLED-1", 100)
	dn1 := &domain.DeliveryNote{
		CompanyID: "co1", DNNumber: "DN-UNBILLED-1", SOID: so1.ID,
		DeliveryDate: time.Now().UTC(), Status: domain.DNDraft,
		Lines: []domain.DNLine{{
			SOLineID: so1.Lines[0].ID, ItemCode: "A", ItemName: "Widget", Unit: "pcs",
			QtyDelivered: 10, UnitPrice: 100, CostPrice: 60,
		}},
	}
	require.NoError(t, svc.CreateDN(ctx, dn1))
	require.NoError(t, svc.PostDN(ctx, dn1.ID))

	// billed: posted DN covered by posted invoice (DNID link)
	so2 := seedSO(t, svc, ctx, "SO-BILLED-1", 100)
	dn2 := &domain.DeliveryNote{
		CompanyID: "co1", DNNumber: "DN-BILLED-1", SOID: so2.ID,
		DeliveryDate: time.Now().UTC(), Status: domain.DNDraft,
		Lines: []domain.DNLine{{
			SOLineID: so2.Lines[0].ID, ItemCode: "B", ItemName: "Gadget", Unit: "pcs",
			QtyDelivered: 10, UnitPrice: 200, CostPrice: 120,
		}},
	}
	require.NoError(t, svc.CreateDN(ctx, dn2))
	require.NoError(t, svc.PostDN(ctx, dn2.ID))
	inv := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-BILLED-1", InvoiceDate: time.Now().UTC(),
		DNID: dn2.ID, SOID: so2.ID,
		CustomerID: "c1", CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines: []domain.InvLine{{ItemName: "Gadget", Unit: "pcs", Quantity: 10, UnitPrice: 200,
			VATRate: 10, RevenueAccount: "5111", VATAccountID: "3331"}},
	}
	seedInvoice(t, svc, ctx, inv)

	// draft DN excluded
	so3 := seedSO(t, svc, ctx, "SO-DRAFT-1", 100)
	dn3 := &domain.DeliveryNote{
		CompanyID: "co1", DNNumber: "DN-DRAFT-1", SOID: so3.ID,
		DeliveryDate: time.Now().UTC(), Status: domain.DNDraft,
		Lines: []domain.DNLine{{
			SOLineID: so3.Lines[0].ID, ItemName: "Widget", Unit: "pcs",
			QtyDelivered: 5, UnitPrice: 100,
		}},
	}
	require.NoError(t, svc.CreateDN(ctx, dn3))

	rpt, err := svc.GetUnbilledDeliveriesReport(ctx, "co1")
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 1)
	assert.Equal(t, "DN-UNBILLED-1", rpt.Rows[0].DNNumber)
	assert.Equal(t, "c1", rpt.Rows[0].CustomerID)
	assert.Equal(t, "TestCo", rpt.Rows[0].CustomerName)
	assert.Equal(t, 1000.0, rpt.Rows[0].Amount)
}

// ─── Bug fix: InvoiceCount counts invoices, not lines ────────────────

func TestGetVATOutputReport_InvoiceCountCountsInvoices(t *testing.T) {
	svc, ctx := setupSaleSvc(t)
	seedCust(t, svc, ctx, "c1", "co1")
	now := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)

	// one invoice with two 10% lines → counts once, not twice
	inv := &domain.CustomerInvoice{
		CompanyID: "co1", InvoiceNumber: "INV-MULTI", InvoiceDate: now,
		CustomerID: "c1", CustomerName: "TestCo", CustomerTaxCode: "1234567890",
		CustomerAddress: "addr", InvoiceType: "domestic", Currency: "VND",
		Status: domain.SInvDraft,
		Lines: []domain.InvLine{
			{ItemName: "a", Unit: "pcs", Quantity: 1, UnitPrice: 100,
				VATRate: 10, RevenueAccount: "5111", VATAccountID: "3331"},
			{ItemName: "b", Unit: "pcs", Quantity: 1, UnitPrice: 50,
				VATRate: 10, RevenueAccount: "5111", VATAccountID: "3331"},
		},
	}
	seedInvoice(t, svc, ctx, inv)

	rpt, err := svc.GetVATOutputReport(ctx, "co1", "", "")
	require.NoError(t, err)
	require.Len(t, rpt.Rows, 1)
	assert.Equal(t, 10.0, rpt.Rows[0].VatRate)
	assert.Equal(t, 150.0, rpt.Rows[0].Subtotal)
	assert.Equal(t, 15.0, rpt.Rows[0].VatAmount)
	assert.Equal(t, 1, rpt.Rows[0].InvoiceCount, "one invoice with two 10% lines counts once")
}
