# GoTax COA Module - UI Templates

---

## 1. COA List View (Account Tree)

```
 ┌─────────────────────────────────────────────────────────────┐
 │ GoTax  │  Dashboard  │  COA ▼  │  Journal  │  Reports      │
 ├─────────────────────────────────────────────────────────────┤
 │ COA Management                                              │
 │                                                              │
 │ ┌─ Search box................. [ 🔍 ] [+ Add Account] [📥]  │
 │ │                                                           │
 │ │ Filters: [All Types ▼] [Active Only ☑] [Parent ▼]        │
 │ │                                                           │
 │ │ ┌───────────────────────────────────────────────────┐    │
 │ │ │ Code    Name                          Type   Status│    │
 │ │ │ ───────────────────────────────────────────────── │    │
 │ │ │ ► 1     TAI SAN (ASSETS)              ASSET   ✓   │    │
 │ │ │   ► 111  Tien mat                      ASSET   ✓   │    │
 │ │ │       1111 Tien mat VND               ASSET   ✓   │    │
 │ │ │       1112 Tien mat USD               ASSET   ✓   │    │
 │ │ │   ► 112  Tien gui khong ky han        ASSET   ✓   │    │
 │ │ │       1121 Tien gui VND               ASSET   ✓   │    │
 │ │ │       1122 Tien gui USD               ASSET   ✓   │    │
 │ │ │   ► 131  Phai thu khach hang          ASSET   🔒   │    │
 │ │ │ ► 3     NO PHAI TRA (LIABILITIES)     LIAB    ✓   │    │
 │ │ │ ► 4     VON CHU SO HUU (EQUITY)       EQUITY  ✓   │    │
 │ │ │ ► 5     DOANH THU (REVENUE)           REV     ✓   │    │
 │ │ │ ► 6-9   CHI PHI (EXPENSES)            EXP     ✓   │    │
 │ │ └───────────────────────────────────────────────────┘    │
 │ │                                                           │
 │ │ [📥 Import] [📤 Export ▼] [📋 Versions] [📄 IGAP]       │
 │ │                                                           │
 │ └─ Status: 220 accounts loaded | v2.1 | 71 L1 accounts ──┘ │
 └─────────────────────────────────────────────────────────────┘
```

---

## 2. Account Detail / Create Form

```
 ┌─────────────────────────────────────────────────────────────┐
 │ GoTax  │  COA  │  Create Account                            │
 ├─────────────────────────────────────────────────────────────┤
 │ Create New Account                           [Save] [Cancel]│
 │                                                              │
 │ ┌─ Account Information ──────────────────────────────────┐  │
 │ │                                                        │  │
 │ │ Account Code *     [ 1113    ] (numeric, min 3 digits) │  │
 │ │ Account Name VN *  [Tien mat EUR                   ]   │  │
 │ │ Account Name EN    [Cash in EUR                    ]   │  │
 │ │                                                        │  │
 │ │ Account Type *     [ASSET ▼]                           │  │
 │ │ Parent Account     [111 - Tien mat ▼]                  │  │
 │ │                                                        │  │
 │ │ Is Parent Account  [☐] (enables child accounts)        │  │
 │ │ Is Active          [☑]                                 │  │
 │ │ Is Foreign Currency [☑] (EUR, USD)                     │  │
 │ │                                                        │  │
 │ │ Detail By          [OBJECT ▼]                          │  │
 │ │                     None / OBJECT / PROJECT / CONTRACT  │  │
 │ │                     COST_ITEM / DEPARTMENT              │  │
 │ │                                                        │  │
 │ │ Freeze Reason      [                              ]    │  │
 │ │ (only when freezing)                                    │  │
 │ │                                                        │  │
 │ │ Note               [                              ]    │  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ ┌─ Analysis Codes ───────────────────────────────────────┐  │
 │ │ Cost Center       [▼ Select or create new...         ]  │  │
 │ │ Profit Center     [▼ Select or create new...         ]  │  │
 │ │ Department        [▼ Select or create new...         ]  │  │
 │ │ Project           [▼ Select or create new...         ]  │  │
 │ └────────────────────────────────────────────────────────┘  │
 └─────────────────────────────────────────────────────────────┘
```

---

## 3. Account Balance Inquiry

```
 ┌─────────────────────────────────────────────────────────────┐
 │ GoTax  │  COA  │  Account Balance                           │
 ├─────────────────────────────────────────────────────────────┤
 │ Account: 1111 - Tien mat VND       Type: ASSET              │
 │ Status: ACTIVE                      Parent: 111 - Tien mat  │
 ├─────────────────────────────────────────────────────────────┤
 │ Period: [2026 ▼] / [07 ▼]  Currency: [All ▼]               │
 │                                                              │
 │ ┌─ Balance Summary ──────────────────────────────────────┐  │
 │ │                                            Amount (VND) │  │
 │ │ Opening Balance Debit                   125,000,000     │  │
 │ │ Opening Balance Credit                          0       │  │
 │ │ Period Debit                              32,500,000    │  │
 │ │ Period Credit                             18,200,000    │  │
 │ │ ─────────────────────────────────────────────────────── │  │
 │ │ Closing Balance Debit                     139,300,000   │  │
 │ │ Closing Balance Credit                            0     │  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ ┌─ Period Activity ──────────────────────────────────────┐  │
 │ │ Date       Voucher   Description          Debit   Credit│  │
 │ │ ──────────────────────────────────────────────────────── │  │
 │ │ 01/07/2026  PT-0001  Opening balance     125M     -     │  │
 │ │ 05/07/2026  TH-0105  Customer payment     -      5.2M  │  │
 │ │ 12/07/2026  TH-0123  Customer payment     -     12.0M  │  │
 │ │ 15/07/2026  CH-0098  Office supplies     3.5M    -     │  │  ▶ click
 │ │ 20/07/2026  TH-0150  Customer payment     -      1.0M  │  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ [📥 Export CSV] [📊 Chart View]                             │
 └─────────────────────────────────────────────────────────────┘

 ┌─ Drill-down: Journal CH-0098 ───────────────────────────┐
 │ Entry: CH-0098 | Date: 15/07/2026 | Type: CHI (Payment) │
 │ Description: Office supplies purchase                    │
 │ ──────────────────────────────────────────────────────── │
 │ Line  Account     Description          Debit    Credit   │
 │ ──────────────────────────────────────────────────────── │
 │ 1     1111       Payment to supplier   3,500,000  -      │
 │ 2     3311       Supplier payable       -       3,500,000│
 └─────────────────────────────────────────────────────────┘
```

---

## 4. COA Import Preview

```
 ┌─────────────────────────────────────────────────────────────┐
 │ GoTax  │  COA  │  Import                                    │
 ├─────────────────────────────────────────────────────────────┤
 │ Import Chart of Accounts                                    │
 │                                                              │
 │ Step 1: Upload File    [Choose File] sample_coa_2026.csv    │
 │ File type: CSV                                          ✓   │
 │ Encoding: [UTF-8 ▼]  |  Rows detected: 156                 │
 │                                                              │
 │ Step 2: Review Preview                                      │
 │                                                              │
 │ ┌─ Import Summary ──────────────────────────────────────┐  │
 │ │ New accounts to create:    42                           │  │
 │ │ Existing accounts to update: 114                        │  │
 │ │ Errors to resolve:          3                           │  │
 │ │ ─────────────────────────────────────────────────────── │  │
 │ │                                                          │  │
 │ │ ⚠  ROW 23: Code "1562" removed in TT99 → merge to 156  │  │
 │ │ ⚠  ROW 78: Parent "611" not found in TT99              │  │
 │ │ ⚠  ROW 91: Code "A001" invalid (non-numeric)           │  │
 │ │                                                          │  │
 │ │ [Download Error Report] [Fix and Re-upload]             │  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ ┌─ Preview (first 20 rows) ─────────────────────────────┐  │
 │ │ #  Action  Code    Name              Type    Parent   │  │
 │ │ ────────────────────────────────────────────────────── │  │
 │ │ 1  NEW     215     Tai san sinh hoc  ASSET   2        │  │
 │ │ 2  NEW     332     Phai tra co tuc   LIAB    3        │  │
 │ │ 3  UPDATE  1111    Tien mat VND      ASSET   111      │  │
 │ │ ...                                                     │  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ [← Back]                            [Confirm Import ►]     │
 └─────────────────────────────────────────────────────────────┘
```

---

## 5. COA Version Comparison

```
 ┌─────────────────────────────────────────────────────────────┐
 │ GoTax  │  COA  │  Version History                           │
 ├─────────────────────────────────────────────────────────────┤
 │ COA Version Comparison                                      │
 │                                                              │
 │ Compare: [v1.0 (2026-01-01) ▼] vs [v2.1 (2026-07-27) ▼]   │
 │                                                              │
 │ ┌─ Changes Summary ──────────────────────────────────────┐  │
 │ │ Added:       6 accounts                                 │  │
 │ │ Removed:     2 accounts                                 │  │
 │ │ Modified:    12 accounts                                │  │
 │ │ ─────────────────────────────────────────────────────── │  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ ┌─ Added Accounts ───────────────────────────────────────┐  │
 │ │ Code    Name                           Since            │  │
 │ │ ─────────────────────────────────────────────────────── │  │
 │ │ 🟢 215   Tai san sinh hoc              v2.0 15/05/2026 │  │
 │ │ 🟢 332   Phai tra co tuc loi nhuan     v2.0 15/05/2026 │  │
 │ │ 🟢 158   Hang hoa kho bao thue         v2.0 15/05/2026 │  │
 │ │ 🟢 171   Giao dich mua lai             v2.0 15/05/2026 │  │
 │ │ 🟢 82112 CP thue TNDN toi cau          v2.0 15/05/2026 │  │
 │ │ 🟢 1362  Phai thu noi bo ty gia        v3.0 27/07/2026 │  │
 │ │ 🟢 3386  Bao hiem that nghiep          v3.0 27/07/2026 │  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ ┌─ Removed Accounts ───────────────────────────────────┐  │
 │ │ Code    Name                           Until           │  │
 │ │ ───────────────────────────────────────────────────── │  │
 │ │ 🔴 611   Mua hang (TT200)              v2.0 15/05/2026│  │
 │ │ 🔴 1562  Hang hoa ban ra (TT200)       v2.0 15/05/2026│  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ ┌─ Modified Accounts ─────────────────────────────────┐  │
 │ │ Code    Attribute     Old                New          │  │
 │ │ ──────────────────────────────────────────────────── │  │
 │ │ 🟡 112   Name          "Tien gui NH"     "TG khong KH"│  │
 │ │ 🟡 242   Name          "CP tra truoc"    "CP cho PB"  │  │
 │ │ 🟡 419   Name          "Co phieu quy"    "CP mua lai" │  │
 │ │ ...                                                    │  │
 │ └────────────────────────────────────────────────────────┘  │
 │                                                              │
 │ [📥 Export Diff as PDF]  [📥 Export Diff as CSV]           │
 └─────────────────────────────────────────────────────────────┘
```

---

## 6. Approval Request Queue

```
 ┌─────────────────────────────────────────────────────────────┐
 │ GoTax  │  COA  │  Pending Approvals                         │
 ├─────────────────────────────────────────────────────────────┤
 │ Pending Approval Requests                    [3 pending]    │
 │                                                              │
 │ ┌──────────────────────────────────────────────────────┐    │
 │ │ #  Type     Account  Requestor  Date        Status   │    │
 │ │ ──────────────────────────────────────────────────── │    │
 │ │ 1  CREATE   82112    Nguyen VA  27/07 09:30  PENDING  │    │  ▶
 │ │ 2  MODIFY   1121     Tran TH    27/07 10:15  PENDING  │    │  ▶
 │ │ 3  DEACT    611      Le QH      27/07 11:00  PENDING  │    │  ▶
 │ └──────────────────────────────────────────────────────┘    │
 │                                                              │
 │ ── Review: #1 CREATE 82112 ────────────────────────────     │
 │                                                              │
 │ Account: 82112 - CP thue TNDN toi thieu toan cau            │
 │ Type: EXPENSE | Parent: 821 - CP thue TNDN                 │
 │ Requested by: Nguyen Van A (Accountant)                     │
 │ Reason: "TT99 requires new account for global minimum tax"   │
 │                                                              │
 │ ╔══════════════════════════════════════════════════════════╗ │
 │ ║  [✓ Approve]    [✗ Reject]    [View Details]            ║ │
 │ ╚══════════════════════════════════════════════════════════╝ │
 │                                                              │
 │ Rejection reason (if rejecting): [                      ]   │
 └─────────────────────────────────────────────────────────────┘
```

---

## 7. IGAP Document Generator

```
 ┌─────────────────────────────────────────────────────────────┐
 │ GoTax  │  COA  │  Generate IGAP                             │
 ├─────────────────────────────────────────────────────────────┤
 │ Generate Internal Accounting Policy (IGAP)                  │
 │                                                              │
 │ Company Information                                         │
 │ ───────────────────────────────────────────────────────────  │
 │ Company Name *  [CONG TY TNHH GOTAX                    ]    │
 │ Tax Code        [0123456789                              ]   │
 │ Address         [Hanoi, Vietnam                          ]   │
 │ Chief Acct      [Nguyen Van A                            ]   │
 │                                                              │
 │ Sections to Include                                         │
 │ ┌──────────────────────────────────────────────────────────┐ │
 │ │ [☑] Chart of Accounts (full list)                       │ │
 │ │ [☑] Account Descriptions                                │ │
 │ │ [☑] Accounting Methods (per account type)               │ │
 │ │ [☑] Voucher Cycle Descriptions                          │ │
 │ │ [☑] Internal Control Procedures                         │ │
 │ │ [☐] Financial Statement Templates                       │ │
 │ └──────────────────────────────────────────────────────────┘ │
 │                                                              │
 │ Version: [2.1]  |  Effective Date: [01/08/2026]             │
 │                                                              │
 │ [ Preview PDF ]  [ Generate ]                               │
 └─────────────────────────────────────────────────────────────┘
```

---

## 8. Mobile Responsive (Account Tree - Mobile)

```
 ┌────────────────┐
 │ GoTax │ ≡      │
 ├────────────────┤
 │ 🔍 Search COA  │
 ├────────────────┤
 │ Filters ▼      │
 ├────────────────┤
 │ ► TAI SAN      │
 │  ▼ 111 TM      │
 │    1111 VND    │
 │    1112 USD    │
 │    1113 EUR 🔒  │
 │  ▼ 112 TGNH    │
 │    1121 VND    │
 │    1122 USD    │
 ├────────────────┤
 │ ► NO PHAI TRA  │
 │ ► VON CSH      │
 │ ► DOANH THU    │
 │ ► CHI PHI      │
 ├────────────────┤
 │ [+ Add] [📥]   │
 └────────────────┘
```

---

## 9. Color Scheme & State Indicators

| State | Color | Icon |
|---|---|---|
| ACTIVE | Green (#22c55e) | ✓ |
| FROZEN | Amber (#f59e0b) | 🔒 |
| INACTIVE | Gray (#9ca3af) | — |
| PENDING (approval) | Blue (#3b82f6) | ⏳ |
| REJECTED | Red (#ef4444) | ✗ |
| Has balance | — | ! |
| Is foreign currency | — | $ |
| Is parent (group) | — | ► |
| Modified in current version | Yellow dot | 🟡 |
| Added in current version | Green dot | 🟢 |
