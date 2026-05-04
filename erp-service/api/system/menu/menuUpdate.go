package menu

import (
	"erp-service/common"
	model "erp-service/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuUpdate 编辑菜单/按钮/API
func MenuUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Fail(c, common.CodeBadRequest, "ID 无效")
		return
	}

	var req model.MenuUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, common.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 检查记录存在
	var menu model.Menu
	if err := common.DB.First(&menu, id).Error; err != nil {
		common.Fail(c, common.CodeNotFound, "记录不存在")
		return
	}

	// code 唯一性(排除自己)
	if req.Code != "" && req.Code != menu.Code {
		var count int64
		common.DB.Model(&model.Menu{}).Where("code = ? AND id <> ?", req.Code, id).Count(&count)
		if count > 0 {
			common.Fail(c, common.CodeBadRequest, "权限编码已存在")
			return
		}
	}

	// 不允许把父节点设为自己或自己的子孙(简化: 禁止自指)
	if req.ParentID == uint(id) {
		common.Fail(c, common.CodeBadRequest, "父节点不能为自己")
		return
	}
	if req.ParentID != 0 {
		var parent model.Menu
		if err := common.DB.First(&parent, req.ParentID).Error; err != nil {
			common.Fail(c, common.CodeBadRequest, "父节点不存在")
			return
		}
	}

	updates := map[string]interface{}{
		"code":      req.Code,
		"name":      req.Name,
		"type":      req.Type,
		"parent_id": req.ParentID,
		"path":      req.Path,
		"icon":      req.Icon,
		"sort":      req.Sort,
		"visible":   req.Visible,
	}
	if err := common.DB.Model(&model.Menu{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		common.Fail(c, common.CodeInternalError, "更新失败")
		return
	}
	// 菜单变更后刷新所有角色权限缓存
	go refreshAllRoles()
	common.SuccessWithMessage(c, "更新成功")
}
