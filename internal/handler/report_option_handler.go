package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/service"
)

type ReportOptionHandler struct {
	svc *service.ReportOptionService
}

func NewReportOptionHandler(svc *service.ReportOptionService) *ReportOptionHandler {
	return &ReportOptionHandler{svc: svc}
}

func RegisterReportOptionRoutes(r *gin.Engine, h *ReportOptionHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	rpt := v1.Group("/report-options")
	{
		rpt.GET("", h.Get)
		rpt.PUT("", h.Update)
	}
}

func (h *ReportOptionHandler) Get(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	opts, err := h.svc.Get(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, opts)
}

func (h *ReportOptionHandler) Update(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	var req service.ReportOptions
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.svc.Update(c.Request.Context(), companyID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
