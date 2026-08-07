# Use Cases — Thủ Kho (Warehouse Keeper) Module

**Version:** 1.0
**Date:** August 2026

---

## UC-001: Assign Warehouse Keeper

**Actor:** Warehouse Manager / Admin
**Precondition:** User exists, Warehouse exists, no active assignment overlap
**Postcondition:** Keeper assigned to warehouse for effective period

### Happy Path
1. Manager navigates to Warehouse > Keeper Assignment
2. Manager clicks "Add Assignment"
3. System shows form: Warehouse (dropdown), User (dropdown), Role (keeper/manager), Effective From, Effective To
4. Manager fills form, clicks "Save"
5. System validates no overlapping active assignment for same warehouse
6. System creates assignment record
7. System displays success message

### Alternative Paths
4a. Manager leaves Effective To empty → System creates open-ended assignment
4b. Manager assigns same user to multiple warehouses → System allows (user can be keeper for multiple warehouses)

### Exception Paths
5a. Overlapping assignment exists → System rejects with error "Warehouse already has active keeper for this period"
5b. Warehouse is inactive → System rejects with error "Cannot assign to inactive warehouse"
5c. User not found → System rejects with error "User not found"

---

## UC-002: Review and Record Slips (Ghi sổ)

**Actor:** Thủ kho (Warehouse Keeper)
**Precondition:** Keeper assigned to warehouse, pending slips exist
**Postcondition:** Slips recorded in Stock Ledger

### Happy Path
1. Keeper logs in, navigates to Warehouse > Keeper > Pending Slips
2. System shows list of unrecorded slips (GRN, DN, Transfers) for assigned warehouse
3. Keeper reviews each slip (item code, qty, date, reference)
4. Keeper selects one or multiple slips (Ctrl+Click or Ctrl+A)
5. Keeper clicks "Ghi sổ" (Record)
6. System prompts: "Record selected slips to Stock Ledger?"
7. Keeper confirms
8. System creates StockLedgerEntry for each slip with recorded_by = keeper, timestamp
9. Slips move from "Pending" to "Recorded" tab
10. System shows success: "3 slips recorded"

### Alternative Paths
3a. Keeper wants to see details → Clicks slip to expand (shows item lines)
4a. Keeper selects only some lines of a slip → System allows partial recording (records selected lines only)
5a. Keeper clicks "Ghi sổ" with no selection → Button disabled, no action

### Exception Paths
7a. Keeper cancels → No recording, slips remain pending
8a. Recording fails (DB error) → System rolls back, shows error "Recording failed, please try again"
8b. Slip already recorded by another keeper → System skips, shows warning "Slip already recorded by [name]"

---

## UC-003: Un-record Slips (Bỏ ghi)

**Actor:** Thủ kho (Warehouse Keeper)
**Precondition:** Slip has been recorded by this keeper, not yet reconciled
**Postcondition:** Slip un-recorded with reason logged

### Happy Path
1. Keeper navigates to Warehouse > Keeper > Recorded Slips
2. System shows list of recorded slips
3. Keeper selects a slip to un-record
4. Keeper clicks "Bỏ ghi" (Un-record)
5. System prompts for reason (required text field)
6. Keeper enters reason, confirms
7. System marks slip as unrecorded, sets unrecorded_by, unrecorded_at, unrecord_reason
8. Slip moves back to "Pending" tab

### Alternative Paths
5a. Keeper cancels → No un-recording

### Exception Paths
6a. Reason is empty → System rejects "Reason required for un-recording"
6b. Period is already closed → System rejects "Cannot un-record in closed period"
6c. Slip has been reconciled → System rejects "Cannot un-record reconciled slip"

---

## UC-004: View Stock Ledger (Sổ kho)

**Actor:** Thủ kho (Warehouse Keeper)
**Precondition:** Keeper assigned to warehouse
**Postcondition:** Ledger displayed

### Happy Path
1. Keeper navigates to Warehouse > Keeper > Stock Ledger
2. System shows filter: Warehouse (pre-selected), Item (optional), Date Range, Voucher Type
3. Keeper applies filters
4. System loads ledger entries sorted by date
5. Each entry shows: Date, Voucher Type, Voucher No, Description, Receipt Qty, Issue Qty, Running Balance
6. If cost visibility enabled: also shows Unit Cost, Total Value

### Alternative Paths
4a. No entries found → System shows empty state "No ledger entries for selected criteria"
4b. Keeper clicks "Print" → System generates PDF in S06-DN format
4c. Keeper clicks "Export" → System downloads Excel file

### Exception Paths
2a. No warehouse assigned → System shows "No warehouse assignment found. Contact administrator."

---

## UC-005: View Stock Card (Thẻ kho)

**Actor:** Thủ kho (Warehouse Keeper)
**Precondition:** Keeper assigned to warehouse, items exist
**Postcondition:** Stock card displayed

### Happy Path
1. Keeper navigates to Warehouse > Keeper > Stock Card
2. System shows filter: Warehouse, Item (required), Period (month/year)
3. Keeper selects item and period
4. System generates stock card:
   - Header: Item code, name, unit, warehouse
   - Opening balance (from previous period closing)
   - Each receipt with date, voucher no, qty, cumulative balance
   - Each issue with date, voucher no, qty, cumulative balance
   - Closing balance
5. Keeper clicks "Print" → System generates PDF

### Alternative Paths
4a. No transactions in period → System shows "No transactions. Opening and closing balances are zero."
4b. Keeper changes period → System reloads card for new period

### Exception Paths
3a. Item not found → System rejects "Item not found"
3b. Period not valid → System rejects "Invalid period format"

---

## UC-006: Physical Inventory Count (Kiểm kê kho)

**Actor:** Thủ kho + Warehouse Manager/Accountant
**Precondition:** Stock Take initiated by manager, keeper assigned
**Postcondition:** Count completed, variance approved, adjustment posted

### Happy Path
1. Manager creates Stock Take (existing workflow) for warehouse
2. System assigns keeper to count (from assignment)
3. Keeper navigates to Warehouse > Keeper > Inventory Count
4. System shows count sheet with items (book qty hidden if blind count)
5. Keeper enters physical count for each item
6. Keeper submits count
7. System calculates variance per item
8. Accountant reviews count results
9. If variance > tolerance → re-count triggered
10. Manager approves final count
11. System posts adjustment to GL and stock balance
12. Biên bản kiểm kê generated with signatures

### Alternative Paths
5a. Keeper cannot count some items → Keeper enters reason, skips item
6a. Keeper saves draft → System saves partial count, keeper can resume
9a. Re-count matches first count → Variance confirmed, proceed to approval
9b. Re-count differs → Third count by manager

### Exception Paths
4a. Count already in progress by another keeper → System blocks "Count already in progress for this warehouse"
4b. Period closed → System blocks "Cannot count in closed period"
8a. Book qty shown despite blind count setting → System error, escalate

---

## UC-007: View Reconciliation Report

**Actor:** Warehouse Manager / Chief Accountant
**Precondition:** Keeper has recorded entries, StockBalance exists
**Postcondition:** Report displayed

### Happy Path
1. Manager navigates to Warehouse > Keeper > Reconciliation
2. System shows filter: Warehouse, Date Range
3. Manager applies filters
4. System generates report:
   - For each item: Keeper Qty (from StockLedgerEntry balance) vs Accounting Qty (from StockBalance)
   - Variance Qty and Value
5. Manager reviews discrepancies
6. Manager clicks "Export" → Excel download

### Alternative Paths
4a. No variances → System shows "All records match"
4b. Keeper hasn't recorded any slips → System shows "No keeper records found for this period"

### Exception Paths
4a. Large variance (>$10M) → System highlights in red, requires investigation note

---

## UC-008: Toggle Keeper Module

**Actor:** System Administrator
**Precondition:** Admin access
**Postcondition:** Module enabled/disabled

### Happy Path
1. Admin navigates to System > Configuration > Warehouse Keeper
2. System shows current config: Module Enabled (toggle), Cost Price Hidden (toggle)
3. Admin toggles module off
4. System confirms "Keeper module disabled. Keeper tabs will be hidden."
5. Keeper-related tabs disappear from warehouse menu

### Alternative Paths
3a. Admin toggles on → Keeper tabs reappear
3b. Admin changes cost visibility → Takes effect immediately

### Exception Paths
4a. Active assignments exist → System warns "Disabling module will hide keeper features but assignments remain"
