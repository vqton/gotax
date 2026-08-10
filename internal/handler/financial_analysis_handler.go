package handler

import (
	"net/http"
	"strconv"

	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

type FinancialAnalysisHandler struct {
	svc *service.FinancialAnalysisService
}

func NewFinancialAnalysisHandler(svc *service.FinancialAnalysisService) *FinancialAnalysisHandler {
	return &FinancialAnalysisHandler{svc: svc}
}

func RegisterFinancialAnalysisRoutes(r *gin.Engine, h *FinancialAnalysisHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1/reports", authMW)
	v1.GET("/financial-ratios", h.GetFinancialRatios)
	v1.GET("/budget-vs-actual", h.GetBudgetVsActual)
}

func (h *FinancialAnalysisHandler) GetFinancialRatios(c *gin.Context) {
	companyID := c.Query("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(2026)))
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	ratios, err := h.svc.CalculateRatios(c.Request.Context(), companyID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ratios)
}

func (h *FinancialAnalysisHandler) GetBudgetVsActual(c *gin.Context) {
	companyID := c.Query("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(2026)))
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	result, err := h.svc.CompareBudgetVsActual(c.Request.Context(), companyID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
