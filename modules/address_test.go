package modules

import (
	"testing"

	"fmt"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/models"
)

func TestAddress_Save(t *testing.T) {
	cli.Init()
	addr, err := Address.Save(0, 3, "邓鸿风", "18111634003", "511321", "二环路1号", models.SureYes)
	fmt.Println(err)
	fmt.Println(addr)
}
