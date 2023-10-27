package service

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/eviltomorrow/king/lib/system"
	"github.com/flosch/pongo2/v6"
)

var templates = make(map[string]*pongo2.Template, 8)

func InitTemplates() error {
	templateDir := filepath.Join(system.Runtime.RootDir, "etc/template")
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		var name = entry.Name()
		if strings.HasSuffix(name, ".html") {
			tpl, err := pongo2.FromFile(filepath.Join(templateDir, name))
			if err != nil {
				return err
			}
			templates[name] = tpl
		}
	}
	return nil
}

func GetTemplate(name string) (*pongo2.Template, bool) {
	if name == "" {
		for _, v := range templates {
			return v, true
		}
	}
	tpl, ok := templates[name]
	return tpl, ok
}
