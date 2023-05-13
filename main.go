package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shusann01116/disco-gpt/discord"
)

func Handle(ctx context.Context, event interface{}) (*discord.DicordResponse, error) {
	fmt.Printf("%v\n", event)
	return &discord.DicordResponse{
		Message: "Hello from lambda",
	}, nil
}

func main() {
	lambda.Start(Handle)
}
