package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Mqtt     *Mqtt           `json:"mqtt"`
	Database *DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
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
