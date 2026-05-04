package system

import "erp-service/common"

// Menu 系统菜单模型
// type: 1-菜单(有路由) 2-目录(无路由,仅分组)
// 权限标识规则:
//   code    → 菜单可见 (查看)
//   code:r  → GET 接口 (只读)
//   code:w  → POST/PUT/DELETE 接口 (操作)
type Menu struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex;size:100"`
	Name     string `json:"name" gorm:"size:100"`
	Type     int8   `json:"type" gorm:"column:type;default:1"` // 1-菜单 2-目录
	ParentID uint   `json:"parentId" gorm:"column:parent_id"`
	Path     string `json:"path" gorm:"size:200"`
	Icon     string `json:"icon" gorm:"size:50"`
	Sort     int    `json:"sort" gorm:"column:sort"`
	Visible  int8   `json:"visible" gorm:"column:visible;default:1"`
}

func (Menu) TableName() string {
	return common.TableSysMenu
}

// RoleMenu 角色-菜单权限关联模型
// PermType: 1-查看(菜单可见) 2-只读(GET接口) 3-操作(写接口)
type RoleMenu struct {
	RoleID   uint `json:"roleId" gorm:"primaryKey;column:role_id"`
	MenuID   uint `json:"menuId" gorm:"primaryKey;column:menu_id"`
	PermType int8 `json:"permType" gorm:"primaryKey;column:perm_type;default:1"`
}

func (RoleMenu) TableName() string {
	return common.TableSysRoleMenu
}

// --- DTO ---

// MenuCreateReq 创建菜单请求
type MenuCreateReq struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Type     int8   `json:"type" binding:"required,oneof=1 2"`
	ParentID uint   `json:"parentId"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	Visible  int8   `json:"visible"`
}

// MenuUpdateReq 编辑菜单请求
type MenuUpdateReq struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Type     int8   `json:"type"`
	ParentID uint   `json:"parentId"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	Visible  int8   `json:"visible"`
}
