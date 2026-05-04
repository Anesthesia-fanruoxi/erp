package role

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// RoleCreate 创建角色(含菜单权限关联)
func RoleCreate(c *gin.Context) {
	var req model.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 检查角色名是否重复
	var count int64
	common.DB.Model(&model.Role{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		common.Fail(c, common.CodeBadRequest, "角色名已存在")
		return
	}

	tx := common.DB.Begin()

	r := model.Role{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := tx.Create(&r).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "创建角色失败")
		return
	}

	// 批量创建菜单权限关联
	links := buildRoleMenuLinks(r.ID, req.MenuPerms)
	if len(links) > 0 {
		if err := tx.Create(&links).Error; err != nil {
			tx.Rollback()
			common.Fail(c, common.CodeInternalError, "创建权限关联失败")
			return
		}
	}

	tx.Commit()
	// 刷新该角色的权限缓存
	go common.RefreshRolePerms(r.ID)
	common.SuccessWithMessage(c, "创建成功")
}

// buildRoleMenuLinks 将前端传入的 menuPerms 展开为 RoleMenu 记录列表
func buildRoleMenuLinks(roleID uint, menuPerms []model.RoleMenuPerm) []model.RoleMenu {
	links := make([]model.RoleMenu, 0)
	for _, mp := range menuPerms {
		for _, pt := range mp.PermTypes {
			links = append(links, model.RoleMenu{
				RoleID:   roleID,
				MenuID:   mp.MenuID,
				PermType: pt,
			})
		}
	}
	return links
}
