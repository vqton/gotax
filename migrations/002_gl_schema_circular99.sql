-- GoTax GL Module Schema v2 — Circular 99/2025/TT-BTC compliant
-- Effective 1/1/2026. Replaces TT200 schema.
-- PostgreSQL 16+

BEGIN;

-- ============================================================
-- EXTENSION
-- ============================================================
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================================
-- USERS / AUTH
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(100) NOT NULL,
    email         VARCHAR(100),
    role          VARCHAR(20) NOT NULL DEFAULT 'accountant'
                  CHECK (role IN ('admin', 'chief_accountant', 'accountant', 'viewer')),
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- ACCOUNTING PERIOD
-- ============================================================
CREATE TABLE IF NOT EXISTS periods (
    id         VARCHAR(20) PRIMARY KEY,
    year       INTEGER NOT NULL,
    month      INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
    start_date DATE NOT NULL,
    end_date   DATE NOT NULL,
    status     VARCHAR(10) NOT NULL DEFAULT 'OPEN'
               CHECK (status IN ('OPEN', 'CLOSING', 'CLOSED', 'LOCKED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (year, month)
);

-- ============================================================
-- CHART OF ACCOUNTS — Circular 99/2025/TT-BTC Appendix II
-- 71 Level-1 accounts, enterprises may add detail accounts
-- ============================================================
CREATE TABLE IF NOT EXISTS accounts (
    code         VARCHAR(20) PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    name2        VARCHAR(255),
    type         VARCHAR(10) NOT NULL
                 CHECK (type IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE')),
    parent_code  VARCHAR(20) REFERENCES accounts(code),
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    is_foreign   BOOLEAN NOT NULL DEFAULT FALSE,
    detail_by    VARCHAR(20)
                 CHECK (detail_by IN ('OBJECT', 'PROJECT', 'CONTRACT', 'COST_ITEM', 'DEPARTMENT')),
    is_parent    BOOLEAN NOT NULL DEFAULT FALSE,
    note         TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accounts_parent ON accounts(parent_code);
CREATE INDEX IF NOT EXISTS idx_accounts_type ON accounts(type);
CREATE INDEX IF NOT EXISTS idx_accounts_active ON accounts(is_active);

-- ============================================================
-- EXCHANGE RATES (for multi-currency)
-- ============================================================
CREATE TABLE IF NOT EXISTS exchange_rates (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    currency_code VARCHAR(3) NOT NULL,
    rate_date     DATE NOT NULL,
    buy_rate      NUMERIC(18,6),
    sell_rate     NUMERIC(18,6),
    average_rate  NUMERIC(18,6) NOT NULL,
    source        VARCHAR(50) DEFAULT 'MANUAL',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (currency_code, rate_date)
);

CREATE INDEX IF NOT EXISTS idx_rates_currency_date ON exchange_rates(currency_code, rate_date);

-- ============================================================
-- JOURNAL ENTRIES (Chung tu ghi so — Circular 99)
-- ============================================================
CREATE TABLE IF NOT EXISTS journal_entries (
    id              VARCHAR(20) PRIMARY KEY,
    company_id      UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    entry_number    VARCHAR(50) NOT NULL,
    voucher_type    VARCHAR(10) NOT NULL DEFAULT 'KHAC'
                    CHECK (voucher_type IN ('THU', 'CHI', 'BAN', 'MUA', 'KHAC', 'KC')),
    entry_date      DATE NOT NULL,
    accounting_date DATE NOT NULL,
    period_id       VARCHAR(20) NOT NULL REFERENCES periods(id),
    description     TEXT NOT NULL,
    status          VARCHAR(15) NOT NULL DEFAULT 'DRAFT'
                    CHECK (status IN ('DRAFT', 'REVIEWING', 'APPROVED', 'POSTED', 'CANCELLED')),
    currency_code   VARCHAR(3) NOT NULL DEFAULT 'VND',
    exchange_rate   NUMERIC(18,6) DEFAULT 1,
    created_by      UUID REFERENCES users(id),
    approved_by     UUID REFERENCES users(id),
    reviewed_by     UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    posted_at       TIMESTAMPTZ,
    approved_at     TIMESTAMPTZ,
    CONSTRAINT chk_posted_status CHECK (
        (status != 'POSTED') OR (posted_at IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_je_date ON journal_entries(entry_date);
CREATE INDEX IF NOT EXISTS idx_je_period ON journal_entries(period_id);
CREATE INDEX IF NOT EXISTS idx_je_status ON journal_entries(status);
CREATE INDEX IF NOT EXISTS idx_je_voucher ON journal_entries(voucher_type);
CREATE INDEX IF NOT EXISTS idx_je_company ON journal_entries(company_id);

-- ============================================================
-- JOURNAL LINES (Chi tiet chung tu)
-- ============================================================
CREATE TABLE IF NOT EXISTS journal_lines (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id        VARCHAR(20) NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    line_number     INTEGER NOT NULL,
    account_code    VARCHAR(20) NOT NULL REFERENCES accounts(code),
    debit_amount    NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (debit_amount >= 0),
    credit_amount   NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (credit_amount >= 0),
    description     TEXT,
    currency_code   VARCHAR(3) DEFAULT 'VND',
    foreign_amount  NUMERIC(18,2),
    exchange_rate   NUMERIC(18,6),
    object_id       VARCHAR(50),
    project_id      VARCHAR(50),
    contract_id     VARCHAR(50),
    cost_item_id    VARCHAR(50),
    department_id   VARCHAR(50),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_at_least_one_amount CHECK (debit_amount > 0 OR credit_amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_jl_entry ON journal_lines(entry_id);
CREATE INDEX IF NOT EXISTS idx_jl_account ON journal_lines(account_code);

-- ============================================================
-- ACCOUNT BALANCES (So du tai khoan — materialized)
-- ============================================================
CREATE TABLE IF NOT EXISTS account_balances (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_code        VARCHAR(20) NOT NULL REFERENCES accounts(code),
    period_id           VARCHAR(20) NOT NULL REFERENCES periods(id),
    open_balance_debit   NUMERIC(18,2) NOT NULL DEFAULT 0,
    open_balance_credit  NUMERIC(18,2) NOT NULL DEFAULT 0,
    period_debit        NUMERIC(18,2) NOT NULL DEFAULT 0,
    period_credit       NUMERIC(18,2) NOT NULL DEFAULT 0,
    last_updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_code, period_id)
);

CREATE INDEX IF NOT EXISTS idx_bal_period ON account_balances(period_id);
CREATE INDEX IF NOT EXISTS idx_bal_account ON account_balances(account_code);

-- ============================================================
-- AUDIT TRAIL (append-only, tamper-evident)
-- ============================================================
CREATE TABLE IF NOT EXISTS audit_log (
    id            BIGSERIAL PRIMARY KEY,
    user_id       UUID REFERENCES users(id),
    username      VARCHAR(50),
    ip_address    VARCHAR(45),
    action        VARCHAR(20) NOT NULL,
    entity_type   VARCHAR(30) NOT NULL,
    entity_id     VARCHAR(50),
    old_value     JSONB,
    new_value     JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_log(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_time ON audit_log(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_log(action);

-- Prevent deletion or update on audit_log
CREATE OR REPLACE FUNCTION f_audit_log_immutable()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'audit_log is append-only: cannot %', TG_OP;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS t_audit_log_immutable ON audit_log;
CREATE TRIGGER t_audit_log_immutable
    BEFORE UPDATE OR DELETE ON audit_log
    FOR EACH ROW EXECUTE FUNCTION f_audit_log_immutable();

-- ============================================================
-- CLOSING ENTRY TEMPLATES
-- ============================================================
CREATE TABLE IF NOT EXISTS closing_templates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name             VARCHAR(100) NOT NULL,
    description      TEXT,
    sequence_order   INTEGER NOT NULL DEFAULT 0,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS closing_template_lines (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id     UUID NOT NULL REFERENCES closing_templates(id) ON DELETE CASCADE,
    line_number     INTEGER NOT NULL,
    debit_account   VARCHAR(20) NOT NULL REFERENCES accounts(code),
    credit_account  VARCHAR(20) NOT NULL REFERENCES accounts(code),
    formula         VARCHAR(50) NOT NULL DEFAULT 'BALANCE',
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- SEED: Circular 99 Chart of Accounts (71 Level-1 accounts)
-- ============================================================
INSERT INTO accounts (code, name, type, is_parent, is_active) VALUES
    -- LOAI TAI SAN (ASSETS)
    ('1', 'Tai san', 'ASSET', TRUE, TRUE),
    ('111', 'Tien mat', 'ASSET', TRUE, TRUE),
    ('1111', 'Tien mat VND', 'ASSET', FALSE, TRUE),
    ('1112', 'Tien mat USD', 'ASSET', FALSE, TRUE),
    ('112', 'Tien gui khong ky han', 'ASSET', TRUE, TRUE),
    ('1121', 'Tien gui VND', 'ASSET', FALSE, TRUE),
    ('1122', 'Tien gui USD', 'ASSET', FALSE, TRUE),
    ('113', 'Tien dang chuyen', 'ASSET', FALSE, TRUE),
    ('121', 'Chung khoan kinh doanh', 'ASSET', TRUE, TRUE),
    ('1211', 'Co phieu', 'ASSET', FALSE, TRUE),
    ('1212', 'Trai phieu', 'ASSET', FALSE, TRUE),
    ('128', 'Dau tu nam giu den ngay dao han', 'ASSET', TRUE, TRUE),
    ('1281', 'Tien gui co ky han', 'ASSET', FALSE, TRUE),
    ('1282', 'Trai phieu', 'ASSET', FALSE, TRUE),
    ('1283', 'Cho vay', 'ASSET', FALSE, TRUE),
    ('1288', 'Dau tu khac', 'ASSET', FALSE, TRUE),
    ('131', 'Phai thu cua khach hang', 'ASSET', TRUE, TRUE),
    ('1311', 'Phai thu KH trong nuoc', 'ASSET', FALSE, TRUE),
    ('1312', 'Phai thu KH nuoc ngoai', 'ASSET', FALSE, TRUE),
    ('133', 'Thue GTGT duoc khau tru', 'ASSET', TRUE, TRUE),
    ('1331', 'Thue GTGT hang hoa DV', 'ASSET', FALSE, TRUE),
    ('1332', 'Thue GTGT TSCD', 'ASSET', FALSE, TRUE),
    ('136', 'Phai thu noi bo', 'ASSET', TRUE, TRUE),
    ('1361', 'Von KD o don vi truc thuoc', 'ASSET', FALSE, TRUE),
    ('1368', 'Phai thu noi bo khac', 'ASSET', FALSE, TRUE),
    ('138', 'Phai thu khac', 'ASSET', TRUE, TRUE),
    ('1381', 'Tai san thieu cho xu ly', 'ASSET', FALSE, TRUE),
    ('1383', 'Thue TDBB hang nhap khau', 'ASSET', FALSE, TRUE),
    ('1385', 'Phai thu ve cho muon TSCD', 'ASSET', FALSE, TRUE),
    ('1388', 'Phai thu khac', 'ASSET', FALSE, TRUE),
    ('141', 'Tam ung', 'ASSET', FALSE, TRUE),
    ('151', 'Hang mua dang di duong', 'ASSET', FALSE, TRUE),
    ('152', 'Nguyen lieu vat lieu', 'ASSET', FALSE, TRUE),
    ('153', 'Cong cu dung cu', 'ASSET', FALSE, TRUE),
    ('154', 'Chi phi SXKD do dang', 'ASSET', FALSE, TRUE),
    ('155', 'San pham', 'ASSET', FALSE, TRUE),
    ('156', 'Hang hoa', 'ASSET', FALSE, TRUE),
    ('157', 'Hang gui ban', 'ASSET', FALSE, TRUE),
    ('158', 'Hang hoa kho bao thuue', 'ASSET', FALSE, TRUE),
    ('171', 'Giao dich mua lai trai phieu CP', 'ASSET', FALSE, TRUE),
    ('211', 'TSCD huu hinh', 'ASSET', TRUE, TRUE),
    ('2111', 'Nha cua vat kien truc', 'ASSET', FALSE, TRUE),
    ('2112', 'May moc thiet bi', 'ASSET', FALSE, TRUE),
    ('2113', 'Phuong tien van tai', 'ASSET', FALSE, TRUE),
    ('2114', 'Thiet bi quan ly', 'ASSET', FALSE, TRUE),
    ('2115', 'Cay lau nam', 'ASSET', FALSE, TRUE),
    ('2118', 'TSCD khac', 'ASSET', FALSE, TRUE),
    ('212', 'TSCD thue tai chinh', 'ASSET', FALSE, TRUE),
    ('213', 'TSCD vo hinh', 'ASSET', TRUE, TRUE),
    ('2131', 'Quyen su dung dat', 'ASSET', FALSE, TRUE),
    ('2132', 'Ban quyen bang sang che', 'ASSET', FALSE, TRUE),
    ('2133', 'Nhan hieu thuong mai', 'ASSET', FALSE, TRUE),
    ('2134', 'Phan mem may tinh', 'ASSET', FALSE, TRUE),
    ('2135', 'Gia tri loi the', 'ASSET', FALSE, TRUE),
    ('214', 'Hao mon TSCD', 'ASSET', TRUE, TRUE),
    ('2141', 'Hao mon TSCD huu hinh', 'ASSET', FALSE, TRUE),
    ('2142', 'Hao mon TSCD thue TC', 'ASSET', FALSE, TRUE),
    ('2143', 'Hao mon TSCD vo hinh', 'ASSET', FALSE, TRUE),
    ('2147', 'Hao mon BDS dau tu', 'ASSET', FALSE, TRUE),
    ('215', 'Tai san sinh hoc', 'ASSET', TRUE, TRUE),
    ('2151', 'Suc vat nuoi SP dinh ky', 'ASSET', FALSE, TRUE),
    ('2152', 'Suc vat nuoi SP mot lan', 'ASSET', FALSE, TRUE),
    ('2153', 'Cay trong thoi vu', 'ASSET', FALSE, TRUE),
    ('217', 'Bat dong san dau tu', 'ASSET', FALSE, TRUE),
    ('221', 'Dau tu vao cong ty con', 'ASSET', FALSE, TRUE),
    ('222', 'Dau tu vao cong ty LK LD', 'ASSET', FALSE, TRUE),
    ('228', 'Dau tu khac', 'ASSET', TRUE, TRUE),
    ('2281', 'Dau tu von vao don vi khac', 'ASSET', FALSE, TRUE),
    ('2288', 'Dau tu khac', 'ASSET', FALSE, TRUE),
    ('229', 'Du phong ton that tai san', 'ASSET', TRUE, TRUE),
    ('2291', 'Du phong CK kinh doanh', 'ASSET', FALSE, TRUE),
    ('2292', 'Du phong ton that dau tu', 'ASSET', FALSE, TRUE),
    ('2293', 'Du phong phai thu kho doi', 'ASSET', FALSE, TRUE),
    ('2294', 'Du phong giam gia hang ton kho', 'ASSET', FALSE, TRUE),
    ('2295', 'Du phong ton that TSSH', 'ASSET', FALSE, TRUE),
    ('241', 'Xay dung co ban do dang', 'ASSET', TRUE, TRUE),
    ('2411', 'Mua sam TSCD', 'ASSET', FALSE, TRUE),
    ('2412', 'Xay dung co ban', 'ASSET', FALSE, TRUE),
    ('2413', 'Sua chua bao duong dinh ky', 'ASSET', FALSE, TRUE),
    ('2414', 'Nang cap cai tao TSCD', 'ASSET', FALSE, TRUE),
    ('242', 'Chi phi cho phan bo', 'ASSET', FALSE, TRUE),
    ('243', 'Tai san thue TN hoan lai', 'ASSET', FALSE, TRUE),
    ('244', 'Cam co the chap ky quy', 'ASSET', FALSE, TRUE),

    -- LOAI NO PHAI TRA (LIABILITIES)
    ('3', 'No phai tra', 'LIABILITY', TRUE, TRUE),
    ('331', 'Phai tra nguoi ban', 'LIABILITY', TRUE, TRUE),
    ('3311', 'Phai tra NB trong nuoc', 'LIABILITY', FALSE, TRUE),
    ('3312', 'Phai tra NB nuoc ngoai', 'LIABILITY', FALSE, TRUE),
    ('332', 'Phai tra co tuc loi nhuan', 'LIABILITY', FALSE, TRUE),
    ('333', 'Thue va cac khoan phai nap NN', 'LIABILITY', TRUE, TRUE),
    ('3331', 'Thue GTGT phai nap', 'LIABILITY', FALSE, TRUE),
    ('33311', 'Thue GTGT dau ra', 'LIABILITY', FALSE, TRUE),
    ('33312', 'Thue GTGT hang nhap khau', 'LIABILITY', FALSE, TRUE),
    ('3332', 'Thue tieu thu dac biet', 'LIABILITY', FALSE, TRUE),
    ('3333', 'Thue XNK', 'LIABILITY', FALSE, TRUE),
    ('3334', 'Thue TNDN', 'LIABILITY', FALSE, TRUE),
    ('3335', 'Thue TNCN', 'LIABILITY', FALSE, TRUE),
    ('3336', 'Thue tai nguyen', 'LIABILITY', FALSE, TRUE),
    ('3337', 'Thue nha dat tien thue dat', 'LIABILITY', FALSE, TRUE),
    ('3338', 'Thue BVMT va thue khac', 'LIABILITY', TRUE, TRUE),
    ('33381', 'Thue BVMT', 'LIABILITY', FALSE, TRUE),
    ('33382', 'Thue khac', 'LIABILITY', FALSE, TRUE),
    ('3339', 'Phi le phi', 'LIABILITY', FALSE, TRUE),
    ('334', 'Phai tra nguoi lao dong', 'LIABILITY', FALSE, TRUE),
    ('335', 'Chi phi phai tra', 'LIABILITY', FALSE, TRUE),
    ('336', 'Phai tra noi bo', 'LIABILITY', TRUE, TRUE),
    ('3361', 'Phai tra noi bo von KD', 'LIABILITY', FALSE, TRUE),
    ('3368', 'Phai tra noi bo khac', 'LIABILITY', FALSE, TRUE),
    ('337', 'Thanh toan theo tien do HD', 'LIABILITY', FALSE, TRUE),
    ('338', 'Phai tra khac', 'LIABILITY', TRUE, TRUE),
    ('3381', 'Tai san thua cho xu ly', 'LIABILITY', FALSE, TRUE),
    ('3382', 'Kinh phi cong doan', 'LIABILITY', FALSE, TRUE),
    ('3383', 'Bao hiem xa hoi', 'LIABILITY', FALSE, TRUE),
    ('3384', 'Bao hiem y te', 'LIABILITY', FALSE, TRUE),
    ('3385', 'Bao hiem that nghiep', 'LIABILITY', FALSE, TRUE),
    ('3387', 'Doanh thu cho phan bo', 'LIABILITY', FALSE, TRUE),
    ('3388', 'Phai tra khac', 'LIABILITY', FALSE, TRUE),
    ('341', 'Vay va no thue TC', 'LIABILITY', TRUE, TRUE),
    ('3411', 'Vay', 'LIABILITY', FALSE, TRUE),
    ('3412', 'No thue tai chinh', 'LIABILITY', FALSE, TRUE),
    ('343', 'Trai phieu phat hanh', 'LIABILITY', TRUE, TRUE),
    ('3431', 'Trai phieu thuong', 'LIABILITY', FALSE, TRUE),
    ('3432', 'Trai phieu chuyen doi', 'LIABILITY', FALSE, TRUE),
    ('344', 'Nhan ky quy ky cuoc', 'LIABILITY', FALSE, TRUE),
    ('347', 'Thue TN hoan lai phai tra', 'LIABILITY', FALSE, TRUE),
    ('352', 'Du phong phai tra', 'LIABILITY', TRUE, TRUE),
    ('3521', 'Du phong BH san pham', 'LIABILITY', FALSE, TRUE),
    ('3522', 'Du phong BH xay lap', 'LIABILITY', FALSE, TRUE),
    ('3523', 'Du phong tai cau truc', 'LIABILITY', FALSE, TRUE),
    ('3525', 'Du phong khac', 'LIABILITY', FALSE, TRUE),
    ('353', 'Quy khen thuong phuc loi', 'LIABILITY', FALSE, TRUE),

    -- LOAI VON CHU SO HUU (EQUITY)
    ('4', 'Von chu so huu', 'EQUITY', TRUE, TRUE),
    ('411', 'Von dau tu cua CSH', 'EQUITY', TRUE, TRUE),
    ('4111', 'Von co phan', 'EQUITY', FALSE, TRUE),
    ('4112', 'Thang du von', 'EQUITY', FALSE, TRUE),
    ('4113', 'Von khac', 'EQUITY', FALSE, TRUE),
    ('412', 'Chenh lech danh gia lai TS', 'EQUITY', FALSE, TRUE),
    ('413', 'Chenh lech ty gia', 'EQUITY', FALSE, TRUE),
    ('414', 'Quy dau tu phat trien', 'EQUITY', FALSE, TRUE),
    ('415', 'Quy du tru bo sung VDL', 'EQUITY', FALSE, TRUE),
    ('417', 'Quy khac', 'EQUITY', FALSE, TRUE),
    ('418', 'Cac quy KHKT', 'EQUITY', FALSE, TRUE),
    ('419', 'Co phieu quy', 'EQUITY', FALSE, TRUE),
    ('421', 'Loi nhuan sau thue', 'EQUITY', TRUE, TRUE),
    ('4211', 'Loi nhuan nam truoc', 'EQUITY', FALSE, TRUE),
    ('4212', 'Loi nhuan nam nay', 'EQUITY', FALSE, TRUE),

    -- LOAI DOANH THU (REVENUE)
    ('5', 'Doanh thu', 'REVENUE', TRUE, TRUE),
    ('511', 'Doanh thu ban hang', 'REVENUE', TRUE, TRUE),
    ('5111', 'DT ban hang hoa', 'REVENUE', FALSE, TRUE),
    ('5112', 'DT cung cap dich vu', 'REVENUE', FALSE, TRUE),
    ('515', 'DT hoat dong tai chinh', 'REVENUE', FALSE, TRUE),
    ('521', 'Cac khoan giam tru DT', 'REVENUE', TRUE, TRUE),
    ('5211', 'Chiet khau thuong mai', 'REVENUE', FALSE, TRUE),
    ('5212', 'Hang ban bi tra lai', 'REVENUE', FALSE, TRUE),
    ('5213', 'Giam gia hang ban', 'REVENUE', FALSE, TRUE),
    ('711', 'Thu nhap khac', 'REVENUE', FALSE, TRUE),

    -- LOAI CHI PHI (EXPENSES)
    ('6', 'Chi phi SXKD', 'EXPENSE', TRUE, TRUE),
    ('621', 'Chi phi NVL truc tiep', 'EXPENSE', FALSE, TRUE),
    ('622', 'Chi phi nhan cong truc tiep', 'EXPENSE', FALSE, TRUE),
    ('623', 'Chi phi su dung may thi cong', 'EXPENSE', TRUE, TRUE),
    ('6231', 'Chi phi nhan cong', 'EXPENSE', FALSE, TRUE),
    ('6232', 'Chi phi vat lieu', 'EXPENSE', FALSE, TRUE),
    ('6233', 'Chi phi dung cu', 'EXPENSE', FALSE, TRUE),
    ('6237', 'Chi phi khac', 'EXPENSE', FALSE, TRUE),
    ('627', 'Chi phi san xuat chung', 'EXPENSE', TRUE, TRUE),
    ('6271', 'Chi phi nhan vien', 'EXPENSE', FALSE, TRUE),
    ('6272', 'Chi phi vat lieu', 'EXPENSE', FALSE, TRUE),
    ('6273', 'Chi phi dung cu', 'EXPENSE', FALSE, TRUE),
    ('6274', 'Chi phi khau hao', 'EXPENSE', FALSE, TRUE),
    ('6275', 'Chi phi thue phi le phi', 'EXPENSE', FALSE, TRUE),
    ('6277', 'Chi phi DV mua ngoai', 'EXPENSE', FALSE, TRUE),
    ('6278', 'Chi phi khac bang tien', 'EXPENSE', FALSE, TRUE),
    ('632', 'Gia von hang ban', 'EXPENSE', FALSE, TRUE),
    ('635', 'Chi phi hoat dong TC', 'EXPENSE', FALSE, TRUE),
    ('641', 'Chi phi ban hang', 'EXPENSE', TRUE, TRUE),
    ('6411', 'CP ban hang - nhan vien', 'EXPENSE', FALSE, TRUE),
    ('6412', 'CP ban hang - vat lieu', 'EXPENSE', FALSE, TRUE),
    ('6413', 'CP ban hang - dung cu', 'EXPENSE', FALSE, TRUE),
    ('6414', 'CP ban hang - khau hao', 'EXPENSE', FALSE, TRUE),
    ('6415', 'CP ban hang - thue phi', 'EXPENSE', FALSE, TRUE),
    ('6417', 'CP ban hang - DV mua ngoai', 'EXPENSE', FALSE, TRUE),
    ('6418', 'CP ban hang - bang tien', 'EXPENSE', FALSE, TRUE),
    ('642', 'Chi phi quan ly DN', 'EXPENSE', TRUE, TRUE),
    ('6421', 'CP QLDN - nhan vien', 'EXPENSE', FALSE, TRUE),
    ('6422', 'CP QLDN - vat lieu', 'EXPENSE', FALSE, TRUE),
    ('6423', 'CP QLDN - dung cu', 'EXPENSE', FALSE, TRUE),
    ('6424', 'CP QLDN - khau hao', 'EXPENSE', FALSE, TRUE),
    ('6425', 'CP QLDN - thue phi', 'EXPENSE', FALSE, TRUE),
    ('6426', 'CP QLDN - du phong', 'EXPENSE', FALSE, TRUE),
    ('6427', 'CP QLDN - DV mua ngoai', 'EXPENSE', FALSE, TRUE),
    ('6428', 'CP QLDN - bang tien', 'EXPENSE', FALSE, TRUE),
    ('811', 'Chi phi khac', 'EXPENSE', FALSE, TRUE),

    -- LOAI XAC DINH KET QUA (DETERMINATION OF RESULTS)
    ('8', 'Xac dinh ket qua', 'EXPENSE', TRUE, TRUE),
    ('821', 'Chi phi thue TNDN', 'EXPENSE', TRUE, TRUE),
    ('8211', 'Chi phi thue TNDN hien hanh', 'EXPENSE', FALSE, TRUE),
    ('82111', 'CP thue TNDN hien hanh', 'EXPENSE', FALSE, TRUE),
    ('82112', 'CP thue TNDN toi thieu toan cau', 'EXPENSE', FALSE, TRUE),
    ('8212', 'CP thue TNDN hoan lai', 'EXPENSE', FALSE, TRUE),
    ('911', 'Xac dinh ket qua KD', 'EXPENSE', FALSE, TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================
-- SEED: Closing Entry Templates
-- ============================================================
WITH t AS (
    INSERT INTO closing_templates (id, name, description, sequence_order)
    VALUES
        (gen_random_uuid(), 'Ket chuyen doanh thu', 'Close revenue to P&L', 1),
        (gen_random_uuid(), 'Ket chuyen chi phi', 'Close expenses to P&L', 2),
        (gen_random_uuid(), 'Ket chuyen loi nhuan', 'Close P&L to retained earnings', 3)
    RETURNING id, name
)
INSERT INTO closing_template_lines (template_id, line_number, debit_account, credit_account, formula, description)
SELECT t.id, 1, '511', '911', 'BALANCE', 'Close revenue' FROM t WHERE t.name = 'Ket chuyen doanh thu'
UNION ALL
SELECT t.id, 2, '515', '911', 'BALANCE', 'Close financial income' FROM t WHERE t.name = 'Ket chuyen doanh thu'
UNION ALL
SELECT t.id, 3, '711', '911', 'BALANCE', 'Close other income' FROM t WHERE t.name = 'Ket chuyen doanh thu'
UNION ALL
SELECT t.id, 1, '911', '632', 'BALANCE', 'Close COGS' FROM t WHERE t.name = 'Ket chuyen chi phi'
UNION ALL
SELECT t.id, 2, '911', '635', 'BALANCE', 'Close financial cost' FROM t WHERE t.name = 'Ket chuyen chi phi'
UNION ALL
SELECT t.id, 3, '911', '641', 'BALANCE', 'Close selling expenses' FROM t WHERE t.name = 'Ket chuyen chi phi'
UNION ALL
SELECT t.id, 4, '911', '642', 'BALANCE', 'Close admin expenses' FROM t WHERE t.name = 'Ket chuyen chi phi'
UNION ALL
SELECT t.id, 5, '911', '811', 'BALANCE', 'Close other expenses' FROM t WHERE t.name = 'Ket chuyen chi phi'
UNION ALL
SELECT t.id, 6, '911', '821', 'BALANCE', 'Close CIT expense' FROM t WHERE t.name = 'Ket chuyen chi phi'
UNION ALL
SELECT t.id, 1, '911', '421', 'BALANCE', 'Close profit to retained earnings' FROM t WHERE t.name = 'Ket chuyen loi nhuan';

-- ============================================================
-- SEED: Default admin user (password: admin123 — change on first login)
-- Note: password hash is bcrypt of 'admin123'
-- ============================================================
INSERT INTO users (username, password_hash, full_name, role)
VALUES ('admin', '$2a$10$bE87GqCBB5/QMO.grdBik.lS4jna0wVGid3B2qdiNc1o6V4Ljl4Ay', 'Administrator', 'admin')
ON CONFLICT (username) DO NOTHING;

COMMIT;