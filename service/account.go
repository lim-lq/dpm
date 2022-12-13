package service

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/errcode"
	"github.com/lim-lq/dpm/metadata"
	"github.com/lim-lq/dpm/models"
	"github.com/lim-lq/dpm/utils"
)

func QueryAccountList(c *gin.Context) {
	cond := metadata.Condition{}
	err := c.ShouldBind(&cond)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("解析参数失败: %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	account := models.AccountManager()
	accList, err := account.GetList(c, &cond)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("查询账号失败: %v", err), errcode.GET_ACCOUNT_ERROR)
		return
	}
	utils.ResponseOK(c, accList)
}

func CreateAccount(c *gin.Context) {
	account := models.AccountManager()
	err := c.ShouldBind(account)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("解析参数失败: %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	// 解析密码
	cipher, err := core.GetRsaClient().Decrypt(account.Password)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("创建账号失败: %v", err), errcode.CREATE_ACCOUNT_ERROR)
		return
	}
	account.Password = utils.StringToMd5String(cipher)
	err = account.Create(c)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("创建账号失败: %v", err), errcode.CREATE_ACCOUNT_ERROR)
		return
	}
	utils.ResponseOK(c, "Success")
}

func UpdateAccount(c *gin.Context) {
	accID, err := strconv.ParseInt(c.Param("accountid"), 10, 64)
	if err != nil {
		utils.ResponseError(c, "参数错误", errcode.PARSE_PARAMETER_ERROR)
		return
	}
	accMgr := models.AccountManager()
	cond := metadata.Condition{
		Filters: metadata.Filters{"id": accID},
	}
	accList, err := accMgr.GetList(c, &cond)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("获取账号失败: %v", err), errcode.GET_ACCOUNT_ERROR)
		return
	}
	if len(accList) == 0 {
		utils.ResponseError(c, "要编辑的账号不存在", errcode.ACCOUNT_NOT_FOUND)
		return
	}
	oldAcc := accList[0]
	newAcc := metadata.MapStr{}
	err = c.ShouldBind(&newAcc)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("解析参数失败: %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	err = oldAcc.Update(c, newAcc)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("更新账号失败: %v", err), errcode.UPDATE_ACCOUNT_ERROR)
		return
	}
	utils.ResponseOK(c, "Success")
}

func DeleteAccount(c *gin.Context) {
	// account := models.AccountManager()
	filters := metadata.Filters{}
	err := c.ShouldBind(&filters)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("解析参数失败: %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	username, ok := filters["username"]
	if !ok {
		utils.ResponseError(c, "请传入参数username", errcode.PARAMETER_MISSING)
		return
	}
	if username == "admin" {
		utils.ResponseError(c, "admin用户不允许删除", errcode.PARSE_PARAMETER_ERROR)
		return
	}
	cond := metadata.Condition{Filters: filters}
	err = models.AccountManager().Delete(c, &cond)
	if err != nil {
		utils.ResponseError(c, "删除失败", errcode.DELETE_ACCOUNT_ERROR)
		return
	}
	utils.ResponseOK(c, "Success")
}

func ChangeAccountPassword(c *gin.Context) {
	accId, err := strconv.ParseInt(c.Param("accountid"), 10, 64)
	if err != nil {
		utils.ResponseError(c, "Account id is wrong format", errcode.PARSE_PARAMETER_ERROR)
		return
	}
	acc := models.Account{BaseModel: models.BaseModel{Id: accId}}
	cipher := map[string]string{}
	err = c.ShouldBind(&cipher)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("Parse cipher body error: %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	// rsa 解密密码
	plain, err := core.GetRsaClient().Decrypt(cipher["cipher"])
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("Decrypt password error: %v", err), errcode.DECRYPT_FAILURE)
		return
	}
	// 密码加密存储
	acc.ChangePassword(c, utils.StringToMd5String(plain))
	utils.ResponseOK(c, "OK")
}
