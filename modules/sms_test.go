package modules

import (
	"testing"

	"github.com/hfdend/cxz/cli"
)

func TestSms_Send(t *testing.T) {
	cli.Init()
	SMS.Send("18111634003", "1234")
}
