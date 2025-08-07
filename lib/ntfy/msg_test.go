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
