# BRD: GoTax General Ledger Module — Production Readiness Assessment

## Document Control
| Field | Value |
|-------|-------|
| Version | 1.0 |
| Date | 2026-07-27 |
| Author | BA Lead + Chief Accountant (20+ yrs) |
| Status | DRAFT |
| Regulatory Basis | Circular 99/2025/TT-BTC (effective 1/1/2026) |

## 1. Executive Summary

### 1.1 Verdict
**The current GL module is NOT production-ready.** It is a functional prototype (v0.5) with fundamental architectural, regulatory, and feature gaps that prevent deployment in any real Vietnamese enterprise.

### 1.2 Critical Findings
| # | Severity | Finding | Impact |
|---|----------|---------|--------|
| 1 | 🔴 CRITICAL | Circular 200 COA (deprecated since 1/1/2026) | Non-compliant, audit failure |
| 2 | 🔴 CRITICAL | In-memory only — no persistent database | Data loss on restart |
| 3 | 🔴 CRITICAL | No financial statements (BS, IS, CF) | Cannot close books |
| 4 | 🔴 CRITICAL | No period close workflow | Cannot lock/finalize periods |
| 5 | 🔴 CRITICAL | No multi-currency | All Vietnamese businesses need it |
| 6 | 🔴 CRITICAL | No audit trail / transaction log | Unauditable |
| 7 | 🔴 CRITICAL | No approval workflow | No internal control |
| 8 | 🟠 HIGH | No e-tax/e-invoice/bank integration | Manual data entry |
| 9 | 🟠 HIGH | No user/role management | Shared secret model only |
| 10 | 🟠 HIGH | No budget tracking | Cannot control spending |

## 2. Regulatory Compliance (Circular 99/2025/TT-BTC)

### 2.1 Regulatory Landscape Change
**Circular 99/2025/TT-BTC** issued 27/10/2025, effective 1/1/2026, REPLACED:
- Circular 200/2014/TT-BTC
- Circular 75/2015/TT-BTC
- Circular 53/2016/TT-BTC
- Circular 195/2012/TT-BTC

### 2.2 Key Changes Affecting Module Design

#### Chart of Accounts Changes
```
Circular 200 (Current Module)  →  Circular 99 (Required)
76 Level-1 accounts               71 Level-1 accounts
No Account 215                    Account 215 - Biological Assets
No Account 332                    Account 332 - Dividends Payable
No Account 1383                   Account 1383 - SCT on imported goods
No Account 2295                   Account 2295 - Provision for biological assets
No Account 82111/82112            Account 82111/82112 - Pillar 2 GMT
No Account 2413/2414              Account 2413/2414 - Repair/Upgrade
Account 1562 abolished            Use Account 156 instead
Account 611 abolished             Use Account 154 instead
Account 3385 abolished            Use Account 138/338
Account 161/441/461/466 abolished Gone
```

#### Renamed Accounts
| Old Name (TT200) | New Name (TT99) |
|------------------|-----------------|
| Tiền gửi Ngân hàng (112) | Tiền gửi không kỳ hạn |
| Thành phẩm (155) | Sản phẩm |
| Chi phí trả trước (242) | Chi phí chờ phân bổ |
| Thặng dư vốn cổ phần (4112) | Thặng dư vốn |
| Bảng cân đối kế toán | Báo cáo tình hình tài chính |

#### Key Regulatory Requirements (Điều 28 TT99)
| Requirement | Current Module Status |
|-------------|----------------------|
| Accounting software must ensure data integrity | ⚠️ Partial (in-memory only) |
| Audit trail for all changes | ❌ Missing |
| Role-based access control | ❌ Missing |
| Ability to export data in standard formats | ❌ Missing |
| Automatic calculation and posting | ✅ Present (basic) |
| Support for internal accounting regulations | ❌ Missing |
| Integration with tax authority systems | ❌ Missing |

### 2.3 Legal References
- **Circular 99/2025/TT-BTC**: Enterprise accounting regime (effective 1/1/2026)
- **Circular 58/2026/TT-BTC**: Accounting regime for micro-enterprises (effective 1/7/2026)
- **Circular 133/2016/TT-BTC**: SME accounting regime (still active for SMEs)
- **Decision 345/QĐ-BTC**: IFRS application roadmap in Vietnam
- **Law on Accounting 2015** (amended 2024): Primary legal framework
- **Decree 70/2025/NĐ-CP**: Electronic invoices and digital transformation

## 3. Market Comparison

### 3.1 MISA AMIS (Latest Version R92 — May 2026)
| Feature | MISA AMIS | Current Module |
|---------|-----------|---------------|
| AI assistant (AVA) | ✅ AVA Accounting AI | ❌ |
| Auto bank transaction accounting | ✅ | ❌ |
| Auto e-invoice processing | ✅ | ❌ |
| E-tax declaration | ✅ | ❌ |
| Multi-branch consolidation | ✅ | ❌ |
| 200+ reports | ✅ | ❌ (trial balance only) |
| Period close | ✅ Automated | ❌ |
| Budget management | ✅ Basic | ❌ |
| Mobile app | ✅ | ❌ |
| Multi-currency | ✅ | ❌ |
| Circular 99 compliance | ✅ Since R89 (Jan 2026) | ❌ |
| Pricing | 2.45M-8.15M VND/yr | N/A |
| Users | 250,000+ enterprises | N/A |

### 3.2 FAST Accounting Online (Latest Version 2026)
| Feature | FAST | Current Module |
|---------|------|---------------|
| FastAI assistant | ✅ | ❌ |
| E-commerce integration (Shopee, TikTok, Lazada) | ✅ | ❌ |
| Bank integration (BIDV iConnect) | ✅ | ❌ |
| 300+ management reports | ✅ | ❌ |
| Circular 99 compliance | ✅ Since Dec 2025 | ❌ |
| 10 subsystems | ✅ | ❌ (1 partial) |
| Cloud-native | ✅ | ❌ (no DB) |
| Pricing | From 1.75M VND/yr | N/A |
| Multi-company | ✅ | ❌ |

### 3.3 BRAVO 10 ERP (Latest Version 2026)
| Feature | BRAVO | Current Module |
|---------|-------|---------------|
| Full ERP (12 modules) | ✅ | ❌ (GL only) |
| VAS + IFRS dual reporting | ✅ | ❌ |
| ISO 27001 certified | ✅ | ❌ |
| Customizable workflow | ✅ | ❌ |
| Cost calculation by multiple methods | ✅ | ❌ |
| Budget tracking | ✅ | ❌ |
| Multi-company consolidation | ✅ | ❌ |
| Circular 99 compliance | ✅ | ❌ |
| Target market | Mid-large enterprises | N/A |

### 3.4 Gap Analysis Summary
| Capability | MISA | FAST | BRAVO | Current Module |
|-----------|------|------|-------|---------------|
| Chart of Accounts | ✅ TT99 | ✅ TT99 | ✅ TT99 | ❌ TT200 |
| Journal Entry | ✅ | ✅ | ✅ | ✅ Basic |
| General Ledger | ✅ | ✅ | ✅ | ✅ Basic |
| Trial Balance | ✅ | ✅ | ✅ | ✅ Basic |
| Financial Statements | ✅ 4 reports | ✅ 4 reports | ✅ VAS+IFRS | ❌ |
| Period Close | ✅ Auto | ✅ Auto | ✅ Auto | ❌ |
| Multi-Currency | ✅ | ✅ | ✅ | ❌ |
| Budget Control | ✅ | ✅ | ✅ Budget+Forecast | ❌ |
| Audit Trail | ✅ Full | ✅ Full | ✅ Full log | ❌ |
| Approval Workflow | ✅ Multi-level | ✅ | ✅ Configurable | ❌ |
| Tax Integration | ✅ GDT | ✅ GDT | ✅ GDT | ❌ |
| Bank Integration | ✅ 20+ banks | ✅ BIDV | ✅ | ❌ |
| E-Invoice | ✅ | ✅ | ✅ | ❌ |
| AI/Automation | ✅ AVA | ✅ FastAI | ✅ AI analytics | ❌ |
| Mobile | ✅ | ✅ Web mobile | ✅ | ❌ |
| Security | ISO/Cloud | Cloud | ISO 27001 | ❌ (no auth) |
| Multi-company | ✅ | ✅ | ✅ | ❌ |

## 4. Detailed BRD — GL Module v1.0

### 4.1 Business Objectives
1. Replace manual Excel-based accounting with automated double-entry system
2. Comply with Circular 99/2025/TT-BTC
3. Produce statutory financial statements (B01-B09)
4. Support period close and audit readiness
5. Enable basic multi-currency operations

### 4.2 Scope (v1.0 MVP)
**IN SCOPE:**
- ✅ Chart of Accounts (Circular 99 compliant)
- ✅ Journal Entry (voucher-based, Chứng từ ghi sổ)
- ✅ General Ledger (Sổ Cái)
- ✅ Sub-ledgers (Sổ chi tiết)
- ✅ Trial Balance (Bảng cân đối số phát sinh)
- ✅ Financial Statements (B01-B04)
- ✅ Period Management (Open/Close/Lock)
- ✅ Multi-currency basic (VND + 1 foreign currency)
- ✅ Audit trail (timestamp + user)
- ✅ Basic approval (Draft → Review → Posted)
- ✅ User/Role management (Admin/Accountant/Viewer)
- ✅ Data export (Excel, PDF)

**OUT OF SCOPE (v1.1+):**
- ❌ E-tax declaration integration
- ❌ E-invoice generation
- ❌ Bank feed integration
- ❌ AI/ML automation
- ❌ Budget management
- ❌ Fixed asset management
- ❌ Inventory costing
- ❌ Payroll
- ❌ IFRS reporting (VAS only)

### 4.3 Functional Requirements

#### FR-01: Chart of Accounts (Hệ thống tài khoản)
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-01.01 | System provides default COA per Circular 99 Appendix II (71 Level-1 accounts) | P0 |
| FR-01.02 | Users may add/modify detail accounts without changing system accounts | P0 |
| FR-01.03 | Account code validation: numeric, hierarchical (e.g., 111 → 1111 → 11111) | P0 |
| FR-01.04 | Account type: ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE | P0 |
| FR-01.05 | Foreign currency flag per account | P1 |
| FR-01.06 | Detail-by tracking: OBJECT, PROJECT, CONTRACT, COST_ITEM, DEPARTMENT | P1 |
| FR-01.07 | Parent/child hierarchy with balance rollup | P0 |
| FR-01.08 | Import/export COA from Excel | P1 |
| FR-01.09 | Print COA listing | P2 |

#### FR-02: Vouchers / Journal Entry (Chứng từ kế toán)
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-02.01 | Create voucher with auto-numbering per type | P0 |
| FR-02.02 | Voucher types: Thu, Chi, Bán hàng, Mua hàng, Khác | P0 |
| FR-02.03 | Double-entry validation (Debit = Credit) | P0 |
| FR-02.04 | Multi-line support (unlimited lines per voucher) | P0 |
| FR-02.05 | Foreign currency amount + exchange rate per line | P1 |
| FR-02.06 | Voucher date, accounting date, period validation | P0 |
| FR-02.07 | Reference fields: invoice number, contract, PO | P1 |
| FR-02.08 | Attach documents (PDF, image) per voucher | P2 |
| FR-02.09 | Print voucher in Circular 99 format | P2 |
| FR-02.10 | Copy/reverse voucher | P1 |

#### FR-03: General Ledger (Sổ Cái & Sổ chi tiết)
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-03.01 | General ledger per account (Sổ Cái) with running balance | P0 |
| FR-03.02 | Sub-ledger per detail dimension (customer, vendor, project) | P1 |
| FR-03.03 | Period filter (month/quarter/year) | P0 |
| FR-03.04 | Export to Excel/PDF | P0 |
| FR-03.05 | Drill-down from balance to journal entry | P1 |

#### FR-04: Trial Balance (Bảng cân đối số phát sinh)
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-04.01 | Trial balance per period with opening/period/closing columns | P0 |
| FR-04.02 | Multi-level consolidation (parent accounts roll up children) | P0 |
| FR-04.03 | Year-to-date cumulative | P0 |
| FR-04.04 | Compare actual vs budget | P2 |

#### FR-05: Financial Statements (Báo cáo tài chính)
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-05.01 | B01-DN: Statement of Financial Position (Bảng cân đối kế toán → Báo cáo tình hình tài chính) | P0 |
| FR-05.02 | B02-DN: Income Statement (Báo cáo kết quả kinh doanh) | P0 |
| FR-05.03 | B03-DN: Cash Flow Statement (Báo cáo lưu chuyển tiền tệ) — Direct/Indirect | P0 |
| FR-05.04 | B09-DN: Notes to Financial Statements (Thuyết minh BCTC) | P1 |
| FR-05.05 | Automatic data pull from GL/TB | P0 |
| FR-05.06 | Print in Circular 99 format | P0 |
| FR-05.07 | Comparative columns (current year vs prior year) | P1 |

#### FR-06: Period Management (Quản lý kỳ kế toán)
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-06.01 | Automatic period creation (12 months per year) | P0 |
| FR-06.02 | OPEN → CLOSING → CLOSED state machine | P0 |
| FR-06.03 | Closing entries (Kết chuyển): 511→911, 911→421, etc. | P0 |
| FR-06.04 | Lock period: no new entries, no edits | P0 |
| FR-06.05 | Year-end close: carry forward balances to new year | P0 |
| FR-06.06 | Restrict: can only post to OPEN period | P0 |
| FR-06.07 | Warning when posting to prior period | P1 |

#### FR-07: Multi-Currency (Đa tiền tệ)
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-07.01 | Base currency: VND (mandatory) | P0 |
| FR-07.02 | Transaction currency: any currency | P1 |
| FR-07.03 | Exchange rate table (daily/monthly rates) | P1 |
| FR-07.04 | Auto-calculate VND equivalent at transaction rate | P1 |
| FR-07.05 | Year-end revaluation (Đánh giá lại cuối kỳ) per Circular 99 | P1 |
| FR-07.06 | Report both original currency and VND amounts | P1 |

#### FR-08: Audit Trail
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-08.01 | Every change logged: WHO, WHAT, WHEN, OLD_VALUE, NEW_VALUE | P0 |
| FR-08.02 | Immutable: once POSTED, no direct edit (only reversal) | P0 |
| FR-08.03 | Audit log viewer with filters | P0 |
| FR-08.04 | Tamper-evident log (append-only, hash chaining) | P2 |

#### FR-09: User & Role Management
| ID | Requirement | Priority |
|----|-------------|----------|
| FR-09.01 | Authentication: username/password (basic) | P0 |
| FR-09.02 | Roles: Admin, ChiefAccountant, Accountant, Viewer | P0 |
| FR-09.03 | Permission: create/edit/approve/post/cancel/delete/view | P0 |
| FR-09.04 | User per company/branch | P1 |
| FR-09.05 | Digital signature integration | P2 |

### 4.4 Business Rules

#### BR-01: Account Rules
1. Account code must be numeric, minimum 3 characters
2. Account code must be unique
3. Account must be ACTIVE to accept postings
4. Parent accounts cannot have direct postings (post to leaf accounts)
5. Parent-child hierarchy: up to 5 levels
6. Account type determines normal balance (Debit: Asset/Expense; Credit: Liability/Equity/Revenue)
7. Cannot delete account with non-zero balance or with children
8. Foreign currency accounts tracked separately

#### BR-02: Journal Entry Rules
1. Every entry must have at least 2 lines
2. Total Debits must equal Total Credits (tolerance: VND 1)
3. Entry date cannot be in the future
4. Entry date must fall within an OPEN period (or with special override)
5. Cannot post to a CLOSED/LOCKED period
6. Once POSTED, entry is immutable (no edit, no delete)
7. To correct a POSTED entry: create reversal entry + new corrected entry
8. Auto-numbering: prefix + year + sequential (e.g., CT-2026-00001)
9. Reference required for: invoice number (if applicable)

#### BR-03: Period Rules
1. Exactly 12 periods per fiscal year
2. Period transition: OPEN → CLOSING → CLOSED
3. CLOSING period: only closing entries allowed (type: K/C)
4. CLOSED period: no entries at all (read-only)
5. Year-end close: zero out revenue/expense accounts to 421
6. Cannot close period if subordinate period is open
7. Opening balance carry-forward mandatory at year start

#### BR-04: Audit Rules
1. No deletion of any posted entry (logical cancel only)
2. All status changes logged with timestamp and user
3. Audit log is append-only
4. Print audit log must show: user, time, action, old/new values

### 4.5 Workflows

#### WF-01: Journal Entry Workflow
```
[Start] → Create Entry (DRAFT) → Review → Approve → Post → [End]
               ↑                                      ↓
               └────── Reject (return to DRAFT) ──────┘
                    ↓
               CANCEL (delete without trace for DRAFT only)

States: DRAFT → REVIEWING → APPROVED → POSTED
Cancelled states: CANCELLED (from DRAFT only)
```

#### WF-02: Period Close Workflow
```
[Start Month N] → Normal Posting (OPEN)
                → Cut-off date reached
                → Post closing entries (CLOSING)
                → Generate Financial Statements
                → Verify Trial Balance = Zero for revenue/expense
                → Generate B01-B09
                → Lock period (CLOSED)
                → Carry forward balances to Month N+1
                → [End]
```

#### WF-03: Correction Workflow
```
[Error found in POSTED entry]
    → Create Reversal Entry (reverses all lines)
    → Post Reversal
    → Create Corrected Entry (correct amounts)
    → Post Corrected Entry
    → Note: both entries appear in audit trail
```

### 4.6 Data Model (Enhanced)

```
Company (1)
  ├── Users (N)
  ├── Chart of Accounts (N)
  ├── Periods (12/Year)
  ├── Journal Entries (N)
  │     └── Journal Lines (N)
  ├── Account Balances (N × Periods)
  ├── Financial Statements (4/Period)
  ├── Audit Log (N)
  ├── Exchange Rates (N)
  └── Closing Entries Templates (N)
```

### 4.7 User Journeys

#### UJ-01: Daily Accountant
1. Log in → Dashboard shows pending approvals
2. Create payment voucher → Enter vendor, amount, account
3. System validates: balance = debit + credit, account active, period open
4. Save as DRAFT → Submit for approval
5. Chief accountant approves → Status changes to POSTED
6. GL and Account Balance updated automatically

#### UJ-02: Month-End Close
1. Chief Accountant initiates period close
2. System changes period to CLOSING
3. Run closing entries (511→911, 911→421, etc.) automatically
4. Generate Trial Balance → Verify all accounts balance
5. Generate B01-B04 → Review
6. Lock period → CLOSED
7. Carry forward balances to next period

#### UJ-03: Audit Inquiry
1. Auditor requests transaction history
2. Export audit log for period (user, time, changes)
3. Drill down from Trial Balance → GL → Journal Entry
4. Verify no modifications to posted entries
5. Print Sổ Cái and Sổ chi tiết in Circular 99 format

### 4.8 Key Performance Indicators

| KPI | Target |
|-----|--------|
| Journal entry creation time | < 30 seconds |
| Trial balance generation | < 5 seconds (10k entries) |
| Financial statement generation | < 10 seconds |
| Period close process | < 1 hour |
| Concurrent users | 50+ |
| Data retention | 10 years |
| Audit log completeness | 100% of changes |
| Circular 99 compliance | 100% |

## 5. Gap Analysis — Current Module vs Production Requirements

| # | Area | Current | Required | Effort | Priority |
|---|------|---------|----------|--------|----------|
| 1 | Persistence | In-memory | PostgreSQL | 2 weeks | P0 |
| 2 | Circular 99 COA | TT200 | Update seed data | 1 day | P0 |
| 3 | Financial Statements | None | B01-B04 | 4 weeks | P0 |
| 4 | Period Close | Status field | Full workflow | 3 weeks | P0 |
| 5 | Multi-currency | None | Full support | 4 weeks | P1 |
| 6 | Audit Trail | None | Append-only log | 2 weeks | P0 |
| 7 | User/RBAC | None | Auth + roles | 2 weeks | P0 |
| 8 | Approval Workflow | Draft→Posted | Multi-step | 2 weeks | P0 |
| 9 | Closing Entries | None | Auto templates | 1 week | P0 |
| 10 | Sub-ledgers | None | Customer/Vendor | 3 weeks | P1 |
| 11 | E-tax integration | None | GDT API | 6 weeks | P2 |
| 12 | Bank integration | None | API | 6 weeks | P2 |
| 13 | E-invoice | None | GTGT | 6 weeks | P2 |
| 14 | Budget management | None | Full | 4 weeks | P2 |
| 15 | Mobile app | None | Basic view | 4 weeks | P3 |

### Estimated Production Timeline
- **Phase 0** (Critical fixes): 3 weeks — COA, DB, Audit, Auth
- **Phase 1** (Core features): 8 weeks — Financial Statements, Period Close, Approval
- **Phase 2** (Enhancement): 6 weeks — Multi-currency, Sub-ledgers, Closing Entries
- **Phase 3** (Integration): 12 weeks — E-tax, Bank, E-invoice
- **Phase 4** (Advanced): 8 weeks — Budget, AI, Mobile

**Total to production: ~37 weeks (9 months) with dedicated team of 4-5 engineers**

## 6. Risk Assessment

### 6.1 Regulatory Risks
| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Circular 99 non-compliance | HIGH | CRITICAL | Immediate COA update, legal review |
| IFRS convergence requirements | MEDIUM | HIGH | Monitor MOF roadmap, build flexible schema |
| Tax law changes (GMT, VAT) | MEDIUM | HIGH | Parameterize tax accounts, easy update path |
| Data retention non-compliance | MEDIUM | MEDIUM | Implement archival strategy from day 1 |

### 6.2 Technical Risks
| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Database migration from in-memory | HIGH | MEDIUM | Start with SQLite for dev, PG for prod |
| Data loss during restart | CERTAIN | HIGH | Add persistence immediately |
| Audit trail tampering | MEDIUM | CRITICAL | Append-only log with hash chain |
| Concurrent access issues | MEDIUM | MEDIUM | Proper transaction management |

### 6.3 Competitive Risks
| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| MISA/FAST more feature-rich | HIGH | MEDIUM | Focus on simplicity + price |
| Users unwilling to switch | HIGH | HIGH | Seamless data import from Excel/MISA |
| Circular 99 updates late | MEDIUM | HIGH | Engage MOF for early notifications |

## 7. Architecture Roadmap

```
Phase 0 (Weeks 1-3):
┌─────────────┐  ┌──────────────┐  ┌────────────┐  ┌─────────────┐
│ PostgreSQL  │  │ Audit Log    │  │ Auth/RBAC  │  │ COA TT99    │
│ Persistence │  │ (append-only)│  │ (JWT/Roles)│  │ (seed data) │
└─────────────┘  └──────────────┘  └────────────┘  └─────────────┘

Phase 1 (Weeks 4-11):
┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ B01-B04      │  │ Period Close │  │ Approval     │  │ Closing      │
│ Statements   │  │ Workflow     │  │ Workflow     │  │ Templates    │
└──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘

Phase 2 (Weeks 12-17):
┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Multi-       │  │ Sub-ledgers  │  │ Exchange     │  │ Year-end     │
│ Currency     │  │ Cust/Vendor  │  │ Rate Mgmt    │  │ Close        │
└──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘

Phase 3 (Weeks 18-29):
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ E-tax (GDT)  │  │ Bank Feed    │  │ E-invoice    │
│ Integration  │  │ Integration  │  │ GTGT         │
└──────────────┘  └──────────────┘  └──────────────┘
```

## 8. Appendices

### A. Circular 99 Account List (71 Level-1 Accounts)
Source: Appendix II, Circular 99/2025/TT-BTC

**Type 1: ASSETS (Tài sản) — 33 accounts**
111, 112, 113, 121, 128, 131, 133, 136, 138, 141, 151, 152, 153, 154, 155, 156, 157, 158, 171, 211, 212, 213, 214, 215, 217, 221, 222, 228, 229, 241, 242, 243, 244

**Type 2: LIABILITIES (Nợ phải trả) — 14 accounts**
331, 332, 333, 334, 335, 336, 337, 338, 341, 343, 344, 347, 352, 353

**Type 3: EQUITY (Vốn chủ sở hữu) — 9 accounts**
411, 412, 413, 414, 415, 417, 418, 419, 421

**Type 4: REVENUE (Doanh thu) — 4 accounts**
511, 515, 521, 711

**Type 5: EXPENSES (Chi phí) — 9 accounts**
621, 622, 623, 627, 632, 635, 641, 642, 811

**Type 6: DETERMINATION OF RESULTS (Xác định KQKD) — 2 accounts**
821, 911

### B. Financial Statement Forms (Circular 99)
| Form | Name | Frequency |
|------|------|-----------|
| B01-DN | Statement of Financial Position | Quarterly/Annually |
| B02-DN | Statement of Profit or Loss | Quarterly/Annually |
| B03-DN | Cash Flow Statement | Annually |
| B09-DN | Notes to Financial Statements | Annually |

### C. Voucher Templates (Circular 99, Appendix I)
| Form | Name |
|------|------|
| 01-TT | Receipt (Phiếu thu) |
| 02-TT | Payment (Phiếu chi) |
| 03-TT | Cash/Cheque (Phiếu nộp tiền/Séc) |
| 04-TT | Petty Cash (Giấy đề nghị tạm ứng) |
| 05-TT | Advance Settlement (Giấy thanh toán tạm ứng) |
| 06-TT | Payroll (Bảng chấm công, Bảng thanh toán lương) |
| 01-BH | Goods Receipt (Phiếu nhập kho) |
| 02-BH | Goods Issue (Phiếu xuất kho) |

### D. Closing Entries (Bút toán kết chuyển)
| # | Entry | Debit | Credit |
|---|-------|-------|--------|
| 1 | Close revenue | 511, 515, 711 | 911 |
| 2 | Close expenses | 911 | 621, 622, 623, 627, 632, 635, 641, 642, 811 |
| 3 | Close other income/expense | 711 | 911 / 911 → 811 |
| 4 | Close profit to retained earnings | 911 | 421 |
| 5 | Close loss to retained earnings | 421 | 911 |
| 6 | VAT deduction transfer | 33311 | 133 |
| 7 | Provision adjustments | 642 | 229 (or reverse) |

---

## Document References
- Circular 99/2025/TT-BTC (MOF, 27/10/2025, effective 1/1/2026)
- Circular 58/2026/TT-BTC (MOF, 25/5/2026, effective 1/7/2026)
- KPMG Vietnam Financial Reporting Alert (Nov 2025)
- Grant Thornton Vietnam IFRS Viewpoint (Dec 2025)
- EY Vietnam Webinar on Circular 99 (Dec 2025)
- Forvis Mazars Vietnam Accounting Update (Mar 2026)
- MISA AMIS Release Notes R85-R92 (2025-2026)
- FAST Accounting Online Product Documentation (2026)
- BRAVO 10 ERP Product Documentation (2026)
