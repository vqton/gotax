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

func setupInvoiceBookTestRouter() (*gin.Engine, *InvoiceBookHandler) {
	gin.SetMode(gin.TestMode)
	repo := repository.NewMemoryInvoiceBookRepo()
	svc := service.NewInvoiceBookService(repo)
	h := NewInvoiceBookHandler(svc)
	r := gin.New()
	return r, h
}

func TestCreateBook(t *testing.T) {
	r, h := setupInvoiceBookTestRouter()
	r.POST("/api/v1/invoice-books", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.CreateBook)

	book := domain.InvoiceBook{
		Name:       "Hóa đơn GTGT",
		Pattern:    "01GTTT",
		Serial:    "AA/25E",
		FromNumber: 1,
		ToNumber:   1000,
	}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books?company_id=CMP001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp domain.InvoiceBook
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
	if resp.NextNumber != 1 {
		t.Errorf("expected next_number=1, got %d", resp.NextNumber)
	}
}

func TestCreateBook_InvalidRange(t *testing.T) {
	r, h := setupInvoiceBookTestRouter()
	r.POST("/api/v1/invoice-books", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.CreateBook)

	book := domain.InvoiceBook{
		Name:       "Bad Book",
		Pattern:    "01GTTT",
		FromNumber: 100,
		ToNumber:   10,
	}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books?company_id=CMP001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAllocateAndReleaseNumber(t *testing.T) {
	r, h := setupInvoiceBookTestRouter()
	r.POST("/api/v1/invoice-books", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.CreateBook)
	r.POST("/api/v1/invoice-books/:id/allocate", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.AllocateNumber)
	r.POST("/api/v1/invoice-numbers/:id/release", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.ReleaseNumber)
	r.GET("/api/v1/invoice-books/:id/available-count", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.GetAvailableCount)

	// Create book
	book := domain.InvoiceBook{Name: "Test", Pattern: "01GTTT", FromNumber: 1, ToNumber: 5}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books?company_id=CMP001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var created domain.InvoiceBook
	json.Unmarshal(w.Body.Bytes(), &created)

	// Allocate number
	req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books/"+created.ID+"/allocate", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var num domain.InvoiceNumber
	json.Unmarshal(w.Body.Bytes(), &num)
	if num.Number != 1 {
		t.Errorf("expected number=1, got %d", num.Number)
	}

	// Check available count
	req = httptest.NewRequest(http.MethodGet, "/api/v1/invoice-books/"+created.ID+"/available-count", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var countResp map[string]int
	json.Unmarshal(w.Body.Bytes(), &countResp)
	if countResp["available"] != 4 {
		t.Errorf("expected 4 available, got %d", countResp["available"])
	}

	// Release number
	req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-numbers/"+num.ID+"/release", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAllocateNumber_BookFull(t *testing.T) {
	r, h := setupInvoiceBookTestRouter()
	r.POST("/api/v1/invoice-books", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.CreateBook)
	r.POST("/api/v1/invoice-books/:id/allocate", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.AllocateNumber)

	// Create book with only 2 numbers
	book := domain.InvoiceBook{Name: "Tiny", Pattern: "01GTTT", FromNumber: 1, ToNumber: 2}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books?company_id=CMP001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var created domain.InvoiceBook
	json.Unmarshal(w.Body.Bytes(), &created)

	// Allocate 2 numbers (should succeed)
	for i := 0; i < 2; i++ {
		req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books/"+created.ID+"/allocate", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("allocate %d: expected 201, got %d", i+1, w.Code)
		}
	}

	// Third allocation should fail
	req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books/"+created.ID+"/allocate", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for full book, got %d", w.Code)
	}
}

func TestMarkMissing(t *testing.T) {
	r, h := setupInvoiceBookTestRouter()
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

	// Create book and allocate
	book := domain.InvoiceBook{Name: "Test", Pattern: "01GTTT", FromNumber: 1, ToNumber: 10}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books?company_id=CMP001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var created domain.InvoiceBook
	json.Unmarshal(w.Body.Bytes(), &created)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books/"+created.ID+"/allocate", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var num domain.InvoiceNumber
	json.Unmarshal(w.Body.Bytes(), &num)

	// Mark as missing
	reqBody, _ := json.Marshal(map[string]string{"reason": "Mất hóa đơn"})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-numbers/"+num.ID+"/missing", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteBook_WithIssuedNumbers(t *testing.T) {
	r, h := setupInvoiceBookTestRouter()
	r.POST("/api/v1/invoice-books", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.CreateBook)
	r.POST("/api/v1/invoice-books/:id/allocate", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.AllocateNumber)
	r.DELETE("/api/v1/invoice-books/:id", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	}, h.DeleteBook)

	// Create and allocate
	book := domain.InvoiceBook{Name: "Test", Pattern: "01GTTT", FromNumber: 1, ToNumber: 10}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books?company_id=CMP001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var created domain.InvoiceBook
	json.Unmarshal(w.Body.Bytes(), &created)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/invoice-books/"+created.ID+"/allocate", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Try to delete — should fail
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/invoice-books/"+created.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
