package auth

import (
	"erp-service/common"
	system "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// DeviceUnbind 解除设备绑定
func DeviceUnbind(c *gin.Context) {
	tokenInfo := common.GetCurrentUser(c)
	if tokenInfo == nil {
		common.Fail(c, common.CodeUnauthorized, "未认证")
		return
	}

	result := common.DB.Where("user_id = ?", tokenInfo.UserID).Delete(&system.DeviceBinding{})
	if result.Error != nil {
		common.Fail(c, common.CodeInternalError, "解绑失败")
		return
	}

	if result.RowsAffected == 0 {
		common.Fail(c, common.CodeBadRequest, "未绑定设备")
		return
	}

	common.SuccessWithMessage(c, "设备解绑成功")
}
