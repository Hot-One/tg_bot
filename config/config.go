package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	BotToken string

	PostgresHost          string
	PostgresUser          string
	PostgresDatabase      string
	PostgresPassword      string
	PostgresPort          int
	PostgresMaxConnection int32
}

func Load() Config {
	cfg := Config{}

	cfg.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "localhost"))
	cfg.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "abdulbosit"))
	cfg.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "bot"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "946236953"))
	cfg.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432))

	cfg.PostgresMaxConnection = cast.ToInt32(getOrReturnDefaultValue("POSTGRES_PORT", 30))

	cfg.BotToken = cast.ToString(getOrReturnDefaultValue("TELEGRAM_BOT_TOKEN", "6680687799:AAEG-kVUEwufsiSPTv47j9kuAvQLmKk6iOI"))

	return cfg

}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
