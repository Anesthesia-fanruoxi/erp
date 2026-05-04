package auth

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// Profile 获取当前用户完整个人信息
func Profile(c *gin.Context) {
	tokenInfo := common.GetCurrentUser(c)
	if tokenInfo == nil {
		common.Fail(c, common.CodeUnauthorized, "未认证")
		return
	}

	var user model.User
	if err := common.DB.First(&user, tokenInfo.UserID).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "用户不存在")
		return
	}

	common.Success(c, gin.H{
		"id":        user.ID,
		"userName":  user.Username,
		"realName":  user.RealName,
		"email":     user.Email,
		"phone":     user.Phone,
		"status":    user.Status,
		"createdAt": user.CreatedAt,
	})
}

// ProfileUpdate 修改当前用户个人信息
func ProfileUpdate(c *gin.Context) {
	tokenInfo := common.GetCurrentUser(c)
	if tokenInfo == nil {
		common.Fail(c, common.CodeUnauthorized, "未认证")
		return
	}

	var req model.ProfileUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{
		"real_name": req.RealName,
		"email":     req.Email,
		"phone":     req.Phone,
	}

	// 如果传了密码则更新密码
	if req.Password != "" {
		hash, err := common.HashPassword(req.Password)
		if err != nil {
			common.Fail(c, common.CodeInternalError, "密码加密失败")
			return
		}
		updates["password_hash"] = hash
	}

	if err := common.DB.Model(&model.User{}).Where("id = ?", tokenInfo.UserID).Updates(updates).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "更新失败")
		return
	}

	common.SuccessWithMessage(c, "更新成功")
}
