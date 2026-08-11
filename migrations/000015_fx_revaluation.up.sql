CREATE TABLE IF NOT EXISTS fx_revaluations (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    revaluation_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    total_gain DOUBLE PRECISION NOT NULL DEFAULT 0,
    total_loss DOUBLE PRECISION NOT NULL DEFAULT 0,
    gl_posted BOOLEAN NOT NULL DEFAULT FALSE,
    gl_posted_at VARCHAR(40) NOT NULL DEFAULT '',
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fxr_company ON fx_revaluations (company_id);
CREATE INDEX IF NOT EXISTS idx_fxr_status ON fx_revaluations (status);

CREATE TABLE IF NOT EXISTS fx_revaluation_lines (
    id VARCHAR(36) PRIMARY KEY,
    revaluation_id VARCHAR(36) NOT NULL REFERENCES fx_revaluations(id) ON DELETE CASCADE,
    invoice_id VARCHAR(36) NOT NULL,
    invoice_number VARCHAR(30) NOT NULL,
    supplier_id VARCHAR(36) NOT NULL,
    supplier_name VARCHAR(255) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    balance_due DOUBLE PRECISION NOT NULL DEFAULT 0,
    original_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
    revaluation_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
    fx_gain DOUBLE PRECISION NOT NULL DEFAULT 0,
    fx_loss DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_fxrl_revaluation ON fx_revaluation_lines (revaluation_id);
CREATE INDEX IF NOT EXISTS idx_fxrl_invoice ON fx_revaluation_lines (invoice_id);
