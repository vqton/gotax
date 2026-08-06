package einvoice

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"gotax/internal/domain"
)

// money marshals as plain decimal (never scientific notation) — tax amounts
// must read as exact VND figures.
type money float64

func (m money) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(strconv.FormatFloat(float64(m), 'f', -1, 64), start)
}

// Outbound TXML per TAX_TEMPLATES.md §5 (Decree 254/2026/ND-CP).
// Root BK namespace; signature (BK:ChuKySo) injected by the signer
// before BK:KyThuat closes. Seller tax code comes from the company
// profile (future) — omitted for now.
type BKHoaDon struct {
	XMLName       xml.Name        `xml:"BK:HoaDon"`
	XmlnsBK       string          `xml:"xmlns:BK,attr"`
	ThongTinChung BKThongTinChung `xml:"BK:ThongTinChung"`
	BenMua        BKParty         `xml:"BK:BenMua"`
	DanhSachHang  BKHangHoaList   `xml:"BK:DanhSachHangHoa"`
	TongTien      BKTongTien      `xml:"BK:TongTien"`
	KyThuat       BKKyThuat       `xml:"BK:KyThuat"`
}

type BKThongTinChung struct {
	MaSoThue    string `xml:"BK:MaSoThue,omitempty"`
	MauSo       string `xml:"BK:MauSo"`
	KyHieu      string `xml:"BK:KyHieu"`
	NgayLap     string `xml:"BK:NgayLap"`
	ThoiDiemLap string `xml:"BK:ThoiDiemLap,omitempty"`
	LoaiHoaDon  string `xml:"BK:LoaiHoaDon"`
	DonViTienTe string `xml:"BK:DonViTienTe"`
}

type BKParty struct {
	Ten      string `xml:"BK:Ten"`
	MaSoThue string `xml:"BK:MaSoThue,omitempty"`
	DiaChi   string `xml:"BK:DiaChi,omitempty"`
	Email    string `xml:"BK:Email,omitempty"`
}

type BKHangHoaList struct {
	Items []BKHangHoa `xml:"BK:HangHoa"`
}

type BKHangHoa struct {
	STT       int    `xml:"BK:STT"`
	TenHang   string `xml:"BK:TenHang"`
	DonViTinh string `xml:"BK:DonViTinh,omitempty"`
	SoLuong   money  `xml:"BK:SoLuong"`
	DonGia    money  `xml:"BK:DonGia"`
	ThanhTien money  `xml:"BK:ThanhTien"`
	ThueSuat  money  `xml:"BK:ThueSuat"`
	ThueGTGT  money  `xml:"BK:ThueGTGT"`
}

type BKTongTien struct {
	TongTienHang      money  `xml:"BK:TongTienHang"`
	TongThueGTGT      money  `xml:"BK:TongThueGTGT"`
	TongTienThanhToan money  `xml:"BK:TongTienThanhToan"`
	SoTienBangChu     string `xml:"BK:SoTienBangChu,omitempty"`
}

type BKKyThuat struct {
	ChuKySo  *BKChuKySo `xml:"BK:ChuKySo,omitempty"`
	MaCuaCQT string     `xml:"BK:MaCuaCQT,omitempty"`
}

type BKChuKySo struct {
	SerialNumber string `xml:"BK:SerialNumber"`
	ThoiDiemKy   string `xml:"BK:ThoiDiemKy"`
	DuLieuKy     string `xml:"BK:DuLieuKy"`
}

const bkNamespace = "http://gdt.gov.vn/HoaDon"

// GenerateTXML renders an outbound invoice as GDT TXML (BK:HoaDon).
// Deterministic output — safe to sign byte-for-byte after canonicalization.
func GenerateTXML(inv *domain.EInvoice) ([]byte, error) {
	if inv == nil {
		return nil, fmt.Errorf("einvoice: nil invoice")
	}
	items := make([]BKHangHoa, 0, len(inv.Lines))
	for i, l := range inv.Lines {
		stt := l.LineNumber
		if stt == 0 {
			stt = i + 1
		}
		items = append(items, BKHangHoa{
			STT:       stt,
			TenHang:   l.Description,
			DonViTinh: l.Unit,
			SoLuong:   money(l.Quantity),
			DonGia:    money(l.UnitPrice),
			ThanhTien: money(l.LineTotal),
			ThueSuat:  money(l.VATRate),
			ThueGTGT:  money(l.VATAmount),
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("einvoice: invoice has no lines")
	}
	doc := BKHoaDon{
		XmlnsBK: bkNamespace,
		ThongTinChung: BKThongTinChung{
			MauSo:       inv.Pattern,
			KyHieu:      inv.Serial,
			NgayLap:     inv.IssueDate,
			LoaiHoaDon:  "BAN_HANG",
			DonViTienTe: inv.CurrencyCode,
		},
		BenMua: BKParty{
			Ten:      inv.BuyerName,
			MaSoThue: inv.BuyerTaxCode,
			DiaChi:   inv.BuyerAddress,
			Email:    inv.BuyerEmail,
		},
		DanhSachHang: BKHangHoaList{Items: items},
		TongTien: BKTongTien{
			TongTienHang:      money(inv.Subtotal),
			TongThueGTGT:      money(inv.VATAmount),
			TongTienThanhToan: money(inv.GrandTotal),
			SoTienBangChu:     AmountInWords(int64(inv.GrandTotal)),
		},
	}
	body, err := xml.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("einvoice: marshal TXML: %w", err)
	}
	out := append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n"), body...)
	return out, nil
}
