package handler

import (
	"net/http"
	"strconv"

	"gotax/internal/domain"
	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

type SaleHandler struct {
	svc *service.SaleService
}

func NewSaleHandler(svc *service.SaleService) *SaleHandler {
	return &SaleHandler{svc: svc}
}

func RegisterSaleRoutes(r *gin.Engine, h *SaleHandler, authMW gin.HandlerFunc) {
	sale := r.Group("/api/v1/sale", authMW)
	{
		customers := sale.Group("/customers")
		{
			customers.POST("", h.CreateCustomer)
			customers.GET("", h.ListCustomers)
			customers.GET("/:id", h.GetCustomer)
			customers.PUT("/:id", h.UpdateCustomer)
			customers.DELETE("/:id", h.DeleteCustomer)
		}
		orders := sale.Group("/orders")
		{
			orders.POST("", h.CreateSO)
			orders.GET("", h.ListSOs)
			orders.GET("/:id", h.GetSO)
			orders.PUT("/:id", h.UpdateSO)
			orders.PATCH("/:id/approve", h.ApproveSO)
			orders.PATCH("/:id/cancel", h.CancelSO)
			orders.PATCH("/:id/close", h.CloseSO)
		}
		deliveries := sale.Group("/deliveries")
		{
			deliveries.POST("", h.CreateDN)
			deliveries.GET("", h.ListDNs)
			deliveries.GET("/:id", h.GetDN)
			deliveries.PUT("/:id", h.UpdateDN)
			deliveries.PATCH("/:id/post", h.PostDN)
			deliveries.PATCH("/:id/cancel", h.CancelDN)
		}
		invoices := sale.Group("/invoices")
		{
			invoices.POST("", h.CreateInvoice)
			invoices.GET("", h.ListInvoices)
			invoices.GET("/:id", h.GetInvoice)
			invoices.PUT("/:id", h.UpdateInvoice)
			invoices.PATCH("/:id/post", h.PostInvoice)
			invoices.PATCH("/:id/cancel", h.CancelInvoice)
		}
		receipts := sale.Group("/receipts")
		{
			receipts.POST("", h.CreateReceipt)
			receipts.GET("", h.ListReceipts)
			receipts.GET("/:id", h.GetReceipt)
			receipts.PATCH("/:id/post", h.PostReceipt)
			receipts.PATCH("/:id/cancel", h.CancelReceipt)
			receipts.POST("/:id/allocate", h.AllocateReceipt)
		}
		creditnotes := sale.Group("/credit-notes")
		{
			creditnotes.POST("", h.CreateCN)
			creditnotes.GET("", h.ListCNs)
			creditnotes.GET("/:id", h.GetCN)
			creditnotes.PATCH("/:id/post", h.PostCN)
			creditnotes.PATCH("/:id/cancel", h.CancelCN)
		}
		ar := sale.Group("/ar")
		{
			ar.GET("/aging", h.GetARAgingReport)
			ar.GET("/summary", h.GetARSummary)
			ar.GET("/statement", h.GetCustomerStatement)
		}
		quotations := sale.Group("/quotations")
		{
			quotations.POST("", h.CreateSQ)
			quotations.GET("", h.ListSQs)
			quotations.GET("/:id", h.GetSQ)
			quotations.PUT("/:id", h.UpdateSQ)
		}
	}
}

// ─── Customer ──────────────────────────────────────────────────────────

func (h *SaleHandler) CreateCustomer(c *gin.Context) {
	var cust domain.Customer
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cust.CompanyID = c.Query("company_id")
	if cust.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	if err := h.svc.CreateCustomer(c.Request.Context(), &cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cust)
}

func (h *SaleHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	cust, err := h.svc.GetCustomer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if cust.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.JSON(http.StatusOK, cust)
}

func (h *SaleHandler) ListCustomers(c *gin.Context) {
	companyID := c.Query("company_id")
	customers, err := h.svc.ListCustomers(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": customers, "total": len(customers)})
}

func (h *SaleHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.svc.GetCustomer(c.Request.Context(), id)
	if err != nil || existing.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	var cust domain.Customer
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cust.ID = id
	cust.CompanyID = existing.CompanyID
	if err := h.svc.UpdateCustomer(c.Request.Context(), &cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cust)
}

func (h *SaleHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.svc.GetCustomer(c.Request.Context(), id)
	if err != nil || existing.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	if err := h.svc.DeleteCustomer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ─── Sales Order ───────────────────────────────────────────────────────

func (h *SaleHandler) CreateSO(c *gin.Context) {
	var so domain.SalesOrder
	if err := c.ShouldBindJSON(&so); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	so.CompanyID = c.Query("company_id")
	so.CreatedBy = GetUserID(c)
	if err := h.svc.CreateSO(c.Request.Context(), &so); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, so)
}

func (h *SaleHandler) requireCompany(c *gin.Context) (string, bool) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return "", false
	}
	return companyID, true
}

func (h *SaleHandler) GetSO(c *gin.Context) {
	id := c.Param("id")
	so, err := h.svc.GetSO(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if so.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "SO not found"})
		return
	}
	c.JSON(http.StatusOK, so)
}

func (h *SaleHandler) ListSOs(c *gin.Context) {
	filter := domain.SalesOrderFilter{
		CompanyID:  c.Query("company_id"),
		CustomerID: c.Query("customer_id"),
		Status:     domain.SOStatus(c.Query("status")),
		FromDate:   c.Query("from_date"),
		ToDate:     c.Query("to_date"),
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	filter.Offset = offset
	filter.Limit = limit
	sos, total, err := h.svc.ListSOs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sos, "total": total})
}

func (h *SaleHandler) UpdateSO(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.svc.GetSO(c.Request.Context(), id)
	if err != nil || existing.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "SO not found"})
		return
	}
	var so domain.SalesOrder
	if err := c.ShouldBindJSON(&so); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	so.ID = id
	so.CompanyID = existing.CompanyID
	if err := h.svc.UpdateSO(c.Request.Context(), &so); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, so)
}

func (h *SaleHandler) ApproveSO(c *gin.Context) {
	id := c.Param("id")
	so, err := h.svc.GetSO(c.Request.Context(), id)
	if err != nil || so.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "SO not found"})
		return
	}
	if err := h.svc.ApproveSO(c.Request.Context(), id, GetUserID(c)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *SaleHandler) CancelSO(c *gin.Context) {
	id := c.Param("id")
	so, err := h.svc.GetSO(c.Request.Context(), id)
	if err != nil || so.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "SO not found"})
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.CancelSO(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

func (h *SaleHandler) CloseSO(c *gin.Context) {
	id := c.Param("id")
	so, err := h.svc.GetSO(c.Request.Context(), id)
	if err != nil || so.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "SO not found"})
		return
	}
	if err := h.svc.CloseSO(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "closed"})
}

// ─── Delivery Note ─────────────────────────────────────────────────────

func (h *SaleHandler) CreateDN(c *gin.Context) {
	var dn domain.DeliveryNote
	if err := c.ShouldBindJSON(&dn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dn.CompanyID = c.Query("company_id")
	dn.CreatedBy = GetUserID(c)
	if err := h.svc.CreateDN(c.Request.Context(), &dn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dn)
}

func (h *SaleHandler) GetDN(c *gin.Context) {
	id := c.Param("id")
	dn, err := h.svc.GetDN(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if dn.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "DN not found"})
		return
	}
	c.JSON(http.StatusOK, dn)
}

func (h *SaleHandler) ListDNs(c *gin.Context) {
	filter := domain.DeliveryNoteFilter{
		CompanyID: c.Query("company_id"),
		SOID:      c.Query("so_id"),
		Status:    domain.DNStatus(c.Query("status")),
		FromDate:  c.Query("from_date"),
		ToDate:    c.Query("to_date"),
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	filter.Offset = offset
	filter.Limit = limit
	dns, total, err := h.svc.ListDNs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": dns, "total": total})
}

func (h *SaleHandler) UpdateDN(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.svc.GetDN(c.Request.Context(), id)
	if err != nil || existing.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "DN not found"})
		return
	}
	var dn domain.DeliveryNote
	if err := c.ShouldBindJSON(&dn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dn.ID = id
	dn.CompanyID = existing.CompanyID
	if err := h.svc.UpdateDN(c.Request.Context(), &dn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dn)
}

func (h *SaleHandler) PostDN(c *gin.Context) {
	id := c.Param("id")
	dn, err := h.svc.GetDN(c.Request.Context(), id)
	if err != nil || dn.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "DN not found"})
		return
	}
	if err := h.svc.PostDN(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posted"})
}

func (h *SaleHandler) CancelDN(c *gin.Context) {
	id := c.Param("id")
	dn, err := h.svc.GetDN(c.Request.Context(), id)
	if err != nil || dn.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "DN not found"})
		return
	}
	if err := h.svc.CancelDN(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

// ─── Invoice ───────────────────────────────────────────────────────────

func (h *SaleHandler) CreateInvoice(c *gin.Context) {
	var inv domain.CustomerInvoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inv.CompanyID = c.Query("company_id")
	inv.CreatedBy = GetUserID(c)
	if err := h.svc.CreateInvoice(c.Request.Context(), &inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, inv)
}

func (h *SaleHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	inv, err := h.svc.GetInvoice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if inv.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *SaleHandler) ListInvoices(c *gin.Context) {
	filter := domain.CustomerInvoiceFilter{
		CompanyID:  c.Query("company_id"),
		CustomerID: c.Query("customer_id"),
		Status:     domain.SaleInvoiceStatus(c.Query("status")),
		FromDate:   c.Query("from_date"),
		ToDate:     c.Query("to_date"),
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	filter.Offset = offset
	filter.Limit = limit
	invs, total, err := h.svc.ListInvoices(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": invs, "total": total})
}

func (h *SaleHandler) UpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.svc.GetInvoice(c.Request.Context(), id)
	if err != nil || existing.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}
	var inv domain.CustomerInvoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inv.ID = id
	inv.CompanyID = existing.CompanyID
	if err := h.svc.UpdateInvoice(c.Request.Context(), &inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *SaleHandler) PostInvoice(c *gin.Context) {
	id := c.Param("id")
	inv, err := h.svc.GetInvoice(c.Request.Context(), id)
	if err != nil || inv.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}
	if err := h.svc.PostInvoice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posted"})
}

func (h *SaleHandler) CancelInvoice(c *gin.Context) {
	id := c.Param("id")
	inv, err := h.svc.GetInvoice(c.Request.Context(), id)
	if err != nil || inv.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}
	if err := h.svc.CancelInvoice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

// ─── Receipt ───────────────────────────────────────────────────────────

func (h *SaleHandler) CreateReceipt(c *gin.Context) {
	var rcpt domain.CustomerReceipt
	if err := c.ShouldBindJSON(&rcpt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rcpt.CompanyID = c.Query("company_id")
	rcpt.CreatedBy = GetUserID(c)
	if err := h.svc.CreateReceipt(c.Request.Context(), &rcpt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rcpt)
}

func (h *SaleHandler) GetReceipt(c *gin.Context) {
	id := c.Param("id")
	rcpt, err := h.svc.GetReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if rcpt.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "receipt not found"})
		return
	}
	c.JSON(http.StatusOK, rcpt)
}

func (h *SaleHandler) ListReceipts(c *gin.Context) {
	filter := domain.ReceiptFilter{
		CompanyID:  c.Query("company_id"),
		CustomerID: c.Query("customer_id"),
		Status:     domain.ReceiptStatus(c.Query("status")),
		FromDate:   c.Query("from_date"),
		ToDate:     c.Query("to_date"),
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	filter.Offset = offset
	filter.Limit = limit
	receipts, total, err := h.svc.ListReceipts(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": receipts, "total": total})
}

func (h *SaleHandler) PostReceipt(c *gin.Context) {
	id := c.Param("id")
	rcpt, err := h.svc.GetReceipt(c.Request.Context(), id)
	if err != nil || rcpt.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "receipt not found"})
		return
	}
	if err := h.svc.PostReceipt(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posted"})
}

func (h *SaleHandler) CancelReceipt(c *gin.Context) {
	id := c.Param("id")
	rcpt, err := h.svc.GetReceipt(c.Request.Context(), id)
	if err != nil || rcpt.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "receipt not found"})
		return
	}
	if err := h.svc.CancelReceipt(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

func (h *SaleHandler) AllocateReceipt(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		InvoiceID string  `json:"invoice_id"`
		Amount    float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rcpt, err := h.svc.GetReceipt(c.Request.Context(), id)
	if err != nil || rcpt.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "receipt not found"})
		return
	}
	if err := h.svc.AllocateReceipt(c.Request.Context(), id, req.InvoiceID, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "allocated"})
}

// ─── Credit Note ───────────────────────────────────────────────────────

func (h *SaleHandler) CreateCN(c *gin.Context) {
	var cn domain.CreditNote
	if err := c.ShouldBindJSON(&cn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cn.CompanyID = c.Query("company_id")
	cn.CreatedBy = GetUserID(c)
	if err := h.svc.CreateCN(c.Request.Context(), &cn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cn)
}

func (h *SaleHandler) GetCN(c *gin.Context) {
	id := c.Param("id")
	cn, err := h.svc.GetCN(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if cn.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "credit note not found"})
		return
	}
	c.JSON(http.StatusOK, cn)
}

func (h *SaleHandler) ListCNs(c *gin.Context) {
	filter := domain.CreditNoteFilter{
		CompanyID:  c.Query("company_id"),
		CustomerID: c.Query("customer_id"),
		Status:     domain.CNStatus(c.Query("status")),
		FromDate:   c.Query("from_date"),
		ToDate:     c.Query("to_date"),
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	filter.Offset = offset
	filter.Limit = limit
	cns, total, err := h.svc.ListCNs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cns, "total": total})
}

func (h *SaleHandler) PostCN(c *gin.Context) {
	id := c.Param("id")
	cn, err := h.svc.GetCN(c.Request.Context(), id)
	if err != nil || cn.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "credit note not found"})
		return
	}
	if err := h.svc.PostCN(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posted"})
}

func (h *SaleHandler) CancelCN(c *gin.Context) {
	id := c.Param("id")
	cn, err := h.svc.GetCN(c.Request.Context(), id)
	if err != nil || cn.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "credit note not found"})
		return
	}
	if err := h.svc.CancelCN(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

// ─── AR Reports ────────────────────────────────────────────────────────

func (h *SaleHandler) GetARAgingReport(c *gin.Context) {
	companyID := c.Query("company_id")
	report, err := h.svc.GetARAgingReport(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *SaleHandler) GetARSummary(c *gin.Context) {
	companyID := c.Query("company_id")
	summary, err := h.svc.GetARSummary(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *SaleHandler) GetCustomerStatement(c *gin.Context) {
	customerID := c.Query("customer_id")
	fromDate := c.DefaultQuery("from_date", "")
	toDate := c.DefaultQuery("to_date", "")
	stmt, err := h.svc.GetCustomerStatement(c.Request.Context(), customerID, fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stmt)
}

// ─── Sales Quotation ───────────────────────────────────────────────────

func (h *SaleHandler) CreateSQ(c *gin.Context) {
	var sq domain.SalesQuotation
	if err := c.ShouldBindJSON(&sq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sq.CompanyID = c.Query("company_id")
	sq.CreatedBy = GetUserID(c)
	if err := h.svc.CreateSQ(c.Request.Context(), &sq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sq)
}

func (h *SaleHandler) GetSQ(c *gin.Context) {
	id := c.Param("id")
	sq, err := h.svc.GetSQ(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if sq.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "quotation not found"})
		return
	}
	c.JSON(http.StatusOK, sq)
}

func (h *SaleHandler) ListSQs(c *gin.Context) {
	companyID := c.Query("company_id")
	sqs, err := h.svc.ListSQs(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sqs, "total": len(sqs)})
}

func (h *SaleHandler) UpdateSQ(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.svc.GetSQ(c.Request.Context(), id)
	if err != nil || existing.CompanyID != c.Query("company_id") {
		c.JSON(http.StatusNotFound, gin.H{"error": "quotation not found"})
		return
	}
	var sq domain.SalesQuotation
	if err := c.ShouldBindJSON(&sq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sq.ID = id
	sq.CompanyID = existing.CompanyID
	if err := h.svc.UpdateSQ(c.Request.Context(), &sq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sq)
}
