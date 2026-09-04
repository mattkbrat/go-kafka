package main

import (
	"encoding/json"
	"log"
	"os"
)

type KafkaConfig struct {
	Url       string `json:"route"`
	Topic     string `json:"topic"`
	GroupId   string `json:"group"`
	Partition int    `json:"partition"`
}

type Config struct {
	Kafka KafkaConfig `json:"kafka"`
}

func ReadConfig(file string) Config {
	configFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to read config file %s", err.Error())
	}

	config := Config{}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatalf("failed to parse config file %s", err.Error())
	}

	return config
}
