# Purchase Module — Functional Specifications

**Version:** 2.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Note:** Models implemented in `internal/domain/models_purchase.go`. See code for latest field definitions.

---

## 1. Data Model

### 1.1 Supplier

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK → company |
| code | VARCHAR(20) | yes | Unique per company |
| name | VARCHAR(255) | yes | |
| tax_code | VARCHAR(20) | yes | MST, validated format |
| address | TEXT | no | |
| phone | VARCHAR(20) | no | |
| email | VARCHAR(100) | no | |
| bank_account_name | VARCHAR(255) | no | For payment |
| bank_account_number | VARCHAR(50) | no | |
| bank_name | VARCHAR(255) | no | |
| payment_terms | VARCHAR(50) | no | net30, net60, net90, COD |
| credit_limit | DECIMAL(18,2) | no | 0 = no limit |
| currency | VARCHAR(3) | no | Default VND |
| supplier_type | VARCHAR(20) | no | domestic, import, both |
| status | VARCHAR(20) | yes | active, suspended, blacklisted |
| notes | TEXT | no | |
| created_at | TIMESTAMP | auto | |
| updated_at | TIMESTAMP | auto | |

### 1.2 Purchase Order

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK → company |
| po_number | VARCHAR(30) | yes | Auto-generated: PO-YYYYMM-XXXX |
| supplier_id | UUID | yes | FK → supplier |
| requisition_id | UUID | no | FK → purchase_requisition |
| order_date | DATE | yes | |
| expected_date | DATE | no | |
| currency | VARCHAR(3) | yes | Default VND |
| exchange_rate | DECIMAL(18,6) | no | If foreign currency |
| payment_terms | VARCHAR(50) | no | |
| delivery_terms | VARCHAR(50) | no | FOB, CIF, DDP, EXW |
| subtotal | DECIMAL(18,2) | auto | Sum of line amounts |
| discount_amount | DECIMAL(18,2) | no | |
| tax_amount | DECIMAL(18,2) | auto | |
| total_amount | DECIMAL(18,2) | auto | |
| status | VARCHAR(20) | yes | draft, approved, sent, partial, received, cancelled, closed |
| approved_by | UUID | no | FK → user |
| approved_at | TIMESTAMP | no | |
| cancelled_reason | TEXT | no | |
| notes | TEXT | no | |
| created_by | UUID | yes | FK → user |
| created_at | TIMESTAMP | auto | |
| updated_at | TIMESTAMP | auto | |

### 1.3 Purchase Order Line

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| po_id | UUID | yes | FK → purchase_order |
| line_number | INT | yes | |
| item_code | VARCHAR(50) | no | |
| item_name | VARCHAR(255) | yes | |
| unit | VARCHAR(20) | yes | kg, pcs, box, m |
| quantity | DECIMAL(18,4) | yes | |
| unit_price | DECIMAL(18,4) | yes | |
| discount_pct | DECIMAL(5,2) | no | |
| vat_rate | DECIMAL(5,2) | yes | 0, 5, 8, 10 |
| vat_type | VARCHAR(10) | yes | VAT_0, VAT_5, VAT_8, VAT_10, NT (non-taxable) |
| account_id | VARCHAR(20) | yes | GL account (152/156/153/642...) |
| vat_account_id | VARCHAR(20) | yes | GL account 1331 |
| line_total | DECIMAL(18,2) | auto | |
| line_vat_amount | DECIMAL(18,2) | auto | |
| received_qty | DECIMAL(18,4) | auto | Cumulative from GRN |
| invoiced_qty | DECIMAL(18,4) | auto | Cumulative from invoices |

### 1.4 Goods Receipt Note

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK → company |
| grn_number | VARCHAR(30) | yes | Auto: GRN-YYYYMM-XXXX |
| po_id | UUID | yes | FK → purchase_order |
| receipt_date | DATE | yes | |
| warehouse | VARCHAR(50) | no | Warehouse code |
| status | VARCHAR(20) | yes | draft, posted, cancelled |
| notes | TEXT | no | |
| created_by | UUID | yes | |
| created_at | TIMESTAMP | auto | |

### 1.5 Goods Receipt Line

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| grn_id | UUID | yes | FK |
| po_line_id | UUID | yes | FK → po_line |
| item_code | VARCHAR(50) | no | |
| item_name | VARCHAR(255) | yes | |
| unit | VARCHAR(20) | yes | |
| quantity_received | DECIMAL(18,4) | yes | |
| quantity_rejected | DECIMAL(18,4) | no | Quality reject |
| unit_price | DECIMAL(18,4) | auto | From PO |
| line_total | DECIMAL(18,2) | auto | |

### 1.6 Supplier Invoice

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK |
| invoice_number | VARCHAR(50) | yes | Supplier's invoice # |
| invoice_date | DATE | yes | |
| po_id | UUID | no | FK |
| grn_id | UUID | no | FK |
| supplier_id | UUID | yes | FK |
| supplier_name | VARCHAR(255) | yes | Denormalized |
| supplier_tax_code | VARCHAR(20) | yes | Denormalized |
| invoice_type | VARCHAR(20) | yes | domestic, import, service |
| currency | VARCHAR(3) | yes | |
| exchange_rate | DECIMAL(18,6) | no | |
| subtotal | DECIMAL(18,2) | auto | |
| discount_amount | DECIMAL(18,2) | no | |
| tax_amount | DECIMAL(18,2) | auto | |
| total_amount | DECIMAL(18,2) | auto | |
| amount_paid | DECIMAL(18,2) | no | |
| balance_due | DECIMAL(18,2) | auto | |
| due_date | DATE | no | From payment terms |
| vat_deduction_status | VARCHAR(20) | yes | pending, claimed, rejected |
| e_invoice_data | JSONB | no | Full Decree 254 XML |
| e_invoice_code | VARCHAR(100) | no | GDT invoice code |
| status | VARCHAR(20) | yes | draft, verified, posted, paid, cancelled |
| gl_posted | BOOLEAN | no | |
| gl_posted_at | TIMESTAMP | no | |
| notes | TEXT | no | |
| created_by | UUID | yes | |
| created_at | TIMESTAMP | auto | |

### 1.7 Supplier Invoice Line

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| invoice_id | UUID | yes | FK |
| po_line_id | UUID | no | FK for 3-way match |
| grn_line_id | UUID | no | FK for 3-way match |
| item_code | VARCHAR(50) | no | |
| item_name | VARCHAR(255) | yes | |
| unit | VARCHAR(20) | yes | |
| quantity | DECIMAL(18,4) | yes | |
| unit_price | DECIMAL(18,4) | yes | |
| vat_rate | DECIMAL(5,2) | yes | |
| vat_type | VARCHAR(10) | yes | |
| line_total | DECIMAL(18,2) | auto | |
| line_vat_amount | DECIMAL(18,2) | auto | |
| account_id | VARCHAR(20) | yes | GL debit account |
| vat_account_id | VARCHAR(20) | yes | 1331 |

### 1.8 AP Transaction

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK |
| supplier_id | UUID | yes | FK |
| invoice_id | UUID | no | FK |
| transaction_type | VARCHAR(20) | yes | invoice, credit_note, payment, prepayment, offset |
| transaction_date | DATE | yes | |
| amount | DECIMAL(18,2) | yes | Positive = increase AP, Negative = decrease AP |
| currency | VARCHAR(3) | yes | |
| reference_type | VARCHAR(20) | no | po, grn, payment |
| reference_id | UUID | no | |
| notes | TEXT | no | |
| created_at | TIMESTAMP | auto | |

### 1.9 Purchase Cost Allocation

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK |
| invoice_id | UUID | yes | FK |
| cost_type | VARCHAR(50) | yes | transport, insurance, customs, inspection |
| cost_amount | DECIMAL(18,2) | yes | |
| allocation_method | VARCHAR(20) | yes | by_qty, by_value, by_weight, by_volume |
| allocated_lines | JSONB | yes | [{line_id, amount}] |
| notes | TEXT | no | |

---

## 2. API Endpoints

### 2.1 Suppliers

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/purchase/suppliers | List suppliers |
| GET | /api/v1/purchase/suppliers/:id | Get supplier |
| POST | /api/v1/purchase/suppliers | Create supplier |
| PUT | /api/v1/purchase/suppliers/:id | Update supplier |
| DELETE | /api/v1/purchase/suppliers/:id | Soft-delete (suspend) |

### 2.2 Purchase Orders

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/purchase/orders | List POs (filterable by status, supplier, date range) |
| GET | /api/v1/purchase/orders/:id | Get PO with lines |
| POST | /api/v1/purchase/orders | Create PO |
| PUT | /api/v1/purchase/orders/:id | Update PO (draft only) |
| PATCH | /api/v1/purchase/orders/:id/approve | Approve PO |
| PATCH | /api/v1/purchase/orders/:id/cancel | Cancel PO |
| PATCH | /api/v1/purchase/orders/:id/close | Close PO |

### 2.3 Goods Receipts

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/purchase/receipts | List GRNs |
| GET | /api/v1/purchase/receipts/:id | Get GRN with lines |
| POST | /api/v1/purchase/receipts | Create GRN from PO |
| PUT | /api/v1/purchase/receipts/:id | Update GRN (draft only) |
| PATCH | /api/v1/purchase/receipts/:id/post | Post GRN (affects inventory) |
| PATCH | /api/v1/purchase/receipts/:id/cancel | Cancel GRN |

### 2.4 Supplier Invoices

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/purchase/invoices | List supplier invoices |
| GET | /api/v1/purchase/invoices/:id | Get invoice with lines |
| POST | /api/v1/purchase/invoices | Create invoice (manual entry) |
| POST | /api/v1/purchase/invoices/e-invoice | Receive e-invoice from GDT XML |
| PUT | /api/v1/purchase/invoices/:id | Update invoice (draft only) |
| PATCH | /api/v1/purchase/invoices/:id/verify | Verify (3-way match) |
| PATCH | /api/v1/purchase/invoices/:id/post | Post invoice (affects GL) |
| PATCH | /api/v1/purchase/invoices/:id/cancel | Cancel invoice |
| PATCH | /api/v1/purchase/invoices/:id/claim-vat | Claim VAT deduction |

### 2.5 AP

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/purchase/ap-aging | AP aging report |
| GET | /api/v1/purchase/ap-aging/supplier/:id | AP aging per supplier |
| GET | /api/v1/purchase/ap-summary | AP summary by supplier |

### 2.6 Reports

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/purchase/reports/s01-dn | Purchase ledger |
| GET | /api/v1/purchase/reports/s02-dn | Supplier detail |
| GET | /api/v1/purchase/reports/s03-dn | Goods purchase ledger |
| GET | /api/v1/purchase/reports/vat-input | VAT input tracking |
| GET | /api/v1/purchase/reports/uninvoiced-receipts | Uninvoiced GRNs |

---

## 3. State Machines

### 3.1 Purchase Order States

```
draft ──→ approved ──→ sent ──→ partial ──→ received ──→ closed
  │         │          │          │
  └──→ cancelled       └──────────┘──→ cancelled
```

### 3.2 Goods Receipt States

```
draft ──→ posted ──→ cancelled
```

### 3.3 Supplier Invoice States

```
draft ──→ verified ──→ posted ──→ paid ──→ closed
  │         │           │
  └──→ cancelled        └──→ cancelled
```

---

## 4. 3-Way Matching Logic

```
PO qty (ordered)  ─┐
                   ├── Match OK if: PO.qty ≥ Invoice.qty ∧ GRN.qty ≥ Invoice.qty
GRN qty (received) ─┤               (tolerance = configurable ±%)
                    │    AND unit_price variance < configurable threshold
Invoice qty        ─┘
```

- Exact match: all 3 quantities + prices match within tolerance → auto-verify
- Partial match: PO.qty > Invoice.qty → flag for review
- Price variance > threshold → flag for AP manager approval
- Over-invoice: Invoice.qty > GRN.qty → reject or hold

---

## 5. GL Posting Map

| Transaction | Debit Account | Credit Account | Condition |
|-------------|---------------|----------------|-----------|
| Goods receipt (inventory) | 152/153/156 | 331 (temporary) | If no invoice |
| Goods receipt (expense) | 641/642/627 | 331 | Services, non-inventory |
| Supplier invoice (goods) | 331 (reverse temp) | 331 (permanent) | Invoice received |
| VAT input | 1331 | 331 | If VAT deductible |
| Payment to supplier | 331 | 111/112 | Full or partial |
| Prepayment | 331 (debit balance) | 111/112 | Advance to supplier |
| Purchase return | 152/156 (reverse) | 331 | Return goods |
| Discount received | 331 | 515 | Early payment discount |
| Import duty | 152/156 | 3333 | Customs duty |
| Import VAT | 1331 (or expense) | 33312 | Import VAT |

---

## 6. Numbering Rules

| Document | Format | Example | Reset |
|----------|--------|---------|-------|
| PO | PO-YYYYMM-XXXX | PO-202607-0001 | Monthly |
| GRN | GRN-YYYYMM-XXXX | GRN-202607-0001 | Monthly |
| Supplier Invoice | (from supplier) | — | N/A |

---

## 7. E-Invoice Receipt (Decree 254/2026)

### 7.1 GDT XML Parse Flow

```
GDT Push → Webhook /api/v1/purchase/invoices/e-invoice
  → Parse XML (GDT standard format)
  → Validate: supplier_tax_code, invoice_number, total_amount, VAT
  → Match to PO by item description or manual
  → Create SupplierInvoice (draft)
  → Notify AP accountant
```

### 7.2 XML Fields Expected

- Invoice number, date, time
- Supplier: name, tax code, address
- Buyer: name, tax code, address
- Line items: description, unit, qty, price, VAT rate, VAT amount
- Totals: subtotal, VAT, grand total
- GDT invoice code (if tax authority coded)
- Digital signature (supplier's)
- QR code (for lookup)

---

## 8. Budget Check (P2)

```
PO Creation → Check budget_consumed + PO.total ≤ budget_allowed
             → If exceed: warn / block (configurable)
```

Budget allocation: by cost center, department, project, or account.