package config

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server
	Log
	Kafka
	Redis
	Etcd
}

type Server struct {
	Port int
}

type Log struct {
	Name  string
	Level int
}

type Kafka struct {
	BrokerAddrs string
}

type Redis struct {
	Host string
}

type Etcd struct {
	Host string
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
