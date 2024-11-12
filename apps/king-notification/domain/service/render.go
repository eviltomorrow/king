package service

import (
	"fmt"
	"path/filepath"

	"github.com/eviltomorrow/king/lib/system"
	"github.com/flosch/pongo2/v6"
)

func GenerateMarketStatusHTMLText(name string, data map[string]any) (string, error) {
	tpl, err := pongo2.FromFile(filepath.Join(system.Directory.EtcDir, "assets", name))
	if err != nil {
		return "", fmt.Errorf("load daily-report.html failure, nest error: %v", err)
	}

	return tpl.Execute(data)
}
