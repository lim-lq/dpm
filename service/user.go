package service

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/core/config"
	"github.com/lim-lq/dpm/errcode"
	"github.com/lim-lq/dpm/models"
	"github.com/lim-lq/dpm/utils"
)

type loginInfo struct {
	Username string
	Password string
}

func UserInfo(c *gin.Context) {
	rediscli := core.GetRedisClient()
	sessionid, _ := c.Cookie(config.SessionName)
	userInfoStr, _ := rediscli.Get(sessionid)
	userInfo := models.Account{}
	err := json.Unmarshal([]byte(userInfoStr), &userInfo)
	if err != nil {
		utils.ResponseError(c, "Get user info failure", 1)
	} else {
		utils.ResponseOK(c, userInfo)
	}
}

func Login(c *gin.Context) {
	logininfo := loginInfo{}
	err := c.BindJSON(&logininfo)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("parameter error - %v", err), 1)
		return
	}
	user := models.AccountManager()
	err = user.SearchByName(c, logininfo.Username)
	if err != nil {
		utils.ResponseError(c, "用户名或密码错误", errcode.LOGIN_INFO_ERROR)
		return
	}
	passPlain, err := core.GetRsaClient().Decrypt(logininfo.Password)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("rsa解密密码失败 - %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	// 对密码做md5加密比较

	if utils.StringToMd5String(passPlain) != user.Password {
		utils.ResponseError(c, "用户名或密码错误", errcode.LOGIN_INFO_ERROR)
		return
	}
	// sessionid, _ := c.Cookie(config.SessionName)
	// rediscli := core.GetRedisClient()
	// infoBytes, err := json.Marshal(&logininfo)
	// if err != nil {
	// 	utils.ResponseError(c, fmt.Sprintf("login error - %v", err), 1)
	// 	return
	// }
	// err = rediscli.Set(sessionid, string(infoBytes), 3600*time.Second)
	// if err != nil {
	// 	utils.ResponseError(c, fmt.Sprintf("login error - %v", err), 1)
	// 	return
	// }
	/* 生成jwt token */
	token, err := core.NewToken(user.Username)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("生成token失败 - %v", err), errcode.GEN_TOKEN_FAILURE)
		return
	}
	utils.ResponseOK(c, map[string]string{"token": token})
}
