package modules

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/payment/wxpay"
	"github.com/hfdend/cxz/utils"
)

type order int

var Order order

type OrderProductInfo struct {
	ProductID int `json:"product_id"`
	Number    int `json:"number"`
}

func (order) GetOrderProducts(orderID string, info []OrderProductInfo, weekNumber int) (price, freight float64, body string, products []*models.OrderProduct, err error) {
	var maxFreight float64
	for _, v := range info {
		if v.Number <= 0 {
			err = errors.New("数量错误")
			return
		}
		var product *models.Product
		if product, err = models.ProductDefault.GetByID(v.ProductID); err != nil {
			return
		} else if product == nil {
			err = errors.New("未找到商品")
			return
		} else if product.Price < 0 {
			err = errors.New("商品金额错误")
			return
		}
		if product.Freight > maxFreight {
			maxFreight = product.Freight
		}
		body += fmt.Sprintf(",%s", product.Name)
		orderProduct := new(models.OrderProduct)
		orderProduct.OrderID = orderID
		orderProduct.ProductID = product.ID
		orderProduct.Number = v.Number
		orderProduct.IPrice = utils.Round(product.Price*float64(v.Number), 2)

		orderProduct.Name = product.Name
		orderProduct.Type = product.Type
		orderProduct.Taste = product.Taste
		orderProduct.Unit = product.Unit
		orderProduct.Price = product.Price
		orderProduct.Image = product.Image
		orderProduct.Mark = product.Mark
		orderProduct.Intro = product.Intro
		products = append(products, orderProduct)
		price += orderProduct.IPrice
	}
	body = strings.TrimLeft(body, ",")
	bodyRune := []rune(body)
	if len(bodyRune) > 32 {
		body = string(bodyRune[0:32])
	}
	freight = utils.Round(maxFreight*float64(weekNumber), 2)
	price = utils.Round(price*float64(weekNumber), 2)
	return
}

func (order) Build(userID, addressID int, info []OrderProductInfo, notice string, weekNumber int) (o *models.Order, err error) {
	if len(info) > 20 {
		err = errors.New("一个订单最多支持20个商品，请分开结算")
		return
	}
	var address *models.Address
	if address, err = models.AddressDefault.GetByID(addressID); err != nil {
		return
	} else if address == nil {
		err = errors.New("请选择收货地址")
		return
	} else if address.UserID != userID {
		err = errors.New("收货地址错误")
		return
	}
	if weekNumber <= 0 {
		weekNumber = 1
	}
	o = new(models.Order)
	o.UserID = userID
	if o.OrderID, err = models.BuildOrderID(); err != nil {
		return
	}
	o.WeekNumber = weekNumber
	o.Status = models.OrderStatusWaiting
	o.DeliveryStatus = models.DeliveryStatusWaiting
	o.ExpTime = time.Now().Add(20 * time.Minute).Unix()
	var (
		body           string
		price, freight float64
		products       []*models.OrderProduct
	)

	if price, freight, body, products, err = Order.GetOrderProducts(o.OrderID, info, weekNumber); err != nil {
		return
	}
	o.Body = body
	o.Freight = freight
	// 金额等于 商品金额 * 期数 + 运费
	o.Price = price
	o.Notice = notice
	o.PaymentPrice = o.Price + o.Freight
	o.OrderAddress = new(models.OrderAddress)
	o.OrderAddress.OrderID = o.OrderID
	o.OrderAddress.AddressID = address.ID
	o.OrderAddress.UserID = address.UserID
	o.OrderAddress.Name = address.Name
	o.OrderAddress.Phone = address.Phone
	o.OrderAddress.DistrictCode = address.DistrictCode
	o.OrderAddress.DistrictName = address.DistrictName
	o.OrderAddress.DetailAddress = address.DetailAddress
	o.OrderProducts = products
	o.ApplyStatus = models.ApplyStatusNil
	db := cli.DB.Begin()
	defer func() {
		if err == nil {
			db.Commit()
		} else {
			db.Rollback()
		}
	}()
	if err = o.Insert(db); err != nil {
		return
	}
	if err = o.OrderAddress.Insert(db); err != nil {
		return
	}
	for _, v := range o.OrderProducts {
		if err = v.Insert(db); err != nil {
			return
		}
	}
	return
}

func (order) GetByID(orderID string, userID int) (o *models.Order, err error) {
	if o, err = models.OrderDefault.GetByOrderIDAndUserID(orderID, userID); err != nil {
		return
	} else if o == nil {
		return
	}
	if o.OrderAddress, err = models.OrderAddressDefault.GetByOrderID(o.OrderID); err != nil {
		return
	}
	if o.OrderProducts, err = models.OrderProductDefault.GetByOrderID(o.OrderID); err != nil {
		return
	}
	if o.OrderPlans, err = models.OrderPlanDefault.GetByOrderID(o.OrderID); err != nil {
		return
	}
	return
}

func (order) GetList(cond models.OrderCondition, pager *models.Pager) (list []*models.Order, err error) {
	if list, err = models.OrderDefault.GetList(cond, pager); err != nil {
		return
	}
	return
}

func (order) GetListDetail(cond models.OrderCondition, pager *models.Pager) (list []*models.Order, err error) {
	if list, err = models.OrderDefault.GetList(cond, pager); err != nil {
		return
	}
	if len(list) == 0 {
		return
	}
	var (
		orderIDs     []string
		addresses    []*models.OrderAddress
		products     []*models.OrderProduct
		plans        []*models.OrderPlan
		addressesMap = map[string]*models.OrderAddress{}
		productsMap  = map[string][]*models.OrderProduct{}
		planMap      = map[string][]*models.OrderPlan{}
	)
	for _, v := range list {
		orderIDs = append(orderIDs, v.OrderID)
	}
	if addresses, err = models.OrderAddressDefault.GetByOrderIDs(orderIDs); err != nil {
		return
	}
	if products, err = models.OrderProductDefault.GetByOrderIDs(orderIDs); err != nil {
		return
	}
	if plans, err = models.OrderPlanDefault.GetByOrderIDs(orderIDs); err != nil {
		return
	}
	for _, v := range addresses {
		addressesMap[v.OrderID] = v
	}
	for _, v := range products {
		productsMap[v.OrderID] = append(productsMap[v.OrderID], v)
	}
	for _, v := range plans {
		planMap[v.OrderID] = append(planMap[v.OrderID], v)
	}
	for _, v := range list {
		v.OrderAddress, _ = addressesMap[v.OrderID]
		v.OrderProducts, _ = productsMap[v.OrderID]
		v.OrderPlans, _ = planMap[v.OrderID]
	}
	return
}

// 小程序支付
func (order) WXAPay(orderID string, user *models.User, ip string) (wxaRequest *wxpay.PayWxaRequest, err error) {
	var order *models.Order
	if order, err = models.OrderDefault.GetByOrderIDAndUserID(orderID, user.ID); err != nil {
		return
	} else if order == nil {
		err = errors.New("订单不存在")
		return
	} else if order.Status != models.OrderStatusWaiting {
		err = errors.New("请勿重复支付")
		return
	} else if order.ApplyStatus != models.ApplyStatusNil {
		err = errors.New("订单正在取消中")
		return
	}
	var c wxpay.PayConfig
	c.AppId = conf.Config.WXPay.AppId
	c.MchId = conf.Config.WXPay.MchId
	c.Key = conf.Config.WXPay.Key
	c.NotifyUrl = conf.Config.WXPay.NotifyUrl
	api := wxpay.NewApi(c)
	api.Logger = logrus.New()
	query := wxpay.NewPayUnifiedOrder()
	query.SetBody(order.Body)
	query.SetOutTradeNo(order.OrderID)
	if conf.Config.WXPay.TestAmount > 0 {
		query.SetTotalFee(int(utils.Round(conf.Config.WXPay.TestAmount*100, 0)))
	} else {
		query.SetTotalFee(int(utils.Round(order.PaymentPrice*100, 0)))
	}
	query.SetOpenId(user.OpenID)
	query.SetTradeType("JSAPI")
	query.SetSpbillCreateIp(ip)
	var result *wxpay.PayResults
	if result, err = api.UnifiedOrder(query, 5*time.Second); err != nil {
		return
	}
	if result.ResultCode != "SUCCESS" {
		err = errors.New(result.ReturnMsg)
		return
	}
	if result.ResultCode != "SUCCESS" {
		err = errors.New(result.ErrCodeDes)
		return
	}
	wxaRequest = new(wxpay.PayWxaRequest)
	wxaRequest.SignType = "MD5"
	wxaRequest.NonceStr = wxpay.GetNonceStr()
	wxaRequest.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	wxaRequest.Package = fmt.Sprintf("prepay_id=%s", result.PrepayId)
	wxaRequest.SetSign(api.Config.AppId, api.Config.Key)
	return
}

func (order) PaymentSuccess(orderID, transactionID string) error {
	db := cli.DB.Begin()
	order, err := models.OrderDefault.GetByOrderIDForUpdate(db, orderID)
	if err != nil {
		db.Rollback()
		return err
	} else if order == nil || order.Status != models.OrderStatusWaiting {
		db.Rollback()
		return nil
	}
	if err := order.ToSuccess(db, transactionID); err != nil {
		db.Rollback()
		return err
	}
	products, err := models.OrderProductDefault.GetByOrderID(orderID)
	if err != nil {
		db.Rollback()
		return err
	}
	y, m, d := time.Now().Date()
	t := time.Date(y, m, d, 0, 0, 0, 0, time.Local)

	// 添加发货计划
	for i := 0; i < order.WeekNumber; i++ {
		op := new(models.OrderPlan)
		op.OrderID = order.OrderID
		op.UserID = order.UserID
		op.Title = order.Body
		if len(products) > 0 {
			op.Image = products[0].Image
			op.SetImageSrc()
		}
		if order.Freight > 0 {
			op.Freight = utils.Round(order.Freight/float64(order.WeekNumber), 2)
		}
		if order.Price > 0 {
			op.Price = utils.Round(order.Price/float64(order.WeekNumber), 2)
		}
		op.Item = i + 1
		op.TotalItem = order.WeekNumber
		if i == 0 {
			op.PlanTime = t.Add(24 * 7 * time.Hour).Unix()
		} else {
			op.PlanTime = t.Add(time.Duration(i+1) * 24 * 7 * time.Hour).Unix()
		}
		op.Status = models.PlanStatusWaiting
		if err = op.Insert(db); err != nil {
			db.Rollback()
			return err
		}
	}
	db.Commit()
	return nil
}

// 订单发货
func (order) Delivery(orderID string, item int, express, waybillNumber string) error {
	db := cli.DB.Begin()
	order, err := models.OrderDefault.GetByOrderIDForUpdate(db, orderID)
	if err != nil {
		db.Rollback()
		return err
	} else if order == nil {
		db.Rollback()
		return nil
	} else if order.DeliveryStatus == models.DeliveryStatusOver {
		db.Rollback()
		return errors.New("订单已发货完成")
	}
	planList, err := models.OrderPlanDefault.GetByOrderIDForUpdate(db, orderID)
	var plan *models.OrderPlan
	if err != nil {
		db.Rollback()
		return err
	}
	for _, v := range planList {
		if v.Item == item {
			plan = v
		} else if v.Item < item && v.Status == models.PlanStatusWaiting {
			db.Rollback()
			return errors.New("请按顺序发货")
		}
	}
	if plan == nil {
		db.Rollback()
		return errors.New("没有找到发货计划")
	} else if plan.Status != models.PlanStatusWaiting {
		db.Rollback()
		return errors.New("计划状态不正确")
	} else if plan.ApplyStatus == models.ApplyStatusWaiting {
		db.Rollback()
		return errors.New("计划正在申请取消不能发货")
	} else if plan.ApplyStatus == models.ApplyStatusSuccess {
		db.Rollback()
		return errors.New("计划已经申请取消不能发货")
	}
	if err := plan.Delivery(db, express, waybillNumber); err != nil {
		db.Rollback()
		return err
	}
	var deliveryStatus models.DeliveryStatus
	if plan.Item != plan.TotalItem {
		deliveryStatus = models.DeliveryStatusIng
	} else {
		deliveryStatus = models.DeliveryStatusOver
	}
	if err := order.UpdateDeliveryStatus(db, deliveryStatus, plan.Item); err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

func (order) PlanDelay(userID int, orderID string, item int, day string) error {
	t, err := time.ParseInLocation("20060102", day, time.Local)
	if err != nil {
		return errors.New("时间格式错误")
	}
	db := cli.DB.Begin()
	plans, err := models.OrderPlanDefault.GetByOrderIDAndUserIDForUpdate(db, orderID, userID)
	if err != nil {
		db.Rollback()
		return err
	}
	if len(plans) == 0 {
		db.Rollback()
		return nil
	}
	var delayIDs []int
	var diff int64
	for _, v := range plans {
		if v.Item == item {
			if v.Status != models.PlanStatusWaiting {
				db.Rollback()
				return errors.New("已发货不能延期")
			}
			if v.ApplyStatus != models.ApplyStatusNil {
				db.Rollback()
				return errors.New("不能延期")
			}
			if v.PlanTime == t.Unix() {
				db.Rollback()
				return nil
			}
			diff = t.Unix() - v.PlanTime
			delayIDs = append(delayIDs, v.ID)
		} else if v.Item > item && v.Status == models.PlanStatusWaiting && v.ApplyStatus == models.ApplyStatusNil {
			delayIDs = append(delayIDs, v.ID)
		}
	}
	// 数据更改
	if err := models.OrderPlanDefault.Delay(db, delayIDs, diff); err != nil {
		db.Rollback()
		return err
	}
	return nil
}
