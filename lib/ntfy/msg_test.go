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
	fmt.Println(ntfy.String())

	data, err := Send(fmt.Sprintf("%s://%s:%d", ntfy.Scheme, ntfy.Server, ntfy.Port), ntfy.Username, ntfy.Password, "", &Msg{
		Title:    "Hi",
		Message:  "This is shepard",
		Priority: 3,
		Tags:     []string{"cow", "bear"},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
}
