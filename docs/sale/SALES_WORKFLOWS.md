# Sale Module — Workflows & Processes

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)

---

## Workflow 1: Order-to-Cash (Full Cycle)

```
┌────────────────────────────────────────────────────────────────────────────┐
│                        ORDER-TO-CASH PROCESS                              │
└────────────────────────────────────────────────────────────────────────────┘

   ┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
   │  SALES     │    │  SALES     │    │ DELIVERY   │    │ CUSTOMER   │
   │ QUOTATION  │───→│  ORDER     │───→│ NOTE       │───→│ INVOICE    │
   │ (P1)       │    │ (P0)       │    │ (P0)       │    │ (P0)       │
   └────────────┘    └────────────┘    └────────────┘    └────────────┘
                                                               │
                         ┌────────────┐    ┌────────────┐      │
                         │  CUSTOMER  │←───│ AR AGING   │←─────┘
                         │  RECEIPT   │    │ MANAGEMENT │
                         │ (P0)       │    │ (P0)       │
                         └────────────┘    └────────────┘
                              │
                              ▼
                         ┌────────────┐
                         │  CREDIT    │
                         │  NOTE /    │ (if return or adjustment)
                         │  SALES     │
                         │  RETURN    │
                         │ (P0)       │
                         └────────────┘

Post-GL:
  DEL → Dr 632 (COGS)    Cr 152/156 (inventory)
  INV → Dr 131 (AR)      Cr 5111 (revenue)
        Dr 131 (AR)      Cr 3331 (VAT output)
  RCP → Dr 112/111       Cr 131 (AR)
  CN  → Dr 5111          Cr 131 (AR)
        Dr 3331           Cr 131 (AR)
```

---

## Workflow 2: Domestic Sale (VAT Taxable)

```
START ──→ Create SO ──→ Approve SO ──→ Confirm SO
  │
  ├──→ [Pick & pack goods]
  │
  ├──→ Deliver goods: Delivery Note ←→ SO (match qty)
  │     └──→ Post DN: Dr 632 (COGS)  Cr 152/156 (inventory)
  │
  ├──→ Create customer invoice (from SO + DN)
  │     └──→ 2-way match: SO × DN × Invoice
  │     └──→ Generate TXML e-invoice
  │     └──→ Digitally sign → Submit to GDT → Get code
  │     └──→ Post invoice: Dr 131  Cr 5111
  │                        Dr 131  Cr 3331
  │
  ├──→ [Payment due date arrives]
  │
  └──→ Receive payment: Dr 112/111  Cr 131
        └──→ Record FX gain/loss if foreign currency
        └──→ Allocate payment to invoice(s)

GL Entries (Domestic, goods 100M + VAT 10% = 110M, cost 70M):

  Step 1 - Delivery (COGS):
    Dr 632 (COGS)         70,000,000
    Cr 156 (Goods)        70,000,000

  Step 2 - Invoice:
    Dr 131 (AR)          110,000,000
    Cr 5111 (Revenue)    100,000,000
    Cr 3331 (VAT output)  10,000,000

  Step 3 - Receipt:
    Dr 112 (Bank)        110,000,000
    Cr 131 (AR)          110,000,000
```

---

## Workflow 3: Export Sale (Non-VAT, Customs)

```
START ──→ Create SO (export) ──→ Approve ──→ Confirm
  │
  ├──→ Prepare export declaration
  │     └──→ Generate packing list
  │
  ├──→ Customs clearance
  │     └──→ Submit export docs to customs
  │     └──→ Receive customs clearance
  │
  ├──→ Ship goods (FOB/CIF terms)
  │     └──→ Delivery Note (export warehouse)
  │     └──→ Post DN: Dr 632  Cr 152/156
  │
  ├──→ Create invoice (VAT 0% rate)
  │     └──→ E-invoice type: export
  │     └──→ GDT submission (VAT 0%)
  │     └──→ Post: Dr 131  Cr 511 (VAT 0%)
  │
  ├──→ [Payment due]
  │
  └──→ Receive foreign currency payment
        └──→ FX gain: Dr 131  Cr 515
        └──→ FX loss: Dr 635  Cr 131

GL Entries (Export $10,000, cost 200M, rate 25,000):

  Step 1 - Delivery (COGS):
    Dr 632 (COGS)        200,000,000
    Cr 156 (Goods)       200,000,000

  Step 2 - Invoice:
    Dr 131 (AR)          250,000,000  ($10,000 × 25,000)
    Cr 5111 (Revenue)    250,000,000

  Step 3 - Receipt (rate 25,500):
    Dr 112 (Bank)        255,000,000  ($10,000 × 25,500)
    Cr 131 (AR)          250,000,000
    Cr 515 (FX gain)       5,000,000
```

---

## Workflow 4: Sales Return / Credit Note

```
Goods delivered, invoice issued
  │
  ├──→ Customer returns goods
  │
  ├──→ [PRE-INVOICE RETURN]
  │     └──→ Return to warehouse (reverse delivery)
  │     └──→ Reverse COGS: Dr 152/156  Cr 632
  │     └──→ Adjust SO delivered qty
  │
  └──→ [POST-INVOICE RETURN]
        └──→ Receive returned goods (return DN)
        └──→ Create credit note
        └──→ Generate negative e-invoice TXML
        └──→ Sign + submit to GDT (credit note type)
        └──→ Post credit note:
              Dr 5111 (Revenue reversal)    Cr 131 (AR)
              Dr 3331 (VAT reversal)        Cr 131 (AR)
              Dr 152/156 (Return to stock)  Cr 632 (COGS reversal)

GL Entries (Return goods 100M, VAT 10M, cost 70M):

  Step 1 - Revenue reversal:
    Dr 5111 (Revenue)    100,000,000
    Cr 131 (AR)          100,000,000

  Step 2 - VAT reversal:
    Dr 3331 (VAT output) 10,000,000
    Cr 131 (AR)           10,000,000

  Step 3 - Inventory return:
    Dr 156 (Goods)       70,000,000
    Cr 632 (COGS)        70,000,000
```

---

## Workflow 5: Corrective/Adjustment Invoice

```
Invoice issued (original)
  │
  ├──→ [INCREASE] Amount needs to increase
  │     └──→ Create adjustment invoice (increase type)
  │     └──→ Reference original invoice number + GDT code
  │     └──→ Post: Dr 131  Cr 5111/3331 (differential)
  │
  └──→ [DECREASE] Amount needs to decrease
        └──→ Create adjustment invoice (decrease type)
        └──→ Reference original invoice number + GDT code
        └──→ Post: Dr 5111/3331  Cr 131 (differential)

  [REPLACEMENT] Full replacement
        └──→ Cancel original invoice on GDT
        └──→ Create new invoice (replacement type)
        └──→ Reference original invoice number + GDT code
        └──→ Submit as replacement to GDT
```

---

## Workflow 6: Prepayment / Customer Deposit

```
Customer pays deposit before delivery
  │
  ├──→ Record receipt (prepayment type)
  │     └──→ Dr 112  Cr 131 (credit balance = advance)
  │
  ├──→ [Optional: Issue deposit e-invoice]
  │     └──→ Decree 254: deposit invoice (advance)
  │
  ├──→ Goods delivered + invoice issued
  │
  └──→ Offset deposit against invoice
        └──→ Dr 131 (credit bal)  Cr 131 (AR invoice)
        └──→ Issue remaining balance invoice

  If deposit was e-invoiced:
        └──→ Final invoice = total - deposit amount
        └──→ Reference deposit invoice in final invoice
```

---

## Workflow 7: E-Invoice Issuance (Decree 254/2026)

```
GoTax System                              GDT Portal
     │                                        │
     │── Invoice created (draft)              │
     │                                        │
     │── Generate TXML XML                    │
     │   (Decree 254 format)                  │
     │                                        │
     │── Digitally sign XML                   │
     │   (XMLDSig, registered cert)           │
     │                                        │
     │── POST /api/v1/invoices ────────────→  │
     │                                        │── Validate XML format
     │                                        │── Verify signature
     │                                        │── Check business rules
     │                                        │── Assign invoice code
     │                                        │── Timestamp (GDT time)
     │                                        │
     │←── Response ──────────────────────────│
     │   OK: invoice_code, issue_date, QR     │
     │   FAIL: error_code, error_message      │
     │                                        │
     │── [If OK]                              │
     │   └──→ Update invoice status = coded   │
     │   └──→ Store invoice_code, QR          │
     │   └──→ Notify customer                 │
     │   └──→ Post to GL                      │
     │                                        │
     │── [If FAIL]                            │
     │   └──→ Log error                       │
     │   └──→ Notify AR accountant            │
     │   └──→ Fix and resubmit                │
```

---

## Workflow 8: Month-End Closing for Sale

```
Month-End Procedures:
  │
  ├──→ 1. Verify all delivery notes have corresponding invoices
  │     └──→ Run "Unbilled Deliveries Report"
  │     └──→ Accrue revenue for goods delivered not invoiced
  │     └──→ Dr 131 (accrued)  Cr 5111 (accrued revenue)
  │     └──→ Reverse next period
  │
  ├──→ 2. Verify all invoices posted to GL
  │     └──→ Run "AR Ledger vs GL" reconciliation
  │     └──→ AR sub-ledger balance = GL 131 balance
  │
  ├──→ 3. AR aging report review
  │     └──→ Confirm aging categories correct
  │     └──→ Flag overdue items for collection
  │     └──→ Review bad debt provision (Account 229)
  │
  ├──→ 4. VAT output reconciliation
  │     └──→ Total VAT output in sales = VAT declaration
  │     └──→ Verify e-invoice codes all match GDT records
  │
  ├──→ 5. FX revaluation (if foreign currency AR)
  │     └──→ Revalue AR at month-end rate
  │     └──→ Post FX gain/loss: 515 or 635
  │
  └──→ 6. Revenue recognition check
        └──→ Verify unearned revenue (3387) correctly released
        └──→ Deferred revenue schedule review

Accrual GL Entry (unbilled delivery):
  Dr 131 (Accrued AR)   100,000,000
  Cr 5111 (Revenue)     100,000,000

  Next month reversal:
  Dr 5111 (Revenue)     100,000,000
  Cr 131 (Accrued AR)   100,000,000
```

---

## Workflow 9: Collection Process

```
┌──────────┐    ┌────────────┐    ┌───────────┐    ┌──────────┐    ┌───────────┐
│ AR AGING │    │ COLLECTION │    │ DUNNING   │    │ RECEIPT  │    │ RECON-    │
│ REVIEW   │───→│ SCHEDULE   │───→│ (if due)   │───→│ PROCESS  │───→│ CILIATION │
└──────────┘    └────────────┘    └───────────┘    └──────────┘    └───────────┘

AR Aging Review (AR Accountant):
  - Run AR aging report
  - Identify overdue invoices
  - Prioritize by amount + aging

Collection Schedule (AR Accountant):
  - Contact customer before due date
  - Confirm payment date
  - Note communication in system

Dunning (if overdue):
  - Send payment reminder (email/SMS)
  - Level 1: 1-7 days overdue (gentle reminder)
  - Level 2: 8-30 days overdue (formal demand)
  - Level 3: 31-60 days overdue (final notice)
  - Level 4: 60+ days (legal action / bad debt write-off)

Receipt (AR Accountant / Treasurer):
  - Receive payment notification (bank/cheque/cash)
  - Record receipt in system
  - Allocate to specific invoice(s)
  - Post to GL

Reconciliation (AR Accountant):
  - Match bank statement with receipt
  - Update invoice status to paid
  - Handle discrepancies / short payments
```
