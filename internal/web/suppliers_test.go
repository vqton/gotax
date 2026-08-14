package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/auth"
	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

type suppliersSetup struct {
	r    *gin.Engine
	repo *repository.MemoryPurchaseRepo
}

func setupSuppliers(t *testing.T) *suppliersSetup {
	t.Helper()
	auth.SetJWTSecret("web-test-secret")

	repo := repository.NewMemoryPurchaseRepo()
	purchaseSvc := service.NewPurchaseService(repo, repo, repo, repo, repo, repo, repo, repo, repo, nil)
	deps := Deps{Purchase: purchaseSvc}

	srv, err := NewServer([]string{"suppliers", "legacy"})
	require.NoError(t, err)

	r := gin.New()
	token, err := auth.GenerateAccessToken(&domain.User{ID: "web-test-user", Username: "admin", Role: domain.UserRoleAdmin})
	require.NoError(t, err)
	// Engine middleware before routes: gin snapshots handlers at registration.
	r.Use(func(c *gin.Context) {
		c.Request.Header.Set("Authorization", "Bearer "+token)
	})
	srv.RegisterPages(r, NewPages(deps), srv.NewActions(deps))

	return &suppliersSetup{r: r, repo: repo}
}

func (s *suppliersSetup) seedSupplier(t *testing.T, code, name string) {
	t.Helper()
	err := s.repo.CreateSupplier(context.Background(), &domain.Supplier{
		CompanyID: "CMP001",
		Code:      code,
		Name:      name,
		TaxCode:   "0123456789",
		Currency:  "VND",
		Status:    domain.SupplierActive,
	})
	require.NoError(t, err)
}

func TestSuppliersPageRender(t *testing.T) {
	s := setupSuppliers(t)
	s.seedSupplier(t, "NCC-TEST", "Công ty cung cấp XYZ")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/suppliers.html?company_id=CMP001", nil)
	s.r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Nhà cung cấp")
	assert.Contains(t, body, "NCC-TEST")
	assert.Contains(t, body, "Công ty cung cấp XYZ")
	assert.NotContains(t, body, "x-data")
}

func TestSuppliersCreateAction(t *testing.T) {
	s := setupSuppliers(t)

	form := url.Values{
		"code":    {"NCC-NEW"},
		"name":    {"Công ty cung cấp mới"},
		"tax_code": {"0987654321"},
		"phone":   {"024 111 222"},
		"email":   {"supply@new.vn"},
		"address": {"Hải Phòng"},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/suppliers/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "NCC-NEW")
	assert.Contains(t, body, "Công ty cung cấp mới")

	sups, _, err := s.repo.ListSuppliers(context.Background(), domain.PurchaseOrderFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	require.Len(t, sups, 1)
	assert.Equal(t, "NCC-NEW", sups[0].Code)
	assert.NotEmpty(t, sups[0].ID)
}

func TestSuppliersCreateActionValidationError(t *testing.T) {
	s := setupSuppliers(t)

	form := url.Values{"code": {"NCC-BAD"}} // missing name + tax_code
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/suppliers/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
	sups, _, err := s.repo.ListSuppliers(context.Background(), domain.PurchaseOrderFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	assert.Empty(t, sups)
}

func TestSuppliersDeleteAction(t *testing.T) {
	s := setupSuppliers(t)
	s.seedSupplier(t, "NCC-DEL", "Công ty cung cấp xóa")
	sups, _, err := s.repo.ListSuppliers(context.Background(), domain.PurchaseOrderFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	require.Len(t, sups, 1)
	id := sups[0].ID

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/suppliers/delete?company_id=CMP001", strings.NewReader("id="+id))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	sups, _, err = s.repo.ListSuppliers(context.Background(), domain.PurchaseOrderFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	assert.Empty(t, sups)
}
