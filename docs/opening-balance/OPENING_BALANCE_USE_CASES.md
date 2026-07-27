# Opening Balance Module — Use Cases

**Version:** 1.0 | **Standard:** Circular 99/2025/TT-BTC

---

## UC-OB-01: Enter Opening Balance for Single Account

**Actor:** Accountant
**Precondition:** Company configured, fiscal year created, period open, account exists

### Happy Path
1. Accountant opens "Opening Balance" screen
2. Selects company, fiscal year, period (auto-default to current open period)
3. Clicks "Add Opening Balance"
4. Selects account from COA tree (leaf accounts only)
5. Enters: debit amount OR credit amount (not both), currency (default VND)
6. System validates: account active, not parent, period open, unique per period
7. System saves OB as DRAFT
8. Message: "Opening balance saved as draft. Submit for approval."

**Alternative Path 1a — Single line with detail breakdown:**
1-5. Same as happy path
6. Clicks "Add Detail Breakdown"
7. Selects entity type (Customer/Supplier/Bank/etc.)
8. Selects entity from dropdown
9. Enters allocated amounts per entity (sum must = total)
10. System validates detail sum = total
11. Saves with detail lines

**Alternative Path 1b — Foreign currency:**
1-4. Same as happy path
5. Selects currency (e.g., USD)
6. System auto-fills exchange rate from ExchangeRate table (last rate before OB date)
7. User confirms or adjusts rate
8. Enters foreign amount
9. System computes VND equivalent
10. Saves

**Exception Path 1c — Period closed/locked:**
- System returns error: "period not open"

**Exception Path 1d — Account has approved OB for same period:**
- System returns error: "opening balance already exists for this account"

**Exception Path 1e — Account is parent type:**
- System returns error: "parent account cannot have opening balance"

**Business Rules:**
- OB-001: One approved OB per account per period per currency
- OB-002: Debit XOR Credit must be > 0
- OB-003: Leaf accounts only (is_parent = false)
- OB-004: Period must be in OPEN status

---

## UC-OB-02: Import Opening Balances from Excel

**Actor:** Accountant, Chief Accountant
**Precondition:** Excel file prepared per template

### Happy Path
1. User opens "Import Opening Balances"
2. Uploads Excel file (.xlsx or .csv)
3. System detects encoding (UTF-8/VNI/TCVN3)
4. System parses and validates each row:
   - Account exists, is leaf, is active
   - Amount format valid
   - Currency valid
   - No duplicate account in file
5. System shows preview: X valid rows, Y errors, total debit = total credit
6. If unbalanced: system highlights difference in red, user can:
   a. Go back and fix file
   b. Add VRR account 999 to balance (with warning)
7. User confirms import
8. System executes in transaction
9. System creates all OBs as DRAFT
10. System shows result: imported=X, errors=Y, batch_id

**Alternative Path 2a — All rows invalid:**
- System rejects import, shows error table with row numbers

**Alternative Path 2b — Partial success:**
- System imports valid rows, creates error file for failed rows
- User can download error report

**Exception Path 2c — File exceeds 10MB:**
- System returns 413, error: "file too large"

**Exception Path 2d — Wrong template:**
- System detects missing required columns, returns format error

**Business Rules:**
- OB-005: Import is transactional within each batch; no partial DB writes
- OB-006: Max 10,000 rows per import
- OB-007: File max 10MB
- OB-008: Import creates DRAFT status entries (approval required separately)

---

## UC-OB-03: Approve Opening Balance (4-Eyes)

**Actor:** Chief Accountant
**Precondition:** OB exists in PENDING status, approver ≠ creator

### Happy Path
1. Chief Accountant opens "Pending Approvals" queue
2. Sees list of OBs pending approval with: account, amount, creator, date
3. Clicks on one to review
4. System shows: account info, current OB, account type, normal balance indicator
5. If detail breakdown exists: shows entity breakdown with totals
6. Chief Accountant verifies accuracy
7. Clicks "Approve"
8. System: status → APPROVED, records approver + timestamp
9. Audit log entry created (APPROVE, OPENING_BALANCE, id)
10. OB now active for reporting

**Alternative Path 3a — Chief Accountant rejects:**
1-5. Same
6. Clicks "Reject"
7. Enters rejection reason (required)
8. System sets status → REJECTED, notifies creator
9. Creator can edit and resubmit

**Alternative Path 3b — Batch approval:**
1. Chief Accountant selects multiple OBs
2. Clicks "Approve Selected"
3. System approves all in transaction
4. Any failure: rollback entire batch

**Exception Path 3c — Self-approval attempted:**
- System rejects: "cannot self-approve opening balances"

**Exception Path 3d — OB already approved:**
- System returns error: "opening balance already approved"

**Business Rules:**
- OB-009: Approver must be different from creator (4-eyes)
- OB-010: Approved OB becomes immutable (no edit, only correction)
- OB-011: Batch approval atomic — all or nothing

---

## UC-OB-04: Correct Approved Opening Balance

**Actor:** Accountant (initiates), Chief Accountant (approves)
**Precondition:** OB exists in APPROVED status

### Happy Path
1. Accountant opens approved OB
2. Clicks "Request Correction"
3. Enters: new debit/credit amounts, reason (required)
4. System validates: reason not empty, new amounts valid
5. System creates correction request:
   - Original OB: status → CORRECTED, marked correction_of = new ID
   - New OB: status → PENDING, linked to original via correction_of
6. Notification sent to Chief Accountant
7. Chief Accountant reviews and approves (see UC-OB-03)
8. Reports now use corrected OB

**Alternative Path 4a — Correction rejected:**
- Original OB stays APPROVED
- Correction request deleted
- Creator notified with reason

**Exception Path 4b — Correction during period close:**
- System blocks: "period is closing, cannot modify opening balances"

**Business Rules:**
- OB-012: Correction preserves full history (no overwrite)
- OB-013: CORRECTED status retains original values for audit
- OB-014: Reason field mandatory for corrections

---

## UC-OB-05: Execute Fiscal Year Carry-Forward

**Actor:** Chief Accountant (or automatic system process)
**Precondition:** Previous fiscal year periods all CLOSED, new year period exists in FUTURE/OPEN

### Happy Path
1. Chief Accountant opens "Year-End Close" wizard
2. Selects: from fiscal year (e.g., 2026) → to fiscal year (e.g., 2027)
3. System analyzes all accounts with closing balances
4. System shows preview:
   - Revenue/expense to be zeroed: X accounts, total Y VND
   - Balance sheet to carry forward: Z accounts
   - Total debit = total credit: verified
5. User confirms
6. System executes:
   a. Creates closing journal entries (zero rev/exp → 421)
   b. Posts closing entries
   c. Creates opening balances for new year (source = CARRY_FORWARD)
   d. Sets new year OBs to APPROVED (auto)
   e. Logs CarryForwardLog
7. System shows success summary
8. New year period now has opening balances ready

**Alternative Path 5a — Partial year (e.g., company starts mid-year):**
- No carry-forward from prior year
- All OBs entered manually or imported

**Exception Path 5b — Prior year not fully closed:**
- System blocks: "close all periods in fiscal year first"

**Exception Path 5c — Carry-forward already executed:**
- System blocks: "carry-forward already completed for this period pair"

**Business Rules:**
- OB-015: Carry-forward is a one-time operation per fiscal year pair
- OB-016: After carry-forward, previous year is LOCKED (no reopening)
- OB-017: Revenue/expense zeroed to Retained Earnings (Account 421)
- OB-018: Off-balance-sheet accounts (00x) carried forward separately

---

## UC-OB-06: Circular 99 Transitional Migration

**Actor:** Chief Accountant, Admin
**Precondition:** Company using Circular 200 accounts, need to migrate to Circular 99

### Happy Path
1. Admin selects "Circular 99 Migration" from Company settings
2. System shows current COA (TT200) with balances
3. System loads mapping table (TT200 → TT99)
4. System shows migration preview:
   - Direct mapping: accounts with same code (e.g., 1111 → 1111)
   - Remove: accounts abolished in TT99 (611, 441, 466, etc.)
   - New accounts: introduced in TT99 (215, 332, etc.)
   - Split: accounts with multiple destinations
5. Chief Accountant reviews and adjusts mappings as needed
6. Confirms migration
7. System executes:
   - Creates transfer journal entries per mapping rules
   - Zeroes abolished accounts
   - Creates opening balances for new TT99 accounts
   - Generates BalanceMigration record
   - Creates audit trail
8. System shows: total migrated, direct, split/merged, manual required

**Alternative Path 6a — Split mapping (e.g., 1562 → merge to 156):**
- System creates 2 lines: debit 1562, credit 156
- Creates journal entry for the transfer

**Exception Path 6b — Unmapped accounts:**
- System flags for manual mapping
- Migration paused until all accounts mapped

**Exception Path 6c — Period not at year-end:**
- System warns: "migration should occur at 31-Dec boundary"
- Allows override with chief accountant approval

**Business Rules:**
- OB-019: Migration date must be 31-Dec-2025 per Circular 99 Art 30
- OB-020: All accounts must be mapped before execution
- OB-021: Migration creates auditable journal entries
- OB-022: Migration is irreversible (support rollback only via reversal entries)

---

## UC-OB-07: View Opening Balance Report

**Actor:** Accountant, Chief Accountant, Auditor
**Precondition:** Company has approved opening balances for selected period

### Happy Path
1. User navigates to Reports → "Opening Balance"
2. Selects company, fiscal year, period
3. System displays:
   - Summary: total debit, total credit, difference, balanced status
   - Table: account code, name, type, normal balance, debit, credit, detail count
   - Grouped by account type (ASSET, LIABILITY, etc.)
4. User can:
   - Filter by account type
   - Search by account code/name
   - Expand to see detail breakdown per account
   - Export to Excel/PDF
   - View audit trail

**Alternative Path 7a — Period has no OBs:**
- Shows empty state: "no opening balances entered for this period"

**Business Rules:**
- OB-023: Only APPROVED OBs appear in reports
- OB-024: Report must balance: total debit = total credit

---

## UC-OB-08: Set Opening Balance Detail Breakdown

**Actor:** Accountant
**Precondition:** OB exists in DRAFT or PENDING status

### Happy Path
1. User opens an OB in DRAFT status
2. Clicks "Detail Breakdown" tab
3. Clicks "Add Detail Line"
4. Selects entity type (e.g., CUSTOMER, SUPPLIER, EMPLOYEE)
5. Selects entity from searchable dropdown
6. Enters: allocated debit/credit amount
7. For inventory: enters quantity + unit price
8. For FA: enters original cost + accumulated depreciation
9. System validates: sum of detail lines = OB total
10. Saves detail lines

**Alternative Path 8a — Detail sum not equal to total:**
- System shows error: "detail sum (X) must equal OB total (Y)"
- User adjusts amounts

**Business Rules:**
- OB-025: Detail sum must exactly equal OB total
- OB-026: Entity must exist in company's master data
- OB-027: Detail lines are immutable after OB approval

---

## UC-OB-09: Rollback Carry-Forward

**Actor:** Chief Accountant (emergency only)
**Precondition:** Carry-forward executed, but new year has ZERO entries posted

### Happy Path
1. Chief Accountant opens "Carry-Forward History"
2. Selects the executed carry-forward
3. Clicks "Rollback"
4. System warns: "This will delete all opening balances created by carry-forward"
5. Chief Accountant confirms with reason
6. System:
   - Deletes OBs with source=CARRY_FORWARD for the target period
   - Reverses closing journal entries
   - Updates CarryForwardLog status → ROLLED_BACK
7. Periods restored to previous state

**Exception Path 9a — Entries posted in new year:**
- System blocks: "cannot rollback, transactions exist in target period"

**Business Rules:**
- OB-028: Rollback only allowed if new year has no posted entries
- OB-029: Rollback requires admin-level authorization

---

## Use Case Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    OPENING BALANCE MODULE                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Accountant              Chief Accountant            Admin        │
│     │                          │                      │           │
│     ├── UC-01 Enter OB         │                      │           │
│     ├── UC-02 Import Excel     │                      │           │
│     ├── UC-04 Request Correct  │                      │           │
│     ├── UC-08 Set Detail       │                      │           │
│     ├── UC-07 View Report      │                      │           │
│     │                          │                      │           │
│     │                          ├── UC-03 Approve OB   │           │
│     │                          ├── UC-05 Carry-Forward│           │
│     │                          ├── UC-06 C99 Migration│           │
│     │                          ├── UC-09 Rollback     │           │
│     │                          ├── UC-07 View Report  │           │
│     │                          │                      │           │
│     │                          │                      ├── UC-06  │
│     │                          │                      ├── UC-09  │
│ ────┴─── AUDITOR ──────────────┴──────────────────────┴───       │
│     │                                                             │
│     └── UC-07 View Report (read-only)                             │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```
