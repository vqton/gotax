package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"

	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupExcelTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	memPurchaseRepo := repository.NewMemoryPurchaseRepo()
	memSaleRepo := repository.NewMemorySaleRepo()
	itemRepo := repository.NewMemoryItemRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	accRepo := repository.NewMemoryAccountRepo()
	perRepo := repository.NewMemoryPeriodRepo()
	importSvc := service.NewExcelImportService(memPurchaseRepo, memSaleRepo, itemRepo)
	exportSvc := service.NewExportService(jeRepo, accRepo, perRepo)
	eh := NewExcelHandler(importSvc, exportSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterExcelRoutes(r, eh, noopMW)
	return r
}

func createTestXlsx(t *testing.T, headers []string, rows [][]string) []byte {
	t.Helper()
	f := excelize.NewFile()
	sheet := "Sheet1"
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for r, row := range rows {
		for c, val := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			f.SetCellValue(sheet, cell, val)
		}
	}
	var buf bytes.Buffer
	f.Write(&buf)
	return buf.Bytes()
}

func TestImportSuppliers_Success(t *testing.T) {
	r := setupExcelTest(t)
	data := createTestXlsx(t,
		[]string{"Code", "Name", "TaxCode", "Address", "Phone", "Email"},
		[][]string{{"SUP001", "Supplier A", "0123456789", "Hanoi", "0901234567", "a@test.com"}},
	)

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, _ := w.CreateFormFile("file", "suppliers.xlsx")
	part.Write(data)
	w.Close()

	w2 := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/excel/import/suppliers?company_id=C001", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestImportSuppliers_MissingCompanyID(t *testing.T) {
	r := setupExcelTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/excel/import/suppliers", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestImportCustomers_Success(t *testing.T) {
	r := setupExcelTest(t)
	data := createTestXlsx(t,
		[]string{"Code", "Name", "TaxCode"},
		[][]string{{"CUST001", "Customer A", "9876543210"}},
	)

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, _ := w.CreateFormFile("file", "customers.xlsx")
	part.Write(data)
	w.Close()

	w2 := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/excel/import/customers?company_id=C001", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestImportItems_Success(t *testing.T) {
	r := setupExcelTest(t)
	data := createTestXlsx(t,
		[]string{"Code", "Name", "Unit"},
		[][]string{{"ITEM001", "Widget", "PCS"}},
	)

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, _ := w.CreateFormFile("file", "items.xlsx")
	part.Write(data)
	w.Close()

	w2 := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/excel/import/items?company_id=C001", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)
}
