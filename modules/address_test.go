package modules

import (
	"testing"

	"fmt"

	"github.com/hfdend/cxz/cli"
)

func TestAddress_Save(t *testing.T) {
	cli.Init()
	addr, err := Address.Save(0, 3, "邓鸿风", "18111634003", "511321", "二环路1号")
	fmt.Println(err)
	fmt.Println(addr)
}
