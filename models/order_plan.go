package models

import (
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

// 计划状态
// 1: 等待发货
// 2: 已发货
// 3: 申请取消
// 4: 已取消
// swagger:model PlanStatus
type PlanStatus int

const (
	PlanStatusWaiting PlanStatus = iota + 1
	PlanStatusDeliveried
	PlanStatusApplyCancel
	PlanStatusCanceled
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
	// 计划状态
	Status PlanStatus `json:"status"`
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

func (OrderPlan) GetByOrderIDs(orderIDs []string) (list []*OrderPlan, err error) {
	err = cli.DB.Where("order_id in (?)", orderIDs).Find(&list).Error
	return
}
