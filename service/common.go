package service

import (
	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/utils"
)

func GetPublicKey(c *gin.Context) {
	rsaCli := core.GetRsaClient()
	utils.ResponseOK(c, rsaCli.GetPublicKeyStr())
}
