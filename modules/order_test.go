package modules

import (
	"fmt"
	"testing"

	"github.com/hfdend/cxz/cli"
)

func TestOrder_Build(t *testing.T) {
	cli.Init()
	order, err := Order.Build(3, 2, []OrderProductInfo{
		{
			ProductID: 10,
			Number:    1,
		},
	})
	fmt.Println(err)
	fmt.Println(order)
}
