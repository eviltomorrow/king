package system

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	Runtime   _runtime
	Machine   _machine
	Network   _network
	Process   _process
	Directory _directory
)

type _runtime struct {
	BootupTime      time.Time
	RunningDuration func() string

	OS   string
	ARCH string
}

type _machine struct {
	Hostname string
}

type _network struct {
	BindIP string
	NatIP  string
}

type _process struct {
	Name string
	Args []string

	Pid  int
	PPid int
}

type _directory struct {
	RootDir string

	BinDir string
	EtcDir string
	UsrDir string
	VarDir string
}

func String() string {
	var data = map[string]interface{}{
		"machine": Machine,
		"network": Network,
		"process": Process,
	}

	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(data)
	return string(buf)

}
