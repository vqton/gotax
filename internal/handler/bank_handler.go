package handler

import (
	"gotax/internal/domain"
	"gotax/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	svc *service.BankService
}

func NewBankHandler(svc *service.BankService) *BankHandler {
	return &BankHandler{svc: svc}
}

func RegisterBankRoutes(r *gin.Engine, h *BankHandler, authMW gin.HandlerFunc) {
	bank := r.Group("/api/v1/bank", authMW)
	{
		statements := bank.Group("/statements")
		{
			statements.POST("/import", h.ImportStatement)
			statements.GET("", h.ListStatements)
			statements.GET("/:id", h.GetStatement)
			statements.DELETE("/:id", h.DeleteStatement)
			statements.GET("/:id/lines", h.GetStatementLines)
		}

		recon := bank.Group("/reconciliations")
		{
			recon.POST("", h.StartReconciliation)
			recon.GET("", h.ListReconciliations)
			recon.GET("/:id", h.GetReconciliation)
			recon.POST("/:id/complete", h.CompleteReconciliation)
			recon.POST("/:id/matches", h.AddMatch)
			recon.GET("/:id/matches", h.GetMatches)
			recon.DELETE("/matches/:id", h.RemoveMatch)
		}

		orders := bank.Group("/payment-orders")
		{
			orders.POST("", h.CreatePaymentOrder)
			orders.GET("", h.ListPaymentOrders)
			orders.GET("/:id", h.GetPaymentOrder)
			orders.POST("/:id/submit", h.SubmitPaymentOrder)
			orders.POST("/:id/approve", h.ApprovePaymentOrder)
			orders.POST("/:id/reject", h.RejectPaymentOrder)
		}

		batches := bank.Group("/batches")
		{
			batches.POST("", h.CreateBatch)
			batches.GET("", h.ListBatches)
			batches.GET("/:id", h.GetBatch)
			batches.POST("/:id/submit", h.SubmitBatch)
			batches.GET("/:id/orders", h.GetBatchOrders)
		}

		loans := bank.Group("/loans")
		{
			loans.POST("", h.CreateLoan)
			loans.GET("", h.ListLoans)
			loans.GET("/:id", h.GetLoan)
			loans.POST("/:id/disburse", h.DisburseLoan)
			loans.GET("/:id/disbursements", h.GetDisbursements)
			loans.POST("/:id/repay", h.MakeRepayment)
			loans.GET("/:id/repayments", h.GetRepayments)
		}

		deposits := bank.Group("/term-deposits")
		{
			deposits.POST("", h.CreateDeposit)
			deposits.GET("", h.ListDeposits)
			deposits.GET("/:id", h.GetDeposit)
			deposits.POST("/:id/mature", h.MatureDeposit)
		}
	}

	reports := r.Group("/api/v1/reports", authMW)
	{
		reports.GET("/bank-ledger", h.GetBankLedger)
		reports.GET("/bank-balance", h.GetBalance)
	}
}

// ─── Statement ──────────────────────────────────────────────────────────

func (h *BankHandler) ImportStatement(c *gin.Context) {
	var req struct {
		Statement *domain.BankStatement    `json:"statement"`
		Lines     []domain.BankStatementLine `json:"lines"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Statement.CompanyID = c.Query("company_id")
	req.Statement.ImportedBy = GetUserID(c)

	if err := h.svc.ImportStatement(c.Request.Context(), req.Statement, req.Lines); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req.Statement)
}

func (h *BankHandler) GetStatement(c *gin.Context) {
	id := c.Param("id")
	s, err := h.svc.GetStatement(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *BankHandler) ListStatements(c *gin.Context) {
	companyID := c.Query("company_id")
	bankAccountID := c.Query("bank_account_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	items, total, err := h.svc.ListStatements(c.Request.Context(), companyID, bankAccountID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (h *BankHandler) DeleteStatement(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteStatement(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BankHandler) GetStatementLines(c *gin.Context) {
	statementID := c.Param("id")
	lines, err := h.svc.GetStatementLines(c.Request.Context(), statementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lines)
}

// ─── Reconciliation ─────────────────────────────────────────────────────

func (h *BankHandler) StartReconciliation(c *gin.Context) {
	var rc domain.BankReconciliation
	if err := c.ShouldBindJSON(&rc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rc.CompanyID = c.Query("company_id")

	if err := h.svc.StartReconciliation(c.Request.Context(), &rc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rc)
}

func (h *BankHandler) GetReconciliation(c *gin.Context) {
	id := c.Param("id")
	rc, err := h.svc.GetReconciliation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rc)
}

func (h *BankHandler) ListReconciliations(c *gin.Context) {
	companyID := c.Query("company_id")
	bankAccountID := c.Query("bank_account_id")
	items, err := h.svc.ListReconciliations(c.Request.Context(), companyID, bankAccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *BankHandler) CompleteReconciliation(c *gin.Context) {
	id := c.Param("id")
	completedBy := GetUserID(c)
	if err := h.svc.CompleteReconciliation(c.Request.Context(), id, completedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "reconciliation completed"})
}

func (h *BankHandler) AddMatch(c *gin.Context) {
	var m domain.BankReconciliationMatch
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddMatch(c.Request.Context(), &m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *BankHandler) GetMatches(c *gin.Context) {
	reconID := c.Param("id")
	items, err := h.svc.GetMatches(c.Request.Context(), reconID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *BankHandler) RemoveMatch(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RemoveMatch(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ─── Payment Orders ─────────────────────────────────────────────────────

func (h *BankHandler) CreatePaymentOrder(c *gin.Context) {
	var po domain.PaymentOrder
	if err := c.ShouldBindJSON(&po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	po.CompanyID = c.Query("company_id")
	po.CreatedBy = GetUserID(c)

	if err := h.svc.CreatePaymentOrder(c.Request.Context(), &po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, po)
}

func (h *BankHandler) GetPaymentOrder(c *gin.Context) {
	id := c.Param("id")
	po, err := h.svc.GetPaymentOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, po)
}

func (h *BankHandler) ListPaymentOrders(c *gin.Context) {
	filter := domain.PaymentOrderFilter{
		CompanyID:   c.Query("company_id"),
		Status:      domain.PaymentOrderStatus(c.Query("status")),
		PaymentType: domain.PaymentOrderType(c.Query("payment_type")),
		FromDate:    c.Query("from_date"),
		ToDate:      c.Query("to_date"),
	}
	filter.Offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	filter.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, total, err := h.svc.ListPaymentOrders(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (h *BankHandler) SubmitPaymentOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.SubmitPaymentOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submitted for approval"})
}

func (h *BankHandler) ApprovePaymentOrder(c *gin.Context) {
	id := c.Param("id")
	approvedBy := GetUserID(c)
	if err := h.svc.ApprovePaymentOrder(c.Request.Context(), id, approvedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *BankHandler) RejectPaymentOrder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RejectPaymentOrder(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

// ─── Payment Batches ────────────────────────────────────────────────────

func (h *BankHandler) CreateBatch(c *gin.Context) {
	var req struct {
		Batch    domain.PaymentOrderBatch `json:"batch"`
		OrderIDs []string                 `json:"order_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Batch.CompanyID = c.Query("company_id")
	req.Batch.CreatedBy = GetUserID(c)

	if err := h.svc.CreateBatch(c.Request.Context(), &req.Batch, req.OrderIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req.Batch)
}

func (h *BankHandler) GetBatch(c *gin.Context) {
	id := c.Param("id")
	b, err := h.svc.GetBatch(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BankHandler) ListBatches(c *gin.Context) {
	companyID := c.Query("company_id")
	items, err := h.svc.ListBatches(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *BankHandler) SubmitBatch(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.SubmitBatch(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "batch submitted"})
}

func (h *BankHandler) GetBatchOrders(c *gin.Context) {
	batchID := c.Param("id")
	orderIDs, err := h.svc.GetBatchOrders(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderIDs)
}

// ─── Loans ─────────────────────────────────────────────────────────────

func (h *BankHandler) CreateLoan(c *gin.Context) {
	var l domain.LoanAgreement
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	l.CompanyID = c.Query("company_id")
	if err := h.svc.CreateLoan(c.Request.Context(), &l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, l)
}

func (h *BankHandler) GetLoan(c *gin.Context) {
	id := c.Param("id")
	l, err := h.svc.GetLoan(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, l)
}

func (h *BankHandler) ListLoans(c *gin.Context) {
	filter := domain.LoanFilter{
		CompanyID: c.Query("company_id"),
		Status:    domain.LoanStatus(c.Query("status")),
		LoanType:  domain.LoanType(c.Query("loan_type")),
	}
	items, err := h.svc.ListLoans(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *BankHandler) DisburseLoan(c *gin.Context) {
	var d domain.LoanDisbursement
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d.LoanID = c.Param("id")
	if err := h.svc.DisburseLoan(c.Request.Context(), &d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *BankHandler) GetDisbursements(c *gin.Context) {
	loanID := c.Param("id")
	items, err := h.svc.GetDisbursements(c.Request.Context(), loanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *BankHandler) MakeRepayment(c *gin.Context) {
	var rp domain.LoanRepayment
	if err := c.ShouldBindJSON(&rp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rp.LoanID = c.Param("id")
	if err := h.svc.MakeRepayment(c.Request.Context(), &rp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rp)
}

func (h *BankHandler) GetRepayments(c *gin.Context) {
	loanID := c.Param("id")
	items, err := h.svc.GetRepayments(c.Request.Context(), loanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ─── Term Deposits ─────────────────────────────────────────────────────

func (h *BankHandler) CreateDeposit(c *gin.Context) {
	var d domain.TermDeposit
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d.CompanyID = c.Query("company_id")
	if err := h.svc.CreateDeposit(c.Request.Context(), &d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *BankHandler) GetDeposit(c *gin.Context) {
	id := c.Param("id")
	d, err := h.svc.GetDeposit(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *BankHandler) ListDeposits(c *gin.Context) {
	companyID := c.Query("company_id")
	items, err := h.svc.ListDeposits(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *BankHandler) MatureDeposit(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.MatureDeposit(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deposit matured"})
}

// ─── Reports ─────────────────────────────────────────────────────────

func (h *BankHandler) GetBankLedger(c *gin.Context) {
	companyID := c.Query("company_id")
	bankAccountID := c.Query("bank_account_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	ledger, err := h.svc.GetBankLedger(c.Request.Context(), companyID, bankAccountID, fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ledger)
}

func (h *BankHandler) GetBalance(c *gin.Context) {
	companyID := c.Query("company_id")
	bankAccountID := c.Query("bank_account_id")

	balance, err := h.svc.GetBalance(c.Request.Context(), companyID, bankAccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
