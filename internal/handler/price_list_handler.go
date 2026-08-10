package handler

import (
	"net/http"

	"gotax/internal/domain"
	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

type PriceListHandler struct {
	svc *service.PriceListService
}

func NewPriceListHandler(svc *service.PriceListService) *PriceListHandler {
	return &PriceListHandler{svc: svc}
}

func RegisterPriceListRoutes(r *gin.Engine, h *PriceListHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	pl := v1.Group("/price-lists")
	{
		pl.POST("", h.CreatePriceList)
		pl.GET("", h.ListPriceLists)
		pl.GET("/:id", h.GetPriceList)
		pl.PUT("/:id", h.UpdatePriceList)
		pl.DELETE("/:id", h.DeletePriceList)
		pl.POST("/:id/lines", h.AddLines)
		pl.GET("/:id/lines", h.GetLines)
		pl.POST("/calculate-price", h.CalculateSellingPrice)
	}
}

func (h *PriceListHandler) CreatePriceList(c *gin.Context) {
	var pl domain.PriceList
	if err := c.ShouldBindJSON(&pl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pl.CompanyID = c.Query("company_id")
	if err := h.svc.CreatePriceList(c.Request.Context(), &pl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pl)
}

func (h *PriceListHandler) GetPriceList(c *gin.Context) {
	pl, err := h.svc.GetPriceList(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pl)
}

func (h *PriceListHandler) ListPriceLists(c *gin.Context) {
	companyID := c.Query("company_id")
	lists, err := h.svc.ListPriceLists(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

func (h *PriceListHandler) UpdatePriceList(c *gin.Context) {
	var pl domain.PriceList
	if err := c.ShouldBindJSON(&pl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pl.ID = c.Param("id")
	if err := h.svc.UpdatePriceList(c.Request.Context(), &pl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pl)
}

func (h *PriceListHandler) DeletePriceList(c *gin.Context) {
	if err := h.svc.DeletePriceList(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *PriceListHandler) AddLines(c *gin.Context) {
	var lines []domain.PriceListLine
	if err := c.ShouldBindJSON(&lines); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddLines(c.Request.Context(), c.Param("id"), lines); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"count": len(lines)})
}

func (h *PriceListHandler) GetLines(c *gin.Context) {
	lines, err := h.svc.GetLines(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lines)
}

func (h *PriceListHandler) CalculateSellingPrice(c *gin.Context) {
	var req struct {
		PriceListID string  `json:"price_list_id"`
		ItemCode    string  `json:"item_code"`
		MarkupPct   float64 `json:"markup_pct"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	price, err := h.svc.CalculateSellingPrice(c.Request.Context(), req.PriceListID, req.ItemCode, req.MarkupPct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"unit_price": price})
}
