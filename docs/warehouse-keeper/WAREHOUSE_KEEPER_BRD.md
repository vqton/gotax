# Business Requirements Document — Thủ Kho (Warehouse Keeper) Module

**Version:** 1.0
**Date:** August 2026
**Author:** BA Lead + Chief Accountant (20+ years combined)
**Status:** DRAFT
**Regulatory Basis:** Luật Kế toán 2015, Thông tư 99/2025/TT-BTC, Decree 123/2020/ND-CP

---

## 1. Executive Summary

Thủ Kho (Warehouse Keeper) is a **role-based sub-module** within the Warehouse module that enforces **separation of duties** between physical goods custody (Thủ kho) and accounting records (Kế toán kho). This is a mandatory internal control principle under Vietnamese Law on Accounting 2015, Article 6 and Circular 99/2025/TT-BTC, Article 3.

**Problem:** GoTax's current warehouse module has no concept of warehouse keeper role. All users see all data, no custody tracking, no keeper sign-off workflow. This violates segregation of duties and makes the system unsuitable for production deployment in Vietnamese enterprises.

**Solution:** Add a Thủ kho sub-module that:
1. Assigns users as warehouse keepers to specific warehouses
2. Provides a keeper-specific review & recording workflow (Ghi sổ)
3. Maintains a separate Stock Ledger (Sổ kho) per keeper
4. Generates reconciliation reports between keeper and accounting
5. Supports physical inventory counts with keeper participation

**Benchmark:** MISA SME 2026 Thủ kho module (module #17/18), MISA AMIS Kế toán Thủ kho workflow.

---

## 2. Business Context

### 2.1 Regulatory Requirements

| Regulation | Article/Section | Requirement |
|------------|----------------|-------------|
| **Luật Kế toán 2015** | Article 6 | Separation of duties: person managing assets must not simultaneously maintain accounting records |
| **Luật Kế toán 2015** | Article 16 | Accounting vouchers must be complete, accurate, timely; each voucher created once per transaction |
| **Circular 99/2025/TT-BTC** | Article 3 | Enterprises must establish internal governance regulations defining rights, obligations, responsibilities of departments/individuals |
| **Circular 99/2025/TT-BTC** | Article 10 | Every economic/financial transaction must have an accounting voucher; voucher signing authority must ensure strict control |
| **Circular 99/2025/TT-BTC** | Article 16 | Inventory accounting: track quantity and value per warehouse, per department |
| **Decree 123/2020/ND-CP** | Article 4 | Enterprise internal control requirements |

### 2.2 Business Drivers

1. **Regulatory compliance** — Law on Accounting 2015 mandates separation of duties
2. **Internal control** — Cross-checking between keeper and accountant detects discrepancies
3. **Audit readiness** — Auditors expect evidence of segregation (Điều 6 Luật Kế toán 2015)
4. **MISA parity** — MISA SME 2026 includes Thủ kho as standard module in all tiers
5. **Fraud prevention** — Single-person custody + recording = undetectable manipulation

### 2.3 Scope

**In Scope:**
- Warehouse keeper assignment to warehouses
- Keeper review & recording workflow (Ghi sổ / Bỏ ghi)
- Stock Ledger (Sổ kho) per keeper
- Stock Card (Thẻ kho) per item per keeper
- Reconciliation report (Đối chiếu kế toán - Thủ kho)
- Physical inventory count with keeper participation
- Keeper-specific reports

**Out of Scope:**
- Warehouse module enhancements (GRN, DN, transfers, adjustments — already PROD)
- Cost price calculation (handled by warehouse valuation module)
- Goods receipt note creation (handled by Purchase module)
- Delivery note creation (handled by Sale module)

---

## 3. User Roles

### 3.1 Thủ kho (Warehouse Keeper)

| Attribute | Value |
|-----------|-------|
| **Role** | Physical goods custodian |
| **Responsibility** | Receive, store, preserve, issue goods; maintain Stock Ledger |
| **Access** | Read-only on source documents; Record/Un-record authority |
| **Data visibility** | Quantities visible; cost prices optionally hidden |
| **Cannot do** | Create/edit/delete PNK/PXK; handle money; create accounting entries |
| **Reports** | Stock Ledger, Stock Card, Inventory Count Minutes, Reconciliation |

### 3.2 Kế toán kho (Warehouse Accountant)

| Attribute | Value |
|-----------|-------|
| **Role** | Accounting records for inventory |
| **Responsibility** | Create journal entries, reconcile accounts, generate financial reports |
| **Access** | Full financial data (quantities + values + costs) |
| **Cannot do** | Physically receive/issue goods; maintain Stock Ledger |

### 3.3 Quản lý kho (Warehouse Manager)

| Attribute | Value |
|-----------|-------|
| **Role** | Supervisory oversight |
| **Responsibility** | Approve transfers, adjustments; review reconciliation reports |
| **Access** | Full data across all warehouses under management |

---

## 4. Functional Requirements

### 4.1 FR-001: Warehouse Keeper Assignment

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-001.1 | System shall allow assigning one or more users as Warehouse Keeper to specific warehouses | MUST |
| FR-001.2 | Assignment includes effective date range (from/to) | MUST |
| FR-001.3 | Only one active keeper per warehouse at a time (during any period) | SHOULD |
| FR-001.4 | Assignment history must be maintained (audit trail) | MUST |
| FR-001.5 | Keeper can only see warehouses they are assigned to | MUST |

### 4.2 FR-002: Slip Review & Recording (Ghi sổ)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-002.1 | System shall present pending receipt/issue slips to assigned keeper | MUST |
| FR-002.2 | Keeper can select one or multiple slips for recording | MUST |
| FR-002.3 | Recording creates a Stock Ledger entry with timestamp and keeper ID | MUST |
| FR-002.4 | Recorded slips show "Đã ghi sổ" status with keeper name and date | MUST |
| FR-002.5 | Keeper can un-record (Bỏ ghi) slips if not yet reconciled | MUST |
| FR-002.6 | Un-recording requires reason and creates audit trail | MUST |
| FR-002.7 | Cost prices can be hidden from keeper view (configurable per company) | SHOULD |

### 4.3 FR-003: Stock Ledger (Sổ kho)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-003.1 | System shall maintain a Stock Ledger per warehouse per item | MUST |
| FR-003.2 | Ledger shows: date, voucher type, voucher no, description, receipt qty, issue qty, running balance | MUST |
| FR-003.3 | Ledger entries are read-only after recording | MUST |
| FR-003.4 | Ledger supports printing in standard format (S06-DN or S07-DN) | MUST |
| FR-003.5 | Ledger can be filtered by date range, item, warehouse | MUST |

### 4.4 FR-004: Stock Card (Thẻ kho)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-004.1 | System shall generate Stock Card per item per warehouse per period | MUST |
| FR-004.2 | Card shows: opening balance, each receipt, each issue, closing balance | MUST |
| FR-004.3 | Card supports per-lot and per-serial tracking (if configured) | SHOULD |
| FR-004.4 | Card can be printed in A4 landscape/portrait format | MUST |

### 4.5 FR-005: Reconciliation Report

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-005.1 | System shall generate reconciliation between Keeper records and Accounting records | MUST |
| FR-005.2 | Report shows: item, warehouse, keeper qty, accounting qty, variance, variance value | MUST |
| FR-005.3 | Report supports date range filtering | MUST |
| FR-005.4 | Report can be exported to Excel | SHOULD |

### 4.6 FR-006: Physical Inventory Count

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-006.1 | System shall support inventory count with keeper participation | MUST |
| FR-006.2 | Count sheet shows book balance (can be hidden for blind count) | MUST |
| FR-006.3 | Keeper enters physical count quantities | MUST |
| FR-006.4 | System calculates variance automatically | MUST |
| FR-006.5 | Variance > tolerance triggers re-count | MUST |
| FR-006.6 | Count minutes (Biên bản kiểm kê) generated with keeper + accountant signatures | MUST |
| FR-006.7 | Approved count posts adjustment to GL and stock balance | MUST |

### 4.7 FR-007: Keeper Reports

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-007.1 | Inventory summary report per keeper warehouse | MUST |
| FR-007.2 | Receipt/Issue detail report per keeper | MUST |
| FR-007.3 | Reconciliation report (Đối chiếu nhập xuất kế toán - thủ kho) | MUST |
| FR-007.4 | Inventory count variance report | MUST |
| FR-007.5 | Warehouse turnover report per keeper | COULD |

---

## 5. Non-Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| NFR-001 | Response time < 200ms for slip list, < 500ms for reports | MUST |
| NFR-002 | Support 50+ concurrent keepers | SHOULD |
| NFR-003 | All keeper actions logged with timestamp, user, IP | MUST |
| NFR-004 | Module can be hidden/disabled via system configuration | MUST |
| NFR-005 | Compatible with existing warehouse module (no breaking changes) | MUST |
| NFR-006 | Works with both PostgreSQL and in-memory backends | MUST |

---

## 6. Business Rules

| ID | Rule | Source |
|----|------|--------|
| BR-001 | Keeper cannot create/edit/delete source documents (PNK/PXK) | Luật Kế toán 2015 Art. 6 |
| BR-002 | Keeper cannot view cost prices unless authorized | Internal control principle |
| BR-003 | One active keeper assignment per warehouse per period | MISA benchmark |
| BR-004 | Recorded slips cannot be modified without un-recording first | Circular 99 Art. 10 |
| BR-005 | Un-recording requires reason and creates audit trail | Internal control |
| BR-006 | Reconciliation report must show all discrepancies | Audit requirement |
| BR-007 | Physical count must have minimum 2 participants (keeper + accountant/manager) | Law on Accounting Art. 6 |
| BR-008 | Count variance > 5% of item value requires re-count | Industry standard |
| BR-009 | Keeper assignment changes must have overlap period (no gap) | Continuity requirement |
| BR-010 | Stock Ledger entries are immutable after period close | Circular 99 |

---

## 7. Assumptions

1. Existing warehouse module (CRUD, GRN, DN, transfers, adjustments, takes) is PROD-ready
2. Users and roles already exist in the auth module
3. Company-scoped data isolation already works
4. GL posting infrastructure is in place (service creates JournalEntries)
5. The module will be toggleable (can be hidden if company doesn't separate roles)

## 8. Dependencies

| Dependency | Type | Impact |
|------------|------|--------|
| Warehouse module | Internal | Uses existing stock balance, inventory transactions, stock take |
| Auth/User module | Internal | Needs user assignment, role checking |
| GL module | Internal | Stock ledger entries may need GL posting |
| Purchase module | Internal | PNK created by purchase flow |
| Sale module | Internal | PXK created by sale flow |
| Frontend (Alpine.js) | Internal | New pages needed |

## 9. Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Role conflict with existing warehouse permissions | Medium | High | Design as overlay on existing RBAC, not replacement |
| Performance impact from dual record-keeping | Low | Medium | Read-only ledger queries, indexed |
| User resistance to additional workflow step | High | Medium | Make module toggleable, provide clear training docs |
| Circular 99 interpretation changes | Low | High | Monitor MOF guidance, modular design for adaptation |

## 10. Success Criteria

1. Keeper can be assigned to warehouse with effective dates
2. Keeper sees only assigned warehouse slips
3. Keeper can record/un-record with audit trail
4. Stock Ledger matches MISA format (S06-DN)
5. Reconciliation report shows keeper vs accounting variances
6. Physical count supports blind count and re-count workflow
7. All tests pass: `go test -count=1 ./...`
8. Module can be toggled on/off without affecting existing warehouse functionality

## 11. Approvals

| Role | Name | Date | Status |
|------|------|------|--------|
| BA Lead | [TBD] | [TBD] | PENDING |
| Chief Accountant | [TBD] | [TBD] | PENDING |
| Tech Lead | [TBD] | [TBD] | PENDING |
