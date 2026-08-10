# Task Breakdown — Admin/Config Modules

**Document ID:** GOTAX-TASKS-ADMIN-001  
**Version:** 1.0  
**Date:** 2026-08-07  
**Total Tasks:** 156  
**Estimated Effort:** 82 days

---

## Module 1: System Options (FR-SYS-001)

### 1.1 Database Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| SYS-DB-001 | Create migration `000035_system_options.up.sql` | 2h | TODO |
| SYS-DB-002 | Create migration `000035_system_options.down.sql` | 1h | TODO |
| SYS-DB-003 | Add indexes for company_id + category | 1h | TODO |

### 1.2 Domain Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| SYS-DOM-001 | Create `models_system.go` with `SystemOption` struct | 2h | TODO |
| SYS-DOM-002 | Add validation tags | 1h | TODO |
| SYS-DOM-003 | Create `SystemOptionRepo` interface in `interfaces.go` | 1h | TODO |
| SYS-DOM-004 | Define error variables in `errors.go` | 0.5h | TODO |

### 1.3 Repository Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| SYS-REP-001 | Implement `pg_system_option.go` | 4h | TODO |
| SYS-REP-002 | Implement `memory_system_option.go` | 3h | TODO |
| SYS-REP-003 | Test PG repo with integration test | 2h | TODO |

### 1.4 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| SYS-SVC-001 | Create `system_option_service.go` | 4h | TODO |
| SYS-SVC-002 | Implement Get/Set by category | 2h | TODO |
| SYS-SVC-003 | Implement defaults initialization | 2h | TODO |
| SYS-SVC-004 | Write service tests | 3h | TODO |

### 1.5 Handler Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| SYS-HDL-001 | Create `system_option_handler.go` | 3h | TODO |
| SYS-HDL-002 | Register routes in `handler.go` | 1h | TODO |
| SYS-HDL-003 | Write handler tests | 3h | TODO |

### 1.6 Integration
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| SYS-INT-001 | Wire repos/services in `main.go` (PG branch) | 1h | TODO |
| SYS-INT-002 | Wire repos/services in `main.go` (memory branch) | 1h | TODO |
| SYS-INT-003 | Add to `RegisterRoutesWithCompany` | 1h | TODO |

**Module Total:** 33.5h (~4.2 days)

---

## Module 2: Voucher Numbering (FR-NUM-001)

### 2.1 Database Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| NUM-DB-001 | Create migration `000036_numbering_rules.up.sql` | 2h | TODO |
| NUM-DB-002 | Create migration `000036_numbering_rules.down.sql` | 1h | TODO |

### 2.2 Domain Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| NUM-DOM-001 | Create `models_numbering.go` | 2h | TODO |
| NUM-DOM-002 | Create `NumberingRuleRepo` interface | 1h | TODO |
| NUM-DOM-003 | Define errors | 0.5h | TODO |

### 2.3 Repository Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| NUM-REP-001 | Implement PG repo | 3h | TODO |
| NUM-REP-002 | Implement Memory repo | 2h | TODO |

### 2.4 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| NUM-SVC-001 | Create `numbering_rule_service.go` | 4h | TODO |
| NUM-SVC-002 | Implement `GetNextNumber` with atomic increment | 3h | TODO |
| NUM-SVC-003 | Write tests | 3h | TODO |

### 2.5 Handler Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| NUM-HDL-001 | Create handler | 3h | TODO |
| NUM-HDL-002 | Write tests | 3h | TODO |

### 2.6 Integration
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| NUM-INT-001 | Wire in main.go | 1h | TODO |

**Module Total:** 29.5h (~3.7 days)

---

## Module 3: Fiscal Year (FR-FY-001)

### 3.1 Database Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| FY-DB-001 | Extend `periods` table with fiscal_year fields | 2h | TODO |

### 3.2 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| FY-SVC-001 | Add fiscal year config to SystemOption | 2h | TODO |
| FY-SVC-002 | Auto-generate 12 periods on FY creation | 4h | TODO |
| FY-SVC-003 | Write tests | 3h | TODO |

### 3.3 Integration
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| FY-INT-001 | Integrate with existing period handler | 2h | TODO |

**Module Total:** 13h (~1.6 days)

---

## Module 4: Number Format (FR-SYS-004)

### 4.1 Utility Package
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| NF-UTL-001 | Create `internal/format/format.go` | 3h | TODO |
| NF-UTL-002 | Implement `FormatNumber` | 2h | TODO |
| NF-UTL-003 | Implement `ParseNumber` | 2h | TODO |
| NF-UTL-004 | Write unit tests | 2h | TODO |

**Module Total:** 9h (~1.1 days)

---

## Module 5: Contracts (FR-CON-001)

### 5.1 Database Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| CON-DB-001 | Create migration `000037_contracts.up.sql` | 2h | TODO |
| CON-DB-002 | Create migration `000037_contracts.down.sql` | 1h | TODO |

### 5.2 Domain Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| CON-DOM-001 | Create `models_contract.go` | 3h | TODO |
| CON-DOM-002 | Create `ContractRepo` interface | 1h | TODO |
| CON-DOM-003 | Define errors | 0.5h | TODO |

### 5.3 Repository Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| CON-REP-001 | Implement PG repo | 4h | TODO |
| CON-REP-002 | Implement Memory repo | 3h | TODO |

### 5.4 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| CON-SVC-001 | Create `contract_service.go` | 4h | TODO |
| CON-SVC-002 | Implement CRUD + status transitions | 3h | TODO |
| CON-SVC-003 | Write tests | 3h | TODO |

### 5.5 Handler Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| CON-HDL-001 | Create handler | 4h | TODO |
| CON-HDL-002 | Register routes | 1h | TODO |
| CON-HDL-003 | Write tests | 4h | TODO |

### 5.6 Integration
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| CON-INT-001 | Wire in main.go | 1h | TODO |

**Module Total:** 34.5h (~4.3 days)

---

## Module 6: Loan Agreements (FR-LOAN-001)

### 6.1 Database Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| LOAN-DB-001 | Create migration `000038_loan_agreements.up.sql` | 2h | TODO |
| LOAN-DB-002 | Create migration `000038_loan_agreements.down.sql` | 1h | TODO |

### 6.2 Domain Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| LOAN-DOM-001 | Create `models_loan.go` | 3h | TODO |
| LOAN-DOM-002 | Create `LoanAgreementRepo` interface | 1h | TODO |
| LOAN-DOM-003 | Define errors | 0.5h | TODO |

### 6.3 Repository Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| LOAN-REP-001 | Implement PG repo | 4h | TODO |
| LOAN-REP-002 | Implement Memory repo | 3h | TODO |

### 6.4 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| LOAN-SVC-001 | Create `loan_service.go` | 4h | TODO |
| LOAN-SVC-002 | Implement amortization schedule | 4h | TODO |
| LOAN-SVC-003 | Implement payment recording | 3h | TODO |
| LOAN-SVC-004 | Write tests | 4h | TODO |

### 6.5 Handler Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| LOAN-HDL-001 | Create handler | 4h | TODO |
| LOAN-HDL-002 | Register routes | 1h | TODO |
| LOAN-HDL-003 | Write tests | 4h | TODO |

### 6.6 Integration
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| LOAN-INT-001 | Wire in main.go | 1h | TODO |

**Module Total:** 39.5h (~4.9 days)

---

## Module 7: E-Banking (FR-EBK-001)

### 7.1 Parser Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| EBK-PAR-001 | Design CSV parser interface | 2h | TODO |
| EBK-PAR-002 | Implement VCB parser | 4h | TODO |
| EBK-PAR-003 | Implement BIDV parser | 3h | TODO |
| EBK-PAR-004 | Implement CTG parser | 3h | TODO |
| EBK-PAR-005 | Implement VTB parser | 3h | TODO |
| EBK-PAR-006 | Implement ACB parser | 3h | TODO |

### 7.2 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| EBK-SVC-001 | Create `bank_import_service.go` | 4h | TODO |
| EBK-SVC-002 | Implement reconciliation logic | 6h | TODO |
| EBK-SVC-003 | Write tests | 4h | TODO |

### 7.3 Handler Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| EBK-HDL-001 | Create handler | 3h | TODO |
| EBK-HDL-002 | Write tests | 3h | TODO |

**Module Total:** 38h (~4.8 days)

---

## Module 8: E-Tax Filing (FR-ETX-001)

### 8.1 XML Generation
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| ETX-XML-001 | VAT declaration XML (Mẫu 01) | 6h | TODO |
| ETX-XML-002 | CIT declaration XML (Mẫu 01) | 6h | TODO |
| ETX-XML-003 | PIT declaration XML | 4h | TODO |

### 8.2 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| ETX-SVC-001 | Create `tax_filing_service.go` | 4h | TODO |
| ETX-SVC-002 | Integrate with GDT client | 4h | TODO |
| ETX-SVC-003 | Write tests | 4h | TODO |

### 8.3 Handler Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| ETX-HDL-001 | Create handler | 3h | TODO |
| ETX-HDL-002 | Write tests | 3h | TODO |

**Module Total:** 34h (~4.3 days)

---

## Module 9: Digital Signature (FR-DSG-001)

### 9.1 Provider Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| DSG-PRV-001 | Design `SignatureProvider` interface | 2h | TODO |
| DSG-PRV-002 | Implement VNPT-CA provider | 6h | TODO |
| DSG-PRV-003 | Implement mock provider for testing | 2h | TODO |

### 9.2 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| DSG-SVC-001 | Create `signature_service.go` | 3h | TODO |
| DSG-SVC-002 | Write tests | 2h | TODO |

**Module Total:** 15h (~1.9 days)

---

## Module 10: Backup & Restore (FR-BKP-001)

### 10.1 Database Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| BKP-DB-001 | Create migration `000039_backups.up.sql` | 2h | TODO |

### 10.2 Service Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| BKP-SVC-001 | Create `backup_service.go` | 4h | TODO |
| BKP-SVC-002 | Implement pg_dump wrapper | 3h | TODO |
| BKP-SVC-003 | Implement restore logic | 4h | TODO |
| BKP-SVC-004 | Write tests | 3h | TODO |

### 10.3 Handler Layer
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| BKP-HDL-001 | Create handler | 3h | TODO |
| BKP-HDL-002 | Write tests | 3h | TODO |

**Module Total:** 22h (~2.8 days)

---

## Module 11: Report Customization (FR-RPT-001)

### 11.1 Integration
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| RPT-INT-001 | Add report options to SystemOption | 2h | TODO |
| RPT-INT-002 | Extend PDF generator | 4h | TODO |
| RPT-INT-003 | Add logo upload | 3h | TODO |
| RPT-INT-004 | Write tests | 2h | TODO |

**Module Total:** 11h (~1.4 days)

---

## Module 12: Multi-branch Enhancements

### 12.1 System Options
| Task ID | Task | Est. Hours | Status |
|---------|------|-----------|--------|
| MB-SYS-001 | Branch-level options support | 4h | TODO |
| MB-SYS-002 | Per-branch numbering rules | 3h | TODO |
| MB-SYS-003 | Branch report filtering | 4h | TODO |
| MB-SYS-004 | Write tests | 3h | TODO |

**Module Total:** 14h (~1.8 days)

---

## Summary by Module

| Module | Tasks | Hours | Days |
|--------|-------|-------|------|
| System Options | 19 | 33.5 | 4.2 |
| Voucher Numbering | 11 | 29.5 | 3.7 |
| Fiscal Year | 4 | 13 | 1.6 |
| Number Format | 4 | 9 | 1.1 |
| Contracts | 12 | 34.5 | 4.3 |
| Loan Agreements | 12 | 39.5 | 4.9 |
| E-Banking | 11 | 38 | 4.8 |
| E-Tax Filing | 8 | 34 | 4.3 |
| Digital Signature | 5 | 15 | 1.9 |
| Backup & Restore | 6 | 22 | 2.8 |
| Report Customization | 4 | 11 | 1.4 |
| Multi-branch | 4 | 14 | 1.8 |
| **TOTAL** | **100** | **293** | **36.6** |

---

## Notes

1. **Hours include testing** — approximately 30% of time is tests
2. **Days assume 8h/day** — adjust for part-time availability
3. **Parallel execution possible** — some modules can overlap
4. **Buffer recommended** — add 20% for unknowns = ~44 days total
