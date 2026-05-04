package router

import (
	audit2 "erp-service/api/audit"
	"erp-service/api/contract"
	"erp-service/api/system/auth"
	"erp-service/api/system/menu"
	"erp-service/api/system/role"
	"erp-service/api/system/user"
	"erp-service/common"

	"github.com/gin-gonic/gin"
)

// SetupRouter 注册所有路由
// 权限码规则:
//
//	code:r → 只读接口 (GET)
//	code:w → 操作接口 (POST/PUT/DELETE)
//	""     → 只需登录,无需特定权限
func SetupRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	// ----------------------------------------------------------------
	// 公开路由 (不需要认证)
	// ----------------------------------------------------------------
	authPublic := v1.Group("/auth")
	{
		authPublic.POST("/login", auth.Login)
		authPublic.POST("/device/login", auth.DeviceLogin)
	}

	// ----------------------------------------------------------------
	// 需要登录的路由
	// ----------------------------------------------------------------
	authorized := v1.Group("")
	authorized.Use(common.AuthMiddleware())
	{
		authGroup := authorized.Group("/auth")
		{
			authGroup.POST("/logout", auth.Logout)
			authGroup.GET("/current", auth.Current)
			authGroup.GET("/menus", auth.Menus)
			authGroup.GET("/profile", auth.Profile)
			authGroup.PUT("/profile",
				common.RouteConfig("", true, "修改个人信息"),
				auth.ProfileUpdate)
			authGroup.POST("/device/bind",
				common.RouteConfig("", true, "绑定设备"),
				auth.DeviceBind)
			authGroup.DELETE("/device/unbind",
				common.RouteConfig("", true, "解绑设备"),
				auth.DeviceUnbind)
		}
	}

	// ----------------------------------------------------------------
	// 需要登录 + 权限码控制的路由（去掉 AdminMiddleware，由 RouteConfig 细粒度控制）
	// ----------------------------------------------------------------
	admin := authorized.Group("")
	{
		// 用户管理
		admin.GET("/users",
			common.RouteConfig("system:user:r", false, ""),
			user.UserList)
		admin.GET("/users/:id",
			common.RouteConfig("system:user:r", false, ""),
			user.UserDetail)
		admin.POST("/users",
			common.RouteConfig("system:user:w", true, "创建用户"),
			user.UserCreate)
		admin.PUT("/users/:id",
			common.RouteConfig("system:user:w", true, "编辑用户"),
			user.UserUpdate)
		admin.DELETE("/users/:id",
			common.RouteConfig("system:user:w", true, "删除用户"),
			user.UserDelete)
		admin.PUT("/users/:id/status",
			common.RouteConfig("system:user:w", true, "修改用户状态"),
			user.UserStatus)

		// 角色管理
		admin.GET("/roles",
			common.RouteConfig("system:role:r", false, ""),
			role.RoleList)
		admin.POST("/roles",
			common.RouteConfig("system:role:w", true, "创建角色"),
			role.RoleCreate)
		admin.PUT("/roles/:id",
			common.RouteConfig("system:role:w", true, "编辑角色"),
			role.RoleUpdate)
		admin.DELETE("/roles/:id",
			common.RouteConfig("system:role:w", true, "删除角色"),
			role.RoleDelete)

		// 菜单管理
		admin.GET("/menus",
			common.RouteConfig("system:menu:r", false, ""),
			menu.MenuList)
		admin.POST("/menus",
			common.RouteConfig("system:menu:w", true, "创建菜单"),
			menu.MenuCreate)
		admin.PUT("/menus/:id",
			common.RouteConfig("system:menu:w", true, "编辑菜单"),
			menu.MenuUpdate)
		admin.DELETE("/menus/:id",
			common.RouteConfig("system:menu:w", true, "删除菜单"),
			menu.MenuDelete)

		// 注册审核
		admin.GET("/registrations",
			common.RouteConfig("system:audit:r", false, ""),
			audit2.AuditList)
		admin.PUT("/registrations/:id/approve",
			common.RouteConfig("system:audit:w", true, "审核通过"),
			audit2.AuditApprove)
		admin.PUT("/registrations/:id/reject",
			common.RouteConfig("system:audit:w", true, "审核拒绝"),
			audit2.AuditReject)

		// 审计日志
		admin.GET("/audit/logs",
			common.RouteConfig("audit:log:r", false, ""),
			audit2.AuditLogList)

		// 合同管理
		admin.GET("/contracts",
			common.RouteConfig("business:contract:r", false, ""),
			contract.ContractList)
		admin.GET("/contracts/:id",
			common.RouteConfig("business:contract:r", false, ""),
			contract.ContractDetail)
		admin.POST("/contracts",
			common.RouteConfig("business:contract:w", true, "创建合同"),
			contract.ContractCreate)
		admin.PUT("/contracts/:id",
			common.RouteConfig("business:contract:w", true, "更新合同"),
			contract.ContractUpdate)
		admin.DELETE("/contracts/:id",
			common.RouteConfig("business:contract:w", true, "删除合同"),
			contract.ContractDelete)
	}
}
