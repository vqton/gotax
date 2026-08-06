package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type CostCenterHandler struct {
	svc *service.CostCenterService
}

func NewCostCenterHandler(svc *service.CostCenterService) *CostCenterHandler {
	return &CostCenterHandler{svc: svc}
}

func RegisterCostCenterRoutes(r *gin.Engine, h *CostCenterHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	ccs := v1.Group("/cost-centers")
	{
		ccs.POST("", h.Create)
		ccs.GET("", h.List)
		ccs.GET("/hierarchy", h.ListHierarchy)
		ccs.GET("/:id", h.Get)
		ccs.PUT("/:id", h.Update)
		ccs.DELETE("/:id", h.Delete)
	}
}

func (h *CostCenterHandler) costCenterError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCostCenterNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCostCenterCodeExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *CostCenterHandler) Create(c *gin.Context) {
	var req domain.CostCenter
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		h.costCenterError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *CostCenterHandler) Get(c *gin.Context) {
	cc, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.costCenterError(c, err)
		return
	}
	c.JSON(http.StatusOK, cc)
}

func (h *CostCenterHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	list, err := h.svc.List(c.Request.Context(), companyID)
	if err != nil {
		h.costCenterError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CostCenterHandler) ListHierarchy(c *gin.Context) {
	companyID := c.Query("company_id")
	list, err := h.svc.ListHierarchy(c.Request.Context(), companyID)
	if err != nil {
		h.costCenterError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CostCenterHandler) Update(c *gin.Context) {
	var req domain.CostCenter
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		h.costCenterError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *CostCenterHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.costCenterError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
