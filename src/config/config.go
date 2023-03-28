package config

import (
	"github.com/spf13/viper"
)

type databaseConfig struct {
	Uri  string
	Name string
}

type apiConfig struct {
	Port     string
	Database databaseConfig
}

var (
	config *apiConfig
)

func GetAPIConfig() (*apiConfig, error) {
	if config != nil {
		return config, nil
	}
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config = &apiConfig{
		Port: viper.GetString("port"),
		Database: databaseConfig{
			Uri:  viper.GetString("mongo.uri"),
			Name: viper.GetString("mongo.name"),
		},
	}
	return config, nil
}
