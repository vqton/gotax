# Sale Module — Data Flows & Entity Relationships

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)

---

## 1. Entity Relationship Diagram (Text)

```
customer
  │
  ├──< sales_quotation (1:N, P1)
  │       │
  │       └──< quotation_line (1:N)
  │
  ├──< sales_order (1:N)
  │       │
  │       └──< so_line (1:N)
  │             │
  │             ├──< delivery_line (1:N)
  │             │       │
  │             │       └──< delivery_note (N:1)
  │             │
  │             └──< invoice_line (1:N)
  │                       │
  │                       └──< customer_invoice (N:1)
  │
  ├──< delivery_note (1:N)
  │
  ├──< customer_invoice (1:N)
  │     │
  │     ├──< credit_note_line (1:N)
  │     │       │
  │     │       └──< credit_note (N:1)
  │     │
  │     └──< receipt_allocation (1:N)
  │             │
  │             └──< customer_receipt (N:1)
  │
  ├──< ar_transaction (1:N)
  │
  └──< customer_receipt (1:N)
```

---

## 2. Core Data Flow: Order-to-Cash

```
                    ┌─────────────────┐
                    │    Company      │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
     ┌────────▼────────┐          ┌─────────▼────────┐
     │    Customer      │          │   GL Account     │
     │  (Customer)      │          │  (Chart of Accts)│
     └────────┬────────┘          └─────────┬────────┘
              │                             │
              ▼                             │
     ┌────────────────┐                     │
     │  Sales Order   │◄────────────────────┘
     │ code, date,    │      (accounts: 5111, 3331,
     │ customer, total│       131, 632 per line)
     └───────┬────────┘
              │
              ▼
     ┌────────────────┐
     │ Delivery Note   │
     │ code, date, qty,│──→ COGS (632)
     │ SO reference    │──→ Inventory reduction (152/156)
     └───────┬────────┘
              │
              ▼
     ┌────────────────┐
     │Customer Invoice│
     │ #, date, amt,  │──→ 2-way match (SO × DN)
     │ VAT, e-invoice  │──→ AR increase (131)
     │ XML, GDT code   │──→ Revenue recognition (5111)
     │                 │──→ VAT output (3331)
     └───────┬────────┘
              │
              ▼
     ┌────────────────┐
     │   AR Ledger    │
     │ by customer,   │
     │ by invoice,    │
     │ aging buckets  │
     └───────┬────────┘
              │
              ▼
     ┌────────────────┐
     │Customer Receipt│──→ Bank (112)
     │ payment from   │──→ AR decrease (131)
     │ customer       │
     └────────────────┘

   On Sales Return:
     ┌────────────────┐
     │  Credit Note   │──→ Revenue reversal (5111)
     │  (negative     │──→ VAT reversal (3331)
     │   e-invoice)   │──→ AR decrease (131)
     │                │──→ Inventory return (152/156 → 632)
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
         ┌──────────────────────────┼───────────────────────────┐
         ▼                          ▼                           ▼
┌─────────────────┐    ┌─────────────────────┐    ┌─────────────────┐
│ Delivery Post   │    │   Invoice Post      │    │  Receipt Post   │
│                 │    │                     │    │                 │
│ Dr 632 (COGS)   │    │ Dr 131 (AR)        │    │ Dr 112/111      │
│ Cr 152/156      │    │ Cr 5111 (Revenue)  │    │ Cr 131 (AR)     │
│                 │    │ Cr 3331 (VAT)      │    │                 │
│                 │    │                     │    │                 │
└─────────────────┘    └─────────────────────┘    └─────────────────┘

       │                       │                           │
       ▼                       ▼                           ▼
┌──────────────────────────────────────────────────────────────────┐
│                        GL Account Balances                       │
│  131 (AR) ↑ │ 5111 (Revenue) ↑ │ 3331 (VAT) ↑ │ 632 (COGS) ↑    │
│  152/156 (Inventory) ↓                                          │
└──────────────────────────────────────────────────────────────────┘
```

---

## 4. E-Invoice Data Flow

```
GoTax System                    GDT Portal                  Customer
     │                              │                         │
     │── Generate TXML XML          │                         │
     │   (Decree 254 format)        │                         │
     │                              │                         │
     │── Digitally sign             │                         │
     │   (XMLDSig, registered cert) │                         │
     │                              │                         │
     │── POST /api/v1/invoices ──→ │                         │
     │                              │── Validate              │
     │                              │── Assign invoice code   │
     │                              │── Timestamp             │
     │                              │                         │
     │←── Response (code) ────────│                         │
     │                              │                         │
     │── Store invoice code         │                         │
     │── Generate PDF/QR            │                         │
     │                              │                         │
     │── Send e-invoice ──────────────────────────────────→ │
     │   (email/portal)             │                         │
     │                              │                         │
     │── Post to GL                 │                         │
```

---

## 5. AR Aging Calculation Flow

```
Customer Invoice Table      Receipt Table           AR Aging Engine
     │                            │                       │
     │── Invoice amount ──────────┤                       │
     │── Amount received  ←───────┤── Receipt amount      │
     │── Credit note amount ←─────┤── Credit note amount  │
     │── Due date                 │                       │
     │                            │                       │── For each invoice:
     │                            │                       │   balance = invoice - receipts - credit_notes
     │                            │                       │   days_overdue = today - due_date
     │                            │                       │   bucket = f(days_overdue):
     │                            │                       │      ≤ 0: "Current"
     │                            │                       │      1-30: "1-30 days"
     │                            │                       │      31-60: "31-60 days"
     │                            │                       │      61-90: "61-90 days"
     │                            │                       │      91-120: "91-120 days"
     │                            │                       │      120+: "120+ days"
     │                            │                       │
     │                            │                       │── Group by customer
     │                            │                       │── Return aging report
```

---

## 6. Integration with Existing GoTax Modules

### 6.1 Bank Module Integration

```
Customer Receipt:
  POST /api/v1/sales/receipts
    → Create ARTransaction (type=receipt)
    → Call Bank module: /api/v1/bank/transactions (if bank transfer)
    → Debit Account 112, Credit Account 131
    → Update invoice amount_received, balance_due
    → If balance_due = 0 → invoice status = paid

Customer Refund (credit note with payment):
  POST /api/v1/sales/credit-notes/:id/refund
    → Create bank transaction (outgoing)
    → Debit 131, Credit 112
```

### 6.2 Company Module Integration

```
Customer is per-company (multi-tenant):
  customer.company_id = company_id from auth middleware
  All sales queries filtered by company context
```

### 6.3 Auth Module Integration

```
Role-based access (see SALES_RULES.md R16):
  Auth middleware extracts user_id, role from JWT
  Sales handler checks role for each operation
  Audit log: every mutation logged with user_id
```

### 6.4 GL Module Integration

```
Auto-posting via existing JournalEntry API:
  Delivery  → POST /api/v1/journal-entries (COGS)
  Invoice   → POST /api/v1/journal-entries (Revenue + AR + VAT)
  Receipt   → POST /api/v1/journal-entries (AR + Bank)
  Credit Note → POST /api/v1/journal-entries (Reversals)
```

### 6.5 Tax Module Integration

```
VAT Output Declaration:
  Tax module reads from:
    - customer_invoice.vat_type per line
    - customer_invoice.tax_amount by VAT rate
  Tax declaration form 01/GTGT:
    - Line 25-27: VAT output on goods/services (rate breakdown)
  Monthly reconciliation:
    - Total 3331 credit = VAT declaration output
```

### 6.6 Cash Module Integration

```
Customer Payment via Petty Cash:
  POST /api/v1/sales/receipts (payment_method = cash)
    → Also creates cash receipt in Cash module
    → Uses existing ReceiptSales type, CounterpartCustomer
```

### 6.7 e-Invoice / Digital Signature Integration

```
E-invoice signing uses existing DigitalSignature model:
  digital_signature_id on customer_invoice
  DigitalSignatures CRUD already exists in Company module
  New: TXML generator + GDT API client
  New: QR code generator
```
