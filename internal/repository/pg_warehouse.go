package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gotax/internal/domain"
)

type PGWarehouseRepo struct {
	pool *pgxpool.Pool
}

func NewPGWarehouseRepo(pool *pgxpool.Pool) *PGWarehouseRepo {
	return &PGWarehouseRepo{pool: pool}
}

// ─── Warehouse ────────────────────────────────────────────────────────

func (r *PGWarehouseRepo) CreateWarehouse(ctx context.Context, w *domain.Warehouse) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO warehouses (id,company_id,code,name,address,manager,is_active,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		w.ID, w.CompanyID, w.Code, w.Name, nullStr(w.Address), nullStr(w.Manager), true, now, now,
	).Scan(&w.ID)
}

func (r *PGWarehouseRepo) scanWarehouse(row pgx.Row) (domain.Warehouse, error) {
	var w domain.Warehouse
	err := row.Scan(&w.ID, &w.CompanyID, &w.Code, &w.Name, &w.Address, &w.Manager, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
	return w, err
}

func (r *PGWarehouseRepo) GetWarehouseByID(ctx context.Context, id string) (*domain.Warehouse, error) {
	w, err := r.scanWarehouse(r.pool.QueryRow(ctx,
		`SELECT id,company_id,code,name,address,manager,is_active,created_at,updated_at FROM warehouses WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrWarehouseNotFound }
		return nil, err
	}
	return &w, nil
}

func (r *PGWarehouseRepo) ListWarehouses(ctx context.Context, companyID string) ([]domain.Warehouse, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,code,name,address,manager,is_active,created_at,updated_at FROM warehouses WHERE company_id=$1 ORDER BY code`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.Warehouse
	for rows.Next() {
		w, err := r.scanWarehouse(rows)
		if err != nil { return nil, err }
		out = append(out, w)
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateWarehouse(ctx context.Context, w *domain.Warehouse) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE warehouses SET code=$1,name=$2,address=$3,manager=$4,is_active=$5,updated_at=$6 WHERE id=$7`,
		w.Code, w.Name, nullStr(w.Address), nullStr(w.Manager), w.IsActive, now, w.ID)
	return err
}

func (r *PGWarehouseRepo) DeleteWarehouse(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM warehouses WHERE id=$1`, id)
	return err
}

func (r *PGWarehouseRepo) GetWarehouseByCode(ctx context.Context, companyID, code string) (*domain.Warehouse, error) {
	w, err := r.scanWarehouse(r.pool.QueryRow(ctx,
		`SELECT id,company_id,code,name,address,manager,is_active,created_at,updated_at FROM warehouses WHERE company_id=$1 AND code=$2`, companyID, code))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrWarehouseNotFound }
		return nil, err
	}
	return &w, nil
}

// ─── Item Category ────────────────────────────────────────────────────

func (r *PGWarehouseRepo) CreateCategory(ctx context.Context, c *domain.ItemCategory) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO item_categories (id,company_id,code,name,description,parent_id,is_active,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		c.ID, c.CompanyID, c.Code, c.Name, nullStr(c.Description), nullStr(c.ParentID), true, now, now,
	).Scan(&c.ID)
}

func (r *PGWarehouseRepo) scanCategory(row pgx.Row) (domain.ItemCategory, error) {
	var c domain.ItemCategory
	err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ParentID, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *PGWarehouseRepo) GetCategoryByID(ctx context.Context, id string) (*domain.ItemCategory, error) {
	c, err := r.scanCategory(r.pool.QueryRow(ctx,
		`SELECT id,company_id,code,name,description,parent_id,is_active,created_at,updated_at FROM item_categories WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrCategoryNotFound }
		return nil, err
	}
	return &c, nil
}

func (r *PGWarehouseRepo) ListCategories(ctx context.Context, companyID string) ([]domain.ItemCategory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,code,name,description,parent_id,is_active,created_at,updated_at FROM item_categories WHERE company_id=$1 ORDER BY code`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.ItemCategory
	for rows.Next() {
		c, err := r.scanCategory(rows)
		if err != nil { return nil, err }
		out = append(out, c)
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateCategory(ctx context.Context, c *domain.ItemCategory) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE item_categories SET code=$1,name=$2,description=$3,parent_id=$4,is_active=$5,updated_at=$6 WHERE id=$7`,
		c.Code, c.Name, nullStr(c.Description), nullStr(c.ParentID), c.IsActive, now, c.ID)
	return err
}

func (r *PGWarehouseRepo) DeleteCategory(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM item_categories WHERE id=$1`, id)
	return err
}

func (r *PGWarehouseRepo) GetCategoryByCode(ctx context.Context, companyID, code string) (*domain.ItemCategory, error) {
	c, err := r.scanCategory(r.pool.QueryRow(ctx,
		`SELECT id,company_id,code,name,description,parent_id,is_active,created_at,updated_at FROM item_categories WHERE company_id=$1 AND code=$2`, companyID, code))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrCategoryNotFound }
		return nil, err
	}
	return &c, nil
}

// ─── Item ─────────────────────────────────────────────────────────────

func (r *PGWarehouseRepo) CreateItem(ctx context.Context, i *domain.Item) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO items (id,company_id,code,name,category_id,unit,base_price,min_stock,max_stock,valuation_method,tax_rate,is_active,notes,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING id`,
		i.ID, i.CompanyID, i.Code, i.Name, nullStr(i.CategoryID), i.Unit, i.BasePrice, i.MinStock, i.MaxStock, string(i.ValuationMethod), i.TaxRate, true, nullStr(i.Notes), now, now,
	).Scan(&i.ID)
}

func (r *PGWarehouseRepo) scanItem(row pgx.Row) (domain.Item, error) {
	var i domain.Item
	var vm string
	err := row.Scan(&i.ID, &i.CompanyID, &i.Code, &i.Name, &i.CategoryID, &i.Unit, &i.BasePrice, &i.MinStock, &i.MaxStock, &vm, &i.TaxRate, &i.IsActive, &i.Notes, &i.CreatedAt, &i.UpdatedAt)
	i.ValuationMethod = domain.ValuationMethod(vm)
	return i, err
}

func (r *PGWarehouseRepo) GetItemByID(ctx context.Context, id string) (*domain.Item, error) {
	i, err := r.scanItem(r.pool.QueryRow(ctx,
		`SELECT id,company_id,code,name,category_id,unit,base_price,min_stock,max_stock,valuation_method,tax_rate,is_active,notes,created_at,updated_at FROM items WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrItemNotFound }
		return nil, err
	}
	return &i, nil
}

func (r *PGWarehouseRepo) ListItems(ctx context.Context, companyID string) ([]domain.Item, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,code,name,category_id,unit,base_price,min_stock,max_stock,valuation_method,tax_rate,is_active,notes,created_at,updated_at FROM items WHERE company_id=$1 ORDER BY code`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.Item
	for rows.Next() {
		i, err := r.scanItem(rows)
		if err != nil { return nil, err }
		out = append(out, i)
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateItem(ctx context.Context, i *domain.Item) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE items SET code=$1,name=$2,category_id=$3,unit=$4,base_price=$5,min_stock=$6,max_stock=$7,valuation_method=$8,tax_rate=$9,is_active=$10,notes=$11,updated_at=$12 WHERE id=$13`,
		i.Code, i.Name, nullStr(i.CategoryID), i.Unit, i.BasePrice, i.MinStock, i.MaxStock, string(i.ValuationMethod), i.TaxRate, i.IsActive, nullStr(i.Notes), now, i.ID)
	return err
}

func (r *PGWarehouseRepo) DeleteItem(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM items WHERE id=$1`, id)
	return err
}

func (r *PGWarehouseRepo) GetItemByCode(ctx context.Context, companyID, code string) (*domain.Item, error) {
	i, err := r.scanItem(r.pool.QueryRow(ctx,
		`SELECT id,company_id,code,name,category_id,unit,base_price,min_stock,max_stock,valuation_method,tax_rate,is_active,notes,created_at,updated_at FROM items WHERE company_id=$1 AND code=$2`, companyID, code))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrItemNotFound }
		return nil, err
	}
	return &i, nil
}

// ─── Stock Balance ────────────────────────────────────────────────────

func (r *PGWarehouseRepo) CreateStockBalance(ctx context.Context, b *domain.StockBalance) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO stock_balances (id,company_id,warehouse_id,item_id,period,quantity,unit_cost,total_cost,last_transaction_at,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id`,
		b.ID, b.CompanyID, b.WarehouseID, b.ItemID, b.Period, b.Quantity, b.UnitCost, b.TotalCost, nullTime(b.LastTransactionAt), now, now,
	).Scan(&b.ID)
}

func (r *PGWarehouseRepo) scanBalance(row pgx.Row) (domain.StockBalance, error) {
	var b domain.StockBalance
	err := row.Scan(&b.ID, &b.CompanyID, &b.WarehouseID, &b.ItemID, &b.Period, &b.Quantity, &b.UnitCost, &b.TotalCost, &b.LastTransactionAt, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func (r *PGWarehouseRepo) GetStockBalanceByID(ctx context.Context, id string) (*domain.StockBalance, error) {
	b, err := r.scanBalance(r.pool.QueryRow(ctx,
		`SELECT id,company_id,warehouse_id,item_id,period,quantity,unit_cost,total_cost,last_transaction_at,created_at,updated_at FROM stock_balances WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrBalanceNotFound }
		return nil, err
	}
	return &b, nil
}

func (r *PGWarehouseRepo) FindStockBalance(ctx context.Context, companyID, warehouseID, itemID, period string) (*domain.StockBalance, error) {
	b, err := r.scanBalance(r.pool.QueryRow(ctx,
		`SELECT id,company_id,warehouse_id,item_id,period,quantity,unit_cost,total_cost,last_transaction_at,created_at,updated_at
		 FROM stock_balances WHERE company_id=$1 AND warehouse_id=$2 AND item_id=$3 AND period=$4`,
		companyID, warehouseID, itemID, period))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrBalanceNotFound }
		return nil, err
	}
	return &b, nil
}

func (r *PGWarehouseRepo) ListStockBalances(ctx context.Context, companyID, warehouseID string) ([]domain.StockBalance, error) {
	q := `SELECT id,company_id,warehouse_id,item_id,period,quantity,unit_cost,total_cost,last_transaction_at,created_at,updated_at
		  FROM stock_balances WHERE company_id=$1`
	args := []any{companyID}
	if warehouseID != "" {
		q += ` AND warehouse_id=$2`
		args = append(args, warehouseID)
	}
	q += ` ORDER BY item_id,period`
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.StockBalance
	for rows.Next() {
		b, err := r.scanBalance(rows)
		if err != nil { return nil, err }
		out = append(out, b)
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateStockBalance(ctx context.Context, b *domain.StockBalance) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE stock_balances SET warehouse_id=$1,item_id=$2,period=$3,quantity=$4,unit_cost=$5,total_cost=$6,last_transaction_at=$7,updated_at=$8 WHERE id=$9`,
		b.WarehouseID, b.ItemID, b.Period, b.Quantity, b.UnitCost, b.TotalCost, nullTime(b.LastTransactionAt), now, b.ID)
	return err
}

func (r *PGWarehouseRepo) UpsertStockBalance(ctx context.Context, b *domain.StockBalance) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO stock_balances (id,company_id,warehouse_id,item_id,period,quantity,unit_cost,total_cost,last_transaction_at,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 ON CONFLICT (id) DO UPDATE SET quantity=$6,unit_cost=$7,total_cost=$8,last_transaction_at=$9,updated_at=$11
		 RETURNING id`,
		b.ID, b.CompanyID, b.WarehouseID, b.ItemID, b.Period, b.Quantity, b.UnitCost, b.TotalCost, nullTime(b.LastTransactionAt), now, now,
	).Scan(&b.ID)
}

// ─── Inventory Transaction ────────────────────────────────────────────

func (r *PGWarehouseRepo) CreateInventoryTransaction(ctx context.Context, t *domain.InventoryTransaction) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO inventory_transactions (id,company_id,warehouse_id,item_id,trans_type,ref_type,ref_id,qty_before,quantity,qty_after,unit_cost,total_cost,created_by,created_at,notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING id`,
		t.ID, t.CompanyID, t.WarehouseID, t.ItemID, string(t.TransType), nullStr(t.RefType), nullStr(t.RefID), t.QtyBefore, t.Quantity, t.QtyAfter, t.UnitCost, t.TotalCost, t.CreatedBy, now, nullStr(t.Notes),
	).Scan(&t.ID)
}

func (r *PGWarehouseRepo) scanTxn(row pgx.Row) (domain.InventoryTransaction, error) {
	var t domain.InventoryTransaction
	var tt string
	err := row.Scan(&t.ID, &t.CompanyID, &t.WarehouseID, &t.ItemID, &tt, &t.RefType, &t.RefID, &t.QtyBefore, &t.Quantity, &t.QtyAfter, &t.UnitCost, &t.TotalCost, &t.CreatedBy, &t.CreatedAt, &t.Notes)
	t.TransType = domain.TransactionType(tt)
	return t, err
}

func (r *PGWarehouseRepo) GetInventoryTransactionByID(ctx context.Context, id string) (*domain.InventoryTransaction, error) {
	t, err := r.scanTxn(r.pool.QueryRow(ctx,
		`SELECT id,company_id,warehouse_id,item_id,trans_type,ref_type,ref_id,qty_before,quantity,qty_after,unit_cost,total_cost,created_by,created_at,notes FROM inventory_transactions WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrTransNotFound }
		return nil, err
	}
	return &t, nil
}

func (r *PGWarehouseRepo) ListInventoryTransactions(ctx context.Context, companyID, warehouseID, itemID string, offset, limit int) ([]domain.InventoryTransaction, int, error) {
	where := ` WHERE company_id=$1`
	args := []any{companyID}
	argIdx := 2
	if warehouseID != "" {
		where += fmt.Sprintf(` AND warehouse_id=$%d`, argIdx)
		args = append(args, warehouseID)
		argIdx++
	}
	if itemID != "" {
		where += fmt.Sprintf(` AND item_id=$%d`, argIdx)
		args = append(args, itemID)
		argIdx++
	}
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory_transactions`+where, args...).Scan(&total)
	if err != nil { return nil, 0, err }
	if limit <= 0 { limit = 50 }
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,warehouse_id,item_id,trans_type,ref_type,ref_id,qty_before,quantity,qty_after,unit_cost,total_cost,created_by,created_at,notes FROM inventory_transactions`+where+` ORDER BY created_at DESC LIMIT $`+fmt.Sprintf("%d", argIdx)+` OFFSET $`+fmt.Sprintf("%d", argIdx+1),
		append(args, limit, offset)...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var out []domain.InventoryTransaction
	for rows.Next() {
		t, err := r.scanTxn(rows)
		if err != nil { return nil, 0, err }
		out = append(out, t)
	}
	return out, total, nil
}

// ─── Stock Transfer ───────────────────────────────────────────────────

func (r *PGWarehouseRepo) CreateStockTransfer(ctx context.Context, t *domain.StockTransfer) error {
	now := time.Now().UTC().Format(time.RFC3339)
	err := r.pool.QueryRow(ctx,
		`INSERT INTO stock_transfers (id,company_id,transfer_number,from_warehouse_id,to_warehouse_id,status,transfer_date,created_by,approved_by,approved_at,completed_by,completed_at,cancelled_reason,notes,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING id`,
		t.ID, t.CompanyID, t.TransferNumber, t.FromWarehouseID, t.ToWarehouseID, string(t.Status), nullStr(t.TransferDate), t.CreatedBy, nullStr(t.ApprovedBy), nullStr(t.ApprovedAt), nullStr(t.CompletedBy), nullStr(t.CompletedAt), nullStr(t.CancelledReason), nullStr(t.Notes), now, now,
	).Scan(&t.ID)
	if err != nil { return err }
	for i := range t.Items {
		t.Items[i].TransferID = t.ID
		if err := r.CreateTransferItem(ctx, &t.Items[i]); err != nil { return err }
	}
	return nil
}

func (r *PGWarehouseRepo) scanTransfer(row pgx.Row) (domain.StockTransfer, error) {
	var t domain.StockTransfer
	var st string
	err := row.Scan(&t.ID, &t.CompanyID, &t.TransferNumber, &t.FromWarehouseID, &t.ToWarehouseID, &st, &t.TransferDate, &t.CreatedBy, &t.ApprovedBy, &t.ApprovedAt, &t.CompletedBy, &t.CompletedAt, &t.CancelledReason, &t.Notes, &t.CreatedAt, &t.UpdatedAt)
	t.Status = domain.TransferStatus(st)
	return t, err
}

func (r *PGWarehouseRepo) GetStockTransferByID(ctx context.Context, id string) (*domain.StockTransfer, error) {
	t, err := r.scanTransfer(r.pool.QueryRow(ctx,
		`SELECT id,company_id,transfer_number,from_warehouse_id,to_warehouse_id,status,transfer_date,created_by,approved_by,approved_at,completed_by,completed_at,cancelled_reason,notes,created_at,updated_at FROM stock_transfers WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrTransferNotFound }
		return nil, err
	}
	items, err := r.GetTransferItems(ctx, id)
	if err != nil { return nil, err }
	t.Items = items
	return &t, nil
}

func (r *PGWarehouseRepo) ListStockTransfers(ctx context.Context, companyID string) ([]domain.StockTransfer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,transfer_number,from_warehouse_id,to_warehouse_id,status,transfer_date,created_by,approved_by,approved_at,completed_by,completed_at,cancelled_reason,notes,created_at,updated_at FROM stock_transfers WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.StockTransfer
	for rows.Next() {
		t, err := r.scanTransfer(rows)
		if err != nil { return nil, err }
		out = append(out, t)
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateStockTransfer(ctx context.Context, t *domain.StockTransfer) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE stock_transfers SET from_warehouse_id=$1,to_warehouse_id=$2,status=$3,transfer_date=$4,approved_by=$5,approved_at=$6,completed_by=$7,completed_at=$8,cancelled_reason=$9,notes=$10,updated_at=$11 WHERE id=$12`,
		t.FromWarehouseID, t.ToWarehouseID, string(t.Status), nullStr(t.TransferDate), nullStr(t.ApprovedBy), nullStr(t.ApprovedAt), nullStr(t.CompletedBy), nullStr(t.CompletedAt), nullStr(t.CancelledReason), nullStr(t.Notes), now, t.ID)
	return err
}

func (r *PGWarehouseRepo) UpdateStockTransferStatus(ctx context.Context, id string, status domain.TransferStatus) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx, `UPDATE stock_transfers SET status=$1,updated_at=$2 WHERE id=$3`, string(status), now, id)
	return err
}

func (r *PGWarehouseRepo) GetTransferItems(ctx context.Context, transferID string) ([]domain.TransferItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,transfer_id,item_id,quantity,unit_cost FROM transfer_items WHERE transfer_id=$1`, transferID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TransferItem
	for rows.Next() {
		var it domain.TransferItem
		if err := rows.Scan(&it.ID, &it.TransferID, &it.ItemID, &it.Quantity, &it.UnitCost); err != nil { return nil, err }
		out = append(out, it)
	}
	return out, nil
}

func (r *PGWarehouseRepo) CreateTransferItem(ctx context.Context, it *domain.TransferItem) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO transfer_items (id,transfer_id,item_id,quantity,unit_cost) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		it.ID, it.TransferID, it.ItemID, it.Quantity, it.UnitCost,
	).Scan(&it.ID)
}

// ─── Stock Adjustment ─────────────────────────────────────────────────

func (r *PGWarehouseRepo) CreateStockAdjustment(ctx context.Context, a *domain.StockAdjustment) error {
	now := time.Now().UTC().Format(time.RFC3339)
	err := r.pool.QueryRow(ctx,
		`INSERT INTO stock_adjustments (id,company_id,warehouse_id,adjustment_number,adj_type,reason,status,created_by,approved_by,approved_at,posted_at,rejected_reason,notes,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING id`,
		a.ID, a.CompanyID, a.WarehouseID, a.AdjustmentNumber, string(a.AdjType), nullStr(a.Reason), string(a.Status), a.CreatedBy, nullStr(a.ApprovedBy), nullStr(a.ApprovedAt), nullStr(a.PostedAt), nullStr(a.RejectedReason), nullStr(a.Notes), now, now,
	).Scan(&a.ID)
	if err != nil { return err }
	for i := range a.Items {
		a.Items[i].AdjustmentID = a.ID
		if err := r.CreateAdjustmentItem(ctx, &a.Items[i]); err != nil { return err }
	}
	return nil
}

func (r *PGWarehouseRepo) scanAdjustment(row pgx.Row) (domain.StockAdjustment, error) {
	var a domain.StockAdjustment
	var at, st string
	err := row.Scan(&a.ID, &a.CompanyID, &a.WarehouseID, &a.AdjustmentNumber, &at, &a.Reason, &st, &a.CreatedBy, &a.ApprovedBy, &a.ApprovedAt, &a.PostedAt, &a.RejectedReason, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
	a.AdjType = domain.AdjType(at)
	a.Status = domain.AdjStatus(st)
	return a, err
}

func (r *PGWarehouseRepo) GetStockAdjustmentByID(ctx context.Context, id string) (*domain.StockAdjustment, error) {
	a, err := r.scanAdjustment(r.pool.QueryRow(ctx,
		`SELECT id,company_id,warehouse_id,adjustment_number,adj_type,reason,status,created_by,approved_by,approved_at,posted_at,rejected_reason,notes,created_at,updated_at FROM stock_adjustments WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrAdjNotFound }
		return nil, err
	}
	items, err := r.GetAdjustmentItems(ctx, id)
	if err != nil { return nil, err }
	a.Items = items
	return &a, nil
}

func (r *PGWarehouseRepo) ListStockAdjustments(ctx context.Context, companyID string) ([]domain.StockAdjustment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,warehouse_id,adjustment_number,adj_type,reason,status,created_by,approved_by,approved_at,posted_at,rejected_reason,notes,created_at,updated_at FROM stock_adjustments WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.StockAdjustment
	for rows.Next() {
		a, err := r.scanAdjustment(rows)
		if err != nil { return nil, err }
		out = append(out, a)
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateStockAdjustment(ctx context.Context, a *domain.StockAdjustment) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE stock_adjustments SET warehouse_id=$1,adj_type=$2,reason=$3,status=$4,approved_by=$5,approved_at=$6,posted_at=$7,rejected_reason=$8,notes=$9,updated_at=$10 WHERE id=$11`,
		a.WarehouseID, string(a.AdjType), nullStr(a.Reason), string(a.Status), nullStr(a.ApprovedBy), nullStr(a.ApprovedAt), nullStr(a.PostedAt), nullStr(a.RejectedReason), nullStr(a.Notes), now, a.ID)
	return err
}

func (r *PGWarehouseRepo) UpdateStockAdjustmentStatus(ctx context.Context, id string, status domain.AdjStatus) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx, `UPDATE stock_adjustments SET status=$1,updated_at=$2 WHERE id=$3`, string(status), now, id)
	return err
}

func (r *PGWarehouseRepo) GetAdjustmentItems(ctx context.Context, adjustmentID string) ([]domain.AdjItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,adjustment_id,item_id,qty_before,qty_after,unit_cost,reason FROM adjustment_items WHERE adjustment_id=$1`, adjustmentID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.AdjItem
	for rows.Next() {
		var it domain.AdjItem
		if err := rows.Scan(&it.ID, &it.AdjustmentID, &it.ItemID, &it.QtyBefore, &it.QtyAfter, &it.UnitCost, &it.Reason); err != nil { return nil, err }
		out = append(out, it)
	}
	return out, nil
}

func (r *PGWarehouseRepo) CreateAdjustmentItem(ctx context.Context, it *domain.AdjItem) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO adjustment_items (id,adjustment_id,item_id,qty_before,qty_after,unit_cost,reason) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		it.ID, it.AdjustmentID, it.ItemID, it.QtyBefore, it.QtyAfter, it.UnitCost, nullStr(it.Reason),
	).Scan(&it.ID)
}

// ─── Stock Take ───────────────────────────────────────────────────────

func (r *PGWarehouseRepo) CreateStockTake(ctx context.Context, t *domain.StockTake) error {
	now := time.Now().UTC().Format(time.RFC3339)
	err := r.pool.QueryRow(ctx,
		`INSERT INTO stock_takes (id,company_id,warehouse_id,take_number,status,take_date,created_by,verified_by,verified_at,posted_at,notes,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`,
		t.ID, t.CompanyID, t.WarehouseID, t.TakeNumber, string(t.Status), nullStr(t.TakeDate), t.CreatedBy, nullStr(t.VerifiedBy), nullStr(t.VerifiedAt), nullStr(t.PostedAt), nullStr(t.Notes), now, now,
	).Scan(&t.ID)
	if err != nil { return err }
	for i := range t.Items {
		t.Items[i].TakeID = t.ID
		if err := r.CreateTakeItem(ctx, &t.Items[i]); err != nil { return err }
	}
	return nil
}

func (r *PGWarehouseRepo) scanTake(row pgx.Row) (domain.StockTake, error) {
	var t domain.StockTake
	var st string
	err := row.Scan(&t.ID, &t.CompanyID, &t.WarehouseID, &t.TakeNumber, &st, &t.TakeDate, &t.CreatedBy, &t.VerifiedBy, &t.VerifiedAt, &t.PostedAt, &t.Notes, &t.CreatedAt, &t.UpdatedAt)
	t.Status = domain.TakeStatus(st)
	return t, err
}

func (r *PGWarehouseRepo) GetStockTakeByID(ctx context.Context, id string) (*domain.StockTake, error) {
	t, err := r.scanTake(r.pool.QueryRow(ctx,
		`SELECT id,company_id,warehouse_id,take_number,status,take_date,created_by,verified_by,verified_at,posted_at,notes,created_at,updated_at FROM stock_takes WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrTakeNotFound }
		return nil, err
	}
	items, err := r.GetTakeItems(ctx, id)
	if err != nil { return nil, err }
	t.Items = items
	return &t, nil
}

func (r *PGWarehouseRepo) ListStockTakes(ctx context.Context, companyID string) ([]domain.StockTake, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,warehouse_id,take_number,status,take_date,created_by,verified_by,verified_at,posted_at,notes,created_at,updated_at FROM stock_takes WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.StockTake
	for rows.Next() {
		t, err := r.scanTake(rows)
		if err != nil { return nil, err }
		out = append(out, t)
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateStockTake(ctx context.Context, t *domain.StockTake) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE stock_takes SET warehouse_id=$1,status=$2,take_date=$3,verified_by=$4,verified_at=$5,posted_at=$6,notes=$7,updated_at=$8 WHERE id=$9`,
		t.WarehouseID, string(t.Status), nullStr(t.TakeDate), nullStr(t.VerifiedBy), nullStr(t.VerifiedAt), nullStr(t.PostedAt), nullStr(t.Notes), now, t.ID)
	return err
}

func (r *PGWarehouseRepo) UpdateStockTakeStatus(ctx context.Context, id string, status domain.TakeStatus) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx, `UPDATE stock_takes SET status=$1,updated_at=$2 WHERE id=$3`, string(status), now, id)
	return err
}

func (r *PGWarehouseRepo) GetTakeItems(ctx context.Context, takeID string) ([]domain.TakeItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,take_id,item_id,expected_qty,actual_qty,unit_cost,variance,notes FROM take_items WHERE take_id=$1`, takeID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.TakeItem
	for rows.Next() {
		var it domain.TakeItem
		if err := rows.Scan(&it.ID, &it.TakeID, &it.ItemID, &it.ExpectedQty, &it.ActualQty, &it.UnitCost, &it.Variance, &it.Notes); err != nil { return nil, err }
		out = append(out, it)
	}
	return out, nil
}

func (r *PGWarehouseRepo) CreateTakeItem(ctx context.Context, it *domain.TakeItem) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO take_items (id,take_id,item_id,expected_qty,actual_qty,unit_cost,variance,notes) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		it.ID, it.TakeID, it.ItemID, it.ExpectedQty, it.ActualQty, it.UnitCost, it.Variance, nullStr(it.Notes),
	).Scan(&it.ID)
}

// ─── Inventory Valuation Run ──────────────────────────────────────────

func (r *PGWarehouseRepo) CreateValuationRun(ctx context.Context, v *domain.InventoryValuationRun) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return r.pool.QueryRow(ctx,
		`INSERT INTO valuation_runs (id,company_id,valuation_date,method,status,created_by,completed_at,error_log,notes,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		v.ID, v.CompanyID, v.ValuationDate, string(v.Method), string(v.Status), v.CreatedBy, nullStr(v.CompletedAt), nullStr(v.ErrorLog), nullStr(v.Notes), now,
	).Scan(&v.ID)
}

func (r *PGWarehouseRepo) scanValuationRun(row pgx.Row) (domain.InventoryValuationRun, error) {
	var v domain.InventoryValuationRun
	var vm, st string
	err := row.Scan(&v.ID, &v.CompanyID, &v.ValuationDate, &vm, &st, &v.CreatedBy, &v.CompletedAt, &v.ErrorLog, &v.Notes, &v.CreatedAt)
	v.Method = domain.ValuationMethod(vm)
	v.Status = domain.ValuationRunStatus(st)
	return v, err
}

func (r *PGWarehouseRepo) GetValuationRunByID(ctx context.Context, id string) (*domain.InventoryValuationRun, error) {
	v, err := r.scanValuationRun(r.pool.QueryRow(ctx,
		`SELECT id,company_id,valuation_date,method,status,created_by,completed_at,error_log,notes,created_at FROM valuation_runs WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows { return nil, domain.ErrValRunNotFound }
		return nil, err
	}
	return &v, nil
}

func (r *PGWarehouseRepo) ListValuationRuns(ctx context.Context, companyID string) ([]domain.InventoryValuationRun, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,company_id,valuation_date,method,status,created_by,completed_at,error_log,notes,created_at FROM valuation_runs WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []domain.InventoryValuationRun
	for rows.Next() {
		v, err := r.scanValuationRun(rows)
		if err != nil { return nil, err }
		out = append(out, v)
	}
	return out, nil
}

func (r *PGWarehouseRepo) UpdateValuationRun(ctx context.Context, v *domain.InventoryValuationRun) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.pool.Exec(ctx,
		`UPDATE valuation_runs SET valuation_date=$1,method=$2,status=$3,completed_at=$4,error_log=$5,notes=$6,created_at=$7 WHERE id=$8`,
		v.ValuationDate, string(v.Method), string(v.Status), nullStr(v.CompletedAt), nullStr(v.ErrorLog), nullStr(v.Notes), now, v.ID)
	return err
}
