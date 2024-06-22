package testconfig

import (
	"os"
)

type Config struct {
	CheckMigratedTests bool
}

func NewConfig() *Config {
	return &Config{
		CheckMigratedTests: os.Getenv("CHECK_MIGRATED_TESTS") == "true",
	}
}
