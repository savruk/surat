package surat

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type ConfigJson struct {
	Backend Backend
}

type Backend struct {
	Protocol string
	Host     string
	Port     uint
}

type Config struct {
	Path    string
	Content []byte
}

func (c *Config) ReadConfig() {
	b, err := ioutil.ReadFile(c.Path)
	if err != nil {
		panic(err)
	}
	c.Content = b
}

func (c *Config) Parse() ConfigJson {
	c.ReadConfig()
	var jsontype ConfigJson
	err := json.Unmarshal(c.Content, &jsontype)
	if err != nil {
		panic(err)
	}
	return jsontype
}

func LoadConfig(path string) ConfigJson {
	abspath, _ := filepath.Abs(path)
	c := &Config{
		Path: abspath,
	}
	parsed := c.Parse()
	return parsed
}
