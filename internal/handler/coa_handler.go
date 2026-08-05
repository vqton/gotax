package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
)

func (h *Handler) CreateApprovalRequest(c *gin.Context) {
	var req domain.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	req.RequestedBy = GetUserID(c)
	if err := h.svc.CreateApprovalRequest(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) ApproveRequest(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note string `json:"note"`
	}
	c.ShouldBindJSON(&req)
	userID := GetUserID(c)
	if err := h.svc.ApproveRequest(c.Request.Context(), id, userID, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request approved"})
}

func (h *Handler) RejectRequest(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "note is required"})
		return
	}
	userID := GetUserID(c)
	if err := h.svc.RejectRequest(c.Request.Context(), id, userID, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request rejected"})
}

func (h *Handler) ListApprovalRequests(c *gin.Context) {
	status := domain.ApprovalStatus(c.Query("status"))
	requests, err := h.svc.GetApprovalRequests(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

func (h *Handler) CreateAccountVersion(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}
	ver, err := h.svc.CreateAccountVersion(c.Request.Context(), req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ver)
}

func (h *Handler) GetVersion(c *gin.Context) {
	versionNumber := c.Param("versionNumber")
	ver, err := h.svc.GetVersion(c.Request.Context(), versionNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ver)
}

func (h *Handler) ListVersions(c *gin.Context) {
	versions, err := h.svc.ListVersions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func (h *Handler) CompareVersions(c *gin.Context) {
	v1 := c.Query("v1")
	v2 := c.Query("v2")
	if v1 == "" || v2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "v1 and v2 query params required"})
		return
	}
	diff, err := h.svc.CompareVersions(c.Request.Context(), v1, v2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, diff)
}

func (h *Handler) CreateAccountAnalysis(c *gin.Context) {
	code := c.Param("code")
	var analysis domain.AccountAnalysis
	if err := c.ShouldBindJSON(&analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	analysis.AccountCode = code
	if err := h.svc.CreateAccountAnalysis(c.Request.Context(), &analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, analysis)
}

func (h *Handler) GetAccountAnalysis(c *gin.Context) {
	code := c.Param("code")
	analysis, err := h.svc.GetAccountAnalysis(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analysis)
}

func (h *Handler) UpdateAccountAnalysis(c *gin.Context) {
	code := c.Param("code")
	var analysis domain.AccountAnalysis
	if err := c.ShouldBindJSON(&analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	analysis.AccountCode = code
	if err := h.svc.UpdateAccountAnalysis(c.Request.Context(), &analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analysis)
}

func (h *Handler) CreateAccountMapping(c *gin.Context) {
	var mapping domain.AccountMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAccountMapping(c.Request.Context(), &mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mapping)
}

func (h *Handler) GetMappingByOldCode(c *gin.Context) {
	oldCode := c.Param("oldCode")
	sourceRegime := c.Query("source_regime")
	if sourceRegime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_regime query param required"})
		return
	}
	mapping, err := h.svc.GetMappingByOldCode(c.Request.Context(), sourceRegime, oldCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mapping)
}

func (h *Handler) ListMappings(c *gin.Context) {
	sourceRegime := c.Query("source_regime")
	targetRegime := c.Query("target_regime")
	mappings, err := h.svc.ListMappings(c.Request.Context(), sourceRegime, targetRegime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mappings)
}

func (h *Handler) CreateIFRSMapping(c *gin.Context) {
	var mapping domain.IFRSMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateIFRSMapping(c.Request.Context(), &mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mapping)
}

func (h *Handler) GetIFRSMapping(c *gin.Context) {
	vasCode := c.Param("vasCode")
	mapping, err := h.svc.GetIFRSMapping(c.Request.Context(), vasCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mapping)
}

func (h *Handler) ListIFRSMappings(c *gin.Context) {
	mappings, err := h.svc.ListIFRSMappings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mappings)
}
