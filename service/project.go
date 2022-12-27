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
	result, err := models.ProjectManager().GetPageList(c, cond)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("Get project list error - %v", err), 1)
	} else {
		utils.ResponseList(c, *result)
	}
}

func CreateProject(c *gin.Context) {
	project := models.ProjectModel{}
	err := c.ShouldBind(&project)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("Parse parameter error - %v", err), errcode.PARSE_PARAMETER_ERROR)
		return
	}
	err = project.Create(c)
	if err != nil {
		utils.ResponseError(c, fmt.Sprintf("Create project error - %v", err), errcode.CREATE_PROJECT_ERROR)
		return
	}
	utils.ResponseOK(c, "Success")
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
