package config

import (
	"github.com/spf13/viper"
)

type ApiConfig struct {
	Port             string
	Database         databaseConfig
	DefaultLanguages string
	Cors             corsConfig
	JWTConfig        JWTConfig
}
type JWTConfig struct {
	PrivateKey             string
	AccessTokenExpiration  int
	RefreshTokenExpiration int
}
type databaseConfig struct {
	Uri  string
	Name string
}
type corsConfig struct {
	AllowOrigin []string
	AllowHeader []string
	AllowMethod []string
}

var (
	config *ApiConfig
)

func GetAPIConfig() (*ApiConfig, error) {
	if config != nil {
		return config, nil
	}
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	jwtConfig := JWTConfig{
		PrivateKey:             viper.GetString("jwt.private_key"),
		AccessTokenExpiration:  viper.GetInt("jwt.access_token_expiration"),
		RefreshTokenExpiration: viper.GetInt("jwt.refresh_token_expiration"),
	}
	config = &ApiConfig{
		Port: viper.GetString("port"),
		Database: databaseConfig{
			Uri:  viper.GetString("mongo.uri"),
			Name: viper.GetString("mongo.name"),
		},
		DefaultLanguages: viper.GetString("defaultLanguages"),
		Cors: corsConfig{
			AllowOrigin: []string{"*"},
			AllowHeader: []string{
				"Accept",
				"Accept-Language",
				"Content-Language",
				"Content-Type",
				"Origin",
				"Authorization",
				"Access-Control-Request-Method",
				"Access-Control-Request-Headers",
				"Access-Control-Allow-Headers",
				"Access-Control-Allow-Origin",
				"Access-Control-Allow-Methods",
				"Access-Control-Allow-Credentials",
				"Access-Control-Expose-Headers",
				"Access-Control-Max-Age",
				"Referer",
				"Host",
				"x-language-code",
				"x-timestamp",
				"x-timezone",
				"x-request-id",
				"user-agent",
			},
			AllowMethod: []string{
				"GET",
				"POST",
				"PUT",
				"DELETE",
			},
		},
		JWTConfig: jwtConfig,
	}
	return config, nil
}
