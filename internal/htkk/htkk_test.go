package htkk

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
)

func testCompany() *domain.Company {
	return &domain.Company{
		TaxCode:     "0100123456",
		LegalNameVN: "CONG TY TNHH ABC",
	}
}

// local-name struct: xml.Unmarshal with no namespace URI matches by local
// name only, so <BK:BoKe> (URI http://gdt.gov.vn/BoKe) parses into this.
type parsedFile struct {
	XMLName xml.Name `xml:"BoKe"`
	TTChung struct {
		MaSoThue        string  `xml:"MaSoThue"`
		TenNguoiNopThue string  `xml:"TenNguoiNopThue"`
		LoaiToKhai      string  `xml:"LoaiToKhai"`
		KyTinhThue      string  `xml:"KyTinhThue"`
		LanDau          int     `xml:"LanDau"`
		NgayTao         string  `xml:"NgayTao"`
	} `xml:"ThongTinChung"`
	DuLieu struct {
		ChiTieu []struct {
			MaChiTieu string  `xml:"MaChiTieu"`
			GiaTri    float64 `xml:"GiaTri"`
		} `xml:"ChiTieu"`
	} `xml:"DuLieu"`
}

func parse(t *testing.T, raw []byte) parsedFile {
	t.Helper()
	var f parsedFile
	require.NoError(t, xml.Unmarshal(raw, &f))
	return f
}

func testDecl(t domain.DeclarationType) *domain.TaxDeclaration {
	return &domain.TaxDeclaration{
		ID:              "DECL-1",
		CompanyID:       "c1",
		DeclarationType: t,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 1},
		Status:          domain.DeclStatusVALIDATED,
		AdjustmentType:  domain.AdjTypeNONE,
		CreatedAt:       "2026-04-15T10:00:00+07:00",
		Lines: []domain.TaxDeclarationLine{
			{LineCode: "10", LineName: "Hang hoa, dich vu ban ra", Amount: 500000000},
			{LineCode: "25", LineName: "Thue GTGT dau ra", Amount: 50000000},
			{LineCode: "33", LineName: "Thue GTGT dau vao", Amount: 30000000},
			{LineCode: "40", LineName: "Thue GTGT phai nop", Amount: 20000000},
		},
	}
}

func TestGenerate_GTGT01(t *testing.T) {
	out, err := Generate(testDecl(domain.DeclTypeGTGT01), testCompany())
	require.NoError(t, err)

	f := parse(t, out)
	assert.Contains(t, string(out), `xmlns:BK="http://gdt.gov.vn/BoKe"`)
	assert.Equal(t, "0100123456", f.TTChung.MaSoThue)
	assert.Equal(t, "CONG TY TNHH ABC", f.TTChung.TenNguoiNopThue)
	assert.Equal(t, "01/GTGT", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026Q1", f.TTChung.KyTinhThue)
	assert.Equal(t, 1, f.TTChung.LanDau)
	// spec §3.1: NgayTao is a date-only field, indicator codes are bracketed
	assert.Equal(t, "2026-04-15", f.TTChung.NgayTao)
	assert.Len(t, f.DuLieu.ChiTieu, 4)
	assert.Equal(t, "[10]", f.DuLieu.ChiTieu[0].MaChiTieu)
	assert.Equal(t, 500000000.0, f.DuLieu.ChiTieu[0].GiaTri)
}

func TestGenerate_BracketedInputIdempotent(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT01)
	d.Lines = []domain.TaxDeclarationLine{{LineCode: "[10]", Amount: 1000000}}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	require.Len(t, f.DuLieu.ChiTieu, 1)
	assert.Equal(t, "[10]", f.DuLieu.ChiTieu[0].MaChiTieu)
}

func TestGenerate_InvalidNgayTao(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT01)
	d.CreatedAt = "not-a-timestamp"
	_, err := Generate(d, testCompany())
	assert.Error(t, err)
}

func TestGenerate_TNDN03(t *testing.T) {
	out, err := Generate(testDecl(domain.DeclTypeTNDN03), testCompany())
	require.NoError(t, err)
	assert.Contains(t, string(out), "03/TNDN")
}

func TestGenerate_MonthlyPeriodKey(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT01)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 3}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	assert.Contains(t, string(out), "<BK:KyTinhThue>2026M3</BK:KyTinhThue>")
}

func TestGenerate_AmendmentSetsLanDau(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT01)
	d.AdjustmentType = domain.AdjTypeAMENDMENT
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, 2, f.TTChung.LanDau)
}

func TestGenerate_GTGT02(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT02)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 6}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "02/GTGT", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026M6", f.TTChung.KyTinhThue)
	assert.Len(t, f.DuLieu.ChiTieu, 4)
}

func TestGenerate_TTDB01(t *testing.T) {
	d := testDecl(domain.DeclTypeTTDB01)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 3}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "01", LineName: "San luong", Amount: 10000},
		{LineCode: "02", LineName: "Gia ban", Amount: 50000},
		{LineCode: "03", LineName: "Tyle (%)", Amount: 10},
		{LineCode: "04", LineName: "Thue phai nop", Amount: 50000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "TTDB/01", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026M3", f.TTChung.KyTinhThue)
	assert.Len(t, f.DuLieu.ChiTieu, 4)
}

func TestGenerate_BVMT01(t *testing.T) {
	d := testDecl(domain.DeclTypeBVMT01)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 3}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "01", LineName: "San luong", Amount: 5000},
		{LineCode: "02", LineName: "Dinh muc", Amount: 1000},
		{LineCode: "03", LineName: "Thue suat", Amount: 50000},
		{LineCode: "04", LineName: "Thue phai nop", Amount: 250000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "BVMT/01", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026M3", f.TTChung.KyTinhThue)
	assert.Len(t, f.DuLieu.ChiTieu, 4)
}

func TestGenerate_NTN01(t *testing.T) {
	d := testDecl(domain.DeclTypeNTNN01)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 3}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "01", LineName: "Thu nhap", Amount: 200000000},
		{LineCode: "02", LineName: "Thue suat", Amount: 5},
		{LineCode: "03", LineName: "Thue phai nop", Amount: 10000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "NTNN/01", f.TTChung.LoaiToKhai)
	assert.Len(t, f.DuLieu.ChiTieu, 3)
}

func TestGenerate_QTTTNCN(t *testing.T) {
	d := testDecl(domain.DeclTypeQTTTNCN)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 2}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "01", LineName: "So nguoi lao dong", Amount: 30},
		{LineCode: "02", LineName: "Tong thu nhap", Amount: 600000000},
		{LineCode: "03", LineName: "Tong khau tru", Amount: 60000000},
		{LineCode: "04", LineName: "Tong thue TNCN", Amount: 54000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "QTT/TNCN", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026Q2", f.TTChung.KyTinhThue)
	assert.Len(t, f.DuLieu.ChiTieu, 4)
}

func TestGenerate_KKTNCN(t *testing.T) {
	d := testDecl(domain.DeclTypeKKTNCN)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 6}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "01", LineName: "So nguoi lao dong", Amount: 25},
		{LineCode: "02", LineName: "Tong thu nhap", Amount: 500000000},
		{LineCode: "03", LineName: "Tong khau tru", Amount: 50000000},
		{LineCode: "04", LineName: "Tong thue TNCN", Amount: 45000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "KK/TNCN", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026M6", f.TTChung.KyTinhThue)
	assert.Len(t, f.DuLieu.ChiTieu, 4)
}

func TestGenerate_TNDN05(t *testing.T) {
	d := testDecl(domain.DeclTypeTNDN05)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeAnnual, PeriodYear: 2026, PeriodNumber: 0}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "04", LineName: "Tong doanh thu", Amount: 800000000},
		{LineCode: "06", LineName: "Doanh thu tinh thue", Amount: 800000000},
		{LineCode: "12", LineName: "Thu nhap chiu thue", Amount: 80000000},
		{LineCode: "13", LineName: "Thue suat (%)", Amount: 20},
		{LineCode: "14", LineName: "Thue TNDN phai nop", Amount: 16000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "05/TNDN", f.TTChung.LoaiToKhai)
	assert.Len(t, f.DuLieu.ChiTieu, 5)
}

func TestGenerate_TNDN06(t *testing.T) {
	d := testDecl(domain.DeclTypeTNDN06)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeAnnual, PeriodYear: 2026, PeriodNumber: 0}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "04", LineName: "Tong doanh thu", Amount: 5000000000},
		{LineCode: "06", LineName: "Doanh thu tinh thue", Amount: 5000000000},
		{LineCode: "12", LineName: "Thu nhap chiu thue", Amount: 500000000},
		{LineCode: "13", LineName: "Thue suat (%)", Amount: 50},
		{LineCode: "14", LineName: "Thue TNDN phai nop", Amount: 250000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "06/TNDN", f.TTChung.LoaiToKhai)
	assert.Len(t, f.DuLieu.ChiTieu, 5)
}

func TestGenerate_TNDN04(t *testing.T) {
	d := testDecl(domain.DeclTypeTNDN04)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeAnnual, PeriodYear: 2026, PeriodNumber: 0}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "04", LineName: "Tong doanh thu", Amount: 2000000000},
		{LineCode: "06", LineName: "Doanh thu tinh thue", Amount: 2000000000},
		{LineCode: "12", LineName: "Thu nhap chiu thue", Amount: 200000000},
		{LineCode: "13", LineName: "Thue suat (%)", Amount: 20},
		{LineCode: "14", LineName: "Thue TNDN phai nop", Amount: 40000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "04/TNDN", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026", f.TTChung.KyTinhThue)
	assert.Len(t, f.DuLieu.ChiTieu, 5)
}

func TestGenerate_TNDN02(t *testing.T) {
	d := testDecl(domain.DeclTypeTNDN02)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 2}
	d.Lines = []domain.TaxDeclarationLine{
		{LineCode: "04", LineName: "Tong doanh thu", Amount: 500000000},
		{LineCode: "06", LineName: "Doanh thu tinh thue", Amount: 500000000},
		{LineCode: "12", LineName: "Thu nhap chiu thue", Amount: 50000000},
		{LineCode: "13", LineName: "Thue suat (%)", Amount: 20},
		{LineCode: "14", LineName: "Thue TNDN phai nop", Amount: 10000000},
	}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "02/TNDN", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026Q2", f.TTChung.KyTinhThue)
	assert.Len(t, f.DuLieu.ChiTieu, 5)
}

func TestGenerate_GTGT05(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT05)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "05/GTGT", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026M1", f.TTChung.KyTinhThue)
}

func TestGenerate_GTGT04(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT04)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 3}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "04/GTGT", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026M3", f.TTChung.KyTinhThue)
}

func TestGenerate_GTGT03(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT03)
	d.TaxPeriod = domain.TaxPeriod{PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 2}
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	f := parse(t, out)
	assert.Equal(t, "03/GTGT", f.TTChung.LoaiToKhai)
	assert.Equal(t, "2026Q2", f.TTChung.KyTinhThue)
	assert.Len(t, f.DuLieu.ChiTieu, 4)
}

func TestGenerate_UnsupportedType(t *testing.T) {
	_, err := Generate(testDecl("INVALID"), testCompany())
	assert.Error(t, err)
}

func TestGenerate_MissingCompany(t *testing.T) {
	_, err := Generate(testDecl(domain.DeclTypeGTGT01), nil)
	assert.Error(t, err)
}

func TestGenerate_NoScientificNotation(t *testing.T) {
	d := testDecl(domain.DeclTypeGTGT01)
	d.Lines = []domain.TaxDeclarationLine{{LineCode: "[10]", Amount: 100000000000}} // 1e11
	out, err := Generate(d, testCompany())
	require.NoError(t, err)
	assert.False(t, strings.Contains(string(out), "e+"))
	assert.Contains(t, string(out), ">100000000000<")
}
