-- Cash Module Schema (Circular 99/2025/TT-BTC compliant)
-- Effective: 2026-07-28

CREATE TABLE IF NOT EXISTS cash_receipts (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL,
    voucher_no VARCHAR(20) NOT NULL,
    voucher_date DATE NOT NULL,
    cash_account_id VARCHAR(20) NOT NULL,
    counterpart_id VARCHAR(50),
    counterpart_name VARCHAR(255),
    counterpart_type VARCHAR(20) NOT NULL DEFAULT 'OTHER',
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    exchange_rate NUMERIC(18,4) NOT NULL DEFAULT 1,
    amount NUMERIC(18,2) NOT NULL,
    amount_vnd NUMERIC(18,2) NOT NULL,
    debit_account_id VARCHAR(20) NOT NULL,
    credit_account_id VARCHAR(20) NOT NULL,
    reason TEXT,
    receipt_type VARCHAR(30) NOT NULL DEFAULT 'OTHER',
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    approved_by VARCHAR(50),
    approved_at TIMESTAMPTZ,
    posted_by VARCHAR(50),
    posted_at TIMESTAMPTZ,
    gl_journal_id VARCHAR(50),
    created_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, voucher_no)
);

CREATE TABLE IF NOT EXISTS cash_payments (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL,
    voucher_no VARCHAR(20) NOT NULL,
    voucher_date DATE NOT NULL,
    cash_account_id VARCHAR(20) NOT NULL,
    payee_id VARCHAR(50),
    payee_name VARCHAR(255),
    payee_type VARCHAR(20) NOT NULL DEFAULT 'OTHER',
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    exchange_rate NUMERIC(18,4) NOT NULL DEFAULT 1,
    amount NUMERIC(18,2) NOT NULL,
    amount_vnd NUMERIC(18,2) NOT NULL,
    debit_account_id VARCHAR(20) NOT NULL,
    credit_account_id VARCHAR(20) NOT NULL,
    reason TEXT,
    payment_type VARCHAR(30) NOT NULL DEFAULT 'OTHER',
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    approved_by VARCHAR(50),
    approved_at TIMESTAMPTZ,
    posted_by VARCHAR(50),
    posted_at TIMESTAMPTZ,
    gl_journal_id VARCHAR(50),
    created_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, voucher_no)
);

CREATE TABLE IF NOT EXISTS cash_transfers (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL,
    transfer_date DATE NOT NULL,
    from_account_id VARCHAR(20) NOT NULL,
    to_account_id VARCHAR(20) NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    exchange_rate NUMERIC(18,4) NOT NULL DEFAULT 1,
    reason TEXT,
    transfer_type VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    source_voucher_id VARCHAR(50),
    dest_voucher_id VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    posted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS petty_cash_funds (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL,
    fund_code VARCHAR(20) NOT NULL,
    fund_name VARCHAR(255) NOT NULL,
    custodian_id VARCHAR(50) NOT NULL,
    initial_amount NUMERIC(18,2) NOT NULL,
    current_balance NUMERIC(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, fund_code)
);

CREATE TABLE IF NOT EXISTS cash_inventories (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL,
    inventory_date DATE NOT NULL,
    cash_account_id VARCHAR(20) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    book_balance NUMERIC(18,2) NOT NULL,
    actual_balance NUMERIC(18,2) NOT NULL,
    difference NUMERIC(18,2) NOT NULL DEFAULT 0,
    difference_type VARCHAR(10) NOT NULL DEFAULT 'none',
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    approved_by VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cash_inventory_details (
    id VARCHAR(50) PRIMARY KEY,
    inventory_id VARCHAR(50) NOT NULL REFERENCES cash_inventories(id),
    denomination NUMERIC(18,2) NOT NULL,
    count INT NOT NULL DEFAULT 0,
    subtotal NUMERIC(18,2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_cash_receipts_company ON cash_receipts(company_id);
CREATE INDEX IF NOT EXISTS idx_cash_receipts_date ON cash_receipts(voucher_date);
CREATE INDEX IF NOT EXISTS idx_cash_receipts_status ON cash_receipts(status);
CREATE INDEX IF NOT EXISTS idx_cash_payments_company ON cash_payments(company_id);
CREATE INDEX IF NOT EXISTS idx_cash_payments_date ON cash_payments(voucher_date);
CREATE INDEX IF NOT EXISTS idx_cash_payments_status ON cash_payments(status);
CREATE INDEX IF NOT EXISTS idx_cash_transfers_company ON cash_transfers(company_id);
CREATE INDEX IF NOT EXISTS idx_petty_cash_company ON petty_cash_funds(company_id);
CREATE INDEX IF NOT EXISTS idx_cash_inventory_company ON cash_inventories(company_id);
CREATE INDEX IF NOT EXISTS idx_cash_inventory_details_inv ON cash_inventory_details(inventory_id);
