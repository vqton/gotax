# TODOs

## Completed
- [x] AUTH module (login, 2FA, refresh tokens, rate limiting, sessions)
- [x] COA module (freeze, approvals, versioning, mappings, IFRS, analysis)
- [x] Opening Balance module + excelize + maroto
- [x] Clean Architecture migration + tax foundation
- [x] Cash module (receipts, payments, transfers, petty cash, inventory, advance settlement)
- [x] Bank module (statements, recon, payment orders, loans, term deposits, PG+memory repos)
- [x] DB migrations: 002 GL, 003 company, 004 bank schemas
- [x] Model files split by bounded context (7 files)
- [x] Company + bank review fixes
- [x] Purchase module: domain models + errors + interfaces (models_purchase.go, 401 lines)
- [x] Purchase module: 005_purchase_schema.sql migration (8 tables + indexes)
- [x] Purchase module: memory_purchase.go written
- [x] Purchase module: pg_purchase.go written
- [x] Purchase module: purchase_service.go written
- [x] Purchase module: purchase_handler.go written (25 endpoints)
- [x] Purchase module: main.go wired (PG + CA paths)

## In Progress
- [ ] **BLOCKER**: Fix interface naming collision — Go 1.26 can't overload `Create`/`GetByID`/`List`/`Update` across 6 repo types by return type alone
  - Must rename interfaces to prefixed names (`CreateSupplier`, `GetPO`, `CreateInvoice`, `GetCostAllocation`)
  - Align both memory_purchase.go + pg_purchase.go to new names
  - Align purchase_service.go call sites to new names
  - Verify `go build ./...` passes

## Blocked
- Build broken: `go build ./...` fails (~100 errors)
- Tests not yet written for purchase module (handler_test.go)

## Next
- [ ] Fix interface names + align repos + service (blocker)
- [ ] `go build ./...` green
- [ ] Write purchase handler tests (TDD)
- [ ] Run full test suite `go test ./...`
- [ ] Review module against MISA/Fast/BravoERP standard
