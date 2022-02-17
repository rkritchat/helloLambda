package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"helloLambda/internal/user"
	"net/http"
)

var service user.Service

func main() {
	service = user.NewService()
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQ: %#v\n", req)
	fmt.Printf("CTX: %#v\n", ctx)
	switch req.HTTPMethod {
	case http.MethodGet:
		return service.GetUser(req)
	case http.MethodPost:
		return service.CreateUser(req)
	case http.MethodPut:
		return service.JustReturnErr(req)
	default:
		fmt.Printf("method is not allowed, method:%v", req.HTTPMethod)
		return service.MethodNotAllowed()
	}
}
