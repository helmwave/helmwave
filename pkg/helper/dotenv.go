package helper

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
)

func Dotenv() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}
	}
	Helm = helm.New() // Recreate helm instance to respect helm variables from .env file
}
