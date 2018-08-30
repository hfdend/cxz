package modules

import (
	"fmt"

	"github.com/GiterLab/aliyun-sms-go-sdk/dysms"
	"github.com/google/uuid"
)

type sms int

var SMS sms

func (sms) Send(phone, code string) error {
	dysms.SetACLClient("LTAIEhBUQ1g5O1ug", "UJp0JlYw2OVoNG2Uoteb8UYX9Fdo12")
	respSendSms := dysms.SendSms(uuid.New().String(), phone, "小馋主", "SMS_137656511", fmt.Sprintf(`{"code":"%s"}`, code))
	v, c, err := respSendSms.Request.Do("")
	fmt.Println(err)
	fmt.Println(c)
	fmt.Println(string(v))
	return nil
}
