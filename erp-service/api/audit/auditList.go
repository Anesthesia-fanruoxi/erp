package audit

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// AuditList 待审核用户列表(分页)
func AuditList(c *gin.Context) {
	var req model.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := common.DB.Model(&model.User{}).Where("status = ?", common.UserStatusPending)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "查询失败")
		return
	}

	var users []model.User
	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&users).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "查询失败")
		return
	}

	common.Success(c, gin.H{
		"list":     users,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}
