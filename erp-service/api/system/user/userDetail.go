package user

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserDetail 用户详情(含单个角色)
func UserDetail(c *gin.Context) {
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

	// 查询用户角色（只取一条）
	var role model.Role
	err = common.DB.Table(common.TableSysRole).
		Joins("JOIN "+common.TableSysUserRole+" ON "+common.TableSysUserRole+".role_id = "+common.TableSysRole+".id").
		Where(common.TableSysUserRole+".user_id = ?", id).
		First(&role).Error

	resp := gin.H{"user": user}
	if err == nil {
		resp["role"] = role
	} else {
		resp["role"] = nil
	}

	common.Success(c, resp)
}
