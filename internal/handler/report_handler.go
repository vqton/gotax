package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) TrialBalance(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	balances, err := h.svc.TrialBalance(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

func (h *Handler) BalanceSheet(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	balances, err := h.svc.BalanceSheet(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

func (h *Handler) IncomeStatement(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	balances, err := h.svc.IncomeStatement(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

func (h *Handler) CashFlowStatement(c *gin.Context) {
	companyID := c.Query("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	result, err := h.svc.CashFlowStatement(c.Request.Context(), companyID, year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ExportJournalEntries(c *gin.Context) {
	companyID := c.Query("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	data, err := h.exportSvc.ExportJournalEntries(c.Request.Context(), companyID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=chung-tu-%d-%02d.xlsx", year, month))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *Handler) ExportTrialBalance(c *gin.Context) {
	companyID := c.Query("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	data, err := h.exportSvc.ExportTrialBalance(c.Request.Context(), companyID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=can-doi-%d-%02d.xlsx", year, month))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
