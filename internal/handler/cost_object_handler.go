package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type CostObjectHandler struct {
	svc *service.CostObjectService
}

func NewCostObjectHandler(svc *service.CostObjectService) *CostObjectHandler {
	return &CostObjectHandler{svc: svc}
}

func RegisterCostObjectRoutes(r *gin.Engine, h *CostObjectHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	co := v1.Group("/cost-objects")
	{
		co.POST("", h.Create)
		co.GET("", h.List)
		co.GET("/:id", h.Get)
		co.PUT("/:id", h.Update)
		co.DELETE("/:id", h.Delete)
	}
}

func (h *CostObjectHandler) costObjectError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCostObjectNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCostObjectCodeExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCostObjectCodeRequired),
		errors.Is(err, domain.ErrCostObjectNameRequired),
		errors.Is(err, domain.ErrCostObjectTypeInvalid),
		errors.Is(err, domain.ErrCostObjectMethodInvalid):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *CostObjectHandler) Create(c *gin.Context) {
	var req domain.CostObject
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		h.costObjectError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *CostObjectHandler) Get(c *gin.Context) {
	co, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.costObjectError(c, err)
		return
	}
	c.JSON(http.StatusOK, co)
}

func (h *CostObjectHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	list, err := h.svc.List(c.Request.Context(), companyID)
	if err != nil {
		h.costObjectError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CostObjectHandler) Update(c *gin.Context) {
	var req domain.CostObject
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		h.costObjectError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *CostObjectHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.costObjectError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
