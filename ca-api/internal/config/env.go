package config

import "github.com/joho/godotenv"

func LoadEnvVariables() {
	godotenv.Load()
}
