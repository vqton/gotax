CREATE TABLE IF NOT EXISTS bank_statements (
    id TEXT PRIMARY KEY,
    company_id TEXT NOT NULL,
    bank_account_id TEXT NOT NULL,
    statement_date TEXT NOT NULL,
    from_date TEXT NOT NULL,
    to_date TEXT NOT NULL,
    opening_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    closing_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_credits NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_debits NUMERIC(18,2) NOT NULL DEFAULT 0,
    line_count INTEGER NOT NULL DEFAULT 0,
    currency TEXT NOT NULL DEFAULT 'VND',
    status TEXT NOT NULL DEFAULT 'IMPORTED',
    import_method TEXT NOT NULL DEFAULT '',
    raw_file_name TEXT NOT NULL DEFAULT '',
    raw_file_hash TEXT NOT NULL DEFAULT '',
    imported_by TEXT NOT NULL DEFAULT '',
    imported_at TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_bank_statements_company ON bank_statements(company_id);
CREATE INDEX IF NOT EXISTS idx_bank_statements_account ON bank_statements(company_id, bank_account_id);

CREATE TABLE IF NOT EXISTS bank_statement_lines (
    id TEXT PRIMARY KEY,
    statement_id TEXT NOT NULL REFERENCES bank_statements(id),
    transaction_date TEXT NOT NULL,
    value_date TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    debit_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    credit_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    balance_after NUMERIC(18,2) NOT NULL DEFAULT 0,
    reference_no TEXT NOT NULL DEFAULT '',
    bank_ref TEXT NOT NULL DEFAULT '',
    counterparty TEXT NOT NULL DEFAULT '',
    counterparty_acc TEXT NOT NULL DEFAULT '',
    counterparty_bank TEXT NOT NULL DEFAULT '',
    raw_data TEXT NOT NULL DEFAULT '',
    match_status TEXT NOT NULL DEFAULT 'PENDING',
    matched_line_id TEXT NOT NULL DEFAULT '',
    matched_at TEXT NOT NULL DEFAULT '',
    matched_by TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_bank_lines_statement ON bank_statement_lines(statement_id);
CREATE INDEX IF NOT EXISTS idx_bank_lines_match ON bank_statement_lines(match_status);

CREATE TABLE IF NOT EXISTS bank_reconciliations (
    id TEXT PRIMARY KEY,
    company_id TEXT NOT NULL,
    bank_account_id TEXT NOT NULL,
    statement_id TEXT NOT NULL DEFAULT '',
    from_date TEXT NOT NULL,
    to_date TEXT NOT NULL,
    opening_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    closing_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    statement_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    difference NUMERIC(18,2) NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'IN_PROGRESS',
    matched_lines INTEGER NOT NULL DEFAULT 0,
    unmatched_lines INTEGER NOT NULL DEFAULT 0,
    write_off_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    completed_by TEXT NOT NULL DEFAULT '',
    completed_at TEXT NOT NULL DEFAULT '',
    reversed_at TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_recon_company ON bank_reconciliations(company_id, bank_account_id);

CREATE TABLE IF NOT EXISTS bank_reconciliation_matches (
    id TEXT PRIMARY KEY,
    reconciliation_id TEXT NOT NULL REFERENCES bank_reconciliations(id),
    statement_line_id TEXT NOT NULL DEFAULT '',
    transaction_type TEXT NOT NULL DEFAULT '',
    transaction_id TEXT NOT NULL DEFAULT '',
    transaction_ref TEXT NOT NULL DEFAULT '',
    match_method TEXT NOT NULL DEFAULT '',
    confidence NUMERIC(5,2) NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_recon_matches_recon ON bank_reconciliation_matches(reconciliation_id);

CREATE TABLE IF NOT EXISTS payment_orders (
    id TEXT PRIMARY KEY,
    company_id TEXT NOT NULL,
    payment_date TEXT NOT NULL,
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    exchange_rate NUMERIC(18,4) NOT NULL DEFAULT 1,
    beneficiary_name TEXT NOT NULL DEFAULT '',
    beneficiary_acc TEXT NOT NULL DEFAULT '',
    beneficiary_bank TEXT NOT NULL DEFAULT '',
    beneficiary_branch TEXT NOT NULL DEFAULT '',
    beneficiary_code TEXT NOT NULL DEFAULT '',
    from_bank_acc_id TEXT NOT NULL DEFAULT '',
    payment_content TEXT NOT NULL DEFAULT '',
    urgent INTEGER NOT NULL DEFAULT 0,
    payment_type TEXT NOT NULL DEFAULT 'OTHER',
    status TEXT NOT NULL DEFAULT 'DRAFT',
    created_by TEXT NOT NULL DEFAULT '',
    approved_by TEXT NOT NULL DEFAULT '',
    approved_at TEXT NOT NULL DEFAULT '',
    submitted_at TEXT NOT NULL DEFAULT '',
    bank_ref TEXT NOT NULL DEFAULT '',
    failure_reason TEXT NOT NULL DEFAULT '',
    error_code TEXT NOT NULL DEFAULT '',
    print_count INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_payment_orders_company ON payment_orders(company_id);
CREATE INDEX IF NOT EXISTS idx_payment_orders_status ON payment_orders(status);
CREATE INDEX IF NOT EXISTS idx_payment_orders_date ON payment_orders(payment_date);

CREATE TABLE IF NOT EXISTS payment_order_batches (
    id TEXT PRIMARY KEY,
    company_id TEXT NOT NULL,
    batch_name TEXT NOT NULL DEFAULT '',
    batch_date TEXT NOT NULL DEFAULT '',
    total_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    order_count INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'DRAFT',
    created_by TEXT NOT NULL DEFAULT '',
    submitted_at TEXT NOT NULL DEFAULT '',
    bank_ref TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_batches_company ON payment_order_batches(company_id);

CREATE TABLE IF NOT EXISTS payment_order_batch_items (
    id TEXT PRIMARY KEY,
    batch_id TEXT NOT NULL REFERENCES payment_order_batches(id),
    order_id TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_batch_items_batch ON payment_order_batch_items(batch_id);

CREATE TABLE IF NOT EXISTS loan_agreements (
    id TEXT PRIMARY KEY,
    company_id TEXT NOT NULL,
    bank_account_id TEXT NOT NULL DEFAULT '',
    contract_no TEXT NOT NULL DEFAULT '',
    loan_type TEXT NOT NULL DEFAULT 'SHORT_TERM',
    principal_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    interest_rate NUMERIC(5,4) NOT NULL DEFAULT 0,
    interest_method TEXT NOT NULL DEFAULT 'FIXED',
    base_rate NUMERIC(5,4) NOT NULL DEFAULT 0,
    margin_rate NUMERIC(5,4) NOT NULL DEFAULT 0,
    disbursed_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    outstanding_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    start_date TEXT NOT NULL DEFAULT '',
    maturity_date TEXT NOT NULL DEFAULT '',
    repayment_method TEXT NOT NULL DEFAULT 'ANNUITY',
    repayment_freq TEXT NOT NULL DEFAULT 'MONTHLY',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_loans_company ON loan_agreements(company_id);

CREATE TABLE IF NOT EXISTS loan_disbursements (
    id TEXT PRIMARY KEY,
    loan_id TEXT NOT NULL REFERENCES loan_agreements(id),
    disbursement_date TEXT NOT NULL DEFAULT '',
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    to_bank_account_id TEXT NOT NULL DEFAULT '',
    reference_no TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_disbursements_loan ON loan_disbursements(loan_id);

CREATE TABLE IF NOT EXISTS loan_repayments (
    id TEXT PRIMARY KEY,
    loan_id TEXT NOT NULL REFERENCES loan_agreements(id),
    repayment_date TEXT NOT NULL DEFAULT '',
    principal_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    interest_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    fee_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    payment_order_id TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'SCHEDULED',
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_repayments_loan ON loan_repayments(loan_id);

CREATE TABLE IF NOT EXISTS term_deposits (
    id TEXT PRIMARY KEY,
    company_id TEXT NOT NULL,
    bank_account_id TEXT NOT NULL DEFAULT '',
    deposit_no TEXT NOT NULL DEFAULT '',
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    interest_rate NUMERIC(5,4) NOT NULL DEFAULT 0,
    term_days INTEGER NOT NULL DEFAULT 0,
    start_date TEXT NOT NULL DEFAULT '',
    maturity_date TEXT NOT NULL DEFAULT '',
    interest_at_maturity NUMERIC(18,2) NOT NULL DEFAULT 0,
    auto_renewal INTEGER NOT NULL DEFAULT 0,
    renewal_term_days INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT '',
    matured_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_deposits_company ON term_deposits(company_id);
