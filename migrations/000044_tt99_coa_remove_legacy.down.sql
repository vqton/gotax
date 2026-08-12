-- ============================================================
-- 000044 down: restore TT200-legacy 415/417 + prior 1385 name
-- ============================================================

UPDATE accounts SET name = 'Phải thu về cho mượn TSCĐ' WHERE code = '1385';

INSERT INTO accounts (code, name, type, is_parent, is_active)
VALUES ('415', 'Quỹ dự trữ bổ sung vốn điều lệ', 'EQUITY', FALSE, TRUE),
       ('417', 'Quỹ khác', 'EQUITY', FALSE, TRUE)
ON CONFLICT (code) DO NOTHING;
