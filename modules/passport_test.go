package modules

import (
	"fmt"
	"testing"

	"github.com/hfdend/cxz/cli"
)

func TestPassport_SendRegisterCode(t *testing.T) {
	cli.Init()
	code, err := Passport.SendRegisterCode("18111634003")
	fmt.Println(err)
	fmt.Println(code)
}
