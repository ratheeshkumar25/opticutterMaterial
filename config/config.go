package config

import "github.com/spf13/viper"

type Config struct {
	Host            string `mapstructure:"HOST"`
	User            string `mapstructure:"USER"`
	Password        string `mapstructure:"PASSWORD"`
	Database        string `mapstructure:"DBNAME"`
	Port            string `mapstructure:"PORT"`
	Sslmode         string `mapstructure:"SSL"`
	GrpcPort        string `mapstructure:"GRPCPORT"`
	UserPort        string `mapstructure:"GRPCUSERPORT"`
	AdminPort       string `mapstructure:"GRPCADMINPORT"`
	REDISHOST       string `mapstructure:"REDISHOST"`
	APIKey          string `mapstructure:"RAZORPAY_KEY_ID"`
	APISecret       string `mapstructure:"RAZORPAY_SECRET"`
	STRIPESECRETKEY string `mapstructure:"STRIPE_SECRET_KEY"`
}

// LoadConfig will load the environment variable to access.
func LoadConfig() *Config {
	var config Config
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.Unmarshal(&config)
	return &config
}
