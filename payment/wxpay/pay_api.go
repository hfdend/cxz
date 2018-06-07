package wxpay

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
)

const (
	URLUnifiedOrder = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	URLMicropay     = "https://api.mch.weixin.qq.com/pay/micropay"
	URLQuery        = "https://api.mch.weixin.qq.com/pay/orderquery"
	URLCancel       = "https://api.mch.weixin.qq.com/secapi/pay/reverse"
)

type PayApi struct {
	id     string
	Config PayConfig
	Logger *logrus.Logger
}

func NewApi(config PayConfig) *PayApi {
	v := new(PayApi)
	v.Config = config
	v.id = uuid.New().String()
	return v
}

// 统一下单，PayUnifiedOrder中out_trade_no、body、total_fee、trade_type必填
// appid、mchid、spbill_create_ip、nonce_str不需要填入
func (api PayApi) UnifiedOrder(args PayUnifiedOrder, timeout time.Duration) (*PayResults, error) {
	api.info(fmt.Sprintf("unified order args: %+v url: %s", args, URLUnifiedOrder))
	var (
		ok        bool
		tradeType string
	)
	if _, ok := args.GetOutTradeNo(); !ok {
		return nil, errors.New("缺少统一支付接口必填参数out_trade_no")
	}
	if _, ok := args.GetBody(); !ok {
		return nil, errors.New("缺少统一支付接口必填参数body")
	}
	if _, ok := args.GetTotalFee(); !ok {
		return nil, errors.New("缺少统一支付接口必填参数total_fee")
	}
	if tradeType, ok = args.GetTradeType(); !ok {
		return nil, errors.New("缺少统一支付接口必填参数trade_type")
	}
	if tradeType == "JSAPI" {
		if _, ok := args.GetOpenId(); !ok {
			return nil, errors.New("统一支付接口中，缺少必填参数openid！trade_type为JSAPI时，openid为必填参数")
		}
	}
	if tradeType == "NATIVE" {
		if _, ok := args.GetProductId(); !ok {
			return nil, errors.New("统一支付接口中，缺少必填参数product_id！trade_type为JSAPI时，product_id为必填参数")
		}
	}
	if _, ok := args.GetNotifyUrl(); !ok {
		args.SetNotifyUrl(api.Config.NotifyUrl)
	}
	args.SetAppId(api.Config.AppId)
	args.SetMchId(api.Config.MchId)
	if _, ok := args.GetSpbillCreateIp(); !ok {
		args.SetSpbillCreateIp("0.0.0.0")
	}
	args.SetNonceStr(GetNonceStr())
	if _, ok := args.GetSignType(); !ok {
		args.SetSignType(SignTypeMD5)
	}
	args.SetSign(api.Config.Key)
	xmlString, err := args.ToXml()
	if err != nil {
		return nil, err
	}
	data, err := api.postXmlCurl(xmlString, URLUnifiedOrder, nil, timeout)
	if err != nil {
		api.error(fmt.Sprintf("data: %v, err: %v", xmlString, err))
		return nil, err
	}
	var result PayResults
	if err := xml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 提交被扫支付API(刷卡支付)
// 收银员使用扫码设备读取微信用户刷卡授权码以后，二维码或条码信息传送至商户收银台，
// 由商户收银台或者商户后台调用该接口发起支付。
// PayWxPayMicroPay中body、out_trade_no、total_fee、auth_code参数必填
// appid、mchid、spbill_create_ip、nonce_str不需要填入
func (api PayApi) MicroPay(args PayUnifiedOrder, timeout time.Duration) (*PayResults, error) {
	api.info(fmt.Sprintf("unified order args: %+v url: %s", args, URLMicropay))
	if _, ok := args.GetBody(); !ok {
		return nil, errors.New("no set body")
	}
	if _, ok := args.GetOutTradeNo(); !ok {
		return nil, errors.New("no set out_trade_no")
	}
	if _, ok := args.GetTotalFee(); !ok {
		return nil, errors.New("not set total_fee")
	}
	if _, ok := args.GetAuthCode(); !ok {
		return nil, errors.New("not set auth_code")
	}
	if _, ok := args.GetSpbillCreateIp(); !ok {
		args.SetSpbillCreateIp("0.0.0.0")
	}
	args.SetAppId(api.Config.AppId)
	args.SetMchId(api.Config.MchId)
	args.SetNonceStr(GetNonceStr())
	if _, ok := args.GetSignType(); !ok {
		args.SetSignType(SignTypeMD5)
	}
	args.SetSign(api.Config.Key)
	xmlString, err := args.ToXml()
	if err != nil {
		return nil, err
	}
	data, err := api.postXmlCurl(xmlString, URLMicropay, nil, timeout)
	if err != nil {
		api.error(fmt.Sprintf("data: %v, err: %v", xmlString, err))
		return nil, err
	}
	var result PayResults
	if err := xml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 扫码支付
func (api PayApi) NativePay(args PayUnifiedOrder, timeout time.Duration) (string, error) {
	args.SetTradeType("NATIVE")
	args.SetNonceStr(GetNonceStr())
	args.SetAppId(api.Config.AppId)
	args.SetMchId(api.Config.MchId)
	if _, ok := args.GetSignType(); !ok {
		args.SetSignType(SignTypeMD5)
	}
	args.SetSign(api.Config.Key)
	result, err := api.UnifiedOrder(args, timeout)
	if err != nil {
		return "", err
	}
	if result.ReturnCode != "SUCCESS" {
		return "", errors.New(result.ReturnMsg)
	}
	if result.ResultCode != "SUCCESS" {
		return "", errors.New(fmt.Sprintf("%s: %s", result.ErrCode, result.ErrCodeDes))
	}
	return result.CodeUrl, nil
}

// 订单查询
func (api PayApi) OrderQuery(args PayOrderQuery, timeout time.Duration) (*QueryResults, error) {
	if _, ok := args.GetOutTradeNo(); !ok {
		if _, ok := args.GetTransactionId(); !ok {
			return nil, errors.New("no set out_trade_no and transaction_id")
		}
	}
	args.SetAppId(api.Config.AppId)
	args.SetMchId(api.Config.MchId)
	args.SetNonceStr(GetNonceStr())
	if _, ok := args.GetSignType(); !ok {
		args.SetSignType(SignTypeMD5)
	}
	args.SetSign(api.Config.Key)
	xmlString, err := args.ToXml()
	if err != nil {
		return nil, err
	}
	data, err := api.postXmlCurl(xmlString, URLQuery, nil, timeout)
	if err != nil {
		return nil, err
	}
	var result QueryResults
	if err := xml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (api PayApi) OrderCancel(args PayOrderCancel, timeout time.Duration) (*CancelResults, error) {
	if _, ok := args.GetOutTradeNo(); !ok {
		if _, ok := args.GetTransactionId(); !ok {
			return nil, errors.New("no set out_trade_no and transaction_id")
		}
	}
	args.SetAppId(api.Config.AppId)
	args.SetMchId(api.Config.MchId)
	args.SetNonceStr(GetNonceStr())
	if _, ok := args.GetSignType(); !ok {
		args.SetSignType(SignTypeMD5)
	}
	args.SetSign(api.Config.Key)
	xmlString, err := args.ToXml()
	if err != nil {
		return nil, err
	}
	tlsConfig, err := api.tlsConfig()
	if err != nil {
		return nil, err
	}
	response, err := api.postXmlCurl(xmlString, URLCancel, tlsConfig, timeout)
	if err != nil {
		return nil, err
	}
	var result CancelResults
	if err := xml.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (api PayApi) tlsConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(api.Config.SSLCertPath, api.Config.SSLKeyPath)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(api.Config.SSLRootCaPath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(data)
	if !ok {
		return nil, errors.New("failed to parse root certificate")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return tlsConfig, nil
}

func (api PayApi) postXmlCurl(xmlString, url string, tlsConfig *tls.Config, timeout time.Duration) ([]byte, error) {
	buff := bytes.NewBufferString(xmlString)
	api.info(fmt.Sprintf("request data: %v", buff))
	fmt.Println(url)
	fmt.Println(xmlString)
	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		api.error(fmt.Sprintf("new request error: %s", err.Error()))
		return nil, err
	}
	cli := &http.Client{}
	cli.Timeout = timeout
	if tlsConfig != nil {
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		}
		cli.Transport = transport
	}
	resp, err := cli.Do(req)
	if err != nil {
		api.error(fmt.Sprintf("do request error: %s", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.error(fmt.Sprintf("read response body error: %s", err.Error()))
		return nil, err
	}
	api.info(fmt.Sprintf("response data: %s", string(bts)))
	return bts, err
}

func (api PayApi) error(s string) {
	if api.Logger != nil {
		api.Logger.Error(fmt.Sprintf("[%s] %s", api.id, s))
	}
}

func (api PayApi) warn(s string) {
	if api.Logger != nil {
		api.Logger.Warn(fmt.Sprintf("[%s] %s", api.id, s))
	}
}

func (api PayApi) info(s string) {
	if api.Logger != nil {
		api.Logger.Info(fmt.Sprintf("[%s] %s", api.id, s))
	}
}

func (api PayApi) getOpenId() string {
	return ""
}

func GetNonceStr() string {
	return fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
