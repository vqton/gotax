CREATE TABLE IF NOT EXISTS warehouse_keeper_assignments (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id VARCHAR(36) NOT NULL,
    warehouse_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'keeper',
    effective_from DATE NOT NULL,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, warehouse_id, effective_from)
);
CREATE INDEX IF NOT EXISTS idx_keeper_assignments_company ON warehouse_keeper_assignments(company_id);
CREATE INDEX IF NOT EXISTS idx_keeper_assignments_warehouse ON warehouse_keeper_assignments(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_keeper_assignments_user ON warehouse_keeper_assignments(user_id);

CREATE TABLE IF NOT EXISTS stock_ledger_entries (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id VARCHAR(36) NOT NULL,
    warehouse_id VARCHAR(36) NOT NULL,
    item_id VARCHAR(36) NOT NULL,
    entry_date DATE NOT NULL,
    voucher_type VARCHAR(30) NOT NULL,
    voucher_no VARCHAR(50),
    voucher_ref_id VARCHAR(36),
    description TEXT,
    receipt_qty NUMERIC(18,4) NOT NULL DEFAULT 0,
    issue_qty NUMERIC(18,4) NOT NULL DEFAULT 0,
    balance_qty NUMERIC(18,4) NOT NULL DEFAULT 0,
    unit_cost NUMERIC(18,4),
    total_value NUMERIC(18,2),
    recorded_by VARCHAR(36) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    unrecorded_by VARCHAR(36),
    unrecorded_at TIMESTAMPTZ,
    unrecord_reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'recorded',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ledger_company_warehouse ON stock_ledger_entries(company_id, warehouse_id);
CREATE INDEX IF NOT EXISTS idx_ledger_item ON stock_ledger_entries(item_id);
CREATE INDEX IF NOT EXISTS idx_ledger_date ON stock_ledger_entries(entry_date);
CREATE INDEX IF NOT EXISTS idx_ledger_voucher_ref ON stock_ledger_entries(voucher_ref_id);
CREATE INDEX IF NOT EXISTS idx_ledger_recorded_by ON stock_ledger_entries(recorded_by);

CREATE TABLE IF NOT EXISTS keeper_inventory_counts (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    stock_take_id VARCHAR(36) NOT NULL,
    keeper_id VARCHAR(36) NOT NULL,
    book_qty_hidden BOOLEAN NOT NULL DEFAULT false,
    keeper_count_date DATE,
    keeper_signature VARCHAR(255),
    accountant_signature VARCHAR(255),
    manager_signature VARCHAR(255),
    status VARCHAR(30) NOT NULL DEFAULT 'pending_keeper',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_keeper_counts_stock_take ON keeper_inventory_counts(stock_take_id);
CREATE INDEX IF NOT EXISTS idx_keeper_counts_keeper ON keeper_inventory_counts(keeper_id);

CREATE TABLE IF NOT EXISTS warehouse_keeper_config (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id VARCHAR(36) NOT NULL,
    module_enabled BOOLEAN NOT NULL DEFAULT true,
    cost_price_hidden_from_keeper BOOLEAN NOT NULL DEFAULT false,
    auto_record_on_grn BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id)
);
CREATE INDEX IF NOT EXISTS idx_keeper_config_company ON warehouse_keeper_config(company_id);
