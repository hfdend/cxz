package models

import (
	"fmt"
	"testing"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/utils"
)

func TestProduct_GetList(t *testing.T) {
	cli.Init()
	cond := ProductCondition{}
	cond.Type = "狗狗湿粮"
	cond.Taste = "土豆牛肉"
	pager := NewPager(1, 20)
	list, err := ProductDefault.GetList(cond, pager)
	fmt.Println(err)
	utils.JSON(list)
}
