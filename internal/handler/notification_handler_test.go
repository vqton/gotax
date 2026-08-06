package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupNotificationTest(t *testing.T) (*gin.Engine, *service.NotificationService, domain.NotificationRepository) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	notifRepo := repository.NewMemoryNotificationRepo()
	notifSvc := service.NewNotificationService(notifRepo)
	nh := NewNotificationHandler(notifSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterNotificationRoutes(r, nh, noopMW)

	return r, notifSvc, notifRepo
}

func seedNotification(t *testing.T, svc *service.NotificationService, companyID, userID string) string {
	t.Helper()
	n := &domain.Notification{
		CompanyID: companyID,
		UserID:    userID,
		Type:      domain.NotifTypeINFO,
		Title:     "Test Notification",
		Message:   "This is a test",
	}
	err := svc.Create(t.Context(), n)
	require.NoError(t, err)
	return n.ID
}

func TestNotificationList(t *testing.T) {
	r, svc, _ := setupNotificationTest(t)
	seedNotification(t, svc, "CMP001", "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["total"])
}

func TestNotificationUnreadCount(t *testing.T) {
	r, svc, _ := setupNotificationTest(t)
	seedNotification(t, svc, "CMP001", "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/unread-count?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["count"])
}

func TestNotificationMarkRead(t *testing.T) {
	r, svc, _ := setupNotificationTest(t)
	id := seedNotification(t, svc, "CMP001", "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/"+id+"/read", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// unread should be 0 now
	count, _ := svc.UnreadCount(t.Context(), "CMP001", "test-user")
	assert.Equal(t, 0, count)
}

func TestNotificationMarkAllRead(t *testing.T) {
	r, svc, _ := setupNotificationTest(t)
	seedNotification(t, svc, "CMP001", "test-user")
	n2 := &domain.Notification{
		CompanyID: "CMP001",
		UserID:    "test-user",
		Type:      domain.NotifTypeWARNING,
		Title:     "Warning",
		Message:   "Something",
	}
	svc.Create(t.Context(), n2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/read-all?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	count, _ := svc.UnreadCount(t.Context(), "CMP001", "test-user")
	assert.Equal(t, 0, count)
}

func TestNotificationDelete(t *testing.T) {
	r, svc, _ := setupNotificationTest(t)
	id := seedNotification(t, svc, "CMP001", "test-user")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/"+id, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	_, err := svc.List(t.Context(), "CMP001", "test-user", 50)
	require.NoError(t, err)
	assert.Equal(t, 0, len([]domain.Notification{}))
}

func TestNotificationDeleteNotFound(t *testing.T) {
	r, _, _ := setupNotificationTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestNotificationMarkReadNotFound(t *testing.T) {
	r, _, _ := setupNotificationTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/nonexistent/read", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestNotificationListEmpty(t *testing.T) {
	r, _, _ := setupNotificationTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?company_id=CMP999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["total"])
}

func TestNotificationListMissingCompanyID(t *testing.T) {
	r, _, _ := setupNotificationTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestNotificationCreateViaService(t *testing.T) {
	_, svc, _ := setupNotificationTest(t)

	n := &domain.Notification{
		CompanyID: "CMP001",
		UserID:    "user1",
		Type:      domain.NotifTypeDUE,
		Title:     "Tax due",
		Message:   "Deadline approaching",
	}
	err := svc.Create(t.Context(), n)
	require.NoError(t, err)
	assert.NotEmpty(t, n.ID)
	assert.NotEmpty(t, n.CreatedAt)
	assert.Equal(t, domain.NotifTypeDUE, n.Type)
}

func TestNotificationTypes(t *testing.T) {
	r, svc, _ := setupNotificationTest(t)

	types := []domain.NotificationType{domain.NotifTypeINFO, domain.NotifTypeWARNING, domain.NotifTypeERROR, domain.NotifTypeDUE}
	for _, typ := range types {
		n := &domain.Notification{
			CompanyID: "CMP001",
			UserID:    "test-user",
			Type:      typ,
			Title:     string(typ),
			Message:   "msg",
		}
		err := svc.Create(t.Context(), n)
		require.NoError(t, err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?company_id=CMP001&limit=100", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(4), resp["total"])
}
