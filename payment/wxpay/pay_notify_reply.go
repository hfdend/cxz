package wxpay

// 回调基础类
type PayNotifyReply interface {
	SetReturnCode(code string) // 设置错误码 FAIL 或者 SUCCESS
	GetReturnCode() string     // 获取错误码 FAIL 或者 SUCCESS
	SetReturnMsg(msg string)   // 设置错误信息
	GetReturnMsg() string      // 获取错误信息
}

func NewPayNotifyReply() PayNotifyReply {
	return NewPayDataBase()
}
