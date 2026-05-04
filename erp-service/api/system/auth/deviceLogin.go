package auth

import (
	"erp-service/common"
	system "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// DeviceLogin 设备免登录
func DeviceLogin(c *gin.Context) {
	var req system.DeviceLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 查询所有设备绑定记录
	var bindings []system.DeviceBinding
	if err := common.DB.Find(&bindings).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "查询设备绑定失败")
		return
	}

	// 遍历比对机器码 (bcrypt 需逐个比对)
	var matchedUserID uint
	for _, b := range bindings {
		if common.CheckMachineCode(req.MachineCode, b.MachineHash) {
			matchedUserID = b.UserID
			break
		}
	}

	if matchedUserID == 0 {
		common.Fail(c, common.CodeBadRequest, "设备未绑定")
		return
	}

	// 查询用户信息
	var user system.User
	if err := common.DB.First(&user, matchedUserID).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "查询用户失败")
		return
	}

	// 检查用户状态
	if user.Status == common.UserStatusPending {
		common.Fail(c, common.CodeForbidden, "账号待审核")
		return
	}
	if user.Status == common.UserStatusDisabled {
		common.Fail(c, common.CodeForbidden, "账号已禁用")
		return
	}

	// 复用公共函数生成Token
	token, err := buildAndGenerateToken(&user)
	if err != nil {
		common.Fail(c, common.CodeInternalError, "登录失败: "+err.Error())
		return
	}

	common.Success(c, system.LoginResp{Token: token})
}
