package modules

import (
	"fmt"
	"testing"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/models"
)

func TestOrder_Build(t *testing.T) {
	cli.Init()
	order, err := Order.Build(3, 2, []OrderProductInfo{
		{
			ProductID: 10,
			Number:    1,
		},
	}, "快点发货", 3)
	fmt.Println(err)
	fmt.Println(order)
}

func TestOrder_WXAPay(t *testing.T) {
	cli.Init()
	user, _ := models.UserDefault.GetByID(3)
	fmt.Println(user)
	res, err := Order.WXAPay("2018060417030300011393", user, "127.0.0.1")
	fmt.Println(err)
	fmt.Println(res)
}

func TestOrder_PaymentSuccess(t *testing.T) {
	cli.Init()
	Order.PaymentSuccess("2018060810504500016033", "4201000147201806075200163613")
}
