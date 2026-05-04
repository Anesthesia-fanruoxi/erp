package contract

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ContractUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "无效ID")
		return
	}

	var req model.ContractUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	var ct model.Contract
	if err := common.DB.First(&ct, id).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "合同不存在")
		return
	}

	tx := common.DB.Begin()

	updates := map[string]interface{}{
		"project_name":  req.ProjectName,
		"order_no":      req.OrderNo,
		"order_date":    req.OrderDate,
		"from_company":  req.FromCompany,
		"to_company":    req.ToCompany,
		"buyer":         req.Buyer,
		"attn":          req.Attn,
		"buyer_email":   req.BuyerEmail,
		"buyer_tel":     req.BuyerTel,
		"attn_tel":      req.AttnTel,
		"total_amount":  req.TotalAmount,
		"delivery_addr": req.DeliveryAddr,
		"remark":        req.Remark,
	}
	if err := tx.Model(&model.Contract{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "更新失败: "+err.Error())
		return
	}

	// 全量替换明细
	if err := tx.Where("contract_id = ?", id).Delete(&model.ContractItem{}).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "更新明细失败")
		return
	}
	if len(req.Items) > 0 {
		items := make([]model.ContractItem, 0, len(req.Items))
		for _, it := range req.Items {
			items = append(items, model.ContractItem{
				ContractID: uint(id),
				Seq:        it.Seq,
				Name:       it.Name,
				Spec:       it.Spec,
				Brand:      it.Brand,
				Qty:        it.Qty,
				Unit:       it.Unit,
				UnitPrice:  it.UnitPrice,
				Amount:     it.Amount,
				Operator:   it.Operator,
				Location:   it.Location,
				Remark:     it.Remark,
			})
		}
		if err := tx.Create(&items).Error; err != nil {
			tx.Rollback()
			common.Fail(c, common.CodeInternalError, "创建明细失败")
			return
		}
	}

	tx.Commit()
	common.SuccessWithMessage(c, "更新成功")
}
