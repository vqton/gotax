# Payroll Module — Implementation Roadmap

**Document ID:** PLAN-PAYROLL-001
**Version:** 1.0
**Date:** 2026-08-04
**Duration:** 16 weeks (4 months)
**Target:** PROD-ready MVP

---

## Executive Summary

Build Vietnamese-compliant payroll module for GoTax from 0% to PROD-ready in 16 weeks. Four phases, each delivering working vertical slices. Dependencies flow bottom-up: schema → models → repos → service → handler → tests → docs.

---

## Architecture Decisions

| Decision | Rationale |
|----------|-----------|
| Separate `internal/payroll/` package | Clean module boundary, follows existing pattern (tax/, company/) |
| Extend Employee model via `employee_payroll_info` table | Avoids breaking existing Employee struct; 1:1 relationship |
| Config-driven rates | Vietnamese laws change frequently; JSON config per company |
| Separate timekeeping module | Timekeeping reusable beyond payroll (attendance, reporting) |
| PDF payslip via maroto/v2 | Already in codebase, no new dependency |
| GL posting via existing service | Reuse `Service.PostJournalEntry()` from GL module |

---

## Dependency Graph

```
Phase 1: Foundation
    │
    ├── Migration schema (DB tables)
    │       │
    │       ├── Domain models (structs)
    │       │       │
    │       │       ├── Repository interfaces
    │       │       │       │
    │       │       │       ├── PG repository
    │       │       │       └── Memory repository
    │       │       │
    │       │       └── Validation functions
    │       │
    │       └── Seed data (rate config, holidays)
    │
    └── Calculation engine (pure functions)
            │
            └── Unit tests (100% coverage)

Phase 2: Timekeeping
    │
    ├── Timekeeping schema + models
    ├── CSV import handler
    ├── Leave management
    └── Integration with payroll calc

Phase 3: Payroll Processing
    │
    ├── Period management
    ├── Payroll run engine
    ├── Payslip generation
    └── GL integration

Phase 4: Declarations & Polish
    │
    ├── D02-TS generation
    ├── 05/KK-TNCN generation
    ├── Approval workflow
    └── Employee self-service
```

---

## Phase Overview

| Phase | Weeks | Focus | Deliverable |
|-------|-------|-------|-------------|
| 1 | 1-4 | Foundation & Calculation Engine | Gross-to-net calc with SI/HI/UI/PIT |
| 2 | 5-8 | Timekeeping & Attendance | Import, OT, night shift, leave |
| 3 | 9-12 | Payroll Processing & GL | Periods, runs, payslips, journal entries |
| 4 | 13-16 | Declarations & Polish | D02-TS, 05/KK-TNCN, approval, self-service |

---

## Risk Register

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Law changes during development | HIGH | MEDIUM | Config-driven rates, update JSON not code |
| Complex overtime rules | MEDIUM | HIGH | Rule engine with exhaustive test cases |
| Multi-region min wage | MEDIUM | LOW | Region-based config, already modeled |
| Foreign employee edge cases | MEDIUM | MEDIUM | Flag-based branching, separate test suite |
| I-VAN integration complexity | LOW | HIGH | Defer to Phase 4, XML generation only |
| GL account code conflicts | LOW | LOW | Use standard Vietnamese COA codes |
| Performance with 500+ employees | LOW | LOW | Batch processing, indexed queries |

---

## Success Criteria

| Metric | Target | Measurement |
|--------|--------|-------------|
| Calculation accuracy | 100% | Match HAPRI calculator for 50 test cases |
| Test coverage | >90% | `go test -cover` |
| Processing speed | <30s for 100 employees | Benchmark test |
| GL posting accuracy | 100% | Reconciliation test |
| Compliance | 100% | Audit against Law 109/2025, Law 41/2024 |

---

## Checkpoint Schedule

| Checkpoint | After Tasks | Gate Criteria |
|------------|-------------|---------------|
| CP-1 | T1.1-T1.5 | Schema + models compile, all tests pass |
| CP-2 | T1.6-T1.10 | Calc engine produces correct gross-to-net |
| CP-3 | T2.1-T2.5 | Timekeeping import + OT calc works |
| CP-4 | T3.1-T3.5 | Full payroll run produces payslips |
| CP-5 | T4.1-T4.5 | D02-TS + 05/KK-TNCN generate correctly |
| CP-6 | T4.6-T4.10 | End-to-end flow: import → calc → approve → GL → payslip |

---

## Resource Allocation

| Role | Allocation | Responsibilities |
|------|------------|------------------|
| Backend Dev 1 | 100% | Core payroll engine, GL integration |
| Backend Dev 2 | 100% | Timekeeping, declarations |
| QA | 50% | Test cases, compliance verification |
| BA/Accountant | 25% | Requirements clarification, UAT |

---

## Sprint Cadence

| Sprint | Weeks | Focus |
|--------|-------|-------|
| Sprint 1 | 1-2 | Schema + Domain models + Repos |
| Sprint 2 | 3-4 | Calculation engine + Unit tests |
| Sprint 3 | 5-6 | Timekeeping import + OT calc |
| Sprint 4 | 7-8 | Leave management + Integration |
| Sprint 5 | 9-10 | Period mgmt + Payroll run |
| Sprint 6 | 11-12 | Payslip + GL posting |
| Sprint 7 | 13-14 | D02-TS + 05/KK-TNCN |
| Sprint 8 | 15-16 | Approval + Self-service + Polish |

---

## Appendix: Detailed Task Reference

See `tasks/todo.md` for full task list with acceptance criteria.
See `docs/payroll/` for complete specifications.
