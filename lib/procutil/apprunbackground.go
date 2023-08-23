package procutil

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/eviltomorrow/king/lib/netutil"
	jsoniter "github.com/json-iterator/go"
)

var (
	pingbackHost = "127.0.0.1"
)

type BootInfo struct {
	ChallengeKey []byte `json:"challenge_key"`
	ListenPort   int    `json:"listen_port"`
}

func (bi *BootInfo) Marshal() ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(bi)
}

func (bi *BootInfo) UnMarshal(buf []byte) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(buf, bi)
}

func RunAppBackground(name string, args []string) error {
	args = append(args, "--ppid", fmt.Sprintf("%d", os.Getpid()))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return runAppStartBackground(cmd)
}

func NotifyStopDaemon(buf []byte) error {
	var bi = &BootInfo{}
	if err := bi.UnMarshal(buf); err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", pingbackHost, bi.ListenPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("dialing confirmation address: %v", err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte(bi.ChallengeKey))
	if err != nil {
		return fmt.Errorf("writing confirmation bytes to %s, nest error: %v", address, err)
	}
	return nil
}

func runAppStartBackground(cmd *exec.Cmd) error {
	port, err := netutil.GetAvailablePort()
	if err != nil {
		return fmt.Errorf("get available port failure, nest error: %v", err)
	}
	address := fmt.Sprintf("%s:%d", pingbackHost, port)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen %s failure, nest error: %v", address, err)
	}
	defer ln.Close()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return err
	}
	go func() {
		var bi = &BootInfo{
			ChallengeKey: key,
			ListenPort:   port,
		}
		buf, err := bi.Marshal()
		if err != nil {
			log.Fatalf("[F] Marshal boot-info failure, nest error: %v", err)
		}
		_, _ = stdin.Write(buf)
		stdin.Close()
	}()

	if err := cmd.Start(); err != nil {
		return err
	}

	success, exit := make(chan struct{}), make(chan error)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Fatalf("[F] Accept tcp data failure, nest error: %v", err)
				}
				break
			}
			err = handlePingbackConn(conn, key)
			if err == nil {
				close(success)
				break
			}
			log.Fatalf("[F] Parse pingback data failure, nest error: %v", err)
		}
	}()

	go func() {
		err := cmd.Wait()
		exit <- err
	}()

	select {
	case <-success:
		printRunning(cmd.Process.Pid)
	case err := <-exit:
		return fmt.Errorf("process exited with error: %v", err)
	}
	return nil
}

func handlePingbackConn(conn net.Conn, expect []byte) error {
	defer conn.Close()
	key, err := io.ReadAll(io.LimitReader(conn, 32))
	if err != nil {
		return err
	}
	if !bytes.Equal(key, expect) {
		return fmt.Errorf("wrong challenge key: %x", key)
	}
	return nil
}
