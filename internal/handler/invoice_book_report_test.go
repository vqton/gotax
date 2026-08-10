package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

func setupInvoiceBookReportTestRouter() (*gin.Engine, *InvoiceBookHandler) {
	gin.SetMode(gin.TestMode)
	repo := repository.NewMemoryInvoiceBookRepo()
	svc := service.NewInvoiceBookService(repo)
	h := NewInvoiceBookHandler(svc)
	r := gin.New()
	return r, h
}

func TestGetBookSummary(t *testing.T) {
	r, h := setupInvoiceBookReportTestRouter()
	r.POST("/api/v1/invoice-books", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.CreateBook)
	r.POST("/api/v1/invoice-books/:id/allocate", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.AllocateNumber)
	r.POST("/api/v1/invoice-numbers/:id/missing", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.MarkMissing)
	r.GET("/api/v1/invoice-books/:id/summary", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.GetBookSummary)

	// Create book
	book := domain.InvoiceBook{Name: "Summary Test", Pattern: "01GTTT", FromNumber: 1, ToNumber: 10}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books?company_id=CMP001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var created domain.InvoiceBook
	json.Unmarshal(w.Body.Bytes(), &created)

	// Allocate 3 numbers
	var allocatedIDs []string
	for i := 0; i < 3; i++ {
		req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books/"+created.ID+"/allocate", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var num domain.InvoiceNumber
		json.Unmarshal(w.Body.Bytes(), &num)
		allocatedIDs = append(allocatedIDs, num.ID)
	}

	// Mark 1 as missing
	reqBody, _ := json.Marshal(map[string]string{"reason": "Mất hóa đơn"})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-numbers/"+allocatedIDs[0]+"/missing", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Get summary
	req = httptest.NewRequest(http.MethodGet, "/api/v1/invoice-books/"+created.ID+"/summary", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var summary map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &summary)
	if summary["issued"].(float64) != 2 {
		t.Errorf("expected 2 issued, got %.0f", summary["issued"].(float64))
	}
	if summary["missing"].(float64) != 1 {
		t.Errorf("expected 1 missing, got %.0f", summary["missing"].(float64))
	}
	if summary["available"].(float64) != 7 {
		t.Errorf("expected 7 available, got %.0f", summary["available"].(float64))
	}
}

func TestGetMissingReport(t *testing.T) {
	r, h := setupInvoiceBookReportTestRouter()
	r.POST("/api/v1/invoice-books", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.CreateBook)
	r.POST("/api/v1/invoice-books/:id/allocate", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.AllocateNumber)
	r.POST("/api/v1/invoice-numbers/:id/missing", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.MarkMissing)
	r.GET("/api/v1/invoice-books/:id/missing-report", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.GetMissingReport)

	// Create book and allocate
	book := domain.InvoiceBook{Name: "Missing Test", Pattern: "01GTTT", FromNumber: 1, ToNumber: 5}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books?company_id=CMP001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var created domain.InvoiceBook
	json.Unmarshal(w.Body.Bytes(), &created)

	// Allocate 2 numbers
	var allocatedIDs []string
	for i := 0; i < 2; i++ {
		req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books/"+created.ID+"/allocate", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var num domain.InvoiceNumber
		json.Unmarshal(w.Body.Bytes(), &num)
		allocatedIDs = append(allocatedIDs, num.ID)
	}

	// Mark both as missing with different reasons
	for i, id := range allocatedIDs {
		reason := "Mất hóa đơn"
		if i == 1 {
			reason = "Hỏa hóa đơn"
		}
		reqBody, _ := json.Marshal(map[string]string{"reason": reason})
		req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-numbers/"+id+"/missing", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	// Get missing report
	req = httptest.NewRequest(http.MethodGet, "/api/v1/invoice-books/"+created.ID+"/missing-report", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var report map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &report)
	if report["total_missing"].(float64) != 2 {
		t.Errorf("expected 2 missing, got %.0f", report["total_missing"].(float64))
	}
}
