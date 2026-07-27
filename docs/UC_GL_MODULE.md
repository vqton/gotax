# GL Module Use Cases

## UC-01: Account Management

### UC-01.01: Create Account
| Field | Value |
|-------|-------|
| **Actor** | Chief Accountant, Admin |
| **Precondition** | Authenticated with CREATE_ACCOUNT permission |
| **Postcondition** | Account created and visible in COA |

#### Happy Path
1. User navigates to Chart of Accounts
2. Clicks "Create Account"
3. Enters: Code, Name, Type, Parent (optional), IsActive, DetailBy
4. System validates: code numeric, ≥3 chars, unique, type valid
5. System creates account → returns 201 Created
6. Account appears in COA listing

#### Alternative Path: Create Child Account
1. User selects parent account → clicks "Add Child"
2. Parent code auto-filled
3. Same validation as Happy Path
4. Account created with parent relationship

#### Exception Path: Duplicate Code
1. User enters code that already exists
2. Returns 409 Conflict: "Account code already exists"
3. Form retains all entered data
4. User must change code

#### Exception Path: Invalid Code
1. User enters non-numeric code or <3 chars
2. Returns 400 Bad Request with validation message
3. Form retains valid fields, highlights invalid

---

### UC-01.02: Update Account
| Field | Value |
|-------|-------|
| **Actor** | Chief Accountant, Admin |
| **Precondition** | Account exists |
| **Postcondition** | Account details updated |

#### Happy Path
1. User searches account → selects from results
2. Modifies fields (Name, IsActive, DetailBy, etc.)
3. Code NOT editable (immutable after creation)
4. Submits → 200 OK
5. Updated data reflected immediately

#### Exception: Account Not Found
1. User attempts to update non-existent code
2. Returns 404 Not Found

#### Exception: Code Mismatch
1. URL code ≠ body code
2. Returns 400 Bad Request

---

### UC-01.03: Delete Account
| Field | Value |
|-------|-------|
| **Actor** | Admin only |
| **Precondition** | Account exists, no children, zero balance |
| **Postcondition** | Account removed from COA |

#### Happy Path
1. User confirms deletion of leaf account (no children)
2. System verifies zero balance in all periods
3. Account deleted → 200 OK

#### Exception: Account Has Children
1. User attempts to delete parent account
2. Returns 409 Conflict: "Account has children"
3. User must delete children first or move them

#### Exception: Account Has Balance
1. Account has non-zero balance in any period
2. Returns 409 Conflict: "Account has balance"
3. Account must be zeroed before deletion

---

## UC-02: Journal Entry

### UC-02.01: Post Journal Entry
| Field | Value |
|-------|-------|
| **Actor** | Accountant |
| **Precondition** | At least 2 active accounts exist, period is OPEN |
| **Postcondition** | Entry POSTED, balances updated |

#### Happy Path
1. User creates new voucher → selects type (Thu/Chi/BH/MH/etc.)
2. Enters: Date, Description, Period, Lines (account, debit, credit, description)
3. System validates:
   - All accounts exist and are ACTIVE
   - Total Debit = Total Credit
   - Date within OPEN period
   - No account with detail-by constraint violated
4. Auto-generates voucher number (prefix + YYYY + sequential)
5. Entry saved as DRAFT → submitted for approval → APPROVED → POSTED
6. Account balances updated
7. Returns 201 Created with entry ID

#### Alternative Path: Save as Draft
1. User enters partial data
2. Saves as DRAFT (not submitted)
3. Draft can be retrieved, edited, and deleted later

#### Alternative Path: Multi-Currency Voucher
1. User selects foreign currency account
2. System prompts for: currency code, exchange rate, foreign amount
3. Auto-calculates VND equivalent
4. Balance check in both currencies

#### Exception: Inactive Account
1. User enters inactive account code
2. Returns 400: "Account is inactive"
3. User must select active account

#### Exception: Unbalanced Entry
1. Total Debit ≠ Total Credit
2. Returns 400: "Total debit must equal total credit"
3. Shows current totals: "Debit: 10,000,000 / Credit: 9,500,000"
4. User must fix discrepancy

#### Exception: Period Closed
1. Entry date falls in CLOSED period
2. Returns 400: "Period is closed"
3. User must change date or reopen period

#### Exception: Account Unknown
1. User enters non-existent account code
2. Returns 400: "Account not found"
3. User must select valid account

---

### UC-02.02: Cancel Entry
| Field | Value |
|-------|-------|
| **Actor** | Chief Accountant, Admin |
| **Precondition** | Entry exists and is POSTED |
| **Postcondition** | Entry status = CANCELLED |

#### Happy Path
1. User searches for posted entry
2. Clicks "Cancel"
3. Confirms cancellation (modal: "This cannot be undone")
4. System validates: entry is POSTED
5. Status changes to CANCELLED
6. User prompted to create reversal entry

#### Exception: Already Cancelled
1. User attempts to cancel already-cancelled entry
2. Returns 409 Conflict

#### Exception: Draft Entry
1. User attempts to cancel DRAFT entry
2. Returns 400: "Draft entry — delete instead of cancel"

---

## UC-03: Trial Balance

### UC-03.01: View Trial Balance
| Field | Value |
|-------|-------|
| **Actor** | All authenticated users |
| **Precondition** | Period exists |
| **Postcondition** | Report displayed |

#### Happy Path
1. User selects year and month
2. System fetches period → validates period exists
3. Aggregates all POSTED entries for period
4. Calculates: Opening Balance + Period Activity = Closing Balance
5. Computes: Total Debits = Total Credits (cross-check)
6. Displays table: Account Code, Name, Opening(Dr/Cr), Period(Dr/Cr), Closing(Dr/Cr)

#### Exception: Period Not Found
1. User selects year/month with no period created
2. Returns 404: "Period not found"

---

## UC-04: Period Management

### UC-04.01: Close Period
| Field | Value |
|-------|-------|
| **Actor** | Chief Accountant, Admin |
| **Precondition** | Period is OPEN, all entries posted |
| **Postcondition** | Period = CLOSED, entries blocked |

#### Happy Path
1. User initiates period close for month N
2. System verifies: no unposted entries in period
3. System changes status: OPEN → CLOSING
4. User runs closing entries (511→911, 911→421)
5. System generates trial balance → verifies zero balance on revenue/expense
6. System generates B01-B04 financial statements
7. Admin reviews and confirms
8. Period → CLOSED
9. Opening balances for period N+1 calculated

#### Exception: Period Already Closed
1. User attempts to close already-closed period
2. Returns 409: "Period already closed"

#### Exception: Unbalanced Entries
1. System finds unbalanced entries during verification
2. Returns list of unbalanced entries
3. User must fix before closing

---

## UC-05: Audit Trail

### UC-05.01: View Audit Log
| Field | Value |
|-------|-------|
| **Actor** | Admin, Auditor |
| **Precondition** | Authenticated |
| **Postcondition** | Log displayed |

#### Happy Path
1. User navigates to Audit Log
2. Filters by: date range, user, action type, entity type
3. System displays chronological log:
   - Timestamp | User | IP | Action | Entity | Old Value | New Value
4. Exportable to Excel/PDF

#### Alternative: Drill-Down
1. From Trial Balance or GL, user clicks "View History"
2. Shows all changes to that specific account/entry
3. Full before/after snapshot

---

## UC-06: Financial Statements

### UC-06.01: Generate B01 — Statement of Financial Position
| Field | Value |
|-------|-------|
| **Actor** | Chief Accountant |
| **Precondition** | Period closed, trial balance verified |
| **Postcondition** | B01-DN generated |

#### Happy Path
1. User selects period (monthly/quarterly/annual)
2. System maps trial balance accounts to B01 line items:
   - ASSETS: Current (100-199) + Non-current (200-299)
   - LIABILITIES: Current (300-399) + Non-current
   - EQUITY: 400-499
3. Calculates: Total Assets = Total Liabilities + Equity
4. Displays B01-DN in Circular 99 format
5. Comparative: current period vs prior period
6. Export to PDF/Excel with Circular 99 template

### UC-06.02: Generate B02 — Income Statement
1. User selects period
2. System maps: Revenue (500-599) - Expenses (600-899) = Profit/Loss
3. Displays by line item per Circular 99

### UC-06.03: Generate B03 — Cash Flow Statement
1. User selects direct or indirect method
2. System computes:
   - Operating activities (indirect: from P&L + working capital changes)
   - Investing activities (fixed asset purchases/sales)
   - Financing activities (loans, dividends, equity changes)

---

## UC-07: Multi-Currency

### UC-07.01: Define Exchange Rate
| Field | Value |
|-------|-------|
| **Actor** | Chief Accountant |
| **Precondition** | Authenticated |
| **Postcondition** | Rate stored for period |

#### Happy Path
1. User navigates to Exchange Rates
2. Selects: Currency, Date, Rate (buy/sell/average)
3. Rate applied to all transactions on that date
4. Auto-calculate VND equivalent for foreign entries

### UC-07.02: Year-End Revaluation
1. User initiates revaluation for period
2. System identifies all foreign currency balances
3. Applies period-end exchange rate
4. Records exchange rate difference (Account 413/635/515)
5. Generates revaluation report

---

## Appendix: State Machines

### Journal Entry States
```
                  ┌─────────┐
                  │  DRAFT  │
                  └────┬────┘
                       │ submit
                       ▼
                 ┌───────────┐
           ┌─────│ REVIEWING │─────┐
           │     └─────┬─────┘     │
           │ approve   │           │ reject
           ▼           │           ▼
     ┌────────┐        │     ┌──────────┐
     │APPROVED│        │     │  DRAFT   │ (returned for fix)
     └───┬────┘        │     └──────────┘
         │ post        │
         ▼             │
    ┌──────────┐       │
    │  POSTED  │       │
    └────┬─────┘       │
         │ cancel      │
         ▼             │
    ┌───────────┐      │
    │ CANCELLED │      │
    └───────────┘      │
                       │ delete (DRAFT only)
                       ▼
                    (removed)
```

### Period States
```
     ┌───────┐
     │ OPEN  │────→ CLOSING ────→ CLOSED
     └───────┘        │              │
          ▲           │              │
          │           ▼              │
          └───── REOPEN ─────────────┘
                     (audit trail logged)
```

### Account Status
```
     ┌──────────┐
     │  ACTIVE  │ ←──→ INACTIVE (toggle)
     └──────────┘
           │
           ▼
       (delete only if: no children, zero balance)
```
