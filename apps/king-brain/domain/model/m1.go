package model

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
)

const NameM1 = ""

func init() {
	service.RegisterModel(NameM1, &service.Model{
		Name: NameM1,
		Desc: "",
		F:    m1,
	})
}

func m1(k *chart.K) (int, bool) {
	if len(k.Candlesticks) == 0 {
		return 0, false
	}

	last := k.Candlesticks[len(k.Candlesticks)-1]

	fmt.Println(last.String())

	if last.Close < last.MA[chart.Ma_50] || last.Close < last.MA[chart.Ma150] || last.Close < last.MA[chart.Ma200] {
		return 0, false
	}
	// score := 0
	// for i := len(k.Candlesticks) - 1; i >= 0; i-- {
	// 	current := k.Candlesticks[i]
	// }
	return 0, true
}
