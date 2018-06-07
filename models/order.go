package models

import (
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

// 订单状态
// 1: 等待付款
// 2: 付款成功
// 3: 订单发货中 (月够订单)
// 4: 发货完成
// swagger:model OrderStatus
type OrderStatus int

const (
	OrderStatusWaitting OrderStatus = iota + 1
	OrderStatusSuccess
	OrderStatusDelivering
	OrderStatusDeliveried
)

// 订单
// swagger:model Order
type Order struct {
	Model
	OrderID string `json:"order_id"`
	UserID  int    `json:"user_id"`
	// 如果是月够的则有值
	PlanID int `json:"plan_id"`
	// 金额
	Price float64 `json:"price"`
	// 运费
	Freight float64 `json:"freight"`
	// 支付金额
	PaymentPrice float64 `json:"payment_price"`
	// 支付方式 1: 微信支付
	PaymentMethod int `json:"payment_method"`
	// 买家留言
	Notice string      `json:"notice"`
	Status OrderStatus `json:"status"`
	// 创建时间
	Created int64 `json:"created"`
	// 支付截止时间
	ExpTime int64 `json:"exp_time"`
	// 支付时间
	PaymentTime int64 `json:"payment_time"`
	UpdateTime  int64 `json:"update_time"`

	OrderProducts []*OrderProduct `json:"order_products" gorm:"-"`
	OrderAddress  *OrderAddress   `json:"order_address" gorm:"-"`
}

type OrderCondition struct {
	UserID    int         `json:"user_id" form:"user_id"`
	OrderID   string      `json:"order_id" form:"order_id"`
	StartTime int64       `json:"start_time" form:"start_time"`
	EndTime   int64       `json:"end_time" form:"end_time"`
	Status    OrderStatus `json:"status" form:"status"`
}

var OrderDefault Order

func (Order) TableName() string {
	return "order"
}

func (o *Order) Insert(db *gorm.DB) error {
	o.Created = time.Now().Unix()
	o.UpdateTime = time.Now().Unix()
	return db.Create(o).Error
}

func (*Order) GetByOrderID(orderID string) (*Order, error) {
	var data Order
	if err := cli.DB.Where("order_id = ?", orderID).Find(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (Order) GetByOrderIDAndUserID(orderID string, userID int) (*Order, error) {
	var data Order
	if err := cli.DB.Where("order_id = ? and user_id = ?", orderID, userID).Find(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (Order) GetList(cond OrderCondition, pager *Pager) (list []*Order, err error) {
	db := cli.DB.Model(Order{})
	if cond.OrderID != "" {
		db = db.Where("order_id = ?", cond.OrderID)
	}
	if cond.UserID != 0 {
		db = db.Where("user_id = ?", cond.UserID)
	}
	if cond.StartTime != 0 {
		db = db.Where("created >= ?", cond.StartTime)
	}
	if cond.EndTime != 0 {
		db = db.Where("created < ?", cond.EndTime)
	}
	if cond.Status != 0 {
		db = db.Where("status = ?", cond.Status)
	}
	db = db.Where("exp_time > ? or exp_time = 0", time.Now())
	if pager != nil {
		if db, err = pager.Exec(db); err != nil {
			return
		}
		if pager.Count == 0 {
			return
		}
	}
	if err = db.Order("id desc").Find(&list).Error; err != nil {
		return
	}
	return
}
