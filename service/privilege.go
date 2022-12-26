package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/errcode"
	"github.com/lim-lq/dpm/metadata"
	"github.com/lim-lq/dpm/models"
	"github.com/lim-lq/dpm/utils"
)

// 注册权限
func RegisterAction(c *gin.Context) {
	actions := []metadata.MapStr{}

	err := c.ShouldBind(&actions)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("解析参数失败: %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	action := models.PrivilegeActionModel{}
	err = action.Update(c, actions)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("更新权限失败: %v", err), errcode.UPDATE_PRIVILEGE_ERROR)
	} else {
		utils.ResponseOK(c, "Success")
	}
}
