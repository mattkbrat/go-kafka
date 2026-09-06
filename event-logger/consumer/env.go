package main

import (
	"encoding/json"
	"log"
	"os"
)

var KafkaRoute = "localhost:9092"
var KafkaPartition = 0

type GcpConfig struct {
	ProviderCertUrl string `json:"auth_provider_x509_cert_url"`
	AuthUri         string `json:"auth_uri"`
	ClientEmail     string `json:"client_email"`
	ClientId        string `json:"client_id"`
	ClientCertUrl   string `json:"client_x509_cert_url"`
	PrivateKey      string `json:"private_key"`
	PrivateKeyId    string `json:"private_key_id"`
	ProjectId       string `json:"project_id"`
	TokenUri        string `json:"token_uri"`
	Type            string `json:"type"`
	Domain          string `json:"universe_domain"`
}

type KafkaConfig struct {
	Url       string `json:"route"`
	Topic     string `json:"topic"`
	GroupId   string `json:"group"`
	Partition int    `json:"partition"`
}

type Config struct {
	Kafka     KafkaConfig `json:"kafka"`
	Gcp       GcpConfig
	DatasetId string `json:"datasetId"`
	TableId   string `json:"tableId"`
}

type ConfigParams struct {
	Gcp    string
	Config string
}

func readGcpConfig(file string) GcpConfig {
	configFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to read config file %s", err.Error())
	}

	gcp := GcpConfig{}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&gcp); err != nil {
		log.Fatalf("failed to parse config file %s", err.Error())
	}

	return gcp
}

func ReadConfig(params *ConfigParams) Config {
	configFile, err := os.Open(params.Config)
	if err != nil {
		log.Fatalf("failed to read config file %s", err.Error())
	}

	config := Config{}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatalf("failed to parse config file %s", err.Error())
	}

	gcp := readGcpConfig(params.Gcp)

	return Config{
		Kafka:     config.Kafka,
		Gcp:       gcp,
		DatasetId: config.DatasetId,
		TableId:   config.TableId,
	}
}
