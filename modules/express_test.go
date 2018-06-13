package modules

import (
	"testing"

	"fmt"

	"github.com/hfdend/cxz/cli"
)

func TestExpress_Query(t *testing.T) {
	cli.Init()
	data, err := Express.Query("9973191074656", "")
	fmt.Println(err)
	fmt.Printf("%+v", data)
}
