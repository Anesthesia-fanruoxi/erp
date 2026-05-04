package contract

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ContractDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效ID")
		return
	}

	var ct model.Contract
	if err := common.DB.First(&ct, id).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "合同不存在")
		return
	}

	var items []model.ContractItem
	common.DB.Where("contract_id = ?", id).Order("seq ASC").Find(&items)

	common.Success(c, model.ContractDetail{Contract: ct, Items: items})
}
