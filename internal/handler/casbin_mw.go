package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gotax/internal/authz"
)

// CasbinMW enforces Casbin permission after authMW has set user_id, tenant_id, role.
// obj format: "<module>:<resource>:<action>" or just "<resource>:<action>"
//
// Usage (after authMW):
//   v1 := r.Group("/api/v1", authMW)
//   v1.Use(CasbinMW())
//   v1.GET("/reports", h.GetReports)
//
// Skip a route:
//   v1.GET("/public", pubHandler, SkipCasbin())
func CasbinMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		tenantID := c.GetString("tenant_id")
		rawObj := c.FullPath() // e.g. /api/v1/journal-entries
		if rawObj == "" {
			rawObj = c.Request.URL.Path
		}
		act := httpMethodToAct(c.Request.Method)

		// Build a resource fingerprint stripping version prefix → journal-entries
		parts := strings.SplitN(rawObj, "/", 4)
		obj := rawObj
		if len(parts) >= 4 {
			obj = "/" + parts[3]
		}

		// Build subject: user_id (flat RBAC) — role bindings resolved by g() rule.
		sub := userID.(string)

		if !authz.Enforce(sub, tenantID+obj, act) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
				"code":  "AUTH_004",
				"required": obj + ":" + act,
				"user": sub,
				"tenant": tenantID,
			})
			return
		}
		c.Next()
	}
}

// RequirePerm is a per-route helper (use as decorator middleware).
// Example: accounts.GET("/:code", RequirePerm("accounts:read"), h.GetAccount)
func RequirePerm(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		tenantID := c.GetString("tenant_id")
		if !authz.Enforce(userID.(string), tenantID+obj, act) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden", "code": "AUTH_004"})
			return
		}
		c.Next()
	}
}

func httpMethodToAct(m string) string {
	switch m {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return m
	}
}
