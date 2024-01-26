package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name                 string `json:"name,omitempty"`
	Admin                bool   `json:"admin,omitempty"`
	RoleId               int    `json:"role_id,omitempty"`
	jwt.RegisteredClaims `json:"registered_claims,omitempty"`
}
