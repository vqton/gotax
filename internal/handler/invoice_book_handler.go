package handler

import (
	"net/http"

	"gotax/internal/domain"
	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

type InvoiceBookHandler struct {
	svc *service.InvoiceBookService
}

func NewInvoiceBookHandler(svc *service.InvoiceBookService) *InvoiceBookHandler {
	return &InvoiceBookHandler{svc: svc}
}

func RegisterInvoiceBookRoutes(r *gin.Engine, h *InvoiceBookHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	{
		books := v1.Group("/invoice-books")
		{
			books.POST("", h.CreateBook)
			books.GET("", h.ListBooks)
			books.GET("/:id", h.GetBook)
			books.PUT("/:id", h.UpdateBook)
			books.DELETE("/:id", h.DeleteBook)
			books.POST("/:id/allocate", h.AllocateNumber)
			books.GET("/:id/numbers", h.GetNumbers)
			books.GET("/:id/available-count", h.GetAvailableCount)
		}
		numbers := v1.Group("/invoice-numbers")
		{
			numbers.POST("/:id/release", h.ReleaseNumber)
			numbers.POST("/:id/missing", h.MarkMissing)
		}
		reports := v1.Group("/invoice-books")
		{
			reports.GET("/:id/summary", h.GetBookSummary)
			reports.GET("/:id/missing-report", h.GetMissingReport)
		}
	}
}

func (h *InvoiceBookHandler) CreateBook(c *gin.Context) {
	var book domain.InvoiceBook
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book.CompanyID = c.Query("company_id")
	if book.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}
	if err := h.svc.CreateBook(c.Request.Context(), &book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
}

func (h *InvoiceBookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")
	book, err := h.svc.GetBook(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *InvoiceBookHandler) ListBooks(c *gin.Context) {
	companyID := c.Query("company_id")
	items, err := h.svc.ListBooks(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *InvoiceBookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book domain.InvoiceBook
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book.ID = id
	if err := h.svc.UpdateBook(c.Request.Context(), &book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *InvoiceBookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteBook(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *InvoiceBookHandler) AllocateNumber(c *gin.Context) {
	bookID := c.Param("id")
	num, err := h.svc.AllocateNumber(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, num)
}

func (h *InvoiceBookHandler) GetNumbers(c *gin.Context) {
	bookID := c.Param("id")
	nums, err := h.svc.GetNumbersByBook(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": nums})
}

func (h *InvoiceBookHandler) GetAvailableCount(c *gin.Context) {
	bookID := c.Param("id")
	count, err := h.svc.GetAvailableCount(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"available": count})
}

func (h *InvoiceBookHandler) ReleaseNumber(c *gin.Context) {
	numberID := c.Param("id")
	if err := h.svc.ReleaseNumber(c.Request.Context(), numberID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "number released"})
}

func (h *InvoiceBookHandler) MarkMissing(c *gin.Context) {
	numberID := c.Param("id")
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason required"})
		return
	}
	if err := h.svc.MarkMissing(c.Request.Context(), numberID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "number marked as missing"})
}

func (h *InvoiceBookHandler) GetBookSummary(c *gin.Context) {
	bookID := c.Param("id")
	book, err := h.svc.GetBook(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	nums, err := h.svc.GetNumbersByBook(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	issued, missing, damaged := 0, 0, 0
	for _, n := range nums {
		switch n.Status {
		case domain.InvNumIssued:
			issued++
		case domain.InvNumMissing:
			missing++
		case domain.InvNumDamaged:
			damaged++
		}
	}
	available := book.ToNumber - book.NextNumber + 1
	if available < 0 {
		available = 0
	}
	c.JSON(http.StatusOK, gin.H{
		"book":      book,
		"total":     book.ToNumber - book.FromNumber + 1,
		"issued":    issued,
		"missing":   missing,
		"damaged":   damaged,
		"available": available,
	})
}

func (h *InvoiceBookHandler) GetMissingReport(c *gin.Context) {
	bookID := c.Param("id")
	book, err := h.svc.GetBook(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	nums, err := h.svc.GetNumbersByBook(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var missingList []gin.H
	for _, n := range nums {
		if n.Status == domain.InvNumMissing || n.Status == domain.InvNumDamaged {
			missingList = append(missingList, gin.H{
				"id":            n.ID,
				"number":        n.Number,
				"status":        n.Status,
				"missing_reason": n.MissingReason,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"book_name":   book.Name,
		"pattern":     book.Pattern,
		"serial":      book.Serial,
		"missing":     missingList,
		"total_missing": len(missingList),
	})
}
