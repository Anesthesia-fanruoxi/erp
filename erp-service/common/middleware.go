package common

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware Token认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header提取 Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			Fail(c, CodeUnauthorized, "缺少认证Token")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			Fail(c, CodeUnauthorized, "认证Token格式错误")
			c.Abort()
			return
		}

		token := parts[1]
		info, err := GetTokenInfo(token)
		if err != nil {
			Fail(c, CodeUnauthorized, "Token无效或已过期")
			c.Abort()
			return
		}

		// 将TokenInfo和token注入Context
		c.Set("tokenInfo", info)
		c.Set("token", token)
		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		info := GetCurrentUser(c)
		if info == nil {
			Fail(c, CodeUnauthorized, "未认证")
			c.Abort()
			return
		}

		isAdmin := false
		for _, role := range info.Roles {
			if role == "admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			Fail(c, CodeForbidden, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware 异常恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				LogError("Panic recovered: %v", err)
				Fail(c, CodeInternalError, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RequestLogMiddleware 请求日志中间件（记录业务响应码和消息，非200时输出ERROR）
func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 包装 ResponseWriter 捕获业务响应码和消息
		rw := &responseWriter{ResponseWriter: c.Writer}
		c.Writer = rw

		c.Next()

		duration := time.Since(start)
		bizCode := rw.bizCode
		if bizCode == 0 {
			bizCode = CodeSuccess
		}

		if bizCode == CodeSuccess {
			LogInfo("%s %s bizCode=%d %v", c.Request.Method, c.Request.URL.Path, bizCode, duration)
		} else {
			LogError("%s %s bizCode=%d msg=%s %v", c.Request.Method, c.Request.URL.Path, bizCode, rw.bizMsg, duration)
		}
	}
}

// GetCurrentUser 从Context获取当前用户信息
func GetCurrentUser(c *gin.Context) *TokenInfo {
	value, exists := c.Get("tokenInfo")
	if !exists {
		return nil
	}
	info, ok := value.(*TokenInfo)
	if !ok {
		return nil
	}
	return info
}
