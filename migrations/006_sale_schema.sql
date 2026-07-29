/*
Sale Module Schema — Order-to-Cash (O2C) — v1
No third-party integrations. GL auto-posting deferred.
Complies: Circular 99/2025, Decree 123/2020, IFRS 15.
*/
CREATE TABLE IF NOT EXISTS customers (
  id                    TEXT    PRIMARY KEY,
  company_id            TEXT    NOT NULL,
  code                  TEXT    NOT NULL,
  name                  TEXT    NOT NULL,
  tax_code              TEXT    NOT NULL DEFAULT '',
  address               TEXT    NOT NULL DEFAULT '',
  phone                 TEXT    NOT NULL DEFAULT '',
  email                 TEXT    NOT NULL DEFAULT '',
  bank_account_name     TEXT    NOT NULL DEFAULT '',
  bank_account_number   TEXT    NOT NULL DEFAULT '',
  bank_name             TEXT    NOT NULL DEFAULT '',
  payment_terms         TEXT    NOT NULL DEFAULT 'net30',
  credit_limit          NUMERIC(18,2)   NOT NULL DEFAULT 0,
  currency              TEXT    NOT NULL DEFAULT 'VND',
  customer_type         TEXT    NOT NULL DEFAULT 'domestic',
  status                TEXT    NOT NULL DEFAULT 'ACTIVE',
  notes                 TEXT    NOT NULL DEFAULT '',
  created_at            TEXT    NOT NULL DEFAULT '',
  updated_at            TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS sales_orders (
  id                  TEXT    PRIMARY KEY,
  company_id          TEXT    NOT NULL,
  so_number           TEXT    NOT NULL,
  quotation_id        TEXT    NOT NULL DEFAULT '',
  customer_id         TEXT    NOT NULL REFERENCES customers(id),
  order_date          TEXT    NOT NULL,
  expected_date       TEXT    NOT NULL DEFAULT '',
  currency            TEXT    NOT NULL DEFAULT 'VND',
  exchange_rate       NUMERIC(18,2)   NOT NULL DEFAULT 1,
  payment_terms       TEXT    NOT NULL DEFAULT '',
  delivery_terms      TEXT    NOT NULL DEFAULT '',
  shipping_address    TEXT    NOT NULL DEFAULT '',
  subtotal            NUMERIC(18,2)   NOT NULL DEFAULT 0,
  discount_amount     NUMERIC(18,2)   NOT NULL DEFAULT 0,
  tax_amount          NUMERIC(18,2)   NOT NULL DEFAULT 0,
  total_amount        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  status              TEXT    NOT NULL DEFAULT 'DRAFT',
  approved_by         TEXT    NOT NULL DEFAULT '',
  approved_at         TEXT    NOT NULL DEFAULT '',
  cancelled_reason    TEXT    NOT NULL DEFAULT '',
  notes               TEXT    NOT NULL DEFAULT '',
  created_by          TEXT    NOT NULL DEFAULT '',
  created_at          TEXT    NOT NULL DEFAULT '',
  updated_at          TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS so_lines (
  id                TEXT    PRIMARY KEY,
  so_id             TEXT    NOT NULL REFERENCES sales_orders(id),
  line_number       INT     NOT NULL DEFAULT 0,
  item_code         TEXT    NOT NULL DEFAULT '',
  item_name         TEXT    NOT NULL,
  unit              TEXT    NOT NULL,
  quantity          NUMERIC(18,2)   NOT NULL DEFAULT 0,
  unit_price        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  discount_pct      NUMERIC(18,2)   NOT NULL DEFAULT 0,
  vat_rate          NUMERIC(18,2)   NOT NULL DEFAULT 0,
  vat_type          TEXT    NOT NULL DEFAULT 'VAT10',
  revenue_account_id TEXT   NOT NULL,
  vat_account_id    TEXT   NOT NULL,
  line_total        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  line_vat_amount   NUMERIC(18,2)   NOT NULL DEFAULT 0,
  delivered_qty     NUMERIC(18,2)   NOT NULL DEFAULT 0,
  invoiced_qty      NUMERIC(18,2)   NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS delivery_notes (
  id                TEXT    PRIMARY KEY,
  company_id        TEXT    NOT NULL,
  dn_number         TEXT    NOT NULL,
  so_id             TEXT    NOT NULL REFERENCES sales_orders(id),
  delivery_date     TEXT    NOT NULL,
  warehouse         TEXT    NOT NULL DEFAULT '',
  shipping_method   TEXT    NOT NULL DEFAULT '',
  carrier_name      TEXT    NOT NULL DEFAULT '',
  tracking_number   TEXT    NOT NULL DEFAULT '',
  delivery_address  TEXT    NOT NULL DEFAULT '',
  status            TEXT    NOT NULL DEFAULT 'DRAFT',
  notes             TEXT    NOT NULL DEFAULT '',
  created_by        TEXT    NOT NULL DEFAULT '',
  created_at        TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS dn_lines (
  id                TEXT    PRIMARY KEY,
  dn_id             TEXT    NOT NULL REFERENCES delivery_notes(id),
  so_line_id        TEXT    NOT NULL DEFAULT '',
  item_code         TEXT    NOT NULL DEFAULT '',
  item_name         TEXT    NOT NULL,
  unit              TEXT    NOT NULL DEFAULT '',
  qty_delivered     NUMERIC(18,2)   NOT NULL DEFAULT 0,
  qty_returned      NUMERIC(18,2)   NOT NULL DEFAULT 0,
  unit_price        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  line_total        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  cost_price        NUMERIC(18,2)   NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS customer_invoices (
  id                    TEXT    PRIMARY KEY,
  company_id            TEXT    NOT NULL,
  invoice_number        TEXT    NOT NULL,
  invoice_date          TEXT    NOT NULL,
  so_id                 TEXT    NOT NULL DEFAULT '',
  dn_id                 TEXT    NOT NULL DEFAULT '',
  customer_id           TEXT    NOT NULL REFERENCES customers(id),
  customer_name         TEXT    NOT NULL,
  customer_tax_code     TEXT    NOT NULL,
  customer_address      TEXT    NOT NULL DEFAULT '',
  invoice_type          TEXT    NOT NULL DEFAULT '',
  currency              TEXT    NOT NULL DEFAULT 'VND',
  exchange_rate         NUMERIC(18,2)   NOT NULL DEFAULT 1,
  subtotal              NUMERIC(18,2)   NOT NULL DEFAULT 0,
  discount_amount       NUMERIC(18,2)   NOT NULL DEFAULT 0,
  tax_amount            NUMERIC(18,2)   NOT NULL DEFAULT 0,
  total_amount          NUMERIC(18,2)   NOT NULL DEFAULT 0,
  amount_received       NUMERIC(18,2)   NOT NULL DEFAULT 0,
  balance_due           NUMERIC(18,2)   NOT NULL DEFAULT 0,
  due_date              TEXT    NOT NULL DEFAULT '',
  invoice_note          TEXT    NOT NULL DEFAULT '',
  e_invoice_data        TEXT    NOT NULL DEFAULT '',
  e_invoice_code        TEXT    NOT NULL DEFAULT '',
  e_invoice_status      TEXT    NOT NULL DEFAULT 'NOT_SENT',
  digital_signature_id  TEXT    NOT NULL DEFAULT '',
  signed_data           TEXT    NOT NULL DEFAULT '',
  gdt_response          TEXT    NOT NULL DEFAULT '',
  original_invoice_id   TEXT    NOT NULL DEFAULT '',
  adjustment_type       TEXT    NOT NULL DEFAULT '',
  status                TEXT    NOT NULL DEFAULT 'DRAFT',
  gl_posted             INT     NOT NULL DEFAULT 0,
  gl_posted_at          TEXT    NOT NULL DEFAULT '',
  notes                 TEXT    NOT NULL DEFAULT '',
  created_by            TEXT    NOT NULL DEFAULT '',
  created_at            TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS inv_lines (
  id                TEXT    PRIMARY KEY,
  invoice_id        TEXT    NOT NULL REFERENCES customer_invoices(id),
  so_line_id        TEXT    NOT NULL DEFAULT '',
  dn_line_id        TEXT    NOT NULL DEFAULT '',
  item_code         TEXT    NOT NULL DEFAULT '',
  item_name         TEXT    NOT NULL,
  unit              TEXT    NOT NULL DEFAULT '',
  quantity          NUMERIC(18,2)   NOT NULL DEFAULT 0,
  unit_price        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  discount_pct      NUMERIC(18,2)   NOT NULL DEFAULT 0,
  vat_rate          NUMERIC(18,2)   NOT NULL DEFAULT 0,
  vat_type          TEXT    NOT NULL DEFAULT 'VAT10',
  line_total        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  line_vat_amount   NUMERIC(18,2)   NOT NULL DEFAULT 0,
  revenue_account_id TEXT   NOT NULL,
  vat_account_id    TEXT   NOT NULL
);

CREATE TABLE IF NOT EXISTS customer_receipts (
  id                  TEXT    PRIMARY KEY,
  company_id          TEXT    NOT NULL,
  receipt_number      TEXT    NOT NULL,
  customer_id         TEXT    NOT NULL REFERENCES customers(id),
  receipt_date        TEXT    NOT NULL,
  payment_method      TEXT    NOT NULL DEFAULT '',
  bank_account_id     TEXT    NOT NULL DEFAULT '',
  currency            TEXT    NOT NULL DEFAULT 'VND',
  exchange_rate       NUMERIC(18,2)   NOT NULL DEFAULT 1,
  amount              NUMERIC(18,2)   NOT NULL DEFAULT 0,
  unallocated_amount  NUMERIC(18,2)   NOT NULL DEFAULT 0,
  reference           TEXT    NOT NULL DEFAULT '',
  notes               TEXT    NOT NULL DEFAULT '',
  status              TEXT    NOT NULL DEFAULT 'DRAFT',
  created_by          TEXT    NOT NULL DEFAULT '',
  created_at          TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS receipt_allocations (
  id                TEXT    PRIMARY KEY,
  receipt_id        TEXT    NOT NULL REFERENCES customer_receipts(id),
  invoice_id        TEXT    NOT NULL REFERENCES customer_invoices(id),
  allocated_amount  NUMERIC(18,2)   NOT NULL DEFAULT 0,
  discount_amount   NUMERIC(18,2)   NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS credit_notes (
  id                  TEXT    PRIMARY KEY,
  company_id          TEXT    NOT NULL,
  cn_number           TEXT    NOT NULL,
  original_invoice_id TEXT    NOT NULL REFERENCES customer_invoices(id),
  customer_id         TEXT    NOT NULL REFERENCES customers(id),
  return_date         TEXT    NOT NULL,
  return_reason       TEXT    NOT NULL DEFAULT '',
  return_type         TEXT    NOT NULL DEFAULT 'FULL',
  dn_id               TEXT    NOT NULL DEFAULT '',
  subtotal            NUMERIC(18,2)   NOT NULL DEFAULT 0,
  tax_amount          NUMERIC(18,2)   NOT NULL DEFAULT 0,
  total_amount        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  e_invoice_data      TEXT    NOT NULL DEFAULT '',
  e_invoice_code      TEXT    NOT NULL DEFAULT '',
  status              TEXT    NOT NULL DEFAULT 'DRAFT',
  gl_posted           INT     NOT NULL DEFAULT 0,
  gl_posted_at        TEXT    NOT NULL DEFAULT '',
  notes               TEXT    NOT NULL DEFAULT '',
  created_by          TEXT    NOT NULL DEFAULT '',
  created_at          TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS cn_lines (
  id              TEXT    PRIMARY KEY,
  cn_id           TEXT    NOT NULL REFERENCES credit_notes(id),
  invoice_line_id TEXT    NOT NULL DEFAULT '',
  item_name       TEXT    NOT NULL,
  unit            TEXT    NOT NULL DEFAULT '',
  quantity        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  unit_price      NUMERIC(18,2)   NOT NULL DEFAULT 0,
  vat_rate        NUMERIC(18,2)   NOT NULL DEFAULT 0,
  line_total      NUMERIC(18,2)   NOT NULL DEFAULT 0,
  line_vat_amount NUMERIC(18,2)   NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS ar_transactions (
  id                TEXT    PRIMARY KEY,
  company_id        TEXT    NOT NULL,
  customer_id       TEXT    NOT NULL REFERENCES customers(id),
  invoice_id        TEXT    NOT NULL DEFAULT '',
  transaction_type  TEXT    NOT NULL,
  transaction_date  TEXT    NOT NULL,
  amount            NUMERIC(18,2)   NOT NULL DEFAULT 0,
  currency          TEXT    NOT NULL DEFAULT 'VND',
  reference_type    TEXT    NOT NULL DEFAULT '',
  reference_id      TEXT    NOT NULL DEFAULT '',
  notes             TEXT    NOT NULL DEFAULT '',
  created_at        TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_customers_company ON customers(company_id);
CREATE INDEX IF NOT EXISTS idx_so_company ON sales_orders(company_id);
CREATE INDEX IF NOT EXISTS idx_so_customer ON sales_orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_so_status ON sales_orders(status);
CREATE INDEX IF NOT EXISTS idx_dn_company ON delivery_notes(company_id);
CREATE INDEX IF NOT EXISTS idx_dn_so ON delivery_notes(so_id);
CREATE INDEX IF NOT EXISTS idx_inv_company ON customer_invoices(company_id);
CREATE INDEX IF NOT EXISTS idx_inv_customer ON customer_invoices(customer_id);
CREATE INDEX IF NOT EXISTS idx_inv_status ON customer_invoices(status);
CREATE INDEX IF NOT EXISTS idx_rcpt_company ON customer_receipts(company_id);
CREATE INDEX IF NOT EXISTS idx_cn_company ON credit_notes(company_id);

CREATE TABLE IF NOT EXISTS sales_quotations (
  id                TEXT    PRIMARY KEY,
  company_id        TEXT    NOT NULL,
  qn_number         TEXT    NOT NULL,
  customer_id       TEXT    NOT NULL DEFAULT '',
  valid_until       TEXT    NOT NULL DEFAULT '',
  status            TEXT    NOT NULL DEFAULT 'DRAFT',
  total_amount      NUMERIC(18,2)   NOT NULL DEFAULT 0,
  created_by        TEXT    NOT NULL DEFAULT '',
  created_at        TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_art_company ON ar_transactions(company_id);
CREATE INDEX IF NOT EXISTS idx_art_customer ON ar_transactions(customer_id);
