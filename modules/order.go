package modules

import (
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

func (order) Build(userId int, info []OrderProductInfo) (o *models.Order, err error) {
	o = new(models.Order)
	o.UserId = userId
	if o.OrderId, err = models.BuildOrderID(); err != nil {
		return
	}
	o.Status = models.OrderStatusWatting
	for _, v := range info {
		if v.Number <= 0 {
			err = errors.New("数量错误")
			return
		}
		var product *models.Product
		if product, err = models.ProductDefault.GetByID(v.ProductID); err != nil {
			return
		} else if product.Price < 0 {
			err = errors.New("商品金额错误")
			return
		}
		orderProduct := new(models.OrderProduct)
		orderProduct.OrderId = o.OrderId
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
	for _, v := range o.OrderProducts {
		if err = v.Insert(db); err != nil {
			return
		}
	}
	return
}
