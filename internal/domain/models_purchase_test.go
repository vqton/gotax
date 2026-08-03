package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupplierValidate_Valid(t *testing.T) {
	s := &Supplier{Code: "S1", Name: "ABC", TaxCode: "TX1"}
	require.NoError(t, s.Validate())
	assert.Equal(t, "VND", s.Currency)
	assert.Equal(t, SupplierActive, s.Status)
}

func TestSupplierValidate_Invalid(t *testing.T) {
	tests := []struct {
		name string
		s    Supplier
		want error
	}{
		{"code", Supplier{Name: "n", TaxCode: "t"}, ErrSupplierCodeRequired},
		{"name", Supplier{Code: "c", TaxCode: "t"}, ErrSupplierNameRequired},
		{"taxcode", Supplier{Code: "c", Name: "n"}, ErrSupplierTaxCodeRequired},
		{"status", Supplier{Code: "c", Name: "n", TaxCode: "t", Status: "weird"}, ErrSupplierStatusInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.s.Validate(), tt.want)
		})
	}
}

func TestPOStatus_ValidTransition(t *testing.T) {
	assert.True(t, POStatusDraft.ValidTransition(POStatusApproved))
	assert.True(t, POStatusApproved.ValidTransition(POStatusSent))
	assert.True(t, POStatusSent.ValidTransition(POStatusClosed))
	assert.False(t, POStatusClosed.ValidTransition(POStatusDraft))
	assert.False(t, POStatusDraft.ValidTransition(POStatusReceived))
	assert.False(t, POStatusReceived.ValidTransition(POStatusCancelled))
}

func TestPurchaseOrderValidate_Valid(t *testing.T) {
	po := &PurchaseOrder{
		PONumber: "PO-1", SupplierID: "sup", OrderDate: time.Now(),
		Lines: []POItem{{ItemName: "W", Unit: "pcs", Quantity: 1, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, po.Validate())
	assert.Equal(t, "VND", po.Currency)
	assert.Equal(t, POStatusDraft, po.Status)
}

func TestPurchaseOrderValidate_Invalid(t *testing.T) {
	tests := []struct {
		name string
		po   PurchaseOrder
		want error
	}{
		{"number", PurchaseOrder{SupplierID: "s", OrderDate: time.Now(), Lines: []POItem{{}}}, ErrPONumberRequired},
		{"supplier", PurchaseOrder{PONumber: "p", OrderDate: time.Now(), Lines: []POItem{{}}}, ErrPOSupplierRequired},
		{"date", PurchaseOrder{PONumber: "p", SupplierID: "s", Lines: []POItem{{}}}, ErrPODateRequired},
		{"lines", PurchaseOrder{PONumber: "p", SupplierID: "s", OrderDate: time.Now()}, ErrPOLinesRequired},
		{"status", PurchaseOrder{PONumber: "p", SupplierID: "s", OrderDate: time.Now(), Status: "bad", Lines: []POItem{{}}}, ErrPOStatusInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.po.Validate(), tt.want)
		})
	}
}

func TestPurchaseOrder_CalculateTotals(t *testing.T) {
	po := &PurchaseOrder{Lines: []POItem{
		{Quantity: 10, UnitPrice: 100000, VATRate: 10, DiscountPct: 10},
		{Quantity: 2, UnitPrice: 50000, VATRate: 8},
	}}
	po.CalculateTotals()
	assert.InDelta(t, 1000000, po.Subtotal, 0.001)
	assert.InDelta(t, 100000, po.DiscountAmount, 0.001)
	assert.InDelta(t, 98000, po.TaxAmount, 0.001)
	assert.InDelta(t, 1098000, po.TotalAmount, 0.001)
}

func TestPOItemValidate(t *testing.T) {
	l := POItem{ItemName: "W", Unit: "pcs", Quantity: 1, AccountID: "152", VATAccountID: "1331"}
	require.NoError(t, l.Validate())
	assert.ErrorIs(t, (&POItem{Unit: "pcs", Quantity: 1, AccountID: "152", VATAccountID: "1331"}).Validate(), ErrPOItemNameRequired)
	assert.ErrorIs(t, (&POItem{ItemName: "W", Quantity: 1, AccountID: "152", VATAccountID: "1331"}).Validate(), ErrPOItemUnitRequired)
	assert.ErrorIs(t, (&POItem{ItemName: "W", Unit: "pcs", Quantity: 0, AccountID: "152", VATAccountID: "1331"}).Validate(), ErrPOItemQuantityRequired)
	assert.ErrorIs(t, (&POItem{ItemName: "W", Unit: "pcs", Quantity: 1, AccountID: "152", VATAccountID: "1331", UnitPrice: -1}).Validate(), ErrPOItemPriceRequired)
	assert.ErrorIs(t, (&POItem{ItemName: "W", Unit: "pcs", Quantity: 1, VATAccountID: "1331"}).Validate(), ErrPOItemAccountRequired)
	assert.ErrorIs(t, (&POItem{ItemName: "W", Unit: "pcs", Quantity: 1, AccountID: "152"}).Validate(), ErrPOItemVATAccountRequired)
}

func TestGRNValidate_Valid(t *testing.T) {
	g := &GRN{GRNNumber: "G1", POID: "po", ReceiptDate: time.Now(), Lines: []GRNItem{{ItemName: "W", POLineID: "l"}}}
	require.NoError(t, g.Validate())
	assert.Equal(t, GRNDraft, g.Status)
}

func TestGRNValidate_Invalid(t *testing.T) {
	tests := []struct {
		name string
		g    GRN
		want error
	}{
		{"number", GRN{POID: "p", ReceiptDate: time.Now(), Lines: []GRNItem{{}}}, ErrGRNNumberRequired},
		{"po", GRN{GRNNumber: "g", ReceiptDate: time.Now(), Lines: []GRNItem{{}}}, ErrGRNPORequired},
		{"date", GRN{GRNNumber: "g", POID: "p", Lines: []GRNItem{{}}}, ErrGRNDateRequired},
		{"lines", GRN{GRNNumber: "g", POID: "p", ReceiptDate: time.Now()}, ErrGRNLinesRequired},
		{"status", GRN{GRNNumber: "g", POID: "p", ReceiptDate: time.Now(), Status: "bad", Lines: []GRNItem{{}}}, ErrGRNStatusInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.g.Validate(), tt.want)
		})
	}
}

func TestGRNItemValidate(t *testing.T) {
	l := GRNItem{ItemName: "W", POLineID: "l", QuantityReceived: 1}
	require.NoError(t, l.Validate())
	assert.ErrorIs(t, (&GRNItem{POLineID: "l"}).Validate(), ErrGRNItemNameRequired)
	assert.ErrorIs(t, (&GRNItem{ItemName: "W"}).Validate(), ErrGRNItemPOLineRequired)
	assert.ErrorIs(t, (&GRNItem{ItemName: "W", POLineID: "l", QuantityReceived: -1}).Validate(), ErrGRNItemQtyRequired)
}

func TestSupplierInvoiceValidate_Valid(t *testing.T) {
	inv := &SupplierInvoice{
		InvoiceNumber: "I1", SupplierID: "s", SupplierName: "ABC", SupplierTaxCode: "TX",
		InvoiceDate: time.Now(), Lines: []SupplierInvoiceLine{{ItemName: "W", Quantity: 1, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, inv.Validate())
	assert.Equal(t, "VND", inv.Currency)
	assert.Equal(t, InvoiceDraft, inv.Status)
}

func TestSupplierInvoiceValidate_Invalid(t *testing.T) {
	tests := []struct {
		name string
		inv  SupplierInvoice
		want error
	}{
		{"number", SupplierInvoice{SupplierID: "s", SupplierName: "n", SupplierTaxCode: "t", InvoiceDate: time.Now(), Lines: []SupplierInvoiceLine{{}}}, ErrInvoiceNumberRequired},
		{"supplier", SupplierInvoice{InvoiceNumber: "i", SupplierName: "n", SupplierTaxCode: "t", InvoiceDate: time.Now(), Lines: []SupplierInvoiceLine{{}}}, ErrInvoiceSupplierRequired},
		{"name", SupplierInvoice{InvoiceNumber: "i", SupplierID: "s", SupplierTaxCode: "t", InvoiceDate: time.Now(), Lines: []SupplierInvoiceLine{{}}}, ErrInvoiceSupplierNameRequired},
		{"taxcode", SupplierInvoice{InvoiceNumber: "i", SupplierID: "s", SupplierName: "n", InvoiceDate: time.Now(), Lines: []SupplierInvoiceLine{{}}}, ErrInvoiceSupplierTaxCodeRequired},
		{"date", SupplierInvoice{InvoiceNumber: "i", SupplierID: "s", SupplierName: "n", SupplierTaxCode: "t", Lines: []SupplierInvoiceLine{{}}}, ErrInvoiceDateRequired},
		{"lines", SupplierInvoice{InvoiceNumber: "i", SupplierID: "s", SupplierName: "n", SupplierTaxCode: "t", InvoiceDate: time.Now()}, ErrInvoiceLinesRequired},
		{"status", SupplierInvoice{InvoiceNumber: "i", SupplierID: "s", SupplierName: "n", SupplierTaxCode: "t", InvoiceDate: time.Now(), Status: "bad", Lines: []SupplierInvoiceLine{{}}}, ErrInvoiceStatusInvalidPurchase},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.inv.Validate(), tt.want)
		})
	}
}

func TestSupplierInvoice_CalculateTotals(t *testing.T) {
	inv := &SupplierInvoice{Lines: []SupplierInvoiceLine{
		{Quantity: 2, UnitPrice: 100000, VATRate: 10},
		{Quantity: 1, UnitPrice: 50000, VATRate: 8},
	}}
	inv.CalculateTotals()
	assert.InDelta(t, 250000, inv.Subtotal, 0.001)
	assert.InDelta(t, 24000, inv.TaxAmount, 0.001)
	assert.InDelta(t, 274000, inv.TotalAmount, 0.001)
}

func TestSupplierInvoiceLineValidate(t *testing.T) {
	l := SupplierInvoiceLine{ItemName: "W", Quantity: 1, AccountID: "152", VATAccountID: "1331"}
	require.NoError(t, l.Validate())
	assert.ErrorIs(t, (&SupplierInvoiceLine{Quantity: 1, AccountID: "152", VATAccountID: "1331"}).Validate(), ErrInvoiceItemNameRequired)
	assert.ErrorIs(t, (&SupplierInvoiceLine{ItemName: "W", Quantity: 0, AccountID: "152", VATAccountID: "1331"}).Validate(), ErrInvoiceItemQtyRequired)
	assert.ErrorIs(t, (&SupplierInvoiceLine{ItemName: "W", Quantity: 1, AccountID: "152", VATAccountID: "1331", UnitPrice: -1}).Validate(), ErrInvoiceItemPriceRequired)
	assert.ErrorIs(t, (&SupplierInvoiceLine{ItemName: "W", Quantity: 1, VATAccountID: "1331"}).Validate(), ErrInvoiceItemAccountRequired)
	assert.ErrorIs(t, (&SupplierInvoiceLine{ItemName: "W", Quantity: 1, AccountID: "152"}).Validate(), ErrInvoiceItemVATAccountRequired)
}

func TestDoubtfulDebtRate_Tiers(t *testing.T) {
	tests := []struct {
		months int
		want   float64
	}{
		{0, 0}, {5, 0}, {6, 0.30}, {11, 0.30}, {12, 0.50}, {23, 0.50},
		{24, 0.70}, {35, 0.70}, {36, 1.00}, {60, 1.00},
	}
	for _, tt := range tests {
		rate := DoubtfulDebtRate(tt.months)
		assert.Equal(t, tt.want, rate, "months=%d", tt.months)
	}
}

func TestDoubtfulDebtProvision_Validate(t *testing.T) {
	valid := &DoubtfulDebtProvision{
		CompanyID: "c1", AsOfDate: "2026-07-31", Status: ProvisionDraft,
		Lines: []DoubtfulDebtProvisionLine{{
			SupplierID: "s", SupplierName: "ABC", AgeMonths: 12,
			RatePct: 50, OutstandingAmount: 1000, ProvisionAmount: 500,
		}},
	}
	require.NoError(t, valid.Validate())
	invalid := *valid
	invalid.AsOfDate = ""
	assert.ErrorIs(t, invalid.Validate(), ErrProvisionDateRequired)
	invalid = *valid
	invalid.Lines = nil
	assert.ErrorIs(t, invalid.Validate(), ErrProvisionNoLines)
}
