package smtp

import (
	"fmt"
	"testing"

	"github.com/eviltomorrow/king/app/king-email/conf"
)

func TestSendEmail(t *testing.T) {
	var path = "/home/shepard/data/king/conf/king-email/etc/smtp.json"
	smtp, err := conf.FindSMTP(path)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(smtp.String())

	if err := SendWithSSL(smtp.Server, smtp.Username, smtp.Password, &Message{
		From: Contact{
			Name:    smtp.Alias,
			Address: smtp.Username,
		},
		To: []Contact{
			{Name: "shepard", Address: "eviltomorrow@163.com"},
		},
		Subject:     "This is one test",
		Body:        "Hello world",
		ContentType: TextPlain,
	}); err != nil {
		t.Fatal(err)
	}
}
