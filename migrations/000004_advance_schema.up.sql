CREATE TABLE IF NOT EXISTS advance_requests (
    id            TEXT PRIMARY KEY,
    company_id    TEXT NOT NULL,
    requestor_id  TEXT NOT NULL,
    requestor_name TEXT NOT NULL DEFAULT '',
    amount        DOUBLE PRECISION NOT NULL,
    amount_vnd    DOUBLE PRECISION NOT NULL DEFAULT 0,
    currency      TEXT NOT NULL DEFAULT 'VND',
    exchange_rate DOUBLE PRECISION NOT NULL DEFAULT 1,
    purpose       TEXT NOT NULL,
    status        TEXT NOT NULL DEFAULT 'DRAFT',
    approved_by   TEXT NOT NULL DEFAULT '',
    approved_at   TEXT NOT NULL DEFAULT '',
    paid_by       TEXT NOT NULL DEFAULT '',
    paid_at       TEXT NOT NULL DEFAULT '',
    gl_journal_id TEXT NOT NULL DEFAULT '',
    created_by    TEXT NOT NULL,
    created_at    TEXT NOT NULL,
    updated_at    TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_advance_company ON advance_requests(company_id);
CREATE INDEX IF NOT EXISTS idx_advance_status ON advance_requests(status);

CREATE TABLE IF NOT EXISTS advance_settlements (
    id               TEXT PRIMARY KEY,
    advance_id       TEXT NOT NULL REFERENCES advance_requests(id),
    company_id       TEXT NOT NULL,
    total_spent      DOUBLE PRECISION NOT NULL DEFAULT 0,
    remaining_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    currency         TEXT NOT NULL DEFAULT 'VND',
    notes            TEXT NOT NULL DEFAULT '',
    status           TEXT NOT NULL DEFAULT 'DRAFT',
    created_at       TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_settlement_advance ON advance_settlements(advance_id);
