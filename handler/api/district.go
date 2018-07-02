package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
)

type district int

var District district

func (district) GetGradation(c *gin.Context) {
	list, err := models.DistrictDefault.GetGradation()
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}
