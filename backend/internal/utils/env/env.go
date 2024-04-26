package env

import (
	"os"

	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/log"
)

func GetEnvOrFatal(key string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		log.Fatalf("Error getting environment variable. Key: %s", key)
	}
	return value
}

func GetEnvOrDefault(key, defaultValue string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}
	return value
}
