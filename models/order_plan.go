package models

import (
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

// 计划状态
// 1: 等待发货
// 2: 已发货
// swagger:model PlanStatus
type PlanStatus int

const (
	PlanStatusWaiting PlanStatus = iota + 1
	PlanStatusDeliveried
)

// 申请取消状态
// swagger:model ApplyStatus
// 1: 未申请
// 2: 申请取消
// 3: 取消成功
// 4: 取消拒绝
type ApplyStatus int

const (
	ApplyStatusNil ApplyStatus = iota + 1
	ApplyStatusWaiting
	ApplyStatusSuccess
	ApplyStatusCancel
)

type OrderPlan struct {
	Model
	OrderID string `json:"order_id"`
	// 第几期
	Item int `json:"item"`
	// 总共几期
	TotalItem int `json:"total_item"`
	// 计划发货时间
	PlanTime int64 `json:"plan_time"`
	// 实际发货时间
	DeliveryTime int64 `json:"delivery_time"`
	// 计划发货状态
	Status PlanStatus `json:"status"`
	// 申请取消状态
	ApplyStatus ApplyStatus `json:"apply_status"`
	// 快递公司
	Express string `json:"express"`
	// 快递单号
	WaybillNumber string `json:"waybill_number"`
	// 取消申请时间
	ApplyCancelTime int64 `json:"apply_cancel_time"`
	// 确认取消时间
	CanceledTime int64 `json:"canceled_time"`
	Created      int64 `json:"created"`
}

var OrderPlanDefault OrderPlan

func (OrderPlan) TableName() string {
	return "order_plan"
}

func (op *OrderPlan) Insert(db *gorm.DB) error {
	op.Created = time.Now().Unix()
	return db.Create(op).Error
}

func (OrderPlan) GetByOrderID(orderID string) (list []*OrderPlan, err error) {
	err = cli.DB.Where("order_id = ?", orderID).Find(&list).Error
	return
}

func (OrderPlan) GetByOrderIDForUpdate(db *gorm.DB, orderID string) (list []*OrderPlan, err error) {
	err = db.Set("gorm:query_option", "FOR UPDATE").Where("order_id = ?", orderID).Find(&list).Error
	return
}

func (OrderPlan) GetByOrderIDAndItemForUpdate(db *gorm.DB, orderID string, item int) (*OrderPlan, error) {
	var data OrderPlan
	err := db.Set("gorm:query_option", "FOR UPDATE").Where("order_id = ? and item = ?", orderID, item).Find(&data).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (OrderPlan) GetByOrderIDs(orderIDs []string) (list []*OrderPlan, err error) {
	err = cli.DB.Where("order_id in (?)", orderIDs).Find(&list).Error
	return
}

func (OrderPlan) HasNoDelivery(db *gorm.DB, orderID string) (bool, error) {
	var data OrderPlan
	if err := db.Where("order_id = ?", orderID).First(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}

func (op *OrderPlan) Delivery(db *gorm.DB, express, waybillNumber string) error {
	data := map[string]interface{}{
		"express":        express,
		"waybill_number": waybillNumber,
		"delivery_time":  time.Now().Unix(),
		"status":         PlanStatusDeliveried,
	}
	return db.Model(op).Update(data).Error
}
