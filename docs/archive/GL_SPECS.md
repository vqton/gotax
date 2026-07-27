# GL Module Specification

## Overview
General Ledger (Sổ Cái) module for GoTax. Covers Chart of Accounts, Journal Entry, Ledger, Trial Balance, and Financial Reports. Follows Circular 200/2014/TT-BTC and Circular 99/2025/TT-BTC.

## Architecture
```
handler (HTTP) → service (business logic) → repository (data access)
```
- **Handler**: Gin HTTP handlers, request/response serialization
- **Service**: Business rules (double-entry validation, account status checks, balance calculation)
- **Repository**: Interface + in-memory implementation (PostgreSQL migration ready)

## Core Entities

### Account (Tai khoan ke toan)
- `code`: Digit string, ≥3 chars, unique (e.g. "1111")
- `name`: Vietnamese name
- `name2`: English/secondary name
- `type`: ASSET | LIABILITY | EQUITY | REVENUE | EXPENSE
- `parent_code`: Hierarchical (e.g. "1111" parent "111")
- `is_active`: Allow/block posting
- `is_foreign`: Foreign currency flag
- `detail_by`: OBJECT | PROJECT | CONTRACT | COST_ITEM | DEPARTMENT
- `is_parent`: Has children accounts

### JournalEntry (Chung tu ghi so)
- `entry_date`: Transaction date
- `period_id`: Accounting period FK
- `description`: Required
- `status`: DRAFT → POSTED → CANCELLED
- `lines`: 2+ lines, debit=credit (double-entry)
- `posted_at`: Timestamp when posted

### JournalLine (Chi tiet chung tu)
- `account_code`: Account FK
- `debit_amount` / `credit_amount`: One must be > 0, both >= 0
- `description`: Per-line description
- `object_id`, `project_id`, `contract_id`, `cost_item_id`, `department_id`: Detail tracking

### Period (Ky ke toan)
- `year`, `month`: Logical key
- `start_date`, `end_date`: Date range
- `status`: OPEN | CLOSED | LOCKED

### AccountBalance (So du TK)
- `open_balance_debit/credit`, `period_debit/credit`
- `Calculate()` derives `total_debit`, `total_credit`, `closing_balance`

## API Endpoints

### Accounts
| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/accounts | Create account |
| GET | /api/v1/accounts | List accounts (?active=true) |
| GET | /api/v1/accounts/:code | Get by code |
| PUT | /api/v1/accounts/:code | Update (code must match body) |
| DELETE | /api/v1/accounts/:code | Delete (fails if children) |

### Journal Entries
| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/journal-entries | Create + post entry |
| GET | /api/v1/journal-entries | List by date range (?from=&to=) |
| GET | /api/v1/journal-entries/:id | Get by ID |
| POST | /api/v1/journal-entries/:id/cancel | Cancel posted entry |

### Reports
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/reports/trial-balance | Trial balance (?year=&month=) |

## Business Rules
1. Account code must be digits, ≥3 characters
2. Account type must be one of 5 types
3. Journal entry must have ≥1 line
4. Total debit must equal total credit (within 0.001 tolerance)
5. All referenced accounts must exist and be active
6. Cannot post already-posted or cancelled entries
7. Cannot cancel draft entries (only posted)
8. Cannot delete account with children
9. Trial balance only counts POSTED entries

## Migration
`migrations/001_gl_schema.sql` contains:
- Tables: periods, accounts, journal_entries, journal_lines, account_balances
- Constraints, indexes, foreign keys
- Seed data: 12 monthly periods for 2026, 40+ accounts per Circular 200

## Deviations from MISA/FAST/BRAVO (Current Version)
| Feature | MISA/FAST/BRAVO | GoTax GL v1 | Planned |
|---------|----------------|-------------|---------|
| Budget tracking | Yes | No | v2 |
| Multi-currency auto-conversion | Yes | Manual fields | v2 |
| Bank/e-tax integration | Yes | No | v3 |
| IFRS dual reporting | BRAVO | No | v2 |
| AI transaction matching | MISA AVA | No | v3 |
| Financial statements (BS/IS/CF) | Yes | Trial balance only | v2 |
| Period close (ket chuyen) | Yes | Status field only | v2 |
| Audit trail | Full log | None | v2 |
| Approval workflow | Yes | Draft→Posted only | v2 |

## Test Coverage
- **Domain models**: 22 test cases (validation, balance calc, normal balance)
- **Service layer**: 19 test cases (CRUD, posting, cancel, trial balance, error paths)
- **Handler layer**: 20 test cases (HTTP success/failure, error mapping)
- **Repository integration**: 17 test cases (CRUD, trial balance, workflow)

Total: 78 test cases, all passing.