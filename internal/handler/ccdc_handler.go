package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type CCDCHandler struct {
	svc service.CCDCServiceInterface
}

func NewCCDCHandler(svc service.CCDCServiceInterface) *CCDCHandler {
	return &CCDCHandler{svc: svc}
}

func RegisterCCDCRoutes(r *gin.Engine, h *CCDCHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	ccdc := v1.Group("/ccdc")
	{
		cats := ccdc.Group("/categories")
		{
			cats.POST("", h.CreateCategory)
			cats.GET("", h.ListCategories)
			cats.GET("/:id", h.GetCategory)
			cats.PUT("/:id", h.UpdateCategory)
			cats.DELETE("/:id", h.DeleteCategory)
		}
		items := ccdc.Group("")
		{
			items.POST("", h.Create)
			items.GET("", h.List)
			items.GET("/:id", h.GetByID)
			items.PUT("/:id", h.Update)
			items.DELETE("/:id", h.Delete)
		}
	}
}

func (h *CCDCHandler) ccdcError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCCDCNotFound),
		errors.Is(err, domain.ErrCCDCCategoryNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCCDCCodeExists),
		errors.Is(err, domain.ErrCCDCCategoryCodeExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCCDCCategoryCodeRequired),
		errors.Is(err, domain.ErrCCDCCategoryNameRequired),
		errors.Is(err, domain.ErrCCDCItemCodeRequired),
		errors.Is(err, domain.ErrCCDCItemNameRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// ─── Categories ──────────────────────────────────────────────────────

func (h *CCDCHandler) CreateCategory(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	var cat domain.ToolEquipmentCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.CompanyID = companyID
	if err := h.svc.CreateCategory(c.Request.Context(), &cat); err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *CCDCHandler) ListCategories(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	cats, err := h.svc.ListCategories(c.Request.Context(), companyID)
	if err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *CCDCHandler) GetCategory(c *gin.Context) {
	cat, err := h.svc.GetCategory(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *CCDCHandler) UpdateCategory(c *gin.Context) {
	var cat domain.ToolEquipmentCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.ID = c.Param("id")
	if err := h.svc.UpdateCategory(c.Request.Context(), &cat); err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *CCDCHandler) DeleteCategory(c *gin.Context) {
	if err := h.svc.DeleteCategory(c.Request.Context(), c.Param("id")); err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ─── Items ───────────────────────────────────────────────────────────

func (h *CCDCHandler) Create(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	var item domain.ToolEquipment
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.CompanyID = companyID
	if err := h.svc.Create(c.Request.Context(), &item); err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *CCDCHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	items, err := h.svc.List(c.Request.Context(), companyID)
	if err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *CCDCHandler) GetByID(c *gin.Context) {
	item, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *CCDCHandler) Update(c *gin.Context) {
	var item domain.ToolEquipment
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.ID = c.Param("id")
	if err := h.svc.Update(c.Request.Context(), &item); err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *CCDCHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.ccdcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
