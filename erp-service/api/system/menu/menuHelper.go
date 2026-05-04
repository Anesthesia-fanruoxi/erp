package menu

import "erp-service/common"

// refreshAllRoles 刷新所有角色权限缓存（菜单变更后调用）
func refreshAllRoles() {
	var roleIDs []uint
	common.DB.Table(common.TableSysRole).Pluck("id", &roleIDs)
	for _, roleID := range roleIDs {
		common.RefreshRolePerms(roleID)
	}
}

// autoGrantAdminMenu 自动给 admin 角色分配新菜单的全部权限并刷新缓存
func autoGrantAdminMenu(menuID uint) {
	// 查询 admin 角色 ID
	var adminRole struct {
		ID uint `gorm:"column:id"`
	}
	if err := common.DB.Table(common.TableSysRole).
		Where("name = ?", "admin").
		First(&adminRole).Error; err != nil {
		common.LogWarn("自动授权: 未找到 admin 角色, err=%v", err)
		return
	}

	// 插入三种权限类型（查看=1, 只读=2, 操作=3），忽略已存在的
	for _, permType := range []int8{1, 2, 3} {
		common.DB.Exec(
			"INSERT IGNORE INTO "+common.TableSysRoleMenu+
				" (role_id, menu_id, perm_type) VALUES (?, ?, ?)",
			adminRole.ID, menuID, permType,
		)
	}

	// 刷新 admin 角色的 Redis 缓存
	if err := common.RefreshRolePerms(adminRole.ID); err != nil {
		common.LogWarn("自动授权: 刷新 admin 权限缓存失败, err=%v", err)
	} else {
		common.LogInfo("自动授权: 已给 admin 角色分配菜单 %d 的全部权限", menuID)
	}
}
