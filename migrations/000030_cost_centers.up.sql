CREATE TABLE IF NOT EXISTS cost_centers (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id VARCHAR(36) NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(36),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_cost_centers_company ON cost_centers(company_id);
