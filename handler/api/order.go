package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
)

type order int

var Order order

func (order) GetList(c *gin.Context) {
	var args struct {
		Condition models.OrderCondition
		Page      int `json:"page" form:"page"`
	}
	if c.Bind(&args) != nil {
		return
	}
	pager := models.NewPager(args.Page, 20)
	if list, err := modules.Order.GetList(args.Condition, pager); err != nil {
		JSON(c, err)
	} else {
		JSON(c, map[string]interface{}{"list": list, "pager": pager})
	}
}
