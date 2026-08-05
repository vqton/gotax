// Package payroll generates Vietnamese social insurance and PIT declaration XML.
// D02-TS: Social insurance declaration (Decision 366/QĐ-BHXH)
// 05/KK-TNCN: Quarterly PIT declaration (Decision 1109/QD-BTC)
// TK3-TS: Employer registration form (Circular 590/2017/TT-BTC)
package payroll

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"gotax/internal/domain"
)

// ─── CSV Import ─────────────────────────────────────────────────

// ParseTimekeepingCSV reads a CSV file and returns timekeeping records.
// Format: employee_code,date,clock_in,clock_out,ot_hours,night_hours,leave_type,notes
func ParseTimekeepingCSV(r io.Reader) ([]domain.TimekeepingRecord, error) {
	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	colIdx := make(map[string]int)
	for i, h := range header {
		colIdx[strings.TrimSpace(h)] = i
	}
	needCols := []string{"employee_code", "date"}
	for _, c := range needCols {
		if _, ok := colIdx[c]; !ok {
			return nil, fmt.Errorf("missing required column: %s", c)
		}
	}
	var records []domain.TimekeepingRecord
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read row: %w", err)
		}
		get := func(key string) string {
			if i, ok := colIdx[key]; ok && i < len(row) {
				return strings.TrimSpace(row[i])
			}
			return ""
		}
		otHours, err := strconv.ParseFloat(get("ot_hours"), 64)
		if err != nil && get("ot_hours") != "" {
			return nil, fmt.Errorf("invalid ot_hours for %s: %w", get("employee_code"), err)
		}
		nightHours, err := strconv.ParseFloat(get("night_hours"), 64)
		if err != nil && get("night_hours") != "" {
			return nil, fmt.Errorf("invalid night_hours for %s: %w", get("employee_code"), err)
		}
		hoursWorked := 8.0
		if get("clock_in") != "" && get("clock_out") != "" {
			hoursWorked = computeHoursWorked(get("clock_in"), get("clock_out"))
		}
		records = append(records, domain.TimekeepingRecord{
			EmployeeID:  get("employee_code"),
			Date:        get("date"),
			ClockIn:     get("clock_in"),
			ClockOut:    get("clock_out"),
			HoursWorked: hoursWorked,
			OTHours:     otHours,
			NightHours:  nightHours,
			LeaveType:   domain.LeaveType(get("leave_type")),
			Notes:       get("notes"),
		})
	}
	return records, nil
}

// computeHoursWorked calculates hours between HH:MM clock_in and clock_out.
func computeHoursWorked(clockIn, clockOut string) float64 {
	inH, inM := 0, 0
	outH, outM := 0, 0
	fmt.Sscanf(clockIn, "%d:%d", &inH, &inM)
	fmt.Sscanf(clockOut, "%d:%d", &outH, &outM)
	minutes := outH*60 + outM - inH*60 - inM
	if minutes < 0 {
		minutes += 24 * 60 // overnight shift
	}
	return float64(minutes) / 60.0
}

// ─── D02-TS: Social Insurance Declaration ───────────────────────

type D02TS struct {
	XMLName          xml.Name       `xml:"D02TS"`
	XMLNS            string         `xml:"xmlns,attr"`
	MaDonVi          string         `xml:"MaDonVi"`
	TenDonVi         string         `xml:"TenDonVi"`
	MaSoThue         string         `xml:"MaSoThue"`
	Nam              int            `xml:"Nam"`
	Quy              int            `xml:"Quy"`
	DanhSachNhanVien []D02TSNhanVien `xml:"DanhSachNhanVien>NhanVien"`
}

type D02TSNhanVien struct {
	MaSoBHXH        string  `xml:"MaSoBHXH"`
	HoVaTen         string  `xml:"HoVaTen"`
	CMND            string  `xml:"CMND"`
	NamSinh         string  `xml:"NamSinh"`
	GioiTinh        string  `xml:"GioiTinh"`
	ChucVu          string  `xml:"ChucVu"`
	MucLuong        float64 `xml:"MucLuong"`
	TienLuongBH     float64 `xml:"TienLuongBH"`
	DongBHXH        float64 `xml:"DongBHXH"`
	DongBHTN        float64 `xml:"DongBHTN"`
	DongBHYT        float64 `xml:"DongBHYT"`
}

func GenerateD02TS(companyName, taxCode string, year, quarter int, employees []domain.EmployeePayrollInfo, runs []domain.PayrollRun) ([]byte, error) {
	d := D02TS{
		XMLNS:    "http://baohiemxahoi.gov.vn/D02TS",
		MaDonVi:  taxCode,
		TenDonVi: companyName,
		MaSoThue: taxCode,
		Nam:      year,
		Quy:      quarter,
	}
	runMap := make(map[string]domain.PayrollRun)
	for _, r := range runs {
		runMap[r.EmployeeID] = r
	}
	for _, emp := range employees {
		run := runMap[emp.EmployeeID]
		nv := D02TSNhanVien{
			HoVaTen:     emp.EmployeeID,
			NamSinh:     emp.ContractStartDate,
			MucLuong:    emp.BaseSalary,
			TienLuongBH: emp.InsuranceBaseSalary,
			DongBHXH:    run.SIDeduction + run.EmployerSI,
			DongBHTN:    run.UIDeduction + run.EmployerUI,
			DongBHYT:    run.HIDeduction + run.EmployerHI,
		}
		d.DanhSachNhanVien = append(d.DanhSachNhanVien, nv)
	}
	return xml.MarshalIndent(d, "", "  ")
}

// ─── 05/KK-TNCN: Quarterly PIT Declaration ─────────────────────

type KKTNCN struct {
	XMLName          xml.Name            `xml:"KK_TNCN"`
	XMLNS            string              `xml:"xmlns,attr"`
	MaSoThue         string              `xml:"MaSoThue"`
	TenDonVi         string              `xml:"TenDonVi"`
	KyTinhThue       string              `xml:"KyTinhThue"`
	DanhSachNhanVien []KKTNCNNhanVien    `xml:"DanhSachNhanVien>NhanVien"`
	TongCong         KKTNCNTongCong      `xml:"TongCong"`
}

type KKTNCNNhanVien struct {
	MaSoThueCN       string  `xml:"MaSoThueCN"`
	HoVaTen          string  `xml:"HoVaTen"`
	SoCMND           string  `xml:"SoCMND"`
	DiaChi           string  `xml:"DiaChi"`
	SoTaiKhoan       string  `xml:"SoTaiKhoan"`
	MaNganHang       string  `xml:"MaNganHang"`
	ThuNhapChiuThue float64 `xml:"ThuNhapChiuThue"`
	GiamTru          float64 `xml:"GiamTru"`
	ThueTNCN         float64 `xml:"ThueTNCN"`
}

// Use domain.PersonalDeductionMonthly (15.5M/month) for deductions.

type KKTNCNTongCong struct {
	TongThuNhap    float64 `xml:"TongThuNhap"`
	TongGiamTru    float64 `xml:"TongGiamTru"`
	TongThueTNCN   float64 `xml:"TongThueTNCN"`
}

func Generate05KKTNCN(taxCode, companyName, periodKey string, employees []domain.EmployeePayrollInfo, runs []domain.PayrollRun) ([]byte, error) {
	d := KKTNCN{
		XMLNS:      "http://mof.gov.vn/KKTNCN",
		MaSoThue:   taxCode,
		TenDonVi:   companyName,
		KyTinhThue: periodKey,
	}
	runMap := make(map[string]domain.PayrollRun)
	for _, r := range runs {
		runMap[r.EmployeeID] = r
	}
	for _, emp := range employees {
		run := runMap[emp.EmployeeID]
		taxable := run.GrossSalary - run.SIDeduction - run.HIDeduction - run.UIDeduction - run.TradeUnionDues
		// Personal deduction: 15.5M/month (Resolution 110/2025/UBTVQH15)
		giamTru := domain.PersonalDeductionMonthly
		nv := KKTNCNNhanVien{
			HoVaTen:          emp.EmployeeID,
			ThuNhapChiuThue: taxable,
			GiamTru:          giamTru,
			ThueTNCN:         run.PITAmount,
		}
		d.DanhSachNhanVien = append(d.DanhSachNhanVien, nv)
		d.TongCong.TongThuNhap += taxable
		d.TongCong.TongThueTNCN += run.PITAmount
	}
	return xml.MarshalIndent(d, "", "  ")
}

// ─── TK3-TS: Employer Registration Form ────────────────────────

type TK3TS struct {
	XMLName          xml.Name          `xml:"TK3TS"`
	XMLNS            string            `xml:"xmlns,attr"`
	MaDonVi          string            `xml:"MaDonVi"`
	TenDonVi         string            `xml:"TenDonVi"`
	MaSoThue         string            `xml:"MaSoThue"`
	DiaChi           string            `xml:"DiaChi"`
	DanhSachNhanVien []TK3TSNhanVien   `xml:"DanhSachNhanVien>NhanVien"`
}

type TK3TSNhanVien struct {
	MaSoBHXH   string  `xml:"MaSoBHXH"`
	HoVaTen    string  `xml:"HoVaTen"`
	NamSinh    string  `xml:"NamSinh"`
	GioiTinh   string  `xml:"GioiTinh"`
	SoCMND     string  `xml:"SoCMND"`
	DiaChi     string  `xml:"DiaChi"`
	NgayVaoLam string  `xml:"NgayVaoLam"`
	ChucVu     string  `xml:"ChucVu"`
	MucLuong   float64 `xml:"MucLuong"`
}

func GenerateTK3TS(companyName, taxCode, address string, employees []domain.EmployeePayrollInfo) ([]byte, error) {
	d := TK3TS{
		XMLNS:    "http://baohiemxahoi.gov.vn/TK3TS",
		MaDonVi:  taxCode,
		TenDonVi: companyName,
		MaSoThue: taxCode,
		DiaChi:   address,
	}
	for _, emp := range employees {
		nv := TK3TSNhanVien{
			HoVaTen:    emp.EmployeeID,
			NamSinh:    emp.ContractStartDate,
			NgayVaoLam: emp.ContractStartDate,
			MucLuong:   emp.BaseSalary,
		}
		d.DanhSachNhanVien = append(d.DanhSachNhanVien, nv)
	}
	return xml.MarshalIndent(d, "", "  ")
}
