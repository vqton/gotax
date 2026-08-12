-- Reverse 000042: restore pre-fix structure (711 under '5', 911 under '8',
-- '8' named "Xác định kết quả kinh doanh", drop categories 7 and 9).
-- Guarded so a 711/911 that moved on since does not get clobbered.
UPDATE accounts SET parent_code = '5' WHERE code = '711' AND parent_code = '7';
UPDATE accounts SET parent_code = '8' WHERE code = '911' AND parent_code = '9';
UPDATE accounts SET name = 'Xác định kết quả kinh doanh' WHERE code = '8';
DELETE FROM accounts WHERE code IN ('7', '9');
