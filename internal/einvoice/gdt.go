// Package einvoice implements GDT (Decree 254/2026) e-invoice XML
// generation and parsing for supplier invoices.
//
// Format authority: docs/purchase/PURCHASE_TEMPLATES.md §8.
package einvoice

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"gotax/internal/domain"
)

const (
	gdtNamespace = "http://gdt.gov.vn/schemas/einvoice/2026"
	dateLayout   = "2006-01-02"
	timeLayout   = "15:04:05"
)

// Default GL accounts applied to parsed lines. GDT payloads carry no account
// mapping; the ledger assignment is an application concern.
const (
	defaultAccountID    = "152"
	defaultVATAccountID = "1331"
)

// GDTInvoice mirrors the GDT e-invoice XML schema.
type GDTInvoice struct {
	XMLName          xml.Name             `xml:"Invoice"`
	Xmlns            string               `xml:"xmlns,attr"`
	InvoiceHeader    GDTInvoiceHeader     `xml:"InvoiceHeader"`
	Seller           GDTParty             `xml:"Seller"`
	Buyer            GDTParty             `xml:"Buyer"`
	InvoiceLines     GDTInvoiceLines      `xml:"InvoiceLines"`
	Summary          GDTSummary           `xml:"Summary"`
	DigitalSignature GDTDigitalSignature  `xml:"DigitalSignature,omitempty"`
	QRCode           string               `xml:"QRCode,omitempty"`
}

type GDTInvoiceHeader struct {
	InvoiceNumber string  `xml:"InvoiceNumber"`
	InvoiceDate   string  `xml:"InvoiceDate"`
	InvoiceTime   string  `xml:"InvoiceTime,omitempty"`
	InvoiceType   string  `xml:"InvoiceType"`
	InvoiceSeries string  `xml:"InvoiceSeries,omitempty"`
	CurrencyCode  string  `xml:"CurrencyCode"`
	ExchangeRate  float64 `xml:"ExchangeRate,omitempty"`
}

type GDTParty struct {
	TaxCode     string `xml:"TaxCode"`
	Name        string `xml:"Name"`
	Address     string `xml:"Address,omitempty"`
	Phone       string `xml:"Phone,omitempty"`
	BankAccount string `xml:"BankAccount,omitempty"`
	BankName    string `xml:"BankName,omitempty"`
}

type GDTInvoiceLines struct {
	Lines []GDTLine `xml:"Line"`
}

type GDTLine struct {
	LineNumber int     `xml:"LineNumber"`
	ItemCode   string  `xml:"ItemCode,omitempty"`
	ItemName   string  `xml:"ItemName"`
	Unit       string  `xml:"Unit,omitempty"`
	Quantity   float64 `xml:"Quantity"`
	UnitPrice  float64 `xml:"UnitPrice"`
	VatRate    float64 `xml:"VatRate"`
	VatAmount  float64 `xml:"VatAmount"`
	LineTotal  float64 `xml:"LineTotal"`
}

type GDTSummary struct {
	Subtotal      float64 `xml:"Subtotal"`
	Discount      float64 `xml:"Discount"`
	VatAmount     float64 `xml:"VatAmount"`
	GrandTotal    float64 `xml:"GrandTotal"`
	AmountInWords string  `xml:"AmountInWords,omitempty"`
}

type GDTDigitalSignature struct {
	Signature        string `xml:"Signature"`
	SignDate         string `xml:"SignDate,omitempty"`
	CertificateSerial string `xml:"CertificateSerial,omitempty"`
}

// Parse converts a GDT e-invoice XML document into a draft SupplierInvoice.
// GL accounts default to 152/1331 on each line.
func Parse(raw []byte) (*domain.SupplierInvoice, error) {
	var g GDTInvoice
	if err := xml.Unmarshal(raw, &g); err != nil {
		return nil, err
	}
	if g.InvoiceHeader.InvoiceNumber == "" {
		return nil, domain.ErrInvoiceNumberRequired
	}
	if g.Seller.Name == "" {
		return nil, domain.ErrInvoiceSupplierNameRequired
	}
	if g.Seller.TaxCode == "" {
		return nil, domain.ErrInvoiceSupplierTaxCodeRequired
	}
	if len(g.InvoiceLines.Lines) == 0 {
		return nil, domain.ErrInvoiceLinesRequired
	}
	invoiceDate, err := parseDate(g.InvoiceHeader.InvoiceDate, g.InvoiceHeader.InvoiceTime)
	if err != nil {
		return nil, domain.ErrInvoiceDateRequired
	}
	currency := g.InvoiceHeader.CurrencyCode
	if currency == "" {
		currency = "VND"
	}
	inv := &domain.SupplierInvoice{
		InvoiceNumber:   g.InvoiceHeader.InvoiceNumber,
		InvoiceDate:     invoiceDate,
		SupplierName:    g.Seller.Name,
		SupplierTaxCode: g.Seller.TaxCode,
		Currency:        currency,
		ExchangeRate:    g.InvoiceHeader.ExchangeRate,
		Subtotal:        g.Summary.Subtotal,
		DiscountAmount:  g.Summary.Discount,
		TaxAmount:       g.Summary.VatAmount,
		TotalAmount:     g.Summary.GrandTotal,
	}
	if inv.TotalAmount == 0 {
		inv.TotalAmount = inv.Subtotal + inv.TaxAmount
	}
	if g.InvoiceHeader.InvoiceType == domain.InvoiceTypeCreditNote {
		inv.InvoiceType = domain.InvoiceTypeCreditNote
	}
	inv.Lines = make([]domain.SupplierInvoiceLine, 0, len(g.InvoiceLines.Lines))
	for _, l := range g.InvoiceLines.Lines {
		line := domain.SupplierInvoiceLine{
			ItemCode:      l.ItemCode,
			ItemName:      l.ItemName,
			Unit:          l.Unit,
			Quantity:      l.Quantity,
			UnitPrice:     l.UnitPrice,
			VATRate:       l.VatRate,
			VATType:       vatRateToType(l.VatRate),
			LineTotal:     l.LineTotal,
			LineVATAmount: l.VatAmount,
			AccountID:     defaultAccountID,
			VATAccountID:  defaultVATAccountID,
		}
		if line.LineTotal == 0 {
			line.LineTotal = l.Quantity * l.UnitPrice
		}
		inv.Lines = append(inv.Lines, line)
	}
	return inv, nil
}

// Generate renders a SupplierInvoice as a GDT e-invoice XML document.
func Generate(inv *domain.SupplierInvoice) ([]byte, error) {
	if inv == nil {
		return nil, errors.New("einvoice: nil invoice")
	}
	currency := inv.Currency
	if currency == "" {
		currency = "VND"
	}
	g := GDTInvoice{
		Xmlns: gdtNamespace,
		InvoiceHeader: GDTInvoiceHeader{
			InvoiceNumber: inv.InvoiceNumber,
			InvoiceDate:   inv.InvoiceDate.Format(dateLayout),
			InvoiceTime:   inv.InvoiceDate.Format(timeLayout),
			InvoiceType:   gdtInvoiceType(inv.InvoiceType),
			CurrencyCode:  currency,
			ExchangeRate:  inv.ExchangeRate,
		},
		Seller: GDTParty{
			TaxCode: inv.SupplierTaxCode,
			Name:    inv.SupplierName,
		},
		Buyer: GDTParty{},
		InvoiceLines: GDTInvoiceLines{
			Lines: make([]GDTLine, 0, len(inv.Lines)),
		},
		Summary: GDTSummary{
			Subtotal:   inv.Subtotal,
			Discount:   inv.DiscountAmount,
			VatAmount:  inv.TaxAmount,
			GrandTotal: inv.TotalAmount,
		},
	}
	for i, l := range inv.Lines {
		g.InvoiceLines.Lines = append(g.InvoiceLines.Lines, GDTLine{
			LineNumber: i + 1,
			ItemCode:   l.ItemCode,
			ItemName:   l.ItemName,
			Unit:       l.Unit,
			Quantity:   l.Quantity,
			UnitPrice:  l.UnitPrice,
			VatRate:    vatTypeToRate(l.VATType, l.VATRate),
			VatAmount:  l.LineVATAmount,
			LineTotal:  l.LineTotal,
		})
	}
	if g.Summary.GrandTotal == 0 {
		g.Summary.GrandTotal = inv.Subtotal + inv.TaxAmount
	}
	return xml.MarshalIndent(g, "", "  ")
}

func vatRateToType(rate float64) domain.VATType {
	switch rate {
	case 0:
		return domain.VAT0
	case 5:
		return domain.VAT5
	case 8:
		return domain.VAT8
	case -1:
		return domain.VATNonTaxable
	case 10:
		return domain.VAT10
	default:
		return domain.VAT10
	}
}

func vatTypeToRate(t domain.VATType, fallback float64) float64 {
	switch t {
	case domain.VAT0:
		return 0
	case domain.VAT5:
		return 5
	case domain.VAT8:
		return 8
	case domain.VAT10:
		return 10
	case domain.VATNonTaxable:
		return -1
	default:
		if fallback != 0 {
			return fallback
		}
		return 10
	}
}

func gdtInvoiceType(invType string) string {
	if invType == domain.InvoiceTypeCreditNote {
		return "credit_note"
	}
	if invType != "" {
		return invType
	}
	return "01GTKT"
}

func parseDate(date, clock string) (time.Time, error) {
	value := date
	if clock != "" {
		value = date + "T" + clock
	}
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid invoice date %q", strings.TrimSpace(date))
}
