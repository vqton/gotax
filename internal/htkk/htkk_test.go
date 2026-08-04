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

func TestGenerate_UnsupportedType(t *testing.T) {
	_, err := Generate(testDecl(domain.DeclTypeTTDB01), testCompany())
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
