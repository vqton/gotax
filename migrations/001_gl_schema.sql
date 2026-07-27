-- GoTax GL Module Schema
-- Based on Circular 200/2014/TT-BTC (updated by Circular 99/2025/TT-BTC)
-- Compatible with PostgreSQL

-- ============================================================
-- ACCOUNTING PERIOD
-- ============================================================
CREATE TABLE IF NOT EXISTS periods (
    id          VARCHAR(20) PRIMARY KEY,
    year        INTEGER NOT NULL,
    month       INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
    start_date  DATE NOT NULL,
    end_date    DATE NOT NULL,
    status      VARCHAR(10) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED', 'LOCKED')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (year, month)
);

-- ============================================================
-- CHART OF ACCOUNTS (He thong tai khoan ke toan)
-- ============================================================
CREATE TABLE IF NOT EXISTS accounts (
    code         VARCHAR(20) PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    name2        VARCHAR(255),           -- English/secondary name
    type         VARCHAR(10) NOT NULL CHECK (type IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE')),
    parent_code  VARCHAR(20) REFERENCES accounts(code),
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    is_foreign   BOOLEAN NOT NULL DEFAULT FALSE,
    detail_by    VARCHAR(20) CHECK (detail_by IN ('OBJECT', 'PROJECT', 'CONTRACT', 'COST_ITEM', 'DEPARTMENT')),
    is_parent    BOOLEAN NOT NULL DEFAULT FALSE,
    arrears_days INTEGER DEFAULT 0,
    note         TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_parent ON accounts(parent_code);
CREATE INDEX idx_accounts_type ON accounts(type);
CREATE INDEX idx_accounts_active ON accounts(is_active);

-- ============================================================
-- JOURNAL ENTRIES (Chung tu ghi so)
-- ============================================================
CREATE TABLE IF NOT EXISTS journal_entries (
    id           VARCHAR(20) PRIMARY KEY,
    entry_number VARCHAR(50) NOT NULL,
    entry_date   DATE NOT NULL,
    period_id    VARCHAR(20) REFERENCES periods(id),
    description  TEXT NOT NULL,
    status       VARCHAR(10) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'POSTED', 'CANCELLED')),
    created_by   VARCHAR(100),
    approved_by  VARCHAR(100),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    posted_at    TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_journal_entries_date ON journal_entries(entry_date);
CREATE INDEX idx_journal_entries_period ON journal_entries(period_id);
CREATE INDEX idx_journal_entries_status ON journal_entries(status);

-- ============================================================
-- JOURNAL LINES (Chi tiet chung tu)
-- ============================================================
CREATE TABLE IF NOT EXISTS journal_lines (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id       VARCHAR(20) NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    line_number    INTEGER NOT NULL,
    account_code   VARCHAR(20) NOT NULL REFERENCES accounts(code),
    debit_amount   NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (debit_amount >= 0),
    credit_amount  NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (credit_amount >= 0),
    description    TEXT,
    object_id      VARCHAR(50),       -- KH, NCC, NV...
    project_id     VARCHAR(50),
    contract_id    VARCHAR(50),
    cost_item_id   VARCHAR(50),
    department_id  VARCHAR(50),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_at_least_one_amount CHECK (debit_amount > 0 OR credit_amount > 0)
);

CREATE INDEX idx_journal_lines_entry ON journal_lines(entry_id);
CREATE INDEX idx_journal_lines_account ON journal_lines(account_code);

-- ============================================================
-- ACCOUNT BALANCES (So du tai khoan - materialized for perf)
-- ============================================================
CREATE TABLE IF NOT EXISTS account_balances (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_code       VARCHAR(20) NOT NULL REFERENCES accounts(code),
    period_id          VARCHAR(20) NOT NULL REFERENCES periods(id),
    open_balance_debit  NUMERIC(18,2) NOT NULL DEFAULT 0,
    open_balance_credit NUMERIC(18,2) NOT NULL DEFAULT 0,
    period_debit       NUMERIC(18,2) NOT NULL DEFAULT 0,
    period_credit      NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_code, period_id)
);

CREATE INDEX idx_balances_period ON account_balances(period_id);
CREATE INDEX idx_balances_account ON account_balances(account_code);

-- ============================================================
-- SEED DEFAULT ACCOUNTS (Circular 200 - abbreviated)
-- ============================================================
INSERT INTO periods (id, year, month, start_date, end_date, status) VALUES
    ('P-2026-01', 2026, 1, '2026-01-01', '2026-01-31', 'OPEN'),
    ('P-2026-02', 2026, 2, '2026-02-01', '2026-02-28', 'OPEN'),
    ('P-2026-03', 2026, 3, '2026-03-01', '2026-03-31', 'OPEN'),
    ('P-2026-04', 2026, 4, '2026-04-01', '2026-04-30', 'OPEN'),
    ('P-2026-05', 2026, 5, '2026-05-01', '2026-05-31', 'OPEN'),
    ('P-2026-06', 2026, 6, '2026-06-01', '2026-06-30', 'OPEN'),
    ('P-2026-07', 2026, 7, '2026-07-01', '2026-07-31', 'OPEN'),
    ('P-2026-08', 2026, 8, '2026-08-01', '2026-08-31', 'OPEN'),
    ('P-2026-09', 2026, 9, '2026-09-01', '2026-09-30', 'OPEN'),
    ('P-2026-10', 2026, 10, '2026-10-01', '2026-10-31', 'OPEN'),
    ('P-2026-11', 2026, 11, '2026-11-01', '2026-11-30', 'OPEN'),
    ('P-2026-12', 2026, 12, '2026-12-01', '2026-12-31', 'OPEN')
ON CONFLICT DO NOTHING;

-- Seed basic chart of accounts (Type 1: ASSET)
INSERT INTO accounts (code, name, type, is_parent, is_active) VALUES
    ('1', 'Tai san ngan han', 'ASSET', TRUE, TRUE),
    ('11', 'Tien va tuong duong tien', 'ASSET', TRUE, TRUE),
    ('111', 'Tien mat', 'ASSET', TRUE, TRUE),
    ('1111', 'Tien mat VND', 'ASSET', FALSE, TRUE),
    ('1112', 'Tien mat USD', 'ASSET', FALSE, TRUE),
    ('112', 'Tien gui ngan hang', 'ASSET', TRUE, TRUE),
    ('1121', 'Tien gui VND', 'ASSET', FALSE, TRUE),
    ('1122', 'Tien gui USD', 'ASSET', FALSE, TRUE),
    ('131', 'Phai thu cua khach hang', 'ASSET', TRUE, TRUE),
    ('1311', 'Phai thu KH trong nuoc', 'ASSET', FALSE, TRUE),
    ('1312', 'Phai thu KH xuat khau', 'ASSET', FALSE, TRUE),
    ('156', 'Hang hoa', 'ASSET', TRUE, TRUE),
    ('1561', 'Hang hoa mua vao', 'ASSET', FALSE, TRUE),
    ('1562', 'Hang hoa ban ra', 'ASSET', FALSE, TRUE),
    ('2', 'Tai san dai han', 'ASSET', TRUE, TRUE),
    ('211', 'Tai san co dinh huu hinh', 'ASSET', TRUE, TRUE),
    ('2111', 'Nha cua, vat kien truc', 'ASSET', FALSE, TRUE),
    ('2112', 'May moc thiet bi', 'ASSET', FALSE, TRUE),
    ('214', 'Hao mon TSCD', 'ASSET', TRUE, TRUE),
    ('2141', 'Hao mon TSCD huu hinh', 'ASSET', FALSE, TRUE),
    ('3', 'No phai tra', 'LIABILITY', TRUE, TRUE),
    ('31', 'No ngan han', 'LIABILITY', TRUE, TRUE),
    ('311', 'Vay ngan hang', 'LIABILITY', TRUE, TRUE),
    ('3111', 'Vay ngan hang ngan han', 'LIABILITY', FALSE, TRUE),
    ('331', 'Phai tra nguoi ban', 'LIABILITY', TRUE, TRUE),
    ('3311', 'Phai tra NB trong nuoc', 'LIABILITY', FALSE, TRUE),
    ('3312', 'Phai tra NB nuoc ngoai', 'LIABILITY', FALSE, TRUE),
    ('333', 'Thue va cac khoan phai tra NN', 'LIABILITY', TRUE, TRUE),
    ('3331', 'Thue GTGT phai nap', 'LIABILITY', FALSE, TRUE),
    ('33311', 'Thue GTGT dau ra', 'LIABILITY', FALSE, TRUE),
    ('33312', 'Thue GTGT hang nhap khau', 'LIABILITY', FALSE, TRUE),
    ('3334', 'Thue TNDN', 'LIABILITY', FALSE, TRUE),
    ('4', 'Nguon von chu so huu', 'EQUITY', TRUE, TRUE),
    ('411', 'Von dau tu cua CSH', 'EQUITY', TRUE, TRUE),
    ('4111', 'Von co phan', 'EQUITY', FALSE, TRUE),
    ('421', 'Loi nhuan sau thue', 'EQUITY', TRUE, TRUE),
    ('4211', 'Loi nhuan nam truoc', 'EQUITY', FALSE, TRUE),
    ('4212', 'Loi nhuan nam nay', 'EQUITY', FALSE, TRUE),
    ('5', 'Doanh thu', 'REVENUE', TRUE, TRUE),
    ('511', 'Doanh thu ban hang', 'REVENUE', TRUE, TRUE),
    ('5111', 'Doanh thu ban hang hoa', 'REVENUE', FALSE, TRUE),
    ('5112', 'Doanh thu cung cap DV', 'REVENUE', FALSE, TRUE),
    ('515', 'Doanh thu hoat dong TC', 'REVENUE', FALSE, TRUE),
    ('6', 'Chi phi san xuat kinh doanh', 'EXPENSE', TRUE, TRUE),
    ('621', 'Chi phi NVL truc tiep', 'EXPENSE', FALSE, TRUE),
    ('622', 'Chi phi nhan cong truc tiep', 'EXPENSE', FALSE, TRUE),
    ('627', 'Chi phi san xuat chung', 'EXPENSE', TRUE, TRUE),
    ('6271', 'Chi phi nhan vien', 'EXPENSE', FALSE, TRUE),
    ('6272', 'Chi phi vat lieu', 'EXPENSE', FALSE, TRUE),
    ('632', 'Gia von hang ban', 'EXPENSE', FALSE, TRUE),
    ('635', 'Chi phi hoat dong TC', 'EXPENSE', FALSE, TRUE),
    ('641', 'Chi phi ban hang', 'EXPENSE', TRUE, TRUE),
    ('6411', 'CP ban hang - nhan vien', 'EXPENSE', FALSE, TRUE),
    ('6412', 'CP ban hang - quang cao', 'EXPENSE', FALSE, TRUE),
    ('642', 'Chi phi quan ly DN', 'EXPENSE', TRUE, TRUE),
    ('6421', 'CP QLDN - nhan vien', 'EXPENSE', FALSE, TRUE),
    ('6422', 'CP QLDN - van phong', 'EXPENSE', FALSE, TRUE),
    ('7', 'Thu nhap khac', 'REVENUE', TRUE, TRUE),
    ('711', 'Thu nhap khac', 'REVENUE', FALSE, TRUE),
    ('8', 'Chi phi khac', 'EXPENSE', TRUE, TRUE),
    ('811', 'Chi phi khac', 'EXPENSE', FALSE, TRUE)
ON CONFLICT DO NOTHING;