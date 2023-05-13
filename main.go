package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shusann01116/disco-gpt/discord"
)

func ensureValue(key string) string {
	if value, available := os.LookupEnv(key); available {
		return value
	}

	panic("Missing environment variable: " + key)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("RequestID:", req.RequestContext.RequestID)
	log.Println("Body size:", len(req.Body))

	var discordReq discord.DiscordRequest
	if err := json.Unmarshal([]byte(req.Body), &discordReq); err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusBadRequest,
		}, err
	}

	timestamp, signature := req.Headers["x-signature-timestamp"], req.Headers["x-signature-ed25519"]
	log.Println("Timestamp:", timestamp)
	log.Println("Signature:", signature)

	log.Println("Body:", req.Body)

	// Verify request
	key, err := hex.DecodeString(ensureValue("DISCORD_PUBLIC_KEY"))
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Failed to decode public key. Please ensure DISCORD_PUBLIC_KEY is a valid hex string of length 64")
	}

	if !discord.VerifyRequest(timestamp, req.Body, signature, ed25519.PublicKey(key)) {
		// return 403
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusForbidden,
		}, fmt.Errorf("Failed to verify request")
	}

	log.Println("Request verified")

	var res events.LambdaFunctionURLResponse
	switch discordReq.Type {
	case 1:
		res = events.LambdaFunctionURLResponse{
			StatusCode: http.StatusOK,
			Body:       "{\"type\":1}",
		}
		log.Println("Response:", res)
		return res, nil
	}

	// Return invalid request
	res = events.LambdaFunctionURLResponse{
		StatusCode: http.StatusBadRequest,
	}
	log.Println("Response:", res)
	return res, fmt.Errorf("Invalid request")
}

func main() {
	lambda.Start(handler)
}
