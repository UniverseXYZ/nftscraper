package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Returns a single config parameter value.
func Get(parameter string) string {

	// the idea here is to change this line to ".env" on deployment to prod
	//and it has to be discussed with @Zero
	//is it normal in go that it's reading the file every time to get the param value?
	err := godotenv.Load("./config/.env.dev")

	if err != nil {
		log.Fatalf("Error loading config! " + err.Error())
	}

	return os.Getenv(parameter)
}