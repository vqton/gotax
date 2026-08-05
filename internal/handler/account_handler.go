package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
)

func (h *Handler) CreateAccount(c *gin.Context) {
	var account domain.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAccount(c.Request.Context(), &account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}

func (h *Handler) GetAccount(c *gin.Context) {
	code := c.Param("code")
	account, err := h.svc.GetAccount(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (h *Handler) ListAccounts(c *gin.Context) {
	activeOnly := c.DefaultQuery("active_only", "false") == "true"
	accounts, err := h.svc.GetAllAccounts(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	code := c.Param("code")
	var account domain.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	account.Code = code
	if err := h.svc.UpdateAccount(c.Request.Context(), &account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	code := c.Param("code")
	if err := h.svc.DeleteAccount(c.Request.Context(), code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}

func (h *Handler) FreezeAccount(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}
	if err := h.svc.FreezeAccount(c.Request.Context(), code, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account frozen"})
}

func (h *Handler) UnfreezeAccount(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}
	if err := h.svc.UnfreezeAccount(c.Request.Context(), code, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account unfrozen"})
}

func (h *Handler) GetAccountBalance(c *gin.Context) {
	code := c.Param("code")
	periodID := c.Query("period_id")
	if periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_id is required"})
		return
	}
	balance, err := h.svc.GetAccountBalance(c.Request.Context(), code, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *Handler) GetAccountBalanceDrillDown(c *gin.Context) {
	code := c.Param("code")
	periodID := c.Query("period_id")
	if periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_id is required"})
		return
	}
	entries, err := h.svc.GetAccountBalanceDrillDown(c.Request.Context(), code, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetAccountUsage(c *gin.Context) {
	code := c.Param("code")
	usage, err := h.svc.GetAccountUsage(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usage)
}
