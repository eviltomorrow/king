package procutil

import (
	"fmt"
	"os"

	"github.com/eviltomorrow/king/lib/fs"
)

func CreatePidFile(path string, pid int) (func() error, error) {
	file, err := fs.NewFlockFile(path)
	if err != nil {
		return nil, err
	}

	file.WriteString(fmt.Sprintf("%d", pid))
	if err := file.Sync(); err != nil {
		file.Close()
		return nil, err
	}

	return func() error {
		if file != nil {
			if err := file.Close(); err != nil {
				return err
			}
			return os.Remove(path)
		}
		return fmt.Errorf("panic: pid file is nil")
	}, nil
}
