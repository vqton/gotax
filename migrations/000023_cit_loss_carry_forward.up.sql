CREATE TABLE IF NOT EXISTS cit_loss_carry_forwards (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    loss_year INTEGER NOT NULL,
    loss_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    used_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    expiry_year INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cit_loss_company_year ON cit_loss_carry_forwards(company_id, loss_year);
