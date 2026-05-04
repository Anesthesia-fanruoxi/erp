package user

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserStatus 启用/禁用用户
func UserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效的用户ID")
		return
	}

	var req model.UserStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 更新状态
	if err := common.DB.Model(&model.User{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		common.LogError("更新用户状态失败 userID=%d status=%d err=%v", id, req.Status, err)
		common.Fail(c, common.CodeInternalError, "更新状态失败: "+err.Error())
		return
	}

	// 禁用时踢人
	if req.Status == common.UserStatusDisabled {
		_ = common.DeleteUserTokens(uint(id))
	}

	common.SuccessWithMessage(c, "操作成功")
}
