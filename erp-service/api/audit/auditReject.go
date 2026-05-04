package audit

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuditReject 审核拒绝
func AuditReject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效的用户ID")
		return
	}

	// 查询用户并确认待审核状态
	var user model.User
	if err := common.DB.First(&user, id).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "用户不存在")
		return
	}
	if user.Status != common.UserStatusPending {
		common.Fail(c, common.CodeBadRequest, "该用户不在待审核状态")
		return
	}

	// 更新为禁用状态(拒绝)
	if err := common.DB.Model(&model.User{}).Where("id = ?", id).Update("status", common.UserStatusDisabled).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "审核操作失败")
		return
	}

	common.SuccessWithMessage(c, "已拒绝")
}
