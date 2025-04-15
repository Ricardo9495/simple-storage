package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		PG          `yaml:"postgres"`
		ApiPort     int64  `env-required:"true" yaml:"api_port" env:"API_PORT"`
		ApiHost     string `env-required:"true" yaml:"api_host" env:"API_HOST"`
		MaxFileSize int64  `env-required:"true" yaml:"max_file_size" env:"MAX_FILE_SIZE"`
		StorageDir  string `env-required:"true" yaml:"storage_dir" env:"STORAGE_DIR"`
	}

	PG struct {
		DBName   string `env-required:"true" yaml:"db_name" env:"PG_DB_NAME"`
		PoolMax  int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		Host     string `env-required:"true" yaml:"host"      env:"PG_HOST"`
		Port     int    `env-required:"true" yaml:"port"      env:"PG_PORT"`
		User     string `env-required:"true" yaml:"user"      env:"PG_USER"`
		Password string `env-required:"true" yaml:"password"      env:"PG_PASSWORD"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	cwd := projectRoot()
	envFilePath := cwd + ".env"

	err := readEnv(envFilePath, cfg)
	if err != nil {
		panic(err)
	}
	cfg.StorageDir = cwd + cfg.StorageDir

	return cfg
}

func readEnv(envFilePath string, cfg *Config) error {
	envFileExists := checkFileExists(envFilePath)

	if envFileExists {
		err := cleanenv.ReadConfig(envFilePath, cfg)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	} else {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {

			if _, statErr := os.Stat(envFilePath + ".example"); statErr == nil {
				return fmt.Errorf("missing environmentvariables: %w\n\nprovide all required environment variables or rename and update .env.example to .env for convinience", err)
			}

			return err
		}
	}
	return nil
}

func checkFileExists(fileName string) bool {
	envFileExists := false
	if _, err := os.Stat(fileName); err == nil {
		envFileExists = true
	}
	return envFileExists
}

func projectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(b)

	return projectRoot + "/../"
}
