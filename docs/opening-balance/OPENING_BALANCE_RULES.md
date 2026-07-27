# Opening Balance Module — Business Rules

**Version:** 1.0

---

## OB-Rules: General

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-R001 | Balance equation | BLOCKER | Total opening debit = Total opening credit for each period |
| OB-R002 | One direction | BLOCKER | Each OB line: debit XOR credit > 0, never both |
| OB-R003 | Leaf account only | BLOCKER | Parent accounts (is_parent=true) cannot have opening balance |
| OB-R004 | Active account | BLOCKER | Account must be ACTIVE status to receive opening balance |
| OB-R005 | Non-frozen account | BLOCKER | FROZEN accounts cannot receive opening balance |
| OB-R006 | Open period | BLOCKER | Period must be in OPEN status (or FUTURE for initial setup) |
| OB-R007 | Unique per period | BLOCKER | One APPROVED OB per account per period per currency |
| OB-R008 | Min amount | BLOCKER | Amount must be > 0 (no zero-balance entries) |
| OB-R009 | Account exists | BLOCKER | Account must exist in company's COA |
| OB-R010 | Currency valid | BLOCKER | Currency code must be ISO 4217 and configured for company |
| OB-R011 | Status transition | BLOCKER | Only valid state transitions allowed (see state machine) |

## OB-Rules: Detail Breakdown

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-D001 | Detail sum = total | BLOCKER | Sum of detail lines must exactly equal OB total |
| OB-D002 | Entity master data | BLOCKER | Entity must exist in company's master data |
| OB-D003 | Detail only for DRAFT/PENDING | BLOCKER | Cannot modify details after approval |
| OB-D004 | Same currency | BLOCKER | Detail line currency must match parent OB currency |
| OB-D005 | No negative detail | BLOCKER | Detail amounts must be >= 0 |
| OB-D006 | Inventory quantity | WARNING | Quantity + unit price should approximate amount |

## OB-Rules: Approval

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-A001 | 4-eyes principle | BLOCKER | Approver must differ from creator |
| OB-A002 | Chief authority | BLOCKER | Only Chief Accountant or Admin can approve |
| OB-A003 | No skip | BLOCKER | DRAFT cannot become APPROVED without going through PENDING |
| OB-A004 | Rejection reason | BLOCKER | Rejection requires reason text |
| OB-A005 | Approval threshold | CONFIG | Balances > threshold (configurable, default 1B VND) require explicit approval |
| OB-A006 | Batch atomic | BLOCKER | Batch approval: all succeed or all rollback |
| OB-A007 | Auto-approve setup | CONFIG | First-time initial setup can auto-approve with chief accountant flag |

## OB-Rules: Carry-Forward

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-C001 | One-time | BLOCKER | Carry-forward can execute once per fiscal year pair |
| OB-C002 | Prior-year closed | BLOCKER | All periods in source fiscal year must be CLOSED |
| OB-C003 | Target period ready | BLOCKER | Target period must exist (FUTURE or OPEN status) |
| OB-C004 | All accounts | BLOCKER | Every balance sheet account with non-zero closing must carry forward |
| OB-C005 | Revenue/expense zero | BLOCKER | All rev/exp accounts zeroed to Account 421 |
| OB-C006 | Off-balance sheet | INFO | Account group 00x carried forward separately |
| OB-C007 | Rollback guard | BLOCKER | Rollback only if target period has ZERO posted entries |
| OB-C008 | Audit mandatory | BLOCKER | Every carried-forward balance logged in audit trail |
| OB-C009 | Exchange rate carry | BLOCKER | Foreign currency balances carry forward with exchange rate at CF date |

## OB-Rules: Correction

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-R001 | No overwrite | BLOCKER | Correction creates new OB record; original marked CORRECTED |
| OB-R002 | Reason mandatory | BLOCKER | Correction reason required |
| OB-R003 | Link preserved | BLOCKER | New OB must reference original via correction_of field |
| OB-R004 | Audit preserved | BLOCKER | Original CORRECTED balance remains visible for audit |
| OB-R005 | Cascade approval | BLOCKER | Correction requires same approval flow as new OB |

## OB-Rules: Import

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-I001 | File size | BLOCKER | Max file size: 10MB |
| OB-I002 | Row limit | BLOCKER | Max rows: 10,000 per import |
| OB-I003 | Transactional | BLOCKER | Import is all-or-nothing per batch |
| OB-I004 | Encoding auto-detect | RECOMMENDED | Support UTF-8, VNI, TCVN3 |
| OB-I005 | Template validation | BLOCKER | Required columns: account code, debit OR credit |
| OB-I006 | Duplicate in file | WARNING | Duplicate account code in file: last row wins |
| OB-I007 | Preview prior | BLOCKER | Import must be previewed before execution |

## OB-Rules: Circular 99 Migration

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-M001 | Full mapping | BLOCKER | All accounts must be mapped before migration |
| OB-M002 | Year-end boundary | WARNING | Migration should occur at fiscal year boundary (31-Dec) |
| OB-M003 | Journal entries | BLOCKER | Migration creates auditable journal entries for each transfer |
| OB-M004 | No partial | BLOCKER | Migration is all-or-nothing |
| OB-M005 | Abolished accounts | BLOCKER | Accounts removed in Circular 99 must be zeroed |
| OB-M006 | New accounts | WARNING | New accounts (215, 332, etc.) must be created if non-zero |
| OB-M007 | Split ratio | BLOCKER | Split mapping requires ratio summing to 1.0 (100%) |

## OB-Rules: Reporting

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-R001 | Approved only | BLOCKER | Only APPROVED OBs appear in financial reports |
| OB-R002 | Period filter | BLOCKER | Reports default to current open period |
| OB-R003 | Balance validation | BLOCKER | Report shows (un)balanced status prominently |
| OB-R004 | Comparative | BLOCKER | Balance Sheet shows "Đầu kỳ" and "Cuối kỳ" columns |
| OB-R005 | Audit export | RECOMMENDED | Report exportable to PDF/Excel for audit sign-off |

## OB-Rules: Security

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-S001 | Read access | BLOCKER | All authenticated users with company access can read OBs |
| OB-S002 | Write access | BLOCKER | Only Accountant role and above can create/edit OBs |
| OB-S003 | Approve access | BLOCKER | Only Chief Accountant and Admin can approve |
| OB-S004 | Delete access | BLOCKER | Only Admin can delete OBs (DRAFT status only) |
| OB-S005 | ABAC sensitive | CONFIG | Accounts flagged as "sensitive" require 2FA for modification |

## OB-Rules: Fiscal Year and Period

| ID | Rule | Severity | Description |
|----|------|----------|-------------|
| OB-F001 | Default OB date | INFO | Opening balance date = accounting start date (company setting) |
| OB-F002 | Short fiscal year | WARNING | Short fiscal year (< 12 months) supported with proportional CF |
| OB-F003 | Mid-year start | INFO | Company starting mid-year: all OBs entered manually per Account |
| OB-F004 | Period transfer | BLOCKER | Balances transfer only between consecutive fiscal years |
| OB-F005 | Multi-FY | WARNING | Different companies in same tenant can have different FY calendars |
