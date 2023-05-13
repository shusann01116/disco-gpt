package discord

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"log"
)

type DiscordRequest struct {
	ApplicationID string        `json:"application_id"`
	Entitlements  []interface{} `json:"entitlements"`
	ID            string        `json:"id"`
	Token         string        `json:"token"`
	Type          int           `json:"type"`
	User          User          `json:"user"`
	Version       int           `json:"version"`
}

type User struct {
	Avatar           string      `json:"avatar"`
	AvatarDecoration interface{} `json:"avatar_decoration"`
	Discriminator    string      `json:"discriminator"`
	DisplayName      interface{} `json:"display_name"`
	GlobalName       interface{} `json:"global_name"`
	ID               string      `json:"id"`
	PublicFlags      int         `json:"public_flags"`
	Username         string      `json:"username"`
}

type DicordResponse struct {
	Message string `json:"message"`
}

func VerifyRequest(timestamp, body, signature string, key ed25519.PublicKey) bool {
	msg := bytes.NewBufferString(timestamp + body)

	sig, err := hex.DecodeString(signature)
	if err != nil {
		log.Println("Failed to decode signature:", err)
		return false
	}

	if len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
		log.Println("Invalid signature")
		return false
	}

	return ed25519.Verify(key, msg.Bytes(), sig)
}
