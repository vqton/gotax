package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/auth"
	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

type customersSetup struct {
	r    *gin.Engine
	repo *repository.MemorySaleRepo
}

// setupCustomers wires the real pages pipeline: memory sale repo → real
// SaleService → Deps → NewServer template sets → RegisterPages, authenticated
// with a real RS256 token through PageAuthMiddleware.
func setupCustomers(t *testing.T) *customersSetup {
	t.Helper()
	auth.SetJWTSecret("web-test-secret")

	repo := repository.NewMemorySaleRepo()
	saleSvc := service.NewSaleService(repo, repo, repo, repo, repo, repo, repo, repo, nil)
	deps := Deps{Sale: saleSvc}

	srv, err := NewServer([]string{"customers", "legacy"})
	require.NoError(t, err)

	r := gin.New()
	token, err := auth.GenerateAccessToken(&domain.User{ID: "web-test-user", Username: "admin", Role: domain.UserRoleAdmin})
	require.NoError(t, err)
	// Engine middleware must be registered before routes: gin snapshots the
	// engine handler chain into route nodes at registration time.
	r.Use(func(c *gin.Context) {
		c.Request.Header.Set("Authorization", "Bearer "+token)
	})
	srv.RegisterPages(r, NewPages(deps), srv.NewActions(deps))

	return &customersSetup{r: r, repo: repo}
}

// Seed one customer through the real service so the page renders repo state,
// not mock data.
func (s *customersSetup) seedCustomer(t *testing.T, code, name string) {
	t.Helper()
	err := s.repo.CreateCustomer(context.Background(), &domain.Customer{
		CompanyID: "CMP001",
		Code:      code,
		Name:      name,
		TaxCode:   "0123456789",
		Currency:  "VND",
		Status:    domain.CustomerActive,
	})
	require.NoError(t, err)
}

func TestCustomersPageRender(t *testing.T) {
	s := setupCustomers(t)
	s.seedCustomer(t, "KH-TEST", "Công ty kiểm thử XYZ")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/customers.html?company_id=CMP001", nil)
	s.r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Khách hàng")
	assert.Contains(t, body, "KH-TEST")
	assert.Contains(t, body, "Công ty kiểm thử XYZ")
	// Converted pages must not carry Alpine directives anymore.
	assert.NotContains(t, body, "x-data")
}
