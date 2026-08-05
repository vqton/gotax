package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
)

func (h *Handler) CreateExchangeRate(c *gin.Context) {
	var rate domain.ExchangeRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateExchangeRate(c.Request.Context(), &rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rate)
}

func (h *Handler) ListExchangeRates(c *gin.Context) {
	rates, err := h.svc.ListExchangeRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rates)
}
