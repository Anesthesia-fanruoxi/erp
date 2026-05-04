package common

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// TokenInfo Redis中存储的用户基础信息（不含权限和菜单）
type TokenInfo struct {
	UserID   uint     `json:"userId"`
	UserName string   `json:"userName"`
	RealName string   `json:"realName"`
	RoleIDs  []uint   `json:"roleIds"` // 角色ID列表
	Roles    []string `json:"roles"`   // 角色名列表
}

// MenuNode 菜单节点
type MenuNode struct {
	ID       uint       `json:"id"`
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Icon     string     `json:"icon"`
	Sort     int        `json:"sort"`
	Children []MenuNode `json:"children,omitempty"`
}

// GenerateToken 生成UUID Token并存入Redis
func GenerateToken(info *TokenInfo) (string, error) {
	token := uuid.New().String()

	jsonData, err := json.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("序列化TokenInfo失败: %w", err)
	}

	ttl := RedisTokenTTL()

	// 存储 token:<uuid> -> JSON
	if err := RDB.Set(Ctx, RedisTokenKey(token), string(jsonData), ttl).Err(); err != nil {
		return "", fmt.Errorf("存储Token到Redis失败: %w", err)
	}

	// 添加到用户Token集合
	if err := RDB.SAdd(Ctx, RedisUserTokensKey(info.UserID), token).Err(); err != nil {
		return "", fmt.Errorf("添加Token到用户集合失败: %w", err)
	}

	return token, nil
}

// GetTokenInfo 从Redis获取用户信息并刷新TTL(滑动过期)
func GetTokenInfo(token string) (*TokenInfo, error) {
	jsonData, err := RDB.Get(Ctx, RedisTokenKey(token)).Result()
	if err != nil {
		return nil, fmt.Errorf("Token不存在或已过期")
	}

	var info TokenInfo
	if err := json.Unmarshal([]byte(jsonData), &info); err != nil {
		return nil, fmt.Errorf("反序列化TokenInfo失败: %w", err)
	}

	// 刷新TTL(滑动过期)
	RDB.Expire(Ctx, RedisTokenKey(token), RedisTokenTTL())

	return &info, nil
}

// DeleteToken 删除指定Token
func DeleteToken(token string) error {
	// 先获取Token信息，用于清理用户Token集合
	info, err := GetTokenInfo(token)
	if err == nil && info != nil {
		RDB.SRem(Ctx, RedisUserTokensKey(info.UserID), token)
	}

	return RDB.Del(Ctx, RedisTokenKey(token)).Err()
}

// DeleteUserTokens 删除用户所有Token(踢人)
func DeleteUserTokens(userId uint) error {
	// 获取用户所有Token
	tokens, err := RDB.SMembers(Ctx, RedisUserTokensKey(userId)).Result()
	if err != nil {
		return fmt.Errorf("获取用户Token列表失败: %w", err)
	}

	// 批量删除所有Token
	for _, token := range tokens {
		RDB.Del(Ctx, RedisTokenKey(token))
	}

	// 删除用户Token集合
	return RDB.Del(Ctx, RedisUserTokensKey(userId)).Err()
}
