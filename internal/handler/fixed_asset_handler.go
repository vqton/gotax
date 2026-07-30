package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type FAHandler struct {
	svc service.FAServiceInterface
}

func NewFAHandler(svc service.FAServiceInterface) *FAHandler {
	return &FAHandler{svc: svc}
}

func RegisterFixedAssetRoutes(r *gin.Engine, h *FAHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	fa := v1.Group("/fixed-assets")
	{
		categories := fa.Group("/categories")
		{
			categories.POST("", h.CreateCategory)
			categories.GET("", h.ListCategories)
			categories.GET("/:id", h.GetCategory)
			categories.PUT("/:id", h.UpdateCategory)
			categories.DELETE("/:id", h.DeleteCategory)
		}

		assets := fa.Group("")
		{
			assets.POST("", h.CreateAsset)
			assets.GET("", h.ListAssets)
			assets.GET("/:id", h.GetAsset)
			assets.PUT("/:id", h.UpdateAsset)
			assets.DELETE("/:id", h.DeleteAsset)
			assets.PATCH("/:id/status", h.ChangeAssetStatus)

			assets.POST("/:id/depreciate", h.RunDepreciationForAsset)
			assets.GET("/:id/depreciation", h.GetDepreciationByAsset)
			assets.GET("/:id/transactions", h.GetTransactions)
			assets.GET("/:id/allocations", h.GetAllocations)
			assets.POST("/:id/allocations", h.SetAllocations)

			assets.POST("/:id/dispose", h.DisposeAsset)
			assets.POST("/:id/sell", h.SellAsset)
			assets.POST("/:id/adjust", h.AdjustAsset)
			assets.POST("/:id/revalue", h.RevalueAsset)
			assets.POST("/:id/impair", h.ImpairAsset)
			assets.POST("/:id/transfer", h.TransferAsset)
			assets.POST("/:id/cip-transfer", h.CIPTransfer)
			assets.POST("/:id/suspend", h.SuspendDepreciation)
			assets.POST("/:id/resume", h.ResumeDepreciation)
		}

		depr := fa.Group("/depreciation")
		{
			depr.POST("/run", h.RunDepreciation)
			depr.GET("/period/:periodID", h.GetDepreciationByPeriod)
		}

		inventory := fa.Group("/inventory")
		{
			inventory.POST("/plans", h.CreateInventoryPlan)
			inventory.GET("/plans", h.ListInventoryPlans)
			inventory.GET("/plans/:id", h.GetInventoryPlan)
			inventory.PUT("/plans/:id", h.UpdateInventoryPlan)
			inventory.POST("/results", h.CreateInventoryResult)
			inventory.GET("/plans/:id/results", h.GetInventoryResults)
		}
	}
}

func (h *FAHandler) faError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrFANotFound),
		errors.Is(err, domain.ErrFACategoryNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrFACodeExists),
		errors.Is(err, domain.ErrFACategoryCodeExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrFANotActive),
		errors.Is(err, domain.ErrFADepreciationExists),
		errors.Is(err, domain.ErrFAAllocationPctSum),
		errors.Is(err, domain.ErrFATransferSameDepartment),
		errors.Is(err, domain.ErrFASuspensionNotApplicable),
		errors.Is(err, domain.ErrFAResumeNotApplicable),
		errors.Is(err, domain.ErrFACIPTransferNotApplicable):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// ─── Categories ──────────────────────────────────────────────────────

func (h *FAHandler) CreateCategory(c *gin.Context) {
	var req domain.FixedAssetCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	if err := h.svc.CreateCategory(c.Request.Context(), &req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *FAHandler) GetCategory(c *gin.Context) {
	cat, err := h.svc.GetCategory(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *FAHandler) ListCategories(c *gin.Context) {
	filter := domain.FACategoryFilter{CompanyID: c.Query("company_id")}
	cats, err := h.svc.ListCategories(c.Request.Context(), filter)
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *FAHandler) UpdateCategory(c *gin.Context) {
	var req domain.FixedAssetCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	if err := h.svc.UpdateCategory(c.Request.Context(), &req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *FAHandler) DeleteCategory(c *gin.Context) {
	if err := h.svc.DeleteCategory(c.Request.Context(), c.Param("id")); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ─── Assets ──────────────────────────────────────────────────────────

func (h *FAHandler) CreateAsset(c *gin.Context) {
	var req domain.FixedAsset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.CreateAsset(c.Request.Context(), &req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *FAHandler) GetAsset(c *gin.Context) {
	a, err := h.svc.GetAsset(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *FAHandler) ListAssets(c *gin.Context) {
	filter := domain.FAListFilter{
		CompanyID: c.Query("company_id"),
		Keyword:   c.Query("keyword"),
	}
	if s := c.Query("status"); s != "" {
		st := domain.FixedAssetStatus(s)
		filter.Status = &st
	}
	if cat := c.Query("category_id"); cat != "" {
		filter.CategoryID = &cat
	}
	if dept := c.Query("department_id"); dept != "" {
		filter.DepartmentID = &dept
	}
	filter.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter.Offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	assets, total, err := h.svc.ListAssets(c.Request.Context(), filter)
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assets, "total": total})
}

func (h *FAHandler) UpdateAsset(c *gin.Context) {
	var req domain.FixedAsset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	req.UpdatedBy = GetUserID(c)
	if err := h.svc.UpdateAsset(c.Request.Context(), &req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *FAHandler) DeleteAsset(c *gin.Context) {
	if err := h.svc.DeleteAsset(c.Request.Context(), c.Param("id")); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *FAHandler) ChangeAssetStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.svc.ChangeAssetStatus(c.Request.Context(), c.Param("id"), domain.FixedAssetStatus(req.Status), GetUserID(c)); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

// ─── Depreciation ────────────────────────────────────────────────────

func (h *FAHandler) RunDepreciation(c *gin.Context) {
	var req domain.FARunDepreciationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	req.CreatedBy = GetUserID(c)
	results, err := h.svc.RunDepreciation(c.Request.Context(), req)
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "count": len(results)})
}

func (h *FAHandler) RunDepreciationForAsset(c *gin.Context) {
	var req domain.FARunDepreciationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	req.CreatedBy = GetUserID(c)
	result, err := h.svc.RunDepreciationForAsset(c.Request.Context(), req, c.Param("id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *FAHandler) GetDepreciationByAsset(c *gin.Context) {
	entries, err := h.svc.GetDepreciationByAsset(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *FAHandler) GetDepreciationByPeriod(c *gin.Context) {
	entries, err := h.svc.GetDepreciationByPeriod(c.Request.Context(), c.Param("periodID"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

// ─── Transactions ────────────────────────────────────────────────────

func (h *FAHandler) GetTransactions(c *gin.Context) {
	txns, err := h.svc.GetTransactions(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, txns)
}

// ─── Allocations ─────────────────────────────────────────────────────

func (h *FAHandler) GetAllocations(c *gin.Context) {
	allocs, err := h.svc.GetAllocations(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, allocs)
}

func (h *FAHandler) SetAllocations(c *gin.Context) {
	var req struct {
		Allocations []domain.FixedAssetAllocation `json:"allocations"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.svc.SetAllocations(c.Request.Context(), c.Param("id"), req.Allocations); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "allocations updated"})
}

// ─── Business Operations ─────────────────────────────────────────────

func (h *FAHandler) DisposeAsset(c *gin.Context) {
	var req domain.FADisposalInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.DisposeAsset(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "asset disposed"})
}

func (h *FAHandler) SellAsset(c *gin.Context) {
	var req domain.FASaleInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.SellAsset(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "asset sold"})
}

func (h *FAHandler) AdjustAsset(c *gin.Context) {
	var req domain.FAAdjustmentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.AdjustAsset(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "asset adjusted"})
}

func (h *FAHandler) RevalueAsset(c *gin.Context) {
	var req domain.FARevaluationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.RevalueAsset(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "asset revalued"})
}

func (h *FAHandler) ImpairAsset(c *gin.Context) {
	var req domain.FAImpairmentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.ImpairAsset(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "asset impaired"})
}

func (h *FAHandler) TransferAsset(c *gin.Context) {
	var req domain.FATransferInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.TransferAsset(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "asset transferred"})
}

func (h *FAHandler) CIPTransfer(c *gin.Context) {
	var req domain.FACIPTransferInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.CIPTransfer(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "CIP transfer completed"})
}

func (h *FAHandler) SuspendDepreciation(c *gin.Context) {
	var req domain.FASuspendInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.SuspendDepreciation(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "depreciation suspended"})
}

func (h *FAHandler) ResumeDepreciation(c *gin.Context) {
	var req domain.FAResumeInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.FixedAssetID = c.Param("id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.ResumeDepreciation(c.Request.Context(), req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "depreciation resumed"})
}

// ─── Inventory ───────────────────────────────────────────────────────

func (h *FAHandler) CreateInventoryPlan(c *gin.Context) {
	var req domain.FixedAssetInventoryPlan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	req.CreatedBy = GetUserID(c)
	if err := h.svc.CreateInventoryPlan(c.Request.Context(), &req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *FAHandler) GetInventoryPlan(c *gin.Context) {
	p, err := h.svc.GetInventoryPlan(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *FAHandler) ListInventoryPlans(c *gin.Context) {
	plans, err := h.svc.ListInventoryPlans(c.Request.Context(), c.Query("company_id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *FAHandler) UpdateInventoryPlan(c *gin.Context) {
	var req domain.FixedAssetInventoryPlan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	if err := h.svc.UpdateInventoryPlan(c.Request.Context(), &req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *FAHandler) CreateInventoryResult(c *gin.Context) {
	var req domain.FixedAssetInventoryResult
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.svc.CreateInventoryResult(c.Request.Context(), &req); err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *FAHandler) GetInventoryResults(c *gin.Context) {
	results, err := h.svc.GetInventoryResults(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.faError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}
