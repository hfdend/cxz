package models

import (
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

type Product struct {
	Model
	Name string `json:"name"`
	// 类型
	Type string `json:"type"`
	// 味道
	Taste string `json:"taste"`
	// 商品规格
	Unit string `json:"unit"`
	// 售价
	Price int64 `json:"price"`
	// 图片
	Image string `json:"image"`
	// 备注
	Mark string `json:"mark"`
	// 介绍
	Intro   string `json:"intro"`
	IsDel   Sure   `json:"is_del"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}

var ProductDefault Product

func (Product) TableName() string {
	return "product"
}

func (p *Product) Save() error {
	if p.Created == 0 {
		p.Created = time.Now().Unix()
	}
	p.Updated = time.Now().Unix()
	return cli.DB.Save(p).Error
}

func (Product) GetByID(id int) (*Product, error) {
	var data Product
	if err := cli.DB.Where("id = ? and is_del = ?", id, SureNo).Find(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (Product) DelByID(id int) error {
	data := map[string]interface{}{
		"is_del":  SureYes,
		"updated": time.Now().Unix(),
	}
	return cli.DB.Where("id = ?", id).Update(data).Error
}

func (Product) GetList(pager *Pager) (list []*Product, err error) {
	err = cli.DB.Where("is_del = ?", SureNo).Find(&list).Error
	return
}
