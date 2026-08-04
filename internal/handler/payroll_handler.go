package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type PayrollHandler struct {
	svc *service.PayrollService
}

func NewPayrollHandler(svc *service.PayrollService) *PayrollHandler {
	return &PayrollHandler{svc: svc}
}

func RegisterPayrollRoutes(r *gin.Engine, h *PayrollHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)
	pw := v1.Group("/payroll")
	{
		// Employee payroll info
		pw.GET("/employees/:id", h.GetEmployeePayrollInfo)
		pw.PUT("/employees/:id", h.UpdateEmployeePayrollInfo)
		pw.GET("/employees/:id/dependants", h.ListDependants)
		pw.POST("/employees/:id/dependants", h.CreateDependant)
		pw.DELETE("/employees/:id/dependants/:did", h.DeleteDependant)

		// Periods
		pw.GET("/periods", h.ListPeriods)
		pw.POST("/periods", h.CreatePeriod)
		pw.GET("/periods/:id", h.GetPeriod)
		pw.POST("/periods/:id/submit", h.SubmitPeriod)
		pw.POST("/periods/:id/approve", h.ApprovePeriod)
		pw.GET("/periods/:id/runs", h.ListRuns)
		pw.GET("/periods/:id/summary", h.GetPeriodSummary)

		// Runs
		pw.PUT("/runs/:id", h.UpdateRun)

		// Timekeeping
		pw.POST("/timekeeping", h.CreateTimekeeping)
		pw.GET("/timekeeping", h.ListTimekeeping)
		pw.POST("/timekeeping/bulk", h.BulkCreateTimekeeping)

		// Leave
		pw.POST("/leave", h.RequestLeave)
		pw.GET("/leave/pending", h.ListPendingLeaveRequests)
		pw.POST("/leave/:id/approve", h.ApproveLeaveRequest)
		pw.POST("/leave/:id/reject", h.RejectLeaveRequest)
		pw.GET("/leave/balance", h.GetLeaveBalance)

		// Payslips
		pw.GET("/payslips", h.ListPayslipsByPeriod)
		pw.GET("/payslips/:id", h.GetPayslip)

		// Config
		pw.GET("/config", h.GetConfig)
		pw.PUT("/config", h.SetConfig)
	}
}

func (h *PayrollHandler) payrollError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrPayrollPeriodNotFound),
		errors.Is(err, domain.ErrPayrollRunNotFound),
		errors.Is(err, domain.ErrPayrollEmployeeNotFound),
		errors.Is(err, domain.ErrPayrollDependantNotFound),
		errors.Is(err, domain.ErrPayrollLeaveNotFound),
		errors.Is(err, domain.ErrPayrollPayslipNotFound),
		errors.Is(err, domain.ErrPayrollConfigNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrPayrollPeriodExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrPayrollPeriodNotDraft),
		errors.Is(err, domain.ErrPayrollLeaveAlreadyProcessed):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// ─── Employee Payroll Info ──────────────────────────────────────

func (h *PayrollHandler) GetEmployeePayrollInfo(c *gin.Context) {
	info, err := h.svc.GetEmployeePayrollInfo(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *PayrollHandler) UpdateEmployeePayrollInfo(c *gin.Context) {
	var info domain.EmployeePayrollInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	info.EmployeeID = c.Param("id")
	if err := h.svc.UpdateEmployeePayrollInfo(c.Request.Context(), &info); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *PayrollHandler) ListDependants(c *gin.Context) {
	deps, err := h.svc.ListDependants(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, deps)
}

func (h *PayrollHandler) CreateDependant(c *gin.Context) {
	var dep domain.Dependant
	if err := c.ShouldBindJSON(&dep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dep.EmployeeID = c.Param("id")
	if err := h.svc.CreateDependant(c.Request.Context(), &dep); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dep)
}

func (h *PayrollHandler) DeleteDependant(c *gin.Context) {
	if err := h.svc.DeleteDependant(c.Request.Context(), c.Param("did")); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ─── Periods ────────────────────────────────────────────────────

func (h *PayrollHandler) ListPeriods(c *gin.Context) {
	companyID := c.Query("company_id")
	periods, err := h.svc.ListPeriods(c.Request.Context(), companyID)
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, periods)
}

func (h *PayrollHandler) CreatePeriod(c *gin.Context) {
	var req struct {
		CompanyID string `json:"company_id"`
		Year      int    `json:"year"`
		Month     int    `json:"month"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	period, err := h.svc.CreatePeriod(c.Request.Context(), req.CompanyID, req.Year, req.Month)
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusCreated, period)
}

func (h *PayrollHandler) GetPeriod(c *gin.Context) {
	period, err := h.svc.GetPeriod(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, period)
}

func (h *PayrollHandler) SubmitPeriod(c *gin.Context) {
	// TODO: implement submit for review workflow
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not yet implemented"})
}

func (h *PayrollHandler) ApprovePeriod(c *gin.Context) {
	approvedBy := c.GetString("username")
	if err := h.svc.ApprovePeriod(c.Request.Context(), c.Param("id"), approvedBy); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

// ─── Runs ───────────────────────────────────────────────────────

func (h *PayrollHandler) ListRuns(c *gin.Context) {
	runs, err := h.svc.ListRunsByPeriod(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, runs)
}

func (h *PayrollHandler) UpdateRun(c *gin.Context) {
	var run domain.PayrollRun
	if err := c.ShouldBindJSON(&run); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	run.ID = c.Param("id")
	if err := h.svc.UpdateRun(c.Request.Context(), &run); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, run)
}

// ─── Timekeeping ────────────────────────────────────────────────

func (h *PayrollHandler) CreateTimekeeping(c *gin.Context) {
	var tk domain.TimekeepingRecord
	if err := c.ShouldBindJSON(&tk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateTimekeeping(c.Request.Context(), &tk); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusCreated, tk)
}

func (h *PayrollHandler) ListTimekeeping(c *gin.Context) {
	employeeID := c.Query("employee_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	records, err := h.svc.ListTimekeeping(c.Request.Context(), employeeID, startDate, endDate)
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *PayrollHandler) BulkCreateTimekeeping(c *gin.Context) {
	var records []domain.TimekeepingRecord
	if err := c.ShouldBindJSON(&records); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.BulkCreateTimekeeping(c.Request.Context(), records); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"count": len(records)})
}

// ─── Leave ──────────────────────────────────────────────────────

func (h *PayrollHandler) RequestLeave(c *gin.Context) {
	var lr domain.LeaveRequest
	if err := c.ShouldBindJSON(&lr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RequestLeave(c.Request.Context(), &lr); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusCreated, lr)
}

func (h *PayrollHandler) ListPendingLeaveRequests(c *gin.Context) {
	companyID := c.Query("company_id")
	requests, err := h.svc.ListPendingLeaveRequests(c.Request.Context(), companyID)
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, requests)
}

func (h *PayrollHandler) ApproveLeaveRequest(c *gin.Context) {
	approvedBy := c.GetString("username")
	if err := h.svc.ApproveLeaveRequest(c.Request.Context(), c.Param("id"), approvedBy); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

func (h *PayrollHandler) RejectLeaveRequest(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RejectLeaveRequest(c.Request.Context(), c.Param("id"), req.Reason); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}

func (h *PayrollHandler) GetLeaveBalance(c *gin.Context) {
	employeeID := c.Query("employee_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", "2026"))
	leaveType := domain.LeaveType(c.DefaultQuery("leave_type", "ANNUAL"))
	balance, err := h.svc.GetLeaveBalance(c.Request.Context(), employeeID, year, leaveType)
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, balance)
}

// ─── Payslips ───────────────────────────────────────────────────

func (h *PayrollHandler) ListPayslipsByPeriod(c *gin.Context) {
	periodID := c.Query("period_id")
	payslips, err := h.svc.ListPayslipsByPeriod(c.Request.Context(), periodID)
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, payslips)
}

func (h *PayrollHandler) GetPayslip(c *gin.Context) {
	payslip, err := h.svc.GetPayslipByRun(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, payslip)
}

// ─── Summary ────────────────────────────────────────────────────

func (h *PayrollHandler) GetPeriodSummary(c *gin.Context) {
	summary, err := h.svc.GetPeriodSummary(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, summary)
}

// ─── Config ─────────────────────────────────────────────────────

func (h *PayrollHandler) GetConfig(c *gin.Context) {
	companyID := c.Query("company_id")
	key := c.Query("key")
	cfg, err := h.svc.GetConfig(c.Request.Context(), companyID, key)
	if err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *PayrollHandler) SetConfig(c *gin.Context) {
	var cfg domain.PayrollConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.SetConfig(c.Request.Context(), &cfg); err != nil {
		h.payrollError(c, err)
		return
	}
	c.JSON(http.StatusOK, cfg)
}
