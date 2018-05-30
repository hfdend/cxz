package modules

import (
	"testing"

	"os"

	"fmt"

	"github.com/hfdend/cxz/cli"
)

func TestFile_UploadToCDN(t *testing.T) {
	cli.Init()
	f, _ := os.OpenFile("/Users/denghongfeng/Desktop/testimg.png", os.O_RDONLY, 0644)
	data, err := File.UploadToCDN(f, "a.png")
	fmt.Println(err)
	fmt.Println(data)
}
