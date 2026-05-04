package user

import (
	"erp-service/common"
	model "erp-service/model/system"

	"github.com/gin-gonic/gin"
)

// UserList 用户列表(分页+搜索)
func UserList(c *gin.Context) {
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

	db := common.DB.Model(&model.User{})

	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("username LIKE ? OR real_name LIKE ?", like, like)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

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

	// 批量查询用户角色（每人一个角色）
	userIDs := make([]uint, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}

	type userRoleRow struct {
		UserID   uint   `gorm:"column:user_id"`
		RoleID   uint   `gorm:"column:role_id"`
		RoleName string `gorm:"column:name"`
	}
	var userRoleRows []userRoleRow
	if len(userIDs) > 0 {
		common.DB.Table(common.TableSysUserRole).
			Select(common.TableSysUserRole+".user_id, "+
				common.TableSysUserRole+".role_id, "+
				common.TableSysRole+".name").
			Joins("JOIN "+common.TableSysRole+" ON "+common.TableSysRole+".id = "+common.TableSysUserRole+".role_id").
			Where(common.TableSysUserRole+".user_id IN ?", userIDs).
			Scan(&userRoleRows)
	}

	// 按用户ID索引（取第一条，每人只有一个角色）
	type roleInfo struct {
		ID   uint
		Name string
	}
	roleMap := make(map[uint]roleInfo)
	for _, row := range userRoleRows {
		if _, exists := roleMap[row.UserID]; !exists {
			roleMap[row.UserID] = roleInfo{ID: row.RoleID, Name: row.RoleName}
		}
	}

	// 组装结果
	list := make([]model.UserListItem, 0, len(users))
	for _, u := range users {
		item := model.UserListItem{User: u}
		if r, ok := roleMap[u.ID]; ok {
			item.RoleName = r.Name
			id := r.ID
			item.RoleID = &id
		}
		list = append(list, item)
	}

	common.Success(c, model.UserListResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}
