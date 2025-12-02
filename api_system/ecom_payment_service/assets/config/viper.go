package config_assets

// write a struct and a function to read the .env using viper

import (
	"time"

	"github.com/spf13/viper"
)

type ReadENV struct {
	DBSource          string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	JWTSecret         string        `mapstructure:"JWT_SECRET"`
	OrderDuration     time.Duration `mapstructure:"ORDER_DURATION"`
	ClientIP          []string      `mapstructure:"CLIENT_IP"`

	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	// Kafka configuration
	KafkaBrokers       string `mapstructure:"KAFKA_BROKERS"`
	KafkaConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP"`

	AccessKeyMoMo string `mapstructure:"ACCESS_KEY_MOMO"`
	SecretKeyMoMo string `mapstructure:"SECRET_KEY_MOMO"`
	RedirectURL   string `mapstructure:"REDIRECTURL"`
	PublicID      string `mapstructure:"PUBLIC_ID"`
	IpnURL        string `mapstructure:"IPNURL"`
	EndPointMoMo  string `mapstructure:"ENDPOINT_MOMO"`
	PlatformID    string `mapstructure:"PLATFORM_ID"`

	// URL service
	URLProductService string `mapstructure:"URL_PRODUCT_SERVICE"`
	URLOrderService   string `mapstructure:"URL_ORDER_SERVICE"`

	// email service
	BrevoAPIKey string `mapstructure:"BREVO_API_KEY"`
	SenderEmail string `mapstructure:"SENDER_EMAIL"`
	SenderName  string `mapstructure:"SENDER_NAME"`
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
