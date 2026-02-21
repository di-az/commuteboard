package main

import (
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	DBName string
	DBHost string
	DBPort string
}

func LoadEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	env := &Env{
		DBName: os.Getenv("DB_NAME"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
	}

	return env, nil
}
