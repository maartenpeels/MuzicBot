package bot

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type env struct {
	Token string
}

var Env *env

func LoadEnv() error {
	// Attempt to load .env file but do not fail if it's missing so we can override with env variable
	_ = godotenv.Load(".env")

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalf("No token provided")
	}

	Env = &env{
		Token: token,
	}

	return nil
}
