package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type BudgetHandler struct {
	svc *service.BudgetService
}

func NewBudgetHandler(svc *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{svc: svc}
}

func RegisterBudgetRoutes(r *gin.Engine, h *BudgetHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	budgets := v1.Group("/budgets")
	{
		budgets.POST("", h.Create)
		budgets.GET("", h.List)
		budgets.GET("/report", h.VarianceReport)
		budgets.POST("/upsert", h.BulkUpsert)
		budgets.GET("/:id", h.Get)
		budgets.PUT("/:id", h.Update)
		budgets.DELETE("/:id", h.Delete)
		budgets.POST("/sync-actuals", h.SyncActuals)
	}
}

func (h *BudgetHandler) budgetError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrBudgetNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrBudgetExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *BudgetHandler) Create(c *gin.Context) {
	var req domain.Budget
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		h.budgetError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *BudgetHandler) Get(c *gin.Context) {
	b, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.budgetError(c, err)
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BudgetHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	year := 0
	if y := c.Query("year"); y != "" {
		year, _ = strconv.Atoi(y)
	}
	list, err := h.svc.List(c.Request.Context(), companyID, year)
	if err != nil {
		h.budgetError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *BudgetHandler) Update(c *gin.Context) {
	var req domain.Budget
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		h.budgetError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *BudgetHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.budgetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *BudgetHandler) BulkUpsert(c *gin.Context) {
	var req []*domain.Budget
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	companyID := c.Query("company_id")
	for _, b := range req {
		if b.CompanyID == "" {
			b.CompanyID = companyID
		}
	}
	count, err := h.svc.BulkUpsert(c.Request.Context(), req)
	if err != nil {
		h.budgetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"upserted": count})
}

func (h *BudgetHandler) VarianceReport(c *gin.Context) {
	companyID := c.Query("company_id")
	year := 0
	if y := c.Query("year"); y != "" {
		year, _ = strconv.Atoi(y)
	}
	month := 0
	if m := c.Query("month"); m != "" {
		month, _ = strconv.Atoi(m)
	}
	report, err := h.svc.VarianceReport(c.Request.Context(), companyID, year, month)
	if err != nil {
		h.budgetError(c, err)
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *BudgetHandler) SyncActuals(c *gin.Context) {
	companyID := c.Query("company_id")
	year := 0
	if y := c.Query("year"); y != "" {
		year, _ = strconv.Atoi(y)
	}
	month := 0
	if m := c.Query("month"); m != "" {
		month, _ = strconv.Atoi(m)
	}
	count, err := h.svc.SyncActuals(c.Request.Context(), companyID, year, month)
	if err != nil {
		h.budgetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": count})
}
