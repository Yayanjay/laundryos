package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/laundryos/backend/pkg/apiresponse"
	"github.com/laundryos/backend/pkg/jwt"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
	TenantIDKey         = "tenant_id"
	RoleKey             = "role"
)

func Auth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiresponse.Unauthorized())
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiresponse.Unauthorized())
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, apiresponse.Error(
					"TOKEN_EXPIRED",
					"Token Kadaluarsa",
					"Token Expired",
					"Silakan refresh token",
					"Please refresh your token",
				))
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiresponse.Unauthorized())
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(TenantIDKey, claims.TenantID)
		c.Set(RoleKey, claims.Role)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(RoleKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, apiresponse.Forbidden())
			return
		}

		userRole := role.(string)
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, apiresponse.Forbidden())
	}
}

func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get(UserIDKey); exists {
		return userID.(string)
	}
	return ""
}

func GetTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get(TenantIDKey); exists {
		return tenantID.(string)
	}
	return ""
}

func GetRole(c *gin.Context) string {
	if role, exists := c.Get(RoleKey); exists {
		return role.(string)
	}
	return ""
}
