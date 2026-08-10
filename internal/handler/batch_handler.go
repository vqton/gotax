package handler

import (
	"net/http"

	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

type BatchHandler struct {
	svc *service.BatchOperationService
}

func NewBatchHandler(svc *service.BatchOperationService) *BatchHandler {
	return &BatchHandler{svc: svc}
}

func RegisterBatchRoutes(r *gin.Engine, h *BatchHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1/journal-entries/batch", authMW)
	v1.POST("/submit", h.BatchSubmit)
	v1.POST("/approve", h.BatchApprove)
	v1.POST("/post", h.BatchPost)
	v1.POST("/cancel", h.BatchCancel)
}

type batchRequest struct {
	IDs []string `json:"ids" binding:"required,min=1"`
}

func (h *BatchHandler) BatchSubmit(c *gin.Context) {
	var req batchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids required"})
		return
	}
	userID := GetUserID(c)
	result, err := h.svc.BatchSubmit(c.Request.Context(), req.IDs, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *BatchHandler) BatchApprove(c *gin.Context) {
	var req batchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids required"})
		return
	}
	userID := GetUserID(c)
	result, err := h.svc.BatchApprove(c.Request.Context(), req.IDs, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *BatchHandler) BatchPost(c *gin.Context) {
	var req batchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids required"})
		return
	}
	result, err := h.svc.BatchPost(c.Request.Context(), req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *BatchHandler) BatchCancel(c *gin.Context) {
	var req batchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids required"})
		return
	}
	result, err := h.svc.BatchCancel(c.Request.Context(), req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
