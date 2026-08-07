package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type CostPoolHandler struct {
	svc *service.CostPoolService
}

func NewCostPoolHandler(svc *service.CostPoolService) *CostPoolHandler {
	return &CostPoolHandler{svc: svc}
}

func RegisterCostPoolRoutes(r *gin.Engine, h *CostPoolHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	cp := v1.Group("/cost-pools")
	{
		cp.POST("", h.Create)
		cp.GET("", h.List)
		cp.GET("/:id", h.Get)
		cp.DELETE("/:id", h.Delete)
		cp.POST("/:id/lines", h.AddLine)
		cp.GET("/:id/lines", h.ListLines)
		cp.POST("/:id/close", h.ClosePool)
	}
}

func (h *CostPoolHandler) costPoolError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCostPoolNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCostPoolAlreadyClosed):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCostPoolNoLines),
		errors.Is(err, domain.ErrCostPoolAccountRequired):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *CostPoolHandler) Create(c *gin.Context) {
	var req domain.CostPool
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		h.costPoolError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *CostPoolHandler) Get(c *gin.Context) {
	cp, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.costPoolError(c, err)
		return
	}
	c.JSON(http.StatusOK, cp)
}

func (h *CostPoolHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	periodID := c.Query("period_id")
	list, err := h.svc.ListByPeriod(c.Request.Context(), companyID, periodID)
	if err != nil {
		h.costPoolError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CostPoolHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.costPoolError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *CostPoolHandler) AddLine(c *gin.Context) {
	var req domain.CostPoolLine
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.PoolID = c.Param("id")
	if err := h.svc.AddLine(c.Request.Context(), &req); err != nil {
		h.costPoolError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *CostPoolHandler) ListLines(c *gin.Context) {
	lines, err := h.svc.ListLines(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.costPoolError(c, err)
		return
	}
	c.JSON(http.StatusOK, lines)
}

func (h *CostPoolHandler) ClosePool(c *gin.Context) {
	if err := h.svc.ClosePool(c.Request.Context(), c.Param("id")); err != nil {
		h.costPoolError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "closed"})
}
