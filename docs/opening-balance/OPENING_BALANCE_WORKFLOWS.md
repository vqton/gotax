# Opening Balance Module — Workflows & Data Flows

**Version:** 1.0

---

## 1. Opening Balance Lifecycle State Machine

```
                    ┌──────────┐
                    │  DRAFT   │
                    └────┬─────┘
                         │ submit
                         ▼
                    ┌──────────┐
         ┌──────────│ PENDING  │──────────┐
         │          └────┬─────┘          │
         │ approve       │       reject   │
         ▼               ▼                ▼
    ┌──────────┐   ┌──────────┐    ┌──────────┐
    │ APPROVED │   │ (active) │    │ REJECTED │
    └────┬─────┘   └──────────┘    └────┬─────┘
         │                              │ resubmit
         │ correct                       ▼
         ▼                          ┌──────────┐
    ┌──────────┐                    │  DRAFT   │
    │CORRECTED │                    └──────────┘
    └──────────┘
```

**Transitions:**
- DRAFT → PENDING: accountant submits for approval
- PENDING → APPROVED: chief accountant approves
- PENDING → REJECTED: chief accountant rejects (with reason)
- REJECTED → DRAFT: accountant edits and resubmits
- APPROVED → CORRECTED: correction request approved
- CORRECTED → (new) PENDING: new version awaiting approval

---

## 2. Initial Setup Workflow (New Company)

```
┌─────────┐    ┌────────────┐    ┌──────────────┐    ┌──────────┐
│Admin    │    │System      │    │Chief Acc      │    │System    │
│Creates  │───►│Generates   │───►│Reviews &      │───►│OB Ready  │
│Company  │    │COA + FY    │    │Enters OBs     │    │          │
└─────────┘    └────────────┘    └──────┬───────┘    └──────────┘
                                        │
                                        ├── Manual entry (single)
                                        ├── Import from Excel
                                        └── Import from old system
                                             │
                                             ▼
                                        ┌────────────┐
                                        │Validate    │
                                        │OB balanced?│
                                        └──────┬─────┘
                                         Yes/  │  No
                                              │
                                              ▼
                                        ┌────────────┐
                                        │ Auto-balance│
                                        │ with 999?   │
                                        └────────────┘
```

---

## 3. Opening Balance Entry Data Flow

```
 UI / API Request                     Service                         Repository                    Database
 ──────────────────            ──────────────────────           ────────────────────           ──────────────
         │                             │                              │                          │
         │ POST opening-balances       │                              │                          │
         ├────────────────────────────►│                              │                          │
         │                             │ Validate account             ├── GetByCode()             │
         │                             ├─────────────────────────────►│──── accounts              │
         │                             │◄─────────────────────────────┤                          │
         │                             │                              │                          │
         │                             │ Validate period              ├── GetPeriodByID()          │
         │                             ├─────────────────────────────►│──── period_v2             │
         │                             │◄─────────────────────────────┤                          │
         │                             │                              │                          │
         │                             │ Check duplicate             ├── GetByAccount()           │
         │                             ├─────────────────────────────►│──── opening_balances      │
         │                             │◄─────────────────────────────┤                          │
         │                             │                              │                          │
         │                             │ Generate ID                  │                          │
         │                             │ Set status = DRAFT           │                          │
         │                             │ Set created_by               │                          │
         │                             │                              │                          │
         │                             ├── Create()                  ├── INSERT                  │
         │                             ├─────────────────────────────►│──── opening_balances      │
         │                             │                              │                          │
         │                             │ If detail lines:            ├── BulkCreateDetails()      │
         │                             ├─────────────────────────────►│──── ob_details            │
         │                             │                              │                          │
         │                             ├── Create()                  ├── INSERT                  │
         │                             ├─────────────────────────────►│──── audit_entries         │
         │                             │                              │                          │
         │  201 Created                │                              │                          │
         │◄────────────────────────────┤                              │                          │
```

---

## 4. Fiscal Year Carry-Forward Data Flow

```
 Chief Accountant                   Service                          Repository                    Database
 ────────────────            ──────────────────────           ────────────────────           ──────────────
         │                             │                              │                          │
         │ POST /carry-forward          │                              │                          │
         ├────────────────────────────►│                              │                          │
         │                             │                              │                          │
         │                             │ Validate:                    │                          │
         │                             │ - All periods closed in FY  │                          │
         │                             │ - No CF already exists      │                          │
         │                             │ - Target period exists      │                          │
         │                             │                              │                          │
         │                             ├── Get trial balance         │                          │
         │                             │   for closing period        │                          │
         │                             ├─────────────────────────────►│                          │
         │                             │                              │                          │
         │   Preview                   │  Compute:                   │                          │
         │◄────────────────────────────┤  - Revenue/expense to close │                          │
         │                             │  - B/S to carry forward    │                          │
         │                             │  - Total balancing check    │                          │
         │                             │                              │                          │
         │  Confirm                    │                              │                          │
         ├────────────────────────────►│                              │                          │
         │                             │                              │                          │
         │                             │ BEGIN TRANSACTION            │                          │
         │                             │                              │                          │
         │                             │ 1. Create closing JEs        ├── Create()               │
         │                             │    (rev/exp → 421)          ├────────────────────────►│── journal_entries
         │                             │                              │                          │
         │                             │ 2. Post closing entries      ├── PostEntry()            │
         │                             ├─────────────────────────────►│                          │
         │                             │                              │                          │
         │                             │ 3. For each B/S account:    ├── Create()               │
         │                             │    Create OB for new period ├────────────────────────►│── opening_balances
         │                             │    (source=CARRY_FORWARD)   │                          │
         │                             │    status=APPROVED          │                          │
         │                             │                              │                          │
         │                             │ 4. Create CarryForwardLog   ├── CreateCarryForward()    │
         │                             ├─────────────────────────────►│── carry_forward_logs    │
         │                             │                              │                          │
         │                             │ COMMIT                       │                          │
         │                             │                              │                          │
         │                             │ SET previous FY → LOCKED    │                          │
         │                             │                              │                          │
         │  Success summary            │                              │                          │
         │◄────────────────────────────┤                              │                          │
```

---

## 5. Circular 99 Migration Data Flow

```
 Admin/Chief Acc                    Service                           Repository                   Database
 ──────────────────            ───────────────────────          ────────────────────          ──────────────
         │                             │                              │                          │
         │ POST /migrations/execute     │                              │                          │
         ├────────────────────────────►│                              │                          │
         │                             │                              │                          │
         │                             │ 1. Load Circular 99 mapping ├── GetAll()                │
         │                             ├─────────────────────────────►│── circular99_mappings   │
         │                             │                              │                          │
         │                             │ 2. Get current OBs          ├── GetByPeriod()          │
         │                             ├─────────────────────────────►│── opening_balances       │
         │                             │                              │                          │
         │                             │ 3. For each OB:             │                          │
         │                             │    - If DIRECT: new code    │                          │
         │                             │    - If REMOVE: zero out    │                          │
         │                             │    - If MERGE: combine      │                          │
         │                             │    - If SPLIT: allocate%    │                          │
         │                             │                              │                          │
         │                             │ 4. Create transfer JEs      ├── Create()               │
         │                             ├─────────────────────────────►│── journal_entries       │
         │                             │                              │                          │
         │                             │ 5. Create new OBs           ├── BulkCreate()            │
         │                             │    under TT99 COA           ├────────────────────────►│── opening_balances
         │                             │                              │                          │
         │                             │ 6. Create migration record  ├── CreateMigration()       │
         │                             ├─────────────────────────────►│── balance_migrations    │
         │                             │                              │                          │
         │  Migration result           │                              │                          │
         │◄────────────────────────────┤                              │                          │
```

---

## 6. Excel Import Data Flow

```
 Accountant                     Service                           Repository                   Database
 ────────────             ───────────────────────          ────────────────────          ──────────────
     │                            │                              │                          │
     │ POST /import (file)        │                              │                          │
     ├───────────────────────────►│                              │                          │
     │                            │                              │                          │
     │                            │ Parse file                   │                          │
     │                            │ Detect encoding              │                          │
     │                            │ Map columns                  │                          │
     │                            │                              │                          │
     │                            │ For each row:                │                          │
     │                            │  - Validate account          │                          │
     │                            │  - Validate amount           │                          │
     │                            │  - Validate currency         │                          │
     │                            │  - Check duplicate           │                          │
     │                            │                              │                          │
     │  Preview (X valid, Y err)  │                              │                          │
     │◄───────────────────────────┤                              │                          │
     │                            │                              │                          │
     │  Confirm                   │                              │                          │
     ├───────────────────────────►│                              │                          │
     │                            │                              │                          │
     │                            │ BEGIN TRANSACTION            │                          │
     │                            │ For each valid row:         ├── BulkCreate()            │
     │                            ├─────────────────────────────►│── opening_balances       │
     │                            │                              │                          │
     │                            │ COMMIT                       │                          │
     │                            │                              │                          │
     │  Result (imported, errors) │                              │                          │
     │◄───────────────────────────┤                              │                          │
```

---

## 7. Correction Workflow

```
 Accountant              Service                    Chief Accountant              Database
 ──────────         ──────────────────             ────────────────           ──────────────
     │                      │                             │                        │
     │ Request correction   │                             │                        │
     ├─────────────────────►│                             │                        │
     │                      │ Validate reason             │                        │
     │                      │ Validate new amounts        │                        │
     │                      │                             │                        │
     │                      │ Set original → CORRECTED   ├── INSERT (new OB)      │
     │                      │ Create new OB → PENDING   ├────────────────────────►│
     │                      │ Link via correction_of     │                        │
     │                      │                             │                        │
     │  "Awaiting approval"  │  Notify approver           │                        │
     │◄─────────────────────┤───────────────────────────►│                        │
     │                      │                             │                        │
     │                      │                             │ Review & Approve       │
     │                      │                             ├────────────────────────►│
     │                      │                             │  Set new OB → APPROVED │
     │                      │◄────────────────────────────┤                        │
     │                      │                             │                        │
     │  Correction approved │  Audit log                  │                        │
     │◄─────────────────────┤─────────────────────────────┤                        │
```

---

## 8. Opening Balance Audit Trail

Every operation on opening balances generates an audit event:

| Event | Audit Action | Entity Type | Captures |
|-------|-------------|-------------|----------|
| Create | CREATE | OPENING_BALANCE | New OB values |
| Submit | SUBMIT | OPENING_BALANCE | Status change |
| Approve | APPROVE | OPENING_BALANCE | Approver, timestamp |
| Reject | REJECT | OPENING_BALANCE | Rejection reason |
| Correct | CORRECT | OPENING_BALANCE | Old vs new values |
| Delete | DELETE | OPENING_BALANCE | Delete reason |
| Carry-forward | CLOSE | FISCAL_YEAR | Summary of CF |
| Import | IMPORT | OPENING_BALANCE | Batch ID, row count |

---

## 9. Event Notifications

| Event | Notification | Channel | Recipient |
|-------|-------------|---------|-----------|
| OB submitted | "X OBs pending your approval" | In-app, email | Chief Accountant |
| OB approved | "OB for account X approved" | In-app | Creator |
| OB rejected | "OB for account X rejected: reason" | In-app | Creator |
| Correction requested | "Correction pending approval" | In-app, email | Chief Accountant |
| Carry-forward | "Fiscal year carry-forward completed" | In-app, email | Chief Accountant |
| Import completed | "Import completed: X success, Y errors" | In-app | Creator |
