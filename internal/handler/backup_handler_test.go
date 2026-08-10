package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupBackupHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryBackupRepo()
	svc := service.NewBackupService(repo, "", "")
	h := NewBackupHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("role", "admin")
		c.Next()
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterBackupRoutes(r, h, noopMW)
	return r
}

func TestBackupList(t *testing.T) {
	r := setupBackupHandlerTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/backups?company_id=comp1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Empty(t, resp)
}
