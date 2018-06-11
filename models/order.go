package models

import (
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

// 订单状态
// 1: 等待付款
// 2: 付款成功
// swagger:model OrderStatus
type OrderStatus int

const (
	OrderStatusWaiting OrderStatus = iota + 1
	OrderStatusSuccess
	OrderStatusDelivering
	OrderStatusDeliveried
)

// 发货状态
// swagger:model DeliveryStatus
// 1: 等待发货
// 2: 部分周期发货
// 2: 发货完成
type DeliveryStatus int

const (
	DeliveryStatusWaiting DeliveryStatus = iota + 1
	DeliveryStatusIng
	DeliveryStatusOver
)

// 订单
// swagger:model Order
type Order struct {
	Model
	OrderID       string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	UserID        int    `json:"user_id"`
	Body          string `json:"body"`
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
	// 购买了几周
	WeekNumber int `json:"week_number"`
	// 发货了几周
	WeekDelivered int `json:"week_delivered"`
	// 买家留言
	Notice string `json:"notice"`
	// 订单支付状态
	Status OrderStatus `json:"status"`
	// 发货状态
	DeliveryStatus DeliveryStatus `json:"delivery_status"`
	// 创建时间
	Created int64 `json:"created"`
	// 支付截止时间
	ExpTime int64 `json:"exp_time"`
	// 支付时间
	PaymentTime int64 `json:"payment_time"`
	UpdateTime  int64 `json:"update_time"`

	// 订单包含的商品
	OrderProducts []*OrderProduct `json:"order_products" gorm:"-"`
	// 订单地址
	OrderAddress *OrderAddress `json:"order_address" gorm:"-"`
	// 订单发货计划
	OrderPlans []*OrderPlan `json:"order_plans" gorm:"-"`
}

type OrderCondition struct {
	UserID         int            `json:"user_id" form:"user_id"`
	OrderID        string         `json:"order_id" form:"order_id"`
	StartTime      int64          `json:"start_time" form:"start_time"`
	EndTime        int64          `json:"end_time" form:"end_time"`
	Status         OrderStatus    `json:"status" form:"status"`
	DeliveryStatus DeliveryStatus `json:"delivery_status" form:"delivery_status"`
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

func (*Order) GetByOrderIDForUpdate(db *gorm.DB, orderID string) (*Order, error) {
	var data Order
	if err := db.Set("gorm:query_option", "FOR UPDATE").Where("order_id = ?", orderID).Find(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (Order) GetByOrderIDAndUserID(orderID string, userID int) (*Order, error) {
	var data Order
	db := cli.DB.Where("order_id = ?", orderID)
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	if err := db.Find(&data).Error; err != nil {
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
	if cond.DeliveryStatus != 0 {
		db = db.Where("delivery_status = ?", cond.DeliveryStatus)
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

func (o *Order) ToSuccess(db *gorm.DB, transactionID string) error {
	data := map[string]interface{}{
		"status":         OrderStatusSuccess,
		"exp_time":       0,
		"payment_time":   time.Now().Unix(),
		"update_time":    time.Now().Unix(),
		"transaction_id": transactionID,
		"payment_method": 1,
	}
	return db.Model(o).Update(data).Error
}

func (o *Order) UpdateDeliveryStatus(db *gorm.DB, deliveryStatus DeliveryStatus, item int) error {
	data := map[string]interface{}{
		"update_time":     time.Now().Unix(),
		"delivery_status": deliveryStatus,
		"week_delivered":  item,
	}
	return db.Model(o).Update(data).Error
}
