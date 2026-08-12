-- ============================================================
-- 000041 DOWN: Revert COA alignment to TT99 Phụ lục II
-- Reverses structure only (account names stay diacritic — reverting
-- names has no audit value and risks breaking running integrations).
-- ============================================================

-- ─── 1. Delete inserted accounts (children before parents; FK on parent_code) ─
-- Guard: never delete an account referenced by journaling data.
DELETE FROM accounts
 WHERE code IN ('215121', '215122')
   AND NOT EXISTS (SELECT 1 FROM journal_lines WHERE account_code IN ('215121', '215122'))
   AND NOT EXISTS (SELECT 1 FROM account_balances WHERE account_code IN ('215121', '215122'))
   AND NOT EXISTS (SELECT 1 FROM account_analysis WHERE account_code IN ('215121', '215122'));

DELETE FROM accounts
 WHERE code IN ('21511', '21512')
   AND NOT EXISTS (SELECT 1 FROM journal_lines WHERE account_code IN ('21511', '21512'))
   AND NOT EXISTS (SELECT 1 FROM account_balances WHERE account_code IN ('21511', '21512'))
   AND NOT EXISTS (SELECT 1 FROM account_analysis WHERE account_code IN ('21511', '21512'));

DELETE FROM accounts
 WHERE code IN ('1362', '1363', '3362', '3363', '3531', '3532', '3533', '3534',
                '3561', '3562', '357', '41111', '41112', '4118',
                '6234', '6238', '5113')
   AND NOT EXISTS (SELECT 1 FROM journal_lines WHERE account_code IN ('1362', '1363', '3362', '3363', '3531', '3532', '3533', '3534', '3561', '3562', '357', '41111', '41112', '4118', '6234', '6238', '5113'))
   AND NOT EXISTS (SELECT 1 FROM account_balances WHERE account_code IN ('1362', '1363', '3362', '3363', '3531', '3532', '3533', '3534', '3561', '3562', '357', '41111', '41112', '4118', '6234', '6238', '5113'))
   AND NOT EXISTS (SELECT 1 FROM account_analysis WHERE account_code IN ('1362', '1363', '3362', '3363', '3531', '3532', '3533', '3534', '3561', '3562', '357', '41111', '41112', '4118', '6234', '6238', '5113'));

-- 356 is parent of 3561/3562 — delete only after children removed.
DELETE FROM accounts
 WHERE code = '356'
   AND NOT EXISTS (SELECT 1 FROM journal_lines WHERE account_code = '356')
   AND NOT EXISTS (SELECT 1 FROM account_balances WHERE account_code = '356')
   AND NOT EXISTS (SELECT 1 FROM account_analysis WHERE account_code = '356');

-- ─── 2. Revert is_parent for accounts whose children were removed ──
UPDATE accounts SET is_parent = FALSE WHERE code IN ('2151', '353', '4111');

-- ─── 3. Revert 3386 -> 3385 (BHTN) ────────────────────────────────
-- If 3385 survived (up-guard blocked renumber), just drop 3386.
DELETE FROM accounts
 WHERE code = '3386'
   AND EXISTS (SELECT 1 FROM accounts WHERE code = '3385')
   AND NOT EXISTS (SELECT 1 FROM journal_lines WHERE account_code = '3386')
   AND NOT EXISTS (SELECT 1 FROM account_balances WHERE account_code = '3386')
   AND NOT EXISTS (SELECT 1 FROM account_analysis WHERE account_code = '3386');

-- Otherwise move 3386 back to 3385.
UPDATE accounts
   SET code = '3385', name = 'Bao hiem that nghiep'
 WHERE code = '3386'
   AND NOT EXISTS (SELECT 1 FROM accounts WHERE code = '3385')
   AND NOT EXISTS (SELECT 1 FROM journal_lines WHERE account_code = '3386')
   AND NOT EXISTS (SELECT 1 FROM account_balances WHERE account_code = '3386')
   AND NOT EXISTS (SELECT 1 FROM account_analysis WHERE account_code = '3386');
