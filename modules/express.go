package modules

import (
	"fmt"
	"net/http"
	"net/url"

	"time"

	"io/ioutil"

	"encoding/json"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
)

type express int

var Express express

func (express) Query(number, company string) (data models.ExpressData, err error) {
	key := fmt.Sprintf("express_cache_%s_%s", number, company)
	err = models.GetByRedis(cli.Redis, key, time.Hour, &data, func() (dataRaw []byte, err error) {
		var data models.ExpressData
		defer func() {
			dataRaw, err = json.Marshal(data)
		}()
		var u *url.URL
		if u, err = url.Parse("https://wuliu.market.alicloudapi.com/kdi"); err != nil {
			return
		}
		values := url.Values{}
		values.Set("no", number)
		u.RawQuery = values.Encode()
		var req *http.Request
		if req, err = http.NewRequest("GET", u.String(), nil); err != nil {
			return
		}
		appCode := conf.Config.Aliyun.Express.AppCode
		req.Header.Set("Authorization", fmt.Sprintf("APPCODE %s", appCode))
		c := &http.Client{}
		c.Timeout = 15 * time.Second
		var resp *http.Response
		if resp, err = c.Do(req); err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			err = errors.New("Authentication failure")
			return
		}
		var raw []byte
		if raw, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}
		var result struct {
			Status  string          `json:"status"`
			Message string          `json:"msg"`
			Result  json.RawMessage `json:"result"`
		}
		if err = json.Unmarshal(raw, &result); err != nil {
			return
		}
		if result.Status != "0" {
			err = errors.New(result.Message)
			return
		}
		if err = json.Unmarshal(result.Result, &data); err != nil {
			return
		}
		return
	})
	return
}
