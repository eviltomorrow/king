package fs

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/eviltomorrow/king/lib/buildinfo"
)

var (
	stderrFileHandler *os.File

	StderrFilePath = fmt.Sprintf("/var/log/%s/panic.log", buildinfo.AppName)
)

func RewriteStderrFile() error {
	file, err := os.OpenFile(StderrFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	stderrFileHandler = file

	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		return err
	}
	runtime.SetFinalizer(stderrFileHandler, func(fd *os.File) {
		fd.Close()
	})

	return nil
}
