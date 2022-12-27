package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/metadata"
)

var timeFileds = map[string]bool{
	"createTime": true,
	"updateTime": true,
}

func DealListCondition(c *gin.Context) (*metadata.Condition, error) {
	request := new(metadata.PageRequest)
	err := c.BindJSON(request)
	if err != nil {
		return nil, err
	}
	cond := new(metadata.Condition)
	for key, value := range request.Filters {
		if _, ok := timeFileds[key]; ok {
			switch value.(type) {
			case string:
				cond.Filters[key], _ = time.Parse("2006-01-02 15:04:05", value.(string))
			default:
				tf := map[string]time.Time{}
				for subkey, subVal := range value.(map[string]interface{}) {
					tf[subkey], _ = time.Parse("2006-01-02 15:04:05", subVal.(string))
				}
				cond.Filters[key] = tf
			}
		} else {
			cond.Filters[key] = value
		}
	}
	cond.Limit = request.PageSize
	cond.Offset = (request.PageNo - 1) * request.PageSize
	return cond, nil
}
