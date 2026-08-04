package repository

import (
	"reflect"
	"testing"
	"time"

	"gotax/internal/domain"
)

// Mapping round-trip tests: domain -> GORM -> domain must preserve all
// fields (PG parity contract, models_gorm_sale*.go == models_sale.go).
// Regression guard for lossy mapping (balance_due, cost_price, status,
// credit_limit, e_invoice fields previously dropped on PG).

func fixedTime() time.Time {
	return time.Date(2026, 8, 4, 10, 30, 0, 0, time.UTC)
}

func TestCustomerMappingRoundTrip(t *testing.T) {
	in := &domain.Customer{
		ID: "c1", CompanyID: "co1", Code: "KH001", Name: "Cty ABC",
		TaxCode: "0123456789", Address: "Hanoi", Phone: "0912345678",
		Email: "a@b.vn", BankAccountName: "ABC", BankAccountNumber: "123456",
		BankName: "VCB", PaymentTerms: "net30", CreditLimit: 1000000,
		Currency: "VND", CustomerType: domain.CustomerExport,
		CustomerGroup: domain.CustGroupWholesale, PriceListID: "pl1",
		Status: domain.CustomerBlacklisted, Notes: "n",
		CreatedBy: "u1", CreatedAt: fixedTime(), UpdatedAt: fixedTime(),
	}
	got := customerFromGORM(customerToGORM(in))
	eq := *in == *got
	if !eq {
		t.Fatalf("customer round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}

func TestSalesOrderMappingRoundTrip(t *testing.T) {
	exp := fixedTime()
	in := &domain.SalesOrder{
		ID: "so1", CompanyID: "co1", SONumber: "SO-2026-00001",
		QuotationID: "sq1", CustomerID: "c1", OrderDate: fixedTime(),
		ExpectedDate: &exp, Currency: "USD", ExchangeRate: 25400,
		PaymentTerms: "net30", DeliveryTerms: "FOB",
		ShippingAddress: "HCM", Subtotal: 1000, DiscountAmount: 50,
		TaxAmount: 100, TotalAmount: 1050, Status: domain.SOApproved,
		ApprovedBy: "u2", ApprovedAt: &exp, CancelledReason: "",
		Notes: "n", CreatedBy: "u1", CreatedAt: fixedTime(), UpdatedAt: fixedTime(),
		Lines: []domain.SOLine{{
			ID: "l1", SOID: "so1", LineNumber: 1, ItemCode: "ITM1",
			ItemName: "Widget", Unit: "pcs", Quantity: 10, UnitPrice: 100,
			DiscountPct: 5, VATRate: 10, VATType: "VAT10",
			RevenueAccount: "5111", VATAccountID: "33311",
			LineTotal: 950, LineVATAmount: 95, DeliveredQty: 5, InvoicedQty: 0,
		}},
	}
	got := soFromGORM(soToGORM(in))
	if len(got.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got.Lines))
	}
	if !reflect.DeepEqual(*in, *got) {
		t.Fatalf("SO round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}

func TestDeliveryNoteMappingRoundTrip(t *testing.T) {
	in := &domain.DeliveryNote{
		ID: "dn1", CompanyID: "co1", DNNumber: "DN-2026-00001",
		SOID: "so1", DeliveryDate: fixedTime(), Warehouse: "WH1",
		ShippingMethod: "road", CarrierName: "GHTK", TrackingNumber: "T1",
		DeliveryAddress: "HN", Status: domain.DNPosted, Notes: "n",
		CreatedBy: "u1", CreatedAt: fixedTime(), UpdatedAt: fixedTime(),
		Lines: []domain.DNLine{{
			ID: "dl1", DNID: "dn1", SOLineID: "l1", ItemCode: "ITM1",
			ItemName: "Widget", Unit: "pcs", QtyDelivered: 5,
			QtyReturned: 1, UnitPrice: 100, LineTotal: 500, CostPrice: 60,
		}},
	}
	got := dnFromGORM(dnToGORM(in))
	if len(got.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got.Lines))
	}
	if !reflect.DeepEqual(*in, *got) {
		t.Fatalf("DN round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}

func TestInvoiceMappingRoundTrip(t *testing.T) {
	due := fixedTime()
	glAt := fixedTime().Add(1 * time.Hour)
	in := &domain.CustomerInvoice{
		ID: "inv1", CompanyID: "co1", InvoiceNumber: "INV-2026-00001",
		InvoiceDate: fixedTime(), SOID: "so1", DNID: "dn1",
		CustomerID: "c1", CustomerName: "Cty ABC", CustomerTaxCode: "0123456789",
		CustomerAddress: "HN", InvoiceType: "standard", Currency: "USD",
		ExchangeRate: 25400, Subtotal: 1000, DiscountAmount: 0,
		TaxAmount: 100, TotalAmount: 1100, AmountReceived: 500,
		BalanceDue: 600, DueDate: &due, InvoiceNote: "in",
		EInvoiceData: "<xml/>", EInvoiceCode: "GDT-1",
		EInvStatus: domain.EInvIssued, DigitalSignatureID: "sig1",
		SignedData: "sig", GDTResponse: "ok", OriginalInvoiceID: "",
		AdjustmentType: "", Status: domain.SInvPosted, GLPosted: true,
		GLPostedAt: &glAt, Notes: "n", CreatedBy: "u1",
		CreatedAt: fixedTime(), UpdatedAt: fixedTime(),
		Lines: []domain.InvLine{{
			ID: "il1", InvoiceID: "inv1", SOLineID: "l1", DNLineID: "dl1",
			ItemCode: "ITM1", ItemName: "Widget", Unit: "pcs", Quantity: 5,
			UnitPrice: 100, DiscountPct: 0, VATRate: 10, VATType: "VAT10",
			LineTotal: 500, LineVATAmount: 50,
			RevenueAccount: "5111", VATAccountID: "33311",
		}},
	}
	got := cinvFromGORM(cinvToGORM(in))
	if len(got.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got.Lines))
	}
	if !reflect.DeepEqual(*in, *got) {
		t.Fatalf("invoice round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}

func TestReceiptMappingRoundTrip(t *testing.T) {
	glAt := fixedTime().Add(2 * time.Hour)
	in := &domain.CustomerReceipt{
		ID: "r1", CompanyID: "co1", ReceiptNumber: "RC-2026-00001",
		CustomerID: "c1", ReceiptDate: fixedTime(),
		PaymentMethod: "bank_transfer", BankAccountID: "ba1",
		Currency: "VND", ExchangeRate: 1, Amount: 500,
		UnallocatedAmount: 200, Reference: "ref1", Notes: "n",
		Status: domain.RcpPosted, GLPosted: true, GLPostedAt: &glAt,
		CreatedBy: "u1", CreatedAt: fixedTime(), UpdatedAt: fixedTime(),
		Allocations: []domain.RcpAllocation{{
			ID: "a1", ReceiptID: "r1", InvoiceID: "inv1",
			AllocatedAmount: 300, DiscountAmount: 0,
		}},
	}
	got := rcptFromGORM(rcptToGORM(in))
	if len(got.Allocations) != 1 {
		t.Fatalf("expected 1 allocation, got %d", len(got.Allocations))
	}
	if !reflect.DeepEqual(*in, *got) {
		t.Fatalf("receipt round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}

func TestCreditNoteMappingRoundTrip(t *testing.T) {
	glAt := fixedTime().Add(3 * time.Hour)
	in := &domain.CreditNote{
		ID: "cn1", CompanyID: "co1", CNNumber: "CN-2026-00001",
		OriginalInvoiceID: "inv1", CustomerID: "c1",
		ReturnDate: fixedTime(), ReturnReason: "defect",
		ReturnType: domain.RetPartial, DNID: "dn1",
		Subtotal: 500, TaxAmount: 50, TotalAmount: 550,
		EInvoiceData: "<xml/>", EInvoiceCode: "GDT-CN1",
		Status: domain.CNPosted, GLPosted: true, GLPostedAt: &glAt,
		Notes: "n", CreatedBy: "u1", CreatedAt: fixedTime(), UpdatedAt: fixedTime(),
		Lines: []domain.CNLine{{
			ID: "cl1", CNID: "cn1", InvLineID: "il1",
			ItemName: "Widget", Unit: "pcs", Quantity: 5, UnitPrice: 100,
			VATRate: 10, LineTotal: 500, LineVATAmount: 50,
		}},
	}
	got := cnFromGORM(cnToGORM(in))
	if len(got.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got.Lines))
	}
	if !reflect.DeepEqual(*in, *got) {
		t.Fatalf("CN round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}

func TestARTransactionMappingRoundTrip(t *testing.T) {
	in := &domain.ARTransaction{
		ID: "at1", CompanyID: "co1", CustomerID: "c1",
		InvoiceID: "inv1", TransactionType: domain.ARTransCreditNote,
		TransactionDate: fixedTime(), Amount: 550, Currency: "VND",
		ReferenceType: "credit_notes", ReferenceID: "cn1",
		Notes: "n", CreatedAt: fixedTime(),
	}
	got := artFromGORM(artToGORM(in))
	if !reflect.DeepEqual(*in, *got) {
		t.Fatalf("AR txn round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}

func TestSalesQuotationMappingRoundTrip(t *testing.T) {
	valid := fixedTime().Add(30 * 24 * time.Hour)
	in := &domain.SalesQuotation{
		ID: "sq1", CompanyID: "co1", QNNumber: "SQ-2026-00001",
		CustomerID: "c1", ValidUntil: &valid, Status: "APPROVED",
		TotalAmount: 1000, CreatedBy: "u1", CreatedAt: fixedTime(),
	}
	got := sqFromGORM(sqToGORM(in))
	if !reflect.DeepEqual(*in, *got) {
		t.Fatalf("SQ round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}

func TestSalesQuotationMappingNilValidUntil(t *testing.T) {
	in := &domain.SalesQuotation{
		ID: "sq2", CompanyID: "co1", QNNumber: "SQ-2026-00002",
		CustomerID: "c1", Status: "DRAFT", TotalAmount: 0,
		CreatedBy: "u1", CreatedAt: fixedTime(),
	}
	got := sqFromGORM(sqToGORM(in))
	if !reflect.DeepEqual(*in, *got) {
		t.Fatalf("SQ nil ValidUntil round-trip mismatch\n in=%+v\ngot=%+v", *in, *got)
	}
}
