package wxpay

type PayOrderCancel interface {
	ToXml() (string, error)
	SetSignType(key SignType)      // 签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
	GetSignType() (SignType, bool) // 签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
	SetAppId(string)
	GetAppId() (string, bool)
	SetMchId(string)
	GetMchId() (string, bool)
	SetTransactionId(string)
	GetTransactionId() (string, bool)
	SetOutTradeNo(string)
	GetOutTradeNo() (string, bool)
	SetNonceStr(str string)      // 设置随机字符串，不长于32位。推荐随机数生成算法
	GetNonceStr() (string, bool) // 获取随机字符串，不长于32位。推荐随机数生成算法的值
	SetSign(key string)          // 设置签名，详见签名生成算法
}

func NewPayOrderCancel() PayOrderCancel {
	return NewPayDataBase()
}
