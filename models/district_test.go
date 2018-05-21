package models

import (
	"testing"

	"fmt"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/utils"
)

func TestDistrict_GetGradation(t *testing.T) {
	cli.Init()
	list, err := DistrictDefault.GetGradation()
	fmt.Println(err)
	utils.JSON(list)
}

func TestDistrict_GetNames(t *testing.T) {
	cli.Init()
	list, err := DistrictDefault.GetNames("511321")
	fmt.Println(err)
	fmt.Println(list)
}
