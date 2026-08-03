# Purchase Module — Business Requirements Document (BRD)

**Version:** 2.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP (amended by 70/2025), Decree 254/2026/ND-CP, IAS 2, IFRS 15, VAS 02 (Inventory)

---

## 1. Executive Summary

Purchase module manages the complete Procure-to-Pay (P2P) cycle: from purchase requisition through supplier selection, PO issuance, goods receipt, invoice verification, AP tracking, and supplier payment. Vietnamese tax compliance requires strict VAT invoice tracking (Decree 123/2020 + TT99 accounts: 151, 152, 153, 156, 331, 133).

**GoTax Purchase module: ~70% complete.** Domain models, repos (PG + memory), service, handlers (28 routes), and migration are built. Remaining gaps: 3-way matching, GL auto-posting, service tests, e-invoice GDT integration. This BRD defines remaining requirements for PROD launch.

---

## 2. Business Objectives

| # | Objective | Success Metric | Priority |
|---|-----------|----------------|----------|
| OBJ-1 | Record all supplier invoices and track AP balances | Zero untracked supplier invoices | P0 |
| OBJ-2 | Automate 3-way matching (PO × Receipt × Invoice) | >95% match rate | P0 |
| OBJ-3 | Generate AP aging reports for cash flow planning | Aging by 30/60/90/120+ days | P0 |
| OBJ-4 | Auto-post purchase to GL accounts (152/156/331/133) | Zero manual GL entries | P0 |
| OBJ-5 | Track VAT input deduction per invoice | Full VAT deduction trail | P0 |
| OBJ-6 | Support import purchase with duty/VAT/special tax | Complete landed cost | P1 |
| OBJ-7 | Support prepayment to suppliers | Track advances vs invoices | P0 |
| OBJ-8 | Generate S01-DN, S02-DN, S03-DN regulatory reports | Circular 99 compliant | P0 |
| OBJ-9 | Receive e-invoices from GDT portal (Decree 254/2026) | Auto-receive supplier e-invoices | P0 |
| OBJ-10 | Supplier evaluation and performance tracking | Scorecard per supplier | P2 |

---

## 3. Stakeholders

| Stakeholder | Role | Key Concern |
|-------------|------|-------------|
| Chief Accountant | Oversees purchase accounting | Accuracy, compliance, audit trail |
| AP Accountant | Processes supplier invoices | Efficiency, 3-way match, payment scheduling |
| Purchasing Manager | Negotiates and places POs | PO tracking, budget control |
| Warehouse Keeper | Receives goods | Receipt accuracy, quality check |
| CFO | Cash flow management | AP aging, payment forecast |
| External Auditor | Verifies purchase transactions | Audit trail, VAT deduction evidence |
| Tax Authority (GDT) | Validates VAT input | E-invoice authenticity, deduction eligibility |
| Customs (if import) | Validates import declarations | Duty, import VAT, HS code |

---

## 4. Functional Requirements

### FR-1: Supplier Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-1.1 | Create supplier with tax code, name, address, phone, bank account | P0 |
| FR-1.2 | Classify supplier: domestic/import, goods/services, one-time/regular | P0 |
| FR-1.3 | Track supplier contract: contract#, start/end date, terms, value | P1 |
| FR-1.4 | Supplier credit limit and payment terms (net 30/60/90, COD) | P0 |
| FR-1.5 | Supplier status: active, suspended, blacklisted | P0 |
| FR-1.6 | Supplier group/segment classification | P2 |
| FR-1.7 | Supplier bank account for payment (linked to Bank module) | P0 |

### FR-2: Purchase Requisition

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-2.1 | Create requisition with items, quantity, estimated price, need-by date | P1 |
| FR-2.2 | Approval workflow: requester → department head → purchasing | P1 |
| FR-2.3 | Convert approved requisition to PO | P1 |
| FR-2.4 | Track requisition status: draft, pending, approved, rejected, ordered | P1 |

### FR-3: Purchase Order

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-3.1 | Create PO from requisition or directly | P0 |
| FR-3.2 | PO with supplier, items, quantity, unit price, VAT rate, delivery date | P0 |
| FR-3.3 | PO terms: payment terms, delivery terms (FOB/CIF/DDP), currency | P0 |
| FR-3.4 | PO approval workflow (value-based: <10M auto, <100M manager, >100M director) | P0 |
| FR-3.5 | PO status: draft, approved, sent, partial, received, cancelled, closed | P0 |
| FR-3.6 | Partial receipt allowed against PO | P0 |
| FR-3.7 | PO amendment/change order with version tracking | P1 |
| FR-3.8 | Budget check against PO (if purchase budgeting enabled) | P2 |

### FR-4: Goods Receipt

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-4.1 | Create goods receipt note (GRN) referencing PO | P0 |
| FR-4.2 | Record actual quantity received, quality ok/reject | P0 |
| FR-4.3 | Link to warehouse/inventory location | P0 |
| FR-4.4 | Over-receipt control: % tolerance (configurable, default 5%) | P0 |
| FR-4.5 | Under-receipt tracking for back-order fulfillment | P0 |
| FR-4.6 | Direct receipt (non-inventory items: services, expenses) | P0 |
| FR-4.7 | Quality inspection hold, reject, return to supplier | P1 |
| FR-4.8 | Auto-post GRN to inventory (152/153/156) with temporary price if uninvoiced | P0 |

### FR-5: Supplier Invoice

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-5.1 | Record supplier invoice referencing PO and/or GRN | P0 |
| FR-5.2 | Capture invoice: number, date, amount, VAT amount, VAT rate | P0 |
| FR-5.3 | Support e-invoice receipt from GDT (XML format per Decree 254) | P0 |
| FR-5.4 | 3-way matching: PO qty vs GRN qty vs invoice qty | P0 |
| FR-5.5 | Handle invoice with no PO (direct invoice for services) | P0 |
| FR-5.6 | Track VAT input deduction status: pending, claimed, rejected | P0 |
| FR-5.7 | Handle corrective/adjustment invoice from supplier | P0 |
| FR-5.8 | Handle credit note (return) from supplier | P0 |
| FR-5.9 | Invoice status: draft, verified, posted, paid, cancelled | P0 |

### FR-6: Accounts Payable

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-6.1 | AP aging: current, 1-30, 31-60, 61-90, 91-120, 120+ days by invoice | P0 |
| FR-6.2 | AP by supplier: total, due, overdue | P0 |
| FR-6.3 | Payment allocation to specific invoice(s) | P0 |
| FR-6.4 | Track prepayment/advance to supplier | P0 |
| FR-6.5 | Offset prepayment against invoice on receipt | P0 |
| FR-6.6 | Payment scheduling based on due date + payment terms | P0 |
| FR-6.7 | Debit balance (advance) tracking per supplier | P0 |
| FR-6.8 | Discount capture (early payment discount to 515 account) | P1 |
| FR-6.9 | Supplier statement reconciliation | P1 |

### FR-7: Purchase Returns

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-7.1 | Return goods to supplier with reason | P1 |
| FR-7.2 | Receive credit note from supplier | P1 |
| FR-7.3 | Auto-adjust AP balance on return | P1 |
| FR-7.4 | Return with/without replacement | P1 |

### FR-8: Reports

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-8.1 | S01-DN: Purchase ledger (So chi tiet mua hang) — per supplier per period | P0 |
| FR-8.2 | S02-DN: Supplier detail ledger (So chi tiet cong no phai tra) | P0 |
| FR-8.3 | S03-DN: Goods purchase ledger (So chi tiet hang hoa) | P0 |
| FR-8.4 | AP aging report (Bang phan tich tuoi no phai tra) | P0 |
| FR-8.5 | Purchase summary by supplier, item, period | P0 |
| FR-8.6 | VAT input tracking report (Bang ke hoa don VAT dau vao) | P0 |
| FR-8.7 | Uninvoiced receipt report (hang ve chua co hoa don) | P0 |
| FR-8.8 | Purchase budget vs actual | P2 |

---

## 5. Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-1 | PO creation response < 2s at 100K POs | P0 |
| NFR-2 | 3-way matching completes < 5s per invoice | P0 |
| NFR-3 | AP aging report loads < 10s for 50K invoices | P0 |
| NFR-4 | Data retention: 10 years per Accounting Law | P0 |
| NFR-5 | Audit trail: every purchase doc state change logged | P0 |
| NFR-6 | E-invoice XML receives within 30s of GDT push | P0 |
| NFR-7 | Multi-tenant isolation: supplier data per company | P0 |

---

## 6. Regulatory Compliance Matrix

| Regulation | Purchase Requirement | Module Impact |
|------------|---------------------|---------------|
| Circular 99/2025/TT-BTC Account 152 | Raw material purchase recording | Auto-posting to 152 |
| Circular 99/2025/TT-BTC Account 156 | Goods purchase recording | Auto-posting to 156 |
| Circular 99/2025/TT-BTC Account 331 | AP tracking per supplier | AP sub-ledger |
| Circular 99/2025/TT-BTC Account 133 | VAT deduction tracking | VAT input tracking |
| Circular 99/2025/TT-BTC Account 151 | Goods in transit | Transit tracking |
| Decree 123/2020 Art 56 | Buyer must request invoice | E-invoice receive workflow |
| Decree 254/2026 | GDT e-invoice XML format | XML parser, validation |
| Decree 70/2025/ND-CP | Return invoice handling | Credit note workflow |
| IAS 2 | Inventory cost = purchase + conversion + other | Cost allocation |
| VAS 02 | Inventory at lower of cost or NRV | NRV check on receipt |
| Decree 23/2025/ND-CP | Digital signature on e-invoice | Signature verification |

---

## 7. Integration Points

| Module | Integration | Direction |
|--------|------------|-----------|
| Inventory (future) | PO → Goods Receipt → Stock increase | Purchase → Inventory |
| Bank (existing) | AP → Payment → Bank transaction | Purchase → Bank |
| Tax (planned) | VAT input from supplier invoices | Purchase → Tax |
| GL (existing) | Auto-post all purchase transactions | Purchase → GL |
| Company (existing) | Supplier table per company | Company → Purchase |
| Cash (future) | Payment from petty cash to supplier | Cash → Purchase |

---

## 8. Assumptions & Constraints

- **ASSUMPTION-1:** GoTax will maintain its own supplier master (not synced from external CRM)
- **ASSUMPTION-2:** E-invoice receipt from GDT requires building XML parser + API client
- **ASSUMPTION-3:** Inventory module is built after/in parallel — Purchase stores GRN qty in purchase tables first
- **ASSUMPTION-4:** Currency = VND by default; multi-currency is P1
- **CONSTRAINT-1:** Must support Circular 99 COA accounts from day 1
- **CONSTRAINT-2:** Must handle Decree 254/2026 e-invoice XML format
- **CONSTRAINT-3:** Must integrate with existing auth/company/GL layers