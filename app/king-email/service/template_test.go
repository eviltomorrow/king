package service

import (
	"fmt"
	"log"
	"testing"

	"github.com/eviltomorrow/king/lib/system"
	"github.com/flosch/pongo2/v6"
)

func TestRenderTemplate(t *testing.T) {
	system.Runtime.RootDir = ".."
	if err := InitTemplates(); err != nil {
		log.Fatal(err)
	}
	tpl, ok := GetTemplate("email-brief.html")
	if !ok {
		log.Fatal("no template")
	}
	content, err := tpl.Execute(pongo2.Context{
		"user":     "Shepard",
		"content1": "同步数据失败",
		"content2": "原因",
		"end":      "Bye",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(content)
}
