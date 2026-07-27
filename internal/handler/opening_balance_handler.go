package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
)

// ─── CRUD ──────────────────────────────────────────────────────────

func (h *Handler) CreateOpeningBalance(c *gin.Context) {
	var ob domain.OpeningBalance
	if err := c.ShouldBindJSON(&ob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ob.CreatedBy = c.GetString("user_id")
	if err := h.svc.CreateOpeningBalance(c.Request.Context(), &ob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ob)
}

func (h *Handler) GetOpeningBalance(c *gin.Context) {
	ob, err := h.svc.GetOpeningBalance(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ob)
}

func (h *Handler) ListOpeningBalances(c *gin.Context) {
	filter := domain.OBListFilter{
		CompanyID:   c.Query("company_id"),
		PeriodID:    c.Query("period_id"),
		Status:      domain.OpeningBalanceStatus(c.Query("status")),
		AccountCode: c.Query("account_code"),
		Currency:    c.Query("currency"),
		SourceType:  c.Query("source_type"),
	}
	list, err := h.svc.ListOpeningBalances(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) UpdateOpeningBalance(c *gin.Context) {
	var ob domain.OpeningBalance
	if err := c.ShouldBindJSON(&ob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ob.ID = c.Param("id")
	if err := h.svc.UpdateOpeningBalance(c.Request.Context(), &ob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ob)
}

func (h *Handler) DeleteOpeningBalance(c *gin.Context) {
	if err := h.svc.DeleteOpeningBalance(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ─── Status lifecycle ──────────────────────────────────────────────

func (h *Handler) SubmitOpeningBalance(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := h.svc.SubmitOpeningBalance(c.Request.Context(), c.Param("id"), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submitted"})
}

func (h *Handler) ApproveOpeningBalance(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := h.svc.ApproveOpeningBalance(c.Request.Context(), c.Param("id"), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *Handler) CorrectOpeningBalance(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	corrected, err := h.svc.CorrectOpeningBalance(c.Request.Context(), c.Param("id"), userID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, corrected)
}

// ─── Details ───────────────────────────────────────────────────────

func (h *Handler) CreateOpeningBalanceDetail(c *gin.Context) {
	var d domain.OpeningBalanceDetail
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d.OpeningBalanceID = c.Param("id")
	if err := h.svc.CreateOpeningBalanceDetail(c.Request.Context(), &d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) GetOpeningBalanceDetails(c *gin.Context) {
	details, err := h.svc.GetOpeningBalanceDetails(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, details)
}

func (h *Handler) DeleteOpeningBalanceDetail(c *gin.Context) {
	if err := h.svc.DeleteOpeningBalanceDetail(c.Request.Context(), c.Param("detailId")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ─── Totals & balance ──────────────────────────────────────────────

func (h *Handler) GetOpeningBalanceTotals(c *gin.Context) {
	debit, credit, err := h.svc.GetOpeningBalanceTotals(c.Request.Context(), c.Query("company_id"), c.Query("period_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_debit": debit, "total_credit": credit})
}

// ─── Carry Forward ─────────────────────────────────────────────────

func (h *Handler) CarryForward(c *gin.Context) {
	var req struct {
		CompanyID      string `json:"company_id"`
		FromPeriodID   string `json:"from_period_id"`
		ToPeriodID     string `json:"to_period_id"`
		FromFiscalYear string `json:"from_fiscal_year"`
		ToFiscalYear   string `json:"to_fiscal_year"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	log, err := h.svc.CarryForward(c.Request.Context(), req.CompanyID, req.FromPeriodID, req.ToPeriodID, req.FromFiscalYear, req.ToFiscalYear, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, log)
}

func (h *Handler) GetCarryForwardLogs(c *gin.Context) {
	logs, err := h.svc.GetCarryForwardLogs(c.Request.Context(), c.Query("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (h *Handler) GetCarryForwardLogByID(c *gin.Context) {
	log, err := h.svc.GetCarryForwardLogByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, log)
}

// ─── Circular 99 Mapping ──────────────────────────────────────────

func (h *Handler) CreateCircular99Mapping(c *gin.Context) {
	var m domain.Circular99Mapping
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateCircular99Mapping(c.Request.Context(), &m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *Handler) ListCircular99Mappings(c *gin.Context) {
	list, err := h.svc.ListCircular99Mappings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetCircular99MappingByOldCode(c *gin.Context) {
	m, err := h.svc.GetCircular99MappingByOldCode(c.Request.Context(), c.Param("oldCode"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// ─── Balance Migration ─────────────────────────────────────────────

func (h *Handler) CreateBalanceMigration(c *gin.Context) {
	var m domain.BalanceMigration
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateBalanceMigration(c.Request.Context(), &m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *Handler) GetBalanceMigrationByID(c *gin.Context) {
	m, err := h.svc.GetBalanceMigrationByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *Handler) ListBalanceMigrations(c *gin.Context) {
	list, err := h.svc.ListBalanceMigrations(c.Request.Context(), c.Query("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ImportOpeningBalances(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read file: " + err.Error()})
		return
	}

	companyID := c.PostForm("company_id")
	periodID := c.PostForm("period_id")
	userID := c.GetString("user_id")

	result, err := h.svc.ImportOpeningBalances(c.Request.Context(), data, companyID, periodID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) DownloadOpeningBalancePDF(c *gin.Context) {
	companyID := c.Query("company_id")
	periodID := c.Query("period_id")
	if companyID == "" || periodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id and period_id required"})
		return
	}
	data, err := h.svc.GenerateOpeningBalancePDF(c.Request.Context(), companyID, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="opening_balance_%s_%s.pdf"`, companyID, periodID))
	c.Data(http.StatusOK, "application/pdf", data)
}
