package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type RecurringHandler struct {
	svc *service.RecurringService
}

func NewRecurringHandler(svc *service.RecurringService) *RecurringHandler {
	return &RecurringHandler{svc: svc}
}

func RegisterRecurringRoutes(r *gin.Engine, h *RecurringHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	rec := v1.Group("/recurring-entries")
	{
		rec.POST("", h.Create)
		rec.GET("", h.List)
		rec.GET("/:id", h.Get)
		rec.PUT("/:id", h.Update)
		rec.DELETE("/:id", h.Delete)
		rec.POST("/:id/run", h.RunNow)
		rec.POST("/process-due", h.ProcessDue)
	}
}

func (h *RecurringHandler) recurringError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrRecurringNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrRecurringNotActive):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrRecurringCompanyRequired),
		errors.Is(err, domain.ErrRecurringTemplateNameRequired),
		errors.Is(err, domain.ErrRecurringFrequencyInvalid),
		errors.Is(err, domain.ErrRecurringDayOfMonthInvalid),
		errors.Is(err, domain.ErrRecurringLinesRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *RecurringHandler) Create(c *gin.Context) {
	var req domain.RecurringEntry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		h.recurringError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *RecurringHandler) Get(c *gin.Context) {
	entry, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.recurringError(c, err)
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *RecurringHandler) List(c *gin.Context) {
	entries, err := h.svc.List(c.Request.Context(), c.Query("company_id"))
	if err != nil {
		h.recurringError(c, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *RecurringHandler) Update(c *gin.Context) {
	var req domain.RecurringEntry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = c.Param("id")
	req.CompanyID = c.Query("company_id")
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		h.recurringError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *RecurringHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.recurringError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *RecurringHandler) RunNow(c *gin.Context) {
	userID, _ := c.Get("user_id")
	je, err := h.svc.RunNow(c.Request.Context(), c.Param("id"), userID.(string))
	if err != nil {
		h.recurringError(c, err)
		return
	}
	c.JSON(http.StatusCreated, je)
}

func (h *RecurringHandler) ProcessDue(c *gin.Context) {
	userID, _ := c.Get("user_id")
	count, err := h.svc.ProcessDue(c.Request.Context(), userID.(string))
	if err != nil {
		h.recurringError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"processed": count})
}
