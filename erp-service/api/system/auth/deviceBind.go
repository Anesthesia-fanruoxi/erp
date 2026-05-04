package auth

import (
	"erp-service/common"
	system "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// DeviceBind 绑定设备
func DeviceBind(c *gin.Context) {
	var req system.DeviceBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	tokenInfo := common.GetCurrentUser(c)
	if tokenInfo == nil {
		common.Fail(c, common.CodeUnauthorized, "未认证")
		return
	}

	// 检查该用户是否已绑定设备
	var count int64
	common.DB.Model(&system.DeviceBinding{}).Where("user_id = ?", tokenInfo.UserID).Count(&count)
	if count > 0 {
		common.Fail(c, common.CodeBadRequest, "该用户已绑定设备，请先解绑")
		return
	}

	// 加密机器码
	hash, err := common.HashMachineCode(req.MachineCode)
	if err != nil {
		common.Fail(c, common.CodeInternalError, "机器码加密失败")
		return
	}

	// 创建绑定记录
	binding := system.DeviceBinding{
		UserID:      tokenInfo.UserID,
		MachineHash: hash,
		DeviceName:  req.DeviceName,
	}
	if err := common.DB.Create(&binding).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "设备绑定失败")
		return
	}

	common.SuccessWithMessage(c, "设备绑定成功")
}
