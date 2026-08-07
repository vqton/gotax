package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type CostingHandler struct {
	periodSvc    *service.CostingPeriodService
	engine       *service.CostingEngine
	costingJESvc *service.CostingJEService
}

func NewCostingHandler(periodSvc *service.CostingPeriodService, engine *service.CostingEngine, costingJESvc *service.CostingJEService) *CostingHandler {
	return &CostingHandler{periodSvc: periodSvc, engine: engine, costingJESvc: costingJESvc}
}

func RegisterCostingRoutes(r *gin.Engine, h *CostingHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	periods := v1.Group("/costing-periods")
	{
		periods.POST("", h.CreatePeriod)
		periods.GET("", h.ListPeriods)
		periods.GET("/:id", h.GetPeriod)
		periods.POST("/:id/close", h.ClosePeriod)
		periods.POST("/:id/close-je", h.ClosePeriodWithJE)
		periods.POST("/:id/reopen", h.ReopenPeriod)
		periods.POST("/:id/collect-materials", h.CollectMaterialCosts)
		periods.POST("/:id/collect-labor", h.CollectLaborCosts)
		periods.POST("/:id/cogs", h.GenerateCOGS)
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

func (h *CostingHandler) ClosePeriodWithJE(c *gin.Context) {
	id := c.Param("id")
	companyID := c.Query("company_id")
	userID := c.GetString("user_id")
	if err := h.costingJESvc.ClosePeriod(c.Request.Context(), companyID, id, userID); err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period closed with journal entries"})
}

func (h *CostingHandler) ReopenPeriod(c *gin.Context) {
	id := c.Param("id")
	companyID := c.Query("company_id")
	if err := h.costingJESvc.ReopenPeriod(c.Request.Context(), companyID, id); err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period reopened"})
}

type collectCostsRequest struct {
	Lines []domain.CostPoolLineInput `json:"lines"`
}

func (h *CostingHandler) CollectMaterialCosts(c *gin.Context) {
	id := c.Param("id")
	companyID := c.Query("company_id")
	var req collectCostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.costingJESvc.CollectMaterialCosts(c.Request.Context(), companyID, id, req.Lines); err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "material costs collected"})
}

func (h *CostingHandler) CollectLaborCosts(c *gin.Context) {
	id := c.Param("id")
	companyID := c.Query("company_id")
	var req collectCostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.costingJESvc.CollectLaborCosts(c.Request.Context(), companyID, id, req.Lines); err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "labor costs collected"})
}

type generateCOGSRequest struct {
	ObjectID string  `json:"object_id"`
	Quantity float64 `json:"quantity"`
}

func (h *CostingHandler) GenerateCOGS(c *gin.Context) {
	id := c.Param("id")
	companyID := c.Query("company_id")
	var req generateCOGSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.costingJESvc.GenerateCOGSEntry(c.Request.Context(), companyID, id, req.ObjectID, req.Quantity); err != nil {
		h.costingError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "COGS entry generated"})
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
