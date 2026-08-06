package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func RegisterNotificationRoutes(r *gin.Engine, h *NotificationHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1/notifications", authMW)
	{
		v1.GET("", h.List)
		v1.GET("/unread-count", h.UnreadCount)
		v1.POST("/read-all", h.MarkAllRead)       // must be before /:id routes
		v1.POST("/:id/read", h.MarkRead)
		v1.DELETE("/:id", h.Delete)
	}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := GetUserID(c)
	companyID := c.Query("company_id")
	if companyID == "" {
		companyID = c.GetHeader("X-Company-ID")
	}
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	notifs, err := h.svc.List(c.Request.Context(), companyID, userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notifications": notifs, "total": len(notifs)})
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := GetUserID(c)
	companyID := c.Query("company_id")
	if companyID == "" {
		companyID = c.GetHeader("X-Company-ID")
	}
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}
	count, err := h.svc.UnreadCount(c.Request.Context(), companyID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.MarkRead(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := GetUserID(c)
	companyID := c.Query("company_id")
	if companyID == "" {
		companyID = c.GetHeader("X-Company-ID")
	}
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}
	if err := h.svc.MarkAllRead(c.Request.Context(), companyID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification deleted"})
}
