package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
)

type freight int

var Freight freight

func (freight) SaveAll(c *gin.Context) {
	var args []*models.Freight
	if c.Bind(&args) != nil {
		return
	}
	if err := modules.Freight.SaveAll(args); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}

func (freight) GetList(c *gin.Context) {
	if list, err := models.FreightDefault.GetAll(); err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}
