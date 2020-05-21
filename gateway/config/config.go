/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package config

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server
	Log
	Channel
}

type Server struct {
	Port    string
	Name    string
	Network string
	Proto   string
	RpcPort string
}

type Log struct {
	Name  string
	Level int
}

type Channel struct {
	BucketNum     int
	OneChannelNum int
	AppChannelNum int
	AllChannelNum int
}

var (
	onceDo sync.Once
)

func Init(filePath string) *Config {
	c := Config{}
	onceDo.Do(func() {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(data, &c); err != nil {
			panic(err)
		}
	})
	return &c
}
