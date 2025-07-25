package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	DB   DBConfig   `yaml:"db"`
	JWT  JWTConfig  `yaml:"jwt"`
	GRPC GRPCConfig `yaml:"grpc"`
}

type DBConfig struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	BusyTimeout string `yaml:"busy_timeout"`
	JournalMode string `yaml:"journal_mode"`
}

type JWTConfig struct {
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	// Path to Config
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	// Check path is exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	// Parse config file
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
