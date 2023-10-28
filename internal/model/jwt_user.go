package model

import "github.com/golang-jwt/jwt/v5"

const (
	AppGuardName = "app"
)

type JwtUser interface {
	GetUid() string
}

// CustomClaims 自定义 Claims
type CustomClaims struct {
	Key string `json:"key,omitempty"`
	jwt.RegisteredClaims
}

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
