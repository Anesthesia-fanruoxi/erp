package audit

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuditApprove 审核通过
func AuditApprove(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效的用户ID")
		return
	}

	var user model.User
	if err := common.DB.First(&user, id).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "用户不存在")
		return
	}
	if user.Status != common.UserStatusPending {
		common.Fail(c, common.CodeBadRequest, "该用户不在待审核状态")
		return
	}

	tx := common.DB.Begin()

	// 激活用户
	if err := tx.Model(&model.User{}).Where("id = ?", id).Update("status", common.UserStatusActive).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "审核操作失败")
		return
	}

	// 分配默认 "user" 角色（先查角色是否存在，再检查是否已绑定）
	var defaultRole model.Role
	if err := common.DB.Where("name = ?", "user").First(&defaultRole).Error; err == nil {
		var count int64
		common.DB.Model(&model.UserRole{}).Where("user_id = ?", id).Count(&count)
		if count == 0 {
			if err := tx.Create(&model.UserRole{
				UserID: uint(id),
				RoleID: defaultRole.ID,
			}).Error; err != nil {
				tx.Rollback()
				common.Fail(c, common.CodeInternalError, "分配角色失败")
				return
			}
		}
	}

	tx.Commit()
	common.SuccessWithMessage(c, "审核通过")
}
