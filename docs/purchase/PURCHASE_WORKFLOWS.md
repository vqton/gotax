# Purchase Module — Workflows & Processes

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)

---

## Workflow 1: Procure-to-Pay (Full Cycle)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        PROCURE-TO-PAY PROCESS                          │
└─────────────────────────────────────────────────────────────────────────┘

   ┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
   │REQUISITION │    │ PURCHASE   │    │ GOODS      │    │ SUPPLIER   │
   │            │───→│ ORDER      │───→│ RECEIPT    │───→│ INVOICE    │
   │(P1)        │    │ (P0)       │    │ (P0)       │    │ (P0)       │
   └────────────┘    └────────────┘    └────────────┘    └────────────┘
                                                               │
                         ┌────────────┐    ┌────────────┐      │
                         │   PAYMENT  │←───│ 3-WAY      │←─────┘
                         │            │    │ MATCH      │
                         │ (P0)       │    │ (P0)       │
                         └────────────┘    └────────────┘

Post-GL:
  GRN  →  Dr 152/156     Cr 331 (temp)
  INV  →  Dr 331 (temp)  Cr 331 (AP)
          Dr 1331         Cr 331 (AP)
  PAY  →  Dr 331 (AP)    Cr 112
```

---

## Workflow 2: Domestic Purchase (VAT Deductible)

```
START ──→ Create PO ──→ Approve PO ──→ Send PO to Supplier
  │
  ├──→ [Supplier delivers goods + issues e-invoice]
  │
  ├──→ Receive goods: GRN ←→ PO (match qty)
  │     └──→ Post GRN: Dr 152/156   Cr 331 (temp price)
  │
  ├──→ Receive e-invoice (from GDT or manual entry)
  │     └──→ 3-way match: PO × GRN × Invoice
  │     └──→ Post invoice: Dr 331(temp)  Cr 331(AP)
  │                        Dr 1331       Cr 331(AP)
  │
  ├──→ Payment due date
  │
  └──→ Pay supplier: Dr 331(AP)  Cr 112
        └──→ Record FX gain/loss if foreign currency

GL Entries (Domestic, VAT 10%):
  Purchase of goods 100M + VAT 10M = 110M

  Step 1 - GRN posting (if goods received before invoice):
    Dr 156 (Goods)      100,000,000
    Cr 331 (AP temp)    100,000,000

  Step 2 - Invoice posting:
    Dr 331 (AP temp)    100,000,000
    Dr 1331 (VAT input)  10,000,000
    Cr 331 (AP)         110,000,000

  Step 3 - Payment:
    Dr 331 (AP)         110,000,000
    Cr 112 (Bank)       110,000,000
```

---

## Workflow 3: Import Purchase

```
START ──→ Create import PO ──→ PO approved
  │
  ├──→ Supplier ships goods (FOB/CIF terms)
  │
  ├──→ Customs declaration
  │     └──→ Pay import duty (Dr 152/156  Cr 3333)
  │     └──→ Pay import VAT (Dr 1331  Cr 33312) if deductible
  │     └──→ Pay special tax (Dr 152/156  Cr 3332) if applicable
  │
  ├──→ Goods arrive at warehouse → GRN
  │     └──→ Post GRN: Dr 156 (including duty + freight)  Cr 331
  │
  ├──→ Record supplier invoice (foreign currency)
  │     └──→ 3-way match with PO
  │     └──→ Post invoice: Dr 331(temp)  Cr 331(AP) at FX rate
  │
  ├──→ Allocate purchase costs (freight, insurance, duty)
  │     └──→ Dr 156 (cost allocation)  Cr 331
  │
  └──→ Pay supplier (FX rate may differ)
        └──→ FX gain: Dr 331  Cr 515
        └──→ FX loss: Dr 635  Cr 331

GL Entries (Import, value $10,000, duty 5%, VAT 10%):
  Assume rate 25,000 VND/USD

  Step 1 - Duty payment:
    Dr 156 (Goods)      12,500,000  (10,000 × 5% × 25,000)
    Cr 112              12,500,000

  Step 2 - GRN (cost = $10,000 × 25,000 + duty):
    Dr 156 (Goods)     262,500,000  (250M + 12.5M)
    Cr 331 (AP temp)   250,000,000  (supplier portion)
    Cr 112 (paid)       12,500,000  (duty)

  Step 3 - Import VAT (deductible):
    Dr 1331              26,250,000  (262.5M × 10%)
    Cr 33312             26,250,000
```

---

## Workflow 4: Purchase Return

```
Goods received (GRN posted)
  │
  ├──→ Quality reject / wrong item / excess
  │
  ├──→ [PRE-INVOICE RETURN]
  │     └──→ Create return GRN (negative quantity)
  │     └──→ Reverse GL: Dr 331(temp)  Cr 152/156
  │     └──→ Adjust PO received qty
  │
  └──→ [POST-INVOICE RETURN]
        └──→ Request credit note from supplier
        └──→ Receive credit note (negative invoice)
        └──→ 3-way match credit note
        └──→ Post credit note:
              Dr 331(AP)  Cr 152/156
              Dr 331(AP)  Cr 1331
        └──→ Adjust AP balance
```

---

## Workflow 5: Month-End Closing for Purchase

```
Month-End Procedures:
  │
  ├──→ 1. Verify all GRNs have corresponding invoices
  │     └──→ Run "Uninvoiced Receipts Report"
  │     └──→ Accrue for goods received not invoiced (Dr 152/156  Cr 338)
  │
  ├──→ 2. Verify all supplier invoices posted to GL
  │     └──→ Run "Supplier Ledger vs GL" reconciliation
  │     └──→ AP sub-ledger balance = GL 331 balance
  │
  ├──→ 3. AP aging report review
  │     └──→ Confirm aging categories correct
  │     └──→ Flag overdue items for follow-up
  │
  ├──→ 4. Supplier statement reconciliation (if received)
  │     └──→ Match supplier statement vs GoTax AP balance
  │     └──→ Investigate discrepancies
  │
  ├──→ 5. VAT input reconciliation
  │     └──→ Total VAT input in purchase = VAT deduction in tax return
  │
  └──→ 6. FX revaluation (if foreign currency AP)
        └──→ Revalue AP at month-end rate
        └──→ Post FX gain/loss: 515 or 635
```

---

## Workflow 6: E-Invoice Receipt from GDT

```
GDT System                               GoTax System
    │                                        │
    │── Supplier issues e-invoice ──→        │
    │                                        │
    │── GDT receives & codes (if required)   │
    │                                        │
    │── GDT pushes to buyer webhook ──────→  │
    │                                        │── Parse XML
    │                                        │── Validate XML
    │                                        │── Find supplier by tax code
    │                                        │── Find matching PO (if any)
    │                                        │── Create draft invoice
    │                                        │── Store raw XML
    │                                        │── Notify AP accountant
    │                                        │
    │←── Response (accepted/rejected) ──────│
    │                                        │
    │                                        │── [AP Accountant reviews]
    │                                        │── 3-way match
    │                                        │── Post invoice → GL
```

---

## Workflow 7: Supplier Payment Process

```
┌──────────┐    ┌───────────┐    ┌───────────┐    ┌───────────┐    ┌──────────┐
│ PAYMENT  │    │ PAYMENT   │    │ APPROVAL  │    │ EXECUTION │    │ RECON-   │
│ PROPOSAL │───→│ SCHEDULE  │───→│           │───→│           │───→│ CILIATION│
│          │    │           │    │           │    │           │    │          │
└──────────┘    └───────────┘    └───────────┘    └───────────┘    └──────────┘

Payment Proposal (AP Accountant):
  - Select supplier(s)
  - Select due invoices
  - Group by due date / priority
  - Calculate total payment amount

Payment Schedule (AP Manager):
  - Review proposal
  - Adjust priorities based on cash position
  - Confirm payment dates

Approval (Chief Accountant / CFO):
  - Approve large payments (> threshold)
  - Verify bank account details
  - Authorize payment

Execution (Treasurer):
  - Generate payment order (UNC/transfer)
  - Submit to bank
  - Record payment in system

Reconciliation (AP Accountant):
  - Match bank statement with payment
  - Close invoice
  - Address failed payments
```