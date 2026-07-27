-- GoTax Company Domain Schema v3
-- Multi-tenant company management, branches, fiscal years, employees, etc.
-- PostgreSQL 16+

BEGIN;

-- ============================================================
-- COMPANIES (doanh nghiep)
-- ============================================================
CREATE TABLE IF NOT EXISTS companies (
    id                     VARCHAR(20) PRIMARY KEY,
    tenant_id              VARCHAR(50) NOT NULL,
    legal_name_vn          VARCHAR(255) NOT NULL,
    legal_name_en          VARCHAR(255),
    short_name             VARCHAR(100),
    legal_form             VARCHAR(20) NOT NULL
                           CHECK (legal_form IN ('LLC_1MEMBER','LLC_2MEMBERS','JSC','SOLE_PROP','PARTNERSHIP','FO','RO')),
    tax_code               VARCHAR(20) NOT NULL,
    business_reg_no        VARCHAR(50),
    business_reg_date      VARCHAR(20),
    business_reg_place     VARCHAR(255),
    reg_address            VARCHAR(255) NOT NULL DEFAULT '',
    reg_province           VARCHAR(100),
    reg_district           VARCHAR(100),
    head_office_address    VARCHAR(255),
    head_office_province   VARCHAR(100),
    head_office_district   VARCHAR(100),
    phone                  VARCHAR(30),
    email                  VARCHAR(100),
    website                VARCHAR(255),
    legal_rep_name         VARCHAR(100),
    legal_rep_title        VARCHAR(100),
    legal_rep_id_number    VARCHAR(30),
    chief_accountant       VARCHAR(100),
    chief_accountant_email VARCHAR(100),
    tax_office_code        VARCHAR(20),
    tax_office_name        VARCHAR(255),
    accounting_regime      VARCHAR(10) NOT NULL
                           CHECK (accounting_regime IN ('TT99','TT133','TT58')),
    fiscal_year_start_month INTEGER NOT NULL DEFAULT 1,
    default_currency       VARCHAR(3) NOT NULL DEFAULT 'VND',
    secondary_currency     VARCHAR(3),
    company_type           VARCHAR(20)
                           CHECK (company_type IN ('MANUFACTURING','TRADING','SERVICE','CONSTRUCTION','AGRICULTURE','OTHER')),
    company_size           VARCHAR(10)
                           CHECK (company_size IN ('MICRO','SMALL','MEDIUM','LARGE')),
    status                 VARCHAR(15) NOT NULL DEFAULT 'ACTIVE'
                           CHECK (status IN ('ACTIVE','SUSPENDED','DISSOLVED','MERGED')),
    parent_company_id      VARCHAR(20),
    logo_url               VARCHAR(500),
    settings               JSONB DEFAULT '{}',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_companies_tenant ON companies(tenant_id);
CREATE INDEX IF NOT EXISTS idx_companies_tax_code ON companies(tax_code);
CREATE INDEX IF NOT EXISTS idx_companies_status ON companies(status);

-- ============================================================
-- COMPANY BRANCHES (chi nhanh)
-- ============================================================
CREATE TABLE IF NOT EXISTS company_branches (
    id              VARCHAR(20) PRIMARY KEY,
    company_id      VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    branch_name     VARCHAR(255) NOT NULL,
    branch_tax_code VARCHAR(20) NOT NULL,
    branch_type     VARCHAR(20) NOT NULL DEFAULT 'BRANCH'
                    CHECK (branch_type IN ('BRANCH','REP_OFFICE','DEPENDENT_UNIT')),
    address         VARCHAR(255),
    phone           VARCHAR(30),
    email           VARCHAR(100),
    manager_name    VARCHAR(100),
    status          VARCHAR(15) NOT NULL DEFAULT 'ACTIVE',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_branches_company ON company_branches(company_id);

-- ============================================================
-- FISCAL YEARS (nam tai chinh)
-- ============================================================
CREATE TABLE IF NOT EXISTS fiscal_years (
    id            VARCHAR(20) PRIMARY KEY,
    company_id    VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    year          INTEGER NOT NULL,
    start_month   INTEGER NOT NULL DEFAULT 1,
    is_short_year BOOLEAN NOT NULL DEFAULT FALSE,
    start_date    VARCHAR(20),
    end_date      VARCHAR(20),
    status        VARCHAR(15) NOT NULL DEFAULT 'ACTIVE',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fy_company ON fiscal_years(company_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_fy_company_year ON fiscal_years(company_id, year);

-- ============================================================
-- PERIODS V2 (ky ke toan — company-managed)
-- ============================================================
CREATE TABLE IF NOT EXISTS periods_v2 (
    id              VARCHAR(20) PRIMARY KEY,
    company_id      VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    fiscal_year_id  VARCHAR(20) NOT NULL REFERENCES fiscal_years(id) ON DELETE CASCADE,
    period_type     VARCHAR(10) NOT NULL DEFAULT 'MONTHLY'
                    CHECK (period_type IN ('MONTHLY','QUARTERLY','ANNUAL')),
    period_number   INTEGER NOT NULL,
    label           VARCHAR(50) NOT NULL,
    start_date      VARCHAR(20) NOT NULL,
    end_date        VARCHAR(20) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'FUTURE'
                    CHECK (status IN ('OPEN','CLOSED','PERMANENTLY_CLOSED','FUTURE')),
    opened_at       VARCHAR(30),
    closed_at       VARCHAR(30),
    closed_by       VARCHAR(100),
    reopened_count  INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_periods_v2_company ON periods_v2(company_id);
CREATE INDEX IF NOT EXISTS idx_periods_v2_fiscal_year ON periods_v2(fiscal_year_id);
CREATE INDEX IF NOT EXISTS idx_periods_v2_status ON periods_v2(status);

-- ============================================================
-- DEPARTMENTS (phong ban)
-- ============================================================
CREATE TABLE IF NOT EXISTS departments (
    id          VARCHAR(20) PRIMARY KEY,
    company_id  VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    code        VARCHAR(50) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    parent_id   VARCHAR(20),
    manager_id  VARCHAR(20),
    status      VARCHAR(15) NOT NULL DEFAULT 'ACTIVE',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_depts_company ON departments(company_id);

-- ============================================================
-- EMPLOYEES (nhan vien)
-- ============================================================
CREATE TABLE IF NOT EXISTS employees (
    id                    VARCHAR(20) PRIMARY KEY,
    company_id            VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_code         VARCHAR(50) NOT NULL,
    full_name             VARCHAR(255) NOT NULL,
    title                 VARCHAR(100),
    email                 VARCHAR(100),
    phone                 VARCHAR(30),
    department_id         VARCHAR(20),
    personal_tax_code     VARCHAR(20),
    social_insurance_no   VARCHAR(20),
    bank_account_no       VARCHAR(30),
    user_id               VARCHAR(50),
    status                VARCHAR(15) NOT NULL DEFAULT 'ACTIVE'
                          CHECK (status IN ('ACTIVE','LEAVE','TERMINATED')),
    hire_date             VARCHAR(20),
    termination_date      VARCHAR(20),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_employees_company ON employees(company_id);
CREATE INDEX IF NOT EXISTS idx_employees_code ON employees(company_id, employee_code);

-- ============================================================
-- COMPANY BANK ACCOUNTS (tai khoan ngan hang)
-- ============================================================
CREATE TABLE IF NOT EXISTS company_bank_accounts (
    id              VARCHAR(20) PRIMARY KEY,
    company_id      VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    bank_code       VARCHAR(20),
    bank_name       VARCHAR(255) NOT NULL,
    branch_name     VARCHAR(255),
    account_number  VARCHAR(50) NOT NULL,
    account_holder  VARCHAR(255) NOT NULL DEFAULT '',
    currency        VARCHAR(3) NOT NULL DEFAULT 'VND',
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    status          VARCHAR(15) NOT NULL DEFAULT 'ACTIVE'
                    CHECK (status IN ('ACTIVE','CLOSED')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bank_accounts_company ON company_bank_accounts(company_id);

-- ============================================================
-- E-INVOICE PATTERNS (mau hoa don dien tu)
-- ============================================================
CREATE TABLE IF NOT EXISTS e_invoice_patterns (
    id              VARCHAR(20) PRIMARY KEY,
    company_id      VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    pattern_code    VARCHAR(50) NOT NULL,
    serial          VARCHAR(20) NOT NULL,
    form            VARCHAR(50),
    invoice_type    VARCHAR(50),
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE'
                    CHECK (status IN ('REGISTERED','ACTIVE','CANCELLED')),
    gdt_status      VARCHAR(50),
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_einvoice_company ON e_invoice_patterns(company_id);

-- ============================================================
-- DIGITAL SIGNATURES (chu ky so)
-- ============================================================
CREATE TABLE IF NOT EXISTS digital_signatures (
    id                    VARCHAR(20) PRIMARY KEY,
    company_id            VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    signature_type        VARCHAR(20) NOT NULL DEFAULT 'USB_TOKEN'
                          CHECK (signature_type IN ('USB_TOKEN','REMOTE_HSM')),
    provider              VARCHAR(100),
    serial_number         VARCHAR(100) NOT NULL,
    owner_name            VARCHAR(255),
    certificate_subject   TEXT,
    certificate_issuer    TEXT,
    valid_from            VARCHAR(20) NOT NULL DEFAULT '',
    valid_to              VARCHAR(20) NOT NULL DEFAULT '',
    status                VARCHAR(15) NOT NULL DEFAULT 'ACTIVE'
                          CHECK (status IN ('ACTIVE','EXPIRED','REVOKED')),
    is_default            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_digsig_company ON digital_signatures(company_id);

-- ============================================================
-- INTEGRATION PROFILES (ket noi voi he thong ben ngoai)
-- ============================================================
CREATE TABLE IF NOT EXISTS integration_profiles (
    id                VARCHAR(20) PRIMARY KEY,
    company_id        VARCHAR(20) NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    integration_type  VARCHAR(20) NOT NULL
                      CHECK (integration_type IN ('GDT','CUSTOMS','BHXH','DVC')),
    endpoint_url      VARCHAR(500),
    status            VARCHAR(20) NOT NULL DEFAULT 'DISCONNECTED'
                      CHECK (status IN ('CONNECTED','DISCONNECTED','ERROR')),
    last_connected_at VARCHAR(30),
    last_error_at     VARCHAR(30),
    last_error_msg    TEXT,
    config            JSONB DEFAULT '{}',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_integrations_company ON integration_profiles(company_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_integrations_type ON integration_profiles(company_id, integration_type);

COMMIT;
