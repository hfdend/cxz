package models

import (
	"strings"
	"time"

	"fmt"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
	"github.com/jinzhu/gorm"
)

type Product struct {
	Model
	Name string `json:"name"`
	// 是否是月够商品
	IsPlan Sure `json:"is_plan"`
	// 类型
	Type string `json:"type"`
	// 味道
	Taste string `json:"taste"`
	// 最小体重
	MinWeight int `json:"min_weight"`
	// 最大体重
	MaxWeight int `json:"max_weight"`
	// 最小年龄
	MinAge int `json:"min_age"`
	// 最大年龄
	MaxAge int `json:"max_age"`
	// 商品规格
	Unit string `json:"unit"`
	// 售价
	Price float64 `json:"price"`
	// 运费
	Freight float64 `json:"freight"`
	// 图片
	Image string `json:"image"`
	// 备注
	Mark string `json:"mark"`
	// 介绍
	Intro   string `json:"intro"`
	IsDel   Sure   `json:"is_del"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`

	ImageSrc string `json:"image_src" gorm:"-"`
}

// 商品搜索条件
// swagger:model ProductCondition
type ProductCondition struct {
	// 种类
	Type string `json:"type" form:"type"`
	// 口味
	Taste string `json:"taste" form:"taste"`
	// 是否是月够商品
	IsPlan Sure `json:"is_plan" form:"is_plan"`
	// 最小体重
	MinWeight int `json:"min_weight" form:"min_weight"`
	// 最大体重
	MaxWeight int `json:"max_weight" form:"max_weight"`
	// 最小年龄
	MinAge int `json:"min_age" form:"min_age"`
	// 最大年龄
	MaxAge int `json:"max_age" form:"max_age"`
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
	data.SetImageSrc()
	return &data, nil
}

func (Product) DelByID(id int) error {
	data := map[string]interface{}{
		"is_del":  SureYes,
		"updated": time.Now().Unix(),
	}
	return cli.DB.Model(Product{}).Where("id = ?", id).Update(data).Error
}

func (Product) GetList(cond ProductCondition, pager *Pager) (list []*Product, err error) {
	db := cli.DB.Model(Product{}).Where("is_del = ?", SureNo)
	if cond.Type != "" {
		db = db.Where("type = ?", cond.Type)
	}
	if cond.Taste != "" {
		db = db.Where("taste = ?", cond.Taste)
	}
	if cond.IsPlan != SureNil {
		db = db.Where("is_plan = ?", cond.IsPlan)
	}
	if cond.MinWeight != 0 {
		db = db.Where("max_weight > ?", cond.MinWeight)
	}
	if cond.MaxWeight != 0 {
		db = db.Where("min_weight <= ?", cond.MaxWeight)
	}
	if cond.MinAge != 0 {
		db = db.Where("max_age > ?", cond.MinAge)
	}
	if cond.MaxAge != 0 {
		db = db.Where("min_age <= ?", cond.MaxAge)
	}
	if pager != nil {
		if db, err = pager.Exec(db); err != nil {
			return
		} else if pager.Count == 0 {
			return
		}
	}
	err = db.Find(&list).Error
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (p *Product) SetImageSrc() {
	if p.Image == "" {
		return
	}
	c := conf.Config.Aliyun.OSS
	p.ImageSrc = fmt.Sprintf("%s/%s", strings.TrimRight(c.Domain, "/"), strings.TrimLeft(p.Image, "/"))
	return
}
