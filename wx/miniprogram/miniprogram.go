package miniprogram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Session struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
}

func GetSession(appID, secret, code string) (session Session, err error) {
	requestUrl := fmt.Sprintf(
		"%s?appid=%s&secret=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		"https://api.weixin.qq.com/sns/jscode2session",
		appID,
		secret,
		code,
	)
	httpClient := &http.Client{}
	httpClient.Timeout = 10 * time.Second
	var resp *http.Response
	if resp, err = httpClient.Get(requestUrl); err != nil {
		return
	}
	defer resp.Body.Close()
	var data struct {
		Session
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	var bts []byte
	if bts, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(bts, &data); err != nil {
		return
	}
	if data.ErrCode != 0 {
		err = fmt.Errorf("code: %d, msg: %s", data.ErrCode, data.ErrMsg)
		return
	}
	session = data.Session
	return
}
