package menu

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuDelete 删除菜单/按钮项
// 规则: 若存在子节点则拒绝; 同时清理角色-菜单关联
func MenuDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "ID 无效")
		return
	}

	var menu model.Menu
	if err := common.DB.First(&menu, id).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "记录不存在")
		return
	}

	// 存在子节点则禁止删除
	var childCount int64
	common.DB.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		common.Fail(c, common.CodeBadRequest, "存在子节点, 请先删除子项")
		return
	}

	tx := common.DB.Begin()
	// 清理角色-菜单关联
	if err := tx.Where("menu_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "清理角色关联失败")
		return
	}
	if err := tx.Delete(&menu).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "删除失败")
		return
	}
	tx.Commit()
	// 菜单删除后刷新所有角色权限缓存
	go refreshAllRoles()
	common.SuccessWithMessage(c, "删除成功")
}
