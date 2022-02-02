package configs

import (
	"hackathon-api/utils"
	"os"
)

func EnvMongoURI() string {
	// Config validation
	mongoURI, errMongoURI := os.LookupEnv("MONGOURI")
	utils.Checkb(errMongoURI, "MONGOURI is required")
	return mongoURI
}
