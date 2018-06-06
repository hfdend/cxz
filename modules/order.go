package modules

import (
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/utils"
)

type order int

var Order order

type OrderProductInfo struct {
	ProductID int `json:"product_id"`
	Number    int `json:"number"`
}

func (order) Build(userID, addressID int, info []OrderProductInfo) (o *models.Order, err error) {
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
	o.UserId = userID
	if o.OrderID, err = models.BuildOrderID(); err != nil {
		return
	}
	o.Status = models.OrderStatusWatting
	o.ExpTime = time.Now().Add(20 * time.Minute).Unix()
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
	o.Fee = 0
	o.PaymentPrice = o.Price + o.Fee
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
