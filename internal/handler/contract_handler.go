package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type ContractHandler struct {
	svc *service.ContractService
}

func NewContractHandler(svc *service.ContractService) *ContractHandler {
	return &ContractHandler{svc: svc}
}

func RegisterContractRoutes(r *gin.Engine, h *ContractHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	cons := v1.Group("/contracts")
	{
		cons.POST("", h.Create)
		cons.GET("", h.List)
		cons.GET("/:id", h.Get)
		cons.PUT("/:id", h.Update)
		cons.DELETE("/:id", h.Delete)
		cons.GET("/:id/payments", h.ListPayments)
		cons.POST("/:id/payments", h.AddPayment)
	}
}

func (h *ContractHandler) contractError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrContractNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrContractExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *ContractHandler) Create(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	var req domain.Contract
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = companyID
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		h.contractError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *ContractHandler) Get(c *gin.Context) {
	contract, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.contractError(c, err)
		return
	}
	c.JSON(http.StatusOK, contract)
}

func (h *ContractHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	list, err := h.svc.List(c.Request.Context(), companyID)
	if err != nil {
		h.contractError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ContractHandler) Update(c *gin.Context) {
	var req domain.Contract
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		h.contractError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *ContractHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.contractError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ContractHandler) AddPayment(c *gin.Context) {
	var req domain.ContractPayment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ContractID = c.Param("id")
	if err := h.svc.AddPayment(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *ContractHandler) ListPayments(c *gin.Context) {
	list, err := h.svc.ListPayments(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
