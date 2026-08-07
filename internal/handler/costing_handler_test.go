package handler

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupCostingTest() (*gin.Engine, *service.CostingPeriodService, *service.CostingEngine, *repository.MemoryCostObjectRepo, *repository.MemoryCostPoolRepo, *repository.MemoryCostPoolLineRepo) {
	gin.SetMode(gin.TestMode)

	periodRepo := repository.NewMemoryCostingPeriodRepo()
	costObjectRepo := repository.NewMemoryCostObjectRepo()
	costPoolRepo := repository.NewMemoryCostPoolRepo()
	costPoolLineRepo := repository.NewMemoryCostPoolLineRepo()
	resultRepo := repository.NewMemoryCostingResultRepo()
	resultLineRepo := repository.NewMemoryCostingResultLineRepo()

	periodSvc := service.NewCostingPeriodService(periodRepo)
	engine := service.NewCostingEngine(periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo)
	costingJESvc := service.NewCostingJEService(periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo, &mockJECreator{})

	h := NewCostingHandler(periodSvc, engine, costingJESvc)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("role", "admin")
		c.Next()
	})
	RegisterCostingRoutes(router, h, func(c *gin.Context) { c.Next() })

	return router, periodSvc, engine, costObjectRepo, costPoolRepo, costPoolLineRepo
}

type mockJECreator struct{}

func (m *mockJECreator) CreateEntry(_ context.Context, _ *domain.JournalEntry, _ string) error {
	return nil
}

func TestCostingCreatePeriod(t *testing.T) {
	router, _, _, _, _, _ := setupCostingTest()

	body := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestCostingCreatePeriod_EmptyCompany(t *testing.T) {
	router, _, _, _, _, _ := setupCostingTest()

	body := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestCostingCreatePeriod_Duplicate(t *testing.T) {
	router, _, _, _, _, _ := setupCostingTest()

	body := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	req2 := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 409, w2.Code)
}

func TestCostingListPeriods(t *testing.T) {
	router, _, _, _, _, _ := setupCostingTest()

	body := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	req2 := httptest.NewRequest("GET", "/api/v1/costing-periods?company_id=COMP1", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)
	assert.Contains(t, w2.Body.String(), "2026")
}

func TestCostingGetPeriod_NotFound(t *testing.T) {
	router, _, _, _, _, _ := setupCostingTest()

	req := httptest.NewRequest("GET", "/api/v1/costing-periods/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestCostingClosePeriod(t *testing.T) {
	router, periodSvc, _, _, _, _ := setupCostingTest()

	body := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	periods, _ := periodSvc.List(context.Background(), "COMP1")
	assert.Greater(t, len(periods), 0)
	periodID := periods[0].ID

	req3 := httptest.NewRequest("POST", "/api/v1/costing-periods/"+periodID+"/close", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 200, w3.Code)
}

func TestCostingClosePeriod_AlreadyClosed(t *testing.T) {
	router, periodSvc, _, _, _, _ := setupCostingTest()

	body := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	periods, _ := periodSvc.List(context.Background(), "COMP1")
	periodID := periods[0].ID

	req2 := httptest.NewRequest("POST", "/api/v1/costing-periods/"+periodID+"/close", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)

	req3 := httptest.NewRequest("POST", "/api/v1/costing-periods/"+periodID+"/close", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 409, w3.Code)
}

func TestCostingRun_HappyPath(t *testing.T) {
	router, periodSvc, _, costObjectRepo, costPoolRepo, costPoolLineRepo := setupCostingTest()

	periodBody := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(periodBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	periods, _ := periodSvc.List(context.Background(), "COMP1")
	periodID := periods[0].ID

	obj := &domain.CostObject{
		ID:            "OBJ-001",
		CompanyID:     "COMP1",
		Code:          "SP001",
		Name:          "Product A",
		Type:          "PRODUCT",
		CostingMethod: "SIMPLE",
		PlanQuantity:  100,
		IsActive:      true,
	}
	_ = costObjectRepo.Create(context.Background(), obj)

	pool := &domain.CostPool{
		ID:            "POOL-001",
		CompanyID:     "COMP1",
		PeriodID:      periodID,
		GLAccountCode: "621",
		Name:          "Direct Materials",
		Status:        "OPEN",
		TotalAmount:   5000000,
	}
	_ = costPoolRepo.Create(context.Background(), pool)

	line := &domain.CostPoolLine{
		ID:          "LINE-001",
		PoolID:      "POOL-001",
		SourceType:  "JOURNAL",
		Description: "Material from GRN-001",
		Amount:      5000000,
	}
	_ = costPoolLineRepo.Create(context.Background(), line)

	runBody := `{"company_id":"COMP1","period_id":"` + periodID + `"}`
	req2 := httptest.NewRequest("POST", "/api/v1/costing/run", strings.NewReader(runBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)
	assert.Contains(t, w2.Body.String(), "costing completed")
}

func TestCostingRun_ClosedPeriod(t *testing.T) {
	router, periodSvc, _, _, _, _ := setupCostingTest()

	periodBody := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(periodBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	periods, _ := periodSvc.List(context.Background(), "COMP1")
	periodID := periods[0].ID

	req2 := httptest.NewRequest("POST", "/api/v1/costing-periods/"+periodID+"/close", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)

	runBody := `{"company_id":"COMP1","period_id":"` + periodID + `"}`
	req3 := httptest.NewRequest("POST", "/api/v1/costing/run", strings.NewReader(runBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 409, w3.Code)
}

func TestCostingRun_NoPeriods(t *testing.T) {
	router, _, _, _, _, _ := setupCostingTest()

	body := `{"company_id":"COMP1","period_id":"nonexistent"}`
	req := httptest.NewRequest("POST", "/api/v1/costing/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
