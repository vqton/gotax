-- O2C/Sale domain schema — 13 tables
-- Circular 99/2025 + Decree 123/2020. Multi-tenant, multi-currency.

-- ─── Customer ──────────────────────────────────────────────────────────

CREATE TABLE customers (
    id               VARCHAR(36) PRIMARY KEY,
    company_id       VARCHAR(36) NOT NULL,
    code             VARCHAR(20) NOT NULL,
    name             VARCHAR(255) NOT NULL,
    tax_code         VARCHAR(20) NOT NULL,
    address          TEXT,
    phone            VARCHAR(20),
    email            VARCHAR(255),
    bank_account_name   VARCHAR(255),
    bank_account_number VARCHAR(50),
    bank_name        VARCHAR(255),
    payment_terms    VARCHAR(50),
    credit_limit     NUMERIC NOT NULL DEFAULT 0,
    currency         VARCHAR(3) DEFAULT 'VND',
    customer_type    VARCHAR(20) DEFAULT 'domestic',
    customer_group   VARCHAR(20),
    price_list_id    VARCHAR(36),
    status           VARCHAR(20) DEFAULT 'ACTIVE',
    notes            TEXT,
    created_by       VARCHAR(36),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (company_id, code)
);
CREATE INDEX idx_customer_company_code ON customers (company_id, code);
CREATE INDEX idx_customer_status ON customers (status);

-- ─── Sales Quotation ───────────────────────────────────────────────────

CREATE TABLE sales_quotations (
    id              VARCHAR(36) PRIMARY KEY,
    company_id      VARCHAR(36) NOT NULL,
    qn_number       VARCHAR(30) NOT NULL,
    customer_id     VARCHAR(36),
    valid_until     DATE,
    status          VARCHAR(20) DEFAULT 'DRAFT',
    total_amount    NUMERIC NOT NULL DEFAULT 0,
    created_by      VARCHAR(36),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (company_id, qn_number)
);
CREATE INDEX idx_sq_company ON sales_quotations (company_id);

-- ─── Sales Order ───────────────────────────────────────────────────────

CREATE TABLE sales_orders (
    id               VARCHAR(36) PRIMARY KEY,
    company_id       VARCHAR(36) NOT NULL,
    so_number        VARCHAR(30) NOT NULL,
    quotation_id     VARCHAR(36),
    customer_id      VARCHAR(36) NOT NULL,
    order_date       DATE NOT NULL,
    expected_date    DATE,
    currency         VARCHAR(3) DEFAULT 'VND',
    exchange_rate    NUMERIC NOT NULL DEFAULT 1,
    payment_terms    VARCHAR(50),
    delivery_terms   VARCHAR(50),
    shipping_address TEXT,
    subtotal         NUMERIC NOT NULL,
    discount_amount  NUMERIC NOT NULL DEFAULT 0,
    tax_amount       NUMERIC NOT NULL DEFAULT 0,
    total_amount     NUMERIC NOT NULL,
    status           VARCHAR(20) DEFAULT 'DRAFT',
    approved_by      VARCHAR(36),
    approved_at      TIMESTAMPTZ,
    cancelled_reason TEXT,
    notes            TEXT,
    created_by       VARCHAR(36),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (company_id, so_number)
);
CREATE INDEX idx_so_company_number ON sales_orders (company_id, so_number);
CREATE INDEX idx_so_customer ON sales_orders (customer_id);

-- ─── SO Lines ──────────────────────────────────────────────────────────

CREATE TABLE so_lines (
    id                VARCHAR(36) PRIMARY KEY,
    so_id             VARCHAR(36) NOT NULL REFERENCES sales_orders(id) ON DELETE CASCADE,
    line_number       INT NOT NULL,
    item_code         VARCHAR(50),
    item_name         VARCHAR(255) NOT NULL,
    unit              VARCHAR(20) NOT NULL,
    quantity          NUMERIC NOT NULL,
    unit_price        NUMERIC NOT NULL,
    discount_pct      NUMERIC NOT NULL DEFAULT 0,
    vat_rate          NUMERIC NOT NULL DEFAULT 0,
    vat_type          VARCHAR(10) DEFAULT 'VAT10',
    revenue_account_id VARCHAR(20) NOT NULL,
    vat_account_id    VARCHAR(20) NOT NULL,
    line_total        NUMERIC NOT NULL,
    line_vat_amount   NUMERIC NOT NULL DEFAULT 0,
    delivered_qty     NUMERIC NOT NULL DEFAULT 0,
    invoiced_qty      NUMERIC NOT NULL DEFAULT 0
);
CREATE INDEX idx_soline_so ON so_lines (so_id);

-- ─── Delivery Note ─────────────────────────────────────────────────────

CREATE TABLE delivery_notes (
    id               VARCHAR(36) PRIMARY KEY,
    company_id       VARCHAR(36) NOT NULL,
    dn_number        VARCHAR(30) NOT NULL,
    so_id            VARCHAR(36) NOT NULL,
    delivery_date    DATE NOT NULL,
    warehouse        VARCHAR(50),
    shipping_method  VARCHAR(50),
    carrier_name     VARCHAR(100),
    tracking_number  VARCHAR(100),
    delivery_address TEXT,
    status           VARCHAR(20) DEFAULT 'DRAFT',
    notes            TEXT,
    tolerance_percent NUMERIC NOT NULL DEFAULT 5,
    created_by       VARCHAR(36),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (company_id, dn_number)
);
CREATE INDEX idx_dn_company_number ON delivery_notes (company_id, dn_number);
CREATE INDEX idx_dn_so ON delivery_notes (so_id);

-- ─── DN Lines ──────────────────────────────────────────────────────────

CREATE TABLE dn_lines (
    id            VARCHAR(36) PRIMARY KEY,
    dn_id         VARCHAR(36) NOT NULL REFERENCES delivery_notes(id) ON DELETE CASCADE,
    so_line_id    VARCHAR(36),
    item_code     VARCHAR(50),
    item_name     VARCHAR(255) NOT NULL,
    unit          VARCHAR(20),
    qty_delivered NUMERIC NOT NULL,
    qty_returned  NUMERIC NOT NULL DEFAULT 0,
    unit_price    NUMERIC NOT NULL,
    line_total    NUMERIC NOT NULL,
    cost_price    NUMERIC NOT NULL DEFAULT 0
);
CREATE INDEX idx_dnline_dn ON dn_lines (dn_id);

-- ─── Customer Invoice ──────────────────────────────────────────────────

CREATE TABLE customer_invoices (
    id                   VARCHAR(36) PRIMARY KEY,
    company_id           VARCHAR(36) NOT NULL,
    invoice_number       VARCHAR(30) NOT NULL,
    invoice_date         DATE NOT NULL,
    so_id                VARCHAR(36),
    dn_id                VARCHAR(36),
    customer_id          VARCHAR(36) NOT NULL,
    customer_name        VARCHAR(255) NOT NULL,
    customer_tax_code    VARCHAR(20) NOT NULL,
    customer_address     TEXT,
    invoice_type         VARCHAR(20),
    currency             VARCHAR(3) DEFAULT 'VND',
    exchange_rate        NUMERIC NOT NULL DEFAULT 1,
    subtotal             NUMERIC NOT NULL,
    discount_amount      NUMERIC NOT NULL DEFAULT 0,
    tax_amount           NUMERIC NOT NULL DEFAULT 0,
    total_amount         NUMERIC NOT NULL,
    amount_received      NUMERIC NOT NULL DEFAULT 0,
    balance_due          NUMERIC NOT NULL DEFAULT 0,
    due_date             DATE,
    invoice_note         TEXT,
    e_invoice_data       TEXT,
    e_invoice_code       VARCHAR(50),
    e_invoice_status     VARCHAR(20) DEFAULT 'pending',
    digital_signature_id VARCHAR(36),
    signed_data          TEXT,
    gdt_response         TEXT,
    original_invoice_id  VARCHAR(36),
    adjustment_type      VARCHAR(20),
    status               VARCHAR(20) DEFAULT 'DRAFT',
    gl_posted            BOOLEAN NOT NULL DEFAULT FALSE,
    gl_posted_at         TIMESTAMPTZ,
    notes                TEXT,
    created_by           VARCHAR(36),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (company_id, invoice_number)
);
CREATE INDEX idx_cinv_company_number ON customer_invoices (company_id, invoice_number);
CREATE INDEX idx_cinv_customer ON customer_invoices (customer_id);

-- ─── Invoice Lines ─────────────────────────────────────────────────────

CREATE TABLE inv_lines (
    id                 VARCHAR(36) PRIMARY KEY,
    invoice_id         VARCHAR(36) NOT NULL REFERENCES customer_invoices(id) ON DELETE CASCADE,
    so_line_id         VARCHAR(36),
    dn_line_id         VARCHAR(36),
    item_code          VARCHAR(50),
    item_name          VARCHAR(255) NOT NULL,
    unit               VARCHAR(20),
    quantity           NUMERIC NOT NULL,
    unit_price         NUMERIC NOT NULL,
    discount_pct       NUMERIC NOT NULL DEFAULT 0,
    vat_rate           NUMERIC NOT NULL DEFAULT 0,
    vat_type           VARCHAR(10) DEFAULT 'VAT10',
    line_total         NUMERIC NOT NULL,
    line_vat_amount    NUMERIC NOT NULL DEFAULT 0,
    revenue_account_id VARCHAR(20) NOT NULL,
    vat_account_id     VARCHAR(20) NOT NULL
);
CREATE INDEX idx_invline_inv ON inv_lines (invoice_id);

-- ─── Customer Receipt ──────────────────────────────────────────────────

CREATE TABLE customer_receipts (
    id                VARCHAR(36) PRIMARY KEY,
    company_id        VARCHAR(36) NOT NULL,
    receipt_number    VARCHAR(30) NOT NULL,
    customer_id       VARCHAR(36) NOT NULL,
    receipt_date      DATE NOT NULL,
    payment_method    VARCHAR(30),
    bank_account_id   VARCHAR(36),
    currency          VARCHAR(3) DEFAULT 'VND',
    exchange_rate     NUMERIC NOT NULL DEFAULT 1,
    amount            NUMERIC NOT NULL,
    unallocated_amount NUMERIC NOT NULL DEFAULT 0,
    reference         VARCHAR(100),
    notes             TEXT,
    status            VARCHAR(20) DEFAULT 'DRAFT',
    gl_posted         BOOLEAN NOT NULL DEFAULT FALSE,
    gl_posted_at      TIMESTAMPTZ,
    created_by        VARCHAR(36),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (company_id, receipt_number)
);
CREATE INDEX idx_rcpt_company_number ON customer_receipts (company_id, receipt_number);
CREATE INDEX idx_rcpt_customer ON customer_receipts (customer_id);

-- ─── Receipt Allocations ───────────────────────────────────────────────

CREATE TABLE rcp_allocations (
    id               VARCHAR(36) PRIMARY KEY,
    receipt_id       VARCHAR(36) NOT NULL REFERENCES customer_receipts(id) ON DELETE CASCADE,
    invoice_id       VARCHAR(36) NOT NULL,
    allocated_amount NUMERIC NOT NULL,
    discount_amount  NUMERIC NOT NULL DEFAULT 0
);
CREATE INDEX idx_rcp_alloc_rcpt ON rcp_allocations (receipt_id);

-- ─── Credit Note ───────────────────────────────────────────────────────

CREATE TABLE credit_notes (
    id                 VARCHAR(36) PRIMARY KEY,
    company_id         VARCHAR(36) NOT NULL,
    cn_number          VARCHAR(30) NOT NULL,
    original_invoice_id VARCHAR(36) NOT NULL,
    customer_id        VARCHAR(36) NOT NULL,
    return_date        DATE NOT NULL,
    return_reason      TEXT,
    return_type        VARCHAR(20) DEFAULT 'partial_return',
    dn_id              VARCHAR(36),
    subtotal           NUMERIC NOT NULL,
    tax_amount         NUMERIC NOT NULL DEFAULT 0,
    total_amount       NUMERIC NOT NULL,
    e_invoice_data     TEXT,
    e_invoice_code     VARCHAR(50),
    status             VARCHAR(20) DEFAULT 'DRAFT',
    gl_posted          BOOLEAN NOT NULL DEFAULT FALSE,
    gl_posted_at       TIMESTAMPTZ,
    notes              TEXT,
    created_by         VARCHAR(36),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (company_id, cn_number)
);
CREATE INDEX idx_cn_company_number ON credit_notes (company_id, cn_number);
CREATE INDEX idx_cn_customer ON credit_notes (customer_id);
CREATE INDEX idx_cn_original_invoice ON credit_notes (original_invoice_id);

-- ─── CN Lines ──────────────────────────────────────────────────────────

CREATE TABLE cn_lines (
    id              VARCHAR(36) PRIMARY KEY,
    cn_id           VARCHAR(36) NOT NULL REFERENCES credit_notes(id) ON DELETE CASCADE,
    invoice_line_id VARCHAR(36),
    item_name       VARCHAR(255) NOT NULL,
    unit            VARCHAR(20),
    quantity        NUMERIC NOT NULL,
    unit_price      NUMERIC NOT NULL,
    vat_rate        NUMERIC NOT NULL DEFAULT 0,
    line_total      NUMERIC NOT NULL,
    line_vat_amount NUMERIC NOT NULL DEFAULT 0
);
CREATE INDEX idx_cnline_cn ON cn_lines (cn_id);

-- ─── AR Transactions ───────────────────────────────────────────────────

CREATE TABLE ar_transactions (
    id                VARCHAR(36) PRIMARY KEY,
    company_id        VARCHAR(36) NOT NULL,
    customer_id       VARCHAR(36) NOT NULL,
    invoice_id        VARCHAR(36),
    transaction_type  VARCHAR(20) NOT NULL,
    transaction_date  DATE NOT NULL,
    amount            NUMERIC NOT NULL,
    currency          VARCHAR(3) DEFAULT 'VND',
    reference_type    VARCHAR(50),
    reference_id      VARCHAR(36),
    notes             TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ar_customer ON ar_transactions (customer_id);
CREATE INDEX idx_ar_invoice ON ar_transactions (invoice_id);
