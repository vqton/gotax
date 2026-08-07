package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type WarehouseKeeperHandler struct {
	svc *service.WarehouseKeeperService
}

func NewWarehouseKeeperHandler(svc *service.WarehouseKeeperService) *WarehouseKeeperHandler {
	return &WarehouseKeeperHandler{svc: svc}
}

func RegisterKeeperRoutes(r *gin.Engine, h *WarehouseKeeperHandler, authMW gin.HandlerFunc) {
	k := r.Group("/api/v1/warehouse/keeper", authMW)
	{
		assignments := k.Group("/assignments")
		{
			assignments.POST("", h.CreateAssignment)
			assignments.GET("", h.ListAssignments)
			assignments.GET("/:id", h.GetAssignment)
			assignments.PUT("/:id", h.UpdateAssignment)
			assignments.DELETE("/:id", h.DeleteAssignment)
		}
		ledger := k.Group("/ledger")
		{
			ledger.GET("", h.ListLedgerEntries)
			ledger.GET("/:id", h.GetLedgerEntry)
			ledger.POST("/record", h.RecordSlips)
			ledger.POST("/unrecord", h.UnrecordEntry)
			ledger.GET("/balance", h.GetLedgerBalance)
		}
		k.GET("/pending-slips", h.GetPendingSlips)
		k.GET("/pending-slips/count", h.GetPendingSlipsCount)
		reconciliation := k.Group("/reconciliation")
		{
			reconciliation.GET("", h.GetReconciliationReport)
		}
		stockCard := k.Group("/stock-card")
		{
			stockCard.GET("", h.GetStockCard)
		}
		reports := k.Group("/reports")
		{
			reports.GET("/inventory-summary", h.GetInventorySummary)
		}
	}
}

// ─── Assignment ─────────────────────────────────────────────────────────────

func (h *WarehouseKeeperHandler) CreateAssignment(c *gin.Context) {
	var req struct {
		WarehouseID   string `json:"warehouse_id" binding:"required"`
		UserID        string `json:"user_id" binding:"required"`
		Role          string `json:"role" binding:"required"`
		EffectiveFrom string `json:"effective_from" binding:"required"`
		EffectiveTo   string `json:"effective_to"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effective_from date"})
		return
	}
	a := &domain.WarehouseKeeperAssignment{
		CompanyID:     c.GetString("company_id"),
		WarehouseID:   req.WarehouseID,
		UserID:        req.UserID,
		Role:          domain.KeeperRole(req.Role),
		EffectiveFrom: from,
	}
	if req.EffectiveTo != "" {
		to, err := time.Parse("2006-01-02", req.EffectiveTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effective_to date"})
			return
		}
		a.EffectiveTo = &to
	}
	userID := c.GetString("user_id")
	if err := h.svc.CreateAssignment(c.Request.Context(), a, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *WarehouseKeeperHandler) ListAssignments(c *gin.Context) {
	companyID := c.GetString("company_id")
	assignments, err := h.svc.ListAssignments(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

func (h *WarehouseKeeperHandler) GetAssignment(c *gin.Context) {
	id := c.Param("id")
	a, err := h.svc.GetAssignment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *WarehouseKeeperHandler) UpdateAssignment(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		WarehouseID   string `json:"warehouse_id"`
		UserID        string `json:"user_id"`
		Role          string `json:"role"`
		EffectiveFrom string `json:"effective_from"`
		EffectiveTo   string `json:"effective_to"`
		IsActive      *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.svc.GetAssignment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if req.WarehouseID != "" {
		a.WarehouseID = req.WarehouseID
	}
	if req.UserID != "" {
		a.UserID = req.UserID
	}
	if req.Role != "" {
		a.Role = domain.KeeperRole(req.Role)
	}
	if req.EffectiveFrom != "" {
		from, err := time.Parse("2006-01-02", req.EffectiveFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effective_from date"})
			return
		}
		a.EffectiveFrom = from
	}
	if req.EffectiveTo != "" {
		to, err := time.Parse("2006-01-02", req.EffectiveTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effective_to date"})
			return
		}
		a.EffectiveTo = &to
	}
	if req.IsActive != nil {
		a.IsActive = *req.IsActive
	}
	if err := h.svc.UpdateAssignment(c.Request.Context(), a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *WarehouseKeeperHandler) DeleteAssignment(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteAssignment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ─── Stock Ledger ───────────────────────────────────────────────────────────

func (h *WarehouseKeeperHandler) ListLedgerEntries(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.Query("warehouse_id")
	if warehouseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "warehouse_id is required"})
		return
	}
	filter := domain.LedgerFilter{
		CompanyID:   companyID,
		WarehouseID: warehouseID,
		ItemID:      c.Query("item_id"),
		VoucherType: c.Query("voucher_type"),
		Status:      c.Query("status"),
		Page:        1,
		PageSize:    50,
	}
	if from := c.Query("from"); from != "" {
		t, _ := time.Parse("2006-01-02", from)
		filter.From = t
	}
	if to := c.Query("to"); to != "" {
		t, _ := time.Parse("2006-01-02", to)
		filter.To = t
	}
	entries, total, err := h.svc.ListLedgerEntries(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": entries, "total": total})
}

func (h *WarehouseKeeperHandler) GetLedgerEntry(c *gin.Context) {
	id := c.Param("id")
	entry, err := h.svc.GetLedgerEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *WarehouseKeeperHandler) RecordSlips(c *gin.Context) {
	companyID := c.GetString("company_id")
	var req struct {
		WarehouseID string `json:"warehouse_id" binding:"required"`
		Slips       []struct {
			ItemID       string  `json:"item_id" binding:"required"`
			VoucherType  string  `json:"voucher_type" binding:"required"`
			VoucherNo    string  `json:"voucher_no"`
			VoucherRefID string  `json:"voucher_ref_id"`
			Description  string  `json:"description"`
			ReceiptQty   float64 `json:"receipt_qty"`
			IssueQty     float64 `json:"issue_qty"`
		} `json:"slips" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var slips []service.RecordSlipRequest
	for _, s := range req.Slips {
		slips = append(slips, service.RecordSlipRequest{
			ItemID:       s.ItemID,
			VoucherType:  domain.LedgerVoucherType(s.VoucherType),
			VoucherNo:    s.VoucherNo,
			VoucherRefID: s.VoucherRefID,
			Description:  s.Description,
			ReceiptQty:   s.ReceiptQty,
			IssueQty:     s.IssueQty,
		})
	}
	userID := c.GetString("user_id")
	if err := h.svc.RecordSlips(c.Request.Context(), companyID, req.WarehouseID, slips, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "recorded"})
}

func (h *WarehouseKeeperHandler) UnrecordEntry(c *gin.Context) {
	var req struct {
		EntryID string `json:"entry_id" binding:"required"`
		Reason  string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	if err := h.svc.UnrecordEntry(c.Request.Context(), req.EntryID, userID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unrecorded"})
}

func (h *WarehouseKeeperHandler) GetLedgerBalance(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.Query("warehouse_id")
	itemID := c.Query("item_id")
	if warehouseID == "" || itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "warehouse_id and item_id are required"})
		return
	}
	balance, err := h.svc.GetLedgerBalance(c.Request.Context(), companyID, warehouseID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// ─── Pending Slips ──────────────────────────────────────────────────────────

func (h *WarehouseKeeperHandler) GetPendingSlips(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.Query("warehouse_id")
	if warehouseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "warehouse_id is required"})
		return
	}
	slips, err := h.svc.GetPendingSlips(c.Request.Context(), companyID, warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": slips})
}

func (h *WarehouseKeeperHandler) GetPendingSlipsCount(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.Query("warehouse_id")
	if warehouseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "warehouse_id is required"})
		return
	}
	count, err := h.svc.GetPendingSlipsCount(c.Request.Context(), companyID, warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// ─── Reconciliation ─────────────────────────────────────────────────────────

func (h *WarehouseKeeperHandler) GetReconciliationReport(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.Query("warehouse_id")
	if warehouseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "warehouse_id is required"})
		return
	}
	now := time.Now()
	from := now.AddDate(0, 0, -30)
	to := now
	if f := c.Query("from"); f != "" {
		from, _ = time.Parse("2006-01-02", f)
	}
	if t := c.Query("to"); t != "" {
		to, _ = time.Parse("2006-01-02", t)
	}
	items, err := h.svc.GetReconciliationReport(c.Request.Context(), companyID, warehouseID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ─── Stock Card ─────────────────────────────────────────────────────────────

func (h *WarehouseKeeperHandler) GetStockCard(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.Query("warehouse_id")
	itemID := c.Query("item_id")
	period := c.Query("period")
	if warehouseID == "" || itemID == "" || period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "warehouse_id, item_id, and period are required"})
		return
	}
	card, err := h.svc.GetStockCard(c.Request.Context(), companyID, warehouseID, itemID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, card)
}

// ─── Keeper Reports ─────────────────────────────────────────────────────────

func (h *WarehouseKeeperHandler) GetInventorySummary(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.Query("warehouse_id")
	if warehouseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "warehouse_id is required"})
		return
	}
	items, err := h.svc.GetInventorySummary(c.Request.Context(), companyID, warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
