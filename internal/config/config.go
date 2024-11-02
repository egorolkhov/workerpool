package config

import (
	"encoding/json"
	"io/ioutil"
)

const filename = "./tmp/config.json"

type Config struct {
	Tasks       []string `json:"tasks"`
	WorkerCount int      `json:"worker_count"`
}

func NewConfig() (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
