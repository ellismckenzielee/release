package main

import (
        "context"
        "fmt"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-sdk-go-v2/aws"
        "github.com/aws/aws-sdk-go-v2/config"
        "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Response struct {
        statusCode int
}

func HandleRequest(ctx context.Context ) (Response, error) {
        cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"))
        if err != nil {
            fmt.Println("unable to load SDK config, %v", err)
        }
    
        // Using the Config value, create the DynamoDB client
        svc := dynamodb.NewFromConfig(cfg)
    
        // Build the request with its input parameters
        resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
            Limit: aws.Int32(5),
        })
        if err != nil {
            fmt.Println("Failed to list tables, %v", err)
        }
    
        fmt.Println("Tables:")
        for _, tableName := range resp.TableNames {
            fmt.Println(tableName)
        }
        response := Response{statusCode: 200}
        return response, nil
}

func main() {
        lambda.Start(HandleRequest)
}