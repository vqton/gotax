// Package htkk generates tax declaration XML per the HTKK (Hệ thống
// Khai thuế qua mạng) envelope documented in docs/tax/TAX_SPECS.md §3.1
// (Circular 80/2021/TT-BTC). Pure encoding/xml — no HTTP, no storage.
package htkk

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gotax/internal/domain"
)

// money marshals as plain decimal (never scientific notation).
type money float64

func (m money) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(strconv.FormatFloat(float64(m), 'f', -1, 64), start)
}

// formCode maps declaration types to HTKK form codes.
func formCode(t domain.DeclarationType) string {
	switch t {
	case domain.DeclTypeGTGT01:
		return "01/GTGT"
	case domain.DeclTypeTNDN03:
		return "03/TNDN"
	default:
		return ""
	}
}

// BoKe is the HTKK file root (TAX_SPECS §3.1).
type BoKe struct {
	XMLName       xml.Name        `xml:"BK:BoKe"`
	XmlnsBK       string          `xml:"xmlns:BK,attr"`
	ThongTinChung BKThongTinChung `xml:"BK:ThongTinChung"`
	DuLieu        BKDuLieu        `xml:"BK:DuLieu"`
}

type BKThongTinChung struct {
	MaSoThue       string `xml:"BK:MaSoThue"`
	TenNguoiNopThue string `xml:"BK:TenNguoiNopThue"`
	LoaiToKhai     string `xml:"BK:LoaiToKhai"`
	KyTinhThue     string `xml:"BK:KyTinhThue"`
	LanDau         int    `xml:"BK:LanDau"`
	NgayTao        string `xml:"BK:NgayTao"`
}

type BKDuLieu struct {
	ChiTieu []BKChiTieu `xml:"BK:ChiTieu"`
}

type BKChiTieu struct {
	MaChiTieu string `xml:"BK:MaChiTieu"`
	GiaTri    money  `xml:"BK:GiaTri"`
}

// periodKey renders "2026Q1" / "2026M3" / "2026" per the form's tax period.
func periodKey(p domain.TaxPeriod) string {
	switch p.PeriodType {
	case domain.PeriodTypeQuarterly:
		return fmt.Sprintf("%dQ%d", p.PeriodYear, p.PeriodNumber)
	case domain.PeriodTypeAnnual:
		return fmt.Sprintf("%d", p.PeriodYear)
	default: // monthly / per-occurrence fall back to month
		return fmt.Sprintf("%dM%d", p.PeriodYear, p.PeriodNumber)
	}
}

// declarationDate renders the creation timestamp as a date-only value, per
// spec §3.1 the HTKK XSD `date` type rejects time components.
func declarationDate(createdAt string) (string, error) {
	for _, layout := range []string{time.RFC3339, time.DateOnly} {
		t, err := time.Parse(layout, createdAt)
		if err == nil {
			return t.Format(time.DateOnly), nil
		}
	}
	return "", fmt.Errorf("unrecognized timestamp %q", createdAt)
}

// indicatorCode renders an indicator code as bracketed ([10], [23]) per the
// form layout. Domain LineCodes are stored unbracketed; idempotent for input
// that already carries brackets.
func indicatorCode(code string) string {
	c := strings.Trim(code, "[]")
	if c == "" {
		return c
	}
	return "[" + c + "]"
}

// Generate renders the declaration file. Header taxpayer data comes from the
// company profile; indicator values from declaration lines ([10], [20a], …).
func Generate(d *domain.TaxDeclaration, c *domain.Company) ([]byte, error) {
	code := formCode(d.DeclarationType)
	if code == "" {
		return nil, fmt.Errorf("htkk: unsupported declaration type %q", d.DeclarationType)
	}
	if c == nil || c.TaxCode == "" {
		return nil, fmt.Errorf("htkk: company tax code required")
	}
	ngayTao, err := declarationDate(d.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("htkk: invalid declaration CreatedAt: %w", err)
	}
	f := BoKe{
		XmlnsBK: "http://gdt.gov.vn/BoKe",
		ThongTinChung: BKThongTinChung{
			MaSoThue:        c.TaxCode,
			TenNguoiNopThue: c.LegalNameVN,
			LoaiToKhai:      code,
			KyTinhThue:      periodKey(d.TaxPeriod),
			LanDau:          1,
			NgayTao:         ngayTao,
		},
	}
	if d.AdjustmentType != domain.AdjTypeNONE {
		f.ThongTinChung.LanDau = 2 // supplemental/amended filing
	}
	for _, l := range d.Lines {
		if mc := indicatorCode(l.LineCode); mc != "" {
			f.DuLieu.ChiTieu = append(f.DuLieu.ChiTieu, BKChiTieu{
				MaChiTieu: mc,
				GiaTri:    money(l.Amount),
			})
		}
	}
	return xml.MarshalIndent(f, "", "  ")
}
