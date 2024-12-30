package flagsutil

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var Opts = &Flags{}

type Flags struct {
	ConfigFile    string `short:"c" long:"config-file" description:"specifying a config file"`
	Daemon        bool   `short:"d" long:"daemon" description:"running in background"`
	EnablePprof   bool   `long:"enable-pprof" description:"enable pprof profiling"`
	PprofAddr     string `long:"pprof-addr" default:":56060" description:"pprof listen addr"`
	DisableStdlog bool   `long:"disable-stdlog" description:"disable standard logging"`

	Version bool `short:"v" long:"version" description:"show version number"`
}

func (f *Flags) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(f)
	if err != nil {
		return fmt.Sprintf("marshal metadata failure, nest error: %v", err)
	}
	return string(buf)
}
