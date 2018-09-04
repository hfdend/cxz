package api

import (
	"fmt"
	"time"

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
		AdminUserID   int    `json:"admin_user_id"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if err := modules.Order.Delivery(args.OrderID, args.Item, args.Express, args.WaybillNumber, args.AdminUserID); err != nil {
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
		JSON(c, map[string]interface{}{"list": list, "pager": pager})
	}
}

func (order) GetNeedSendListExport(c *gin.Context) {
	data, err := modules.Order.GetNeedSendListExport()
	if err != nil {
		JSON(c, err)
		return
	}
	filename := fmt.Sprintf("发货导出-%s.xlsx", time.Now().Format("20060102150405"))
	h := c.Writer.Header()
	h.Set("Content-Disposition", "attachment; charset=GBK; filename="+filename)
	h.Set("Content-Type", "application/vnd.ms-excel")
	c.Writer.Write(data)
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

func (order) UpdateAddress(c *gin.Context) {
	var args struct {
		ID            int    `json:"id"`
		DetailAddress string `json:"detail_address"`
		DistrictName  string `json:"district_name"`
		Name          string `json:"name"`
		OrderId       string `json:"order_id"`
		Phone         string `json:"phone"`
	}
	if c.Bind(&args) != nil {
		return
	}
	user := getUser(c)
	if err := modules.Order.UpdateAddress(user.ID, args.ID, args.DetailAddress, args.Name, args.Phone, args.DistrictName, args.OrderId); err != nil {
		JSON(c, err)
	} else {
		JSON(c, nil)
	}
}

// swagger:parameters Order_QueryExpress
type OrderQueryExpressArgs struct {
	// 订单ID
	OrderID string `json:"order_id" form:"order_id"`
	// 物流单号
	Number string `json:"number" form:"number"`
	// 快递公司
	Company string `json:"company" form:"company"`
}

type OrderQueryExpressResp struct {
	ExpressData models.ExpressData   `json:"express_data"`
	Address     *models.OrderAddress `json:"address"`
}

func (order) QueryExpress(c *gin.Context) {
	var args struct {
		// 订单ID
		OrderID string `json:"order_id" form:"order_id"`
		// 物流单号
		Number string `json:"number" form:"number"`
		// 快递公司
		Company string `json:"company" form:"company"`
	}
	var resp OrderQueryExpressResp
	var err error
	if c.Bind(&args) != nil {
		return
	}
	if resp.ExpressData, err = modules.Express.Query(args.Number, args.Company); err != nil {
		JSON(c, err)
		return
	}
	if resp.Address, err = models.OrderAddressDefault.GetByOrderID(args.OrderID); err != nil {
		JSON(c, err)
		return
	}
	JSON(c, resp)
}
