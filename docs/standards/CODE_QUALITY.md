# Code Quality Standard

**Author:** System Analyst Lead (20+ yrs)  
**Sources:** dev.to, Stack Overflow, Medium, JetBrains Go Modern Guidelines, CISQ ISO 5055, industry code review research 2025-2026  
**Applicable to:** GoTax Go codebase (Go 1.26.5)

---

## 1. Static Analysis Toolchain

Run before every commit and in CI:

| Tool | Scope | Command |
|------|-------|---------|
| `go vet` | Compiler-level bugs, nil deref, goroutine leaks, escape analysis | `go vet ./...` |
| `go test -race` | Data race detection | `go test -race ./...` |
| `staticcheck` | Unused code, dead branches, style violations | `staticcheck ./...` |
| `gosec` | Security vulns (injection, hardcoded creds, weak crypto) | `gosec ./...` |
| `govulncheck` | Known vulns in dependencies | `govulncheck ./...` |
| `ineffassign` | Unused variable assignments | `ineffassign ./...` |
| `golangci-lint` | Meta-linter, runs all above in one pass | `golangci-lint run` |

**Rule:** Zero warnings in CI. No PR merged with `vet` or `staticcheck` warnings.

---

## 2. Code Review Checklist (8 Dimensions)

Review every PR against these dimensions, priority-ordered:

### 2.1 Correctness & Logic

- [ ] Edge cases handled: empty input, nil, overflow, zero values
- [ ] No silent ignored errors (`_` for error returns only when intentional and documented)
- [ ] Boundary conditions tested (off-by-one, empty slice, max int)
- [ ] Data races impossible (no unprotected shared mutable state)
- [ ] Context properly propagated through entire call chain

### 2.2 Security

- [ ] All external inputs validated and sanitized at handler boundary
- [ ] SQL queries use parameterized statements (never string interpolation)
- [ ] No secrets (passwords, tokens, keys) in logs, error messages, or URLs
- [ ] Token/signature comparisons use constant-time (`crypto/subtle`)
- [ ] Rate limiting on auth and mutation endpoints
- [ ] Dependencies scanned for known vulns (`govulncheck`)

### 2.3 Error Handling

- [ ] Errors wrapped with context: `fmt.Errorf("doing X: %w", err)`
- [ ] Callers use `errors.Is` / `errors.As`, never string comparison on `.Error()`
- [ ] No panic outside package init. Recover only at top-level HTTP handler
- [ ] Sentry-level errors (DB, network) return 500, not leak stack to client
- [ ] Validation errors return 400 with actionable message

### 2.4 Concurrency & Performance

- [ ] Goroutine lifecycle managed (no leaked goroutines)
- [ ] `sync.WaitGroup` or errgroup for fan-out; never raw `go` without tracking
- [ ] Mutex covers minimal critical section, not whole function
- [ ] Hot paths: no allocation in tight loops, pre-slice capacity
- [ ] `-race` clean on all tests

### 2.5 Observability

- [ ] Structured logging (key-value pairs, not `fmt.Sprintf`)
- [ ] Every external call (DB, API) has context timeout
- [ ] Metrics for: request count, latency p50/p99, error rate
- [ ] Trace context propagated across service boundaries

### 2.6 Maintainability

- [ ] Functions small, single-purpose, < 40 lines
- [ ] No boolean trap parameters — use enum or option struct
- [ ] Package names short, lowercase, no underscores
- [ ] Interfaces defined by consumer, not producer. 1-3 methods max
- [ ] No circular imports. Dependency direction: handler → service → repository

### 2.7 Testing

- [ ] Unit tests for all exported functions
- [ ] Table-driven tests for multiple cases
- [ ] Edge cases covered (not just happy path)
- [ ] Service tests use in-memory repos (no DB dependency)
- [ ] Handler tests use `httptest.NewRecorder` + mock auth middleware
- [ ] Race detection enabled (`go test -race`)

### 2.8 Go Idiom Compliance (1.26)

- [ ] `range-over-func` used where iterable pattern emerges (not overused)
- [ ] `slices`, `maps` stdlib used (no hand-rolled sort/search)
- [ ] `cmp.Or` for zero-value fallback
- [ ] `errors.Join` for multi-error aggregation
- [ ] `log/slog` for structured logging (not `log`)
- [ ] `encoding/json` with struct tags, not reflection hacks

---

## 3. Architecture Standards

### 3.1 Layer Rules

| Layer | Responsibility | Depends on | Forbidden |
|-------|---------------|-----------|-----------|
| Handler | Parse request, call service, return JSON | Service | Business logic, DB access |
| Service | Business rules, validation, orchestration | Domain interfaces | HTTP, raw DB |
| Repository | Data access (2 impls: PG + memory) | Domain models | Business logic |
| Domain | Models, interfaces, errors | Nothing | Anything |

### 3.2 Package Boundaries

- No circular imports (enforced by `go vet`)
- Domain package has zero imports from project
- Interface defined where consumed (service), not where implemented (repo)
- New entity = interface in domain + 2 repo impls + service method + handler + test

### 3.3 Configuration

- Env vars at startup. No config files shipped with binary
- `JWT_SECRET` required. Server panics if unset
- `DATABASE_URL` optional — set = PG backend, unset = in-memory

---

## 4. Production Readiness Review (PRR)

Before any deploy to production:

- [ ] Static analysis suite passes (vet + race + staticcheck + gosec)
- [ ] All tests pass (unit + integration)
- [ ] No P0/P1 findings in security scan
- [ ] Metrics and logging confirmed emitting
- [ ] Runbook exists for this service/feature
- [ ] Rollback plan documented
- [ ] Load test completed (if performance-sensitive path changed)
- [ ] Migration SQL reviewed (backward-compatible, no data loss)

---

## 5. Measuring Code Quality

| Metric | Target | Tool |
|--------|--------|------|
| Cyclomatic complexity per func | ≤ 10 | `gocyclo` |
| Test coverage | ≥ 80% | `go test -cover` |
| Lines per function | ≤ 40 | `gocyclo`, `cloc` |
| Duplication | ≤ 5% | `staticcheck` |
| 'go vet' warnings | 0 | `go vet` |
| Known vuln dependencies | 0 | `govulncheck` |

---

## 6. Weak/Outdated Practices — Removed

These do not appear in this standard because research and Go 1.26 evolution made them obsolete:

- ~~Error strings compared with `==`~~ → use `errors.Is` / `errors.As`
- ~~Raw `go` without errgroup or WaitGroup~~ → always track goroutine lifecycle
- ~~`log.Fatal` in non-main packages~~ → return error, let main decide
- ~~`ioutil` package~~ → moved to `io` / `os` (gone since Go 1.19, still seen in old guides)
- ~~Hand-written sort/search~~ → use `slices` / `sort` stdlib
- ~~`context.TODO` in production code~~ → use `context.Background` or pass real context
- ~~Mock services in handler tests~~ → real in-memory repos are cheaper and test more
- ~~Single flat package~~ → Clean Architecture layers with explicit boundaries

---

## Related

- **`AGENTS.md`** (repo root) — operational commands, build/test/run steps, testing conventions, adding-feature workflow, architecture gotchas. Code quality standard covers _what_ to enforce; AGENTS.md covers _how_ to run this project.

## References

- [JetBrains Go Modern Guidelines](https://github.com/JetBrains/go-modern-guidelines) (2026)
- [CISQ ISO 5055 — Software Quality Measures](https://www.it-cisq.org/standards/code-quality-standards/)
- [Go Static Analysis Tools — 20 Tools Survey](https://www.in-com.com/blog/write-better-go-code-20-static-analysis-tools-that-catch-bugs-before-you-do/)
- [Kodus Go Code Review Best Practices (2026)](https://kodus.io/en/golang-code-review-practices/)
- [Go Toolchain Upgrades: Vet Improvements](https://medium.com/@backendbyeli/go-toolchain-upgrades-static-analysis-vet-enhancements-that-ship-safer-code-7512836db7d0)
