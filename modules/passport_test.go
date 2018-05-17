package modules

import (
	"testing"
	"github.com/hfdend/cxr/cli"
	"fmt"
)

func TestPassport_SendRegisterCode(t *testing.T) {
	cli.Init()
	code, err := Passport.SendRegisterCode("18111634003")
	fmt.Println(err)
	fmt.Println(code)
}