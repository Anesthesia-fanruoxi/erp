package auth

import (
	"erp-service/common"
	system "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// Login 账号密码登录
func Login(c *gin.Context) {
	var req system.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	var user system.User
	if err := common.DB.Where("username = ?", req.UserName).First(&user).Error; err != nil {
		common.Fail(c, common.CodeBadRequest, "用户名或密码错误")
		return
	}

	if !common.CheckPassword(req.Password, user.PasswordHash) {
		common.Fail(c, common.CodeBadRequest, "用户名或密码错误")
		return
	}

	if user.Status == common.UserStatusPending {
		common.Fail(c, common.CodeForbidden, "账号待审核")
		return
	}
	if user.Status == common.UserStatusDisabled {
		common.Fail(c, common.CodeForbidden, "账号已禁用")
		return
	}

	token, err := buildAndGenerateToken(&user)
	if err != nil {
		common.Fail(c, common.CodeInternalError, "登录失败: "+err.Error())
		return
	}

	common.Success(c, system.LoginResp{Token: token})
}

// buildAndGenerateToken 查询用户角色并生成Token（权限从Redis角色缓存读取）
func buildAndGenerateToken(user *system.User) (string, error) {
	// 查询用户角色
	var roles []system.Role
	common.DB.Table(common.TableSysRole).
		Joins("JOIN "+common.TableSysUserRole+" ON "+common.TableSysUserRole+".role_id = "+common.TableSysRole+".id").
		Where(common.TableSysUserRole+".user_id = ?", user.ID).
		Find(&roles)

	roleIDs := make([]uint, 0, len(roles))
	roleNames := make([]string, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
		roleNames = append(roleNames, r.Name)
	}

	// Token 只存基础信息 + 角色，权限从 Redis 角色缓存动态读取
	tokenInfo := &common.TokenInfo{
		UserID:   user.ID,
		UserName: user.Username,
		RealName: user.RealName,
		RoleIDs:  roleIDs,
		Roles:    roleNames,
	}

	return common.GenerateToken(tokenInfo)
}

// DeviceLogin 复用此函数（设备免登录）
func buildTokenForUser(user *system.User) (string, error) {
	return buildAndGenerateToken(user)
}
