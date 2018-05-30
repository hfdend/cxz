package miniprogram

import (
	"fmt"
	"testing"
)

func TestGetSession(t *testing.T) {
	s, err := GetSession("wx56fb16e23ab0442c", "f8e756dfc1ef1a9b9c0ee00e54b71cdd", "0711bZof1rBAmz0Dzcof1lKMof11bZon")
	fmt.Println(err)
	fmt.Println(s)
}
