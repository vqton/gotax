# Implementation Plan: GoTax Pending Features

## Overview

Vietnamese tax-compliant GL API needs 12 modules implemented. Current state: purchase module has interface naming collision breaking build. Plan: fix blocker first, then implement modules from simple to complex, core to edge.

## Architecture Decisions

- **Follow AGENTS.md step order**: Interface → GORM model → Repository → Service → Handler → Wiring → Tests
- **TDD approach**: Write failing test first, then implementation
- **Two backends**: Always implement PG + memory repos
- **Vertical slicing**: Complete feature paths, not horizontal layers
- **Review after each phase**: Code review before moving to next phase

## Task List

### Phase 0: Fix Immediate Blocker (1 day)

- [ ] Task 0.1: Rename purchase interfaces to prefixed names
- [ ] Task 0.2: Align memory_purchase.go to new interface names
- [ ] Task 0.3: Align pg_purchase.go to new interface names
- [ ] Task 0.4: Align purchase_service.go call sites
- [ ] Task 0.5: Verify `go build ./...` passes
- [ ] Task 0.6: Write purchase handler tests (TDD)
- [ ] Task 0.7: Run full test suite `go test -count=1 ./...`

**Checkpoint: Phase 0**
- [ ] Build green
- [ ] Purchase module tests pass
- [ ] Code review

### Phase 1: Foundation Modules (5 days)

#### Module 1: Number Format (1 day)
- [ ] Task 1.1: Create `internal/format/format.go` with `FormatNumber`
- [ ] Task 1.2: Implement `ParseNumber`
- [ ] Task 1.3: Write unit tests
- [ ] Task 1.4: Integrate with existing handlers

#### Module 2: System Options (2 days)
- [ ] Task 2.1: Create migration `000035_system_options.up.sql` + `.down.sql`
- [ ] Task 2.2: Create `models_system.go` with `SystemOption` struct
- [ ] Task 2.3: Create `SystemOptionRepo` interface in `interfaces.go`
- [ ] Task 2.4: Define error variables in `errors.go`
- [ ] Task 2.5: Implement `pg_system_option.go`
- [ ] Task 2.6: Implement `memory_system_option.go`
- [ ] Task 2.7: Create `system_option_service.go`
- [ ] Task 2.8: Implement Get/Set by category
- [ ] Task 2.9: Implement defaults initialization
- [ ] Task 2.10: Create `system_option_handler.go`
- [ ] Task 2.11: Register routes in `handler.go`
- [ ] Task 2.12: Wire repos/services in `main.go` (PG + memory)
- [ ] Task 2.13: Write service tests
- [ ] Task 2.14: Write handler tests

#### Module 3: Voucher Numbering (2 days)
- [ ] Task 3.1: Create migration `000036_numbering_rules.up.sql` + `.down.sql`
- [ ] Task 3.2: Create `models_numbering.go`
- [ ] Task 3.3: Create `NumberingRuleRepo` interface
- [ ] Task 3.4: Define errors
- [ ] Task 3.5: Implement PG repo
- [ ] Task 3.6: Implement Memory repo
- [ ] Task 3.7: Create `numbering_rule_service.go`
- [ ] Task 3.8: Implement `GetNextNumber` with atomic increment
- [ ] Task 3.9: Create handler
- [ ] Task 3.10: Wire in main.go
- [ ] Task 3.11: Write tests

**Checkpoint: Phase 1**
- [ ] All foundation modules working
- [ ] Tests pass
- [ ] Code review

### Phase 2: Core Business Modules (8 days)

#### Module 4: Fiscal Year (2 days)
- [ ] Task 4.1: Extend `periods` table with fiscal_year fields
- [ ] Task 4.2: Add fiscal year config to SystemOption
- [ ] Task 4.3: Auto-generate 12 periods on FY creation
- [ ] Task 4.4: Integrate with existing period handler
- [ ] Task 4.5: Write tests

#### Module 5: Contracts (3 days)
- [ ] Task 5.1: Create migration `000037_contracts.up.sql` + `.down.sql`
- [ ] Task 5.2: Create `models_contract.go`
- [ ] Task 5.3: Create `ContractRepo` interface
- [ ] Task 5.4: Define errors
- [ ] Task 5.5: Implement PG repo
- [ ] Task 5.6: Implement Memory repo
- [ ] Task 5.7: Create `contract_service.go`
- [ ] Task 5.8: Implement CRUD + status transitions
- [ ] Task 5.9: Create handler
- [ ] Task 5.10: Register routes
- [ ] Task 5.11: Wire in main.go
- [ ] Task 5.12: Write tests

#### Module 6: Loan Agreements (3 days)
- [ ] Task 6.1: Create migration `000038_loan_agreements.up.sql` + `.down.sql`
- [ ] Task 6.2: Create `models_loan.go`
- [ ] Task 6.3: Create `LoanAgreementRepo` interface
- [ ] Task 6.4: Define errors
- [ ] Task 6.5: Implement PG repo
- [ ] Task 6.6: Implement Memory repo
- [ ] Task 6.7: Create `loan_service.go`
- [ ] Task 6.8: Implement amortization schedule
- [ ] Task 6.9: Implement payment recording
- [ ] Task 6.10: Create handler
- [ ] Task 6.11: Register routes
- [ ] Task 6.12: Wire in main.go
- [ ] Task 6.13: Write tests

**Checkpoint: Phase 2**
- [ ] Core business modules working
- [ ] Tests pass
- [ ] Code review

### Phase 3: Integration Modules (7 days)

#### Module 7: E-Banking (4 days)
- [ ] Task 7.1: Design CSV parser interface
- [ ] Task 7.2: Implement VCB parser
- [ ] Task 7.3: Implement BIDV parser
- [ ] Task 7.4: Implement CTG parser
- [ ] Task 7.5: Implement VTB parser
- [ ] Task 7.6: Implement ACB parser
- [ ] Task 7.7: Create `bank_import_service.go`
- [ ] Task 7.8: Implement reconciliation logic
- [ ] Task 7.9: Create handler
- [ ] Task 7.10: Write tests

#### Module 8: E-Tax Filing (2 days)
- [ ] Task 8.1: VAT declaration XML (Mẫu 01)
- [ ] Task 8.2: CIT declaration XML (Mẫu 01)
- [ ] Task 8.3: PIT declaration XML
- [ ] Task 8.4: Create `tax_filing_service.go`
- [ ] Task 8.5: Integrate with GDT client
- [ ] Task 8.6: Create handler
- [ ] Task 8.7: Write tests

#### Module 9: Digital Signature (1 day)
- [ ] Task 9.1: Design `SignatureProvider` interface
- [ ] Task 9.2: Implement VNPT-CA provider
- [ ] Task 9.3: Implement mock provider for testing
- [ ] Task 9.4: Create `signature_service.go`
- [ ] Task 9.5: Write tests

**Checkpoint: Phase 3**
- [ ] Integration modules working
- [ ] Tests pass
- [ ] Code review

### Phase 4: Polish Modules (4 days)

#### Module 10: Backup & Restore (2 days)
- [ ] Task 10.1: Create migration `000039_backups.up.sql`
- [ ] Task 10.2: Create `backup_service.go`
- [ ] Task 10.3: Implement pg_dump wrapper
- [ ] Task 10.4: Implement restore logic
- [ ] Task 10.5: Create handler
- [ ] Task 10.6: Write tests

#### Module 11: Report Customization (1 day)
- [ ] Task 11.1: Add report options to SystemOption
- [ ] Task 11.2: Extend PDF generator
- [ ] Task 11.3: Add logo upload
- [ ] Task 11.4: Write tests

#### Module 12: Multi-branch (1 day)
- [ ] Task 12.1: Branch-level options support
- [ ] Task 12.2: Per-branch numbering rules
- [ ] Task 12.3: Branch report filtering
- [ ] Task 12.4: Write tests

**Checkpoint: Phase 4**
- [ ] All modules complete
- [ ] Full test suite passing
- [ ] Final code review

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Interface naming collision | High | Fix immediately, rename to prefixed names |
| Build broken | High | Verify after each change |
| Missing tests | Medium | TDD approach, write tests with implementation |
| Complex business logic | Medium | Follow Circular 99/Decree 123 regulations |
| Performance issues | Low | Profile after implementation |

## Open Questions

1. Should we implement all E-Banking parsers or start with VCB only?
2. Digital Signature: use real VNPT-CA or mock only?
3. Backup: implement pg_dump wrapper or use Go native?

## Verification

After each phase:
1. `go build ./...` passes
2. `go vet ./...` passes
3. `go test -count=1 ./...` passes
4. Code review completed
5. Manual verification of key workflows
