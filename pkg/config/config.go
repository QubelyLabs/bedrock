package config

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Port        string `mapstructure:"PORT"`
	GoEnv       string `mapstructure:"GO_ENV"`

	// Database configuration
	DBDialect  string `mapstructure:"DB_DIALECT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASS"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSync     bool   `mapstructure:"DB_SYNC"`
	DBLog      bool   `mapstructure:"DB_LOG"`

	// Throttle settings
	TTL   int `mapstructure:"THROTTLE_TTL"`
	Limit int `mapstructure:"THROTTLE_LIMIT"`

	// HTTP settings
	Timeout      int `mapstructure:"HTTP_TIMEOUT"`
	MaxRedirects int `mapstructure:"HTTP_MAX_REDIRECTS"`

	// Redis configuration
	RedisDB       int    `mapstructure:"REDIS_DB"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPassword string `mapstructure:"REDIS_PASS"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`

	// Switch settings
	TrafficLog string `mapstructure:"TRAFFIC_LOG_SWITCH"`
	Shutdown   string `mapstructure:"SHUTDOWN_SWITCH"`
}
