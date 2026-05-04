package common

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ================================================================
// 路由元数据 + 中间件
// 每个路由可挂载三个参数:
//   - permCode : 权限码 (空字符串表示只需登录即可)
//   - audit    : 是否记录审计日志
//   - action   : 审计操作描述 (如 "创建用户")
// ================================================================

// RouteConfig 返回一个同时处理权限校验和审计记录的中间件链
// 用法: router.POST("/users", RouteConfig("system:user:w", true, "创建用户"), handler)
func RouteConfig(permCode string, audit bool, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ---- 1. 权限校验 ----
		if permCode != "" {
			info := GetCurrentUser(c)
			if info == nil {
				Fail(c, CodeUnauthorized, "未认证")
				c.Abort()
				return
			}
			// 从 Redis 角色权限缓存检查
			if !HasPermission(info.RoleIDs, permCode) {
				Fail(c, CodeForbidden, "无权限: "+permCode)
				c.Abort()
				return
			}
		}

		// ---- 2. 审计记录 ----
		if audit {
			// 读取请求体(需要先缓存,否则 handler 里读不到)
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			// 脱敏: 去掉 password / machineCode 字段
			bodyStr := desensitizeBody(bodyBytes)

			// 替换响应 Writer 以捕获业务响应码
			rw := &responseWriter{ResponseWriter: c.Writer}
			c.Writer = rw

			start := time.Now()
			c.Next()
			duration := time.Since(start).Milliseconds()

			// 异步写入审计日志
			info := GetCurrentUser(c)
			userID := uint(0)
			userName := "anonymous"
			if info != nil {
				userID = info.UserID
				userName = info.UserName
			}

			go writeAuditLog(userID, userName, action,
				c.Request.Method, c.Request.URL.Path,
				c.Request.URL.RawQuery, bodyStr,
				rw.bizCode, ClientIP(c), duration)
			return
		}

		c.Next()
	}
}

// PermMiddleware 单独的权限校验中间件(不含审计)
func PermMiddleware(permCode string) gin.HandlerFunc {
	return RouteConfig(permCode, false, "")
}

// ----------------------------------------------------------------
// 内部工具
// ----------------------------------------------------------------

// responseWriter 包装 gin.ResponseWriter 以捕获业务响应码和消息
type responseWriter struct {
	gin.ResponseWriter
	bizCode int
	bizMsg  string
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(b, &resp); err == nil && resp.Code != 0 {
		rw.bizCode = resp.Code
		rw.bizMsg = resp.Message
	}
	return rw.ResponseWriter.Write(b)
}

// desensitizeBody 对请求体 JSON 进行脱敏处理
var sensitiveFields = []string{"password", "machineCode", "passwordHash"}

func desensitizeBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		// 非 JSON 直接截断返回
		if len(body) > 500 {
			return string(body[:500]) + "..."
		}
		return string(body)
	}
	for _, field := range sensitiveFields {
		if _, ok := m[field]; ok {
			m[field] = "***"
		}
	}
	result, _ := json.Marshal(m)
	return string(result)
}

// ClientIP 获取真实客户端 IP
func ClientIP(c *gin.Context) string {
	// 优先从代理头获取
	for _, header := range []string{"X-Real-IP", "X-Forwarded-For"} {
		if ip := c.GetHeader(header); ip != "" {
			return strings.Split(ip, ",")[0]
		}
	}
	return c.RemoteIP()
}

// writeAuditLog 异步写入审计日志(避免循环依赖,直接操作 DB)
func writeAuditLog(userID uint, userName, action, method, path, query, body string, statusCode int, ip string, duration int64) {
	// 使用 map 避免 import model 包造成循环依赖
	DB.Table(TableSysAuditLog).Create(map[string]interface{}{
		"user_id":     userID,
		"user_name":   userName,
		"action":      action,
		"method":      method,
		"path":        path,
		"query":       query,
		"body":        body,
		"status_code": statusCode,
		"ip":          ip,
		"duration":    duration,
	})
}
