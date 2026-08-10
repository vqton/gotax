package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/service"
)

type BackupHandler struct {
	svc *service.BackupService
}

func NewBackupHandler(svc *service.BackupService) *BackupHandler {
	return &BackupHandler{svc: svc}
}

func RegisterBackupRoutes(r *gin.Engine, h *BackupHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	bkps := v1.Group("/backups")
	{
		bkps.POST("", h.Create)
		bkps.GET("", h.List)
		bkps.GET("/:id", h.Get)
	}
}

func (h *BackupHandler) Create(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	userID, _ := c.Get("user_id")
	b, err := h.svc.CreateBackup(c.Request.Context(), companyID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *BackupHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	list, err := h.svc.List(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *BackupHandler) Get(c *gin.Context) {
	b, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}
