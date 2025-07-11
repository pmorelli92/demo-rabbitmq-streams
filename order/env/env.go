package env

import (
	"errors"
	"os"
	"strconv"
)

type Env struct {
	HTTPServerAddress string
	RabbitMQHost      string
	RabbitMQPort      int
	RabbitMQUser      string
	RabbitMQPassword  string
	ConsumerName      string
	DB                DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func ParseEnv() (Env, error) {
	e := Env{}
	err := errors.Join(
		parseKey("HTTP_SERVER_ADDRESS", &e.HTTPServerAddress),
		parseKey("RABBITMQ_HOST", &e.RabbitMQHost),
		parseInt("RABBITMQ_PORT", &e.RabbitMQPort),
		parseKey("RABBITMQ_USER", &e.RabbitMQUser),
		parseKey("RABBITMQ_PASSWORD", &e.RabbitMQPassword),
		parseKey("CONSUMER_NAME", &e.ConsumerName),
		parseKey("DB_HOST", &e.DB.Host),
		parseInt("DB_PORT", &e.DB.Port),
		parseKey("DB_USER", &e.DB.User),
		parseKey("DB_PASSWORD", &e.DB.Password),
		parseKey("DB_NAME", &e.DB.Database),
		parseKey("DB_SSL_MODE", &e.DB.SSLMode),
	)
	return e, err
}

func parseKey(key string, target *string) error {
	value := os.Getenv(key)
	if value == "" {
		return errors.New("missing required environment variable: " + key)
	}
	*target = value
	return nil
}

func parseInt(key string, target *int) error {
	value := os.Getenv(key)
	if value == "" {
		return errors.New("missing required environment variable: " + key)
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("invalid integer value for " + key + ": " + value)
	}
	*target = parsed
	return nil
}
