package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
	"github.com/jinzhu/gorm"
)

type OrderProduct struct {
	Model
	OrderID string `json:"order_id"`
	// 此项价格
	IPrice float64 `json:"iprice"`
	// 商品数量
	Number int `json:"number"`

	ProductID int `json:"product_id"`

	Name   string `json:"name"`
	IsPlan Sure   `json:"is_plan"`
	// 类型
	Type string `json:"type"`
	// 味道
	Taste string `json:"taste"`
	// 商品规格
	Unit string `json:"unit"`
	// 售价
	Price float64 `json:"price"`
	// 图片
	Image string `json:"image"`
	// 备注
	Mark string `json:"mark"`
	// 介绍
	Intro   string `json:"intro"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`

	ImageSrc string `json:"image_src" gorm:"-"`
}

func (OrderProduct) TableName() string {
	return "order_product"
}

var OrderProductDefault OrderProduct

func (op *OrderProduct) Insert(db *gorm.DB) error {
	op.Created = time.Now().Unix()
	op.Updated = time.Now().Unix()
	return db.Create(op).Error
}

func (OrderProduct) GetByOrderID(orderID string) (list []*OrderProduct, err error) {
	if err = cli.DB.Model(OrderProduct{}).Where("order_id = ?", orderID).Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (OrderProduct) GetByOrderIDs(orderIDs []string) (list []*OrderProduct, err error) {
	if err = cli.DB.Model(OrderProduct{}).Where("order_id in (?)", orderIDs).Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (op *OrderProduct) SetImageSrc() {
	if op.Image == "" {
		return
	}
	c := conf.Config.Aliyun.OSS
	op.ImageSrc = fmt.Sprintf("%s/%s", strings.TrimRight(c.Domain, "/"), strings.TrimLeft(op.Image, "/"))
	return
}
