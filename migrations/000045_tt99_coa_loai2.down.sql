-- ============================================================
-- 000045 down: restore pre-migration loại layout
-- ============================================================

UPDATE accounts SET name = 'Tài sản' WHERE code = '1' AND name = 'Tài sản ngắn hạn';

UPDATE accounts SET parent_code = '1'
 WHERE parent_code = '2' AND code IN
 ('171','211','212','213','214','215','217','221','222','228','229',
  '241','242','243','244');

DELETE FROM accounts WHERE code = '2';
