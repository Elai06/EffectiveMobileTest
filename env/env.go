package env

import (
	"effectiveMobile/internal/utils"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	HttpPort     string        `env:"HANDLER_PORT" default:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"10s"`
	DBPort       string        `env:"DB_PORT" default:""`
	DBHost       string        `env:"DB_HOST" default:""`
	DBName       string        `env:"DB_NAME" default:""`
	DBUser       string        `env:"DB_USER" default:""`
	DBPass       string        `env:"DB_PASSWORD" default:""`
	DBSSlMode    string        `env:"DB_SSLMODE" default:""`
	DebugMode    bool          `env:"DEBUG_MODE" default:"false"`
}

func LoadConfig() (Config, error) {
	utils.InfoLog("Loading env variables")
	cfg := Config{}
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return cfg, err
	}

	utils.DebugLog("Config is loaded")

	return cfg, nil
}
