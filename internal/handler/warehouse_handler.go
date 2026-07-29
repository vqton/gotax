package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type WarehouseHandler struct {
	svc *service.WarehouseService
}

func NewWarehouseHandler(svc *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{svc: svc}
}

func RegisterWarehouseRoutes(r *gin.Engine, h *WarehouseHandler, authMW gin.HandlerFunc) {
	wh := r.Group("/api/v1/warehouse", authMW)
	{
		warehouses := wh.Group("/warehouses")
		{
			warehouses.POST("", h.CreateWarehouse)
			warehouses.GET("", h.ListWarehouses)
			warehouses.GET("/:id", h.GetWarehouse)
			warehouses.PUT("/:id", h.UpdateWarehouse)
			warehouses.DELETE("/:id", h.DeleteWarehouse)
		}
		cats := wh.Group("/categories")
		{
			cats.POST("", h.CreateCategory)
			cats.GET("", h.ListCategories)
			cats.GET("/:id", h.GetCategory)
			cats.PUT("/:id", h.UpdateCategory)
			cats.DELETE("/:id", h.DeleteCategory)
		}
		items := wh.Group("/items")
		{
			items.POST("", h.CreateItem)
			items.GET("", h.ListItems)
			items.GET("/:id", h.GetItem)
			items.PUT("/:id", h.UpdateItem)
			items.DELETE("/:id", h.DeleteItem)
		}
		balances := wh.Group("/balances")
		{
			balances.GET("", h.ListStockBalances)
			balances.GET("/:id", h.GetStockBalance)
			balances.GET("/find", h.FindStockBalance)
		}
		txns := wh.Group("/transactions")
		{
			txns.GET("", h.ListInventoryTransactions)
		}
		transfers := wh.Group("/transfers")
		{
			transfers.POST("", h.CreateStockTransfer)
			transfers.GET("", h.ListStockTransfers)
			transfers.GET("/:id", h.GetStockTransfer)
			transfers.PUT("/:id", h.UpdateStockTransfer)
			transfers.PATCH("/:id/submit", h.SubmitStockTransfer)
			transfers.PATCH("/:id/approve", h.ApproveStockTransfer)
			transfers.PATCH("/:id/transfer", h.TransferStockTransfer)
			transfers.PATCH("/:id/complete", h.CompleteStockTransfer)
			transfers.PATCH("/:id/cancel", h.CancelStockTransfer)
		}
		adjustments := wh.Group("/adjustments")
		{
			adjustments.POST("", h.CreateStockAdjustment)
			adjustments.GET("", h.ListStockAdjustments)
			adjustments.GET("/:id", h.GetStockAdjustment)
			adjustments.PUT("/:id", h.UpdateStockAdjustment)
			adjustments.PATCH("/:id/submit", h.SubmitStockAdjustment)
			adjustments.PATCH("/:id/approve", h.ApproveStockAdjustment)
			adjustments.PATCH("/:id/post", h.PostStockAdjustment)
			adjustments.PATCH("/:id/reject", h.RejectStockAdjustment)
		}
		takes := wh.Group("/takes")
		{
			takes.POST("", h.CreateStockTake)
			takes.GET("", h.ListStockTakes)
			takes.GET("/:id", h.GetStockTake)
			takes.PUT("/:id", h.UpdateStockTake)
			takes.PATCH("/:id/start", h.StartStockTake)
			takes.PATCH("/:id/complete", h.CompleteStockTake)
			takes.PATCH("/:id/verify", h.VerifyStockTake)
			takes.PATCH("/:id/post", h.PostStockTake)
		}
		valuations := wh.Group("/valuations")
		{
			valuations.POST("", h.CreateValuationRun)
			valuations.GET("", h.ListValuationRuns)
			valuations.GET("/:id", h.GetValuationRun)
			valuations.POST("/:id/run", h.RunValuation)
		}
	}
}

// ─── Warehouse ────────────────────────────────────────────────────────

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var w domain.Warehouse
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	w.CompanyID = c.Query("company_id")
	if w.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	if err := h.svc.CreateWarehouse(c.Request.Context(), &w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *WarehouseHandler) GetWarehouse(c *gin.Context) {
	id := c.Param("id")
	w, err := h.svc.GetWarehouse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *WarehouseHandler) ListWarehouses(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, err := h.svc.ListWarehouses(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	id := c.Param("id")
	var w domain.Warehouse
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	w.ID = id
	if err := h.svc.UpdateWarehouse(c.Request.Context(), &w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteWarehouse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ─── Category ─────────────────────────────────────────────────────────

func (h *WarehouseHandler) CreateCategory(c *gin.Context) {
	var cat domain.ItemCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.CompanyID = c.Query("company_id")
	if cat.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	if err := h.svc.CreateCategory(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *WarehouseHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	cat, err := h.svc.GetCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *WarehouseHandler) ListCategories(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, err := h.svc.ListCategories(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *WarehouseHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat domain.ItemCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.ID = id
	if err := h.svc.UpdateCategory(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *WarehouseHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ─── Item ─────────────────────────────────────────────────────────────

func (h *WarehouseHandler) CreateItem(c *gin.Context) {
	var i domain.Item
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	i.CompanyID = c.Query("company_id")
	if i.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	if err := h.svc.CreateItem(c.Request.Context(), &i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, i)
}

func (h *WarehouseHandler) GetItem(c *gin.Context) {
	id := c.Param("id")
	i, err := h.svc.GetItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, i)
}

func (h *WarehouseHandler) ListItems(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, err := h.svc.ListItems(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *WarehouseHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var i domain.Item
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	i.ID = id
	if err := h.svc.UpdateItem(c.Request.Context(), &i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, i)
}

func (h *WarehouseHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ─── Stock Balance ────────────────────────────────────────────────────

func (h *WarehouseHandler) GetStockBalance(c *gin.Context) {
	id := c.Param("id")
	b, err := h.svc.GetStockBalance(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *WarehouseHandler) FindStockBalance(c *gin.Context) {
	companyID := c.Query("company_id")
	warehouseID := c.Query("warehouse_id")
	itemID := c.Query("item_id")
	period := c.Query("period")
	if companyID == "" || warehouseID == "" || itemID == "" || period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id, warehouse_id, item_id, period query params required"})
		return
	}
	b, err := h.svc.FindStockBalance(c.Request.Context(), companyID, warehouseID, itemID, period)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *WarehouseHandler) ListStockBalances(c *gin.Context) {
	companyID := c.Query("company_id")
	warehouseID := c.Query("warehouse_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, err := h.svc.ListStockBalances(c.Request.Context(), companyID, warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ─── Inventory Transaction ────────────────────────────────────────────

func (h *WarehouseHandler) ListInventoryTransactions(c *gin.Context) {
	companyID := c.Query("company_id")
	warehouseID := c.Query("warehouse_id")
	itemID := c.Query("item_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, total, err := h.svc.ListInventoryTransactions(c.Request.Context(), companyID, warehouseID, itemID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "total": total})
}

// ─── Stock Transfer ───────────────────────────────────────────────────

func (h *WarehouseHandler) CreateStockTransfer(c *gin.Context) {
	var t domain.StockTransfer
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.CompanyID = c.Query("company_id")
	if t.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	t.CreatedBy = c.GetString("user_id")
	if err := h.svc.CreateStockTransfer(c.Request.Context(), &t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *WarehouseHandler) GetStockTransfer(c *gin.Context) {
	id := c.Param("id")
	t, err := h.svc.GetStockTransfer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *WarehouseHandler) ListStockTransfers(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, err := h.svc.ListStockTransfers(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *WarehouseHandler) UpdateStockTransfer(c *gin.Context) {
	id := c.Param("id")
	var t domain.StockTransfer
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ID = id
	if err := h.svc.UpdateStockTransfer(c.Request.Context(), &t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *WarehouseHandler) SubmitStockTransfer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.SubmitStockTransfer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submitted"})
}

func (h *WarehouseHandler) ApproveStockTransfer(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	if err := h.svc.ApproveStockTransfer(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *WarehouseHandler) TransferStockTransfer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.TransferStockTransfer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transferred"})
}

func (h *WarehouseHandler) CompleteStockTransfer(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	if err := h.svc.CompleteStockTransfer(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "completed"})
}

func (h *WarehouseHandler) CancelStockTransfer(c *gin.Context) {
	id := c.Param("id")
	var req struct{ Reason string `json:"reason"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Reason = ""
	}
	if err := h.svc.CancelStockTransfer(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

// ─── Stock Adjustment ─────────────────────────────────────────────────

func (h *WarehouseHandler) CreateStockAdjustment(c *gin.Context) {
	var a domain.StockAdjustment
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.CompanyID = c.Query("company_id")
	if a.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	a.CreatedBy = c.GetString("user_id")
	if err := h.svc.CreateStockAdjustment(c.Request.Context(), &a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *WarehouseHandler) GetStockAdjustment(c *gin.Context) {
	id := c.Param("id")
	a, err := h.svc.GetStockAdjustment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *WarehouseHandler) ListStockAdjustments(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, err := h.svc.ListStockAdjustments(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *WarehouseHandler) UpdateStockAdjustment(c *gin.Context) {
	id := c.Param("id")
	var a domain.StockAdjustment
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.ID = id
	if err := h.svc.UpdateStockAdjustment(c.Request.Context(), &a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *WarehouseHandler) SubmitStockAdjustment(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.SubmitStockAdjustment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submitted"})
}

func (h *WarehouseHandler) ApproveStockAdjustment(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	if err := h.svc.ApproveStockAdjustment(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *WarehouseHandler) PostStockAdjustment(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.PostStockAdjustment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posted"})
}

func (h *WarehouseHandler) RejectStockAdjustment(c *gin.Context) {
	id := c.Param("id")
	var req struct{ Reason string `json:"reason"` }
	c.ShouldBindJSON(&req)
	if err := h.svc.RejectStockAdjustment(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

// ─── Stock Take ───────────────────────────────────────────────────────

func (h *WarehouseHandler) CreateStockTake(c *gin.Context) {
	var t domain.StockTake
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.CompanyID = c.Query("company_id")
	if t.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	t.CreatedBy = c.GetString("user_id")
	if err := h.svc.CreateStockTake(c.Request.Context(), &t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *WarehouseHandler) GetStockTake(c *gin.Context) {
	id := c.Param("id")
	t, err := h.svc.GetStockTake(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *WarehouseHandler) ListStockTakes(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, err := h.svc.ListStockTakes(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *WarehouseHandler) UpdateStockTake(c *gin.Context) {
	id := c.Param("id")
	var t domain.StockTake
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ID = id
	if err := h.svc.UpdateStockTake(c.Request.Context(), &t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *WarehouseHandler) StartStockTake(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.StartStockTake(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "started"})
}

func (h *WarehouseHandler) CompleteStockTake(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.CompleteStockTake(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "completed"})
}

func (h *WarehouseHandler) VerifyStockTake(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	if err := h.svc.VerifyStockTake(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "verified"})
}

func (h *WarehouseHandler) PostStockTake(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.PostStockTake(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posted"})
}

// ─── Valuation Run ────────────────────────────────────────────────────

func (h *WarehouseHandler) CreateValuationRun(c *gin.Context) {
	var v domain.InventoryValuationRun
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v.CompanyID = c.Query("company_id")
	if v.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	v.CreatedBy = c.GetString("user_id")
	if err := h.svc.CreateValuationRun(c.Request.Context(), &v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *WarehouseHandler) GetValuationRun(c *gin.Context) {
	id := c.Param("id")
	v, err := h.svc.GetValuationRun(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *WarehouseHandler) ListValuationRuns(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	list, err := h.svc.ListValuationRuns(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *WarehouseHandler) RunValuation(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RunValuation(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "valuation completed"})
}
