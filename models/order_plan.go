package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
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
type ApplyStatus int

const (
	ApplyStatusNil ApplyStatus = iota + 1
	ApplyStatusWaiting
	ApplyStatusSuccess
)

type OrderPlan struct {
	Model
	OrderID string `json:"order_id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	// 第几期
	Item int `json:"item"`
	// 金额
	Price float64 `json:"price"`
	// 运费
	Freight float64 `json:"freight"`
	// 退款金额
	RefundAmount float64 `json:"refund_amount"`
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
	// 发货的后台人员ID
	AdminUserID int `json:"admin_user_id"`
	// 第一个商品的图片
	Image string `json:"image"`
	// 取消申请时间
	ApplyCancelTime int64 `json:"apply_cancel_time"`
	// 确认取消时间
	CanceledTime int64 `json:"canceled_time"`
	Created      int64 `json:"created"`
	Updated      int64 `json:"updated"`

	ImageSrc string `json:"image_src" gorm:"-"`
}

var OrderPlanDefault OrderPlan

func (OrderPlan) TableName() string {
	return "order_plan"
}

func (op *OrderPlan) Insert(db *gorm.DB) error {
	op.Created = time.Now().Unix()
	op.Updated = time.Now().Unix()
	return db.Create(op).Error
}

// 按需要发货的期数获取期
func (op *OrderPlan) GetNeedSendList(pager *Pager) (list []*OrderPlan, err error) {
	db := cli.DB.Model(OrderPlan{})
	db = db.Where("status = ? and apply_status = ?", PlanStatusWaiting, ApplyStatusNil).Order("plan_time asc")
	if pager != nil {
		if db, err = pager.Exec(db); err != nil {
			return
		} else if pager.Count == 0 {
			return
		}
	}
	if err = db.Find(&list).Error; err != nil {
		return
	}
	return
}

func (OrderPlan) GetByOrderID(orderID string) (list []*OrderPlan, err error) {
	if err = cli.DB.Where("order_id = ?", orderID).Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (OrderPlan) GetByOrderIDForUpdate(db *gorm.DB, orderID string) (list []*OrderPlan, err error) {
	if err = db.Set("gorm:query_option", "FOR UPDATE").Where("order_id = ?", orderID).Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (OrderPlan) GetByOrderIDAndUserIDForUpdate(db *gorm.DB, orderID string, userID int) (list []*OrderPlan, err error) {
	if err = db.Set("gorm:query_option", "FOR UPDATE").Where("order_id = ? and user_id = ?", orderID, userID).Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (OrderPlan) Delay(db *gorm.DB, ids []int, diff int64) error {
	data := map[string]interface{}{
		"plan_time": gorm.Expr("plan_time + ?", diff),
		"updated":   time.Now().Unix(),
	}
	return db.Model(OrderPlan{}).Update(data).Error
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
	data.SetImageSrc()
	return &data, nil
}

func (OrderPlan) GetByOrderIDAndItemsForUpdate(db *gorm.DB, orderID string, items []int) (list []*OrderPlan, err error) {
	err = db.Set("gorm:query_option", "FOR UPDATE").Where("order_id = ? and item in (?)", orderID, items).Find(&list).Error
	if err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (OrderPlan) CancelOrder(db *gorm.DB, ids []int) error {
	data := map[string]interface{}{
		"apply_status":  ApplyStatusSuccess,
		"canceled_time": time.Now().Unix(),
		"updated":       time.Now().Unix(),
	}
	return db.Model(OrderPlan{}).Where("id in (?)", ids).Update(data).Error
}

func (OrderPlan) GetByOrderIDs(orderIDs []string) (list []*OrderPlan, err error) {
	if err = cli.DB.Where("order_id in (?)", orderIDs).Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
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

func (OrderPlan) GetByUserID(userID int) (list []*OrderPlan, err error) {
	if err = cli.DB.Where("user_id = ?", userID).Order("plan_time asc, id asc").Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (op *OrderPlan) Delivery(db *gorm.DB, express, waybillNumber string, adminUserID int) error {
	data := map[string]interface{}{
		"express":        express,
		"waybill_number": waybillNumber,
		"delivery_time":  time.Now().Unix(),
		"updated":        time.Now().Unix(),
		"status":         PlanStatusDeliveried,
		"admin_user_id":  adminUserID,
	}
	return db.Model(op).Update(data).Error
}

func (op *OrderPlan) SetImageSrc() {
	if op.Image == "" {
		return
	}
	c := conf.Config.Aliyun.OSS
	op.ImageSrc = fmt.Sprintf("%s/%s", strings.TrimRight(c.Domain, "/"), strings.TrimLeft(op.Image, "/"))
	return
}
