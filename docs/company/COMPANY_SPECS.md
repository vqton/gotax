# GoTax Company Module — Technical Specifications

## Version: 1.0 | Date: 2026-07-27 | Status: DRAFT

---

## 1. Data Model

### 1.1 companies
```sql
CREATE TABLE companies (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,

    -- Legal Identity
    legal_name_vn       VARCHAR(500) NOT NULL,
    legal_name_en       VARCHAR(500),
    short_name          VARCHAR(200),
    legal_form          VARCHAR(50) NOT NULL,  -- LLC_1MEMBER, LLC_2MEMBERS, JSC, SOLE_PROP, PARTNERSHIP, FO, RO
    tax_code            VARCHAR(20) NOT NULL,  -- MST: 10 or 13 digits
    business_reg_no     VARCHAR(50),           -- MSDN / Giay phep DKDN
    business_reg_date   DATE,
    business_reg_place  VARCHAR(200),

    -- Address
    reg_address         TEXT NOT NULL,
    reg_address_province VARCHAR(100),
    reg_address_district VARCHAR(100),
    head_office_address TEXT,
    head_office_province VARCHAR(100),
    head_office_district VARCHAR(100),

    -- Contact
    phone               VARCHAR(30),
    email               VARCHAR(200),
    website             VARCHAR(200),

    -- Representatives
    legal_rep_name      VARCHAR(200),
    legal_rep_title     VARCHAR(200),
    legal_rep_id_number VARCHAR(50),
    legal_rep_id_date   DATE,
    legal_rep_id_place  VARCHAR(200),
    chief_accountant_name   VARCHAR(200),
    chief_accountant_email  VARCHAR(200),

    -- Tax Office
    tax_office_code     VARCHAR(10),           -- Ma Chi cuc Thue quan ly
    tax_office_name     VARCHAR(200),

    -- Accounting Configuration
    accounting_regime   VARCHAR(10) NOT NULL DEFAULT 'TT99',  -- TT99, TT133, TT58
    fiscal_year_start_month INTEGER NOT NULL DEFAULT 1,       -- 1-12
    default_currency    VARCHAR(3) NOT NULL DEFAULT 'VND',
    secondary_currency  VARCHAR(3),            -- Optional: USD, EUR

    -- Company Type
    company_type        VARCHAR(50),           -- MANUFACTURING, TRADING, SERVICE, CONSTRUCTION, AGRI, OTHER
    company_size        VARCHAR(20),           -- MICRO, SMALL, MEDIUM, LARGE

    -- Status
    status              VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE, SUSPENDED, DISSOLVED, MERGED
    status_reason       TEXT,

    -- Parent
    parent_company_id   UUID REFERENCES companies(id),
    hierarchy_path      LTREE,                -- PostgreSQL ltree for hierarchy queries

    -- Metadata
    logo_url            TEXT,
    settings            JSONB DEFAULT '{}',    -- Company-specific preferences
    metadata            JSONB DEFAULT '{}',

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    archived_at         TIMESTAMPTZ,

    -- Constraints
    UNIQUE(tenant_id, tax_code),
    UNIQUE(tenant_id, business_reg_no)
);

CREATE INDEX idx_companies_tenant ON companies(tenant_id);
CREATE INDEX idx_companies_hierarchy ON companies USING GIST(hierarchy_path);
CREATE INDEX idx_companies_status ON companies(status);
CREATE INDEX idx_companies_tax_code ON companies(tax_code);
```

### 1.2 company_branches
```sql
CREATE TABLE company_branches (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    branch_name     VARCHAR(500) NOT NULL,
    branch_tax_code VARCHAR(20) NOT NULL,  -- MST chi nhanh: parent_MST + '-' + XXX
    branch_type     VARCHAR(50) NOT NULL DEFAULT 'BRANCH',  -- BRANCH, REP_OFFICE, DEPENDENT_UNIT
    address         TEXT,
    phone           VARCHAR(30),
    email           VARCHAR(200),
    manager_name    VARCHAR(200),
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(company_id, branch_tax_code)
);

CREATE INDEX idx_branches_company ON company_branches(company_id);
```

### 1.3 fiscal_years
```sql
CREATE TYPE period_status AS ENUM ('OPEN', 'CLOSED', 'PERMANENTLY_CLOSED', 'FUTURE');

CREATE TABLE fiscal_years (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    year            INTEGER NOT NULL,
    start_month     INTEGER NOT NULL DEFAULT 1,  -- 1-12
    is_short_year   BOOLEAN NOT NULL DEFAULT FALSE,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'OPEN',  -- OPEN, CLOSED

    UNIQUE(company_id, year)
);

CREATE TABLE periods (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    fiscal_year_id  UUID NOT NULL REFERENCES fiscal_years(id),
    period_type     VARCHAR(10) NOT NULL,        -- MONTHLY, QUARTERLY, ANNUAL
    period_number   INTEGER NOT NULL,            -- 1-12 for monthly, 1-4 for quarterly
    label           VARCHAR(50) NOT NULL,        -- "Thang 01/2026"
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    status          period_status NOT NULL DEFAULT 'FUTURE',
    opened_at       TIMESTAMPTZ,
    closed_at       TIMESTAMPTZ,
    closed_by       UUID,                        -- user_id
    reopened_count  INTEGER NOT NULL DEFAULT 0,

    UNIQUE(company_id, fiscal_year_id, period_type, period_number)
);

CREATE INDEX idx_periods_company ON periods(company_id);
CREATE INDEX idx_periods_fiscal_year ON periods(fiscal_year_id);
CREATE INDEX idx_periods_status ON periods(company_id, status);
```

### 1.4 departments
```sql
CREATE TABLE departments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    code            VARCHAR(20) NOT NULL,
    name            VARCHAR(200) NOT NULL,
    parent_id       UUID REFERENCES departments(id),
    manager_id      UUID,              -- references employees(id)
    hierarchy_path  LTREE,
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(company_id, code)
);

CREATE INDEX idx_depts_company ON departments(company_id);
CREATE INDEX idx_depts_hierarchy ON departments USING GIST(hierarchy_path);
```

### 1.5 employees
```sql
CREATE TABLE employees (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    employee_code   VARCHAR(30) NOT NULL,
    full_name       VARCHAR(200) NOT NULL,
    title           VARCHAR(200),
    email           VARCHAR(200),
    phone           VARCHAR(30),
    department_id   UUID REFERENCES departments(id),
    personal_tax_code VARCHAR(20),     -- MST ca nhan
    social_insurance_no VARCHAR(20),   -- BHXH
    bank_account_no VARCHAR(50),
    bank_id         UUID,              -- references company_bank_accounts
    user_id         UUID,              -- optional link to system user (1:1)
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE, LEAVE, TERMINATED
    hire_date       DATE,
    termination_date DATE,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(company_id, employee_code)
);

CREATE INDEX idx_employees_company ON employees(company_id);
CREATE INDEX idx_employees_department ON employees(department_id);
CREATE INDEX idx_employees_user ON employees(user_id);
```

### 1.6 company_bank_accounts
```sql
CREATE TABLE company_bank_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    bank_code       VARCHAR(10),         -- Napas bank code
    bank_name       VARCHAR(200) NOT NULL,
    branch_name     VARCHAR(200),
    account_number  VARCHAR(50) NOT NULL,
    account_holder  VARCHAR(200) NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'VND',
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE, CLOSED
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(company_id, account_number)
);

CREATE INDEX idx_bank_accounts_company ON company_bank_accounts(company_id);
```

### 1.7 einvoice_patterns
```sql
CREATE TABLE einvoice_patterns (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    pattern_code    VARCHAR(20) NOT NULL,  -- Mau so HD (e.g., "01GTKT0/001")
    serial          VARCHAR(20) NOT NULL,  -- Ky hieu (e.g., "AA/26E")
    form            VARCHAR(100),          -- Kieu HD
    invoice_type    VARCHAR(50),           -- GTGT, EXPORT, RETAIL, SERVICE
    status          VARCHAR(20) NOT NULL DEFAULT 'REGISTERED',  -- REGISTERED, ACTIVE, CANCELLED
    gdt_synced_at   TIMESTAMPTZ,
    gdt_status      VARCHAR(50),           -- GDT portal sync status
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(company_id, pattern_code, serial)
);

CREATE INDEX idx_einvoice_company ON einvoice_patterns(company_id);
```

### 1.8 digital_signatures
```sql
CREATE TABLE digital_signatures (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    signature_type  VARCHAR(20) NOT NULL,   -- USB_TOKEN, REMOTE_HSM
    provider        VARCHAR(100),           -- ViettelCA, VNPTCA, CloudCA, etc.
    serial_number   VARCHAR(200) NOT NULL,
    owner_name      VARCHAR(200),
    certificate_subject TEXT,
    certificate_issuer  TEXT,
    valid_from      DATE NOT NULL,
    valid_to        DATE NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE, EXPIRED, REVOKED
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    config          JSONB,                 -- Provider-specific config (encrypted)
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(company_id, serial_number)
);

CREATE INDEX idx_signatures_company ON digital_signatures(company_id);
CREATE INDEX idx_signatures_expiry ON digital_signatures(valid_to)
    WHERE status = 'ACTIVE';
```

### 1.9 integration_profiles
```sql
CREATE TABLE integration_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    integration_type VARCHAR(50) NOT NULL,  -- GDT, CUSTOMS, BHXH, DVC
    endpoint_url    VARCHAR(500),
    credentials     BYTEA,                  -- AES-256-GCM encrypted
    config          JSONB DEFAULT '{}',
    status          VARCHAR(20) NOT NULL DEFAULT 'DISCONNECTED',  -- CONNECTED, DISCONNECTED, ERROR
    last_connected_at TIMESTAMPTZ,
    last_error_at   TIMESTAMPTZ,
    last_error_msg  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(company_id, integration_type)
);

CREATE INDEX idx_integrations_company ON integration_profiles(company_id);
```

---

## 2. API Endpoints

### 2.1 Company Management
| Method | Path | Description | Permission |
|--------|------|-------------|------------|
| POST   | /api/v1/companies | Create company | company:create |
| GET    | /api/v1/companies | List companies (user's accessible) | company:list |
| GET    | /api/v1/companies/:id | Get company detail | company:read |
| PUT    | /api/v1/companies/:id | Update company | company:write |
| PATCH  | /api/v1/companies/:id | Partial update | company:write |
| DELETE | /api/v1/companies/:id | Deactivate company (soft) | company:delete |
| POST   | /api/v1/companies/switch | Switch current company context | company:switch |
| GET    | /api/v1/companies/:id/hierarchy | Get company hierarchy tree | company:read |

### 2.2 Branch Management
| Method | Path | Description | Permission |
|--------|------|-------------|------------|
| POST   | /api/v1/companies/:id/branches | Create branch | company:branch:create |
| GET    | /api/v1/companies/:id/branches | List branches | company:read |
| PUT    | /api/v1/branches/:id | Update branch | company:write |
| DELETE | /api/v1/branches/:id | Deactivate branch | company:write |

### 2.3 Fiscal Year / Period
| Method | Path | Description | Permission |
|--------|------|-------------|------------|
| POST   | /api/v1/companies/:id/fiscal-years | Configure fiscal year | company:fiscalyear:write |
| GET    | /api/v1/companies/:id/fiscal-years | List fiscal years | company:read |
| GET    | /api/v1/companies/:id/periods | List periods | company:read |
| POST   | /api/v1/companies/:id/periods/:pid/close | Close period | company:period:close |
| POST   | /api/v1/companies/:id/periods/:pid/reopen | Reopen period | company:period:reopen |
| POST   | /api/v1/companies/:id/periods/:pid/permanent-close | Permanently close | company:period:permanent-close |

### 2.4 Department / Employee
| Method | Path | Description | Permission |
|--------|------|-------------|------------|
| POST   | /api/v1/companies/:id/departments | Create department | company:dept:create |
| GET    | /api/v1/companies/:id/departments | List departments | company:read |
| POST   | /api/v1/companies/:id/employees | Create employee | company:employee:create |
| GET    | /api/v1/companies/:id/employees | List employees | company:list |
| GET    | /api/v1/employees/:id | Get employee | company:read |
| PUT    | /api/v1/employees/:id | Update employee | company:employee:write |
| DELETE | /api/v1/employees/:id | Deactivate employee | company:employee:write |

### 2.5 Bank Accounts
| Method | Path | Description | Permission |
|--------|------|-------------|------------|
| POST   | /api/v1/companies/:id/bank-accounts | Register bank account | company:bank:write |
| GET    | /api/v1/companies/:id/bank-accounts | List bank accounts | company:read |
| PUT    | /api/v1/bank-accounts/:id | Update bank account | company:bank:write |
| DELETE | /api/v1/bank-accounts/:id | Close bank account | company:bank:write |
| POST   | /api/v1/bank-accounts/:id/verify | Verify account | company:bank:verify |

### 2.6 E-Invoice Patterns
| Method | Path | Description | Permission |
|--------|------|-------------|------------|
| POST   | /api/v1/companies/:id/einvoice-patterns | Register pattern | company:einvoice:write |
| GET    | /api/v1/companies/:id/einvoice-patterns | List patterns | company:read |
| POST   | /api/v1/einvoice-patterns/:id/sync | Sync with GDT | company:einvoice:write |

### 2.7 Digital Signatures
| Method | Path | Description | Permission |
|--------|------|-------------|------------|
| POST   | /api/v1/companies/:id/digital-signatures | Register signature | company:signature:write |
| GET    | /api/v1/companies/:id/digital-signatures | List signatures | company:read |
| POST   | /api/v1/digital-signatures/:id/test | Test signature | company:signature:write |

### 2.8 Integration Profiles
| Method | Path | Description | Permission |
|--------|------|-------------|------------|
| POST   | /api/v1/companies/:id/integrations | Configure integration | company:integration:write |
| GET    | /api/v1/companies/:id/integrations | List integrations | company:read |
| POST   | /api/v1/integrations/:id/test | Test connection | company:integration:write |
| PUT    | /api/v1/integrations/:id | Update credentials | company:integration:write |

---

## 3. Service Layer Interface

```go
type CompanyService interface {
    // Company CRUD
    Create(ctx context.Context, tenantID uuid.UUID, input CreateCompanyInput) (*Company, error)
    Get(ctx context.Context, companyID uuid.UUID) (*Company, error)
    List(ctx context.Context, tenantID uuid.UUID, filter CompanyFilter) ([]*Company, error)
    Update(ctx context.Context, companyID uuid.UUID, input UpdateCompanyInput) (*Company, error)
    Deactivate(ctx context.Context, companyID uuid.UUID, reason string) error
    SwitchCompany(ctx context.Context, userID, companyID uuid.UUID) (*CompanyContext, error)

    // Hierarchy
    GetHierarchy(ctx context.Context, companyID uuid.UUID) (*CompanyNode, error)
    AddBranch(ctx context.Context, companyID uuid.UUID, input BranchInput) (*Branch, error)

    // Fiscal Year
    CreateFiscalYear(ctx context.Context, companyID uuid.UUID, input FiscalYearInput) (*FiscalYear, error)
    ListPeriods(ctx context.Context, companyID uuid.UUID) ([]*Period, error)
    ClosePeriod(ctx context.Context, companyID, periodID, userID uuid.UUID) error
    ReopenPeriod(ctx context.Context, companyID, periodID, userID uuid.UUID) error
    PermanentClosePeriod(ctx context.Context, companyID, periodID, userID uuid.UUID) error

    // Employee
    CreateEmployee(ctx context.Context, companyID uuid.UUID, input EmployeeInput) (*Employee, error)
    ListEmployees(ctx context.Context, companyID uuid.UUID, filter EmployeeFilter) ([]*Employee, error)

    // Department
    CreateDepartment(ctx context.Context, companyID uuid.UUID, input DepartmentInput) (*Department, error)

    // Bank Account
    RegisterBankAccount(ctx context.Context, companyID uuid.UUID, input BankAccountInput) (*BankAccount, error)

    // E-Invoice
    RegisterEInvoicePattern(ctx context.Context, companyID uuid.UUID, input EInvoiceInput) (*EInvoicePattern, error)

    // Digital Signature
    RegisterSignature(ctx context.Context, companyID uuid.UUID, input SignatureInput) (*DigitalSignature, error)

    // Integration
    ConfigureIntegration(ctx context.Context, companyID uuid.UUID, input IntegrationInput) (*IntegrationProfile, error)
    TestIntegration(ctx context.Context, integrationID uuid.UUID) error
}
```

---

## 4. Security

### 4.1 Data Isolation
- Every table has `company_id` column
- Every query includes `WHERE company_id = current_context`
- RLS policies on all tables
- API never accepts `company_id` from request body — extracted from JWT
- Separate database user per tenant (optional, for high-security deployments)

### 4.2 Encryption
- Integration credentials: AES-256-GCM at rest
- Employee PII: field-level encryption (name, personal_tax_code, bank_account)
- Audit log: immutable, append-only

### 4.3 Audit Trail
```sql
CREATE TABLE company_audit_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL,
    actor_id        UUID NOT NULL,
    action          VARCHAR(50) NOT NULL,  -- CREATE, UPDATE, DEACTIVATE, CLOSE_PERIOD, REOPEN_PERIOD, SWITCH_COMPANY
    entity_type     VARCHAR(50) NOT NULL,  -- COMPANY, BRANCH, EMPLOYEE, DEPARTMENT, BANK_ACCOUNT, EINVOICE, SIGNATURE, INTEGRATION
    entity_id       UUID,
    old_values      JSONB,
    new_values      JSONB,
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Retention: 5 years per Law on Accounting Art. 13
```

---

## 5. Performance Targets

| Operation | Target p50 | Target p95 | Target p99 |
|-----------|-----------|-----------|-----------|
| Company CRUD | 50ms | 150ms | 300ms |
| Period close | 200ms | 500ms | 1s |
| Employee list (1000 records) | 100ms | 300ms | 500ms |
| Company switch | 100ms | 300ms | 500ms |
| Hierarchy tree (100 nodes) | 200ms | 500ms | 1s |
| Integration test | 2s | 5s | 10s |
