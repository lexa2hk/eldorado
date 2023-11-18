package config

import (
	"errors"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type APIConfig struct {
	Env             string   `yaml:"env"`
	AuthServiceAddr string   `yaml:"auth_service_addr" env:"AUTH_SERVICE_ADDR"`
	Postgres        postgres `yaml:"postgres"`
	Server          server   `yaml:"server"`
	Redis           redis    `yaml:"redis"`
}

type StatisticServiceConfig struct {
	Env      string   `yaml:"env"`
	Schedule string   `yaml:"schedule"`
	RabbitMQ rabbitmq `yaml:"rabbitmq"`
	Postgres postgres `yaml:"postgres"`
}

type EmailServiceConfig struct {
	Env      string   `yaml:"env"`
	RabbitMQ rabbitmq `yaml:"rabbitmq"`
	SMTP     smtp     `yaml:"smtp"`
}

type AuthServiceConfig struct {
	Env          string   `yaml:"env"`
	Port         string   `env:"PORT"`
	AccessCreds  rsacreds `yaml:"access"`
	RefreshCreds rsacreds `yaml:"refresh"`
	Postgres     postgres `yaml:"postgres"`
	Redis        redis    `yaml:"redis"`
}

type rsacreds struct {
	PrivateKey string        `yaml:"private_key" env:"RSA_PRIVATE_KEY"`
	PublicKey  string        `yaml:"public_key" env:"RSA_PUBLIC_KEY"`
	Expires    time.Duration `yaml:"expires"`
}

type server struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type redis struct {
	URL string        `yaml:"url" env:"REDIS_URL"`
	TTL time.Duration `yaml:"cache_ttl" env:"REDIS_CACHE_TTL" env-default:"30m"`
}

type postgres struct {
	URL string `yaml:"url" env:"PG_URL"`
}

type rabbitmq struct {
	URL       string `yaml:"url" env:"RABBITMQ_URL"`
	QueueName string `yaml:"queue_name"`
}

type smtp struct {
	Host     string `yaml:"host" env:"SMTP_HOST"`
	Port     string `yaml:"port" env:"SMTP_PORT"`
	Email    string `env:"SMTP_EMAIL"`
	Password string `env:"SMTP_PASS"`
}

func LoadEmailSenderConfig() (*EmailServiceConfig, error) {
	path, err := configPath("EMAIL_SENDER_CONFIG_PATH")
	if err != nil {
		return nil, err
	}

	var cfg EmailServiceConfig
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadStatisticServiceConfig() (*StatisticServiceConfig, error) {
	path, err := configPath("STATISTIC_CONFIG_PATH")
	if err != nil {
		return nil, err
	}

	var cfg StatisticServiceConfig
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, err
}

func LoadAuthServiceConfig() (*AuthServiceConfig, error) {
	path, err := configPath("AUTH_SERVICE_CONFIG_PATH")
	if err != nil {
		return nil, err
	}

	var cfg AuthServiceConfig
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, err
}

func LoadApiConfig() (*APIConfig, error) {
	path, err := configPath("API_CONFIG_PATH")
	if err != nil {
		return nil, err
	}

	var cfg APIConfig
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func configPath(env string) (string, error) {
	path := os.Getenv(env)
	if path == "" {
		return "", errors.New("path to config file not set")
	}

	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	return path, nil
}
