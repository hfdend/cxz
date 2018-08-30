package modules

import (
	"fmt"
	"testing"

	"github.com/hfdend/cxz/utils"

	"github.com/hfdend/cxz/models"

	"github.com/hfdend/cxz/cli"
)

func TestProduct_GetList(t *testing.T) {
	cli.Init()
	// type=狗狗湿粮&taste=三文鱼美毛&page=1&min_weight=2&max_weight=5&min_age=10&max_age=9999&is_plan=1
	cond := models.ProductCondition{}
	cond.Type = "狗狗湿粮"
	cond.Taste = "三文鱼美毛"
	cond.MinWeight = 0
	cond.MaxWeight = 2
	cond.MinAge = 10
	cond.MaxAge = 9999
	list, err := Product.GetList(cond, nil)
	fmt.Println(err)
	utils.JSON(list)
}
