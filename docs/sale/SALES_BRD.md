# Sale Module — Business Requirements Document (BRD)

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP (amended by 70/2025), Decree 254/2026/ND-CP, IFRS 15 (Revenue from Contracts with Customers), VAS 14 (Revenue), VAS 08 (Sales Return), Circular 200/2014/TT-BTC

---

## 1. Executive Summary

Sale module manages the complete Order-to-Cash (O2C) cycle: from customer inquiry through sales order, delivery, invoicing, AR tracking, and customer payment collection. Vietnamese tax compliance requires strict e-invoice issuance (Decree 123 + Decree 254), VAT output tracking, and revenue recognition per Circular 99 accounts (511, 131, 3331, 521, 632).

**GoTax currently has ZERO sale functionality.** This BRD defines requirements for building P0 (must-have) features for PROD launch.

---

## 2. Business Objectives

| # | Objective | Success Metric | Priority |
|---|-----------|----------------|----------|
| OBJ-1 | Record all sales transactions and track AR balances | Zero untracked customer invoices | P0 |
| OBJ-2 | Automate e-invoice issuance (TXML → sign → GDT submit) | >95% e-invoice auto-issuance rate | P0 |
| OBJ-3 | Generate AR aging reports for cash flow management | Aging by 30/60/90/120+ days | P0 |
| OBJ-4 | Auto-post sales to GL accounts (511/131/3331/632) | Zero manual GL entries for standard sales | P0 |
| OBJ-5 | Track VAT output per invoice | Full VAT output trail for tax return | P0 |
| OBJ-6 | Support sales returns with credit note workflow | Accurate revenue reversal and AR adjustment | P0 |
| OBJ-7 | Track customer prepayments and deposits | Link prepayments to invoices automatically | P0 |
| OBJ-8 | Generate S01-BH, S02-BH, S03-BH regulatory reports | Circular 99 compliant sales ledgers | P0 |
| OBJ-9 | Issue e-invoices to GDT portal (Decree 254/2026) | Fully compliant e-invoice pipeline | P0 |
| OBJ-10 | Customer credit limit management | Block shipments/orders exceeding limit | P1 |
| OBJ-11 | Sales quotation and conversion tracking | Quote-to-order conversion rate | P2 |
| OBJ-12 | Sales commission tracking | Accurate commission calculation per invoice | P2 |

---

## 3. Stakeholders

| Stakeholder | Role | Key Concern |
|-------------|------|-------------|
| Chief Accountant | Oversees revenue accounting | Accuracy, compliance, audit trail |
| AR Accountant | Processes customer invoices and collections | Efficiency, e-invoice, AR aging |
| Sales Manager | Manages customer relationships | Order tracking, credit limits |
| Warehouse Keeper | Delivers goods | Delivery accuracy, stock availability |
| CFO | Cash flow management | AR aging, collection forecast, DSO |
| External Auditor | Verifies revenue transactions | Audit trail, revenue recognition timing |
| Tax Authority (GDT) | Validates VAT output | E-invoice authenticity, timely issuance |
| Customer | Receives goods and invoices | Accurate billing, e-invoice receipt |

---

## 4. Functional Requirements

### FR-1: Customer Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-1.1 | Create customer with tax code, name, address, phone, bank account | P0 |
| FR-1.2 | Classify customer: domestic/export, goods/services, one-time/regular | P0 |
| FR-1.3 | Track customer contract: contract#, start/end date, terms, value | P1 |
| FR-1.4 | Customer credit limit and payment terms (net 15/30/45/60, COD) | P0 |
| FR-1.5 | Customer status: active, suspended, blacklisted | P0 |
| FR-1.6 | Customer group/segment classification (retail/wholesale/distributor) | P1 |
| FR-1.7 | Customer bank account for refunds | P0 |
| FR-1.8 | Customer price list / tiered pricing | P1 |

### FR-2: Sales Quotation (Optional, P1)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-2.1 | Create quotation with items, qty, unit price, discount, VAT | P1 |
| FR-2.2 | Quotation approval if above threshold | P1 |
| FR-2.3 | Convert approved quotation to sales order | P1 |
| FR-2.4 | Quotation expiry and version tracking | P1 |
| FR-2.5 | Quote-to-order conversion rate reporting | P2 |

### FR-3: Sales Order

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-3.1 | Create sales order from quotation or directly | P0 |
| FR-3.2 | Sales order with customer, items, qty, unit price, discount, VAT, delivery date | P0 |
| FR-3.3 | Order terms: payment terms, delivery terms (FOB/CIF/DDP), currency, FX rate | P0 |
| FR-3.4 | Sales order approval workflow (value-based: <10M auto, <100M manager, >100M director) | P0 |
| FR-3.5 | SO status: draft, approved, confirmed, processing, delivered, invoiced, cancelled, closed | P0 |
| FR-3.6 | Partial delivery allowed against SO | P0 |
| FR-3.7 | SO amendment with version tracking | P1 |
| FR-3.8 | Credit check against customer limit before approval | P1 |
| FR-3.9 | Inventory availability check before confirmation | P1 |

### FR-4: Delivery / Goods Delivery

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-4.1 | Create delivery note referencing SO | P0 |
| FR-4.2 | Record actual quantity delivered vs ordered | P0 |
| FR-4.3 | Over-delivery control: % tolerance (configurable, default 5%) | P0 |
| FR-4.4 | Under-delivery tracking for back-order fulfillment | P0 |
| FR-4.5 | Direct delivery (service/non-inventory items) | P0 |
| FR-4.6 | Delivery address, shipping method, carrier tracking | P1 |
| FR-4.7 | Auto-post delivery to COGS (632) and reduce inventory (152/156) | P0 |
| FR-4.8 | e-Delivery note (Packing List, Delivery Order) generation | P1 |

### FR-5: Customer Invoice

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-5.1 | Create invoice referencing SO and/or delivery note | P0 |
| FR-5.2 | Capture: invoice number (auto), date, customer, items, qty, price, VAT | P0 |
| FR-5.3 | Generate e-invoice XML per Decree 254 format | P0 |
| FR-5.4 | Digitally sign e-invoice before GDT submission | P0 |
| FR-5.5 | Submit e-invoice to GDT portal | P0 |
| FR-5.6 | Receive GDT invoice code and update status | P0 |
| FR-5.7 | Invoice status: draft, signed, submitted, coded, issued, cancelled, replaced | P0 |
| FR-5.8 | Support corrective/adjustment invoice (increase/decrease) | P0 |
| FR-5.9 | Support consolidated invoice (multiple deliveries in one invoice) | P1 |
| FR-5.10 | Auto-reverse accrual if revenue recognized before invoicing | P1 |

### FR-6: Accounts Receivable

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-6.1 | AR aging: current, 1-30, 31-60, 61-90, 91-120, 120+ days by invoice | P0 |
| FR-6.2 | AR by customer: total, due, overdue | P0 |
| FR-6.3 | Payment allocation to specific invoice(s) | P0 |
| FR-6.4 | Track customer prepayment/deposit | P0 |
| FR-6.5 | Offset prepayment against invoice on issuance | P0 |
| FR-6.6 | Collection scheduling based on due date | P0 |
| FR-6.7 | Credit balance (customer overpayment) tracking | P0 |
| FR-6.8 | Early payment discount tracking (customer discount) | P1 |
| FR-6.9 | Customer statement generation | P1 |
| FR-6.10 | DSO (Days Sales Outstanding) calculation | P2 |

### FR-7: Customer Receipt / Payment Collection

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-7.1 | Record customer payment (cash/bank/cheque) | P0 |
| FR-7.2 | Allocate payment to one or multiple invoices | P0 |
| FR-7.3 | Partial payment handling | P0 |
| FR-7.4 | Overpayment handling (credit balance) | P0 |
| FR-7.5 | Payment in foreign currency with FX gain/loss | P0 |
| FR-7.6 | Bank integration for auto-reconciliation (linked to Bank module) | P1 |
| FR-7.7 | Payment reminder / dunning letters | P2 |

### FR-8: Sales Return / Credit Note

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-8.1 | Record sales return with reason (defect/wrong item/excess/other) | P0 |
| FR-8.2 | Issue credit note (negative e-invoice) per Decree 254 | P0 |
| FR-8.3 | Auto-reverse revenue (511) and VAT output (3331) | P0 |
| FR-8.4 | Auto-adjust AR balance on credit note | P0 |
| FR-8.5 | Return goods to inventory (reverse COGS) | P0 |
| FR-8.6 | Return with/without replacement | P1 |

### FR-9: Reports

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-9.1 | S01-BH: Sales ledger (So chi tiet ban hang) — per customer per period | P0 |
| FR-9.2 | S02-BH: Customer detail ledger (So chi tiet cong no phai thu) | P0 |
| FR-9.3 | S03-BH: Goods sales ledger (So chi tiet hang hoa ban) | P0 |
| FR-9.4 | AR aging report (Bang phan tich tuoi no phai thu) | P0 |
| FR-9.5 | Sales summary by customer, item, period | P0 |
| FR-9.6 | VAT output tracking report (Bang ke hoa don VAT dau ra) | P0 |
| FR-9.7 | Revenue recognition schedule | P1 |
| FR-9.8 | Unbilled delivery report (hang da giao chua xuat hoa don) | P0 |
| FR-9.9 | Customer statement of account | P1 |

---

## 5. Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-1 | Sales order creation response < 2s at 100K orders | P0 |
| NFR-2 | E-invoice generation + signing < 5s per invoice | P0 |
| NFR-3 | AR aging report loads < 10s for 50K invoices | P0 |
| NFR-4 | Data retention: 10 years per Accounting Law | P0 |
| NFR-5 | Audit trail: every sales document state change logged | P0 |
| NFR-6 | E-invoice XML generated within 1s for immediate issuance | P0 |
| NFR-7 | Multi-tenant isolation: customer data per company | P0 |
| NFR-8 | Concurrent invoice issuance: support 100+ invoices/min | P1 |

---

## 6. Regulatory Compliance Matrix

| Regulation | Sale Requirement | Module Impact |
|------------|------------------|---------------|
| Circular 99/2025/TT-BTC Account 511 | Revenue recording by type (5111/5112/5113) | Auto-posting to revenue accounts |
| Circular 99/2025/TT-BTC Account 131 | AR tracking per customer | AR sub-ledger |
| Circular 99/2025/TT-BTC Account 3331 | VAT output tracking | VAT output per invoice |
| Circular 99/2025/TT-BTC Account 521 | Trade discounts/returns/allowances | Revenue adjustment accounts |
| Circular 99/2025/TT-BTC Account 632 | COGS recognition | Inventory → COGS on delivery |
| Circular 99/2025/TT-BTC Account 3387 | Unearned revenue | Deferred revenue tracking |
| Decree 123/2020 Art 48 | Mandatory e-invoice issuance within 24h | E-invoice pipeline |
| Decree 254/2026 | E-invoice XML format | TXML generator, signer, submitter |
| Decree 23/2025/ND-CP | Digital signature on e-invoice | Signature integration |
| Decree 70/2025/ND-CP | Corrective invoice handling | Adjustment invoice workflow |
| IFRS 15 | 5-step revenue recognition | Performance obligation tracking |
| VAS 14 | Revenue recognition timing | Revenue recognition rules engine |
| Law on VAT 13/2024/QH15 | VAT output declaration | VAT output tracking per rate |

---

## 7. Integration Points

| Module | Integration | Direction |
|--------|------------|-----------|
| Inventory (future) | SO → Delivery → Stock decrease | Sale → Inventory |
| Bank (existing) | Receipt → Payment allocation → Bank transaction | Sale → Bank |
| Tax (planned) | VAT output from invoices | Sale → Tax |
| GL (existing) | Auto-post all sales transactions | Sale → GL |
| Company (existing) | Customer table per company | Company → Sale |
| Cash (existing) | Customer payment via petty cash | Cash → Sale |
| Purchase (future) | Inter-company sales → Purchase | Sale → Purchase |

---

## 8. Assumptions & Constraints

- **ASSUMPTION-1:** GoTax will maintain its own customer master (not synced from external CRM)
- **ASSUMPTION-2:** E-invoice issuance to GDT requires building XML generator + signer + API client
- **ASSUMPTION-3:** Inventory module is built after/in parallel — Sale stores delivery qty in sale tables first
- **ASSUMPTION-4:** Digital signature for e-invoice uses existing DigitalSignature model in company module
- **ASSUMPTION-5:** Export sales follow separate workflow (non-VAT, customs declaration)
- **CONSTRAINT-1:** E-invoice must be issued within 24 hours of delivery (Decree 123 Art 48)
- **CONSTRAINT-2:** Invoice must be signed by registered digital certificate
- **CONSTRAINT-3:** Corrective invoices must reference original invoice number and GDT code
