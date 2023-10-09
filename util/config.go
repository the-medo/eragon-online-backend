package util

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	Environment            string        `mapstructure:"ENVIRONMENT"`
	DBDriver               string        `mapstructure:"DB_DRIVER"`
	DBSource               string        `mapstructure:"DB_SOURCE"`
	MigrationURL           string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress      string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress      string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration   time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddress           string        `mapstructure:"REDIS_ADDRESS"`
	EmailSenderName        string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress     string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	SmtpEndpoint           string        `mapstructure:"SMTP_ENDPOINT"`
	SmtpUsername           string        `mapstructure:"SMTP_USERNAME"`
	SmtpPassword           string        `mapstructure:"SMTP_PASSWORD"`
	CookieDomain           string        `mapstructure:"COOKIE_DOMAIN"`
	FullDomain             string        `mapstructure:"FULL_DOMAIN"`
	CloudflareApiToken     string        `mapstructure:"CLOUDFLARE_API_TOKEN"`
	CloudflareAccountId    string        `mapstructure:"CLOUDFLARE_ACCOUNT_ID"`
	SentryDsn              string        `mapstructure:"SENTRY_DSN"`
	SentryTracesSampleRate float64       `mapstructure:"SENTRY_TRACES_SAMPLE_RATE"`
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
	_, err = os.Stat(localConfigPath)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if !os.IsNotExist(err) {
		viper.SetConfigName("app.env.local")
		err = viper.MergeInConfig()
		if err != nil {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
