# Implementation Plan: GoTax Remaining Modules

## Overview

Four modules remain from TASK_BREAKDOWN_ADMIN.md: E-Banking, E-Tax Filing, Digital Signature, Multi-branch. Build green, tests pass. Implement from simple to complex, core to edge.

## Architecture Decisions

- **Follow AGENTS.md step order**: Interface → GORM model → Repository → Service → Handler → Wiring → Tests
- **TDD approach**: Write failing test first, then implementation
- **Two backends**: Always implement PG + memory repos
- **Vertical slicing**: Complete feature paths, not horizontal layers
- **Review after each module**: Code review before moving to next

## Task List

### Module 1: E-Banking (CSV Parsers)

**Simplest module — pure parsing, no complex business logic.**

#### Task 1.1: Design CSV parser interface
- **Description:** Create interface for bank CSV parsers
- **Acceptance criteria:**
  - [ ] Interface defined in `internal/domain/interfaces.go`
  - [ ] Common methods: `Parse`, `Validate`, `GetBankCode`
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** None
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** S

#### Task 1.2: Implement VCB parser
- **Description:** Implement Vietcombank CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented in `internal/einvoice/parser_vcb.go`
  - [ ] Handles VCB CSV format
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/einvoice/` passes
- **Dependencies:** Task 1.1
- **Files likely touched:**
  - `internal/einvoice/parser_vcb.go`
  - `internal/einvoice/parser_vcb_test.go`
- **Estimated scope:** M

#### Task 1.3: Implement BIDV parser
- **Description:** Implement BIDV CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles BIDV CSV format
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/einvoice/` passes
- **Dependencies:** Task 1.1
- **Files likely touched:**
  - `internal/einvoice/parser_bidv.go`
  - `internal/einvoice/parser_bidv_test.go`
- **Estimated scope:** M

#### Task 1.4: Implement CTG parser
- **Description:** Implement VietinBank CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles CTG CSV format
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/einvoice/` passes
- **Dependencies:** Task 1.1
- **Files likely touched:**
  - `internal/einvoice/parser_ctg.go`
  - `internal/einvoice/parser_ctg_test.go`
- **Estimated scope:** M

#### Task 1.5: Implement VTB parser
- **Description:** Implement VPBank CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles VTB CSV format
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/einvoice/` passes
- **Dependencies:** Task 1.1
- **Files likely touched:**
  - `internal/einvoice/parser_vtb.go`
  - `internal/einvoice/parser_vtb_test.go`
- **Estimated scope:** M

#### Task 1.6: Implement ACB parser
- **Description:** Implement ACB CSV parser
- **Acceptance criteria:**
  - [ ] Parser implemented
  - [ ] Handles ACB CSV format
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/einvoice/` passes
- **Dependencies:** Task 1.1
- **Files likely touched:**
  - `internal/einvoice/parser_acb.go`
  - `internal/einvoice/parser_acb_test.go`
- **Estimated scope:** M

#### Task 1.7: Create bank import service
- **Description:** Create service for bank import
- **Acceptance criteria:**
  - [ ] Service implemented in `internal/service/bank_import_service.go`
  - [ ] Parser selection logic
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/service/` passes
- **Dependencies:** Task 1.2-1.6
- **Files likely touched:**
  - `internal/service/bank_import_service.go`
  - `internal/service/bank_import_service_test.go`
- **Estimated scope:** M

#### Task 1.8: Create bank import handler
- **Description:** Create HTTP handler for bank import
- **Acceptance criteria:**
  - [ ] Handler implemented in `internal/handler/bank_import_handler.go`
  - [ ] Routes registered
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/handler/` passes
- **Dependencies:** Task 1.7
- **Files likely touched:**
  - `internal/handler/bank_import_handler.go`
  - `internal/handler/bank_import_handler_test.go`
- **Estimated scope:** M

#### Task 1.9: Wire in main.go
- **Description:** Wire bank import dependencies in main.go
- **Acceptance criteria:**
  - [ ] PG branch wired
  - [ ] Memory branch wired
  - [ ] Service initialized
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 1.8
- **Files likely touched:**
  - `main.go`
- **Estimated scope:** S

**Checkpoint: E-Banking Module**
- [ ] All parsers working
- [ ] Service and handler tests pass
- [ ] Code review

---

### Module 2: Digital Signature

**Simple module — provider pattern, minimal business logic.**

#### Task 2.1: Design SignatureProvider interface
- **Description:** Create interface for digital signature providers
- **Acceptance criteria:**
  - [ ] Interface defined in `internal/domain/interfaces.go`
  - [ ] Methods: `Sign`, `Verify`, `GetProviderName`
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** None
- **Files likely touched:**
  - `internal/domain/interfaces.go`
- **Estimated scope:** S

#### Task 2.2: Implement mock provider
- **Description:** Implement mock signature provider for testing
- **Acceptance criteria:**
  - [ ] Mock provider implemented in `internal/xmldsig/mock_provider.go`
  - [ ] Returns predictable results
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/xmldsig/` passes
- **Dependencies:** Task 2.1
- **Files likely touched:**
  - `internal/xmldsig/mock_provider.go`
  - `internal/xmldsig/mock_provider_test.go`
- **Estimated scope:** S

#### Task 2.3: Implement VNPT-CA provider
- **Description:** Implement VNPT-CA signature provider
- **Acceptance criteria:**
  - [ ] Provider implemented in `internal/xmldsig/vnpt_ca.go`
  - [ ] Handles VNPT-CA protocol
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/xmldsig/` passes
- **Dependencies:** Task 2.1
- **Files likely touched:**
  - `internal/xmldsig/vnpt_ca.go`
  - `internal/xmldsig/vnpt_ca_test.go`
- **Estimated scope:** L

#### Task 2.4: Create signature service
- **Description:** Create service for digital signatures
- **Acceptance criteria:**
  - [ ] Service implemented in `internal/service/signature_service.go`
  - [ ] Provider selection logic
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/service/` passes
- **Dependencies:** Task 2.2, 2.3
- **Files likely touched:**
  - `internal/service/signature_service.go`
  - `internal/service/signature_service_test.go`
- **Estimated scope:** M

**Checkpoint: Digital Signature Module**
- [ ] Providers working
- [ ] Service tests pass
- [ ] Code review

---

### Module 3: E-Tax Filing

**Complex module — XML generation, GDT integration.**

#### Task 3.1: Implement VAT declaration XML
- **Description:** Implement VAT declaration XML generation (Mẫu 01)
- **Acceptance criteria:**
  - [ ] XML generation in `internal/htkk/vat_declaration.go`
  - [ ] Follows Mẫu 01 format
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/htkk/` passes
- **Dependencies:** None
- **Files likely touched:**
  - `internal/htkk/vat_declaration.go`
  - `internal/htkk/vat_declaration_test.go`
- **Estimated scope:** M

#### Task 3.2: Implement CIT declaration XML
- **Description:** Implement CIT declaration XML generation (Mẫu 01)
- **Acceptance criteria:**
  - [ ] XML generation in `internal/htkk/cit_declaration.go`
  - [ ] Follows Mẫu 01 format
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/htkk/` passes
- **Dependencies:** Task 3.1
- **Files likely touched:**
  - `internal/htkk/cit_declaration.go`
  - `internal/htkk/cit_declaration_test.go`
- **Estimated scope:** M

#### Task 3.3: Implement PIT declaration XML
- **Description:** Implement PIT declaration XML generation
- **Acceptance criteria:**
  - [ ] XML generation in `internal/htkk/pit_declaration.go`
  - [ ] Follows PIT format
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/htkk/` passes
- **Dependencies:** Task 3.2
- **Files likely touched:**
  - `internal/htkk/pit_declaration.go`
  - `internal/htkk/pit_declaration_test.go`
- **Estimated scope:** M

#### Task 3.4: Create tax filing service
- **Description:** Create service for tax filing
- **Acceptance criteria:**
  - [ ] Service implemented in `internal/service/tax_filing_service.go`
  - [ ] XML generation logic
  - [ ] GDT client integration
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/service/` passes
- **Dependencies:** Task 3.1-3.3
- **Files likely touched:**
  - `internal/service/tax_filing_service.go`
  - `internal/service/tax_filing_service_test.go`
- **Estimated scope:** L

#### Task 3.5: Create tax filing handler
- **Description:** Create HTTP handler for tax filing
- **Acceptance criteria:**
  - [ ] Handler implemented in `internal/handler/tax_filing_handler.go`
  - [ ] Routes registered
  - [ ] Unit tests pass
- **Verification:**
  - [ ] `go test -v ./internal/handler/` passes
- **Dependencies:** Task 3.4
- **Files likely touched:**
  - `internal/handler/tax_filing_handler.go`
  - `internal/handler/tax_filing_handler_test.go`
- **Estimated scope:** M

**Checkpoint: E-Tax Filing Module**
- [ ] XML generation working
- [ ] Service and handler tests pass
- [ ] Code review

---

### Module 4: Multi-branch

**Complex module — cross-cutting concerns, system options extension.**

#### Task 4.1: Add branch-level options support
- **Description:** Add branch-level system options
- **Acceptance criteria:**
  - [ ] Branch options schema in system options
  - [ ] Override logic
  - [ ] Inheritance from company
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** None
- **Files likely touched:**
  - `internal/service/system_option_service.go`
- **Estimated scope:** M

#### Task 4.2: Implement per-branch numbering rules
- **Description:** Implement per-branch numbering rules
- **Acceptance criteria:**
  - [ ] Branch-specific rules
  - [ ] Fallback to company rules
  - [ ] Atomic generation
- **Verification:**
  - [ ] Unit tests pass
- **Dependencies:** Task 4.1
- **Files likely touched:**
  - `internal/service/numbering_rule_service.go`
- **Estimated scope:** M

#### Task 4.3: Add branch report filtering
- **Description:** Add branch filtering to reports
- **Acceptance criteria:**
  - [ ] Branch parameter support
  - [ ] Filtering logic
  - [ ] Performance optimization
- **Verification:**
  - [ ] `go build ./...` passes
- **Dependencies:** Task 4.2
- **Files likely touched:**
  - `internal/service/report_service.go`
- **Estimated scope:** M

#### Task 4.4: Write tests
- **Description:** Write tests for multi-branch functionality
- **Acceptance criteria:**
  - [ ] Service tests
  - [ ] Handler tests
  - [ ] Edge cases covered
- **Verification:**
  - [ ] `go test -count=1 ./...` passes
- **Dependencies:** Task 4.3
- **Files likely touched:**
  - `internal/service/system_option_service_test.go`
  - `internal/service/numbering_rule_service_test.go`
  - `internal/service/report_service_test.go`
- **Estimated scope:** M

**Checkpoint: Multi-branch Module**
- [ ] All multi-branch features working
- [ ] Tests pass
- [ ] Code review

---

### Final Checkpoint

- [ ] All modules complete
- [ ] Full test suite passing
- [ ] Build green
- [ ] Final code review
- [ ] Git sync

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Complex XML formats | Medium | Follow existing HTKK patterns |
| GDT API changes | Low | Use existing GDT client |
| Performance issues | Low | Profile after implementation |

## Verification

After each module:
1. `go build ./...` passes
2. `go vet ./...` passes
3. `go test -count=1 ./...` passes
4. Code review completed
5. Manual verification of key workflows
