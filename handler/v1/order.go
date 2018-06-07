package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
)

type order int

var Order order

// swagger:parameters Order_GetFreight
type OrderGetFreightArgs struct {
	// in: body
	Body struct {
		ProductInfo []modules.OrderProductInfo `json:"product_info"`
		WeekNumber  int                        `json:"week_number"`
	}
}

// 运费
// swagger:response OrderGetFreightResp
type OrderGetFreightResp struct {
	// in: body
	Body struct {
		// 运费
		Freight float64 `json:"freight"`
	}
}

// swagger:route POST /order/products/freight Order_GetFreight
// 获取商品运费
// responses:
//    200: OrderGetFreightResp
func (order) GetFreight(c *gin.Context) {
	var args OrderGetFreightArgs
	var resp OrderGetFreightResp
	if c.Bind(args.Body) != nil {
		return
	}
	resp.Body.Freight = 0
	JSON(c, resp.Body)
}

// swagger:parameters Order_Build
type OrderBuildArgs struct {
	// in: body
	Body struct {
		AddressID   int                        `json:"address_id"`
		ProductInfo []modules.OrderProductInfo `json:"product_info"`
		// 买家留言
		Notice     string `json:"notice"`
		WeekNumber int    `json:"week_number"`
	}
}

// 订单详情
// swagger:response OrderBuildResp
type OrderBuildResp struct {
	// in: body
	Body *models.Order
}

// swagger:route POST /order/build 订单 Order_Build
// 下单
// responses:
//     200: OrderBuildResp
func (order) Build(c *gin.Context) {
	var args OrderBuildArgs
	if c.Bind(&args.Body) != nil {
		return
	}
	user := GetUser(c)
	order, err := modules.Order.Build(user.ID, args.Body.AddressID, args.Body.ProductInfo, args.Body.Notice, args.Body.WeekNumber)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, order)
	}
}

// swagger:parameters Order_GetByOrderID
type OrderGetByOrderIDArgs struct {
	OrderID string `json:"order_id" form:"order_id"`
}

// 订单详情
// swagger:response OrderGetByOrderIDResp
type OrderGetByOrderIDResp struct {
	// in: body
	Body *models.Order
}

// swagger:route GET /order/detail 订单 Order_GetByOrderID
// 订单详情
// responses:
//     200: OrderGetByOrderIDResp
func (order) GetByOrderID(c *gin.Context) {
	var args OrderGetByOrderIDArgs
	if c.Bind(&args) != nil {
		return
	}
	user := GetUser(c)
	order, err := modules.Order.GetByID(args.OrderID, user.ID)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, order)
	}
}

// swagger:parameters Order_GetList
type OrderGetListArgs struct {
	Page int `json:"page" form:"page"`
}

// 订单列表
// swagger:response OrderGetListResp
type OrderGetListResp struct {
	// in: body
	Body struct {
		List  []*models.Order `json:"list"`
		Pager *models.Pager   `json:"pager"`
	}
}

// swagger:route GET /order/list 订单 Order_GetList
// 订单列表
// responses:
//     200: OrderGetListResp
func (order) GetList(c *gin.Context) {
	var args OrderGetListArgs
	var resp OrderGetListResp
	var err error
	if c.Bind(&args) != nil {
		return
	}
	user := GetUser(c)
	cond := models.OrderCondition{}
	cond.UserID = user.ID
	resp.Body.Pager = models.NewPager(args.Page, 20)
	if resp.Body.List, err = modules.Order.GetListDetail(cond, resp.Body.Pager); err != nil {
		JSON(c, err)
	} else {
		JSON(c, resp.Body)
	}
}
