package role

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleUpdate 编辑角色(含菜单权限关联)
func RoleUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效的角色ID")
		return
	}

	var req model.RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	var r model.Role
	if err := common.DB.First(&r, id).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "角色不存在")
		return
	}

	tx := common.DB.Begin()

	// 更新基本信息
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}
	if err := tx.Model(&model.Role{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "更新角色失败")
		return
	}

	// 全量替换菜单权限关联
	if req.MenuPerms != nil {
		if err := tx.Where("role_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			tx.Rollback()
			common.Fail(c, common.CodeInternalError, "更新权限关联失败")
			return
		}

		links := buildRoleMenuLinks(uint(id), req.MenuPerms)
		if len(links) > 0 {
			if err := tx.Create(&links).Error; err != nil {
				tx.Rollback()
				common.Fail(c, common.CodeInternalError, "创建权限关联失败")
				return
			}
		}
	}

	tx.Commit()
	// 刷新该角色的权限缓存
	go common.RefreshRolePerms(uint(id))
	common.SuccessWithMessage(c, "更新成功")
}
