package cli

import (
	"sync"

	"github.com/hfdend/cxr/conf"
)

var once sync.Once

// Init 初始化操作
func Init() {
	once.Do(func() {
		conf.Init()
		InitializeLogger()
		InitMysql()
		InitializeRedis()
	})
}
