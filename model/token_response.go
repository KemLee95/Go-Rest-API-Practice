package model

type TokenResponse struct {
	ExpireAt int64  `json:"expire_at"`
	Token    string `json:"token"`
}
