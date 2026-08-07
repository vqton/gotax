# Technical Specifications — Thủ Kho (Warehouse Keeper) Module

**Version:** 1.0
**Date:** August 2026
**Circular:** 99/2025/TT-BTC, Law on Accounting 2015

---

## 1. Domain Model

### WarehouseKeeperAssignment

```
WarehouseKeeperAssignment {
  ID              UUID
  CompanyID       UUID
  WarehouseID     UUID          // FK -> Warehouse
  UserID          UUID          // FK -> User
  Role            enum(keeper,manager)
  EffectiveFrom   date
  EffectiveTo     date          // nullable = open-ended
  IsActive        bool
  CreatedBy       UUID
  CreatedAt       datetime
  UpdatedAt       datetime
}
```

### StockLedgerEntry

```
StockLedgerEntry {
  ID              UUID
  CompanyID       UUID
  WarehouseID     UUID
  ItemID          UUID
  EntryDate       date
  VoucherType     enum(receipt,issue,transfer_in,transfer_out,adjustment_in,adjustment_out,count_gain,count_loss)
  VoucherNo       string        // Reference to source document
  VoucherRefID    UUID          // FK -> source (GRN/DN/Transfer/Adjustment/StockTake)
  Description     text
  ReceiptQty      decimal(18,4) // Positive for inbound
  IssueQty        decimal(18,4) // Positive for outbound
  BalanceQty      decimal(18,4) // Running balance after this entry
  UnitCost        decimal(18,4) // Optional, hidden from keeper if configured
  TotalValue      decimal(18,2) // Optional, hidden from keeper if configured
  RecordedBy      UUID          // FK -> User (the keeper)
  RecordedAt      datetime
  UnrecordedBy    UUID          // FK -> User (nullable, set on un-record)
  UnrecordedAt    datetime
  UnrecordReason  text
  Status          enum(recorded,unrecorded)
  CreatedAt       datetime
}
```

### KeeperReconciliationReport

```
// Virtual entity — generated on-demand, not persisted
KeeperReconciliationItem {
  ItemID          UUID
  ItemCode        string
  ItemName        string
  WarehouseID     UUID
  WarehouseName   string
  KeeperQty       decimal(18,4) // From StockLedgerEntry balance
  AccountingQty   decimal(18,4) // From StockBalance.quantity
  VarianceQty     decimal(18,4)
  UnitCost        decimal(18,4)
  VarianceValue   decimal(18,2)
  LastUpdated     datetime
}
```

### KeeperInventoryCount

```
// Extends existing StockTake with keeper-specific fields
KeeperInventoryCount {
  // Inherits all StockTake fields
  StockTakeID     UUID          // FK -> StockTake
  KeeperID        UUID          // FK -> User (assigned keeper)
  BookQtyHidden   bool          // blind count flag
  KeeperCountDate date
  KeeperSignature string        // digital signature or name
  AccountantSignature string
  ManagerSignature string
  Status          enum(pending_keeper, pending_accountant, pending_manager, completed)
}
```

---

## 2. Repository Interface

```go
type WarehouseKeeperRepository interface {
  // Assignment
  CreateAssignment(ctx context.Context, a *domain.WarehouseKeeperAssignment) error
  GetAssignment(ctx context.Context, id uuid.UUID) (*domain.WarehouseKeeperAssignment, error)
  ListAssignments(ctx context.Context, companyID uuid.UUID) ([]domain.WarehouseKeeperAssignment, error)
  GetActiveAssignment(ctx context.Context, companyID, warehouseID uuid.UUID, date time.Time) (*domain.WarehouseKeeperAssignment, error)
  UpdateAssignment(ctx context.Context, a *domain.WarehouseKeeperAssignment) error
  DeleteAssignment(ctx context.Context, id uuid.UUID) error

  // Stock Ledger
  CreateLedgerEntry(ctx context.Context, e *domain.StockLedgerEntry) error
  GetLedgerEntry(ctx context.Context, id uuid.UUID) (*domain.StockLedgerEntry, error)
  ListLedgerEntries(ctx context.Context, filter domain.LedgerFilter) ([]domain.StockLedgerEntry, int, error)
  UnrecordLedgerEntry(ctx context.Context, id uuid.UUID, unrecordedBy uuid.UUID, reason string) error
  GetLedgerBalance(ctx context.Context, companyID, warehouseID, itemID uuid.UUID) (decimal.Decimal, error)

  // Reconciliation
  GetReconciliationReport(ctx context.Context, companyID, warehouseID uuid.UUID, from, to time.Time) ([]domain.KeeperReconciliationItem, error)

  // Keeper Stock Card
  GetStockCard(ctx context.Context, companyID, warehouseID, itemID uuid.UUID, period string) (*domain.StockCard, error)

  // Keeper Reports
  GetKeeperInventorySummary(ctx context.Context, companyID, warehouseID uuid.UUID) ([]domain.KeeperInventorySummaryItem, error)
}

type LedgerFilter struct {
  CompanyID   uuid.UUID
  WarehouseID uuid.UUID
  ItemID      uuid.UUID   // optional
  From        time.Time
  To          time.Time
  VoucherType string      // optional
  Status      string      // optional (recorded/unrecorded)
  Page        int
  PageSize    int
}
```

---

## 3. API Endpoints

### 3.1 Assignment Management

```
POST   /api/v1/warehouse/keeper/assignments         Create assignment
GET    /api/v1/warehouse/keeper/assignments          List assignments
GET    /api/v1/warehouse/keeper/assignments/:id      Get assignment detail
PUT    /api/v1/warehouse/keeper/assignments/:id      Update assignment
DELETE /api/v1/warehouse/keeper/assignments/:id      Delete assignment (soft)
```

### 3.2 Stock Ledger

```
GET    /api/v1/warehouse/keeper/ledger               List ledger entries (filter: warehouse, item, date range, status)
GET    /api/v1/warehouse/keeper/ledger/:id           Get ledger entry detail
POST   /api/v1/warehouse/keeper/ledger/record        Record slip(s) to ledger (bulk)
POST   /api/v1/warehouse/keeper/ledger/unrecord      Un-record ledger entry (with reason)
GET    /api/v1/warehouse/keeper/ledger/balance       Get current balance per item per warehouse
```

### 3.3 Pending Slips

```
GET    /api/v1/warehouse/keeper/pending-slips        List slips awaiting keeper review
GET    /api/v1/warehouse/keeper/pending-slips/count   Count pending slips
```

### 3.4 Reconciliation

```
GET    /api/v1/warehouse/keeper/reconciliation       Get reconciliation report (query: warehouse, from, to)
GET    /api/v1/warehouse/keeper/reconciliation/export  Export reconciliation to Excel
```

### 3.5 Stock Card

```
GET    /api/v1/warehouse/keeper/stock-card           Get stock card (query: warehouse, item, period)
GET    /api/v1/warehouse/keeper/stock-card/print     Print stock card (PDF)
```

### 3.6 Keeper Reports

```
GET    /api/v1/warehouse/keeper/reports/inventory-summary    Inventory summary per warehouse
GET    /api/v1/warehouse/keeper/reports/receipt-issue         Receipt/Issue detail report
GET    /api/v1/warehouse/keeper/reports/count-variance        Count variance report
```

### 3.7 Module Toggle

```
GET    /api/v1/warehouse/keeper/config              Get keeper module config
PUT    /api/v1/warehouse/keeper/config              Update config (enable/disable, cost visibility)
```

---

## 4. Database Schema (PostgreSQL)

```sql
CREATE TABLE warehouse_keeper_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL REFERENCES companies(id),
  warehouse_id UUID NOT NULL REFERENCES warehouses(id),
  user_id UUID NOT NULL REFERENCES users(id),
  role VARCHAR(20) NOT NULL DEFAULT 'keeper',
  effective_from DATE NOT NULL,
  effective_to DATE,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(company_id, warehouse_id, effective_from)
);

CREATE TABLE stock_ledger_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL REFERENCES companies(id),
  warehouse_id UUID NOT NULL REFERENCES warehouses(id),
  item_id UUID NOT NULL REFERENCES items(id),
  entry_date DATE NOT NULL,
  voucher_type VARCHAR(30) NOT NULL,
  voucher_no VARCHAR(50),
  voucher_ref_id UUID,
  description TEXT,
  receipt_qty NUMERIC(18,4) NOT NULL DEFAULT 0,
  issue_qty NUMERIC(18,4) NOT NULL DEFAULT 0,
  balance_qty NUMERIC(18,4) NOT NULL DEFAULT 0,
  unit_cost NUMERIC(18,4),
  total_value NUMERIC(18,2),
  recorded_by UUID NOT NULL REFERENCES users(id),
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  unrecorded_by UUID REFERENCES users(id),
  unrecorded_at TIMESTAMPTZ,
  unrecord_reason TEXT,
  status VARCHAR(20) NOT NULL DEFAULT 'recorded',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ledger_company_warehouse ON stock_ledger_entries(company_id, warehouse_id);
CREATE INDEX idx_ledger_item ON stock_ledger_entries(item_id);
CREATE INDEX idx_ledger_date ON stock_ledger_entries(entry_date);
CREATE INDEX idx_ledger_voucher_ref ON stock_ledger_entries(voucher_ref_id);
CREATE INDEX idx_ledger_recorded_by ON stock_ledger_entries(recorded_by);

CREATE TABLE keeper_inventory_counts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  stock_take_id UUID NOT NULL REFERENCES stock_takes(id),
  keeper_id UUID NOT NULL REFERENCES users(id),
  book_qty_hidden BOOLEAN NOT NULL DEFAULT false,
  keeper_count_date DATE,
  keeper_signature VARCHAR(255),
  accountant_signature VARCHAR(255),
  manager_signature VARCHAR(255),
  status VARCHAR(30) NOT NULL DEFAULT 'pending_keeper',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE warehouse_keeper_config (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL REFERENCES companies(id),
  module_enabled BOOLEAN NOT NULL DEFAULT true,
  cost_price_hidden_from_keeper BOOLEAN NOT NULL DEFAULT false,
  auto_record_on_grn BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(company_id)
);
```

---

## 5. Validation Rules

| Field | Rule |
|-------|------|
| WarehouseID | Must reference an active Warehouse |
| UserID | Must reference an active User with keeper role |
| EffectiveFrom | Cannot be in the past for new assignments |
| EffectiveTo | Must be >= EffectiveFrom if set |
| Overlap | Cannot have two active assignments for same warehouse on same date |
| VoucherType | Must be from predefined enum |
| ReceiptQty | >= 0 |
| IssueQty | >= 0 |
| BalanceQty | Must equal previous balance + receipt - issue |
| UnrecordReason | Required when un-recording |
| Status | Cannot un-record an already unrecorded entry |

---

## 6. GL Posting Integration

### 6.1 Stock Ledger Recording

Recording a slip to the Stock Ledger does NOT create GL entries. GL entries are created by the source module (Purchase for GRN, Sale for DN, etc.). The Stock Ledger is a **parallel record** maintained by the keeper for cross-referencing.

### 6.2 Un-recording

Un-recording a Stock Ledger entry does NOT reverse GL entries. It only marks the keeper's record as "unrecorded" with a reason. The source document remains unchanged.

### 6.3 Inventory Count Adjustment

When a physical count is completed and approved, the adjustment posts to GL:
- Gain: `Dr 152/156 Cr 632/711`
- Loss: `Dr 632/811 Cr 152/156`

This is handled by the existing Stock Take workflow in the warehouse service.

---

## 7. State Machines

### 7.1 Assignment Lifecycle

```
CREATE → ACTIVE → EXPIRED (auto on EffectiveTo)
       → CANCELLED (manual)
```

### 7.2 Ledger Entry Lifecycle

```
PENDING → RECORDED → UNRECORDED (with reason)
```

### 7.3 Inventory Count (Keeper Flow)

```
PLANNED → KEEPER_COUNTING → ACCOUNTANT_REVIEWING → MANAGER_APPROVING → COMPLETED
```

---

## 8. Integration Points

| Module | Integration | Direction |
|--------|------------|-----------|
| Warehouse | Read StockBalance, InventoryTransaction | Inbound |
| Warehouse | Read StockTake for count workflow | Inbound |
| Purchase | Read GRN (Phiếu nhập kho) for keeper review | Inbound |
| Sale | Read DN (Phiếu xuất kho) for keeper review | Inbound |
| Auth | User assignment, role checking | Inbound |
| GL | No direct posting (source modules handle GL) | None |

---

## 9. Report Formats

### 9.1 Sổ kho (Stock Ledger) — TT 99 Format

```
┌─────────────────────────────────────────────────────────────────┐
│                    SỔ KHO (S06-DN)                              │
│ Nhóm Vật tư, hàng hóa: [Category]                             │
│ Kho: [Warehouse]    Đơn vị tính: [Unit]                        │
│ Mã, tên vật tư: [Item Code] - [Item Name]                      │
├──────┬──────────┬──────────────┬──────────┬──────────┬──────────┤
│ Ngày │ Số증 từ  │ Diễn giải    │ Nhập     │ Xuất     │ Tồn     │
├──────┼──────────┼──────────────┼──────────┼──────────┼──────────┤
│ ...  │ ...      │ ...          │ ...      │ ...      │ ...      │
├──────┼──────────┼──────────────┼──────────┼──────────┼──────────┤
│      │          │ Tổng phát sinh│ ...     │ ...      │ ...      │
│      │          │ Tồn kỳ trước │         │          │ ...      │
│      │          │ Tồn cuối kỳ  │         │          │ ...      │
└──────┴──────────┴──────────────┴──────────┴──────────┴──────────┘
```

### 9.2 Biên bản kiểm kê (Inventory Count Minutes)

```
┌─────────────────────────────────────────────────────────────────┐
│               BIÊN BẢN KIỂM KÊ KHO                             │
│ Ngày kiểm kê: [Date]     Kho: [Warehouse]                      │
│ Lý do: [Reason]                                                │
├──────┬──────────┬──────────┬──────────┬──────────┬──────────────┤
│ STT  │ Mã hàng  │ Tên hàng │ Tồn sách │ Tồn thực │ Chênh lệch  │
├──────┼──────────┼──────────┼──────────┼──────────┼──────────────┤
│ ...  │ ...      │ ...      │ ...      │ ...      │ ...          │
└──────┴──────────┴──────────┴──────────┴──────────┴──────────────┘
│ Thủ kho ký: _______    Kế toán ký: _______    QL kho ký: _____ │
```

---

## 10. Open Questions

1. Should auto-recording be supported (system records on GRN/DN post, keeper just verifies)?
2. Should keeper see cost prices by default or hidden by default?
3. Should Stock Ledger entries be printed per-item or per-period?
4. How to handle multi-currency items in the keeper view?
5. Should the keeper module affect existing warehouse RBAC or overlay on top?
