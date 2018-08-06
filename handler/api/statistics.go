package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/modules"
)

type statistics int

var Statistics statistics

func (statistics) Statistics(c *gin.Context) {
	var args struct {
		StartTime int64 `json:"start_time" form:"start_time"`
		EndTime   int64 `json:"end_time" form:"end_time"`
	}
	if c.Bind(&args) != nil {
		return
	}
	data, err := modules.Statistics.Statistics(args.StartTime, args.EndTime)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, data)
	}
}
