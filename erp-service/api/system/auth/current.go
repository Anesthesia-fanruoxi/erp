package auth

import (
	"erp-service/common"

	"github.com/gin-gonic/gin"
)

// Current 获取当前用户信息
func Current(c *gin.Context) {
	tokenInfo := common.GetCurrentUser(c)
	if tokenInfo == nil {
		common.Fail(c, common.CodeUnauthorized, "未认证")
		return
	}

	common.Success(c, gin.H{
		"id":       tokenInfo.UserID,
		"userName": tokenInfo.UserName,
		"realName": tokenInfo.RealName,
		"roles":    tokenInfo.Roles,
	})
}
