package ntfy

import (
	"fmt"
	"testing"

	"github.com/eviltomorrow/king/apps/king-notification/conf"
)

func TestSend(t *testing.T) {
	path := "../../apps/king-notification/conf/etc/ntfy.json"
	ntfy, err := conf.LoadNTFY(path)
	if err != nil {
		t.Fatal(err)
	}

	data, err := Send(fmt.Sprintf("%s://%s:%d", ntfy.Scheme, ntfy.Server, ntfy.Port), ntfy.Username, ntfy.Password, "SrxOPwCBiRWZUOq0", &Msg{
		Title: "Hi",
		Message: `
亲爱的 Liarsa,	
以下是你收到的一份简报:	
	上证指数：⬇️ 3286.41 (-1.07%)
	深证成指：⬇️ 10543.33 (-1.32%)
	创业板指：⬇️ 2177.31 (-2.31%)
	科创50:  ⬇️ 975.37 (-1.12%)
	北证50:  ⬆️ 1317.41 (3.68%)

其中，今日正常交易的股票数量为 5393 家，上涨 1038 家，占比 19.25%, 下跌 4281 家，占比 79.38%。同时，涨幅超过 7% 的有 246 家，占比 4.56%；跌幅超过 7% 的有 99 家，占比 1.84%。	

问候,
The King Team.

金钱永不眠！`,
		Priority: 3,
		Tags:     []string{"简报", "股票"},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
}
