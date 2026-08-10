# GoTax vs MISA SME 2026 — Comprehensive Gap Analysis

**Date:** 2026-08-07  
**Purpose:** Full comparison of GoTax backend modules against MISA SME 2026 (latest version R6)  
**Scope:** All 19 MISA modules + Admin/Config system features

---

## Executive Summary

GoTax has **17/19 core accounting modules** implemented. Missing: Hợp đồng (Contracts), Khế ước vay (Loan Agreements). Admin/Config depth is ~40% of MISA SME 2026. E-Banking and E-Tax integrations are stub-only.

**Total Gap Score:**
- Core modules: 89% complete (17/19)
- Admin/Config: 40% complete
- Integration (E-Banking/E-Tax): 15% complete
- **Overall: ~72% of MISA SME 2026 parity**

---

## 1. Core Accounting Modules

| # | Module (MISA) | GoTax Status | Gap |
|---|--------------|-------------|-----|
| 1 | Hệ thống tài khoản (COA) | ✅ Full | None — COA versions, mappings, IFRS, freeze/unfreeze |
| 2 | Bút toán sổ cái (Journal) | ✅ Full | None — create/submit/review/approve/post/cancel |
| 3 | Số dư đầu kỳ (Opening Balance) | ✅ Full | None — import, correct, PDF export |
| 4 | Kết chuyển (Carry Forward) | ✅ Full | None — auto-close, logs |
| 5 | Báo cáo tài chính (Financial Reports) | ✅ Full | None — TB, BS, IS, journal export |
| 6 | Thuế (Tax) | ✅ Full | None — VAT, CIT, PIT, e-invoice |
| 7 | Mua hàng (Purchase) | ✅ Full | None — PO, GRN, AP, supplier inv |
| 8 | Bán hàng (Sale) | ✅ Full | None — SO, DN, AR, customer inv |
| 9 | Kho (Warehouse) | ✅ Full | None — receipts, issues, transfers, balances |
| 10 | Tiền mặt (Cash) | ✅ Full | None — receipt, payment, bank transfer |
| 11 | Ngân hàng (Bank) | ✅ Full | None — transaction import, reconciliation |
| 12 | Tài sản cố định (Fixed Asset) | ✅ Full | None — register, depreciation, disposal |
| 13 | Tiền lương (Payroll) | ✅ Full | None — grade, run, slip, SI/PIT |
| 14 | Ngân sách (Budget) | ✅ Full | None — budget vs actual |
| 15 | CCDC (Tools & Supplies) | ✅ Full | None — register, allocation |
| 16 | Trung tâm chi phí (Cost Centers) | ✅ Full | None — allocation |
| 17 | Giá thành (Cost Accounting) | ✅ Full | None — 6 methods, WIP, COGS, JE |
| 18 | **Hợp đồng (Contracts)** | ❌ **Missing** | **No contract tracking** |
| 19 | **Khế ước vay (Loan Agreements)** | ❌ **Missing** | **No loan principal/interest tracking** |

---

## 2. Admin/Config System (Hệ thống)

### 2.1 MISA SME 2026 System Features

| Feature | Description | GoTax Status | Gap |
|---------|------------|-------------|-----|
| **Thiết lập ngày hạch toán** | Set accounting start date per company | ❌ Missing | No accounting date config |
| **Quản lý tài nguyên** | Resource management (branches, departments) | ⚠️ Partial | Company exists, no branch/department tree |
| **Quản lý người dùng** | User CRUD + assign roles + branch | ⚠️ Partial | User CRUD exists, no branch assignment UI |
| **Vai trò và quyền hạn** | Role-based access + fine-grained permissions | ✅ Good | Casbin RBAC (admin/chief/user) |
| **Quản lý người dùng Mobile** | Mobile user management + payroll view | ❌ Missing | No mobile user concept |
| **Thiết lập gửi email** | SMTP config + email sending | ❌ Missing | No email system |
| **Thiết lập mẫu email** | Email templates | ❌ Missing | No email templates |
| **Nhật ký gửi email** | Email sending log | ❌ Missing | No email logging |
| **Nhắc nhở** | Reminders/notifications | ❌ Missing | Basic notification exists |
| **Thiết lập ký số** | Digital signature config | ❌ Missing | No digital sig config |
| **Xem nhật ký truy cập** | Access log viewer | ✅ Good | Audit middleware logs all |
| **Thiết lập tùy chọn** | System-wide options | ❌ **Critical gap** | See Section 2.2 |

### 2.2 System Options (Thiết lập tùy chọn) — CRITICAL MISSING

This is the **biggest gap**. MISA SME 2026 has a comprehensive system options panel:

| Option Category | MISA Features | GoTax Status |
|----------------|--------------|-------------|
| **Tùy chọn riêng** (Personal) | UI color, search filter, email address | ❌ Missing |
| **Tùy chọn chung** (Global) | Accounting mode, fiscal year, currency, branch config | ❌ Missing |
| **Báo cáo chứng từ** (Reports) | Company info on reports, logo, font, alignment | ❌ Missing |
| **Hóa đơn** (Invoice) | Invoice management, numbering, format | ❌ Missing |
| **Vật tư hàng hóa** (Inventory) | Costing method (FIFO/AVG/etc), branch costing | ❌ Missing |
| **Mua hàng** (Purchase) | Purchase options | ❌ Missing |
| **Bán hàng** (Sale) | Sale options | ❌ Missing |
| **Tiền lương** (Payroll) | Payroll options | ❌ Missing |
| **Định dạng số** (Number format) | Decimal places, separator, negative display | ❌ Missing |
| **Quy tắc đánh số CT** (Voucher numbering) | Auto-numbering rules per voucher type | ❌ Missing |
| **Ẩn/hiện nghiệp vụ** (Show/hide features) | Toggle feature visibility | ❌ Missing |
| **Sao lưu dữ liệu** (Backup) | Local + Google Drive backup | ❌ Missing |

### 2.3 Other Admin Features

| Feature | MISA | GoTax |
|---------|------|-------|
| Backup/restore | ✅ Local + Google Drive | ❌ Missing |
| Data import/export | ✅ Excel, CSV | ⚠️ Partial (journal export only) |
| Multi-branch support | ✅ Full branch tree | ⚠️ Partial (company_id based) |
| Fiscal year management | ✅ Configurable | ❌ Missing |
| Numbering rules | ✅ Per voucher type, per branch | ❌ Missing |
| Currency config | ✅ Multi-currency with rate management | ✅ Exchange rates exist |
| Report customization | ✅ Font, logo, alignment, headers | ❌ Missing |

---

## 3. Integration Modules

| Integration | MISA SME 2026 | GoTax Status |
|------------|--------------|-------------|
| **E-Banking** | ✅ MISA BankHub — auto-import bank statements from 40+ banks | ⚠️ Stub only |
| **E-Tax** | ✅ MISA mTax — auto-file PIT/CIT/BVAT to GDT | ⚠️ GDT client exists |
| **E-Invoice** | ✅ Multi-provider (meInvoice, VNPT, Viettel, BKAV) | ⚠️ Parse+generate XML only |
| **Digital Signature** | ✅ VNPT-CA, Viettel-CA, BKAV-CA | ❌ Missing |
| **POS Integration** | ✅ Cash register integration | ❌ Missing |
| **HRM Integration** | ✅ MISA HRM sync | ❌ Missing |

---

## 4. Priority Assessment

### Tier 1 — Critical (Accounting Compliance)
1. **System Options** — Without this, GoTax can't be configured per company needs
2. **Voucher Numbering Rules** — Required for audit trail compliance
3. **Fiscal Year Config** — Currently hardcoded, needs to be configurable
4. **Number Format Config** — Vietnamese format requirements

### Tier 2 — High (Business Operations)
5. **Contracts (Hợp đồng)** — Needed for project-based businesses
6. **Loan Agreements (Khế ước vay)** — Needed for financial management
7. **Backup/Restore** — Data safety requirement
8. **Report Customization** — Company info, logo on reports

### Tier 3 — Medium (Integration)
9. **E-Banking Integration** — Bank statement auto-import
10. **E-Tax Filing** — Auto-file to GDT
11. **Digital Signature** — For e-invoice signing
12. **Multi-branch enhancements** — Branch-level options

### Tier 4 — Low (Nice-to-have)
13. Email system (SMTP config, templates)
14. Mobile user management
15. Reminder/notification system
16. POS integration
17. HRM integration

---

## 5. Effort Estimates

| Module | Backend | Frontend | Tests | Total Days |
|--------|---------|----------|-------|-----------|
| System Options | 5 | 3 | 2 | **10** |
| Voucher Numbering | 3 | 2 | 1 | **6** |
| Fiscal Year Config | 2 | 1 | 1 | **4** |
| Number Format | 1 | 1 | 1 | **3** |
| Contracts | 4 | 3 | 2 | **9** |
| Loan Agreements | 4 | 3 | 2 | **9** |
| Backup/Restore | 3 | 1 | 1 | **5** |
| Report Customization | 2 | 2 | 1 | **5** |
| E-Banking | 5 | 3 | 2 | **10** |
| E-Tax Filing | 5 | 3 | 2 | **10** |
| Digital Signature | 3 | 1 | 1 | **5** |
| Multi-branch Enhance | 3 | 2 | 1 | **6** |
| **TOTAL** | | | | **~82 days** |

---

## 6. Recommendation

**Phase 1 (Weeks 1-3): System Foundation**
- System Options + Voucher Numbering + Fiscal Year + Number Format
- This unlocks all other modules' configurability

**Phase 2 (Weeks 4-6): Missing Business Modules**
- Contracts + Loan Agreements
- These complete the 19-module parity

**Phase 3 (Weeks 7-9): Integration**
- E-Banking + E-Tax + Digital Signature
- These enable real-world production use

**Phase 4 (Weeks 10-12): Polish**
- Backup/Restore + Report Customization + Multi-branch
- These improve UX and operational readiness
