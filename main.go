package main

import (
	"context"
	"github.com/HoseonYim/isfree-backend/database"
	"github.com/HoseonYim/isfree-backend/router"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

var (
	fiberLambda *fiberadapter.FiberLambda
)

// init the Fiber Server
func init() {
	log.Printf("Fiber cold start")

	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return
	}
	database.DB = dynamodb.New(awsSession)

	app := fiber.New()
	router.Initialize(app)
	fiberLambda = fiberadapter.New(app)
}

// Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
