package model

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Configuration model
type Configuration struct {
	HttpPort        string
	KafkaConnection string
}

// GetConfiguration method basically populate configuration information from .env and return Configuration model
func GetConfiguration() Configuration {
	err := godotenv.Load("./.env")

	if err != nil {
		logrus.Error("Error loading .env file")
	}

	configuration := Configuration{
		os.Getenv("HTTP_PORT"),
		os.Getenv("KAFKA_CONNECTION"),
	}

	return configuration
}
