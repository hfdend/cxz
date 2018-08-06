package modules

import (
	"fmt"
	"testing"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/utils"
)

func TestStatistics_Statistics(t *testing.T) {
	cli.Init()
	data, err := Statistics.Statistics(0, 0)
	fmt.Println(err)
	utils.JSON(data)
}
