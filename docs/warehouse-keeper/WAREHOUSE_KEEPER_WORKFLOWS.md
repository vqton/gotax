# Workflows — Thủ Kho (Warehouse Keeper) Module

**Version:** 1.0
**Date:** August 2026

---

## WF-001: Keeper Assignment Workflow

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Admin   │────>│  Create  │────>│  Active  │
│  creates │     │  Record  │     │  Period  │
└──────────┘     └──────────┘     └──────────┘
                                      │
                                      │ EffectiveTo reached
                                      ▼
                                   ┌──────────┐
                                   │ Expired  │
                                   └──────────┘
                                      or
                                   ┌──────────┐
                                   │Cancelled │
                                   └──────────┘
```

**Steps:**
1. Admin creates assignment (warehouse, user, role, effective dates)
2. System validates no overlap
3. Assignment becomes ACTIVE immediately (or on EffectiveFrom if future)
4. On EffectiveTo date, assignment auto-expires
5. Admin can cancel assignment early (sets EffectiveTo = today)

---

## WF-002: Slip Recording Workflow (Ghi sổ)

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ Purchase │────>│  GRN     │────>│ Pending  │────>│  Keeper  │
│ /Sale    │     │  Created │     │  Slips   │     │  Reviews │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                                                       │
                                              ┌────────┴────────┐
                                              │                 │
                                              ▼                 ▼
                                        ┌──────────┐     ┌──────────┐
                                        │  Record  │     │   Skip   │
                                        │ (Ghi sổ) │     │ (Bỏ qua) │
                                        └──────────┘     └──────────┘
                                              │
                                              ▼
                                        ┌──────────┐
                                        │ Recorded │
                                        │ in Ledger│
                                        └──────────┘
                                              │
                                              ▼
                                        ┌──────────┐
                                        │Reconciled│
                                        │(optional)│
                                        └──────────┘
```

**Steps:**
1. Purchase/Sales module creates GRN/DN
2. Slip appears in keeper's pending queue
3. Keeper reviews slip details (items, quantities)
4. Keeper selects slips and clicks "Ghi sổ"
5. System creates StockLedgerEntry with keeper ID and timestamp
6. Slip status changes to "Recorded"
7. Cross-reference report can now compare keeper vs accounting records

**Timing:** Within 1 business day of slip creation (recommended)

---

## WF-003: Un-recording Workflow (Bỏ ghi)

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Keeper  │────>│  Select  │────>│  Enter   │────>│  System  │
│ requests │     │  Recorded│     │  Reason  │     │  Process │
│  unrecord│     │  Slip    │     │          │     │          │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                                                       │
                                              ┌────────┴────────┐
                                              │                 │
                                              ▼                 ▼
                                        ┌──────────┐     ┌──────────┐
                                        │ Success  │     │  Reject  │
                                        │ Slip back│     │ (closed  │
                                        │ to pendng│     │  period) │
                                        └──────────┘     └──────────┘
```

**Steps:**
1. Keeper identifies slip that was recorded in error
2. Keeper selects recorded slip, clicks "Bỏ ghi"
3. System requires reason (mandatory text)
4. System validates: period not closed, slip not reconciled
5. System marks entry as unrecorded
6. Slip returns to pending queue

**Guards:**
- Cannot un-record in closed period
- Cannot un-record reconciled entries
- Reason is mandatory
- Creates audit trail (unrecorded_by, unrecorded_at, reason)

---

## WF-004: Physical Inventory Count (Keeper Flow)

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Manager  │────>│  Create  │────>│  Keeper  │
│ creates  │     │Stock Take│     │ assigned │
│ count    │     │          │     │          │
└──────────┘     └──────────┘     └──────────┘
                                      │
                                      ▼
                                ┌──────────┐
                                │  Keeper  │
                                │  Counts  │
                                │ (blind   │
                                │  or open)│
                                └──────────┘
                                      │
                                ┌─────┴─────┐
                                │           │
                                ▼           ▼
                          ┌──────────┐ ┌──────────┐
                          │ Variance │ │  No Var  │
                          │ > tol    │ │  → Done  │
                          └──────────┘ └──────────┘
                                │
                                ▼
                          ┌──────────┐
                          │ Re-count │
                          │ (2nd     │
                          │  person) │
                          └──────────┘
                                │
                          ┌─────┴─────┐
                          │           │
                          ▼           ▼
                    ┌──────────┐ ┌──────────┐
                    │  Match   │ │  Differ  │
                    │  → Done  │ │  → 3rd   │
                    └──────────┘ └──────────┘
                                      │
                                      ▼
                                ┌──────────┐
                                │ Manager  │
                                │ Approve  │
                                └──────────┘
                                      │
                                      ▼
                                ┌──────────┐
                                │ Post GL  │
                                │ Adjust   │
                                └──────────┘
```

**Steps:**
1. Manager creates Stock Take for warehouse
2. System assigns keeper from active assignment
3. Keeper opens count sheet
4. If blind count: book quantities hidden
5. Keeper enters physical quantities for each item
6. System calculates variance per item
7. If variance > tolerance (±1 unit or ±0.5%): re-count triggered
8. Second counter (different person) re-counts
9. If still different: manager counts (3rd count)
10. Manager reviews and approves
11. System posts GL adjustment (gain: Dr 152 Cr 632/711, loss: Dr 632/811 Cr 152)
12. Biên bản kiểm kê generated

---

## WF-005: Reconciliation Workflow

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Manager  │────>│  Select  │────>│  System  │
│ requests │     │  Period  │     │  Generate│
│ report   │     │  & WH    │     │  Report  │
└──────────┘     └──────────┘     └──────────┘
                                      │
                                ┌─────┴─────┐
                                │           │
                                ▼           ▼
                          ┌──────────┐ ┌──────────┐
                          │  Var = 0 │ │ Var > 0  │
                          │ "All     │ │ Show     │
                          │  match"  │ │ details  │
                          └──────────┘ └──────────┘
                                      │
                                      ▼
                                ┌──────────┐
                                │  Export  │
                                │  Excel   │
                                └──────────┘
```

**Steps:**
1. Manager selects warehouse and date range
2. System queries StockLedgerEntry (keeper records) and StockBalance (accounting records)
3. System joins on item_id and warehouse_id
4. Calculates variance per item
5. Displays report with highlighting for significant variances
6. Manager can export to Excel

---

## WF-006: Daily Keeper Routine

```
┌─────────────────────────────────────────────────────────────┐
│                    DAILY ROUTINE                            │
├─────────────────────────────────────────────────────────────┤
│ Morning:                                                    │
│   1. Login as Keeper                                        │
│   2. Check pending slips count                              │
│   3. Review new GRN/DN from yesterday                       │
│   4. Record approved slips to Stock Ledger                  │
│                                                             │
│ During Day:                                                 │
│   5. Receive new slips (auto-refresh or manual)             │
│   6. Verify against physical goods (if delivered)           │
│   7. Record or flag discrepancies                           │
│                                                             │
│ End of Day:                                                 │
│   8. Review Stock Ledger balance                            │
│   9. Check for unrecorded slips                             │
│  10. Export daily summary (optional)                        │
│                                                             │
│ Monthly:                                                    │
│  11. Participate in inventory count                         │
│  12. Review reconciliation report                           │
│  13. Sign off on count minutes                              │
└─────────────────────────────────────────────────────────────┘
```
