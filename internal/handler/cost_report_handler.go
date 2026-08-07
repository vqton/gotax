package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/service"
)

type CostReportHandler struct {
	reportSvc *service.CostReportService
}

func NewCostReportHandler(reportSvc *service.CostReportService) *CostReportHandler {
	return &CostReportHandler{reportSvc: reportSvc}
}

func RegisterCostReportRoutes(r *gin.Engine, h *CostReportHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	reports := v1.Group("/cost-reports")
	{
		reports.GET("/calculation-sheet", h.GetCostCalculationSheet)
		reports.GET("/summary", h.GetCostSummary)
		reports.GET("/wip-valuation", h.GetWIPValuation)
		reports.GET("/variance-analysis", h.GetVarianceAnalysis)
	}
}

func (h *CostReportHandler) GetCostCalculationSheet(c *gin.Context) {
	companyID := c.Query("company_id")
	periodID := c.Query("period_id")
	objectID := c.Query("object_id")
	if companyID == "" || periodID == "" || objectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id, period_id, and object_id are required"})
		return
	}
	sheet, err := h.reportSvc.GetCostCalculationSheet(c.Request.Context(), companyID, periodID, objectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sheet)
}

func (h *CostReportHandler) GetCostSummary(c *gin.Context) {
	companyID := c.Query("company_id")
	periodID := c.Query("period_id")
	if companyID == "" || periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id and period_id are required"})
		return
	}
	summary, err := h.reportSvc.GetCostSummary(c.Request.Context(), companyID, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *CostReportHandler) GetWIPValuation(c *gin.Context) {
	companyID := c.Query("company_id")
	periodID := c.Query("period_id")
	if companyID == "" || periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id and period_id are required"})
		return
	}
	report, err := h.reportSvc.GetWIPValuation(c.Request.Context(), companyID, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *CostReportHandler) GetVarianceAnalysis(c *gin.Context) {
	companyID := c.Query("company_id")
	periodID := c.Query("period_id")
	if companyID == "" || periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id and period_id are required"})
		return
	}
	analysis, err := h.reportSvc.GetVarianceAnalysis(c.Request.Context(), companyID, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analysis)
}
