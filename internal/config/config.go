package config

import (
	"os"
	"strings"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBScheme   string
	Port       string
}

func LoadEnv(logger *utils.Logger) error {
	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal("Error getting current working directory: %v", err)
		return err
	}

	//Крайне кривая реализация, но работает
	envFilePath := wd + "/../../internal/config/.env"
	envFilePath = strings.Replace(envFilePath, "\\", "/", -1)

	logger.Debugf("Loading .env file from path: %s", envFilePath)

	err = godotenv.Load(envFilePath)
	if err != nil {
		logger.Fatal("Error loading .env file: %v", err)
		return err
	}

	return nil
}

func GetConfig() *Config {
	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBScheme:   os.Getenv("DB_SCHEME"),
		Port:       os.Getenv("PORT"),
	}
}
