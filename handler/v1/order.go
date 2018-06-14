package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
	"github.com/hfdend/cxz/payment/wxpay"
)

type order int

var Order order

// 订阅列表
// swagger:response GetOrderPlanListResp
type GetOrderPlanListResp struct {
	// in: body
	Body []*models.OrderPlan
}

// swagger:route GET /order/plans 订单 Order_GetOrderPlanList
// 获取订单订阅
// responses:
//     200: GetOrderPlanListResp
func (order) GetOrderPlanList(c *gin.Context) {
	user := GetUser(c)
	if list, err := models.OrderPlanDefault.GetByUserID(user.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}

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
		// 商品价格
		Price float64 `json:"price"`
	}
}

// swagger:route POST /order/products/freight Order_GetFreight
// 获取商品运费
// responses:
//    200: OrderGetFreightResp
func (order) GetFreight(c *gin.Context) {
	var args OrderGetFreightArgs
	var resp OrderGetFreightResp
	if c.Bind(&args.Body) != nil {
		return
	}
	if price, freight, _, _, err := modules.Order.GetOrderProducts("", args.Body.ProductInfo, args.Body.WeekNumber); err != nil {
		JSON(c, err)
	} else {
		resp.Body.Freight = freight
		resp.Body.Price = price
		JSON(c, resp.Body)
	}
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

// swagger:parameters Order_WXAPayment
type OrderWXAPaymentArgs struct {
	// in: body
	Body struct {
		OrderID string `json:"order_id"`
	}
}

// 小程序支付
// swagger:response OrderWXAPaymentResp
type OrderWXAPaymentResp struct {
	// in: body
	Body *wxpay.PayWxaRequest
}

// swagger:route POST /order/wxapayment 订单 Order_WXAPayment
// 小程序支付
// responses:
//     200: OrderWXAPaymentResp
func (order) WXAPayment(c *gin.Context) {
	var args OrderWXAPaymentArgs
	var resp OrderWXAPaymentResp
	var err error
	if c.Bind(&args.Body) != nil {
		return
	}
	body := args.Body
	user := GetUser(c)
	if resp.Body, err = modules.Order.WXAPay(body.OrderID, user, c.ClientIP()); err != nil {
		JSON(c, err)
	} else {
		JSON(c, resp.Body)
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

// 物流信息
// swagger:response OrderQueryExpressResp
type OrderQueryExpressResp struct {
	// in: body
	Body struct {
		ExpressData models.ExpressData   `json:"express_data"`
		Address     *models.OrderAddress `json:"address"`
	}
}

// swagger:route /order/query/express 订单 Order_QueryExpress
// 查询物流单号
// responses:
//     200: OrderQueryExpressResp
func (order) QueryExpress(c *gin.Context) {
	var args OrderQueryExpressArgs
	var resp OrderQueryExpressResp
	var err error
	if c.Bind(&args) != nil {
		return
	}
	if resp.Body.ExpressData, err = modules.Express.Query(args.Number, args.Company); err != nil {
		JSON(c, err)
		return
	}
	if resp.Body.Address, err = models.OrderAddressDefault.GetByOrderID(args.OrderID); err != nil {
		JSON(c, err)
		return
	}
	JSON(c, resp.Body)
}

// swagger:parameters Order_PlanDelay
type OrderPlanDelayArgs struct {
	// in: body
	Body struct {
		OrderID string `json:"order_id"`
		Item    int    `json:"item"`
		// 日期，格式20180301
		Day string `json:"day"`
	}
}

// swagger:route POST /order/plan/delay 订单 Order_PlanDelay
// 发货计划推迟
// responses:
//     200: SUCCESS
func (order) PlanDelay(c *gin.Context) {
	var args OrderPlanDelayArgs
	if c.Bind(&args.Body) != nil {
		return
	}
	user := GetUser(c)
	if err := modules.Order.PlanDelay(user.ID, args.Body.OrderID, args.Body.Item, args.Body.Day); err != nil {
		JSON(c, err)
	} else {
		JSON(c, SUCCESS)
	}
}
