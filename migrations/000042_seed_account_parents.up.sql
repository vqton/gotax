-- 000042_seed_account_parents
--
-- Populate parent_code for the seeded COA. Migrations 000001-000041 never
-- set parent_code: the seed insert lists code/name/type/is_parent only, so
-- every account is a root and the account tree (UI, GetChildren guards,
-- delete cascade checks) is flat.
--
-- Rule: parent = longest OTHER existing account code that is a string
-- prefix (TT99 hierarchy cấp 1 → 2 → 3 → 4). Type rows '1'-'6' become the
-- parents of level-1 accounts, mirroring the MISA SME tree grouping.
-- Examples: 111→1, 1111→111, 21511→2151, 215121→21512, 4118→411, 5211→521.
-- Accounts created via API keep whatever parent_code the form set.

UPDATE accounts a
SET parent_code = p.code
FROM accounts p
WHERE p.code <> a.code
  AND a.code LIKE p.code || '%'
  AND p.code = (
    SELECT x.code FROM accounts x
    WHERE x.code <> a.code
      AND a.code LIKE x.code || '%'
    ORDER BY length(x.code) DESC, x.code
    LIMIT 1
  );

-- Level-1 accounts with no existing prefix row: group by account type
-- under the matching synthetic type row (seed has no '2'/'7'/'9' rows).
-- length(code) > 1 excludes the type rows themselves ('1','3','4','5','6','8')
-- which would otherwise self-parent.
-- 2xx long-term assets → '1' (Tài sản); 711 thu nhập khác → '5' (Doanh thu,
-- TT99 has no Loại 7, 711 belongs to Doanh thu); 911 XĐKQKD → '8'
-- (Xác định kết quả kinh doanh, same group as 821).
UPDATE accounts SET parent_code = '1' WHERE parent_code IS NULL AND length(code) > 1 AND type = 'ASSET';
UPDATE accounts SET parent_code = '5' WHERE parent_code IS NULL AND length(code) > 1 AND type = 'REVENUE';
UPDATE accounts SET parent_code = '8' WHERE parent_code IS NULL AND length(code) > 1 AND type = 'EXPENSE';
