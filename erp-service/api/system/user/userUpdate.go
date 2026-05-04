package user

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserUpdate 编辑用户信息(含单个角色)
func UserUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效的用户ID")
		return
	}

	var req model.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	var u model.User
	if err := common.DB.First(&u, id).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "用户不存在")
		return
	}

	tx := common.DB.Begin()

	updates := map[string]interface{}{
		"real_name": req.RealName,
		"email":     req.Email,
		"phone":     req.Phone,
	}
	if req.Password != "" {
		hash, err := common.HashPassword(req.Password)
		if err != nil {
			tx.Rollback()
			common.Fail(c, common.CodeInternalError, "密码加密失败")
			return
		}
		updates["password_hash"] = hash
	}

	if err := tx.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "更新失败")
		return
	}

	// 更新角色：req.RoleID != nil 时才处理
	// RoleID=0 表示清除角色，RoleID>0 表示设置新角色
	if req.RoleID != nil {
		// 先删除旧角色关联
		if err := tx.Where("user_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
			tx.Rollback()
			common.Fail(c, common.CodeInternalError, "更新角色失败")
			return
		}
		// 分配新角色
		if *req.RoleID > 0 {
			userRole := model.UserRole{UserID: uint(id), RoleID: *req.RoleID}
			if err := tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				common.Fail(c, common.CodeInternalError, "分配角色失败")
				return
			}
		}
	}

	tx.Commit()
	common.SuccessWithMessage(c, "更新成功")
}
