-- ============================================================
-- 000045: Complete loại headers per TT99 Appendix II
-- Reference: Circular 99/2025/TT-BTC Appendix II
-- Appendix II groups accounts under 9 loại (categories):
--   Loại 1 - Tài sản ngắn hạn   (111..158, 18 accounts)
--   Loại 2 - Tài sản dài hạn    (171, 211..244, 15 accounts)
--   Loại 3..9 as previously seeded
-- Seed created loại rows only for families that had one; loại 2
-- was missing and its 15 accounts were wrongly parented to loại 1.
-- This migration creates loại 2, reparents 171/211..244, and
-- restores the official loại 1 name.
-- ============================================================

-- ─── 1. Create missing loại 2 ────────────────────────────────────────
INSERT INTO accounts (code, name, type, is_parent, is_active, status)
SELECT '2', 'Tài sản dài hạn', 'ASSET', TRUE, TRUE, 'ACTIVE'
WHERE NOT EXISTS (SELECT 1 FROM accounts WHERE code = '2');

-- ─── 2. Reparent long-term accounts from loại 1 to loại 2 ───────────
UPDATE accounts SET parent_code = '2'
 WHERE parent_code = '1' AND code IN
 ('171','211','212','213','214','215','217','221','222','228','229',
  '241','242','243','244');

-- ─── 3. Restore official loại 1 name ────────────────────────────────
UPDATE accounts SET name = 'Tài sản ngắn hạn' WHERE code = '1' AND name = 'Tài sản';
