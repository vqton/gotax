# Thủ kho (Warehouse Keeper) Module — Research Findings

## 1. Concept & Role Definition

**Thủ kho** (Warehouse Keeper / Storekeeper) is the person responsible for managing physical goods in the warehouse — receiving, storing, preserving, and issuing inventory. This is a **distinct role** from **Kế toán kho** (Warehouse Accountant), who handles the accounting/bookkeeping side.

| Aspect | Thủ kho (Warehouse Keeper) | Kế toán kho (Warehouse Accountant) |
|--------|---------------------------|-------------------------------------|
| Focus | Physical goods management | Accounting records & bookkeeping |
| Documents | Signs receipts, issues goods, maintains Sổ kho (stock ledger) | Creates journal entries, reconciles accounts |
| Access | See quantities, may hide cost prices | Full financial data |
| Separation | MISA explicitly separates these roles | If one person does both, no independent cross-check |

**Key principle (MISA):** When both roles are combined in one person, the system cannot provide independent cross-referencing data to detect discrepancies. The software should support role separation.

---

## 2. MISA SME 2026 — Thủ kho Module

**Source:** https://helpsme.misa.vn/2026/ac/thu-kho/ (updated 06/11/2025)

MISA SME 2026 treats **Thủ kho as module #17 of 18 modules** (Standard/Professional/Enterprise tiers). It is a separate, toggleable module from the main Kho (Warehouse) module.

### 2.1 Module Purpose
- Manage warehouse receipt/issue slips
- Record slips into the **Sổ kho** (Stock Ledger)
- Check **Biên bản kiểm kê** (Inventory Count Minutes)
- View warehouse reports

### 2.2 Workflow

**Daily:**
1. Login as Thủ kho
2. Go to tab **"Đề nghị nhập, xuất kho"** (Request for receipt/issue)
3. Select PNK/PXK slips (Ctrl+A for all, Ctrl+Click for selective)
4. Click **"Ghi sổ"** (Record to Stock Ledger)
5. Choose posting date (document date or custom date)
6. Confirm — system auto-records all selected slips into Thủ kho's Stock Ledger

**Monthly/Quarterly:**
- Generate **Báo cáo quản trị** (Management reports)

**Annually:**
- **Kiểm kê kho** (Physical inventory count)
- Review **Biên bản kiểm kê** (Inventory count minutes)

### 2.3 Key Features
- **Separate module** — can be hidden if company doesn't separate roles
- **Inbound/outbound slip review** — Thủ kho receives slips from Purchasing/Sales departments
- **Stock Ledger (Sổ kho)** — Thủ kho maintains their own ledger separate from accounting
- **Print PNK/PXK** — for physical reconciliation with actual goods
- **Internal transfer slips** — for cross-warehouse transfers
- **Cost price visibility control** — can hide cost prices from Thủ kho (only show quantities)
- **Inventory count reconciliation** — cross-reference between Thủ kho and Kế toán kho reports

### 2.4 Reports Available
- Báo cáo tổng hợp tồn kho (Inventory summary)
- Báo cáo đối chiếu kho bán hàng (Sales warehouse reconciliation)
- Báo cáo đối chiếu kho mua hàng (Purchase warehouse reconciliation)
- Thủ kho-specific reports in **Báo cáo → Báo cáo → Nhóm báo cáo Thủ kho**

### 2.5 Toggle Control
- Navigate to **Hệ thống → Tuỳ chọn → Ẩn/hiện nghiệp vụ**
- Uncheck **"Thủ kho có tham gia hệ thống"** to hide the module
- When hidden, Thủ kho reports still accessible via Reports tab

### 2.6 Role/Permission
- Separate login role: "Thủ kho"
- Can view quantities but optionally hide cost prices
- Cannot create purchase/sales documents — only review and record

### 2.7 Pricing
- Included in Standard (13 modules), Professional (15 modules), Enterprise (18 modules)
- MISA SME.NET 2026 Standard: 13,950,000 VND (one-time)
- Professional: 16,950,000 VND
- Enterprise: 19,950,000 VND

---

## 3. MISA AMIS Kế toán — Thủ kho Module

**Source:** https://helpact.misa.vn/kb/html_26000000/ (updated 27/05/2026)

### 3.1 Workflow (One-way process)
1. **Kế toán** creates receipt/issue slips in the system
2. Slips are transferred to **Thủ kho** for review and recording to stock ledger
3. If Kế toán modifies a slip, Thủ kho must un-record and re-record

### 3.2 Cross-referencing Report
- **"Đối chiếu nhập xuất giữa kế toán và thủ kho"** — reconciliation report between accounting and warehouse keeper
- If one person does both roles → no independent cross-check data exists

### 3.3 Access Control
- Login → click dropdown next company name → select "Thủ kho" mode
- Separate view with limited data access

---

## 4. FAST Business Online — Inventory Module

**Source:** https://fbohelp.fast.com.vn/fbohelp/ton-kho/ and https://fast.com.vn

FAST does **NOT have a separate "Thủ kho" module** like MISA. Instead, it has a unified **Fast Inventory** module with:

### 4.1 Features
- Receipt (Nhập kho): finished goods, raw materials from workshops, other
- Issue (Xuất kho): materials for production, internal consumption, auto from finished goods, value adjustment, other
- Transfer (Điều chuyển kho)
- **Lot tracking** — configurable per item
- **Inventory valuation methods:** Monthly average, Moving average (daily), FIFO, Specific identification
- **Barcode integration**
- Real-time stock check

### 4.2 Role Separation
- Achieved through **access control** (user permissions per screen/action)
- Not a separate module — same system, different permission levels
- Can restrict: view, add, edit, delete per user/group per screen or field

### 4.3 Key Difference from MISA
- FAST = unified inventory module with role-based access
- MISA = separate Thủ kho module with explicit workflow separation
- FAST approach is more standard for international ERP; MISA approach reflects Vietnamese regulatory tradition

---

## 5. Bravo ERP — Warehouse Module

**Source:** https://www.bravo.com.vn and https://bravo.com.vn/en/product/basic-modules/inventory-management

### 5.1 Overview
Bravo treats warehouse as a **basic module** within their 12-module ERP system. No separate "Thủ kho" module — uses role-based permissions.

### 5.2 Features
- Control value and quantity of materials/goods in stock
- Track rotation, usage, prevent damage in storage
- Create and control receipt/issue slips for production and circulation
- **Barcode tracking** for incoming/outgoing goods
- Multi-warehouse support
- Item classification by groups
- Serial number and lot tracking
- Date-based tracking

### 5.3 Integration
- Tightly integrated with: Accounting, Sales, Purchasing, Manufacturing
- Real-time data sync across modules
- No separate Thủ kho workflow — accounting and warehouse are in one system

---

## 6. Tryton ERP — Stock Module

**Source:** https://docs.tryton.org/latest/ and https://tryton-doc-text.readthedocs.io

### 6.1 Architecture
Tryton's Stock module defines fundamentals for all stock management:

**Location Types:**
- Storage (real places)
- Warehouse (meta-locations with input, storage, picking, output, lost-and-found)
- Customer (virtual)
- Supplier (virtual)
- Lost And Found (inventory corrections)
- Drop (drop shipping intermediary)
- Production
- View (logical grouping)

**Core Entities:**
- **Move** — product movement between locations
- **Supplier Shipment** — incoming from supplier (Draft → Received → Done)
- **Customer Shipment** — outgoing to customer (Draft → Waiting → Assigned → Packed → Done)
- **Internal Shipment** — between company locations
- **Inventory** — control and update stock levels

### 6.2 No Separate Thủ kho Role
Tryton follows international ERP convention: warehouse keeper is just a user role with specific permissions, not a separate module.

---

## 7. Vietnamese Regulatory Requirements

### 7.1 Thông tư 99/2025/TT-BTC (Circular 99)

**Source:** ketoanleanh.edu.vn, thuvienphapluat.vn, easybooks.vn

**Effective:** 01/01/2026 (replaces Circular 200/2014)

**Inventory-related provisions:**
- Inventory = assets purchased for sale, production, or service delivery
- Includes: goods in transit, raw materials, tools, WIP, finished goods, goods on consignment, bonded warehouse materials
- Must track **both quantity and value** for each warehouse, each department
- **Internal control requirement:** Enterprises must establish internal governance and control regulations defining rights, obligations, and responsibilities of departments/individuals
- **Voucher design:** Enterprises have FULL autonomy to design warehouse receipt/issue vouchers (no mandatory government forms) — but must comply with Article 16 of Law on Accounting 2015
- **Accounting system:** Only level-1 accounts are mandatory; enterprises freely design sub-accounts for warehouse management

### 7.2 Key Warehouse Documents (Circular 99)

| Document | Purpose | Creator | Keeper |
|----------|---------|---------|--------|
| **Phiếu nhập kho** (PNK) | Record goods received into warehouse | Purchasing/Production/Accounting | Thủ kho keeps Liên 2 for stock card |
| **Phiếu xuất kho** (PXK) | Record goods issued from warehouse | Sales/Production/Accounting | Thủ kho keeps Liên 2 for stock card |
| **Phiếu xuất điều chuyển** | Internal transfer between warehouses | Any authorized person | Both warehouse keepers |
| **Biên bản kiểm kê** | Physical inventory count results | Counting committee | Thủ kho + Accountant |
| **Sổ kho** (Stock Ledger) | Daily record of all movements | Thủ kho | Thủ kho maintains |
| **Thẻ kho** (Stock Card) | Per-item quantity tracking | Thủ kho | Thủ kho maintains |

### 7.3 Phiếu Nhập Kho (Receipt Voucher) — Structure

**2 copies** (purchased goods): Liên 1 (originator), Liên 2 (Thủ kho → Accountant)
**3 copies** (self-produced): + Liên 3 (delivery person)

**Required fields:**
- Date, slip number
- Item code, name, description
- Quantity, unit of measure, unit price, total value
- Supplier, source
- Warehouse location
- Signatures: Creator, Thủ kho, Warehouse Manager, Accountant Chief, Director

### 7.4 Thủ kho Responsibilities (Legal/Regulatory)

Based on cross-referenced Vietnamese sources:

1. **Physical custody** — responsible for all goods in their warehouse
2. **Slip verification** — check that PNK/PXK match actual goods
3. **Stock ledger maintenance** — record daily movements in Sổ kho
4. **Stock card (Thẻ kho)** — per-item quantity tracking
5. **Inventory counts** — participate in and sign Biên bản kiểm kê
6. **Discrepancy reporting** — immediately report any shortages, surpluses, damage
7. **Goods preservation** — proper storage, protection from damage/theft
8. **Signing authority** — co-sign receipts with delivery person; co-sign issues with recipient
9. **Report generation** — supply data for warehouse reports
10. **No financial authority** — Thủ kho does NOT handle money or create accounting entries

### 7.5 Separation of Duties

Vietnamese accounting tradition (per Law on Accounting 2015 and Circular 99) **strongly recommends** separation between:
- **Thủ kho** (physical goods custody)
- **Kế toán kho** (accounting records)
- **Thủ quỹ** (cash custody)

This is a **key internal control principle**. When combined:
- No independent cross-check possible
- Higher fraud/misstatement risk
- Auditor concern

---

## 8. Key Differences: Thủ kho Module vs General Warehouse Module

| Dimension | Thủ kho Module | General Warehouse Module |
|-----------|---------------|------------------------|
| **Primary user** | Warehouse keeper (physical custody role) | All warehouse staff, accountants, managers |
| **Core function** | Review & record receipt/issue slips into stock ledger | Full inventory lifecycle management |
| **Documents** | PNK, PXK, Sổ kho, Biên bản kiểm kê | All warehouse documents + valuation + costing |
| **Valuation** | Typically quantities only (cost hidden) | Full cost tracking, FIFO/AVG/etc. |
| **Integration** | Standalone review step in workflow | Integrated with purchasing, sales, manufacturing |
| **Access level** | Read-only review + record authority | Full CRUD on inventory |
| **Reporting** | Stock ledger, count minutes, reconciliation | Full suite: stock aging, turnover, ABC analysis |
| **Purpose** | Segregation of duties (internal control) | Operational efficiency |

---

## 9. Implementation Recommendations for GoTax

### 9.1 Feature Set (Based on MISA Pattern)

**Core screens:**
1. **Đề nghị nhập/xuất kho** — pending slips awaiting Thủ kho review
2. **Sổ kho** — Thủ kho's stock ledger (read-only after recording)
3. **Biên bản kiểm kê** — inventory count minutes
4. **Báo cáo Thủ kho** — warehouse keeper reports

**Workflow:**
1. Accounting/Purchasing/Sales creates PNK/PXK
2. Slip appears in Thủ kho's queue
3. Thủ kho reviews, verifies against physical goods
4. Thủ kho clicks "Ghi sổ" → recorded in Sổ kho
5. Cross-reference report shows discrepancies between Thủ kho and Accounting

**Access control:**
- Separate "Thủ kho" role/permission set
- Can view quantities; optionally hide cost prices
- Cannot create/edit/delete source documents
- Can only "record" (ghi sổ) or "un-record" (bỏ ghi)

### 9.2 Data Model Additions

**Entities needed:**
- `StockLedger` — Thủ kho's recording log (date, slip reference, recorded_by, status)
- `InventoryCount` — Biên bản kiểm kê with item-level counts
- `StockReconciliation` — Cross-reference between Thủ kho and Accounting records

**Fields on existing entities:**
- PNK/PXK: add `recorded_by_thu_kho` (boolean), `thu_kho_recorded_at` (timestamp)
- Add `cost_visibility` setting per user role (show/hide unit cost)

### 9.3 Reports

1. **Sổ kho tổng hợp** — consolidated stock ledger
2. **Báo cáo tồn kho Thủ kho** — inventory summary per Thủ kho view
3. **Đối chiếu nhập xuất kế toán - thủ kho** — reconciliation between accounting and keeper
4. **Biên bản kiểm kê** — count minutes with variance handling

### 9.4 Regulatory Compliance

- **Circular 99 compliance:** Support custom voucher design, sub-account flexibility
- **Law on Accounting 2015:** Maintain separation of duties, proper documentation
- **Internal control:** Cross-referencing capability is mandatory

---

## 10. Sources

| Source | URL | Key Finding |
|--------|-----|-------------|
| MISA SME 2026 Help | helpsme.misa.vn/2026/ac/thu-kho/ | Separate module #17, workflow, toggle control |
| MISA AMIS Kế toán | helpact.misa.vn/kb/html_26000000/ | One-way workflow, cross-referencing |
| MISA SME Overview | sme.misa.vn/48682/ | Thủ kho receives slips online from Purchasing/Sales |
| MISA Academy | academy.misa.vn | Thủ kho section in Kho module training |
| MISA Pricing | misavietnam.com | Thủ kho included in all tiers |
| FAST Business Online | fbohelp.fast.com.vn/fbohelp/ton-kho/ | No separate Thủ kho module, unified inventory |
| FAST General Features | fast.com.vn/fast-business-general-features | Role-based access, not separate module |
| BRAVO ERP | bravo.com.vn/en/product/basic-modules/inventory-management | Barcode tracking, no separate Thủ kho |
| Tryton Stock Module | docs.tryton.org/latest/ | Location-based architecture, no Thủ kho concept |
| Circular 99/2025 | thuvienphapluat.vn, ketoanleanh.edu.vn | Internal control, voucher autonomy, inventory accounting |
| Law on Accounting 2015 | thuvienphapluat.vn | Separation of duties principle |
| TopCV — Thủ kho Role | topcv.vn/thu-kho-la-gi | Job description, career path, responsibilities |
| Arito — Thủ kho | arito.vn/thu-kho-la-gi/ | Warehouse keeper vs accountant distinction |
| Bizzi — Kế toán kho | bizzi.vn/ke-toan-kho | Accountant vs keeper separation |
| InterLOG — PNK | interlogistics.com.vn | Receipt voucher structure, Thủ kho signing authority |
| EasyBooks — Circular 99 | easybooks.vn | Inventory accounting under Circular 99 |
