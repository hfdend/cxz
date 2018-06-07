package wxpay

import (
	"fmt"
	"testing"
)

func TestPayDataBase_ToXml(t *testing.T) {
	data := NewPayDataBase()
	data.SetMchId("acb")
	data.SetTotalFee(111)

	d := NewPayDataBase()

	ss, _ := data.ToXml()
	fmt.Println(ss)
	err := d.Init([]byte(ss))
	fmt.Println(err)
	fmt.Println(d.Values)
}

func TestPayDataBase_GetSign(t *testing.T) {
	data := NewPayDataBase()
	data.Values["appid"] = "wx0fe121c961c7b4c8"
	data.Values["partnerid"] = "1473208202"
	data.Values["prepayid"] = "wx2018011115362649b74244050919735240"
	data.Values["package"] = "Sign=WXPay"
	data.Values["timestamp"] = "1515656186"
	data.SetNonceStr("864849E4FDBE80DAD891F6EC4D336C9D")
	data.SetSign("snailpetggloevedogscatsweixinpay")
	fmt.Printf("%+v\n", data)
}
