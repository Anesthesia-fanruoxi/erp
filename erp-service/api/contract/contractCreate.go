package contract

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

func ContractCreate(c *gin.Context) {
	var req model.ContractCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	info := common.GetCurrentUser(c)
	createdBy := uint(0)
	if info != nil {
		createdBy = info.UserID
	}

	tx := common.DB.Begin()

	ct := model.Contract{
		ProjectName:  req.ProjectName,
		OrderNo:      req.OrderNo,
		OrderDate:    req.OrderDate,
		FromCompany:  req.FromCompany,
		ToCompany:    req.ToCompany,
		Buyer:        req.Buyer,
		Attn:         req.Attn,
		BuyerEmail:   req.BuyerEmail,
		BuyerTel:     req.BuyerTel,
		AttnTel:      req.AttnTel,
		TotalAmount:  req.TotalAmount,
		DeliveryAddr: req.DeliveryAddr,
		Remark:       req.Remark,
		CreatedBy:    createdBy,
	}
	if err := tx.Create(&ct).Error; err != nil {
		tx.Rollback()
		common.Fail(c, common.CodeInternalError, "创建合同失败: "+err.Error())
		return
	}

	if len(req.Items) > 0 {
		items := make([]model.ContractItem, 0, len(req.Items))
		for _, it := range req.Items {
			items = append(items, model.ContractItem{
				ContractID: ct.ID,
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
			common.Fail(c, common.CodeInternalError, "创建明细失败: "+err.Error())
			return
		}
	}

	tx.Commit()
	common.Success(c, gin.H{"id": ct.ID})
}
