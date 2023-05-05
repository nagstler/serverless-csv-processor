package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	HeaderMapping map[string]string `json:"header_mapping"`
	SQSQueueURL   string            `json:"sqs_queue_url"`
	ArchiveFolder string            `json:"archive_folder"`
}

func LoadConfig(reader io.Reader) (*Config, error) {
	var config Config
	err := json.NewDecoder(reader).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadConfigFromFile(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return LoadConfig(file)
}
