package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	DefaultListId string
}

func GetConfig() (Config, error) {

	conf := Config{}

	b, err := ioutil.ReadFile(getDefaultConfigFilePath())
	if err != nil {
		return conf, err
	}
	if err = json.Unmarshal(b, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func SetConfig(config Config) error {

	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(getDefaultConfigFilePath(), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getDefaultConfigFilePath() string {
	homeDirPath := os.Getenv("HOME")
	return fmt.Sprintf("%s/.quick-task-creator/%s", homeDirPath, "config.json")
}
