CREATE TABLE IF NOT EXISTS budgets (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id VARCHAR(36) NOT NULL,
    account_code VARCHAR(20) NOT NULL,
    period_year INT NOT NULL,
    period_month INT NOT NULL,
    budgeted DOUBLE PRECISION NOT NULL DEFAULT 0,
    actual DOUBLE PRECISION NOT NULL DEFAULT 0,
    variance DOUBLE PRECISION NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, account_code, period_year, period_month)
);
CREATE INDEX IF NOT EXISTS idx_budgets_company_year ON budgets(company_id, period_year);
