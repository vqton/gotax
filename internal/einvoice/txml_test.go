package einvoice

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
)

func testEInvoice() *domain.EInvoice {
	return &domain.EInvoice{
		ID: "EINV-1", CompanyID: "c1",
		Pattern: "01GTKT0/001", Serial: "AA/26E", InvoiceNumber: 1,
		InvoiceType:  domain.EInvTypeORIGINAL,
		BuyerName:    "CONG TY TNHH ABC",
		BuyerTaxCode: "0987654321",
		BuyerAddress: "123 Nguyen Hue, Q1, HCMC",
		CurrencyCode: "VND",
		IssueDate:    "2026-04-15",
		Subtotal:     100000000, VATAmount: 10000000, GrandTotal: 110000000,
		Lines: []domain.EInvoiceLine{{
			LineNumber: 1, Description: "May tinh xach tay", Unit: "Chiec",
			Quantity: 5, UnitPrice: 20000000, LineTotal: 100000000,
			VATRate: 10, VATAmount: 10000000,
		}},
	}
}

func TestGenerateTXML(t *testing.T) {
	out, err := GenerateTXML(testEInvoice())
	require.NoError(t, err)
	s := string(out)

	assert.True(t, strings.HasPrefix(s, `<?xml version="1.0" encoding="UTF-8"?>`))
	assert.Contains(t, s, `<BK:HoaDon xmlns:BK="http://gdt.gov.vn/HoaDon">`)
	assert.Contains(t, s, `<BK:MauSo>01GTKT0/001</BK:MauSo>`)
	assert.Contains(t, s, `<BK:KyHieu>AA/26E</BK:KyHieu>`)
	assert.Contains(t, s, `<BK:NgayLap>2026-04-15</BK:NgayLap>`)
	assert.Contains(t, s, `<BK:Ten>CONG TY TNHH ABC</BK:Ten>`)
	assert.Contains(t, s, `<BK:MaSoThue>0987654321</BK:MaSoThue>`)
	assert.Contains(t, s, `<BK:SoLuong>5</BK:SoLuong>`)
	assert.Contains(t, s, `<BK:TongTienThanhToan>110000000</BK:TongTienThanhToan>`)
	assert.NotContains(t, s, "e+0") // plain decimal, never scientific

	// well-formed XML (token walk)
	dec := xml.NewDecoder(strings.NewReader(s))
	for {
		_, err := dec.Token()
		if err != nil {
			require.ErrorIs(t, err, io.EOF)
			break
		}
	}

	// deterministic
	out2, err := GenerateTXML(testEInvoice())
	require.NoError(t, err)
	assert.Equal(t, string(out), string(out2))
}

func TestGenerateTXML_NilAndNoLines(t *testing.T) {
	_, err := GenerateTXML(nil)
	assert.Error(t, err)
	inv := testEInvoice()
	inv.Lines = nil
	_, err = GenerateTXML(inv)
	assert.Error(t, err)
}

func TestGenerateTXML_ZeroLineNumberFallsBack(t *testing.T) {
	inv := testEInvoice()
	inv.Lines[0].LineNumber = 0
	out, err := GenerateTXML(inv)
	require.NoError(t, err)
	assert.Contains(t, string(out), `<BK:STT>1</BK:STT>`)
}
