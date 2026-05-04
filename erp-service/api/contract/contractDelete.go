package contract

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ContractDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效ID")
		return
	}

	tx := common.DB.Begin()
	if err := tx.Where("contract_id = ?", id).Delete(&model.ContractItem{}).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "删除明细失败")
		return
	}
	if err := tx.Delete(&model.Contract{}, id).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "删除合同失败")
		return
	}
	tx.Commit()
	common.SuccessWithMessage(c, "删除成功")
}
