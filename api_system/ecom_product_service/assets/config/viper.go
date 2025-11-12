package config_assets

// write a struct and a function to read the .env using viper

import (
	"github.com/spf13/viper"
)

type ReadENV struct {
	DBSource          string   `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string   `mapstructure:"HTTP_SERVER_ADDRESS"`
	JWTSecret         string   `mapstructure:"JWT_SECRET"`
	ClientIP          []string `mapstructure:"CLIENT_IP"`
	RedisAddress      string   `mapstructure:"REDIS_ADDRESS"`
	// // URL service
	ImagePath string `mapstructure:"IMAGE_PATH"`
}

func LoadConfig(path string) (config ReadENV, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
