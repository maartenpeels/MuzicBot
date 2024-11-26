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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalf("No token provided")
	}

	Env = &env{
		Token: token,
	}

	return nil
}
