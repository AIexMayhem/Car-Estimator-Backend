package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                   string
	Port                  string
	ProfileServiceAddr    string
	PredictionServiceAddr string
	FeedServiceAddr       string
}

func Load(envfile string) (*Config, error) {
	if err := godotenv.Load(envfile); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	return &Config{
		Env:                   os.Getenv("MODE"),
		Port:                  os.Getenv("SERVE_PORT"),
		ProfileServiceAddr: os.Getenv("PROFILE_SERVICE_ADDR"),
		PredictionServiceAddr: os.Getenv("PREDICTION_SERVICE_ADDR"),
		FeedServiceAddr: os.Getenv("FEED_SERVICE_ADDR"),
	}, nil
}
