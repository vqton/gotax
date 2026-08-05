# Payroll Module — Templates & Configuration Reference

**Document ID:** TPL-PAYROLL-001
**Version:** 1.0
**Date:** 2026-08-04

---

## 1. Timekeeping Import Template

### CSV Format

```csv
employee_code,date,clock_in,clock_out,ot_hours,night_hours,leave_type,notes
E001,2026-08-01,08:00,17:00,0,0,,
E001,2026-08-02,08:00,17:00,2,0,,
E001,2026-08-03,22:00,06:00,0,8,,
E002,2026-08-01,08:00,17:00,0,0,ANNUAL,Leave approved
E003,2026-08-01,08:00,12:00,0,0,SICK,Sick leave certificate #123
```

### Column Definitions

| Column | Required | Format | Description |
|--------|----------|--------|-------------|
| employee_code | YES | String | Employee code from GoTax |
| date | YES | YYYY-MM-DD | Work date |
| clock_in | NO | HH:MM | Clock in time |
| clock_out | NO | HH:MM | Clock out time |
| ot_hours | NO | Decimal | Overtime hours |
| night_hours | NO | Decimal | Night shift hours (22:00-06:00) |
| leave_type | NO | Enum | ANNUAL, SICK, MATERNITY, UNPAID, COMPENSATORY |
| notes | NO | String | Additional notes |

---

## 2. Employee Payroll Info Import Template

### CSV Format

```csv
employee_code,contract_type,salary_type,base_salary,salary_coefficient,position_allowance,responsibility_allowance,seniority_allowance,other_allowances,insurance_base_salary,region,is_foreign_employee,is_trade_union_member,is_high_tech_talent,bank_account_no,bank_code
E001,INDEFINITE,TIME_BASED,15000000,0,1000000,500000,0,0,15000000,I,FALSE,TRUE,FALSE,1234567890,VCB
E002,INDEFINITE,COEFFICIENT,0,2.34,2000000,1000000,500000,0,25000000,I,FALSE,FALSE,FALSE,0987654321,CTG
E003,DEFINITE,TIME_BASED,20000000,0,1500000,0,0,0,20000000,II,TRUE,FALSE,FALSE,1122334455,VTB
```

---

## 3. Dependant Registration Template

### CSV Format

```csv
employee_code,full_name,relationship,date_of_birth,tax_code,is_disabled
E001,Nguyen Thi B,CHILD,2015-03-15,,
E001,Nguyen Van C,CHILD,2018-07-20,,
E002,Tran Thi D,SPOUSE,1990-01-10,1234567890,FALSE
```

---

## 4. Salary Component Templates

### Standard Vietnamese Payroll Components

| Code | Name | Type | Taxable | Insurable | Default |
|------|------|------|---------|-----------|---------|
| BS | Lương cơ bản | INCOME | Yes | Yes | - |
| PC_GT | Phụ cấp đi lại | INCOME | Yes | Yes | 500,000 |
| PC_AT | Phụ cấp ăn trưa | INCOME | Yes | Yes | 500,000 |
| PC_DT | Phụ cấp điện thoại | INCOME | Yes | No | 200,000 |
| PC_CN | Phụ cấp chức vụ | INCOME | Yes | Yes | - |
| PC_TN | Phụ cấp thâm niên | INCOME | Yes | Yes | - |
| PC_KV | Phụ cấp khu vực | INCOME | Yes | Yes | - |
| LT | Lương tăng ca | INCOME | Yes | No | - |
| LC | Lương ca đêm | INCOME | Yes | No | - |
| PHEP | Tiền phép | INCOME | Yes | Yes | - |
| TET | Thưởng Tết | INCOME | Yes | Yes | - |
| TNKL | Thưởng KPI | INCOME | Yes | Yes | - |

### Standard Deduction Components

| Code | Name | Type | Rate | Cap |
|------|------|------|------|-----|
| BHXH_NLĐ | BHXH phần NLĐ | DEDUCTION | 8% | 20 × base salary |
| BHYT_NLĐ | BHYT phần NLĐ | DEDUCTION | 1.5% | 20 × base salary |
| BHTN_NLĐ | BHTN phần NLĐ | DEDUCTION | 1% | 20 × regional min |
| CĐCS | Công đoàn cơ sở | DEDUCTION | 1% | 253,000 |
| TNCN | Thuế TNCN | DEDUCTION | Progressive | - |

### Standard Employer Cost Components

| Code | Name | Type | Rate | Cap |
|------|------|------|------|-----|
| BHXH_DNSH | BHXH phần使用lao động | EMPLOYER_COST | 17.5% | 20 × base salary |
| BHYT_DNSH | BHYT phần使用lao động | EMPLOYER_COST | 3% | 20 × base salary |
| BHTN_DNSH | BHTN phần使用lao động | EMPLOYER_COST | 1% | 20 × regional min |
| CĐCS_DNSH | CĐCS phần使用lao động | EMPLOYER_COST | 1% | 253,000 |

---

## 5. Configuration Parameters

### 5.1 Rate Configuration

```json
{
  "effective_from": "2026-01-01",
  "base_salary": {
    "2026-01-01_to_2026-06-30": 2340000,
    "2026-07-01_to_2026-12-31": 2530000
  },
  "si_hi_cap_multiplier": 20,
  "si_hi_cap": {
    "2026-01-01_to_2026-06-30": 46800000,
    "2026-07-01_to_2026-12-31": 50600000
  },
  "regional_minimum_wages": {
    "I": 5310000,
    "II": 4730000,
    "III": 4140000,
    "IV": 3700000
  },
  "ui_cap_multiplier": 20,
  "employee_rates": {
    "si": 0.08,
    "hi": 0.015,
    "ui": 0.01,
    "trade_union": 0.01
  },
  "employer_rates": {
    "si_retirement": 0.14,
    "si_sickness": 0.03,
    "si_accident": 0.005,
    "hi": 0.03,
    "ui": 0.01,
    "trade_union": 0.01
  },
  "employer_rates_foreign": {
    "si_retirement": 0.14,
    "si_sickness": 0.03,
    "si_accident": 0.005,
    "hi": 0.03,
    "ui": 0,
    "trade_union": 0.01
  },
  "pit_deductions": {
    "personal": 15500000,
    "per_dependant": 6200000,
    "supplementary_pension_cap": 3000000,
    "life_insurance_cap": 0
  },
  "pit_brackets": {
    "old_7bracket_2026H1": [
      {"min": 0, "max": 5000000, "rate": 0.05, "constant": 0},
      {"min": 5000000, "max": 10000000, "rate": 0.10, "constant": 250000},
      {"min": 10000000, "max": 18000000, "rate": 0.15, "constant": 750000},
      {"min": 18000000, "max": 32000000, "rate": 0.20, "constant": 1650000},
      {"min": 32000000, "max": 52000000, "rate": 0.25, "constant": 3250000},
      {"min": 52000000, "max": 80000000, "rate": 0.30, "constant": 5850000},
      {"min": 80000000, "max": null, "rate": 0.35, "constant": 9850000}
    ],
    "new_5bracket_2026H2": [
      {"min": 0, "max": 10000000, "rate": 0.05, "constant": 0},
      {"min": 10000000, "max": 30000000, "rate": 0.10, "constant": 500000},
      {"min": 30000000, "max": 60000000, "rate": 0.20, "constant": 3500000},
      {"min": 60000000, "max": 100000000, "rate": 0.30, "constant": 9500000},
      {"min": 100000000, "max": null, "rate": 0.35, "constant": 14500000}
    ]
  },
  "trade_union_cap": {
    "2026-01-01_to_2026-06-30": 234000,
    "2026-07-01_to_2026-12-31": 253000
  },
  "overtime_rates": {
    "weekday": 1.5,
    "rest_day": 2.0,
    "holiday": 3.0,
    "night_premium": 0.30,
    "night_ot_additional": 0.20
  },
  "leave_entitlements": {
    "annual": 12,
    "annual_per_5_years": 1,
    "sick_rate": 0.75,
    "maternity_months": 6,
    "paternity_days_normal": 5,
    "paternity_days_cesarean": 14
  }
}
```

### 5.2 Holiday Calendar (2026)

Per Labour Code 2019 Art. 112, Decree 128/2025/NĐ-CP, Official Letter 9859/VPCP-KGVX:

```json
{
  "year": 2026,
  "statutory_days": 11,
  "holidays": [
    {"date": "2026-01-01", "name": "New Year's Day (Tết Dương lịch)", "days": 1},
    {"date": "2026-02-16", "name": "Tet Eve (29 Tet)", "days": 1},
    {"date": "2026-02-17", "name": "Tet Nguyen Dan (Mung 1)", "days": 1},
    {"date": "2026-02-18", "name": "Tet Holiday (Mung 2)", "days": 1},
    {"date": "2026-02-19", "name": "Tet Holiday (Mung 3)", "days": 1},
    {"date": "2026-02-20", "name": "Tet Holiday (Mung 4)", "days": 1},
    {"date": "2026-04-26", "name": "Hung Kings Commemoration (Giỗ Tổ Hùng Vương)", "days": 1},
    {"date": "2026-04-30", "name": "Reunification Day (Ngày Giải phóng)", "days": 1},
    {"date": "2026-05-01", "name": "International Labour Day (Ngày Quốc tế Lao động)", "days": 1},
    {"date": "2026-09-01", "name": "National Day (Ngày Quốc khánh)", "days": 1},
    {"date": "2026-09-02", "name": "National Day (Ngày Quốc khánh)", "days": 1}
  ],
  "notes": {
    "tet_full_break": "14 Feb (Sat) - 22 Feb (Sun) = 9 days incl. weekends",
    "national_day_swap": "31 Aug (Mon) swapped to 22 Aug (Sat). Break: 29 Aug - 02 Sep",
    "private_sector_tet": "Employer chooses: 1+4 or 2+3 or 3+2 days (Decree 128/2025 Art. 7)",
    "private_sector_national_day": "2 Sep + 1 of (1 Sep or 3 Sep)"
  }
}
```

---

## 6. Payslip Template

### PDF Layout

```
┌─────────────────────────────────────────────────────────────────┐
│                    [COMPANY LOGO]                                │
│                    [COMPANY NAME]                                │
│                    [COMPANY ADDRESS]                             │
├─────────────────────────────────────────────────────────────────┤
│                         PHIẾU LƯƠNG                              │
│                    Tháng 07/2026                                 │
├─────────────────────────────────────────────────────────────────┤
│  Mã NV: E001                                                   │
│  Họ tên: Nguyễn Văn A                                          │
│  Phòng ban: Phòng Kế toán                                      │
│  Chức vụ: Kế toán viên                                         │
│  Mã số thuế: 1234567890                                        │
│  Số BHXH: BN123456789                                          │
│  Tài khoản NH: 1234567890 - VCB                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  I. CÁC KHOẢN THU NHẬP                                          │
│  ─────────────────────────────────────────────────────────────   │
│  Lương cơ bản .................... 15,000,000                   │
│  Phụ cấp đi lại .................  1,000,000                   │
│  Phụ cấp ăn trưa ................  1,000,000                   │
│  Lương tăng ca ..................    500,000                   │
│  Lương ca đêm ...................    200,000                   │
│  ─────────────────────────────────────────────────────────────   │
│  TỔNG THU NHẬP .................. 17,700,000                   │
│                                                                   │
│  II. CÁC KHOẢN KHẤU TRỪ                                         │
│  ─────────────────────────────────────────────────────────────   │
│  BHXH (8%) ....................... (1,240,000)                  │
│  BHYT (1.5%) .....................   (232,500)                  │
│  BHTN (1%) .......................   (155,000)                  │
│  Công đoàn (1%) ..................    (15,500)                  │
│  Thuế TNCN .......................   (275,000)                  │
│  ─────────────────────────────────────────────────────────────   │
│  TỔNG KHẤU TRỪ ................. (1,918,000)                  │
│                                                                   │
│  ═════════════════════════════════════════════════════════════   │
│  THỰC LĨNH ...................... 15,782,000                   │
│  ═════════════════════════════════════════════════════════════   │
│                                                                   │
│  III. PHÍA DOANH NGHIỆP                                         │
│  ─────────────────────────────────────────────────────────────   │
│  BHXH (17.5%) ...................  2,170,000                   │
│  BHYT (3%) ......................    465,000                   │
│  BHTN (1%) ......................    155,000                   │
│  Công đoàn (1%) ..................    15,500                   │
│  ─────────────────────────────────────────────────────────────   │
│  TỔNG CHI PHÍ NHÂN SỰ ........... 18,587,500                 │
│                                                                   │
├─────────────────────────────────────────────────────────────────┤
│  Ghi chú:                                                       │
│  - Ngày trả lương: 28/07/2026                                   │
│  - Kỳ tính thuế: Q3/2026                                        │
│  - NLĐ vui lòng kiểm tra trong vòng 5 ngày làm việc            │
│                                                                   │
│  Người lập                        Kế toán trưởng                │
│  (Ký, ghi rõ họ tên)             (Ký, ghi rõ họ tên)          │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 7. GL Account Codes

### Payroll-Related Accounts (Circular 99/2025/TT-BTC)

| Code | Name | Type | Description |
|------|------|------|-------------|
| 3331 | Phải trả người lao động | Liability | Salary payable |
| 3332 | Phải trả BHXH (DNSH) | Liability | Employer SI payable |
| 3333 | Phải trả BHYT (DNSH) | Liability | Employer HI payable |
| 3334 | Phải trả BHTN (DNSH) | Liability | Employer UI payable |
| 3335 | Phải trả BHXH (NLĐ) | Liability | Employee SI payable |
| 3336 | Phải trả BHYT (NLĐ) | Liability | Employee HI payable |
| 3337 | Phải trả BHTN (NLĐ) | Liability | Employee UI payable |
| 3338 | Phải trả thuế TNCN | Liability | PIT withholding payable |
| 3339 | Phải trả công đoàn | Liability | Trade union payable |
| 6421 | Chi phí lương | Expense | Salary expense |
| 6422 | Chi phí lương tăng ca | Expense | OT expense |
| 6423 | Chi phí phụ cấp | Expense | Allowance expense |
| 6424 | Chi phí bảo hiểm | Expense | Insurance expense |
| 6425 | Chi phí thưởng | Expense | Bonus expense |

---

## 8. Declaration Form References

### 8.1 D02-TS (Social Insurance Registration List)

**Purpose:** Register new employees, adjust contributions, terminate coverage
**Frequency:** Monthly (when changes occur)
**Submission:** Electronic via I-VAN or paper to district BHXH
**Deadline:** Within 30 days of contract signing

**Key Fields:**
- Unit name and code
- Employee name, SI number, date of birth
- Position/occupation
- Salary and allowances
- Contract type
- Effective date
- Type of change (new, adjustment, reduction)

### 8.2 TK3-TS (Employer Registration)

**Purpose:** Register employer for SI/HI participation
**Frequency:** One-time (with changes as needed)
**Submission:** Electronic or paper

**Key Fields:**
- Enterprise name, tax code, address
- Number of employees
- SI/HI unit code (assigned by BHXH)

### 8.3 05/KK-TNCN (Quarterly PIT Declaration)

**Purpose:** Report and remit PIT withheld from salaries
**Frequency:** Quarterly (from Q2/2026)
**Submission:** Electronic via eTax
**Deadline:** Last day of month following quarter end

**Key Fields:**
- Period (quarter/year)
- Total income per employee
- Total insurance deductions
- Total PIT withheld
- Exempt income (OT, night shift, leave)

---

## 9. Bank Salary File Format

### CSV Format (VCB-compatible)

```csv
STT,So tai khoan,Ho ten,So tien,Ghi chu
1,1234567890,Nguyen Van A,15782000,Luong 07/2026
2,0987654321,Tran Thi B,18500000,Luong 07/2026
3,1122334455,Le Van C,12300000,Luong 07/2026
```

### Fields

| Column | Description |
|--------|-------------|
| STT | Sequence number |
| So tai khoan | Bank account number |
| Ho ten | Employee full name |
| So tien | Net pay amount (VND) |
| Ghi chu | Reference (payslip number) |

---

## 10. Leave Types Configuration

| Code | Name | Paid | Max Days/Year | Accumulates | Encashable |
|------|------|------|---------------|-------------|------------|
| ANNUAL | Phép năm | Yes | 12 | Yes | Yes |
| SICK | Ốm đau | Yes (75%) | 30 | No | No |
| MATERNITY | Thai sản | Yes (100%) | 180 | No | No |
| PATERNITY | Thai sản (NĐ) | Yes (100%) | 5-14 | No | No |
| UNPAID | Không lương | No | Unlimited | No | No |
| COMPENSATORY | Nghỉ bù | Yes | Per policy | Yes | No |
| MARRIAGE | Kết hôn | Yes | 3 | No | No |
| BEREAVEMENT | Giỗ tang | Yes | 3 | No | No |
| FAMILY_CARE | Giữ con nhỏ | Yes | Per law | No | No |
