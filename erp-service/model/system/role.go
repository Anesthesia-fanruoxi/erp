package system

import "erp-service/common"

// Role 角色表模型
type Role struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name" gorm:"uniqueIndex;size:50"`
	Description string           `json:"description" gorm:"size:255"`
	CreatedAt   common.MilliTime `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (Role) TableName() string {
	return common.TableSysRole
}

// UserRole 用户-角色关联模型
type UserRole struct {
	UserID uint `json:"userId" gorm:"primaryKey;column:user_id"`
	RoleID uint `json:"roleId" gorm:"primaryKey;column:role_id"`
}

func (UserRole) TableName() string {
	return common.TableSysUserRole
}

// --- DTO ---

// RoleMenuPerm 角色菜单权限项 (前端传入)
// PermTypes: 勾选的权限类型列表, 如 [1,2,3] 表示查看+只读+操作
type RoleMenuPerm struct {
	MenuID    uint   `json:"menuId" binding:"required"`
	PermTypes []int8 `json:"permTypes"` // 1-查看 2-只读 3-操作
}

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	MenuPerms   []RoleMenuPerm `json:"menuPerms"` // 菜单权限列表
}

// RoleUpdateReq 编辑角色请求
type RoleUpdateReq struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	MenuPerms   []RoleMenuPerm `json:"menuPerms"`
}
