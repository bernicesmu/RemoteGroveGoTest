package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	// initialising dynamodb
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_REGION")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Read AWS credentials and region from environment variables
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	// Configure AWS SDK using the loaded credentials and region
	// cfg is an instance of aws.Config
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Error configuring AWS SDK: %v", err)
	}

	dynamodb.NewFromConfig(cfg)

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
	}
	response, err := GetAllCartsHandler(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	// allCartsJSON, err := json.Marshal(AllCarts)
	// if err != nil {
	// 	t.Fatalf("Failed to marshal AllCarts: %v", err)
	// }
	// assert.Equal(t, string(allCartsJSON), response.Body)
}
