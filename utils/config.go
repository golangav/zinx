package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

type ConfigObject struct {
	Name           string
	TcpServer      ziface.IServer
	Host           string
	Port           int
	Version        string
	MaxConn        int
	MaxPackageSize uint32
}

// 定义一个全局的对外ConfigObj
var ConfigObj *ConfigObject

func (c *ConfigObject) Reload() {
	data, err := ioutil.ReadFile("../config/config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &ConfigObj)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，出初始化ConfigObj
func init() {
	ConfigObj = &ConfigObject{
		Name:           "zinx server",
		Version:        "0.5",
		Host:           "127.0.0.1",
		Port:           8009,
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	ConfigObj.Reload()
}
