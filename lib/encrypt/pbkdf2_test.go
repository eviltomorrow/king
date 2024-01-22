package encrypt

import (
	"fmt"
	"testing"
)

func TestKey(t *testing.T) {
	v, err := Key("hello world")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}
