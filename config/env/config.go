package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	TEMPORAL_HOSTPORT  = "TEMPORAL_HOSTPORT"
	TEMPORAL_NAMESPACE = "TEMPORAL_NAMESPACE"
	TEMPORAL_TASKQUEUE = "TEMPORAL_TASKQUEUE"

	TERRAFORM_STATE_DIR = "TERRAFORM_STATE_DIR"

	TERRAFORM_STATE_BUCKET = "TERRAFORM_STATE_BUCKET"

	GITHUB_TOKEN = "GITHUB_TOKEN"
)

type (
	Config struct {
		TemporalHostPort  string
		TemporalNamespace string
		TemporalTaskQueue string
		TfState           TfState
		GithubToken       string
	}

	TfState struct {
		Region   string // aws region
		Bucket   string // s3 bucket name
		DynamoDB string
	}
)

func MustGetConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		TemporalHostPort:  getOrDefault(TEMPORAL_HOSTPORT, "127.0.0.1:7233"),
		TemporalNamespace: getOrDefault(TEMPORAL_NAMESPACE, "default"),
		TemporalTaskQueue: getOrDefault(TEMPORAL_TASKQUEUE, "meetupgo"),
		TfState:           getTfState(),
		GithubToken:       getOrDefault(GITHUB_TOKEN, ""),
	}
}

func getOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getTfState() TfState {
	return TfState{
		Bucket: getOrDefault(TERRAFORM_STATE_BUCKET, "meetupgo"),
	}
}
