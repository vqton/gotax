CREATE TABLE IF NOT EXISTS tax_declarations (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    declaration_type VARCHAR(20) NOT NULL,
    tax_period_period_year INTEGER NOT NULL,
    tax_period_period_number INTEGER NOT NULL,
    tax_period_period_type VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    adjustment_type VARCHAR(20) NOT NULL DEFAULT 'REGULAR',
    version INTEGER NOT NULL DEFAULT 1,
    submitted_at TIMESTAMPTZ,
    submitted_by VARCHAR(36),
    acknowledged_at TIMESTAMPTZ,
    acknowledgement_ref VARCHAR(100),
    declaration_xml TEXT,
    gdt_response_xml TEXT,
    previous_declaration_id VARCHAR(36),
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_taxdecl_company_type ON tax_declarations(company_id, declaration_type);
CREATE INDEX IF NOT EXISTS idx_taxdecl_status ON tax_declarations(status);

CREATE TABLE IF NOT EXISTS tax_declaration_lines (
    id VARCHAR(36) PRIMARY KEY,
    declaration_id VARCHAR(36) NOT NULL REFERENCES tax_declarations(id) ON DELETE CASCADE,
    line_code VARCHAR(20) NOT NULL,
    line_name VARCHAR(255) NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    source_type VARCHAR(30) NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_taxdecl_line ON tax_declaration_lines(declaration_id, line_code);

CREATE TABLE IF NOT EXISTS tax_rates (
    id VARCHAR(36) PRIMARY KEY,
    tax_type VARCHAR(20) NOT NULL,
    rate_code VARCHAR(20) NOT NULL,
    rate_name VARCHAR(255) NOT NULL,
    rate_type VARCHAR(30) NOT NULL,
    rate_value DOUBLE PRECISION,
    effective_from DATE NOT NULL,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    legal_ref VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_taxrate_code ON tax_rates(rate_code);
CREATE INDEX IF NOT EXISTS idx_taxrate_type ON tax_rates(tax_type);
CREATE INDEX IF NOT EXISTS idx_taxrate_active ON tax_rates(is_active);

CREATE TABLE IF NOT EXISTS tax_payments (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    declaration_id VARCHAR(36) REFERENCES tax_declarations(id),
    tax_type VARCHAR(20) NOT NULL,
    period_year INTEGER NOT NULL,
    period_number INTEGER NOT NULL,
    declared_amount DOUBLE PRECISION NOT NULL,
    paid_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    payment_date DATE,
    due_date DATE NOT NULL,
    payment_ref VARCHAR(100),
    payment_method VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    late_days INTEGER NOT NULL DEFAULT 0,
    late_interest DOUBLE PRECISION NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_taxpay_company ON tax_payments(company_id);
CREATE INDEX IF NOT EXISTS idx_taxpay_declaration ON tax_payments(declaration_id);
CREATE INDEX IF NOT EXISTS idx_taxpay_status ON tax_payments(status);

CREATE TABLE IF NOT EXISTS e_invoices (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    pattern VARCHAR(10) NOT NULL,
    serial VARCHAR(20) NOT NULL,
    invoice_number INTEGER NOT NULL DEFAULT 0,
    invoice_type VARCHAR(10) NOT NULL,
    buyer_name VARCHAR(255) NOT NULL,
    buyer_tax_code VARCHAR(20),
    buyer_address TEXT,
    buyer_email VARCHAR(255),
    currency_code VARCHAR(3) NOT NULL DEFAULT 'VND',
    exchange_rate DOUBLE PRECISION NOT NULL DEFAULT 1,
    subtotal DOUBLE PRECISION NOT NULL,
    vat_amount DOUBLE PRECISION NOT NULL,
    grand_total DOUBLE PRECISION NOT NULL,
    issue_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    xml_body TEXT,
    signed_xml TEXT,
    signing_date VARCHAR(30),
    digital_signature_id VARCHAR(36),
    journal_entry_id VARCHAR(36),
    cancelled_at VARCHAR(30),
    cancel_reason TEXT,
    original_invoice_id VARCHAR(36),
    gdt_response TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_einv_status ON e_invoices(status);
CREATE INDEX IF NOT EXISTS idx_einv_journal ON e_invoices(journal_entry_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_einv_company ON e_invoices(company_id, pattern);

CREATE TABLE IF NOT EXISTS e_invoice_lines (
    id BIGSERIAL PRIMARY KEY,
    e_invoice_id VARCHAR(36) NOT NULL REFERENCES e_invoices(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    unit VARCHAR(20),
    quantity DOUBLE PRECISION NOT NULL,
    unit_price DOUBLE PRECISION NOT NULL,
    line_total DOUBLE PRECISION NOT NULL,
    vat_rate DOUBLE PRECISION NOT NULL,
    vat_amount DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_einvline_inv ON e_invoice_lines(e_invoice_id, line_number);

CREATE TABLE IF NOT EXISTS tax_calendar (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    tax_type VARCHAR(20) NOT NULL,
    period_year INTEGER NOT NULL,
    period_number INTEGER NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_taxcal_company_period ON tax_calendar(company_id, tax_type, period_year, period_number);
CREATE INDEX IF NOT EXISTS idx_taxcal_status ON tax_calendar(status);

CREATE TABLE IF NOT EXISTS tax_alerts (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    alert_type VARCHAR(30) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    due_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_taxalert_company ON tax_alerts(company_id);
CREATE INDEX IF NOT EXISTS idx_taxalert_type ON tax_alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_taxalert_read ON tax_alerts(is_read);

CREATE TABLE IF NOT EXISTS tax_audit_cases (
    id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    case_number VARCHAR(30) NOT NULL,
    audit_type VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    open_date DATE NOT NULL,
    close_date DATE,
    auditor_name VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_taxaudit_case ON tax_audit_cases(case_number);
CREATE INDEX IF NOT EXISTS idx_taxaudit_company ON tax_audit_cases(company_id);
CREATE INDEX IF NOT EXISTS idx_taxaudit_status ON tax_audit_cases(status);
