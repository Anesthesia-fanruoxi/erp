package menu

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// MenuCreate 创建菜单/按钮/API 权限项
func MenuCreate(c *gin.Context) {
	var req model.MenuCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 校验 code 唯一
	var count int64
	common.DB.Model(&model.Menu{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		common.Fail(c, common.CodeBadRequest, "权限编码已存在")
		return
	}

	// 校验 parentId 合法性
	if req.ParentID != 0 {
		var parent model.Menu
		if err := common.DB.First(&parent, req.ParentID).Error; err != nil {
			common.Fail(c, common.CodeBadRequest, "父节点不存在")
			return
		}
	}

	if req.Visible == 0 {
		req.Visible = 1
	}

	menu := model.Menu{
		Code:     req.Code,
		Name:     req.Name,
		Type:     req.Type,
		ParentID: req.ParentID,
		Path:     req.Path,
		Icon:     req.Icon,
		Sort:     req.Sort,
		Visible:  req.Visible,
	}
	if err := common.DB.Create(&menu).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "创建失败")
		return
	}

	// 自动给 admin 角色分配该菜单的全部权限（查看+只读+操作）
	go autoGrantAdminMenu(menu.ID)

	common.Success(c, menu)
}
