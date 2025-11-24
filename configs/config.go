package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Email    EmailConfig    `mapstructure:"email"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	Timeout int    `mapstructure:"timeout"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type EmailConfig struct {
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("server.port", ":8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)

	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Override with environment variables if they exist
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = ":" + port
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}

	return &config, nil
}

// GetDSN returns database connection string
func (d *DatabaseConfig) GetDSN() string {
	return "host=" + d.Host +
		" port=" + strconv.Itoa(d.Port) +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Name +
		" sslmode=" + d.SSLMode
}
