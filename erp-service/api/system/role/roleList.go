package role

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// RoleMenuPermItem 角色菜单权限项(返回给前端)
type RoleMenuPermItem struct {
	MenuID    uint   `json:"menuId"`
	PermTypes []int8 `json:"permTypes"`
}

// RoleListItem 角色列表项
type RoleListItem struct {
	model.Role
	MenuPerms []RoleMenuPermItem `json:"menuPerms"`
}

// RoleList 角色列表(含每个菜单的权限类型)
func RoleList(c *gin.Context) {
	var roles []model.Role
	if err := common.DB.Order("id ASC").Find(&roles).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "查询角色失败")
		return
	}

	// 批量查询所有角色的菜单权限关联
	roleIDs := make([]uint, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}

	var roleMenus []model.RoleMenu
	if len(roleIDs) > 0 {
		common.DB.Where("role_id IN ?", roleIDs).Find(&roleMenus)
	}

	// 按 roleID -> menuID -> []permType 分组
	type key struct{ roleID, menuID uint }
	permMap := make(map[key][]int8)
	for _, rm := range roleMenus {
		k := key{rm.RoleID, rm.MenuID}
		permMap[k] = append(permMap[k], rm.PermType)
	}

	// 组装结果
	result := make([]RoleListItem, 0, len(roles))
	for _, r := range roles {
		// 收集该角色涉及的所有 menuID
		menuIDSet := make(map[uint]struct{})
		for k := range permMap {
			if k.roleID == r.ID {
				menuIDSet[k.menuID] = struct{}{}
			}
		}

		menuPerms := make([]RoleMenuPermItem, 0, len(menuIDSet))
		for menuID := range menuIDSet {
			menuPerms = append(menuPerms, RoleMenuPermItem{
				MenuID:    menuID,
				PermTypes: permMap[key{r.ID, menuID}],
			})
		}

		result = append(result, RoleListItem{
			Role:      r,
			MenuPerms: menuPerms,
		})
	}

	common.Success(c, result)
}
