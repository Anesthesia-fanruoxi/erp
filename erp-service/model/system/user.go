package system

import (
	"erp-service/common"

	"gorm.io/gorm"
)

// User 用户表模型
type User struct {
	ID           uint             `json:"id" gorm:"primaryKey"`
	Username     string           `json:"userName" gorm:"column:username;uniqueIndex;size:50"`
	PasswordHash string           `json:"-" gorm:"column:password_hash;size:255"`
	RealName     string           `json:"realName" gorm:"column:real_name;size:50"`
	Email        string           `json:"email" gorm:"column:email;size:100"`
	Phone        string           `json:"phone" gorm:"column:phone;size:20"`
	Status       int8             `json:"status" gorm:"column:status;default:0"`
	CreatedAt    common.MilliTime `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    common.MilliTime `json:"updatedAt" gorm:"column:updated_at;->"`
	DeletedAt    gorm.DeletedAt   `json:"-" gorm:"column:deleted_at;index"`
}

func (User) TableName() string {
	return common.TableSysUser
}

// --- 请求 DTO ---

// UserCreateReq 创建用户请求 (每个用户只能绑定一个角色)
type UserCreateReq struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"realName" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleID   *uint  `json:"roleId"` // 单个角色ID，nil 表示不分配
}

// UserUpdateReq 编辑用户请求
type UserUpdateReq struct {
	RealName string `json:"realName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"` // 可选,不传则不修改密码
	RoleID   *uint  `json:"roleId"`   // 可选,nil 表示不修改角色，0 表示清除角色
}

// ProfileUpdateReq 个人信息修改请求
type ProfileUpdateReq struct {
	RealName string `json:"realName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// UserStatusReq 启用/禁用请求
type UserStatusReq struct {
	Status int8 `json:"status" binding:"required,oneof=1 2"`
}

// UserListReq 用户列表查询请求
type UserListReq struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   *int8  `form:"status" json:"status"`
}

// --- 响应 DTO ---

// UserListItem 用户列表项(含单个角色)
type UserListItem struct {
	User
	RoleName string `json:"roleName"` // 角色名，无角色时为空字符串
	RoleID   *uint  `json:"roleId"`   // 角色ID，无角色时为 null
}

// UserListResp 用户列表响应
type UserListResp struct {
	List     []UserListItem `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}
