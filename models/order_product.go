package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OrderProduct struct {
	Model
	OrderId string `json:"order_id"`
	// 此项价格
	IPrice float64 `json:"iprice"`
	// 商品数量
	Number int `json:"number"`

	ProductID int `json:"product_id"`

	Name string `json:"name"`
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

func (op *OrderProduct) Insert(db *gorm.DB) error {
	op.Created = time.Now().Unix()
	op.Updated = time.Now().Unix()
	return db.Create(op).Error
}
