package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/metadata"
)

func ResponseOK(c *gin.Context, msg interface{}) {
	c.JSON(http.StatusOK, metadata.Response{
		BaseResponse: metadata.BaseResponse{
			Code:   0,
			Status: "success",
		},
		Info: msg,
	})
}

func ResponseError(c *gin.Context, msg interface{}, errcode int64) {
	c.JSON(http.StatusOK, metadata.Response{
		BaseResponse: metadata.BaseResponse{
			Code:   errcode,
			Status: "failure",
		},
		Info: msg,
	})
}
