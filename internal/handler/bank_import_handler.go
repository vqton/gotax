package handler

import (
	"io"
	"net/http"

	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

type BankImportHandler struct {
	svc *service.BankImportService
}

func NewBankImportHandler(svc *service.BankImportService) *BankImportHandler {
	return &BankImportHandler{svc: svc}
}

func RegisterBankImportRoutes(r *gin.Engine, h *BankImportHandler, authMW gin.HandlerFunc) {
	bank := r.Group("/api/v1/bank", authMW)
	{
		imports := bank.Group("/imports")
		{
			imports.POST("", h.ImportCSV)
			imports.GET("", h.ListImports)
			imports.GET("/:id", h.GetImport)
			imports.GET("/:id/transactions", h.GetTransactions)
		}
	}
}

func (h *BankImportHandler) ImportCSV(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}
	bankCode := c.PostForm("bank_code")
	if bankCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bank_code required"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	imp, err := h.svc.Import(c.Request.Context(), companyID, bankCode, file.Filename, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, imp)
}

func (h *BankImportHandler) GetImport(c *gin.Context) {
	id := c.Param("id")
	imp, err := h.svc.GetImport(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, imp)
}

func (h *BankImportHandler) ListImports(c *gin.Context) {
	companyID := c.Query("company_id")
	items, err := h.svc.ListImports(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *BankImportHandler) GetTransactions(c *gin.Context) {
	id := c.Param("id")
	txns, err := h.svc.GetTransactions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": txns})
}
