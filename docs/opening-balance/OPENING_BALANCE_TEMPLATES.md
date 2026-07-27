# Opening Balance Module — UI Templates

**Version:** 1.0

---

## 1. Opening Balance Dashboard

```
┌──────────────────────────────────────────────────────────────┐
│ GoTax │ Dashboard │ GL │ Opening Balance                      │
├──────────────────────────────────────────────────────────────┤
│ Opening Balance Management                                    │
│                                                               │
│ Company: [Cong ty TNHH GoTax ▼]   Period: [2026 / 07 ▼]      │
│                                                               │
│ ┌─ Balance Summary ────────────────────────────────────────┐  │
│ │ Total Debit:  2,450,000,000 VND      ✓ BALANCED          │  │
│ │ Total Credit: 2,450,000,000 VND                          │  │
│ │ Accounts:     45 of 189 have balances                     │  │
│ │ Status:       38 APPROVED | 5 PENDING | 2 DRAFT          │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ [+ Add Balance] [📥 Import Excel] [📤 Export] [✓ Validate]   │
│                                                               │
│ ┌─ Opening Balances ──────────────────────────────────────┐  │
│ │ Code     Name                  Type    Debit     Credit  │  │
│ │ ─────────────────────────────────────────────────────── │  │
│ │ ► 1      TAI SAN               ASSET                    │  │
│ │   1111   Tien mat VND          ASSET   125,000,000  -   │  │  ✓
│ │   1121   TGNH VND              ASSET   500,000,000  -   │  │  ✓
│ │   1122   TGNH USD              ASSET   250,000,000  -   │  │  ✓  🔍
│ │   1311   Phai thu KH VND       ASSET   350,000,000  -   │  │  ⏳  🔍
│ │ ► 3      NO PHAI TRA           LIAB                      │  │
│ │   3311   Phai tra NCC VND      LIAB     -    200,000,000│  │  ✓
│ │ ► 4      VON CHU SO HUU        EQUITY                    │  │
│ │   4111   Von dau tu CSH       EQUITY    -  1,000,000,000│  │  ✓
│ │   4211   LN chua phan phoi     EQUITY    -    500,000,000│  │  ⏳
│ │ ► 0      TAI KHOAN NGOAI BANG OFF-BS                    │  │
│ │   001    Tai san thue ngoai    OFF-BS   25,000,000  -   │  │  ✓
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ Status: ✓ APPROVED | ⏳ PENDING | ✏️ DRAFT                    │
│ 🔍 = has detail breakdown                                      │
└──────────────────────────────────────────────────────────────┘
```

---

## 2. Add/Edit Opening Balance Form

```
┌──────────────────────────────────────────────────────────────┐
│ GoTax │ Opening Balance │ Add Balance                         │
├──────────────────────────────────────────────────────────────┤
│ Add Opening Balance                          [Save] [Cancel]  │
│                                                               │
│ ┌─ Balance Information ────────────────────────────────────┐  │
│ │ Account *          [1111 - Tien mat VND ▼]               │  │
│ │                     (leaf accounts only)                  │  │
│ │                                                           │  │
│ │ Currency           [VND ▼]    Exchange Rate: [1.0000]     │  │
│ │                    (auto-filled if foreign)               │  │
│ │                                                           │  │
│ │ Debit Amount       [125,000,000 █]   VND                  │  │
│ │ Credit Amount      [0                ]   VND              │  │
│ │                     (enter debit OR credit)               │  │
│ │                                                           │  │
│ │ Foreign Amount     [             ]   (if foreign currency)│  │
│ │                                                           │  │
│ │ Source             [MANUAL ▼]                             │  │
│ │ Note               [Opening balance from old system   ]   │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ ┌─ Detail Breakdown (optional) ────────────────────────────┐  │
│ │ [Add Detail Line]                                         │  │
│ │                                                           │  │
│ │ Entity Type   Entity ID   Entity Name    Amount    Action │  │
│ │ ─────────────────────────────────────────────────────────  │  │
│ │ CUSTOMER      KH001      Cong ty A      200,000,000  🗑  │  │
│ │ CUSTOMER      KH002      Cong ty B      150,000,000  🗑  │  │
│ │ ─────────────────────────────────────────────────────────  │  │
│ │ Total:                                  350,000,000       │  │
│ │ OB Total:                               350,000,000       │  │
│ │ Difference:                                      0  ✓     │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ ┌─ Account Info (read-only) ──────────────────────────────┐  │
│ │ Code: 1311    Name: Phai thu khach hang                 │  │
│ │ Type: ASSET   Normal Balance: DEBIT                     │  │
│ │ Status: ACTIVE                                          │  │
│ └──────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

---

## 3. Import Preview Screen

```
┌──────────────────────────────────────────────────────────────┐
│ GoTax │ Opening Balance │ Import                              │
├──────────────────────────────────────────────────────────────┤
│ Import Opening Balances from Excel                            │
│                                                               │
│ Step 1: Upload File    [Choose File] ob_template_2026.xlsx    │
│ File type: XLSX ✓    Encoding: [UTF-8 ▼]                     │
│ Rows detected: 156                                           │
│                                                               │
│ ┌─ Import Summary ────────────────────────────────────────┐  │
│ │ Valid rows to import:          153                       │  │
│ │ Rows with errors:               3                        │  │
│ │ New OBs to create:            153                        │  │
│ │ ──────────────────────────────────────────────────────── │  │
│ │ ⚠  Row 12: Account 611 abolished in Circular 99         │  │
│ │ ⚠  Row 45: Amount 0 — skipping zero-row                │  │
│ │ ⚠  Row 78: Account 9999 not found in COA               │  │
│ │                                                          │  │
│ │ [Download Error Report] [Fix and Re-upload]              │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ ┌─ Balance Check ─────────────────────────────────────────┐  │
│ │ Total Debit:  2,450,000,000                               │  │
│ │ Total Credit: 2,450,000,000                               │  │
│ │ Difference:              0  ✓ BALANCED                    │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ ┌─ Preview (first 10 rows) ──────────────────────────────┐  │
│ │ #  Action  Code    Name         Debit      Credit       │  │
│ │ ───────────────────────────────────────────────────────  │  │
│ │ 1  NEW     1111   Tien mat     125,000,000  -          │  │
│ │ 2  NEW     1121   TGNH VND     500,000,000  -          │  │
│ │ 3  NEW     1311   Phai thu KH  350,000,000  -          │  │
│ │ 4  NEW     3311   Phai tra NCC -            200,000,000│  │
│ │ ...                                                      │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ [← Back]                              [✓ Confirm Import]     │
└──────────────────────────────────────────────────────────────┘
```

---

## 4. Approval Queue

```
┌──────────────────────────────────────────────────────────────┐
│ GoTax │ Opening Balance │ Pending Approvals                    │
├──────────────────────────────────────────────────────────────┤
│ Pending Opening Balance Approvals           [12 pending]      │
│                                                               │
│ ┌─────────────────────────────────────────────────────────┐  │
│ │ #  Code    Name             Amount   Requester    Date  │  │
│ │ ───────────────────────────────────────────────────────  │  │
│ │ 1  4211    LN chua PP       500M     Nguyen VA   07/25  │  │  ▶
│ │ 2  1311    Phai thu KH      350M     Tran TH     07/25  │  │  ▶
│ │ 3  2111    TSCD huu hinh    2.5B     Le QH      07/26  │  │  ▶
│ │ 4  3311    Phai tra NCC     200M     Nguyen VA   07/26  │  │  ▶
│ └─────────────────────────────────────────────────────────┘  │
│                                                               │
│ ── Review: #1 Account 4211 ─────────────────────────────     │
│                                                               │
│ Account: 4211 - Loi nhuan chua phan phoi                      │
│ Type: EQUITY | Normal Balance: CREDIT                          │
│ Requested by: Nguyen Van A (Accountant)                       │
│ Reason: "Opening balance from FY2025 carry-forward"            │
│                                                               │
│ Debit: 0 | Credit: 500,000,000 VND                            │
│                                                               │
│ ╔═══════════════════════════════════════════════════════════╗  │
│ ║  [✓ Approve]    [✗ Reject]    [View Account Details]     ║  │
│ ╚═══════════════════════════════════════════════════════════╝  │
│                                                               │
│ Rejection reason (if rejecting): [                      ]    │
│                                                               │
│ [✓] Also approve 3 selected below (batch approve)            │
└──────────────────────────────────────────────────────────────┘
```

---

## 5. Carry-Forward Wizard

```
┌──────────────────────────────────────────────────────────────┐
│ GoTax │ Opening Balance │ Year-End Close                      │
├──────────────────────────────────────────────────────────────┤
│ Fiscal Year Carry-Forward Wizard                              │
│                                                               │
│ Step 1: Select Source and Target                               │
│ ────────────────────────────────────────────────────────────  │
│ From Fiscal Year:  [2026 ▼]   (status: READY FOR CLOSE)      │
│ To Fiscal Year:    [2027 ▼]   (status: FUTURE)                │
│                                                               │
│ Step 2: Review Preview                                        │
│ ────────────────────────────────────────────────────────────  │
│                                                               │
│ ┌─ Closing Summary ───────────────────────────────────────┐  │
│ │ REVENUE & EXPENSE ACCOUNTS                               │  │
│ │ ────────────────────────────────────────────────────────  │  │
│ │ Revenue (Group 5) total:    15,000,000,000               │  │
│ │ Expense (Group 6,7,8) total: 12,500,000,000              │  │
│ │ Net profit to 421:           2,500,000,000               │  │
│ │                                                           │  │
│ │ BALANCE SHEET ACCOUNTS (carried forward)                 │  │
│ │ ────────────────────────────────────────────────────────  │  │
│ │ Total ASSET (1,2):          5,000,000,000                │  │
│ │ Total LIABILITY (3):        2,000,000,000                │  │
│ │ Total EQUITY (4):           3,000,000,000                │  │
│ │                                                           │  │
│ │ CHECK:         5,000,000,000 = 5,000,000,000  ✓          │  │
│ │ (ASSET = LIABILITY + EQUITY)                              │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ ⚠ Warning: After carry-forward, FY2026 will be LOCKED.       │
│ This operation CANNOT be reversed if any entries are         │
│ posted in FY2027.                                             │
│                                                               │
│ Reason for close: [Annual closing - fiscal year 2026      ]  │
│                                                               │
│ [← Back]                    [✓ Execute Carry-Forward]        │
└──────────────────────────────────────────────────────────────┘
```

---

## 6. Detail Breakdown Form

```
┌──────────────────────────────────────────────────────────────┐
│ GoTax │ Opening Balance │ Detail Breakdown                    │
├──────────────────────────────────────────────────────────────┤
│ Account: 1311 - Phai thu khach hang    Total: 350,000,000    │
│                                                               │
│ ┌─ Detail Lines ──────────────────────────────────────────┐  │
│ │ [Add Entity]  [Import from Excel]                        │  │
│ │                                                          │  │
│ │ #  Entity ID  Entity Name     Amount     Aging   Action  │  │
│ │ ───────────────────────────────────────────────────────  │  │
│ │ 1  KH001     Cong ty TNHH A   200,000,000  45d    🗑    │  │
│ │ 2  KH002     Cong ty CP B     100,000,000  30d    🗑    │  │
│ │ 3  KH003     DNTN C            50,000,000  90d    🗑    │  │
│ │ ───────────────────────────────────────────────────────  │  │
│ │ Total:                       350,000,000                 │  │
│ │ OB Total:                    350,000,000                 │  │
│ │ Difference:                           0   ✓             │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ ┌─ Filters ───────────────────────────────────────────────┐  │
│ │ [All Entities ▼] [Search...]                            │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ [Save Details] [Cancel]                                       │
└──────────────────────────────────────────────────────────────┘
```

---

## 7. Circular 99 Migration Screen

```
┌──────────────────────────────────────────────────────────────┐
│ GoTax │ Company Settings │ Circular 99 Migration               │
├──────────────────────────────────────────────────────────────┤
│ Circular 99 Transitional Migration                            │
│                                                               │
│ Current Regime: Circular 200/2014/TT-BTC                      │
│ Target Regime:  Circular 99/2025/TT-BTC                       │
│ Recommendation: Execute at 31-Dec-2025                        │
│                                                               │
│ ┌─ Migration Preview ─────────────────────────────────────┐  │
│ │ Mapping Status: 4,050 accounts to process                 │  │
│ │                                                           │  │
│ │ ┌─ Direct Mapping (3,800 accounts) ────────────────────┐ │  │
│ │ │ Code unchanged, balance transfers directly              │ │  │
│ │ │ Example: 1111 → 1111, 1121 → 1121                    │ │  │
│ │ └──────────────────────────────────────────────────────┘ │  │
│ │                                                           │  │
│ │ ┌─ Removed Accounts (87 accounts) ─────────────────────┐ │  │
│ │ │ ⚠ Must zero out and transfer:                          │ │  │
│ │ │ 441 → 4118 (Capital construction → Other capital)      │ │  │
│ │ │ 466 → 4118 (Fixed asset fund → Other capital)          │ │  │
│ │ │ 611 → 632  (Purchases → COGS)                          │ │  │
│ │ │ 1562 → 156 (Merchandise → merge)                       │ │  │
│ │ │ ...                                                     │ │  │
│ │ └──────────────────────────────────────────────────────┘ │  │
│ │                                                           │  │
│ │ ┌─ New Accounts (163 accounts) ────────────────────────┐ │  │
│ │ │ 🟢 Must create with zero or migrated balance:          │ │  │
│ │ │ 215 - Tai san sinh hoc                                │ │  │
│ │ │ 332 - Phai tra co tuc, loi nhuan                      │ │  │
│ │ │ 82112 - CP thue TNDN toi thieu                        │ │  │
│ │ │ 2295 - DP giam gia TSSH                               │ │  │
│ │ └──────────────────────────────────────────────────────┘ │  │
│ │                                                           │  │
│ │ [View Full Mapping Table] [Edit Mappings]                │  │
│ └──────────────────────────────────────────────────────────┘  │
│                                                               │
│ ╔═══════════════════════════════════════════════════════════╗  │
│ ║  [✗ Cancel]             [Execute Migration]               ║  │
│ ╚═══════════════════════════════════════════════════════════╝  │
└──────────────────────────────────────────────────────────────┘
```

---

## 8. Mobile View (Opening Balance List)

```
┌────────────────────┐
│ GoTax  │ ≡ OB     │
├────────────────────┤
│ 🔍 Search account │
├────────────────────┤
│ Filters: [Period] │
├────────────────────┤
│ Summary: ✓ BALANCED│
│ D: 2.45B C: 2.45B │
├────────────────────┤
│ TAI SAN (ASSET)  ▶ │
│  1111  125M      ✓ │
│  1121  500M      ✓ │
│  1122  250M  ✓ 🔍 │
│  1311  350M   ⏳ 🔍 │
├────────────────────┤
│ NO PHAI TRA ▶     │
│  3311  200M      ✓ │
├────────────────────┤
│ VON CSH     ▶     │
│  4111  1B        ✓ │
│  4211  500M    ⏳  │
├────────────────────┤
│ [+ Add] [📥] [✓]  │
└────────────────────┘
```

---

## 9. Color Scheme & Indicators

| State | Color | Icon | Meaning |
|-------|-------|------|---------|
| DRAFT | Gray (#9ca3af) | ✏️ | Not yet submitted |
| PENDING | Blue (#3b82f6) | ⏳ | Awaiting approval |
| APPROVED | Green (#22c55e) | ✓ | Active in reports |
| REJECTED | Red (#ef4444) | ✗ | Needs correction |
| CORRECTED | Amber (#f59e0b) | 🔄 | Superseded by new version |
| Has detail | — | 🔍 | Click to see breakdown |
| Balanced | Green | ✓ | Total debit = total credit |
| Unbalanced | Red | ✗ | Total debit ≠ total credit |
