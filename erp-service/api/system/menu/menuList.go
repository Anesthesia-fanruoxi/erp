package menu

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// MenuList 获取全部菜单/权限项列表(平铺, 前端自行组装)
func MenuList(c *gin.Context) {
	var list []model.Menu
	if err := common.DB.Order("parent_id ASC, sort ASC, id ASC").Find(&list).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "查询菜单失败")
		return
	}
	common.Success(c, list)
}
