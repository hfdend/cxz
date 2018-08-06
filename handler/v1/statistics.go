package v1

import "github.com/gin-gonic/gin"

type statistics int

func (statistics) Statistics(c *gin.Context) {
	var args struct {
		StartTime int64
		EndTime   int64
	}
}
