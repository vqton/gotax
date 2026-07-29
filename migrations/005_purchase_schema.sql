/*
Purchase Module Schema — Procure-to-Pay (P2P) — v1
No third-party integrations. GL auto-posting deferred to Phase 3.
Complies: Circular 99/2025, Decree 123/2020, Decree 254/2026.
*/
CREATE TABLE IF NOT EXISTS suppliers (
  id              TEXT    PRIMARY KEY,
  company_id      TEXT    NOT NULL,
  code            TEXT    NOT NULL,
  name            TEXT    NOT NULL,
  tax_code        TEXT    NOT NULL,
  address         TEXT    NOT NULL DEFAULT '',
  phone           TEXT    NOT NULL DEFAULT '',
  email           TEXT    NOT NULL DEFAULT '',
  bank_account_name  TEXT    NOT NULL DEFAULT '',
  bank_account_number TEXT    NOT NULL DEFAULT '',
  bank_name        TEXT    NOT NULL DEFAULT '',
  payment_terms    TEXT    NOT NULL DEFAULT 'net30',
  credit_limit     REAL   NOT NULL DEFAULT 0,
  currency         TEXT    NOT NULL DEFAULT 'VND',
  supplier_type    TEXT    NOT NULL DEFAULT 'domestic',
  status           TEXT    NOT NULL DEFAULT 'ACTIVE',
  notes            TEXT    NOT NULL DEFAULT '',
  created_at       TEXT    NOT NULL DEFAULT '',
  updated_at       TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS purchase_orders (
  id              TEXT    PRIMARY KEY,
  company_id      TEXT    NOT NULL,
  po_number       TEXT    NOT NULL,
  supplier_id     TEXT    NOT NULL REFERENCES suppliers(id),
  requisition_id  TEXT    NOT NULL DEFAULT '',
  order_date      TEXT    NOT NULL,
  expected_date   TEXT    NOT NULL DEFAULT '',
  currency        TEXT    NOT NULL DEFAULT 'VND',
  exchange_rate   REAL   NOT NULL DEFAULT 1,
  payment_terms   TEXT    NOT NULL DEFAULT '',
  delivery_terms  TEXT    NOT NULL DEFAULT '',
  subtotal        REAL   NOT NULL DEFAULT 0,
  discount_amount REAL   NOT NULL DEFAULT 0,
  tax_amount      REAL   NOT NULL DEFAULT 0,
  total_amount    REAL   NOT NULL DEFAULT 0,
  status          TEXT    NOT NULL DEFAULT 'DRAFT',
  approved_by     TEXT    NOT NULL DEFAULT '',
  approved_at     TEXT    NOT NULL DEFAULT '',
  cancelled_reason TEXT   NOT NULL DEFAULT '',
  notes           TEXT    NOT NULL DEFAULT '',
  created_by      TEXT    NOT NULL,
  created_at      TEXT    NOT NULL DEFAULT '',
  updated_at      TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS po_lines (
  id             TEXT    PRIMARY KEY,
  po_id          TEXT    NOT NULL REFERENCES purchase_orders(id),
  line_number    INTEGER NOT NULL,
  item_code      TEXT    NOT NULL DEFAULT '',
  item_name      TEXT    NOT NULL,
  unit           TEXT    NOT NULL,
  quantity       REAL   NOT NULL DEFAULT 0,
  unit_price     REAL   NOT NULL DEFAULT 0,
  discount_pct   REAL   NOT NULL DEFAULT 0,
  vat_rate       REAL   NOT NULL DEFAULT 0,
  vat_type       TEXT    NOT NULL DEFAULT 'VAT_0',
  account_id     TEXT    NOT NULL,
  vat_account_id TEXT    NOT NULL,
  line_total     REAL   NOT NULL DEFAULT 0,
  line_vat_amount REAL  NOT NULL DEFAULT 0,
  received_qty   REAL   NOT NULL DEFAULT 0,
  invoiced_qty   REAL   NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS goods_receipt_notes (
  id           TEXT    PRIMARY KEY,
  company_id   TEXT    NOT NULL,
  grn_number   TEXT    NOT NULL,
  po_id        TEXT    NOT NULL REFERENCES purchase_orders(id),
  receipt_date TEXT    NOT NULL,
  warehouse    TEXT    NOT NULL DEFAULT '',
  status       TEXT    NOT NULL DEFAULT 'DRAFT',
  notes        TEXT    NOT NULL DEFAULT '',
  created_by   TEXT    NOT NULL,
  created_at   TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS grn_lines (
  id                TEXT    PRIMARY KEY,
  grn_id            TEXT    NOT NULL REFERENCES goods_receipt_notes(id),
  po_line_id        TEXT    NOT NULL,
  item_code         TEXT    NOT NULL DEFAULT '',
  item_name         TEXT    NOT NULL,
  unit              TEXT    NOT NULL DEFAULT '',
  quantity_received REAL   NOT NULL DEFAULT 0,
  quantity_rejected REAL   NOT NULL DEFAULT 0,
  unit_price        REAL   NOT NULL DEFAULT 0,
  line_total        REAL   NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS supplier_invoices (
  id                    TEXT    PRIMARY KEY,
  company_id            TEXT    NOT NULL,
  invoice_number        TEXT    NOT NULL,
  invoice_date          TEXT    NOT NULL,
  po_id                 TEXT    NOT NULL DEFAULT '',
  grn_id                TEXT    NOT NULL DEFAULT '',
  supplier_id           TEXT    NOT NULL REFERENCES suppliers(id),
  supplier_name         TEXT    NOT NULL,
  supplier_tax_code     TEXT    NOT NULL,
  invoice_type          TEXT    NOT NULL DEFAULT 'domestic',
  currency              TEXT    NOT NULL DEFAULT 'VND',
  exchange_rate         REAL   NOT NULL DEFAULT 1,
  subtotal              REAL   NOT NULL DEFAULT 0,
  discount_amount       REAL   NOT NULL DEFAULT 0,
  tax_amount            REAL   NOT NULL DEFAULT 0,
  total_amount          REAL   NOT NULL DEFAULT 0,
  amount_paid           REAL   NOT NULL DEFAULT 0,
  balance_due           REAL   NOT NULL DEFAULT 0,
  due_date              TEXT    NOT NULL DEFAULT '',
  vat_deduction_status  TEXT    NOT NULL DEFAULT 'pending',
  e_invoice_data        TEXT    NOT NULL DEFAULT '',
  e_invoice_code        TEXT    NOT NULL DEFAULT '',
  status                TEXT    NOT NULL DEFAULT 'DRAFT',
  gl_posted             INTEGER NOT NULL DEFAULT 0,
  gl_posted_at          TEXT    NOT NULL DEFAULT '',
  notes                 TEXT    NOT NULL DEFAULT '',
  created_by            TEXT    NOT NULL,
  created_at            TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS invoice_lines (
  id            TEXT    PRIMARY KEY,
  invoice_id    TEXT    NOT NULL REFERENCES supplier_invoices(id),
  po_line_id    TEXT    NOT NULL DEFAULT '',
  grn_line_id   TEXT    NOT NULL DEFAULT '',
  item_code     TEXT    NOT NULL DEFAULT '',
  item_name     TEXT    NOT NULL,
  unit          TEXT    NOT NULL DEFAULT '',
  quantity      REAL   NOT NULL DEFAULT 0,
  unit_price    REAL   NOT NULL DEFAULT 0,
  vat_rate      REAL   NOT NULL DEFAULT 0,
  vat_type      TEXT    NOT NULL DEFAULT 'VAT_0',
  line_total    REAL   NOT NULL DEFAULT 0,
  line_vat_amount REAL NOT NULL DEFAULT 0,
  account_id    TEXT    NOT NULL,
  vat_account_id TEXT   NOT NULL
);

CREATE TABLE IF NOT EXISTS ap_transactions (
  id              TEXT    PRIMARY KEY,
  company_id      TEXT    NOT NULL,
  supplier_id     TEXT    NOT NULL REFERENCES suppliers(id),
  invoice_id      TEXT    NOT NULL DEFAULT '',
  transaction_type TEXT  NOT NULL,
  transaction_date  TEXT  NOT NULL,
  amount           REAL  NOT NULL DEFAULT 0,
  currency         TEXT  NOT NULL DEFAULT 'VND',
  reference_type   TEXT  NOT NULL DEFAULT '',
  reference_id     TEXT  NOT NULL DEFAULT '',
  notes            TEXT  NOT NULL DEFAULT '',
  created_at       TEXT  NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS cost_allocations (
  id               TEXT    PRIMARY KEY,
  company_id       TEXT    NOT NULL,
  invoice_id       TEXT    NOT NULL DEFAULT '',
  cost_type        TEXT    NOT NULL DEFAULT '',
  cost_amount      REAL   NOT NULL DEFAULT 0,
  allocation_method TEXT   NOT NULL DEFAULT '',
  allocated_lines  TEXT    NOT NULL DEFAULT '',
  notes            TEXT    NOT NULL DEFAULT '',
  created_at       TEXT    NOT NULL DEFAULT ''
);

-- Suppliers: per-company lookup, code uniqueness per company
CREATE INDEX IF NOT EXISTS idx_suppliers_company ON suppliers(company_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_suppliers_company_code ON suppliers(company_id, code);

-- Purchase Orders
CREATE INDEX IF NOT EXISTS idx_po_company      ON purchase_orders(company_id);
CREATE INDEX IF NOT EXISTS idx_po_supplier     ON purchase_orders(supplier_id);
CREATE INDEX IF NOT EXISTS idx_po_number       ON purchase_orders(po_number);
CREATE INDEX IF NOT EXISTS idx_po_status       ON purchase_orders(status);
CREATE INDEX IF NOT EXISTS idx_po_lines_po     ON po_lines(po_id);

-- GRNs
CREATE INDEX IF NOT EXISTS idx_grn_company  ON goods_receipt_notes(company_id);
CREATE INDEX IF NOT EXISTS idx_grn_po       ON goods_receipt_notes(po_id);
CREATE INDEX IF NOT EXISTS idx_grn_status   ON goods_receipt_notes(status);
CREATE INDEX IF NOT EXISTS idx_grn_number   ON goods_receipt_notes(grn_number);
CREATE INDEX IF NOT EXISTS idx_grn_lines_grn ON grn_lines(grn_id);

-- Supplier Invoices
CREATE INDEX IF NOT EXISTS idx_inv_company  ON supplier_invoices(company_id);
CREATE INDEX IF NOT EXISTS idx_inv_supplier ON supplier_invoices(supplier_id);
CREATE INDEX IF NOT EXISTS idx_inv_number   ON supplier_invoices(invoice_number);
CREATE INDEX IF NOT EXISTS idx_inv_status   ON supplier_invoices(status);
CREATE INDEX IF NOT EXISTS idx_inv_po       ON supplier_invoices(po_id);
CREATE INDEX IF NOT EXISTS idx_inv_grn      ON supplier_invoices(grn_id);
CREATE INDEX IF NOT EXISTS idx_inv_lines_inv ON invoice_lines(invoice_id);

-- AP Transactions
CREATE INDEX IF NOT EXISTS idx_apt_company  ON ap_transactions(company_id);
CREATE INDEX IF NOT EXISTS idx_apt_supplier ON ap_transactions(supplier_id);
CREATE INDEX IF NOT EXISTS idx_apt_invoice  ON ap_transactions(invoice_id);

-- Cost Allocations
CREATE INDEX IF NOT EXISTS idx_costalloc_invoice ON cost_allocations(invoice_id);
