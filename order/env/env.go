package env

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Env struct {
	HTTPAddress         string
	HTTPTimeout         time.Duration
	RabbitMQ            RabbitMQ
	DBConnectionDSN     string
	DBConnectionMigrate string
}

type RabbitMQ struct {
	Host         string
	Port         int
	User         string
	Password     string
	ConsumerName string
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
	dbConn := dbConnection{}

	err := errors.Join(
		parseKey("HTTP_ADDRESS", &e.HTTPAddress),
		parseDuration("HTTP_TIMEOUT", &e.HTTPTimeout),
		parseKey("RABBITMQ_HOST", &e.RabbitMQ.Host),
		parseInt("RABBITMQ_PORT", &e.RabbitMQ.Port),
		parseKey("RABBITMQ_USER", &e.RabbitMQ.User),
		parseKey("RABBITMQ_PASSWORD", &e.RabbitMQ.Password),
		parseKey("RABBITMQ_CONSUMER_NAME", &e.RabbitMQ.ConsumerName),
		parseKey("DB_HOST", &dbConn.host),
		parseKey("DB_USER", &dbConn.user),
		parseKey("DB_PASSWORD", &dbConn.password),
		parseKey("DB_PORT", &dbConn.port),
		parseKey("DB_NAME", &dbConn.name),
	)

	e.DBConnectionMigrate = dbConn.getConn("pgx5")
	e.DBConnectionDSN = dbConn.getConn("postgresql")
	return e, err
}

func parseKey(envKey string, field *string) error {
	val := os.Getenv(envKey)
	if val == "" {
		return fmt.Errorf("%s env var is not set", envKey)
	}

	*field = val
	return nil
}

func parseInt(envKey string, field *int) error {
	value := os.Getenv(envKey)
	if value == "" {
		return fmt.Errorf("%s env var is not set", envKey)
	}

	parsed, err := strconv.Atoi(value)
	*field = parsed
	return err
}

func parseDuration(envKey string, field *time.Duration) error {
	durationStr := ""
	if err := parseKey(envKey, &durationStr); err != nil {
		return err
	}

	var err error
	*field, err = time.ParseDuration(durationStr)
	return err
}

type dbConnection struct {
	host     string
	user     string
	password string
	port     string
	name     string
}

func (db dbConnection) getConn(driver string) string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		driver,
		db.user,
		url.QueryEscape(db.password),
		db.host,
		db.port,
		db.name,
	)
}
