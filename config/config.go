package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// Config represents config info
var Config ConfList

// ConfList has contents of config.ini
type ConfList struct {
	DBdriver string
	DBname   string
	Port     int
	IP       string
}

// InitConfig initializes config settings
func InitConfig() {
	conf, err := ini.Load("config.ini")
	if err != nil {
		logrus.Warnf("init file open error: %v", err)
	}

	Config = ConfList{
		DBdriver: conf.Section("db").Key("driver").String(),
		DBname:   conf.Section("db").Key("name").String(),
		Port:     conf.Section("web").Key("port").MustInt(),
		IP:       conf.Section("web").Key("ip").String(),
	}
}
