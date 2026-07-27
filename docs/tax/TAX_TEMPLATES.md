# Tax Module — Declaration Forms & Templates

**Role:** BA Lead (20+ yrs) + Chief Accountant (20+ yrs)
**Date:** 2026-07-27
**Reference:** Circular 80/2021/TT-BTC, Circular 99/2025/TT-BTC

---

## 1. VAT Declaration Forms

### 01/GTGT — VAT Deduction Method (Monthly/Quarterly)

**Form structure per Circular 80/2021/TT-BTC Appendix II:**

| Indicator | Code | Field Name | Source |
|-----------|------|------------|--------|
| A | | **Thong tin chung** | |
| | [01] | Ky ke khai (Period) | Auto from period |
| | [02] | Lan dau / Bo sung | First/Amendment flag |
| B | | **Hoa don, chung tu mua vao** (Input) | |
| | [10] | Hang hoa, DV mua vao | Total purchases (Cr 331, 111, 112) |
| | [11] | Gia tri hang hoa nhap khau | Import purchases (156, 152) |
| | [12] | Thue GTGT hang nhap khau | Import VAT (33312) |
| | [13] | Gia tri hang hoa mua vao trong nuoc | Domestic purchases |
| | [14] | Thue GTGT dau vao | Input VAT (1331) |
| | [15] | Thue GTGT TSCD | Input VAT on FA (1332) |
| | [16] | Tong thue GTGT duoc khau tru | [14] + [15] |
| C | | **Hoa don, chung tu ban ra** (Output) | |
| | [20] | Doanh thu ban hang | Total sales (511) |
| | [21] | Thue GTGT dau ra | Output VAT (33311) |
| | [22] | Dieu chinh thue GTGT | Adjustment |
| | [23] | Tong thue GTGT dau ra | [21] + [22] |
| D | | **Xac dinh so thue phai nap** | |
| | [30] | Tong thue GTGT dau ra - dau vao | [23] - [16] |
| | [31] | Thue GTGT phai nap trong ky | If [30] > 0 = [30] |
| | [32] | Thue GTGT con duoc khau tru chuyen ky sau | If [30] < 0 = abs([30]) |

### 02/GTGT — VAT Direct Method

| Indicator | Code | Field Name |
|-----------|------|------------|
| | [01] | Gia tri hang hoa ban ra (Revenue + VAT incl.) |
| | [02] | Ty le GTGT tren doanh thu |
| | [03] | Gia tri gia tang tinh thue |
| | [04] | Thue GTGT phai nap = [03] × Rate |

### 03/GTGT — VAT for Gold/Silver/Precious Stones

| Indicator | Code | Field Name |
|-----------|------|------------|
| | [01] | Gia tri ban ra |
| | [02] | Gia tri mua vao tuong ung |
| | [03] | Gia tri gia tang = [01] - [02] |
| | [04] | Thue GTGT = [03] × 10% |

### Purchase Ledger (Bang ke mua vao)

| Column | Description |
|--------|-------------|
| STT | Serial number |
| Ten nguoi ban | Seller name |
| MST | Seller tax code |
| So hoa don | Invoice number |
| Ngay hoa don | Invoice date |
| Mat hang | Description |
| Gia tri truoc thue | Amount before VAT |
| Thue suat GTGT | VAT rate (0%, 5%, 8%, 10%) |
| Thue GTGT | VAT amount |
| Ghi chu | Notes |

### Sales Ledger (Bang ke ban ra)

| Column | Description |
|--------|-------------|
| STT | Serial number |
| Ten nguoi mua | Buyer name |
| MST | Buyer tax code |
| So hoa don | Invoice number |
| Ngay hoa don | Invoice date |
| Mat hang | Description |
| Gia tri truoc thue | Amount before VAT |
| Thue suat GTGT | VAT rate (0%, 5%, 8%, 10%) |
| Thue GTGT | VAT amount |
| Ghi chu | Notes |

---

## 2. CIT Declaration Forms

### 03/TNDN — CIT Annual Finalization

**Structure per Circular 80/2021/TT-BTC Appendix II (updated by Circular 40/2025/TT-BTC):**

| Indicator | Code | Field Name | Formula |
|-----------|------|------------|---------|
| I | | **Doanh thu, chi phi** | |
| | [01] | Doanh thu ban hang hoa, DV | 511 |
| | [02] | Doanh thu hoat dong tai chinh | 515 |
| | [03] | Thu nhap khac | 711 |
| | [04] | Tong doanh thu | [01]+[02]+[03] |
| | [05] | Cac khoan giam tru doanh thu | |
| | [06] | Doanh thu tinh thue | [04]-[05] |
| | [07] | Chi phi hop le duoc tru | 632+641+642+635+811 |
| | [08] | Thu nhap khac chiu thue | |
| | [09] | Loi nhap truoc thue | [06]-[07]+[08] |
| II | | **Xac dinh thue** | |
| | [10] | Thu nhap mien thue | |
| | [11] | Khoan lo duoc chuyen | Loss carry-forward |
| | [12] | Thu nhap tinh thue | [09]-[10]-[11] |
| | [13] | Thue suat (%) | 20% (or incentive rate) |
| | [14] | Thue TNDN phai nap | [12]×[13] |
| | [15] | Uu dai thue TNDN | Incentive reduction |
| | [16] | Thue TNDN con phai nap | [14]-[15] |
| III | | **Thue da nop** | |
| | [17] | Thue da tam nop | Sum of 04/TNDN filings |
| | [18] | Thue con phai nop | [16]-[17] (if >0) |
| | [19] | Thue nop thua | [17]-[16] (if >0) |

### Appendices to 03/TNDN

| Appendix | Code | Description |
|----------|------|-------------|
| 03-1A | KHCN | Science & technology fund |
| 03-2 | CHUYEN LO | Loss carry-forward detail |
| 03-3A | UUDAI DAU TU | Investment incentives |
| 03-3B | UUDAI MO RONG | Expansion project incentives |
| 03-3C | UUDAI CNCAO | High-tech zones |
| 03-3D | UUDAI KHCN | Science-tech enterprise |
| 03-4A | CHUYEN GIA | Transfer pricing disclosure |
| 03-5 | THOA THUAN | Tax treaty benefits |
| 03-6 | TRICH LAP QUY | Fund allocation |
| 03-7 | TSCD THUE | Assets subject to tax |

### 04/TNDN — Quarterly CIT Provisional

| Indicator | Code | Field Name |
|-----------|------|------------|
| | [01] | Doanh thu trong ky |
| | [02] | Chi phi trong ky |
| | [03] | Loi nhap = [01] - [02] |
| | [04] | Dieu chinh lam tang loi nhap |
| | [05] | Dieu chinh lam giam loi nhap |
| | [06] | Thu nhap tinh thue = [03] + [04] - [05] |
| | [07] | Thue suat (%) |
| | [08] | Thue TNDN phai nap = [06] × [07] |

---

## 3. PIT Declaration Forms

### 05/KK-TNCN — Monthly/Quarterly PIT Declaration

| Indicator | Code | Field Name |
|-----------|------|------------|
| I | | **Thu nhap chiu thue** |
| | [01] | Tong thu nhap chiu thue | Gross income |
| | [02] | Thu nhap mien thue | Exempt income |
| | [03] | Thu nhap tinh thue | [01] - [02] |
| II | | **Khoan giam tru** | |
| | [04] | Giam tru ban than | Personal deduction |
| | [05] | Giam tru nguoi phu thuoc | Dependant deductions |
| | [06] | Giam tru bao hiem | Insurance: BHXH+BHYT+BHTN |
| | [07] | Giam tru khac | Other deductions |
| | [08] | Tong giam tru | [04]+[05]+[06]+[07] |
| III | | **Thue** | |
| | [09] | Thu nhap tinh thue sau giam tru | [03] - [08] |
| | [10] | Thue TNCN phai nop | Apply progressive table |
| | [11] | So nguoi phu thuoc | Count of dependants |

### 02/QTT-TNCN — Annual PIT Finalization

Similar to 05/KK-TNCN but annual. Additional fields:
- [12] Thue da khau tru trong nam (PIT withheld during year)
- [13] Thue con phai nop / nop thua (Balance due/refund)
- [14] So thang giam tru (Months of personal deduction)
- [15] So thang giam tru nguoi phu thuoc (Months of dependant deduction)

### 05-1/BK-TNCN — Employee Detail Appendix

| Column | Description |
|--------|-------------|
| STT | Sequence |
| Ho va ten | Full name |
| MST | Tax code |
| Quoc tich | Nationality |
| Ngay sinh | DOB |
| Gioi tinh | Gender |
| CMND/HC | ID/Passport |
| Ngay cu tru | Residency start |
| Thu nhap chiu thue | Taxable income |
| Thu nhap mien thue | Exempt income |
| Giam tru ban than | Personal deduction |
| Giam tru nguoi phu thuoc | Dependant count |
| So ma so thue nguoi phu thuoc | Dependant tax codes |
| Giam tru bao hiem | Insurance deductions |
| Thue TNCN | PIT calculated |
| Thue da khau tru | PIT withheld |

---

## 4. Other Tax Forms

### 01/TTDB — Special Consumption Tax

| Indicator | Description |
|-----------|-------------|
| [01] | Doanh so ban hang chua co thue TTDB |
| [02] | Thue suat TTDB |
| [03] | Thue TTDB phai nap |
| [04] | Thue TTDB mua vao duoc khau tru |
| [05] | Thue TTDB con phai nap |

### 01/BVMT — Environmental Protection Tax

| Indicator | Description |
|-----------|-------------|
| Line | Pollutant type (xang, diesel, than, etc.) |
| Quantity | Volume in physical units |
| Rate per unit | VND/unit per Law on Environmental Tax |
| Tax payable | Quantity × Rate |

### FCT Forms (Foreign Contractor Tax)

**01/NTNN — FCT Declaration**

| Indicator | Description |
|-----------|-------------|
| [01] | Doanh thu chiu thue (Gross revenue) |
| [02] | VAT rate (%) |
| [03] | VAT payable = [01] × [02] |
| [04] | CIT rate (%) |
| [05] | CIT payable = [01] × [04] |
| [06] | Total tax = [03] + [05] |

FCT rates per service type:
| Service Type | CIT Rate | VAT Rate |
|-------------|----------|----------|
| Royalties | 10% | 5% |
| Services | 5% | 5% |
| Leasing | 5% | 5% |
| Interest | 5% | N/A |
| Insurance | 5% | 5% |
| Transport | 2% | 3% |
| Other | 2% | 5% |
| E-commerce (since 2025) | 2% | 8-10% |

---

## 5. E-Invoice TXML Format (Decree 254/2026/ND-CP)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<BK:HoaDon xmlns:BK="http://gdt.gov.vn/HoaDon">
  <BK:ThongTinChung>
    <BK:MaSoThue>0123456789</BK:MaSoThue>
    <BK:MauSo>01GTKT0/001</BK:MauSo>
    <BK:KyHieu>AA/22E</BK:KyHieu>
    <BK:NgayLap>2026-04-15</BK:NgayLap>
    <BK:ThoiDiemLap>2026-04-15T14:30:00+07:00</BK:ThoiDiemLap>
    <BK:LoaiHoaDon>BAN_HANG</BK:LoaiHoaDon>
    <BK:DonViTienTe>VND</BK:DonViTienTe>
  </BK:ThongTinChung>
  <BK:BenMua>
    <BK:Ten>CONG TY TNHH ABC</BK:Ten>
    <BK:MaSoThue>0987654321</BK:MaSoThue>
    <BK:DiaChi>123 Nguyen Hue, Q1, HCMC</BK:DiaChi>
    <BK:Email>info@abc.com</BK:Email>
  </BK:BenMua>
  <BK:DanhSachHangHoa>
    <BK:HangHoa>
      <BK:STT>1</BK:STT>
      <BK:TenHang>May tinh xach tay</BK:TenHang>
      <BK:DonViTinh>Chiec</BK:DonViTinh>
      <BK:SoLuong>5</BK:SoLuong>
      <BK:DonGia>20000000</BK:DonGia>
      <BK:ThanhTien>100000000</BK:ThanhTien>
      <BK:ThueSuat>10</BK:ThueSuat>
      <BK:ThueGTGT>10000000</BK:ThueGTGT>
    </BK:HangHoa>
  </BK:DanhSachHangHoa>
  <BK:TongTien>
    <BK:TongTienHang>100000000</BK:TongTienHang>
    <BK:TongThueGTGT>10000000</BK:TongThueGTGT>
    <BK:TongTienThanhToan>110000000</BK:TongTienThanhToan>
    <BK:SoTienBangChu>Mot tram muoi trieu dong</BK:SoTienBangChu>
  </BK:TongTien>
  <BK:KyThuat>
    <BK:ChuKySo>
      <BK:SerialNumber>AB1234567890</BK:SerialNumber>
      <BK:ThoiDiemKy>2026-04-15T14:30:05+07:00</BK:ThoiDiemKy>
      <BK:DuLieuKy>Base64signaturedata...</BK:DuLieuKy>
    </BK:ChuKySo>
    <BK:MaCuaCQT>A1B2C3D4-E5F6-7890</BK:MaCuaCQT>
  </BK:KyThuat>
</BK:HoaDon>
```

---

## 6. Tax Payment Order Format

```
KHO BAC NHÀ NƯỚC
LENH NOP TIEN VAO NGAN SACH NHA NUOC

□ Nop thue phat sinh
□ Nop thue theo thong bao/thong bao cua CQT

────────────────────────────────────
Ma NDKT (revenue code): XXXX
Ma Chuong (chapter): YYY
Ma DVBHNS: ZZZ

So tien bang so: 124,000,000
So tien bang chu: Mot tram hai bon trieu dong

Don vi nop:
  Ten: CONG TY TNHH ABC
  MST: 0123456789

Noi dung: Nop thue TNDN nam 2025

────────────────────────────────────
So tai khoan NSNN: 7111.1112.123456
Tai Kho BNN: KBNN TP HCM
```

Revenue codes (Ma NDKT) mapping:

| Tax Type | Code |
|----------|------|
| VAT | 1701 |
| CIT (domestic) | 1801 |
| PIT (resident) | 1901 |
| PIT (non-resident) | 1902 |
| TTDB | 1501 |
| BVMT | 1711 |
| FCT (CIT) | 1851 |
| FCT (VAT) | 1751 |
| Resource tax | 1601 |
| Land tax | 1401 |
| Import tax | 1301 |
| Export tax | 1201 |
| License tax | 1101 |
