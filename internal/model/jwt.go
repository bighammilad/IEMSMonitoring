package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name                 string `json:"name,omitempty"`
	RoleId               int    `json:"role,omitempty"`
	jwt.RegisteredClaims `json:"registered_claims,omitempty"`
}
