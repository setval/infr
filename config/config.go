package config

import (
	"context"
	"os"

	"github.com/sethvargo/go-envconfig"
	"github.com/subosito/gotenv"
)

func New[T any]() (*T, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}
	return parseEnv[T]()
}

func parseEnv[T any]() (*T, error) {
	var c T
	ctx := context.Background()
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func loadEnv() error {
	if _, err := os.Stat(".env.local"); err == nil {
		if err := gotenv.Load(".env.local"); err != nil {
			return err
		}
	} else {
		if _, err := os.Stat(".env"); err == nil {
			if err := gotenv.Load(".env"); err != nil {
				return err
			}
		}
	}
	return nil
}
