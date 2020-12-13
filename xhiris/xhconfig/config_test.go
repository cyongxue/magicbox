package xhconfig

import (
	"fmt"
	"github.com/cyongxue/magicbox/xhiris/xhcrypt"
	"log"
	"testing"
)

// 服务的配置，启动只有一个
var SrvConfig = ServerConfig{
	Log:   &Log{},
	PProf: &PProf{},
}

type ServerConfig struct {
	Log   *Log
	PProf *PProf
}

type Log struct {
	Dir     string
	Prefix  string
	Level   string
	Console bool
}

func (l *Log) Parse(container ConfContainer, runMode string) error {
	l.Dir = container.String(fmt.Sprintf("%s::%s", runMode, "log.dir"), xhcrypt.ConfigAes)
	l.Prefix = container.String(fmt.Sprintf("%s::%s", runMode, "log.prefix"), xhcrypt.ConfigAes)
	l.Level = container.String(fmt.Sprintf("%s::%s", runMode, "log.level"), xhcrypt.ConfigAes)
	l.Console = container.DefaultBool(fmt.Sprintf("%s::%s", runMode, "log.console"), true)
	return nil
}

type PProf struct {
	PProfEnable bool
	PProfAddr   string
}

func (p *PProf) Parse(container ConfContainer, runMode string) error {
	p.PProfEnable = container.DefaultBool(fmt.Sprintf("%s::%s", runMode, "pprof.enable"), true)
	p.PProfAddr = container.String(fmt.Sprintf("%s::%s", runMode, "pprof.addr"), xhcrypt.ConfigAes)
	return nil
}

func TestConfigParse(t *testing.T) {
	key := "lmN4dtPyeC5r29DYBLl0P0OoA4Afy/2UnCg0zd+hHhg="
	oldKeys := make(map[string][]byte)
	if err := xhcrypt.Init("=DhaV=", []byte(key), oldKeys); err != nil {
		fmt.Println(err.Error())
		return
	}

	fileName := "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhconfig/config/config.conf"
	if err := ConfigParse(&SrvConfig, fileName); err != nil {
		log.Print(fmt.Errorf("config load error, configFile=%s, msg=%s", "config.conf", err.Error()))
		return
	}
}
