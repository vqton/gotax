package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type CostingHandler struct {
	periodSvc *service.CostingPeriodService
	engine    *service.CostingEngine
}

func NewCostingHandler(periodSvc *service.CostingPeriodService, engine *service.CostingEngine) *CostingHandler {
	return &CostingHandler{periodSvc: periodSvc, engine: engine}
}

func RegisterCostingRoutes(r *gin.Engine, h *CostingHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	periods := v1.Group("/costing-periods")
	{
		periods.POST("", h.CreatePeriod)
		periods.GET("", h.ListPeriods)
		periods.GET("/:id", h.GetPeriod)
		periods.POST("/:id/close", h.ClosePeriod)
	}

	v1.POST("/costing/run", h.RunCosting)
}

func (h *CostingHandler) costingError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCostingPeriodNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCostingPeriodExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCostingPeriodAlreadyClosed):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrCostPoolNoLines):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *CostingHandler) CreatePeriod(c *gin.Context) {
	var req domain.CostingPeriod
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	if err := h.periodSvc.Create(c.Request.Context(), &req); err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *CostingHandler) ListPeriods(c *gin.Context) {
	companyID := c.Query("company_id")
	periods, err := h.periodSvc.List(c.Request.Context(), companyID)
	if err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, periods)
}

func (h *CostingHandler) GetPeriod(c *gin.Context) {
	id := c.Param("id")
	period, err := h.periodSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, period)
}

func (h *CostingHandler) ClosePeriod(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	if err := h.periodSvc.Close(c.Request.Context(), id, userID); err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period closed"})
}

type runCostingRequest struct {
	CompanyID string `json:"company_id"`
	PeriodID  string `json:"period_id"`
}

func (h *CostingHandler) RunCosting(c *gin.Context) {
	var req runCostingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	if req.PeriodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_id is required"})
		return
	}

	if err := h.engine.RunCosting(c.Request.Context(), req.CompanyID, req.PeriodID); err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "costing completed"})
}
