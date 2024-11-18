package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct{
	Env string `yaml:"env" env:"ENV" envDefault:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"10s"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("CANNOT READ ENV FILE")
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatalf("ERROR CONF PATH %s IS NOT SET", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil{
		log.Fatalf("Config: %s read faild. %s", configPath, err)
	}


	return &cfg
}