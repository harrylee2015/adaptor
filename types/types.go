package types

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

//配置文件模板
type Config struct {
	JsonRpc []string `yaml:"jsonrpc"`
}

func ReadConf(path string) *Config {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var c Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &c
}
