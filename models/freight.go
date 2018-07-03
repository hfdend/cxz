package models

import (
	"fmt"

	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

type Freight struct {
	Model
	// 地区code
	Code string `json:"code"`
	Name string `json:"name"`
	// 单期运费-自选
	Amount float64 `json:"amount"`
	// 订单免多少免运费(-1表示不设置此条件)-自选
	OrderFree float64 `json:"order_free"`
	// 单期免多少免运费(-1表示不设置此条件)-自选
	PhaseFree float64 `json:"phase_free"`
	// 单期运费-月够
	PlanAmount float64 `json:"plan_amount"`
	// 订单免多少免运费(-1表示不设置此条件)-月够
	PlanOrderFree float64 `json:"plan_order_free"`
	// 单期免多少免运费(-1表示不设置此条件)-月够
	PlanPhaseFree float64 `json:"plan_phase_free"`
	Created       int64   `json:"created"`
}

var FreightDefault Freight

func (Freight) TableName() string {
	return "freight"
}

func (Freight) GetAll() (list []*Freight, err error) {
	err = cli.DB.Find(&list).Error
	return
}

func (Freight) GetByCode(code string) (*Freight, error) {
	var data Freight
	if err := cli.DB.Where("code = ?", code).Find(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (f *Freight) Save() error {
	f.Created = time.Now().Unix()
	return cli.DB.Save(f).Error
}

func (Freight) Truncate() error {
	return cli.DB.Exec(fmt.Sprintf("truncate table %s", Freight{}.TableName())).Error
}
