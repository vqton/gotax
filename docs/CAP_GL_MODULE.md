# GL Module — Critical Action Plan (CAP)

## Stop-Gap: What to Fix NOW Before Any Production Use

### Critical (Week 1-2)
```bash
# 1. Replace Circular 200 COA → Circular 99
# 2. Add PostgreSQL persistence driver
# 3. Add audit log table + append-only middleware
# 4. Add JWT auth + RBAC middleware
```
- Update `migrations/001_gl_schema.sql` per Appendix B of BRD
- 71 Level-1 accounts per Circular 99, NOT 76 per Circular 200
- Add audit_log table with: id, user_id, action, entity_type, entity_id, old_value, new_value, ip_address, created_at
- Add users table with roles: admin, chief_accountant, accountant, viewer

### High Priority (Week 3-6)
- Financial Statements (B01-B04) with auto data pull from trial balance
- Period close workflow (OPEN→CLOSING→CLOSED with auto closing entries)
- Approval workflow (DRAFT→REVIEWING→APPROVED→POSTED)
- Closing entry templates (511→911, 911→421, etc.)

### Medium Priority (Week 7-12)
- Multi-currency with exchange rate table and year-end revaluation
- Sub-ledgers (customer/vendor/project tracking)
- Enhanced reporting (Sổ Cái, Sổ chi tiết in Circular 99 format)

## Go/No-Go Checklist
```
[ ] Circular 99 compliant chart of accounts
[ ] PostgreSQL persistence (no data loss)
[ ] Audit trail (100% of changes logged)
[ ] Authentication + Role-based access control
[ ] Financial Statements B01-B04
[ ] Period close workflow
[ ] Multi-currency transactions
[ ] Sub-ledger tracking
[ ] 50+ concurrent users tested
[ ] Data export (Excel, PDF)
[ ] DR plan (backup, restore tested)
→ ALL items must pass before PROD
```

## Current Module Score
```
Domain Models:     ██████████  10/10
Validation Logic:  ██████████  10/10
Unit Tests:        ██████████  10/10
API Layer:         ████████░░  8/10 (missing period + auth)
Business Logic:    ███████░░░  7/10 (basic posting, no workflow)
Data Persistence:  ██░░░░░░░░  2/10 (in-memory only)
Security:          ██░░░░░░░░  2/10 (no auth, no audit, no RBAC)
Feature Complete:  ██░░░░░░░░  2/10 (GL only, no FS/close/multicurrency)
Regulatory:        █░░░░░░░░░  1/10 (Circular 200 instead of 99)
Production Ready:  █░░░░░░░░░  1/10

TOTAL:                  43/100
PASSING THRESHOLD:      80/100
GAP:                    37 points
```

## Files Created
| File | Description |
|------|-------------|
| `docs/BRD_GL_MODULE.md` | Full BRD: regulatory analysis, market comparison, 70+ functional requirements, 4+ business rules, 3 workflows, 3 user journeys, architecture roadmap |
| `docs/UC_GL_MODULE.md` | 12 use cases with happy/alternative/exception paths, state machines for Journal Entry, Period, Account |
| `docs/CAP_GL_MODULE.md` | This file — critical action plan & go/no-go checklist |
| `docs/GL_SPECS.md` | Updated spec reflecting gap analysis |
