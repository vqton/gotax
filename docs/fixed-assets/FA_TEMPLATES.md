# Fixed Asset Module — Document Templates

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Compliant with:** Circular 99/2025/TT-BTC, Circular 200/2014/TT-BTC, Decree 123/2020, IAS 16, VAS 03

---

## 1. FA Registration Form (Phiếu đăng ký TSCĐ)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       PHIẾU ĐĂNG KÝ TÀI SẢN CỐ ĐỊNH                          │
│                       FIXED ASSET REGISTRATION                              │
│                                                                             │
│  Số/No: FA-REG-202607-0001                       Ngày/Date: 15/07/2026       │
├─────────────────────────────────────────────────────────────────────────────┤
│  THÔNG TIN CHUNG / GENERAL INFORMATION                                      │
│                                                                             │
│  Mã TSCĐ / Asset Code:     FIT-2026-001                                     │
│  Tên TSCĐ / Asset Name:    Máy tính Dell Latitude 5420                      │
│  Loại TSCĐ / Category:     Máy móc thiết bị / Machinery & Equipment        │
│  Nhóm TSCĐ / Group:        MMTB-03 / Machinery-Equipment-Group-03          │
│  Ngày mua / Acq. Date:     10/07/2026                                       │
│  Ngày đưa vào sử dụng / In-service Date:  15/07/2026                        │
│                                                                             │
│  GIÁ TRỊ / VALUATION                                                        │
│                                                                             │
│  Nguyên giá / Original Cost:                  25,000,000                     │
│  Giá trị thanh lý / Residual Value:            1,250,000                     │
│  Giá trị tính khấu hao / Depreciable Amount:  23,750,000                    │
│                                                                             │
│  KHẤU HAO / DEPRECIATION                                                     │
│                                                                             │
│  Thời gian sử dụng / Useful Life:              48 tháng / months             │
│  Phương pháp khấu hao / Depreciation Method:   Đường thẳng / Straight-Line  │
│  Tỷ lệ khấu hao / Depreciation Rate:           25%/năm (year)               │
│  Mức KH tháng / Monthly Depreciation:            494,792                     │
│                                                                             │
│  THÔNG TIN KHÁC / OTHER INFORMATION                                          │
│                                                                             │
│  Phòng ban sử dụng / Department:          Phòng CNTT / IT Dept.             │
│  Địa điểm sử dụng / Location:             Lầu 3 - VP HCM / 3rd Floor HCMC   │
│  Nhà cung cấp / Supplier:                Công ty XYZ / XYZ Corp.           │
│  Số hợp đồng / Contract No:               HD-MT-2026-123                    │
│  Số hóa đơn / Invoice No:                INV-54321                          │
│  Nguồn hình thành / Source:              Mua sắm / Purchase                 │
│  Số serial / Serial No:                   DL5420-XYZ123456                  │
│  Nhà sản xuất / Manufacturer:            Dell Inc.                          │
│  Năm sản xuất / Manuf. Year:             2026                               │
│  Nước sản xuất / Country of Origin:      China                              │
│                                                                             │
│  TÀI KHOẢN HẠCH TOÁN / GL ACCOUNTS                                          │
│                                                                             │
│  TK TSCĐ / Asset Account:                 2112 (Máy móc thiết bị)           │
│  TK Khấu hao / Depreciation Account:      2142 (HMMT thiết bị)              │
│  TK Chi phí / Expense Account:            6424 (CP quản lý - khấu hao)      │
├─────────────────────────────────────────────────────────────────────────────┤
│ Người lập              Kế toán trưởng           Thủ trưởng đơn vị           │
│ (Prepared by)           (Chief Accountant)       (Director)                  │
│  ___________            ___________              ___________                 │
│  Ngày/Date: ___         Ngày/Date: ___           Ngày/Date: ___             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Sample JSON Payload (POST /api/v1/fixed-assets):**

```json
{
  "code": "FIT-2026-001",
  "name": "Máy tính Dell Latitude 5420",
  "category_id": "9b3f5c71-2a8e-4d1c-9f0a-5e7b3c2d1f0e",
  "acquisition_date": "2026-07-10T00:00:00Z",
  "original_cost": 25000000,
  "residual_value": 1250000,
  "useful_life_months": 48,
  "depreciation_method": "STRAIGHT_LINE",
  "department_id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
  "location": "Lầu 3 - VP HCM",
  "supplier_id": "d4e5f6a7-b8c9-0123-def4-567890abcdef",
  "contract_no": "HD-MT-2026-123",
  "invoice_id": "e7f8a9b0-c1d2-3456-7890-abcdef123456",
  "source": "PURCHASE",
  "serial_no": "DL5420-XYZ123456",
  "manufacturer": "Dell Inc.",
  "manufacture_year": 2026,
  "country_of_origin": "China",
  "technical_specs": "Core i7-1365U, 16GB RAM, 512GB SSD, 14-inch FHD",
  "asset_account_id": "2112",
  "depreciation_account_id": "2142",
  "expense_account_id": "6424"
}
```

---

## 2. FA Tag/Label (Thẻ TSCĐ)

Per Circular 200/2014/TT-BTC — Form S01-TSCĐ (Mẫu số 01-TSCĐ).

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          THẺ TÀI SẢN CỐ ĐỊNH                                │
│                          FIXED ASSET CARD                                   │
│                                                                             │
│  Mẫu số / Form: 01-TSCĐ (TT 200/2014)                                      │
│  Ban hành theo / Issued per: Circular 200/2014/TT-BTC                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Mã TSCĐ / Asset Code:           FIT-2026-001                               │
│  Tên TSCĐ / Asset Name:          Máy tính Dell Latitude 5420                │
│  Loại TSCĐ / Category:           Máy móc thiết bị                           │
│  Số hiệu TK TSCĐ / Asset GL:     2112                                       │
│  Số hiệu TK HM / Depr. GL:       2142                                       │
│  Số hiệu TK CP / Expense GL:     6424                                       │
│                                                                             │
├─┬───────────────────────────────────────────────────────────────────────────┤
│ │ 1. NGÀY THÁNG NGHIỆP THU / ACQUISITION DETAILS                           │
│ ├───────────────────────────────────────────────────────────────────────────┤
│ │ Ngày mua / Acquisition Date:      10/07/2026                              │
│ │ Ngày đưa vào sử dụng / In-Service: 15/07/2026                            │
│ │ Nguồn hình thành / Source:         Mua sắm / Purchase                     │
│ │                                                                           │
│ │ - Nguyên giá / Original Cost:            25,000,000                       │
│ │ - Giá trị thanh lý / Residual Value:      1,250,000                       │
│ │ - Giá trị tính KH / Depreciable Amount: 23,750,000                        │
│ │ - Chi phí vận chuyển / Transport Cost:        0                           │
│ │ - Chi phí lắp đặt / Installation Cost:        0                           │
│ │ - Thuế / Tax:                                     0                       │
│ │                                                                           │
│ │ Chứng từ từng lần tăng / Supporting Documents:                            │
│ │   HĐ mua số: INV-54321 ngày 10/07/2026                                    │
│ │   Hợp đồng số: HD-MT-2026-123                                             │
├─┼───────────────────────────────────────────────────────────────────────────┤
│ │ 2. THEO DÕI KHẤU HAO / DEPRECIATION TRACKING                              │
│ ├───────────────────────────────────────────────────────────────────────────┤
│ │ Thời gian sử dụng / Useful Life:          48 tháng                        │
│ │ Phương pháp tính KH / Method:             Đường thẳng / Straight-Line     │
│ │ Tỷ lệ KH / Depreciation Rate:             25%/năm                         │
│ │ Mức KH tháng / Monthly Depr.:               494,792                       │
│ │                                                                           │
│ │ Các lần điều chỉnh / Adjustments:                                          │
│ │ ┌──────┬──────────┬────────────┬─────────────┬───────────────────────────┐│
│ │ │ Lần  │ Ngày     │ Nguyên giá │ Số KH lũy   │ Ghi chú                  ││
│ │ │ (#)  │ (Date)   │ Cost       │ Acc. Depr.  │ (Notes)                  ││
│ │ ├──────┼──────────┼────────────┼─────────────┼───────────────────────────┤│
│ │ │  1   │15/07/2026│ 25,000,000 │        0    │ Nhập mới / New            ││
│ │ │      │          │            │             │                           ││
│ │ │      │          │            │             │                           ││
│ │ └──────┴──────────┴────────────┴─────────────┴───────────────────────────┘│
├─┴───────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ 3. GIẢM TSCĐ / DISPOSAL (khi bàn giao/thanh lý)                            │
│                                                                             │
│ ┌──────┬──────────┬────────────┬─────────────┬─────────────────────────────┐│
│ │ Lý do│ Chứng từ │ Ngày       │ Nguyên giá  │ Số KH lũy kế               ││
│ │(Reason)│ (Doc)  │ (Date)     │ Cost        │ Acc. Depreciation           ││
│ ├──────┼──────────┼────────────┼─────────────┼─────────────────────────────┤│
│ │      │          │            │             │                             ││
│ │      │          │            │             │                             ││
│ └──────┴──────────┴────────────┴─────────────┴─────────────────────────────┘│
├─────────────────────────────────────────────────────────────────────────────┤
│ PHỤ TRÁCH BỘ PHẬN        KẾ TOÁN TRƯỞNG           THỦ TRƯỞNG ĐƠN VỊ       │
│ (Department Head)         (Chief Accountant)         (Director)              │
│  ___________              ___________                ___________             │
│  Ngày/Date: ___           Ngày/Date: ___             Ngày/Date: ___         │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Print Format Specification:**

| Property | Specification |
|----------|--------------|
| Paper size | A5 (148mm × 210mm) or A4 folded |
| Orientation | Portrait |
| Font | Unicode (Arial/Time New Roman), size 10-11pt |
| Header | Bold, centered, 14pt |
| Border | 1pt solid line, 0.5pt inner grid |
| Vietnamese | Left, above English translation in parentheses |
| Amounts | Right-aligned, comma-separated thousands |
| Copies | 2 bản (02 copies): Lưu phòng TCKT + Lưu hồ sơ TSCĐ |
| Color | White paper, black ink |
| Binding | Pre-numbered, bound in register book |

---

## 3. Depreciation Schedule (Bảng tính khấu hao TSCĐ)

```
┌───────────────────────────────────────────────────────────────────────────────────────────────┐
│                     BẢNG TÍNH KHẤU HAO TÀI SẢN CỐ ĐỊNH                                        │
│                     DEPRECIATION SCHEDULE                                                      │
│                                                                                               │
│  Kỳ / Period: Tháng 07/2026                              Ngày in / Print Date: 31/07/2026     │
│  Đơn vị / Unit: VND                                                                           │
│  ĐVT / Currency: Đồng / VND                                                                   │
├───────────────────────────────────────────────────────────────────────────────────────────────┤
│ Phòng ban / Department: Công nghệ thông tin (CNTT) / IT                                       │
├────┬──────────┬────────────────────┬───────────┬──────┬──────────┬────────────┬────────┬───────┤
│ STT│ Mã TSCĐ  │ Tên TSCĐ           │ Nguyên giá│ Số   │ Phương   │ KH tháng   │ KH lũy │ GTCL  │
│    │ (Code)   │ (Asset Name)       │ (Cost)    │ năm  │ pháp KH  │ này        │ kế     │(Carr. │
│    │          │                    │           │(Life)│ (Method) │(Month Depr)│(Accum) │Amount)│
├────┼──────────┼────────────────────┼───────────┼──────┼──────────┼────────────┼────────┼───────┤
│  1 │ FIT-001  │ Dell Latitude 5420 │ 25,000,000│  4   │ SL       │    494,792 │      0 │ 25,000,000│
│  2 │ FIT-002  │ MacBook Pro 14"    │ 45,000,000│  5   │ SL       │    729,000 │ 29,160,000│15,840,000│
│  3 │ FIT-003  │ Server Dell R750   │380,000,000│  6   │ SL       │  6,333,333 │253,333,333│126,666,667│
├────┼──────────┼────────────────────┼───────────┼──────┼──────────┼────────────┼────────┼───────┤
│    │ Cộng / Total: IT Dept.       │450,000,000│      │          │  7,557,125 │282,493,333│167,506,667│
├────┴──────────┴────────────────────┴───────────┴──────┴──────────┴────────────┴────────┴───────┤
│ Phòng ban / Department: Kế toán / Accounting                                                   │
├────┬──────────┬────────────────────┬───────────┬──────┬──────────┬────────────┬────────┬───────┤
│  4 │ KT-001   │ Máy in Canon IR    │ 35,000,000│  5   │ DB 2x    │  1,166,667 │ 11,666,667│ 8,400,000│
│  5 │ KT-002   │ Điều hòa Mitsubishi│ 28,000,000│  6   │ SL       │    388,889 │ 15,555,556│12,444,444│
├────┼──────────┼────────────────────┼───────────┼──────┼──────────┼────────────┼────────┼───────┤
│    │ Cộng / Total: Accounting Dept.│ 63,000,000│      │          │  1,555,556 │ 27,222,223│20,844,444│
├────┴──────────┴────────────────────┴───────────┴──────┴──────────┴────────────┴────────┴───────┤
│                                                                                                │
│ TỔNG CỘNG TOÀN CÔNG TY / GRAND TOTAL                                  │                     │
├────┬──────────────────────────────────────────┬─────────────────────────┼──────────────────────┤
│    │ Nguyên giá / Total Original Cost         │              513,000,000 │                      │
│    │ Khấu hao trong kỳ / Depr. This Period    │                9,112,681 │                      │
│    │ Khấu hao lũy kế / Accumulated Depr.      │              309,715,556 │                      │
│    │ Giá trị còn lại / Carrying Amount        │              203,284,444 │                      │
└────┴──────────────────────────────────────────┴─────────────────────────┴──────────────────────┘
```

**Sample JSON Response (GET /api/v1/reports/fa-depreciation-schedule):**

```json
{
  "period": {"year": 2026, "month": 7},
  "rows": [
    {
      "seq": 1,
      "fixed_asset_id": "a1b2...",
      "code": "FIT-001",
      "name": "Dell Latitude 5420",
      "department": "IT",
      "category": "Máy móc thiết bị",
      "original_cost": 25000000,
      "useful_life_months": 48,
      "depreciation_method": "STRAIGHT_LINE",
      "monthly_depreciation": 494792,
      "accumulated_depreciation": 0,
      "carrying_amount": 25000000,
      "depreciation_start_date": "2026-07-15",
      "depreciation_end_date": "2030-07-15"
    }
  ],
  "summary": {
    "total_original_cost": 513000000,
    "total_monthly_depreciation": 9112681,
    "total_accumulated_depreciation": 309715556,
    "total_carrying_amount": 203284444
  }
}
```

**Excel Export Column Mapping:**

| Column | Header (VN) | Header (EN) | Type | Format |
|--------|-------------|-------------|------|--------|
| A | STT | No. | Integer | 1, 2, 3... |
| B | Mã TSCĐ | Asset Code | Text | |
| C | Tên TSCĐ | Asset Name | Text | |
| D | Phòng ban | Department | Text | |
| E | Loại TSCĐ | Category | Text | |
| F | Nguyên giá | Original Cost | Number | #,##0 |
| G | Thời gian SD (tháng) | Life (months) | Integer | |
| H | Phương pháp KH | Depr. Method | Text | SL / DB / PB |
| I | KH tháng này | Month Depr. | Number | #,##0 |
| J | KH lũy kế | Accumulated | Number | #,##0 |
| K | Giá trị còn lại | Carrying Amount | Number | #,##0 |
| L | Ngày bắt đầu KH | Depr. Start | Date | DD/MM/YYYY |
| M | Ngày kết thúc KH | Depr. End | Date | DD/MM/YYYY |

---

## 4. FA Movement Schedule (Báo cáo tăng giảm TSCĐ)

Per VAS 03 (Vietnamese Accounting Standard No. 03 — Tangible/Intangible Fixed Assets) disclosure format.

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                     BÁO CÁO TĂNG, GIẢM TÀI SẢN CỐ ĐỊNH                                         │
│                     FIXED ASSET MOVEMENT SCHEDULE                                               │
│                                                                                                  │
│  Kỳ / Period: Năm 2026 / Year 2026                                Đơn vị / Unit: VND            │
├──────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                  │
│  I. TÀI SẢN CỐ ĐỊNH HỮU HÌNH / TANGIBLE FIXED ASSETS (TK 211)                                  │
├────────┬────────────┬────────────┬────────────┬────────────┬────────────┬────────────┬──────────┤
│ Chỉ    │ Nhà cửa,   │ Máy móc,  │ Phương tiện│ Thiết bị,  │ Cây xanh,  │ TSCĐ hữu   │ TỔNG    │
│ tiêu   │ vật kiến   │ thiết bị   │ vận chuyển │ dụng cụ    │ súc vật,   │ hình khác  │ CỘNG    │
│ (Item) │ trúc       │ (2112)     │ (2113)     │ quản lý    │ (2115)     │ (2118)     │(Total)  │
│        │ (2111)     │            │            │ (2114)     │            │            │         │
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│        │         A. NGUYÊN GIÁ / ORIGINAL COST                                                  │
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ SDĐK   │5,000,000,000│2,500,000,000│800,000,000│600,000,000│  50,000,000│ 100,000,000│9,050,000,000│
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ Tăng   │         0  │ 450,000,000│         0  │ 120,000,000│         0  │         0  │ 570,000,000│
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│  - Mua │         0  │ 380,000,000│         0  │ 120,000,000│         0  │         0  │ 500,000,000│
│  - XDCB│         0  │         0  │         0  │         0  │         0  │         0  │         0  │
│  - Cấp │         0  │         0  │         0  │         0  │         0  │         0  │         0  │
│  - Khác│         0  │  70,000,000│         0  │         0  │         0  │         0  │  70,000,000│
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ Giảm   │         0  │         0  │         0  │         0  │         0  │         0  │         0  │
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│  - Thanh lý│      0  │         0  │         0  │         0  │         0  │         0  │         0  │
│  - Bán  │         0  │         0  │         0  │         0  │         0  │         0  │         0  │
│  - Khác│         0  │         0  │         0  │         0  │         0  │         0  │         0  │
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ SDCK   │5,000,000,000│2,950,000,000│800,000,000│720,000,000│  50,000,000│ 100,000,000│9,620,000,000│
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│        │      B. GIÁ TRỊ HAO MÒN / ACCUMULATED DEPRECIATION (TK 214)                           │
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ SDĐK   │1,500,000,000│1,200,000,000│400,000,000│250,000,000│  20,000,000│  60,000,000│3,430,000,000│
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ KH tăng│  83,333,333│ 100,000,000│ 13,333,333│  20,000,000│   1,666,667│   3,333,333│ 221,666,667│
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ KH giảm│         0  │         0  │         0  │         0  │         0  │         0  │         0  │
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ SDCK   │1,583,333,333│1,300,000,000│413,333,333│270,000,000│  21,666,667│  63,333,333│3,651,666,667│
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│        │   C. GIÁ TRỊ CÒN LẠI / CARRYING AMOUNT                                                │
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ GTCL ĐK│3,500,000,000│1,300,000,000│400,000,000│350,000,000│  30,000,000│  40,000,000│5,620,000,000│
├────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────┤
│ GTCL CK│3,416,666,667│1,650,000,000│386,666,667│450,000,000│  28,333,333│  36,666,667│5,968,333,333│
└────────┴────────────┴────────────┴────────────┴────────────┴────────────┴────────────┴──────────┘

┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│  II. TÀI SẢN CỐ ĐỊNH VÔ HÌNH / INTANGIBLE FIXED ASSETS (TK 213)                                │
├────────────────────────┬────────────┬────────────┬────────────┬────────────┬────────────────────┤
│ Chỉ tiêu (Item)        │ Quyền SD   │ Bản quyền, │ Phần mềm   │ Nhãn hiệu, │ TỔNG CỘNG         │
│                        │ đất (2131) │ bằng sáng  │ máy tính   │ thương     │ (Total)            │
│                        │            │ chế (2132) │ (2133)     │ hiệu (2134)│                    │
├────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────────────┤
│ A. NGUYÊN GIÁ          │            │            │            │            │                    │
│ Số dư đầu kỳ           │200,000,000 │   500,000  │1,200,000,000│   300,000  │1,400,800,000       │
│ Tăng trong kỳ          │         0  │         0  │    50,000,000│         0  │   50,000,000       │
│ Giảm trong kỳ          │         0  │         0  │         0   │         0  │         0           │
│ Số dư cuối kỳ          │200,000,000 │   500,000  │1,250,000,000│   300,000  │1,450,800,000       │
├────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────────────┤
│ B. HAO MÒN             │            │            │            │            │                    │
│ Số dư đầu kỳ           │  50,000,000│   500,000  │   500,000,000│   300,000 │  550,800,000       │
│ KH trong kỳ            │   2,000,000│         0  │    20,833,333│         0  │   22,833,333       │
│ Giảm                   │         0  │         0  │         0   │         0  │         0           │
│ Số dư cuối kỳ          │  52,000,000│   500,000  │   520,833,333│   300,000 │  573,633,333       │
├────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────────────┤
│ C. GIÁ TRỊ CÒN LẠI     │            │            │            │            │                    │
│ GTCL đầu kỳ            │150,000,000 │         0  │   700,000,000│         0  │  850,000,000       │
│ GTCL cuối kỳ           │148,000,000 │         0  │   729,166,667│         0  │  877,166,667       │
└────────────────────────┴────────────┴────────────┴────────────┴────────────┴────────────────────┘
```

**Sample JSON Response (GET /api/v1/reports/fa-movement):**

```json
{
  "period": "2026",
  "tangible": {
    "by_type": {
      "2111": {
        "name": "Nhà cửa, vật kiến trúc",
        "opening_cost": 5000000000,
        "increases": {
          "purchase": 0,
          "construction": 0,
          "donation": 0,
          "other": 0,
          "total": 0
        },
        "decreases": {
          "disposal": 0,
          "sale": 0,
          "other": 0,
          "total": 0
        },
        "closing_cost": 5000000000,
        "opening_depreciation": 1500000000,
        "depreciation_increase": 83333333,
        "depreciation_decrease": 0,
        "closing_depreciation": 1583333333,
        "opening_carrying": 3500000000,
        "closing_carrying": 3416666667
      }
    },
    "summary": {
      "opening_cost": 9050000000,
      "increase_total": 570000000,
      "decrease_total": 0,
      "closing_cost": 9620000000,
      "opening_depreciation": 3430000000,
      "closing_depreciation": 3651666667,
      "opening_carrying": 5620000000,
      "closing_carrying": 5968333333
    }
  },
  "intangible": {
    "by_type": {
      "2131": {
        "name": "Quyền sử dụng đất",
        "opening_cost": 200000000,
        "increases": {"total": 0},
        "decreases": {"total": 0},
        "closing_cost": 200000000,
        "opening_depreciation": 50000000,
        "closing_depreciation": 52000000,
        "opening_carrying": 150000000,
        "closing_carrying": 148000000
      }
    }
  }
}
```

---

## 5. FA Register (Sổ TSCĐ)

Per Circular 200/2014/TT-BTC, form S02-TSCĐ — individual asset tracking from acquisition to disposal.

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                           SỔ TÀI SẢN CỐ ĐỊNH / FIXED ASSET REGISTER                             │
│                                                                                                  │
│  Mã TSCĐ / Asset Code: FIT-2026-001                  Loại / Category: Máy móc thiết bị           │
│  Tên TSCĐ / Asset Name: Máy tính Dell Latitude 5420  TK TSCĐ / Asset GL: 2112                    │
│  Nguyên giá / Cost: 25,000,000                        TK HM / Depr. GL: 2142                     │
│  Thời gian SD / Life: 48 tháng                        TK CP / Expense GL: 6424                   │
├───────┬────────────┬──────────────┬───────────┬────────────┬────────────┬────────────┬──────────┤
│ Ngày  │ Số CT      │ Diễn giải    │ Loại phát │ Nguyên giá │ Nguyên giá │ KH lũy kế  │ GTCL     │
│(Date) │ (Doc#)     │ (Desc.)      │ sinh      │ Tăng       │ Giảm       │ (Acc. Depr)│ (Carr.)  │
│       │            │              │(Trx Type) │(Cost Incr.)│(Cost Decr.)│            │          │
├───────┼────────────┼──────────────┼───────────┼────────────┼────────────┼────────────┼──────────┤
│01/07  │            │ Dư đầu kỳ    │ OPENING   │ 25,000,000 │          0 │          0 │ 25,000,000│
│       │            │ (Opening)    │           │            │            │            │          │
├───────┼────────────┼──────────────┼───────────┼────────────┼────────────┼────────────┼──────────┤
│15/07  │ FA-REG-001│ Đăng ký TSCĐ │ ACQUISITION│25,000,000 │          0 │          0 │ 25,000,000│
│       │            │ (Registration)│           │            │            │            │          │
├───────┼────────────┼──────────────┼───────────┼────────────┼────────────┼────────────┼──────────┤
│31/07  │ DEP-072026│ KH tháng 07   │ DEPRECIATION│        0 │          0 │    494,792 │ 24,505,208│
│       │            │ (Depr. Jul)  │           │            │            │            │          │
├───────┼────────────┼──────────────┼───────────┼────────────┼────────────┼────────────┼──────────┤
│31/08  │ DEP-082026│ KH tháng 08   │ DEPRECIATION│        0 │          0 │    989,584 │ 24,010,416│
│       │            │ (Depr. Aug)  │           │            │            │            │          │
├───────┼────────────┼──────────────┼───────────┼────────────┼────────────┼────────────┼──────────┤
│  ...  │     ...    │     ...      │    ...    │    ...     │    ...     │    ...     │   ...    │
├───────┼────────────┼──────────────┼───────────┼────────────┼────────────┼────────────┼──────────┤
│31/12  │ DEP-122030│ KH tháng cuối │ DEPRECIATION│        0 │          0 │ 23,750,000 │  1,250,000│
│       │            │ (Final Depr.)│           │            │            │            │          │
├───────┼────────────┼──────────────┼───────────┼────────────┼────────────┼────────────┼──────────┤
│01/01  │ LIQ-0001   │ Thanh lý     │ DISPOSAL   │        0  │ 25,000,000│ 23,750,000 │         0 │
│2031   │            │ (Disposal)   │           │            │            │            │          │
├───────┼────────────┼──────────────┼───────────┼────────────┼────────────┼────────────┼──────────┤
│       │            │ Dư cuối kỳ   │ CLOSING    │        0  │          0 │          0 │         0 │
│       │            │ (Closing)    │           │            │            │            │          │
└───────┴────────────┴──────────────┴───────────┴────────────┴────────────┴────────────┴──────────┘
```

**Sample JSON Response (GET /api/v1/reports/fa-register?asset_id=a1b2...):**

```json
{
  "asset": {
    "id": "a1b2c3d4-...",
    "code": "FIT-2026-001",
    "name": "Máy tính Dell Latitude 5420",
    "category": "Máy móc thiết bị",
    "original_cost": 25000000,
    "useful_life_months": 48,
    "depreciation_method": "STRAIGHT_LINE",
    "asset_account_id": "2112",
    "depreciation_account_id": "2142",
    "expense_account_id": "6424"
  },
  "entries": [
    {
      "date": "2026-07-15",
      "doc_no": "FA-REG-001",
      "description": "Đăng ký TSCĐ",
      "transaction_type": "ACQUISITION",
      "cost_increase": 25000000,
      "cost_decrease": 0,
      "accumulated_depreciation": 0,
      "carrying_amount": 25000000
    },
    {
      "date": "2026-07-31",
      "doc_no": "DEP-072026",
      "description": "KH tháng 07",
      "transaction_type": "DEPRECIATION",
      "cost_increase": 0,
      "cost_decrease": 0,
      "accumulated_depreciation": 494792,
      "carrying_amount": 24505208
    }
  ]
}
```

---

## 6. FA Disposal Certificate (Biên bản thanh lý TSCĐ)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                   BIÊN BẢN THANH LÝ TÀI SẢN CỐ ĐỊNH                         │
│                   FIXED ASSET DISPOSAL CERTIFICATE                         │
│                                                                             │
│  Số/No: FA-LIQ-202607-0001                    Ngày/Date: 31/07/2026         │
├─────────────────────────────────────────────────────────────────────────────┤
│  THÔNG TIN TSCĐ / ASSET INFORMATION                                         │
│                                                                             │
│  Mã TSCĐ / Asset Code:      FIT-2019-002                                    │
│  Tên TSCĐ / Asset Name:     Máy chủ HP ProLiant DL380 Gen10                 │
│  Loại TSCĐ / Category:      Máy móc, thiết bị                               │
│  Số Serial / Serial No:     HP-DL380-SN987654                               │
│  Năm đưa vào SD / In-Service: 2019                                          │
│                                                                             │
│  GIÁ TRỊ / VALUATION (tại ngày thanh lý / at disposal date)                 │
│                                                                             │
│  Nguyên giá / Original Cost:                        280,000,000             │
│  Khấu hao lũy kế / Accumulated Depreciation:        280,000,000             │
│  Giá trị còn lại / Carrying Amount:                         0               │
│  Giá trị thu hồi / Proceeds:                           5,000,000            │
│  Chi phí thanh lý / Disposal Cost:                     2,000,000            │
│  Lãi/(Lỗ) thanh lý / Gain/(Loss):                      3,000,000            │
│                                                                             │
│  LÝ DO THANH LÝ / DISPOSAL REASON: Hết thời gian sử dụng / Fully depreciated│
│                                                                             │
│  HẠCH TOÁN / GL POSTING                                                     │
│                                                                             │
│  Nợ 2142 (HM thiết bị)       280,000,000    (giảm KH lũy kế)               │
│  Nợ 811 (Chi phí khác)               0      (lỗ thanh lý nếu có)            │
│  Có 2112 (Máy móc TB)         280,000,000    (giảm nguyên giá)              │
│                                                                             │
│  Nợ 111/112 (Tiền)             5,000,000    (thu hồi từ thanh lý)           │
│  Có 711 (Thu nhập khác)        5,000,000    (thu nhập thanh lý)             │
│                                                                             │
│  Nợ 811 (Chi phí khác)        2,000,000     (chi phí thanh lý)              │
│  Có 111/112 (Tiền)            2,000,000     (chi phí thanh lý)              │
├─────────────────────────────────────────────────────────────────────────────┤
│  THÀNH PHẦN THAM DỰ / PARTICIPANTS                                          │
│                                                                             │
│  Ông/Bà / Mr/Ms: Nguyễn Văn A      Chức danh: Giám đốc / Director          │
│  Ông/Bà / Mr/Ms: Trần Thị B        Chức danh: Kế toán trưởng / Chief Acc.  │
│  Ông/Bà / Mr/Ms: Lê Văn C          Chức danh: Phụ trách CNTT / IT Manager  │
├─────────────────────────────────────────────────────────────────────────────┤
│  ĐẠI DIỆN HỘI ĐỒNG      KẾ TOÁN TRƯỞNG        NGƯỜI LẬP BIÊN BẢN          │
│  (Board Representative)  (Chief Accountant)     (Preparer)                   │
│  ___________             ___________            ___________                  │
│  Ngày/Date: ___          Ngày/Date: ___         Ngày/Date: ___              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Sample JSON Payload (PATCH /api/v1/fixed-assets/:id/dispose):**

```json
{
  "disposal_date": "2026-07-31",
  "reason": "FULLY_DEPRECIATED",
  "proceeds": 5000000,
  "disposal_cost": 2000000,
  "description": "Thanh lý máy chủ HP DL380 hết thời gian sử dụng",
  "gl_posting": {
    "debit_2142": 280000000,
    "credit_2112": 280000000,
    "debit_111": 5000000,
    "credit_711": 5000000,
    "debit_811": 2000000,
    "credit_111": 2000000
  }
}
```

---

## 7. FA Transfer Certificate (Biên bản bàn giao TSCĐ)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                   BIÊN BẢN BÀN GIAO TÀI SẢN CỐ ĐỊNH                        │
│                   FIXED ASSET TRANSFER CERTIFICATE                         │
│                                                                             │
│  Số/No: FA-TRF-202607-0001                    Ngày/Date: 31/07/2026         │
├─────────────────────────────────────────────────────────────────────────────┤
│  THÔNG TIN TSCĐ / ASSET INFORMATION                                         │
│                                                                             │
│  Mã TSCĐ / Asset Code:      FIT-2026-001                                    │
│  Tên TSCĐ / Asset Name:     Máy tính Dell Latitude 5420                     │
│  Loại TSCĐ / Category:      Máy móc, thiết bị                               │
│  Số Serial / Serial No:     DL5420-XYZ123456                                │
│  Nguyên giá / Original Cost:                    25,000,000                   │
│  Giá trị còn lại / Carrying Amount:            25,000,000                   │
│                                                                             │
│  BÀN GIAO / TRANSFER                                                        │
│                                                                             │
│  Bộ phận bàn giao / Source Department:  Phòng CNTT / IT Dept.              │
│  Bộ phận nhận / Target Department:      Phòng Kế toán / Accounting Dept.   │
│  Ngày hiệu lực / Effective Date:        01/08/2026                          │
│  Lý do chuyển / Transfer Reason:        Thay đổi cơ cấu phòng ban           │
│                                         / Department restructuring          │
│                                                                             │
│  HẠCH TOÁN / GL POSTING                                                     │
│                                                                             │
│  Ghi chú: Giao dịch này KHÔNG ảnh hưởng GL / This transfer has NO GL impact │
│  Chỉ thay đổi thông tin quản lý nội bộ / Internal allocation change only.  │
│  Chi phí khấu hao sẽ hạch toán vào TK 6424 (mới) từ kỳ tiếp theo           │
│  Depreciation expense will be posted to new expense account next period.    │
│                                                                             │
│  TÌNH TRẠNG TSCĐ / ASSET CONDITION AT TRANSFER                              │
│  ☐ Mới 100% / New           ☒ Tốt / Good                                   │
│  ☐ Đã qua sử dụng ổn / Fair  ☐ Hư hỏng / Damaged                           │
├─────────────────────────────────────────────────────────────────────────────┤
│  BÊN BÀN GIAO / TRANSFEROR                    BÊN NHẬN / TRANSFEREE         │
│                                                                             │
│  Đại diện / Representative:                    Đại diện / Representative:   │
│  ___________                                   ___________                  │
│  Chức vụ / Title:                              Chức vụ / Title:             │
│  ___________                                   ___________                  │
│                                                                             │
│  KẾ TOÁN TRƯỞNG                                THỦ TRƯỞNG ĐƠN VỊ           │
│  (Chief Accountant)                            (Director)                    │
│  ___________                                   ___________                  │
│  Ngày/Date: ___                                Ngày/Date: ___              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Sample JSON Payload (POST /api/v1/fixed-assets/:id/transfer):**

```json
{
  "source_department_id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
  "target_department_id": "f0e1d2c3-b4a5-6789-0123-456789abcdef",
  "effective_date": "2026-08-01",
  "reason": "Department restructuring",
  "new_expense_account_id": "6424",
  "notes": "Transfer from IT to Accounting per reorg decision #REORG-2026-07"
}
```

---

## 8. FA Inventory Tag (Thẻ kiểm kê TSCĐ)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                   THẺ KIỂM KÊ TÀI SẢN CỐ ĐỊNH                              │
│                   FIXED ASSET INVENTORY TAG                                │
│                                                                             │
│  Số/No: FA-INV-202607-0001                      Ngày/Date: 31/07/2026      │
│                                                                             │
│  Kế hoạch kiểm kê / Inventory Plan: INV-PLAN-2026-Q3                      │
│  Địa điểm kiểm kê / Location: Toàn bộ công ty / All company               │
│  Người kiểm kê / Inventory Team: Nguyễn Văn A, Trần Thị B                  │
├──────┬──────────┬────────────────────┬────────────────────┬──────┬─────────┤
│ Mã   │ Tên TSCĐ │ Vị trí dự kiến     │ Vị trí thực tế    │ KQ   │ Ghi chú │
│ TSCĐ │ (Name)   │ (Expected Loc.)   │ (Actual Loc.)     │(Disc)│ (Notes) │
├──────┼──────────┼────────────────────┼────────────────────┼──────┼─────────┤
│FIT-  │Máy tính  │Lầu 3 - Phòng 305  │Lầu 3 - Phòng 305  │ OK   │         │
│001  │Dell 5420 │IT Dept.           │IT Dept.           │      │         │
├──────┼──────────┼────────────────────┼────────────────────┼──────┼─────────┤
│FIT-  │MacBook   │Lầu 3 - Phòng 302  │Lầu 2 - Phòng 202  │ MOVED│ Cần điều│
│002  │Pro 14"   │Marketing          │Sales Dept.        │      │ chỉnh   │
├──────┼──────────┼────────────────────┼────────────────────┼──────┼─────────┤
│FIT-  │Server    │Phòng Server - Lầu1│Phòng Server - Lầu1│ OK   │    OK   │
│003  │Dell R750 │                   │                   │      │         │
├──────┼──────────┼────────────────────┼────────────────────┼──────┼─────────┤
│KT-   │Máy in    │Lầu 2 - Phòng 201  │Không tìm thấy     │MISSING│Báo mất  │
│001  │Canon IR  │                   │                   │      │ ngày 05/07│
├──────┼──────────┼────────────────────┼────────────────────┼──────┼─────────┤
│      │          │                    │                    │      │         │
│      │          │                    │                    │      │         │
├──────┴──────────┴────────────────────┴────────────────────┴──────┴─────────┤
│                                                                             │
│ TỔNG HỢP / SUMMARY                                                          │
│   Tổng số TSCĐ kiểm kê / Total Inventoried:      48                         │
│   Khớp / Matched (OK):                          42                         │
│   Sai vị trí / Moved (MOVED):                    4                         │
│   Thiếu / Missing (MISSING):                     1                         │
│   Hỏng / Damaged (DAMAGED):                      1                         │
│   Chưa đăng ký / Unregistered:                   0                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ NGƯỜI KIỂM KÊ          KẾ TOÁN TRƯỞNG          THỦ TRƯỞNG ĐƠN VỊ          │
│ (Inventory Taker)       (Chief Accountant)       (Director)                  │
│  ___________            ___________              ___________                 │
│  Ngày/Date: ___         Ngày/Date: ___           Ngày/Date: ___             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Discrepancy Types:**

| Code | Vietnamese | English | Action |
|------|-----------|---------|--------|
| OK | Khớp | Matched | None |
| MOVED | Sai vị trí | Location mismatch | Update location in system |
| MISSING | Thiếu | Missing | Investigation, write-off procedure |
| DAMAGED | Hỏng | Damaged | Repair or disposal assessment |
| UNREGISTERED | Chưa đăng ký | Unregistered | Create FA registration |

**Sample JSON for Recording Result (POST /api/v1/fixed-assets/inventory/results):**

```json
{
  "plan_id": "plan-uuid-here",
  "results": [
    {
      "fixed_asset_id": "asset-uuid-001",
      "expected_location": "Lầu 3 - Phòng 305, IT Dept.",
      "actual_location": "Lầu 3 - Phòng 305, IT Dept.",
      "expected_status": "ACTIVE",
      "actual_status": "ACTIVE",
      "discrepancy": "OK",
      "notes": ""
    },
    {
      "fixed_asset_id": "asset-uuid-002",
      "expected_location": "Lầu 3 - Phòng 302, Marketing",
      "actual_location": "Lầu 2 - Phòng 202, Sales",
      "expected_status": "ACTIVE",
      "actual_status": "ACTIVE",
      "discrepancy": "MOVED",
      "notes": "Cần cập nhật vị trí trong hệ thống"
    },
    {
      "fixed_asset_id": "asset-uuid-004",
      "expected_location": "Lầu 2 - Phòng 201",
      "actual_location": "",
      "expected_status": "ACTIVE",
      "actual_status": "MISSING",
      "discrepancy": "MISSING",
      "notes": "Báo mất ngày 05/07/2026 - chờ xử lý"
    }
  ]
}
```

---

## 9. E-Invoice XML Template for FA Sale

Per GDT e-invoice spec (Decree 123/2020, Decree 254/2026). Special fields for FA disposition sale.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Invoice xmlns="http://gdt.gov.vn/schemas/einvoice/2026">
  <InvoiceHeader>
    <InvoiceNumber>INV-FA-202607-0001</InvoiceNumber>
    <InvoiceDate>2026-07-31</InvoiceDate>
    <InvoiceTime>10:30:00</InvoiceTime>
    <InvoiceType>01GTKT</InvoiceType>
    <InvoiceSeries>AA/26E</InvoiceSeries>
    <CurrencyCode>VND</CurrencyCode>
    <ExchangeRate>1</ExchangeRate>
    <InvoiceCategory>FIXED_ASSET_SALE</InvoiceCategory>
    <!-- Hóa đơn bán TSCĐ / FA Sale Invoice -->
  </InvoiceHeader>
  <Seller>
    <TaxCode>0123456789</TaxCode>
    <Name>Công ty ABC</Name>
    <Address>123 Đường ABC, Quận 1, TP.HCM</Address>
    <Phone>0283822111</Phone>
    <BankAccount>0123456789</BankAccount>
    <BankName>Vietcombank</BankName>
  </Seller>
  <Buyer>
    <TaxCode>9876543210</TaxCode>
    <Name>Công ty XYZ</Name>
    <Address>456 Đường DEF, Quận 1, TP.HCM</Address>
    <Phone>0283822333</Phone>
  </Buyer>
  <InvoiceLines>
    <Line>
      <LineNumber>1</LineNumber>
      <ItemType>FIXED_ASSET</ItemType>
      <AssetCode>FIT-2019-002</AssetCode>
      <!-- Mã TSCĐ / Fixed Asset Code -->
      <SerialNo>HP-DL380-SN987654</SerialNo>
      <!-- Số serial TSCĐ / Asset Serial Number -->
      <ItemCode>FA-HPSERVER</ItemCode>
      <ItemName>Máy chủ HP ProLiant DL380 Gen10 đã qua sử dụng</ItemName>
      <!-- Tên hàng ghi rõ "đã qua sử dụng" / Name must state "used" -->
      <Unit>Cái</Unit>
      <Quantity>1</Quantity>
      <UnitPrice>5000000</UnitPrice>
      <VatRate>10</VatRate>
      <!-- Thuế suất VAT 10% cho bán TSCĐ (nếu không thuộc đối tượng chịu thuế thì 0) -->
      <VatAmount>500000</VatAmount>
      <LineTotal>5000000</LineTotal>
      <!-- FA-specific metadata -->
      <FixedAssetInfo>
        <OriginalCost>280000000</OriginalCost>
        <!-- Nguyên giá / Original Cost -->
        <AccumulatedDepreciation>280000000</AccumulatedDepreciation>
        <!-- Khấu hao lũy kế / Accumulated Depreciation -->
        <CarryingAmount>0</CarryingAmount>
        <!-- Giá trị còn lại / Carrying Amount -->
        <AcquisitionDate>2019-05-15</AcquisitionDate>
        <!-- Ngày mua / Acquisition Date -->
        <DisposalReason>FULLY_DEPRECIATED</DisposalReason>
        <!-- Lý do thanh lý / Disposal Reason -->
        <DisposalDecisionNo>QD-LIQ-2026-001</DisposalDecisionNo>
        <!-- Số quyết định thanh lý / Disposal Decision Number -->
      </FixedAssetInfo>
      <TaxExemptCategory xsi:nil="true"></TaxExemptCategory>
    </Line>
  </InvoiceLines>
  <Summary>
    <Subtotal>5000000</Subtotal>
    <Discount>0</Discount>
    <VatAmount>500000</VatAmount>
    <GrandTotal>5500000</GrandTotal>
    <AmountInWords>Năm triệu năm trăm nghìn đồng</AmountInWords>
  </Summary>
  <DigitalSignature>
    <Signature>base64-encoded-signature</Signature>
    <SignDate>2026-07-31</SignDate>
    <CertificateSerial>SN-123456</CertificateSerial>
  </DigitalSignature>
  <QRCode>base64-qr-code-image</QRCode>
</Invoice>
```

**Key FA-Specific XML Elements:**

| Element | Required | Description |
|---------|----------|-------------|
| `InvoiceCategory` | Yes | Must be `FIXED_ASSET_SALE` to distinguish from regular sale |
| `ItemType` | Yes | Must be `FIXED_ASSET` |
| `AssetCode` | Yes | Internal FA code for cross-reference |
| `SerialNo` | Yes | Physical serial number of the asset |
| `ItemName` | Yes | Must include "đã qua sử dụng" (used) to indicate second-hand |
| `FixedAssetInfo.OriginalCost` | Yes | Original cost at acquisition (informational) |
| `FixedAssetInfo.AccumulatedDepreciation` | Yes | Total depreciation accumulated (informational) |
| `FixedAssetInfo.CarryingAmount` | Yes | Book value at sale date (informational) |
| `FixedAssetInfo.DisposalReason` | Yes | Reason for disposal/sale |
| `FixedAssetInfo.DisposalDecisionNo` | Conditional | Required if disposal is by company decision |

**GL Posting for FA Sale:**

```
Scenario: Proceeds > Carrying Amount (gain)
  Nợ 112 (Tiền gửi)                  5,000,000   (giá bán chưa VAT / sale price excl. VAT)
  Nợ 2142 (HM TSCĐ)             280,000,000   (xóa KH lũy kế / remove acc. dep.)
  Có 2112 (Nguyên giá TSCĐ)      280,000,000   (xóa nguyên giá / remove cost)
  Có 711 (Thu nhập khác)            5,000,000   (lãi thanh lý / disposal gain)

  Nợ 111/112 (Tiền)                    500,000   (VAT output)
  Có 3331 (Thuế GTGT đầu ra)           500,000   (VAT payable)
```

---

## 10. FA Import Excel Template

**File Name:** `fa_import_template.xlsx`

**Sheet Name:** `FA_Import`

**Column Specification:**

| Col | Field (EN) | Trường (VN) | Type | Max Len | Required | Validation Rules |
|-----|-----------|-------------|------|---------|----------|-----------------|
| A | `code` | Mã TSCĐ | Text | 50 | Yes | Unique per company. Regex: `/^[A-Z0-9\-_\/]+$/`. No spaces. |
| B | `name` | Tên TSCĐ | Text | 500 | Yes | Must not be empty. Min 2 chars. |
| C | `category_code` | Mã loại TSCĐ | Text | 20 | Yes | Must match existing category in `fixed_asset_categories` by company. |
| D | `acquisition_date` | Ngày mua | Date | — | Yes | Format: `DD/MM/YYYY`. Must not be future date. Must be ≤ in-service date. |
| E | `original_cost` | Nguyên giá | Number | — | Yes | Must be > 0. Max 2 decimal places. VND. |
| F | `useful_life_months` | Thời gian SD (tháng) | Integer | — | Yes | Must be > 0. Common values: 12, 24, 36, 48, 60, 72, 84, 96, 120, 180, 240. |
| G | `depreciation_method` | Phương pháp KH | Text | 30 | Yes | One of: `STRAIGHT_LINE`, `DECLINING_BALANCE`, `PRODUCTION_BASED`. |
| H | `residual_value` | Giá trị thanh lý | Number | — | Yes | Must be ≥ 0. Must be < original_cost. Default: 0. |
| I | `department_code` | Mã phòng ban | Text | 50 | Yes | Must match existing department in `departments` by company. |
| J | `location` | Địa điểm sử dụng | Text | 255 | No | Free text. E.g. "Lầu 3 - VP HCM". |
| K | `supplier_code` | Mã nhà cung cấp | Text | 50 | No | Must match existing supplier in `suppliers` by company. |
| L | `contract_no` | Số hợp đồng | Text | 100 | No | Free text. |
| M | `serial_no` | Số serial | Text | 100 | No | Free text. Recommend unique per asset. |
| N | `manufacturer` | Nhà sản xuất | Text | 255 | No | Free text. |
| O | `manufacture_year` | Năm sản xuất | Integer | — | No | Must be ≥ 1900 and ≤ current year. |
| P | `country_of_origin` | Nước sản xuất | Text | 100 | No | Free text. |
| Q | `source` | Nguồn hình thành | Text | 30 | Yes | One of: `PURCHASE`, `CONSTRUCTION`, `LEASE`, `DONATION`, `CAPITAL_CONTRIBUTION`, `EXCHANGE`. |
| R | `asset_account_id` | TK TSCĐ | Text | 20 | Yes | Must be valid COA account. Default per category. E.g. 2112. |
| S | `depreciation_account_id` | TK khấu hao | Text | 20 | Yes | Must be valid COA account. Default per category. E.g. 2142. |
| T | `expense_account_id` | TK chi phí KH | Text | 20 | Yes | Must be valid COA account. Default per category. E.g. 6424. |
| U | `notes` | Ghi chú | Text | 500 | No | Free text. |

**Sample Excel Rows:**

| A | B | C | D | E | F | G | H | I | J | K | L | M | N | O | P | Q | R | S | T | U |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| `FIT-2026-001` | `Máy tính Dell Latitude 5420` | `MMTB-03` | `15/07/2026` | `25000000` | `48` | `STRAIGHT_LINE` | `1250000` | `IT-DEPT` | `Lầu 3 - VP HCM` | `NCC001` | `HD-MT-123` | `DL5420-XYZ123456` | `Dell Inc.` | `2026` | `China` | `PURCHASE` | `2112` | `2142` | `6424` | `Nhập theo đợt đầu tư 2026` |
| `FIT-2026-002` | `MacBook Pro 14"` | `MMTB-03` | `20/07/2026` | `45000000` | `60` | `STRAIGHT_LINE` | `0` | `MKT-DEPT` | `Lầu 2 - Marketing` | `NCC001` | `HD-MT-124` | `MBP14-ABC987` | `Apple Inc.` | `2026` | `China` | `PURCHASE` | `2112` | `2142` | `6424` | `` |
| `HD-2026-001` | `Nhà xưởng sản xuất` | `NHA-01` | `01/06/2026` | `5000000000` | `240` | `STRAIGHT_LINE` | `250000000` | `PROD-DEPT` | `KCN Bình Dương` | `NCC002` | `HD-XD-001` | `` | `Cty Xây dựng ABC` | `2026` | `Vietnam` | `CONSTRUCTION` | `2111` | `2141` | `6274` | `Công trình xây dựng nhà xưởng` |

**Validation Error Response Example:**

```json
{
  "total_rows": 50,
  "success_count": 48,
  "error_count": 2,
  "errors": [
    {
      "row": 3,
      "field": "original_cost",
      "value": "abc",
      "error": "Invalid number format. Must be a positive numeric value."
    },
    {
      "row": 7,
      "field": "category_code",
      "value": "INVALID-CAT",
      "error": "Category code not found for this company."
    },
    {
      "row": 7,
      "field": "useful_life_months",
      "value": "0",
      "error": "Useful life must be greater than 0."
    }
  ]
}
```

**Import Processing Rules:**

1. Validate all rows before processing any — fail fast on first error per row, collect all row errors.
2. If `category_code` resolves to a category with default accounts, allow omitting `asset_account_id`, `depreciation_account_id`, `expense_account_id` (use category defaults).
3. `acquisition_date` must be ≤ `depreciation_start_date` (if provided separately). If no `depreciation_start_date`, default to following month start.
4. Auto-generate `status = DRAFT` for all imported assets (requires manual activation).
5. Auto-generate `company_id`, `created_by`, `created_at` from import context.
6. If `carrying_amount` is provided, it's informational only — computed as `original_cost - accumulated_depreciation`.
7. If `depreciation_start_date` is empty, derive as: `acquisition_date + 1 month, first day`.
8. If `depreciation_end_date` is empty, derive as: `depreciation_start_date + useful_life_months - 1 day`.
9. Transaction log: create `FATrxAcquisition` record for each imported asset referencing the import batch.

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-07-30 | BA Lead + Chief Accountant | Initial release — 10 templates |
