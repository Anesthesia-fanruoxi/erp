package common

// ================================================================
// 数据库表名统一定义
// 规范: 所有系统表以 sys_ 开头
// 使用方式:
//   - Model 的 TableName() 方法直接返回此常量
//   - 直接使用 Table()/JOIN 的地方也应引用此常量
// 如需改表名, 仅需修改此处一处
// ================================================================

const (
	// 系统用户表
	TableSysUser = "sys_user"
	// 系统角色表
	TableSysRole = "sys_role"
	// 用户-角色关联表
	TableSysUserRole = "sys_user_role"
	// 系统菜单/权限项表 (菜单/按钮/API 三合一)
	TableSysMenu = "sys_menu"
	// 角色-菜单关联表
	TableSysRoleMenu = "sys_role_menu"
	// 设备绑定表
	TableSysDeviceBinding = "sys_device_binding"
	// 操作审计日志表
	TableSysAuditLog = "sys_audit_log"
	// 合同主表
	TableContract = "biz_contract"
	// 合同明细表
	TableContractItem = "biz_contract_item"
)
