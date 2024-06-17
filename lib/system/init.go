package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/eviltomorrow/king/lib/netutil"
	"github.com/eviltomorrow/king/lib/timeutil"
)

func InitRuntime() error {
	executePath, err := os.Executable()
	if err != nil {
		return err
	}
	executePath, err = filepath.Abs(executePath)
	if err != nil {
		return err
	}

	Directory.BinDir = filepath.Dir(executePath)
	if !strings.HasPrefix(Directory.BinDir, "/bin") {
		Directory.RootDir = filepath.Dir(Directory.BinDir)
	} else {
		Directory.RootDir = Directory.BinDir
	}
	Directory.EtcDir = filepath.Join(Directory.RootDir, "/etc")
	Directory.UsrDir = filepath.Join(Directory.RootDir, "/usr")
	Directory.VarDir = filepath.Join(Directory.RootDir, "/var")

	Process.Name = filepath.Base(executePath)
	Process.Args = os.Args[1:]
	Process.Pid = os.Getpid()
	Process.PPid = os.Getppid()

	ipv4, err := netutil.GetInterfaceIPv4First()
	if err != nil {
		return err
	}
	Network.BindIP = ipv4

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	Machine.Hostname = hostname

	now := time.Now()
	Runtime.BootupTime = now
	Runtime.RunningDuration = func() string {
		return timeutil.FormatDuration(time.Since(now))
	}
	Runtime.OS = runtime.GOOS
	Runtime.ARCH = runtime.GOARCH

	return nil
}
