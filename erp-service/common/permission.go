package common

import (
	"encoding/json"
	"fmt"
	"sort"
)

// ================================================================
// 角色权限预热 / 刷新 / 查询
//
// Redis 结构:
//   role_perms:<roleId>  → Set<permCode>   角色权限码集合
//   role_menus:<roleId>  → JSON<[]MenuNode> 角色菜单树
//
// 权限码规则:
//   perm_type=1(查看) → code        菜单可见
//   perm_type=2(只读) → code:r      GET 接口
//   perm_type=3(操作) → code:w      写接口
// ================================================================

// roleMenuRow 查询角色菜单权限的中间结构
type roleMenuRow struct {
	MenuID   uint   `gorm:"column:menu_id"`
	PermType int8   `gorm:"column:perm_type"`
	Code     string `gorm:"column:code"`
	Name     string `gorm:"column:name"`
	Type     int8   `gorm:"column:type"`
	ParentID uint   `gorm:"column:parent_id"`
	Path     string `gorm:"column:path"`
	Icon     string `gorm:"column:icon"`
	Sort     int    `gorm:"column:sort"`
	Visible  int8   `gorm:"column:visible"`
}

// WarmUpAllRolePerms 启动时预热所有角色权限到 Redis
func WarmUpAllRolePerms() error {
	// 查询所有角色 ID
	var roleIDs []uint
	if err := DB.Table(TableSysRole).Pluck("id", &roleIDs).Error; err != nil {
		return fmt.Errorf("查询角色列表失败: %w", err)
	}

	for _, roleID := range roleIDs {
		if err := RefreshRolePerms(roleID); err != nil {
			LogWarn("预热角色 %d 权限失败: %v", roleID, err)
		}
	}

	LogInfo("角色权限预热完成, 共 %d 个角色", len(roleIDs))
	return nil
}

// RefreshRolePerms 刷新单个角色的权限缓存（角色权限变更时调用）
func RefreshRolePerms(roleID uint) error {
	rows, err := queryRoleMenuRows(roleID)
	if err != nil {
		return err
	}

	// 构建权限码集合
	permSet := make(map[string]struct{})
	menuIDSet := make(map[uint]struct{})
	for _, row := range rows {
		switch row.PermType {
		case 1: // 查看
			permSet[row.Code] = struct{}{}
			menuIDSet[row.MenuID] = struct{}{}
		case 2: // 只读
			permSet[row.Code+":r"] = struct{}{}
		case 3: // 操作
			permSet[row.Code+":w"] = struct{}{}
		}
	}

	// 写入权限码集合到 Redis (先删后写)
	key := RedisRolePermsKey(roleID)
	pipe := RDB.Pipeline()
	pipe.Del(Ctx, key)
	if len(permSet) > 0 {
		members := make([]interface{}, 0, len(permSet))
		for code := range permSet {
			members = append(members, code)
		}
		pipe.SAdd(Ctx, key, members...)
	}
	if _, err := pipe.Exec(Ctx); err != nil {
		return fmt.Errorf("写入角色权限到Redis失败: %w", err)
	}

	// 构建并写入菜单树
	menuTree := buildRoleMenuTree(rows, menuIDSet)
	menuJSON, err := json.Marshal(menuTree)
	if err != nil {
		return fmt.Errorf("序列化菜单树失败: %w", err)
	}
	if err := RDB.Set(Ctx, RedisRoleMenusKey(roleID), string(menuJSON), 0).Err(); err != nil {
		return fmt.Errorf("写入角色菜单树到Redis失败: %w", err)
	}

	return nil
}

// GetRolePerms 从 Redis 获取角色权限码集合
func GetRolePerms(roleID uint) ([]string, error) {
	perms, err := RDB.SMembers(Ctx, RedisRolePermsKey(roleID)).Result()
	if err != nil {
		return nil, fmt.Errorf("获取角色权限失败: %w", err)
	}
	return perms, nil
}

// GetRoleMenus 从 Redis 获取角色菜单树
func GetRoleMenus(roleID uint) ([]MenuNode, error) {
	data, err := RDB.Get(Ctx, RedisRoleMenusKey(roleID)).Result()
	if err != nil {
		return nil, fmt.Errorf("获取角色菜单失败: %w", err)
	}
	var menus []MenuNode
	if err := json.Unmarshal([]byte(data), &menus); err != nil {
		return nil, fmt.Errorf("反序列化菜单树失败: %w", err)
	}
	return menus, nil
}

// GetUserPermsAndMenus 合并用户所有角色的权限码和菜单树
func GetUserPermsAndMenus(roleIDs []uint) ([]string, []MenuNode) {
	permSet := make(map[string]struct{})
	// 用 map 去重菜单节点（按 ID）
	menuMap := make(map[uint]MenuNode)

	for _, roleID := range roleIDs {
		// 合并权限码
		perms, err := GetRolePerms(roleID)
		if err != nil {
			LogWarn("获取角色 %d 权限失败: %v", roleID, err)
			continue
		}
		for _, p := range perms {
			permSet[p] = struct{}{}
		}

		// 合并菜单（取顶层节点，按 ID 去重）
		menus, err := GetRoleMenus(roleID)
		if err != nil {
			LogWarn("获取角色 %d 菜单失败: %v", roleID, err)
			continue
		}
		mergeMenuNodes(menus, menuMap)
	}

	// 权限码转切片
	permCodes := make([]string, 0, len(permSet))
	for code := range permSet {
		permCodes = append(permCodes, code)
	}

	// 菜单 map 转有序切片
	menuList := make([]MenuNode, 0, len(menuMap))
	for _, node := range menuMap {
		menuList = append(menuList, node)
	}
	sort.Slice(menuList, func(i, j int) bool {
		return menuList[i].Sort < menuList[j].Sort
	})

	return permCodes, menuList
}

// HasPermission 检查角色列表是否包含指定权限码
func HasPermission(roleIDs []uint, permCode string) bool {
	for _, roleID := range roleIDs {
		ok, _ := RDB.SIsMember(Ctx, RedisRolePermsKey(roleID), permCode).Result()
		if ok {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------
// 内部工具
// ----------------------------------------------------------------

// queryRoleMenuRows 查询角色的菜单权限关联行
func queryRoleMenuRows(roleID uint) ([]roleMenuRow, error) {
	var rows []roleMenuRow
	err := DB.Table(TableSysRoleMenu).
		Select(TableSysRoleMenu+".menu_id, "+
			TableSysRoleMenu+".perm_type, "+
			TableSysMenu+".code, "+
			TableSysMenu+".name, "+
			TableSysMenu+".type, "+
			TableSysMenu+".parent_id, "+
			TableSysMenu+".path, "+
			TableSysMenu+".icon, "+
			TableSysMenu+".sort, "+
			TableSysMenu+".visible").
		Joins("JOIN "+TableSysMenu+" ON "+TableSysMenu+".id = "+TableSysRoleMenu+".menu_id").
		Where(TableSysRoleMenu+".role_id = ?", roleID).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("查询角色菜单权限失败: %w", err)
	}
	return rows, nil
}

// buildRoleMenuTree 从查询行构建菜单树（只含有查看权限的菜单）
func buildRoleMenuTree(rows []roleMenuRow, menuIDSet map[uint]struct{}) []MenuNode {
	// 收集有查看权限且可见的菜单
	type menuItem struct {
		ID       uint
		Name     string
		Type     int8
		ParentID uint
		Path     string
		Icon     string
		Sort     int
	}
	itemMap := make(map[uint]menuItem)
	for _, row := range rows {
		if row.PermType == 1 && row.Visible == 1 {
			if _, ok := itemMap[row.MenuID]; !ok {
				itemMap[row.MenuID] = menuItem{
					ID: row.MenuID, Name: row.Name, Type: row.Type,
					ParentID: row.ParentID, Path: row.Path,
					Icon: row.Icon, Sort: row.Sort,
				}
			}
		}
	}

	// 补全缺失的父目录节点
	needParents := make([]uint, 0)
	for _, item := range itemMap {
		if item.ParentID != 0 {
			if _, ok := itemMap[item.ParentID]; !ok {
				needParents = append(needParents, item.ParentID)
			}
		}
	}
	if len(needParents) > 0 {
		type parentRow struct {
			ID       uint   `gorm:"column:id"`
			Name     string `gorm:"column:name"`
			Type     int8   `gorm:"column:type"`
			ParentID uint   `gorm:"column:parent_id"`
			Path     string `gorm:"column:path"`
			Icon     string `gorm:"column:icon"`
			Sort     int    `gorm:"column:sort"`
		}
		var parents []parentRow
		DB.Table(TableSysMenu).
			Select("id, name, type, parent_id, path, icon, sort").
			Where("id IN ?", needParents).
			Scan(&parents)
		for _, p := range parents {
			itemMap[p.ID] = menuItem{
				ID: p.ID, Name: p.Name, Type: p.Type,
				ParentID: p.ParentID, Path: p.Path,
				Icon: p.Icon, Sort: p.Sort,
			}
		}
	}

	// 按 parentID 分组
	childrenMap := make(map[uint][]menuItem)
	for _, item := range itemMap {
		childrenMap[item.ParentID] = append(childrenMap[item.ParentID], item)
	}

	var build func(parentID uint) []MenuNode
	build = func(parentID uint) []MenuNode {
		children := childrenMap[parentID]
		if len(children) == 0 {
			return nil
		}
		sort.Slice(children, func(i, j int) bool {
			return children[i].Sort < children[j].Sort
		})
		nodes := make([]MenuNode, 0, len(children))
		for _, c := range children {
			nodes = append(nodes, MenuNode{
				ID:       c.ID,
				Name:     c.Name,
				Path:     c.Path,
				Icon:     c.Icon,
				Sort:     c.Sort,
				Children: build(c.ID),
			})
		}
		return nodes
	}

	return build(0)
}

// mergeMenuNodes 将菜单节点列表合并到 map（递归去重并保持排序）
func mergeMenuNodes(nodes []MenuNode, m map[uint]MenuNode) {
	for _, node := range nodes {
		if existing, ok := m[node.ID]; ok {
			// 已存在则合并子节点
			childMap := make(map[uint]MenuNode)
			mergeMenuNodes(existing.Children, childMap)
			mergeMenuNodes(node.Children, childMap)
			children := make([]MenuNode, 0, len(childMap))
			for _, c := range childMap {
				children = append(children, c)
			}
			// 子节点按 Sort 排序
			sort.Slice(children, func(i, j int) bool {
				return children[i].Sort < children[j].Sort
			})
			existing.Children = children
			m[node.ID] = existing
		} else {
			// 递归排序子节点
			sortMenuChildren(&node)
			m[node.ID] = node
		}
	}
}

// sortMenuChildren 递归排序菜单节点的子节点
func sortMenuChildren(node *MenuNode) {
	if len(node.Children) == 0 {
		return
	}
	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Sort < node.Children[j].Sort
	})
	for i := range node.Children {
		sortMenuChildren(&node.Children[i])
	}
}
