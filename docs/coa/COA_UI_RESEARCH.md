# COA UI/UX Research — Chart of Accounts Screen Design

Date: 2026-08-12
Purpose: evidence base for redesigning `web/app/coa.html` (GoTax Hệ thống tài khoản) for accountant ergonomics. Primary sources = vendor help centers + product docs; community sources tagged.

## Sources & findings

### MISA SME 2026 / MIMOSA / AMIS — Vietnamese market reference
Source: https://helpsme.misa.vn/2026/kb/khai-bao-tai-khoan/ (Khai báo tài khoản), https://helpmimosaonline.misa.vn/kb/he-thong-tai-khoan/ (Danh mục Hệ thống tài khoản)

- Entry: menu **Danh mục → Tài khoản → Hệ thống tài khoản**.
- Toolbar operations: **Thêm**, **Sửa** (double-click or Xem/Sửa arrow), **Ngừng theo dõi/Theo dõi**, **Nhân bản**, **Chuyển tài khoản hạch toán**, export. Save button = **Cất**.
- Required fields marked `*` ("thông tin có ký hiệu (*) kế toán sẽ bắt buộc phải khai báo").
- Foreign-currency flag: **Có hạch toán ngoại tệ** checkbox; 1112/1122/1132 forced on and cannot be unchecked.
- **Theo dõi chi tiết theo** (detail-by): Đối tượng (KH/NCC/NV), ngân hàng, khế ước vay, CTMT/DA, đơn đặt hàng, hợp đồng bán/mua, khoản mục CP, đơn vị, mã thống kê — each with **cảnh báo** (warn only) vs **bắt buộc nhập** (block) mode.
- Parent selection: new child account picks **TK tổng hợp** (parent). Posting to tổng hợp accounts forbidden.
- Safety: adding a child to an account **with phát sinh** prompts a warning — Yes moves all transactions to the new child, No aborts.
- Sibling-code rule: "1111 is child of 111, cannot add 11110 as sibling (prefix collision) — 111A allowed".
- Stop-tracking **cascades**: stopping a parent stops all children; stopping a child leaves parent active (follow means parent is followed by default).
- Opening balances must be deleted before toggling detail-by.
- MIMOSA R7 (https://helpmimosaonline.misa.vn/kb/ke-toan-muon-vua-xem-danh-muc-tai-khoan-vat-tu-hang-hoa-cong-cu-dung-cu-tai-san-co-dinh-can-bo-vua-xem-luon-thong-tin-chi-tiet-lien-quan-tren-cung-giao-dien-de-tiet-kiem-thoi-gian-lam-viec/): **split-pane master-detail** — select a list item, see its full details in the same screen.

### QuickBooks Online — chart of accounts
Source: https://quickbooks.intuit.com/learn-support/en-us/help-article/chart-accounts/learn-chart-accounts-quickbooks-online/L2yc6KBob_US_en_US

- Grid columns: **Name, Account Type, Detail Type, QuickBooks Balance, Bank Balance**.
- Buttons above the chart: **batch edit multiple accounts, print, edit displayed columns** (column chooser).
- In-app search bar to navigate. Detail View (https://www.cleverence.com/articles/quickbooks-documentation/chart-of-accounts-detail-view-quickbooks-intuit-5268/) supports **customize columns (drag-drop), filtering, sorting**.

### Xero
Sources: https://www.xero.com/ca/glossary/chart-of-accounts/, https://central.xero.com/s/article/Chart-of-accounts-templates-in-Xero-HQ-UK

- Each account: **name, brief description, GL code**; account **type** (Current Asset, Fixed Asset, Expense, Revenue, Equity...) determines report placement.
- No parent/child hierarchy in Xero — grouping happens in the report layout designer. Implication: hierarchy display (tree/indent) is a differentiator, not a given.

### NetSuite — structure best practices
Sources: https://www.netsuite.com/portal/resource/articles/accounting/chart-of-accounts.shtml, https://www.houseblend.io/articles/netsuite-chart-of-accounts-best-practices

- **Avoid deeply nested trees** — hard to manage, obscure totals. Practical max 2–4 levels.
- First digit of code = primary category (1 assets, 2 liabilities...), subsequent digits = subcategory + account.
- Dimensions over account proliferation.

### COA structure/UX guidance (consulting + community)
Sources: https://www.precisionfinancialllc.com/blog/how-to-structure-your-chart, https://www.accountingcoach.com/chart-of-accounts/explanation, https://plotpath.com/how-to-set-up-chart-of-accounts-in-quickbooks-and-xero, https://ramp.com/blog/chart-of-accounts-template

- **Flat hierarchy**: functional grouping, ≤2–4 levels; don't over-categorize by vendor (merge "AdWords"/"Bing" → "Online Advertising").
- Leave **gaps in numbering** (increments of 10/100) for later inserts.
- Normal balance explicit: **assets/expenses debit; liabilities/equity/revenue credit** (matches GoTax derived NormalBalance).
- Don't delete accounts with history — **inactivate** them.
- Debit/credit increase column + description column standard in COA tables (AccountingCoach).
- Batch import via CSV/Excel; batch edit (QBO).

## What accountants need from a COA screen (synthesis)

1. **Dense, keyboard-usable grid** — professionals scan hundreds of rows; single-select rows, dblclick/Enter = edit, Delete = delete-with-confirm.
2. **Hierarchy always visible** — loại grouping headers (1–9) + cấp 1 + indented children; collapse/expand; grouping must not hide totals.
3. **Columns that match how accountants think**: Số hiệu TK (monospace), Tên, Loại TK, Tính chất (Nợ/Có), Theo dõi chi tiết, Ngoại tệ, Trạng thái.
4. **Immediate search** across code (prefix) + name + English name, keeping ancestors visible.
5. **Safety rails**: confirm before delete, block deleting accounts with posted usage, warn on stopping tracking, cascade rules explicit.
6. **Batch-ish power**: duplicate (Nhân bản), export Excel/CSV, print.
7. **Status always visible** (Đang sử dụng / Ngừng theo dõi / Không dùng) + counts in the footer/status bar.
8. **Split-pane master-detail** (MISA R7) reduces jumping between list and form.

## Design implications for GoTax

- Constraints: Alpine.js static page, **no CSS framework** (vanilla CSS only until a framework lands), Vietnamese locale, Circular 99/2025 COA (9 loại 1-char grouping headers, cấp 1 = 3-digit, 221 accounts), existing two-pane tree+grid layout, app shell injected by `mountAppShell` (sidebar/topbar currently unstyled after Tailwind removal).
- Adopt: MISA toolbar semantics (Thêm/Sửa/Xóa/Nhân bản/Ngừng theo dõi + Xuất Excel/In), dblclick+Enter edit, Delete key with confirm, status filter dropdown, `*` required markers, forced-foreign 1112/1122/1132 notice, parent preselection when adding a child, status footer with counts, vanilla `web/static/css/coa.css` (page-scoped, survives framework swap), shell styles included (sidebar/topbar/hamburger) since Tailwind classes are inert.
- Defer (needs backend): Chuyển tài khoản hạch toán (transfer with số dư/phát sinh migration), batch edit, column chooser, cảnh báo vs bắt buộc modes for detail-by.
