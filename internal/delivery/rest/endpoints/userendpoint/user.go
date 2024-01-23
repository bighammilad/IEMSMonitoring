package userendpoint

import (
	"monitoring/internal/usecase/useruc"

	"github.com/golang-jwt/jwt/v4"
)

type UserEndpoint struct {
	Name  string `json:"name,omitempty"`
	Admin bool   `json:"admin,omitempty"`
	jwt.RegisteredClaims
	UserUC *useruc.UserUsecase `json:"repo,omitempty"`
}
