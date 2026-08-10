CREATE TABLE IF NOT EXISTS invoice_books (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    pattern VARCHAR(50) NOT NULL,
    serial VARCHAR(20),
    from_number INTEGER NOT NULL,
    to_number INTEGER NOT NULL,
    next_number INTEGER NOT NULL,
    used_count INTEGER DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ib_company ON invoice_books(company_id);

CREATE TABLE IF NOT EXISTS invoice_numbers (
    id VARCHAR(36) PRIMARY KEY,
    book_id VARCHAR(36) NOT NULL,
    number INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
    invoice_id VARCHAR(36),
    issued_at TIMESTAMPTZ,
    missing_reason VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_in_book ON invoice_numbers(book_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_in_book_number ON invoice_numbers(book_id, number);
