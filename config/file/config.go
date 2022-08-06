package file

import (
	"encoding/json"
	"io/ioutil"
)

type appConfig struct {
	APISettings APISettings `json:"api-config"`
}

type APISettings struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func GetConfig() (*appConfig, error) {
	appConf := &appConfig{}
	if err := appConf.LoadConfigFromFile(); err != nil {
		return nil, err
	}

	return appConf, nil
}

func (c *appConfig) LoadConfigFromFile() error {
	jsonFile, err := ioutil.ReadFile("./config/file/config.json")
	if err != nil {
		panic(err)
	}
	unmarshalErr := json.Unmarshal(jsonFile, &c)
	if unmarshalErr != nil {
		return nil
	}

	return nil
}
