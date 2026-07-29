/*
Warehouse Module Schema — Inventory Management — v1
Complies: Circular 99/2025/TT-BTC, Decree 123/2020.
No third-party integrations. GL auto-posting deferred.
*/
CREATE TABLE IF NOT EXISTS warehouses (
  id          TEXT    PRIMARY KEY,
  company_id  TEXT    NOT NULL,
  code        TEXT    NOT NULL,
  name        TEXT    NOT NULL,
  address     TEXT    NOT NULL DEFAULT '',
  manager     TEXT    NOT NULL DEFAULT '',
  is_active   INTEGER NOT NULL DEFAULT 1,
  created_at  TEXT    NOT NULL DEFAULT '',
  updated_at  TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS item_categories (
  id          TEXT    PRIMARY KEY,
  company_id  TEXT    NOT NULL,
  code        TEXT    NOT NULL,
  name        TEXT    NOT NULL,
  description TEXT    NOT NULL DEFAULT '',
  parent_id   TEXT    NOT NULL DEFAULT '',
  is_active   INTEGER NOT NULL DEFAULT 1,
  created_at  TEXT    NOT NULL DEFAULT '',
  updated_at  TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS items (
  id                TEXT    PRIMARY KEY,
  company_id        TEXT    NOT NULL,
  code              TEXT    NOT NULL,
  name              TEXT    NOT NULL,
  category_id       TEXT    NOT NULL DEFAULT '',
  unit              TEXT    NOT NULL DEFAULT '',
  base_price        NUMERIC(18,2) NOT NULL DEFAULT 0,
  min_stock         NUMERIC(18,2) NOT NULL DEFAULT 0,
  max_stock         NUMERIC(18,2) NOT NULL DEFAULT 0,
  valuation_method  TEXT    NOT NULL DEFAULT 'weighted_average',
  tax_rate          NUMERIC(5,2)  NOT NULL DEFAULT 0,
  is_active         INTEGER NOT NULL DEFAULT 1,
  notes             TEXT    NOT NULL DEFAULT '',
  created_at        TEXT    NOT NULL DEFAULT '',
  updated_at        TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS stock_balances (
  id                  TEXT    PRIMARY KEY,
  company_id          TEXT    NOT NULL,
  warehouse_id        TEXT    NOT NULL REFERENCES warehouses(id),
  item_id             TEXT    NOT NULL REFERENCES items(id),
  period              TEXT    NOT NULL,
  quantity            NUMERIC(18,4) NOT NULL DEFAULT 0,
  unit_cost           NUMERIC(18,2) NOT NULL DEFAULT 0,
  total_cost          NUMERIC(18,2) NOT NULL DEFAULT 0,
  last_transaction_at TEXT    NOT NULL DEFAULT '',
  created_at          TEXT    NOT NULL DEFAULT '',
  updated_at          TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS inventory_transactions (
  id          TEXT    PRIMARY KEY,
  company_id  TEXT    NOT NULL,
  warehouse_id TEXT   NOT NULL,
  item_id     TEXT    NOT NULL,
  trans_type  TEXT    NOT NULL,
  ref_type    TEXT    NOT NULL DEFAULT '',
  ref_id      TEXT    NOT NULL DEFAULT '',
  qty_before  NUMERIC(18,4) NOT NULL DEFAULT 0,
  quantity    NUMERIC(18,4) NOT NULL DEFAULT 0,
  qty_after   NUMERIC(18,4) NOT NULL DEFAULT 0,
  unit_cost   NUMERIC(18,2) NOT NULL DEFAULT 0,
  total_cost  NUMERIC(18,2) NOT NULL DEFAULT 0,
  created_by  TEXT    NOT NULL DEFAULT '',
  created_at  TEXT    NOT NULL DEFAULT '',
  notes       TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS stock_transfers (
  id                TEXT    PRIMARY KEY,
  company_id        TEXT    NOT NULL,
  transfer_number   TEXT    NOT NULL,
  from_warehouse_id TEXT    NOT NULL,
  to_warehouse_id   TEXT    NOT NULL,
  status            TEXT    NOT NULL DEFAULT 'DRAFT',
  transfer_date     TEXT    NOT NULL DEFAULT '',
  created_by        TEXT    NOT NULL DEFAULT '',
  approved_by       TEXT    NOT NULL DEFAULT '',
  approved_at       TEXT    NOT NULL DEFAULT '',
  completed_by      TEXT    NOT NULL DEFAULT '',
  completed_at      TEXT    NOT NULL DEFAULT '',
  cancelled_reason  TEXT    NOT NULL DEFAULT '',
  notes             TEXT    NOT NULL DEFAULT '',
  created_at        TEXT    NOT NULL DEFAULT '',
  updated_at        TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS transfer_items (
  id          TEXT    PRIMARY KEY,
  transfer_id TEXT    NOT NULL REFERENCES stock_transfers(id),
  item_id     TEXT    NOT NULL,
  quantity    NUMERIC(18,4) NOT NULL DEFAULT 0,
  unit_cost   NUMERIC(18,2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS stock_adjustments (
  id                 TEXT    PRIMARY KEY,
  company_id         TEXT    NOT NULL,
  warehouse_id       TEXT    NOT NULL,
  adjustment_number  TEXT    NOT NULL,
  adj_type           TEXT    NOT NULL,
  reason             TEXT    NOT NULL DEFAULT '',
  status             TEXT    NOT NULL DEFAULT 'DRAFT',
  created_by         TEXT    NOT NULL DEFAULT '',
  approved_by        TEXT    NOT NULL DEFAULT '',
  approved_at        TEXT    NOT NULL DEFAULT '',
  posted_at          TEXT    NOT NULL DEFAULT '',
  rejected_reason    TEXT    NOT NULL DEFAULT '',
  notes              TEXT    NOT NULL DEFAULT '',
  created_at         TEXT    NOT NULL DEFAULT '',
  updated_at         TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS adjustment_items (
  id             TEXT    PRIMARY KEY,
  adjustment_id  TEXT    NOT NULL REFERENCES stock_adjustments(id),
  item_id        TEXT    NOT NULL,
  qty_before     NUMERIC(18,4) NOT NULL DEFAULT 0,
  qty_after      NUMERIC(18,4) NOT NULL DEFAULT 0,
  unit_cost      NUMERIC(18,2) NOT NULL DEFAULT 0,
  reason         TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS stock_takes (
  id            TEXT    PRIMARY KEY,
  company_id    TEXT    NOT NULL,
  warehouse_id  TEXT    NOT NULL,
  take_number   TEXT    NOT NULL,
  status        TEXT    NOT NULL DEFAULT 'PLANNING',
  take_date     TEXT    NOT NULL DEFAULT '',
  created_by    TEXT    NOT NULL DEFAULT '',
  verified_by   TEXT    NOT NULL DEFAULT '',
  verified_at   TEXT    NOT NULL DEFAULT '',
  posted_at     TEXT    NOT NULL DEFAULT '',
  notes         TEXT    NOT NULL DEFAULT '',
  created_at    TEXT    NOT NULL DEFAULT '',
  updated_at    TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS take_items (
  id            TEXT    PRIMARY KEY,
  take_id       TEXT    NOT NULL REFERENCES stock_takes(id),
  item_id       TEXT    NOT NULL,
  expected_qty  NUMERIC(18,4) NOT NULL DEFAULT 0,
  actual_qty    NUMERIC(18,4) NOT NULL DEFAULT 0,
  unit_cost     NUMERIC(18,2) NOT NULL DEFAULT 0,
  variance      NUMERIC(18,4) NOT NULL DEFAULT 0,
  notes         TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS valuation_runs (
  id              TEXT    PRIMARY KEY,
  company_id      TEXT    NOT NULL,
  valuation_date  TEXT    NOT NULL,
  method          TEXT    NOT NULL DEFAULT 'weighted_average',
  status          TEXT    NOT NULL DEFAULT 'PENDING',
  created_by      TEXT    NOT NULL DEFAULT '',
  completed_at    TEXT    NOT NULL DEFAULT '',
  error_log       TEXT    NOT NULL DEFAULT '',
  notes           TEXT    NOT NULL DEFAULT '',
  created_at      TEXT    NOT NULL DEFAULT ''
);
