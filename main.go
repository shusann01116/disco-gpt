package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shusann01116/disco-gpt/discord"
)

func Handle(ctx context.Context, event interface{}) (*discord.DicordResponse, error) {
	fmt.Printf("%v\n", event)

	e, ok := event.(map[string]interface{})
	if !ok {
		log.Println("Error: failed to parse event")
		return nil, fmt.Errorf("failed to parse event")
	}

	log.Println(e["body"])

	var req discord.DiscordRequest
	if err := json.Unmarshal([]byte(e["body"].(string)), &req); err != nil {
		log.Println("Error: failed to parse request body")
		return nil, fmt.Errorf("failed to parse request body")
	}

	log.Printf("ID: %v, User: %v\n", req.ID, req.User)

	return &discord.DicordResponse{
		Message: "Hello from lambda",
	}, nil
}

func main() {
	lambda.Start(Handle)
}
