# Test Strategy

**Author:** QA Lead (20+ yrs)  
**Sources:** dev.to, Stack Overflow, Medium, Go docs, Go wiki, Keploy, Testcontainers, QASkills.sh, Kodus, industry research 2025-2026  
**Applies to:** GoTax (Go 1.26.5, Gin, pgx/v5)

---

## 1. Test Pyramid

```
        /\          E2E / UI: 0 (API backend, no UI)
       /  \         Integration: target 15%
      /    \
     /      \       Service (real repos): target 35%
    /________\
   /__________\     Unit (domain + logic): target 50%
```

Current GoTax: 112 tests across 3 layers — matches pyramid.

| Layer | Location | Count | Runs without DB |
|-------|----------|-------|-----------------|
| Unit | `domain/models_test.go` | 22 | Yes |
| Service | `service/service_test.go` | 34 | Yes (in-memory repos) |
| Handler | `handler/handler_test.go` + `company_handler_test.go` | 56 | Yes (in-memory repos + httptest) |

**Rule:** Every layer uses in-memory repos. No DB, no network, no external dependencies. This keeps tests fast, deterministic, parallelizable.

---

## 2. Go Testing Idioms

### 2.1 Table-Driven Tests

Standard for all multi-case tests. Structure:

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid code", "1111", false},
        {"empty code", "", true},
        {"too long", "12345678901", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
            }
        })
    }
}
```

- Use `t.Run` for named subtests (isolated failures, filterable with `-run`)
- `tt := tt` NOT needed (fixed in Go 1.22)
- Use testify `assert`/`require` for cleaner assertions

### 2.2 Subtest Selection

```sh
go test -v -run 'TestValidate/valid_code'    # single subtest
go test -v -run 'TestValidate/'               # all subtests
```

### 2.3 Parallel Execution

```go
func TestSomething(t *testing.T) {
    t.Parallel()
    // ...
}
```

Run with `go test -parallel N` or default `GOMAXPROCS`.

### 2.4 Cleanup

```go
func TestWithCleanup(t *testing.T) {
    db := setupDB()
    t.Cleanup(func() { db.Close() })
}
```

### 2.5 TestMain for Package-Level Setup

```go
func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    teardown()
    os.Exit(code)
}
```

Not currently used in GoTax (not needed yet — in-memory repos have no setup cost).

---

## 3. Layer-Specific Patterns

### 3.1 Domain Tests (`domain/models_test.go`)

- Validate struct validation methods
- Pure unit tests, zero dependencies
- Table-driven for multiple error/valid cases
- Test exported Validate() on every domain entity

### 3.2 Service Tests (`service/service_test.go`)

- Real in-memory repos wired in `setupService()`
- Full journal lifecycle: draft → submit → approve → post → cancel
- Auth flows: login, 2FA, refresh, lockout, change password
- Reports: trial balance, account balance, drill-down
- COA operations: versioning, mappings, IFRS, analysis
- Edge cases: duplicate accounts, invalid status transitions, missing entities

### 3.3 Handler Tests (`handler/handler_test.go` + `company_handler_test.go`)

- Real in-memory repos + real service (NO mock services)
- Mock auth middleware sets `user_id`, `username`, `role` in gin.Context
- `httptest.NewRecorder` + `httptest.NewRequest`
- Test: happy path, validation errors, not-found, auth-required, missing params
- Company domain tested separately via `RegisterCompanyRoutes`

---

## 4. Coverage Targets

| Metric | Current | Target |
|--------|---------|--------|
| Package coverage | — | ≥ 80% |
| Critical paths (journal lifecycle, auth) | — | ≥ 90% |
| Handlers | tested via service layer | ≥ 70% |

Commands:

```sh
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 5. Race Detection

Run on every test suite:

```sh
go test -race ./...
```

Rule: zero data races in CI. A race is a bug, not a flake. Fix the shared-state access.

---

## 6. Fuzzing

Go 1.18+ native fuzzing for parsers, validators, and any function accepting untrusted input.

```go
func FuzzParseTaxCode(f *testing.F) {
    f.Add("1234567890")
    f.Add("")
    f.Add("abc")
    f.Fuzz(func(t *testing.T, s string) {
        // must never panic
        _ = domain.ValidateTaxCode(s)
    })
}
```

Run:

```sh
go test -fuzz=FuzzParseTaxCode -fuzztime=30s ./internal/domain/
```

---

## 7. Benchmark Tests

For hot paths:

```go
func BenchmarkValidate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Validate("1111")
    }
}
```

Run:

```sh
go test -bench=. -benchmem ./...
```

---

## 8. What NOT to Do (Outdated/Weak Practices)

| Practice | Why Removed |
|----------|-------------|
| `gomock` codegen | Deprecated. Interface-based mocks in same package are simpler |
| Ginkgo/Gomega BDD | Unnecessary indirection. Go `testing` + testify is sufficient |
| Mock services with function fields | Current GoTax uses real in-memory repos — better coverage, less maintenance |
| `tt := tt` inside loop | Fixed in Go 1.22. Not needed |
| `ioutil` package | Removed in Go 1.19. Use `io`/`os` |
| Selenium/E2E UI tests | API backend has no UI. Contract tests cover integration |
| Shared mutable test data | Each test creates own fixtures. No test pollution |
| `context.TODO()` in production tests | Use `context.Background()` |
| External test runners | `go test` is sufficient. No mocha/jest-style runners |

---

## 9. Future: Integration Tests with Testcontainers

When GoTax needs PG-specific tests (e.g., migration correctness, PG-specific features):

```go
func TestMain(m *testing.M) {
    pg, _ := testcontainers.GenericContainer(context.Background(),
        testcontainers.GenericContainerRequest{
            ContainerRequest: testcontainers.ContainerRequest{
                Image: "postgres:16-alpine",
            },
            Started: true,
        })
    os.Setenv("DATABASE_URL", pg.ConnectionString())
    code := m.Run()
    pg.Terminate()
    os.Exit(code)
}
```

Not implemented yet. In-memory repos cover all current needs.

---

## 10. CI Gates

| Gate | Command | Blocking |
|------|---------|----------|
| Unit tests | `go test ./...` | Yes |
| Race detection | `go test -race ./...` | Yes |
| Coverage threshold | `go test -cover ./...` | Warning (≥80% target) |
| Fuzzing | `go test -fuzztime=30s ./...` | Nightly |
| Build | `go build ./...` | Yes |
| Vet | `go vet ./...` | Yes |

---

## Related

- **`AGENTS.md`** (repo root) — commands, test-running patterns, adding-feature workflow
- **`docs/standards/CODE_QUALITY.md`** — static analysis toolchain, review checklist, PRR gate, quality metrics

---

## References

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests)
- [Go Fuzzing](https://go.dev/doc/fuzz/)
- [Go httptest Package](https://pkg.go.dev/net/http/httptest)
- [Testcontainers for Go](https://github.com/testcontainers/testcontainers-go)
- [Go Code Coverage](https://go.dev/blog/cover)
- [Kodus Go Code Review Best Practices 2026](https://kodus.io/en/golang-code-review-practices/)
- [QASkills.sh Go Testing Guide 2026](https://qaskills.sh/blog/go-testing-tutorial-table-driven-tests-2026)
- [Go Unit Testing Structure & Best Practices (Rost Glukhov 2025)](https://dev.to/rosgluk/go-unit-testing-structure-best-practices-1b5n)
