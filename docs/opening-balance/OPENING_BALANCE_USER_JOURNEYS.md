# Opening Balance Module — User Journeys

**Version:** 1.0

---

## Journey 1: New Enterprise Onboarding (First-time Setup)

**Persona:** Nguyen Van A — Chief Accountant, 15 years experience
**Company:** Cong ty TNHH San Xuat ABC — newly subscribed to GoTax
**Goal:** Set up opening balances correctly so financial statements are accurate

### Scenario
ABC Manufacturing has been using Excel + MISA for 5 years. They switch to GoTax starting August 2026. They need opening balances as of 1 Aug 2026.

### Journey Steps

**Day 1: Prepare Data**
1. Chief Accountant exports trial balance from MISA as of 31 July 2026
2. Verifies: total debit = total credit in old system
3. Identifies accounts needing detail breakdown:
   - Account 112 (Bank): 2 bank accounts (VCB 500M, CTG 300M)
   - Account 131 (Receivables): 50 customers totaling 1.2B
   - Account 152 (Raw materials): 150 items with quantities
   - Account 211 (Fixed assets): 12 assets with cost & depreciation
4. Prepares Excel file per GoTax import template (200 rows)

**Day 2: Import**
1. Logs into GoTax → Opening Balance → Import Excel
2. Uploads file → System shows preview: 195 valid, 5 errors
3. Fixes errors (2 account codes wrong, 3 amounts negative)
4. Re-uploads → Preview: 200 valid, 0 errors, total debit = total credit
5. Confirms import → System creates 200 OBs as DRAFT
6. Reviews all 200 entries in the list view

**Day 3: Approve**
1. Goes to Approval Queue
2. Reviews each OB (spot-checks 20 accounts for accuracy)
3. Batch approves all 200 OBs
4. Runs Opening Balance Report → verifies debit = credit = 2.45B
5. Opens Trial Balance → confirms opening balances appear correctly
6. Opens Balance Sheet (B01-DN) → "Đầu kỳ" column populated

**Outcome:** Company ready for daily accounting in GoTax. All reports accurate.

---

## Journey 2: Fiscal Year Carry-Forward (Year-End)

**Persona:** Tran Thi B — Chief Accountant, 20 years experience
**Company:** Cong ty CP Thuong Mai XYZ
**Goal:** Close FY2026 and carry forward balances to FY2027

### Scenario
Year-end 31 Dec 2026. All 12 periods have been reconciled. Time to close the year.

### Journey Steps

**Phase 1: Pre-Close (Dec 26-30)**
1. Runs Trial Balance for December → verifies debit = credit
2. Runs Balance Sheet → checks all accounts have correct balances
3. Reviews revenue/expense accounts (5xx, 6xx, 7xx, 8xx)
   - Total revenue: 15B
   - Total expenses: 12.5B
   - Net profit: 2.5B (will zero to Account 421)
4. Ensures all journal entries for December are posted
5. Closes Period 12 (December) → status = CLOSED

**Phase 2: Execute (Dec 31)**
1. Opens "Year-End Close" wizard
2. Selects: From 2026 → To 2027
3. System shows preview:
   - Revenue/expense to be zeroed: 35 accounts, 15B
   - Balance sheet to carry forward: 50 accounts, 5B
   - Total debit = total credit: ✓ verified
4. Enters reason: "Annual closing FY2026"
5. Confirms → System executes:
   - Creates 2 closing journal entries (revenue → 421, expense → 421)
   - Posts closing entries automatically
   - Creates 50 opening balances for FY2027 (source: CARRY_FORWARD)
   - All auto-approved
6. System shows success summary:
   - Closing entries: CT-2026-99901, CT-2026-99902
   - Opening balances created: 50 accounts, 5B
   - FY2026 now LOCKED

**Phase 3: Post-Close Verification (Jan 1-2)**
1. Opens FY2027 Opening Balance Report
2. Verifies: Account 1111 opening = 125M (matches Dec 2026 closing)
3. Verifies: Account 421 opening = 3B (includes net profit carry)
4. Opens FY2027 Period 1 → status = OPEN, ready for posting
5. Signs off on year-end close report for board

**Outcome:** Clean year-end close. FY2027 starts with correct opening balances.

---

## Journey 3: Circular 99 Migration (Regulatory Transition)

**Persona:** Le Van C — Chief Accountant, 25 years experience + Admin
**Company:** Cong ty TNHH FDI Dau tu Nuoc ngoai
**Goal:** Migrate from Circular 200 COA to Circular 99 COA by 1 Jan 2026

### Scenario
It's December 2025. Circular 99/2025/TT-BTC takes effect 1 Jan 2026. The company currently uses Circular 200 COA with 150 accounts and 4.5B in balances.

### Journey Steps

**Phase 1: Preparation (Nov 2025)**
1. Studies Circular 99 changes (KPMG, PWC, EY guidance)
2. Downloads GoTax Circular 99 mapping table
3. Identifies affected accounts:
   - Remove: 611 (Purchases), 441 (Capital construction), 466 (FA fund)
   - New: 215 (Biological assets), 332 (Dividends payable)
   - Change: 1562 → merge to 156
4. Creates internal memo for board on COA changes

**Phase 2: Migration (Dec 31, 2025)**
1. Runs final trial balance under Circular 200
2. Opens Circular 99 Migration screen in GoTax
3. System shows mapping preview:
   - Direct: 135 accounts (no change)
   - Remove: 8 accounts (need zeroing)
   - New: 5 accounts (create with balances)
   - Split: 2 accounts (need allocation)
4. Reviews and adjusts mappings:
   - Account 611 balance (500M) → reclassify to 632
   - Account 441 balance (2B) → transfer to 4118
   - Account 466 balance (300M) → transfer to 4118
5. Confirms migration
6. System executes:
   - Creates 10 transfer journal entries
   - Zeroes 8 abolished accounts
   - Creates 5 new accounts with opening balances
   - Generates migration report
7. Verifies: total opening debit = total opening credit = 4.5B

**Phase 3: Validation (Jan 1-15, 2026)**
1. Runs Balance Sheet under new COA
2. Compares: total assets before = total assets after
3. Runs Trial Balance → correct
4. Prepares Circular 99 compliance statement
5. Archives Circular 200 version for audit trail
6. Submits migration report to board + external auditor

**Outcome:** Full Circular 99 compliance by effective date. Zero disruption to operations.

---

## Journey 4: Opening Balance Error Correction

**Persona:** Pham Thi D — Accountant, 5 years experience
**Company:** Any
**Goal:** Fix an opening balance error discovered after approval

### Scenario
Accountant discovers that Account 1121 (Bank VCB) opening balance was entered as 550M instead of the correct 505M — a 45M overstatement.

### Journey Steps

**Step 1: Discovery**
1. Accountant receives bank statement showing actual balance: 505M
2. Compares with GoTax opening balance: 550M
3. Confirms: entry error in initial setup

**Step 2: Correction Request**
1. Opens the approved OB for Account 1121 (status: APPROVED)
2. Clicks "Request Correction"
3. Enters:
   - Current debit: 550M
   - Correct debit: 505M
   - Reason: "Bank statement shows 505M, was entered as 550M"
4. Submits → System creates:
   - Original OB: 550M → status CORRECTED
   - New OB: 505M → status PENDING

**Step 3: Approval**
1. Chief Accountant receives notification
2. Reviews: old 550M vs new 505M, reads reason
3. Verifies against bank statement (scanned copy attached to reason)
4. Approves the correction
5. System updates: New OB → APPROVED, total debit adjusted

**Step 4: Verification**
1. Accountant re-runs Trial Balance
2. Total debit now: 2,405,000,000 (reduced by 45M)
3. Bank account balance matches statement
4. Signs off on correction

**Outcome:** Error corrected with full audit trail. Original value preserved.

---

## Journey 5: Mid-Year Implementation (Company Starts Using GoTax in June)

**Persona:** Hoang Van E — Accountant, 8 years
**Company:** New GoTax customer starting mid-fiscal year

### Scenario
Company starts using GoTax on 1 June 2026. They have 5 months of transactions already posted in their old system.

### Journey Steps

**Step 1: Data Preparation**
1. Exports trial balance from old system as of 31 May 2026
2. Extracts detail breakdowns:
   - Receivables by customer (30 customers)
   - Payables by supplier (15 suppliers)
   - Bank balances (3 bank accounts)
   - Inventory quantities and values
   - Fixed asset register
3. Prepares 5 Excel files (accounts, AR, AP, inventory, FA)

**Step 2: Import + Detail Setup**
1. Imports account balances via Opening Balance Import
2. System creates 40 OB entries (DRAFT)
3. For Account 131 (1.2B total):
   - Adds detail breakdown: 30 customer lines
   - Each line: customer ID, name, amount
4. For Account 152 (500M total):
   - Adds detail breakdown: 150 inventory items
   - Each line: item code, quantity, unit price

**Step 3: Verification**
1. Runs validation: total debit = total credit
2. Compares with old system trial balance
3. Matches: opening balances correct ✓

**Step 4: Approve + Go Live**
1. Chief Accountant approves all OBs
2. System ready for June transactions
3. Accountant starts posting daily entries
4. All reports correct: opening balances + June activity = correct closing

**Outcome:** Smooth mid-year transition. No data loss. Full audit trail.

---

## Journey 6: Audit — Opening Balance Verification

**Persona:** Robert F — External Auditor (PWC Vietnam)
**Goal:** Verify opening balances as part of annual audit

### Scenario
External audit of FY2026 financial statements. Auditor needs to verify opening balances.

### Journey Steps

**Step 1: Request**
1. Auditor requests: Opening Balance Report as of 1 Jan 2026
2. Company exports from GoTax → Opening Balance Report PDF

**Step 2: Verify**
1. Auditor checks: total debit = total credit
2. Traces 5 accounts to prior year audited financial statements:
   - Account 1111: 125M = closing balance FY2025 ✓
   - Account 1121: 500M = closing balance FY2025 ✓
   - Account 421: 3B = retained earnings per prior year audit ✓
3. Checks: no unexplained changes in opening balances

**Step 3: Audit Trail**
1. Requests: Audit Trail for Opening Balance changes
2. System exports: all 200 OB entries with creator, approver, timestamps
3. Auditor verifies:
   - Every OB has creator ≠ approver (4-eyes)
   - No unauthorized modifications
   - Corrections have reason + approval

**Step 4: Sign-Off**
1. Auditor satisfied with opening balance accuracy
2. Issues clean audit opinion
3. Notes: "Opening balances verified. GoTax system controls adequate."

**Outcome:** Clean audit. Opening balances pass external scrutiny.
