package common

import (
	"fmt"
	"time"

	"erp-service/config"
)

// ================================================================
// Redis Key 统一管理
// 规范:
//   - 所有 Key 前缀集中定义, 避免散落在业务代码
//   - 统一通过构造函数生成最终 Key, 避免拼写错误
//   - 每类 Key 的 TTL(生命周期) 也统一在此定义/派生, 便于全局调整
// 命名约定:
//   - 前缀以冒号 ":" 结尾, 遵循 Redis 约定 (如 "token:xxx")
// ================================================================

// ---- Key 前缀 ----
const (
	// RedisKeyPrefixToken Token 主体存储: token:<uuid> -> TokenInfo JSON
	RedisKeyPrefixToken = "token:"

	// RedisKeyPrefixUserTokens 用户持有的 token 集合: user_tokens:<userId> -> Set<token>
	RedisKeyPrefixUserTokens = "user_tokens:"

	// RedisKeyPrefixRolePerms 角色权限集合: role_perms:<roleId> -> Set<permCode>
	RedisKeyPrefixRolePerms = "role_perms:"

	// RedisKeyPrefixRoleMenus 角色菜单树: role_menus:<roleId> -> JSON
	RedisKeyPrefixRoleMenus = "role_menus:"
)

// ---- Key 构造函数 ----

// RedisTokenKey 构造 token 详情 Key: token:<uuid>
func RedisTokenKey(token string) string {
	return RedisKeyPrefixToken + token
}

// RedisUserTokensKey 构造用户 Token 集合 Key: user_tokens:<userId>
func RedisUserTokensKey(userID uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefixUserTokens, userID)
}

// RedisRolePermsKey 构造角色权限集合 Key: role_perms:<roleId>
func RedisRolePermsKey(roleID uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefixRolePerms, roleID)
}

// RedisRoleMenusKey 构造角色菜单树 Key: role_menus:<roleId>
func RedisRoleMenusKey(roleID uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefixRoleMenus, roleID)
}

// ---- 生命周期 (TTL) ----

// RedisTokenTTL Token 的滑动过期时长
// 取自 config.AppConfig.Token.ExpireTime (秒)
func RedisTokenTTL() time.Duration {
	return time.Duration(config.AppConfig.Token.ExpireTime) * time.Second
}
