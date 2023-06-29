//go:build darwin

package procutil

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	HomeDir = "."
)

func BackgroundStartup(name string, args []string, reader io.Reader, writer io.WriteCloser) error {
	var data = make([]string, 0, len(args)+1)
	data = append(data, name)
	data = append(data, args...)
	var cmd = &exec.Cmd{
		Path:  name,
		Args:  data,
		Stdin: reader,
	}
	cmd.Dir = HomeDir
	cmd.Env = os.Environ()

	if writer == nil {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var (
		ch     = make(chan os.Signal, 1)
		errsig = make(chan error, 1)
	)
	go func() {
		errsig <- cmd.Wait()
	}()

	signal.Notify(ch, os.Interrupt, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case err := <-errsig:
		if err != nil {
			return err
		}
		printOK(cmd.Process.Pid)
		return nil

	case sig := <-ch:
		switch sig {
		case syscall.SIGUSR1:
			printOK(cmd.Process.Pid)
		default:
		}
		return nil
	}
}
