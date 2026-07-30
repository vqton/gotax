CREATE TABLE IF NOT EXISTS fixed_asset_categories (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    code VARCHAR(30) NOT NULL,
    name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(36),
    level INTEGER NOT NULL DEFAULT 1,
    default_useful_life_months INTEGER NOT NULL DEFAULT 120,
    default_depreciation_method VARCHAR(30) NOT NULL DEFAULT 'STRAIGHT_LINE',
    asset_account_id VARCHAR(36) NOT NULL,
    depreciation_account_id VARCHAR(36) NOT NULL,
    expense_account_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_facat_company_code ON fixed_asset_categories(company_id, code);
CREATE INDEX IF NOT EXISTS idx_facat_parent ON fixed_asset_categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_facat_level ON fixed_asset_categories(level);

CREATE TABLE IF NOT EXISTS fixed_assets (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category_id VARCHAR(36) NOT NULL REFERENCES fixed_asset_categories(id),
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    acquisition_date DATE NOT NULL,
    original_cost DOUBLE PRECISION NOT NULL DEFAULT 0,
    accumulated_depreciation DOUBLE PRECISION NOT NULL DEFAULT 0,
    residual_value DOUBLE PRECISION NOT NULL DEFAULT 0,
    carrying_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    useful_life_months INTEGER NOT NULL,
    depreciation_method VARCHAR(30) NOT NULL DEFAULT 'STRAIGHT_LINE',
    depreciation_start_date DATE,
    depreciation_end_date DATE,
    department_id VARCHAR(36) NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT '',
    user_id VARCHAR(36),
    supplier_id VARCHAR(36),
    contract_no VARCHAR(50),
    invoice_id VARCHAR(36),
    serial_no VARCHAR(100),
    manufacturer VARCHAR(255),
    manufacture_year INTEGER DEFAULT 0,
    country_of_origin VARCHAR(100),
    technical_specs TEXT,
    notes TEXT,
    source VARCHAR(30) NOT NULL DEFAULT 'PURCHASE',
    cip_account_id VARCHAR(36),
    asset_account_id VARCHAR(36) NOT NULL,
    depreciation_account_id VARCHAR(36) NOT NULL,
    expense_account_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(36) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(36) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_fa_company_code ON fixed_assets(company_id, code);
CREATE INDEX IF NOT EXISTS idx_fa_category ON fixed_assets(category_id);
CREATE INDEX IF NOT EXISTS idx_fa_status ON fixed_assets(status);
CREATE INDEX IF NOT EXISTS idx_fa_department ON fixed_assets(department_id);
CREATE INDEX IF NOT EXISTS idx_fa_source ON fixed_assets(source);

CREATE TABLE IF NOT EXISTS depreciation_entries (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    fixed_asset_id VARCHAR(36) NOT NULL REFERENCES fixed_assets(id) ON DELETE CASCADE,
    period_id VARCHAR(36) NOT NULL,
    period_year INTEGER NOT NULL,
    period_month INTEGER NOT NULL,
    depreciation_amount DOUBLE PRECISION NOT NULL,
    accumulated_after DOUBLE PRECISION NOT NULL,
    carrying_amount_after DOUBLE PRECISION NOT NULL,
    gl_posted BOOLEAN NOT NULL DEFAULT FALSE,
    gl_journal_entry_id VARCHAR(36),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(36) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_depr_asset_period ON depreciation_entries(fixed_asset_id, period_id);
CREATE INDEX IF NOT EXISTS idx_depr_period ON depreciation_entries(period_id);
CREATE INDEX IF NOT EXISTS idx_depr_gl_posted ON depreciation_entries(gl_posted);
CREATE INDEX IF NOT EXISTS idx_depr_company ON depreciation_entries(company_id);

CREATE TABLE IF NOT EXISTS fixed_asset_transactions (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    fixed_asset_id VARCHAR(36) NOT NULL REFERENCES fixed_assets(id) ON DELETE CASCADE,
    transaction_type VARCHAR(30) NOT NULL,
    transaction_date DATE NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    old_value DOUBLE PRECISION DEFAULT 0,
    new_value DOUBLE PRECISION DEFAULT 0,
    description TEXT,
    gl_journal_id VARCHAR(36),
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fat_asset ON fixed_asset_transactions(fixed_asset_id);
CREATE INDEX IF NOT EXISTS idx_fat_type ON fixed_asset_transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_fat_date ON fixed_asset_transactions(transaction_date);

CREATE TABLE IF NOT EXISTS fixed_asset_allocations (
    id VARCHAR(36) PRIMARY KEY,
    fixed_asset_id VARCHAR(36) NOT NULL REFERENCES fixed_assets(id) ON DELETE CASCADE,
    department_id VARCHAR(36) NOT NULL,
    allocation_pct DOUBLE PRECISION NOT NULL,
    expense_account_id VARCHAR(36) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_faalloc_asset ON fixed_asset_allocations(fixed_asset_id);
CREATE INDEX IF NOT EXISTS idx_faalloc_dept ON fixed_asset_allocations(department_id);

CREATE TABLE IF NOT EXISTS fixed_asset_inventory_plans (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    plan_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    notes TEXT,
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_faip_company ON fixed_asset_inventory_plans(company_id);
CREATE INDEX IF NOT EXISTS idx_faip_status ON fixed_asset_inventory_plans(status);
CREATE INDEX IF NOT EXISTS idx_faip_date ON fixed_asset_inventory_plans(plan_date);

CREATE TABLE IF NOT EXISTS fixed_asset_inventory_results (
    id VARCHAR(36) PRIMARY KEY,
    plan_id VARCHAR(36) NOT NULL REFERENCES fixed_asset_inventory_plans(id) ON DELETE CASCADE,
    fixed_asset_id VARCHAR(36) NOT NULL REFERENCES fixed_assets(id),
    expected_location VARCHAR(255),
    actual_location VARCHAR(255),
    expected_status VARCHAR(30),
    actual_status VARCHAR(30),
    discrepancy VARCHAR(50) NOT NULL,
    notes TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_fair_plan_asset ON fixed_asset_inventory_results(plan_id, fixed_asset_id);
CREATE INDEX IF NOT EXISTS idx_fair_discrepancy ON fixed_asset_inventory_results(discrepancy);
