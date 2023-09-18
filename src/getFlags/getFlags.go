package main

import (
        "context"
        "fmt"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-sdk-go-v2/aws"
        "github.com/aws/aws-sdk-go-v2/config"
        "github.com/aws/aws-sdk-go-v2/service/dynamodb"
        "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

)

type Response struct {
        statusCode int
}

type Query struct {
    TableName string
    KeyConditionExpression string 
}

func HandleRequest(ctx context.Context ) (Response, error) {
        cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"))
        if err != nil {
            fmt.Println("unable to load SDK config, %v", err)
        }

        tableName := "release-table"
        keyCond := expression.Key("ClientId").Equal(expression.Value("client-1"))
        expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()

        fmt.Print(tableName)

    
        // Using the Config value, create the DynamoDB client
        svc := dynamodb.NewFromConfig(cfg)
    
        // Build the request with its input parameters
        resp, err := svc.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
		})

        if err != nil {
            fmt.Println("Failed to read items, %v", err)
        }
    
        for _, flag := range resp.Items {
            fmt.Println(flag)
        }
        fmt.Print("Complete")
        response := Response{statusCode: 200}
        return response, nil
}

func main() {
        lambda.Start(HandleRequest)
}