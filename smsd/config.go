package smsd

import (
	"encoding/json"
	"log"
	"os"
)

type ApiConfig struct {
	CustomerKey       string `json:"customer_key"`
	CustomerKeySecret string `json:"customer_key_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

type BeanstalkConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	Api       ApiConfig       `json:"api"`
	Beanstalk BeanstalkConfig `json:"beanstalk"`
	DBUrl     string          `json:"database_url"`
}

func GetConfig(filename string) (*Config, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		log.Printf("Cannot read config file: %v", err)
		return nil, err
	}
	dec := json.NewDecoder(configFile)
	config = new(Config)
	err = dec.Decode(config)
	if err != nil {
		log.Printf("Cannot decode config file: %v", err)
		return nil, err
	}
	return config, nil
}
