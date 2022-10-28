package service

import (
	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/utils"
)

func Ping(c *gin.Context) {
	utils.ResponseOK(c, "pong")
}
