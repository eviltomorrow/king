package metadata

import (
	"path/filepath"

	"github.com/eviltomorrow/king/apps/king-collector/conf"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "metadata",
	Short: "Metadata root operation",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cfg = conf.Default

func loadConfig() error {
	return cfg.LoadFile(filepath.Join(system.Runtime.RootDir, "/etc/config.toml"))
}
