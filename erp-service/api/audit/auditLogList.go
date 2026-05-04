package audit

import (
	"time"

	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// AuditLogList 审计日志列表(分页+多条件筛选)
func AuditLogList(c *gin.Context) {
	var req model.AuditLogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	db := common.DB.Model(&model.AuditLog{})

	// 按操作人筛选
	if req.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+req.UserName+"%")
	}

	// 按操作描述筛选
	if req.Action != "" {
		db = db.Where("action LIKE ?", "%"+req.Action+"%")
	}

	// 按时间范围筛选
	if req.StartTime > 0 {
		db = db.Where("created_at >= ?", time.UnixMilli(req.StartTime))
	}
	if req.EndTime > 0 {
		db = db.Where("created_at <= ?", time.UnixMilli(req.EndTime))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "查询失败")
		return
	}

	var list []model.AuditLog
	offset := (req.Page - 1) * req.PageSize
	if err := db.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&list).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "查询失败")
		return
	}

	common.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}
