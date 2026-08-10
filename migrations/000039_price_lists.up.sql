CREATE TABLE IF NOT EXISTS price_lists (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    currency VARCHAR(3) DEFAULT 'VND',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_pl_company_code UNIQUE (company_id, code)
);

CREATE INDEX IF NOT EXISTS idx_pl_company ON price_lists(company_id);

CREATE TABLE IF NOT EXISTS price_list_lines (
    id VARCHAR(36) PRIMARY KEY,
    price_list_id VARCHAR(36) NOT NULL REFERENCES price_lists(id) ON DELETE CASCADE,
    item_code VARCHAR(50) NOT NULL,
    item_name VARCHAR(255) DEFAULT '',
    unit VARCHAR(20) DEFAULT '',
    unit_price NUMERIC(18,2) NOT NULL DEFAULT 0,
    vat_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    min_quantity NUMERIC(12,2) NOT NULL DEFAULT 0,
    effective_from VARCHAR(10) DEFAULT '',
    effective_to VARCHAR(10) DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_pll_pl ON price_list_lines(price_list_id);
