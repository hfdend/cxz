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
	if list, err := modules.Order.GetListDetail(args.Condition, pager); err != nil {
		JSON(c, err)
	} else {
		JSON(c, map[string]interface{}{"list": list, "pager": pager})
	}
}

func (order) GetByOrderID(c *gin.Context) {
	var args struct {
		OrderID string `json:"order_id" form:"order_id"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if order, err := modules.Order.GetByID(args.OrderID, 0); err != nil {
		JSON(c, err)
	} else {
		JSON(c, order)
	}
}

func (order) Delivery(c *gin.Context) {
	var args struct {
		OrderID       string `json:"order_id"`
		Item          int    `json:"item"`
		Express       string `json:"express"`
		WaybillNumber string `json:"waybill_number"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if err := modules.Order.Delivery(args.OrderID, args.Item, args.Express, args.WaybillNumber); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}

func (order) GetNeedSendList(c *gin.Context) {
	var args struct {
		Page int `json:"page" form:"page"`
	}
	if c.Bind(&args) != nil {
		return
	}
	pager := models.NewPager(args.Page, 20)
	list, err := models.OrderPlanDefault.GetNeedSendList(pager)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}

func (order) CancelOrder(c *gin.Context) {
	var args struct {
		OrderID      string  `json:"order_id"`
		Items        []int   `json:"items"`
		RefundAmount float64 `json:"refund_amount"`
	}
	if c.Bind(&args) != nil {
		return
	}
	user := getUser(c)
	if err := modules.Order.CancelOrder(user.ID, args.OrderID, args.Items, args.RefundAmount); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}
