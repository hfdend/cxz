package modules

import (
	"fmt"
	"time"

	"strings"

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

func (order) Build(userID, addressID int, info []OrderProductInfo, notice string, weekNumber int) (o *models.Order, err error) {
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
	o = new(models.Order)
	o.UserID = userID
	if o.OrderID, err = models.BuildOrderID(); err != nil {
		return
	}
	o.Status = models.OrderStatusWaitting
	o.ExpTime = time.Now().Add(20 * time.Minute).Unix()
	var body string
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
		body += fmt.Sprintf(",%s", product.Name)
		orderProduct := new(models.OrderProduct)
		orderProduct.OrderID = o.OrderID
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
		o.OrderProducts = append(o.OrderProducts, orderProduct)
		o.Price += orderProduct.IPrice
	}
	body = strings.TrimLeft(body, ",")
	bodyRune := []rune(body)
	if len(bodyRune) > 32 {
		o.Body = string(bodyRune[0:32])
	} else {
		o.Body = body
	}
	o.Notice = notice
	o.Freight = 0
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
		addressesMap = map[string]*models.OrderAddress{}
		productsMap  = map[string][]*models.OrderProduct{}
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
	for _, v := range addresses {
		addressesMap[v.OrderID] = v
	}
	for _, v := range products {
		productsMap[v.OrderID] = append(productsMap[v.OrderID], v)
	}
	for _, v := range list {
		v.OrderAddress, _ = addressesMap[v.OrderID]
		v.OrderProducts, _ = productsMap[v.OrderID]
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
	//order, err := models.OrderDefault.GetByOrderID(orderID)
	return nil
}
