package util

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	SmtpEndpoint         string        `mapstructure:"SMTP_ENDPOINT"`
	SmtpUsername         string        `mapstructure:"SMTP_USERNAME"`
	SmtpPassword         string        `mapstructure:"SMTP_PASSWORD"`
	CookieDomain         string        `mapstructure:"COOKIE_DOMAIN"`
	FullDomain           string        `mapstructure:"FULL_DOMAIN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.MergeInConfig()
	if err != nil {
		return
	}

	localConfigPath := path + "/app.env.local"
	if _, err := os.Stat(localConfigPath); err == nil {
		viper.SetConfigName("app.env.local")
		err = viper.MergeInConfig()
		if err != nil && !os.IsNotExist(err) {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
