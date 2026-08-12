-- ============================================================
-- 000044: Remove non-TT99 cấp 1 accounts; fix 1385 name
-- Reference: Circular 99/2025/TT-BTC Appendix II (effective 01/01/2026)
-- Appendix II defines the closed cấp 1 list (71 accounts). Per Điều 11
-- enterprises may open detail (cấp 2+) freely but NOT invent cấp 1.
--   415 "Quỹ dự trữ bổ sung vốn điều lệ" — TT200 legacy, NOT in Appendix II
--   417 "Quỹ khác"                       — TT200 legacy, NOT in Appendix II
-- Transition per TT99 Điều 29/Điều 35: balances map to 418 (các quỹ khác)
-- or 4118 (vốn khác) at the enterprise's discretion. Rows unreferenced by
-- any GL data are deleted; referenced rows are deactivated (kept for
-- historical data, hidden from active COA).
-- Also: 1385 kept as enterprise detail (Điều 11) with its standard name
-- restored — TT200 1385 = "Phải thu về cổ phần hóa" (TT99 retains
-- equitization accounting provisions, e.g. Điều 35/38/45).
-- ============================================================

-- ─── 1. Drop unreferenced 415/417 ────────────────────────────────────
DELETE FROM accounts
 WHERE code IN ('415', '417')
   AND NOT EXISTS (SELECT 1 FROM journal_lines WHERE account_code IN ('415', '417'))
   AND NOT EXISTS (SELECT 1 FROM account_balances WHERE account_code IN ('415', '417'))
   AND NOT EXISTS (SELECT 1 FROM account_analysis WHERE account_code IN ('415', '417'))
   AND NOT EXISTS (SELECT 1 FROM closing_template_lines
                   WHERE debit_account IN ('415', '417') OR credit_account IN ('415', '417'));

-- ─── 2. Deactivate any referenced leftovers ───────────────────────────
UPDATE accounts
   SET is_active = FALSE, status = 'INACTIVE'
 WHERE code IN ('415', '417') AND is_active = TRUE;

-- ─── 3. Restore standard 1385 name ───────────────────────────────────
UPDATE accounts SET name = 'Phải thu về cổ phần hóa' WHERE code = '1385';
