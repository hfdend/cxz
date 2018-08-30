package modules

import (
	"fmt"
	"testing"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
)

func TestSms_Send(t *testing.T) {
	cli.Init()
	fmt.Println(conf.Config.Aliyun.SMS.Test)
	//SMS.Send("18111634003", "1234")
}
