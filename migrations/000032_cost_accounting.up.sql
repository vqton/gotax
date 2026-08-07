-- Cost Accounting module: CostObject, CostPool, CostPoolLine, CostingPeriod, CostingResult, CostingResultLine
-- Circular 99/2025/TT-BTC: TK 154 (WIP), TK 621-627 (cost pools), TK 632 (COGS)

CREATE TABLE IF NOT EXISTS cost_objects (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id      TEXT NOT NULL REFERENCES companies(id),
    code            TEXT NOT NULL,
    name            TEXT NOT NULL,
    type            TEXT NOT NULL DEFAULT 'PRODUCT',
    costing_method  TEXT NOT NULL DEFAULT 'SIMPLE',
    cost_center_id  TEXT,
    unit_of_measure TEXT,
    standard_cost   NUMERIC(18,2) DEFAULT 0,
    plan_quantity   NUMERIC(18,2) DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);

CREATE INDEX IF NOT EXISTS idx_cost_objects_company ON cost_objects(company_id);

CREATE TABLE IF NOT EXISTS cost_pools (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id      TEXT NOT NULL REFERENCES companies(id),
    period_id       TEXT NOT NULL,
    gl_account_code TEXT NOT NULL,
    name            TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'OPEN',
    total_amount    NUMERIC(18,2) DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cost_pools_company_period ON cost_pools(company_id, period_id);

CREATE TABLE IF NOT EXISTS cost_pool_lines (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    pool_id         TEXT NOT NULL REFERENCES cost_pools(id) ON DELETE CASCADE,
    source_type     TEXT NOT NULL,
    source_id       TEXT,
    description     TEXT NOT NULL DEFAULT '',
    amount          NUMERIC(18,2) NOT NULL DEFAULT 0,
    cost_center_id  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cost_pool_lines_pool ON cost_pool_lines(pool_id);

CREATE TABLE IF NOT EXISTS costing_periods (
    id          TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id  TEXT NOT NULL REFERENCES companies(id),
    year        INT NOT NULL,
    month       INT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'OPEN',
    closed_by   TEXT,
    closed_at   TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, year, month)
);

CREATE INDEX IF NOT EXISTS idx_costing_periods_company ON costing_periods(company_id);

CREATE TABLE IF NOT EXISTS costing_results (
    id               TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    company_id       TEXT NOT NULL REFERENCES companies(id),
    period_id        TEXT NOT NULL REFERENCES costing_periods(id),
    cost_object_id   TEXT NOT NULL REFERENCES cost_objects(id),
    costing_method   TEXT NOT NULL,
    total_direct_mat NUMERIC(18,2) DEFAULT 0,
    total_direct_lab NUMERIC(18,2) DEFAULT 0,
    total_overhead   NUMERIC(18,2) DEFAULT 0,
    total_cost       NUMERIC(18,2) DEFAULT 0,
    output_quantity  NUMERIC(18,2) DEFAULT 0,
    unit_cost        NUMERIC(18,4) DEFAULT 0,
    wip_beginning    NUMERIC(18,2) DEFAULT 0,
    wip_ending       NUMERIC(18,2) DEFAULT 0,
    status           TEXT NOT NULL DEFAULT 'DRAFT',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(period_id, cost_object_id)
);

CREATE INDEX IF NOT EXISTS idx_costing_results_period ON costing_results(period_id);
CREATE INDEX IF NOT EXISTS idx_costing_results_company ON costing_results(company_id);

CREATE TABLE IF NOT EXISTS costing_result_lines (
    id               TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    result_id        TEXT NOT NULL REFERENCES costing_results(id) ON DELETE CASCADE,
    cost_category    TEXT NOT NULL,
    gl_account_code  TEXT NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    planned_amount   NUMERIC(18,2) DEFAULT 0,
    actual_amount    NUMERIC(18,2) DEFAULT 0,
    allocated_amount NUMERIC(18,2) DEFAULT 0,
    coefficient      NUMERIC(18,6) DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_costing_result_lines_result ON costing_result_lines(result_id);
