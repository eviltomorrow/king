package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/procutil"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/spf13/cobra"
)

var StopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running app",
	Run: func(cmd *cobra.Command, args []string) {
		var pidFile = filepath.Join(system.Directory.VarDir, fmt.Sprintf("/run/%s.pid", buildinfo.AppName))
		if err := procutil.StopProcessWithPidFile(pidFile); err != nil {
			log.Fatalf("[F] Stop process with pidfile failure, nest error: %v", err)
		}
	},
}
