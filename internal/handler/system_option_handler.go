package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

// ─── SystemOptionHandler ───────────────────────────────────────────

type SystemOptionHandler struct {
	svc *service.SystemOptionService
}

func NewSystemOptionHandler(svc *service.SystemOptionService) *SystemOptionHandler {
	return &SystemOptionHandler{svc: svc}
}

func RegisterSystemOptionRoutes(r *gin.Engine, h *SystemOptionHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	opts := v1.Group("/system-options")
	{
		opts.PUT("", h.Upsert)
		opts.GET("", h.GetAll)
		opts.GET("/:category", h.GetByCategory)
		opts.DELETE("/:category/:key", h.Delete)
	}
}

func (h *SystemOptionHandler) Upsert(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	var req domain.SystemOption
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = companyID
	if err := h.svc.Upsert(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *SystemOptionHandler) GetByCategory(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	category := c.Param("category")
	list, err := h.svc.GetByCategory(c.Request.Context(), companyID, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *SystemOptionHandler) GetAll(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	list, err := h.svc.GetAll(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *SystemOptionHandler) Delete(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	category := c.Param("category")
	key := c.Param("key")
	if err := h.svc.Delete(c.Request.Context(), companyID, category, key); err != nil {
		if errors.Is(err, domain.ErrSystemOptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ─── NumberingRuleHandler ──────────────────────────────────────────

type NumberingRuleHandler struct {
	svc *service.NumberingRuleService
}

func NewNumberingRuleHandler(svc *service.NumberingRuleService) *NumberingRuleHandler {
	return &NumberingRuleHandler{svc: svc}
}

func RegisterNumberingRuleRoutes(r *gin.Engine, h *NumberingRuleHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	rules := v1.Group("/numbering-rules")
	{
		rules.POST("", h.Create)
		rules.GET("", h.List)
	}
	v1.GET("/numbering-rules/:id", h.Get)
	v1.PUT("/numbering-rules/:id", h.Update)
	v1.DELETE("/numbering-rules/:id", h.Delete)
	v1.GET("/numbering-rules/next/:voucherType", h.GetNext)
}

func (h *NumberingRuleHandler) numberingError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNumberingRuleNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrNumberingRuleExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *NumberingRuleHandler) Create(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	var req domain.NumberingRule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = companyID
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		h.numberingError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *NumberingRuleHandler) Get(c *gin.Context) {
	rule, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.numberingError(c, err)
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *NumberingRuleHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	list, err := h.svc.List(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *NumberingRuleHandler) Update(c *gin.Context) {
	var req domain.NumberingRule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		h.numberingError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *NumberingRuleHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.numberingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *NumberingRuleHandler) GetNext(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	voucherType := c.Param("voucherType")
	num, err := h.svc.GetNextNumber(c.Request.Context(), companyID, voucherType)
	if err != nil {
		h.numberingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"voucher_type": voucherType, "current_num": num})
}
