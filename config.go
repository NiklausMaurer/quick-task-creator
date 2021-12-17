package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DefaultListId string
}

func GetConfig(configFilePath string) (Config, error) {

	conf := Config{}

	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return conf, err
	}
	if err = json.Unmarshal(b, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
