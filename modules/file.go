package modules

import (
	"fmt"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hfdend/cxz/conf"
)

type FileCDN struct {
	Domain string `json:"domain"`
	Key    string `json:"key"`
	Url    string `json:"url"`
}

type file int

var File file

func (file) UploadToCDN(reader io.Reader, key string) (data FileCDN, err error) {
	c := conf.Config.Aliyun.OSS
	var client *oss.Client
	if client, err = oss.New(c.Endpoint, c.AccessKeyID, c.AccessKeySecret); err != nil {
		return
	}
	var bucket *oss.Bucket
	if bucket, err = client.Bucket(c.Bucket); err != nil {
		return
	}
	if err = bucket.PutObject(key, reader); err != nil {
		return
	}
	data.Domain = strings.TrimRight(c.Domain, "/")
	data.Key = key
	data.Url = fmt.Sprintf("%s/%s", data.Domain, key)
	return
}
