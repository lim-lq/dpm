package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/core/config"
	"github.com/lim-lq/dpm/metadata"
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
		c.JSON(http.StatusOK, metadata.Response{
			BaseResponse: metadata.BaseResponse{
				Code:   1,
				Status: "failure",
			},
			Info: "Get user info failure",
		})
	} else {
		c.JSON(http.StatusOK, metadata.Response{
			BaseResponse: metadata.BaseResponse{
				Code:   0,
				Status: "success",
			},
			Info: userInfo,
		})
	}
}

func Login(c *gin.Context) {
	logininfo := loginInfo{}
	err := c.BindJSON(&logininfo)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("parameter error - %v", err), 1)
		return
	}
	if logininfo.Username == "admin" && logininfo.Password == "123456" {
		sessionid, _ := c.Cookie(config.SessionName)
		rediscli := core.GetRedisClient()
		infoBytes, err := json.Marshal(&logininfo)
		if err != nil {
			utils.ResponseError(c, fmt.Sprintf("login error - %v", err), 1)
			return
		}
		err = rediscli.Set(sessionid, string(infoBytes), 3600*time.Second)
		if err != nil {
			utils.ResponseError(c, fmt.Sprintf("login error - %v", err), 1)
			return
		}
	} else {
		utils.ResponseError(c, "username or password error", 1)
		return
	}
	utils.ResponseOK(c, "Login success")
}
