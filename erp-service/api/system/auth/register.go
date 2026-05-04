package auth

import (
	"erp-service/common"
	system "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req system.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var count int64
	common.DB.Model(&system.User{}).Where("username = ?", req.UserName).Count(&count)
	if count > 0 {
		common.Fail(c, common.CodeBadRequest, "用户名已存在")
		return
	}

	// 加密密码
	hash, err := common.HashPassword(req.Password)
	if err != nil {
		common.Fail(c, common.CodeInternalError, "密码加密失败")
		return
	}

	// 创建用户 status=0 待审核
	user := system.User{
		Username:     req.UserName,
		PasswordHash: hash,
		RealName:     req.RealName,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       common.UserStatusPending,
	}
	if err := common.DB.Create(&user).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "创建用户失败")
		return
	}

	common.SuccessWithMessage(c, "注册成功，请等待审核")
}
