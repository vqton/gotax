# User Journeys — Thủ Kho (Warehouse Keeper) Module

**Version:** 1.0
**Date:** August 2026

---

## UJ-001: First-Time Setup — Admin Assigns Keeper

**Persona:** System Administrator (Nguyễn Văn Admin)
**Context:** Company just started using GoTax. Warehouse module is active. Need to assign keeper.

```
Step 1: Login → Dashboard
Step 2: Navigate to Warehouse → Keeper → Assignments
Step 3: Click "Add Assignment"
Step 4: Select warehouse "Kho chính" from dropdown
Step 5: Select user "Trần Thị Bình" (Thủ kho) from dropdown
Step 6: Set Role = "Keeper"
Step 7: Set Effective From = "2026-08-01"
Step 8: Leave Effective To empty (open-ended)
Step 9: Click "Save"
Step 10: System confirms "Assignment created successfully"
Step 11: Trần Thị Bình can now login and see keeper features
```

**Time:** 2 minutes
**Skill level:** Admin

---

## UJ-002: Daily Slip Recording — Keeper Reviews and Records

**Persona:** Warehouse Keeper (Trần Thị Bình)
**Context:** 3 new GRNs from yesterday's deliveries need to be recorded.

```
Step 1: Login as Trần Thị Bình
Step 2: Navigate to Warehouse → Keeper → Pending Slips
Step 3: See 3 pending slips (GRN-202608-001, GRN-202608-002, GRN-202608-003)
Step 4: Click GRN-202608-001 to expand → Review items (ABC-001 x100, ABC-002 x50)
Step 5: Verify quantities match physical delivery (check delivery note)
Step 6: Click checkbox next to all 3 slips (Ctrl+A)
Step 7: Click "Ghi sổ" button
Step 8: System confirms "Record selected slips?"
Step 9: Click "Confirm"
Step 10: System shows "3 slips recorded successfully"
Step 11: Slips move to "Recorded" tab
Step 12: Check Stock Ledger → See new entries with timestamp and your name
```

**Time:** 5 minutes
**Skill level:** Keeper
**Prerequisite:** Keeper assigned to warehouse

---

## UJ-003: Correcting a Mistake — Keeper Un-records

**Persona:** Warehouse Keeper (Trần Thị Bình)
**Context:** Recorded wrong slip (GRN-202608-002 was for different warehouse).

```
Step 1: Navigate to Warehouse → Keeper → Recorded Slips
Step 2: Find GRN-202608-002 in the list
Step 3: Click "Bỏ ghi" (Un-record)
Step 4: System prompts for reason
Step 5: Enter "Recorded wrong warehouse — this GRN is for Kho A, not Kho chính"
Step 6: Click "Confirm"
Step 7: System processes un-recording
Step 8: Slip returns to "Pending" tab
Step 9: Stock Ledger entry marked as unrecorded with reason
```

**Time:** 2 minutes
**Skill level:** Keeper
**Guard:** Cannot un-record if period is closed

---

## UJ-004: Monthly Inventory Count — Keeper Participates

**Persona:** Warehouse Keeper (Trần Thị Bình) + Warehouse Manager (Lê Văn Cường)
**Context:** Monthly inventory count for August 2026.

```
Step 1: Manager creates Stock Take (existing workflow)
Step 2: System assigns Trần Thị Bình as keeper for count
Step 3: Trần Thị Bình logs in → Warehouse → Keeper → Inventory Count
Step 4: Sees count sheet for Kho chính (20 items)
Step 5: Book quantities are hidden (blind count)
Step 6: Physically counts each item, enters quantities
Step 7: Saves draft (can resume later)
Step 8: Submits count when done
Step 9: System calculates variance: Item ABC-001 shows -5 units
Step 10: Manager reviews, triggers re-count for ABC-001
Step 11: Second person re-counts ABC-001 → confirms -5
Step 12: Manager approves count
Step 13: System posts adjustment: Dr 632 Cr 152 (loss)
Step 14: Biên bản kiểm kê generated with 3 signatures
Step 15: Trần Thị Bình, Lê Văn Cường, and Kế toán all sign
```

**Time:** 2-4 hours (depends on item count)
**Skill level:** Keeper + Manager
**Prerequisite:** Stock Take created by manager

---

## UJ-005: Reconciliation Review — Manager Checks Keeper vs Accounting

**Persona:** Warehouse Manager (Lê Văn Cường)
**Context:** Month-end reconciliation between keeper and accounting records.

```
Step 1: Navigate to Warehouse → Keeper → Reconciliation
Step 2: Select warehouse "Kho chính"
Step 3: Select period "08/2026"
Step 4: Click "Generate Report"
Step 5: System shows 15 items with keeper qty vs accounting qty
Step 6: 13 items match (green)
Step 7: Item ABC-001: Keeper 95, Accounting 100, Variance -5 (red)
Step 8: Item XYZ-003: Keeper 200, Accounting 198, Variance +2 (yellow)
Step 9: Click ABC-001 → See detail: difference from unrecorded GRN
Step 10: Click "Export to Excel" → Download report
Step 11: Follow up with Trần Thị Bình to resolve ABC-001 discrepancy
```

**Time:** 10 minutes
**Skill level:** Manager
**Prerequisite:** Keeper has recorded entries for the period

---

## UJ-006: Keeper Views Stock Ledger — End of Day Check

**Persona:** Warehouse Keeper (Trần Thị Bình)
**Context:** End of day, checking stock ledger balance.

```
Step 1: Navigate to Warehouse → Keeper → Stock Ledger
Step 2: Warehouse pre-selected (Kho chính)
Step 3: Date range defaults to today
Step 4: See 5 entries: 3 receipts, 2 issues
Step 5: Check running balance for each item
Step 6: Click "Print" → PDF generated in S06-DN format
Step 7: Review PDF, confirm all entries match physical records
```

**Time:** 5 minutes
**Skill level:** Keeper

---

## UJ-007: Admin Disables Keeper Module

**Persona:** System Administrator
**Context:** Small company where same person handles both keeper and accounting.

```
Step 1: Navigate to System → Configuration → Warehouse Keeper
Step 2: Toggle "Module Enabled" to OFF
Step 3: System warns "Keeper features will be hidden. Data preserved."
Step 4: Click "Confirm"
Step 5: Keeper tabs disappear from Warehouse menu
Step 6: Existing assignment records preserved in database
Step 7: Later: can re-enable anytime
```

**Time:** 1 minute
**Skill level:** Admin
