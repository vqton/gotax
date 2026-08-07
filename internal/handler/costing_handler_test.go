package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupCostingTest() (*gin.Engine, *service.CostingPeriodService) {
	gin.SetMode(gin.TestMode)

	periodRepo := repository.NewMemoryCostingPeriodRepo()
	costObjectRepo := repository.NewMemoryCostObjectRepo()
	costPoolRepo := repository.NewMemoryCostPoolRepo()
	costPoolLineRepo := repository.NewMemoryCostPoolLineRepo()
	resultRepo := repository.NewMemoryCostingResultRepo()
	resultLineRepo := repository.NewMemoryCostingResultLineRepo()

	periodSvc := service.NewCostingPeriodService(periodRepo)
	engine := service.NewCostingEngine(periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo)

	h := NewCostingHandler(periodSvc, engine)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("role", "admin")
		c.Next()
	})
	RegisterCostingRoutes(router, h, func(c *gin.Context) { c.Next() })

	return router, periodSvc
}

func TestCostingCreatePeriod(t *testing.T) {
	router, _ := setupCostingTest()

	body := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestCostingCreatePeriod_Duplicate(t *testing.T) {
	router, _ := setupCostingTest()

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
	router, _ := setupCostingTest()

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

func TestCostingClosePeriod(t *testing.T) {
	router, periodSvc := setupCostingTest()

	body := `{"year":2026,"month":8}`
	req := httptest.NewRequest("POST", "/api/v1/costing-periods?company_id=COMP1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	periods, _ := periodSvc.List(nil, "COMP1")
	assert.Greater(t, len(periods), 0)
	periodID := periods[0].ID

	req3 := httptest.NewRequest("POST", "/api/v1/costing-periods/"+periodID+"/close", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 200, w3.Code)
}

func TestCostingRun_NoPeriods(t *testing.T) {
	router, _ := setupCostingTest()

	body := `{"company_id":"COMP1","period_id":"nonexistent"}`
	req := httptest.NewRequest("POST", "/api/v1/costing/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
