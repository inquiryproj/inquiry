package factory

import (
	"time"

	"github.com/caarlos0/env/v9"
)

// LogLevel is the level of the logs.
type LogLevel string

// Log levels.
const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogFormat is the format of the logs.
type LogFormat string

// Log formats.
const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

// Config is the configuration for the application.
type Config struct {
	LogLevel  LogLevel  `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat LogFormat `env:"LOG_FORMAT" envDefault:"json"`

	RepositoryConfig RepositoryConfig
	ServerConfig     ServerConfig
}

// RepositoryType is the type of repository.
type RepositoryType string

// String returns the string representation of the repository type.
func (r RepositoryType) String() string {
	return string(r)
}

// Repository types.
const (
	RepositoryTypeSQLite RepositoryType = "sqlite"
)

// RepositoryConfig is the configuration for the repository.
type RepositoryConfig struct {
	RepositoryType RepositoryType `env:"REPOSITORY_TYPE" envDefault:"sqlite"`
	DSN            string         `env:"REPOSITORY_DSN" envDefault:"inquiry.db"`
}

// ServerConfig is the configuration for the server.
type ServerConfig struct {
	Port          int           `env:"API_PORT" envDefault:"3000"`
	ShutdownDelay time.Duration `env:"API_SHUTDOWN_DELAY" envDefault:"0s"`
	AuthEnabled   bool          `env:"API_AUTH_ENABLED" envDefault:"true"`
	APIKey        string        `env:"API_KEY" envDefault:""`
}

// NewConfig creates a new Config instance.
func NewConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
