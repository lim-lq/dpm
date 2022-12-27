package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/metadata"
)

func baseResponse(c *gin.Context, httpStatus int, msg interface{}, code int64) {
	var status string
	if code == 0 {
		status = "success"
	} else {
		status = "failure"
	}
	c.JSON(httpStatus, metadata.Response{
		BaseResponse: metadata.BaseResponse{
			Code:   code,
			Status: status,
		},
		Info: msg,
	})
}

func ResponseOK(c *gin.Context, msg interface{}) {
	baseResponse(c, http.StatusOK, msg, 0)
}

func ResponseList(c *gin.Context, msg metadata.PageData) {
	c.JSON(http.StatusOK, metadata.PageListResponse{
		BaseResponse: metadata.BaseResponse{
			Code:   0,
			Status: "success",
		},
		Info: msg,
	})
}

func ResponseError(c *gin.Context, msg interface{}, errcode int64) {
	baseResponse(c, http.StatusOK, msg, errcode)
}

func UnauthedError(c *gin.Context, msg interface{}) {
	baseResponse(c, http.StatusUnauthorized, msg, -1)
}

func ForbiddenError(c *gin.Context, msg interface{}) {
	baseResponse(c, http.StatusForbidden, msg, -1)
}
