package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type CompanyHandler struct {
	svc service.CompanyService
}

func NewCompanyHandler(svc service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

func RegisterCompanyRoutes(r *gin.Engine, h *CompanyHandler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)

	companies := v1.Group("/companies")
	{
		companies.POST("", h.CreateCompany)
		companies.GET("", h.ListCompanies)
		companies.GET("/:companyID", h.GetCompany)
		companies.GET("/tax-code/:taxCode", h.GetCompanyByTaxCode)
		companies.PUT("/:companyID", h.UpdateCompany)
		companies.DELETE("/:companyID", h.DeactivateCompany)
		companies.GET("/:companyID/hierarchy", h.GetCompanyHierarchy)
		companies.POST("/:companyID/switch", h.SwitchCompany)

		branches := companies.Group("/:companyID/branches")
		{
			branches.POST("", h.CreateBranch)
			branches.GET("", h.ListBranches)
		}

		fiscalYears := companies.Group("/:companyID/fiscal-years")
		{
			fiscalYears.POST("", h.CreateFiscalYear)
			fiscalYears.GET("", h.ListFiscalYears)
		}

		periods := companies.Group("/:companyID/periods")
		{
			periods.POST("/generate", h.GeneratePeriods)
			periods.POST("/:periodID/close", h.ClosePeriod)
			periods.POST("/:periodID/reopen", h.ReopenPeriod)
			periods.POST("/:periodID/permanent-close", h.PermanentClosePeriod)
			periods.GET("/current", h.GetCurrentPeriod)
		}

		departments := companies.Group("/:companyID/departments")
		{
			departments.POST("", h.CreateDepartment)
			departments.GET("", h.ListDepartments)
		}

		employees := companies.Group("/:companyID/employees")
		{
			employees.POST("", h.CreateEmployee)
			employees.GET("", h.ListEmployees)
		}

		bankAccounts := companies.Group("/:companyID/bank-accounts")
		{
			bankAccounts.POST("", h.CreateBankAccount)
			bankAccounts.GET("", h.ListBankAccounts)
		}

		einvoicePatterns := companies.Group("/:companyID/einvoice-patterns")
		{
			einvoicePatterns.POST("", h.RegisterEInvoicePattern)
			einvoicePatterns.GET("", h.ListEInvoicePatterns)
		}

		digitalSignatures := companies.Group("/:companyID/digital-signatures")
		{
			digitalSignatures.POST("", h.RegisterDigitalSignature)
			digitalSignatures.GET("", h.ListDigitalSignatures)
		}

		integrations := companies.Group("/:companyID/integrations")
		{
			integrations.POST("", h.CreateIntegrationProfile)
			integrations.GET("", h.ListIntegrationProfiles)
		}
	}

	get := v1.Group("")
	{
		get.GET("/branches/:id", h.GetBranch)
		get.PUT("/branches/:id", h.UpdateBranch)
		get.DELETE("/branches/:id", h.DeactivateBranch)

		get.GET("/fiscal-years/:id", h.GetFiscalYear)

		get.GET("/departments/:id", h.GetDepartment)
		get.PUT("/departments/:id", h.UpdateDepartment)

		get.GET("/employees/:id", h.GetEmployee)
		get.PUT("/employees/:id", h.UpdateEmployee)
		get.DELETE("/employees/:id", h.DeactivateEmployee)

		get.GET("/bank-accounts/:id", h.GetBankAccount)
		get.PUT("/bank-accounts/:id", h.UpdateBankAccount)
		get.DELETE("/bank-accounts/:id", h.DeactivateBankAccount)

		get.GET("/einvoice-patterns/:id", h.GetEInvoicePattern)

		get.GET("/digital-signatures/:id", h.GetDigitalSignature)

		get.GET("/integrations/:id", h.GetIntegrationProfile)
		get.PUT("/integrations/:id", h.UpdateIntegrationProfile)
		get.POST("/integrations/:id/test", h.TestIntegration)
	}
}

// ─── Company ────────────────────────────────────────────────────────

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var comp domain.Company
	if err := c.ShouldBindJSON(&comp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.CreateCompany(c.Request.Context(), &comp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comp)
}

func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	tenantID := GetUserID(c)
	companies, err := h.svc.ListCompanies(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, companies)
}

func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id := c.Param("companyID")
	company, err := h.svc.GetCompany(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) GetCompanyByTaxCode(c *gin.Context) {
	taxCode := c.Param("taxCode")
	tenantID := GetUserID(c)
	company, err := h.svc.GetCompanyByTaxCode(c.Request.Context(), tenantID, taxCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id := c.Param("companyID")
	var comp domain.Company
	if err := c.ShouldBindJSON(&comp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	comp.ID = id
	if err := h.svc.UpdateCompany(c.Request.Context(), &comp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comp)
}

func (h *CompanyHandler) DeactivateCompany(c *gin.Context) {
	id := c.Param("companyID")
	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)
	if err := h.svc.DeactivateCompany(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "company deactivated"})
}

func (h *CompanyHandler) GetCompanyHierarchy(c *gin.Context) {
	id := c.Param("companyID")
	hierarchy, err := h.svc.GetCompanyHierarchy(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hierarchy)
}

func (h *CompanyHandler) SwitchCompany(c *gin.Context) {
	companyID := c.Param("companyID")
	userID := GetUserID(c)
	ctxOut, err := h.svc.SwitchCompany(c.Request.Context(), userID, companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ctxOut)
}

// ─── Branch ────────────────────────────────────────────────────────

func (h *CompanyHandler) CreateBranch(c *gin.Context) {
	companyID := c.Param("companyID")
	var branch domain.CompanyBranch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	branch.CompanyID = companyID
	if err := h.svc.CreateBranch(c.Request.Context(), &branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, branch)
}

func (h *CompanyHandler) GetBranch(c *gin.Context) {
	id := c.Param("id")
	branch, err := h.svc.GetBranch(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, branch)
}

func (h *CompanyHandler) ListBranches(c *gin.Context) {
	companyID := c.Param("companyID")
	branches, err := h.svc.ListBranches(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, branches)
}

func (h *CompanyHandler) UpdateBranch(c *gin.Context) {
	id := c.Param("id")
	var branch domain.CompanyBranch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	branch.ID = id
	if err := h.svc.UpdateBranch(c.Request.Context(), &branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, branch)
}

func (h *CompanyHandler) DeactivateBranch(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeactivateBranch(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "branch deactivated"})
}

// ─── Fiscal Year ───────────────────────────────────────────────────

func (h *CompanyHandler) CreateFiscalYear(c *gin.Context) {
	companyID := c.Param("companyID")
	var fy domain.FiscalYear
	if err := c.ShouldBindJSON(&fy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	fy.CompanyID = companyID
	if err := h.svc.CreateFiscalYear(c.Request.Context(), &fy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, fy)
}

func (h *CompanyHandler) GetFiscalYear(c *gin.Context) {
	id := c.Param("id")
	fy, err := h.svc.GetFiscalYear(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fy)
}

func (h *CompanyHandler) ListFiscalYears(c *gin.Context) {
	companyID := c.Param("companyID")
	fys, err := h.svc.ListFiscalYears(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fys)
}

// ─── Period V2 ─────────────────────────────────────────────────────

func (h *CompanyHandler) GeneratePeriods(c *gin.Context) {
	companyID := c.Param("companyID")
	var fy domain.FiscalYear
	if err := c.ShouldBindJSON(&fy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	fy.CompanyID = companyID
	periods, err := h.svc.GeneratePeriods(c.Request.Context(), &fy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, periods)
}

func (h *CompanyHandler) ClosePeriod(c *gin.Context) {
	companyID := c.Param("companyID")
	periodID := c.Param("periodID")
	userID := GetUserID(c)
	if err := h.svc.ClosePeriod(c.Request.Context(), companyID, periodID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period closed"})
}

func (h *CompanyHandler) ReopenPeriod(c *gin.Context) {
	companyID := c.Param("companyID")
	periodID := c.Param("periodID")
	userID := GetUserID(c)
	if err := h.svc.ReopenPeriod(c.Request.Context(), companyID, periodID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period reopened"})
}

func (h *CompanyHandler) PermanentClosePeriod(c *gin.Context) {
	companyID := c.Param("companyID")
	periodID := c.Param("periodID")
	userID := GetUserID(c)
	if err := h.svc.PermanentClosePeriod(c.Request.Context(), companyID, periodID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period permanently closed"})
}

func (h *CompanyHandler) GetCurrentPeriod(c *gin.Context) {
	companyID := c.Param("companyID")
	period, err := h.svc.GetCurrentPeriod(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, period)
}

// ─── Department ────────────────────────────────────────────────────

func (h *CompanyHandler) CreateDepartment(c *gin.Context) {
	companyID := c.Param("companyID")
	var dept domain.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	dept.CompanyID = companyID
	if err := h.svc.CreateDepartment(c.Request.Context(), &dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dept)
}

func (h *CompanyHandler) GetDepartment(c *gin.Context) {
	id := c.Param("id")
	dept, err := h.svc.GetDepartment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dept)
}

func (h *CompanyHandler) ListDepartments(c *gin.Context) {
	companyID := c.Param("companyID")
	depts, err := h.svc.ListDepartments(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, depts)
}

func (h *CompanyHandler) UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var dept domain.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	dept.ID = id
	if err := h.svc.UpdateDepartment(c.Request.Context(), &dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dept)
}

// ─── Employee ──────────────────────────────────────────────────────

func (h *CompanyHandler) CreateEmployee(c *gin.Context) {
	companyID := c.Param("companyID")
	var emp domain.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	emp.CompanyID = companyID
	if err := h.svc.CreateEmployee(c.Request.Context(), &emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, emp)
}

func (h *CompanyHandler) GetEmployee(c *gin.Context) {
	id := c.Param("id")
	emp, err := h.svc.GetEmployee(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func (h *CompanyHandler) ListEmployees(c *gin.Context) {
	companyID := c.Param("companyID")
	emps, err := h.svc.ListEmployees(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emps)
}

func (h *CompanyHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var emp domain.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	emp.ID = id
	if err := h.svc.UpdateEmployee(c.Request.Context(), &emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func (h *CompanyHandler) DeactivateEmployee(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeactivateEmployee(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "employee deactivated"})
}

// ─── Bank Account ──────────────────────────────────────────────────

func (h *CompanyHandler) CreateBankAccount(c *gin.Context) {
	companyID := c.Param("companyID")
	var ba domain.CompanyBankAccount
	if err := c.ShouldBindJSON(&ba); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	ba.CompanyID = companyID
	if err := h.svc.CreateBankAccount(c.Request.Context(), &ba); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ba)
}

func (h *CompanyHandler) GetBankAccount(c *gin.Context) {
	id := c.Param("id")
	ba, err := h.svc.GetBankAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ba)
}

func (h *CompanyHandler) ListBankAccounts(c *gin.Context) {
	companyID := c.Param("companyID")
	bas, err := h.svc.ListBankAccounts(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bas)
}

func (h *CompanyHandler) UpdateBankAccount(c *gin.Context) {
	id := c.Param("id")
	var ba domain.CompanyBankAccount
	if err := c.ShouldBindJSON(&ba); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	ba.ID = id
	if err := h.svc.UpdateBankAccount(c.Request.Context(), &ba); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ba)
}

func (h *CompanyHandler) DeactivateBankAccount(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeactivateBankAccount(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bank account deactivated"})
}

// ─── E-Invoice ─────────────────────────────────────────────────────

func (h *CompanyHandler) RegisterEInvoicePattern(c *gin.Context) {
	companyID := c.Param("companyID")
	var inv domain.EInvoicePattern
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	inv.CompanyID = companyID
	if err := h.svc.RegisterEInvoicePattern(c.Request.Context(), &inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, inv)
}

func (h *CompanyHandler) GetEInvoicePattern(c *gin.Context) {
	id := c.Param("id")
	inv, err := h.svc.GetEInvoicePattern(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *CompanyHandler) ListEInvoicePatterns(c *gin.Context) {
	companyID := c.Param("companyID")
	invs, err := h.svc.ListEInvoicePatterns(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, invs)
}

// ─── Digital Signature ─────────────────────────────────────────────

func (h *CompanyHandler) RegisterDigitalSignature(c *gin.Context) {
	companyID := c.Param("companyID")
	var sig domain.DigitalSignature
	if err := c.ShouldBindJSON(&sig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	sig.CompanyID = companyID
	if err := h.svc.RegisterDigitalSignature(c.Request.Context(), &sig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sig)
}

func (h *CompanyHandler) GetDigitalSignature(c *gin.Context) {
	id := c.Param("id")
	sig, err := h.svc.GetDigitalSignature(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sig)
}

func (h *CompanyHandler) ListDigitalSignatures(c *gin.Context) {
	companyID := c.Param("companyID")
	sigs, err := h.svc.ListDigitalSignatures(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sigs)
}

// ─── Integration ───────────────────────────────────────────────────

func (h *CompanyHandler) CreateIntegrationProfile(c *gin.Context) {
	companyID := c.Param("companyID")
	var prof domain.IntegrationProfile
	if err := c.ShouldBindJSON(&prof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	prof.CompanyID = companyID
	if err := h.svc.CreateIntegrationProfile(c.Request.Context(), &prof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, prof)
}

func (h *CompanyHandler) GetIntegrationProfile(c *gin.Context) {
	id := c.Param("id")
	prof, err := h.svc.GetIntegrationProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prof)
}

func (h *CompanyHandler) ListIntegrationProfiles(c *gin.Context) {
	companyID := c.Param("companyID")
	profs, err := h.svc.ListIntegrationProfiles(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profs)
}

func (h *CompanyHandler) UpdateIntegrationProfile(c *gin.Context) {
	id := c.Param("id")
	var prof domain.IntegrationProfile
	if err := c.ShouldBindJSON(&prof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	prof.ID = id
	if err := h.svc.UpdateIntegrationProfile(c.Request.Context(), &prof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prof)
}

func (h *CompanyHandler) TestIntegration(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.TestIntegration(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "integration test passed"})
}
