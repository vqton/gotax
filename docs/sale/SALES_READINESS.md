# Sale Module — Readiness Assessment

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)

---

## 1. Current State

**VERDICT: Sale module ~0% complete. NOT production-ready.**

What exists: **NOTHING**. Zero code, zero schema, zero endpoints.

Existing related code that can be leveraged:
- `digital_signature` CRUD in Company module (for e-invoice signing)
- `bank_account` CRUD in Company module (for customer bank info)
- `EInvoiceType`, `EInvoiceID` enums in Tax module models
- `CounterpartCustomer`, `ReceiptSales`, `ReceiptCustomerPayment` in Cash module
- `DetailCustomer` in Opening Balance module (for AR opening balance)

---

## 2. Gap Analysis vs Enterprise ERPs

| Capability | MISA | FAST | Bravo ERP | GoTax (Current) | Priority |
|------------|------|------|-----------|-----------------|----------|
| Customer master | ✅ | ✅ | ✅ | ✅ | P0 |
| Customer classification | ✅ | ✅ | ✅ | ❌ | P0 |
| Sales quotation | ✅ | ✅ | ✅ | ✅ | P2 |
| Sales order | ✅ | ✅ | ✅ | ✅ | P0 |
| SO approval workflow | ✅ | ✅ | ✅ | ✅ | P0 |
| Credit limit check | ✅ | ✅ | ✅ | ✅ | P1 |
| Inventory availability | ✅ | ✅ | ✅ | ❌ | P1 |
| Delivery note | ✅ | ✅ | ✅ | ✅ | P0 |
| COGS auto-posting | ✅ | ✅ | ✅ | ✅ | P0 |
| E-invoice (TXML + GDT) | ✅ | ✅ | ✅ | ❌ | P0 |
| E-invoice signing (XMLDSig) | ✅ | ✅ | ✅ | ❌ | P0 |
| VAT output tracking | ✅ | ✅ | ✅ | ✅ | P0 |
| AR management | ✅ | ✅ | ✅ | ✅ | P0 |
| AR aging report | ✅ | ✅ | ✅ | ✅ | P0 |
| Customer receipt | ✅ | ✅ | ✅ | ✅ | P0 |
| Payment allocation | ✅ | ✅ | ✅ | ✅ | P0 |
| Credit note / sales return | ✅ | ✅ | ✅ | ✅ | P0 |
| Corrective invoice | ✅ | ✅ | ✅ | ❌ | P0 |
| Prepayment / deposit | ✅ | ✅ | ✅ | ❌ | P0 |
| Dunning / reminder | ✅ | ✅ | ✅ | ❌ | P2 |
| Customer statement | ✅ | ✅ | ✅ | ✅ | P1 |
| S01-BH / S02-BH / S03-BH | ✅ | ✅ | ✅ | ✅ | P0 |
| Sales by item / customer | ✅ | ✅ | ✅ | ✅ | P1 |
| Price list / tiered pricing | ✅ | ✅ | ✅ | ❌ | P1 |
| Sales commission | ✅ | ✅ | ✅ | ❌ | P2 |
| DSO calculation | ✅ | ✅ | ✅ | ❌ | P2 |
| Bank integration (auto-rec) | ✅ | ✅ | ✅ | ❌ | P1 |
| E-invoice QR code | ✅ | ✅ | ✅ | ❌ | P0 |
| Multi-currency AR | ✅ | ✅ | ✅ | ❌ | P0 |

---

## 3. Build Priority vs MISA/FAST/Bravo

| Phase | Features | Timeline | Comparable to |
|-------|----------|----------|---------------|
| P0 (MUST) | Customer, SO, Delivery, Invoice, AR, Receipt, Credit Note, E-invoice, Reports | Phase 1 | MISA SME basic |
| P1 (SHOULD) | Quotation, Price List, Credit Limit, Customer Statement, Bank Integration, Deposit | Phase 2 | FAST basic |
| P2 (COULD) | Commission, Dunning, DSO, Dashboard, Sales Forecast, CRM Lite | Phase 3 | Bravo / MISA Premium |

---

## 4. Key Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| GDT e-invoice API integration complex | High - blocks P0 | Start early, build mock GDT for testing |
| E-invoice XML format changes | Medium - rework | Abstract XML generator behind interface |
| Digital signature integration | High - security critical | Use existing DigitalSignature model; audit thoroughly |
| 24h e-invoice issuance SLA | Medium - operational | Monitoring dashboard; queue for retry on GDT failure |
| Multi-currency FX revaluation | Medium - accounting | Reuse FX logic from existing GL service |
| Inventory module not ready | Medium - delivery blocked | Store delivery qty in sale tables first; post COGS to GL directly |

---

## 5. Recommended Build Order

1. **Customer master** (CRUD + validation) — standalone
2. **Sales Order** (create, approve, confirm, cancel) — standalone
3. **Delivery Note** (create, post to GL, link to SO) — depends on SO
4. **Customer Invoice** (create, 2-way match, e-invoice pipeline) — depends on SO + DN
5. **E-invoice pipeline** (TXML gen → sign → GDT submit → code) — depends on Invoice
6. **AR** (aging, sub-ledger, reports) — depends on Invoice
7. **Customer Receipt** (record, allocate, post) — depends on Invoice
8. **Credit Note** (return, reversal, e-invoice) — depends on Invoice
9. **Reports** (S01-BH, S02-BH, S03-BH, VAT output, unbilled) — depends on all above
10. **P1/P2 features** — incremental after P0 stable

Total estimated effort: **8-12 weeks** for P0 with 1-2 developers + 1 BA/Accountant.
