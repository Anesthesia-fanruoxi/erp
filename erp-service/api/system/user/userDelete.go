package user

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserDelete 删除用户(软删除+踢人)
func UserDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效的用户ID")
		return
	}

	// 软删除用户
	if err := common.DB.Delete(&model.User{}, id).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "删除失败")
		return
	}

	// 踢掉该用户的所有登录
	_ = common.DeleteUserTokens(uint(id))

	common.SuccessWithMessage(c, "删除成功")
}
