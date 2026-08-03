CREATE TABLE IF NOT EXISTS doubtful_debt_provisions (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    as_of_date DATE NOT NULL,
    total_outstanding DOUBLE PRECISION NOT NULL DEFAULT 0,
    total_provision DOUBLE PRECISION NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ddp_company ON doubtful_debt_provisions(company_id);
CREATE INDEX IF NOT EXISTS idx_ddp_status ON doubtful_debt_provisions(status);
CREATE INDEX IF NOT EXISTS idx_ddp_date ON doubtful_debt_provisions(as_of_date);

CREATE TABLE IF NOT EXISTS doubtful_debt_provision_lines (
    id VARCHAR(36) PRIMARY KEY,
    provision_id VARCHAR(36) NOT NULL REFERENCES doubtful_debt_provisions(id) ON DELETE CASCADE,
    supplier_id VARCHAR(36) NOT NULL,
    supplier_name VARCHAR(255) NOT NULL,
    tax_code VARCHAR(50),
    outstanding_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    age_months INTEGER NOT NULL DEFAULT 0,
    rate_pct DOUBLE PRECISION NOT NULL DEFAULT 0,
    provision_amount DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_ddpl_provision ON doubtful_debt_provision_lines(provision_id);
CREATE INDEX IF NOT EXISTS idx_ddpl_supplier ON doubtful_debt_provision_lines(supplier_id);
