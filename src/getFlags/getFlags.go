package main

import (
        "context"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-sdk-go-v2/service/dynamodb"
        "github.com/aws/aws-sdk-go/aws/session"

)

type Response struct {
        statusCode int
}

func HandleRequest(ctx context.Context ) (Response, error) {
        sess, err := session.NewSession()
        response := Response{statusCode: 200}
        return response, nil
}

func main() {
        lambda.Start(HandleRequest)
}