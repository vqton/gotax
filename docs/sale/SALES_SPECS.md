# Sale Module — Functional Specifications

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)

---

## 1. Data Model

### 1.1 Customer

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK → company |
| code | VARCHAR(20) | yes | Unique per company |
| name | VARCHAR(255) | yes | |
| tax_code | VARCHAR(20) | yes | MST, validated format (10/13/14 digits) |
| address | TEXT | no | |
| phone | VARCHAR(20) | no | |
| email | VARCHAR(100) | no | |
| bank_account_name | VARCHAR(255) | no | For refunds |
| bank_account_number | VARCHAR(50) | no | |
| bank_name | VARCHAR(255) | no | |
| payment_terms | VARCHAR(50) | no | net15, net30, net45, net60, COD |
| credit_limit | DECIMAL(18,2) | no | 0 = no limit |
| currency | VARCHAR(3) | no | Default VND |
| customer_type | VARCHAR(20) | no | domestic, export, both |
| customer_group | VARCHAR(50) | no | retail, wholesale, distributor, agent |
| price_list_id | UUID | no | FK → price_list (tiered pricing) |
| status | VARCHAR(20) | yes | active, suspended, blacklisted |
| notes | TEXT | no | |
| created_by | UUID | yes | FK → user |
| created_at | TIMESTAMP | auto | |
| updated_at | TIMESTAMP | auto | |

### 1.2 Sales Order

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK → company |
| so_number | VARCHAR(30) | yes | Auto-generated: SO-YYYYMM-XXXX |
| quotation_id | UUID | no | FK → sales_quotation |
| customer_id | UUID | yes | FK → customer |
| order_date | DATE | yes | |
| expected_date | DATE | no | Expected delivery date |
| currency | VARCHAR(3) | yes | Default VND |
| exchange_rate | DECIMAL(18,6) | no | If foreign currency |
| payment_terms | VARCHAR(50) | no | |
| delivery_terms | VARCHAR(50) | no | FOB, CIF, DDP, EXW |
| shipping_address | TEXT | no | |
| subtotal | DECIMAL(18,2) | auto | Sum of line amounts (before discount) |
| discount_amount | DECIMAL(18,2) | no | Order-level discount |
| tax_amount | DECIMAL(18,2) | auto | Total VAT output |
| total_amount | DECIMAL(18,2) | auto | |
| status | VARCHAR(20) | yes | draft, approved, confirmed, processing, delivered, invoiced, cancelled, closed |
| approved_by | UUID | no | FK → user |
| approved_at | TIMESTAMP | no | |
| cancelled_reason | TEXT | no | |
| notes | TEXT | no | |
| created_by | UUID | yes | FK → user |
| created_at | TIMESTAMP | auto | |
| updated_at | TIMESTAMP | auto | |

### 1.3 Sales Order Line

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| so_id | UUID | yes | FK → sales_order |
| line_number | INT | yes | |
| item_code | VARCHAR(50) | no | |
| item_name | VARCHAR(255) | yes | |
| unit | VARCHAR(20) | yes | kg, pcs, box, m |
| quantity | DECIMAL(18,4) | yes | |
| unit_price | DECIMAL(18,4) | yes | |
| discount_pct | DECIMAL(5,2) | no | Line-level discount % |
| vat_rate | DECIMAL(5,2) | yes | 0, 5, 8, 10 |
| vat_type | VARCHAR(10) | yes | VAT_0, VAT_5, VAT_8, VAT_10, NT (non-taxable), KCT (not deductible) |
| revenue_account_id | VARCHAR(20) | yes | GL account (5111/5112/5113) |
| vat_account_id | VARCHAR(20) | yes | GL account 3331 |
| line_total | DECIMAL(18,2) | auto | quantity × unit_price × (1 - discount_pct) |
| line_vat_amount | DECIMAL(18,2) | auto | |
| delivered_qty | DECIMAL(18,4) | auto | Cumulative from delivery notes |
| invoiced_qty | DECIMAL(18,4) | auto | Cumulative from invoices |

### 1.4 Delivery Note

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK → company |
| dn_number | VARCHAR(30) | yes | Auto: DN-YYYYMM-XXXX |
| so_id | UUID | yes | FK → sales_order |
| delivery_date | DATE | yes | |
| warehouse | VARCHAR(50) | no | Warehouse code |
| shipping_method | VARCHAR(50) | no | |
| carrier_name | VARCHAR(100) | no | |
| tracking_number | VARCHAR(100) | no | |
| delivery_address | TEXT | no | |
| status | VARCHAR(20) | yes | draft, posted, cancelled |
| notes | TEXT | no | |
| created_by | UUID | yes | |
| created_at | TIMESTAMP | auto | |

### 1.5 Delivery Line

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| dn_id | UUID | yes | FK |
| so_line_id | UUID | yes | FK → so_line |
| item_code | VARCHAR(50) | no | |
| item_name | VARCHAR(255) | yes | |
| unit | VARCHAR(20) | yes | |
| quantity_delivered | DECIMAL(18,4) | yes | |
| quantity_returned | DECIMAL(18,4) | no | Post-delivery return |
| unit_price | DECIMAL(18,4) | auto | From SO |
| line_total | DECIMAL(18,2) | auto | |
| cost_price | DECIMAL(18,4) | no | For COGS calculation |

### 1.6 Customer Invoice

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK |
| invoice_number | VARCHAR(30) | yes | Auto: INV-YYYYMM-XXXX |
| invoice_date | DATE | yes | |
| so_id | UUID | no | FK → sales_order |
| dn_id | UUID | no | FK → delivery_note |
| customer_id | UUID | yes | FK |
| customer_name | VARCHAR(255) | yes | Denormalized |
| customer_tax_code | VARCHAR(20) | yes | Denormalized |
| customer_address | TEXT | yes | Denormalized for e-invoice |
| invoice_type | VARCHAR(20) | yes | domestic, export, service |
| currency | VARCHAR(3) | yes | |
| exchange_rate | DECIMAL(18,6) | no | |
| subtotal | DECIMAL(18,2) | auto | |
| discount_amount | DECIMAL(18,2) | no | |
| tax_amount | DECIMAL(18,2) | auto | |
| total_amount | DECIMAL(18,2) | auto | |
| amount_received | DECIMAL(18,2) | no | Cumulative payments |
| balance_due | DECIMAL(18,2) | auto | |
| due_date | DATE | no | From payment terms |
| invoice_note | TEXT | no | Printed on invoice |
| e_invoice_data | TEXT | no | Full TXML XML |
| e_invoice_code | VARCHAR(100) | no | GDT invoice code |
| e_invoice_status | VARCHAR(20) | yes | pending, signed, submitted, coded, issued, cancelled, replaced |
| digital_signature_id | UUID | no | FK → digital_signature |
| signed_data | TEXT | no | Signed XML |
| gdt_response | JSONB | no | Raw GDT API response |
| original_invoice_id | UUID | no | For corrective/replacement |
| adjustment_type | VARCHAR(20) | no | increase, decrease, replacement |
| status | VARCHAR(20) | yes | draft, posted, paid, cancelled |
| gl_posted | BOOLEAN | no | |
| gl_posted_at | TIMESTAMP | no | |
| notes | TEXT | no | |
| created_by | UUID | yes | |
| created_at | TIMESTAMP | auto | |

### 1.7 Invoice Line

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| invoice_id | UUID | yes | FK |
| so_line_id | UUID | no | FK for 2-way match |
| dn_line_id | UUID | no | FK for 2-way match |
| item_code | VARCHAR(50) | no | |
| item_name | VARCHAR(255) | yes | |
| unit | VARCHAR(20) | yes | |
| quantity | DECIMAL(18,4) | yes | |
| unit_price | DECIMAL(18,4) | yes | |
| discount_pct | DECIMAL(5,2) | no | |
| vat_rate | DECIMAL(5,2) | yes | |
| vat_type | VARCHAR(10) | yes | |
| line_total | DECIMAL(18,2) | auto | |
| line_vat_amount | DECIMAL(18,2) | auto | |
| revenue_account_id | VARCHAR(20) | yes | GL credit account (5111/5112/5113) |
| vat_account_id | VARCHAR(20) | yes | 3331 |

### 1.8 AR Transaction

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK |
| customer_id | UUID | yes | FK |
| invoice_id | UUID | no | FK |
| transaction_type | VARCHAR(20) | yes | invoice, credit_note, receipt, prepayment, offset |
| transaction_date | DATE | yes | |
| amount | DECIMAL(18,2) | yes | Positive = increase AR, Negative = decrease AR |
| currency | VARCHAR(3) | yes | |
| reference_type | VARCHAR(20) | no | so, dn, receipt |
| reference_id | UUID | no | |
| notes | TEXT | no | |
| created_by | UUID | yes | |
| created_at | TIMESTAMP | auto | |

### 1.9 Customer Receipt

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK |
| receipt_number | VARCHAR(30) | yes | Auto: RCP-YYYYMM-XXXX |
| customer_id | UUID | yes | FK |
| receipt_date | DATE | yes | |
| payment_method | VARCHAR(20) | yes | cash, bank_transfer, cheque, credit_card |
| bank_account_id | UUID | no | FK → bank_account (if bank transfer) |
| currency | VARCHAR(3) | yes | |
| exchange_rate | DECIMAL(18,6) | no | |
| amount | DECIMAL(18,2) | yes | |
| unallocated_amount | DECIMAL(18,2) | auto | After invoice allocation |
| reference | VARCHAR(100) | no | Bank reference/transaction ID |
| notes | TEXT | no | |
| status | VARCHAR(20) | yes | draft, posted, reconciled, cancelled |
| created_by | UUID | yes | |
| created_at | TIMESTAMP | auto | |

### 1.10 Receipt Allocation

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| receipt_id | UUID | yes | FK |
| invoice_id | UUID | yes | FK |
| allocated_amount | DECIMAL(18,2) | yes | |
| discount_amount | DECIMAL(18,2) | no | Early payment discount (521) |

### 1.11 Credit Note (Sales Return)

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK |
| cn_number | VARCHAR(30) | yes | Auto: CN-YYYYMM-XXXX |
| original_invoice_id | UUID | yes | FK → customer_invoice |
| customer_id | UUID | yes | FK |
| return_date | DATE | yes | |
| return_reason | VARCHAR(200) | yes | defect, wrong_item, excess, other |
| return_type | VARCHAR(20) | yes | full_return, partial_return, price_adjustment |
| dn_id | UUID | no | FK → delivery_note (if goods returned) |
| subtotal | DECIMAL(18,2) | auto | |
| tax_amount | DECIMAL(18,2) | auto | |
| total_amount | DECIMAL(18,2) | auto | |
| e_invoice_data | TEXT | no | Negative e-invoice XML |
| e_invoice_code | VARCHAR(100) | no | GDT code for credit note |
| status | VARCHAR(20) | yes | draft, posted, cancelled |
| gl_posted | BOOLEAN | no | |
| notes | TEXT | no | |
| created_by | UUID | yes | |
| created_at | TIMESTAMP | auto | |

### 1.12 Credit Note Line

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| cn_id | UUID | yes | FK |
| invoice_line_id | UUID | yes | FK → invoice_line (original) |
| item_name | VARCHAR(255) | yes | |
| unit | VARCHAR(20) | yes | |
| quantity | DECIMAL(18,4) | yes | Negative quantity |
| unit_price | DECIMAL(18,4) | yes | |
| vat_rate | DECIMAL(5,2) | yes | |
| line_total | DECIMAL(18,2) | auto | |
| line_vat_amount | DECIMAL(18,2) | auto | |

### 1.13 Sales Quotation (P1)

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| id | UUID | auto | PK |
| company_id | UUID | yes | FK |
| qn_number | VARCHAR(30) | yes | Auto: QN-YYYYMM-XXXX |
| customer_id | UUID | yes | FK |
| valid_until | DATE | no | |
| status | VARCHAR(20) | yes | draft, sent, accepted, rejected, expired, converted |
| total_amount | DECIMAL(18,2) | auto | |
| created_by | UUID | yes | |
| created_at | TIMESTAMP | auto | |

---

## 2. API Endpoints

### 2.1 Customers

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/customers | List customers |
| GET | /api/v1/sales/customers/:id | Get customer |
| POST | /api/v1/sales/customers | Create customer |
| PUT | /api/v1/sales/customers/:id | Update customer |
| DELETE | /api/v1/sales/customers/:id | Soft-delete (suspend) |

### 2.2 Sales Orders

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/orders | List SOs (filterable by status, customer, date range) |
| GET | /api/v1/sales/orders/:id | Get SO with lines |
| POST | /api/v1/sales/orders | Create SO |
| POST | /api/v1/sales/orders/from-quotation/:id | Create SO from quotation |
| PUT | /api/v1/sales/orders/:id | Update SO (draft only) |
| PATCH | /api/v1/sales/orders/:id/approve | Approve SO |
| PATCH | /api/v1/sales/orders/:id/confirm | Confirm SO (reserve inventory) |
| PATCH | /api/v1/sales/orders/:id/cancel | Cancel SO |
| PATCH | /api/v1/sales/orders/:id/close | Close SO |

### 2.3 Delivery Notes

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/deliveries | List delivery notes |
| GET | /api/v1/sales/deliveries/:id | Get delivery note with lines |
| POST | /api/v1/sales/deliveries | Create delivery note from SO |
| PUT | /api/v1/sales/deliveries/:id | Update delivery note (draft only) |
| PATCH | /api/v1/sales/deliveries/:id/post | Post delivery note (affects COGS + inventory) |
| PATCH | /api/v1/sales/deliveries/:id/cancel | Cancel delivery note |

### 2.4 Customer Invoices

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/invoices | List invoices |
| GET | /api/v1/sales/invoices/:id | Get invoice with lines |
| POST | /api/v1/sales/invoices | Create invoice (from SO or delivery) |
| POST | /api/v1/sales/invoices/:id/sign | Digitally sign e-invoice |
| POST | /api/v1/sales/invoices/:id/submit | Submit e-invoice to GDT |
| PUT | /api/v1/sales/invoices/:id | Update invoice (draft only) |
| PATCH | /api/v1/sales/invoices/:id/post | Post invoice (affects GL) |
| PATCH | /api/v1/sales/invoices/:id/cancel | Cancel invoice |
| PATCH | /api/v1/sales/invoices/:id/replace | Replace invoice (corrective) |

### 2.5 Customer Receipts

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/receipts | List receipts |
| GET | /api/v1/sales/receipts/:id | Get receipt with allocations |
| POST | /api/v1/sales/receipts | Record customer payment |
| PUT | /api/v1/sales/receipts/:id | Update receipt (draft only) |
| PATCH | /api/v1/sales/receipts/:id/post | Post receipt (affects AR + GL) |
| PATCH | /api/v1/sales/receipts/:id/allocate | Allocate payment to invoices |
| PATCH | /api/v1/sales/receipts/:id/cancel | Cancel receipt |

### 2.6 Credit Notes

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/credit-notes | List credit notes |
| GET | /api/v1/sales/credit-notes/:id | Get credit note with lines |
| POST | /api/v1/sales/credit-notes | Create credit note |
| POST | /api/v1/sales/credit-notes/:id/sign | Sign credit note e-invoice |
| POST | /api/v1/sales/credit-notes/:id/submit | Submit credit note to GDT |
| PATCH | /api/v1/sales/credit-notes/:id/post | Post credit note (affects GL) |
| PATCH | /api/v1/sales/credit-notes/:id/cancel | Cancel credit note |

### 2.7 AR

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/ar-aging | AR aging report |
| GET | /api/v1/sales/ar-aging/customer/:id | AR aging per customer |
| GET | /api/v1/sales/ar-summary | AR summary by customer |
| GET | /api/v1/sales/customer-statement/:id | Customer statement of account |

### 2.8 Sales Quotations (P1)

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/quotations | List quotations |
| GET | /api/v1/sales/quotations/:id | Get quotation |
| POST | /api/v1/sales/quotations | Create quotation |
| PATCH | /api/v1/sales/quotations/:id/convert | Convert quotation to SO |

### 2.9 Reports

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sales/reports/s01-bh | Sales ledger (per customer) |
| GET | /api/v1/sales/reports/s02-bh | Customer detail ledger |
| GET | /api/v1/sales/reports/s03-bh | Goods sales ledger |
| GET | /api/v1/sales/reports/vat-output | VAT output tracking |
| GET | /api/v1/sales/reports/unbilled-deliveries | Unbilled delivery report |

---

## 3. State Machines

### 3.1 Sales Order States

```
draft ──→ approved ──→ confirmed ──→ processing ──→ delivered ──→ invoiced ──→ closed
  │         │             │              │               │
  └──→ cancelled         └──────────────┴───────────────┴──→ cancelled
```

- **draft**: Being created, editable
- **approved**: Approved by authorized person
- **confirmed**: Credit check passed, inventory reserved
- **processing**: Partially delivered
- **delivered**: Fully delivered
- **invoiced**: Customer invoice created
- **cancelled**: Cancelled before fulfillment
- **closed**: All obligations complete

### 3.2 Delivery Note States

```
draft ──→ posted ──→ cancelled
```

- **draft**: Being created
- **posted**: Goods delivered, COGS posted
- **cancelled**: Reversed

### 3.3 Customer Invoice States

```
draft ──→ signed ──→ submitted ──→ coded ──→ issued ──→ posted ──→ paid ──→ closed
  │         │          │             │         │          │
  └──→ cancelled       └─────────────┴─────────┴──────────┴──→ cancelled
  └──→ replaced ──→ new invoice referencing this invoice
```

- **draft**: Being created, editable
- **signed**: Digitally signed
- **submitted**: Sent to GDT portal
- **coded**: GDT returned invoice code
- **issued**: Legally valid e-invoice
- **posted**: GL entries generated
- **paid**: Fully paid by customer
- **cancelled**: Per regulatory cancellation
- **replaced**: Replaced by corrective invoice

### 3.4 Credit Note States

```
draft ──→ signed ──→ submitted ──→ posted ──→ closed
  │
  └──→ cancelled
```

### 3.5 Receipt States

```
draft ──→ posted ──→ reconciled ──→ cancelled
```

---

## 4. 2-Way Matching Logic

```
SO qty (ordered) ─┐
                   ├── Match OK if SO.qty ≥ Invoice.qty ∧ DN.qty ≥ Invoice.qty
DN qty (delivered) ─┤           (tolerance = configurable ±%)
                    │    AND unit_price variance < configurable threshold
Invoice qty        ─┘
```

- **Exact match**: all quantities + prices match → auto-verify
- **Partial match**: SO.qty > Invoice.qty → flag for review
- **Price variance** > threshold → flag for AR manager approval
- **Over-invoice**: Invoice.qty > DN.qty → reject or hold
- **Under-invoice**: Invoice.qty < DN.qty → allow (partial billing)

---

## 5. GL Posting Map

| Transaction | Debit Account | Credit Account | Condition |
|-------------|---------------|----------------|-----------|
| Delivery (goods) | 632 (COGS) | 152/156 (inventory) | Cost price |
| Invoice (goods, domestic) | 131 (AR) | 5111 (revenue) | Revenue recognition |
| Invoice (services) | 131 (AR) | 5112 (revenue) | Revenue recognition |
| VAT output | 131 (AR) | 3331 (VAT output) | If VAT taxable |
| Customer receipt | 111/112 | 131 (AR) | Full or partial |
| Prepayment from customer | 111/112 | 131 (credit balance) | Deposit received |
| Prepayment offset | 131 (credit balance) | 131 (AR) | Against invoice |
| Credit note (return) | 5111/5112 | 131 (AR) | Revenue reversal |
| VAT reversal on credit note | 3331 | 131 (AR) | VAT output reversal |
| COGS reversal on return | 152/156 | 632 | Goods returned to inventory |
| Trade discount | 5211 | 131 (AR) | Commercial discount |
| Sales allowance | 5212 | 131 (AR) | Post-sale allowance |
| Early payment discount | 5213 | 131 (AR) | Discount granted |
| Unearned revenue (deposit) | 111/112 | 3387 | If deferred revenue |
| Revenue recognition from deposit | 3387 | 5111 | When obligation satisfied |
| FX loss on AR | 635 | 131 | AR revaluation loss |
| FX gain on AR | 131 | 515 | AR revaluation gain |

---

## 6. Numbering Rules

| Document | Format | Example | Reset |
|----------|--------|---------|-------|
| Sales Order | SO-YYYYMM-XXXX | SO-202607-0001 | Monthly |
| Delivery Note | DN-YYYYMM-XXXX | DN-202607-0001 | Monthly |
| Customer Invoice | INV-YYYYMM-XXXX | INV-202607-0001 | Monthly |
| Credit Note | CN-YYYYMM-XXXX | CN-202607-0001 | Monthly |
| Customer Receipt | RCP-YYYYMM-XXXX | RCP-202607-0001 | Monthly |
| Quotation | QN-YYYYMM-XXXX | QN-202607-0001 | Monthly |

---

## 7. E-Invoice Issuance (Decree 254/2026)

### 7.1 E-Invoice Pipeline

```
Invoice Created (draft)
  → Generate TXML (Decree 254 format)
  → Digitally sign XML
  → Submit to GDT API (POST /api/v1/invoices)
  → Receive GDT response:
       - OK: invoice_code + issue_date + QR code
       - FAIL: error → fix + resubmit
  → Update invoice with GDT code
  → Send e-invoice to customer (email/portal)
  → Post to GL
```

### 7.2 TXML XML Fields Generated

- Invoice number, date, time (UTC+7)
- Seller: name, tax code, address, phone
- Buyer: name, tax code, address, bank account
- Line items: description, unit, qty, price, VAT rate/amount, discount
- Totals: subtotal, discount, VAT, grand total (in words)
- Signature: seller's digital signature (XMLDSig)
- GDT invoice code (assigned by GDT)
- QR code: GDT lookup URL + invoice info
- Converted amount (if foreign currency)

### 7.3 GDT API Integration

| Operation | Endpoint | Description |
|-----------|----------|-------------|
| Submit Invoice | POST /api/v1/invoices | Submit signed e-invoice |
| Cancel Invoice | POST /api/v1/invoices/{code}/cancel | Cancel issued e-invoice |
| Replace Invoice | POST /api/v1/invoices/{code}/replace | Replace with corrective invoice |
| Query Status | GET /api/v1/invoices/{code} | Check invoice status |
| Download PDF | GET /api/v1/invoices/{code}/pdf | Get issued invoice PDF |

---

## 8. Credit Check (P1)

```
SO Confirmation → Check customer credit_limit
                 → Sum of outstanding AR + SO.total ≤ credit_limit
                 → If exceed → warn or block (configurable)
                 → If no limit → auto-pass
```
