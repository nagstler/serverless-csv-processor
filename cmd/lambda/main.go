package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nagstler/serverless-csv-processor/internal/handler"
)

func main() {
	lambda.Start(handler.HandleS3Event)
}
