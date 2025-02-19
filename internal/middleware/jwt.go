package middleware

import (
	v1 "bk/api/v1"
	"bk/pkg/jwt"
	"bk/pkg/log"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StrictAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			logger.WithContext(ctx).Warn("No token", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}), zap.Error(err))
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		// 获取客户端真实IP
		clientIP := getClientIP(ctx.Request)

		// IP绑定验证
		if claims.ClientIP != clientIP {
			// 将失效token加入黑名单
			if err := j.AddToBlacklist(ctx, tokenString, claims.ExpiresAt); err != nil {
				logger.Error("加入黑名单失败", zap.Error(err))
			}
			logger.WithContext(ctx).Warn("IP changed",
				zap.String("tokenIP", claims.ClientIP),
				zap.String("currentIP", clientIP))
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		// 新增黑名单检查
		if j.IsInBlacklist(ctx, tokenString) {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		// 后台接口需要admin权限
		if claims.UserType != "admin" {
			v1.HandleError(ctx, http.StatusForbidden, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func NoStrictAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			tokenString, _ = ctx.Cookie("accessToken")
		}
		if tokenString == "" {
			tokenString = ctx.Query("accessToken")
		}
		if tokenString == "" {
			ctx.Next()
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}
		fmt.Println(claims.UserId)
		if claims.UserId != "BqnsOt4iZ9" {
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func recoveryLoggerFunc(ctx *gin.Context, logger *log.Logger) {
	if userInfo, ok := ctx.MustGet("claims").(*jwt.MyCustomClaims); ok {
		logger.WithValue(ctx, zap.String("UserId", userInfo.UserId))
	}
}

// 获取客户端IP
func getClientIP(r *http.Request) string {
	// Try Forwarded header (RFC 7239 standard)
	if forwarded := r.Header.Get("Forwarded"); forwarded != "" {
		if parts := strings.Split(forwarded, ";"); len(parts) > 0 {
			for _, part := range parts {
				part = strings.TrimSpace(parts[0])
				if strings.HasPrefix(part, "for=") {
					ip := strings.TrimPrefix(part, "for=")
					ip = strings.Trim(ip, `"`) // Remove quotes if present
					return strings.Split(ip, ",")[0]
				}
			}
		}
	}

	// Fallback to X-Forwarded-For (common de facto standard)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}

	// Final fallback to remote address
	return strings.Split(r.RemoteAddr, ":")[0] // Remove port number
}
