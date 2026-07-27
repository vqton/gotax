# GoTax COA Module - Use Cases

## Version: 1.0 | Standard: Circular 99/2025/TT-BTC

---

## UC-01: Create New Account

**Actor:** Accountant (Ke toan vien)
**Precondition:** User authenticated with `account:write` permission

### Happy Path
1. User opens Create Account form
2. User enters: code, name, type, parent (optional), detail_by, is_foreign
3. System validates: code format (numeric), uniqueness, parent exists, type valid
4. If account has no balance: status = ACTIVE, no approval needed (for accountants)
5. If account has balance (modification): routes to approval queue (UC-06)
6. System creates account, logs audit entry (CREATE, ACCOUNT, code)
7. System returns created account with 201

**Alternative Path 1a — Code already exists:**
- System returns 409 CONFLICT, error: "account code already exists"

**Alternative Path 1b — Parent does not exist:**
- System returns 400, error: "parent account not found"

**Alternative Path 1c — Parent is not a parent-type account:**
- System returns 400, error: "parent account is marked as detail account"

**Alternative Path 1d — Code violates hierarchy (e.g. 1111 under 112):**
- System returns 400, error: "account code not in valid hierarchy under parent"

**Exception Path 1e — Database error:**
- System returns 500, rolls back, logs error, alerts admin

**Business Rules:**
- Account code must be numeric, min 3 chars
- Type must be one of: ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
- Detail-by: OBJECT, PROJECT, CONTRACT, COST_ITEM, DEPARTMENT, or empty
- Parent must be marked `is_parent = true`

---

## UC-02: Import COA from CSV

**Actor:** Chief Accountant (Ke toan truong) or Admin
**Precondition:** User authenticated with `account:import` permission

### Happy Path
1. User uploads CSV file (columns: code, name, type, parent_code, is_parent, is_active, is_foreign, detail_by)
2. System validates file format (headers match, no encoding issues)
3. System validates each row:
   - Code format, uniqueness (in file + DB), hierarchy depth
   - Parent exists in file or DB
   - Type valid
   - Is_parent flag consistent (parent if has children)
4. System shows preview: X accounts to create, Y to update, Z errors
5. User confirms import
6. System executes in transaction: creates/updates all valid accounts
7. System logs every change in audit trail (IMPORT, ACCOUNT, batch_id)
8. System returns success with summary: created=X, updated=Y, errors=Z

**Alternative Path 2a — File has encoding errors:**
- System detects non-UTF-8, offers encoding conversion (UTF-8, UTF-16, VNI, TCVN3)
- Returns 400 if all formats fail

**Alternative Path 2b — Some rows fail validation:**
- System shows errors per row, user can fix and retry failed rows only
- Valid rows are imported, failed rows are skipped

**Alternative Path 2c — Import old TT200 accounts:**
- User provides mapping file (old_code → new_code)
- System auto-maps TT200 accounts to TT99 equivalents
- System shows mapping preview before import

**Exception Path 2d — Transaction failure during import:**
- System rolls back entire import
- Returns 500 with error details

**Exception Path 2e — File exceeds 10MB:**
- System rejects with 413, error: "file too large, max 10MB"

**Business Rules:**
- Duplicate codes in file: last row wins (with warning)
- Cannot modify account type if account has balance (route to approval)
- Import log retained for 5 years per VSA

---

## UC-03: View Account Balance with Drill-Down

**Actor:** Accountant, Chief Accountant, Auditor
**Precondition:** User authenticated with `account:read` permission

### Happy Path
1. User selects account (from tree or search)
2. User selects period (year + month)
3. User optionally selects currency (VND, USD, etc.)
4. System displays:
   - Account info (code, name, type, parent)
   - Opening balance (debit + credit)
   - Period activity (debit + credit)
   - Closing balance (calculated: open + period)
   - Running balance by day (chart + table)
5. User clicks on period debit amount → system drills down to journal entries
6. User clicks on specific journal entry → system shows journal lines
7. User can export balance view to CSV

**Alternative Path 3a — Account has no activity in period:**
- System shows zero balances, no drill-down available

**Alternative Path 3b — Account is frozen:**
- System shows frozen indicator, balance as of freeze date

**Exception Path 3c — Account not found:**
- System returns 404

**Exception Path 3d — Period not found:**
- System returns 400, error: "period does not exist"

**Business Rules:**
- Balance = opening_debit - opening_credit + period_debit - period_credit
- For ASSET/EXPENSE accounts: normal balance = DEBIT
- For LIABILITY/EQUITY/REVENUE accounts: normal balance = CREDIT
- Running balance resets each fiscal year

---

## UC-04: Export COA

**Actor:** Accountant, Chief Accountant
**Precondition:** User authenticated with `account:export` permission

### Happy Path (CSV)
1. User opens Export COA form
2. User selects format: CSV
3. User optionally filters: active only, by type, by parent
4. User selects columns to include
5. System generates CSV with header row and all accounts
6. System returns CSV file download

### Happy Path (Excel)
1. Same as CSV but format = Excel (.xlsx)
2. System generates formatted Excel:
   - Sheet 1: Full COA (filterable columns)
   - Sheet 2: Account by type (pivot-style)
   - Sheet 3: Account hierarchy (tree view)
3. System returns Excel file download

### Happy Path (PDF — IGAP)
1. User selects PDF with "Official Chart" option
2. System generates PDF:
   - Company header, title "HE THONG TAI KHOAN KE TOAN"
   - Table of accounts grouped by type
   - Account descriptions
   - Issue date, version number
   - Authorized signature block
3. System returns PDF download

**Exception Path 4a — No accounts found matching filter:**
- System returns empty file with header row

**Exception Path 4b — Export takes > 10 seconds:**
- System switches to async mode, queues export, notifies user when ready

---

## UC-05: Account Freeze/Unfreeze

**Actor:** Chief Accountant, Admin
**Precondition:** User authenticated with `account:freeze` permission

### Happy Path (Freeze)
1. User selects account
2. User selects "Freeze Account"
3. User enters reason
4. System validates: account is not already frozen
5. System sets account status = FROZEN
6. System logs audit entry (FREEZE, ACCOUNT, code, reason)
7. System prevents any new journal postings referencing this account
8. All existing balances remain intact
9. Existing posted entries remain visible

### Happy Path (Unfreeze)
1. User selects frozen account
2. User selects "Unfreeze Account"
3. User enters reason
4. System sets account status = ACTIVE
5. System logs audit entry (UNFREEZE, ACCOUNT, code, reason)
6. Account available for new postings

**Exception Path 5a — Account already frozen (on freeze):**
- System returns 409, error: "account is already frozen"

**Exception Path 5b — Account not frozen (on unfreeze):**
- System returns 409, error: "account is not frozen"

**Exception Path 5c — Account has pending approval (on freeze):**
- System returns 409, error: "account has pending changes, complete approval first"

**Business Rules:**
- Freeze does NOT affect historical data
- Freeze is independent of period close
- Freeze reason is mandatory (for audit compliance)

---

## UC-06: Account Change Approval Workflow (4-Eyes)

**Actor:** Chief Accountant (approver), Accountant (requester)
**Precondition:** Requester authenticated; approver authenticated with `account:approve` permission

### Happy Path
1. Accountant requests to create/modify/deactivate account
2. System detects change requires approval (account has balance or type change)
3. System creates approval request: proposed change, reason, requester, timestamp
4. System notifies Chief Accountant
5. Chief Accountant reviews: sees current vs proposed values (diff view)
6. Chief Accountant approves
7. System applies change, logs audit (APPROVE, ACCOUNT, with old/new values)
8. System notifies requester of approval

**Alternative Path 6a — Chief Accountant rejects:**
- System marks approval as REJECTED with rejection reason
- System notifies requester
- No change applied

**Alternative Path 6b — Approval auto-expires (48h):**
- System marks approval as EXPIRED
- System notifies requester to resubmit

**Alternative Path 6c — Requester cancels pending request:**
- System marks approval as CANCELLED
- System notifies approver

**Exception Path 6d — Approver is same as requester:**
- System rejects: "cannot self-approve account changes"

**Exception Path 6e — Account changed by another request while pending:**
- System detects conflict, invalidates pending request, notifies requester

**Business Rules:**
- Changes with zero balance accounts: no approval needed (except type change)
- Type change always requires approval regardless of balance
- Deactivation requires approval if account has balance or is parent
- Approval chain by jurisdiction

---

## UC-07: Compare COA Versions

**Actor:** Chief Accountant, Auditor
**Precondition:** User authenticated with `account:version` permission

### Happy Path
1. User selects "Version History"
2. System shows timeline of COA snapshots with timestamps and change reasons
3. User selects two versions (v1 and v2) to compare
4. System generates diff:
   - Added accounts (in v2, not in v1)
   - Removed accounts (in v1, not in v2)
   - Modified accounts (attributes changed)
   - Side-by-side view with color coding (green=added, red=removed, yellow=modified)
5. User can export diff to PDF

**Alternative Path 7a — No previous version:**
- System shows message: "first COA version, no comparison available"

**Exception Path 7b — Selected versions are identical:**
- System shows: "no differences between selected versions"

---

## UC-08: Generate Internal Accounting Policy (IGAP)

**Actor:** Chief Accountant
**Precondition:** User authenticated with `account:igap` permission

### Happy Path
1. User selects "Generate IGAP"
2. User enters: company name, tax code, address, chief accountant name
3. User selects sections to include:
   - COA full list
   - Account descriptions
   - Accounting methods (per account type)
   - Voucher cycle descriptions
   - Internal control procedures
4. User selects version number
5. System generates PDF with:
   - Cover page with company info
   - Table of contents
   - COA table grouped by type (Loai)
   - Account-level notes (detail_by, is_foreign, etc.)
   - Accounting policy statements
   - Approval and effective date
   - Signature block for Chief Accountant + Legal Representative
6. System returns PDF
7. System saves IGAP as an auditable document version

**Exception Path 8a — Template missing:**
- System returns 500, error: "IGAP template not configured"

---

## UC-09: Link Analysis Code to Account

**Actor:** Accountant
**Precondition:** User authenticated with `account:analysis` permission

### Happy Path
1. User selects account
2. User opens "Analysis Codes" tab
3. User assigns:
   - Cost center (dropdown from cost center list)
   - Profit center (dropdown)
   - Department (from org hierarchy)
   - Project (from project list)
4. System saves analysis code assignment
5. System logs audit (UPDATE, ACCOUNT_ANALYSIS, code)

**Alternative Path 9a — Analysis code does not exist:**
- System provides option to create new analysis code inline

---

## UC-10: Migrate from Old COA (TT200 → TT99)

**Actor:** Admin, Chief Accountant
**Precondition:** System initialized with TT200 COA or user has TT200 export file

### Happy Path
1. User provides old COA (TT200 format) via file upload or DB
2. System loads system mapping table: TT200 account → TT99 account
3. System shows mapping preview:
   - Accounts with direct mapping (same code)
   - Accounts needing mapping (code changed, e.g., 1562 → merged)
   - Accounts needing new code (new TT99 accounts like 215, 332)
   - Removed accounts (611, 631, 441, 461, 466) flagged
4. User reviews and confirms mapping
5. System creates journal entries to zero out old accounts and transfer to new
6. System archives old COA as version "TT200-legacy"
7. System activates new TT99 COA
8. System logs complete migration audit trail

**Exception Path 10a — Old COA has unposted entries:**
- System requires all entries posted before migration

**Exception Path 10b — Mapping ambiguous:**
- System flags for manual review, requires chief accountant decision

---

## Use Case Diagram Summary

```
┌─────────────────────────────────────────────────────────────┐
│                     COA MODULE                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Accountant             Chief Accountant       Admin         │
│     │                       │                    │           │
│     ├── UC-01 Create Acc    │                    │           │
│     ├── UC-03 View Balance──┤                    │           │
│     ├── UC-04 Export COA    │                    │           │
│     ├── UC-09 Analysis Code │                    │           │
│     │                       ├── UC-02 Import COA  │           │
│     │                       ├── UC-05 Freeze Acc  │           │
│     │                       ├── UC-06 Approve     │           │
│     │                       ├── UC-07 Compare     │           │
│     │                       ├── UC-08 Gen IGAP    │           │
│     │                       ├── UC-10 Migrate     ├── UC-10   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```
