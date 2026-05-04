package auth

import (
	"erp-service/common"

	"github.com/gin-gonic/gin"
)

// Menus 获取当前用户菜单和权限（从Redis角色缓存动态合并）
func Menus(c *gin.Context) {
	tokenInfo := common.GetCurrentUser(c)
	if tokenInfo == nil {
		common.Fail(c, common.CodeUnauthorized, "未认证")
		return
	}

	permissions, menus := common.GetUserPermsAndMenus(tokenInfo.RoleIDs)

	common.Success(c, gin.H{
		"menus":       menus,
		"permissions": permissions,
	})
}
