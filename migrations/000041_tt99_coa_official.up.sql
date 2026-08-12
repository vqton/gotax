-- ============================================================
-- 000041: Align COA seed to official TT99 Phụ lục II
-- Reference: Circular 99/2025/TT-BTC (effective 01/01/2026),
--   Appendix II - Enterprise Chart of Accounts (issued 27/10/2025).
-- What this migration does:
--   1. Inserts official TT99 accounts missing from the 000001 seed:
--      1362, 1363, 3362, 3363, 3531-3534, 356, 3561, 3562, 357,
--      4118, 41111, 41112, 6234, 6238, 21511, 21512, 215121,
--      215122 (TT99 added funds/loan-cost detail + biological asset depth)
--   2. Renumbers 3385 -> 3386 for BHTN (TT99 moved social insurance
--      unemployment to 3386; 3385 no longer exists in Appendix II)
--   3. Rewrites account names to official diacritic Vietnamese per
--      Appendix II. Non-diacritic names fail audit/tax inspection.
--   4. Fixes is_parent flags for accounts that gained children.
--   5. Inserts 5113 as enterprise-detail (TT99 Art. 11) — NOT in
--      Appendix II (TT99 lists 511/521 without cấp 2 detail), but
--      required by tax_service.go vatRevenueAccounts.
-- Kept (enterprise-detail per TT99 Art. 11, used by app services):
--   1111/1112/1121/1122/1211/1212/1311/1312/1385/2111-2118/2131-2135/
--   3311/3312/415/417/5111/5112/5211-5213 — legal as self-added detail;
--   deleting would break PG FK references in journal/import workflows.
-- Note: TT99 has NO Loại 0 (off-balance 001-009) — eliminated from
--   Appendix II. Do not re-add.
-- ============================================================

-- ─── 1. Insert missing official TT99 accounts ─────────────────────
INSERT INTO accounts (code, name, type, is_parent, is_active) VALUES
    -- TK 1362/1363 - intra-company receivables (TT99 new detail)
    ('1362', 'Phải thu nội bộ về chênh lệch tỷ giá', 'ASSET', FALSE, TRUE),
    ('1363', 'Phải thu nội bộ về chi phí đi vay đủ điều kiện được vốn hoá', 'ASSET', FALSE, TRUE),
    -- TK 2151 - biological asset maturity detail (cấp 3/4)
    ('21511', 'Súc vật nuôi cho sản phẩm định kỳ chưa đạt đến giai đoạn trưởng thành', 'ASSET', FALSE, TRUE),
    ('21512', 'Súc vật nuôi cho sản phẩm định kỳ đạt đến giai đoạn trưởng thành', 'ASSET', TRUE, TRUE),
    ('215121', 'Nguyên giá', 'ASSET', FALSE, TRUE),
    ('215122', 'Giá trị khấu hao lũy kế', 'ASSET', FALSE, TRUE),
    -- TK 3362/3363 - intra-company payables (TT99 new detail)
    ('3362', 'Phải trả nội bộ về chênh lệch tỷ giá', 'LIABILITY', FALSE, TRUE),
    ('3363', 'Phải trả nội bộ về chi phí đi vay đủ điều kiện được vốn hoá', 'LIABILITY', FALSE, TRUE),
    -- TK 353 - bonus & welfare fund detail (TT99 lists 3531-3534)
    ('3531', 'Quỹ khen thưởng', 'LIABILITY', FALSE, TRUE),
    ('3532', 'Quỹ phúc lợi', 'LIABILITY', FALSE, TRUE),
    ('3533', 'Quỹ phúc lợi đã hình thành TSCĐ', 'LIABILITY', FALSE, TRUE),
    ('3534', 'Quỹ thưởng ban quản lý điều hành công ty', 'LIABILITY', FALSE, TRUE),
    -- TK 356 - science & technology fund (TT99 keeps TT200 356)
    ('356', 'Quỹ phát triển khoa học và công nghệ', 'LIABILITY', TRUE, TRUE),
    ('3561', 'Quỹ phát triển khoa học và công nghệ', 'LIABILITY', FALSE, TRUE),
    ('3562', 'Quỹ phát triển khoa học và công nghệ đã hình thành tài sản', 'LIABILITY', FALSE, TRUE),
    -- TK 357 - price stabilization fund (TT99 keeps TT200 357)
    ('357', 'Quỹ bình ổn giá', 'LIABILITY', FALSE, TRUE),
    -- TK 4111/4118 - equity detail
    ('41111', 'Cổ phiếu phổ thông có quyền biểu quyết', 'EQUITY', FALSE, TRUE),
    ('41112', 'Cổ phiếu ưu đãi', 'EQUITY', FALSE, TRUE),
    ('4118', 'Vốn khác', 'EQUITY', FALSE, TRUE),
    -- TK 623 - construction machine cost detail
    ('6234', 'Chi phí khấu hao máy thi công', 'EXPENSE', FALSE, TRUE),
    ('6238', 'Chi phí bằng tiền khác', 'EXPENSE', FALSE, TRUE)
ON CONFLICT (code) DO NOTHING;

-- ─── 1b. Enterprise-detail insert (TT99 Art. 11) ──────────────────
-- 5113 absent from Appendix II (TT99 has no 511 cấp 2), but
-- tax_service.go vatRevenueAccounts hardcodes it.
INSERT INTO accounts (code, name, type, is_parent, is_active)
VALUES ('5113', 'Doanh thu cung cấp dịch vụ', 'REVENUE', FALSE, TRUE)
ON CONFLICT (code) DO NOTHING;

-- ─── 2. Renumber 3385 -> 3386 (BHTN) ──────────────────────────────
-- 3385 is unused by app code. Guard against any referenced row.
UPDATE accounts
   SET code = '3386', name = 'Bảo hiểm thất nghiệp'
 WHERE code = '3385'
   AND NOT EXISTS (SELECT 1 FROM journal_lines WHERE account_code = '3385')
   AND NOT EXISTS (SELECT 1 FROM account_balances WHERE account_code = '3385')
   AND NOT EXISTS (SELECT 1 FROM account_analysis WHERE account_code = '3385')
   AND NOT EXISTS (SELECT 1 FROM closing_template_lines
                   WHERE debit_account = '3385' OR credit_account = '3385');

-- If 3385 was absent or guard blocked the renumber, ensure 3386 exists.
INSERT INTO accounts (code, name, type, is_parent, is_active)
VALUES ('3386', 'Bảo hiểm thất nghiệp', 'LIABILITY', FALSE, TRUE)
ON CONFLICT (code) DO NOTHING;

-- ─── 3. Fix is_parent for accounts that gained children ───────────
UPDATE accounts SET is_parent = TRUE WHERE code IN ('2151', '353', '4111');
-- 356 inserted as parent above; 21512 inserted as parent above.

-- ─── 4. Rewrite names to official diacritic Vietnamese ────────────
UPDATE accounts SET name = CASE code
    -- Loại 1-2: TÀI SẢN
    WHEN '1' THEN 'Tài sản'
    WHEN '111' THEN 'Tiền mặt'
    WHEN '1111' THEN 'Tiền mặt Việt Nam đồng'
    WHEN '1112' THEN 'Tiền mặt ngoại tệ'
    WHEN '112' THEN 'Tiền gửi không kỳ hạn'
    WHEN '1121' THEN 'Tiền gửi Việt Nam đồng'
    WHEN '1122' THEN 'Tiền gửi ngoại tệ'
    WHEN '113' THEN 'Tiền đang chuyển'
    WHEN '121' THEN 'Chứng khoán kinh doanh'
    WHEN '1211' THEN 'Cổ phiếu'
    WHEN '1212' THEN 'Trái phiếu'
    WHEN '128' THEN 'Đầu tư nắm giữ đến ngày đáo hạn'
    WHEN '1281' THEN 'Tiền gửi có kỳ hạn'
    WHEN '1282' THEN 'Trái phiếu'
    WHEN '1283' THEN 'Cho vay'
    WHEN '1288' THEN 'Các khoản đầu tư khác nắm giữ đến ngày đáo hạn'
    WHEN '131' THEN 'Phải thu của khách hàng'
    WHEN '1311' THEN 'Phải thu của khách hàng trong nước'
    WHEN '1312' THEN 'Phải thu của khách hàng nước ngoài'
    WHEN '133' THEN 'Thuế GTGT được khấu trừ'
    WHEN '1331' THEN 'Thuế GTGT được khấu trừ của hàng hóa, dịch vụ'
    WHEN '1332' THEN 'Thuế GTGT được khấu trừ của TSCĐ'
    WHEN '136' THEN 'Phải thu nội bộ'
    WHEN '1361' THEN 'Vốn kinh doanh ở đơn vị trực thuộc'
    WHEN '1368' THEN 'Phải thu nội bộ khác'
    WHEN '138' THEN 'Phải thu khác'
    WHEN '1381' THEN 'Tài sản thiếu chờ xử lý'
    WHEN '1383' THEN 'Thuế TTĐB của hàng nhập khẩu'
    WHEN '1385' THEN 'Phải thu về cho mượn TSCĐ'
    WHEN '1388' THEN 'Phải thu khác'
    WHEN '141' THEN 'Tạm ứng'
    WHEN '151' THEN 'Hàng mua đang đi đường'
    WHEN '152' THEN 'Nguyên liệu, vật liệu'
    WHEN '153' THEN 'Công cụ, dụng cụ'
    WHEN '154' THEN 'Chi phí sản xuất, kinh doanh dở dang'
    WHEN '155' THEN 'Sản phẩm'
    WHEN '156' THEN 'Hàng hóa'
    WHEN '157' THEN 'Hàng gửi đi bán'
    WHEN '158' THEN 'Nguyên liệu, vật tư tại kho bảo thuế'
    WHEN '171' THEN 'Giao dịch mua, bán lại trái phiếu chính phủ'
    WHEN '211' THEN 'Tài sản cố định hữu hình'
    WHEN '2111' THEN 'Nhà cửa, vật kiến trúc'
    WHEN '2112' THEN 'Máy móc, thiết bị'
    WHEN '2113' THEN 'Phương tiện vận tải'
    WHEN '2114' THEN 'Thiết bị quản lý'
    WHEN '2115' THEN 'Cây lâu năm'
    WHEN '2118' THEN 'TSCĐ khác'
    WHEN '212' THEN 'Tài sản cố định thuê tài chính'
    WHEN '213' THEN 'Tài sản cố định vô hình'
    WHEN '2131' THEN 'Quyền sử dụng đất'
    WHEN '2132' THEN 'Bản quyền, bằng sáng chế'
    WHEN '2133' THEN 'Nhãn hiệu thương mại'
    WHEN '2134' THEN 'Phần mềm máy tính'
    WHEN '2135' THEN 'Giá trị lợi thế'
    WHEN '214' THEN 'Hao mòn tài sản cố định'
    WHEN '2141' THEN 'Hao mòn TSCĐ hữu hình'
    WHEN '2142' THEN 'Hao mòn TSCĐ thuê tài chính'
    WHEN '2143' THEN 'Hao mòn TSCĐ vô hình'
    WHEN '2147' THEN 'Hao mòn bất động sản đầu tư'
    WHEN '215' THEN 'Tài sản sinh học'
    WHEN '2151' THEN 'Súc vật nuôi cho sản phẩm định kỳ'
    WHEN '2152' THEN 'Súc vật nuôi lấy sản phẩm một lần'
    WHEN '2153' THEN 'Cây trồng theo mùa vụ hoặc lấy sản phẩm một lần'
    WHEN '217' THEN 'Bất động sản đầu tư'
    WHEN '221' THEN 'Đầu tư vào công ty con'
    WHEN '222' THEN 'Đầu tư vào công ty liên doanh, liên kết'
    WHEN '228' THEN 'Đầu tư khác'
    WHEN '2281' THEN 'Đầu tư góp vốn vào đơn vị khác'
    WHEN '2288' THEN 'Đầu tư khác'
    WHEN '229' THEN 'Dự phòng tổn thất tài sản'
    WHEN '2291' THEN 'Dự phòng giảm giá chứng khoán kinh doanh'
    WHEN '2292' THEN 'Dự phòng tổn thất đầu tư vào đơn vị khác'
    WHEN '2293' THEN 'Dự phòng phải thu khó đòi'
    WHEN '2294' THEN 'Dự phòng giảm giá hàng tồn kho'
    WHEN '2295' THEN 'Dự phòng tổn thất tài sản sinh học'
    WHEN '241' THEN 'Xây dựng cơ bản dở dang'
    WHEN '2411' THEN 'Mua sắm TSCĐ'
    WHEN '2412' THEN 'Xây dựng cơ bản'
    WHEN '2413' THEN 'Sửa chữa, bảo dưỡng định kỳ TSCĐ'
    WHEN '2414' THEN 'Nâng cấp, cải tạo TSCĐ'
    WHEN '242' THEN 'Chi phí chờ phân bổ'
    WHEN '243' THEN 'Tài sản thuế thu nhập hoãn lại'
    WHEN '244' THEN 'Ký quỹ, ký cược'
    -- Loại 3: NỢ PHẢI TRẢ
    WHEN '3' THEN 'Nợ phải trả'
    WHEN '331' THEN 'Phải trả cho người bán'
    WHEN '3311' THEN 'Phải trả cho người bán trong nước'
    WHEN '3312' THEN 'Phải trả cho người bán nước ngoài'
    WHEN '332' THEN 'Phải trả cổ tức, lợi nhuận'
    WHEN '333' THEN 'Thuế và các khoản phải nộp Nhà nước'
    WHEN '3331' THEN 'Thuế giá trị gia tăng phải nộp'
    WHEN '33311' THEN 'Thuế GTGT đầu ra'
    WHEN '33312' THEN 'Thuế GTGT hàng nhập khẩu'
    WHEN '3332' THEN 'Thuế tiêu thụ đặc biệt'
    WHEN '3333' THEN 'Thuế xuất, nhập khẩu'
    WHEN '3334' THEN 'Thuế thu nhập doanh nghiệp'
    WHEN '3335' THEN 'Thuế thu nhập cá nhân'
    WHEN '3336' THEN 'Thuế tài nguyên'
    WHEN '3337' THEN 'Thuế nhà đất, tiền thuê đất'
    WHEN '3338' THEN 'Thuế bảo vệ môi trường và các loại thuế khác'
    WHEN '33381' THEN 'Thuế bảo vệ môi trường'
    WHEN '33382' THEN 'Các loại thuế khác'
    WHEN '3339' THEN 'Phí, lệ phí và các khoản phải nộp khác'
    WHEN '334' THEN 'Phải trả người lao động'
    WHEN '335' THEN 'Chi phí phải trả'
    WHEN '336' THEN 'Phải trả nội bộ'
    WHEN '3361' THEN 'Phải trả nội bộ về vốn kinh doanh'
    WHEN '3368' THEN 'Phải trả nội bộ khác'
    WHEN '337' THEN 'Thanh toán theo tiến độ hợp đồng xây dựng'
    WHEN '338' THEN 'Phải trả, phải nộp khác'
    WHEN '3381' THEN 'Tài sản thừa chờ giải quyết'
    WHEN '3382' THEN 'Kinh phí công đoàn'
    WHEN '3383' THEN 'Bảo hiểm xã hội'
    WHEN '3384' THEN 'Bảo hiểm y tế'
    WHEN '3387' THEN 'Doanh thu chờ phân bổ'
    WHEN '3388' THEN 'Phải trả, phải nộp khác'
    WHEN '341' THEN 'Vay và nợ thuê tài chính'
    WHEN '3411' THEN 'Các khoản đi vay'
    WHEN '3412' THEN 'Nợ thuê tài chính'
    WHEN '343' THEN 'Trái phiếu phát hành'
    WHEN '3431' THEN 'Trái phiếu thường'
    WHEN '3432' THEN 'Trái phiếu chuyển đổi'
    WHEN '344' THEN 'Nhận ký quỹ, ký cược'
    WHEN '347' THEN 'Thuế thu nhập hoãn lại phải trả'
    WHEN '352' THEN 'Dự phòng phải trả'
    WHEN '3521' THEN 'Dự phòng bảo hành sản phẩm, hàng hóa'
    WHEN '3522' THEN 'Dự phòng bảo hành công trình xây dựng'
    WHEN '3523' THEN 'Dự phòng tái cơ cấu doanh nghiệp'
    WHEN '3525' THEN 'Dự phòng phải trả khác'
    WHEN '353' THEN 'Quỹ khen thưởng, phúc lợi'
    -- Loại 4: VỐN CHỦ SỞ HỮU
    WHEN '4' THEN 'Vốn chủ sở hữu'
    WHEN '411' THEN 'Vốn đầu tư của chủ sở hữu'
    WHEN '4111' THEN 'Vốn góp của chủ sở hữu'
    WHEN '4112' THEN 'Thặng dư vốn'
    WHEN '4113' THEN 'Quyền chọn chuyển đổi trái phiếu'
    WHEN '412' THEN 'Chênh lệch đánh giá lại tài sản'
    WHEN '413' THEN 'Chênh lệch tỷ giá hối đoái'
    WHEN '414' THEN 'Quỹ đầu tư phát triển'
    WHEN '415' THEN 'Quỹ dự trữ bổ sung vốn điều lệ'
    WHEN '417' THEN 'Quỹ khác'
    WHEN '418' THEN 'Các quỹ khác thuộc vốn chủ sở hữu'
    WHEN '419' THEN 'Cổ phiếu mua lại của chính mình'
    WHEN '421' THEN 'Lợi nhuận sau thuế chưa phân phối'
    WHEN '4211' THEN 'Lợi nhuận sau thuế chưa phân phối lũy kế đến cuối năm trước'
    WHEN '4212' THEN 'Lợi nhuận sau thuế chưa phân phối năm nay'
    -- Loại 5: DOANH THU
    WHEN '5' THEN 'Doanh thu'
    WHEN '511' THEN 'Doanh thu bán hàng và cung cấp dịch vụ'
    WHEN '5111' THEN 'Doanh thu bán hàng hóa'
    WHEN '5112' THEN 'Doanh thu bán thành phẩm'
    WHEN '515' THEN 'Doanh thu hoạt động tài chính'
    WHEN '521' THEN 'Các khoản giảm trừ doanh thu'
    WHEN '5211' THEN 'Chiết khấu thương mại'
    WHEN '5212' THEN 'Hàng bán bị trả lại'
    WHEN '5213' THEN 'Giảm giá hàng bán'
    WHEN '711' THEN 'Thu nhập khác'
    -- Loại 6: CHI PHÍ SẢN XUẤT, KINH DOANH
    WHEN '6' THEN 'Chi phí sản xuất, kinh doanh'
    WHEN '621' THEN 'Chi phí nguyên liệu, vật liệu trực tiếp'
    WHEN '622' THEN 'Chi phí nhân công trực tiếp'
    WHEN '623' THEN 'Chi phí sử dụng máy thi công'
    WHEN '6231' THEN 'Chi phí nhân công'
    WHEN '6232' THEN 'Chi phí vật liệu'
    WHEN '6233' THEN 'Chi phí dụng cụ sản xuất'
    WHEN '6237' THEN 'Chi phí dịch vụ mua ngoài'
    WHEN '627' THEN 'Chi phí sản xuất chung'
    WHEN '6271' THEN 'Chi phí nhân viên phân xưởng'
    WHEN '6272' THEN 'Chi phí vật liệu'
    WHEN '6273' THEN 'Chi phí dụng cụ sản xuất'
    WHEN '6274' THEN 'Chi phí khấu hao TSCĐ'
    WHEN '6275' THEN 'Thuế, phí, lệ phí'
    WHEN '6277' THEN 'Chi phí dịch vụ mua ngoài'
    WHEN '6278' THEN 'Chi phí bằng tiền khác'
    WHEN '632' THEN 'Giá vốn hàng bán'
    WHEN '635' THEN 'Chi phí tài chính'
    WHEN '641' THEN 'Chi phí bán hàng'
    WHEN '6411' THEN 'Chi phí nhân viên'
    WHEN '6412' THEN 'Chi phí vật liệu, bao bì'
    WHEN '6413' THEN 'Chi phí dụng cụ, đồ dùng'
    WHEN '6414' THEN 'Chi phí khấu hao TSCĐ'
    WHEN '6415' THEN 'Thuế, phí, lệ phí'
    WHEN '6417' THEN 'Chi phí dịch vụ mua ngoài'
    WHEN '6418' THEN 'Chi phí bằng tiền khác'
    WHEN '642' THEN 'Chi phí quản lý doanh nghiệp'
    WHEN '6421' THEN 'Chi phí nhân viên quản lý'
    WHEN '6422' THEN 'Chi phí vật liệu quản lý'
    WHEN '6423' THEN 'Chi phí đồ dùng văn phòng'
    WHEN '6424' THEN 'Chi phí khấu hao TSCĐ'
    WHEN '6425' THEN 'Thuế, phí và lệ phí'
    WHEN '6426' THEN 'Chi phí dự phòng'
    WHEN '6427' THEN 'Chi phí dịch vụ mua ngoài'
    WHEN '6428' THEN 'Chi phí bằng tiền khác'
    -- Loại 8-9: CHI PHÍ KHÁC / XÁC ĐỊNH KẾT QUẢ
    WHEN '8' THEN 'Xác định kết quả kinh doanh'
    WHEN '811' THEN 'Chi phí khác'
    WHEN '821' THEN 'Chi phí thuế thu nhập doanh nghiệp'
    WHEN '8211' THEN 'Chi phí thuế TNDN hiện hành'
    WHEN '82111' THEN 'Chi phí thuế thu nhập doanh nghiệp hiện hành theo quy định của Luật thuế thu nhập doanh nghiệp'
    WHEN '82112' THEN 'Chi phí thuế thu nhập doanh nghiệp bổ sung theo quy định về thuế tối thiểu toàn cầu'
    WHEN '8212' THEN 'Chi phí thuế TNDN hoãn lại'
    WHEN '911' THEN 'Xác định kết quả kinh doanh'
    ELSE name
  END
WHERE code IN (
    '1','111','1111','1112','112','1121','1122','113','121','1211','1212',
    '128','1281','1282','1283','1288','131','1311','1312','133','1331','1332',
    '136','1361','1368','138','1381','1383','1385','1388','141','151','152','153',
    '154','155','156','157','158','171','211','2111','2112','2113','2114','2115',
    '2118','212','213','2131','2132','2133','2134','2135','214','2141','2142','2143',
    '2147','215','2151','2152','2153','217','221','222','228','2281','2288','229',
    '2291','2292','2293','2294','2295','241','2411','2412','2413','2414','242','243','244',
    '3','331','3311','3312','332','333','3331','33311','33312','3332','3333','3334',
    '3335','3336','3337','3338','33381','33382','3339','334','335','336','3361','3368',
    '337','338','3381','3382','3383','3384','3387','3388','341','3411','3412','343',
    '3431','3432','344','347','352','3521','3522','3523','3525','353',
    '4','411','4111','4112','4113','412','413','414','415','417','418','419','421',
    '4211','4212',
    '5','511','5111','5112','515','521','5211','5212','5213','711',
    '6','621','622','623','6231','6232','6233','6237','627','6271','6272','6273',
    '6274','6275','6277','6278','632','635','641','6411','6412','6413','6414','6415',
    '6417','6418','642','6421','6422','6423','6424','6425','6426','6427','6428',
    '8','811','821','8211','82111','82112','8212','911'
);
