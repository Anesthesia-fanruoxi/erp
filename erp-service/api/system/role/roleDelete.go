package role

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleDelete 删除角色
func RoleDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效的角色ID")
		return
	}

	// 检查是否有用户正在使用该角色
	var userCount int64
	common.DB.Model(&model.UserRole{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		common.Fail(c, common.CodeBadRequest, "该角色下还有用户, 无法删除")
		return
	}

	tx := common.DB.Begin()

	// 删除角色-菜单关联
	if err := tx.Where("role_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "删除菜单关联失败")
		return
	}

	// 删除角色
	if err := tx.Delete(&model.Role{}, id).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "删除角色失败")
		return
	}

	tx.Commit()
	// 清理该角色的权限缓存
	common.RDB.Del(common.Ctx, common.RedisRolePermsKey(uint(id)))
	common.RDB.Del(common.Ctx, common.RedisRoleMenusKey(uint(id)))
	common.SuccessWithMessage(c, "删除成功")
}
