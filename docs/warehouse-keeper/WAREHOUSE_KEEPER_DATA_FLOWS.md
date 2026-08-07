# Data Flows — Thủ Kho (Warehouse Keeper) Module

**Version:** 1.0
**Date:** August 2026

---

## DF-001: Slip Recording Data Flow

```
Purchase Module                Warehouse Keeper Module              GL Module
     │                                │                                │
     │  Create GRN                    │                                │
     │───────────────────────────────>│                                │
     │                                │  Pending slips appear          │
     │                                │  in keeper queue               │
     │                                │                                │
     │                                │  Keeper reviews                │
     │                                │  Selects slips                 │
     │                                │  Clicks "Ghi sổ"              │
     │                                │                                │
     │                                │  Create StockLedgerEntry       │
     │                                │  (recorded_by, timestamp)      │
     │                                │                                │
     │                                │  Update StockBalance           │
     │                                │  (quantity += received)         │
     │                                │                                │
     │                                │  Create InventoryTransaction   │
     │                                │  (type=GRN_PURCHASE)           │
     │                                │                                │
     │                                │  GL Posting (by Purchase)      │
     │                                │  Dr 152 Cr 331                 │
     │                                │                                │
```

**Key point:** Keeper recording is a **parallel record** — it does NOT trigger GL posting. GL posting is handled by the source module (Purchase/Sale). The keeper's Stock Ledger is for cross-referencing only.

---

## DF-002: Un-recording Data Flow

```
Keeper                          System                         Source Module
  │                                │                                │
  │  Request un-record             │                                │
  │  (slip_id, reason)             │                                │
  │───────────────────────────────>│                                │
  │                                │  Validate:                     │
  │                                │  - Period not closed           │
  │                                │  - Not reconciled              │
  │                                │  - User authorized             │
  │                                │                                │
  │                                │  Update StockLedgerEntry:      │
  │                                │  status = "unrecorded"         │
  │                                │  unrecorded_by = user_id       │
  │                                │  unrecorded_at = now()         │
  │                                │  unrecord_reason = reason      │
  │                                │                                │
  │                                │  Slip returns to pending       │
  │                                │                                │
  │                                │  NOTE: No GL reversal          │
  │                                │  NOTE: StockBalance unchanged  │
```

---

## DF-003: Reconciliation Data Flow

```
StockLedgerEntries (Keeper)     StockBalance (Accounting)    Reconciliation Report
        │                              │                           │
        │  Query:                      │  Query:                   │
        │  WHERE warehouse_id = ?      │  WHERE warehouse_id = ?   │
        │  AND status = 'recorded'     │                           │
        │                              │                           │
        │  GROUP BY item_id            │  GROUP BY item_id         │
        │  SUM(balance_qty)            │  SUM(quantity)            │
        │                              │                           │
        ├──────────────────────────────┼──────────────────────────>│
        │                              │                           │
        │                              │  For each item:           │
        │                              │  keeper_qty = ledger sum  │
        │                              │  accounting_qty = bal sum │
        │                              │  variance = difference    │
        │                              │  value = variance * cost  │
        │                              │                           │
```

---

## DF-004: Inventory Count Data Flow

```
StockTake (Manager)       Keeper Count          System Processing
       │                       │                       │
       │  Create count         │                       │
       │  for warehouse        │                       │
       │──────────────────────>│                       │
       │                       │  Open count sheet     │
       │                       │  Enter physical qty   │
       │                       │──────────────────────>│
       │                       │                       │
       │                       │  Calculate variance   │
       │                       │  per item             │
       │                       │                       │
       │                       │  If variance > tol:   │
       │                       │  Trigger re-count     │
       │                       │<──────────────────────│
       │                       │                       │
       │  Approve count        │                       │
       │<──────────────────────│                       │
       │                       │                       │
       │  Post adjustment      │                       │
       │──────────────────────────────────────────────>│
       │                       │                       │
       │                       │  Gain: Dr 152 Cr 632  │
       │                       │  Loss: Dr 632 Cr 152  │
       │                       │                       │
       │                       │  Update StockBalance  │
       │                       │  Create InvTransaction│
```

---

## DF-005: Module Configuration Data Flow

```
Admin                         Config Store              Keeper Module
  │                                │                        │
  │  Toggle module                 │                        │
  │  on/off                        │                        │
  │───────────────────────────────>│                        │
  │                                │  Update config         │
  │                                │  module_enabled = bool │
  │                                │                        │
  │                                │  Notify module         │
  │                                │───────────────────────>│
  │                                │                        │
  │                                │  Show/hide tabs        │
  │                                │  based on config       │
```
