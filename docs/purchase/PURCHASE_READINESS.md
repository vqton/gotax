# Purchase Module — PROD Readiness Assessment

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP (amended by 70/2025), Decree 254/2026/ND-CP, IAS 2, IFRS 15

---

## VERDICT: NOT PRODUCTION READY — ~0% COMPLETE

| Dimension | Score | Detail |
|-----------|-------|--------|
| Code implementation | **0%** | No purchase models, repos, services, handlers exist |
| Database schema | **0%** | No purchase/AP tables in migrations |
| API endpoints | **0%** | No purchase routes registered |
| Business logic | **0%** | No P2P workflow, no AP aging, no tax calc |
| E-invoice integration | **0%** | No Decree 123/254 e-invoice receipt handling |
| Tax compliance | **0%** | No VAT deduction tracking, no import duty calc |
| Inventory integration | **0%** | No 152/153/156/151 account posting |
| AP management | **0%** | No 331 tracking, no payment scheduling |
| Multi-currency | **0%** | No FX revaluation for AP |
| Reporting | **0%** | No S01-DN, S02-DN, S03-DN, AP aging reports |

---

## What EXISTS in current GoTax that purchase reuses

| Component | File | Reuse For |
|-----------|------|-----------|
| `Company` domain model | `models_company.go` | Supplier table (extends Company or separate) |
| `ExchangeRate` | `models.go:295` | FX conversion for import purchases |
| `JournalEntry` posting engine | `service.go` | Auto-posting purchase → GL |
| `Period` management | `models.go:114` | Period-locking for purchase docs |
| Auth middleware | `auth.go` | Role-based access (AP clerk, AP manager, chief accountant) |
| Company multi-tenant | `company_handler.go` | Supplier-per-company isolation |
| PG pool + migrations | `db/pg.go` | Purchase tables migration |
| Memory repos pattern | `memory*.go` | Purchase memory repos for tests |
| Swagger annotations | `handler.go` | Purchase endpoint API docs |
| Bank module (112) | `models_bank.go` | Payment to suppliers via bank |

---

## Gap Analysis: MISA/Fast/Bravo vs GoTax

| Capability | MISA AMIS | Fast Accounting | Bravo ERP | GoTax (Current) | GoTax (Target) |
|------------|-----------|----------------|-----------|-----------------|----------------|
| Purchase requisition | Yes | Yes | Yes | No | P1 |
| Purchase order (PO) | Yes | Yes | Yes | No | P0 |
| Goods receipt note | Yes | Yes | Yes | No | P0 |
| Supplier invoice recording | Yes | Yes | Yes | No | P0 |
| Return to supplier | Yes | Yes | Yes | No | P1 |
| AP aging by invoice | Yes | Yes | Yes | No | P0 |
| AP aging by due date | Yes | Yes | Yes | No | P0 |
| Payment allocation | Yes | Yes | Yes | No | P0 |
| Prepayment tracking | Yes | Yes | Yes | No | P0 |
| Domestic purchase (VAT) | Yes | Yes | Yes | No | P0 |
| Import purchase (duty) | Yes | Yes | Yes | No | P1 |
| Purchase cost allocation | Yes | Yes | Yes | No | P1 |
| Supplier evaluation | No | Limited | Yes | No | P2 |
| E-invoice receipt (GDT) | Yes | Yes | Via API | No | P0 |
| 3-way matching | Yes | Yes | Yes | No | P0 |
| Multi-currency AP | Yes | Yes | Yes | No | P1 |
| Purchase budgeting | Yes | No | Yes | No | P2 |
| S01-DN (Purchase ledger) | Yes | Yes | Yes | No | P0 |
| S02-DN (Supplier detail) | Yes | Yes | Yes | No | P0 |
| S03-DN (Goods ledger) | Yes | Yes | Yes | No | P0 |

---

## Critical Path to PROD (Month-by-Month)

### Phase 1: Foundation (Month 1) — P0
- Database schema: suppliers, POs, goods receipt, supplier invoices, AP
- Domain models: Supplier, PurchaseOrder, GoodsReceipt, SupplierInvoice, APTransaction
- PG + memory repos for all entities
- CRUD endpoints for suppliers, POs, goods receipts, supplier invoices

### Phase 2: Core P2P (Month 2) — P0
- Full procure-to-pay flow: requisition → PO → receipt → invoice → payment
- 3-way matching engine (PO × receipt × invoice)
- AP aging + payment scheduling
- Auto-posting to GL (152/156/331/133)
- VAT deduction tracking per invoice

### Phase 3: Compliance (Month 3) — P0
- E-invoice receipt from GDT (Decree 254/2026)
- Import purchase with duty/VAT calc
- S01-DN, S02-DN, S03-DN reports
- Return to supplier workflow
- Prepayment tracking

### Phase 4: Enhancement (Month 4+) — P1/P2
- Purchase requisition approval workflow
- Supplier evaluation
- Purchase budgeting
- Multi-currency AP with auto revaluation
- Purchase cost allocation
- Supplier portal integration

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-----------|--------|------------|
| Decree 254/2026 e-invoice format changes | Medium | High | Build adapter layer, abstract GDT API |
| TT99 account changes for purchase | Low | Medium | Already on TT99 COA |
| Inventory valuation method changes | Low | Medium | Support FIFO, Weighted Avg, Specific ID, Std Cost |
| Multi-currency AP revaluation | Medium | Medium | Use existing ExchangeRate framework |
| Performance at scale (100K+ POs) | Low | Medium | PG indexing, pagination from day 1 |

---

## References
- MISA AMIS Mua hang: https://amis.misa.vn/amis-mua-hang/
- Fast Accounting Purchase: https://fast.com.vn/phan-mem-ke-toan-fast-accounting-phan-he-ke-toan-mua-hang-va-cong-no-phai-tra/
- Bravo ERP Purchase: https://www.bravo.com.vn/san-pham/phan-he-co-ban/quan-ly-mua-hang
- Tryton Purchase: https://docs.tryton.org/8.0/modules-purchase/index.html
- Circular 99/2025/TT-BTC: https://congbao.chinhphu.vn/van-ban/thong-tu-so-99-2025-tt-btc-46529.htm
- Decree 254/2026/ND-CP: E-invoice rules
- IAS 2 Inventories: https://www.ifrs.org/issued-standards/list-of-standards/ias-2-inventories/