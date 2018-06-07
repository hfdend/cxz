package wxpay

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SignType string

const (
	SignTypeMD5  SignType = "MD5"
	SignTypeHMAC SignType = "HMAC-SHA256"
)

type DataValues map[string]interface{}

func (s DataValues) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "xml"
	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(toString(value)), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (s *DataValues) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var key, val string
	for {
		t, err := d.Token()
		if err == io.EOF { // found end of element
			break
		}
		if err != nil {
			return err
		}

		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			key = token.Name.Local
			//for _, attr := range token.Attr {
			//	attrName := attr.Name.Local
			//	attrValue := attr.Value
			//	fmt.Printf("An attribute is: %s %s\n", attrName, attrValue)
			//}
			// 处理元素结束（标签）
		case xml.EndElement:
			if token.Name == start.Name {
				break
			}
			(*s)[key] = val
			//fmt.Printf("Token of '%s' end\n", token.Name.Local)
			// 处理字符数据（这里就是元素的文本）
		case xml.CharData:
			val = string([]byte(token))
			//fmt.Printf("This is the content: %v\n", content)
		default:
			// ...
		}
	}
	return nil
}

type payDataBase struct {
	Values DataValues
}

func NewPayDataBase() *payDataBase {
	base := new(payDataBase)
	base.Values = DataValues{}
	return base
}

func (w *payDataBase) SetData(key string, value interface{}) {
	w.Values[key] = value
}

func (w *payDataBase) SetSign(key string) {
	sign := w.MakeSign(key)
	w.Values["sign"] = sign
}

// 设置是否支持信用卡
func (w *payDataBase) SetNoCredit(sure bool) {
	if sure {
		w.Values["limit_pay"] = "no_credit"
	} else {
		delete(w.Values, "limit_pay")
	}
}

func (w *payDataBase) GetSignType() (SignType, bool) {
	t, ok := w.Values["sign_type"]
	if ok {
		s := fmt.Sprintf("%v", t)
		return SignType(s), true
	} else {
		return "", false
	}
}

func (w *payDataBase) SetSignType(t SignType) {
	w.Values["sign_type"] = fmt.Sprintf("%v", t)
}

func (w *payDataBase) GetSign() string {
	s, _ := w.Values["sign"]
	return fmt.Sprintf("%v", s)
}

func (w *payDataBase) IsSignSet() bool {
	_, ok := w.Values["sign"]
	return ok
}

func (w *payDataBase) GetString(key string) (string, bool) {
	if v, ok := w.Values[key]; ok {
		return fmt.Sprintf("%v", v), true
	} else {
		return "", false
	}
}

func (w *payDataBase) GetInt(key string) (int, bool) {
	if v, ok := w.Values[key]; ok {
		if i, err := strconv.Atoi(fmt.Sprintf("%v", v)); err != nil {
			return 0, false
		} else {
			return i, true
		}
	} else {
		return 0, false
	}
}

// 生成签名
// 签名，本函数不覆盖sign成员变量，如要设置签名需要调用SetSign方法赋值
func (w *payDataBase) MakeSign(key string) string {
	t, _ := w.GetSignType()
	switch t {
	default:
		str := w.toUrlParams() + "&key=" + key
		str = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(str))))
		return str
	case SignTypeHMAC:
		str := w.toUrlParams() + "&key=" + key
		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte(str))
		str = strings.ToUpper(fmt.Sprintf("%x", mac.Sum(nil)))
		return str
	}
}

func (w *payDataBase) GetValues() DataValues {
	return w.Values
}

func (w *payDataBase) toUrlParams() string {
	buff := ""
	var keys []string
	for k := range w.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := w.Values[k]
		vv := toString(v)
		if k != "sign" && v != "" {
			buff += k + "=" + vv + "&"
		}
	}
	buff = strings.TrimRight(buff, "&")
	return buff
}

// 设置微信分配的公众账号ID
func (w *payDataBase) SetAppId(appId string) {
	w.SetData("appid", appId)
}

// 获取微信分配的公众账号ID的值
func (w *payDataBase) GetAppId() (string, bool) {
	return w.GetString("appid")
}

// 设置微信支付分配的商户号
func (w *payDataBase) SetMchId(mchId string) {
	w.SetData("mch_id", mchId)
}

// 获取微信支付分配的商户号的值
func (w *payDataBase) GetMchId() (string, bool) {
	return w.GetString("mch_id")
}

// 设置商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
func (w *payDataBase) SetOutTradeNo(tradeNo string) {
	w.SetData("out_trade_no", tradeNo)
}

// 获取商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号的值
func (w *payDataBase) GetOutTradeNo() (string, bool) {
	return w.GetString("out_trade_no")
}

// 设置微信支付分配的终端设备号，商户自定义
func (w *payDataBase) SetDeviceInfo(device string) {
	w.SetData("device_info", device)
}

// 获取微信支付分配的终端设备号，商户自定义的值
func (w *payDataBase) GetDeviceInfo() (string, bool) {
	return w.GetString("device_info")
}

// 设置随机字符串，不长于32位。推荐随机数生成算法
func (w *payDataBase) SetNonceStr(str string) {
	w.SetData("nonce_str", str)
}

// 获取随机字符串，不长于32位。推荐随机数生成算法的值
func (w *payDataBase) GetNonceStr() (string, bool) {
	return w.GetString("nonce_str")
}

// 设置商品或支付单简要描述
func (w *payDataBase) SetBody(body string) {
	w.SetData("body", body)
}

// 获取商品或支付单简要描述的值
func (w *payDataBase) GetBody() (string, bool) {
	return w.GetString("body")
}

// 设置商品名称明细列表
func (w *payDataBase) SetDetail(detail string) {
	w.SetData("detail", detail)
}

// 获取商品名称明细列表的值
func (w *payDataBase) GetDetail() (string, bool) {
	return w.GetString("detail")
}

// 设置附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
func (w *payDataBase) SetAttach(attach string) {
	w.SetData("attach", attach)
}

// 设置附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
func (w *payDataBase) GetAttach() (string, bool) {
	return w.GetString("attach")
}

// 设置符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
func (w *payDataBase) SetFeeType(freeType string) {
	w.SetData("fee_type", freeType)
}

// 获取符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型的值
func (w *payDataBase) GetFeeType() (string, bool) {
	return w.GetString("fee_type")
}

// 设置订单总金额，只能为整数，详见支付金额
func (w *payDataBase) SetTotalFee(totalFree int) {
	w.SetData("total_fee", totalFree)
}

func (w *payDataBase) GetTotalFee() (int, bool) {
	return w.GetInt("total_fee")
}

// 设置订单总金额，只能为整数，详见支付金额
func (w *payDataBase) SetAuthCode(authCode string) {
	w.SetData("auth_code", authCode)
}

func (w *payDataBase) GetAuthCode() (string, bool) {
	return w.GetString("auth_code")
}

// 设置APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
func (w *payDataBase) SetSpbillCreateIp(ip string) {
	w.SetData("spbill_create_ip", ip)
}

// 获取APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
func (w *payDataBase) GetSpbillCreateIp() (string, bool) {
	return w.GetString("spbill_create_ip")
}

// 设置订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
func (w *payDataBase) SetTimeStart(t time.Time) {
	w.SetData("time_start", t.Format("20060102150405"))
}

// 获取订单生成时间
func (w *payDataBase) GetTimeStart() (time.Time, bool) {
	var tt time.Time
	var err error
	if t, ok := w.GetString("time_start"); ok {
		if tt, err = time.ParseInLocation("20060102150405", t, time.Local); err != nil {
			return tt, false
		} else {
			return tt, true
		}
	}
	return tt, false
}

func (w *payDataBase) GetTimeEnd() (string, bool) {
	return w.GetString("time_end")
}

// 设置订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则
func (w *payDataBase) SetTimeExpire(t time.Time) {
	w.SetData("time_expire", t.Format("20060102150405"))
}

// 获取订单失效时间
func (w *payDataBase) GetTimeExpire() (time.Time, bool) {
	var tt time.Time
	var err error
	if t, ok := w.GetString("time_expire"); ok {
		if tt, err = time.ParseInLocation("20060102150405", t, time.Local); err != nil {
			return tt, false
		} else {
			return tt, true
		}
	}
	return tt, false
}

// 设置商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
func (w *payDataBase) SetGoodsTag(val string) {
	w.SetData("goods_tag", val)
}

// 获取商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠的值
func (w *payDataBase) GetGoodsTag() (string, bool) {
	return w.GetString("goods_tag")
}

// 设置接收微信支付异步通知回调地址
func (w *payDataBase) SetNotifyUrl(val string) {
	w.SetData("notify_url", val)
}

// 获取接收微信支付异步通知回调地址
func (w *payDataBase) GetNotifyUrl() (string, bool) {
	return w.GetString("notify_url")
}

// 设置取值如下：JSAPI，NATIVE，APP，详细说明见参数规定
func (w *payDataBase) SetTradeType(tradeType string) {
	w.SetData("trade_type", tradeType)
}

// 获取取值如下：JSAPI，NATIVE，APP，详细说明见参数规定的值
func (w *payDataBase) GetTradeType() (string, bool) {
	return w.GetString("trade_type")
}

// 设置trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。
func (w *payDataBase) SetProductId(productId string) {
	w.SetData("product_id", productId)
}

// 获取trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。的值
func (w *payDataBase) GetProductId() (string, bool) {
	return w.GetString("product_id")
}

// 设置trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。下单前需要调用【网页授权获取用户信息】接口获取到用户的Openid。
func (w *payDataBase) SetOpenId(openId string) {
	w.SetData("openid", openId)
}

// 获取trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。下单前需要调用【网页授权获取用户信息】接口获取到用户的Openid。 的值
func (w *payDataBase) GetOpenId() (string, bool) {
	return w.GetString("openid")
}

func (w *payDataBase) GetBankType() (string, bool) {
	return w.GetString("bank_type")
}

// 设置错误码 FAIL 或者 SUCCESS
func (w *payDataBase) SetReturnCode(code string) {
	w.SetData("return_code", code)
}

// 获取错误码 FAIL 或者 SUCCESS
func (w *payDataBase) GetReturnCode() string {
	code, _ := w.Values["return_code"]
	return fmt.Sprintf("%v", code)
}

// 获取错误码 FAIL 或者 SUCCESS
func (w *payDataBase) GetResultCode() string {
	code, _ := w.Values["result_code"]
	return fmt.Sprintf("%v", code)
}

// 设置错误信息
func (w *payDataBase) SetReturnMsg(msg string) {
	w.SetData("return_msg", msg)
}

// 获取错误信息
func (w *payDataBase) GetReturnMsg() string {
	msg, _ := w.Values["return_msg"]
	return fmt.Sprintf("%v", msg)
}

// 设置微信的订单号，优先使用
func (w *payDataBase) SetTransactionId(id string) {
	w.SetData("transaction_id", id)
}

// 获取微信的订单号，优先使用的值
func (w *payDataBase) GetTransactionId() (string, bool) {
	return w.GetString("transaction_id")
}

// 检测签名
func (w *payDataBase) CheckSign(key string) error {
	if w.IsSignSet() {
		return ErrorSign
	}
	sign := w.MakeSign(key)
	if w.GetSign() != sign {
		return ErrorSign
	}
	return nil
}

// 使用数组初始化
func (w *payDataBase) FromArray(ary map[string]interface{}) {
	w.Values = ary
}

// 使用数组初始化对象
// 如果传了key则验证签名
func (w *payDataBase) InitFromArray(ary map[string]interface{}, key string) (*payDataBase, error) {
	w.FromArray(ary)
	if key != "" {
		if err := w.CheckSign(key); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (w *payDataBase) Init(xmlData []byte) error {
	return xml.Unmarshal(xmlData, &w.Values)
}

func (w *payDataBase) ToXml() (string, error) {
	if w.Values == nil {
		return "", nil
	}
	bts, err := xml.Marshal(w.Values)
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

func round(v float64, places int) float64 {
	s := fmt.Sprintf(fmt.Sprintf("%%.%df", places), v)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func toString(v interface{}) string {
	var vv string
	switch v.(type) {
	case int8, int16, int32, int, int64, uint8, uint16, uint32, uint, uint64, string:
		vv = fmt.Sprintf("%v", v)
	case float32:
		vv = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
	case float64:
		vv = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	}
	return vv
}
