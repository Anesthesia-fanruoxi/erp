package user

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// UserCreate 管理员创建用户(直接激活+分配单个角色)
func UserCreate(c *gin.Context) {
	var req model.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var existing model.User
	result := common.DB.Unscoped().Where("username = ?", req.UserName).First(&existing)
	if result.Error == nil {
		if existing.DeletedAt.Valid {
			// 已软删除的同名用户，物理删除旧记录及其角色关联
			common.DB.Unscoped().Where("user_id = ?", existing.ID).Delete(&model.UserRole{})
			common.DB.Unscoped().Delete(&existing)
		} else {
			common.Fail(c, common.CodeBadRequest, "用户名已存在")
			return
		}
	}

	hash, err := common.HashPassword(req.Password)
	if err != nil {
		common.Fail(c, common.CodeInternalError, "密码加密失败")
		return
	}

	tx := common.DB.Begin()

	u := model.User{
		Username:     req.UserName,
		PasswordHash: hash,
		RealName:     req.RealName,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       common.UserStatusActive,
	}
	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "创建用户失败")
		return
	}

	// 分配单个角色
	if req.RoleID != nil && *req.RoleID > 0 {
		userRole := model.UserRole{UserID: u.ID, RoleID: *req.RoleID}
		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			common.Fail(c, common.CodeInternalError, "分配角色失败")
			return
		}
	}

	tx.Commit()
	common.SuccessWithMessage(c, "创建成功")
}
