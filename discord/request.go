package discord

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
