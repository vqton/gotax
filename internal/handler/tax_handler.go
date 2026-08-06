package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type TaxHandler struct {
	svc   service.TaxServiceInterface
	audit domain.AuditLogRepository
}

func NewTaxHandler(svc service.TaxServiceInterface, audit domain.AuditLogRepository) *TaxHandler {
	return &TaxHandler{svc: svc, audit: audit}
}

func (h *TaxHandler) logAudit(c *gin.Context, action domain.AuditAction, entityType, entityID string) {
	_ = h.audit.Create(c.Request.Context(), &domain.AuditEntry{
		UserID:     GetUserID(c),
		Username:   GetUsername(c),
		IPAddress:  c.ClientIP(),
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
	})
}

func RegisterTaxRoutes(r *gin.Engine, h *TaxHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	tax := v1.Group("/tax")
	{
		declarations := tax.Group("/declarations")
		{
			declarations.POST("", h.CreateDeclaration)
			declarations.POST("/generate", h.GenerateDeclaration)
			declarations.GET("", h.ListDeclarations)
			declarations.GET("/:id", h.GetDeclaration)
			declarations.PUT("/:id", h.UpdateDeclaration)
			declarations.POST("/:id/submit", h.SubmitDeclaration)
			declarations.POST("/:id/check-status", h.CheckDeclarationStatus)
			declarations.POST("/:id/acknowledge", h.AcknowledgeDeclaration)
			declarations.POST("/:id/reject", h.RejectDeclaration)
			declarations.POST("/:id/cancel", h.CancelDeclaration)
			declarations.POST("/:id/amend", h.AmendDeclaration)
		}
		rates := tax.Group("/rates")
		{
			rates.POST("", h.CreateRate)
			rates.GET("", h.ListRates)
			rates.GET("/:id", h.GetRate)
			rates.PUT("/:id", h.UpdateRate)
		}
	payments := tax.Group("/payments")
	{
		payments.POST("", h.CreatePayment)
		payments.GET("", h.ListPayments)
		payments.GET("/summary", h.GetPaymentSummary)
		payments.GET("/:id", h.GetPayment)
		payments.POST("/:id/record", h.RecordPayment)
	}
		invoices := tax.Group("/e-invoices")
		{
			invoices.POST("", h.CreateEInvoice)
			invoices.GET("", h.ListEInvoices)
			invoices.GET("/:id", h.GetEInvoice)
			invoices.POST("/:id/issue", h.IssueEInvoice)
			invoices.POST("/:id/status", h.CheckInvoiceStatus)
			invoices.POST("/:id/cancel", h.CancelEInvoice)
			invoices.POST("/:id/amend", h.CreateAmendmentInvoice)
		}
		calendar := tax.Group("/calendar")
		{
			calendar.POST("", h.CreateCalendarEntry)
			calendar.GET("/company/:companyID", h.GetCalendarByCompany)
			calendar.GET("/period/:companyID/:year/:number", h.GetCalendarByPeriod)
			calendar.GET("/:id", h.GetCalendarEntry)
		calendar.POST("/scan-overdue", h.ScanOverdueCalendars)
		calendar.POST("/generate-alerts", h.GenerateDeadlineAlerts)
	}
		alerts := tax.Group("/alerts")
		{
			alerts.POST("", h.CreateAlert)
			alerts.GET("", h.ListAlerts)
			alerts.GET("/:id", h.GetAlert)
		}
		audit := tax.Group("/audit-cases")
		{
			audit.POST("", h.CreateAuditCase)
			audit.GET("", h.ListAuditCases)
			audit.GET("/:id", h.GetAuditCase)
			audit.POST("/:id/close", h.CloseAuditCase)
		}
		calc := tax.Group("/calculate")
		{
			calc.POST("/vat", h.CalculateVAT)
			calc.POST("/cit", h.CalculateCIT)
			calc.POST("/cit/provisional", h.CalculateQuarterlyProvisional)
			calc.POST("/pit", h.CalculatePIT)
		}
	reconcile := tax.Group("/reconcile")
	{
		reconcile.POST("/vat", h.ReconcileVAT)
		reconcile.POST("/cit", h.ReconcileCIT)
	}
	penalty := tax.Group("/penalty")
	{
		penalty.POST("/calculate", h.CalculatePenalty)
	}
	}
}

func (h *TaxHandler) taxError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrDeclarationNotFound),
		errors.Is(err, domain.ErrTaxRateNotFound),
		errors.Is(err, domain.ErrTaxPaymentNotFound),
		errors.Is(err, domain.ErrInvoiceNotFound),
		errors.Is(err, domain.ErrCalendarNotFound),
		errors.Is(err, domain.ErrTaxAlertNotFound),
		errors.Is(err, domain.ErrAuditCaseNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrDeclarationNotEditable),
		errors.Is(err, domain.ErrDeclarationAlreadySubmitted),
		errors.Is(err, domain.ErrDeclarationPeriodAlreadyDeclared),
		errors.Is(err, domain.ErrInvoiceStatusInvalid),
		errors.Is(err, domain.ErrDuplicateDeclaration):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrGDTRejected),
		errors.Is(err, domain.ErrGDTInvalidTaxCode):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrGDTUnauthorized),
		errors.Is(err, domain.ErrGDTUnavailable):
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// ─── Declarations ──────────────────────────────────────────────────────

func (h *TaxHandler) CreateDeclaration(c *gin.Context) {
	var d domain.TaxDeclaration
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	userID := c.GetString("user_id")
	d.CreatedBy = userID
	if err := h.svc.CreateDeclaration(c.Request.Context(), &d); err != nil {
		h.taxError(c, err)
		return
	}
	h.logAudit(c, domain.AuditActionCreate, "tax_declaration", d.ID)
	c.JSON(http.StatusCreated, d)
}

func (h *TaxHandler) GenerateDeclaration(c *gin.Context) {
	var req struct {
		CompanyID       string                 `json:"company_id"`
		DeclarationType domain.DeclarationType `json:"declaration_type"`
		TaxPeriod       domain.TaxPeriod       `json:"tax_period"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	decl, err := h.svc.GenerateDeclaration(c.Request.Context(), req.CompanyID,
		req.DeclarationType, req.TaxPeriod, c.GetString("user_id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, decl)
}

func (h *TaxHandler) GetDeclaration(c *gin.Context) {
	d, err := h.svc.GetDeclaration(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *TaxHandler) ListDeclarations(c *gin.Context) {
	filter := domain.TaxDeclarationFilter{
		CompanyID:       c.Query("company_id"),
		DeclarationType: domain.DeclarationType(c.Query("declaration_type")),
		Status:          domain.DeclarationStatus(c.Query("status")),
	}
	if year, err := strconv.Atoi(c.Query("period_year")); err == nil {
		filter.PeriodYear = year
	}
	if num, err := strconv.Atoi(c.Query("period_number")); err == nil {
		filter.PeriodNumber = num
	}
	declarations, err := h.svc.ListDeclarations(c.Request.Context(), filter)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, declarations)
}

func (h *TaxHandler) UpdateDeclaration(c *gin.Context) {
	var d domain.TaxDeclaration
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	d.ID = c.Param("id")
	if err := h.svc.CreateDeclaration(c.Request.Context(), &d); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *TaxHandler) SubmitDeclaration(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.SubmitDeclaration(c.Request.Context(), id, userID); err != nil {
		h.taxError(c, err)
		return
	}
	h.logAudit(c, domain.AuditActionPost, "tax_declaration", id)
	c.JSON(http.StatusOK, gin.H{"message": "declaration submitted to GDT"})
}

func (h *TaxHandler) CheckDeclarationStatus(c *gin.Context) {
	if err := h.svc.CheckDeclarationStatus(c.Request.Context(), c.Param("id")); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "declaration status updated"})
}

func (h *TaxHandler) ReconcileVAT(c *gin.Context) {
	var req struct {
		CompanyID string           `json:"company_id"`
		Period    domain.TaxPeriod `json:"period"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	res, err := h.svc.ReconcileVAT(c.Request.Context(), req.CompanyID, req.Period)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *TaxHandler) ReconcileCIT(c *gin.Context) {
	var req struct {
		CompanyID string           `json:"company_id"`
		Period    domain.TaxPeriod `json:"period"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	res, err := h.svc.ReconcileCIT(c.Request.Context(), req.CompanyID, req.Period)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *TaxHandler) CalculatePenalty(c *gin.Context) {
	var req struct {
		DeclarationDue string `json:"declaration_due"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	res, err := h.svc.CalculatePenalty(c.Request.Context(), req.DeclarationDue)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *TaxHandler) AcknowledgeDeclaration(c *gin.Context) {
	var req struct {
		Reference string `json:"reference"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reference required"})
		return
	}
	if err := h.svc.AcknowledgeDeclaration(c.Request.Context(), c.Param("id"), req.Reference); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "declaration acknowledged"})
}

func (h *TaxHandler) RejectDeclaration(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason required"})
		return
	}
	if err := h.svc.RejectDeclaration(c.Request.Context(), c.Param("id"), req.Reason); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "declaration rejected"})
}

func (h *TaxHandler) CancelDeclaration(c *gin.Context) {
	if err := h.svc.CancelDeclaration(c.Request.Context(), c.Param("id")); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "declaration cancelled"})
}

func (h *TaxHandler) AmendDeclaration(c *gin.Context) {
	var req struct {
		Lines []domain.TaxDeclarationLine `json:"lines"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	amended, err := h.svc.AmendDeclaration(c.Request.Context(), c.Param("id"), req.Lines)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, amended)
}

// ─── Rates ─────────────────────────────────────────────────────────────

func (h *TaxHandler) CreateRate(c *gin.Context) {
	var rate domain.TaxRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateRate(c.Request.Context(), &rate); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, rate)
}

func (h *TaxHandler) GetRate(c *gin.Context) {
	rate, err := h.svc.GetRate(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, rate)
}

func (h *TaxHandler) ListRates(c *gin.Context) {
	filter := domain.TaxRateFilter{
		TaxType: domain.TaxType(c.Query("tax_type")),
	}
	rates, err := h.svc.ListRates(c.Request.Context(), filter)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, rates)
}

func (h *TaxHandler) UpdateRate(c *gin.Context) {
	var rate domain.TaxRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	rate.ID = c.Param("id")
	if err := h.svc.UpdateRate(c.Request.Context(), &rate); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, rate)
}

// ─── Payments ──────────────────────────────────────────────────────────

func (h *TaxHandler) CreatePayment(c *gin.Context) {
	var p domain.TaxPayment
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreatePayment(c.Request.Context(), &p); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *TaxHandler) GetPayment(c *gin.Context) {
	p, err := h.svc.GetPayment(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *TaxHandler) ListPayments(c *gin.Context) {
	filter := domain.PaymentFilter{
		CompanyID: c.Query("company_id"),
		TaxType:   domain.TaxType(c.Query("tax_type")),
		Status:    domain.PaymentStatus(c.Query("status")),
	}
	if year, err := strconv.Atoi(c.Query("period_year")); err == nil {
		filter.PeriodYear = year
	}
	if num, err := strconv.Atoi(c.Query("period_number")); err == nil {
		filter.PeriodNumber = num
	}
	payments, err := h.svc.ListPayments(c.Request.Context(), filter)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, payments)
}

func (h *TaxHandler) GetPaymentSummary(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}
	summary, err := h.svc.GetPaymentSummary(c.Request.Context(), companyID)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *TaxHandler) RecordPayment(c *gin.Context) {
	var req struct {
		Amount float64 `json:"amount"`
		Date   string  `json:"date"`
		Ref    string  `json:"reference"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	id := c.Param("id")
	if err := h.svc.RecordPayment(c.Request.Context(), id, req.Amount, req.Date, req.Ref); err != nil {
		h.taxError(c, err)
		return
	}
	h.logAudit(c, domain.AuditActionPost, "tax_payment", id)
	c.JSON(http.StatusOK, gin.H{"message": "payment recorded"})
}

// ─── E-Invoices ────────────────────────────────────────────────────────

func (h *TaxHandler) CreateEInvoice(c *gin.Context) {
	var req domain.EInvoiceInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	inv := &domain.EInvoice{
		CompanyID:          req.CompanyID,
		Pattern:            req.Pattern,
		Serial:             req.Serial,
		InvoiceType:        req.InvoiceType,
		BuyerName:          req.BuyerName,
		BuyerTaxCode:       req.BuyerTaxCode,
		BuyerAddress:       req.BuyerAddress,
		BuyerEmail:         req.BuyerEmail,
		CurrencyCode:       req.CurrencyCode,
		ExchangeRate:       req.ExchangeRate,
		Lines:              req.Lines,
		DigitalSignatureID: req.DigitalSignatureID,
		IssueDate:          req.IssueDate,
		Status:             domain.EInvStatusDRAFT,
	}
	for i := range inv.Lines {
		inv.Lines[i].LineNumber = i + 1
		inv.Lines[i].VATAmount = inv.Lines[i].LineTotal * inv.Lines[i].VATRate / 100.0
		inv.Subtotal += inv.Lines[i].LineTotal
		inv.VATAmount += inv.Lines[i].VATAmount
	}
	inv.GrandTotal = inv.Subtotal + inv.VATAmount
	if err := h.svc.CreateEInvoice(c.Request.Context(), inv); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, inv)
}

func (h *TaxHandler) GetEInvoice(c *gin.Context) {
	inv, err := h.svc.GetEInvoice(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *TaxHandler) ListEInvoices(c *gin.Context) {
	filter := domain.EInvoiceFilter{
		CompanyID: c.Query("company_id"),
		Status:    domain.EInvLifecycleStatus(c.Query("status")),
		FromDate:  c.Query("from_date"),
		ToDate:    c.Query("to_date"),
	}
	invoices, err := h.svc.ListEInvoices(c.Request.Context(), filter)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, invoices)
}

func (h *TaxHandler) IssueEInvoice(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.IssueEInvoice(c.Request.Context(), id); err != nil {
		h.taxError(c, err)
		return
	}
	h.logAudit(c, domain.AuditActionPost, "e_invoice", id)
	c.JSON(http.StatusOK, gin.H{"message": "invoice submitted to GDT"})
}

func (h *TaxHandler) CheckInvoiceStatus(c *gin.Context) {
	if err := h.svc.CheckInvoiceStatus(c.Request.Context(), c.Param("id")); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "invoice status updated"})
}

func (h *TaxHandler) CancelEInvoice(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason required"})
		return
	}
	id := c.Param("id")
	if err := h.svc.CancelEInvoice(c.Request.Context(), id, req.Reason); err != nil {
		h.taxError(c, err)
		return
	}
	h.logAudit(c, domain.AuditActionCancel, "e_invoice", id)
	c.JSON(http.StatusOK, gin.H{"message": "invoice cancelled"})
}

func (h *TaxHandler) CreateAmendmentInvoice(c *gin.Context) {
	var req struct {
		InvoiceType string               `json:"invoice_type" binding:"required"`
		Lines       []domain.EInvoiceLine `json:"lines" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	invType := domain.EInvoiceType(req.InvoiceType)
	inv, err := h.svc.CreateAmendmentInvoice(c.Request.Context(), c.Param("id"), invType, req.Lines, GetUserID(c))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, inv)
}

// ─── Calendar ──────────────────────────────────────────────────────────

func (h *TaxHandler) CreateCalendarEntry(c *gin.Context) {
	var cal domain.TaxCalendar
	if err := c.ShouldBindJSON(&cal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateCalendarEntry(c.Request.Context(), &cal); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, cal)
}

func (h *TaxHandler) GetCalendarEntry(c *gin.Context) {
	cal, err := h.svc.GetCalendarEntry(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, cal)
}

func (h *TaxHandler) GetCalendarByPeriod(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))
	number, _ := strconv.Atoi(c.Param("number"))
	entries, err := h.svc.GetCalendarByPeriod(c.Request.Context(), c.Param("companyID"), year, number)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *TaxHandler) GetCalendarByCompany(c *gin.Context) {
	entries, err := h.svc.GetCalendarByCompany(c.Request.Context(), c.Param("companyID"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *TaxHandler) ScanOverdueCalendars(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}
	count, err := h.svc.ScanOverdueCalendars(c.Request.Context(), companyID)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"overdue_count": count})
}

func (h *TaxHandler) GenerateDeadlineAlerts(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}
	daysAhead := 7
	if d := c.Query("days"); d != "" {
		if v, err := strconv.Atoi(d); err == nil && v > 0 {
			daysAhead = v
		}
	}
	count, err := h.svc.GenerateDeadlineAlerts(c.Request.Context(), companyID, daysAhead)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"alerts_created": count})
}

// ─── Alerts ────────────────────────────────────────────────────────────

func (h *TaxHandler) CreateAlert(c *gin.Context) {
	var a domain.TaxAlert
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAlert(c.Request.Context(), &a); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *TaxHandler) GetAlert(c *gin.Context) {
	a, err := h.svc.GetAlert(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *TaxHandler) ListAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	alerts, err := h.svc.ListAlerts(c.Request.Context(), c.Query("company_id"), limit)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, alerts)
}

// ─── Audit Cases ───────────────────────────────────────────────────────

func (h *TaxHandler) CreateAuditCase(c *gin.Context) {
	var a domain.TaxAuditCase
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateAuditCase(c.Request.Context(), &a); err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *TaxHandler) GetAuditCase(c *gin.Context) {
	a, err := h.svc.GetAuditCase(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *TaxHandler) ListAuditCases(c *gin.Context) {
	cases, err := h.svc.ListAuditCases(c.Request.Context(), c.Query("company_id"))
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, cases)
}

func (h *TaxHandler) CloseAuditCase(c *gin.Context) {
	var req struct {
		Findings string  `json:"findings"`
		Penalty  float64 `json:"penalty_amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	id := c.Param("id")
	if err := h.svc.CloseAuditCase(c.Request.Context(), id, req.Findings, req.Penalty); err != nil {
		h.taxError(c, err)
		return
	}
	h.logAudit(c, domain.AuditActionClose, "tax_audit_case", id)
	c.JSON(http.StatusOK, gin.H{"message": "audit case closed"})
}

// ─── Calculations ──────────────────────────────────────────────────────

func (h *TaxHandler) CalculateVAT(c *gin.Context) {
	var req struct {
		CompanyID string                `json:"company_id"`
		Period    domain.TaxPeriod      `json:"period"`
		Entries   []domain.JournalEntry `json:"entries"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	result, err := h.svc.CalculateVAT(c.Request.Context(), req.CompanyID, req.Period, req.Entries)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TaxHandler) CalculateCIT(c *gin.Context) {
	var req struct {
		CompanyID string                `json:"company_id"`
		Year      int                   `json:"year"`
		Entries   []domain.JournalEntry `json:"entries"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	result, err := h.svc.CalculateCIT(c.Request.Context(), req.CompanyID, req.Year, req.Entries)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TaxHandler) CalculatePIT(c *gin.Context) {
	var req struct {
		CompanyID string                    `json:"company_id"`
		Period    domain.TaxPeriod          `json:"period"`
		Employees []domain.PITEmployeeInput `json:"employees"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	result, err := h.svc.CalculatePIT(c.Request.Context(), req.CompanyID, req.Period, req.Employees)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TaxHandler) CalculateQuarterlyProvisional(c *gin.Context) {
	var req struct {
		CompanyID string                `json:"company_id"`
		Year      int                   `json:"year"`
		Quarter   int                   `json:"quarter"`
		Entries   []domain.JournalEntry `json:"entries"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	result, err := h.svc.CalculateQuarterlyProvisional(c.Request.Context(), req.CompanyID, req.Year, req.Quarter, req.Entries)
	if err != nil {
		h.taxError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
