package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gotax/internal/domain"
	"gotax/internal/service"
)

type DepartmentHandler struct {
	svc service.CompanyService
}

func NewDepartmentHandler(svc service.CompanyService) *DepartmentHandler {
	return &DepartmentHandler{svc: svc}
}

func RegisterDepartmentRoutes(r *gin.Engine, h *DepartmentHandler, authMW gin.HandlerFunc) {
	v1 := r.Group("/api/v1", authMW)

	depts := v1.Group("/departments")
	{
		depts.POST("", h.Create)
		depts.GET("", h.List)
		depts.GET("/:id", h.Get)
		depts.PUT("/:id", h.Update)
		depts.DELETE("/:id", h.Deactivate)
	}
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

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

func (h *DepartmentHandler) List(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

	depts, err := h.svc.ListDepartments(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, depts)
}

func (h *DepartmentHandler) Get(c *gin.Context) {
	id := c.Param("id")
	dept, err := h.svc.GetDepartment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dept)
}

func (h *DepartmentHandler) Update(c *gin.Context) {
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

func (h *DepartmentHandler) Deactivate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeactivateDepartment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "department deactivated"})
}
