package v1

import (
	"github.com/hfdend/cxz/handler"
)

var JSON = handler.JSON
var Bind = handler.Bind

// 成功返回 success 字符串
// swagger:response SUCCESS
type Success string

const (
	SUCCESS Success = "success"
)
