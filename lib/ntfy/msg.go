package ntfy

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/eviltomorrow/king/lib/http/client"
	jsoniter "github.com/json-iterator/go"
)

func Send(ntfyAddress string, username, password string, topic string, msg *Msg) (string, error) {
	return send(ntfyAddress, username, password, topic, msg)
}

func send(ntfyAddress string, username, password string, topic string, msg *Msg) (string, error) {
	var (
		url  = fmt.Sprintf("%s/%s", ntfyAddress, topic)
		auth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	)

	header := map[string]string{
		"Title":         msg.Title,
		"Priority":      fmt.Sprintf("%d", msg.Priority),
		"Tags":          strings.Join(msg.Tags, ","),
		"Attach":        msg.Attach,
		"Authorization": fmt.Sprintf("Basic %s", auth),
	}

	return client.Post(url, 20*time.Second, header, strings.NewReader(msg.Message))
}

// Message msg
type Msg struct {
	Title   string
	Message string

	Priority int64
	Tags     []string
	Attach   string
	// Actions  []Action
}

type Action struct {
	Action  string
	Label   string
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

func (m *Msg) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return string(buf)
}
