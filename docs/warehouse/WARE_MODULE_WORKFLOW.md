# Ware Module — End-to-End Enterprise Workflow

**Version:** v2.0  
**Benchmark:** MISA AMIS Kho hàng (06/2026), FAST Business Online Fast Inventory (03/2026), Bravo 10 ERP, Odoo Inventory 19.0, SAP MM-IM  
**Standard:** **Thông tư 99/2025/TT-BTC** (effective 01/01/2026 — replaces 200/2014), IAS 2, VAS 02, Decree 123/2020/ND-CP, Decree 254/2026/ND-CP  
**Key Regulatory Change:** **TK 611 eliminated** — direct posting to inventory accounts (151/152/153/156) per TT 99/2025  
**Status:** DRAFT (0% implemented)

---

## Table of Contents

1. [Master Data Setup](#1-master-data-setup)
2. [Opening Balance (Tồn đầu kỳ)](#2-opening-balance)
3. [Goods Receipt — Purchase (Nhập kho mua hàng)](#3-goods-receipt--purchase)
4. [Goods Receipt — Production (Nhập kho thành phẩm)](#4-goods-receipt--production)
5. [Goods Receipt — Sales Return (Nhập kho hàng bán trả lại)](#5-goods-receipt--sales-return)
6. [Goods Receipt — Direct / Other (Nhập kho khác)](#6-goods-receipt--direct--other)
7. [Goods Issue — Sales Delivery (Xuất kho bán hàng)](#7-goods-issue--sales-delivery)
8. [Goods Issue — Production Consumption (Xuất kho sản xuất)](#8-goods-issue--production-consumption)
9. [Goods Issue — Return to Supplier (Xuất kho trả lại NCC)](#9-goods-issue--return-to-supplier)
10. [Goods Issue — Write-Off / Disposal (Xuất kho hủy bỏ)](#10-goods-issue--write-off--disposal)
11. [Goods Issue — Sample / Internal (Xuất kho mẫu / nội bộ)](#11-goods-issue--sample--internal)
12. [Stock Transfer (Điều chuyển kho)](#12-stock-transfer)
13. [Stock Adjustment (Điều chỉnh tồn kho)](#13-stock-adjustment)
14. [Cost Revaluation (Đánh giá lại giá vốn)](#14-cost-revaluation)
15. [Stock Take / Physical Count (Kiểm kê kho)](#15-stock-take--physical-count)
16. [Inventory Valuation (Tính giá hàng tồn kho)](#16-inventory-valuation)
17. [Period-End Provision — TK 2294](#17-period-end-provision--tk-2294)
18. [GL Posting Reference (Circular 99/2025)](#18-gl-posting-reference-circular-992025)
19. [COA Accounts Mapping (TT 99/2025)](#19-coa-accounts-mapping-tt-992025)
20. [State Machine Reference](#20-state-machine-reference)

---

## 1. Master Data Setup

### 1.1 Create Warehouse (Khai báo kho)

```
CreateWarehouse → SuspendWarehouse
```

#### Step: Create Warehouse

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: WAREHOUSE_ACTIVE. Trigger: New physical/virtual storage location needed. |
| **Actor / Persona** | Warehouse Manager, System Admin |
| **Input Documents / Data** | Warehouse code (auto WH-XXXXX or manual), name, address, type (PHYSICAL/VIRTUAL/TRANSIT/BONDED), manager, allow_negative_stock flag, default cost_method. |
| **Business Validations & Rules** | Code unique per company. Alphanumeric, uppercase. Type must exist. At least one warehouse required per company (auto-create "Kho chính" on company setup). Cost method inherits: Company → Warehouse → Item. Cannot set allow_negative_stock=true for BONDED type. |
| **Output / Artifact Generated** | `Warehouse` record in DB. First warehouse for company auto-set as default. |
| **Exception / Failure Handling** | Duplicate code → reject with error "Mã kho đã tồn tại". Invalid type → reject. |

#### Step: Suspend Warehouse

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: ACTIVE. Post: SUSPENDED. Trigger: Warehouse closure, relocation. |
| **Actor / Persona** | Warehouse Manager, Director |
| **Input Documents / Data** | Warehouse ID, reason, date effective. |
| **Business Validations & Rules** | Stock on hand must be ZERO for all items. Open transfers for this warehouse must be completed or cancelled. No open stock takes. |
| **Output / Artifact Generated** | Warehouse status → SUSPENDED. All items' committed_qty released. |
| **Exception / Failure Handling** | Non-zero stock → block, require prior stock transfer to another warehouse. Open transfers → list them, block until resolved. |

### 1.2 Create Item (Khai báo vật tư, hàng hóa)

```
CreateItem → UpdateItem → SuspendItem
```

#### Step: Create Item

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: ITEM_ACTIVE. Trigger: New product/material/tool purchased or manufactured. |
| **Actor / Persona** | Inventory Planner, Warehouse Manager, Chief Accountant (for GL account mapping) |
| **Input Documents / Data** | Code (SKU, auto or manual), barcode (EAN/GTIN optional), name, category_id, unit, cost_method (WEIGHTED_AVG/FIFO/SPECIFIC_ID/STANDARD), standard_cost (if STANDARD), account_id (152/153/155/156 default inventory TK), cogs_account_id (632 default), revenue_account_id (5111), vat_rate, valuation_method, min/max stock, reorder_point, is_stockable, is_serialized, is_lot_tracked, shelf_life_days. |
| **Business Validations & Rules** | Code unique per company. account_id must be inventory type (151-158). Cost method must match company/warehouse hierarchy or override. If is_serialized → serial required per unit on receipt. If is_lot_tracked → lot# + expiry date required. Non-stock (is_stockable=false) skips inventory posting. Standard_cost required if cost_method=STANDARD. |
| **Output / Artifact Generated** | `Item` record in DB. Default GL account mapping stored for auto-account-determination. |
| **Exception / Failure Handling** | Duplicate barcode → warning (not block). Invalid account_id → block. Code too long (>50 chars) → reject. |

#### Step: Suspend Item

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: ACTIVE. Post: INACTIVE. Trigger: Item discontinued, obsolete. |
| **Actor / Persona** | Warehouse Manager, Chief Accountant |
| **Input Documents / Data** | Item ID, reason. |
| **Business Validations & Rules** | Stock on hand must be ZERO across all warehouses. No open POs/SOs referencing this item. |
| **Output / Artifact Generated** | Item → INACTIVE. Cannot be used on new transactions. |
| **Exception / Failure Handling** | Non-zero stock → block. Open orders → warn with order list. |

### 1.3 Create Item Category (Nhóm hàng)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: CATEGORY_CREATED. Trigger: Classification needed for ABC analysis, reporting. |
| **Actor / Persona** | Warehouse Manager, Inventory Planner |
| **Input Documents / Data** | Code, name, parent_id (for tree), abc_class (A/B/C), is_active. |
| **Business Validations & Rules** | Code unique per company. Parent category must exist and be active. Circular parent reference forbidden. |
| **Output / Artifact Generated** | `ItemCategory` record. Tree structure for hierarchical reporting. |
| **Exception / Failure Handling** | Circular parent → reject. Parent inactive → warn but allow. |

### 1.4 Configure Cost Method (Phương pháp tính giá)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: COST_METHOD_CONFIGURED. Trigger: Company setup, fiscal year start, accounting policy change. |
| **Actor / Persona** | Chief Accountant, CFO |
| **Input Documents / Data** | Method selection per company per item: WEIGHTED_AVG / FIFO / SPECIFIC_ID / STANDARD. Effective date. |
| **Business Validations & Rules** | **Circular 99/2025**: Methods allowed — Weighted Avg (periodic or perpetual), FIFO, Specific ID, retail method. LIFO prohibited. Method must be applied consistently across periods. Change requires disclosure in financial statements. Method determined per item, not per company (items in same company can use different methods per TT 99). |
| **Output / Artifact Generated** | `CostMethod` config stored on Item and Warehouse records. Valuation engine reads method at calculation time. |
| **Exception / Failure Handling** | Method change mid-period → flag for chief accountant. Inconsistent application across periods → audit warning. |

---

## 2. Opening Balance (Tồn đầu kỳ)

### Workflow

```
IMPORT → DRAFT → APPROVE → POST
```

### Step 2.1: Import/Enter Opening Balance

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: INIT. Post: DRAFT. Trigger: Fiscal year start, initial system setup, data migration. |
| **Actor / Persona** | Warehouse Keeper (enter), Chief Accountant (validate) |
| **Input Documents / Data** | Item code, warehouse, quantity, unit_price, lot# (if tracked), serial# (if tracked), expiry date (if tracked), total_value = qty × unit_price. Can be imported via CSV/Excel or manual entry per item per warehouse. |
| **Business Validations & Rules** | Items must be `is_stockable=true` and ACTIVE. Warehouse must be ACTIVE. Quantity > 0. Unit price ≥ 0. Duplicate (item+warehouse) → sum or reject (configurable). Total value per company within configurable threshold. |
| **Output / Artifact Generated** | `OpeningBalance` record (DRAFT) with `EntityType=INVENTORY_ITEM`. `OpeningBalanceDetail` per item+warehouse with Quantity, UnitPrice. |
| **Exception / Failure Handling** | Item not found → skip with error log. Item inactive → warn. Quantity=0 → skip row. Price negative → flag for review. |

### Step 2.2: Validate Totals

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: DRAFT (validated). Trigger: All rows imported, user clicks Validate. |
| **Actor / Persona** | System (auto), Chief Accountant (review) |
| **Input Documents / Data** | All `OpeningBalanceDetail` rows for this batch. |
| **Business Validations & Rules** | **Multi-company balancing per TT 99/2025**: Opening entries must be balanced (total debit = total credit) if posting to GL. Counterpart account required (3388 temporary or 419 retained earnings). Period must be OPEN (not CLOSED/LOCKED). |
| **Output / Artifact Generated** | Validation summary: total items, total quantity, total value, pass/fail status. |
| **Exception / Failure Handling** | Unbalanced (debit ≠ credit) → block with variance amount. Period CLOSED/LOCKED → require period reopen or post to next period. |

### Step 2.3: Approve

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT (validated). Post: APPROVED. Trigger: Chief Accountant approves. |
| **Actor / Persona** | Chief Accountant (approver). **Segregation: Creator ≠ Approver enforced.** |
| **Input Documents / Data** | Opening balance report with totals, validation pass. |
| **Business Validations & Rules** | Approver must have different ID from creator. Total value must match validated totals. |
| **Output / Artifact Generated** | `OpeningBalance` status → APPROVED. `approved_by`, `approved_at` set. |
| **Exception / Failure Handling** | Creator = Approver → block. Value changes between validate and approve → re-validation required. |

### Step 2.4: Post to GL & Stock Balance

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: APPROVED. Post: POSTED. Trigger: Approval confirmed. |
| **Actor / Persona** | System (auto-post on approve) |
| **Input Documents / Data** | Approved `OpeningBalance` with all details. |
| **Business Validations & Rules** | Account 151/152/153/156 must exist and be ACTIVE in COA. Counterpart account (3388/419) must exist. |
| **Output / Artifact Generated** | **Per Circular 99/2025 (direct posting, no TK 611):** `Dr 152/153/156 Cr 3388` (temporary offset) or `Dr 152/153/156 Cr 419` (retained earnings). `StockBalance` records created per item+warehouse: `quantity_on_hand = qty`, `stock_value = qty × unit_price`, `average_cost = unit_price`. |
| **Exception / Failure Handling** | Account frozen → block. GL posting fails → rollback stock balance changes. |

---

## 3. Goods Receipt — Purchase (Nhập kho mua hàng)

### MISA Multi-Step Benchmark

MISA AMIS (06/2026): 1-step (receive→stock) or multi-step (receive→quality→stock). GoTax implements configurable multi-step receipts per MISA standard.

```
PO → [GRN: DRAFT] → RECEIVE (1-step) or RECEIVE→QC→STOCK (3-step) → POST
```

### Step 3.1: Create GRN from PO

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: PO SENT or PARTIAL. Post: GRN DRAFT. Trigger: Supplier delivers goods with delivery note. |
| **Actor / Persona** | Warehouse Keeper, AP Clerk |
| **Input Documents / Data** | PO reference, GRN number (auto GRN-YYYYMM-XXXX), receipt_date, warehouse_id, supplier delivery note #/date. Per line: po_line_id, item_code, qty_received, qty_rejected, lot# (if tracked), serial# range (if serialized), expiry_date (if tracked). |
| **Business Validations & Rules** | PO must be in APPROVED/SENT/PARTIAL state. Received qty ≤ PO remaining qty + tolerance (default +5%). Warehouse must be ACTIVE. Items must match PO lines (code must match). Unit must match PO. Lot# required if item is_lot_tracked. Serial# required if is_serialized. |
| **Output / Artifact Generated** | `GRN` record (DRAFT). `GRNItem` per line with received/rejected qty and unit_price from PO. PO line `received_qty` cumulatively tracked (not yet incremented until POST). |
| **Exception / Failure Handling** | Over-receipt > 5% → block, require manager override with reason. Item not match PO → reject line. PO cancelled/closed → block. Duplicate supplier delivery note → warn (configurable block). |

### Step 3.2 (Optional): Multi-Step Receipt — Receive Goods

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: GRN DRAFT. Post: RECEIVED (goods in input area). Trigger: Physical goods arrive at receiving dock. |
| **Actor / Persona** | Warehouse Keeper (receiving clerk) |
| **Input Documents / Data** | Physical count, condition check, barcode/QR scan, weighbridge ticket (if bulk). |
| **Business Validations & Rules** | Count matches GRN line qty. Condition acceptable. For import: customs declaration #, duty/tax amounts available. |
| **Output / Artifact Generated** | Receiving confirmation. Status update on GRN. Goods physically in input/quality hold area (not yet available for sale). |
| **Exception / Failure Handling** | Count mismatch → adjust GRN line. Damage → reject line or route to inspection. No customs doc (import) → hold for customs clearance. |

### Step 3.3 (Optional): Multi-Step Receipt — Quality Check

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: RECEIVED (input area). Post: QC_PASSED or QC_FAILED. Trigger: Quality inspection required (config per item/warehouse). |
| **Actor / Persona** | Quality Control Inspector |
| **Input Documents / Data** | Inspection checklist, sample test results, photos (if damage). |
| **Business Validations & Rules** | Pass/fail decision per line or per batch. QC result must be documented. Re-inspection allowed (max 2 times). **MISA AMIS benchmark**: QC step configurable per warehouse, define "Là bước kiểm tra chất lượng" flag. |
| **Output / Artifact Generated** | QC report. If PASS → goods move to stock. If FAIL → route to quarantine/rejected area. |
| **Exception / Failure Handling** | QC fail → full reject (return to supplier) or partial reject (accept good units, return defective). QC fail with no supplier → write-off or rework. |

### Step 3.4: Post GRN (Confirm Goods Receipt)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT (or QC_PASSED). Post: POSTED. Trigger: Goods accepted into stock. |
| **Actor / Persona** | Warehouse Keeper (post), System (auto-generate GL) |
| **Input Documents / Data** | Final GRN with confirmed quantities, warehouse. |
| **Business Validations & Rules** | All items must have a cost (unit_price from PO or manual entry for non-PO). PO line `received_qty` incremented. PO auto-transition to PARTIAL (if partial) or RECEIVED (if complete). GRN cannot be modified after POST. |
| **Output / Artifact Generated** | **StockBalance update**: `quantity_on_hand += qty_received`, `on_order_qty -= qty_received`, `stock_value += qty_received × unit_price`. **InventoryTransaction** log (type=GRN_PURCHASE). **GL Entry per Circular 99/2025 (direct, no TK 611):** `Dr 152/153/156 Cr 331 (temp)`. PO status recalculated. PO line `received_qty` updated. |
| **Exception / Failure Handling** | Unit_price=0 → warn "Giá nhập bằng 0, kiểm tra lại". PO auto-close when all lines fully received → status RECEIVED. Negative stock after receipt → shouldn't happen (receipt only increases stock). |

### GL Posting Detail (TT 99/2025 — No TK 611)

```
Goods for resale:         Dr 1561     Cr 331 (temp)
Raw materials:            Dr 152      Cr 331 (temp)
Tools/instruments:        Dr 153      Cr 331 (temp)
Import duty added:        Dr 152/156  Cr 3333
Import VAT (reclaimable): Dr 1331     Cr 33312
Transport cost added:     Dr 152/156  Cr 331/111/112
Trade discount received:  Dr 331      Cr 156/152/632 (allocated per remaining/sold ratio)
```

**Key change vs TT 200**: TT 200 used `Dr 152 Cr 611` at purchase, then `Dr 611 Cr 331` at receipt. TT 99/2025 eliminates TK 611 entirely — direct `Dr 152 Cr 331`.

### Step 3.5: Cancel GRN

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: POSTED. Post: CANCELLED. Trigger: Goods defective, returned to supplier before invoicing. |
| **Actor / Persona** | Warehouse Keeper, AP Manager (approve) |
| **Input Documents / Data** | GRN ID, cancel reason, return qty (may be partial). |
| **Business Validations & Rules** | PO must not be fully invoiced (invoice qty > 0 for returned lines). Return qty ≤ original received qty - previously returned. |
| **Output / Artifact Generated** | Reverse GRN: `StockBalance -= qty_received`, reverse GL entry (`Dr 331(temp) Cr 152/156`). PO line `received_qty` decremented. PO status may revert from RECEIVED→PARTIAL. |
| **Exception / Failure Handling** | Already invoiced → block, require supplier credit note instead. Quantity already consumed/manufactured → block cancel. |

---

## 4. Goods Receipt — Production (Nhập kho thành phẩm)

```
ProdOrderCompleted → CreateReceipt → POST
```

### Step 4.1: Create Production Receipt

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: PRODUCTION_ORDER.COMPLETED. Post: DRAFT. Trigger: Manufacturing order finished, goods move to warehouse. |
| **Actor / Persona** | Production Supervisor, Warehouse Keeper |
| **Input Documents / Data** | Production order ref, item codes, qty completed, unit cost (from production cost roll-up), warehouse, lot# if tracked. |
| **Business Validations & Rules** | Production order must be COMPLETED or IN_PROGRESS status. All raw materials must be issued (consumption recorded). Cost allocation complete from production costing engine. Quantity > 0. |
| **Output / Artifact Generated** | Production receipt doc (DRAFT). |
| **Exception / Failure Handling** | Production order not complete → block. Raw materials not fully issued → warn, partial receipt allowed. Cost not yet calculated → use estimated cost, flag for adjustment. |

### Step 4.2: Post Production Receipt

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: POSTED. Trigger: Finished goods accepted into stock. |
| **Actor / Persona** | Warehouse Keeper, System |
| **Input Documents / Data** | Final receipt doc with confirmed quantities. |
| **Business Validations & Rules** | Cost > 0. Account 155 (finished goods) or relevant inventory account must be active. |
| **Output / Artifact Generated** | **StockBalance**: `+qty`, `+value` at unit cost. **GL**: `Dr 155/156 Cr 154 (WIP)`. WIP account (154) decreased by production cost. |
| **Exception / Failure Handling** | Cost = 0 → flag for cost revaluation. WIP account insufficient balance → warn. |

---

## 5. Goods Receipt — Sales Return (Nhập kho hàng bán trả lại)

```
CreditNoteCreated → InspectReturnedGoods → POST
```

### Step 5.1: Create Return Receipt from Credit Note

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: CREDIT_NOTE POSTED. Post: DRAFT. Trigger: Customer returns goods, credit note issued. |
| **Actor / Persona** | Sales Clerk, Warehouse Keeper |
| **Input Documents / Data** | Credit Note reference, original DN line, item code, return qty, condition (sellable/damaged), warehouse. |
| **Business Validations & Rules** | Return qty ≤ original delivered qty - previously returned for this DN line. Item must match original delivery. Credit note must be POSTED status. |
| **Output / Artifact Generated** | Return receipt (DRAFT). |
| **Exception / Failure Handling** | Qty exceeds delivered → block. Credit note not posted → block. |

### Step 5.2: Inspect and Post Return Receipt

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: POSTED. Trigger: Returned goods inspected and accepted. |
| **Actor / Persona** | Warehouse Keeper (inspect), Quality Control (if damaged) |
| **Input Documents / Data** | Condition assessment (sellable / damaged / obsolete). |
| **Business Validations & Rules** | **Per Circular 99/2025**: Return valued at original cost price from DN line (not selling price). If damaged → route to damaged stock or write-off. |
| **Output / Artifact Generated** | **StockBalance**: `+qty`, `+value` at original cost. **GL per TT 99**: `Dr 152/156 Cr 632` (reverse COGS at original cost). InventoryTransaction (GRN_RETURN). |
| **Exception / Failure Handling** | Damaged goods → separate stock location. Cost price unavailable → use average cost at time of return. |

---

## 6. Goods Receipt — Direct / Other (Nhập kho khác)

### Bravo Benchmark

Bravo 10 supports: Phiếu nhập mua, nhập khẩu, nhập thành phẩm, **nhập khác** (donation, samples, owner contribution, found surplus).

```
CreateUnplannedReceipt → Approve(if needed) → POST
```

### Step 6.1: Create Direct Receipt

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: DRAFT. Trigger: Goods received without PO — donation, sample, owner contribution, branch transfer, surplus found. |
| **Actor / Persona** | Warehouse Keeper |
| **Input Documents / Data** | Item codes, quantities, unit_price, counterpart account (111/112/338/411/711), reason, reference document#. |
| **Business Validations & Rules** | Reason required from approved list: DONATION / SAMPLE / OWNER_CONTRIBUTION / BRANCH_TRANSFER / SURPLUS / OTHER. Counterpart account must be appropriate for reason. |
| **Output / Artifact Generated** | Direct receipt doc (DRAFT). |
| **Exception / Failure Handling** | No counterpart account → block. Invalid reason → block. |

### Step 6.2: Approve (if value > threshold)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: APPROVED. Trigger: Value exceeds auto-approve threshold. |
| **Actor / Persona** | Warehouse Manager (≤10M VND), Chief Accountant (>10M VND) |
| **Input Documents / Data** | Receipt details, evidence document. |
| **Business Validations & Rules** | Thresholds: <10M auto, 10-100M Warehouse Manager, >100M Chief Accountant. |
| **Output / Artifact Generated** | Approval record. Status→APPROVED. |
| **Exception / Failure Handling** | Threshold exceeded without approval → block. Approver cannot be creator. |

### Step 6.3: Post Direct Receipt

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT or APPROVED. Post: POSTED. |
| **Actor / Persona** | Warehouse Keeper, System |
| **Input Documents / Data** | Approved receipt. |
| **Business Validations & Rules** | Account determination per counterpart type. |
| **Output / Artifact Generated** | **Per Circular 99/2025 (direct to inventory):** `Dr 152/153/156 Cr [counterpart_account]`. Counterpart by reason: DONATION→Cr 711, OWNER_CONTRIBUTION→Cr 4111, BRANCH_TRANSFER→Cr 336, SURPLUS→Cr 3381 (awaiting resolution) or 711. StockBalance: +qty, +value. |
| **Exception / Failure Handling** | Counterpart account frozen → block. GL posting fails → rollback. |

---

## 7. Goods Issue — Sales Delivery (Xuất kho bán hàng)

### MISA Multi-Step Benchmark

MISA AMIS (06/2026): 1-step (pick→ship) or multi-step (pick→pack→ship with quality hold). Configurable per warehouse.

```
SO Confirmed → [DN: DRAFT] → Pick → Pack(opt) → POST → GL
```

### Step 7.1: Create Delivery Note from SO

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: SO CONFIRMED or PROCESSING. Post: DN DRAFT. Trigger: Goods ready to deliver, warehouse receives picking request. |
| **Actor / Persona** | Sales Clerk, Warehouse Keeper |
| **Input Documents / Data** | SO reference, DN number (auto DN-YYYYMM-XXXX), warehouse, shipping info (carrier, tracking#, delivery address). Per line: so_line_id, item_code, qty_delivered, warehouse. |
| **Business Validations & Rules** | SO must be CONFIRMED or PROCESSING. Delivered qty ≤ SO remaining qty + tolerance (default +5%). Item must match SO line. `delivered_qty` tracked cumulatively on SO line. Warehouse must be ACTIVE and allow outbound. |
| **Output / Artifact Generated** | `DN` record (DRAFT). `DNLine` per item. SO status auto-updated when fully delivered. |
| **Exception / Failure Handling** | SO not confirmed → block. Over-delivery > 5% → require manager approval. Invalid warehouse for sales region → warn. |

### Step 7.2 (Optional): Validate Stock Availability

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DN DRAFT. Post: STOCK_CHECKED. Trigger: Pre-pick stock check. |
| **Actor / Persona** | System (auto) |
| **Input Documents / Data** | DN lines, current StockBalance, committed_qty from other SOs. |
| **Business Validations & Rules** | `available = quantity_on_hand - committed_qty` ≥ delivered qty. If `allow_negative_stock=false` on warehouse and insufficient → block. **FAST/Bravo benchmark**: configurable per warehouse "Cho phép xuất khi chưa đủ hàng". If allowed and insufficient → proceed with warning, resulting negative balance. |
| **Output / Artifact Generated** | Availability validation: PASS (green), WARN (yellow, partial), BLOCK (red). |
| **Exception / Failure Handling** | Insufficient stock + allow_negative=false → block with shortage report. Insufficient + allow_negative=true → warn, create backorder for remaining qty. Partial fulfillment → split DN, create backorder SO for remainder. |

### Step 7.3: Pick Goods (MISA multi-step)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT/STOCK_CHECKED. Post: PICKED. Trigger: Warehouse picker collects goods from shelves. |
| **Actor / Persona** | Warehouse Keeper (picker) |
| **Input Documents / Data** | Picking list (item code, location, bin, qty to pick), barcode scanner, lot#/serial# confirmation. |
| **Business Validations & Rules** | Picked qty = DN qty (allow short pick with reason). Serial# scanned for serialized items (must match expected). FIFO/FEFO picking rule enforced (pick oldest lots first for FIFO items). **Odoo benchmark**: removal strategies — FIFO, LIFO, FEFO configurable per product. |
| **Output / Artifact Generated** | Picking confirmation. Items moved to staging/dispatch area. `committed_qty` consumed. |
| **Exception / Failure Handling** | Short pick → partial DN, backorder remainder. Wrong serial# scanned → alert, require correction. Item not found at location → trigger cycle count. |

### Step 7.4 (Optional): Pack Goods

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: PICKED. Post: PACKED. Trigger: Goods packed for shipment. |
| **Actor / Persona** | Warehouse Keeper (packer) |
| **Input Documents / Data** | Packing list, package weights/dimensions, carrier label. |
| **Business Validations & Rules** | All picked items accounted for. Dangerous goods -> special handling flag. |
| **Output / Artifact Generated** | Packing confirmation. Shipping label generated. |
| **Exception / Failure Handling** | Item missing → return to pick step. Damaged during packing → replace picked item or remove from delivery. |

### Step 7.5: Post DN (Confirm Delivery + COGS)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT (or PACKED). Post: POSTED. Trigger: Goods leave warehouse, delivered to customer/carrier. |
| **Actor / Persona** | Warehouse Keeper (post), System (GL auto-post) |
| **Input Documents / Data** | Final delivery confirmation, receiver signature (or carrier acknowledgment), delivery date. |
| **Business Validations & Rules** | Cost price must be determinable per item (FIFO/WAC/Specific ID at time of delivery). `qty_delivered` must be > 0. GL posting mandatory. No post-back after POST. |
| **Output / Artifact Generated** | **StockBalance**: `quantity_on_hand -= qty`, `stock_value -= qty × cost_price`. **GL per TT 99/2025**: `Dr 632 Cr 152/156` (COGS at cost price). SO line `delivered_qty` incremented. `InventoryTransaction` (GI_SALES) logged. SO status auto-advances: if all lines fulfilled → DELIVERED. |
| **Exception / Failure Handling** | Cost price unavailable → use weighted average, flag for review. Zero cost price → block, require cost revaluation. Negative stock after issue → warn (if allowed) or block (if not). |

### Step 7.6: Cancel DN

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: POSTED. Post: CANCELLED. Trigger: Delivery cancelled, goods returned (pre-invoice return). |
| **Actor / Persona** | Warehouse Keeper, Sales Manager |
| **Input Documents / Data** | DN ID, cancel reason. |
| **Business Validations & Rules** | Invoice must not exist for this DN. Return qty ≤ original delivered qty. |
| **Output / Artifact Generated** | Reverse GL: `Dr 152/156 Cr 632`. StockBalance restored. SO line `delivered_qty` decremented. SO status may revert to previous state. |
| **Exception / Failure Handling** | Already invoiced → block, require credit note instead. Goods already consumed/production → block. |

---

## 8. Goods Issue — Production Consumption (Xuất kho sản xuất)

```
ProdOrderReleased → MaterialRequest → POST
```

### Step 8.1: Create Material Issue from Production Order

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: PRODUCTION_ORDER.RELEASED or IN_PROGRESS. Post: DRAFT. Trigger: Production requires raw materials. |
| **Actor / Persona** | Production Supervisor, Warehouse Keeper |
| **Input Documents / Data** | Production order ref, BOM item list, planned qty, actual qty to issue, warehouse. |
| **Business Validations & Rules** | Production order must be RELEASED/IN_PROGRESS. Material list must match BOM (with substitution tolerance). Stock available ≥ issue qty. Account 152/153 (raw materials/tools) must be active. |
| **Output / Artifact Generated** | Material issue slip (DRAFT). |
| **Exception / Failure Handling** | BOM mismatch → manager approval for substitution. Stock insufficient → partial issue, production order flagged for material shortage. |

### Step 8.2: Post Material Issue

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: POSTED. Trigger: Materials physically issued from warehouse to production. |
| **Actor / Persona** | Warehouse Keeper, System |
| **Input Documents / Data** | Final issue quantities confirmed. |
| **Business Validations & Rules** | Cost determined per cost method. |
| **Output / Artifact Generated** | **StockBalance**: `-qty`, `-value` at cost. **GL per TT 99/2025**: `Dr 154 (WIP) Cr 152/153`. InventoryTransaction (GI_PRODUCTION). |
| **Exception / Failure Handling** | Cost = 0 → use standard cost or average. WIP account invalid → block. |

---

## 9. Goods Issue — Return to Supplier (Xuất kho trả lại NCC)

```
SupplierAgreesReturn → CreateReturnGRN → POST
```

### Step 9.1: Create Return GRN

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: Original GRN POSTED. Post: RETURN_GRN DRAFT. Trigger: Defective/incorrect goods sent back to supplier before or after invoice. |
| **Actor / Persona** | Warehouse Keeper, AP Clerk |
| **Input Documents / Data** | Original GRN reference, item codes, return qty, return reason (defect/wrong_item/excess/other), supplier acknowledgment. |
| **Business Validations & Rules** | Return qty ≤ original received qty - previously returned. Original GRN must be POSTED. If already invoiced → must link to supplier credit note. |
| **Output / Artifact Generated** | Return GRN (DRAFT, type=RETURN). |
| **Exception / Failure Handling** | Already invoiced + no credit note → block (need AP adjustment first). Return qty > available → block. |

### Step 9.2: Post Return to Supplier

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: POSTED. Trigger: Goods shipped back to supplier. |
| **Actor / Persona** | Warehouse Keeper, System |
| **Input Documents / Data** | Return confirmation, shipping proof. |
| **Business Validations & Rules** | GL posting depends on invoice status. |
| **Output / Artifact Generated** | **Uninvoiced (pre-invoice):** Reverse original GRN: `Dr 331(temp) Cr 152/156`. **Invoiced (post-invoice):** Must use credit note instead: `Dr 331(AP) Cr 152/156 + Cr 1331`. StockBalance: `-qty`, `-value`. |
| **Exception / Failure Handling** | Post-invoice return without credit note → block. GL posting fails → rollback. |

---

## 10. Goods Issue — Write-Off / Disposal (Xuất kho hủy bỏ)

```
Request → Approve → POST
```

### Step 10.1: Create Write-Off Request

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: DRAFT. Trigger: Goods damaged, expired, stolen, obsolete. |
| **Actor / Persona** | Warehouse Keeper |
| **Input Documents / Data** | Item codes, quantities, write-off reason (DAMAGE/EXPIRY/THEFT/OBSOLETE/OTHER), estimated cost, evidence (photos, inspection report), proposed GL account. |
| **Business Validations & Rules** | Stock available ≥ write-off qty. Reason required. Evidence required for THEFT/DAMAGE (police report for theft). |
| **Output / Artifact Generated** | Write-off request (DRAFT). |
| **Exception / Failure Handling** | Insufficient stock → block. No reason → block. |

### Step 10.2: Approve Write-Off

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: APPROVED. Trigger: Manager reviews and approves. |
| **Actor / Persona** | Warehouse Manager (≤10M VND), Chief Accountant (10M-100M), Director (>100M) |
| **Input Documents / Data** | Write-off request with evidence, GL account recommendation. |
| **Business Validations & Rules** | Approval thresholds enforced. Segregation: Approver ≠ Creator. **Bravo benchmark**: multi-level approval based on value. |
| **Output / Artifact Generated** | Approval chain completed. Request → APPROVED. |
| **Exception / Failure Handling** | Threshold exceeded without proper approval → block. Approver same as creator → block. |

### Step 10.3: Post Write-Off

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: APPROVED. Post: POSTED. |
| **Actor / Persona** | System, Chief Accountant (review) |
| **Input Documents / Data** | Approved write-off. |
| **Business Validations & Rules** | GL account per reason per **Circular 99/2025**: DAMAGE(632), EXPIRY(632), THEFT(811), OBSOLETE(632), OTHER(642). |
| **Output / Artifact Generated** | **GL per TT 99 (direct posting):** `Dr 632/642/811 Cr 152/153/156`. StockBalance: `-qty`, `-value` at cost. InventoryTransaction (GI_WRITEOFF). |
| **Exception / Failure Handling** | GL account frozen → block. Large write-off triggers insurance claim record (if THEFT). |

---

## 11. Goods Issue — Sample / Internal (Xuất kho mẫu / nội bộ)

```
Request → Approve → POST
```

### Step 11.1: Create Internal Issue Request

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: DRAFT. Trigger: Department needs goods for samples, marketing, charity, internal use. |
| **Actor / Persona** | Department Requester (any dept), Warehouse Keeper |
| **Input Documents / Data** | Item codes, qty, purpose (SAMPLE/MARKETING/CHARITY/INTERNAL_USE/ADMIN), department ID, cost_center, counterpart account. |
| **Business Validations & Rules** | Purpose required. Department required. Counterpart account must match purpose per policy. Stock available ≥ qty. |
| **Output / Artifact Generated** | Internal issue request (DRAFT). |
| **Exception / Failure Handling** | Purpose not in approved list → block. No counterpart account → block. Insufficient stock → warn. |

### Step 11.2: Post Internal Issue

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT (auto-approved for low value). Post: POSTED. |
| **Actor / Persona** | Warehouse Keeper, System |
| **Input Documents / Data** | Issue confirmation, department acknowledgment. |
| **Business Validations & Rules** | GL account per purpose. |
| **Output / Artifact Generated** | **GL per TT 99**: SAMPLES→`Dr 6417 Cr 152/156`, MARKETING→`Dr 6417 Cr 152/156`, CHARITY→`Dr 811 Cr 152/156`, INTERNAL_USE→`Dr 642/627 Cr 152/156`, ADMIN→`Dr 642 Cr 152/156`. StockBalance: `-qty`, `-value` at cost. |
| **Exception / Failure Handling** | Wrong GL account → block. Value > monthly budget (if configured) → warn. |

---

## 12. Stock Transfer (Điều chuyển kho)

### FAST/Bravo Benchmark

FAST: 2-step transfer (request → confirm). Bravo: Phiếu điều chuyển kho + Phiếu xuất điều chuyển vị trí.

```
TransferRequest → Approve → Pick → Dispatch(Intransit) → Receive → POST
```

### Step 12.1: Create Transfer Request

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: REQUESTED. Trigger: Stock surplus at WH A, shortage at WH B; warehouse consolidation. |
| **Actor / Persona** | Warehouse Manager, Inventory Planner |
| **Input Documents / Data** | Source warehouse, destination warehouse, item codes, qty, expected date, reason. |
| **Business Validations & Rules** | Source ≠ Destination. Both warehouses ACTIVE. Items exist in source warehouse. Stock available ≥ qty (if allow_negative=false). |
| **Output / Artifact Generated** | `StockTransfer` (REQUESTED). Source WH: `committed_qty += qty`. |
| **Exception / Failure Handling** | Source = Dest → reject. Invalid item → reject. Insufficient stock → partial transfer allowed. |

### Step 12.2: Approve Transfer

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: REQUESTED. Post: APPROVED. |
| **Actor / Persona** | Warehouse Manager (≤ threshold), Chief Accountant (> threshold) |
| **Input Documents / Data** | Transfer request. |
| **Business Validations & Rules** | Segregation: Approver ≠ Requester. Value threshold (configurable). |
| **Output / Artifact Generated** | Status → APPROVED. |
| **Exception / Failure Handling** | Value exceeds approver's authority → escalate. |

### Step 12.3: Pick Goods (Source Warehouse)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: APPROVED. Post: PICKED. |
| **Actor / Persona** | Warehouse Keeper (source WH) |
| **Input Documents / Data** | Picking slip, barcode scan, lot#/serial# confirmation. |
| **Business Validations & Rules** | Picked qty ≤ requested qty. Serial# scanned if serialized. FIFO picking for FIFO items. |
| **Output / Artifact Generated** | Status → PICKED. `committed_qty` consumed. Physical goods ready at dispatch area. |
| **Exception / Failure Handling** | Short pick → partial transfer, remaining cancelled or backordered. |

### Step 12.4: Dispatch / Ship (In Transit)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: PICKED. Post: IN_TRANSIT. |
| **Actor / Persona** | Warehouse Keeper (source), Shipper/Logistics |
| **Input Documents / Data** | Dispatch confirmation, carrier name, tracking#, shipping date. |
| **Business Validations & Rules** | All picked items accounted for before dispatch. |
| **Output / Artifact Generated** | Status → IN_TRANSIT. **Stock (TT 99):** If transit tracking enabled → `Dr 151 (goods in transit) Cr 152/156`. If simple transfer → `Dr 152/156 (dest) Cr 152/156 (src)`. Source WH stock decreased. |
| **Exception / Failure Handling** | Damage during loading → remove from transfer, create write-off. |

### Step 12.5: Receive at Destination

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: IN_TRANSIT. Post: RECEIVED or COMPLETED. |
| **Actor / Persona** | Warehouse Keeper (destination WH) |
| **Input Documents / Data** | Receiving confirmation, qty received, qty damaged_in_transit, lot#/serial# verification. |
| **Business Validations & Rules** | Received qty ≤ shipped qty. Damaged qty ≤ received qty. Lot#/serial# match shipped records. |
| **Output / Artifact Generated** | Status → RECEIVED (if partial) or COMPLETED (if fully received). Destination WH stock increased. Transit account reversed: `Dr 152/156 Cr 151`. InventoryTransaction (TRANSFER_RECEIVE). |
| **Exception / Failure Handling** | Short receipt → adjust transfer, document missing qty. Damaged in transit → file carrier claim: `Dr 1388 (claim receivable) Cr 151 (transit)`. |

---

## 13. Stock Adjustment (Điều chỉnh tồn kho)

```
Request → Approve → POST
```

### Step 13.1: Create Adjustment Request

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: DRAFT. Trigger: Discrepancy found (count variance, data error, rounding). |
| **Actor / Persona** | Warehouse Keeper, Inventory Controller |
| **Input Documents / Data** | Item code, warehouse, old_qty, new_qty, old_value, new_value, reason, reference (count sheet#, audit ref#). |
| **Business Validations & Rules** | Qty change ≠ 0. Reason required from approved list: COUNTING_ERROR / DATA_ENTRY / ROUNDING / DAMAGE_FOUND / THEFT_FOUND / OTHER. Value change must be consistent with qty change (unit cost preserved unless revaluation). |
| **Output / Artifact Generated** | Adjustment request (DRAFT). |
| **Exception / Failure Handling** | No qty change → reject ("No adjustment needed"). No reason → block. |

### Step 13.2: Approve Adjustment

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: APPROVED. |
| **Actor / Persona** | Warehouse Manager (≤5M), Chief Accountant (>5M) |
| **Input Documents / Data** | Adjustment details, investigation notes. |
| **Business Validations & Rules** | Approval thresholds enforced. Frequent adjustments on same item → flag for audit. |
| **Output / Artifact Generated** | Status → APPROVED. |
| **Exception / Failure Handling** | Same item adjusted 3+ times this month → escalate to Chief Accountant. |

### Step 13.3: Post Adjustment

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: APPROVED. Post: POSTED. |
| **Actor / Persona** | System, Chief Accountant (post) |
| **Input Documents / Data** | Approved adjustment. |
| **Business Validations & Rules** | GL account per direction + reason. |
| **Output / Artifact Generated** | **Increase (qty gain):** `Dr 152/156 Cr 3381` (unresolved) → later to 711 or 632 if resolved. **Decrease (qty loss):** `Dr 632/642/811 Cr 152/156`. StockBalance corrected. InventoryTransaction (ADJUST_INCREASE or ADJUST_DECREASE). |
| **Exception / Failure Handling** | Causing negative stock → flag. Large adjustment > threshold → require second approval. |

---

## 14. Cost Revaluation (Đánh giá lại giá vốn)

```
Request → POST
```

### Step 14.1: Create Revaluation Request

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: DRAFT. Trigger: Standard cost update, FX rate change for import goods, audit correction. |
| **Actor / Persona** | Chief Accountant (solely) |
| **Input Documents / Data** | Item code, warehouse (or all warehouses), old_unit_cost, new_unit_cost, reason, effective date. |
| **Business Validations & Rules** | Must affect all stock of this item in company (or per warehouse if configured). New cost ≥ 0. Reason required. Only Chief Accountant role can create revaluation. |
| **Output / Artifact Generated** | Revaluation record (DRAFT). |
| **Exception / Failure Handling** | Old cost = new cost → reject ("No change"). Negative new cost → reject. Non-chief-accountant attempting → 403 Forbidden. |

### Step 14.2: Post Revaluation

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: DRAFT. Post: POSTED. |
| **Actor / Persona** | System (auto-post by Chief Accountant) |
| **Input Documents / Data** | Revaluation details. |
| **Business Validations & Rules** | Only Standard Cost items should be revalued regularly. FIFO/WAC revaluation only for audit corrections (with audit ref). |
| **Output / Artifact Generated** | **Per TT 99/2025**: Financial adjustment through TK 611-type function → `Dr 152/156 Cr 611` (cost increase) or `Dr 611 Cr 152/156` (cost decrease). StockBalance: value recalculated. New average cost set. Note: TT 99 eliminated TK 611 for routine purchase, but cost variance account retained for standard cost variance. |
| **Exception / Failure Handling** | Revaluation > 10% of stock value → CFO notification. |

---

## 15. Stock Take / Physical Count (Kiểm kê kho)

### Full Workflow

```
PLAN → COUNT → RECOUNT(if needed) → APPROVE_VARIANCE → POST_ADJUSTMENT → COMPLETE
```

### Step 15.1: Plan Count

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: none. Post: PLANNED. Trigger: Period-end closing, schedule (monthly A/quarterly B/yearly C classes), ad-hoc request. |
| **Actor / Persona** | Warehouse Manager, Chief Accountant (approve plan) |
| **Input Documents / Data** | Warehouse list, item categories (ABC class), count date/time, assigned counters (1st + optional 2nd counter), count type (full/cycle/ad-hoc). |
| **Business Validations & Rules** | No other count in progress for same warehouse. ABC class frequency: A=monthly, B=quarterly, C=yearly. Book qty snapshot frozen at plan creation time. Period must be OPEN. |
| **Output / Artifact Generated** | `StockTake` record (PLANNED). Book qty frozen per item+warehouse. Count sheets generated (paper/PDF/scanner format). |
| **Exception / Failure Handling** | Count already in progress for warehouse → block. Period CLOSED → block. |

### Step 15.2: Execute Count

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: PLANNED. Post: IN_PROGRESS. Trigger: Count sheets distributed, physical counting begins. |
| **Actor / Persona** | Warehouse Keeper (Counter 1), optional Supervisor (Counter 2 for A-class items) |
| **Input Documents / Data** | Physical quantities per item+location. Barcode/RFID scans. Lot#/serial# verification. Date/time of count. |
| **Business Validations & Rules** | **MISA AMIS benchmark**: blind count (book qty hidden) or open count (book shown) — configurable. Items on plan = items counted. Missing items on plan → investigate. Count per location must be complete before moving to next. |
| **Output / Artifact Generated** | Physical qty recorded per item. Count record with timestamp, counter ID. Variance calculated: `variance_qty = book_qty - physical_qty`, `variance_value = variance_qty × unit_cost`. |
| **Exception / Failure Handling** | Item found not on plan → add to count sheet (investigate origin). Item blocked/quarantined → count separately with note. Scanner battery failure → paper backup. |

### Step 15.3: Re-Count (if variance > tolerance)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: IN_PROGRESS (variance found). Post: IN_PROGRESS (recounted). Trigger: Variance exceeds configurable tolerance. |
| **Actor / Persona** | Supervisor, Second Counter (different from first) |
| **Input Documents / Data** | Items with variance > tolerance, fresh count sheet. |
| **Business Validations & Rules** | **Tolerance**: ±1 unit or ±0.5% of value (whichever larger) per item per warehouse. Items within tolerance auto-accepted. If 2nd count matches 1st → confirm variance. If 2nd count differs from 1st → 3rd count by manager. Max 3 counts. |
| **Output / Artifact Generated** | 2nd count record. If match 1st → variance confirmed. If differ → 3rd count triggered. |
| **Exception / Failure Handling** | 3 counts all different → force close by manager with override. Variance within tolerance → skip to auto-accept. |

### Step 15.4: Approve Variance

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: IN_PROGRESS (counted). Post: APPROVED. Trigger: All variances reviewed, reason assigned. |
| **Actor / Persona** | Warehouse Manager (≤ threshold), Chief Accountant (> threshold) |
| **Input Documents / Data** | Count variance report, investigation notes, reason classification: COUNTING_ERROR / DAMAGE / THEFT / SYSTEM_ERROR / PRIOR_ADJUSTMENT / OTHER. |
| **Business Validations & Rules** | Reason required for each variance line. Segregation: Approver ≠ Counter. Threshold for escalation: 10M VND total variance → Chief Accountant. 100M VND → CFO. |
| **Output / Artifact Generated** | Variance approved. `StockTake. approved_by` set. Variance lines marked for posting. |
| **Exception / Failure Handling** | Total variance exceeds threshold without proper approval level → escalate. Reason = THEFT → insurance claim process triggered. |

### Step 15.5: Post Adjustment + Complete

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: APPROVED. Post: COMPLETED. |
| **Actor / Persona** | System (auto-generate adjustment docs), Chief Accountant (review) |
| **Input Documents / Data** | Approved variance lines. |
| **Business Validations & Rules** | GL account per variance reason per **TT 99/2025**: COUNTING_ERROR→632, DAMAGE→632, THEFT→811, SYSTEM_ERROR→632/711, PRIOR_ADJ→adjust prior period. |
| **Output / Artifact Generated** | **Adjustment notes auto-generated** (MISA/FAST/Bravo benchmark). **GL entries**: GAIN→`Dr 152/156 Cr 632`(or 711 if income). LOSS→`Dr 632/811 Cr 152/156`. StockBalance corrected. `StockTake` → COMPLETED. Closed period lock enforced. |
| **Exception / Failure Handling** | Period already closed → post as prior-period adjustment. Any unresolved variances → warn before final close. |

---

## 16. Inventory Valuation (Tính giá hàng tồn kho)

### Periodic Calculation Run

```
LOCK → LOAD → CALCULATE → APPLY → POST_ADJUSTMENT → COMPLETE
```

### Step 16.1: Lock Transactions

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: Period CLOSING started. Post: TRANSACTIONS_LOCKED. Trigger: Period-end valuation run initiated. |
| **Actor / Persona** | System (auto), Chief Accountant (initiate) |
| **Input Documents / Data** | Company ID, period (year+month), warehouse list. |
| **Business Validations & Rules** | All GRNs/DNs/Transfers for period must be POSTED (not DRAFT). No open stock takes. Period not already valued. |
| **Output / Artifact Generated** | Transaction lock. No new inventory movements allowed for this period until valuation completes. |
| **Exception / Failure Handling** | Open DRAFT documents → block, list open docs for user to post/cancel. Stock take in progress → block. |

### Step 16.2: Calculate Cost by Method

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: TRANSACTIONS_LOCKED. Post: COST_CALCULATED. |
| **Actor / Persona** | System (engine) |
| **Input Documents / Data** | All inventory transactions for period, opening balances, receipts, issues. Per-item cost method from Item master. |
| **Business Validations & Rules** | **Per method**: |
| | **WEIGHTED_AVG (periodic):** `Unit Cost = (Open Value + Receipt Value) / (Open Qty + Receipt Qty)`. Apply to all issues in period. |
| | **WEIGHTED_AVG (perpetual/moving):** `New Avg = (Current Value + Receipt Value) / (Current Qty + Receipt Qty)`. Cost per issue varies by date. |
| | **FIFO:** Sort receipt layers by date. Issue from earliest layer. Remaining layers = closing stock value. |
| | **SPECIFIC_ID:** Match issue to specific receipt via serial# or lot#. Exact cost. |
| | **STANDARD:** Issue at standard cost. Variance accumulated. |
| **Output / Artifact Generated** | Per-item cost calculation report. Closing stock value by item+warehouse. COGS for period. |
| **Exception / Failure Handling** | Missing cost on receipt → use last purchase price, flag. FIFO layers exhausted (negative stock) → use last receipt cost, flag. Zero quantity but non-zero value → adjust to zero, post to 632. |

### Step 16.3: Post Cost Adjustment

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: COST_CALCULATED. Post: COMPLETED. Trigger: Calculated cost differs from provisional cost at issue time. |
| **Actor / Persona** | System (auto) |
| **Input Documents / Data** | Difference between provisional COGS and calculated COGS. |
| **Business Validations & Rules** | Immaterial variance (< 0.1% of COGS) → skip. |
| **Output / Artifact Generated** | **Per TT 99/2025 (direct adjustment):** If calculated COGS > provisional → `Dr 632 Cr 152/156`. If calculated < provisional → reverse. StockBalance closing value updated. `InventoryValuationRun` record completed. |
| **Exception / Failure Handling** | Period-end close blocked until valuation completed. Immaterial variance → skip posting (configurable). |

### Step 16.4: Rollback (if needed)

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: COMPLETED. Post: ROLLED_BACK. Trigger: Period reopened for correction. |
| **Actor / Persona** | Chief Accountant (only) |
| **Input Documents / Data** | Reversal request, reason. |
| **Business Validations & Rules** | Only if period not yet closed. All downstream reports must be re-run. |
| **Output / Artifact Generated** | All cost adjustment entries reversed. Valuation run marked ROLLED_BACK. Transactions unlocked. |
| **Exception / Failure Handling** | Period already closed → prior-period adjustment required. |

---

## 17. Period-End Provision — TK 2294 (Dự phòng giảm giá hàng tồn kho)

```
CalculateNRV → CompareCost → PostProvision → COMPLETE
```

### Step 17.1: Calculate Net Realizable Value

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: Inventory valuation completed. Post: NRV_CALCULATED. Trigger: Period-end, before financial statements. |
| **Actor / Persona** | Chief Accountant, System |
| **Input Documents / Data** | Per item: current selling price, estimated cost to complete, estimated selling cost. |
| **Business Validations & Rules** | **VAS 02, IAS 2**: NRV = estimated selling price - estimated costs to complete - estimated selling costs. NRV calculated per item, not aggregated. |
| **Output / Artifact Generated** | NRV per item. Comparison: cost > NRV → provision needed. |
| **Exception / Failure Handling** | No selling price data for slow-moving items → use last sale price or market data. |

### Step 17.2: Post Provision

| Dimension | Detail |
|---|---|
| **State & Lifecycle Trigger** | Pre: NRV_CALCULATED. Post: PROVISION_POSTED. |
| **Actor / Persona** | System (auto), Chief Accountant (approve) |
| **Input Documents / Data** | Provision calculation: `provision = (cost - NRV) × qty_on_hand`. Compare with existing provision balance on 2294. |
| **Business Validations & Rules** | Provision increase: `Dr 632 Cr 2294`. Provision decrease (reversal): `Dr 2294 Cr 632`. 2294 is contra-asset (credit balance). |
| **Output / Artifact Generated** | **GL per TT 99/2025**: `Dr 632 Cr 2294` (increase) or reverse. 2294 balance updated. Financial statements: inventory at lower of cost or NRV. |
| **Exception / Failure Handling** | Existing provision > needed → reverse to income. Zero provision for zero-stock items → skip. |

---

## 18. GL Posting Reference (Circular 99/2025)

**KEY CHANGE**: TT 99/2025 eliminates TK 611. All inventory transactions post directly to inventory accounts (151-158).

### Movement Type → GL Posting Matrix

| # | Movement Type | Direction | Debit | Credit | Condition |
|---|---|---|---|---|---|
| 1 | GRN_PURCHASE (goods) | IN | 152/153/156 | 331 (temp) | PO-based, before invoice |
| 2 | GRN_PURCHASE (expense) | IN | 641/642/627 | 331 (temp) | Non-stock PO items |
| 3 | GRN_PRODUCTION | IN | 155/156 | 154 (WIP) | Finished goods from production |
| 4 | GRN_RETURN (sales return) | IN | 152/156 | 632 | Reverse COGS at original cost |
| 5 | GRN_DIRECT (donation) | IN | 152/156 | 711 | Counterpart = 711 |
| 6 | GRN_DIRECT (owner) | IN | 152/156 | 4111 | Counterpart = equity |
| 7 | GRN_DIRECT (purchase) | IN | 152/156 | 111/112/338 | Direct payment |
| 8 | GRN_DIRECT (surplus) | IN | 152/156 | 3381 | Awaiting resolution |
| 9 | GI_SALES | OUT | 632 | 152/156 | COGS at cost |
| 10 | GI_PRODUCTION | OUT | 154 | 152/153 | Raw materials to WIP |
| 11 | GI_RETURN_SUPPLIER (uninv) | OUT | 331 (temp) | 152/156 | Reverse GRN, pre-invoice |
| 12 | GI_RETURN_SUPPLIER (inv) | OUT | 331 (AP) | 152/156 + 1331 | Via credit note, post-invoice |
| 13 | GI_WRITEOFF (damage) | OUT | 632 | 152/156 | Normal loss |
| 14 | GI_WRITEOFF (theft) | OUT | 811 | 152/156 | Other expense |
| 15 | GI_WRITEOFF (admin) | OUT | 642 | 152/156 | Admin expense |
| 16 | GI_SAMPLE | OUT | 6417 | 152/156 | Marketing samples |
| 17 | GI_DONATION | OUT | 811 | 152/156 | Charity |
| 18 | GI_INTERNAL | OUT | 641/642/627 | 152/156 | Internal use |
| 19 | TRANSFER_SHIP (transit) | TRANSIT | 151 | 152/156 | Goods in transit |
| 20 | TRANSFER_RECEIVE | TRANSIT | 152/156 | 151 | Transit → destination |
| 21 | TRANSFER_SIMPLE | TRANSFER | 152/156 (dest) | 152/156 (src) | Same-account transfer |
| 22 | ADJUST_INCREASE | IN | 152/156 | 3381/632/711 | Per reason |
| 23 | ADJUST_DECREASE | OUT | 632/642/811 | 152/156 | Per reason |
| 24 | REVALUATION_UP | VALUE | 152/156 | 611 | Standard cost increase |
| 25 | REVALUATION_DOWN | VALUE | 611 | 152/156 | Standard cost decrease |
| 26 | STOCKTAKE_GAIN | IN | 152/156 | 632/711 | Count surplus |
| 27 | STOCKTAKE_LOSS | OUT | 632/642/811 | 152/156 | Count shortage |
| 28 | PROVISION_INCREASE | VALUE | 632 | 2294 | NRV provision |
| 29 | PROVISION_DECREASE | VALUE | 2294 | 632 | NRV reversal |
| 30 | OPENING_BALANCE | IN | 152/156 | 3388/419 | Initial stock setup |

---

## 19. COA Accounts Mapping (TT 99/2025)

### Inventory Balance Sheet Accounts

| TK | Name (VN) | Name (EN) | Type | Normal | Note |
|---|---|---|---|---|---|
| **151** | Hàng mua đang đi đường | Goods in transit | Asset | Debit | Auto-cleared next period |
| **152** | Nguyên liệu, vật liệu | Raw materials | Asset | Debit | NVL for production |
| **153** | Công cụ, dụng cụ | Tools & instruments | Asset | Debit | Low-value tools |
| **154** | CPSXKD dở dang | WIP | Asset | Debit | Production in progress |
| **155** | Sản phẩm | Finished goods | Asset | Debit | Self-manufactured |
| **1561** | Hàng hóa - giá mua | Merchandise - purchase price | Asset | Debit | Goods for resale |
| **1562** | Hàng hóa - chi phí mua | Merchandise - transport cost | Asset | Debit | Allocated freight/duty |
| **157** | Hàng gửi đi bán | Consignment goods | Asset | Debit | Sent to agents/consignees |
| **158** | NVL tại kho bảo thuế | Bonded warehouse goods | Asset | Debit | **New in TT 99/2025** |
| **2294** | Dự phòng giảm giá HTK | Inventory provision | Contra-asset | Credit | NRV provision |

### Inventory Income Statement Accounts

| TK | Name (VN) | Name (EN) | Type | Normal |
|---|---|---|---|---|
| **611** | Chênh lệch giá | Cost variance | Expense | Debit |
| **632** | Giá vốn hàng bán | COGS | Expense | Debit |
| **6417** | Chi phí quảng cáo, khuyến mại | Marketing/samples | Expense | Debit |
| **642** | Chi phí quản lý | Admin expense | Expense | Debit |
| **711** | Thu nhập khác | Other income | Revenue | Credit |
| **811** | Chi phí khác | Other expense | Expense | Debit |
| **3381** | Tài sản thừa chờ xử lý | Surplus awaiting resolution | Liability | Credit |
| **3388** | Phải thu khác | Other receivables | Asset | Debit |
| **419** | Lợi nhuận chưa phân phối | Retained earnings | Equity | Credit/Cr |

---

## 20. State Machine Reference

### Warehouse

```
ACTIVE ←→ SUSPENDED (only if stock=0)
```

### Item

```
ACTIVE ←→ INACTIVE (only if stock=0 in all warehouses)
```

### GRN (Goods Receipt Note)

```
DRAFT → POSTED → CANCELLED
(POSTED: stock +, GL posted)
```

### DN (Delivery Note)

```
DRAFT → POSTED → CANCELLED
(POSTED: stock -, COGS posted)
```

### Stock Transfer

```
REQUESTED → APPROVED → PICKED → IN_TRANSIT → RECEIVED → COMPLETED
    │                                                   
    └──→ CANCELLED (any pre-IN_TRANSIT state)
         (CANCELLED after PICKED → restore source committed_qty)
```

### Stock Adjustment

```
DRAFT → APPROVED → POSTED
```

### Stock Take

```
PLANNED → IN_PROGRESS → APPROVED → COMPLETED
    │                      │
    └──→ CANCELLED         └──→ CORRECTED (if re-count needed after approval)
```

### Cost Revaluation

```
DRAFT → POSTED
```

### Inventory Valuation Run

```
PENDING → IN_PROGRESS → COMPLETED
   │                        │
   └──→ FAILED              └──→ ROLLED_BACK (if period reopened)
```

---

## Appendix A: Integration Map

| Module | Integration Point | Warehouse Effect |
|---|---|---|
| **Purchase** | GRN posted → Warehouse | Stock +, on_order -, GL posted |
| **Purchase** | Return to supplier | Stock -, GL reversed |
| **Sales** | DN posted → Warehouse | Stock -, COGS posted |
| **Sales** | Credit note (return) | Stock +, COGS reversed |
| **GL** | Auto-account determination | Per movement type (table in §18) |
| **OB** | OpeningBalanceDetail with INVENTORY_ITEM | Initial stock + value |
| **Period Close** | Valuation run | Cost adjustment, provision 2294 |
| **Tax** | Inventory movement + VAT rate | Input VAT, import duty tracking |
| **Company** | Multi-tenant filter | All queries by company_id |
| **Auth** | Role-based access | SoD enforced: creator ≠ approver |

---

## Appendix B: Data Model Entities (Proposed)

Refer to `internal/domain/models_warehouse.go` (new file):

- `Warehouse` — code, name, type, allow_negative_stock, default_cost_method
- `Item` — code, barcode, cost_method, account_id, cogs_account_id, is_serialized, is_lot_tracked
- `ItemCategory` — code, name, parent_id, abc_class
- `StockBalance` — item_id, warehouse_id, qty_on_hand, committed_qty, on_order_qty, avg_cost, stock_value
- `InventoryTransaction` — movement_type, direction, qty, unit_cost, total_value, ref_type, ref_id
- `StockTransfer` — source_wh, dest_wh, status, lines with requested/picked/received/damaged qty
- `StockAdjustment` — reason, status, lines with old/new qty/value
- `StockTake` — warehouse, count_date, status, lines with book/physical/variance qty
- `InventoryValuationRun` — period, status, completed_at

Each entity has corresponding repository interface (in `interfaces.go`) with PG + memory implementations.

---

*Document version: v2.0 — July 2026*  
*Benchmark sources: MISA AMIS Kho hàng (helpamis.misa.vn, 06/2026), FAST Business Online Fast Inventory (fast.com.vn, 03/2026), Bravo 10 ERP (bravo.com.vn), Odoo 19.0 Inventory, SAP MM-IM*  
*Regulatory: Thông tư 99/2025/TT-BTC (effective 01/01/2026), IAS 2, VAS 02, Decree 123/2020/ND-CP*
