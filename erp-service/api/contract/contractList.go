package contract

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

func ContractList(c *gin.Context) {
	var req model.ContractListReq
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

	db := common.DB.Model(&model.Contract{})
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("project_name LIKE ? OR from_company LIKE ? OR to_company LIKE ?", like, like, like)
	}
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}

	var total int64
	db.Count(&total)

	var list []model.Contract
	offset := (req.Page - 1) * req.PageSize
	db.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&list)

	common.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}
