# GoTax Company Module — Templates

## Version: 1.0 | Date: 2026-07-27

---

## 1. Company Registration Form

```json
{
  "legal_name_vn": "CONG TY TRACH NHIEM HUU HAN ABC",
  "legal_name_en": "ABC COMPANY LIMITED",
  "short_name": "ABC Co., Ltd",
  "legal_form": "LLC_1MEMBER",
  "tax_code": "0123456789",
  "business_reg_no": "0101234567",
  "business_reg_date": "2026-01-15",
  "business_reg_place": "So Ke hoach Dau tu TP Ha Noi",

  "reg_address": "So 1, Duong A, Phuong B, Quan C, TP Ha Noi",
  "reg_address_province": "Ha Noi",
  "reg_address_district": "Quan C",
  "head_office_address": "Tang 5, Toa nha XYZ, So 10, Duong D, Quan E, TP Ha Noi",
  "head_office_province": "Ha Noi",
  "head_office_district": "Quan E",

  "phone": "+84 24 1234 5678",
  "email": "info@abc.vn",
  "website": "https://abc.vn",

  "legal_rep_name": "Nguyen Van A",
  "legal_rep_title": "Giam doc",
  "legal_rep_id_number": "0123456789",
  "legal_rep_id_date": "2010-05-20",
  "legal_rep_id_place": "CA TP Ha Noi",

  "chief_accountant_name": "Tran Thi B",
  "chief_accountant_email": "accountant@abc.vn",

  "tax_office_code": "01",
  "tax_office_name": "Chi cuc Thue quan Dong Da",

  "accounting_regime": "TT99",
  "fiscal_year_start_month": 1,
  "default_currency": "VND",
  "secondary_currency": "USD",

  "company_type": "TRADING",
  "company_size": "SMALL"
}
```

---

## 2. Company Profile Response

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "tenant_id": "t0000-0000-0000-000000000001",
  "legal_name_vn": "CONG TY TRACH NHIEM HUU HAN ABC",
  "legal_name_en": "ABC COMPANY LIMITED",
  "short_name": "ABC Co., Ltd",
  "legal_form": "LLC_1MEMBER",
  "legal_form_label": "Cong ty TNHH 1 thanh vien",
  "tax_code": "0123456789",
  "business_reg_no": "0101234567",
  "business_reg_date": "2026-01-15",
  "business_reg_place": "So Ke hoach Dau tu TP Ha Noi",
  "reg_address": "So 1, Duong A, Phuong B, Quan C, TP Ha Noi",
  "head_office_address": "Tang 5, Toa nha XYZ, So 10, Duong D, Quan E, TP Ha Noi",
  "phone": "+84 24 1234 5678",
  "email": "info@abc.vn",
  "website": "https://abc.vn",
  "legal_rep_name": "Nguyen Van A",
  "legal_rep_title": "Giam doc",
  "chief_accountant_name": "Tran Thi B",
  "chief_accountant_email": "accountant@abc.vn",
  "tax_office_code": "01",
  "tax_office_name": "Chi cuc Thue quan Dong Da",
  "accounting_regime": "TT99",
  "fiscal_year_start_month": 1,
  "default_currency": "VND",
  "secondary_currency": "USD",
  "company_type": "TRADING",
  "company_size": "SMALL",
  "status": "ACTIVE",
  "parent_company_id": null,
  "hierarchy_path": "a1b2c3d4",
  "logo_url": "https://cdn.gotax.vn/logos/a1b2c3d4.png",
  "statistics": {
    "branch_count": 2,
    "employee_count": 25,
    "department_count": 5,
    "bank_account_count": 3,
    "einvoice_pattern_count": 2,
    "digital_signature_count": 1,
    "current_fiscal_year": 2026,
    "current_period": "Thang 07/2026",
    "current_period_status": "OPEN"
  },
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-07-27T10:00:00Z"
}
```

---

## 3. Company Hierarchy Tree

```json
{
  "company": {
    "id": "a1b2",
    "legal_name_vn": "CONG TY TNHH ABC",
    "tax_code": "0123456789",
    "legal_form": "LLC_1MEMBER",
    "type": "HOLDING",
    "status": "ACTIVE"
  },
  "children": [
    {
      "id": "c3d4",
      "legal_name_vn": "CHI NHANH HA NOI - CONG TY TNHH ABC",
      "tax_code": "0123456789-001",
      "type": "BRANCH",
      "status": "ACTIVE",
      "children": [
        {
          "id": "e5f6",
          "name": "Phong Kinh doanh",
          "type": "DEPARTMENT",
          "manager": "Nguyen Van C"
        },
        {
          "id": "g7h8",
          "name": "Phong Ke toan",
          "type": "DEPARTMENT",
          "manager": "Tran Thi D"
        }
      ]
    },
    {
      "id": "i9j0",
      "legal_name_vn": "CONG TY TNHH ABC LOGISTICS",
      "tax_code": "0987654321",
      "legal_form": "LLC_1MEMBER",
      "type": "SUBSIDIARY",
      "ownership_pct": 80,
      "status": "ACTIVE",
      "children": []
    }
  ]
}
```

---

## 4. Fiscal Year & Period Template

```json
{
  "fiscal_year": {
    "id": "fy-uuid",
    "company_id": "a1b2c3d4",
    "year": 2026,
    "start_month": 1,
    "is_short_year": false,
    "start_date": "2026-01-01",
    "end_date": "2026-12-31",
    "status": "OPEN"
  },
  "periods": [
    {
      "id": "p01-uuid",
      "period_type": "MONTHLY",
      "period_number": 1,
      "label": "Thang 01/2026",
      "start_date": "2026-01-01",
      "end_date": "2026-01-31",
      "status": "CLOSED",
      "closed_at": "2026-02-05T17:00:00Z"
    },
    {
      "id": "p02-uuid",
      "period_type": "MONTHLY",
      "period_number": 2,
      "label": "Thang 02/2026",
      "start_date": "2026-02-01",
      "end_date": "2026-02-28",
      "status": "CLOSED",
      "closed_at": "2026-03-05T17:00:00Z"
    },
    {
      "id": "p03-uuid",
      "period_type": "MONTHLY",
      "period_number": 3,
      "label": "Thang 03/2026",
      "start_date": "2026-03-01",
      "end_date": "2026-03-31",
      "status": "CLOSED",
      "closed_at": "2026-04-06T17:00:00Z"
    },
    {
      "id": "p04-uuid",
      "period_type": "MONTHLY",
      "period_number": 4,
      "label": "Thang 04/2026",
      "start_date": "2026-04-01",
      "end_date": "2026-04-30",
      "status": "CLOSED",
      "closed_at": "2026-05-05T17:00:00Z"
    },
    {
      "id": "p05-uuid",
      "period_type": "MONTHLY",
      "period_number": 5,
      "label": "Thang 05/2026",
      "start_date": "2026-05-01",
      "end_date": "2026-05-31",
      "status": "CLOSED",
      "closed_at": "2026-06-04T17:00:00Z"
    },
    {
      "id": "p06-uuid",
      "period_type": "MONTHLY",
      "period_number": 6,
      "label": "Thang 06/2026",
      "start_date": "2026-06-01",
      "end_date": "2026-06-30",
      "status": "CLOSED",
      "closed_at": "2026-07-04T17:00:00Z"
    },
    {
      "id": "p07-uuid",
      "period_type": "MONTHLY",
      "period_number": 7,
      "label": "Thang 07/2026",
      "start_date": "2026-07-01",
      "end_date": "2026-07-31",
      "status": "OPEN"
    },
    {
      "id": "p08-uuid",
      "period_type": "MONTHLY",
      "period_number": 8,
      "label": "Thang 08/2026",
      "start_date": "2026-08-01",
      "end_date": "2026-08-31",
      "status": "FUTURE"
    },
    {
      "id": "p09-uuid",
      "period_type": "MONTHLY",
      "period_number": 9,
      "label": "Thang 09/2026",
      "start_date": "2026-09-01",
      "end_date": "2026-09-30",
      "status": "FUTURE"
    },
    {
      "id": "p10-uuid",
      "period_type": "MONTHLY",
      "period_number": 10,
      "label": "Thang 10/2026",
      "start_date": "2026-10-01",
      "end_date": "2026-10-31",
      "status": "FUTURE"
    },
    {
      "id": "p11-uuid",
      "period_type": "MONTHLY",
      "period_number": 11,
      "label": "Thang 11/2026",
      "start_date": "2026-11-01",
      "end_date": "2026-11-30",
      "status": "FUTURE"
    },
    {
      "id": "p12-uuid",
      "period_type": "MONTHLY",
      "period_number": 12,
      "label": "Thang 12/2026",
      "start_date": "2026-12-01",
      "end_date": "2026-12-31",
      "status": "FUTURE"
    },
    {
      "id": "q1-uuid",
      "period_type": "QUARTERLY",
      "period_number": 1,
      "label": "Quy 01/2026",
      "start_date": "2026-01-01",
      "end_date": "2026-03-31",
      "status": "CLOSED"
    },
    {
      "id": "q2-uuid",
      "period_type": "QUARTERLY",
      "period_number": 2,
      "label": "Quy 02/2026",
      "start_date": "2026-04-01",
      "end_date": "2026-06-30",
      "status": "CLOSED"
    },
    {
      "id": "q3-uuid",
      "period_type": "QUARTERLY",
      "period_number": 3,
      "label": "Quy 03/2026",
      "start_date": "2026-07-01",
      "end_date": "2026-09-30",
      "status": "FUTURE"
    },
    {
      "id": "q4-uuid",
      "period_type": "QUARTERLY",
      "period_number": 4,
      "label": "Quy 04/2026",
      "start_date": "2026-10-01",
      "end_date": "2026-12-31",
      "status": "FUTURE"
    }
  ]
}
```

---

## 5. E-Invoice Pattern Registration

```json
{
  "pattern_code": "01GTKT0/001",
  "serial": "AA/26E",
  "form": "Hoa don gia tri gia tang",
  "invoice_type": "GTGT",
  "description": "Hoa don GTGT ban ra cho khach hang thong thuong"
}
```

### Common Pattern Codes (Circular 99 + Decree 123/2020)
| Code | Description | Type |
|------|-------------|------|
| 01GTKT0/001 | Hoa don GTGT (Value-added invoice) | GTGT |
| 02GTTT0/001 | Hoa don ban hang (Sales invoice) | RETAIL |
| 03XKNB0/001 | Hoa don xuat khau (Export invoice) | EXPORT |
| 04CKTB0/001 | Hoa don cho ky tinh thue (Non-taxable) | SERVICE |
| 07KPTQ0/001 | Hoa don khu phi thuong mai (Non-commercial) | OTHER |

---

## 6. Digital Signature Registration

### USB Token
```json
{
  "signature_type": "USB_TOKEN",
  "serial_number": "4A:5B:6C:7D:8E:9F:00:11:22:33",
  "owner_name": "Nguyen Van A - Giam doc",
  "provider": "ViettelCA",
  "valid_from": "2026-01-01",
  "valid_to": "2027-12-31",
  "is_default": true
}
```

### Remote HSM
```json
{
  "signature_type": "REMOTE_HSM",
  "serial_number": "HSM-VTT-2026-001234",
  "owner_name": "Cong ty TNHH ABC",
  "provider": "ViettelCA",
  "valid_from": "2026-01-01",
  "valid_to": "2027-12-31",
  "is_default": true,
  "config": {
    "api_endpoint": "https://api.ca.viettel.vn/signing",
    "certificate_alias": "abc-company-2026"
  }
}
```

---

## 7. Bank Account Registration

```json
{
  "bank_code": "VCB",
  "bank_name": "Ngan hang TMCP Ngoai thuong Viet Nam (Vietcombank)",
  "branch_name": "So giao dich So 1, Hoan Kiem, Ha Noi",
  "account_number": "0011001234567",
  "account_holder": "CONG TY TNHH ABC",
  "currency": "VND",
  "is_default": true
}
```

### Common Vietnamese Bank Codes (Napas)
| Code | Bank Name |
|------|-----------|
| VCB | Vietcombank |
| CTG | VietinBank |
| BIDV | BIDV |
| VIB | VIB |
| TCB | Techcombank |
| ACB | ACB |
| MB | MB Bank |
| SHB | SHB |
| HDB | HDBank |
| VPB | VPBank |
| LPB | LienVietPostBank |
| MSB | MSB |
| STB | Sacombank |
| EIB | Eximbank |
| OCB | OCB |
| TPBank | TPBank |
| BVB | BaoViet Bank |
| SCB | SCB |
| NCB | NCB |
| KLB | KienLong Bank |

---

## 8. Employee Import CSV Template

```csv
employee_code,full_name,title,department_code,email,phone,personal_tax_code,social_insurance_no,hire_date
NV001,Nguyen Van A,Giam doc,SALES,nguyenvana@abc.vn,0901234567,0987654321,BHXH001,2020-01-01
NV002,Tran Thi B,Ke toan truong,ACCOUNTING,tranthib@abc.vn,0901234568,0987654322,BHXH002,2021-03-15
NV003,Le Van C,Nhan vien ban hang,SALES,levanc@abc.vn,0901234569,0987654323,BHXH003,2022-06-01
NV004,Pham Thi D,Nhan vien ke toan,ACCOUNTING,phamthid@abc.vn,0901234570,0987654324,BHXH004,2023-02-01
NV005,Hoang Van E,Nhan vien kho,WAREHOUSE,hoangvane@abc.vn,0901234571,0987654325,BHXH005,2023-09-01
```

---

## 9. Company Context Response (Post-Switch)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "company": {
    "id": "a1b2c3d4",
    "legal_name_vn": "CONG TY TNHH ABC",
    "tax_code": "0123456789",
    "accounting_regime": "TT99",
    "default_currency": "VND",
    "fiscal_year": 2026,
    "current_period": {
      "id": "p07-uuid",
      "label": "Thang 07/2026",
      "status": "OPEN"
    }
  },
  "permissions": [
    "journal:write",
    "journal:read",
    "account:read",
    "report:read",
    "company:read"
  ]
}
```

---

## 10. Integration Profile (GDT)

```json
{
  "company_id": "a1b2c3d4",
  "integration_type": "GDT",
  "endpoint_url": "https://ihtkk.gdt.gov.vn",
  "config": {
    "tax_office_code": "01",
    "tax_code": "0123456789",
    "username": "mst0123456789",
    "connection_timeout_secs": 30,
    "retry_count": 3,
    "retry_delay_secs": 5
  },
  "status": "CONNECTED",
  "last_connected_at": "2026-07-27T09:00:00Z"
}
```

Note: Credentials field not shown in API response (returned as masked: `****`). Stored encrypted at rest.
