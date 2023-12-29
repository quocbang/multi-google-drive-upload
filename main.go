package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	APIKey string `envconfig:"API_KEY"`
}

func loadEnv() (*Env, error) {
	// load default .env file, ignore the error
	_ = godotenv.Load()

	env := new(Env)
	err := envconfig.Process("", env)
	if err != nil {
		return nil, fmt.Errorf("load config error: %v", err)
	}

	return env, nil
}

func main() {
	_, err := loadEnv()
	if err != nil {
		log.Fatalf("failed to load env, error: %v", err)
	}
}
