package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gotax/internal/service"
)

type FiscalYearHandler struct {
	svc *service.FiscalYearService
}

func NewFiscalYearHandler(svc *service.FiscalYearService) *FiscalYearHandler {
	return &FiscalYearHandler{svc: svc}
}

func RegisterFiscalYearRoutes(r *gin.Engine, h *FiscalYearHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	fy := v1.Group("/fiscal-years")
	{
		fy.POST("/:year/periods", h.CreateYear)
		fy.GET("/periods", h.ListPeriods)
	}
}

func (h *FiscalYearHandler) CreateYear(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	periods, err := h.svc.CreateYearFromOptions(c.Request.Context(), companyID, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"periods": periods, "count": len(periods)})
}

func (h *FiscalYearHandler) ListPeriods(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	periods, err := h.svc.ListPeriods(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, periods)
}
