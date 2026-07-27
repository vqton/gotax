# GoTax COA Module - Technical Specifications

## Version: 1.0 | Stack: Go + Gin + PostgreSQL 16+ | Circular 99/2025/TT-BTC

---

## 1. Architecture

```
Client Layer (Web/Mobile/API)
        │ HTTPS/TLS 1.3
        ▼
 ┌─────────────────────────────────────────────────────┐
 │                API Gateway / Auth Layer              │
 │  JWT validation → RBAC check → Rate limit           │
 └─────────────────────┬───────────────────────────────┘
                       ▼
 ┌─────────────────────────────────────────────────────┐
 │              COA Service Layer                       │
 │                                                      │
 │  ┌─────────────┐  ┌──────────────┐  ┌───────────┐  │
 │  │ COA Handler  │  │ COA Service  │  │ Validator  │  │
 │  │ (HTTP/REST)  │──► (Business)   │──► (Rules)    │  │
 │  └─────────────┘  └──────┬───────┘  └───────────┘  │
 │                          │                          │
 │  ┌───────────────────────┴───────────────────────┐  │
 │  │         Repository Layer                       │  │
 │  │  ┌─────────┐ ┌──────────┐ ┌────────────────┐  │  │
 │  │  │ Account │ │ Version  │ │ Approval       │  │  │
 │  │  │ Repo    │ │ Repo     │ │ Repo           │  │  │
 │  │  └─────────┘ └──────────┘ └────────────────┘  │  │
 │  │  ┌─────────┐ ┌──────────┐ ┌────────────────┐  │  │
 │  │  │ Mapping │ │ IFRS     │ │ Analysis       │  │  │
 │  │  │ Repo    │ │ Repo     │ │ Repo           │  │  │
 │  │  └─────────┘ └──────────┘ └────────────────┘  │  │
 │  └───────────────────────────────────────────────┘  │
 └──────────────────────────┬──────────────────────────┘
                            ▼
 ┌─────────────────────────────────────────────────────┐
 │              Data Layer                              │
 │  PostgreSQL 16+ (Primary)                           │
 │  Redis (Cache for COA tree)                         │
 │  S3/MinIO (Export files, IGAP docs)                 │
 └─────────────────────────────────────────────────────┘
```

---

## 2. Database Schema (New Tables)

### account_versions
```sql
CREATE TABLE account_versions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version_number  VARCHAR(20) NOT NULL,
    snapshot        JSONB NOT NULL,           -- full COA snapshot
    change_summary  TEXT,                     -- human-readable change log
    change_reason   TEXT,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (version_number)
);
```

### approval_requests
```sql
CREATE TABLE approval_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID REFERENCES tenants(id),
    entity_type     VARCHAR(30) NOT NULL,     -- 'ACCOUNT'
    entity_id       VARCHAR(50) NOT NULL,
    request_type    VARCHAR(20) NOT NULL,     -- 'CREATE','MODIFY','DEACTIVATE'
    old_value       JSONB,
    new_value       JSONB,
    reason          TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING'
                    CHECK (status IN ('PENDING','APPROVED','REJECTED','CANCELLED','EXPIRED')),
    requested_by    UUID REFERENCES users(id),
    reviewed_by     UUID REFERENCES users(id),
    review_note     TEXT,
    expires_at      TIMESTAMPTZ,              -- 48h auto-expiry
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at     TIMESTAMPTZ
);
```

### account_mappings
```sql
CREATE TABLE account_mappings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_regime   VARCHAR(20) NOT NULL,     -- 'TT200','TT133','USER_DEFINED'
    target_regime   VARCHAR(20) NOT NULL,     -- 'TT99','TT133'
    old_code        VARCHAR(20) NOT NULL,
    new_code        VARCHAR(20) NOT NULL REFERENCES accounts(code),
    mapping_type    VARCHAR(20) NOT NULL      -- 'DIRECT','MERGE','SPLIT','CUSTOM'
                    CHECK (mapping_type IN ('DIRECT','MERGE','SPLIT','CUSTOM')),
    split_ratio     NUMERIC(5,4),             -- for SPLIT type
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (source_regime, target_regime, old_code)
);
```

### account_ifrs_mapping
```sql
CREATE TABLE account_ifrs_mapping (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vas_code        VARCHAR(20) NOT NULL REFERENCES accounts(code),
    ifrs_code       VARCHAR(20) NOT NULL,
    ifrs_name       VARCHAR(255),
    reclassification_rule TEXT,                -- SQL expression or description
    adjustment_type VARCHAR(20),               -- 'RECLASSIFY','ADJUST','BOTH'
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### account_analysis
```sql
CREATE TABLE account_analysis (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_code    VARCHAR(20) NOT NULL REFERENCES accounts(code),
    cost_center_id  VARCHAR(50),
    profit_center_id VARCHAR(50),
    department_id   VARCHAR(50),
    project_id      VARCHAR(50),
    custom_dim1     VARCHAR(50),
    custom_dim2     VARCHAR(50),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_code)
);
```

### Schema Migration: accounts table (add new columns)
```sql
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE'
    CHECK (status IN ('ACTIVE', 'FROZEN', 'INACTIVE'));
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS freeze_reason TEXT;
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS version_number VARCHAR(20);
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS effective_from DATE;
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS effective_to DATE;
```

---

## 3. Account Code Validation Rules (Go Implementation)

```go
func ValidateAccountCode(code, parentCode string, existingCodes map[string]bool) error {
    // R001: numeric only
    for _, c := range code {
        if !unicode.IsDigit(c) {
            return ErrAccountCodeInvalid
        }
    }
    // R002: min 3, max 20
    if len(code) < 3 || len(code) > 20 {
        return ErrAccountCodeInvalid
    }
    // R003: first digit valid
    first := code[0]
    switch first {
    case '0': // Loai 0: off-balance sheet
    case '1','2': // Loai 1,2: ASSET
    case '3': // Loai 3: LIABILITY
    case '4': // Loai 4: EQUITY
    case '5','7': // Loai 5,7: REVENUE
    case '6','8','9': // Loai 6,8,9: EXPENSE
    default:
        return ErrAccountCodeInvalid
    }
    // R004: parent prefix
    if parentCode != "" && !strings.HasPrefix(code, parentCode) {
        return ErrAccountCodeInvalidHierarchy
    }
    // R005: uniqueness
    if existingCodes[code] {
        return ErrAccountCodeExists
    }
    return nil
}

func ValidateAccountType(type AccountType) error {
    switch type {
    case AccountTypeAsset, AccountTypeLiability, AccountTypeEquity,
         AccountTypeRevenue, AccountTypeExpense:
        return nil
    default:
        return ErrAccountInvalidType
    }
}
```

---

## 4. COA Import Engine Design

```go
type ImportRow struct {
    Code       string  `csv:"code"`
    Name       string  `csv:"name"`
    Type       string  `csv:"type"`
    ParentCode string  `csv:"parent_code"`
    IsParent   bool    `csv:"is_parent"`
    IsActive   bool    `csv:"is_active"`
    IsForeign  bool    `csv:"is_foreign"`
    DetailBy   string  `csv:"detail_by"`
}

type ImportResult struct {
    Created    int      `json:"created"`
    Updated    int      `json:"updated"`
    Errors     []ImportError `json:"errors,omitempty"`
    BatchID    string   `json:"batch_id"`
}

// ImportFlow:
// 1. Parse → Validate headers → Validate rows (parallel)
// 2. Build dependency graph (parents before children)
// 3. Sort by hierarchy depth (BFS)
// 4. BEGIN TX → execute UPSERT per account → COMMIT
// 5. On failure: ROLLBACK, return error
// 6. Log audit for each row
type ImportEngine struct {
    parser     FileParser       // CSV or Excel
    validator  AccountValidator
    repo       AccountRepository
    auditRepo  AuditLogRepository
    maxRows    int
    maxSize    int64
}
```

### Import Performance Targets
| Operation | Target |
|---|---|
| Parse 1000 rows | < 1s |
| Validate 1000 rows | < 2s |
| Import 1000 rows (transactional) | < 3s |
| Full import 10000 rows | < 30s |

---

## 5. COA Versioning Strategy

```go
type COAVersionService struct {
    repo    AccountRepository
    verRepo VersionRepository
    snapFn  func() ([]Account, error)  // snapshot current COA
}

// On every COA change:
// 1. Apply change to accounts table
// 2. Snapshot complete COA (SELECT *)
// 3. Diff vs previous snapshot
// 4. Generate version metadata:
//    - Major bump: structural change (add/remove accounts)
//    - Minor bump: attribute change (rename, detail_by)
//    - Patch: no version bump (is_active only)
// 5. Store snapshot as JSONB in account_versions
// 6. Update accounts.version_number

// Version ID format: "COA-V{major}.{minor}.{yyyymmdd}"
// Example: "COA-V2.1.20260727"

// As-of-date query:
// SELECT * FROM account_versions
// WHERE effective_at <= $1
// ORDER BY effective_at DESC
// LIMIT 1
```

---

## 6. Approval Workflow State Machine

```
                     ┌──────────┐
                     │ PENDING  │
                     └────┬─────┘
                          │
            ┌─────────────┼─────────────┐
            │             │             │
            ▼             ▼             ▼
      ┌─────────┐  ┌─────────┐  ┌─────────┐
      │APPROVED │  │REJECTED │  │CANCELLED│
      └────┬────┘  └─────────┘  └─────────┘
           │
           ▼
    Apply change to
    accounts table
           │
           ▼
    Log audit entry
```

---

## 7. Export Engine

```go
type ExportFormat string
const (
    ExportCSV   ExportFormat = "csv"
    ExportExcel ExportFormat = "xlsx"
    ExportPDF   ExportFormat = "pdf"   // IGAP document
)

type ExportEngine struct {
    accountRepo AccountRepository
    templateRenderer *TemplateRenderer  // for PDF/IGAP
}

// CSV: streaming writer, 5000 rows per chunk
// Excel: excelize library, multi-sheet
// PDF: go-wkhtmltopdf or similar, HTML template → PDF

// IGAP PDF template includes:
// - Company logo and header
// - "INTERNAL ACCOUNTING POLICY" title
// - Version number and effective date
// - COA table grouped by Loai (type)
// - Account descriptions
// - Signatory block
```

---

## 8. Caching Strategy

| Data | Cache | TTL | Invalidation |
|---|---|---|---|
| Full COA tree | Redis | 1 hour | On COA change |
| Account by code | Redis | 30 min | On account update |
| COA version list | Redis | 1 hour | On version create |
| Account hierarchy | Redis | 1 hour | On add/remove |
| Import preview | In-memory | Session | On import confirm/discard |

---

## 9. Error Codes

| Code | HTTP | Description |
|---|---|---|
| COA_001 | 400 | Invalid account code format |
| COA_002 | 400 | Invalid account type |
| COA_003 | 400 | Parent account not found |
| COA_004 | 400 | Parent is not a group account |
| COA_005 | 400 | Invalid hierarchy (code prefix mismatch) |
| COA_006 | 409 | Account code already exists |
| COA_007 | 400 | Cannot modify account type (has balance) |
| COA_008 | 409 | Account has children, cannot delete |
| COA_009 | 409 | Account has balance, cannot delete |
| COA_010 | 409 | Account is frozen |
| COA_011 | 400 | Account not found |
| COA_012 | 400 | Invalid import file format |
| COA_013 | 413 | Import file too large |
| COA_014 | 400 | Import validation failed |
| COA_015 | 500 | Import transaction failed |
| COA_016 | 400 | Approval request not found |
| COA_017 | 409 | Approval request already processed |
| COA_018 | 400 | Cannot self-approve |
| COA_019 | 400 | Version not found |
| COA_020 | 400 | Account status transition invalid |
| COA_021 | 409 | Account already frozen |
| COA_022 | 409 | Account not frozen |
| COA_023 | 400 | IGAP template not configured |

---

## 10. Dependencies

### Existing (in go.mod)
- github.com/gin-gonic/gin v1.12.0
- github.com/jackc/pgx/v5 v5.10.0
- github.com/golang-jwt/jwt/v5 v5.3.1

### New Dependencies Required
```go
import (
    "github.com/xuri/excelize/v2"       // Excel read/write
    "github.com/jszwec/csvutil"         // CSV marshal/unmarshal
    "github.com/SebastiaanKlippert/go-wkhtmltopdf"  // PDF generation
    "github.com/redis/go-redis/v9"      // COA caching
    "github.com/google/uuid"            // UUID gen
    "github.com/goccy/go-json"          // JSONB handling (already indirect)
)
```

---

## 11. Testing Strategy

| Test Type | Scope | Framework |
|---|---|---|
| Unit (Account validation) | All account rules R001-R043 | go test + testify |
| Unit (Import parser) | CSV/Excel parse + validate | go test + testify |
| Unit (Version service) | Snapshot + diff + compare | go test + testify |
| Unit (Approval workflow) | State machine transitions | go test + testify |
| Integration (PG repos) | CRUD + balance calc + import | testcontainers-postgres |
| Integration (Import E2E) | Full import → verify DB | testcontainers-postgres |
| E2E (API) | Full REST API coverage | httptest + testify |
| Performance | Import throughput, query latency | benchstat |
