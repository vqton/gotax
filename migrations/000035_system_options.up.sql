CREATE TABLE IF NOT EXISTS system_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    category VARCHAR(50) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value JSONB NOT NULL DEFAULT '"{}"',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, category, key)
);

CREATE INDEX IF NOT EXISTS idx_system_options_company ON system_options(company_id);

CREATE TABLE IF NOT EXISTS numbering_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    voucher_type VARCHAR(50) NOT NULL,
    prefix VARCHAR(20) DEFAULT '',
    suffix VARCHAR(20) DEFAULT '',
    number_length INT NOT NULL DEFAULT 5,
    scope VARCHAR(20) NOT NULL DEFAULT 'company',
    reset_rule VARCHAR(20) NOT NULL DEFAULT 'never',
    current_num INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, voucher_type)
);

CREATE INDEX IF NOT EXISTS idx_numbering_rules_company ON numbering_rules(company_id);

CREATE TABLE IF NOT EXISTS backup_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    filename VARCHAR(255) NOT NULL,
    file_size BIGINT DEFAULT 0,
    backup_type VARCHAR(20) NOT NULL DEFAULT 'manual',
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_backup_records_company ON backup_records(company_id);
