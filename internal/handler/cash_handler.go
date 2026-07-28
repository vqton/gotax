package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type CashHandler struct {
	svc service.Service
}

func NewCashHandler(svc service.Service) *CashHandler {
	return &CashHandler{svc: svc}
}

func RegisterCashRoutes(r *gin.Engine, h *CashHandler, authMW gin.HandlerFunc) {
	cash := r.Group("/api/v1/cash", authMW)
	{
		receipts := cash.Group("/receipts")
		{
			receipts.POST("", h.CreateCashReceipt)
			receipts.GET("", h.ListCashReceipts)
			receipts.GET("/:id", h.GetCashReceipt)
			receipts.PUT("/:id", h.UpdateCashReceipt)
			receipts.DELETE("/:id", h.DeleteCashReceipt)
			receipts.POST("/:id/submit", h.SubmitCashReceipt)
			receipts.POST("/:id/approve", h.ApproveCashReceipt)
			receipts.POST("/:id/reject", h.RejectCashReceipt)
			receipts.POST("/:id/post", h.PostCashReceipt)
			receipts.GET("/:id/print", h.PrintCashReceipt)
		}

		payments := cash.Group("/payments")
		{
			payments.POST("", h.CreateCashPayment)
			payments.GET("", h.ListCashPayments)
			payments.GET("/:id", h.GetCashPayment)
			payments.PUT("/:id", h.UpdateCashPayment)
			payments.DELETE("/:id", h.DeleteCashPayment)
			payments.POST("/:id/submit", h.SubmitCashPayment)
			payments.POST("/:id/approve", h.ApproveCashPayment)
			payments.POST("/:id/reject", h.RejectCashPayment)
			payments.POST("/:id/post", h.PostCashPayment)
			payments.GET("/:id/print", h.PrintCashPayment)
		}

		transfers := cash.Group("/transfers")
		{
			transfers.POST("", h.CreateCashTransfer)
			transfers.GET("", h.ListCashTransfers)
		}

		cash.GET("/cash-book", h.GetCashBook)
		cash.GET("/balance", h.GetCashBalance)

		petty := cash.Group("/petty-cash")
		{
			petty.POST("", h.CreatePettyCashFund)
			petty.GET("", h.ListPettyCashFunds)
			petty.GET("/:id", h.GetPettyCashFund)
		}

		inventory := cash.Group("/inventory")
		{
			inventory.POST("", h.CreateCashInventory)
			inventory.GET("", h.ListCashInventories)
			inventory.GET("/:id", h.GetCashInventory)
		}

		advances := cash.Group("/advances")
		{
			advances.POST("", h.CreateAdvance)
			advances.GET("", h.ListAdvances)
			advances.GET("/:id", h.GetAdvance)
			advances.PUT("/:id", h.UpdateAdvance)
			advances.POST("/:id/approve", h.ApproveAdvance)
			advances.POST("/:id/reject", h.RejectAdvance)
			advances.POST("/:id/pay", h.PayAdvance)
			advances.POST("/:id/settle", h.SettleAdvance)
		}
	}

	reports := r.Group("/api/v1/reports", authMW)
	{
		reports.GET("/cash-flow", h.CashFlowStatement)
		reports.GET("/cash-book", h.CashBookReport)
	}
}

// ─── Cash Receipts ──────────────────────────────────────────────────

func (ch *CashHandler) CreateCashReceipt(c *gin.Context) {
	var req domain.CashReceipt
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.CreatedBy = GetUserID(c)
	if err := ch.svc.CreateCashReceipt(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (ch *CashHandler) GetCashReceipt(c *gin.Context) {
	id := c.Param("id")
	receipt, err := ch.svc.GetCashReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, receipt)
}

func (ch *CashHandler) ListCashReceipts(c *gin.Context) {
	filter := domain.CashReceiptFilter{
		CompanyID:   c.Query("company_id"),
		ReceiptType: domain.ReceiptType(c.Query("receipt_type")),
		Currency:    c.Query("currency"),
		Status:      domain.CashStatus(c.Query("status")),
		FromDate:    c.Query("from_date"),
		ToDate:      c.Query("to_date"),
	}
	filter.Offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	filter.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	if filter.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

	receipts, total, err := ch.svc.ListCashReceipts(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": receipts, "total": total})
}

func (ch *CashHandler) UpdateCashReceipt(c *gin.Context) {
	id := c.Param("id")
	var req domain.CashReceipt
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := ch.svc.UpdateCashReceipt(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (ch *CashHandler) DeleteCashReceipt(c *gin.Context) {
	id := c.Param("id")
	if err := ch.svc.DeleteCashReceipt(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "receipt deleted"})
}

func (ch *CashHandler) SubmitCashReceipt(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.SubmitCashReceipt(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "receipt submitted"})
}

func (ch *CashHandler) ApproveCashReceipt(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.ApproveCashReceipt(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "receipt approved"})
}

func (ch *CashHandler) RejectCashReceipt(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.RejectCashReceipt(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "receipt rejected"})
}

func (ch *CashHandler) PostCashReceipt(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.PostCashReceipt(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "receipt posted"})
}

// ─── Cash Payments ──────────────────────────────────────────────────

func (ch *CashHandler) CreateCashPayment(c *gin.Context) {
	var req domain.CashPayment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.CreatedBy = GetUserID(c)
	if err := ch.svc.CreateCashPayment(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (ch *CashHandler) GetCashPayment(c *gin.Context) {
	id := c.Param("id")
	payment, err := ch.svc.GetCashPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (ch *CashHandler) ListCashPayments(c *gin.Context) {
	filter := domain.CashPaymentFilter{
		CompanyID:   c.Query("company_id"),
		PaymentType: domain.PaymentType(c.Query("payment_type")),
		Currency:    c.Query("currency"),
		Status:      domain.CashStatus(c.Query("status")),
		FromDate:    c.Query("from_date"),
		ToDate:      c.Query("to_date"),
	}
	filter.Offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	filter.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	if filter.CompanyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

	payments, total, err := ch.svc.ListCashPayments(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": payments, "total": total})
}

func (ch *CashHandler) UpdateCashPayment(c *gin.Context) {
	id := c.Param("id")
	var req domain.CashPayment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := ch.svc.UpdateCashPayment(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (ch *CashHandler) DeleteCashPayment(c *gin.Context) {
	id := c.Param("id")
	if err := ch.svc.DeleteCashPayment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment deleted"})
}

func (ch *CashHandler) SubmitCashPayment(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.SubmitCashPayment(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment submitted"})
}

func (ch *CashHandler) ApproveCashPayment(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.ApproveCashPayment(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment approved"})
}

func (ch *CashHandler) RejectCashPayment(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.RejectCashPayment(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment rejected"})
}

func (ch *CashHandler) PostCashPayment(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.PostCashPayment(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment posted"})
}

// ─── Cash Transfers ─────────────────────────────────────────────────

func (ch *CashHandler) CreateCashTransfer(c *gin.Context) {
	var req domain.CashTransfer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ch.svc.CreateCashTransfer(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (ch *CashHandler) ListCashTransfers(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	transfers, err := ch.svc.GetCashTransfers(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transfers)
}

// ─── Cash Book & Balance ────────────────────────────────────────────

func (ch *CashHandler) GetCashBook(c *gin.Context) {
	companyID := c.Query("company_id")
	currency := c.DefaultQuery("currency", "VND")
	accountID := c.Query("account_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if companyID == "" || accountID == "" || fromDate == "" || toDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id, account_id, from_date, to_date are required"})
		return
	}

	cb, err := ch.svc.GetCashBook(c.Request.Context(), companyID, currency, accountID, fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cb)
}

func (ch *CashHandler) GetCashBalance(c *gin.Context) {
	companyID := c.Query("company_id")
	accountID := c.Query("account_id")

	if companyID == "" || accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id and account_id are required"})
		return
	}

	balance, err := ch.svc.GetCashBalance(c.Request.Context(), companyID, accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// ─── Petty Cash ─────────────────────────────────────────────────────

func (ch *CashHandler) CreatePettyCashFund(c *gin.Context) {
	var req domain.PettyCashFund
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ch.svc.CreatePettyCashFund(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (ch *CashHandler) GetPettyCashFund(c *gin.Context) {
	id := c.Param("id")
	fund, err := ch.svc.GetPettyCashFund(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fund)
}

func (ch *CashHandler) ListPettyCashFunds(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	funds, err := ch.svc.ListPettyCashFunds(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, funds)
}

// ─── Cash Inventory ─────────────────────────────────────────────────

func (ch *CashHandler) CreateCashInventory(c *gin.Context) {
	var req domain.CashInventory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ch.svc.CreateCashInventory(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (ch *CashHandler) GetCashInventory(c *gin.Context) {
	id := c.Param("id")
	inv, err := ch.svc.GetCashInventory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (ch *CashHandler) ListCashInventories(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	invs, err := ch.svc.ListCashInventories(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, invs)
}

// ─── Cash Reports ───────────────────────────────────────────────────

func (ch *CashHandler) CashFlowStatement(c *gin.Context) {
	companyID := c.Query("company_id")
	currency := c.DefaultQuery("currency", "VND")
	accountID := c.Query("account_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if companyID == "" || accountID == "" || fromDate == "" || toDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id, account_id, from_date, to_date are required"})
		return
	}

	stmt, err := ch.svc.GetCashFlowStatement(c.Request.Context(), companyID, currency, accountID, fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report_type": "b03_dn_cash_flow", "report": stmt})
}

func (ch *CashHandler) CashBookReport(c *gin.Context) {
	companyID := c.Query("company_id")
	currency := c.DefaultQuery("currency", "VND")
	accountID := c.Query("account_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if companyID == "" || accountID == "" || fromDate == "" || toDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id, account_id, from_date, to_date are required"})
		return
	}

	cb, err := ch.svc.GetCashBook(c.Request.Context(), companyID, currency, accountID, fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cb)
}

// ─── Voucher Forms ──────────────────────────────────────────────────

func (ch *CashHandler) PrintCashReceipt(c *gin.Context) {
	id := c.Param("id")
	receipt, err := ch.svc.GetCashReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"form_type": "M01-TT", "voucher": receipt})
}

func (ch *CashHandler) PrintCashPayment(c *gin.Context) {
	id := c.Param("id")
	payment, err := ch.svc.GetCashPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"form_type": "M02-TT", "voucher": payment})
}

// ─── Advance Request / Settlement ─────────────────────────────────────

func (ch *CashHandler) CreateAdvance(c *gin.Context) {
	var req domain.AdvanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.RequestorID = GetUserID(c)
	req.CreatedBy = GetUserID(c)
	if err := ch.svc.CreateAdvance(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (ch *CashHandler) GetAdvance(c *gin.Context) {
	id := c.Param("id")
	a, err := ch.svc.GetAdvance(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (ch *CashHandler) ListAdvances(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}
	advances, err := ch.svc.ListAdvances(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, advances)
}

func (ch *CashHandler) UpdateAdvance(c *gin.Context) {
	id := c.Param("id")
	var req domain.AdvanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := ch.svc.UpdateAdvance(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (ch *CashHandler) ApproveAdvance(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.ApproveAdvance(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "advance approved"})
}

func (ch *CashHandler) RejectAdvance(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.RejectAdvance(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "advance rejected"})
}

func (ch *CashHandler) PayAdvance(c *gin.Context) {
	id := c.Param("id")
	userID := GetUserID(c)
	if err := ch.svc.PayAdvance(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "advance paid"})
}

func (ch *CashHandler) SettleAdvance(c *gin.Context) {
	id := c.Param("id")
	// settlement info can come in request body
	var req struct {
		SettlementID string `json:"settlement_id"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.SettlementID != "" {
		if err := ch.svc.SettleAdvance(c.Request.Context(), id, req.SettlementID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := ch.svc.SettleAdvance(c.Request.Context(), id, ""); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "advance settled"})
}
