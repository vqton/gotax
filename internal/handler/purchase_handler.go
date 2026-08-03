package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"gotax/internal/domain"
	"gotax/internal/service"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	svc *service.PurchaseService
}

func NewPurchaseHandler(svc *service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{svc: svc}
}

func RegisterPurchaseRoutes(r *gin.Engine, h *PurchaseHandler, authMW gin.HandlerFunc) {
	purchase := r.Group("/api/v1/purchase", authMW)
	{
		suppliers := purchase.Group("/suppliers")
		{
			suppliers.POST("", h.CreateSupplier)
			suppliers.GET("", h.ListSuppliers)
			suppliers.GET("/:id", h.GetSupplier)
			suppliers.PUT("/:id", h.UpdateSupplier)
			suppliers.DELETE("/:id", h.DeleteSupplier)
		}
		orders := purchase.Group("/orders")
		{
			orders.POST("", h.CreatePO)
			orders.GET("", h.ListPOs)
			orders.GET("/:id", h.GetPO)
			orders.PUT("/:id", h.UpdatePO)
			orders.PATCH("/:id/approve", h.ApprovePO)
			orders.PATCH("/:id/cancel", h.CancelPO)
			orders.PATCH("/:id/close", h.ClosePO)
		}
		requisitions := purchase.Group("/requisitions")
		{
			requisitions.POST("", h.CreateRequisition)
			requisitions.GET("", h.ListRequisitions)
			requisitions.GET("/:id", h.GetRequisition)
			requisitions.PUT("/:id", h.UpdateRequisition)
			requisitions.DELETE("/:id", h.DeleteRequisition)
			requisitions.PATCH("/:id/submit", h.SubmitRequisition)
			requisitions.PATCH("/:id/approve", h.ApproveRequisition)
			requisitions.PATCH("/:id/reject", h.RejectRequisition)
			requisitions.POST("/:id/convert-to-po", h.ConvertRequisitionToPO)
		}
		receipts := purchase.Group("/receipts")
		{
			receipts.POST("", h.CreateGRN)
			receipts.GET("", h.ListGRNs)
			receipts.GET("/:id", h.GetGRN)
			receipts.PUT("/:id", h.UpdateGRN)
			receipts.PATCH("/:id/post", h.PostGRN)
			receipts.PATCH("/:id/cancel", h.CancelGRN)
		}
		invoices := purchase.Group("/invoices")
		{
			invoices.POST("", h.CreateInvoice)
			invoices.GET("", h.ListInvoices)
			invoices.GET("/:id", h.GetInvoice)
			invoices.PUT("/:id", h.UpdateInvoice)
			invoices.PATCH("/:id/verify", h.VerifyInvoice)
			invoices.PATCH("/:id/post", h.PostInvoice)
			invoices.PATCH("/:id/cancel", h.CancelInvoice)
			invoices.PATCH("/:id/claim-vat", h.ClaimVAT)
		}
		ap := purchase.Group("/ap")
		{
			ap.GET("/aging", h.GetAPAgingReport)
			ap.GET("/summary", h.GetAPSummary)
		}
		provisions := purchase.Group("/provisions")
		{
			provisions.GET("/calculate", h.CalculateDoubtfulDebtProvision)
			provisions.POST("", h.CreateDoubtfulDebtProvision)
			provisions.GET("", h.ListDoubtfulDebtProvisions)
			provisions.GET("/:id", h.GetDoubtfulDebtProvision)
		}
		reports := purchase.Group("/reports")
		{
			reports.GET("/s01-dn", h.GetPurchaseLedger)
			reports.GET("/s02-dn", h.GetSupplierLedger)
			reports.GET("/s03-dn", h.GetGoodsPurchaseReport)
			reports.GET("/vat-input", h.GetVATInputReport)
			reports.GET("/uninvoiced-receipts", h.GetUninvoicedReceipts)
		}
	}
}

// ─── Supplier Handlers ───────────────────────────────────────────────────

func (h *PurchaseHandler) CreateSupplier(c *gin.Context) {
	var sup domain.Supplier
	if err := c.ShouldBindJSON(&sup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sup.CompanyID = c.Query("company_id")
	if sup.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	if err := h.svc.CreateSupplier(c.Request.Context(), &sup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sup)
}

func (h *PurchaseHandler) GetSupplier(c *gin.Context) {
	id := c.Param("id")
	sup, err := h.svc.GetSupplier(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sup)
}

func (h *PurchaseHandler) ListSuppliers(c *gin.Context) {
	companyID := c.Query("company_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	suppliers, total, err := h.svc.ListSuppliers(c.Request.Context(), companyID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": suppliers, "total": total})
}

func (h *PurchaseHandler) UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	var sup domain.Supplier
	if err := c.ShouldBindJSON(&sup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sup.ID = id
	sup.CompanyID = c.Query("company_id")
	if err := h.svc.UpdateSupplier(c.Request.Context(), &sup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sup)
}

func (h *PurchaseHandler) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteSupplier(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ─── Purchase Order Handlers ─────────────────────────────────────────────

func (h *PurchaseHandler) CreatePO(c *gin.Context) {
	var po domain.PurchaseOrder
	if err := c.ShouldBindJSON(&po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	po.CompanyID = c.Query("company_id")
	po.CreatedBy = GetUserID(c)
	if err := h.svc.CreatePO(c.Request.Context(), &po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, po)
}

func (h *PurchaseHandler) GetPO(c *gin.Context) {
	id := c.Param("id")
	po, err := h.svc.GetPO(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, po)
}

func (h *PurchaseHandler) ListPOs(c *gin.Context) {
	filter := domain.PurchaseOrderFilter{
		CompanyID:  c.Query("company_id"),
		SupplierID: c.Query("supplier_id"),
		Status:     domain.POStatus(c.Query("status")),
		FromDate:   c.Query("from_date"),
		ToDate:     c.Query("to_date"),
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	filter.Offset = offset
	filter.Limit = limit
	pos, total, err := h.svc.ListPOs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pos, "total": total})
}

func (h *PurchaseHandler) UpdatePO(c *gin.Context) {
	id := c.Param("id")
	var po domain.PurchaseOrder
	if err := c.ShouldBindJSON(&po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	po.ID = id
	if err := h.svc.UpdatePO(c.Request.Context(), &po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, po)
}

func (h *PurchaseHandler) ApprovePO(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ApprovePO(c.Request.Context(), id, GetUserID(c)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *PurchaseHandler) CancelPO(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.CancelPO(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

func (h *PurchaseHandler) ClosePO(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ClosePO(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "closed"})
}

// ─── GRN Handlers ────────────────────────────────────────────────────────

func (h *PurchaseHandler) CreateGRN(c *gin.Context) {
	var grn domain.GRN
	if err := c.ShouldBindJSON(&grn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grn.CompanyID = c.Query("company_id")
	grn.CreatedBy = GetUserID(c)
	if err := h.svc.CreateGRN(c.Request.Context(), &grn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, grn)
}

func (h *PurchaseHandler) GetGRN(c *gin.Context) {
	id := c.Param("id")
	grn, err := h.svc.GetGRN(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grn)
}

func (h *PurchaseHandler) ListGRNs(c *gin.Context) {
	filter := domain.GRNFilter{
		CompanyID: c.Query("company_id"),
		POID:      c.Query("po_id"),
		Status:    domain.GRNStatus(c.Query("status")),
		FromDate:  c.Query("from_date"),
		ToDate:    c.Query("to_date"),
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	filter.Offset = offset
	filter.Limit = limit
	grns, total, err := h.svc.ListGRNs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": grns, "total": total})
}

func (h *PurchaseHandler) UpdateGRN(c *gin.Context) {
	id := c.Param("id")
	var grn domain.GRN
	if err := c.ShouldBindJSON(&grn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grn.ID = id
	if err := h.svc.UpdateGRN(c.Request.Context(), &grn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grn)
}

func (h *PurchaseHandler) PostGRN(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.PostGRN(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posted"})
}

func (h *PurchaseHandler) CancelGRN(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.CancelGRN(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

// ─── Invoice Handlers ────────────────────────────────────────────────────

func (h *PurchaseHandler) CreateInvoice(c *gin.Context) {
	var inv domain.SupplierInvoice
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

func (h *PurchaseHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	inv, err := h.svc.GetInvoice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *PurchaseHandler) ListInvoices(c *gin.Context) {
	filter := domain.SupplierInvoiceFilter{
		CompanyID:  c.Query("company_id"),
		SupplierID: c.Query("supplier_id"),
		Status:     domain.InvoiceStatus(c.Query("status")),
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

func (h *PurchaseHandler) UpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	var inv domain.SupplierInvoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inv.ID = id
	if err := h.svc.UpdateInvoice(c.Request.Context(), &inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *PurchaseHandler) VerifyInvoice(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.VerifyInvoice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "verified"})
}

func (h *PurchaseHandler) PostInvoice(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.PostInvoice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "posted"})
}

func (h *PurchaseHandler) CancelInvoice(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.CancelInvoice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

func (h *PurchaseHandler) ClaimVAT(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ClaimVAT(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "VAT claimed"})
}

// ─── AP Report Handlers ──────────────────────────────────────────────────

func (h *PurchaseHandler) GetAPAgingReport(c *gin.Context) {
	companyID := c.Query("company_id")
	report, err := h.svc.GetAPAgingReport(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *PurchaseHandler) GetAPSummary(c *gin.Context) {
	companyID := c.Query("company_id")
	summary, err := h.svc.GetAPSummary(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

// ─── Doubtful Debt Provision Handlers ───────────────────────────────────

func (h *PurchaseHandler) CalculateDoubtfulDebtProvision(c *gin.Context) {
	companyID := c.Query("company_id")
	asOf := c.Query("as_of_date")
	if asOf == "" {
		asOf = time.Now().Format(time.DateOnly)
	}
	asOfDate, err := time.Parse(time.DateOnly, asOf)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrProvisionDateRequired.Error()})
		return
	}
	lines, err := h.svc.CalculateDoubtfulDebtProvision(c.Request.Context(), companyID, asOfDate)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrProvisionNoPrepayments) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"as_of_date": asOf, "lines": lines})
}

func (h *PurchaseHandler) CreateDoubtfulDebtProvision(c *gin.Context) {
	var prov domain.DoubtfulDebtProvision
	if err := c.ShouldBindJSON(&prov); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	prov.CompanyID = c.Query("company_id")
	prov.CreatedBy = c.GetString("user_id")
	if err := h.svc.CreateDoubtfulDebtProvision(c.Request.Context(), &prov); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, domain.ErrProvisionNoLines) || errors.Is(err, domain.ErrProvisionDateRequired) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, prov)
}

func (h *PurchaseHandler) GetDoubtfulDebtProvision(c *gin.Context) {
	prov, err := h.svc.GetDoubtfulDebtProvision(c.Request.Context(), c.Param("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrProvisionNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prov)
}

func (h *PurchaseHandler) ListDoubtfulDebtProvisions(c *gin.Context) {
	companyID := c.Query("company_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	provisions, total, err := h.svc.ListDoubtfulDebtProvisions(c.Request.Context(), companyID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": provisions, "total": total})
}

// ─── Report Handlers ─────────────────────────────────────────────────────

func (h *PurchaseHandler) GetPurchaseLedger(c *gin.Context) {
	companyID := c.Query("company_id")
	rpt, err := h.svc.GetPurchaseLedger(c.Request.Context(), companyID, c.Query("from_date"), c.Query("to_date"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rpt)
}

func (h *PurchaseHandler) GetSupplierLedger(c *gin.Context) {
	companyID := c.Query("company_id")
	supplierID := c.Query("supplier_id")
	if supplierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supplier_id query param required"})
		return
	}
	rpt, err := h.svc.GetSupplierLedger(c.Request.Context(), companyID, supplierID, c.Query("from_date"), c.Query("to_date"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rpt)
}

func (h *PurchaseHandler) GetGoodsPurchaseReport(c *gin.Context) {
	companyID := c.Query("company_id")
	rpt, err := h.svc.GetGoodsPurchaseReport(c.Request.Context(), companyID, c.Query("from_date"), c.Query("to_date"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rpt)
}

func (h *PurchaseHandler) GetVATInputReport(c *gin.Context) {
	companyID := c.Query("company_id")
	rpt, err := h.svc.GetVATInputReport(c.Request.Context(), companyID, c.Query("from_date"), c.Query("to_date"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rpt)
}

func (h *PurchaseHandler) GetUninvoicedReceipts(c *gin.Context) {
	companyID := c.Query("company_id")
	rows, err := h.svc.GetUninvoicedReceipts(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// ─── Requisition Handlers ────────────────────────────────────────────────

func (h *PurchaseHandler) CreateRequisition(c *gin.Context) {
	var req domain.PurchaseRequisition
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.CompanyID = c.Query("company_id")
	if req.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id query param required"})
		return
	}
	req.CreatedBy = GetUserID(c)
	if err := h.svc.CreateRequisition(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *PurchaseHandler) GetRequisition(c *gin.Context) {
	id := c.Param("id")
	req, err := h.svc.GetRequisition(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *PurchaseHandler) ListRequisitions(c *gin.Context) {
	companyID := c.Query("company_id")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	filter := domain.RequisitionFilter{CompanyID: companyID, Limit: limit, Offset: offset}
	if st := c.Query("status"); st != "" {
		filter.Status = domain.RequisitionStatus(st)
	}
	if req := c.Query("requester_id"); req != "" {
		filter.RequesterID = req
	}
	list, total, err := h.svc.ListRequisitions(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "total": total})
}

func (h *PurchaseHandler) UpdateRequisition(c *gin.Context) {
	id := c.Param("id")
	var req domain.PurchaseRequisition
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := h.svc.UpdateRequisition(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *PurchaseHandler) DeleteRequisition(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRequisition(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *PurchaseHandler) SubmitRequisition(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.SubmitRequisition(c.Request.Context(), id, GetUserID(c)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submitted"})
}

func (h *PurchaseHandler) ApproveRequisition(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ApproveRequisition(c.Request.Context(), id, GetUserID(c)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *PurchaseHandler) RejectRequisition(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&body)
	if err := h.svc.RejectRequisition(c.Request.Context(), id, body.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

func (h *PurchaseHandler) ConvertRequisitionToPO(c *gin.Context) {
	id := c.Param("id")
	supplierID := c.Query("supplier_id")
	if supplierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supplier_id query param required"})
		return
	}
	po, err := h.svc.ConvertRequisitionToPO(c.Request.Context(), id, supplierID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, po)
}
