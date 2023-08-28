package main

import (
        "context"
        "github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
        statusCode int
}

func HandleRequest(ctx context.Context ) (Response, error) {
        response := Response{statusCode: 200}
        return response, nil
}

func main() {
        lambda.Start(HandleRequest)
}