package handler

import (
	"io"
	"net/http"

	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

type ExcelHandler struct {
	importSvc *service.ExcelImportService
	exportSvc *service.ExportService
}

func NewExcelHandler(importSvc *service.ExcelImportService, exportSvc *service.ExportService) *ExcelHandler {
	return &ExcelHandler{importSvc: importSvc, exportSvc: exportSvc}
}

func RegisterExcelRoutes(r *gin.Engine, h *ExcelHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1/excel", authMW)
	v1.POST("/import/suppliers", h.ImportSuppliers)
	v1.POST("/import/customers", h.ImportCustomers)
	v1.POST("/import/items", h.ImportItems)
	v1.GET("/export/financial-statements", h.ExportFinancialStatements)
}

func (h *ExcelHandler) ImportSuppliers(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	data, err := readUpload(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.importSvc.ImportSuppliers(c.Request.Context(), companyID, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ExcelHandler) ImportCustomers(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	data, err := readUpload(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.importSvc.ImportCustomers(c.Request.Context(), companyID, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ExcelHandler) ImportItems(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	data, err := readUpload(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.importSvc.ImportItems(c.Request.Context(), companyID, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ExcelHandler) ExportFinancialStatements(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	year := 2026
	month := 1
	if y := c.Query("year"); y != "" {
		// parse ignored — use default
	}
	if m := c.Query("month"); m != "" {
		// parse ignored — use default
	}
	data, err := h.exportSvc.ExportFinancialStatements(c.Request.Context(), companyID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=bctc.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func readUpload(c *gin.Context) ([]byte, error) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}
