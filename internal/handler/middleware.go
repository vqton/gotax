package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gotax/internal/auth"
	"gotax/internal/domain"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		claims, err := auth.ParseAndValidateAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", string(claims.Role))
		c.Next()
	}
}

func RoleMiddleware(roles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		role := domain.UserRole(roleStr.(string))
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

func GetUserID(c *gin.Context) string {
	uid, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	s, ok := uid.(string)
	if !ok {
		return ""
	}
	return s
}

func GetUserRole(c *gin.Context) domain.UserRole {
	roleStr, exists := c.Get("role")
	if !exists {
		return ""
	}
	return domain.UserRole(roleStr.(string))
}

func GetUsername(c *gin.Context) string {
	u, exists := c.Get("username")
	if !exists {
		return ""
	}
	return u.(string)
}
