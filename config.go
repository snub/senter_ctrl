package main

import (
	"bytes"
	"encoding/json"
	senter "git.oneiros.ml/senter/senter.git"
	"io/ioutil"
	"os"
)

type Config struct {
	Mqtt     *Mqtt                  `json:"mqtt"`
	Database *senter.DatabaseConfig `json:"database"`
}

type Mqtt struct {
	Broker string `json:"broker"`
}

//var defaultConfig *Config

//func init() {
//defaultConfig := &Config{
//	&Mqtt{"tcp://localhost:1883"},
//	&DatabaseConfig{"localhost", "3306", "root", "qwerty123", "sensor_center"}}
//}

func LoadConfig(fileName string) (*Config, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Println("unable to open config: ", err)
		return nil, err
	}

	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		logger.Println("unable to parse config: ", err)
		return nil, err
	}
	return temp, nil
}

func (c *Config) SaveToFile(fileName string) {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		logger.Printf("unable to marchall config json: %s\n", err)
	}

	var buffer bytes.Buffer
	err = json.Indent(&buffer, jsonBytes, "", "  ")
	if err != nil {
		logger.Printf("unable to indent json: %s\n", err)
	}

	err = ioutil.WriteFile(fileName, buffer.Bytes(), os.ModePerm)
	if err != nil {
		logger.Printf("unable to write json to file: %s\n", err)
	}
}
