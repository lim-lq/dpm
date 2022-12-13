package service

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/errcode"
	"github.com/lim-lq/dpm/metadata"
	"github.com/lim-lq/dpm/models"
	"github.com/lim-lq/dpm/utils"
)

func ProjectList(c *gin.Context) {
	cond, err := utils.DealListCondition(c)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("Post data parse error - %v", err), 1)
		return
	}
	projects, err := models.ProjectManager().GetList(c, cond)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("Get project list error - %v", err), 1)
	} else {
		utils.ResponseOK(c, projects)
	}
}

func ProjectDetail(c *gin.Context) {
	projectid, err := strconv.ParseInt(c.Param("projectid"), 10, 64)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("parse project id error - %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	cond := metadata.Condition{
		Page: metadata.Page{
			Limit:  1,
			Offset: 0,
		},
		Filters: map[string]interface{}{
			"id": projectid,
		},
	}
	projectInst := models.ProjectManager()
	err = projectInst.GetDetail(c, &cond)

	if err != nil {
		utils.ResponseError(c, fmt.Sprintf(
			"Do not found the project with specify project id %d error - %v", projectid, err), errcode.PROJECT_NOT_FOUND)
		return
	}
	utils.ResponseOK(c, projectInst)
}