package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	ACCOUNTSID string `mapstructure:"ACCOUNTSID"`
	SERVICESID string `mapstructure:"SERVICESID"`
	AUTHTOKEN  string `mapstructure:"AUTHTOKEN"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", "ACCOUNTS_ID", "SERVICES_ID", "AUTH_TOKEN",
}

func LoadConfig() (Config, error) {
	var config Config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading the env file")
	}
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {

			return config, err
		}
	}

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("error loading the env file")
	// }

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
