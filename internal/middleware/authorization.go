package middleware

import (
	v1 "bk/api/v1"
	"bk/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			v1.HandleError(c, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			c.Abort()
			return
		}

		customClaims := claims.(*jwt.MyCustomClaims)
		if customClaims.UserType != "admin" {
			v1.HandleError(c, http.StatusForbidden, v1.ErrUnauthorized, nil)
			c.Abort()
			return
		}
	}
}
