# Payroll Module — Process Flows & User Journeys

**Document ID:** PROC-PAYROLL-001
**Version:** 1.0
**Date:** 2026-08-04

---

## 1. End-to-End Payroll Process

### 1.1 Monthly Payroll Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        PAYROLL LIFECYCLE                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐              │
│  │  DRAFT  │───▶│PROCESSING│───▶│ APPROVED│───▶│  PAID   │              │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘              │
│       │              │              │              │                      │
│       ▼              ▼              ▼              ▼                      │
│  Import data    Calculate      GL posting     Bank file                 │
│  Update info    Review         Payslips       Confirm                   │
│  Validate       Adjust         Lock period    Archive                   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Status Transitions

| From | To | Trigger | Action |
|------|----|---------|--------|
| DRAFT | PROCESSING | Clerk clicks "Submit for Review" | Calculation locked |
| PROCESSING | APPROVED | Approver clicks "Approve" | GL entries posted, payslips generated |
| PROCESSING | DRAFT | Approver clicks "Reject" | Returned to clerk for revision |
| APPROVED | PAID | Clerk confirms payment | Payment recorded |
| PAID | CLOSED | System auto-close (month-end) | Period archived |

---

## 2. Detailed Process Maps

### 2.1 Employee Onboarding for Payroll

```
┌─────────────────────────────────────────────────────────────────┐
│              EMPLOYEE PAYROLL ONBOARDING                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  New Employee Hired                                               │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ HR Creates Employee Record          │                         │
│  │ (Company Module)                    │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Payroll Clerk Creates               │                         │
│  │ EmployeePayrollInfo                 │                         │
│  │ - Contract type                     │                         │
│  │ - Salary type & amount              │                         │
│  │ - Allowances                        │                         │
│  │ - Insurance base                    │                         │
│  │ - Region                            │                         │
│  │ - Bank account                      │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Register for Social Insurance       │                         │
│  │ - Generate TK1-TS (if new SI code)  │                         │
│  │ - Add to D02-TS list                │                         │
│  │ - Submit to BHXH within 30 days     │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Set Up Leave Balance                │                         │
│  │ - Annual leave: 12 days             │                         │
│  │ - Sick leave: per policy            │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Employee Ready for Payroll          │                         │
│  └─────────────────────────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Timekeeping Import Process

```
┌─────────────────────────────────────────────────────────────────┐
│              TIMEKEEPING IMPORT PROCESS                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  External Timekeeping System                                      │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ Export CSV/Excel                    │                         │
│  │ Columns:                            │                         │
│  │ - Employee Code                     │                         │
│  │ - Date                              │                         │
│  │ - Clock In / Clock Out              │                         │
│  │ - OT Hours                          │                         │
│  │ - Leave Type (if applicable)        │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ GoTax: Payroll → Import Timekeeping │                         │
│  │ - Upload file                       │                         │
│  │ - System validates format           │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│       ┌─────────┴─────────┐                                       │
│       ▼                   ▼                                       │
│  ┌─────────┐        ┌─────────┐                                  │
│  │ SUCCESS │        │ ERROR   │                                  │
│  └────┬────┘        └────┬────┘                                  │
│       │                   │                                       │
│       ▼                   ▼                                       │
│  ┌─────────────┐   ┌─────────────┐                               │
│  │ Records     │   │ Error Log   │                               │
│  │ Created     │   │ Shown       │                               │
│  │ Summary     │   │ Clerk Fixes │                               │
│  │ Displayed   │   │ & Re-import │                               │
│  └─────────────┘   └─────────────┘                               │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.3 Payroll Calculation Process

```
┌─────────────────────────────────────────────────────────────────┐
│              PAYROLL CALCULATION PROCESS                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Clerk Clicks "Calculate Payroll"                                 │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ Step 1: Validate Inputs             │                         │
│  │ - All employees have salary info?   │                         │
│  │ - Timekeeping data complete?        │                         │
│  │ - Insurance base set?               │                         │
│  │ - Region assigned?                  │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│       ┌─────────┴─────────┐                                       │
│       ▼                   ▼                                       │
│  ┌─────────┐        ┌─────────┐                                  │
│  │ PASS    │        │ FAIL    │                                  │
│  └────┬────┘        └────┬────┘                                  │
│       │                   │                                       │
│       ▼                   ▼                                       │
│  ┌─────────────┐   ┌─────────────┐                               │
│  │ Step 2:     │   │ Show errors │                               │
│  │ Calculate   │   │ Clerk fixes │                               │
│  │ for each    │   └─────────────┘                               │
│  │ employee    │                                                  │
│  └──────┬──────┘                                                  │
│         │                                                         │
│         ▼                                                         │
│  ┌─────────────────────────────────────┐                         │
│  │ For Each Employee:                  │                         │
│  │                                     │                         │
│  │ 2a. Calculate Income                │                         │
│  │     Base + OT + Night + Leave +     │                         │
│  │     Holiday + Allowances + Bonus    │                         │
│  │                                     │                         │
│  │ 2b. Calculate Insurance Base        │                         │
│  │     MIN(salary, 20 × ref_level)     │                         │
│  │                                     │                         │
│  │ 2c. Calculate Employee Deductions   │                         │
│  │     SI(8%) + HI(1.5%) + UI(1%)     │                         │
│  │                                     │                         │
│  │ 2d. Calculate PIT                   │                         │
│  │     gross - insurance - deductions  │                         │
│  │     Apply 5-bracket schedule        │                         │
│  │                                     │                         │
│  │ 2e. Calculate Net Pay               │                         │
│  │     gross - all deductions          │                         │
│  │                                     │                         │
│  │ 2f. Calculate Employer Costs        │                         │
│  │     SI(17.5%) + HI(3%) + UI(1%)    │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Step 3: Generate Summary            │                         │
│  │ - Total gross salary                │                         │
│  │ - Total deductions                  │                         │
│  │ - Total net pay                     │                         │
│  │ - Total employer cost               │                         │
│  │ - Per-department breakdown          │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Step 4: Flag Warnings               │                         │
│  │ - Salary below minimum wage         │                         │
│  │ - OT exceeds 40h/month              │                         │
│  │ - Missing insurance info            │                         │
│  │ - Employee missing SI number        │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Payroll Ready for Review            │                         │
│  └─────────────────────────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.4 Approval & GL Posting Process

```
┌─────────────────────────────────────────────────────────────────┐
│              APPROVAL & GL POSTING PROCESS                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Clerk Submits for Review                                        │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ Chief Accountant Reviews            │                         │
│  │ - Check payroll summary             │                         │
│  │ - Compare with budget               │                         │
│  │ - Review flagged items              │                         │
│  │ - Spot-check individual payslips    │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│       ┌─────────┴─────────┐                                       │
│       ▼                   ▼                                       │
│  ┌─────────┐        ┌─────────┐                                  │
│  │ APPROVE │        │ REJECT  │                                  │
│  └────┬────┘        └────┬────┘                                  │
│       │                   │                                       │
│       ▼                   ▼                                       │
│  ┌─────────────┐   ┌─────────────┐                               │
│  │ System:     │   │ Return to   │                               │
│  │ - Lock      │   │ Clerk with  │                               │
│  │   period    │   │ comments    │                               │
│  │ - Post GL   │   └─────────────┘                               │
│  │   entries   │                                                  │
│  │ - Generate  │                                                  │
│  │   payslips  │                                                  │
│  │ - Update    │                                                  │
│  │   leave     │                                                  │
│  │   balances  │                                                  │
│  └──────┬──────┘                                                  │
│         │                                                         │
│         ▼                                                         │
│  ┌─────────────────────────────────────┐                         │
│  │ Journal Entries Posted to GL        │                         │
│  │                                     │                         │
│  │ Dr. 6421 Salary Expense     XXX    │                         │
│  │ Dr. 6422 OT Expense         XXX    │                         │
│  │ Dr. 6424 Insurance Expense  XXX    │                         │
│  │     Cr. 3331 Salary Payable    XXX │                         │
│  │     Cr. 3332 SI Payable (ER)  XXX │                         │
│  │     Cr. 3333 HI Payable (ER)  XXX │                         │
│  │     Cr. 3334 UI Payable (ER)  XXX │                         │
│  │     Cr. 3335 SI Payable (EE)  XXX │                         │
│  │     Cr. 3336 HI Payable (EE)  XXX │                         │
│  │     Cr. 3337 UI Payable (EE)  XXX │                         │
│  │     Cr. 3338 PIT Payable      XXX │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Payslips Generated (PDF)            │                         │
│  │ - Per employee                      │                         │
│  │ - Email distribution                │                         │
│  │ - Self-service portal updated       │                         │
│  └─────────────────────────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.5 Payment Process

```
┌─────────────────────────────────────────────────────────────────┐
│              SALARY PAYMENT PROCESS                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Payroll Approved                                                │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ Generate Bank Salary File           │                         │
│  │ - Format: CSV/XML (bank-specific)   │                         │
│  │ - Fields:                           │                         │
│  │   - Employee bank account           │                         │
│  │   - Amount (net pay)                │                         │
│  │   - Reference (payslip number)      │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Upload to Bank Portal               │                         │
│  │ - Login to bank                     │                         │
│  │ - Upload salary file                │                         │
│  │ - Confirm batch                     │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Bank Processes Payments             │                         │
│  │ - Transfers to employee accounts    │                         │
│  │ - Returns confirmation file         │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ GoTax: Record Payment               │                         │
│  │ - Update payroll status → PAID      │                         │
│  │ - Post payment journal entry:       │                         │
│  │   Dr. 3331 Salary Payable   XXX    │                         │
│  │   Cr. 1111 Bank Account     XXX    │                         │
│  │ - Send payment notification         │                         │
│  └─────────────────────────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. User Journeys

### 3.1 Payroll Clerk — First-Time Setup

```
┌─────────────────────────────────────────────────────────────────┐
│     PAYROLL CLERK — FIRST-TIME SETUP JOURNEY                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Step 1: Configure Payroll Parameters                            │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Navigate to Payroll → Configuration                 │        │
│  │ ☐ Set base salary (VND 2,530,000 from Jul 2026)      │        │
│  │ ☐ Set regional minimum wages                          │        │
│  │ ☐ Set insurance rates (SI 8%/17.5%, HI 1.5%/3%,     │        │
│  │   UI 1%/1%)                                          │        │
│  │ ☐ Set SI/HI cap (VND 50,600,000)                     │        │
│  │ ☐ Set UI cap per region                               │        │
│  │ ☐ Set PIT deductions (VND 15.5M personal,            │        │
│  │   VND 6.2M per dependant)                             │        │
│  │ ☐ Set trade union rate (1%, cap VND 253,000)          │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 2: Create Salary Components                                │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Base Salary (INCOME, FIXED)                         │        │
│  │ ☐ Transport Allowance (INCOME, FIXED)                 │        │
│  │ ☐ Meal Allowance (INCOME, FIXED)                      │        │
│  │ ☐ Phone Allowance (INCOME, FIXED)                     │        │
│  │ ☐ Responsibility Allowance (INCOME, FIXED)            │        │
│  │ ☐ Seniority Allowance (INCOME, FIXED)                 │        │
│  │ ☐ Overtime Pay (INCOME, FORMULA)                      │        │
│  │ ☐ Night Shift Pay (INCOME, FORMULA)                   │        │
│  │ ☐ Performance Bonus (INCOME, FIXED)                   │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 3: Create Salary Templates                                 │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ "Office Staff" template                             │        │
│  │   - Base Salary + Transport + Meal + Phone            │        │
│  │ ☐ "Factory Worker" template                           │        │
│  │   - Base Salary + Transport + Meal + Responsibility   │        │
│  │ ☐ "Manager" template                                  │        │
│  │   - Base Salary + Transport + Meal + Phone +          │        │
│  │     Responsibility + Seniority                        │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 4: Update Employee Payroll Info                            │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ For each employee:                                  │        │
│  │   - Set contract type                                 │        │
│  │   - Set salary type and amount                        │        │
│  │   - Apply salary template                             │        │
│  │   - Set insurance base salary                         │        │
│  │   - Set region                                        │        │
│  │   - Add dependants                                    │        │
│  │   - Set bank account                                  │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 5: Import Historical Data                                  │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Import current leave balances                       │        │
│  │ ☐ Import last 3 months payslips (optional)            │        │
│  │ ☐ Verify insurance registration data                  │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 6: Test Run                                               │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Create test period (current month)                  │        │
│  │ ☐ Import test timekeeping data                        │        │
│  │ ☐ Run calculation                                     │        │
│  │ ☐ Verify results against manual calculation           │        │
│  │ ☐ Check payslip output                                │        │
│  │ ☐ Verify GL journal entries                           │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 Payroll Clerk — Monthly Routine

```
┌─────────────────────────────────────────────────────────────────┐
│     PAYROLL CLERK — MONTHLY ROUTINE JOURNEY                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Day 1: Start of Month                                          │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Login to GoTax                                     │        │
│  │ ☐ Navigate to Payroll → Periods                       │        │
│  │ ☐ Create new period (e.g., August 2026)               │        │
│  │ ☐ Period status: DRAFT                                │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Day 2-3: Data Collection                                        │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Export timekeeping from attendance system            │        │
│  │ ☐ Import to GoTax (CSV upload)                         │        │
│  │ ☐ Review import summary                                │        │
│  │ ☐ Fix any errors                                       │        │
│  │ ☐ Update employee changes:                             │        │
│  │   - New hires (create payroll info)                    │        │
│  │   - Terminations (final pay calculation)               │        │
│  │   - Salary changes (update payroll info)               │        │
│  │   - Transfers (update department)                      │        │
│  │ ☐ Collect pending leave requests                       │        │
│  │ ☐ Approve/reject leave requests                        │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Day 4-5: Calculation                                            │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Click "Calculate Payroll"                           │        │
│  │ ☐ Wait for calculation (should be < 30 seconds)       │        │
│  │ ☐ Review summary:                                      │        │
│  │   - Total gross: VND XXX                              │        │
│  │   - Total deductions: VND XXX                         │        │
│  │   - Total net pay: VND XXX                            │        │
│  │   - Total employer cost: VND XXX                      │        │
│  │ ☐ Review warnings:                                     │        │
│  │   - Salary below minimum wage?                         │        │
│  │   - OT exceeds limit?                                  │        │
│  │   - Missing insurance info?                            │        │
│  │ ☐ Click individual employees to review detail          │        │
│  │ ☐ Make adjustments if needed (with reason)             │        │
│  │ ☐ Click "Submit for Review"                            │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Day 6-7: Approval                                               │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Notify Chief Accountant for review                  │        │
│  │ ☐ Wait for approval                                   │        │
│  │ ☐ If rejected: fix issues and resubmit                 │        │
│  │ ☐ If approved: system posts GL entries                 │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Day 8-10: Payment                                               │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Generate bank salary file                           │        │
│  │ ☐ Upload to bank portal                               │        │
│  │ ☐ Confirm payment completion                          │        │
│  │ ☐ Record payment in GoTax                             │        │
│  │ ☐ System posts payment journal entry                  │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Day 11-15: Declarations & Closing                               │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ ☐ Generate D02-TS (if new employees this month)       │        │
│  │ ☐ Generate 05/KK-TNCN (quarterly)                     │        │
│  │ ☐ Submit declarations electronically                  │        │
│  │ ☐ Archive payslips                                    │        │
│  │ ☐ Period auto-closes at month-end                      │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 3.3 Employee — Payslip Viewing Journey

```
┌─────────────────────────────────────────────────────────────────┐
│     EMPLOYEE — PAYSLIP VIEWING JOURNEY                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Step 1: Login                                                   │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ Employee opens GoTax self-service portal               │        │
│  │ Enters credentials and 2FA code                        │        │
│  │ Dashboard loads                                        │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 2: Navigate to Payslips                                   │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ Click "My Payslips" in menu                           │        │
│  │ List of payslips by month displayed                   │        │
│  │ Current month highlighted                             │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 3: View Payslip                                           │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ Click on month (e.g., "July 2026")                    │        │
│  │ Payslip detail displays:                               │        │
│  │                                                        │        │
│  │ INCOME                                                 │        │
│  │   Base Salary:        15,000,000                      │        │
│  │   OT Pay:              1,200,000                      │        │
│  │   Night Shift:           300,000                      │        │
│  │   Transport Allow:      500,000                       │        │
│  │   Meal Allow:           500,000                       │        │
│  │   ─────────────────────────────                       │        │
│  │   GROSS SALARY:        17,500,000                     │        │
│  │                                                        │        │
│  │ DEDUCTIONS                                            │        │
│  │   Social Insurance:    (1,282,500)                    │        │
│  │   Health Insurance:      (191,250)                    │        │
│  │   Unemployment Ins:      (150,000)                    │        │
│  │   Trade Union:             (12,500)                   │        │
│  │   Personal Income Tax:      (0)                      │        │
│  │   ─────────────────────────────                       │        │
│  │   TOTAL DEDUCTIONS:   (1,636,250)                     │        │
│  │                                                        │        │
│  │ ══════════════════════════════════                     │        │
│  │ NET PAY:              15,863,750                      │        │
│  │ ══════════════════════════════════                     │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 4: Download PDF                                            │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ Click "Download PDF" button                           │        │
│  │ PDF generated and downloaded                          │        │
│  │ Employee saves for personal records                    │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
│  Step 5: Dispute (if needed)                                    │
│  ┌──────────────────────────────────────────────────────┐        │
│  │ Employee notices incorrect OT hours                    │        │
│  │ Clicks "Dispute" button                               │        │
│  │ Enters reason: "OT on July 5 should be 4h, not 2h"   │        │
│  │ Submits dispute                                       │        │
│  │ System sends notification to Payroll Clerk            │        │
│  │ Clerk investigates and resolves                       │        │
│  │ Employee receives resolution notification             │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 4. Exception Handling Flows

### 4.1 Employee Terminated Mid-Month

```
┌─────────────────────────────────────────────────────────────────┐
│     TERMINATION MID-MONTH FLOW                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Employee Terminated (e.g., July 15)                             │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ Calculate Final Pay:                │                         │
│  │                                     │                         │
│  │ 1. Prorated salary (15/26 days)     │                         │
│  │    = 15,000,000 × 15/26             │                         │
│  │    = 8,653,846                      │                         │
│  │                                     │                         │
│  │ 2. Unused annual leave payout       │                         │
│  │    = (15,000,000/26) × remaining    │                         │
│  │    days                             │                         │
│  │                                     │                         │
│  │ 3. Severance pay (if 12+ months)    │                         │
│  │    = 0.5 month × years of service   │                         │
│  │                                     │                         │
│  │ 4. Insurance for final month        │                         │
│  │    - Still required to contribute    │                         │
│  │                                     │                         │
│  │ 5. PIT on final pay                 │                         │
│  │    - Include all termination income  │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Generate D02-TS with "Termination"  │                         │
│  │ - Mark employee as reduced          │                         │
│  │ - Submit to BHXH                    │                         │
│  └─────────────────────────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Retroactive Pay Adjustment

```
┌─────────────────────────────────────────────────────────────────┐
│     RETROACTIVE PAY ADJUSTMENT FLOW                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Scenario: Employee promoted June 1, but payroll processed       │
│  in July without updated salary                                  │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ Clerk Creates Adjustment:           │                         │
│  │ - Employee: Nguyen Van A            │                         │
│  │ - Period: June 2026                 │                         │
│  │ - Type: Salary increase             │                         │
│  │ - Old salary: 15,000,000            │                         │
│  │ - New salary: 20,000,000            │                         │
│  │ - Effective: June 1, 2026           │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ System Calculates Arrears:          │                         │
│  │                                     │                         │
│  │ Gross arrears:                      │                         │
│  │   (20M - 15M) × 1 month = 5,000,000│                         │
│  │                                     │                         │
│  │ Insurance arrears:                  │                         │
│  │   SI: 5M × 8% = 400,000            │                         │
│  │   HI: 5M × 1.5% = 75,000           │                         │
│  │   UI: 5M × 1% = 50,000             │                         │
│  │                                     │                         │
│  │ PIT arrears:                        │                         │
│  │   Recalculate June PIT with new     │                         │
│  │   salary, subtract already paid     │                         │
│  │                                     │                         │
│  │ Net arrears to employee:            │                         │
│  │   5,000,000 - 525,000 - pit_diff    │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Arrears Added to Current Payroll    │                         │
│  │ - Line item: "Arrears - June 2026"  │                         │
│  │ - Included in current month gross   │                         │
│  │ - Insurance and PIT calculated      │                         │
│  └─────────────────────────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 4.3 Minimum Wage Violation

```
┌─────────────────────────────────────────────────────────────────┐
│     MINIMUM WAGE VIOLATION FLOW                                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  System Detects: Employee salary < regional minimum              │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ Warning Displayed:                  │                         │
│  │ "Employee E001 (Nguyen Van A)      │                         │
│  │  salary VND 4,500,000 is below     │                         │
│  │  Region I minimum VND 5,310,000"   │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│       ┌─────────┴─────────┐                                       │
│       ▼                   ▼                                       │
│  ┌─────────┐        ┌─────────┐                                  │
│  │ FIX     │        │ OVERRIDE│                                  │
│  └────┬────┘        └────┬────┘                                  │
│       │                   │                                       │
│       ▼                   ▼                                       │
│  ┌─────────────┐   ┌─────────────┐                               │
│  │ Clerk       │   │ Clerk       │                               │
│  │ updates     │   │ overrides   │                               │
│  │ salary to   │   │ with        │                               │
│  │ minimum     │   │ justification│                              │
│  │ wage        │   │ (legal      │                               │
│  │             │   │  exception) │                               │
│  └─────────────┘   └─────────────┘                               │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 5. Integration Flows

### 5.1 GL Integration

```
┌─────────────────────────────────────────────────────────────────┐
│     GL INTEGRATION FLOW                                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Payroll Approved                                                │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ System Generates Journal Entries:   │                         │
│  │                                     │                         │
│  │ Entry 1: Salary Accrual             │                         │
│  │   Dr. 6421 Salary Expense           │                         │
│  │   Dr. 6422 OT Expense               │                         │
│  │   Dr. 6423 Allowance Expense        │                         │
│  │       Cr. 3331 Salary Payable       │                         │
│  │                                     │                         │
│  │ Entry 2: Employer Insurance         │                         │
│  │   Dr. 6424 Insurance Expense        │                         │
│  │       Cr. 3332 SI Payable (ER)      │                         │
│  │       Cr. 3333 HI Payable (ER)      │                         │
│  │       Cr. 3334 UI Payable (ER)      │                         │
│  │                                     │                         │
│  │ Entry 3: Employee Deductions        │                         │
│  │       Cr. 3335 SI Payable (EE)      │                         │
│  │       Cr. 3336 HI Payable (EE)      │                         │
│  │       Cr. 3337 UI Payable (EE)      │                         │
│  │       Cr. 3338 PIT Payable          │                         │
│  │   (offset against Salary Payable)   │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Entries Posted to GL                │                         │
│  │ - Period: July 2026                 │                         │
│  │ - Reference: PAYROLL-2026-07        │                         │
│  │ - Status: Posted                    │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ GL Shows:                           │                         │
│  │ - Salary Expense: VND XXX          │                         │
│  │ - Insurance Expense: VND XXX       │                         │
│  │ - Salary Payable: VND XXX          │                         │
│  │ - Insurance Payable: VND XXX       │                         │
│  │ - PIT Payable: VND XXX             │                         │
│  └─────────────────────────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 Social Insurance Declaration Flow

```
┌─────────────────────────────────────────────────────────────────┐
│     SI DECLARATION FLOW                                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Monthly Trigger (or when employee changes occur)                │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────┐                         │
│  │ System Auto-Detects Changes:        │                         │
│  │                                     │                         │
│  │ New employees (need SI registration)│                         │
│  │ - Generate TK1-TS for each          │                         │
│  │ - Add to D02-TS list                │                         │
│  │                                     │                         │
│  │ Salary changes (need adjustment)    │                         │
│  │ - Add to D02-TS with "Adjustment"   │                         │
│  │                                     │                         │
│  │ Terminated employees                │                         │
│  │ - Add to D02-TS with "Reduction"    │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Clerk Reviews D02-TS:               │                         │
│  │ - Verify employee data              │                         │
│  │ - Confirm insurance base amounts    │                         │
│  │ - Approve declaration               │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ System Generates XML:               │                         │
│  │ - Format: BHXH electronic format    │                         │
│  │ - Includes: D02-TS + TK1-TS         │                         │
│  │ - Digital signature applied         │                         │
│  └──────────────┬──────────────────────┘                         │
│                 │                                                 │
│                 ▼                                                 │
│  ┌─────────────────────────────────────┐                         │
│  │ Submit via I-VAN:                   │                         │
│  │ - Upload XML to I-VAN portal        │                         │
│  │ - Submit to BHXH                    │                         │
│  │ - Receive confirmation              │                         │
│  │ - Archive receipt                   │                         │
│  └─────────────────────────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 6. Dashboard Views

### 6.1 Payroll Dashboard

```
┌─────────────────────────────────────────────────────────────────┐
│                    PAYROLL DASHBOARD                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ Current Period: August 2026                    Status: DRAFT │    │
│  ├─────────────────────────────────────────────────────────┤    │
│  │                                                         │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐│    │
│  │  │Employees │  │ Gross    │  │ Net Pay  │  │Employer  ││    │
│  │  │    85    │  │1,250M    │  │1,050M    │  │Cost 280M ││    │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘│    │
│  │                                                         │    │
│  │  ┌──────────────────────────────────────────────────┐  │    │
│  │  │ Warnings (3)                                      │  │    │
│  │  │ ⚠ E001: Salary below minimum (Region I)          │  │    │
│  │  │ ⚠ E045: OT 42h exceeds 40h/month limit          │  │    │
│  │  │ ⚠ E072: Missing insurance base salary            │  │    │
│  │  └──────────────────────────────────────────────────┘  │    │
│  │                                                         │    │
│  │  ┌──────────────────────────────────────────────────┐  │    │
│  │  │ Quick Actions                                     │  │    │
│  │  │ [Import Timekeeping] [Calculate] [Submit Review]  │  │    │
│  │  └──────────────────────────────────────────────────┘  │    │
│  │                                                         │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                   │
│  Recent Periods                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ July 2026    │ APPROVED │ 85 employees │ 1,200M gross   │    │
│  │ June 2026    │ PAID     │ 83 employees │ 1,180M gross   │    │
│  │ May 2026     │ CLOSED   │ 82 employees │ 1,150M gross   │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 7. Report Templates

### 7.1 Monthly Payroll Summary

```
┌─────────────────────────────────────────────────────────────────┐
│                 MONTHLY PAYROLL SUMMARY                          │
│                 Period: July 2026                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  SUMMARY                                                         │
│  ─────────────────────────────────────────────                   │
│  Total Employees:           85                                   │
│  Active Employees:          82                                   │
│  New Hires:                  2                                   │
│  Terminations:               1                                   │
│                                                                   │
│  INCOME                                                          │
│  ─────────────────────────────────────────────                   │
│  Total Base Salary:     1,275,000,000                           │
│  Total OT Pay:            125,000,000                           │
│  Total Night Shift:        15,000,000                           │
│  Total Allowances:         85,000,000                           │
│  Total Bonuses:            50,000,000                           │
│  ─────────────────────────────────────────────                   │
│  TOTAL GROSS:           1,550,000,000                           │
│                                                                   │
│  DEDUCTIONS                                                       │
│  ─────────────────────────────────────────────                   │
│  Employee SI:             124,000,000                           │
│  Employee HI:              23,250,000                           │
│  Employee UI:              15,500,000                           │
│  Trade Union:               1,550,000                           │
│  PIT:                      65,000,000                           │
│  ─────────────────────────────────────────────                   │
│  TOTAL DEDUCTIONS:       229,300,000                           │
│                                                                   │
│  NET PAY:            1,320,700,000                               │
│                                                                   │
│  EMPLOYER COSTS                                                   │
│  ─────────────────────────────────────────────                   │
│  Employer SI:            271,250,000                           │
│  Employer HI:             46,500,000                           │
│  Employer UI:             15,500,000                           │
│  Employer Trade Union:     1,550,000                           │
│  ─────────────────────────────────────────────                   │
│  TOTAL EMPLOYER COST:   1,821,500,000                           │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```
