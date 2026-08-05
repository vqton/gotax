package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
)

func (h *Handler) CreatePeriod(c *gin.Context) {
	var period domain.Period
	if err := c.ShouldBindJSON(&period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreatePeriod(c.Request.Context(), &period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, period)
}

func (h *Handler) GetPeriod(c *gin.Context) {
	id := c.Param("id")
	period, err := h.svc.GetPeriod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, period)
}

func (h *Handler) ListPeriods(c *gin.Context) {
	periods, err := h.svc.GetAllPeriods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, periods)
}

func (h *Handler) ClosePeriod(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ClosePeriod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period closed"})
}

func (h *Handler) ReopenPeriod(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ReopenPeriod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period reopened"})
}
