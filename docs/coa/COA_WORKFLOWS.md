# GoTax COA Module - Workflows & Data Flows

---

## 1. Account Lifecycle State Machine

```
                  ┌──────────┐
                  │ PENDING  │ (pending approval)
                  └────┬─────┘
                       │ approve
                       ▼
     ┌──────────┐  ┌──────┐  ┌──────────┐
     │ DRAFT    │──► ACTIVE│──► FROZEN   │
     │ (unused) │  └──┬───┘  └──────────┘
     └──────────┘     │        │
                      │        │ unfreeze
                      ▼        ▼
                  ┌──────────┐
                  │INACTIVE  │ (soft-deleted)
                  └──────────┘
```

**Transitions:**
- DRAFT → ACTIVE: on creation (no approval needed for zero-balance)
- DRAFT → PENDING: on creation requires approval → ACTIVE after approve
- ACTIVE → FROZEN: chief accountant freeze (direct)
- FROZEN → ACTIVE: chief accountant unfreeze (direct)
- ACTIVE → INACTIVE: deactivation (approval needed if balance > 0)
- FROZEN → INACTIVE: deactivation (unfreeze first)
- Any → INACTIVE: system on period end (zero-balance auto)

---

## 2. Account Creation Workflow

```
 ┌────────┐    ┌──────────┐    ┌──────────┐    ┌────────┐
 │User    │    │System    │    │System    │    │System  │
 │Request │───►│Validate  │───►│Check     │───►│Create  │
 │Create  │    │          │    │Balance=0?│    │Account │
 └────────┘    └──────────┘    └─────┬────┘    └───┬────┘
                                     │              │
                            ┌────────┴──────┐       │
                            │Yes            │No     │
                            ▼               ▼       │
                     ┌──────────┐    ┌──────────┐   │
                     │Direct    │    │Send to   │   │
                     │ACTIVE    │    │Approval  │   │
                     └──────────┘    │Queue     │   │
                                     └────┬─────┘   │
                                          ▼         │
                                    ┌──────────┐     │
                                    │Chief Acc │─────┤
                                    │Approve?  │     │
                                    └──┬───┬───┘     │
                              Yes/     │   │No       │
                              Approve  │   ▼         │
                                       │  Rejected   │
                                       ▼  (no op)    │
                                  ┌──────────┐       │
                                  │ACTIVE    │       │
                                  └──────────┘       │
                                                     ▼
                                              ┌──────────┐
                                              │Audit Log │
                                              │Recorded  │
                                              └──────────┘
```

---

## 3. COA Import Data Flow

```
 External CSV/Excel             System                    Database
 ──────────────────       ──────────────────       ─────────────────
         │                        │                        │
         │  Upload file           │                        │
         ├───────────────────────►│                        │
         │                        │                        │
         │                        ├── Parse headers         │
         │                        ├── Validate encoding     │
         │                        ├── Validate rows         │
         │                        │                        │
         │  Preview (X/Y/Z)       │                        │
         │◄───────────────────────┤                        │
         │                        │                        │
         │  Confirm import        │                        │
         ├───────────────────────►│                        │
         │                        │                        │
         │                        ├── BEGIN TRANSACTION    │
         │                        ├── For each valid row:  │
         │                        │   ├── Check parent     │
         │                        │   ├── Check duplicate  │
         │                        │   ├── Check hierarchy  │
         │                        │   ├── UPSERT account   │────► accounts
         │                        │   └── Audit log row    │────► audit_log
         │                        ├── COMMIT               │
         │                        │                        │
         │  Success response      │                        │
         │◄───────────────────────┤                        │
```

---

## 4. Account Balance Inquiry Data Flow

```
 User Request                         System                    Database
 ─────────────────            ────────────────────          ─────────────────
         │                           │                           │
         │ GET /accounts/{code}      │                           │
         │ /balance?period=P-2026-07 │                           │
         ├──────────────────────────►│                           │
         │                           ├── Get account info        │────► accounts
         │                           ├── Get opening balance     │────► account_balances
         │                           │   (from prev period)      │
         │                           ├── Get period activity     │────► journal_lines
         │                           │   (SUM debit/credit)      │     + journal_entries
         │                           ├── Calculate closing       │
         │                           │   balance                 │
         │                           │                           │
         │  Balance response         │                           │
         │◄──────────────────────────┤                           │
```

**Calculation:**
```
opening_debit   = balance(prev_period).closing_debit
opening_credit  = balance(prev_period).closing_credit
period_debit    = SUM(line.debit) WHERE account=code AND period=X AND status=POSTED
period_credit   = SUM(line.credit) WHERE account=code AND period=X AND status=POSTED
closing_balance = (opening_debit + period_debit) - (opening_credit + period_credit)
```

---

## 5. Account Approval Workflow (4-Eyes)

```
 Requester                 System                    Approver
 ─────────            ──────────────────           ───────────
     │                        │                         │
     │  Change request        │                         │
     ├───────────────────────►│                         │
     │                        ├── Create approval       │
     │                        │   record (PENDING)      │
     │                        ├── Notify approver       │
     │                        ├────────────────────────►│
     │                        │                         │
     │  "Awaiting approval"   │                         ├── Review diff
     │◄───────────────────────┤                         │   (current vs proposed)
     │                        │                         │
     │                        │                         ├── Approve or Reject
     │                        │                         ├────────────────►
     │                        │                         │
     │                        ├── If APPROVED:          │
     │                        │   ├── Apply change      │
     │                        │   └── Log audit         │
     │                        │                         │
     │                        ├── If REJECTED:          │
     │                        │   └── Log with reason   │
     │                        │                         │
     │  Result notification   │                         │
     │◄───────────────────────┤                         │
```

---

## 6. COA Data Model (Entity Relationship)

```
 ┌─────────────┐       ┌───────────────────┐
 │   tenants   │       │  accounts          │
 │─────────────│       │───────────────────│
 │ id (UUID)   │──┐    │ code (PK, VARCHAR) │
 │ tax_code    │  │    │ name               │
 │ name        │  │    │ name2              │
 └─────────────┘  │    │ type (ENUM)        │
                  │    │ parent_code (FK)   │── self-ref
                  │    │ is_active          │
                  │    │ is_parent          │
                  │    │ is_foreign         │
                  │    │ detail_by (ENUM)   │
                  │    │ status (active/    │
                  │    │   frozen/inactive) │
                  │    │ freeze_reason      │
                  │    │ version            │
                  │    │ effective_from     │
                  │    │ effective_to       │
                  │    │ created_at         │
                  │    │ updated_at         │
                  │    └────────┬───────────┘
                  │             │
                  │    ┌───────┴────────────┐
                  │    │ account_analysis    │
                  │    │────────────────────│
                  │    │ account_code (FK)   │
                  │    │ cost_center_id      │
                  │    │ profit_center_id    │
                  │    │ department_id       │
                  │    │ project_id          │
                  │    │ custom_dimension    │
                  │    └────────────────────┘
                  │
                  │    ┌───────────────────┐
                  ├────│ account_versions   │
                  │    │───────────────────│
                  │    │ id (PK)           │
                  │    │ version_number     │
                  │    │ snapshot (JSONB)   │
                  │    │ change_reason      │
                  │    │ created_by         │
                  │    │ created_at         │
                  │    └───────────────────┘
                  │
                  │    ┌───────────────────┐
                  ├────│ account_mappings   │
                  │    │───────────────────│
                  │    │ old_code           │
                  │    │ new_code (FK)      │
                  │    │ mapping_type       │
                  │    │ effective_from     │
                  │    └───────────────────┘
                  │
                  │    ┌───────────────────┐
                  └────│ account_ifrs_map   │
                       │───────────────────│
                       │ vas_code (FK)      │
                       │ ifrs_code          │
                       │ reclassification   │
                       │ adjustment_rule    │
                       └───────────────────┘
```

---

## 7. Business Rules Engine

### Account Code Rules
| Rule | Description |
|---|---|
| R001 | Code must be numeric only |
| R002 | Code min 3 characters, max 20 |
| R003 | First digit must match account type (1-2=ASSET, 3=LIABILITY, 4=EQUITY, 5-7=REVENUE, 6-8=EXPENSE) |
| R004 | Parent's code must be prefix of child code |
| R005 | Code must be unique within tenant |
| R006 | Code cannot be changed once journal entries reference it |

### Account Hierarchy Rules
| Rule | Description |
|---|---|
| R010 | A parent account cannot have direct journal postings |
| R011 | Only leaf accounts (is_parent=false) can have journal entries |
| R012 | Max hierarchy depth: unlimited (practical limit: 10) |
| R013 | Circular parent reference prohibited |
| R014 | Account type must match parent type (or be compatible) |

### Account State Rules
| Rule | Description |
|---|---|
| R020 | FROZEN accounts: no new debit/credit postings, balance reads OK |
| R021 | INACTIVE accounts: hidden from default views, history preserved |
| R022 | Cannot delete account with balance > 0 (must freeze then zero-balance) |
| R023 | Cannot delete account that is parent of existing accounts |
| R024 | Account type change: always requires approval |

### Import/Export Rules
| Rule | Description |
|---|---|
| R030 | Import is transactional: all-or-nothing |
| R031 | Import preserves existing accounts (upsert by code) |
| R032 | Import file max 10MB, max 10000 rows |
| R033 | Import preview is advisory, not a lock |
| R034 | Export respects user's permission scope |

### Versioning Rules
| Rule | Description |
|---|---|
| R040 | Version snapshotted on every change that modifies account structure |
| R041 | Name-only changes: versioned but minor bump |
| R042 | Structure changes (add/remove/reparent): major version bump |
| R043 | Versions are immutable once created |

---

## 8. User Journeys

### Journey 1: New Enterprise Onboarding
```
Day 1: Admin provisions GoTax tenant
     → System seeds Circular 99 COA (71 L1 accounts + sub-accounts)
     → System generates initial version (v1.0)

Day 1: Chief Accountant reviews standard COA
     → Modifies account names to match enterprise custom naming
     → Adds enterprise-specific detail accounts (level 4-5)
     → Generates IGAP document
     → Freezes accounts not relevant to this enterprise

Day 2: Import legacy data from old system
     → Uploads TT200 COA export
     → System maps TT200 → TT99 accounts
     → Chief Accountant reviews mapping
     → Confirms migration
     → System journals transfer entries
     → New COA version (v2.0) with migration history

Day 3: Accountants begin daily operations
     → Create journal entries using new COA
     → View account balances with drill-down
     → Export COA for reference
```

### Journey 2: Monthly COA Review
```
Day 1: Chief Accountant runs "Account Usage Report"
     → Sees accounts with zero activity for 3+ months
     → Freezes unused accounts
     → Creates approval requests for accounts needing modification

Day 2: Accountant receives freeze notifications
     → Reviews impact on recurring entries
     → Updates recurring entry templates if needed

Day 3: Chief Accountant reviews COA changes
     → Compares current COA to previous month version
     → Generates IGAP update if changes made
     → Signs off on monthly COA review
```

### Journey 3: Audit Response
```
Auditor requests COA change history for FY2026
     → Chief Accountant exports COA version comparison
     → Shows every account creation, modification, deactivation
     → Shows approval records for each change
     → Shows IGAP document with effective dates
     → Exports account balance snapshots at each period end
     → Auditor validates compliance with Circular 99
```

---

## 9. API Endpoints (COA-Specific)

| Method | Path | Description |
|---|---|---|
| GET | /api/v1/coa/accounts | List accounts (filterable: active, type, parent) |
| POST | /api/v1/coa/accounts | Create account |
| GET | /api/v1/coa/accounts/:code | Get account by code |
| PUT | /api/v1/coa/accounts/:code | Update account |
| DELETE | /api/v1/coa/accounts/:code | Deactivate account |
| POST | /api/v1/coa/accounts/:code/freeze | Freeze account |
| POST | /api/v1/coa/accounts/:code/unfreeze | Unfreeze account |
| GET | /api/v1/coa/accounts/:code/balance | Get account balance |
| GET | /api/v1/coa/accounts/:code/balance/drilldown | Drill-down to journals |
| GET | /api/v1/coa/accounts/:code/history | Account change history |
| POST | /api/v1/coa/import | Import COA (CSV/Excel) |
| GET | /api/v1/coa/export | Export COA (CSV/Excel/PDF) |
| GET | /api/v1/coa/versions | List COA versions |
| GET | /api/v1/coa/versions/:v1/compare/:v2 | Compare versions |
| POST | /api/v1/coa/approvals | Create approval request |
| GET | /api/v1/coa/approvals/:id | Get approval request |
| POST | /api/v1/coa/approvals/:id/approve | Approve |
| POST | /api/v1/coa/approvals/:id/reject | Reject |
| POST | /api/v1/coa/igap/generate | Generate IGAP PDF |
| GET | /api/v1/coa/mappings | List COA mappings (TT200→TT99) |
| POST | /api/v1/coa/mappings | Create mapping entry |
| GET | /api/v1/coa/analysis/:code | Get analysis codes for account |
| PUT | /api/v1/coa/analysis/:code | Update analysis codes |
