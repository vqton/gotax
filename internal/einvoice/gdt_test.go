package einvoice

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
)

const validXML = `<?xml version="1.0" encoding="UTF-8"?>
<Invoice xmlns="http://gdt.gov.vn/schemas/einvoice/2026">
  <InvoiceHeader>
    <InvoiceNumber>INV-12345</InvoiceNumber>
    <InvoiceDate>2026-07-20</InvoiceDate>
    <InvoiceTime>14:30:00</InvoiceTime>
    <InvoiceType>01GTKT</InvoiceType>
    <InvoiceSeries>AA/26E</InvoiceSeries>
    <CurrencyCode>VND</CurrencyCode>
    <ExchangeRate>1</ExchangeRate>
  </InvoiceHeader>
  <Seller>
    <TaxCode>9876543210</TaxCode>
    <Name>Cong ty XYZ</Name>
    <Address>456 Duong DEF, Quan 1, TP.HCM</Address>
  </Seller>
  <Buyer>
    <TaxCode>0123456789</TaxCode>
    <Name>Cong ty ABC</Name>
    <Address>123 Duong ABC, Quan 1, TP.HCM</Address>
  </Buyer>
  <InvoiceLines>
    <Line>
      <LineNumber>1</LineNumber>
      <ItemCode>SP001</ItemCode>
      <ItemName>May tinh Dell 5420</ItemName>
      <Unit>Cai</Unit>
      <Quantity>10</Quantity>
      <UnitPrice>15000000</UnitPrice>
      <VatRate>8</VatRate>
      <VatAmount>12000000</VatAmount>
      <LineTotal>150000000</LineTotal>
    </Line>
  </InvoiceLines>
  <Summary>
    <Subtotal>150000000</Subtotal>
    <Discount>0</Discount>
    <VatAmount>12000000</VatAmount>
    <GrandTotal>162000000</GrandTotal>
  </Summary>
</Invoice>
`

func TestParseValidInvoice(t *testing.T) {
	inv, err := Parse([]byte(validXML))
	require.NoError(t, err)
	assert.Equal(t, "INV-12345", inv.InvoiceNumber)
	assert.Equal(t, "Cong ty XYZ", inv.SupplierName)
	assert.Equal(t, "9876543210", inv.SupplierTaxCode)
	assert.Equal(t, "VND", inv.Currency)
	require.Len(t, inv.Lines, 1)
	line := inv.Lines[0]
	assert.Equal(t, "May tinh Dell 5420", line.ItemName)
	assert.Equal(t, "Cai", line.Unit)
	assert.Equal(t, 10.0, line.Quantity)
	assert.Equal(t, 15000000.0, line.UnitPrice)
	assert.Equal(t, 8.0, line.VATRate)
	assert.Equal(t, domain.VAT8, line.VATType)
	assert.Equal(t, 150000000.0, line.LineTotal)
	assert.Equal(t, 12000000.0, line.LineVATAmount)
	assert.Equal(t, "152", line.AccountID)
	assert.Equal(t, "1331", line.VATAccountID)
	assert.Equal(t, 150000000.0, inv.Subtotal)
	assert.Equal(t, 12000000.0, inv.TaxAmount)
	assert.Equal(t, 162000000.0, inv.TotalAmount)
	assert.Equal(t, "2026-07-20", inv.InvoiceDate.Format("2006-01-02"))
}

func TestParseFCInvoice(t *testing.T) {
	fc := strings.Replace(validXML, "<CurrencyCode>VND</CurrencyCode>", "<CurrencyCode>USD</CurrencyCode>", 1)
	fc = strings.Replace(fc, "<ExchangeRate>1</ExchangeRate>", "<ExchangeRate>25000</ExchangeRate>", 1)
	inv, err := Parse([]byte(fc))
	require.NoError(t, err)
	assert.Equal(t, "USD", inv.Currency)
	assert.Equal(t, 25000.0, inv.ExchangeRate)
}

func TestParseCreditNote(t *testing.T) {
	cn := strings.Replace(validXML, "<InvoiceType>01GTKT</InvoiceType>", "<InvoiceType>credit_note</InvoiceType>", 1)
	inv, err := Parse([]byte(cn))
	require.NoError(t, err)
	assert.Equal(t, domain.InvoiceTypeCreditNote, inv.InvoiceType)
}

func TestParseVATTypes(t *testing.T) {
	for rate, want := range map[string]domain.VATType{"0": domain.VAT0, "5": domain.VAT5, "8": domain.VAT8, "10": domain.VAT10, "-1": domain.VATNonTaxable} {
		x := strings.Replace(validXML, "<VatRate>8</VatRate>", "<VatRate>"+rate+"</VatRate>", 1)
		inv, err := Parse([]byte(x))
		require.NoError(t, err, rate)
		assert.Equal(t, want, inv.Lines[0].VATType, rate)
	}
}

func TestParseMalformed(t *testing.T) {
	_, err := Parse([]byte("<Invoice><InvoiceHeader>"))
	assert.Error(t, err)
}

func TestParseMissingRequired(t *testing.T) {
	_, err := Parse([]byte(`<Invoice xmlns="http://gdt.gov.vn/schemas/einvoice/2026">
  <InvoiceHeader><InvoiceNumber>INV-X</InvoiceNumber><InvoiceDate>2026-07-20</InvoiceDate></InvoiceHeader>
  <Seller><Name>Sup</Name></Seller>
  <InvoiceLines><Line><ItemName>W</ItemName></Line></InvoiceLines>
</Invoice>`))
	assert.ErrorIs(t, err, domain.ErrInvoiceSupplierTaxCodeRequired)
}

func TestGenerateParseRoundTrip(t *testing.T) {
	inv := &domain.SupplierInvoice{
		InvoiceNumber:   "INV-RT1",
		InvoiceDate:     mustTime(t, "2026-07-21"),
		SupplierName:    "Cong ty Roundtrip",
		SupplierTaxCode: "111222333",
		Currency:        "USD",
		ExchangeRate:    25000,
		Subtotal:        2000,
		TaxAmount:       100,
		TotalAmount:     2100,
		DiscountAmount:  0,
		Lines: []domain.SupplierInvoiceLine{
			{ItemCode: "A1", ItemName: "Item A", Unit: "pcs", Quantity: 2, UnitPrice: 1000, VATRate: 5, VATType: domain.VAT5, LineTotal: 2000, LineVATAmount: 100, AccountID: "152", VATAccountID: "1331"},
		},
	}
	raw, err := Generate(inv)
	require.NoError(t, err)
	assert.Contains(t, string(raw), `xmlns="http://gdt.gov.vn/schemas/einvoice/2026"`)
	assert.Contains(t, string(raw), "<InvoiceNumber>INV-RT1</InvoiceNumber>")

	back, err := Parse(raw)
	require.NoError(t, err)
	assert.Equal(t, inv.InvoiceNumber, back.InvoiceNumber)
	assert.Equal(t, inv.SupplierName, back.SupplierName)
	assert.Equal(t, inv.SupplierTaxCode, back.SupplierTaxCode)
	assert.Equal(t, inv.Currency, back.Currency)
	assert.Equal(t, inv.ExchangeRate, back.ExchangeRate)
	require.Len(t, back.Lines, 1)
	assert.Equal(t, inv.Lines[0].ItemName, back.Lines[0].ItemName)
	assert.Equal(t, domain.VAT5, back.Lines[0].VATType)
	assert.Equal(t, 2000.0, back.Subtotal)
	assert.Equal(t, 2100.0, back.TotalAmount)
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	dd, err := time.Parse("2006-01-02", s)
	require.NoError(t, err)
	return dd
}
