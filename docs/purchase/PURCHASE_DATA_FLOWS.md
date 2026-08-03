# Purchase Module — Data Flows & Entity Relationships

**Version:** 2.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Note:** Data flows for CRUD operations implemented. GL posting and inventory posting flows pending.

---

## 1. Entity Relationship Diagram (Text)

```
supplier
  │
  ├──< purchase_order (1:N)
  │       │
  │       └──< po_line (1:N)
  │             │
  │             ├──< goods_receipt_line (1:N)
  │             │       │
  │             │       └──< goods_receipt_note (N:1)
  │             │
  │             └──< supplier_invoice_line (1:N)
  │                       │
  │                       └──< supplier_invoice (N:1)
  │
  ├──< goods_receipt_note (1:N, via PO)
  │
  ├──< supplier_invoice (1:N)
  │
  ├──< ap_transaction (1:N)
  │
  └──< purchase_cost_allocation (1:N)
```

---

## 2. Core Data Flow: Procure-to-Pay

```
                    ┌─────────────────┐
                    │    Company      │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
     ┌────────▼────────┐          ┌─────────▼────────┐
     │    Supplier      │          │   GL Account     │
     │  (Supplier)      │          │  (Chart of Accts)│
     └────────┬────────┘          └─────────┬────────┘
              │                             │
              ▼                             │
     ┌────────────────┐                     │
     │ Purchase Order │◄────────────────────┘
     │ code, date,    │      (accounts: 152, 156,
     │ supplier, total│       331, 1331 per line)
     └───────┬────────┘
              │
              ▼
     ┌────────────────┐
     │  Goods Receipt  │
     │ code, date, qty,│──→ Inventory (152/156)
     │ PO reference    │──→ AP temp (331)
     └───────┬────────┘
              │
              ▼
     ┌────────────────┐
     │ Supplier Invoice│
     │ #, date, amount,│──→ 3-way match (PO × GRN)
     │ VAT, e-invoice  │──→ AP permanent (331)
     │ XML             │──→ VAT input (1331)
     └───────┬────────┘
              │
              ▼
     ┌────────────────┐
     │   AP Ledger    │
     │ by supplier,   │
     │ by invoice,    │
     │ aging buckets  │
     └───────┬────────┘
              │
              ▼
     ┌────────────────┐
     │   Payment      │──→ Bank (112)
     │ to supplier    │──→ AP decrease (331)
     └────────────────┘
```

---

## 3. GL Integration Data Flow

```
                         ┌─────────────────────┐
                         │    Journal Entry     │
                         │    (GL Service)      │
                         └──────────┬──────────┘
                                    │
         ┌──────────────────────────┼──────────────────────────┐
         ▼                          ▼                          ▼
┌─────────────────┐    ┌─────────────────────┐    ┌─────────────────┐
│    GRN Post     │    │   Invoice Post      │    │   Payment Post  │
│                 │    │                     │    │                 │
│ Dr 152/156      │    │ Dr 331 (temp rev)  │    │ Dr 331 (AP)     │
│ Cr 331 (temp)   │    │ Dr 1331            │    │ Cr 112          │
│                 │    │ Cr 331 (AP)        │    │                 │
│                 │    │                     │    │                 │
└─────────────────┘    └─────────────────────┘    └─────────────────┘
```

---

## 4. E-Invoice Data Flow

```
Supplier System                GDT Portal               GoTax System
     │                            │                         │
     │── XML Invoice ──→         │                         │
     │                            │── Validate signature    │
     │                            │── Assign invoice code   │
     │                            │── Push to buyer ──────→ │
     │                            │                         │── Parse XML
     │                            │                         │── Extract:
     │                            │                         │   - supplier_tax_code
     │                            │                         │   - invoice_number
     │                            │                         │   - date
     │                            │                         │   - items[]
     │                            │                         │   - totals
     │                            │                         │   - signature
     │                            │                         │
     │                            │                         │── Match supplier
     │                            │                         │── Match PO (by item)
     │                            │                         │── Store raw XML
     │                            │                         │── Create draft invoice
     │                            │                         │── Notify AP clerk
     │                            │                         │
     │                            │                         │── [AP Clerk reviews]
     │                            │                         │── Verify + Post
     │                            │                         │── Update GL
```

---

## 5. AP Aging Calculation Flow

```
Supplier Invoice Table         Payment Table           AP Aging Engine
     │                            │                         │
     │── Invoice amount ──────────┤                         │
     │── Amount paid    ←─────────┤── Payment amount        │
     │── Due date                 │                         │
     │                            │                         │── For each invoice:
     │                            │                         │   balance = amount - sum(payments)
     │                            │                         │   days_overdue = today - due_date
     │                            │                         │   bucket = f(days_overdue):
     │                            │                         │     0: "Current"
     │                            │                         │     1-30: "1-30 days"
     │                            │                         │     31-60: "31-60 days"
     │                            │                         │     61-90: "61-90 days"
     │                            │                         │     91-120: "91-120 days"
     │                            │                         │     120+: "120+ days"
     │                            │                         │
     │                            │                         │── Group by supplier
     │                            │                         │── Return aging report
```

---

## 6. Integration with Existing GoTax Modules

### 6.1 Bank Module Integration

```
Purchase Payment:
  POST /api/v1/purchase/invoices/:id/pay
    → Create APTransaction (type=payment)
    → Call Bank module: /api/v1/bank/transactions
    → Debit Account 331, Credit Account 112
    → Update invoice amount_paid, balance_due
    → If balance_due = 0 → invoice status = paid
```

### 6.2 Company Module Integration

```
Supplier is per-company (multi-tenant):
  supplier.company_id = company_id from auth middleware
  All purchase queries filtered by company context
```

### 6.3 Auth Module Integration

```
Role-based access (see PURCHASE_RULES.md R15):
  Auth middleware extracts user_id, role from JWT
  Purchase handler checks role for each operation
  Audit log: every mutation logged with user_id
```

### 6.4 GL Module Integration

```
Purchase → GL posting via existing service.PostJournalEntry():
  GRN:  Dr 152/156     Cr 331 (temp)
  INV:  Dr 331 (temp)  Cr 331 (AP)
        Dr 1331         Cr 331 (AP)
  PAY:  Dr 331 (AP)    Cr 112
```

---

## 7. Data Volume Estimates (SME)

| Entity | Monthly Volume | Annual Volume | Retention |
|--------|---------------|---------------|-----------|
| Suppliers | 10-50 new | 100-500 | Active |
| POs | 100-500 | 1,200-6,000 | 10 years |
| GRNs | 100-500 | 1,200-6,000 | 10 years |
| Supplier invoices | 100-500 | 1,200-6,000 | 10 years |
| AP transactions | 200-1,000 | 2,400-12,000 | 10 years |
| E-invoice XML storage | 100-500 | 1,200-6,000 | 10 years |

**Estimate at 5 years:** ~60K invoices, ~120K AP transactions, ~100MB XML storage.
**Scale target:** 10x for mid-market (600K invoices).