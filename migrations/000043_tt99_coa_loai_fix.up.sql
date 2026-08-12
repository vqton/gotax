-- ============================================================
-- 000043: Fix COA category structure per TT99 Phụ lục II
-- Reference: Circular 99/2025/TT-BTC Appendix II (effective 01/01/2026)
-- Appendix II defines 8 LOẠI section headers (text, no codes):
--   TÀI SẢN / NỢ PHẢI TRẢ / VỐN CHỦ SỞ HỮU / DOANH THU /
--   CHI PHÍ SẢN XUẤT KINH DOANH / THU NHẬP KHÁC / CHI PHÍ KHÁC /
--   XÁC ĐỊNH KẾT QUẢ KINH DOANH
-- 000042 seeded parent_code with two wrong groupings: 711 under '5'
-- (Doanh thu) and 911 under '8' (claiming "TT99 has no Loại 7").
-- TT99 Appendix II DOES list LOẠI THU NHẬP KHÁC (711) and a separate
-- XÁC ĐỊNH KẾT QUẢ KINH DOANH section (911). Category '8' was also
-- misnamed "Xác định kết quả kinh doanh" — that is 911's name; LOẠI 8
-- (811, 821) is CHI PHÍ KHÁC.
-- Fixes: add missing LOẠI 7 + 9, rename '8', reparent 711 -> 7, 911 -> 9.
-- ============================================================

-- ─── 1. Add missing LOẠI 7 (Thu nhập khác) and LOẠI 9 (Xác định KQKD) ──
INSERT INTO accounts (code, name, type, is_parent, is_active) VALUES
    ('7', 'Thu nhập khác', 'REVENUE', TRUE, TRUE),
    ('9', 'Xác định kết quả kinh doanh', 'EXPENSE', TRUE, TRUE)
ON CONFLICT (code) DO NOTHING;

-- ─── 2. Fix '8' name: LOẠI 8 is Chi phí khác (811, 821) ───────────────
UPDATE accounts SET name = 'Chi phí khác' WHERE code = '8';

-- ─── 3. Reparent per Appendix II: 711 in LOẠI Thu nhập khác, ───────────
--     911 in LOẠI Xác định kết quả kinh doanh
UPDATE accounts SET parent_code = '7' WHERE code = '711';
UPDATE accounts SET parent_code = '9' WHERE code = '911';
