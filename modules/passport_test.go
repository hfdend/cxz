package modules

import (
	"fmt"
	"testing"

	"github.com/hfdend/cxz/cli"
)

func TestPassport_SendRegisterCode(t *testing.T) {
	cli.Init()
	code, err := Passport.SendRegisterCode("15198200257")
	fmt.Println(err)
	fmt.Println(code)
}
