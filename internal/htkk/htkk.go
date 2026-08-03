// Package htkk generates tax declaration XML per the HTKK (Hệ thống
// Khai thuế qua mạng) envelope documented in docs/tax/TAX_SPECS.md §3.1
// (Circular 80/2021/TT-BTC). Pure encoding/xml — no HTTP, no storage.
package htkk

import (
	"encoding/xml"
	"fmt"
	"strconv"

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
	f := BoKe{
		XmlnsBK: "http://gdt.gov.vn/BoKe",
		ThongTinChung: BKThongTinChung{
			MaSoThue:        c.TaxCode,
			TenNguoiNopThue: c.LegalNameVN,
			LoaiToKhai:      code,
			KyTinhThue:      periodKey(d.TaxPeriod),
			LanDau:          1,
			NgayTao:         d.CreatedAt,
		},
	}
	if d.AdjustmentType != domain.AdjTypeNONE {
		f.ThongTinChung.LanDau = 2 // supplemental/amended filing
	}
	for _, l := range d.Lines {
		f.DuLieu.ChiTieu = append(f.DuLieu.ChiTieu, BKChiTieu{
			MaChiTieu: l.LineCode,
			GiaTri:    money(l.Amount),
		})
	}
	return xml.MarshalIndent(f, "", "  ")
}
