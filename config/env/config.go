package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const (
	redis_database_url   = "REDIS_URL"
	redis_host           = "REDIS_HOST"
	redis_port           = "REDIS_PORT"
	redis_container_name = "REDIS_CONTAINER_NAME"

	mysql_database_url   = "MYSQL_URL"
	mysql_host           = "MYSQL_HOST"
	mysql_port           = "MYSQL_PORT"
	mysql_db_name        = "MySQL_DB_NAME"
	mysql_user           = "MYSQL_USER"
	mysql_password       = "MYSQL_PASSWORD"
	mysql_container_name = "MYSQL_CONTAINER_NAME"
)

func init() {
	if err := godotenv.Load("./config/env/.env"); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}
}

type RedisConfig struct {
	RedisUrl       string
	Host           string
	Port           string
	RedisContainer string
}

type MysqlConfig struct {
	MysqlUrl       string
	Host           string
	Port           string
	User           string
	DbName         string
	Password       string
	MysqlContainer string
}

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		RedisUrl:       os.Getenv(redis_database_url),
		Host:           os.Getenv(redis_host),
		Port:           os.Getenv(redis_port),
		RedisContainer: os.Getenv(redis_container_name),
	}
}

func GetMysqlConfig() *MysqlConfig {
	return &MysqlConfig{
		MysqlUrl:       os.Getenv(mysql_database_url),
		Host:           os.Getenv(mysql_host),
		Port:           os.Getenv(mysql_port),
		DbName:         os.Getenv(mysql_db_name),
		User:           os.Getenv(mysql_user),
		Password:       os.Getenv(mysql_password),
		MysqlContainer: os.Getenv(mysql_container_name),
	}
}
