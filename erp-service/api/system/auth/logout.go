package auth

import (
	"erp-service/common"

	"github.com/gin-gonic/gin"
)

// Logout 退出登录
func Logout(c *gin.Context) {
	// 从上下文获取 token（由 AuthMiddleware 设置）
	tokenValue, exists := c.Get("token")
	if !exists {
		common.Fail(c, common.CodeUnauthorized, "未认证")
		return
	}
	token, ok := tokenValue.(string)
	if !ok || token == "" {
		common.Fail(c, common.CodeUnauthorized, "未认证")
		return
	}

	if err := common.DeleteToken(token); err != nil {
		common.Fail(c, common.CodeInternalError, "退出失败")
		return
	}

	common.SuccessWithMessage(c, "退出成功")
}
